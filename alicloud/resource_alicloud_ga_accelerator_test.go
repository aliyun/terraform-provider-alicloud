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
		action := "ListAccelerators"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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

			conn, err := client.NewGaplusClient()
			if err != nil {
				return WrapError(err)
			}
			for {
				action := "ListIpSets"
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				var resp interface{}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					action := "DeleteIpSets"
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(1*time.Minute, func() *resource.RetryError {
						request["ClientToken"] = buildClientToken("DeleteIpSet")
						resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				var resp interface{}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)
					action := "DeleteEndpointGroup"
					wait := incrementalWait(3*time.Second, 3*time.Second)
					err = resource.Retry(1*time.Minute, func() *resource.RetryError {
						resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

func TestAccAlicloudGaAccelerator_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence)
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
					"spec":            "1",
					"auto_use_coupon": "true",
					"duration":        "1",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Accelerator",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec":            "1",
						"auto_use_coupon": "true",
						"duration":        "1",
						"tags.%":          "2",
						"tags.Created":    "TF",
						"tags.For":        "Accelerator",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "accelerator_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "accelerator_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": `2`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name,
					"description":      "accelerator",
					"spec":             `1`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name,
						"description":      "accelerator",
						"spec":             "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF_Update",
						"For":     "Accelerator_Update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF_Update",
						"tags.For":     "Accelerator_Update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "auto_use_coupon"},
			},
		},
	})
}

func TestAccAlicloudGaAccelerator_basic01(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaAccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence)
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
					"accelerator_name":       name,
					"description":            "accelerator",
					"spec":                   "1",
					"auto_use_coupon":        "true",
					"duration":               "1",
					"bandwidth_billing_type": "CDT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name":       name,
						"description":            "accelerator",
						"spec":                   "1",
						"auto_use_coupon":        "true",
						"duration":               "1",
						"bandwidth_billing_type": "CDT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "accelerator_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "accelerator_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": `2`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name,
					"description":      "accelerator",
					"spec":             `1`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name,
						"description":      "accelerator",
						"spec":             "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "auto_use_coupon"},
			},
		},
	})
}

var AlicloudGaAcceleratorMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudGaAcceleratorBasicDependence(name string) string {
	return ""
}

func TestUnitAlicloudGaAccelerator(t *testing.T) {
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
		err := resourceAlicloudCrChartRepositoryRead(d, rawClient)
		patchDoRequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Ga Accelerator. >>> Resource test cases, automatically generated.
// Case 3688
func TestAccAlicloudGaAccelerator_basic3688(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap3688)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence3688)
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
					"duration":      "1",
					"pricing_cycle": "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"duration":      "1",
						"pricing_cycle": "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type":   "POSTPAY",
					"duration":               "1",
					"auto_pay":               "true",
					"bandwidth_billing_type": "CDT",
					"pricing_cycle":          "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":   "POSTPAY",
						"duration":               "1",
						"auto_pay":               "true",
						"bandwidth_billing_type": "CDT",
						"pricing_cycle":          "Month",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudGaAcceleratorMap3688 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudGaAcceleratorBasicDependence3688(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3179
func TestAccAlicloudGaAccelerator_basic3179(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap3179)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence3179)
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
					"spec": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_billing_type": "CDT",
					"accelerator_name":       name + "_update",
					"auto_pay":               "true",
					"spec":                   "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type": "CDT",
						"accelerator_name":       name + "_update",
						"auto_pay":               "true",
						"spec":                   "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudGaAcceleratorMap3179 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudGaAcceleratorBasicDependence3179(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 1858
func TestAccAlicloudGaAccelerator_basic1858(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap1858)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence1858)
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
					"pricing_cycle":    "Month",
					"duration":         "2",
					"accelerator_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pricing_cycle":    "Month",
						"duration":         "2",
						"accelerator_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name + "_update",
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
					"accelerator_name": name + "_update",
					"auto_use_coupon":  "false",
					"pricing_cycle":    "Month",
					"duration":         "2",
					"spec":             "1",
					"auto_pay":         "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name": name + "_update",
						"auto_use_coupon":  "false",
						"pricing_cycle":    "Month",
						"duration":         "2",
						"spec":             "1",
						"auto_pay":         "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudGaAcceleratorMap1858 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudGaAcceleratorBasicDependence1858(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 678
func TestAccAlicloudGaAccelerator_basic678(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap678)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence678)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudGaAcceleratorMap678 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudGaAcceleratorBasicDependence678(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 225
func TestAccAlicloudGaAccelerator_basic225(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap225)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence225)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

var AlicloudGaAcceleratorMap225 = map[string]string{
	"payment_type": CHECKSET,
	"status":       CHECKSET,
	"create_time":  CHECKSET,
}

func AlicloudGaAcceleratorBasicDependence225(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3688  twin
func TestAccAlicloudGaAccelerator_basic3688_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap3688)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence3688)
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
					"instance_charge_type":   "POSTPAY",
					"duration":               "1",
					"auto_pay":               "true",
					"bandwidth_billing_type": "CDT",
					"pricing_cycle":          "Month",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":   "POSTPAY",
						"duration":               "1",
						"auto_pay":               "true",
						"bandwidth_billing_type": "CDT",
						"pricing_cycle":          "Month",
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

// Case 3179  twin
func TestAccAlicloudGaAccelerator_basic3179_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap3179)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence3179)
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
					"accelerator_name":       name,
					"auto_pay":               "true",
					"spec":                   "5",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_billing_type": "CDT",
						"accelerator_name":       name,
						"auto_pay":               "true",
						"spec":                   "5",
						"tags.%":                 "2",
						"tags.Created":           "TF",
						"tags.For":               "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

// Case 1858  twin
func TestAccAlicloudGaAccelerator_basic1858_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap1858)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence1858)
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
					"accelerator_name":    name,
					"auto_use_coupon":     "false",
					"pricing_cycle":       "Month",
					"duration":            "2",
					"spec":                "2",
					"auto_pay":            "false",
					"description":         "Descriptive information of the global acceleration instance1",
					"auto_renew":          "true",
					"renewal_status":      "NotRenewal",
					"auto_renew_duration": "2",
					"resource_group_id":   "rg-aek2xl5qajpkquq",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_name":    name,
						"auto_use_coupon":     "false",
						"pricing_cycle":       "Month",
						"duration":            "2",
						"spec":                "2",
						"auto_pay":            "false",
						"description":         "Descriptive information of the global acceleration instance1",
						"auto_renew":          "true",
						"renewal_status":      "NotRenewal",
						"auto_renew_duration": "2",
						"resource_group_id":   "rg-aek2xl5qajpkquq",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

// Case 678  twin
func TestAccAlicloudGaAccelerator_basic678_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap678)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence678)
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

// Case 225  twin
func TestAccAlicloudGaAccelerator_basic225_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_accelerator.default"
	ra := resourceAttrInit(resourceId, AlicloudGaAcceleratorMap225)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaAccelerator")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaaccelerator%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaAcceleratorBasicDependence225)
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "auto_renew_duration", "auto_use_coupon", "ddos_region_id", "duration", "instance_charge_type", "pricing_cycle", "promotion_option_no"},
			},
		},
	})
}

// Test Ga Accelerator. <<< Resource test cases, automatically generated.
