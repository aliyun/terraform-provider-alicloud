#!/bin/bash
# Test harness for bootstrap/log.sh
# Tests: seen, run_done, escalate

set -u

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

PASS=0
FAIL=0

assert_pass() {
    echo "PASS: $1"
    PASS=$((PASS + 1))
}

assert_fail() {
    echo "FAIL: $1"
    FAIL=$((FAIL + 1))
}

# ---------------------------------------------------------------------------
# Use a temp dir so tests are isolated and repo is not littered
# ---------------------------------------------------------------------------
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

# Redirect runs/ and escalation/ to temp dir by overriding the path via env
export JARVIS_RUNS_DIR="$tmpdir/runs"
export JARVIS_ESCALATION_DIR="$tmpdir/escalation"

mkdir -p "$JARVIS_RUNS_DIR"
mkdir -p "$JARVIS_ESCALATION_DIR"

# Source the library
# shellcheck source=../bootstrap/log.sh
# We need to override repo_root so log.sh uses our temp dirs
# log.sh reads JARVIS_RUNS_DIR and JARVIS_ESCALATION_DIR env vars if set

source "$proj_root/bootstrap/log.sh"

# ---------------------------------------------------------------------------
# Test 1: seen returns non-zero for unknown id
# ---------------------------------------------------------------------------
echo "=== Test 1: seen returns non-zero for unknown id ==="
if seen "TASK-UNKNOWN-9999"; then
    assert_fail "seen should return non-zero for unknown id, got 0"
else
    assert_pass "seen returns non-zero for unknown id"
fi

# ---------------------------------------------------------------------------
# Test 2: run_done writes a file in runs/
# ---------------------------------------------------------------------------
echo "=== Test 2: run_done writes runs/ file ==="
run_done "TASK-001" "processed successfully"

runs_file_count=$(ls "$JARVIS_RUNS_DIR"/*-TASK-001.md 2>/dev/null | wc -l | tr -d ' ')
if [ "$runs_file_count" -ge 1 ]; then
    assert_pass "run_done created runs/ file for TASK-001"
else
    assert_fail "run_done did not create runs/ file for TASK-001 (found $runs_file_count files)"
fi

# ---------------------------------------------------------------------------
# Test 3: after run_done, seen returns 0
# ---------------------------------------------------------------------------
echo "=== Test 3: seen returns 0 after run_done ==="
if seen "TASK-001"; then
    assert_pass "seen returns 0 after run_done"
else
    assert_fail "seen should return 0 after run_done, got non-zero"
fi

# ---------------------------------------------------------------------------
# Test 4: escalate writes escalation/<id>.md
# ---------------------------------------------------------------------------
echo "=== Test 4: escalate writes escalation file ==="
escalate "TASK-002" "low priority, needs manual review"

if [ -f "$JARVIS_ESCALATION_DIR/TASK-002.md" ]; then
    assert_pass "escalate created escalation/TASK-002.md"
else
    assert_fail "escalate did not create escalation/TASK-002.md"
fi

# ---------------------------------------------------------------------------
# Test 5: runs/ file contains the id and summary
# ---------------------------------------------------------------------------
echo "=== Test 5: runs/ file content check ==="
runs_file=$(ls "$JARVIS_RUNS_DIR"/*-TASK-001.md 2>/dev/null | head -1)
if grep -q "TASK-001" "$runs_file" 2>/dev/null; then
    assert_pass "runs/ file contains id"
else
    assert_fail "runs/ file does not contain id"
fi

if grep -q "processed successfully" "$runs_file" 2>/dev/null; then
    assert_pass "runs/ file contains summary"
else
    assert_fail "runs/ file does not contain summary"
fi

# ---------------------------------------------------------------------------
# Test 6: escalation file contains id and reason
# ---------------------------------------------------------------------------
echo "=== Test 6: escalation file content check ==="
if grep -q "TASK-002" "$JARVIS_ESCALATION_DIR/TASK-002.md" 2>/dev/null; then
    assert_pass "escalation file contains id"
else
    assert_fail "escalation file does not contain id"
fi

if grep -q "low priority" "$JARVIS_ESCALATION_DIR/TASK-002.md" 2>/dev/null; then
    assert_pass "escalation file contains reason"
else
    assert_fail "escalation file does not contain reason"
fi

# ---------------------------------------------------------------------------
# Test 7: runs/ filename uses date format YYYY-MM-DD
# ---------------------------------------------------------------------------
echo "=== Test 7: runs/ filename date format ==="
runs_file=$(ls "$JARVIS_RUNS_DIR"/*-TASK-001.md 2>/dev/null | head -1)
basename_file=$(basename "$runs_file" 2>/dev/null)
if echo "$basename_file" | grep -qE '^[0-9]{4}-[0-9]{2}-[0-9]{2}-TASK-001\.md$'; then
    assert_pass "runs/ filename has correct date format"
else
    assert_fail "runs/ filename does not match expected pattern, got: $basename_file"
fi

# ---------------------------------------------------------------------------
# Test 8: direct dispatch — bash log.sh seen ID (non-zero for unknown)
# ---------------------------------------------------------------------------
echo "=== Test 8: direct dispatch seen (non-zero for unknown) ==="
if JARVIS_RUNS_DIR="$JARVIS_RUNS_DIR" JARVIS_ESCALATION_DIR="$JARVIS_ESCALATION_DIR" \
   bash "$proj_root/bootstrap/log.sh" seen "TASK-NEVER-SEEN-XYZ"; then
    assert_fail "direct dispatch seen should return non-zero for unknown id"
else
    assert_pass "direct dispatch seen returns non-zero for unknown id"
fi

# ---------------------------------------------------------------------------
# Test 9: direct dispatch — bash log.sh run_done then seen
# ---------------------------------------------------------------------------
echo "=== Test 9: direct dispatch run_done then seen ==="
JARVIS_RUNS_DIR="$JARVIS_RUNS_DIR" JARVIS_ESCALATION_DIR="$JARVIS_ESCALATION_DIR" \
    bash "$proj_root/bootstrap/log.sh" run_done "TASK-003" "dispatched ok"

if JARVIS_RUNS_DIR="$JARVIS_RUNS_DIR" JARVIS_ESCALATION_DIR="$JARVIS_ESCALATION_DIR" \
   bash "$proj_root/bootstrap/log.sh" seen "TASK-003"; then
    assert_pass "direct dispatch: seen returns 0 after run_done TASK-003"
else
    assert_fail "direct dispatch: seen should return 0 after run_done TASK-003"
fi

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo "=== Summary ==="
echo "PASS: $PASS  FAIL: $FAIL"
echo ""

if [ "$FAIL" -gt 0 ]; then
    echo "TESTS FAILED"
    exit 1
fi

echo "All tests passed"
exit 0
