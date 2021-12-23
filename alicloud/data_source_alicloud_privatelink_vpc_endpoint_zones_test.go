package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// Currently, Private network slb can only be created through the console.
func SkipTestAccAlicloudPrivatelinkVpcEndpointZonesDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_vpc_endpoint_zones.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointZones%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivatelinkVpcEndpointZonesDependence)

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"endpoint_id": "${alicloud_privatelink_vpc_endpoint_zone.default.endpoint_id}",
			"status":      "Wait",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"endpoint_id": "${alicloud_privatelink_vpc_endpoint_zone.default.endpoint_id}",
			"status":      "Connected",
		}),
	}

	var existPrivatelinkVpcEndpointZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"ids.0":               CHECKSET,
			"zones.#":             "1",
			"zones.0.id":          CHECKSET,
			"zones.0.eni_id":      CHECKSET,
			"zones.0.eni_ip":      CHECKSET,
			"zones.0.status":      "Wait",
			"zones.0.vswitch_id":  CHECKSET,
			"zones.0.zone_domain": CHECKSET,
			"zones.0.zone_id":     CHECKSET,
		}
	}

	var fakePrivatelinkVpcEndpointZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var PrivatelinkVpcEndpointZonesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPrivatelinkVpcEndpointZonesMapFunc,
		fakeMapFunc:  fakePrivatelinkVpcEndpointZonesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
	}

	PrivatelinkVpcEndpointZonesInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, statusConf)
}

func dataSourcePrivatelinkVpcEndpointZonesDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	 name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
	 	 is_default = true
	}

	resource "alicloud_security_group" "default" {
	 name = "%[1]s"
	 description = "privatelink test security group"
	 vpc_id = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	service_description = "test for privatelink connection"
	connect_bandwidth = 103
	auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint_service_resource" "default" {
	 service_id    =  "${alicloud_privatelink_vpc_endpoint_service.default.id}"
	 resource_id   =  "lb-gw8nuxxxxxx"
	 resource_type = "slb"
	}
	resource "alicloud_privatelink_vpc_endpoint" "default" {
	 service_id = alicloud_privatelink_vpc_endpoint_service_resource.default.service_id
	 vpc_id = data.alicloud_vpcs.default.ids.0
	 security_group_id = [alicloud_security_group.default.id]
	 vpc_endpoint_name = "%[1]s"
	 depends_on = [alicloud_privatelink_vpc_endpoint_service.default]
	}
	resource "alicloud_privatelink_vpc_endpoint_zone" "default" {
	 endpoint_id =  alicloud_privatelink_vpc_endpoint.default.id
	 vswitch_id  =  data.alicloud_vswitches.default.ids.0
	 zone_id     =  "eu-central-1a"
	}
   `, name)
}
