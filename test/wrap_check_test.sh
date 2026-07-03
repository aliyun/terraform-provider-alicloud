#!/usr/bin/env bash
# test/wrap_check_test.sh — unit tests for wrap-check.sh
#
# Scenarios:
#   1. done:false id with no runs/ file → exit 2
#   2. done:false id with a runs/ file → exit 0
#   3. done:true id with no runs/ file → exit 0 (done:true is exempt)
#   4. no claims file → exit 0
#
# Run: bash test/wrap_check_test.sh
# Prints PASS and exits 0 on success; prints FAIL and exits 1 on any failure.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/../bootstrap" && pwd)"

WRAP_CHECK="$BOOTSTRAP_DIR/wrap-check.sh"

# ---------------------------------------------------------------------------
# Fake a1 stub (must be in PATH to prevent any accidental a1 call)
# ---------------------------------------------------------------------------
FAKE_BIN_DIR="$(mktemp -d)"
cat > "$FAKE_BIN_DIR/a1" <<'EOF'
#!/usr/bin/env bash
echo "a1: should not be called in wrap-check tests" >&2
exit 99
EOF
chmod +x "$FAKE_BIN_DIR/a1"
export PATH="$FAKE_BIN_DIR:$PATH"

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
pass_count=0
fail_count=0

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

make_jarvis_root() {
    local root
    root="$(mktemp -d)"
    mkdir -p "$root/.my-day" "$root/runs"
    # minimal config/pools.json so log.sh/wrap-check.sh can find it
    mkdir -p "$root/config"
    cp "$BOOTSTRAP_DIR/../config/pools.json" "$root/config/pools.json"
    echo "$root"
}

today="$(date -u +%F)"

# ---------------------------------------------------------------------------
# Declare all temp roots up front so the trap never references an unset variable
# ---------------------------------------------------------------------------
ROOT1=""
ROOT2=""
ROOT3=""
ROOT4=""
ROOT5=""
ROOT6=""
ROOT7=""
ROOT8=""
ROOT9=""
ROOT10=""
ROOT11=""
ROOT12=""
ROOT13=""
trap 'rm -rf "$FAKE_BIN_DIR" "$ROOT1" "$ROOT2" "$ROOT3" "$ROOT4" "$ROOT5" "$ROOT6" "$ROOT7" "$ROOT8" "$ROOT9" "$ROOT10" "$ROOT11" "$ROOT12" "$ROOT13" 2>/dev/null; exit' INT TERM EXIT

# ---------------------------------------------------------------------------
# Test 1: done:false id, no runs/ file → exit 2
# ---------------------------------------------------------------------------
echo "Test 1: done:false + no run file → exit 2"
ROOT1="$(make_jarvis_root)"

claims_file="$ROOT1/.my-day/claims-${today}.json"
printf '[{"id":"WI-T1","done":false}]\n' > "$claims_file"

assert_exit_code "missing run file → exit 2" 2 \
    env JARVIS_ROOT="$ROOT1" JARVIS_RUNS_DIR="$ROOT1/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 2: done:false id, runs/ file present → exit 0
# ---------------------------------------------------------------------------
echo "Test 2: done:false + run file present → exit 0"
ROOT2="$(make_jarvis_root)"
claims_file2="$ROOT2/.my-day/claims-${today}.json"
printf '[{"id":"WI-T2","done":false}]\n' > "$claims_file2"
# Create the run file that log.sh seen would match
touch "$ROOT2/runs/${today}-WI-T2.md"

assert_exit_code "run file present → exit 0" 0 \
    env JARVIS_ROOT="$ROOT2" JARVIS_RUNS_DIR="$ROOT2/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 3: done:true id, no runs/ file → still exit 0 (done items exempt)
# ---------------------------------------------------------------------------
echo "Test 3: done:true + no run file → exit 0"
ROOT3="$(make_jarvis_root)"
claims_file3="$ROOT3/.my-day/claims-${today}.json"
printf '[{"id":"WI-T3","done":true}]\n' > "$claims_file3"

