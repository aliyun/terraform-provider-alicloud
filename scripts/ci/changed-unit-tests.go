// CI_UNIT_TESTS_V1: run changed-file unit tests inside ci-check-quick.
// This tool uses the Go parser and build constraints, never filename arguments to go test.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type unitCase struct {
	Package string `json:"package"`
	Name    string `json:"name"`
	Result  string `json:"result"`
}

type unitSummary struct {
	Schema            int        `json:"schema"`
	Base              string     `json:"base"`
	Status            string     `json:"status"`
	Reason            string     `json:"reason,omitempty"`
	ChangedFiles      []string   `json:"changedFiles"`
	BroadenedPackages []string   `json:"broadenedPackages"`
	Cases             []unitCase `json:"cases"`
	Passed            int        `json:"passed"`
	Failed            int        `json:"failed"`
	Skipped           int        `json:"skipped"`
	Missing           int        `json:"missing"`
}

func changedFiles(base string) ([]string, error) {
	if _, err := exec.Command("git", "rev-parse", "--verify", base+"^{commit}").Output(); err != nil {
		return nil, fmt.Errorf("invalid CI_CHECK_BASE: %w", err)
	}
	b, err := exec.Command("git", "diff", "--name-only", "-z", "--no-renames", base, "--").Output()
	if err != nil {
		return nil, err
	}
	files := strings.Split(strings.TrimSuffix(string(b), "\x00"), "\x00")
	if len(files) == 1 && files[0] == "" {
		files = []string{}
	}
	return files, nil
}

func unitNames(file string) ([]string, error) {
	ok, err := build.Default.MatchFile(filepath.Dir(file), filepath.Base(file))
	if err != nil || !ok {
		return nil, err
	}
	f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "TestUnit") {
			names = append(names, fn.Name.Name)
		}
	}
	return names, nil
}

// registrationMaps finds the two literal registration maps in Provider. Anything
// outside this known shape uses the shared-code fallback instead of guessing.
func registrationMaps(f *ast.File) (map[string]*ast.CompositeLit, bool) {
	maps := map[string]*ast.CompositeLit{}
	providers := 0
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != "Provider" {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok || nodeText(lit.Type) != "schema.Provider" {
				return true
			}
			providers++
			for _, elt := range lit.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok || (key.Name != "ResourcesMap" && key.Name != "DataSourcesMap") {
					continue
				}
				value, ok := kv.Value.(*ast.CompositeLit)
				if !ok || nodeText(value.Type) != "map[string]*schema.Resource" || maps[key.Name] != nil {
					providers += 2 // invalidate dynamic maps and duplicate fields
					continue
				}
				maps[key.Name] = value
			}
			return true
		})
	}
	return maps, providers == 1 && len(maps) == 2
}

func nodeText(n ast.Node) string {
	var b bytes.Buffer
	if err := format.Node(&b, token.NewFileSet(), n); err != nil {
		return ""
	}
	return b.String()
}

func registrationEntries(lit *ast.CompositeLit) (map[string]ast.Expr, bool) {
	entries := map[string]ast.Expr{}
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			return nil, false
		}
		key, ok := kv.Key.(*ast.BasicLit)
		if !ok || key.Kind != token.STRING {
			return nil, false
		}
		name, err := strconv.Unquote(key.Value)
		if err != nil || entries[name] != nil {
			return nil, false
		}
		entries[name] = kv.Value
	}
	return entries, true
}

// registrationChanges proves that the entire file differs only by added/removed
// registrations. Replacing an existing key's factory is deliberately shared code.
func registrationChanges(before, after []byte) ([]ast.Expr, bool) {
	// Compiler/build directives affect execution even though ordinary comments do not.
	var directives [2][]string
	for i, source := range [][]byte{before, after} {
		for _, line := range strings.Split(string(source), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "//go:") || strings.HasPrefix(line, "// +build") {
				directives[i] = append(directives[i], line)
			}
		}
	}
	if strings.Join(directives[0], "\n") != strings.Join(directives[1], "\n") {
		return nil, false
	}
	oldFile, err := parser.ParseFile(token.NewFileSet(), "provider.go", before, 0)
	if err != nil {
		return nil, false
	}
	newFile, err := parser.ParseFile(token.NewFileSet(), "provider.go", after, 0)
	if err != nil {
		return nil, false
	}
	oldMaps, ok := registrationMaps(oldFile)
	if !ok {
		return nil, false
	}
	newMaps, ok := registrationMaps(newFile)
	if !ok {
		return nil, false
	}
	var changed []ast.Expr
	for _, name := range []string{"ResourcesMap", "DataSourcesMap"} {
		oldEntries, ok := registrationEntries(oldMaps[name])
		if !ok {
			return nil, false
		}
		newEntries, ok := registrationEntries(newMaps[name])
		if !ok {
			return nil, false
		}
		for key, oldValue := range oldEntries {
			if newValue, exists := newEntries[key]; exists {
				if nodeText(oldValue) != nodeText(newValue) {
					return nil, false
				}
			} else {
				changed = append(changed, oldValue)
			}
		}
		for key, value := range newEntries {
			if _, exists := oldEntries[key]; !exists {
				changed = append(changed, value)
			}
		}
		oldMaps[name].Elts = nil
		newMaps[name].Elts = nil
	}
	return changed, nodeText(oldFile) == nodeText(newFile)
}

