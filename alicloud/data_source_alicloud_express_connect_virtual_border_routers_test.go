package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudExpressConnectVirtualBorderRoutersDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "data.alicloud_express_connect_virtual_border_routers.default"
	rand := acctest.RandIntRange(1, 2999)
	name := fmt.Sprintf("tf-testAccExpressConnectVirtualBorderRoutersTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectVirtualBorderRoutersDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_virtual_border_router.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
			"status": "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_express_connect_virtual_border_router.default.id}"},
			"status": "terminated",
		}),
	}
	filterConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
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
		fakeConfig: testAccConfig(map[string]interface{}{
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
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
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
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alicloud_express_connect_virtual_border_router.default.id}-fake"},
			"status":     "terminated",
		}),
	}
	var existExpressConnectVirtualBorderRoutersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                         "1",
			"ids.0":                                         CHECKSET,
			"names.#":                                       "1",
			"names.0":                                       name,
			"routers.#":                                     "1",
			"routers.0.id":                                  CHECKSET,
			"routers.0.access_point_id":                     CHECKSET,
			"routers.0.activation_time":                     CHECKSET,
			"routers.0.circuit_code":                        "",
			"routers.0.cloud_box_instance_id":               "",
			"routers.0.create_time":                         CHECKSET,
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
			"routers.0.physical_connection_business_status": CHECKSET,
			"routers.0.physical_connection_id":              CHECKSET,
			"routers.0.physical_connection_owner_uid":       CHECKSET,
			"routers.0.physical_connection_status":          CHECKSET,
			"routers.0.recovery_time":                       "",
			"routers.0.route_table_id":                      CHECKSET,
			"routers.0.termination_time":                    "",
			"routers.0.virtual_border_router_id":            CHECKSET,
			"routers.0.virtual_border_router_name":          name,
			"routers.0.vlan_id":                             CHECKSET,
			"routers.0.vlan_interface_id":                   CHECKSET,
		}
	}

	var fakeExpressConnectVirtualBorderRoutersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"routers.#": "0",
		}
	}

	var ExpressConnectVirtualBorderRoutersInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressConnectVirtualBorderRoutersMapFunc,
		fakeMapFunc:  fakeExpressConnectVirtualBorderRoutersMapFunc,
	}
	ExpressConnectVirtualBorderRoutersInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, filterConf, allConf)
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
}`, name, acctest.RandIntRange(1, 2999))
}
