#!/bin/bash
# bootstrap/plan.sh — emit an execution plan and enforce the supervised auth gate.
#
# Reads mode from autonomy.md (default: supervised).
# Gets work items from scan.sh, filters out already-seen ids (via log.sh seen).
# Writes a plan to runs/plan-<UTCdate>.md.
#
# supervised → print "待授权:逐条 可行?允许?授权?" and EXIT 2  (awaiting human auth)
# unattended → EXIT 0
#
# NEVER calls any a1 write command — this script is read-only.
#
# Environment overrides (for testing):
#   JARVIS_RUNS_DIR        — override default <repo_root>/runs
#   JARVIS_AUTONOMY_FILE   — override autonomy.md path

set -uo pipefail

# ---------------------------------------------------------------------------
# Path resolution
# ---------------------------------------------------------------------------
_plan_repo_root() {
    git -C "$(dirname "${BASH_SOURCE[0]}")" rev-parse --show-toplevel 2>/dev/null \
        || git rev-parse --show-toplevel 2>/dev/null
}

repo_root="$(_plan_repo_root)"
bootstrap_dir="$repo_root/bootstrap"
runs_dir="${JARVIS_RUNS_DIR:-$repo_root/runs}"
autonomy_file="${JARVIS_AUTONOMY_FILE:-$repo_root/autonomy.md}"

mkdir -p "$runs_dir"

# ---------------------------------------------------------------------------
# Read mode from autonomy.md json block  (default: supervised)
# ---------------------------------------------------------------------------
mode="supervised"
if [ -f "$autonomy_file" ]; then
    extracted=$(grep -o '"mode":"[a-z]*"' "$autonomy_file" 2>/dev/null | head -1 | grep -o '[a-z]*"$' | tr -d '"')
    if [ -n "$extracted" ]; then
        mode="$extracted"
    fi
fi

# ---------------------------------------------------------------------------
# Get items from scan.sh
# ---------------------------------------------------------------------------
scan_script="$bootstrap_dir/scan.sh"
# Allow override via PATH (tests may inject a stub)
if command -v scan.sh > /dev/null 2>&1; then
    raw_items="$(scan.sh 2>/dev/null)"
else
    raw_items="$(bash "$scan_script" 2>/dev/null)"
fi

# ---------------------------------------------------------------------------
# Filter out already-seen ids using log.sh seen
# ---------------------------------------------------------------------------
log_script="$bootstrap_dir/log.sh"

# Build filtered item list as a JSON array of objects
filtered_items="[]"

item_count=$(echo "$raw_items" | jq 'length' 2>/dev/null || echo 0)

i=0
while [ "$i" -lt "$item_count" ]; do
    id=$(echo "$raw_items" | jq -r ".[$i].id")
    # Check if seen — use direct dispatch of log.sh
    if JARVIS_RUNS_DIR="$runs_dir" bash "$log_script" seen "$id" 2>/dev/null; then
        # Already seen — skip
        i=$((i + 1))
        continue
    fi
    # Append to filtered list
    item=$(echo "$raw_items" | jq ".[$i]")
    filtered_items=$(echo "$filtered_items" | jq ". + [$item]")
    i=$((i + 1))
done

# ---------------------------------------------------------------------------
# Build plan markdown
# ---------------------------------------------------------------------------
utc_date=$(date -u +%F)
plan_file="$runs_dir/plan-${utc_date}.md"

{
    echo "# Jarvis 执行计划 — ${utc_date}"
    echo ""
    echo "**模式:** $mode"
    echo ""
    echo "| ID | 标题 | 拟动作 | 置信 | auto/stop | 不可逆点 |"
    echo "|-----|------|--------|------|-----------|---------|"

    filtered_count=$(echo "$filtered_items" | jq 'length' 2>/dev/null || echo 0)

    j=0
    while [ "$j" -lt "$filtered_count" ]; do
        item_id=$(echo "$filtered_items" | jq -r ".[$j].id")
        item_title=$(echo "$filtered_items" | jq -r ".[$j].title")
        item_type=$(echo "$filtered_items" | jq -r ".[$j].type")
        item_status=$(echo "$filtered_items" | jq -r ".[$j].status")

        # Determine action based on type
        case "$item_type" in
            bug)      action="reply + create_req" ;;
            task)     action="reply" ;;
            req)      action="create_req + create_cr" ;;
            *)        action="reply" ;;
        esac

        # Confidence: default low_conf unless explicitly mapped
        # (stub field — real implementation would query OpenAPI/source)
        confidence="low_conf"
        auto_stop="stop"

        # Irreversibility point
        irreversible="create_cr / release_prod"

        echo "| $item_id | $item_title | $action | $confidence | $auto_stop | $irreversible |"

        j=$((j + 1))
    done

    echo ""
    echo "---"
    echo ""
    echo "待授权:逐条 可行?允许?授权?"
    echo ""
    echo "_计划生成时间: $(date -u +%Y-%m-%dT%H:%M:%SZ)_"
} > "$plan_file"

# ---------------------------------------------------------------------------
# Print plan to stdout
# ---------------------------------------------------------------------------
cat "$plan_file"

# ---------------------------------------------------------------------------
# Auth gate
# ---------------------------------------------------------------------------
if [ "$mode" = "supervised" ]; then
    exit 2
fi

# unattended
exit 0
