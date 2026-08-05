#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

active_rules=(
  ".claude/skills/aone-triage/SKILL.md"
  ".agents/skills/aone-triage/SKILL.md"
  ".claude/skills/aone-triage/references/tf-customer-request-routing.md"
  ".agents/skills/aone-triage/references/tf-customer-request-routing.md"
  ".claude/skills/aone-triage/references/team-roster.md"
  ".agents/skills/aone-triage/references/team-roster.md"
  ".claude/skills/aone-triage/references/templates.md"
  ".agents/skills/aone-triage/references/templates.md"
  ".claude/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
  ".agents/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
  ".claude/skills/provider-resource-dev/SKILL.md"
  ".agents/skills/provider-resource-dev/SKILL.md"
  ".claude/skills/provider-resource-dev/references/zhenyuan-verification.md"
  ".agents/skills/provider-resource-dev/references/zhenyuan-verification.md"
  ".claude/skills/invoke-terraform-acc-test-remote/SKILL.md"
  ".agents/skills/invoke-terraform-acc-test-remote/SKILL.md"
  ".claude/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
  ".agents/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
  ".claude/skills/terraform-provider-release/SKILL.md"
  ".agents/skills/terraform-provider-release/SKILL.md"
  ".claude/skills/terraform-provider-release/references/cloudspec-new-resource-infer.md"
  ".agents/skills/terraform-provider-release/references/cloudspec-new-resource-infer.md"
  ".claude/skills/terraform-provider-release/references/requirement-to-openapi-discovery.md"
  ".agents/skills/terraform-provider-release/references/requirement-to-openapi-discovery.md"
  "loops/aone-triage.md"
  "loops/persona-collab.md"
  ".claude/agents/terraform-pd.md"
  ".codex/agents/terraform-pd.toml"
  ".claude/agents/terraform-rd.md"
  ".codex/agents/terraform-rd.toml"
  ".claude/agents/terraform-qa.md"
  ".codex/agents/terraform-qa.toml"
  "CLAUDE.md"
  "AGENTS.md"
  "autonomy.md"
)

require_terms() {
  local file="$1"
  shift
  local term
  for term in "$@"; do
    grep -Fq -- "$term" "$file" || {
      echo "zhenyuan_self_close_rules_test: missing '$term' in $file" >&2
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
      echo "zhenyuan_self_close_rules_test: stale '$term' in $file" >&2
      exit 1
    fi
  done
}

for rel in "${active_rules[@]}"; do
  test -f "$repo_root/$rel"
done

# Retired Zhenyuan identity and project protections from the original suite.
active_paths=()
for rel in "${active_rules[@]}"; do
  active_paths+=("$repo_root/$rel")
done
for forbidden in \
  'WORKER_1783326253279' \
  '镇元 agent' \
  '镇元agent' \
  '谜拟'; do
  if grep -Fn "$forbidden" "${active_paths[@]}"; then
    echo "zhenyuan_self_close_rules_test: active rules still route through '$forbidden'" >&2
    exit 1
  fi
done

while IFS= read -r hit; do
  case "$hit" in
    *删除*|*移除*|*不存在*|*不得*|*不创建*) ;;
    *)
      echo "zhenyuan_self_close_rules_test: positive upstream.cloudspec_gap use: $hit" >&2
      exit 1
      ;;
  esac
done < <(grep -Fn 'upstream.cloudspec_gap' "${active_paths[@]}" || true)

while IFS= read -r hit; do
  case "$hit" in
    *不*|*禁止*|*不得*|*旧*|*回退*) ;;
    *)
      echo "zhenyuan_self_close_rules_test: positive 2165097 route: $hit" >&2
      exit 1
      ;;
  esac
done < <(grep -Fn '2165097' "${active_paths[@]}" || true)

jq -e '
  (.upstream.cloudspec_gap? == null)
  and (.upstream.cloudspec_docs_quality.project == 2169561)
  and (.upstream.cloudspec_docs_quality.assignee == 373108)
  and (.upstream.cloudspec_docs_quality.access == "submit_only")
' "$repo_root/config/pools.json" >/dev/null

jq -e '
  .contacts[]
  | select(.id == "WORKER_1783326253279")
  | .legacy_inbound_only == true
' "$repo_root/config/contacts.json" >/dev/null
jq -e \
  '.agent_fallbacks.WORKER_1783326253279 == "479782"' \
  "$repo_root/config/contacts.json" >/dev/null
