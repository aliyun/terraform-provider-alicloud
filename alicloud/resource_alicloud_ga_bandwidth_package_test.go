package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea-rpc/client"
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
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_use_coupon", "auto_pay", "duration"},
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

func TestUnitAlicloudGaBandwidthPackage(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ga_bandwidth_package"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_bandwidth_package"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"auto_pay":                  true,
		"auto_use_coupon":           true,
		"bandwidth":                 10,
		"bandwidth_type":            "CreateBandwidthPackageValue",
		"billing_type":              "CreateBandwidthPackageValue",
		"cbn_geographic_region_ida": "CreateBandwidthPackageValue",
		"cbn_geographic_region_idb": "CreateBandwidthPackageValue",
		"duration":                  "CreateBandwidthPackageValue",
		"payment_type":              "CreateBandwidthPackageValue",
		"ratio":                     10,
		"type":                      "CreateBandwidthPackageValue",
		"auto_renew_duration":       1,
		"renewal_status":            "CreateBandwidthPackageValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeBandwidthPackage
		"Bandwidth":              10,
		"BandwidthPackageId":     "MockBandwidthPackageId",
		"BandwidthType":          "CreateBandwidthPackageValue",
		"BillingType":            "CreateBandwidthPackageValue",
		"CbnGeographicRegionIdA": "CreateBandwidthPackageValue",
		"CbnGeographicRegionIdB": "CreateBandwidthPackageValue",
		"CreateTime":             "DefaultValue",
		"ExpiredTime":            "DefaultValue",
		"ChargeType":             "CreateBandwidthPackageValue",
		"RegionId":               "CreateBandwidthPackageValue",
		"State":                  "active",
		"Ratio":                  10,
		"Type":                   "CreateBandwidthPackageValue",
		"AutoRenewDuration":      1,
		"RenewalStatus":          "CreateBandwidthPackageValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateBandwidthPackage
		"BandwidthPackageId": "MockBandwidthPackageId",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_bandwidth_package", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudGaBandwidthPackageCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeBandwidthPackage Response
		"BandwidthPackageId": "MockBandwidthPackageId",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateBandwidthPackage" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudGaBandwidthPackageCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_bandwidth_package"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudGaBandwidthPackageUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateBandwidthPackagaAutoRenewAttribute
	attributesDiff := map[string]interface{}{
		"auto_renew_duration": 2,
		"renewal_status":      "UpdateBandwidthPackagaAutoRenewAttributeValue",
	}
	diff, err := newInstanceDiff("alicloud_ga_bandwidth_package", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_bandwidth_package"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeBandwidthPackage Response
		"AutoRenewDuration": 2,
		"RenewalStatus":     "UpdateBandwidthPackagaAutoRenewAttributeValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateBandwidthPackagaAutoRenewAttribute" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudGaBandwidthPackageUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_bandwidth_package"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// UpdateBandwidthPackage
	attributesDiff = map[string]interface{}{
		"bandwidth":              15,
		"bandwidth_package_name": "UpdateBandwidthPackageValue",
		"bandwidth_type":         "UpdateBandwidthPackageValue",
		"description":            "UpdateBandwidthPackageValue",
	}
	diff, err = newInstanceDiff("alicloud_ga_bandwidth_package", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_bandwidth_package"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeBandwidthPackage Response
		"Bandwidth":     15,
		"Name":          "UpdateBandwidthPackageValue",
		"BandwidthType": "UpdateBandwidthPackageValue",
		"Description":   "UpdateBandwidthPackageValue",
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "NotExist.BandwidthPackage", "StateError.BandwidthPackage", "UpgradeError.BandwidthPackage", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateBandwidthPackage" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudGaBandwidthPackageUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_bandwidth_package"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeBandwidthPackage" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudGaBandwidthPackageRead(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudGaBandwidthPackageDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteBandwidthPackage" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "DeleteBandwidthPackage" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudGaBandwidthPackageDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
