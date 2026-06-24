#!/usr/bin/env bash
# bootstrap/claim.sh – claim/release primitive for concurrent multi-machine agents.
# Usage:
#   claim.sh claim   <workitem_id> <project_id>
#   claim.sh release <workitem_id> <project_id>
#
# claim:   tags the workitem jarvis-claimed, posts a timestamped comment, then reads
#          back the tag list to verify we won any concurrent race. Exits 0 on success,
#          1 if the workitem is not found in the readback (lost race).
# release: tags the workitem jarvis-done. Exits 0.
#
# Reads tag names from config/pools.json (.claim.tag and .claim.done_tag).
# Respects JARVIS_ROOT env override for the repo root.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
pools_cfg="$jarvis_root/config/pools.json"

if [ ! -f "$pools_cfg" ]; then
    echo "claim.sh: config/pools.json not found at $pools_cfg" >&2
    exit 1
fi

# Read tag names from config
CLAIM_TAG=$(jq -r '.claim.tag' "$pools_cfg")
DONE_TAG=$(jq -r '.claim.done_tag' "$pools_cfg")

cmd="${1:-}"
workitem_id="${2:-}"
project_id="${3:-}"

if [ -z "$cmd" ] || [ -z "$workitem_id" ] || [ -z "$project_id" ]; then
    echo "Usage: claim.sh <claim|release> <workitem_id> <project_id>" >&2
    exit 1
fi

myday_dir="$jarvis_root/.my-day"
today_date="$(date -u +%F)"
ledger_file="$myday_dir/claims-${today_date}.json"

# Atomically append/update entry in the claims ledger.
# Usage: _ledger_upsert <id> <done_value>   (done_value: true|false)
_ledger_upsert() {
    local id="$1"
    local done_val="$2"
    mkdir -p "$myday_dir"
    # Create file with empty array if absent
    if [ ! -f "$ledger_file" ]; then
        printf '[]' > "$ledger_file"
    fi
    local tmp
    tmp="$(mktemp "$myday_dir/.claims-tmp.XXXXXX")"
    # If id already exists: update done field; otherwise append new entry
    jq --arg id "$id" --argjson done "$done_val" \
        'if any(.[]; .id == $id) then
             map(if .id == $id then .done = $done else . end)
         else
             . + [{"id": $id, "done": $done}]
         end' \
        "$ledger_file" > "$tmp" && mv "$tmp" "$ledger_file"
}

case "$cmd" in
    claim)
        # 1. Tag the workitem as claimed
        a1 project workitem update "$workitem_id" --tag "$CLAIM_TAG"

        # 2. Post a timestamped comment identifying this machine
        utcnow=$(date -u +%Y-%m-%dT%H:%M:%SZ)
        a1 project workitem comment create "$workitem_id" -m "jarvis-claim $(hostname) $utcnow"

        # 3. Readback (with retries for tag-index lag): own id under claimed tag = won
        for _ in 1 2 3; do
            readback=$(a1 project workitem list --project "$project_id" --tag "$CLAIM_TAG" -f json 2>/dev/null || echo "[]")
            if echo "$readback" | jq -e --arg id "$workitem_id" 'any(.[]; ((.identifier // .id)|tostring)==$id)' > /dev/null 2>&1; then
                echo "claim.sh: claimed workitem $workitem_id in project $project_id"
                _ledger_upsert "$workitem_id" false
                exit 0
            fi
            sleep 2
        done
        echo "claim.sh: lost race for workitem $workitem_id (not found in readback)" >&2
        exit 1
        ;;

    release)
        # Tag as done (no untag exists; done_tag supersedes claimed)
        a1 project workitem update "$workitem_id" --tag "$DONE_TAG"
        echo "claim.sh: released workitem $workitem_id in project $project_id"
        _ledger_upsert "$workitem_id" true
        exit 0
        ;;

    *)
        echo "claim.sh: unknown command '$cmd'. Use claim or release." >&2
        exit 1
        ;;
esac
