package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionHoneyPotDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honey_pot.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honey_pot.default.id}_fake"]`,
		}),
	}

	NodeIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_threat_detection_honey_pot.default.id}"]`,
			"node_id": `"${alicloud_threat_detection_honey_pot.default.node_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_threat_detection_honey_pot.default.id}_fake"]`,
			"node_id": `"${alicloud_threat_detection_honey_pot.default.node_id}"`,
		}),
	}
	HoneypotNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_threat_detection_honey_pot.default.id}"]`,
			"honeypot_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_threat_detection_honey_pot.default.id}_fake"]`,
			"honeypot_name": `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_threat_detection_honey_pot.default.id}"]`,
			"node_id":       `"${alicloud_threat_detection_honey_pot.default.node_id}"`,
			"honeypot_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_threat_detection_honey_pot.default.id}_fake"]`,
			"node_id":       `"${alicloud_threat_detection_honey_pot.default.node_id}"`,
			"honeypot_name": `"${var.name}_fake"`,
		}),
	}

	ThreatDetectionHoneyPotCheckInfo.dataSourceTestCheck(t, rand, idsConf, NodeIdConf, HoneypotNameConf, allConf)
}

var existThreatDetectionHoneyPotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pots.#":                     "1",
		"pots.0.id":                  CHECKSET,
		"pots.0.honeypot_id":         CHECKSET,
		"pots.0.honeypot_image_id":   CHECKSET,
		"pots.0.honeypot_image_name": CHECKSET,
		"pots.0.honeypot_name":       CHECKSET,
		"pots.0.node_id":             CHECKSET,
		"pots.0.state.#":             "1",
		"pots.0.status":              CHECKSET,
	}
}

var fakeThreatDetectionHoneyPotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pots.#": "0",
	}
}

var ThreatDetectionHoneyPotCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_honey_pots.default",
	existMapFunc: existThreatDetectionHoneyPotMapFunc,
	fakeMapFunc:  fakeThreatDetectionHoneyPotMapFunc,
}

func testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionHoneyPot%d"
}

data "alicloud_threat_detection_honeypot_images" "default" {
  name_regex = "^ruoyi"
}

resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name           = var.name
  available_probe_num = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}

resource "alicloud_threat_detection_honey_pot" "default" {
  honeypot_image_name = "ruoyi"
  honeypot_image_id   = data.alicloud_threat_detection_honeypot_images.default.images.0.honeypot_image_id
  honeypot_name       = var.name
  node_id             = alicloud_threat_detection_honeypot_node.default.id
}

data "alicloud_threat_detection_honey_pots" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
