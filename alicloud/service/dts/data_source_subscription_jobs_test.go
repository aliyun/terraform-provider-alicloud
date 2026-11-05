package dts_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const dtsSubscriptionJobsDataSourceId = "data.alicloud_dts_subscription_jobs.default"

func TestAccAlicloudDTSSubscriptionJobsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// subscriptionJobidconf
			{
				Config: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
					"ids": `["${alicloud_dts_subscription_job.default.id}"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSubscriptionJobExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
					"ids": `["${alicloud_dts_subscription_job.default.id}_fake"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSubscriptionJobFakeChecks()...),
			},
			// nameRegexConf
			{
				Config: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
					"name_regex": `"${alicloud_dts_subscription_job.default.dts_job_name}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSubscriptionJobExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
					"name_regex": `"${alicloud_dts_subscription_job.default.dts_job_name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSubscriptionJobFakeChecks()...),
			},
			// allConf
			{
				Config: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
					"ids":        `["${alicloud_dts_subscription_job.default.id}"]`,
					"name_regex": `"${alicloud_dts_subscription_job.default.dts_job_name}"`,
					"status":     `"${alicloud_dts_subscription_job.default.status}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSubscriptionJobExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsSubscriptionJobSourceConfig(rand, map[string]string{
					"ids":        `["${alicloud_dts_subscription_job.default.id}"]`,
					"name_regex": `"${alicloud_dts_subscription_job.default.dts_job_name}_fake"`,
					"status":     `"Abnormal"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSubscriptionJobFakeChecks()...),
			},
		},
	})
}

func dtsSubscriptionJobExistChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dtsSubscriptionJobsDataSourceId, "jobs.#", "1"),
		resource.TestCheckResourceAttrSet(dtsSubscriptionJobsDataSourceId, "jobs.0.dts_job_id"),
	}
}

func dtsSubscriptionJobFakeChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dtsSubscriptionJobsDataSourceId, "jobs.#", "0"),
	}
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
}

data "alicloud_dts_subscription_jobs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
