#!/bin/bash
# Test harness for bootstrap/log.sh: run_done + seen.

set -u

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT
export JARVIS_RUNS_DIR="$tmpdir/runs"
mkdir -p "$JARVIS_RUNS_DIR"

PASS=0
FAIL=0
pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# shellcheck source=../bootstrap/log.sh
source "$proj_root/bootstrap/log.sh"

if seen "TASK-UNKNOWN-9999"; then
    fail "seen should reject an unknown id"
else
    pass "seen rejects an unknown id"
fi

run_done "TASK-001" "processed successfully"
runs_file="$(find "$JARVIS_RUNS_DIR" -maxdepth 1 -name '*-TASK-001.md' -print -quit)"
[ -n "$runs_file" ] && pass "run_done creates a run record" || fail "run_done did not create a run record"
grep -q 'TASK-001' "$runs_file" && pass "run record contains id" || fail "run record misses id"
grep -q 'processed successfully' "$runs_file" && pass "run record contains summary" || fail "run record misses summary"
grep -q '^\*\*state:\*\* pending' "$runs_file" && pass "run_done defaults to pending" || fail "default state is not pending"
basename "$runs_file" | grep -qE '^[0-9]{4}-[0-9]{2}-[0-9]{2}-TASK-001\.md$' \
    && pass "run record filename uses UTC date" || fail "run record filename is malformed"
seen "TASK-001" && pass "seen finds a completed run" || fail "seen missed completed run"

JARVIS_RUNS_DIR="$JARVIS_RUNS_DIR" bash "$proj_root/bootstrap/log.sh" run_done "TASK-002" "merged" merged
mfile="$(find "$JARVIS_RUNS_DIR" -maxdepth 1 -name '*-TASK-002.md' -print -quit)"
grep -q '^\*\*state:\*\* merged' "$mfile" && pass "direct dispatch honors merged state" || fail "merged state was lost"
JARVIS_RUNS_DIR="$JARVIS_RUNS_DIR" bash "$proj_root/bootstrap/log.sh" seen "TASK-002" \
    && pass "direct dispatch seen finds TASK-002" || fail "direct dispatch seen missed TASK-002"

echo "PASS: $PASS  FAIL: $FAIL"
[ "$FAIL" -eq 0 ]
