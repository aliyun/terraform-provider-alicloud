package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudExpressConnectPhysicalConnectionsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "data.alicloud_express_connect_physical_connections.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccExpressConnectPhysicalConnectionsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectPhysicalConnectionsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_physical_connection.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_physical_connection.default.id}-fake"},
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
			"ids":    []string{"${alicloud_express_connect_physical_connection.default.id}"},
			"status": "Allocated",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_express_connect_physical_connection.default.id}"},
			"status": "Approved",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_express_connect_physical_connection.default.id}"},
			"status":     "Allocated",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alicloud_express_connect_physical_connection.default.id}-fake"},
			"status":     "Approved",
		}),
	}
	var existExpressConnectPhysicalConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                          "1",
			"ids.0":                                          CHECKSET,
			"names.#":                                        "1",
			"names.0":                                        name,
			"connections.#":                                  "1",
			"connections.0.id":                               CHECKSET,
			"connections.0.access_point_id":                  CHECKSET,
			"connections.0.ad_location":                      CHECKSET,
			"connections.0.bandwidth":                        CHECKSET,
			"connections.0.business_status":                  CHECKSET,
			"connections.0.circuit_code":                     "",
			"connections.0.create_time":                      CHECKSET,
			"connections.0.description":                      CHECKSET,
			"connections.0.enabled_time":                     "",
			"connections.0.end_time":                         "",
			"connections.0.has_reservation_data":             CHECKSET,
			"connections.0.line_operator":                    CHECKSET,
			"connections.0.loa_status":                       "",
			"connections.0.payment_type":                     "",
			"connections.0.peer_location":                    CHECKSET,
			"connections.0.physical_connection_id":           CHECKSET,
			"connections.0.physical_connection_name":         CHECKSET,
			"connections.0.port_number":                      CHECKSET,
			"connections.0.port_type":                        CHECKSET,
			"connections.0.redundant_physical_connection_id": "",
			"connections.0.reservation_active_time":          "",
			"connections.0.reservation_internet_charge_type": "",
			"connections.0.reservation_order_type":           "",
			"connections.0.spec":                             CHECKSET,
			"connections.0.status":                           CHECKSET,
			"connections.0.type":                             CHECKSET,
		}
	}

	var fakeExpressConnectPhysicalConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"names.#":       "0",
			"connections.#": "0",
		}
	}

	var ExpressConnectPhysicalConnectionsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressConnectPhysicalConnectionsMapFunc,
		fakeMapFunc:  fakeExpressConnectPhysicalConnectionsMapFunc,
	}

	ExpressConnectPhysicalConnectionsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, statusConf, allConf)
}

func dataSourceExpressConnectPhysicalConnectionsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_express_connect_physical_connection" "default" {
  access_point_id          = "%s"
  line_operator            = "CT"
  peer_location            = var.name
  physical_connection_name = var.name
  type                     = "VPC"
  description              = "my domestic connection"
  port_type                = "1000Base-LX"
  bandwidth                = 100
}`, name, getAccessPointId())
}
