#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

for skill in \
  "$repo_root/.agents/skills/writing-weekly-report/SKILL.md" \
  "$repo_root/.claude/skills/writing-weekly-report/SKILL.md"; do
  grep -Fq '个人归因门：参与且产生结果' "$skill"
  grep -Fq '先确定具体汇报主体' "$skill"
  grep -Fq 'staff list <姓名>' "$skill"
  grep -Fq '本人参与证据' "$skill"
  grep -Fq '后续结果证据' "$skill"
  grep -Fq '因果时序成立' "$skill"
  grep -Fq '服务账号的执行结果不自动归属于汇报主体' "$skill"
  grep -Fq '只有本人有代码提交证据时，才写“开发”“实现”' "$skill"
  grep -Fq '项目看板和标题搜索只生成候选清单' "$skill"
  grep -Fq 'sanitize_round_trip=True' "$skill"
  grep -Fq '没有明确量化依据时，不主观修改 KR 百分比' "$skill"

  if grep -Fq '钉钉文档读写能力（`dingtalk-doc-rw`）**当前未落地**' "$skill"; then
    echo "stale DingTalk missing-capability statement remains in $skill" >&2
    exit 1
  fi

  if grep -Eq '辰羿|320687|ChenHanZhang|chenhanzhang' "$skill"; then
    echo "hard-coded report subject remains in $skill" >&2
    exit 1
  fi
done

echo "writing weekly report attribution rules: PASS"
