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
    "2169561" \
    "纯 datasource source-only 契约" \
    "仅涉及 \`data.alicloud_xxx\` 的查询、过滤、分页、输出字段或 Read" \
    "resource+datasource 混合诉求、G Provider 全局改造、手写 resource D 均不属于 pure datasource" \
    "紧急源单 assignee=新山（521957）" \
    "非紧急源单 assignee=过载（484483）" \
    "Jarvis/TerraformRD 在源单直接开发" \
    "严禁为 pure datasource create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766" \
    "历史 relation 只读保留" \
    "不删、不迁、不关、不改派" \
    "不是开发、完成或 blocker 门" \
    "允许引用已有 PR 防重复" \
    "RD route phase 只幂等同步源单 assignee + per-type progress_status" \
    "bridge executor 独占源单 claim/唯一回复/tag/release/finish" \
    "CI pending/fail 或 QA fail 均回 RD 修复" \
    "open PR + QA pass 时源单 release，不 finish" \
    "G 与所有非-datasource D 保留 528766" \
    "I/E/D-临钧/A/F/H 不变" \
    "G / 紧急非-datasource D 的双 owner 契约" \
    "源客户主单 assignee 保持新山（521957）" \
    "528766 研发关联单 assignee 固定过载（484483）" \
    "同题单原地复用并幂等改派过载" \
    "healthy existing claim 不抢占" \
    "Jarvis/TerraformRD claim 并尝试修复" \
    "relation/assignee/status 不是完成信号" \
    "无 PR/CI/QA 完成信号必须继续开发" \
    "build/test/CI 失败只走 RD ↔ QA 修复闭环" \
    "不得转交新山" \
    "missing_capability / retry exhausted" \
    "blocked / SUSPENDED" \
    "不得 finish" \
    "源单与实际 claim 的 528766 各自最多一次聚合 bookend" \
    "PR 未合并只 release"; do
    if ! grep -q "$term" "$skill"; then
      echo "provider_internal_aone_rules_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done

  for stale in \
    "Step 3 · 分支 D-新山" \
    "指派其他人的关联单不 claim" \
    "紧急普通 D（纯 datasource；或" \
    "G/紧急普通 D hard gate 只新增双 owner"; do
    if grep -Fq "$stale" "$skill"; then
      echo "provider_internal_aone_rules_test: stale handoff rule '$stale' in $skill" >&2
      exit 1
    fi
  done
done

for verification in \
  "$repo_root/.agents/skills/provider-resource-dev/references/zhenyuan-verification.md" \
  "$repo_root/.claude/skills/provider-resource-dev/references/zhenyuan-verification.md"; do
  for term in \
    "纯 datasource source-only 契约" \
    "仅涉及 \`data.alicloud_xxx\` 的查询、过滤、分页、输出字段或 Read" \
    "resource+datasource 混合诉求、G Provider 全局改造、手写 resource D 均不属于 pure datasource" \
    "紧急源单 assignee=新山（521957）" \
    "非紧急源单 assignee=过载（484483）" \
    "严禁为 pure datasource create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766" \
    "历史 relation 只读保留" \
    "RD route phase 只幂等同步源单 assignee + per-type progress_status" \
    "bridge executor 独占源单 claim/唯一回复/tag/release/finish" \
    "CI pending/fail 或 QA fail 均回 RD 修复" \
    "open PR + QA pass 时源单 release，不 finish" \
    "G 与所有非-datasource D 保留 528766" \
    "I/E/D-临钧/A/F/H 不变" \
    "G / 紧急非-datasource D 的双 owner 契约" \
    "源客户主单 assignee 保持新山（521957）" \
    "528766 研发关联单 assignee 固定过载（484483）" \
    "同题 528766 已指派新山时原地复用" \
    "healthy existing claim 不抢占" \
    "Jarvis/TerraformRD claim 并尝试修复" \
    "无 PR/CI/QA 完成信号必须继续 RD" \
    "build/test/CI 失败" \
    "不得转交新山" \
    "blocked / SUSPENDED" \
    "D-临钧/A/F/H/非紧急非-datasource D 边界不变"; do
    if ! grep -Fq "$term" "$verification"; then
      echo "provider_internal_aone_rules_test: missing '$term' in $verification" >&2
      exit 1
    fi
  done

  for stale in \
    "D-新山" \
    "紧急普通 D（纯 datasource；或" \
    "无该注释(手写)**或纯 datasource 问题**"; do
    if grep -Fq "$stale" "$verification"; then
      echo "provider_internal_aone_rules_test: stale datasource carrier rule '$stale' in $verification" >&2
      exit 1
    fi
  done
done

for rd_agent in \
  "$repo_root/.claude/agents/terraform-rd.md" \
  "$repo_root/.codex/agents/terraform-rd.toml"; do
  for term in \
    '纯 datasource source-only 契约' \
    'RD route phase 只幂等同步源单 assignee + per-type progress_status' \
    'bridge executor 独占源单 claim/唯一回复/tag/release/finish' \
    '严禁为 pure datasource create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766' \
    '历史 relation 只读保留' \
    'CI pending/fail 或 QA fail 均回 RD 修复' \
    'open PR + QA pass 时源单 release，不 finish' \
    '同一 terraform-rd Task 在 dev 前进入 route-finalizer phase' \
    '仅 materialize/claim 528766' \
    '不得评论源工单' \
    'claim 成功后才切 dev' \
    '最终 finalizer 再 bookend 该研发单' \
    '非紧急非-datasource D 与 I 的 Provider docs 紧急兜底腿保持既有内部 claim/bookend'; do
    if ! grep -Fq "$term" "$rd_agent"; then
      echo "provider_internal_aone_rules_test: missing RD route phase '$term' in $rd_agent" >&2
      exit 1
    fi
  done

  if grep -Fq '开发阶段不写 Aone/钉钉；主处理 run 只有 finalizer 可执行最终外部动作' \
      "$rd_agent"; then
    echo "provider_internal_aone_rules_test: blanket RD external-write ban remains in $rd_agent" >&2
    exit 1
  fi
done

for qa_agent in \
  "$repo_root/.claude/agents/terraform-qa.md" \
  "$repo_root/.codex/agents/terraform-qa.toml"; do
  for term in \
    '纯 datasource source-only 契约' \
    'CI pending/fail 或 QA fail 均回 RD 修复' \
    '不得标为 blocked' \
    'open PR + QA pass 时源单 release，不 finish' \
    'CI fail 或 pending 都返回 `status: fail`' \
    '`next=terraform-rd/fix`' \
    'CI 未就绪不得标为 blocked' \
    'blocked 仅用于 `missing_capability`、`retry exhausted`、明确外部依赖或人工决策'; do
    if ! grep -Fq "$term" "$qa_agent"; then
      echo "provider_internal_aone_rules_test: missing QA state rule '$term' in $qa_agent" >&2
      exit 1
    fi
  done

  for stale in \
    '红或 pending 直接返回 `blocked`' \
    '依赖或 CI 未就绪'; do
    if grep -Fq "$stale" "$qa_agent"; then
      echo "provider_internal_aone_rules_test: stale QA blocked rule '$stale' in $qa_agent" >&2
      exit 1
    fi
  done
done

echo "provider_internal_aone_rules_test: PASS"
