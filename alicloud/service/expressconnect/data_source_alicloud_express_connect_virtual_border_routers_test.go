package expressconnect_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudExpressConnectVirtualBorderRoutersDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "data.alicloud_express_connect_virtual_border_routers.default"
	rand := tfacctest.RandIntRange(1, 2999)
	name := fmt.Sprintf("tf-testAccExpressConnectVirtualBorderRoutersTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectVirtualBorderRoutersDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_virtual_border_router.default.id}-fake"},
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
			"status": "active",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
			"status": "terminated",
		}),
	}
	filterConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"filter": []map[string]interface{}{
				{
					"key":    "PhysicalConnectionId",
					"values": []string{"${alicloud_express_connect_virtual_border_router.default.physical_connection_id}"},
				},
				{
					"key":    "VbrId",
					"values": []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"filter": []map[string]interface{}{
				{
					"key":    "PhysicalConnectionId",
					"values": []string{"${alicloud_express_connect_virtual_border_router.default.physical_connection_id}-fake"},
				},
				{
					"key":    "VbrId",
					"values": []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
				},
			},
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
			"status":     "active",
			"filter": []map[string]interface{}{
				{
					"key":    "PhysicalConnectionId",
					"values": []string{"${alicloud_express_connect_virtual_border_router.default.physical_connection_id}"},
				},
				{
					"key":    "VbrId",
					"values": []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alicloud_express_connect_virtual_border_router.default.id}-fake"},
			"status":     "terminated",
		}),
	}
	var existExpressConnectVirtualBorderRoutersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                         "1",
			"ids.0":                                         acctest.CHECKSET,
			"names.#":                                       "1",
			"names.0":                                       name,
			"routers.#":                                     "1",
			"routers.0.id":                                  acctest.CHECKSET,
			"routers.0.access_point_id":                     acctest.CHECKSET,
			"routers.0.activation_time":                     acctest.CHECKSET,
			"routers.0.circuit_code":                        "",
			"routers.0.cloud_box_instance_id":               "",
			"routers.0.create_time":                         acctest.CHECKSET,
			"routers.0.description":                         "",
			"routers.0.detect_multiplier":                   "10",
			"routers.0.ecc_id":                              "",
			"routers.0.enable_ipv6":                         "false",
			"routers.0.status":                              "active",
			"routers.0.local_gateway_ip":                    "10.0.0.1",
			"routers.0.local_ipv6_gateway_ip":               "",
			"routers.0.min_rx_interval":                     "1000",
			"routers.0.min_tx_interval":                     "1000",
			"routers.0.payment_vbr_expire_time":             "",
			"routers.0.peer_gateway_ip":                     "10.0.0.2",
			"routers.0.peer_ipv6_gateway_ip":                "",
			"routers.0.peering_ipv6_subnet_mask":            "",
			"routers.0.peering_subnet_mask":                 "255.255.255.252",
			"routers.0.physical_connection_business_status": acctest.CHECKSET,
			"routers.0.physical_connection_id":              acctest.CHECKSET,
			"routers.0.physical_connection_owner_uid":       acctest.CHECKSET,
			"routers.0.physical_connection_status":          acctest.CHECKSET,
			"routers.0.recovery_time":                       "",
			"routers.0.route_table_id":                      acctest.CHECKSET,
			"routers.0.termination_time":                    "",
			"routers.0.virtual_border_router_id":            acctest.CHECKSET,
			"routers.0.virtual_border_router_name":          name,
			"routers.0.vlan_id":                             acctest.CHECKSET,
			"routers.0.vlan_interface_id":                   acctest.CHECKSET,
		}
	}

	var fakeExpressConnectVirtualBorderRoutersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"routers.#": "0",
		}
	}

	var ExpressConnectVirtualBorderRoutersInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existExpressConnectVirtualBorderRoutersMapFunc,
		FakeMapFunc:  fakeExpressConnectVirtualBorderRoutersMapFunc,
	}
	ExpressConnectVirtualBorderRoutersInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, filterConf, allConf)
}

func dataSourceExpressConnectVirtualBorderRoutersDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
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
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}`, name, tfacctest.RandIntRange(1, 2999))
}
