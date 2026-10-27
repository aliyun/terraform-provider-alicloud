#!/usr/bin/env bash

# Fixture tests for scripts/ci/scoping-init-regression.sh.
#
# #11312 wired scripts/ci/initialization-check_test.sh into the "Test change
# range and lint scoping" workflow step unconditionally. PR heads forked before
# #11312 lack the file, so the bare call exits 127 against the new baseline.
# These fixtures prove the boundary-gated runner:
#   - runs the real test (and propagates its exit code) when the head has it;
#   - explicit SKIP (exit 0, no exec) when head and merge-base both lack it
#     (legitimate pre-#11312 head);
#   - FAIL (exit 1, no bypass) when head lacks it but the trusted merge-base
#     has it (PR deleted/tampered);
#   - fails closed (exit 2) on misconfiguration.

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
helper="$script_dir/scoping-init-regression.sh"

die() {
        echo "FAIL: $*" >&2
        exit 1
}

commit_file() {
        local repository="$1"
        local path="$2"
        local content="$3"
        local subject="$4"
        mkdir -p "$(dirname "$repository/$path")"
        printf '%s\n' "$content" >"$repository/$path"
        git -C "$repository" add "$path"
        git -C "$repository" commit -q -m "$subject"
}

# Build a fork checkout whose PR head is `head_sha`, branched from `base_sha` on
# a trusted bare base repo, then run the helper from inside the checkout.
# $1 temp_dir, $2 base_setup (with_script|without_script), $3 head_action
# (delete_script|none), $4 extra env for the helper invocation.
run_in_fork() {
        local temp_dir="$1" base_setup="$2" head_action="$3"
        shift 3
        local extra_env=("$@")

        # Derive a unique tag from the fixture shape so concurrent/sequential
        # run_in_fork calls never reuse the same source/base.git/fork.git/checkout
        # subdirs (the prior collision left a stale 'feature' branch and existing
        # clone destinations, cascading into 'not our ref' and unbound-variable
        # failures that masked the real boundary result).
        local tag="${base_setup}-${head_action}"
        local source="$temp_dir/${tag}-source"
        local base_bare="$temp_dir/${tag}-base.git"
        local fork_bare="$temp_dir/${tag}-fork.git"
        local checkout="$temp_dir/${tag}-checkout"

        git init -q -b master "$source"
        git -C "$source" config core.hooksPath /dev/null
        git -C "$source" config user.name "CI Boundary Test"
        git -C "$source" config user.email "ci-boundary-test@example.invalid"

        # Common ancestor commit (always present, no script yet).
        commit_file "$source" alicloud/common.go 'package alicloud' common-ancestor

        if [[ "$base_setup" == "with_script" ]]; then
                commit_file "$source" scripts/ci/initialization-check_test.sh \
                        '#!/usr/bin/env bash
echo "real-init-test-ran"
exit 0' add-init-test
                chmod +x "$source/scripts/ci/initialization-check_test.sh"
                git -C "$source" add scripts/ci/initialization-check_test.sh
                git -C "$source" commit -q -m 'chmod init test'
        fi
        local base_sha
        base_sha="$(git -C "$source" rev-parse HEAD)"

        git -C "$source" switch -q -c feature
        if [[ "$head_action" == "delete_script" ]]; then
                git -C "$source" rm -q scripts/ci/initialization-check_test.sh
                git -C "$source" commit -q -m 'PR deletes init test'
        fi
        commit_file "$source" alicloud/head.go 'package alicloud' pr-head
        local head_sha
        head_sha="$(git -C "$source" rev-parse HEAD)"

        # Trusted bare base (canonical repo master) + bare fork (head origin).
        # The runner fetches ci-base with --filter=blob:none, so the bare repos
        # must advertise filter support (and SHA fetch) over the file transport.
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --bare --branch master "file://$source" "$base_bare"
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --bare --branch feature "file://$source" "$fork_bare"
        git -C "$base_bare" config uploadpack.allowFilter true
        git -C "$base_bare" config uploadpack.allowAnySHA1InWant true
        git -C "$fork_bare" config uploadpack.allowFilter true
        git -C "$fork_bare" config uploadpack.allowAnySHA1InWant true
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --branch feature "file://$fork_bare" "$checkout"
        git -C "$checkout" config protocol.file.allow always
        git -C "$checkout" checkout -q "$head_sha"

        # The helper must not see the real repo's init test; the checkout tree
        # is the PR head, so the file presence is whatever the fork provides.
        set +e
        output="$(cd "$checkout" && env "${extra_env[@]}" \
                PR_BASE_SHA="$base_sha" PR_BASE_REF=master \
                CI_BASE_REPOSITORY_URL="file://$base_bare" \
                "$helper" 2>&1)"
        rc=$?
        set -e
        printf '%s\n' "$rc"
        printf '%s\n' "$output"
}

