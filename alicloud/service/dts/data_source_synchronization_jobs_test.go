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

const dtsSynchronizationJobsDataSourceId = "data.alicloud_dts_synchronization_jobs.default"

func TestAccAlicloudDTSSynchronizationJobsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// synchronizationJobidconf
			{
				Config: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
					"ids": `["${alicloud_dts_synchronization_job.default.id}"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSynchronizationJobExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
					"ids": `["${alicloud_dts_synchronization_job.default.id}_fake"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSynchronizationJobFakeChecks()...),
			},
			// nameRegexConf
			{
				Config: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
					"name_regex": `"${alicloud_dts_synchronization_job.default.dts_job_name}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSynchronizationJobExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
					"name_regex": `"${alicloud_dts_synchronization_job.default.dts_job_name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSynchronizationJobFakeChecks()...),
			},
			// statusConf
			{
				Config: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
					"ids":    `["${alicloud_dts_synchronization_job.default.id}"]`,
					"status": `"${alicloud_dts_synchronization_job.default.status}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSynchronizationJobExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
					"ids":    `["${alicloud_dts_synchronization_job.default.id}"]`,
					"status": `"Failed"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSynchronizationJobFakeChecks()...),
			},
			// allConf
			{
				Config: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
					"ids":        `["${alicloud_dts_synchronization_job.default.id}"]`,
					"name_regex": `"${alicloud_dts_synchronization_job.default.dts_job_name}"`,
					"status":     `"${alicloud_dts_synchronization_job.default.status}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSynchronizationJobExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
					"ids":        `["${alicloud_dts_synchronization_job.default.id}"]`,
					"name_regex": `"${alicloud_dts_synchronization_job.default.dts_job_name}"`,
					"status":     `"Failed"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsSynchronizationJobFakeChecks()...),
			},
		},
	})
}

func dtsSynchronizationJobExistChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dtsSynchronizationJobsDataSourceId, "jobs.#", "1"),
		resource.TestCheckResourceAttrSet(dtsSynchronizationJobsDataSourceId, "jobs.0.dts_job_id"),
	}
}

func dtsSynchronizationJobFakeChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dtsSynchronizationJobsDataSourceId, "jobs.#", "0"),
	}
}

func testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccDtsSynchronizationJobs%d"
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

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "source" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}
resource "alicloud_db_instance" "dest" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}

resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                        = "PayAsYouGo"
  source_endpoint_engine_name         = "MySQL"
  source_endpoint_region              = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name    = "MySQL"
  destination_endpoint_region         = data.alicloud_regions.default.regions.0.id
  instance_class                      = "small"
  sync_architecture                   = "oneway"
}


resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.dest.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  db_instance_id      = alicloud_db_instance.dest.id
  account_name        = "tftestdts"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_account.account.instance_id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}

resource "alicloud_db_database" "db_r" {
  count       = 2
  instance_id = alicloud_db_instance.source.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account_r" {
  db_instance_id      = alicloud_db_instance.source.id
  account_name        = "tftestdts"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege_r" {
  instance_id  = alicloud_db_account.account_r.instance_id
  account_name = alicloud_db_account.account_r.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db_r.*.name
}

resource "alicloud_dts_synchronization_job" "default" {
  dts_instance_id                     = alicloud_dts_synchronization_instance.default.id
  dts_job_name                        = "tf-testAccCase1-%d"
  source_endpoint_instance_type       = "RDS"
  source_endpoint_instance_id         = alicloud_db_instance.source.id
  source_endpoint_engine_name         = "MySQL"
  source_endpoint_region              = data.alicloud_regions.default.regions.0.id
  source_endpoint_database_name       = "tfaccountpri_0"
  source_endpoint_user_name           = "tftestdts"
  source_endpoint_password            = "Test12345"
  destination_endpoint_instance_type  = "RDS"
  destination_endpoint_instance_id    = alicloud_db_instance.dest.id
  destination_endpoint_engine_name    = "MySQL"
  destination_endpoint_region         = data.alicloud_regions.default.regions.0.id
  destination_endpoint_database_name  = "tfaccountpri_0"
  destination_endpoint_user_name      = "tftestdts"
  destination_endpoint_password       = "Test12345"
  db_list                             = "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}"
  depends_on                          = [alicloud_db_database.db, alicloud_db_database.db_r, alicloud_db_account_privilege.privilege, alicloud_db_account_privilege.privilege_r]
  structure_initialization            = "true"
  data_initialization                 = "true"
  data_synchronization                = "true"
}

data "alicloud_dts_synchronization_jobs" "default" {
%s
}
`, rand, rand, strings.Join(pairs, "\n   "))
	return config
}
