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
	resource.AddTestSweepers("alicloud_ga_accelerator", &resource.Sweeper{
		Name: "alicloud_ga_accelerator",
		F:    testSweepGaAccelerator,
	})
}

func testSweepGaAccelerator(region string) error {
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
		action := "ListAccelerators"
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.Accelerators", response)
		if err != nil {
			log.Println(err)
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			accelerator := v.(map[string]interface{})
			acceleratorName := fmt.Sprint(accelerator["Name"])
			acceleratorId := fmt.Sprint(accelerator["AcceleratorId"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(acceleratorName, prefix) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ga accelerator: %s(%s) ", acceleratorId, acceleratorName)
				continue
			}
			log.Printf("[Info] Delete Ga accelerator: %s(%s)", acceleratorId, acceleratorName)
			request := make(map[string]interface{})
			request["AcceleratorId"] = acceleratorId
			request["RegionId"] = client.RegionId
			request["PageSize"] = helper.PageSizeLarge
			request["PageNumber"] = 1

			var err error
			for {
				action := "ListIpSets"
				var resp interface{}
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(1*time.Minute, func() *retry.RetryError {
					response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
					if err != nil {
						if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet"}) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					resp, _ = jsonpath.Get("$.IpSets", response)
					return nil
				})

				for _, v := range resp.([]interface{}) {
					request := map[string]interface{}{
						"IpSetId":  v.(map[string]interface{})["IpSetId"],
						"RegionId": client.RegionId,
					}
					action := "DeleteIpSets"
					wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
					err = retry.Retry(1*time.Minute, func() *retry.RetryError {
						request["ClientToken"] = helper.BuildClientToken("DeleteIpSet")
						response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
						if err != nil {
							if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet"}) {
								wait()
								return retry.RetryableError(err)
							}
							return retry.NonRetryableError(err)
						}
						logs.AddDebug(action, resp, request)
						return nil
					})
					if err != nil {
						log.Printf("[ERROR] Deleting ip set %s got an error: %s", request["IpSetId"], err)
					}
				}
				break
			}

			for {
				action := "ListEndpointGroups"
				var resp interface{}
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(1*time.Minute, func() *retry.RetryError {
					response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
					if err != nil {
						if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet"}) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					resp, _ = jsonpath.Get("$.EndpointGroups", response)
					return nil
				})

				for _, v := range resp.([]interface{}) {
					request := map[string]interface{}{
						"EndpointGroupId": v.(map[string]interface{})["EndpointGroupId"],
						"AcceleratorId":   acceleratorId,
					}
					action := "DeleteEndpointGroup"
					wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
					err = retry.Retry(1*time.Minute, func() *retry.RetryError {
						response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
						if err != nil {
							if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup"}) {
								wait()
								return retry.RetryableError(err)
							}
							return retry.NonRetryableError(err)
						}
						logs.AddDebug(action, resp, request)
						return nil
					})
					if err != nil {
						log.Printf("[ERROR] Deleting endpoint group %s got an error: %s", request["EndpointGroupId"], err)
					}
				}
				break
			}

			bandwidthPackageIds := make([]string, 0)
			if v, ok := accelerator["BasicBandwidthPackage"].(map[string]interface{}); ok && len(v) > 0 {
				bandwidthPackageIds = append(bandwidthPackageIds, fmt.Sprint(v["InstanceId"]))
			}
			if v, ok := accelerator["CrossDomainBandwidthPackage"].(map[string]interface{}); ok && len(v) > 0 {
				bandwidthPackageIds = append(bandwidthPackageIds, fmt.Sprint(v["InstanceId"]))
			}
			for _, bandwidthPackageId := range bandwidthPackageIds {
				action := "BandwidthPackageRemoveAccelerator"
				request := map[string]interface{}{
					"AcceleratorId":      acceleratorId,
					"BandwidthPackageId": bandwidthPackageId,
					"RegionId":           client.RegionId,
				}
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
					log.Printf("[ERROR] Removing bandwidth package %s got an error: %s", bandwidthPackageId, err)
				}
			}

			for _, bandwidthPackageId := range bandwidthPackageIds {
				action := "DeleteBandwidthPackage"
				request := map[string]interface{}{
					"BandwidthPackageId": bandwidthPackageId,
					"RegionId":           client.RegionId,
				}
				request["ClientToken"] = helper.BuildClientToken("DeleteBandwidthPackage")
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
				logs.AddDebug(action, response, request)
				if err != nil {
					log.Printf("[ERROR] Deleting bandwidth package %s got an error: %s", bandwidthPackageId, err)
				}
			}

		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudGaAccelerator_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"spec":     "1",
					"duration": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":     "1",
						"duration": "1",
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
					"accelerator_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name,
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
					"spec": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Accelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Accelerator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaAccelerator_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"spec":                   "1",
					"bandwidth_billing_type": "BandwidthPackage",
					"payment_type":           "Subscription",
					"duration":               "1",
					"pricing_cycle":          "Month",
					"auto_use_coupon":        "false",
					"renewal_status":         "AutoRenewal",
					"auto_renew_duration":    "1",
					"promotion_option_no":    "1234",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"accelerator_name":       name,
					"description":            name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Accelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":                   "1",
						"bandwidth_billing_type": "BandwidthPackage",
						"payment_type":           "Subscription",
						"duration":               "1",
						"pricing_cycle":          "Month",
						"auto_use_coupon":        "false",
						"renewal_status":         "AutoRenewal",
						"auto_renew_duration":    "1",
						"promotion_option_no":    "1234",
						"resource_group_id":      acctest.CHECKSET,
						"accelerator_name":       name,
						"description":            name,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Accelerator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaAccelerator_basic1(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_billing_type": "CDT",
					"payment_type":           "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type": "CDT",
						"payment_type":           "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_border_status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_border_status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cross_border_mode": "bgpPro",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cross_border_mode": "bgpPro",
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
					"accelerator_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name,
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
						"For":     "Accelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Accelerator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaAccelerator_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_billing_type": "CDT",
					"payment_type":           "PayAsYouGo",
					"cross_border_status":    "true",
					"cross_border_mode":      "bgpPro",
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					"accelerator_name":       name,
					"description":            name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Accelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type": "CDT",
						"payment_type":           "PayAsYouGo",
						"cross_border_status":    "true",
						"cross_border_mode":      "bgpPro",
						"resource_group_id":      acctest.CHECKSET,
						"accelerator_name":       name,
						"description":            name,
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Accelerator",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

func TestAccAliCloudGaAccelerator_withAnycast(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaAccelerator")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_billing_type": "CDT",
					"payment_type":           "PayAsYouGo",
					"bandwidth":              "200",
					"ip_set_config": []map[string]interface{}{
						{
							"access_mode": "Anycast",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type":      "CDT",
						"payment_type":                "PayAsYouGo",
						"bandwidth":                   "200",
						"ip_set_config.#":             "1",
						"ip_set_config.0.access_mode": "Anycast",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_billing_type": "CDT",
					"payment_type":           "PayAsYouGo",
					"bandwidth":              "250",
					"ip_set_config": []map[string]interface{}{
						{
							"access_mode": "Anycast",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type":      "CDT",
						"payment_type":                "PayAsYouGo",
						"bandwidth":                   "250",
						"ip_set_config.#":             "1",
						"ip_set_config.0.access_mode": "Anycast",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "auto_use_coupon", "promotion_option_no"},
			},
		},
	})
}

var AliCloudGaAcceleratorMap0 = map[string]string{
	"bandwidth_billing_type": acctest.CHECKSET,
	"payment_type":           acctest.CHECKSET,
	"resource_group_id":      acctest.CHECKSET,
	"status":                 acctest.CHECKSET,
}

func AliCloudGaAcceleratorBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
	}
`)
}
