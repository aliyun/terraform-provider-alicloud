package resourcemanager_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_resource_manager_handshake", &resource.Sweeper{
		Name: "alicloud_resource_manager_handshake",
		F:    testSweepResourceManagerHandshake,
	})
}

func testSweepResourceManagerHandshake(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	action := "ListHandshakesForResourceDirectory"
	request := make(map[string]interface{})
	request["PageSize"] = helper.PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	var handshakeIds []string

	for {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"EntityNotExists.ResourceDirectory"}) {
				return nil
			}
			log.Printf("[ERROR] Failed to retrieve resoure manager handshake in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Handshakes.Handshake", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Handshakes.Handshake", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			// Skip Invalid handshake.
			if v, ok := item["Status"].(string); ok && v == "Pending" {
				handshakeIds = append(handshakeIds, item["HandshakeId"].(string))
			}
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, handshakeId := range handshakeIds {
		log.Printf("[INFO] Delete resource manager handshake %s ", handshakeId)

		request := map[string]interface{}{
			"HandshakeId": handshakeId,
		}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
			log.Printf("[ERROR] Failed to delete resource manager handshake (%s): %s", handshakeId, err)
		}
	}
	return nil
}

func TestAccAliCloudResourceManagerHandshake_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_handshake.default"
	ra := acctest.ResourceAttrInit(resourceId, ResourceManagerHandshakeMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourcemanagerService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerHandshake")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerHandshake%d", rand)
	testAccUseResourceManagerHandshakeManagementAccount(t)
	invitedAccountId := testAccResolveResourceManagerHandshakeInvitedAccountId(t)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, ResourceManagerHandshakeBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
			testAccPreCheckResourceManagerHandshakeAccountDetached(t, invitedAccountId)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"target_entity": invitedAccountId,
					"target_type":   "Account",
					"note":          "test resource manager handshake",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_entity": acctest.CHECKSET,
						"target_type":   "Account",
						"note":          "test resource manager handshake",
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

var ResourceManagerHandshakeMap = map[string]string{}

func ResourceManagerHandshakeBasicdependence(name string) string {
	return resourceManagerHandshakeManagementProviderDependence(name)
}

// lintignore: R001
func TestUnitAlicloudResourceManagerHandshake(t *testing.T) {
	p := acctest.Provider.ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_resource_manager_handshake"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_resource_manager_handshake"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"note":          "CreateHandshakeValue",
		"target_entity": "CreateHandshakeValue",
		"target_type":   "CreateHandshakeValue",
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
		// GetFolder
		"Handshake": map[string]interface{}{
			"HandshakeId":         "MockHandshakeId",
			"ExpireTime":          "CreateHandshakeValue",
			"MasterAccountId":     "CreateHandshakeValue",
			"MasterAccountName":   "CreateHandshakeValue",
			"ModifyTime":          "CreateHandshakeValue",
			"Note":                "CreateHandshakeValue",
			"ResourceDirectoryId": "CreateHandshakeValue",
			"Status":              "CreateHandshakeValue",
			"TargetEntity":        "CreateHandshakeValue",
			"TargetType":          "CreateHandshakeValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateFolder
		"Handshake": map[string]interface{}{
			"HandshakeId": "MockHandshakeId",
		},
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_resource_manager_handshake", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewResourcemanagerClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    helper.String("loadEndpoint error"),
				Data:    helper.String("loadEndpoint error"),
				Message: helper.String("loadEndpoint error"),
			}
		})
		err := resourcemanager.ResourceHandshakeCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{}
		errorCodes := []string{"NonRetryableError", "LimitExceeded.InvitationRate", "LimitExceeded.SameTargetInvitationRate", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "InviteAccountToResourceDirectory" {
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
			err := resourcemanager.ResourceHandshakeCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_resource_manager_handshake"].Schema).Data(dInit.State(), nil)
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

	// Read
	t.Run("Read", func(t *testing.T) {
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "{}"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "GetHandshake" {
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
			err := resourcemanager.ResourceHandshakeRead(dExisted, rawClient)
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
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewResourcemanagerClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    helper.String("loadEndpoint error"),
				Data:    helper.String("loadEndpoint error"),
				Message: helper.String("loadEndpoint error"),
			}
		})
		err := resourcemanager.ResourceHandshakeDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "EntityNotExists.Handshake", "HandshakeStatusMismatch"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CancelHandshake" {
					switch errorCode {
					case "NonRetryableError", "EntityNotExists.Handshake", "HandshakeStatusMismatch":
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
			err := resourcemanager.ResourceHandshakeDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "EntityNotExists.Handshake", "HandshakeStatusMismatch":
				assert.Nil(t, err)
			}
		}
	})
}

// Test ResourceManager Handshake. >>> Resource test cases, automatically generated.
// Case Handshake-资源测试用例 11272
func TestAccAliCloudResourceManagerHandshake_basic11272(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_handshake.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudResourceManagerHandshakeMap11272)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerHandshake")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccUseResourceManagerHandshakeManagementAccount(t)
	invitedAccountId := testAccResolveResourceManagerHandshakeInvitedAccountId(t)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudResourceManagerHandshakeBasicDependence11272)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
			testAccPreCheckResourceManagerHandshakeAccountDetached(t, invitedAccountId)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"note":          "test",
					"target_type":   "Account",
					"target_entity": invitedAccountId,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"note":          "test",
						"target_type":   "Account",
						"target_entity": acctest.CHECKSET,
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

var AlicloudResourceManagerHandshakeMap11272 = map[string]string{
	"status":                acctest.CHECKSET,
	"modify_time":           acctest.CHECKSET,
	"resource_directory_id": acctest.CHECKSET,
	"create_time":           acctest.CHECKSET,
	"master_account_name":   acctest.CHECKSET,
	"master_account_id":     acctest.CHECKSET,
	"expire_time":           acctest.CHECKSET,
}

func AlicloudResourceManagerHandshakeBasicDependence11272(name string) string {
	return resourceManagerHandshakeManagementProviderDependence(name)
}

func resourceManagerHandshakeManagementProviderDependence(name string) string {
	ak, sk := testAccResourceManagerHandshakeManagementCredentialValues()
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  region     = "cn-hangzhou"
  access_key = "%s"
  secret_key = "%s"
}
`, name, ak, sk)
}

// Test ResourceManager Handshake. <<< Resource test cases, automatically generated.
