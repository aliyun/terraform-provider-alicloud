package alicloud

import (
	"encoding/json"
	"fmt"
	"reflect"
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
				// Strict import check: the Read restores version from currentVersion
				// and version_description from the matching entry of versions[], so
				// the imported state must match the applied state attribute-for-
				// attribute and the post-import plan stays empty.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAccAliCloudExpandAgentloopEvaluatorConfig(t *testing.T) {
	raw := map[string]interface{}{
		"prompt":          `{"json":"example"}`,
		"variables":       `[{"name":"input"}]`,
		"variablesPadded": ` [1] `,
		"outputSchema":    `{"type":"object"}`,
		"faasConfig":      "not-json",
		"parameters":      `[broken`,
		"codeType":        "python",
		"number":          "10",
		"nonString":       1,
	}
	got := expandAgentloopEvaluatorConfig(raw)

	if got["prompt"] != `{"json":"example"}` {
		t.Errorf("prompt = %#v, want verbatim string (not a structured key)", got["prompt"])
	}
	if want := []interface{}{map[string]interface{}{"name": "input"}}; !reflect.DeepEqual(got["variables"], want) {
		t.Errorf("variables = %#v, want %#v", got["variables"], want)
	}
	if want := map[string]interface{}{"type": "object"}; !reflect.DeepEqual(got["outputSchema"], want) {
		t.Errorf("outputSchema = %#v, want %#v", got["outputSchema"], want)
	}
	if got["faasConfig"] != "not-json" {
		t.Errorf("faasConfig = %#v, want verbatim string", got["faasConfig"])
	}
	if got["parameters"] != `[broken` {
		t.Errorf("parameters = %#v, want verbatim invalid JSON", got["parameters"])
	}
	if got["codeType"] != "python" {
		t.Errorf("codeType = %#v, want verbatim string", got["codeType"])
	}
	if got["number"] != "10" {
		t.Errorf("number = %#v, want verbatim string (scalar JSON is not expanded)", got["number"])
	}
	if got["nonString"] != 1 {
		t.Errorf("nonString = %#v, want 1", got["nonString"])
	}
	// "variablesPadded" is not in the structured-keys whitelist, so it must be
	// sent verbatim even though it parses as JSON.
	if got["variablesPadded"] != ` [1] ` {
		t.Errorf("variablesPadded = %#v, want verbatim string (not whitelisted)", got["variablesPadded"])
	}

	if got := expandAgentloopEvaluatorConfig("not-a-map"); len(got) != 0 {
		t.Errorf("non map input: got %v entries, want empty map", len(got))
	}

	// A structured key carrying two JSON values must stay verbatim instead of
	// being silently truncated to the first document.
	multi := expandAgentloopEvaluatorConfig(map[string]interface{}{
		"variables": `[{"name":"a"}] [{"name":"b"}]`,
	})
	if multi["variables"] != `[{"name":"a"}] [{"name":"b"}]` {
		t.Errorf("variables multi-doc = %#v, want verbatim multi-document string", multi["variables"])
	}
}

func TestAccAliCloudFlattenAgentloopEvaluatorConfig(t *testing.T) {
	raw := map[string]interface{}{
		"prompt":       "text",
		"variables":    []interface{}{"a"},
		"outputSchema": map[string]interface{}{"type": "object"},
		"bigNumber":    []interface{}{json.Number("9007199254740993")},
	}
	got := flattenAgentloopEvaluatorConfig(raw, nil)
	if got["prompt"] != "text" {
		t.Errorf("prompt = %q, want verbatim string", got["prompt"])
	}
	if got["variables"] != `["a"]` {
		t.Errorf("variables = %q, want serialized array", got["variables"])
	}
	if got["outputSchema"] != `{"type":"object"}` {
		t.Errorf("outputSchema = %q, want serialized object", got["outputSchema"])
	}
	// json.Number marshals back to its exact literal without float64 rounding.
	if got["bigNumber"] != `[9007199254740993]` {
		t.Errorf("bigNumber = %q, want exact large-integer literal", got["bigNumber"])
	}

	if got := flattenAgentloopEvaluatorConfig("not-a-map", nil); len(got) != 0 {
		t.Errorf("non map input: got %v entries, want empty map", len(got))
	}

	// Semantically equal JSON keeps the configured text verbatim, so a
	// differently formatted configuration does not become a perpetual diff;
	// semantically different JSON falls back to the canonical API form.
	configured := map[string]interface{}{
		"variables":    `[ "a" ]`,
		"outputSchema": `{"type":"object","extra":true}`,
	}
	gotConfigured := flattenAgentloopEvaluatorConfig(raw, configured)
	if gotConfigured["variables"] != `[ "a" ]` {
		t.Errorf("variables (configured) = %q, want configured text kept verbatim", gotConfigured["variables"])
	}
	if gotConfigured["outputSchema"] != `{"type":"object"}` {
		t.Errorf("outputSchema (configured) = %q, want canonical API serialization", gotConfigured["outputSchema"])
	}
}
