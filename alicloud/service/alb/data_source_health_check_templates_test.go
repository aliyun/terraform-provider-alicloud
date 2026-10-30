package alb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudALBHealthCheckTemplatesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_health_check_template.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_health_check_template.default.id}_fake"]`,
		}),
	}

	alb_healthCheckTemplateIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"health_check_template_ids": `["${alicloud_alb_health_check_template.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"health_check_template_ids": `["${alicloud_alb_health_check_template.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_health_check_template.default.health_check_template_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_health_check_template.default.health_check_template_name}_fake"`,
		}),
	}

	alb_healthCheckTemplatenameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"health_check_template_name": `"${alicloud_alb_health_check_template.default.health_check_template_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"health_check_template_name": `"${alicloud_alb_health_check_template.default.health_check_template_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"ids":                        `["${alicloud_alb_health_check_template.default.id}"]`,
			"health_check_template_ids":  `["${alicloud_alb_health_check_template.default.id}"]`,
			"name_regex":                 `"${alicloud_alb_health_check_template.default.health_check_template_name}"`,
			"health_check_template_name": `"${alicloud_alb_health_check_template.default.health_check_template_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"ids":                        `["${alicloud_alb_health_check_template.default.id}_fake"]`,
			"health_check_template_ids":  `["${alicloud_alb_health_check_template.default.id}_fake"]`,
			"name_regex":                 `"${alicloud_alb_health_check_template.default.health_check_template_name}_fake"`,
			"health_check_template_name": `"${alicloud_alb_health_check_template.default.health_check_template_name}_fake"`,
		}),
	}

	var existDataAlicloudAlbHealthCheckTemplatesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"templates.#":                            "1",
			"templates.0.health_check_template_name": fmt.Sprintf("tf-testAccAlbHealthCheckTemplate%d", rand),
		}
	}
	var fakeDataAlicloudAlbHealthCheckTemplatesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"templates.#": "0",
		}
	}
	var alicloudAlbHealthCheckTemplateCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_alb_health_check_templates.default",
		ExistMapFunc: existDataAlicloudAlbHealthCheckTemplatesSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudAlbHealthCheckTemplatesSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbHealthCheckTemplateCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, alb_healthCheckTemplateIdsConf, idsConf, nameRegexConf, alb_healthCheckTemplatenameConf, allConf)
}
func testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAlbHealthCheckTemplate%d"
}
resource "alicloud_alb_health_check_template" "default" {
	health_check_template_name = var.name
}
data "alicloud_alb_health_check_templates" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
