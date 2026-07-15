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
trap 'rm -rf "$tmpbin" "$tmpconfig"; rm -f "$tmplog" "$tmpstate" "$tmpcapture" "$tmpgetcnt" "${tmpstatuscap:-}"' EXIT

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
    "progress_status": "处理中",
    "start_statuses": ["待处理", "新建", "New", "待认领", "Reopen"],
    "ttl_min": 45
  },
  "pools": {
    "with_override": { "project": 2100304, "name": "override pool", "done_status": "已完成", "progress_status": "处理中" },
    "cat_override":  { "project": 2100305, "name": "per-category pool", "done_status": {"需求": "已发布", "功能缺陷": "Fixed"}, "progress_status": {"需求": "问题解决中", "功能缺陷": "Open"} },
    "no_override":   { "project": 1086837, "name": "plain pool" }
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
#   A1_STATUS_CAPTURE – file to record the --status value of an update (optional)
#   A1_REJECT_STATUS  – if set and equals the update's --status value, exit 1 (enum mismatch sim)
#   A1_COMMENTS       – shared comments JSON file (array of {"content":...}); comment
#                       create appends, comment list echoes it. Models one Aone workitem's
#                       comment stream shared across multiple "machines" for arbitration.
#   A1_WTYPE          – workitem type displayValue returned by `workitem get` in a
#                       fields[] entry identifier=="workitemType" (default 需求). Lets
#                       finish exercise per-category done_status selection.
#   A1_STATUS         – current status displayValue returned by `workitem get` in a
#                       fields[] entry identifier=="status" (default 待处理, a start state).
#                       Drives claim's status-advance path (advance only from a start state).
# A --tag-less update (e.g. finish's separate --status call) must NOT clobber A1_STATE tags.
# ---------------------------------------------------------------------------
cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
echo "$@" >> "$A1_LOG"
if [ "$1 $2 $3 $4" = "project workitem comment create" ]; then
    args=("$@"); content=""; i=0
    while [ $i -lt ${#args[@]} ]; do
        if [ "${args[$i]}" = "-m" ]; then j=$((i + 1)); content="${args[$j]}"; fi
        i=$((i + 1))
    done
    f="${A1_COMMENTS:-}"
    if [ -n "$f" ]; then
        [ -s "$f" ] || printf '[]' > "$f"
        ctmp=$(mktemp); jq --arg c "$content" '. + [{"content":$c}]' "$f" > "$ctmp" && mv "$ctmp" "$f"
    fi
    exit 0
fi
if [ "$1 $2 $3 $4" = "project workitem comment list" ]; then
    if [ -n "${A1_COMMENTS:-}" ] && [ -s "$A1_COMMENTS" ]; then cat "$A1_COMMENTS"; else echo '[]'; fi
    exit 0
fi
if [ "$1 $2 $3" = "project workitem get" ]; then
    if [ "${A1_GET_FAIL:-}" = "all" ]; then exit 1; fi
    if [ "${A1_GET_FAIL:-}" = "first" ]; then
        c=$(cat "$A1_GETCNT" 2>/dev/null || echo 0); c=$((c + 1)); echo "$c" > "$A1_GETCNT"
        if [ "$c" = "1" ]; then exit 1; fi
    fi
    tags=$(cat "$A1_STATE" 2>/dev/null || echo "")
    # Real Aone multiList tag fields join displayValue/value with a comma-SPACE, and `value`
    # carries the tag IDs (which _get_tag_pairs preserves by, so cross-project tags survive).
    # The stub stores the tag set comma-only in A1_STATE; convert to comma-space here and let
    # NAMES double as their own IDs (value == displayValue) so the by-ID preservation path
    # round-trips and the name-based assertions below still hold.
    tags_ml="${tags//,/, }"
    wtype="${A1_WTYPE:-需求}"
    status="${A1_STATUS:-待处理}"
    printf '{"fields":[{"identifier":"tag","displayValue":"%s","value":"%s","format":"multiList"},{"identifier":"workitemType","displayValue":"%s"},{"identifier":"status","displayValue":"%s"}]}\n' "$tags_ml" "$tags_ml" "$wtype" "$status"
    exit 0
fi
if [ "$1 $2 $3" = "project workitem update" ]; then
    args=("$@"); tagval=""; statusval=""; has_tag=""; has_status=""; i=0
    while [ $i -lt ${#args[@]} ]; do
        if [ "${args[$i]}" = "--tag" ]; then j=$((i + 1)); tagval="${args[$j]}"; has_tag=1; fi
        if [ "${args[$i]}" = "--status" ]; then j=$((i + 1)); statusval="${args[$j]}"; has_status=1; fi
        i=$((i + 1))
    done
    if [ -n "$has_status" ]; then
        [ -n "${A1_STATUS_CAPTURE:-}" ] && printf '%s' "$statusval" > "$A1_STATUS_CAPTURE"
        if [ -n "${A1_REJECT_STATUS:-}" ] && [ "$statusval" = "$A1_REJECT_STATUS" ]; then
            echo "Error: unsupported target status \"$statusval\"" >&2; exit 1
        fi
    fi
    if [ -n "$has_tag" ]; then
        printf '%s' "$tagval" > "$A1_CAPTURE"
        if [ -z "${A1_UPDATE_NOOP:-}" ]; then printf '%s' "$tagval" > "$A1_STATE"; fi
    fi
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"
export JARVIS_A1=a1   # 走 PATH 上的 a1 stub,不经真 bin/a1id

tmpstatuscap=$(mktemp)
export A1_LOG="$tmplog" A1_STATE="$tmpstate" A1_CAPTURE="$tmpcapture" A1_GETCNT="$tmpgetcnt" A1_STATUS_CAPTURE="$tmpstatuscap"

WORKITEM_ID="9001"
PROJECT_ID="1086837"

# run_claim <subcmd> — resets log/capture/getcnt, runs claim.sh, sets $out/$rc
run_claim() {
    : > "$tmplog"; : > "$tmpcapture"; : > "$tmpgetcnt"
    out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" JARVIS_CLAIM_READBACK_SLEEP=0 \
        bash "$proj_root/bootstrap/claim.sh" "$1" "$WORKITEM_ID" "$PROJECT_ID" 2>&1)
    rc=$?
}

# run_claim_prog <project> — claim exercising the status-advance path. Resets log/capture/
# getcnt AND the status-capture file, so each case sees a clean status write (or the absence
# of one). Callers set A1_STATUS / A1_WTYPE / A1_REJECT_STATUS / JARVIS_CLAIM_PROGRESS in the
# environment beforehand. Sets $out/$rc; the written status lands in $tmpstatuscap.
run_claim_prog() {
    : > "$tmplog"; : > "$tmpcapture"; : > "$tmpgetcnt"; : > "$tmpstatuscap"; : > "$tmpstate"
    out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" JARVIS_CLAIM_READBACK_SLEEP=0 \
        bash "$proj_root/bootstrap/claim.sh" claim "$WORKITEM_ID" "$1" 2>&1)
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
# Test 14: finish uses per-pool done_status override (project 2100304 → 已完成)
# ---------------------------------------------------------------------------
echo "=== Test 14: finish uses per-pool done_status override ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS COORD_ID
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 2100304 2>&1); rc=$?
cap=$(cat "$tmpstatuscap" 2>/dev/null); state=$(cat "$tmpstate" 2>/dev/null)
echo "rc=$rc status_set='$cap' tags='$state'"
if [ "$cap" = "已完成" ]; then assert_pass "per-pool: status set to 已完成"; else assert_fail "per-pool status should be 已完成, got '$cap'"; fi
if printf '%s' "$state" | grep -q "jarvis-done"; then assert_pass "per-pool: tagged jarvis-done"; else assert_fail "per-pool: expected jarvis-done, got '$state'"; fi

# ---------------------------------------------------------------------------
# Test 15: finish falls back to global done_status when pool has no override
# ---------------------------------------------------------------------------
echo "=== Test 15: finish falls back to global done_status ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
unset A1_REJECT_STATUS
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 1086837 2>&1); rc=$?
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc status_set='$cap'"
if [ "$cap" = "已发布待需求排期" ]; then assert_pass "global fallback: status set to 已发布待需求排期"; else assert_fail "global fallback wrong, got '$cap'"; fi

# ---------------------------------------------------------------------------
# Test 16: rejected status → finish DOWNGRADES jarvis-done → jarvis-idle + escalates.
# 一致性闸门:done_status 落不到合法完成态时,继续挂 jarvis-done 会造成「标签 done、真源没完」
# 的黑洞(_decide 永久 skip)。故降级为 jarvis-idle(交 idle 门/Revisit 兜底) + 写 escalation,
# 仍 exit 0 + warns(不 fatal,bookend 流程正常收尾)。
# ---------------------------------------------------------------------------
echo "=== Test 16: rejected status → downgrade jarvis-done→jarvis-idle + escalate ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
rm -f "$tmpconfig/escalation/$WORKITEM_ID.md"
export A1_REJECT_STATUS="已发布待需求排期" A1_STATUS="处理中"
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 1086837 2>&1); rc=$?
state=$(cat "$tmpstate" 2>/dev/null)
unset A1_REJECT_STATUS A1_STATUS
echo "rc=$rc tags='$state'"; echo "out: $out"
if [ "$rc" = "0" ]; then assert_pass "rejected status: finish still exits 0"; else assert_fail "rejected status: exit $rc"; fi
if printf '%s' "$state" | grep -q "jarvis-idle" && ! printf '%s' "$state" | grep -q "jarvis-done"; then
    assert_pass "rejected status: downgraded jarvis-done→jarvis-idle (no black hole), got '$state'"
else
    assert_fail "rejected status: should downgrade to jarvis-idle without jarvis-done, got '$state'"
fi
if printf '%s' "$out" | grep -qi "warning"; then assert_pass "rejected status: emits warning"; else assert_fail "rejected status: no warning emitted"; fi
if [ -f "$tmpconfig/escalation/$WORKITEM_ID.md" ] && grep -q "finish_status_unresolved" "$tmpconfig/escalation/$WORKITEM_ID.md"; then
    assert_pass "rejected status: escalation file written with finish_status_unresolved"
else
    assert_fail "rejected status: expected escalation/$WORKITEM_ID.md with finish_status_unresolved"
fi

# ---------------------------------------------------------------------------
# Test 16b: successful status write KEEPS jarvis-done (no spurious downgrade).
# ---------------------------------------------------------------------------
echo "=== Test 16b: finish with accepted status keeps jarvis-done ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS COORD_ID
export A1_STATUS="处理中"
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 1086837 2>&1); rc=$?
state=$(cat "$tmpstate" 2>/dev/null)
unset A1_STATUS
echo "rc=$rc tags='$state'"
if [ "$rc" = "0" ]; then assert_pass "accepted status: finish exits 0"; else assert_fail "accepted status: exit $rc"; fi
if printf '%s' "$state" | grep -q "jarvis-done" && ! printf '%s' "$state" | grep -q "jarvis-idle"; then
    assert_pass "accepted status: keeps jarvis-done, no downgrade ($state)"
else
    assert_fail "accepted status: should keep jarvis-done without idle, got '$state'"
fi

# ---------------------------------------------------------------------------
# Test 17: per-category done_status — bug (功能缺陷) → Fixed
# The pool's done_status is an object keyed by workitem type; finish must read the
# workitem type (fields[] identifier==workitemType displayValue) and select its status.
# ---------------------------------------------------------------------------
echo "=== Test 17: per-category done_status — 功能缺陷 → Fixed ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS COORD_ID
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" A1_WTYPE=功能缺陷 \
    bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 2100305 2>&1); rc=$?
cap=$(cat "$tmpstatuscap" 2>/dev/null); state=$(cat "$tmpstate" 2>/dev/null)
echo "rc=$rc status_set='$cap' tags='$state'"
if [ "$cap" = "Fixed" ]; then assert_pass "per-category 功能缺陷: status set to Fixed"; else assert_fail "per-category 功能缺陷 status should be Fixed, got '$cap'"; fi
if printf '%s' "$state" | grep -q "jarvis-done"; then assert_pass "per-category 功能缺陷: tagged jarvis-done"; else assert_fail "per-category 功能缺陷: expected jarvis-done, got '$state'"; fi

# ---------------------------------------------------------------------------
# Test 18: per-category done_status — req (需求) → 已发布
# ---------------------------------------------------------------------------
echo "=== Test 18: per-category done_status — 需求 → 已发布 ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS COORD_ID
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" A1_WTYPE=需求 \
    bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 2100305 2>&1); rc=$?
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc status_set='$cap'"
if [ "$cap" = "已发布" ]; then assert_pass "per-category 需求: status set to 已发布"; else assert_fail "per-category 需求 status should be 已发布, got '$cap'"; fi

# ---------------------------------------------------------------------------
# Test 19: per-category unknown type falls back to first object value (avoid empty)
# ---------------------------------------------------------------------------
echo "=== Test 19: per-category unknown type → first object value fallback ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS COORD_ID
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" A1_WTYPE=任务 \
    bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 2100305 2>&1); rc=$?
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc status_set='$cap'"
if [ "$cap" = "已发布" ]; then assert_pass "per-category unknown type falls back to first value (已发布)"; else assert_fail "per-category unknown type should fall back to first value 已发布, got '$cap'"; fi

# ===========================================================================
# Cross-machine claim arbitration (JARVIS_CLAIM_SETTLE>0).
# Two "machines" (distinct JARVIS_SELF_HOST) claim the SAME workitem, sharing one
# A1_STATE (tag set) and one A1_COMMENTS (comment stream) = one Aone workitem. Run
# sequentially so hostA's claim comment carries the earlier UTC timestamp (hostA
# starts and finishes — incl. its settle sleep — before hostB begins), modelling
# near-simultaneous claims with a deterministic global order. hostA (earliest) must
# win; hostB (later) must detect the earlier claim, stand down, and roll back to idle.
# ===========================================================================
echo "=== Test 17: cross-machine arbitration — earliest host wins, later stands down ==="
ARB_STATE=$(mktemp); ARB_COMMENTS=$(mktemp)
ARB_LOGA=$(mktemp); ARB_CAPA=$(mktemp); ARB_CNTA=$(mktemp)
ARB_LOGB=$(mktemp); ARB_CAPB=$(mktemp); ARB_CNTB=$(mktemp)
printf '' > "$ARB_STATE"          # start with no tags on the shared workitem
printf '[]' > "$ARB_COMMENTS"     # empty comment stream
ARB_ID="8001"
unset COORD_ID CLAUDE_CODE_SESSION_ID

# hostA claims first (earliest timestamp). JARVIS_ROOT=$tmpconfig so config/pools.json resolves.
outA=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" JARVIS_CLAIM_READBACK_SLEEP=0 \
    JARVIS_CLAIM_SETTLE=1 JARVIS_SELF_HOST=host-A \
    A1_LOG="$ARB_LOGA" A1_STATE="$ARB_STATE" A1_CAPTURE="$ARB_CAPA" A1_GETCNT="$ARB_CNTA" \
    A1_COMMENTS="$ARB_COMMENTS" \
    bash "$proj_root/bootstrap/claim.sh" claim "$ARB_ID" "$PROJECT_ID" 2>&1); rcA=$?

# hostB claims the same item afterwards (later timestamp) → must lose to host-A
outB=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" JARVIS_CLAIM_READBACK_SLEEP=0 \
    JARVIS_CLAIM_SETTLE=1 JARVIS_SELF_HOST=host-B \
    A1_LOG="$ARB_LOGB" A1_STATE="$ARB_STATE" A1_CAPTURE="$ARB_CAPB" A1_GETCNT="$ARB_CNTB" \
    A1_COMMENTS="$ARB_COMMENTS" \
    bash "$proj_root/bootstrap/claim.sh" claim "$ARB_ID" "$PROJECT_ID" 2>&1); rcB=$?

echo "hostA rc=$rcA out=$outA"
echo "hostB rc=$rcB out=$outB"
echo "capB(last --tag)=$(cat "$ARB_CAPB")"
echo "comments=$(cat "$ARB_COMMENTS")"

# Exactly one winner
winners=0
[ "$rcA" -eq 0 ] && winners=$((winners + 1))
[ "$rcB" -eq 0 ] && winners=$((winners + 1))
if [ "$winners" -eq 1 ]; then assert_pass "exactly one host wins the claim"; else assert_fail "expected exactly one winner, rcA=$rcA rcB=$rcB"; fi

# host-A (earliest) wins
if [ "$rcA" -eq 0 ]; then assert_pass "host-A (earliest claim) exits 0 (won)"; else assert_fail "host-A should win, got rc=$rcA"; fi

# host-B (later) stands down
if [ "$rcB" -eq 1 ]; then assert_pass "host-B (later claim) exits 1 (stood down)"; else assert_fail "host-B should lose, got rc=$rcB"; fi

# host-B stood down WITHOUT clobbering the shared tag set: on loss it must NOT write idle
# (that would erase the winner's claimed, whole-set replace, no per-owner identity). Its
# last --tag write stays its own step-1 claimed (idempotent with the winner's).
capB=$(cat "$ARB_CAPB")
if echo "$capB" | grep -q "jarvis-claimed" && ! echo "$capB" | grep -q "jarvis-idle"; then
    assert_pass "host-B did not roll back to idle on loss (capB=$capB)"
else
    assert_fail "host-B must not write idle on loss (would clobber winner), got '$capB'"
fi

# CORE invariant: the winner's claim survives on the SHARED workitem — the arbitration
# loser must not have removed jarvis-claimed from the shared tag set.
sharedB=$(cat "$ARB_STATE")
if echo "$sharedB" | grep -q "jarvis-claimed"; then
    assert_pass "winner keeps jarvis-claimed on shared workitem (loser did not clobber)"
else
    assert_fail "shared tag set lost jarvis-claimed after loser stood down, got '$sharedB'"
fi

# host-B's stderr explains the stand-down
if echo "$outB" | grep -qi "arbitration"; then assert_pass "host-B emits arbitration stand-down message"; else assert_fail "host-B should mention arbitration, got: $outB"; fi

# Both machines wrote reconcile-parseable claim comments (nonce after the UTC timestamp,
# so reconcile.sh's existing regex still captures the timestamp).
recon_ts=$(python3 -c "
import json,sys,re
cs=json.load(open('$ARB_COMMENTS'))
pat=re.compile(r'jarvis-claim\s+\S+\s+(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)')
n=sum(1 for c in cs if pat.search(c.get('content','')))
print(n)
" 2>/dev/null || echo 0)
if [ "$recon_ts" = "2" ]; then
    assert_pass "both claim comments match reconcile.sh timestamp regex (nonce non-breaking)"
else
    assert_fail "expected 2 reconcile-parseable comments, got $recon_ts"
fi

# Nonces are distinct (8-hex tail after the timestamp)
distinct_nonces=$(python3 -c "
import json,re
cs=json.load(open('$ARB_COMMENTS'))
pat=re.compile(r'jarvis-claim\s+\S+\s+\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z\s+(\S+)')
ns=set()
for c in cs:
    m=pat.search(c.get('content',''))
    if m: ns.add(m.group(1))
print(len(ns))
" 2>/dev/null || echo 0)
if [ "$distinct_nonces" = "2" ]; then assert_pass "two distinct nonces present"; else assert_fail "expected 2 distinct nonces, got $distinct_nonces"; fi

rm -f "$ARB_STATE" "$ARB_COMMENTS" "$ARB_LOGA" "$ARB_CAPA" "$ARB_CNTA" "$ARB_LOGB" "$ARB_CAPB" "$ARB_CNTB"

# ===========================================================================
# claim → advance Aone status to the pool's in-progress status (best-effort, non-fatal,
# only from a starting state). A1_STATUS drives the workitem's CURRENT status; the write
# lands in $tmpstatuscap (reset per run by run_claim_prog).
# ===========================================================================

# ---------------------------------------------------------------------------
# Test 20: advance from a start state (string pool) — 待处理 → pool progress_status
# ---------------------------------------------------------------------------
echo "=== Test 20: claim advances status from start state (string pool → 处理中) ==="
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS A1_WTYPE COORD_ID JARVIS_CLAIM_PROGRESS
export A1_STATUS=待处理
run_claim_prog 2100304          # with_override, progress_status="处理中"
unset A1_STATUS
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc status_set='$cap'"
if [ "$rc" = "0" ]; then assert_pass "claim from start state exits 0"; else assert_fail "claim from start state exit $rc"; fi
if [ "$cap" = "处理中" ]; then assert_pass "string pool: status advanced to 处理中"; else assert_fail "string pool status should be 处理中, got '$cap'"; fi

# ---------------------------------------------------------------------------
# Test 21: advance per-category (object pool) — 功能缺陷 → Open, 需求 → 问题解决中
# ---------------------------------------------------------------------------
echo "=== Test 21: claim advances per-category status (object pool) ==="
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS COORD_ID JARVIS_CLAIM_PROGRESS
export A1_STATUS=待处理 A1_WTYPE=功能缺陷
run_claim_prog 2100305          # cat_override, progress_status={需求:问题解决中,功能缺陷:Open}
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc(bug) status_set='$cap'"
if [ "$cap" = "Open" ]; then assert_pass "object pool 功能缺陷: status advanced to Open"; else assert_fail "object pool 功能缺陷 status should be Open, got '$cap'"; fi

export A1_WTYPE=需求
run_claim_prog 2100305
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc(req) status_set='$cap'"
if [ "$cap" = "问题解决中" ]; then assert_pass "object pool 需求: status advanced to 问题解决中"; else assert_fail "object pool 需求 status should be 问题解决中, got '$cap'"; fi
unset A1_STATUS A1_WTYPE

# ---------------------------------------------------------------------------
# Test 22: NO advance when current status is NOT a start state (no backward move)
# ---------------------------------------------------------------------------
echo "=== Test 22: claim does NOT advance status when not in start_statuses ==="
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS A1_WTYPE COORD_ID JARVIS_CLAIM_PROGRESS
export A1_STATUS=问题解决中     # mid-flight, not a start state
run_claim_prog 2100304
unset A1_STATUS
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc status_set='$cap'"
if [ "$rc" = "0" ]; then assert_pass "claim from non-start state exits 0"; else assert_fail "claim from non-start state exit $rc"; fi
if [ -z "$cap" ]; then assert_pass "non-start state: no status write (stays empty)"; else assert_fail "non-start state should not write status, got '$cap'"; fi

# ---------------------------------------------------------------------------
# Test 23: rejected status is non-fatal — claim still exits 0 + WARN, tag still lands
# ---------------------------------------------------------------------------
echo "=== Test 23: claim resilient to a rejected status ==="
unset A1_GET_FAIL A1_UPDATE_NOOP A1_WTYPE COORD_ID JARVIS_CLAIM_PROGRESS
export A1_STATUS=待处理 A1_REJECT_STATUS=处理中
run_claim_prog 2100304          # string pool → progress 处理中, which a1 stub rejects
unset A1_STATUS A1_REJECT_STATUS
state=$(cat "$tmpstate" 2>/dev/null)
echo "rc=$rc out='$out' tags='$state'"
if [ "$rc" = "0" ]; then assert_pass "rejected status: claim still exits 0"; else assert_fail "rejected status: claim exit $rc"; fi
if printf '%s' "$out" | grep -qi "warning"; then assert_pass "rejected status: emits warning"; else assert_fail "rejected status: no warning emitted"; fi
if printf '%s' "$state" | grep -q "jarvis-claimed"; then assert_pass "rejected status: still tagged jarvis-claimed (lock intact)"; else assert_fail "rejected status: claimed tag lost, got '$state'"; fi

# ---------------------------------------------------------------------------
# Test 24: JARVIS_CLAIM_PROGRESS=0 disables the status advance
# ---------------------------------------------------------------------------
echo "=== Test 24: JARVIS_CLAIM_PROGRESS=0 disables status advance ==="
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS A1_WTYPE COORD_ID
export A1_STATUS=待处理 JARVIS_CLAIM_PROGRESS=0
run_claim_prog 2100304
unset A1_STATUS JARVIS_CLAIM_PROGRESS
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc status_set='$cap'"
if [ "$rc" = "0" ]; then assert_pass "JARVIS_CLAIM_PROGRESS=0: claim exits 0"; else assert_fail "JARVIS_CLAIM_PROGRESS=0: claim exit $rc"; fi
if [ -z "$cap" ]; then assert_pass "JARVIS_CLAIM_PROGRESS=0: no status write"; else assert_fail "JARVIS_CLAIM_PROGRESS=0 should not write status, got '$cap'"; fi

# ---------------------------------------------------------------------------
# Test 25: no per-pool progress_status → global .claim.progress_status fallback
# ---------------------------------------------------------------------------
echo "=== Test 25: claim falls back to global progress_status ==="
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS A1_WTYPE COORD_ID JARVIS_CLAIM_PROGRESS
export A1_STATUS=新建            # another start state
run_claim_prog 1086837          # no_override pool → global 处理中
unset A1_STATUS
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc status_set='$cap'"
if [ "$cap" = "处理中" ]; then assert_pass "global fallback: status advanced to 处理中"; else assert_fail "global fallback status should be 处理中, got '$cap'"; fi

# ---------------------------------------------------------------------------
# Test 26: finish status_override (4th arg) — PrWatch passes 已完成; overrides the
# resolved per-pool/global done_status. Backward compat: 3-arg finish unchanged (Test 14/15).
# ---------------------------------------------------------------------------
echo "=== Test 26: finish status_override (4th arg) wins over resolved done_status ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS A1_STATUS COORD_ID
# project 1086837 (no per-pool override) would normally resolve to global 已发布待需求排期;
# the 4th arg 已完成 must take precedence.
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 1086837 已完成 2>&1); rc=$?
cap=$(cat "$tmpstatuscap" 2>/dev/null); state=$(cat "$tmpstate" 2>/dev/null)
echo "rc=$rc status_set='$cap' tags='$state'"
if [ "$rc" = "0" ]; then assert_pass "status_override: finish exits 0"; else assert_fail "status_override: exit $rc"; fi
if [ "$cap" = "已完成" ]; then assert_pass "status_override: status set to 已完成 (override wins)"; else assert_fail "status_override: expected 已完成, got '$cap'"; fi
if printf '%s' "$state" | grep -q "jarvis-done"; then assert_pass "status_override: tagged jarvis-done"; else assert_fail "status_override: expected jarvis-done, got '$state'"; fi

# ---------------------------------------------------------------------------
# Test 27: backward compat — 3-arg finish (no override) still resolves to done_status.
# Empty/absent 4th arg must NOT hijack the resolved status.
# ---------------------------------------------------------------------------
echo "=== Test 27: 3-arg finish unchanged when no override passed ==="
printf 'jarvis-claimed' > "$tmpstate"; : > "$tmpstatuscap"
unset A1_GET_FAIL A1_UPDATE_NOOP A1_REJECT_STATUS A1_STATUS COORD_ID
out=$(PATH="$tmpbin:$PATH" JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/claim.sh" finish "$WORKITEM_ID" 1086837 2>&1); rc=$?
cap=$(cat "$tmpstatuscap" 2>/dev/null)
echo "rc=$rc status_set='$cap'"
if [ "$cap" = "已发布待需求排期" ]; then assert_pass "no-override: still resolves to global done_status (backward compatible)"; else assert_fail "no-override: expected 已发布待需求排期, got '$cap'"; fi

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
