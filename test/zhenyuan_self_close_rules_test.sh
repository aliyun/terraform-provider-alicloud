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

for rel in "${active_rules[@]}"; do
  test -f "$repo_root/$rel"
done

for forbidden in \
  'WORKER_1783326253279' \
  '镇元 agent' \
  '镇元agent' \
  '谜拟'; do
  if grep -Fn "$forbidden" "${active_rules[@]/#/$repo_root/}"; then
    echo "zhenyuan_self_close_rules_test: active rules still route through '$forbidden'" >&2
    exit 1
  fi
done

# The removed config key may appear only in explicit deletion/negative policy text.
while IFS= read -r hit; do
  case "$hit" in
    *删除*|*移除*|*不存在*|*不得*|*不创建*) ;;
    *)
      echo "zhenyuan_self_close_rules_test: active rule positively uses upstream.cloudspec_gap: $hit" >&2
      exit 1
      ;;
  esac
done < <(grep -Fn 'upstream.cloudspec_gap' "${active_rules[@]/#/$repo_root/}" || true)

# The retired project id may appear only in an explicit negative/legacy warning.  A
# blanket literal ban would reject the desired policy text ("不创建 2165097").
while IFS= read -r hit; do
  case "$hit" in
    *不*|*禁止*|*不得*|*旧*|*回退*) ;;
    *)
      echo "zhenyuan_self_close_rules_test: active rule positively routes through 2165097: $hit" >&2
      exit 1
      ;;
  esac
done < <(grep -Fn '2165097' "${active_rules[@]/#/$repo_root/}" || true)

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
jq -e '
  .agent_fallbacks.WORKER_1783326253279 == "479782"
' "$repo_root/config/contacts.json" >/dev/null
grep -Fq 'legacy inbound only' "$repo_root/config/contacts.json"

routing="$repo_root/.claude/skills/aone-triage/references/tf-customer-request-routing.md"
main_skill="$repo_root/.claude/skills/aone-triage/SKILL.md"
templates="$repo_root/.claude/skills/aone-triage/references/templates.md"
verification="$repo_root/.claude/skills/provider-resource-dev/references/zhenyuan-verification.md"
team_roster="$repo_root/.claude/skills/aone-triage/references/team-roster.md"
pd_agent="$repo_root/.claude/agents/terraform-pd.md"
pd_agent_mirror="$repo_root/.codex/agents/terraform-pd.toml"
release_loop="$repo_root/.claude/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
rd_agent="$repo_root/.claude/agents/terraform-rd.md"
rd_agent_mirror="$repo_root/.codex/agents/terraform-rd.toml"
qa_agent="$repo_root/.claude/agents/terraform-qa.md"
qa_agent_mirror="$repo_root/.codex/agents/terraform-qa.toml"
runtime_orchestrator="$repo_root/bridge/aone_tasks.py"
bot_runtime="$repo_root/bridge/jarvis_dingtalk_bot.py"
persona_collab="$repo_root/loops/persona-collab.md"

for phrase in \
  '分支 E — CloudSpec 结构 metadata 原主单自闭环' \
  '字段集合、类型、约束、CRUD' \
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
  'E → D-临钧' \
  'createBuildTaskV2' \
  '临钧（429768）' \
  'pre 未收敛不得触发 Acube' \
  '不得由 E 直接执行 Provider PR/CI/ACC' \
  '不得在 E 完成后直接 release/idle' \
  '已有正确 relation/taskId/aoneId 时只查询/复用' \
  '只允许分支 E 进入此转换' \
  '不得泛化到 A/F/G/H/I、纯 datasource 或纯手写 Provider-only bug' \
  'missing_capability' \
  'blocked' \
  '不得 finish'; do
  grep -Fq "$phrase" "$routing"
done

for phrase in \
  '分支 I — CloudSpec 文档文本 metadata' \
  'resource/property/operation description' \
  '字段解释、NOTE 与枚举文案' \
  '不改变字段集合、类型、约束或 CRUD' \
  '2169561' \
  '念依（373108）' \
  'upstream.cloudspec_docs_quality' \
  'submit_only' \
  '创建或复用' \
  '分池防重' \
  '独立 528766 紧急兜底腿' \
  '一个池已有 relation 不能抑制另一个池的缺失补建'; do
  grep -Fq "$phrase" "$routing"
  grep -Fq "$phrase" "$verification"
done

for phrase in \
  'CloudSpec 结构 OK 三条件' \
  'CloudSpec 文档文本 metadata' \
  '念依' \
  '373108' \
  '结构 metadata' \
  'E → D-临钧'; do
  grep -Fq "$phrase" "$team_roster"
done

