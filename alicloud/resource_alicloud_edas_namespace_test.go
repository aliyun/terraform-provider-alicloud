package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	roa "github.com/alibabacloud-go/tea-roa/client"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_edas_namespace",
		&resource.Sweeper{
			Name: "alicloud_edas_namespace",
			F:    testSweepEdasNamespace,
		})
}

func testSweepEdasNamespace(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.EdasSupportedRegions) {
		log.Printf("[INFO] Skipping Edas Namespace unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "/pop/v5/user_region_defs"
	request := map[string]*string{}

	var response map[string]interface{}
	conn, err := aliyunClient.NewEdasClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-08-01"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	if fmt.Sprint(response["Code"]) != "200" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}

	resp, err := jsonpath.Get("$.UserDefineRegionList.UserDefineRegionEntity", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.UserDefineRegionList.UserDefineRegionEntity", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["RegionName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Edas Namespace: %s", item["RegionName"].(string))
			continue
		}
		action := "/pop/v5/user_region_def"
		request := map[string]*string{
			"Id": StringPointer(item["Id"].(string)),
		}
		_, err = conn.DoRequest(StringPointer("2017-08-01"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Edas Namespace (%s): %s", item["RegionName"].(string), err)
		}
		log.Printf("[INFO] Delete Edas Namespace success: %s ", item["RegionName"].(string))
	}

	return nil
}

func TestAccAlicloudEDASNamespace_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_edas_namespace.default"
	checkoutSupportedRegions(t, true, connectivity.EdasSupportedRegions)
	ra := resourceAttrInit(resourceId, AlicloudEDASNamespaceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEdasNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEDASNamespaceBasicDependence0)
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
					"debug_enable":         "false",
					"description":          "${var.name}",
					"namespace_name":       "${var.name}",
					"namespace_logical_id": "${var.logical_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"debug_enable":         "false",
						"description":          name,
						"namespace_name":       name,
						"namespace_logical_id": defaultRegionToTest + ":" + name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespace_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"debug_enable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"debug_enable": "true",
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

var AlicloudEDASNamespaceMap0 = map[string]string{
	"namespace_logical_id": CHECKSET,
	"namespace_name":       CHECKSET,
	"debug_enable":         CHECKSET,
}

func AlicloudEDASNamespaceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
variable "logical_id" {
  default = "%s:%s"
}
`, name, defaultRegionToTest, name)
}

func TestUnitAccAlicloudEDASNamespace(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_edas_namespace"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_edas_namespace"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"debug_enable":         false,
		"description":          "CreateEdasNamespaceValue",
		"namespace_name":       "CreateEdasNamespaceValue",
		"namespace_logical_id": "CreateEdasNamespaceValue",
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
		"body": map[string]interface{}{
			"UserDefineRegionList": map[string]interface{}{
				"UserDefineRegionEntity": []interface{}{
					map[string]interface{}{
						"RegionName":   "CreateEdasNamespaceValue",
						"Description":  "CreateEdasNamespaceValue",
						"UserId":       "CreateEdasNamespaceValue",
						"DebugEnable":  false,
						"Id":           "CreateEdasNamespaceValue",
						"RegionId":     "CreateEdasNamespaceValue",
						"BelongRegion": "CreateEdasNamespaceValue",
					},
				},
			},
			"Code": 200,
		},
	}
	CreateMockResponse := map[string]interface{}{
		"body": map[string]interface{}{
			"UserDefineRegionEntity": map[string]interface{}{
				"Id": "CreateEdasNamespaceValue",
			},
			"Code": 200,
		},
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_edas_namespace", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEdasClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEdasNamespaceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, action *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "/pop/v5/user_region_def" {
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
		err := resourceAlicloudEdasNamespaceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_edas_namespace"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEdasClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEdasNamespaceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"debug_enable":   true,
		"description":    "UpdateEdasNamespaceValue",
		"namespace_name": "UpdateEdasNamespaceValue",
	}
	diff, err := newInstanceDiff("alicloud_edas_namespace", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_edas_namespace"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"body": map[string]interface{}{
			"UserDefineRegionList": map[string]interface{}{
				"UserDefineRegionEntity": []interface{}{
					map[string]interface{}{
						"RegionName":  "UpdateEdasNamespaceValue",
						"Description": "UpdateEdasNamespaceValue",
						"DebugEnable": true,
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, action *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "/pop/v5/user_region_def" {
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
		err := resourceAlicloudEdasNamespaceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_edas_namespace"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_edas_namespace", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_edas_namespace"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, action *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "/pop/v5/user_region_defs" {
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
		err := resourceAlicloudEdasNamespaceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEdasClient", func(_ *connectivity.AliyunClient) (*roa.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEdasNamespaceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_edas_namespace", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_edas_namespace"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "notFoundBody", "failCode"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest", func(_ *roa.Client, _ *string, _ *string, _ *string, _ *string, action *string, _ map[string]*string, _ map[string]*string, _ interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "/pop/v5/user_region_def" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"body": map[string]interface{}{
								"Code": 200,
							},
						}
						return ReadMockResponse, nil
					}
					if errorCodes[retryIndex] == "notFoundBody" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					if errorCodes[retryIndex] == "failCode" {
						ReadMockResponse = map[string]interface{}{
							"body": map[string]interface{}{
								"Code": 400,
							},
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEdasNamespaceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "notFoundBody":
			assert.NotNil(t, err)
		case "failCode":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
