#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

require_terms() {
  local file="$1"
  shift
  local term
  for term in "$@"; do
    grep -Fq "$term" "$file" || {
      echo "provider_internal_aone_rules_test: missing '$term' in $file" >&2
      exit 1
    }
  done
}

forbid_terms() {
  local file="$1"
  shift
  local term
  for term in "$@"; do
    if grep -Fq "$term" "$file"; then
      echo "provider_internal_aone_rules_test: stale '$term' in $file" >&2
      exit 1
    fi
  done
}

provider_skills=(
  "$repo_root/.agents/skills/provider-resource-dev/SKILL.md"
  "$repo_root/.claude/skills/provider-resource-dev/SKILL.md"
)
verification_refs=(
  "$repo_root/.agents/skills/provider-resource-dev/references/zhenyuan-verification.md"
  "$repo_root/.claude/skills/provider-resource-dev/references/zhenyuan-verification.md"
)

for skill in "${provider_skills[@]}"; do
  require_terms "$skill" \
    'Aone source-only 路由与同步' \
    '分支 D/E/G 与 pure datasource 均在**源工单**直接开发' \
    'create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766' \
    '历史 relation' \
    '不是开发、完成、阻塞或 observe 门' \
    'pure datasource' \
    '紧急源单 assignee=新山（521957）' \
    '非紧急=过载（484483）' \
    'D 手写 resource' \
    '生成器/release 路径=临钧' \
    'G Provider 全局改造' \
    'G 不发新增路由私信' \
    'verification_mode: cloudspec_pre' \
    '在同一源单上下文继续 Provider' \
    'PR CI' \
    '远程 ACC' \
    'python3 -m bridge.terraform_route_notify' \
    'handwritten-urgent|handwritten-normal|generated' \
    'terraform-route:d:<subtype>:owner:<staffId>' \
    'pending 不阻塞' \
    'posted/suppressed 不重发' \
    'post_uncertain 复用同一 receipt' \
    '持久化失败不得宣称通知完成' \
    'CI pending/fail 或 QA fail 均回 RD 修复' \
    'open PR + QA pass 时 release 源单、不 finish' \
    'prod/online' \
    'master/main merge/push' \
    '正式发布仍是人工硬门' \
    '2169561' \
    '念依（373108）' \
    '独立 528766 兜底腿' \
    '过载（484483）' \
    'H 仍进入 528766 并指派夏节' \
    'PD/QA 不外写' \
    'bridge executor'

  forbid_terms "$skill" \
    'G / 紧急非-datasource D 的双 owner 契约' \
    '528766 研发关联单 assignee 固定过载（484483）' \
    'route-finalizer phase' \
    'E → D-临钧' \
    'hand off through Acube' \
    '已有正确 relation/taskId/aoneId 时只查询/复用' \
    '源单与实际 claim 的 528766 各自最多一次聚合 bookend'
done

for verification in "${verification_refs[@]}"; do
  require_terms "$verification" \
    '前置短路 · 纯 datasource 问题不查镇元' \
    '仅涉及 `data.alicloud_xxx` 的查询、过滤、分页、输出字段或 Read' \
    'resource+datasource 混合诉求、G Provider 全局改造、手写 resource D 均不属于 pure datasource' \
    '分支 D/G 与 pure datasource：源单直接开发' \
    'create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766' \
    '历史 relation' \
    '不是开发、完成、阻塞或 observe 门' \
    '紧急新山（521957）' \
    '非紧急过载（484483）' \
    'G Provider 全局改造的源单 owner=新山（521957）' \
    'G 不新增路由私信' \
    '生成器/release 路径' \
    '源单 owner=临钧（429768）' \
    '手写 resource' \
    '源单 owner=新山（521957）' \
    '源单 owner=过载（484483）' \
    'python3 -m bridge.terraform_route_notify' \
    'terraform-route:d:<subtype>:owner:<staffId>' \
    'pending 不阻塞开发' \
    'post_uncertain 复用同一 receipt' \
    'pre Meta 收敛后由 QA' \
    'verification_mode: cloudspec_pre' \
    '在同一源单上下文继续 Provider' \
    'PR CI' \
    '远程 ACC' \
    'createBuildTaskV2' \
    'open PR + QA pass 时 release 源单、不 finish' \
    '2169561/念依' \
    '独立 528766' \
    'H 仍进入 528766/夏节' \
    'A/F 保持原路由' \
    'amp publish prod' \
    'prod/online' \
    'master/main merge/push' \
    '人工硬门'

  forbid_terms "$verification" \
    'G / 紧急非-datasource D 的双 owner 契约' \
    '528766 研发关联单 assignee 固定过载（484483）' \
    'route-finalizer phase' \
    'E → D-临钧' \
    '同题 528766 已指派新山时原地复用' \
    '不得由 E 直接执行 Provider PR/CI/ACC'
