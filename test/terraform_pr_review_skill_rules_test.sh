#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

for skill in \
  "$repo_root/.agents/skills/terraform-pr-review/SKILL.md" \
  "$repo_root/.claude/skills/terraform-pr-review/SKILL.md"; do
  for term in \
    "Pull Request Max Commits" \
    "git rev-list --count" \
    "保持单提交" \
    "gh run view <run_id> --job <job_id> --log" \
    "同一个 workflow 里的其它 job" \
    "内部 Aone" \
    "tf_provider" \
    "528766" \
    "禁止同步到 tf_customer"; do
    if ! grep -q "$term" "$skill"; then
      echo "terraform_pr_review_skill_rules_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done
done

echo "terraform_pr_review_skill_rules_test: PASS"
