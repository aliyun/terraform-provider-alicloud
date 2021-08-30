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
	resource.AddTestSweepers("alicloud_config_compliance_pack", &resource.Sweeper{
		Name: "alicloud_config_compliance_pack",
		F:    testSweepConfigCompliancePack,
	})
}

func testSweepConfigCompliancePack(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
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
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	action := "ListCompliancePacks"
	var response map[string]interface{}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
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
			log.Printf("[ERROR] Failed To Delete Compliance Packs (%s): %v", strings.Join(compliancePackIds, ","), err)
		}
	}
	return nil
}

func TestAccAlicloudConfigCompliancePack_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigcompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigCompliancePackBasicDependence0)
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
					"compliance_pack_template_id": "ct-3d20ff4e06a30027f76e",
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
						"compliance_pack_name":        name,
						"compliance_pack_template_id": "ct-3d20ff4e06a30027f76e",
						"config_rules.#":              "1",
						"description":                 name,
						"risk_level":                  "1",
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

func TestAccAlicloudConfigCompliancePack_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_config_compliance_pack.default"
	ra := resourceAttrInit(resourceId, AlicloudConfigCompliancePackMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigCompliancePack")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sconfigcompliancepack%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudConfigCompliancePackBasicDependence0)
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
					"compliance_pack_template_id": "ct-3d20ff4e06a30027f76e",
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
						"compliance_pack_name":        name,
						"compliance_pack_template_id": "ct-3d20ff4e06a30027f76e",
						"config_rules.#":              "1",
						"description":                 name,
						"risk_level":                  "1",
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

var AlicloudConfigCompliancePackMap0 = map[string]string{
	"compliance_pack_name":        CHECKSET,
	"compliance_pack_template_id": CHECKSET,
	"config_rules.#":              "1",
	"description":                 CHECKSET,
	"risk_level":                  "1",
	"status":                      CHECKSET,
}

func AlicloudConfigCompliancePackBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

`, name)
}
