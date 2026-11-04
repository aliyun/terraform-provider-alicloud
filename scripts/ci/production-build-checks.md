# Production build & dependency-graph checks

These checks close the gap where product-migration PRs could pass CI without
ever proving the root provider main links or that the production dependency
graph stays free of the SDK v2 `helper/resource` test package.

## What each check does

| Entry | Command | Proves |
| --- | --- | --- |
| `make provider-build` | `go build -o /dev/null .` | The root `main` package **and every transitive production dependency it links** (including the split `alicloud/service/*` subpackages) compile and link. This is a real `go build`, not a renamed `vet`. |
| `make check-production-deps` | `scripts/ci/check-production-deps.sh` → `go list -deps -f '{{.ImportPath}}' .` | The **complete production dependency graph** of root `alicloud` does not contain `github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource` (direct or transitive). This is a Go-toolchain graph query, not a grep over a diff. |

### Why `helper/resource` must stay out of the production graph

The test line runs `terraform-plugin-testing`. Both
`terraform-plugin-sdk/v2/helper/resource` and
`terraform-plugin-testing/helper/resource` register a `-sweep` flag in `init()`.
A provider binary that links both dies at startup with
`flag redefined: sweep`. Production code uses
`terraform-plugin-sdk/v2/helper/retry` (flag-free) for retry helpers; tests
use `terraform-plugin-testing/helper/resource` for the harness.

## How this differs from `test-framework-guard.sh`

`test-framework-guard.sh` scans only **newly added** import lines in the PR
change range (a diff-level check). It stops new violations from creeping into
edited files but cannot prove the **whole** production graph is clean — a file
that has not been migrated and is not in the diff could still carry the import.

`check-production-deps.sh` queries the **entire** production dependency graph of
root `.` (main) with `go list -deps`, so already-migrated and not-yet-migrated
files are both covered. If any production source path introduces the forbidden
package, the link graph of `.` contains it and the check fails, regardless of
whether the offending file is in the PR diff.

## Where the checks run

- **GitHub CI** — `ProductionBuild` job in
  `.github/workflows/acctest-terraform-basic.yml` runs the underlying
  `go build -o /dev/null .` and `scripts/ci/check-production-deps.sh` directly
  (not via `make`, because the route runner for a fork PR can be `ubuntu-latest`,
  where `make` is not guaranteed to be on PATH). It runs when the route
  classifier marks the
  PR relevant (product-Go changes) OR when a `BuildInputs` job detects that the
  change range touches build-input files (CI workflows, GNUmakefile,
  `scripts/ci/*`, the local-ci-check and `*check.sh` scripts) — so a pure
  CI/Make/script PR is itself covered instead of being skipped. The job is
  additive: the existing `Compile` (vet) and `TestFrameworkGuard` jobs are
  unchanged.
- **Local / FC quick check** — `make ci-check-quick`
  (`scripts/local-ci-check.sh --quick`) runs `make check-production-deps` as a
  new `Part 1.6`, so the PR前 double-check (`submit_quick`) and every product
  migration PR reuse it automatically. The quick path skips the build
  (`--skip-build`), so `provider-build` is NOT exercised here; the production
  build evidence comes from the GitHub CI `ProductionBuild` job below.
- **Remote init probe** — `probe_init` compiles+inits a single service package
  (`alicloud/service/cloudauth`) via a no-such-test-case `go test` run and
  returns INITIALIZATION_ONLY (`no tests to run`). It does NOT compile or link
  the root provider main; do not cite it as root-main build evidence. The
  root-main build is proved by `make provider-build` in the `ProductionBuild`
  CI job.

## Reuse for product-migration tasks

Product-migration PRs reuse these checks with **no extra wiring**, provided
the entries are already present in the migration PR's tested revision and
applicable CI configuration:

1. When the migration PR's head already contains this base PR's changes, its
   GitHub CI runs `ProductionBuild` (`provider-build` + `check-production-deps`)
   on the relevant path.
2. The PR前 Workflow double-check (`submit_quick` → `make ci-check-quick`) runs
   `check-production-deps` once the migration head contains the `Part 1.6`
   wiring.
3. The production build (`make provider-build` / `go build -o /dev/null .`) is
   the entry that proves the root main links; the cloudauth init probe only
   covers one service package, so do not treat it as root-main build evidence.

