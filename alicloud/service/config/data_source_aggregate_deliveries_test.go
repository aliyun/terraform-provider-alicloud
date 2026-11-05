package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudConfigAggregateDeliveriesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregate_delivery.default.delivery_channel_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregate_delivery.default.delivery_channel_name}_fake"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_aggregate_delivery.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_aggregate_delivery.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_aggregate_delivery.default.id}"]`,
			"status": `"1"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_aggregate_delivery.default.id}"]`,
			"status": `"0"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_config_aggregate_delivery.default.id}"]`,
			"name_regex": `"${alicloud_config_aggregate_delivery.default.delivery_channel_name}"`,
			"status":     `"1"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregateDeliveriesDataSourceName(rand, map[string]string{
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
			"deliveries.0.id":            acctest.CHECKSET,
			"deliveries.0.account_id":    acctest.CHECKSET,
			"deliveries.0.aggregator_id": acctest.CHECKSET,
			"deliveries.0.delivery_channel_assume_role_arn":       acctest.CHECKSET,
			"deliveries.0.configuration_item_change_notification": "true",
			"deliveries.0.configuration_snapshot":                 "false",
			"deliveries.0.delivery_channel_condition":             "",
			"deliveries.0.delivery_channel_id":                    acctest.CHECKSET,
			"deliveries.0.delivery_channel_name":                  fmt.Sprintf("tf-testaccaggregatedelivery-%d", rand),
			"deliveries.0.delivery_channel_target_arn":            acctest.CHECKSET,
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
	var alicloudConfigAggregateDeliveriesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_config_aggregate_deliveries.default",
		ExistMapFunc: existAlicloudConfigAggregateDeliveriesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudConfigAggregateDeliveriesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckEnterpriseAccountEnabled(t)
	}
	alicloudConfigAggregateDeliveriesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)
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
  aggregator_id = alicloud_config_aggregate_delivery.default.aggregator_id
  %s
}
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
