package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudConfigDeliveriesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery.default.delivery_channel_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery.default.delivery_channel_name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_delivery.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_delivery.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_delivery.default.id}"]`,
			"status": `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_delivery.default.id}"]`,
			"status": `"0"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery.default.delivery_channel_name}"`,
			"ids":        `["${alicloud_config_delivery.default.id}"]`,
			"status":     `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery.default.delivery_channel_name}_fake"`,
			"ids":        `["${alicloud_config_delivery.default.id}_fake"]`,
			"status":     `"0"`,
		}),
	}
	var existAlicloudConfigDeliveriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "1",
			"names.#":      "1",
			"deliveries.#": "1",
			"deliveries.0.configuration_item_change_notification": "true",
			"deliveries.0.configuration_snapshot":                 "false",
			"deliveries.0.account_id":                             CHECKSET,
			"deliveries.0.delivery_channel_id":                    CHECKSET,
			"deliveries.0.delivery_channel_assume_role_arn":       CHECKSET,
			"deliveries.0.delivery_channel_condition":             "",
			"deliveries.0.oversized_data_oss_target_arn":          "",
			"deliveries.0.id":                                     CHECKSET,
			"deliveries.0.delivery_channel_name":                  fmt.Sprintf("tf-testaccconfigdelivery%d", rand),
			"deliveries.0.delivery_channel_target_arn":            CHECKSET,
			"deliveries.0.delivery_channel_type":                  "SLS",
			"deliveries.0.description":                            fmt.Sprintf("tf-testaccconfigdelivery%d", rand),
			"deliveries.0.non_compliant_notification":             "true",
			"deliveries.0.status":                                 "1",
		}
	}
	var fakeAlicloudConfigDeliveriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudConfigDeliveriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_config_deliveries.default",
		existMapFunc: existAlicloudConfigDeliveriesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudConfigDeliveriesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudConfigDeliveriesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudConfigDeliveriesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccconfigdelivery%d"
}
data "alicloud_account" "this" {}
locals {
  uid          	   = data.alicloud_account.this.id
  sls	       	   = format("acs:log:%[2]s:%%s:project/%%s/logstore/%%s",local.uid,alicloud_log_project.this.name,alicloud_log_store.this.name)
}
resource "alicloud_log_project" "this" {
  name = var.name
}
resource "alicloud_log_store" "this" {
  name = var.name
  project = alicloud_log_project.this.name
}
resource "alicloud_config_delivery" "default" {
	configuration_item_change_notification = true
	non_compliant_notification = true
	delivery_channel_name = var.name
	delivery_channel_target_arn = local.sls
	delivery_channel_type = "SLS"
	description = var.name
}

data "alicloud_config_deliveries" "default" {	
	%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