func registrationTests(base, file string) ([]string, bool) {
	before, err := exec.Command("git", "show", base+":"+filepath.ToSlash(file)).Output()
	if err != nil {
		return nil, false
	}
	after, err := os.ReadFile(file)
	if err != nil {
		return nil, false
	}
	changed, ok := registrationChanges(before, after)
	if !ok {
		return nil, false
	}
	// Resolve actual factory declarations, so aliases select the implementation's
	// companion rather than guessing a filename from the registration key.
	factories := map[string]bool{}
	for _, expr := range changed {
		call, ok := expr.(*ast.CallExpr)
		if !ok || len(call.Args) != 0 || call.Ellipsis.IsValid() {
			return nil, false
		}
		fn, ok := call.Fun.(*ast.Ident)
		if !ok { // qualified/dynamic factories retain package-wide coverage
			return nil, false
		}
		factories[fn.Name] = true
	}
	if len(factories) == 0 {
		return nil, true
	}
	paths, err := filepath.Glob(filepath.Join(filepath.Dir(file), "*.go"))
	if err != nil {
		return nil, false
	}
	matches := map[string][]string{}
	for _, path := range paths {
		name := filepath.Base(path)
		if strings.HasSuffix(name, "_test.go") || (!strings.HasPrefix(name, "resource_") && !strings.HasPrefix(name, "data_source_")) {
			continue
		}
		active, err := build.Default.MatchFile(filepath.Dir(path), name)
		if err != nil || !active {
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return nil, false
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Recv == nil && factories[fn.Name.Name] {
				matches[fn.Name.Name] = append(matches[fn.Name.Name], strings.TrimSuffix(path, ".go")+"_test.go")
			}
		}
	}
	var tests []string
	for factory := range factories {
		if len(matches[factory]) != 1 {
			return nil, false
		}
		test := matches[factory][0]
		if _, err := os.Stat(test); err != nil {
			return nil, false
		}
		tests = append(tests, test)
	}
	sort.Strings(tests)
	return tests, true
}

// Resource/data-source implementations have companion tests. Every other source
// file is shared helper/service code and broadens to all TestUnit* in its package.
func selectCases(files []string, bases ...string) ([]unitCase, []string, error) {
	selected := map[string]map[string]bool{}
	broad := map[string]bool{}
	allPackages := false
	for _, file := range files {
		if file == "go.mod" || file == "go.sum" {
			allPackages = true
		}
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		dir := filepath.Dir(file)
		if selected[dir] == nil {
			selected[dir] = map[string]bool{}
		}
		if file == "alicloud/provider.go" && len(bases) == 1 {
			if tests, ok := registrationTests(bases[0], file); ok {
				for _, test := range tests {
					selected[dir][test] = true
				}
				continue
			}
		}
		name := filepath.Base(file)
		if strings.HasSuffix(name, "_test.go") {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				broad[dir] = true
			} else {
				selected[dir][file] = true
			}
		} else if strings.HasPrefix(name, "resource_") || strings.HasPrefix(name, "data_source_") {
			companion := strings.TrimSuffix(file, ".go") + "_test.go"
			if _, err := os.Stat(companion); err == nil {
				selected[dir][companion] = true
			} else if !os.IsNotExist(err) {
				return nil, nil, err
			}
		} else {
			broad[dir] = true
		}
	}
	if allPackages {
		err := filepath.WalkDir(".", func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && (d.Name() == ".git" || d.Name() == "vendor") {
				return filepath.SkipDir
			}
			if !d.IsDir() && strings.HasSuffix(p, "_test.go") {
				broad[filepath.Dir(p)] = true
			}
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}
	var broadened []string
	for dir := range broad {
		broadened = append(broadened, "./"+filepath.ToSlash(dir))
		matches, err := filepath.Glob(filepath.Join(dir, "*_test.go"))
		if err != nil {
			return nil, nil, err
		}
		if selected[dir] == nil {
			selected[dir] = map[string]bool{}
		}
		for _, file := range matches {
			selected[dir][file] = true
		}
	}
	var cases []unitCase
	seen := map[string]bool{}
	for dir, tests := range selected {
		for file := range tests {
			names, err := unitNames(file)
			if err != nil {
				return nil, nil, err
			}
			for _, name := range names {
				pkg := "./" + filepath.ToSlash(dir)
				key := pkg + "\x00" + name
				if !seen[key] {
					cases = append(cases, unitCase{pkg, name, "MISSING"})
					seen[key] = true
				}
			}
		}
	}
	sort.Slice(cases, func(i, j int) bool {
		if cases[i].Package != cases[j].Package {
			return cases[i].Package < cases[j].Package
		}
		return cases[i].Name < cases[j].Name
	})
	sort.Strings(broadened)
	return cases, broadened, nil
}

