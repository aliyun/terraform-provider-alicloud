// Package splitlayout maps product and provider/core source files to the
// registrations they implement.
//
// A split-layout path cannot yield the resource name on its own: the product
// slug only coincidentally matches the resource-name prefix (alicloud_cen_*
// will live in service/cbn). The registration in the product's register.go is
// the source of truth: it pairs each TypeName with the Factory constructor,
// and that constructor is declared in exactly one implementation file.
// Resolve joins those two facts.
//
// Registrations come in two dialects. The SDKv2 dialect fills
// conns.ServicePackage SDKResources/SDKDataSources; the framework dialect
// fills Resources/DataSources/EphemeralResources/ListResources/Actions. Both
// name legitimate provider entities: SDKv2 entries feed the diff checkers
// that run against the provider runtime, framework entries feed the
// AST-based breaking-change checker and the explicit out-of-scope logs of
// the runtime-based checkers.
package splitlayout

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	// serviceSourceRegex matches product and provider/core source paths for
	// resource, data-source, ephemeral, action and list-resource
	// implementation and test files. The product segment is a literal
	// character class so a diff path can never smuggle in path separators or
	// traversal segments.
	serviceSourceRegex = regexp.MustCompile(`^alicloud/(?:service/([0-9A-Za-z_-]+)|(provider/core))/((?:resource|data_source|ephemeral|action|list)_[0-9A-Za-z_]+?)(_test)?\.go$`)
	// servicePathRegex matches files directly inside a product or core package.
	servicePathRegex = regexp.MustCompile(`^alicloud/(?:service/([0-9A-Za-z_-]+)|(provider/core))/[0-9A-Za-z_.-]+$`)
)

// Kind classifies a registration by the kind of provider entity it names.
type Kind string

const (
	KindResource     Kind = "resource"
	KindDataSource   Kind = "data_source"
	KindEphemeral    Kind = "ephemeral"
	KindListResource Kind = "list"
	KindAction       Kind = "action"
)

// Entry is one registration declared by a register.go.
type Entry struct {
	TypeName string // e.g. "alicloud_cloudauth_face_config"
	Factory  string // constructor identifier, e.g. "ResourceFaceConfig"
	Kind     Kind   // which provider entity the registration names
}

// IsResource reports whether the entry names a managed resource, as opposed
// to a data source, ephemeral, list resource or action.
func (e Entry) IsResource() bool { return e.Kind == KindResource }

// Resolved joins a split-layout source file to its registration.
type Resolved struct {
	TypeName string
	ImplPath string
	TestPath string
	Kind     Kind
}

// IsResource reports whether the resolved registration names a managed
// resource.
func (r Resolved) IsResource() bool { return r.Kind == KindResource }

// IsSplitSource reports whether path is a split-layout resource, data-source,
// ephemeral, action or list-resource source file (implementation or test),
// and returns the product directory and the implementation basename stem
// (without the _test suffix). Provider-level core files use "provider/core"
// as their package key, rather than a product name.
func IsSplitSource(path string) (product, stem string, isTest bool, ok bool) {
	m := serviceSourceRegex.FindStringSubmatch(path)
	if m == nil {
		return "", "", false, false
	}
	product = m[1]
	if product == "" {
		product = m[2]
	}
	return product, m[3], m[4] == "_test", true
}

// ServiceProduct returns the package key for a file directly inside a product
// directory or provider/core. Core uses the distinct key "provider/core", so it
// cannot collide with a product named core. Other provider internals are excluded.
func ServiceProduct(path string) (string, bool) {
	m := servicePathRegex.FindStringSubmatch(path)
	if m == nil {
		return "", false
	}
	if m[1] == "" {
		return m[2], true
	}
	return m[1], true
}

// PackageDir maps a package key to its source directory. provider/core is
// reserved for provider-level entities; product keys keep their service path.
func PackageDir(product string) string {
	if product == "provider/core" {
		return "alicloud/provider/core"
	}
	return "alicloud/service/" + product
}

