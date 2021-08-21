package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostUserGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_user_group.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostUserGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostUserGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhostusergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostUserGroupBasicDependence0)
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
					"user_group_name": "tf-testAcc-0T2Sep=samLLheEIbZ",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":     CHECKSET,
						"user_group_name": "tf-testAcc-0T2Sep=samLLheEIbZ",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "tf-testAcc-6ke&*Cfo/6lOS@jj.o#KRgf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "tf-testAcc-6ke&*Cfo/6lOS@jj.o#KRgf",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_group_name": "tf-testAcc-r]L_,Zap@tCCFG2L8<xI5~IA",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_group_name": "tf-testAcc-r]L_,Zap@tCCFG2L8<xI5~IA",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment":         "tf-testAcc-I[{q!E,7l?W{1)Uf7<,wz]",
					"user_group_name": "tf-testAcc--rw|aCqxJ@ILzv_OOSedz?",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment":         "tf-testAcc-I[{q!E,7l?W{1)Uf7<,wz]",
						"user_group_name": "tf-testAcc--rw|aCqxJ@ILzv_OOSedz?",
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

var AlicloudBastionhostUserGroupMap0 = map[string]string{
	"comment":         "",
	"user_group_id":   CHECKSET,
	"instance_id":     CHECKSET,
	"user_group_name": "tf-testAcc-0T2Sep=samLLheEIbZ",
}

func AlicloudBastionhostUserGroupBasicDependence0(name string) string {
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
