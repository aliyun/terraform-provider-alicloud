package alicloud

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	tfSchema "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
									"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
									"start_time":  "2026-02-01T00:00:00Z",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"context_type":                 "experience",
						"description":                  "terraform自动化测试描述更新",
						"config.#":                     "1",
						"config.0.service_names.#":     "1",
						"config.0.service_names.0":     "svc-tf-acc-updated",
						"config.0.source.0.start_time": "2026-02-01T00:00:00Z",
						"config.0.metadata_field.%":    "2",
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
						"config.0.service_names.0": "svc-tf-acc-updated",
						"config.0.service_names.1": "svc-tf-acc",
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
				// Applying the reordered service_names must then converge with an
				// empty post-apply plan; serviceNames are read back in API order.
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
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// service_names and source are read back; metadataField is
				// write-only (the Get API returns internal keys instead of the
				// user mapping), so it cannot be recovered on import.
				ImportStateVerifyIgnore: []string{"config.0.metadata_field"},
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

`, name)
}

// Test Agentloop ContextStore. <<< Resource test cases, automatically generated.

func TestAccAliCloudAgentloopContextStoreFlattenConfig(t *testing.T) {
	apiConfig := map[string]interface{}{
		"metadataField": map[string]interface{}{
			"defaultRecallPolicyId":     "p-1",
			"experienceMetadataVersion": "v1",
		},
		"miningInterval": "1d",
		"serviceNames":   []interface{}{"svc-b", "svc-a"},
		"source": map[string]interface{}{
			"agentSpace": "as",
			"startTime":  "2026-01-01T00:00:00Z",
		},
	}

	// Import: no prior state, metadata_field cannot be recovered.
	got := flattenAgentloopContextStoreConfig(apiConfig, []interface{}{})
	want := []map[string]interface{}{{
		"service_names": []interface{}{"svc-b", "svc-a"},
		"source": []interface{}{map[string]interface{}{
			"agent_space": "as",
			"start_time":  "2026-01-01T00:00:00Z",
		}},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("import flatten = %#v, want %#v", got, want)
	}

	// Refresh: the configured mapping and an equivalent time spelling are kept.
	prior := []interface{}{map[string]interface{}{
		"metadata_field": map[string]interface{}{"userId": "user_id"},
		"service_names":  []interface{}{"svc-b", "svc-a"},
		"source": []interface{}{map[string]interface{}{
			"agent_space": "as",
			"start_time":  "2026-01-01T08:00:00+08:00",
		}},
	}}
	got = flattenAgentloopContextStoreConfig(apiConfig, prior)
	if m := got[0]["metadata_field"]; !reflect.DeepEqual(m, map[string]interface{}{"userId": "user_id"}) {
		t.Errorf("metadata_field = %#v, want prior mapping", m)
	}
	source := got[0]["source"].([]interface{})[0].(map[string]interface{})
	if source["start_time"] != "2026-01-01T08:00:00+08:00" {
		t.Errorf("start_time = %#v, want configured spelling", source["start_time"])
	}

	// Drift: a different instant is surfaced.
	prior[0].(map[string]interface{})["source"] = []interface{}{map[string]interface{}{
		"agent_space": "as",
		"start_time":  "2026-02-01T00:00:00Z",
	}}
	got = flattenAgentloopContextStoreConfig(apiConfig, prior)
	source = got[0]["source"].([]interface{})[0].(map[string]interface{})
	if source["start_time"] != "2026-01-01T00:00:00Z" {
		t.Errorf("start_time = %#v, want API value", source["start_time"])
	}

	if got := flattenAgentloopContextStoreConfig(map[string]interface{}{"miningInterval": "1d"}, nil); len(got) != 0 {
		t.Errorf("server-only config = %#v, want empty", got)
	}
	if got := flattenAgentloopContextStoreConfig(nil, nil); len(got) != 0 {
		t.Errorf("nil config = %#v, want empty", got)
	}
}

// TestAccAliCloudAgentloopContextStoreImportPlan builds the state produced by
// an import Read and checks that re-applying the original service_names and
// source plans no replacement.
func TestAccAliCloudAgentloopContextStoreImportPlan(t *testing.T) {
	r := resourceAliCloudAgentloopContextStore()
	imported := map[string]interface{}{
		"agent_space":        "as",
		"context_store_name": "cs",
		"context_type":       "experience",
	}
	data := tfSchema.TestResourceDataRaw(t, r.Schema, imported)
	data.SetId("as:cs")
	if err := data.Set("config", flattenAgentloopContextStoreConfig(map[string]interface{}{
		"serviceNames": []interface{}{"svc"},
		"source":       map[string]interface{}{"agentSpace": "as", "startTime": "2026-01-01T00:00:00Z"},
	}, nil)); err != nil {
		t.Fatalf("set config: %v", err)
	}
	state := data.State()

	cfg := map[string]interface{}{
		"agent_space":        "as",
		"context_store_name": "cs",
		"context_type":       "experience",
		"config": []interface{}{map[string]interface{}{
			"service_names": []interface{}{"svc"},
			"source": []interface{}{map[string]interface{}{
				"agent_space": "as",
				"start_time":  "2026-01-01T00:00:00Z",
			}},
		}},
	}
	diff, err := r.Diff(state, terraform.NewResourceConfigRaw(cfg), nil)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}
	if diff != nil && diff.RequiresNew() {
		t.Fatalf("imported config must not require replacement: %#v", diff.Attributes)
	}

	cfg["config"] = []interface{}{map[string]interface{}{
		"service_names": []interface{}{"svc-other"},
	}}
	diff, err = r.Diff(state, terraform.NewResourceConfigRaw(cfg), nil)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}
	if diff == nil || !diff.RequiresNew() {
		t.Fatalf("service_names drift must require replacement, got %#v", diff)
	}
}
