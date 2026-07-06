#!/usr/bin/env bash
# bootstrap/mirror.sh — .claude ↔ .agents 双端镜像统一入口
#
# 合并了原三脚本(P1.b):
#   sync-to-codex.sh       → mirror.sh to-codex  [--all | <file>...]
#   sync-to-claude.sh      → mirror.sh to-claude [--all | <file>...]
#   skills-mirror-check.sh → mirror.sh check     [<file>...]
#
# sed 规则仍在 bootstrap/skills-mirror-lib.sh(source-only)。
#
# 注:AGENTS.md 已入库跟踪(内容 = mirror_sed_claude_to_codex(CLAUDE.md)),pre-commit
# 钩子在 CLAUDE.md staged 时自动重生成+git add,PostToolUse 实时 sync + preflight 兜底;
# check 对 CLAUDE.md↔AGENTS.md 双向配对校验 disk 一致性(round-trip 恒等,见 test)。
#
# 用法:
#   mirror.sh to-codex <file>...    # Claude → Codex 单文件 sync(hook 场景)
#   mirror.sh to-codex --all        # 全量 sync
#   mirror.sh to-claude <file>...   # 反向单文件
#   mirror.sh to-claude --all       # 反向全量
#   mirror.sh check                 # 全量 drift check(CLI verify)
#   mirror.sh check <file>...       # 指定文件 check(pre-commit 门禁场景)
#
# 修复漂移:
#   bash bootstrap/mirror.sh to-codex --all      # 以 .claude 为准 → .agents
#   bash bootstrap/mirror.sh to-claude --all     # 以 .agents 为准 → .claude
#
# to-codex/to-claude 无匹配路径静默 exit 0(不打断 hook 场景)。
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
jarvis_root="${JARVIS_ROOT:-$(cd "$script_dir/.." && pwd)}"

# shellcheck source=bootstrap/skills-mirror-lib.sh
source "$script_dir/skills-mirror-lib.sh"

# ============ sync core(to-codex + to-claude 共用) ============

