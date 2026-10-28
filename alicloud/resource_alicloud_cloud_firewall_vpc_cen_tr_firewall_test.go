package alicloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				// Re-plan the same managed config: on the fixed (Computed) schema the
				// server-assigned firewall_vpc_id/firewall_vswitch_id must NOT produce a
				// removal diff, so the plan is empty (no replacement). On the buggy
				// (non-Computed) schema the Read-backfill yields a removal diff that
				// ForceNew turns into destroy+recreate, making the plan non-empty.
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
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// Case VpcCenTrFirewall manual route_mode
func TestAccAliCloudCloudFirewallVpcCenTrFirewall_manual(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcCenTrFirewallMapManual)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallManualDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_description":      "VpcCenTrFirewall manual created by terraform",
					"region_no":                 "${var.region}",
					"route_mode":                "manual",
					"cen_id":                    "${alicloud_cen_instance.cen.id}",
					"firewall_vpc_id":           "${alicloud_vpc.vpc1.id}",
					"firewall_vswitch_id":       "${alicloud_vswitch.vpc1vsw1.id}",
					"transit_router_id":         "${alicloud_cen_transit_router.tr.transit_router_id}",
					"firewall_name":             "${var.firewall_name}",
					"tr_attachment_master_zone": "${alicloud_vswitch.vpc1vsw1.zone_id}",
					"tr_attachment_slave_zone":  "${alicloud_vswitch.vpc1vsw2.zone_id}",
					"depends_on": []string{
						"alicloud_cen_transit_router_vpc_attachment.tr-vpc1",
						"alicloud_route_table_attachment.rt-att",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"firewall_description":      "VpcCenTrFirewall manual created by terraform",
						"region_no":                 CHECKSET,
						"route_mode":                "manual",
						"cen_id":                    CHECKSET,
						"firewall_vpc_id":           CHECKSET,
						"firewall_vswitch_id":       CHECKSET,
						"transit_router_id":         CHECKSET,
						"firewall_name":             CHECKSET,
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

// firewall_vpc_attachment_id (TrAttachmentId) is only returned by
// DescribeTrFirewallsV2Detail in automatic (managed) mode; in manual mode the
// firewall VPC's TR attachment is user-managed, so the field is empty and is
// not asserted here.
var AlicloudCloudFirewallVpcCenTrFirewallMapManual = map[string]string{
	"firewall_eni_vpc_id": CHECKSET,
	"firewall_eni_id":     CHECKSET,
	"status":              CHECKSET,
	"firewall_vpc_id":     CHECKSET,
	"firewall_vswitch_id": CHECKSET,
}

// AlicloudCloudFirewallVpcCenTrFirewallManualDependence builds the manual-mode
// fixture by reusing the cn-hangzhou BasicDependence3609 setup (whose TR-valid
// zones cn-hangzhou-h/i are shared with the managed/semi-managed tests) and
// appending a custom route table + attachment so the user VPC supports subnet
// routing, which manual mode requires. The firewall's tr_attachment_master_zone
// and tr_attachment_slave_zone are sourced directly from the vSwitch zone_id
// attributes (vpc1vsw1/vpc1vsw2) that BasicDependence3609 also feeds into the TR
// VPC attachment's zone_mappings, so the firewall zones are always the same
// zones the TR attachment is deployed in — keeping TR/vSwitch/Firewall
// consistent. zone_mappings is a TypeSet and cannot be indexed by position in
// HCL, so the vSwitch attributes are referenced instead. No data source returns
// TR-valid zones in Cloud Firewall supported regions: eu-central-1, where
// data.alicloud_zones.ids.0 is TR-valid, is excluded from
// CloudFirewallVpcFirewallCenSupportRegions, so the vSwitch zones reuse the
// proven managed fixture instead.
func AlicloudCloudFirewallVpcCenTrFirewallManualDependence(name string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_route_table" "rt" {
  vpc_id = alicloud_vpc.vpc1.id
  name   = "${var.name}-rt"
}

resource "alicloud_route_table_attachment" "rt-att" {
  route_table_id = alicloud_route_table.rt.id
  vswitch_id     = alicloud_vswitch.vpc1vsw1.id
}`, AlicloudCloudFirewallVpcCenTrFirewallBasicDependence3609(name))
}

// Case VpcCenTrFirewall semi-managed route_mode (手动模式托管VPC). In
// semi-managed mode Cloud Firewall still creates and manages the firewall
// VPC, so only firewall_vpc_cidr is required; firewall_vpc_id /
// firewall_vswitch_id are server-assigned and backfilled (Computed), exactly
// like managed mode. Reuses the cn-hangzhou BasicDependence3609 setup.
func TestAccAliCloudCloudFirewallVpcCenTrFirewall_semimanaged(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcCenTrFirewallMapSemiManaged)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewall")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallBasicDependence3609)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_description":      "VpcCenTrFirewall semi-managed created by terraform",
					"region_no":                 "${var.region}",
					"route_mode":                "semi-managed",
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
						"firewall_description":      "VpcCenTrFirewall semi-managed created by terraform",
						"region_no":                 CHECKSET,
						"route_mode":                "semi-managed",
						"cen_id":                    CHECKSET,
						"firewall_vpc_cidr":         CHECKSET,
						"firewall_vpc_id":           CHECKSET,
						"firewall_vswitch_id":       CHECKSET,
						"transit_router_id":         CHECKSET,
						"tr_attachment_master_cidr": CHECKSET,
						"firewall_name":             CHECKSET,
						"firewall_subnet_cidr":      CHECKSET,
						"tr_attachment_slave_cidr":  CHECKSET,
					}),
				),
			},
			{
				// Re-plan the same semi-managed config: the server-assigned
				// firewall_vpc_id / firewall_vswitch_id (Computed) must NOT
				// produce a removal diff, so the plan is empty (no
				// replacement). Mirrors the managed-mode guard in basic3609.
				Config: testAccConfig(map[string]interface{}{
					"firewall_description":      "VpcCenTrFirewall semi-managed created by terraform",
					"region_no":                 "${var.region}",
					"route_mode":                "semi-managed",
					"cen_id":                    "${alicloud_cen_instance.cen.id}",
					"firewall_vpc_cidr":         "${var.firewall_vpc_cidr}",
					"transit_router_id":         "${alicloud_cen_transit_router.tr.transit_router_id}",
					"tr_attachment_master_cidr": "${var.tr_attachment_master_cidr}",
					"firewall_name":             "${var.firewall_name}",
					"firewall_subnet_cidr":      "${var.firewall_subnet_cidr}",
					"tr_attachment_slave_cidr":  "${var.tr_attachment_slave_cidr}",
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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

// semi-managed (手动模式托管VPC) behaves like managed for Read: Cloud
// Firewall creates and manages the firewall VPC, so firewall_vpc_id /
// firewall_vswitch_id are server-assigned and backfilled (Computed), and
// firewall_vpc_attachment_id (TrAttachmentId) is returned.
var AlicloudCloudFirewallVpcCenTrFirewallMapSemiManaged = map[string]string{
	"firewall_eni_vpc_id":        CHECKSET,
	"firewall_eni_id":            CHECKSET,
	"status":                     CHECKSET,
	"firewall_vpc_id":            CHECKSET,
	"firewall_vswitch_id":        CHECKSET,
	"firewall_vpc_attachment_id": CHECKSET,
}

// Test CloudFirewall VpcCenTrFirewall. <<< Resource test cases, automatically generated.

// TestUnitCloudFirewallVpcCenTrFirewallManagedOmitComputedGuard is a cloud-free
// regression guard for the managed-mode ForceNew anti-pattern on
// firewall_vpc_id / firewall_vswitch_id.
//
// In managed (automatic) mode the user omits these fields; Cloud Firewall
// auto-assigns the firewall ENI VPC/vSwitch and DescribeTrFirewallsV2Detail
// returns them as FirewallEniVpcId / FirewallEniVswitchId. Read() backfills them
// into state unconditionally. With Optional+ForceNew but WITHOUT Computed, the
// backfilled value yields a removal diff on every plan and ForceNew turns it
// into a destroy+recreate — a silent destructive change for existing managed
// resources. Marking them Computed lets the SDK suppress the diff when the
// configuration omits them, matching the Computed modeling of
// firewall_eni_vpc_id / firewall_vpc_attachment_id.
//
// On the buggy head (no Computed) this guard fails; on the fixed head it passes.
func TestUnitCloudFirewallVpcCenTrFirewallManagedOmitComputedGuard(t *testing.T) {
	schemaMap := Provider().ResourcesMap["alicloud_cloud_firewall_vpc_cen_tr_firewall"].Schema
	for _, field := range []string{"firewall_vpc_id", "firewall_vswitch_id"} {
		s, ok := schemaMap[field]
		if !ok {
			t.Fatalf("schema missing field %q", field)
		}
		if !s.Optional {
			t.Errorf("%q must be Optional (manual route_mode sets it in config)", field)
		}
		if !s.ForceNew {
			t.Errorf("%q must be ForceNew (ENI VPC/vSwitch cannot change after create)", field)
		}
		if !s.Computed {
			t.Errorf("%q must be Computed to suppress the managed-mode Read-backfill diff; without Computed the backfilled value produces a removal diff that ForceNew turns into destroy+recreate", field)
		}
	}
}

// TestUnitCloudFirewallVpcCenTrFirewallManagedOmitNoForceNew is a cloud-free
// behavioral reproduction of the managed-mode ForceNew anti-pattern. It drives
// the SDK's real Resource.Diff over a managed-mode configuration that omits
// firewall_vpc_id / firewall_vswitch_id against a prior state in which Read()
// has backfilled the server-assigned values (FirewallEniVpcId /
// FirewallEniVswitchId). With Optional+ForceNew and no Computed, the omitted
// config produces a removal diff that ForceNew turns into destroy+recreate, so
// diff is non-nil and diff.RequiresNew() is true. With Computed the SDK
// suppresses that diff entirely — a nil diff means no change and no
// replacement, which is the correct fixed behavior. On the buggy head (no
// Computed) this test FAILS (RequiresNew); on the fixed head (Computed) it
// PASSES (nil or no-RequiresNew diff) — the requested positive evidence that
// the bug is real and the fix removes the replacement.
func TestUnitCloudFirewallVpcCenTrFirewallManagedOmitNoForceNew(t *testing.T) {
	r := Provider().ResourcesMap["alicloud_cloud_firewall_vpc_cen_tr_firewall"]
	if r == nil {
		t.Fatal("resource alicloud_cloud_firewall_vpc_cen_tr_firewall not registered")
	}
	// Prior state after Read(): managed mode, user omitted firewall_vpc_id and
	// firewall_vswitch_id, server auto-assigned them and Read() backfilled.
	state := &terraform.InstanceState{
		ID: "fw-test-id",
		Attributes: map[string]string{
			"id":                         "fw-test-id",
			"cen_id":                     "cen-xxx",
			"firewall_name":              "tf-test",
			"firewall_vpc_cidr":          "192.168.3.0/24",
			"firewall_subnet_cidr":       "192.168.3.0/25",
			"region_no":                  "cn-hangzhou",
			"route_mode":                 "managed",
			"transit_router_id":          "tr-xxx",
			"tr_attachment_master_cidr":  "192.168.3.192/26",
			"tr_attachment_slave_cidr":   "192.168.3.128/26",
			"firewall_vpc_id":            "vpc-server-assigned",
			"firewall_vswitch_id":        "vsw-server-assigned",
			"firewall_eni_vpc_id":        "vpc-server-assigned",
			"firewall_eni_id":            "eni-xxx",
			"firewall_vpc_attachment_id": "tr-att-xxx",
			"status":                     "Ready",
		},
	}
	// Managed-mode config deliberately omits firewall_vpc_id and
	// firewall_vswitch_id (server-assigned, not user-supplied).
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"cen_id":                    "cen-xxx",
		"firewall_name":             "tf-test",
		"firewall_vpc_cidr":         "192.168.3.0/24",
		"firewall_subnet_cidr":      "192.168.3.0/25",
		"region_no":                 "cn-hangzhou",
		"route_mode":                "managed",
		"transit_router_id":         "tr-xxx",
		"tr_attachment_master_cidr": "192.168.3.192/26",
		"tr_attachment_slave_cidr":  "192.168.3.128/26",
	})
	diff, err := r.Diff(context.Background(), state, config, nil)
	if err != nil {
		t.Fatalf("Resource.Diff returned error: %v", err)
	}
	// On the fixed (Computed) schema, when the managed-mode config omits
	// firewall_vpc_id/firewall_vswitch_id and state holds server-assigned
	// values, the SDK suppresses the diff entirely — a nil diff means no
	// change and no replacement, which is the correct fixed behavior. On
	// the buggy (non-Computed) schema the backfilled values yield a
	// removal diff that ForceNew turns into destroy+recreate, so diff is
	// non-nil and diff.RequiresNew() is true. Either a nil diff or a
	// non-nil diff without RequiresNew is a PASS; only RequiresNew proves
	// the bug.
	if diff != nil && diff.RequiresNew() {
		t.Errorf("managed-mode omit of firewall_vpc_id/firewall_vswitch_id must NOT require replacement; diff.RequiresNew()=true means the Read-backfilled server-assigned values yield a removal diff that ForceNew turns into destroy+recreate. Fix: mark firewall_vpc_id/firewall_vswitch_id as Computed so the SDK suppresses the diff when config omits them and state holds a server-assigned value.")
	}
}

// TestUnitCloudFirewallVpcCenTrFirewallSemiManagedOmitNoForceNew is the
// semi-managed equivalent of TestUnitCloudFirewallVpcCenTrFirewallManagedOmitNoForceNew.
// It drives the SDK's real Resource.Diff over a semi-managed configuration
// (route_mode="semi-managed") that omits firewall_vpc_id / firewall_vswitch_id
// against a prior state in which Read() has backfilled the server-assigned
// values. semi-managed (手动模式托管VPC) behaves like managed: Cloud Firewall
// creates and manages the firewall VPC, so these fields are server-assigned and
// Computed. With Computed the SDK suppresses the omission diff (nil/no-RequiresNew
// diff = PASS); without Computed the backfilled values would yield a removal diff
// that ForceNew turns into destroy+recreate (RequiresNew = FAIL). It also confirms
// route_mode="semi-managed" is accepted by the schema ValidateFunc, so the
// enum-acceptance and no-replacement behavior are both proven cloud-free.
func TestUnitCloudFirewallVpcCenTrFirewallSemiManagedOmitNoForceNew(t *testing.T) {
	r := Provider().ResourcesMap["alicloud_cloud_firewall_vpc_cen_tr_firewall"]
	if r == nil {
		t.Fatal("resource alicloud_cloud_firewall_vpc_cen_tr_firewall not registered")
	}
	// Prior state after Read(): semi-managed mode, user omitted firewall_vpc_id
	// and firewall_vswitch_id, server auto-assigned them and Read() backfilled.
	state := &terraform.InstanceState{
		ID: "fw-test-id",
		Attributes: map[string]string{
			"id":                         "fw-test-id",
			"cen_id":                     "cen-xxx",
			"firewall_name":              "tf-test",
			"firewall_vpc_cidr":          "192.168.3.0/24",
			"firewall_subnet_cidr":       "192.168.3.0/25",
			"region_no":                  "cn-hangzhou",
			"route_mode":                 "semi-managed",
			"transit_router_id":          "tr-xxx",
			"tr_attachment_master_cidr":  "192.168.3.192/26",
			"tr_attachment_slave_cidr":   "192.168.3.128/26",
			"firewall_vpc_id":            "vpc-server-assigned",
			"firewall_vswitch_id":        "vsw-server-assigned",
			"firewall_eni_vpc_id":        "vpc-server-assigned",
			"firewall_eni_id":            "eni-xxx",
			"firewall_vpc_attachment_id": "tr-att-xxx",
			"status":                     "Ready",
		},
	}
	// semi-managed-mode config deliberately omits firewall_vpc_id and
	// firewall_vswitch_id (server-assigned, not user-supplied).
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"cen_id":                    "cen-xxx",
		"firewall_name":             "tf-test",
		"firewall_vpc_cidr":         "192.168.3.0/24",
		"firewall_subnet_cidr":      "192.168.3.0/25",
		"region_no":                 "cn-hangzhou",
		"route_mode":                "semi-managed",
		"transit_router_id":         "tr-xxx",
		"tr_attachment_master_cidr": "192.168.3.192/26",
		"tr_attachment_slave_cidr":  "192.168.3.128/26",
	})
	diff, err := r.Diff(context.Background(), state, config, nil)
	if err != nil {
		t.Fatalf("Resource.Diff returned error: %v", err)
	}
	if diff != nil && diff.RequiresNew() {
		t.Errorf("semi-managed-mode omit of firewall_vpc_id/firewall_vswitch_id must NOT require replacement; diff.RequiresNew()=true means the Read-backfilled server-assigned values yield a removal diff that ForceNew turns into destroy+recreate. Fix: mark firewall_vpc_id/firewall_vswitch_id as Computed so the SDK suppresses the diff when config omits them and state holds a server-assigned value.")
	}
}

// TestUnitCloudFirewallVpcCenTrFirewallRouteModeContract is a cloud-free guard
// for the Create-time route_mode required-field contract (see
// validateTrFirewallRouteModeContract). It proves the validation rejects a manual
// configuration that omits the four required manual fields with a precise error,
// and accepts a fully-specified manual config and the managed/semi-managed
// firewall_vpc_cidr-required configs. This is the non-billable negative example
// for the manual zone-required check: without it, Cloud Firewall only returns a
// generic ErrorParameters (400) after the API call.
func TestUnitCloudFirewallVpcCenTrFirewallRouteModeContract(t *testing.T) {
	// manual missing every required field -> error
	if err := validateTrFirewallRouteModeContract("manual", "", "", "", "", ""); err == nil {
		t.Errorf("manual route_mode with empty firewall_vpc_id/firewall_vswitch_id/tr_attachment_master_zone/tr_attachment_slave_zone must fail validation")
	}
	// manual with vpc/vswitch but missing the zones -> error on the zone fields
	err := validateTrFirewallRouteModeContract("manual", "vpc-1", "vsw-1", "", "", "")
	if err == nil {
		t.Errorf("manual route_mode with empty tr_attachment_master_zone/tr_attachment_slave_zone must fail validation")
	}
	// manual fully specified -> no error
	if err := validateTrFirewallRouteModeContract("manual", "vpc-1", "vsw-1", "cn-hangzhou-h", "cn-hangzhou-i", ""); err != nil {
		t.Errorf("manual route_mode with all required fields must pass validation, got: %v", err)
	}
	// managed missing firewall_vpc_cidr -> error
	if err := validateTrFirewallRouteModeContract("managed", "", "", "", "", ""); err == nil {
		t.Errorf("managed route_mode with empty firewall_vpc_cidr must fail validation")
	}
	// semi-managed missing firewall_vpc_cidr -> error
	if err := validateTrFirewallRouteModeContract("semi-managed", "", "", "", "", ""); err == nil {
		t.Errorf("semi-managed route_mode with empty firewall_vpc_cidr must fail validation")
	}
	// managed with firewall_vpc_cidr -> no error
	if err := validateTrFirewallRouteModeContract("managed", "", "", "", "", "192.168.3.0/24"); err != nil {
		t.Errorf("managed route_mode with firewall_vpc_cidr must pass validation, got: %v", err)
	}
	// semi-managed with firewall_vpc_cidr -> no error
	if err := validateTrFirewallRouteModeContract("semi-managed", "", "", "", "", "192.168.3.0/24"); err != nil {
		t.Errorf("semi-managed route_mode with firewall_vpc_cidr must pass validation, got: %v", err)
	}
}
