package main

import (
	"reflect"
	"strings"
	"testing"
)

func layoutFixture() map[string]string {
	return map[string]string{
		"alicloud/service/demo/register.go": `package demo
func ServicePackage() conns.ServicePackage { return conns.ServicePackage{
 Resources: []conns.Resource{{TypeName: "alicloud_demo", Factory: NewDemo}, {TypeName: "alicloud_demo_alias", Factory: NewDemo}},
} }`,
		"alicloud/service/demo/resource_demo.go": `package demo
type Resource struct{}
func NewDemo() Resource { return Resource{} }`,
		"alicloud/service/demo/resource_demo_test.go": `package demo
import "testing"
func TestAccDemo(t *testing.T) {}
func TestUnitDemo(t *testing.T) {}`,
	}
}

func TestLayoutNewFilesHaveActionableDestinations(t *testing.T) {
	for _, tc := range []struct{ file, code, target string }{
		{"resource_demo_helpers.go", "LAYOUT001", "resource_demo.go"},
		{"resource_demo_decode.go", "LAYOUT001", "resource_demo.go"},
		{"resource_demo_unit_test.go", "LAYOUT002", "resource_demo_test.go"},
		{"resource_demo_helpers_test.go", "LAYOUT002", "resource_demo_test.go"},
		{"resource_unknown.go", "LAYOUT004", "register.go"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			files := layoutFixture()
			before := fixture(t, files)
			file := "alicloud/service/demo/" + tc.file
			files[file] = "package demo; func helper() {}"
			issues := layoutIssues(before, fixture(t, files), []string{file})
			if len(issues) != 1 || issues[0].Code != tc.code || !strings.Contains(issues[0].Fix, tc.target) || issues[0].Guide != layoutGuide || issues[0].Line != 1 {
				t.Fatalf("not actionable: %+v", issues)
			}
		})
	}
}

func TestLayoutGrandfathersExistingSplitFiles(t *testing.T) {
	files := layoutFixture()
	file := "alicloud/service/demo/resource_demo_helpers.go"
	files[file] = "package demo; func helper() { old() }"
	before := fixture(t, files)
	files[file] = "package demo; func helper() { changed() }"
	if got := layoutIssues(before, fixture(t, files), []string{file}); len(got) != 0 {
		t.Fatal(got)
	}
	delete(files, file)
	if got := layoutIssues(before, fixture(t, files), []string{file}); len(got) != 0 {
		t.Fatal(got)
	}
}

func TestLayoutKeepsInternalUnitExternalAcceptanceBoundary(t *testing.T) {
	for _, external := range []bool{false, true} {
		files := layoutFixture()
		if external {
			files["alicloud/service/demo/resource_demo_test.go"] = "package demo_test; import \"testing\"; func TestAccDemo(t *testing.T) {}"
		}
		before := fixture(t, files)
		file := "alicloud/service/demo/resource_demo_unit_test.go"
		files[file] = "package demo; import \"testing\"; func TestUnitDemo(t *testing.T) {}"
		issues := layoutIssues(before, fixture(t, files), []string{file})
		if (len(issues) == 0) != external {
			t.Fatalf("external=%v: %+v", external, issues)
		}
	}
}

func TestLayoutRegisteredHelpersNameAndSharedFilesAreLegal(t *testing.T) {
	files := layoutFixture()
	before := fixture(t, files)
	files["alicloud/service/demo/register.go"] = strings.Replace(files["alicloud/service/demo/register.go"], "TypeName: \"alicloud_demo\"", "TypeName: \"alicloud_demo_helpers\"", 1)
	file := "alicloud/service/demo/resource_demo_helpers.go"
	files[file] = files["alicloud/service/demo/resource_demo.go"]
	delete(files, "alicloud/service/demo/resource_demo.go")
	files["alicloud/service/demo/service.go"] = "package demo; func Shared() {}"
	files["alicloud/service/demo/service_query.go"] = "package demo; func Query() {}"
	if got := layoutIssues(before, fixture(t, files), []string{file, "alicloud/service/demo/service.go", "alicloud/service/demo/service_query.go"}); len(got) != 0 {
		t.Fatal(got)
	}
}

