#!/usr/bin/env bash
# bootstrap/sync-to-claude.sh —— 把 Codex 侧改动镜像到 Claude 侧（双向对偶脚本，反向）
#
# 触发场景（由 .codex/hooks.json PostToolUse hook 调用，Edit/Write 后自动跑）：
#   .agents/skills/<name>/**  → .claude/skills/<name>/**
#   AGENTS.md                 → CLAUDE.md
#
# 关键词替换（Codex 侧文本 → Claude 侧文本；与 sync-to-codex.sh 互为逆）：
#   Codex                           → Claude Code
#   Co-Authored-By: Codex           → Co-Authored-By: Claude
#   .Codex/agents/                  → .claude/agents/
#   AGENTS.md                       → CLAUDE.md
#   codex-guide                     → claude-code-guide
#
# 使用：
#   sync-to-claude.sh <file_path>            单文件同步（hook 场景）
#   sync-to-claude.sh --all                   全量同步（初次或漂移修复）
#
# 输入无匹配路径则静默 exit 0（不打断 Edit/Write）。
# stdout 打印同步动作（一行一条），便于 hook 日志追踪；stderr 保留给错误。

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"

# 共享 sed transform 规则(sync-to-codex / sync-to-claude / skills-mirror-check 复用同一份)
# shellcheck source=bootstrap/skills-mirror-lib.sh
source "$script_dir/skills-mirror-lib.sh"

# Codex → Claude 关键词替换
_sed_transform() { mirror_sed_codex_to_claude; }

_transform_file() {
    local src="$1" dst="$2" transform="${3:-1}"
    mkdir -p "$(dirname "$dst")"
    if [ "$transform" = "1" ]; then
        _sed_transform < "$src" > "$dst"
    else
        cp "$src" "$dst"
    fi
    echo "sync: $src → $dst"
}

# 判断路径类型 + 计算目标路径
_map_target() {
    local abs="$1"
    local rel="${abs#$jarvis_root/}"
    case "$rel" in
        .agents/skills/*)
            echo "$jarvis_root/.claude/skills/${rel#.agents/skills/}"
            ;;
        AGENTS.md)
            echo "$jarvis_root/CLAUDE.md"
            ;;
        *)
            echo ""
            ;;
    esac
}

_sync_one() {
    local src="$1"
    [[ "$src" != /* ]] && src="$jarvis_root/$src"
    [ -f "$src" ] || return 0
    local dst
    dst="$(_map_target "$src")"
    [ -z "$dst" ] && return 0
    _transform_file "$src" "$dst" 1
}

_sync_all() {
    if [ -d "$jarvis_root/.agents/skills" ]; then
        while IFS= read -r f; do
            _sync_one "$f"
        done < <(find "$jarvis_root/.agents/skills" -type f)
    fi
    [ -f "$jarvis_root/AGENTS.md" ] && _sync_one "$jarvis_root/AGENTS.md"
}

case "${1:-}" in
    ""|-h|--help)
        sed -n '2,20p' "$0"
        exit 0
        ;;
    --all)
        _sync_all
        ;;
    *)
        _sync_one "$1"
        ;;
esac
