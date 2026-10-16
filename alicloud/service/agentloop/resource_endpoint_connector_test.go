package agentloop_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Agentloop EndpointConnector. >>> Resource test cases, automatically generated.
// Case EndpointConnector全生命周期 11080
func TestAccAliCloudAgentloopEndpointConnector_basic11080(t *testing.T) {
	resourceId := "alicloud_agentloop_endpoint_connector.default"
	rac := resourceAttrCheckInit(resourceId, AlicloudAgentloopEndpointConnectorMap11080, (*agentloop.Service).DescribeEndpointConnector)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccconnector%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopEndpointConnectorBasicDependence11080)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
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
						"qpsLimit":  "50",
						"modelList": "[{\\\"model\\\":\\\"qwen-plus\\\"}]",
					},
					"tags": []string{"tag1", "tag2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":          name + "-as",
						"type":                 "model_service",
						"name":                 name,
						"endpoint":             "https://dashscope.aliyuncs.com/compatible-mode/v1",
						"alias":                name + "alias",
						"description":          "terraform自动化测试描述",
						"headers.#":            "1",
						"headers.0.key":        "X-Header-1",
						"headers.0.value":      "value1",
						"properties.%":         "3",
						"properties.protocol":  "openai-compatible",
						"properties.qpsLimit":  "50",
						"properties.modelList": "[{\"model\":\"qwen-plus\"}]",
						"tags.#":               "2",
						"tags.0":               "tag1",
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
						"qpsLimit":  "80",
						"modelList": "[{\\\"model\\\":\\\"qwen-max\\\"}]",
					},
					"tags": []string{"tag2", "tag3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name + "update",
						"endpoint":             "https://dashscope.aliyuncs.com/compatible-mode/v2",
						"alias":                name + "aliasupdate",
						"description":          "terraform自动化测试描述更新",
						"headers.#":            "1",
						"headers.0.key":        "X-Header-2",
						"headers.0.value":      "value2",
						"properties.%":         "3",
						"properties.qpsLimit":  "80",
						"properties.modelList": "[{\"model\":\"qwen-max\"}]",
						"tags.#":               "2",
						"tags.0":               "tag2",
						"tags.1":               "tag3",
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
						"headers.#":       "2",
						"headers.0.key":   "X-Header-2",
						"headers.0.value": "value2",
						"headers.1.key":   "X-Header-1",
						"headers.1.value": "value1",
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
				// empty post-apply plan: the Get API returns headers in the
				// submitted order.
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
						"headers.#":       "2",
						"headers.0.key":   "X-Header-1",
						"headers.0.value": "value1",
						"headers.1.key":   "X-Header-2",
						"headers.1.value": "value2",
					}),
				),
			},
			{
				// TypeList order coverage for tags: a reorder-only plan against
				// the previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"tags": []string{"tag3", "tag2"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered tags updates the connector in place and
				// must then converge with an empty post-apply plan.
				Config: testAccConfig(map[string]interface{}{
					"tags": []string{"tag3", "tag2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.#": "2",
						"tags.0": "tag3",
						"tags.1": "tag2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// credential is desensitized by the Get API and cannot be
				// imported; headers, properties and tags are read back.
				ImportStateVerifyIgnore: []string{"credential"},
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
