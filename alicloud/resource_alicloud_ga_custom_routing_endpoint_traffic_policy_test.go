package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaCustomRoutingEndpointTrafficPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_custom_routing_endpoint_traffic_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudGaCustomRoutingEndpointTrafficPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaCustomRoutingEndpointTrafficPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sGaCustomRoutingEndpointTrafficPolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaCustomRoutingEndpointTrafficPolicyBasicDependence0)
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
					"endpoint_id": "${alicloud_ga_custom_routing_endpoint.default.custom_routing_endpoint_id}",
					"address":     "192.168.192.2",
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "1",
							"to_port":   "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_id":             CHECKSET,
						"address":                 "192.168.192.2",
						"port_ranges.#":           "1",
						"port_ranges.0.from_port": "1",
						"port_ranges.0.to_port":   "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address": "192.168.192.6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address": "192.168.192.6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_ranges": []map[string]interface{}{
						{
							"from_port": "2",
							"to_port":   "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_ranges.#":           "1",
						"port_ranges.0.from_port": "2",
						"port_ranges.0.to_port":   "6",
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

var AlicloudGaCustomRoutingEndpointTrafficPolicyMap = map[string]string{
	"accelerator_id":    CHECKSET,
	"listener_id":       CHECKSET,
	"endpoint_group_id": CHECKSET,
	"custom_routing_endpoint_traffic_policy_id": CHECKSET,
	"status": CHECKSET,
}

func AlicloudGaCustomRoutingEndpointTrafficPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING-Ipv6"
	}

	data "alicloud_vswitches" "default" {
  		name_regex = "default-zone-g_28"
  		vpc_id     = data.alicloud_vpcs.default.ids.0
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth      = 100
  		type           = "Basic"
  		bandwidth_type = "Basic"
  		payment_type   = "PayAsYouGo"
  		billing_type   = "PayBy95"
  		ratio          = 30
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		accelerator_id       = data.alicloud_ga_accelerators.default.accelerators.0.id
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  		listener_type  = "CustomRouting"
  		port_ranges {
    		from_port = 10000
    		to_port   = 26000
  		}
	}

	resource "alicloud_ga_custom_routing_endpoint_group" "default" {
  		accelerator_id                     = alicloud_ga_listener.default.accelerator_id
  		listener_id                        = alicloud_ga_listener.default.id
  		endpoint_group_region              = "%s"
  		custom_routing_endpoint_group_name = var.name
  		description                        = var.name
	}

	resource "alicloud_ga_custom_routing_endpoint_group_destination" "default" {
  		endpoint_group_id = alicloud_ga_custom_routing_endpoint_group.default.id
  		protocols         = ["tcp"]
  		from_port         = 1
  		to_port           = 10
	}

	resource "alicloud_ga_custom_routing_endpoint" "default" {
  		endpoint_group_id          = alicloud_ga_custom_routing_endpoint_group_destination.default.endpoint_group_id
  		endpoint                   = data.alicloud_vswitches.default.ids.0
  		type                       = "PrivateSubNet"
  		traffic_to_endpoint_policy = "AllowAll"
	}
`, name, defaultRegionToTest)
}
