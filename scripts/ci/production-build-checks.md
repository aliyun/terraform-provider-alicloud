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
