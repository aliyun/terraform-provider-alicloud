package alicloud

import (
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_dms_enterprise_user", &resource.Sweeper{
		Name: "alicloud_dms_enterprise_user",
		F:    testSweepDMSEnterpriseUsers,
	})
}

func testSweepDMSEnterpriseUsers(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"testacc",
	}
	request := map[string]interface{}{
		"UserState":  "NORMAL",
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	var response map[string]interface{}
	action := "ListUsers"
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_users", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.UserList.User", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.UserList.User", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if _, ok := item["NickName"]; !ok {
				skip = false
			} else {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprintf("%v", item["NickName"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DMS Enterprise User: %v", item["NickName"])
				continue
			}
			action := "DeleteUser"
			request := map[string]interface{}{
				"Uid": item["Uid"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DMS Enterprise User (%v): %s", item["NickName"], err)
				continue
			}

			log.Printf("[INFO] Delete DMS Enterprise User Success: %v ", item["NickName"])
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDMSEnterpriseUser_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_user.default"
	ra := resourceAttrInit(resourceId, DmsEnterpriseUserMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Dms_enterpriseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDmsEnterpriseUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccDmsEnterpriseUser%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DmsEnterpriseUserBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"uid":        "${alicloud_ram_user.user.id}",
					"nick_name":  name,
					"role_names": []string{"DBA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name":    name,
						"role_names.#": "1",
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
					"max_execute_count": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_execute_count": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_result_count": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_result_count": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nick_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_names": []string{"USER"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "NORMAL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_execute_count": "1000",
					"max_result_count":  "1000",
					"nick_name":         name + "change",
					"role_names":        []string{"DBA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_execute_count": "1000",
						"max_result_count":  "1000",
						"nick_name":         name + "change",
						"role_names.#":      "1",
					}),
				),
			},
		},
	})
}

var DmsEnterpriseUserMap = map[string]string{
	"status": CHECKSET,
}

func DmsEnterpriseUserBasicdependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name         = "%s"
	  display_name = "user_display_name"
	  mobile       = "86-18688888888"
	  email        = "hello.uuu@aaa.com"
	  comments     = "yoyoyo"
	}`, name)
}

func TestAccAlicloudDMSEnterpriseUser_unit(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"config_rule_trigger_types":  "CreateConfigRuleValue",
		"description":                "CreateConfigRuleValue",
		"exclude_resource_ids_scope": "CreateConfigRuleValue",
		"input_parameters": map[string]interface{}{
			"vpcIds": "CreateConfigRuleValue",
		},
		"maximum_execution_frequency": "CreateConfigRuleValue",
		"region_ids_scope":            "CreateConfigRuleValue",
		"resource_group_ids_scope":    "CreateConfigRuleValue",
		"resource_types_scope":        []interface{}{"CreateConfigRuleValue0", "CreateConfigRuleValue1"},
		"risk_level":                  10,
		"rule_name":                   "CreateConfigRuleValue",
		"source_identifier":           "CreateConfigRuleValue",
		"source_owner":                "CreateConfigRuleValue",
		"tag_key_scope":               "CreateConfigRuleValue",
		"tag_value_scope":             "CreateConfigRuleValue",
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
		// GetConfigRule
		"ConfigRule": map[string]interface{}{
			"ConfigRuleId":            "CreateConfigRuleValue",
			"Description":             "CreateConfigRuleValue",
			"ExcludeResourceIdsScope": "CreateConfigRuleValue",
			"InputParameters": map[string]interface{}{
				"vpcIds": "CreateConfigRuleValue",
			},
			"MaximumExecutionFrequency": "CreateConfigRuleValue",
			"RegionIdsScope":            "CreateConfigRuleValue",
			"ResourceGroupIdsScope":     "CreateConfigRuleValue",
			"Scope": map[string]interface{}{
				"ComplianceResourceTypes": "CreateConfigRuleValue",
			},
			"Source": map[string]interface{}{
				"Owner":      "CreateConfigRuleValue",
				"Identifier": "CreateConfigRuleValue",
				"SourceDetails": []interface{}{
					map[string]interface{}{
						"MessageType": "CreateConfigRuleValue",
					},
				},
			},
			"RiskLevel":       10,
			"ConfigRuleName":  "CreateConfigRuleValue",
			"ConfigRuleState": "ACTIVE",
			"TagKeyScope":     "CreateConfigRuleValue",
			"TagValueScope":   "CreateConfigRuleValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateConfigRule
		"ConfigRuleId": "CreateConfigRuleValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_config_rule", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigRuleCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleId": "CreateConfigRuleValue",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateConfigRule" {
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
		err := resourceAlicloudConfigRuleCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigRuleUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateConfigRule
	attributesDiff := map[string]interface{}{
		"config_rule_trigger_types":  "UpdateConfigRuleValue",
		"description":                "UpdateConfigRuleValue",
		"resource_types_scope":       []interface{}{"UpdateConfigRuleValue"},
		"risk_level":                 15,
		"exclude_resource_ids_scope": "UpdateConfigRuleValue",
		"input_parameters": map[string]interface{}{
			"vpcIds": "UpdateConfigRuleValue",
		},
		"maximum_execution_frequency": "UpdateConfigRuleValue",
		"region_ids_scope":            "UpdateConfigRuleValue",
		"resource_group_ids_scope":    "UpdateConfigRuleValue",
		"tag_key_scope":               "UpdateConfigRuleValue",
		"tag_value_scope":             "UpdateConfigRuleValue",
	}
	diff, err := newInstanceDiff("alicloud_config_rule", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetConfigRule Response
		"ConfigRule": map[string]interface{}{
			"Source": map[string]interface{}{
				"SourceDetails": []interface{}{
					map[string]interface{}{
						"MessageType": "UpdateConfigRuleValue",
					},
				},
			},
			"Description":             "UpdateConfigRuleValue",
			"ResourceTypesScope":      []interface{}{"UpdateConfigRuleValue"},
			"RiskLevel":               15,
			"ExcludeResourceIdsScope": "UpdateConfigRuleValue",
			"InputParameters": map[string]interface{}{
				"vpcIds": "UpdateConfigRuleValue",
			},
			"MaximumExecutionFrequency": "UpdateConfigRuleValue",
			"RegionIdsScope":            "UpdateConfigRuleValue",
			"ResourceGroupIdsScope":     "UpdateConfigRuleValue",
			"TagKeyScope":               "UpdateConfigRuleValue",
			"TagValueScope":             "UpdateConfigRuleValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateConfigRule" {
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
		err := resourceAlicloudConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// StopConfigRules
	attributesDiff = map[string]interface{}{
		"status": "INACTIVE",
	}
	diff, err = newInstanceDiff("alicloud_config_rule", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleState": "INACTIVE",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "StopConfigRules" {
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
		err := resourceAlicloudConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ActiveConfigRule
	attributesDiff = map[string]interface{}{
		"status": "ACTIVE",
	}
	diff, err = newInstanceDiff("alicloud_config_rule", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleState": "ACTIVE",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ActiveConfigRules" {
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
		err := resourceAlicloudConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "ConfigRuleNotExists", "Invalid.ConfigRuleId.Value", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetConfigRule" {
				switch errorCode {
				case "{}", "ConfigRuleNotExists", "Invalid.ConfigRuleId.Value":
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
		err := resourceAlicloudConfigRuleRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}", "ConfigRuleNotExists", "Invalid.ConfigRuleId.Value":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigRuleDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "ConfigRuleNotExists"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteConfigRules" {
				switch errorCode {
				case "NonRetryableError", "ConfigRuleNotExists":
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
		err := resourceAlicloudConfigRuleDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "ConfigRuleNotExists":
			assert.Nil(t, err)
		}
	}

}
