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
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func testSweepCloudSSODirectoryAccessAssignment(region, directoryId string) error {
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
	action := "ListAccessAssignments"
	request := map[string]interface{}{}
	request["DirectoryId"] = directoryId

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

	resp, err := jsonpath.Get("$.AccessAssignments", response)
	if formatInt(response["TotalCounts"]) != 0 && err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.AccessAssignments", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["AccessConfigurationName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Cloud SSO AccessConfigurationName: %s", item["AccessConfigurationName"].(string))
			continue
		}
		action := "DeleteAccessAssignment"
		req := map[string]interface{}{
			"DirectoryId":           directoryId,
			"AccessConfigurationId": item["AccessConfigurationId"],
			"TargetType":            item["TargetType"],
			"TargetId":              item["TargetId"],
			"PrincipalType":         item["PrincipalType"],
			"PrincipalId":           item["PrincipalId"],
		}
		_, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, req, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cloud SSO AccessAssignment (%s): %s", item["AccessConfigurationName"].(string), err)
		}
		log.Printf("[INFO] Delete Cloud SSO AccessAssignment success: %s ", item["AccessConfigurationName"].(string))
	}
	return nil
}

func TestUnitAliCloudCloudSSOAccessAssignment(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_cloud_sso_access_assignment"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cloud_sso_access_assignment"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"directory_id":            "CreateAccessAssignmentValue",
		"access_configuration_id": "CreateAccessAssignmentValue",
		"principal_id":            "CreateAccessAssignmentValue",
		"principal_type":          "CreateAccessAssignmentValue",
		"target_id":               "CreateAccessAssignmentValue",
		"target_type":             "CreateAccessAssignmentValue",
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
		// ListAccessAssignments
		"AccessAssignments": []interface{}{
			map[string]interface{}{
				"PrincipalId":           "CreateAccessAssignmentValue",
				"DirectoryId":           "CreateAccessAssignmentValue",
				"AccessConfigurationId": "CreateAccessAssignmentValue",
				"TargetType":            "CreateAccessAssignmentValue",
				"TargetId":              "CreateAccessAssignmentValue",
				"PrincipalType":         "CreateAccessAssignmentValue",
			},
		},
		"Task": map[string]interface{}{
			"TaskId": "CreateAccessAssignmentValue",
		},
		"TaskStatus": map[string]interface{}{
			"Status": "Success",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateAccessAssignment
		"Task": map[string]interface{}{
			"TaskId": "CreateAccessAssignmentValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_cloud_sso_access_assignment", errorCode))
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
		err := resourceAliCloudCloudSSOAccessAssignmentCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{
			// ListAccessAssignments Response
			"Task": map[string]interface{}{
				"TaskId": "CreateAccessAssignmentValue",
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CreateAccessAssignment" {
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
			err := resourceAliCloudCloudSSOAccessAssignmentCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cloud_sso_access_assignment"].Schema).Data(dInit.State(), nil)
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
		err := resourceAliCloudCloudSSOAccessAssignmentUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	// Read
	t.Run("Read", func(t *testing.T) {
		attributesDiff := map[string]interface{}{}
		diff, err := newInstanceDiff("alicloud_cloud_sso_access_assignment", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_sso_access_assignment"].Schema).Data(dInit.State(), diff)
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "{}"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "ListAccessAssignments" {
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
			err := resourceAliCloudCloudSSOAccessAssignmentRead(dExisted, rawClient)
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
		err := resourceAliCloudCloudSSOAccessAssignmentDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		attributesDiff := map[string]interface{}{}
		diff, err := newInstanceDiff("alicloud_cloud_sso_access_assignment", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_sso_access_assignment"].Schema).Data(dInit.State(), diff)
		errorCodes := []string{"NonRetryableError", "Throttling", "OperationConflict.Task", "DeletionConflict.AccessConfigurationProvisioning.AccessAssignment", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteAccessAssignment" {
					switch errorCode {
					case "NonRetryableError":
						return failedResponseMock(errorCode)
					default:
						retryIndex++
						if errorCodes[retryIndex] == "nil" {
							ReadMockResponse = map[string]interface{}{
								"Task": map[string]interface{}{
									"TaskId": "CreateAccessAssignmentValue",
								},
								"TaskStatus": map[string]interface{}{
									"Status": "Success",
								},
							}
							return ReadMockResponse, nil
						}
						return failedResponseMock(errorCodes[retryIndex])
					}
				}
				return ReadMockResponse, nil
			})
			err := resourceAliCloudCloudSSOAccessAssignmentDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "nil":
				assert.Nil(t, err)
			}

		}
	})
}

// Test CloudSSO AccessAssignment. >>> Resource test cases, automatically generated.
// Case AccessAssignment 10051
func TestAccAliCloudCloudSSOAccessAssignment_basic10051(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOAccessAssignmentMap10051)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSSOAccessAssignment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOAccessAssignmentBasicDependence10051)
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
					"directory_id":            "${local.directory_id}",
					"principal_id":            "${alicloud_cloud_sso_user.default.user_id}",
					"target_type":             "RD-Account",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"principal_type":          "User",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            CHECKSET,
						"access_configuration_id": CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               CHECKSET,
						"principal_type":          "User",
						"principal_id":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

