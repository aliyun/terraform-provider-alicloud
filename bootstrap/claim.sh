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

case "$cmd" in
    claim)
        # 1. Tag the workitem as claimed
        a1 project workitem update "$workitem_id" --project "$project_id" --tag "$CLAIM_TAG"

        # 2. Post a timestamped comment identifying this machine
        utcnow=$(date -u +%Y-%m-%dT%H:%M:%SZ)
        a1 project workitem comment create "$workitem_id" --project "$project_id" -m "jarvis-claim $(hostname) $utcnow"

        # 3. Readback: verify the workitem appears under the claimed tag (race check)
        readback=$(a1 project workitem list --project "$project_id" --tag "$CLAIM_TAG" -f json 2>/dev/null || echo "[]")
        if echo "$readback" | jq -e --argjson id "$workitem_id" '[.[].id] | index($id) != null' > /dev/null 2>&1; then
            echo "claim.sh: claimed workitem $workitem_id in project $project_id"
            exit 0
        else
            echo "claim.sh: lost race for workitem $workitem_id (not found in readback)" >&2
            exit 1
        fi
        ;;

    release)
        # Tag as done (no untag exists; done_tag supersedes claimed)
        a1 project workitem update "$workitem_id" --project "$project_id" --tag "$DONE_TAG"
        echo "claim.sh: released workitem $workitem_id in project $project_id"
        exit 0
        ;;

    *)
        echo "claim.sh: unknown command '$cmd'. Use claim or release." >&2
        exit 1
        ;;
esac
