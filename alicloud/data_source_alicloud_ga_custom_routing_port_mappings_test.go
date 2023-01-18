package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaCustomRoutingPortMappingsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	listenerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint.default.listener_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint.default.listener_id}_fake"`,
		}),
	}
	endpointGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"status": `"allow"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"status": `"deny"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint.default.listener_id}"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}"`,
			"status":            `"allow"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint.default.listener_id}_fake"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}_fake"`,
			"status":            `"deny"`,
		}),
	}
	var existAlicloudGaCustomRoutingPortMappingsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"custom_routing_port_mappings.#":                                         CHECKSET,
			"custom_routing_port_mappings.0.accelerator_id":                          CHECKSET,
			"custom_routing_port_mappings.0.listener_id":                             CHECKSET,
			"custom_routing_port_mappings.0.endpoint_group_id":                       CHECKSET,
			"custom_routing_port_mappings.0.endpoint_id":                             CHECKSET,
			"custom_routing_port_mappings.0.accelerator_port":                        CHECKSET,
			"custom_routing_port_mappings.0.vswitch":                                 CHECKSET,
			"custom_routing_port_mappings.0.endpoint_group_region":                   "cn-hangzhou",
			"custom_routing_port_mappings.0.protocols.#":                             "1",
			"custom_routing_port_mappings.0.protocols.0":                             "TCP",
			"custom_routing_port_mappings.0.destination_socket_address.#":            "1",
			"custom_routing_port_mappings.0.destination_socket_address.0.ip_address": CHECKSET,
			"custom_routing_port_mappings.0.destination_socket_address.0.port":       CHECKSET,
			"custom_routing_port_mappings.0.status":                                  "allow",
		}
	}
	var fakeAlicloudGaCustomRoutingPortMappingsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"custom_routing_port_mappings.#": "0",
		}
	}
	var alicloudGaCustomRoutingPortMappingsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_custom_routing_port_mappings.default",
		existMapFunc: existAlicloudGaCustomRoutingPortMappingsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaCustomRoutingPortMappingsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGaCustomRoutingPortMappingsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, listenerIdConf, endpointGroupIdConf, statusConf, allConf)
}

func testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaCustomRoutingPortMappings-%d"
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

	resource "alicloud_ga_custom_routing_endpoint_group_destination" "default" {
  		endpoint_group_id = alicloud_ga_custom_routing_endpoint_group.default.id
  		protocols         = ["TCP"]
  		from_port         = 1
  		to_port           = 10
	}
	
	resource "alicloud_ga_custom_routing_endpoint" "default" {
  		endpoint_group_id          = alicloud_ga_custom_routing_endpoint_group_destination.default.endpoint_group_id
  		endpoint                   = data.alicloud_vswitches.default.ids.0
  		type                       = "PrivateSubNet"
  		traffic_to_endpoint_policy = "AllowAll"
	}

	data "alicloud_ga_custom_routing_port_mappings" "default" {
  		accelerator_id = alicloud_ga_custom_routing_endpoint.default.accelerator_id
		%s
	}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
