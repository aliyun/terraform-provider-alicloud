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
		{"test-like factory variable", old + "\nvar TestAccFactories = factory()\nfunc factory() int { return 1 }\nfunc TestUnitFactories(t *testing.T) { use(TestAccFactories) }\n", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := selectScope(fixture(t, map[string]string{file: old}), fixture(t, map[string]string{file: tc.source}), []string{file})
			if !reflect.DeepEqual(result.Packages["alicloud/service/demo"], tc.want) || len(result.Blocked) != 0 {
				t.Fatalf("unexpected scope: %+v", result)
			}
		})
	}
}

func TestAccScopeUnusedPrivateTestFunction(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/data_source_demo_test.go"
	const helper = "\nfunc privateHelper() { use(\"old\") }"
	const base = "package demo\nfunc TestAccDemo(t *testing.T) {}"
	for _, tc := range []struct {
		name, before, after string
		wantTest, blocked   bool
	}{
		{"added alongside new AccTest", "package demo", base + helper, true, false},
		{"edited", base + helper, base + strings.Replace(helper, "old", "new", 1), false, false},
		{"removed", base + helper, base, false, false},
		{"old caller", strings.Replace(base, "{}", "{ privateHelper() }", 1) + helper, base, true, false},
		{"function value caller", base + helper, strings.Replace(base, "{}", "{ use(privateHelper) }", 1) + helper, true, false},
		{"reference without a test", base, base + helper + "\nvar callback = privateHelper", true, false},
		{"method", base, base + "\ntype mock struct{}\nfunc (mock) privateHelper() {}", true, false},
		{"initializer", base, base + "\nvar config = loadConfig()", true, false},
		{"directive", base, base + "\n  //go:noinline" + helper, true, true},
		{"unresolved receiver", base, base + "\ntype mock struct{}\nfunc (mock) call() {}\nfunc TestAccMock(t *testing.T) { makeMock().call() }", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := selectScope(fixture(t, map[string]string{file: tc.before}),
				fixture(t, map[string]string{file: tc.after}), []string{file})
			gotTest := false
			for _, name := range result.Packages[dir] {
				gotTest = gotTest || name == "TestAccDemo"
			}
			if gotTest != tc.wantTest || (len(result.Blocked) != 0) != tc.blocked {
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
			if tc.name == "removed test" && !reflect.DeepEqual(result.AccOnlyBlockers, result.Blocked) {
				t.Fatalf("removed AccTest should only block acceptance: %+v", result)
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

func TestAccScopeSAEFrameworkAndRecorder(t *testing.T) {
	const file = "alicloud/service/sae/resource_demo_test.go"
	const source = `package sae_test
import (
 "testing"
 state "github.com/hashicorp/terraform-plugin-testing/terraform"
 fw "github.com/hashicorp/terraform-plugin-framework/provider"
 res "github.com/hashicorp/terraform-plugin-framework/resource"
 server "github.com/hashicorp/terraform-plugin-framework/providerserver"
)
type recorder struct{}
type client struct{}
func setup() (*client, *recorder) { return nil, nil }
func (r *recorder) captured() { recordOld() }
type provider struct{}
func (*provider) Configure(ctx Context, req fw.ConfigureRequest, resp *fw.ConfigureResponse) { providerOld() }
func (*provider) Resources(ctx Context) []Factory { resourcesOld(); return nil }
type resource struct{}
func (*resource) Configure(ctx Context, req res.ConfigureRequest, resp *res.ConfigureResponse) { resourceOld() }
func TestUnitProtocol(t *testing.T) { server.NewProtocol5(&provider{})() }
func TestAccRecorder(t *testing.T) { _, r := setup(); r.captured() }
func TestAccState(t *testing.T) { var s *state.State; root := s.RootModule(); use(root.Resources["x"]) }
func TestAccResource(t *testing.T) { r, ok := fromFactory().(res.ResourceWithConfigure); if ok { r.Configure(ctx, res.ConfigureRequest{}, nil) } }
`
	for _, tc := range []struct {
		change string
		want   []string
	}{
		{"recordOld()", []string{"TestAccRecorder"}},
		{"providerOld()", nil},
		{"resourcesOld()", nil},
		{"resourceOld()", []string{"TestAccResource"}},
	} {
		t.Run(tc.change, func(t *testing.T) {
			before := fixture(t, map[string]string{file: source})
			after := fixture(t, map[string]string{file: strings.Replace(source, tc.change, "updated()", 1)})
			scope := selectScope(before, after, []string{file})
			if len(scope.Blocked) != 0 || len(scope.FullPackages) != 0 || !reflect.DeepEqual(scope.Packages["alicloud/service/sae"], tc.want) {
				t.Fatalf("incorrect framework/recorder impact: %+v", scope)
			}
		})
	}
}

func TestProtocolOptionalCallbacksReachAcceptanceTests(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/resource_demo_test.go"
	const template = `package demo_test
import (
 fw "github.com/hashicorp/terraform-plugin-framework/provider"
 server "github.com/hashicorp/terraform-plugin-framework/providerserver"
)
type probe struct{ fw.Provider }
type other struct{ fw.Provider }
func (*probe) CALLBACK() { original() }
func TestUnitDirect(t *testing.T) { p := &probe{}; p.CALLBACK() }
func TestAccRegistered(t *testing.T) { server.CONSTRUCTOR(&probe{})() }
func TestAccUnrelated(t *testing.T) { server.CONSTRUCTOR(&other{})() }`
	for _, ctor := range []string{"NewProtocol5", "NewProtocol6", "NewProtocol5WithError", "NewProtocol6WithError"} {
		for _, method := range []string{"Functions", "EphemeralResources", "ListResources", "Actions", "StateStores", "MetaSchema", "ConfigValidators", "ValidateConfig"} {
			t.Run(ctor+"/"+method, func(t *testing.T) {
				source := strings.NewReplacer("CALLBACK", method, "CONSTRUCTOR", ctor).Replace(template)
				before := fixture(t, map[string]string{file: source})
				after := fixture(t, map[string]string{file: strings.Replace(source, "original()", "changed()", 1)})
				s := selectScope(before, after, []string{file})
				if len(s.Blocked) != 0 || len(s.FullPackages) != 0 || !reflect.DeepEqual(s.Packages[dir], []string{"TestAccRegistered"}) {
					t.Fatalf("optional callback lost or matched an unrelated provider: %+v", s)
				}
			})
		}
	}
}

func TestOptionalCallbackViaEmbeddedProviderOrFactory(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/resource_demo_test.go"
	const source = `package demo_test
import (
 "context"
 "testing"
 "github.com/hashicorp/terraform-plugin-framework/function"
 fw "github.com/hashicorp/terraform-plugin-framework/provider"
 server "github.com/hashicorp/terraform-plugin-framework/providerserver"
)
type probe struct{ fw.Provider }
func (*probe) Functions(context.Context) []func() function.Function { return nil }
type wrapper struct{ *probe }
func factory() fw.Provider { return &probe{} }
func TestUnitDirect(t *testing.T) { p := &probe{}; p.Functions(context.Background()) }
func TestAccRegistered(t *testing.T) { server.NewProtocol6(ARG)() }
func TestAccSibling(t *testing.T) {}`
	for _, arg := range []string{"&wrapper{}", "factory()"} {
		t.Run(arg, func(t *testing.T) {
			old := strings.Replace(source, "ARG", arg, 1)
			newSource := strings.Replace(old, "return nil", "return []func() function.Function{}", 1)
			scope := selectScope(fixture(t, map[string]string{file: old}), fixture(t, map[string]string{file: newSource}), []string{file})
			if len(scope.Blocked) != 0 || len(scope.FullPackages[dir]) == 0 ||
				!reflect.DeepEqual(scope.Packages[dir], []string{"TestAccRegistered", "TestAccSibling"}) {
				t.Fatalf("unresolved provider must include the entire candidate suite: %+v", scope)
			}
		})
	}
}

func TestChainedFrameworkRPCSelectsTheFailingAcceptanceCase(t *testing.T) {
	const dir = "alicloud/service/probe"
	const file = dir + "/resource_demo_test.go"
	const source = `package probe_test
import (
 "context"
 "testing"
 fw "github.com/hashicorp/terraform-plugin-framework/provider"
 "github.com/hashicorp/terraform-plugin-framework/function"
 "github.com/hashicorp/terraform-plugin-framework/providerserver"
 "github.com/hashicorp/terraform-plugin-go/tfprotov5"
)
type probe struct { fw.Provider; called bool }
func (*probe) Metadata(ctx context.Context, req fw.MetadataRequest, resp *fw.MetadataResponse) { resp.TypeName = "probe" }
func (p *probe) Functions(context.Context) []func() function.Function { p.called = true; return nil }
func TestUnitDirect(t *testing.T) {
 p := &probe{}
 if len(p.Functions(context.Background())) != 0 { t.Fatal("unexpected function") }
}
func TestAccProtocol(t *testing.T) {
 p := &probe{}
 _, err := providerserver.NewProtocol5(p)().GetFunctions(context.Background(), &tfprotov5.GetFunctionsRequest{})
 if err != nil || !p.called { t.Fatal("callback not observed", err) }
}`
	before := fixture(t, map[string]string{file: source})
	after := fixture(t, map[string]string{file: strings.Replace(source, "p.called = true", "p.called = false", 1)})
	scope := selectScope(before, after, []string{file})
	if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages[dir], []string{"TestAccProtocol"}) {
		t.Fatalf("chained RPC lost its provider callback: %+v", scope)
	}
}

func TestUnitHelperSelectionUsesRunnableGoEntrypoints(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/resource_demo_test.go"
	const source = `package demo
import "testing"
type T = testing.T
type testArg = *testing.T
type probe struct{}
func (probe) check() { old() }
func TESTNAME(t ARG) { probe{}.check() }
func TestAccSibling(t *testing.T) {}`
	for _, tc := range []struct {
		name, arg string
		unit      bool
	}{
		{"TestUnitAlias", "*T", true},
		{"TestUnitQualified", "*testing.T", true},
		{"TestUnitInvalid", "testArg", false},
		{"Testhelper", "*T", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			old := strings.NewReplacer("TESTNAME", tc.name, "ARG", tc.arg).Replace(source)
			changed := strings.Replace(old, "old()", "changed()", 1)
			scope := selectScope(fixture(t, map[string]string{file: old}), fixture(t, map[string]string{file: changed}), []string{file})
			if len(scope.Blocked) != 0 || (len(scope.FullPackages[dir]) == 0) != tc.unit {
				t.Fatalf("incorrect unit-only classification: %+v", scope)
			}
		})
	}
}

func TestFactoryReturningExternalInterfaceKeepsMethodConsumers(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/resource_demo_test.go"
	const source = `package demo_test
import (
 "context"
 "testing"
 res "github.com/hashicorp/terraform-plugin-framework/resource"
)
type probe struct{ res.Resource }
func (*probe) Metadata(ctx context.Context, req res.MetadataRequest, resp *res.MetadataResponse) { resp.TypeName = "old" }
func newThing() (int, res.Resource) { return 1, &probe{} }
func TestUnitDirect(t *testing.T) { p := &probe{}; p.Metadata(context.Background(), res.MetadataRequest{}, &res.MetadataResponse{}) }
func TestAccViaInterface(t *testing.T) { _, r := newThing(); r.Metadata(context.Background(), res.MetadataRequest{}, &res.MetadataResponse{}) }
func TestAccSibling(t *testing.T) {}`
	before := fixture(t, map[string]string{file: source})
	after := fixture(t, map[string]string{file: strings.Replace(source, `"old"`, `"new"`, 1)})
	scope := selectScope(before, after, []string{file})
	if len(scope.Blocked) != 0 || len(scope.FullPackages[dir]) == 0 ||
		!reflect.DeepEqual(scope.Packages[dir], []string{"TestAccSibling", "TestAccViaInterface"}) {
		t.Fatalf("external interface result silently lost its consumers: %+v", scope)
	}
}

func TestExternalFactoryDoesNotConflateSchemaFieldsAndMethodContracts(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/resource_demo_test.go"
	const source = `package demo_test
import (
 "context"
 "testing"
 fw "github.com/hashicorp/terraform-plugin-framework/provider"
 ds "github.com/hashicorp/terraform-plugin-framework/datasource"
 sdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)
type probe struct{ fw.Provider }
func (*probe) Schema(ctx context.Context, req fw.SchemaRequest, resp *fw.SchemaResponse) { old() }
func TestUnitProvider(t *testing.T) { p := &probe{}; p.Schema(context.Background(), fw.SchemaRequest{}, &fw.SchemaResponse{}) }
func sdkResource() *sdk.Resource { return &sdk.Resource{} }
func dataSource() ds.DataSource { return nil }
func TestAccSDKField(t *testing.T) { _ = sdkResource().Schema }
func TestAccDataSource(t *testing.T) { dataSource().Schema(context.Background(), ds.SchemaRequest{}, &ds.SchemaResponse{}) }`
	before := fixture(t, map[string]string{file: source})
	after := fixture(t, map[string]string{file: strings.Replace(source, "old()", "changed()", 1)})
	scope := selectScope(before, after, []string{file})
	if len(scope.Blocked) != 0 || len(scope.Packages) != 0 {
		t.Fatalf("unrelated Schema contracts selected: %+v", scope)
	}
}

func TestRepeatedAmbiguityProducesOnePackageReason(t *testing.T) {
	const dir = "alicloud/service/demo"
	a, b := dir+"/service_a.go", dir+"/service_b.go"
	files := map[string]string{
		a:                     "package demo; type A struct{}; func (A) Probe() int {return 1}",
		b:                     "package demo; type B struct{}; func (B) Probe() int {return 1}",
		dir + "/case_test.go": "package demo; func pick() interface{Probe() int} {return A{}}; func TestAccDemo(t *testing.T) {_ = pick().Probe()}",
	}
	before := fixture(t, files)
	for _, f := range []string{a, b} {
		files[f] = strings.Replace(files[f], "return 1", "return 2", 1)
	}
	scope := selectScope(before, fixture(t, files), []string{a, b})
	if reasons := scope.FullPackages[dir]; len(reasons) != 1 || !strings.Contains(reasons[0], "Probe") {
		t.Fatalf("expected one shared ambiguity: %v", reasons)
	}
}

func TestExpandedTestPackageKeepsDeletedCaseBlocked(t *testing.T) {
	const file = "alicloud/service/one/resource_demo_test.go"
	files := map[string]string{
		file:                                 "package one; func TestAccRemoved(t *testing.T) {}; func TestAccRemaining(t *testing.T) {}",
		"alicloud/service/two/other_test.go": `package two; import _ "github.com/aliyun/terraform-provider-alicloud/alicloud/service/one"; func TestAccUnrelated(t *testing.T) {}`,
	}
	before := fixture(t, files)
	files[file] = "package one; type helper struct{}; func (helper) Unknown() {}; func TestAccRemaining(t *testing.T) {}"
	s := selectScope(before, fixture(t, files), []string{file})
	if len(s.Blocked) != 1 || !strings.Contains(s.Blocked[0], "removed/renamed AccTest") ||
		len(s.FullPackages) != 1 || !reflect.DeepEqual(s.Packages, map[string][]string{"alicloud/service/one": {"TestAccRemaining"}}) {
		t.Fatalf("incorrect test-only fallback: %+v", s)
	}
}

func TestFallbackRetainsOldImportersAndStopsAtTestBinaries(t *testing.T) {
	const source = "alicloud/service/product/service.go"
	const bridge = "alicloud/service/bridge/imports.go"
	files := map[string]string{
		source:                              "package product; type Public struct{}; func Changed() int { return 1 }",
		bridge:                              `package bridge; import p "github.com/aliyun/terraform-provider-alicloud/alicloud/service/product"; type Public = p.Public`,
		"alicloud/service/client/client.go": "package client; type Public struct{}",
		"alicloud/service/client/client_test.go": `package client_test
import ("testing"; b "github.com/aliyun/terraform-provider-alicloud/alicloud/service/bridge")
type testArg = *testing.T
func TestAccClient(t testArg) { _ = b.Public{} }`,
		"alicloud/service/unrelated/demo_test.go": `package unrelated; import ("testing"; c "github.com/aliyun/terraform-provider-alicloud/alicloud/service/client"); func TestAccUnrelated(t *testing.T) { _ = c.Public{} }`,
	}
	before := fixture(t, files)
	files[source] = "package product; type Public struct{}; func Changed() int { return 2 }"
	files[bridge] = "package bridge; type Public struct{}"
	after := fixture(t, files)
	scope := selectScope(before, after, []string{source, bridge})
	want := map[string][]string{"alicloud/service/client": {"TestAccClient"}}
	if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages, want) ||
		!strings.Contains(strings.Join(scope.FullPackages["alicloud/service/client"], " "), "service.go:Changed") {
		t.Fatalf("lost old importer or followed a test-only import transitively: %+v", scope)
	}
}

func TestAccScopeUnknownReceiverBlocksCrossServiceExpansion(t *testing.T) {
	files := map[string]string{
		"alicloud/service/demo/service.go": `package demo
type Service struct{}
func (s *Service) Describe() string { return "old" }`,
		"alicloud/service/demo/resource_demo.go": `package demo
func resourceDemo() { svc := factory(); svc.Describe() }`,
		"alicloud/service/demo/resource_demo_test.go": `package demo_test
func TestAccDemo(t *testing.T) {}
func TestAccOtherInPackage(t *testing.T) {}`,
		"alicloud/service/wrapper/service.go":     `package wrapper; import _ "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"`,
		"alicloud/service/consumer/demo_test.go":  `package consumer; import _ "github.com/aliyun/terraform-provider-alicloud/alicloud/service/wrapper"; func TestAccConsumer(t *testing.T) {}`,
		"alicloud/service/unrelated/demo_test.go": `package unrelated; func TestAccUnrelated(t *testing.T) {}`,
	}
	before := fixture(t, files)
	file := "alicloud/service/demo/service.go"
	files[file] = strings.Replace(files[file], `"old"`, `"new"`, 1)
	result := selectScope(before, fixture(t, files), []string{file})
	message := strings.Join(result.Blocked, "\n")
	if len(result.Packages) != 0 || len(result.FullPackages) != 0 ||
		!strings.Contains(message, "would expand to 2 service packages") || !strings.Contains(message, "resource_demo.go:2") ||
		!strings.Contains(message, "receiver type could not be resolved") ||
		!strings.Contains(message, "Service.Describe (alicloud/service/demo/service.go)") {
		t.Fatalf("unknown receiver incorrectly treated as proven: %+v", result)
	}
}

func TestResourceSweeperDoesNotConnectOtherInitializers(t *testing.T) {
	const dir = "alicloud/service/demo/"
	const svc = dir + "service.go"
	const test = dir + "resource_demo_test.go"
	const registration = `func init() { res.AddTestSweepers("alicloud_demo", &res.Sweeper{Name:"alicloud_demo", F:sweepDemo}) }`
	files := map[string]string{
		dir + "register.go":      `package demo; func ServicePackage() conns.ServicePackage { return conns.ServicePackage{SDKResources: []conns.SDKResource{{TypeName:"alicloud_demo",Factory: NewDemo}}} }`,
		dir + "resource_demo.go": `package demo; func NewDemo() *Resource { return nil }`,
		svc:                      `package demo; func Sweep() { old() }`,
		test: `package demo_test
import (d "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"; res "github.com/hashicorp/terraform-plugin-testing/helper/resource")
func init() {}
` + registration + `
func sweepDemo() { d.Sweep() }
func TestAccDemo(t *testing.T) {}`,
		dir + "resource_other_test.go":            `package demo_test; func init() {}; func TestAccOther(t *testing.T) {}`,
		"alicloud/service/unrelated/case_test.go": `package unrelated; import _ "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"; func TestAccUnrelated(t *testing.T) {}`,
	}
	before := fixture(t, files)
	files[svc] = strings.Replace(files[svc], "old()", "changed()", 1)
	s := selectScope(before, fixture(t, files), []string{svc})
	if len(s.Blocked) != 0 || len(s.FullPackages) != 0 || !reflect.DeepEqual(s.ResourceFiles, []string{dir + "resource_demo.go"}) ||
		!reflect.DeepEqual(s.Packages["alicloud/service/demo"], []string{"TestAccDemo"}) {
		t.Fatalf("sweeper must select the registered resource: %+v", s)
	}
	before = fixture(t, files)
	files[test] = strings.Replace(files[test], registration, strings.Replace(registration, "{ res.", "{ sideEffect(); res.", 1), 1)
	s = selectScope(before, fixture(t, files), []string{test})
	if len(s.FullPackages) != 1 || !reflect.DeepEqual(s.Packages["alicloud/service/demo"], []string{"TestAccDemo", "TestAccOther"}) {
		t.Fatalf("unknown initialization must retain the package suite: %+v", s)
	}
}

func TestFallbackInServiceWithMultipleResourceOwners(t *testing.T) {
	const dir = "alicloud/service/demo"
	const helper = dir + "/service.go"
	files := map[string]string{
		helper:                        "package demo; func Changed() { old() }",
		dir + "/resource_one.go":      "package demo; func One() {}",
		dir + "/resource_one_test.go": "package demo_test; func TestAccOne(t *testing.T) {}",
		dir + "/resource_two.go":      "package demo; func Two() {}",
		dir + "/resource_two_test.go": "package demo_test; func TestAccTwo(t *testing.T) {}",
	}
	before := fixture(t, files)
	files[helper] = strings.Replace(files[helper], "old()", "changed()", 1)
	s := selectScope(before, fixture(t, files), []string{helper})
	if len(s.Packages) != 0 || len(s.FullPackages) != 0 || !strings.Contains(strings.Join(s.Blocked, " "), "unproven resource suites") {
		t.Fatalf("speculative multi-resource fallback must block: %+v", s)
	}
	testHelper := dir + "/mock_test.go"
	before = fixture(t, files)
	files[testHelper] = "package demo_test; var mock = unknownFactory()"
	s = selectScope(before, fixture(t, files), []string{testHelper})
	if len(s.Blocked) == 0 || !reflect.DeepEqual(s.AccOnlyBlockers, s.Blocked) {
		t.Fatalf("test-only ambiguity must not block documentation: %+v", s)
	}
	files[dir+"/resource_one.go"] = "package demo; func One() { Changed() }"
	files[dir+"/resource_two.go"] = "package demo; func Two() { Changed() }"
	files["alicloud/service/other/resource_three.go"] = `package other; import d "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"; func Three() { d.Changed() }`
	files["alicloud/service/other/resource_three_test.go"] = "package other; func TestAccThree(t *testing.T) {}"
	before = fixture(t, files)
	files[helper] = strings.Replace(files[helper], "changed()", "newBody()", 1)
	s = selectScope(before, fixture(t, files), []string{helper})
	want := map[string][]string{dir: {"TestAccOne", "TestAccTwo"}, "alicloud/service/other": {"TestAccThree"}}
	if len(s.Blocked) != 0 || len(s.FullPackages) != 0 || !reflect.DeepEqual(s.Packages, want) {
		t.Fatalf("all proven same- and cross-package consumers must run: %+v", s)
	}
}

func TestInternalTestExportIsVisibleOnlyInItsExternalTestPackage(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/export_test.go"
	const use = `import d "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"; `
	files := map[string]string{
		file:                                  "package demo; var Exported = internal",
		dir + "/demo_test.go":                 "package demo_test; " + use + "func TestAccActual(t *testing.T) { d.Exported() }; func TestAccOther(t *testing.T) {}",
		"alicloud/service/other/demo_test.go": "package other; " + use + "func TestAccElsewhere(t *testing.T) { d.Exported() }",
	}
	before := fixture(t, files)
	files[file] = strings.Replace(files[file], "internal", "changed", 1)
	s := selectScope(before, fixture(t, files), []string{file})
	if len(s.Blocked) != 0 || len(s.FullPackages) != 0 || !reflect.DeepEqual(s.Packages, map[string][]string{dir: {"TestAccActual"}}) {
		t.Fatalf("wrong external test consumers: %+v", s)
	}
}

func TestProtocolMethodOnUnitOnlyReceiver(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/demo_test.go"
	const src = `package demo_test
type stub struct{}
func (*stub) ApplyResourceChange() { old() }
func TestUnitProtocol(t *testing.T) { runServer(&stub{}) }
func TestAccOther(t *testing.T) {}`
	for _, sibling := range []string{"",
		"func TestAccViaUnit(t *testing.T) { TestUnitProtocol(t) }",
		"func TestAccAssert(t *testing.T) { _ = server().(*stub) }",
		"var shared = &stub{}; func TestAccShared(t *testing.T) { runServer(shared) }",
	} {
		files := map[string]string{file: src, dir + "/other_test.go": "package demo_test\n" + sibling}
		before := fixture(t, files)
		files[file] = strings.Replace(src, "old()", "changed()", 1)
		s := selectScope(before, fixture(t, files), []string{file})
		if len(s.Blocked) != 0 || (len(s.Packages[dir]) != 0) != (sibling != "") {
			t.Fatalf("incorrect protocol callback scope with sibling %q: %+v", sibling, s)
		}
	}
}

func TestAccScopeResolvesNamedPointerAssertion(t *testing.T) {
	const dir = "alicloud/service/demo"
	const file = dir + "/resource_demo_test.go"
	const source = `package demo
import (
	"sync"
	"testing"
)
type fixtureManager struct{}
func (*fixtureManager) cleanup() int { return 1 }
type otherManager struct{}
func (*otherManager) cleanup() int { return 2 }
var fixtures sync.Map
func TestAccDemo(t *testing.T) {
	if value, ok := fixtures.Load(t); ok {
		_ = value.(*fixtureManager).cleanup()
	}
}
func TestAccOther(t *testing.T) {
	if value, ok := fixtures.Load(t); ok {
		_ = value.(*otherManager).cleanup()
	}
}`
	before := fixture(t, map[string]string{file: source})
	after := fixture(t, map[string]string{file: strings.Replace(source, "return 1", "return 3", 1)})
	result := selectScope(before, after, []string{file})
	if len(result.Blocked) != 0 || len(result.FullPackages) != 0 ||
		!reflect.DeepEqual(result.Packages[dir], []string{"TestAccDemo"}) {
		t.Fatalf("pointer assertion did not isolate the affected test: %+v", result)
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
	if len(result.Packages) != 0 || len(result.Blocked) != 0 {
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

func TestAccScopeUsesConstructorResultsForSameNamedMethods(t *testing.T) {
	files := map[string]string{
		"alicloud/service/cr/service.go": `package cr
type CrServiceV2 struct{}
func NewCrServiceV2() *CrServiceV2 { return &CrServiceV2{} }
func (s *CrServiceV2) Describe() string { return "old" }`,
		"alicloud/service/cr/resource_cr.go": `package cr
func resourceCr() { svc := NewCrServiceV2(); svc.Describe() }`,
		"alicloud/service/cr/resource_cr_test.go": `package cr
import "testing"
func TestAccCr(t *testing.T) {}
func TestUnitCr(t *testing.T) { svc := NewCrServiceV2(); svc.Describe() }`,
		"alicloud/service/amqpopen/service.go": `package amqpopen
type AmqpServiceV2 struct{}
func NewAmqpServiceV2() *AmqpServiceV2 { return &AmqpServiceV2{} }
func (s *AmqpServiceV2) Describe() string { return "unrelated" }`,
		"alicloud/service/amqpopen/resource_amqp.go": `package amqpopen
func resourceAmqp() { svc := NewAmqpServiceV2(); svc.Describe() }`,
		"alicloud/service/amqpopen/resource_amqp_test.go": `package amqpopen
import "testing"
func TestAccAmqp(t *testing.T) {}`,
	}
	before := fixture(t, files)
	file := "alicloud/service/cr/service.go"
	files[file] = strings.Replace(files[file], `"old"`, `"new"`, 1)
	result := selectScope(before, fixture(t, files), []string{file})
	want := map[string][]string{"alicloud/service/cr": {"TestAccCr"}}
	if len(result.Blocked) != 0 || !reflect.DeepEqual(result.Packages, want) {
		t.Fatalf("constructor receiver scope = %+v, want only %v", result, want)
	}
}

func TestNewStaticRegistration(t *testing.T) {
	const dir = "alicloud/service/demo"
	const register = dir + "/register.go"
	const source = `package demo
import c "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
func ServicePackage() c.ServicePackage {
	return c.ServicePackage{Name: "Demo",
		Resources: []c.Resource{{TypeName: "alicloud_demo", Factory: NewDemo}},
		DataSources: []c.DataSource{{TypeName: "alicloud_demos", Factory: NewDemos}},
	}
}`
	files := map[string]string{
		dir + "/resource_demo.go":          `package demo; import "example.invalid/resource"; func NewDemo() resource.Resource { return nil }`,
		dir + "/resource_demo_test.go":     "package demo; func TestAccDemo(t *testing.T) {}",
		dir + "/data_source_demos.go":      `package demo; import "example.invalid/datasource"; func NewDemos() datasource.DataSource { return nil }`,
		dir + "/data_source_demos_test.go": "package demo; func TestAccDemos(t *testing.T) {}",
	}
	before := fixture(t, files)
	files[register] = source
	after := fixture(t, files)
	scope := selectScope(before, after, []string{register})
	if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages, map[string][]string{dir: {"TestAccDemo", "TestAccDemos"}}) ||
		!reflect.DeepEqual(scope.ResourceFiles, []string{dir + "/data_source_demos.go", dir + "/resource_demo.go"}) {
		t.Fatalf("new static registration: %+v", scope)
	}
	for _, tc := range []struct{ old, new, blocker string }{
		{`Name: "Demo"`, "Name: productName()", "unsupported registration change"},
		{"return c.ServicePackage", "initialize(); return c.ServicePackage", "unsupported registration change"},
		{"Factory: NewDemos", "Factory: selectFactory()", "unsupported registration change"},
		{"Factory: NewDemos", "Factory: Missing", "cannot resolve AccTest scope"},
	} {
		files[register] = strings.Replace(source, tc.old, tc.new, 1)
		scope := selectScope(before, fixture(t, files), []string{register})
		if !strings.Contains(strings.Join(scope.Blocked, "\n"), tc.blocker) {
			t.Errorf("%s: missing blocker %q: %+v", tc.new, tc.blocker, scope)
		}
	}
	files[register] = strings.Replace(source, `"Demo"`, `"Renamed"`, 1)
	if scope := selectScope(after, fixture(t, files), []string{register}); len(scope.Blocked) == 0 {
		t.Fatal("changing an existing product name requires review")
	}
}

func TestRegistrationAndHelperSelectOnlyActualDataSource(t *testing.T) {
	const dir = "alicloud/service/ecs/"
	const register = dir + "register.go"
	const helper = dir + "data_source_framework_helpers.go"
	const implementation = dir + "data_source_image_pipeline_executions.go"
	const shared = "alicloud/connectivity/client.go"
	files := map[string]string{
		register: `package ecs
func ServicePackage() conns.ServicePackage {
	return conns.ServicePackage{DataSources: []conns.DataSource{
		{TypeName: "alicloud_other", Factory: NewOther},
	}}
}`,
		helper:         `package ecs; func decode() { old() }`,
		implementation: `package ecs; import "example.invalid/datasource"; func NewExecutions() datasource.DataSource { decode(); return nil }`,
		implementation[:len(implementation)-3] + "_test.go": `package ecs; func TestAccExecutions(t *testing.T) {}`,
		dir + "data_source_other.go":                        `package ecs; import "example.invalid/datasource"; func NewOther() datasource.DataSource { return nil }`,
		dir + "data_source_other_test.go":                   `package ecs; func TestAccOther(t *testing.T) {}`,
		shared:                                              `package connectivity; func Do() { old() }`,
	}
	before := fixture(t, files)
	files[register] = strings.Replace(files[register], "{TypeName:", "{TypeName: \"alicloud_executions\", Factory: NewExecutions}, {TypeName:", 1)
	files[helper] = strings.Replace(files[helper], "old()", "changed()", 1)
	files[shared] = strings.Replace(files[shared], "old()", "changed()", 1)
	after := fixture(t, files)
	for _, changed := range [][]string{{register}, {helper}, {register, helper, shared}} {
		result := selectScope(before, after, changed)
		if len(result.Blocked) != 0 || !reflect.DeepEqual(result.Packages, map[string][]string{dir[:len(dir)-1]: {"TestAccExecutions"}}) ||
			!reflect.DeepEqual(result.ResourceFiles, []string{implementation}) {
			t.Fatalf("changes %v selected %+v", changed, result)
		}
	}
	result := selectScope(before, after, []string{shared})
	if len(result.Blocked) != 0 || len(result.Packages) != 0 || len(result.ResourceFiles) != 0 {
		t.Fatalf("shared change selected %+v", result)
	}
}

func TestAccScopeResolvesEmbeddedLocalInterfaceAndStructField(t *testing.T) {
	const service = "alicloud/service/demo/service.go"
	files := map[string]string{
		service: `package demo
type Reader interface { Describe(id string) error }
type Service interface { Reader; Update(string) error }
type Client struct{}
func (*Client) Update(string) error { return nil }
func (*Client) Describe(string) error { return old() }
type Other struct{}
func (*Other) Describe(string) error { return nil }
type Wrapper struct { client *Client }
func (w Wrapper) Read() error { return w.client.Describe("id") }
func (Wrapper) Update(string) error { return nil }
func (w Wrapper) Describe(id string) error { return w.Read() }`,
		"alicloud/service/demo/resource_demo.go": `package demo
func resourceDemo(service Service) { service.Describe("id") }`,
		"alicloud/service/demo/resource_demo_test.go": `package demo
import "testing"
func TestAccDemo(t *testing.T) {}`,
		"alicloud/service/demo/resource_other.go": `package demo
func resourceOther(client *Other) { client.Describe("id") }`,
		"alicloud/service/demo/resource_other_test.go": `package demo
import "testing"
func TestAccOther(t *testing.T) {}`,
	}
	before := fixture(t, files)
	files[service] = strings.Replace(files[service], "old()", "changed()", 1)
	result := selectScope(before, fixture(t, files), []string{service})
	if len(result.Blocked) != 0 || !reflect.DeepEqual(result.Packages["alicloud/service/demo"], []string{"TestAccDemo"}) {
		t.Fatalf("unexpected interface impact: %+v", result)
	}
}

func TestAccScopeTracesHTTPTestHandlerAndItsHelpers(t *testing.T) {
	const file = "alicloud/service/demo/resource_demo_test.go"
	source := `package demo
import (
	"net/http"
	servers "net/http/httptest"
	"testing"
)
type fixture struct{}
func (f *fixture) ServeHTTP(w http.ResponseWriter, r *http.Request) { body() }
func body() { old() }
func read(f *fixture) { servers.NewServer(f) }
func TestUnitDemo(t *testing.T) { read(&fixture{}) }
func TestAccDemo(t *testing.T) {}
`
	before := fixture(t, map[string]string{file: source})
	changed := strings.NewReplacer("body() }", "body(); extra() }", "old()", "updated()").Replace(source)
	result := selectScope(before, fixture(t, map[string]string{file: changed}), []string{file})
	if len(result.Blocked) != 0 || len(result.Packages) != 0 {
		t.Fatalf("handler is only used by the unit test: %+v", result)
	}

	shadowed := strings.ReplaceAll(changed, "servers.NewServer(f)", "servers := otherFactory(); servers.NewServer(f)")
	// Compare two snapshots with the shadow in place to avoid passing via the old edge.
	prior := strings.ReplaceAll(source, "servers.NewServer(f)", "servers := otherFactory(); servers.NewServer(f)")
	result = selectScope(fixture(t, map[string]string{file: prior}), fixture(t, map[string]string{file: shadowed}), []string{file})
	if len(result.Blocked) != 0 || len(result.FullPackages["alicloud/service/demo"]) == 0 {
		t.Fatalf("shadowed alias established handler coverage: %+v", result)
	}
}

func TestAccScopeReportsIndirectCoverageGapWithoutBlockingOtherCases(t *testing.T) {
	const dir = "alicloud/service/rds/"
	const serviceFile = dir + "service.go"
	const legacy = dir + "resource_alicloud_rds_ddr_instance.go"
	sources := map[string]string{
		serviceFile:                              `package rds; type RdsService struct{}; func (s *RdsService) GetSecurityIps() int { return 1 }`,
		legacy:                                   `package rds; func readLegacy(s *RdsService) { s.GetSecurityIps() }`,
		dir + "resource_alicloud_db_instance.go": `package rds; func readInstance(s *RdsService) { s.GetSecurityIps() }`,
		dir + "resource_alicloud_db_instance_test.go": `package rds; func TestAccInstance(t *testing.T) {}`,
	}
	before := fixture(t, sources)
	sources[serviceFile] = strings.Replace(sources[serviceFile], "return 1", "return 2", 1)
	indirect := selectScope(before, fixture(t, sources), []string{serviceFile})
	if !reflect.DeepEqual(indirect.Packages[strings.TrimSuffix(dir, "/")], []string{"TestAccInstance"}) ||
		len(indirect.Blocked) != 0 || len(indirect.UnverifiedCoverage) != 1 ||
		!strings.Contains(indirect.UnverifiedCoverage[0], legacy) {
		t.Fatalf("expected runnable case and nonblocking legacy coverage diagnostic, got %+v", indirect)
	}

	sources[legacy] = strings.Replace(sources[legacy], "s.GetSecurityIps()", "s.GetSecurityIps(); s.GetSecurityIps()", 1)
	direct := selectScope(before, fixture(t, sources), []string{serviceFile, legacy})
	if len(direct.Blocked) != 1 || len(direct.CoverageBlockers) != 1 || len(direct.UnverifiedCoverage) != 0 ||
		!strings.Contains(direct.Blocked[0], legacy) {
		t.Fatalf("expected missing coverage on directly modified legacy resource to block, got %+v", direct)
	}
}

func TestFirstStaticListAndCompanionRepair(t *testing.T) {
	const dir = "alicloud/service/demo/"
	const register = dir + "register.go"
	const impl = dir + "data_source_things.go"
	const misplaced = dir + "resource_thing_test.go"
	const companion = dir + "data_source_things_test.go"
	const list = `DataSources: []conns.DataSource{{TypeName:"alicloud_things",Factory: NewThings}},`
	files := map[string]string{
		register:  `package demo; func ServicePackage() conns.ServicePackage { return conns.ServicePackage{Name:"Demo"} }`,
		impl:      `package demo; import "example.invalid/datasource"; func NewThings() datasource.DataSource { return nil }`,
		misplaced: `package demo_test; import "github.com/hashicorp/terraform-plugin-testing/helper/resource"; func TestAccThings(t *testing.T) { resource.TestCheckResourceAttr("data.alicloud_things.default","id","1") }`,
	}
	before := fixture(t, files)
	original := files[register]
	files[register] = strings.Replace(original, "Name:", list+"Name:", 1)
	scope := selectScope(before, fixture(t, files), []string{register})
	if len(scope.Blocked) != 1 || !strings.Contains(scope.Blocked[0], companion) || !strings.Contains(scope.Blocked[0], misplaced+":TestAccThings") || !strings.Contains(scope.Blocked[0], "data.alicloud_things") {
		t.Fatalf("missing actionable coverage diagnostic: %+v", scope)
	}
	files[companion] = files[misplaced]
	delete(files, misplaced)
	scope = selectScope(before, fixture(t, files), []string{register})
	if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages["alicloud/service/demo"], []string{"TestAccThings"}) {
		t.Fatalf("companion repair: %+v", scope)
	}
	files[register] = strings.Replace(original, "Name:", "DataSources: dynamic(), Name:", 1)
	scope = selectScope(before, fixture(t, files), []string{register})
	if len(scope.Blocked) == 0 || !strings.Contains(strings.Join(scope.Blocked, " "), "ServicePackage (DataSources)") {
		t.Fatalf("dynamic field: %+v", scope)
	}
	files[register] = strings.Replace(original, "Name:", list+"Name:", 1)
	files[register] = strings.Replace(files[register], `"Demo"`, `"other"`, 1)
	scope = selectScope(before, fixture(t, files), []string{register})
	if len(scope.Blocked) == 0 {
		t.Fatalf("mixed edit was not blocked: %+v", scope)
	}
}

func TestAccScopeTracesResourceCheckInitDescribeMethod(t *testing.T) {
	const dir = "alicloud/service/demo"
	const svcFile = dir + "/cross_service.go"
	files := map[string]string{
		svcFile: `package demo
type RamService struct{}
func NewRamService() *RamService { return &RamService{} }
func (s *RamService) DescribeRamRoleAttachment(id string) error { return oldRam() }
type Service struct{}
func NewService() *Service { return &Service{} }
func (s *Service) DescribeDemoItem(id string) error { return oldItem() }`,
		dir + "/resource_ram_test.go": `package demo_test
import (
	"testing"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"
	tfa "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
)
func TestAccRamAttachment(t *testing.T) {
	var v map[string]interface{}
	rc := resourceCheckInitWithDescribeMethod("alicloud_ram_role_attachment.default", &v, func() interface{} {
		return demo.NewRamService()
	}, "DescribeRamRoleAttachment")
	_ = rc
}
func TestAccDemoItem(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_demo_item.default"
	serviceFunc := func() interface{} {
		return demo.NewService()
	}
	rc := tfa.ResourceCheckInit(resourceId, &v, serviceFunc)
	_ = rc
}
func TestAccUnrelated(t *testing.T) {}`,
	}
	before := fixture(t, files)

	afterRam := map[string]string{}
	for k, v := range files {
		afterRam[k] = v
	}
	afterRam[svcFile] = strings.Replace(files[svcFile], "oldRam()", "newRam()", 1)
	scope := selectScope(before, fixture(t, afterRam), []string{svcFile})
	if len(scope.Blocked) != 0 || len(scope.FullPackages) != 0 ||
		!reflect.DeepEqual(scope.Packages, map[string][]string{dir: {"TestAccRamAttachment"}}) {
		t.Fatalf("resourceCheckInitWithDescribeMethod was not traced accurately: %+v", scope)
	}

	afterItem := map[string]string{}
	for k, v := range files {
		afterItem[k] = v
	}
	afterItem[svcFile] = strings.Replace(files[svcFile], "oldItem()", "newItem()", 1)
	scope = selectScope(before, fixture(t, afterItem), []string{svcFile})
	if len(scope.Blocked) != 0 || len(scope.FullPackages) != 0 ||
		!reflect.DeepEqual(scope.Packages, map[string][]string{dir: {"TestAccDemoItem"}}) {
		t.Fatalf("acctest.ResourceCheckInit was not traced accurately: %+v", scope)
	}
}

func TestAccScopeStandaloneSweepPackageDoesNotExpandOrBlock(t *testing.T) {
	const dir = "alicloud/service/demo"
	const svc = dir + "/service_alicloud_demo.go"
	const sweeper = "alicloud/sweep/sweeper_demo_test.go"
	const sweepInit = `func init() {
	res.AddTestSweepers("alicloud_demo", &res.Sweeper{
		Name: "alicloud_demo",
		F:    testSweepDemo,
		Dependencies: []string{"alicloud_dep"},
	})
}`
	files := map[string]string{
		dir + "/register.go":           `package demo; func ServicePackage() conns.ServicePackage { return conns.ServicePackage{SDKResources: []conns.SDKResource{{TypeName:"alicloud_demo",Factory: NewDemo}}} }`,
		dir + "/resource_demo.go":      `package demo; func NewDemo() *Resource { return nil }`,
		dir + "/resource_demo_test.go": `package demo_test; func TestAccDemo(t *testing.T) {}`,
		svc: `package demo
type Service struct{}
func NewService() *Service { return &Service{} }
func (s *Service) SweepDemo() { oldSweep() }`,
		sweeper: `package sweep_test
import (
	demo "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"
	res "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)
` + sweepInit + `
func testSweepDemo(region string) error { demo.NewService().SweepDemo(); return nil }`,
	}
	before := fixture(t, files)
	after := map[string]string{}
	for k, v := range files {
		after[k] = v
	}
	after[svc] = strings.Replace(files[svc], "oldSweep()", "newSweep()", 1)
	scope := selectScope(before, fixture(t, after), []string{svc})
	if len(scope.Blocked) != 0 || len(scope.Packages) != 0 || len(scope.FullPackages) != 0 {
		t.Fatalf("standalone sweeper must not block or expand to resource AccTests: %+v", scope)
	}

	// If a resource also calls SweepDemo, that resource's AccTests must still be selected.
	withResourceCaller := map[string]string{}
	for k, v := range files {
		withResourceCaller[k] = v
	}
	withResourceCaller[dir+"/resource_demo.go"] = `package demo; func NewDemo() *Resource { NewService().SweepDemo(); return nil }`
	beforeWithRes := fixture(t, withResourceCaller)
	withResourceCaller[svc] = strings.Replace(withResourceCaller[svc], "oldSweep()", "newSweep()", 1)
	scope = selectScope(beforeWithRes, fixture(t, withResourceCaller), []string{svc})
	if len(scope.Blocked) != 0 || !reflect.DeepEqual(scope.Packages, map[string][]string{dir: {"TestAccDemo"}}) {
		t.Fatalf("resource caller must still be selected alongside standalone sweeper: %+v", scope)
	}

	// Side effects or dynamic dependencies in alicloud/sweep init() must still block.
	for _, tc := range []struct {
		name, old, repl string
	}{
		{"side effect", "res.AddTestSweepers", "sideEffect(); res.AddTestSweepers"},
		{"dynamic deps", `[]string{"alicloud_dep"}`, "dynamicDeps()"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dirty := map[string]string{}
			for k, v := range files {
				dirty[k] = v
			}
			dirty[sweeper] = strings.Replace(files[sweeper], tc.old, tc.repl, 1)
			dirtyBefore := fixture(t, dirty)
			dirty[svc] = strings.Replace(dirty[svc], "oldSweep()", "newSweep()", 1)
			s := selectScope(dirtyBefore, fixture(t, dirty), []string{svc})
			if len(s.Blocked) == 0 {
				t.Fatalf("non-pure sweep init (%s) must remain blocked: %+v", tc.name, s)
			}
		})
	}
}

func TestAccScopeSkipsUnreferencedServiceMethodsSafely(t *testing.T) {
	const dir = "alicloud/service/demo"
	const svc = dir + "/service_alicloud_demo.go"
	base := map[string]string{
		svc: `package demo
type Service struct{}
func (s *Service) DescribeDead() string { return "old" }
func (s *Service) WaitForDead() string { return s.DescribeDead() }
func (s *Service) String() string { return "oldStr" }`,
		dir + "/resource_one.go":                  `package demo; func One() {}`,
		dir + "/resource_one_test.go":             `package demo_test; func TestAccOne(t *testing.T) {}`,
		dir + "/resource_two.go":                  `package demo; func Two() {}`,
		dir + "/resource_two_test.go":             `package demo_test; func TestAccTwo(t *testing.T) {}`,
		"alicloud/service/other/resource_test.go": `package other; import _ "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"; func TestAccOther(t *testing.T) {}`,
	}
	before := fixture(t, base)

	// Transitive dead service method chain with 0 external callers is skipped without expanding scope.
	deadAfter := map[string]string{}
	for k, v := range base {
		deadAfter[k] = v
	}
	deadAfter[svc] = strings.Replace(base[svc], `"old"`, `"new"`, 1)
	scope := selectScope(before, fixture(t, deadAfter), []string{svc})
	if len(scope.Blocked) != 0 || len(scope.Packages) != 0 || len(scope.FullPackages) != 0 ||
		!reflect.DeepEqual(scope.UnreferencedDeclarations, []string{svc + ":DescribeDead"}) {
		t.Fatalf("unreferenced dead service method should be skipped cleanly: %+v", scope)
	}

	// Standard/external interface method names (e.g. String) must NOT be waived.
	strAfter := map[string]string{}
	for k, v := range base {
		strAfter[k] = v
	}
	strAfter[svc] = strings.Replace(base[svc], `"oldStr"`, `"newStr"`, 1)
	scope = selectScope(before, fixture(t, strAfter), []string{svc})
	if len(scope.Blocked) == 0 {
		t.Fatalf("external interface method String() must not be waived as dead: %+v", scope)
	}

	// Unresolved/ambiguous caller on the same method name must still block.
	ambig := map[string]string{}
	for k, v := range base {
		ambig[k] = v
	}
	ambig[dir+"/resource_one.go"] = `package demo; func One() { unknown().DescribeDead() }`
	beforeAmbig := fixture(t, ambig)
	ambig[svc] = strings.Replace(ambig[svc], `"old"`, `"new"`, 1)
	scope = selectScope(beforeAmbig, fixture(t, ambig), []string{svc})
	if len(scope.Blocked) == 0 {
		t.Fatalf("ambiguous caller must not be waived as dead: %+v", scope)
	}

	// Service method referenced only by a disabled SkipTestAcc* helper chain is skipped cleanly.
	skipOnly := map[string]string{}
	for k, v := range base {
		skipOnly[k] = v
	}
	skipOnly[dir+"/resource_one_test.go"] = `package demo_test
import demo "github.com/aliyun/terraform-provider-alicloud/alicloud/service/demo"
func checkDeadDestroy() { (&demo.Service{}).DescribeDead() }
func SkipTestAccDead(t *testing.T) { checkDeadDestroy() }
func TestAccOne(t *testing.T) {}`
	beforeSkip := fixture(t, skipOnly)
	skipOnly[svc] = strings.Replace(skipOnly[svc], `"old"`, `"new"`, 1)
	scope = selectScope(beforeSkip, fixture(t, skipOnly), []string{svc})
	if len(scope.Blocked) != 0 || len(scope.Packages) != 0 ||
		!reflect.DeepEqual(scope.UnreferencedDeclarations, []string{svc + ":DescribeDead"}) {
		t.Fatalf("dead service method only reached by SkipTestAcc* should be skipped: %+v", scope)
	}
}

func TestAccScopeResolvesFactoryAliasesAndConcreteStructMismatch(t *testing.T) {
	const dir = "alicloud/service/vpc"
	const svc = dir + "/service_alicloud_vpc.go"
	files := map[string]string{
		svc: `package vpc
type VpcService struct{}
func NewVpcService() *VpcService { return &VpcService{} }
func (s *VpcService) WaitForRouteEntry() error { return oldWait() }
func (s *VpcService) DescribeRouteEntry() error { return oldDesc() }
type VpcServiceV2 struct{}
func NewVpcServiceV2() *VpcServiceV2 { return &VpcServiceV2{} }
type OtherService struct{}
func (s *OtherService) DescribeLegacy() error { return oldLegacy() }`,
		dir + "/acc_typealiases_test.go": `package vpc_test
import vpc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
var NewVpcService = vpc.NewVpcService
var NewVpcServiceV2 = vpc.NewVpcServiceV2`,
		dir + "/resource_route_entry_test.go": `package vpc_test
import (
	"testing"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
)
func TestAccRouteEntry_wait(t *testing.T) {
	svc := NewVpcService()
	_ = svc.WaitForRouteEntry()
}
func TestAccRouteEntry_desc(t *testing.T) {
	var v map[string]interface{}
	_ = acctest.ResourceCheckInitWithDescribeMethod("alicloud_route_entry.default", &v, func() interface{} {
		return NewVpcService()
	}, "DescribeRouteEntry")
}
func TestAccRouteEntry_typoOnV2(t *testing.T) {
	var v map[string]interface{}
	_ = acctest.ResourceCheckInitWithDescribeMethod("alicloud_legacy.default", &v, func() interface{} {
		return NewVpcServiceV2()
	}, "DescribeLegacy")
}
func TestAccRouteEntry_unrelated(t *testing.T) {}`,
	}
	before := fixture(t, files)

	afterWait := map[string]string{}
	for k, v := range files {
		afterWait[k] = v
	}
	afterWait[svc] = strings.Replace(files[svc], "oldWait()", "newWait()", 1)
	scope := selectScope(before, fixture(t, afterWait), []string{svc})
	if len(scope.Blocked) != 0 || len(scope.FullPackages) != 0 ||
		!reflect.DeepEqual(scope.Packages, map[string][]string{dir: {"TestAccRouteEntry_wait"}}) {
		t.Fatalf("factory alias direct call was not resolved: %+v", scope)
	}

	afterDesc := map[string]string{}
	for k, v := range files {
		afterDesc[k] = v
	}
	afterDesc[svc] = strings.Replace(files[svc], "oldDesc()", "newDesc()", 1)
	scope = selectScope(before, fixture(t, afterDesc), []string{svc})
	if len(scope.Blocked) != 0 || len(scope.FullPackages) != 0 ||
		!reflect.DeepEqual(scope.Packages, map[string][]string{dir: {"TestAccRouteEntry_desc"}}) {
		t.Fatalf("factory alias inside ResourceCheckInitWithDescribeMethod was not resolved: %+v", scope)
	}

	// A concrete struct (*VpcServiceV2) missing DescribeLegacy must NOT create an ambiguity edge onto OtherService.DescribeLegacy.
	afterLegacy := map[string]string{}
	for k, v := range files {
		afterLegacy[k] = v
	}
	afterLegacy[svc] = strings.Replace(files[svc], "oldLegacy()", "newLegacy()", 1)
	scope = selectScope(before, fixture(t, afterLegacy), []string{svc})
	if len(scope.Blocked) != 0 || len(scope.Packages) != 0 ||
		!reflect.DeepEqual(scope.UnreferencedDeclarations, []string{svc + ":DescribeLegacy"}) {
		t.Fatalf("concrete struct mismatch leaked ambiguity onto another receiver: %+v", scope)
	}
}