func TestAccAliCloudCloudSSOAccessAssignment_basic10051_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOAccessAssignmentMap10051)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSSOAccessAssignment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOAccessAssignmentBasicDependence10051)
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
					"directory_id":            "${local.directory_id}",
					"principal_id":            "${alicloud_cloud_sso_user.default.user_id}",
					"target_type":             "RD-Account",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"principal_type":          "User",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
					"deprovision_strategy":    "DeprovisionForLastAccessAssignmentOnAccount",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            CHECKSET,
						"access_configuration_id": CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               CHECKSET,
						"principal_type":          "User",
						"principal_id":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

var AliCloudCloudSSOAccessAssignmentMap10051 = map[string]string{
	"create_time": CHECKSET,
}

func AliCloudCloudSSOAccessAssignmentBasicDependence10051(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_directories" "default" {
	}

	data "alicloud_cloud_sso_directories" "default" {
	}

	resource "alicloud_cloud_sso_directory" "default" {
  		count          = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  		directory_name = var.name
	}

	resource "alicloud_cloud_sso_group" "default" {
  		directory_id = local.directory_id
  		group_name   = var.name
  		description  = var.name
	}

	resource "alicloud_cloud_sso_user" "default" {
  		directory_id = local.directory_id
  		user_name    = var.name
	}

	resource "alicloud_cloud_sso_access_configuration" "default" {
  		directory_id              = local.directory_id
  		access_configuration_name = var.name
	}

	locals {
  		directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
	}
`, name)
}

// Case AccessAssignment_online 10355

func TestAccAliCloudCloudSSOAccessAssignment_basic10355(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOAccessAssignmentMap10051)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSSOAccessAssignment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOAccessAssignmentBasicDependence10051)
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
					"directory_id":            "${local.directory_id}",
					"principal_id":            "${alicloud_cloud_sso_group.default.group_id}",
					"target_type":             "RD-Account",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"principal_type":          "Group",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            CHECKSET,
						"access_configuration_id": CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               CHECKSET,
						"principal_type":          "Group",
						"principal_id":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

func TestAccAliCloudCloudSSOAccessAssignment_basic10355_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudSSOAccessAssignmentMap10051)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudSSOServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSSOAccessAssignment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudSSOAccessAssignmentBasicDependence10051)
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
					"directory_id":            "${local.directory_id}",
					"principal_id":            "${alicloud_cloud_sso_group.default.group_id}",
					"target_type":             "RD-Account",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"principal_type":          "Group",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
					"deprovision_strategy":    "DeprovisionForLastAccessAssignmentOnAccount",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            CHECKSET,
						"access_configuration_id": CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               CHECKSET,
						"principal_type":          "Group",
						"principal_id":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

// Test CloudSSO AccessAssignment. <<< Resource test cases, automatically generated.
