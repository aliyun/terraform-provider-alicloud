package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBHealthCheckTemplatesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_health_check_template.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_health_check_template.default.id}_fake"]`,
		}),
	}

	alb_healthCheckTemplateIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"health_check_template_ids": `["${alicloud_alb_health_check_template.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"health_check_template_ids": `["${alicloud_alb_health_check_template.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_health_check_template.default.health_check_template_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_health_check_template.default.health_check_template_name}_fake"`,
		}),
	}

	alb_healthCheckTemplatenameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"health_check_template_name": `"${alicloud_alb_health_check_template.default.health_check_template_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"health_check_template_name": `"${alicloud_alb_health_check_template.default.health_check_template_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
			"ids":                        `["${alicloud_alb_health_check_template.default.id}"]`,
			"health_check_template_ids":  `["${alicloud_alb_health_check_template.default.id}"]`,
			"name_regex":                 `"${alicloud_alb_health_check_template.default.health_check_template_name}"`,
			"health_check_template_name": `"${alicloud_alb_health_check_template.default.health_check_template_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbHealthCheckTemplateDataSourceName(rand, map[string]string{
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
	var alicloudAlbHealthCheckTemplateCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_health_check_templates.default",
		existMapFunc: existDataAlicloudAlbHealthCheckTemplatesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbHealthCheckTemplatesSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbHealthCheckTemplateCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, alb_healthCheckTemplateIdsConf, idsConf, nameRegexConf, alb_healthCheckTemplatenameConf, allConf)
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
