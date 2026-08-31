package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace Model. >>> Resource test cases, automatically generated.

// Case 模型测试用例_20241223 9669
func TestAccAliCloudPaiWorkspaceModel_basic9669(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_model.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceModelMap9669)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceModel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceModelBasicDependence9669)
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
					"model_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"model_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility": "PUBLIC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility": "PUBLIC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain": "aigc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain": "aigc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"extra_info": map[string]interface{}{
						"test": "15",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"extra_info.%":    "1",
						"extra_info.test": "15",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"key":   "base_model",
							"value": "SD 1.5",
						},
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
							"key":   "base_model",
							"value": "SD 1.5",
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
							"key":   "base_model",
							"value": "SD 1.5",
						},
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
					"model_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"model_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"model_doc": "https://***.md",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"model_doc": "https://***.md",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"model_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"model_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"model_type": "Checkpoint",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"model_type": "Checkpoint",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_number": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_number": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"origin": "Civitai",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"origin": "Civitai",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task": "text-to-image-synthesis",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task": "text-to-image-synthesis",
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

func TestAccAliCloudPaiWorkspaceModel_basic9669_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_model.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceModelMap9669)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceModel")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceModelBasicDependence9669)
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
					"model_name":        name,
					"workspace_id":      "${alicloud_pai_workspace_workspace.defaultDI9fsL.id}",
					"origin":            "Civitai",
					"task":              "text-to-image-synthesis",
					"accessibility":     "PUBLIC",
					"model_type":        "Checkpoint",
					"order_number":      "1",
					"model_description": "ModelDescription.",
					"model_doc":         "https://eas-***.oss-cn-hangzhou.aliyuncs.com/s**.safetensors",
					"domain":            "aigc",
					"labels": []map[string]interface{}{
						{
							"key":   "base_model",
							"value": "SD 1.5",
						},
						{
							"key":   "k1",
							"value": "v1",
						},
						{
							"key":   "k2",
							"value": "v2",
						},
					},
					"extra_info": map[string]interface{}{
						"test": "15",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"model_name":        name,
						"workspace_id":      CHECKSET,
						"origin":            "Civitai",
						"task":              "text-to-image-synthesis",
						"accessibility":     "PUBLIC",
						"model_type":        "Checkpoint",
						"order_number":      "1",
						"model_description": "ModelDescription.",
						"model_doc":         "https://eas-***.oss-cn-hangzhou.aliyuncs.com/s**.safetensors",
						"domain":            "aigc",
						"labels.#":          "3",
						"extra_info.%":      "1",
						"extra_info.test":   "15",
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

var AliCloudPaiWorkspaceModelMap9669 = map[string]string{}

func AliCloudPaiWorkspaceModelBasicDependence9669(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultDI9fsL" {
  description    = "826"
  workspace_name = var.name
  env_types      = ["prod"]
  display_name   = var.name
}


`, name)
}

// Test PaiWorkspace Model. <<< Resource test cases, automatically generated.
