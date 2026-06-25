#!/usr/bin/env bash
# bootstrap/wrap-check.sh — claim-ledger integrity check for Stop hooks.
#
# Scans ALL .my-day/claims-*.json files (not just today's) so that a claim
# opened before midnight is not silently skipped the next day.
# For each entry with done:false across all ledger files, checks if a runs/
# file exists via log.sh seen.
# If any unclosed claim has no run record → prints the offending ids and exits 2.
# Otherwise exits 0.
#
# Exit codes:
#   0 — all open claims have corresponding run records (or no claims files)
#   2 — one or more open claims are missing run records
#
# Respects:
#   JARVIS_ROOT       — repo root (default: directory above this script)
#   JARVIS_RUNS_DIR   — runs directory (default: ${JARVIS_ROOT}/runs)

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
export JARVIS_ROOT="$jarvis_root"
export JARVIS_RUNS_DIR="${JARVIS_RUNS_DIR:-$jarvis_root/runs}"

# Source log.sh so we can call seen() directly (mirrors sweep.sh pattern)
# shellcheck source=bootstrap/log.sh
source "$script_dir/log.sh"

myday_dir="$jarvis_root/.my-day"

# Collect all ledger files across all dates
ledger_files=()
if [ -d "$myday_dir" ]; then
    while IFS= read -r -d '' f; do
        ledger_files+=("$f")
    done < <(find "$myday_dir" -maxdepth 1 -name 'claims-*.json' -print0 2>/dev/null)
fi

# Touched ledgers: wrap.sh records every sync/done id here. Catches the blind spot
# where a ticket was handled but NEVER claimed → no claims-*.json entry → vacuous pass.
touched_files=()
while IFS= read -r -d '' f; do touched_files+=("$f"); done \
    < <(find "$myday_dir" -maxdepth 1 -name 'touched-*.json' -print0 2>/dev/null)

# No ledgers at all → nothing to check
if [ "${#ledger_files[@]}" -eq 0 ] && [ "${#touched_files[@]}" -eq 0 ]; then
    exit 0
fi

# Merge: open claims (done==false) + every touched id, deduplicate. Each must have a run_done.
open_ids=()
while IFS= read -r id; do
    [ -n "$id" ] && open_ids+=("$id")
done < <(
    {
        for ledger_file in ${ledger_files[@]+"${ledger_files[@]}"}; do
            jq -r '.[] | select(.done == false) | .id' "$ledger_file" 2>/dev/null
        done
        for tf in ${touched_files[@]+"${touched_files[@]}"}; do
            jq -r '.[]' "$tf" 2>/dev/null
        done
    } | sort -u
)

if [ "${#open_ids[@]}" -eq 0 ]; then
    exit 0
fi

# For each open id, use log.sh seen as the authoritative run-exists check
missing=()
for id in "${open_ids[@]}"; do
    if ! seen "$id"; then
        missing+=("$id")
    fi
done

if [ "${#missing[@]}" -eq 0 ]; then
    exit 0
fi

echo "wrap-check: claimed-or-touched workitems with no run_done (未收尾):" >&2
for id in "${missing[@]}"; do
    echo "  $id" >&2
done
exit 2
