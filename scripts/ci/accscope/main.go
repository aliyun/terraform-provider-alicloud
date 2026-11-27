// ACC impact selection from Git snapshots. This program never builds Provider packages.
package main

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const module = "github.com/aliyun/terraform-provider-alicloud/"
const framework = "github.com/hashicorp/terraform-plugin-framework/"

type reference struct {
	pkg, name, receiver string
	method, factory     bool
	optionalCallback    bool
	pointer             bool
	line                int
	field               string
	resultIndex         int
}

type declaration struct {
	id, pkg, name, file, receiver, hash string
	refs                                []reference
	results                             []reference
	factoryAlias                        reference
	signature                           string
	requestPackage                      string
	interfaceMethods                    []string
	structFields                        map[string]reference
	mayPromoteOptionalCallbacks         bool
	assets                              []string
	test                                bool
	isFunc                              bool
	acc, unit                           bool
	assertedTypes                       []string
	privateTestFunc                     bool
	sweeperAddress                      string
	staticRegistration                  bool
}

type snapshot struct {
	decls              map[string]*declaration
	files              map[string]bool
	assets             map[string]string
	hazards            map[string]string
	registrationIssues map[string]string
	registrations      map[string]map[string]string
	packages           map[string]string
	importers          map[string]map[string]bool
}

type accScope struct {
	Packages                 map[string][]string `json:"packages"`
	ResourceFiles            []string            `json:"resourceFiles"`
	Reasons                  map[string][]string `json:"reasons"`
	Blocked                  []string            `json:"blocked"`
	UnverifiedCoverage       []string            `json:"unverifiedCoverage,omitempty"`
	UnreferencedDeclarations []string            `json:"unreferencedDeclarations,omitempty"`
	AccOnlyBlockers          []string            `json:"accOnlyBlockers,omitempty"`
	CoverageBlockers         []string            `json:"coverageBlockers"`
	LayoutIssues             []layoutIssue       `json:"layoutIssues,omitempty"`
	FullPackages             map[string][]string `json:"fullPackages,omitempty"`
}

func command(args ...string) ([]byte, error) {
	return exec.Command("git", args...).Output()
}

func nodeText(fset *token.FileSet, node ast.Node) string {
	var b bytes.Buffer
	_ = format.Node(&b, fset, node)
	return b.String()
}

func receiverName(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return receiverName(x.X)
	case *ast.IndexExpr:
		return receiverName(x.X)
	case *ast.IndexListExpr:
		return receiverName(x.X)
	}
	return ""
}

func typeReference(expr ast.Expr, pkg string, imports map[string]string) reference {
	switch x := expr.(type) {
	case *ast.Ident:
		return reference{pkg: pkg, name: x.Name}
	case *ast.StarExpr:
		ref := typeReference(x.X, pkg, imports)
		ref.pointer = true
		return ref
	case *ast.SelectorExpr:
		if alias, ok := x.X.(*ast.Ident); ok {
			if imported := imports[alias.Name]; imported != "" {
				return reference{pkg: strings.TrimPrefix(imported, module), name: x.Sel.Name}
			}
		}
	}
	return reference{}
}

func assignedType(values []ast.Expr, index int, pkg string, imports map[string]string, depth int) reference {
	if len(values) == 1 && index > 0 {
		ref := valueType(values[0], pkg, imports, depth+1)
		if ref.factory {
			ref.resultIndex = index
			return ref
		}
		return reference{}
	}
	if index < len(values) {
		return valueType(values[index], pkg, imports, depth+1)
	}
	return reference{}
}

func configureInterface(ref reference) bool {
	return ref.pkg == framework+"resource" && ref.name == "ResourceWithConfigure" ||
		ref.pkg == framework+"datasource" && ref.name == "DataSourceWithConfigure"
}

func frameworkRequestMethod(ref reference) bool {
	typ := reference{pkg: ref.pkg, name: ref.receiver}
	if ref.name == "Configure" {
		return configureInterface(typ)
	}
	return ref.name == "Schema" && (typ.pkg == framework+"resource" && typ.name == "Resource" ||
		typ.pkg == framework+"datasource" && typ.name == "DataSource" || typ.pkg == framework+"provider" && typ.name == "Provider")
}

func valueType(expr ast.Expr, pkg string, imports map[string]string, depth int) reference {
	if depth > 8 {
		return reference{}
	}
	switch x := expr.(type) {
	case *ast.CompositeLit:
		return typeReference(x.Type, pkg, imports)
	case *ast.TypeAssertExpr:
		typ := typeReference(x.Type, pkg, imports)
		if (typ.pointer && typ.name != "") || configureInterface(typ) {
			return typ
		}
	case *ast.ParenExpr:
		return valueType(x.X, pkg, imports, depth+1)
	case *ast.UnaryExpr:
		return valueType(x.X, pkg, imports, depth+1)
	case *ast.CallExpr:
		switch fn := x.Fun.(type) {
		case *ast.Ident:
			if fn.Obj != nil {
				switch fn.Obj.Decl.(type) {
				case *ast.FuncDecl, *ast.ValueSpec:
				default:
					return reference{}
				}
			}
			if fn.Obj == nil && fn.Name == "new" && len(x.Args) == 1 {
				return typeReference(x.Args[0], pkg, imports)
			}
			return reference{pkg: pkg, name: fn.Name, factory: true}
		case *ast.SelectorExpr:
			if fn.Sel.Name == "RootModule" && len(x.Args) == 0 {
				state := valueType(fn.X, pkg, imports, depth+1)
				if !state.factory && state.field == "" && state.name == "State" && (state.pkg == "github.com/hashicorp/terraform-plugin-testing/terraform" ||
					state.pkg == "github.com/hashicorp/terraform-plugin-sdk/v2/terraform") {
					return reference{pkg: state.pkg, name: "ModuleState"}
				}
			}
			if alias, ok := fn.X.(*ast.Ident); ok && alias.Obj == nil {
				if imported := imports[alias.Name]; strings.HasPrefix(imported, module+"alicloud/") {
					return reference{pkg: strings.TrimPrefix(imported, module), name: fn.Sel.Name, factory: true}
				}
			}
		}
	case *ast.SelectorExpr:
		target := valueType(x.X, pkg, imports, depth+1)
		if !target.factory && target.name != "" && target.field == "" {
			target.field = x.Sel.Name
			return target
		}
	case *ast.Ident:
		if x.Obj != nil {
			switch decl := x.Obj.Decl.(type) {
			case *ast.Field:
				return typeReference(decl.Type, pkg, imports)
			case *ast.ValueSpec:
				if decl.Type != nil {
					return typeReference(decl.Type, pkg, imports)
				}
				for i, name := range decl.Names {
					if name.Name == x.Name {
						return assignedType(decl.Values, i, pkg, imports, depth)
					}
				}
			case *ast.AssignStmt:
				for i, name := range decl.Lhs {
					if id, ok := name.(*ast.Ident); ok && id.Name == x.Name {
						return assignedType(decl.Rhs, i, pkg, imports, depth)
					}
				}
			}
		}
	}
	return reference{}
}

