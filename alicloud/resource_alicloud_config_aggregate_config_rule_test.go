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
	resource.AddTestSweepers("alicloud_config_aggregate_config_rule", &resource.Sweeper{
		Name: "alicloud_config_aggregate_config_rule",
		F:    testSweepConfigAggregateConfigRule,
	})
}

func testSweepConfigAggregateConfigRule(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}

	// Get all AggregatorId
	aggregatorIds := make([]string, 0)
	action := "ListAggregators"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-09-07"), StringPointer("AK"), request, nil, &runtime)
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
			log.Println("List Config Aggregator Failed!", err)
			return nil
		}
		resp, err := jsonpath.Get("$.AggregatorsResult.Aggregators", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AggregatorsResult.Aggregators", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["AggregatorName"])), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Config Aggregator: %v (%v)", item["AggregatorName"], item["AggregatorId"])
				continue
			}
			aggregatorIds = append(aggregatorIds, fmt.Sprint(item["AggregatorId"]))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	for _, aggregatorId := range aggregatorIds {
		configRuleIds := make([]string, 0)
		action := "ListAggregateConfigRules"
		var response map[string]interface{}
		request := map[string]interface{}{
			"AggregatorId": aggregatorId,
			"PageSize":     PageSizeLarge,
			"PageNumber":   1,
		}
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-09-07"), StringPointer("AK"), request, nil, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			resp, err := jsonpath.Get("$.ConfigRules.ConfigRuleList", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ConfigRules.ConfigRuleList", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				skip := true
				item := v.(map[string]interface{})
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["ConfigRuleName"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Aggregate Config Rule: %v (%v)", item["ConfigRuleName"], item["ConfigRuleId"])
					continue
				}
				configRuleIds = append(configRuleIds, fmt.Sprint(item["ConfigRuleId"]))
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}

		if len(configRuleIds) > 0 {
			log.Printf("[INFO] Deleting Aggregate Config Rules: (%s)", strings.Join(configRuleIds, ","))
			action = "DeleteAggregateConfigRules"
			deleteRequest := map[string]interface{}{
				"AggregatorId":  aggregatorId,
				"ConfigRuleIds": strings.Join(configRuleIds, ","),
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, deleteRequest, &util.RuntimeOptions{})
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
				log.Printf("[ERROR] Failed To Delete Aggregate Config Rules (%s): %v", strings.Join(configRuleIds, ","), err)
				continue
			}
		}
	}
	return nil
}

func TestAccAlicloudConfigAggregateConfigRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_config_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregateConfigRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateConfigRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregateconfigrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregateConfigRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregate_config_rule_name": "${var.name}",
					"aggregator_id":              "${data.alicloud_config_aggregators.default.ids.0}",
					"config_rule_trigger_types":  "ConfigurationItemChangeNotification",
					"source_owner":               "ALIYUN",
					"source_identifier":          "ecs-cpu-min-count-limit",
					"risk_level":                 `1`,
					"resource_types_scope":       []string{"ACS::ECS::Instance"},
					"input_parameters": map[string]string{
						"cpuCount": "4",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregate_config_rule_name": name,
						"risk_level":                 "1",
						"resource_types_scope.#":     "1",
						"config_rule_trigger_types":  "ConfigurationItemChangeNotification",
						"source_owner":               "ALIYUN",
						"source_identifier":          "ecs-cpu-min-count-limit",
						"input_parameters.%":         "1",
						"input_parameters.cpuCount":  "4",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_parameters": map[string]string{
						"cpuCount": "3",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_parameters.%":        "1",
						"input_parameters.cpuCount": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level": `2`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_ids_scope": "${data.alicloud_instances.default.ids.0}",
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
					"tag_key_scope": "tfTest_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_key_scope": "tfTest_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_value_scope": "tfTest 123 update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_value_scope": "tfTest 123 update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"input_parameters": map[string]string{
						"cpuCount": "4",
					},
					"description":                name,
					"risk_level":                 `1`,
					"exclude_resource_ids_scope": "${data.alicloud_instances.default.ids.1}",
					"region_ids_scope":           "cn-shanghai",
					"resource_group_ids_scope":   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tag_key_scope":              "tftest",
					"tag_value_scope":            "tfTest 123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"input_parameters.%":         "1",
						"input_parameters.cpuCount":  "4",
						"description":                name,
						"risk_level":                 "1",
						"exclude_resource_ids_scope": CHECKSET,
						"region_ids_scope":           "cn-shanghai",
						"resource_group_ids_scope":   CHECKSET,
						"tag_key_scope":              "tftest",
						"tag_value_scope":            "tfTest 123",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudConfigAggregateConfigRule_status(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_config_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregateConfigRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateConfigRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregateconfigrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregateConfigRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregate_config_rule_name": "${var.name}",
					"aggregator_id":              "${data.alicloud_config_aggregators.default.ids.0}",
					"config_rule_trigger_types":  "ConfigurationItemChangeNotification",
					"source_owner":               "ALIYUN",
					"source_identifier":          "ecs-cpu-min-count-limit",
					"risk_level":                 `1`,
					"resource_types_scope":       []string{"ACS::ECS::Instance"},
					"input_parameters": map[string]string{
						"cpuCount": "4",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregate_config_rule_name": name,
						"risk_level":                 "1",
						"resource_types_scope.#":     "1",
						"config_rule_trigger_types":  "ConfigurationItemChangeNotification",
						"source_owner":               "ALIYUN",
						"source_identifier":          "ecs-cpu-min-count-limit",
						"input_parameters.%":         "1",
						"input_parameters.cpuCount":  "4",
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
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudConfigAggregateConfigRuleMap0 = map[string]string{
	"aggregate_config_rule_name":  CHECKSET,
	"aggregator_id":               CHECKSET,
	"config_rule_trigger_types":   "ConfigurationItemChangeNotification",
	"description":                 "",
	"input_parameters.%":          "1",
	"input_parameters.cpuCount":   "4",
	"source_owner":                "ALIYUN",
	"source_identifier":           "ecs-cpu-min-count-limit",
	"region_ids_scope":            "",
	"resource_group_ids_scope":    "",
	"risk_level":                  "1",
	"exclude_resource_ids_scope":  "",
	"resource_types_scope.#":      "1",
	"maximum_execution_frequency": "",
	"tag_key_scope":               "",
	"tag_value_scope":             "",
	"status":                      "ACTIVE",
}

func AlicloudConfigAggregateConfigRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_instances" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_config_aggregators" "default" {}

`, name)
}

// Test this case need use a custom `source_identifier`
func SkipTestAccAlicloudConfigAggregateConfigRule_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_config_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregateConfigRuleMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateConfigRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregateconfigrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregateConfigRuleBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregate_config_rule_name": "${var.name}",
					"aggregator_id":              "${data.alicloud_config_aggregators.default.ids.0}",
					"config_rule_trigger_types":  "ConfigurationItemChangeNotification",
					"source_owner":               "CUSTOM_FC",
					"source_identifier":          "*** your_fc_function_arn ***",
					"risk_level":                 `1`,
					"resource_types_scope":       []string{"ACS::ECS::Instance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregate_config_rule_name": name,
						"risk_level":                 "1",
						"resource_types_scope.#":     "1",
						"config_rule_trigger_types":  "ConfigurationItemChangeNotification",
						"source_owner":               "CUSTOM_FC",
						"source_identifier":          "*** your_fc_function_arn ***",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level": `2`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_ids_scope": "${data.alicloud_instances.default.ids.0}",
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
					"tag_key_scope": "tfTest_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_key_scope": "tfTest_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_value_scope": "tfTest 123 update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_value_scope": "tfTest 123 update",
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
					"maximum_execution_frequency": "Six_Hours",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maximum_execution_frequency": "Six_Hours",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_types_scope": []string{"ACS::VPC::VSwitch"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_types_scope.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                 name,
					"risk_level":                  `1`,
					"exclude_resource_ids_scope":  "${data.alicloud_instances.default.ids.1}",
					"region_ids_scope":            "cn-shanghai",
					"resource_group_ids_scope":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tag_key_scope":               "tftest",
					"tag_value_scope":             "tfTest 123",
					"maximum_execution_frequency": "Three_Hours",
					"config_rule_trigger_types":   "ConfigurationItemChangeNotification",
					"resource_types_scope":        []string{"ACS::ECS::Instance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                 name,
						"risk_level":                  "1",
						"exclude_resource_ids_scope":  CHECKSET,
						"region_ids_scope":            "cn-shanghai",
						"resource_group_ids_scope":    CHECKSET,
						"tag_key_scope":               "tftest",
						"tag_value_scope":             "tfTest 123",
						"maximum_execution_frequency": "Three_Hours",
						"config_rule_trigger_types":   "ConfigurationItemChangeNotification",
						"resource_types_scope.#":      "1",
					}),
				),
			},
		},
	})
}

var AlicloudConfigAggregateConfigRuleMap1 = map[string]string{
	"aggregate_config_rule_name":  CHECKSET,
	"aggregator_id":               CHECKSET,
	"config_rule_trigger_types":   "ConfigurationItemChangeNotification",
	"description":                 "",
	"source_owner":                "ALIYUN",
	"source_identifier":           "*** your_fc_function_arn ***",
	"region_ids_scope":            "",
	"resource_group_ids_scope":    "",
	"risk_level":                  "1",
	"exclude_resource_ids_scope":  "",
	"resource_types_scope.#":      "1",
	"maximum_execution_frequency": "",
	"tag_key_scope":               "",
	"tag_value_scope":             "",
	"status":                      "ACTIVE",
}

func AlicloudConfigAggregateConfigRuleBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_instances" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_config_aggregators" "default" {}

`, name)
}

