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
	resource.AddTestSweepers("alicloud_config_compliance_pack", &resource.Sweeper{
		Name: "alicloud_config_compliance_pack",
		F:    testSweepConfigCompliancePack,
	})
}

func testSweepConfigCompliancePack(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	compliancePackIds := make([]string, 0)

	action := "ListCompliancePacks"
	var response map[string]interface{}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
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
			log.Printf("[ERROR] Failed To List Compliance Packs: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.CompliancePacksResult.CompliancePacks", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CompliancePacksResult.CompliancePacks", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			skip := true
			item := v.(map[string]interface{})
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["CompliancePackName"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Compliance Pack: %v (%v)", item["CompliancePackName"], item["CompliancePackId"])
					continue
				}
			}
			compliancePackIds = append(compliancePackIds, fmt.Sprint(item["CompliancePackId"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if len(compliancePackIds) > 0 {
		log.Printf("[INFO] Deleting Compliance Packs: (%s)", strings.Join(compliancePackIds, ","))
		action = "DeleteCompliancePacks"
		deleteRequest := map[string]interface{}{
			"CompliancePackIds": strings.Join(compliancePackIds, ","),
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
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
			log.Printf("[ERROR] Failed To Delete Compliance Packs (%s): %v", strings.Join(compliancePackIds, ","), err)
		}
	}
	return nil
}

func TestAccAliCloudConfigCompliancePack_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigcompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigCompliancePackBasicDependence0)
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
					"compliance_pack_name": name,
					"description":          name,
					"risk_level":           "1",
					"config_rules": []map[string]interface{}{
						{
							"managed_rule_identifier": "oss-bucket-public-read-prohibited",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compliance_pack_name": name,
						"description":          name,
						"risk_level":           "1",
						"config_rules.#":       "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"compliance_pack_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compliance_pack_name": name + "_update",
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

	prefixes := []string{
		"oss-",
		"ecs-",
		"governance-",
	}
	testSweepConfigRuleByPrefixes(defaultRegionToTest, prefixes)
}

func TestAccAliCloudConfigCompliancePack_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigcompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigCompliancePackBasicDependence0)
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
					"compliance_pack_name":        name,
					"description":                 name,
					"risk_level":                  "1",
					"compliance_pack_template_id": "${var.compliance_pack_template_id}",
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
						"compliance_pack_name":        name,
						"description":                 name,
						"risk_level":                  "1",
						"compliance_pack_template_id": CHECKSET,
						"config_rules.#":              "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudConfigCompliancePack_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AliCloudConfigCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigcompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudConfigCompliancePackBasicDependence1)
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
					"compliance_pack_name": name,
					"description":          name,
					"risk_level":           "1",
					"config_rule_ids": []map[string]interface{}{
						{
							"config_rule_id": "${alicloud_config_rule.default.0.id}",
						},
						{
							"config_rule_id": "${alicloud_config_rule.default.1.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compliance_pack_name": name,
						"description":          name,
						"risk_level":           "1",
						"config_rule_ids.#":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_rule_ids": []map[string]interface{}{
						{
							"config_rule_id": "${alicloud_config_rule.default.0.id}",
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
							"config_rule_id": "${alicloud_config_rule.default.0.id}",
						},
						{
							"config_rule_id": "${alicloud_config_rule.default.1.id}",
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

var AliCloudConfigCompliancePackMap0 = map[string]string{
	"status": CHECKSET,
}

func AliCloudConfigCompliancePackBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	variable "compliance_pack_template_id" {
		default = "ct-3d20ff4e06a30027f76e"
	}
`, name)
}

func AliCloudConfigCompliancePackBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
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

	resource "alicloud_config_rule" "default" {
  		count                      = 2
  		rule_name                  = var.name
  		description                = var.name
  		source_identifier          = "ecs-instances-in-vpc"
  		source_owner               = "ALIYUN"
  		resource_types_scope       = ["ACS::ECS::Instance"]
  		risk_level                 = 1
  		config_rule_trigger_types  = "ConfigurationItemChangeNotification"
		tag_key_scope              = "tfTest"
		tag_value_scope            = "tfTest 123"
		resource_group_ids_scope   = data.alicloud_resource_manager_resource_groups.default.ids.0
		exclude_resource_ids_scope = alicloud_instance.default.id
		region_ids_scope           = "%s"
		input_parameters = {
			vpcIds = data.alicloud_vpcs.default.ids.0
		}
	}
`, name, defaultRegionToTest)
}

func TestUnitAliCloudConfigCompliancePack(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_config_compliance_pack"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_config_compliance_pack"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"compliance_pack_name":        "CreateCompliancePackValue",
		"compliance_pack_template_id": "CreateCompliancePackValue",
		"config_rules": []interface{}{
			map[string]interface{}{
				"managed_rule_identifier": "CreateCompliancePackValue",
				"config_rule_parameters": []interface{}{
					map[string]interface{}{
						"parameter_name":  "CreateCompliancePackValue",
						"parameter_value": "CreateCompliancePackValue",
					},
				},
			},
		},
		"description": "CreateCompliancePackValue",
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
		// GetCompliancePack
		"CompliancePack": map[string]interface{}{
			"ConfigRules": []interface{}{
				map[string]interface{}{
					"ManagedRuleIdentifier": "CreateCompliancePackValue",
					"ConfigRuleParameters": []interface{}{
						map[string]interface{}{
							"ParameterName":  "CreateCompliancePackValue",
							"ParameterValue": "CreateCompliancePackValue",
						},
					},
				},
			},
			"CompliancePackName":       "CreateCompliancePackValue",
			"CompliancePackTemplateId": "CreateCompliancePackValue",
			"Description":              "CreateCompliancePackValue",
			"RiskLevel":                1,
			"Status":                   "ACTIVE",
		},
	}
	CreateMockResponse := map[string]interface{}{
		// CreateCompliancePack
		"CompliancePackId": "CreateCompliancePackValue",
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
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_config_compliance_pack", errorCode))
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
	err = resourceAliCloudConfigCompliancePackCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff := map[string]interface{}{
		// GetCompliancePack Response
		"CompliancePack": map[string]interface{}{
			"CompliancePackId": "CreateCompliancePackValue",
		},
	}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateCompliancePack" {
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
		err := resourceAliCloudConfigCompliancePackCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_compliance_pack"].Schema).Data(dInit.State(), nil)
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
	err = resourceAliCloudConfigCompliancePackUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateCompliancePack
	attributesDiff := map[string]interface{}{
		"compliance_pack_name": "UpdateCompliancePackValue",
		"description":          "UpdateCompliancePackValue",
		"risk_level":           2,
		"config_rules": []interface{}{
			map[string]interface{}{
				"managed_rule_identifier": "UpdateCompliancePackValue",
				"config_rule_parameters": []interface{}{
					map[string]interface{}{
						"parameter_name":  "UpdateCompliancePackValue",
						"parameter_value": "UpdateCompliancePackValue",
					},
				},
			},
		},
	}
	diff, err := newInstanceDiff("alicloud_config_compliance_pack", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_config_compliance_pack"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetCompliancePack Response
		"CompliancePack": map[string]interface{}{
			"CompliancePackName": "UpdateCompliancePackValue",
			"Description":        "UpdateCompliancePackValue",
			"RiskLevel":          2,
			"ConfigRules": []interface{}{
				map[string]interface{}{
					"ManagedRuleIdentifier": "UpdateCompliancePackValue",
					"ConfigRuleParameters": []interface{}{
						map[string]interface{}{
							"ParameterName":  "UpdateCompliancePackValue",
							"ParameterValue": "UpdateCompliancePackValue",
						},
					},
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "CompliancePackAlreadyPending", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateCompliancePack" {
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
		err := resourceAliCloudConfigCompliancePackUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_config_compliance_pack"].Schema).Data(dExisted.State(), nil)
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

	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "Invalid.CompliancePackId.Value", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetCompliancePack" {
				switch errorCode {
				case "{}", "Invalid.CompliancePackId.Value":
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
		err := resourceAliCloudConfigCompliancePackRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}", "Invalid.CompliancePackId.Value":
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
	err = resourceAliCloudConfigCompliancePackDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteCompliancePacks" {
				switch errorCode {
				case "NonRetryableError":
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
		err := resourceAliCloudConfigCompliancePackDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		}
	}

}
