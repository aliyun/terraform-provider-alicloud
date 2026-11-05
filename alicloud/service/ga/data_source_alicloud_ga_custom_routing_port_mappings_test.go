package ga_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudGaCustomRoutingPortMappingsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	listenerIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint.default.listener_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_custom_routing_endpoint.default.listener_id}_fake"`,
		}),
	}
	endpointGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"status": `"allow"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"status": `"deny"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint.default.listener_id}"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}"`,
			"status":            `"allow"`,
		}),
		FakeConfig: testAccCheckAlicloudGaCustomRoutingPortMappingsDataSourceName(rand, map[string]string{
			"listener_id":       `"${alicloud_ga_custom_routing_endpoint.default.listener_id}_fake"`,
			"endpoint_group_id": `"${alicloud_ga_custom_routing_endpoint.default.endpoint_group_id}_fake"`,
			"status":            `"deny"`,
		}),
	}
	var existAlicloudGaCustomRoutingPortMappingsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"custom_routing_port_mappings.#":                                         acctest.CHECKSET,
			"custom_routing_port_mappings.0.accelerator_id":                          acctest.CHECKSET,
			"custom_routing_port_mappings.0.listener_id":                             acctest.CHECKSET,
			"custom_routing_port_mappings.0.endpoint_group_id":                       acctest.CHECKSET,
			"custom_routing_port_mappings.0.endpoint_id":                             acctest.CHECKSET,
			"custom_routing_port_mappings.0.accelerator_port":                        acctest.CHECKSET,
			"custom_routing_port_mappings.0.vswitch":                                 acctest.CHECKSET,
			"custom_routing_port_mappings.0.endpoint_group_region":                   "cn-hangzhou",
			"custom_routing_port_mappings.0.protocols.#":                             "1",
			"custom_routing_port_mappings.0.protocols.0":                             "TCP",
			"custom_routing_port_mappings.0.destination_socket_address.#":            "1",
			"custom_routing_port_mappings.0.destination_socket_address.0.ip_address": acctest.CHECKSET,
			"custom_routing_port_mappings.0.destination_socket_address.0.port":       acctest.CHECKSET,
			"custom_routing_port_mappings.0.status":                                  "allow",
		}
	}
	var fakeAlicloudGaCustomRoutingPortMappingsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"custom_routing_port_mappings.#": "0",
		}
	}
	var alicloudGaCustomRoutingPortMappingsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ga_custom_routing_port_mappings.default",
		ExistMapFunc: existAlicloudGaCustomRoutingPortMappingsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudGaCustomRoutingPortMappingsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudGaCustomRoutingPortMappingsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, listenerIdConf, endpointGroupIdConf, statusConf, allConf)
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
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
