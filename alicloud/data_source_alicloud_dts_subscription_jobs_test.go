package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDTSSubscriptionJobsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	subscriptionJobidconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dts_subscription_job.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dts_subscription_job.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dts_subscription_job.default.dts_job_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dts_subscription_job.default.dts_job_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dts_subscription_job.default.id}"]`,
			"name_regex": `"${alicloud_dts_subscription_job.default.dts_job_name}"`,
			"status":     `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dts_subscription_job.default.id}"]`,
			"name_regex": `"${alicloud_dts_subscription_job.default.dts_job_name}_fake"`,
			"status":     `"Abnormal"`,
		}),
	}

	DtsSubscriptionJobCheckInfo.dataSourceTestCheck(t, rand, subscriptionJobidconf, nameRegexConf, allConf)
}

func testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccDtsSubscriptionJobs%d"
}
variable "region_id" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "5.6"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "5.6"
	instance_charge_type = "PostPaid"
}

resource "alicloud_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
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
    source_endpoint_region              = var.region_id
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

data "alicloud_dts_subscription_jobs" "default" {
%s
}
`, rand, os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, "\n   "))
	return config
}

var existDtsSubscriptionJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#":            "1",
		"jobs.0.dts_job_id": CHECKSET,
	}
}

var fakeDtsSubscriptionJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#": "0",
	}
}

var DtsSubscriptionJobCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dts_subscription_jobs.default",
	existMapFunc: existDtsSubscriptionJobMapFunc,
	fakeMapFunc:  fakeDtsSubscriptionJobMapFunc,
}
