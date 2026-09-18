package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Agentloop ContextStore. >>> Resource test cases, automatically generated.
// Case ContextStore全生命周期 11077
func TestAccAliCloudAgentloopContextStore_basic11077(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_agentloop_context_store.default"
	ra := resourceAttrInit(resourceId, AlicloudAgentloopContextStoreMap11077)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AgentloopServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAgentloopContextStore")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccctxstore%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopContextStoreBasicDependence11077)
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
									"project":    "${alicloud_log_project.default.project_name}",
									"logstore":   "${alicloud_log_store.default.logstore_name}",
									"start_time": "2026-01-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":        name + "-as",
						"context_store_name": name,
						"context_type":       "experience",
						"description":        "terraform自动化测试描述",
						"config.#":           "1",
					}),
				),
			},
			{
				// Update-only step: only description changes (config keeps the
				// values merged from the previous step), so this exercises the
				// ContextStoreUpdate path instead of a ForceNew recreation.
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform自动化测试描述仅更新",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform自动化测试描述仅更新",
						"config.#":    "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform自动化测试描述更新",
					"config": []map[string]interface{}{
						{
							"metadata_field": map[string]interface{}{
								"userId":    "user_id_updated",
								"sessionId": "session_id",
							},
							"service_names": []string{"svc-tf-acc-updated"},
							"source": []map[string]interface{}{
								{
									"project":    "${alicloud_log_project.default2.project_name}",
									"logstore":   "${alicloud_log_store.default2.logstore_name}",
									"start_time": "2026-02-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"context_type": "experience",
						"description":  "terraform自动化测试描述更新",
						"config.#":     "1",
					}),
				),
			},
			{
				// TypeList order coverage baseline: two distinct service_names
				// (config is ForceNew, so this step recreates the store).
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
									"project":    "${alicloud_log_project.default2.project_name}",
									"logstore":   "${alicloud_log_store.default2.logstore_name}",
									"start_time": "2026-02-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.#": "1",
					}),
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
									"project":    "${alicloud_log_project.default2.project_name}",
									"logstore":   "${alicloud_log_store.default2.logstore_name}",
									"start_time": "2026-02-01T00:00:00Z",
								},
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered service_names must then converge with an
				// empty post-apply plan. The backend never reads config back, so
				// convergence is guaranteed at the state level.
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
									"project":    "${alicloud_log_project.default2.project_name}",
									"logstore":   "${alicloud_log_store.default2.logstore_name}",
									"start_time": "2026-02-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// The config block is ForceNew and never read back: the backend
				// consumes metadataField, silently ignores serviceNames/source on
				// update, and the Get API does not echo the user-supplied values.
				ImportStateVerifyIgnore: []string{"config"},
			},
		},
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

resource "alicloud_log_project" "default" {
    project_name = "${var.name}-log"
}

resource "alicloud_log_store" "default" {
    project_name  = alicloud_log_project.default.project_name
    logstore_name = "${var.name}-store"
}

resource "alicloud_log_project" "default2" {
    project_name = "${var.name}-log2"
}

resource "alicloud_log_store" "default2" {
    project_name  = alicloud_log_project.default2.project_name
    logstore_name = "${var.name}-store2"
}

`, name)
}

// Test Agentloop ContextStore. <<< Resource test cases, automatically generated.
