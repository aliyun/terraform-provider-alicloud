package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// Skipped: The resource and its apis have deprecated
func SkipTestAccAliCloudConfigDeliveryChannelsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_config_delivery_channels.example"

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_delivery_channel.example.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_delivery_channel.example.id}_fake"]`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"status":     `1`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"status":     `0`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"ids":        `["${alicloud_config_delivery_channel.example.id}"]`,
			"status":     `1`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveryChannelsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}_fake"`,
			"ids":        `["${alicloud_config_delivery_channel.example.id}_fake"]`,
			"status":     `0`,
		}),
	}

	var ExistConfigDeliveryChannelsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"channels.#": "1",
			"names.#":    "1",
			"ids.#":      "1",
			"channels.0.delivery_channel_assume_role_arn": acctest.CHECKSET,
			"channels.0.delivery_channel_condition":       acctest.CHECKSET,
			"channels.0.id":                               acctest.CHECKSET,
			"channels.0.delivery_channel_id":              acctest.CHECKSET,
			"channels.0.delivery_channel_name":            fmt.Sprintf("tf-testAccConfigDeliveryChannels%d", rand),
			"channels.0.delivery_channel_target_arn":      acctest.CHECKSET,
			"channels.0.delivery_channel_type":            "MNS",
			"channels.0.description":                      fmt.Sprintf("tf-testAccConfigDeliveryChannels%d", rand),
			"channels.0.status":                           "1",
		}
	}

	var FakeConfigDeliveryChannelsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"channels.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var configDeliveryChannelsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: ExistConfigDeliveryChannelsRecordsMapFunc,
		FakeMapFunc:  FakeConfigDeliveryChannelsRecordsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
	}

	configDeliveryChannelsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)

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
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}

// TestAccAliCloudConfigDeliveryChannelsDataSource exercises the read path of the (deprecated
// but still functional) alicloud_config_delivery_channels data source. The data source lists
// delivery channels via the Config 2019-01-08 DescribeDeliveryChannels API, so the fixture
// creates the channel through the deprecated alicloud_config_delivery_channel resource, whose
// PutDeliveryChannel create and DescribeDeliveryChannels read share that same 2019 API version.
// A channel created via the non-deprecated alicloud_config_delivery resource (2020-09-07 API)
// is not visible to this 2019 list API, which is why the fixture pairs the deprecated resource
// with the deprecated data source. An MNS topic target is used: MNS topics are lightweight and
// avoid the OSS "bucket not exists" and SLS "project quota 300" failures that the OSS/SLS
// delivery fixtures hit on this CI account, so the read path is covered without extra quota
// dependencies. The delivery_channel_assume_role_arn references the existing
// AliyunServiceRoleForConfig service-linked role (read-only lookup, no role creation), which is
// present once Cloud Config is enabled on the account.
func TestAccAliCloudConfigDeliveryChannelsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudConfigSupportedRegions)
	resourceId := "data.alicloud_config_delivery_channels.example"

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}_fake"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_delivery_channel.example.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_delivery_channel.example.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"status":     `1`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"status":     `0`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}"`,
			"ids":        `["${alicloud_config_delivery_channel.example.id}"]`,
			"status":     `1`,
		}),
		FakeConfig: testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_delivery_channel.example.delivery_channel_name}_fake"`,
			"ids":        `["${alicloud_config_delivery_channel.example.id}_fake"]`,
			"status":     `0`,
		}),
	}

	var ExistConfigDeliveryChannelsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"channels.#":                     "1",
			"names.#":                        "1",
			"ids.#":                          "1",
			"channels.0.id":                  acctest.CHECKSET,
			"channels.0.delivery_channel_id": acctest.CHECKSET,
			"channels.0.delivery_channel_assume_role_arn": acctest.CHECKSET,
			"channels.0.delivery_channel_name":            fmt.Sprintf("tf-testaccconfigdeliverychannels%d", rand),
			"channels.0.delivery_channel_target_arn":      acctest.CHECKSET,
			"channels.0.delivery_channel_type":            "MNS",
			"channels.0.description":                      fmt.Sprintf("tf-testaccconfigdeliverychannels%d", rand),
			"channels.0.status":                           "1",
		}
	}

	var FakeConfigDeliveryChannelsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"channels.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var configDeliveryChannelsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: ExistConfigDeliveryChannelsMapFunc,
		FakeMapFunc:  FakeConfigDeliveryChannelsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
	}

	configDeliveryChannelsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)
}

func testAccCheckAlicloudConfigDeliveryChannelsMNSSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testaccconfigdeliverychannels%d"
}
data "alicloud_account" "this" {}
data "alicloud_ram_roles" "this" {
  name_regex = "^AliyunServiceRoleForConfig$"
}
locals {
  uid       = data.alicloud_account.this.id
  role_arn  = data.alicloud_ram_roles.this.roles.0.arn
  mns       = format("acs:mns:%[2]s:%%s:/topics/%%s", local.uid, alicloud_mns_topic.default.name)
}
resource "alicloud_mns_topic" "default" {
  name = var.name
}
resource "alicloud_config_delivery_channel" "example" {
  delivery_channel_assume_role_arn = local.role_arn
  delivery_channel_target_arn      = local.mns
  delivery_channel_type            = "MNS"
  delivery_channel_name            = var.name
  description                      = var.name
  status                           = 1
}

data "alicloud_config_delivery_channels" "example" {
%s
}
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
