// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudFirewall VpcCenTrFirewallPolicy. >>> Resource test cases, automatically generated.
// Case VpcCenTrFirewallPolicy-fullmesh(SrcCandidateList&gt;=3)_副本1746585232432 10781
func TestAccAliCloudCloudFirewallVpcCenTrFirewallPolicy_basic10781(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcCenTrFirewallPolicyMap10781)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewallPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence10781)
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
					"src_candidate_list": []map[string]interface{}{
						{
							"candidate_id":   "${alicloud_express_connect_router_express_connect_router.ExpressConnectRouter.id}",
							"candidate_type": "ECR",
						},
						{
							"candidate_id":   "${alicloud_vpc.vpc1.id}",
							"candidate_type": "VPC",
						},
					},
					"policy_type":        "fullmesh",
					"policy_description": "111111",
					"firewall_id":        "${alicloud_cloud_firewall_vpc_cen_tr_firewall.VpcCenTrFirewall.id}",
					"policy_name":        "222222",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"src_candidate_list.#": "2",
						"policy_type":          "fullmesh",
						"policy_description":   CHECKSET,
						"firewall_id":          CHECKSET,
						"policy_name":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"src_candidate_list": []map[string]interface{}{
						{
							"candidate_id":   "${alicloud_vpc.vpc1.id}",
							"candidate_type": "VPC",
						},
						{
							"candidate_id":   "${alicloud_vpc.vpc2.id}",
							"candidate_type": "VPC",
						},
						{
							"candidate_id":   "${alicloud_vpc.vpc3.id}",
							"candidate_type": "VPC",
						},
						{
							"candidate_id":   "${alicloud_express_connect_router_express_connect_router.ExpressConnectRouter.id}",
							"candidate_type": "ECR",
						},
					},
					"should_recover":      "true",
					"dest_candidate_list": []map[string]interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"src_candidate_list.#":  "4",
						"should_recover":        CHECKSET,
						"dest_candidate_list.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "closed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "closed",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "should_recover"},
			},
		},
	})
}

var AlicloudCloudFirewallVpcCenTrFirewallPolicyMap10781 = map[string]string{
	"tr_firewall_route_policy_id": CHECKSET,
}

func AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence10781(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone4" {
  default = "cn-hangzhou-k"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

variable "zone3" {
  default = "cn-hangzhou-j"
}

data "alicloud_cen_transit_router_service" "default" {
	enable = "On"
}

resource "alicloud_cen_instance" "cen" {
  description       = "yqc-test"
  cen_instance_name = "yqc-test-CenInstance"
}

resource "alicloud_cen_transit_router" "tr" {
  cen_id              = alicloud_cen_instance.cen.id
  transit_router_name = "yqc-test-TransitRouter"
}

resource "alicloud_express_connect_router_express_connect_router" "ExpressConnectRouter" {
  ecr_name         = "yqc-test-ecr"
  alibaba_side_asn = "65520"
  description      = "22222"
}

resource "alicloud_express_connect_router_tr_association" "ExpressConnectRouterTrAssociation" {
  association_region_id = var.region
  ecr_id                = alicloud_express_connect_router_express_connect_router.ExpressConnectRouter.id
  cen_id                = alicloud_cen_instance.cen.id
  transit_router_id     = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_cen_transit_router_ecr_attachment" "ExpressConnectRouterTrAssociation" {
  ecr_id                                = alicloud_express_connect_router_express_connect_router.ExpressConnectRouter.id
  cen_id                                = alicloud_cen_instance.cen.id
  transit_router_ecr_attachment_name    = "yqc-test-TransitRouterEcrAttachmentName"
  transit_router_attachment_description = "yqc-test-TransitRouterAttachmentDescription"
  transit_router_id                     = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_vpc" "vpc1" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-test-01"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  vpc_id     = alicloud_vpc.vpc1.id
  cidr_block = "172.16.1.0/24"
  zone_id    = var.zone1
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id     = alicloud_vpc.vpc1.id
  cidr_block = "172.16.2.0/24"
  zone_id    = var.zone2
}

resource "alicloud_vpc" "vpc2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-test-02"
}

