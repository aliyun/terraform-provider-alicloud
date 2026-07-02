#!/bin/bash
# scan.sh – pull assigned Aone work items per pool, emit [{id,title,type,status,pool,category}] JSON.
# Returns ALL assigned items incl. jarvis-claimed (board.sh needs them for 进行中/inflight).
# Uses pool-scoped --project queries; no claim-tag exclusion (dedup is downstream, not here).
# Each pool scanned thrice (--category req,bug,task); rows stamped category:"req|bug|task".
# Writes .my-day/scan.json AND echoes to stdout. 30min TTL: serve cached scan.json if fresh,
# unless --force (mirror preflight.sh; JARVIS_SCAN_TTL=0 forces too). Empty/failing pools
# skipped (non-fatal). Exits non-zero only on fatal errors.

set -uo pipefail

# Determine repo root: allow override via JARVIS_ROOT (used in tests), else derive via git-common-dir.
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"

# 30min TTL gate: serve cached scan.json if younger than TTL, unless --force (or JARVIS_SCAN_TTL=0).
# 走 cache.sh age 消除内嵌 stat -f %m/-c %Y 跨平台重实现(P1.d)。
out_f="$jarvis_root/.my-day/scan.json"
ttl="${JARVIS_SCAN_TTL:-1800}"   # 30min
[ "${1:-}" = "--force" ] && ttl=0
if [ "$ttl" -gt 0 ] && [ -s "$out_f" ]; then
  age=$(bash "$script_dir/cache.sh" age "$out_f" 2>/dev/null) || age=""
  if [ -n "$age" ] && [ "$age" -lt "$ttl" ]; then
    echo "scan.sh: skip (scan.json < $((ttl/60))min old; --force to rescan)" >&2
    cat "$out_f"; exit 0
  fi
fi

# Verify we can authenticate.
account=$(a1 auth whoami | awk '/Account:/{print $2}')
if [ -z "$account" ]; then
  echo "scan.sh: could not determine account from 'a1 auth whoami'" >&2
  exit 1
fi

# Pools come from config/pools.json. claim_tag is no longer used to filter scan output
# (claimed items are kept for the board); the key stays in config for claim.sh/triage dedup.
pools_cfg="$jarvis_root/config/pools.json"

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
    local filter="" pat pool_out="[]" cat page pg n
    # NOTE: jarvis-claimed items are intentionally KEPT (board.sh maps them → 进行中/inflight).
    # The old "NOT tag=$claim_tag" exclusion was triage-loop dedup, not for the board, and broke
    # 进行中 (always empty). If a triage caller needs dedup, filter on tag downstream, not in scan.
    [ -n "$exclude_status" ] && { [ -n "$filter" ] && filter="$filter AND "; filter="${filter}NOT status=$exclude_status"; }
    if [ -n "$exclude_title" ]; then
      IFS=',' read -ra _pats <<< "$exclude_title"
      for pat in "${_pats[@]}"; do [ -n "$filter" ] && filter="$filter AND "; filter="${filter}subject!~$pat"; done
    fi
    # Three categories: req,bug,task. --category makes categoryIdentifier authoritative; stamp literal.
    for cat in req bug task; do
      page=1
      while :; do
        if [ -n "$filter" ]; then
          pg=$(a1 project workitem list --project "$pool_project" --assignee "$account" --category "$cat" --columns id,title,status,priority,tag,type,category --filter "$filter" --page "$page" --page-size "$PAGE_SIZE" -f json 2>/dev/null) || true
        else
          pg=$(a1 project workitem list --project "$pool_project" --assignee "$account" --category "$cat" --columns id,title,status,priority,tag,type,category --page "$page" --page-size "$PAGE_SIZE" -f json 2>/dev/null) || true
        fi
        n=$(echo "$pg" | jq 'length' 2>/dev/null); [ -z "$n" ] && break
        pg=$(jq --arg c "$cat" '[.[] | .category=$c]' <<<"$pg" 2>/dev/null) || pg="[]"
        pool_out=$(jq -s 'add' <<<"$pool_out"$'\n'"$pg" 2>/dev/null) || pool_out="[]"
        [ "$n" -lt "$PAGE_SIZE" ] && break
        page=$((page+1))
      done
    done
    echo "$pool_out" | jq --arg pool "$pool_key" --arg proj "$pool_project" '[.[] | {id:.identifier,title:.subject,type:(.categoryIdentifier // .workitemType),status,pool:$pool,pool_project:$proj,priority,tag,category}]'
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
  result=$(jq -s 'add // []' "$tmpd"/*.json 2>/dev/null) || result="[]"
else
  # No pools configured: fall back to assignee-based global list (category unstamped).
  # No claim-tag exclusion — keep jarvis-claimed so the board can show 进行中.
  result=$(a1 project workitem list --assignee "$account" --columns id,title,status,priority,tag,type -f json \
    | jq '[.[] | {id: .identifier, title: .subject, type: (.categoryIdentifier // .workitemType), status, priority, tag, category: null}]') || result="[]"
fi

# Persist scan.json atomically (temp+mv, no torn file) and echo to stdout.
mkdir -p "$(dirname "$out_f")"
tmp="$out_f.$$.tmp"; printf '%s' "$result" > "$tmp" && mv -f "$tmp" "$out_f" || rm -f "$tmp"
printf '%s\n' "$result"
