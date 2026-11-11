# CI change scopes

`check_scope.py` validates `DIFF_BASE` and `DIFF_HEAD` before selecting work. Both
must be full, locally available commit SHAs, HEAD must match, the checkout must
not be shallow, and base must be an ancestor of head. The workflow resolver
supplies the PR merge-base rather than only the last commit. Git failures fail
checks; they never become an empty successful selection.

Rename detection is disabled for selection so a move contributes both its old
path and its destination. File-based checks inspect existing destination files.
Package-based checks load complete Go packages, including required imports.
Deleted files and rename sources still select their old package when other Go
files survive there, so removal of a referenced declaration cannot silently skip
package validation.

- Errcheck keeps its analyzers and reviewdog integration, selects changed
  non-test Go packages, and compares findings against the validated base.
  Shared module/tool configuration changes and unusual shell flag paths fall
  back to the existing full scope. Test-only changes remain covered by vet,
  rather than changing errcheck's existing `--tests=false` policy.
- Formatter checks selected provider Go files without `go mod tidy`, retains
  the existing goimports version, and propagates formatter/tool failures.
- Vet uses the validated range, including test-only package changes. Shared
  inputs select the full provider and any changed packages outside it.
- Consistency, attribute testing coverage, and TypeList coverage keep their
  original Go checker logic and full diff. The lightweight wrapper only avoids
  provider compilation when changes are clearly unrelated. All provider Go
  files, shared helpers, checker sources, and CI changes remain relevant.
- Documentation spelling and Terraform block formatting are file-scoped on
  ordinary PRs. Master and shared documentation configuration/tool changes
  retain their original full checks. Content and Markdown lint remain diff-based
  to avoid newly enforcing historical findings in unchanged pages.
  The original global `.md` link check remains global to avoid losing checks
  on unchanged sources, and changed `.markdown` sources are also checked.
  Link checks make at most three attempts when all reported failures have no
  HTTP response and include a transient timeout, connection reset or DNS retry
  error. Every attempt remains in the log. HTTP failures, tool failures and
  failures after the last attempt still fail the check.

Older PR heads missing these wrappers use workflow fallbacks that execute the
previous checker instead of skipping it. Direct invocations of the existing Go
checkers and `scripts/vetcheck.sh` remain available.

Run `python3 -m unittest discover -s scripts/ci -p test_check_scope.py` for real
Git-history regression coverage. `smoke_ci_checks.py` additionally requires
Go, goimports, golangci-lint, and the documentation tools, and demonstrates both
success and deliberately rejected changes. The fork validation workflow pins
those tools and runs this smoke test without upstream self-hosted runners.

These changes do not remove Go dependency loading, and changes in the large
root provider package can still be expensive. They do not replace cloud
acceptance tests or expand the runtime schema checkers' existing capabilities.

## Code Integration Checks on ACube preproduction

`.aoneci/acctest-terraform-integration.yml` retains both job names and runs
`scripts/ci/acube-integration.py`. It freezes the complete Code change range,
imports head/base via full SHA, groups unique AccTests by package (concurrency 2),
and runs each affected documentation example with `hclAction=example` including
Import/generated config, deprecated attributes and cleanup. The two jobs report independently,
so an AccTest failure cannot hide the Example result. The YAML declares the MR
`accepted` event and master pushes. Provider paths include all of `alicloud/**`,
including shared test-support packages and fixtures. GitHub workflows are unchanged.

Configure the protected repository secret `ACUBE_PROVIDER_RUN_TOKEN` for
`https://pre-acube.aliyun-inc.com`; cloud account keys stay in FC. The authenticated
capabilities check must pass before submissions. Deploy the compatible server and
isolated preproduction FC before enabling the YAML. Keep the ACube receipt IDs
and private `.acube-ci/` evidence; the same CI Run/plan/item key returns the original
receipt after lost responses or process restarts. A new CI Run is an explicit new
attempt. Do not rerun until the original cloud execution/cleanup is reconciled.

Every terminal task prints a redacted result to the job log and appends the same
evidence to the Aone Markdown summary. Failures include the remote error, HCL
phases, cleanup status, unsuccessful AccTests and a bounded diagnostic excerpt
from remote stdout. Log download failures preserve the original task failure.
Both jobs publish only `.acube-ci/diagnostics/` as downloadable artifacts, even
when the check fails; original requests and signed download URLs remain private.

SKIP, missing/duplicate cases, mismatched sources and unconfirmed cleanup fail
closed and remain visible. Neither a PlanOnly case nor top-level ACube success
proves complete acceptance. Required external prerequisites must be resolved or
explicitly accepted; this adapter does not silently waive them.

