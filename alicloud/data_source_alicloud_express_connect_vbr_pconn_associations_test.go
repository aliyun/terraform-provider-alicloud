package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudExpressConnectVbrPconnAssociationDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVbrPconnAssociationSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_express_connect_vbr_pconn_association.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVbrPconnAssociationSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_express_connect_vbr_pconn_association.default.id}_fake"]`,
		}),
	}
	vbrIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVbrPconnAssociationSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_express_connect_vbr_pconn_association.default.id}"]`,
			"vbr_id": `"${alicloud_express_connect_vbr_pconn_association.default.vbr_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVbrPconnAssociationSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_express_connect_vbr_pconn_association.default.id}"]`,
			"vbr_id": `"${alicloud_express_connect_vbr_pconn_association.default.vbr_id}_fake"`,
		}),
	}

	ExpressConnectVbrPconnAssociationCheckInfo.dataSourceTestCheck(t, rand, idsConf, vbrIdConf)
}

var existExpressConnectVbrPconnAssociationMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                   "1",
		"associations.#":                          "1",
		"associations.0.id":                       CHECKSET,
		"associations.0.circuit_code":             "",
		"associations.0.enable_ipv6":              CHECKSET,
		"associations.0.local_gateway_ip":         CHECKSET,
		"associations.0.local_ipv6_gateway_ip":    "",
		"associations.0.peer_gateway_ip":          CHECKSET,
		"associations.0.peer_ipv6_gateway_ip":     "",
		"associations.0.peering_ipv6_subnet_mask": "",
		"associations.0.peering_subnet_mask":      CHECKSET,
		"associations.0.physical_connection_id":   CHECKSET,
		"associations.0.status":                   CHECKSET,
		"associations.0.vbr_id":                   CHECKSET,
		"associations.0.vlan_id":                  CHECKSET,
	}
}

var fakeExpressConnectVbrPconnAssociationMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":          "0",
		"associations.#": "0",
	}
}

var ExpressConnectVbrPconnAssociationCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_express_connect_vbr_pconn_associations.default",
	existMapFunc: existExpressConnectVbrPconnAssociationMapFunc,
	fakeMapFunc:  fakeExpressConnectVbrPconnAssociationMapFunc,
}

func testAccCheckAlicloudExpressConnectVbrPconnAssociationSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccExpressConnectVbrPconnAssociation%d"
}

data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = 100
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_express_connect_vbr_pconn_association" "default" {
  peer_gateway_ip        = "10.0.0.6"
  local_gateway_ip       = "10.0.0.5"
  physical_connection_id = data.alicloud_express_connect_physical_connections.nameRegex.connections.2.id
  vbr_id                 = alicloud_express_connect_virtual_border_router.default.id
  peering_subnet_mask    = "255.255.255.252"
  vlan_id                = 1122
  enable_ipv6            = false
}

data "alicloud_express_connect_vbr_pconn_associations" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
