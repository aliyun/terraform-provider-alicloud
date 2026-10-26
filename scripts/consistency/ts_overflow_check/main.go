// Package main: timestamp_overflow_check is a pure go/ast static checker that prevents
// regressions of the 32-bit (386/arm) integer-overflow panic in timestamp attributes.
//
// Background: schema.TypeInt uses Go's native int (platform width). On 386/arm the SDK
// helper/schema read path calls strconv.ParseInt(value, 0, 0) then int(v), panicking on
// 13-digit millisecond / 19-digit nanosecond timestamps. Fix: TypeInt -> TypeString for
// timestamp fields. This checker guards the three regression vectors:
//
//	R1 (hard-fail, no baseline): a d.Set("field", expr) whose value derives from
//	time.Time.Unix{Milli,Micro,Nano}() while the field's schema is TypeInt. This is the
//	exact panic pattern. R1 detects the Unix-time selector via go/ast (SelectorExpr with
//	Sel.Name in {UnixMilli,UnixMicro,UnixNano}), NOT by regex over an ast.Fprint dump
//	(which never matches the structured output). After the type fix it must be zero
//	repo-wide.
//
//	R2 (hard-fail + documented exemptions): a TypeInt attribute whose documentation
//	declares a timestamp keyword AND a sub-second unit (ms/millisecond/ns/nanosecond/μs/
//	microsecond). Exemptions carry a reason and are keyed by "<doc-resource>::<path>",
//	where <path> is the dotted nested path (e.g. day_periods.start), so a bare field
//	name shared across sibling blocks (start/end in time_periods vs day_periods) is
//	disambiguated. The doc file is website/docs/{r,d}/<name>.html.markdown; the
//	description is the attribute's own list-item line, not a fixed-width window that
//	bleeds into neighbouring attributes.
//
//	R3 (snapshot baseline): a TypeInt attribute whose name matches the timestamp-shaped
//	regex (?:^|_)(time|timestamp|date|at) excluding the duration blacklist
//	(timeout|interval|duration|period|retention|ttl|wait|delay|idle|keep_alive|elapsed|
//	cost|window|span). The post-fix snapshot is stored in timestamp_overflow_baseline.txt;
//	the checker fails on NEW (current \ baseline) hits and reports cleanups
//	(baseline \ current) as a suggestion to refresh the baseline. R3 uses bare field
//	names (the schema field name) so the baseline is stable across nesting changes.
//
// The checker does NOT import any alicloud package; it parses source files with go/parser
// and go/ast only (equivalent to gofmt — pure syntax analysis, no Provider compilation).
// It eats a newline-separated list of changed .go file paths via -fileNames (the output of
// `git diff --name-only`), checking only changed resources. When -fileNames is empty it
// ships immediately. Pass -all to scan the whole alicloud/ tree; -all compares against the
// baseline without writing. Pass -update-baseline (implies -all) to rewrite the baseline;
// only this flag writes the baseline file, so a plain -all can never silently mask new hits.
//
//nolint:all
package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// config holds the resolved flags so the check logic is unit-testable without touching
// the process-global flag set or calling os.Exit.
type config struct {
	fileNames string
	all       bool
	baseline  string
	updateBL  bool
	verbose   bool
}

// R3 name shape: any underscore-or-start delimited segment equal to time/timestamp/date/at.
var r3NameSeg = regexp.MustCompile(`^(time|timestamp|date|at)$`)
var r3Black = regexp.MustCompile(`timeout|interval|duration|period|retention|ttl|wait|delay|idle|keep_?alive|elapsed|cost|window|span`)

// R2 doc keywords: an explicit "timestamp" mention (the strongest doc signal that a field
// holds an absolute epoch value, the only shape that overflows 386) plus a sub-second unit.
// Bare "time"/"date" are intentionally NOT matched: they appear in duration descriptions
// ("time-out time", "time interval", "monitoring timeout time") which are small bounded
// values that never overflow, so matching them would flag the whole repo as false positives.
// Fields named time/timestamp/date/at are still guarded by R3 regardless of doc wording.
var r2TSKeyword = regexp.MustCompile(`(?i)timestamp`)
var r2SubSecond = regexp.MustCompile(`(?i)millisecond|nanosecond|microsecond|\bms\b|\bns\b|\bμs\b|\bus\b`)

