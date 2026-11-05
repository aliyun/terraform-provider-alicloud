package dfs_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// lintignore: R001
func TestUnitAliCloudDFSMountPoint(t *testing.T) {
	resourceSchema := dfs.ResourceMountPoint()
	dInit, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(resourceSchema.Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"network_type":    "CreateMountPointValue",
		"vpc_id":          "CreateMountPointValue",
		"vswitch_id":      "CreateMountPointValue",
		"file_system_id":  "CreateMountPointValue",
		"description":     "CreateMountPointValue",
		"access_group_id": "CreateMountPointValue",
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
		// GetMountPoint
		"MountPoint": map[string]interface{}{
			"FileSystemId":  "CreateMountPointValue",
			"AccessGroupId": "CreateMountPointValue",
			"Description":   "CreateMountPointValue",
			"NetworkType":   "CreateMountPointValue",
			"Status":        "CreateMountPointValue",
			"VpcId":         "CreateMountPointValue",
			"VSwitchId":     "CreateMountPointValue",
		},
		"MountPointId": "CreateMountPointValue",
	}
	CreateMockResponse := map[string]interface{}{
		"MountPointId": "CreateMountPointValue",
	}

	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       tea.String(errorCode),
			Data:       tea.String(errorCode),
			Message:    tea.String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("alicloud_dfs_mount_point", errorCode))
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
	ReadMockResponseDiff := map[string]interface{}{
		"MountPointId": "CreateMountPointValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateMountPoint" {
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
	//ModifyMountPoint
	attributesDiff := map[string]interface{}{
		"access_group_id": "ModifyMountPointValue",
		"description":     "ModifyMountPointValue",
		"status":          "ModifyMountPointValue",
	}
	diff, err := acctest.NewInstanceDiff("alicloud_dfs_mount_point", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetMountPoint Response
		"MountPoint": map[string]interface{}{
			"AccessGroupId": "ModifyMountPointValue",
			"Description":   "ModifyMountPointValue",
			"Status":        "ModifyMountPointValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyMountPoint" {
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

	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_dfs_mount_point", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetMountPoint" {
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
	attributesDiff = map[string]interface{}{}
	diff, err = acctest.NewInstanceDiff("alicloud_dfs_mount_point", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(resourceSchema.Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidParameter.MountPointNotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteMountPoint" {
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
		case "InvalidParameter.MountPointNotFound":
			assert.Nil(t, err)
		}
	}
}

// Test Dfs MountPoint. >>> Resource test cases, automatically generated.
// Case MountPoint资源测试用例_2.0 5564
func TestAccAliCloudDfsMountPoint_basic5564(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_mount_point.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDfsMountPointMap5564)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dfs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDfsMountPoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsmountpoint%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDfsMountPointBasicDependence5564)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":          "${alicloud_vpc.DefaultVPCRMC.id}",
					"network_type":    "VPC",
					"vswitch_id":      "${alicloud_vswitch.DefaultVSwitch.id}",
					"file_system_id":  "${alicloud_dfs_file_system.DefaultFs.id}",
					"access_group_id": "${alicloud_dfs_access_group.DefaultAccessGroupRMC.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":          acctest.CHECKSET,
						"network_type":    "VPC",
						"vswitch_id":      acctest.CHECKSET,
						"file_system_id":  acctest.CHECKSET,
						"access_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alias_prefix": "MPAliasRMCTest39",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alias_prefix": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "mountpoint RMC test case",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "mountpoint RMC test case",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "mountpoint RMC test case fix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "mountpoint RMC test case fix",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_id": "${alicloud_dfs_access_group.UpdateAccessGroup.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "mountpoint RMC test case",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "mountpoint RMC test case",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_id": "${alicloud_dfs_access_group.DefaultAccessGroupRMC.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "mountpoint RMC test case fix",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "mountpoint RMC test case fix",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_id": "${alicloud_dfs_access_group.UpdateAccessGroup.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":          "${alicloud_vpc.DefaultVPCRMC.id}",
					"description":     "mountpoint RMC test case",
					"network_type":    "VPC",
					"vswitch_id":      "${alicloud_vswitch.DefaultVSwitch.id}",
					"file_system_id":  "${alicloud_dfs_file_system.DefaultFs.id}",
					"access_group_id": "${alicloud_dfs_access_group.DefaultAccessGroupRMC.id}",
					"alias_prefix":    "MPAliasRMCTest39",
					"status":          "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":          acctest.CHECKSET,
						"description":     "mountpoint RMC test case",
						"network_type":    "VPC",
						"vswitch_id":      acctest.CHECKSET,
						"file_system_id":  acctest.CHECKSET,
						"access_group_id": acctest.CHECKSET,
						"alias_prefix":    acctest.CHECKSET,
						"status":          "Active",
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

var AliCloudDfsMountPointMap5564 = map[string]string{
	"status":         acctest.CHECKSET,
	"mount_point_id": acctest.CHECKSET,
	"create_time":    acctest.CHECKSET,
	"region_id":      acctest.CHECKSET,
}

func AliCloudDfsMountPointBasicDependence5564(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "DefaultVPCRMC" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name

}

resource "alicloud_vswitch" "DefaultVSwitch" {
  description  = "rmc test"
  vpc_id       = alicloud_vpc.DefaultVPCRMC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
  zone_id = "cn-hangzhou-e"
}

resource "alicloud_dfs_access_group" "DefaultAccessGroupRMC" {
  description       = "AccessGroup resource manager center test case for mp"
  network_type      = "VPC"
  access_group_name = var.name

}

resource "alicloud_dfs_access_group" "UpdateAccessGroup" {
  description       = "Second AccessGroup resource manager center test case for mp"
  network_type      = "VPC"
  access_group_name = join("-", [var.name, "0"])

}

resource "alicloud_dfs_file_system" "DefaultFs" {
  space_capacity       = "1024"
  description          = "for mountpoint RMC test"
  storage_type         = "STANDARD"
  zone_id              = "cn-hangzhou-e"
  protocol_type        = "PANGU"
  data_redundancy_type = "LRS"
  file_system_name     = var.name

}


`, name)
}

// Case MountPoint资源测试用例_2.0 5564  twin
func TestAccAliCloudDfsMountPoint_basic5564_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dfs_mount_point.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDfsMountPointMap5564)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dfs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDfsMountPoint")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdfsmountpoint%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDfsMountPointBasicDependence5564)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":          "${alicloud_vpc.DefaultVPCRMC.id}",
					"description":     "mountpoint RMC test case fix",
					"network_type":    "VPC",
					"vswitch_id":      "${alicloud_vswitch.DefaultVSwitch.id}",
					"file_system_id":  "${alicloud_dfs_file_system.DefaultFs.id}",
					"access_group_id": "${alicloud_dfs_access_group.UpdateAccessGroup.id}",
					"status":          "Inactive",
					"alias_prefix":    "MPAliasRMCTest14",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":          acctest.CHECKSET,
						"description":     "mountpoint RMC test case fix",
						"network_type":    "VPC",
						"vswitch_id":      acctest.CHECKSET,
						"file_system_id":  acctest.CHECKSET,
						"access_group_id": acctest.CHECKSET,
						"status":          "Inactive",
						"alias_prefix":    acctest.CHECKSET,
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

// Test Dfs MountPoint. <<< Resource test cases, automatically generated.
