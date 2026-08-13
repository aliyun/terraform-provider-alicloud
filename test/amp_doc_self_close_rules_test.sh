#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

active_routes=(
  AGENTS.md
  CLAUDE.md
  autonomy.md
  loops/aone-triage.md
  loops/persona-collab.md
  bridge/aone_tasks.py
  .claude/agents/terraform-pd.md
  .claude/agents/terraform-rd.md
  .codex/agents/terraform-pd.toml
  .codex/agents/terraform-rd.toml
  .claude/skills/aone-triage/SKILL.md
  .agents/skills/aone-triage/SKILL.md
  .claude/skills/aone-triage/references/tf-customer-request-routing.md
  .agents/skills/aone-triage/references/tf-customer-request-routing.md
)

if rg -n 'cloudspec_docs_quality|2169561|念依' "${active_routes[@]}"; then
  echo "amp_doc_self_close_rules_test: active I route still hands off to legacy docs pool" >&2
  exit 1
fi

stale_i_provider_fallback='I 的 Provider (public )?docs 紧急兜底|I public-docs 紧急腿|I/H.{0,80}(合法|claim|bookend)|claim.{0,80}I/H|Provider markdown.{0,120}(只能作为|紧急兜底)'
if rg -n -U "$stale_i_provider_fallback" "${active_routes[@]}"; then
  echo "amp_doc_self_close_rules_test: active I route still authorizes a Provider docs 528766 fallback" >&2
  exit 1
fi

jq -e '(.upstream.cloudspec_docs_quality? == null)' config/pools.json >/dev/null

for path in \
  .agents/skills/amp-doc-backend/SKILL.md \
  .claude/skills/amp-doc-backend/SKILL.md; do
  test -f "$path"
  grep -Fq '/usr/bin/python3 -I <jarvis-root>/bootstrap/amp_safe.py' "$path"
  grep -Fq 'doc recommend-resource' "$path"
  grep -Fq 'doc get' "$path"
  grep -Fq 'doc submit-audit' "$path"
  grep -Fq '已进入审核，尚未正式发布' "$path"
  grep -Fq '发布状态未验证' "$path"
  grep -Fq '只允许更新 `description`、`title`、`example`、`enumValueTitles`' "$path"
  grep -Eq '枚举(值集合|集合).*(分支 E|转.*E)' "$path"
  if rg -n '权威新枚举|enum（仅|`enum`.*(白名单|允许|可编辑)' "$path"; then
    echo "amp_doc_self_close_rules_test: enum values remain editable in $path" >&2
    exit 1
  fi
  if rg -n '(^|[[:space:]`])amp[[:space:]]' "$path"; then
    echo "amp_doc_self_close_rules_test: raw amp invocation in $path" >&2
    exit 1
  fi
done

test ! -e .agents/skills/amp-doc-backend/.amp/context.yaml
test ! -e .claude/skills/amp-doc-backend/.amp/context.yaml
cmp -s .agents/skills/amp-doc-backend/SKILL.md \
  .claude/skills/amp-doc-backend/SKILL.md
cmp -s .agents/skills/amp-doc-backend/references/metadata.md \
  .claude/skills/amp-doc-backend/references/metadata.md

for path in \
  .agents/skills/amp-doc-backend/references/metadata.md \
  .claude/skills/amp-doc-backend/references/metadata.md; do
  grep -Eq '枚举(值集合|集合).*(分支 E|转.*E)' "$path"
  if rg -n '权威新枚举|enum（仅|`enum`.*(白名单|允许|可编辑)' "$path"; then
    echo "amp_doc_self_close_rules_test: enum values remain editable in $path" >&2
    exit 1
  fi
done

grep -Fq 'recommend-resource --env online' AGENTS.md
grep -Fq 'API/struct 固定走' AGENTS.md
grep -Fq '绝不宣称正式发布' AGENTS.md
grep -Eq '源单.*(自闭环|直接执行)' AGENTS.md
grep -Fq '分支 I 只指 resource/API/struct' bridge/aone_tasks.py
grep -Fq 'amp-doc-backend 自闭环' bridge/aone_tasks.py
grep -Fq 'recommend-resource --env online' .claude/agents/terraform-rd.md
grep -Fq 'submit-audit→get-audit-url' .claude/agents/terraform-rd.md

echo "amp_doc_self_close_rules_test: PASS"
