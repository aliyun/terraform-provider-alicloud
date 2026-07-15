#!/usr/bin/env bash
# bootstrap/agents-mirror-lib.sh — .claude/agents/*.md → .codex/agents/*.toml 生成规则。
# **不是可执行脚本;只通过 source 引入**(与 skills-mirror-lib.sh 同款约定)。
#
# 真源与方向:.claude/agents/<name>.md 单向生成 .codex/agents/<name>.toml——
#   - name / description ← md frontmatter(description 的 YAML `>-` 折叠语义 = 各行以单空格连接)
#   - developer_instructions ← frontmatter 之后的正文(verbatim;TOML """ 多行串,转义 `\` 与 `"""`)
#   - frontmatter 的 tools/skills/model 是 Claude 侧专属字段,不进 toml
# 反向(toml→md)不支持:toml 是生成物,改 md 后重生成即可。
#
# 消费者:
#   bootstrap/mirror.sh {to-codex|check}   (含 PostToolUse hook / pre-commit 门禁链路)
#   test/codex_agents_sync_test.sh

# 输出 <md> 生成的 toml 内容到 stdout;frontmatter 缺失/未闭合/无 name 时退非零。
agents_md_to_toml() {
    local md="$1"
    python3 - "$md" <<'PY'
import sys

path = sys.argv[1]
lines = open(path, encoding="utf-8").read().split("\n")
if not lines or lines[0].strip() != "---":
    sys.stderr.write("agents-mirror: %s 缺 frontmatter\n" % path); sys.exit(1)
fm, i = [], 1
while i < len(lines) and lines[i].strip() != "---":
    fm.append(lines[i]); i += 1
if i >= len(lines):
    sys.stderr.write("agents-mirror: %s frontmatter 未闭合\n" % path); sys.exit(1)
body = "\n".join(lines[i + 1:]).lstrip("\n")
if not body.endswith("\n"):
    body += "\n"

name, desc_lines, j = "", [], 0
while j < len(fm):
    line = fm[j]
    if line.startswith("name:"):
        name = line.split(":", 1)[1].strip()
    elif line.startswith("description:"):
        val = line.split(":", 1)[1].strip()
        if val in (">-", ">", "|", "|-"):
            j += 1
            while j < len(fm) and (fm[j].startswith("  ") or fm[j].strip() == ""):
                if fm[j].strip():
                    desc_lines.append(fm[j].strip())
                j += 1
            continue
        desc_lines = [val]
    j += 1
if not name:
    sys.stderr.write("agents-mirror: %s frontmatter 无 name\n" % path); sys.exit(1)
description = " ".join(desc_lines)

esc = lambda s: s.replace("\\", "\\\\").replace('"', '\\"')
body = body.replace("\\", "\\\\").replace('"""', '""\\"')
sys.stdout.write(
    'name = "%s"\n' % esc(name)
    + 'description = "%s"\n' % esc(description)
    + 'developer_instructions = """\n%s"""\n' % body
)
PY
}
