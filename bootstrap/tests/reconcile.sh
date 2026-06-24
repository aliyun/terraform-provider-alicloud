#!/usr/bin/env bash
# bootstrap/tests/reconcile.sh — unit tests for reconcile.sh
#
# Scenarios:
#   1. One claimed id with a runs/ file → reconcile calls claim.sh release, prints RECONCILED: <id>
#   2. One claimed id with NO runs/ file → skip (sweep handles it), prints RECONCILED: none
#   3. No claimed items at all → prints RECONCILED: none
#
# Run: bash bootstrap/tests/reconcile.sh
# Prints PASS and exits 0 on success; prints FAIL and exits 1 on any failure.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

RECONCILE="$BOOTSTRAP_DIR/reconcile.sh"

# ---------------------------------------------------------------------------
# Helpers
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

make_jarvis_root() {
    local root
    root="$(mktemp -d)"
    mkdir -p "$root/runs" "$root/.my-day" "$root/config" "$root/escalation"
    cp "$BOOTSTRAP_DIR/../config/pools.json" "$root/config/pools.json"
    echo "$root"
}

today="$(date -u +%F)"

# ---------------------------------------------------------------------------
# Temp roots declared up front so the trap always has them
# ---------------------------------------------------------------------------
ROOT1=""
ROOT2=""
ROOT3=""
FAKE_BIN_DIR=""

cleanup() {
    rm -rf "$FAKE_BIN_DIR" "$ROOT1" "$ROOT2" "$ROOT3" 2>/dev/null || true
}
trap cleanup INT TERM EXIT

# ---------------------------------------------------------------------------
# Build shared stub bin dir
# The stubs write their invocations to <root>/stub-calls.log
# ---------------------------------------------------------------------------
FAKE_BIN_DIR="$(mktemp -d)"

# a1 stub: when called as "workitem list ... -f json" return our fake payload.
# The payload is read from the env var A1_LIST_JSON (default: []).
cat > "$FAKE_BIN_DIR/a1" <<'STUB'
#!/usr/bin/env bash
# Record every invocation
echo "a1 $*" >> "${STUB_LOG:-/dev/null}"
# Only respond to "project workitem list ... -f json"
if [[ "$*" == *"workitem list"* && "$*" == *"-f json"* ]]; then
    echo "${A1_LIST_JSON:-[]}"
    exit 0
fi
exit 0
STUB
chmod +x "$FAKE_BIN_DIR/a1"

# claim.sh stub: record invocations, always succeed.
# Injected via RECONCILE_CLAIM_CMD (not PATH) to avoid the ambient-PATH footgun.
cat > "$FAKE_BIN_DIR/claim.sh" <<'STUB'
#!/usr/bin/env bash
echo "claim.sh $*" >> "${STUB_LOG:-/dev/null}"
exit 0
STUB
chmod +x "$FAKE_BIN_DIR/claim.sh"

export PATH="$FAKE_BIN_DIR:$PATH"
FAKE_BIN="$FAKE_BIN_DIR"

# ---------------------------------------------------------------------------
# Test 1: claimed id WITH a runs/ file → RECONCILED: <id>
# ---------------------------------------------------------------------------
echo "Test 1: claimed id with runs/ file → RECONCILED: <id> + claim.sh release called"
ROOT1="$(make_jarvis_root)"
STUB_LOG1="$ROOT1/stub-calls.log"

# Fake a1 returning a single claimed work item
ITEM_ID="WI-999"
export A1_LIST_JSON="[{\"identifier\":\"$ITEM_ID\"}]"
export STUB_LOG="$STUB_LOG1"

# Create the run file so log.sh seen succeeds
touch "$ROOT1/runs/${today}-${ITEM_ID}.md"

output=$(env \
    JARVIS_ROOT="$ROOT1" \
    JARVIS_RUNS_DIR="$ROOT1/runs" \
    JARVIS_ESCALATION_DIR="$ROOT1/escalation" \
    STUB_LOG="$STUB_LOG1" \
    A1_LIST_JSON="$A1_LIST_JSON" \
    RECONCILE_CLAIM_CMD="$FAKE_BIN/claim.sh" \
    bash "$RECONCILE" 2>/dev/null)

# Check output contains RECONCILED: WI-999
if echo "$output" | grep -q "RECONCILED: $ITEM_ID"; then
    assert_pass "output contains RECONCILED: $ITEM_ID"
else
    assert_fail "output should contain RECONCILED: $ITEM_ID" "got: $output"
fi

# Check claim.sh release was recorded in stub log
if grep -q "claim.sh release $ITEM_ID" "$STUB_LOG1" 2>/dev/null; then
    assert_pass "claim.sh release $ITEM_ID was called"
else
    assert_fail "claim.sh release $ITEM_ID should have been called" "stub log: $(cat "$STUB_LOG1" 2>/dev/null)"
fi

# ---------------------------------------------------------------------------
# Test 2: claimed id with NO runs/ file → skip, RECONCILED: none
# ---------------------------------------------------------------------------
echo "Test 2: claimed id with no runs/ file → RECONCILED: none (sweep handles it)"
ROOT2="$(make_jarvis_root)"
STUB_LOG2="$ROOT2/stub-calls.log"

# Use same ITEM_ID but no run file created
export A1_LIST_JSON="[{\"identifier\":\"$ITEM_ID\"}]"

output2=$(env \
    JARVIS_ROOT="$ROOT2" \
    JARVIS_RUNS_DIR="$ROOT2/runs" \
    JARVIS_ESCALATION_DIR="$ROOT2/escalation" \
    STUB_LOG="$STUB_LOG2" \
    A1_LIST_JSON="$A1_LIST_JSON" \
    RECONCILE_CLAIM_CMD="$FAKE_BIN/claim.sh" \
    bash "$RECONCILE" 2>/dev/null)

if echo "$output2" | grep -q "RECONCILED: none"; then
    assert_pass "output contains RECONCILED: none when no run file"
else
    assert_fail "output should contain RECONCILED: none" "got: $output2"
fi

# claim.sh release should NOT have been called
if grep -q "claim.sh release" "$STUB_LOG2" 2>/dev/null; then
    assert_fail "claim.sh release should NOT be called when no run file"
else
    assert_pass "claim.sh release was correctly not called"
fi

# ---------------------------------------------------------------------------
# Test 3: no claimed items at all → RECONCILED: none
# ---------------------------------------------------------------------------
echo "Test 3: no claimed items → RECONCILED: none"
ROOT3="$(make_jarvis_root)"

export A1_LIST_JSON="[]"

output3=$(env \
    JARVIS_ROOT="$ROOT3" \
    JARVIS_RUNS_DIR="$ROOT3/runs" \
    JARVIS_ESCALATION_DIR="$ROOT3/escalation" \
    A1_LIST_JSON="$A1_LIST_JSON" \
    bash "$RECONCILE" 2>/dev/null)

if echo "$output3" | grep -q "RECONCILED: none"; then
    assert_pass "empty list → RECONCILED: none"
else
    assert_fail "empty list should produce RECONCILED: none" "got: $output3"
fi

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
