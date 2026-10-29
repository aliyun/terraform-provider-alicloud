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
func TestUnitAliCloudGaForwardingRule(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_ga_forwarding_rule"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ga_forwarding_rule"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"rule_conditions": []map[string]interface{}{
			{
				"rule_condition_type": "CreateForwardingRulesValue",
				"host_config": []map[string]interface{}{
					{
						"values": []string{"CreateForwardingRulesValue"},
					},
				},
			},
		},
		"priority": 1000,
		"rule_actions": []map[string]interface{}{
			{
				"order":            20,
				"rule_action_type": "CreateForwardingRulesValue",
				"forward_group_config": []map[string]interface{}{
					{
						"server_group_tuples": []map[string]interface{}{
							{
								"endpoint_group_id": "CreateForwardingRulesValue",
							},
						},
					},
				},
			},
		},
		"accelerator_id": "CreateForwardingRulesValue",
		"listener_id":    "CreateForwardingRulesValue",
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
		// ListForwardingRules
		"ForwardingRules": []interface{}{
			map[string]interface{}{
				"Priority":             1000,
				"ForwardingRuleId":     "CreateForwardingRulesValue",
				"ForwardingRuleName":   "CreateForwardingRulesValue",
				"ForwardingRuleStatus": "active",
				"RuleConditions": []interface{}{
					map[string]interface{}{
						"RuleConditionType": "CreateForwardingRulesValue",
						"PathConfig": map[string]interface{}{
							"Values": "CreateForwardingRulesValue",
						},
						"HostConfig": map[string]interface{}{
							"Values": "CreateForwardingRulesValue",
						},
					},
				},
				"RuleActions": []interface{}{
					map[string]interface{}{
						"Order":          20,
						"RuleActionType": "CreateForwardingRulesValue",
						"ForwardGroupConfig": map[string]interface{}{
							"ServerGroupTuples": []interface{}{
								map[string]interface{}{
									"EndpointGroupId": "CreateForwardingRulesValue",
								},
							},
						},
					},
				},
			},
		},
		"ForwardingRuleId": "CreateForwardingRulesValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateForwardingRules
		"ForwardingRules": []interface{}{
			map[string]interface{}{
				"ForwardingRuleId": "CreateForwardingRulesValue",
			},
		},
		"ForwardingRuleId": "CreateForwardingRulesValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ga_forwarding_rule", errorCode))
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
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGaForwardingRuleCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListForwardingRules Response
		"ForwardingRules": []interface{}{
			map[string]interface{}{
				"ForwardingRuleId": "CreateForwardingRulesValue",
			},
		},
		"ForwardingRuleId": "CreateForwardingRulesValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "StateError.Accelerator", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateForwardingRules" {
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
		err := resourceAliCloudGaForwardingRuleCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_forwarding_rule"].Schema).Data(dInit.State(), nil)
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
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGaForwardingRuleUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateForwardingRules
	attributesDiff := map[string]interface{}{
		"forwarding_rule_name": "UpdateForwardingRulesValue",
		"rule_conditions": []map[string]interface{}{
			{
				"rule_condition_type": "UpdateForwardingRulesValue",
				"host_config": []map[string]interface{}{
					{
						"values": []string{"UpdateForwardingRulesValue"},
					},
				},
			},
		},
		"priority": 2000,
		"rule_actions": []map[string]interface{}{
			{
				"order":            30,
				"rule_action_type": "UpdateForwardingRulesValue",
				"forward_group_config": []map[string]interface{}{
					{
						"server_group_tuples": []map[string]interface{}{
							{
								"endpoint_group_id": "UpdateForwardingRulesValue",
							},
						},
					},
				},
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_ga_forwarding_rule", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_forwarding_rule"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListForwardingRules
		"ForwardingRules": []interface{}{
			map[string]interface{}{
				"Priority":             2000,
				"ForwardingRuleName":   "UpdateForwardingRulesValue",
				"ForwardingRuleStatus": "active",
				"RuleConditions": []interface{}{
					map[string]interface{}{
						"RuleConditionType": "UpdateForwardingRulesValue",
						"PathConfig": map[string]interface{}{
							"Values": "UpdateForwardingRulesValue",
						},
						"HostConfig": map[string]interface{}{
							"Values": "UpdateForwardingRulesValue",
						},
					},
				},
				"RuleActions": []interface{}{
					map[string]interface{}{
						"Order":          30,
						"RuleActionType": "UpdateForwardingRulesValue",
						"ForwardGroupConfig": map[string]interface{}{
							"ServerGroupTuples": []interface{}{
								map[string]interface{}{
									"EndpointGroupId": "UpdateForwardingRulesValue",
								},
							},
						},
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateForwardingRules" {
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
		err := resourceAliCloudGaForwardingRuleUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ga_forwarding_rule"].Schema).Data(dExisted.State(), nil)
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
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ga_forwarding_rule", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_forwarding_rule"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ListForwardingRules" {
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
		err := resourceAliCloudGaForwardingRuleRead(dExisted, rawClient)
		patches.Reset()
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
			Code:       tea.String("loadEndpoint error"),
			Data:       tea.String("loadEndpoint error"),
			Message:    tea.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudGaForwardingRuleDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ga_forwarding_rule", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ga_forwarding_rule"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "StateError.Accelerator", "StateError.ForwardingRule", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteForwardingRules" {
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
			return ReadMockResponse, nil
		})
		err := resourceAliCloudGaForwardingRuleDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
