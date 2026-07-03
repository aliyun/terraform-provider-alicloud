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
# Owner-scoped block (cap-claim-ledger-owner-scoping Phase 1):
#   self = ${COORD_ID:-}. Each open/touched id's owner is looked up from the
#   claims ledger (id→owner map, first non-empty owner wins).
#     · owner non-empty AND != self → NOT counted as block; a WARN line is
#       printed for visibility and reconcile is left to clean up. A session must
#       not be blocked by another named instance's in-flight claim.
#     · owner == self, OR owner empty/missing (legacy entry, interactive session)
#       → original behavior: require a run_done, else count as missing (no
#       regression — you are still held to account for your own & ownerless claims).
#   Phase 1 deliberately does NOT consult coord.sh dead: a foreign instance's
#   claim is skipped regardless of liveness (dead-instance leftovers are handed to
#   reconcile). The coord-dead linkage is Phase 2.
#
# Exit codes:
#   0 — all in-scope open claims have corresponding run records (or no claims files)
#   2 — one or more in-scope open claims are missing run records
#
# Respects:
#   JARVIS_ROOT       — repo root (default: directory above this script)
#   JARVIS_RUNS_DIR   — runs directory (default: ${JARVIS_ROOT}/runs)
#   COORD_ID          — this instance's coord id (default: empty = interactive)

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/lib.sh
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"
export JARVIS_ROOT="$jarvis_root"
export JARVIS_RUNS_DIR="$(lib_runs_dir)"

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

# Build an id→owner lookup from all ledger files: one "id<TAB>owner" line per entry
# that carries a non-empty owner. Legacy entries without an owner field, or with an
# empty owner, are omitted → they resolve to "" on lookup (bash 3.2: no assoc array,
# so a flat temp file + awk first-match is used for a stable, portable map).
self="${COORD_ID:-}"
owner_map="$(mktemp)"
trap 'rm -f "$owner_map"' EXIT
for ledger_file in ${ledger_files[@]+"${ledger_files[@]}"}; do
    jq -r '.[] | select((.owner // "") != "") | "\(.id)\t\(.owner)"' "$ledger_file" 2>/dev/null
done > "$owner_map"

# owner_of <id> — echoes the first non-empty owner recorded for id, or "" if none.
owner_of() {
    awk -F'\t' -v id="$1" '$1 == id { print $2; exit }' "$owner_map"
}

# For each open id, scope by owner then use log.sh seen as the run-exists check.
missing=()
for id in "${open_ids[@]}"; do
    owner="$(owner_of "$id")"
    if [ -n "$owner" ] && [ "$owner" != "$self" ]; then
        # Owned by another named instance → not this session's problem. Warn for
        # visibility; reconcile is responsible for dead-instance leftovers.
        echo "wrap-check: skip $id (owned by other instance $owner)" >&2
        continue
    fi
    # owner == self, or owner empty/missing (legacy/interactive) → require run_done.
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
