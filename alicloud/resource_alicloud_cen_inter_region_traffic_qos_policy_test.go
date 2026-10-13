package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAliCloudCenInterRegionTrafficQosPolicy_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_cen_inter_region_traffic_qos_policy.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCenInterRegionTrafficQosPolicyMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenInterRegionTrafficQosPolicy-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCenInterRegionTrafficQosPolicyBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactoriesAlternate(),
		CheckDestroy:      testAccCheckCenInterRegionTrafficQosPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_id":                           "${alicloud_cen_transit_router.hz.transit_router_id}",
					"transit_router_attachment_id":                "${alicloud_cen_transit_router_peer_attachment.default.transit_router_attachment_id}",
					"inter_region_traffic_qos_policy_name":        name,
					"inter_region_traffic_qos_policy_description": name,
					"bandwidth_guarantee_mode":                    "byBandwidthPercent",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_id":                           CHECKSET,
						"transit_router_attachment_id":                CHECKSET,
						"inter_region_traffic_qos_policy_name":        name,
						"inter_region_traffic_qos_policy_description": name,
						"bandwidth_guarantee_mode":                    "byBandwidthPercent",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inter_region_traffic_qos_policy_name":        name + "_update",
					"inter_region_traffic_qos_policy_description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inter_region_traffic_qos_policy_name":        name + "_update",
						"inter_region_traffic_qos_policy_description": name + "_update",
					}),
				),
			},
		},
	})
}

func testAccCheckCenInterRegionTrafficQosPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
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

var resourceAlicloudCenInterRegionTrafficQosPolicyMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudCenInterRegionTrafficQosPolicyBasicDependence(name string) string {
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
`, name, configAlternateRegionProvider("cn-beijing"))
}
