package gpdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudGpdbLogbackupDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"start_time":     `"2026-06-29T02:00Z"`,
			"end_time":       `"2026-07-09T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2022-12-13T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
	}

	pageConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"page_number":    `"1"`,
			"page_size":      `"50"`,
			"start_time":     `"2026-06-29T02:00Z"`,
			"end_time":       `"2026-07-09T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"page_number":    `"2"`,
			"page_size":      `"50"`,
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2022-12-13T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"start_time":     `"2026-06-29T02:00Z"`,
			"end_time":       `"2026-07-09T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2022-12-13T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
	}

	GpdbLogbackupCheckInfo.DataSourceTestCheck(t, rand, idsConf, pageConf, allConf)
}

var existGpdbLogbackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"logbackups.#":                acctest.CHECKSET,
		"logbackups.0.record_total":   acctest.CHECKSET,
		"logbackups.0.segment_name":   acctest.CHECKSET,
		"logbackups.0.db_instance_id": acctest.CHECKSET,
		"logbackups.0.log_time":       acctest.CHECKSET,
		"logbackups.0.log_file_size":  acctest.CHECKSET,
		"logbackups.0.log_backup_id":  acctest.CHECKSET,
		"logbackups.0.log_file_name":  acctest.CHECKSET,
	}
}

var fakeGpdbLogbackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"logbackups.#": "0",
	}
}

var GpdbLogbackupCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_gpdb_log_backups.default",
	ExistMapFunc: existGpdbLogbackupMapFunc,
	FakeMapFunc:  fakeGpdbLogbackupMapFunc,
}

func testAccCheckAlicloudGpdbLogbackupSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccGpdbLogbackup%d"
}

data "alicloud_gpdb_instances" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_gpdb_log_backups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
