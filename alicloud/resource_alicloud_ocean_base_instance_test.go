package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudOceanBaseInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_instance.default"
	checkoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudOceanBaseInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OceanBaseProService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOceanBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOceanBaseInstanceBasicDependence0)
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
					"instance_name":      "${var.name}",
					"series":             "normal",
					"auto_renew":         "false",
					"disk_size":          "100",
					"payment_type":       "PayAsYouGo",
					"instance_class":     "8C32GB",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"zones":              []string{"${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 2]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 3]}", "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 4]}"},
					"backup_retain_mode": "delete_all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":     name,
						"series":            "normal",
						"auto_renew":        "false",
						"disk_size":         "100",
						"instance_class":    "8C32GB",
						"resource_group_id": CHECKSET,
						"zones.#":           "3",
						"payment_type":      "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_num": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_num": "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size":      "200",
					"instance_class": "14C70GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size":      "200",
						"instance_class": "14C70GB",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period_unit", "auto_renew_period", "backup_retain_mode", "period"},
			},
		},
	})
}

var AlicloudOceanBaseInstanceMap0 = map[string]string{
	"status":  CHECKSET,
	"zones.#": CHECKSET,
}

func AlicloudOceanBaseInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {}


data "alicloud_resource_manager_resource_groups" "default"{}

`, name)
}