# 计算 sync 目标路径。$1=源(绝对/相对); $2=方向(claude_to_codex|codex_to_claude)
_map_target() {
    local src="$1" direction="$2"
    [[ "$src" != /* ]] && src="$jarvis_root/$src"
    local rel="${src#$jarvis_root/}"
    case "$direction" in
        claude_to_codex)
            case "$rel" in
                .claude/skills/*) echo "$jarvis_root/.agents/skills/${rel#.claude/skills/}" ;;
                CLAUDE.md)        echo "$jarvis_root/AGENTS.md" ;;
                *) echo "" ;;
            esac ;;
        codex_to_claude)
            case "$rel" in
                .agents/skills/*) echo "$jarvis_root/.claude/skills/${rel#.agents/skills/}" ;;
                AGENTS.md)        echo "$jarvis_root/CLAUDE.md" ;;
                *) echo "" ;;
            esac ;;
    esac
}

_transform_file() {
    local src="$1" dst="$2" direction="$3"
    mkdir -p "$(dirname "$dst")"
    case "$direction" in
        claude_to_codex) mirror_sed_claude_to_codex < "$src" > "$dst" ;;
        codex_to_claude) mirror_sed_codex_to_claude < "$src" > "$dst" ;;
    esac
    echo "sync: $src → $dst"
}

_sync_one() {
    local src="$1" direction="$2"
    [[ "$src" != /* ]] && src="$jarvis_root/$src"
    [ -f "$src" ] || return 0
    local dst
    dst="$(_map_target "$src" "$direction")"
    [ -z "$dst" ] && return 0
    _transform_file "$src" "$dst" "$direction"
}

_sync_all() {
    local direction="$1"
    local src_root src_prefix top_md
    case "$direction" in
        claude_to_codex) src_root="$jarvis_root/.claude/skills"; top_md="$jarvis_root/CLAUDE.md" ;;
        codex_to_claude) src_root="$jarvis_root/.agents/skills"; top_md="$jarvis_root/AGENTS.md" ;;
    esac
    if [ -d "$src_root" ]; then
        while IFS= read -r f; do
            _sync_one "$f" "$direction"
        done < <(find "$src_root" -type f)
    fi
    [ -f "$top_md" ] && _sync_one "$top_md" "$direction"
}

# ============ check core ============

# 输入 mirror 源文件,输出目标路径 + transform 方向;无映射时返回 1
_map_pair() {
    local src="$1"
    [[ "$src" != /* ]] && src="$jarvis_root/$src"
    local rel="${src#$jarvis_root/}"
    case "$rel" in
        .claude/skills/*) echo "$jarvis_root/.agents/skills/${rel#.claude/skills/} claude_to_codex" ;;
        CLAUDE.md)        echo "$jarvis_root/AGENTS.md claude_to_codex" ;;
        .agents/skills/*) echo "$jarvis_root/.claude/skills/${rel#.agents/skills/} codex_to_claude" ;;
        AGENTS.md)        echo "$jarvis_root/CLAUDE.md codex_to_claude" ;;
        *) return 1 ;;
    esac
}

# 0 通过 / 1 drift(并打印) / 2 非 mirror 路径(静默)
_check_pair() {
    local src="$1"
    [[ "$src" != /* ]] && src="$jarvis_root/$src"
    [ -f "$src" ] || return 0

    local mapping
    mapping="$(_map_pair "$src")" || return 2
    local dst="${mapping% *}"
    local direction="${mapping##* }"

    if [ ! -f "$dst" ]; then
        # 通用:生成型 untracked 镜像目标(.gitignore 标注)缺失 ≠ drift——fresh worktree/
        # codex 首检出尚未按需生成时本就没有该文件;存在才走下方 sed 比对 catch 真 drift。
        # tracked 镜像(已入库的 AGENTS.md、.claude/.agents 双份 skills)缺失仍报 drift。
        # AGENTS.md 现已入库跟踪、不再命中本跳过分支;分支保留通用性防未来新增生成型镜像回退。
        if git -C "$jarvis_root" check-ignore -q "$dst" 2>/dev/null; then
            return 0
        fi
        echo "mirror-missing: $dst" >&2
        echo "  (expected mirror of $src)" >&2
        return 1
    fi

    local expected actual
    case "$direction" in
        claude_to_codex) expected="$(mirror_sed_claude_to_codex < "$src")" ;;
        codex_to_claude) expected="$(mirror_sed_codex_to_claude < "$src")" ;;
    esac
    actual="$(cat "$dst")"

    if [ "$expected" != "$actual" ]; then
        echo "mirror-drift: $dst" >&2
        echo "  (source: $src, direction: $direction)" >&2
        return 1
    fi
    return 0
}

_cmd_check() {
    local drift=0 rc
    if [ "$#" -eq 0 ]; then
        # 全量:两侧 skills + 顶层 md;累积全部 drift(不首错即退,
        # 用 process substitution 避免管道子 shell 吞变量)
        while IFS= read -r f; do
            _check_pair "$f"
            rc=$?
            [ "$rc" -eq 1 ] && drift=1
        done < <(
            [ -d "$jarvis_root/.claude/skills" ] && find "$jarvis_root/.claude/skills" -type f
            [ -d "$jarvis_root/.agents/skills" ] && find "$jarvis_root/.agents/skills" -type f
            [ -f "$jarvis_root/CLAUDE.md" ] && echo "$jarvis_root/CLAUDE.md"
            [ -f "$jarvis_root/AGENTS.md" ] && echo "$jarvis_root/AGENTS.md"
        )
    else
        # 指定文件(pre-commit hook 场景)
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
            echo "在 Edit/Write/MultiEdit 后跑 mirror.sh to-{codex,claude} 实时双向同步 +"
            echo "token 替换。本次 drift 通常是 Bash cp/sed/echo>、外部编辑器、其它 agent、"
            echo "或 hook 静默失败造成。"
            echo ""
            echo "修复:"
            echo "  bash bootstrap/mirror.sh to-codex --all      # 以 .claude 为准 → .agents"
            echo "  bash bootstrap/mirror.sh to-claude --all     # 以 .agents 为准 → .claude"
            echo "然后 git add 变更并重跑 commit。"
        } >&2
        return 1
    fi
    return 0
}

# ============ CLI dispatch ============

_help() { sed -n '2,22p' "$0"; }

cmd="${1:-}"; shift || true

case "$cmd" in
    to-codex)
        case "${1:-}" in
            "")    _help; exit 0 ;;
            --all) _sync_all claude_to_codex ;;
            *)     for f in "$@"; do _sync_one "$f" claude_to_codex; done ;;
        esac ;;
    to-claude)
        case "${1:-}" in
            "")    _help; exit 0 ;;
            --all) _sync_all codex_to_claude ;;
            *)     for f in "$@"; do _sync_one "$f" codex_to_claude; done ;;
        esac ;;
    check)
        _cmd_check "$@" ;;
    ""|-h|--help)
        _help ;;
    *)
        echo "mirror: unknown command '$cmd'" >&2
        _help
        exit 1 ;;
esac
