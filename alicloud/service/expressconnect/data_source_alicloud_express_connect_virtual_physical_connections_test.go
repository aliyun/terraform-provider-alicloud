package expressconnect_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudExpressConnectVirtualPhysicalConnectionDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	testAccPreCheckWithExpressConnectUidSetting(t)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_express_connect_virtual_physical_connection.default.virtual_physical_connection_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_express_connect_virtual_physical_connection.default.virtual_physical_connection_name}_fake"`,
		}),
	}

	isConfirmedConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"is_confirmed": `false`,
		}),
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"is_confirmed": `true`,
		}),
	}
	businessStatusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"business_status": `"Normal"`,
		}),
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"business_status": `"FinancialLocked"`,
		}),
	}
	virtualPhysicalConnectionIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                             `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"virtual_physical_connection_ids": `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                             `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"virtual_physical_connection_ids": `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
		}),
	}
	virtualPhysicalConnectionStatusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"virtual_physical_connection_status": `"UnConfirmed"`,
		}),
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                                `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"virtual_physical_connection_status": `"Deleted"`,
		}),
	}
	VlanIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"vlan_ids": `[789]`,
		}),
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"vlan_ids": `[1]`,
		}),
	}

	VpconnAliUidConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"vpconn_ali_uid": `"${var.vpconn_ali_uid}"`,
		}),
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_express_connect_virtual_physical_connection.default.id}_fake"]`,
			"vpconn_ali_uid": `"1"`,
		}),
	}
	ParentPhysicalConnectionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
			"ids":                           `["${alicloud_express_connect_virtual_physical_connection.default.id}"]`,
			"parent_physical_connection_id": `"${data.alicloud_express_connect_physical_connections.default.ids.0}"`,
		}),
		FakeConfig: "",
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
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
		FakeConfig: testAccCheckAlicloudExpressConnectVirtualPhysicalConnectionSourceConfig(rand, map[string]string{
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

	ExpressConnectVirtualPhysicalConnectionCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, isConfirmedConf, businessStatusConf, virtualPhysicalConnectionIdsConf, virtualPhysicalConnectionStatusConf, VlanIdsConf, VpconnAliUidConf, ParentPhysicalConnectionIdConf, allConf)
}

var existExpressConnectVirtualPhysicalConnectionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                            "1",
		"names.#":                                          "1",
		"connections.#":                                    "1",
		"connections.0.id":                                 acctest.CHECKSET,
		"connections.0.access_point_id":                    acctest.CHECKSET,
		"connections.0.ad_location":                        acctest.CHECKSET,
		"connections.0.bandwidth":                          acctest.CHECKSET,
		"connections.0.business_status":                    acctest.CHECKSET,
		"connections.0.circuit_code":                       "",
		"connections.0.create_time":                        "",
		"connections.0.description":                        acctest.CHECKSET,
		"connections.0.enabled_time":                       "",
		"connections.0.end_time":                           "",
		"connections.0.expect_spec":                        "",
		"connections.0.line_operator":                      acctest.CHECKSET,
		"connections.0.loa_status":                         "",
		"connections.0.order_mode":                         acctest.CHECKSET,
		"connections.0.parent_physical_connection_ali_uid": acctest.CHECKSET,
		"connections.0.parent_physical_connection_id":      acctest.CHECKSET,
		"connections.0.peer_location":                      acctest.CHECKSET,
		"connections.0.port_number":                        "",
		"connections.0.port_type":                          acctest.CHECKSET,
		"connections.0.redundant_physical_connection_id":   "",
		"connections.0.resource_group_id":                  acctest.CHECKSET,
		"connections.0.spec":                               "50M",
		"connections.0.status":                             acctest.CHECKSET,
		"connections.0.virtual_physical_connection_id":     acctest.CHECKSET,
		"connections.0.virtual_physical_connection_name":   acctest.CHECKSET,
		"connections.0.virtual_physical_connection_status": acctest.CHECKSET,
		"connections.0.vlan_id":                            "789",
		"connections.0.vpconn_ali_uid":                     acctest.CHECKSET,
	}
}

var fakeExpressConnectVirtualPhysicalConnectionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":         "0",
		"names.#":       "0",
		"connections.#": "0",
	}
}

var ExpressConnectVirtualPhysicalConnectionCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_express_connect_virtual_physical_connections.default",
	ExistMapFunc: existExpressConnectVirtualPhysicalConnectionMapFunc,
	FakeMapFunc:  fakeExpressConnectVirtualPhysicalConnectionMapFunc,
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
