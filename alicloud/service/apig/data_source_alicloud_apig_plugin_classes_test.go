package apig_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudApigPluginClassesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_plugin_class.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_plugin_class.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_plugin_class.default.id}"]`,
			"name_regex": `"${alicloud_apig_plugin_class.default.plugin_class_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_plugin_class.default.id}"]`,
			"name_regex": `"${alicloud_apig_plugin_class.default.plugin_class_name}_fake"`,
		}),
	}

	detailsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_class.default.id}"]`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_class.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_class.default.id}"]`,
			"name_regex":     `"${alicloud_apig_plugin_class.default.plugin_class_name}"`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_class.default.id}_fake"]`,
			"name_regex":     `"${alicloud_apig_plugin_class.default.plugin_class_name}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	ApigPluginClassesCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, detailsConf, allConf)
}

var existApigPluginClassesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"classes.#":                   "1",
		"classes.0.id":                acctest.CHECKSET,
		"classes.0.plugin_class_id":   acctest.CHECKSET,
		"classes.0.plugin_class_name": acctest.CHECKSET,
		"classes.0.description":       acctest.CHECKSET,
		"classes.0.version":           acctest.CHECKSET,
		"ids.#":                       "1",
		"names.#":                     "1",
	}
}

var fakeApigPluginClassesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"classes.#": "0",
		"ids.#":     "0",
		"names.#":   "0",
	}
}

var ApigPluginClassesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_apig_plugin_classes.default",
	ExistMapFunc: existApigPluginClassesMapFunc,
	FakeMapFunc:  fakeApigPluginClassesMapFunc,
}

func testAccCheckAlicloudApigPluginClassesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
    default = "tfaccapig%d"
}

resource "alicloud_apig_plugin_class" "default" {
  wasm_url            = "https://example.com/plugin.wasm"
  description         = "A test plugin class for CloudSpec coverage"
  version_description = "Initial version for testing"
  plugin_class_name   = var.name
  version             = "1.0.2"
  alias               = "test-plugin-alias"
  execute_priority    = "1"
  wasm_language       = "TinyGo"
  execute_stage       = "UNSPECIFIED_PHASE"
}

data "alicloud_apig_plugin_classes" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
