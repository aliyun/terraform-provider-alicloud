package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHostGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_group.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostgroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostGroupBasicDependence0)
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
					"instance_id":     "${alicloud_bastionhost_instance.default.id}",
					"host_group_name": "tf-testaccHostGroupName12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":     CHECKSET,
						"host_group_name": "tf-testaccHostGroupName12345",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "tf-testaccHostGroupComment12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "tf-testaccHostGroupComment12345",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_group_name": "tf-testaccHostGroupName12345-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_group_name": "tf-testaccHostGroupName12345-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":         "tf-testaccHostGroupComment12345update",
					"host_group_name": "tf-testaccHostGroupName12345",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":         "tf-testaccHostGroupComment12345update",
						"host_group_name": "tf-testaccHostGroupName12345",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudBastionhostHostGroupMap0 = map[string]string{
	"host_group_id": CHECKSET,
	"instance_id":   CHECKSET,
}

func AlicloudBastionhostHostGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
 available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
 name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
 zone_id = local.zone_id
 vpc_id  = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_vswitch" "this" {
 count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id       = data.alicloud_vpcs.default.ids.0
 zone_id      = data.alicloud_zones.default.ids.0
 cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
resource "alicloud_security_group" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 name   = var.name
}
locals {
 vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
 zone_id     = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
}
resource "alicloud_bastionhost_instance" "default" {
 description        = var.name
 license_code       = "bhah_ent_50_asset"
 period             = "1"
 vswitch_id         = local.vswitch_id
 security_group_ids = [alicloud_security_group.default.id]
}
`, name)
}
