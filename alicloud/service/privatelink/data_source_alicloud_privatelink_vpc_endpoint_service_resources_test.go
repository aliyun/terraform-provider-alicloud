package privatelink_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The service's backend SLB is self-created (intranet Classic SLB,
// PayByCLCU/slb.lcu.elastic), so the data source lists a real resource.
func TestAccAlicloudPrivatelinkVpcEndpointServiceResourcesDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_vpc_endpoint_service_resources.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointServiceResources%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivatelinkVpcEndpointServiceResourcesDependence)

	existConfig := testAccConfig(map[string]interface{}{
		"service_id": "${alicloud_privatelink_vpc_endpoint_service_resource.default.service_id}",
	})

	// ids/count/resource_type stay in the check map; resource_id is tied to the
	// real fixture SLB id via TestCheckResourceAttrPair instead of a hardcoded
	// placeholder, so the data source is verified against a real backend.
	existMap := map[string]string{
		"ids.#":                     "1",
		"ids.0":                     acctest.CHECKSET,
		"resources.#":               "1",
		"resources.0.id":            acctest.CHECKSET,
		"resources.0.resource_type": "slb",
	}
	testAccCheck := acctest.ResourceAttrInit(resourceId, existMap).ResourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.SlbPrivateNetSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: existConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
					resource.TestCheckResourceAttrPair(resourceId, "resources.0.resource_id", "alicloud_slb_load_balancer.default", "id"),
				),
			},
		},
	})
}

func dataSourcePrivatelinkVpcEndpointServiceResourcesDependence(name string) string {
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

	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	  service_description    = "test for privatelink connection"
	  service_resource_type  = "slb"
	  auto_accept_connection = false
	}

	resource "alicloud_slb_load_balancer" "default" {
	  load_balancer_name   = "%[1]s"
	  load_balancer_spec   = "slb.lcu.elastic"
	  address_type         = "intranet"
	  instance_charge_type = "PayByCLCU"
	  # eu-west-1 defaults intranet SLB delete_protection to "on"; setting "on"
	  # (the non-default value) makes Create send DeleteProtection="on", so the
	  # server value matches the config and the refresh plan stays empty. This
	  # slb.lcu.elastic intranet SLB pattern is proven in task 22152.
	  delete_protection = "on"
	  vswitch_id       = alicloud_vswitch.default.id
	  master_zone_id   = data.alicloud_slb_zones.default.zones.0.id
	  slave_zone_id    = data.alicloud_slb_zones.default.zones.1.id
	}

	resource "alicloud_privatelink_vpc_endpoint_service_resource" "default" {
	 service_id    = alicloud_privatelink_vpc_endpoint_service.default.id
	 resource_id   = alicloud_slb_load_balancer.default.id
	 resource_type = "slb"
	}
`, name)
}
