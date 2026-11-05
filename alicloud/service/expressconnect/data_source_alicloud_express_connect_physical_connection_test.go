package expressconnect_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudExpressConnectPhysicalConnectionsDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	resourceId := "data.alicloud_express_connect_physical_connections.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccExpressConnectPhysicalConnectionsTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectPhysicalConnectionsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_physical_connection.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_physical_connection.default.id}-fake"},
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
			"ids":    []string{"${alicloud_express_connect_physical_connection.default.id}"},
			"status": "Allocated",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_express_connect_physical_connection.default.id}"},
			"status": "Approved",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_express_connect_physical_connection.default.id}"},
			"status":     "Allocated",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alicloud_express_connect_physical_connection.default.id}-fake"},
			"status":     "Approved",
		}),
	}
	var existExpressConnectPhysicalConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                          "1",
			"ids.0":                                          acctest.CHECKSET,
			"names.#":                                        "1",
			"names.0":                                        name,
			"connections.#":                                  "1",
			"connections.0.id":                               acctest.CHECKSET,
			"connections.0.access_point_id":                  acctest.CHECKSET,
			"connections.0.ad_location":                      acctest.CHECKSET,
			"connections.0.bandwidth":                        acctest.CHECKSET,
			"connections.0.business_status":                  acctest.CHECKSET,
			"connections.0.circuit_code":                     "",
			"connections.0.create_time":                      acctest.CHECKSET,
			"connections.0.description":                      acctest.CHECKSET,
			"connections.0.enabled_time":                     "",
			"connections.0.end_time":                         "",
			"connections.0.has_reservation_data":             acctest.CHECKSET,
			"connections.0.line_operator":                    acctest.CHECKSET,
			"connections.0.loa_status":                       "",
			"connections.0.payment_type":                     "",
			"connections.0.peer_location":                    acctest.CHECKSET,
			"connections.0.physical_connection_id":           acctest.CHECKSET,
			"connections.0.physical_connection_name":         acctest.CHECKSET,
			"connections.0.port_number":                      acctest.CHECKSET,
			"connections.0.port_type":                        acctest.CHECKSET,
			"connections.0.redundant_physical_connection_id": "",
			"connections.0.reservation_active_time":          "",
			"connections.0.reservation_internet_charge_type": "",
			"connections.0.reservation_order_type":           "",
			"connections.0.spec":                             acctest.CHECKSET,
			"connections.0.status":                           acctest.CHECKSET,
			"connections.0.type":                             acctest.CHECKSET,
		}
	}

	var fakeExpressConnectPhysicalConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"names.#":       "0",
			"connections.#": "0",
		}
	}

	var ExpressConnectPhysicalConnectionsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existExpressConnectPhysicalConnectionsMapFunc,
		FakeMapFunc:  fakeExpressConnectPhysicalConnectionsMapFunc,
	}

	ExpressConnectPhysicalConnectionsInfo.DataSourceTestCheck(t, 0, idsConf, nameRegexConf, statusConf, allConf)
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
