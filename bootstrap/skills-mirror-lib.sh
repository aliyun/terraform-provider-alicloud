#!/usr/bin/env bash
# bootstrap/skills-mirror-lib.sh — .claude/skills ↔ .agents/skills 双份镜像的
# 共享 sed transform 规则。**不是可执行脚本;只通过 source 引入**。
#
# 消费者:
#   bootstrap/sync-to-codex.sh       (Claude → Codex 单向)
#   bootstrap/sync-to-claude.sh      (Codex → Claude 单向)
#   bootstrap/skills-mirror-check.sh (双向 in-memory 对比)
#
# 关键词替换规则(顺序敏感:更长的先替换避免误伤;两个方向互为逆):
#   Claude Code             ↔ Codex
#   Co-Authored-By: Claude  ↔ Co-Authored-By: Codex
#   .claude/agents/         ↔ .Codex/agents/
#   .claude/skills/         ↔ .Codex/skills/     (含 ~/.claude/skills/ 等硬编码路径)
#   claude-code-guide       ↔ codex-guide
#   CLAUDE.md               ↔ AGENTS.md
#
# 使用:
#   source "$(dirname "${BASH_SOURCE[0]}")/skills-mirror-lib.sh"
#   ... | mirror_sed_claude_to_codex
#   ... | mirror_sed_codex_to_claude

# Claude 侧文本 → Codex 侧文本
mirror_sed_claude_to_codex() {
    sed \
        -e 's|Claude Code|Codex|g' \
        -e 's|Co-Authored-By: Claude|Co-Authored-By: Codex|g' \
        -e 's|\.claude/agents/|.Codex/agents/|g' \
        -e 's|\.claude/skills/|.Codex/skills/|g' \
        -e 's|claude-code-guide|codex-guide|g' \
        -e 's|CLAUDE\.md|AGENTS.md|g'
}

# Codex 侧文本 → Claude 侧文本(mirror_sed_claude_to_codex 的逆)
mirror_sed_codex_to_claude() {
    sed \
        -e 's|Co-Authored-By: Codex|Co-Authored-By: Claude|g' \
        -e 's|\.Codex/agents/|.claude/agents/|g' \
        -e 's|\.Codex/skills/|.claude/skills/|g' \
        -e 's|codex-guide|claude-code-guide|g' \
        -e 's|AGENTS\.md|CLAUDE.md|g' \
        -e 's|Codex|Claude Code|g'
}
