package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDtsConsumerChannelsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	checkoutSupportedRegions(t, true, connectivity.DTSSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsConsumerChannelsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dts_consumer_channel.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDtsConsumerChannelsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dts_consumer_channel.default.id}_fake"]`,
		}),
	}
	var existAlicloudDtsConsumerChannelsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"channels.#":                          "1",
			"channels.0.consumer_group_name":      fmt.Sprintf("tf_testAcc%d", rand),
			"channels.0.consumer_group_user_name": fmt.Sprintf("tf_testAcc%d", rand),
			"channels.0.message_delay":            CHECKSET,
			"channels.0.id":                       CHECKSET,
			"channels.0.consumption_checkpoint":   "",
			"channels.0.consumer_group_id":        CHECKSET,
			"channels.0.unconsumed_data":          CHECKSET,
		}
	}
	var fakeAlicloudDtsConsumerChannelsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDtsConsumerChannelsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dts_consumer_channels.default",
		existMapFunc: existAlicloudDtsConsumerChannelsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDtsConsumerChannelsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudDtsConsumerChannelsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudDtsConsumerChannelsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf_testAcc%d"
}
variable "region" {
  default = "%s"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_db_zones" "default" {}

data "alicloud_db_instance_classes" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    =  data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage = "10"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  instance_name    = var.name
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.instance.id
  name        = join("", [var.database_name, count.index])
}

resource "alicloud_db_account" "account" {
  db_instance_id      = alicloud_db_instance.instance.id
  account_name        = var.database_name
  account_password    = var.password
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_instance.instance.id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}

resource "alicloud_dts_subscription_job" "default" {
  dts_job_name                       = var.name
  payment_type                       = "PayAsYouGo"
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = var.region
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.instance.id
  source_endpoint_database_name      = var.database_name
  source_endpoint_user_name          = var.database_name
  source_endpoint_password           = var.password
  subscription_instance_network_type = "vpc"
  db_list                            = "{\"dtstestdata\":{\"name\":\"tftestdatabase\",\"all\":true}}"
  subscription_instance_vpc_id       = data.alicloud_vpcs.default.ids[0]
  subscription_instance_vswitch_id   = data.alicloud_vswitches.default.ids[0]
  status                             = "Normal"
  depends_on = [alicloud_db_account_privilege.privilege]
}

resource "alicloud_dts_consumer_channel" "default" {
  dts_instance_id          = alicloud_dts_subscription_job.default.dts_instance_id
  consumer_group_name      = var.name
  consumer_group_user_name = var.name
  consumer_group_password  = var.password
}

data "alicloud_dts_consumer_channels" "default" {
  dts_instance_id = alicloud_dts_subscription_job.default.dts_instance_id
  %s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
