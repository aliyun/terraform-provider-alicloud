// Package coverage checks for executable acceptance-test steps, without running cloud operations.
package coverage

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build/constraint"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Check requires an isolated TypeList reorder plan, followed by a convergent apply.
func Check(filename string, src []byte, resourceName string, r *schema.Resource, aliases ...string) []string {
	candidates := Candidates(r)
	f, err := parser.ParseFile(token.NewFileSet(), filename, src, parser.ParseComments)
	if err != nil {
		return []string{err.Error()}
	}
	var issues []string
	for _, group := range f.Comments {
		if group.Pos() > f.Package {
			break
		}
		for _, c := range group.List {
			if constraint.IsGoBuild(c.Text) || constraint.IsPlusBuild(c.Text) {
				issues = append(issues, "build-constrained test files cannot establish acceptance test discovery")
			}
		}
	}
	resourceNames := map[string]bool{resourceName: true}
	for _, name := range aliases {
		resourceNames[name] = true
	}
	covered := map[string]bool{}
	sdk := importName(f, "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource")
	testing := importName(f, "testing")
	if sdk != "" && testing != "" {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || !strings.HasPrefix(fn.Name.Name, "TestAcc") || fn.Body == nil || fn.Recv != nil || fn.Type.Params.NumFields() != 1 {
				continue
			}
			param := fn.Type.Params.List[0]
			pointer, ok := param.Type.(*ast.StarExpr)
			if !ok || len(param.Names) != 1 || !selector(pointer.X, testing, "T") {
				continue
			}
			inspectTest(fn, sdk, param.Names[0].Name, resourceNames, candidates, covered)
		}
	}
	paths := make([]string, 0, len(candidates))
	for path := range candidates {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if !covered[path] {
			issues = append(issues, fmt.Sprintf("%s: TypeList requires adjacent apply A -> reordered PlanOnly with ExpectNonEmptyPlan:true -> apply B with an empty post-apply plan (at least two distinct members); see scripts/collection-order/README.md", path))
		}
	}
	return issues
}

// Candidates includes multi-member configurable TypeList fields, including nested lists.
func Candidates(r *schema.Resource) map[string]bool {
	out := map[string]bool{}
	collect(r.Schema, "", out)
	return out
}

func collect(fields map[string]*schema.Schema, prefix string, out map[string]bool) {
	for name, s := range fields {
		if s == nil || (!s.Optional && !s.Required) {
			continue
		}
		path := prefix + name
		if s.Type == schema.TypeList || s.Type == schema.TypeSet {
			if s.Type == schema.TypeList && s.MaxItems != 1 {
				out[path] = true
			}
			if nested, ok := s.Elem.(*schema.Resource); ok {
				collect(nested.Schema, path+".", out)
			}
		}
	}
}

func importName(f *ast.File, path string) string {
	for _, imp := range f.Imports {
		if text(imp.Path) == path {
			if imp.Name != nil {
				if imp.Name.Name == "." || imp.Name.Name == "_" {
					return ""
				}
				return imp.Name.Name
			}
			parts := strings.Split(path, "/")
			return parts[len(parts)-1]
		}
	}
	return ""
}
func selector(expr ast.Expr, pkg, name string) bool {
	s, ok := expr.(*ast.SelectorExpr)
	return ok && ident(s.X) == pkg && s.Sel.Name == name
}
func ident(expr ast.Expr) string {
	if i, ok := expr.(*ast.Ident); ok {
		return i.Name
	}
	return ""
}
func text(expr ast.Expr) string {
	if s, ok := expr.(*ast.BasicLit); ok && s.Kind == token.STRING {
		v, _ := strconv.Unquote(s.Value)
		return v
	}
	return ""
}
func fields(expr ast.Expr) map[string]ast.Expr {
	lit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	result := map[string]ast.Expr{}
	for _, e := range lit.Elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			return nil
		}
		key := ident(kv.Key)
		if key == "" {
			key = text(kv.Key)
		}
		if key == "" {
			return nil
		}
		if _, exists := result[key]; exists {
			return nil
		}
		result[key] = kv.Value
	}
	return result
}
func falseOrAbsent(expr ast.Expr) bool { return expr == nil || ident(expr) == "false" }

type builder struct {
	address string
	config  map[string]interface{}
}
type snapshot struct {
	builder string
	config  map[string]interface{}
	paths   []string
}

