# TypeList Order Coverage

`TypeList Order Coverage` requires resources touched by a PR to test that reordering
a configurable multi-member `TypeList` produces a diff and then converges after
applying that order. It reads registered Provider schemas and inspects real
acceptance-test steps. It does **not** run cloud operations: the acceptance runner
must execute the test to verify apply, refresh and post-apply convergence.
The implementation retains its historical directory name, `scripts/collection-order`.

## Scope

The workflow uses `scripts/ci/resolve-change-range.sh` and the full PR merge-base
to head diff, including earlier commits. It checks changed
`alicloud/resource_alicloud_*.go` implementations and `_test.go` files, and new or
changed `provider.go` resource registrations. Constructor identities resolve
aliases and renamed files. Registrations sharing one constructor share its ACC
coverage. Deleted resources are skipped only if no longer registered; missing
registrations, mappings and required test files fail closed.

Shared service/helper-only changes are outside this direct resource-file scope.
Review them manually to identify affected resources, and include and run the affected
resource's acceptance test when changing list handling in a shared helper. A PR
without selected resources exits successfully.

Candidates are configurable (`Optional` or `Required`) multi-member `TypeList`
fields, including nested paths. Removed fields, computed-only fields/subtrees,
`MaxItems: 1` containers themselves, and `TypeSet` itself are excluded. Configurable
`TypeList` children inside a `TypeSet` or single-item container are still checked.

For a list nested inside a multi-member set, keep the outer set members in the
same literal order in all three configurations and reorder only the inner list.
The checker compares outer containers by index; it does not match set members by
identity. Review the correspondence between those members manually.

The SDK defines [TypeList as ordered and TypeSet as unordered](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-types#typelist).
Ordinary `TypeSet` attributes do not need a repeated per-resource reorder test.
Custom set hashing or conversions can have targeted tests when needed. Comments
cannot exempt a `TypeList` from the required acceptance steps.

## Required acceptance steps

Use at least two distinct, valid members and three adjacent steps in an actual
`TestAcc...(*testing.T)` function in `resource_alicloud_<name>_test.go`:

1. Apply configuration A successfully.
2. Change only the target list's order to B, with `PlanOnly: true` and
   `ExpectNonEmptyPlan: true`.
3. Apply the same complete configuration B normally, with
   `ExpectNonEmptyPlan: false` or omitted. The SDK then requires the post-apply
   plan to be empty, including its refresh check.

```go
func TestAccExample_memberOrder(t *testing.T) {
    resourceId := "alicloud_example.default"
    name := fmt.Sprintf("tf-testacc-order-%d", acctest.RandInt())
    testAccConfig := resourceTestAccConfigFunc(resourceId, name, exampleDependence)
    resource.Test(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        Providers: testAccProviders,
        CheckDestroy: checkDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccConfig(map[string]interface{}{
                    // Include the required creation attributes here.
                    "member_ids": []string{
                        "${alicloud_example_member.first.id}",
                        "${alicloud_example_member.second.id}",
                    },
                }),
            },
            {
                Config: testAccConfig(map[string]interface{}{
                    "member_ids": []string{
                        "${alicloud_example_member.second.id}",
                        "${alicloud_example_member.first.id}",
                    },
                }),
                PlanOnly: true,
                ExpectNonEmptyPlan: true,
            },
            {
                Config: testAccConfig(map[string]interface{}{
                    "member_ids": []string{
                        "${alicloud_example_member.second.id}",
                        "${alicloud_example_member.first.id}",
                    },
                }),
                ExpectNonEmptyPlan: false,
            },
        },
    })
}
```

A nonempty reorder plan alone is insufficient: list schema semantics can create
a diff even when the remote API discards order. Applying B and requiring an empty
post-apply plan catches that failure to converge. Do not set
`ExpectNonEmptyPlan: true` on either apply step or apply B before checking its
plan. Each candidate path needs its own isolated sequence; nested lists of maps
and slices are supported. Review still establishes the API's ordering contract.

The checker parses Go AST rather than searching for names or comments. It follows
the config builder's cumulative top-level map updates, including `REMOVEKEY`,
`CLEARMAP` and `CLEARLIST`, and compares complete snapshots. Membership changes,
other configuration changes, a missing final apply, and `lifecycle.ignore_changes`
cannot count. Skipped, error-expected, import, destroy, refresh and plan-only steps
cannot stand in for a normal apply. Do not set `Taint` on any of the three steps,
including the reorder plan; replacement caused by schema `ForceNew` is allowed.

Supported proof is deliberately narrow: inline `resource.Test`/`ParallelTest`,
literal `TestCase.Steps`, a directly initialized `resourceTestAccConfigFunc`, and
literal map/slice configuration containers. Scalar expressions are compared
syntactically; the acceptance run verifies their runtime values. Dynamic HCL,
helper-generated steps, build-constrained test files, conditional runners, builder
aliases and builder calls outside Config do not count as coverage. Use the form
above when a test cannot be recognized. This is a coverage check, not a Go
interpreter or a substitute for reviewing helper/dependency behavior.

## Migrating existing tests

Keep existing `fmt.Sprintf`-based tests and add a separate three-step `TestAcc`
function in the same resource test file, using the example above. Move the target
resource's attributes into its config maps; extract only dependency HCL into
`func exampleDependence(name string) string`. The dependency function may still
use `fmt.Sprintf`, but must not declare the target resource a second time:

```go
func exampleDependence(name string) string {
    return fmt.Sprintf(`
resource "alicloud_example_member" "first" {
  name = "%[1]s-first"
}
resource "alicloud_example_member" "second" {
  name = "%[1]s-second"
}
`, name)
}
```

The example resource types, fields and `checkDestroy` are placeholders: retain
the real resource's dependencies, required attributes and destroy check. For
example, the legacy `network_acl_entries` tests can keep their existing HCL
functions while a new test extracts the VPC and network ACL blocks as dependencies
and supplies the entries through literal config maps. No whole-file rewrite is
needed. A direct `Config: fmt.Sprintf(...)` does not establish order coverage.

For a region prerequisite that is equivalent to the existing region helper, use
the resource's supported regions in `PreCheck`, following the existing tests:

```go
PreCheck: func() {
    testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
},
```

Use the actual supported region list; the helper can select a fallback region,
so it does not replace an arbitrary condition with identical behavior. Inline
`if` guards and runners inside control flow remain unsupported. Explicit
`t.Skip`, `t.Skipf`, `t.SkipNow`, `t.Fatal`, `t.Fatalf` or `t.FailNow` anywhere in
the test function, including local callbacks and `Check`, disqualify that test.
This conservatively includes unreachable calls; helper bodies are not interpreted.
The acceptance runner must report a real PASS for the reorder test in a supported
environment: SKIP is not PASS. The static node does not verify the API's remote
ordering behavior or replace that acceptance run.

## Local checks

```sh
go test ./scripts/collection-order/internal/coverage
go run ./scripts/collection-order -resource=alicloud_example
# Or use full SHAs, with HEAD checked out at the PR head:
go run ./scripts/collection-order -base="$DIFF_BASE" -head="$DIFF_HEAD"
```

The workflow uses the existing pinned runner policy, exact-head checkout and
read-only token permissions. It requires no cloud credentials and does not write
back to the PR.
