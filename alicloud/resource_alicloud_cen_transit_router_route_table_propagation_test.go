package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

/**
This resource has buried point data.
VBR is buried point data.
*/
func SkipTestAccAlicloudCenTransitRouterRouteTablePropagation_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_route_table_propagation.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterRouteTablePropagationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterRouteTablePropagation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterRouteTablePropagation%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterRouteTablePropagationBasicDependence0)
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
					"transit_router_attachment_id":  "${alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id}",
					"transit_router_route_table_id": "${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_id":  CHECKSET,
						"transit_router_route_table_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudCenTransitRouterRouteTablePropagationMap0 = map[string]string{
	"dry_run":                       NOSET,
	"status":                        CHECKSET,
	"transit_router_attachment_id":  CHECKSET,
	"transit_router_route_table_id": CHECKSET,
}

func AlicloudCenTransitRouterRouteTablePropagationBasicDependence0(name string) string {
	return fmt.Sprintf(`

variable "name" {	
	default = "%s"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  cen_id = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vbr_id = "vbr-j6cd9pm9y6d6e20atoi6w"
  auto_publish_route_enabled = true
  transit_router_attachment_name = var.name
  transit_router_attachment_description = "tf-test"
}
resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}
`, name)
}
