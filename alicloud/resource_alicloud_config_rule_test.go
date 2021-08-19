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