grep -Fq 'legacy inbound only' "$repo_root/config/contacts.json"

routing="$repo_root/.claude/skills/aone-triage/references/tf-customer-request-routing.md"
route_mirror="$repo_root/.agents/skills/aone-triage/references/tf-customer-request-routing.md"
main_skill="$repo_root/.claude/skills/aone-triage/SKILL.md"
main_skill_mirror="$repo_root/.agents/skills/aone-triage/SKILL.md"
templates="$repo_root/.claude/skills/aone-triage/references/templates.md"
templates_mirror="$repo_root/.agents/skills/aone-triage/references/templates.md"
verification="$repo_root/.claude/skills/provider-resource-dev/references/zhenyuan-verification.md"
team_roster="$repo_root/.claude/skills/aone-triage/references/team-roster.md"
team_roster_mirror="$repo_root/.agents/skills/aone-triage/references/team-roster.md"
release_loop="$repo_root/.claude/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
pd_agent="$repo_root/.claude/agents/terraform-pd.md"
pd_mirror="$repo_root/.codex/agents/terraform-pd.toml"
rd_agent="$repo_root/.claude/agents/terraform-rd.md"
rd_mirror="$repo_root/.codex/agents/terraform-rd.toml"
qa_agent="$repo_root/.claude/agents/terraform-qa.md"
qa_mirror="$repo_root/.codex/agents/terraform-qa.toml"
runtime="$repo_root/bridge/aone_tasks.py"
persona="$repo_root/loops/persona-collab.md"
acube="$repo_root/.claude/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
acube_mirror="$repo_root/.agents/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"

# Branch I remains text-only, submit-only, and independently deduplicated.
for file in "$routing" "$route_mirror" "$team_roster" "$team_roster_mirror"; do
  require_terms "$file" \
    'CloudSpec 文档文本 metadata' \
    'resource/property/operation description' \
    '字段解释、NOTE' \
    '枚举文案' \
    '不改变字段集合' \
    '类型' \
    '约束' \
    'CRUD' \
    '2169561' \
    '念依' \
    '373108' \
    'submit_only' \
    '分池防重' \
    '独立' \
    '528766'
done
require_terms "$routing" \
  '一个池已有 relation 不能抑制另一个池的缺失补建' \
  'CloudSpec 文档源正确，Provider 本地文档生成/展示偏差' \
  '分支 D'

# I/D/E classification order and reverse prohibitions.
for file in "$templates" "$templates_mirror" "$pd_agent" "$pd_mirror"; do
  require_terms "$file" \
    'CloudSpec 文档文本 metadata' \
    '分支 I' \
    'Provider 本地文档生成/展示偏差' \
    '分支 D' \
    'CloudSpec 结构 metadata' \
    '分支 E'
  forbid_terms "$file" \
    'CloudSpec 文档文本 metadata → 分支 E' \
    'Provider 本地文档生成/展示偏差 → 分支 I'
done
grep -Eq 'CloudSpec 文档文本 metadata.*分支 I' "$templates"
grep -Eq 'CloudSpec 文档源正确，Provider 本地文档生成/展示偏差.*分支 D' "$templates"
grep -Eq 'CloudSpec 结构 metadata.*分支 E' "$templates"

# Team roster keeps I/H owners and rejects the old mixed-owner schema.
require_terms "$team_roster" \
  '| 场景 | 路由/承接关系 | 工号/项目 |' \
  'Provider 侧全局改造 G' \
  '源单新山；TerraformRD 源单直办，不建 528766、不发 route DM' \
  'pure datasource source-only（紧急）' \
  'pure datasource source-only（非紧急）' \
  'CloudSpec 文档文本 metadata（I）' \
  '念依（2169561 submit_only）' \
  'CloudSpec 结构 metadata（E）' \
  'pre QA 后源单继续 Provider dev/CI/远程 ACC/PR' \
  'D 生成/发布腿（含 E pre 收敛后）' \
  '源单临钧' \
  'D 手写 resource 紧急' \
  '源单新山' \
  'D 手写 resource 非紧急' \
  '源单过载' \
  'NPE 兜底' \
  '夏节' \
  '401498'
forbid_terms "$team_roster" \
  '| 场景 | 源客户主单 owner | 下游/研发 owner |' \
  '源单新山；528766 过载' \
  '521957 / 484483' \
  'E → D-临钧'

if grep -Fq 'OK 四条件' "$routing" "$verification" "$team_roster"; then
  echo 'zhenyuan_self_close_rules_test: stale four-condition CloudSpec OK contract' >&2
  exit 1
