#!/bin/bash
# Test harness for bootstrap/plan.sh
# Tests: scan→filter seen→plan file created, supervised exit 2, no a1 write mutation

set -u

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

tmpbin=$(mktemp -d)
tmpruns=$(mktemp -d)
trap 'rm -rf "$tmpbin" "$tmpruns"' EXIT

export JARVIS_RUNS_DIR="$tmpruns"

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
# Stub scan.sh: emits 2 items — WI-100 (high conf) and WI-200 (low conf)
# WI-200 is already in runs/ (seen), so only WI-100 should appear in the plan
# ---------------------------------------------------------------------------
cat > "$tmpbin/scan.sh" << 'STUB'
#!/bin/bash
echo '[{"id":"WI-100","title":"Fix login bug","type":"bug","status":"open"},{"id":"WI-200","title":"Old item already done","type":"task","status":"closed"}]'
STUB
chmod +x "$tmpbin/scan.sh"

# Stub a1: must NOT be called for any write command
cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
# Record any call for later assertion
echo "$@" >> "$JARVIS_RUNS_DIR/a1_calls.log"
# If called with a write-like subcommand, flag it
case "${2:-}${3:-}" in
    workitemupdate|commentcreate|workitemcreate|replycomment)
        echo "ERROR: a1 write command invoked: $*" >&2
        exit 99
        ;;
