#!/bin/bash
# scan.sh – pull assigned Aone work items, emit [{id,title,type,status}] JSON to stdout.
# Exits non-zero if a1 fails.

set -uo pipefail

# Determine repo root: allow override via JARVIS_ROOT (used in tests), else derive from script location.
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"

account=$(a1 auth whoami | awk '/Account:/{print $2}')
if [ -z "$account" ]; then
  echo "scan.sh: could not determine account from 'a1 auth whoami'" >&2
  exit 1
fi

# Read claim tag from config/pools.json if present; silently skip if missing or no tag.
claim_tag=""
pools_cfg="$jarvis_root/config/pools.json"
if [ -f "$pools_cfg" ]; then
  claim_tag=$(jq -r '.claim.tag // empty' "$pools_cfg" 2>/dev/null || true)
fi

if [ -n "$claim_tag" ]; then
  a1 project workitem list --assignee "$account" --filter "NOT tag=$claim_tag" -f json \
    | jq '[.[] | {id: .identifier, title: .subject, type: .categoryIdentifier, status}]'
else
  a1 project workitem list --assignee "$account" -f json \
    | jq '[.[] | {id: .identifier, title: .subject, type: .categoryIdentifier, status}]'
fi