// Read the static registration slices and leave other ServicePackage logic in
// the declaration fingerprint so mixed edits cannot pass as mapping-only changes.
func normalizeRegistration(fn *ast.FuncDecl, imports map[string]string) (map[string]string, string, bool) {
	if fn.Body == nil || len(fn.Body.List) != 1 {
		return nil, "return expression", false
	}
	ret, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return nil, "return expression", false
	}
	body, ok := ret.Results[0].(*ast.CompositeLit)
	if !ok {
		return nil, "return expression", false
	}
	mappings := map[string]string{}
	var remaining []ast.Expr
	literalName := false
	for _, entry := range body.Elts {
		field, ok := entry.(*ast.KeyValueExpr)
		if !ok {
			return nil, "return expression", false
		}
		kind := receiverName(field.Key)
		switch kind {
		case "Resources", "SDKResources", "DataSources", "SDKDataSources", "EphemeralResources", "Actions", "ListResources":
		default:
			if kind == "Name" {
				lit, ok := field.Value.(*ast.BasicLit)
				literalName = ok && lit.Kind == token.STRING
			}
			remaining = append(remaining, entry)
			continue
		}
		list, ok := field.Value.(*ast.CompositeLit)
		if !ok {
			return nil, kind, false
		}
		slice, ok := list.Type.(*ast.ArrayType)
		if !ok || slice.Len != nil {
			return nil, kind, false
		}
		elem, ok := slice.Elt.(*ast.SelectorExpr)
		owner, qualified := body.Type.(*ast.SelectorExpr)
		if !ok || !qualified || owner.Sel.Name != "ServicePackage" || receiverName(elem.X) != receiverName(owner.X) || elem.Sel.Name != strings.TrimSuffix(kind, "s") {
			return nil, kind, false
		}
		for _, element := range list.Elts {
			item, ok := element.(*ast.CompositeLit)
			if !ok || item.Type != nil || len(item.Elts) != 2 {
				return nil, kind, false
			}
			var name, factory string
			for _, attr := range item.Elts {
				pair, ok := attr.(*ast.KeyValueExpr)
				if !ok {
					return nil, kind, false
				}
				switch receiverName(pair.Key) {
				case "TypeName":
					if lit, ok := pair.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						name, _ = strconv.Unquote(lit.Value)
					}
				case "Factory":
					if ident, ok := pair.Value.(*ast.Ident); ok {
						factory = ident.Name
					}
				}
			}
			key := kind + ":" + name
			if name == "" || factory == "" || mappings[key] != "" {
				return nil, kind, false
			}
			mappings[key] = factory
		}
	}
	body.Elts = remaining
	// A new ServicePackage has no previous fingerprint. Admit it only when the
	// factory mappings and a literal product name account for the entire body.
	owner, ok := body.Type.(*ast.SelectorExpr)
	static := len(mappings) > 0 && len(remaining) == 1 && literalName && ok &&
		imports[receiverName(owner.X)] == module+"alicloud/provider/conns"
	return mappings, "", static
}

