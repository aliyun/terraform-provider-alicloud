//nolint:all
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/waigani/diffparser"
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

// scanCoverageFixture writes src to a temporary Go test file and runs the
// production step scanner scanTestCaseSteps against it, returning the
// captured per-function step state. Tests use it to assert step attribution
// (stepPlanOnly / stepStr) directly for layouts whose plan-only flag has no
// effect on the scored coverage sets (e.g. a helper-Config step).
func scanCoverageFixture(t *testing.T, src string) ResourceTest {
	t.Helper()
	filePath := filepath.Join(t.TempDir(), "resource_test.go")
	if err := os.WriteFile(filePath, []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	ignoreSet := mapset.NewSet()
	resourceTest, ok := scanTestCaseSteps(filePath, &ignoreSet)
	if !ok {
		t.Fatalf("scanner should read the fixture file")
	}
	return resourceTest
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

// PlanOnly modify-coverage semantics (Aone 88271295 residual defect): a
// PlanOnly step only refreshes the plan, it never applies, so a field whose
// value changes only inside a PlanOnly step must NOT enter testModifySet —
// while real apply steps still count, PlanOnly: false stays a normal step,
// and the plan-only values must not clobber the last applied value that
// later apply steps are compared against.

// onlyPlanNotCountedSrc — apply {a:1,b:1}; plan-only {a:9,b:1}; apply {a:1,b:1}.
// "a" changes only inside the plan-only step and is reverted afterwards, so
// nothing actually applied ever differs: the modify set must stay empty. The
// plan-only step also proves adjacent Config segmentation still parses (its
// non-empty Config is captured independently, no bracket/JSON error).
var onlyPlanNotCountedSrc = fixtureTab(`func TestAccAliCloudOnlyPlanNotCounted_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "9",
»»»»"b": "1",
»»»}),
»»»PlanOnly: true,
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerPlanOnlyStepExcludedFromModifyCoverage(t *testing.T) {
	ok, mustSet, modifySet := parseCoverageFixture(t, onlyPlanNotCountedSrc)
	if !ok {
		t.Fatalf("scanner should parse apply/plan-only/apply steps independently; got parse failure — steps likely merged")
	}
	// Plan-only steps still contribute attribute (must-set) coverage.
	for _, attr := range []string{"a", "b"} {
		if !mustSet.Contains(attr) {
			t.Fatalf("expected %q captured including from the plan-only step; got %v", attr, mustSet.ToSlice())
		}
		if modifySet.Contains(attr) {
			t.Fatalf("field %q changed only inside a PlanOnly step must not count as modification coverage; got modify set %v", attr, modifySet.ToSlice())
		}
	}
	if modifySet.Cardinality() != 0 {
		t.Fatalf("expected empty modify set (no actually-applied field changed); got %v", modifySet.ToSlice())
	}
}

// realApplyStillCountedSrc — apply {a:1,b:1}; plan-only {a:9,b:1}; apply {a:3,b:1}.
// "a" differs between the two real apply steps (1 -> 3) even though the
// plan-only step also touched it: the modification must still be counted.
var realApplyStillCountedSrc = fixtureTab(`func TestAccAliCloudRealApplyStillCounted_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "9",
»»»»"b": "1",
»»»}),
»»»PlanOnly: true,
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "3",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerRealApplyStillCountedAcrossPlanOnly(t *testing.T) {
	ok, _, modifySet := parseCoverageFixture(t, realApplyStillCountedSrc)
	if !ok {
		t.Fatalf("scanner should parse apply/plan-only/apply steps independently; got parse failure")
	}
	if !modifySet.Contains("a") {
		t.Fatalf("field 'a' changed between two real apply steps must still count as modification coverage; got %v", modifySet.ToSlice())
	}
	if modifySet.Contains("b") {
		t.Fatalf("field 'b' never changed in an applied step; got %v", modifySet.ToSlice())
	}
}

// planOnlyFalseSrc — apply {a:1}; apply {a:2} with an explicit PlanOnly: false.
// A false flag keeps the step a normal apply step, so the change still counts.
var planOnlyFalseSrc = fixtureTab(`func TestAccAliCloudPlanOnlyFalse_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "2",
»»»}),
»»»PlanOnly: false,
»»},
»}
}
`)

func TestUnitScannerPlanOnlyFalseStepStillCounted(t *testing.T) {
	ok, _, modifySet := parseCoverageFixture(t, planOnlyFalseSrc)
	if !ok {
		t.Fatalf("scanner should parse a PlanOnly:false step like a normal step; got parse failure")
	}
	if !modifySet.Contains("a") {
		t.Fatalf("explicit PlanOnly:false is a normal apply step; its field change must count as modification coverage; got %v", modifySet.ToSlice())
	}
}

// planOnlyAfterSentinelSrc — apply {a:1,b:1}; {a:2,b:1} with ExpectNonEmptyPlan
// closing the capture before the PlanOnly line (field order inside a TestStep
// literal is free). The step is still plan-only: "a" changed only there must
// not be counted.
var planOnlyAfterSentinelSrc = fixtureTab(`func TestAccAliCloudPlanOnlyAfterSentinel_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "2",
»»»»"b": "1",
»»»}),
»»»ExpectNonEmptyPlan: false,
»»»PlanOnly: true,
»»},
»}
}
`)

func TestUnitScannerPlanOnlyAfterSentinelNotCounted(t *testing.T) {
	ok, _, modifySet := parseCoverageFixture(t, planOnlyAfterSentinelSrc)
	if !ok {
		t.Fatalf("scanner should parse a step whose PlanOnly line follows an ExpectNonEmptyPlan sentinel; got parse failure")
	}
	if modifySet.Contains("a") {
		t.Fatalf("field 'a' changed only in a plan-only step (PlanOnly written after another sentinel) must not count; got %v", modifySet.ToSlice())
	}
	if modifySet.Cardinality() != 0 {
		t.Fatalf("expected empty modify set; got %v", modifySet.ToSlice())
	}
}

// planOnlyFirstStepSrc — plan-only {a:9,b:9}; apply {a:1,b:1}. The plan-only
// step never applies and must not establish a value baseline either, so the
// first real apply step is still the baseline: no field is a modification and
// the plan-only values must not fake "first occurrence in a later step".
var planOnlyFirstStepSrc = fixtureTab(`func TestAccAliCloudPlanOnlyFirstStep_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "9",
»»»»"b": "9",
»»»}),
»»»PlanOnly: true,
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerPlanOnlyFirstStepNotFalseBaseline(t *testing.T) {
	ok, mustSet, modifySet := parseCoverageFixture(t, planOnlyFirstStepSrc)
	if !ok {
		t.Fatalf("scanner should parse a leading plan-only step independently; got parse failure")
	}
	for _, attr := range []string{"a", "b"} {
		if !mustSet.Contains(attr) {
			t.Fatalf("expected %q captured including from the leading plan-only step; got %v", attr, mustSet.ToSlice())
		}
		if modifySet.Contains(attr) {
			t.Fatalf("field %q first set by the first real apply step after a leading plan-only step is a baseline, not a modification; got %v", attr, modifySet.ToSlice())
		}
	}
	if modifySet.Cardinality() != 0 {
		t.Fatalf("expected empty modify set (no apply step pair differs); got %v", modifySet.ToSlice())
	}
}

// configlessPlanOnlySrc — apply {a:1}; a step holding only ResourceName and
// PlanOnly (no Config, so nothing is captured for it); apply {a:3}. The
// Config-less step must not spill its plan-only flag onto the next captured
// step: "a" still changes between two real apply steps and must count.
var configlessPlanOnlySrc = fixtureTab(`func TestAccAliCloudConfiglessPlanOnly_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»ResourceName: "alicloud_example_foo",
»»»PlanOnly: true,
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "3",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerConfiglessPlanOnlyStepNotSpilledForward(t *testing.T) {
	ok, _, modifySet := parseCoverageFixture(t, configlessPlanOnlySrc)
	if !ok {
		t.Fatalf("scanner should parse a step without Config fields; got parse failure")
	}
	if !modifySet.Contains("a") {
		t.Fatalf("a Config-less plan-only step must not spill its flag onto the next captured apply step; 'a' changed 1 -> 3 between real applies; got %v", modifySet.ToSlice())
	}
	if modifySet.Cardinality() != 1 {
		t.Fatalf("expected exactly 'a' in the modify set; got %v", modifySet.ToSlice())
	}
}

// PlanOnly written BEFORE the Config line (field order inside a TestStep
// literal is free): the flag is deferred (planOnlyPending) until the step's
// Config capture really closes, then lands on that same step. These fixtures
// cover the deferral's main path and its two known layout boundaries.

// planOnlyBeforeConfigSrc — apply {a:1,b:1}; plan-only {a:9,b:9} with the
// PlanOnly line written before the Config line; apply {a:3,b:1}. The pending
// flag must land on the plan-only step: "a" still counts (1 -> 3 between real
// applies) while "b" (1 -> 9 -> 1, reverted only inside the plan-only step)
// must not. If the pending attribution were lost, the plan-only step would
// count as an apply and "b" would enter the modify set; if it landed on the
// wrong step, "a" would stop counting.
var planOnlyBeforeConfigSrc = fixtureTab(`func TestAccAliCloudPlanOnlyBeforeConfig_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»PlanOnly: true,
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "9",
»»»»"b": "9",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "3",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerPlanOnlyBeforeConfigDeferred(t *testing.T) {
	ok, mustSet, modifySet := parseCoverageFixture(t, planOnlyBeforeConfigSrc)
	if !ok {
		t.Fatalf("scanner should parse a step whose PlanOnly line precedes its Config line; got parse failure — steps likely merged")
	}
	if !mustSet.Contains("a") || !mustSet.Contains("b") {
		t.Fatalf("expected 'a' and 'b' captured including from the plan-only step; got %v", mustSet.ToSlice())
	}
	if !modifySet.Contains("a") {
		t.Fatalf("field 'a' changed between two real apply steps around a deferred plan-only step must still count; got %v", modifySet.ToSlice())
	}
	if modifySet.Contains("b") {
		t.Fatalf("field 'b' changed only inside a plan-only step (PlanOnly written before Config) must not count; got %v", modifySet.ToSlice())
	}
	if modifySet.Cardinality() != 1 {
		t.Fatalf("expected exactly 'a' in the modify set; got %v", modifySet.ToSlice())
	}
	// Structural pin of the deferral main path: the pending flag lands on the
	// step whose capture just closed, not on its neighbours.
	resourceTest := scanCoverageFixture(t, planOnlyBeforeConfigSrc)
	f := resourceTest.funcs["TestAccAliCloudPlanOnlyBeforeConfig_basic"]
	if !f.stepPlanOnly[1] || f.stepPlanOnly[0] || f.stepPlanOnly[2] {
		t.Fatalf("plan-only flag deferred from a PlanOnly-before-Config line must land on step 1 only; got %v", f.stepPlanOnly)
	}
}

// planOnlyBeforeConfigSameLineOpenSrc — same body, but each step literal
// opens on the `resource.TestStep{` line itself instead of a bare `{`
// (repo layouts do both). The step-literal open must be recognized so the
// PlanOnly-before-Config line defers instead of mislabeling the previous
// step (or being dropped when it is the first step).
var planOnlyBeforeConfigSameLineOpenSrc = fixtureTab(`func TestAccAliCloudPlanOnlySameLineOpen_basic(t *testing.T) {
»cases := []resource.TestStep{
»»resource.TestStep{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»resource.TestStep{
»»»PlanOnly: true,
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "9",
»»»»"b": "9",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»resource.TestStep{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "3",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerPlanOnlyBeforeConfigSameLineStepOpen(t *testing.T) {
	ok, mustSet, modifySet := parseCoverageFixture(t, planOnlyBeforeConfigSameLineOpenSrc)
	if !ok {
		t.Fatalf("scanner should parse steps whose literal opens on the resource.TestStep{ line; got parse failure")
	}
	if !mustSet.Contains("a") || !mustSet.Contains("b") {
		t.Fatalf("expected 'a' and 'b' captured; got %v", mustSet.ToSlice())
	}
	if !modifySet.Contains("a") {
		t.Fatalf("field 'a' changed between two real apply steps must still count; got %v", modifySet.ToSlice())
	}
	if modifySet.Contains("b") {
		t.Fatalf("field 'b' changed only inside the plan-only step of a same-line-opened literal must not count; got %v", modifySet.ToSlice())
	}
	resourceTest := scanCoverageFixture(t, planOnlyBeforeConfigSameLineOpenSrc)
	f := resourceTest.funcs["TestAccAliCloudPlanOnlySameLineOpen_basic"]
	if !f.stepPlanOnly[1] || f.stepPlanOnly[0] || f.stepPlanOnly[2] {
		t.Fatalf("plan-only flag must land on step 1 only when the literal opens on the resource.TestStep{ line; got %v", f.stepPlanOnly)
	}
}

// planOnlyBeforeConfigNestedBlockSrc — same body, with a nested block
// (ConfigPlanChecks) written between the PlanOnly line and the Config line.
// The nested block's closing `}` must not be mistaken for the step literal's
// own closing brace, which would drop the pending attribution.
var planOnlyBeforeConfigNestedBlockSrc = fixtureTab(`func TestAccAliCloudPlanOnlyNestedBlock_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»PlanOnly: true,
»»»ConfigPlanChecks: resource.ConfigPlanChecks{
»»»»PreApply: []resource.PlanCheck{
»»»»»resource.TestCheckResourceAttr("alicloud_example_foo", "name", "test"),
»»»»},
»»»},
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "9",
»»»»"b": "9",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "3",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»}
}
`)

func TestUnitScannerPlanOnlyBeforeConfigNestedBlockKeepsPending(t *testing.T) {
	ok, mustSet, modifySet := parseCoverageFixture(t, planOnlyBeforeConfigNestedBlockSrc)
	if !ok {
		t.Fatalf("scanner should parse a plan-only step with a nested block before its Config; got parse failure")
	}
	if !mustSet.Contains("a") || !mustSet.Contains("b") {
		t.Fatalf("expected 'a' and 'b' captured; got %v", mustSet.ToSlice())
	}
	if !modifySet.Contains("a") {
		t.Fatalf("field 'a' changed between two real apply steps must still count; got %v", modifySet.ToSlice())
	}
	if modifySet.Contains("b") {
		t.Fatalf("field 'b' changed only inside the plan-only step must not count even with a nested block before Config; got %v", modifySet.ToSlice())
	}
	resourceTest := scanCoverageFixture(t, planOnlyBeforeConfigNestedBlockSrc)
	f := resourceTest.funcs["TestAccAliCloudPlanOnlyNestedBlock_basic"]
	if !f.stepPlanOnly[1] || f.stepPlanOnly[0] || f.stepPlanOnly[2] {
		t.Fatalf("nested block braces between the PlanOnly line and the Config line must not drop the pending plan-only flag; got %v", f.stepPlanOnly)
	}
}

// checkFallbackPlanOnlySrc — a helper-Config step (no inline map) whose Check
// never opens a map either: the Check-block fallback is armed and waits. A
// `PlanOnly: true` field of that same step arrives while the fallback is
// still waiting and must mark that step instead of being swallowed as
// check-map text. The plan-only flag on a helper-Config step has no scoring
// effect (its Config carries no readable attributes), so this is asserted
// structurally against the captured step state.
var checkFallbackPlanOnlySrc = fixtureTab(`func TestAccAliCloudCheckFallbackPlanOnly_basic(t *testing.T) {
»cases := []resource.TestStep{
»»{
»»»Config: testAccConfig(map[string]interface{}{
»»»»"a": "1",
»»»»"b": "1",
»»»}),
»»»Check: resource.ComposeAggregateTestCheckFunc(),
»»},
»»{
»»»Config: testAccExampleConfig(name),
»»»Check: resource.ComposeAggregateTestCheckFunc(
»»»»resource.TestCheckResourceAttr("alicloud_example_foo", "name", "test"),
»»»),
»»»PlanOnly: true,
»»},
»}
}
`)

func TestUnitScannerCheckFallbackWaitingPlanOnly(t *testing.T) {
	ok, mustSet, modifySet := parseCoverageFixture(t, checkFallbackPlanOnlySrc)
	if !ok {
		t.Fatalf("scanner should parse a helper-Config step with a Check fallback; got parse failure")
	}
	if !mustSet.Contains("a") || !mustSet.Contains("b") {
		t.Fatalf("expected 'a' and 'b' captured from the apply step; got %v", mustSet.ToSlice())
	}
	if modifySet.Cardinality() != 0 {
		t.Fatalf("expected empty modify set (helper-Config step contributes no readable attributes); got %v", modifySet.ToSlice())
	}
	resourceTest := scanCoverageFixture(t, checkFallbackPlanOnlySrc)
	f := resourceTest.funcs["TestAccAliCloudCheckFallbackPlanOnly_basic"]
	if !f.stepPlanOnly[1] {
		t.Fatalf("a PlanOnly: true line arriving while the Check fallback is still waiting for its map must mark that step as plan-only; got %v", f.stepPlanOnly)
	}
	if f.stepPlanOnly[0] {
		t.Fatalf("the leading apply step must not be mislabeled plan-only; got %v", f.stepPlanOnly)
	}
}

func TestUnitHasCoverageRelevantTestChanges(t *testing.T) {
	const path = "alicloud/service/example/resource_example_test.go"
	testCases := []struct {
		name  string
		hunks string
		want  bool
	}{
		{
			name: "imports and non-ACC unit test provider adaptation",
			hunks: `@@ -11,7 +11,7 @@ import (
 "github.com/alibabacloud-go/tea/tea"
-	alicloud "github.com/aliyun/terraform-provider-alicloud/alicloud"
 tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
+	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
 )
@@ -26,7 +26,7 @@ func exampleConfig() string {
 func TestUnitExample(t *testing.T) {
-	p := alicloud.Provider().ResourcesMap
+	p := sdkv2.Provider().ResourcesMap
 dInit, _ := schema.InternalMap(p["alicloud_example"].Schema).Data(nil, nil)
 }
`,
			want: false,
		},
		{
			name: "imports and ACC provider factory adaptation",
			hunks: `@@ -5,6 +5,6 @@ import (
-	alicloud "github.com/aliyun/terraform-provider-alicloud/alicloud"
+	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
 )
@@ -18,7 +18,7 @@ func TestAccAliCloudExample(t *testing.T) {
 resource.Test(t, resource.TestCase{
-		ProviderFactories: alicloud.ProviderFactories,
+		ProviderFactories: sdkv2.ProviderFactories,
 Steps: []resource.TestStep{},
 })
`,
			want: false,
		},
		{
			name: "current acctest SDK provider factory adaptation",
			hunks: `@@ -18,7 +18,7 @@ func TestAccAliCloudExample(t *testing.T) {
 resource.Test(t, resource.TestCase{
-		ProviderFactories: alicloud.ProviderFactories,
+		ProviderFactories: acctest.SDKProviderFactories,
 Steps: []resource.TestStep{},
 })
`,
			want: false,
		},
		{
			name: "current acctest protocol factory adaptation",
			hunks: `@@ -18,7 +18,7 @@ func TestAccAliCloudExample(t *testing.T) {
 resource.Test(t, resource.TestCase{
-		ProviderFactories: alicloud.ProviderFactories,
+		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
 Steps: []resource.TestStep{},
 })
`,
			want: false,
		},
		{
			name: "current acctest SDK providers adaptation",
			hunks: `@@ -18,7 +18,7 @@ func TestAccAliCloudExample(t *testing.T) {
 resource.Test(t, resource.TestCase{
-		Providers: alicloud.Providers,
+		Providers: acctest.SDKProviders,
 Steps: []resource.TestStep{},
 })
`,
			want: false,
		},
		{
			name: "current acctest provider variable adaptation",
			hunks: `@@ -40,7 +40,7 @@ func exampleProviderResources() {
 func exampleProviderResources() {
-	p := alicloud.Provider().ResourcesMap
+	p := acctest.Provider.ResourcesMap
 }
`,
			want: false,
		},
		{
			name: "unknown provider qualifier is not mechanical adaptation",
			hunks: `@@ -18,7 +18,7 @@ func TestAccAliCloudExample(t *testing.T) {
 resource.Test(t, resource.TestCase{
-		ProviderFactories: alicloud.ProviderFactories,
+		ProviderFactories: other.ProviderFactories,
 Steps: []resource.TestStep{},
 })
`,
			want: true,
		},
		{
			name: "ACC declaration context does not trigger mechanical adaptation",
			hunks: `@@ -18,7 +18,7 @@
 func TestAccAliCloudExample(t *testing.T) {
 resource.Test(t, resource.TestCase{
-		ProviderFactories: alicloud.ProviderFactories,
+		ProviderFactories: sdkv2.ProviderFactories,
 Steps: []resource.TestStep{},
 })
`,
			want: false,
		},
		{
			name: "ordinary Go unit test provider adaptation",
			hunks: `@@ -40,7 +40,7 @@
 func TestExample(t *testing.T) {
-	p := alicloud.Provider().ResourcesMap
+	p := sdkv2.Provider().ResourcesMap
 }
`,
			want: false,
		},
		{
			name: "blank lines only",
			hunks: `@@ -40,2 +40,2 @@ func exampleConfig() string {
-
+
 `,
			want: false,
		},
		{
			name: "provider adaptation mixed with ACC config change",
			hunks: `@@ -18,9 +18,10 @@ func TestAccAliCloudExample(t *testing.T) {
 resource.Test(t, resource.TestCase{
-		ProviderFactories: alicloud.ProviderFactories,
+		ProviderFactories: sdkv2.ProviderFactories,
 Steps: []resource.TestStep{
 {
-				Config: testAccConfig(map[string]interface{}{}),
+				Config: testAccConfig(map[string]interface{}{
+					"sts_role_arn": "acs:ram::123456789012:role/example",
+				}),
 },
 },
 })
`,
			want: true,
		},
		{
			name: "stale unit test hunk header does not hide ACC config",
			hunks: `@@ -40,7 +40,7 @@ func TestUnitExample(t *testing.T) {
 }

 var testAccExampleConfig = map[string]interface{}{
-	"name": "old",
+	"name": "new",
 }
`,
			want: true,
		},
		{
			name: "unit test declaration does not hide later ACC config",
			hunks: `@@ -40,12 +40,12 @@ func TestUnitExample(t *testing.T) {
 func TestUnitExample(t *testing.T) {
-	p := alicloud.Provider().ResourcesMap
+	p := sdkv2.Provider().ResourcesMap
 }

 var testAccExampleConfig = map[string]interface{}{
-	"name": "old",
+	"name": "new",
 }
`,
			want: true,
		},
		{
			name: "non-mechanical provider factory change",
			hunks: `@@ -18,7 +18,7 @@ func TestAccAliCloudExample(t *testing.T) {
 resource.Test(t, resource.TestCase{
-		ProviderFactories: alicloud.ProviderFactories,
+		ProviderFactories: testAccCustomProviderFactories,
 Steps: []resource.TestStep{},
 })
`,
			want: true,
		},
		{
			name: "Go build directive change",
			hunks: `@@ -1,5 +1,5 @@
-//go:build unit
+//go:build !unit

 package example
`,
			want: true,
		},
		{
			name: "TestMain provider setup change",
			hunks: `@@ -20,7 +20,7 @@ func TestMain(m *testing.M) {
 func TestMain(m *testing.M) {
-	setupTestProvider(alicloud.Provider())
+	setupTestProvider(sdkv2.Provider())
 }
`,
			want: true,
		},
		{
			name: "provider name in comments",
			hunks: `@@ -118,7 +118,7 @@ func testSweepExample(region string) error {
 // factories for the cross-region tests. Each factory wraps a
-// fresh alicloud.Provider() (independent SDK Meta) with
+// fresh sdkv2.Provider() (independent SDK Meta) with
 // provider.ProtoV5ProviderServerFactory.
 func testAccExampleProviderFactories() {
`,
			want: false,
		},
		{
			name: "shared ACC config helper change",
			hunks: `@@ -40,7 +40,7 @@ func testAccExampleConfig(name string) string {
 return fmt.Sprintf(` + "`" + `
 resource "alicloud_example" "default" {
-	name = "old"
+	name = "new"
 }
 ` + "`" + `, name)
`,
			want: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := parseCoverageTestDiff(t, path, testCase.hunks)
			if got := hasCoverageRelevantTestChanges(file); got != testCase.want {
				t.Fatalf("hasCoverageRelevantTestChanges() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestUnitHasCoverageRelevantTestChangesWithoutHunks(t *testing.T) {
	file := &diffparser.DiffFile{Mode: diffparser.MODIFIED}
	if !hasCoverageRelevantTestChanges(file) {
		t.Fatal("modified test file without hunks must fail towards coverage checking")
	}
}

func parseCoverageTestDiff(t *testing.T, path, hunks string) *diffparser.DiffFile {
	t.Helper()
	raw := "diff --git a/" + path + " b/" + path + "\n--- a/" + path + "\n+++ b/" + path + "\n" + hunks
	diff, err := diffparser.Parse(raw)
	if err != nil {
		t.Fatalf("parse diff: %v", err)
	}
	if len(diff.Files) != 1 {
		t.Fatalf("parsed %d files, want 1", len(diff.Files))
	}
	return diff.Files[0]
}

// ===== indexed ImportStateVerifyIgnore normalization =====
//
// The CI gate compares modification coverage against schema leaf paths, while
// ImportStateVerifyIgnore entries legitimately carry the real indexed state
// keys (e.g. "db_cluster_ip_array.0.white_list_type") the SDK import
// comparison matches via strings.HasPrefix. normalizeIgnoreSet aligns the two
// spellings for the CI set comparisons only — the original strings stay in
// testIgnoreSet so runtime behavior does not change — and rejects entries
// whose index-like segments are not legal non-negative integers.

// TestUnitNormalizeIgnoreEntry — legal indexes (0, multi-digit, nested) are
// stripped so the entry aligns with the schema leaf path; field names that
// merely contain digits are never stripped; a signed (e.g. negative) index
// segment makes the whole entry invalid: it is returned unnormalized and
// reported instead of being silently aligned like a legal index.
func TestUnitNormalizeIgnoreEntry(t *testing.T) {
	cases := []struct {
		entry string
		want  string
		valid bool
	}{
		{"block.0.leaf", "block.leaf", true},                // index 0
		{"block.42.leaf", "block.leaf", true},               // multi-digit index
		{"block.007.leaf", "block.leaf", true},              // leading zeros are digits
		{"outer.12.inner.0.leaf", "outer.inner.leaf", true}, // nested indexes
		{"db_cluster_ip_array.0.white_list_type", "db_cluster_ip_array.white_list_type", true},
		{"block.leaf", "block.leaf", true},               // no index: untouched
		{"plain", "plain", true},                         // not dotted: untouched
		{"zones.zone_2.name", "zones.zone_2.name", true}, // field name with digit: not stripped
		{"ip6_entry.prefix", "ip6_entry.prefix", true},   // field name with digit: not stripped
		{"block.-1.leaf", "block.-1.leaf", false},        // negative index: rejected
		{"block.+1.leaf", "block.+1.leaf", false},        // signed index: rejected
		{"block.-1.0.leaf", "block.-1.0.leaf", false},    // negative among legal indexes: whole entry rejected
	}
	for _, c := range cases {
		got, valid := normalizeIgnoreEntry(c.entry)
		if got != c.want || valid != c.valid {
			t.Errorf("normalizeIgnoreEntry(%q) = (%q, %v), want (%q, %v)", c.entry, got, valid, c.want, c.valid)
		}
	}
}

// TestUnitNormalizeIgnoreSetKeepsOriginalAndStripped — the CI-side view keeps
// the original spelling (the SDK import comparison still matches it against
// real state keys) and additionally contributes the schema leaf path for the
// set comparisons. An invalid-index entry is returned separately and never
// normalized. Parent paths and sibling fields are not added: an entry for
// "block.0.leaf" contributes exactly "block.0.leaf" and "block.leaf", nothing
// else — the exemption is not widened.
func TestUnitNormalizeIgnoreSetKeepsOriginalAndStripped(t *testing.T) {
	normalized, invalid := normalizeIgnoreSet(setOf("block.0.leaf", "block.-1.leaf", "plain"))
	for _, want := range []string{"block.0.leaf", "block.leaf", "plain"} {
		if !normalized.Contains(want) {
			t.Errorf("normalized view should keep %q, got %v", want, normalized.ToSlice())
		}
	}
	// Parent paths and sibling fields are never added: an entry for
	// "block.0.leaf" contributes exactly "block.0.leaf" and "block.leaf",
	// nothing else — the exemption is not widened.
	for _, notWant := range []string{"block", "block.0", "block.sibling", "block.0.sibling"} {
		if normalized.Contains(notWant) {
			t.Errorf("normalized view must not add %q (exemption widening), got %v", notWant, normalized.ToSlice())
		}
	}
	// The illegal entry itself stays verbatim (raw spellings are preserved so
	// the gate fails on the flagged entry rather than silently dropping it),
	// and it is never de-indexed: it is reported through the invalid list,
	// where checkAttributeSet marks the ignore set illegal and fails.
	if !normalized.Contains("block.-1.leaf") {
		t.Errorf("normalized view keeps the raw illegal entry verbatim, missing %q, got %v", "block.-1.leaf", normalized.ToSlice())
	}
	if len(invalid) != 1 || invalid[0] != "block.-1.leaf" {
		t.Errorf("expected [block.-1.leaf] reported invalid, got %v", invalid)
	}
}

// writeCoverageImplFixture writes a stub implementation file the
// neverReadBackAttributes scanner reads. The stub need not compile: the
// scanner only string-scans it for d.Set("<root>" occurrences.
func writeCoverageImplFixture(t *testing.T, setsRoots ...string) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("// stub impl fixture for the coverage gate end-to-end test\npackage fixture\n\nfunc setAttrs(d interface{ Set(string, interface{}) }) {\n")
	for _, root := range setsRoots {
		b.WriteString("\td.Set(\"" + root + "\", \"v\")\n")
	}
	b.WriteString("}\n")
	implPath := filepath.Join(t.TempDir(), "resource_alicloud_coverage_fixture.go")
	if err := os.WriteFile(implPath, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write impl fixture: %v", err)
	}
	return implPath
}

// TestUnitIndexedImportIgnoreAlignsWithSchemaLeaf — END-TO-END regression
// through the REAL production gate checkAttributeSet, not a re-derived set
// formula. "block.leaf" is never modified in the test config; the original
// defect reported it as missing modification coverage whenever the only
// ImportStateVerifyIgnore entries were the real indexed state keys.
//
//	A: indexed ignore "block.0.leaf" offsets the schema leaf
//	   "block.leaf" → PASS (no schema-only string needed in the ignore set).
//	B: without any indexed ignore the leaf is still reported missing → FAIL.
//	C: a parent-path ignore "block" does NOT offset the leaf → FAIL.
//	D: a negative-index ignore "block.-1.leaf" is rejected, not aligned → FAIL.
//	E: a real missing modification on another (sibling) field stays FAIL even
//	   with the indexed ignore present → the exemption is not widened.
func TestUnitIndexedImportIgnoreAlignsWithSchemaLeaf(t *testing.T) {
	implPath := writeCoverageImplFixture(t, "block", "sibling", "other_leaf")

	// All attrs Optional; block is a container whose leaf is only exercised
	// through import. "other_leaf" is a real sibling field.
	schemaAllSet := setOf("block", "block.leaf", "sibling", "other_leaf")
	schemaMustSet := setOf("block", "block.leaf", "sibling", "other_leaf")
	schemaModifySet := setOf("block", "block.leaf", "sibling", "other_leaf")
	schemaForceNewSet := mapset.NewSet()
	schemaDeprecatedSet := mapset.NewSet()
	testMustSet := setOf("block", "block.leaf", "sibling", "other_leaf")

	check := func(ignoreSet mapset.Set, testModifySet mapset.Set) bool {
		return checkAttributeSet("coverage_fixture", "resource",
			schemaMustSet, testMustSet, schemaModifySet, testModifySet,
			schemaForceNewSet, schemaAllSet, ignoreSet, schemaDeprecatedSet, implPath)
	}

	leafUnmodified := setOf("block", "sibling", "other_leaf") // block.leaf never modified in config

	// A: the real indexed state key alone is enough — the schema-only string
	// the old gate demanded is no longer required.
	t.Run("indexed_ignore_offsets_schema_leaf", func(t *testing.T) {
		if !check(setOf("block.0.leaf", "block.1.leaf"), leafUnmodified) {
			t.Fatalf("indexed ImportStateVerifyIgnore should offset the schema leaf path; checkAttributeSet returned FAIL")
		}
	})

	// B: no ignore → the never-modified leaf is genuinely missing coverage.
	t.Run("no_ignore_still_fails", func(t *testing.T) {
		if check(mapset.NewSet(), leafUnmodified) {
			t.Fatalf("an unmodified, unignored leaf must stay a coverage failure; checkAttributeSet returned PASS")
		}
	})

	// C: ignoring only the parent path must not exempt the leaf.
	t.Run("parent_path_not_widened", func(t *testing.T) {
		if check(setOf("block"), leafUnmodified) {
			t.Fatalf("parent-path ignore must not offset the leaf's modification coverage; checkAttributeSet returned PASS")
		}
	})

	// D: an invalid (negative) index is rejected instead of being aligned
	// with the legal index form.
	t.Run("negative_index_rejected", func(t *testing.T) {
		if check(setOf("block.-1.leaf"), leafUnmodified) {
			t.Fatalf("negative-index ignore must be rejected, not normalized; checkAttributeSet returned PASS")
		}
	})

	// E: with the indexed ignore present, a genuinely missing modification on
	// a sibling field still fails.
	t.Run("sibling_missing_still_fails", func(t *testing.T) {
		if check(setOf("block.0.leaf"), setOf("block", "block.leaf")) {
			t.Fatalf("a real missing modification on a sibling field must stay blocked; checkAttributeSet returned PASS")
		}
	})
}

// TestUnitIndexedForceNewIgnoreMatchesSchemaPath — the explicit comparison
// between an indexed ForceNew ignore and its schema path form: both spellings
// must produce the identical gate outcome on every side of the ForceNew
// gate. When Read writes the container back (d.Set("block")), both the
// indexed "block.0.fn_attr" and the schema-path "block.fn_attr" are illegal
// in ImportStateVerifyIgnore; when the container is never read back
// (neverReadBack), both are legitimately exempted. The indexed spelling can
// neither escape the gate nor gain an exemption the plain spelling lacks.
func TestUnitIndexedForceNewIgnoreMatchesSchemaPath(t *testing.T) {
	schemaAllSet := setOf("block", "block.fn_attr")
	schemaMustSet := setOf("block", "block.fn_attr")
	schemaModifySet := setOf("block") // fn_attr is ForceNew: not a modification target
	schemaForceNewSet := setOf("block.fn_attr")
	schemaDeprecatedSet := mapset.NewSet()
	testMustSet := setOf("block", "block.fn_attr")
	testModifySet := setOf("block")

	run := func(implRoots []string, ignoreSet mapset.Set) bool {
		implPath := writeCoverageImplFixture(t, implRoots...)
		return checkAttributeSet("coverage_fixture", "resource",
			schemaMustSet, testMustSet, schemaModifySet, testModifySet,
			schemaForceNewSet, schemaAllSet, ignoreSet, schemaDeprecatedSet, implPath)
	}

	// Read writes "block" back: a ForceNew attr inside it must not be ignored.
	t.Run("read_back_both_spellings_illegal", func(t *testing.T) {
		if run([]string{"block"}, setOf("block.0.fn_attr")) {
			t.Fatalf("indexed ForceNew ignore (block.0.fn_attr) must be flagged like its schema path; checkAttributeSet returned PASS")
		}
		if run([]string{"block"}, setOf("block.fn_attr")) {
			t.Fatalf("schema-path ForceNew ignore (block.fn_attr) must be flagged; checkAttributeSet returned PASS")
		}
	})

	// Container never read back: both spellings are legitimately exempt.
	t.Run("never_read_back_both_spellings_exempt", func(t *testing.T) {
		if !run(nil, setOf("block.0.fn_attr")) {
			t.Fatalf("indexed never-read-back ForceNew ignore must be exempt like its schema path; checkAttributeSet returned FAIL")
		}
		if !run(nil, setOf("block.fn_attr")) {
			t.Fatalf("schema-path never-read-back ForceNew ignore must stay exempt; checkAttributeSet returned FAIL")
		}
	})
}

// TestUnitIndexedApiNonRoundTripIgnoreExempted — END-TO-END regression for
// the alicloud_instance stock ACC shape. The stock test
// TestAccAliCloudECSInstanceCpuOptionsAccelerators ignores the indexed state
// key "cpu_options.0.accelerators" because DescribeInstances never echoes
// CpuOptions.Accelerators while Read preserves the configured members from
// prior state via d.Set("cpu_options", ...). Since the indexed-ignore
// normalization aligns that key with the schema leaf path
// "cpu_options.accelerators" (a ForceNew leaf under a read-back container),
// only the apiNonRoundTrip registry entry clears the ForceNew gate for it.
//
//	A: the indexed ignore of the registered non-round-trip ForceNew leaf
//	   → PASS (the registry exemption applies after normalization).
//	B: an unregistered ForceNew sibling under the same read-back container
//	   → FAIL (the exemption is not widened to the whole container).
func TestUnitIndexedApiNonRoundTripIgnoreExempted(t *testing.T) {
	// Read writes the cpu_options container back, so neverReadBackAttributes
	// exempts nothing here; the registry is the only exemption path.
	implPath := writeCoverageImplFixture(t, "cpu_options")

	schemaAllSet := setOf("cpu_options", "cpu_options.core_count", "cpu_options.accelerators", "cpu_options.topology_type")
	schemaMustSet := setOf("cpu_options", "cpu_options.core_count", "cpu_options.accelerators", "cpu_options.topology_type")
	schemaModifySet := setOf("cpu_options", "cpu_options.core_count", "cpu_options.accelerators", "cpu_options.topology_type")
	schemaForceNewSet := setOf("cpu_options.accelerators", "cpu_options.topology_type")
	schemaDeprecatedSet := mapset.NewSet()
	testMustSet := setOf("cpu_options", "cpu_options.core_count", "cpu_options.accelerators", "cpu_options.topology_type")
	testModifySet := setOf("cpu_options", "cpu_options.core_count", "cpu_options.accelerators", "cpu_options.topology_type")

	check := func(ignoreSet mapset.Set) bool {
		return checkAttributeSet("instance", "resource",
			schemaMustSet, testMustSet, schemaModifySet, testModifySet,
			schemaForceNewSet, schemaAllSet, ignoreSet, schemaDeprecatedSet, implPath)
	}

	// A: the stock indexed ignore (as written in
	// resource_alicloud_instance_test.go) normalizes onto the registered leaf.
	t.Run("indexed_registered_ignore_legal", func(t *testing.T) {
		if !check(setOf("cpu_options.0.accelerators")) {
			t.Fatalf("indexed ignore of the registered non-round-trip cpu_options.accelerators must clear the ForceNew gate; checkAttributeSet returned FAIL")
		}
	})

	// B: an unregistered ForceNew sibling under the same container stays illegal.
	t.Run("unregistered_sibling_still_illegal", func(t *testing.T) {
		if check(setOf("cpu_options.0.topology_type")) {
			t.Fatalf("an unregistered ForceNew ignore under the same read-back container must stay illegal; checkAttributeSet returned PASS")
		}
	})
}
