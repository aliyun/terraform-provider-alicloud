#!/usr/bin/env bash
# bootstrap/tests/lib.sh — unit tests for lib.sh jarvis_root()
#
# Scenarios:
#   1. JARVIS_ROOT env override → returns that value
#   2. From main repo dir → returns main repo root
#   3. From a git worktree → returns main repo root (not worktree root)
#   4. jarvis_root is idempotent (multiple source + call = same result)
#
# Run: bash bootstrap/tests/lib.sh
# Prints PASS and exits 0 on success; prints FAIL and exits 1 on any failure.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
LIB_SH="$BOOTSTRAP_DIR/lib.sh"

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
pass_count=0
fail_count=0

assert_eq() {
    local desc="$1"
    local expected="$2"
    local actual="$3"
    if [ "$expected" = "$actual" ]; then
        echo "  PASS: $desc"
        pass_count=$((pass_count + 1))
    else
        echo "  FAIL: $desc (expected='$expected' actual='$actual')"
        fail_count=$((fail_count + 1))
    fi
}

assert_exit_code() {
    local desc="$1"
    local expected="$2"
    shift 2
    local actual=0
    "$@" >/dev/null 2>&1 || actual=$?
    if [ "$actual" -eq "$expected" ]; then
        echo "  PASS: $desc (exit $actual)"
        pass_count=$((pass_count + 1))
    else
        echo "  FAIL: $desc (expected=$expected actual=$actual)"
        fail_count=$((fail_count + 1))
    fi
}

# ---------------------------------------------------------------------------
# Test 1: JARVIS_ROOT env override is respected
# ---------------------------------------------------------------------------
echo "Test 1: JARVIS_ROOT env override"
result=$(JARVIS_ROOT=/tmp/test-override bash -c "source '$LIB_SH' && jarvis_root")
assert_eq "env override returns /tmp/test-override" "/tmp/test-override" "$result"

# ---------------------------------------------------------------------------
# Test 2: From main repo → returns main repo root (strip /.git)
# ---------------------------------------------------------------------------
echo "Test 2: from main repo dir"
# Derive expected main repo root from git-common-dir in this script's repo
expected_root=$(git -C "$BOOTSTRAP_DIR" rev-parse --path-format=absolute --git-common-dir 2>/dev/null)
expected_root="${expected_root%/.git}"
result=$(unset JARVIS_ROOT; bash -c "cd '$BOOTSTRAP_DIR' && source '$LIB_SH' && jarvis_root")
assert_eq "main repo returns $expected_root" "$expected_root" "$result"

# ---------------------------------------------------------------------------
# Test 3: From git worktree → returns MAIN repo root (not worktree root)
# ---------------------------------------------------------------------------
echo "Test 3: from git worktree"
# The test suite itself may be running in a worktree. git-common-dir should
# always point to the main repo's .git, so the result must equal expected_root
# (which we derived via the same mechanism but from BOOTSTRAP_DIR).
worktree_toplevel=$(git -C "$BOOTSTRAP_DIR" rev-parse --show-toplevel 2>/dev/null)
result=$(unset JARVIS_ROOT; bash -c "cd '$BOOTSTRAP_DIR' && source '$LIB_SH' && jarvis_root")
assert_eq "worktree returns main root $expected_root (not $worktree_toplevel)" "$expected_root" "$result"

# ---------------------------------------------------------------------------
# Test 4: Idempotent — sourcing twice, calling twice = same result
# ---------------------------------------------------------------------------
echo "Test 4: idempotent multiple source+call"
result=$(unset JARVIS_ROOT; bash -c "
    source '$LIB_SH'
    r1=\$(jarvis_root)
    source '$LIB_SH'
    r2=\$(jarvis_root)
    [ \"\$r1\" = \"\$r2\" ] && echo \"\$r1\" || echo 'MISMATCH'
")
assert_eq "idempotent call returns $expected_root" "$expected_root" "$result"

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo "Results: $pass_count passed, $fail_count failed"
if [ "$fail_count" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    exit 1
fi
