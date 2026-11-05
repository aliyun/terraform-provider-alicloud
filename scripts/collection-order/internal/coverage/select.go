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

// Select maps both changed implementations and tests to current registrations. A deleted
// implementation is ignored only when its constructors are no longer registered.
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
			// A *_test.go file whose mechanically trimmed implementation never
			// existed is a same-package helper (validation, unit, ...), not the
			// main test of a missing resource. The registry references
			// constructors rather than filenames, and main.go reads every
			// non-test resource_*.go file into fileFunctions, so the helper's
			// real implementation is already present. Associate the helper with
			// that implementation so the resource's main coverage check runs
			// instead of a git show on a path that never existed; a test with no
			// matching implementation still falls through to oldSource below and
			// fails closed.
			if _, ok := fileFunctions[implementation]; !ok {
				if actual, ok := helperImplementation(path, fileFunctions); ok {
					implementation = actual
				}
			}
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
			removed := false
			for _, constructor := range before {
				removed = removed || names[constructor]
			}
			if !removed {
				return nil, fmt.Errorf("cannot establish deletion of registered resource %s", path)
			}
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

// helperImplementation associates a *_test.go file whose mechanically trimmed
// implementation path never existed with the actual implementation file in the
// same package directory. The registry references constructors, not filenames,
// and main.go reads every non-test resource_*.go file into fileFunctions, so a
// helper such as resource_example_validation_test.go has its real
// implementation resource_example.go already present. Matching the longest
// implementation stem that prefixes the test stem keeps
// resource_example_extra_validation_test.go with resource_example_extra.go
// rather than resource_example.go. A main test (resource_example_test.go) is
// never a helper: its trimmed path is the implementation itself, which exists
// in fileFunctions, so this function is only reached for helpers.
func helperImplementation(testPath string, fileFunctions map[string]map[string]bool) (string, bool) {
	testDir, testBase := splitPath(testPath)
	if !strings.HasSuffix(testBase, "_test.go") {
		return "", false
	}
	testStem := strings.TrimSuffix(testBase, "_test.go")
	var best string
	bestLen := -1
	for implPath := range fileFunctions {
		implDir, implBase := splitPath(implPath)
		if implDir != testDir {
			continue
		}
		implStem := strings.TrimSuffix(implBase, ".go")
		if implStem == testStem || !strings.HasPrefix(testStem, implStem+"_") {
			continue
		}
		if len(implStem) > bestLen {
			best = implPath
			bestLen = len(implStem)
		}
	}
	return best, best != ""
}

// splitPath returns the directory and base name of a slash-separated
// repository-relative path.
func splitPath(path string) (dir, base string) {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[:i], path[i+1:]
	}
	return "", path
}