fi
require_terms "$team_roster" \
  'CloudSpec 结构 OK 三条件' \
  '测试覆盖度 100%'

# Canned missing-input gate remains fail-closed.
for file in "$routing" "$route_mirror" "$main_skill" "$main_skill_mirror" "$pd_agent" "$pd_mirror"; do
  require_terms "$file" \
    'Canned 缺参前置门' \
    '等待补料' \
    '不得进入正式路由' \
    '只读安全查证' \
    '不得据此形成正式路由结论'
done
for file in "$templates" "$templates_mirror"; do
  require_terms "$file" \
    'Canned 补料等待骨架' \
    '当前不进入正式路由或开发' \
    '等待补料'
done

# Single writer / executor-only bookend ownership remains explicit.
for file in \
  "$routing" "$route_mirror" \
  "$main_skill" "$main_skill_mirror" \
  "$templates" "$templates_mirror" \
  "$pd_agent" "$pd_mirror" \
  "$rd_agent" "$rd_mirror" \
  "$qa_agent" "$qa_mirror" \
  "$team_roster" "$team_roster_mirror"; do
  require_terms "$file" 'single-writer'
done
for file in \
  "$routing" "$route_mirror" \
  "$main_skill" "$main_skill_mirror" \
  "$team_roster" "$team_roster_mirror"; do
  require_terms "$file" 'executor 只负责原主单 bookend'
done
if grep -En \
  'executor.{0,40}执行(上述)?(结构化动作|结构化副作用|状态动作)|finalizer/executor|executor 或 terraform-rd' \
  "$routing" "$route_mirror" "$main_skill" "$main_skill_mirror"; then
  echo 'zhenyuan_self_close_rules_test: downstream actions regressed to executor' >&2
  exit 1
fi

# CloudSpec feature/build/check/pre chain and production hard gates.
for file in "$routing" "$route_mirror" "$release_loop" "$rd_agent" "$rd_mirror"; do
  require_terms "$file" \
    'CloudSpec' \
    'pre Meta 收敛' \
    'Provider' \
    'prod/online' \
    'master/main' \
    '人工硬门'
done
require_terms "$routing" \
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
  'cloudspec_pre_verify' \
  'pre QA pass 后返回 `next=terraform-rd/dev`' \
  '同一源单上下文继续 Provider' \
  '远程 AccTest' \
  '不调用 Acube `createBuildTaskV2`' \
  '不创建/复用/关联/claim/bookend'
require_terms "$release_loop" \
  '只运行一次 `aliyun cspec build`' \
  '逐个、前台、串行' \
  '禁止后台或多 Agent' \
  '全部资源 check 通过后' \
  'verification_mode: cloudspec_pre' \
  'next=terraform-rd/dev' \
  'Provider PR' \
  '远程 ACC' \
  'E 禁止调用 Acube `createBuildTaskV2`'

# PD/RD/QA next schemas remain precise and clarify stays retired.
for file in "$pd_agent" "$pd_mirror"; do
  require_terms "$file" \
    'action: dev | acc_verify | cloudspec_pre_verify | finalize' \
    'requested_external_actions: []' \
    'next=terraform-rd/dev'
  forbid_terms "$file" 'clarify'
done
for file in "$rd_agent" "$rd_mirror"; do
  require_terms "$file" \
    '逐个、前台、' \
    '禁止后台或多 Agent' \
    '全部 check 通过后' \
    'role: terraform-rd | terraform-qa | terraform-rd-finalizer' \
    'action: fix | acc_verify | cloudspec_pre_verify | finalize' \
    '文档源证据不足时返回 `status: blocked`、`next=terraform-rd-finalizer/finalize`' \
    'next=terraform-qa/cloudspec_pre_verify'
  forbid_terms "$file" \
    'clarify' \
    'pre_handoff' \
    'E → D-临钧'
done
for file in "$qa_agent" "$qa_mirror"; do
  require_terms "$file" \
    '逐个串行 check' \
    '后台或多 Agent' \
    '全部 check 通过后' \
    'verification_mode: cloudspec_pre' \
    'build/check/pre Meta 收敛' \
    '不运行远程 AccTest' \
    'next=terraform-rd/dev' \
    '同一源单上下文继续 Provider dev/CI/PR' \
    'CI fail 或 pending 都返回 `status: fail`' \
    'next=terraform-rd/fix' \
    'CI 未就绪不得标为 blocked' \
    'blocked 仅用于 `missing_capability`、`retry exhausted`、明确外部依赖或人工决策'
  forbid_terms "$file" \
    'pre_handoff' \
    'E → D-临钧' \
    'clarify'