func TestLayoutNewLongEntrypointPreservesPublicNames(t *testing.T) {
	files := layoutFixture()
	before := fixture(t, map[string]string{})
	file := "alicloud/service/demo/resource_alicloud_demo.go"
	files[file] = files["alicloud/service/demo/resource_demo.go"]
	delete(files, "alicloud/service/demo/resource_demo.go")
	issues := layoutIssues(before, fixture(t, files), []string{file})
	if len(issues) != 1 || issues[0].Code != "LAYOUT003" || !strings.Contains(issues[0].Fix, "Preserve every registered TypeName") {
		t.Fatal(issues)
	}
}

func TestLayoutUnitSuffixCanBePartOfRegisteredEntityName(t *testing.T) {
	files := layoutFixture()
	before := fixture(t, map[string]string{})
	for _, suffix := range []string{".go", "_test.go"} {
		old := "alicloud/service/demo/resource_demo" + suffix
		files["alicloud/service/demo/resource_storage_unit"+suffix] = files[old]
		delete(files, old)
	}
	file := "alicloud/service/demo/resource_storage_unit_test.go"
	if issues := layoutIssues(before, fixture(t, files), []string{file}); len(issues) != 0 {
		t.Fatalf("real companion mistaken for split unit test: %+v", issues)
	}
}

func TestScopeHelperSpellingDoesNotCreateEntities(t *testing.T) {
	for _, name := range []string{"resource_demo_helpers.go", "data_source_alicloud_demo_helpers.go", "data_source_framework_helpers.go", "decode.go"} {
		t.Run(name, func(t *testing.T) {
			files := layoutFixture()
			file := "alicloud/service/demo/" + name
			files[file] = "package demo; func decode() { old() }"
			files["alicloud/service/demo/resource_demo.go"] = strings.Replace(files["alicloud/service/demo/resource_demo.go"], "return Resource{}", "decode(); return Resource{}", 1)
			before := fixture(t, files)
			files[file] = "package demo; func decode() { changed() }"
			scope := selectScope(before, fixture(t, files), []string{file})
			if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages["alicloud/service/demo"], []string{"TestAccDemo"}) || !reflect.DeepEqual(scope.ResourceFiles, []string{"alicloud/service/demo/resource_demo.go"}) {
				t.Fatalf("%+v", scope)
			}
		})
	}
}

func TestScopeUnknownHelperSelectsFullPackage(t *testing.T) {
	files := layoutFixture()
	file := "alicloud/service/demo/data_source_alicloud_demo_helpers.go"
	files[file] = "package demo; func decode() { old() }"
	before := fixture(t, files)
	files[file] = "package demo; func decode() { changed() }"
	scope := selectScope(before, fixture(t, files), []string{file})
	if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages["alicloud/service/demo"], []string{"TestAccDemo"}) ||
		!strings.Contains(strings.Join(scope.FullPackages["alicloud/service/demo"], " "), "production helper consumers") {
		t.Fatalf("%+v", scope)
	}
}

func TestLayoutPrivateSplitIsDetectedWithoutHelperSuffix(t *testing.T) {
	files := layoutFixture()
	before := fixture(t, files)
	file := "alicloud/service/demo/decoding.go"
	files[file] = "package demo; func decode() {}"
	files["alicloud/service/demo/resource_demo.go"] = strings.Replace(files["alicloud/service/demo/resource_demo.go"], "return Resource{}", "decode(); return Resource{}", 1)
	issues := layoutIssues(before, fixture(t, files), []string{file})
	if len(issues) != 1 || issues[0].Code != "LAYOUT001" || !strings.Contains(issues[0].Fix, "resource_demo.go") {
		t.Fatalf("unnamed split escaped policy: %+v", issues)
	}
}