// unixTimeMethods is the set of time.Time selectors that feed a platform-width int from
// an absolute epoch and trigger the 386/arm panic when written into a TypeInt field.
var unixTimeMethods = map[string]bool{"UnixMilli": true, "UnixMicro": true, "UnixNano": true}

// r2Exemptions lists TypeInt attributes whose docs declare a sub-second timestamp but are
// intentionally kept TypeInt with a recorded reason. Key: "<doc-resource>::<path>" where
// <path> is the dotted nested schema path (e.g. day_periods.start).
var r2Exemptions = map[string]string{
	"alicloud_wafv3_defense_rule::day_periods.start": "bounded [0, 86400000) ms-of-day; fits int32, not overflow-prone; API uses ms-of-day not absolute time",
	"alicloud_wafv3_defense_rule::day_periods.end":   "bounded [0, 86400000) ms-of-day; fits int32, not overflow-prone; API uses ms-of-day not absolute time",
}

type fieldInfo struct {
	name string // bare schema field name (matches d.Set key and the R3 baseline)
	path string // dotted nested path (parent.child); used for R2 exemption keys
	typ  string // Int / String / Bool / ...
}

type fileResult struct {
	file       string
	fields     []fieldInfo // bare-name collection; used for R1 typeIntFields and R3
	pathFields []fieldInfo // path-aware collection; used for R2
	r1Hits     []r1Hit
	r3Hits     []string // field names matching R3, TypeInt
}

