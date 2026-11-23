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
)

const module = "github.com/aliyun/terraform-provider-alicloud/"

type reference struct {
	pkg, name, receiver string
	method, factory     bool
	line                int
	field               string
}

type declaration struct {
	id, pkg, name, file, receiver, hash string
	refs                                []reference
	result                              reference
	signature                           string
	interfaceMethods                    []string
	structFields                        map[string]reference
	assets                              []string
	test                                bool
	acc, unit                           bool
	assertedTypes                       []string
	privateTestFunc                     bool
}

type snapshot struct {
	decls              map[string]*declaration
	files              map[string]bool
	assets             map[string]string
	hazards            map[string]string
	registrationIssues map[string]string
	registrations      map[string]map[string]string
	packages           map[string]string
}

type accScope struct {
	Packages           map[string][]string `json:"packages"`
	ResourceFiles      []string            `json:"resourceFiles"`
	Reasons            map[string][]string `json:"reasons"`
	Blocked            []string            `json:"blocked"`
	UnverifiedCoverage []string            `json:"unverifiedCoverage,omitempty"`
	CoverageBlockers   []string            `json:"coverageBlockers"`
	LayoutIssues       []layoutIssue       `json:"layoutIssues,omitempty"`
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
		return typeReference(x.X, pkg, imports)
	case *ast.SelectorExpr:
		if alias, ok := x.X.(*ast.Ident); ok {
			if imported := imports[alias.Name]; imported != "" {
				return reference{pkg: strings.TrimPrefix(imported, module), name: x.Sel.Name}
			}
		}
	}
	return reference{}
}

func valueType(expr ast.Expr, pkg string, imports map[string]string, depth int) reference {
	if depth > 8 {
		return reference{}
	}
	switch x := expr.(type) {
	case *ast.CompositeLit:
		return typeReference(x.Type, pkg, imports)
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
					if name.Name == x.Name && i < len(decl.Values) {
						return valueType(decl.Values[i], pkg, imports, depth+1)
					}
				}
			case *ast.AssignStmt:
				for i, name := range decl.Lhs {
					if id, ok := name.(*ast.Ident); ok && id.Name == x.Name && i < len(decl.Rhs) {
						return valueType(decl.Rhs[i], pkg, imports, depth+1)
					}
				}
			}
		}
	}
	return reference{}
}

