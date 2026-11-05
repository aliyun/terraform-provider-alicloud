//nolint:all
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	mapset "github.com/deckarep/golang-set"
)

// apiNonRoundTripAttributes is the narrow exemption added to recognize
// ImportStateVerifyIgnore attributes whose values the API does NOT round-trip
// (Read writes them back via d.Set, but the API returns a value != configured).
// These tests prove the exemption is registration-gated and does NOT blanket
// every attribute that uses d.Set.
//
// Functions are named TestUnit* so the ci-check-quick changed-unit-tests stage
// (scripts/ci/changed-unit-tests.go) collects and runs them; only functions
// whose names have the TestUnit prefix are selected.

func setOf(items ...string) mapset.Set {
	s := mapset.NewSet()
	for _, i := range items {
		s.Add(i)
	}
	return s
}

// TestUnitApiNonRoundTripRegisteredAttrExempted — POSITIVE fixture: the one
// registered (resource, attr) pair (aligreen_oss_stock_task.end_date), whose
// configured value the API does not round-trip (Read writes EndDate back via
// d.Set, but Describe returns the current time, not the configured end_date),
// is exempted when it is in the resource's ImportStateVerifyIgnore.
func TestUnitApiNonRoundTripRegisteredAttrExempted(t *testing.T) {
	ignoreSet := setOf("end_date")
	got := apiNonRoundTripAttributes("aligreen_oss_stock_task", ignoreSet)
	if got.Cardinality() != 1 || !got.Contains("end_date") {
		t.Fatalf("expected end_date exempted, got %v", got.ToSlice())
	}
}

// TestUnitApiNonRoundTripAttrNotInIgnoreSet — NEGATIVE fixture 1: the registered
// attr is exempted only if it is actually in ImportStateVerifyIgnore. If the
// resource dropped it from the ignore set (an ACC probe that removes end_date
// from ImportStateVerifyIgnore fails ImportStateVerify), the exemption returns
// nothing — proving the exemption does not silently add attrs to the ignore set.
func TestUnitApiNonRoundTripAttrNotInIgnoreSet(t *testing.T) {
	ignoreSet := setOf("other_attr") // end_date not ignored here
	got := apiNonRoundTripAttributes("aligreen_oss_stock_task", ignoreSet)
	if got.Cardinality() != 0 {
		t.Fatalf("expected empty (end_date not in ignoreSet), got %v", got.ToSlice())
	}
}

// TestUnitApiNonRoundTripUnregisteredResource — NEGATIVE fixture 2: an
// unregistered resource fails closed (empty exemption), so the original FAIL
// behavior is preserved. This is the core safeguard against blanket exemption.
func TestUnitApiNonRoundTripUnregisteredResource(t *testing.T) {
	ignoreSet := setOf("end_date", "some_force_new")
	got := apiNonRoundTripAttributes("ecs_instance", ignoreSet) // not in registry
	if got.Cardinality() != 0 {
		t.Fatalf("expected empty for unregistered resource, got %v", got.ToSlice())
	}
}

// TestUnitApiNonRoundTripOnlyRegisteredAttrExempted — NEGATIVE fixture 3: when a
// registered resource's ignore set also contains an unregistered attr, only
// the registered attr is exempted — the unregistered one is NOT. This proves
// "not every d.Set attribute is exempted": a sibling ForceNew attr that lacks
// proven API non-round-trip evidence stays a violation.
func TestUnitApiNonRoundTripOnlyRegisteredAttrExempted(t *testing.T) {
	ignoreSet := setOf("end_date", "another_force_new")
	got := apiNonRoundTripAttributes("aligreen_oss_stock_task", ignoreSet)
	if got.Cardinality() != 1 || !got.Contains("end_date") || got.Contains("another_force_new") {
		t.Fatalf("expected only end_date exempted, got %v", got.ToSlice())
	}
}

