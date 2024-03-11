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

func TestAccAlicloudQuotasQuotaAlarm_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_alarm.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasQuotaAlarmMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasQuotaAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccQuotasQuotaAlarmTest%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasQuotaAlarmBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_alarm_name":  name,
					"product_code":      "ecs",
					"quota_action_code": "q_prepaid-instance-count-per-once-purchase",
					"threshold":         "100",
					"quota_dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name":   name,
						"product_code":       "ecs",
						"quota_action_code":  "q_prepaid-instance-count-per-once-purchase",
						"threshold":          "100",
						"quota_dimensions.#": "1",
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
					"quota_alarm_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_alarm_name": name,
					"threshold":        "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name": name,
						"threshold":        "100",
					}),
				),
			},
		},
	})
}

var AlicloudQuotasQuotaAlarmMap = map[string]string{}

func AlicloudQuotasQuotaAlarmBasicDependence(name string) string {
	return ""
}

func TestUnitAlicloudQuotasQuotaAlarm(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_quotas_quota_alarm"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_quotas_quota_alarm"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"quota_alarm_name":  "CreateQuotaAlarmValue",
		"product_code":      "CreateQuotaAlarmValue",
		"quota_action_code": "CreateQuotaAlarmValue",
		"threshold":         100,
		"quota_dimensions": []map[string]interface{}{
			{
				"key":   "CreateQuotaAlarmValue",
				"value": "CreateQuotaAlarmValue",
			},
		},
		"threshold_percent": 100,
		"web_hook":          "CreateQuotaAlarmValue",
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
		// GetQuotaAlarm
		"QuotaAlarm": map[string]interface{}{
			"ProductCode":     "CreateQuotaAlarmValue",
			"QuotaActionCode": "CreateQuotaAlarmValue",
			"AlarmName":       "CreateQuotaAlarmValue",
			"QuotaDimension": map[string]interface{}{
				"CreateQuotaAlarmValue": "CreateQuotaAlarmValue",
			},
			"Threshold":        100,
			"ThresholdPercent": 100,
		},
		"AlarmId": "CreateQuotaAlarmValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateQuotaAlarm
		"AlarmId": "CreateQuotaAlarmValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_quotas_quota_alarm", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewQuotasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudQuotasQuotaAlarmCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetQuotaAlarm Response
		"AlarmId": "CreateQuotaAlarmValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateQuotaAlarm" {
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
		err := resourceAliCloudQuotasQuotaAlarmCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_quotas_quota_alarm"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewQuotasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudQuotasQuotaAlarmUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateQuotaAlarm
	attributesDiff := map[string]interface{}{
		"quota_alarm_name":  "UpdateQuotaAlarmValue",
		"threshold":         200,
		"threshold_percent": 200,
		"web_hook":          "UpdateQuotaAlarmValue",
	}
	diff, err := newInstanceDiff("alicloud_quotas_quota_alarm", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_quotas_quota_alarm"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetQuotaAlarm Response
		"QuotaAlarm": map[string]interface{}{
			"AlarmName":        "UpdateQuotaAlarmValue",
			"Threshold":        200,
			"ThresholdPercent": 200,
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateQuotaAlarm" {
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
		err := resourceAliCloudQuotasQuotaAlarmUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_quotas_quota_alarm"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetQuotaAlarm" {
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
		err := resourceAliCloudQuotasQuotaAlarmRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewQuotasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudQuotasQuotaAlarmDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteQuotaAlarm" {
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
		err := resourceAliCloudQuotasQuotaAlarmDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}

func TestAccAlicloudQuotasQuotaAlarm_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_alarm.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasQuotaAlarmMap2901)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasQuotaAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotasquotaalarm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasQuotaAlarmBasicDependence2901)
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
					"quota_action_code": "ecs.hfg7.xlarge",
					"quota_dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
						{
							"key":   "chargeType",
							"value": "PostPaid",
						},
						{
							"key":   "zoneId",
							"value": "cn-hangzhou-k",
						},
						{
							"key":   "networkType",
							"value": "vpc",
						},
						{
							"key":   "resourceType",
							"value": "InstanceType",
						},
					},
					"threshold_percent": "80",
					"product_code":      "ecs-spec",
					"quota_alarm_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code":  "ecs.hfg7.xlarge",
						"quota_dimensions.#": "5",
						"threshold_percent":  "80",
						"product_code":       "ecs-spec",
						"quota_alarm_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"web_hook": "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"web_hook": "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_type": "used",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_type": "used",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_percent": "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_percent": "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_alarm_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_percent": "20",
					"threshold_type":    "usable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_percent": "20",
						"threshold_type":    "usable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "ecs.hfg7.xlarge",
					"quota_dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
						{
							"key":   "chargeType",
							"value": "PostPaid",
						},
						{
							"key":   "zoneId",
							"value": "cn-hangzhou-k",
						},
						{
							"key":   "networkType",
							"value": "vpc",
						},
						{
							"key":   "resourceType",
							"value": "InstanceType",
						},
					},
					"threshold_percent": "80",
					"product_code":      "ecs-spec",
					"quota_alarm_name":  name + "_update",
					"web_hook":          "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
					"threshold_type":    "used",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code":  "ecs.hfg7.xlarge",
						"quota_dimensions.#": "5",
						"threshold_percent":  "80",
						"product_code":       "ecs-spec",
						"quota_alarm_name":   name + "_update",
						"web_hook":           "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
						"threshold_type":     "used",
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