// Read the static registration slices and leave other ServicePackage logic in
// the declaration fingerprint so mixed edits cannot pass as mapping-only changes.
func normalizeRegistration(fn *ast.FuncDecl) (map[string]string, string) {
	if fn.Body == nil || len(fn.Body.List) != 1 {
		return nil, "return expression"
	}
	ret, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return nil, "return expression"
	}
	body, ok := ret.Results[0].(*ast.CompositeLit)
	if !ok {
		return nil, "return expression"
	}
	mappings := map[string]string{}
	var remaining []ast.Expr
	for _, entry := range body.Elts {
		field, ok := entry.(*ast.KeyValueExpr)
		if !ok {
			return nil, "return expression"
		}
		kind := receiverName(field.Key)
		switch kind {
		case "Resources", "SDKResources", "DataSources", "SDKDataSources", "EphemeralResources", "Actions", "ListResources":
		default:
			remaining = append(remaining, entry)
			continue
		}
		list, ok := field.Value.(*ast.CompositeLit)
		if !ok {
			return nil, kind
		}
		slice, ok := list.Type.(*ast.ArrayType)
		if !ok || slice.Len != nil {
			return nil, kind
		}
		elem, ok := slice.Elt.(*ast.SelectorExpr)
		owner, qualified := body.Type.(*ast.SelectorExpr)
		if !ok || !qualified || owner.Sel.Name != "ServicePackage" || receiverName(elem.X) != receiverName(owner.X) || elem.Sel.Name != strings.TrimSuffix(kind, "s") {
			return nil, kind
		}
		for _, element := range list.Elts {
			item, ok := element.(*ast.CompositeLit)
			if !ok || item.Type != nil || len(item.Elts) != 2 {
				return nil, kind
			}
			var name, factory string
			for _, attr := range item.Elts {
				pair, ok := attr.(*ast.KeyValueExpr)
				if !ok {
					return nil, kind
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
				return nil, kind
			}
			mappings[key] = factory
		}
	}
	body.Elts = remaining
	return mappings, ""
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
	add := func(name, receiver string, node ast.Node) {
		key := name
		if receiver != "" {
			key = receiver + "." + name
		}
		if name == "init" {
			key += "@" + file
		}
		d := &declaration{id: pkg + ":" + key + "@" + file, pkg: pkg, name: name, receiver: receiver,
			file: file, test: strings.HasSuffix(file, "_test.go")}
		switch decl := node.(type) {
		case *ast.FuncDecl:
			d.privateTestFunc = d.test && decl.Recv == nil && decl.Body != nil &&
				!ast.IsExported(name) && name != "init" && name != "_"
			if receiver != "" {
				d.signature = methodSignature(fset, decl.Type)
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
					for _, name := range field.Names {
						d.structFields[name.Name] = typeReference(field.Type, pkg, imports)
					}
				}
			}
		case *ast.ValueSpec:
			for i, ident := range decl.Names {
				if ident.Name == name && i < len(decl.Values) {
					if fn, ok := decl.Values[i].(*ast.FuncLit); ok && fn.Type.Results != nil && len(fn.Type.Results.List) == 1 &&
						len(fn.Type.Results.List[0].Names) <= 1 {
						d.result = typeReference(fn.Type.Results.List[0].Type, pkg, imports)
					}
				}
			}
		}
		if fn, ok := node.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Type.Results != nil {
			results := fn.Type.Results.List
			if len(results) == 1 && len(results[0].Names) <= 1 {
				d.result = typeReference(results[0].Type, pkg, imports)
			}
		}
		if fn, ok := node.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Type.Results == nil && len(fn.Type.Params.List) == 1 {
			if pointer, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr); ok {
				if typ, ok := pointer.X.(*ast.SelectorExpr); ok && typ.Sel.Name == "T" {
					d.acc = strings.HasPrefix(name, "TestAcc")
					d.unit = strings.HasPrefix(name, "Test") && !d.acc
				}
			}
		}
		usedImports := map[string]string{}
		ast.Inspect(node, func(n ast.Node) bool {
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
				if target.name == "" {
					target.pkg = pkg
				}
				d.refs = append(d.refs, reference{pkg: target.pkg, name: x.Sel.Name,
					receiver: target.name, method: true, factory: target.factory, field: target.field, line: fset.Position(x.Sel.Pos()).Line})
				ast.Inspect(x.X, func(child ast.Node) bool {
					if id, ok := child.(*ast.Ident); ok && (id.Obj == nil || top[id.Obj.Decl.(ast.Node)]) {
						d.refs = append(d.refs, reference{pkg: pkg, name: id.Name})
					}
					return true
				})
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
					switch fn.Sel.Name {
					case "NewServer", "NewTLSServer", "NewUnstartedServer":
						if alias, ok := fn.X.(*ast.Ident); ok && alias.Obj == nil && imports[alias.Name] == "net/http/httptest" {
							handler := valueType(x.Args[0], pkg, imports, 0)
							if handler.name != "" {
								// httptest accepts http.Handler and invokes ServeHTTP through the interface.
								d.refs = append(d.refs, reference{pkg: handler.pkg, name: "ServeHTTP",
									receiver: handler.name, method: true, factory: handler.factory,
									field: handler.field, line: fset.Position(fn.Pos()).Line})
							}
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
		})
		// Import changes can rebind an unchanged selector; include that binding in the fingerprint.
		bindings, _ := json.Marshal(usedImports)
		d.hash = fmt.Sprintf("%x", sha256.Sum256([]byte(nodeText(fset, node)+directives+string(bindings))))
		s.decls[d.id] = d
	}
	for _, decl := range f.Decls {
		switch x := decl.(type) {
		case *ast.FuncDecl:
			receiver := ""
			if x.Recv != nil {
				receiver = receiverName(x.Recv.List[0].Type)
			}
			if path.Base(file) == "register.go" && x.Name.Name == "ServicePackage" && receiver == "" {
				if mappings, field := normalizeRegistration(x); mappings != nil {
					if s.registrations == nil {
						s.registrations = map[string]map[string]string{}
					}
					s.registrations[file] = mappings
				} else {
					if s.registrationIssues == nil {
						s.registrationIssues = map[string]string{}
					}
					s.registrationIssues[file] = field
				}
			}
			add(x.Name.Name, receiver, x)
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
			}
			index[d.pkg+":"+d.receiver+"."+d.name] = append(index[d.pkg+":"+d.receiver+"."+d.name], id)
			if !d.test {
				qualified := strings.Split(d.pkg, "#")[0] + ":" + d.receiver + "." + d.name
				index[qualified] = append(index[qualified], id)
			}
		}
		index[key] = append(index[key], id)
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
			if ref.factory {
				factories := index[ref.pkg+":"+ref.receiver]
				result := reference{}
				if len(factories) == 1 {
					result = s.decls[factories[0]].result
				}
				ref.pkg, ref.receiver = result.pkg, result.name
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
			unresolved, interfaceCall := false, false
			if ref.method {
				if ref.receiver != "" {
					matches = index[ref.pkg+":"+ref.receiver+"."+ref.name]
					if implementations := interfaces[ref.pkg+":"+ref.receiver]; len(implementations) > 0 {
						interfaceCall = true
						for _, receiver := range implementations {
							matches = append(matches, index[receiver+"."+ref.name]...)
						}
					}
					if len(matches) == 0 && (strings.HasPrefix(ref.pkg, "alicloud/") || strings.HasPrefix(ref.pkg, "alicloud#")) {
						matches = index["method:"+ref.name]
						unresolved = true
					}
				} else {
					candidates := map[string]bool{}
					for _, target := range append(append([]string{}, index["method:"+ref.name]...), index[key]...) {
						candidates[target] = true
					}
					matches = sortedKeys(candidates)
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

func selectScope(before, after *snapshot, files []string) accScope {
	result := accScope{Packages: map[string][]string{}, ResourceFiles: []string{}, Reasons: map[string][]string{}, Blocked: []string{}}
	selected, resources, blocked := map[string]bool{}, map[string]bool{}, map[string]bool{}
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
	selectResource := func(file, reason string) {
		if moved := movedFiles[file]; moved != "" {
			file = moved
		}
		resources[file] = true
		stem := strings.TrimSuffix(file, ".go")
		found := false
		for _, d := range after.decls {
			if !d.test || !d.acc {
				continue
			}
			if d.file == stem+"_test.go" || (strings.HasPrefix(d.file, stem+"_") &&
				!after.files[strings.TrimSuffix(d.file, "_test.go")+".go"]) {
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
			if !d.test && d.receiver == "" && d.result.name != "" {
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
	graphs := []map[string][]edge{reverseGraph(before), reverseGraph(after)}
	hasLinkname := false
	for _, snap := range []*snapshot{before, after} {
		for _, directives := range snap.hazards {
			hasLinkname = hasLinkname || strings.Contains(directives, "//go:linkname")
		}
	}
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
	impact := func(seed, reason string, requireResolved bool) (bool, bool) {
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
			if d := after.decls[id]; d != nil {
				if d.test && d.acc {
					selectTest(d, reason)
					covered = true
				}
				if d.test && d.unit {
					unit = true
				}
				if resourceFile(d.file, registered) {
					selectResource(d.file, reason)
					covered = true
					// Provider factories register all resources; the implementation boundary is the type scope.
					continue
				}
			}
			for _, graph := range graphs {
				for _, e := range graph[id] {
					if e.ambiguity != "" {
						if requireResolved {
							blocked[e.ambiguity+"; manual impact review required"] = true
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
				blocked["removed/renamed AccTest "+reason+"; confirm replacement coverage"] = true
				continue
			}
			if !d.test && d.name == "init" {
				blocked["test-support package initialization changed: "+reason+"; manual importer impact review required"] = true
			} else if d.name == "TestMain" || d.name == "init" {
				for _, test := range after.decls {
					if path.Dir(test.file) == path.Dir(d.file) {
						selectTest(test, reason)
					}
				}
			} else if !unusedTestFunc(old, current) {
				acc, unit := impact(id, reason, true)
				if !acc && !unit && !strings.HasPrefix(d.name, "Test") && !strings.HasPrefix(d.name, "Benchmark") && !strings.HasPrefix(d.name, "Fuzz") {
					blocked["cannot resolve test helper consumers for "+reason+"; manual impact review required"] = true
				}
			}
		} else if resourceFile(d.file, registered) {
			directResources[d.file] = true
			selectResource(d.file, reason)
		} else if service(d) {
			if covered, unit := impact(id, reason, true); !covered && !unit {
				blocked["cannot resolve service impact for "+reason+"; manual ACC scope review required"] = true
			}
		} else if path.Base(d.file) == "register.go" {
			if field := after.registrationIssues[d.file]; field != "" {
				reason += " (" + field + ")"
			}
			blocked["unsupported registration change: "+reason+"; manual impact review required"] = true
		} else if strings.HasPrefix(dir, "alicloud/service/") {
			// Do not turn unknown helper consumers into a successful empty scope.
			if covered, unit := impact(id, reason, true); !covered && !unit {
				blocked["cannot resolve production helper consumers for "+reason+"; manual ACC scope review required"] = true
			}
		}
	}
	for file := range changedFiles {
		if !strings.HasPrefix(file, "alicloud/") || strings.HasSuffix(file, ".go") {
			continue
		}
		if before.assets[file] != after.assets[file] {
			acc, unit := impact("asset:"+file, file, true)
			if !acc && !unit && !strings.HasSuffix(file, ".md") && !strings.HasSuffix(file, ".markdown") {
				blocked["cannot resolve test fixture consumers for "+file+"; manual impact review required"] = true
			}
		}
	}
	byPackage := map[string]map[string]bool{}
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
