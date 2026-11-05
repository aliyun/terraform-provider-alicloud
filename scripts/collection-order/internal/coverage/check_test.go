package coverage

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func fixture(steps string) string {
	return `package alicloud
import (
 "testing"
 "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)
func TestAccExample(t *testing.T) {
 resourceId := "alicloud_example.default"
 testAccConfig := resourceTestAccConfigFunc(resourceId, "test", dependence)
 resource.Test(t, resource.TestCase{Steps: []resource.TestStep{` + steps + `}})
}`
}

func TestListReorderCoverage(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true}}}
	apply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"},"name":"same"})},`
	plan := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	reorderedApply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}})},`
	valid := fixture(apply + plan + reorderedApply)
	cases := []struct {
		name, src string
		pass      bool
	}{
		{"diff and convergent apply", valid, true},
		{"sdk v2 harness import still discovered", strings.Replace(valid, "terraform-plugin-testing", "terraform-plugin-sdk/v2", 1), true},
		{"convergent apply without plan", fixture(apply + reorderedApply), true},
		{"explicit convergent apply", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", "{ExpectNonEmptyPlan:false,Config:", 1)), true},
		{"apply retained builder snapshot", fixture(apply + plan + `{Config:testAccConfig(map[string]interface{}{})},`), true},
		{"clean plan only step between applies", fixture(apply + strings.Replace(plan, "ExpectNonEmptyPlan:true", "ExpectNonEmptyPlan:false", 1) + reorderedApply), true},
		{"plan only step without ExpectNonEmptyPlan", fixture(apply + strings.Replace(plan, ",ExpectNonEmptyPlan:true", "", 1) + reorderedApply), true},
		{"repeated plan only step", fixture(apply + plan + plan + reorderedApply), true},
		{"diff without apply", fixture(apply + plan), false},
		{"nonconvergent apply", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", "{ExpectNonEmptyPlan:true,Config:", 1)), false},
		{"apply original order", fixture(apply + plan + apply), false},
		{"apply changed membership", fixture(apply + plan + strings.Replace(reorderedApply, `"b","a"`, `"b","c"`, 1)), false},
		{"apply changed sibling", fixture(apply + plan + strings.Replace(reorderedApply, `"ids":`, `"name":"changed","ids":`, 1)), false},
		{"plan changed membership", fixture(apply + strings.Replace(plan, `"b","a"`, `"b","c"`, 1) + reorderedApply), true},
		{"plan changed sibling", fixture(apply + strings.Replace(plan, `"ids":`, `"name":"changed","ids":`, 1) + reorderedApply), false},
		{"unchanged plan step", fixture(apply + strings.Replace(plan, `"b","a"`, `"a","b"`, 1) + reorderedApply), true},
		{"intervening import", fixture(apply + plan + `{ImportState:true},` + reorderedApply), false},
		{"apply is plan only", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", "{PlanOnly:true,Config:", 1)), false},
		{"apply error", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", "{ExpectError:err,Config:", 1)), false},
		{"apply skipped", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", "{SkipFunc:skip,Config:", 1)), false},
		{"tainted initial apply", fixture(strings.Replace(apply, "{Config:", `{Taint:[]string{"alicloud_example.default"},Config:`, 1) + plan + reorderedApply), false},
		{"tainted reorder plan", fixture(apply + strings.Replace(plan, "{Config:", `{Taint:[]string{"alicloud_example.default"},Config:`, 1) + reorderedApply), false},
		{"tainted final apply", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", `{Taint:[]string{"alicloud_example.default"},Config:`, 1)), false},
		{"initial Check skips", fixture(strings.Replace(apply, "{Config:", `{Check:func(*terraform.State)error{t.Skip("disabled");return nil},Config:`, 1) + plan + reorderedApply), false},
		{"final Check skips", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", `{Check:func(*terraform.State)error{t.SkipNow();return nil},Config:`, 1)), false},
		{"local PreCheck skips", strings.Replace(strings.Replace(valid, "resource.Test(t,", `preCheck:=func(){t.Skip("disabled")}; resource.Test(t,`, 1), "resource.TestCase{Steps:", "resource.TestCase{PreCheck:preCheck,Steps:", 1), false},
		{"region PreCheck helper", strings.Replace(valid, "resource.TestCase{Steps:", `resource.TestCase{PreCheck:func(){testAccPreCheckWithRegions(t,true,[]connectivity.Region{connectivity.Beijing})},Steps:`, 1), true},
		{"inline region guard unsupported", strings.Replace(valid, "resource.Test(t,", `if !isRegionSupported(){t.Skip("region")}; resource.Test(t,`, 1), false},
		{"initial apply nonconvergent", fixture(strings.Replace(apply, "{Config:", "{ExpectNonEmptyPlan:true,Config:", 1) + plan + reorderedApply), false},
		{"comment cannot waive coverage", "// collection-order: ordered ids -- API ordering rationale.\n" + fixture(apply), false},
		{"duplicates only", strings.ReplaceAll(valid, `"b"`, `"a"`), false},
		{"comment only", fixture(apply + "/*" + plan + reorderedApply + "*/"), false},
		{"go build header", "//go:build ignore\n\n" + valid, false},
		{"legacy build header", "// +build ignore\n\n" + valid, false},
		{"go build mention after package", strings.Replace(valid, "func TestAcc", "//go:build example\nfunc TestAcc", 1), true},
		{"legacy build mention after package", strings.Replace(valid, "func TestAcc", "// +build example\nfunc TestAcc", 1), true},
		{"ordinary header comment", "// go:build example\n\n" + valid, true},
		{"other resource", strings.Replace(valid, "alicloud_example.default", "alicloud_other.default", 1), false},
		{"not TestAcc", strings.Replace(valid, "TestAccExample", "TestUnitExample", 1), false},
		{"unused steps", strings.Replace(valid, "resource.Test(t, resource.TestCase{Steps:", "unused(resource.TestCase{Steps:", 1), false},
		{"conditional call", strings.Replace(valid, "resource.Test(t,", "if false { return }; resource.Test(t,", 1), false},
		{"removed before plan", fixture(apply + `{Config:testAccConfig(map[string]interface{}{"ids":REMOVEKEY})},` + plan + reorderedApply), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			issues := Check("test.go", []byte(tc.src), nil, "alicloud_example", r)
			if (len(issues) == 0) != tc.pass {
				t.Fatalf("issues=%v want pass=%v", issues, tc.pass)
			}
		})
	}
}

func TestCandidatesAndNestedList(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{
		"outer": {Type: schema.TypeSet, Optional: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"ids": {Type: schema.TypeList, Required: true},
		}}},
		"plain_set": {Type: schema.TypeSet, Optional: true},
		"single":    {Type: schema.TypeList, Optional: true, MaxItems: 1},
		"output":    {Type: schema.TypeList, Computed: true},
		//"removed":         {Type: schema.TypeList, Optional: true, Removed: "gone"},
		"computed_parent": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"child": {Type: schema.TypeList, Optional: true}}}},
	}}
	apply := `{Config:testAccConfig(map[string]interface{}{"outer":[]map[string]interface{}{{"ids":[]int{1,2},"name":"same"}}})},`
	plan := `{Config:testAccConfig(map[string]interface{}{"outer":[]map[string]interface{}{{"ids":[]int{2,1},"name":"same"}}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	final := `{Config:testAccConfig(map[string]interface{}{"outer":[]map[string]interface{}{{"ids":[]int{2,1},"name":"same"}}})},`
	if issues := Check("test.go", []byte(fixture(apply)), nil, "alicloud_example", r); len(issues) != 1 || !strings.Contains(issues[0], "outer.ids") {
		t.Fatalf("nested List not required: %v", issues)
	}
	if issues := Check("test.go", []byte(fixture(apply+plan+final)), nil, "alicloud_example", r); len(issues) > 0 {
		t.Fatal(issues)
	}
	r.Schema["outer"].Type = schema.TypeList
	r.Schema["outer"].MaxItems = 1
	if issues := Check("test.go", []byte(fixture(apply+plan+final)), nil, "alicloud_example", r); len(issues) > 0 {
		t.Fatal(issues)
	}
	delete(r.Schema, "outer")
	if issues := Check("test.go", []byte(fixture("")), nil, "alicloud_example", r); len(issues) > 0 {
		t.Fatalf("Set itself must not require reorder coverage: %v", issues)
	}
}

func TestNoHiddenBuilderMutation(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true}}}
	apply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"},"name":"old"})},`
	plan := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	final := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}})},`
	valid := fixture(apply + plan + final)
	cases := []string{
		fixture(strings.Replace(apply, `"ids":`, `"lifecycle":[]map[string]interface{}{{"ignore_changes":[]string{"ids"}}},"ids":`, 1) + plan + final),
		strings.Replace(valid, "resource.Test(t,", `t.FailNow(); resource.Test(t,`, 1),
		strings.Replace(valid, "resource.TestCase{Steps:", `resource.TestCase{PreCheck:func(){t.Skip("disabled")},Steps:`, 1),
		strings.Replace(valid, "resource.Test(t,", `alias:=testAccConfig; alias(map[string]interface{}{"name":"new"}); resource.Test(t,`, 1),
		strings.Replace(valid, "resource.Test(t,", `testAccConfig(map[string]interface{}{"name":"new"}); resource.Test(t,`, 1),
		fixture(strings.Replace(apply, "{Config:", `{Check:makeCheck(testAccConfig(map[string]interface{}{"name":"new"})),Config:`, 1) + plan + final),
		fixture(apply + strings.Replace(plan, "PlanOnly:true", "PlanOnly:true,SkipFunc:skip", 1) + final),
		fixture(strings.Replace(apply, "{Config:", "{RefreshState:true,Config:", 1) + plan + final),
		strings.Replace(valid, "resource.Test(t,", `t.Skip("disabled"); resource.Test(t,`, 1),
		strings.Replace(valid, "resource.Test(t,", `return; resource.Test(t,`, 1),
	}
	for i, src := range cases {
		if issues := Check("test.go", []byte(src), nil, "alicloud_example", r); len(issues) == 0 {
			t.Errorf("case %d incorrectly passed", i)
		}
	}
}

func TestConstructorAliasesShareCoverage(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true}}}
	src := fixture(`{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"}})},
 {Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}}),PlanOnly:true,ExpectNonEmptyPlan:true},
 {Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}})},`)
	if issues := Check("test.go", []byte(src), nil, "alicloud_old_example", r, "alicloud_example"); len(issues) > 0 {
		t.Fatal(issues)
	}
}

func exampleOrderSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return old == new
}

func suppressedFixture(steps, unit string) string {
	return `package alicloud
import (
 "testing"
 "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)
` + unit + `
func TestAccExample(t *testing.T) {
 resourceId := "alicloud_example.default"
 testAccConfig := resourceTestAccConfigFunc(resourceId, "test", dependence)
 resource.Test(t, resource.TestCase{Steps: []resource.TestStep{` + steps + `}})
}`
}

func exampleOrderUnitNet() string {
	return `func TestUnitExampleOrderSuppress(t *testing.T) {
 testCases := []struct {
  name     string
  old      string
  new      string
  expected bool
 }{
  {"ReorderOnly", "a,b", "b,a", true},
  {"MemberAdded", "a,b", "a,b,c", false},
 }
 for _, tc := range testCases {
  if got := exampleOrderSuppressFunc("ids", tc.old, tc.new, nil); got != tc.expected {
   t.Fatalf("case %s: got %v", tc.name, got)
  }
 }
}`
}

func TestSuppressedListCoverage(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true, DiffSuppressFunc: exampleOrderSuppressFunc}}}
	apply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"},"name":"same"})},`
	reorderedApply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"},"name":"same"})},`
	reorderedPlan := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	unit := exampleOrderUnitNet()
	cases := []struct {
		name, src string
		pass      bool
	}{
		{"suppressed reorder with bounded net", suppressedFixture(apply+reorderedApply, unit), true},
		{"no reordered apply", suppressedFixture(apply, unit), false},
		{"reordered plan only expects clean plan", suppressedFixture(apply+strings.Replace(reorderedApply, "{Config:", "{PlanOnly:true,Config:", 1), unit), true},
		{"reordered plan only without prior apply", suppressedFixture(strings.Replace(reorderedApply, "{Config:", "{PlanOnly:true,Config:", 1), unit), false},
		{"reordered apply nonconvergent", suppressedFixture(apply+strings.Replace(reorderedApply, "{Config:", "{ExpectNonEmptyPlan:true,Config:", 1), unit), false},
		{"reordered apply changed membership", suppressedFixture(apply+strings.Replace(reorderedApply, `"b","a"`, `"b","c"`, 1), unit), false},
		{"reordered apply changed sibling", suppressedFixture(apply+strings.Replace(reorderedApply, `"ids":`, `"name":"changed","ids":`, 1), unit), false},
		{"plan-only diff step between applies", suppressedFixture(apply+reorderedPlan+reorderedApply, unit), true},
		{"missing net", suppressedFixture(apply+reorderedApply, ""), false},
		{"net without reorder case", suppressedFixture(apply+reorderedApply, strings.Replace(unit, `{"ReorderOnly", "a,b", "b,a", true},`, `{"SameOrder", "a,b", "a,b", true},`, 1)), false},
		{"net without membership case", suppressedFixture(apply+reorderedApply, strings.Replace(unit, `{"MemberAdded", "a,b", "a,b,c", false},`, `{"SameOrder", "a,b", "a,b", false},`, 1)), false},
		{"net reorder case unsuppressed", suppressedFixture(apply+reorderedApply, strings.Replace(unit, `"b,a", true},`, `"b,a", false},`, 1)), false},
		{"net membership case suppressed", suppressedFixture(apply+reorderedApply, strings.Replace(unit, `"a,b,c", false},`, `"a,b,c", true},`, 1)), false},
		{"bypass marker", "// collection-order: bypass ids -- pre-existing attribute outside this change.\n" + suppressedFixture(apply, ""), true},
		{"bypass marker without reason", "// collection-order: bypass ids\n" + suppressedFixture(apply, ""), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			issues := Check("test.go", []byte(tc.src), nil, "alicloud_example", r)
			if (len(issues) == 0) != tc.pass {
				t.Fatalf("issues=%v want pass=%v", issues, tc.pass)
			}
		})
	}
}

func TestSuppressedCandidateNames(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{
		"ids":  {Type: schema.TypeList, Optional: true, DiffSuppressFunc: exampleOrderSuppressFunc},
		"anon": {Type: schema.TypeList, Optional: true, DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return true }},
		"set":  {Type: schema.TypeSet, Optional: true, DiffSuppressFunc: exampleOrderSuppressFunc},
	}}
	suppressed := Suppressed(r)
	if len(suppressed) != 2 || suppressed["ids"] != "exampleOrderSuppressFunc" || suppressed["anon"] != "" {
		t.Fatalf("unexpected suppressed map: %v", suppressed)
	}
}

func TestAnonymousSuppressorRejected(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true, DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return true }}}}
	issues := Check("test.go", []byte(fixture(`{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"}})},`)), nil, "alicloud_example", r)
	if len(issues) != 1 || !strings.Contains(issues[0], "anonymous DiffSuppressFunc") {
		t.Fatalf("issues=%v", issues)
	}
}

