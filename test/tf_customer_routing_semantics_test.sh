#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

source_route="$repo_root/.claude/skills/aone-triage/references/tf-customer-request-routing.md"
mirror_route="$repo_root/.agents/skills/aone-triage/references/tf-customer-request-routing.md"
source_acube="$repo_root/.claude/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
mirror_acube="$repo_root/.agents/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
source_batch="$repo_root/.claude/skills/aone-triage/references/batch-bookend-template.md"
mirror_batch="$repo_root/.agents/skills/aone-triage/references/batch-bookend-template.md"
pools="$repo_root/config/pools.json"

require_terms() {
  local file="$1"
  shift
  local term
  for term in "$@"; do
    grep -Fq -- "$term" "$file" || {
      echo "tf_customer_routing_semantics_test: missing '$term' in $file" >&2
      exit 1
    }
  done
}

forbid_terms() {
  local file="$1"
  shift
  local term
  for term in "$@"; do
    if grep -Fq -- "$term" "$file"; then
      echo "tf_customer_routing_semantics_test: stale '$term' in $file" >&2
      exit 1
    fi
  done
}

for route in "$source_route" "$mirror_route"; do
  test -f "$route"

  # Existing-related, source-only and ownership invariants.
  require_terms "$route" \
    'Existing-related 状态机' \
    '分支 A 不要求关联单' \
    '纯 datasource source-only 契约' \
    '仅涉及 `data.alicloud_xxx` 的查询、过滤、分页、输出字段或 Read' \
    'resource+datasource 混合诉求、G Provider 全局改造、手写 resource D 均不属于 pure datasource' \
    '紧急源单 assignee=新山（521957）' \
    '非紧急源单 assignee=过载（484483）' \
    'Jarvis/TerraformRD 在源单直接开发' \
    '严禁为 pure datasource create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766' \
    '历史 relation 只读保留' \
    '不删、不迁、不关、不改派' \
    '允许引用已有 PR 防重复' \
    'RD route phase 只幂等同步源单 assignee + per-type progress_status' \
    'bridge executor 独占源单 claim/唯一回复/tag/release/finish' \
    'CI pending/fail 或 QA fail 均回 RD 修复' \
    'open PR + QA pass 时源单 release，不 finish' \
    'D/E/G 同样严禁 528766 承载' \
    'D/G source-only + D route DM 契约' \
    '手写 resource 紧急源单 assignee=新山（521957）' \
    '手写非紧急 assignee=过载（484483）' \
    '生成/发布腿（含 E pre 收敛后）assignee=临钧（429768）' \
    'G Provider 全局改造源单 assignee=新山（521957）' \
    '不发送新增 route DM' \
    'D/E/G 均由 Jarvis/TerraformRD 在源单上下文主动开发' \
    'create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766' \
    '不得因 owner/status、通知或历史 relation 观察等待' \
    'python3 -m bridge.terraform_route_notify' \
    'handwritten-urgent|handwritten-normal|generated' \
    'terraform-route:d:<subtype>:owner:<staffId>' \
    'ticket 参与 ledger id' \
    'durable pending 不阻断' \
    'posted/suppressed 不重发' \
    'post_uncertain 保持同一 receipt' \
    'ledger 无法持久化不得' \
    'open PR + QA pass 时源单 release，不 finish' \
    'prod/online' \
    'master/main merge/push' \
    '正式发布仍是人工硬门'

  # I/D/E classification and reverse protections.
  require_terms "$route" \
    '分支 I — CloudSpec 文档文本 metadata' \
    'text-only' \
    'resource/property/operation description' \
    '字段解释、NOTE 与枚举文案' \
    '不改变字段集合、类型、约束或 CRUD' \
    '2169561' \
    '念依（373108）' \
    'submit_only' \
    '独立 528766 紧急兜底腿' \
    '分池防重' \
    '一个池已有 relation 不能抑制另一个池的缺失补建' \
    'CloudSpec 文档源正确，Provider 本地文档生成/展示偏差' \
    '分支 D' \
    '分支 E — CloudSpec 结构 metadata 原主单自闭环' \
    'requested_external_actions: []' \
    'next=terraform-rd/dev' \
    'cloudspec-amp-workflow' \
    'cloudspec-idl-guide' \
    'cloudspec-resource-edit' \
    'cloudspec-operation-edit' \
    'cloudspec-build-fix' \
    'cloudspec-norm-check-fix' \
    'AMP 返回的 SSH URL' \
    'task 专属 feature 分支' \
    'aliyun cspec build' \
    'aliyun cspec check' \
    'amp publish pre --dry-run' \
    'amp publish pre' \
    'pre Meta 收敛' \
    'next=terraform-qa/cloudspec_pre_verify' \
    'pre QA pass 后返回 `next=terraform-rd/dev`' \
    '同一源单上下文继续 Provider' \
    '远程 AccTest' \
    '不调用 Acube `createBuildTaskV2`' \
    '不创建/复用/关联/claim/bookend' \
    'H 仍创建/复用 528766→夏节' \
    'jarvis-npe'

  # Single writer, per-type progress and H tag merge protection.
  require_terms "$route" \
    'terraform-rd finalizer 是 downstream single-writer' \
    'executor 只负责原主单 bookend' \
    'config/pools.json' \
    '.pools.tf_customer.progress_status[workitemType]' \
    '需求问题 → `问题解决中`' \
    '功能缺陷/线上问题 → `Open`' \
    '性能瓶颈 → `Open`' \
    '任务 → `处理中`' \
    '未映射的 workitemType 必须 blocked' \
    'H 分支标签合并保护' \
    '`.fields[].tag.value`' \
    'merge existing tag IDs' \
    '保留 `jarvis-idle`' \
    '保留全部业务 tag' \
    '禁止裸名称覆盖'

  # Non-PR terminal mapping and customer-no-response triple gate.
  require_terms "$route" \
    '非 PR 终局结果表' \
    '方案功能已存在' \
    '已拒绝' \
    '已发布待需求方验收' \
    'Fixed' \
    '保持最后处理人' \
    '.pools.tf_customer.pr_merged_status' \
    '已合入主线' \
    'merged 即已发布待验收' \
    '先查当前合法枚举' \
    '客户未响应三重门' \
    '`.claim.done_statuses`' \
    '`.pools.tf_customer.exclude_status`' \
    '不得先写 `客户未响应` 再调用 finish' \
    '无法选择则返回 blocked'

  # Canned gate remains independent of the D/E/G change.
  require_terms "$route" \
    'Canned 缺参前置门' \
    '等待补料' \
    '不得进入正式路由' \
    '只读安全查证' \
    '不得据此形成正式路由结论'

  forbid_terms "$route" \
    'G / 紧急非-datasource D 的双 owner 契约' \
    '528766 研发关联单 assignee 固定过载（484483）' \
    'route-finalizer phase' \
    'E → D-临钧' \
    'pre_handoff' \
    '不得由 E 直接执行 Provider PR/CI/ACC' \
    '映射字段一致后才进入观察等待' \
    'D-新山'