done

for rd_agent in \
  "$repo_root/.claude/agents/terraform-rd.md" \
  "$repo_root/.codex/agents/terraform-rd.toml"; do
  require_terms "$rd_agent" \
    'PD/QA 不外写' \
    'single-writer' \
    'executor' \
    'Canned 缺参前置门' \
    '文档源证据不足时返回 `status: blocked`、`next=terraform-rd-finalizer/finalize`' \
    'role: terraform-rd | terraform-qa | terraform-rd-finalizer' \
    'action: fix | acc_verify | cloudspec_pre_verify | finalize' \
    'CloudSpec 文档文本 metadata' \
    '2169561' \
    'submit_only' \
    '两腿分池' \
    '一个池已有 relation 不能抑制另一个池的缺失补建' \
    '独立 528766 紧急兜底腿' \
    'CloudSpec 结构 metadata 原主单自闭环' \
    'task 专属 feature 分支' \
    'AMP 返回的 SSH URL' \
    'build/check/publish pre' \
    'pre Meta 收敛' \
    'next=terraform-qa/cloudspec_pre_verify' \
    '同一源单上下文继续 Provider dev/CI/PR' \
    '远程 AccTest' \
    '不创建/复用 528766' \
    '纯 datasource source-only 契约' \
    'RD route phase 只幂等同步源单 assignee + per-type progress_status' \
    'bridge executor' \
    'D/G source-only + D route DM 契约' \
    'python3 -m bridge.terraform_route_notify' \
    'terraform-route:d:<subtype>:owner:<staffId>' \
    'G 不发新增 route DM' \
    'open PR + QA pass 时源单 release 不 finish' \
    'prod/online' \
    'master/main merge/push' \
    '人工硬门'

  forbid_terms "$rd_agent" \
    'route-finalizer phase' \
    '仅 materialize/claim 528766' \
    'claim 成功后才切 dev' \
    'E → D-临钧' \
    'clarify'
done

for qa_agent in \
  "$repo_root/.claude/agents/terraform-qa.md" \
  "$repo_root/.codex/agents/terraform-qa.toml"; do
  require_terms "$qa_agent" \
    'PD/QA 不外写' \
    'single-writer' \
    '纯 datasource source-only 契约' \
    'D/G source-only 验收契约' \
    'CI pending/fail 或 QA fail 均回 RD 修复' \
    'next=terraform-rd/fix' \
    '不得标为 blocked' \
    'open PR + QA pass 时源单 release，不 finish' \
    'CI fail 或 pending 都返回 `status: fail`' \
    'CI 未就绪不得标为 blocked' \
    'blocked 仅用于 `missing_capability`、`retry exhausted`、明确外部依赖或人工决策' \
    'verification_mode: cloudspec_pre' \
    'build/check/pre Meta 收敛' \
    '不运行远程 AccTest' \
    'next=terraform-rd/dev' \
    '同一源单上下文继续 Provider dev/CI/PR' \
    'QA 不调用 `createBuildTaskV2`' \
    '不创建/关联/指派 528766'

  forbid_terms "$qa_agent" \
    '红或 pending 直接返回 `blocked`' \
    'E → D-临钧' \
    'pre_handoff' \
    'clarify'
done

echo "provider_internal_aone_rules_test: PASS"
