package alicloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// lintignore: AT001
func TestAccAliCloudReservedInstanceBasic(t *testing.T) {
	var v ecs.ReservedInstance

	resourceId := "alicloud_reserved_instance.default"
	ra := resourceAttrInit(resourceId, testAccReservedInstanceCheckMap)

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeReservedInstance")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsReservedInstanceConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceReservedInstanceBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{

				Config: testAccConfig(map[string]interface{}{
					"instance_type":     "ecs.g6.large",
					"instance_amount":   "1",
					"period_unit":       "Month",
					"offering_type":     "All Upfront",
					"name":              name,
					"description":       "ReservedInstance",
					"zone_id":           "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}",
					"scope":             "Zone",
					"period":            "1",
					"platform":          "Linux",
					"renewal_status":    "AutoRenewal",
					"auto_renew_period": "36",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"tags":              map[string]interface{}{"Created": "TF", "Foo": "Bar"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "ecs.g6.large",
						"instance_amount":   "1",
						"period_unit":       "Month",
						"offering_type":     "All Upfront",
						"name":              name,
						"description":       "ReservedInstance",
						"zone_id":           CHECKSET,
						"scope":             "Zone",
						"platform":          "Linux",
						"renewal_status":    "AutoRenewal",
						"auto_renew_period": "36",
						"resource_group_id": CHECKSET,
						"tags.%":            "2",
						"tags.Created":      "TF",
						"tags.Foo":          "Bar",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":     "true",
					"renewal_status": "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":     "true",
						"renewal_status": "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":     "false",
					"renewal_status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":     "false",
						"renewal_status": "Normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "ReservedInstanceChange",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "ReservedInstanceChange",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name,
					"description": "ReservedInstance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": "ReservedInstance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]interface{}{"Created": "TF1", "Foo": "Bar1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF1",
						"tags.Foo":     "Bar1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_period": "36",
					"renewal_status":    "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_period": "36",
						"renewal_status":    "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renewal_status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renewal_status": "Normal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "period_unit", "auto_renew_period"},
			},
		},
	})
}

