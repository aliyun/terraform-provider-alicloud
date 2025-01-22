package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCenTransitRouterMulticastDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_transit_router_multicast_domain.default"
	ra := resourceAttrInit(resourceId, resourceAliCloudCenTransitRouterMulticastDomainMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterMulticastDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterMulticastDomain-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAliCloudCenTransitRouterMulticastDomainBasicDependence)
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
					"transit_router_id": "${alicloud_cen_transit_router.default.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_multicast_domain_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_multicast_domain_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_multicast_domain_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_multicast_domain_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options": []map[string]interface{}{
						{
							"igmpv2_support": "enable",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudCenTransitRouterMulticastDomain_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_transit_router_multicast_domain.default"
	ra := resourceAttrInit(resourceId, resourceAliCloudCenTransitRouterMulticastDomainMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterMulticastDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterMulticastDomain-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAliCloudCenTransitRouterMulticastDomainBasicDependence)
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
					"options": []map[string]interface{}{
						{
							"igmpv2_support": "enable",
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_id":                           CHECKSET,
						"transit_router_multicast_domain_name":        name,
						"transit_router_multicast_domain_description": name,
						"options.#":    "1",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
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

var resourceAliCloudCenTransitRouterMulticastDomainMap = map[string]string{
	"options.#": CHECKSET,
	"region_id": CHECKSET,
	"status":    CHECKSET,
}

func resourceAliCloudCenTransitRouterMulticastDomainBasicDependence(name string) string {
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