type r1Hit struct {
	file  string
	field string
	expr  string
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// parseFlags resolves args into a config; it returns exit code 2 on usage error.
func parseFlags(args []string, stderr io.Writer) (config, int) {
	cfg := config{
		baseline: "scripts/consistency/timestamp_overflow_baseline.txt",
	}
	fs := flag.NewFlagSet("ts_overflow_check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&cfg.fileNames, "fileNames", "", "path to a file containing newline-separated .go file paths to check (output of `git diff --name-only`)")
	fs.BoolVar(&cfg.all, "all", false, "scan the whole alicloud/ tree and compare against the baseline (does not write)")
	fs.StringVar(&cfg.baseline, "baseline", cfg.baseline, "path to the R3 baseline snapshot file")
	fs.BoolVar(&cfg.updateBL, "update-baseline", false, "scan the whole alicloud/ tree and rewrite the R3 baseline file from the current scan (use after a deliberate cleanup; commit the result)")
	fs.BoolVar(&cfg.verbose, "v", false, "verbose: print every hit")
	if err := fs.Parse(args); err != nil {
		return cfg, 2
	}
	return cfg, 0
}

// run is the testable entry point: it parses args, runs the check and returns the exit code.
func run(args []string, stdout, stderr io.Writer) int {
	cfg, code := parseFlags(args, stderr)
	if code != 0 {
		return code
	}
	return check(cfg, stdout, stderr)
}

func check(cfg config, stdout, stderr io.Writer) int {
	files := resolveFiles(cfg)
	if len(files) == 0 {
		fmt.Fprintln(stdout, "timestamp_overflow_check: no changed .go files; shipped.")
		return 0
	}

	var res []fileResult
	r3Current := map[string]bool{} // "file::field"
	for _, f := range files {
		fr := analyze(f)
		res = append(res, fr)
		for _, n := range fr.r3Hits {
			r3Current[f+"::"+n] = true
		}
	}

	exit := 0

	// R1 + R2: hard fail.
	var r1, r2 []string
	for _, fr := range res {
		for _, h := range fr.r1Hits {
			r1 = append(r1, fmt.Sprintf("R1 %s: d.Set(%q, %s) — TypeInt field fed from time.Unix{Milli,Micro,Nano}(); convert to TypeString + fmt.Sprint", h.file, h.field, h.expr))
		}
		r2 = append(r2, checkR2(fr, cfg)...)
	}
	for _, m := range r1 {
		fmt.Fprintln(stdout, m)
	}
	for _, m := range r2 {
		fmt.Fprintln(stdout, m)
	}
	if len(r1) > 0 || len(r2) > 0 {
		exit = 1
	}

	// R3 baseline rewrite is gated on -update-baseline only; a plain -all never writes,
	// so it cannot silently mask new hits.
	if cfg.updateBL {
		writeBaseline(cfg.baseline, r3Current)
		fmt.Fprintf(stdout, "R3 baseline written: %d entries across %d files -> %s\n", len(r3Current), len(files), cfg.baseline)
		return exit
	}

	changedSet := map[string]bool{}
	for _, f := range files {
		changedSet[f] = true
	}
	bl := loadBaseline(cfg.baseline)
	var news, cleanup []string
	for k := range r3Current {
		if !bl[k] {
			news = append(news, k)
		}
	}
	for k := range bl {
		if !r3Current[k] {
			// In -fileNames mode the scan only covers changed files; baseline entries
			// from unscanned files are not real cleanups, so restrict the report to
			// changed files. In -all mode changedSet covers all files and every cleanup
			// reports.
			if len(changedSet) > 0 {
				f := k
				if i := strings.Index(f, "::"); i >= 0 {
					f = f[:i]
				}
				if !changedSet[f] {
					continue
				}
			}
			cleanup = append(cleanup, k)
		}
	}
	sort.Strings(news)
	sort.Strings(cleanup)
	if len(news) > 0 {
		fmt.Fprintf(stdout, "R3 NEW TypeInt timestamp-named attributes (%d) — convert to TypeString or, if intentionally TypeInt, record a reason:\n", len(news))
		for _, k := range news {
			fmt.Fprintf(stdout, "  + %s\n", k)
		}
		exit = 1
	}
	if len(cleanup) > 0 {
		fmt.Fprintf(stdout, "R3 cleanups (%d) — refresh the baseline with `timestamp_overflow_check -update-baseline` and commit:\n", len(cleanup))
		for _, k := range cleanup {
			fmt.Fprintf(stdout, "  - %s\n", k)
		}
	}
	if cfg.verbose {
		for _, fr := range res {
			for _, h := range fr.r3Hits {
				fmt.Fprintf(stdout, "R3 hit %s::%s\n", fr.file, h)
			}
		}
	}

	if exit == 0 {
		fmt.Fprintf(stdout, "timestamp_overflow_check: PASS (R1=%d R2=%d R3_new=%d R3_cleanup=%d)\n", len(r1), len(r2), len(news), len(cleanup))
	} else {
		fmt.Fprintf(stdout, "timestamp_overflow_check: FAIL (R1=%d R2=%d R3_new=%d)\n", len(r1), len(r2), len(news))
	}
	return exit
}

func resolveFiles(cfg config) []string {
	if cfg.all || cfg.updateBL {
		matches, _ := filepath.Glob("alicloud/*.go")
		return matches
	}
	if cfg.fileNames == "" {
		return nil
	}
	b, err := os.ReadFile(cfg.fileNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read -fileNames: %v\n", err)
		return nil
	}
	raw := strings.Split(string(b), "\n")
	// The CI orchestrator (scripts/ci/run-provider-check.py) and the local parity script
	// feed the FULL `git diff` as -fileNames, mirroring consistency_check.go. A plain
	// newline-separated path list (`git diff --name-only`) is also accepted. Detect which
	// form we got so a full diff is never mistaken for an empty path list (which would
	// silently ship and mask every regression).
	isDiff := false
	for _, ln := range raw {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "diff --git ") || strings.HasPrefix(t, "+++ ") || strings.HasPrefix(t, "--- ") {
			isDiff = true
			break
		}
	}
	var out []string
	seen := map[string]bool{}
	for _, ln := range raw {
		t := strings.TrimSpace(ln)
		var p string
		if isDiff {
			// Only the new-path line (+++ b/<path>) uniquely identifies a changed file;
			// --- a/ lines would otherwise duplicate renames and deletions.
			switch {
			case strings.HasPrefix(t, "+++ b/"):
				p = strings.TrimPrefix(t, "+++ b/")
			case strings.HasPrefix(t, "+++ "):
				p = strings.TrimPrefix(t, "+++ ")
			default:
				continue
			}
		} else {
			p = t
		}
		if p == "" || !strings.HasPrefix(p, "alicloud/") || !strings.HasSuffix(p, ".go") {
			continue
		}
		if !strings.HasPrefix(p, "alicloud/resource_") && !strings.HasPrefix(p, "alicloud/data_source_") {
			continue
		}
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}

// analyze parses one file and collects schema field types, R1 hits, and R3 hits.
func analyze(file string) fileResult {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", file, err)
		return fileResult{file: file}
	}
	src, _ := os.ReadFile(file)
	return analyzeAST(file, fset, f, src)
}

// analyzeAST is the file-agnostic core of analyze, factored so tests can feed inline source.
func analyzeAST(file string, fset *token.FileSet, f *ast.File, src []byte) fileResult {
	fr := fileResult{file: file}
	// Bare-name collection: walk every composite literal keyed by a string literal; for
	// each, find the Type field's selector (schema.TypeX) at the SAME nesting level. This
	// is kept stable (bare names) so the R3 baseline is byte-stable across nesting changes.
	ast.Inspect(f, func(n ast.Node) bool {
		cl, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		for _, el := range cl.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.BasicLit)
			if !ok || key.Kind != token.STRING {
				continue
			}
			name := unquote(key.Value)
			val := kv.Value
			// schema fields are written as &schema.Schema{...}; unwrap the leading &.
			if unary, ok := val.(*ast.UnaryExpr); ok {
				val = unary.X
			}
			valCL, ok := val.(*ast.CompositeLit)
			if !ok {
				continue
			}
			t := typeOfSchema(valCL)
			if t != "" {
				fr.fields = append(fr.fields, fieldInfo{name: name, typ: t})
			}
		}
		return true
	})
	// Path-aware collection for R2: dotted nested paths (parent.child) so the exemption
	// key disambiguates a bare field name shared across sibling blocks (start/end).
	fr.pathFields = collectPathFields(f)
	// R3: TypeInt + name shape (bare name).
	for _, fi := range fr.fields {
		if fi.typ == "Int" && r3Keep(fi.name) {
			fr.r3Hits = append(fr.r3Hits, fi.name)
		}
	}
	// R1: d.Set("field", expr) where expr calls time.Unix{Milli,Micro,Nano}() and field is TypeInt.
	typeIntFields := map[string]bool{}
	for _, fi := range fr.fields {
		if fi.typ == "Int" {
			typeIntFields[fi.name] = true
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok || se.Sel == nil || se.Sel.Name != "Set" {
			return true
		}
		if len(ce.Args) < 2 {
			return true
		}
		keyLit, ok := ce.Args[0].(*ast.BasicLit)
		if !ok || keyLit.Kind != token.STRING {
			return true
		}
		field := unquote(keyLit.Value)
		if !typeIntFields[field] {
			return true
		}
		if method := unixTimeCallIn(ce.Args[1]); method != "" {
			fr.r1Hits = append(fr.r1Hits, r1Hit{file: file, field: field, expr: nodeSource(fset, ce.Args[1], src)})
		}
		return true
	})
	return fr
}

