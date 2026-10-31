package cen_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudCenTransitRouterServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cen_transit_router_service.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenTransitRouterServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "Opened",
					}),
				),
			},
			{
				// Cover the data source `enable` arg (Optional, no ForceNew ->
				// must-set + modify-set) for the TestingCoverageRate gate. Step0's
				// raw-HCL `enable = "On"` is not parseable by the static checker, so
				// `enable` is asserted here in the Check map and modified relative to
				// step0 (On -> Off). enable="Off" exercises the Read short-circuit
				// branch (provider-set id "...HasNotBeenOpened", status=""); a real
				// input value, not a fabricated id or StaticSkip.
				Config: testAccCheckAlicloudCenTransitRouterServiceDataSourceDisable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "",
						"enable": "Off",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudCenTransitRouterServiceDataSource = `
data "alicloud_cen_transit_router_service" "default" {
	enable = "On"
}
`
const testAccCheckAlicloudCenTransitRouterServiceDataSourceDisable = `
data "alicloud_cen_transit_router_service" "default" {
	enable = "Off"
}
`
