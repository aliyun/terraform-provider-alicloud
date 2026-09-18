package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Agentloop EndpointConnector. >>> Resource test cases, automatically generated.
// Case EndpointConnector全生命周期 11080
func TestAccAliCloudAgentloopEndpointConnector_basic11080(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_agentloop_endpoint_connector.default"
	ra := resourceAttrInit(resourceId, AlicloudAgentloopEndpointConnectorMap11080)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AgentloopServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAgentloopEndpointConnector")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccconnector%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopEndpointConnectorBasicDependence11080)
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
					"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
					"type":        "model_service",
					"name":        name,
					"endpoint":    "https://dashscope.aliyuncs.com/compatible-mode/v1",
					"credential": map[string]interface{}{
						"api_key": "sk-test-1234567890",
					},
					"alias":       name + "alias",
					"description": "terraform自动化测试描述",
					"headers": []map[string]interface{}{
						{
							"key":   "X-Header-1",
							"value": "value1",
						},
					},
					"properties": map[string]interface{}{
						"protocol":  "openai-compatible",
						"property1": "value1",
						"modelList": "[{\\\"model\\\":\\\"qwen-plus\\\"}]",
					},
					"tags": []string{"tag1", "tag2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space": name + "-as",
						"type":        "model_service",
						"name":        name,
						"endpoint":    "https://dashscope.aliyuncs.com/compatible-mode/v1",
						"alias":       name + "alias",
						"description": "terraform自动化测试描述",
						"headers.#":   "1",
						"tags.#":      "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     name + "update",
					"endpoint": "https://dashscope.aliyuncs.com/compatible-mode/v2",
					"credential": map[string]interface{}{
						"api_key": "sk-test-0987654321",
					},
					"alias":       name + "aliasupdate",
					"description": "terraform自动化测试描述更新",
					"headers": []map[string]interface{}{
						{
							"key":   "X-Header-2",
							"value": "value2",
						},
					},
					"properties": map[string]interface{}{
						"protocol":  "openai-compatible",
						"property1": "value2",
						"modelList": "[{\\\"model\\\":\\\"qwen-max\\\"}]",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name + "update",
						"endpoint":    "https://dashscope.aliyuncs.com/compatible-mode/v2",
						"alias":       name + "aliasupdate",
						"description": "terraform自动化测试描述更新",
						"headers.#":   "1",
					}),
				),
			},
			{
				// TypeList order coverage baseline for headers: two distinct
				// entries (headers is updated in place).
				Config: testAccConfig(map[string]interface{}{
					"headers": []map[string]interface{}{
						{
							"key":   "X-Header-2",
							"value": "value2",
						},
						{
							"key":   "X-Header-1",
							"value": "value1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"headers.#": "2",
					}),
				),
			},
			{
				// TypeList order coverage: a reorder-only plan against the
				// previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"headers": []map[string]interface{}{
						{
							"key":   "X-Header-1",
							"value": "value1",
						},
						{
							"key":   "X-Header-2",
							"value": "value2",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered headers must then converge with an
				// empty post-apply plan. The Get API never returns headers, so
				// convergence is guaranteed at the state level.
				Config: testAccConfig(map[string]interface{}{
					"headers": []map[string]interface{}{
						{
							"key":   "X-Header-1",
							"value": "value1",
						},
						{
							"key":   "X-Header-2",
							"value": "value2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"headers.#": "2",
					}),
				),
			},
			{
				// TypeList order coverage for tags (ForceNew): a reorder-only
				// plan against the previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"tags": []string{"tag2", "tag1"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered tags recreates the connector and must
				// converge with an empty post-apply plan.
				Config: testAccConfig(map[string]interface{}{
					"tags": []string{"tag2", "tag1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// credential is desensitized; properties gets server-side default keys
				// injected; headers and tags are never returned by the Get API.
				ImportStateVerifyIgnore: []string{"credential", "properties", "headers", "tags"},
			},
		},
	})
}

var AlicloudAgentloopEndpointConnectorMap11080 = map[string]string{
	"connector_id": CHECKSET,
	"region_id":    CHECKSET,
	"created_at":   CHECKSET,
	"updated_at":   CHECKSET,
}

func AlicloudAgentloopEndpointConnectorBasicDependence11080(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
    agent_space = "${var.name}-as"
}

`, name)
}

// Test Agentloop EndpointConnector. <<< Resource test cases, automatically generated.
