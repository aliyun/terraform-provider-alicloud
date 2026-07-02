#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$repo_root/bootstrap/skills-mirror-lib.sh"

for rel in \
  "SKILL.md" \
  "references/templates.md" \
  "references/delivery-aliyun-automation-agent.md" \
  "references/delivery-cloudspec.md" \
  "scripts/sync-provider.sh"; do
  expected="$tmpdir/${rel//\//__}"
  mirror_sed_claude_to_codex < "$repo_root/.claude/skills/aone-triage/$rel" > "$expected"
  diff -u \
    "$repo_root/.agents/skills/aone-triage/$rel" \
    "$expected"
done

echo "aone_triage_templates_sync_test: PASS"
