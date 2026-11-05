package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionAntiBruteForceRulesDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_anti_brute_force_rule.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_anti_brute_force_rule.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_anti_brute_force_rule.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_anti_brute_force_rule.default.id}_fake"]`,
		}),
	}

	ThreatDetectionAntiBruteForceRuleCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

var existThreatDetectionAntiBruteForceRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                            "1",
		"rules.0.id":                         acctest.CHECKSET,
		"rules.0.anti_brute_force_rule_id":   acctest.CHECKSET,
		"rules.0.anti_brute_force_rule_name": acctest.CHECKSET,
		"rules.0.default_rule":               acctest.CHECKSET,
		"rules.0.fail_count":                 acctest.CHECKSET,
		"rules.0.forbidden_time":             acctest.CHECKSET,
		"rules.0.span":                       acctest.CHECKSET,
		"rules.0.uuid_list.#":                "1",
	}
}

var fakeThreatDetectionAntiBruteForceRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
	}
}

var ThreatDetectionAntiBruteForceRuleCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_anti_brute_force_rules.default",
	ExistMapFunc: existThreatDetectionAntiBruteForceRuleMapFunc,
	FakeMapFunc:  fakeThreatDetectionAntiBruteForceRuleMapFunc,
}

func testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionAntiBruteForceRule%d"
}

data "alicloud_threat_detection_assets" "default" {
    machine_types = "ecs"
}

resource "alicloud_threat_detection_anti_brute_force_rule" "default" {
  anti_brute_force_rule_name = var.name
  forbidden_time             = 360
  uuid_list = [
  "${data.alicloud_threat_detection_assets.default.assets.0.uuid}"]
  fail_count = 80
  span       = 10
}

data "alicloud_threat_detection_anti_brute_force_rules" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
