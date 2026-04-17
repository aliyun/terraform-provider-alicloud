package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccAliCloudEcsDiskDefaultKmsKeyId_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_default_kms_key_id.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskDefaultKmsKeyIdMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDiskDefaultKMSKeyId")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsDiskDefaultKmsKeyId%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskDefaultKmsKeyIdBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEcsDiskDefaultKmsKeyIdDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_key_id": "${alicloud_kms_key.example.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_key_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_key_id": "${alicloud_kms_key.example2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_key_id": CHECKSET,
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

func testAccCheckEcsDiskDefaultKmsKeyIdDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_disk_default_kms_key_id" {
			continue
		}

		// Check the encryption status to determine if the resource is "destroyed"
		// When the default KMS key is reset, encryption by default should be disabled
		object, err := ecsServiceV2.DescribeEcsDiskEncryptionByDefaultStatus(rs.Primary.ID)
		if err != nil {
			// If we can't get the encryption status, consider it as destroyed
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		// Check if encryption is disabled (this means the default KMS key is reset)
		if encrypted, ok := object["Encrypted"].(bool); ok && !encrypted {
			continue // Resource is properly "destroyed" (encryption disabled, meaning no default KMS key)
		}

		// If encryption is still enabled, check if there's actually a KMS key set
		kmsObject, err := ecsServiceV2.DescribeEcsDiskDefaultKMSKeyId(rs.Primary.ID)
		if err != nil {
			// If describe KMS key fails, consider it as destroyed
			continue
		}

		// If there's still a KMS key set, the resource wasn't properly destroyed
		if kmsKeyId, ok := kmsObject["KMSKeyId"]; ok && kmsKeyId != nil && kmsKeyId != "" {
			return WrapError(Error("ECS disk default KMS key ID is still set to %v in region %s", kmsKeyId, rs.Primary.ID))
		}
	}

	return nil
}

var AliCloudEcsDiskDefaultKmsKeyIdMap = map[string]string{}

func AliCloudEcsDiskDefaultKmsKeyIdBasicDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = "cn-hangzhou-i"
	}

	data "alicloud_vswitches" "default2" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = "cn-hangzhou-j"
	}

	# Create KMS instance
	resource "alicloud_kms_instance" "example" {
		product_version = "3"
		vpc_num         = "1"
		key_num         = "1000"
		secret_num      = "100"
		spec            = "1000"
		vpc_id          = data.alicloud_vpcs.default.ids.0
		vswitch_ids = [
			data.alicloud_vswitches.default.ids.0,
			data.alicloud_vswitches.default2.ids.0
		]
		zone_ids = [
			"cn-hangzhou-i",
			"cn-hangzhou-j"
		]
		payment_type                = "PayAsYouGo"
		force_delete_without_backup = "true"

		timeouts {
			delete = "20m"
		}
	}

	# Create a KMS key in the instance
	resource "alicloud_kms_key" "example" {
		description            = "KMS key for ECS disk encryption"
		pending_window_in_days = 7
		key_usage              = "ENCRYPT/DECRYPT"
		key_spec               = "Aliyun_AES_256"
		dkms_instance_id       = alicloud_kms_instance.example.id

		timeouts {
			delete = "20m"
		}
	}

	# Create a second KMS key for testing update
	resource "alicloud_kms_key" "example2" {
		description            = "Second KMS key for ECS disk encryption update test"
		pending_window_in_days = 7
		key_usage              = "ENCRYPT/DECRYPT"
		key_spec               = "Aliyun_AES_256"
		dkms_instance_id       = alicloud_kms_instance.example.id

		timeouts {
			delete = "20m"
		}
	}

	# Enable ECS disk encryption by default first
	resource "alicloud_ecs_disk_encryption_by_default" "example" {
		enabled = true
	}

`, name)
}

func TestUnitAliCloudEcsDiskDefaultKmsKeyId(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	var dExisted2Data *schema.ResourceData
	d, _ := schema.InternalMap(p["alicloud_ecs_disk_default_kms_key_id"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ecs_disk_default_kms_key_id"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	kmsKeyId := "7906979c-8e06-46a2-be2d-68e3ccbc****"
	d.Set("kms_key_id", kmsKeyId)
	dCreate.Set("kms_key_id", kmsKeyId)
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// DescribeDiskDefaultKMSKeyId
		"KMSKeyId": kmsKeyId,
	}
	CreateMockResponse := map[string]interface{}{
		// ModifyDiskDefaultKMSKeyId
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
	err = resourceAliCloudEcsDiskDefaultKmsKeyIdCreate(dCreate, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
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
		err := resourceAliCloudEcsDiskDefaultKmsKeyIdCreate(dCreate, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk_default_kms_key_id"].Schema).Data(dCreate.State(), nil)
			_ = dCompare.Set("kms_key_id", attributes2["kms_key_id"])
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
	err = resourceAliCloudEcsDiskDefaultKmsKeyIdUpdate(d, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDiskDefaultKMSKeyId
	newKmsKeyId := "7906979c-8e06-46a2-be2d-68e3ccbc****"
	attributesDiff := map[string]interface{}{
		"kms_key_id": newKmsKeyId,
	}
	diff, err := newInstanceDiff("alicloud_ecs_disk_default_kms_key_id", attributes2, attributesDiff, dCreate.State())
	if err != nil {
		t.Error(err)
	}
	dExisted2Data, _ = schema.InternalMap(p["alicloud_ecs_disk_default_kms_key_id"].Schema).Data(dCreate.State(), diff)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDiskDefaultKMSKeyId Response
		"KMSKeyId": newKmsKeyId,
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
		err := resourceAliCloudEcsDiskDefaultKmsKeyIdUpdate(dExisted2Data, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk_default_kms_key_id"].Schema).Data(dExisted2Data.State(), nil)
			_ = dCompare.Set("kms_key_id", attributesDiff["kms_key_id"])
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
		err := resourceAliCloudEcsDiskDefaultKmsKeyIdRead(d, rawClient)
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
	err = resourceAliCloudEcsDiskDefaultKmsKeyIdDelete(d, rawClient)
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
		err := resourceAliCloudEcsDiskDefaultKmsKeyIdDelete(d, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
		}
	}
}

var attributes2 = map[string]interface{}{
	"kms_key_id": "7906979c-8e06-46a2-be2d-68e3ccbc****",
}
