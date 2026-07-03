#!/usr/bin/env bash
# test/claim_test.sh – TDD tests for bootstrap/claim.sh
#
# Stateful a1 stub: `workitem get -f json` echoes the current tag set from a state
# file; `workitem update --tag <csv>` captures the written list and (unless NOOP)
# persists it back to the state file — faithfully modelling a1's whole-set tag
# replace so we can assert both the point-read readback and tag preservation.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

tmpbin=$(mktemp -d)
tmpconfig=$(mktemp -d)
tmplog=$(mktemp)
tmpstate=$(mktemp)
tmpcapture=$(mktemp)
tmpgetcnt=$(mktemp)
trap 'rm -rf "$tmpbin" "$tmpconfig"; rm -f "$tmplog" "$tmpstate" "$tmpcapture" "$tmpgetcnt"' EXIT

PASS=0
FAIL=0

assert_pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
assert_fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# ---------------------------------------------------------------------------
# tmp pools.json with a full claim section (tag/idle_tag/done_tag/done_status)
# ---------------------------------------------------------------------------
mkdir -p "$tmpconfig/config"
cat > "$tmpconfig/config/pools.json" << 'JSON'
{
  "claim": {
    "tag": "jarvis-claimed",
    "idle_tag": "jarvis-idle",
    "done_tag": "jarvis-done",
    "done_status": "已发布待需求排期",
    "ttl_min": 45
  }
}
JSON

# ---------------------------------------------------------------------------
# Stateful a1 stub. Reads A1_* env vars set per invocation:
#   A1_LOG        – append every call here
#   A1_STATE      – current tag CSV (get reads, update persists)
#   A1_CAPTURE    – last update's --tag value
#   A1_GETCNT     – counter file for GET_FAIL=first
#   A1_GET_FAIL   – "all" | "first" | unset
#   A1_UPDATE_NOOP – "1" → update captures but does NOT persist (lost-race sim)
# ---------------------------------------------------------------------------
cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
echo "$@" >> "$A1_LOG"
if [ "$1 $2 $3" = "project workitem get" ]; then
    if [ "${A1_GET_FAIL:-}" = "all" ]; then exit 1; fi
    if [ "${A1_GET_FAIL:-}" = "first" ]; then
        c=$(cat "$A1_GETCNT" 2>/dev/null || echo 0); c=$((c + 1)); echo "$c" > "$A1_GETCNT"
        if [ "$c" = "1" ]; then exit 1; fi
    fi
    tags=$(cat "$A1_STATE" 2>/dev/null || echo "")
    printf '{"fields":[{"identifier":"tag","displayValue":"%s","format":"multiList"}]}\n' "$tags"
    exit 0
