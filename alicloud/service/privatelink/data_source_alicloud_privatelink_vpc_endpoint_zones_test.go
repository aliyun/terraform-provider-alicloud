package privatelink_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// The service's backend SLB is self-created (intranet Classic SLB), so the
// data source lists a real endpoint zone. With auto_accept_connection=false the
// added zone stays in the Wait state (the resource Create waits for
// ZoneStatus in {Wait, Connected}), which is what the exist step filters on;
// the Connected filter returns no zones, exercising the empty-result path.
func TestAccAlicloudPrivatelinkVpcEndpointZonesDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_vpc_endpoint_zones.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointZones%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivatelinkVpcEndpointZonesDependence)

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"endpoint_id": "${alicloud_privatelink_vpc_endpoint_zone.default.endpoint_id}",
			"status":      "Wait",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"endpoint_id": "${alicloud_privatelink_vpc_endpoint_zone.default.endpoint_id}",
			"status":      "Connected",
		}),
	}

	var existPrivatelinkVpcEndpointZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"ids.0":               acctest.CHECKSET,
			"zones.#":             "1",
			"zones.0.id":          acctest.CHECKSET,
			"zones.0.eni_id":      acctest.CHECKSET,
			"zones.0.eni_ip":      acctest.CHECKSET,
			"zones.0.status":      "Wait",
			"zones.0.vswitch_id":  acctest.CHECKSET,
			"zones.0.zone_domain": acctest.CHECKSET,
			"zones.0.zone_id":     acctest.CHECKSET,
		}
	}

	var fakePrivatelinkVpcEndpointZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var PrivatelinkVpcEndpointZonesInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existPrivatelinkVpcEndpointZonesMapFunc,
		FakeMapFunc:  fakePrivatelinkVpcEndpointZonesMapFunc,
	}

	// Pinned to the region where the intranet Classic SLB + endpoint zone chain
	// is proven (SlbPrivateNetSupportRegions), reusing the same region condition
	// as the VpcEndpointZone resource tests instead of the hardcoded eu-central-1a.
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.SlbPrivateNetSupportRegions)
	}

	PrivatelinkVpcEndpointZonesInfo.DataSourceTestCheckWithPreCheck(t, 0, preCheck, statusConf)
}

func dataSourcePrivatelinkVpcEndpointZonesDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_slb_zones" "default" {}

	resource "alicloud_vpc" "default" {
	  description = "test-terraform-service"
	  cidr_block  = "10.0.0.0/8"
	  vpc_name    = "%[1]s"
	}

	resource "alicloud_vswitch" "default" {
	  vpc_id     = alicloud_vpc.default.id
	  zone_id    = data.alicloud_slb_zones.default.zones.0.id
	  cidr_block = "10.1.0.0/16"
	}

	resource "alicloud_slb_load_balancer" "default" {
	  load_balancer_name   = "%[1]s"
	  load_balancer_spec   = "slb.s2.small"
	  address_type         = "intranet"
	  instance_charge_type = "PayBySpec"
	  delete_protection    = "off"
	  vswitch_id           = alicloud_vswitch.default.id
	  master_zone_id       = data.alicloud_slb_zones.default.zones.0.id
	  slave_zone_id        = data.alicloud_slb_zones.default.zones.1.id
	  lifecycle {
	    ignore_changes = [delete_protection]
	  }
	}

	resource "alicloud_security_group" "default" {
	 name        = "%[1]s"
	 description = "privatelink test security group"
	 vpc_id      = alicloud_vpc.default.id
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	 service_description    = "test for privatelink connection"
	 service_resource_type  = "slb"
	 auto_accept_connection = false
	}
	resource "alicloud_privatelink_vpc_endpoint_service_resource" "default" {
	 service_id    = "${alicloud_privatelink_vpc_endpoint_service.default.id}"
	 resource_id   = "${alicloud_slb_load_balancer.default.id}"
	 resource_type = "slb"
	}
	resource "alicloud_privatelink_vpc_endpoint" "default" {
	 service_id         = alicloud_privatelink_vpc_endpoint_service_resource.default.service_id
	 vpc_id             = alicloud_vpc.default.id
	 security_group_ids = [alicloud_security_group.default.id]
	 vpc_endpoint_name  = "%[1]s"
	 depends_on         = [alicloud_privatelink_vpc_endpoint_service.default]
	}
	resource "alicloud_privatelink_vpc_endpoint_zone" "default" {
	 endpoint_id = alicloud_privatelink_vpc_endpoint.default.id
	 vswitch_id  = alicloud_vswitch.default.id
	 zone_id     = alicloud_vswitch.default.zone_id
	}
   `, name)
}
