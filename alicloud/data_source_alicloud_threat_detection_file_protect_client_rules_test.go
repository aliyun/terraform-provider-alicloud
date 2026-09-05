package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudThreatDetectionFileProtectClientRulesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids": `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids": `["${alicloud_threat_detection_file_protect_client_rule.default.id}_fake"]`,
		}),
	}

	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":       `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"rule_name": `"${alicloud_threat_detection_file_protect_client_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":       `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"rule_name": `"${alicloud_threat_detection_file_protect_client_rule.default.rule_name}_fake"`,
		}),
	}

	ruleActionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":         `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"rule_action": `"pass"`,
		}),
		fakeConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":         `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"rule_action": `"block"`,
		}),
	}

	platformConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":      `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"platform": `"linux"`,
		}),
		fakeConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":      `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"platform": `"windows"`,
		}),
	}

	alertLevelConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":         `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"alert_level": `"1"`,
		}),
		fakeConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":         `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"alert_level": `"3"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":         `["${alicloud_threat_detection_file_protect_client_rule.default.id}"]`,
			"rule_name":   `"${alicloud_threat_detection_file_protect_client_rule.default.rule_name}"`,
			"rule_action": `"pass"`,
			"platform":    `"linux"`,
			"alert_level": `"1"`,
		}),
		fakeConfig: testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name, map[string]string{
			"ids":         `["${alicloud_threat_detection_file_protect_client_rule.default.id}_fake"]`,
			"rule_name":   `"${alicloud_threat_detection_file_protect_client_rule.default.rule_name}_fake"`,
			"rule_action": `"block"`,
			"platform":    `"windows"`,
			"alert_level": `"3"`,
		}),
	}

	var existDataMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"rules.#":                 "1",
			"rules.0.id":              CHECKSET,
			"rules.0.rule_name":       name,
			"rules.0.rule_action":     "pass",
			"rules.0.status":          "0",
			"rules.0.platform":        "linux",
			"rules.0.alert_level":     "1",
			"rules.0.file_paths.#":    "2",
			"rules.0.file_ops.#":      "1",
			"rules.0.proc_paths.#":    "1",
			"rules.0.file_types.#":    "1",
			"rules.0.exclude_users.#": "1",
		}
	}

	var fakeDataMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"rules.#": "0",
		}
	}

	var fileProtectClientRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_threat_detection_file_protect_client_rules.default",
		existMapFunc: existDataMapFunc,
		fakeMapFunc:  fakeDataMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		testAccPreCheck(t)
	}

	fileProtectClientRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, ruleNameConf, ruleActionConf, platformConf, alertLevelConf, allConf)
}

func testAccCheckAliCloudThreatDetectionFileProtectClientRulesDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_threat_detection_file_protect_client_rule" "default" {
  rule_name     = var.name
  rule_action   = "pass"
  status        = 0
  platform      = "linux"
  alert_level   = 1
  file_paths    = ["/opt/a", "/tmp/d"]
  file_ops      = ["READ"]
  proc_paths    = ["/usr/bin/java"]
  file_types    = ["sh"]
  exclude_users = ["root"]
}

data "alicloud_threat_detection_file_protect_client_rules" "default" {
  %s
}
`, name, strings.Join(pairs, "\n  "))
	return config
}
