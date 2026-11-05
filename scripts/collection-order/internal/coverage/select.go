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
	return registryFor(src, "ResourcesMap")
}

// DataSourcesRegistry reads the data-source constructor identities from the
// provider.go DataSourcesMap literal, the same way Registry reads
// ResourcesMap. Data sources are read-only (Read only), so their constructor
// resolves the schema the order-coverage checker inspects the same way it does
// for managed resources.
func DataSourcesRegistry(src []byte) (map[string]string, error) {
	return registryFor(src, "DataSourcesMap")
}

// registryFor parses one provider entity map literal (ResourcesMap or
// DataSourcesMap) into a type-name → constructor-identifier index. The map
// name is parameterised so the same logic serves resources and data sources;
// a flat provider.go registers both kinds at the root.
func registryFor(src []byte, mapName string) (map[string]string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "provider.go", src, 0)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	var parseErr error
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok || ident(kv.Key) != mapName {
			return true
		}
		if found {
			parseErr = fmt.Errorf("multiple %s declarations", mapName)
			return false
		}
		found = true
		entries := fields(kv.Value)
		if entries == nil {
			parseErr = fmt.Errorf("%s must be a literal map", mapName)
			return false
		}
		for name, expr := range entries {
			call, ok := expr.(*ast.CallExpr)
			if !ok || ident(call.Fun) == "" || len(call.Args) != 0 {
				parseErr = fmt.Errorf("unresolved constructor for %s", name)
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
		return nil, fmt.Errorf("%s is missing", mapName)
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
	return unionRegistry(providerSrc, registers, splitlayout.KindResource)
}

// DataSourcesRegistryUnion merges the provider.go DataSourcesMap literal with
// the SDKv2 data-source registrations declared by the split service packages,
// the same way RegistryUnion merges managed resources. Data sources are a
// distinct registration kind (SDKDataSources / DataSourcesMap), so they never
// collide with a same-named resource: the two registries stay separate and
// the order-coverage checker runs each against its own runtime schema map.
func DataSourcesRegistryUnion(providerSrc []byte, registers map[string][]byte) (map[string]string, error) {
	return unionRegistry(providerSrc, registers, splitlayout.KindDataSource)
}

// unionRegistry reads the root map literal for kind and folds in the matching
// SDKv2 registrations from the split service packages. Only entries of the
// requested kind contribute; the other kinds (a register's SDKDataSources
// when kind is resource, or its SDKResources when kind is data source) stay
// out, so a product's data sources never leak into the resource registry and
// vice versa.
func unionRegistry(providerSrc []byte, registers map[string][]byte, kind splitlayout.Kind) (map[string]string, error) {
	root, err := registryFor(providerSrc, rootMapName(kind))
	if err != nil {
		return nil, err
	}
	for path, src := range registers {
		product, ok := splitlayout.IsServiceRegister(path)
		if !ok {
			continue
		}
		for _, entry := range splitlayout.SDKRegisterEntries(src) {
			if entry.Kind != kind {
				continue
			}
			qualified := product + "." + entry.Factory
			if previous, dup := root[entry.TypeName]; dup {
				return nil, fmt.Errorf("duplicate %s declaration for %s: %s and %s", kind, entry.TypeName, previous, qualified)
			}
			root[entry.TypeName] = qualified
		}
	}
	return root, nil
}

// rootMapName maps a registration kind to the provider.go map literal that
// carries its flat-layout aliases.
func rootMapName(kind splitlayout.Kind) string {
	if kind == splitlayout.KindDataSource {
		return "DataSourcesMap"
	}
	return "ResourcesMap"
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

// Scope narrows Select and Check to one provider entity kind. ResourceScope
// selects managed-resource implementations and tests; DataSourceScope selects
// data-source ones. The flat file prefix and split-layout stem differ between
// the two (resource_alicloud_/resource_ vs data_source_alicloud_/data_source_),
// so each scope only ever matches its own implementation files; a data source
// changed in the same diff never leaks into a resource Select and vice versa.
// readOnly records that a data source has no write lifecycle: its Read is
// order-insensitive, so a convergent reorder of a configurable TypeList filter
// is trivially satisfied without an explicit reorder test.
type Scope struct {
	label      string
	flatPrefix string
	splitStem  string
	readOnly   bool
}

// Label returns the human-readable scope name used in diagnostics. It is the
// exported accessor for the unexported label field, so callers outside the
// coverage package (such as the collection-order main command) can include the
// scope in error messages without reaching into the struct fields directly.
func (s Scope) Label() string { return s.label }

var (
	// ResourceScope selects managed-resource implementations and tests.
	ResourceScope = Scope{label: "resource", flatPrefix: "alicloud/resource_alicloud_", splitStem: "resource_", readOnly: false}
	// DataSourceScope selects data-source implementations and tests. Data
	// sources are read-only, so their TypeList filter order is semantically
	// irrelevant to the API result and the convergent reorder is trivially
	// satisfied by the read-only lifecycle.
	DataSourceScope = Scope{label: "data source", flatPrefix: "alicloud/data_source_alicloud_", splitStem: "data_source_", readOnly: true}
)

// Select maps both changed implementations and tests to current registrations. A deleted
// implementation is ignored only when its constructors are no longer registered. It
// selects managed resources; data sources use SelectScope with DataSourceScope.
func Select(changed []string, before, after map[string]string, sources map[string][]byte, oldSource func(string) ([]byte, error)) (map[string]string, error) {
	return SelectScope(changed, before, after, sources, oldSource, ResourceScope)
}

// SelectScope maps changed implementations and tests to current registrations
// for the given scope. The body is entity-kind-agnostic: a data source's
// constructor resolves its schema the same way a resource's does, and the
// convergence check in Check reads the test steps by builder variable rather
// than address. Data sources are read-only, so the only coverage target is
// the schema's configurable multi-member TypeList attributes (the order in
// which filter blocks such as ids/enabled_filters are configured); the check
// is advisory for data sources just as it is for resources without such
// attributes. A deleted implementation is ignored only when its
// constructors are no longer registered.
func SelectScope(changed []string, before, after map[string]string, sources map[string][]byte, oldSource func(string) ([]byte, error), scope Scope) (map[string]string, error) {
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
			return fmt.Errorf("cannot locate implementation for registered %s %s (%s)", scope.label, name, after[name])
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
		// untouched products contribute nothing. Only the scope's kind
		// contributes, because unionRegistry filtered by kind already.
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
		flat := strings.HasPrefix(path, scope.flatPrefix) && strings.HasSuffix(path, ".go")
		_, stem, _, isSplit := splitlayout.IsSplitSource(path)
		isSplit = isSplit && strings.HasPrefix(stem, scope.splitStem)
		if !flat && !isSplit {
			// The other entity kind's files, docs and framework-only
			// service files stay outside this scope, as they did in the
			// flat layout. A data source's Select never matches a resource_
			// file and a resource's Select never matches a data_source_ file,
			// so the two scopes do not steal each other's changes.
			continue
		}
		implementation := strings.TrimSuffix(path, "_test.go")
		if implementation != path {
			implementation += ".go"
			// A *_test.go file whose mechanically trimmed implementation never
			// existed is a same-package helper (validation, unit, ...), not the
			// main test of a missing resource. The registry references
			// constructors rather than filenames, and main.go reads every
			// non-test implementation file into fileFunctions, so the helper's
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
				return nil, fmt.Errorf("cannot resolve changed %s %s: %w", scope.label, path, err)
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
			return nil, fmt.Errorf("changed %s file %s has no current registration", scope.label, path)
		}
		if !matched && !exists && flat {
			removed := false
			for _, constructor := range before {
				removed = removed || names[constructor]
			}
			if !removed {
				return nil, fmt.Errorf("cannot establish deletion of registered %s %s", scope.label, path)
			}
		}
		// A changed split file with no registration is a framework resource,
		// framework data source or helper outside this scope: the product's
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
