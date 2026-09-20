package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Agentloop EvaluatorSkill. >>> Resource test cases, automatically generated.
// Case EvaluatorSkill全生命周期 11083
func TestAccAliCloudAgentloopEvaluatorSkill_basic11083(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_agentloop_evaluator_skill.default"
	ra := resourceAttrInit(resourceId, AlicloudAgentloopEvaluatorSkillMap11083)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AgentloopServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAgentloopEvaluatorSkill")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccevalskill%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopEvaluatorSkillBasicDependence11083)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":  "${alicloud_agentloop_agent_space.default.agent_space}",
					"name":         "${alicloud_agentloop_evaluator.default.name}",
					"skill_name":   name,
					"description":  "terraform自动化测试描述",
					"display_name": name,
					"enable":       "true",
					"files": []map[string]interface{}{
						{
							"name": "SKILL.md",
							// The backend validates that the SKILL.md frontmatter
							// 'name' equals skill_name, so inject the random name.
							"content": fmt.Sprintf("---\\nname: %s\\ndescription: terraform acc test skill v1\\n---\\n# Test Skill V1", name),
							"remark":  "初始版本",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":  name + "-as",
						"name":         name + "eval",
						"skill_name":   name,
						"description":  "terraform自动化测试描述",
						"display_name": name,
						"enable":       "true",
						"files.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					// description/display_name/enable are ForceNew: they must keep
					// the same values as step 0, otherwise the skill is recreated
					// and the cumulative checks for them (still expecting the step
					// 0 values) fail on the blank replacement.
					"description":  "terraform自动化测试描述",
					"display_name": name,
					"enable":       "true",
					"files": []map[string]interface{}{
						{
							"name":    "SKILL.md",
							"content": fmt.Sprintf("---\\nname: %s\\ndescription: terraform acc test skill v2\\n---\\n# Test Skill V2", name),
							"remark":  "更新版本",
						},
						{
							// A second file under references/ (the backend only
							// allows SKILL.md, references/ and scripts/ entries),
							// so files.name takes a distinct value in this step.
							"name":    "references/detail.md",
							"content": "# Detail Reference V2",
							"remark":  "新增引用",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"files.#": "2",
					}),
				),
			},
			{
				// TypeList order coverage: a reorder-only plan against the
				// previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"files": []map[string]interface{}{
						{
							"name":    "references/detail.md",
							"content": "# Detail Reference V2",
							"remark":  "新增引用",
						},
						{
							"name":    "SKILL.md",
							"content": fmt.Sprintf("---\\nname: %s\\ndescription: terraform acc test skill v2\\n---\\n# Test Skill V2", name),
							"remark":  "更新版本",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered files must then converge with an empty
				// post-apply plan, proving the API preserves the order.
				Config: testAccConfig(map[string]interface{}{
					"files": []map[string]interface{}{
						{
							"name":    "references/detail.md",
							"content": "# Detail Reference V2",
							"remark":  "新增引用",
						},
						{
							"name":    "SKILL.md",
							"content": fmt.Sprintf("---\\nname: %s\\ndescription: terraform acc test skill v2\\n---\\n# Test Skill V2", name),
							"remark":  "更新版本",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"files.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "terraform自动化测试描述",
					"display_name": name,
					"enable":       "true",
					"files": []map[string]interface{}{
						{
							// Configuration order deliberately reversed against the
							// name-sorted order (references/detail.md before
							// SKILL.md). The default post-step plan check fails if
							// the Read produces a perpetual diff by re-sorting.
							"name":    "references/detail.md",
							"content": "# Detail Reference V3",
							"remark":  "引用更新",
						},
						{
							"name":    "SKILL.md",
							"content": fmt.Sprintf("---\\nname: %s\\ndescription: terraform acc test skill v3\\n---\\n# Test Skill V3", name),
							"remark":  "再次更新",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"files.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// GetEvaluatorSkill does not return files[].remark (the Read
				// preserves them from the configuration, which is not available
				// on import) and does not guarantee the file order, so the files
				// attribute cannot be verified on import.
				ImportStateVerifyIgnore: []string{"files"},
			},
		},
	})
}

var AlicloudAgentloopEvaluatorSkillMap11083 = map[string]string{
	"created_at": CHECKSET,
	"updated_at": CHECKSET,
}

func AlicloudAgentloopEvaluatorSkillBasicDependence11083(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
    agent_space = "${var.name}-as"
}

resource "alicloud_agentloop_evaluator" "default" {
    agent_space = alicloud_agentloop_agent_space.default.agent_space
    name        = "${var.name}eval"
    type        = "AGENT"
    metric_name = "tf_test_metric"
    version     = "v1"
}

`, name)
}

// Test Agentloop EvaluatorSkill. <<< Resource test cases, automatically generated.
