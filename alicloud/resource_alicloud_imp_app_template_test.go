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

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudIMPAppTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_imp_app_template.default"
	ra := resourceAttrInit(resourceId, AlicloudIMPAppTemplateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ImpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeImpAppTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%simpapptemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudIMPAppTemplateBasicDependence0)
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
					"scene":             "business",
					"integration_mode":  "paasSDK",
					"component_list":    []string{"component.live"},
					"app_template_name": "tf_testAcc_GWcpq51dSi5td18Qd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scene":             "business",
						"integration_mode":  "paasSDK",
						"component_list.#":  "1",
						"app_template_name": "tf_testAcc_GWcpq51dSi5td18Qd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_template_name": "tf_testAcc_IN1u0gHPAo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_template_name": "tf_testAcc_IN1u0gHPAo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_list": []map[string]interface{}{
						{
							"key":   "config.appCallbackAuthKey",
							"value": "tf-testAcc-jdD4qhGOujVlYcCUqTDUumAV",
						},
						{
							"key":   "config.appCallbackUrl",
							"value": "http://aliyun.com/tf-testAcc-jdD4qhGOujVlYcCUqTDUumAV",
						},
						{
							"key":   "config.callbackClass.live",
							"value": "config.callbackEvent.liveStatus",
						},
						{
							"key":   "config.callbackClass.user",
							"value": "config.callbackEvent.userEnterRoom",
						},
						{
							"key":   "config.livePullDomain",
							"value": "tf-testAcc-jdD4qhGOujVlYcCUqTDUumAV.com",
						},
						{
							"key":   "config.livePushDomain",
							"value": "tf-testAcc-jdD4qhGOujVlYcCUqTD.com",
						},
						{
							"key":   "config.multipleClientsLogin",
							"value": "false",
						},
						{
							"key":   "config.regionId",
							"value": "cn-hangzhou",
						},
						{
							"key":   "config.streamChangeCallbackUrl",
							"value": "https://aliyun.com/tf-testAcc-jdD4qhGOujVlYcCU",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_list.#": "9",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_template_name": "tf_testAcc_tqPHQU5xU",
					"config_list": []map[string]interface{}{
						{
							"key":   "config.appCallbackAuthKey",
							"value": "tf-testAcc-jdD4qhGOxxxxxxxx",
						},
						{
							"key":   "config.appCallbackUrl",
							"value": "http://aliyun.com/tf-testAcc-jdD4qhGOxxxxxxxx",
						},
						{
							"key":   "config.callbackClass.live",
							"value": "config.callbackEvent.liveStatusUpdate",
						},
						{
							"key":   "config.callbackClass.user",
							"value": "config.callbackEvent.userEnterRoomUpdate",
						},
						{
							"key":   "config.livePullDomain",
							"value": "tf-testAcc-jdD4qhGOxxxxxxxx.com",
						},
						{
							"key":   "config.livePushDomain",
							"value": "tf-testAcc-jdD4qhGOxxxxxxxx.com",
						},
						{
							"key":   "config.multipleClientsLogin",
							"value": "true",
						},
						{
							"key":   "config.regionId",
							"value": "cn-shanghai",
						},
						{
							"key":   "config.streamChangeCallbackUrl",
							"value": "https://aliyun.com/tf-testAcc-jdD4qhGOxxxxxxxx",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_template_name": "tf_testAcc_tqPHQU5xU",
						"config_list.#":     "9",
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

var AlicloudIMPAppTemplateMap0 = map[string]string{
	"component_list.#": CHECKSET,
	"config_list.#":    CHECKSET,
	"status":           CHECKSET,
}

func AlicloudIMPAppTemplateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAlicloudImpAppTemplate(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_imp_app_template"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_imp_app_template"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"app_template_name": "CreateAppTemplateValue",
		"component_list":    []interface{}{"CreateAppTemplateValue0", "CreateAppTemplateValue1"},
		"integration_mode":  "paasSDK",
		"scene":             "business",
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
		// GetAppTemplate
		"Result": map[string]interface{}{
			"AppTemplateName": "CreateAppTemplateValue",
			"ComponentList":   "CreateAppTemplateValue",
			"IntegrationMode": "paasSDK",
			"Scene":           "business",
			"Status":          "CreateAppTemplateValue",
			"ConfigList": []interface{}{
				map[string]interface{}{
					"Key":   "CreateAppTemplateValue",
					"Value": "CreateAppTemplateValue",
				},
			},
			"AppTemplateId": "CreateAppTemplateValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateAppTemplate
		"Result": map[string]interface{}{
			"AppTemplateId": "CreateAppTemplateValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_imp_app_template", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImpClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudImpAppTemplateCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		"Result": map[string]interface{}{
			"AppTemplateId": "CreateAppTemplateValue",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateAppTemplate" {
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
		err := resourceAlicloudImpAppTemplateCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_imp_app_template"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImpClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudImpAppTemplateUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateAppTemplate
	attributesDiff := map[string]interface{}{
		"app_template_name": "UpdateAppTemplateValue",
	}
	diff, err := newInstanceDiff("alicloud_imp_app_template", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_imp_app_template"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAppTemplate Response
		"Result": map[string]interface{}{
			"AppTemplateName": "UpdateAppTemplateValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateAppTemplate" {
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
		err := resourceAlicloudImpAppTemplateUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_imp_app_template"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// UpdateAppTemplateConfig
	attributesDiff = map[string]interface{}{
		"config_list": []interface{}{
			map[string]interface{}{
				"key":   "UpdateAppTemplateConfig",
				"value": "UpdateAppTemplateConfig",
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_imp_app_template", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_imp_app_template"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Result": map[string]interface{}{
			"ConfigList": []interface{}{
				map[string]interface{}{
					"Key":   "UpdateAppTemplateConfig",
					"Value": "UpdateAppTemplateConfig",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateAppTemplateConfig" {
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
		err := resourceAlicloudImpAppTemplateUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_imp_app_template"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidAppTemplateId.App.NotFound", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetAppTemplate" {
				switch errorCode {
				case "{}", "InvalidAppTemplateId.App.NotFound":
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
		err := resourceAlicloudImpAppTemplateRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}", "InvalidAppTemplateId.App.NotFound":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImpClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudImpAppTemplateDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidAppTemplateId.App.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteAppTemplate" {
				switch errorCode {
				case "NonRetryableError", "InvalidAppTemplateId.App.NotFound":
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
		err := resourceAlicloudImpAppTemplateDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidAppTemplateId.App.NotFound":
			assert.Nil(t, err)
		}
	}
}
