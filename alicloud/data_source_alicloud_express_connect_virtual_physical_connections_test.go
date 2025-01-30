package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudExpressConnectVirtualPhysicalConnectionDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	testAccPreCheckWithExpressConnectUidSetting(t)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_express_connect_virtual_physical_connection.default.virtual_physical_connection_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_express_connect_virtual_physical_connection.default.virtual_physical_connection_name}_fake"`,
		}),
	}

	isConfirmedConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"is_confirmed": `false`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"is_confirmed": `true`,
		}),
	}
	businessStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"business_status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"business_status": `"FinancialLocked"`,
		}),
	}
	virtualPhysicalConnectionIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                             `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"virtual_physical_connection_ids": `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                             `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"virtual_physical_connection_ids": `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
		}),
	}
	virtualPhysicalConnectionStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"virtual_physical_connection_status": `"UnConfirmed"`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"virtual_physical_connection_status": `"Deleted"`,
		}),
	}
	VlanIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"vlan_ids": `[789]`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"vlan_ids": `[1]`,
		}),
	}

	VpconnAliUidConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"vpconn_ali_uid": `"${var.vpconn_ali_uid}"`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"vpconn_ali_uid": `"1"`,
		}),
	}
	ParentPhysicalConnectionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                           `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"parent_physical_connection_id": `"${data.alicloud_express_connect_physical_connections.default.ids.0}"`,
		}),
		fakeConfig: "",
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"vlan_ids":                           `[789]`,
			"vpconn_ali_uid":                     `"${var.vpconn_ali_uid}"`,
			"parent_physical_connection_id":      `"${data.alicloud_express_connect_physical_connections.default.ids.0}"`,
			"virtual_physical_connection_status": `"UnConfirmed"`,
			"virtual_physical_connection_ids":    `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"name_regex":                         `"${alicloud_express_connect_virtual_physical_connection.default.virtual_physical_connection_name}"`,
			"business_status":                    `"Normal"`,
			"is_confirmed":                       `false`,
		}),
		fakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"vlan_ids":                           `[1]`,
			"vpconn_ali_uid":                     `"1"`,
			"virtual_physical_connection_status": `"Deleted"`,
			"virtual_physical_connection_ids":    `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"name_regex":                         `"${alicloud_express_connect_virtual_physical_connection.default.virtual_physical_connection_name}_fake"`,
			"business_status":                    `"FinancialLocked"`,
			"is_confirmed":                       `true`,
		}),
	}

	ExpressConnectVirtualPhysicalConnectionCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, isConfirmedConf, businessStatusConf, virtualPhysicalConnectionIdsConf, virtualPhysicalConnectionStatusConf, VlanIdsConf, VpconnAliUidConf, ParentPhysicalConnectionIdConf, allConf)
}

var existExpressConnectVirtualPhysicalConnectionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                            "1",
		"names.#":                                          "1",
		"connections.#":                                    "1",
		"connections.0.id":                                 CHECKSET,
		"connections.0.access_point_id":                    CHECKSET,
		"connections.0.ad_location":                        CHECKSET,
		"connections.0.bandwidth":                          CHECKSET,
		"connections.0.business_status":                    CHECKSET,
		"connections.0.circuit_code":                       "",
		"connections.0.create_time":                        "",
		"connections.0.description":                        CHECKSET,
		"connections.0.enabled_time":                       "",
		"connections.0.end_time":                           "",
		"connections.0.expect_spec":                        "",
		"connections.0.line_operator":                      CHECKSET,
		"connections.0.loa_status":                         "",
		"connections.0.order_mode":                         CHECKSET,
		"connections.0.parent_physical_connection_ali_uid": CHECKSET,
		"connections.0.parent_physical_connection_id":      CHECKSET,
		"connections.0.peer_location":                      CHECKSET,
		"connections.0.port_number":                        "",
		"connections.0.port_type":                          CHECKSET,
		"connections.0.redundant_physical_connection_id":   "",
		"connections.0.resource_group_id":                  CHECKSET,
		"connections.0.spec":                               "50M",
		"connections.0.status":                             CHECKSET,
		"connections.0.virtual_physical_connection_id":     CHECKSET,
		"connections.0.virtual_physical_connection_name":   CHECKSET,
		"connections.0.virtual_physical_connection_status": CHECKSET,
		"connections.0.vlan_id":                            "789",
		"connections.0.vpconn_ali_uid":                     CHECKSET,
	}
}

var fakeExpressConnectVirtualPhysicalConnectionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":         "0",
		"names.#":       "0",
		"connections.#": "0",
	}
}

var ExpressConnectVirtualPhysicalConnectionCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_express_connect_virtual_physical_connections.default",
	existMapFunc: existExpressConnectVirtualPhysicalConnectionMapFunc,
	fakeMapFunc:  fakeExpressConnectVirtualPhysicalConnectionMapFunc,
}

func testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccExpressConnectVirtualPhysicalConnection%d"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}

variable "vpconn_ali_uid" {
	default = "%s"
}

resource "alicloud_express_connect_virtual_physical_connection" "default" {
  virtual_physical_connection_name = var.name
  description                      = var.name
  order_mode                       = "PayByPhysicalConnectionOwner"
  parent_physical_connection_id    = data.alicloud_express_connect_physical_connections.default.ids.0
  spec                             = "50M"
  vlan_id                          = 789
  vpconn_ali_uid                   = var.vpconn_ali_uid
}

data "alicloud_express_connect_virtual_physical_connections" "default" {
%s
}
`, rand, os.Getenv("ALICLOUD_EXPRESS_CONNECT_UID"), strings.Join(pairs, "\n   "))
	return config
}
