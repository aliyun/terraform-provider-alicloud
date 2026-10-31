package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionBaselineStrategiesDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_baseline_strategy.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_baseline_strategy.default.id}_fake"]`,
		}),
	}

	CustomTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_baseline_strategy.default.id}"]`,
			"custom_type": `"custom"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_baseline_strategy.default.id}_fake"]`,
			"custom_type": `"custom_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_baseline_strategy.default.id}"]`,
			"custom_type": `"custom"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionBaselineStrategySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_baseline_strategy.default.id}_fake"]`,
			"custom_type": `"custom_fake"`,
		}),
	}

	ThreatDetectionBaselineStrategyCheckInfo.DataSourceTestCheck(t, rand, idsConf, CustomTypeConf, allConf)
}

var existThreatDetectionBaselineStrategyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"strategies.#":                        "1",
		"strategies.0.id":                     acctest.CHECKSET,
		"strategies.0.baseline_strategy_id":   acctest.CHECKSET,
		"strategies.0.baseline_strategy_name": acctest.CHECKSET,
		"strategies.0.custom_type":            acctest.CHECKSET,
		"strategies.0.cycle_days":             acctest.CHECKSET,
		"strategies.0.cycle_start_time":       acctest.CHECKSET,
		"strategies.0.end_time":               acctest.CHECKSET,
		"strategies.0.start_time":             acctest.CHECKSET,
		"strategies.0.target_type":            acctest.CHECKSET,
	}
}

var fakeThreatDetectionBaselineStrategyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"strategies.#": "0",
	}
}

var ThreatDetectionBaselineStrategyCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_baseline_strategies.default",
	ExistMapFunc: existThreatDetectionBaselineStrategyMapFunc,
	FakeMapFunc:  fakeThreatDetectionBaselineStrategyMapFunc,
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
