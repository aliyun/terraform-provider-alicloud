package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudConfigAggregateDeliveriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregate_delivery.default.delivery_channel_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregate_delivery.default.delivery_channel_name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_aggregate_delivery.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_aggregate_delivery.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_aggregate_delivery.default.id}"]`,
			"status": `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_aggregate_delivery.default.id}"]`,
			"status": `"0"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_config_aggregate_delivery.default.id}"]`,
			"name_regex": `"${alicloud_config_aggregate_delivery.default.delivery_channel_name}"`,
			"status":     `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_config_aggregate_delivery.default.id}_fake"]`,
			"name_regex": `"${alicloud_config_aggregate_delivery.default.delivery_channel_name}_fake"`,
			"status":     `"0"`,
		}),
	}
	var existAlicloudConfigAggregateDeliveriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                    "1",
			"ids.#":                      "1",
			"deliveries.#":               "1",
			"deliveries.0.id":            CHECKSET,
			"deliveries.0.account_id":    CHECKSET,
			"deliveries.0.aggregator_id": CHECKSET,
			"deliveries.0.delivery_channel_assume_role_arn":       CHECKSET,
			"deliveries.0.configuration_item_change_notification": "true",
			"deliveries.0.configuration_snapshot":                 "false",
			"deliveries.0.delivery_channel_condition":             "",
			"deliveries.0.delivery_channel_id":                    CHECKSET,
			"deliveries.0.delivery_channel_name":                  fmt.Sprintf("tf-testaccaggregatedelivery-%d", rand),
			"deliveries.0.delivery_channel_target_arn":            CHECKSET,
			"deliveries.0.delivery_channel_type":                  "SLS",
			"deliveries.0.description":                            fmt.Sprintf("tf-testaccaggregatedelivery-%d", rand),
			"deliveries.0.non_compliant_notification":             "true",
			"deliveries.0.oversized_data_oss_target_arn":          "",
			"deliveries.0.status":                                 "1",
		}
	}
	var fakeAlicloudConfigAggregateDeliveriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudConfigAggregateDeliveriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_config_aggregate_deliveries.default",
		existMapFunc: existAlicloudConfigAggregateDeliveriesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudConfigAggregateDeliveriesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
	}
	alicloudConfigAggregateDeliveriesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testaccaggregatedelivery-%d"
}
data "alicloud_account" "this" {}
data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}
resource "alicloud_config_aggregator" "default" {
	aggregator_accounts {
		account_id   =  data.alicloud_resource_manager_accounts.default.accounts.1.account_id
		account_name =  data.alicloud_resource_manager_accounts.default.accounts.1.display_name
		account_type = "ResourceDirectory"
	}
	aggregator_name = var.name
	description = var.name
}
locals {
  uid = data.alicloud_account.this.id
  sls = format("acs:log:%[2]s:%%s:project/%%s/logstore/%%s", local.uid, alicloud_log_project.this.name, alicloud_log_store.this.name)
}
resource "alicloud_log_project" "this" {
  name = var.name
}
resource "alicloud_log_store" "this" {
  name    = var.name
  project = alicloud_log_project.this.name
}
resource "alicloud_config_aggregate_delivery" "default" {
  aggregator_id                          = alicloud_config_aggregator.default.id
  configuration_item_change_notification = true
  delivery_channel_name                  = var.name
  delivery_channel_target_arn            = local.sls
  delivery_channel_type                  = "SLS"
  description                            = var.name
  non_compliant_notification             = true
}

data "alicloud_config_aggregate_deliveries" "default" {
  aggregator_id = alicloud_config_aggregator.default.id
  %s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