resource "alicloud_vswitch" "vpc2vsw1" {
  vpc_id     = alicloud_vpc.vpc2.id
  zone_id    = var.zone1
  cidr_block = "172.16.3.0/24"
}

resource "alicloud_vswitch" "vpc2vsw2" {
  vpc_id     = alicloud_vpc.vpc2.id
  cidr_block = "172.16.4.0/24"
  zone_id    = var.zone2
}

resource "alicloud_vpc" "vpc3" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-test-03"
}

resource "alicloud_vswitch" "vpc3vsw1" {
  vpc_id     = alicloud_vpc.vpc3.id
  zone_id    = var.zone1
  cidr_block = "172.17.1.0/24"
}

resource "alicloud_vswitch" "vpc3vsw2" {
  vpc_id     = alicloud_vpc.vpc3.id
  cidr_block = "172.17.2.0/24"
  zone_id    = var.zone2
}

resource "alicloud_vpc" "vpc4" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-test-04"
}

resource "alicloud_vswitch" "vpc4vsw1" {
  vpc_id     = alicloud_vpc.vpc4.id
  zone_id    = var.zone1
  cidr_block = "172.16.8.0/24"
}

resource "alicloud_vswitch" "vpc4vsw2" {
  vpc_id     = alicloud_vpc.vpc4.id
  zone_id    = var.zone2
  cidr_block = "172.16.9.0/24"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  vpc_id = alicloud_vpc.vpc1.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
  }
  cen_id                             = alicloud_cen_instance.cen.id
  transit_router_id                  = alicloud_cen_transit_router.tr.transit_router_id
  auto_publish_route_enabled         = true
  transit_router_vpc_attachment_name = "TransitRouterVpcAttachmentName-1"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  auto_publish_route_enabled = true
  vpc_id                     = alicloud_vpc.vpc2.id
  cen_id                     = alicloud_cen_instance.cen.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = alicloud_vswitch.vpc2vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = alicloud_vswitch.vpc2vsw2.zone_id
  }
  transit_router_id                  = alicloud_cen_transit_router.tr.transit_router_id
  transit_router_vpc_attachment_name = "TransitRouterVpcAttachmentName-2"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc3" {
  auto_publish_route_enabled = true
  vpc_id                     = alicloud_vpc.vpc3.id
  cen_id                     = alicloud_cen_instance.cen.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc3vsw1.id
    zone_id    = alicloud_vswitch.vpc3vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc3vsw2.id
    zone_id    = alicloud_vswitch.vpc3vsw2.zone_id
  }
  transit_router_id                  = alicloud_cen_transit_router.tr.transit_router_id
  transit_router_vpc_attachment_name = "TransitRouterVpcAttachmentName-3"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc4" {
  vpc_id = alicloud_vpc.vpc4.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc4vsw1.id
    zone_id    = alicloud_vswitch.vpc4vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc4vsw2.id
    zone_id    = alicloud_vswitch.vpc4vsw2.zone_id
  }
  cen_id                                = alicloud_cen_instance.cen.id
  transit_router_id                     = alicloud_cen_transit_router.tr.transit_router_id
  transit_router_vpc_attachment_name    = "TransitRouterVpcAttachmentName-4"
  auto_publish_route_enabled            = true
  transit_router_attachment_description = "TransitRouterAttachmentDescription4"
}

resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "VpcCenTrFirewall" {
  route_mode                = "managed"
  region_no                 = var.region
  firewall_description      = "VpcCenTrFirewall created by terraform"
  tr_attachment_master_zone = var.zone1
  firewall_name             = "yqc-test-Firewall"
  tr_attachment_master_cidr = "10.0.2.0/24"
  firewall_subnet_cidr      = "10.0.1.0/24"
  cen_id                    = alicloud_cen_instance.cen.id
  tr_attachment_slave_cidr  = "10.0.3.0/24"
  tr_attachment_slave_zone  = var.zone2
  firewall_vpc_cidr         = "10.0.0.0/16"
  transit_router_id         = alicloud_cen_transit_router_vpc_attachment.tr-vpc4.transit_router_id

  depends_on = [
    alicloud_cen_transit_router_vpc_attachment.tr-vpc1,
    alicloud_cen_transit_router_vpc_attachment.tr-vpc2,
    alicloud_cen_transit_router_vpc_attachment.tr-vpc3,
    alicloud_cen_transit_router_vpc_attachment.tr-vpc4,
    alicloud_cen_transit_router_ecr_attachment.ExpressConnectRouterTrAssociation,
    alicloud_express_connect_router_tr_association.ExpressConnectRouterTrAssociation,
    alicloud_cen_transit_router_route_table.TransitRouterRouteTable,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation1,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation2,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation3,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation4,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation5,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation1,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation2,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation3,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation4,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation5,
  ]
}

resource "alicloud_cen_transit_router_route_table" "TransitRouterRouteTable" {
  transit_router_route_table_description = "111"
  transit_router_route_table_name        = "222"
  transit_router_id                      = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation1" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation1" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation2" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc2.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation2" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc2.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation3" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc3.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation3" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc3.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation4" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc4.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation4" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc4.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation5" {
  transit_router_attachment_id  = alicloud_cen_transit_router_ecr_attachment.ExpressConnectRouterTrAssociation.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation5" {
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
  transit_router_attachment_id  = alicloud_cen_transit_router_ecr_attachment.ExpressConnectRouterTrAssociation.id
}


`, name)
}

// Case VpcCenTrFirewallPolicy-end_to_end_(dest all)_预发测试 10600
func TestAccAliCloudCloudFirewallVpcCenTrFirewallPolicy_basic10600(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcCenTrFirewallPolicyMap10600)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewallPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence10600)
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
					"src_candidate_list": []map[string]interface{}{
						{
							"candidate_id":   "${alicloud_vpc.vpc1.id}",
							"candidate_type": "VPC",
						},
					},
					"policy_type":        "end_to_end",
					"policy_description": "test-dolicyDescription",
					"firewall_id":        "${alicloud_cloud_firewall_vpc_cen_tr_firewall.VpcCenTrFirewall.id}",
					"policy_name":        "yqc-test",
					"dest_candidate_list": []map[string]interface{}{
						{
							"candidate_type": "ALL",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"src_candidate_list.#":  "1",
						"policy_type":           "end_to_end",
						"policy_description":    "test-dolicyDescription",
						"firewall_id":           CHECKSET,
						"policy_name":           "yqc-test",
						"dest_candidate_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":         "closed",
					"should_recover": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         "closed",
						"should_recover": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "should_recover"},
			},
		},
	})
}

var AlicloudCloudFirewallVpcCenTrFirewallPolicyMap10600 = map[string]string{
	"tr_firewall_route_policy_id": CHECKSET,
}

func AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence10600(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

resource "alicloud_cen_instance" "cen" {
  description       = "yqc-test"
  cen_instance_name = "yqc-test-CenInstance"
}

resource "alicloud_cen_transit_router" "tr" {
  cen_id              = alicloud_cen_instance.cen.id
  transit_router_name = "yqc-test-TransitRouter"
}

resource "alicloud_vpc" "vpc1" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-test-01"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  vpc_id     = alicloud_vpc.vpc1.id
  cidr_block = "172.16.1.0/24"
  zone_id    = var.zone1
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id     = alicloud_vpc.vpc1.id
  cidr_block = "172.16.2.0/24"
  zone_id    = var.zone2
}

resource "alicloud_vpc" "vpc2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-test-02"
}

resource "alicloud_vswitch" "vpc2vsw1" {
  vpc_id     = alicloud_vpc.vpc2.id
  zone_id    = var.zone1
  cidr_block = "172.16.3.0/24"
}

resource "alicloud_vswitch" "vpc2vsw2" {
  vpc_id     = alicloud_vpc.vpc2.id
  cidr_block = "172.16.4.0/24"
  zone_id    = var.zone2
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  vpc_id = alicloud_vpc.vpc1.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
  }
  cen_id                     = alicloud_cen_instance.cen.id
  transit_router_id          = alicloud_cen_transit_router.tr.transit_router_id
  auto_publish_route_enabled = true
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  auto_publish_route_enabled = true
  vpc_id                     = alicloud_vpc.vpc2.id
  cen_id                     = alicloud_cen_instance.cen.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = alicloud_vswitch.vpc2vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = alicloud_vswitch.vpc2vsw2.zone_id
  }
  transit_router_id = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_cen_transit_router_route_table" "TransitRouterRouteTable" {
  transit_router_route_table_description = "111"
  transit_router_route_table_name        = "222"
  transit_router_id                      = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation1" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation1" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation2" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc2.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation2" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc2.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "VpcCenTrFirewall" {
  route_mode                = "managed"
  region_no                 = var.region
  firewall_description      = "VpcCenTrFirewall created by terraform"
  tr_attachment_master_zone = var.zone1
  firewall_name             = "yqc-test-Firewall"
  tr_attachment_master_cidr = "10.0.2.0/24"
  firewall_subnet_cidr      = "10.0.1.0/24"
  cen_id                    = alicloud_cen_instance.cen.id
  tr_attachment_slave_cidr  = "10.0.3.0/24"
  tr_attachment_slave_zone  = var.zone2
  firewall_vpc_cidr         = "10.0.0.0/16"
  transit_router_id         = alicloud_cen_transit_router.tr.transit_router_id
  depends_on = [
    alicloud_cen_transit_router_vpc_attachment.tr-vpc1,
    alicloud_cen_transit_router_vpc_attachment.tr-vpc2,
    alicloud_cen_transit_router_route_table.TransitRouterRouteTable,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation1,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation2,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation1,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation2,
  ]
}


`, name)
}

