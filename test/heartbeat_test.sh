#!/usr/bin/env bash
# test/heartbeat_test.sh — unit tests for bootstrap/heartbeat.sh
#
# Registers an instance, starts a short-lived follow_pid, runs heartbeat sidecar
# with HB_INT=1, verifies hb file mtime advances, then kills follow_pid and
# asserts the heartbeat loop exits within 2 seconds.
#
# Run: bash test/heartbeat_test.sh
# Prints PASS/FAIL per assertion; exits 0 on all-pass, 1 on any failure.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/../bootstrap" && pwd)"

COORD="$BOOTSTRAP_DIR/coord.sh"
HB="$BOOTSTRAP_DIR/heartbeat.sh"

export JARVIS_ROOT
JARVIS_ROOT="$(mktemp -d)"

follow_pid=0
hb_pid=0

cleanup() {
    [ "$follow_pid" -gt 0 ] && kill "$follow_pid" 2>/dev/null || true
    [ "$hb_pid"    -gt 0 ] && kill "$hb_pid"    2>/dev/null || true
    wait 2>/dev/null || true
    rm -rf "$JARVIS_ROOT"
}
trap cleanup EXIT

pass=0; fail=0
ck() {
    if [ "$2" = "$3" ]; then
        echo "PASS $1"; pass=$((pass+1))
    else
        echo "FAIL $1: got='$2' want='$3'"; fail=$((fail+1))
    fi
}

# ── Register an instance ────────────────────────────────────────────────────
id=$(bash "$COORD" register triage)
hb_file="$JARVIS_ROOT/.my-day/instances/$id.hb"

# Record initial mtime of the hb file (created by register)
mtime_before=$(stat -f %m "$hb_file" 2>/dev/null || stat -c %Y "$hb_file")

# ── Start a short-lived process to follow ───────────────────────────────────
sleep 5 & follow_pid=$!

# ── Start heartbeat sidecar with HB_INT=1 ───────────────────────────────────
HB_INT=1 bash "$HB" "$id" "$follow_pid" &
hb_pid=$!

# ── Wait for at least one heartbeat cycle ───────────────────────────────────
sleep 2

# ── Assert hb file mtime advanced ───────────────────────────────────────────
mtime_after=$(stat -f %m "$hb_file" 2>/dev/null || stat -c %Y "$hb_file")
ck "hb-mtime-advanced" "$([ "$mtime_after" -gt "$mtime_before" ] && echo y || echo n)" y

# ── Kill the followed process ────────────────────────────────────────────────
kill "$follow_pid" 2>/dev/null || true
wait "$follow_pid" 2>/dev/null || true
follow_pid=0   # prevent double-kill in cleanup

# ── Wait for heartbeat loop to notice and exit ───────────────────────────────
sleep 2

# ── Assert heartbeat loop exited ─────────────────────────────────────────────
# Use ps state to distinguish zombie/gone (done) from still-running
state=$(ps -p "$hb_pid" -o state= 2>/dev/null | head -1 | tr -d '[:space:]')
if [ -z "$state" ] || [[ "$state" == Z* ]]; then
    hb_alive=n   # process gone or zombie — effectively exited
else
    hb_alive=y   # still running
fi
ck "hb-exited" "$hb_alive" n

wait "$hb_pid" 2>/dev/null || true
hb_pid=0   # prevent double-kill in cleanup

# ── Summary ──────────────────────────────────────────────────────────────────
echo ""
echo "Results: $pass passed, $fail failed"
[ "$fail" = 0 ] && echo ALLPASS && exit 0 || exit 1
