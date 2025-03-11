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
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ecs_disk", &resource.Sweeper{
		Name: "alicloud_ecs_disk",
		F:    testAliCloudEcsDisk,
	})
}

func testAliCloudEcsDisk(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}

	var response map[string]interface{}
	for {
		action := "DescribeDisks"
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %s", action, err)
			return nil
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Disks.Disk", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Disks.Disk", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DiskName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Disk: %s", item["DiskName"].(string))
				continue
			}
			action = "DeleteDisk"
			request := map[string]interface{}{
				"DiskId":   item["DiskId"],
				"RegionId": client.RegionId,
			}
			_, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Disk (%s): %s", item["DiskName"].(string), err)
			}
			log.Printf("[INFO] Delete Disk success: %s ", item["DiskName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestUnitECSDisk(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"advanced_features": "CreateDiskValue",
		"disk_name":         "CreateDiskValue",
		"description":       "CreateDiskValue",
		"zone_id":           "CreateDiskValue",
		"size":              500,
		"payment_type":      "PayAsYouGo",
		"category":          "CreateDiskValue",
		"resource_group_id": "CreateDiskValue",
		"snapshot_id":       "CreateDiskValue",
		"tags": map[string]string{
			"TagResourcesValue_1": "CreateDiskValue",
			"TagResourcesValue_2": "CreateDiskValue",
		},
		"encrypt_algorithm":                  "CreateDiskValue",
		"instance_id":                        "CreateDiskValue",
		"kms_key_id":                         "CreateDiskValue",
		"performance_level":                  "CreateDiskValue",
		"storage_set_id":                     "CreateDiskValue",
		"storage_set_partition_number":       1,
		"dedicated_block_storage_cluster_id": "CreateDiskValue",
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
		// DescribeDisks
		"Disks": map[string]interface{}{
			"Disk": []interface{}{
				map[string]interface{}{
					"DiskId":             "CreateDiskValue",
					"Category":           "CreateDiskValue",
					"DeleteAutoSnapshot": true,
					"DeleteWithInstance": true,
					"Description":        "CreateDiskValue",
					"DiskName":           "CreateDiskValue",
					"EnableAutoSnapshot": false,
					"Encrypted":          "CreateDiskValue",
					"InstanceId":         "CreateDiskValue",
					"KMSKeyId":           "CreateDiskValue",
					"DiskChargeType":     "PostPaid",
					"PerformanceLevel":   "CreateDiskValue",
					"ResourceGroupId":    "CreateDiskValue",
					"Size":               500,
					"SourceSnapshotId":   "CreateDiskValue",
					"Status":             "Available",
					"ZoneId":             "CreateDiskValue",
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "TagResourcesValue_1",
								"Value": "CreateDiskValue",
							},
							map[string]interface{}{
								"Key":   "TagResourcesValue_2",
								"Value": "CreateDiskValue",
							},
						},
					},
				},
			},
		},
		"DiskId": "CreateDiskValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDisk
		"DiskId": "CreateDiskValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ecs_disk", errorCode))
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
	err = resourceAliCloudEcsDiskCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDisks Response
		"DiskId": "CreateDiskValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDisk" {
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
		err := resourceAliCloudEcsDiskCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
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
	err = resourceAliCloudEcsDiskUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ResizeDisk
	attributesDiff := map[string]interface{}{
		"size": 600,
		"type": "ResizeDiskValue",
	}
	diff, err := newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDisks Response
		"Disks": map[string]interface{}{
			"Disk": []interface{}{
				map[string]interface{}{
					"Size": 600,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ResizeDisk" {
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
		err := resourceAliCloudEcsDiskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// JoinResourceGroup
	attributesDiff = map[string]interface{}{
		"resource_group_id": "JoinResourceGroupValue",
	}
	diff, err = newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDisks Response
		"Disks": map[string]interface{}{
			"Disk": []interface{}{
				map[string]interface{}{
					"ResourceGroupId": "JoinResourceGroupValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "JoinResourceGroup" {
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
		err := resourceAliCloudEcsDiskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// ModifyDiskSpec
	attributesDiff = map[string]interface{}{
		"category":          "ModifyDiskSpecValue",
		"performance_level": "ModifyDiskSpecValue",
		"dry_run":           false,
	}
	diff, err = newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDisks Response
		"Disks": map[string]interface{}{
			"Disk": []interface{}{
				map[string]interface{}{
					"Category":         "ModifyDiskSpecValue",
					"PerformanceLevel": "ModifyDiskSpecValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDiskSpec" {
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
		err := resourceAliCloudEcsDiskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// ModifyDiskChargeType
	attributesDiff = map[string]interface{}{
		"instance_id":  "ModifyDiskChargeTypeValue",
		"payment_type": "Subscription",
	}
	diff, err = newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDisks Response
		"Disks": map[string]interface{}{
			"Disk": []interface{}{
				map[string]interface{}{
					"InstanceId":     "ModifyDiskChargeTypeValue",
					"DiskChargeType": "PrePaid",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDiskChargeType" {
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
		err := resourceAliCloudEcsDiskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// ModifyDiskAttribute
	attributesDiff = map[string]interface{}{
		"delete_auto_snapshot": false,
		"delete_with_instance": false,
		"disk_name":            "ModifyDiskAttributeValue",
		"name":                 "ModifyDiskAttributeValue",
		"description":          "ModifyDiskAttributeValue",
		"enable_auto_snapshot": true,
	}
	diff, err = newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDisks Response
		"Disks": map[string]interface{}{
			"Disk": []interface{}{
				map[string]interface{}{
					"DeleteAutoSnapshot": false,
					"DeleteWithInstance": false,
					"Description":        "ModifyDiskAttributeValue",
					"DiskName":           "ModifyDiskAttributeValue",
					"EnableAutoSnapshot": true,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDiskAttribute" {
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
		err := resourceAliCloudEcsDiskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}
	// TagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"TagResourcesValue_1": "TagResourcesValue_1",
			"TagResourcesValue_2": "TagResourcesValue_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDisks Response
		"Disks": map[string]interface{}{
			"Disk": []interface{}{
				map[string]interface{}{
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "TagResourcesValue_1",
								"Value": "TagResourcesValue_1",
							},
							map[string]interface{}{
								"Key":   "TagResourcesValue_2",
								"Value": "TagResourcesValue_2",
							},
						},
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "TagResources" {
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
		err := resourceAliCloudEcsDiskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// UntagResources
	attributesDiff = map[string]interface{}{
		"tags": map[string]interface{}{
			"UntagResourcesValue3_1": "UnTagResourcesValue3_1",
			"UntagResourcesValue3_2": "UnTagResourcesValue3_2",
		},
	}
	diff, err = newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDisks Response
		"Disks": map[string]interface{}{
			"Disk": []interface{}{
				map[string]interface{}{
					"Tags": map[string]interface{}{
						"Tag": []interface{}{
							map[string]interface{}{
								"Key":   "UntagResourcesValue3_1",
								"Value": "UnTagResourcesValue3_1",
							},
							map[string]interface{}{
								"Key":   "UntagResourcesValue3_2",
								"Value": "UnTagResourcesValue3_2",
							},
						},
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UntagResources" {
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
		err := resourceAliCloudEcsDiskUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dExisted.State(), nil)
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
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDisks" {
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
		err := resourceAliCloudEcsDiskRead(dExisted, rawClient)
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
	err = resourceAliCloudEcsDiskDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_ecs_disk", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_ecs_disk"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "IncorrectDiskStatus.Initializing", "nil", "InvalidDiskId.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDisk" {
				switch errorCode {
				case "NonRetryableError", "InvalidDiskId.NotFound":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Disks": map[string]interface{}{
								"Disk": []interface{}{
									map[string]interface{}{
										"DiskId": "CreateDiskValue",
									},
								},
							},
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAliCloudEcsDiskDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidDiskId.NotFound":
			assert.Nil(t, err)
		}
	}
}

// Test Ecs Disk. >>> Resource test cases, automatically generated.
// Case cloud_efficiency生命周期 8017
func TestAccAliCloudECSDisk_basic8017(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8017)
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
					"size":    "500",
					"zone_id": "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size":    "500",
						"zone_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category": "cloud_essd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "cloud_essd",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"performance_level": "PL2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"performance_level": "PL2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8017_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8017)
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
					"category":             "cloud_efficiency",
					"delete_auto_snapshot": "true",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"enable_auto_snapshot": "false",
					"payment_type":         "PayAsYouGo",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size":        "500",
					"snapshot_id": "${alicloud_ecs_snapshot.default.id}",
					"zone_id":     "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":             "cloud_efficiency",
						"delete_auto_snapshot": "true",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"enable_auto_snapshot": "false",
						"payment_type":         "PayAsYouGo",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
						"snapshot_id":          CHECKSET,
						"zone_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

var AliCloudEcsDiskMap8017 = map[string]string{
	"create_time":          CHECKSET,
	"delete_with_instance": CHECKSET,
	"enable_auto_snapshot": CHECKSET,
	"multi_attach":         CHECKSET,
	"payment_type":         CHECKSET,
	"region_id":            CHECKSET,
	"status":               CHECKSET,
}

func AliCloudEcsDiskBasicDependence8017(name string) string {
	return fmt.Sprintf(`
	variable "name" {
 		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
 		available_resource_creation = "VSwitch"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name        = var.name
  		description = "New security group"
  		vpc_id      = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_disk" "default" {
  		name              = var.name
  		availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		category          = "cloud_efficiency"
  		size              = "20"
	}

	data "alicloud_images" "default" {
  		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		owners        = "system"
  		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	resource "alicloud_instance" "default" {
  		availability_zone = data.alicloud_zones.default.zones.0.id
  		instance_name     = var.name
  		host_name         = "tf-testAcc"
  		image_id          = data.alicloud_images.default.images.0.id
  		instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups   = [alicloud_security_group.default.id]
  		vswitch_id        = data.alicloud_vswitches.default.ids.0
	}

	resource "alicloud_disk_attachment" "default" {
  		disk_id     = alicloud_disk.default.id
  		instance_id = alicloud_instance.default.id
	}

	resource "alicloud_ecs_snapshot" "default" {
  		category       = "standard"
  		disk_id        = alicloud_disk_attachment.default.disk_id
  		retention_days = "20"
  		snapshot_name  = var.name
	}
`, name)
}

// Case cloud_essd生命周期, 适配废弃字段availability_zone, name 8018
func TestAccAliCloudECSDisk_basic8018(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
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
					"category":          "cloud_essd",
					"size":              "500",
					"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":          "cloud_essd",
						"size":              "500",
						"availability_zone": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"performance_level": "PL2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"performance_level": "PL2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8018_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
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
					"category":             "cloud_essd",
					"delete_auto_snapshot": "true",
					"description":          name,
					"name":                 name,
					"enable_auto_snapshot": "false",
					"encrypted":            "true",
					"kms_key_id":           "${alicloud_kms_key.default.id}",
					"multi_attach":         "Enabled",
					"payment_type":         "PayAsYouGo",
					"performance_level":    "PL2",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size":              "500",
					"availability_zone": "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":             "cloud_essd",
						"delete_auto_snapshot": "true",
						"description":          name,
						"name":                 name,
						"enable_auto_snapshot": "false",
						"encrypted":            "true",
						"kms_key_id":           CHECKSET,
						"multi_attach":         "Enabled",
						"payment_type":         "PayAsYouGo",
						"performance_level":    "PL2",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
						"availability_zone":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

// Case cloud_auto生命周期 8019
func TestAccAliCloudECSDisk_basic8019(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
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
					"category": "cloud_auto",
					"size":     "500",
					"zone_id":  "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "cloud_auto",
						"size":     "500",
						"zone_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"provisioned_iops": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"provisioned_iops": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8019_twin_excluding_kms_key_id(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
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
					"bursting_enabled":     "true",
					"category":             "cloud_auto",
					"delete_auto_snapshot": "true",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"enable_auto_snapshot": "false",
					"encrypted":            "true",
					"payment_type":         "PayAsYouGo",
					"provisioned_iops":     "100",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size":    "500",
					"zone_id": "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled":     "true",
						"category":             "cloud_auto",
						"delete_auto_snapshot": "true",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"enable_auto_snapshot": "false",
						"encrypted":            "true",
						"kms_key_id":           CHECKSET,
						"payment_type":         "PayAsYouGo",
						"provisioned_iops":     "100",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
						"zone_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8019_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8018)
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
					"bursting_enabled":     "true",
					"category":             "cloud_auto",
					"delete_auto_snapshot": "true",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"enable_auto_snapshot": "false",
					"encrypted":            "true",
					"kms_key_id":           "${alicloud_kms_key.default.id}",
					"payment_type":         "PayAsYouGo",
					"provisioned_iops":     "100",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size":    "500",
					"zone_id": "${data.alicloud_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled":     "true",
						"category":             "cloud_auto",
						"delete_auto_snapshot": "true",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"enable_auto_snapshot": "false",
						"encrypted":            "true",
						"kms_key_id":           CHECKSET,
						"payment_type":         "PayAsYouGo",
						"provisioned_iops":     "100",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
						"zone_id":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

func AliCloudEcsDiskBasicDependence8018(name string) string {
	return fmt.Sprintf(`
	variable "name" {
 		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_zones" "default" {
 		available_resource_creation = "VSwitch"
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		pending_window_in_days = "7"
  		key_state              = "Enabled"
	}
`, name)
}

// Case 预付费生命周期 8020
func TestAccAliCloudECSDisk_basic8020(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8020)
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
					"instance_id":  "${alicloud_instance.default.id}",
					"payment_type": "Subscription",
					"size":         "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"payment_type": "Subscription",
						"size":         "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_auto_snapshot": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_auto_snapshot": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_snapshot": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_snapshot": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

func TestAccAliCloudECSDisk_basic8020_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskMap8017)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskBasicDependence8020)
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
					"category":             "cloud_efficiency",
					"delete_auto_snapshot": "true",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"enable_auto_snapshot": "false",
					"instance_id":          "${alicloud_instance.default.id}",
					"payment_type":         "Subscription",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"size": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":             "cloud_efficiency",
						"delete_auto_snapshot": "true",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"enable_auto_snapshot": "false",
						"instance_id":          CHECKSET,
						"payment_type":         "Subscription",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
						"size":                 "500",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"advanced_features", "dry_run", "encrypt_algorithm", "type"},
			},
		},
	})
}

func AliCloudEcsDiskBasicDependence8020(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
  		status = "OK"
	}

	data "alicloud_instance_types" "default" {
		instance_charge_type = "PrePaid"
	}

	data "alicloud_images" "default" {
		name_regex    = "^ubuntu_[0-9]+_[0-9]+_x64*"
		owners        = "system"
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
		vpc_id  = data.alicloud_vpcs.default.ids.0
		zone_id = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
	}

	resource "alicloud_security_group" "default" {
		name   = var.name
		vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
	}

	resource "alicloud_instance" "default" {
		image_id                      = data.alicloud_images.default.images.0.id
		security_groups               = [alicloud_security_group.default.id]
		instance_type                 = data.alicloud_instance_types.default.instance_types.0.id
		system_disk_category          = "cloud_efficiency"
		instance_name                 = var.name
		spot_strategy                 = "NoSpot"
		spot_price_limit              = "0"
		security_enhancement_strategy = "Active"
		user_data                     = "I_am_user_data"
		instance_charge_type          = "PrePaid"
		period                        = 1
		vswitch_id                    = data.alicloud_vswitches.default.ids.0
		force_delete                  = true
	}
`, name)
}

// Test Ecs Disk. <<< Resource test cases, automatically generated.
