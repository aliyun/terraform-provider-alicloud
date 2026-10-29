//nolint:all
package main

import (
	"fmt"
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

// writeImplStub writes a minimal impl fixture that calls d.Set for each root
// attribute name, so neverReadBackAttributes reports those roots as read back
// (not neverReadBack). The stub need not compile; neverReadBackAttributes only
// string-scans the source for `d.Set("<root>"`.
func writeImplStub(t *testing.T, roots ...string) string {
	t.Helper()
	implPath := filepath.Join(t.TempDir(), "resource_alicloud_coverage_gate_test.go")
	var b strings.Builder
	b.WriteString("// stub impl fixture for coverage gate tests\npackage coveragegate\n\n")
	b.WriteString("func setAttrs(d interface{ Set(string, interface{}) }) {\n")
	for _, r := range roots {
		fmt.Fprintf(&b, "\td.Set(%q, \"v\")\n", r)
	}
	b.WriteString("}\n")
	if err := os.WriteFile(implPath, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write impl fixture: %v", err)
	}
	return implPath
}

// TestUnitStripNumericIndexPartsForCoverage proves the index-stripping helper
// used by the coverage fix drops only purely-numeric path segments. Single
// digits (0/1), multi-digit (10/99) and nested indices are removed; field
// names that merely contain digits (zone_id, ipv6_address) are preserved.
func TestUnitStripNumericIndexPartsForCoverage(t *testing.T) {
	cases := map[string]string{
		"db_cluster_ip_array.0.white_list_type": "db_cluster_ip_array.white_list_type",
		"db_cluster_ip_array.1.white_list_type": "db_cluster_ip_array.white_list_type",
		"a.10.b.0.c":                            "a.b.c",
		"x.99.y":                                "x.y",
		"nic.ipv6.0.addr":                       "nic.ipv6.addr",
		// Field names containing digits are NOT indices and must be kept.
		"zone_id":            "zone_id",
		"nic.ipv6_address":   "nic.ipv6_address",
		"vpc_id.0.subnet_id": "vpc_id.subnet_id",
	}
	for in, want := range cases {
		if got := stripNumericIndexParts(in); got != want {
			t.Errorf("stripNumericIndexParts(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestUnitNormalizeIgnoreSetStripsIndices proves the normalized ignore set
// deduplicates indexed entries to their schema-path form while leaving
// non-indexed entries unchanged.
func TestUnitNormalizeIgnoreSetStripsIndices(t *testing.T) {
	got := normalizeIgnoreSet(setOf(
		"db_cluster_ip_array.0.white_list_type",
		"db_cluster_ip_array.1.white_list_type",
		"end_date",
	))
	if got.Cardinality() != 2 || !got.Contains("db_cluster_ip_array.white_list_type") || !got.Contains("end_date") {
		t.Fatalf("expected {db_cluster_ip_array.white_list_type, end_date}, got %v", got.ToSlice())
	}
}

// TestUnitIndexedImportIgnoreNotFalselyReportedUnmodified — AC1 end-to-end
// through the REAL production gate checkAttributeSet. Reproduces the reported
// PolarDB DBCluster shape: a leaf attribute (db_cluster_ip_array.white_list_type)
// is in the modification-coverage schema set but NOT exercised by any test
// config (testModifySet empty), and the only ImportStateVerifyIgnore entries
// are the real indexed runtime paths (.0./.1.). Before the fix the indexed
// ignores never matched the schema-path modification key, so the leaf stayed
// in the diff and the coverage rate was falsely 98.48% / FAIL. With index
// normalization the indexed ignores deduct the schema-path leaf and the gate
// returns true without requiring a schema-only string the SDK cannot match.
func TestUnitIndexedImportIgnoreNotFalselyReportedUnmodified(t *testing.T) {
	implPath := writeImplStub(t, "db_cluster_ip_array")
	schemaAllSet := setOf("db_cluster_ip_array", "db_cluster_ip_array.white_list_type")
	schemaMustSet := setOf("db_cluster_ip_array.white_list_type")
	schemaModifySet := setOf("db_cluster_ip_array.white_list_type")
	schemaForceNewSet := mapset.NewSet()
	schemaDeprecatedSet := mapset.NewSet()
	testMustSet := setOf("db_cluster_ip_array.white_list_type")
	testModifySet := mapset.NewSet()
	testIgnoreSet := setOf("db_cluster_ip_array.0.white_list_type", "db_cluster_ip_array.1.white_list_type")

	got := checkAttributeSet("polardb_dbcluster", "resource",
		schemaMustSet, testMustSet, schemaModifySet, testModifySet,
		schemaForceNewSet, schemaAllSet, testIgnoreSet, schemaDeprecatedSet, implPath)
	if !got {
		t.Fatalf("indexed ImportStateVerifyIgnore should deduct the schema-path leaf; checkAttributeSet returned FAIL (false unmodified report)")
	}
}

// TestUnitNonIgnoredUnmodifiedFieldStillFailsCoverage — AC3: a genuinely
// unmodified, NON-ignored field still fails the coverage gate. The indexed
// ignore deducts only its own schema-path leaf; a sibling leaf that has no
// test modification and is not ignored must remain in the diff.
func TestUnitNonIgnoredUnmodifiedFieldStillFailsCoverage(t *testing.T) {
	implPath := writeImplStub(t, "db_cluster_ip_array")
	schemaAllSet := setOf("db_cluster_ip_array", "db_cluster_ip_array.white_list_type", "db_cluster_ip_array.other_type")
	schemaMustSet := setOf("db_cluster_ip_array.white_list_type", "db_cluster_ip_array.other_type")
	schemaModifySet := setOf("db_cluster_ip_array.white_list_type", "db_cluster_ip_array.other_type")
	schemaForceNewSet := mapset.NewSet()
	schemaDeprecatedSet := mapset.NewSet()
	testMustSet := setOf("db_cluster_ip_array.white_list_type", "db_cluster_ip_array.other_type")
	testModifySet := mapset.NewSet()
	// Only white_list_type is ignored (indexed); other_type is not ignored and
	// not modified, so it must stay uncovered.
	testIgnoreSet := setOf("db_cluster_ip_array.0.white_list_type")

	got := checkAttributeSet("polardb_dbcluster", "resource",
		schemaMustSet, testMustSet, schemaModifySet, testModifySet,
		schemaForceNewSet, schemaAllSet, testIgnoreSet, schemaDeprecatedSet, implPath)
	if got {
		t.Fatalf("non-ignored unmodified field (db_cluster_ip_array.other_type) must fail the coverage gate; checkAttributeSet returned PASS")
	}
}

// TestUnitIllegalIndexedImportIgnoreStillFails — AC3: an ImportStateVerifyIgnore
// entry whose schema-path form is not part of the active schema is illegal and
// still fails the gate (isIgnoreLegal=false). Normalization does not legalize
// ignores that map to no real attribute.
func TestUnitIllegalIndexedImportIgnoreStillFails(t *testing.T) {
	implPath := writeImplStub(t, "real_attr")
	schemaAllSet := setOf("real_attr")
	schemaMustSet := setOf("real_attr")
	schemaModifySet := setOf("real_attr")
	schemaForceNewSet := mapset.NewSet()
	schemaDeprecatedSet := mapset.NewSet()
	testMustSet := setOf("real_attr")
	testModifySet := setOf("real_attr")
	// nonexistent_field.foo is not in the schema → illegal ignore.
	testIgnoreSet := setOf("nonexistent_field.0.foo")

	got := checkAttributeSet("coverage_gate_fixture", "resource",
		schemaMustSet, testMustSet, schemaModifySet, testModifySet,
		schemaForceNewSet, schemaAllSet, testIgnoreSet, schemaDeprecatedSet, implPath)
	if got {
		t.Fatalf("illegal ImportStateVerifyIgnore entry mapping to no schema attribute must fail; checkAttributeSet returned PASS")
	}
}

// TestUnitNestedMultiDigitIndexedIgnoreCovered — AC2: multi-digit and nested
// indexed ignore paths normalize to their schema-path form and deduct the
// matching modification keys, giving a consistent coverage conclusion.
func TestUnitNestedMultiDigitIndexedIgnoreCovered(t *testing.T) {
	implPath := writeImplStub(t, "a", "x")
	schemaAllSet := setOf("a", "a.b", "a.b.c", "x", "x.y")
	schemaMustSet := setOf("a.b.c", "x.y")
	schemaModifySet := setOf("a.b.c", "x.y")
	schemaForceNewSet := mapset.NewSet()
	schemaDeprecatedSet := mapset.NewSet()
	testMustSet := setOf("a.b.c", "x.y")
	testModifySet := mapset.NewSet()
	// Multi-digit (10, 99) and nested (a.10.b.0.c) indices normalize to a.b.c / x.y.
	testIgnoreSet := setOf("a.10.b.0.c", "x.99.y")

	got := checkAttributeSet("coverage_gate_fixture", "resource",
		schemaMustSet, testMustSet, schemaModifySet, testModifySet,
		schemaForceNewSet, schemaAllSet, testIgnoreSet, schemaDeprecatedSet, implPath)
	if !got {
		t.Fatalf("multi-digit/nested indexed ignores should normalize and deduct schema-path leaves; checkAttributeSet returned FAIL")
	}
}
