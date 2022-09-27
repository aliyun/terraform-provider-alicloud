package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

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
	resource.AddTestSweepers(
		"alicloud_dfs_file_system",
		&resource.Sweeper{
			Name: "alicloud_dfs_file_system",
			F:    testSweepDFSFileSystem,
		})
}

func testSweepDFSFileSystem(region string) error {
	rawClient, err := sharedClientForRegionWithBackendRegions(region, true, connectivity.DfsSupportRegions)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := map[string]interface{}{
		"InputRegionId": client.RegionId,
	}

	action := "ListFileSystems"
	conn, err := client.NewAlidfsClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	var response map[string]interface{}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.FileSystems", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.FileSystems", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["FileSystemName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping DFS FileSystem: %s", item["FileSystemName"].(string))
			continue
		}

		action := "DeleteFileSystem"
		request := map[string]interface{}{
			"FileSystemId":  item["FileSystemId"].(string),
			"InputRegionId": client.RegionId,
		}

		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to delete DFS FileSystem (%s): %s", item["FileSystemName"].(string), err)
		}
		log.Printf("[INFO] Delete  DFS FileSystem success: %s ", item["FileSystemName"].(string))
	}

	return nil
}

func TestAccAlicloudDFSFileSystem_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_file_system.default"
	checkoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudDFSFileSystemMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsfilesystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDFSFileSystemBasicDependence0)
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
					"storage_type":                     "${local.storage_type}",
					"zone_id":                          "${local.zone_id}",
					"protocol_type":                    "HDFS",
					"description":                      name,
					"file_system_name":                 name,
					"space_capacity":                   "1024",
					"throughput_mode":                  "Provisioned",
					"provisioned_throughput_in_mi_bps": "512",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_type":                     CHECKSET,
						"zone_id":                          CHECKSET,
						"protocol_type":                    "HDFS",
						"description":                      name,
						"file_system_name":                 name,
						"space_capacity":                   "1024",
						"throughput_mode":                  "Provisioned",
						"provisioned_throughput_in_mi_bps": "512",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"space_capacity": "2048",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"space_capacity": "2048",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      name,
					"file_system_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      name,
						"file_system_name": name,
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

var AlicloudDFSFileSystemMap0 = map[string]string{}

func AlicloudDFSFileSystemBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_dfs_zones" "default" {}
locals {
  zone_id      = data.alicloud_dfs_zones.default.zones.0.zone_id
  storage_type = data.alicloud_dfs_zones.default.zones.0.options.0.storage_type
}
`, name)
}

func TestUnitAlicloudDFSFileSystem(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_dfs_file_system"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_dfs_file_system"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"storage_type":                     "CreateFileSystemValue",
		"zone_id":                          "CreateFileSystemValue",
		"protocol_type":                    "CreateFileSystemValue",
		"description":                      "CreateFileSystemValue",
		"file_system_name":                 "CreateFileSystemValue",
		"space_capacity":                   1024,
		"throughput_mode":                  "CreateFileSystemValue",
		"provisioned_throughput_in_mi_bps": 512,
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
		// GetFileSystem
		"FileSystem": map[string]interface{}{
			"Description":                  "CreateFileSystemValue",
			"FileSystemName":               "CreateFileSystemValue",
			"ProtocolType":                 "CreateFileSystemValue",
			"ProvisionedThroughputInMiBps": 512,
			"SpaceCapacity":                1024,
			"StorageType":                  "CreateFileSystemValue",
			"ThroughputMode":               "CreateFileSystemValue",
			"ZoneId":                       "CreateFileSystemValue",
		},
		"FileSystemId": "CreateFileSystemValue",
	}
	CreateMockResponse := map[string]interface{}{
		"FileSystemId": "CreateFileSystemValue",
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_dfs_file_system", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidfsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudDfsFileSystemCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateFileSystem" {
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
		err := resourceAlicloudDfsFileSystemCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dfs_file_system"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidfsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudDfsFileSystemUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	//ModifyAccessGroup
	attributesDiff := map[string]interface{}{
		"file_system_name":                 "ModifyFileSystemValue",
		"description":                      "ModifyFileSystemValue",
		"provisioned_throughput_in_mi_bps": 256,
		"space_capacity":                   512,
		"throughput_mode":                  "ModifyFileSystemValue",
	}
	diff, err := newInstanceDiff("alicloud_dfs_file_system", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dfs_file_system"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetFileSystem Response
		"FileSystem": map[string]interface{}{
			"Description":                  "ModifyFileSystemValue",
			"FileSystemName":               "ModifyFileSystemValue",
			"ProvisionedThroughputInMiBps": 256,
			"SpaceCapacity":                512,
			"ThroughputMode":               "ModifyFileSystemValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyFileSystem" {
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
		err := resourceAlicloudDfsFileSystemUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dfs_file_system"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetFileSystem" {
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
		err := resourceAlicloudDfsFileSystemRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidfsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudDfsFileSystemDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidParameter.FileSystemNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteFileSystem" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{
							"Success": true,
						}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudDfsFileSystemDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidParameter.FileSystemNotFound":
			assert.Nil(t, err)
		}
	}
}
