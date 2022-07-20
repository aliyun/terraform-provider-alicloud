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
	resource.AddTestSweepers("alicloud_config_rule", &resource.Sweeper{
		Name: "alicloud_config_rule",
		F:    testSweepConfigRule,
	})
}

func testSweepConfigRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-test",
	}

	request := make(map[string]interface{})
	var response map[string]interface{}
	action := "ListConfigRules"
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var ruleIds []string
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-01-08"), StringPointer("AK"), request, nil, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"Throttling.User"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve config rule in service list: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.ConfigRules.ConfigRuleList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ConfigRules.ConfigRuleList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			ruleIds = append(ruleIds, item["ConfigRuleName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	for _, ruleId := range ruleIds {
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(ruleId), strings.ToLower(prefix)) {
				skip = false
			}
		}
		if skip {
			log.Printf("[INFO] Skipping config rule: %s ", ruleId)
			continue
		}
		action = "DeleteConfigRules"
		request := map[string]interface{}{
			"ConfigRuleIds": ruleId,
		}
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve config rule (%s): %s", ruleId, err)
			continue
		}
		log.Printf("[INFO] Delete config rule success: %s ", ruleId)
	}
	return nil
}

func TestAccAlicloudConfigRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_rule.default"
	ra := resourceAttrInit(resourceId, ConfigRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigRuleBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name":                 name,
					"risk_level":                "1",
					"resource_types_scope":      []string{"ACS::ECS::Instance"},
					"config_rule_trigger_types": "ConfigurationItemChangeNotification",
					"source_identifier":         "ecs-instances-in-vpc",
					"source_owner":              "ALIYUN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":                 name,
						"risk_level":                "1",
						"resource_types_scope.#":    "1",
						"config_rule_trigger_types": "ConfigurationItemChangeNotification",
						"source_identifier":         "ecs-instances-in-vpc",
						"source_owner":              "ALIYUN",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_rule",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_rule",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_parameters": map[string]string{
						"vpcIds": "${data.alicloud_instances.default.instances[0].vpc_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_parameters.%": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_ids_scope": "${data.alicloud_instances.default.instances[0].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_ids_scope": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"region_ids_scope": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_ids_scope": "cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_ids_scope": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_ids_scope": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_key_scope": "tfTest123Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_key_scope": "tfTest123Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_value_scope": "tfTest 123 Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_value_scope": "tfTest 123 Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name":                  name,
					"risk_level":                 "1",
					"source_identifier":          "ecs-instances-in-vpc",
					"source_owner":               "ALIYUN",
					"region_ids_scope":           "cn-beijing",
					"resource_group_ids_scope":   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"exclude_resource_ids_scope": "${data.alicloud_instances.default.instances[1].id}",
					"tag_key_scope":              "tfTest123",
					"tag_value_scope":            "tfTest 123 Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":                  name,
						"risk_level":                 "1",
						"source_identifier":          "ecs-instances-in-vpc",
						"source_owner":               "ALIYUN",
						"region_ids_scope":           "cn-beijing",
						"resource_group_ids_scope":   CHECKSET,
						"exclude_resource_ids_scope": CHECKSET,
						"tag_key_scope":              "tfTest123",
						"tag_value_scope":            "tfTest 123 Update",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudConfigRule_status(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_rule.default"
	ra := resourceAttrInit(resourceId, ConfigRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigRuleBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name":                 name,
					"risk_level":                "1",
					"resource_types_scope":      []string{"ACS::ECS::Instance"},
					"config_rule_trigger_types": "ConfigurationItemChangeNotification",
					"source_identifier":         "ecs-instances-in-vpc",
					"source_owner":              "ALIYUN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":                 name,
						"risk_level":                "1",
						"resource_types_scope.#":    "1",
						"config_rule_trigger_types": "ConfigurationItemChangeNotification",
						"source_identifier":         "ecs-instances-in-vpc",
						"source_owner":              "ALIYUN",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "INACTIVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "INACTIVE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ACTIVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ACTIVE",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

var ConfigRuleMap = map[string]string{}

func ConfigRuleBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_instances" "default"{}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

`, name)
}

// Test this case need use a custom `source_identifier`
func SkipTestAccAlicloudConfigRule_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_rule.default"
	ra := resourceAttrInit(resourceId, ConfigRuleMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccConfigRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ConfigRuleBasicdependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name":                 name,
					"risk_level":                "1",
					"resource_types_scope":      []string{"ACS::ECS::Instance"},
					"config_rule_trigger_types": "ConfigurationItemChangeNotification",
					"source_identifier":         "acs:fc:cn-shanghai:1009318965****:services/customer-demo.LATEST/functions/ApprovedAimsByIdPython",
					"source_owner":              "CUSTOM_FC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":                 name,
						"risk_level":                "1",
						"resource_types_scope.#":    "1",
						"config_rule_trigger_types": "ConfigurationItemChangeNotification",
						"source_identifier":         "acs:fc:cn-shanghai:1009318965****:services/customer-demo.LATEST/functions/ApprovedAimsByIdPython",
						"source_owner":              "CUSTOM_FC",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test_rule",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test_rule",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_parameters": map[string]string{
						"vpcIds": "${data.alicloud_instances.default.instances[0].vpc_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_parameters.%": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_ids_scope": "${data.alicloud_instances.default.instances[0].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_ids_scope": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_types_scope": []string{"ACS::OSS::Bucket"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_types_scope.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_rule_trigger_types": "ScheduledNotification",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_trigger_types": "ScheduledNotification",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maximum_execution_frequency": "One_Hour",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maximum_execution_frequency": "One_Hour",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"region_ids_scope": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_ids_scope": "cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_ids_scope": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_ids_scope": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_key_scope": "tfTest123Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_key_scope": "tfTest123Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_value_scope": "tfTest 123 Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_value_scope": "tfTest 123 Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level":                "1",
					"config_rule_trigger_types": "ConfigurationItemChangeNotification",
					"resource_types_scope":      []string{"ACS::ECS::Instance"},
					"region_ids_scope":          "cn-beijing",
					"resource_group_ids_scope":  "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tag_key_scope":             "tfTest123",
					"tag_value_scope":           "tfTest 123 Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level":                "1",
						"config_rule_trigger_types": "ConfigurationItemChangeNotification",
						"resource_types_scope.#":    "1",
						"region_ids_scope":          "cn-beijing",
						"resource_group_ids_scope":  CHECKSET,
						"tag_key_scope":             "tfTest123",
						"tag_value_scope":           "tfTest 123 Update",
					}),
				),
			},
		},
	})
}

var ConfigRuleMap1 = map[string]string{}

func ConfigRuleBasicdependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_instances" "default"{}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

`, name)
}

func TestUnitAlicloudConfigRule(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"config_rule_trigger_types":  "CreateConfigRuleValue",
		"description":                "CreateConfigRuleValue",
		"exclude_resource_ids_scope": "CreateConfigRuleValue",
		"input_parameters": map[string]interface{}{
			"vpcIds": "CreateConfigRuleValue",
		},
		"maximum_execution_frequency": "CreateConfigRuleValue",
		"region_ids_scope":            "CreateConfigRuleValue",
		"resource_group_ids_scope":    "CreateConfigRuleValue",
		"resource_types_scope":        []interface{}{"CreateConfigRuleValue0", "CreateConfigRuleValue1"},
		"risk_level":                  10,
		"rule_name":                   "CreateConfigRuleValue",
		"source_identifier":           "CreateConfigRuleValue",
		"source_owner":                "CreateConfigRuleValue",
		"tag_key_scope":               "CreateConfigRuleValue",
		"tag_value_scope":             "CreateConfigRuleValue",
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
		// GetConfigRule
		"ConfigRule": map[string]interface{}{
			"ConfigRuleId":            "CreateConfigRuleValue",
			"Description":             "CreateConfigRuleValue",
			"ExcludeResourceIdsScope": "CreateConfigRuleValue",
			"InputParameters": map[string]interface{}{
				"vpcIds": "CreateConfigRuleValue",
			},
			"MaximumExecutionFrequency": "CreateConfigRuleValue",
			"RegionIdsScope":            "CreateConfigRuleValue",
			"ResourceGroupIdsScope":     "CreateConfigRuleValue",
			"Scope": map[string]interface{}{
				"ComplianceResourceTypes": "CreateConfigRuleValue",
			},
			"Source": map[string]interface{}{
				"Owner":      "CreateConfigRuleValue",
				"Identifier": "CreateConfigRuleValue",
				"SourceDetails": []interface{}{
					map[string]interface{}{
						"MessageType": "CreateConfigRuleValue",
					},
				},
			},
			"RiskLevel":       10,
			"ConfigRuleName":  "CreateConfigRuleValue",
			"ConfigRuleState": "ACTIVE",
			"TagKeyScope":     "CreateConfigRuleValue",
			"TagValueScope":   "CreateConfigRuleValue",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateConfigRule
		"ConfigRuleId": "CreateConfigRuleValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_config_rule", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigRuleCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleId": "CreateConfigRuleValue",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateConfigRule" {
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
		err := resourceAlicloudConfigRuleCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigRuleUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateConfigRule
	attributesDiff := map[string]interface{}{
		"config_rule_trigger_types":  "UpdateConfigRuleValue",
		"description":                "UpdateConfigRuleValue",
		"resource_types_scope":       []interface{}{"UpdateConfigRuleValue"},
		"risk_level":                 15,
		"exclude_resource_ids_scope": "UpdateConfigRuleValue",
		"input_parameters": map[string]interface{}{
			"vpcIds": "UpdateConfigRuleValue",
		},
		"maximum_execution_frequency": "UpdateConfigRuleValue",
		"region_ids_scope":            "UpdateConfigRuleValue",
		"resource_group_ids_scope":    "UpdateConfigRuleValue",
		"tag_key_scope":               "UpdateConfigRuleValue",
		"tag_value_scope":             "UpdateConfigRuleValue",
	}
	diff, err := newInstanceDiff("alicloud_config_rule", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetConfigRule Response
		"ConfigRule": map[string]interface{}{
			"Source": map[string]interface{}{
				"SourceDetails": []interface{}{
					map[string]interface{}{
						"MessageType": "UpdateConfigRuleValue",
					},
				},
			},
			"Description":             "UpdateConfigRuleValue",
			"ResourceTypesScope":      []interface{}{"UpdateConfigRuleValue"},
			"RiskLevel":               15,
			"ExcludeResourceIdsScope": "UpdateConfigRuleValue",
			"InputParameters": map[string]interface{}{
				"vpcIds": "UpdateConfigRuleValue",
			},
			"MaximumExecutionFrequency": "UpdateConfigRuleValue",
			"RegionIdsScope":            "UpdateConfigRuleValue",
			"ResourceGroupIdsScope":     "UpdateConfigRuleValue",
			"TagKeyScope":               "UpdateConfigRuleValue",
			"TagValueScope":             "UpdateConfigRuleValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateConfigRule" {
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
		err := resourceAlicloudConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// StopConfigRules
	attributesDiff = map[string]interface{}{
		"status": "INACTIVE",
	}
	diff, err = newInstanceDiff("alicloud_config_rule", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleState": "INACTIVE",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "StopConfigRules" {
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
		err := resourceAlicloudConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()

		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// ActiveConfigRule
	attributesDiff = map[string]interface{}{
		"status": "ACTIVE",
	}
	diff, err = newInstanceDiff("alicloud_config_rule", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleState": "ACTIVE",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ActiveConfigRules" {
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
		err := resourceAlicloudConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_rule"].Schema).Data(dExisted.State(), nil)
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
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "ConfigRuleNotExists", "Invalid.ConfigRuleId.Value", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetConfigRule" {
				switch errorCode {
				case "{}", "ConfigRuleNotExists", "Invalid.ConfigRuleId.Value":
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
		err := resourceAlicloudConfigRuleRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}", "ConfigRuleNotExists", "Invalid.ConfigRuleId.Value":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigRuleDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "ConfigRuleNotExists"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteConfigRules" {
				switch errorCode {
				case "NonRetryableError", "ConfigRuleNotExists":
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
		err := resourceAlicloudConfigRuleDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "ConfigRuleNotExists":
			assert.Nil(t, err)
		}
	}

}
