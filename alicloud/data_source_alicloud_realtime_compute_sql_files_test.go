package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudRealtimeComputeSqlFileDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
	rand := acctest.RandIntRange(10000, 99999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRealtimeComputeSqlFileSourceConfig(rand),
		fakeConfig:  "",
	}

	RealtimeComputeSqlFileCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existRealtimeComputeSqlFileMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                   "1",
		"names.#":                 "1",
		"sql_files.#":             "1",
		"sql_files.0.id":          CHECKSET,
		"sql_files.0.workspace":   CHECKSET,
		"sql_files.0.namespace":   CHECKSET,
		"sql_files.0.sql_file_id": CHECKSET,
		"sql_files.0.name":        CHECKSET,
	}
}

var fakeRealtimeComputeSqlFileMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"names.#":     "0",
		"sql_files.#": "0",
	}
}

var RealtimeComputeSqlFileCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_realtime_compute_sql_files.default",
	existMapFunc: existRealtimeComputeSqlFileMapFunc,
	fakeMapFunc:  fakeRealtimeComputeSqlFileMapFunc,
}

func testAccCheckAlicloudRealtimeComputeSqlFileSourceConfig(rand int) string {
	name := fmt.Sprintf("tfaccrealtimecompute%d", rand)
	return fmt.Sprintf(`
variable "name" {
	default = "%[1]s"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "test-tf-vpc-sql-file-ds"
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "test-tf-vswitch-sql-file-ds"
}

resource "alicloud_oss_bucket" "default" {
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = ["${alicloud_vswitch.default.id}"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id     = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_sql_file" "default" {
  workspace  = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace  = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name       = var.name
  sql_script = "CREATE TABLE datagen (id VARCHAR, name VARCHAR) WITH ('connector' = 'datagen'); INSERT INTO blackhole SELECT * FROM datagen;"
  description = "test sql file for datasource"
  session_cluster_name = "test-session-cluster"
  parent_id   = "0"
}

data "alicloud_realtime_compute_sql_files" "default" {
  workspace      = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace      = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  ids            = [alicloud_realtime_compute_sql_file.default.id]
  enable_details = true
}
`, name)
}