// collectPathFields walks schema map literals and returns field infos with a dotted path of
// the form "<immediate-parent>.<field>" (two segments; top-level fields carry only their own
// name). This is the path shape used in R2 exemption keys (e.g. day_periods.start), so a bare
// field name shared across sibling blocks — start/end in time_periods vs day_periods — is
// disambiguated by its immediate parent without tracking the full root-to-leaf path. Each
// schema map is processed once (the visitor returns false for the subtree so nested maps
// reached via Elem recursion are not re-visited), preventing duplicate bare entries.
func collectPathFields(f *ast.File) []fieldInfo {
	var out []fieldInfo
	var walkMap func(mapCL *ast.CompositeLit, prefix string)
	walkMap = func(mapCL *ast.CompositeLit, prefix string) {
		for _, el := range mapCL.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.BasicLit)
			if !ok || key.Kind != token.STRING {
				continue
			}
			name := unquote(key.Value)
			full := name
			if prefix != "" {
				full = prefix + "." + name
			}
			val := kv.Value
			if unary, ok := val.(*ast.UnaryExpr); ok {
				val = unary.X
			}
			schCL, ok := val.(*ast.CompositeLit)
			if !ok {
				continue
			}
			t := typeOfSchema(schCL)
			if t != "" {
				out = append(out, fieldInfo{name: name, path: full, typ: t})
			}
			// Recurse into Elem: &schema.Resource{ Schema: <map> }.
			for _, e := range schCL.Elts {
				ekv, ok := e.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				eident, ok := ekv.Key.(*ast.Ident)
				if !ok || eident.Name != "Elem" {
					continue
				}
				ev := ekv.Value
				if unary, ok := ev.(*ast.UnaryExpr); ok {
					ev = unary.X
				}
				resCL, ok := ev.(*ast.CompositeLit)
				if !ok {
					continue
				}
				for _, re := range resCL.Elts {
					rkv, ok := re.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					rident, ok := rkv.Key.(*ast.Ident)
					if !ok || rident.Name != "Schema" {
						continue
					}
					rv := rkv.Value
					if unary, ok := rv.(*ast.UnaryExpr); ok {
						rv = unary.X
					}
					rmap, ok := rv.(*ast.CompositeLit)
					if !ok {
						continue
					}
					walkMap(rmap, name)
				}
			}
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		mapCL, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		if _, ok := mapCL.Type.(*ast.MapType); !ok {
			return true
		}
		if !isSchemaMap(mapCL) {
			return true
		}
		walkMap(mapCL, "")
		return false // do not let ast.Inspect re-visit nested schema maps
	})
	return out
}

