package hbase_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

var existMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":              acctest.CHECKSET,
		"instances.0.id":           acctest.CHECKSET,
		"instances.0.name":         acctest.CHECKSET,
		"instances.0.region_id":    acctest.CHECKSET,
		"instances.0.zone_id":      acctest.CHECKSET,
		"instances.0.status":       acctest.CHECKSET,
		"instances.0.tags.%":       "2",
		"instances.0.tags.Created": "TF",
		"instances.0.tags.For":     "acceptance test",
		"ids.#":                    "1",
		"ids.0":                    acctest.CHECKSET,
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

var checkInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_hbase_instances.default",
	ExistMapFunc: existMapFunc,
	FakeMapFunc:  fakeMapFunc,
}

func TestAccAlicloudHBaseInstancesDataSourceNewInstance(t *testing.T) {
	rand := tfacctest.RandInt()
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"ids": `["${alicloud_hbase_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"ids": `["${alicloud_hbase_instance.default.id}_fake"]`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		FakeConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
			"ids":        `["${alicloud_hbase_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudHBaseDataSourceConfigNewInstance(rand, map[string]string{
			"name_regex": `"${alicloud_hbase_instance.default.name}"`,
			"ids":        `["${alicloud_hbase_instance.default.id}_fake"]`,
		}),
	}

	checkInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, tagsConf, allConf)
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

data "alicloud_hbase_zones" "default" {}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_hbase_zones.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  # Test-infra adaptation: offset 8 persistently overlaps an existing subnet
  # in the shared ^default-NODELETING$ VPC; offset 200 avoids it. The hbase
  # data-source assertions are unchanged.
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 200)
  zone_id = data.alicloud_hbase_zones.default.ids.0
  vswitch_name              = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_hbase_instance" "default" {
  name = var.name
  engine = "hbaseue"
  engine_version = "2.0"
  master_instance_type = "hbase.sn2.2xlarge"
  core_instance_type = "hbase.sn2.2xlarge"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  pay_type = "PostPaid"
  duration = 1
  auto_renew = false
  vpc_id = data.alicloud_vpcs.default.ids.0
  vswitch_id = local.vswitch_id
  cold_storage_size = 0
  deletion_protection = false
  immediate_delete_flag = true
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

data "alicloud_hbase_instances" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
