#!/usr/bin/env bash

# Boundary-gated runner for scripts/ci/initialization-check_test.sh.
#
# #11312 added scripts/ci/initialization-check_test.sh and wired it into the
# "Test change range and lint scoping" step of
# .github/workflows/acctest-terraform-lint.yml, which checks out the PR head and
# runs the script unconditionally. PR heads forked before #11312 do not contain
# that file, so the bare call exits 127 (command not found) against the new
# workflow baseline and produces a noise failure after the resolve-change-range
# and run-tflint regressions pass.
#
# This wrapper keeps the regression real without faking a pass or letting a PR
# bypass it. The decision uses the trusted PR merge-base derived from the GitHub
# event's base/head (NOT the PR-controlled head history), so a PR cannot forge a
# SKIP by rewriting its own history:
#   - head contains the test script -> run it; a real failure propagates.
#   - head lacks the script, but the merge-base contains it -> the PR
#     deleted/tampered it -> FAIL (no bypass).
#   - head lacks the script and the merge-base also lacks it -> the head forked
#     before #11312 -> explicit SKIP with a notice (never a faked PASS; exit 0).
#
# The merge-base presence is checked with `git ls-tree`, which reads tree
# objects (present in a blobless `filter: blob:none` checkout) and never lazy-
# loads a blob, so the tamper signal cannot be masked by a missing blob.
#
# Usage (from the workflow step, after checkout of the PR head):
#   PR_BASE_SHA=<base.sha> PR_BASE_REF=<base.ref> \
#   CI_BASE_REPOSITORY_URL=<base.repo.clone_url> \
#     scripts/ci/scoping-init-regression.sh
#
# Exit codes:
#   0  regression ran and passed, or legitimately skipped (pre-#11312 head)
#   1  regression failed, or the PR deleted/tampered the script (no bypass)
#   2  misconfiguration / unable to verify the boundary

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
init_test="$repo_root/scripts/ci/initialization-check_test.sh"

die() {
        echo "scoping-init-regression: $*" >&2
        exit 2
}

# Case 1: the test script is present in the checked-out head -> run it for real.
if [[ -x "$init_test" ]]; then
        exec "$init_test"
fi
# A present-but-non-executable file is suspicious (tamper); do not silently skip.
if [[ -e "$init_test" ]]; then
        echo "::error::scripts/ci/initialization-check_test.sh exists in the PR head but is not executable; cannot run the regression (possible tamper)." >&2
        exit 1
fi

# Case 2: the script is absent from the head. Decide via the trusted merge-base.
base_sha="${PR_BASE_SHA:-}"
base_ref="${PR_BASE_REF:-}"
base_url="${CI_BASE_REPOSITORY_URL:-}"

[[ "$base_sha" =~ ^[0-9a-f]{40}$ ]] ||
        die "PR_BASE_SHA must be a full 40-char commit SHA (got: ${base_sha:-<empty>})"
[[ -n "$base_ref" ]] || die "PR_BASE_REF is required"
[[ -n "$base_url" && "$base_url" != -* ]] ||
        die "CI_BASE_REPOSITORY_URL is invalid"
git check-ref-format "refs/heads/$base_ref" >/dev/null 2>&1 ||
        die "base ref is invalid: $base_ref"

if [[ "$(git rev-parse --is-shallow-repository)" == "true" ]]; then
        die "pull request checkout is shallow; use fetch-depth: 0"
fi

# Register the trusted base URL as a promisor and fetch its complete commit
# graph (blobless) so the PR base commit is resolvable and base-only blobs can
# lazy-load later. This mirrors scripts/ci/resolve-change-range.sh; the later
# "Resolve change range" step re-runs the same idempotent config/fetch.
git config --local --replace-all remote.ci-base.url "$base_url"
git config --local remote.ci-base.promisor true
git config --local remote.ci-base.partialclonefilter blob:none
fetched_base_ref="refs/remotes/ci-base/base"
git fetch --no-tags --no-recurse-submodules --force --filter=blob:none \
        ci-base "+refs/heads/$base_ref:$fetched_base_ref" >/dev/null

git cat-file -e "${base_sha}^{commit}" 2>/dev/null ||
        die "PR base commit is unavailable after fetching the base branch"
git merge-base --is-ancestor "$base_sha" "$fetched_base_ref" ||
        die "PR base commit is not part of the fetched base branch"

mb="$(git merge-base "$base_sha" HEAD)" ||
        die "PR base and head do not have a common ancestor"

# `git ls-tree` reads the merge-base tree only (present in a blobless checkout);
# a present entry means the script existed at the divergence point.
if git ls-tree "$mb" -- scripts/ci/initialization-check_test.sh | grep -q .; then
        echo "::error::scripts/ci/initialization-check_test.sh exists at the trusted merge-base ($mb) but is absent from the PR head; deletion or tampering cannot bypass the regression." >&2
        exit 1
fi

# The script is absent from both the head and the merge-base: the head forked
# before #11312 introduced it. Explicit SKIP — not a faked PASS.
echo "::notice::scripts/ci/initialization-check_test.sh is absent from the PR head and its trusted merge-base ($mb); the head predates #11312. Skipping the initialization-check regression (explicit SKIP, not a faked PASS)."
echo "SKIP - initialization-check_test.sh: head predates #11312 (absent at trusted merge-base $mb)"
exit 0
