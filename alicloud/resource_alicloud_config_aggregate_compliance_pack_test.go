package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_config_aggregate_compliance_pack", &resource.Sweeper{
		Name: "alicloud_config_aggregate_compliance_pack",
		F:    testSweepConfigAggregateCompliancePack,
	})
}

func testSweepConfigAggregateCompliancePack(region string) error {
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
				log.Printf("[INFO] Skipping Aggregate: %v (%v)", item["AggregatorName"], item["AggregatorId"])
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

	// Delete Aggregate Compliance Packs
	for _, aggregatorId := range aggregatorIds {
		action := "ListAggregateCompliancePacks"
		request := map[string]interface{}{
			"AggregatorId": aggregatorId,
			"PageSize":     PageSizeLarge,
			"PageNumber":   1,
		}
		compliancePackIds := make([]string, 0)
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
				log.Printf("[ERROR] Failed To List Aggregate Compliance Packs : %s", err)
			}
			resp, err := jsonpath.Get("$.CompliancePacksResult.CompliancePacks", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CompliancePacksResult.CompliancePacks", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				skip := true
				item := v.(map[string]interface{})
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["CompliancePackName"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Aggregate Compliance Pack: %v (%v)", item["CompliancePackName"], item["CompliancePackId"])
					continue
				}
				compliancePackIds = append(compliancePackIds, fmt.Sprint(item["CompliancePackId"]))
			}
			if len(result) < PageSizeLarge {
				break
			}
			request["PageNumber"] = request["PageNumber"].(int) + 1
		}

		if len(compliancePackIds) > 0 {
			log.Printf("[INFO] Deleting Aggregate Compliance Packs: (%s)", strings.Join(compliancePackIds, ","))
			action = "DeleteAggregateCompliancePacks"
			deleteRequest := map[string]interface{}{
				"CompliancePackIds": strings.Join(compliancePackIds, ","),
				"AggregatorId":      aggregatorId,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(time.Minute*10, func() *resource.RetryError {
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
				log.Printf("[ERROR] Failed To Delete Aggregate Compliance Packs (%s): %v", strings.Join(compliancePackIds, ","), err)
				continue
			}
		}
	}
	return nil
}

func TestAccAlicloudConfigAggregateCompliancePack_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregateCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatecompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregateCompliancePackBasicDependence0)
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
					"aggregator_id":                  "${data.alicloud_config_aggregators.default.ids.0}",
					"aggregate_compliance_pack_name": name,
					"compliance_pack_template_id":    "ct-3d20ff4e06a30027f76e",
					"config_rules": []map[string]interface{}{
						{
							"managed_rule_identifier": "ecs-snapshot-retention-days",
							"config_rule_parameters": []map[string]interface{}{
								{
									"parameter_name":  "days",
									"parameter_value": "7",
								},
							},
						},
					},
					"description": name,
					"risk_level":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_id":                  CHECKSET,
						"aggregate_compliance_pack_name": name,
						"compliance_pack_template_id":    "ct-3d20ff4e06a30027f76e",
						"config_rules.#":                 "1",
						"description":                    name,
						"risk_level":                     "1",
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
					"config_rules": []map[string]interface{}{
						{
							"managed_rule_identifier": "ecs-snapshot-retention-days",
							"config_rule_parameters": []map[string]interface{}{
								{
									"parameter_name":  "days",
									"parameter_value": "7",
								},
							},
						},
						{
							"managed_rule_identifier": "ecs-instance-expired-check",
							"config_rule_parameters": []map[string]interface{}{
								{
									"parameter_name":  "days",
									"parameter_value": "60",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rules.#": "2",
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
					"config_rules": []map[string]interface{}{
						{
							"managed_rule_identifier": "ecs-snapshot-retention-days",
							"config_rule_parameters": []map[string]interface{}{
								{
									"parameter_name":  "days",
									"parameter_value": "7",
								},
							},
						},
					},
					"description": name,
					"risk_level":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rules.#": "1",
						"description":    name,
						"risk_level":     "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudConfigAggregateCompliancePack_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregateCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatecompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregateCompliancePackBasicDependence0)
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
					"aggregator_id":                  "${data.alicloud_config_aggregators.default.ids.0}",
					"aggregate_compliance_pack_name": name,
					"compliance_pack_template_id":    "ct-3d20ff4e06a30027f76e",
					"config_rules": []map[string]interface{}{
						{
							"managed_rule_identifier": "oss-bucket-public-read-prohibited",
						},
					},
					"description": name,
					"risk_level":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_id":                  CHECKSET,
						"aggregate_compliance_pack_name": name,
						"compliance_pack_template_id":    "ct-3d20ff4e06a30027f76e",
						"config_rules.#":                 "1",
						"description":                    name,
						"risk_level":                     "1",
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

func TestAccAlicloudConfigAggregateCompliancePack_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigAggregateCompliancePackMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatecompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigAggregateCompliancePackBasicDependence1)
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
					"aggregator_id":                  "${data.alicloud_config_aggregators.default.ids.0}",
					"aggregate_compliance_pack_name": name,
					"config_rules": []map[string]interface{}{
						{
							"config_rule_id": "${alicloud_config_aggregate_config_rule.default.config_rule_id}",
							"managed_rule_identifier": "${alicloud_config_aggregate_config_rule.default.source_identifier}",
							"config_rule_parameters": []map[string]interface{}{
								{
									"parameter_name":  "cpuCount",
									"parameter_value": "4",
								},
							},
						},
					},
					"description": name,
					"risk_level":  "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_id":                  CHECKSET,
						"aggregate_compliance_pack_name": name,
						"config_rules.#":                 "1",
						"description":                    name,
						"risk_level":                     "1",
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

var AlicloudConfigAggregateCompliancePackMap0 = map[string]string{
	"aggregator_id":                  CHECKSET,
	"aggregate_compliance_pack_name": CHECKSET,
	"compliance_pack_template_id":    CHECKSET,
	"config_rules.#":                 "1",
	"description":                    CHECKSET,
	"risk_level":                     "1",
	"status":                         CHECKSET,
}

var AlicloudConfigAggregateCompliancePackMap1 = map[string]string{}

func AlicloudConfigAggregateCompliancePackBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_config_aggregators" "default" {}

`, name)
}

func AlicloudConfigAggregateCompliancePackBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_instances" "default" {}

data "alicloud_config_aggregators" "default" {}

resource "alicloud_config_aggregate_config_rule" "default" {
  aggregator_id              = data.alicloud_config_aggregators.default.ids.0
  aggregate_config_rule_name = var.name
  source_owner               = "ALIYUN"
  source_identifier    		= "ecs-cpu-min-count-limit"
  config_rule_trigger_types = "ConfigurationItemChangeNotification"
  resource_types_scope      = ["ACS::ECS::Instance"]
  risk_level                = 1
  description                = var.name
  exclude_resource_ids_scope = data.alicloud_instances.default.ids.0
  input_parameters = {
    cpuCount = "4",
  }
  region_ids_scope         = "cn-hangzhou"
  resource_group_ids_scope = data.alicloud_resource_manager_resource_groups.default.ids.0
  tag_key_scope            = "tFTest"
  tag_value_scope          = "forTF 123"
}

`, name)
}