// TestUnitApiNonRoundTripForceNewButIgnoreEndToEnd — END-TO-END regression
// through the REAL production gate checkAttributeSet, not a re-derived
// Intersect/Difference formula inside the test. A temporary impl fixture
// writes d.Set("end_date", ...) so neverReadBackAttributes reports end_date
// as read back (no neverReadBack false positive); the apiNonRoundTrip
// exemption is therefore the ONLY thing that can clear a registered
// ForceNew+ignored attr. If the production gate stopped applying the
// exemption, scenario A would FAIL — a genuine gate regression, not a
// tautology. This is why the test asserts on checkAttributeSet's return
// value rather than recomputing the set arithmetic itself.
//
//	A: registered aligreen_oss_stock_task.end_date (ForceNew, in
//	   ImportStateVerifyIgnore, Read writes it back via d.Set, API does not
//	   round-trip the configured value) → checkAttributeSet returns true.
//	B: an unregistered ForceNew attr in ImportStateVerifyIgnore → returns
//	   false (still illegal, fail-closed).
func TestUnitApiNonRoundTripForceNewButIgnoreEndToEnd(t *testing.T) {
	// Temporary impl fixture: Read writes BOTH attrs back via d.Set, so
	// neither is neverReadBack. The only path that can exempt end_date is the
	// apiNonRoundTrip registry; another_force_new has no entry, so it stays a
	// violation. neverReadBackAttributes only string-scans this source for
	// `d.Set("end_date"` / `d.Set("another_force_new"`, so the stub need not
	// compile.
	implPath := filepath.Join(t.TempDir(), "resource_alicloud_aligreen_oss_stock_task.go")
	implSrc := []byte(`// stub impl fixture for the coverage gate end-to-end test
package aligreen

func setAttrs(d interface{ Set(string, interface{}) }) {
	d.Set("end_date", "2027-10-27 00:00:00 +0800")
	d.Set("another_force_new", "v")
}
`)
	if err := os.WriteFile(implPath, implSrc, 0o644); err != nil {
		t.Fatalf("write impl fixture: %v", err)
	}

	// Both attrs are Optional+ForceNew. schema/test must and modify sets are
	// fully covered, so isFullCover and isAllModified are true; the return
	// value reflects isIgnoreLegal — i.e. the exemption's effect — alone.
	schemaAllSet := setOf("end_date", "another_force_new")
	schemaMustSet := setOf("end_date", "another_force_new")
	schemaModifySet := setOf("end_date", "another_force_new")
	schemaForceNewSet := setOf("end_date", "another_force_new")
	schemaDeprecatedSet := mapset.NewSet()
	testMustSet := setOf("end_date", "another_force_new")
	testModifySet := setOf("end_date", "another_force_new")

	check := func(ignoreSet mapset.Set) bool {
		return checkAttributeSet("aligreen_oss_stock_task", "resource",
			schemaMustSet, testMustSet, schemaModifySet, testModifySet,
			schemaForceNewSet, schemaAllSet, ignoreSet, schemaDeprecatedSet, implPath)
	}

	// Scenario A: the registered end_date is ForceNew and in
	// ImportStateVerifyIgnore; Read writes it back (not neverReadBack) and it
	// is in the apiNonRoundTrip registry → forceNewButIgnore is empty →
	// isIgnoreLegal stays true → checkAttributeSet returns true (PASS).
	t.Run("registered_end_date_legal", func(t *testing.T) {
		if !check(setOf("end_date")) {
			t.Fatalf("registered end_date (proven non-round-trip, read back via d.Set) should be legal; checkAttributeSet returned FAIL")
		}
	})

	// Scenario B: an unregistered ForceNew attr in ImportStateVerifyIgnore.
	// apiNonRoundTripAttributes returns empty for it (not registered) and
	// neverReadBackAttributes does not exempt it (d.Set present), so it stays
	// in forceNewButIgnore → isIgnoreLegal=false → checkAttributeSet returns
	// false (FAIL).
	t.Run("unregistered_force_new_still_illegal", func(t *testing.T) {
		if check(setOf("another_force_new")) {
			t.Fatalf("unregistered another_force_new in ImportStateVerifyIgnore should remain illegal; checkAttributeSet returned PASS")
		}
	})
}

// fixtureTab expands the » marker into a single tab so fixture source can be
// written readably inside a raw string literal while still matching the
// scanner's tab-anchored regexes (configRegex / stepEndSentinelRegex key on
// \t, not \s). The marker is replaced before the fixture is written to the
// temp file the scanner reads, so » never reaches the scanner.
func fixtureTab(s string) string {
	return strings.ReplaceAll(s, "»", "\t")
}

