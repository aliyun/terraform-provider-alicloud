package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fixture(t *testing.T, sources map[string]string) *snapshot {
	t.Helper()
	s := &snapshot{decls: map[string]*declaration{}, files: map[string]bool{}, assets: map[string]string{}}
	for file, source := range sources {
		s.files[file] = true
		if strings.HasSuffix(file, ".go") {
			if err := parseSource(s, file, []byte(source)); err != nil {
				t.Fatal(err)
			}
		} else {
			s.assets[file] = source
		}
	}
	return s
}

func TestAccScopeChangedTestAndTransitiveHelper(t *testing.T) {
	file := "alicloud/service/demo/resource_demo_test.go"
	old := `package demo_test
import "fmt"
func TestAccDemo_basic(t *testing.T) { config() }
func TestAccDemo_complete(t *testing.T) { unrelated() }
func config() string { return wrapper() }
func wrapper() string { return fmt.Sprintf("old") }
func unrelated() string { return "fixed" }
func TestUnitDemo(t *testing.T) { unrelated() }
`
	for _, tc := range []struct {
		name, source string
		want         []string
	}{
		{"one test", strings.Replace(old, "config() }", "config(); check() }", 1), []string{"TestAccDemo_basic"}},
		{"transitive helper", strings.Replace(old, `Sprintf("old")`, `Sprintf("new")`, 1), []string{"TestAccDemo_basic"}},
		{"unit only", strings.Replace(old, "func TestUnitDemo(t *testing.T) { unrelated() }", "func TestUnitDemo(t *testing.T) { unrelated(); check() }", 1), nil},
		{"comment only", old + "// ordinary comment\n", nil},
		{"new import for new test", strings.Replace(old, "import \"fmt\"", "import (\"fmt\"; \"os\")", 1) + "func TestAccDemo_new(t *testing.T) { os.Getenv(\"TEST\") }", []string{"TestAccDemo_new"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := selectScope(fixture(t, map[string]string{file: old}), fixture(t, map[string]string{file: tc.source}), []string{file})
			if !reflect.DeepEqual(result.Packages["alicloud/service/demo"], tc.want) || len(result.Blocked) != 0 {
				t.Fatalf("unexpected scope: %+v", result)
			}
		})
	}
}

func TestAccScopeServiceAndImplementation(t *testing.T) {
	files := map[string]string{
		"alicloud/service/demo/service_demo.go": `package demo
type DemoService struct{}
func (s DemoService) DescribeDemo() string { return "old" }
func (s DemoService) DescribeOther() string { return "unrelated" }`,
		"alicloud/service/demo/resource_demo.go": `package demo
func resourceDemo() { resourceReadDemo() }
func resourceReadDemo() { s := DemoService{}; s.DescribeDemo() }`,
		"alicloud/service/demo/resource_other.go": `package demo
func resourceOther() { s := DemoService{}; s.DescribeOther() }`,
		"alicloud/service/demo/resource_demo_test.go": `package demo_test
func TestAccDemo_basic(t *testing.T) {}
func TestAccDemo_complete(t *testing.T) {}`,
		"alicloud/service/demo/resource_other_test.go": `package demo_test
func TestAccOther_basic(t *testing.T) {}`,
	}
	before := fixture(t, files)
	file := "alicloud/service/demo/service_demo.go"
	files[file] = strings.Replace(files[file], `"old"`, `"new"`, 1)
	result := selectScope(before, fixture(t, files), []string{file})
	if !reflect.DeepEqual(result.Packages["alicloud/service/demo"], []string{"TestAccDemo_basic", "TestAccDemo_complete"}) || len(result.Blocked) != 0 {
		t.Fatalf("service selected unrelated cases or missed callers: %+v", result)
	}
	if !reflect.DeepEqual(result.ResourceFiles, []string{"alicloud/service/demo/resource_demo.go"}) {
		t.Fatalf("wrong document scope: %+v", result)
	}
}

