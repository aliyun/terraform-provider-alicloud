package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenTransitRouterMulticastDomainAssociation_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_transit_router_multicast_domain_association.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCenTransitRouterMulticastDomainAssociationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterMulticastDomainAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterMulticastDomainAssociation-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCenTransitRouterMulticastDomainAssociationBasicDependence)
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
					"transit_router_multicast_domain_id": "${alicloud_cen_transit_router_multicast_domain.default.id}",
					"transit_router_attachment_id":       "${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}",
					"vswitch_id":                         "${data.alicloud_vswitches.default.vswitches.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_multicast_domain_id": CHECKSET,
						"transit_router_attachment_id":       CHECKSET,
						"vswitch_id":                         CHECKSET,
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

var resourceAlicloudCenTransitRouterMulticastDomainAssociationMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudCenTransitRouterMulticastDomainAssociationBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		name_regex = "default-zone-i"
  		vpc_id     = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id            = alicloud_cen_instance.default.id
  		support_multicast = true
	}

	resource "alicloud_cen_transit_router_multicast_domain" "default" {
  		transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	}

	resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  		cen_id            = alicloud_cen_transit_router.default.cen_id
  		transit_router_id = alicloud_cen_transit_router_multicast_domain.default.transit_router_id
  		vpc_id            = data.alicloud_vpcs.default.ids.0
  		zone_mappings {
    		zone_id    = data.alicloud_vswitches.default.vswitches.0.zone_id
    		vswitch_id = data.alicloud_vswitches.default.vswitches.0.id
  		}
	}
`, name)
}
