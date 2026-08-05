#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

require_text() {
  local file="$1"
  local text="$2"
  if ! grep -Fq -- "$text" "$file"; then
    echo "provider_go_validation_policy_test: missing '$text' in $file" >&2
    exit 1
  fi
}

for file in \
  "$repo_root/.claude/agents/terraform-rd.md" \
  "$repo_root/.codex/agents/terraform-rd.toml" \
  "$repo_root/.claude/skills/terraform-provider-release/SKILL.md" \
  "$repo_root/.agents/skills/terraform-provider-release/SKILL.md" \
  "$repo_root/.claude/skills/provider-resource-dev/SKILL.md" \
  "$repo_root/.agents/skills/provider-resource-dev/SKILL.md"; do
  require_text "$file" "go vet -p=1 ./alicloud"
  require_text "$file" "docs-only"
done

for file in \
  "$repo_root/.claude/agents/terraform-rd.md" \
  "$repo_root/.codex/agents/terraform-rd.toml" \
  "$repo_root/.claude/skills/terraform-provider-release/SKILL.md" \
  "$repo_root/.agents/skills/terraform-provider-release/SKILL.md" \
  "$repo_root/.claude/skills/provider-resource-dev/SKILL.md" \
  "$repo_root/.agents/skills/provider-resource-dev/SKILL.md"; do
  require_text "$file" "每个连贯"
  require_text "$file" "首次 push"
done

if rg -n 'ops\.build|Local compilation is strictly forbidden|Never compile locally' \
  "$repo_root/.claude/skills/terraform-provider-release/SKILL.md" \
  "$repo_root/.agents/skills/terraform-provider-release/SKILL.md"; then
  echo "provider_go_validation_policy_test: found obsolete validation wording" >&2
  exit 1
fi

jq -e '
  .workspaces.terraform_provider.validation.strategy == "staged" and
  .workspaces.terraform_provider.ops.vet == "go vet -p=1 ./alicloud" and
  .workspaces.terraform_provider.validation.edit_batch.docs_only == [] and
  .workspaces.terraform_provider.validation.local_forbidden ==
    ["go vet ./...","go build ./..."] and
  .workspaces.automation_platform_runtime.ops ==
    {"build":"go build ./...","test":"go test ./..."} and
  (.workspaces.automation_platform_runtime | has("validation") | not) and
  (.workspaces.terraformer | has("validation") | not)
' "$repo_root/config/workspaces.json" >/dev/null

echo "provider_go_validation_policy_test: PASS"
