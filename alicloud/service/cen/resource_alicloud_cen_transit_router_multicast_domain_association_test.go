package cen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCenTransitRouterMulticastDomainAssociation_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_transit_router_multicast_domain_association.default"
	ra := acctest.ResourceAttrInit(resourceId, resourceAliCloudCenTransitRouterMulticastDomainAssociationMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCbnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenTransitRouterMulticastDomainAssociation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterMulticastDomainAssociation-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAliCloudCenTransitRouterMulticastDomainAssociationBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_multicast_domain_id": "${alicloud_cen_transit_router_multicast_domain.default.id}",
					"transit_router_attachment_id":       "${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}",
					"vswitch_id":                         "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_multicast_domain_id": acctest.CHECKSET,
						"transit_router_attachment_id":       acctest.CHECKSET,
						"vswitch_id":                         acctest.CHECKSET,
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
	"status": acctest.CHECKSET,
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
