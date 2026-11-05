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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/directmail"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_direct_mail_receivers", &resource.Sweeper{
		Name: "alicloud_direct_mail_receivers",
		F:    testSweepDirectMailReceivers,
	})
}

func testSweepDirectMailReceivers(region string) error {
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
	action := "QueryReceiverByParam"
	request := map[string]interface{}{
		"PageSize": helper.PageSizeLarge,
		"PageNo":   1,
	}
	Ids := make([]string, 0)
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
		if err != nil {
			log.Println("List Direct Mail Receivers Failed!", err)
			return nil
		}
		resp, err := jsonpath.Get("$.data.receiver", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.data.receiver", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["ReceiversName"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Direct Mail Receivers: %v (%v)", item["ReceiversName"], item["ReceiverId"])
				continue
			}
			Ids = append(Ids, fmt.Sprint(item["ReceiverId"]))
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNo"] = request["PageNo"].(int) + 1
	}
	for _, Id := range Ids {
		log.Printf("[INFO] Deleting Direct Mail Receivers: (%s)", Id)
		action := "DeleteReceiver"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(time.Minute*9, func() *retry.RetryError {
			_, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed To Delete Direct Mail Receivers : %s", err)
			continue
		}
		log.Printf("[INFO] Delete Direct Mail Receivers Success : %s", Id)
	}
	return nil
}

func TestAccAliCloudDirectMailReceivers_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_direct_mail_receivers.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDirectMailReceiversMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return directmail.NewDmService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDirectMailReceivers")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDirectMailReceiversBasicDependence0)
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
					"receivers_name":  "${var.name}",
					"receivers_alias": name + "@onaliyun.com",
					"description":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"receivers_name":  name,
						"receivers_alias": name + "@onaliyun.com",
						"description":     name,
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

var AlicloudDirectMailReceiversMap0 = map[string]string{}

func AlicloudDirectMailReceiversBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDirectMailReceivers(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	acctest.CheckoutSupportedRegions(t, true, connectivity.DmSupportRegions)
	dInit, _ := schema.InternalMap(p["alicloud_direct_mail_receivers"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_direct_mail_receivers"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"receivers_name":  "CreateReceiverValue",
		"receivers_alias": "CreateReceiverValue",
		"description":     "CreateReceiverValue",
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
		// QueryReceiverByParam
		"data": map[string]interface{}{
			"receiver": []interface{}{
				map[string]interface{}{
					"Desc":            "CreateReceiverValue",
					"ReceiversAlias":  "CreateReceiverValue",
					"ReceiversName":   "CreateReceiverValue",
					"ReceiversStatus": 1,
					"ReceiverId":      "CreateReceiverValue",
				},
			},
		},
		"ReceiverId": "CreateReceiverValue",
	}
	CreateMockResponse := map[string]interface{}{
		"ReceiverId": "CreateReceiverValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_direct_mail_receivers", errorCode))
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
	err = directmail.ResourceReceiversCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		"MailAddressId": "CreateReceiverValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateReceiver" {
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
		err := directmail.ResourceReceiversCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_direct_mail_receivers"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "QueryReceiverByParam" {
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
		err := directmail.ResourceReceiversRead(dExisted, rawClient)
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
	err = directmail.ResourceReceiversDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteReceiver" {
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
		err := directmail.ResourceReceiversDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
