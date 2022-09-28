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

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEhpcJobTemplate_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ehpc_job_template.default"
	ra := resourceAttrInit(resourceId, AlicloudEhpcJobTemplateMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EhpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEhpcJobTemplate")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sehpcjobtemplate%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEhpcJobTemplateBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"job_template_name": "JobTemplateNameT",
					"command_line":      "./LammpsTest/lammps.pbs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_template_name": "JobTemplateNameT",
						"command_line":      "./LammpsTest/lammps.pbs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"priority": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stderr_redirect_path": "./LammpsTest1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stderr_redirect_path": "./LammpsTest1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"clock_time": "12:00:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"clock_time": "12:00:00",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gpu": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gpu": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"runas_user": "user1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"runas_user": "user1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"thread": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"thread": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"job_template_name": "JobTemplateNameH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"job_template_name": "JobTemplateNameH",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"variables": "[{Name:,Value:},{Name:,Value:}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"variables": "[{Name:,Value:},{Name:,Value:}]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"re_runable": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"re_runable": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"command_line": "./LammpsTestOne/lammps.pbs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_line": "./LammpsTestOne/lammps.pbs",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mem": "1GB",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mem": "1GB",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stdout_redirect_path": "./LammpsTest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stdout_redirect_path": "./LammpsTest",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"array_request": "1-10:2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"array_request": "1-10:2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"queue": "workq",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"queue": "workq",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"package_path": "./jobfolderOne",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_path": "./jobfolderOne",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task":                 "4",
					"priority":             "2",
					"stderr_redirect_path": "./LammpsTestT",
					"node":                 "4",
					"clock_time":           "14:00:00",
					"gpu":                  "3",
					"runas_user":           "user3",
					"thread":               "3",
					"job_template_name":    "JobTemplateNameY",
					"variables":            "[{Demo:,Test:},{Test:,Demo:}]",
					"re_runable":           "true",
					"command_line":         "./LammpsTestT/lammps.pbs",
					"mem":                  "3GB",
					"stdout_redirect_path": "./LammpsTestH",
					"array_request":        "1-12:2",
					"queue":                "workq",
					"package_path":         "./jobfolderT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task":                 "4",
						"priority":             "2",
						"stderr_redirect_path": "./LammpsTestT",
						"node":                 "4",
						"clock_time":           "14:00:00",
						"gpu":                  "3",
						"runas_user":           "user3",
						"thread":               "3",
						"job_template_name":    "JobTemplateNameY",
						"variables":            "[{Demo:,Test:},{Test:,Demo:}]",
						"re_runable":           "true",
						"command_line":         "./LammpsTestT/lammps.pbs",
						"mem":                  "3GB",
						"stdout_redirect_path": "./LammpsTestH",
						"array_request":        "1-12:2",
						"queue":                "workq",
						"package_path":         "./jobfolderT",
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

var AlicloudEhpcJobTemplateMap0 = map[string]string{
	"command_line":      "./LammpsTest/lammps.pbs",
	"job_template_name": "JobTemplateName",
}

func AlicloudEhpcJobTemplateBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

func TestUnitAlicloudEhpcJobTemplate(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ehpc_job_template"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ehpc_job_template"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"array_request":        "CreateJobTemplateValue",
		"clock_time":           "CreateJobTemplateValue",
		"command_line":         "CreateJobTemplateValue",
		"gpu":                  128,
		"job_template_name":    "CreateJobTemplateValue",
		"mem":                  "CreateJobTemplateValue",
		"node":                 1,
		"package_path":         "CreateJobTemplateValue",
		"priority":             1,
		"queue":                "CreateJobTemplateValue",
		"re_runable":           false,
		"runas_user":           "CreateJobTemplateValue",
		"stderr_redirect_path": "CreateJobTemplateValue",
		"stdout_redirect_path": "CreateJobTemplateValue",
		"task":                 1,
		"thread":               1,
		"variables":            "CreateJobTemplateValue",
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
		// ListJobTemplates
		"Templates": map[string]interface{}{
			"JobTemplates": []interface{}{
				map[string]interface{}{
					"ArrayRequest":       "CreateJobTemplateValue",
					"ClockTime":          "CreateJobTemplateValue",
					"CommandLine":        "CreateJobTemplateValue",
					"Gpu":                128,
					"Mem":                "CreateJobTemplateValue",
					"Node":               1,
					"PackagePath":        "CreateJobTemplateValue",
					"Priority":           1,
					"Queue":              "CreateJobTemplateValue",
					"Name":               "CreateJobTemplateValue",
					"ReRunable":          "false",
					"RunasUser":          "CreateJobTemplateValue",
					"StderrRedirectPath": "CreateJobTemplateValue",
					"StdoutRedirectPath": "CreateJobTemplateValue",
					"Task":               1,
					"Thread":             1,
					"Variables":          "CreateJobTemplateValue",
					"Id":                 "CreateJobTemplateValue",
				},
			},
		},
		"TemplateId": "CreateJobTemplateValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateJobTemplate
		"TemplateId": "CreateJobTemplateValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ehpc_job_template", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcJobTemplateCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// ListJobTemplates Response
		"Templates": map[string]interface{}{
			"JobTemplates": []interface{}{
				map[string]interface{}{
					"Id": "CreateJobTemplateValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateJobTemplate" {
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
		err := resourceAlicloudEhpcJobTemplateCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ehpc_job_template"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcJobTemplateUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// EditJobTemplate
	attributesDiff := map[string]interface{}{
		"array_request":        "EditJobTemplateValue",
		"clock_time":           "EditJobTemplateValue",
		"command_line":         "EditJobTemplateValue",
		"gpu":                  256,
		"job_template_name":    "EditJobTemplateValue",
		"mem":                  "EditJobTemplateValue",
		"node":                 2,
		"package_path":         "EditJobTemplateValue",
		"priority":             2,
		"queue":                "EditJobTemplateValue",
		"re_runable":           true,
		"runas_user":           "EditJobTemplateValue",
		"stderr_redirect_path": "EditJobTemplateValue",
		"stdout_redirect_path": "EditJobTemplateValue",
		"task":                 2,
		"thread":               2,
		"variables":            "EditJobTemplateValue",
	}
	diff, err := newInstanceDiff("alicloud_ehpc_job_template", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ehpc_job_template"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// ListJobTemplates Response
		"Templates": map[string]interface{}{
			"JobTemplates": []interface{}{
				map[string]interface{}{
					"ArrayRequest":       "EditJobTemplateValue",
					"ClockTime":          "EditJobTemplateValue",
					"CommandLine":        "EditJobTemplateValue",
					"Gpu":                256,
					"Mem":                "EditJobTemplateValue",
					"Node":               2,
					"PackagePath":        "EditJobTemplateValue",
					"Priority":           2,
					"Queue":              "EditJobTemplateValue",
					"Name":               "EditJobTemplateValue",
					"ReRunable":          "true",
					"RunasUser":          "EditJobTemplateValue",
					"StderrRedirectPath": "EditJobTemplateValue",
					"StdoutRedirectPath": "EditJobTemplateValue",
					"Task":               2,
					"Thread":             2,
					"Variables":          "EditJobTemplateValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "EditJobTemplate" {
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
		err := resourceAlicloudEhpcJobTemplateUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ehpc_job_template"].Schema).Data(dExisted.State(), nil)
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
			if *action == "ListJobTemplates" {
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
		err := resourceAlicloudEhpcJobTemplateRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudEhpcJobTemplateDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteJobTemplates" {
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
		err := resourceAlicloudEhpcJobTemplateDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}

}
