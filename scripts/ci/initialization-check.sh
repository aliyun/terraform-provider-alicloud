#!/usr/bin/env bash

# Cloudauth no-tests initialization check.
#
# Compiles, links and initializes the standard alicloud/service/cloudauth test
# binary WITHOUT running any acceptance test, by running:
#   go test ./alicloud/service/cloudauth -run '^$' -count=1 -v
#
# `-run '^$'` matches no Test* function, so no test body executes and no cloud
# resource is created; the test binary still compiles, links and runs its
# init() (and TestMain, if any) before reporting "no tests to run".
# `-count=1` defeats the Go test cache so a cached pass cannot mask a real
# compile/link/init failure.
#
# This is the local-ci-check equivalent of the remote `probe_init` init probe
# (see scripts/ci/production-build-checks.md). It returns INITIALIZATION_ONLY
# and countsAsAccTest=false: it proves ONE service test binary builds+inits,
# NOT the root main build (use `make provider-build`) and NOT a real ACC (use
# the cloud acceptance suite). It exists so `ci-check-quick` no longer leaves
# the init probe to a separate, repeatedly-respawned task.
#
# A compile, link or init failure (non-zero exit, panic, build-failed) makes
# this check AND the enclosing quick/CI run fail. A cached result is never
# accepted as a pass: `-count=1` forces the real binary execution, and the
# pass decision additionally requires a genuine "no test ran" signal (no
# `=== RUN` / `--- PASS` / `--- FAIL` / `--- SKIP` lines) so a stub or a
# misconfigured filter that accidentally runs a test is rejected instead of
# being faked as PASS.
#
# Usage:
#   scripts/ci/initialization-check.sh
#   INIT_CHECK_PACKAGE=./alicloud/service/cloudauth scripts/ci/initialization-check.sh
#   INIT_CHECK_GO_TEST_CMD='<cmd>' scripts/ci/initialization-check.sh
#       # fixed command (fixture tests) — mirrors PROD_DEPS_INPUT so the
#       # pass/fail decision logic runs without a Go toolchain
#   CI_INIT_SOURCE_SHA=<sha> CI_INIT_SOURCE_ID=<id> CI_INIT_TASK_ID=<id> \
#       scripts/ci/initialization-check.sh
#       # bind the receipt to the outer real Provider Run identifiers; never
#       # guessed — only emitted when the outer run provides them
#
# Exit codes:
#   0  test binary compiled, linked and initialized; no test ran
#   1  compile/link/init failed, or a test ran, or no positive signal
#   2  misconfiguration

set -euo pipefail

INIT_CHECK_PACKAGE="${INIT_CHECK_PACKAGE:-./alicloud/service/cloudauth}"
RECEIPT_TAG="::ci:initialization-only::"

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
project_root="$(cd "$script_dir/../.." && pwd)"

die() {
        echo "initialization-check: $*" >&2
        exit 2
}

# Always print the exact command and the original stdout/stderr verbatim, so a
# failure is diagnosable from the raw output and the precise invocation — never
# replaced by go build / vet / a static scan summary.
echo "▶ Initialization check: ${INIT_CHECK_PACKAGE}"
if [ -z "${INIT_CHECK_GO_TEST_CMD:-}" ]; then
        GO_TEST_DISPLAY="go test ${INIT_CHECK_PACKAGE} -run '^\$' -count=1 -v"
else
        GO_TEST_DISPLAY="${INIT_CHECK_GO_TEST_CMD}"
fi
echo "Command: ${GO_TEST_DISPLAY}"

cd "$project_root"

if [ -z "${INIT_CHECK_GO_TEST_CMD:-}" ]; then
        set +e
        output="$(go test "${INIT_CHECK_PACKAGE}" -run '^$' -count=1 -v 2>&1)"
        rc=$?
        set -e
else
        set +e
        output="$(eval "${INIT_CHECK_GO_TEST_CMD}" 2>&1)"
        rc=$?
        set -e
fi

# Original stdout/stderr, verbatim.
printf '%s\n' "$output"

emit_receipt() { # status
        local status="$1"
        local line="${RECEIPT_TAG} status=${status} counts_as_acc_test=false package=${INIT_CHECK_PACKAGE}"
        # Bind to the outer real Provider Run identifiers ONLY when the outer
        # run provides them. Never guessed, never fabricated.
        [ -n "${CI_INIT_SOURCE_SHA:-}" ] && line="${line} source_sha=${CI_INIT_SOURCE_SHA}"
        [ -n "${CI_INIT_SOURCE_ID:-}" ]  && line="${line} source_id=${CI_INIT_SOURCE_ID}"
        [ -n "${CI_INIT_TASK_ID:-}" ]    && line="${line} task_id=${CI_INIT_TASK_ID}"
        printf '%s\n' "$line"
}

# Fail fast: a non-zero exit is a compile, link, init or runtime failure.
if [ "$rc" -ne 0 ]; then
        emit_receipt fail
        echo "✗ FAILED: initialization check exited ${rc}" >&2
        exit 1
fi

# A pass requires a genuine no-tests signal: no test body may have executed.
if printf '%s\n' "$output" | grep -qE '^(=== RUN|--- PASS|--- FAIL|--- SKIP)'; then
        emit_receipt fail
        echo "✗ FAILED: a test body executed (=== RUN / --- PASS / --- FAIL / --- SKIP present); expected no tests to run" >&2
        exit 1
fi

# Reject "no test files" — the test binary was never built, so initialization
# was not actually exercised.
if printf '%s\n' "$output" | grep -qE '^\?[[:space:]]'; then
        emit_receipt fail
        echo "✗ FAILED: no test files in package; test binary not built" >&2
        exit 1
fi

# Positive signal: the package compiled + linked + initialized and reported ok.
if ! printf '%s\n' "$output" | grep -qE '^(ok[[:space:]]|PASS[[:space:]]*$)'; then
        emit_receipt fail
        echo "✗ FAILED: no ok/PASS signal from the test binary" >&2
        exit 1
fi

emit_receipt pass
echo "✓ PASSED: initialization check (no tests to run, exit 0)"
exit 0
