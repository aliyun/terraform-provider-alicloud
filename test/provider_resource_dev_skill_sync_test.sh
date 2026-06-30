#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

diff -u \
  "$repo_root/.agents/skills/provider-resource-dev/SKILL.md" \
  "$repo_root/.claude/skills/provider-resource-dev/SKILL.md"

echo "provider_resource_dev_skill_sync_test: PASS"