### ACC selection aligned with the agent workflow

The scope follows the workflow's
[shared test standard](https://code.alibaba-inc.com/terraflow/terraform-agent-workflow/blob/9a15deb6f798a2f46cfb8c5ff557e0eff635328e/skills/provider-resource-standards/references/test-standards.md#L185).
`accscope/main.go` parses committed base/head Go declarations using only the Go
standard library. It does not load dependencies or compile Provider packages.
The CI adapter runs it with `GO111MODULE=off` and `GOWORK=off`, independently of
the candidate Provider's Go module/toolchain requirements.

- Resource/data-source implementation changes select the type's acceptance
  functions in its companion and supplemental test files, excluding separately
  implemented types. Sources and selectors retain their actual Go package.
- A test-only change selects newly added or changed `TestAcc*` declarations.
  Test helpers/config constants are followed transitively to their consumers;
  an unrelated case in the same file is not selected. A change confined to
  `TestUnit*` does not request cloud ACC. Changing TestMain/test initialization
  selects the package's acceptance cases.
- Service changes follow declarations and method references in both snapshots,
  including removed call edges. Reached resource implementations select that
  type's cases and examples. Statically declared receiver types distinguish
  same-name methods; uncertain receivers require review.
- Shared test helpers in `_test.go` and `alicloud/acctest/` are traced across
  packages and import aliases. Literal package-relative fixture paths, including
  static `filepath.Join` arguments, select their actual consumers.
- Deletions and both paths of renames participate. Removed acceptance cases,
  unresolved helper/service/fixture consumers and unknown build/receiver effects
  block rather than become a successful empty selection. Source-only changes to
  CI tooling do not themselves request Provider ACC.

`.acube-ci/plan.json` records the exact cases, affected resource files, selection
reasons and blockers. A blocked plan writes `diagnostics/scope.json` and fails
before **any** ACube capability/source/task requests. Test-only changes do not
automatically submit the resource's documentation examples. Scope regression
tests are part of `CI Adapter Regression`; they use synthetic source and real
local Git snapshots without cloud calls or Provider builds.

### Documentation scenes and source evidence

The documentation job follows the workflow's
[document example contract](https://code.alibaba-inc.com/terraflow/terraform-agent-workflow/blob/9a15deb6f798a2f46cfb8c5ff557e0eff635328e/skills/invoke-terraform-acc-test-remote/references/document-examples.md).
Direct documentation changes select changed scenes, assertions and declared
inputs. Implementation/service impact selects all scenes of the affected type.
Unchanged HCL under a prose-only edit does not run again; the report explicitly
keeps wording assertions NOT_VERIFIED for independent QA. Test-only changes
continue to select ACC without selecting documentation lifecycles.

Before execution, review a version-1 JSON contract at
`ci/assets/document-examples/r/<resource>.json` (`d/<data_source>.json` for data
sources). The matching page is `website/docs/r/<resource>.html.markdown`. For
example, a shared Provider fence and two independent resource fences use cases
with blocks [1, 2] and [1, 3]. A change to the shared fence selects both.
Each case requires a unique id, 1-based HCL blocks, originalSha256, ready=true
and nonempty assertions. Each native assertion is {"id": "...", "check": "..."},
where check is lifecycle, import or cleanup. Other API/business assertions
remain blocked for QA rather than being inferred from native task success.

Example contract for a page whose first fence contains the complete scene:

```json
{
  "version": 1,
  "cases": [{
    "id": "basic",
    "blocks": [1],
    "originalSha256": "<SHA-256 of the exact first HCL fence contents>",
    "ready": true,
    "assertions": [
      {"id": "apply-and-converge", "check": "lifecycle"},
      {"id": "import-no-changes", "check": "import"},
      {"id": "destroy-and-empty-state", "check": "cleanup"}
    ]
  }]
}
```

Optional sourceFiles and fixtures contain literal tracked repository-relative
.tf/.hcl paths. The original scene is the selected fences, followed by sourceFiles,
joined in declared order with one newline. Its originalSha256 is SHA-256 of the
exact UTF-8 bytes (including original line endings). Fixtures are appended to the
test copy. Optional mainTf replaces the test copy after explicit adaptation; it
never replaces the original evidence. Dependencies shared across cases are
declared in every consuming case. Every extracted HCL fence and linked HCL
source must be accounted for. The full change range follows declared inputs,
including input-only edits, deletions and both paths of renames. CI triggers
include scenario contracts and all .tf/.hcl paths.

No scene is automatically ready merely because a fence was extracted. Missing
contracts, unresolved grouping/inputs, stale source hashes, unsupported assertions
and removed scenes fail before any ACube requests for the documentation job.
A page without executable examples needs an explicit version=1 contract with
cases=[] and a nonempty notApplicableReason. A removed page requires scope review.
Known unlabelled/indented HCL is blocked for explicit transcription into a
reviewable fenced/source-file scene.

The original scene must explicitly declare exactly one alicloud Provider source,
aliyun/alicloud or registry.terraform.io/aliyun/alicloud. The adapter does not
inject a missing declaration. Published version constraints can remain in the
source, while the reviewed test mainTf must omit them or pin the mirror version
0.0.1 / = 0.0.1. Both configurations are checked separately. Extra/implicit
Providers, modules, remote state, required variables and file-dependent or obvious
placeholder inputs block. This conservative declaration parser is not a complete
HCL parser; readiness still requires review, followed by remote init/validate.

The private plan retains originalMainTf, source spans/hashes, fixture hashes,
adaptationDiff and failed scene records. Public diagnostics expose only hashes,
locations and native assertion results, with proseAssertions=NOT_VERIFIED.
Successful native evidence must match sourceId, binary SHA-256 and mainTf hash;
include every build/init/validate/plan/apply/second-plan/Import/destroy/state-list/
plan-after-destroy phase in order; show no Import mutations and confirmed
cleanup. Lifecycle PASS does not prove arbitrary wording or remote deletion
semantics. Such observations still need the workflow's independent QA/API evidence.

### ACC evidence reuse and recovery

Both server capability names remain supported: accTestCaseReuse=exact-source-v1
and accTestContentReuse=whole-source-and-comments-v1. The adapter queries
POST /acctest/reusable-cases with the full frozen candidate. The service decides
source equivalence and runner validity; changing executable source does not
become reusable merely because a case name is unchanged. A case independently
passing in a failed original batch may still be reused when the service returns
it; the adapter does not demand whole-batch success.

All hits create no task. Partial hits submit only missing cases and retain the
full candidateRequest, reduced request and original case-to-task references.
A fresh reuse lookup at completion must still prove every reused case valid;
revocation/query failure fails closed without an automatic rerun. Original and
current task references are recorded separately. Reuse diagnostics cover only
the reused subset; the remainder has its own native result.

The ACC key depends on the run, package and frozen ACC request, independently of
documentation inventory/rules. Documentation blockers do not prevent the ACC
job from querying reuse. Same-run recovery preserves the original subset and
receipt, including UNKNOWN submission intent, instead of recomputing a different
remainder. Reusing an output directory with different source/input identity is
rejected; retain the old evidence and use a new directory after reconciliation.

Run `python3 -m unittest discover -s scripts/ci -p 'test_acube_*.py'` for the
selection, document contract, native evidence and reuse/recovery regressions.

### Explicit remaining differences

- General production helpers, request/retry/authentication code and module
  dependency changes require a maintainer's ACC scope decision under the shared
  standard. This adapter has no authoritative input for that decision: it emits
  a blocked plan, not automatic PASS or an implicit waiver. Provider registration
  changes also require review when their type scope cannot be inferred here.
- This is a source reference analysis, not a typed runtime call graph. Reaching a
  resource implementation conservatively selects its type's tests; it cannot
  prove which lifecycle steps exercise a changed branch. Factory/interface
  dispatch, reflection, generated/build-specific code, package initialization and
  dynamically constructed fixture paths may require manual impact review.
  Nonstandard naming/layout or dynamic test discovery is not silently expanded
  to all Provider resources.
- Arbitrary prose assertions, negative examples, business outputs and remote
  deletion semantics require QA/API observations. CI validates the declared
  native checks and records the remaining semantic scope as NOT_VERIFIED;
  it cannot establish those facts from lifecycle phases alone.
- ACube preproduction, fixed case concurrency 2, server-authoritative source
  equivalence/reuse, log-warning policy and trigger events remain platform differences.
  A CI success is not the workflow's independent QA acceptance or merge approval.

#### RocketMQ subscription acceptance

The two `basic4128` subscription cases stay in the complete plan. The adapter
submits ordinary cases with concurrency 2, then submits these two as a separate
serial batch. It requires `rocketmq-refund-v1` release receipts for both cases in
addition to native case-level PASS; SKIP and a merely accepted refund cannot pass.
The preproduction runner must enable the verified test account and the matching
Provider test-only cleanup hook. Production subscription deletion and existing
GitHub workflows are unchanged. Recovery uses the private task-specific cleanup
intent after the original FC invocation has stopped, without replaying the test.