fi
if [ "$1 $2 $3" = "project workitem update" ]; then
    args=("$@"); tagval=""; i=0
    while [ $i -lt ${#args[@]} ]; do
        if [ "${args[$i]}" = "--tag" ]; then j=$((i + 1)); tagval="${args[$j]}"; fi
        i=$((i + 1))
    done
    printf '%s' "$tagval" > "$A1_CAPTURE"
    if [ -z "${A1_UPDATE_NOOP:-}" ]; then printf '%s' "$tagval" > "$A1_STATE"; fi
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"

export A1_LOG="$tmplog" A1_STATE="$tmpstate" A1_CAPTURE="$tmpcapture" A1_GETCNT="$tmpgetcnt"

WORKITEM_ID="9001"
PROJECT_ID="1086837"

# run_claim <subcmd> — resets log/capture/getcnt, runs claim.sh, sets $out/$rc
run_claim() {
    : > "$tmplog"; : > "$tmpcapture"; : > "$tmpgetcnt"
    out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" JARVIS_CLAIM_READBACK_SLEEP=0 \
        bash "$proj_root/bootstrap/claim.sh" "$1" "$WORKITEM_ID" "$PROJECT_ID" 2>&1)
    rc=$?
}

# ---------------------------------------------------------------------------
# Test 1: readback via point-read (workitem get) sees claimed → exit 0
# ---------------------------------------------------------------------------
echo "=== Test 1: point-read readback sees claimed → exit 0 ==="
printf '' > "$tmpstate"           # start with no tags
unset A1_GET_FAIL A1_UPDATE_NOOP
run_claim claim
echo "Output: $out"; echo "Exit: $rc"; echo "a1 calls:"; cat "$tmplog"; echo

if [ "$rc" -eq 0 ]; then
    assert_pass "claim exits 0 when point-read readback shows claimed tag"
else
    assert_fail "claim should exit 0 when readback shows claimed, got $rc"
fi
if grep -q "project workitem get $WORKITEM_ID" "$tmplog"; then
    assert_pass "claim readback uses point-read: a1 project workitem get <id>"
else
    assert_fail "claim readback should use 'workitem get' (log: $(cat "$tmplog"))"
fi
if grep -q "project workitem list" "$tmplog"; then
    assert_fail "claim must NOT use 'workitem list' readback anymore (log: $(cat "$tmplog"))"
else
    assert_pass "claim no longer uses search-index 'workitem list' readback"
fi
if grep -q "project workitem update $WORKITEM_ID" "$tmplog" && [ "$(cat "$tmpcapture")" = "jarvis-claimed" ]; then
    assert_pass "claim writes --tag jarvis-claimed (no prior tags)"
else
    assert_fail "claim should write --tag jarvis-claimed, capture=$(cat "$tmpcapture")"
fi

# ---------------------------------------------------------------------------
# Test 2: point-read never shows claimed (update no-op) → exit 1 (lost race)
# ---------------------------------------------------------------------------
echo "=== Test 2: point-read never shows claimed → exit 1 (lost race) ==="
printf '' > "$tmpstate"
unset A1_GET_FAIL
export A1_UPDATE_NOOP=1
run_claim claim
unset A1_UPDATE_NOOP
echo "Output: $out"; echo "Exit: $rc"; echo

if [ "$rc" -ne 0 ]; then
    assert_pass "claim exits nonzero when claimed tag never visible in point-read (lost race)"
else
    assert_fail "claim should exit nonzero when readback lacks claimed, got $rc"
fi
if grep -q "project workitem list" "$tmplog"; then
    assert_fail "lost-race path must not use 'workitem list' (log: $(cat "$tmplog"))"
else
    assert_pass "lost-race path uses point-read only"
fi

# ---------------------------------------------------------------------------
# Test 3: tag preservation on CLAIM — jarvis-probe survives, claimed added
# ---------------------------------------------------------------------------
echo "=== Test 3: claim preserves jarvis-probe, adds jarvis-claimed ==="
printf 'jarvis-probe' > "$tmpstate"
unset A1_GET_FAIL A1_UPDATE_NOOP
run_claim claim
cap=$(cat "$tmpcapture")
echo "Captured --tag: $cap"; echo "Exit: $rc"; echo

if [ "$rc" -eq 0 ]; then assert_pass "claim (with prior probe) exits 0"; else assert_fail "claim exit $rc"; fi
if echo "$cap" | grep -q "jarvis-probe" && echo "$cap" | grep -q "jarvis-claimed"; then
    assert_pass "claim --tag full list preserves jarvis-probe AND adds jarvis-claimed ($cap)"
else
    assert_fail "claim --tag should contain both jarvis-probe and jarvis-claimed, got '$cap'"
fi

# ---------------------------------------------------------------------------
# Test 4: tag preservation on RELEASE — probe survives, claimed→idle
# ---------------------------------------------------------------------------
echo "=== Test 4: release preserves jarvis-probe, swaps claimed→idle ==="
printf 'jarvis-probe,jarvis-claimed' > "$tmpstate"
unset A1_GET_FAIL A1_UPDATE_NOOP
run_claim release
cap=$(cat "$tmpcapture")
echo "Captured --tag: $cap"; echo "Exit: $rc"; echo

if [ "$rc" -eq 0 ]; then assert_pass "release exits 0"; else assert_fail "release exit $rc"; fi
if echo "$cap" | grep -q "jarvis-probe" && echo "$cap" | grep -q "jarvis-idle" && ! echo "$cap" | grep -q "jarvis-claimed"; then
    assert_pass "release --tag = probe + idle, claimed dropped ($cap)"
else
    assert_fail "release --tag should be probe+idle without claimed, got '$cap'"
fi

# ---------------------------------------------------------------------------
# Test 5: tag preservation on FINISH — probe survives, done added, status set
# ---------------------------------------------------------------------------
echo "=== Test 5: finish preserves jarvis-probe, adds done, sets status ==="
printf 'jarvis-probe,jarvis-claimed' > "$tmpstate"
unset A1_GET_FAIL A1_UPDATE_NOOP
run_claim finish
cap=$(cat "$tmpcapture")
echo "Captured --tag: $cap"; echo "Exit: $rc"; echo "a1 calls:"; cat "$tmplog"; echo

if [ "$rc" -eq 0 ]; then assert_pass "finish exits 0"; else assert_fail "finish exit $rc"; fi
if echo "$cap" | grep -q "jarvis-probe" && echo "$cap" | grep -q "jarvis-done" \
   && ! echo "$cap" | grep -q "jarvis-claimed" && ! echo "$cap" | grep -q "jarvis-idle"; then
    assert_pass "finish --tag = probe + done, claimed/idle dropped ($cap)"
else
    assert_fail "finish --tag should be probe+done without claimed/idle, got '$cap'"
fi
if grep -q -- "--status 已发布待需求排期" "$tmplog"; then
    assert_pass "finish passes --status 已发布待需求排期"
else
    assert_fail "finish should pass --status done_status (log: $(cat "$tmplog"))"
fi

# ---------------------------------------------------------------------------
# Test 6: get-fail degrade — warning + legacy bare-tag write, still exits 0
# (exercised on release: no readback, so it isolates the write-degrade branch)
# ---------------------------------------------------------------------------
echo "=== Test 6: get-fail → warning + legacy single --tag write (release) ==="
printf 'jarvis-probe,jarvis-claimed' > "$tmpstate"
export A1_GET_FAIL=all
unset A1_UPDATE_NOOP
run_claim release
unset A1_GET_FAIL
cap=$(cat "$tmpcapture")
echo "Output: $out"; echo "Captured --tag: $cap"; echo "Exit: $rc"; echo

if [ "$rc" -eq 0 ]; then assert_pass "release still exits 0 when tags unreadable"; else assert_fail "release exit $rc"; fi
if echo "$out" | grep -qi "could not read existing tags"; then
    assert_pass "degrade emits stderr warning about unreadable tags"
else
    assert_fail "degrade should warn about unreadable tags, output: $out"
fi
if [ "$cap" = "jarvis-idle" ]; then
    assert_pass "degrade falls back to legacy bare --tag jarvis-idle (no merge)"
else
    assert_fail "degrade should write only jarvis-idle, got '$cap'"
fi

# ---------------------------------------------------------------------------
# Test 7: claim degrade on merge-read (first get fails) then readback confirms
# ---------------------------------------------------------------------------
echo "=== Test 7: claim degrade — first get fails, legacy write, readback wins ==="
printf '' > "$tmpstate"
export A1_GET_FAIL=first
unset A1_UPDATE_NOOP
run_claim claim
unset A1_GET_FAIL
cap=$(cat "$tmpcapture")
echo "Output: $out"; echo "Captured --tag: $cap"; echo "Exit: $rc"; echo

if [ "$rc" -eq 0 ]; then assert_pass "claim recovers via readback after merge-read get failure"; else assert_fail "claim exit $rc"; fi
if echo "$out" | grep -qi "could not read existing tags"; then
    assert_pass "claim merge-read failure emits warning"
else
    assert_fail "claim merge-read failure should warn, output: $out"
fi
if [ "$cap" = "jarvis-claimed" ]; then
    assert_pass "claim degrade writes legacy bare --tag jarvis-claimed"
else
    assert_fail "claim degrade should write only jarvis-claimed, got '$cap'"
fi

# ---------------------------------------------------------------------------
# Test 8: ledger records owner=COORD_ID on claim (INSERT)
# ---------------------------------------------------------------------------
echo "=== Test 8: claim writes owner=COORD_ID into ledger ==="
today="$(date -u +%F)"
ledger="$tmpconfig/.my-day/claims-${today}.json"
rm -f "$ledger"                    # start from a clean ledger for the owner cases
printf '' > "$tmpstate"
unset A1_GET_FAIL A1_UPDATE_NOOP
export COORD_ID="inst-A"
run_claim claim
echo "Exit: $rc"; echo "Ledger: $(cat "$ledger" 2>/dev/null)"; echo

led_owner=$(jq -r --arg id "$WORKITEM_ID" '.[] | select(.id==$id) | .owner' "$ledger" 2>/dev/null)
led_done=$(jq -r --arg id "$WORKITEM_ID" '.[] | select(.id==$id) | .done' "$ledger" 2>/dev/null)
if [ "$led_owner" = "inst-A" ]; then
    assert_pass "claim ledger entry owner == COORD_ID (inst-A)"
else
    assert_fail "claim ledger owner should be inst-A, got '$led_owner'"
fi
if [ "$led_done" = "false" ]; then
    assert_pass "claim ledger entry done == false"
else
    assert_fail "claim ledger done should be false, got '$led_done'"
fi

# ---------------------------------------------------------------------------
# Test 9: finish from a DIFFERENT instance does NOT overwrite original owner (coalesce)
# ---------------------------------------------------------------------------
echo "=== Test 9: finish by other instance keeps original owner (coalesce) ==="
printf 'jarvis-probe,jarvis-claimed' > "$tmpstate"
unset A1_GET_FAIL A1_UPDATE_NOOP
export COORD_ID="inst-B"           # a different instance closes it
run_claim finish
echo "Exit: $rc"; echo "Ledger: $(cat "$ledger" 2>/dev/null)"; echo

led_owner=$(jq -r --arg id "$WORKITEM_ID" '.[] | select(.id==$id) | .owner' "$ledger" 2>/dev/null)
led_done=$(jq -r --arg id "$WORKITEM_ID" '.[] | select(.id==$id) | .done' "$ledger" 2>/dev/null)
if [ "$led_owner" = "inst-A" ]; then
    assert_pass "finish by inst-B preserves original owner inst-A (coalesce, no overwrite)"
else
    assert_fail "finish should keep owner inst-A, got '$led_owner'"
fi
if [ "$led_done" = "true" ]; then
    assert_pass "finish sets ledger done == true"
else
    assert_fail "finish should set done true, got '$led_done'"
fi

# ---------------------------------------------------------------------------
# Test 10: UPDATE backfills owner when the original entry had an empty owner
# ---------------------------------------------------------------------------
echo "=== Test 10: update backfills owner when original owner was empty ==="
printf '[{"id":"%s","done":false,"owner":""}]' "$WORKITEM_ID" > "$ledger"
printf 'jarvis-probe,jarvis-claimed' > "$tmpstate"
unset A1_GET_FAIL A1_UPDATE_NOOP
export COORD_ID="inst-C"
run_claim release
echo "Exit: $rc"; echo "Ledger: $(cat "$ledger" 2>/dev/null)"; echo
unset COORD_ID

led_owner=$(jq -r --arg id "$WORKITEM_ID" '.[] | select(.id==$id) | .owner' "$ledger" 2>/dev/null)
if [ "$led_owner" = "inst-C" ]; then
    assert_pass "update backfills owner (empty → inst-C)"
else
    assert_fail "update should backfill empty owner to inst-C, got '$led_owner'"
fi

# ---------------------------------------------------------------------------
# Test 11 (D2): no COORD_ID but CLAUDE_CODE_SESSION_ID set → owner = cc-<sid>
# ---------------------------------------------------------------------------
echo "=== Test 11: interactive session owner = cc-<CLAUDE_CODE_SESSION_ID> ==="
rm -f "$ledger"
printf '' > "$tmpstate"
unset A1_GET_FAIL A1_UPDATE_NOOP COORD_ID
export CLAUDE_CODE_SESSION_ID="sess-D2"
run_claim claim
echo "Exit: $rc"; echo "Ledger: $(cat "$ledger" 2>/dev/null)"; echo
led_owner=$(jq -r --arg id "$WORKITEM_ID" '.[] | select(.id==$id) | .owner' "$ledger" 2>/dev/null)
if [ "$led_owner" = "cc-sess-D2" ]; then
    assert_pass "interactive claim owner == cc-sess-D2"
else
    assert_fail "interactive owner should be cc-sess-D2, got '$led_owner'"
fi
unset CLAUDE_CODE_SESSION_ID

# ---------------------------------------------------------------------------
# Test 12 (D2): COORD_ID takes precedence over CLAUDE_CODE_SESSION_ID
# ---------------------------------------------------------------------------
echo "=== Test 12: COORD_ID wins over CLAUDE_CODE_SESSION_ID ==="
rm -f "$ledger"
printf '' > "$tmpstate"
unset A1_GET_FAIL A1_UPDATE_NOOP
export COORD_ID="inst-P" CLAUDE_CODE_SESSION_ID="sess-Q"
run_claim claim
echo "Exit: $rc"; echo "Ledger: $(cat "$ledger" 2>/dev/null)"; echo
led_owner=$(jq -r --arg id "$WORKITEM_ID" '.[] | select(.id==$id) | .owner' "$ledger" 2>/dev/null)
if [ "$led_owner" = "inst-P" ]; then
    assert_pass "COORD_ID precedence: owner == inst-P (not cc-sess-Q)"
else
    assert_fail "COORD_ID should win, got '$led_owner'"
fi
unset COORD_ID CLAUDE_CODE_SESSION_ID

# ---------------------------------------------------------------------------
# Test 13 (D1): concurrent claims of distinct ids → no lost update (mkdir lock)
# ---------------------------------------------------------------------------
echo "=== Test 13: concurrent ledger upserts keep all entries ==="
rm -f "$ledger"
unset COORD_ID CLAUDE_CODE_SESSION_ID A1_GET_FAIL A1_UPDATE_NOOP
N=12
for i in $(seq 1 "$N"); do
    PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" JARVIS_CLAIM_READBACK_SLEEP=0 \
        A1_LOG="$(mktemp)" A1_STATE="$(mktemp)" A1_CAPTURE="$(mktemp)" A1_GETCNT="$(mktemp)" \
        bash "$proj_root/bootstrap/claim.sh" claim "77${i}" "$PROJECT_ID" >/dev/null 2>&1 &
done
wait
cnt=$(jq 'length' "$ledger" 2>/dev/null || echo 0)
if [ "$cnt" = "$N" ]; then
    assert_pass "concurrent upserts: all $N entries present (no lost update)"
else
    assert_fail "lost update under concurrency: $cnt/$N entries in ledger"
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
