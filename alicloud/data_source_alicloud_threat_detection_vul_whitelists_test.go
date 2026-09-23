package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionVulWhitelistDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vul_whitelist.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vul_whitelist.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vul_whitelist.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vul_whitelist.default.id}_fake"]`,
		}),
	}
	ThreatDetectionVulWhitelistCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existThreatDetectionVulWhitelistMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                         "1",
		"whitelists.#":                  "1",
		"whitelists.0.id":               CHECKSET,
		"whitelists.0.vul_whitelist_id": CHECKSET,
		"whitelists.0.whitelist":        "[{\"aliasName\":\"RHSA-2021:2260: libwebp 安全更新\",\"name\":\"RHSA-2021:2260: libwebp 安全更新\",\"type\":\"cve\"}]",
		"whitelists.0.target_info":      "{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678]}",
		"whitelists.0.reason":           CHECKSET,
	}
}

var fakeThreatDetectionVulWhitelistMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":        "0",
		"whitelists.#": "0",
	}
}

var ThreatDetectionVulWhitelistCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_vul_whitelists.default",
	existMapFunc: existThreatDetectionVulWhitelistMapFunc,
	fakeMapFunc:  fakeThreatDetectionVulWhitelistMapFunc,
}

func testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccThreatDetectionVulWhitelist%d"
	}

	resource "alicloud_threat_detection_vul_whitelist" "default" {
  		whitelist   = "[{\"aliasName\":\"RHSA-2021:2260: libwebp 安全更新\",\"name\":\"RHSA-2021:2260: libwebp 安全更新\",\"type\":\"cve\"}]"
  		target_info = "{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678]}"
  		reason      = var.name
	}

	data "alicloud_threat_detection_vul_whitelists" "default" {
		%s
	}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
