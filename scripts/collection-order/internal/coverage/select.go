package coverage

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
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
		if !strings.HasPrefix(path, "alicloud/resource_alicloud_") || !strings.HasSuffix(path, ".go") {
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
			names, err = functionNames(implementation, src)
			if err != nil {
				return nil, err
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
		if !matched && exists {
			return nil, fmt.Errorf("changed resource file %s has no current ResourcesMap registration", path)
		}
		if !matched && !exists {
			removed := false
			for _, constructor := range before {
				removed = removed || names[constructor]
			}
			if !removed {
				return nil, fmt.Errorf("cannot establish deletion of registered resource %s", path)
			}
		}
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
