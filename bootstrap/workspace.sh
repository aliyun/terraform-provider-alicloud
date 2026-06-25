#!/bin/bash
# workspace.sh – resolve a workspace path for jarvis, per CLAUDE.md 工作纪律 #3.
# Truth = config/workspaces.json (checkin, no machine paths) deep-merged with
# config/workspaces.local.json (gitignored, machine-local; wins). Per repo:
#   (a) merged .path 存在且目录在 → 用之
#   (b) 否则 ${JARVIS_WORKSPACE_ROOT:-~/workspace}/<repo> 存在 → 用之
#   (c) 否则有 git_url → 打印 clone 目标(ROOT/repo)供调用方 clone
#   (d) 无 git_url → escalate(missing_capability)
# Read-only. Commands: dir <key> | config [<key>] | list. Exit 4 = missing_capability.
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
base="$jarvis_root/config/workspaces.json"
local="$jarvis_root/config/workspaces.local.json"
ws_root="${JARVIS_WORKSPACE_ROOT:-$HOME/workspace}"

[ -f "$base" ] || { echo "workspace.sh: $base not found" >&2; exit 1; }

# Deep-merge base * local (local wins); local missing → base alone.
merged() {
  if [ -f "$local" ]; then jq -s '.[0] * .[1]' "$base" "$local"
  else cat "$base"; fi
}

usage() { echo "usage: workspace.sh {dir <key>|config [key]|list}" >&2; exit 2; }

cmd="${1:-}"; key="${2:-}"
case "$cmd" in
  list)   merged | jq -r '.workspaces | keys[]' ;;
  config) merged | { [ -n "$key" ] && jq -e ".workspaces[\"$key\"]" || jq '.workspaces'; } ;;
  dir)
    [ -n "$key" ] || usage
    w=$(merged | jq -e ".workspaces[\"$key\"]" 2>/dev/null) \
      || { echo "workspace.sh: unknown workspace '$key' (escalate: missing_capability)" >&2; exit 4; }
    repo=$(jq -r '.repo // empty' <<<"$w")
    path=$(jq -r '.path // empty' <<<"$w")
    url=$(jq -r '.git_url // empty' <<<"$w")
    # (a) merged path
    if [ -n "$path" ] && [ -d "$path" ]; then echo "$path"; exit 0; fi
    # (b) ROOT/repo
    if [ -n "$repo" ] && [ -d "$ws_root/$repo" ]; then echo "$ws_root/$repo"; exit 0; fi
    # (c) clone target
    if [ -n "$url" ]; then
      echo "$ws_root/$repo"
      echo "workspace.sh: '$key' not on disk; clone $url -> $ws_root/$repo" >&2
      exit 0
    fi
    # (d) no path, no url
    echo "workspace.sh: '$key' has no path and no git_url (escalate: missing_capability)" >&2
    exit 4 ;;
  *) usage ;;
esac
