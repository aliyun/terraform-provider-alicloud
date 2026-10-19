package hcltestconfig

import (
	"reflect"
	"sort"
	"testing"
)

// sorted returns a sorted copy of s so attribute-set comparisons are
// order-independent.
func sorted(s []string) []string {
	out := append([]string(nil), s...)
	sort.Strings(out)
	return out
}

// fullCoverageTestSrc mirrors the SecurityGroupSnapshotPolicy split-layout
// Framework test: a fmt.Sprintf HCL helper setting both Optional+Computed
// attributes, and three TestAcc steps that each vary one of them (status then
// name). Both attributes are set in every step and modified across steps, so
// the resolved coverage is complete:
//
//	must   = {security_group_snapshot_policy_name, status}
//	modify = {security_group_snapshot_policy_name, status}
//
// This is the positive case: the real measurement passes when checked against
// the synthesized schema (must ⊇ schemaMustSet, modify ⊇ schemaModifySet).
// The backtick delimiters of the helper's raw string are spliced in via
// string concatenation because a Go raw string literal cannot contain a
// backtick. Config: is on its own line, matching the real test layout the
// configLineRegex anchors on.
const fullCoverageTestSrc = `package ecs_test

import (
	"fmt"
	"testing"
)

func TestAccAliCloudEcsSecurityGroupSnapshotPolicy_basic(t *testing.T) {
	name := "tf-testacc-sgsp"
	renamed := "tf-testacc-sgsp-r"
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig(name, "enable"),
			},
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig(name, "disable"),
			},
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig(renamed, "disable"),
			},
		},
	})
}

func testAccSecurityGroupSnapshotPolicyConfig(name, status string) string {
	return fmt.Sprintf(` + "`" + `
resource "alicloud_security_group_snapshot_policy" "default" {
  security_group_snapshot_policy_name = "%[1]s"
  status                               = "%[2]s"
}
` + "`" + `, name, status)
}
`

// TestResolve_FullCoveragePass is the positive case: the helper sets both
// schema must-set attributes and the three steps modify both of them, so
// Resolve reports a complete must/modify set. Against the synthesized schema
// this coverage passes instead of being silently skipped as "unmeasurable".
func TestUnitResolve_FullCoveragePass(t *testing.T) {
	measured, result, err := Resolve([]byte(fullCoverageTestSrc), "alicloud_security_group_snapshot_policy")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if !measured {
		t.Fatalf("Resolve: expected measured=true for a test with resolvable HCL helpers, got false (would silently skip instead of measuring)")
	}
	wantMust := []string{"security_group_snapshot_policy_name", "status"}
	wantModify := []string{"security_group_snapshot_policy_name", "status"}
	if got := sorted(result.Must); !reflect.DeepEqual(got, wantMust) {
		t.Fatalf("must = %v; want %v (full coverage requires both schema must-set attributes)", got, wantMust)
	}
	if got := sorted(result.Modify); !reflect.DeepEqual(got, wantModify) {
		t.Fatalf("modify = %v; want %v (both attributes change value across the three steps)", got, wantModify)
	}
}

// missingAttrTestSrc is the negative case: the helper omits
// security_group_snapshot_policy_name (only sets status). Resolve reports
// must={status}, which is missing a schema must-set attribute, so coverage
// fails instead of silently passing. The measurement surfaces the gap.
const missingAttrTestSrc = `package ecs_test

import (
	"fmt"
	"testing"
)

func TestAccAliCloudEcsSecurityGroupSnapshotPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig("enable"),
			},
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig("disable"),
			},
		},
	})
}

func testAccSecurityGroupSnapshotPolicyConfig(status string) string {
	return fmt.Sprintf(` + "`" + `
resource "alicloud_security_group_snapshot_policy" "default" {
  status = "%[1]s"
}
` + "`" + `, status)
}
`

// TestResolve_MissingAttributeFailsCoverage is the negative case: the helper
// only sets status, omitting security_group_snapshot_policy_name. Resolve
// reports must={status}, which is missing a schema must-set attribute — the
// measurement surfaces the gap instead of silently passing it.
func TestUnitResolve_MissingAttributeFailsCoverage(t *testing.T) {
	measured, result, err := Resolve([]byte(missingAttrTestSrc), "alicloud_security_group_snapshot_policy")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if !measured {
		t.Fatalf("Resolve: expected measured=true for a test with a resolvable (if incomplete) HCL helper")
	}
	wantMust := []string{"status"}
	if got := sorted(result.Must); !reflect.DeepEqual(got, wantMust) {
		t.Fatalf("must = %v; want %v (security_group_snapshot_policy_name must be missing)", got, wantMust)
	}
	// status does change across the two steps, so it is modified; the gap is
	// only in must-set, proving the measurement distinguishes the two layers.
	wantModify := []string{"status"}
	if got := sorted(result.Modify); !reflect.DeepEqual(got, wantModify) {
		t.Fatalf("modify = %v; want %v", got, wantModify)
	}
}

