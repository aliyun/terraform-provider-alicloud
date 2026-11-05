package ddoscoo_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDdosCooInstanceDataSource_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_ddoscoo_instances.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceDdoscooInstanceConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddoscoo_instance.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddoscoo_instance.default.name}-fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ddoscoo_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ddoscoo_instance.default.id}-fake"},
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddoscoo_instance.default.name}",
			"ids":        []string{"${alicloud_ddoscoo_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddoscoo_instance.default.name}-fake",
			"ids":        []string{"${alicloud_ddoscoo_instance.default.id}"},
		}),
	}

	var existDdoscooInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         acctest.CHECKSET,
			"names.#":                       "1",
			"names.0":                       fmt.Sprintf("tf_testAcc%d", rand),
			"instances.#":                   "1",
			"instances.0.name":              fmt.Sprintf("tf_testAcc%d", rand),
			"instances.0.bandwidth":         "30",
			"instances.0.base_bandwidth":    "30",
			"instances.0.service_bandwidth": "100",
			"instances.0.port_count":        "50",
			"instances.0.domain_count":      "50",
			"instances.0.remark":            acctest.CHECKSET,
			"instances.0.ip_mode":           acctest.CHECKSET,
			"instances.0.debt_status":       acctest.CHECKSET,
			"instances.0.edition":           acctest.CHECKSET,
			"instances.0.status":            acctest.CHECKSET,
			"instances.0.ip_version":        acctest.CHECKSET,
			"instances.0.enabled":           acctest.CHECKSET,
			"instances.0.expire_time":       acctest.CHECKSET,
			"instances.0.create_time":       acctest.CHECKSET,
		}
	}

	var fakeDdoscooInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var ddoscooInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDdoscooInstanceMapFunc,
		FakeMapFunc:  fakeDdoscooInstanceMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckPrePaidResources(t)
		acctest.PreCheckWithTime(t, []int{15})
		acctest.PreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)

	}
	ddoscooInstanceCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
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
