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

func TestAccAlicloudCenInterRegionTrafficQosPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_inter_region_traffic_qos_policy.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCenInterRegionTrafficQosPolicyMap)
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenInterRegionTrafficQosPolicy-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCenInterRegionTrafficQosPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInterRegionTrafficQosPolicyDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_id":                           "${alicloud_cen_transit_router.hz.transit_router_id}",
					"transit_router_attachment_id":                "${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}",
					"inter_region_traffic_qos_policy_name":        name,
					"inter_region_traffic_qos_policy_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInterRegionTrafficQosPolicyExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"transit_router_id":                           CHECKSET,
						"transit_router_attachment_id":                CHECKSET,
						"inter_region_traffic_qos_policy_name":        name,
						"inter_region_traffic_qos_policy_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inter_region_traffic_qos_policy_name":        name + "_update",
					"inter_region_traffic_qos_policy_description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInterRegionTrafficQosPolicyExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"inter_region_traffic_qos_policy_name":        name + "_update",
						"inter_region_traffic_qos_policy_description": name + "_update",
					}),
				),
			},
		},
	})
}

func testAccCheckCenInterRegionTrafficQosPolicyDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenInterRegionTrafficQosPolicyDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenInterRegionTrafficQosPolicyDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_inter_region_traffic_qos_policy" {
			continue
		}
		resp, err := cbnService.DescribeCenInterRegionTrafficQosPolicy(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Cen Inter Region Traffic Qos Policy still exist,  ID %s ", fmt.Sprint(resp["TrafficQosPolicyId"]))
		}
	}

	return nil
}

func testAccCheckCenInterRegionTrafficQosPolicyExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_cen_inter_region_traffic_qos_policy ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			cbnService := CbnService{client}

			resp, err := cbnService.DescribeCenInterRegionTrafficQosPolicy(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}
			res = resp
			return nil
		}
		return fmt.Errorf("alicloud_cen_inter_region_traffic_qos_policy not found")
	}
}

var resourceAlicloudCenInterRegionTrafficQosPolicyMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudCenInterRegionTrafficQosPolicyBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	provider "alicloud" {
  		alias  = "bj"
  		region = "cn-beijing"
	}

	provider "alicloud" {
  		alias  = "hz"
  		region = "cn-hangzhou"
	}

	resource "alicloud_cen_instance" "default" {
  		provider          = alicloud.hz
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_bandwidth_package" "default" {
  		provider               = alicloud.hz
  		bandwidth              = 5
  		geographic_region_a_id = "China"
  		geographic_region_b_id = "China"
	}

	resource "alicloud_cen_bandwidth_package_attachment" "default" {
  		provider             = alicloud.hz
  		instance_id          = alicloud_cen_instance.default.id
  		bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
	}

	resource "alicloud_cen_transit_router" "hz" {
  		provider = alicloud.hz
  		cen_id   = alicloud_cen_bandwidth_package_attachment.default.instance_id
	}

	resource "alicloud_cen_transit_router" "bj" {
  		provider = alicloud.bj
  		cen_id   = alicloud_cen_transit_router.hz.cen_id
	}

	resource "alicloud_cen_transit_router_peer_attachment" "default" {
  		provider                      = alicloud.hz
  		cen_id                        = alicloud_cen_instance.default.id
  		transit_router_id             = alicloud_cen_transit_router.hz.transit_router_id
  		peer_transit_router_region_id = "cn-beijing"
  		peer_transit_router_id        = alicloud_cen_transit_router.bj.transit_router_id
  		cen_bandwidth_package_id      = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  		bandwidth                     = 5
	}
`, name)
}
