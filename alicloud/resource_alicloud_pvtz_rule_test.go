package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPrivateZoneRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pvtz_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudPrivateZoneRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePvtzRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPrivateZoneRuleBasicDependence0)
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
					"zone_name":   name,
					"rule_name":   name,
					"type":        "OUTBOUND",
					"endpoint_id": "${alicloud_pvtz_endpoint.default.id}",
					"forward_ips": []map[string]interface{}{
						{
							"ip":   "114.114.114.114",
							"port": "8081",
						},
						{
							"ip":   "114.114.114.115",
							"port": "8087",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_name":     name,
						"rule_name":     name,
						"type":          "OUTBOUND",
						"endpoint_id":   CHECKSET,
						"forward_ips.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forward_ips": []map[string]interface{}{
						{
							"ip":   "114.114.114.114",
							"port": "8085",
						},
						{
							"ip":   "112.114.114.115",
							"port": "8089",
						},
						{
							"ip":   "119.114.114.115",
							"port": "8012",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forward_ips.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": name,
					"forward_ips": []map[string]interface{}{
						{
							"ip":   "114.114.114.114",
							"port": "8081",
						},
						{
							"ip":   "114.114.114.115",
							"port": "8087",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":     name,
						"forward_ips.#": "2",
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

var AlicloudPrivateZoneRuleMap0 = map[string]string{}

func AlicloudPrivateZoneRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}


data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  count      = 2
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, count.index)
  zone_id    = data.alicloud_pvtz_resolver_zones.default.zones[count.index].zone_id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
  name   = var.name
}

resource "alicloud_pvtz_endpoint" "default" {
  endpoint_name     = var.name
  security_group_id = alicloud_security_group.default.id
  vpc_id            = alicloud_vpc.default.id
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

`, name, defaultRegionToTest)
}
