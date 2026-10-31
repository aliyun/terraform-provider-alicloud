package coverage

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/splitlayout"
)

// Registry reads constructor identities, not an inferred resource name from a filename.
func Registry(src []byte) (map[string]string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "provider.go", src, 0)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	var parseErr error
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok || ident(kv.Key) != "ResourcesMap" {
			return true
		}
		if found {
			parseErr = fmt.Errorf("multiple ResourcesMap declarations")
			return false
		}
		found = true
		entries := fields(kv.Value)
		if entries == nil {
			parseErr = fmt.Errorf("ResourcesMap must be a literal map")
			return false
		}
		for name, expr := range entries {
			call, ok := expr.(*ast.CallExpr)
			if !ok || ident(call.Fun) == "" || len(call.Args) != 0 {
				parseErr = fmt.Errorf("unresolved resource constructor for %s", name)
				continue
			}
			result[name] = ident(call.Fun)
		}
		return false
	})
	if parseErr != nil {
		return nil, parseErr
	}
	if !found {
		return nil, fmt.Errorf("ResourcesMap is missing")
	}
	return result, nil
}

// RegistryUnion merges the provider.go ResourcesMap literal with the SDKv2
// registrations declared by the split service packages. Split constructors
// carry a product qualifier ("cloudauth.ResourceFaceConfig") because
// top-level names are only unique per package; the flat literal keeps bare
// names. A declaration appearing in both sources is the same startup error
// assembleProvider reports.
func RegistryUnion(providerSrc []byte, registers map[string][]byte) (map[string]string, error) {
	reg, err := Registry(providerSrc)
	if err != nil {
		return nil, err
	}
	for path, src := range registers {
		product, ok := splitlayout.IsServiceRegister(path)
		if !ok {
			continue
		}
		for _, entry := range splitlayout.SDKRegisterEntries(src) {
			if entry.Kind != splitlayout.KindResource {
				continue
			}
			qualified := product + "." + entry.Factory
			if previous, dup := reg[entry.TypeName]; dup {
				return nil, fmt.Errorf("duplicate resource declaration for %s: %s and %s", entry.TypeName, previous, qualified)
			}
			reg[entry.TypeName] = qualified
		}
	}
	return reg, nil
}

// qualify keys a file's top-level names the way the registry references
// them: bare in the flat alicloud package, product-qualified in the split
// service packages, where names are only unique per package.
func qualify(path string, names map[string]bool) map[string]bool {
	product, ok := splitlayout.ServiceProduct(path)
	if !ok {
		return names
	}
	out := make(map[string]bool, len(names))
	for name := range names {
		out[product+"."+name] = true
	}
	return out
}

// Select maps both changed implementations and tests to current registrations.
// A deleted flat implementation is ignored when its constructors were never
// registered (a dead factory's file moved with the product split) or are still
// registered under moved type names; it fails closed when a constructor was
// registered and its type name is no longer registered, distinguishing a real
// registration deletion from an unregistered relocation.
func Select(changed []string, before, after map[string]string, sources map[string][]byte, oldSource func(string) ([]byte, error)) (map[string]string, error) {
	functions := map[string]string{}
	fileFunctions := map[string]map[string]bool{}
	for path, src := range sources {
		names, err := functionNames(path, src)
		if err != nil {
			return nil, err
		}
		names = qualify(path, names)
		fileFunctions[path] = names
		for name := range names {
			if previous, ok := functions[name]; ok {
				return nil, fmt.Errorf("duplicate function %s in %s and %s", name, previous, path)
			}
			functions[name] = path
		}
	}
	selected := map[string]string{}
	add := func(name string) error {
		path := functions[after[name]]
		if path == "" {
			return fmt.Errorf("cannot locate implementation for registered resource %s (%s)", name, after[name])
		}
		selected[name] = strings.TrimSuffix(path, ".go") + "_test.go"
		return nil
	}
	for _, path := range changed {
		if path == "alicloud/provider.go" {
			for name, constructor := range after {
				if before[name] != constructor {
					if err := add(name); err != nil {
						return nil, err
					}
				}
			}
			continue
		}
		// A product register.go is a registration source the same way
		// provider.go is: select every registration the diff touched. The
		// union registry keeps split constructors product-qualified, so
		// untouched products contribute nothing.
		if _, ok := splitlayout.IsServiceRegister(path); ok {
			for name, constructor := range after {
				if before[name] != constructor {
					if err := add(name); err != nil {
						return nil, err
					}
				}
			}
			continue
		}
		flat := strings.HasPrefix(path, "alicloud/resource_alicloud_") && strings.HasSuffix(path, ".go")
		_, stem, _, isSplit := splitlayout.IsSplitSource(path)
		isSplit = isSplit && strings.HasPrefix(stem, "resource_")
		if !flat && !isSplit {
			// Data sources, docs and framework-only service files stay
			// outside the managed-resource scope, as they were in the flat
			// layout.
			continue
		}
		implementation := strings.TrimSuffix(path, "_test.go")
		if implementation != path {
			implementation += ".go"
		}
		names, exists := fileFunctions[implementation]
		if !exists {
			src, err := oldSource(implementation)
			if err != nil {
				return nil, fmt.Errorf("cannot resolve changed resource %s: %w", path, err)
			}
			base, err := functionNames(implementation, src)
			if err != nil {
				return nil, err
			}
			names = qualify(implementation, base)
		}
		matched := false
		for name, constructor := range after {
			if names[constructor] {
				matched = true
				if err := add(name); err != nil {
					return nil, err
				}
			}
		}
		if !matched && exists && flat {
			return nil, fmt.Errorf("changed resource file %s has no current ResourcesMap registration", path)
		}
		if !matched && !exists && flat {
			// The flat implementation file is gone from the current snapshot.
			// Its constructors are bare, pre-split names tied to type names
			// through the before registry; the type name — not the factory
			// name, which a split move renames and qualifies — is the stable
			// link. Distinguish the relocations a filename-guessing gate
			// conflated:
			//   - a constructor never registered in before: a dead,
			//     unregistered factory's file moved with the product split
			//     (alicloud_havip's resourceAlicloudHavip was never
			//     registered; the registered resourceAliCloudVpcHaVip lives in
			//     a different file). Not a deletion — no error, no selection.
			//   - a constructor registered in before whose type name is still
			//     in after: a registered implementation moved, often with an
			//     export rename so bare and qualified factories no longer
			//     match. The moved file carries selection; no error here.
			//   - a constructor registered in before whose type name is no
			//     longer in after: a real public registration deletion. Fail
			//     closed rather than waive it as a relocation.
			for typeName, constructor := range before {
				if !names[constructor] {
					continue
				}
				if _, ok := after[typeName]; !ok {
					return nil, fmt.Errorf("registered resource %s removed without a moved registration in %s", typeName, path)
				}
			}
			// No before-registered constructor declared by this file had its type
			// name removed from after: the deletion is a dead-factory relocation,
			// not a registration deletion.
		}
		// A changed split file with no registration is a framework resource,
		// data source or helper outside this checker: the product's
		// register.go is the source of truth for the SDKv2 scope.
	}
	return selected, nil
}

func functionNames(filename string, src []byte) (map[string]bool, error) {
	f, err := parser.ParseFile(token.NewFileSet(), filename, src, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil {
			names[fn.Name.Name] = true
		}
	}
	return names, nil
}

// Names gives deterministic diagnostics and command output.
func Names(selected map[string]string) []string {
	names := make([]string, 0, len(selected))
	for name := range selected {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
