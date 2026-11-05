package vpcpeer_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
)

func TestAccAlicloudVpcPeerConnectionsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_peer_connection.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_peer_connection.default.id}_fake"]`,
		}),
	}
	peerConnectionNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_vpc_peer_connection.default.id}"]`,
			"peer_connection_name": `"${alicloud_vpc_peer_connection.default.peer_connection_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_vpc_peer_connection.default.id}"]`,
			"peer_connection_name": `"${alicloud_vpc_peer_connection.default.peer_connection_name}_fake"`,
		}),
	}
	vpcIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_peer_connection.default.id}"]`,
			"vpc_id": `"${alicloud_vpc_peer_connection.default.vpc_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_peer_connection.default.id}"]`,
			"vpc_id": `"${alicloud_vpc_peer_connection.default.vpc_id}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_peer_connection.default.peer_connection_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_peer_connection.default.peer_connection_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_peer_connection.default.id}"]`,
			"status": `"${alicloud_vpc_peer_connection.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_peer_connection.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_vpc_peer_connection.default.id}"]`,
			"name_regex":           `"${alicloud_vpc_peer_connection.default.peer_connection_name}"`,
			"peer_connection_name": `"${alicloud_vpc_peer_connection.default.peer_connection_name}"`,
			"status":               `"${alicloud_vpc_peer_connection.default.status}"`,
			"vpc_id":               `"${alicloud_vpc_peer_connection.default.vpc_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_vpc_peer_connection.default.id}_fake"]`,
			"name_regex":           `"${alicloud_vpc_peer_connection.default.peer_connection_name}_fake"`,
			"peer_connection_name": `"${alicloud_vpc_peer_connection.default.peer_connection_name}_fake"`,
			"status":               `"Creating"`,
			"vpc_id":               `"${alicloud_vpc_peer_connection.default.vpc_id}_fake"`,
		}),
	}
	var existAlicloudVpcPeerConnectionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"connections.#":                      "1",
			"connections.0.accepting_ali_uid":    acctest.CHECKSET,
			"connections.0.accepting_region_id":  acctest.DefaultRegionToTest,
			"connections.0.accepting_vpc_id":     acctest.CHECKSET,
			"connections.0.bandwidth":            acctest.CHECKSET,
			"connections.0.description":          fmt.Sprintf("tf-testAccPeerConnection-%d", rand),
			"connections.0.peer_connection_name": fmt.Sprintf("tf-testAccPeerConnection-%d", rand),
			"connections.0.vpc_id":               acctest.CHECKSET,
			"connections.0.create_time":          acctest.CHECKSET,
			"connections.0.id":                   acctest.CHECKSET,
			"connections.0.peer_connection_id":   acctest.CHECKSET,
			"connections.0.status":               acctest.CHECKSET,
		}
	}
	var fakeAlicloudVpcPeerConnectionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVpcPeerConnectionsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_vpc_peer_connections.default",
		ExistMapFunc: existAlicloudVpcPeerConnectionsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudVpcPeerConnectionsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudVpcPeerConnectionsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, peerConnectionNameConf, vpcIdConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcPeerConnectionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccPeerConnection-%d"
}
data "alicloud_account" "default" {}

resource "alicloud_vpc" "local" {
	vpc_name = var.name
}

resource "alicloud_vpc" "peer" {
	vpc_name = var.name
}
resource "alicloud_vpc_peer_connection" "default" {
  peer_connection_name        = var.name
  vpc_id              = alicloud_vpc.local.id
  accepting_ali_uid   = data.alicloud_account.default.id
  accepting_region_id = "%s"
  accepting_vpc_id    = alicloud_vpc.peer.id
  description         = var.name
}

data "alicloud_vpc_peer_connections" "default" {
	%s
}
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