assert_exit_code "done:true exempt → exit 0" 0 \
    env JARVIS_ROOT="$ROOT3" JARVIS_RUNS_DIR="$ROOT3/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 4: no claims file at all → exit 0
# ---------------------------------------------------------------------------
echo "Test 4: no claims file → exit 0"
ROOT4="$(make_jarvis_root)"
# don't write any claims file

assert_exit_code "no claims file → exit 0" 0 \
    env JARVIS_ROOT="$ROOT4" JARVIS_RUNS_DIR="$ROOT4/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 5: claim in YESTERDAY's file, unclosed, no runs/ file → exit 2
# (midnight-orphan edge case: wrap-check must scan all claims-*.json, not just today's)
# ---------------------------------------------------------------------------
echo "Test 5: yesterday's claims file, done:false + no run file → exit 2"
ROOT5="$(make_jarvis_root)"
yesterday="$(date -u -v-1d +%F 2>/dev/null || date -u --date='1 day ago' +%F)"
yesterday_claims="$ROOT5/.my-day/claims-${yesterday}.json"
printf '[{"id":"WI-T5","done":false}]\n' > "$yesterday_claims"
# No runs/ file for WI-T5 — should trigger exit 2

assert_exit_code "yesterday unclosed claim → exit 2" 2 \
    env JARVIS_ROOT="$ROOT5" JARVIS_RUNS_DIR="$ROOT5/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# assert_stderr_contains — runs cmd, asserts stderr contains a substring
# ---------------------------------------------------------------------------
assert_stderr_contains() {
    local desc="$1"
    local needle="$2"
    shift 2
    local err
    err="$("$@" 2>&1 1>/dev/null)" || true
    if printf '%s' "$err" | grep -qF -- "$needle"; then
        echo "  PASS: $desc"
        pass_count=$((pass_count + 1))
    else
        echo "  FAIL: $desc (stderr missing '$needle': $err)"
        fail_count=$((fail_count + 1))
    fi
}

# ===========================================================================
# Owner-scoping matrix (cap-claim-ledger-owner-scoping Phase 1)
# self is set via COORD_ID; a claim's owner scopes whether the Stop hook blocks.
# ===========================================================================
SELF="inst-self-1"
OTHER="inst-other-2"

# ---------------------------------------------------------------------------
# Test 6: owner == self, no run_done → block (exit 2)
# ---------------------------------------------------------------------------
echo "Test 6: owner==self + no run file → exit 2 (own claim, blocks)"
ROOT6="$(make_jarvis_root)"
printf '[{"id":"WI-T6","done":false,"owner":"%s"}]\n' "$SELF" \
    > "$ROOT6/.my-day/claims-${today}.json"
assert_exit_code "own claim, no run_done → exit 2" 2 \
    env COORD_ID="$SELF" JARVIS_ROOT="$ROOT6" JARVIS_RUNS_DIR="$ROOT6/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 7: owner == self, has run_done → pass (exit 0)
# ---------------------------------------------------------------------------
echo "Test 7: owner==self + run file present → exit 0"
ROOT7="$(make_jarvis_root)"
printf '[{"id":"WI-T7","done":false,"owner":"%s"}]\n' "$SELF" \
    > "$ROOT7/.my-day/claims-${today}.json"
touch "$ROOT7/runs/${today}-WI-T7.md"
assert_exit_code "own claim, run_done present → exit 0" 0 \
    env COORD_ID="$SELF" JARVIS_ROOT="$ROOT7" JARVIS_RUNS_DIR="$ROOT7/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 8: owner == other instance, no run_done → skip (exit 0) + WARN on stderr
# ---------------------------------------------------------------------------
echo "Test 8: owner==other + no run file → exit 0 (skip) + WARN"
ROOT8="$(make_jarvis_root)"
printf '[{"id":"WI-T8","done":false,"owner":"%s"}]\n' "$OTHER" \
    > "$ROOT8/.my-day/claims-${today}.json"
assert_exit_code "foreign claim, no run_done → skip exit 0" 0 \
    env COORD_ID="$SELF" JARVIS_ROOT="$ROOT8" JARVIS_RUNS_DIR="$ROOT8/runs" bash "$WRAP_CHECK"
