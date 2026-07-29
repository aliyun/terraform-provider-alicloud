#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

files=(
  "$repo_root/.agents/skills/terraform-provider-release/SKILL.md"
  "$repo_root/.claude/skills/terraform-provider-release/SKILL.md"
  "$repo_root/.agents/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
  "$repo_root/.claude/skills/terraform-provider-release/references/cloudspec-pre-resource-loop.md"
  "$repo_root/.agents/skills/invoke-terraform-acc-test-remote/SKILL.md"
  "$repo_root/.claude/skills/invoke-terraform-acc-test-remote/SKILL.md"
)

for file in "${files[@]}"; do
  for term in \
    'verification_mode: cloudspec_pre' \
    '源单' \
    'Provider' \
    '远程 ACC' \
    'createBuildTaskV2' \
    '528766'; do
    grep -Fq "$term" "$file" || {
      echo "terraform_provider_release_cloudspec_test: missing '$term' in $file" >&2
      exit 1
    }
  done
  if grep -Fq 'E → D-临钧' "$file"; then
    echo "terraform_provider_release_cloudspec_test: stale E handoff in $file" >&2
    exit 1
  fi
  if grep -Fq 'hand off through Acube' "$file"; then
    echo "terraform_provider_release_cloudspec_test: executable Acube handoff remains in $file" >&2
    exit 1
  fi
done

for file in \
  "$repo_root/.agents/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md" \
  "$repo_root/.claude/skills/aone-triage/references/acube-createBuildTaskV2-workflow.md"; do
  grep -Fq '历史只读' "$file"
  grep -Fq 'D/E/G' "$file"
  grep -Fq '历史 relation/PR 防重取证' "$file"
  grep -Fq 'I/H' "$file"
  for forbidden in \
    'curl ' \
    'POST /' \
    '```bash' \
    'relation add ' \
    'workitem update ' \
    '--assignee ' \
    'createBuildTaskV2`：'; do
    if grep -Fq -- "$forbidden" "$file"; then
      echo "terraform_provider_release_cloudspec_test: executable '$forbidden' remains in $file" >&2
      exit 1
    fi
  done
done

post_pr_context="$repo_root/bootstrap/post-pr-context.sh"
# shellcheck source=/dev/null
source "$post_pr_context"
marker="$(mktemp)"
trap 'rm -f "$marker"' EXIT
printf '%s\n' \
  '{"pid":321,"pgid":654,"policy_revision":"terraform-rd-single-writer-v4","kind":"pr_ci_fix"}' \
  >"$marker"
_jarvis_post_pr_marker_valid "$marker" 321 654 || {
  echo "terraform_provider_release_cloudspec_test: v4 marker compatibility missing" >&2
  exit 1
}
exec_lineage_source="$(
  sed -n '/^jarvis_post_pr_exec_lineage_active()/,/^}/p' "$post_pr_context"
)"
v4_exec_count="$(printf '%s\n' "$exec_lineage_source" |
  grep -Fc -- '--policy-revision terraform-rd-single-writer-v4' || true)"
if [ "$v4_exec_count" -lt 2 ]; then
  echo "terraform_provider_release_cloudspec_test: v4 exec lineage compatibility missing" >&2
  exit 1
fi
grep -Fq -- '--policy-revision terraform-rd-single-writer-v5' \
  <<<"$exec_lineage_source"

echo "terraform_provider_release_cloudspec_test: PASS"