done

cmp -s "$source_route" "$mirror_route" || {
  echo "tf_customer_routing_semantics_test: routing mirrors drifted" >&2
  exit 1
}

for batch in "$source_batch" "$mirror_batch"; do
  require_terms "$batch" \
    '骨架 C · pure datasource source-only 源单路由' \
    'PURE_DATASOURCE_ASSIGNEE' \
    'URGENT_ASSIGNEE=521957' \
    'NONURGENT_ASSIGNEE=484483' \
    '从 pools.json progress_status[workitemType] 解析' \
    'SOURCE_ROUTE_DRIFT' \
    '只幂等同步源单 assignee + per-type progress_status' \
    'bridge executor 独占源单 claim/唯一回复/tag/release/finish' \
    '历史 relation 只读保留' \
    '禁止 create/reuse-as-carrier/reassign/relation/claim/wrap/release/finish 528766' \
    '骨架 D · D/G source-only + D route DM' \
    'D/E/G 严禁对 528766 执行 create/reuse/reassign/relation/claim/wrap/release/finish' \
    'd:handwritten-urgent) SOURCE_ASSIGNEE=521957' \
    'd:handwritten-normal) SOURCE_ASSIGNEE=484483' \
    'd:generated) SOURCE_ASSIGNEE=429768' \
    'g:) SOURCE_ASSIGNEE=521957' \
    'python3 -m bridge.terraform_route_notify' \
    '--ticket "$SRC"' \
    '--subtype "$SUBTYPE"' \
    'executor 独占源工单' \
    '最后交 AONE_RESULT' \
    'release 源单，不 finish'

  pure_section="$(awk '
    /^## 骨架 C · pure datasource source-only 源单路由$/ { capture = 1 }
    capture && /^## 骨架 D ·/ { exit }
    capture { print }
  ' "$batch")"
  dg_section="$(awk '
    /^## 骨架 D · D\/G source-only \+ D route DM$/ { capture = 1 }
    capture && /^## / && !/^## 骨架 D ·/ { exit }
    capture { print }
  ' "$batch")"
  test -n "$pure_section"
  test -n "$dg_section"
  for section in "$pure_section" "$dg_section"; do
    for forbidden in \
      'RELATED_ID=' \
      'RELATED_PROJECT=' \
      'relation add' \
      'workitem create' \
      'claim "$RELATED_ID"' \
      'wrap.sh done "$RELATED_ID"' \
      'release "$RELATED_ID"'; do
      if grep -Fq -- "$forbidden" <<<"$section"; then
        echo "tf_customer_routing_semantics_test: source-only carrier action '$forbidden' in $batch" >&2
        exit 1
      fi
    done
  done

  forbid_terms "$batch" \
    '骨架 D · G / 紧急非-datasource D 双 owner' \
    'RELATED_ASSIGNEE=484483' \
    'claim "$RELATED_ID" "$RELATED_PROJECT"' \
    'wrap.sh done "$RELATED_ID"'
done

for acube in "$source_acube" "$mirror_acube"; do
  require_terms "$acube" \
    '历史只读' \
    'D/E/G 禁令' \
    '调用 `createBuildTaskV2`' \
    '创建、复用或承载 528766' \
    '新建或补写 relation' \
    '改派源单或历史关联单' \
    'claim、wrap、release、finish 或其它 bookend' \
    '历史 relation/PR 防重取证' \
    'taskId' \
    'aoneId' \
    'I/H 边界' \
    '不是 I/H 的执行入口'
  forbid_terms "$acube" \
    'POST /' \
    '```bash' \
    'curl ' \
    'relation add ' \
    'workitem update ' \
    '--assignee '
done

cmp -s "$source_acube" "$mirror_acube" || {
  echo "tf_customer_routing_semantics_test: Acube mirrors drifted" >&2
  exit 1
}

test "$(jq -r '.pools.tf_customer.progress_status["性能瓶颈"]' "$pools")" = 'Open'
test "$(jq -r '.pools.tf_customer.done_status["性能瓶颈"]' "$pools")" = 'Fixed'
jq -e '
  (.upstream.cloudspec_gap? == null)
  and (.upstream.cloudspec_docs_quality.project == 2169561)
  and (.upstream.cloudspec_docs_quality.assignee == 373108)
  and (.upstream.cloudspec_docs_quality.access == "submit_only")
' "$pools" >/dev/null

echo "tf_customer_routing_semantics_test: PASS"
