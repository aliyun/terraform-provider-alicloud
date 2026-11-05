package nlb

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_nlb_server_group",
		&resource.Sweeper{
			Name: "alicloud_nlb_server_group",
			F:    testSweepNlbServerGroup,
		})
}

func testSweepNlbServerGroup(region string) error {
	rawClient, err := nlbTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListServerGroups"
	request := map[string]interface{}{}
	request["RegionId"] = aliyunClient.RegionId

	request["MaxResults"] = helper.PageSizeLarge

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = aliyunClient.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.ServerGroups", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.ServerGroups", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["ServerGroupName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Nlb Server Group: %s", item["ServerGroupName"].(string))
				continue
			}
			action := "DeleteServerGroup"
			request := map[string]interface{}{
				"ServerGroupId": item["ServerGroupId"],
				"RegionId":      aliyunClient.RegionId,
			}
			_, err = aliyunClient.RpcPost("Nlb", "2022-04-30", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Nlb Server Group (%s): %s", item["ServerGroupName"].(string), err)
			}
			log.Printf("[INFO] Delete Nlb Server Group success: %s ", item["ServerGroupName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// lintignore: R001
func TestUnitAccAliCloudNlbServerGroup(t *testing.T) {
	dInit, _ := schema.InternalMap(ResourceServerGroup().Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(ResourceServerGroup().Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"resource_group_id": "CreateNlbServerGroupValue",
		"server_group_name": "CreateNlbServerGroupValue",
		"server_group_type": "CreateNlbServerGroupValue",
		"vpc_id":            "CreateNlbServerGroupValue",
		"scheduler":         "CreateNlbServerGroupValue",
		"protocol":          "CreateNlbServerGroupValue",
		"health_check": []map[string]interface{}{
			{
				"health_check_enabled":         true,
				"health_check_type":            "CreateNlbServerGroupValue",
				"health_check_connect_port":    0,
				"healthy_threshold":            2,
				"health_check_url":             "CreateNlbServerGroupValue",
				"health_check_domain":          "CreateNlbServerGroupValue",
				"unhealthy_threshold":          2,
				"health_check_connect_timeout": 5,
				"health_check_interval":        10,
				"http_check_method":            "CreateNlbServerGroupValue",
				"health_check_http_code":       []string{"CreateNlbServerGroupValue"},
			},
		},
		"connection_drain":           true,
		"connection_drain_timeout":   60,
		"address_ip_version":         "CreateNlbServerGroupValue",
		"ip_version_affinity_mode":   "CreateNlbServerGroupValue",
		"preserve_client_ip_enabled": true,
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
	rawClient, err := nlbTestClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"RegionId":         "CreateNlbServerGroupValue",
				"ServerGroupId":    "CreateNlbServerGroupValue",
				"ServerGroupName":  "CreateNlbServerGroupValue",
				"ServerGroupType":  "CreateNlbServerGroupValue",
				"AddressIPVersion": "CreateNlbServerGroupValue",
				"VpcId":            "CreateNlbServerGroupValue",
				"Scheduler":        "CreateNlbServerGroupValue",
				"Protocol":         "CreateNlbServerGroupValue",
				"HealthCheck": map[string]interface{}{
					"HealthCheckEnabled":        true,
					"HealthCheckType":           "CreateNlbServerGroupValue",
					"HealthCheckConnectPort":    0,
					"HealthyThreshold":          2,
					"UnhealthyThreshold":        2,
					"HealthCheckConnectTimeout": 5,
					"HealthCheckInterval":       10,
					"HealthCheckDomain":         "CreateNlbServerGroupValue",
					"HealthCheckUrl":            "CreateNlbServerGroupValue",
					"HealthCheckHttpCode":       []string{"CreateNlbServerGroupValue"},
					"HttpCheckMethod":           "CreateNlbServerGroupValue",
					"HealthCheckHttpVersion":    "CreateNlbServerGroupValue",
				},
				"ConnectionDrainEnabled":  true,
				"ConnectionDrainTimeout":  60,
				"PreserveClientIpEnabled": true,
				"IpVersionAffinityMode":   "CreateNlbServerGroupValue",
				"ResourceGroupId":         "CreateNlbServerGroupValue",
				"ServerGroupStatus":       "Available",
				"ServerCount":             2,
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"ServerGroupId": "CreateNlbServerGroupValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_nlb_server_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceServerGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateServerGroup" {
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
		err := resourceServerGroupCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceServerGroup().Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceServerGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"scheduler":                "UpdateNlbServerGroupValue",
		"connection_drain":         false,
		"connection_drain_timeout": 100,
		"health_check": []map[string]interface{}{
			{
				"health_check_enabled":         true,
				"health_check_type":            "UpdateNlbServerGroupValue",
				"health_check_connect_port":    0,
				"healthy_threshold":            2,
				"health_check_url":             "UpdateNlbServerGroupValue",
				"health_check_domain":          "UpdateNlbServerGroupValue",
				"unhealthy_threshold":          2,
				"health_check_connect_timeout": 5,
				"health_check_interval":        10,
				"http_check_method":            "UpdateNlbServerGroupValue",
				"health_check_http_code":       []string{"UpdateNlbServerGroupValue"},
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceServerGroup().Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"Scheduler":              "UpdateNlbServerGroupValue",
				"ConnectionDrainEnabled": false,
				"ConnectionDrainTimeout": 100,
				"HealthCheck": map[string]interface{}{
					"HealthCheckEnabled":        true,
					"HealthCheckType":           "UpdateNlbServerGroupValue",
					"HealthCheckConnectPort":    0,
					"HealthyThreshold":          2,
					"UnhealthyThreshold":        2,
					"HealthCheckConnectTimeout": 5,
					"HealthCheckInterval":       10,
					"HealthCheckDomain":         "UpdateNlbServerGroupValue",
					"HealthCheckUrl":            "UpdateNlbServerGroupValue",
					"HealthCheckHttpCode":       []string{"UpdateNlbServerGroupValue"},
					"HttpCheckMethod":           "UpdateNlbServerGroupValue",
					"HealthCheckHttpVersion":    "UpdateNlbServerGroupValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateServerGroupAttribute" {
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
		err := resourceServerGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceServerGroup().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ServerGroupName
	attributesDiff = map[string]interface{}{
		"server_group_name": "UpdateNlbServerGroupValue",
	}
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceServerGroup().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"ServerGroupName": "UpdateNlbServerGroupValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateServerGroupAttribute" {
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
		err := resourceServerGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceServerGroup().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// TagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"TagResourcesValue_1": "TagResourcesValue_1",
			"TagResourcesValue_2": "TagResourcesValue_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceServerGroup().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"Tags": []interface{}{
					map[string]interface{}{
						"Key":   "TagResourcesValue_1",
						"Value": "TagResourcesValue_1",
					},
					map[string]interface{}{
						"Key":   "TagResourcesValue_2",
						"Value": "TagResourcesValue_2",
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TagResources" {
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
		err := resourceServerGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceServerGroup().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UntagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"UntagResourcesValue3_1": "UnTagResourcesValue3_1",
			"UntagResourcesValue3_2": "UnTagResourcesValue3_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceServerGroup().Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"ServerGroups": []interface{}{
			map[string]interface{}{
				"Tags": []interface{}{
					map[string]interface{}{
						"Key":   "UntagResourcesValue3_1",
						"Value": "UnTagResourcesValue3_1",
					},
					map[string]interface{}{
						"Key":   "UntagResourcesValue3_2",
						"Value": "UnTagResourcesValue3_2",
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UntagResources" {
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
		err := resourceServerGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(ResourceServerGroup().Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceServerGroup().Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListServerGroups" {
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
		err := resourceServerGroupRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewNlbClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceServerGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_nlb_server_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(ResourceServerGroup().Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteServerGroup" {
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
			if *action == "ListServerGroups" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceServerGroupDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test Nlb ServerGroup. >>> Resource test cases, automatically generated.
// Case 3734

func TestUnitNlbServerGroupHealthCheckHttpCodeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{"ReorderOnly", "http_2xx,http_3xx", "http_3xx,http_2xx", true},
		{"ReorderDuplicates", "http_2xx,http_2xx,http_3xx", "http_3xx,http_2xx,http_2xx", true},
		{"MemberAdded", "http_2xx,http_3xx", "http_2xx,http_3xx,http_4xx", false},
		{"MemberRemoved", "http_2xx,http_3xx,http_4xx", "http_2xx,http_3xx", false},
		{"MemberReplaced", "http_2xx", "http_3xx", false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := nlbServerGroupHealthCheckHttpCodeDiffSuppressFunc("health_check.0.health_check_http_code.0", tc.old, tc.new, nil); got != tc.expected {
				t.Fatalf("case %s: suppressor = %v, want %v", tc.name, got, tc.expected)
			}
			diff, err := ResourceServerGroup().Diff(
				context.Background(),
				testNlbServerGroupHealthCheckHttpCodeState(tc.old),
				terraform.NewResourceConfigRaw(testNlbServerGroupHealthCheckHttpCodeConfig(tc.new)),
				nil)
			if err != nil {
				t.Fatalf("case %s: diff: %s", tc.name, err)
			}
			var elementDiffs bool
			if diff != nil {
				for k := range diff.Attributes {
					if strings.HasPrefix(k, "health_check.0.health_check_http_code.") {
						elementDiffs = true
						break
					}
				}
			}
			if elementDiffs == tc.expected {
				t.Fatalf("case %s: instance diff carries health_check_http_code changes = %v, want suppressed = %v", tc.name, elementDiffs, tc.expected)
			}
		})
	}
}

func testNlbServerGroupHealthCheckHttpCodeState(old string) *terraform.InstanceState {
	members := strings.Split(old, ",")
	attributes := map[string]string{
		"health_check.#":                          "1",
		"health_check.0.health_check_enabled":     "true",
		"health_check.0.health_check_type":        "HTTP",
		"health_check.0.health_check_http_code.#": fmt.Sprintf("%d", len(members)),
	}
	for i, member := range members {
		attributes[fmt.Sprintf("health_check.0.health_check_http_code.%d", i)] = member
	}
	return &terraform.InstanceState{ID: "ngp-unit", Attributes: attributes}
}

func testNlbServerGroupHealthCheckHttpCodeConfig(new string) map[string]interface{} {
	members := strings.Split(new, ",")
	values := make([]interface{}, len(members))
	for i, member := range members {
		values[i] = member
	}
	return map[string]interface{}{
		"health_check": []interface{}{
			map[string]interface{}{
				"health_check_enabled":   "true",
				"health_check_type":      "HTTP",
				"health_check_http_code": values,
			},
		},
	}
}

// Test Nlb ServerGroup. <<< Resource test cases, automatically generated.
