// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsGatewayQosDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_gateway_qos.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
		}),
	}

	GatewayQosTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}"]`,
			"gateway_qos_type": `"Nat"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
			"gateway_qos_type": `"LoadBalancer"`,
		}),
	}
	GatewayQosNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}"]`,
			"gateway_qos_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
			"gateway_qos_name": `"${var.name}_fake"`,
		}),
	}
	NetworkIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_gateway_qos.default.id}"]`,
			"network_id": `"${alicloud_ens_network.defaultC7YqlT.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
			"network_id": `"${alicloud_ens_network.defaultC7YqlT.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}"]`,
			"gateway_qos_type": `"Nat"`,

			"gateway_qos_name": `"${var.name}"`,

			"network_id": `"${alicloud_ens_network.defaultC7YqlT.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
			"gateway_qos_type": `"LoadBalancer"`,

			"gateway_qos_name": `"${var.name}_fake"`,

			"network_id": `"${alicloud_ens_network.defaultC7YqlT.id}"`,
		}),
	}

	EnsGatewayQosCheckInfo.dataSourceTestCheck(t, rand, idsConf, GatewayQosTypeConf, GatewayQosNameConf, NetworkIdConf, allConf)
}

var existEnsGatewayQosMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"qos.#":                  "1",
		"qos.0.status":           CHECKSET,
		"qos.0.gateway_qos_name": CHECKSET,
		"qos.0.bandwidth_in":     CHECKSET,
		"qos.0.gateway_qos_type": CHECKSET,
		"qos.0.network_id":       CHECKSET,
		"qos.0.bandwidth_out":    CHECKSET,
		"qos.0.gateway_qos_id":   CHECKSET,
		"qos.0.creation_time":    CHECKSET,
		"qos.0.ens_region_id":    CHECKSET,
	}
}

var fakeEnsGatewayQosMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"qos.#": "0",
	}
}

var EnsGatewayQosCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_gateway_qos.default",
	existMapFunc: existEnsGatewayQosMapFunc,
	fakeMapFunc:  fakeEnsGatewayQosMapFunc,
}

func testAccCheckAlicloudEnsGatewayQosSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEnsGatewayQos%d"
}
variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "defaultC7YqlT" {
  network_name  = "镇元-网关限速测试使用"
  cidr_block    = "10.0.0.0/10"
  ens_region_id = var.ens_region_id
}


resource "alicloud_ens_gateway_qos" "default" {
  gateway_qos_name = var.name
  bandwidth_in     = "10"
  gateway_qos_type = "Nat"
  network_id       = alicloud_ens_network.defaultC7YqlT.id
  bandwidth_out    = "20"
}

data "alicloud_ens_gateway_qos" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
