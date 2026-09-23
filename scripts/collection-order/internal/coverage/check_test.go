package coverage

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func fixture(steps string) string {
	return `package alicloud
import (
 "testing"
 "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		{"explicit convergent apply", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", "{ExpectNonEmptyPlan:false,Config:", 1)), true},
		{"apply retained builder snapshot", fixture(apply + plan + `{Config:testAccConfig(map[string]interface{}{})},`), true},
		{"empty plan", fixture(apply + strings.Replace(plan, "ExpectNonEmptyPlan:true", "ExpectNonEmptyPlan:false", 1) + reorderedApply), false},
		{"default empty plan", fixture(apply + strings.Replace(plan, ",ExpectNonEmptyPlan:true", "", 1) + reorderedApply), false},
		{"diff without apply", fixture(apply + plan), false},
		{"nonconvergent apply", fixture(apply + plan + strings.Replace(reorderedApply, "{Config:", "{ExpectNonEmptyPlan:true,Config:", 1)), false},
		{"apply original order", fixture(apply + plan + apply), false},
		{"apply changed membership", fixture(apply + plan + strings.Replace(reorderedApply, `"b","a"`, `"b","c"`, 1)), false},
		{"apply changed sibling", fixture(apply + plan + strings.Replace(reorderedApply, `"ids":`, `"name":"changed","ids":`, 1)), false},
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
		{"unchanged", fixture(apply + strings.Replace(plan, `"b","a"`, `"a","b"`, 1) + reorderedApply), false},
		{"plan changed membership", fixture(apply + strings.Replace(plan, `"b","a"`, `"b","c"`, 1) + reorderedApply), false},
		{"duplicates only", strings.ReplaceAll(valid, `"b"`, `"a"`), false},
		{"apply without plan", fixture(apply + reorderedApply), false},
		{"comment only", fixture(apply + "/*" + plan + reorderedApply + "*/"), false},
		{"go build header", "//go:build ignore\n\n" + valid, false},
		{"legacy build header", "// +build ignore\n\n" + valid, false},
		{"go build mention after package", strings.Replace(valid, "func TestAcc", "//go:build example\nfunc TestAcc", 1), true},
		{"legacy build mention after package", strings.Replace(valid, "func TestAcc", "// +build example\nfunc TestAcc", 1), true},
		{"ordinary header comment", "// go:build example\n\n" + valid, true},
		{"plan changed sibling", fixture(apply + strings.Replace(plan, `"ids":`, `"name":"changed","ids":`, 1) + reorderedApply), false},
		{"other resource", strings.Replace(valid, "alicloud_example.default", "alicloud_other.default", 1), false},
		{"not TestAcc", strings.Replace(valid, "TestAccExample", "TestUnitExample", 1), false},
		{"unused steps", strings.Replace(valid, "resource.Test(t, resource.TestCase{Steps:", "unused(resource.TestCase{Steps:", 1), false},
		{"conditional call", strings.Replace(valid, "resource.Test(t,", "if false { return }; resource.Test(t,", 1), false},
		{"removed before plan", fixture(apply + `{Config:testAccConfig(map[string]interface{}{"ids":REMOVEKEY})},` + plan + reorderedApply), false},
		{"intervening plan", fixture(apply + plan + plan + reorderedApply), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			issues := Check("test.go", []byte(tc.src), "alicloud_example", r)
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
		"plain_set":       {Type: schema.TypeSet, Optional: true},
		"single":          {Type: schema.TypeList, Optional: true, MaxItems: 1},
		"output":          {Type: schema.TypeList, Computed: true},
		"removed":         {Type: schema.TypeList, Optional: true, Removed: "gone"},
		"computed_parent": {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: map[string]*schema.Schema{"child": {Type: schema.TypeList, Optional: true}}}},
	}}
	apply := `{Config:testAccConfig(map[string]interface{}{"outer":[]map[string]interface{}{{"ids":[]int{1,2},"name":"same"}}})},`
	plan := `{Config:testAccConfig(map[string]interface{}{"outer":[]map[string]interface{}{{"ids":[]int{2,1},"name":"same"}}}),PlanOnly:true,ExpectNonEmptyPlan:true},`
	final := `{Config:testAccConfig(map[string]interface{}{"outer":[]map[string]interface{}{{"ids":[]int{2,1},"name":"same"}}})},`
	if issues := Check("test.go", []byte(fixture(apply)), "alicloud_example", r); len(issues) != 1 || !strings.Contains(issues[0], "outer.ids") {
		t.Fatalf("nested List not required: %v", issues)
	}
	if issues := Check("test.go", []byte(fixture(apply+plan+final)), "alicloud_example", r); len(issues) > 0 {
		t.Fatal(issues)
	}
	r.Schema["outer"].Type = schema.TypeList
	r.Schema["outer"].MaxItems = 1
	if issues := Check("test.go", []byte(fixture(apply+plan+final)), "alicloud_example", r); len(issues) > 0 {
		t.Fatal(issues)
	}
	delete(r.Schema, "outer")
	if issues := Check("test.go", []byte(fixture("")), "alicloud_example", r); len(issues) > 0 {
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
		if issues := Check("test.go", []byte(src), "alicloud_example", r); len(issues) == 0 {
			t.Errorf("case %d incorrectly passed", i)
		}
	}
}

func TestConstructorAliasesShareCoverage(t *testing.T) {
	r := &schema.Resource{Schema: map[string]*schema.Schema{"ids": {Type: schema.TypeList, Optional: true}}}
	src := fixture(`{Config:testAccConfig(map[string]interface{}{"ids":[]string{"a","b"}})},
 {Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}}),PlanOnly:true,ExpectNonEmptyPlan:true},
 {Config:testAccConfig(map[string]interface{}{"ids":[]string{"b","a"}})},`)
	if issues := Check("test.go", []byte(src), "alicloud_old_example", r, "alicloud_example"); len(issues) > 0 {
		t.Fatal(issues)
	}
}
