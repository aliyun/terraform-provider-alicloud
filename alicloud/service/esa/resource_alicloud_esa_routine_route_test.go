package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA RoutineRoute. >>> Resource test cases, automatically generated.
// Case 0
func TestAccAliCloudEsaRoutineRoute_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_route.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRoutineRouteMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRoutineRoute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%serr%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRoutineRouteBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"site_id":      acctest.CHECKSET,
						"routine_name": acctest.CHECKSET,
						"route_name":   name,
						"route_enable": "on",
						"rule":         acctest.CHECKSET,
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
						"rule": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRoutineRouteMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRoutineRoute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%serr%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRoutineRouteBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"site_id":      acctest.CHECKSET,
						"routine_name": acctest.CHECKSET,
						"route_name":   name,
						"route_enable": "on",
						"rule":         acctest.CHECKSET,
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
	"bypass":    acctest.CHECKSET,
	"config_id": acctest.CHECKSET,
	"fallback":  acctest.CHECKSET,
	"sequence":  acctest.CHECKSET,
}

func AliCloudEsaRoutineRouteBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_routine" "default" {
  name = var.name
}
`, name)
}

// Test ESA RoutineRoute. <<< Resource test cases, automatically generated.

// ESA enforces a server-side optimistic lock on the same SiteId +
// RoutineName; concurrent RoutineRoute writes surface as LockFailed /
// Site.ServiceBusy / TooManyRequests, which helper.NeedRetry() does not cover.
// Without the errs.IsExpectedErrors(...) wrap on Create / Update / Delete,
// terraform apply with N routes under one routine drops mid-flight; with
// the wrap, the retry+helper.IncrementalWait loop absorbs the collision and
// apply/destroy both go green.
func TestAccAliCloudEsaRoutineRoute_lockRetry(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_routine_route.concurrent"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRoutineRouteMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRoutineRoute")
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%serrlk%d", acctest.DefaultRegionToTest, rand)
	const routeCount = 20

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