func inspectTest(fn *ast.FuncDecl, sdk, tvar string, resourceNames map[string]bool, candidates map[string]bool, covered map[string]bool) {
	if skipsTest(fn.Body, tvar) {
		return
	}
	stringsByName := map[string]string{}
	builders := map[string]*builder{}
	for _, stmt := range fn.Body.List {
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
				continue
			}
			name := ident(s.Lhs[0])
			if name == "" {
				continue
			}
			delete(stringsByName, name)
			delete(builders, name)
			if value := text(s.Rhs[0]); value != "" {
				stringsByName[name] = value
			}
			call, ok := s.Rhs[0].(*ast.CallExpr)
			if ok && ident(call.Fun) == "resourceTestAccConfigFunc" && len(call.Args) == 3 {
				address := text(call.Args[0])
				if address == "" {
					address = stringsByName[ident(call.Args[0])]
				}
				parts := strings.Split(address, ".")
				if len(parts) == 2 && resourceNames[parts[0]] && parts[1] != "" {
					builders[name] = &builder{address: address, config: map[string]interface{}{}}
				}
			} else {
				invalidateBuilderUses(s.Rhs[0], builders)
			}
		case *ast.DeclStmt:
			// Other local declarations do not establish supported config builders.
			invalidateBuilderUses(s, builders)
		case *ast.ExprStmt:
			call, ok := s.X.(*ast.CallExpr)
			if !ok {
				return
			}
			if selector(call.Fun, sdk, "Test") || selector(call.Fun, sdk, "ParallelTest") {
				if len(call.Args) == 2 && ident(call.Args[0]) == tvar {
					tc, ok := call.Args[1].(*ast.CompositeLit)
					if ok && selector(tc.Type, sdk, "TestCase") {
						attrs := fields(tc)
						for name, expr := range attrs {
							if name != "Steps" {
								invalidateBuilderUses(expr, builders)
							}
						}
						if attrs["ErrorCheck"] == nil && falseOrAbsent(attrs["IsUnitTest"]) {
							inspectSteps(attrs["Steps"], sdk, builders, candidates, covered)
						}
					}
				}
				// Do not infer state across separate test runners.
				return
			}
			invalidateBuilderUses(s, builders)
		default:
			return
		}
	}
}

func skipsTest(expr ast.Node, tvar string) bool {
	skipped := false
	if expr == nil {
		return false
	}
	ast.Inspect(expr, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			for _, method := range []string{"Skip", "Skipf", "SkipNow", "Fatal", "Fatalf", "FailNow"} {
				if selector(call.Fun, tvar, method) {
					skipped = true
				}
			}
		}
		return true
	})
	return skipped
}

func invalidateBuilderUses(node ast.Node, builders map[string]*builder) {
	ast.Inspect(node, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			delete(builders, id.Name)
		}
		return true
	})
}

func inspectSteps(expr ast.Expr, sdk string, builders map[string]*builder, candidates map[string]bool, covered map[string]bool) {
	steps, ok := expr.(*ast.CompositeLit)
	if !ok {
		return
	}
	typ, ok := steps.Type.(*ast.ArrayType)
	if !ok || !selector(typ.Elt, sdk, "TestStep") {
		return
	}
	var previous, pending *snapshot
	for _, step := range steps.Elts {
		applied, planned := previous, pending
		previous, pending = nil, nil
		attrs := fields(step)
		for name, expr := range attrs {
			if name != "Config" {
				invalidateBuilderUses(expr, builders)
			}
		}
		config, ok := attrs["Config"].(*ast.CallExpr)
		if !ok || len(config.Args) != 1 {
			continue
		}
		key := ident(config.Fun)
		b := builders[key]
		if b == nil {
			continue
		}
		invalidateBuilderUses(config.Args[0], builders)
		if builders[key] == nil {
			continue
		}
		delta, ok := value(config.Args[0]).(map[string]interface{})
		if !ok {
			delete(builders, key)
			continue
		}
		merged := map[string]interface{}{}
		for k, v := range b.config {
			merged[k] = v
		}
		for k, v := range delta {
			switch v {
			case sentinel("REMOVEKEY"), "#REMOVEKEY":
				delete(merged, k)
			case sentinel("CLEARLIST"), "#CLEARLIST":
				merged[k] = []interface{}{}
			case sentinel("CLEARMAP"), "#CLEARMAP":
				merged[k] = map[string]interface{}{}
			default:
				merged[k] = v
			}
		}
		b.config = merged
		clean := !ignoresChanges(merged["lifecycle"]) && attrs["ExpectError"] == nil && falseOrAbsent(attrs["Destroy"]) && falseOrAbsent(attrs["ImportState"]) && attrs["Taint"] == nil && attrs["SkipFunc"] == nil && attrs["PreConfig"] == nil && attrs["PreventPostDestroyRefresh"] == nil && falseOrAbsent(attrs["RefreshState"])
		if clean && falseOrAbsent(attrs["PlanOnly"]) && falseOrAbsent(attrs["ExpectNonEmptyPlan"]) {
			if planned != nil && planned.builder == key && reflect.DeepEqual(planned.config, merged) {
				for _, path := range planned.paths {
					covered[path] = true
				}
			}
			previous = &snapshot{builder: key, config: merged}
		}
		if clean && ident(attrs["PlanOnly"]) == "true" && ident(attrs["ExpectNonEmptyPlan"]) == "true" && applied != nil && applied.builder == key {
			var paths []string
			for path := range candidates {
				equal, reordered := permutationAt(applied.config, merged, strings.Split(path, "."))
				if equal && reordered {
					paths = append(paths, path)
				}
			}
			pending = &snapshot{builder: key, config: merged, paths: paths}
		}
	}
}

