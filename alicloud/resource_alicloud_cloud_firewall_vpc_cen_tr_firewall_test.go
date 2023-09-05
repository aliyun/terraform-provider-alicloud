package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall VpcCenTrFirewall. >>> Resource test cases, automatically generated.
// Case 3609
func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3609(t *testing.T) {
	var v map[string]interface{}
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
			testAccPreCheckWithRegions(t, true, connectivity.CENVpcCenTrFirewallSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_name":             "${var.firewall_name}",
					"region_no":                 "${var.region}",
					"route_mode":                "managed",
					"cen_id":                    "${alicloud_cen_transit_router_vpc_attachment.tr-vpc2.cen_id}",
					"transit_router_id":         "${alicloud_cen_transit_router.tr.transit_router_id}",
					"tr_attachment_slave_cidr":  "${var.tr_attachment_slave_cidr}",
					"firewall_vpc_cidr":         "${var.firewall_vpc_cidr}",
					"tr_attachment_master_cidr": "${var.tr_attachment_master_cidr}",
					"firewall_subnet_cidr":      "${var.firewall_subnet_cidr}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
					"firewall_name": "${var.firewall_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"firewall_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_name": "${var.firewall_name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"firewall_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_description":      "VpcCenTrFirewall created by terraform",
					"region_no":                 "${var.region}",
					"route_mode":                "managed",
					"cen_id":                    "${alicloud_cen_instance.cen.id}",
					"firewall_vpc_cidr":         "${var.firewall_vpc_cidr}",
					"transit_router_id":         "${alicloud_cen_transit_router.tr.transit_router_id}",
					"tr_attachment_master_cidr": "${var.tr_attachment_master_cidr}",
					"firewall_name":             "${var.firewall_name}",
					"firewall_subnet_cidr":      "${var.firewall_subnet_cidr}",
					"tr_attachment_slave_cidr":  "${var.tr_attachment_slave_cidr}",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCloudFirewallVpcCenTrFirewallMap3609 = map[string]string{
	"status": CHECKSET,
}

func AlicloudCloudFirewallVpcCenTrFirewallBasicDependence3609(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "description" {
  default = "Created by Terraform"
}

variable "firewall_name" {
  default = "tf-example"
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
  default = "tf-example-1"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_cen_instance" "cen" {
  description       = "terraform example"
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "tr" {
  transit_router_name        = "${var.name}1"
  transit_router_description = "tr-created-by-terraform"
  cen_id                     = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  description = "created by terraform"
  cidr_block  = "192.168.1.0/24"
  vpc_name    = "${var.name}2"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  cidr_block   = "192.168.1.0/25"
  vswitch_name = "${var.name}3"
  vpc_id       = alicloud_vpc.vpc1.id
  zone_id      = var.zone1
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id       = alicloud_vpc.vpc1.id
  cidr_block   = "192.168.1.128/26"
  vswitch_name = "${var.name}4"
  zone_id      = var.zone2
}

resource "alicloud_vpc" "vpc2" {
  description = "created by terraform"
  cidr_block  = "192.168.2.0/24"
  vpc_name    = "${var.name}5"
}

resource "alicloud_vswitch" "vpc2vsw1" {
  cidr_block   = "192.168.2.0/25"
  vswitch_name = "${var.name}6"
  vpc_id       = alicloud_vpc.vpc2.id
  zone_id      = var.zone1
}

resource "alicloud_vswitch" "vpc2vsw2" {
  cidr_block   = "192.168.2.128/26"
  vswitch_name = "${var.name}7"
  vpc_id       = alicloud_vpc.vpc2.id
  zone_id      = var.zone2
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = var.zone1
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
    zone_id    = var.zone2
  }
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_transit_router.tr.cen_id
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  vpc_id = alicloud_vpc.vpc2.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = var.zone1
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = var.zone2
  }
  cen_id = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.cen_id
}


`, name)
}

// Case 3609  twin
func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3609_twin(t *testing.T) {
	var v map[string]interface{}
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
			testAccPreCheckWithRegions(t, true, connectivity.CENVpcCenTrFirewallSupportRegions)
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
					"cen_id":                    "${alicloud_cen_transit_router_vpc_attachment.tr-vpc2.cen_id}",
					"firewall_vpc_cidr":         "${var.firewall_vpc_cidr}",
					"transit_router_id":         "${alicloud_cen_transit_router.tr.transit_router_id}",
					"tr_attachment_master_cidr": "${var.tr_attachment_master_cidr}",
					"firewall_name":             "${var.firewall_name_update}",
					"firewall_subnet_cidr":      "${var.firewall_subnet_cidr}",
					"tr_attachment_slave_cidr":  "${var.tr_attachment_slave_cidr}",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test CloudFirewall VpcCenTrFirewall. <<< Resource test cases, automatically generated.
