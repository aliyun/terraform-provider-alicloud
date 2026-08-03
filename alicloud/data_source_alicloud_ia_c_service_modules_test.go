package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudIaCServiceModuleDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc%siacmodule%d", defaultRegionToTest, rand)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIaCServiceModuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ia_c_service_module.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudIaCServiceModuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ia_c_service_module.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIaCServiceModuleSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf("\"^%s$\"", name),
		}),
		fakeConfig: testAccCheckAliCloudIaCServiceModuleSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf("\"^%s_fake$\"", name),
		}),
	}

	moduleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIaCServiceModuleSourceConfig(rand, map[string]string{
			"module_name": `"${alicloud_ia_c_service_module.default.module_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudIaCServiceModuleSourceConfig(rand, map[string]string{
			"module_name": `"${alicloud_ia_c_service_module.default.module_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudIaCServiceModuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ia_c_service_module.default.id}"]`,
			"name_regex":  fmt.Sprintf("\"^%s$\"", name),
			"module_name": `"${alicloud_ia_c_service_module.default.module_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudIaCServiceModuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ia_c_service_module.default.id}_fake"]`,
			"name_regex":  fmt.Sprintf("\"^%s_fake$\"", name),
			"module_name": `"${alicloud_ia_c_service_module.default.module_name}_fake"`,
		}),
	}

	IaCServiceModuleCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, moduleNameConf, allConf)
}

var existIaCServiceModuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"modules.#":                  "1",
		"modules.0.module_id":        CHECKSET,
		"modules.0.module_name":      CHECKSET,
		"modules.0.description":      CHECKSET,
		"modules.0.source":           CHECKSET,
		"modules.0.status":           CHECKSET,
		"modules.0.create_time":      CHECKSET,
		"modules.0.tags.#":           "1",
		"modules.0.tags.0.tag_key":   CHECKSET,
		"modules.0.tags.0.tag_value": CHECKSET,
	}
}

var fakeIaCServiceModuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"modules.#": "0",
	}
}

var IaCServiceModuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ia_c_service_modules.default",
	existMapFunc: existIaCServiceModuleMapFunc,
	fakeMapFunc:  fakeIaCServiceModuleMapFunc,
}

func testAccCheckAliCloudIaCServiceModuleSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	name := fmt.Sprintf("tf-testacc%siacmodule%d", defaultRegionToTest, rand)
	config := fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_ia_c_service_module" "default" {
  module_name      = var.name
  source           = "Registry"
  source_path      = "alibaba/security-group:2.4.1"
  version_strategy = "Manual"
  description      = var.name
  tags {
    tag_key   = "Created"
    tag_value = "TF"
  }
}

data "alicloud_ia_c_service_modules" "default" {
%s
}
`, name, strings.Join(pairs, "\n   "))
	return config
}
