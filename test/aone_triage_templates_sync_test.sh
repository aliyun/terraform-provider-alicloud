#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

diff -u \
  "$repo_root/.agents/skills/aone-triage/references/templates.md" \
  "$repo_root/.claude/skills/aone-triage/references/templates.md"

echo "aone_triage_templates_sync_test: PASS"
