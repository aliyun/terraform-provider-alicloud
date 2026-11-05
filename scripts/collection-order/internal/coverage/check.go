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
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Check requires convergent reorder coverage for every candidate TypeList: an
// apply of the original member order followed by a convergent apply of a member
// permutation. Attributes carrying a DiffSuppressFunc may evidence the reorder
// with a clean PlanOnly step instead, and their suppressor must be named and
// bounded by a unit net. Only convergence is asserted: whether a reorder diffs
// depends on the attribute's semantics and the API's ordering contract, which
// review establishes.
func Check(filename string, src []byte, testSources map[string][]byte, resourceName string, r *schema.Resource, aliases ...string) []string {
	candidates := Candidates(r)
	suppressed := Suppressed(r)
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
	bypassed := markerPaths(f, bypassMarker)
	nets := map[string]suppressNet{}
	suppressorNames := map[string]bool{}
	for _, name := range suppressed {
		if name != "" {
			suppressorNames[name] = true
		}
	}
	for name := range suppressorNames {
		nets[name] = resolveSuppressNet(f, testSources, name)
	}
	resourceNames := map[string]bool{resourceName: true}
	for _, name := range aliases {
		resourceNames[name] = true
	}
	covered := map[string]bool{}
	suppressedCovered := map[string]bool{}
	// The harness import flipped to terraform-plugin-testing; accept either
	// so pre-flip-shaped sources keep being inspected.
	sdk := importName(f, "github.com/hashicorp/terraform-plugin-testing/helper/resource")
	if sdk == "" {
		sdk = importName(f, "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource")
	}
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
			inspectTest(fn, sdk, param.Names[0].Name, resourceNames, candidates, suppressed, covered, suppressedCovered)
		}
	}
	paths := make([]string, 0, len(candidates))
	for path := range candidates {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if bypassed[path] {
			continue
		}
		name, isSuppressed := suppressed[path]
		if !isSuppressed {
			if !covered[path] {
				issues = append(issues, fmt.Sprintf("%s: TypeList requires an apply of the original order followed by a convergent apply with the members reordered (at least two distinct members); see scripts/collection-order/README.md", path))
			}
			continue
		}
		if name == "" {
			issues = append(issues, fmt.Sprintf("%s: anonymous DiffSuppressFunc cannot establish the suppression unit net; use a named top-level function or remove the suppressor; see scripts/collection-order/README.md", path))
			continue
		}
		var missing []string
		if !suppressedCovered[path] {
			missing = append(missing, "a reordered convergent apply, or a reordered PlanOnly step with a clean plan, after an apply of the original order")
		}
		if !nets[name].reorder {
			missing = append(missing, "a suppressor unit case where only member order changes and suppression is expected")
		}
		if !nets[name].membership {
			missing = append(missing, "a suppressor unit case where membership changes and suppression is not expected")
		}
		if len(missing) > 0 {
			issues = append(issues, fmt.Sprintf("%s: suppressed TypeList requires %s; see scripts/collection-order/README.md", path, strings.Join(missing, "; ")))
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
}

func inspectTest(fn *ast.FuncDecl, sdk, tvar string, resourceNames map[string]bool, candidates map[string]bool, suppressed map[string]string, covered map[string]bool, suppressedCovered map[string]bool) {
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
			if ok && isConfigBuilder(call.Fun) && len(call.Args) == 3 {
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
							inspectSteps(attrs["Steps"], sdk, builders, candidates, suppressed, covered, suppressedCovered)
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

func inspectSteps(expr ast.Expr, sdk string, builders map[string]*builder, candidates map[string]bool, suppressed map[string]string, covered map[string]bool, suppressedCovered map[string]bool) {
	steps, ok := expr.(*ast.CompositeLit)
	if !ok {
		return
	}
	typ, ok := steps.Type.(*ast.ArrayType)
	if !ok || !selector(typ.Elt, sdk, "TestStep") {
		return
	}
	var previous *snapshot
	for _, step := range steps.Elts {
		applied := previous
		attrs := fields(step)
		for name, expr := range attrs {
			if name != "Config" {
				invalidateBuilderUses(expr, builders)
			}
		}
		config, ok := attrs["Config"].(*ast.CallExpr)
		if !ok || len(config.Args) != 1 {
			previous = nil
			continue
		}
		key := ident(config.Fun)
		b := builders[key]
		if b == nil {
			previous = nil
			continue
		}
		invalidateBuilderUses(config.Args[0], builders)
		if builders[key] == nil {
			previous = nil
			continue
		}
		delta, ok := value(config.Args[0]).(map[string]interface{})
		if !ok {
			delete(builders, key)
			previous = nil
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
		if !clean {
			previous = nil
			continue
		}
		if ident(attrs["PlanOnly"]) == "true" {
			// Plan-only steps observe the plan without applying, so they leave the
			// last applied config in place. A clean plan over a reordered config is
			// reorder evidence for suppressed attributes; ExpectNonEmptyPlan is a
			// diff assertion, not convergence, and establishes nothing here.
			if falseOrAbsent(attrs["ExpectNonEmptyPlan"]) && applied != nil && applied.builder == key {
				for path := range suppressed {
					equal, reordered := permutationAt(applied.config, merged, strings.Split(path, "."))
					if equal && reordered {
						suppressedCovered[path] = true
					}
				}
			}
			continue
		}
		if !falseOrAbsent(attrs["ExpectNonEmptyPlan"]) {
			// An apply expecting a non-empty plan asserts non-convergence.
			previous = nil
			continue
		}
		if applied != nil && applied.builder == key {
			// A convergent apply of a member permutation of the previously applied
			// config is the reorder coverage: the empty post-apply plan the SDK
			// requires proves the resource converged on the reordered members.
			for path := range candidates {
				equal, reordered := permutationAt(applied.config, merged, strings.Split(path, "."))
				if equal && reordered {
					covered[path] = true
					suppressedCovered[path] = true
				}
			}
		}
		previous = &snapshot{builder: key, config: merged}
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

// Suppressed maps candidate paths that declare a DiffSuppressFunc to the suppressor's
// function name. Anonymous functions and bound methods map to the empty string: they
// cannot carry a unit net.
func Suppressed(r *schema.Resource) map[string]string {
	out := map[string]string{}
	suppressedAt(r.Schema, "", out)
	return out
}

func suppressedAt(fields map[string]*schema.Schema, prefix string, out map[string]string) {
	for name, s := range fields {
		if s == nil || (!s.Optional && !s.Required) {
			continue
		}
		path := prefix + name
		if s.Type == schema.TypeList || s.Type == schema.TypeSet {
			if s.Type == schema.TypeList && s.MaxItems != 1 && s.DiffSuppressFunc != nil {
				out[path] = suppressorName(s.DiffSuppressFunc)
			}
			if nested, ok := s.Elem.(*schema.Resource); ok {
				suppressedAt(nested.Schema, path+".", out)
			}
		}
	}
}

func suppressorName(fn schema.SchemaDiffSuppressFunc) string {
	value := runtime.FuncForPC(reflect.ValueOf(fn).Pointer())
	if value == nil {
		return ""
	}
	full := value.Name()
	if strings.Contains(full, ".func") {
		return "" // closures carry no stable name for a unit net
	}
	name := full
	if i := strings.LastIndexByte(name, '.'); i >= 0 {
		name = name[i+1:]
	}
	if strings.Contains(name, "-") || !token.IsIdentifier(name) {
		return ""
	}
	return name
}

var bypassMarker = regexp.MustCompile(`^//\s*collection-order:\s*bypass\s+([A-Za-z0-9_.]+)\s+--\s*(\S.*)$`)

// markerPaths collects the attribute paths carrying a structured marker comment. The
// reason is required and stays reviewable in the diff; free-form comments do not match.
func markerPaths(f *ast.File, marker *regexp.Regexp) map[string]bool {
	out := map[string]bool{}
	for _, group := range f.Comments {
		for _, c := range group.List {
			if m := marker.FindStringSubmatch(c.Text); m != nil {
				out[m[1]] = true
			}
		}
	}
	return out
}

// suppressNet records the bounded-suppression unit cases established for one
// suppressor function.
type suppressNet struct {
	reorder    bool
	membership bool
}

// resolveSuppressNet reads the house unit-test form: a test function ranging over a
// literal []struct table, calling the suppressor with the row's serialized old and new
// values and comparing the result with the row's expected value. The unit net may live
// in a sibling *_unit_test.go file in the same package rather than the main *_test.go
// file that carries the TestAcc steps; testSources carries those sibling sources so the
// net is still discovered when the split-package layout separates them.
func resolveSuppressNet(f *ast.File, testSources map[string][]byte, name string) suppressNet {
	var net suppressNet
	files := []*ast.File{f}
	for _, src := range testSources {
		tf, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
		if err != nil {
			continue
		}
		files = append(files, tf)
	}
	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			for _, row := range unitRows(fn, name) {
				oldMembers, newMembers := splitMembers(row.oldValue), splitMembers(row.newValue)
				switch {
				case sameMembers(oldMembers, newMembers) && row.oldValue != row.newValue:
					if row.expected {
						net.reorder = true
					}
				case !sameMembers(oldMembers, newMembers):
					if !row.expected {
						net.membership = true
					}
				}
			}
		}
	}
	return net
}

type unitRow struct {
	oldValue, newValue string
	expected           bool
}

// unitRows resolves table rows whose old, new and expected fields are literals.
func unitRows(fn *ast.FuncDecl, name string) []unitRow {
	var rows []unitRow
	tables := map[string]*ast.CompositeLit{}
	for _, stmt := range fn.Body.List {
		s, ok := stmt.(*ast.AssignStmt)
		if !ok || len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			continue
		}
		if id := ident(s.Lhs[0]); id != "" {
			if lit, ok := s.Rhs[0].(*ast.CompositeLit); ok {
				tables[id] = lit
			} else {
				delete(tables, id)
			}
		}
	}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		r, ok := n.(*ast.RangeStmt)
		if !ok {
			return true
		}
		table, ok := r.X.(*ast.CompositeLit)
		if !ok {
			table, ok = tables[ident(r.X)]
		}
		v := ident(r.Value)
		if !ok || v == "" {
			return true
		}
		oldField, newField, expectedField := unitCallFields(r.Body, v, name)
		if oldField == "" || newField == "" || expectedField == "" {
			return true
		}
		fieldOrder := structFieldNames(table)
		for _, el := range table.Elts {
			row, ok := el.(*ast.CompositeLit)
			if !ok {
				continue
			}
			values := rowValues(row, fieldOrder)
			if values == nil {
				continue
			}
			oldExpr, hasOld := values[oldField]
			newExpr, hasNew := values[newField]
			expectedExpr, hasExpected := values[expectedField]
			if !hasOld || !hasNew || !hasExpected || !isStringLiteral(oldExpr) || !isStringLiteral(newExpr) {
				continue
			}
			want, ok := boolIdent(expectedExpr)
			if !ok {
				continue
			}
			rows = append(rows, unitRow{oldValue: text(oldExpr), newValue: text(newExpr), expected: want})
		}
		return true
	})
	return rows
}

func structFieldNames(table *ast.CompositeLit) []string {
	arr, ok := table.Type.(*ast.ArrayType)
	if !ok {
		return nil
	}
	st, ok := arr.Elt.(*ast.StructType)
	if !ok {
		return nil
	}
	names := make([]string, 0, len(st.Fields.List))
	for _, field := range st.Fields.List {
		if len(field.Names) != 1 {
			return nil
		}
		names = append(names, field.Names[0].Name)
	}
	return names
}

// rowValues accepts keyed and positional row literals.
func rowValues(row *ast.CompositeLit, fieldOrder []string) map[string]ast.Expr {
	if len(row.Elts) == 0 {
		return nil
	}
	values := map[string]ast.Expr{}
	if _, ok := row.Elts[0].(*ast.KeyValueExpr); ok {
		for _, el := range row.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
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
			if _, exists := values[key]; exists {
				return nil
			}
			values[key] = kv.Value
		}
		return values
	}
	if len(row.Elts) != len(fieldOrder) {
		return nil
	}
	for i, el := range row.Elts {
		values[fieldOrder[i]] = el
	}
	return values
}

