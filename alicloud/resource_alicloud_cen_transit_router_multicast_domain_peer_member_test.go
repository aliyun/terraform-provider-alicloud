package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCenTransitRouterMulticastDomainPeerMember_basic1905(t *testing.T) {
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
	ra := resourceAttrInit(resourceId, AliCloudCenTransitRouterMulticastDomainPeerMemberMap1905)
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
						"transit_router_multicast_domain_id":      CHECKSET,
						"peer_transit_router_multicast_domain_id": CHECKSET,
						"group_ip_address":                        CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudCenTransitRouterMulticastDomainPeerMemberMap1905 = map[string]string{}

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
  		provider          = alicloud.hz
  		cen_id            = alicloud_cen_bandwidth_package_attachment.default.instance_id
  		support_multicast = true
	}

	resource "alicloud_cen_transit_router" "peer" {
  		provider          = alicloud.qd
  		cen_id            = alicloud_cen_bandwidth_package_attachment.default.instance_id
  		support_multicast = true
	}

	resource "alicloud_cen_transit_router_peer_attachment" "default" {
  		provider                              = alicloud.hz
  		cen_id                                = alicloud_cen_bandwidth_package_attachment.default.instance_id
  		transit_router_id                     = alicloud_cen_transit_router.default.transit_router_id
  		peer_transit_router_id                = alicloud_cen_transit_router.peer.transit_router_id
  		peer_transit_router_region_id         = "cn-qingdao"
  		cen_bandwidth_package_id              = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  		bandwidth                             = 5
  		transit_router_attachment_description = var.name
  		transit_router_attachment_name        = var.name
	}

	resource "alicloud_cen_transit_router_multicast_domain" "default" {
  		provider                                    = alicloud.hz
  		transit_router_id                           = alicloud_cen_transit_router_peer_attachment.default.transit_router_id
  		transit_router_multicast_domain_name        = var.name
  		transit_router_multicast_domain_description = var.name
	}

	resource "alicloud_cen_transit_router_multicast_domain" "peer" {
  		provider                                    = alicloud.qd
  		transit_router_id                           = alicloud_cen_transit_router_peer_attachment.default.peer_transit_router_id
  		transit_router_multicast_domain_name        = var.name
  		transit_router_multicast_domain_description = var.name
	}

	resource "alicloud_cen_transit_router_multicast_domain_peer_member" "default" {
  		provider                                = alicloud.hz
  		transit_router_multicast_domain_id      = alicloud_cen_transit_router_multicast_domain.default.id
  		peer_transit_router_multicast_domain_id = alicloud_cen_transit_router_multicast_domain.peer.id
  		group_ip_address                        = "239.0.0.8"
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
