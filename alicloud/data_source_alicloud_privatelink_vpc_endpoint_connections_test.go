package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPrivatelinkVpcEndpointConnectionsDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_vpc_endpoint_connections.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointConnections%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivatelinkVpcEndpointConnectionsDependence)

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_id":  "${alicloud_privatelink_vpc_endpoint_connection.default.service_id}",
			"endpoint_id": "${alicloud_privatelink_vpc_endpoint_connection.default.endpoint_id}",
			"status":      "Connected",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_id":  "${alicloud_privatelink_vpc_endpoint_connection.default.service_id}",
			"endpoint_id": "${alicloud_privatelink_vpc_endpoint_connection.default.endpoint_id}",
			"status":      "Pending",
		}),
	}

	var existPrivatelinkVpcEndpointConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     CHECKSET,
			"connections.#":             "1",
			"connections.0.id":          CHECKSET,
			"connections.0.bandwidth":   "1024",
			"connections.0.endpoint_id": CHECKSET,
			"connections.0.status":      CHECKSET,
		}
	}

	var fakePrivatelinkVpcEndpointConnectionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"connections.#": "0",
		}
	}

	var PrivatelinkVpcEndpointConnectionsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPrivatelinkVpcEndpointConnectionsMapFunc,
		fakeMapFunc:  fakePrivatelinkVpcEndpointConnectionsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
	}

	PrivatelinkVpcEndpointConnectionsInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, statusConf)
}

func dataSourcePrivatelinkVpcEndpointConnectionsDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	 name_regex = "default-NODELETING"
	}
	resource "alicloud_security_group" "default" {
	 name = "tf-testacc-forprivatelink"
	 description = "privatelink test security group"
	 vpc_id = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	service_description = "test for privatelink connection"
	connect_bandwidth = 103
	auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service.default.id
	 vpc_id = data.alicloud_vpcs.default.ids.0
	 security_group_ids = [alicloud_security_group.default.id]
	 vpc_endpoint_name = "testformaintf"
	 depends_on = [alicloud_privatelink_vpc_endpoint_service.default]
	}
	resource "alicloud_privatelink_vpc_endpoint_connection" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service.default.id
	 endpoint_id = alicloud_privatelink_vpc_endpoint.default.id
	 bandwidth = "1024"
	}
	`)
}
