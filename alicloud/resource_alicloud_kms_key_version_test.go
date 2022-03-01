package alicloud

import (
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKmsKeyVersion_basic(t *testing.T) {
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

func TestAccAlicloudKmsKeyVersion_unit(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_kms_key_version"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_kms_key_version"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"key_id": "key_id",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"KeyVersion": map[string]interface{}{
			"KeyVersionId": "MockKeyVersionId",
			"KeyId":        "key_id",
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_kms_key_version", "MockKeyVersionId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:    String(errorCode),
				Data:    String(errorCode),
				Message: String(errorCode),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:    String("loadEndpoint error"),
				Data:    String("loadEndpoint error"),
				Message: String("loadEndpoint error"),
			}
		})
		err := resourceAlicloudKmsKeyVersionCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudKmsKeyVersionCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAlicloudKmsKeyVersionCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId(fmt.Sprint("key_id", ":", "MockKeyVersionId"))

	// Delete
	t.Run("DeleteClientNormal", func(t *testing.T) {
		err := resourceAlicloudKmsKeyVersionDelete(d, rawClient)
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeKmsKeyVersionNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudKmsKeyVersionRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeKmsKeyVersionAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAlicloudKmsKeyVersionRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}
