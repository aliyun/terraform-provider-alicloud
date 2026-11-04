#!/usr/bin/env bash

# Fixture tests for scripts/ci/initialization-check.sh. Each case feeds a stub
# `go test` command via INIT_CHECK_GO_TEST_CMD (mirroring PROD_DEPS_INPUT) so
# the pass/fail decision logic runs without a Go toolchain, plus source-level
# assertions on the cache flags, single-probe wiring and receipt format.
#
# The genuine `go test` execution is validated separately by a real
# `ci-check-quick` run (remote quick / `make ci-check-quick`); these fixtures
# prove the wrapper rejects compile failure, panic, an accidentally-run test,
# a no-test-files package and a bare ok/PASS without a genuine no-tests
# signal, and accepts only a real no-tests `ok` line carrying the
# `[no tests to run]` marker (plus the `testing: warning: no tests to run`
# line).

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
check="$script_dir/initialization-check.sh"
local_ci="$script_dir/../local-ci-check.sh"

die() {
        echo "FAIL: $*" >&2
        exit 1
}

# run_case <name> <stub-cmd> <expected-exit> <expected-substring>
run_case() {
        local name="$1" stub="$2" expect_exit="$3" expect_match="$4"
        set +e
        output="$(INIT_CHECK_PACKAGE=./a/pkg INIT_CHECK_GO_TEST_CMD="$stub" "$check" 2>&1)"
        rc=$?
        set -e
        [[ "$rc" -eq "$expect_exit" ]] ||
                die "$name: expected exit $expect_exit, got $rc; output: $output"
        [[ -z "$expect_match" || "$output" == *"$expect_match"* ]] ||
                die "$name: expected output to contain '$expect_match'; got: $output"
        echo "ok - $name"
}

# 1. Genuine no-tests pass: the real go-test output shape — a `testing:
#    warning: no tests to run` line, a bare PASS line, and an `ok` line
#    carrying the `[no tests to run]` marker. exit 0.
run_case "no-tests pass" \
        $'printf \'testing: warning: no tests to run\\nPASS\\nok\\t./a/pkg\\t0.001s\\t[no tests to run]\\n\'; exit 0' \
        0 "no tests to run"

# 2. Compile failure propagates: original stderr printed, check fails.
run_case "compile failure fails" \
        $'printf \'# ./a/pkg\\n./x_test.go:1:1: expected declaration\\nFAIL\\t./a/pkg [build failed]\\n\' >&2; exit 1' \
        1 "expected declaration"

# 3. init() panic propagates as a failure.
run_case "init panic fails" \
        $'printf \'panic: runtime error: index out of range\\nFAIL\\t./a/pkg [exited]\\n\' >&2; exit 1' \
        1 "panic"

# 4. A test body executed (=== RUN present) is rejected even with exit 0:
#    success must genuinely come from a no-tests run, never faked as PASS.
run_case "ran-a-test rejected" \
        $'printf \'=== RUN   TestFake\\n--- PASS: TestFake (0.00s)\\nPASS\\nok\\t./a/pkg\\t0.001s\\n\'; exit 0' \
        1 "a test body executed"

# 5. "no test files" means the binary was never built — initialization not
#    exercised, so fail rather than fake a pass.
run_case "no-test-files rejected" \
        $'printf \'?\\t./a/pkg\\t[no test files]\\n\'; exit 0' \
        1 "no test files"

# 6. No positive ok/PASS signal (exit 0 but garbage) is rejected.
run_case "no-ok-signal rejected" \
        $'printf \'random output\\n\'; exit 0' \
        1 "no ok/PASS signal"

# 7. Bare ok/PASS WITHOUT a genuine no-tests marker is rejected: a stub can
#    print "ok"/"PASS" without ever running the test binary, so the wrapper
#    must not report a pass from a bare ok/PASS line alone. This is the gap
#    the tightened judgment closes (an ok/PASS line is necessary but not
#    sufficient; the real `[no tests to run]` / testing warning signal is
#    required).
run_case "bare ok/pass without no-tests signal rejected" \
        $'printf \'PASS\\nok\\t./a/pkg\\t0.001s\\n\'; exit 0' \
        1 "no genuine no-tests signal"

# 8. Receipt is machine-readable and marks INITIALIZATION_ONLY / not an ACC.
set +e
receipt="$(INIT_CHECK_PACKAGE=./a/pkg \
        INIT_CHECK_GO_TEST_CMD=$'printf \'testing: warning: no tests to run\\nPASS\\nok\\t./a/pkg\\t0.001s\\t[no tests to run]\\n\'; exit 0' \
        "$check" 2>&1)"
