package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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
