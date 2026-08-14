package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsNotifyStrategiesDataSource(t *testing.T) {
	supportedRegions := []connectivity.Region{connectivity.Hangzhou}
	testAccPreCheckWithRegions(t, true, supportedRegions)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := testAccCmsUniqueResourceName()
	notifyTarget := testAccCmsCreateContact(t, rand, supportedRegions)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(name, notifyTarget, map[string]string{
			"ids": `["${alicloud_cms_notify_strategy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(name, notifyTarget, map[string]string{
			"ids": `["${alicloud_cms_notify_strategy.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(name, notifyTarget, map[string]string{
			"name_regex": `"${alicloud_cms_notify_strategy.default.notify_strategy_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(name, notifyTarget, map[string]string{
			"name_regex": `"${alicloud_cms_notify_strategy.default.notify_strategy_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(name, notifyTarget, map[string]string{
			"ids":         `["${alicloud_cms_notify_strategy.default.id}"]`,
			"name_regex":  `"${alicloud_cms_notify_strategy.default.notify_strategy_name}"`,
			"output_file": `"./tf-testacc-cms-notify-strategies.txt"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(name, notifyTarget, map[string]string{
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
		"notify_strategies.0.custom_template_entries.#": "0",
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

func testAccCheckAliCloudCmsNotifyStrategiesSourceConfig(name string, notifyTarget *testAccCmsContact, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = %q
}

resource "alicloud_cms_notify_strategy" "default" {
  notify_strategy_name = var.name
  workspace            = %q
  description          = "tf-acc cms notify strategy datasource test"
  grouping_setting {
    grouping_keys = ["severity"]
    period_min    = 5
    silence_sec   = 300
    times         = 1
  }
  routes {
    channels {
      channel_type         = "CONTACT"
      receivers            = [%q]
      enabled_sub_channels = ["EMAIL"]
    }
    effect_time_range {
      day_in_week          = [1, 2, 3, 4, 5]
      start_time_in_minute = 0
      end_time_in_minute   = 1439
      time_zone            = "Asia/Shanghai"
    }
    severities = ["CRITICAL"]
  }
}

data "alicloud_cms_notify_strategies" "default" {
%s
}
`, name, notifyTarget.Workspace, notifyTarget.Identifier, strings.Join(pairs, "\n   "))
	return config
}
