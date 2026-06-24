#!/usr/bin/env bash
# bootstrap/reconcile.sh — collection drift reconciler
#
# Finds work items left jarvis-claimed but never released, and closes them.
#
# Logic:
#   For each project pool, list all claimed (jarvis-claimed) items that are
#   NOT yet marked done (NOT tag=jarvis-done). For each such item:
#     - If a local runs/<date>-<id>.md exists (log.sh seen) → work finished but
#       release was missed. Call claim.sh release to apply jarvis-done tag.
#       Print: RECONCILED: <id>
#     - If NO run file → leave it alone (sweep.sh handles stale-escalation).
#
# End: print RECONCILED: none if zero items were reconciled.
#
# Usage: bash bootstrap/reconcile.sh
#
# Environment overrides (for testing):
#   JARVIS_ROOT           — repo root (default: git rev-parse --show-toplevel)
#   JARVIS_RUNS_DIR       — runs dir (default: <JARVIS_ROOT>/runs)
#   JARVIS_ESCALATION_DIR — escalation dir (default: <JARVIS_ROOT>/escalation)
#
# Read-mostly; only writes are missed claim.sh release calls.

set -uo pipefail

# ---------------------------------------------------------------------------
# Resolve JARVIS_ROOT
# ---------------------------------------------------------------------------
if [ -z "${JARVIS_ROOT:-}" ]; then
    JARVIS_ROOT="$(git -C "$(dirname "${BASH_SOURCE[0]}")" rev-parse --show-toplevel 2>/dev/null \
        || git rev-parse --show-toplevel)"
fi
export JARVIS_ROOT

# ---------------------------------------------------------------------------
# Resolve script dir and co-located helpers
# ---------------------------------------------------------------------------
_reconcile_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source log.sh for the `seen` function
# shellcheck source=bootstrap/log.sh
source "$_reconcile_dir/log.sh"

# ---------------------------------------------------------------------------
# Read config from pools.json
# ---------------------------------------------------------------------------
POOLS_JSON="$JARVIS_ROOT/config/pools.json"

CLAIM_TAG="$(jq -r '.claim.tag' "$POOLS_JSON")"
DONE_TAG="$(jq -r '.claim.done_tag' "$POOLS_JSON")"

# Collect all project IDs from pools (mirrors sweep.sh pattern)
PROJECTS="$(jq -r '.pools[].project' "$POOLS_JSON")"

# ---------------------------------------------------------------------------
# Main reconcile loop
# ---------------------------------------------------------------------------
RECONCILED_IDS=()

while IFS= read -r project; do
    [ -z "$project" ] && continue

    # List all jarvis-claimed items in this project, excluding those marked done
    claimed_json=$(a1 project workitem list \
        --project "$project" \
        --tag "$CLAIM_TAG" \
        --filter "NOT tag=$DONE_TAG" \
        -f json 2>/dev/null || echo "[]")

    # Extract item identifiers (use 'identifier' key, fall back to 'id')
    item_ids=$(printf '%s' "$claimed_json" | python3 -c "
import json, sys
try:
    items = json.loads(sys.stdin.read())
    for item in items:
        item_id = item.get('identifier', item.get('id', ''))
        if item_id:
            print(str(item_id))
except Exception:
    pass
" 2>/dev/null || true)

    while IFS= read -r item_id; do
        [ -z "$item_id" ] && continue

        # Check if a run file exists: work finished but release was missed
        if seen "$item_id"; then
            # Release the claim: apply jarvis-done tag.
            # Prefer claim.sh found on PATH (allows test stubs); fall back to co-located script.
            _claim_cmd="$_reconcile_dir/claim.sh"
            if command -v claim.sh &>/dev/null; then
                _claim_cmd="claim.sh"
            fi
            $_claim_cmd release "$item_id" "$project"
            echo "RECONCILED: $item_id"
            RECONCILED_IDS+=("$item_id")
        fi
        # If no run file: leave it — sweep.sh handles stale escalation

    done <<< "$item_ids"

done <<< "$PROJECTS"

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
if [ "${#RECONCILED_IDS[@]}" -eq 0 ]; then
    echo "RECONCILED: none"
fi

exit 0
