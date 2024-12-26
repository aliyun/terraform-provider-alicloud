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
	resource.AddTestSweepers("alicloud_config_aggregate_compliance_pack", &resource.Sweeper{
		Name: "alicloud_config_aggregate_compliance_pack",
		F:    testSweepConfigAggregateCompliancePack,
	})
}

func testSweepConfigAggregateCompliancePack(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	// Get all AggregatorId
	aggregatorIds := make([]string, 0)
	action := "ListAggregators"
	request := map[string]interface{}{
		"MaxResults": PageSizeLarge,
	}
	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
			if !sweepAll() {
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
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				response, err = client.RpcGet("Config", "2020-09-07", action, request, nil)
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
				response, err = client.RpcPost("Config", "2020-09-07", action, nil, deleteRequest, false)
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

func TestAccAliCloudConfigAggregateCompliancePack_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregateCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatecompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregateCompliancePackBasicDependence0)
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
					"aggregator_id":                  "${alicloud_config_aggregator.default.id}",
					"aggregate_compliance_pack_name": name,
					"description":                    name,
					"risk_level":                     "1",
					"config_rules": []map[string]interface{}{
						{
							"managed_rule_identifier": "oss-bucket-public-read-prohibited",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_id":                  CHECKSET,
						"aggregate_compliance_pack_name": name,
						"description":                    name,
						"risk_level":                     "1",
						"config_rules.#":                 "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"aggregate_compliance_pack_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregate_compliance_pack_name": name + "_update",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rules.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudConfigAggregateCompliancePack_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregateCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatecompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregateCompliancePackBasicDependence0)
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
					"aggregator_id":                  "${alicloud_config_aggregator.default.id}",
					"aggregate_compliance_pack_name": name,
					"description":                    name,
					"risk_level":                     "1",
					"compliance_pack_template_id":    "${var.compliance_pack_template_id}",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_id":                  CHECKSET,
						"aggregate_compliance_pack_name": name,
						"description":                    name,
						"risk_level":                     "1",
						"compliance_pack_template_id":    CHECKSET,
						"config_rules.#":                 "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudConfigAggregateCompliancePack_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_aggregate_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigAggregateCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigAggregateCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigaggregatecompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigAggregateCompliancePackBasicDependence1)
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
					"aggregator_id":                  "${alicloud_config_aggregator.default.id}",
					"aggregate_compliance_pack_name": name,
					"description":                    name,
					"risk_level":                     "1",
					"config_rule_ids": []map[string]interface{}{
						{
							"config_rule_id": "${alicloud_config_aggregate_config_rule.default.0.config_rule_id}",
						},
						{
							"config_rule_id": "${alicloud_config_aggregate_config_rule.default.1.config_rule_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"aggregator_id":                  CHECKSET,
						"aggregate_compliance_pack_name": name,
						"description":                    name,
						"risk_level":                     "1",
						"config_rule_ids.#":              "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_rule_ids": []map[string]interface{}{
						{
							"config_rule_id": "${alicloud_config_aggregate_config_rule.default.0.config_rule_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_rule_ids": []map[string]interface{}{
						{
							"config_rule_id": "${alicloud_config_aggregate_config_rule.default.0.config_rule_id}",
						},
						{
							"config_rule_id": "${alicloud_config_aggregate_config_rule.default.1.config_rule_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_rule_ids": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_rule_ids.#": "0",
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

var AliCloudConfigAggregateCompliancePackMap0 = map[string]string{
	"aggregator_compliance_pack_id": CHECKSET,
	"status":                        CHECKSET,
}

func AliCloudConfigAggregateCompliancePackBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}
	
	variable "compliance_pack_template_id" {
		default = "ct-3d20ff4e06a30027f76e"
	}

	data "alicloud_resource_manager_accounts" "default" {
	  status = "CreateSuccess"
	}

	resource "alicloud_config_aggregator" "default" {
	  aggregator_accounts {
		account_id   = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
		account_name = data.alicloud_resource_manager_accounts.default.accounts.0.display_name
		account_type = "ResourceDirectory"
	  }
	  aggregator_name = var.name
	  description     = var.name
	  aggregator_type = "CUSTOM"
	}
	
	resource "alicloud_config_aggregate_config_rule" "default" {
	  aggregate_config_rule_name = "contains-tag"
	  aggregator_id              = alicloud_config_aggregator.default.id
	  config_rule_trigger_types  = "ConfigurationItemChangeNotification"
	  source_owner               = "ALIYUN"
	  source_identifier          = "contains-tag"
	  description                = var.name
	  risk_level                 = 1
	  resource_types_scope       = ["ACS::ECS::Instance"]
	  input_parameters = {
		key   = "example"
		value = "example"
	  }
	}
`, name)
}

func AliCloudConfigAggregateCompliancePackBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_resource_manager_accounts" "default" {
	  status = "CreateSuccess"
	}

	resource "alicloud_config_aggregator" "default" {
	  aggregator_accounts {
		account_id   = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
		account_name = data.alicloud_resource_manager_accounts.default.accounts.0.display_name
		account_type = "ResourceDirectory"
	  }
	  aggregator_name = var.name
	  description     = var.name
	  aggregator_type = "CUSTOM"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		instance_name              = var.name
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_zones.default.zones.0.id
  		instance_charge_type       = "PostPaid"
  		password                   = "YourPassword12345!"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = data.alicloud_vswitches.default.ids.0
	}

	resource "alicloud_config_aggregate_config_rule" "default" {
  		count                      = 2
  		aggregator_id              = alicloud_config_aggregator.default.id
  		aggregate_config_rule_name = var.name
  		source_owner               = "ALIYUN"
  		source_identifier          = "ecs-cpu-min-count-limit"
  		config_rule_trigger_types  = "ConfigurationItemChangeNotification"
  		resource_types_scope       = ["ACS::ECS::Instance"]
  		risk_level                 = 1
  		description                = var.name
  		exclude_resource_ids_scope = alicloud_instance.default.id
  		input_parameters = {
    		cpuCount = "4",
  		}
  		region_ids_scope         = "%s"
  		resource_group_ids_scope = data.alicloud_resource_manager_resource_groups.default.ids.0
  		tag_key_scope            = "tFTest"
  		tag_value_scope          = "forTF 123"
	}
`, name, defaultRegionToTest)
}

func TestUnitAliCloudConfigAggregateCompliancePack(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"aggregate_compliance_pack_name": "CreateAggregateCompliancePackValue",
		"aggregator_id":                  "CreateAggregateCompliancePackValue",
		"compliance_pack_template_id":    "CreateAggregateCompliancePackValue",
		"config_rules": []map[string]interface{}{
			{
				"managed_rule_identifier": "CreateAggregateCompliancePackValue",
				"config_rule_parameters": []map[string]interface{}{
					{
						"parameter_name":  "CreateAggregateCompliancePackValue",
						"parameter_value": "CreateAggregateCompliancePackValue",
					},
				},
			},
		},
		"description": "CreateAggregateCompliancePackValue",
		"risk_level":  1,
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
		// GetAggregateCompliancePack
		"CompliancePack": map[string]interface{}{
			"AggregatorId":             "CreateAggregateCompliancePackValue",
			"CompliancePackName":       "CreateAggregateCompliancePackValue",
			"CompliancePackTemplateId": "CreateAggregateCompliancePackValue",
			"ConfigRules": []interface{}{
				map[string]interface{}{
					"ConfigRuleParameters": []interface{}{
						map[string]interface{}{
							"ParameterName":  "CreateAggregateCompliancePackValue",
							"ParameterValue": "CreateAggregateCompliancePackValue",
						},
					},
					"ManagedRuleIdentifier": "CreateAggregateCompliancePackValue",
				},
			},
			"Description": "CreateAggregateCompliancePackValue",
			"RiskLevel":   1,
			"Status":      "ACTIVE",
		},
		"CompliancePackId": "CreateAggregateCompliancePackValue",
	}
	CreateMockResponse := map[string]interface{}{
		// CreateAggregateCompliancePack
		"CompliancePackId": "CreateAggregateCompliancePackValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_config_aggregate_compliance_pack", errorCode))
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
	err = resourceAliCloudConfigAggregateCompliancePackCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetAggregateCompliancePack Response
		"CompliancePackId": "CreateAggregateCompliancePackValue",
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateAggregateCompliancePack" {
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
		err := resourceAliCloudConfigAggregateCompliancePackCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(dInit.State(), nil)
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
	err = resourceAliCloudConfigAggregateCompliancePackUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateAggregateCompliancePack
	attributesDiff := map[string]interface{}{
		"config_rules": []map[string]interface{}{
			{
				"managed_rule_identifier": "UpdateAggregateCompliancePackValue",
				"config_rule_parameters": []map[string]interface{}{
					{
						"parameter_name":  "UpdateAggregateCompliancePackValue",
						"parameter_value": "UpdateAggregateCompliancePackValue",
					},
				},
			},
		},
		"description":                    "UpdateAggregateCompliancePackValue",
		"risk_level":                     2,
		"aggregate_compliance_pack_name": "UpdateAggregateCompliancePackValue",
	}
	diff, err := newInstanceDiff("alicloud_config_aggregate_compliance_pack", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAggregateCompliancePack Response
		"CompliancePack": map[string]interface{}{
			"ConfigRules": []interface{}{
				map[string]interface{}{
					"ConfigRuleParameters": []interface{}{
						map[string]interface{}{
							"ParameterName":  "UpdateAggregateCompliancePackValue",
							"ParameterValue": "UpdateAggregateCompliancePackValue",
						},
					},
					"ManagedRuleIdentifier": "UpdateAggregateCompliancePackValue",
				},
			},

			"Description":        "UpdateAggregateCompliancePackValue",
			"RiskLevel":          2,
			"CompliancePackName": "UpdateAggregateCompliancePackValue",
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateAggregateCompliancePack" {
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
		err := resourceAliCloudConfigAggregateCompliancePackUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				_ = dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// AttachAggregateConfigRuleToCompliancePack
	attributesDiff = map[string]interface{}{
		"config_rule_ids": []map[string]interface{}{
			{
				"config_rule_id": "AttachAggregateConfigRuleToCompliancePackValue",
			},
		},
	}
	diff, err = newInstanceDiff("alicloud_config_aggregate_compliance_pack", attributes, attributesDiff, dExisted.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(dExisted.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetAggregateCompliancePack Response
		"CompliancePack": map[string]interface{}{
			"ConfigRules": []interface{}{
				map[string]interface{}{
					"ConfigRuleId": "AttachAggregateConfigRuleToCompliancePackValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "AttachAggregateConfigRuleToCompliancePack" {
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
		err := resourceAliCloudConfigAggregateCompliancePackUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(dExisted.State(), nil)
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
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_config_aggregate_compliance_pack", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetAggregateCompliancePack" {
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
		err := resourceAliCloudConfigAggregateCompliancePackRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
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
	err = resourceAliCloudConfigAggregateCompliancePackDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	attributesDiff = map[string]interface{}{}
	diff, err = newInstanceDiff("alicloud_config_aggregate_compliance_pack", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_aggregate_compliance_pack"].Schema).Data(dInit.State(), diff)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "Invalid.AggregatorId.Value", "Invalid.CompliancePackId.Value"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteAggregateCompliancePacks" {
				switch errorCode {
				case "NonRetryableError", "Invalid.AggregatorId.Value", "Invalid.CompliancePackId.Value":
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
		err := resourceAliCloudConfigAggregateCompliancePackDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "Invalid.AggregatorId.Value", "Invalid.CompliancePackId.Value":
			assert.Nil(t, err)
		}
	}
}
