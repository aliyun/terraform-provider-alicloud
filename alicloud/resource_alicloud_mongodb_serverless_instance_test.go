package alicloud

import (
	"fmt"
	"log"
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
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMongoDBServerlessInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_serverless_instance.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBServerlessSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudMongoDBServerlessInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbServerlessInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbserverlessinstance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongoDBServerlessInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password":    "Abc12345",
					"db_instance_storage": "5",
					"capacity_unit":       "100",
					"engine_version":      "4.2",
					"vswitch_id":          "${local.vswitch_id}",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":             "${data.alicloud_mongodb_zones.default.zones.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "5",
						"capacity_unit":       "100",
						"engine_version":      "4.2",
						"vswitch_id":          CHECKSET,
						"vpc_id":              CHECKSET,
						"zone_id":             CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "MongodbServerlessInstance",
						"For":     "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "MongodbServerlessInstance",
						"tags.For":     "TF",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "02:00Z",
					"maintain_end_time":   "03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "02:00Z",
						"maintain_end_time":   "03:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test",
							"security_ip_group_name":      "test",
							"security_ip_list":            "192.168.0.1",
						},
						{
							"security_ip_group_attribute": "test1",
							"security_ip_group_name":      "test1",
							"security_ip_list":            "192.168.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_groups.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity_unit":       "2000",
					"db_instance_storage": "20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity_unit":       "2000",
						"db_instance_storage": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"capacity_unit":           "100",
					"db_instance_storage":     "5",
					"db_instance_description": "${var.name}_update",
					"maintain_start_time":     "01:00Z",
					"maintain_end_time":       "02:00Z",
					"tags": map[string]string{
						"Created": "MongodbServerlessInstance1",
						"For":     "TF1",
					},
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test3",
							"security_ip_group_name":      "test3",
							"security_ip_list":            "192.168.0.3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity_unit":           "100",
						"db_instance_storage":     "5",
						"db_instance_description": name + "_update",
						"maintain_start_time":     "01:00Z",
						"maintain_end_time":       "02:00Z",
						"tags.%":                  "2",
						"tags.Created":            "MongodbServerlessInstance1",
						"tags.For":                "TF1",
						"security_ip_groups.#":    "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "period", "period_price_type", "account_password"},
			},
		},
	})
}

func TestAccAlicloudMongoDBServerlessInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_serverless_instance.default"
	checkoutSupportedRegions(t, true, connectivity.MongoDBServerlessSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudMongoDBServerlessInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbServerlessInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-mongodbserverlessinstance-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMongoDBServerlessInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":              "false",
					"account_password":        "Abc12345",
					"capacity_unit":           "100",
					"db_instance_storage":     "5",
					"storage_engine":          "WiredTiger",
					"engine":                  "MongoDB",
					"resource_group_id":       "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"engine_version":          "4.2",
					"db_instance_description": "${var.name}",
					"vpc_id":                  "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":                 "${data.alicloud_mongodb_zones.default.zones.0.id}",
					"vswitch_id":              "${local.vswitch_id}",
					"period":                  "1",
					"period_price_type":       "Day",
					"tags": map[string]string{
						"Created": "MongodbServerlessInstance",
						"For":     "TF",
					},
					"security_ip_groups": []map[string]interface{}{
						{
							"security_ip_group_attribute": "test",
							"security_ip_group_name":      "test",
							"security_ip_list":            "192.168.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"capacity_unit":           "100",
						"db_instance_storage":     "5",
						"engine":                  "MongoDB",
						"engine_version":          "4.2",
						"storage_engine":          "WiredTiger",
						"db_instance_description": name,
						"vswitch_id":              CHECKSET,
						"tags.%":                  "2",
						"tags.Created":            "MongodbServerlessInstance",
						"tags.For":                "TF",
						"security_ip_groups.#":    "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "period", "period_price_type", "account_password"},
			},
		},
	})
}

var AlicloudMongoDBServerlessInstanceMap0 = map[string]string{
	"auto_pay": NOSET,
	"status":   CHECKSET,
}

func AlicloudMongoDBServerlessInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_vswitch" "default" {
	count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	vswitch_name = var.name
	vpc_id       = data.alicloud_vpcs.default.ids.0
    zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
    cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
}

