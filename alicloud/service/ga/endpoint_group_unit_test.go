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
func TestUnitAliCloudGaEndpointGroup(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"accelerator_id":        "CreateEndpointGroupValue",
		"description":           "CreateEndpointGroupValue",
		"endpoint_group_region": "CreateEndpointGroupValue",
		"endpoint_configurations": []map[string]interface{}{
			{
				"enable_clientip_preservation": true,
				"endpoint":                     "CreateEndpointGroupValue",
				"type":                         "PublicIp",
				"weight":                       20,
			},
		},
		"endpoint_group_type":           "CreateEndpointGroupValue",
		"listener_id":                   "CreateEndpointGroupValue",
		"endpoint_request_protocol":     "CreateEndpointGroupValue",
		"health_check_interval_seconds": 3,
		"health_check_path":             "CreateEndpointGroupValue",
		"health_check_port":             20,
		"health_check_protocol":         "CreateEndpointGroupValue",
		"name":                          "CreateEndpointGroupValue",
		"port_overrides": []map[string]interface{}{
			{
				"endpoint_port": 10,
				"listener_port": 60,
			},
		},
		"threshold_count":    3,
		"traffic_percentage": 20,
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
		// DescribeEndpointGroup
		"AcceleratorId": "CreateEndpointGroupValue",
		"Description":   "CreateEndpointGroupValue",
		"EndpointConfigurations": []interface{}{
			map[string]interface{}{
				"EnableClientIPPreservation": true,
				"Endpoint":                   "CreateEndpointGroupValue",
				"Type":                       "PublicIp",
				"Weight":                     20,
			},
		},
		"EndpointGroupRegion":        "CreateEndpointGroupValue",
		"EndpointGroupType":          "CreateEndpointGroupValue",
		"HealthCheckIntervalSeconds": 3,
		"HealthCheckPath":            "CreateEndpointGroupValue",
		"HealthCheckPort":            20,
		"HealthCheckProtocol":        "CreateEndpointGroupValue",
		"ListenerId":                 "CreateEndpointGroupValue",
		"EndpointRequestProtocol":    "CreateEndpointGroupValue",
		"Name":                       "CreateEndpointGroupValue",
		"PortOverrides": []interface{}{
			map[string]interface{}{
				"EndpointPort": 10,
				"ListenerPort": 60,
			},
		},
		"State":             "active",
		"ThresholdCount":    3,
		"TrafficPercentage": 20,
	}
	CreateMockResponse := map[string]interface{}{
		// CreateEndpointGroup
		"EndpointGroupId": "CreateEndpointGroupValue",
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ga_endpoint_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGaplusClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeEndpointGroup Response
		"EndpointGroupId": "CreateEndpointGroupValue",
	}
	errorCodes := []string{"NonRetryableError", "GA_NOT_STEADY", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateEndpointGroup" {
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
		err := resourceAliCloudGaEndpointGroupCreate(dInit, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dInit.State(), nil)
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
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateEndpointGroup
	attributesDiff := map[string]interface{}{
		"description": "UpdateEndpointGroup",
		"endpoint_configurations": []map[string]interface{}{
			{
				"enable_clientip_preservation": false,
				"endpoint":                     "UpdateEndpointGroup",
				"type":                         "UpdateEndpointGroup",
				"weight":                       30,
			},
		},
		"endpoint_request_protocol":     "UpdateEndpointGroup",
		"health_check_interval_seconds": 4,
		"health_check_path":             "UpdateEndpointGroup",
		"health_check_port":             30,
		"health_check_protocol":         "UpdateEndpointGroup",
		"name":                          "UpdateEndpointGroup",
		"port_overrides": []map[string]interface{}{
			{
				"endpoint_port": 20,
				"listener_port": 70,
			},
		},
		"threshold_count":    4,
		"traffic_percentage": 30,
	}
	diff, err := newInstanceDiff("alicloud_ga_endpoint_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeEndpointGroup Response
		"Description": "UpdateEndpointGroup",
		"EndpointConfigurations": []interface{}{
			map[string]interface{}{
				"EnableClientIPPreservation": false,
				"Endpoint":                   "UpdateEndpointGroup",
				"Type":                       "UpdateEndpointGroup",
				"Weight":                     30,
			},
		},
		"HealthCheckIntervalSeconds": 4,
		"HealthCheckPath":            "UpdateEndpointGroup",
		"HealthCheckPort":            30,
		"HealthCheckProtocol":        "UpdateEndpointGroup",
		"EndpointRequestProtocol":    "UpdateEndpointGroup",
		"Name":                       "UpdateEndpointGroup",
		"PortOverrides": []interface{}{
			map[string]interface{}{
				"EndpointPort": 20,
				"ListenerPort": 70,
			},
		},
		"State":             "active",
		"ThresholdCount":    4,
		"TrafficPercentage": 30,
	}
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateEndpointGroup" {
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
		err := resourceAliCloudGaEndpointGroupUpdate(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_endpoint_group"].Schema).Data(dExisted.State(), nil)
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
			if *action == "DescribeEndpointGroup" {
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
		err := resourceAliCloudGaEndpointGroupRead(dExisted, rawClient)
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
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceAliCloudGaEndpointGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "StateError.Accelerator", "StateError.EndPointGroup", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteEndpointGroup" {
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
			if *action == "DeleteEndpointGroup" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaEndpointGroupDelete(dExisted, rawClient)
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
