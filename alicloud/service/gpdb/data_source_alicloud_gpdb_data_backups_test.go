package gpdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudGpdbDatabackupDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbDatabackupSourceConfig(rand, map[string]string{
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbDatabackupSourceConfig(rand, map[string]string{
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGpdbDatabackupSourceConfig(rand, map[string]string{
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}"]`,
		}),
		FakeConfig: testAccCheckAlicloudGpdbDatabackupSourceConfig(rand, map[string]string{
			"db_instance_id": `"${data.alicloud_gpdb_instances.default.ids.0}"`,
			"ids":            `["${data.alicloud_gpdb_instances.default.ids.0}_fake"]`,
		}),
	}

	GpdbDatabackupCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

var existGpdbDatabackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"backups.#":                         "1",
		"backups.0.backup_method":           acctest.CHECKSET,
		"backups.0.status":                  acctest.CHECKSET,
		"backups.0.backup_size":             acctest.CHECKSET,
		"backups.0.backup_mode":             acctest.CHECKSET,
		"backups.0.backup_set_id":           acctest.CHECKSET,
		"backups.0.backup_end_time_local":   acctest.CHECKSET,
		"backups.0.db_instance_id":          acctest.CHECKSET,
		"backups.0.backup_end_time":         acctest.CHECKSET,
		"backups.0.consistent_time":         acctest.CHECKSET,
		"backups.0.backup_start_time":       acctest.CHECKSET,
		"backups.0.data_type":               acctest.CHECKSET,
		"backups.0.backup_start_time_local": acctest.CHECKSET,
	}
}

var fakeGpdbDatabackupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"backups.#": "0",
	}
}

var GpdbDatabackupCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_gpdb_data_backups.default",
	ExistMapFunc: existGpdbDatabackupMapFunc,
	FakeMapFunc:  fakeGpdbDatabackupMapFunc,
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
