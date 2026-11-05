package dfs_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	gomonkey "github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dfs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("alicloud_dfs_access_group", &resource.Sweeper{
		Name: "alicloud_dfs_access_group",
		F:    testSweepDFSAccessGroup,
	})
}

func testSweepDFSAccessGroup(region string) error {
	rawClient, err := sweep.SharedClientForRegionWithBackendRegions(region, true, connectivity.DfsSupportRegions)
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

	action := "ListAccessGroups"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}

	resp, err := jsonpath.Get("$.AccessGroups", response)
	if err != nil {
		log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.AccessGroups", action, err)
		return nil
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(item["AccessGroupName"].(string)), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping DFS AccessGroup: %s", item["AccessGroupName"].(string))
			continue
		}

		action := "DeleteAccessGroup"
		request := map[string]interface{}{
			"AccessGroupId": item["AccessGroupId"].(string),
			"InputRegionId": client.RegionId,
		}

		_, err = client.RpcPost("DFS", "2018-06-20", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to delete DFS AccessGroup (%s): %s", item["AccessGroupName"].(string), err)
		}
		log.Printf("[INFO] Delete  DFS AccessGroup success: %s ", item["AccessGroupName"].(string))
	}

	return nil
}

func TestAccAliCloudDfsAccessGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_access_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDFSAccessGroupMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dfs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDfsAccessGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsaccessgroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDFSAccessGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DfsSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network_type":      "VPC",
					"description":       "${var.name}_Desc",
					"access_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":      "VPC",
						"description":       name + "_Desc",
						"access_group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_Desc_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_Desc_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${var.name}_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.name}_Desc",
					"access_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       name + "_Desc",
						"access_group_name": name,
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

var AlicloudDFSAccessGroupMap0 = map[string]string{
	"network_type": "VPC",
}

func AlicloudDFSAccessGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}

// lintignore: R001
func TestUnitAlicloudDFSAccessGroup(t *testing.T) {
	resourceSchema := dfs.ResourceAccessGroup()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"network_type":      "CreateAccessGroupValue",
		"description":       "CreateAccessGroupValue",
		"access_group_name": "CreateAccessGroupValue",
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
		// GetAccessGroup
		"AccessGroup": map[string]interface{}{
			"AccessGroupName": "CreateAccessGroupValue",
			"Description":     "CreateAccessGroupValue",
			"NetworkType":     "CreateAccessGroupValue",
		},
		"AccessGroupId": "CreateAccessGroupValue",
	}
	CreateMockResponse := map[string]interface{}{
		"AccessGroupId": "CreateAccessGroupValue",
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dfs_access_group", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			helper.MapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAlidfsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceSchema.Create(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateAccessGroup" {
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
		err := resourceSchema.Create(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), nil)
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
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceSchema.Update(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	//ModifyAccessGroup
	attributesDiff := map[string]interface{}{
		"access_group_name": "ModifyAccessGroupValue",
		"description":       "ModifyAccessGroupValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_dfs_access_group", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAccessGroup Response
		"AccessGroup": map[string]interface{}{
			"AccessGroupName": "ModifyAccessGroupValue",
			"Description":     "ModifyAccessGroupValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyAccessGroup" {
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
		err := resourceSchema.Update(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(resourceSchema.Schema).Data(dExisted.State(), nil)
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
			if *action == "GetAccessGroup" {
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
		err := resourceSchema.Read(dExisted, rawClient)
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
			Code:    tea.String("loadEndpoint error"),
			Data:    tea.String("loadEndpoint error"),
			Message: tea.String("loadEndpoint error"),
		}
	})
	err = resourceSchema.Delete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidParameter.AccessGroupNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteAccessGroup" {
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
		err := resourceSchema.Delete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidParameter.AccessGroupNotFound":
			assert.Nil(t, err)
		}
	}
}

// Test Dfs AccessGroup. >>> Resource test cases, automatically generated.
// Case AccessGroup资源测试用例 5176
func TestAccAliCloudDfsAccessGroup_basic5176(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_access_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDfsAccessGroupMap5176)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dfs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDfsAccessGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsaccessgroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDfsAccessGroupBasicDependence5176)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network_type":      "VPC",
					"access_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_type":      "VPC",
						"access_group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "AccessGroup resource manager center test case",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "AccessGroup resource manager center test case",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "AccessGroup resource manager center test case fix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "AccessGroup resource manager center test case fix",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "AccessGroup resource manager center test case",
					"network_type":      "VPC",
					"access_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "AccessGroup resource manager center test case",
						"network_type":      "VPC",
						"access_group_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDfsAccessGroupMap5176 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudDfsAccessGroupBasicDependence5176(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case AccessGroup资源测试用例 5176  twin
func TestAccAliCloudDfsAccessGroup_basic5176_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_access_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDfsAccessGroupMap5176)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dfs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDfsAccessGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsaccessgroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDfsAccessGroupBasicDependence5176)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "AccessGroup resource manager center test case fix",
					"network_type":      "VPC",
					"access_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "AccessGroup resource manager center test case fix",
						"network_type":      "VPC",
						"access_group_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Test Dfs AccessGroup. <<< Resource test cases, automatically generated.
