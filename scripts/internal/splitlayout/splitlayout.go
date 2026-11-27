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
	"fmt"
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

const (
	// ProviderSourcePath is the current SDKv2 provider implementation path.
	ProviderSourcePath = "alicloud/provider/sdkv2/provider.go"
	// LegacyProviderSourcePath is the historical root provider implementation path.
	LegacyProviderSourcePath = "alicloud/provider.go"
)

// ProviderSourcePaths returns the current SDKv2 implementation path first and
// the historical root path second for compatibility with older base commits.
func ProviderSourcePaths() []string {
	return []string{ProviderSourcePath, LegacyProviderSourcePath}
}

// IsProviderSource reports whether path holds a full provider implementation
// in either the current SDKv2 layout or the historical root layout.
func IsProviderSource(path string) bool {
	return path == ProviderSourcePath || path == LegacyProviderSourcePath
}

// ReadProviderSource reads the SDKv2 implementation from either the current or
// historical layout. The root compatibility forwarding file is never a fallback
// when the SDKv2 implementation exists.
func ReadProviderSource() ([]byte, error) {
	return ReadProviderSourceFrom(os.ReadFile)
}

// ReadProviderSourceFrom reads the current SDKv2 implementation path first,
// then the historical root path. It is also used for base-revision git reads.
func ReadProviderSourceFrom(read func(string) ([]byte, error)) ([]byte, error) {
	var err error
	for _, path := range ProviderSourcePaths() {
		var src []byte
		src, err = read(path)
		if err == nil {
			return src, nil
		}
	}
	return nil, err
}

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
	// Factory is the constructor the registration names. The resolution is
	// file-grained: the package tolerates one file declaring several entities'
	// constructors, and then every entity of that file resolves to the same
	// ImplPath while its schema lives in one constructor's declaration — a
	// caller that parses the whole file gets the union of the entities'
	// attributes. The factory lets such a caller detect the collision instead
	// of trusting the file alone (see FrameworkFactoriesDeclaredIn).
	Factory string
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
// declared in one dialect, folding a parse failure into no entries. The
// lookup scans tolerate an unparsable register because a broken register
// fails the compile first and the scans only pick entries; the one caller
// that must tell a parse failure from "declares nothing" — the
// file-granularity guard — uses registerEntriesErr instead.
func registerEntries(registerSrc []byte, framework bool) []Entry {
	entries, _ := registerEntriesErr(registerSrc, framework)
	return entries
}

// registerEntriesErr is registerEntries with the parse failure kept as an
// error instead of an empty result: the caller can fail closed on "the
// register could not be read" rather than treat it as "no entries declared",
// which are different verdicts for a guard.
func registerEntriesErr(registerSrc []byte, framework bool) ([]Entry, error) {
	f, err := parser.ParseFile(token.NewFileSet(), "register.go", registerSrc, parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("parsing register.go: %w", err)
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
	return entries, nil
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

// ProviderRegisteredNames parses the SDKv2 provider source and returns the
// entries registered in its DataSourcesMap and ResourcesMap literals — the
// public aliases registered directly by the provider rather than in a product
// register.go. The base title index needs them
// so a migration that moves such an alias into a service register announces
// it as modified, not as new: the alias type name need not match any flat
// filename stem (Quotas #11288 registers alicloud_quotas_application_info at
// the root), so the filename-derived index would otherwise miss it.
//
// Each entry carries the kind its map literal belongs to — DataSourcesMap
// yields KindDataSource and ResourcesMap yields KindResource — so the title
// gate keeps a root DataSource distinct from a same-named service Resource;
// a name-only index would let a same-name Resource/DataSource collision waive
// the wrong declaration. The Factory field is not populated: the title gate
// only needs the type name and kind, and the provider.go value expressions
// are not register factories.
//
// Provider-level registrations live in the SDKv2 provider source, which the
// title gate reads at the base commit the same way it reads every register.go;
// this function only walks the two top-level provider map literals, so nested
// schema maps (map[string]*schema.Schema) are never mistaken for provider entities.
func ProviderRegisteredNames(providerSrc []byte) []Entry {
	f, err := parser.ParseFile(token.NewFileSet(), "provider.go", providerSrc, parser.SkipObjectResolution)
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
		if !ok || key.Name != "DataSourcesMap" && key.Name != "ResourcesMap" {
			return true
		}
		kind := KindResource
		if key.Name == "DataSourcesMap" {
			kind = KindDataSource
		}
		lit, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			return true
		}
		for _, elt := range lit.Elts {
			pair, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			if name := stringLit(pair.Key); name != "" {
				entries = append(entries, Entry{TypeName: name, Kind: kind})
			}
		}
		return true
	})
	return entries
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
			Factory:  e.Factory,
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
//
// A name alone can name two entities: a product may register the same type
// name as a managed resource and as a data source (the SAE NLB pair), and this
// lookup returns whichever entry the register order happens to reach first. A
// caller that knows the kind from the diff it is checking must use
// FrameworkFindByKindAndTypeName so the kind selects the entity.
func FrameworkFindByTypeName(typeName string) *Resolved {
	return findByTypeName(typeName, true)
}

