package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCenTransitRouterMulticastDomainAssociation_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_transit_router_multicast_domain_association.default"
	ra := resourceAttrInit(resourceId, resourceAliCloudCenTransitRouterMulticastDomainAssociationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterMulticastDomainAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterMulticastDomainAssociation-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAliCloudCenTransitRouterMulticastDomainAssociationBasicDependence)
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
					"vswitch_id":                         "${alicloud_vswitch.default.id}",
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

var resourceAliCloudCenTransitRouterMulticastDomainAssociationMap = map[string]string{
	"status": CHECKSET,
}

func resourceAliCloudCenTransitRouterMulticastDomainAssociationBasicDependence(name string) string {
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
  		vpc_id            = alicloud_vswitch.default.vpc_id
  		zone_mappings {
    		zone_id    = alicloud_vswitch.default.zone_id
    		vswitch_id = alicloud_vswitch.default.id
  		}
	}
`, name)
}
