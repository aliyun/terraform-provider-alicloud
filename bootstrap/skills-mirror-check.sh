#!/usr/bin/env bash
# bootstrap/skills-mirror-check.sh — .claude/skills 与 .agents/skills 双份镜像兜底门禁
#
# 主线机制:PostToolUse hook(.claude/settings.json + .codex/hooks.json)在每次
# Edit/Write/MultiEdit 后跑 sync-to-{codex,claude}.sh,实时双向同步 + agent-specific
# token 替换:
#   Claude Code             ↔ Codex
#   CLAUDE.md               ↔ AGENTS.md
#   .claude/agents/         ↔ .Codex/agents/
#   claude-code-guide       ↔ codex-guide
#   Co-Authored-By: Claude  ↔ Co-Authored-By: Codex
#
# 本脚本**只兜底**主线没接住的场景:Bash cp/sed/echo>、外部编辑器、其它 agent、
# hook 静默失败、跨会话遗留。**不改任何文件**——只做 check,drift 时打印列表 + exit 1。
#
# 用法:
#   bash bootstrap/skills-mirror-check.sh                    # 全量 check
#   bash bootstrap/skills-mirror-check.sh <file> [file...]   # 只 check 指定文件对
#
# SED 规则来自共享 lib(bootstrap/skills-mirror-lib.sh),与 sync-to-codex.sh
# 和 sync-to-claude.sh 复用同一份 source of truth。
#
# 修复:
#   bash bootstrap/sync-to-codex.sh --all    # 以 .claude 为准全量推到 .agents
#   bash bootstrap/sync-to-claude.sh --all   # 以 .agents 为准全量推到 .claude
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="$(cd "$script_dir/.." && pwd)"

# 共享 sed transform 规则(sync-to-codex / sync-to-claude / skills-mirror-check 复用同一份)
# shellcheck source=bootstrap/skills-mirror-lib.sh
source "$script_dir/skills-mirror-lib.sh"

# 局部别名:保留旧函数名让下文 _check_pair 里的 case 无需改
_sed_claude_to_codex() { mirror_sed_claude_to_codex; }
_sed_codex_to_claude() { mirror_sed_codex_to_claude; }

# 输入 mirror 源文件(绝对或相对),输出:目标绝对路径 + transform 方向
# 无映射时输出空,返回码 1
_map() {
    local src="$1"
    [[ "$src" != /* ]] && src="$jarvis_root/$src"
    local rel="${src#$jarvis_root/}"
    case "$rel" in
        .claude/skills/*)  echo "$jarvis_root/.agents/skills/${rel#.claude/skills/} claude_to_codex" ;;
        CLAUDE.md)         echo "$jarvis_root/AGENTS.md claude_to_codex" ;;
        .agents/skills/*)  echo "$jarvis_root/.claude/skills/${rel#.agents/skills/} codex_to_claude" ;;
        AGENTS.md)         echo "$jarvis_root/CLAUDE.md codex_to_claude" ;;
        *) return 1 ;;
    esac
}

# 检查一对文件的 mirror 关系
# 返回:0 通过 / 1 drift(并打印) / 2 非 mirror 路径(静默)
_check_pair() {
    local src="$1"
    [[ "$src" != /* ]] && src="$jarvis_root/$src"
    [ -f "$src" ] || return 0    # 文件被删/不存在,不 flag(sync 脚本也会跳过)

    local mapping
    mapping="$(_map "$src")" || return 2
    local dst="${mapping% *}"
    local direction="${mapping##* }"

    if [ ! -f "$dst" ]; then
        echo "mirror-missing: $dst" >&2
        echo "  (expected mirror of $src)" >&2
        return 1
    fi

    local expected actual
    case "$direction" in
        claude_to_codex) expected="$(_sed_claude_to_codex < "$src")" ;;
        codex_to_claude) expected="$(_sed_codex_to_claude < "$src")" ;;
    esac
    actual="$(cat "$dst")"

    if [ "$expected" != "$actual" ]; then
        echo "mirror-drift: $dst" >&2
        echo "  (source: $src, direction: $direction)" >&2
        return 1
    fi
    return 0
}

drift=0

if [ "$#" -eq 0 ]; then
    # 全量 check:遍历两侧 skills + 顶层 md
    {
        [ -d "$jarvis_root/.claude/skills" ] && find "$jarvis_root/.claude/skills" -type f
        [ -d "$jarvis_root/.agents/skills" ] && find "$jarvis_root/.agents/skills" -type f
        [ -f "$jarvis_root/CLAUDE.md" ] && echo "$jarvis_root/CLAUDE.md"
        [ -f "$jarvis_root/AGENTS.md" ] && echo "$jarvis_root/AGENTS.md"
    } | while IFS= read -r f; do
        _check_pair "$f"
        rc=$?
        [ "$rc" -eq 1 ] && exit 1
    done || drift=1
else
    # 指定文件模式(pre-commit hook 用):只 check 传入的路径
    for f in "$@"; do
        _check_pair "$f"
        rc=$?
        [ "$rc" -eq 1 ] && drift=1
    done
fi

if [ "$drift" -eq 1 ]; then
    {
        echo ""
        echo "skills mirror drift detected between .claude and .agents."
        echo "主线机制是 PostToolUse hook(.claude/settings.json + .codex/hooks.json)"
        echo "在 Edit/Write/MultiEdit 后跑 sync-to-{codex,claude}.sh 实时双向同步 + token"
        echo "替换。本次 drift 通常是 Bash cp/sed/echo>、外部编辑器、其它 agent、或"
        echo "hook 静默失败造成。"
        echo ""
        echo "修复:"
        echo "  bash bootstrap/sync-to-codex.sh --all      # 以 .claude 为准 → .agents"
        echo "  bash bootstrap/sync-to-claude.sh --all     # 以 .agents 为准 → .claude"
        echo "然后 git add 变更并重跑 commit。"
    } >&2
    exit 1
fi
exit 0
