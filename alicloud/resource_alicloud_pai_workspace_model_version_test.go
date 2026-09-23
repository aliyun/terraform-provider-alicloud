// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace ModelVersion. >>> Resource test cases, automatically generated.
// Case ModelVersion_20241223_副本1740484379759 10362
func TestAccAliCloudPaiWorkspaceModelVersion_basic10362(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_model_version.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceModelVersionMap10362)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceModelVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceModelVersionBasicDependence10362)
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
					"uri":      "oss://hz-test-0701.oss-cn-hangzhou-internal.aliyuncs.com/checkpoints/",
					"model_id": "${alicloud_pai_workspace_model.defaultsHptEL.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uri":      "oss://hz-test-0701.oss-cn-hangzhou-internal.aliyuncs.com/checkpoints/",
						"model_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"approval_status": "Approved",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"approval_status": "Approved",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"extra_info": map[string]interface{}{
						"test": "ExtraInfo",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"extra_info.%":    "1",
						"extra_info.test": "ExtraInfo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inference_spec": map[string]interface{}{
						"test": "InferenceSpec",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inference_spec.%":    "1",
						"inference_spec.test": "InferenceSpec",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"key":   "k1",
							"value": "v1",
						},
						{
							"key":   "k2",
							"value": "v2",
						},
						{
							"key":   "k3",
							"value": "v3",
						},
						{
							"key":   "k4",
							"value": "v4",
						},
						{
							"key":   "k5",
							"value": "v5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"key":   "k1",
							"value": "v1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"key":   "k1",
							"value": "v1",
						},
						{
							"key":   "k2",
							"value": "v2",
						},
						{
							"key":   "k3",
							"value": "v3",
						},
						{
							"key":   "k4",
							"value": "v4",
						},
						{
							"key":   "k5",
							"value": "v5",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"metrics": map[string]interface{}{
						"test": "Metrics",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metrics.%":    "1",
						"metrics.test": "Metrics",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": `{\"test\":\"options\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_id": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_id": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type": "TrainingService",
					"source_id":   "region=$${region},workspaceId=$${workspace_id},kind=TrainingJob,id=job-id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type": "TrainingService",
						"source_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"training_spec": map[string]interface{}{
						"test": "TrainingSpec",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"training_spec.%":    "1",
						"training_spec.test": "TrainingSpec",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_description": name,
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

func TestAccAliCloudPaiWorkspaceModelVersion_basic10362_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_model_version.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceModelVersionMap10362)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceModelVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceModelVersionBasicDependence10362)
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
					"approval_status": "Approved",
					"extra_info": map[string]interface{}{
						"test": "ExtraInfo",
					},
					"format_type":    "SavedModel",
					"framework_type": "PyTorch",
					"inference_spec": map[string]interface{}{
						"test": "InferenceSpec",
					},
					"labels": []map[string]interface{}{
						{
							"key":   "k1",
							"value": "v1",
						},
						{
							"key":   "k2",
							"value": "v2",
						},
						{
							"key":   "k3",
							"value": "v3",
						},
						{
							"key":   "k4",
							"value": "v4",
						},
						{
							"key":   "k5",
							"value": "v5",
						},
					},
					"metrics": map[string]interface{}{
						"test": "Metrics",
					},
					"model_id":    "${alicloud_pai_workspace_model.defaultsHptEL.id}",
					"options":     `{\"test\":\"options\"}`,
					"source_type": "TrainingService",
					"source_id":   "region=$${region},workspaceId=$${workspace_id},kind=TrainingJob,id=job-id",
					"training_spec": map[string]interface{}{
						"test": "TrainingSpec",
					},
					"uri":                 "oss://hz-test-0701.oss-cn-hangzhou-internal.aliyuncs.com/checkpoints/",
					"version_description": name,
					"version_name":        "1.0.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"approval_status":     "Approved",
						"extra_info.%":        "1",
						"extra_info.test":     "ExtraInfo",
						"format_type":         "SavedModel",
						"framework_type":      "PyTorch",
						"inference_spec.%":    "1",
						"inference_spec.test": "InferenceSpec",
						"labels.#":            "5",
						"metrics.%":           "1",
						"metrics.test":        "Metrics",
						"model_id":            CHECKSET,
						"options":             CHECKSET,
						"source_type":         "TrainingService",
						"source_id":           CHECKSET,
						"training_spec.%":     "1",
						"training_spec.test":  "TrainingSpec",
						"uri":                 "oss://hz-test-0701.oss-cn-hangzhou-internal.aliyuncs.com/checkpoints/",
						"version_description": name,
						"version_name":        "1.0.0",
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

var AliCloudPaiWorkspaceModelVersionMap10362 = map[string]string{
	"approval_status": CHECKSET,
	"source_type":     CHECKSET,
	"version_name":    CHECKSET,
}

func AliCloudPaiWorkspaceModelVersionBasicDependence10362(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultDI9fsL" {
  description    = "382"
  workspace_name = var.name
  env_types      = ["prod"]
  display_name   = var.name
}

resource "alicloud_pai_workspace_model" "defaultsHptEL" {
  model_name        = var.name
  workspace_id      = alicloud_pai_workspace_workspace.defaultDI9fsL.id
  origin            = "Civitai"
  task              = "text-to-image-synthesis"
  accessibility     = "PRIVATE"
  model_type        = "Checkpoint"
  order_number      = "1"
  model_description = "ModelDescription."
  model_doc         = "https://eas-***.oss-cn-hangzhou.aliyuncs.com/s**.safetensors"
  domain            = "aigc"
  labels {
    key   = "base_model"
    value = "SD 1.5"
  }
  extra_info = {
    test = "15"
  }
}


`, name)
}

// Test PaiWorkspace ModelVersion. <<< Resource test cases, automatically generated.
