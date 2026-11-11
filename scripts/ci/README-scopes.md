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

#### Explicit remaining differences

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
- Direct documentation edits still run all fenced blocks on that page and
  inject the existing provider/fixture setup. Full scenario dependency/wording
  assertions and the workflow's stricter original-HCL contract remain separate.
- ACube preproduction, fixed case concurrency 2, the existing exact-source reuse
  protocol, log-warning policy and trigger events remain platform differences.
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