assert_stderr_contains "foreign claim emits skip/WARN line" "skip WI-T8" \
    env COORD_ID="$SELF" JARVIS_ROOT="$ROOT8" JARVIS_RUNS_DIR="$ROOT8/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 9: owner empty (legacy entry), no run_done → block (exit 2, no regression)
# ---------------------------------------------------------------------------
echo "Test 9: owner empty/legacy + no run file → exit 2 (no regression)"
ROOT9="$(make_jarvis_root)"
# legacy shape: no owner field at all
printf '[{"id":"WI-T9","done":false}]\n' \
    > "$ROOT9/.my-day/claims-${today}.json"
assert_exit_code "legacy ownerless claim, no run_done → exit 2" 2 \
    env COORD_ID="$SELF" JARVIS_ROOT="$ROOT9" JARVIS_RUNS_DIR="$ROOT9/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 10: touched id whose claim owner is a foreign instance → skip (exit 0)
# (touched-*.json has no owner field; owner is resolved from the claims map)
# ---------------------------------------------------------------------------
echo "Test 10: touched id + foreign claim owner → exit 0 (skip via claims map)"
ROOT10="$(make_jarvis_root)"
# claim closed (done:true) but still records the foreign owner; touched re-surfaces id
printf '[{"id":"WI-T10","done":true,"owner":"%s"}]\n' "$OTHER" \
    > "$ROOT10/.my-day/claims-${today}.json"
printf '["WI-T10"]\n' > "$ROOT10/.my-day/touched-${today}.json"
assert_exit_code "touched id owned by other instance → skip exit 0" 0 \
    env COORD_ID="$SELF" JARVIS_ROOT="$ROOT10" JARVIS_RUNS_DIR="$ROOT10/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 11: touched id with no claim entry (owner unknown="") → block (exit 2)
# ---------------------------------------------------------------------------
echo "Test 11: touched id, no claim/owner + no run file → exit 2 (treated as ownerless)"
ROOT11="$(make_jarvis_root)"
printf '["WI-T11"]\n' > "$ROOT11/.my-day/touched-${today}.json"
assert_exit_code "touched ownerless id, no run_done → exit 2" 2 \
    env COORD_ID="$SELF" JARVIS_ROOT="$ROOT11" JARVIS_RUNS_DIR="$ROOT11/runs" bash "$WRAP_CHECK"

# ===========================================================================
# D2: interactive sessions derive owner from CLAUDE_CODE_SESSION_ID (cc-<sid>)
# via coord_self(). Two different sessions get distinct owners and must not block
# each other; a session is still held to account for its own claims.
# ===========================================================================

# ---------------------------------------------------------------------------
# Test 12: claim owned by session A; session B (no COORD_ID, different sid) → skip (0)
# ---------------------------------------------------------------------------
echo "Test 12: interactive — other session's claim → exit 0 (skip)"
ROOT12="$(make_jarvis_root)"
printf '[{"id":"WI-T12","done":false,"owner":"cc-sess-A"}]\n' \
    > "$ROOT12/.my-day/claims-${today}.json"
assert_exit_code "interactive foreign-session claim → skip exit 0" 0 \
    env CLAUDE_CODE_SESSION_ID="sess-B" JARVIS_ROOT="$ROOT12" JARVIS_RUNS_DIR="$ROOT12/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 13: session A runs wrap-check on its OWN claim, no run_done → block (exit 2)
# ---------------------------------------------------------------------------
echo "Test 13: interactive — own session claim, no run file → exit 2 (no regression)"
ROOT13="$(make_jarvis_root)"
printf '[{"id":"WI-T13","done":false,"owner":"cc-sess-A"}]\n' \
    > "$ROOT13/.my-day/claims-${today}.json"
assert_exit_code "interactive own-session claim, no run_done → exit 2" 2 \
    env CLAUDE_CODE_SESSION_ID="sess-A" JARVIS_ROOT="$ROOT13" JARVIS_RUNS_DIR="$ROOT13/runs" bash "$WRAP_CHECK"

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