// isSchemaMap reports whether a map composite literal has at least one schema.Schema
// element (a field whose value composite carries a schema.TypeX selector). This avoids
// treating unrelated map[string]X literals as schema maps.
func isSchemaMap(mapCL *ast.CompositeLit) bool {
	for _, el := range mapCL.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		val := kv.Value
		if unary, ok := val.(*ast.UnaryExpr); ok {
			val = unary.X
		}
		schCL, ok := val.(*ast.CompositeLit)
		if !ok {
			continue
		}
		if typeOfSchema(schCL) != "" {
			return true
		}
	}
	return false
}

// typeOfSchema inspects a schema field composite literal for its Type selector at the
// immediate level (schema.TypeInt etc.), ignoring nested Elem resources.
func typeOfSchema(cl *ast.CompositeLit) string {
	for _, el := range cl.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kv.Key.(*ast.Ident)
		if !ok || ident.Name != "Type" {
			continue
		}
		sel, ok := kv.Value.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil {
			return ""
		}
		// sel.Sel.Name is the full identifier (TypeInt, TypeString, ...); strip the
		// "Type" prefix so callers compare against "Int"/"String"/"Bool"/"Set"/"List".
		name := sel.Sel.Name
		if strings.HasPrefix(name, "Type") {
			name = name[len("Type"):]
		}
		return name
	}
	return ""
}

// unixTimeCallIn reports the name of the first time.Unix{Milli,Micro,Nano}() selector
// found anywhere in expr (via go/ast SelectorExpr inspection), or "" if none. This replaces
// the previous ast.Fprint + regex approach, which never matched the structured dump.
func unixTimeCallIn(expr ast.Node) string {
	found := ""
	ast.Inspect(expr, func(n ast.Node) bool {
		se, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if se.Sel != nil && unixTimeMethods[se.Sel.Name] {
			found = se.Sel.Name
			return false
		}
		return true
	})
	return found
}

