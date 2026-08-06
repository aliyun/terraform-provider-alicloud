package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsNotifyStrategiesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_notify_strategy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_notify_strategy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_notify_strategy.default.notify_strategy_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_notify_strategy.default.notify_strategy_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cms_notify_strategy.default.id}"]`,
			"name_regex":  `"${alicloud_cms_notify_strategy.default.notify_strategy_name}"`,
			"output_file": `"./tf-testacc-cms-notify-strategies.txt"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_notify_strategy.default.id}"]`,
			"workspace": `"tf-testacc-nonexistent-workspace"`,
		}),
	}
	CmsNotifyStrategiesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existCmsNotifyStrategiesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                  "1",
		"names.#":                                "1",
		"notify_strategies.#":                    "1",
		"notify_strategies.0.notify_strategy_id": CHECKSET,
		"notify_strategies.0.notify_strategy_name":      CHECKSET,
		"notify_strategies.0.enable":                    CHECKSET,
		"notify_strategies.0.create_time":               CHECKSET,
		"notify_strategies.0.update_time":               CHECKSET,
		"notify_strategies.0.grouping_setting.#":        "1",
		"notify_strategies.0.routes.#":                  "1",
		"notify_strategies.0.custom_template_entries.#": CHECKSET,
	}
}

var fakeCmsNotifyStrategiesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":               "0",
		"names.#":             "0",
		"notify_strategies.#": "0",
	}
}

var CmsNotifyStrategiesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cms_notify_strategies.default",
	existMapFunc: existCmsNotifyStrategiesMapFunc,
	fakeMapFunc:  fakeCmsNotifyStrategiesMapFunc,
}

func testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCmsNotifyStrategy%d"
}

resource "alicloud_cms_notify_strategy" "default" {
  notify_strategy_name = var.name
  description          = "tf-acc cms notify strategy datasource test"
  grouping_setting {
    grouping_keys = ["severity"]
    period_min    = 5
    silence_sec   = 300
    times         = 1
  }
  routes {
    channels {
      channel_type = "DING"
      receivers    = ["CONTACT"]
    }
    effect_time_range {
      day_in_week          = [1, 2, 3, 4, 5]
      start_time_in_minute = 0
      end_time_in_minute   = 1439
      time_zone            = "Asia/Shanghai"
    }
    severities = ["CRITICAL"]
  }
  custom_template_entries {
    target_type   = "MAIL"
    template_uuid = "${var.name}-tpl-mail"
  }
}

data "alicloud_cms_notify_strategies" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
