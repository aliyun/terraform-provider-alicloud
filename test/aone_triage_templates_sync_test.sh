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
  "references/delivery-aliyun-automation-platform.md" \
  "references/delivery-cloudspec.md" \
  "scripts/sync-provider.sh"; do
  expected="$tmpdir/${rel//\//__}"
  mirror_sed_claude_to_codex < "$repo_root/.claude/skills/aone-triage/$rel" > "$expected"
  diff -u \
    "$repo_root/.agents/skills/aone-triage/$rel" \
    "$expected"
done

skill="$repo_root/.claude/skills/aone-triage/SKILL.md"
delivery="$repo_root/.claude/skills/aone-triage/references/delivery-aliyun-automation-platform.md"
adhoc="$repo_root/loops/adhoc-intake.md"

grep -Fq 'references/delivery-aliyun-automation-platform.md' "$skill"
grep -Fq '1091779' "$skill"
grep -Fq 'automation_platform' "$skill"

for token in \
  '1091779' \
  'automation_platform' \
  '172823' \
  '2156624' \
  'WORKER_1782379562571' \
  'prestage 66' \
  'prod 67' \
  'cloudspec build' \
  'systemlog-prod' \
  'systemlog-pre' \
  'aliyun-automation-platform-1252907582134651-pop-aliyun-cn' \
  'sls-log-query-aliyun-automation-platform' \
  'automation_platform_frontend' \
  'automation_platform_runtime' \
  'automation_platform_function_test' \
  'automation_platform_api' \
  'automation_platform_api_inner'; do
  grep -Fq "$token" "$delivery"
done

for excluded in 'Agent portal' 'AgentRuntime' 'PlayGround' 'FC sandbox' 'WebSocket/STS' 'aliyun-automation-agent'; do
  grep -Fq "$excluded" "$delivery"
done

grep -Fq '直接 URL/ID 处理不受 assignedTo 限制' "$delivery"
grep -Fq '正式发布必须取得明确的人工批准' "$delivery"
grep -Fq 'tf_customer / tf_provider / mcp_server / automation_platform / api_toolkit' "$adhoc"
grep -Fq 'bootstrap/workspace.sh dir' "$adhoc"
if grep -Fq 'tf_provider / tf_customer / mcp_server / cloudspec / api_toolkit' "$adhoc"; then
  echo 'adhoc-intake still exposes obsolete cloudspec pool candidate' >&2
  exit 1
fi

echo "aone_triage_templates_sync_test: PASS"
