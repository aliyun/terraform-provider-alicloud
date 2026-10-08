package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
							"check_id":   370,
							"section_id": 515,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time":          "18",
						"enable_auto_check": "true",
						"vendors.#":         "1",
						"cycle_days.#":      "3",
						"enable_add_check":  "true",
						"start_time":        "12",
						"configure":         "not",
						"system_config":     "false",
					}),
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
							"check_id":   23,
							"section_id": 11,
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
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"vendors": []string{
						"ALIYUN"},
					"cycle_days": []string{
						"4"},
					"enable_add_check": "true",
					"start_time":       "12",
					"configure":        "not",
					"system_config":    "false",
					// Config order deliberately differs from any sorted API
					// response order (check_id-first or section_id-first) to
					// reproduce the user drift scenario.
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time":          "18",
						"enable_auto_check": "true",
						"vendors.#":         "1",
						"cycle_days.#":      "1",
						"enable_add_check":  "true",
						"start_time":        "12",
						"configure":         "not",
						"system_config":     "false",
						"selected_checks.#": "2",
					}),
				),
			},
			{
				// Drift regression: re-planning the identical config must stay
				// empty even when the GetCheckConfig response order differs from
				// the config order (response order has no API contract).
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"vendors": []string{
						"ALIYUN"},
					"cycle_days": []string{
						"4"},
					"enable_add_check": "true",
					"start_time":       "12",
					"configure":        "not",
					"system_config":    "false",
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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

// Required by the TypeList Order Coverage gate (scripts/collection-order): proves
// that reordering the configurable cycle_days TypeList produces a diff and then
// converges after applying the reordered configuration. Members are weekdays 1-7
// per the ChangeCheckConfig API.
//
// selected_checks is pinned in every step: this resource is an account-level
// singleton whose Delete is a no-op, so a config without selected_checks leaves
// checks from earlier steps on the API side and the plan never converges.
//
// lintignore: AT001
func TestAccAliCloudThreatDetectionCheckConfig_cycleDaysOrder12031(t *testing.T) {
	resourceId := "alicloud_threat_detection_check_config.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, fmt.Sprintf("tfaccthreatdetection%d", acctest.RandIntRange(10000, 99999)), AlicloudThreatDetectionCheckConfigBasicDependence12031)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"enable_add_check":  "true",
					"start_time":        "12",
					"cycle_days": []string{
						"1", "3", "5"},
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cycle_days.#", "3"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"enable_add_check":  "true",
					"start_time":        "12",
					"cycle_days": []string{
						"5", "3", "1"},
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"enable_add_check":  "true",
					"start_time":        "12",
					"cycle_days": []string{
						"5", "3", "1"},
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cycle_days.#", "3"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "2"),
				),
			},
		},
	})
}

// Required by the TypeList Order Coverage gate (scripts/collection-order): proves
// that reordering the configurable vendors TypeList produces a diff and then
// converges after applying the reordered configuration. Members use the vendor
// values documented by the ChangeCheckConfig API (ALIYUN, TENCENT, AWS); vendors
// are not read back by the provider, so state mirrors the applied order.
//
// selected_checks is pinned in every step for the same singleton-residue reason
// as the cycle_days order test above.
//
// lintignore: AT001
func TestAccAliCloudThreatDetectionCheckConfig_vendorsOrder12031(t *testing.T) {
	resourceId := "alicloud_threat_detection_check_config.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, fmt.Sprintf("tfaccthreatdetection%d", acctest.RandIntRange(10000, 99999)), AlicloudThreatDetectionCheckConfigBasicDependence12031)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"enable_add_check":  "true",
					"start_time":        "12",
					"cycle_days": []string{
						"7"},
					"vendors": []string{
						"ALIYUN", "TENCENT"},
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cycle_days.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "vendors.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "2"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"enable_add_check":  "true",
					"start_time":        "12",
					"cycle_days": []string{
						"7"},
					"vendors": []string{
						"TENCENT", "ALIYUN"},
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"enable_add_check":  "true",
					"start_time":        "12",
					"cycle_days": []string{
						"7"},
					"vendors": []string{
						"TENCENT", "ALIYUN"},
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "cycle_days.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "vendors.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "selected_checks.#", "2"),
				),
			},
		},
	})
}

// Unit regression: the Read-side sort must produce a section_id-then-check_id
// deterministic order regardless of the API response order, and the set hash
// must be stable across int/float64 encodings while remaining distinct per
// (section_id, check_id) pair. Drift scenario: a config ordered
// [(370,515),(23,11)] and a GetCheckConfig response in a different order used
// to produce persistent plan diffs with the old TypeList + unsorted Read.
func TestUnitAliCloudThreatDetectionCheckConfigSelectedChecksOrder(t *testing.T) {
	input := []map[string]interface{}{
		{"check_id": 370, "section_id": 515},
		{"check_id": 23, "section_id": 11},
		{"check_id": float64(24), "section_id": float64(11)},
	}
	sortThreatDetectionCheckConfigSelectedChecks(input)

	expect := []struct {
		checkID   int
		sectionID int
	}{
		{23, 11},
		{24, 11},
		{370, 515},
	}
	for i, e := range expect {
		gotCheck := formatInt(input[i]["check_id"])
		gotSection := formatInt(input[i]["section_id"])
		if gotCheck != e.checkID || gotSection != e.sectionID {
			t.Fatalf("unexpected sorted order at %d: got (section=%d, check=%d), want (section=%d, check=%d)",
				i, gotSection, gotCheck, e.sectionID, e.checkID)
		}
	}

	hashInt := resourceAliCloudThreatDetectionCheckConfigSelectedChecksHash(map[string]interface{}{
		"check_id": 370, "section_id": 515,
	})
	hashFloat := resourceAliCloudThreatDetectionCheckConfigSelectedChecksHash(map[string]interface{}{
		"check_id": float64(370), "section_id": float64(515),
	})
	if hashInt != hashFloat {
		t.Fatalf("hash must normalize int and float64 identically: int=%d float64=%d", hashInt, hashFloat)
	}

	hashSwapped := resourceAliCloudThreatDetectionCheckConfigSelectedChecksHash(map[string]interface{}{
		"check_id": 515, "section_id": 370,
	})
	if hashInt == hashSwapped {
		t.Fatalf("hash must distinguish (section=515,check=370) from (section=370,check=515): both=%d", hashInt)
	}
}

// Test ThreatDetection CheckConfig. <<< Resource test cases, automatically generated.
