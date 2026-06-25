#!/usr/bin/env bash
# bootstrap/week.sh — 周回顾工具
#
# Cross-pool, lists jarvis-done workitems modified in the last N days (default 7)
# and produces a markdown checklist with Aone links, then cross-checks against
# local runs/<date>-<id>.md ledger to surface drift.
#
# Usage:
#   bash bootstrap/week.sh            # last 7 days
#   bash bootstrap/week.sh --days 14  # last 14 days
#
# Environment overrides:
#   JARVIS_ROOT       — repo root (default: git rev-parse --show-toplevel)
#   JARVIS_RUNS_DIR   — runs dir (default: <JARVIS_ROOT>/runs)
#   JARVIS_WEEK_DAYS  — window in days (default 7; --days wins over env)
#
# Read-only on Aone; only reads local runs/. Pure bash + jq + a1, no network deps.

set -uo pipefail

# ---------------------------------------------------------------------------
# Resolve JARVIS_ROOT
# ---------------------------------------------------------------------------
if [ -z "${JARVIS_ROOT:-}" ]; then
    JARVIS_ROOT="$(git -C "$(dirname "${BASH_SOURCE[0]}")" rev-parse --show-toplevel 2>/dev/null \
        || git rev-parse --show-toplevel)"
fi
export JARVIS_ROOT

RUNS_DIR="${JARVIS_RUNS_DIR:-$JARVIS_ROOT/runs}"
POOLS_JSON="$JARVIS_ROOT/config/pools.json"

# ---------------------------------------------------------------------------
# Window: default 7, env JARVIS_WEEK_DAYS, --days N (flag wins)
# ---------------------------------------------------------------------------
DAYS="${JARVIS_WEEK_DAYS:-7}"
while [ "$#" -gt 0 ]; do
    case "$1" in
        --days) DAYS="${2:-7}"; shift 2 ;;
        --days=*) DAYS="${1#*=}"; shift ;;
        *) shift ;;
    esac
done

DONE_TAG="$(jq -r '.claim.done_tag' "$POOLS_JSON")"

# Unique project ids across pools (skip pools without a project field)
PROJECTS="$(jq -r '.pools[].project | select(. != null)' "$POOLS_JSON" | sort -u)"

# Aone link format mirrors plan.sh / board.sh
URL_FMT="https://project.aone.alibaba-inc.com/v2/project/%s/req/%s"

# ---------------------------------------------------------------------------
# Cutoff epoch: NOW - DAYS days
# ---------------------------------------------------------------------------
NOW_EPOCH=$(date -u +%s)
CUTOFF_EPOCH=$(( NOW_EPOCH - DAYS * 86400 ))

# Parse "YYYY-MM-DD HH:MM" (a1 gmtModified, local) → epoch. Empty on failure.
parse_epoch() {
    local ts="$1" e=""
    if e=$(date -j -f "%Y-%m-%d %H:%M" "$ts" +%s 2>/dev/null); then echo "$e"; return; fi
    if e=$(date -d "$ts" +%s 2>/dev/null); then echo "$e"; return; fi
    if e=$(python3 -c "
import sys, time, datetime
try:
    print(int(time.mktime(datetime.datetime.strptime('$ts', '%Y-%m-%d %H:%M').timetuple())))
except Exception:
    sys.exit(1)
" 2>/dev/null); then echo "$e"; return; fi
    echo ""
}

echo "# 本周回顾 — 近 ${DAYS} 天 jarvis-done"
echo ""

DONE_IDS=()

while IFS= read -r project; do
    [ -z "$project" ] && continue

    json=$(a1 project workitem list --project "$project" --tag "$DONE_TAG" -f json 2>/dev/null || echo "[]")

    # id<TAB>gmtModified<TAB>subject per line, jarvis-done only
    rows=$(printf '%s' "$json" | jq -r --arg t "$DONE_TAG" '
        .[]? | select((.tag // "") | test($t)) |
        "\(.identifier)\t\(.gmtModified // "")\t\(.subject // "")"' 2>/dev/null || true)

    while IFS=$'\t' read -r id mtime subject; do
        [ -z "$id" ] && continue
        [ "$id" = "null" ] && continue
        epoch=$(parse_epoch "$mtime")
        # keep when unparsable (don't silently drop) or within window
        if [ -n "$epoch" ] && [ "$epoch" -lt "$CUTOFF_EPOCH" ]; then
            continue
        fi
        url=$(printf "$URL_FMT" "$project" "$id")
        echo "- [$id] $subject — $url"
        DONE_IDS+=("$id")
    done <<< "$rows"

done <<< "$PROJECTS"

[ "${#DONE_IDS[@]}" -eq 0 ] && echo "_(无)_"

# ---------------------------------------------------------------------------
# Cross-check local runs/<date>-<id>.md ledger vs jarvis-done
# ---------------------------------------------------------------------------
LEDGER_IDS=()
if [ -d "$RUNS_DIR" ]; then
    while IFS= read -r f; do
        bn="$(basename "$f")"
        lid="${bn##*-}"; lid="${lid%.md}"
        [ -n "$lid" ] && LEDGER_IDS+=("$lid")
    done < <(ls "$RUNS_DIR"/*-*.md 2>/dev/null)
fi

in_list() { local x="$1"; shift; for e in "$@"; do [ "$e" = "$x" ] && return 0; done; return 1; }

echo ""
echo "## 漂移核对"
echo ""
echo "### 有台账但没 jarvis-done"
miss_done=0
if [ "${#LEDGER_IDS[@]}" -gt 0 ]; then
    for lid in "${LEDGER_IDS[@]}"; do
        if [ "${#DONE_IDS[@]}" -eq 0 ] || ! in_list "$lid" "${DONE_IDS[@]}"; then
            echo "- $lid (runs/ 有,在线无 jarvis-done)"
            miss_done=1
        fi
    done
fi
[ "$miss_done" -eq 0 ] && echo "_(无)_"

echo ""
echo "### 有 done 但本地无台账"
miss_led=0
if [ "${#DONE_IDS[@]}" -gt 0 ]; then
    for did in "${DONE_IDS[@]}"; do
        if [ "${#LEDGER_IDS[@]}" -eq 0 ] || ! in_list "$did" "${LEDGER_IDS[@]}"; then
            echo "- $did (在线 jarvis-done,runs/ 无台账)"
            miss_led=1
        fi
    done
fi
[ "$miss_led" -eq 0 ] && echo "_(无)_"

exit 0
