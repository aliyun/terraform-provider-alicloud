#!/usr/bin/env bash
# bootstrap/tests/watchdog.sh — unit tests for bootstrap/reconcile.sh orphan(原 watchdog.sh)
#
# Sets JARVIS_ROOT + JARVIS_ESCALATION_DIR to temp dirs, creates a task owned
# by dead instance "nohost-999999", runs reconcile.sh orphan, and asserts that
# an escalation file is written for the orphaned task.
#
# Run: bash bootstrap/tests/watchdog.sh
# Prints PASS/FAIL per assertion; exits 0 on all-pass, 1 on any failure.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

WATCHDOG="$BOOTSTRAP_DIR/reconcile.sh"

# ── Isolated temp roots ──────────────────────────────────────────────────────
export JARVIS_ROOT
JARVIS_ROOT="$(mktemp -d)"
export JARVIS_ESCALATION_DIR
JARVIS_ESCALATION_DIR="$(mktemp -d)"

trap 'rm -rf "$JARVIS_ROOT" "$JARVIS_ESCALATION_DIR"' EXIT

pass=0; fail=0
ck() {
    if [ "$2" = "$3" ]; then
        echo "PASS $1"; pass=$((pass+1))
    else
        echo "FAIL $1: got='$2' want='$3'"; fail=$((fail+1))
    fi
}

# ── Create a task owned by a dead instance ───────────────────────────────────
# "nohost-999999" is dead: coord.sh dead checks for instances/nohost-999999.hb
# which does not exist → exits 0 (dead).
AID="watchdog-test-aid-1"
DEAD_OWNER="nohost-999999"

mkdir -p "$JARVIS_ROOT/.my-day/tasks"
printf '{"aone_id":"%s","owner_instance":"%s","stage":"coding","worktree":"","branch":"","repo":"","updated":"2024-01-01T00:00:00Z"}' \
    "$AID" "$DEAD_OWNER" > "$JARVIS_ROOT/.my-day/tasks/$AID.json"

# ── Run watchdog ─────────────────────────────────────────────────────────────
bash "$WATCHDOG" orphan

# ── Assert escalation file was created ──────────────────────────────────────
esc_file="$JARVIS_ESCALATION_DIR/$AID.md"
ck "escalation-file-exists" "$([ -f "$esc_file" ] && echo y || echo n)" y

# Verify content mentions the expected reason
if [ -f "$esc_file" ]; then
    ck "escalation-reason-present" \
        "$(grep -c 'owner dead' "$esc_file" || echo 0)" 1
fi

# ── I4: dedup — second run must NOT append a duplicate escalation ─────────────
AID2="watchdog-test-aid-2"
mkdir -p "$JARVIS_ROOT/.my-day/tasks"
printf '{"aone_id":"%s","owner_instance":"%s","stage":"coding","worktree":"","branch":"","repo":"","updated":"2024-01-01T00:00:00Z"}' \
    "$AID2" "$DEAD_OWNER" > "$JARVIS_ROOT/.my-day/tasks/$AID2.json"

bash "$WATCHDOG" orphan   # first run — creates escalation file
esc2="$JARVIS_ESCALATION_DIR/$AID2.md"
size_after_first=$(wc -c < "$esc2" | tr -d ' ')

bash "$WATCHDOG" orphan   # second run — must be a no-op (file already exists)
size_after_second=$(wc -c < "$esc2" | tr -d ' ')

ck "dedup-no-append" "$size_after_first" "$size_after_second"

# ── Summary ──────────────────────────────────────────────────────────────────
echo ""
echo "Results: $pass passed, $fail failed"
[ "$fail" = 0 ] && echo ALLPASS && exit 0 || exit 1