// Case VpcCenTrFirewallPolicy-one_to_one(add status ) 10769
func TestAccAliCloudCloudFirewallVpcCenTrFirewallPolicy_basic10769(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_cen_tr_firewall_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallVpcCenTrFirewallPolicyMap10769)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudFirewallServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcCenTrFirewallPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence10769)
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
					"src_candidate_list": []map[string]interface{}{
						{
							"candidate_id":   "${alicloud_vpc.vpc1.id}",
							"candidate_type": "VPC",
						},
					},
					"policy_type":        "one_to_one",
					"policy_description": "111111",
					"firewall_id":        "${alicloud_cloud_firewall_vpc_cen_tr_firewall.VpcCenTrFirewall.id}",
					"policy_name":        "222222",
					"dest_candidate_list": []map[string]interface{}{
						{
							"candidate_type": "VPC",
							"candidate_id":   "${alicloud_vpc.vpc2.id}",
						},
					},
					"status": "closed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"src_candidate_list.#":  "1",
						"policy_type":           "one_to_one",
						"policy_description":    CHECKSET,
						"firewall_id":           CHECKSET,
						"policy_name":           CHECKSET,
						"dest_candidate_list.#": "1",
						"status":                "closed",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":         "opened",
					"should_recover": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":         "opened",
						"should_recover": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang", "should_recover"},
			},
		},
	})
}

var AlicloudCloudFirewallVpcCenTrFirewallPolicyMap10769 = map[string]string{
	"tr_firewall_route_policy_id": CHECKSET,
}

func AlicloudCloudFirewallVpcCenTrFirewallPolicyBasicDependence10769(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone4" {
  default = "cn-hangzhou-k"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

variable "zone3" {
  default = "cn-hangzhou-j"
}

resource "alicloud_cen_instance" "cen" {
  description       = "yqc-test"
  cen_instance_name = "yqc-test-CenInstance"
}

resource "alicloud_cen_transit_router" "tr" {
  cen_id              = alicloud_cen_instance.cen.id
  transit_router_name = "yqc-test-TransitRouter"
}

resource "alicloud_vpc" "vpc1" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-test-01"
}

