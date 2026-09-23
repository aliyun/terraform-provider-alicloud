package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGpdbLogbackupDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2024-12-12T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2022-12-13T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
	}

	pageConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"page_number":    `"1"`,
			"page_size":      `"50"`,
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2024-12-12T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"page_number":    `"2"`,
			"page_size":      `"50"`,
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2022-12-13T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2024-12-12T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbLogbackupSourceConfig(rand, map[string]string{
			"start_time":     `"2022-12-12T02:00Z"`,
			"end_time":       `"2022-12-13T02:00Z"`,
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
	}

	GpdbLogbackupCheckInfo.dataSourceTestCheck(t, rand, idsConf, pageConf, allConf)
}

var existGpdbLogbackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"logbackups.#":                CHECKSET,
		"logbackups.0.record_total":   CHECKSET,
		"logbackups.0.segment_name":   CHECKSET,
		"logbackups.0.db_instance_id": CHECKSET,
		"logbackups.0.log_time":       CHECKSET,
		"logbackups.0.log_file_size":  CHECKSET,
		"logbackups.0.log_backup_id":  CHECKSET,
		"logbackups.0.log_file_name":  CHECKSET,
	}
}

var fakeGpdbLogbackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"logbackups.#": "0",
	}
}

var GpdbLogbackupCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_gpdb_log_backups.default",
	existMapFunc: existGpdbLogbackupMapFunc,
	fakeMapFunc:  fakeGpdbLogbackupMapFunc,
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
