package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsNetworkRouteTablesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)

	routeTableNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNetworkRouteTablesSourceConfig(rand, map[string]string{
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNetworkRouteTablesSourceConfig(rand, map[string]string{
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNetworkRouteTablesSourceConfig(rand, map[string]string{
			"network_id": `"${alicloud_ens_network_route_table.default.network_id}"`,
		}),
		// The network_id filter returns the network's auto-created system router table
		// (AssociateType=VSwitch, Type=System), which carries an empty route_table_name,
		// not the created Gateway route table. Override the shared existMapFunc's
		// route_table_name assertion so the network_id step checks the real returned value.
		existChangMap: map[string]string{
			"tables.0.route_table_name": "",
		},
		fakeConfig: testAccCheckAlicloudEnsNetworkRouteTablesSourceConfig(rand, map[string]string{
			"network_id": `"fake-network-id"`,
		}),
	}

	EnsNetworkRouteTableCheckInfo.dataSourceTestCheck(t, rand, routeTableNameConf, allConf)
}

func testAccCheckAlicloudEnsNetworkRouteTablesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-ens-net-route-tables-ds-%d"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_network_route_table" "default" {
  associate_type                 = "Gateway"
  description                    = "tf-test-route-table-ds"
  is_default_gateway_route_table = true
  network_id                     = alicloud_ens_network.default.id
  route_table_name               = var.name
}

data "alicloud_ens_network_route_tables" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existEnsNetworkRouteTableMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"tables.#":                  "1",
		"tables.0.route_table_name": fmt.Sprintf("tf-testacc-ens-net-route-tables-ds-%d", rand),
	}
}

var fakeEnsNetworkRouteTableMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"tables.#": "0",
	}
}

var EnsNetworkRouteTableCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_network_route_tables.default",
	existMapFunc: existEnsNetworkRouteTableMapFunc,
	fakeMapFunc:  fakeEnsNetworkRouteTableMapFunc,
}
