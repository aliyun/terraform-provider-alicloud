# CI fixture compatibility with older PR branches

A workflow can advance on master while checking out a PR head that predates its
new regression scripts. Run `run-ci-fixtures.py` after `resolve-change-range.sh`:
the latter validates the full history, exact event base SHA, fetched base ancestry
and merge base. The workflow loads a missing runner from that exact base SHA.

Complete HEAD fixtures execute from HEAD, including any PR changes and failures.
For a legacy PR with missing fixtures/dependencies, the runner archives a coherent
exact-base snapshot into RUNNER_TEMP, executes the real tests, and removes it.
The provider checkout, index and source files remain unchanged. A temporary Git
root lets Bash fixtures locate their siblings without binding to the PR tree.

Missing fixtures on push, invalid revisions, explicit deletions, or PR changes to
CI fixture inputs with missing dependencies fail clearly. Python discovery must
find at least one test. No check is replaced with a success marker or a skip.

## Workflow audit

| Step | Finding and handling |
| --- | --- |
| tflint range/lint/initialization regressions | New initialization fixtures were called unconditionally. Use the complete lint fixture profile, including local-ci-check.sh and its initialization probe. |
| Document regression | unittest discovery could pass with zero tests on old heads. Use the document profile and reject zero tests. |
| Consistency scope regression | Missing helpers previously skipped the tests. Use the scope profile. |
| Basic writable-cache regression | Missing tests previously skipped. Use the writable-cache profile. |
| Compact Go cache round trip | Already archives coherent exact-base workflow/test inputs and rejects changed or deleted inputs. Preserve that existing path. |
| Cache permission preparation / seed restore | Existing exact-base/administrator fallbacks remain unchanged. |
| PR title checker | Already loads its checker and sibling dependencies from the exact PR base. Preserve it. |
| Provider checks / build-input detector / production dependency guard | Present in the inspected #11277 head; keep checking PR source with existing range/scope checks. |
| CI Policy Test | Runs on CI-path changes, testing the proposed HEAD; add compatibility runner/tests to its path filters and regression command. |

Regression tests execute the actual inline bootstrap in real divergent Git
histories, including lazy blobless clones, a base newer than the old head, a
later branch tip, fixture failures, explicit deletions, missing dependencies,
zero tests and unchanged PR worktrees.