// RegisterPaths lists the supported product and provider-level registrations.
func RegisterPaths() ([]string, error) {
	var paths []string
	for _, pattern := range []string{"alicloud/service/*/register.go", "alicloud/provider/core/register.go"} {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		paths = append(paths, matches...)
	}
	sort.Strings(paths)
	return paths, nil
}

// IsServiceRegister reports whether path is a supported register.go, and
// returns its package key. A register.go is a registration source the way
// provider.go is.
func IsServiceRegister(path string) (string, bool) {
	product, ok := ServiceProduct(path)
	if !ok || !strings.HasSuffix(path, "/register.go") {
		return "", false
	}
	return product, true
}

// registerKinds maps the conns.ServicePackage registration fields to the kind
// of entity they name and the dialect they belong to. Functions are not
// listed: a provider function has no service-package implementation file to
// resolve.
var registerKinds = map[string]struct {
	kind      Kind
	framework bool
}{
	"SDKResources":       {KindResource, false},
	"SDKDataSources":     {KindDataSource, false},
	"Resources":          {KindResource, true},
	"DataSources":        {KindDataSource, true},
	"EphemeralResources": {KindEphemeral, true},
	"ListResources":      {KindListResource, true},
	"Actions":            {KindAction, true},
}

// registerEntries parses one register.go source and returns the entries
// declared in one dialect.
func registerEntries(registerSrc []byte, framework bool) []Entry {
	f, err := parser.ParseFile(token.NewFileSet(), "register.go", registerSrc, parser.SkipObjectResolution)
	if err != nil {
		return nil
	}
	var entries []Entry
	ast.Inspect(f, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			return true
		}
		meta, ok := registerKinds[key.Name]
		if !ok || meta.framework != framework {
			return true
		}
		lit, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			return true
		}
		for _, elt := range lit.Elts {
			entry, ok := elt.(*ast.CompositeLit)
			if !ok {
				continue
			}
			var e Entry
			for _, fld := range entry.Elts {
				fkv, ok := fld.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				switch identName(fkv.Key) {
				case "TypeName":
					e.TypeName = stringLit(fkv.Value)
				case "Factory":
					e.Factory = identName(fkv.Value)
				}
			}
			if e.TypeName != "" && e.Factory != "" {
				e.Kind = meta.kind
				entries = append(entries, e)
			}
		}
		return true
	})
	return entries
}

// SDKRegisterEntries parses one register.go source and returns its SDKv2
// entries. Framework registrations are ignored: they are outside the SDKv2
// diff checkers' scope.
func SDKRegisterEntries(registerSrc []byte) []Entry {
	return registerEntries(registerSrc, false)
}

// FrameworkRegisterEntries parses one register.go source and returns its
// framework-native entries: managed resources, data sources, ephemerals,
// list resources and actions.
func FrameworkRegisterEntries(registerSrc []byte) []Entry {
	return registerEntries(registerSrc, true)
}

// FrameworkRegisteredIn reports whether one register.go source declares
// typeName through a framework-native registration (Resources, DataSources,
// EphemeralResources, ListResources or Actions).
func FrameworkRegisteredIn(registerSrc []byte, typeName string) bool {
	for _, e := range FrameworkRegisterEntries(registerSrc) {
		if e.TypeName == typeName {
			return true
		}
	}
	return false
}

// FrameworkRegistered reports whether any product register.go declares
// typeName as framework-native. It lets the SDKv2 diff checkers tell a
// framework registration (outside their scope) from an unregistered name.
func FrameworkRegistered(typeName string) bool {
	registers, err := RegisterPaths()
	if err != nil {
		return false
	}
	for _, register := range registers {
		src, err := os.ReadFile(register)
		if err != nil {
			continue
		}
		if FrameworkRegisteredIn(src, typeName) {
			return true
		}
	}
	return false
}

