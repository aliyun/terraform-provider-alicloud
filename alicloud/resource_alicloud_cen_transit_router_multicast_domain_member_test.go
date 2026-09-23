package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCenTransitRouterMulticastDomainMember_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CENTransitRouterMulticastDomainMemberSupportRegions)
	resourceId := "alicloud_cen_transit_router_multicast_domain_member.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterMulticastDomainMemberMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterMulticastDomainMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCenTransitRouterMulticastDomainMember%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterMulticastDomainMemberBasicDependence0)
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
					"transit_router_multicast_domain_id": "${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id}",
					"group_ip_address":                   "239.0.0.8",
					"network_interface_id":               "${alicloud_ecs_network_interface.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_multicast_domain_id": CHECKSET,
						"group_ip_address":                   CHECKSET,
						"network_interface_id":               CHECKSET,
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

func TestAccAliCloudCenTransitRouterMulticastDomainMember_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CENTransitRouterMulticastDomainMemberSupportRegions)
	resourceId := "alicloud_cen_transit_router_multicast_domain_member.default"
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterMulticastDomainMemberMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterMulticastDomainMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCenTransitRouterMulticastDomainMember%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCenTransitRouterMulticastDomainMemberBasicDependence0)
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
					"transit_router_multicast_domain_id": "${alicloud_cen_transit_router_multicast_domain_association.default.transit_router_multicast_domain_id}",
					"group_ip_address":                   "239.0.0.8",
					"network_interface_id":               "${alicloud_ecs_network_interface.default.id}",
					"vpc_id":                             "${alicloud_cen_transit_router_vpc_attachment.default.vpc_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_multicast_domain_id": CHECKSET,
						"group_ip_address":                   CHECKSET,
						"network_interface_id":               CHECKSET,
						"vpc_id":                             CHECKSET,
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

var AliCloudCenTransitRouterMulticastDomainMemberMap0 = map[string]string{}

func AliCloudCenTransitRouterMulticastDomainMemberBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = "cn-hangzhou-i"
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vswitch.default.vpc_id
	}

	resource "alicloud_ecs_network_interface" "default" {
  		vswitch_id         = alicloud_vswitch.default.id
  		security_group_ids = [alicloud_security_group.default.id]
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_bandwidth_package" "default" {
  		bandwidth                  = 5
  		cen_bandwidth_package_name = var.name
  		geographic_region_a_id     = "China"
  		geographic_region_b_id     = "China"
	}

	resource "alicloud_cen_bandwidth_package_attachment" "default" {
  		instance_id          = alicloud_cen_instance.default.id
  		bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id            = alicloud_cen_bandwidth_package_attachment.default.instance_id
  		support_multicast = true
	}

	resource "alicloud_cen_transit_router_multicast_domain" "default" {
  		transit_router_id                           = alicloud_cen_transit_router.default.transit_router_id
  		transit_router_multicast_domain_name        = var.name
  		transit_router_multicast_domain_description = var.name
	}

	resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  		cen_id                                = alicloud_cen_transit_router.default.cen_id
  		transit_router_id                     = alicloud_cen_transit_router_multicast_domain.default.transit_router_id
  		vpc_id                                = alicloud_vswitch.default.vpc_id
  		transit_router_attachment_description = var.name
  		transit_router_attachment_name        = var.name
  		zone_mappings {
    		vswitch_id = alicloud_vswitch.default.id
    		zone_id    = alicloud_vswitch.default.zone_id
  		}
	}

	resource "alicloud_cen_transit_router_multicast_domain_association" "default" {
  		transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.default.id
  		transit_router_attachment_id       = alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
  		vswitch_id                         = alicloud_vswitch.default.id
	}
`, name)
}