// Test Quotas QuotaAlarm. >>> Resource test cases, automatically generated.
// Case 2901
func TestAccAlicloudQuotasQuotaAlarm_basic2901(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_alarm.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasQuotaAlarmMap2901)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasQuotaAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotasquotaalarm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasQuotaAlarmBasicDependence2901)
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
					"quota_action_code": "q_desktop-count",
					"quota_dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
					"threshold_percent": "80",
					"product_code":      "gws",
					"quota_alarm_name":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code":  "q_desktop-count",
						"quota_dimensions.#": "1",
						"threshold_percent":  "80",
						"product_code":       "gws",
						"quota_alarm_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"web_hook": "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"web_hook": "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_type": "used",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_type": "used",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_percent": "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_percent": "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_alarm_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold_percent": "20",
					"threshold_type":    "usable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold_percent": "20",
						"threshold_type":    "usable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "q_desktop-count",
					"quota_dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
					"threshold_percent": "80",
					"product_code":      "gws",
					"quota_alarm_name":  name + "_update",
					"web_hook":          "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
					"threshold_type":    "used",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code":  "q_desktop-count",
						"quota_dimensions.#": "1",
						"threshold_percent":  "80",
						"product_code":       "gws",
						"quota_alarm_name":   name + "_update",
						"web_hook":           "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
						"threshold_type":     "used",
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

var AlicloudQuotasQuotaAlarmMap2901 = map[string]string{
	"create_time":    CHECKSET,
	"threshold_type": "used",
}

func AlicloudQuotasQuotaAlarmBasicDependence2901(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 2936
func TestAccAlicloudQuotasQuotaAlarm_basic2936(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_alarm.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasQuotaAlarmMap2936)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasQuotaAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotasquotaalarm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasQuotaAlarmBasicDependence2936)
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
					"quota_action_code": "q_user_poc_money_consumption",
					"product_code":      "computenest",
					"quota_alarm_name":  name,
					"threshold":         "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "q_user_poc_money_consumption",
						"product_code":      "computenest",
						"quota_alarm_name":  name,
						"threshold":         "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_alarm_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold": "1001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold": "1001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "q_user_poc_money_consumption",
					"product_code":      "computenest",
					"quota_alarm_name":  name + "_update",
					"threshold":         "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "q_user_poc_money_consumption",
						"product_code":      "computenest",
						"quota_alarm_name":  name + "_update",
						"threshold":         "1000",
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

var AlicloudQuotasQuotaAlarmMap2936 = map[string]string{
	"create_time":    CHECKSET,
	"threshold_type": "used",
}

func AlicloudQuotasQuotaAlarmBasicDependence2936(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 2901  twin
func TestAccAlicloudQuotasQuotaAlarm_basic2901_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_alarm.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasQuotaAlarmMap2901)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasQuotaAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotasquotaalarm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasQuotaAlarmBasicDependence2901)
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
					"quota_action_code": "q_desktop-count",
					"quota_dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
					"threshold_percent": "20",
					"product_code":      "gws",
					"quota_alarm_name":  name,
					"web_hook":          "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
					"threshold_type":    "usable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code":  "q_desktop-count",
						"quota_dimensions.#": "1",
						"threshold_percent":  "20",
						"product_code":       "gws",
						"quota_alarm_name":   name,
						"web_hook":           "https://oapi.dingtalk.com/robot/send?access_token=0a09bd617f43d07e8607b258c6cdffbacf0e023f1bbe46cfeb0265127802bf43",
						"threshold_type":     "usable",
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

// Case 2936  twin
func TestAccAlicloudQuotasQuotaAlarm_basic2936_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_alarm.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasQuotaAlarmMap2936)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasQuotaAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotasquotaalarm%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasQuotaAlarmBasicDependence2936)
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
					"quota_action_code": "q_user_poc_money_consumption",
					"product_code":      "computenest",
					"quota_alarm_name":  name,
					"threshold":         "1001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "q_user_poc_money_consumption",
						"product_code":      "computenest",
						"quota_alarm_name":  name,
						"threshold":         "1001",
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

// Test Quotas QuotaAlarm. <<< Resource test cases, automatically generated.