test_present_runs_and_propagates() {
        local temp_dir="$1"
        local source="$temp_dir/present-src"
        local checkout="$temp_dir/present-checkout"

        git init -q -b master "$source"
        git -C "$source" config core.hooksPath /dev/null
        git -C "$source" config user.name "CI Boundary Test"
        git -C "$source" config user.email "ci-boundary-test@example.invalid"
        commit_file "$source" alicloud/main.go 'package alicloud' base
        # Executable stub init test whose exit code is controlled by STUB_EXIT.
        mkdir -p "$source/scripts/ci"
        printf '%s\n' '#!/usr/bin/env bash' \
                'echo "stub-init-test-ran rc=${STUB_EXIT:-0}"' \
                'exit "${STUB_EXIT:-0}"' >"$source/scripts/ci/initialization-check_test.sh"
        chmod +x "$source/scripts/ci/initialization-check_test.sh"
        git -C "$source" add scripts/ci/initialization-check_test.sh
        git -C "$source" commit -q -m 'add stub init test'
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --branch master "file://$source" "$checkout"

        local out rc
        set +e
        out="$(cd "$checkout" && STUB_EXIT=0 "$helper" 2>&1)"
        rc=$?
        set -e
        [[ "$rc" -eq 0 ]] || die "present/exit0: expected 0, got $rc; out: $out"
        [[ "$out" == *"stub-init-test-ran rc=0"* ]] ||
                die "present/exit0: helper did not exec the stub; out: $out"

        set +e
        out="$(cd "$checkout" && STUB_EXIT=1 "$helper" 2>&1)"
        rc=$?
        set -e
        [[ "$rc" -eq 1 ]] || die "present/exit1: expected 1 (propagation), got $rc; out: $out"
        [[ "$out" == *"stub-init-test-ran rc=1"* ]] ||
                die "present/exit1: helper did not exec the stub; out: $out"

        echo "ok - present head runs the real test and propagates exit code"
}

test_absent_at_head_and_merge_base_skips() {
        local temp_dir="$1" result out rc
        # Base predates #11312 (no script); head branched from it, still none.
        result="$(run_in_fork "$temp_dir" without_script none)" \
                || die "absent-skip: fixture setup failed: $result"
        rc="$(printf '%s\n' "$result" | sed -n '1p')"
        out="$(printf '%s\n' "$result" | tail -n +2)"
        [[ "$rc" =~ ^[0-9]+$ ]] ||
                die "absent-skip: fixture result was polluted (no clean rc line): $result"
        [[ "$rc" -eq 0 ]] ||
                die "absent-skip: expected exit 0, got $rc; out: $out"
        [[ "$out" == *"SKIP - initialization-check_test.sh"* ]] ||
                die "absent-skip: expected explicit SKIP line; out: $out"
        [[ "$out" == *"predates #11312"* ]] ||
                die "absent-skip: expected pre-#11312 reason; out: $out"
        [[ "$out" != *"stub-init-test-ran"* && "$out" != *"real-init-test-ran"* ]] ||
                die "absent-skip: must not run any init test; out: $out"
        echo "ok - absent at head and merge-base: explicit SKIP (exit 0), no exec"
}

test_absent_at_head_present_at_merge_base_fails() {
        local temp_dir="$1" result out rc
        # Base has the script (#11312 merged); the PR head deletes it -> tamper.
        result="$(run_in_fork "$temp_dir" with_script delete_script)" \
                || die "tamper: fixture setup failed: $result"
        rc="$(printf '%s\n' "$result" | sed -n '1p')"
        out="$(printf '%s\n' "$result" | tail -n +2)"
        [[ "$rc" =~ ^[0-9]+$ ]] ||
                die "tamper: fixture result was polluted (no clean rc line): $result"
        [[ "$rc" -eq 1 ]] ||
                die "tamper: expected exit 1 (no bypass), got $rc; out: $out"
        [[ "$out" == *"merge-base"* && "$out" == *"absent from the PR head"* ]] ||
                die "tamper: expected merge-base/tamper error; out: $out"
        [[ "$out" != *"SKIP"* ]] ||
                die "tamper: must not SKIP; out: $out"
        echo "ok - absent at head but present at merge-base: FAIL (no bypass)"
}

