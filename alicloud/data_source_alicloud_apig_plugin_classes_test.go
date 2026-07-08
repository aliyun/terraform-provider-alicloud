package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudApigPluginClassesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_plugin_class.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_plugin_class.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_plugin_class.default.id}"]`,
			"name_regex": `"${alicloud_apig_plugin_class.default.plugin_class_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_plugin_class.default.id}"]`,
			"name_regex": `"${alicloud_apig_plugin_class.default.plugin_class_name}_fake"`,
		}),
	}

	detailsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_class.default.id}"]`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_class.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_class.default.id}"]`,
			"name_regex":     `"${alicloud_apig_plugin_class.default.plugin_class_name}"`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginClassesSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_plugin_class.default.id}_fake"]`,
			"name_regex":     `"${alicloud_apig_plugin_class.default.plugin_class_name}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	ApigPluginClassesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, detailsConf, allConf)
}

var existApigPluginClassesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"classes.#":                   "1",
		"classes.0.id":                CHECKSET,
		"classes.0.plugin_class_id":   CHECKSET,
		"classes.0.plugin_class_name": CHECKSET,
		"classes.0.description":       CHECKSET,
		"classes.0.version":           CHECKSET,
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

var ApigPluginClassesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_plugin_classes.default",
	existMapFunc: existApigPluginClassesMapFunc,
	fakeMapFunc:  fakeApigPluginClassesMapFunc,
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