func TestAccScopeSharedImportedHelperAndFixture(t *testing.T) {
	files := map[string]string{
		"alicloud/acctest/config.go": `package acctest
func Config() string { return "old" }`,
		"alicloud/service/demo/resource_demo_test.go": `package demo_test
import shared "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
func TestAccDemo_basic(t *testing.T) { shared.Config() }
func TestAccDemo_fixture(t *testing.T) { os.ReadFile(filepath.Join("testdata", "demo.tf")) }
func TestAccDemo_other(t *testing.T) {}`,
		"alicloud/service/demo/testdata/demo.tf": "old fixture",
	}
	before := fixture(t, files)
	file := "alicloud/acctest/config.go"
	files[file] = strings.Replace(files[file], `"old"`, `"new"`, 1)
	result := selectScope(before, fixture(t, files), []string{file})
	// Production files in the explicit test-support packages are test helpers.
	if len(result.Blocked) != 0 || !reflect.DeepEqual(result.Packages["alicloud/service/demo"], []string{"TestAccDemo_basic"}) {
		t.Fatalf("imported test helper was not traced: %+v", result)
	}
	before = fixture(t, files)
	file = "alicloud/service/demo/testdata/demo.tf"
	files[file] = "new fixture"
	result = selectScope(before, fixture(t, files), []string{file})
	if len(result.Blocked) != 0 || !reflect.DeepEqual(result.Packages["alicloud/service/demo"], []string{"TestAccDemo_fixture"}) {
		t.Fatalf("fixture was not traced: %+v", result)
	}
}

func TestAccScopeDoesNotWaiveUnresolvedOrDeletedCoverage(t *testing.T) {
	for _, tc := range []struct{ name, file, before, after string }{
		{"public helper", "alicloud/helper/request.go", `package helper; func Request() { old() }`, `package helper; func Request() { changed() }`},
		{"service without caller", "alicloud/service/demo/service_demo.go", `package demo; func Describe() { old() }`, `package demo; func Describe() { changed() }`},
		{"removed test", "alicloud/resource_demo_test.go", `package alicloud; func TestAccDemo(t *testing.T) {}`, `package alicloud`},
		{"removed implementation", "alicloud/resource_demo.go", `package alicloud; func resourceDemo() {}`, `package alicloud`},
		{"side-effect import", "alicloud/import_test.go", `package alicloud`, "package alicloud\nimport _ \"example.invalid/setup\""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := selectScope(fixture(t, map[string]string{tc.file: tc.before}), fixture(t, map[string]string{tc.file: tc.after}), []string{tc.file})
			if len(result.Blocked) == 0 {
				t.Fatalf("unresolved change silently waived: %+v", result)
			}
		})
	}
}

func TestAccScopeIndependentConfigConstantsDoNotSelectUnchangedCases(t *testing.T) {
	file := "alicloud/resource_demo_test.go"
	old := `package alicloud
const (
  configA = "old"
  configB = "unchanged"
)
func TestAccA(t *testing.T) { use(configA) }
func TestAccB(t *testing.T) { use(configB) }`
	result := selectScope(fixture(t, map[string]string{file: old}),
		fixture(t, map[string]string{file: strings.Replace(old, `"old"`, `"new"`, 1)}), []string{file})
	if len(result.Blocked) != 0 || !reflect.DeepEqual(result.Packages["alicloud"], []string{"TestAccA"}) {
		t.Fatalf("independent configs incorrectly grouped: %+v", result)
	}
}

func TestAccScopeUnknownReceiverIsBlocked(t *testing.T) {
	files := map[string]string{
		"alicloud/service/demo/service.go": `package demo
type Service struct{}
func (s *Service) Describe() string { return "old" }`,
		"alicloud/service/demo/resource_demo.go": `package demo
func resourceDemo() { svc := factory(); svc.Describe() }`,
		"alicloud/service/demo/resource_demo_test.go": `package demo_test
func TestAccDemo(t *testing.T) {}`,
	}
	before := fixture(t, files)
	file := "alicloud/service/demo/service.go"
	files[file] = strings.Replace(files[file], `"old"`, `"new"`, 1)
	result := selectScope(before, fixture(t, files), []string{file})
	if len(result.Blocked) == 0 || !strings.Contains(strings.Join(result.Blocked, "\n"), "receiver") {
		t.Fatalf("unknown receiver incorrectly treated as proven: %+v", result)
	}
}

func TestAccScopeTestOnlyTypesCannotBeCalledByOtherPackageTests(t *testing.T) {
	files := map[string]string{
		"alicloud/service/demo/resource_demo_test.go": `package demo_test
type provider struct{}
func (p provider) Schema() { old() }
func TestAccDemo(t *testing.T) { p := provider{}; p.Schema() }
var TestAccFactories = factory
func factory() {}`,
		"alicloud/service/other/resource_other_test.go": `package other_test
func TestAccUnrelated(t *testing.T) { unknown.Schema() }`,
	}
	before := fixture(t, files)
	file := "alicloud/service/demo/resource_demo_test.go"
	files[file] = strings.Replace(files[file], "old()", "changed()", 1)
	result := selectScope(before, fixture(t, files), []string{file})
	if len(result.Blocked) != 0 || !reflect.DeepEqual(result.Packages, map[string][]string{"alicloud/service/demo": {"TestAccDemo"}}) {
		t.Fatalf("test-only method escaped its package: %+v", result)
	}
}