func runCases(cases []unitCase) error {
	for start := 0; start < len(cases); {
		end := start + 1
		for end < len(cases) && cases[end].Package == cases[start].Package {
			end++
		}
		names := []string{}
		byName := map[string]int{}
		for i := start; i < end; i++ {
			names = append(names, regexp.QuoteMeta(cases[i].Name))
			byName[cases[i].Name] = i
		}
		args := []string{"test", "-json", "-count=1", "-p=1", cases[start].Package, "-run=^(" + strings.Join(names, "|") + ")$", "-timeout=10m"}
		fmt.Printf("Changed unit command: go %s\n", strings.Join(args, " "))
		cmd := exec.Command("go", args...)
		cmd.Env = []string{}
		for _, e := range os.Environ() {
			if !strings.HasPrefix(e, "TF_ACC=") {
				cmd.Env = append(cmd.Env, e)
			}
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		cmd.Stderr = os.Stderr
		if err = cmd.Start(); err != nil {
			return err
		}
		decoder := json.NewDecoder(stdout)
		for {
			var event struct {
				Action string
				Test   string
				Output string
			}
			if err = decoder.Decode(&event); err == io.EOF {
				break
			} else if err != nil {
				cmd.Process.Kill()
				cmd.Wait()
				return err
			}
			fmt.Print(event.Output)
			if i, ok := byName[strings.SplitN(event.Test, "/", 2)[0]]; ok {
				switch event.Action {
				case "pass":
					if event.Test == cases[i].Name && cases[i].Result == "MISSING" {
						cases[i].Result = "PASS"
					}
				case "fail":
					cases[i].Result = "FAIL"
				case "skip":
					if cases[i].Result != "FAIL" {
						cases[i].Result = "SKIP"
					}
				}
			}
		}
		if err = cmd.Wait(); err != nil {
			return err
		}
		start = end
	}
	return nil
}

func main() {
	s := unitSummary{Schema: 1, Base: os.Getenv("CI_CHECK_BASE"), Status: "failed", Cases: []unitCase{}, ChangedFiles: []string{}, BroadenedPackages: []string{}}
	var err error
	if s.Base == "" {
		err = fmt.Errorf("CI_CHECK_BASE is required for changed unit tests")
	} else {
		s.ChangedFiles, err = changedFiles(s.Base)
	}
	if err == nil {
		s.Base = strings.TrimSpace(string(mustGitBase(s.Base)))
		s.Cases, s.BroadenedPackages, err = selectCases(s.ChangedFiles, s.Base)
	}
	if s.Cases == nil {
		s.Cases = []unitCase{}
	}
	if s.BroadenedPackages == nil {
		s.BroadenedPackages = []string{}
	}
	if err == nil && len(s.Cases) == 0 {
		s.Status = "not_applicable"
		s.Reason = "no_matching_unit_tests"
		goChange := false
		for _, f := range s.ChangedFiles {
			if strings.HasSuffix(f, ".go") || f == "go.mod" || f == "go.sum" {
				goChange = true
			}
		}
		if !goChange {
			s.Reason = "no_go_changes"
		}
	} else if err == nil {
		err = runCases(s.Cases)
	}
	finishUnitSummary(&s, err)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	b, _ := json.Marshal(s)
	fmt.Printf("CI_UNIT_TESTS_V1 %s\n", b)
	if s.Status == "failed" {
		os.Exit(1)
	}
}

// A skipped unit test is a completed, non-error outcome for quick CI.
// Keep its SKIP evidence while allowing the combined CI/unit-test gate to pass.
func finishUnitSummary(s *unitSummary, err error) {
	s.Passed, s.Failed, s.Skipped, s.Missing = 0, 0, 0, 0
	for _, c := range s.Cases {
		switch c.Result {
		case "PASS":
			s.Passed++
		case "FAIL":
			s.Failed++
		case "SKIP":
			s.Skipped++
		default:
			s.Missing++
		}
	}
	if err != nil || s.Failed+s.Missing > 0 {
		s.Status = "failed"
	} else if len(s.Cases) > 0 {
		s.Status = "passed"
	}
}

func mustGitBase(base string) []byte {
	b, _ := exec.Command("git", "rev-parse", "--verify", base+"^{commit}").Output()
	return b
}
