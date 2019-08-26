package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDdoscooInstanceDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_ddoscoo_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceDdoscooInstanceConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddoscoo_instance.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddoscoo_instance.default.name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ddoscoo_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ddoscoo_instance.default.id}-fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddoscoo_instance.default.name}",
			"ids":        []string{"${alicloud_ddoscoo_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddoscoo_instance.default.name}-fake",
			"ids":        []string{"${alicloud_ddoscoo_instance.default.id}"},
		}),
	}

	var existDdoscooInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         CHECKSET,
			"names.#":                       "1",
			"names.0":                       fmt.Sprintf("tf_testAcc%d", rand),
			"instances.#":                   "1",
			"instances.0.name":              fmt.Sprintf("tf_testAcc%d", rand),
			"instances.0.bandwidth":         "30",
			"instances.0.base_bandwidth":    "30",
			"instances.0.service_bandwidth": "100",
			"instances.0.port_count":        "50",
			"instances.0.domain_count":      "50",
		}
	}

	var fakeDdoscooInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var ddoscooInstanceCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDdoscooInstanceMapFunc,
		fakeMapFunc:  fakeDdoscooInstanceMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)

	}
	ddoscooInstanceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceDdoscooInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddoscoo_instance" "default" {
      name                    = "%s"
      bandwidth               = "30"
      base_bandwidth          = "30"
      service_bandwidth       = "100"
      port_count              = "50"
      domain_count            = "50"
	}`, name)
}
