package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaCustomRoutingEndpointGroupDestination_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_custom_routing_endpoint_group_destination.default"
	ra := resourceAttrInit(resourceId, AlicloudGaCustomRoutingEndpointGroupDestinationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaCustomRoutingEndpointGroupDestination")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sGaCustomRoutingEndpointGroupDestination%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaCustomRoutingEndpointGroupDestinationBasicDependence0)
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
					"endpoint_group_id": "${alicloud_ga_custom_routing_endpoint_group.default.id}",
					"protocols":         []string{"tcp"},
					"from_port":         "1",
					"to_port":           "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint_group_id": CHECKSET,
						"protocols.#":       "1",
						"protocols.0":       "tcp",
						"from_port":         "1",
						"to_port":           "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{"udp"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "1",
						"protocols.0": "udp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocols": []string{"tcp", "udp"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocols.#": "2",
						"protocols.0": "tcp",
						"protocols.1": "udp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"from_port": "2",
					"to_port":   "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"from_port": "2",
						"to_port":   "3",
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

var AlicloudGaCustomRoutingEndpointGroupDestinationMap = map[string]string{
	"accelerator_id": CHECKSET,
	"listener_id":    CHECKSET,
	"custom_routing_endpoint_group_destination_id": CHECKSET,
	"status": CHECKSET,
}

func AlicloudGaCustomRoutingEndpointGroupDestinationBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}
	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
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
    		to_port   = 16000
  		}
	}
	resource "alicloud_ga_custom_routing_endpoint_group" "default" {
  		accelerator_id                     = alicloud_ga_listener.default.accelerator_id
  		listener_id                        = alicloud_ga_listener.default.id
  		endpoint_group_region              = "%s"
  		custom_routing_endpoint_group_name = var.name
  		description                        = var.name
	}
`, name, defaultRegionToTest)
}
