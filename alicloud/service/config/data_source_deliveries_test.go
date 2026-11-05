package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudConfigDeliveriesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery.default.delivery_channel_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery.default.delivery_channel_name}_fake"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_delivery.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_delivery.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_delivery.default.id}"]`,
			"status": `"1"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_delivery.default.id}"]`,
			"status": `"0"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery.default.delivery_channel_name}"`,
			"ids":        `["${alicloud_config_delivery.default.id}"]`,
			"status":     `"1"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveriesDataSourceName(rand, map[string]string{
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
			"deliveries.0.account_id":                             acctest.CHECKSET,
			"deliveries.0.delivery_channel_id":                    acctest.CHECKSET,
			"deliveries.0.delivery_channel_assume_role_arn":       acctest.CHECKSET,
			"deliveries.0.delivery_channel_condition":             "",
			"deliveries.0.oversized_data_oss_target_arn":          "",
			"deliveries.0.id":                                     acctest.CHECKSET,
			"deliveries.0.delivery_channel_name":                  fmt.Sprintf("tf-testaccconfigdelivery%d", rand),
			"deliveries.0.delivery_channel_target_arn":            acctest.CHECKSET,
			"deliveries.0.delivery_channel_type":                  "MNS",
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
	var alicloudConfigDeliveriesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_config_deliveries.default",
		ExistMapFunc: existAlicloudConfigDeliveriesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudConfigDeliveriesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudConfigDeliveriesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)
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
  mns	       	   = format("acs:mns:%[2]s:%%s:/topics/%%s",local.uid,alicloud_mns_topic.this.name)
}
resource "alicloud_mns_topic" "this" {
  name = var.name
}
resource "alicloud_config_delivery" "default" {
	configuration_item_change_notification = true
	non_compliant_notification = true
	delivery_channel_name = var.name
	delivery_channel_target_arn = local.mns
	delivery_channel_type = "MNS"
	description = var.name
}

data "alicloud_config_deliveries" "default" {	
	%s
}
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
