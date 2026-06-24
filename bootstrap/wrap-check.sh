#!/usr/bin/env bash
# bootstrap/wrap-check.sh — claim-ledger integrity check for Stop hooks.
#
# Reads today's claims-<date>.json from ${JARVIS_ROOT}/.my-day/.
# For each entry with done:false, checks if a runs/ file exists via log.sh seen.
# If any unclosed claim has no run record → prints the offending ids and exits 2.
# Otherwise exits 0.
#
# Exit codes:
#   0 — all open claims have corresponding run records (or no claims file)
#   2 — one or more open claims are missing run records
#
# Respects:
#   JARVIS_ROOT       — repo root (default: directory above this script)
#   JARVIS_RUNS_DIR   — runs directory (default: ${JARVIS_ROOT}/runs)

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
runs_dir="${JARVIS_RUNS_DIR:-$jarvis_root/runs}"

today_date="$(date -u +%F)"
ledger_file="$jarvis_root/.my-day/claims-${today_date}.json"

# No claims file for today → nothing to check
if [ ! -f "$ledger_file" ]; then
    exit 0
fi

# Extract all ids where done == false
open_ids=()
while IFS= read -r id; do
    [ -n "$id" ] && open_ids+=("$id")
done < <(jq -r '.[] | select(.done == false) | .id' "$ledger_file" 2>/dev/null)

if [ "${#open_ids[@]}" -eq 0 ]; then
    exit 0
fi

# For each open id, check if a run file exists
missing=()
for id in "${open_ids[@]}"; do
    # Match pattern: runs/<date>-<id>.md  (any date)
    match_count="$(ls "$runs_dir"/*-"${id}".md 2>/dev/null | wc -l | tr -d ' ')"
    if [ "$match_count" -eq 0 ]; then
        missing+=("$id")
    fi
done

if [ "${#missing[@]}" -eq 0 ]; then
    exit 0
fi

echo "wrap-check: unclosed claims with no run record:" >&2
for id in "${missing[@]}"; do
    echo "  $id" >&2
done
exit 2
