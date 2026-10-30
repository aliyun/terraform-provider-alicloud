package ga_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudGaCustomRoutingEndpointGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_custom_routing_endpoint_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_custom_routing_endpoint_group.default.custom_routing_endpoint_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_custom_routing_endpoint_group.default.custom_routing_endpoint_group_name}_fake"`,
		}),
	}
	listenerIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint_group.default.listener_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint_group.default.listener_id}_fake"`,
		}),
	}
	endpointGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group.default.id}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"status": `"active"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"status": `"deleting"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ga_custom_routing_endpoint_group.default.id}"]`,
			"name_regex":        `"${alicloud_ga_custom_routing_endpoint_group.default.custom_routing_endpoint_group_name}"`,
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint_group.default.listener_id}"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint_group.default.id}"`,
			"status":            `"active"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingEndpointGroupsDataSourceName(rand, map[string]string{
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
			"groups.0.id":                       acctest.CHECKSET,
			"groups.0.endpoint_group_id":        acctest.CHECKSET,
			"groups.0.accelerator_id":           acctest.CHECKSET,
			"groups.0.listener_id":              acctest.CHECKSET,
			"groups.0.endpoint_group_region":    acctest.DefaultRegionToTest,
			"groups.0.endpoint_group_ip_list.#": acctest.CHECKSET,
			"groups.0.endpoint_group_unconfirmed_ip_list.#": acctest.NOSET,
			"groups.0.custom_routing_endpoint_group_name":   acctest.CHECKSET,
			"groups.0.description":                          acctest.CHECKSET,
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
	var alicloudGaCustomRoutingEndpointGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ga_custom_routing_endpoint_groups.default",
		ExistMapFunc: existAlicloudGaCustomRoutingEndpointGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudGaCustomRoutingEndpointGroupsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudGaCustomRoutingEndpointGroupsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, listenerIdConf, endpointGroupIdConf, statusConf, allConf)
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
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
