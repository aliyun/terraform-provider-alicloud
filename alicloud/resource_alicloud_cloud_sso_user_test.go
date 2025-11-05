package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func testSweepCloudSsoDirectoryUser(region, directoryId string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"",
	}
	action := "ListUsers"
	request := map[string]interface{}{}
	request["DirectoryId"] = directoryId
	request["MaxResults"] = 100

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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

	resp, err := jsonpath.Get("$.Users", response)
	if formatInt(response["TotalCounts"]) != 0 && err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Users", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["UserName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Cloud Sso User: %s", item["UserName"].(string))
			continue
		}
		action := "DeleteUser"
		req := map[string]interface{}{
			"UserId":      item["UserId"],
			"DirectoryId": directoryId,
		}
		_, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, req, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cloud Sso User (%s): %s", item["UserName"].(string), err)
		}
		log.Printf("[INFO] Delete Cloud Sso User success: %s ", item["UserName"].(string))
	}
	return nil
}

func TestUnitAliCloudCloudSSOUser(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cloud_sso_user"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cloud_sso_user"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"description":  "CreateUserValue",
		"display_name": "CreateUserValue",
		"directory_id": "CreateUserValue",
		"email":        "CreateUserValue",
		"first_name":   "CreateUserValue",
		"last_name":    "CreateUserValue",
		"status":       "Disabled",
		"user_name":    "CreateUserValue",
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
		// GetUser
		"User": map[string]interface{}{
			"UserId":      "CreateUserValue",
			"Description": "CreateUserValue",
			"DisplayName": "CreateUserValue",
			"Email":       "CreateUserValue",
			"DirectoryId": "CreateUserValue",
			"FirstName":   "CreateUserValue",
			"LastName":    "CreateUserValue",
			"Status":      "Disabled",
			"UserName":    "CreateUserValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateUser
		"User": map[string]interface{}{
			"UserId": "CreateUserValue",
		},
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cloud_sso_user", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudssoClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCloudSsoUserCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{
			// GetUser Response
			"User": map[string]interface{}{
				"UserId": "CreateUserValue",
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CreateUser" {
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
			err := resourceAliCloudCloudSsoUserCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cloud_sso_user"].Schema).Data(dInit.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Update
	t.Run("Update", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudssoClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCloudSsoUserUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		// UpdateUserStatus
		attributesDiff := map[string]interface{}{
			"status": "Enabled",
		}
		diff, err := newInstanceDiff("alicloud_cloud_sso_user", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_sso_user"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// GetUser Response
			"User": map[string]interface{}{
				"Status": "Enabled",
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "UpdateUserStatus" {
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
			err := resourceAliCloudCloudSsoUserUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cloud_sso_user"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}

		// UpdateUser
		attributesDiff = map[string]interface{}{
			"description":  "UpdateUserValue",
			"display_name": "UpdateUserValue",
			"email":        "UpdateUserValue",
			"first_name":   "UpdateUserValue",
			"last_name":    "UpdateUserValue",
		}
		diff, err = newInstanceDiff("alicloud_cloud_sso_user", attributes, attributesDiff, dExisted.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_sso_user"].Schema).Data(dExisted.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			// GetUser Response
			"User": map[string]interface{}{
				"Description": "UpdateUserValue",
				"DisplayName": "UpdateUserValue",
				"Email":       "UpdateUserValue",
				"FirstName":   "UpdateUserValue",
				"LastName":    "UpdateUserValue",
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
			err := resourceAliCloudCloudSsoUserUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cloud_sso_user"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					_ = dCompare.Set(key, value)
				}
				assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
			}
			if retryIndex >= len(errorCodes)-1 {
				break
			}
		}
	})

	// Read
	t.Run("Read", func(t *testing.T) {
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "{}"}
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
			err := resourceAliCloudCloudSsoUserRead(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "{}":
				assert.Nil(t, err)
			}
		}
	})

	// Delete
	t.Run("Delete", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudssoClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCloudSsoUserDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		attributesDiff := map[string]interface{}{}
		diff, err := newInstanceDiff("alicloud_cloud_sso_user", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_sso_user"].Schema).Data(dInit.State(), diff)
		errorCodes := []string{"NonRetryableError", "Throttling", "DeletionConflict.User.AccessAssigment", "nil", "EntityNotExists.User"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteUser" {
					switch errorCode {
					case "NonRetryableError", "EntityNotExists.User":
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
			err := resourceAliCloudCloudSsoUserDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "EntityNotExists.User":
				assert.Nil(t, err)
			}

		}
	})
}

// Test CloudSso User. >>> Resource test cases, automatically generated.
// Case User_pre_tag 10696
func TestAccAliCloudCloudSsoUser_basic10696(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_user.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSsoUserMap10696)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSsoUserBasicDependence10696)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id": "${local.directory_id}",
					"user_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id": CHECKSET,
						"user_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "cloud_sso_user@aliyun.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "cloud_sso_user@aliyun.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"first_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"first_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"last_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"last_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_authentication_settings": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_authentication_settings": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_authentication_settings": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_authentication_settings": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Yourpassword123!",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Yourpassword123!",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAliCloudCloudSsoUser_basic10696_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_user.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSsoUserMap10696)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSsoUserBasicDependence10696)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":                "${local.directory_id}",
					"user_name":                   name,
					"description":                 name,
					"display_name":                name,
					"email":                       "cloud_sso_user@aliyun.com",
					"first_name":                  name,
					"last_name":                   name,
					"mfa_authentication_settings": "Enabled",
					"password":                    "Yourpassword123!",
					"status":                      "Enabled",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":                CHECKSET,
						"user_name":                   name,
						"description":                 name,
						"display_name":                name,
						"email":                       "cloud_sso_user@aliyun.com",
						"first_name":                  name,
						"last_name":                   name,
						"mfa_authentication_settings": "Enabled",
						"password":                    "Yourpassword123!",
						"status":                      "Enabled",
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var AliCloudCloudSsoUserMap10696 = map[string]string{
	"create_time": CHECKSET,
	"user_id":     CHECKSET,
	"status":      CHECKSET,
}

func AliCloudCloudSsoUserBasicDependence10696(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_cloud_sso_directories" "default" {
	}

	resource "alicloud_cloud_sso_directory" "default" {
  		count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  		directory_name = var.name
	}

	locals {
  		directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
	}
`, name)
}

// Test CloudSso User. <<< Resource test cases, automatically generated.
