#!/bin/bash
# pools.sh – per-pool inspection. Reads config/pools.json, queries each pool,
# prints active/total counts grouped by line. Active = total minus terminal states.
# Failing pool prints ERR (no abort). Read-only. Exits 0.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
pools_cfg="$jarvis_root/config/pools.json"

if [ ! -f "$pools_cfg" ]; then
    echo "pools.sh: config/pools.json not found at $pools_cfg" >&2
    exit 1
fi

# Terminal (closed-ish) statuses excluded from "active".
CLOSED='["Closed","Fixed","已完成","已发布","已取消","已拒绝","需求撤回","验收通过","已发布待需求方验收"]'

t_active=0; t_total=0
SIZE=200
while IFS=$'\t' read -r line name project; do
    json="[]"; page=1
    while :; do
        pg=$(a1 project workitem list --project "$project" --page "$page" --page-size "$SIZE" -f json 2>/dev/null)
        n=$(jq 'length' <<<"$pg" 2>/dev/null); [ -z "$n" ] && { json=""; break; }
        json=$(jq -s 'add' <<<"$json"$'\n'"$pg" 2>/dev/null)
        [ "$n" -lt "$SIZE" ] && break
        page=$((page+1))
    done
    total=$(jq 'length' <<<"$json" 2>/dev/null)
    if [ -z "$total" ]; then
        echo "$line/$name $project ERR"; continue
    fi
    active=$(jq --argjson c "$CLOSED" '[.[]|select(.status as $s|$c|index($s)|not)]|length' <<<"$json" 2>/dev/null)
    echo "$line/$name $project active=$active total=$total"
    t_active=$((t_active + active)); t_total=$((t_total + total))
done < <(jq -r '.pools | to_entries[] | [.value.line, .key, (.value.project|tostring)] | @tsv' "$pools_cfg")

echo "TOTAL active=$t_active total=$t_total"
exit 0