for phrase in \
  '| 场景 | 路由/承接关系 | 工号/项目 |' \
  '| **Provider 侧全局改造**' \
  '源单新山；528766 过载并由 TerraformRD 内部开发' \
  '521957 / 484483' \
  '| **CloudSpec 文档文本 metadata（I）**' \
  '| 念依（2169561 submit_only） | 373108 |' \
  '| **CloudSpec 结构 metadata（E）**' \
  '| open-jarvis → 临钧 | 原主单内部执行 → 429768 |' \
  '| **NPE 兜底**' \
  '| 夏节 | 401498 |'; do
  grep -Fq "$phrase" "$team_roster"
done
if grep -Fq '| 场景 | 源客户主单 owner | 下游/研发 owner |' "$team_roster"; then
  echo 'zhenyuan_self_close_rules_test: mixed-owner roster schema remains' >&2
  exit 1
fi

if grep -Fq 'OK 四条件' "$routing" "$verification" "$team_roster"; then
  echo 'zhenyuan_self_close_rules_test: stale four-condition CloudSpec OK contract' >&2
  exit 1
fi

awk '
  index($0, "text-only") && index($0, "分支 I") && !text_split { text_split = NR }
  index($0, "CloudSpec 结构 OK（三条件全满足）") && !ok_check { ok_check = NR }
  END { exit !(text_split && ok_check && text_split < ok_check) }
' "$routing"

for phrase in \
  'CloudSpec 文档文本 metadata' \
  '分支 I' \
  'Provider 本地文档生成/展示偏差' \
  '分支 D' \
  'CloudSpec 结构 metadata' \
  '分支 E'; do
  grep -Fq "$phrase" "$templates"
  grep -Fq "$phrase" "$pd_agent"
done

grep -Eq 'CloudSpec 文档文本 metadata.*分支 I' "$templates"
grep -Eq 'CloudSpec 文档源正确，Provider 本地文档生成/展示偏差.*分支 D' "$templates"
grep -Eq 'CloudSpec 结构 metadata.*分支 E' "$templates"
if grep -Eq 'CloudSpec 文档文本 metadata.*分支 E|Provider 本地文档生成/展示偏差.*分支 I' "$templates"; then
  echo 'zhenyuan_self_close_rules_test: document source/local rendering routes are reversed' >&2
  exit 1
fi

for file in "$routing" "$main_skill" "$pd_agent" "$rd_agent" "$qa_agent"; do
  for phrase in \
    'PD/QA 不外写' \
    'single-writer'; do
    grep -Fq "$phrase" "$file"
  done
done

ownership_files=(
  "$repo_root/.claude/skills/aone-triage/SKILL.md"
  "$repo_root/.agents/skills/aone-triage/SKILL.md"
  "$repo_root/.claude/skills/aone-triage/references/team-roster.md"
  "$repo_root/.agents/skills/aone-triage/references/team-roster.md"
  "$repo_root/.claude/skills/aone-triage/references/tf-customer-request-routing.md"
  "$repo_root/.agents/skills/aone-triage/references/tf-customer-request-routing.md"
  "$repo_root/.claude/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
  "$repo_root/.agents/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"
)
for file in "${ownership_files[@]}"; do
  grep -Fq 'terraform-rd finalizer 是 downstream single-writer' "$file"
  grep -Fq 'executor 只负责原主单 bookend' "$file"
done
if grep -En \
  'executor.{0,40}执行(上述)?(结构化动作|结构化副作用|状态动作)|finalizer/executor|executor 或 terraform-rd|仅是独立 finalizer 示例|绝不能运行这些命令|不得在 executor 托管 run 中执行' \
  "${ownership_files[@]}"; then
  echo 'zhenyuan_self_close_rules_test: downstream single-writer ownership regressed to executor' >&2
  exit 1
fi

for file in "$routing" "$main_skill" "$pd_agent"; do
  for phrase in \
    'Canned 缺参前置门' \
    '等待补料' \
    '不得进入正式路由' \
    '只读安全查证' \
    '不得据此形成正式路由结论'; do
    grep -Fq "$phrase" "$file"
  done
done

for phrase in \
  'Canned 补料等待骨架' \
  '当前不进入正式路由或开发' \
  '等待补料'; do
  grep -Fq "$phrase" "$templates"
done

for phrase in \
  'amp publish prod' \
  'prod/online' \
  'master/main' \
  '人工硬门' \
  'pre 成功'; do
  grep -Fq "$phrase" "$release_loop" "$rd_agent"
done

for file in "$rd_agent" "$rd_agent_mirror"; do
  grep -Fq '文档源证据不足时返回 `status: blocked`、`next=terraform-rd-finalizer/finalize`' "$file"
  grep -Fq 'role: terraform-rd | terraform-qa | terraform-rd-finalizer' "$file"
  grep -Fq 'action: fix | acc_verify | cloudspec_pre_verify | finalize' "$file"
  if grep -Fq 'clarify' "$file"; then
    echo "zhenyuan_self_close_rules_test: stale RD clarify route in $file" >&2
    exit 1
  fi
