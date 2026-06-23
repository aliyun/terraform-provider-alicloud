#!/usr/bin/env bash
# bootstrap/sweep.sh — stale-claim sweeper
#
# Finds jarvis-claimed workitems whose claim timestamp is older than claim.ttl_min
# and escalates them (writes escalation/<id>.md) via log.sh.
#
# Usage: bash bootstrap/sweep.sh
#
# Environment overrides (for testing):
#   JARVIS_ROOT           — repo root (default: git rev-parse --show-toplevel)
#   JARVIS_ESCALATION_DIR — escalation dir (default: <repo_root>/escalation)
#
# Read-only on Aone; only writes local escalation/ files.

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
# Source log.sh (provides escalate function)
# Prefer the real script directory so log.sh is always co-located with sweep.sh.
# ---------------------------------------------------------------------------
_sweep_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/log.sh
source "$_sweep_dir/log.sh"

# ---------------------------------------------------------------------------
# Read config
# ---------------------------------------------------------------------------
POOLS_JSON="$JARVIS_ROOT/config/pools.json"

TTL_MIN="$(python3 -c "import json,sys; d=json.load(open('$POOLS_JSON')); print(d['claim']['ttl_min'])")"

# Collect all project IDs from pools
PROJECTS="$(python3 -c "
import json, sys
d = json.load(open('$POOLS_JSON'))
projects = [str(p['project']) for p in d['pools'].values() if 'project' in p]
print('\n'.join(projects))
")"

# ---------------------------------------------------------------------------
# UTC epoch helper: parse ISO8601 UTC timestamp → epoch seconds
# Returns empty string on parse failure.
# ---------------------------------------------------------------------------
parse_utc_epoch() {
    local ts="$1"
    local epoch=""

    # macOS (BSD date)
    if epoch=$(date -u -j -f "%Y-%m-%dT%H:%M:%SZ" "$ts" +%s 2>/dev/null); then
        echo "$epoch"
        return
    fi

    # GNU date (Linux)
    if epoch=$(date -u -d "$ts" +%s 2>/dev/null); then
        echo "$epoch"
        return
    fi

    # Fallback: python3
    if epoch=$(python3 -c "
import sys, calendar, datetime
try:
    dt = datetime.datetime.strptime('$ts', '%Y-%m-%dT%H:%M:%SZ')
    print(int(calendar.timegm(dt.timetuple())))
except Exception:
    sys.exit(1)
" 2>/dev/null); then
        echo "$epoch"
        return
    fi

    # Parse failed
    echo ""
}

# ---------------------------------------------------------------------------
# Main sweep loop
# ---------------------------------------------------------------------------
NOW_EPOCH=$(date -u +%s)
TTL_SECS=$((TTL_MIN * 60))

STALE_IDS=()

while IFS= read -r project; do
    [ -z "$project" ] && continue

    # List all jarvis-claimed items in this project
    claimed_json=$(a1 project workitem list --project "$project" --tag jarvis-claimed -f json 2>/dev/null || echo "[]")

    # Extract item IDs (support both string and numeric ids)
    item_ids=$(printf '%s' "$claimed_json" | python3 -c "
import json, sys
try:
    items = json.loads(sys.stdin.read())
    for item in items:
        print(str(item.get('id', '')))
except Exception:
    pass
" 2>/dev/null || true)

    while IFS= read -r item_id; do
        [ -z "$item_id" ] && continue

        # Get comments for this item
        comments_json=$(a1 project workitem comment list "$item_id" -f json 2>/dev/null || echo "[]")

        # Find the latest jarvis-claim timestamp in comments
        latest_ts=$(printf '%s' "$comments_json" | python3 -c "
import json, sys, re
try:
    comments = json.loads(sys.stdin.read())
    pattern = re.compile(r'jarvis-claim\s+\S+\s+(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)')
    timestamps = []
    for c in comments:
        content = c.get('content', '')
        m = pattern.search(content)
        if m:
            timestamps.append(m.group(1))
    if timestamps:
        # Return lexicographically latest (ISO8601 sorts correctly)
        print(sorted(timestamps)[-1])
except Exception:
    pass
" 2>/dev/null || true)

        if [ -z "$latest_ts" ]; then
            # No jarvis-claim comment found; skip gracefully
            continue
        fi

        claim_epoch=$(parse_utc_epoch "$latest_ts")
        if [ -z "$claim_epoch" ]; then
            # Failed to parse timestamp; skip gracefully
            echo "WARN: could not parse claim timestamp '$latest_ts' for item $item_id, skipping" >&2
            continue
        fi

        age_secs=$(( NOW_EPOCH - claim_epoch ))

        if [ "$age_secs" -gt "$TTL_SECS" ]; then
            age_min=$(( age_secs / 60 ))
            escalate "$item_id" "stale claim >${TTL_MIN}min (age: ${age_min}min)"
            STALE_IDS+=("$item_id")
        fi

    done <<< "$item_ids"

done <<< "$PROJECTS"

# ---------------------------------------------------------------------------
# Print stale ids
# ---------------------------------------------------------------------------
if [ "${#STALE_IDS[@]}" -gt 0 ]; then
    echo "STALE_CLAIMS: ${STALE_IDS[*]}"
else
    echo "STALE_CLAIMS: none"
fi

exit 0
