package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionWebLockConfigsDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_web_lock_config.default.id}_fake"]`,
		}),
	}

	langConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
			"lang": `"zh"`,
		}),
		FakeConfig: "",
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
			"status": `"off"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
			"status": `"on"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_threat_detection_web_lock_config.default.id}"]`,
			"status": `"off"`,
			"lang":   `"zh"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionWebLockConfigSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_threat_detection_web_lock_config.default.id}_fake"]`,
			"status": `"on"`,
		}),
	}

	ThreatDetectionWebLockConfigCheckInfo.DataSourceTestCheck(t, rand, idsConf, langConf, statusConf, allConf)
}

var existThreatDetectionWebLockConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                         "1",
		"configs.#":                     "1",
		"configs.0.id":                  acctest.CHECKSET,
		"configs.0.defence_mode":        "audit",
		"configs.0.dir":                 "/tmp/",
		"configs.0.exclusive_dir":       "",
		"configs.0.exclusive_file":      "",
		"configs.0.exclusive_file_type": "",
		"configs.0.inclusive_file_type": "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg",
		"configs.0.local_backup_dir":    "/usr/local/aegis/bak",
		"configs.0.mode":                "whitelist",
		"configs.0.uuid":                acctest.CHECKSET,
	}
}

var fakeThreatDetectionWebLockConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":     "0",
		"configs.#": "0",
	}
}

var ThreatDetectionWebLockConfigCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_web_lock_configs.default",
	ExistMapFunc: existThreatDetectionWebLockConfigMapFunc,
	FakeMapFunc:  fakeThreatDetectionWebLockConfigMapFunc,
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

data "alicloud_threat_detection_assets" "default" {
  machine_types = "ecs"
}

resource "alicloud_threat_detection_web_lock_config" "default" {
  inclusive_file_type = "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg"
  uuid                = data.alicloud_threat_detection_assets.default.ids.0
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