done

# Pure datasource contract remains stable across core routing surfaces.
pure_files=(
  "$routing" "$route_mirror"
  "$main_skill" "$main_skill_mirror"
  "$team_roster" "$team_roster_mirror"
  "$rd_agent" "$rd_mirror"
  "$qa_agent" "$qa_mirror"
  "$repo_root/AGENTS.md" "$repo_root/CLAUDE.md"
)
for file in "${pure_files[@]}"; do
  require_terms "$file" \
    '纯 datasource source-only 契约' \
    '528766' \
    '历史 relation' \
    '只读' \
    'open PR + QA pass 时源单 release，不 finish'
done
for file in \
  "$routing" "$route_mirror" \
  "$main_skill" "$main_skill_mirror" \
  "$team_roster" "$team_roster_mirror" \
  "$rd_agent" "$rd_mirror" \
  "$repo_root/AGENTS.md" "$repo_root/CLAUDE.md"; do
  require_terms "$file" \
    '新山（521957）' \
    '过载（484483）' \
    'Jarvis/TerraformRD 在源单直接开发'
done
for file in "$qa_agent" "$qa_mirror"; do
  require_terms "$file" \
    '紧急源单 owner 新山（521957）' \
    '非紧急源单 owner 过载（484483）'
done

# D/G source-only and typed-DM replaces the old dual-owner carrier gate.
dg_files=(
  "$routing" "$route_mirror"
  "$main_skill" "$main_skill_mirror"
  "$team_roster" "$team_roster_mirror"
  "$rd_agent" "$rd_mirror"
  "$repo_root/AGENTS.md" "$repo_root/CLAUDE.md"
  "$persona"
)
for file in "${dg_files[@]}"; do
  require_terms "$file" \
    'D/G source-only + D route DM 契约' \
    'D/E/G' \
    '528766' \
    '历史 relation'
  forbid_terms "$file" \
    'G / 紧急非-datasource D 的双 owner 契约' \
    '528766 研发关联单 assignee 固定过载（484483）' \
    'route-finalizer phase' \
    'E → D-临钧'
done
for file in "$qa_agent" "$qa_mirror"; do
  require_terms "$file" \
    'D/G source-only 验收契约' \
    '新山（521957）' \
    '过载（484483）' \
    '临钧' \
    'QA 不写 owner/通知' \
    '不检查 528766 carrier' \
    '历史 relation'
  forbid_terms "$file" \
    'G / 紧急非-datasource D 的双 owner 契约' \
    'route-finalizer phase' \
    'E → D-临钧'
done
require_terms "$routing" \
  'python3 -m bridge.terraform_route_notify' \
  'terraform-route:d:<subtype>:owner:<staffId>' \
  'durable pending 不阻断' \
  'G Provider 全局改造源单 assignee=新山（521957）' \
  '不发送新增 route DM'
require_terms "$runtime" \
  'D/G source-only runtime hard gate' \
  'python3 -m bridge.terraform_route_notify' \
  'terraform-route:d:<subtype>:owner:<staffId>' \
  '逐个、前台、串行 check' \
  '禁止后台或多 Agent' \
  '全部 check 通过后' \
  '通过 cloudspec_pre_verify 后在同一源单上下文继续 Provider'
forbid_terms "$runtime" \
  'E → D-临钧' \
  'route-finalizer phase'

# H tag merge protection remains active and independent.
require_terms "$routing" \
  'H 分支标签合并保护' \
  '`.fields[].tag.value`' \
  'merge existing tag IDs' \
  '保留 `jarvis-idle`' \
  '保留全部业务 tag' \
  '禁止裸名称覆盖' \
  'H NPE 兜底' \
  '夏节（401498）'

# Acube page is non-executable history only; I/H must use active skills.
for file in "$acube" "$acube_mirror"; do
  require_terms "$file" \
    '历史只读' \
    'D/E/G 禁令' \
    '历史 relation/PR 防重取证' \
    'I/H 边界' \
    '不是 I/H 的执行入口'
  forbid_terms "$file" \
    'POST /' \
    '```bash' \
    'curl ' \
    'relation add ' \
    'workitem update ' \
    '--assignee '
done

echo "zhenyuan_self_close_rules_test: PASS"
