package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionAntiBruteForceRuleDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_anti_brute_force_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_anti_brute_force_rule.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_anti_brute_force_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionAntiBruteForceRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_anti_brute_force_rule.default.id}_fake"]`,
		}),
	}

	ThreatDetectionAntiBruteForceRuleCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existThreatDetectionAntiBruteForceRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                            "1",
		"rules.0.id":                         CHECKSET,
		"rules.0.anti_brute_force_rule_id":   CHECKSET,
		"rules.0.anti_brute_force_rule_name": CHECKSET,
		"rules.0.default_rule":               CHECKSET,
		"rules.0.fail_count":                 CHECKSET,
		"rules.0.forbidden_time":             CHECKSET,
		"rules.0.span":                       CHECKSET,
		"rules.0.uuid_list.#":                "1",
	}
}

var fakeThreatDetectionAntiBruteForceRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
	}
}

var ThreatDetectionAntiBruteForceRuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_anti_brute_force_rules.default",
	existMapFunc: existThreatDetectionAntiBruteForceRuleMapFunc,
	fakeMapFunc:  fakeThreatDetectionAntiBruteForceRuleMapFunc,
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
    ids = ["79d76eac-055a-492a-a5c8-eef3bac80c15"]
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