test_misconfig_fails_closed() {
        local temp_dir="$1" source="$temp_dir/misconfig-src"
        local checkout="$temp_dir/misconfig-checkout"

        git init -q -b master "$source"
        git -C "$source" config core.hooksPath /dev/null
        git -C "$source" config user.name "CI Boundary Test"
        git -C "$source" config user.email "ci-boundary-test@example.invalid"
        commit_file "$source" alicloud/main.go 'package alicloud' base
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --branch master "file://$source" "$checkout"

        # Script absent and no valid PR_BASE_SHA -> fail closed (exit 2).
        local out rc
        set +e
        out="$(cd "$checkout" && PR_BASE_SHA= PR_BASE_REF=master \
                CI_BASE_REPOSITORY_URL="file://$source" "$helper" 2>&1)"
        rc=$?
        set -e
        [[ "$rc" -eq 2 ]] || die "misconfig: expected exit 2, got $rc; out: $out"
        [[ "$out" == *"PR_BASE_SHA"* ]] ||
                die "misconfig: expected PR_BASE_SHA reason; out: $out"
        echo "ok - misconfiguration fails closed (exit 2)"
}

# Workflow-entry-level fixture (review inputs #1/#2/#3/#4). The workflow step
# checks out the PR HEAD, then the new merge-commit baseline materializes BOTH
# scoping-init-regression.sh and scoping-init-regression_test.sh from a trusted,
# immutable SHA (GitHub-event base = master tip, content-addressed) into a temp
# tool dir and executes both. A head that predates this PR lacks those entries,
# so a bare call would itself exit 127. The materialized runner SKIPs on such a
# head (boundary-checked against the trusted merge-base). But the materialized
# TEST must also execute successfully on that old-head CWD: it must NOT assume
# the checked-out HEAD is the tool source -- on a real pre-#11312 head, HEAD
# lacks both the new runner and the init-test add history, so `git log` from HEAD
# and `git show HEAD:...` would fail (review input #4 reproduction). The workflow
# passes the trusted tool-source SHA as SCOPING_TOOL_SOURCE_SHA; the test uses
# that SHA (not HEAD) to locate the new entry and the init-test add commit. This
# fixture:
#   - materializes BOTH entries from the trusted tool source;
#   - runs the runner on a REAL pre-#11312 head CWD -> SKIP (no 127);
#   - runs the materialized TEST on that same old-head CWD -> PASS (exit 0),
#     proving the COMPLETE workflow step (both entries) succeeds on an old head,
#     not just the runner alone (a depth guard bounds re-entry);
#   - still FAILs a head that deletes the init test and a head whose init test
#     fails for real.
test_workflow_entry_bootstrap_against_old_head() {
        local temp_dir="$1"
        # Depth guard: when this fixture runs inside a materialized _test that is
        # already executing on an old-head CWD (the very scenario under test),
        # re-running the _test would recurse unboundedly. At depth>0 we still run
        # the runner (proving SKIP is stable from the nested context) but do not
        # re-invoke the _test.
        local depth="${SCOPING_ENTRY_DEPTH:-0}"
        # The test is materialized into a temp tool dir (not a git repo); use the
        # workflow step's CWD (the checked-out head repo root), not $script_dir,
        # for git operations against the real repository.
        local repo_root
        repo_root="$(git rev-parse --show-toplevel)"

        # Trusted immutable SHA carrying this PR's new entry. The workflow passes
        # it as SCOPING_TOOL_SOURCE_SHA (the same SHA used to materialize the
        # runner+test). Do NOT assume current HEAD is the tool source: on a real
        # pre-#11312 head, HEAD lacks the entry and the init-test add history, so
        # `git log` from HEAD and `git show HEAD:...` would fail. Default to HEAD
        # only for local dev runs (where HEAD carries the entry) and verify the
        # resolved SHA actually contains the entry.
        local tool_source_sha="${SCOPING_TOOL_SOURCE_SHA:-$(git -C "$repo_root" rev-parse HEAD)}"
        [[ "$tool_source_sha" =~ ^[0-9a-f]{40}$ ]] \
                || die "entry: tool source SHA must be 40 hex chars (got: ${tool_source_sha:-<empty>})"
        git -C "$repo_root" cat-file -e "${tool_source_sha}:scripts/ci/scoping-init-regression.sh" 2>/dev/null \
                || die "entry: tool source SHA ${tool_source_sha} does not carry scripts/ci/scoping-init-regression.sh; set SCOPING_TOOL_SOURCE_SHA to the trusted base"

        # Real pre-#11312 commit: parent of the commit that added the init test,
        # located in the TRUSTED TOOL SOURCE history (not current HEAD, which on
        # a real old head predates #11312 and has no such add commit -> empty ->
        # `${add_commit}^` -> fatal ambiguous '^' -> die).
        local add_commit old_head_sha
        add_commit="$(git -C "$repo_root" log --diff-filter=A --format=%H -1 \
                "$tool_source_sha" -- scripts/ci/initialization-check_test.sh)" \
                || die "entry: cannot locate the init-test add commit in the tool source history"
        [[ -n "$add_commit" ]] \
                || die "entry: init-test add commit is empty in the tool source history (tool source predates #11312?)"
        old_head_sha="$(git -C "$repo_root" rev-parse "${add_commit}^")" \
                || die "entry: cannot resolve pre-#11312 parent commit"
        git -C "$repo_root" ls-tree "$old_head_sha" -- scripts/ci/initialization-check_test.sh \
                | grep -q . && die "entry: old head must lack the init test"
        git -C "$repo_root" ls-tree "$old_head_sha" -- scripts/ci/scoping-init-regression.sh \
                | grep -q . && die "entry: old head must lack the scoping runner"

        # Trusted bare base whose master carries the new entry (stand-in for
        # post-merge master). Blobless so the full object download is avoided;
        # the boundary path only needs commits+trees. The runner fetches ci-base
        # with --filter=blob:none, so advertise filter support over file://.
        # Clone from the stable main repo path (SCOPING_ENTRY_MAIN_REPO), not the
        # CWD: at depth>0 the CWD is a worktree whose file:// clone path is not a
        # transport-safe repository root.
        local main_repo="${SCOPING_ENTRY_MAIN_REPO:-$repo_root}"
        # $main_repo is a blobless partial checkout (actions/checkout filter:
        # blob:none, fetch-depth 0). Cloning it over file:// with --filter=blob:none
        # needs the source to advertise filter support (uploadpack.allowFilter),
        # else the filter is ignored and the full clone demands blobs the partial
        # checkout must lazy-fetch from its promisor (origin) -> "could not fetch
        # ... from promisor remote" -> exit 128. Enable filter advertisement on the
        # source for this fixture (local config on the ephemeral CI checkout);
        # the blobless clone then takes only commits+trees, which the partial
        # checkout holds locally.
        git -C "$main_repo" config uploadpack.allowFilter true
        git -C "$main_repo" config uploadpack.allowAnySHA1InWant true
        local trusted_bare="$temp_dir/entry-trusted.git"
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --bare --filter=blob:none "file://$main_repo" "$trusted_bare"
        git -C "$trusted_bare" branch -f master "$tool_source_sha"
        git -C "$trusted_bare" config uploadpack.allowFilter true
        git -C "$trusted_bare" config uploadpack.allowAnySHA1InWant true
        # Break the promisor chain: the blobless clone's origin remote points
        # back at the partial checkout (promisor -> origin/GitHub). When the
        # runner later fetches ci-base=file://$trusted_bare, an upload-pack walk
        # that touched a missing blob would lazy-fetch up that chain and fail
        # ("could not fetch ... from promisor remote"). The runner fetch is
        # blobless (commits+trees only) and trusted_bare already holds them, so
        # detach the promisor remote: any genuinely missing object now fails
        # closed locally instead of triggering the unreachable chain.
        git -C "$trusted_bare" remote remove origin 2>/dev/null || true

        # Materialize BOTH the runner and its test from the trusted tool source
        # SHA into a temp tool dir -- the exact bootstrap the workflow step
        # performs -- NOT a bare call from the head (which would exit 127 against
        # an old head).
        local entry_tool="$temp_dir/entry-tool"
        mkdir -p "$entry_tool"
        git -C "$repo_root" show "${tool_source_sha}:scripts/ci/scoping-init-regression.sh" \
                >"$entry_tool/scoping-init-regression.sh" 2>/dev/null \
                || die "entry: cannot materialize runner from tool source"
        git -C "$repo_root" show "${tool_source_sha}:scripts/ci/scoping-init-regression_test.sh" \
                >"$entry_tool/scoping-init-regression_test.sh" 2>/dev/null \
                || die "entry: cannot materialize test from tool source"
        chmod +x "$entry_tool/scoping-init-regression.sh" "$entry_tool/scoping-init-regression_test.sh"
        local entry_runner="$entry_tool/scoping-init-regression.sh"
        local entry_test="$entry_tool/scoping-init-regression_test.sh"

        # Head under test: a worktree at the real pre-#11312 commit. --no-checkout
        # is enough: the old head lacks the init test, so case 1 is skipped and
        # the boundary path uses git objects (commits/trees), not working-tree
        # files.
        local head_wt="$temp_dir/entry-old-head"
        git -C "$repo_root" worktree add -q --no-checkout --detach "$head_wt" "$old_head_sha" \
                || die "entry: cannot create worktree at the old head"

        local out rc

        # Sub-case 1: old head lacks the entry -> materialized runner SKIPs
        # (exit 0), NO exit 127 / "command not found".
        set +e
        out="$(cd "$head_wt" && PR_BASE_SHA="$tool_source_sha" PR_BASE_REF=master \
                CI_BASE_REPOSITORY_URL="file://$trusted_bare" \
                SCOPING_TOOL_SOURCE_SHA="$tool_source_sha" \
                "$entry_runner" 2>&1)"
        rc=$?
        set -e
        [[ "$rc" -eq 0 ]] ||
                die "entry/runner: expected SKIP (exit 0), got $rc; out: $out"
        [[ "$out" != *"command not found"* && "$out" != *"not found"* ]] ||
                die "entry/runner: must not exit 127 (entry not in head); out: $out"
        [[ "$out" == *"SKIP - initialization-check_test.sh"* ]] ||
                die "entry/runner: expected explicit SKIP; out: $out"
        echo "ok - workflow entry: runner SKIPs on real pre-#11312 head (no 127)"

        if [[ "$depth" -eq 0 ]]; then
                # Sub-case 2: the materialized TEST executes on the old-head CWD
                # and PASSes (exit 0). This proves the COMPLETE workflow step
                # (both materialized entries) succeeds on a real old head -- not
                # just the runner alone (review input #3/#4). The nested _test
                # re-runs every fixture (present/absent/tamper/misconfig all
                # build synthetic repos, so they are CWD-independent) and the
                # nested entry fixture (depth 1) re-runs the runner SKIP but does
                # not re-invoke the _test, so re-entry is bounded. SCOPING_ENTRY_
                # MAIN_REPO points the nested clone at the stable main repo.
                set +e
                out="$(cd "$head_wt" && PR_BASE_SHA="$tool_source_sha" PR_BASE_REF=master \
                        CI_BASE_REPOSITORY_URL="file://$trusted_bare" \
                        SCOPING_TOOL_SOURCE_SHA="$tool_source_sha" \
                        SCOPING_ENTRY_MAIN_REPO="$main_repo" \
                        SCOPING_ENTRY_DEPTH=$((depth+1)) \
                        "$entry_test" 2>&1)"
                rc=$?
                set -e
                [[ "$rc" -eq 0 ]] ||
                        die "entry/test: materialized _test must pass on old-head CWD (got $rc); out: $out"
                echo "ok - workflow entry: materialized _test passes on real pre-#11312 head CWD"
        else
                echo "ok - workflow entry (nested depth $depth): runner SKIP re-verified, _test re-entry guarded"
        fi

        # Sub-case 3: new-version head that deletes the init test -> the
        # materialized (trusted) runner FAILs (no bypass). Synthetic repo whose
        # master has the init test, feature head deletes it.
        local t2_src="$temp_dir/entry-tamper-src" t2_base t2_head
        local t2_bare="$temp_dir/entry-tamper-base.git"
        local t2_checkout="$temp_dir/entry-tamper-checkout"
        git init -q -b master "$t2_src"
        git -C "$t2_src" config core.hooksPath /dev/null
        git -C "$t2_src" config user.name "CI Entry Test"
        git -C "$t2_src" config user.email "ci-entry-test@example.invalid"
        commit_file "$t2_src" alicloud/common.go 'package alicloud' base
        commit_file "$t2_src" scripts/ci/initialization-check_test.sh \
                '#!/usr/bin/env bash