func TestUnitAlicloudConfigAggregateConfigRule(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"aggregate_config_rule_name":  "CreateAggregateConfigRuleValue",
		"aggregator_id":               "CreateAggregateConfigRuleValue",
		"config_rule_trigger_types":   "CreateAggregateConfigRuleValue",
		"description":                 "CreateAggregateConfigRuleValue",
		"exclude_resource_ids_scope":  "CreateAggregateConfigRuleValue",
		"maximum_execution_frequency": "CreateAggregateConfigRuleValue",
		"region_ids_scope":            "CreateAggregateConfigRuleValue",
		"resource_group_ids_scope":    "CreateAggregateConfigRuleValue",
		"source_owner":                "CreateAggregateConfigRuleValue",
		"source_identifier":           "CreateAggregateConfigRuleValue",
		"risk_level":                  1,
		"resource_types_scope":        []string{"CreateAggregateConfigRuleValue"},
		"input_parameters": map[string]string{
			"cpuCount": "4",
		},
		"tag_key_scope":   "CreateAggregateConfigRuleValue",
		"tag_value_scope": "CreateAggregateConfigRuleValue",
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
		// GetAggregateConfigRule
		"ConfigRule": map[string]interface{}{
			"AggregatorId":              "CreateAggregateConfigRuleValue",
			"ConfigRuleId":              "CreateAggregateConfigRuleValue",
			"ConfigRuleName":            "CreateAggregateConfigRuleValue",
			"ConfigRuleTriggerTypes":    "CreateAggregateConfigRuleValue",
			"Description":               "CreateAggregateConfigRuleValue",
			"ExcludeResourceIdsScope":   "CreateAggregateConfigRuleValue",
			"InputParameters":           "CreateAggregateConfigRuleValue",
			"MaximumExecutionFrequency": "CreateAggregateConfigRuleValue",
			"RegionIdsScope":            "CreateAggregateConfigRuleValue",
			"ResourceGroupIdsScope":     "CreateAggregateConfigRuleValue",
			"Scope": map[string]interface{}{
				"ComplianceResourceTypes": "CreateAggregateConfigRuleValue",
			},
			"RiskLevel": 1,
			"Source": map[string]interface{}{
				"Identifier": "CreateAggregateConfigRuleValue",
				"Owner":      "CreateAggregateConfigRuleValue",
			},
			"ConfigRuleState": "INACTIVE",
			"TagKeyScope":     "CreateAggregateConfigRuleValue",
			"TagValueScope":   "CreateAggregateConfigRuleValue",
		},
		"ConfigRuleId": "CreateAggregateConfigRuleValue",
	}
	CreateMockResponse := map[string]interface{}{
		//CreateAggregateConfigRule
		"ConfigRuleId": "CreateAggregateConfigRuleValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_config_aggregate_config_rule", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigAggregateConfigRuleCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetConfigRule Response
		"ConfigRuleId": "CreateAggregateConfigRuleValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateAggregateConfigRule" {
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
		err := resourceAlicloudConfigAggregateConfigRuleCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dInit.State(), nil)
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
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigAggregateConfigRuleUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	//UpdateAggregateConfigRule
	attributesDiff := map[string]interface{}{
		"config_rule_trigger_types":  "UpdateAggregateConfigRuleValue",
		"resource_types_scope":       []string{"UpdateAggregateConfigRuleValue"},
		"risk_level":                 2,
		"description":                "UpdateAggregateConfigRuleValue",
		"exclude_resource_ids_scope": "UpdateAggregateConfigRuleValue",
		"input_parameters": map[string]string{
			"cpuCount": "8",
		},
		"maximum_execution_frequency": "UpdateAggregateConfigRuleValue",
		"region_ids_scope":            "UpdateAggregateConfigRuleValue",
		"resource_group_ids_scope":    "UpdateAggregateConfigRuleValue",
		"tag_key_scope":               "UpdateAggregateConfigRuleValue",
		"tag_value_scope":             "UpdateAggregateConfigRuleValue",
	}
	diff, err := newInstanceDiff("alicloud_config_aggregate_config_rule", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAggregateConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleTriggerTypes": "UpdateAggregateConfigRuleValue",
			"Scope": map[string]interface{}{
				"ComplianceResourceTypes": "UpdateAggregateConfigRuleValue",
			},
			"RiskLevel":                 2,
			"Description":               "UpdateAggregateConfigRuleValue",
			"ExcludeResourceIdsScope":   "UpdateAggregateConfigRuleValue",
			"InputParameters":           "UpdateAggregateConfigRuleValue",
			"MaximumExecutionFrequency": "UpdateAggregateConfigRuleValue",
			"RegionIdsScope":            "UpdateAggregateConfigRuleValue",
			"ResourceGroupIdsScope":     "UpdateAggregateConfigRuleValue",
			"TagKeyScope":               "UpdateAggregateConfigRuleValue",
			"TagValueScope":             "UpdateAggregateConfigRuleValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateAggregateConfigRule" {
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
		err := resourceAlicloudConfigAggregateConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	//ActiveAggregateConfigRules
	attributesDiff = map[string]interface{}{
		"status": "ACTIVE",
	}
	diff, err = newInstanceDiff("alicloud_config_aggregate_config_rule", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAggregateConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleState": "ACTIVE",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "ActiveAggregateConfigRules" {
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
		err := resourceAlicloudConfigAggregateConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	//DeactiveAggregateConfigRules
	attributesDiff = map[string]interface{}{
		"status": "INACTIVE",
	}
	diff, err = newInstanceDiff("alicloud_config_aggregate_config_rule", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAggregateConfigRule Response
		"ConfigRule": map[string]interface{}{
			"ConfigRuleState": "INACTIVE",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeactiveAggregateConfigRules" {
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
		err := resourceAlicloudConfigAggregateConfigRuleUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	//Read
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_config_aggregate_config_rule", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetAggregateConfigRule" {
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
		err := resourceAlicloudConfigAggregateConfigRuleRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewConfigClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:       String("loadEndpoint error"),
			Data:       String("loadEndpoint error"),
			Message:    String("loadEndpoint error"),
			StatusCode: tea.Int(400),
		}
	})
	err = resourceAlicloudConfigAggregateConfigRuleDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_config_aggregate_config_rule", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_config_rule"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteAggregateConfigRules" {
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
		err := resourceAlicloudConfigAggregateConfigRuleDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