// parseCoverageFixture writes src to a temporary Go test file and runs the
// production scanner getTestCaseAttr against it, returning whether parsing
// succeeded plus the captured must/modify attribute sets. Tests assert on both
// Config-step segmentation (success/failure) and attribute coverage.
//
// Fixture functions inside src are named TestAccAliCloud* (NOT TestUnit*)
// because the scanner skips TestUnit* functions (unitFuncRegex); the actual
// test functions here keep the TestUnit* prefix so ci-check-quick collects them.
func parseCoverageFixture(t *testing.T, src string) (ok bool, mustSet, modifySet mapset.Set) {
	t.Helper()
	filePath := filepath.Join(t.TempDir(), "resource_test.go")
	if err := os.WriteFile(filePath, []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	mustSet = mapset.NewSet()
	modifySet = mapset.NewSet()
	ignoreSet := mapset.NewSet()
	ok = getTestCaseAttr(filePath, "alicloud_test_resource", &mustSet, &modifySet, &ignoreSet)
	return ok, mustSet, modifySet
}

// AC1: a PlanOnly step (with an empty-map Config) must end Config capture so
// the next (nested) Config step is segmented independently. Before the fix
// PlanOnly was not a sentinel: the PlanOnly line, the step boundary `},`/`{`
// and the next Config line were all appended to the same configStr, producing
// an unbalanced-bracket blob that bracketMatch rejected (false). The nested
// config also proves list/map type declarations (matchMap) still parse.
var planOnlyEmptyMapThenNestedSrc = fixtureTab(`func TestAccAliCloudPlanOnly_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{}),
»»»PlanOnly: true,
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"name": "test",
»»»»"items": []map[string]interface{}{
»»»»»{
»»»»»»"key": "value",
»»»»»},
»»»»},
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerPlanOnlyEndsNestedConfigCapture(t *testing.T) {
	ok, mustSet, _ := parseCoverageFixture(t, planOnlyEmptyMapThenNestedSrc)
	if !ok {
		t.Fatalf("scanner should parse PlanOnly step + adjacent nested Config independently (no bracket/JSON error); got parse failure — steps likely merged")
	}
	for _, attr := range []string{"name", "items", "items.key"} {
		if !mustSet.Contains(attr) {
			t.Fatalf("expected %q captured from the nested Config step; got %v", attr, mustSet.ToSlice())
		}
	}
}

// AC1 (flat variant): same segmentation with a simple flat Config so the
// regression is not pinned to a single config shape.
var planOnlyEmptyMapThenFlatSrc = fixtureTab(`func TestAccAliCloudPlanOnly_flat(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{}),
»»»PlanOnly: true,
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"name": "test",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerPlanOnlyEndsFlatConfigCapture(t *testing.T) {
	ok, mustSet, _ := parseCoverageFixture(t, planOnlyEmptyMapThenFlatSrc)
	if !ok {
		t.Fatalf("scanner should parse PlanOnly step + adjacent flat Config independently; got parse failure — steps likely merged")
	}
	if !mustSet.Contains("name") {
		t.Fatalf("expected 'name' captured from the flat Config step; got %v", mustSet.ToSlice())
	}
}

// AC2: consecutive PlanOnly-only steps must not merge; the final Config step
// (with Check) still segments and captures its attributes.
var continuousPlanOnlyThenConfigSrc = fixtureTab(`func TestAccAliCloudContinuousPlanOnly_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{}),
»»»PlanOnly: true,
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{}),
»»»PlanOnly: true,
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"name": "final",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerContinuousPlanOnlyNotMerged(t *testing.T) {
	ok, mustSet, _ := parseCoverageFixture(t, continuousPlanOnlyThenConfigSrc)
	if !ok {
		t.Fatalf("scanner should parse consecutive PlanOnly steps without merging; got parse failure")
	}
	if !mustSet.Contains("name") {
		t.Fatalf("expected 'name' captured from the final Config step; got %v", mustSet.ToSlice())
	}
}

// AC3: the PlanOnly addition must not disturb existing sentinels — Check,
// ExpectError and ExpectNonEmptyPlan steps still segment independently.
var sentinelStepsSrc = fixtureTab(`func TestAccAliCloudSentinels_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"name": "a",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"name": "b",
»»»}),
»»»ExpectError: regexp.MustCompile("error"),
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"name": "c",
»»»}),
»»»ExpectNonEmptyPlan: true,
»»},
»}
}
`)

func TestUnitScannerSentinelsStillIndependent(t *testing.T) {
	ok, mustSet, _ := parseCoverageFixture(t, sentinelStepsSrc)
	if !ok {
		t.Fatalf("scanner should parse Check/ExpectError/ExpectNonEmptyPlan steps independently; got parse failure")
	}
	if !mustSet.Contains("name") {
		t.Fatalf("expected 'name' captured across sentinel steps; got %v", mustSet.ToSlice())
	}
}

// AC4: a genuinely illegal config (unbalanced brackets — an extra closing
// brace) must still be blocked after the PlanOnly fix; the scanner must not
// swallow the error or skip all checks when it cannot parse.
var illegalBracketConfigSrc = fixtureTab(`func TestAccAliCloudIllegalBracket_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"name": "test",
»»»}}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerIllegalConfigStillBlocked(t *testing.T) {
	ok, _, _ := parseCoverageFixture(t, illegalBracketConfigSrc)
	if ok {
		t.Fatalf("scanner must still block genuinely illegal (unbalanced bracket) configs after the PlanOnly fix; got parse success")
	}
}
