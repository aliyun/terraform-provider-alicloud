// Package ehpc. Internal test package (package ehpc) so the TestUnit
// gomonkey tests can call the unexported resourceAlicloudEhpc* CRUD
// functions directly, matching the migrated-subpackage convention
// (see alicloud/service/ecs/disk_unit_test.go).
package ehpc

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAlicloudEhpcJobTemplate(t *testing.T) {
	p := unitTestProviderResources()
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
	rawClient, err := sweep.SharedClientForRegion(region)
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
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_ehpc_job_template", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEhpcClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
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
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
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
