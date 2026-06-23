#!/bin/bash
# pools.sh – per-pool open-item inspection.
# Reads config/pools.json, queries each pool via a1, prints open counts.
# Skips failing pools (prints ERR) so one bad pool doesn't abort the rest.
# Prints a grand total and exits 0.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"

pools_cfg="$jarvis_root/config/pools.json"

if [ ! -f "$pools_cfg" ]; then
    echo "pools.sh: config/pools.json not found at $pools_cfg" >&2
    exit 1
fi

total=0

# Iterate over each pool name; extract project id from JSON.
while IFS=$'\t' read -r name project; do
    count=$(a1 project workitem list --project "$project" -f json 2>/dev/null | jq length 2>/dev/null)
    if [ -z "$count" ]; then
        echo "$name $project ERR"
    else
        echo "$name $project open=$count"
        total=$((total + count))
    fi
done < <(jq -r '.pools | to_entries[] | [.key, (.value.project | tostring)] | @tsv' "$pools_cfg")

echo "total=$total"
exit 0