func TestMixedSuppressedAndPlainCoverage(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{
		"ids":   {Type: schema.TypeList, Optional: true, DiffSuppressFunc: exampleOrderSuppressFunc},
		"plain": {Type: schema.TypeList, Optional: true},
	}}
	step1 := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"},"plain":[]string{"x","y"},"name":"same"})},`
	step2 := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}})},`
	step3 := `{Config:testAccConfig(map[string]interface{}{"plain":[]string{"y","x"}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	step4 := `{Config:testAccConfig(map[string]interface{}{"plain":[]string{"y","x"}})},`
	unit := exampleOrderUnitNet()
	if issues := Check("test.go", []byte(suppressedFixture(step1+step2+step3+step4, unit)), nil, "alicloud_example", r); len(issues) > 0 {
		t.Fatal(issues)
	}
	issues := Check("test.go", []byte(suppressedFixture(step1+step2, unit)), nil, "alicloud_example", r)
	if len(issues) != 1 || !strings.Contains(issues[0], "plain") {
		t.Fatalf("issues=%v", issues)
	}
}

// qualifiedFixture mirrors fixture() but uses the split-package qualified
// constructor acctest.ResourceTestAccConfig, the form the alicloud/service/*
// packages use after the relocation. The checker must recognize it as readily
// as the bare legacy resourceTestAccConfigFunc; the package alias is a per-file
// choice and is not constrained.
func qualifiedFixture(steps string) string {
	return `package alicloud_test
import (
 "testing"
 acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
 "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)
func TestAccExample(t *testing.T) {
 resourceId := "alicloud_example.default"
 testAccConfig := acctest.ResourceTestAccConfig(resourceId, "test", dependence)
 resource.Test(t, resource.TestCase{Steps: []resource.TestStep{` + steps + `}})
}`
}

// unitNetFile wraps a unit-net function as a complete sibling *_unit_test.go
// source in the same package, the split-package layout that separates the
// suppressor's table test from the TestAcc steps.
func unitNetFile(unit string) string {
	return `package alicloud_test
import (
 "testing"
)
` + unit
}

func TestQualifiedConstructorReorderCoverage(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true}}}
	apply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"},"name":"same"})},`
	plan := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	reorderedApply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}})},`
	src := qualifiedFixture(apply + plan + reorderedApply)
	if issues := Check("test.go", []byte(src), nil, "alicloud_example", r); len(issues) != 0 {
		t.Fatalf("qualified constructor not recognized as a config builder: %v", issues)
	}
	// A bare legacy constructor in the same shape still passes (regression guard).
	if issues := Check("test.go", []byte(fixture(apply+plan+reorderedApply)), nil, "alicloud_example", r); len(issues) != 0 {
		t.Fatalf("bare constructor regressed: %v", issues)
	}
}

func TestSuppressedUnitNetInSiblingFile(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true, DiffSuppressFunc: exampleOrderSuppressFunc}}}
	apply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"},"name":"same"})},`
	reorderedApply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"},"name":"same"})},`
	// Main *_test.go file: TestAcc steps with the qualified constructor, no inline unit net.
	mainSrc := qualifiedFixture(apply + reorderedApply)
	// Sibling *_unit_test.go file: carries the suppressor unit net.
	siblingSrc := unitNetFile(exampleOrderUnitNet())
	testSources := map[string][]byte{"example_unit_test.go": []byte(siblingSrc)}
	if issues := Check("test.go", []byte(mainSrc), testSources, "alicloud_example", r); len(issues) != 0 {
		t.Fatalf("expected pass with unit net in sibling file, got %v", issues)
	}
	// Without the sibling file the net is missing and coverage fails.
	if issues := Check("test.go", []byte(mainSrc), nil, "alicloud_example", r); len(issues) == 0 {
		t.Fatal("expected failure when unit net sibling file is absent")
	}
	// A net lacking the reorder case still fails even with the sibling file.
	badNet := strings.Replace(exampleOrderUnitNet(), `{"ReorderOnly", "a,b", "b,a", true},`, `{"SameOrder", "a,b", "a,b", true},`, 1)
	badSources := map[string][]byte{"example_unit_test.go": []byte(unitNetFile(badNet))}
	if issues := Check("test.go", []byte(mainSrc), badSources, "alicloud_example", r); len(issues) == 0 {
		t.Fatal("expected failure when sibling net lacks the reorder case")
	}
}

