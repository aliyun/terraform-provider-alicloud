package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenTransitRouter_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitrouter%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterBasicDependence0)
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
					"cen_id":              "${alicloud_cen_instance.default.id}",
					"transit_router_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":              CHECKSET,
						"transit_router_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_description": "deds",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_description": "deds",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_description": "desd",
					"transit_router_name":        "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_description": "desd",
						"transit_router_name":        name,
					}),
				),
			},
		},
	})
}

var AlicloudCenTransitRouterMap0 = map[string]string{
	"cen_id":                     CHECKSET,
	"dry_run":                    NOSET,
	"status":                     CHECKSET,
	"transit_router_description": CHECKSET,
	"transit_router_name":        CHECKSET,
}

func AlicloudCenTransitRouterBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
			default = "%s"
	}

	resource "alicloud_cen_instance" "default" {
		cen_instance_name = var.name
	}
	`, name)
}
