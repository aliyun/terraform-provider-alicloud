# Release Pipeline

Automated, gated release pipeline for `terraform-provider-alicloud`. It runs
the pre-release test gate, generates and submits the changelog, pauses for a
mandatory human approval, and then creates the release tag. The actual
goreleaser build (multi-arch binaries, `SHA256SUMS`, GPG signature, registry
manifest, draft release) is performed by the existing
[release.yml](../../.github/workflows/release.yml) workflow, triggered by the
`v2.*` tag this pipeline creates. The [`.goreleaser.yml`](../../.goreleaser.yml)
contract is intentionally left untouched.

## Pipeline at a glance

```
prepare ──▶ pre-release-tests ──▶ generate-changelog ──▶ submit-changelog
        │           │
        ├──▶ module-tests
        ├──▶ acc-tests (optional)
        │
        └──▶ release-approval (human gate)
                    │
                    ▼
              tag-and-release ──▶ tag v2.X.Y ──▶ release.yml (goreleaser)
```

All jobs use self-hosted runners registered in
[.github/actionlint.yaml](../../.github/actionlint.yaml) plus the dedicated
`terraform-releaser` runner for the goreleaser build.

## Files

| File | Role |
|------|------|
| `.github/workflows/release-pipeline.yml` | The orchestration workflow. |
| `scripts/release/release_common.sh` | Version validation, last-tag resolution, change-range base computation, module-affected detection. |
| `scripts/release/generate_changelog.py` | CHANGELOG generator (`--update` merges into Unreleased; `--release` finalizes the version). Pure Python stdlib. |
| `scripts/release/README.md` | This document. |

## Trigger

`workflow_dispatch` only. From the GitHub Actions UI (or `gh workflow run`),
provide:

| Input | Type | Required | Meaning |
|-------|------|----------|---------|
| `version` | string | yes | Target version, e.g. `2.2.0` (with or without `v` prefix). |
| `run_acctest` | boolean | no | Default `false`. When `true`, run the acceptance-test job before approval. |
| `acctest_cases` | string | no | Comma-separated `TestAcc` function names; empty runs the full suite. |

## Stages

### 1. prepare
Validates the version (`rc_validate_version`), locates the most recent `v*`
release tag (`rc_last_release_tag`), computes the merge-base of that tag and
HEAD (`rc_release_diff_base`) and exposes it as `diff_base`. The same base
anchors `CI_CHECK_BASE` for `make ci-check-quick` and `--from` for the
changelog generator. When no prior tag exists (maiden release), the root
commit is used so the whole history is the change range.

### 2. pre-release-tests (AC1)
Runs `CI_CHECK_BASE=<diff_base> make ci-check-quick`. This reuses the existing
quick CI (`scripts/local-ci-check.sh --quick`) which performs quality checks
and the **changed-commit unit tests**: changed `_test.go` files select their
`TestUnit*` functions, changed `resource_*.go` / `data_source_*.go` select
companion tests, and shared-helper changes broaden to package-wide `TestUnit*`
(see `scripts/ci/changed-unit-tests.md`). This covers the requirement to run
"all changed commits' test cases before release".

### 3. module-tests (AC2, conditional)
Runs only when `rc_module_tests_affected` detects changes under
`examples/`, `modules/`, or `website/`. Performs `terraform fmt -check
-recursive examples/ modules/` — a static, credential-free module gate.
Example tests that create real cloud resources are not run here; they belong to
the optional acceptance-test stage so the pipeline can proceed without cloud
credentials when only formatting is at stake.

### 4. acc-tests (optional)
Runs only when `run_acctest=true`. Guards on the `ALIYUN_ACCESS_KEY_ID` /
`ALIYUN_ACCESS_KEY_SECRET` secrets; fails fast if they are absent. Executes
`TF_ACC=1 go test` against the requested cases (or the full `./alicloud/...`
suite). This is the cloud-backed acceptance gate; enable it for releases that
touch resource CRUD behaviour.

