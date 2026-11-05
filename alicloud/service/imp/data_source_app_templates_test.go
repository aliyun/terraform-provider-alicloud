package imp_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudIMPAppTemplatesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 9999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_imp_app_template.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_imp_app_template.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_imp_app_template.default.id}"]`,
			"name_regex": `"${alicloud_imp_app_template.default.app_template_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_imp_app_template.default.id}"]`,
			"name_regex": `"${alicloud_imp_app_template.default.app_template_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_imp_app_template.default.id}"]`,
			"status": `"unattached"`,
		}),
		FakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_imp_app_template.default.id}"]`,
			"status": `"attached"`,
		}),
	}
	outputFile := filepath.Join(t.TempDir(), "imp_app_templates.txt")
	outputFileConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_imp_app_template.default.id}"]`,
			"output_file": fmt.Sprintf("%q", outputFile),
		}),
		FakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_imp_app_template.default.id}_fake"]`,
			"output_file": fmt.Sprintf("%q", outputFile),
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_imp_app_template.default.id}"]`,
			"name_regex": `"${alicloud_imp_app_template.default.app_template_name}"`,
			"status":     `"unattached"`,
		}),
		FakeConfig: testAccCheckAlicloudIMPAppTemplateDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_imp_app_template.default.id}"]`,
			"name_regex": `"${alicloud_imp_app_template.default.app_template_name}_fake"`,
			"status":     `"attached"`,
		}),
	}
	var existDataAlicloudIMPAppTemplatesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"templates.#":                      "1",
			"templates.0.app_template_creator": acctest.CHECKSET,
			"templates.0.id":                   acctest.CHECKSET,
			"templates.0.app_template_id":      acctest.CHECKSET,
			"templates.0.app_template_name":    acctest.CHECKSET,
			"templates.0.component_list.#":     "2",
			"templates.0.config_list.#":        acctest.CHECKSET,
			"templates.0.create_time":          acctest.CHECKSET,
			"templates.0.integration_mode":     "paasSDK",
			"templates.0.scene":                "business",
			"templates.0.sdk_info":             acctest.CHECKSET,
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
	var alicloudIMPAppTemplateCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_imp_app_templates.default",
		ExistMapFunc: existDataAlicloudIMPAppTemplatesSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudIMPAppTemplatesSourceNameMapFunc,
	}
	alicloudIMPAppTemplateCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, outputFileConf, allConf)
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
