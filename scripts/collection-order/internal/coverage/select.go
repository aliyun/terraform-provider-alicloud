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
	ast.Inspect(f, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok || ident(kv.Key) != "ResourcesMap" {
			return true
		}
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
	if len(result) == 0 {
		return nil, fmt.Errorf("ResourcesMap is missing or empty")
	}
	return result, nil
}

// RegistryUnion merges the provider.go ResourcesMap literal with the SDKv2
// registrations declared by the split service packages. Split constructors
// carry a product qualifier ("cloudauth.ResourceFaceConfig") because
// top-level names are only unique per package; the flat literal keeps bare
// names. A declaration appearing in both sources is the same startup error
// assembleProvider reports.
//
// Framework-native registrations are intentionally excluded: this returns the
// SDKv2 resource registry, and a framework resource must not leak into it (the
// runtime provider validates that SDK and framework TypeNames never collide).
// The collection-order check reaches framework resources separately through
// FrameworkRegistry/CombinedRegistry, synthesizing their schema via
// FrameworkResource, so they are checked with the same fail-closed semantics
// without being conflated with the SDK resource registry.
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

// FrameworkRegistry collects the framework-native resource registrations from
// the split service packages into a map kept separate from the SDK resource
// registry (RegistryUnion). Separating them preserves the SDK-only invariant
// the resource registry holds, while still letting the collection-order check
// apply to framework resources: CombinedRegistry overlays this on the SDK
// registry, and for a framework TypeName the runtime ResourcesMap lookup
// returns nil so the check synthesizes the real schema via FrameworkResource.
// A duplicate framework declaration is an error, not a silent overwrite.
func FrameworkRegistry(registers map[string][]byte) (map[string]string, error) {
	reg := map[string]string{}
	for path, src := range registers {
		product, ok := splitlayout.IsServiceRegister(path)
		if !ok {
			continue
		}
		for _, entry := range splitlayout.FrameworkRegisterEntries(src) {
			if entry.Kind != splitlayout.KindResource {
				continue
			}
			qualified := product + "." + entry.Factory
			if previous, dup := reg[entry.TypeName]; dup {
				return nil, fmt.Errorf("duplicate framework resource declaration for %s: %s and %s", entry.TypeName, previous, qualified)
			}
			reg[entry.TypeName] = qualified
		}
	}
	return reg, nil
}

