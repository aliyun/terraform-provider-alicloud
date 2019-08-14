package alicloud

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudOtsTablesDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_tables.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testAcc%d", rand),
		dataSourceOtsTablesConfigDependence)

	instanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_table.default.instance_name}",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_table.default.instance_name}",
			"ids":           []string{"${alicloud_ots_table.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_table.default.instance_name}",
			"ids":           []string{"${alicloud_ots_table.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_table.default.instance_name}",
			"name_regex":    "${alicloud_ots_table.default.table_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_table.default.instance_name}",
			"name_regex":    "${alicloud_ots_table.default.table_name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_table.default.instance_name}",
			"ids":           []string{"${alicloud_ots_table.default.id}"},
			"name_regex":    "${alicloud_ots_table.default.table_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_table.default.instance_name}",
			"ids":           []string{"${alicloud_ots_table.default.id}"},
			"name_regex":    "${alicloud_ots_table.default.table_name}-fake",
		}),
	}

	var existOtsTablesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                "1",
			"names.0":                CHECKSET,
			"tables.#":               "1",
			"tables.0.table_name":    CHECKSET,
			"tables.0.instance_name": CHECKSET,
			"tables.0.primary_key.#": "2",
			"tables.0.time_to_live":  "-1",
			"tables.0.max_version":   "1",
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

func dataSourceOtsTablesConfigDependence(name string) string {
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
	  time_to_live = -1
	  max_version = 1
	}
	`, name)
}
