// Package flatlayout joins flat alicloud/(resource|data_source)_*.go
// implementation files to the provider registrations they implement.
//
// A flat filename is not a reliable entity name: shared plumbing files
// (data_source_alicloud_common.go) declare no constructor at all, and stock
// files can declare a constructor registered under a different TypeName
// (resource_alicloud_edas_application_package_attachment.go implements
// alicloud_edas_application_deployment). The registration in
// alicloud/provider.go is the source of truth: each ResourcesMap or
// DataSourcesMap entry pairs a TypeName with the constructor it calls, and
// that constructor is declared in exactly one implementation file. Join
// matches those two facts, the same principle scripts/internal/splitlayout
// applies to the split layout.
package flatlayout

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"sort"
	"strconv"
)

// flatSourceRegex guards the paths Join reads: they come from a diff file,
// so a literal character class keeps a path from smuggling in separators or
// traversal segments.
var flatSourceRegex = regexp.MustCompile(`^alicloud/(resource|data_source)[0-9A-Za-z_]*\.go$`)

// Registration is one provider.go map entry: a TypeName paired with the
// constructor the map literal calls.
type Registration struct {
	TypeName   string // e.g. "alicloud_edas_application_deployment"
	Factory    string // e.g. "resourceAlicloudEdasApplicationPackageAttachment"
	DataSource bool   // registered in DataSourcesMap rather than ResourcesMap
}

// Result joins one flat implementation file to the provider registrations.
type Result struct {
	// Registrations holds the provider.go registrations whose constructors
	// the file declares — one per TypeName, so a constructor backing an
	// alias pair yields two.
	Registrations []Registration
	// Unregistered holds constructor-shaped functions the file declares
	// that no provider.go entry calls — legacy V0 schema helpers, dead
	// code, or a genuinely missing registration.
	Unregistered []string
}

// ImplPathOf returns the implementation path of a flat source or test
// path: a _test.go suffix is stripped and .go appended; any other path is
// returned unchanged.
func ImplPathOf(path string) string {
	if n := len(path); n > len("_test.go") && path[n-len("_test.go"):] == "_test.go" {
		return path[:n-len("_test.go")] + ".go"
	}
	return path
}

// ProviderRegistrations parses the alicloud/provider.go source and returns
// the constructor → registration join table built from the ResourcesMap and
// DataSourcesMap literals. A constructor called in either literal is a
// registered entity; flat entities are registered nowhere else (the service
// packages' register.go files feed only the split layout). A constructor can
// back several TypeNames (deprecated aliases such as alicloud_havip and
// alicloud_vpc_ha_vip), so the table holds one Registration per TypeName.
func ProviderRegistrations(providerSrc []byte) map[string][]Registration {
	f, err := parser.ParseFile(token.NewFileSet(), "provider.go", providerSrc, parser.SkipObjectResolution)
	if err != nil {
		return nil
	}
	regs := map[string][]Registration{}
	ast.Inspect(f, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok || (key.Name != "ResourcesMap" && key.Name != "DataSourcesMap") {
			return true
		}
		lit, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			return true
		}
		for _, elt := range lit.Elts {
			entry, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			typeName := stringLit(entry.Key)
			call, ok := entry.Value.(*ast.CallExpr)
			if !ok || typeName == "" {
				continue
			}
			ctor, ok := call.Fun.(*ast.Ident)
			if !ok || ctor.Name == "" {
				continue
			}
			regs[ctor.Name] = append(regs[ctor.Name], Registration{
				TypeName:   typeName,
				Factory:    ctor.Name,
				DataSource: key.Name == "DataSourcesMap",
			})
		}
		return true
	})
	return regs
}

// Join reads the flat implementation file at path (relative to the repo
// root, like the diff paths) and joins the constructor-shaped functions it
// declares against the alicloud/provider.go map literals. A file declaring
// no constructor at all yields an empty Result: it is shared plumbing, not
// an entity source.
func Join(path string) (Result, error) {
	providerSrc, err := os.ReadFile("alicloud/provider.go")
	if err != nil {
		return Result{}, err
	}
	return JoinWith(path, ProviderRegistrations(providerSrc))
}

// JoinWith is Join against an already-parsed registration table.
func JoinWith(path string, regs map[string][]Registration) (Result, error) {
	if !flatSourceRegex.MatchString(path) {
		return Result{}, os.ErrInvalid
	}
	src, err := os.ReadFile(path)
	if err != nil {
		return Result{}, err
	}
	return JoinSrc(path, src, regs)
}

// JoinSrc is the pure join of one flat source file against a registration
// table: the constructor-shaped functions it declares, each resolved to the
// registrations calling it.
func JoinSrc(filename string, src []byte, regs map[string][]Registration) (Result, error) {
	ctors, err := declaredConstructors(filename, src)
	if err != nil {
		return Result{}, err
	}
	var result Result
	for _, name := range ctors {
		if rs, ok := regs[name]; ok {
			result.Registrations = append(result.Registrations, rs...)
		} else {
			result.Unregistered = append(result.Unregistered, name)
		}
	}
	return result, nil
}

// declaredConstructors returns the sorted names of the top-level functions
// in src that have the constructor shape: no receiver, no parameters and a
// single *schema.Resource result. CRUD handlers take arguments and return
// an error; schema helpers return other types, so both fall outside.
func declaredConstructors(filename string, src []byte) ([]string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), filename, src, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Type.Params == nil || len(fn.Type.Params.List) != 0 {
			continue
		}
		if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
			continue
		}
		if !isSchemaResourcePtr(fn.Type.Results.List[0].Type) {
			continue
		}
		names = append(names, fn.Name.Name)
	}
	// Sort so callers log stable, testable output.
	sort.Strings(names)
	return names, nil
}

// isSchemaResourcePtr reports whether expr is *schema.Resource.
func isSchemaResourcePtr(expr ast.Expr) bool {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Resource" {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	return ok && ident.Name == "schema"
}

func stringLit(expr ast.Expr) string {
	if s, ok := expr.(*ast.BasicLit); ok && s.Kind == token.STRING {
		if v, err := strconv.Unquote(s.Value); err == nil {
			return v
		}
	}
	return ""
}
