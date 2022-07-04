package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGaForwardingRule_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_forwarding_rule.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaForwardingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudGaListener%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaForwardingRuleBasicDependence)
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
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type": "Host",
							"host_config": []map[string]interface{}{
								{
									"values": []string{"www.test.com"},
								},
							},
						},
					},
					"priority": "1000",
					"rule_actions": []map[string]interface{}{
						{
							"order":            "20",
							"rule_action_type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"endpoint_group_id": "${alicloud_ga_endpoint_group.default.id}",
										},
									},
								},
							},
						},
					},
					"accelerator_id": "${alicloud_ga_endpoint_group.default.accelerator_id}",
					"listener_id":    "${alicloud_ga_listener.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_conditions.#": "1",
						"priority":          "1000",
						"rule_actions.#":    "1",
						"accelerator_id":    CHECKSET,
						"listener_id":       CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forwarding_rule_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forwarding_rule_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forwarding_rule_name": name + "update",
					"rule_conditions": []map[string]interface{}{
						{
							"rule_condition_type": "Host",
							"host_config": []map[string]interface{}{
								{
									"values": []string{"www.test3.com"},
								},
							},
						},
					},
					"priority": "2000",
					"rule_actions": []map[string]interface{}{
						{
							"order":            "30",
							"rule_action_type": "ForwardGroup",
							"forward_group_config": []map[string]interface{}{
								{
									"server_group_tuples": []map[string]interface{}{
										{
											"endpoint_group_id": "${alicloud_ga_endpoint_group.default.id}",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forwarding_rule_name": name + "update",
						"rule_conditions.#":    "1",
						"priority":             "2000",
						"rule_actions.#":       "1",
					}),
				),
			},
		},
	})
}

func AlicloudGaForwardingRuleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default  = "%s"
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

resource "alicloud_ga_listener" "default" {
  port_ranges{
    from_port="70"
    to_port="70"
  }
  accelerator_id=alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  client_affinity="SOURCE_IP"
  protocol="HTTP"
  name=var.name
}

resource "alicloud_eip_address" "default" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  address_name = var.name
}

resource "alicloud_ga_endpoint_group" "default" {
  accelerator_id=alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  endpoint_configurations{
    endpoint=alicloud_eip_address.default.ip_address
    type="PublicIp"
    weight="20"
  }
  description=var.name
  name=var.name
  threshold_count=4
  endpoint_group_region="%s"
  health_check_interval_seconds="3"
  health_check_path="/healthcheck"
  health_check_port="9999"
  health_check_protocol="http"
  port_overrides{
    endpoint_port="10"
    listener_port="70"
  }
  traffic_percentage=20
  listener_id=alicloud_ga_listener.default.id
  endpoint_group_type = "virtual"
}
`, name, defaultRegionToTest)
}

func TestUnitAlicloudGaForwardingRule(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
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
	rawClient, err := sharedClientForRegion(region)
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
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ga_forwarding_rule", errorCode))
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGaForwardingRuleCreate(dInit, rawClient)
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
		err := resourceAlicloudGaForwardingRuleCreate(dInit, rawClient)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGaForwardingRuleUpdate(dExisted, rawClient)
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
		err := resourceAlicloudGaForwardingRuleUpdate(dExisted, rawClient)
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
		err := resourceAlicloudGaForwardingRuleRead(dExisted, rawClient)
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
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGaForwardingRuleDelete(dExisted, rawClient)
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
		err := resourceAlicloudGaForwardingRuleDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