resource "alicloud_vswitch" "vpc1vsw1" {
  vpc_id     = alicloud_vpc.vpc1.id
  cidr_block = "172.16.1.0/24"
  zone_id    = var.zone1
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id     = alicloud_vpc.vpc1.id
  cidr_block = "172.16.2.0/24"
  zone_id    = var.zone2
}

resource "alicloud_vpc" "vpc2" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "yqc-vpc-test-02"
}

resource "alicloud_vswitch" "vpc2vsw1" {
  vpc_id     = alicloud_vpc.vpc2.id
  zone_id    = var.zone1
  cidr_block = "172.16.3.0/24"
}

resource "alicloud_vswitch" "vpc2vsw2" {
  vpc_id     = alicloud_vpc.vpc2.id
  cidr_block = "172.16.4.0/24"
  zone_id    = var.zone2
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  vpc_id = alicloud_vpc.vpc1.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = alicloud_vswitch.vpc1vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
    zone_id    = alicloud_vswitch.vpc1vsw2.zone_id
  }
  cen_id                             = alicloud_cen_instance.cen.id
  transit_router_id                  = alicloud_cen_transit_router.tr.transit_router_id
  auto_publish_route_enabled         = true
  transit_router_vpc_attachment_name = "TransitRouterVpcAttachmentName-1"
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc2" {
  auto_publish_route_enabled = true
  vpc_id                     = alicloud_vpc.vpc2.id
  cen_id                     = alicloud_cen_instance.cen.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw1.id
    zone_id    = alicloud_vswitch.vpc2vsw1.zone_id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc2vsw2.id
    zone_id    = alicloud_vswitch.vpc2vsw2.zone_id
  }
  transit_router_id                  = alicloud_cen_transit_router.tr.transit_router_id
  transit_router_vpc_attachment_name = "TransitRouterVpcAttachmentName-2"
}

resource "alicloud_cen_transit_router_route_table" "TransitRouterRouteTable" {
  transit_router_route_table_description = "111"
  transit_router_route_table_name        = "222"
  transit_router_id                      = alicloud_cen_transit_router.tr.transit_router_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation1" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation1" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc1.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_association" "TransitRouterRouteTableAssociation2" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc2.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cen_transit_router_route_table_propagation" "TransitRouterRouteTablePropagation2" {
  transit_router_attachment_id  = alicloud_cen_transit_router_vpc_attachment.tr-vpc2.id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.TransitRouterRouteTable.transit_router_route_table_id
}

resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "VpcCenTrFirewall" {
  route_mode                = "managed"
  region_no                 = var.region
  firewall_description      = "VpcCenTrFirewall created by terraform"
  tr_attachment_master_zone = var.zone1
  firewall_name             = "yqc-test-Firewall"
  tr_attachment_master_cidr = "10.0.2.0/24"
  firewall_subnet_cidr      = "10.0.1.0/24"
  cen_id                    = alicloud_cen_instance.cen.id
  tr_attachment_slave_cidr  = "10.0.3.0/24"
  tr_attachment_slave_zone  = var.zone2
  firewall_vpc_cidr         = "10.0.0.0/16"
  transit_router_id         = alicloud_cen_transit_router.tr.transit_router_id
  depends_on = [
    alicloud_cen_transit_router_vpc_attachment.tr-vpc1,
    alicloud_cen_transit_router_vpc_attachment.tr-vpc2,
    alicloud_cen_transit_router_route_table.TransitRouterRouteTable,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation1,
    alicloud_cen_transit_router_route_table_association.TransitRouterRouteTableAssociation2,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation1,
    alicloud_cen_transit_router_route_table_propagation.TransitRouterRouteTablePropagation2,
  ]
}


`, name)
}

// Test CloudFirewall VpcCenTrFirewallPolicy. <<< Resource test cases, automatically generated.
