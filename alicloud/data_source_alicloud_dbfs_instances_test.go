package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDBFSInstancesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_dbfs_instances.default"
	name := fmt.Sprintf("tf-testacc-dbfsInstance%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDbfsInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dbfs_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dbfs_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dbfs_instance.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_dbfs_instance.default.id}"},
			"status": "unattached",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_dbfs_instance.default.id}_fake"},
			"status": "attaching",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dbfs_instance.default.instance_name}",
			"ids":        []string{"${alicloud_dbfs_instance.default.id}"},
			"status":     "unattached",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dbfs_instance.default.instance_name}",
			"ids":        []string{"${alicloud_dbfs_instance.default.id}"},
			"status":     "attaching",
		}),
	}
	var existDbfsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"instances.#":               "1",
			"instances.0.instance_name": name,
		}
	}

	var fakeDbfsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"names.#":     "0",
			"ids.#":       "0",
		}
	}

	var DbfsInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDbfsInstancesMapFunc,
		fakeMapFunc:  fakeDbfsInstancesMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.DBFSSystemSupportRegions)
	}
	DbfsInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)
}

func dataSourceDbfsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		resource "alicloud_dbfs_instance" "default" {
		  category          = "standard"
		  zone_id           = "cn-hangzhou-i"
		  performance_level = "PL1"
		  instance_name     = var.name
		  size              = 100
		}
		`, name)
}
