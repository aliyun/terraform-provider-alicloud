#!/usr/bin/env bash
# test/pr_watch_test.sh — hermetic tests for bootstrap/pr-watch.sh (方案A PR 观察登记表)。
# JARVIS_ROOT → tmp，绝不碰真仓 .my-day。风格照 coord_test.sh。
set -uo pipefail
D="$(cd "$(dirname "${BASH_SOURCE[0]}")/../bootstrap" && pwd)"; PW="$D/pr-watch.sh"
export JARVIS_ROOT="$(mktemp -d)"; F="$JARVIS_ROOT/.my-day/bridge/pr-watch.json"
pass=0; fail=0
ck(){ [ "$2" = "$3" ] && { echo "PASS $1"; pass=$((pass+1)); } || { echo "FAIL $1: '$2' != '$3'"; fail=$((fail+1)); }; }

PR1="https://github.com/aliyun/terraform-provider-alicloud/pull/1234"
PR2="https://github.com/aliyun/terraform-provider-alicloud/pull/56"

# --- list on missing file → exit 0, empty ---
out="$(bash "$PW" list)"; ck list-missing-rc $? 0
ck list-missing-empty "$out" ""

# --- add on missing file seeds an OBJECT ({}), not an array ([]) ---
bash "$PW" add 83902495 "$PR1" 528766 >/dev/null; ck add1-rc $? 0
ck seed-is-object "$(jq -r 'type' "$F")" object
ck add1-has-key "$(jq -r 'has("83902495")' "$F")" true
ck add1-pr_url "$(jq -r '."83902495".pr_url' "$F")" "$PR1"
ck add1-project "$(jq -r '."83902495".project' "$F")" 528766
ck add1-has-ts "$(jq -r '."83902495" | has("submitted_at")' "$F")" true

# --- second add: object grows to two keys ---
bash "$PW" add 99 "$PR2" 1086837 >/dev/null; ck add2-rc $? 0
ck two-keys "$(jq -r 'keys | length' "$F")" 2

# --- list emits TSV rows (<ticket>\t<pr_url>\t<project>) ---
rows="$(bash "$PW" list | sort)"
ck list-line-99 "$(printf '%s\n' "$rows" | grep -c $'^99\t'"$PR2"$'\t1086837$')" 1
ck list-line-main "$(printf '%s\n' "$rows" | grep -c "^83902495"$'\t'"$PR1"$'\t528766$')" 1
# TSV = exactly 3 tab-separated columns per row
ck list-tsv-cols "$(bash "$PW" list | awk -F'\t' 'NF!=3{print "bad"}' | grep -c bad)" 0

# --- remove deletes the key (and only that key) ---
bash "$PW" remove 99 >/dev/null; ck rm-rc $? 0
ck rm-key-gone "$(jq -r 'has("99")' "$F")" false
ck rm-other-kept "$(jq -r 'has("83902495")' "$F")" true
ck rm-still-object "$(jq -r 'type' "$F")" object

# --- url validation: bare number (would be mis-resolved by gh to the wrong repo) → non-zero ---
bash "$PW" add 1 12345 proj >/dev/null 2>&1; ck bad-url-bare-number-rc "$([ $? -ne 0 ] && echo nonzero)" nonzero
ck bad-url-not-registered "$(jq -r 'has("1")' "$F")" false
# a non-github host is also rejected
bash "$PW" add 1 "https://gitlab.com/o/r/pull/5" proj >/dev/null 2>&1; ck bad-url-host-rc "$([ $? -ne 0 ] && echo nonzero)" nonzero
# a valid full github PR url is accepted
bash "$PW" add 2 "https://github.com/o/r/pull/7" proj >/dev/null 2>&1; ck good-url-rc $? 0

# --- missing args → non-zero usage ---
bash "$PW" add onlyone >/dev/null 2>&1; ck add-missing-args "$([ $? -ne 0 ] && echo nonzero)" nonzero
bash "$PW" >/dev/null 2>&1; ck no-subcmd "$([ $? -ne 0 ] && echo nonzero)" nonzero

# --- concurrency: parallel add + remove (mkdir-lock) → file stays valid JSON, no lost/corrupt ---
# fresh registry seeded with id0; then race add id1/id2/id3 against remove id0.
CJ="$JARVIS_ROOT/.my-day/bridge/pr-watch.json"; rm -f "$CJ"
bash "$PW" add id0 "https://github.com/o/r/pull/1" p >/dev/null
bash "$PW" add id1 "https://github.com/o/r/pull/2" p >/dev/null &
bash "$PW" add id2 "https://github.com/o/r/pull/3" p >/dev/null &
bash "$PW" add id3 "https://github.com/o/r/pull/4" p >/dev/null &
bash "$PW" remove id0 >/dev/null &
wait
jq . "$CJ" >/dev/null 2>&1; ck concurrent-valid-json $? 0
ck concurrent-id1 "$(jq -r 'has("id1")' "$CJ")" true
ck concurrent-id2 "$(jq -r 'has("id2")' "$CJ")" true
ck concurrent-id3 "$(jq -r 'has("id3")' "$CJ")" true
ck concurrent-id0-removed "$(jq -r 'has("id0")' "$CJ")" false
# lock dir must not be left behind after all ops settle
ck concurrent-no-stale-lock "$([ -e "$JARVIS_ROOT/.my-day/bridge/.pr-watch.lock" ] && echo present || echo clean)" clean

echo ""
echo "=== Summary: PASS=$pass FAIL=$fail ==="
[ "$fail" = 0 ] && { echo "All tests passed"; exit 0; } || { echo "TESTS FAILED"; exit 1; }
