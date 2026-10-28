package cloudsso

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
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	sweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestUnitAliCloudCloudSSOSCIMServerCredential mock-driven unit test; lives in the internal
// cloudsso package so it can call the unexported CRUD funcs.
func TestUnitAliCloudCloudSSOSCIMServerCredential(t *testing.T) {
	p := unitTestProviderResources()
	dInit, _ := schema.InternalMap(p["alicloud_cloud_sso_scim_server_credential"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_cloud_sso_scim_server_credential"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"directory_id": "CreateSCIMServerCredentialValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value) // lintignore: R001
		assert.Nil(t, err)
		err = dExisted.Set(key, value) // lintignore: R001
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
		// ListSCIMServerCredentials
		"SCIMServerCredentials": []interface{}{
			map[string]interface{}{
				"CredentialId": "CreateSCIMServerCredentialValue",
				"Status":       "Disabled",
				"DirectoryId":  "CreateSCIMServerCredentialValue",
			},
		},
		"SCIMServerCredential": map[string]interface{}{
			"CredentialId": "CreateSCIMServerCredentialValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateSCIMServerCredential
		"SCIMServerCredential": map[string]interface{}{
			"CredentialId": "CreateSCIMServerCredentialValue",
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
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_cloud_sso_scim_server_credential", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	t.Run("Create", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewCloudssoClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCloudSsoScimServerCredentialCreate(dInit, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		ReadMockResponseDiff = map[string]interface{}{
			"SCIMServerCredential": map[string]interface{}{
				"CredentialId": "CreateSCIMServerCredentialValue",
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1 // a counter used to cover retry scenario; the same below
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "CreateSCIMServerCredential" {
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
			err := resourceAliCloudCloudSsoScimServerCredentialCreate(dInit, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cloud_sso_scim_server_credential"].Schema).Data(dInit.State(), nil)
				for key, value := range attributes {
					dCompare.Set(key, value) // lintignore: R001
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
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCloudSsoScimServerCredentialUpdate(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		// UpdateSCIMServerCredentialStatus
		attributesDiff := map[string]interface{}{
			"status": "Enabled",
		}
		diff, err := newInstanceDiff("alicloud_cloud_sso_scim_server_credential", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_sso_scim_server_credential"].Schema).Data(dInit.State(), diff)
		ReadMockResponseDiff = map[string]interface{}{
			"SCIMServerCredentials": []interface{}{
				map[string]interface{}{
					"Status": "Enabled",
				},
			},
		}
		errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "UpdateSCIMServerCredentialStatus" {
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
			err := resourceAliCloudCloudSsoScimServerCredentialUpdate(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
				dCompare, _ := schema.InternalMap(p["alicloud_cloud_sso_scim_server_credential"].Schema).Data(dExisted.State(), nil)
				for key, value := range attributes {
					dCompare.Set(key, value) // lintignore: R001
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
				if *action == "ListSCIMServerCredentials" {
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
			err := resourceAliCloudCloudSsoScimServerCredentialRead(dExisted, rawClient)
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
				Code:       helper.String("loadEndpoint error"),
				Data:       helper.String("loadEndpoint error"),
				Message:    helper.String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudCloudSsoScimServerCredentialDelete(dExisted, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
		attributesDiff := map[string]interface{}{}
		diff, err := newInstanceDiff("alicloud_cloud_sso_scim_server_credential", attributes, attributesDiff, dInit.State())
		if err != nil {
			t.Error(err)
		}
		dExisted, _ = schema.InternalMap(p["alicloud_cloud_sso_scim_server_credential"].Schema).Data(dInit.State(), diff)
		errorCodes := []string{"NonRetryableError", "Throttling", "nil", "EntityNotExists.SCIMCredential"}
		for index, errorCode := range errorCodes {
			retryIndex := index - 1
			patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "DeleteSCIMServerCredential" {
					switch errorCode {
					case "NonRetryableError", "EntityNotExists.SCIMCredential":
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
			err := resourceAliCloudCloudSsoScimServerCredentialDelete(dExisted, rawClient)
			patches.Reset()
			switch errorCode {
			case "NonRetryableError":
				assert.NotNil(t, err)
			case "EntityNotExists.SCIMCredential":
				assert.Nil(t, err)
			}

		}
	})
}
