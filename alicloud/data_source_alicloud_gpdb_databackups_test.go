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
		"databackups.#":                         "1",
		"databackups.0.backup_method":           CHECKSET,
		"databackups.0.status":                  CHECKSET,
		"databackups.0.backup_size":             CHECKSET,
		"databackups.0.backup_mode":             CHECKSET,
		"databackups.0.backup_set_id":           CHECKSET,
		"databackups.0.backup_end_time_local":   CHECKSET,
		"databackups.0.db_instance_id":          CHECKSET,
		"databackups.0.backup_end_time":         CHECKSET,
		"databackups.0.consistent_time":         CHECKSET,
		"databackups.0.backup_start_time":       CHECKSET,
		"databackups.0.data_type":               CHECKSET,
		"databackups.0.backup_start_time_local": CHECKSET,
	}
}

var fakeGpdbDatabackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databackups.#": "0",
	}
}

var GpdbDatabackupCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_gpdb_databackups.default",
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

data "alicloud_gpdb_databackups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
