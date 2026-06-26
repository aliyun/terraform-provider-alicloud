#!/usr/bin/env bash
# bootstrap/tests/triage_one.sh — unit tests for triage-one.sh
#
# Scenarios:
#   1. Happy path: claim succeeds, wrap done succeeds → claim,done,release in order + prints DONE
#   2. Missing status arg (only 4 args) → exit 1
#   3. wrap done fails → escalate called, release NOT called, exit 1
#   4. Claim lost race (exit 1) → prints SKIP, exits 0 (no escalate, no release, no done)
#
# Run: bash bootstrap/tests/triage_one.sh
# Prints PASS and exits 0 on success; prints FAIL and exits 1 on any failure.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

TRIAGE_ONE="$BOOTSTRAP_DIR/triage-one.sh"

# ---------------------------------------------------------------------------
# Test helpers
# ---------------------------------------------------------------------------
pass_count=0
fail_count=0

assert_pass() {
    local desc="$1"
    echo "  PASS: $desc"
    pass_count=$((pass_count + 1))
}

assert_fail() {
    local desc="$1"
    local detail="${2:-}"
    echo "  FAIL: $desc${detail:+ — $detail}"
    fail_count=$((fail_count + 1))
}

# make_stub_dir: creates a temp dir with stub scripts, returns the path via stdout.
# Stubs write their invocations to a STUB_LOG env var path.
STUB_DIR=""
JARVIS_ROOT_DIR=""

cleanup() {
    rm -rf "$STUB_DIR" "$JARVIS_ROOT_DIR" 2>/dev/null || true
}
trap cleanup INT TERM EXIT

STUB_DIR="$(mktemp -d)"
JARVIS_ROOT_DIR="$(mktemp -d)"

# Minimal pools.json so log.sh / wrap.sh don't bail
mkdir -p "$JARVIS_ROOT_DIR/config" "$JARVIS_ROOT_DIR/runs" "$JARVIS_ROOT_DIR/escalation"
cp "$BOOTSTRAP_DIR/../config/pools.json" "$JARVIS_ROOT_DIR/config/pools.json"

# ---------------------------------------------------------------------------
# Build stub claim.sh — behaviour controlled by STUB_CLAIM_EXIT (default 0)
# ---------------------------------------------------------------------------
cat > "$STUB_DIR/claim.sh" <<'STUB'
#!/usr/bin/env bash
echo "claim $*" >> "${STUB_LOG:?STUB_LOG not set}"
exit "${STUB_CLAIM_EXIT:-0}"
STUB
chmod +x "$STUB_DIR/claim.sh"

# ---------------------------------------------------------------------------
# Build stub wrap.sh — behaviour controlled by STUB_WRAP_EXIT (default 0)
# ---------------------------------------------------------------------------
cat > "$STUB_DIR/wrap.sh" <<'STUB'
#!/usr/bin/env bash
echo "wrap $*" >> "${STUB_LOG:?STUB_LOG not set}"
exit "${STUB_WRAP_EXIT:-0}"
STUB
chmod +x "$STUB_DIR/wrap.sh"

# ---------------------------------------------------------------------------
# Build stub log.sh — always succeeds, records calls
# ---------------------------------------------------------------------------
cat > "$STUB_DIR/log.sh" <<'STUB'
#!/usr/bin/env bash
echo "log $*" >> "${STUB_LOG:?STUB_LOG not set}"
exit 0
STUB
chmod +x "$STUB_DIR/log.sh"

# ---------------------------------------------------------------------------
# Helper: run triage-one.sh with given env overrides + args
# ---------------------------------------------------------------------------
run_triage_one() {
    local log_file="$1"
    shift
    env \
        JARVIS_ROOT="$JARVIS_ROOT_DIR" \
        JARVIS_RUNS_DIR="$JARVIS_ROOT_DIR/runs" \
        JARVIS_ESCALATION_DIR="$JARVIS_ROOT_DIR/escalation" \
        TRIAGE_CLAIM_CMD="$STUB_DIR/claim.sh" \
        TRIAGE_WRAP_CMD="$STUB_DIR/wrap.sh" \
        TRIAGE_LOG_CMD="$STUB_DIR/log.sh" \
        STUB_LOG="$log_file" \
        "$@"
}

