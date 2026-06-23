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
  # Pool-scoped scan: one call per pool, merged via jq -s add.
  # Results are emitted newline-separated so jq -s can consume them.
  pool_results=""

  # Iterate over pools: for each pool emit a JSON array (or skip on failure).
  while IFS=$'\t' read -r pool_key pool_project pool_name; do
    if [ -n "$claim_tag" ]; then
      pool_out=$(a1 project workitem list \
        --project "$pool_project" \
        --filter "NOT tag=$claim_tag" \
        --page-size 300 \
        -f json 2>/dev/null) || true
    else
      pool_out=$(a1 project workitem list \
        --project "$pool_project" \
        --page-size 300 \
        -f json 2>/dev/null) || true
    fi

    # Skip empty or failed pool output.
    if [ -z "$pool_out" ]; then
      continue
    fi

    # Transform items for this pool, adding the pool key as the pool field.
    transformed=$(echo "$pool_out" | jq --arg pool "$pool_key" \
      '[.[] | {id: .identifier, title: .subject, type: .categoryIdentifier, status, pool: $pool}]' \
      2>/dev/null) || true

    if [ -z "$transformed" ] || [ "$transformed" = "null" ]; then
      continue
    fi

    if [ -n "$pool_results" ]; then
      pool_results="${pool_results}
${transformed}"
    else
      pool_results="${transformed}"
    fi
  done < <(jq -r '
    .pools // {} | to_entries[] |
    [.key, (.value.project | tostring), (.value.name // .key)] |
    @tsv
  ' "$pools_cfg")

  if [ -z "$pool_results" ]; then
    echo "[]"
  else
    echo "$pool_results" | jq -s 'add // []'
  fi
else
  # No pools configured: fall back to assignee-based global list.
  if [ -n "$claim_tag" ]; then
    a1 project workitem list --assignee "$account" --filter "NOT tag=$claim_tag" -f json \
      | jq '[.[] | {id: .identifier, title: .subject, type: .categoryIdentifier, status}]'
  else
    a1 project workitem list --assignee "$account" -f json \
      | jq '[.[] | {id: .identifier, title: .subject, type: .categoryIdentifier, status}]'
  fi
fi
