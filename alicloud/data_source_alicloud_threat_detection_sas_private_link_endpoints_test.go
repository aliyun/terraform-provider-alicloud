package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudThreatDetectionSasPrivateLinkEndpointsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionSasPrivateLinkEndpointSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_sas_private_link_endpoint.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionSasPrivateLinkEndpointSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_sas_private_link_endpoint.default.id}_fake"]`,
		}),
	}

	nodeNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionSasPrivateLinkEndpointSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_sas_private_link_endpoint.default.id}"]`,
			"node_name":      `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionSasPrivateLinkEndpointSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_sas_private_link_endpoint.default.id}_fake"]`,
			"node_name":      `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionSasPrivateLinkEndpointSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_sas_private_link_endpoint.default.id}"]`,
			"node_name":      `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionSasPrivateLinkEndpointSourceConfig(rand, map[string]string{
			"enable_details": `"true"`,
			"ids":            `["${alicloud_threat_detection_sas_private_link_endpoint.default.id}_fake"]`,
			"node_name":      `"${var.name}_fake"`,
		}),
	}

	ThreatDetectionSasPrivateLinkEndpointCheckInfo.dataSourceTestCheck(t, rand, idsConf, nodeNameConf, allConf)
}

var existThreatDetectionSasPrivateLinkEndpointMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"endpoints.#":                   "1",
		"endpoints.0.id":                CHECKSET,
		"endpoints.0.node_name":         CHECKSET,
		"endpoints.0.vpc_id":            CHECKSET,
		"endpoints.0.security_group_id": CHECKSET,
		"endpoints.0.status":            CHECKSET,
		"endpoints.0.update_domain":     CHECKSET,
		"endpoints.0.jsrv_domain":       CHECKSET,
		"endpoints.0.region_id":         CHECKSET,
		"endpoints.0.zones.#":           "1",
	}
}

var fakeThreatDetectionSasPrivateLinkEndpointMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"endpoints.#": "0",
	}
}

var ThreatDetectionSasPrivateLinkEndpointCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_threat_detection_sas_private_link_endpoints.default",
	existMapFunc: existThreatDetectionSasPrivateLinkEndpointMapFunc,
	fakeMapFunc:  fakeThreatDetectionSasPrivateLinkEndpointMapFunc,
}

func testAccCheckAlicloudThreatDetectionSasPrivateLinkEndpointSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
    default = "tf-testAccThreatDetectionSasPrivateLinkEndpoint%d"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/21"
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
  name   = var.name
}

resource "alicloud_threat_detection_sas_private_link_endpoint" "default" {
  node_name         = var.name
  vpc_id            = alicloud_vpc.default.id
  security_group_id = alicloud_security_group.default.id
  zones {
    v_switch_id = alicloud_vswitch.default.id
    zone_id     = data.alicloud_zones.default.zones.0.id
  }
}

data "alicloud_threat_detection_sas_private_link_endpoints" "default" {
	%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
