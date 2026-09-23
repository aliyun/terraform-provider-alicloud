package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionBaselineStrategyDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_baseline_strategy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_baseline_strategy.default.id}_fake"]`,
		}),
	}

	CustomTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_baseline_strategy.default.id}"]`,
			"custom_type": `"custom"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_baseline_strategy.default.id}_fake"]`,
			"custom_type": `"custom_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_baseline_strategy.default.id}"]`,
			"custom_type": `"custom"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_baseline_strategy.default.id}_fake"]`,
			"custom_type": `"custom_fake"`,
		}),
	}

	ThreatDetectionBaselineStrategyCheckInfo.dataSourceTestCheck(t, rand, idsConf, CustomTypeConf, allConf)
}

var existThreatDetectionBaselineStrategyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"strategies.#":                        "1",
		"strategies.0.id":                     CHECKSET,
		"strategies.0.baseline_strategy_id":   CHECKSET,
		"strategies.0.baseline_strategy_name": CHECKSET,
		"strategies.0.custom_type":            CHECKSET,
		"strategies.0.cycle_days":             CHECKSET,
		"strategies.0.cycle_start_time":       CHECKSET,
		"strategies.0.end_time":               CHECKSET,
		"strategies.0.start_time":             CHECKSET,
		"strategies.0.target_type":            CHECKSET,
	}
}

var fakeThreatDetectionBaselineStrategyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"strategies.#": "0",
	}
}

var ThreatDetectionBaselineStrategyCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_baseline_strategies.default",
	existMapFunc: existThreatDetectionBaselineStrategyMapFunc,
	fakeMapFunc:  fakeThreatDetectionBaselineStrategyMapFunc,
}

func testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionBaselineStrategy%d"
}

resource "alicloud_threat_detection_baseline_strategy" "default" {
  custom_type            = "custom"
  end_time               = "08:00:00"
  baseline_strategy_name = var.name
  cycle_days             = 3
  target_type            = "groupId"
  start_time             = "05:00:00"
  risk_sub_type_name     = "hc_exploit_redis"
}

data "alicloud_threat_detection_baseline_strategies" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
