#!/usr/bin/env bash
# bootstrap/cache.sh — TTL 文件缓存（rebuildable，落 .my-day/cache，已 gitignore）。
# Aone 仍是唯一真源；缓存只为省 a1 往返/装机自检，过期即重取。
# Usage:
#   cache.sh get  <key> <ttl_sec> -- <cmd...>  — 命中(未过期非空)直出，否则跑 cmd、存其 stdout、直出
#   cache.sh bust <key>                        — 删该 key（写操作后调用，丢陈旧详情）
#   cache.sh fresh <key> <ttl_sec>             — 仅判时效：新鲜退 0、否则退 1（不取内容）
#   cache.sh path <key>                        — 打印缓存文件路径
# 命中走缓存的 get 会在 stderr 注 "cache hit <key>"；cmd 失败/空输出不落盘，下次重取。
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"
cache_dir="${JARVIS_CACHE_DIR:-$jarvis_root/.my-day/cache}"

_safe_key() { printf '%s' "$1" | tr -c 'A-Za-z0-9._-' '_'; }
_cache_file() { echo "$cache_dir/$(_safe_key "$1")"; }

# mtime 距今秒数（macOS stat -f %m / GNU stat -c %Y）
_age_sec() {
    local f="$1" m
    m=$(stat -f %m "$f" 2>/dev/null || stat -c %Y "$f" 2>/dev/null) || return 1
    echo $(( $(date +%s) - m ))
}

cmd="${1:-}"
case "$cmd" in
    fresh)
        f="$(_cache_file "${2:-}")"; ttl="${3:-0}"
        [ -s "$f" ] || exit 1
        age=$(_age_sec "$f") || exit 1
        [ "$age" -lt "$ttl" ] && exit 0 || exit 1
        ;;
    path) _cache_file "${2:-}" ;;
    bust) rm -f "$(_cache_file "${2:-}")" 2>/dev/null; exit 0 ;;
    get)
        key="${2:-}"; ttl="${3:-0}"; shift 3
        [ "${1:-}" = "--" ] && shift
        f="$(_cache_file "$key")"
        if [ -s "$f" ] && age=$(_age_sec "$f") && [ "$age" -lt "$ttl" ]; then
            echo "cache hit $key (${age}s/${ttl}s)" >&2
            cat "$f"; exit 0
        fi
        mkdir -p "$cache_dir"
        out=$("$@"); rc=$?
        [ "$rc" -eq 0 ] && [ -n "$out" ] && printf '%s' "$out" > "$f"
        printf '%s' "$out"
        exit "$rc"
        ;;
    *) echo "Usage: cache.sh {get|bust|fresh|path} <key> [ttl] [-- cmd...]" >&2; exit 2 ;;
esac
