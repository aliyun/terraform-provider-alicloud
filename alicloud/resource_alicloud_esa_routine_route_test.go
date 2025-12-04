package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA RoutineRoute. >>> Resource test cases, automatically generated.
// Case resource_RoutineRoute_test
func TestAccAliCloudESARoutineRouteresource_RoutineRoute_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_route.default"
	ra := resourceAttrInit(resourceId, AliCloudESARoutineRouteresource_RoutineRoute_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutineRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARoutineRoute%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESARoutineRouteresource_RoutineRoute_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"routine_name": "${alicloud_esa_routine.default.name}",
					"site_id":      "${alicloud_esa_site.default.id}",
					"bypass":       "off",
					"route_name":   "example_routine",
					"route_enable": "on",
					"rule":         "(http.host eq \\\"video.example1.com\\\")",
					"sequence":     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bypass":       "on",
					"route_name":   "test_routine_test2",
					"route_enable": "off",
					"rule":         "(http.host eq \\\"video.example2.com\\\")",
					"sequence":     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fallback": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fallback": "on",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESARoutineRouteresource_RoutineRoute_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESARoutineRouteresource_RoutineRoute_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_routine" "default" {
  description = "example-routine2"
  name        = "example-routine2"
}

`, name)
}

// Test ESA RoutineRoute. <<< Resource test cases, automatically generated.
