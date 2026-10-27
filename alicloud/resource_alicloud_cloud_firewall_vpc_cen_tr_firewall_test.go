package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudFirewall VpcCenTrFirewall. >>> Resource test cases, automatically generated.
// Case VpcCenTrFirewall全生命周期测试_副本1689148389922 3609
func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3609(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallVpcCenTrFirewallMap3609)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallVpcCenTrFirewallBasicDependence3609)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                    "${alicloud_cen_transit_router.default.cen_id}",
					"firewall_name":             name,
					"region_no":                 "${data.alicloud_regions.default.ids.0}",
					"route_mode":                "managed",
					"transit_router_id":         "${alicloud_cen_transit_router.default.transit_router_id}",
					"firewall_vpc_cidr":         "192.168.3.0/24",
					"tr_attachment_master_cidr": "192.168.3.192/26",
					"firewall_subnet_cidr":      "192.168.3.0/25",
					"tr_attachment_slave_cidr":  "192.168.3.128/26",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                    CHECKSET,
						"firewall_name":             name,
						"region_no":                 CHECKSET,
						"route_mode":                "managed",
						"transit_router_id":         CHECKSET,
						"firewall_vpc_cidr":         "192.168.3.0/24",
						"tr_attachment_master_cidr": "192.168.3.192/26",
						"firewall_subnet_cidr":      "192.168.3.0/25",
						"tr_attachment_slave_cidr":  "192.168.3.128/26",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"firewall_name": name + "update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3609_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallVpcCenTrFirewallMap3609)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallVpcCenTrFirewallBasicDependence3609)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                    "${alicloud_cen_transit_router.default.cen_id}",
					"firewall_name":             name,
					"region_no":                 "${data.alicloud_regions.default.ids.0}",
					"route_mode":                "managed",
					"transit_router_id":         "${alicloud_cen_transit_router.default.transit_router_id}",
					"firewall_description":      name,
					"firewall_vpc_cidr":         "192.168.3.0/24",
					"tr_attachment_master_cidr": "192.168.3.192/26",
					"firewall_subnet_cidr":      "192.168.3.0/25",
					"tr_attachment_slave_cidr":  "192.168.3.128/26",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                    CHECKSET,
						"firewall_name":             name,
						"region_no":                 CHECKSET,
						"route_mode":                "managed",
						"transit_router_id":         CHECKSET,
						"firewall_description":      name,
						"firewall_vpc_cidr":         "192.168.3.0/24",
						"tr_attachment_master_cidr": "192.168.3.192/26",
						"firewall_subnet_cidr":      "192.168.3.0/25",
						"tr_attachment_slave_cidr":  "192.168.3.128/26",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3610(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallVpcCenTrFirewallMap3610)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallVpcCenTrFirewallBasicDependence3610)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                    "${alicloud_cen_transit_router_vpc_attachment.default.cen_id}",
					"firewall_name":             name,
					"region_no":                 "${data.alicloud_regions.default.ids.0}",
					"route_mode":                "manual",
					"transit_router_id":         "${alicloud_cen_transit_router.default.transit_router_id}",
					"firewall_vpc_id":           "${alicloud_vswitch.master1.vpc_id}",
					"firewall_vswitch_id":       "${alicloud_vswitch.master1.id}",
					"tr_attachment_master_zone": "${alicloud_vswitch.master1.zone_id}",
					"tr_attachment_slave_zone":  "${alicloud_vswitch.master2.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                    CHECKSET,
						"firewall_name":             name,
						"region_no":                 CHECKSET,
						"route_mode":                "manual",
						"transit_router_id":         CHECKSET,
						"firewall_vpc_id":           CHECKSET,
						"firewall_vswitch_id":       CHECKSET,
						"tr_attachment_master_zone": CHECKSET,
						"tr_attachment_slave_zone":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"firewall_name": name + "update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3610_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallVpcCenTrFirewallMap3610)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallVpcCenTrFirewallBasicDependence3610)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                    "${alicloud_cen_transit_router_vpc_attachment.default.cen_id}",
					"firewall_name":             name,
					"region_no":                 "${data.alicloud_regions.default.ids.0}",
					"route_mode":                "manual",
					"transit_router_id":         "${alicloud_cen_transit_router.default.transit_router_id}",
					"firewall_description":      name,
					"firewall_vpc_id":           "${alicloud_vswitch.master1.vpc_id}",
					"firewall_vswitch_id":       "${alicloud_vswitch.master1.id}",
					"tr_attachment_master_zone": "${alicloud_vswitch.master1.zone_id}",
					"tr_attachment_slave_zone":  "${alicloud_vswitch.master2.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                    CHECKSET,
						"firewall_name":             name,
						"region_no":                 CHECKSET,
						"route_mode":                "manual",
						"transit_router_id":         CHECKSET,
						"firewall_description":      name,
						"firewall_vpc_id":           CHECKSET,
						"firewall_vswitch_id":       CHECKSET,
						"tr_attachment_master_zone": CHECKSET,
						"tr_attachment_slave_zone":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3611(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallVpcCenTrFirewallMap3609)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallVpcCenTrFirewallBasicDependence3609)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":            "${alicloud_cen_transit_router.default.cen_id}",
					"firewall_name":     name,
					"region_no":         "${data.alicloud_regions.default.ids.0}",
					"route_mode":        "semi-managed",
					"transit_router_id": "${alicloud_cen_transit_router.default.transit_router_id}",
					"firewall_vpc_cidr": "192.168.3.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":            CHECKSET,
						"firewall_name":     name,
						"region_no":         CHECKSET,
						"route_mode":        "semi-managed",
						"transit_router_id": CHECKSET,
						"firewall_vpc_cidr": "192.168.3.0/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"firewall_name": name + "update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func TestAccAliCloudCloudFirewallVpcCenTrFirewall_basic3611_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallVpcCenTrFirewallMap3609)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallVpcCenTrFirewallBasicDependence3609)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":               "${alicloud_cen_transit_router.default.cen_id}",
					"firewall_name":        name,
					"region_no":            "${data.alicloud_regions.default.ids.0}",
					"route_mode":           "semi-managed",
					"transit_router_id":    "${alicloud_cen_transit_router.default.transit_router_id}",
					"firewall_description": name,
					"firewall_vpc_cidr":    "192.168.3.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":               CHECKSET,
						"firewall_name":        name,
						"region_no":            CHECKSET,
						"route_mode":           "semi-managed",
						"transit_router_id":    CHECKSET,
						"firewall_description": name,
						"firewall_vpc_cidr":    "192.168.3.0/24",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

