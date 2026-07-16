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
  "references/alicloud-resource-development.md"; do
  claude_file="$repo_root/.claude/skills/terraformer-resource-dev/$rel"
  codex_file="$repo_root/.agents/skills/terraformer-resource-dev/$rel"
  test -f "$claude_file"
  test -f "$codex_file"

  expected="$tmpdir/${rel//\//__}"
  mirror_sed_codex_to_claude < "$codex_file" > "$expected"
  diff -u "$expected" "$claude_file"
done

expected_layout=$'SKILL.md\nreferences/alicloud-resource-development.md'
for skill_root in \
  "$repo_root/.claude/skills/terraformer-resource-dev" \
  "$repo_root/.agents/skills/terraformer-resource-dev"; do
  actual_layout="$(find "$skill_root" -type f | sed "s#^$skill_root/##" | LC_ALL=C sort)"
  if [[ "$actual_layout" != "$expected_layout" ]]; then
    echo "terraformer_resource_dev_skill_rules_test: unexpected layout in $skill_root" >&2
    diff -u <(printf '%s\n' "$expected_layout") <(printf '%s\n' "$actual_layout") >&2 || true
    exit 1
  fi
done

for skill in \
  "$repo_root/.claude/skills/terraformer-resource-dev/SKILL.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/SKILL.md"; do
  for term in \
    "description: Use when developing, diagnosing, or fixing an Alibaba Cloud resource in Terraformer" \
    "bootstrap/workspace.sh dir terraformer" \
    "stop and escalate missing_capability" \
    "aone-triage" \
    "loops/adhoc-intake.md" \
    "bootstrap/claim.sh claim" \
    "bootstrap/wrap.sh sync <id>" \
    "bootstrap/wrap.sh done" \
    "bootstrap/claim.sh release" \
    "references/alicloud-resource-development.md" \
    "terraform-rd" \
    "terraform-qa" \
    "InitResources" \
    "Do not produce or infer resource relationships"; do
    if ! grep -Fq -- "$term" "$skill"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done
done

for reference in \
  "$repo_root/.claude/skills/terraformer-resource-dev/references/alicloud-resource-development.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/references/alicloud-resource-development.md"; do
  for term in \
    "A. Direct full List with a single-field Import ID" \
    "B. One List returns every multipart-ID segment" \
    "C. Parent-child traversal" \
    "D. Complete enumeration is unavailable" \
    'd.SetId(...)' \
    'ParseResourceId(...)' \
    "A multipart Import ID does not by itself require parent traversal" \
    "A Data Source may require the parent ID" \
    "Reset pagination for every parent" \
    "For token pagination, stop when the returned next token is empty regardless of page length" \
    "For page-number pagination" \
    "Do not produce or infer connections" \
    "does not block core discovery and Import ID support" \
    "go test ./providers/alicloud" \
    "go test ./..." \
    "/tmp/terraformer"; do
    if ! grep -Fq -- "$term" "$reference"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $reference" >&2
      exit 1
    fi
  done
done

evaluation_report="$repo_root/docs/superpowers/reports/2026-07-16-terraformer-resource-dev-forward-evaluation.md"
test -f "$evaluation_report"
for term in \
  "Scenario A — PASS" \
  "Scenario B — PASS" \
  "Scenario C — PASS"; do
  if ! grep -Fq -- "$term" "$evaluation_report"; then
    echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $evaluation_report" >&2
    exit 1
  fi
done

echo "terraformer_resource_dev_skill_rules_test: PASS"
