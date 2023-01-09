package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaCustomRoutingEndpointGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_custom_routing_endpoint_group.default.custom_routing_endpoint_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_custom_routing_endpoint_group.default.custom_routing_endpoint_group_name}_fake"`,
		}),
	}
	listenerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint_group.default.listener_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint_group.default.listener_id}_fake"`,
		}),
	}
	endpointGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group.default.id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"status": `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"status": `"deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint_group.default.id}"]`,
			"name_regex":        `"${alicloud_ga_custom_routing_endpoint_group.default.custom_routing_endpoint_group_name}"`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint_group.default.listener_id}"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group.default.id}"`,
			"status":            `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint_group.default.id}_fake"]`,
			"name_regex":        `"${alicloud_ga_custom_routing_endpoint_group.default.custom_routing_endpoint_group_name}_fake"`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint_group.default.listener_id}_fake"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group.default.id}_fake"`,
			"status":            `"deleting"`,
		}),
	}
	var existAlicloudGaCustomRoutingEndpointGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"groups.#":                          "1",
			"groups.0.id":                       CHECKSET,
			"groups.0.endpoint_group_id":        CHECKSET,
			"groups.0.accelerator_id":           CHECKSET,
			"groups.0.listener_id":              CHECKSET,
			"groups.0.endpoint_group_region":    defaultRegionToTest,
			"groups.0.endpoint_group_ip_list.#": CHECKSET,
			"groups.0.endpoint_group_unconfirmed_ip_list.#": NOSET,
			"groups.0.custom_routing_endpoint_group_name":   CHECKSET,
			"groups.0.description":                          CHECKSET,
			"groups.0.status":                               "active",
		}
	}
	var fakeAlicloudGaCustomRoutingEndpointGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}
	var alicloudGaCustomRoutingEndpointGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_custom_routing_endpoint_groups.default",
		existMapFunc: existAlicloudGaCustomRoutingEndpointGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaCustomRoutingEndpointGroupsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGaCustomRoutingEndpointGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, listenerIdConf, endpointGroupIdConf, statusConf, allConf)
}

func testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaCustomRoutingEndpointGroup-%d"
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
	
	data "alicloud_ga_custom_routing_endpoint_groups" "default" {
  		accelerator_id = alicloud_ga_custom_routing_endpoint_group.default.accelerator_id
		%s
	}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
