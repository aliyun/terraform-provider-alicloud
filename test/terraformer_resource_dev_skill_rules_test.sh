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

for unexpected in \
  "$repo_root/.claude/skills/terraformer-resource-dev/agents/openai.yaml" \
  "$repo_root/.agents/skills/terraformer-resource-dev/agents/openai.yaml"; do
  if [[ -e "$unexpected" ]]; then
    echo "terraformer_resource_dev_skill_rules_test: unexpected optional metadata $unexpected" >&2
    exit 1
  fi
done

for skill in \
  "$repo_root/.claude/skills/terraformer-resource-dev/SKILL.md" \
  "$repo_root/.agents/skills/terraformer-resource-dev/SKILL.md"; do
  for term in \
    "description: Use when developing, diagnosing, or fixing an Alibaba Cloud resource in Terraformer" \
    "bootstrap/workspace.sh dir terraformer" \
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
    "A. Direct full List" \
    "B. One List returns every composite-ID segment" \
    "C. Parent-child traversal" \
    "D. Complete enumeration is unavailable" \
    'd.SetId(...)' \
    'ParseResourceId(...)' \
    "A multipart Import ID does not by itself require parent traversal" \
    "A Data Source may require the parent ID" \
    "Reset pagination for every parent" \
    "Do not produce or infer connections" \
    "go test ./providers/alicloud" \
    "go test ./..." \
    "/tmp/terraformer"; do
    if ! grep -Fq -- "$term" "$reference"; then
      echo "terraformer_resource_dev_skill_rules_test: missing '$term' in $reference" >&2
      exit 1
    fi
  done
done

echo "terraformer_resource_dev_skill_rules_test: PASS"