# ===========================================================================
# Test 1: Happy path — claim succeeds, wrap done succeeds
#   Expected: claim claim, wrap done, claim release in that order; prints DONE:<id>
# ===========================================================================
echo "Test 1: happy path → claim,done,release in order + DONE output"

LOG1="$(mktemp)"

output1=$(run_triage_one "$LOG1" \
    STUB_CLAIM_EXIT=0 \
    STUB_WRAP_EXIT=0 \
    bash "$TRIAGE_ONE" WI-001 pool-a proj-1 "all done" closed 2>/dev/null)

exit1=$?

# Check exit 0
if [ "$exit1" -eq 0 ]; then
    assert_pass "exit code is 0"
else
    assert_fail "exit code should be 0" "got $exit1"
fi

# Check DONE: WI-001 in output
if echo "$output1" | grep -q "DONE: WI-001"; then
    assert_pass "output contains DONE: WI-001"
else
    assert_fail "output should contain DONE: WI-001" "got: $output1"
fi

# Check order: claim claim, then wrap done, then claim release
claim_line=$(grep -n "^claim claim" "$LOG1" | head -1 | cut -d: -f1)
wrap_line=$(grep -n "^wrap done" "$LOG1" | head -1 | cut -d: -f1)
release_line=$(grep -n "^claim release" "$LOG1" | head -1 | cut -d: -f1)

if [ -n "$claim_line" ] && [ -n "$wrap_line" ] && [ -n "$release_line" ]; then
    if [ "$claim_line" -lt "$wrap_line" ] && [ "$wrap_line" -lt "$release_line" ]; then
        assert_pass "call order is claim→done→release"
    else
        assert_fail "call order wrong" "claim=$claim_line wrap=$wrap_line release=$release_line"
    fi
else
    assert_fail "not all three calls found in log" "log: $(cat "$LOG1")"
fi

# Verify correct args forwarded to wrap done
if grep -q "^wrap done WI-001 all done closed$" "$LOG1"; then
    assert_pass "wrap done received correct args (id summary status)"
else
    assert_fail "wrap done args mismatch" "log: $(cat "$LOG1")"
fi

rm -f "$LOG1"

# ===========================================================================
# Test 2: Missing status arg (only 4 args given) → exit 1
# ===========================================================================
echo "Test 2: missing status arg → exit 1"

LOG2="$(mktemp)"

run_triage_one "$LOG2" \
    STUB_CLAIM_EXIT=0 \
    STUB_WRAP_EXIT=0 \
    bash "$TRIAGE_ONE" WI-002 pool-a proj-1 "summary only" \
    >/dev/null 2>/dev/null && exit2=0 || exit2=$?

if [ "$exit2" -ne 0 ]; then
    assert_pass "missing arg → non-zero exit"
else
    assert_fail "missing arg should exit non-zero" "got 0"
fi

rm -f "$LOG2"

# ===========================================================================
# Test 3: wrap done fails → escalate called, release NOT called, exit 1
# ===========================================================================
echo "Test 3: wrap done fails → escalate, no release, exit 1"

LOG3="$(mktemp)"

run_triage_one "$LOG3" \
    STUB_CLAIM_EXIT=0 \
    STUB_WRAP_EXIT=1 \
    bash "$TRIAGE_ONE" WI-003 pool-a proj-1 "failed summary" closed \
    >/dev/null 2>/dev/null && exit3=0 || exit3=$?

if [ "$exit3" -ne 0 ]; then
    assert_pass "wrap failure → non-zero exit"
else
    assert_fail "wrap failure should exit non-zero" "got 0"
fi

