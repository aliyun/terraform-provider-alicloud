package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsNetworkPeerConnectionsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNetworkPeerConnectionsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_network_peer_connection.default.network_peer_connection_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNetworkPeerConnectionsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_network_peer_connection.default.network_peer_connection_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNetworkPeerConnectionsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_network_peer_connection.default.network_peer_connection_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNetworkPeerConnectionsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ens_network_peer_connection.default.network_peer_connection_name}_fake"`,
		}),
	}

	EnsNetworkPeerConnectionCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, allConf)
}

func testAccCheckAlicloudEnsNetworkPeerConnectionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
    default = "tf-testAccEnsNetworkPeerConnectionsDataSource%d"
}

resource "alicloud_ens_network" "default" {
  network_name  = "${var.name}-net1"
  description   = var.name
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "accepting" {
  network_name  = "${var.name}-net2"
  description   = var.name
  cidr_block    = "192.168.3.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network_peer_connection" "default" {
  network_id                    = alicloud_ens_network.default.id
  accepting_network_id          = alicloud_ens_network.accepting.id
  network_peer_connection_name  = var.name
  description                   = var.name
}

data "alicloud_ens_network_peer_connections" "default" {
    %s
}
`, rand, strings.Join(pairs, "\n    "))
	return config
}

var existEnsNetworkPeerConnectionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"connections.#": "1",
		"connections.0.network_peer_connection_name": fmt.Sprintf("tf-testAccEnsNetworkPeerConnectionsDataSource%d", rand),
	}
}

var fakeEnsNetworkPeerConnectionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"connections.#": "0",
	}
}

var EnsNetworkPeerConnectionCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_network_peer_connections.default",
	existMapFunc: existEnsNetworkPeerConnectionMapFunc,
	fakeMapFunc:  fakeEnsNetworkPeerConnectionMapFunc,
}
