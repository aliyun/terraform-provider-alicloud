package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudVirtualBorderRoutersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	physicalConnectionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"physical_connection_id": `"${alicloud_virtual_border_router.default.physical_connection_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"physical_connection_id": `"${alicloud_virtual_border_router.default.physical_connection_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"physical_connection_id": `"${alicloud_virtual_border_router.default.physical_connection_id}"`,
			"status":                 `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"physical_connection_id": `"${alicloud_virtual_border_router.default.physical_connection_id}_fake"`,
			"status":                 `"terminated"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_virtual_border_router.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_virtual_border_router.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_virtual_border_router.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_virtual_border_router.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"physical_connection_id": `"${alicloud_virtual_border_router.default.physical_connection_id}"`,
			"name_regex":             `"${alicloud_virtual_border_router.default.name}"`,
			"ids":                    `["${alicloud_virtual_border_router.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand, map[string]string{
			"physical_connection_id": `"${alicloud_virtual_border_router.default.physical_connection_id}_fake"`,
			"name_regex":             `"${alicloud_virtual_border_router.default.name}_fake"`,
			"ids":                    `["${alicloud_virtual_border_router.default.id}_fake"]`,
		}),
	}
	var virtualBorderRoutersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_virtual_border_routers.default",
		existMapFunc: existVirtualBorderRoutersMapFunc,
		fakeMapFunc:  fakeVirtualBorderRoutersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithPhysicalConnectionSetting(t)
	}
	virtualBorderRoutersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck,
		physicalConnectionIdConf, statusConf, nameRegexConf,
		idsConf, allConf)
}

func testAccCheckAlicloudVirtualBorderRoutersDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCheckAlicloudVirtualBorderRoutersDataSourceBasic-%d"
}

resource "alicloud_virtual_border_router" "default" {
	physical_connection_id = "%s"
	vlan_id                = 2500
	local_gateway_ip       = "10.0.0.3"
	peer_gateway_ip        = "10.0.0.4"
	peering_subnet_mask    = "255.255.255.0"
	name                   = "${var.name}"
	description            = "example"
}

data "alicloud_virtual_border_routers" "default" {
	%s
}`, rand, os.Getenv("ALICLOUD_PHYSICAL_CONNECTION_ID"), strings.Join(pairs, "\n  "))
	return config
}

var existVirtualBorderRoutersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vbrs.0.id":                            CHECKSET,
		"vbrs.0.status":                        "active",
		"vbrs.0.name":                          fmt.Sprintf("tf-testAccCheckAlicloudVirtualBorderRoutersDataSourceBasic-%d", rand),
		"vbrs.0.description":                   "example",
		"vbrs.0.vlan_id":                       "2500",
		"vbrs.0.route_table_id":                CHECKSET,
		"vbrs.0.vlan_interface_id":             CHECKSET,
		"vbrs.0.local_gateway_ip":              "10.0.0.3",
		"vbrs.0.peer_gateway_ip":               "10.0.0.4",
		"vbrs.0.peering_subnet_mask":           "255.255.255.0",
		"vbrs.0.physical_connection_id":        CHECKSET,
		"vbrs.0.physical_connection_owner_uid": CHECKSET,
		"vbrs.0.access_point_id":               CHECKSET,
		"vbrs.0.creation_time":                 CHECKSET,
		"vbrs.0.circuit_code":                  "",
	}
}

var fakeVirtualBorderRoutersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"vbrs.#":  "0",
	}
}
