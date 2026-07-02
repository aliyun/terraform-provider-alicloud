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
    "cloudspec_gap"; do
    if ! grep -q "$term" "$skill"; then
      echo "provider_internal_aone_rules_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done
done

echo "provider_internal_aone_rules_test: PASS"
