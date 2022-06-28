package alicloud

import (
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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// The quota product does not support deletion, so skip the test.
func SkipTestAccAlicloudQuotasQuotaApplication_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_application.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasQuotaApplicationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasQuotaApplication")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AlicloudQuotasQuotaApplicationBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"notice_type":       "0",
					"desire_value":      "60",
					"product_code":      "ess",
					"quota_action_code": "q_db_instance",
					"reason":            "For Terraform Test",
					"dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notice_type":       "0",
						"desire_value":      "60",
						"product_code":      "ess",
						"quota_action_code": "q_db_instance",
						"reason":            "For Terraform Test",
						"dimensions.#":      "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"quota_category"},
			},
		},
	})
}

var AlicloudQuotasQuotaApplicationMap = map[string]string{
	"notice_type": "0",
	"status":      CHECKSET,
}

func AlicloudQuotasQuotaApplicationBasicDependence(name string) string {
	return ""
}

func TestUnitAlicloudQuotasQuotaApplication(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_quotas_quota_application"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_quotas_quota_application"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"notice_type":       0,
		"desire_value":      60,
		"product_code":      "CreateQuotaApplicationValue",
		"quota_action_code": "CreateQuotaApplicationValue",
		"reason":            "CreateQuotaApplicationValue",
		"dimensions": []map[string]interface{}{
			{
				"key":   "CreateQuotaApplicationValue",
				"value": "CreateQuotaApplicationValue",
			},
		},
		"audit_mode":     "CreateQuotaApplicationValue",
		"quota_category": "CreateQuotaApplicationValue",
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
		// GetQuotaApplication
		"QuotaApplication": map[string]interface{}{
			"ApproveValue": "CreateQuotaApplicationValue",
			"AuditReason":  "CreateQuotaApplicationValue",
			"DesireValue":  60,
			"Dimension": map[string]interface{}{
				"CreateQuotaApplicationValue": "CreateQuotaApplicationValue",
			},
			"EffectiveTime":    "CreateQuotaApplicationValue",
			"ExpireTime":       "CreateQuotaApplicationValue",
			"NoticeType":       0,
			"ProductCode":      "CreateQuotaApplicationValue",
			"QuotaActionCode":  "CreateQuotaApplicationValue",
			"QuotaDescription": "CreateQuotaApplicationValue",
			"QuotaName":        "CreateQuotaApplicationValue",
			"QuotaUnit":        "CreateQuotaApplicationValue",
			"Reason":           "CreateQuotaApplicationValue",
			"Status":           "CreateQuotaApplicationValue",
		},
		"ApplicationId": "CreateQuotaApplicationValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateQuotaApplication
		"ApplicationId": "CreateQuotaApplicationValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_quotas_quota_application", errorCode))
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
	err = resourceAlicloudQuotasQuotaApplicationCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetQuotaApplication Response
		"ApplicationId": "CreateQuotaApplicationValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateQuotaApplication" {
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
		err := resourceAlicloudQuotasQuotaApplicationCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_quotas_quota_application"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetQuotaApplication" {
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
		err := resourceAlicloudQuotasQuotaApplicationRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudQuotasQuotaApplicationDelete(dExisted, rawClient)
	assert.Nil(t, err)

}
