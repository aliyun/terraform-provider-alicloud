#!/usr/bin/env bash
# bootstrap/cache.sh — TTL 文件缓存（rebuildable，落 .my-day/cache，已 gitignore）。
# Aone 仍是唯一真源；缓存只为省 a1 往返/装机自检，过期即重取。
# Usage:
#   cache.sh get  <key> <ttl_sec> -- <cmd...>  — 命中(未过期非空)直出，否则跑 cmd、存其 stdout、直出
#   cache.sh bust <key>                        — 删该 key（写操作后调用，丢陈旧详情）
#   cache.sh fresh <key> <ttl_sec>             — 仅判时效：新鲜退 0、否则退 1（不取内容）
#   cache.sh path <key>                        — 打印缓存文件路径
#   cache.sh age  <file_path>                  — 打印任意文件 mtime 距今秒数(跨平台,消除 stat -f %m/-c %Y 重实现);文件不存在退 1
# 命中走缓存的 get 会在 stderr 注 "cache hit <key>"；cmd 失败/空输出不落盘，下次重取。
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"
cache_dir="${JARVIS_CACHE_DIR:-$jarvis_root/.my-day/cache}"

_safe_key() { printf '%s' "$1" | tr -c 'A-Za-z0-9._-' '_'; }
_cache_file() { echo "$cache_dir/$(_safe_key "$1")"; }

# mtime 距今秒数（GNU stat -c %Y 优先 / macOS stat -f %m 兜底）。
# 顺序有讲究: GNU 的 `stat -f` 是文件系统模式, 失败时会把 "File: ..." 整块
# 打到 stdout 再非零退出, `$(A || B)` 会把垃圾和 B 的输出一起捕获, 下游算术
# 直接炸 "File: unbound variable"; 而 BSD 的 `stat -c` 失败时 stdout 干净。
_age_sec() {
    local f="$1" m
    m=$(stat -c %Y "$f" 2>/dev/null || stat -f %m "$f" 2>/dev/null) || return 1
    case "$m" in ''|*[!0-9]*) return 1;; esac
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
    age)
        f="${2:-}"
        [ -f "$f" ] || exit 1
        _age_sec "$f" || exit 1
        ;;
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
        # 原子落盘:temp + mv,并发同 key 不撕裂(失败/空不写,下次重取)
        if [ "$rc" -eq 0 ] && [ -n "$out" ]; then
            tmp="$f.$$.tmp"; printf '%s' "$out" > "$tmp" && mv -f "$tmp" "$f" || rm -f "$tmp"
        fi
        printf '%s' "$out"
        exit "$rc"
        ;;
    *) echo "Usage: cache.sh {get|bust|fresh|path|age} <key|path> [ttl] [-- cmd...]" >&2; exit 2 ;;
esac
