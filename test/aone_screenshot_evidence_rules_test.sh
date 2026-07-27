#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

for skill in \
  "$repo_root/.claude/skills/aone-triage/SKILL.md" \
  "$repo_root/.agents/skills/aone-triage/SKILL.md"; do
  for term in \
    "OpenAPI、CloudSpec/ACube 映射、Provider 源码三层查证必须各有截图证据" \
    "visual_evidence_manifest" \
    "不得传 \`--comment\`" \
    "禁止静默省略"; do
    grep -Fq "$term" "$skill" || {
      echo "aone_screenshot_evidence_rules_test: missing '$term' in $skill" >&2
      exit 1
    }
  done
  if grep -Fq "非 Terraform 流程可调" "$skill"; then
    echo "aone_screenshot_evidence_rules_test: stale Terraform exclusion in $skill" >&2
    exit 1
  fi
done

for skill in \
  "$repo_root/.claude/skills/screenshot-evidence/SKILL.md" \
  "$repo_root/.agents/skills/screenshot-evidence/SKILL.md"; do
  for term in \
    "CloudSpec/ACube 映射" \
    "evidence-manifest.md" \
    "visual_evidence_manifest" \
    "validate-manifest.py" \
    "AONE_RESULT.reply_body" \
    "Terraform finalizer：只取预览 URL" \
    "禁止产生第二条 Aone 评论"; do
    grep -Fq "$term" "$skill" || {
      echo "aone_screenshot_evidence_rules_test: missing '$term' in $skill" >&2
      exit 1
    }
  done
done

for file in \
  "$repo_root/.claude/agents/terraform-pd.md" \
  "$repo_root/.codex/agents/terraform-pd.toml" \
  "$repo_root/.claude/agents/terraform-rd.md" \
  "$repo_root/.codex/agents/terraform-rd.toml" \
  "$repo_root/loops/persona-collab.md"; do
  grep -Fq "visual_evidence_manifest" "$file" || {
    echo "aone_screenshot_evidence_rules_test: missing visual_evidence_manifest in $file" >&2
    exit 1
  }
done

for term in \
  "Terraform 三层可视化证据契约" \
  "visual_evidence_manifest" \
  "AONE_RESULT.reply_body" \
  '严禁传 `--comment`'; do
  grep -Fq "$term" "$repo_root/bridge/aone_tasks.py" || {
    echo "aone_screenshot_evidence_rules_test: bridge prompt missing '$term'" >&2
    exit 1
  }
done

validator="$repo_root/.claude/skills/screenshot-evidence/scripts/validate-manifest.py"
mkdir -p "$tmpdir/screenshots"
: > "$tmpdir/screenshots/openapi.png"
: > "$tmpdir/screenshots/provider.png"
valid_manifest="$tmpdir/valid-manifest.md"
printf '%s\n' \
  '| layer | result | screenshot | source | note |' \
  '|---|---|---|---|---|' \
  "| OpenAPI | pass | $tmpdir/screenshots/openapi.png | https://example.test/openapi | field exists |" \
  '| CloudSpec/ACube | n-a | N/A | command: metadata lookup | missing_capability: page unavailable |' \
  "| Provider | fail | $tmpdir/screenshots/provider.png | provider.go:42 | field missing |" \
  > "$valid_manifest"
python3 "$validator" "$valid_manifest" >/dev/null

invalid_manifest="$tmpdir/invalid-manifest.md"
printf '%s\n' \
  '| layer | result | screenshot | source | note |' \
  '|---|---|---|---|---|' \
  "| OpenAPI | pass | $tmpdir/screenshots/missing.png | https://example.test/openapi | field exists |" \
  '| CloudSpec/ACube | n-a | N/A | command: metadata lookup | missing_capability: page unavailable |' \
  > "$invalid_manifest"
if python3 "$validator" "$invalid_manifest" >/dev/null 2>&1; then
  echo "aone_screenshot_evidence_rules_test: validator accepted missing screenshot/layer" >&2
  exit 1
fi

grep -Fq "RD finalizer 可统一上传一次但不得传 \`--comment\`" \
  "$repo_root/loops/aone-triage.md"
if grep -Fq "仅非 Terraform 流程可上传" "$repo_root/loops/aone-triage.md"; then
  echo "aone_screenshot_evidence_rules_test: stale upload ban in loops/aone-triage.md" >&2
  exit 1
fi

bash "$repo_root/bootstrap/mirror.sh" check \
  "$repo_root/.claude/skills/aone-triage/SKILL.md" \
  "$repo_root/.claude/skills/screenshot-evidence/SKILL.md" \
  "$repo_root/.claude/skills/screenshot-evidence/scripts/validate-manifest.py" \
  "$repo_root/.claude/agents/terraform-pd.md" \
  "$repo_root/.claude/agents/terraform-rd.md"

echo "aone_screenshot_evidence_rules_test: PASS"
