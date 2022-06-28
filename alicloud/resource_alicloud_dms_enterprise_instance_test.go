package alicloud

import (
	"encoding/json"
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
	resource.AddTestSweepers("alicloud_dms_enterprise_instance", &resource.Sweeper{
		Name: "alicloud_dms_enterprise_instance",
		F:    testSweepDMSEnterpriseInstances,
	})
}

func testSweepDMSEnterpriseInstances(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"testacc",
	}
	request := map[string]interface{}{
		"InstanceState": "NORMAL",
		"PageSize":      PageSizeXLarge,
		"PageNumber":    1,
	}
	var response map[string]interface{}
	action := "ListInstances"
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_instances", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.InstanceList.Instance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceList.Instance", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			id := item["Host"].(string) + ":" + item["Port"].(json.Number).String()

			skip := true

			for _, prefix := range prefixes {
				if item["InstanceAlias"] != nil {
					if strings.HasPrefix(strings.ToLower(item["InstanceAlias"].(string)), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			}
			if skip || item["InstanceAlias"] == nil {
				log.Printf("[INFO] Skipping DMS Enterprise Instances: %s", id)
				continue
			}
			action := "DeleteInstance"
			request := map[string]interface{}{
				"Host": item["Host"],
				"Port": item["Port"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DMS Enterprise Instance (%s (%s)): %s", item["InstanceAlias"].(string), id, err)
				continue
			}
			log.Printf("[INFO] Delete DMS Enterprise Instance Success: %s ", item["InstanceAlias"].(string))
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudDMSEnterprise(t *testing.T) {
	resourceId := "alicloud_dms_enterprise_instance.default"
	var v map[string]interface{}
	ra := resourceAttrInit(resourceId, testAccCheckKeyValueInMapsForDMS)

	serviceFunc := func() interface{} {
		return &Dms_enterpriseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccDmsEnterpriseInstance%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDmsConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dba_uid":           "${tonumber(data.alicloud_account.current.id)}",
					"host":              "${alicloud_db_instance.instance.connection_string}",
					"port":              "3306",
					"network_type":      "VPC",
					"safe_rule":         "自由操作",
					"tid":               "${data.alicloud_dms_user_tenants.default.ids.0}",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "test",
					"database_user":     "${alicloud_db_account.account.name}",
					"database_password": "${alicloud_db_account.account.password}",
					"instance_alias":    name,
					"query_timeout":     "70",
					"export_timeout":    "2000",
					"ecs_region":        os.Getenv("ALICLOUD_REGION"),
					"ddl_online":        "0",
					"use_dsql":          "0",
					"data_link_name":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dba_uid":         CHECKSET,
						"host":            CHECKSET,
						"port":            "3306",
						"network_type":    "VPC",
						"safe_rule":       "自由操作",
						"tid":             CHECKSET,
						"instance_type":   "mysql",
						"instance_source": "RDS",
						"env_type":        "test",
						"database_user":   CHECKSET,
						"instance_alias":  name,
						"query_timeout":   "70",
						"export_timeout":  "2000",
						"ecs_region":      os.Getenv("ALICLOUD_REGION"),
						"ddl_online":      "0",
						"use_dsql":        "0",
						"data_link_name":  "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_password", "dba_uid", "network_type", "port", "safe_rule", "tid"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"env_type": "dev",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"env_type": "dev",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_timeout": "77",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_timeout": "77",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"dba_uid":           "${tonumber(data.alicloud_account.current.id)}",
					"host":              "${alicloud_db_instance.instance.connection_string}",
					"port":              "3306",
					"network_type":      "VPC",
					"safe_rule":         "自由操作",
					"tid":               "${data.alicloud_dms_user_tenants.default.ids.0}",
					"instance_type":     "mysql",
					"instance_source":   "RDS",
					"env_type":          "test",
					"database_user":     "${alicloud_db_account.account.name}",
					"database_password": "${alicloud_db_account.account.password}",
					"instance_alias":    name,
					"query_timeout":     "70",
					"export_timeout":    "2000",
					"ecs_region":        os.Getenv("ALICLOUD_REGION"),
					"ddl_online":        "0",
					"use_dsql":          "0",
					"data_link_name":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dba_uid":         CHECKSET,
						"host":            CHECKSET,
						"port":            "3306",
						"network_type":    "VPC",
						"safe_rule":       "自由操作",
						"tid":             CHECKSET,
						"instance_type":   "mysql",
						"instance_source": "RDS",
						"env_type":        "test",
						"database_user":   CHECKSET,
						"instance_alias":  name,
						"query_timeout":   "70",
						"export_timeout":  "2000",
						"ecs_region":      os.Getenv("ALICLOUD_REGION"),
						"ddl_online":      "0",
						"use_dsql":        "0",
						"data_link_name":  "",
					}),
				),
			},
		},
	})
}

