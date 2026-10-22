package agentloop_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test Agentloop ContextStore. >>> Resource test cases, automatically generated.
// Case ContextStore全生命周期 11077
func TestAccAliCloudAgentloopContextStore_basic11077(t *testing.T) {
	resourceId := "alicloud_agentloop_context_store.default"
	rac := resourceAttrCheckInit(resourceId, AlicloudAgentloopContextStoreMap11077, (*agentloop.Service).DescribeContextStore)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccctxstore%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopContextStoreBasicDependence11077)
	// The ID is reused after a replacement, so create_time tells an in-place
	// update apart from a recreation.
	var createTime string
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
					"agent_space":        "${alicloud_agentloop_agent_space.default.agent_space}",
					"context_store_name": name,
					"context_type":       "experience",
					"description":        "terraform自动化测试描述",
					"config": []map[string]interface{}{
						{
							"metadata_field": map[string]interface{}{
								"userId":    "user_id",
								"sessionId": "session_id",
							},
							"service_names": []string{"svc-tf-acc"},
							"source": []map[string]interface{}{
								{
									"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
									"start_time":  "2026-01-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":                   name + "-as",
						"context_store_name":            name,
						"context_type":                  "experience",
						"description":                   "terraform自动化测试描述",
						"config.#":                      "1",
						"config.0.service_names.#":      "1",
						"config.0.service_names.0":      "svc-tf-acc",
						"config.0.source.#":             "1",
						"config.0.source.0.agent_space": name + "-as",
						"config.0.source.0.start_time":  "2026-01-01T00:00:00Z",
						"config.0.metadata_field.%":     "2",
					}),
					checkContextStoreCreateTime(resourceId, &createTime, false),
				),
			},
			{
				// Update-only step: only description changes (config keeps the
				// values merged from the previous step).
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform自动化测试描述仅更新",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform自动化测试描述仅更新",
						"config.#":    "1",
					}),
					checkContextStoreCreateTime(resourceId, &createTime, false),
				),
			},
			{
				// In-place config update: service_names and metadata_field are
				// replaced as a whole while source is unchanged, so the store
				// keeps its create_time.
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform自动化测试描述更新",
					"config": []map[string]interface{}{
						{
							"metadata_field": map[string]interface{}{
								"userId": "user_id_updated",
							},
							"service_names": []string{"svc-tf-acc-updated"},
							"source": []map[string]interface{}{
								{
									"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
									"start_time":  "2026-01-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"context_type":                   "experience",
						"description":                    "terraform自动化测试描述更新",
						"config.#":                       "1",
						"config.0.service_names.#":       "1",
						"config.0.service_names.0":       "svc-tf-acc-updated",
						"config.0.source.0.start_time":   "2026-01-01T00:00:00Z",
						"config.0.metadata_field.%":      "1",
						"config.0.metadata_field.userId": "user_id_updated",
					}),
					checkContextStoreMetadataUpdate(resourceId),
					checkContextStoreCreateTime(resourceId, &createTime, false),
				),
			},
			{
				// source is immutable after creation, so changing start_time
				// recreates the store. Two distinct service_names also set the
				// TypeList order coverage baseline.
				Config: testAccConfig(map[string]interface{}{
					"config": []map[string]interface{}{
						{
							"metadata_field": map[string]interface{}{
								"userId":    "user_id_updated",
								"sessionId": "session_id",
							},
							"service_names": []string{"svc-tf-acc-updated", "svc-tf-acc"},
							"source": []map[string]interface{}{
								{
									"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
									"start_time":  "2026-02-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.#":                     "1",
						"config.0.service_names.#":     "2",
						"config.0.service_names.0":     "svc-tf-acc-updated",
						"config.0.service_names.1":     "svc-tf-acc",
						"config.0.source.0.start_time": "2026-02-01T00:00:00Z",
						"config.0.metadata_field.%":    "2",
					}),
					checkContextStoreCreateTime(resourceId, &createTime, true),
				),
			},
			{
				// TypeList order coverage: a reorder-only plan against the
				// previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"config": []map[string]interface{}{
						{
							"metadata_field": map[string]interface{}{
								"userId":    "user_id_updated",
								"sessionId": "session_id",
							},
							"service_names": []string{"svc-tf-acc", "svc-tf-acc-updated"},
							"source": []map[string]interface{}{
								{
									"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
									"start_time":  "2026-02-01T00:00:00Z",
								},
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered service_names updates the store in place
				// and must converge with an empty post-apply plan; serviceNames
				// are read back in API order.
				Config: testAccConfig(map[string]interface{}{
					"config": []map[string]interface{}{
						{
							"metadata_field": map[string]interface{}{
								"userId":    "user_id_updated",
								"sessionId": "session_id",
							},
							"service_names": []string{"svc-tf-acc", "svc-tf-acc-updated"},
							"source": []map[string]interface{}{
								{
									"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
									"start_time":  "2026-02-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.#":                 "1",
						"config.0.service_names.#": "2",
						"config.0.service_names.0": "svc-tf-acc",
						"config.0.service_names.1": "svc-tf-acc-updated",
					}),
					checkContextStoreCreateTime(resourceId, &createTime, false),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// service_names and source are read back; the Get API mixes the
				// metadataField mapping with internal keys, so it cannot be
				// recovered on import.
				ImportStateVerifyIgnore: []string{"config.0.metadata_field"},
			},
		},
	})
}

// 直接回读 API，避免仅检查保留在 state 中的 metadata_field 而漏掉更新请求缺失。
func checkContextStoreMetadataUpdate(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok || rs.Primary == nil {
			return fmt.Errorf("resource %s missing from state", resourceId)
		}
		object, err := testAccService().DescribeContextStore(rs.Primary.ID)
		if err != nil {
			return err
		}
		config, _ := object["config"].(map[string]interface{})
		fields, _ := config["metadataField"].(map[string]interface{})
		if fields["userId"] != "user_id_updated" {
			return fmt.Errorf("API metadataField.userId = %v, want user_id_updated", fields["userId"])
		}
		if _, ok := fields["sessionId"]; ok {
			return fmt.Errorf("API metadataField.sessionId still present after whole-map replacement")
		}
		return nil
	}
}

// checkContextStoreCreateTime records create_time on the first call and then
// checks that it changed (replaced) or stayed the same (updated in place).
func checkContextStoreCreateTime(resourceId string, createTime *string, replaced bool) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(resourceId, "create_time", func(v string) error {
		prev := *createTime
		*createTime = v
		if prev == "" {
			return nil
		}
		if replaced && v == prev {
			return fmt.Errorf("create_time %s unchanged, want the Context Store to be recreated", v)
		}
		if !replaced && v != prev {
			return fmt.Errorf("create_time changed from %s to %s, want an in-place update", prev, v)
		}
		return nil
	})
}

var AlicloudAgentloopContextStoreMap11077 = map[string]string{
	"region_id":   CHECKSET,
	"create_time": CHECKSET,
	"update_time": CHECKSET,
	"status":      CHECKSET,
}

func AlicloudAgentloopContextStoreBasicDependence11077(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
    agent_space = "${var.name}-as"
}

`, name)
}

// Test Agentloop ContextStore. <<< Resource test cases, automatically generated.
