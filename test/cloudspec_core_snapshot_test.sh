#!/usr/bin/env bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
lock="$repo_root/config/cloudspec-core.lock.json"

required_skills=(
  cloudspec-amp-workflow
  cloudspec-idl-guide
  cloudspec-resource-edit
  cloudspec-operation-edit
  cloudspec-flag-mode-edit
  cloudspec-build-fix
  cloudspec-norm-check-fix
  cloudspec-shared-knowledge
)

jq -e '
  .plugin == "cloudspec-core" and
  .version == "1.3.0" and
  .source.commit == "3b64ce4ad81b3d28eb563997eaa8e29241ba8b77" and
  .source.skillsTree == "f35fdf99abbd7d80136cee8a0d43760f806c8a06" and
  .distribution == "vendored-skill-snapshot" and
  .hooksIncluded == false
' "$lock" >/dev/null

for skill in "${required_skills[@]}"; do
  test -f "$repo_root/.claude/skills/$skill/SKILL.md"
  test -f "$repo_root/.agents/skills/$skill/SKILL.md"
  jq -e --arg skill "$skill" '.skills | index($skill) != null' "$lock" >/dev/null
done

# 机器人只拿业务技能，不把正式插件默认 telemetry hook 带进仓库。
test ! -e "$repo_root/.claude/plugins/cloudspec-core/hooks"
test ! -e "$repo_root/.agents/plugins/cloudspec-core/hooks"
grep -q $'^amp\t' "$repo_root/bootstrap/deps.lock"

bash "$repo_root/bootstrap/cloudspec-core.sh" check >/dev/null
bash "$repo_root/bootstrap/mirror.sh" check >/dev/null

echo "cloudspec_core_snapshot_test: PASS"
