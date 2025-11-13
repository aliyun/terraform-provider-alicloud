package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/PaesslerAG/jsonpath"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_polardb_parameter_group", &resource.Sweeper{
		Name: "alicloud_polardb_parameter_group",
		F:    testSweepPolarDBParameterGroup,
	})
}

func testSweepPolarDBParameterGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeParameterGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var response map[string]interface{}
	polarDBParameterGroupIds := make([]string, 0)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_parameter_group", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ParameterGroups", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ParameterGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		skip := true
		item := v.(map[string]interface{})
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["ParameterGroupName"])), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping PolarDBParameterGroup Instance: %v", item["ParameterGroupName"])
			continue
		}
		polarDBParameterGroupIds = append(polarDBParameterGroupIds, fmt.Sprint(item["ParameterGroupId"]))
	}

	for _, id := range polarDBParameterGroupIds {
		log.Printf("[INFO] Deleting PolarDBParameterGroup Instance: %s", id)
		deleteAction := "DeleteParameterGroup"
		if err != nil {
			return WrapError(err)
		}
		request = map[string]interface{}{
			"RegionId":         client.RegionId,
			"ParameterGroupId": id,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, err = client.RpcPost("polardb", "2017-08-01", deleteAction, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete PolarDBParameterGroup Instance (%s): %s", polarDBParameterGroupIds, err)
		}
	}
	return nil
}

func TestAccAliCloudPolarDBParameterGroup_basic00(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_parameter_group.default"
	ra := resourceAttrInit(resourceId, resourceAliCloudPolarDbParameterGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDBParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sPolarDBParameterGroup-name%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAliCloudPolarDbParameterGroupBasicDependence)
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
					"name":       "tf_testAcc",
					"db_type":    "MySQL",
					"db_version": "8.0",
					"parameters": []map[string]interface{}{
						{
							"param_name":  "wait_timeout",
							"param_value": "86400",
						},
						{
							"param_name":  "innodb_old_blocks_time",
							"param_value": "1000",
						},
					},
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         "tf_testAcc",
						"db_type":      "MySQL",
						"db_version":   "8.0",
						"parameters.#": "2",
						"description":  name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var resourceAliCloudPolarDbParameterGroupMap = map[string]string{}

func resourceAliCloudPolarDbParameterGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}

func TestUnitAlicloudPolarDBParameterGroup(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_polardb_parameter_group"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_polardb_parameter_group"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"name":       "CreatePolarDBParameterGroup",
		"db_type":    "CreatePolarDBParameterGroup",
		"db_version": "CreatePolarDBParameterGroup",
		"parameters": []interface{}{
			map[string]interface{}{
				"param_name":  "CreatePolarDBParameterGroup",
				"param_value": "CreatePolarDBParameterGroup",
			},
		},
		"description": "CreatePolarDBParameterGroup",
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
		// DescribeParameterGroup
		"ParameterGroup": []interface{}{
			map[string]interface{}{
				"ParameterGroupName": "CreatePolarDBParameterGroup",
				"DBType":             "CreatePolarDBParameterGroup",
				"DBVersion":          "CreatePolarDBParameterGroup",
				"ParameterDetail": []interface{}{
					map[string]interface{}{
						"ParamName":  "CreatePolarDBParameterGroup",
						"ParamValue": "CreatePolarDBParameterGroup",
					},
				},
				"ParameterGroupDesc": "CreatePolarDBParameterGroup",
			},
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_polardb_parameter_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPolarDBClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPolarDbParameterGroupCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateParameterGroup" {
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
		err := resourceAliCloudPolarDbParameterGroupCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_polardb_parameter_group"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
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
			if *action == "DescribeParameterGroup" {
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
		err := resourceAliCloudPolarDbParameterGroupRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewPolarDBClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudPolarDbParameterGroupDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteParameterGroup" {
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
			if *action == "DescribeParameterGroup" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudPolarDbParameterGroupDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

// Test PolarDb ParameterGroup. >>> Resource test cases, automatically generated.
// Case 属性全覆盖_多参数 11763
func TestAccAliCloudPolarDbParameterGroup_basic11763(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlicloudPolarDbParameterGroupMap11763)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolarDbParameterGroupBasicDependence11763)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_version":           "8.0",
					"parameter_group_name": name,
					"db_type":              "MySQL",
					"parameters": []map[string]interface{}{
						{
							"param_name":  "wait_timeout",
							"param_value": "86400",
						},
						{
							"param_name":  "innodb_old_blocks_time",
							"param_value": "1000",
						},
						{
							"param_name":  "default_time_zone",
							"param_value": "SYSTEM",
						},
					},
					"description": "tf_testAcc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_version":           CHECKSET,
						"parameter_group_name": name,
						"db_type":              "MySQL",
						"parameters.#":         "3",
						"description":          "tf_testAcc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudPolarDbParameterGroupMap11763 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPolarDbParameterGroupBasicDependence11763(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test PolarDb ParameterGroup. <<< Resource test cases, automatically generated.
