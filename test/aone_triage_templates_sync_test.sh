#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

for rel in \
  "SKILL.md" \
  "references/templates.md" \
  "references/delivery-aliyun-automation-agent.md" \
  "references/delivery-cloudspec.md" \
  "scripts/sync-provider.sh"; do
  diff -u \
    "$repo_root/.agents/skills/aone-triage/$rel" \
    "$repo_root/.claude/skills/aone-triage/$rel"
done

echo "aone_triage_templates_sync_test: PASS"
