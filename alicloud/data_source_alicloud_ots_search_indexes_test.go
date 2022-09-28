package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOtsSearchIndexesDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_search_indexes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testAcc%d", rand),
		dataSourceOtsSearchIndexesConfigDependence)

	instanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_search_index.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_search_index.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"name_regex":    "${alicloud_ots_search_index.default.index_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"name_regex":    "${alicloud_ots_search_index.default.index_name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_search_index.default.id}"},
			"name_regex":    "${alicloud_ots_search_index.default.index_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_search_index.default.instance_name}",
			"table_name":    "${alicloud_ots_search_index.default.table_name}",
			"ids":           []string{"${alicloud_ots_search_index.default.id}"},
			"name_regex":    "${alicloud_ots_search_index.default.index_name}-fake",
		}),
	}

	var existOtsTablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                             "1",
			"names.0":                             CHECKSET,
			"indexes.#":                           "1",
			"indexes.0.id":                        CHECKSET,
			"indexes.0.instance_name":             CHECKSET,
			"indexes.0.table_name":                CHECKSET,
			"indexes.0.index_name":                CHECKSET,
			"indexes.0.create_time":               CHECKSET,
			"indexes.0.time_to_live":              CHECKSET,
			"indexes.0.sync_phase":                CHECKSET,
			"indexes.0.current_sync_timestamp":    CHECKSET,
			"indexes.0.schema":                    CHECKSET,
			"indexes.0.storage_size":              CHECKSET,
			"indexes.0.row_count":                 CHECKSET,
			"indexes.0.reserved_read_cu":          CHECKSET,
			"indexes.0.metering_last_update_time": CHECKSET,
		}
	}

	var fakeOtsTablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":  "0",
			"tables.#": "0",
		}
	}

	var otsTablesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtsTablesMapFunc,
		fakeMapFunc:  fakeOtsTablesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
	}
	otsTablesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, instanceNameConf, idsConf, nameRegexConf, allConf)
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
	  table_name = "${var.name}"
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
