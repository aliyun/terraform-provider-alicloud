package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionVulWhitelistsDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vul_whitelist.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vul_whitelist.default.id}_fake"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vul_whitelist.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionVulWhitelistSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_vul_whitelist.default.id}_fake"]`,
		}),
	}
	ThreatDetectionVulWhitelistCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

var existThreatDetectionVulWhitelistMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                         "1",
		"whitelists.#":                  "1",
		"whitelists.0.id":               acctest.CHECKSET,
		"whitelists.0.vul_whitelist_id": acctest.CHECKSET,
		"whitelists.0.whitelist":        "[{\"aliasName\":\"RHSA-2021:2260: libwebp 安全更新\",\"name\":\"RHSA-2021:2260: libwebp 安全更新\",\"type\":\"cve\"}]",
		"whitelists.0.target_info":      "{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678]}",
		"whitelists.0.reason":           acctest.CHECKSET,
	}
}

var fakeThreatDetectionVulWhitelistMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":        "0",
		"whitelists.#": "0",
	}
}

var ThreatDetectionVulWhitelistCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_vul_whitelists.default",
	ExistMapFunc: existThreatDetectionVulWhitelistMapFunc,
	FakeMapFunc:  fakeThreatDetectionVulWhitelistMapFunc,
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
