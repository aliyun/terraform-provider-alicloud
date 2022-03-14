package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ga_bandwidth_package", &resource.Sweeper{
		Name: "alicloud_ga_bandwidth_package",
		F:    testSweepGaBandwidthPackage,
		Dependencies: []string{
			"alicloud_ga_accelerator",
		},
	})
}

func testSweepGaBandwidthPackage(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}

	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		action := "ListBandwidthPackages"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.BandwidthPackages", response)
		if err != nil {
			log.Println(err)
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			bandwidthPackage := v.(map[string]interface{})
			bandwidthPackageName := fmt.Sprint(bandwidthPackage["Name"])
			bandwidthPackageId := fmt.Sprint(bandwidthPackage["BandwidthPackageId"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(bandwidthPackageName, prefix) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ga bandwidth package: %s(%s) ", bandwidthPackageId, bandwidthPackageName)
				continue
			}
			log.Printf("[Info] Delete Ga bandwidth package: %s(%s)", bandwidthPackageId, bandwidthPackageName)
			action := "DeleteBandwidthPackage"
			request := map[string]interface{}{
				"BandwidthPackageId": bandwidthPackageId,
				"RegionId":           client.RegionId,
			}
			request["ClientToken"] = buildClientToken("DeleteBandwidthPackage")
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(1*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if IsExpectedErrors(err, []string{"StateError.BandwidthPackage", "StateError.Accelerator"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Deleting bandwidth package %s got an error: %s", bandwidthPackageId, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudGaBandwidthPackage_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlicloudGaBandwidthPackageMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaBandwidthPackage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaBandwidthPackageBasicDependence)
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
					"bandwidth":      `100`,
					"type":           "Basic",
					"bandwidth_type": "Basic",
					"billing_type":   "PayBy95",
					"payment_type":   "PayAsYouGo",
					"ratio":          "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":      "100",
						"type":           "Basic",
						"bandwidth_type": "Basic",
						"billing_type":   "PayBy95",
						"payment_type":   "PayAsYouGo",
						"ratio":          "30",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_type", "payment_type", "ratio", "auto_use_coupon"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_type": "Enhanced",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_type": "Enhanced",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "bandwidthpackageDescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "bandwidthpackageDescription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": "${var.name}",
					"description":            "bandwidthpackage",
					"bandwidth":              "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": name,
						"description":            "bandwidthpackage",
						"bandwidth":              "50",
					}),
				),
			},
		},
	})
}

var AlicloudGaBandwidthPackageMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudGaBandwidthPackageBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
