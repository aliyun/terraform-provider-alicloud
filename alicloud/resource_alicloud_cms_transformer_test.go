package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsTransformer_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_transformer.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsTransformerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsTransformer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsTransformerBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"transformer_name": name,
					"description":      "Transformer managed by Terraform",
					"quit_after_match": "false",
					"sort_id":          "1",
					"enable":           "true",
					"workspace":        "${alicloud_cms_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transformer_name": name,
						"description":      "Transformer managed by Terraform",
						"quit_after_match": "false",
						"sort_id":          "1",
						"enable":           "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transformer_name": name + "-updated",
					"description":      "Updated transformer description",
					"quit_after_match": "true",
					"sort_id":          "2",
					"enable":           "false",
					"workspace":        "${alicloud_cms_workspace.default.id}",
					"actions": []map[string]interface{}{
						{
							"label_key": "event.name",
							"source":    "$.data.eventName",
							"target":    "DATA",
							"type":      "SET_FIELD",
							"value":     "updated-value",
							"variable":  "eventName",
							"mapping": map[string]interface{}{
								"key1": "value1updated",
							},
							"filter_setting": []map[string]interface{}{
								{
									"relation":   "OR",
									"expression": "$.level",
									"conditions": []map[string]interface{}{
										{
											"field": "level",
											"op":    "NE",
											"value": "INFO",
										},
									},
								},
							},
						},
						{
							"label_key": "event.name",
							"source":    "$.data.eventName",
							"target":    "SUBJECT",
							"type":      "REGEXP_EXTRACT",
							"reg_exp":   "^updated.*",
						},
					},
					"filter_setting": []map[string]interface{}{
						{
							"relation":   "OR",
							"expression": "$.level",
							"conditions": []map[string]interface{}{
								{
									"field": "level",
									"op":    "NE",
									"value": "INFO",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transformer_name": name + "-updated",
						"description":      "Updated transformer description",
						"quit_after_match": "true",
						"sort_id":          "2",
						"enable":           "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transformer_name": name,
					"description":      REMOVEKEY,
					"quit_after_match": "false",
					"sort_id":          "3",
					"enable":           "true",
					"workspace":        "${alicloud_cms_workspace.default.id}",
					"actions":          []interface{}{},
					"filter_setting":   []interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transformer_name": name,
						"sort_id":          "3",
						"enable":           "true",
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

func TestAccAliCloudCmsTransformer_minimum(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_transformer.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsTransformerMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsTransformer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsTransformerBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"transformer_name": name,
					"workspace":        "${alicloud_cms_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transformer_name": name,
					}),
				),
			},
		},
	})
}

func AlicloudCmsTransformerBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}
`, name)
}

var AlicloudCmsTransformerMap = map[string]string{
	"transformer_name": CHECKSET,
	"workspace":        CHECKSET,
}
