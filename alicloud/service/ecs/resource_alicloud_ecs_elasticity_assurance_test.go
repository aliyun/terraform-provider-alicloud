package ecs_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"

	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// Case 1
func TestAccAliCloudEcsElasticityAssurance_basic1716(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsElasticityAssuranceMap1716)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount":                     "1",
					"description":                         "before",
					"zone_ids":                            []string{"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"},
					"private_pool_options_name":           name,
					"period":                              "1",
					"private_pool_options_match_criteria": "Open",
					"instance_type":                       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"period_unit":                         "Month",
					"auto_renew":                          "true",
					"auto_renew_period":                   "2",
					"auto_renew_period_unit":              "Year",
					"assurance_times":                     "Unlimited",
					"start_time":                          time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z"),
					"resource_group_id":                   "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount":                     "1",
						"description":                         "before",
						"zone_ids.#":                          "1",
						"private_pool_options_name":           name,
						"period":                              "1",
						"private_pool_options_match_criteria": "Open",
						"instance_type.#":                     "1",
						"period_unit":                         "Month",
						"auto_renew":                          "true",
						"auto_renew_period":                   "2",
						"auto_renew_period_unit":              "Year",
						"assurance_times":                     "Unlimited",
						"start_time":                          CHECKSET,
						"resource_group_id":                   CHECKSET,
						"tags.%":                              "2",
						"tags.Created":                        "TF",
						"tags.For":                            "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

func TestAccAliCloudEcsElasticityAssurance_basic1717(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsElasticityAssuranceMap1716)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_ids":          []string{"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"},
					"instance_type":     []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"instance_amount":   "1",
					"auto_renew":        "true",
					"auto_renew_period": "1",
					"period_unit":       "Month",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_ids.#":        "1",
						"instance_type.#":   "1",
						"instance_amount":   "1",
						"auto_renew":        "true",
						"auto_renew_period": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_pool_options_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_pool_options_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "after_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "after_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":             "true",
					"auto_renew_period":      "2",
					"auto_renew_period_unit": "Year",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":             "true",
						"auto_renew_period":      "2",
						"auto_renew_period_unit": "Year",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

var AliCloudEcsElasticityAssuranceMap1716 = map[string]string{}

func AliCloudEcsElasticityAssuranceBasicDependence1716(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}

