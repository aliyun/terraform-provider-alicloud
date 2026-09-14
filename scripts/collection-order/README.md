# Collection ordering CI

The **Collection Order** PR check uses the real provider resource's SDK
`Resource.Diff` without credentials. It checks that permuting a semantically
unordered collection does not change its plan, while adding, removing or
replacing members does. Ordered lists must retain their ordering diff.

Run from the repository root:

```sh
go test -p=1 ./scripts/collection-order -count=1
go run ./scripts/collection-order -base origin/master -head HEAD
```

Omit `-head` to include tracked working-tree changes. CI uses the repository's
existing change-range resolver and pinned runner policy for fork PRs.

## Required fixtures

Every configurable multi-member `TypeList` or `TypeSet` in a directly edited
resource or data-source implementation, or an added/changed provider registration,
needs a `cases.json` entry. Nested collections are included;
computed-only fields and single-member containers are not reordered. The
children of a single-member block are still checked. Existing registered cases
run on every PR, including changes to the checker itself.

Coverage includes the complete directly changed implementation, not only schema
lines. Filenames normally match registration names; four exceptions are mapped
explicitly:

* Data source `ess_lifecyclehooks` → `ess_lifecycle_hooks`.
* Data source `yundun_dbaudit_instances` → `yundun_dbaudit_instance`.
* Resource `ess_lifecyclehook` → `ess_lifecycle_hook`.
* Resource `edas_application_package_attachment` → `edas_application_deployment`.

`alicloud/data_source_alicloud_common.go` is explicitly excluded because it only
contains shared hash/schema helpers. All other unresolved implementation files
fail; only a missing file whose former registration was removed is skipped.
The regression test classifies every non-test resource/data-source filename
against the real provider registry. Git renames are expanded into deletion and
addition paths, so moving a file does not remove coverage while its resource
remains registered.

Shared helper changes do not demand new fixtures or automatically promote
unrelated historical defects to failures. Reviewers must still require relevant
fixtures when a shared schema helper introduces a collection; this check does
not compare schema definitions.

An entry has:

* `resource`: a provider resource name, or `data.alicloud_...` for a data source.
* `path`: the collection path, including fixture indices for nested blocks,
  for example `rules.0.targets`.
* `order`: `unordered` or `ordered` for a list. Sets are automatically unordered.
* `reason`: explain the API's ordering semantics. Being a `TypeList` alone is
  not evidence that its order is meaningful.
* `config`: valid sibling configuration and parent blocks needed by the SDK.
* `members`: at least two distinct, valid values. Supply additional valid values
  when available, including enough to satisfy `MinItems`. All but the last form
  the initial collection, adjusted to at least two members and the schema's size
  bounds. An unused member supplies addition and replacement controls. Four
  values allow a three-member initial collection and an extra rotation check;
  `MaxItems=2` still uses a two-member initial collection.

SDK validation rejects invalid fixtures. Both state and configuration orders
are exercised, with the complete ancestor containers so outer set hashes are
preserved. Membership checks run in both directions. With only two valid values
`A` and `B`, reordering uses `[A,B]` and `[B,A]`; membership controls use
`[A,B]` versus `[A]` and replacement uses `[A]` versus `[B]`, when those lengths
are allowed. Add/remove controls forbidden by schema size bounds print `N/A`
with those bounds. Missing a legal replacement sample fails: for a fixed-size
collection, provide at least one distinct member beyond its size. The checker
never invents duplicate members or treats missing data as proof of a finite
API domain. A truly immutable domain needs a dedicated constraint-specific
verification before it can be supported; there is no fixture exemption flag.
Hash collisions or a suppressor that hides legal membership changes fail the
controls that run. Failures identify the resource, attribute and SDK
diff, including replacement behavior.
The full root diff is inspected, including changed parent and sibling fields.
Only identical computed unknowns from the unchanged control and unchanged input
entries repeated by SDK v1 replacement planning are ignored.

The checker keeps field suppressors and the root `CustomizeDiff`. A resource
with `CustomizeDiff` additionally requires an audited `nil_meta_reason`
explaining why its fixture exercises the relevant callback behavior without a
configured client. Do not declare this if the callback requires credentials or
provider features: add a dedicated offline setup for that resource first.

## Existing defects and limits

Two exact known defects are currently recorded in `main.go`:
`alicloud_ecd_desktop_group.end_user_ids` produces an ordering diff, and
`alicloud_ecd_desktop.end_user_ids` produces a replacement diff. These checks
still execute and print **KNOWN DRIFT** warnings, counted separately from
passes. Direct changes or changed registrations for either resource fail on that
defect. Two audited product-file dependencies also enforce it:

* `service_alicloud_ecd.go`: `DescribeEcdDesktop` returns `EndUserIds` to desktop Read.
* `service_alicloud_ecd_v2.go`: `DescribeDesktopGroupDescribeUsersInGroup` collects
  and sorts `EndUserIds` for desktop-group Read.

These are conservative file boundaries within ECD, not an inferred call graph.
Any edit to the listed service file selects its associated defect; unrelated
VPC/ECS services and common helpers do not. Update the entries if those Read
helpers move. Every fixture still runs on every PR. Once fixed, the check
requires removal of its baseline entry.
Fixtures cannot add new baseline exemptions, and new unregistered collections
fail. An empty or baseline-only suite cannot pass.

This is an offline SDK schema planning check. Canonical state models a
normalized API response; the checker does **not** call Configure, Read, CRUD,
or a cloud API, and does not prove real API ordering, pagination, Terraform Core
apply/refresh behavior, or an end-to-end drift-free lifecycle. Cases provide
explicit ordering contracts; the checker cannot infer them from API docs.
Use acceptance tests for those lifecycle guarantees. This CI addition creates
no cloud resources and requires no acceptance-test credentials.
