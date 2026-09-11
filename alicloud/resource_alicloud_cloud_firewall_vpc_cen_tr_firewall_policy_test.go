package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAliCloudCloudFirewallVpcCenTrFirewallPolicy_basic0 — full lifecycle:
// create with fullmesh + src candidates, update should_recover/status/dest_candidate_list, import.
func TestAccAliCloudCloudFirewallVpcCenTrFirewallPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall_policy.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"status":                      CHECKSET,
		"tr_firewall_route_policy_id": CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewallPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccfwpolicy%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence0)
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
					"firewall_id":        "${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}",
					"policy_type":        "fullmesh",
					"policy_name":        name,
					"policy_description": "tf-test-policy-desc",
					"src_candidate_list": []interface{}{
						map[string]interface{}{
							"candidate_id":   "${alicloud_vpc.vpc1.id}",
							"candidate_type": "VPC",
						},
						map[string]interface{}{
							"candidate_id":   "${alicloud_vpc.vpc2.id}",
							"candidate_type": "VPC",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_type":                 "fullmesh",
						"policy_name":                 CHECKSET,
						"policy_description":          "tf-test-policy-desc",
						"firewall_id":                 CHECKSET,
						"status":                      CHECKSET,
						"tr_firewall_route_policy_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"should_recover": "true",
					"status":         "open",
					"src_candidate_list": []interface{}{
						map[string]interface{}{
							"candidate_id":   "${alicloud_express_connect_router_express_connect_router.ecr.id}",
							"candidate_type": "ECR",
						},
					},
					"dest_candidate_list": []interface{}{
						map[string]interface{}{
							"candidate_id":   "${alicloud_vpc.vpc1.id}",
							"candidate_type": "VPC",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"should_recover": "true",
						"status":         CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"should_recover", "status", "src_candidate_list"},
			},
		},
	})
}

// TestAccAliCloudCloudFirewallVpcCenTrFirewallPolicy_dataSource — datasource query by firewall_id.
func TestAccAliCloudCloudFirewallVpcCenTrFirewallPolicy_dataSource(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall_policy.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"status":                      CHECKSET,
		"tr_firewall_route_policy_id": CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewallPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccfwpolicyds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"firewall_id":        "${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}",
					"policy_type":        "fullmesh",
					"policy_name":        name,
					"policy_description": "tf-test-policy-desc",
					"src_candidate_list": []interface{}{
						map[string]interface{}{
							"candidate_id":   "${alicloud_vpc.vpc1.id}",
							"candidate_type": "VPC",
						},
						map[string]interface{}{
							"candidate_id":   "${alicloud_vpc.vpc2.id}",
							"candidate_type": "VPC",
						},
					},
				}) + fmt.Sprintf(`
data "alicloud_cloud_firewall_vpc_cen_tr_firewall_policies" "default" {
  firewall_id = alicloud_cloud_firewall_vpc_cen_tr_firewall_policy.default.firewall_id
  ids         = [alicloud_cloud_firewall_vpc_cen_tr_firewall_policy.default.id]
}
`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
					resource.TestCheckResourceAttr("data.alicloud_cloud_firewall_vpc_cen_tr_firewall_policies.default", "policies.#", "1"),
				),
			},
		},
	})
}

func AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alicloud_cen_instance" "cen" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "tr" {
  cen_id = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}-vpc1"
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
  vpc_name   = "${var.name}-vpc2"
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
  vpc_id                     = alicloud_vpc.vpc1.id
  cen_id                     = alicloud_cen_instance.cen.id
  transit_router_id          = alicloud_cen_transit_router.tr.transit_router_id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
  }
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  auto_publish_route_enabled = false
  vpc_id                     = alicloud_vpc.vpc2.id
  cen_id                     = alicloud_cen_instance.cen.id
  transit_router_id          = alicloud_cen_transit_router.tr.transit_router_id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = alicloud_vswitch.vpc2vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = alicloud_vswitch.vpc2vsw2.zone_id
  }
}

resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "default" {
  route_mode                = "managed"
  region_no                 = var.region
  firewall_description      = "VpcCenTrFirewall created by terraform"
  tr_attachment_master_zone = var.zone1
  firewall_name             = "${var.name}-fw"
  tr_attachment_master_cidr = "10.0.2.0/24"
  firewall_subnet_cidr      = "10.0.1.0/24"
  cen_id                    = alicloud_cen_instance.cen.id
  tr_attachment_slave_cidr  = "10.0.3.0/24"
  tr_attachment_slave_zone  = var.zone2
  firewall_vpc_cidr         = "10.0.0.0/16"
  transit_router_id         = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_express_connect_router_express_connect_router" "ecr" {
  ecr_name         = "${var.name}-ecr"
  alibaba_side_asn = "64514"
  description      = "tf-test-ecr"
}

resource "alicloud_cen_transit_router_ecr_attachment" "ecr-att" {
  ecr_id            = alicloud_express_connect_router_express_connect_router.ecr.id
  cen_id            = alicloud_cen_instance.cen.id
  transit_router_id = alicloud_cen_transit_router.tr.transit_router_id
}
`, name)
}