// dataSourceFixture mirrors fixture() but uses the three-part data-source
// address "data.<type>.<name>" that a data source's TestAcc carries. The
// convergence logic keys by builder variable, so a three-part address flows
// through the same reorder check a two-part resource address does.
func dataSourceFixture(steps string) string {
	return `package alicloud
import (
 "testing"
 "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)
func TestAccExample(t *testing.T) {
 resourceId := "data.alicloud_example.default"
 testAccConfig := resourceTestAccConfigFunc(resourceId, "test", dependence)
 resource.Test(t, resource.TestCase{Steps: []resource.TestStep{` + steps + `}})
}`
}

// TestUnitDataSourceReadOnlyCheck covers the read-only equivalent check target
// for a data source (CheckScope with DataSourceScope): a read-only schema (no
// Create/Update/Delete) passes with or without an explicit reorder test and
// even with no test file, because Read is order-insensitive; a data source
// carrying any write lifecycle is anomalous and fails. The TestUnit prefix makes
// ci-check-quick's changed-unit-tests stage execute these regressions.
func TestUnitDataSourceReadOnlyCheck(t *testing.T) {
	readOnly := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true}}}
	apply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"},"name":"same"})},`
	plan := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	reorderedApply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}})},`

	// No reorder test: passes under read-only semantics.
	if issues := CheckScope("test.go", []byte(dataSourceFixture(apply)), nil, "alicloud_example", readOnly, DataSourceScope); len(issues) != 0 {
		t.Fatalf("read-only data source without reorder test failed: %v", issues)
	}
	// A convergent reorder test is still recognised; the check passes either way.
	if issues := CheckScope("test.go", []byte(dataSourceFixture(apply+plan+reorderedApply)), nil, "alicloud_example", readOnly, DataSourceScope); len(issues) != 0 {
		t.Fatalf("read-only data source with reorder test failed: %v", issues)
	}
	// No test file at all: the read-only lifecycle alone is the check target.
	if issues := CheckScope("test.go", []byte(`package alicloud`), nil, "alicloud_example", readOnly, DataSourceScope); len(issues) != 0 {
		t.Fatalf("read-only data source with no test failed: %v", issues)
	}
	// A data source carrying a write lifecycle is anomalous: the read-only
	// waiver does not apply and the check fails for each CRUD callback.
	writeCases := []struct {
		name string
		r    *schema.Resource
	}{
		{"Create", &schema.Resource{Schema: readOnly.Schema, Create: func(*schema.ResourceData, interface{}) error { return nil }}},
		{"Update", &schema.Resource{Schema: readOnly.Schema, Update: func(*schema.ResourceData, interface{}) error { return nil }}},
		{"Delete", &schema.Resource{Schema: readOnly.Schema, Delete: func(*schema.ResourceData, interface{}) error { return nil }}},
	}
	for _, tc := range writeCases {
		t.Run(tc.name, func(t *testing.T) {
			issues := CheckScope("test.go", []byte(dataSourceFixture(apply)), nil, "alicloud_example", tc.r, DataSourceScope)
			if len(issues) != 1 || !strings.Contains(issues[0], "write lifecycle") {
				t.Fatalf("write lifecycle should fail: %v", issues)
			}
		})
	}
}

// TestUnitDataSourceAddress covers the three-part "data.<type>.<name>" address
// matching in inspectTest. Under the resource scope (read/write), coverage is
// enforced, so a convergent reorder over a data-source-style builder passes
// only when the three-part address is recognised; a wrong type segment is not
// captured and coverage is not established. The TestUnit prefix makes
// ci-check-quick's changed-unit-tests stage execute this regression.
func TestUnitDataSourceAddress(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true}}}
	apply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"},"name":"same"})},`
	plan := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	reorderedApply := `{Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}})},`
	src := dataSourceFixture(apply + plan + reorderedApply)
	// The three-part data.<type>.<name> address is recognised the same way a
	// two-part resource address is, so a convergent reorder over a data-source
	// builder establishes coverage under the resource scope.
	if issues := Check("test.go", []byte(src), nil, "alicloud_example", r); len(issues) != 0 {
		t.Fatalf("three-part data-source address not recognised: %v", issues)
	}
	// The type segment must match the entity under check; a different type's
	// three-part builder is not captured and coverage is not established.
	wrong := strings.Replace(src, "data.alicloud_example.default", "data.alicloud_other.default", 1)
	if issues := Check("test.go", []byte(wrong), nil, "alicloud_example", r); len(issues) == 0 {
		t.Fatal("wrong-type three-part address should not establish coverage")
	}
}
