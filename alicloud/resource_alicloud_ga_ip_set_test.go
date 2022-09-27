package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"log"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaIpSet_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_ip_set.default"
	ra := resourceAttrInit(resourceId, AlicloudGaIpSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaIpSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudGaIpSetBasicDependence)
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
					"accelerate_region_id": defaultRegionToTest,
					"bandwidth":            "5",
					"accelerator_id":       "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerate_region_id": defaultRegionToTest,
						"bandwidth":            "5",
						"accelerator_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accelerator_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": `10`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
		},
	})
}

var AlicloudGaIpSetMap = map[string]string{
	"ip_address_list.#": CHECKSET,
	"ip_version":        "IPv4",
	"status":            "active",
}

func AlicloudGaIpSetBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_bandwidth_package" "default" {
   	bandwidth              =  100
  	type                   = "Basic"
  	bandwidth_type         = "Basic"
	payment_type           = "PayAsYouGo"
  	billing_type           = "PayBy95"
	ratio       = 30
	bandwidth_package_name = var.name
    auto_pay               = true
    auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
	// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
	accelerator_id = data.alicloud_ga_accelerators.default.ids.0
	bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}
`, name)
}

func TestUnitAlicloudGaIpSet(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ga_ip_set"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_ip_set"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id":       "CreateIpSetsValue",
		"accelerate_region_id": "CreateIpSetsValue",
		"ip_version":           "CreateIpSetsValue",
		"bandwidth":            10,
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
		// DescribeIpSet
		"AccelerateRegionId": "CreateIpSetsValue",
		"IpAddressList":      "CreateIpSetsValue",
		"IpVersion":          "CreateIpSetsValue",
		"Bandwidth":          10,
		"State":              "active",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateIpSets
		"IpSets": []interface{}{
			map[string]interface{}{
				"IpSetId": "MockIpSetId",
			},
		},
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_ip_set", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudGaIpSetCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{}
		errorCodes := []string{"NonRetryableError", "Throttling", "StateError.Accelerator", "StateError.IpSet", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CreateIpSets" {
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
			err := resourceAlicloudGaIpSetCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_ga_ip_set"].Schema).Data(dInit.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudGaIpSetUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		// UpdateIpSet
		attributesDiff := map[string]interface{}{
			"bandwidth": 15,
		}
		diff, err := newInstanceDiff("alicloud_ga_ip_set", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_ga_ip_set"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeIpSet Response
			"Bandwidth": 15,
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "StateError.Accelerator", "StateError.IpSet", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "UpdateIpSet" {
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
			err := resourceAlicloudGaIpSetUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_ga_ip_set"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})
	// Read
	t.Run("Read", func(t *testing.T) {
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "UnknownError", "{}"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DescribeIpSet" {
					switch errorCode {
					case "{}", "UnknownError":
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
			err := resourceAlicloudGaIpSetRead(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "{}", "UnknownError":
				assert.Nil(t, err)
			}
		}
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudGaIpSetDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		errorCodes := []string{"NonRetryableError", "Throttling", "StateError.Accelerator", "StateError.IpSet", "nil", "NotExist.IpSets"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteIpSets" {
					switch errorCode {
					case "NonRetryableError", "NotExist.IpSets":
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
				return ReadMockResponse, nil
			})
			err := resourceAlicloudGaIpSetDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "NotExist.IpSets":
				assert.Nil(t, err)
			}
		}
	})
}
