package alicloud

import (
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudEcsDiskEncryptionByDefault_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_encryption_by_default.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskEncryptionByDefaultMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDiskEncryptionByDefaultStatus")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", AliCloudEcsDiskEncryptionByDefaultBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEcsDiskEncryptionByDefaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
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

func testAccCheckEcsDiskEncryptionByDefaultDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_disk_encryption_by_default" {
			continue
		}

		// For this region-level setting resource, we check that encryption is disabled
		object, err := ecsServiceV2.DescribeEcsDiskEncryptionByDefaultStatus(rs.Primary.ID)
		if err != nil {
			// If we can't get the status, consider it as destroyed
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		// Check if encryption is disabled (this means the resource is in "destroyed" state)
		if encrypted, ok := object["Encrypted"].(bool); ok && !encrypted {
			continue // Resource is properly "destroyed" (encryption disabled)
		}

		// If encryption is still enabled, the resource wasn't properly destroyed
		return WrapError(Error("ECS disk encryption by default is still enabled in region %s", rs.Primary.ID))
	}

	return nil
}

var AliCloudEcsDiskEncryptionByDefaultMap = map[string]string{}

func AliCloudEcsDiskEncryptionByDefaultBasicDependence(name string) string {
	return ""
}

func TestUnitAliCloudEcsDiskEncryptionByDefault(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	var dExisted2Data *schema.ResourceData
	d, _ := schema.InternalMap(p["alicloud_ecs_disk_encryption_by_default"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ecs_disk_encryption_by_default"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	d.Set("enabled", false)
	dCreate.Set("enabled", true)
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeDiskEncryptionByDefaultStatus
		"Encrypted": false,
	}
	CreateMockResponse := map[string]interface{}{
		// EnableDiskEncryptionByDefault
	}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudEcsDiskEncryptionByDefaultCreate(dCreate, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDiskEncryptionByDefaultStatus Response
		"Encrypted": true,
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; -1 means no retry scenario; 0 means retry based on error code; 1 means retry based on status code
		if errorCode == "nil" {
			retryIndex = -1
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryIndex == 0 {
				retryIndex++
				return failedResponseMock(errorCode)
			} else if retryIndex == 1 {
				retryIndex++
				return successResponseMock(CreateMockResponse)
			}
			return CreateMockResponse, nil
		})
		err := resourceAliCloudEcsDiskEncryptionByDefaultCreate(dCreate, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk_encryption_by_default"].Schema).Data(dCreate.State(), nil)
			_ = dCompare.Set("enabled", attributes["enabled"])
			assert.Equal(t, dCompare.State().Attributes, dCreate.State().Attributes)
		}
		if retryIndex >= 0 {
			patches.Reset()
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudEcsDiskEncryptionByDefaultUpdate(d, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// DisableDiskEncryptionByDefault
	attributesDiff := map[string]interface{}{
		"enabled": false,
	}
	diff, err := newInstanceDiff("alicloud_ecs_disk_encryption_by_default", attributes, attributesDiff, dCreate.State())
	if err != nil {
		t.Error(err)
	}
	dExisted2Data, _ = schema.InternalMap(p["alicloud_ecs_disk_encryption_by_default"].Schema).Data(dCreate.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDiskEncryptionByDefaultStatus Response
		"Encrypted": false,
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		if errorCode == "nil" {
			retryIndex = -1
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryIndex == 0 {
				retryIndex++
				return failedResponseMock(errorCode)
			} else if retryIndex == 1 {
				retryIndex++
				return successResponseMock(ReadMockResponseDiff)
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudEcsDiskEncryptionByDefaultUpdate(dExisted2Data, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk_encryption_by_default"].Schema).Data(dExisted2Data.State(), nil)
			_ = dCompare.Set("enabled", attributesDiff["enabled"])
			assert.Equal(t, dCompare.State().Attributes, dExisted2Data.State().Attributes)
		}
		if retryIndex >= 0 {
			patches.Reset()
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		if errorCode == "nil" {
			retryIndex = -1
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryIndex == 0 {
				retryIndex++
				return failedResponseMock(errorCode)
			} else if retryIndex == 1 {
				retryIndex++
				return successResponseMock(ReadMockResponse)
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudEcsDiskEncryptionByDefaultRead(d, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEcsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAliCloudEcsDiskEncryptionByDefaultDelete(d, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		if errorCode == "nil" {
			retryIndex = -1
		}
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryIndex == 0 {
				retryIndex++
				return failedResponseMock(errorCode)
			} else if retryIndex == 1 {
				retryIndex++
				return successResponseMock(ReadMockResponse)
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudEcsDiskEncryptionByDefaultDelete(d, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
		}
	}
}

var attributes = map[string]interface{}{
	"enabled": true,
}
