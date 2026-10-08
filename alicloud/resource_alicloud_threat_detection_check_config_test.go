package alicloud

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Case 检查项扫描策略配置20251209005 12031
// lintignore: AT001
func TestAccAliCloudThreatDetectionCheckConfig_basic12031(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_check_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionCheckConfigMap12031)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionCheckConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionCheckConfigBasicDependence12031)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
			testAccThreatDetectionCheckConfigBackup(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"vendors": []string{
						"ALIYUN"},
					"cycle_days": []string{
						"7", "1", "2"},
					"enable_add_check": "true",
					"start_time":       "12",
					"configure":        "not",
					"system_config":    "false",
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   848,
							"section_id": 62,
						},
						{
							"check_id":   648,
							"section_id": 68,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time":          "18",
						"enable_auto_check": "true",
						"cycle_days.#":      "3",
						"enable_add_check":  "true",
						"start_time":        "12",
						"configure":         "not",
						"system_config":     "false",
					}),
					resource.TestCheckResourceAttr(resourceId, "vendors.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.0.check_id", "848"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.0.section_id", "62"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.1.check_id", "648"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.1.section_id", "68"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configure": "not",
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configure": "not",
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configure": "not",
				}),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   648,
							"section_id": 68,
						},
						{
							"check_id":   848,
							"section_id": 62,
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configure": "not",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.0.check_id", "648"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.0.section_id", "68"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.1.check_id", "848"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.1.section_id", "62"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cycle_days": []string{
						"1", "7", "2",
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configure": "not",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cycle_days.#", "3"),
					resource.TestCheckResourceAttr(resourceId, "cycle_days.0", "1"),
					resource.TestCheckResourceAttr(resourceId, "cycle_days.1", "7"),
					resource.TestCheckResourceAttr(resourceId, "cycle_days.2", "2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vendors": []string{
						"ALIYUN", "AWS",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vendors.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "vendors.0", "ALIYUN"),
					resource.TestCheckResourceAttr(resourceId, "vendors.1", "AWS"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vendors": []string{
						"AWS", "ALIYUN",
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configure": "not",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "vendors.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "vendors.0", "AWS"),
					resource.TestCheckResourceAttr(resourceId, "vendors.1", "ALIYUN"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "12",
					"enable_auto_check": "false",
					"enable_add_check":  "false",
					"start_time":        "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time":          "12",
						"enable_auto_check": "false",
						"enable_add_check":  "false",
						"start_time":        "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_check": "true",
					"enable_add_check":  "true",
					"cycle_days": []string{
						"4"},
					"configure": "not",
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   648,
							"section_id": 68,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_check": "true",
						"cycle_days.#":      "1",
						"enable_add_check":  "true",
						"configure":         "not",
					}),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.0.check_id", "648"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.0.section_id", "68"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   648,
							"section_id": 68,
						},
						{
							"check_id":   848,
							"section_id": 62,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.1.check_id", "848"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.1.section_id", "62"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   648,
							"section_id": 68,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.0.check_id", "648"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.0.section_id", "68"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"configure", "system_config", "vendors"},
			},
		},
	})
}

var AlicloudThreatDetectionCheckConfigMap12031 = map[string]string{}

