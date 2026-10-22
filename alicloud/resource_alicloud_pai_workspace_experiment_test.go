package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckPaiWorkspaceExperimentLabelsEmpty(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		object, err := (&PaiWorkspaceServiceV2{client}).DescribePaiWorkspaceExperiment(rs.Primary.ID)
		if err != nil {
			return err
		}
		if labels, ok := object["Labels"].([]interface{}); ok && len(labels) > 0 {
			return fmt.Errorf("expected no labels to be stored for %s, got %#v", rs.Primary.ID, labels)
		}
		return nil
	}
}

// Test PaiWorkspace Experiment. >>> Resource test cases, automatically generated.
// Case 实验测试_副本1732006685057 9003
func TestAccAliCloudPaiWorkspaceExperiment_basic9003(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_experiment.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceExperimentMap9003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceExperiment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceExperimentBasicDependence9003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Regression for QA F-1: an empty labels block used to reach Apply with a
			// nil element and panic the Create path. The nil-element guard now skips
			// it, labels.key stays Optional (no required/non-empty validation), and
			// the create must succeed without sending an empty Labels payload.
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PRIVATE",
					"artifact_uri":    "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/test/",
					"experiment_name": name,
					"workspace_id":    "${alicloud_pai_workspace_workspace.defaultDI9fsL.id}",
					"labels":          []map[string]interface{}{{}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PRIVATE",
						"artifact_uri":    "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/test/",
						"experiment_name": name,
						"workspace_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "env",
							"value": "prod",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
						"labels.#":        "1",
						"labels.0.key":    "env",
						"labels.0.value":  "prod",
					}),
				),
			},
			// Regression for review P2 (comment 43917): switching to an empty labels
			// block must delete the previous labels without sending a
			// {"Key":"","Value":""} payload.
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels":          []map[string]interface{}{{}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
						"labels.#":        REMOVEKEY,
						"labels.0.key":    REMOVEKEY,
						"labels.0.value":  REMOVEKEY,
					}),
					testAccCheckPaiWorkspaceExperimentLabelsEmpty(resourceId),
				),
			},
			// Regression for review P2 (comment 43917): a block that only sets value
			// has no usable key and must be skipped instead of sending Key="".
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels": []map[string]interface{}{
						{
							"value": "dev",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
					}),
					testAccCheckPaiWorkspaceExperimentLabelsEmpty(resourceId),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "alpha",
							"value": "one",
						},
						{
							"key":   "beta",
							"value": "two",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
						"labels.#":        "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "beta",
							"value": "two",
						},
						{
							"key":   "alpha",
							"value": "one",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "beta",
							"value": "two",
						},
						{
							"key":   "alpha",
							"value": "one",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
						"labels.#":        "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels":          REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
						"labels.#":        "0",
					}),
				),
			},
			// Regression for QA F-2: the delete-label request path must encode the key,
			// otherwise a key like "team#component" truncates the URL at the fragment.
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels": []map[string]interface{}{
						{
							"key":   "team#component",
							"value": "ops",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
						"labels.#":        "1",
						"labels.0.key":    "team#component",
						"labels.0.value":  "ops",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels":          REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
						"labels.#":        "0",
						"labels.0.key":    REMOVEKEY,
						"labels.0.value":  REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudPaiWorkspaceExperimentMap9003 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceExperimentBasicDependence9003(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultDI9fsL" {
  description    = "902"
  display_name   = "test_pop_800"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Case 实验测试 6043
func TestAccAliCloudPaiWorkspaceExperiment_basic6043(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_experiment.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceExperimentMap6043)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceExperiment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceExperimentBasicDependence6043)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PRIVATE",
					"artifact_uri":    "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/test/",
					"experiment_name": name,
					"workspace_id":    "${alicloud_pai_workspace_workspace.defaultDI9fsL.id}",
					// Regression for review P2 (comment 43917): a create whose labels
					// block only sets value has no usable key and must not send Key="".
					"labels": []map[string]interface{}{
						{
							"value": "dev",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PRIVATE",
						"artifact_uri":    "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/test/",
						"experiment_name": name,
						"workspace_id":    CHECKSET,
					}),
					testAccCheckPaiWorkspaceExperimentLabelsEmpty(resourceId),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
					"labels":          REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
						"labels.#":        "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudPaiWorkspaceExperimentMap6043 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceExperimentBasicDependence6043(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultDI9fsL" {
  description    = "418"
  display_name   = "test_pop_923"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

func TestUnitPaiWorkspaceExperimentLabelsExpand(t *testing.T) {
	testCases := []struct {
		name     string
		raw      interface{}
		expected int
	}{
		{"nil input", nil, 0},
		{"empty block", []interface{}{map[string]interface{}{}}, 0},
		{"value only", []interface{}{map[string]interface{}{"value": "dev"}}, 0},
		{"empty key", []interface{}{map[string]interface{}{"key": "", "value": "dev"}}, 0},
		{"nil element", []interface{}{nil}, 0},
		{"non-map element", []interface{}{"label"}, 0},
		{"valid key", []interface{}{map[string]interface{}{"key": "env", "value": "prod"}}, 1},
		{"mixed", []interface{}{map[string]interface{}{"value": "dev"}, map[string]interface{}{"key": "env", "value": "prod"}}, 1},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := expandPaiWorkspaceExperimentLabels(tc.raw)
			if len(got) != tc.expected {
				t.Fatalf("expected %d labels, got %d: %#v", tc.expected, len(got), got)
			}
			for _, item := range got {
				label, ok := item.(map[string]interface{})
				if !ok {
					t.Fatalf("expected a label map, got %#v", item)
				}
				if label["Key"] != "env" || label["Value"] != "prod" {
					t.Fatalf("unexpected label payload: %#v", label)
				}
			}
		})
	}
}

// Test PaiWorkspace Experiment. <<< Resource test cases, automatically generated.
