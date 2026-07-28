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
  ".claude/skills/provider-resource-dev/SKILL.md"
  ".agents/skills/provider-resource-dev/SKILL.md"
  ".claude/skills/provider-resource-dev/references/zhenyuan-verification.md"
  ".agents/skills/provider-resource-dev/references/zhenyuan-verification.md"
  ".claude/skills/invoke-terraform-acc-test-remote/SKILL.md"
  ".agents/skills/invoke-terraform-acc-test-remote/SKILL.md"
  ".claude/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
  ".agents/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
  "loops/aone-triage.md"
  "loops/persona-collab.md"
  ".claude/agents/terraform-pd.md"
  ".codex/agents/terraform-pd.toml"
  ".claude/agents/terraform-rd.md"
  ".codex/agents/terraform-rd.toml"
  "CLAUDE.md"
  "AGENTS.md"
  "autonomy.md"
)

for rel in "${active_rules[@]}"; do
  test -f "$repo_root/$rel"
done

for forbidden in \
  '2165097' \
  '2169561' \
  'upstream.cloudspec_gap' \
  'upstream.cloudspec_docs_quality' \
  'WORKER_1783326253279' \
  '镇元 agent' \
  '镇元agent' \
  '谜拟' \
  '念依'; do
  if grep -Fn "$forbidden" "${active_rules[@]/#/$repo_root/}"; then
    echo "zhenyuan_self_close_rules_test: active rules still route through '$forbidden'" >&2
    exit 1
  fi
done

jq -e '
  (.upstream.cloudspec_gap? == null)
  and (.upstream.cloudspec_docs_quality? == null)
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

for phrase in \
  'CloudSpec 原主单自闭环' \
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
  'missing_capability' \
  'blocked' \
  'release/idle' \
  '不得 finish'; do
  grep -Fq "$phrase" "$routing"
done

for phrase in \
  '在 CloudSpec OK 判定前短路' \
  '不得被 schema、properties 或 CoverageScore 全绿覆盖' \
  '文档源正确性' \
  'CloudSpec 文档源正确，Provider 本地文档生成/展示偏差' \
  '只有 CloudSpec 文档源正确'; do
  grep -Fq "$phrase" "$routing"
  grep -Fq "$phrase" "$verification"
done

for phrase in \
  '镇元 OK 四条件' \
  '文档源正确性' \
  'CloudSpec 文档源正确、Provider 本地文档生成/展示偏差' \
  '直接进入分支 E'; do
  grep -Fq "$phrase" "$team_roster"
done

if grep -Fq 'OK 三条件' "$routing" "$verification" "$team_roster"; then
  echo 'zhenyuan_self_close_rules_test: stale three-condition CloudSpec OK contract' >&2
  exit 1
fi

awk '
  index($0, "CloudSpec 文档源错误") && !source_error { source_error = NR }
  index($0, "CloudSpec OK（四条件全满足）") && !ok_check { ok_check = NR }
  END { exit !(source_error && ok_check && source_error < ok_check) }
' "$routing"

for phrase in \
  'CloudSpec 文档源错误' \
  '分支 E' \
  'Provider 本地文档生成/展示偏差' \
  '分支 D'; do
  grep -Fq "$phrase" "$templates"
  grep -Fq "$phrase" "$pd_agent"
done

grep -Eq 'CloudSpec 文档源错误.*分支 E' "$templates"
grep -Eq 'CloudSpec 文档源正确，Provider 本地文档生成/展示偏差.*分支 D' "$templates"
if grep -Eq 'CloudSpec 文档源错误.*分支 D|Provider 本地文档生成/展示偏差.*分支 E' "$templates"; then
  echo 'zhenyuan_self_close_rules_test: document source/local rendering routes are reversed' >&2
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
  grep -Fq 'role: terraform-qa | terraform-rd-finalizer' "$file"
  grep -Fq 'action: acc_verify | finalize' "$file"
  if grep -Fq 'clarify' "$file"; then
    echo "zhenyuan_self_close_rules_test: stale RD clarify route in $file" >&2
    exit 1
  fi
done

for file in "$pd_agent" "$pd_agent_mirror"; do
  grep -Fq 'action: dev | acc_verify | finalize' "$file"
  if grep -Fq 'clarify' "$file"; then
    echo "zhenyuan_self_close_rules_test: stale PD clarify action in $file" >&2
    exit 1
  fi
done

echo "zhenyuan_self_close_rules_test: PASS"
