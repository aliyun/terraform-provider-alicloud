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
		return fmt.Errorf("error getting AliCloud client: %s", err)
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

func TestUnitAliCloudDFSFileSystem(t *testing.T) {
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
	err = resourceAliCloudDfsFileSystemCreate(dInit, rawClient)
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
		err := resourceAliCloudDfsFileSystemCreate(dInit, rawClient)
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
	err = resourceAliCloudDfsFileSystemUpdate(dExisted, rawClient)
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
		err := resourceAliCloudDfsFileSystemUpdate(dExisted, rawClient)
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
		err := resourceAliCloudDfsFileSystemRead(dExisted, rawClient)
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
	err = resourceAliCloudDfsFileSystemDelete(dExisted, rawClient)
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
		err := resourceAliCloudDfsFileSystemDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidParameter.FileSystemNotFound":
			assert.Nil(t, err)
		}
	}
}

// Test Dfs FileSystem. >>> Resource test cases, automatically generated.
// Case FileSystem资源测试用例_增加StorageType覆盖率 5910
func TestAccAliCloudDfsFileSystem_basic5910(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudDfsFileSystemMap5910)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsfilesystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDfsFileSystemBasicDependence5910)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"space_capacity":                   "1024",
					"description":                      "ResourceManagerCenterFsTestCase",
					"storage_type":                     "PERFORMANCE",
					"zone_id":                          "cn-hangzhou-b",
					"protocol_type":                    "PANGU",
					"file_system_name":                 name,
					"data_redundancy_type":             "LRS",
					"provisioned_throughput_in_mi_bps": "1000",
					"throughput_mode":                  "Provisioned",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"space_capacity":                   "1024",
						"description":                      "ResourceManagerCenterFsTestCase",
						"storage_type":                     "PERFORMANCE",
						"zone_id":                          "cn-hangzhou-b",
						"protocol_type":                    "PANGU",
						"file_system_name":                 name,
						"data_redundancy_type":             "LRS",
						"provisioned_throughput_in_mi_bps": "1000",
						"throughput_mode":                  "Provisioned",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_redundancy_type", "dedicated_cluster_id", "partition_number", "storage_set_name"},
			},
		},
	})
}

var AliCloudDfsFileSystemMap5910 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudDfsFileSystemBasicDependence5910(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}
`, name)
}

// Case FileSystem资源测试用例 5175
func TestAccAliCloudDfsFileSystem_basic5175(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_file_system.default"
	ra := resourceAttrInit(resourceId, AliCloudDfsFileSystemMap5175)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DfsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDfsFileSystem")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsfilesystem%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDfsFileSystemBasicDependence5175)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"space_capacity":                   "1024",
					"description":                      "ResourceManagerCenterFsTestCase",
					"storage_type":                     "STANDARD",
					"zone_id":                          "cn-hangzhou-e",
					"protocol_type":                    "PANGU",
					"file_system_name":                 name,
					"data_redundancy_type":             "LRS",
					"provisioned_throughput_in_mi_bps": "1000",
					"throughput_mode":                  "Provisioned",
					"partition_number":                 "0",
					"storage_set_name":                 "RMCTestStorageSet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"space_capacity":                   "1024",
						"description":                      "ResourceManagerCenterFsTestCase",
						"storage_type":                     "STANDARD",
						"zone_id":                          "cn-hangzhou-e",
						"protocol_type":                    "PANGU",
						"file_system_name":                 name,
						"data_redundancy_type":             "LRS",
						"provisioned_throughput_in_mi_bps": "1000",
						"throughput_mode":                  "Provisioned",
						"partition_number":                 "0",
						"storage_set_name":                 "RMCTestStorageSet",
					}),
				),
			},
			//provisioned_throughput_in_mi_bps can only be modified once a day
			{
				Config: testAccConfig(map[string]interface{}{
					"space_capacity":   "1026",
					"description":      "ResourceManagerCenterTestCase-fix",
					"file_system_name": name + "_update",
					//"provisioned_throughput_in_mi_bps": "1010",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"space_capacity":   "1026",
						"description":      "ResourceManagerCenterTestCase-fix",
						"file_system_name": name + "_update",
						//"provisioned_throughput_in_mi_bps": "1010",
					}),
				),
			},
			//throughput_mode can only be modified once a day
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"throughput_mode": "Standard",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"throughput_mode": "Standard",
			//		}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_redundancy_type", "dedicated_cluster_id", "partition_number", "storage_set_name"},
			},
		},
	})
}

var AliCloudDfsFileSystemMap5175 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudDfsFileSystemBasicDependence5175(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}
`, name)
}

// Test Dfs FileSystem. <<< Resource test cases, automatically generated.