// FrameworkFindByKindAndTypeName reverse-resolves a full framework type name
// together with the kind of entity it names. Matching both is what keeps a
// same-named resource and data source apart: the (KindResource, name) pair
// resolves to the file declaring the resource factory and the
// (KindDataSource, name) pair to the data source one, whichever order the
// register declares them in. It returns nil when no register declares that
// (kind, type name) pair, or when the matching factory is not declared in any
// implementation file.
func FrameworkFindByKindAndTypeName(kind Kind, typeName string) *Resolved {
	return findByKind(typeName, true, &kind)
}

// FrameworkFactoriesDeclaredIn returns the framework registrations whose
// Factory constructors are declared in the given source. The resolution is
// file-grained, so a caller that scopes anything to "the implementation file
// of one entity" — the consistency gate parses the whole file for the
// entity's schema — must first confirm the file declares only that entity:
// with two or more entries returned, the entities' schemas are interleaved in
// one file and cannot be isolated from each other. The entries are counted
// per distinct factory, not per registration: alias registrations reuse one
// factory for several type names (alicloud_disk/alicloud_ecs_disk, the
// sslcertificateservice legacy aliases), and a file declaring one constructor
// still exposes exactly that one schema no matter how many type names point
// at it. An empty result means the file declares no framework factory at
// all; an error means the guard could not prove either way — the source or
// the product's register failed to parse or read — and the caller must fail
// closed on it rather than treat the outcome as "no conflict".
func FrameworkFactoriesDeclaredIn(filename string, src []byte) ([]Entry, error) {
	funcs, err := declaredFuncs(filename, src)
	if err != nil {
		return nil, err
	}
	// The scan is scoped to the file's own product. A register entry counts
	// only when its factory is declared in the registering product's package:
	// the ServicePackage literal references the constructor as a
	// package-local identifier, and a function name is unique inside one
	// package, so only the register of the product this file belongs to can
	// hold entries whose factories this file declares. Matching names across
	// every product's register instead would let another product's same-named
	// factory — Go function names are not repo-unique — inflate the result and
	// fail a perfectly isolated entity, the same wrong-entity class this
	// package exists to eliminate. A file outside any product directory cannot
	// declare a registered framework factory at all, so it scans no register.
	product, ok := ServiceProduct(filename)
	if !ok {
		return nil, nil
	}
	registerSrc, err := os.ReadFile(PackageDir(product) + "/register.go")
	if err != nil {
		return nil, err
	}
	// A register that fails to parse must surface as an error, not as "no
	// entries": the caller counts factories off this result, and a broken
	// register folding into emptiness would turn the guard into a silent pass.
	entries, err := registerEntriesErr(registerSrc, true)
	if err != nil {
		return nil, fmt.Errorf("reading the framework registrations of %s: %w", PackageDir(product), err)
	}
	var declared []Entry
	for _, e := range entries {
		if funcs[e.Factory] {
			declared = append(declared, e)
		}
	}
	return declared, nil
}

func findByTypeName(typeName string, framework bool) *Resolved {
	return findByKind(typeName, framework, nil)
}

// findByKind reverse-resolves a type name across every register in one
// dialect, optionally constrained to one kind. The kind is optional because
// the name-only lookups are the SDKv2 contract and the base-revision reader
// has no kind to offer; a caller that does know the kind must pass it so a
// same-named resource and data source resolve to their own implementations.
func findByKind(typeName string, framework bool, kind *Kind) *Resolved {
	registers, err := RegisterPaths()
	if err != nil {
		return nil
	}
	for _, register := range registers {
		src, err := os.ReadFile(register)
		if err != nil {
			continue
		}
		product, _ := IsServiceRegister(register)
		for _, e := range registerEntries(src, framework) {
			if e.TypeName != typeName {
				continue
			}
			if kind != nil && e.Kind != *kind {
				continue
			}
			if r := locateFactory(product, e); r != nil {
				return r
			}
		}
	}
	return nil
}

// locateFactory scans a product's non-test, non-register source files for the
// top-level declaration of one registration's Factory constructor and joins
// the registration and the declaring file into a Resolved. It returns nil when
// the constructor is not declared in any of them.
func locateFactory(product string, e Entry) *Resolved {
	impls, err := filepath.Glob(PackageDir(product) + "/*.go")
	if err != nil {
		return nil
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
				Factory:  e.Factory,
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
