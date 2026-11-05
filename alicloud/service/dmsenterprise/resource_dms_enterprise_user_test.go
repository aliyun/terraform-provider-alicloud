package dmsenterprise_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	gomonkey "github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	alicloud "github.com/aliyun/terraform-provider-alicloud/alicloud"
	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	dmsenterprise "github.com/aliyun/terraform-provider-alicloud/alicloud/service/dmsenterprise"
	tfsweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_dms_enterprise_user", &resource.Sweeper{
		Name: "alicloud_dms_enterprise_user",
		F:    testSweepDMSEnterpriseUsers,
	})
}

func testSweepDMSEnterpriseUsers(region string) error {

	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return tferrs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"testacc",
	}
	request := map[string]interface{}{
		"UserState":  "NORMAL",
		"PageSize":   tfhelper.PageSizeXLarge,
		"PageNumber": 1,
	}
	var response map[string]interface{}
	action := "ListUsers"
	for {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, true)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.DataDefaultErrorMsg, "alicloud_dms_enterprise_users", action, tferrs.AlibabaCloudSdkGoERROR)
		}
		tflogs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.UserList.User", response)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, action, "$.UserList.User", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if _, ok := item["NickName"]; !ok {
				skip = false
			} else {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprintf("%v", item["NickName"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DMS Enterprise User: %v", item["NickName"])
				continue
			}
			action := "DeleteUser"
			request := map[string]interface{}{
				"Uid": item["Uid"],
			}
			_, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete DMS Enterprise User (%v): %s", item["NickName"], err)
				continue
			}

			log.Printf("[INFO] Delete DMS Enterprise User Success: %v ", item["NickName"])
		}
		if len(result) < tfhelper.PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAliCloudDMSEnterpriseUser_basic(t *testing.T) {
	tfacctest.CheckoutSupportedRegions(t, true, connectivity.DMSEnterpriseSupportRegions)
	var v map[string]interface{}
	resourceId := "alicloud_dms_enterprise_user.default"
	ra := tfacctest.ResourceAttrInit(resourceId, DmsEnterpriseUserMap)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dmsenterprise.NewDmsEnterpriseService(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDmsEnterpriseUser")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccDmsEnterpriseUser%d", rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, DmsEnterpriseUserBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheck(t)
		},

		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"uid":        "${alicloud_ram_user.user.id}",
					"nick_name":  name,
					"role_names": []string{"DBA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name":    name,
						"role_names.#": "1",
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
					"max_execute_count": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_execute_count": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_result_count": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_result_count": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nick_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "NORMAL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_execute_count": "1000",
					"max_result_count":  "1000",
					"nick_name":         name + "change",
					"role_names":        []string{"DBA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_execute_count": "1000",
						"max_result_count":  "1000",
						"nick_name":         name + "change",
						"role_names.#":      "1",
					}),
				),
			},
		},
	})
}

var DmsEnterpriseUserMap = map[string]string{
	"status": tfacctest.CHECKSET,
}

func DmsEnterpriseUserBasicdependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name         = "%s"
	  display_name = "user_display_name"
	  mobile       = "86-18688888888"
	  email        = "hello.uuu@aaa.com"
	  comments     = "yoyoyo"
	}`, name)
}

// lintignore: R001
func TestUnitAlicloudDMSEnterpriseUser(t *testing.T) {
	p := alicloud.Provider().ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"mobile":     "RegisterUserValue",
		"role_names": []string{"RegisterUserValue"},
		"tid":        1,
		"uid":        "RegisterUserValue",
		"user_name":  "RegisterUserValue",
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
	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// GetUser
		"User": map[string]interface{}{
			"Uid":    "RegisterUserValue",
			"Mobile": "RegisterUserValue",
			"RoleNameList": map[string]interface{}{
				"RoleNames": []string{"RegisterUserValue"},
			},
			"State":    "RegisterUserValue",
			"NickName": "RegisterUserValue",
		},
	}
	CreateMockResponse := map[string]interface{}{}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tfhelper.String(errorCode),
			Data:       tfhelper.String(errorCode),
			Message:    tfhelper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, tferrs.GetNotFoundErrorFromString(tferrs.GetNotFoundMessage("alicloud_dms_enterprise_user", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			tfhelper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmsenterpriseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tfhelper.String("loadEndpoint error"),
			Data:       tfhelper.String("loadEndpoint error"),
			Message:    tfhelper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = dmsenterprise.ResourceAlicloudDmsEnterpriseUserCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RegisterUser" {
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
		err := dmsenterprise.ResourceAlicloudDmsEnterpriseUserCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmsenterpriseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tfhelper.String("loadEndpoint error"),
			Data:       tfhelper.String("loadEndpoint error"),
			Message:    tfhelper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = dmsenterprise.ResourceAlicloudDmsEnterpriseUserUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateUser
	attributesDiff := map[string]interface{}{
		"mobile":            "UpdateUserValue",
		"role_names":        []string{"UpdateUserValue"},
		"user_name":         "UpdateUserValue",
		"nick_name":         "UpdateUserValue",
		"max_execute_count": 1,
		"max_result_count":  1,
		"tid":               1,
	}
	diff, err := tfacctest.NewInstanceDiff("alicloud_dms_enterprise_user", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetUser Response
		"User": map[string]interface{}{
			"Mobile": "UpdateUserValue",
			"RoleNameList": map[string]interface{}{
				"RoleNames": []string{"UpdateUserValue"},
			},
			"NickName": "UpdateUserValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateUser" {
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
		err := dmsenterprise.ResourceAlicloudDmsEnterpriseUserUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// DisableUser
	attributesDiff = map[string]interface{}{
		"status": "DISABLE",
		"tid":    2,
	}
	diff, err = tfacctest.NewInstanceDiff("alicloud_dms_enterprise_user", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetUser Response
		"User": map[string]interface{}{
			"State": "DISABLE",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DisableUser" {
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
		err := dmsenterprise.ResourceAlicloudDmsEnterpriseUserUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// EnableUser
	attributesDiff = map[string]interface{}{
		"status": "NORMAL",
		"tid":    2,
	}
	diff, err = tfacctest.NewInstanceDiff("alicloud_dms_enterprise_user", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetUser Response
		"User": map[string]interface{}{
			"State": "NORMAL",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "EnableUser" {
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
		err := dmsenterprise.ResourceAlicloudDmsEnterpriseUserUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dms_enterprise_user"].Schema).Data(dExisted.State(), nil)
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetUser" {
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
		err := dmsenterprise.ResourceAlicloudDmsEnterpriseUserRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmsenterpriseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       tfhelper.String("loadEndpoint error"),
			Data:       tfhelper.String("loadEndpoint error"),
			Message:    tfhelper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = dmsenterprise.ResourceAlicloudDmsEnterpriseUserDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteUser" {
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
		err := dmsenterprise.ResourceAlicloudDmsEnterpriseUserDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
