#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
skill="$repo_root/.claude/skills/terraform-provider-release/SKILL.md"
reference="$repo_root/.claude/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
acube_workflow="$repo_root/.claude/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
new_resource_reference="$repo_root/.claude/skills/terraform-provider-release/references/cloudspec-new-resource-infer.md"
discovery_reference="$repo_root/.claude/skills/terraform-provider-release/references/requirement-to-openapi-discovery.md"
acc_skill="$repo_root/.claude/skills/invoke-terraform-acc-test-remote/SKILL.md"
agent="$repo_root/.claude/agents/terraform-rd.md"

test -f "$reference"

# 接单初期必须以 pre 为真源做需求一致性闸门，歧义时四人同时介入。
grep -q 'Step 1.5: CloudSpec Pre-Resource Alignment Gate' "$skill"
grep -q 'cloudspec-pre-resource-loop.md' "$skill"
grep -q -- '--env pre' "$reference"
for contact in '辰羿(320687)' '临钧(429768)' '过载(484483)' '原根(265607)'; do
  grep -q "$contact" "$reference"
done

# 分支 E 只修到 pre 收敛并强制交接 D-临钧；普通 D 的 pre 生成路径仍保留。
grep -q 'amp publish pre --dry-run' "$reference"
grep -q 'amp publish pre' "$reference"
grep -q 'cloudspec-resource-edit' "$reference"
grep -q 'cloudspec-operation-edit' "$reference"
grep -q 'cloudspec-build-fix' "$reference"
grep -q 'cloudspec-norm-check-fix' "$reference"
grep -q 'CloudSpec 结构 metadata 原主单自闭环' "$reference"
grep -q 'E → D-临钧' "$reference"
grep -q 'createBuildTaskV2' "$reference"
grep -q 'pre Meta 收敛' "$reference"
grep -q 'pre 未收敛不得触发 Acube' "$reference"
grep -q '不得由 E 直接执行 Provider PR/CI/ACC' "$reference"
grep -q '不得在 E 完成后直接 release/idle' "$reference"
grep -q '已有正确 relation/taskId/aoneId 时只查询/复用' "$reference"
grep -q '不得 finish' "$reference"
grep -q 'prod/online' "$reference"
grep -q 'master/main' "$reference"
grep -q '普通分支 D' "$reference"

for file in "$skill" "$new_resource_reference" "$discovery_reference" "$acc_skill"; do
  grep -q '分支 E' "$file"
  grep -q 'pre Meta 收敛' "$file"
  grep -q 'E → D-临钧' "$file"
  grep -q 'createBuildTaskV2' "$file"
  grep -q '不得由 E 直接执行 Provider PR/CI/ACC' "$file"
  grep -q '普通分支 D' "$file"
done

for phrase in \
  '分支 E 的 pre Meta 已收敛' \
  'verification_mode: cloudspec_pre' \
  'terraform-rd finalizer 是 downstream single-writer' \
  'executor 只负责原主单 bookend' \
  '已有正确 528766 relation' \
  'createBuildTaskV2' \
  '只查询/复用'; do
  grep -q "$phrase" "$acube_workflow"
done

# RD 机器人必须显式知道仓库已提供 CloudSpec 技能。
grep -q 'cloudspec-amp-workflow' "$agent"
grep -q 'cloudspec-resource-edit' "$agent"
grep -q 'cloudspec-operation-edit' "$agent"

echo "terraform_provider_release_cloudspec_test: PASS"
