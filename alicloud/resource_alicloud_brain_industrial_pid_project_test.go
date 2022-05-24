package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_brain_industrial_pid_project",
		&resource.Sweeper{
			Name: "alicloud_brain_industrial_pid_project",
			F:    testSweepBrainIndustrialPidProject,
		})
}

func testSweepBrainIndustrialPidProject(region string) error {
	if !testSweepPreCheckWithRegions(region, false, connectivity.BrainIndustrialRegions) {
		log.Printf("[INFO] Skipping Brain Industrial unsupported region: %s", region)
		return nil
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"CurrentPage": 1,
		"PageSize":    20,
	}
	var response map[string]interface{}
	action := "ListPidProjects"
	conn, err := client.NewAistudioClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, _ = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, request, &runtime)
		if fmt.Sprintf(`%v`, response["Code"]) != "200" {
			log.Println(fmt.Errorf("%s failed: %v", action, response))
			return nil
		}
		resp, err := jsonpath.Get("$.PidProjectList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PidProjectList", response)
		}
		for _, v := range resp.([]interface{}) {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["PidProjectName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Brain Industrial: %s", item["PidProjectName"].(string))
				continue
			}
			actionDelete := "DeletePidProject"
			requestDelete := map[string]interface{}{
				"PidProjectId": item["PidProjectId"],
			}
			response, err = conn.DoRequest(StringPointer(actionDelete), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, requestDelete, &util.RuntimeOptions{})
			if fmt.Sprintf(`%v`, response["Code"]) == "200" {
				log.Printf("[INFO] Delete Brain Industrial Project success: %s ", item["PidProjectName"].(string))
			} else if fmt.Sprintf(`%v`, response["Code"]) == "-100" && strings.Contains(response["Message"].(string), "存在回路") {
				log.Printf("[INFO] Firstly, Delete Loop belongs to Project")
				actionLoopList := "ListPidLoops"
				requestLoopList := map[string]interface{}{
					"PidProjectId": item["PidProjectId"],
					"PageSize":     20,
					"CurrentPage":  1,
				}
				for {
					runtime := util.RuntimeOptions{}
					runtime.SetAutoretry(true)

					responseLoop, _ := conn.DoRequest(StringPointer(actionLoopList), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, requestLoopList, &util.RuntimeOptions{})
					respLoop, _ := jsonpath.Get("$.PidLoopList", responseLoop)

					for _, v := range respLoop.([]interface{}) {
						itemLoop := v.(map[string]interface{})
						actionLoopDelete := "DeletePidLoop"
						requestLoopDelete := map[string]interface{}{
							"PidLoopId": itemLoop["PidLoopId"],
						}
						responseLoopDelete, _ := conn.DoRequest(StringPointer(actionLoopDelete), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, requestLoopDelete, &util.RuntimeOptions{})
						if fmt.Sprintf(`%v`, responseLoopDelete["Code"]) != "200" {
							log.Printf("[ERROR] Failed to delete Brain Industrial Loop (%s): %s", itemLoop["PidLoopId"].(string), responseLoopDelete["Message"].(string))
						} else {
							log.Printf("[INFO] Delete Brain Industrial Loop success (%s): %s", itemLoop["PidLoopId"].(string), responseLoopDelete["Message"].(string))
						}
					}
					if len(respLoop.([]interface{})) < request["PageSize"].(int) {
						break
					}
					request["CurrentPage"] = request["CurrentPage"].(int) + 1
				}
				log.Printf("[INFO] Delete Loop Done, Then delete Brain Industrial Project again")
				responseAgain, _ := conn.DoRequest(StringPointer(actionDelete), nil, StringPointer("POST"), StringPointer("2020-09-20"), StringPointer("AK"), nil, requestDelete, &util.RuntimeOptions{})
				if fmt.Sprintf(`%v`, responseAgain["Code"]) != "200" {
					log.Printf("[ERROR] Failed to again delete Brain Industrial Project  (%s): %s", item["PidProjectName"].(string), responseAgain["Message"].(string))
				} else {
					log.Printf("[INFO] Delete Brain Industrial Project again Success(%s): %s", item["PidProjectName"].(string), responseAgain["Message"].(string))
				}
			} else if fmt.Sprintf(`%v`, response["Code"]) != "200" {
				log.Printf("[ERROR] Failed to delete Brain Industrial Project (%s): %s", item["PidProjectName"].(string), response["Message"].(string))
			}
		}
		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	return nil
}

