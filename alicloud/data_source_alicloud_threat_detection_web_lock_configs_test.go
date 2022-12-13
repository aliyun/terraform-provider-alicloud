package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionWebLockConfigDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_web_lock_config.default.id}_fake"]`,
		}),
	}

	langConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
			"lang": `"zh"`,
		}),
		fakeConfig: "",
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
			"status": `"off"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
			"status": `"on"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
			"status": `"off"`,
			"lang":   `"zh"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_threat_detection_web_lock_config.default.id}_fake"]`,
			"status": `"on"`,
		}),
	}

	ThreatDetectionWebLockConfigCheckInfo.dataSourceTestCheck(t, rand, idsConf, langConf, statusConf, allConf)
}

var existThreatDetectionWebLockConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                         "1",
		"configs.#":                     "1",
		"configs.0.id":                  CHECKSET,
		"configs.0.defence_mode":        "audit",
		"configs.0.dir":                 "/tmp/",
		"configs.0.exclusive_dir":       "",
		"configs.0.exclusive_file":      "",
		"configs.0.exclusive_file_type": "",
		"configs.0.inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg",
		"configs.0.local_backup_dir":    "/usr/local/aegis/bak",
		"configs.0.mode":                "whitelist",
		"configs.0.uuid":                CHECKSET,
	}
}

var fakeThreatDetectionWebLockConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"configs.#": "0",
	}
}

var ThreatDetectionWebLockConfigCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_web_lock_configs.default",
	existMapFunc: existThreatDetectionWebLockConfigMapFunc,
	fakeMapFunc:  fakeThreatDetectionWebLockConfigMapFunc,
}

func testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionWebLockConfig%d"
}


resource "alicloud_threat_detection_web_lock_config" "default" {
  inclusive_file_type = "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg"
  uuid                = "08ae7194-4c88-4d7c-a774-0839cbe63b7c"
  mode                = "whitelist"
  local_backup_dir    = "/usr/local/aegis/bak"
  dir                 = "/tmp/"
  defence_mode        = "audit"
}

data "alicloud_threat_detection_web_lock_configs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
