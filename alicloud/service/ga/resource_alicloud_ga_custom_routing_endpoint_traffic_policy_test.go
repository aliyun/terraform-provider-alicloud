package ga_test

// collection-order: bypass port_ranges -- pre-existing multi-member TypeList attribute left untouched by this migration; convergent reorder coverage is pre-existing debt outside this change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudGaCustomRoutingEndpointTrafficPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_custom_routing_endpoint_traffic_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaCustomRoutingEndpointTrafficPolicyMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaCustomRoutingEndpointTrafficPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sGaCustomRoutingEndpointTrafficPolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaCustomRoutingEndpointTrafficPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"endpoint_id":             acctest.CHECKSET,
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

var AliCloudGaCustomRoutingEndpointTrafficPolicyMap = map[string]string{
	"accelerator_id":    acctest.CHECKSET,
	"listener_id":       acctest.CHECKSET,
	"endpoint_group_id": acctest.CHECKSET,
	"custom_routing_endpoint_traffic_policy_id": acctest.CHECKSET,
	"status": acctest.CHECKSET,
}

func AliCloudGaCustomRoutingEndpointTrafficPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
        bandwidth_billing_type = "BandwidthPackage"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name       = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.ids.0
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
  		protocols         = ["TCP"]
  		from_port         = 1
  		to_port           = 10
	}

	resource "alicloud_ga_custom_routing_endpoint" "default" {
  		endpoint_group_id          = alicloud_ga_custom_routing_endpoint_group_destination.default.endpoint_group_id
  		endpoint                   = alicloud_vswitch.default.id
  		type                       = "PrivateSubNet"
  		traffic_to_endpoint_policy = "AllowAll"
	}
`, name, acctest.DefaultRegionToTest)
}
