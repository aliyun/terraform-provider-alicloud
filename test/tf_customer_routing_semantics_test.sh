#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

source_route="$repo_root/.claude/skills/aone-triage/references/tf-customer-request-routing.md"
mirror_route="$repo_root/.agents/skills/aone-triage/references/tf-customer-request-routing.md"
source_acube_workflow="$repo_root/.claude/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
mirror_acube_workflow="$repo_root/.agents/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
pools_config="$repo_root/config/pools.json"

for file in "$source_route" "$mirror_route"; do
  test -f "$file"

  for phrase in \
    'Existing-related 状态机' \
    '分支 A 不要求关联单' \
    'D/G/H 目标池缺失' \
    '禁止重复触发 `createBuildTaskV2`' \
    '也禁止重复改派或重复回复' \
    'relation 齐全但源单映射字段漂移' \
    '只幂等修复 assignee/status' \
    '不 create、不触发 Acube、不重复阶段回复' \
    '映射字段一致后才进入观察等待' \
    '观察等待' \
    '不评论、不改状态、不改派、不 create、不触发 Acube' \
    '距上次实质进展 ≥8 天' \
    '追料/补料' \
    '终局收敛' \
    '缺关联单优先于 8 天催办'; do
    grep -Fq "$phrase" "$file"
  done

  for phrase in \
    '分支 I — CloudSpec 文档文本 metadata' \
    'text-only' \
    '2169561' \
    '念依（373108）' \
    '独立 528766 紧急兜底腿' \
    '分池防重' \
    '分支 E — CloudSpec 结构 metadata 原主单自闭环' \
    'E → D-临钧' \
    'pre 未收敛不得触发 Acube' \
    '不得由 E 直接执行 Provider PR/CI/ACC' \
    '不得在 E 完成后直接 release/idle' \
    '不得泛化到 A/F/G/H/I、纯 datasource 或纯手写 Provider-only bug'; do
    grep -Fq "$phrase" "$file"
  done

  for phrase in \
    'config/pools.json' \
    'progress_status' \
    '需求问题 → `问题解决中`' \
    '功能缺陷/线上问题 → `Open`' \
    '任务 → `处理中`' \
    '同步源单 assignee'; do
    grep -Fq "$phrase" "$file"
  done

  for phrase in \
    '非 PR 终局结果表' \
    '已发布待需求方验收' \
    '已拒绝' \
    '方案功能已存在' \
    'Fixed' \
    '保持最后处理人' \
    'pr_merged_status' \
    '已合入主线' \
    'merged 即已发布待验收' \
    '先查当前合法枚举' \
    '`.claim.done_statuses`' \
    '`.pools.tf_customer.exclude_status`' \
    '三重门' \
    '不得先写 `客户未响应` 再调用 finish' \
    '无法选择则返回 blocked'; do
    grep -Fq "$phrase" "$file"
  done

  for phrase in \
    '--project 528766' \
    'category 跟随源单 workitemType' \
    '缺陷类型强制紧急' \
    '源单 DDL - 2 天' \
    '不得早于 tomorrow' \
    '无 DDL 时 today + 3 天' \
    '计划开始日期' \
    '计划截止日期' \
    '实际工时=0' \
    'Terraform需求类型'; do
    grep -Fq -- "$phrase" "$file"
  done

  for phrase in \
    '[acube-createBuildTaskV2-workflow.md](./acube-createBuildTaskV2-workflow.md)' \
    '`queryAoneByTaskId`' \
    '60 秒' \
    'executor 托管' \
    'RD finalizer'; do
    grep -Fq "$phrase" "$file"
  done

  for phrase in \
    'H 分支标签合并保护' \
    '`.fields[].tag.value`' \
    'merge existing tag IDs' \
    '保留 `jarvis-idle`' \
    '保留全部业务 tag' \
    '禁止裸名称覆盖'; do
    grep -Fq "$phrase" "$file"
  done
done

for file in "$source_acube_workflow" "$mirror_acube_workflow"; do
  test -f "$file"

  for phrase in \
    'Existing-related 状态机' \
    '已有 `taskId/aoneId`' \
    '只查询/复用' \
    'executor 托管' \
    'terraform-rd finalizer' \
    '禁止回退 jarvis' \
    'config/pools.json' \
    '.pools.tf_customer.progress_status[workitemType]' \
    '缺 mapping' \
    'blocked' \
    'relation 只写一次' \
    'assignee=临钧（429768）' \
    'AONE_RESULT' \
    '最终聚合'; do
    grep -Fq "$phrase" "$file"
  done

  if grep -Fq 'bin/a1id --' "$file" ||
     grep -Fq 'login jarvis' "$file" ||
     grep -Fq -- '--status 问题解决中' "$file"; then
    echo "stale Acube write path remains in $file" >&2
    exit 1
  fi
done

test "$(jq -r '.pools.tf_customer.progress_status["性能瓶颈"]' "$pools_config")" = 'Open'
test "$(jq -r '.pools.tf_customer.done_status["性能瓶颈"]' "$pools_config")" = 'Fixed'
jq -e '
  (.upstream.cloudspec_gap? == null)
  and (.upstream.cloudspec_docs_quality.project == 2169561)
  and (.upstream.cloudspec_docs_quality.assignee == 373108)
  and (.upstream.cloudspec_docs_quality.access == "submit_only")
' "$pools_config" >/dev/null

echo "tf_customer_routing_semantics_test: PASS"