// Resolve joins a split-layout source file to the SDKv2 registrations it
// implements, by matching the constructors declared in its implementation
// file against the register's Factory fields. An empty result means the file
// is not an SDKv2 resource or data source (framework-native product, shared
// helper): callers keep skipping such files, as they did before the split.
func Resolve(product, stem string) ([]Resolved, error) {
	return resolve(product, stem, false)
}

// FrameworkResolve joins a split-layout source file to the framework-native
// registrations it implements — managed resources, data sources, ephemerals,
// list resources and actions — the same way Resolve does for the SDKv2
// dialect.
func FrameworkResolve(product, stem string) ([]Resolved, error) {
	return resolve(product, stem, true)
}

func resolve(product, stem string, framework bool) ([]Resolved, error) {
	registerSrc, err := os.ReadFile(PackageDir(product) + "/register.go")
	if err != nil {
		return nil, err
	}
	entries := registerEntries(registerSrc, framework)
	if len(entries) == 0 {
		return nil, nil
	}
	implPath := PackageDir(product) + "/" + stem + ".go"
	implSrc, err := os.ReadFile(implPath)
	if err != nil {
		return nil, err
	}
	funcs, err := declaredFuncs(implPath, implSrc)
	if err != nil {
		return nil, err
	}
	testPath := strings.TrimSuffix(implPath, ".go") + "_test.go"
	var out []Resolved
	for _, e := range entries {
		if !funcs[e.Factory] {
			continue
		}
		out = append(out, Resolved{
			TypeName: e.TypeName,
			ImplPath: implPath,
			TestPath: testPath,
			Kind:     e.Kind,
		})
	}
	return out, nil
}

// FindByTypeName reverse-resolves a full SDKv2 type name (e.g.
// "alicloud_cloudauth_face_config") to its split-layout paths. It returns
// nil when the name is not registered in any service package.
func FindByTypeName(typeName string) *Resolved {
	return findByTypeName(typeName, false)
}

// FrameworkFindByTypeName reverse-resolves a full framework type name — a
// managed resource, data source, ephemeral, list resource or action — to its
// split-layout paths. It returns nil when the name is not registered as
// framework-native in any service package.
func FrameworkFindByTypeName(typeName string) *Resolved {
	return findByTypeName(typeName, true)
}

func findByTypeName(typeName string, framework bool) *Resolved {
	registers, err := RegisterPaths()
	if err != nil {
		return nil
	}
	for _, register := range registers {
		src, err := os.ReadFile(register)
		if err != nil {
			continue
		}
		for _, e := range registerEntries(src, framework) {
			if e.TypeName != typeName {
				continue
			}
			product, _ := IsServiceRegister(register)
			impls, err := filepath.Glob(PackageDir(product) + "/*.go")
			if err != nil {
				continue
			}
			for _, impl := range impls {
				if strings.HasSuffix(impl, "_test.go") || strings.HasSuffix(impl, "register.go") {
					continue
				}
				src, err := os.ReadFile(impl)
				if err != nil {
					continue
				}
				funcs, err := declaredFuncs(impl, src)
				if err != nil {
					continue
				}
				if funcs[e.Factory] {
					return &Resolved{
						TypeName: e.TypeName,
						ImplPath: impl,
						TestPath: strings.TrimSuffix(impl, ".go") + "_test.go",
						Kind:     e.Kind,
					}
				}
			}
		}
	}
	return nil
}

// declaredFuncs returns the top-level function names declared in src.
func declaredFuncs(filename string, src []byte) (map[string]bool, error) {
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

func identName(expr ast.Expr) string {
	if i, ok := expr.(*ast.Ident); ok {
		return i.Name
	}
	return ""
}

func stringLit(expr ast.Expr) string {
	if s, ok := expr.(*ast.BasicLit); ok && s.Kind == token.STRING {
		if v, err := strconv.Unquote(s.Value); err == nil {
			return v
		}
	}
	return ""
}
