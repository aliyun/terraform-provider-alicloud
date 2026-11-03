# Multi-commit pull requests

PRs may contain multiple commits. CI evaluates the net change from the PR
base/head merge base to the event head, not just the last commit. Squashing
at merge time is a separate maintainer decision.

## Audit of PR check entrypoints

- `pull_requests.yml`: removes the `Pull Request Max Commits` job. The title
  checker already receives `DIFF_BASE`/`DIFF_HEAD` from `resolve-change-range.sh`.
- `acctest-terraform-basic.yml`: basic checks, formatting and compatibility
  receive the complete range. The compatibility CLI previously ignored it
  internally and compared `HEAD^..HEAD`; it now honors it in flat, split SDKv2,
  framework and direct-resource paths, failing closed on missing/invalid CI
  revisions. The Python wrapper and direct-Go fallback both inherit the range.
- Lint, vet, document, consistency, test coverage and collection-order workflows
  already consume the resolved range, either directly or via `check_scope.py`.
  Their selection regression suite includes changes before the last commit,
  additions, modifications, renames and deletions.
- The legacy `field_check.sh` entrypoint now uses the same explicit range for
  both its diff and resource selection.
- Standalone local helpers (`local-ci-check.sh`, `local-example-check.sh`,
  `gofmtcheck.sh`, `goimportscheck.sh`) are not the GitHub PR workflow entrypoints.
  Their intentional local last-commit defaults remain; remote quick checks must
  supply `CI_CHECK_BASE`, as before. The resolver's new-branch push fallback is
  also separate from its PR merge-base calculation.

No provider behavior or cloud acceptance-test execution is changed here.
The master branch protection inspected for this change requires `Basic Policy
Gate`, not the removed commit-count check. Repository protection is not changed
by this PR.

## Regression validation

`CI Policy Test / multi-commit-pr-test` runs on relevant PR and master changes:

- Real diverged/multi-commit Git histories verify merge-base selection, inclusion
  of earlier commits, exclusion of base-only changes, blobless forks and rejection
  of shallow or unrelated histories (`resolve-change-range_test.sh`).
- The real compatibility executable detects an Optional-to-Required change in
  an earlier commit followed by a documentation-only commit, including split
  registrations and direct-resource mode; the old last-commit scope misses it.
- A later revert is accepted based on the net diff. Missing, partial and invalid
  CI revisions are rejected; explicit CLI ranges and legacy field selection work.
- Existing scope, compatibility and title regression suites also run.
