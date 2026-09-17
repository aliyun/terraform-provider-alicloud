package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Agentloop Evaluator. >>> Resource test cases, automatically generated.
// Case Evaluator全生命周期 11082
func TestAccAliCloudAgentloopEvaluator_basic11082(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_agentloop_evaluator.default"
	ra := resourceAttrInit(resourceId, AlicloudAgentloopEvaluatorMap11082)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AgentloopServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAgentloopEvaluator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccevaluator%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopEvaluatorBasicDependence11082)
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
					"agent_space":         "${alicloud_agentloop_agent_space.default.agent_space}",
					"name":                name,
					"type":                "AGENT",
					"metric_name":         "tf_test_metric",
					"version":             "v1",
					"version_description": "初始版本描述",
					"display_name":        name,
					"description":         "terraform自动化测试描述",
					// Deliberately reverse-sorted: the Read flattens annotations in
					// whatever order the API returns, so a non-lexicographic
					// configuration order fails the post-step plan check if the
					// API reorders them.
					"annotations": []string{"anno2", "anno1"},
					"config": map[string]interface{}{
						"config1": "value1",
					},
					"properties": map[string]interface{}{
						"property1": "value1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":   name + "-as",
						"name":          name,
						"type":          "AGENT",
						"metric_name":   "tf_test_metric",
						"version":       "v1",
						"display_name":  name,
						"description":   "terraform自动化测试描述",
						"annotations.#": "2",
					}),
				),
			},
			{
				// TypeList order coverage: a reorder-only plan against the
				// previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"annotations": []string{"anno1", "anno2"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered annotations must then converge with an
				// empty post-apply plan, proving the API preserves the order.
				Config: testAccConfig(map[string]interface{}{
					"annotations": []string{"anno1", "anno2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"annotations.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version":      "v2",
					"display_name": name + "update",
					"description":  "terraform自动化测试描述更新",
					"annotations":  []string{"anno3"},
					"config": map[string]interface{}{
						"config1": "value2",
					},
					"properties": map[string]interface{}{
						"property1": "value2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version":       "v2",
						"display_name":  name + "update",
						"description":   "terraform自动化测试描述更新",
						"annotations.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// version_description and version are never returned by the Get API
				// (the response only carries currentVersion/latestVersion/versions),
				// so they cannot be verified on import.
				ImportStateVerifyIgnore: []string{"version_description", "version"},
			},
		},
	})
}

var AlicloudAgentloopEvaluatorMap11082 = map[string]string{
	"created_at":      CHECKSET,
	"updated_at":      CHECKSET,
	"current_version": CHECKSET,
}

func AlicloudAgentloopEvaluatorBasicDependence11082(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
    agent_space = "${var.name}-as"
}

`, name)
}

// Test Agentloop Evaluator. <<< Resource test cases, automatically generated.
