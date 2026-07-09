#!/usr/bin/env bash
# bootstrap/aone-get.sh — 取 Aone 工单详情，3h TTL 缓存（需求池 3 小时内基本无变）。
# 包一层 `a1 project workitem get <id> -f json`，省重复往返；写操作后由 wrap.sh bust。
# Usage: aone-get.sh <id>            （默认 -f json）
#        aone-get.sh <id> --raw      （透传 a1 文本，不入缓存）
# 强制重取:JARVIS_CACHE_TTL=0；缓存路径键 wi-<id>。
# a1 一律走 bin/a1id --(CLAUDE.md 身份纪律 #6):默认 jarvis 身份,不吃环境 ambient 登录。
# JARVIS_A1 覆盖(测试打桩用),与 wrap.sh / claim.sh 同款。
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="$(cd "$script_dir/.." && pwd)"
A1="${JARVIS_A1:-$jarvis_root/bin/a1id --}"
id="${1:?Usage: aone-get.sh <id> [--raw]}"
ttl="${JARVIS_CACHE_TTL:-10800}"   # 3h

if [ "${2:-}" = "--raw" ]; then
    $A1 project workitem get "$id"; exit $?
fi
bash "$script_dir/cache.sh" get "wi-$id" "$ttl" -- $A1 project workitem get "$id" -f json
