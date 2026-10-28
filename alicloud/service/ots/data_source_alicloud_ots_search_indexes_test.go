package ots_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudOtsSearchIndexesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_search_indexes.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testAcc%d", rand),
		dataSourceOtsSearchIndexesConfigDependence)

	instanceNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_search_index.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_search_index.default.id}-fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"name_regex":    "${alicloud_ots_search_index.default.index_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"name_regex":    "${alicloud_ots_search_index.default.index_name}-fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_search_index.default.id}"},
			"name_regex":    "${alicloud_ots_search_index.default.index_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_search_index.default.id}"},
			"name_regex":    "${alicloud_ots_search_index.default.index_name}-fake",
		}),
	}

	var existOtsTablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                             "1",
			"names.0":                             acctest.CHECKSET,
			"indexes.#":                           "1",
			"indexes.0.id":                        acctest.CHECKSET,
			"indexes.0.instance_name":             acctest.CHECKSET,
			"indexes.0.table_name":                acctest.CHECKSET,
			"indexes.0.index_name":                acctest.CHECKSET,
			"indexes.0.create_time":               acctest.CHECKSET,
			"indexes.0.time_to_live":              acctest.CHECKSET,
			"indexes.0.sync_phase":                acctest.CHECKSET,
			"indexes.0.current_sync_timestamp":    acctest.CHECKSET,
			"indexes.0.schema":                    acctest.CHECKSET,
			"indexes.0.storage_size":              acctest.CHECKSET,
			"indexes.0.row_count":                 acctest.CHECKSET,
			"indexes.0.reserved_read_cu":          acctest.CHECKSET,
			"indexes.0.metering_last_update_time": acctest.CHECKSET,
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

func dataSourceOtsSearchIndexesConfigDependence(name string) string {
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

	resource "alicloud_ots_search_index" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${alicloud_ots_table.default.table_name}"
	  index_name = "index1"
	  time_to_live = -1

      schema {
          field_schema {
            field_name = "col1"
            field_type = "Text"
            is_array = false
            index = true
            analyzer = "Split"
            store = true
          }
          field_schema {
            field_name =  "col2"
             field_type = "Long"
             enable_sort_and_agg = true
          }
       

          field_schema {
            field_name =  "pk1"
            field_type = "Long"
            
          }
          field_schema {
            field_name =  "pk2"
            field_type = "Text"
            
          }


          index_setting {
            routing_fields = [ "pk1", "pk2"]
          }

          index_sort {
            sorter {
              sorter_type = "PrimaryKeySort"
              order = "Asc"
            }
            sorter {
              sorter_type = "FieldSort"
              order = "Desc"
              field_name =  "col2"
              mode = "Max"
            }
          }
    }


	}
	`, name)
}