// unitCallFields finds the table fields used as the suppressor's old and new arguments
// and the field its result is compared with.
func unitCallFields(body ast.Node, v, name string) (string, string, string) {
	var oldField, newField, expectedField string
	results := map[string]bool{}
	ast.Inspect(body, func(n ast.Node) bool {
		switch e := n.(type) {
		case *ast.CallExpr:
			if len(e.Args) >= 3 && callName(e.Fun) == name {
				if field := fieldOf(e.Args[1], v); field != "" {
					oldField = field
				}
				if field := fieldOf(e.Args[2], v); field != "" {
					newField = field
				}
			}
		case *ast.AssignStmt:
			if len(e.Lhs) != 1 || len(e.Rhs) != 1 {
				return true
			}
			id := ident(e.Lhs[0])
			if id == "" {
				return true
			}
			if isSuppressorCall(e.Rhs[0], name) {
				results[id] = true
			} else {
				delete(results, id)
			}
		case *ast.BinaryExpr:
			if e.Op != token.EQL && e.Op != token.NEQ {
				return true
			}
			var field string
			switch {
			case isSuppressorCall(e.X, name) || results[ident(e.X)]:
				field = fieldOf(e.Y, v)
			case isSuppressorCall(e.Y, name) || results[ident(e.Y)]:
				field = fieldOf(e.X, v)
			}
			if field != "" {
				expectedField = field
			}
		}
		return true
	})
	return oldField, newField, expectedField
}

