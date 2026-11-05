package imp_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/imp"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// collection-order: bypass component_list -- pre-existing attribute outside this change (pure directory migration; schema unchanged, original flat layout carried no reorder coverage and was outside this gate's scope)
func TestAccAliCloudIMPAppTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_imp_app_template.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudIMPAppTemplateMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return imp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAppTemplate")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%simpapptemplate%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudIMPAppTemplateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
							"value": "tf-testAcc-jdD4qhGOujVlYcCUqTDUumAV.com",
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
	"component_list.#": acctest.CHECKSET,
	"config_list.#":    acctest.CHECKSET,
	"status":           acctest.CHECKSET,
}

func AlicloudIMPAppTemplateBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_imp_app_template": imp.ResourceAppTemplate(),
	}
}

// lintignore: R001
func TestUnitAlicloudImpAppTemplate(t *testing.T) {
	p := unitTestProviderResources()
	resourceSchema := p["alicloud_imp_app_template"]
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_imp_app_template", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImpClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = resourceSchema.Create(dInit, rawClient)
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
		err := resourceSchema.Create(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), nil)
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
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = resourceSchema.Update(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateAppTemplate
	attributesDiff := map[string]interface{}{
		"app_template_name": "UpdateAppTemplateValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_imp_app_template", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
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
		err := resourceSchema.Update(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dExisted.State(), nil)
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
	diff, err = acctest.NewInstanceDiff("alicloud_imp_app_template", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dExisted.State(), diff)
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
		err := resourceSchema.Update(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dExisted.State(), nil)
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
		err := resourceSchema.Read(dExisted, rawClient)
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
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = resourceSchema.Delete(dExisted, rawClient)
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
		err := resourceSchema.Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidAppTemplateId.App.NotFound":
			assert.Nil(t, err)
		}
	}
}
