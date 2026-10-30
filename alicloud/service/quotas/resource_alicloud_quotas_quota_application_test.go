package quotas_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/quotas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// The quota product does not support deletion, so skip the test.
// lintignore: AT001
func SkipTestAccAliCloudQuotasQuotaApplication_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_application.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudQuotasQuotaApplicationMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return quotas.NewQuotasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeQuotasQuotaApplication")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", AlicloudQuotasQuotaApplicationBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
	"status":      acctest.CHECKSET,
}

func AlicloudQuotasQuotaApplicationBasicDependence(name string) string {
	return ""
}

// lintignore: R001
func TestUnitAlicloudQuotasQuotaApplication(t *testing.T) {
	p := quotasTestProviderResources()
	resourceSchema := p["alicloud_quotas_quota_application"]
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_quotas_quota_application", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewQuotasClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSchema.Create(dInit, rawClient)
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
		err := resourceSchema.Read(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceSchema.Delete(dExisted, rawClient)
	assert.Nil(t, err)

}

// Test Quotas QuotaApplication. >>> Resource test cases, automatically generated.
// Case 3294
// lintignore: AT001
func SkipTestAccAliCloudQuotasQuotaApplication_basic3294(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_application.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudQuotasQuotaApplicationMap3294)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return quotas.NewQuotasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeQuotasQuotaApplication")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotasquotaapplication%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudQuotasQuotaApplicationBasicDependence3294)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "q_share_image-count",
					"product_code":      "gws",
					"quota_category":    "CommonQuota",
					"notice_type":       "3",
					"dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
					"desire_value": "53",
					"reason":       "测试",
					"env_language": "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "q_share_image-count",
						"product_code":      "gws",
						"quota_category":    "CommonQuota",
						"notice_type":       "3",
						"dimensions.#":      "1",
						"desire_value":      "53",
						"reason":            "测试",
						"env_language":      "zh",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"audit_mode", "env_language", "quota_category"},
			},
		},
	})
}

var AlicloudQuotasQuotaApplicationMap3294 = map[string]string{
	"status":            acctest.CHECKSET,
	"quota_description": acctest.CHECKSET,
	"create_time":       acctest.CHECKSET,
	"approve_value":     acctest.CHECKSET,
	"quota_name":        acctest.CHECKSET,
	"notice_type":       acctest.CHECKSET,
}

func AlicloudQuotasQuotaApplicationBasicDependence3294(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Case 3289
// lintignore: AT001
func TestAccAliCloudQuotasQuotaApplication_basic3289(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_application.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudQuotasQuotaApplicationMap3289)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return quotas.NewQuotasServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeQuotasQuotaApplication")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%squotasquotaapplication%d", acctest.DefaultRegionToTest, rand)
	currentTime := time.Now()
	sixMonthsLater := currentTime.AddDate(0, 6, 0)
	expireTime := sixMonthsLater.Format("2006-01-02T15:04:05Z")
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudQuotasQuotaApplicationBasicDependence3289)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
					"audit_mode":        "Sync",
					"effective_time":    "2023-05-22T16:00:00Z",
					"product_code":      "vpc",
					"quota_category":    "WhiteListLabel",
					"notice_type":       "3",
					"expire_time":       expireTime,
					"desire_value":      "1",
					"reason":            "",
					"env_language":      "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_action_code": "vpc_whitelist/ha_vip_whitelist",
						"audit_mode":        "Sync",
						"effective_time":    "2023-05-22T16:00:00Z",
						"product_code":      "vpc",
						"quota_category":    "WhiteListLabel",
						"notice_type":       "3",
						"expire_time":       expireTime,
						"desire_value":      "1",
						"reason":            "",
						"env_language":      "zh",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"audit_mode", "env_language", "quota_category", "status"},
			},
		},
	})
}

var AlicloudQuotasQuotaApplicationMap3289 = map[string]string{
	"status":            acctest.CHECKSET,
	"quota_description": acctest.CHECKSET,
	"create_time":       acctest.CHECKSET,
	"audit_reason":      acctest.CHECKSET,
	"approve_value":     acctest.CHECKSET,
	"quota_name":        acctest.CHECKSET,
	"notice_type":       acctest.CHECKSET,
}

func AlicloudQuotasQuotaApplicationBasicDependence3289(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test Quotas QuotaApplication. <<< Resource test cases, automatically generated.
