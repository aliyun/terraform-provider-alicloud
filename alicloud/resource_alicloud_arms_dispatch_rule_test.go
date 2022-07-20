package alicloud

import (
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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_arms_dispatch_rule", &resource.Sweeper{
		Name: "alicloud_arms_dispatch_rule",
		F:    testSweepArmsDispatchRule,
	})
}

func testSweepArmsDispatchRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	action := "ListDispatchRule"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		log.Printf("[ERROR] %s failed: %v", action, err)
		return nil
	}
	resp, err := jsonpath.Get("$.DispatchRules", response)
	if err != nil {
		log.Printf("[ERROR] %v", WrapError(err))
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		name := fmt.Sprint(item["Name"])
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping dispatch rule: %s ", name)
			continue
		}
		log.Printf("[INFO] delete dispatch rule: %s ", name)

		action = "DeleteDispatchRule"
		request = map[string]interface{}{
			"Id":       fmt.Sprint(item["RuleId"]),
			"RegionId": client.RegionId,
		}
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s failed: %v", action, err)
		}
	}
	return nil
}

func TestAccAlicloudARMSDispatchRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_dispatch_rule.default"
	ra := resourceAttrInit(resourceId, ArmsDispatchRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsDispatchRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsDispatchRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ArmsDispatchRuleBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dispatch_rule_name": "${var.name}",
					"group_rules": []map[string]interface{}{
						{
							"group_wait_time": "5",
							"group_interval":  "15",
							"grouping_fields": []string{"alertname"},
							"repeat_interval": "61",
						},
					},
					"dispatch_type": "CREATE_ALERT",
					"label_match_expression_grid": []map[string]interface{}{
						{
							"label_match_expression_groups": []map[string]interface{}{
								{
									"label_match_expressions": []map[string]interface{}{
										{
											"key":      "_aliyun_arms_involvedObject_kind",
											"value":    "app",
											"operator": "eq",
										},
									},
								},
							},
						},
					},
					"notify_rules": []map[string]interface{}{
						{
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_id": "${alicloud_arms_alert_contact.default.id}",
									"notify_type":      "ARMS_CONTACT",
									"name":             "${var.name}",
								},
							},
							"notify_channels": []string{"dingTalk", "wechat"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dispatch_rule_name":            CHECKSET,
						"group_rules.#":                 "1",
						"dispatch_type":                 "CREATE_ALERT",
						"label_match_expression_grid.#": "1",
						"notify_rules.#":                "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dispatch_rule_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dispatch_rule_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_rules": []map[string]interface{}{
						{
							"group_wait_time": "10",
							"group_interval":  "25",
							"grouping_fields": []string{"alertname2"},
							"repeat_interval": "70",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_rules.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_rules": []map[string]interface{}{
						{
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_id": "${alicloud_arms_alert_contact.default.id}",
									"notify_type":      "ARMS_CONTACT",
									"name":             "${var.name}",
								},
								{
									"notify_object_id": "${alicloud_arms_alert_contact_group.default.id}",
									"notify_type":      "ARMS_CONTACT_GROUP",
									"name":             "${var.name}",
								},
							},
							"notify_channels": []string{"dingTalk"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_rules.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"label_match_expression_grid": []map[string]interface{}{
						{
							"label_match_expression_groups": []map[string]interface{}{
								{
									"label_match_expressions": []map[string]interface{}{
										{
											"key":      "_aliyun_arms_involvedObject_kind",
											"value":    "app",
											"operator": "eq",
										},
										{
											"key":      "_aliyun_arms_alert_name",
											"value":    "tf-testaccapp",
											"operator": "eq",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"label_match_expression_grid.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dispatch_type": "DISCARD_ALERT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dispatch_type": "DISCARD_ALERT",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dispatch_rule_name": "${var.name}",
					"group_rules": []map[string]interface{}{
						{
							"group_wait_time": "5",
							"group_interval":  "15",
							"grouping_fields": []string{"alertname"},
							"repeat_interval": "80",
						},
					},
					"dispatch_type": "CREATE_ALERT",
					"label_match_expression_grid": []map[string]interface{}{
						{
							"label_match_expression_groups": []map[string]interface{}{
								{
									"label_match_expressions": []map[string]interface{}{
										{
											"key":      "_aliyun_arms_involvedObject_kind",
											"value":    "app",
											"operator": "eq",
										},
									},
								},
							},
						},
					},
					"notify_rules": []map[string]interface{}{
						{
							"notify_objects": []map[string]interface{}{
								{
									"notify_object_id": "${alicloud_arms_alert_contact.default.id}",
									"notify_type":      "ARMS_CONTACT",
									"name":             "${var.name}",
								},
							},
							"notify_channels": []string{"dingTalk", "wechat"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dispatch_rule_name":            CHECKSET,
						"group_rules.#":                 "1",
						"dispatch_type":                 "CREATE_ALERT",
						"label_match_expression_grid.#": "1",
						"notify_rules.#":                "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dispatch_type"},
			},
		},
	})
}

var ArmsDispatchRuleMap = map[string]string{
	"status": CHECKSET,
}

func ArmsDispatchRuleBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_arms_alert_contact" "default" {
  alert_contact_name = "${var.name}"
  email = "${var.name}@aaa.com"
}
resource "alicloud_arms_alert_contact_group" "default" {
  alert_contact_group_name = "${var.name}"
  contact_ids = [alicloud_arms_alert_contact.default.id]
}
`, name)
}

func TestUnitAlicloudARMSDispatchRule(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_arms_dispatch_rule"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_arms_dispatch_rule"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"dispatch_rule_name": "CreateDispatchRuleValue",
		"dispatch_type":      "CreateDispatchRuleValue",
		"group_rules": []interface{}{
			map[string]interface{}{
				"group_interval":  10,
				"group_wait_time": 10,
				"grouping_fields": []interface{}{"CreateDispatchRuleValue0", "CreateDispatchRuleValue1"},
				"repeat_interval": 10,
			},
		},
		"is_recover": true,
		"label_match_expression_grid": []interface{}{
			map[string]interface{}{
				"label_match_expression_groups": []interface{}{
					map[string]interface{}{
						"label_match_expressions": []interface{}{
							map[string]interface{}{
								"key":      "CreateDispatchRuleValue",
								"operator": "CreateDispatchRuleValue",
								"value":    "CreateDispatchRuleValue",
							},
						},
					},
				},
			},
		},
		"notify_rules": []interface{}{
			map[string]interface{}{
				"notify_channels": []interface{}{"CreateDispatchRuleValue0", "CreateDispatchRuleValue1"},
				"notify_objects": []interface{}{
					map[string]interface{}{
						"name":             "CreateDispatchRuleValue",
						"notify_object_id": "10",
						"notify_type":      "CreateDispatchRuleValue",
					},
				},
			},
		},
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
		// DescribeDispatchRule
		"DispatchRule": map[string]interface{}{
			"RuleId":       10,
			"Name":         "CreateDispatchRuleValue",
			"DispatchType": "CreateDispatchRuleValue",
			"GroupRules": []interface{}{
				map[string]interface{}{
					"GroupInterval": 10,
					"GroupWaitTime": 10,
					"GroupingFields": []interface{}{
						"CreateDispatchRuleValue0",
						"CreateDispatchRuleValue1",
					},
					"RepeatInterval": 10,
				},
			},
			"IsRecover": true,
			"LabelMatchExpressionGrid": map[string]interface{}{
				"LabelMatchExpressionGroups": []interface{}{
					map[string]interface{}{
						"LabelMatchExpressions": []interface{}{
							map[string]interface{}{
								"Key":      "CreateDispatchRuleValue",
								"Operator": "CreateDispatchRuleValue",
								"Value":    "CreateDispatchRuleValue",
							},
						},
					},
				},
			},
			"NotifyRules": []interface{}{
				map[string]interface{}{
					"NotifyChannels": []interface{}{
						"CreateDispatchRuleValue0",
						"CreateDispatchRuleValue1",
					},
					"NotifyObjects": []interface{}{
						map[string]interface{}{
							"Name":           "CreateDispatchRuleValue",
							"NotifyObjectId": 10,
							"NotifyType":     "CreateDispatchRuleValue",
						},
					},
				},
			},
			"State": "DefaultValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDispatchRule
		"DispatchRuleId": 10,
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_arms_dispatch_rule", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudArmsDispatchRuleCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeDispatchRule Response
			"DispatchRule": map[string]interface{}{
				"RuleId": 10,
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CreateDispatchRule" {
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
			err := resourceAlicloudArmsDispatchRuleCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_arms_dispatch_rule"].Schema).Data(dInit.State(), nil)
				for key, value := range attributes {
					dCompare.Set(key, value)
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudArmsDispatchRuleUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		// UpdateDispatchRule
		attributesDiff := map[string]interface{}{
			"dispatch_rule_name": "UpdateDispatchRuleValue",
			"dispatch_type":      "UpdateDispatchRuleValue",
			"group_rules": []interface{}{
				map[string]interface{}{
					"group_id":        15,
					"group_interval":  15,
					"group_wait_time": 15,
					"grouping_fields": []interface{}{"UpdateDispatchRuleValue3"},
					"repeat_interval": 15,
				},
			},
			"is_recover": false,
			"label_match_expression_grid": []interface{}{
				map[string]interface{}{
					"label_match_expression_groups": []interface{}{
						map[string]interface{}{
							"label_match_expressions": []interface{}{
								map[string]interface{}{
									"key":      "UpdateDispatchRuleValue",
									"operator": "UpdateDispatchRuleValue",
									"value":    "UpdateDispatchRuleValue",
								},
							},
						},
					},
				},
			},
			"notify_rules": []interface{}{
				map[string]interface{}{
					"notify_channels": []interface{}{"UpdateDispatchRuleValue3"},
					"notify_objects": []interface{}{
						map[string]interface{}{
							"name":             "UpdateDispatchRuleValue",
							"notify_object_id": "15",
							"notify_type":      "UpdateDispatchRuleValue",
						},
					},
				},
			},
		}
		diff, err := newInstanceDiff("alicloud_arms_dispatch_rule", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_arms_dispatch_rule"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// DescribeDispatchRule Response
			"DispatchRule": map[string]interface{}{
				"Name":         "UpdateDispatchRuleValue",
				"DispatchType": "UpdateDispatchRuleValue",
				"GroupRules": []interface{}{
					map[string]interface{}{
						"GroupId":       15,
						"GroupInterval": 15,
						"GroupWaitTime": 15,
						"GroupingFields": []interface{}{
							"UpdateDispatchRuleValue3",
						},
						"RepeatInterval": 15,
					},
				},
				"IsRecover": false,
				"LabelMatchExpressionGrid": map[string]interface{}{
					"LabelMatchExpressionGroups": []interface{}{
						map[string]interface{}{
							"LabelMatchExpressions": []interface{}{
								map[string]interface{}{
									"Key":      "UpdateDispatchRuleValue",
									"Operator": "UpdateDispatchRuleValue",
									"Value":    "UpdateDispatchRuleValue",
								},
							},
						},
					},
				},
				"NotifyRules": []interface{}{
					map[string]interface{}{
						"NotifyChannels": []interface{}{
							"UpdateDispatchRuleValue3",
						},
						"NotifyObjects": []interface{}{
							map[string]interface{}{
								"Name":           "UpdateDispatchRuleValue",
								"NotifyObjectId": 15,
								"NotifyType":     "UpdateDispatchRuleValue",
							},
						},
					},
				},
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "UpdateDispatchRule" {
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
			err := resourceAlicloudArmsDispatchRuleUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_arms_dispatch_rule"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					dCompare.Set(key, value)
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
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "{}"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DescribeDispatchRule" {
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
			err := resourceAlicloudArmsDispatchRuleRead(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "{}":
				assert.Nil(t, err)
			}
		}
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewArmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudArmsDispatchRuleDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteDispatchRule" {
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
			err := resourceAlicloudArmsDispatchRuleDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			}
		}
	})
}