func r3Keep(name string) bool {
	for _, seg := range strings.Split(name, "_") {
		if r3NameSeg.MatchString(seg) {
			if r3Black.MatchString(name) {
				return false
			}
			return true
		}
	}
	return false
}

// checkR2: for each TypeInt field (by dotted path) of this file whose doc declares a
// timestamp keyword AND a sub-second unit, report a violation unless exempted.
func checkR2(fr fileResult, cfg config) []string {
	docPath, docRes := docFor(fr.file)
	if docPath == "" {
		return nil
	}
	// _test.go files have no website doc and no user-facing schema contract; R2 does not
	// apply to them. Skip before constructing a doc path so the path never points at a
	// non-existent *_test.html.markdown.
	if strings.HasSuffix(fr.file, "_test.go") {
		return nil
	}
	body, err := os.ReadFile(docPath)
	if err != nil {
		// Doc absent for a non-test resource/data_source file: R2 cannot be evaluated
		// here. With the corrected .html.markdown path this only affects files that
		// genuinely lack a doc; R1/R3 still apply. Skip rather than emitting per-field
		// noise, but surface the gap once when verbose.
		if cfg.verbose {
			fmt.Printf("R2 %s: doc not found at %s; skipping R2 for this file\n", fr.file, docPath)
		}
		return nil
	}
	return checkR2Doc(fr, docRes, string(body), cfg.verbose)
}

// checkR2Doc is the disk-free core of checkR2, factored so tests can feed inline doc text.
func checkR2Doc(fr fileResult, docRes, text string, verbose bool) []string {
	var out []string
	for _, fi := range fr.pathFields {
		if fi.typ != "Int" {
			continue
		}
		child := fi.path
		if i := strings.LastIndex(child, "."); i >= 0 {
			child = child[i+1:]
		}
		if !fieldDocMentionsTimestamp(text, child) {
			continue
		}
		key := docRes + "::" + fi.path
		if reason, ok := r2Exemptions[key]; ok {
			if verbose {
				fmt.Printf("R2 exempt %s: %s\n", key, reason)
			}
			continue // documented exemption
		}
		out = append(out, fmt.Sprintf("R2 %s: TypeInt attribute %q (path %q) documents a sub-second timestamp; convert to TypeString or add an exemption with a reason (exemption key: %s)", fr.file, fi.name, fi.path, key))
	}
	return out
}

// fieldDocMentionsTimestamp reports whether the doc contains an attribute list-item line
// for field whose description mentions a timestamp keyword AND a sub-second unit. The
// description is the field's own list-item line plus continuation lines up to the next
// list item or blank line, so a neighbouring attribute's description (e.g. image_create
// next to image_size) does not bleed into the window.
func fieldDocMentionsTimestamp(doc, field string) bool {
	lines := strings.Split(doc, "\n")
	for i := 0; i < len(lines); i++ {
		desc, ok := listItemFieldDesc(lines[i], field)
		if !ok {
			continue
		}
		window := desc
		for j := i + 1; j < len(lines); j++ {
			nxt := strings.TrimSpace(lines[j])
			if nxt == "" || isListItem(nxt) {
				break
			}
			window += " " + nxt
		}
		if r2TSKeyword.MatchString(window) && r2SubSecond.MatchString(window) {
			return true
		}
	}
	return false
}

// listItemFieldDesc extracts the description text that follows `field` in a markdown
// list-item line of the form `* `field` - desc` / `- `field` - desc` (also accepting an
// em dash). ok is false if the line is not a list item naming field.
func listItemFieldDesc(ln, field string) (desc string, ok bool) {
	t := strings.TrimSpace(ln)
	if t == "" || (t[0] != '*' && t[0] != '-') {
		return "", false
	}
	rest := strings.TrimSpace(t[1:])
	if !strings.HasPrefix(rest, "`") {
		return "", false
	}
	end := strings.Index(rest[1:], "`")
	if end < 0 {
		return "", false
	}
	name := rest[1 : 1+end]
	if name != field {
		return "", false
	}
	after := strings.TrimSpace(rest[2+end:])
	after = strings.TrimPrefix(after, "-")
	after = strings.TrimPrefix(after, "—")
	return strings.TrimSpace(after), true
}

