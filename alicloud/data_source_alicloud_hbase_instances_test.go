package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

var existMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":           CHECKSET,
		"instances.0.id":        CHECKSET,
		"instances.0.name":      CHECKSET,
		"instances.0.region_id": CHECKSET,
		"instances.0.zone_id":   CHECKSET,
		"instances.0.engine":    "hbase",
		"instances.0.status":    CHECKSET,
		"ids.#":                 "1",
		"ids.0":                 CHECKSET,
		"names.#":               "1",
	}
}
var fakeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
		"ids.#":       "0",
		"names.#":     "0",
	}
}

var checkInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbase_instances.default",
	existMapFunc: existMapFunc,
	fakeMapFunc:  fakeMapFunc,
}

func TestAccAlicloudHBaseInstancesDataSourceNewInstance(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"ids": `["${alicloud_hbase_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"ids": `["${alicloud_hbase_instance.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
			"ids":        `["${alicloud_hbase_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
			"ids":        `["${alicloud_hbase_instance.default.id}_fake"]`,
		}),
	}

	checkInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

// new a instance
func testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccHBaseInstance_datasource_%d"
}

data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name = var.name
  engine_version = "2.0"
  master_instance_type = "hbase.n1.medium"
  core_instance_type = "hbase.n1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  pay_type = "Postpaid"
  duration = 0
  auto_renew = "false"
  vswitch_id = "vsw-wz9iqvmkdua0svi31ox61"
  is_cold_storage = "false"
}

data "alicloud_hbase_instances" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
