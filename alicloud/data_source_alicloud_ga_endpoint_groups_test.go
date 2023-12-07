package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudGaEndpointGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_endpoint_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_endpoint_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_endpoint_group.default.name}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_endpoint_group.default.name}_fake"`,
		}),
	}
	listenerIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_endpoint_group.default.listener_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"listener_id": `"${alicloud_ga_endpoint_group.default.listener_id}_fake"`,
		}),
	}
	endpointGroupTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"endpoint_group_type": `"default"`,
		}),
		fakeConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"endpoint_group_type": `"virtual"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"status": `"active"`,
		}),
		fakeConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"status": `"configuring"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"ids":                 `["${alicloud_ga_endpoint_group.default.id}"]`,
			"name_regex":          `"${alicloud_ga_endpoint_group.default.name}"`,
			"listener_id":         `"${alicloud_ga_endpoint_group.default.listener_id}"`,
			"endpoint_group_type": `"default"`,
			"status":              `"active"`,
		}),
		fakeConfig: testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand, map[string]string{
			"ids":                 `["${alicloud_ga_endpoint_group.default.id}_fake"]`,
			"name_regex":          `"${alicloud_ga_endpoint_group.default.name}_fake"`,
			"listener_id":         `"${alicloud_ga_endpoint_group.default.listener_id}_fake"`,
			"endpoint_group_type": `"virtual"`,
			"status":              `"configuring"`,
		}),
	}
	var existAliCloudGaEndpointGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"groups.#":                               "1",
			"groups.0.id":                            CHECKSET,
			"groups.0.endpoint_group_id":             CHECKSET,
			"groups.0.listener_id":                   CHECKSET,
			"groups.0.endpoint_group_region":         CHECKSET,
			"groups.0.name":                          CHECKSET,
			"groups.0.description":                   CHECKSET,
			"groups.0.health_check_interval_seconds": CHECKSET,
			"groups.0.health_check_path":             CHECKSET,
			"groups.0.health_check_port":             CHECKSET,
			"groups.0.health_check_protocol":         CHECKSET,
			"groups.0.threshold_count":               CHECKSET,
			"groups.0.traffic_percentage":            CHECKSET,
			"groups.0.endpoint_group_ip_list.#":      CHECKSET,
			"groups.0.status":                        "active",
			"groups.0.port_overrides.#":              "1",
			"groups.0.endpoint_configurations.#":     "1",
		}
	}
	var fakeAliCloudGaEndpointGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}
	var alicloudGaEndpointGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_endpoint_groups.default",
		existMapFunc: existAliCloudGaEndpointGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudGaEndpointGroupsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGaEndpointGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, listenerIdConf, endpointGroupTypeConf, statusConf, allConf)
}

func testAccCheckAliCloudGaEndpointGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaEndpointGroup-%d"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
	}

	resource "alicloud_ga_bandwidth_package" "default" {
  		bandwidth              = 100
  		type                   = "Basic"
  		bandwidth_type         = "Basic"
  		payment_type           = "PayAsYouGo"
  		billing_type           = "PayBy95"
  		ratio                  = 30
  		bandwidth_package_name = var.name
  		auto_pay               = true
  		auto_use_coupon        = true
	}

	resource "alicloud_ga_bandwidth_package_attachment" "default" {
		accelerator_id       = data.alicloud_ga_accelerators.default.ids.0
		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}

	resource "alicloud_ga_listener" "default" {
		accelerator_id  = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
		client_affinity = "SOURCE_IP"
		protocol        = "UDP"
		name            = var.name
  		port_ranges {
    		from_port = "60"
    		to_port   = "70"
  		}
	}

	resource "alicloud_eip_address" "default" {
  		bandwidth            = "10"
  		internet_charge_type = "PayByBandwidth"
  		address_name         = var.name
	}

	resource "alicloud_ga_endpoint_group" "default" {
  		accelerator_id = alicloud_ga_listener.default.accelerator_id
		listener_id        = alicloud_ga_listener.default.id
  		description                   = var.name
  		name                          = var.name
  		threshold_count               = 4
  		traffic_percentage = 20
  		endpoint_group_region         = "%s"
  		health_check_interval_seconds = "3"
  		health_check_path             = "/healthcheck"
  		health_check_port             = "9999"
  		health_check_protocol         = "http"
  		port_overrides {
    		endpoint_port = "10"
    		listener_port = "60"
  		}
  		endpoint_configurations {
    		endpoint = alicloud_eip_address.default.ip_address
    		type     = "PublicIp"
    		weight   = "20"
  		}
	}

	data "alicloud_ga_endpoint_groups" "default" {
  		accelerator_id = alicloud_ga_endpoint_group.default.accelerator_id
		%s
	}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
