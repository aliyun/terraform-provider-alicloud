package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEsaSiteEnvironmentsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_site_environments.default"
	name := fmt.Sprintf("Environment%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaSiteEnvironmentsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_site_environment.default.site_id}",
			"ids":     []string{"${alicloud_esa_site_environment.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_site_environment.default.site_id}",
			"ids":     []string{"${alicloud_esa_site_environment.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_site_environment.default.site_id}",
			"name_regex": "${alicloud_esa_site_environment.default.environment_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_site_environment.default.site_id}",
			"name_regex": "${alicloud_esa_site_environment.default.environment_name}_fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_site_environment.default.site_id}",
			"ids":        []string{"${alicloud_esa_site_environment.default.id}"},
			"name_regex": "${alicloud_esa_site_environment.default.environment_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_site_environment.default.site_id}",
			"ids":        []string{"${alicloud_esa_site_environment.default.id}_fake"},
			"name_regex": "${alicloud_esa_site_environment.default.environment_name}_fake",
		}),
	}

	var existAliCloudEsaSiteEnvironmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"environments.#":                  "1",
			"environments.0.id":               acctest.CHECKSET,
			"environments.0.environment_name": acctest.CHECKSET,
			"environments.0.rule":             acctest.CHECKSET,
			"environments.0.read_only":        acctest.CHECKSET,
			"environments.0.site_version":     acctest.CHECKSET,
			"environments.0.pre_site_version": acctest.CHECKSET,
			"environments.0.priority":         acctest.CHECKSET,
			"environments.0.is_default":       acctest.CHECKSET,
			"environments.0.create_time":      acctest.CHECKSET,
			"environments.0.update_time":      acctest.CHECKSET,
		}
	}

	var fakeAliCloudEsaSiteEnvironmentsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"environments.#": "0",
		}
	}

	var aliCloudEsaSiteEnvironmentsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_esa_site_environments.default",
		ExistMapFunc: existAliCloudEsaSiteEnvironmentsMapFunc,
		FakeMapFunc:  fakeAliCloudEsaSiteEnvironmentsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaSiteEnvironmentsInfo.DataSourceTestCheckWithProtoV5PreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceEsaSiteEnvironmentsConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}

resource "alicloud_esa_site_environment" "default" {
  site_id               = data.alicloud_esa_sites.default.sites.0.id
  environment_name      = var.name
  rule                  = "(ip.src in {1.1.1.1})"
  next_environment_name = "Production"
  read_only             = true
}
`, name)
}
