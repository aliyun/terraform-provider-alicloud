package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// Skipped: The resource and its apis have deprecated
func SkipTestAccAliCloudConfigDeliveryChannelsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_config_delivery_channels.example"

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_delivery_channel.example.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_delivery_channel.example.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"status":     `1`,
		}),
		fakeConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"status":     `0`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"ids":        `["${alicloud_config_delivery_channel.example.id}"]`,
			"status":     `1`,
		}),
		fakeConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}_fake"`,
			"ids":        `["${alicloud_config_delivery_channel.example.id}_fake"]`,
			"status":     `0`,
		}),
	}

	var existConfigDeliveryChannelsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"channels.#": "1",
			"names.#":    "1",
			"ids.#":      "1",
			"channels.0.delivery_channel_assume_role_arn": CHECKSET,
			"channels.0.delivery_channel_condition":       CHECKSET,
			"channels.0.id":                               CHECKSET,
			"channels.0.delivery_channel_id":              CHECKSET,
			"channels.0.delivery_channel_name":            fmt.Sprintf("tf-testAccConfigDeliveryChannels%d", rand),
			"channels.0.delivery_channel_target_arn":      CHECKSET,
			"channels.0.delivery_channel_type":            "MNS",
			"channels.0.description":                      fmt.Sprintf("tf-testAccConfigDeliveryChannels%d", rand),
			"channels.0.status":                           "1",
		}
	}

	var fakeConfigDeliveryChannelsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"channels.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var configDeliveryChannelsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existConfigDeliveryChannelsRecordsMapFunc,
		fakeMapFunc:  fakeConfigDeliveryChannelsRecordsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
	}

	configDeliveryChannelsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)

}

func testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccConfigDeliveryChannels%d"
}

locals {
  uid          	   = data.alicloud_account.this.id
  role_arn         = data.alicloud_ram_roles.this.roles.0.arn
  mns	       	   = format("acs:oss:%[2]s:%%s:/topics/%%s",local.uid,alicloud_mns_topic.default.name)
}

resource "alicloud_mns_topic" "default" {
  name = var.name
}

data "alicloud_account" "this" {}

data "alicloud_ram_roles" "this" {
  name_regex = "^AliyunServiceRoleForConfig$"
}

resource "alicloud_config_delivery_channel" "example" {
  delivery_channel_assume_role_arn = local.role_arn
  delivery_channel_target_arn      = local.mns
  delivery_channel_type            = "MNS"
  description                      = var.name
  status                           = 1
  delivery_channel_name            = var.name
}

data "alicloud_config_delivery_channels" "example"{
%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