func AlicloudThreatDetectionCheckConfigBasicDependence12031(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test ThreatDetection CheckConfig. <<< Resource test cases, automatically generated.

// Check configuration is an account singleton and Delete only detaches state.
// Restore the readable configuration even when an acceptance step fails.
func testAccThreatDetectionCheckConfigBackup(t *testing.T) {
	meta, err := sharedClientForRegion("cn-hangzhou")
	if err != nil {
		t.Fatal("cannot configure check configuration backup client")
	}
	client := meta.(*connectivity.AliyunClient)
	if client.RegionId != "cn-hangzhou" {
		t.Fatal("unexpected check configuration test region")
	}
	identity, err := client.GetCallerIdentity()
	if err != nil || identity == nil || identity.AccountId == "" {
		t.Fatal("cannot verify check configuration test identity")
	}
	service := ThreatDetectionServiceV2{client}
	before, err := service.DescribeThreatDetectionCheckConfig("")
	if err != nil {
		t.Fatal("cannot back up check configuration")
	}
	fields := []string{"CycleDays", "EnableAutoCheck", "EnableAddCheck", "StartTime", "EndTime"}
	for _, field := range fields {
		if value, ok := before[field]; !ok || value == nil {
			t.Fatal("check configuration backup is incomplete")
		}
	}
	if value, ok := before["SelectedChecks"]; !ok || value == nil {
		t.Fatal("check configuration selection backup is incomplete")
	}
	selections := func(value interface{}) []interface{} {
		result := []interface{}{}
		for _, raw := range convertToInterfaceArray(value) {
			item := raw.(map[string]interface{})
			result = append(result, map[string]interface{}{"CheckId": item["CheckId"], "SectionId": item["SectionId"]})
		}
		return result
	}
	normalize := func(value interface{}, checks bool) []string {
		result := []string{}
		for _, raw := range convertToInterfaceArray(value) {
			if checks {
				item := raw.(map[string]interface{})
				result = append(result, fmt.Sprintf("%v:%v", item["CheckId"], item["SectionId"]))
			} else {
				result = append(result, fmt.Sprint(raw))
			}
		}
		sort.Strings(result)
		return result
	}
	restore := func(request map[string]interface{}) error {
		return resource.Retry(time.Minute, func() *resource.RetryError {
			_, err := client.RpcPost("Sas", "2018-12-03", "ChangeCheckConfig", map[string]interface{}{}, request, true)
			if err != nil {
				if NeedRetry(err) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
	}
	t.Cleanup(func() {
		current, err := service.DescribeThreatDetectionCheckConfig("")
		if err != nil {
			t.Error("cannot read check configuration for restoration")
			return
		}
		currentPairs := normalize(current["SelectedChecks"], true)
		if !reflect.DeepEqual(currentPairs, normalize(before["SelectedChecks"], true)) {
			for _, pair := range currentPairs {
				if pair != "848:62" && pair != "648:68" {
					t.Error("check configuration changed outside the test; refusing to overwrite it")
					return
				}
			}
		}
		request := map[string]interface{}{"RegionId": client.RegionId, "RemovedCheck": selections(current["SelectedChecks"])}
		for _, field := range fields {
			request[field] = before[field]
		}
		if err = restore(request); err != nil {
			t.Error("cannot clear test check configuration during restoration")
			return
		}
		delete(request, "RemovedCheck")
		request["AddedCheck"] = selections(before["SelectedChecks"])
		if err = restore(request); err != nil {
			t.Error("cannot restore original check configuration")
			return
		}
		restored, err := service.DescribeThreatDetectionCheckConfig("")
		if err != nil {
			t.Error("cannot verify restored check configuration")
			return
		}
		for _, field := range fields {
			same := fmt.Sprint(before[field]) == fmt.Sprint(restored[field])
			if field == "CycleDays" {
				same = reflect.DeepEqual(normalize(before[field], false), normalize(restored[field], false))
			}
			if !same {
				t.Error("restored check configuration differs from backup")
				return
			}
		}
		if !reflect.DeepEqual(normalize(before["SelectedChecks"], true), normalize(restored["SelectedChecks"], true)) {
			t.Error("restored check selections differ from backup")
		}
	})
}

func TestUnitThreatDetectionCheckConfigSelectedChecksDiff(t *testing.T) {
	pair := func(check, section interface{}) interface{} {
		return map[string]interface{}{"check_id": check, "section_id": section}
	}
	a, b, c := pair(848, 62), pair(648, 68), pair(848, 68)
	cases := []struct {
		name     string
		old, new interface{}
		changed  bool
	}{
		{"configured order change", []interface{}{a, b}, []interface{}{b, a}, true},
		{"same check different section", []interface{}{a, c}, []interface{}{c, a}, true},
		{"duplicate pairs", []interface{}{a, a, b}, []interface{}{b, a, a}, true},
		{"duplicate count changed", []interface{}{a, a, b}, []interface{}{a, b, b}, true},
		{"addition", []interface{}{a}, []interface{}{a, b}, true},
		{"removal", []interface{}{a, b}, []interface{}{b}, true},
		{"value change", []interface{}{a, b}, []interface{}{a, c}, true},
		{"clear", []interface{}{a, b}, []interface{}{}, true},
		{"empty", []interface{}{}, []interface{}{}, false},
		{"whole list unknown", []interface{}{a, b}, "74D93920-ED26-11E3-AC10-0800200C9A66", true},
		{"unknown reordered members", []interface{}{pair(0, 62), b}, []interface{}{b, pair("74D93920-ED26-11E3-AC10-0800200C9A66", 62)}, true},
		// The SDK returns no diff when an unknown integer normalizes to the
		// existing zero before comparing the list. Preserve that SDK baseline.
		{"unknown zero", []interface{}{pair(0, 62), b}, []interface{}{pair("74D93920-ED26-11E3-AC10-0800200C9A66", 62), b}, false},
		{"unknown section zero", []interface{}{pair(848, 0), b}, []interface{}{pair(848, "74D93920-ED26-11E3-AC10-0800200C9A66"), b}, false},
		{"unknown check", []interface{}{a, b}, []interface{}{pair("74D93920-ED26-11E3-AC10-0800200C9A66", 62), b}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := resourceAliCloudThreatDetectionCheckConfig()
			d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"selected_checks": tc.old})
			d.SetId("check-config")
			diff, err := r.Diff(d.State(), terraform.NewResourceConfigRaw(map[string]interface{}{"selected_checks": tc.new}), nil)
			if err != nil {
				t.Fatal(err)
			}
			changed := diff != nil && !diff.Empty()
			if changed != tc.changed {
				t.Fatalf("changed=%v want=%v diff=%#v", changed, tc.changed, diff)
			}
			if tc.name == "unknown reordered members" {
				for key, want := range map[string]string{"selected_checks.0.section_id": "68", "selected_checks.1.section_id": "62"} {
					if diff == nil || diff.Attributes[key] == nil || diff.Attributes[key].New != want {
						t.Fatalf("known section diff suppressed: %s want %s diff=%#v", key, want, diff)
					}
				}
			}
		})
	}
}

func TestUnitThreatDetectionCheckConfigReadOrder(t *testing.T) {
	pair := func(check, section interface{}) map[string]interface{} {
		return map[string]interface{}{"check_id": check, "section_id": section}
	}
	a, b, c := pair(848, 62), pair(648, 68), pair(848, 68)
	remote := []map[string]interface{}{b, c, a, a}
	got := threatDetectionCheckConfigOrder(remote, []interface{}{a, a, b, pair(999, 62)})
	want := []map[string]interface{}{a, a, b, c}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
	// API numbers arrive as float64, whereas persisted schema integers are ints.
	got = threatDetectionCheckConfigOrder([]map[string]interface{}{pair(float64(648), float64(68)), pair(float64(848), float64(62))}, []interface{}{a, b})
	if got[0]["check_id"] != float64(848) {
		t.Fatalf("API numeric values were not matched: %v", got)
	}
	if got := threatDetectionCheckConfigOrder(nil, []interface{}{a}); len(got) != 0 {
		t.Fatalf("removed member retained: %v", got)
	}
}