### 5. generate-changelog (AC3)
Runs `generate_changelog.py --update` to merge freshly parsed commit subjects
into the `## <version> (Unreleased)` section of `CHANGELOG.md`. Entries are
categorised from PR-title prefixes (`New Resource:` / `New Data Source:` /
`New Ephemeral Resource:` → FEATURES; `resource/alicloud_*` /
`data-source/alicloud_*` / `docs:` → ENHANCEMENTS; `fix(...)` / fix hints →
BUG FIXES; `ci:` / `chore:` skipped) and formatted as
`- <summary>. ([#PR](https://github.com/aliyun/terraform-provider-alicloud/issues/PR))`,
matching both the existing CHANGELOG.md and the v2.x release body. Existing
manual entries are preserved; duplicate PR numbers are deduplicated.

### 6. submit-changelog (AC4)
Commits the updated `CHANGELOG.md` to a `changelog/v<version>` branch and opens
(or updates) a PR titled `ci: update CHANGELOG for v<version>` targeting the
release branch. The PR is left open for the human gate to merge.

### 7. release-approval (AC5, human gate)
`environment: release-approval`. **This is the mandatory human approval.**
Configure this GitHub environment with required reviewers in
*Settings → Environments → release-approval*. The pipeline pauses here and
never proceeds to tagging without an explicit approval. Skipped optional jobs
(no module changes, no changelog diff, `run_acctest=false`) do not block this
gate — the condition is `!cancelled() && !failure()`, which treats a skipped
optional job as a pass.

### 8. tag-and-release (AC6)
After approval, merges the changelog PR (squash, then delete branch) into the
release branch, then creates the annotated tag `v<version>` on the merge
commit and pushes it. The `v2.*` tag push triggers the existing
[release.yml](../../.github/workflows/release.yml), which runs goreleaser on
`terraform-releaser` to build the multi-arch binaries, `SHA256SUMS`, GPG
signature and registry manifest, and create the draft release.

## Prerequisites (one-time repository setup)

1. **GitHub environment `release-approval`** — *Settings → Environments → New
   environment*. Add the required reviewers who may approve a release. This is
   the human gate; without it the workflow cannot publish.

2. **Self-hosted runners** with labels `terraform-ci-route`,
   `terraform-ci-heavy`, and `terraform-ci-trusted` (already registered per
   `.github/actionlint.yaml`), plus the existing `terraform-releaser` runner
   for the goreleaser build.

3. **(Optional) ACC secrets** — add `ALIYUN_ACCESS_KEY_ID` and
   `ALIYUN_ACCESS_KEY_SECRET` (and optionally `ACC_REGION` variable) as
   repository secrets to enable the optional acceptance-test stage. Leave
   them unset to skip ACC; the pipeline still runs the unit and module gates.

4. **(Optional) `V2_RELEASE_BRANCH`** repository variable — defaults to
   `master`. Set it only if the v2 release line moves to a different branch
   (must match the value used by `release.yml`).

No changes to `release.yml` or `.goreleaser.yml` are required. The pipeline
extends the release process **before** the tag; the goreleaser contract after
the tag is untouched.

## Relationship to existing CI

- **`release.yml`** — untouched. Triggered by the `v2.*` tag this pipeline
  creates. Performs the goreleaser build.
- **`pull_requests.yml`** — untouched. Runs per-PR title validation and quick
  CI on `terraform-ci-policy`. This pipeline reuses the same `ci-check-quick`
  target and `scripts/ci/resolve-change-range.sh` semantics, but anchored to
  the release range (last tag → HEAD) rather than a single PR.
- **`scripts/changelog-links.sh`** — the generator emits entries that are
  already link-ready (`[#PR](url)`), so `changelog-links.sh` is a no-op for
  generated entries while still working for manually added `[GH-nnnn]`
  shorthand.

## Local / non-sandbox compile policy

This pipeline runs on the provider repository's own self-hosted GitHub Actions
runners (the same infrastructure that already runs `pull_requests.yml` and
`release.yml`), not on a Jarvis RD sandbox. Provider `go vet` / `go build` /
`go test` execute there as part of the repo's CI, consistent with the existing
workflows. The Jarvis RD development flow continues to use the authorised
remote Provider Run v2 entry for pre-PR checks; this pipeline is the
release-time gate, not a replacement for per-PR development verification.

## Usage examples

```bash
# Dry run the changelog generator locally (no write):
python3 scripts/release/generate_changelog.py --from v2.1.0 --version 2.2.0

# Trigger the pipeline via gh:
gh workflow run release-pipeline.yml \
  -f version=2.2.0 \
  -f run_acctest=false
```
