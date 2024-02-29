package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOtsInstancesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%d", rand),
		dataSourceOtsInstancesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ots_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ots_instance.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ots_instance.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ots_instance.default.name}-fake",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": "${alicloud_ots_instance.default.tags}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF-fake",
				"For":     "acceptance test fake",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ots_instance.default.id}"},
			"name_regex": "${alicloud_ots_instance.default.name}",
			"tags":       "${alicloud_ots_instance.default.tags}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ots_instance.default.id}"},
			"name_regex": "${alicloud_ots_instance.default.name}",
			"tags": map[string]string{
				"Created": "TF-fake",
				"For":     "acceptance test fake",
			},
		}),
	}

	var existOtsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                       "1",
			"names.0":                       fmt.Sprintf("tf-testAcc%d", rand),
			"instances.#":                   "1",
			"instances.0.name":              fmt.Sprintf("tf-testAcc%d", rand),
			"instances.0.id":                fmt.Sprintf("tf-testAcc%d", rand),
			"instances.0.status":            string(Running),
			"instances.0.cluster_type":      CHECKSET,
			"instances.0.create_time":       CHECKSET,
			"instances.0.user_id":           CHECKSET,
			"instances.0.description":       fmt.Sprintf("tf-testAcc%d", rand),
			"instances.0.table_quota":       CHECKSET,
			"instances.0.resource_group_id": CHECKSET,
			"instances.0.tags.%":            "2",
		}
	}

	var fakeOtsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":  "0",
			"topics.#": "0",
		}
	}

	var otsInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtsInstancesMapFunc,
		fakeMapFunc:  fakeOtsInstancesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
	}
	otsInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, tagsConf, allConf)
}

func dataSourceOtsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alicloud_ots_instance" "default" {
	  name = "${var.name}"
	  description = "${var.name}"
	  instance_type = "Capacity"
	  tags = {
		Created = "TF-${var.name}"
		For = "acceptance test"
	  }
	}
	`, name)
}
