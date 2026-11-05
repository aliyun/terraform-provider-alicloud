package esa_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudEsaLoadBalancerRegionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEsaLoadBalancerRegionsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_esa_load_balancer_regions.default", "id"),
					resource.TestCheckResourceAttr("data.alicloud_esa_load_balancer_regions.default", "regions.#", acctest.CHECKSET),
					resource.TestCheckResourceAttrSet("data.alicloud_esa_load_balancer_regions.default", "regions.0.region_code"),
					resource.TestCheckResourceAttrSet("data.alicloud_esa_load_balancer_regions.default", "regions.0.region_en_name"),
				),
			},
		},
	})
}

const testAccEsaLoadBalancerRegionsDataSourceConfig = `
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_esa_load_balancer_regions" "default" {}
`
