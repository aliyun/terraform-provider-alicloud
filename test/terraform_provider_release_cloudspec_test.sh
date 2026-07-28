#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
skill="$repo_root/.claude/skills/terraform-provider-release/SKILL.md"
reference="$repo_root/.claude/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
agent="$repo_root/.claude/agents/terraform-rd.md"

test -f "$reference"

# 接单初期必须以 pre 为真源做需求一致性闸门，歧义时四人同时介入。
grep -q 'Step 1.5: CloudSpec Pre-Resource Alignment Gate' "$skill"
grep -q 'cloudspec-pre-resource-loop.md' "$skill"
grep -q -- '--env pre' "$reference"
for contact in '辰羿(320687)' '临钧(429768)' '过载(484483)' '原根(265607)'; do
  grep -q "$contact" "$reference"
done

# 测试期定义问题必须在 pre 修复发布后，再强制从 pre 元数据生成。
grep -q 'amp publish pre --dry-run' "$reference"
grep -q 'amp publish pre' "$reference"
grep -q 'GENERATOR_META_ENV=pre' "$reference"
grep -q '禁止回退 online' "$reference"
grep -q 'cloudspec-resource-edit' "$reference"
grep -q 'cloudspec-operation-edit' "$reference"
grep -q 'cloudspec-build-fix' "$reference"
grep -q 'cloudspec-norm-check-fix' "$reference"
grep -q 'CloudSpec 原主单自闭环' "$reference"
grep -q 'release/idle' "$reference"
grep -q '不得 finish' "$reference"
grep -q 'prod/online' "$reference"
grep -q 'master/main' "$reference"

# RD 机器人必须显式知道仓库已提供 CloudSpec 技能。
grep -q 'cloudspec-amp-workflow' "$agent"
grep -q 'cloudspec-resource-edit' "$agent"
grep -q 'cloudspec-operation-edit' "$agent"

echo "terraform_provider_release_cloudspec_test: PASS"