# escalate must have been called for WI-003
if grep -q "^log escalate WI-003" "$LOG3"; then
    assert_pass "log escalate called for WI-003"
else
    assert_fail "log escalate should be called" "log: $(cat "$LOG3")"
fi

# release must NOT have been called
if grep -q "^claim release" "$LOG3"; then
    assert_fail "claim release should NOT be called when wrap fails" "log: $(cat "$LOG3")"
else
    assert_pass "claim release was correctly not called"
fi

rm -f "$LOG3"

# ===========================================================================
# Test 4: Claim lost race (claim.sh exits 1) → prints SKIP, exits 0
# ===========================================================================
echo "Test 4: claim lost race → SKIP output + exit 0"

LOG4="$(mktemp)"

output4=$(run_triage_one "$LOG4" \
    STUB_CLAIM_EXIT=1 \
    STUB_WRAP_EXIT=0 \
    bash "$TRIAGE_ONE" WI-004 pool-a proj-1 "some summary" closed 2>/dev/null)

exit4=$?

if [ "$exit4" -eq 0 ]; then
    assert_pass "lost race → exit 0"
else
    assert_fail "lost race should exit 0" "got $exit4"
fi

if echo "$output4" | grep -q "SKIP"; then
    assert_pass "lost race → prints SKIP"
else
    assert_fail "lost race should print SKIP" "got: $output4"
fi

# wrap, escalate, release must NOT be called
if grep -q "^wrap" "$LOG4" 2>/dev/null; then
    assert_fail "wrap should NOT be called on lost race" "log: $(cat "$LOG4")"
else
    assert_pass "wrap was correctly not called"
fi

if grep -q "^claim release" "$LOG4" 2>/dev/null; then
    assert_fail "release should NOT be called on lost race" "log: $(cat "$LOG4")"
else
    assert_pass "release was correctly not called"
fi

rm -f "$LOG4"

# ===========================================================================
# Test 5 (C1): triage-one registers a real COORD_ID so the checkpoint
#              owner_instance is non-empty (orphan-visible on crash).
# ===========================================================================
echo "Test 5: COORD_ID registered → checkpoint owner_instance non-empty"

LOG5="$(mktemp)"
COORD_ROOT="$(mktemp -d)"
mkdir -p "$COORD_ROOT/.my-day/instances" "$COORD_ROOT/.my-day/tasks"

output5=$(env \
    JARVIS_ROOT="$COORD_ROOT" \
    TRIAGE_CLAIM_CMD="$STUB_DIR/claim.sh" \
    TRIAGE_WRAP_CMD="$STUB_DIR/wrap.sh" \
    TRIAGE_LOG_CMD="$STUB_DIR/log.sh" \
    STUB_CLAIM_EXIT=0 \
    STUB_WRAP_EXIT=0 \
    STUB_LOG="$LOG5" \
    bash "$TRIAGE_ONE" WI-005 pool-a proj-1 "all done" closed 2>/dev/null)

exit5=$?

if [ "$exit5" -eq 0 ]; then
    assert_pass "Test5 exit code is 0"
else
    assert_fail "Test5 exit code should be 0" "got $exit5"
fi

task_file="$COORD_ROOT/.my-day/tasks/WI-005.json"
if [ -f "$task_file" ]; then
    owner=$(jq -r .owner_instance "$task_file" 2>/dev/null || echo "")
    if [ -n "$owner" ] && [ "$owner" != "null" ]; then
        assert_pass "checkpoint owner_instance is non-empty: $owner"
    else
        assert_fail "checkpoint owner_instance should be non-empty" "got: '$owner'"
    fi
else
    assert_fail "checkpoint task file should exist" "$task_file missing"
fi

rm -f "$LOG5"
rm -rf "$COORD_ROOT"

# ===========================================================================
# Summary
# ===========================================================================
echo ""
echo "Results: $pass_count passed, $fail_count failed"
if [ "$fail_count" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    exit 1
fi
