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

func TestAccAlicloudResourceManagerAccount_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_account.default"
	ra := resourceAttrInit(resourceId, ResourceManagerAccountMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerAccountBasicdependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckEnterpriseAccountEnabled(t)
			testAccPreCheck(t)
			// The created account resource cannot be deleted, and the created dependent resource folder will also be deleted. Therefore, the existing folder is specified in the environment variable for testing. If it is not specified, skip the test.
			testAccPreCheckWithResourceManagerFloderIdSetting(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
					"folder_id":    os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID1"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": name,
						"folder_id":    os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID1"),
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
					"display_name": fmt.Sprintf("%s1", name),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": fmt.Sprintf("%s1", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"folder_id": os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID2"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_id": os.Getenv("ALICLOUD_RESOURCE_MANAGER_FOLDER_ID2"),
					}),
				),
			},
		},
	})
}

var ResourceManagerAccountMap = map[string]string{
	"join_method":           CHECKSET,
	"join_time":             CHECKSET,
	"modify_time":           CHECKSET,
	"resource_directory_id": CHECKSET,
	"status":                CHECKSET,
	"type":                  CHECKSET,
}

func ResourceManagerAccountBasicdependence(name string) string {
	return ""
}

func TestUnitAlicloudResourceManagerAccount(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_resource_manager_account"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_resource_manager_account"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"display_name":        "CreateResourceAccountValue",
		"folder_id":           "CreateResourceAccountValue",
		"account_name_prefix": "CreateResourceAccountValue",
		"payer_account_id":    "CreateResourceAccountValue",
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
		// GetAccount
		"Account": map[string]interface{}{
			"AccountId":           "CreateResourceAccountValue",
			"DisplayName":         "CreateResourceAccountValue",
			"FolderId":            "CreateResourceAccountValue",
			"JoinMethod":          "CreateResourceAccountValue",
			"JoinTime":            "CreateResourceAccountValue",
			"ModifyTime":          "CreateResourceAccountValue",
			"ResourceDirectoryId": "CreateResourceAccountValue",
			"Status":              "CreateResourceAccountValue",
			"Type":                "CreateResourceAccountValue",
		},
		"PayerAccountId": "CreateResourceAccountValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateResourceAccount
		"Account": map[string]interface{}{
			"AccountId": "CreateResourceAccountValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_resource_manager_account", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewResourcemanagerClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudResourceManagerAccountCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetAccount Response
		"Account": map[string]interface{}{
			"AccountId": "CreateResourceAccountValue",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateResourceAccount" {
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
		err := resourceAlicloudResourceManagerAccountCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_resource_manager_account"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewResourcemanagerClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudResourceManagerAccountUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// MoveAccount
	attributesDiff := map[string]interface{}{
		"folder_id": "MoveAccountValue",
	}
	diff, err := newInstanceDiff("alicloud_resource_manager_account", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_resource_manager_account"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAccount Response
		"Account": map[string]interface{}{
			"FolderId": "MoveAccountValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "MoveAccount" {
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
		err := resourceAlicloudResourceManagerAccountUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_resource_manager_account"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UpdateAccount
	attributesDiff = map[string]interface{}{
		"display_name": "UpdateAccountValue",
	}
	diff, err = newInstanceDiff("alicloud_resource_manager_account", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_resource_manager_account"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAccount Response
		"Account": map[string]interface{}{
			"DisplayName": "UpdateAccountValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateAccount" {
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
		err := resourceAlicloudResourceManagerAccountUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_resource_manager_account"].Schema).Data(dExisted.State(), nil)
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
			if *action == "GetAccount" {
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
		err := resourceAlicloudResourceManagerAccountRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudResourceManagerAccountDelete(dExisted, rawClient)
	assert.Nil(t, err)

}
