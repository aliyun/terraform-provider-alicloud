#!/usr/bin/env bash
set -uo pipefail
D="$(cd "$(dirname "${BASH_SOURCE[0]}")/../bootstrap" && pwd)"; COORD="$D/coord.sh"
export JARVIS_ROOT="$(mktemp -d)"; pass=0; fail=0
ck(){ [ "$2" = "$3" ] && { echo "PASS $1"; pass=$((pass+1)); } || { echo "FAIL $1: $2 != $3"; fail=$((fail+1)); }; }

# Test that register creates the right files (uses coord.sh's own $$, pid exits after)
reg_id=$(bash "$COORD" register triage)
ck reg-file "$([ -f "$JARVIS_ROOT/.my-day/instances/$reg_id.json" ] && echo y)" y
ck reg-role "$(jq -r .role "$JARVIS_ROOT/.my-day/instances/$reg_id.json")" triage

# alive-self: use the test script's OWN PID (stays alive for the whole test).
# coord.sh register embeds its subshell PID which exits; create hb manually instead.
live_id="$(hostname)-$$"
: > "$JARVIS_ROOT/.my-day/instances/$live_id.hb"
bash "$COORD" dead "$live_id"; ck alive-self $? 1   # self pid alive -> not dead

ck dead-missing "$(bash "$COORD" dead nohost-999999; echo $?)" 0

# Checkpoint with live_id so 9001 is NOT an orphan during orphan-list test
COORD_ID=$live_id bash "$COORD" checkpoint 9001 coding /wt b1 repoX
ck cp-stage "$(jq -r .stage "$JARVIS_ROOT/.my-day/tasks/9001.json")" coding
ck cp-owner "$(jq -r .owner_instance "$JARVIS_ROOT/.my-day/tasks/9001.json")" "$live_id"

# Only 9002 (dead owner) should appear; 9001 (live owner) must not
COORD_ID=nohost-999999 bash "$COORD" checkpoint 9002 coding
ck orphan-list "$(bash "$COORD" list-orphans)" 9002

COORD_ID=$live_id bash "$COORD" adopt 9002
ck adopt-owner "$(jq -r .owner_instance "$JARVIS_ROOT/.my-day/tasks/9002.json")" "$live_id"

# C1: empty owner_instance → task is invisible to list-orphans (does NOT surface)
COORD_ID="" bash "$COORD" checkpoint 9003 coding /wt b1 repoX
ck empty-owner-invisible "$(bash "$COORD" list-orphans | grep -c 9003 || true)" 0

# C1: dead-owner task IS surfaced by list-orphans (non-empty dead owner)
COORD_ID=nohost-999999 bash "$COORD" checkpoint 9004 coding /wt b1 repoX
ck dead-owner-surfaces "$(bash "$COORD" list-orphans | grep -c 9004 || true)" 1

# I2: stale hb file (mtime > TTL) → dead even if pid happened to be recycled
_stale_id="stalehost-$$"
mkdir -p "$JARVIS_ROOT/.my-day/instances"
: > "$JARVIS_ROOT/.my-day/instances/$_stale_id.hb"
touch -t 197001010000 "$JARVIS_ROOT/.my-day/instances/$_stale_id.hb" 2>/dev/null \
  || touch -d "1970-01-01 00:00:00" "$JARVIS_ROOT/.my-day/instances/$_stale_id.hb" 2>/dev/null || true
bash "$COORD" dead "$_stale_id"; ck dead-stale-hb $? 0   # ancient hb → dead regardless of pid

[ "$fail" = 0 ] && echo ALLPASS || exit 1
