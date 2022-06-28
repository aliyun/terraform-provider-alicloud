package alicloud

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKMSKeyVersion_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key_version.default"
	ra := resourceAttrInit(resourceId, KmsKeyVersionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKeyVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", resourceKMSKeyVersionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckKMSForKeyIdImport(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_id": os.Getenv("ALICLOUD_KMS_KEY_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckKmsKeyVersionExists(resourceId, &l),
					testAccCheck(KmsKeyVersionMap),
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

var KmsKeyVersionMap = map[string]string{
	"key_version_id": CHECKSET,
	"key_id":         CHECKSET,
}

//func testAccCheckKmsKeyVersionExists(n string, kv *kms.KeyVersion) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return WrapError(fmt.Errorf("Not found: %s", n))
//		}
//
//		if rs.Primary.ID == "" {
//			return WrapError(Error("No Key Version ID is set"))
//		}
//
//		client := testAccProvider.Meta().(*connectivity.AliyunClient)
//
//		request := kms.CreateListKeyVersionsRequest()
//		request.KeyId = rs.Primary.Attributes["key_id"]
//
//		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
//			return kmsClient.ListKeyVersions(request)
//		})
//
//		if err == nil {
//			response, _ := raw.(*kms.ListKeyVersionsResponse)
//			if len(response.KeyVersions.KeyVersion) > 0 {
//				for _, v := range response.KeyVersions.KeyVersion {
//					if v.KeyVersionId == strings.Split(rs.Primary.ID, ":")[1] {
//						*kv = v
//						return nil
//					}
//				}
//			}
//			return WrapError(fmt.Errorf("Error finding key version %s", rs.Primary.ID))
//		}
//		return WrapError(err)
//	}
//}

func resourceKMSKeyVersionConfigDependence(name string) string {
	return ""
}

func TestUnitAlicloudKMSKeyVersion(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_kms_key_version"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_kms_key_version"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"key_id": "CreateKeyVersionValue",
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
		// DescribeKeyVersion
		"KeyVersion": map[string]interface{}{
			"KeyId":        "CreateKeyVersionValue",
			"KeyVersionId": "CreateKeyVersionValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateKeyVersion
		"KeyVersion": map[string]interface{}{
			"KeyId":        "CreateKeyVersionValue",
			"KeyVersionId": "CreateKeyVersionValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_kms_key_version", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudKmsKeyVersionCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeKeyVersion Response
		"KeyVersion": map[string]interface{}{
			"KeyId":        "CreateKeyVersionValue",
			"KeyVersionId": "CreateKeyVersionValue",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateKeyVersion" {
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
		err := resourceAlicloudKmsKeyVersionCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_kms_key_version"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	attributesDiff := map[string]interface{}{}
	diff, err := newInstanceDiff("alicloud_kms_key_version", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_kms_key_version"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeKeyVersion" {
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
		err := resourceAlicloudKmsKeyVersionRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudKmsKeyVersionDelete(dExisted, rawClient)
	assert.Nil(t, err)

}
