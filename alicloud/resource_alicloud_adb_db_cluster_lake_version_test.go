package alicloud

import (
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudADBDBClusterLakeVersion_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ADBDBClusterLakeVersionSupportRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdbDbClusterLakeVersionBasicDependence0)
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
					"vpc_id":                        "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":                    "${data.alicloud_vswitches.default.ids.0}",
					"compute_resource":              "16ACU",
					"storage_resource":              "0ACU",
					"db_cluster_version":            "5.0",
					"payment_type":                  "PayAsYouGo",
					"enable_default_resource_group": "false",
					"zone_id":                       "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":             CHECKSET,
						"vswitch_id":         CHECKSET,
						"compute_resource":   "16ACU",
						"storage_resource":   "0ACU",
						"db_cluster_version": "5.0",
						"payment_type":       "PayAsYouGo",
						"zone_id":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_default_resource_group"},
			},
		},
	})
}

var AlicloudAdbDbClusterLakeVersionMap0 = map[string]string{
	"resource_group_id": CHECKSET,
}

func AlicloudAdbDbClusterLakeVersionBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]
}
`, name)
}

func TestAccAlicloudADBDBClusterLakeVersion_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ADBDBClusterLakeVersionSupportRegions)
	resourceId := "alicloud_adb_db_cluster_lake_version.default"
	ra := resourceAttrInit(resourceId, AlicloudAdbDbClusterLakeVersionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAdbDbClusterLakeVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sadbdbclusterlakeversion%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAdbDbClusterLakeVersionBasicDependence0)
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
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":         "${data.alicloud_vswitches.default.ids.0}",
					"compute_resource":   "16ACU",
					"storage_resource":   "0ACU",
					"db_cluster_version": "5.0",
					"payment_type":       "PayAsYouGo",
					"zone_id":            "${data.alicloud_zones.default.ids[length(data.alicloud_zones.default.ids) - 1]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":             CHECKSET,
						"vswitch_id":         CHECKSET,
						"compute_resource":   "16ACU",
						"storage_resource":   "0ACU",
						"db_cluster_version": "5.0",
						"payment_type":       "PayAsYouGo",
						"zone_id":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_resource": "32ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_resource": "32ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"storage_resource": "24ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"storage_resource": "24ACU",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_resource": "48ACU",
					"storage_resource": "48ACU",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_resource": "48ACU",
						"storage_resource": "48ACU",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_default_resource_group"},
			},
		},
	})
}

func TestUnitAlicloudAdbDbClusterLakeVersion(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"compute_resource":              "CreateDBClusterValue",
		"db_cluster_version":            "5.0",
		"enable_default_resource_group": true,
		"payment_type":                  "CreateDBClusterValue",
		"storage_resource":              "CreateDBClusterValue",
		"vswitch_id":                    "CreateDBClusterValue",
		"zone_id":                       "CreateDBClusterValue",
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
		// DescribeDBClusterAttribute
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"CommodityCode":    "DefaultValue",
					"ComputeResource":  "CreateDBClusterValue",
					"ConnectionString": "DefaultValue",
					"CreationTime":     "DefaultValue",
					"DBClusterId":      "CreateDBClusterValue",
					"DBVersion":        "5.0",
					"Engine":           "DefaultValue",
					"EngineVersion":    "DefaultValue",
					"ExpireTime":       "DefaultValue",
					"Expired":          "DefaultValue",
					"LockMode":         "DefaultValue",
					"LockReason":       "DefaultValue",
					"PayType":          "CreateDBClusterValue",
					"Port":             "DefaultValue",
					"ResourceGroupId":  "DefaultValue",
					"DBClusterStatus":  "Running",
					"StorageResource":  "CreateDBClusterValue",
					"VPCId":            "CreateDBClusterValue",
					"VSwitchId":        "CreateDBClusterValue",
					"ZoneId":           "CreateDBClusterValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateDBCluster
		"DBClusterId":     "CreateDBClusterValue",
		"OrderId":         "MockValue",
		"RequestId":       "MockValue",
		"ResourceGroupId": "MockValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_adb_db_cluster_lake_version", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudAdbDbClusterLakeVersionCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDBClusterAttribute Response
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"DBClusterId": "CreateDBClusterValue",
				},
			},
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBCluster" {
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
		err := resourceAlicloudAdbDbClusterLakeVersionCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudAdbDbClusterLakeVersionUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDBCluster
	attributesDiff := map[string]interface{}{
		"compute_resource": "ModifyDBClusterValue",
		"storage_resource": "ModifyDBClusterValue",
	}
	diff, err := newInstanceDiff("alicloud_adb_db_cluster_lake_version", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBClusterAttribute Response
		"Items": map[string]interface{}{
			"DBCluster": []interface{}{
				map[string]interface{}{
					"ComputeResource": "ModifyDBClusterValue",
					"StorageResource": "ModifyDBClusterValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBCluster" {
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
		err := resourceAlicloudAdbDbClusterLakeVersionUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_adb_db_cluster_lake_version"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidDBCluster.NotFound", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBClusterAttribute" {
				switch errorCode {
				case "{}", "InvalidDBCluster.NotFound":
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
		err := resourceAlicloudAdbDbClusterLakeVersionRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}", "InvalidDBCluster.NotFound":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewAdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudAdbDbClusterLakeVersionDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InvalidDBCluster.NotFound"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBCluster" {
				switch errorCode {
				case "NonRetryableError", "InvalidDBCluster.NotFound":
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
		err := resourceAlicloudAdbDbClusterLakeVersionDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InvalidDBCluster.NotFound":
			assert.Nil(t, err)
		}
	}
}
