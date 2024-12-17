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

func TestAccAliCloudCenInterRegionTrafficQosQueue_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_inter_region_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCenInterRegionTrafficQosQueueMap)
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
	name := fmt.Sprintf("tf-testAccCenInterRegionTrafficQosQueue-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCenInterRegionTrafficQosQueueBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInterRegionTrafficQosQueueDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"remain_bandwidth_percent": "20",
					"traffic_qos_policy_id":    "${alicloud_cen_inter_region_traffic_qos_policy.default.id}",
					"dscps":                    []string{"2", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInterRegionTrafficQosQueueExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"remain_bandwidth_percent": "20",
						"traffic_qos_policy_id":    CHECKSET,
						"dscps.#":                  "2",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"remain_bandwidth_percent":                   "10",
					"inter_region_traffic_qos_queue_name":        name + "_update",
					"inter_region_traffic_qos_queue_description": "testDescription",
					"dscps": []string{"4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInterRegionTrafficQosQueueExistsWithProviders(resourceId, v, &providers),
					testAccCheck(map[string]string{
						"remain_bandwidth_percent":                   "10",
						"inter_region_traffic_qos_queue_name":        name + "_update",
						"inter_region_traffic_qos_queue_description": "testDescription",
						"dscps.#": "1",
					}),
				),
			},
		},
	})
}

func testAccCheckCenInterRegionTrafficQosQueueDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenInterRegionTrafficQosQueueDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenInterRegionTrafficQosQueueDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_inter_region_traffic_qos_queue" {
			continue
		}
		resp, err := cbnService.DescribeCenInterRegionTrafficQosQueue(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("Cen Inter Region Traffic Qos Policy still exist,  ID %s ", fmt.Sprint(resp["TrafficQosQueueId"]))
		}
	}

	return nil
}

func testAccCheckCenInterRegionTrafficQosQueueExistsWithProviders(n string, res map[string]interface{}, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No alicloud_cen_inter_region_traffic_qos_queue ID is set. ")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			cbnService := CbnService{client}

			resp, err := cbnService.DescribeCenInterRegionTrafficQosQueue(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}
			res = resp
			return nil
		}
		return fmt.Errorf("alicloud_cen_inter_region_traffic_qos_queue not found")
	}
}

var resourceAlicloudCenInterRegionTrafficQosQueueMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudCenInterRegionTrafficQosQueueBasicDependence(name string) string {
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

	resource "alicloud_cen_inter_region_traffic_qos_policy" "default" {
  		provider                                    = alicloud.hz
  		transit_router_id                           = alicloud_cen_transit_router.hz.transit_router_id
  		transit_router_attachment_id                = alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id
  		inter_region_traffic_qos_policy_name        = var.name
  		inter_region_traffic_qos_policy_description = var.name
	}

`, name)
}