func TestScopeFunctionAdapterMockIsUnitOnly(t *testing.T) {
	files := layoutFixture()
	files["alicloud/service/demo/resource_demo.go"] += `
type Caller interface { RpcPostWithContext(ctx context.Context) error }
func call(c Caller) error { return c.RpcPostWithContext(context.Background()) }`
	file := "alicloud/service/demo/resource_demo_test.go"
	files[file] = `package demo
import ("testing"; "context")
type rpcFunc func(context.Context) error
func(f rpcFunc) RpcPostWithContext(ctx context.Context) error { return f(ctx) }
func TestUnitDeadline(t *testing.T) { call(rpcFunc(func(context.Context) error { return nil })) }`
	before := fixture(t, files)
	files[file] = strings.Replace(files[file], "return f(ctx)", "changed(); return f(ctx)", 1)
	scope := selectScope(before, fixture(t, files), []string{file})
	if len(scope.Blocked) != 0 || len(scope.Packages) != 0 {
		t.Fatalf("mock incorrectly needs ACC: %+v", scope)
	}
}

func TestScopeRenameAndHelperMergeKeepCoverage(t *testing.T) {
	for _, rename := range []bool{false, true} {
		files := layoutFixture()
		helper := "alicloud/service/demo/resource_demo_helpers.go"
		files[helper] = "package demo; func decode() {}"
		files["alicloud/service/demo/resource_demo.go"] = strings.Replace(files["alicloud/service/demo/resource_demo.go"], "return Resource{}", "decode(); return Resource{}", 1)
		before := fixture(t, files)
		files["alicloud/service/demo/resource_demo.go"] += "\nfunc decode() {}"
		delete(files, helper)
		changed := []string{helper, "alicloud/service/demo/resource_demo.go"}
		if rename {
			for _, suffix := range []string{".go", "_test.go"} {
				old := "alicloud/service/demo/resource_demo" + suffix
				new := "alicloud/service/demo/resource_local" + suffix
				files[new] = files[old]
				delete(files, old)
				changed = append(changed, old, new)
			}
		}
		scope := selectScope(before, fixture(t, files), changed)
		if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages["alicloud/service/demo"], []string{"TestAccDemo"}) {
			t.Fatalf("rename=%v %+v", rename, scope)
		}
	}
}

func TestScopeSharedHelperRetainsEveryRegisteredConsumer(t *testing.T) {
	files := layoutFixture()
	register := "alicloud/service/demo/register.go"
	files[register] = strings.Replace(files[register], "Resources: []conns.Resource{", "Resources: []conns.Resource{{TypeName: \"alicloud_other\", Factory: NewOther},", 1)
	files["alicloud/service/demo/resource_demo.go"] += "\nfunc readDemo() { shared() }"
	files["alicloud/service/demo/resource_other.go"] = "package demo; func NewOther() Resource { shared(); return Resource{} }"
	files["alicloud/service/demo/resource_other_test.go"] = "package demo; import \"testing\"; func TestAccOther(t *testing.T) {}"
	file := "alicloud/service/demo/service.go"
	files[file] = "package demo; func shared() { old() }"
	before := fixture(t, files)
	files[file] = "package demo; func shared() { changed() }"
	scope := selectScope(before, fixture(t, files), []string{file})
	if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages["alicloud/service/demo"], []string{"TestAccDemo", "TestAccOther"}) {
		t.Fatalf("%+v", scope)
	}
}

func TestScopeInvalidReplacementDoesNotHideRemovedAccTest(t *testing.T) {
	files := layoutFixture()
	before := fixture(t, files)
	old := "alicloud/service/demo/resource_demo_test.go"
	delete(files, old)
	files["alicloud/service/demo/other_test.go"] = "package demo; func TestAccDemo() {}"
	scope := selectScope(before, fixture(t, files), []string{old, "alicloud/service/demo/other_test.go"})
	if !strings.Contains(strings.Join(scope.Blocked, " "), "removed/renamed AccTest") {
		t.Fatalf("invalid replacement waived: %+v", scope)
	}
}
