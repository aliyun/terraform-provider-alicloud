package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenTransitRouteTableAggregation_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CenTransitRouteTableAggregationSupportRegions)
	resourceId := "alicloud_cen_transit_route_table_aggregation.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudCenTransitRouteTableAggregationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouteTableAggregation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccCenTransitRouteTableAggregation-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudCenTransitRouteTableAggregationBasicDependence)
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
					"transit_route_table_id":                      "${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}",
					"transit_route_table_aggregation_cidr":        "10.0.0.0/8",
					"transit_route_table_aggregation_scope":       "VPC",
					"transit_route_table_aggregation_name":        name,
					"transit_route_table_aggregation_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_route_table_id":                      CHECKSET,
						"transit_route_table_aggregation_cidr":        "10.0.0.0/8",
						"transit_route_table_aggregation_scope":       "VPC",
						"transit_route_table_aggregation_name":        name,
						"transit_route_table_aggregation_description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var resourceAlicloudCenTransitRouteTableAggregationMap = map[string]string{
	"status": CHECKSET,
}

func resourceAlicloudCenTransitRouteTableAggregationBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id = alicloud_cen_instance.default.id
	}

	resource "alicloud_cen_transit_router_route_table" "default" {
  		transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	}
`, name)
}
