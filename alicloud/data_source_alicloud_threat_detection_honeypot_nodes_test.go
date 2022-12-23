package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionHoneypotNodeDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honeypot_node.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honeypot_node.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_threat_detection_honeypot_node.default.node_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_threat_detection_honeypot_node.default.node_name}_fake"`,
		}),
	}
	NodeNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"node_name": `"${alicloud_threat_detection_honeypot_node.default.node_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"node_name": `"${alicloud_threat_detection_honeypot_node.default.node_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_threat_detection_honeypot_node.default.id}"]`,
			"name_regex": `"${alicloud_threat_detection_honeypot_node.default.node_name}"`,
			"node_name":  `"${alicloud_threat_detection_honeypot_node.default.node_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_threat_detection_honeypot_node.default.id}_fake"]`,
			"name_regex": `"${alicloud_threat_detection_honeypot_node.default.node_name}_fake"`,
			"node_name":  `"${alicloud_threat_detection_honeypot_node.default.node_name}_fake"`,
		}),
	}

	ThreatDetectionHoneypotNodeCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, NodeNameConf, allConf)
}

var existThreatDetectionHoneypotNodeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                  "1",
		"names.#":                                "1",
		"nodes.#":                                "1",
		"nodes.0.id":                             CHECKSET,
		"nodes.0.allow_honeypot_access_internet": CHECKSET,
		"nodes.0.available_probe_num":            "20",
		"nodes.0.create_time":                    CHECKSET,
		"nodes.0.node_id":                        CHECKSET,
		"nodes.0.node_name":                      CHECKSET,
		"nodes.0.security_group_probe_ip_list.#": "1",
		"nodes.0.security_group_probe_ip_list.0": "0.0.0.0/0",
		"nodes.0.status":                         CHECKSET,
	}
}

var fakeThreatDetectionHoneypotNodeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"nodes.#": "0",
	}
}

var ThreatDetectionHoneypotNodeCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_honeypot_nodes.default",
	existMapFunc: existThreatDetectionHoneypotNodeMapFunc,
	fakeMapFunc:  fakeThreatDetectionHoneypotNodeMapFunc,
}

func testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccThreatDetectionHoneypotNode%d"
}


resource "alicloud_threat_detection_honeypot_node" "default" {
  node_name           = var.name
  available_probe_num = 20
  security_group_probe_ip_list = ["0.0.0.0/0"]
}

data "alicloud_threat_detection_honeypot_nodes" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
