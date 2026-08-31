package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall VpcCenTrFirewall. >>> Resource test cases, automatically generated.
// Case VpcCenTrFirewall全生命周期测试_副本1689148389922 3609
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
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallBasicDependence3609)
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

var AlicloudCloudFirewallVpcCenTrFirewallMap3609 = map[string]string{
	"firewall_eni_vpc_id":        CHECKSET,
	"firewall_eni_id":            CHECKSET,
	"status":                     CHECKSET,
	"firewall_vpc_attachment_id": CHECKSET,
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

resource "alicloud_cen_instance" "cen" {
  description       = "terraform test"
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "tr" {
  support_multicast          = false
  transit_router_name        = var.name
  transit_router_description = "tr-created-by-terraform"
  cen_id                     = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  description = "created by terraform"
  cidr_block  = "192.168.1.0/24"
  vpc_name    = "${var.name}1"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  cidr_block   = "192.168.1.0/25"
  vswitch_name = "${var.name}11"
  vpc_id       = alicloud_vpc.vpc1.id
  zone_id      = var.zone1
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id       = alicloud_vpc.vpc1.id
  cidr_block   = "192.168.1.128/26"
  vswitch_name = "${var.name}12"
  zone_id      = var.zone2
}

resource "alicloud_vpc" "vpc2" {
  description = "created by terraform"
  cidr_block  = "192.168.2.0/24"
  vpc_name    = "${var.name}2"
}

resource "alicloud_vswitch" "vpc2vsw1" {
  cidr_block   = "192.168.2.0/25"
  vswitch_name = "${var.name}21"
  vpc_id       = alicloud_vpc.vpc2.id
  zone_id      = var.zone1
}

resource "alicloud_vswitch" "vpc2vsw2" {
  cidr_block   = "192.168.2.128/26"
  vswitch_name = "${var.name}22"
  vpc_id       = alicloud_vpc.vpc2.id
  zone_id      = var.zone2
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  auto_publish_route_enabled = false
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
  }
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  auto_publish_route_enabled = false
  vpc_id                     = alicloud_vpc.vpc2.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = alicloud_vswitch.vpc2vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = alicloud_vswitch.vpc2vsw2.zone_id
  }
  cen_id = alicloud_cen_instance.cen.id
}`, name)
}

// Case VpcCenTrFirewall全生命周期测试-001 7158
func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic7158(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcCenTrFirewallMap7158)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallBasicDependence7158)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
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
					"cen_id":                    "${alicloud_cen_instance.cen.id}",
					"firewall_vpc_cidr":         "${var.firewall_vpc_cidr}",
					"transit_router_id":         "${alicloud_cen_transit_router.tr.transit_router_id}",
					"tr_attachment_master_cidr": "${var.tr_attachment_master_cidr}",
					"firewall_name":             "${var.firewall_name}",
					"firewall_subnet_cidr":      "${var.firewall_subnet_cidr}",
					"tr_attachment_slave_cidr":  "${var.tr_attachment_slave_cidr}",
					"tr_attachment_master_zone": "${var.zone1}",
					"tr_attachment_slave_zone":  "${var.zone2}",
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
						"tr_attachment_master_zone": CHECKSET,
						"tr_attachment_slave_zone":  CHECKSET,
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

var AlicloudCloudFirewallVpcCenTrFirewallMap7158 = map[string]string{
	"firewall_eni_vpc_id":        CHECKSET,
	"firewall_eni_id":            CHECKSET,
	"status":                     CHECKSET,
	"firewall_vpc_attachment_id": CHECKSET,
}

func AlicloudCloudFirewallVpcCenTrFirewallBasicDependence7158(name string) string {
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
  default = "10.0.2.0/24"
}

variable "firewall_subnet_cidr" {
  default = "10.0.1.0/24"
}

variable "region" {
  default = "cn-shenzhen"
}

variable "tr_attachment_slave_cidr" {
  default = "10.0.3.0/24"
}

variable "firewall_vpc_cidr" {
  default = "10.0.0.0/16"
}

variable "zone1" {
  default = "cn-shenzhen-d"
}

variable "firewall_name_update" {
  default = "tf-test-1"
}

variable "zone2" {
  default = "cn-shenzhen-e"
}

resource "alicloud_cen_instance" "cen" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "tr" {
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}01"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  vpc_id     = alicloud_vpc.vpc1.id
  zone_id    = var.zone1
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id     = alicloud_vpc.vpc1.id
  zone_id    = var.zone2
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_vpc" "vpc2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}02"
}

resource "alicloud_vswitch" "vpc2vsw1" {
  vpc_id     = alicloud_vpc.vpc2.id
  zone_id    = var.zone1
  cidr_block = "172.16.3.0/24"
}

resource "alicloud_vswitch" "vpc2vsw2" {
  vpc_id     = alicloud_vpc.vpc2.id
  zone_id    = var.zone2
  cidr_block = "172.16.4.0/24"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  auto_publish_route_enabled = false
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
  }
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  auto_publish_route_enabled = false
  vpc_id                     = alicloud_vpc.vpc2.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = alicloud_vswitch.vpc2vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = alicloud_vswitch.vpc2vsw2.zone_id
  }
  cen_id = alicloud_cen_instance.cen.id
}`, name)
}

// Test CloudFirewall VpcCenTrFirewall. <<< Resource test cases, automatically generated.
