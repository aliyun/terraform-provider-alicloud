package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// Currently, Private network slb can only be created through the console.
func SkipTestAccAlicloudPrivatelinkVpcEndpointServiceResourcesDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_vpc_endpoint_service_resources.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointServiceResources%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivatelinkVpcEndpointServiceResourcesDependence)

	serviceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"service_id": "${alicloud_privatelink_vpc_endpoint_service_resource.default.service_id}",
		}),
		fakeConfig: "",
	}

	var existPrivatelinkVpcEndpointServiceResourcesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     CHECKSET,
			"resources.#":               "1",
			"resources.0.id":            CHECKSET,
			"resources.0.resource_id":   "lb-gw8nuym5xxxxx",
			"resources.0.resource_type": "slb",
		}
	}

	var fakePrivatelinkVpcEndpointServiceResourcesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}

	var PrivatelinkVpcEndpointServiceResourcesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPrivatelinkVpcEndpointServiceResourcesMapFunc,
		fakeMapFunc:  fakePrivatelinkVpcEndpointServiceResourcesMapFunc,
	}

	PrivatelinkVpcEndpointServiceResourcesInfo.dataSourceTestCheck(t, 0, serviceIdConf)
}

func dataSourcePrivatelinkVpcEndpointServiceResourcesDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	 name_regex = "default-NODELETING"
	}
	resource "alicloud_security_group" "default" {
	 name = "tf-testAcc-for-privatelink"
	 description = "privatelink test security group"
	 vpc_id = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	service_description = "test for privatelink connection"
	connect_bandwidth = 103
	auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint_service_resource" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service.default.id
	// Currently, Private network slb can only be created through the console.
     resource_id =  "lb-gw8nuym5xxxxx"
     resource_type = "slb"
	}
	`)
}