var testAccCheckKeyValueInMapsForDMS = map[string]string{}

func resourceDmsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_account" "current" {
	}
	data "alicloud_dms_user_tenants" "default" {
		status = "ACTIVE"
	}
	
	data "alicloud_db_zones" "default"{
		engine = "MySQL"
		engine_version = "8.0"
		instance_charge_type = "PostPaid"
		category = "HighAvailability"
		db_instance_storage_type = "cloud_essd"
	}
	
	data "alicloud_db_instance_classes" "default" {
		zone_id = data.alicloud_db_zones.default.zones.0.id
		engine = "MySQL"
		engine_version = "8.0"
		category = "HighAvailability"
		db_instance_storage_type = "cloud_essd"
		instance_charge_type = "PostPaid"
	}
	
	data "alicloud_vpcs" "default" {
	 name_regex = "^default-NODELETING"
	}
	data "alicloud_vswitches" "default" {
	  vpc_id = data.alicloud_vpcs.default.ids.0
	  zone_id = data.alicloud_db_zones.default.zones.0.id
	}
	
	resource "alicloud_security_group" "default" {
		name = var.name
		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	
	resource "alicloud_db_instance" "instance" {
		engine = "MySQL"
		engine_version = "8.0"
		db_instance_storage_type = "cloud_essd"
		instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
		instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
		vswitch_id       = data.alicloud_vswitches.default.ids.0
		instance_name    = var.name
		security_ips     = ["100.104.5.0/24","192.168.0.6"]
		tags = {
			"key1" = "value1"
			"key2" = "value2"
		}
	}
	
	resource "alicloud_db_account" "account" {
	instance_id = "${alicloud_db_instance.instance.id}"
	name        = "tftestnormal"
	password    = "Test12345"
	type        = "Normal"
	}`, name)
}

func TestUnitAlicloudDMSEnterprise(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_dms_enterprise_instance"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_dms_enterprise_instance"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"dba_uid":           1,
		"host":              "RegisterInstanceValue",
		"port":              3306,
		"network_type":      "RegisterInstanceValue",
		"safe_rule":         "RegisterInstanceValue",
		"tid":               1,
		"instance_type":     "RegisterInstanceValue",
		"instance_source":   "RegisterInstanceValue",
		"env_type":          "RegisterInstanceValue",
		"database_user":     "RegisterInstanceValue",
		"database_password": "RegisterInstanceValue",
		"instance_alias":    "RegisterInstanceValue",
		"query_timeout":     70,
		"export_timeout":    2000,
		"ecs_region":        "RegisterInstanceValue",
		"ddl_online":        1,
		"use_dsql":          1,
		"data_link_name":    "RegisterInstanceValue",
		"ecs_instance_id":   "RegisterInstanceValue",
		"skip_test":         false,
		"vpc_id":            "RegisterInstanceValue",
		"sid":               "RegisterInstanceValue",
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
		// GetInstance
		"Instance": map[string]interface{}{
			"Host":           "RegisterInstanceValue",
			"Port":           3306,
			"DataLinkName":   "RegisterInstanceValue",
			"DatabaseUser":   "RegisterInstanceValue",
			"DbaId":          "RegisterInstanceValue",
			"DdlOnline":      1,
			"EcsInstanceId":  "RegisterInstanceValue",
			"EcsRegion":      "RegisterInstanceValue",
			"EnvType":        "RegisterInstanceValue",
			"ExportTimeout":  2000,
			"InstanceId":     "RegisterInstanceValue",
			"InstanceAlias":  "RegisterInstanceValue",
			"InstanceSource": "RegisterInstanceValue",
			"InstanceType":   "RegisterInstanceValue",
			"QueryTimeout":   70,
			"SafeRuleId":     "RegisterInstanceValue",
			"Sid":            "RegisterInstanceValue",
			"State":          "RegisterInstanceValue",
			"UseDsql":        1,
			"VpcId":          "RegisterInstanceValue",
			"DbaNickName":    "RegisterInstanceValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// RegisterInstance
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_dms_enterprise_instance", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmsenterpriseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudDmsEnterpriseInstanceCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetInstance Response
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "RegisterInstanceFailure", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "RegisterInstance" {
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
		err := resourceAlicloudDmsEnterpriseInstanceCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dms_enterprise_instance"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmsenterpriseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudDmsEnterpriseInstanceUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateInstance
	attributesDiff := map[string]interface{}{
		"dba_uid":           1,
		"dba_id":            "UpdateInstanceValue",
		"network_type":      "UpdateInstanceValue",
		"safe_rule":         "UpdateInstanceValue",
		"tid":               1,
		"instance_type":     "UpdateInstanceValue",
		"instance_source":   "UpdateInstanceValue",
		"env_type":          "UpdateInstanceValue",
		"database_user":     "UpdateInstanceValue",
		"database_password": "UpdateInstanceValue",
		"instance_alias":    "UpdateInstanceValue",
		"query_timeout":     80,
		"export_timeout":    1800,
		"ecs_region":        "UpdateInstanceValue",
		"ddl_online":        2,
		"use_dsql":          1,
		"data_link_name":    "UpdateInstanceValue",
		"instance_id":       "UpdateInstanceValue",
		"instance_name":     "UpdateInstanceValue",
		"safe_rule_id":      "UpdateInstanceValue",
		"ecs_instance_id":   "UpdateInstanceValue",
		"sid":               "UpdateInstanceValue",
		"vpc_id":            "UpdateInstanceValue",
	}
	diff, err := newInstanceDiff("alicloud_dms_enterprise_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dms_enterprise_instance"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetInstance Response
		"Instance": map[string]interface{}{
			"DbaId":          "UpdateInstanceValue",
			"SafeRuleId":     "UpdateInstanceValue",
			"InstanceType":   "UpdateInstanceValue",
			"InstanceSource": "UpdateInstanceValue",
			"EnvType":        "UpdateInstanceValue",
			"DatabaseUser":   "UpdateInstanceValue",
			"InstanceAlias":  "UpdateInstanceValue",
			"QueryTimeout":   80,
			"ExportTimeout":  1800,
			"EcsRegion":      "UpdateInstanceValue",
			"DdlOnline":      2,
			"UseDsql":        1,
			"DataLinkName":   "UpdateInstanceValue",
			"InstanceId":     "UpdateInstanceValue",
			"VpcId":          "UpdateInstanceValue",
			"Sid":            "UpdateInstanceValue",
			"EcsInstanceId":  "UpdateInstanceValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateInstance" {
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
		err := resourceAlicloudDmsEnterpriseInstanceUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_dms_enterprise_instance"].Schema).Data(dExisted.State(), nil)
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
	diff, err = newInstanceDiff("alicloud_dms_enterprise_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dms_enterprise_instance"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetInstance" {
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
		err := resourceAlicloudDmsEnterpriseInstanceRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewDmsenterpriseClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudDmsEnterpriseInstanceDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_dms_enterprise_instance", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_dms_enterprise_instance"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "InstanceNoEnoughNumber"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteInstance" {
				switch errorCode {
				case "NonRetryableError", "InstanceNoEnoughNumber":
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
		err := resourceAlicloudDmsEnterpriseInstanceDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "InstanceNoEnoughNumber":
			assert.Nil(t, err)
		}
	}

}
