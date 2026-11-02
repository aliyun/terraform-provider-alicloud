package directmail_test

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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/directmail"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_direct_mail_mail_address",
		&resource.Sweeper{
			Name: "alicloud_direct_mail_mail_address",
			F:    testSweepDirectMailAddress,
		})
}
func testSweepDirectMailAddress(region string) error {
	rawClient, err := sweep.SharedClientForRegionWithBackendRegions(region, true, connectivity.DmSupportRegions)
	if err != nil {
		log.Printf("error getting Alicloud client: %s", err)
		return nil
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "QueryMailAddressByParam"
	request := map[string]interface{}{
		"PageSize": helper.PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch DirectMail Mail Address: %s", errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_direct_mail_mail_addresss", action, errs.AlibabaCloudSdkGoERROR))
			return nil
		}
		v, err := jsonpath.Get("$.data.mailAddress", response)
		if err != nil {
			log.Printf("[ERROR] Failed to parse DirectMail Mail Address: %s", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Groups", response))
			return nil
		}
		if len(v.([]interface{})) < 1 {
			log.Printf("[ERROR] Failed to fetch DirectMail Mail Address: %s", errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_direct_mail_mail_addresss", action, errs.AlibabaCloudSdkGoERROR))
			return nil
		}

		result, _ := v.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["AccountName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DirectMail Mail Address: %s", item["AccountName"].(string))
				continue
			}

			action := "DeleteMailAddress"
			request := map[string]interface{}{
				"MailAddressId": item["MailAddressId"],
			}
			_, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DirectMail Mail Address (%s): %s", item["AccountName"].(string), err)
			}
			log.Printf("[INFO] Delete DirectMail Mail Address success: %s ", item["AccountName"].(string))
		}

		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil

}

// this resource depends on user's account which tf account not support
func SkipTestAccAliCloudDirectMailMailAddress_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_mail_address.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDirectMailMailAddressMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return directmail.NewDmService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDirectMailMailAddress")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdirectmailmailaddress%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDirectMailMailAddressBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DmSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"sendtype":     "batch",
					"account_name": name + "s@xxx.changes.com.cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sendtype":     "batch",
						"account_name": name + "s@xxx.changes.com.cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reply_address": "r@xxx.changes.com.cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reply_address": "r@xxx.changes.com.cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Tf12345678password",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Tf12345678password",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reply_address": "r2@xxx.changes.com.cn",
					"password":      "Tf987654321password",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reply_address": "r2@xxx.changes.com.cn",
						"password":      "Tf987654321password",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AlicloudDirectMailMailAddressMap0 = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudDirectMailMailAddressBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDirectMailMailAddress(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	acctest.CheckoutSupportedRegions(t, true, connectivity.DmSupportRegions)
	dInit, _ := schema.InternalMap(p["alicloud_direct_mail_mail_address"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_direct_mail_mail_address"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"account_name":  "CreateMailAddressValue",
		"reply_address": "CreateMailAddressValue",
		"sendtype":      "CreateMailAddressValue",
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
		// QueryMailAddressByParam
		"data": map[string]interface{}{
			"mailAddress": []interface{}{
				map[string]interface{}{
					"AccountName":   "CreateMailAddressValue",
					"ReplyAddress":  "CreateMailAddressValue",
					"Sendtype":      "CreateMailAddressValue",
					"AccountStatus": "CreateMailAddressValue",
					"MailAddressId": "CreateMailAddressValue",
				},
			},
		},
		"MailAddressId": "CreateMailAddressValue",
	}
	CreateMockResponse := map[string]interface{}{
		"MailAddressId": "CreateMailAddressValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_direct_mail_mail_address", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = directmail.ResourceMailAddressCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		"MailAddressId": "CreateMailAddressValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateMailAddress" {
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
		err := directmail.ResourceMailAddressCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_direct_mail_mail_address"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = directmail.ResourceMailAddressUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	//ModifyMailAddress
	attributesDiff := map[string]interface{}{
		"password":      "ModifyMailAddressValue",
		"reply_address": "ModifyMailAddressValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_direct_mail_mail_address", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_direct_mail_mail_address"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// QueryMailAddressByParam Response
		"data": map[string]interface{}{
			"mailAddress": []interface{}{
				map[string]interface{}{
					"ReplyAddress": "ModifyMailAddressValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyMailAddress" {
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
		err := directmail.ResourceMailAddressUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_direct_mail_mail_address"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "QueryMailAddressByParam" {
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
		err := directmail.ResourceMailAddressRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    helper.String("loadEndpoint error"),
			Data:    helper.String("loadEndpoint error"),
			Message: helper.String("loadEndpoint error"),
		}
	})
	err = directmail.ResourceMailAddressDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidMailAddressId.Malformed"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteMailAddress" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Success": true,
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := directmail.ResourceMailAddressDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidMailAddressId.Malformed":
			assert.Nil(t, err)
		}
	}
}