// Match direct SDK assertions on literal addresses only. A match is a repair hint, never coverage.
func checkedType(call *ast.CallExpr, imports map[string]string) string {
	fn, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	pkg := imports[receiverName(fn.X)]
	if pkg != "github.com/hashicorp/terraform-plugin-testing/helper/resource" && pkg != "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource" {
		return ""
	}
	switch fn.Sel.Name {
	case "TestCheckResourceAttr", "TestCheckResourceAttrSet", "TestCheckResourceAttrPair", "TestMatchResourceAttr":
	default:
		return ""
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	addr, _ := strconv.Unquote(lit.Value)
	parts := strings.Split(addr, ".")
	if len(parts) == 3 && parts[0] == "data" {
		return "data." + parts[1]
	}
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}

func methodSignature(fset *token.FileSet, fn *ast.FuncType) string {
	var parts []string
	for _, fields := range []*ast.FieldList{fn.Params, fn.Results} {
		var types []string
		if fields != nil {
			for _, f := range fields.List {
				for count := max(1, len(f.Names)); count > 0; count-- {
					types = append(types, nodeText(fset, f.Type))
				}
			}
		}
		parts = append(parts, "("+strings.Join(types, ",")+")")
	}
	return strings.Join(parts, "")
}

func stringLiteralValue(expr ast.Expr, depth int) string {
	if depth > 8 {
		return ""
	}
	switch x := expr.(type) {
	case *ast.BasicLit:
		if x.Kind == token.STRING {
			value, _ := strconv.Unquote(x.Value)
			return value
		}
	case *ast.Ident:
		if x.Obj == nil {
			return ""
		}
		switch decl := x.Obj.Decl.(type) {
		case *ast.ValueSpec:
			for i, name := range decl.Names {
				if name.Name == x.Name && i < len(decl.Values) {
					return stringLiteralValue(decl.Values[i], depth+1)
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range decl.Lhs {
				if id, ok := lhs.(*ast.Ident); ok && id.Name == x.Name && i < len(decl.Rhs) {
					return stringLiteralValue(decl.Rhs[i], depth+1)
				}
			}
		}
	}
	return ""
}

func returnedServiceType(expr ast.Expr, pkg string, imports map[string]string, depth int) reference {
	if depth > 8 {
		return reference{}
	}
	switch x := expr.(type) {
	case *ast.FuncLit:
		if x.Body == nil {
			return reference{}
		}
		var found reference
		valid := true
		ast.Inspect(x.Body, func(n ast.Node) bool {
			if !valid || n == nil {
				return false
			}
			if _, nested := n.(*ast.FuncLit); nested {
				return false
			}
			ret, ok := n.(*ast.ReturnStmt)
			if !ok {
				return true
			}
			if len(ret.Results) != 1 {
				valid = false
				return false
			}
			ref := valueType(ret.Results[0], pkg, imports, depth+1)
			if ref.name == "" {
				valid = false
				return false
			}
			if found.name != "" && found != ref {
				valid = false
				return false
			}
			found = ref
			return false
		})
		if valid {
			return found
		}
	case *ast.Ident:
		if x.Obj == nil {
			return reference{}
		}
		switch decl := x.Obj.Decl.(type) {
		case *ast.ValueSpec:
			for i, name := range decl.Names {
				if name.Name == x.Name && i < len(decl.Values) {
					return returnedServiceType(decl.Values[i], pkg, imports, depth+1)
				}
			}
		case *ast.AssignStmt:
			for i, lhs := range decl.Lhs {
				if id, ok := lhs.(*ast.Ident); ok && id.Name == x.Name && i < len(decl.Rhs) {
					return returnedServiceType(decl.Rhs[i], pkg, imports, depth+1)
				}
			}
		}
	}
	return reference{}
}

func describeMethodForResourceID(resourceID string) string {
	start := strings.Index(resourceID, "alicloud_")
	if start < 0 {
		return ""
	}
	start += len("alicloud_")
	dot := strings.Index(resourceID[start:], ".")
	if dot <= 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Describe")
	for _, part := range strings.Split(resourceID[start:start+dot], "_") {
		if len(part) > 0 {
			b.WriteString(strings.ToUpper(part[:1]))
			b.WriteString(part[1:])
		}
	}
	return b.String()
}

func isStringSliceLiteral(expr ast.Expr) bool {
	list, ok := expr.(*ast.CompositeLit)
	if !ok {
		return false
	}
	slice, ok := list.Type.(*ast.ArrayType)
	if !ok || slice.Len != nil {
		return false
	}
	elt, ok := slice.Elt.(*ast.Ident)
	if !ok || elt.Name != "string" {
		return false
	}
	for _, item := range list.Elts {
		lit, ok := item.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return false
		}
	}
	return true
}

// A literal SDK sweeper registration only installs a callback for its named
// resource. Other initializer side effects must keep package-wide coverage.
func registeredSweeper(fn *ast.FuncDecl, imports map[string]string) string {
	if fn.Body == nil || len(fn.Body.List) != 1 {
		return ""
	}
	stmt, ok := fn.Body.List[0].(*ast.ExprStmt)
	if !ok {
		return ""
	}
	call, ok := stmt.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 2 {
		return ""
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "AddTestSweepers" {
		return ""
	}
	alias, ok := sel.X.(*ast.Ident)
	if !ok || alias.Obj != nil || (imports[alias.Name] != "github.com/hashicorp/terraform-plugin-testing/helper/resource" &&
		imports[alias.Name] != "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource") {
		return ""
	}
	name, ok := call.Args[0].(*ast.BasicLit)
	if !ok || name.Kind != token.STRING {
		return ""
	}
	address, _ := strconv.Unquote(name.Value)
	ptr, ok := call.Args[1].(*ast.UnaryExpr)
	if !ok || ptr.Op != token.AND {
		return ""
	}
	spec, ok := ptr.X.(*ast.CompositeLit)
	if !ok || len(spec.Elts) < 2 || len(spec.Elts) > 3 {
		return ""
	}
	typ, ok := spec.Type.(*ast.SelectorExpr)
	if !ok || receiverName(typ.X) != alias.Name || typ.Sel.Name != "Sweeper" {
		return ""
	}
	hasName, hasCallback, hasDeps := false, false, false
	for _, elt := range spec.Elts {
		field, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			return ""
		}
		switch receiverName(field.Key) {
		case "Name":
			lit, ok := field.Value.(*ast.BasicLit)
			if hasName || !ok || lit.Value != name.Value {
				return ""
			}
			hasName = true
		case "F":
			if _, ok := field.Value.(*ast.Ident); hasCallback || !ok {
				return ""
			}
			hasCallback = true
		case "Dependencies":
			if hasDeps || !isStringSliceLiteral(field.Value) {
				return ""
			}
			hasDeps = true
		default:
			return ""
		}
	}
	if hasName && hasCallback {
		return address
	}
	return ""
}

func parseSource(s *snapshot, file string, data []byte) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, data, 0)
	if err != nil {
		return err
	}
	pkg := path.Dir(file) + "#" + f.Name.Name
	if s.packages == nil {
		s.packages = map[string]string{}
	}
	s.packages[file] = f.Name.Name
	imports := map[string]string{}
	for _, imp := range f.Imports {
		value, _ := strconv.Unquote(imp.Path.Value)
		alias := path.Base(value)
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		imports[alias] = value
		if strings.HasPrefix(value, module) {
			if s.importers == nil {
				s.importers = map[string]map[string]bool{}
			}
			dependency := strings.TrimPrefix(value, module)
			if s.importers[dependency] == nil {
				s.importers[dependency] = map[string]bool{}
			}
			s.importers[dependency][file] = true
		}
	}
	directives := ""
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimLeft(line, " \t")
		if strings.HasPrefix(line, "//go:") || strings.HasPrefix(line, "// +build ") {
			directives += line + "\n"
		}
	}
	for _, alias := range []string{".", "_"} {
		if imported := imports[alias]; imported != "" {
			directives += alias + ":" + imported + "\n"
		}
	}
	if s.hazards == nil {
		s.hazards = map[string]string{}
	}
	if directives != "" {
		s.hazards[file] = directives
	}
	top := map[ast.Node]bool{}
	for _, decl := range f.Decls {
		top[decl] = true
		if gen, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				top[spec] = true
			}
		}
	}
	initIndex := 0
	add := func(name, receiver string, node ast.Node) *declaration {
		key := name
		if receiver != "" {
			key = receiver + "." + name
		}
		if name == "init" && receiver == "" {
			key += "@" + strconv.Itoa(initIndex)
			initIndex++
		}
		d := &declaration{id: pkg + ":" + key + "@" + file, pkg: pkg, name: name, receiver: receiver,
			file: file, test: strings.HasSuffix(file, "_test.go")}
		var factory *ast.FuncType
		switch decl := node.(type) {
		case *ast.FuncDecl:
			d.isFunc = true
			if d.test && receiver == "" && name == "init" {
				d.sweeperAddress = registeredSweeper(decl, imports)
			}
			if decl.Recv == nil {
				factory = decl.Type
			}
			d.privateTestFunc = d.test && decl.Recv == nil && decl.Body != nil &&
				!ast.IsExported(name) && name != "init" && name != "_"
			if receiver != "" {
				d.signature = methodSignature(fset, decl.Type)
				if (name == "Configure" || name == "Schema") && len(decl.Type.Params.List) == 3 {
					req := typeReference(decl.Type.Params.List[1].Type, pkg, imports)
					if req.name == name+"Request" && strings.HasPrefix(req.pkg, framework) {
						d.requestPackage = req.pkg
					}
				}
			}
		case *ast.TypeSpec:
			switch typ := decl.Type.(type) {
			case *ast.InterfaceType:
				d.interfaceMethods = []string{}
				for _, field := range typ.Methods.List {
					if fn, ok := field.Type.(*ast.FuncType); ok && len(field.Names) == 1 {
						d.interfaceMethods = append(d.interfaceMethods, field.Names[0].Name+methodSignature(fset, fn))
					} else if embedded, ok := field.Type.(*ast.Ident); ok {
						d.interfaceMethods = append(d.interfaceMethods, "@"+embedded.Name)
					} else {
						d.interfaceMethods = nil
						break
					}
				}
			case *ast.StructType:
				d.structFields = map[string]reference{}
				for _, field := range typ.Fields.List {
					// Embedding the base Framework Provider promotes only its required methods.
					if len(field.Names) == 0 {
						embed := typeReference(field.Type, pkg, imports)
						if embed.pkg != framework+"provider" || embed.name != "Provider" {
							d.mayPromoteOptionalCallbacks = true
						}
					}
					for _, name := range field.Names {
						d.structFields[name.Name] = typeReference(field.Type, pkg, imports)
					}
				}
			}
		case *ast.ValueSpec:
			for i, ident := range decl.Names {
				if ident.Name == name && i < len(decl.Values) {
					switch val := decl.Values[i].(type) {
					case *ast.FuncLit:
						factory = val.Type
					case *ast.SelectorExpr:
						if alias, ok := val.X.(*ast.Ident); ok && alias.Obj == nil {
							if imported := imports[alias.Name]; strings.HasPrefix(imported, module+"alicloud") {
								d.factoryAlias = reference{pkg: strings.TrimPrefix(imported, module), name: val.Sel.Name}
							}
						}
					case *ast.Ident:
						if val.Obj == nil || top[val.Obj.Decl.(ast.Node)] {
							d.factoryAlias = reference{pkg: pkg, name: val.Name}
						}
					}
				}
			}
		}
		if factory != nil && factory.Results != nil {
			for _, field := range factory.Results.List {
				for i := 0; i < max(1, len(field.Names)); i++ {
					d.results = append(d.results, typeReference(field.Type, pkg, imports))
				}
			}
		}
		if fn, ok := node.(*ast.FuncDecl); ok && d.test && fn.Recv == nil && fn.Type.Results == nil && fn.Type.Params.NumFields() == 1 {
			// Include candidate AccTests for runner validation. Unit-only evidence must
			// match go test discovery, which accepts *T and *pkg.T, not an arbitrary alias.
			d.acc = strings.HasPrefix(name, "TestAcc")
			if suffix, test := strings.CutPrefix(name, "Test"); test && !d.acc && fn.Type.TypeParams.NumFields() == 0 {
				first, _ := utf8.DecodeRuneInString(suffix)
				if pointer, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr); ok && !unicode.IsLower(first) {
					switch typ := pointer.X.(type) {
					case *ast.Ident:
						d.unit = typ.Name == "T"
					case *ast.SelectorExpr:
						d.unit = typ.Sel.Name == "T"
					}
				}
			}
		}
		usedImports := map[string]string{}
		var visit func(ast.Node) bool
		visit = func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				if ident, ok := x.X.(*ast.Ident); ok {
					if imported, found := imports[ident.Name]; found {
						usedImports[ident.Name] = imported
						if strings.HasPrefix(imported, module+"alicloud") {
							d.refs = append(d.refs, reference{pkg: strings.TrimPrefix(imported, module), name: x.Sel.Name})
						}
						return false
					}
				}
				target := valueType(x.X, pkg, imports, 0)
				if x.Sel.Name != "Configure" && configureInterface(target) {
					target = reference{}
				}
				if target.name == "" {
					target.pkg = pkg
				}
				d.refs = append(d.refs, reference{pkg: target.pkg, name: x.Sel.Name,
					receiver: target.name, method: true, factory: target.factory, resultIndex: target.resultIndex, field: target.field, line: fset.Position(x.Sel.Pos()).Line})
				// A chained receiver can contain calls with their own dependencies.
				// Visit only X: Sel names a member, not a package-level declaration.
				ast.Inspect(x.X, visit)
				return false
			case *ast.Ident:
				if x.Obj == nil || top[x.Obj.Decl.(ast.Node)] {
					d.refs = append(d.refs, reference{pkg: pkg, name: x.Name})
					if imported := imports["."]; x.Obj == nil && strings.HasPrefix(imported, module+"alicloud") {
						d.refs = append(d.refs, reference{pkg: strings.TrimPrefix(imported, module), name: x.Name})
					}
				}
			case *ast.BasicLit:
				if x.Kind == token.STRING {
					value, _ := strconv.Unquote(x.Value)
					d.assets = append(d.assets, path.Clean(path.Join(path.Dir(file), value)))
				}
			case *ast.CallExpr:
				if d.acc && len(x.Args) > 0 {
					if address := checkedType(x, imports); address != "" {
						d.assertedTypes = append(d.assertedTypes, address)
					}
				}
				if fn, ok := x.Fun.(*ast.SelectorExpr); ok && len(x.Args) == 1 && d.test {
					if alias, ok := fn.X.(*ast.Ident); ok && alias.Obj == nil && imports[alias.Name] == framework+"providerserver" {
						switch fn.Sel.Name {
						case "NewProtocol5", "NewProtocol6", "NewProtocol5WithError", "NewProtocol6WithError":
							provider := valueType(x.Args[0], pkg, imports, 0)
							addCallback := func(method string, optional bool) {
								d.refs = append(d.refs, reference{pkg: provider.pkg, name: method,
									receiver: provider.name, method: true, optionalCallback: optional, factory: provider.factory, field: provider.field,
									resultIndex: provider.resultIndex, line: fset.Position(fn.Pos()).Line})
							}
							for _, method := range []string{"Metadata", "Schema", "Configure", "Resources", "DataSources"} {
								addCallback(method, false)
							}
							for _, method := range []string{"Actions", "ConfigValidators", "EphemeralResources", "Functions", "ListResources", "MetaSchema", "StateStores", "ValidateConfig"} {
								addCallback(method, true)
							}
						}
					}
					switch fn.Sel.Name {
					case "NewServer", "NewTLSServer", "NewUnstartedServer":
						if alias, ok := fn.X.(*ast.Ident); ok && alias.Obj == nil && imports[alias.Name] == "net/http/httptest" {
							handler := valueType(x.Args[0], pkg, imports, 0)
							if handler.name != "" {
								// httptest accepts http.Handler and invokes ServeHTTP through the interface.
								d.refs = append(d.refs, reference{pkg: handler.pkg, name: "ServeHTTP",
									receiver: handler.name, method: true, factory: handler.factory,
									field: handler.field, resultIndex: handler.resultIndex, line: fset.Position(fn.Pos()).Line})
							}
						}
					}
				}
				if d.test {
					callee := ""
					switch fn := x.Fun.(type) {
					case *ast.Ident:
						if fn.Name == "resourceCheckInitWithDescribeMethod" || fn.Name == "resourceCheckInit" {
							callee = fn.Name
						}
					case *ast.SelectorExpr:
						if alias, ok := fn.X.(*ast.Ident); ok && alias.Obj == nil && imports[alias.Name] == module+"alicloud/acctest" {
							if fn.Sel.Name == "ResourceCheckInitWithDescribeMethod" || fn.Sel.Name == "ResourceCheckInit" {
								callee = fn.Sel.Name
							}
						}
					}
					method := ""
					if (callee == "resourceCheckInitWithDescribeMethod" || callee == "ResourceCheckInitWithDescribeMethod") && len(x.Args) >= 4 {
						method = stringLiteralValue(x.Args[3], 0)
					} else if (callee == "resourceCheckInit" || callee == "ResourceCheckInit") && len(x.Args) >= 3 {
						method = describeMethodForResourceID(stringLiteralValue(x.Args[0], 0))
					}
					if method != "" {
						svc := returnedServiceType(x.Args[2], pkg, imports, 0)
						if svc.name != "" {
							if svc.pkg == "" {
								svc.pkg = pkg
							}
							d.refs = append(d.refs, reference{pkg: svc.pkg, name: method,
								receiver: svc.name, method: true, factory: svc.factory,
								field: svc.field, resultIndex: svc.resultIndex, line: fset.Position(x.Pos()).Line})
						}
					}
				}
				if fn, ok := x.Fun.(*ast.SelectorExpr); ok && fn.Sel.Name == "Join" {
					var parts []string
					for _, arg := range x.Args {
						literal, ok := arg.(*ast.BasicLit)
						if !ok || literal.Kind != token.STRING {
							parts = nil
							break
						}
						value, _ := strconv.Unquote(literal.Value)
						parts = append(parts, value)
					}
					if len(parts) > 0 {
						d.assets = append(d.assets, path.Clean(path.Join(path.Dir(file), path.Join(parts...))))
					}
				}
			}
			return true
		}
		ast.Inspect(node, visit)
		// Import changes can rebind an unchanged selector; include that binding in the fingerprint.
		bindings, _ := json.Marshal(usedImports)
		d.hash = fmt.Sprintf("%x", sha256.Sum256([]byte(nodeText(fset, node)+directives+string(bindings))))
		s.decls[d.id] = d
		return d
	}
	for _, decl := range f.Decls {
		switch x := decl.(type) {
		case *ast.FuncDecl:
			receiver := ""
			if x.Recv != nil {
				receiver = receiverName(x.Recv.List[0].Type)
			}
			staticRegistration := false
			if path.Base(file) == "register.go" && x.Name.Name == "ServicePackage" && receiver == "" {
				if mappings, field, static := normalizeRegistration(x, imports); mappings != nil {
					if s.registrations == nil {
						s.registrations = map[string]map[string]string{}
					}
					s.registrations[file] = mappings
					staticRegistration = static
				} else {
					if s.registrationIssues == nil {
						s.registrationIssues = map[string]string{}
					}
					s.registrationIssues[file] = field
				}
			}
			add(x.Name.Name, receiver, x).staticRegistration = staticRegistration
		case *ast.GenDecl:
			inherited := false
			if x.Tok == token.CONST {
				ast.Inspect(x, func(n ast.Node) bool {
					if id, ok := n.(*ast.Ident); ok && id.Name == "iota" {
						inherited = true
					}
					if value, ok := n.(*ast.ValueSpec); ok && len(value.Values) == 0 {
						inherited = true
					}
					return true
				})
			}
			for _, spec := range x.Specs {
				switch value := spec.(type) {
				case *ast.ValueSpec:
					for _, name := range value.Names {
						if inherited {
							add(name.Name, "", x)
						} else {
							add(name.Name, "", value)
						}
					}
				case *ast.TypeSpec:
					add(value.Name.Name, "", value)
				}
			}
		}
	}
	return nil
}

