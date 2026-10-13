package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDataWorksTenantResourceGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_tenant_resource_group.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksTenantResourceGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksTenantResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_trg%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksTenantResourceGroupBasicDependence)
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
					"tenant_resource_group_name":        name,
					"tenant_resource_group_description": "test_description",
					"payment_type":                      "PostPaid",
					"vpc_id":                            "${alicloud_vpc.default.id}",
					"vswitch_id":                        "${alicloud_vswitch.default.id}",
					"aliyun_resource_group_id":          "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"payment_duration":                  0,
					"payment_duration_unit":             "Month",
					"auto_renew_enabled":                false,
					"spec":                              0,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "TenantRG",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tenant_resource_group_name":        name,
						"tenant_resource_group_description": "test_description",
						"payment_type":                      "PostPaid",
						"tags.%":                            "2",
						"tags.Created":                      "TF",
						"tags.For":                          "TenantRG",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tenant_resource_group_description": "test_description_updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tenant_resource_group_description": "test_description_updated",
						"tags.For":                          "TenantRG",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": CLEARMAP,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"payment_duration_unit",
				},
			},
		},
	})
}

var AlicloudDataWorksTenantResourceGroupMap = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudDataWorksTenantResourceGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "default" {
  vpc_name   = format("%%s_vpc", var.name)
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = "cn-hangzhou-h"
  vswitch_name = format("%%s_vsw", var.name)
}

data "alicloud_resource_manager_resource_groups" "default" {}
`, name)
}
