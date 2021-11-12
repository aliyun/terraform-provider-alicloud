package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPrivateZoneRuleAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_rule_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateZoneRuleAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePvtzRuleAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateZoneRuleAttachmentBasicDependence0)
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
					"rule_id": "${alicloud_pvtz_rule.default.id}",
					"vpcs": []map[string]interface{}{
						{
							"vpc_id":    "${alicloud_vpc.default.0.id}",
							"region_id": "${var.region}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_id": CHECKSET,
						"vpcs.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpcs": []map[string]interface{}{
						{
							"vpc_id":    "${alicloud_vpc.default.0.id}",
							"region_id": "${var.region}",
						},
						{
							"vpc_id":    "${alicloud_vpc.default.1.id}",
							"region_id": "${var.region}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpcs.#": "2",
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

var AlicloudPrivateZoneRuleAttachmentMap0 = map[string]string{}

func AlicloudPrivateZoneRuleAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

resource "alicloud_vpc" "default" {
  count      = 3
  vpc_name   = var.name
  cidr_block = cidrsubnet("172.16.0.0/16", 8, count.index)
}

data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}

resource "alicloud_vswitch" "default" {
  count      = 2
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.2.cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.2.id
  name   = var.name
}

resource "alicloud_pvtz_endpoint" "default" {
  endpoint_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.2.id
  vpc_region_id     = var.region
  ip_configs {
    zone_id    = alicloud_vswitch.default[0].zone_id
    cidr_block = alicloud_vswitch.default[0].cidr_block
    vswitch_id = alicloud_vswitch.default[0].id
  }
  ip_configs {
    zone_id    = alicloud_vswitch.default[1].zone_id
    cidr_block = alicloud_vswitch.default[1].cidr_block
    vswitch_id = alicloud_vswitch.default[1].id
  }

}

resource "alicloud_pvtz_rule" "default" {
  endpoint_id = alicloud_pvtz_endpoint.default.id
  rule_name   = var.name
  type        = "OUTBOUND"
  zone_name   = var.name
  forward_ips {
    ip   = "114.114.114.114"
    port = 8080
  }
}

`, name, defaultRegionToTest)
}
