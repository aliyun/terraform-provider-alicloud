#!/usr/bin/env bash
# bootstrap/refresh.sh — "立即同步" entry: force-rescan Aone then rebuild the board.
# scan.sh --force (bypass 30min TTL) → board-html.sh → print done + scan.json category split.
set -uo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"

echo "refresh: rescanning Aone (--force)…" >&2
bash "$script_dir/scan.sh" --force >/dev/null || { echo "refresh: scan failed" >&2; exit 1; }

echo "refresh: rebuilding board…" >&2
bash "$script_dir/board-html.sh" >/dev/null || { echo "refresh: board build failed" >&2; exit 1; }

scan_f="$root/.my-day/scan.json"
total=$(jq 'length' "$scan_f" 2>/dev/null || echo 0)
dist=$(jq -r 'group_by(.category)|map("\(.[0].category // "—"):\(length)")|join(" ")' "$scan_f" 2>/dev/null)
echo "refresh: done — $total items ($dist)"
