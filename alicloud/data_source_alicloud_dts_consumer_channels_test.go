package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudDTSConsumerChannelsDataSource(t *testing.T) {
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

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = var.name
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.instance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  db_instance_id      = alicloud_db_instance.instance.id
  account_name        = "tftestprivilege"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_instance.instance.id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}

resource "alicloud_dts_subscription_job" "default" {
    dts_job_name                        = var.name
    payment_type                        = "PayAsYouGo"
    source_endpoint_engine_name         = "MySQL"
    source_endpoint_region              = data.alicloud_regions.default.regions.0.id
    source_endpoint_instance_type       = "RDS"
    source_endpoint_instance_id         = alicloud_db_instance.instance.id
    source_endpoint_database_name       = "tfaccountpri_0"
    source_endpoint_user_name           = "tftestprivilege"
    source_endpoint_password            = "Test12345"
    db_list                             =  <<EOF
        {"dtstestdata": {"name": "tfaccountpri_0", "all": true}}
    EOF
    subscription_instance_network_type  = "vpc"
    subscription_instance_vpc_id        = data.alicloud_vpcs.default.ids[0]
    subscription_instance_vswitch_id    = data.alicloud_vswitches.default.ids[0]
    status                              = "Normal"
}

resource "alicloud_dts_consumer_channel" "default" {
  dts_instance_id          = alicloud_dts_subscription_job.default.dts_instance_id
  consumer_group_name      = var.name
  consumer_group_user_name = var.name
  consumer_group_password  = "tftestAcc123"
}

data "alicloud_dts_consumer_channels" "default" {
  dts_instance_id = alicloud_dts_subscription_job.default.dts_instance_id
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
