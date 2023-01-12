package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaCustomRoutingEndpointsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint.default.id}_fake"]`,
		}),
	}
	listenerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint.default.listener_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint.default.listener_id}_fake"`,
		}),
	}
	endpointGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint.default.id}"]`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint.default.listener_id}"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint.default.id}_fake"]`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint.default.listener_id}_fake"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}_fake"`,
		}),
	}
	var existAlicloudGaCustomRoutingEndpointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"custom_routing_endpoints.#":    "1",
			"custom_routing_endpoints.0.id": CHECKSET,
			"custom_routing_endpoints.0.endpoint_group_id":          CHECKSET,
			"custom_routing_endpoints.0.custom_routing_endpoint_id": CHECKSET,
			"custom_routing_endpoints.0.accelerator_id":             CHECKSET,
			"custom_routing_endpoints.0.listener_id":                CHECKSET,
			"custom_routing_endpoints.0.endpoint":                   CHECKSET,
			"custom_routing_endpoints.0.type":                       "PrivateSubNet",
			"custom_routing_endpoints.0.traffic_to_endpoint_policy": "DenyAll",
		}
	}
	var fakeAlicloudGaCustomRoutingEndpointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "0",
			"custom_routing_endpoints.#": "0",
		}
	}
	var alicloudGaCustomRoutingEndpointsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_custom_routing_endpoints.default",
		existMapFunc: existAlicloudGaCustomRoutingEndpointsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaCustomRoutingEndpointsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGaCustomRoutingEndpointsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, listenerIdConf, endpointGroupIdConf, allConf)
}

func testAccCheckAlicloudGaCustomRoutingEndpointsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaCustomRoutingEndpoint-%d"
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

	resource "alicloud_ga_custom_routing_endpoint" "default" {
  		endpoint_group_id          = alicloud_ga_custom_routing_endpoint_group.default.id
  		endpoint                   = data.alicloud_vswitches.default.ids.0
  		type                       = "PrivateSubNet"
  		traffic_to_endpoint_policy = "DenyAll"
	}

	data "alicloud_ga_custom_routing_endpoints" "default" {
  		accelerator_id = alicloud_ga_custom_routing_endpoint.default.accelerator_id
		%s
	}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