esac
# auth whoami — stub it for any potential use
if [ "${1:-}" = "auth" ] && [ "${2:-}" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"

export PATH="$tmpbin:$PATH"

# Pre-seed WI-200 as already seen (simulate a prior run_done)
touch "$tmpruns/2026-01-01-WI-200.md"

# ---------------------------------------------------------------------------
# Test 1: supervised mode — plan file created, exit code 2
# ---------------------------------------------------------------------------
echo "=== Test 1: supervised mode creates plan and exits 2 ==="

output=$(JARVIS_RUNS_DIR="$tmpruns" bash "$proj_root/bootstrap/plan.sh" 2>&1)
exit_code=$?

echo "Exit code: $exit_code"

if [ "$exit_code" -eq 2 ]; then
    assert_pass "supervised mode exits with code 2 (awaiting authorisation)"
else
    assert_fail "supervised mode should exit 2, got $exit_code"
fi

# ---------------------------------------------------------------------------
# Test 2: plan file written to runs/
# ---------------------------------------------------------------------------
echo "=== Test 2: plan file created in runs/ ==="

plan_count=$(ls "$tmpruns"/plan-*.md 2>/dev/null | wc -l | tr -d ' ')
if [ "$plan_count" -ge 1 ]; then
    assert_pass "plan file created in runs/"
else
    assert_fail "no plan-*.md found in $tmpruns (count=$plan_count)"
fi

# ---------------------------------------------------------------------------
# Test 3: unseen item (WI-100) appears in plan
# ---------------------------------------------------------------------------
echo "=== Test 3: unseen item WI-100 appears in plan ==="

plan_file=$(ls "$tmpruns"/plan-*.md 2>/dev/null | head -1)
if grep -q "WI-100" "$plan_file" 2>/dev/null; then
    assert_pass "unseen item WI-100 is in the plan"
else
    assert_fail "WI-100 not found in plan: $plan_file"
fi

# ---------------------------------------------------------------------------
# Test 4: seen item (WI-200) does NOT appear in plan
# ---------------------------------------------------------------------------
echo "=== Test 4: seen item WI-200 excluded from plan ==="

if grep -q "WI-200" "$plan_file" 2>/dev/null; then
    assert_fail "WI-200 should be excluded (already seen) but found in plan"
else
    assert_pass "WI-200 correctly excluded from plan"
fi

# ---------------------------------------------------------------------------
# Test 5: plan file contains authorisation prompt
# ---------------------------------------------------------------------------
echo "=== Test 5: plan file contains 待授权 prompt ==="

if grep -q "待授权" "$plan_file" 2>/dev/null; then
    assert_pass "plan file contains 待授权 authorisation prompt"
else
    assert_fail "plan file missing 待授权 prompt"
fi

# ---------------------------------------------------------------------------
# Test 6: stdout also contains 待授权 prompt in supervised mode
# ---------------------------------------------------------------------------
echo "=== Test 6: stdout contains 待授权 prompt ==="

if echo "$output" | grep -q "待授权"; then
    assert_pass "stdout contains 待授权 authorisation prompt"
else
    assert_fail "stdout missing 待授权 prompt"
fi

# ---------------------------------------------------------------------------
# Test 7: no a1 write commands were invoked
# ---------------------------------------------------------------------------
echo "=== Test 7: no a1 write mutations ==="

if grep -qE "workitem update|comment create|workitem create|reply" "$tmpruns/a1_calls.log" 2>/dev/null; then
    assert_fail "a1 write command was invoked (mutation detected)"
else
    assert_pass "no a1 write mutations detected"
fi

# ---------------------------------------------------------------------------
# Test 8: plan contains action/confidence/auto-or-stop fields
# ---------------------------------------------------------------------------
echo "=== Test 8: plan file contains required columns (id, title, action, confidence) ==="

if grep -qiE "WI-100" "$plan_file" 2>/dev/null && \
   grep -qiE "(low_conf|high_conf|stop|auto|escalate)" "$plan_file" 2>/dev/null; then
    assert_pass "plan file has id and confidence markers"
else
    assert_fail "plan file missing required columns — check action/confidence fields"
fi

# ---------------------------------------------------------------------------
# Test 9: unattended mode exits 0 (no wait prompt needed)
# ---------------------------------------------------------------------------
echo "=== Test 9: unattended mode exits 0 ==="

# Stub autonomy.md in tmp that says unattended
tmpaut=$(mktemp -d)
cat > "$tmpaut/autonomy.md" << 'EOF'
```json
{"mode":"unattended","auto":["reply","create_req","tag","create_cr","worktree","prestage"],"stop":["release_prod"],"escalate_if":["low_conf","verify_fail","redline"]}
```
EOF

JARVIS_RUNS_DIR="$tmpruns" JARVIS_AUTONOMY_FILE="$tmpaut/autonomy.md" \
    bash "$proj_root/bootstrap/plan.sh" > /dev/null 2>&1
exit_code_u=$?

if [ "$exit_code_u" -eq 0 ]; then
    assert_pass "unattended mode exits 0"
else
    assert_fail "unattended mode should exit 0, got $exit_code_u"
fi

rm -rf "$tmpaut"

# ---------------------------------------------------------------------------
# Test 10: plan.sh source contains NO a1 write verbs
# ---------------------------------------------------------------------------
echo "=== Test 10: plan.sh source is read-only (no a1 write verbs) ==="

write_hits=$(grep -E "a1 (project|.* )(create|update|comment|relation|release)" \
    "$proj_root/bootstrap/plan.sh" 2>/dev/null || true)

if [ -z "$write_hits" ]; then
    assert_pass "plan.sh source contains no a1 write verbs"
else
    assert_fail "plan.sh source contains a1 write verbs: $write_hits"
fi

# ---------------------------------------------------------------------------
# Test 11: garbage mode → exit 2 (default-deny)
# ---------------------------------------------------------------------------
echo "=== Test 11: garbage mode exits 2 (default-deny) ==="

tmpgarbage=$(mktemp -d)
cat > "$tmpgarbage/autonomy.md" << 'EOF'
```json
{"mode":"GARBAGE_VALUE","auto":[],"stop":[]}
```
EOF

JARVIS_RUNS_DIR="$tmpruns" JARVIS_AUTONOMY_FILE="$tmpgarbage/autonomy.md" \
    bash "$proj_root/bootstrap/plan.sh" > /dev/null 2>&1
exit_code_g=$?

if [ "$exit_code_g" -eq 2 ]; then
    assert_pass "garbage mode exits 2 (default-deny gate)"
else
    assert_fail "garbage mode should exit 2, got $exit_code_g"
fi

rm -rf "$tmpgarbage"

# ---------------------------------------------------------------------------
# Test 12: stdin JSON → plan shows only unseen items, scan.sh not called
# pipe 3 items; WI-300 already seen; expect WI-400 and WI-500 in plan, not WI-300
# ---------------------------------------------------------------------------
echo "=== Test 12: stdin JSON consumed; scan.sh NOT invoked; 2 of 3 items shown ==="

# Fresh runs dir for isolation
tmpruns2=$(mktemp -d)
trap 'rm -rf "$tmpruns2"' EXIT

# Pre-seed WI-300 as already seen
touch "$tmpruns2/2026-01-01-WI-300.md"

# Intercept scan.sh: write a marker file if invoked from within plan.sh
cat > "$tmpbin/scan.sh" << 'STUB'
#!/bin/bash
echo "SCAN_CALLED" >> "${JARVIS_RUNS_DIR:-/tmp}/scan_called.marker"
echo '[{"id":"SCAN-ITEM","title":"Should not appear","type":"bug","status":"open"}]'
STUB
chmod +x "$tmpbin/scan.sh"

stdin_json='[{"id":"WI-300","title":"Already done","type":"task","status":"closed"},{"id":"WI-400","title":"New bug","type":"bug","status":"open"},{"id":"WI-500","title":"New req","type":"req","status":"open"}]'

stdin_output=$(echo "$stdin_json" | JARVIS_RUNS_DIR="$tmpruns2" bash "$proj_root/bootstrap/plan.sh" 2>&1)
stdin_exit=$?

# scan.sh must NOT have been called
if [ -f "$tmpruns2/scan_called.marker" ]; then
    assert_fail "scan.sh was invoked despite stdin being provided"
else
    assert_pass "scan.sh not invoked when stdin JSON provided"
fi

# WI-400 and WI-500 should appear in the plan output
if echo "$stdin_output" | grep -q "WI-400" && echo "$stdin_output" | grep -q "WI-500"; then
    assert_pass "unseen items WI-400 and WI-500 appear in stdin-driven plan"
else
    assert_fail "expected WI-400 and WI-500 in plan but got: $(echo "$stdin_output" | head -5)"
fi

# WI-300 (seen) must NOT appear
if echo "$stdin_output" | grep -q "WI-300"; then
    assert_fail "seen item WI-300 should be excluded from stdin-driven plan"
else
    assert_pass "seen item WI-300 correctly excluded from stdin-driven plan"
fi

# exactly 2 new items (not 3)
item_count_in_plan=$(echo "$stdin_output" | grep -cE "WI-(400|500)" || true)
if [ "$item_count_in_plan" -eq 2 ]; then
    assert_pass "exactly 2 new items shown (not 3)"
else
    assert_fail "expected 2 new items in plan, found $item_count_in_plan"
fi

rm -rf "$tmpruns2"

# ---------------------------------------------------------------------------
# Test 13: plan output includes a priority column
# ---------------------------------------------------------------------------
echo "=== Test 13: plan file/output contains priority column ==="

plan_file_main=$(ls "$tmpruns"/plan-*.md 2>/dev/null | head -1)
if grep -qi "priority\|优先" "$plan_file_main" 2>/dev/null; then
    assert_pass "plan file contains priority column"
else
    assert_fail "plan file missing priority column"
fi

# ---------------------------------------------------------------------------
# Test 14: batch dedup — no per-item bash log.sh fork (strace/ps not available on macOS;
# verify indirectly: plan.sh source must not contain 'bash.*log.sh seen')
# ---------------------------------------------------------------------------
echo "=== Test 14: plan.sh source does not fork bash log.sh seen per-item ==="

if grep -qE 'bash.*log\.sh.*seen' "$proj_root/bootstrap/plan.sh" 2>/dev/null; then
    assert_fail "plan.sh still contains per-item 'bash log.sh seen' fork"
else
    assert_pass "no per-item bash log.sh seen fork in plan.sh source"
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
