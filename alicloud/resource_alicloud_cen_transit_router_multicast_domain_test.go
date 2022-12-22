package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenTransitRouterMulticastDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_transit_router_multicast_domain.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCenTransitRouterMulticastDomainMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterMulticastDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterMulticastDomain-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCenTransitRouterMulticastDomainBasicDependence)
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
					"transit_router_id":                           "${alicloud_cen_transit_router.default.transit_router_id}",
					"transit_router_multicast_domain_name":        name,
					"transit_router_multicast_domain_description": name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "TransitRouterMulticastDomain",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_id":                           CHECKSET,
						"transit_router_multicast_domain_name":        name,
						"transit_router_multicast_domain_description": name,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "TransitRouterMulticastDomain",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_multicast_domain_name":        name + "_update",
					"transit_router_multicast_domain_description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_multicast_domain_name":        name + "_update",
						"transit_router_multicast_domain_description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "TransitRouterMulticastDomain_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "TransitRouterMulticastDomain_Update",
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

var resourceAlicloudCenTransitRouterMulticastDomainMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudCenTransitRouterMulticastDomainBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id            = alicloud_cen_instance.default.id
  		support_multicast = true
	}
`, name)
}
