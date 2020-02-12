package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

var existMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":              CHECKSET,
		"instances.0.id":           CHECKSET,
		"instances.0.name":         CHECKSET,
		"instances.0.region_id":    CHECKSET,
		"instances.0.zone_id":      CHECKSET,
		"instances.0.status":       CHECKSET,
		"instances.0.tags.%":       "2",
		"instances.0.tags.Created": "TF",
		"instances.0.tags.For":     "acceptance test",
		"ids.#":                    "1",
		"ids.0":                    CHECKSET,
		"names.#":                  "1",
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

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
			"tags":       `{Created = "TF1"}`,
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

	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
	}
	checkInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, tagsConf, allConf)
}

// new a instance config
func testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccHBaseInstance_datasource_%d"
}

data "alicloud_vpcs" "default" {
  is_default = "true"
}
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "tf_testAccHBase_vpc"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, %d)}"
}

resource "alicloud_hbase_instance" "default" {
  name = var.name
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  pay_type = "PostPaid"
  duration = 1
  auto_renew = false
  vswitch_id = "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]}"
  cold_storage_size = 0
  deletion_protection = false
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

data "alicloud_hbase_instances" "default" {
  %s
}
`, rand, acctest.RandIntRange(10, 100), strings.Join(pairs, "\n  "))
	return config
}
