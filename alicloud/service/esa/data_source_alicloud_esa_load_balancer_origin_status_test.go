package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEsaLoadBalancerOriginStatusDataSource_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacc%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEsaLoadBalancerOriginStatusDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_esa_load_balancer_origin_status.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_esa_load_balancer_origin_status.default", "origin_status.#", acctest.CHECKSET),
					resource.TestCheckResourceAttrSet("data.alicloud_esa_load_balancer_origin_status.default", "origin_status.0.load_balancer_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_esa_load_balancer_origin_status.default", "origin_status.0.origin_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_esa_load_balancer_origin_status.default", "origin_status.0.status"),
					resource.TestCheckResourceAttrSet("data.alicloud_esa_load_balancer_origin_status.pool", "origin_status.0.status"),
				),
			},
		},
	})
}

func testAccEsaLoadBalancerOriginStatusDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "alicloud_esa_load_balancer" "default" {
  load_balancer_name = "lb.%s.com"
  fallback_pool      = alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_2.origin_pool_id
  default_pools      = [alicloud_esa_origin_pool.resource_OriginPool_LoadBalancer_test_1_1.origin_pool_id]
  site_id            = alicloud_esa_site.resource_Site_OriginPool_test.id
  steering_policy    = "geo"
  description        = "%s"
  enabled            = true
  monitor {
    type              = "ICMP Ping"
    timeout           = 5
    monitoring_region = "Global"
    consecutive_up    = 3
    consecutive_down  = 5
    interval          = 60
  }
}

data "alicloud_esa_load_balancer_origin_status" "default" {
  site_id           = alicloud_esa_site.resource_Site_OriginPool_test.id
  load_balancer_ids = [alicloud_esa_load_balancer.default.load_balancer_id]
}

data "alicloud_esa_load_balancer_origin_status" "pool" {
  site_id           = alicloud_esa_site.resource_Site_OriginPool_test.id
  load_balancer_ids = [alicloud_esa_load_balancer.default.load_balancer_id]
  pool_type         = "default_pool"
}
`, AliCloudEsaLoadBalancerBasicDependence11682(name), name, name)
}
