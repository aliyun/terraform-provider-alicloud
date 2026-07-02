#!/bin/bash
# bootstrap/plan.sh — emit an execution plan and enforce the supervised auth gate.
#
# Reads mode from autonomy.md (default: supervised).
# Gets work items from stdin JSON if provided, else falls back to scan.sh.
# Filters out already-seen ids via a SINGLE batch jq pass (no per-item fork).
# Writes a plan to runs/plan-<UTCdate>.md (includes priority column).
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
# Path resolution (P1.e:走 lib.sh 统一,消除 _plan_repo_root 样板)
# ---------------------------------------------------------------------------
# shellcheck source=bootstrap/lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"
repo_root="$(jarvis_root)"
bootstrap_dir="$repo_root/bootstrap"
runs_dir="$(lib_runs_dir)"
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
# INPUT: read items JSON from stdin if present, else fall back to scan.sh once
# ---------------------------------------------------------------------------
if [ -t 0 ]; then
    # stdin is a terminal — not piped, run scan.sh
    stdin_data=""
else
    # stdin is a pipe or file — try to read it
    stdin_data=$(cat)
fi

if [ -n "${stdin_data:-}" ] && echo "$stdin_data" | jq -e 'type == "array"' >/dev/null 2>&1; then
    # Valid JSON array from stdin — use it directly, no scan.sh
    raw_items="$stdin_data"
else
    # No stdin JSON — run scan.sh once
    scan_script="$bootstrap_dir/scan.sh"
    if command -v scan.sh > /dev/null 2>&1; then
        raw_items="$(scan.sh 2>/dev/null)"
    else
        raw_items="$(bash "$scan_script" 2>/dev/null)"
    fi
fi

# ---------------------------------------------------------------------------
# DEDUP BATCH: list runs/ ONCE, build seen-id set, filter via single jq pass
# ---------------------------------------------------------------------------
# Collect all seen ids from runs/*-<id>.md filenames in one ls + sed pass
seen_ids_json="[]"
if ls "$runs_dir"/*-*.md >/dev/null 2>&1; then
    # Extract ids from filenames: YYYY-MM-DD-<id>.md → <id>
    # Pattern: anything after the date prefix (4d-2d-2d-)
    seen_ids_json=$(ls "$runs_dir"/*-*.md 2>/dev/null \
        | sed 's|.*/[0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}-\(.*\)\.md$|\1|' \
        | jq -R . | jq -s '.')
fi

# Single jq pass: reject items whose id is in the seen set
filtered_items=$(echo "$raw_items" \
    | jq --argjson seen "$seen_ids_json" \
         '[.[] | select(.id as $id | $seen | index($id) | not)]' \
    2>/dev/null || echo "[]")

# ---------------------------------------------------------------------------
# Build plan markdown (includes priority column)
# ---------------------------------------------------------------------------
utc_date=$(date -u +%F)
plan_file="$runs_dir/plan-${utc_date}.md"

{
    echo "# Jarvis 执行计划 — ${utc_date}"
    echo ""
    echo "**模式:** $mode"
    echo ""
    echo "| ID | 标题 | 优先级 | 拟动作 | 置信 | auto/stop | 不可逆点 |"
    echo "|-----|------|--------|--------|------|-----------|---------|"

    # pool→project_id map from config/pools.json (for clickable Aone links; fallback {} on read fail)
    pool_proj_json=$(jq -c '.pools | map_values(.project)' "$repo_root/config/pools.json" 2>/dev/null || echo '{}')

    # Single jq pass builds all rows (type→action map; confidence stub; ID as Aone link)
    echo "$filtered_items" | jq -r --argjson pp "$pool_proj_json" '.[] |
      (if (.type|test("bug|Bug")) then "reply + create_req"
       elif (.type|test("req|Req|需求")) then "create_req + create_cr"
       else "reply" end) as $a |
      (($pp[.pool] // 0) | tostring) as $proj |
      "| [\(.id)](https://project.aone.alibaba-inc.com/v2/project/\($proj)/req/\(.id)) | \(.title) | \(.priority // "P2") | \($a) | low_conf | escalate | create_cr / release_prod |"'

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
# Auth gate  (default-deny: exit 0 ONLY for exact "unattended")
# ---------------------------------------------------------------------------
if [ "$mode" = "unattended" ]; then
    exit 0
fi

# supervised, empty, garbage, missing file → all await human auth
echo "待授权:逐条 可行?允许?授权?" >&2
exit 2
