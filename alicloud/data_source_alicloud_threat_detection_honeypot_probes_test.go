package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionHoneypotProbeDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_honeypot_probe.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_honeypot_probe.default.id}_fake"]`,
		}),
	}

	ProbeTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_honeypot_probe.default.id}"]`,
			"probe_type":     `"host_probe"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_honeypot_probe.default.id}_fake"]`,
			"probe_type":     `"host_probe_fake"`,
		}),
	}
	DisplayNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_honeypot_probe.default.id}"]`,
			"display_name":   `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_honeypot_probe.default.id}_fake"]`,
			"display_name":   `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_honeypot_probe.default.id}"]`,
			"probe_type":     `"host_probe"`,
			"display_name":   `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_honeypot_probe.default.id}_fake"]`,
			"probe_type":     `"host_probe_fake"`,
			"display_name":   `"${var.name}_fake"`,
		}),
	}

	ThreatDetectionHoneypotProbeCheckInfo.dataSourceTestCheck(t, rand, idsConf, ProbeTypeConf, DisplayNameConf, allConf)
}

var existThreatDetectionHoneypotProbeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"probes.#":                   "1",
		"probes.0.id":                CHECKSET,
		"probes.0.arp":               CHECKSET,
		"probes.0.control_node_id":   CHECKSET,
		"probes.0.display_name":      CHECKSET,
		"probes.0.honeypot_probe_id": CHECKSET,
		"probes.0.ping":              CHECKSET,
		"probes.0.probe_type":        CHECKSET,
		"probes.0.status":            CHECKSET,
		"probes.0.uuid":              CHECKSET,
	}
}

var fakeThreatDetectionHoneypotProbeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"probes.#": "0",
	}
}

var ThreatDetectionHoneypotProbeCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_honeypot_probes.default",
	existMapFunc: existThreatDetectionHoneypotProbeMapFunc,
	fakeMapFunc:  fakeThreatDetectionHoneypotProbeMapFunc,
}

func testAccCheckAlicloudThreatDetectionHoneypotProbeSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionHoneypotProbe%d"
}

data "alicloud_threat_detection_assets" "default" {
    machine_types = "ecs"
    ids = ["e52c7872-29d1-4aa1-9908-0299abd53606"]
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

resource "alicloud_threat_detection_honeypot_probe" "default" {
  uuid            = data.alicloud_threat_detection_assets.default.assets.0.uuid
  probe_type      = "host_probe"
  control_node_id = alicloud_threat_detection_honeypot_node.default.id
  ping            = true
  honeypot_bind_list {
    bind_port_list {
      start_port = 80
      end_port   = 80
    }
    honeypot_id = alicloud_threat_detection_honey_pot.default.id
  }
  display_name = var.name
  arp          = true
}

data "alicloud_threat_detection_honeypot_probes" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