locals {
  vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.default.*.id, [""])[0]
}
`, name)
}

func TestUnitAlicloudMongoDBServerlessInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"auto_renew":              false,
		"account_password":        "CreateServerlessDBInstanceValue",
		"capacity_unit":           100,
		"db_instance_storage":     5,
		"storage_engine":          "WiredTiger",
		"engine":                  "CreateServerlessDBInstanceValue",
		"resource_group_id":       "CreateServerlessDBInstanceValue",
		"engine_version":          "CreateServerlessDBInstanceValue",
		"db_instance_description": "CreateServerlessDBInstanceValue",
		"vpc_id":                  "CreateServerlessDBInstanceValue",
		"zone_id":                 "CreateServerlessDBInstanceValue",
		"vswitch_id":              "CreateServerlessDBInstanceValue",
		"period":                  1,
		"period_price_type":       "CreateServerlessDBInstanceValue",
		"tags": map[string]string{
			"Created": "CreateServerlessDBInstanceValue",
			"For":     "CreateServerlessDBInstanceValue",
		},
		"security_ip_groups": []map[string]interface{}{
			{
				"security_ip_group_attribute": "CreateServerlessDBInstanceValue",
				"security_ip_group_name":      "CreateServerlessDBInstanceValue",
				"security_ip_list":            "CreateServerlessDBInstanceValue",
			},
		},
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
		// DescribeDBInstanceAttribute
		"DBInstances": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"CapacityUnit":          100,
					"DBInstanceStatus":      "Running",
					"DBInstanceId":          "CreateServerlessDBInstanceValue",
					"DBInstanceDescription": "CreateServerlessDBInstanceValue",
					"DBInstanceStorage":     5,
					"Engine":                "CreateServerlessDBInstanceValue",
					"EngineVersion":         "CreateServerlessDBInstanceValue",
					"MaintainEndTime":       "CreateServerlessDBInstanceValue",
					"MaintainStartTime":     "CreateServerlessDBInstanceValue",
					"ResourceGroupId":       "CreateServerlessDBInstanceValue",
					"StorageEngine":         "wiredTiger",
					"VSwitchId":             "CreateServerlessDBInstanceValue",
					"ZoneId":                "CreateServerlessDBInstanceValue",
					"VPCId":                 "CreateServerlessDBInstanceValue",
				},
			},
		},
		"SecurityIpGroups": map[string]interface{}{
			"SecurityIpGroup": []interface{}{
				map[string]interface{}{
					"SecurityIpGroupAttribute": "CreateServerlessDBInstanceValue",
					"SecurityIpGroupName":      "CreateServerlessDBInstanceValue",
					"SecurityIpList":           "CreateServerlessDBInstanceValue",
				},
			},
		},
		"TagResources": map[string]interface{}{
			"TagResource": []interface{}{
				map[string]interface{}{
					"TagKey":   "Created",
					"TagValue": "CreateServerlessDBInstanceValue",
				},
				map[string]interface{}{
					"TagKey":   "For",
					"TagValue": "CreateServerlessDBInstanceValue",
				},
			},
		},
		"DBInstanceId": "CreateServerlessDBInstanceValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateServerlessDBInstance
		"DBInstanceId": "CreateServerlessDBInstanceValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_mongodb_serverless_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudMongodbServerlessInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"DBInstanceId": "CreateServerlessDBInstanceValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateServerlessDBInstance" {
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
		err := resourceAlicloudMongodbServerlessInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudMongodbServerlessInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// ModifyDBInstanceDescription
	attributesDiff := map[string]interface{}{
		"db_instance_description": "ModifyDBInstanceDescriptionValue",
	}
	diff, err := newInstanceDiff("alicloud_mongodb_serverless_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// DescribeDBInstanceAttribute Response
		"DBInstances": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"DBInstanceDescription": "ModifyDBInstanceDescriptionValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBInstanceDescription" {
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
		err := resourceAlicloudMongodbServerlessInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), nil)
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
		"tags": map[string]string{
			"Created": "UpdateValue",
			"For":     "UpdateValue",
		},
	}
	diff, err = newInstanceDiff("alicloud_mongodb_serverless_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"TagResources": map[string]interface{}{
			"TagResource": []interface{}{
				map[string]interface{}{
					"TagKey":   "Created",
					"TagValue": "UpdateValue",
				},
				map[string]interface{}{
					"TagKey":   "For",
					"TagValue": "UpdateValue",
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
		err := resourceAlicloudMongodbServerlessInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyDBInstanceMaintainTime
	attributesDiff = map[string]interface{}{
		"maintain_start_time": "ModifyDBInstanceMaintainTimeValue",
		"maintain_end_time":   "ModifyDBInstanceMaintainTimeValue",
	}
	diff, err = newInstanceDiff("alicloud_mongodb_serverless_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"DBInstances": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"MaintainEndTime":   "ModifyDBInstanceMaintainTimeValue",
					"MaintainStartTime": "ModifyDBInstanceMaintainTimeValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBInstanceMaintainTime" {
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
		err := resourceAlicloudMongodbServerlessInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifyDBInstanceMaintainTime
	attributesDiff = map[string]interface{}{
		"db_instance_storage": 10,
		"capacity_unit":       200,
	}
	diff, err = newInstanceDiff("alicloud_mongodb_serverless_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"DBInstances": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"CapacityUnit":      200,
					"DBInstanceStorage": 10,
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifyDBInstanceSpec" {
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
		err := resourceAlicloudMongodbServerlessInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ModifySecurityIps
	attributesDiff = map[string]interface{}{
		"security_ip_groups": []map[string]interface{}{
			{
				"security_ip_group_attribute": "ModifySecurityIpsValue",
				"security_ip_group_name":      "ModifySecurityIpsValue",
				"security_ip_list":            "ModifySecurityIpsValue",
			},
			{
				"security_ip_group_attribute": "ModifySecurityIpsValue",
				"security_ip_group_name":      "ModifySecurityIpsValue",
				"security_ip_list":            "ModifySecurityIpsValue",
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_mongodb_serverless_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"SecurityIpGroups": map[string]interface{}{
			"SecurityIpGroup": []interface{}{
				map[string]interface{}{
					"SecurityIpGroupAttribute": "ModifySecurityIpsValue",
					"SecurityIpGroupName":      "ModifySecurityIpsValue",
					"SecurityIpList":           "ModifySecurityIpsValue",
				},
				map[string]interface{}{
					"SecurityIpGroupAttribute": "ModifySecurityIpsValue",
					"SecurityIpGroupName":      "ModifySecurityIpsValue",
					"SecurityIpList":           "ModifySecurityIpsValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ModifySecurityIps" {
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
		err := resourceAlicloudMongodbServerlessInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_mongodb_serverless_instance"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DescribeDBInstanceAttribute" {
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
		err := resourceAlicloudMongodbServerlessInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	err = resourceAlicloudMongodbServerlessInstanceDelete(dExisted, rawClient)
	assert.Nil(t, err)

}
