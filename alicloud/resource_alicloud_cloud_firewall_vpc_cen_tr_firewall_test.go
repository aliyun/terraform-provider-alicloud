package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall VpcCenTrFirewall. >>> Resource test cases, automatically generated.
var AlicloudCloudFirewallVpcCenTrFirewallMap3609 = map[string]string{
	"status": CHECKSET,
}

func AlicloudCloudFirewallVpcCenTrFirewallBasicDependence3609(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "description" {
  default = "Created by Terraform"
}

variable "firewall_name" {
  default = "tf-test"
}

variable "tr_attachment_master_cidr" {
  default = "192.168.3.192/26"
}

variable "firewall_subnet_cidr" {
  default = "192.168.3.0/25"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "tr_attachment_slave_cidr" {
  default = "192.168.3.128/26"
}

variable "firewall_vpc_cidr" {
  default = "192.168.3.0/24"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "firewall_name_update" {
  default = "tf-test-1"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_cen_instance" "cen" {
  description       = "terraform test"
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "tr" {
  transit_router_name        = var.name
  transit_router_description = "tr-created-by-terraform"
  cen_id                     = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  description = "created by terraform"
  cidr_block  = "192.168.1.0/24"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vpc1vsw1" {
  cidr_block   = "192.168.1.0/25"
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.vpc1.id
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id       = alicloud_vpc.vpc1.id
  cidr_block   = "192.168.1.128/26"
  vswitch_name = var.name
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.vpc1.id
  route_table_name = var.name
  description      = var.name
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
  }
  zone_mappings {
    zone_id    = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
  }
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_instance.cen.id
  transit_router_id = alicloud_cen_transit_router.tr.transit_router_id
  depends_on = [alicloud_route_table.foo]
}



`, name)
}

// Case VpcCenTrFirewall全生命周期测试_副本1689148389922 3609  raw
func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3609_raw(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcCenTrFirewallMap3609)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallvpccentrfirewall%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallBasicDependence3609)
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
					"firewall_description":      "VpcCenTrFirewall created by terraform",
					"region_no":                 "${var.region}",
					"route_mode":                "managed",
					"cen_id":                    "${alicloud_cen_transit_router_vpc_attachment.tr-vpc1.cen_id}",
					"firewall_vpc_cidr":         "${var.firewall_vpc_cidr}",
					"transit_router_id":         "${alicloud_cen_transit_router.tr.transit_router_id}",
					"tr_attachment_master_cidr": "${var.tr_attachment_master_cidr}",
					"firewall_name":             name,
					"firewall_subnet_cidr":      "${var.firewall_subnet_cidr}",
					"tr_attachment_slave_cidr":  "${var.tr_attachment_slave_cidr}",
					"tr_attachment_master_zone": "${data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]}",
					"tr_attachment_slave_zone":  "${data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"firewall_description":      "VpcCenTrFirewall created by terraform",
						"region_no":                 CHECKSET,
						"route_mode":                "managed",
						"cen_id":                    CHECKSET,
						"firewall_vpc_cidr":         CHECKSET,
						"transit_router_id":         CHECKSET,
						"tr_attachment_master_cidr": CHECKSET,
						"firewall_name":             CHECKSET,
						"firewall_subnet_cidr":      CHECKSET,
						"tr_attachment_slave_cidr":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_name": "${var.firewall_name_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"firewall_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tr_attachment_master_zone", "tr_attachment_slave_zone"},
			},
		},
	})
}

// Test CloudFirewall VpcCenTrFirewall. <<< Resource test cases, automatically generated.
