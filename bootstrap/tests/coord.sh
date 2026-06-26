#!/usr/bin/env bash
set -uo pipefail
D="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"; COORD="$D/coord.sh"
export JARVIS_ROOT="$(mktemp -d)"; pass=0; fail=0
ck(){ [ "$2" = "$3" ] && { echo "PASS $1"; pass=$((pass+1)); } || { echo "FAIL $1: $2 != $3"; fail=$((fail+1)); }; }
id=$(bash "$COORD" register triage); ck reg-file "$([ -f "$JARVIS_ROOT/.my-day/instances/$id.json" ] && echo y)" y
ck reg-role "$(jq -r .role "$JARVIS_ROOT/.my-day/instances/$id.json")" triage
bash "$COORD" dead "$id"; ck alive-self $? 1   # self pid alive -> not dead
ck dead-missing "$(bash "$COORD" dead nohost-999999; echo $?)" 0
COORD_ID=$id bash "$COORD" checkpoint 9001 coding /wt b1 repoX; ck cp-stage "$(jq -r .stage "$JARVIS_ROOT/.my-day/tasks/9001.json")" coding; ck cp-owner "$(jq -r .owner_instance "$JARVIS_ROOT/.my-day/tasks/9001.json")" "$id"
[ "$fail" = 0 ] && echo ALLPASS || exit 1