data "alicloud_instance_types" "default" {
	instance_type_family = "ecs.c6"
}
`, name)
}

// Case 3: time-division elasticity assurance (PeriodUnit=Day + RecurrenceRules)
func TestAccAliCloudEcsElasticityAssurance_timeDivision(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsElasticityAssuranceMap1716)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEcsElasticityAssurance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount":                     "1",
					"zone_ids":                            []string{"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"},
					"instance_type":                       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"private_pool_options_name":           name,
					"private_pool_options_match_criteria": "Open",
					"period":                              "7",
					"period_unit":                         "Day",
					"assurance_times":                     "Unlimited",
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "8",
							"end_hour":         "12",
						},
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "13",
							"end_hour":         "17",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_amount":                     "1",
						"zone_ids.#":                          "1",
						"instance_type.#":                     "1",
						"private_pool_options_name":           name,
						"private_pool_options_match_criteria": "Open",
						"period":                              "7",
						"period_unit":                         "Day",
						"assurance_times":                     "Unlimited",
						"recurrence_rules.#":                  "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "recurrence_rules.*", map[string]string{
						"recurrence_type":  "Daily",
						"recurrence_value": "1",
						"start_hour":       "8",
						"end_hour":         "12",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "recurrence_rules.*", map[string]string{
						"recurrence_type":  "Daily",
						"recurrence_value": "1",
						"start_hour":       "13",
						"end_hour":         "17",
					}),
				),
			},
			{
				// The rules are stored as a set, so declaring the same rules in a different order must
				// produce an empty plan instead of proposing to replace the assurance.
				Config: testAccConfig(map[string]interface{}{
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "13",
							"end_hour":         "17",
						},
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "8",
							"end_hour":         "12",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Changing the content of a rule must still propose a replacement, because
				// recurrence_rules is ForceNew. PlanOnly steps skip the pre-apply hook, so the
				// action is asserted on the non-refresh plan instead.
				Config: testAccConfig(map[string]interface{}{
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "13",
							"end_hour":         "17",
						},
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "8",
							"end_hour":         "18",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionReplace),
					},
				},
			},
			{
				// The reordered rules converge after apply without any residual diff.
				Config: testAccConfig(map[string]interface{}{
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "13",
							"end_hour":         "17",
						},
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "8",
							"end_hour":         "12",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recurrence_rules.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
		},
	})
}

// Case 4: invalid recurrence_rules must be rejected at plan time without crashing the provider
func TestAccAliCloudEcsElasticityAssurance_timeDivisionInvalidRules(t *testing.T) {
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEcsElasticityAssurance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount":                     "1",
					"zone_ids":                            []string{"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"},
					"instance_type":                       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"private_pool_options_name":           name,
					"private_pool_options_match_criteria": "Open",
					"period":                              "1",
					"period_unit":                         "Day",
					"assurance_times":                     "Unlimited",
					"recurrence_rules": []map[string]interface{}{
						{
							"start_hour": "8",
						},
					},
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s).*recurrence_rules\.0: (start_hour|end_hour|recurrence_type|recurrence_value) must be specified.*`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "8",
							"end_hour":         "10",
						},
					},
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s).*at least 4 hours.*`),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period_unit": "Month",
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "8",
							"end_hour":         "12",
						},
					},
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s).*recurrence_rules can only be specified when period_unit is Day.*`),
			},
			{
				// recurrence_value is validated against the documented format of each recurrence_type
				Config: testAccConfig(map[string]interface{}{
					"period_unit": "Day",
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "32",
							"start_hour":       "8",
							"end_hour":         "12",
						},
					},
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s).*recurrence_value.*between 1 and 31.*`),
			},
			{
				// period is validated per period_unit: Month accepts 1 to 9 only
				Config: testAccConfig(map[string]interface{}{
					"period_unit":      "Month",
					"period":           "30",
					"recurrence_rules": REMOVEKEY,
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s).*when period_unit is Month, the value must be between 1 and 9.*`),
			},
			{
				// Weekly recurrence_value accepts days 0 to 6 only
				Config: testAccConfig(map[string]interface{}{
					"period_unit": "Day",
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Weekly",
							"recurrence_value": "7",
							"start_hour":       "8",
							"end_hour":         "12",
						},
					},
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s).*recurrence_value.*between 0 and 6.*`),
			},
			{
				// Monthly recurrence_value must be an ascending A-B range
				Config: testAccConfig(map[string]interface{}{
					"period_unit": "Day",
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Monthly",
							"recurrence_value": "5-1",
							"start_hour":       "8",
							"end_hour":         "12",
						},
					},
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s).*recurrence_value.*1 <= A <= B <= 31.*`),
			},
		},
	})
}

// Case 5: plan-time validation is deferred while period_unit or the rule hour bounds are unknown
func TestAccAliCloudEcsElasticityAssurance_timeDivisionUnknownValues(t *testing.T) {
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEcsElasticityAssurance%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceUnknownDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_amount":                     "1",
					"zone_ids":                            []string{"${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"},
					"instance_type":                       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
					"private_pool_options_name":           name,
					"private_pool_options_match_criteria": "Open",
					"period":                              "1",
					"period_unit":                         "${alicloud_vpc.default.id}",
					"assurance_times":                     "Unlimited",
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "${length(alicloud_vpc.default.id)}",
							"end_hour":         "${length(alicloud_vpc.default.id) + 24}",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period_unit": "Day",
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "${alicloud_vpc.default.id}",
							"end_hour":         "12",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// The values resolve only at apply time, so the plan-time check yields and Create
				// re-runs the same validation, which rejects the 1 hour span before calling the API.
				Config: testAccConfig(map[string]interface{}{
					"period_unit": "Day",
					"period":      "7",
					"recurrence_rules": []map[string]interface{}{
						{
							"recurrence_type":  "Daily",
							"recurrence_value": "1",
							"start_hour":       "${8 + length(alicloud_vpc.default.id) * 0}",
							"end_hour":         "${9 + length(alicloud_vpc.default.id) * 0}",
						},
					},
				}),
				ExpectError: regexp.MustCompile(`(?s).*at least 4 hours.*`),
			},
		},
	})
}

