package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEnsNetworkRouteTablesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	routeTableIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"route_table_id": `"${alicloud_ens_network_route_table.default.route_table_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"route_table_id": `"${alicloud_ens_network_route_table.default.route_table_id}_fake"`,
		}),
	}

	routeTableNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}_fake"`,
		}),
	}

	networkIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"network_id": `"${alicloud_ens_network_route_table.default.network_id}"`,
		}),
		existChangMap: map[string]string{
			"tables.#": "2",
		},
		fakeConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"network_id": `"${alicloud_ens_network_route_table.default.network_id}_fake"`,
		}),
	}

	associateTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"associate_type": `"${alicloud_ens_network_route_table.default.associate_type}"`,
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"associate_type": `"InvalidType"`,
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}_fake"`,
		}),
	}

	routeTableTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"route_table_type": `"${alicloud_ens_network_route_table.default.route_table_type}"`,
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"route_table_type": `"InvalidType"`,
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"route_table_id": `"${alicloud_ens_network_route_table.default.route_table_id}"`,
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}"`,
			"network_id": `"${alicloud_ens_network_route_table.default.network_id}"`,
			"associate_type": `"${alicloud_ens_network_route_table.default.associate_type}"`,
		}),
		fakeConfig: testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand, map[string]string{
			"route_table_id": `"${alicloud_ens_network_route_table.default.route_table_id}_fake"`,
			"route_table_name": `"${alicloud_ens_network_route_table.default.route_table_name}_fake"`,
			"network_id": `"${alicloud_ens_network_route_table.default.network_id}_fake"`,
			"associate_type": `"InvalidType"`,
		}),
	}

	EnsNetworkRouteTableCheckInfo.dataSourceTestCheck(t, rand, routeTableIdConf, routeTableNameConf, networkIdConf, associateTypeConf, routeTableTypeConf, allConf)
}

func testAccCheckAliCloudEnsNetworkRouteTableSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEnsNetworkRouteTablesDataSource%d"
}

variable "ens_region_id" {
  default = "cn-hangzhou-31"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_network_route_table" "default" {
  network_id       = alicloud_ens_network.default.id
  associate_type   = "Gateway"
  route_table_name = var.name
  description      = var.name
}

data "alicloud_ens_network_route_tables" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existEnsNetworkRouteTableMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"tables.#":                          "1",
		"tables.0.route_table_name":         fmt.Sprintf("tf-testAccEnsNetworkRouteTablesDataSource%d", rand),
		"tables.0.description":              fmt.Sprintf("tf-testAccEnsNetworkRouteTablesDataSource%d", rand),
		"tables.0.associate_type":           "Gateway",
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