func readSnapshot(ref string) (*snapshot, error) {
	s := &snapshot{decls: map[string]*declaration{}, files: map[string]bool{}, assets: map[string]string{}}
	cmd := exec.Command("git", "archive", ref, "alicloud/")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	reader := tar.NewReader(pipe)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			return nil, err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		data, err := io.ReadAll(reader)
		if err == nil {
			s.files[header.Name] = true
			if strings.HasSuffix(header.Name, ".go") {
				err = parseSource(s, header.Name, data)
			} else {
				s.assets[header.Name] = fmt.Sprintf("%x", sha256.Sum256(data))
			}
		}
		if err != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			return nil, err
		}
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("git archive %s: %w: %s", ref, err, strings.TrimSpace(stderr.String()))
	}
	return s, nil
}

func resourceFile(file string, registered map[string]bool) bool {
	if registered[file] {
		return true
	}
	if strings.HasSuffix(file, "_test.go") {
		return false
	}
	// In a registered service package only actual factory files are entrypoints.
	// A helper's spelling must never manufacture another Terraform entity.
	if registered[path.Dir(file)+"/"] {
		return false
	}
	for _, prefix := range []string{"resource_", "data_source_alicloud_", "ephemeral_", "action_", "list_"} {
		if strings.HasPrefix(path.Base(file), prefix) {
			return true
		}
	}
	return false
}

func service(d *declaration) bool {
	name := path.Base(d.file)
	return strings.HasPrefix(name, "service_") || name == "service.go" || strings.HasSuffix(name, "_service.go") ||
		(strings.HasPrefix(d.file, "alicloud/service/") && strings.HasSuffix(d.receiver, "Service"))
}

type edge struct {
	caller, ambiguity string
}

