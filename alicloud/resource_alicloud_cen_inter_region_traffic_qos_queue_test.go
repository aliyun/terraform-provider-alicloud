package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAliCloudCenInterRegionTrafficQosQueue_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_inter_region_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCenInterRegionTrafficQosQueueMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenInterRegionTrafficQosQueue-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCenInterRegionTrafficQosQueueBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckCenInterRegionTrafficQosQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"remain_bandwidth_percent": "20",
					"traffic_qos_policy_id":    "${alicloud_cen_inter_region_traffic_qos_policy.default.id}",
					"dscps":                    []string{"2", "3"},
				}),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccAliCloudCenInterRegionTrafficQosQueue_basic1(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_inter_region_traffic_qos_queue.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCenInterRegionTrafficQosQueueMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenInterRegionTrafficQosQueue-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCenInterRegionTrafficQosQueueBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckCenInterRegionTrafficQosQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"traffic_qos_policy_id": "${alicloud_cen_inter_region_traffic_qos_policy.default.id}",
					"dscps":                 []string{"2", "3"},
					"bandwidth":             "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"traffic_qos_policy_id": CHECKSET,
						"dscps.#":               "2",
						"bandwidth":             "1",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":                                  "2",
					"inter_region_traffic_qos_queue_name":        name + "_update",
					"inter_region_traffic_qos_queue_description": "testDescription",
					"dscps": []string{"4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":                                  "2",
						"inter_region_traffic_qos_queue_name":        name + "_update",
						"inter_region_traffic_qos_queue_description": "testDescription",
						"dscps.#": "1",
					}),
				),
			},
		},
	})
}

func testAccCheckCenInterRegionTrafficQosQueueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
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
			return fmt.Errorf("Cen Inter Region Traffic Qos Queue still exist,  ID %s ", fmt.Sprint(resp["TrafficQosQueueId"]))
		}
	}

	return nil
}

var resourceAlicloudCenInterRegionTrafficQosQueueMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudCenInterRegionTrafficQosQueueBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	%s

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_bandwidth_package" "default" {
  		bandwidth              = 5
  		geographic_region_a_id = "China"
  		geographic_region_b_id = "China"
	}

	resource "alicloud_cen_bandwidth_package_attachment" "default" {
  		instance_id          = alicloud_cen_instance.default.id
  		bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
	}

	resource "alicloud_cen_transit_router" "hz" {
  		cen_id   = alicloud_cen_bandwidth_package_attachment.default.instance_id
	}

	resource "alicloud_cen_transit_router" "bj" {
  		provider = alicloudalt
  		cen_id   = alicloud_cen_transit_router.hz.cen_id
	}

	resource "alicloud_cen_transit_router_peer_attachment" "default" {
  		cen_id                        = alicloud_cen_instance.default.id
  		transit_router_id             = alicloud_cen_transit_router.hz.transit_router_id
  		peer_transit_router_region_id = "cn-beijing"
  		peer_transit_router_id        = alicloud_cen_transit_router.bj.transit_router_id
  		cen_bandwidth_package_id      = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  		bandwidth                     = 5
	}

	resource "alicloud_cen_inter_region_traffic_qos_policy" "default" {
  		transit_router_id                           = alicloud_cen_transit_router.hz.transit_router_id
  		transit_router_attachment_id                = alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id
  		inter_region_traffic_qos_policy_name        = var.name
  		inter_region_traffic_qos_policy_description = var.name
	}

`, name, configAlternateRegionProvider("cn-beijing"))
}

func resourceAlicloudCenInterRegionTrafficQosQueueBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	%s

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_bandwidth_package" "default" {
  		bandwidth              = 5
  		geographic_region_a_id = "China"
  		geographic_region_b_id = "China"
	}

	resource "alicloud_cen_bandwidth_package_attachment" "default" {
  		instance_id          = alicloud_cen_instance.default.id
  		bandwidth_package_id = alicloud_cen_bandwidth_package.default.id
	}

	resource "alicloud_cen_transit_router" "hz" {
  		cen_id   = alicloud_cen_bandwidth_package_attachment.default.instance_id
	}

	resource "alicloud_cen_transit_router" "bj" {
  		provider = alicloudalt
  		cen_id   = alicloud_cen_transit_router.hz.cen_id
	}

	resource "alicloud_cen_transit_router_peer_attachment" "default" {
  		cen_id                        = alicloud_cen_instance.default.id
  		transit_router_id             = alicloud_cen_transit_router.hz.transit_router_id
  		peer_transit_router_region_id = "cn-beijing"
  		peer_transit_router_id        = alicloud_cen_transit_router.bj.transit_router_id
  		cen_bandwidth_package_id      = alicloud_cen_bandwidth_package_attachment.default.bandwidth_package_id
  		bandwidth                     = 5
	}

	resource "alicloud_cen_inter_region_traffic_qos_policy" "default" {
  		transit_router_id                           = alicloud_cen_transit_router.hz.transit_router_id
  		transit_router_attachment_id                = alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id
  		inter_region_traffic_qos_policy_name        = var.name
  		inter_region_traffic_qos_policy_description = var.name
        bandwidth_guarantee_mode = "byBandwidth"
	}

`, name, configAlternateRegionProvider("cn-beijing"))
}