func TestAccAlicloudBrainIndustrialPidProject_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_brain_industrial_pid_project.default"
	ra := resourceAttrInit(resourceId, AlicloudBrainIndustrialPidProjectMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Brain_industrialService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBrainIndustrialPidProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBrainIndustrialPidProjectBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organization_id": "${alicloud_brain_industrial_pid_organization.default.id}",
					"pid_project_desc":    "tf test",
					"pid_project_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_id": CHECKSET,
						"pid_project_desc":    "tf test",
						"pid_project_name":    name,
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
					"pid_organization_id": "${alicloud_brain_industrial_pid_organization.update.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_project_desc": "tf test update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_project_desc": "tf test update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_project_name": "tf-testAccUp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_project_name": "tf-testAccUp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organization_id": "${alicloud_brain_industrial_pid_organization.default.id}",
					"pid_project_desc":    "tf test",
					"pid_project_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_id": CHECKSET,
						"pid_project_desc":    "tf test",
						"pid_project_name":    name,
					}),
				),
			},
		},
	})
}

var AlicloudBrainIndustrialPidProjectMap = map[string]string{}

func AlicloudBrainIndustrialPidProjectBasicDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_brain_industrial_pid_organization" "default" {
		pid_organization_name = "%s"
	}
	resource "alicloud_brain_industrial_pid_organization" "update" {
		pid_organization_name = "tf-testAccUp"
	}`, name)
}

func TestUnitAlicloudBrainIndustrialPidProject(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_brain_industrial_pid_project"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_brain_industrial_pid_project"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"pid_organization_id": "CreatePidProjectValue",
		"pid_project_name":    "CreatePidProjectValue",
		"pid_project_desc":    "CreatePidProjectValue",
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
		// ListPidProjects
		"PidProjectList": []interface{}{
			map[string]interface{}{
				"PidOrganisationId": "CreatePidProjectValue",
				"PidProjectDesc":    "CreatePidProjectValue",
				"PidProjectName":    "CreatePidProjectValue",
				"PidProjectId":      "CreatePidProjectValue",
			},
		},
		"PidProjectId": "CreatePidProjectValue",
		"Message":      "CreatePidProjectValue",
		"Code":         "200",
	}
	CreateMockResponse := map[string]interface{}{
		// CreatePidProject
		"PidProjectId": "CreatePidProjectValue",
		"Message":      "CreatePidProjectValue",
		"Code":         "200",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_brain_industrial_pid_project", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAistudioClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudBrainIndustrialPidProjectCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListPidProjects Response
		"PidProjectId": "CreatePidProjectValue",
		"Message":      "CreatePidProjectValue",
		"Code":         "200",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreatePidProject" {
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
		err := resourceAlicloudBrainIndustrialPidProjectCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_brain_industrial_pid_project"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAistudioClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudBrainIndustrialPidProjectUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdatePidProject
	attributesDiff := map[string]interface{}{
		"pid_organization_id": "UpdatePidProjectValue",
		"pid_project_name":    "UpdatePidProjectValue",
		"pid_project_desc":    "UpdatePidProjectValue",
	}
	diff, err := newInstanceDiff("alicloud_brain_industrial_pid_project", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_brain_industrial_pid_project"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListPidProjects Response
		"PidProjectList": []interface{}{
			map[string]interface{}{
				"PidOrganisationId": "UpdatePidProjectValue",
				"PidProjectDesc":    "UpdatePidProjectValue",
				"PidProjectName":    "UpdatePidProjectValue",
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdatePidProject" {
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
		err := resourceAlicloudBrainIndustrialPidProjectUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_brain_industrial_pid_project"].Schema).Data(dExisted.State(), nil)
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
			if *action == "ListPidProjects" {
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
		err := resourceAlicloudBrainIndustrialPidProjectRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAistudioClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudBrainIndustrialPidProjectDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeletePidProject" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Code": "200",
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudBrainIndustrialPidProjectDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
