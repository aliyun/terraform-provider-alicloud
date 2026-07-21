package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGpdbDatabackupDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDatabackupSourceConfig(rand, map[string]string{
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDatabackupSourceConfig(rand, map[string]string{
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDatabackupSourceConfig(rand, map[string]string{
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDatabackupSourceConfig(rand, map[string]string{
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}_fake"]`,
		}),
	}

	GpdbDatabackupCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existGpdbDatabackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"backups.#":                         "1",
		"backups.0.backup_method":           CHECKSET,
		"backups.0.status":                  CHECKSET,
		"backups.0.backup_size":             CHECKSET,
		"backups.0.backup_mode":             CHECKSET,
		"backups.0.backup_set_id":           CHECKSET,
		"backups.0.backup_end_time_local":   CHECKSET,
		"backups.0.db_instance_id":          CHECKSET,
		"backups.0.backup_end_time":         CHECKSET,
		"backups.0.consistent_time":         CHECKSET,
		"backups.0.backup_start_time":       CHECKSET,
		"backups.0.data_type":               CHECKSET,
		"backups.0.backup_start_time_local": CHECKSET,
	}
}

var fakeGpdbDatabackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"backups.#": "0",
	}
}

var GpdbDatabackupCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_gpdb_data_backups.default",
	existMapFunc: existGpdbDatabackupMapFunc,
	fakeMapFunc:  fakeGpdbDatabackupMapFunc,
}

func testAccCheckAlicloudGpdbDatabackupSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccGpdbDatabackup%d"
}
data "alicloud_gpdb_instances" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_gpdb_data_backups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