A migration PR whose head does NOT yet contain these checks (e.g. an older
product branch predating this base PR) is not covered by results from this
base branch — the checks must be present in the actual tested revision. Once
present, a migration PR that introduces a stray
`terraform-plugin-sdk/v2/helper/resource` production import — or breaks the
root link — is caught at PR前 / CI time, instead of forcing QA to re-collect
root-main build and dependency-graph evidence per resource PR.

## Local usage

```bash
make provider-build            # real production build of root main
make check-production-deps     # production dependency graph assertion
make ci-check-quick            # quick path now includes the dep check
PROD_DEPS_INPUT=<file> scripts/ci/check-production-deps.sh   # fixed graph (tests)
FORBIDDEN_DEPS="<path> ..." scripts/ci/check-production-deps.sh  # extend scope
```

The fixture tests live in `check-production-deps_test.sh` and exercise the
assertion logic without a Go toolchain via `PROD_DEPS_INPUT`.

## Cloudauth no-tests initialization check

`scripts/ci/initialization-check.sh` is the local-ci-check equivalent of the
remote `probe_init` init probe. It closes the gap where `ci-check-quick` left
the cloudauth test-binary initialization to a separately respawned probe task,
and where some QA treated `go build` / `go vet` success as proof the test
binary had executed.

| Entry | Command | Proves |
| --- | --- | --- |
| `make check-initialization` | `scripts/ci/initialization-check.sh` → `go test ./alicloud/service/cloudauth -run '^$' -count=1 -v` | The standard cloudauth test binary **compiles, links and initializes** (`init()`/`TestMain`) with **no test body running** — no cloud resource is created. INITIALIZATION_ONLY, `countsAsAccTest=false`. |

### Why `-run '^$'` and `-count=1`

- `-run '^$'` matches no `Test*` function, so `go test` builds the test
  binary, links it, runs `init()` (and `TestMain`, if any) and then reports
  "no tests to run". No acceptance case executes and no resource is created.
- `-count=1` disables the Go test result cache, so a cached pass cannot mask
  a real compile/link/init failure — the test binary is actually executed.
- A pass additionally requires a genuine go-test "no tests to run" signal:
  the wrapper rejects any output containing `=== RUN` / `--- PASS` /
  `--- FAIL` / `--- SKIP` (a test body ran) or `?` (no test files; binary
  never built), AND requires the real no-tests marker — `testing: warning:
  no tests to run` and/or the `[no tests to run]` suffix on the `ok` line —
  that go test only emits when the binary truly compiled, linked and
  initialized with no `Test*` function matched. A bare `ok`/`PASS` line
  without that marker is rejected, so a stub or a misconfigured filter cannot
  fake PASS. A non-zero exit (compile error, link error, init panic,
  build-failed) fails the check and the enclosing quick/CI run, with the
  original stdout/stderr and the exact command printed.

### Scope (what it is NOT)

- NOT the root main build — use `make provider-build` (`go build -o /dev/null .`).
- NOT a real ACC — use the cloud acceptance suite. The receipt explicitly
  marks `countsAsAccTest=false`.
- NOT per-resource — one standard service test package, run once per quick.

### Where it runs

- **Local / FC quick check** — `make ci-check-quick`
  (`scripts/local-ci-check.sh --quick`) runs the check as `Part 5b`,
  unconditionally (not gated by `--skip-build`), so `--quick` does not bypass
  it and it runs exactly once per quick.
- **Standalone** — `make check-initialization` / `scripts/ci/initialization-check.sh`.
- **Remote init probe** — the separate `probe_init` Provider Run entry remains
  the standalone remote probe; this CI entry lets the PR前 quick
  (`submit_quick`) collect the init evidence in the same quick run instead of
  respawning a probe task per migration.

### Local usage

```bash
make check-initialization                                  # real cloudauth init check
make ci-check-quick                                        # quick now includes Part 5b
scripts/ci/initialization-check.sh                        # direct
INIT_CHECK_PACKAGE=./alicloud/service/ecs scripts/ci/initialization-check.sh  # other package
# fixture tests (no Go toolchain, via a stub command):
scripts/ci/initialization-check_test.sh
CI_INIT_CHECK_REAL=1 scripts/ci/initialization-check_test.sh  # + real go test smoke (opt-in)
```

The fixture tests live in `initialization-check_test.sh` and exercise the
pass/fail decision logic (compile failure, init panic, ran-a-test rejection,
no-test-files rejection, receipt format, single-probe wiring) without a Go
toolchain via `INIT_CHECK_GO_TEST_CMD`, mirroring `PROD_DEPS_INPUT`. The
genuine `go test` execution is validated by a real `ci-check-quick` run.
