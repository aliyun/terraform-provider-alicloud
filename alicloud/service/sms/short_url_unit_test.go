package sms

import (
	"log"
	"os"
	"reflect"
	"strings"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAccAlicloudSmsShortUrl(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_sms_short_url"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_sms_short_url"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"short_url_name": "CreateSmsShortUrlValue",
		"source_url":     "CreateSmsShortUrlValue",
		"effective_days": 30,
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
	rawClient, err := unitTestSharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}

	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"Data": map[string]interface{}{
			"SourceUrl":          "CreateSmsShortUrlValue",
			"ShortUrlName":       "CreateSmsShortUrlValue",
			"ShortUrl":           "CreateSmsShortUrlValue",
			"ExpireDate":         "2022-08-19 00:00:00",
			"PageViewCount":      "0",
			"CreateDate":         "2022-07-20 10:03:51",
			"ShortUrlStatus":     "audit",
			"UniqueVisitorCount": "0",
		},
		"Code": "OK",
	}
	CreateMockResponse := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       helper.String(errorCode),
			Data:       helper.String(errorCode),
			Message:    helper.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_sms_short_url", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}
	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDysmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSmsShortUrlCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AddShortUrl" {
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
		err := resourceSmsShortUrlCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_sms_short_url"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_sms_short_url", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_sms_short_url"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "QueryShortUrl" {
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
		err := resourceSmsShortUrlRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDysmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("loadEndpoint error"),
			Data:       helper.String("loadEndpoint error"),
			Message:    helper.String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceSmsShortUrlDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_sms_short_url", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_sms_short_url"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteShortUrl" {
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
		err := resourceSmsShortUrlDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}

func unitTestProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"alicloud_sms_short_url": ResourceSmsShortUrl(),
	}
}

// unitTestSharedClientForRegion delegates to the shared sweeper client so the
// unit test reuses the same per-region endpoint-cache isolation as the root
// package instead of carrying a package-local copy of the client builder.
func unitTestSharedClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}

// newInstanceDiff mirrors the root helper for the migrated resource, so the
// unit test stays in the internal sms package (acctest import would
// create a cycle) while still exercising the SDKv2 diff machinery.
// lintignore: R001
func newInstanceDiff(resourceName string, attributes, attributesDiff map[string]interface{}, state *terraform.InstanceState) (*terraform.InstanceDiff, error) {
	p := unitTestProviderResources()
	dOld, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	dNew, _ := schema.InternalMap(p[resourceName].Schema).Data(state, nil)
	for key, value := range attributes {
		err := dOld.Set(key, value)
		if err != nil {
			return nil, errs.WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}
	for key, value := range attributesDiff {
		attributes[key] = value
	}

	for key, value := range attributes {
		err := dNew.Set(key, value)
		if err != nil {
			return nil, errs.WrapErrorf(err, "[ERROR] the field %s setting error.", key)
		}
	}

	diff := terraform.NewInstanceDiff()
	objectKey := ""
	for _, key := range helper.MapSort(dNew.State().Attributes) {
		newValue := dNew.State().Attributes[key]
		if objectKey != "" && !strings.HasPrefix(key, objectKey) {
			objectKey = ""
		}
		if objectKey == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					objectKey = strings.TrimSuffix(key, suffix)
					break
				}
			}
		}
		oldValue, ok := dOld.State().Attributes[key]
		if ok && oldValue == newValue {
			continue
		}
		if oldValue == "" {
			for _, suffix := range []string{"#", "%"} {
				if strings.HasSuffix(key, suffix) {
					oldValue = "0"
				}
			}
		}
		diff.Attributes[key] = &terraform.ResourceAttrDiff{Old: oldValue, New: newValue}

		if objectKey != "" {
			for removeKey, removeValue := range dOld.State().Attributes {
				if strings.HasPrefix(removeKey, objectKey) {
					if _, ok := dNew.State().Attributes[removeKey]; !ok {
						if len(strings.Split(removeKey, ".")) > 2 {
							delete(diff.Attributes, removeKey)
						} else {
							diff.Attributes[removeKey] = &terraform.ResourceAttrDiff{Old: removeValue, New: ""}
						}
					}
				}
			}
			objectKey = ""
		}
	}
	return diff, nil
}
