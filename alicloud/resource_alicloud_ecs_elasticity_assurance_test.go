package alicloud

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func TestAccAliCloudEcsElasticityAssurance_basic1716(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsElasticityAssuranceMap1716)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsElasticityAssurance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
				// The rules are stored as a set: declaring them in a different order must converge
				// without a diff instead of proposing to replace the assurance.
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
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceBasicDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		Providers: testAccProviders,
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
		},
	})
}

// Case 5: plan-time validation is deferred while period_unit or the rule hour bounds are unknown
func TestAccAliCloudEcsElasticityAssurance_timeDivisionUnknownValues(t *testing.T) {
	resourceId := "alicloud_ecs_elasticity_assurance.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sEcsElasticityAssurance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsElasticityAssuranceUnknownDependence1716)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		Providers: testAccProviders,
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

func TestUnitAliCloudEcsElasticityAssuranceRecurrenceRules(t *testing.T) {
	t.Run("nil element is rejected without panic", func(t *testing.T) {
		if _, err := expandEcsElasticityAssuranceRecurrenceRules([]interface{}{nil}); err == nil {
			t.Fatal("expected an error for a nil recurrence_rules element, got nil")
		}
		if err := validateEcsElasticityAssuranceRecurrenceRules("Day", []interface{}{nil}); err == nil {
			t.Fatal("expected an error for a nil recurrence_rules element, got nil")
		}
	})
	t.Run("non object element is rejected", func(t *testing.T) {
		if _, err := expandEcsElasticityAssuranceRecurrenceRules([]interface{}{"invalid"}); err == nil {
			t.Fatal("expected an error for a non object recurrence_rules element, got nil")
		}
	})
	t.Run("valid rules are expanded", func(t *testing.T) {
		rules := []interface{}{
			map[string]interface{}{
				"start_hour":       8,
				"end_hour":         12,
				"recurrence_type":  "Daily",
				"recurrence_value": "1",
			},
		}
		if err := validateEcsElasticityAssuranceRecurrenceRules("Day", rules); err != nil {
			t.Fatalf("unexpected validation error: %s", err)
		}
		expanded, err := expandEcsElasticityAssuranceRecurrenceRules(rules)
		if err != nil {
			t.Fatalf("unexpected expansion error: %s", err)
		}
		rule := expanded[0].(map[string]interface{})
		if rule["StartHour"] != 8 || rule["EndHour"] != 12 || rule["RecurrenceType"] != "Daily" || rule["RecurrenceValue"] != "1" {
			t.Fatalf("unexpected expanded rule: %#v", rule)
		}
	})
	t.Run("empty rules require period unit Day", func(t *testing.T) {
		if err := validateEcsElasticityAssuranceRecurrenceRules("Year", nil); err != nil {
			t.Fatalf("unexpected validation error: %s", err)
		}
		if err := validateEcsElasticityAssuranceRecurrenceRules("Day", nil); err == nil {
			t.Fatal("expected an error when period_unit is Day and recurrence_rules is empty, got nil")
		}
	})
	t.Run("rules require period unit Day", func(t *testing.T) {
		rules := []interface{}{
			map[string]interface{}{
				"start_hour":       8,
				"end_hour":         12,
				"recurrence_type":  "Daily",
				"recurrence_value": "1",
			},
		}
		if err := validateEcsElasticityAssuranceRecurrenceRules("Month", rules); err == nil {
			t.Fatal("expected an error when recurrence_rules is set without period_unit Day, got nil")
		}
	})
	t.Run("hour span shorter than four hours is rejected", func(t *testing.T) {
		rules := []interface{}{
			map[string]interface{}{
				"start_hour":       8,
				"end_hour":         10,
				"recurrence_type":  "Daily",
				"recurrence_value": "1",
			},
		}
		if err := validateEcsElasticityAssuranceRecurrenceRules("Day", rules); err == nil {
			t.Fatal("expected an error for an hour span shorter than four hours, got nil")
		}
	})
	t.Run("recurrence type read back from the api uses the documented casing", func(t *testing.T) {
		for _, tc := range []struct{ api, expected string }{
			{"DAILY", "Daily"},
			{"WEEKLY", "Weekly"},
			{"MONTHLY", "Monthly"},
		} {
			if got := convertEcsElasticityAssuranceRecurrenceTypeResponse(tc.api); got != tc.expected {
				t.Fatalf("expected %q to be normalized to %q, got %#v", tc.api, tc.expected, got)
			}
		}
	})
	t.Run("period is validated per period unit", func(t *testing.T) {
		for _, tc := range []struct {
			unit    string
			period  int
			wantErr bool
		}{
			{"Day", 7, false},
			{"Day", 365, false},
			{"Day", 366, true},
			{"Month", 9, false},
			{"Month", 10, true},
			{"Year", 5, false},
			{"Year", 6, true},
			{"", 5, false},
		} {
			err := validateEcsElasticityAssurancePeriod(tc.unit, tc.period)
			if tc.wantErr && err == nil {
				t.Fatalf("expected period %d with unit %q to be rejected", tc.period, tc.unit)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for period %d with unit %q: %s", tc.period, tc.unit, err)
			}
		}
	})
	t.Run("recurrence value format follows the recurrence type", func(t *testing.T) {
		for _, tc := range []struct {
			recurrenceType string
			value          string
			wantErr        bool
		}{
			{"Daily", "1", false},
			{"Daily", "31", false},
			{"Daily", "0", true},
			{"Daily", "32", true},
			{"Daily", "1,2", true},
			{"Weekly", "1,2", false},
			{"Weekly", "0,6", false},
			{"Weekly", "7", true},
			{"Monthly", "1-5", false},
			{"Monthly", "5-1", true},
			{"Monthly", "1", true},
			{"Monthly", "1-32", true},
		} {
			err := validateEcsElasticityAssuranceRecurrenceValue(tc.recurrenceType, tc.value)
			if tc.wantErr && err == nil {
				t.Fatalf("expected %q with recurrence_type %q to be rejected", tc.value, tc.recurrenceType)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for %q with recurrence_type %q: %s", tc.value, tc.recurrenceType, err)
			}
		}
	})
	t.Run("out of range recurrence value is rejected by the rule validator", func(t *testing.T) {
		rules := []interface{}{
			map[string]interface{}{
				"start_hour":       8,
				"end_hour":         12,
				"recurrence_type":  "Daily",
				"recurrence_value": "32",
			},
		}
		if err := validateEcsElasticityAssuranceRecurrenceRules("Day", rules); err == nil {
			t.Fatal("expected an error for an out of range recurrence_value, got nil")
		}
	})
}
