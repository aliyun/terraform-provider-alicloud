#!/usr/bin/env bash
set -uo pipefail
d="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"; R="${JARVIS_ROOT:-$d}"
I="$R/.my-day/instances"; T="$R/.my-day/tasks"; TTL="${COORD_TTL:-180}"
mkdir -p "$I" "$T"; cmd="${1:-}"
case "$cmd" in
 register) role="${2:-triage}"; id="$(hostname)-$$"; umask 077
   printf '{"id":"%s","role":"%s","pid":%s,"host":"%s","started":"%s","task":null}' \
     "$id" "$role" "$$" "$(hostname)" "$(date -u +%FT%TZ)" > "$I/$id.json"; : > "$I/$id.hb"; echo "$id";;
 heartbeat) : > "$I/${2}.hb";;
 dead) f="$I/${2}.hb"; pid="${2##*-}"; [ -f "$f" ] || exit 0
   kill -0 "$pid" 2>/dev/null && exit 1
   m=$(stat -f %m "$f" 2>/dev/null||stat -c %Y "$f"); [ $(( $(date +%s)-m )) -gt "$TTL" ] && exit 0 || exit 1;;
 *) echo "usage: coord.sh {register|heartbeat|dead}" >&2; exit 2;; esac