func AliCloudEcsElasticityAssuranceUnknownDependence1716(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_instance_types" "default" {
	instance_type_family = "ecs.c6"
}

resource "alicloud_vpc" "default" {
	vpc_name   = var.name
	cidr_block = "172.16.0.0/12"
}
`, name)
}

// The recurrence_rules boundary checks below drive the registered resource through its public
// factory and plugin-testing interfaces, so they stay in this file per the product package test
// standard instead of introducing a separate unit test file.

func elasticityAssuranceUnitConfig(periodUnit string, period int, rules ...map[string]interface{}) map[string]interface{} {
	config := map[string]interface{}{
		"instance_amount": 1,
		"instance_type":   []interface{}{"ecs.c6.large"},
		"zone_ids":        []interface{}{"cn-hangzhou-h"},
	}
	if periodUnit != "" {
		config["period_unit"] = periodUnit
	}
	if period > 0 {
		config["period"] = period
	}
	if len(rules) > 0 {
		rawRules := make([]interface{}, 0, len(rules))
		for _, rule := range rules {
			rawRules = append(rawRules, rule)
		}
		config["recurrence_rules"] = rawRules
	}
	return config
}

func elasticityAssuranceUnitRule(startHour, endHour int, recurrenceType, recurrenceValue string) map[string]interface{} {
	rule := map[string]interface{}{
		"start_hour":       startHour,
		"end_hour":         endHour,
		"recurrence_type":  recurrenceType,
		"recurrence_value": recurrenceValue,
	}
	if startHour < 0 {
		rule["start_hour"] = nil
	}
	if endHour < 0 {
		rule["end_hour"] = nil
	}
	if recurrenceType == "" {
		rule["recurrence_type"] = nil
	}
	if recurrenceValue == "" {
		rule["recurrence_value"] = nil
	}
	return rule
}

func TestUnitAliCloudEcsElasticityAssuranceSchema(t *testing.T) {
	elasticityAssuranceResource := acctest.Provider.ResourcesMap["alicloud_ecs_elasticity_assurance"]
	if elasticityAssuranceResource == nil {
		t.Fatal("the elasticity assurance resource is not registered in the provider")
	}
	for _, testCase := range []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{name: "valid day configuration", config: elasticityAssuranceUnitConfig("Day", 7, elasticityAssuranceUnitRule(8, 12, "Daily", "1"))},
		{name: "start hour zero is valid", config: elasticityAssuranceUnitConfig("Day", 7, elasticityAssuranceUnitRule(0, 4, "Daily", "1"))},
		{name: "day period above the schema range", config: elasticityAssuranceUnitConfig("Day", 366), wantErr: true},
		{name: "period unit must be a known value", config: elasticityAssuranceUnitConfig("Week", 7), wantErr: true},
		{name: "hour must be within a day", config: elasticityAssuranceUnitConfig("Day", 7, elasticityAssuranceUnitRule(8, 25, "Daily", "1")), wantErr: true},
		{name: "required attributes are enforced", config: map[string]interface{}{"period_unit": "Day"}, wantErr: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			diags := elasticityAssuranceResource.Validate(sdktf.NewResourceConfigRaw(testCase.config))
			if testCase.wantErr != diags.HasError() {
				t.Fatalf("unexpected validation result: %v", diags)
			}
		})
	}
}
