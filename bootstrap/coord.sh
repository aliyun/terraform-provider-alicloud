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
 checkpoint) aid="$2"; st="$3"; wt="${4:-}"; br="${5:-}"; rp="${6:-}"; umask 077; tmp=$(mktemp "$T/.t.XXXX"); printf '{"aone_id":"%s","owner_instance":"%s","stage":"%s","worktree":"%s","branch":"%s","repo":"%s","updated":"%s"}' "$aid" "${COORD_ID:-}" "$st" "$wt" "$br" "$rp" "$(date -u +%FT%TZ)" >"$tmp" && mv "$tmp" "$T/$aid.json";;
 list-orphans) for f in "$T"/*.json; do [ -e "$f" ]||continue; o=$(jq -r .owner_instance "$f"); [ -n "$o" ]&&[ "$o" != null ] && bash "$0" dead "$o" && jq -r .aone_id "$f"; done;;
 adopt) f="$T/$2.json"; [ -f "$f" ]||exit 1; jq --arg i "${COORD_ID:-}" '.owner_instance=$i' "$f">"$f.t"&&mv "$f.t" "$f";;
 *) echo "usage: coord.sh {register|heartbeat|dead|checkpoint|list-orphans|adopt}" >&2; exit 2;; esac
