package ots_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudOtsSecondaryIndexesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_secondary_indexes.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testAcc%d", rand),
		dataSourceOtsSecondaryIndexesConfigDependence)

	instanceNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_secondary_index.default.instance_name}",
			"table_name":    "${alicloud_ots_secondary_index.default.table_name}",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_secondary_index.default.instance_name}",
			"table_name":    "${alicloud_ots_secondary_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_secondary_index.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_secondary_index.default.instance_name}",
			"table_name":    "${alicloud_ots_secondary_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_secondary_index.default.id}-fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_secondary_index.default.instance_name}",
			"table_name":    "${alicloud_ots_secondary_index.default.table_name}",
			"name_regex":    "${alicloud_ots_secondary_index.default.index_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_secondary_index.default.instance_name}",
			"table_name":    "${alicloud_ots_secondary_index.default.table_name}",
			"name_regex":    "${alicloud_ots_secondary_index.default.index_name}-fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_secondary_index.default.instance_name}",
			"table_name":    "${alicloud_ots_secondary_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_secondary_index.default.id}"},
			"name_regex":    "${alicloud_ots_secondary_index.default.index_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_secondary_index.default.instance_name}",
			"table_name":    "${alicloud_ots_secondary_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_secondary_index.default.id}"},
			"name_regex":    "${alicloud_ots_secondary_index.default.index_name}-fake",
		}),
	}

	var existOtsTablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                 "1",
			"names.0":                 acctest.CHECKSET,
			"indexes.#":               "1",
			"indexes.0.table_name":    acctest.CHECKSET,
			"indexes.0.instance_name": acctest.CHECKSET,

			"indexes.0.primary_keys.#":    "2",
			"indexes.0.defined_columns.#": "2",
		}
	}

	var fakeOtsTablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":  "0",
			"tables.#": "0",
		}
	}

	var otsTablesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existOtsTablesMapFunc,
		FakeMapFunc:  fakeOtsTablesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
	}
	otsTablesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, instanceNameConf, idsConf, nameRegexConf, allConf)
}

func dataSourceOtsSecondaryIndexesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alicloud_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${var.name}"
	  primary_key {
          name = "pk1"
	      type = "Integer"
	  }
	  primary_key {
          name = "pk2"
          type = "String"
      }
	  defined_column {
          name = "col1"
          type = "String"
      }
	  defined_column {
          name = "col2"
          type = "Integer"
      }
	  time_to_live = -1
	  max_version = 1
	}

	resource "alicloud_ots_secondary_index" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${var.name}"
	  index_name = "index1"
	  index_type = "Global"
	  include_base_data = true
	  primary_keys = ["pk2", "pk1"]
	  defined_columns = ["col2", "col1"]
	}
	`, name)
}
