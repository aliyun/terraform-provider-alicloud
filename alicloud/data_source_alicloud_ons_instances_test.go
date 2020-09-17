package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOnsInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ons_instances.default"
	name := fmt.Sprintf("tf-testacc%sonsinstance%v", defaultRegionToTest, rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOnsInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ons_instance.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ons_instance.default.name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ons_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ons_instance.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_ons_instance.default.id}"},
			"status": "5",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_ons_instance.default.id}"},
			"status": "0",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ons_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ons_instance.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ons_instance.default.id}"},
			"name_regex": "${alicloud_ons_instance.default.name}",
			"status":     "5",
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ons_instance.default.id}_fake"},
			"name_regex": "${alicloud_ons_instance.default.name}",
			"status":     "0",
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	var existOnsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"instances.#":                 "1",
			"names.#":                     "1",
			"instances.0.instance_status": "5",
			"instances.0.status":          "5",
			"instances.0.release_time":    CHECKSET,
			"instances.0.instance_type":   "1",
			"instances.0.instance_name":   fmt.Sprintf("tf-testacc%sonsinstance%v", defaultRegionToTest, rand),
		}
	}

	var fakeOnsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
			"names.#":     "0",
		}
	}

	var onsRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOnsInstancesMapFunc,
		fakeMapFunc:  fakeOnsInstancesMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	onsRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, tagsConf, allConf)
}

func dataSourceOnsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}

resource "alicloud_ons_instance" "default" {
  name   = "${var.name}"
  remark = "default-remark"
  tags 	 = {
	Created = "TF"
	For 	= "acceptance test"
  }
}

`, name)
}
