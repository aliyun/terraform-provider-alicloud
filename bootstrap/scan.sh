#!/bin/bash
# scan.sh – pull assigned Aone work items per pool, emit [{id,title,type,status,pool}] JSON to stdout.
# Uses pool-scoped --project queries so claim tag filter (--filter NOT tag=<tag>) works.
# Empty or failing pools are skipped (non-fatal). Exits non-zero only on fatal errors.

set -uo pipefail

# Determine repo root: allow override via JARVIS_ROOT (used in tests), else derive from script location.
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"

# Verify we can authenticate.
account=$(a1 auth whoami | awk '/Account:/{print $2}')
if [ -z "$account" ]; then
  echo "scan.sh: could not determine account from 'a1 auth whoami'" >&2
  exit 1
fi

# Read claim tag and pools from config/pools.json if present.
claim_tag=""
pools_cfg="$jarvis_root/config/pools.json"
if [ -f "$pools_cfg" ]; then
  claim_tag=$(jq -r '.claim.tag // empty' "$pools_cfg" 2>/dev/null || true)
fi

# Check whether pools[].project entries exist.
has_pools=false
if [ -f "$pools_cfg" ]; then
  pool_count=$(jq '[.pools // {} | to_entries[] | .value.project] | length' "$pools_cfg" 2>/dev/null || echo 0)
  if [ "$pool_count" -gt 0 ] 2>/dev/null; then
    has_pools=true
  fi
fi

if $has_pools; then
  # Pool-scoped scan: pools fetched in parallel, paginated (page-size 1000), merged.
  PAGE_SIZE=1000

  fetch_pool() {  # args: key project status_csv title_csv → prints transformed JSON array
    local pool_key="$1" pool_project="$2" exclude_status="$3" exclude_title="$4"
    local filter="" pat pool_out="[]" page=1 pg n
    [ -n "$claim_tag" ] && filter="NOT tag=$claim_tag"
    [ -n "$exclude_status" ] && { [ -n "$filter" ] && filter="$filter AND "; filter="${filter}NOT status=$exclude_status"; }
    if [ -n "$exclude_title" ]; then
      IFS=',' read -ra _pats <<< "$exclude_title"
      for pat in "${_pats[@]}"; do [ -n "$filter" ] && filter="$filter AND "; filter="${filter}subject!~$pat"; done
    fi
    while :; do
      if [ -n "$filter" ]; then
        pg=$(a1 project workitem list --project "$pool_project" --assignee "$account" --columns id,title,status,priority,tag,type --filter "$filter" --page "$page" --page-size "$PAGE_SIZE" -f json 2>/dev/null) || true
      else
        pg=$(a1 project workitem list --project "$pool_project" --assignee "$account" --columns id,title,status,priority,tag,type --page "$page" --page-size "$PAGE_SIZE" -f json 2>/dev/null) || true
      fi
      n=$(echo "$pg" | jq 'length' 2>/dev/null); [ -z "$n" ] && break
      pool_out=$(jq -s 'add' <<<"$pool_out"$'\n'"$pg" 2>/dev/null) || pool_out="[]"
      [ "$n" -lt "$PAGE_SIZE" ] && break
      page=$((page+1))
    done
    echo "$pool_out" | jq --arg pool "$pool_key" '[.[] | {id:.identifier,title:.subject,type:.categoryIdentifier,status,pool:$pool,priority,tag}]'
  }

  tmpd=$(mktemp -d); trap 'rm -rf "$tmpd"' EXIT
  while IFS=$'\t' read -r pool_key pool_project pool_name exclude_status exclude_title; do
    fetch_pool "$pool_key" "$pool_project" "$exclude_status" "$exclude_title" > "$tmpd/$pool_key.json" 2>/dev/null &
  done < <(jq -r '
    .pools // {} | to_entries[] |
    [.key, (.value.project | tostring), (.value.name // .key), ((.value.exclude_status // [])|join(",")), ((.value.exclude_title // [])|join(","))] |
    @tsv
  ' "$pools_cfg")
  wait
  jq -s 'add // []' "$tmpd"/*.json 2>/dev/null || echo "[]"
else
  # No pools configured: fall back to assignee-based global list.
  if [ -n "$claim_tag" ]; then
    a1 project workitem list --assignee "$account" --columns id,title,status,priority,tag,type --filter "NOT tag=$claim_tag" -f json \
      | jq '[.[] | {id: .identifier, title: .subject, type: (.categoryIdentifier // .workitemType), status, priority, tag}]'
  else
    a1 project workitem list --assignee "$account" --columns id,title,status,priority,tag,type -f json \
      | jq '[.[] | {id: .identifier, title: .subject, type: (.categoryIdentifier // .workitemType), status, priority, tag}]'
  fi
fi