func isListItem(s string) bool {
	return len(s) > 0 && (s[0] == '*' || s[0] == '-')
}

// docFor maps a changed .go file to its website doc path and resource name.
// resource_alicloud_foo.go -> website/docs/r/alicloud_foo.html.markdown, resource "alicloud_foo"
// data_source_alicloud_foo.go -> website/docs/d/alicloud_foo.html.markdown
// docFor maps a changed .go file to its website doc path and resource name.
// resource_alicloud_foo.go -> website/docs/r/foo.html.markdown, resource "alicloud_foo"
// data_source_alicloud_foo.go -> website/docs/d/foo.html.markdown
//
// res keeps the alicloud_ prefix (it is the Terraform resource type and the exemption-key
// identifier); the website doc FILE name omits the alicloud_ prefix, so the path strips it.
// The previous form pointed at alicloud_foo.html.markdown, which does not exist — R2 then
// silently skipped every file, masking real violations.
func docFor(file string) (path, res string) {
	base := filepath.Base(file)
	if strings.HasPrefix(base, "resource_") {
		res = strings.TrimPrefix(base, "resource_")
		res = strings.TrimSuffix(res, ".go")
		docName := strings.TrimPrefix(res, "alicloud_")
		return filepath.Join("website", "docs", "r", docName+".html.markdown"), res
	}
	if strings.HasPrefix(base, "data_source_") {
		res = strings.TrimPrefix(base, "data_source_")
		res = strings.TrimSuffix(res, ".go")
		docName := strings.TrimPrefix(res, "alicloud_")
		return filepath.Join("website", "docs", "d", docName+".html.markdown"), res
	}
	return "", ""
}

func loadBaseline(path string) map[string]bool {
	out := map[string]bool{}
	b, err := os.ReadFile(path)
	if err != nil {
		return out // absent baseline: no guard yet; R3 only reports new as everything.
	}
	s := bufio.NewScanner(strings.NewReader(string(b)))
	for s.Scan() {
		ln := strings.TrimSpace(s.Text())
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		out[ln] = true
	}
	return out
}

func writeBaseline(path string, current map[string]bool) {
	keys := make([]string, 0, len(current))
	for k := range current {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	w, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create baseline: %v\n", err)
		return
	}
	defer w.Close()
	fmt.Fprintln(w, "# R3 baseline: TypeInt attributes whose name matches the timestamp shape")
	fmt.Fprintln(w, "# (segment time|timestamp|date|at, excluding the duration blacklist).")
	fmt.Fprintln(w, "# Generated by `go run scripts/consistency/ts_overflow_check/main.go -update-baseline`.")
	fmt.Fprintln(w, "# Format: <file>::<field>. Do not hand-edit; regenerate after a deliberate cleanup.")
	for _, k := range keys {
		fmt.Fprintln(w, k)
	}
}

func unquote(s string) string {
	if len(s) >= 2 && (s[0] == '"' || s[0] == '`') {
		s = s[1 : len(s)-1]
	}
	return s
}

// nodeSource returns the raw source text of n by slicing the file bytes using fset
// position offsets. It returns "" when src is unavailable (e.g. tests with inline AST),
// so R1 diagnostics degrade gracefully instead of relying on ast.Fprint.
func nodeSource(fset *token.FileSet, n ast.Node, src []byte) string {
	if n == nil || len(src) == 0 {
		return ""
	}
	sp := fset.Position(n.Pos()).Offset
	ep := fset.Position(n.End()).Offset
	if sp < 0 || ep < sp || ep > len(src) {
		return ""
	}
	return string(src[sp:ep])
}
