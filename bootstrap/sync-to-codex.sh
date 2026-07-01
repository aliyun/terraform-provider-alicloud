#!/usr/bin/env bash
# bootstrap/sync-to-codex.sh —— 把 Claude 侧改动镜像到 Codex 侧（两平台并行策略文档）
#
# 触发场景（由 .claude/settings.json PostToolUse hook 调用，Edit/Write 后自动跑）：
#   .claude/skills/<name>/**  → .agents/skills/<name>/**
#   CLAUDE.md                 → AGENTS.md
#
# 关键词替换（Claude 侧文本 → Codex 侧文本）：
#   Claude Code                     → Codex
#   Co-Authored-By: Claude          → Co-Authored-By: Codex
#   .claude/agents/                 → .Codex/agents/
#   CLAUDE.md                       → AGENTS.md
#   claude-code-guide               → codex-guide
#
# 使用：
#   sync-to-codex.sh <file_path>            单文件同步（hook 场景）
#   sync-to-codex.sh --all                   全量同步（初次或漂移修复）
#
# 输入无匹配路径则静默 exit 0（不打断 Edit/Write）。
# stdout 打印同步动作（一行一条），便于 hook 日志追踪；stderr 保留给错误。

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"

# Claude → Codex 关键词替换（顺序敏感：更长的先替换避免误伤）
_sed_transform() {
    sed \
        -e 's|Claude Code|Codex|g' \
        -e 's|Co-Authored-By: Claude|Co-Authored-By: Codex|g' \
        -e 's|\.claude/agents/|.Codex/agents/|g' \
        -e 's|claude-code-guide|codex-guide|g' \
        -e 's|CLAUDE\.md|AGENTS.md|g'
}

# 从源文件生成目标文件（做替换）；--no-transform 时纯 cp（用于 assets 等不该改的内容）
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
# 输入：绝对路径或仓库相对路径
# 输出：目标绝对路径（找不到映射则空）
_map_target() {
    local abs="$1"
    # 归一化为仓根相对路径
    local rel="${abs#$jarvis_root/}"
    case "$rel" in
        .claude/skills/*)
            echo "$jarvis_root/.agents/skills/${rel#.claude/skills/}"
            ;;
        CLAUDE.md)
            echo "$jarvis_root/AGENTS.md"
            ;;
        *)
            echo ""
            ;;
    esac
}

_sync_one() {
    local src="$1"
    # 相对路径也接受
    [[ "$src" != /* ]] && src="$jarvis_root/$src"
    [ -f "$src" ] || return 0   # 文件被删了不 sync（保留 codex 侧原样，人工判断）
    local dst
    dst="$(_map_target "$src")"
    [ -z "$dst" ] && return 0   # 非同步范围内的文件（大多数编辑），静默跳过
    _transform_file "$src" "$dst" 1
}

# 全量同步：扫 .claude/skills/**/* 和 CLAUDE.md
_sync_all() {
    # skills：目录树保持一致，SKILL.md 和内部 references / scripts 都同步
    if [ -d "$jarvis_root/.claude/skills" ]; then
        while IFS= read -r f; do
            _sync_one "$f"
        done < <(find "$jarvis_root/.claude/skills" -type f)
    fi
    # 顶层 CLAUDE.md
    [ -f "$jarvis_root/CLAUDE.md" ] && _sync_one "$jarvis_root/CLAUDE.md"
}

# ---------- entry ----------
case "${1:-}" in
    ""|-h|--help)
        sed -n '2,17p' "$0"
        exit 0
        ;;
    --all)
        _sync_all
        ;;
    *)
        _sync_one "$1"
        ;;
esac