func callName(fun ast.Expr) string {
	if s, ok := fun.(*ast.SelectorExpr); ok {
		return s.Sel.Name
	}
	return ident(fun)
}

// isConfigBuilder recognizes both the bare legacy constructor
// (resourceTestAccConfigFunc) and the split-package qualified constructor
// (acctest.ResourceTestAccConfig). The selector's package alias is not
// constrained because the import alias is a per-file choice; the
// constructor name is unique within the provider.
func isConfigBuilder(fun ast.Expr) bool {
	if id, ok := fun.(*ast.Ident); ok {
		return id.Name == "resourceTestAccConfigFunc"
	}
	if s, ok := fun.(*ast.SelectorExpr); ok {
		return s.Sel.Name == "ResourceTestAccConfig"
	}
	return false
}

func isSuppressorCall(expr ast.Expr, name string) bool {
	call, ok := expr.(*ast.CallExpr)
	return ok && callName(call.Fun) == name
}

func fieldOf(expr ast.Expr, v string) string {
	if s, ok := expr.(*ast.SelectorExpr); ok && ident(s.X) == v {
		return s.Sel.Name
	}
	return ""
}

func boolIdent(expr ast.Expr) (value, ok bool) {
	if id, isIdent := expr.(*ast.Ident); isIdent && (id.Name == "true" || id.Name == "false") {
		return id.Name == "true", true
	}
	return false, false
}

func isStringLiteral(expr ast.Expr) bool {
	s, ok := expr.(*ast.BasicLit)
	return ok && s.Kind == token.STRING
}

// splitMembers splits the serialized old/new diff strings the SDK passes to a
// DiffSuppressFunc.
func splitMembers(s string) []string {
	return strings.Split(s, ",")
}

func sameMembers(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	counts := map[string]int{}
	for _, v := range a {
		counts[v]++
	}
	for _, v := range b {
		counts[v]--
		if counts[v] < 0 {
			return false
		}
	}
	return true
}
