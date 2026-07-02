#!/usr/bin/env bash
# test/claim_test.sh – TDD tests for bootstrap/claim.sh
# Stubs a1 to record calls; asserts claim/release invoke the right subcommands.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

tmpbin=$(mktemp -d)
tmpconfig=$(mktemp -d)
tmplog=$(mktemp)
trap 'rm -rf "$tmpbin" "$tmpconfig" "$tmplog"' EXIT

PASS=0
FAIL=0

assert_pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
assert_fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# ---------------------------------------------------------------------------
# Build tmp pools.json with claim section
# ---------------------------------------------------------------------------
mkdir -p "$tmpconfig/config"
cat > "$tmpconfig/config/pools.json" << 'JSON'
{
  "claim": { "tag": "jarvis-claimed", "done_tag": "jarvis-done", "ttl_min": 45 }
}
JSON

# ---------------------------------------------------------------------------
# Test 1: claim <id> <project> → calls update --tag jarvis-claimed
# ---------------------------------------------------------------------------
echo "=== Test 1: claim calls update --tag jarvis-claimed ==="

: > "$tmplog"
WORKITEM_ID="9001"
PROJECT_ID="1086837"

# Stub: list returns JSON array containing the workitem id (claim succeeds)
cat > "$tmpbin/a1" << STUB
#!/bin/bash
echo "\$@" >> "$tmplog"
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo '[{"id":"9001"}]'
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"

output1=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" claim "$WORKITEM_ID" "$PROJECT_ID" 2>&1)
exit1=$?

echo "Output: $output1"
echo "Exit code: $exit1"
echo "Recorded a1 calls:"
cat "$tmplog"
echo ""

# Must call: a1 project workitem update <id> --project <project> --tag jarvis-claimed
if grep -q "project workitem update $WORKITEM_ID" "$tmplog" && grep -q -- "--tag jarvis-claimed" "$tmplog"; then
    assert_pass "claim calls: a1 project workitem update <id> --tag jarvis-claimed"
else
    assert_fail "claim should call: a1 project workitem update <id> --tag jarvis-claimed (log: $(cat "$tmplog"))"
fi

# Must call: a1 project workitem comment create <id> -m "jarvis-claim ..."
if grep -q "project workitem comment create $WORKITEM_ID" "$tmplog" && grep -q "jarvis-claim" "$tmplog"; then
    assert_pass "claim calls: a1 project workitem comment create with jarvis-claim"
else
    assert_fail "claim should call: a1 project workitem comment create <id> -m jarvis-claim... (log: $(cat "$tmplog"))"
fi

# Comment must include hostname
EXPECTED_HOST=$(hostname)
if grep -q "jarvis-claim $EXPECTED_HOST" "$tmplog"; then
    assert_pass "claim comment contains hostname ($EXPECTED_HOST)"
else
    assert_fail "claim comment should contain hostname $EXPECTED_HOST (log: $(cat "$tmplog"))"
fi

# Readback list must pass --tag jarvis-claimed
if grep -q "project workitem list" "$tmplog" && grep -q -- "--tag jarvis-claimed" "$tmplog"; then
    assert_pass "claim readback list passes --tag jarvis-claimed"
else
    assert_fail "claim readback should pass --tag jarvis-claimed to list (log: $(cat "$tmplog"))"
fi

# Exit 0 when readback confirms our id (now as string)
if [ "$exit1" -eq 0 ]; then
    assert_pass "claim exits 0 when readback confirms claimed id (string match)"
else
    assert_fail "claim should exit 0 when readback confirms id, got $exit1"
fi

# ---------------------------------------------------------------------------
# Test 2: claim readback-not-mine → exits nonzero (lost race simulation)
# ---------------------------------------------------------------------------
echo "=== Test 2: claim readback-not-mine → nonzero exit (lost race) ==="

: > "$tmplog"

# Stub: list returns array WITHOUT our workitem id → lost the race
cat > "$tmpbin/a1" << STUB
#!/bin/bash
echo "\$@" >> "$tmplog"
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo '[{"id":8888}]'
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"

output2=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" claim "$WORKITEM_ID" "$PROJECT_ID" 2>&1)
exit2=$?

echo "Output: $output2"
echo "Exit code: $exit2"
echo ""

if [ "$exit2" -ne 0 ]; then
    assert_pass "claim exits nonzero when readback does not include our id (lost race)"
else
    assert_fail "claim should exit nonzero when readback does not include our id, got $exit2"
fi

# ---------------------------------------------------------------------------
# Test 3: release <id> <project> → calls update --tag jarvis-done, exits 0
# ---------------------------------------------------------------------------
echo "=== Test 3: release calls update --tag jarvis-done, exits 0 ==="

: > "$tmplog"

cat > "$tmpbin/a1" << STUB
#!/bin/bash
echo "\$@" >> "$tmplog"
exit 0
STUB
chmod +x "$tmpbin/a1"

output3=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" release "$WORKITEM_ID" "$PROJECT_ID" 2>&1)
exit3=$?

echo "Output: $output3"
echo "Exit code: $exit3"
echo "Recorded a1 calls:"
cat "$tmplog"
echo ""

if grep -q "project workitem update $WORKITEM_ID" "$tmplog" && grep -q -- "--tag jarvis-done" "$tmplog"; then
    assert_pass "release calls: a1 project workitem update <id> --tag jarvis-done"
else
    assert_fail "release should call: a1 project workitem update <id> --tag jarvis-done (log: $(cat "$tmplog"))"
fi

if [ "$exit3" -eq 0 ]; then
    assert_pass "release exits 0"
else
    assert_fail "release should exit 0, got $exit3"
fi

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo "=== Summary ==="
echo "PASS: $PASS  FAIL: $FAIL"

if [ "$FAIL" -gt 0 ]; then
    echo "TESTS FAILED"
    exit 1
fi

echo "All tests passed"
exit 0
