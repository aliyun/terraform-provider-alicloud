package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterInterfaceDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectRouterInterfaceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_express_connect_router_interface.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectRouterInterfaceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_express_connect_router_interface.default.id}_fake"]`,
		}),
	}
	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectRouterInterfaceSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_express_connect_router_interface.default.router_interface_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectRouterInterfaceSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_express_connect_router_interface.default.router_interface_name}_fake"`,
		}),
	}

	includeReservationDataConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectRouterInterfaceSourceConfig(rand, map[string]string{
			"ids":                      `["${alicloud_express_connect_router_interface.default.id}"]`,
			"include_reservation_data": `"false"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectRouterInterfaceSourceConfig(rand, map[string]string{
			"ids":                      `["${alicloud_express_connect_router_interface.default.id}"]`,
			"name_regex":               `"${alicloud_express_connect_router_interface.default.router_interface_name}"`,
			"include_reservation_data": `"false"`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectRouterInterfaceSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_express_connect_router_interface.default.id}_fake"]`,
			"name_regex": `"${alicloud_express_connect_router_interface.default.router_interface_name}_fake"`,
		}),
	}

	ExpressConnectRouterInterfaceCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConfig, includeReservationDataConfig, allConf)
}

var existExpressConnectRouterInterfaceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                           "1",
		"names.#":                                         "1",
		"interfaces.#":                                    "1",
		"interfaces.0.id":                                 CHECKSET,
		"interfaces.0.access_point_id":                    "",
		"interfaces.0.bandwidth":                          CHECKSET,
		"interfaces.0.business_status":                    CHECKSET,
		"interfaces.0.connected_time":                     "",
		"interfaces.0.create_time":                        CHECKSET,
		"interfaces.0.cross_border":                       CHECKSET,
		"interfaces.0.description":                        fmt.Sprintf("tf-testAccRouterInterface%d", rand),
		"interfaces.0.end_time":                           "",
		"interfaces.0.has_reservation_data":               CHECKSET,
		"interfaces.0.hc_rate":                            CHECKSET,
		"interfaces.0.hc_threshold":                       "",
		"interfaces.0.health_check_source_ip":             "",
		"interfaces.0.health_check_target_ip":             "",
		"interfaces.0.opposite_access_point_id":           "",
		"interfaces.0.opposite_bandwidth":                 CHECKSET,
		"interfaces.0.opposite_interface_business_status": CHECKSET,
		"interfaces.0.opposite_interface_id":              "",
		"interfaces.0.opposite_interface_owner_id":        "",
		"interfaces.0.opposite_interface_spec":            CHECKSET,
		"interfaces.0.opposite_interface_status":          "",
		"interfaces.0.opposite_region_id":                 CHECKSET,
		"interfaces.0.opposite_router_id":                 "",
		"interfaces.0.opposite_router_type":               CHECKSET,
		"interfaces.0.opposite_vpc_instance_id":           "",
		"interfaces.0.payment_type":                       CHECKSET,
		"interfaces.0.reservation_active_time":            "",
		"interfaces.0.reservation_bandwidth":              "",
		"interfaces.0.reservation_internet_charge_type":   "",
		"interfaces.0.reservation_order_type":             "",
		"interfaces.0.role":                               "InitiatingSide",
		"interfaces.0.router_id":                          CHECKSET,
		"interfaces.0.router_interface_id":                CHECKSET,
		"interfaces.0.router_interface_name":              CHECKSET,
		"interfaces.0.router_type":                        "VRouter",
		"interfaces.0.spec":                               "Mini.2",
		"interfaces.0.status":                             CHECKSET,
		"interfaces.0.vpc_instance_id":                    CHECKSET,
	}
}

var fakeExpressConnectRouterInterfaceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":        "0",
		"names.#":      "0",
		"interfaces.#": "0",
	}
}

var ExpressConnectRouterInterfaceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_express_connect_router_interfaces.default",
	existMapFunc: existExpressConnectRouterInterfaceMapFunc,
	fakeMapFunc:  fakeExpressConnectRouterInterfaceMapFunc,
}

func testAccCheckAlicloudExpressConnectRouterInterfaceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRouterInterface%d"
}

resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_express_connect_router_interface" "default" {
  description           = var.name
  opposite_region_id    = "%s"
  router_id             = alicloud_vpc.default.router_id
  role                  = "InitiatingSide"
  router_type           = "VRouter"
  payment_type          = "PayAsYouGo"
  router_interface_name = var.name
  spec                  = "Mini.2"
}

data "alicloud_express_connect_router_interfaces" "default" {
%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
