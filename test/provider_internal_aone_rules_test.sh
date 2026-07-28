#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

for skill in \
  "$repo_root/.agents/skills/provider-resource-dev/SKILL.md" \
  "$repo_root/.claude/skills/provider-resource-dev/SKILL.md"; do
  for term in \
    "terraform-alicloud" \
    "528766" \
    "WORKER_1782379562571" \
    "双向关联" \
    "客户主单" \
    "CloudSpec 结构 metadata 原主单自闭环" \
    "E → D-临钧" \
    "createBuildTaskV2" \
    "pre 未收敛不得触发 Acube" \
    "不得由 E 直接执行 Provider PR/CI/ACC" \
    "cloudspec-pre-resource-loop.md" \
    "CloudSpec 文档文本 metadata" \
    "2169561"; do
    if ! grep -q "$term" "$skill"; then
      echo "provider_internal_aone_rules_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done
done

echo "provider_internal_aone_rules_test: PASS"
