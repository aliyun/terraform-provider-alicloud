package ga

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// White-box (internal package ga) TestUnit cases moved out of the
// external ga_test file so they can call the unexported resource CRUD
// funcs. Shared local helpers (unitTestProviderResources /
// newInstanceDiff / unitTestSharedClientForRegion) live in
// unit_test_helpers_test.go; see that file for why acctest cannot be
// imported here.

// lintignore: R001
func TestUnitAliCloudGaIpSet(t *testing.T) {
	p := unitTestProviderResources()
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
	rawClient, err := unitTestSharedClientForRegion(region)
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
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ga_ip_set", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    tea.String("loadEndpoint error"),
				Data:    tea.String("loadEndpoint error"),
				Message: tea.String("loadEndpoint error"),
			}
		})
		err := resourceAliCloudGaIpSetCreate(dInit, rawClient)
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
			err := resourceAliCloudGaIpSetCreate(dInit, rawClient)
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
				Code:    tea.String("loadEndpoint error"),
				Data:    tea.String("loadEndpoint error"),
				Message: tea.String("loadEndpoint error"),
			}
		})
		err := resourceAliCloudGaIpSetUpdate(dExisted, rawClient)
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
			err := resourceAliCloudGaIpSetUpdate(dExisted, rawClient)
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
			err := resourceAliCloudGaIpSetRead(dExisted, rawClient)
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
				Code:    tea.String("loadEndpoint error"),
				Data:    tea.String("loadEndpoint error"),
				Message: tea.String("loadEndpoint error"),
			}
		})
		err := resourceAliCloudGaIpSetDelete(dExisted, rawClient)
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
			err := resourceAliCloudGaIpSetDelete(dExisted, rawClient)
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
