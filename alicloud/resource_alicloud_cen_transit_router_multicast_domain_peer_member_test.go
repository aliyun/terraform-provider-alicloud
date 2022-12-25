package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudCenTransitRouterMulticastDomainPeerMember_basic1905(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_multicast_domain_peer_member.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterMulticastDomainPeerMemberMap1905)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.CENTransitRouterMulticastDomainPeerMemberSupportRegions)
	name := fmt.Sprintf("tf-testacc%sCenTransitRouterMulticastDomainPeerMember%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		IDRefreshName:     resourceId,
		CheckDestroy:      testAccCheckCenInterRegionTransitRouterMulticastDomainPeerMemberDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccCenTransitRouterMulticastDomainPeerMemberCreateConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInterRegionTransitRouterMulticastDomainPeerMemberExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"peer_transit_router_multicast_domain_id": CHECKSET,
						"transit_router_multicast_domain_id":      CHECKSET,
						"group_ip_address":                        CHECKSET,
					}),
				),
			},
		},
	})
}

func testAccCenTransitRouterMulticastDomainPeerMemberCreateConfig(rand string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
provider "alicloud" {
  alias  = "hz"
  region = "cn-hangzhou"
}
provider "alicloud" {
  alias  = "qd"
  region = "cn-qingdao"
}
data "alicloud_cen_instances" "default" {
  provider   = alicloud.hz
  name_regex = "no-deleting-cen"
}
data "alicloud_cen_transit_routers" "default" {
  provider   = alicloud.hz
  cen_id     = data.alicloud_cen_instances.default.instances.0.id
  name_regex = "no-deleting-cen"
}
data "alicloud_cen_transit_routers" "peer" {
  provider   = alicloud.qd
  cen_id     = data.alicloud_cen_instances.default.instances.0.id
  name_regex = "qingdao-no-delete-cen"
}
data "alicloud_cen_transit_router_multicast_domains" "default" {
  provider          = alicloud.hz
  transit_router_id = data.alicloud_cen_transit_routers.default.transit_routers.0.transit_router_id
  name_regex        = "no-deleting-cen"
}
data "alicloud_cen_transit_router_multicast_domains" "peer" {
  provider          = alicloud.qd
  transit_router_id = data.alicloud_cen_transit_routers.peer.transit_routers.0.transit_router_id
  name_regex        = "no-deleting-cen"
}
resource "alicloud_cen_transit_router_multicast_domain_peer_member" "default" {
  transit_router_multicast_domain_id      = data.alicloud_cen_transit_router_multicast_domains.default.ids.0
  peer_transit_router_multicast_domain_id = data.alicloud_cen_transit_router_multicast_domains.peer.ids.0
  group_ip_address                        = "239.1.1.1"
  provider                                = alicloud.hz
}

`, rand)
}

func testAccCheckCenInterRegionTransitRouterMulticastDomainPeerMemberExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no alicloud_cen_inter_region_transit_router_multicast_domain_peer_member ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			cbnService := CbnService{client}

			resp, err := cbnService.DescribeCenTransitRouterMulticastDomainPeerMember(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}
			res = resp
			return nil
		}
		return fmt.Errorf("alicloud_cen_inter_region_transit_router_multicast_domain_peer_member not found")
	}
}

func testAccCheckCenInterRegionTransitRouterMulticastDomainPeerMemberDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenInterRegionTransitRouterMulticastDomainPeerMemberDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenInterRegionTransitRouterMulticastDomainPeerMemberDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "resource_alicloud_cen_transit_router_multicast_domain_peer_member" {
			continue
		}
		resp, err := cbnService.DescribeCenTransitRouterMulticastDomainPeerMember(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Cen Inter Region TransitRouterMulticastDomainPeerMember still exist,  ID %s ", fmt.Sprint(resp["TransitRouterMulticastDomainId"]))
		}
	}

	return nil
}

var AlicloudCenTransitRouterMulticastDomainPeerMemberMap1905 = map[string]string{}
