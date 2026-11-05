package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionHoneyPotsDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honey_pot.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honey_pot.default.id}_fake"]`,
		}),
	}

	NodeIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_threat_detection_honey_pot.default.id}"]`,
			"node_id": `"${alicloud_threat_detection_honey_pot.default.node_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_threat_detection_honey_pot.default.id}_fake"]`,
			"node_id": `"${alicloud_threat_detection_honey_pot.default.node_id}"`,
		}),
	}
	HoneypotNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_threat_detection_honey_pot.default.id}"]`,
			"honeypot_name": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_threat_detection_honey_pot.default.id}_fake"]`,
			"honeypot_name": `"${var.name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_threat_detection_honey_pot.default.id}"]`,
			"node_id":       `"${alicloud_threat_detection_honey_pot.default.node_id}"`,
			"honeypot_name": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneyPotSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_threat_detection_honey_pot.default.id}_fake"]`,
			"node_id":       `"${alicloud_threat_detection_honey_pot.default.node_id}"`,
			"honeypot_name": `"${var.name}_fake"`,
		}),
	}

	ThreatDetectionHoneyPotCheckInfo.DataSourceTestCheck(t, rand, idsConf, NodeIdConf, HoneypotNameConf, allConf)
}

var existThreatDetectionHoneyPotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pots.#":                     "1",
		"pots.0.id":                  acctest.CHECKSET,
		"pots.0.honeypot_id":         acctest.CHECKSET,
		"pots.0.honeypot_image_id":   acctest.CHECKSET,
		"pots.0.honeypot_image_name": acctest.CHECKSET,
		"pots.0.honeypot_name":       acctest.CHECKSET,
		"pots.0.node_id":             acctest.CHECKSET,
		"pots.0.state.#":             "1",
		"pots.0.status":              acctest.CHECKSET,
	}
}

var fakeThreatDetectionHoneyPotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pots.#": "0",
	}
}

var ThreatDetectionHoneyPotCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_honey_pots.default",
	ExistMapFunc: existThreatDetectionHoneyPotMapFunc,
	FakeMapFunc:  fakeThreatDetectionHoneyPotMapFunc,
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
