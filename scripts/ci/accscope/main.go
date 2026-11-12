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
	method              bool
}

type declaration struct {
	id, pkg, name, file, receiver, hash string
	refs                                []reference
	assets                              []string
	test                                bool
	acc, unit                           bool
}

type snapshot struct {
	decls   map[string]*declaration
	files   map[string]bool
	assets  map[string]string
	hazards map[string]string
}

type accScope struct {
	Packages      map[string][]string `json:"packages"`
	ResourceFiles []string            `json:"resourceFiles"`
	Reasons       map[string][]string `json:"reasons"`
	Blocked       []string            `json:"blocked"`
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
		if fn, ok := x.Fun.(*ast.Ident); ok && fn.Name == "new" && len(x.Args) == 1 {
			return typeReference(x.Args[0], pkg, imports)
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

func parseSource(s *snapshot, file string, data []byte) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, data, 0)
	if err != nil {
		return err
	}
	pkg := path.Dir(file) + "#" + f.Name.Name
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
				d.refs = append(d.refs, reference{pkg: target.pkg, name: x.Sel.Name, receiver: target.name, method: true})
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

func resourceFile(file string) bool {
	if strings.HasSuffix(file, "_test.go") {
		return false
	}
	for _, prefix := range []string{"resource_", "data_source_", "ephemeral_", "action_", "list_"} {
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
	for id, d := range s.decls {
		key := d.pkg + ":" + d.name
		if d.receiver != "" {
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
	graph := map[string][]edge{}
	for id, d := range s.decls {
		for _, ref := range d.refs {
			key := ref.pkg + ":" + ref.name
			if ref.method {
				key = ref.pkg + ":*." + ref.name
			}
			matches := index[key]
			unresolved := false
			if ref.method {
				if ref.receiver != "" {
					matches = index[ref.pkg+":"+ref.receiver+"."+ref.name]
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
			for _, target := range matches {
				ambiguity := ""
				if len(matches) > 1 {
					ambiguity = "ambiguous receiver/build variants for " + key + " in " + d.file
				} else if unresolved {
					ambiguity = "unresolved cross-package receiver for " + key + " in " + d.file
				}
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
			blocked["No AccTest cases for affected resource: "+file+"; confirm replacement/coverage"] = true
		}
	}
	graphs := []map[string][]edge{reverseGraph(before), reverseGraph(after)}
	impact := func(seed, reason string) (bool, bool) {
		seen := map[string]bool{}
		queue := []string{seed}
		covered := false
		unit := false
		for len(queue) > 0 {
			id := queue[0]
			queue = queue[1:]
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
				if resourceFile(d.file) {
					selectResource(d.file, reason)
					covered = true
					// Provider factories register all resources; the implementation boundary is the type scope.
					continue
				}
			}
			for _, graph := range graphs {
				for _, e := range graph[id] {
					if e.ambiguity != "" {
						blocked[e.ambiguity+"; manual impact review required"] = true
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
		if d.test || strings.HasPrefix(d.file, "alicloud/acctest/") {
			if d.acc && current == nil {
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
			} else {
				acc, unit := impact(id, reason)
				if !acc && !unit && !strings.HasPrefix(d.name, "Test") && !strings.HasPrefix(d.name, "Benchmark") && !strings.HasPrefix(d.name, "Fuzz") {
					blocked["cannot resolve test helper consumers for "+reason+"; manual impact review required"] = true
				}
			}
		} else if resourceFile(d.file) {
			selectResource(d.file, reason)
		} else if service(d) {
			if covered, _ := impact(id, reason); !covered {
				blocked["cannot resolve service impact for "+reason+"; manual ACC scope review required"] = true
			}
		} else {
			blocked["manual ACC scope decision required for public code: "+reason] = true
		}
	}
	for file := range changedFiles {
		if !strings.HasPrefix(file, "alicloud/") || strings.HasSuffix(file, ".go") {
			continue
		}
		if before.assets[file] != after.assets[file] {
			acc, unit := impact("asset:"+file, file)
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
	result.ResourceFiles, result.Blocked = sortedKeys(resources), sortedKeys(blocked)
	return result
}

func main() {
	base, head := flag.String("base", "", "comparison SHA"), flag.String("head", "", "candidate SHA")
	flag.Parse()
	if *base == "" || *head == "" || strings.HasPrefix(*base, "-") || strings.HasPrefix(*head, "-") {
		fmt.Fprintln(os.Stderr, "both -base and -head are required")
		os.Exit(1)
	}
	files, err := command("diff", "--no-renames", "--name-only", *base, *head)
	if err == nil {
		var relevant []string
		for _, file := range strings.Split(strings.TrimSpace(string(files)), "\n") {
			if strings.HasPrefix(file, "alicloud/") || file == "go.mod" || file == "go.sum" {
				relevant = append(relevant, file)
			}
		}
		if len(relevant) == 0 {
			empty := &snapshot{decls: map[string]*declaration{}, files: map[string]bool{}, assets: map[string]string{}}
			err = json.NewEncoder(os.Stdout).Encode(selectScope(empty, empty, nil))
		} else {
			var before, after *snapshot
			before, err = readSnapshot(*base)
			if err == nil {
				after, err = readSnapshot(*head)
			}
			if err == nil {
				err = json.NewEncoder(os.Stdout).Encode(selectScope(before, after, relevant))
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot establish ACC scope:", err)
		os.Exit(1)
	}
}