done

for file in "$pd_agent" "$pd_agent_mirror"; do
  grep -Fq 'action: dev | acc_verify | cloudspec_pre_verify | finalize' "$file"
  if grep -Fq 'clarify' "$file"; then
    echo "zhenyuan_self_close_rules_test: stale PD clarify action in $file" >&2
    exit 1
  fi
done

for file in "$qa_agent" "$qa_agent_mirror"; do
  for phrase in \
    'verification_mode: cloudspec_pre' \
    'build/check/pre Meta 收敛' \
    '不运行远程 AccTest' \
    '不得要求 Provider PR/CI/ACC' \
    'E → D-临钧' \
    'action: fix | pre_handoff | finalize'; do
    grep -Fq "$phrase" "$file"
  done
done

for phrase in \
  'cloudspec_pre_verify' \
  'pre_handoff'; do
  grep -Fq "$phrase" "$persona_collab"
  grep -Fq "$phrase" "$runtime_orchestrator"
  grep -Fq "$phrase" "$bot_runtime"
done

for phrase in \
  '分支 I' \
  '2169561' \
  '分池防重' \
  '分支 E' \
  'pre Meta 收敛' \
  'E → D-临钧' \
  'createBuildTaskV2' \
  '不得由 E 直接执行 Provider PR/CI/ACC' \
  '普通分支 D'; do
  grep -Fq "$phrase" "$runtime_orchestrator"
done

g_urgent_route_files=(
  "$repo_root/.claude/skills/aone-triage/SKILL.md"
  "$repo_root/.agents/skills/aone-triage/SKILL.md"
  "$repo_root/.claude/skills/aone-triage/references/team-roster.md"
  "$repo_root/.agents/skills/aone-triage/references/team-roster.md"
  "$repo_root/.claude/skills/aone-triage/references/tf-customer-request-routing.md"
  "$repo_root/.agents/skills/aone-triage/references/tf-customer-request-routing.md"
  "$repo_root/.claude/skills/provider-resource-dev/SKILL.md"
  "$repo_root/.agents/skills/provider-resource-dev/SKILL.md"
  "$repo_root/.claude/skills/provider-resource-dev/references/zhenyuan-verification.md"
  "$repo_root/.agents/skills/provider-resource-dev/references/zhenyuan-verification.md"
)
for file in "${g_urgent_route_files[@]}"; do
  for phrase in \
    'G / 紧急普通 D 的双 owner 契约' \
    '源客户主单 assignee 保持新山（521957）' \
    '528766 研发关联单 assignee 固定过载（484483）' \
    'healthy existing claim 不抢占'; do
    grep -Fq "$phrase" "$file"
  done
done

for file in \
  "$repo_root/loops/aone-triage.md" \
  "$repo_root/loops/persona-collab.md" \
  "$pd_agent" "$pd_agent_mirror" \
  "$rd_agent" "$rd_agent_mirror" \
  "$qa_agent" "$qa_agent_mirror" \
  "$repo_root/CLAUDE.md" "$repo_root/AGENTS.md"; do
  grep -Fq 'G / 紧急普通 D 的双 owner 契约' "$file"
done

for phrase in \
  '源工单禁止模型直接 claim/wrap/release/评论' \
  '源工单禁令不约束按既有契约由内部链承接的 528766' \
  '非紧急 D 与 I 的 Provider docs 紧急兜底腿' \
  '528766 由 RD finalizer claim/bookend' \
  '源工单仍由 executor bookend'; do
  grep -Fq "$phrase" "$persona_collab"
done
if grep -Fq '控制面 Task 的模型 run 不执行 `claim.sh`、`wrap.sh`、`release` 或直接评论' \
     "$persona_collab" ||
   grep -Fq 'G/紧急普通 D 的 528766 是唯一例外' "$persona_collab"; then
  echo 'zhenyuan_self_close_rules_test: blanket model-run bookend ban remains' >&2
  exit 1
fi

for file in "$main_skill" "$routing" "$repo_root/CLAUDE.md"; do
  grep -Fq '源工单禁令不约束按既有契约由内部链承接的 528766' "$file"
  grep -Fq '非紧急 D 与 I 的 Provider docs 紧急兜底腿' "$file"
  if grep -Fq 'G/紧急普通 D 的 528766 是唯一例外' "$file"; then
    echo "zhenyuan_self_close_rules_test: existing internal 528766 path regressed in $file" >&2
    exit 1
  fi
done

for file in \
  "$routing" \
  "$main_skill" \
  "$verification" \
  "$repo_root/.claude/skills/provider-resource-dev/SKILL.md"; do
  if grep -Fq 'D-新山' "$file"; then
    echo "zhenyuan_self_close_rules_test: stale G/urgent-D handoff remains in $file" >&2
    exit 1
  fi
done

echo "zhenyuan_self_close_rules_test: PASS"
