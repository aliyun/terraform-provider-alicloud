package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudCenTransitRouterMulticastDomainMember_basic1906(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_multicast_domain_member.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterMulticastDomainMemberMap1906)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterMulticastDomainMember")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.CENTransitRouterMulticastDomainMemberSupportRegions)
	name := fmt.Sprintf("tf-testacc%sCenTransitRouterMulticastDomainMember%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterMulticastDomainMemberBasicDependence1906)
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
					"vpc_id":                             "${data.alicloud_vpcs.default.ids.0}",
					"transit_router_multicast_domain_id": "${data.alicloud_cen_transit_router_multicast_domains.default.ids.0}",
					"network_interface_id":               "${alicloud_ecs_network_interface.default.id}",
					"group_ip_address":                   "224.0.0.8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":                             CHECKSET,
						"transit_router_multicast_domain_id": CHECKSET,
						"network_interface_id":               CHECKSET,
						"group_ip_address":                   CHECKSET,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenTransitRouterMulticastDomainMemberMap1906 = map[string]string{}

func AlicloudCenTransitRouterMulticastDomainMemberBasicDependence1906(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = data.alicloud_vswitches.default.ids.0
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic test"
  primary_ip_address     = cidrhost(data.alicloud_vswitches.default.vswitches.0.cidr_block, 100)
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
data "alicloud_cen_instances" "default" {
  name_regex = "no-deleting-cen"
}
data "alicloud_cen_transit_routers" "default" {
  cen_id     = data.alicloud_cen_instances.default.instances.0.id
  name_regex = "no-deleting-cen"
}
data "alicloud_cen_transit_router_multicast_domains" "default" {
  transit_router_id = data.alicloud_cen_transit_routers.default.transit_routers.0.transit_router_id
  name_regex        = "no-deleting-cen"
}

`, name)
}
