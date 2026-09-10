// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms PrometheusView. >>> Resource test cases, automatically generated.
// Case prom-view 9280
func TestAccAliCloudCmsPrometheusView_basic9280(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_prometheus_view.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsPrometheusViewMap9280)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsPrometheusView")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsPrometheusViewBasicDependence9280)
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
					"prometheus_view_name": name,
					"version":              "V2",
					"prometheus_instances": []map[string]interface{}{
						{
							"prometheus_instance_id": "${alicloud_cms_prometheus_instance.default.0.id}",
							"region_id":              "${alicloud_cms_prometheus_instance.default.0.region_id}",
							"user_id":                "${data.alicloud_account.default.id}",
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_view_name":   name,
						"version":                "V2",
						"prometheus_instances.#": "1",
						"workspace":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_instances": []map[string]interface{}{
						{
							"prometheus_instance_id": "${alicloud_cms_prometheus_instance.default.0.id}",
							"region_id":              "${alicloud_cms_prometheus_instance.default.0.region_id}",
							"user_id":                "${data.alicloud_account.default.id}",
						},
						{
							"prometheus_instance_id": "${alicloud_cms_prometheus_instance.default.1.id}",
							"region_id":              "${alicloud_cms_prometheus_instance.default.1.region_id}",
							"user_id":                "${data.alicloud_account.default.id}",
						},
						{
							"prometheus_instance_id": "${alicloud_cms_prometheus_instance.default.2.id}",
							"region_id":              "${alicloud_cms_prometheus_instance.default.2.region_id}",
							"user_id":                "${data.alicloud_account.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_instances.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_instances": []map[string]interface{}{
						{
							"prometheus_instance_id": "${alicloud_cms_prometheus_instance.default.1.id}",
							"region_id":              "${alicloud_cms_prometheus_instance.default.1.region_id}",
							"user_id":                "${data.alicloud_account.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_instances.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_free_read_policy": "1.1.1.1",
					"enable_auth_free_read": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_free_read_policy": "1.1.1.1",
						"enable_auth_free_read": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auth_free_read_policy": "2.2.2.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auth_free_read_policy": "2.2.2.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auth_free_read": "false",
					"auth_free_read_policy": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auth_free_read": "false",
						"auth_free_read_policy": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus_view_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_view_name": name + "_update",
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

func TestAccAliCloudCmsPrometheusView_basic9280_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_prometheus_view.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsPrometheusViewMap9280)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsPrometheusView")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsPrometheusViewBasicDependence9280)
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
					"prometheus_view_name": name,
					"version":              "V2",
					"prometheus_instances": []map[string]interface{}{
						{
							"prometheus_instance_id": "${alicloud_cms_prometheus_instance.default.0.id}",
							"region_id":              "${alicloud_cms_prometheus_instance.default.0.region_id}",
							"user_id":                "${data.alicloud_account.default.id}",
						},
						{
							"prometheus_instance_id": "${alicloud_cms_prometheus_instance.default.1.id}",
							"region_id":              "${alicloud_cms_prometheus_instance.default.1.region_id}",
							"user_id":                "${data.alicloud_account.default.id}",
						},
						{
							"prometheus_instance_id": "${alicloud_cms_prometheus_instance.default.2.id}",
							"region_id":              "${alicloud_cms_prometheus_instance.default.2.region_id}",
							"user_id":                "${data.alicloud_account.default.id}",
						},
					},
					"workspace":             "${alicloud_cms_workspace.default.id}",
					"auth_free_read_policy": "1.1.1.1",
					"enable_auth_free_read": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prometheus_view_name":   name,
						"version":                "V2",
						"prometheus_instances.#": "3",
						"workspace":              CHECKSET,
						"auth_free_read_policy":  "1.1.1.1",
						"enable_auth_free_read":  "true",
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

var AliCloudCmsPrometheusViewMap9280 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudCmsPrometheusViewBasicDependence9280(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "default" {
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_prometheus_instance" "default" {
  count                    = 3
  prometheus_instance_name = "${var.name}_${count.index}"
  workspace                = alicloud_cms_workspace.default.id
}
`, name)
}

// Test Cms PrometheusView. <<< Resource test cases, automatically generated.
