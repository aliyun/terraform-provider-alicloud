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
	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ros_stack_group",
		&resource.Sweeper{
			Name: "alicloud_ros_stack_group",
			F:    testSweepRosStackGroup,
		})
}

func testSweepRosStackGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   region,
		"Status":     "ACTIVE",
	}
	var response map[string]interface{}
	action := "ListStackGroups"
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ros_stack_group", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.StackGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.StackGroups", response)
		}
		sweeped := false
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["StackGroupName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ros StackGroup: %s", item["StackGroupName"].(string))
				continue
			}
			sweeped = true
			action := "DeleteStackGroup"
			request := map[string]interface{}{
				"StackGroupName": item["StackGroupName"],
				"RegionId":       region,
			}
			_, err = client.RpcPost("ROS", "2019-09-10", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ros StackGroup (%s): %s", item["StackGroupName"].(string), err)
			}
			if sweeped {
				// Waiting 5 seconds to ensure Ros StackGroup have been deleted.
				time.Sleep(5 * time.Second)
			}
			log.Printf("[INFO] Delete Ros StackGroup success: %s ", item["StackGroupName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudROSStackGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_stack_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRosStackGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosStackGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudRosStackGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosStackGroupBasicDependence)
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
					"stack_group_name": name,
					"template_body":    `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "VpcName",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "InstanceType",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stack_group_name": name,
						"template_body":    CHECKSET,
						"parameters.#":     "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_ids", "operation_description", "template_body", "operation_preferences", "region_ids", "template_url", "template_version"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_body": `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Description\" : \"模板描述信息，可用于说明模板的适用场景、架构说明等。\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "tf-testacc",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "ECS",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_body": CHECKSET,
						"parameters.#":  "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"template_body": `{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}`,
					"parameters": []map[string]interface{}{
						{
							"parameter_key":   "VpcName",
							"parameter_value": "VpcName",
						},
						{
							"parameter_key":   "InstanceType",
							"parameter_value": "InstanceType",
						},
					},
					"description": "test for tf",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"template_body": CHECKSET,
						"parameters.#":  "2",
						"description":   "test for tf",
					}),
				),
			},
		},
	})
}

var AlicloudRosStackGroupMap = map[string]string{
	"stack_group_id": CHECKSET,
	"status":         CHECKSET,
}

func AlicloudRosStackGroupBasicDependence(name string) string {
	return ""
}

func TestUnitAlicloudRosStackGroup(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ros_stack_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ros_stack_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"stack_group_name": "CreateStackGroupValue",
		"template_body":    "CreateStackGroupValue",
		"parameters": []map[string]interface{}{
			{
				"parameter_key":   "CreateStackGroupValue",
				"parameter_value": "CreateStackGroupValue",
			},
		},
		"administration_role_name": "CreateStackGroupValue",
		"description":              "CreateStackGroupValue",
		"execution_role_name":      "CreateStackGroupValue",
		"template_url":             "CreateStackGroupValue",
		"template_version":         "CreateStackGroupValue",
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
		"StackGroup": map[string]interface{}{
			"StackGroupName":         "CreateStackGroupValue",
			"AdministrationRoleName": "CreateStackGroupValue",
			"Description":            "CreateStackGroupValue",
			"ExecutionRoleName":      "CreateStackGroupValue",
			"Parameters": []interface{}{
				map[string]interface{}{
					"ParameterKey":   "CreateStackGroupValue",
					"ParameterValue": "CreateStackGroupValue",
				},
			},
			"StackGroupId": "CreateStackGroupValue",
			"Status":       "ACTIVE",
			"TemplateBody": "CreateStackGroupValue",
		},
	}
	CreateMockResponse := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ros_stack_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudRosStackGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateStackGroup" {
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
		err := resourceAlicloudRosStackGroupCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ros_stack_group"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudRosStackGroupUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"administration_role_name": "UpdateStackGroupValue",
		"description":              "UpdateStackGroupValue",
		"execution_role_name":      "UpdateStackGroupValue",
		"parameters": []map[string]interface{}{
			{
				"parameter_key":   "UpdateStackGroupValue",
				"parameter_value": "UpdateStackGroupValue",
			},
		},
		"template_body":         "UpdateStackGroupValue",
		"account_ids":           "UpdateStackGroupValue",
		"operation_description": "UpdateStackGroupValue",
		"operation_preferences": "UpdateStackGroupValue",
		"region_ids":            "UpdateStackGroupValue",
		"template_url":          "UpdateStackGroupValue",
		"template_version":      "UpdateStackGroupValue",
	}
	diff, err := newInstanceDiff("alicloud_ros_stack_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ros_stack_group"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"StackGroup": map[string]interface{}{
			"StackGroupName":         "UpdateStackGroupValue",
			"AdministrationRoleName": "UpdateStackGroupValue",
			"Description":            "UpdateStackGroupValue",
			"ExecutionRoleName":      "UpdateStackGroupValue",
			"Parameters": []interface{}{
				map[string]interface{}{
					"ParameterKey":   "UpdateStackGroupValue",
					"ParameterValue": "UpdateStackGroupValue",
				},
			},
			"StackGroupId": "UpdateStackGroupValue",
			"Status":       "ACTIVE",
			"TemplateBody": "UpdateStackGroupValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateStackGroup" {
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
		err := resourceAlicloudRosStackGroupUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ros_stack_group"].Schema).Data(dExisted.State(), nil)
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
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetStackGroup" {
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
		err := resourceAlicloudRosStackGroupRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewRosClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudRosStackGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteStackGroup" {
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
		err := resourceAlicloudRosStackGroupDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
