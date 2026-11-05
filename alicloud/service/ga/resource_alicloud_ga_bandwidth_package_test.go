package ga_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	rawClient, err := gaTestClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting AliCloud client.")
	}

	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	for {
		action := "ListBandwidthPackages"
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		logs.AddDebug(action, response, request)

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
			request["ClientToken"] = helper.BuildClientToken("DeleteBandwidthPackage")
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(1*time.Minute, func() *retry.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
				if err != nil {
					if errs.IsExpectedErrors(err, []string{"StateError.BandwidthPackage", "StateError.Accelerator"}) || helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				log.Printf("[ERROR] Deleting bandwidth package %s got an error: %s", bandwidthPackageId, err)
			}
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudGaBandwidthPackage_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBandwidthPackageMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBandwidthPackage")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaBandwidthPackage%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBandwidthPackageBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cbn_geographic_region_ida": "defaultRegionToTest",
					"cbn_geographic_region_idb": "defaultRegionToTest",
					"bandwidth":                 "100",
					"type":                      "Basic",
					"bandwidth_type":            "Basic",
					"auto_pay":                  "true",
					"duration":                  "1",
					"promotion_option_no":       "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":      "100",
						"type":           "Basic",
						"bandwidth_type": "Basic",
						"auto_pay":       "true",
						"duration":       "1",
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
					"renewal_status": "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_duration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_duration": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "BandwidthPackage",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "BandwidthPackage",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaBandwidthPackage_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBandwidthPackageMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBandwidthPackage")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaBandwidthPackage%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBandwidthPackageBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":              "100",
					"type":                   "Basic",
					"bandwidth_type":         "Basic",
					"payment_type":           "Subscription",
					"auto_pay":               "true",
					"duration":               "1",
					"auto_use_coupon":        "false",
					"renewal_status":         "AutoRenewal",
					"auto_renew_duration":    "1",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"bandwidth_package_name": name,
					"description":            name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "BandwidthPackage",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":              "100",
						"type":                   "Basic",
						"bandwidth_type":         "Basic",
						"payment_type":           "Subscription",
						"renewal_status":         "AutoRenewal",
						"auto_renew_duration":    "1",
						"resource_group_id":      acctest.CHECKSET,
						"bandwidth_package_name": name,
						"description":            name,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "BandwidthPackage",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaBandwidthPackage_basic1(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBandwidthPackageMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBandwidthPackage")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaBandwidthPackage%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBandwidthPackageBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":      "100",
					"type":           "Basic",
					"bandwidth_type": "Basic",
					"payment_type":   "PayAsYouGo",
					"billing_type":   "PayBy95",
					"ratio":          "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":      "100",
						"type":           "Basic",
						"bandwidth_type": "Basic",
						"payment_type":   "PayAsYouGo",
						"billing_type":   "PayBy95",
						"ratio":          "30",
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
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "BandwidthPackage",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "BandwidthPackage",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaBandwidthPackage_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_bandwidth_package.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaBandwidthPackageMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBandwidthPackage")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaBandwidthPackage%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBandwidthPackageBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth":              "100",
					"type":                   "Basic",
					"bandwidth_type":         "Basic",
					"payment_type":           "PayAsYouGo",
					"billing_type":           "PayBy95",
					"ratio":                  "30",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"bandwidth_package_name": name,
					"description":            name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "BandwidthPackage",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth":              "100",
						"type":                   "Basic",
						"bandwidth_type":         "Basic",
						"payment_type":           "PayAsYouGo",
						"billing_type":           "PayBy95",
						"ratio":                  "30",
						"resource_group_id":      acctest.CHECKSET,
						"bandwidth_package_name": name,
						"description":            name,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "BandwidthPackage",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "duration", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

var AliCloudGaBandwidthPackageMap0 = map[string]string{
	"resource_group_id": acctest.CHECKSET,
	"status":            acctest.CHECKSET,
}

func AliCloudGaBandwidthPackageBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
	}
`)
}
