package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsMetricRuleTemplatesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacccloudmonitorservicemetricruletemplate%d", rand)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"ids": `["${alicloud_cms_metric_rule_template.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"ids": `["${alicloud_cms_metric_rule_template.default.id}_fake"]`,
		}),
	}

	templateIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"template_id": `"${alicloud_cms_metric_rule_template.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"template_id": `"${alicloud_cms_metric_rule_template.default.id}001"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}_fake"`,
		}),
	}

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"keyword": `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"keyword": `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}_fake"`,
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"metric_rule_template_name": `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"metric_rule_template_name": `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"ids":                       `["${alicloud_cms_metric_rule_template.default.id}"]`,
			"template_id":               `"${alicloud_cms_metric_rule_template.default.id}"`,
			"name_regex":                `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}"`,
			"keyword":                   `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}"`,
			"metric_rule_template_name": `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name, map[string]string{
			"ids":                       `["${alicloud_cms_metric_rule_template.default.id}"]`,
			"template_id":               `"${alicloud_cms_metric_rule_template.default.id}001"`,
			"name_regex":                `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}"`,
			"keyword":                   `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}_fake"`,
			"metric_rule_template_name": `"${alicloud_cms_metric_rule_template.default.metric_rule_template_name}_fake"`,
		}),
	}

	var existAlicloudCmsMetricRuleTemplatesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"templates.#":                           "1",
			"templates.0.id":                        CHECKSET,
			"templates.0.metric_rule_template_name": fmt.Sprintf("tf-testacccloudmonitorservicemetricruletemplate%d", rand),
		}
	}
	var fakeAlicloudCmsMetricRuleTemplatesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCmsMetricRuleTemplatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_metric_rule_templates.default",
		existMapFunc: existAlicloudCmsMetricRuleTemplatesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsMetricRuleTemplatesDataSourceNameMapFunc,
	}
	alicloudCmsMetricRuleTemplatesCheckInfo.dataSourceTestCheck(t, rand, idsConf, templateIdConf, nameRegexConf, keywordConf, nameConf, allConf)
}
func testAccCheckAlicloudCmsMetricRuleTemplatesDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "%s"
}

resource "alicloud_cms_metric_rule_template" "default" {
	description = var.name
	metric_rule_template_name = var.name
	alert_templates {
		category = "ecs"
		metric_name = "cpu_total"
		namespace = "acs_ecs_dashboard"
		rule_name = var.name
		escalations {
			critical {
				comparison_operator = "GreaterThanThreshold"
				statistics = "Average"
				threshold = "90"
				times = "3"
			}
			info {
				comparison_operator = "GreaterThanThreshold"
				statistics = "Average"
				threshold = "90"
				times = "3"
			}
			warn {
				comparison_operator = "GreaterThanThreshold"
				statistics = "Average"
				threshold = "90"
				times = "3"
			}
		}
	}
}


data "alicloud_cms_metric_rule_templates" "default" {	
	%s
}
`, name, strings.Join(pairs, " \n "))
	return config
}