func reverseGraph(s *snapshot) map[string][]edge {
	index := map[string][]string{}
	methods := map[string]map[string]bool{}
	for id, d := range s.decls {
		// Go package initializers cannot be referenced by name. Keep their
		// outgoing dependencies without creating edges between distinct inits.
		if d.name == "init" && d.receiver == "" {
			continue
		}
		key := d.pkg + ":" + d.name
		if d.receiver != "" {
			receiver := d.pkg + ":" + d.receiver
			if methods[receiver] == nil {
				methods[receiver] = map[string]bool{}
			}
			methods[receiver][d.name+d.signature] = true
			key = d.pkg + ":*." + d.name
			if !d.test {
				index["method:"+d.name] = append(index["method:"+d.name], id)
			} else {
				key := "test-method:" + path.Dir(d.file) + ":" + d.name
				index[key] = append(index[key], id)
				index["test-method:*:"+d.name] = append(index["test-method:*:"+d.name], id)
			}
			index[d.pkg+":"+d.receiver+"."+d.name] = append(index[d.pkg+":"+d.receiver+"."+d.name], id)
			if !d.test {
				qualified := strings.Split(d.pkg, "#")[0] + ":" + d.receiver + "." + d.name
				index[qualified] = append(index[qualified], id)
			}
		}
		index[key] = append(index[key], id)
		if d.test && d.receiver == "" && ast.IsExported(d.name) && !strings.HasSuffix(d.pkg, "_test") {
			// go test makes exported internal-test declarations visible to
			// the external test package in this directory only.
			index["test-export:"+d.pkg+":"+d.name] = append(index["test-export:"+d.pkg+":"+d.name], id)
		}
		if d.receiver == "" && !d.test {
			index[strings.Split(d.pkg, "#")[0]+":"+d.name] = append(index[strings.Split(d.pkg, "#")[0]+":"+d.name], id)
		}
	}
	var interfaceMethods func(string, map[string]bool) ([]string, bool)
	interfaceMethods = func(key string, active map[string]bool) ([]string, bool) {
		ids := index[key]
		if len(ids) != 1 || active[key] || s.decls[ids[0]].interfaceMethods == nil {
			return nil, false
		}
		active[key] = true
		defer delete(active, key)
		d := s.decls[ids[0]]
		var signatures []string
		for _, method := range d.interfaceMethods {
			if strings.HasPrefix(method, "@") {
				embedded, ok := interfaceMethods(d.pkg+":"+method[1:], active)
				if !ok {
					return nil, false
				}
				signatures = append(signatures, embedded...)
			} else {
				signatures = append(signatures, method)
			}
		}
		return signatures, len(signatures) != 0
	}
	interfaces := map[string][]string{}
	for _, d := range s.decls {
		if d.interfaceMethods == nil {
			continue
		}
		key := d.pkg + ":" + d.name
		signatures, ok := interfaceMethods(key, map[string]bool{})
		if !ok {
			continue
		}
		for receiver, implemented := range methods {
			if !strings.HasPrefix(receiver, strings.Split(d.pkg, "#")[0]+"#") {
				continue
			}
			matches := true
			for _, signature := range signatures {
				if !implemented[signature] {
					matches = false
					break
				}
			}
			if matches {
				interfaces[key] = append(interfaces[key], receiver)
			}
		}
	}
	graph := map[string][]edge{}
	for id, d := range s.decls {
		for _, ref := range d.refs {
			externalPointer := false
			if ref.factory {
				targetPkg, targetName := ref.pkg, ref.receiver
				result := reference{}
				for depth := 0; depth < 5; depth++ {
					factories := index[targetPkg+":"+targetName]
					if len(factories) != 1 {
						break
					}
					factory := s.decls[factories[0]]
					if factory.factoryAlias.name != "" {
						targetPkg, targetName = factory.factoryAlias.pkg, factory.factoryAlias.name
						continue
					}
					if ref.resultIndex < len(factory.results) {
						result = factory.results[ref.resultIndex]
					}
					break
				}
				ref.pkg, ref.receiver = result.pkg, result.name
				// A pointer to an external named type cannot be an interface receiver.
				externalPointer = result.pointer && result.pkg != "" && result.pkg != "alicloud" &&
					!strings.HasPrefix(result.pkg, "alicloud/") && !strings.HasPrefix(result.pkg, "alicloud#")
				if ref.receiver == "" {
					ref.pkg = d.pkg
				}
			}
			if ref.field != "" {
				types := index[ref.pkg+":"+ref.receiver]
				resolved := reference{pkg: d.pkg}
				if len(types) == 1 {
					if field, ok := s.decls[types[0]].structFields[ref.field]; ok {
						resolved = field
					}
				}
				ref.pkg, ref.receiver = resolved.pkg, resolved.name
			}
			key := ref.pkg + ":" + ref.name
			if ref.method {
				key = ref.pkg + ":*." + ref.name
			}
			matches := index[key]
			if !ref.method && d.test && strings.HasSuffix(d.pkg, "_test") && ref.pkg == path.Dir(d.file) {
				matches = append(append([]string{}, matches...), index["test-export:"+strings.TrimSuffix(d.pkg, "_test")+":"+ref.name]...)
			}
			unresolved, interfaceCall := false, false
			if ref.method {
				testDir := "*"
				if d.test {
					testDir = path.Dir(d.file)
				}
				testMethods := index["test-method:"+testDir+":"+ref.name]
				if ref.receiver != "" {
					matches = index[ref.pkg+":"+ref.receiver+"."+ref.name]
					if frameworkRequestMethod(ref) {
						// Provider, resource and data source take distinct named request types.
						interfaceCall = true
						for _, candidates := range [][]string{index["method:"+ref.name], testMethods} {
							for _, target := range candidates {
								if kind := s.decls[target].requestPackage; kind == "" || kind == ref.pkg {
									matches = append(matches, target)
								}
							}
						}
					}
					if implementations := interfaces[ref.pkg+":"+ref.receiver]; len(implementations) > 0 {
						interfaceCall = true
						for _, receiver := range implementations {
							matches = append(matches, index[receiver+"."+ref.name]...)
						}
					}
					concreteStruct := false
					if len(matches) == 0 {
						types := index[ref.pkg+":"+ref.receiver]
						if len(types) == 1 {
							typ := s.decls[types[0]]
							concreteStruct = typ.structFields != nil && !typ.mayPromoteOptionalCallbacks
						}
					}
					if ref.optionalCallback && len(matches) == 0 {
						if concreteStruct {
							continue
						}
					} else if !ref.factory && len(matches) == 0 && (strings.HasPrefix(ref.pkg, "alicloud/") || strings.HasPrefix(ref.pkg, "alicloud#")) {
						if !concreteStruct {
							matches = index["method:"+ref.name]
							unresolved = true
						}
					}
					if ref.factory && len(matches) == 0 && concreteStruct {
						continue
					}
				} else if !ref.optionalCallback && !ref.factory {
					candidates := map[string]bool{}
					for _, target := range append(append([]string{}, index["method:"+ref.name]...), index[key]...) {
						candidates[target] = true
					}
					matches = sortedKeys(candidates)
					unresolved = true
				}
				// An unindexed factory result may be an interface implemented by this repo.
				dynamic := ref.optionalCallback || ref.factory && !externalPointer
				if dynamic && !interfaceCall && (ref.receiver == "" || len(matches) == 0) {
					// A production caller may receive test implementations from importing test binaries.
					matches = append([]string{}, index["method:"+ref.name]...)
					matches = append(matches, testMethods...)
					unresolved = true
				}
			}
			ambiguity := ""
			if ref.method && !interfaceCall && (len(matches) > 1 || unresolved) {
				candidates := map[string]bool{}
				for _, target := range matches {
					method := s.decls[target]
					candidates[fmt.Sprintf("%s.%s (%s)", method.receiver, method.name, method.file)] = true
				}
				reason := "receiver type could not be resolved"
				if ref.receiver != "" && !unresolved {
					reason = "multiple declarations match the receiver (check build variants)"
				}
				ambiguity = fmt.Sprintf("%s:%d: cannot determine AccTest impact of %s: %s; candidates: %s",
					d.file, ref.line, ref.name, reason, strings.Join(sortedKeys(candidates), ", "))
			} else if len(matches) > 1 && !interfaceCall {
				ambiguity = "multiple declarations match " + key + " in " + d.file + " (check build variants)"
			}
			for _, target := range matches {
				if target != id {
					graph[target] = append(graph[target], edge{caller: id, ambiguity: ambiguity})
				}
			}
		}
		for _, file := range d.assets {
			if _, found := s.assets[file]; found {
				graph["asset:"+file] = append(graph["asset:"+file], edge{caller: id})
			}
		}
	}
	return graph
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// A _test.go file is compiled only into its directory's test binary. A production
// file can reach other binaries through imports, including imports removed in head.
func possibleConsumers(before, after *snapshot, file string) (tests, production map[string]bool) {
	dir := path.Dir(file)
	tests, production = map[string]bool{dir: true}, map[string]bool{}
	if strings.HasSuffix(file, "_test.go") {
		return
	}
	queue := []string{dir}
	// Only alicloud/ was indexed. An in-module import outside that tree can
	// transitively import the changed package, so include its importers too.
	for _, snap := range []*snapshot{before, after} {
		for dependency := range snap.importers {
			if dependency != "alicloud" && !strings.HasPrefix(dependency, "alicloud/") {
				queue = append(queue, dependency)
			}
		}
	}
	for len(queue) > 0 {
		pkg := queue[0]
		queue = queue[1:]
		if production[pkg] {
			continue
		}
		production[pkg] = true
		for _, snap := range []*snapshot{before, after} {
			for importer := range snap.importers[pkg] {
				dir := path.Dir(importer)
				tests[dir] = true
				if !strings.HasSuffix(importer, "_test.go") {
					queue = append(queue, dir)
				}
			}
		}
	}
	return
}

func selectScope(before, after *snapshot, files []string) accScope {
	result := accScope{Packages: map[string][]string{}, ResourceFiles: []string{}, Reasons: map[string][]string{}, Blocked: []string{}}
	selected, resources, blocked := map[string]bool{}, map[string]bool{}, map[string]bool{}
	accOnly := map[string]bool{}
	blockAcc := func(reason string) {
		accOnly[reason] = true
		blocked[reason] = true
	}
	missingCoverage := map[string]map[string]bool{}
	directResources := map[string]bool{}
	// File moves do not remove a uniquely identified symbol. Preserve package,
	// receiver and declaration name; ambiguous moves still require review.
	relocated := map[string]*declaration{}
	afterSymbols := map[string][]*declaration{}
	for _, d := range after.decls {
		key := d.pkg + ":" + d.receiver + ":" + d.name
		afterSymbols[key] = append(afterSymbols[key], d)
	}
	for id, d := range before.decls {
		if !after.files[d.file] {
			matches := afterSymbols[d.pkg+":"+d.receiver+":"+d.name]
			if len(matches) == 1 && matches[0].file != d.file && matches[0].test == d.test && (!d.acc || matches[0].acc) {
				relocated[id] = matches[0]
			}
		}
	}
	movedFiles := map[string]string{}
	beforeEntries, afterEntries := factoryFiles(before), factoryFiles(after)
	beforeFactories := map[string]bool{}
	for file, mappings := range before.registrations {
		for _, factory := range mappings {
			beforeFactories[path.Dir(file)+":"+factory] = true
		}
	}
	moveTargets, unresolvedMove := map[string]map[string]bool{}, map[string]bool{}
	for id, d := range before.decls {
		if !after.files[d.file] && beforeEntries[d.file] && d.receiver == "" && beforeFactories[path.Dir(d.file)+":"+d.name] {
			if moveTargets[d.file] == nil {
				moveTargets[d.file] = map[string]bool{}
			}
			if target := relocated[id]; target != nil && afterEntries[target.file] {
				moveTargets[d.file][target.file] = true
			} else {
				unresolvedMove[d.file] = true
			}
		}
	}
	for file, targets := range moveTargets {
		if len(targets) == 1 && !unresolvedMove[file] {
			movedFiles[file] = sortedKeys(targets)[0]
		}
	}
	changedFiles := map[string]bool{}
	for _, file := range files {
		changedFiles[file] = true
		if file == "go.mod" || file == "go.sum" {
			blocked["manual ACC scope decision required for shared dependency change: "+file] = true
		}
		if before.hazards[file] != after.hazards[file] || strings.Contains(after.hazards[file], "//go:build") || strings.Contains(after.hazards[file], "// +build") {
			blocked["build/directive/side-effect import scope requires manual review: "+file] = true
		}
	}
	selectTest := func(d *declaration, reason string) {
		if d.test && d.acc {
			selected[d.id] = true
			key := path.Dir(d.file) + ":" + d.name
			result.Reasons[key] = append(result.Reasons[key], reason)
		}
	}
	companion := func(file, testFile string) bool {
		stem := strings.TrimSuffix(file, ".go")
		return testFile == stem+"_test.go" || (strings.HasPrefix(testFile, stem+"_") &&
			!after.files[strings.TrimSuffix(testFile, "_test.go")+".go"])
	}
	selectResource := func(file, reason string) {
		if moved := movedFiles[file]; moved != "" {
			file = moved
		}
		resources[file] = true
		found := false
		for _, d := range after.decls {
			if !d.test || !d.acc {
				continue
			}
			if companion(file, d.file) {
				selectTest(d, reason)
				found = true
			}
		}
		if !found {
			if missingCoverage[file] == nil {
				missingCoverage[file] = map[string]bool{}
			}
			missingCoverage[file][reason] = true
		}
	}
	registered := map[string]bool{}
	entryAddresses := map[string]map[string]bool{}
	for _, snap := range []*snapshot{before, after} {
		for file := range snap.files {
			if path.Base(file) == "register.go" {
				registered[path.Dir(file)+"/"] = true
			}
		}
		factories := map[string][]string{}
		for _, d := range snap.decls {
			if !d.test && d.receiver == "" && len(d.results) == 1 && d.results[0].name != "" {
				key := path.Dir(d.file) + ":" + d.name
				factories[key] = append(factories[key], d.file)
			}
		}
		for file, mappings := range snap.registrations {
			other := before.registrations[file]
			if snap == before {
				other = after.registrations[file]
			}
			for name, factory := range mappings {
				targets := factories[path.Dir(file)+":"+factory]
				if len(targets) == 1 {
					registered[targets[0]] = true
					kind, typeName, _ := strings.Cut(name, ":")
					if snap == after {
						if entryAddresses[targets[0]] == nil {
							entryAddresses[targets[0]] = map[string]bool{}
						}
						if strings.HasSuffix(kind, "DataSources") {
							typeName = "data." + typeName
						}
						entryAddresses[targets[0]][typeName] = true
					}
				}
				if !changedFiles[file] || other[name] == factory {
					continue
				}
				reason := file + ": registration " + name + " -> " + factory
				if len(targets) != 1 {
					blocked["cannot resolve AccTest scope for "+reason] = true
					continue
				}
				directResources[targets[0]] = true
				selectResource(targets[0], reason)
			}
		}
	}
	snapshots := []*snapshot{before, after}
	graphs := []map[string][]edge{reverseGraph(before), reverseGraph(after)}
	uncertain := map[string]map[string]bool{}
	expand := func(file, reason string) {
		if uncertain[file] == nil {
			uncertain[file] = map[string]bool{}
		}
		uncertain[file][reason] = true
	}
	packageWide := map[string]bool{}
	selectTestInit := func(d *declaration, reason string) {
		owner := strings.TrimSuffix(d.file, "_test.go") + ".go"
		if entryAddresses[owner][d.sweeperAddress] {
			selectResource(owner, reason)
			return
		}
		if d.sweeperAddress != "" && path.Dir(d.file) == "alicloud/sweep" {
			return
		}
		packageWide[path.Dir(d.file)] = true
		expand(d.file, "test package initialization depends on "+reason)
	}
	hasLinkname := false
	for _, snap := range []*snapshot{before, after} {
		for _, directives := range snap.hazards {
			hasLinkname = hasLinkname || strings.Contains(directives, "//go:linkname")
		}
	}
	structReceivers := [2]map[string]bool{{}, {}}
	for i, snap := range snapshots {
		for _, decl := range snap.decls {
			if !decl.test && decl.receiver == "" && decl.structFields != nil {
				structReceivers[i][decl.pkg+":"+decl.name] = true
			}
		}
	}
	externalInterfaceMethods := map[string]bool{
		"String": true, "GoString": true, "Format": true, "Error": true, "Unwrap": true,
		"Is": true, "As": true, "MarshalJSON": true, "UnmarshalJSON": true, "MarshalText": true,
		"UnmarshalText": true, "MarshalBinary": true, "UnmarshalBinary": true, "Close": true,
		"Read": true, "Write": true, "ReadFrom": true, "WriteTo": true, "Seek": true,
		"ServeHTTP": true, "Len": true, "Less": true, "Swap": true, "Metadata": true,
		"Schema": true, "Configure": true, "Create": true, "Update": true, "Delete": true,
		"ImportState": true, "ModifyPlan": true, "ValidateConfig": true, "UpgradeState": true,
		"MoveState": true, "ConfigValidators": true, "Resources": true, "DataSources": true,
		"EphemeralResources": true, "Actions": true, "Functions": true, "ListResources": true,
		"MetaSchema": true, "StateStores": true, "Equal": true, "Type": true,
		"Description": true, "MarkdownDescription": true,
	}
	isDeadServiceMethodNode := func(i int, d *declaration) bool {
		return !d.test && !resourceFile(d.file, registered) && service(d) && d.receiver != "" &&
			snapshots[i].hazards[d.file] == "" && !externalInterfaceMethods[d.name] &&
			structReceivers[i][d.pkg+":"+d.receiver]
	}
	isDeadTestFuncNode := func(i int, d *declaration) bool {
		return d.test && d.isFunc && d.receiver == "" && !d.acc && !d.unit &&
			d.name != "init" && d.name != "TestMain" &&
			!strings.HasPrefix(d.name, "Test") && !strings.HasPrefix(d.name, "Benchmark") &&
			!strings.HasPrefix(d.name, "Fuzz") && !strings.HasPrefix(d.name, "Example") &&
			snapshots[i].hazards[d.file] == ""
	}
	unreferencedServiceMethod := func(seed string) bool {
		if hasLinkname {
			return false
		}
		seen := map[string]bool{}
		queue := []string{seed}
		for len(queue) > 0 {
			id := queue[0]
			queue = queue[1:]
			if moved := relocated[id]; moved != nil {
				id = moved.id
			}
			if seen[id] {
				continue
			}
			seen[id] = true
			foundAny := false
			for i, snap := range snapshots {
				d := snap.decls[id]
				if d == nil {
					continue
				}
				foundAny = true
				if id == seed {
					if !isDeadServiceMethodNode(i, d) {
						return false
					}
				} else if !isDeadServiceMethodNode(i, d) && !isDeadTestFuncNode(i, d) {
					return false
				}
				for _, e := range graphs[i][id] {
					if e.ambiguity != "" {
						return false
					}
					queue = append(queue, e.caller)
				}
			}
			if !foundAny {
				return false
			}
		}
		return true
	}
	unreferenced := map[string]bool{}
	unusedTestFunc := func(old, current *declaration) bool {
		// Methods may be called through interfaces, and package-level values can run at init.
		if hasLinkname {
			return false
		}
		for i, d := range []*declaration{old, current} {
			if d == nil {
				continue
			}
			snap := before
			if i == 1 {
				snap = after
			}
			if !d.privateTestFunc || snap.hazards[d.file] != "" || len(graphs[i][d.id]) != 0 {
				return false
			}
		}
		return true
	}
	// Protocol servers invoke methods through interfaces. A private test
	// receiver constructed exclusively by standalone unit tests cannot be
	// delivered to an acceptance test. Type references in helpers, fixtures,
	// initializers, or tests invoked by other functions keep normal analysis.
	unitTestMethod := func(old, current *declaration) bool {
		if hasLinkname {
			return false
		}
		for i, d := range []*declaration{old, current} {
			if d == nil {
				continue
			}
			if !d.test || !strings.HasSuffix(d.pkg, "_test") || d.receiver == "" ||
				ast.IsExported(d.receiver) || snapshots[i].hazards[d.file] != "" {
				return false
			}
			// Only certify a receiver declared in the same test file.
			typ := snapshots[i].decls[d.pkg+":"+d.receiver+"@"+d.file]
			if typ == nil || typ.structFields == nil {
				return false
			}
			found := false
			for _, e := range graphs[i][typ.id] {
				use := snapshots[i].decls[e.caller]
				if use.receiver == d.receiver && use.pkg == d.pkg {
					continue
				}
				if e.ambiguity != "" || !use.unit || len(graphs[i][use.id]) != 0 {
					return false
				}
				found = true
			}
			if !found {
				return false
			}
		}
		return true
	}
	impact := func(seed, reason, sourceFile string) (bool, bool) {
		seen := map[string]bool{}
		queue := []string{seed}
		covered := false
		unit := false
		for len(queue) > 0 {
			id := queue[0]
			queue = queue[1:]
			if moved := relocated[id]; moved != nil {
				id = moved.id
			}
			if seen[id] {
				continue
			}
			seen[id] = true
			d := after.decls[id]
			if d == nil {
				d = before.decls[id]
			}
			if d != nil && d.test && d.receiver != "" && unitTestMethod(before.decls[id], after.decls[id]) {
				unit = true
				continue
			}
			if d != nil && d.test && (d.name == "init" || d.name == "TestMain") && d.receiver == "" {
				selectTestInit(d, reason)
				covered = true
				continue
			}
			// A removed test helper may have had only unit callers that
			// were removed in the same change. That evidence applies to
			// the deleted helper, never to surviving production code.
			if d != nil && d.unit && (after.decls[id] != nil ||
				(after.decls[seed] == nil && before.decls[seed] != nil && before.decls[seed].privateTestFunc)) {
				unit = true
			}
			if d := after.decls[id]; d != nil {
				if d.test && d.acc {
					selectTest(d, reason)
					covered = true
				}
				if resourceFile(d.file, registered) {
					selectResource(d.file, reason)
					covered = true
					// Provider factories register all resources; the implementation boundary is the type scope.
					continue
				}
			}
			for i, graph := range graphs {
				for _, e := range graph[id] {
					if e.ambiguity != "" {
						// A possible call from a unit test to a private test mock
						// is bounded by that unit function. Keep walking in case
						// an acceptance test explicitly invokes the unit function.
						target := snapshots[i].decls[id]
						caller := snapshots[i].decls[e.caller]
						if target != nil && target.test && target.receiver != "" && !ast.IsExported(target.name) &&
							caller.unit && caller.pkg == target.pkg {
							queue = append(queue, e.caller)
							continue
						}
						if sourceFile == "" {
							// Asset seeds have no Go declaration to bound consumers.
							blocked[e.ambiguity+"; manual impact review required"] = true
						} else {
							expand(sourceFile, e.ambiguity)
							covered = true
						}
						continue
					}
					queue = append(queue, e.caller)
				}
			}
		}
		return covered, unit
	}
	all := map[string]*declaration{}
	for id, d := range before.decls {
		all[id] = d
	}
	for id, d := range after.decls {
		all[id] = d
	}
	for id, d := range all {
		old, current := before.decls[id], after.decls[id]
		if !changedFiles[d.file] || old != nil && current != nil && old.hash == current.hash && old.file == current.file {
			continue
		}
		reason := d.file + ":" + d.name
		dir := path.Dir(d.file)
		if strings.HasPrefix(dir, "alicloud/") && !strings.HasPrefix(dir, "alicloud/service/") &&
			!strings.HasPrefix(dir, "alicloud/acctest") {
			continue
		}
		if d.test || strings.HasPrefix(d.file, "alicloud/acctest/") {
			if d.acc && current == nil && relocated[id] == nil {
				blockAcc("removed/renamed AccTest " + reason + "; confirm replacement coverage")
				continue
			}
			if !d.test && d.name == "init" {
				blockAcc("test-support package initialization changed: " + reason + "; manual importer impact review required")
			} else if d.name == "init" && d.receiver == "" {
				for _, init := range []*declaration{old, current} {
					if init != nil {
						selectTestInit(init, reason)
					}
				}
			} else if d.name == "TestMain" {
				packageWide[dir] = true
				for _, test := range after.decls {
					if path.Dir(test.file) == path.Dir(d.file) {
						selectTest(test, reason)
					}
				}
			} else if !unusedTestFunc(old, current) {
				acc, unit := impact(id, reason, d.file)
				if !acc && !unit && !strings.HasPrefix(d.name, "Test") && !strings.HasPrefix(d.name, "Benchmark") && !strings.HasPrefix(d.name, "Fuzz") {
					expand(d.file, "cannot resolve test helper consumers for "+reason)
				}
			}
		} else if resourceFile(d.file, registered) {
			directResources[d.file] = true
			selectResource(d.file, reason)
		} else if service(d) {
			if covered, unit := impact(id, reason, d.file); !covered && !unit {
				if unreferencedServiceMethod(id) {
					unreferenced[reason] = true
				} else {
					expand(d.file, "cannot resolve service impact for "+reason)
				}
			}
		} else if path.Base(d.file) == "register.go" {
			if old == nil && current != nil && current.staticRegistration && after.hazards[d.file] == "" {
				continue // Each new mapping was already checked and its factory's tests selected.
			}
			if field := after.registrationIssues[d.file]; field != "" {
				reason += " (" + field + ")"
			}
			blocked["unsupported registration change: "+reason+"; manual impact review required"] = true
		} else if strings.HasPrefix(dir, "alicloud/service/") {
			// Do not turn unknown helper consumers into a successful empty scope.
			if covered, unit := impact(id, reason, d.file); !covered && !unit {
				expand(d.file, "cannot resolve production helper consumers for "+reason)
			}
		}
	}
	for file := range changedFiles {
		if !strings.HasPrefix(file, "alicloud/") || strings.HasSuffix(file, ".go") {
			continue
		}
		if before.assets[file] != after.assets[file] {
			acc, unit := impact("asset:"+file, file, "")
			if !acc && !unit && !strings.HasSuffix(file, ".md") && !strings.HasSuffix(file, ".markdown") {
				blocked["cannot resolve test fixture consumers for "+file+"; manual impact review required"] = true
			}
		}
	}
	byPackage := map[string]map[string]bool{}
	// Expand only selection uncertainty. Missing/deleted coverage, layout and
	// build constraints remain independent blockers, even with a nonempty suite.
	if len(uncertain) > 0 {
		result.FullPackages = map[string][]string{}
		fullReasons := map[string]map[string]bool{}
		cases := map[string][]*declaration{}
		for _, d := range after.decls {
			if d.test && d.acc {
				dir := path.Dir(d.file)
				cases[dir] = append(cases[dir], d)
			}
		}
		// Capture proven resource scope before considering any fallback.
		// An earlier speculative expansion must never validate a later one.
		owners := map[string]int{}
		unproven := map[string][]string{}
		for file := range after.files {
			if !resourceFile(file, registered) {
				continue
			}
			dir := path.Dir(file)
			for _, test := range cases[dir] {
				if companion(file, test.file) {
					owners[dir]++
					if !resources[file] {
						unproven[dir] = append(unproven[dir], file)
					}
					break
				}
			}
		}
		for file, reasons := range uncertain {
			block := func(message string) {
				if strings.HasSuffix(file, "_test.go") {
					blockAcc(message)
				} else {
					blocked[message] = true
				}
			}
			tests, production := possibleConsumers(before, after, file)
			services := map[string]bool{}
			for dir := range tests {
				if strings.HasPrefix(dir, "alicloud/service/") && len(cases[dir]) != 0 {
					services[dir] = true
				}
			}
			if len(services) > 1 {
				block(fmt.Sprintf("%s: unresolved impact would expand to %d service packages; resolve the caller or review the affected resources before dispatch: %s",
					file, len(services), strings.Join(sortedKeys(reasons), "; ")))
				continue
			}
			mixed := false
			for dir := range tests {
				if !packageWide[dir] && owners[dir] > 1 && len(unproven[dir]) != 0 {
					sort.Strings(unproven[dir])
					block(fmt.Sprintf("%s: unresolved impact would include unproven resource suites in %s; verify the consumers before dispatch (%d unproven, e.g. %s): %s",
						file, dir, len(unproven[dir]), unproven[dir][0], strings.Join(sortedKeys(reasons), "; ")))
					mixed = true
				}
			}
			if mixed {
				continue
			}
			found := false
			for dir := range tests {
				if len(cases[dir]) == 0 {
					continue
				}
				found = true
				if fullReasons[dir] == nil {
					fullReasons[dir] = map[string]bool{}
				}
				for reason := range reasons {
					fullReasons[dir][reason] = true
				}
				for _, test := range cases[dir] {
					selected[test.id] = true
				}
			}
			if !found || hasLinkname {
				for reason := range reasons {
					block(reason + "; manual impact review required")
				}
			}
			for f := range after.files {
				if production[path.Dir(f)] && resourceFile(f, registered) {
					resources[f] = true
				}
			}
		}
		for dir, reasons := range fullReasons {
			result.FullPackages[dir] = sortedKeys(reasons)
		}
	}
	for id := range selected {
		d := after.decls[id]
		pkg := path.Dir(d.file)
		if byPackage[pkg] == nil {
			byPackage[pkg] = map[string]bool{}
		}
		byPackage[pkg][d.name] = true
	}
	for pkg, cases := range byPackage {
		result.Packages[pkg] = sortedKeys(cases)
	}
	for key, reasons := range result.Reasons {
		unique := map[string]bool{}
		for _, reason := range reasons {
			unique[reason] = true
		}
		result.Reasons[key] = sortedKeys(unique)
	}
	coverageBlockers, unverifiedCoverage := map[string]bool{}, map[string]bool{}
	for file, reasons := range missingCoverage {
		if !directResources[file] {
			message := "Indirectly affected resource has no runnable AccTest: " + file +
				" (selected by " + strings.Join(sortedKeys(reasons), ", ") + ")"
			unverifiedCoverage[message] = true
			continue
		}
		message := "No AccTest cases for affected resource: " + file +
			" (selected by " + strings.Join(sortedKeys(reasons), ", ") + ")"
		if addresses := entryAddresses[file]; len(addresses) > 0 {
			expected := strings.TrimSuffix(file, ".go") + "_test.go"
			message = "No AccTest cases for affected resource: " + file + " (" + strings.Join(sortedKeys(addresses), ", ") + "). Fix: Create or update " + expected + " with a runnable TestAcc case."
			candidates := map[string]bool{}
			for _, test := range after.decls {
				if path.Dir(test.file) == path.Dir(file) {
					for _, address := range test.assertedTypes {
						if addresses[address] {
							candidates[test.file+":"+test.name] = true
						}
					}
				}
			}
			if len(candidates) > 0 {
				message += " Existing assertions found in " + strings.Join(sortedKeys(candidates), ", ") + ". Move the applicable case into " + expected + " preserving its package, name and assertions, or split out dedicated coverage if the case also covers another entrypoint."
			}
		}
		coverageBlockers[message] = true
		blocked[message] = true
	}
	result.UnverifiedCoverage = sortedKeys(unverifiedCoverage)
	result.UnreferencedDeclarations = sortedKeys(unreferenced)
	result.AccOnlyBlockers = sortedKeys(accOnly)
	result.CoverageBlockers = sortedKeys(coverageBlockers)
	result.ResourceFiles, result.Blocked = sortedKeys(resources), sortedKeys(blocked)
	return result
}

func main() {
	base, head := flag.String("base", "", "comparison SHA"), flag.String("head", "", "candidate SHA")
	layoutOnly := flag.Bool("layout-only", false, "validate new service file organization without planning cloud tests")
	worktree := flag.Bool("worktree", false, "include local/untracked changes; only valid with -layout-only")
	flag.Parse()
	if *base == "" || (!*worktree && *head == "") || strings.HasPrefix(*base, "-") || strings.HasPrefix(*head, "-") || (*worktree && !*layoutOnly) {
		fmt.Fprintln(os.Stderr, "use -base <commit> -head <commit>, or -base <commit> -worktree -layout-only; see "+layoutGuide)
		os.Exit(1)
	}
	resolvedBase, resolveErr := command("rev-parse", "--verify", *base+"^{commit}")
	if resolveErr != nil {
		fmt.Fprintln(os.Stderr, "cannot establish ACC scope: comparison base is unavailable; fetch the CI base commit and retry; see "+layoutGuide)
		os.Exit(1)
	}
	*base = strings.TrimSpace(string(resolvedBase))
	diffArgs := []string{"diff", "--no-renames", "--name-only", *base}
	if !*worktree {
		diffArgs = append(diffArgs, *head)
	}
	diffArgs = append(diffArgs, "--")
	files, err := command(diffArgs...)
	if err == nil && *worktree {
		var untracked []byte
		untracked, err = command("ls-files", "--others", "--exclude-standard", "--", "alicloud/")
		files = append(append(files, '\n'), untracked...)
	}
	failedLayout := false
	emit := func(before, after *snapshot, relevant []string) error {
		issues := layoutIssues(before, after, relevant)
		if *layoutOnly {
			for _, issue := range issues {
				fmt.Fprintln(os.Stderr, issue.String())
			}
			failedLayout = len(issues) != 0
			if failedLayout {
				fmt.Fprintf(os.Stderr, "Local recheck (no cloud calls): GO111MODULE=off GOWORK=off go run scripts/ci/accscope/main.go scripts/ci/accscope/layout.go -base %s -worktree -layout-only\n", *base)
			}
			return json.NewEncoder(os.Stdout).Encode(struct {
				Issues []layoutIssue `json:"layoutIssues"`
			}{issues})
		}
		scope := selectScope(before, after, relevant)
		scope.LayoutIssues = issues
		for _, issue := range issues {
			scope.Blocked = append(scope.Blocked, issue.String())
		}
		return json.NewEncoder(os.Stdout).Encode(scope)
	}
	if err == nil {
		var relevant []string
		for _, file := range strings.Split(strings.TrimSpace(string(files)), "\n") {
			if strings.HasPrefix(file, "alicloud/") || file == "go.mod" || file == "go.sum" {
				relevant = append(relevant, file)
			}
		}
		if len(relevant) == 0 {
			empty := &snapshot{decls: map[string]*declaration{}, files: map[string]bool{}, assets: map[string]string{}}
			err = emit(empty, empty, nil)
		} else {
			var before, after *snapshot
			before, err = readSnapshot(*base)
			if err == nil {
				if *worktree {
					after, err = readWorkingSnapshot()
				} else {
					after, err = readSnapshot(*head)
				}
			}
			if err == nil {
				err = emit(before, after, relevant)
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot establish ACC scope:", err)
		os.Exit(1)
	}
	if failedLayout {
		os.Exit(1)
	}
}
