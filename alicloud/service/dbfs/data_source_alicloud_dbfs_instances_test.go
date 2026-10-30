package dbfs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDBFSInstancesDataSource(t *testing.T) {
	t.Skip("This resource has been deprecated since v1.279.0 and will be removed in the future. See: https://help.aliyun.com/en/dbfs/")
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_dbfs_instances.default"
	name := fmt.Sprintf("tf-testacc-dbfsInstance%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceDbfsInstancesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dbfs_instance.default.instance_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dbfs_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_dbfs_instance.default.id}_fake"},
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_dbfs_instance.default.id}"},
			"status": "unattached",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_dbfs_instance.default.id}_fake"},
			"status": "attaching",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_dbfs_instance.default.instance_name}",
			"ids":        []string{"${alicloud_dbfs_instance.default.id}"},
			"status":     "unattached",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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

	var DbfsInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDbfsInstancesMapFunc,
		FakeMapFunc:  fakeDbfsInstancesMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.DBFSSystemSupportRegions)
	}
	DbfsInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, allConf)
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