// CombinedRegistry returns the SDK resource registry (RegistryUnion) overlaid
// with the framework-native registry (FrameworkRegistry). It is what the
// collection-order check selects and checks against: SDK resources via their
// runtime schema, framework resources via the synthesized schema when the
// runtime ResourcesMap lookup returns nil. A TypeName appearing in both is an
// error: the runtime provider validates that SDK and framework TypeNames never
// collide, so a collision here means that validation was bypassed.
func CombinedRegistry(providerSrc []byte, registers map[string][]byte) (map[string]string, error) {
	sdk, err := RegistryUnion(providerSrc, registers)
	if err != nil {
		return nil, err
	}
	framework, err := FrameworkRegistry(registers)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(sdk)+len(framework))
	for k, v := range sdk {
		out[k] = v
	}
	for k, v := range framework {
		if _, dup := out[k]; dup {
			return nil, fmt.Errorf("TypeName %s registered in both SDK and framework registries", k)
		}
		out[k] = v
	}
	return out, nil
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
		isTest := implementation != path
		if isTest {
			implementation += ".go"
		}
		names, exists := fileFunctions[implementation]
		helperResolved := false
		if !exists {
			src, err := oldSource(implementation)
			if err != nil {
				// A framework-native resource may keep its ACC test in the flat
				// alicloud/ directory (package alicloud_test) so `go test
				// ./alicloud` discovers it, while its real implementation lives
				// in alicloud/service/<product>/. The flat impl path the standard
				// pairing derives (resource_alicloud_<short>.go) exists at neither
				// revision; resolve it to the split-layout impl
				// (alicloud/service/<product>/resource_<short>.go) so the
				// collection-order check still applies to the registered resource.
				//
				// A helper test (resource_foo_crud_test.go,
				// resource_foo_internal_test.go) derives a resource_foo_crud.go
				// impl that exists at neither revision: its real impl is
				// resource_foo.go. Resolve it by stripping the single helper
				// suffix segment so the check still applies to the registered
				// resource the test belongs to — coverage is preserved, not
				// blanket-ignored. The repo convention is a single helper suffix
				// segment; an orphan test whose one-strip candidate is absent,
				// or whose resolved impl is not registered, fails closed rather
				// than silently skipping a changed test.
				if !isTest {
					return nil, fmt.Errorf("cannot resolve changed resource %s: %w", path, err)
				}
				if resolved, ok := flatTestSplitImpl(implementation, fileFunctions); ok {
					implementation = resolved
				} else if resolved, ok := helperTestImpl(implementation, fileFunctions); ok {
					implementation = resolved
				} else {
					return nil, fmt.Errorf("cannot resolve changed resource %s: %w", path, err)
				}
				names = fileFunctions[implementation]
				exists = true
				helperResolved = true
			} else {
				base, err := functionNames(implementation, src)
				if err != nil {
					return nil, err
				}
				names = qualify(implementation, base)
			}
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
		if !matched && helperResolved {
			return nil, fmt.Errorf("changed helper test %s resolves to %s, which is not a registered resource", path, implementation)
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

// helperTestImpl resolves a changed helper test file (e.g.
// alicloud/service/ecs/resource_security_group_snapshot_policy_crud_test.go)
// to the registered resource's implementation file by stripping the single
// helper suffix segment from the derived impl stem. The standard pairing
// (resource_foo_test.go → resource_foo.go) is handled before this is reached:
// this only runs when the derived impl (resource_foo_crud.go) is absent from
// the working tree and the base revision, i.e. the changed file is a
// CRUD/internal helper test whose real impl is resource_foo.go. Stripping is
// limited to one segment and never proceeds past the leading "resource" stem,
// so a different, shorter registered resource cannot be matched by accident;
// a test whose one-strip candidate is absent returns ("", false) so the caller
// fails closed instead of silently skipping a changed test.
func helperTestImpl(impl string, fileFunctions map[string]map[string]bool) (string, bool) {
	slash := strings.LastIndex(impl, "/")
	if slash < 0 {
		return "", false
	}
	dir := impl[:slash]
	name := strings.TrimSuffix(impl[slash+1:], ".go")
	i := strings.LastIndex(name, "_")
	if i <= 0 {
		return "", false
	}
	stripped := name[:i]
	if stripped == "resource" {
		return "", false
	}
	candidate := dir + "/" + stripped + ".go"
	if candidate == impl {
		return "", false
	}
	if _, ok := fileFunctions[candidate]; ok {
		return candidate, true
	}
	return "", false
}

// flatTestSplitImpl resolves a flat-layout test file (e.g.
// alicloud/resource_alicloud_security_group_snapshot_policy_test.go)
// to the split-layout implementation file when the resource is
// framework-native and keeps its ACC test in the flat alicloud/
// directory (package alicloud_test) so `go test ./alicloud` discovers
// it. The flat impl path the standard pairing derives
// (resource_alicloud_<short>.go) exists at neither revision; the real
// impl is alicloud/service/<product>/resource_<short>.go. The short
// name is derived by stripping the flat prefix/suffix and searched in
// the working-tree fileFunctions; a mismatch returns ("", false) so
// the caller's existing fail-closed path applies.
func flatTestSplitImpl(flatImpl string, fileFunctions map[string]map[string]bool) (string, bool) {
	const flatPrefix = "alicloud/resource_alicloud_"
	if !strings.HasPrefix(flatImpl, flatPrefix) || !strings.HasSuffix(flatImpl, ".go") {
		return "", false
	}
	short := strings.TrimSuffix(flatImpl[len(flatPrefix):], ".go")
	if short == "" {
		return "", false
	}
	suffix := "/resource_" + short + ".go"
	for path := range fileFunctions {
		if strings.HasPrefix(path, "alicloud/service/") && strings.HasSuffix(path, suffix) {
			return path, true
		}
	}
	return "", false
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
