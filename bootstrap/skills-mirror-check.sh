#!/usr/bin/env bash
# bootstrap/skills-mirror-check.sh — .claude/skills 与 .agents/skills 双份镜像门禁
#
# 背景:本仓 skills 手工双份维护——.claude/skills(Claude Code 加载)与
# .agents/skills(codex 加载)历史上并行副本。任何单向改动会导致漂移;参见工单
# 83718139 教训:改 .claude 忘改 .agents 造成 codex 侧路由用旧版规则。
#
# 但两份**不是完全 mirror**——SKILL.md 入口文案、agent-specific 引用(AGENTS.md
# vs CLAUDE.md)、脚本路径(~/.Codex/ vs ~/.claude/)、AI 水印(Codex vs Claude Code)
# 都是必然的合法分歧,不能强制一致。
#
# 所以只对 **allowlist**(bootstrap/skills-mirror-allowlist)里列出的文件做严格
# 比对。allowlist 只应含"两侧完全同内容"的纯规则/脚本;其它文件由人工同步。
#
# 用法:
#   bash bootstrap/skills-mirror-check.sh                   # = check
#   bash bootstrap/skills-mirror-check.sh check             # drift → exit 1 + diff
#   bash bootstrap/skills-mirror-check.sh sync-claude-to-agents
#                                                           # 强制以 .claude 覆盖 .agents(仅 allowlist 文件)
#
# 由 pre-commit hook(bootstrap/git-hooks/pre-commit)调用 check 模式。
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"

CLAUDE_DIR="$repo_root/.claude/skills"
AGENTS_DIR="$repo_root/.agents/skills"
ALLOWLIST="$script_dir/skills-mirror-allowlist"

if [ ! -d "$CLAUDE_DIR" ] || [ ! -d "$AGENTS_DIR" ]; then
    exit 0    # 任一侧不存在 → 本机不启用双份,跳过
fi
if [ ! -f "$ALLOWLIST" ]; then
    echo "skills-mirror-check: allowlist missing: $ALLOWLIST" >&2
    exit 0    # 无 allowlist 时静默跳过,避免误伤
fi

cmd="${1:-check}"
drift=0

while IFS= read -r rel; do
    # 跳过空行/注释
    rel="${rel%%#*}"
    rel="$(echo "$rel" | tr -d '[:space:]')"
    [ -z "$rel" ] && continue

    claude_file="$CLAUDE_DIR/$rel"
    agents_file="$AGENTS_DIR/$rel"

    # 两侧任一缺失 → drift
    if [ ! -f "$claude_file" ] || [ ! -f "$agents_file" ]; then
        drift=1
        if [ "$cmd" = "sync-claude-to-agents" ] && [ -f "$claude_file" ]; then
            echo "sync: create $rel  (.claude → .agents)"
            mkdir -p "$(dirname "$agents_file")"
            cp "$claude_file" "$agents_file"
        else
            echo "mirror-missing: $rel (claude=$([ -f "$claude_file" ] && echo Y || echo N) agents=$([ -f "$agents_file" ] && echo Y || echo N))" >&2
        fi
        continue
    fi

    if cmp -s "$claude_file" "$agents_file"; then
        continue    # 已一致
    fi

    drift=1
    if [ "$cmd" = "sync-claude-to-agents" ]; then
        echo "sync: $rel  (.claude → .agents)"
        cp "$claude_file" "$agents_file"
    else
        echo "mirror-drift: $rel" >&2
    fi
done < "$ALLOWLIST"

if [ "$drift" -eq 1 ] && [ "$cmd" = "check" ]; then
    {
        echo ""
        echo "skills allowlist mirror out of sync between .claude and .agents."
        echo "Files under check are listed in bootstrap/skills-mirror-allowlist."
        echo "Resolve by one of:"
        echo "  1) bash bootstrap/skills-mirror-check.sh sync-claude-to-agents"
        echo "     (以 .claude 为准强覆盖 .agents 侧 allowlist 文件)"
        echo "  2) 手工同步差异,如果 codex 侧有意保留独有内容"
        echo "然后 git add 变更并重跑 commit。"
        echo ""
        echo "Note:allowlist 只应含 agent-specific 无差异的纯内容文件"
        echo "(SKILL.md 入口文案、AGENTS.md/CLAUDE.md 引用、~/.Codex/vs ~/.claude/"
        echo " 路径、Codex/Claude 水印所在文件禁止入 allowlist)。"
    } >&2
    exit 1
fi

exit 0
