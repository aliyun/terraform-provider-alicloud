package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudIMPAppTemplatesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 9999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_imp_app_template.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_imp_app_template.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_imp_app_template.default.id}"]`,
			"name_regex": `"${alicloud_imp_app_template.default.app_template_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_imp_app_template.default.id}"]`,
			"name_regex": `"${alicloud_imp_app_template.default.app_template_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_imp_app_template.default.id}"]`,
			"status": `"unattached"`,
		}),
		fakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_imp_app_template.default.id}"]`,
			"status": `"attached"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_imp_app_template.default.id}"]`,
			"name_regex": `"${alicloud_imp_app_template.default.app_template_name}"`,
			"status":     `"unattached"`,
		}),
		fakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_imp_app_template.default.id}"]`,
			"name_regex": `"${alicloud_imp_app_template.default.app_template_name}_fake"`,
			"status":     `"attached"`,
		}),
	}
	var existDataAlicloudIMPAppTemplatesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"templates.#":                      "1",
			"templates.0.app_template_creator": CHECKSET,
			"templates.0.id":                   CHECKSET,
			"templates.0.app_template_id":      CHECKSET,
			"templates.0.app_template_name":    CHECKSET,
			"templates.0.component_list.#":     "2",
			"templates.0.config_list.#":        CHECKSET,
			"templates.0.create_time":          CHECKSET,
			"templates.0.integration_mode":     "paasSDK",
			"templates.0.scene":                "business",
			"templates.0.sdk_info":             CHECKSET,
			"templates.0.standard_room_info":   "",
			"templates.0.status":               "unattached",
		}
	}
	var fakeDataAlicloudIMPAppTemplatesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"templates.#": "0",
		}
	}
	var alicloudIMPAppTemplateCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_imp_app_templates.default",
		existMapFunc: existDataAlicloudIMPAppTemplatesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudIMPAppTemplatesSourceNameMapFunc,
	}
	alicloudIMPAppTemplateCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudIMPAppTemplateDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf_testAccIMPAppTemplate%d"
}


resource "alicloud_imp_app_template" "default" {
  app_template_name = var.name
  component_list    = ["component.live", "component.liveRecord"]
  integration_mode  = "paasSDK"
  scene             = "business"
}

data "alicloud_imp_app_templates" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
