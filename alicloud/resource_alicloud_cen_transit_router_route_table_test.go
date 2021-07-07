package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenTransitRouterRouteTable_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_route_table.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterRouteTableMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterRouteTable")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitrouterroutetable%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterRouteTableBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_id":                      "${alicloud_cen_transit_router.default.transit_router_id}",
					"transit_router_route_table_name":        name,
					"transit_router_route_table_description": "description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_id":                      CHECKSET,
						"transit_router_route_table_name":        name,
						"transit_router_route_table_description": "description",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "transit_router_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_route_table_description": "desp1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_route_table_description": "desp1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_route_table_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_route_table_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_route_table_description": "desp",
					"transit_router_route_table_name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_route_table_description": "desp",
						"transit_router_route_table_name":        name,
					}),
				),
			},
		},
	})
}

var AlicloudCenTransitRouterRouteTableMap = map[string]string{
	"dry_run":                                NOSET,
	"status":                                 CHECKSET,
	"transit_router_id":                      CHECKSET,
	"transit_router_route_table_description": CHECKSET,
	"transit_router_route_table_name":        CHECKSET,
}

func AlicloudCenTransitRouterRouteTableBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}
`, name)
}
