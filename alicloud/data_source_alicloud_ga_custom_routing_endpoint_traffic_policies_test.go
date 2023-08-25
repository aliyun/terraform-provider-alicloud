package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint_traffic_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint_traffic_policy.default.id}_fake"]`,
		}),
	}
	listenerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.listener_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.listener_id}_fake"`,
		}),
	}
	endpointGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.endpoint_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.endpoint_group_id}_fake"`,
		}),
	}
	endpointIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"endpoint_id": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.endpoint_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"endpoint_id": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.endpoint_id}_fake"`,
		}),
	}
	addressConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"address": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.address}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"address": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.address}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint_traffic_policy.default.id}"]`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.listener_id}"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.endpoint_group_id}"`,
			"endpoint_id":       `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.endpoint_id}"`,
			"address":           `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.address}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint_traffic_policy.default.id}_fake"]`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.listener_id}_fake"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.endpoint_group_id}_fake"`,
			"endpoint_id":       `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.endpoint_id}_fake"`,
			"address":           `"${alicloud_ga_custom_routing_endpoint_traffic_policy.default.address}_fake"`,
		}),
	}
	var existAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "1",
			"custom_routing_endpoint_traffic_policies.#":                                           "1",
			"custom_routing_endpoint_traffic_policies.0.id":                                        CHECKSET,
			"custom_routing_endpoint_traffic_policies.0.endpoint_id":                               CHECKSET,
			"custom_routing_endpoint_traffic_policies.0.custom_routing_endpoint_traffic_policy_id": CHECKSET,
			"custom_routing_endpoint_traffic_policies.0.accelerator_id":                            CHECKSET,
			"custom_routing_endpoint_traffic_policies.0.listener_id":                               CHECKSET,
			"custom_routing_endpoint_traffic_policies.0.endpoint_group_id":                         CHECKSET,
			"custom_routing_endpoint_traffic_policies.0.address":                                   "192.168.192.2",
			"custom_routing_endpoint_traffic_policies.0.port_ranges.#":                             "1",
			"custom_routing_endpoint_traffic_policies.0.port_ranges.0.from_port":                   "1",
			"custom_routing_endpoint_traffic_policies.0.port_ranges.0.to_port":                     "2",
		}
	}
	var fakeAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
			"custom_routing_endpoint_traffic_policies.#": "0",
		}
	}
	var alicloudGaCustomRoutingEndpointTrafficPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_custom_routing_endpoint_traffic_policies.default",
		existMapFunc: existAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGaCustomRoutingEndpointTrafficPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, listenerIdConf, endpointGroupIdConf, endpointIdConf, addressConf, allConf)
}

func testAccCheckAliCloudGaCustomRoutingEndpointTrafficPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaCustomRoutingEndpointTrafficPolicy-%d"
	}

	data "alicloud_zones" "default" {
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
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

	resource "alicloud_ga_custom_routing_endpoint_traffic_policy" "default" {
  		endpoint_id = alicloud_ga_custom_routing_endpoint.default.custom_routing_endpoint_id
  		address     = "192.168.192.2"
  		port_ranges {
    		from_port = 1
    		to_port   = 2
  		}
	}

	data "alicloud_ga_custom_routing_endpoint_traffic_policies" "default" {
  		accelerator_id = alicloud_ga_custom_routing_endpoint_traffic_policy.default.accelerator_id
		%s
	}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
