package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ga_accelerator", &resource.Sweeper{
		Name: "alicloud_ga_accelerator",
		F:    testSweepGaAccelerator,
	})
}

func testSweepGaAccelerator(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
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

	for {
		action := "ListAccelerators"
		response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %v", action, err)
			break
		}
		addDebug(action, response, request)

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
			request["PageSize"] = PageSizeLarge
			request["PageNumber"] = 1

			var err error
			for {
				action := "ListIpSets"
				var resp interface{}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
					if err != nil {
						if NeedRetry(err) || IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet"}) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
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
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(1*time.Minute, func() *resource.RetryError {
						request["ClientToken"] = buildClientToken("DeleteIpSet")
						response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
						if err != nil {
							if NeedRetry(err) || IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet"}) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						addDebug(action, resp, request)
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
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err := client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
					if err != nil {
						if NeedRetry(err) || IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.IpSet"}) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
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
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(1*time.Minute, func() *resource.RetryError {
						response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
						if err != nil {
							if NeedRetry(err) || IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup"}) {
								wait()
								return resource.RetryableError(err)
							}
							return resource.NonRetryableError(err)
						}
						addDebug(action, resp, request)
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
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
					log.Printf("[ERROR] Removing bandwidth package %s got an error: %s", bandwidthPackageId, err)
				}
			}

			for _, bandwidthPackageId := range bandwidthPackageIds {
				action := "DeleteBandwidthPackage"
				request := map[string]interface{}{
					"BandwidthPackageId": bandwidthPackageId,
					"RegionId":           client.RegionId,
				}
				request["ClientToken"] = buildClientToken("DeleteBandwidthPackage")
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
				addDebug(action, response, request)
				if err != nil {
					log.Printf("[ERROR] Deleting bandwidth package %s got an error: %s", bandwidthPackageId, err)
				}
			}

		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudGaAccelerator_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
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
						"resource_group_id": CHECKSET,
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
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
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
						"resource_group_id":      CHECKSET,
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
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
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
						"resource_group_id": CHECKSET,
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
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AliCloudGaAcceleratorMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaAcceleratorBasicDependence0)
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
						"resource_group_id":      CHECKSET,
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

var AliCloudGaAcceleratorMap0 = map[string]string{
	"bandwidth_billing_type": CHECKSET,
	"payment_type":           CHECKSET,
	"resource_group_id":      CHECKSET,
	"status":                 CHECKSET,
}

func AliCloudGaAcceleratorBasicDependence0(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default" {
	}
`)
}

func TestUnitAliCloudGaAccelerator(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"auto_use_coupon":     true,
		"duration":            9,
		"spec":                "1",
		"accelerator_name":    "accelerator_name",
		"auto_renew_duration": 1,
		"description":         "description",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		//DescribeGaAccelerator
		"Name":        "accelerator_name",
		"State":       "active",
		"Description": "description",
		"Spec":        "1",

		//DescribeAcceleratorAutoRenewAttribute
		"AutoRenewDuration": 1,
		"RenewalStatus":     false,
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_accelerator", "MockId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			result["AcceleratorId"] = "MockId"
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudGaAcceleratorCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		d.Set("duration", 24)
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudGaAcceleratorCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockId")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudGaAcceleratorUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateAcceleratorAutoRenewAttributeAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"auto_renew_duration", "renewal_status"} {
			switch p["alicloud_ga_accelerator"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudGaAcceleratorUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateAcceleratorAutoRenewAttributeNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"auto_renew_duration", "renewal_status"} {
			switch p["alicloud_ga_accelerator"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("UpdateAcceleratorAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"accelerator_name", "description", "spec", "auto_use_coupon"} {
			switch p["alicloud_ga_accelerator"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["Normal"]("")
		})
		err := resourceAliCloudGaAcceleratorUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateAcceleratorNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"accelerator_name", "description", "spec", "auto_use_coupon"} {
			switch p["alicloud_ga_accelerator"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ga_accelerator"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("RetryError")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeGaAcceleratorNotFound", func(t *testing.T) {
		patchDoRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorRead(d, rawClient)
		patchDoRequest.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadDescribeCrChartNamespaceUnexpectedError", func(t *testing.T) {
		patchDoRequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := false
			noRetryFlag := true
			if NotFoundFlag {
				return responseMock["NotFoundError"]("CHART_REPO_NOT_EXIST")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudGaAcceleratorRead(d, rawClient)
		patchDoRequest.Reset()
		assert.NotNil(t, err)
	})
}
