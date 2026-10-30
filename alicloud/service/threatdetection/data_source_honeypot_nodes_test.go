package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudThreatDetectionHoneypotNodesDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honeypot_node.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_honeypot_node.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_threat_detection_honeypot_node.default.node_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_threat_detection_honeypot_node.default.node_name}_fake"`,
		}),
	}
	NodeNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"node_name": `"${alicloud_threat_detection_honeypot_node.default.node_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"node_name": `"${alicloud_threat_detection_honeypot_node.default.node_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_threat_detection_honeypot_node.default.id}"]`,
			"name_regex": `"${alicloud_threat_detection_honeypot_node.default.node_name}"`,
			"node_name":  `"${alicloud_threat_detection_honeypot_node.default.node_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudThreatDetectionHoneypotNodeSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_threat_detection_honeypot_node.default.id}_fake"]`,
			"name_regex": `"${alicloud_threat_detection_honeypot_node.default.node_name}_fake"`,
			"node_name":  `"${alicloud_threat_detection_honeypot_node.default.node_name}_fake"`,
		}),
	}

	ThreatDetectionHoneypotNodeCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, NodeNameConf, allConf)
}

var existThreatDetectionHoneypotNodeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                  "1",
		"names.#":                                "1",
		"nodes.#":                                "1",
		"nodes.0.id":                             acctest.CHECKSET,
		"nodes.0.allow_honeypot_access_internet": acctest.CHECKSET,
		"nodes.0.available_probe_num":            "20",
		"nodes.0.create_time":                    acctest.CHECKSET,
		"nodes.0.node_id":                        acctest.CHECKSET,
		"nodes.0.node_name":                      acctest.CHECKSET,
		"nodes.0.security_group_probe_ip_list.#": "1",
		"nodes.0.security_group_probe_ip_list.0": "0.0.0.0/0",
		"nodes.0.status":                         acctest.CHECKSET,
	}
}

var fakeThreatDetectionHoneypotNodeMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"nodes.#": "0",
	}
}

var ThreatDetectionHoneypotNodeCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_threat_detection_honeypot_nodes.default",
	ExistMapFunc: existThreatDetectionHoneypotNodeMapFunc,
	FakeMapFunc:  fakeThreatDetectionHoneypotNodeMapFunc,
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