// Ignore rules can hide exactly the drift that the permutation step must detect.
func ignoresChanges(lifecycle interface{}) bool {
	switch v := lifecycle.(type) {
	case nil:
		return false
	case map[string]interface{}:
		_, ignored := v["ignore_changes"]
		return ignored
	case []interface{}:
		for _, block := range v {
			if ignoresChanges(block) {
				return true
			}
		}
		return false
	default:
		return true // A dynamic lifecycle value cannot prove that changes are observed.
	}
}

type sentinel string
type opaque string

// Literal containers make membership reviewable. Scalar expressions are compared syntactically;
// actual dependency values and remote behavior are verified when the acceptance test runs.
func value(expr ast.Expr) interface{} {
	switch e := expr.(type) {
	case *ast.CompositeLit:
		switch e.Type.(type) {
		case *ast.MapType:
			out := map[string]interface{}{}
			for _, el := range e.Elts {
				kv, ok := el.(*ast.KeyValueExpr)
				if !ok {
					return nil
				}
				k := text(kv.Key)
				if k == "" {
					return nil
				}
				if _, exists := out[k]; exists {
					return nil
				}
				v := value(kv.Value)
				if v == nil {
					return nil
				}
				out[k] = v
			}
			return out
		case *ast.ArrayType:
			out := []interface{}{}
			for _, el := range e.Elts {
				v := value(el)
				if v == nil {
					return nil
				}
				out = append(out, v)
			}
			return out
		case nil:
			// Go permits elided map types inside a slice literal.
			out := map[string]interface{}{}
			for _, el := range e.Elts {
				kv, ok := el.(*ast.KeyValueExpr)
				if !ok {
					return nil
				}
				k := text(kv.Key)
				v := value(kv.Value)
				if k == "" || v == nil {
					return nil
				}
				out[k] = v
			}
			return out
		}
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return text(e)
		}
	case *ast.Ident:
		if e.Name == "REMOVEKEY" || e.Name == "CLEARLIST" || e.Name == "CLEARMAP" {
			return sentinel(e.Name)
		}
	}
	if expr == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, token.NewFileSet(), expr); err != nil {
		return nil
	}
	return opaque(buf.String())
}

// permutationAt also checks that every value outside the nominated collection is unchanged.
func permutationAt(a, b interface{}, path []string) (bool, bool) {
	if len(path) == 0 {
		aa, aok := a.([]interface{})
		bb, bok := b.([]interface{})
		if !aok || !bok || len(aa) < 2 || len(aa) != len(bb) {
			return reflect.DeepEqual(a, b), false
		}
		if reflect.DeepEqual(aa, bb) {
			return true, false
		}
		counts := map[string]int{}
		for _, v := range aa {
			counts[fmt.Sprintf("%#v", v)]++
		}
		if len(counts) < 2 {
			return false, false
		}
		for _, v := range bb {
			counts[fmt.Sprintf("%#v", v)]--
		}
		for _, n := range counts {
			if n != 0 {
				return false, false
			}
		}
		return true, true
	}
	if aa, ok := a.(map[string]interface{}); ok {
		bb, ok := b.(map[string]interface{})
		if !ok || len(aa) != len(bb) {
			return false, false
		}
		reordered := false
		for k, av := range aa {
			bv, exists := bb[k]
			if !exists {
				return false, false
			}
			if k == path[0] {
				equal, changed := permutationAt(av, bv, path[1:])
				if !equal {
					return false, false
				}
				reordered = reordered || changed
			} else if !reflect.DeepEqual(av, bv) {
				return false, false
			}
		}
		return true, reordered
	}
	if aa, ok := a.([]interface{}); ok {
		bb, ok := b.([]interface{})
		if !ok || len(aa) != len(bb) {
			return false, false
		}
		reordered := false
		for i := range aa {
			equal, changed := permutationAt(aa[i], bb[i], path)
			if !equal {
				return false, false
			}
			reordered = reordered || changed
		}
		return true, reordered
	}
	return reflect.DeepEqual(a, b), false
}
