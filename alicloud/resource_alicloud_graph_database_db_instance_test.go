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
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_graph_database_db_instance",
		&resource.Sweeper{
			Name: "alicloud_graph_database_db_instance",
			F:    testSweepGraphDatabaseDbInstance,
		})
}

func testSweepGraphDatabaseDbInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeDBInstances"
	request := map[string]interface{}{}

	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("gdb", "2019-09-03", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Items.DBInstance", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["DBInstanceDescription"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if item["DBInstanceStatus"].(string) != "Running" {
					skip = true
				}
				if skip {
					log.Printf("[INFO] Skipping Graph Database DbInstance: %s", item["DBInstanceDescription"].(string))
					continue
				}
			}
			action := "DeleteDBInstance"
			deleteRequest := map[string]interface{}{
				"DBInstanceId": item["DBInstanceId"],
			}
			_, err = client.RpcPost("gdb", "2019-09-03", action, nil, deleteRequest, true)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Graph Database DbInstance (%s): %s", item["DBInstanceDescription"].(string), err)
			}
			log.Printf("[INFO] Delete Graph Database DbInstance success: %s ", item["DBInstanceDescription"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudGraphDatabaseDbInstance_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_graph_database_db_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudGraphDatabaseDbInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGraphDatabaseDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgraphdatabasedbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGraphDatabaseDbInstanceBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.GraphDatabaseSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class":            "gdb.r.xlarge",
					"db_instance_network_type": "vpc",
					"db_version":               "1.0",
					"db_instance_category":     "HA",
					"db_instance_storage_type": "cloud_essd",
					"db_node_storage":          "50",
					"payment_type":             "PayAsYouGo",
					"db_instance_description":  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":            "gdb.r.xlarge",
						"db_instance_network_type": "vpc",
						"db_version":               "1.0",
						"db_instance_category":     "HA",
						"db_instance_storage_type": "cloud_essd",
						"db_node_storage":          "50",
						"payment_type":             "PayAsYouGo",
						"db_instance_description":  name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_ip_array": []map[string]interface{}{
						{
							"db_instance_ip_array_name": "default",
							"security_ips":              "127.0.0.2",
						},
						{
							"db_instance_ip_array_name": "tftest",
							"security_ips":              "192.168.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name + "_update",
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
					"db_node_storage": "80",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_storage": "80",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class": "gdb.r.2xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class": "gdb.r.2xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "${var.name}",
					"db_node_storage":         "100",
					"db_node_class":           "gdb.r.xlarge",
					"db_instance_ip_array": []map[string]interface{}{
						{
							"db_instance_ip_array_name": "default",
							"security_ips":              "127.0.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name,
						"db_node_storage":         "100",
						"db_node_class":           "gdb.r.xlarge",
						"db_instance_ip_array.#":  "1",
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

func TestAccAlicloudGraphDatabaseDbInstance_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_graph_database_db_instance.default"
	checkoutSupportedRegions(t, true, connectivity.GraphDatabaseDbInstanceSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGraphDatabaseDbInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGraphDatabaseDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgraphdatabasedbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGraphDatabaseDbInstanceBasicDependence1)
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
					"db_node_class":            "gdb.r.xlarge",
					"db_instance_network_type": "vpc",
					"db_version":               "1.0",
					"db_instance_category":     "HA",
					"db_instance_storage_type": "cloud_essd",
					"db_node_storage":          "50",
					"payment_type":             "PayAsYouGo",
					"db_instance_description":  "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":                  "${local.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":            "gdb.r.xlarge",
						"db_instance_network_type": "vpc",
						"db_version":               "1.0",
						"db_instance_category":     "HA",
						"db_instance_storage_type": "cloud_essd",
						"db_node_storage":          "50",
						"payment_type":             "PayAsYouGo",
						"db_instance_description":  name,
						"vswitch_id":               CHECKSET,
						"vpc_id":                   CHECKSET,
						"zone_id":                  CHECKSET,
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

func TestAccAlicloudGraphDatabaseDbInstance_single(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_graph_database_db_instance.default"
	checkoutSupportedRegions(t, true, connectivity.GraphDatabaseDbInstanceSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGraphDatabaseDbInstanceMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGraphDatabaseDbInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgraphdatabasedbinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGraphDatabaseDbInstanceBasicDependence1)
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
					"db_node_class":            "gdb.r.xlarge_basic",
					"db_instance_network_type": "vpc",
					"db_version":               "1.0",
					"db_instance_category":     "SINGLE",
					"db_instance_storage_type": "cloud_essd",
					"db_node_storage":          "50",
					"payment_type":             "PayAsYouGo",
					"db_instance_description":  "${var.name}",
					"vswitch_id":               "${data.alicloud_vswitches.default.ids.0}",
					"vpc_id":                   "${data.alicloud_vpcs.default.ids.0}",
					"zone_id":                  "${local.zone_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":            "gdb.r.xlarge_basic",
						"db_instance_network_type": "vpc",
						"db_version":               "1.0",
						"db_instance_category":     "SINGLE",
						"db_instance_storage_type": "cloud_essd",
						"db_node_storage":          "50",
						"payment_type":             "PayAsYouGo",
						"db_instance_description":  name,
						"vswitch_id":               CHECKSET,
						"vpc_id":                   CHECKSET,
						"zone_id":                  CHECKSET,
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

var AlicloudGraphDatabaseDbInstanceMap0 = map[string]string{
	"connection_string": CHECKSET,
	"port":              CHECKSET,
}

func AlicloudGraphDatabaseDbInstanceBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}`, name)
}

// The instance requires the same zone as the vswitch, but currently the instance does not support zone query.
func AlicloudGraphDatabaseDbInstanceBasicDependence1(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
locals {
  zone_id = "cn-hangzhou-h"
}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = local.zone_id
}`, name)
}

func TestUnitAlicloudGraphDatabaseDbInstance(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"db_node_class":            "CreateDBInstanceValue",
		"db_instance_network_type": "CreateDBInstanceValue",
		"db_version":               "CreateDBInstanceValue",
		"db_instance_category":     "CreateDBInstanceValue",
		"db_instance_storage_type": "CreateDBInstanceValue",
		"db_node_storage":          50,
		"payment_type":             "PayAsYouGo",
		"db_instance_description":  "CreateDBInstanceValue",
		"vswitch_id":               "CreateDBInstanceValue",
		"vpc_id":                   "CreateDBInstanceValue",
		"zone_id":                  "CreateDBInstanceValue",
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
		"Items": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"Category":              "CreateDBInstanceValue",
					"DBInstanceDescription": "CreateDBInstanceValue",
					"DBInstanceNetworkType": "CreateDBInstanceValue",
					"DBInstanceStorageType": "CreateDBInstanceValue",
					"DBNodeClass":           "CreateDBInstanceValue",
					"DBNodeStorage":         50,
					"DBVersion":             "CreateDBInstanceValue",
					"PayType":               "Postpaid",
					"DBInstanceStatus":      "Running",
					"VSwitchId":             "CreateDBInstanceValue",
					"ZoneId":                "CreateDBInstanceValue",
					"VpcId":                 "CreateDBInstanceValue",
					"DBInstanceId":          "CreateDBInstanceValue",
				},
			},
			"DBInstanceIPArray": []interface{}{
				map[string]interface{}{
					"SecurityIps":                "CreateDBInstanceValue",
					"DBInstanceIPArrayName":      "CreateDBInstanceValue",
					"DBInstanceIPArrayAttribute": "CreateDBInstanceValue",
				},
			},
		},
	}
	CreateMockResponse := map[string]interface{}{
		"DBInstanceId": "CreateDBInstanceValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_graph_database_db_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGraphDatabaseDbInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		"DBInstanceId": "CreateDBInstanceValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateDBInstance" {
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
		err := resourceAlicloudGraphDatabaseDbInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGraphDatabaseDbInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff := map[string]interface{}{
		"db_instance_description": "ModifyDBInstanceDescriptionValue",
	}
	diff, err := newInstanceDiff("alicloud_graph_database_db_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Items": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"DBInstanceDescription": "ModifyDBInstanceDescriptionValue",
					"DBInstanceStatus":      "Running",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
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
		err := resourceAlicloudGraphDatabaseDbInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	attributesDiff = map[string]interface{}{
		"db_instance_ip_array": []map[string]interface{}{
			{
				"db_instance_ip_array_name":      "ModifyDBInstanceAccessWhiteList",
				"security_ips":                   "ModifyDBInstanceAccessWhiteList",
				"db_instance_ip_array_attribute": "ModifyDBInstanceAccessWhiteList",
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_graph_database_db_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Items": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"DBInstanceIPArray": []interface{}{
						map[string]interface{}{
							"SecurityIps":                "ModifyDBInstanceAccessWhiteList",
							"DBInstanceIPArrayName":      "ModifyDBInstanceAccessWhiteList",
							"DBInstanceIPArrayAttribute": "ModifyDBInstanceAccessWhiteList",
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
			if *action == "ModifyDBInstanceAccessWhiteList" {
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
		err := resourceAlicloudGraphDatabaseDbInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	attributesDiff = map[string]interface{}{
		"db_node_class":            "ModifyDBInstanceSpecValue",
		"db_node_storage":          70,
		"db_instance_storage_type": "ModifyDBInstanceSpecValue",
	}
	diff, err = newInstanceDiff("alicloud_graph_database_db_instance", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		"Items": map[string]interface{}{
			"DBInstance": []interface{}{
				map[string]interface{}{
					"DBInstanceStorageType": "ModifyDBInstanceSpecValue",
					"DBNodeClass":           "ModifyDBInstanceSpecValue",
					"DBNodeStorage":         70,
					"DBInstanceStatus":      "Running",
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
		err := resourceAlicloudGraphDatabaseDbInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_graph_database_db_instance"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
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
		err := resourceAlicloudGraphDatabaseDbInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewGdsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudGraphDatabaseDbInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteDBInstance" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					if errorCodes[retryIndex] == "InvalidDBInstance.NotFound" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			if *action == "DescribeDBInstanceAttribute" {
				return notFoundResponseMock("{}")
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudGraphDatabaseDbInstanceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