// lintignore: AT001
func TestAccAliCloudReservedInstanceBasic1(t *testing.T) {
	var v ecs.ReservedInstance

	resourceId := "alicloud_reserved_instance.default"
	ra := resourceAttrInit(resourceId, testAccReservedInstanceCheckMap)

	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeReservedInstance")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsReservedInstanceConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceReservedInstanceBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{

				Config: testAccConfig(map[string]interface{}{
					"instance_type":          "${data.alicloud_instance_types.default.instance_types.0.id}",
					"instance_amount":        "1",
					"period":                 "1",
					"period_unit":            "Month",
					"offering_type":          "All Upfront",
					"reserved_instance_name": name,
					"description":            "ReservedInstance",
					"zone_id":                "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}",
					"scope":                  "Zone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":          CHECKSET,
						"instance_amount":        "1",
						"period":                 "1",
						"period_unit":            "Month",
						"offering_type":          "All Upfront",
						"reserved_instance_name": name,
						"description":            "ReservedInstance",
						"zone_id":                CHECKSET,
						"scope":                  "Zone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reserved_instance_name": name + "change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reserved_instance_name": name + "change",
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

var testAccReservedInstanceCheckMap = map[string]string{
	"allocation_status": "",
	"create_time":       CHECKSET,
	"expired_time":      CHECKSET,
	"start_time":        CHECKSET,
	"status":            CHECKSET,
	"operation_locks.#": "0",
}

func resourceReservedInstanceBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_instance_types" "default" {
		instance_type_family = "ecs.g7"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}
`, name)
}

// ---------------------------------------------------------------------------
// No-cloud compatibility regression tests
//
// These tests verify the compatibility contract between `auto_renew` (new
// boolean field) and `renewal_status` (legacy string field) without making
// real API calls.
//
// Contract under test:
//  1. Only legacy field configured → legacy field controls Create/Update.
//  2. Only auto_renew configured → true/false both change renewal status.
//  3. Read back-fill of auto_renew must NOT steal control from legacy config.
//  4. Both fields explicitly set → auto_renew wins; no perpetual diff.
//  5. auto_renew_period DiffSuppress consistent with effective switch.
//  6. Default/unknown and import → no active change to cloud renewal status.
// ---------------------------------------------------------------------------

// --- Core logic: autoRenewFromRawConfig with cty values ---------------------
//
// The production autoRenewFromConfig reads d.GetRawConfig(), which the SDK
// only populates during a real Terraform run (not in unit-test ResourceData
// built via Data(state, diff)).  We therefore test the pure inner function
// directly with cty values that mirror what GetRawConfig returns in
// production: an object whose attributes are cty.True, cty.False, null,
// or cty.UnknownVal(bool) (during plan when the value is computed but not
// yet known).

func TestAccAliCloudReservedInstanceAutoRenewFromRawConfig(t *testing.T) {
	cases := []struct {
		name      string
		rawConfig cty.Value
		wantVal   bool
		wantExpl  bool
	}{
		{
			name:      "explicit_true",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.True}),
			wantVal:   true,
			wantExpl:  true,
		},
		{
			name:      "explicit_false",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.False}),
			wantVal:   false,
			wantExpl:  true,
		},
		{
			name:      "absent_null",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.NullVal(cty.Bool)}),
			wantVal:   false,
			wantExpl:  false,
		},
		{
			name:      "unknown_value",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.UnknownVal(cty.Bool)}),
			wantVal:   false,
			wantExpl:  false,
		},
		{
			name:      "null_config_import",
			rawConfig: cty.NullVal(cty.EmptyObject),
			wantVal:   false,
			wantExpl:  false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, explicit := autoRenewFromRawConfig(tc.rawConfig)
			if explicit != tc.wantExpl {
				t.Fatalf("autoRenewFromRawConfig: explicit = %v, want %v", explicit, tc.wantExpl)
			}
			if v != tc.wantVal {
				t.Fatalf("autoRenewFromRawConfig: value = %v, want %v", v, tc.wantVal)
			}
		})
	}
}

// --- Diff behaviour tests ---------------------------------------------------
//
// These tests call r.Diff(ctx, state, config) and inspect the resulting
// terraform.InstanceDiff.  state.RawConfig is populated with the cty mirror of
// the config (riCtyConfig) so it propagates to diff.RawConfig via
// schemaMap.Diff, and the DiffSuppressFuncs read GetRawConfig() exactly the way
// Terraform core delivers it at plan time — so suppression is observable here,
// not only in the pure-function tests below.

// testRIState builds a simulated InstanceState.  Required field instance_type
// is pre-filled so schema Diff does not report missing-attribute errors.
func testRIState(attrs map[string]string) *terraform.InstanceState {
	if attrs == nil {
		attrs = map[string]string{}
	}
	if _, ok := attrs["instance_type"]; !ok {
		attrs["instance_type"] = "ecs.g6.large"
	}
	return &terraform.InstanceState{
		ID:         "ecsri-test-id",
		Attributes: attrs,
	}
}

func testRIDiff(t *testing.T, stateAttrs map[string]string, configRaw map[string]interface{}) *terraform.InstanceDiff {
	t.Helper()
	r := resourceAliCloudReservedInstance()
	state := testRIState(stateAttrs)
	// Complete legal user config: instance_type is Required+ForceNew and Terraform core
	// always carries it into a real plan-time config, so callers must include it. A partial
	// config (instance_type absent) flips the diff to RequiresNew; the SDK's
	// handleRequiresNew re-diff then nils d.state so GetRawConfig() returns null inside
	// DiffSuppressFunc, masking the renewal_status suppression — a test artifact, not a
	// production scenario. No user-explicit field is back-filled from state; the invariant
	// is asserted here so an incomplete config fails loudly instead of silently masking.
	state.RawConfig = riCtyConfig(configRaw)
	config := terraform.NewResourceConfigRaw(configRaw)
	diff, err := r.Diff(context.Background(), state, config, nil)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	if diff != nil {
		if diff.Destroy {
			t.Fatalf("complete legal config must not Destroy")
		}
		if diff.RequiresNew() {
			t.Fatalf("complete legal config must not RequireNew (instance_type must be present); a RequiresNew diff is not a production plan")
		}
	}
	return diff
}

// diffHasChange returns true when the diff carries a real Old→New change
// (not just NewComputed) for the given attribute key.
func diffHasChange(diff *terraform.InstanceDiff, key string) bool {
	if diff == nil || diff.Attributes == nil {
		return false
	}
	rd, ok := diff.Attributes[key]
	if !ok {
		return false
	}
	return rd.Old != rd.New && !rd.NewComputed
}

// diffKeyPresent returns true when the key exists in the diff's Attributes,
// regardless of whether it's a computed or real change.
func diffKeyPresent(diff *terraform.InstanceDiff, key string) bool {
	if diff == nil || diff.Attributes == nil {
		return false
	}
	_, ok := diff.Attributes[key]
	return ok
}

// Scenario 1: legacy双向 — when the user changes only renewal_status (legacy
// field) while state has a computed auto_renew from Read, the renewal_status
// change is NOT suppressed — the user's explicit legacy config controls the
// update.  Tested in both directions (AutoRenewal→Normal and Normal→AutoRenewal).

func TestAccAliCloudReservedInstanceDiffLegacyOnlyARToNormal(t *testing.T) {
	diff := testRIDiff(t,
		map[string]string{
			"renewal_status":    "AutoRenewal",
			"auto_renew":        "true",
			"auto_renew_period": "36",
		},
		map[string]interface{}{
			"instance_type":  "ecs.g6.large",
			"renewal_status": "Normal",
		},
	)
	if !diffHasChange(diff, "renewal_status") {
		t.Fatal("renewal_status change should NOT be suppressed (auto_renew not in config)")
	}
}

func TestAccAliCloudReservedInstanceDiffLegacyOnlyNormalToAR(t *testing.T) {
	diff := testRIDiff(t,
		map[string]string{
			"renewal_status":    "Normal",
			"auto_renew":        "false",
			"auto_renew_period": "12",
		},
		map[string]interface{}{
			"instance_type":  "ecs.g6.large",
			"renewal_status": "AutoRenewal",
		},
	)
	if !diffHasChange(diff, "renewal_status") {
		t.Fatal("renewal_status change should NOT be suppressed (auto_renew not in config)")
	}
}

// Scenario 2: explicit auto_renew — when both auto_renew and renewal_status are
// explicitly set, auto_renew takes precedence.  The renewal_status diff must be
// suppressed even when the config renewal_status value DIFFERS from state (real
// conflict, not a stable-state no-op).  This is the critical regression: without
// the IsKnown fix in autoRenewFromRawConfig, an unknown computed auto_renew
// value would be misclassified as explicit=false, incorrectly suppressing the
// renewal_status diff during plan.

func TestAccAliCloudReservedInstanceDiffConflictDualFields(t *testing.T) {
	diff := testRIDiff(t,
		map[string]string{
			"renewal_status":    "AutoRenewal",
			"auto_renew":        "false",
			"auto_renew_period": "36",
		},
		map[string]interface{}{
			// Genuine contradiction: auto_renew=true (on) contradicts
			// renewal_status=Normal (off). auto_renew must win, so the
			// provider sends RenewalStatus=AutoRenewal on update (verified
			// end-to-end by TestUnitReservedInstanceUpdate_ConflictAutoRenewTrueRenewalStatusNormal).
			"instance_type":  "ecs.g6.large",
			"auto_renew":     true,
			"renewal_status": "Normal",
		},
	)
	// auto_renew has no DiffSuppressFunc: state=false -> config=true is a real
	// diff (not a stable-state no-op like the previous version of this test,
	// which used auto_renew=true + renewal_status=AutoRenewal — mutually
	// consistent, not a real conflict).
	if !diffHasChange(diff, "auto_renew") {
		t.Fatal("auto_renew diff should be present: state=false, config=true is a real conflict")
	}
	// renewal_status would change AutoRenewal->Normal, but auto_renew is
	// explicit so renewalStatusDiffSuppress (reading GetRawConfig, now populated
	// via state.RawConfig) suppresses it.  With all=false the suppressed attr is
	// skipped from diff.Attributes entirely, so the key is absent.
	if diffKeyPresent(diff, "renewal_status") {
		t.Fatalf("renewal_status diff should be suppressed: auto_renew wins the conflict (attrs: %+v)", diff.Attributes)
	}
}

// Scenario 3: period-only — auto_renew_period diffs are suppressed when the
// effective renewal status is Normal (auto-renewal inactive), and NOT
// suppressed when AutoRenewal (active).  Tested via the fallback path
// (auto_renew not in config → GetRawConfig carries auto_renew=null →
// autoRenewFromRawConfig reports not-explicit → autoRenewPeriodDiffSuppress
// falls through to d.Get("renewal_status") from state).

func TestAccAliCloudReservedInstanceDiffPeriodSuppressedWhenNormal(t *testing.T) {
	diff := testRIDiff(t,
		map[string]string{
			"renewal_status":    "Normal",
			"auto_renew":        "false",
			"auto_renew_period": "12",
		},
		map[string]interface{}{
			"instance_type":     "ecs.g6.large",
			"auto_renew_period": 36,
		},
	)
	if diffKeyPresent(diff, "auto_renew_period") {
		t.Fatal("auto_renew_period diff should be suppressed when renewal_status=Normal")
	}
}

func TestAccAliCloudReservedInstanceDiffPeriodNotSuppressedWhenAR(t *testing.T) {
	diff := testRIDiff(t,
		map[string]string{
			"renewal_status":    "AutoRenewal",
			"auto_renew":        "true",
			"auto_renew_period": "36",
		},
		map[string]interface{}{
			"instance_type":     "ecs.g6.large",
			"auto_renew_period": 12,
		},
	)
	if !diffHasChange(diff, "auto_renew_period") {
		t.Fatal("auto_renew_period change should NOT be suppressed when renewal_status=AutoRenewal")
	}
}

// Scenario 4: Read后冲突下一plan — after import/Read, state has computed
// auto_renew and auto_renew_period.  A subsequent plan with a matching
// renewal_status config does not produce spurious diffs for the renewal fields.

func TestAccAliCloudReservedInstanceDiffPlanNoDiffAfterImport(t *testing.T) {
	diff := testRIDiff(t,
		map[string]string{
			"renewal_status":    "AutoRenewal",
			"auto_renew":        "true",
			"auto_renew_period": "36",
		},
		map[string]interface{}{
			"instance_type":  "ecs.g6.large",
			"renewal_status": "AutoRenewal",
		},
	)
	for _, key := range []string{"renewal_status", "auto_renew", "auto_renew_period"} {
		if diffHasChange(diff, key) {
			t.Fatalf("%s has spurious diff: %+v", key, diff.Attributes[key])
		}
	}
}

// Scenario 5: 导入no-op — a plan with only renewal_status matching state and
// no other renewal fields does not emit renewal-related diffs (configuration
// defaults must not change cloud renewal state).

func TestAccAliCloudReservedInstanceDiffPlanNoDiffDefault(t *testing.T) {
	diff := testRIDiff(t,
		map[string]string{
			"renewal_status":    "Normal",
			"auto_renew":        "false",
			"auto_renew_period": "12",
		},
		map[string]interface{}{
			"instance_type":  "ecs.g6.large",
			"renewal_status": "Normal",
		},
	)
	for _, key := range []string{"renewal_status", "auto_renew", "auto_renew_period"} {
		if diffHasChange(diff, key) {
			t.Fatalf("%s has spurious diff: %+v", key, diff.Attributes[key])
		}
	}
}

// --- Pure-function DiffSuppress tests ---------------------------------------
//
// The SDK does not populate GetRawConfig() on ResourceData built via
// Data(state, diff).  We test the extracted pure functions directly with cty
// values that mirror what GetRawConfig returns in a real Terraform plan,
// including the unknown-value case that the IsKnown fix addresses.

func TestAccAliCloudReservedInstanceRenewalStatusDiffSuppress(t *testing.T) {
	cases := []struct {
		name      string
		rawConfig cty.Value
		want      bool
	}{
		{
			name:      "suppressed_when_auto_renew_explicit_true",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.True}),
			want:      true,
		},
		{
			name:      "suppressed_when_auto_renew_explicit_false",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.False}),
			want:      true,
		},
		{
			name:      "not_suppressed_when_auto_renew_absent",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.NullVal(cty.Bool)}),
			want:      false,
		},
		{
			name:      "not_suppressed_when_auto_renew_unknown",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.UnknownVal(cty.Bool)}),
			want:      false,
		},
		{
			name:      "not_suppressed_when_config_null_import",
			rawConfig: cty.NullVal(cty.EmptyObject),
			want:      false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := renewalStatusDiffSuppress(tc.rawConfig); got != tc.want {
				t.Fatalf("renewalStatusDiffSuppress = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAccAliCloudReservedInstanceAutoRenewPeriodDiffSuppress(t *testing.T) {
	cases := []struct {
		name          string
		rawConfig     cty.Value
		renewalStatus string
		want          bool
	}{
		{
			name:          "suppressed_when_normal_and_no_auto_renew",
			rawConfig:     cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.NullVal(cty.Bool)}),
			renewalStatus: "Normal",
			want:          true,
		},
		{
			name:          "not_suppressed_when_autorenewal_and_no_auto_renew",
			rawConfig:     cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.NullVal(cty.Bool)}),
			renewalStatus: "AutoRenewal",
			want:          false,
		},
		{
			name:          "suppressed_when_explicit_auto_renew_false",
			rawConfig:     cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.False}),
			renewalStatus: "AutoRenewal",
			want:          true,
		},
		{
			name:          "not_suppressed_when_explicit_auto_renew_true",
			rawConfig:     cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.True}),
			renewalStatus: "Normal",
			want:          false,
		},
		{
			name:          "suppressed_when_unknown_auto_renew_and_normal",
			rawConfig:     cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.UnknownVal(cty.Bool)}),
			renewalStatus: "Normal",
			want:          true,
		},
		{
			name:          "not_suppressed_when_unknown_auto_renew_and_autorenewal",
			rawConfig:     cty.ObjectVal(map[string]cty.Value{"auto_renew": cty.UnknownVal(cty.Bool)}),
			renewalStatus: "AutoRenewal",
			want:          false,
		},
		{
			name:          "suppressed_when_null_config_and_normal_state",
			rawConfig:     cty.NullVal(cty.EmptyObject),
			renewalStatus: "Normal",
			want:          true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := autoRenewPeriodDiffSuppress(tc.rawConfig, tc.renewalStatus); got != tc.want {
				t.Fatalf("autoRenewPeriodDiffSuppress = %v, want %v", got, tc.want)
			}
		})
	}
}

// The no-cloud DiffSuppress/Diff tests in resource_alicloud_reserved_instance_test.go
// only exercise the helper functions and the SDK r.Diff() harness; they never invoke
// the real resourceAliCloudReservedInstance Create/Update/Read. The tests below build a
// real AliyunClient whose ECS endpoint is an httptest.Server, construct
// production-faithful schema.ResourceData (state + config + RawConfig populated the way
// Terraform core delivers it at plan/apply time) and call the actual CRUD functions,
// asserting the HTTP query parameters the provider emits — without any real cloud cost.
//
// The mock is STATEFUL: ModifyReservedInstanceAutoRenewAttribute updates the mock cloud
// renewal state from the real request (partial update; unsent fields keep their value,
// mirroring the API), and DescribeReservedInstanceAutoRenewAttribute reads it back. This
// lets a single client+resource exercise a real Update -> (auto) Read -> next SDKDiff
// closed loop, proving the Update actually changed the cloud state and the next plan is
// stable — rather than pre-setting the desired response to fake a fix.

const riMockID = "ri-mock"

// riDescribeInstancesJSON is the fixed DescribeReservedInstances response: a single
// Active reserved instance with the instance fields Read copies into state.
const riDescribeInstancesJSON = `{"RequestId":"mock","TotalCount":1,"ReservedInstances":{"ReservedInstance":[{"ReservedInstanceId":"ri-mock","Status":"Active","InstanceType":"ecs.g6.large","InstanceAmount":1,"Scope":"Region","OfferingType":"All Upfront","Platform":"Linux","OperationLocks":{"OperationLock":[]},"Tags":{"Tag":[]}}]}}`

// riCall captures one request the provider sent to the mock ECS endpoint.
type riCall struct {
	action string
	form   url.Values
}

// riMockState is the mutable cloud renewal state of the single test reserved instance.
// ModifyReservedInstanceAutoRenewAttribute updates it from the real request (only the
// fields the client actually sent); DescribeReservedInstanceAutoRenewAttribute reads it
// back. This mirrors the real API's partial-update semantics, so a period-only or
// explicit-false update cannot be "proven" by pre-setting the desired response.
type riMockState struct {
	mu          sync.Mutex
	renewStatus string
	duration    int
	periodUnit  string
}

// riStatefulClient starts a stateful httptest.Server and returns an AliyunClient routed
// at it, the mutable mock state, and a recorder of every captured request. The initial
// renewStatus/duration/periodUnit represent the resource's current cloud state before
// the test acts on it; Create/Modify overwrite only the fields the real request sends.
func riStatefulClient(t *testing.T, renewStatus string, duration int, periodUnit string) (*connectivity.AliyunClient, *riMockState, *[]riCall) {
	t.Helper()
	st := &riMockState{renewStatus: renewStatus, duration: duration, periodUnit: periodUnit}
	var calls []riCall
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		form := url.Values{}
		for k, vs := range r.Form {
			form[k] = vs
		}
		// tea-rpc carries the signed parameters on the URL query for POST.
		for k, vs := range r.URL.Query() {
			if _, ok := form[k]; !ok {
				form[k] = vs
			}
		}
		action := form.Get("Action")
		st.mu.Lock()
		calls = append(calls, riCall{action: action, form: form})
		// Partial update from the captured real request: only sent fields change.
		switch action {
		case "PurchaseReservedInstancesOffering":
			if v, ok := form["AutoRenew"]; ok {
				if v[0] == "true" {
					st.renewStatus = "AutoRenewal"
				} else {
					st.renewStatus = "Normal"
				}
			}
			if v, ok := form["AutoRenewPeriod"]; ok {
				n, _ := strconv.Atoi(v[0])
				if n == 1 {
					st.duration = 1
					st.periodUnit = "Month"
				} else {
					st.duration = n / 12
					st.periodUnit = "Year"
				}
			}
		case "ModifyReservedInstanceAutoRenewAttribute":
			if v, ok := form["RenewalStatus"]; ok {
				st.renewStatus = v[0]
			}
			if v, ok := form["Period"]; ok {
				st.duration, _ = strconv.Atoi(v[0])
			}
			if v, ok := form["PeriodUnit"]; ok {
				st.periodUnit = v[0]
			}
		}
		rs, dur, pu := st.renewStatus, st.duration, st.periodUnit
		st.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "PurchaseReservedInstancesOffering":
			_, _ = io.WriteString(w, fmt.Sprintf(`{"ReservedInstanceIdSets":{"ReservedInstanceId":[%q]},"RequestId":"mock"}`, riMockID))
		case "DescribeReservedInstances":
			_, _ = io.WriteString(w, riDescribeInstancesJSON)
		case "DescribeReservedInstanceAutoRenewAttribute":
			_, _ = io.WriteString(w, fmt.Sprintf(`{"RequestId":"mock","ReservedInstanceRenewAttributes":{"ReservedInstanceRenewAttribute":[{"ReservedInstanceId":%q,"RenewalStatus":%q,"Duration":%d,"PeriodUnit":%q}]}}`, riMockID, rs, dur, pu))
		default:
			_, _ = io.WriteString(w, `{"RequestId":"mock"}`)
		}
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("create credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("init client: %v", err)
	}
	endpoints.Store("ecs", endpoint)
	return client, st, &calls
}

func riActions(calls *[]riCall) []string {
	var acts []string
	for _, c := range *calls {
		acts = append(acts, c.action)
	}
	return acts
}

func riRequireCall(t *testing.T, calls *[]riCall, action string) riCall {
	t.Helper()
	for _, c := range *calls {
		if c.action == action {
			return c
		}
	}
	t.Fatalf("expected an HTTP %s call, got actions: %v", action, riActions(calls))
	return riCall{}
}

func riFormEq(t *testing.T, form url.Values, key, want string) {
	t.Helper()
	if got := form.Get(key); got != want {
		t.Fatalf("%s = %q, want %q (full form: %v)", key, got, want, form)
	}
}

// riFormAbsent asserts the key is not present in the captured form. Used for
// compatibility semantics: a period-only update must not tamper with a non-explicit
// renewal status, so RenewalStatus is absent from the request.
func riFormAbsent(t *testing.T, form url.Values, key string) {
	t.Helper()
	if vs, ok := form[key]; ok {
		t.Fatalf("%s should be absent from the form (got %q); unsent fields must not be tampered with (full form: %v)", key, vs, form)
	}
}

// riCtyConfig builds a cty object mirroring GetRawConfig() as Terraform core delivers it
// at plan time. auto_renew and renewal_status are always present (null when the user did
// not set them) so autoRenewFromRawConfig's GetAttr("auto_renew") and
// renewalStatusFromConfig's GetAttr("renewal_status") behave exactly like production
// instead of erroring on a missing attribute.
func riCtyConfig(configRaw map[string]interface{}) cty.Value {
	vals := map[string]cty.Value{}
	if v, ok := configRaw["auto_renew"]; ok {
		vals["auto_renew"] = cty.BoolVal(v.(bool))
	} else {
		vals["auto_renew"] = cty.NullVal(cty.Bool)
	}
	if v, ok := configRaw["renewal_status"]; ok {
		vals["renewal_status"] = cty.StringVal(v.(string))
	} else {
		vals["renewal_status"] = cty.NullVal(cty.String)
	}
	for k, v := range configRaw {
		if k == "auto_renew" || k == "renewal_status" {
			continue
		}
		switch tv := v.(type) {
		case bool:
			vals[k] = cty.BoolVal(tv)
		case string:
			vals[k] = cty.StringVal(tv)
		case int:
			vals[k] = cty.NumberIntVal(int64(tv))
		}
	}
	return cty.ObjectVal(vals)
}

// riCreateData builds the ResourceData passed to Create: state carries the configured
// values (Create reads via d.GetOk) and RawConfig carries the user's explicit config so
// autoRenewFromConfig's GetRawConfig() is populated.
func riCreateData(t *testing.T, configRaw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := resourceAliCloudReservedInstance()
	attrs := map[string]string{"instance_type": "ecs.g6.large"}
	for k, v := range configRaw {
		attrs[k] = fmt.Sprint(v)
	}
	state := &terraform.InstanceState{ID: "", Attributes: attrs, RawConfig: riCtyConfig(configRaw)}
	return r.Data(state)
}

// riUpdateData builds the ResourceData passed to Update from a COMPLETE legal user config
// (instance_type is Required+ForceNew; Terraform core always carries it into a real
// plan-time config). It asserts the resulting diff neither Destroys nor RequiresNew at
// the Update entry: a RequiresNew diff is a recreate, not a production Update apply, and
// would nil d.state inside the SDK's handleRequiresNew re-diff, breaking the
// GetRawConfig()-based DiffSuppress. Callers must pass a complete config (instance_type
// present); no user-explicit field is back-filled from state.
func riUpdateData(t *testing.T, stateAttrs map[string]string, configRaw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := resourceAliCloudReservedInstance()
	state := testRIState(stateAttrs)
	state.ID = riMockID
	state.RawConfig = riCtyConfig(configRaw)
	config := terraform.NewResourceConfigRaw(configRaw)
	diff, err := r.Diff(context.Background(), state, config, nil)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if diff != nil {
		if diff.Destroy {
			t.Fatalf("Update entry: diff must not Destroy")
		}
		if diff.RequiresNew() {
			t.Fatalf("Update entry: complete legal config must not RequireNew (instance_type must be present); a RequiresNew diff is not a production Update apply")
		}
	}
	d, err := schema.InternalMap(r.SchemaMap()).Data(state, diff)
	if err != nil {
		t.Fatalf("Data: %v", err)
	}
	return d
}

// riReadData builds the ResourceData passed to Read for an existing instance. Not used
// for import — import uses an ID-only empty state (see TestUnitReservedInstanceImport_NoOp).
func riReadData(t *testing.T, stateAttrs map[string]string) *schema.ResourceData {
	r := resourceAliCloudReservedInstance()
	state := testRIState(stateAttrs)
	state.ID = riMockID
	return r.Data(state)
}

// riDiff builds a production-faithful diff from a COMPLETE legal config: state.RawConfig is
// populated with the config cty so DiffSuppressFuncs read GetRawConfig() the way Terraform
// core does at plan time. No state field is back-filled into config — callers pass the
// full user config.
func riDiff(t *testing.T, state *terraform.InstanceState, configRaw map[string]interface{}) *terraform.InstanceDiff {
	t.Helper()
	r := resourceAliCloudReservedInstance()
	if state == nil {
		state = testRIState(nil)
	}
	state.RawConfig = riCtyConfig(configRaw)
	config := terraform.NewResourceConfigRaw(configRaw)
	diff, err := r.Diff(context.Background(), state, config, nil)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	return diff
}

// riAssertStable asserts a post-Update next-SDKDiff is stable: no Destroy, no RequiresNew,
// and no spurious change for the listed renewal keys. Used after a real Update -> Read
// closed loop on the same stateful client+resource.
func riAssertStable(t *testing.T, diff *terraform.InstanceDiff, keys []string) {
	t.Helper()
	if diff != nil {
		if diff.Destroy {
			t.Fatalf("post-Update next diff must not Destroy")
		}
		if diff.RequiresNew() {
			t.Fatalf("post-Update next diff must not RequireNew (complete legal config)")
		}
	}
	for _, k := range keys {
		if diffHasChange(diff, k) {
			t.Fatalf("%s has a spurious diff after Update->Read: %+v", k, diff.Attributes[k])
		}
	}
}

// --- Create: explicit auto_renew and legacy renewal_status -------------------

func TestUnitReservedInstanceCreate_AutoRenewExplicitTrue(t *testing.T) {
	client, _, calls := riStatefulClient(t, "Normal", 1, "Year")
	d := riCreateData(t, map[string]interface{}{
		"instance_type":   "ecs.g6.large",
		"offering_type":   "All Upfront",
		"period":          1,
		"period_unit":     "Month",
		"platform":        "Linux",
		"scope":           "Region",
		"instance_amount": 1,
		"auto_renew":      true,
	})
	if err := resourceAliCloudReservedInstanceCreate(d, client); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if d.Id() != riMockID {
		t.Fatalf("Create did not set the reserved instance id: got %q", d.Id())
	}
	c := riRequireCall(t, calls, "PurchaseReservedInstancesOffering")
	riFormEq(t, c.form, "AutoRenew", "true")
}

func TestUnitReservedInstanceCreate_AutoRenewExplicitFalse(t *testing.T) {
	client, _, calls := riStatefulClient(t, "Normal", 1, "Year")
	d := riCreateData(t, map[string]interface{}{
		"instance_type":   "ecs.g6.large",
		"offering_type":   "All Upfront",
		"period":          1,
		"period_unit":     "Month",
		"platform":        "Linux",
		"scope":           "Region",
		"instance_amount": 1,
		"auto_renew":      false,
	})
	if err := resourceAliCloudReservedInstanceCreate(d, client); err != nil {
		t.Fatalf("Create: %v", err)
	}
	c := riRequireCall(t, calls, "PurchaseReservedInstancesOffering")
	riFormEq(t, c.form, "AutoRenew", "false")
}

func TestUnitReservedInstanceCreate_LegacyRenewalStatusAutoRenewal(t *testing.T) {
	client, _, calls := riStatefulClient(t, "Normal", 1, "Year")
	// auto_renew is intentionally absent: the legacy renewal_status field drives AutoRenew,
	// exercising the else branch in Create.
	d := riCreateData(t, map[string]interface{}{
		"instance_type":   "ecs.g6.large",
		"offering_type":   "All Upfront",
		"period":          1,
		"period_unit":     "Month",
		"platform":        "Linux",
		"scope":           "Region",
		"instance_amount": 1,
		"renewal_status":  "AutoRenewal",
	})
	if err := resourceAliCloudReservedInstanceCreate(d, client); err != nil {
		t.Fatalf("Create: %v", err)
	}
	c := riRequireCall(t, calls, "PurchaseReservedInstancesOffering")
	riFormEq(t, c.form, "AutoRenew", "true")
}

// --- Update: real Update -> Read -> next SDKDiff closed loops ----------------
//
// Each Update test drives a complete legal user config (instance_type Required, plus the
// renewal fields under test) through the production Update on one stateful client+resource.
// Update sends a real Modify (asserted), the stateful mock reflects it, Update's own Read
// re-reads the updated cloud state, and the next SDKDiff with the SAME config must be
// stable: no spurious change, no Destroy, no RequiresNew.

func TestUnitReservedInstanceUpdate_LegacyAutoRenewalToNormal(t *testing.T) {
	// Before cloud state (real): AutoRenewal, Duration 3 / Year (=> auto_renew_period 36).
	client, st, calls := riStatefulClient(t, "AutoRenewal", 3, "Year")
	newConfig := map[string]interface{}{
		"instance_type":     "ecs.g6.large",
		"renewal_status":    "Normal",
		"auto_renew_period": 12,
	}
	d := riUpdateData(t,
		map[string]string{
			"renewal_status":    "AutoRenewal",
			"auto_renew":        "true",
			"auto_renew_period": "36",
		},
		newConfig,
	)
	if err := resourceAliCloudReservedInstanceUpdate(d, client); err != nil {
		t.Fatalf("Update: %v", err)
	}
	// auto_renew is not in config (not explicit) so the legacy renewal_status=Normal
	// drives RenewalStatus. Under Normal the auto_renew_period change (36->12) is
	// suppressed by autoRenewPeriodDiffSuppress, so the writer keeps the state value
	// 36 and GetOk("auto_renew_period") returns (36, true) => Period=36/12=3, not the
	// configured 1. This is the authorized compat behavior: autoRenewPeriodDiffSuppress
	// keeps the prior period 36, the actual request carries Period=3/PeriodUnit=Year,
	// and the stateful mock's renewal period stays unchanged (Duration 3/Year), so the
	// closed loop is stable. This is observed mock behavior, not a claim that the real
	// cloud ignores Period under Normal.
	c := riRequireCall(t, calls, "ModifyReservedInstanceAutoRenewAttribute")
	riFormEq(t, c.form, "RenewalStatus", "Normal")
	riFormEq(t, c.form, "Period", "3")
	riFormEq(t, c.form, "PeriodUnit", "Year")
	st.mu.Lock()
	gotStatus, gotDur := st.renewStatus, st.duration
	st.mu.Unlock()
	if gotStatus != "Normal" || gotDur != 3 {
		t.Fatalf("stateful mock after Modify = (%q,%d), want (Normal,3)", gotStatus, gotDur)
	}
	post := d.State()
	if got := post.Attributes["renewal_status"]; got != "Normal" {
		t.Fatalf("post-Read renewal_status = %q, want Normal", got)
	}
	if got := post.Attributes["auto_renew"]; got != "false" {
		t.Fatalf("post-Read auto_renew = %q, want false", got)
	}
	if got := post.Attributes["auto_renew_period"]; got != "36" {
		t.Fatalf("post-Read auto_renew_period = %q, want 36 (Duration 3/Year from the unchanged mock)", got)
	}
	riAssertStable(t, riDiff(t, post, newConfig), []string{"renewal_status", "auto_renew", "auto_renew_period"})
}

func TestUnitReservedInstanceUpdate_LegacyNormalToAutoRenewal(t *testing.T) {
	client, st, calls := riStatefulClient(t, "Normal", 1, "Year")
	newConfig := map[string]interface{}{
		"instance_type":     "ecs.g6.large",
		"renewal_status":    "AutoRenewal",
		"auto_renew_period": 36,
	}
	d := riUpdateData(t,
		map[string]string{
			"renewal_status":    "Normal",
			"auto_renew":        "false",
			"auto_renew_period": "12",
		},
		newConfig,
	)
	if err := resourceAliCloudReservedInstanceUpdate(d, client); err != nil {
		t.Fatalf("Update: %v", err)
	}
	c := riRequireCall(t, calls, "ModifyReservedInstanceAutoRenewAttribute")
	riFormEq(t, c.form, "RenewalStatus", "AutoRenewal")
	riFormEq(t, c.form, "Period", "3")
	riFormEq(t, c.form, "PeriodUnit", "Year")
	st.mu.Lock()
	gotStatus, gotDur := st.renewStatus, st.duration
	st.mu.Unlock()
	if gotStatus != "AutoRenewal" || gotDur != 3 {
		t.Fatalf("stateful mock after Modify = (%q,%d), want (AutoRenewal,3)", gotStatus, gotDur)
	}
	post := d.State()
	if got := post.Attributes["renewal_status"]; got != "AutoRenewal" {
		t.Fatalf("post-Read renewal_status = %q, want AutoRenewal", got)
	}
	if got := post.Attributes["auto_renew"]; got != "true" {
		t.Fatalf("post-Read auto_renew = %q, want true", got)
	}
	if got := post.Attributes["auto_renew_period"]; got != "36" {
		t.Fatalf("post-Read auto_renew_period = %q, want 36", got)
	}
	riAssertStable(t, riDiff(t, post, newConfig), []string{"renewal_status", "auto_renew", "auto_renew_period"})
}

func TestUnitReservedInstanceUpdate_ExplicitAutoRenewFalse(t *testing.T) {
	// Before cloud state: AutoRenewal, Duration 3 / Year (=> 36). The user changes ONLY
	// auto_renew to false. Because auto_renew is explicit-false, autoRenewPeriodDiffSuppress
	// suppresses the period diff, so the writer keeps the state value 36 and GetOk returns
	// (36, true) => the provider sends Period=3, PeriodUnit=Year alongside RenewalStatus=
	// Normal. autoRenewPeriodDiffSuppress keeps the prior period 36, the actual request
	// carries Period=3/PeriodUnit=Year, and the stateful mock proves the renewal period is
	// left UNCHANGED (still 3/Year) — observed mock behavior, not a claim that the real
	// cloud ignores Period under Normal. A non-stateful mock reading a fixed canned
	// response could not prove the period is preserved.
	client, st, calls := riStatefulClient(t, "AutoRenewal", 3, "Year")
	newConfig := map[string]interface{}{
		"instance_type": "ecs.g6.large",
		"auto_renew":    false,
	}
	d := riUpdateData(t,
		map[string]string{
			"renewal_status":    "AutoRenewal",
			"auto_renew":        "true",
			"auto_renew_period": "36",
		},
		newConfig,
	)
	if err := resourceAliCloudReservedInstanceUpdate(d, client); err != nil {
		t.Fatalf("Update: %v", err)
	}
	c := riRequireCall(t, calls, "ModifyReservedInstanceAutoRenewAttribute")
	riFormEq(t, c.form, "RenewalStatus", "Normal")
	riFormEq(t, c.form, "Period", "3")
	riFormEq(t, c.form, "PeriodUnit", "Year")
	st.mu.Lock()
	gotStatus, gotDur, gotPU := st.renewStatus, st.duration, st.periodUnit
	st.mu.Unlock()
	if gotStatus != "Normal" || gotDur != 3 || gotPU != "Year" {
		t.Fatalf("stateful mock after Modify = (%q,%d,%q), want (Normal,3,Year) — period must be unchanged", gotStatus, gotDur, gotPU)
	}
	post := d.State()
	if got := post.Attributes["renewal_status"]; got != "Normal" {
		t.Fatalf("post-Read renewal_status = %q, want Normal", got)
	}
	if got := post.Attributes["auto_renew"]; got != "false" {
		t.Fatalf("post-Read auto_renew = %q, want false", got)
	}
	if got := post.Attributes["auto_renew_period"]; got != "36" {
		t.Fatalf("post-Read auto_renew_period = %q, want 36 (period untouched)", got)
	}
	riAssertStable(t, riDiff(t, post, newConfig), []string{"renewal_status", "auto_renew", "auto_renew_period"})
}

func TestUnitReservedInstanceUpdate_PeriodOnlyPreservesAutoRenewal(t *testing.T) {
	// Initial cloud state (legacy/computed): AutoRenewal, Duration 3 / Year (=> 36). The
	// config carries ONLY auto_renew_period=12 — neither renewal_status nor auto_renew is
	// explicit. Per ModifyReservedInstanceAutoRenewAttribute semantics (RenewalStatus is
	// optional; the API gives no documented default-override rule, and Period/PeriodUnit are
	// the period parameters), a period-only update sends just the period and must NOT tamper
	// with the non-explicit renewal status: RenewalStatus is absent so the cloud keeps its
	// current renewal state (preserved as AutoRenewal) instead of forcing it.
	// https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-modifyreservedinstanceautorenewattribute
	client, st, calls := riStatefulClient(t, "AutoRenewal", 3, "Year")
	newConfig := map[string]interface{}{
		"instance_type":     "ecs.g6.large",
		"auto_renew_period": 12,
	}
	d := riUpdateData(t,
		map[string]string{
			"renewal_status":    "AutoRenewal",
			"auto_renew":        "true",
			"auto_renew_period": "36",
		},
		newConfig,
	)
	if err := resourceAliCloudReservedInstanceUpdate(d, client); err != nil {
		t.Fatalf("Update: %v", err)
	}
	c := riRequireCall(t, calls, "ModifyReservedInstanceAutoRenewAttribute")
	riFormEq(t, c.form, "Period", "1")
	riFormEq(t, c.form, "PeriodUnit", "Year")
	// Neither auto_renew nor renewal_status is explicit => RenewalStatus is absent.
	riFormAbsent(t, c.form, "RenewalStatus")
	// Stateful mock: period updated to 1/Year; renewal status preserved as AutoRenewal.
	st.mu.Lock()
	gotStatus, gotDur := st.renewStatus, st.duration
	st.mu.Unlock()
	if gotStatus != "AutoRenewal" || gotDur != 1 {
		t.Fatalf("stateful mock after period-only Modify = (%q,%d), want (AutoRenewal,1) — status must be preserved", gotStatus, gotDur)
	}
	post := d.State()
	if got := post.Attributes["renewal_status"]; got != "AutoRenewal" {
		t.Fatalf("post-Read renewal_status = %q, want AutoRenewal (preserved by period-only update)", got)
	}
	if got := post.Attributes["auto_renew"]; got != "true" {
		t.Fatalf("post-Read auto_renew = %q, want true (preserved)", got)
	}
	if got := post.Attributes["auto_renew_period"]; got != "12" {
		t.Fatalf("post-Read auto_renew_period = %q, want 12", got)
	}
	riAssertStable(t, riDiff(t, post, newConfig), []string{"renewal_status", "auto_renew", "auto_renew_period"})
}

// Period-only with a STALE prior state (state renewal_status=AutoRenewal) while the cloud's
// current renewal status is Normal (changed server-side). The period-only request must NOT
// send RenewalStatus, so the cloud keeps Normal — proving the period-only update does not
// accidentally overwrite a server-side renewal-status change. Before the
// renewalStatusFromConfig fix, d.GetOk("renewal_status") returned the stale state value
// AutoRenewal (exists=true for the Computed field), the request re-sent RenewalStatus=
// AutoRenewal, and the cloud's Normal was silently overwritten. This test fails without the
// fix (cloud becomes AutoRenewal) and passes with it (cloud stays Normal).
func TestUnitReservedInstanceUpdate_PeriodOnlyPreservesCloudRenewalStatus(t *testing.T) {
	// Cloud current state: Normal, Duration 1 / Year (changed server-side after the state
	// was last read).
	client, st, calls := riStatefulClient(t, "Normal", 1, "Year")
	newConfig := map[string]interface{}{
		"instance_type":     "ecs.g6.large",
		"auto_renew_period": 12,
	}
	// STALE prior state: still reflects AutoRenewal / period 36 (not refreshed).
	d := riUpdateData(t,
		map[string]string{
			"renewal_status":    "AutoRenewal",
			"auto_renew":        "true",
			"auto_renew_period": "36",
		},
		newConfig,
	)
	if err := resourceAliCloudReservedInstanceUpdate(d, client); err != nil {
		t.Fatalf("Update: %v", err)
	}
	c := riRequireCall(t, calls, "ModifyReservedInstanceAutoRenewAttribute")
	riFormEq(t, c.form, "Period", "1")
	riFormEq(t, c.form, "PeriodUnit", "Year")
	// Neither auto_renew nor renewal_status is explicit => RenewalStatus is absent, so the
	// cloud keeps its current Normal instead of being overwritten by the stale state.
	riFormAbsent(t, c.form, "RenewalStatus")
	st.mu.Lock()
	gotStatus, gotDur := st.renewStatus, st.duration
	st.mu.Unlock()
	if gotStatus != "Normal" || gotDur != 1 {
		t.Fatalf("stateful mock after period-only Modify = (%q,%d), want (Normal,1) — cloud renewal status must be preserved, not overwritten by stale state", gotStatus, gotDur)
	}
	post := d.State()
	if got := post.Attributes["renewal_status"]; got != "Normal" {
		t.Fatalf("post-Read renewal_status = %q, want Normal (cloud state preserved)", got)
	}
	if got := post.Attributes["auto_renew"]; got != "false" {
		t.Fatalf("post-Read auto_renew = %q, want false (derived from preserved Normal)", got)
	}
	if got := post.Attributes["auto_renew_period"]; got != "12" {
		t.Fatalf("post-Read auto_renew_period = %q, want 12", got)
	}
	riAssertStable(t, riDiff(t, post, newConfig), []string{"renewal_status", "auto_renew", "auto_renew_period"})
}

func TestUnitReservedInstanceUpdate_ConflictAutoRenewTrueRenewalStatusNormal(t *testing.T) {
	// Real contradiction: auto_renew=true wants AutoRenewal while the legacy
	// renewal_status=Normal wants Normal. auto_renew must win (sent RenewalStatus=
	// AutoRenewal). The closed loop proves that after the cloud reflects AutoRenewal, the
	// next plan with the SAME contradictory config is stable: renewal_status would change
	// AutoRenewal->Normal but, because auto_renew is explicit, renewalStatusDiffSuppress
	// suppresses it — no perpetual diff.
	client, st, calls := riStatefulClient(t, "Normal", 1, "Year")
	newConfig := map[string]interface{}{
		"instance_type":     "ecs.g6.large",
		"auto_renew":        true,
		"renewal_status":    "Normal",
		"auto_renew_period": 36,
	}
	d := riUpdateData(t,
		map[string]string{
			"renewal_status":    "Normal",
			"auto_renew":        "false",
			"auto_renew_period": "12",
		},
		newConfig,
	)
	if err := resourceAliCloudReservedInstanceUpdate(d, client); err != nil {
		t.Fatalf("Update: %v", err)
	}
	c := riRequireCall(t, calls, "ModifyReservedInstanceAutoRenewAttribute")
	riFormEq(t, c.form, "RenewalStatus", "AutoRenewal")
	riFormEq(t, c.form, "Period", "3")
	riFormEq(t, c.form, "PeriodUnit", "Year")
	st.mu.Lock()
	gotStatus, gotDur := st.renewStatus, st.duration
	st.mu.Unlock()
	if gotStatus != "AutoRenewal" || gotDur != 3 {
		t.Fatalf("stateful mock after Modify = (%q,%d), want (AutoRenewal,3)", gotStatus, gotDur)
	}
	post := d.State()
	if got := post.Attributes["renewal_status"]; got != "AutoRenewal" {
		t.Fatalf("post-Read renewal_status = %q, want AutoRenewal (auto_renew wins)", got)
	}
	if got := post.Attributes["auto_renew"]; got != "true" {
		t.Fatalf("post-Read auto_renew = %q, want true", got)
	}
	if got := post.Attributes["auto_renew_period"]; got != "36" {
		t.Fatalf("post-Read auto_renew_period = %q, want 36", got)
	}
	riAssertStable(t, riDiff(t, post, newConfig), []string{"renewal_status", "auto_renew", "auto_renew_period"})
}

// --- Read: RenewalStatus maps back to auto_renew / auto_renew_period --------

func TestUnitReservedInstanceRead_MapsRenewalStatusToAutoRenew(t *testing.T) {
	t.Run("auto_renewal", func(t *testing.T) {
		client, _, _ := riStatefulClient(t, "AutoRenewal", 3, "Year")
		d := riReadData(t, nil)
		if err := resourceAliCloudReservedInstanceRead(d, client); err != nil {
			t.Fatalf("Read: %v", err)
		}
		if got := d.Get("renewal_status").(string); got != "AutoRenewal" {
			t.Fatalf("renewal_status = %q, want AutoRenewal", got)
		}
		if got := d.Get("auto_renew").(bool); !got {
			t.Fatalf("auto_renew = %v, want true (derived from AutoRenewal)", got)
		}
		if got := d.Get("auto_renew_period").(int); got != 36 {
			t.Fatalf("auto_renew_period = %d, want 36", got)
		}
	})
	t.Run("normal", func(t *testing.T) {
		client, _, _ := riStatefulClient(t, "Normal", 1, "Year")
		d := riReadData(t, nil)
		if err := resourceAliCloudReservedInstanceRead(d, client); err != nil {
			t.Fatalf("Read: %v", err)
		}
		if got := d.Get("renewal_status").(string); got != "Normal" {
			t.Fatalf("renewal_status = %q, want Normal", got)
		}
		if got := d.Get("auto_renew").(bool); got {
			t.Fatalf("auto_renew = %v, want false (derived from Normal)", got)
		}
	})
}

// --- Read-computed does not mask a legacy update (plan-level) ----------------
//
// After Read, state carries computed auto_renew (from RenewalStatus). A subsequent plan
// where the user changes ONLY the legacy renewal_status (auto_renew not in config) must
// NOT suppress that change: the computed Read value must not steal control from the
// user's explicit legacy config. This complements the pure-function tests above.

func TestUnitReservedInstanceReadPlan_LegacyNotSuppressed(t *testing.T) {
	// Read returns AutoRenewal, so state = auto_renew=true (computed). The next plan sets
	// renewal_status=Normal (legacy) with auto_renew absent: the renewal_status change is a
	// real change and must NOT be suppressed.
	client, _, _ := riStatefulClient(t, "AutoRenewal", 3, "Year")
	d := riReadData(t, nil)
	if err := resourceAliCloudReservedInstanceRead(d, client); err != nil {
		t.Fatalf("Read: %v", err)
	}
	readState := d.State()
	diff := riDiff(t, readState, map[string]interface{}{
		"instance_type":  "ecs.g6.large",
		"renewal_status": "Normal",
	})
	riAssertStable(t, diff, nil)
	if !diffHasChange(diff, "renewal_status") {
		t.Fatal("renewal_status change must NOT be suppressed: auto_renew is not in config, so the legacy field controls the update")
	}
}

// --- Import: real Importer.State from ID-only empty state -> Read -> empty diff

func TestUnitReservedInstanceImport_NoOp(t *testing.T) {
	// Import: ID-only empty state (no pre-filled attrs, no RawConfig), exactly as Terraform
	// core delivers to ImportStatePassthrough. The real Importer.State is invoked, then a
	// real Read HTTP populates the resource from the stateful mock. A complete legal config
	// (instance_type Required; renewal fields are computed and absent) must produce a full
	// SDKDiff that is empty (nil or Empty): no spurious change, no Destroy, no RequiresNew —
	// not merely "no change on the 3 renewal fields".
	client, _, _ := riStatefulClient(t, "AutoRenewal", 3, "Year")
	r := resourceAliCloudReservedInstance()
	importD := r.Data(&terraform.InstanceState{ID: riMockID})
	states, err := r.Importer.State(importD, client)
	if err != nil {
		t.Fatalf("Importer.State: %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected 1 imported state, got %d", len(states))
	}
	if err := resourceAliCloudReservedInstanceRead(states[0], client); err != nil {
		t.Fatalf("Read after import: %v", err)
	}
	post := states[0].State()
	diff := riDiff(t, post, map[string]interface{}{"instance_type": "ecs.g6.large"})
	if diff != nil && !diff.Empty() {
		t.Fatalf("import no-op: full SDKDiff must be empty, got attrs: %+v", diff.Attributes)
	}
}