echo "real-init-test-ran"
exit 0' add-init-test
        chmod +x "$t2_src/scripts/ci/initialization-check_test.sh"
        git -C "$t2_src" add scripts/ci/initialization-check_test.sh
        git -C "$t2_src" commit -q -m 'chmod init test'
        t2_base="$(git -C "$t2_src" rev-parse HEAD)"
        git -C "$t2_src" switch -q -c feature
        git -C "$t2_src" rm -q scripts/ci/initialization-check_test.sh
        git -C "$t2_src" commit -q -m 'PR deletes init test'
        t2_head="$(git -C "$t2_src" rev-parse HEAD)"
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --bare --branch master "file://$t2_src" "$t2_bare"
        git -C "$t2_bare" config uploadpack.allowFilter true
        git -C "$t2_bare" config uploadpack.allowAnySHA1InWant true
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --branch feature "file://$t2_src" "$t2_checkout"
        git -C "$t2_checkout" checkout -q "$t2_head"
        set +e
        out="$(cd "$t2_checkout" && PR_BASE_SHA="$t2_base" PR_BASE_REF=master \
                CI_BASE_REPOSITORY_URL="file://$t2_bare" \
                "$entry_runner" 2>&1)"
        rc=$?
        set -e
        [[ "$rc" -eq 1 ]] ||
                die "entry/tamper: expected FAIL (exit 1, no bypass), got $rc; out: $out"
        [[ "$out" == *"merge-base"* ]] ||
                die "entry/tamper: expected merge-base tamper signal; out: $out"
        echo "ok - workflow entry: new-version head deleting init test still FAILs (no bypass)"

        # Sub-case 4: head whose init test fails for real -> materialized runner
        # execs it and the failure propagates (blocks, not skipped).
        local t3_src="$temp_dir/entry-fail-src" t3_checkout="$temp_dir/entry-fail-checkout"
        git init -q -b master "$t3_src"
        git -C "$t3_src" config core.hooksPath /dev/null
        git -C "$t3_src" config user.name "CI Entry Test"
        git -C "$t3_src" config user.email "ci-entry-test@example.invalid"
        commit_file "$t3_src" alicloud/common.go 'package alicloud' base
        mkdir -p "$t3_src/scripts/ci"
        printf '%s\n' '#!/usr/bin/env bash' 'echo "failing-init-test-ran"' 'exit 1' \
                >"$t3_src/scripts/ci/initialization-check_test.sh"
        chmod +x "$t3_src/scripts/ci/initialization-check_test.sh"
        git -C "$t3_src" add scripts/ci/initialization-check_test.sh
        git -C "$t3_src" commit -q -m 'add failing init test'
        git -c protocol.file.allow=always -c core.hooksPath=/dev/null clone -q \
                --branch master "file://$t3_src" "$t3_checkout"
        set +e
        out="$(cd "$t3_checkout" && PR_BASE_SHA= PR_BASE_REF=master \
                CI_BASE_REPOSITORY_URL="file://$t3_src" \
                "$entry_runner" 2>&1)"
        rc=$?
        set -e
        [[ "$rc" -eq 1 ]] ||
                die "entry/real-failure: expected exit 1 (failure propagates), got $rc; out: $out"
        [[ "$out" == *"failing-init-test-ran"* ]] ||
                die "entry/real-failure: runner must exec the failing init test; out: $out"
        echo "ok - workflow entry: real init-test failure still blocks (propagates)"

        # Release the worktree so the EXIT trap can rm -rf the temp dir cleanly.
        git -C "$repo_root" worktree remove -q --force "$head_wt" 2>/dev/null || true
        git -C "$repo_root" worktree prune -q || true
}

main() {
        [[ -x "$helper" ]] || die "missing executable helper: $helper"

        local temp_dir
        temp_dir="$(mktemp -d)"
        TEST_TEMP_DIR="$temp_dir"
        trap 'rm -rf "$TEST_TEMP_DIR"' EXIT

        test_present_runs_and_propagates "$temp_dir"
        test_absent_at_head_and_merge_base_skips "$temp_dir"
        test_absent_at_head_present_at_merge_base_fails "$temp_dir"
        test_misconfig_fails_closed "$temp_dir"
        test_workflow_entry_bootstrap_against_old_head "$temp_dir"
        echo "PASS: scoping-init-regression boundary tests"
}

main "$@"