var AliCloudCloudFirewallVpcCenTrFirewallMap3609 = map[string]string{
	"firewall_eni_vpc_id":        CHECKSET,
	"firewall_eni_id":            CHECKSET,
	"status":                     CHECKSET,
	"firewall_vpc_attachment_id": CHECKSET,
}

var AliCloudCloudFirewallVpcCenTrFirewallMap3610 = map[string]string{
	"firewall_eni_vpc_id": CHECKSET,
	"firewall_eni_id":     CHECKSET,
	"status":              CHECKSET,
}

func AliCloudCloudFirewallVpcCenTrFirewallBasicDependence3609(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_cen_instance" "default" {
  description       = "terraform test"
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  support_multicast          = false
  transit_router_name        = var.name
  transit_router_description = "tr-created-by-terraform"
  cen_id                     = alicloud_cen_instance.default.id
}

resource "alicloud_vpc" "master" {
 description = "created by terraform"
 cidr_block  = "192.168.1.0/24"
 vpc_name    = "${var.name}1"
}

resource "alicloud_vswitch" "master1" {
 cidr_block   = "192.168.1.0/25"
 vswitch_name = "${var.name}11"
 vpc_id       = alicloud_vpc.master.id
 zone_id      = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "master2" {
 vpc_id       = alicloud_vpc.master.id
 cidr_block   = "192.168.1.128/26"
 vswitch_name = "${var.name}12"
 zone_id      = "cn-hangzhou-i"
}

resource "alicloud_vpc" "slave" {
description = "created by terraform"
cidr_block  = "192.168.2.0/24"
vpc_name    = "${var.name}2"
}

resource "alicloud_vswitch" "slave1" {
cidr_block   = "192.168.2.0/25"
vswitch_name = "${var.name}21"
vpc_id       = alicloud_vpc.slave.id
zone_id      = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "slave2" {
cidr_block   = "192.168.2.128/26"
vswitch_name = "${var.name}22"
vpc_id       = alicloud_vpc.slave.id
zone_id      = "cn-hangzhou-i"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr_master" {
 auto_publish_route_enabled = false
 zone_mappings {
   vswitch_id = alicloud_vswitch.master1.id
   zone_id    = alicloud_vswitch.master1.zone_id
 }
 zone_mappings {
   zone_id    = alicloud_vswitch.master2.zone_id
   vswitch_id = alicloud_vswitch.master2.id
 }
 vpc_id = alicloud_vpc.master.id
 cen_id = alicloud_cen_transit_router.default.cen_id
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr_slave" {
auto_publish_route_enabled = false
vpc_id                     = alicloud_vpc.slave.id
zone_mappings {
  vswitch_id = alicloud_vswitch.slave1.id
  zone_id    = alicloud_vswitch.slave1.zone_id
}
zone_mappings {
  vswitch_id = alicloud_vswitch.slave2.id
  zone_id    = alicloud_vswitch.slave2.zone_id
}
cen_id = alicloud_cen_transit_router.default.cen_id
}
`, name)
}

func AliCloudCloudFirewallVpcCenTrFirewallBasicDependence3610(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_cen_instance" "default" {
  description       = "terraform test"
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  support_multicast          = false
  transit_router_name        = var.name
  transit_router_description = "tr-created-by-terraform"
  cen_id                     = alicloud_cen_instance.default.id
}

resource "alicloud_vpc" "master" {
 description = "created by terraform"
 cidr_block  = "192.168.1.0/24"
 vpc_name    = "${var.name}1"
}

resource "alicloud_vswitch" "master1" {
 cidr_block   = "192.168.1.0/25"
 vswitch_name = "${var.name}11"
 vpc_id       = alicloud_vpc.master.id
 zone_id      = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "master2" {
 vpc_id       = alicloud_vpc.master.id
 cidr_block   = "192.168.1.128/26"
 vswitch_name = "${var.name}12"
 zone_id      = "cn-hangzhou-i"
}

resource "alicloud_route_table" "default" {
  vpc_id           = alicloud_vpc.master.id
}

resource "alicloud_route_table_attachment" "default" {
  vswitch_id     = alicloud_vswitch.master1.id
  route_table_id = alicloud_route_table.default.id
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
 auto_publish_route_enabled = false
 zone_mappings {
   vswitch_id = alicloud_route_table_attachment.default.vswitch_id
   zone_id    = alicloud_vswitch.master1.zone_id
 }
 zone_mappings {
   zone_id    = alicloud_vswitch.master2.zone_id
   vswitch_id = alicloud_vswitch.master2.id
 }
 vpc_id = alicloud_vpc.master.id
 cen_id = alicloud_cen_transit_router.default.cen_id
}
`, name)
}

// Test CloudFirewall VpcCenTrFirewall. <<< Resource test cases, automatically generated.
