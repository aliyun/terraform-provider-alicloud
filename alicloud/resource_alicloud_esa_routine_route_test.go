package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ESA RoutineRoute. >>> Resource test cases, automatically generated.
// Case 0
func TestAccAliCloudEsaRoutineRoute_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_route.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaRoutineRouteMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutineRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%serr%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaRoutineRouteBasicDependence0)
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
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"routine_name": "${alicloud_esa_routine.default.id}",
					"route_name":   name,
					"route_enable": "on",
					"rule":         "(http.host eq \\\"video.example1.com\\\")",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":      CHECKSET,
						"routine_name": CHECKSET,
						"route_name":   name,
						"route_enable": "on",
						"rule":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "(http.host eq \\\"video.example2.com\\\")",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_enable": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_enable": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bypass": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bypass": "on",
					}),
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
				Config: testAccConfig(map[string]interface{}{
					"sequence": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sequence": "1",
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

func TestAccAliCloudEsaRoutineRoute_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_route.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaRoutineRouteMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutineRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%serr%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaRoutineRouteBasicDependence0)
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
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.id}",
					"routine_name": "${alicloud_esa_routine.default.id}",
					"route_name":   name,
					"route_enable": "on",
					"rule":         "(http.host eq \\\"video.example1.com\\\")",
					"bypass":       "on",
					"fallback":     "on",
					"sequence":     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":      CHECKSET,
						"routine_name": CHECKSET,
						"route_name":   name,
						"route_enable": "on",
						"rule":         CHECKSET,
						"bypass":       "on",
						"fallback":     "on",
						"sequence":     "1",
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

var AliCloudEsaRoutineRouteMap0 = map[string]string{
	"bypass":    CHECKSET,
	"config_id": CHECKSET,
	"fallback":  CHECKSET,
	"sequence":  CHECKSET,
}

func AliCloudEsaRoutineRouteBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_routine" "default" {
  name = var.name
}
`, name)
}

// Test ESA RoutineRoute. <<< Resource test cases, automatically generated.

// ESA enforces a server-side optimistic lock on the same SiteId +
// RoutineName; concurrent RoutineRoute writes surface as LockFailed /
// Site.ServiceBusy / TooManyRequests, which NeedRetry() does not cover.
// Without the IsExpectedErrors(...) wrap on Create / Update / Delete,
// terraform apply with N routes under one routine drops mid-flight; with
// the wrap, the retry+incrementalWait loop absorbs the collision and
// apply/destroy both go green.
func TestAccAliCloudEsaRoutineRoute_lockRetry(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_route.concurrent"
	ra := resourceAttrInit(resourceId, AliCloudEsaRoutineRouteMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaRoutineRoute")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%serrlk%d", defaultRegionToTest, rand)
	const routeCount = 20

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create N routes concurrently under one Routine — stresses Create-side LockFailed.
				Config: testAccAliCloudEsaRoutineRouteLockRetryConfig(name, routeCount, "on"),
				Check: resource.ComposeTestCheckFunc(
					checkAliCloudEsaRoutineRoutesConcurrent(routeCount),
				),
			},
			{
				// Update all N routes concurrently — stresses Update-side LockFailed.
				Config: testAccAliCloudEsaRoutineRouteLockRetryConfig(name, routeCount, "off"),
				Check: resource.ComposeTestCheckFunc(
					checkAliCloudEsaRoutineRoutesConcurrent(routeCount),
				),
			},
		},
		// Framework's post-test destroy runs Delete on all N routes concurrently —
		// implicit stress on Delete-side LockFailed. Any leftover on the routine
		// after test tear-down fails the run.
	})
}

func testAccAliCloudEsaRoutineRouteLockRetryConfig(name string, routeCount int, enable string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_routine" "default" {
  name = var.name
}

resource "alicloud_esa_routine_route" "concurrent" {
  count        = %d
  site_id      = data.alicloud_esa_sites.default.sites.0.id
  routine_name = alicloud_esa_routine.default.id
  route_name   = "${var.name}-${count.index}"
  route_enable = "%s"
  rule         = "(http.host eq \"concurrent${count.index}.example.com\")"
}
`, name, routeCount, enable)
}

func checkAliCloudEsaRoutineRoutesConcurrent(routeCount int) resource.TestCheckFunc {
	checks := make([]resource.TestCheckFunc, 0, routeCount)
	for i := 0; i < routeCount; i++ {
		checks = append(checks, resource.TestCheckResourceAttrSet(
			fmt.Sprintf("alicloud_esa_routine_route.concurrent.%d", i), "config_id",
		))
	}
	return resource.ComposeTestCheckFunc(checks...)
}
