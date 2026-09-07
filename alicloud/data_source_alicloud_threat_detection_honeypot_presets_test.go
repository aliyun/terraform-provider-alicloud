package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionHoneypotPresetDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honeypot_preset.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honeypot_preset.default.id}_fake"]`,
		}),
	}
	PresetNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_honeypot_preset.default.id}"]`,
			"preset_name": `"${alicloud_threat_detection_honeypot_preset.default.preset_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_threat_detection_honeypot_preset.default.id}"]`,
			"preset_name": `"${alicloud_threat_detection_honeypot_preset.default.preset_name}_fake"`,
		}),
	}
	NodeIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_threat_detection_honeypot_preset.default.id}"]`,
			"node_id": `"${alicloud_threat_detection_honeypot_preset.default.node_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_threat_detection_honeypot_preset.default.id}"]`,
			"node_id": `"${alicloud_threat_detection_honeypot_preset.default.node_id}_fake"`,
		}),
	}
	HoneypotImageNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids":                 `["${alicloud_threat_detection_honeypot_preset.default.id}"]`,
			"honeypot_image_name": `"shiro"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids":                 `["${alicloud_threat_detection_honeypot_preset.default.id}_fake"]`,
			"honeypot_image_name": `"metabase"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids":                 `["${alicloud_threat_detection_honeypot_preset.default.id}"]`,
			"preset_name":         `"${alicloud_threat_detection_honeypot_preset.default.preset_name}"`,
			"node_id":             `"${alicloud_threat_detection_honeypot_preset.default.node_id}"`,
			"honeypot_image_name": `"shiro"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand, map[string]string{
			"ids":                 `["${alicloud_threat_detection_honeypot_preset.default.id}"]`,
			"preset_name":         `"${alicloud_threat_detection_honeypot_preset.default.preset_name}_fake"`,
			"node_id":             `"${alicloud_threat_detection_honeypot_preset.default.node_id}_fake"`,
			"honeypot_image_name": `"metabase"`,
		}),
	}

	ThreatDetectionHoneypotPresetCheckInfo.dataSourceTestCheck(t, rand, idsConf, PresetNameConf, NodeIdConf, HoneypotImageNameConf, allConf)
}

var existThreatDetectionHoneypotPresetMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"presets.#":                        "1",
		"presets.0.id":                     CHECKSET,
		"presets.0.honeypot_image_name":    "shiro",
		"presets.0.honeypot_preset_id":     CHECKSET,
		"presets.0.meta.#":                 "1",
		"presets.0.meta.0.portrait_option": "true",
		"presets.0.meta.0.burp":            "open",
		"presets.0.node_id":                CHECKSET,
		"presets.0.preset_name":            CHECKSET,
	}
}

var fakeThreatDetectionHoneypotPresetMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"presets.#": "0",
	}
}

var ThreatDetectionHoneypotPresetCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_honeypot_presets.default",
	existMapFunc: existThreatDetectionHoneypotPresetMapFunc,
	fakeMapFunc:  fakeThreatDetectionHoneypotPresetMapFunc,
}

func testAccCheckAlicloudThreatDetectionHoneypotPresetSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionHoneypotPreset%d"
}

resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name           = var.name
  available_probe_num = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}

resource "alicloud_threat_detection_honeypot_preset" "default" {
  honeypot_image_name = "shiro"
  meta {
	portrait_option = true
	burp = "open"
  }
  node_id             = alicloud_threat_detection_honeypot_node.default.id
  preset_name         = var.name
}

data "alicloud_threat_detection_honeypot_presets" "default" {
enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
