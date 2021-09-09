package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHostAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostAttachmentBasicDependence0)
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
					"instance_id":   "${alicloud_bastionhost_host_group.default.instance_id}",
					"host_group_id": "${alicloud_bastionhost_host_group.default.host_group_id}",
					"host_id":       "${alicloud_bastionhost_host.default.host_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

var AlicloudBastionhostHostAttachmentMap0 = map[string]string{
	"host_id":       CHECKSET,
	"instance_id":   CHECKSET,
	"host_group_id": CHECKSET,
}

func AlicloudBastionhostHostAttachmentBasicDependence0(name string) string {
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
resource "alicloud_bastionhost_host" "default" {
 instance_id          = alicloud_bastionhost_instance.default.id
 host_name            = var.name
 active_address_type  = "Private"
 host_private_address = "172.16.0.10"
 os_type              = "Linux"
 source               = "Local"
}
resource "alicloud_bastionhost_host_group" "default" {
 instance_id          = alicloud_bastionhost_instance.default.id
 host_group_name = var.name
}
`, name)
}