func TestAccScopeTestMainAndLocalShadow(t *testing.T) {
	file := "alicloud/resource_demo_test.go"
	old := `package alicloud
func TestMain(m *testing.M) { setup() }
func shared() string { return "old" }
func TestAccDemo(t *testing.T) { shared := func(){}; shared() }
func TestAccOther(t *testing.T) {}`
	before := fixture(t, map[string]string{file: old})
	result := selectScope(before, fixture(t, map[string]string{file: strings.Replace(old, `"old"`, `"new"`, 1)}), []string{file})
	if len(result.Packages) != 0 || len(result.Blocked) != 1 {
		t.Fatalf("local variable mistaken for shared helper: %+v", result)
	}
	result = selectScope(before, fixture(t, map[string]string{file: strings.Replace(old, "setup()", "newSetup()", 1)}), []string{file})
	if !reflect.DeepEqual(result.Packages["alicloud"], []string{"TestAccDemo", "TestAccOther"}) {
		t.Fatalf("TestMain did not select package cases: %+v", result)
	}
}

func TestAccScopeReadsCommittedSnapshotsAndDeletedPaths(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	run := func(args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgsign=false"}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %s: %v", args, out, err)
		}
		return strings.TrimSpace(string(out))
	}
	run("init", "-q")
	run("config", "user.name", "Scope fixture")
	run("config", "user.email", "scope@example.invalid")
	file := "alicloud/resource_demo_test.go"
	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
		t.Fatal(err)
	}
	old := "package alicloud\nfunc TestAccDemo(t *testing.T) { old() }\nfunc TestAccUnchanged(t *testing.T) {}\n"
	if err := os.WriteFile(file, []byte(old), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("alicloud/keep.go", []byte("package alicloud\n"), 0600); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-qm", "base")
	base := run("rev-parse", "HEAD")
	if err := os.WriteFile(file, []byte(strings.Replace(old, "old()", "newCall()", 1)), 0600); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-qm", "change one test")
	head := run("rev-parse", "HEAD")
	before, err := readSnapshot(base)
	if err != nil {
		t.Fatal(err)
	}
	after, err := readSnapshot(head)
	if err != nil {
		t.Fatal(err)
	}
	files, err := command("diff", "--no-renames", "--name-only", base, head)
	if err != nil {
		t.Fatal(err)
	}
	result := selectScope(before, after, strings.Fields(string(files)))
	if !reflect.DeepEqual(result.Packages["alicloud"], []string{"TestAccDemo"}) || len(result.Blocked) != 0 {
		t.Fatalf("committed change lost exact selection: %+v", result)
	}
	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-qm", "remove coverage")
	deleted := run("rev-parse", "HEAD")
	files, err = command("diff", "--no-renames", "--name-only", head, "HEAD")
	if err != nil || !strings.Contains(string(files), file) {
		t.Fatalf("deleted path omitted: %s %v", files, err)
	}
	after, err = readSnapshot(deleted)
	if err != nil {
		t.Fatal(err)
	}
	result = selectScope(before, after, strings.Fields(string(files)))
	if len(result.Blocked) != 2 {
		t.Fatalf("deleted cases were silently waived: %+v", result)
	}
}

func TestAccScopeResolvesCrossPackageReceiversWithoutSelectingOtherTypes(t *testing.T) {
	files := map[string]string{
		"alicloud/service/demo/service.go": `package demo
type Service struct{}
func (s *Service) Describe() string { return "old" }
type OtherService struct{}
func (s OtherService) Describe() string { return "unrelated" }`,
		"alicloud/service/consumer/resource_demo.go": `package consumer
import api "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"
func resourceDemo() { svc := &api.Service{}; svc.Describe() }`,
		"alicloud/service/consumer/resource_other.go": `package consumer
import api "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"
func resourceOther() { svc := api.OtherService{}; svc.Describe() }`,
		"alicloud/service/consumer/resource_demo_test.go": `package consumer_test
func TestAccDemo(t *testing.T) {}`,
		"alicloud/service/consumer/resource_other_test.go": `package consumer_test
func TestAccOther(t *testing.T) {}`,
	}
	before := fixture(t, files)
	file := "alicloud/service/demo/service.go"
	files[file] = strings.Replace(files[file], `"old"`, `"new"`, 1)
	result := selectScope(before, fixture(t, files), []string{file})
	if len(result.Blocked) != 0 || !reflect.DeepEqual(result.Packages["alicloud/service/consumer"], []string{"TestAccDemo"}) {
		t.Fatalf("static receiver was not resolved: %+v", result)
	}
}