rc=$?
set -e
[[ "$rc" -eq 0 ]] || die "receipt case: expected exit 0, got $rc; output: $receipt"
echo "$receipt" | grep -qE '^::ci:initialization-only:: status=pass counts_as_acc_test=false package=\./a/pkg$' ||
        die "receipt format wrong; got: $receipt"
# No outer Provider Run identifiers were provided, so none are guessed.
echo "$receipt" | grep -q ' source_id=' && die "receipt guessed a source_id"
echo "$receipt" | grep -q ' task_id=' && die "receipt guessed a task_id"
echo "$receipt" | grep -q ' source_sha=' && die "receipt guessed a source_sha"
echo "ok - receipt format (no guessed identifiers)"

# 9. Outer Provider Run identifiers are bound only when provided (never guessed).
set +e
receipt="$(INIT_CHECK_PACKAGE=./a/pkg \
        CI_INIT_SOURCE_SHA=abc123 CI_INIT_SOURCE_ID=sid42 CI_INIT_TASK_ID=tid7 \
        INIT_CHECK_GO_TEST_CMD=$'printf \'testing: warning: no tests to run\\nPASS\\nok\\t./a/pkg\\t0.001s\\t[no tests to run]\\n\'; exit 0' \
        "$check" 2>&1)"
rc=$?
set -e
[[ "$rc" -eq 0 ]] || die "bind case: expected exit 0, got $rc; output: $receipt"
echo "$receipt" | grep -q 'source_sha=abc123' || die "receipt missing source_sha"
echo "$receipt" | grep -q 'source_id=sid42' || die "receipt missing source_id"
echo "$receipt" | grep -q 'task_id=tid7' || die "receipt missing task_id"
echo "ok - receipt binds outer identifiers when provided"

# 10. The default command defeats the cache (-count=1) and matches no test
#     (-run '^$') with verbose output (-v), so a cached pass cannot skip the
#     real test binary.
grep -q -- "-count=1" "$check" || die "default command missing -count=1"
grep -q -- "-run '^\$'" "$check" || die "default command missing -run '^\$'"
grep -q -- "-v" "$check" || die "default command missing -v"
echo "ok - default command defeats cache and matches no test"

# 11. The probe runs exactly once per quick, and --quick does not bypass it:
#     local-ci-check.sh invokes initialization-check.sh exactly once (count
#     actual command invocations, not doc mentions), outside any
#     SKIP_BUILD/CHANGED_UNIT_TESTS gate that quick would skip.
invocations="$(grep -nE 'ci/initialization-check\.sh' "$local_ci" | grep -vE '^[0-9]+:[[:space:]]*#' | wc -l | tr -d ' ')"
[[ "$invocations" -eq 1 ]] ||
        die "expected exactly one initialization-check.sh invocation in local-ci-check.sh, got $invocations"
# The call must not live inside the SKIP_BUILD build-check block (quick sets
# --skip-build, so a call gated by SKIP_BUILD=false would be bypassed).
call_line="$(grep -nE 'ci/initialization-check\.sh' "$local_ci" | grep -vE '^[0-9]+:[[:space:]]*#' | head -1)"
call_lineno="${call_line%%:*}"
# Find the nearest preceding "if [ \"\$SKIP_BUILD\" = false ]" build gate and
# ensure the init call is not between it and its closing fi.
gate_line="$(grep -nE 'if \[ "\$SKIP_BUILD" = false \]' "$local_ci" | head -1 || true)"
if [ -n "$gate_line" ]; then
        gate_lineno="${gate_line%%:*}"
        [[ "$call_lineno" -gt "$gate_lineno" ]] ||
                die "init check call (line $call_lineno) precedes the SKIP_BUILD gate (line $gate_lineno); expected to be outside it"
fi
echo "ok - single probe, not gated by SKIP_BUILD (quick does not bypass)"

# Optional real-toolchain smoke test. Off by default: the deterministic stub
# cases above cover the decision logic; a real `go test` run is the remote
# quick run's job. Enable explicitly when a Go toolchain is available.
if [ "${CI_INIT_CHECK_REAL:-0}" = "1" ] && command -v go >/dev/null 2>&1; then
        set +e
        real_out="$(INIT_CHECK_PACKAGE=./alicloud/service/cloudauth "$check" 2>&1)"
        real_rc=$?
        set -e
        [[ "$real_rc" -eq 0 ]] ||
                die "real cloudauth init check failed (rc=$real_rc); output: $real_out"
        echo "$real_out" | grep -q 'no tests to run' ||
                die "real run did not report a no-tests pass; output: $real_out"
        echo "ok - real cloudauth no-tests initialization"
fi

echo
echo "all initialization-check fixtures passed"
