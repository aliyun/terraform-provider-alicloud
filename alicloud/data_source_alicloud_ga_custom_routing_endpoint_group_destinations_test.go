package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaCustomRoutingEndpointGroupDestinationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint_group_destination.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint_group_destination.default.id}_fake"]`,
		}),
	}
	listenerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint_group_destination.default.listener_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint_group_destination.default.listener_id}_fake"`,
		}),
	}
	endpointGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group_destination.default.endpoint_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group_destination.default.endpoint_group_id}_fake"`,
		}),
	}
	protocolsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"protocols": `["TCP"]`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"protocols": `["UDP"]`,
		}),
	}
	fromPortConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"from_port": `1`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"from_port": `2`,
		}),
	}
	toPortConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"to_port": `2`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"to_port": `1`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint_group_destination.default.id}"]`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint_group_destination.default.listener_id}"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group_destination.default.endpoint_group_id}"`,
			"protocols":         `["TCP"]`,
			"from_port":         `1`,
			"to_port":           `2`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint_group_destination.default.id}_fake"]`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint_group_destination.default.listener_id}_fake"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group_destination.default.endpoint_group_id}_fake"`,
			"protocols":         `["UDP"]`,
			"from_port":         `2`,
			"to_port":           `1`,
		}),
	}
	var existAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "1",
			"custom_routing_endpoint_group_destinations.#":                                              "1",
			"custom_routing_endpoint_group_destinations.0.id":                                           CHECKSET,
			"custom_routing_endpoint_group_destinations.0.endpoint_group_id":                            CHECKSET,
			"custom_routing_endpoint_group_destinations.0.custom_routing_endpoint_group_destination_id": CHECKSET,
			"custom_routing_endpoint_group_destinations.0.accelerator_id":                               CHECKSET,
			"custom_routing_endpoint_group_destinations.0.listener_id":                                  CHECKSET,
			"custom_routing_endpoint_group_destinations.0.protocols.#":                                  "1",
			"custom_routing_endpoint_group_destinations.0.protocols.0":                                  "TCP",
			"custom_routing_endpoint_group_destinations.0.from_port":                                    "1",
			"custom_routing_endpoint_group_destinations.0.to_port":                                      "2",
		}
	}
	var fakeAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
			"custom_routing_endpoint_group_destinations.#": "0",
		}
	}
	var alicloudGaCustomRoutingEndpointGroupDestinationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_custom_routing_endpoint_group_destinations.default",
		existMapFunc: existAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGaCustomRoutingEndpointGroupDestinationsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, listenerIdConf, endpointGroupIdConf, protocolsConf, fromPortConf, toPortConf, allConf)
}

func testAccCheckAlicloudGaCustomRoutingEndpointGroupDestinationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaCustomRoutingEndpointGroupDestination-%d"
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

	resource "alicloud_ga_custom_routing_endpoint_group_destination" "default" {
  		endpoint_group_id = alicloud_ga_custom_routing_endpoint_group.default.id
  		protocols         = ["TCP"]
  		from_port         = 1
  		to_port           = 2
	}

	data "alicloud_ga_custom_routing_endpoint_group_destinations" "default" {
  		accelerator_id = alicloud_ga_custom_routing_endpoint_group_destination.default.accelerator_id
		%s
	}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
