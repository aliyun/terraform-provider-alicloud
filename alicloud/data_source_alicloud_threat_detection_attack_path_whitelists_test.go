// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudThreatDetectionAttackPathWhitelistDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_attack_path_whitelist.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_attack_path_whitelist.default.id}_fake"]`,
		}),
	}

	WhitelistNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_threat_detection_attack_path_whitelist.default.id}"]`,
			"whitelist_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_threat_detection_attack_path_whitelist.default.id}_fake"]`,
			"whitelist_name": `"${var.name}_fake"`,
		}),
	}
	PathTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_threat_detection_attack_path_whitelist.default.id}"]`,
			"path_type": `"role_escalation"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_threat_detection_attack_path_whitelist.default.id}_fake"]`,
			"path_type": `"role_escalation_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_threat_detection_attack_path_whitelist.default.id}"]`,
			"whitelist_name": `"${var.name}"`,

			"path_type": `"role_escalation"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_threat_detection_attack_path_whitelist.default.id}_fake"]`,
			"whitelist_name": `"${var.name}_fake"`,

			"path_type": `"role_escalation_fake"`,
		}),
	}

	ThreatDetectionAttackPathWhitelistCheckInfo.dataSourceTestCheck(t, rand, idsConf, WhitelistNameConf, PathTypeConf, allConf)
}

var existThreatDetectionAttackPathWhitelistMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"whitelists.#":                          "1",
		"whitelists.0.whitelist_name":           CHECKSET,
		"whitelists.0.attack_path_whitelist_id": CHECKSET,
		"whitelists.0.remark":                   CHECKSET,
		"whitelists.0.path_type":                CHECKSET,
		"whitelists.0.whitelist_type":           CHECKSET,
		"whitelists.0.path_name":                CHECKSET,
	}
}

var fakeThreatDetectionAttackPathWhitelistMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"whitelists.#": "0",
	}
}

var ThreatDetectionAttackPathWhitelistCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_attack_path_whitelists.default",
	existMapFunc: existThreatDetectionAttackPathWhitelistMapFunc,
	fakeMapFunc:  fakeThreatDetectionAttackPathWhitelistMapFunc,
}

func testAccCheckAlicloudThreatDetectionAttackPathWhitelistSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionAttackPathWhitelist%d"
}


resource "alicloud_threat_detection_attack_path_whitelist" "default" {
  path_type      = "role_escalation"
  whitelist_type = "PART_ASSET"
  whitelist_name = var.name
  path_name      = "ecs_get_credential_by_create_login_profile"
  remark         = var.name
  attack_path_asset_list {
    instance_id    = "AliyunYundunSASReadOnlyAccess::System"
    region_id      = "cn-hangzhou"
    vendor         = 0
    asset_type     = 15
    asset_sub_type = 2
    node_type      = "end"
  }
}

data "alicloud_threat_detection_attack_path_whitelists" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
