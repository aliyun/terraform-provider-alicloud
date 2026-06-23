#!/usr/bin/env bash
# bootstrap/workspace.sh – project-workspace registry lookups. Reads config/workspaces.json.
# Subcommands:
#   resolve <repo|pool|line>   prints the workspace path (errors if no path)
#   ops <repo> <build|test|vet|fmt>  prints that op command
#   ensure <repo>              exit 0 if path exists else clone hint + exit 3
#   list                       prints all repos + paths (tab-separated)
#
# Config lookup honors JARVIS_WORKSPACES_FILE (default config/workspaces.json under
# JARVIS_ROOT). Read-only.

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
ws_cfg="${JARVIS_WORKSPACES_FILE:-$jarvis_root/config/workspaces.json}"

if [ ! -f "$ws_cfg" ]; then
    echo "workspace.sh: workspaces file not found at $ws_cfg" >&2
    exit 1
fi

# Match a workspace by repo name, a pool, or a line/key; emit its path (may be empty).
ws_path() {
    local key="$1"
    jq -r --arg k "$key" '
        .workspaces
        | to_entries
        | map(select(
            .key == $k
            or (.value.repo // "") == $k
            or (.value.pool // "") == $k
            or ((.value.pools // []) | index($k))
          ))
        | .[0].value.path // ""
    ' "$ws_cfg"
}

cmd="${1:-}"

case "$cmd" in
    resolve)
        key="${2:-}"
        if [ -z "$key" ]; then echo "Usage: workspace.sh resolve <repo|pool|line>" >&2; exit 1; fi
        path=$(ws_path "$key")
        if [ -z "$path" ]; then
            echo "workspace.sh: no path for '$key'" >&2; exit 1
        fi
        echo "$path"
        ;;

    ops)
        repo="${2:-}"; op="${3:-}"
        if [ -z "$repo" ] || [ -z "$op" ]; then
            echo "Usage: workspace.sh ops <repo> <build|test|vet|fmt>" >&2; exit 1
        fi
        c=$(jq -r --arg r "$repo" --arg o "$op" '
            .workspaces | to_entries
            | map(select(.key==$r or (.value.repo // "")==$r))
            | .[0].value.ops[$o] // ""
        ' "$ws_cfg")
        if [ -z "$c" ]; then echo "workspace.sh: no op '$op' for '$repo'" >&2; exit 1; fi
        echo "$c"
        ;;

    ensure)
        repo="${2:-}"
        if [ -z "$repo" ]; then echo "Usage: workspace.sh ensure <repo>" >&2; exit 1; fi
        path=$(ws_path "$repo")
        if [ -n "$path" ] && [ -d "$path" ]; then exit 0; fi
        echo "workspace.sh: '$repo' not present; clone to ${path:-<configure path>}" >&2
        exit 3
        ;;

    list)
        jq -r '.workspaces | to_entries[] | [(.value.repo // .key), (.value.path // "")] | @tsv' "$ws_cfg"
        ;;

    *)
        echo "Usage: workspace.sh <resolve|ops|ensure|list> ..." >&2
        exit 1
        ;;
esac