// noHelperTestSrc has a TestAcc function but no fmt.Sprintf HCL helper that
// declares the resource. Coverage is genuinely unmeasurable: Resolve reports
// measured=false (nil error), so the caller fail-closes instead of silently
// passing.
const noHelperTestSrc = `package ecs_test

import "testing"

func TestAccAliCloudEcsSecurityGroupSnapshotPolicy_basic(t *testing.T) {
	// No fmt.Sprintf HCL config helper declares this resource.
	t.Log("placeholder")
}
`

// TestResolve_NoHelperIsUnmeasurable proves that a split-layout test with no
// resolvable HCL helper is genuinely unmeasurable: Resolve returns
// measured=false with nil error. The caller must fail closed (return false)
// rather than silently passing — this is the F5 fix's central contract.
func TestUnitResolve_NoHelperIsUnmeasurable(t *testing.T) {
	measured, _, err := Resolve([]byte(noHelperTestSrc), "alicloud_security_group_snapshot_policy")
	if err != nil {
		t.Fatalf("Resolve: unexpected error for an unmeasurable test (should be measured=false, nil error): %v", err)
	}
	if measured {
		t.Fatalf("Resolve: expected measured=false (no HCL helper); the caller must fail closed, not silently pass")
	}
}

// unterminatedCallSrc has a Config call whose argument list has no matching
// close paren on the line. splitCallArgs fails closed with an error rather
// than returning partial args that would mislead coverage measurement.
const unterminatedCallSrc = `package ecs_test

import (
	"fmt"
	"testing"
)

func TestAccAliCloudEcsSecurityGroupSnapshotPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig(name, "enable"
			},
		},
	})
}

func testAccSecurityGroupSnapshotPolicyConfig(name, status string) string {
	return fmt.Sprintf(` + "`" + `
resource "alicloud_security_group_snapshot_policy" "default" {
  security_group_snapshot_policy_name = "%[1]s"
  status                               = "%[2]s"
}
` + "`" + `, name, status)
}
`

// TestResolve_UnterminatedCallFailsClosed proves a malformed Config call (no
// matching close paren) is reported as an error so the caller fail-closes
// rather than silently passing on partial argument data.
func TestUnitResolve_UnterminatedCallFailsClosed(t *testing.T) {
	_, _, err := Resolve([]byte(unterminatedCallSrc), "alicloud_security_group_snapshot_policy")
	if err == nil {
		t.Fatalf("Resolve: expected an error for an unterminated Config call (must fail closed, not silently pass)")
	}
}

// literalAttrTestSrc has a helper that sets status via a placeholder (so it is
// modified across steps) and name as a literal (so it is set but never
// modified). Resolve must distinguish the two: name is in must but not modify;
// status is in both. This proves the measurement does not blanket-mark
// everything as modified (which would game the modify coverage) nor drop
// literal attributes from must.
const literalAttrTestSrc = `package ecs_test

import (
	"fmt"
	"testing"
)

func TestAccAliCloudEcsSecurityGroupSnapshotPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig("enable"),
			},
			{
				Config: testAccSecurityGroupSnapshotPolicyConfig("disable"),
			},
		},
	})
}

func testAccSecurityGroupSnapshotPolicyConfig(status string) string {
	return fmt.Sprintf(` + "`" + `
resource "alicloud_security_group_snapshot_policy" "default" {
  security_group_snapshot_policy_name = "fixed"
  status                               = "%[1]s"
}
` + "`" + `, status)
}
`

// TestResolve_LiteralAttributeSetButNotModified proves a literal attribute is
// in must (it is set in the config) but not in modify (constant across steps),
// while a placeholder attribute that varies is in both. This guards against
// blanket-marking all attributes as modified (gaming the modify coverage) or
// dropping literals from must.
func TestUnitResolve_LiteralAttributeSetButNotModified(t *testing.T) {
	measured, result, err := Resolve([]byte(literalAttrTestSrc), "alicloud_security_group_snapshot_policy")
	if err != nil {
		t.Fatalf("Resolve: unexpected error: %v", err)
	}
	if !measured {
		t.Fatalf("Resolve: expected measured=true")
	}
	wantMust := []string{"security_group_snapshot_policy_name", "status"}
	if got := sorted(result.Must); !reflect.DeepEqual(got, wantMust) {
		t.Fatalf("must = %v; want %v (literal attr is still set)", got, wantMust)
	}
	wantModify := []string{"status"}
	if got := sorted(result.Modify); !reflect.DeepEqual(got, wantModify) {
		t.Fatalf("modify = %v; want %v (literal attr must not be modified; only the varying placeholder is)", got, wantModify)
	}
}
