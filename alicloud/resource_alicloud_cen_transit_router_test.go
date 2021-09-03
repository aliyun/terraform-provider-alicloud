package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_transit_router", &resource.Sweeper{
		Name: "alicloud_cen_transit_router",
		F:    testSweepCenTransitRouters,
	})
}

func testSweepCenTransitRouters(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	for _, cenId := range sweepCenInstanceIds {
		action := "ListTransitRouters"
		request := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["CenId"] = cenId
		request["PageSize"] = PageSizeLarge
		request["PageNumber"] = 1
		var response map[string]interface{}
		conn, err := client.NewCbnClient()
		if err != nil {
			return WrapError(err)
		}
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(2*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] %s failed: %v", action, err)
				return nil
			}
			resp, err := jsonpath.Get("$.TransitRouters", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouters", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				name := fmt.Sprint(item["TransitRouterName"])
				id := fmt.Sprint(item["TransitRouterId"])
				skip := true
				for _, prefix := range prefixes {
					if strings.HasPrefix(name, prefix) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[DEBUG] Skipping the tr %s", name)
					continue
				}

				action := "ListTransitRouterRouteTables"
				request := make(map[string]interface{})
				request["RegionId"] = client.RegionId
				request["TransitRouterId"] = id
				request["PageSize"] = PageSizeLarge
				request["PageNumber"] = 1
				var response map[string]interface{}
				conn, err := client.NewCbnClient()
				if err != nil {
					return WrapError(err)
				}
				for {
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(2*time.Minute, func() *resource.RetryError {
						response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
						if err != nil {
							if NeedRetry(err) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						addDebug(action, response, request)
						return nil
					})
					if err != nil {
						log.Printf("[ERROR] %s failed: %v", action, err)
						return nil
					}
					resp, err := jsonpath.Get("$.RouterTableList.RouterTableListType", response)
					if err != nil {
						return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TransitRouters", response)
					}
					result, _ := resp.([]interface{})
					for _, v := range result {
						item := v.(map[string]interface{})
						name := fmt.Sprint(item["RouteTableName"])
						id := fmt.Sprint(item["RouteTableId"])
						skip := true
						for _, prefix := range prefixes {
							if strings.HasPrefix(name, prefix) {
								skip = false
								break
							}
						}
						if skip {
							log.Printf("[DEBUG] Skipping the tr route table %s", name)
							continue
						}
						action := "DeleteTransitRouterRouteTable"
						log.Printf("[DEBUG] %s %s", action, name)
						request := map[string]interface{}{
							"TransitRouterRouteTableId": id,
						}
						wait := incrementalWait(3*time.Second, 5*time.Second)
						err = resource.Retry(1*time.Minute, func() *resource.RetryError {
							response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
							if err != nil {
								if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
									wait()
									return resource.RetryableError(err)
								}
								return resource.NonRetryableError(err)
							}
							return nil
						})
						addDebug(action, response, request)
						if err != nil {
							log.Printf("[ERROR] %s failed %v", action, err)
						}
					}
					if len(result) < PageSizeLarge {
						break
					}
					request["PageNumber"] = request["PageNumber"].(int) + 1
				}

				action = "DeleteTransitRouter"
				log.Printf("[DEBUG] %s %s", action, name)

				request = map[string]interface{}{
					"TransitRouterId": id,
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					log.Printf("[ERROR] %s failed %v", action, err)
				}
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}
	}
	return nil
}

func TestAccAlicloudCenTransitRouter_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitrouter%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                     "${alicloud_cen_instance.default.id}",
					"transit_router_name":        "${var.name}",
					"transit_router_description": "tf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                     CHECKSET,
						"transit_router_name":        name,
						"transit_router_description": "tf",
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
					"transit_router_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_name": name + "update",
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

var AlicloudCenTransitRouterMap = map[string]string{
	"cen_id":                     CHECKSET,
	"dry_run":                    NOSET,
	"status":                     CHECKSET,
	"transit_router_description": CHECKSET,
	"transit_router_name":        CHECKSET,
}

func AlicloudCenTransitRouterBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_cen_instance" "default" {
		cen_instance_name = var.name
	}
	`, name)
}
