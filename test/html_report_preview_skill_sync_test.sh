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
  "agents/openai.yaml"; do
  expected="$tmpdir/${rel//\//__}"
  mirror_sed_codex_to_claude < "$repo_root/.agents/skills/html-report-preview/$rel" > "$expected"
  diff -u \
    "$expected" \
    "$repo_root/.claude/skills/html-report-preview/$rel"
done

for skill in \
  "$repo_root/.agents/skills/html-report-preview/SKILL.md" \
  "$repo_root/.claude/skills/html-report-preview/SKILL.md"; do
  for term in \
    "bootstrap/html-report-preview.sh" \
    "JARVIS_HTML_REPORT_TOKEN" \
    "JARVIS_HTML_REPORT_BASE_URL" \
    "/api/reports/aone" \
    "/reports/aone" \
    "from-aone" \
    "--comment" \
    "Terraform single-writer boundary" \
    "AONE_RESULT.reply_body" \
    "without \`--comment\`" \
    "must not call \`wrap.sh\`" \
    "buc_required" \
    "Never make report screenshots public-read" \
    "--acl private" \
    "oss sign" \
    "--timeout 15768000" \
    "GET-signed URL" \
    "rgv587_flag:sm" \
    "inline script" \
    "禁止规避" \
    "platform_blocked" \
    "never echoed" \
    "success:false,status:failed"; do
    if ! grep -q -- "$term" "$skill"; then
      echo "html_report_preview_skill_sync_test: missing '$term' in $skill" >&2
      exit 1
    fi
  done

  if grep -q -- "JARVIS_HTML_REPORT_TOKEN=" "$skill"; then
    echo "html_report_preview_skill_sync_test: skill must not contain plaintext token assignment: $skill" >&2
    exit 1
  fi

  for forbidden in \
    "--acl public-read" \
    "as a public-read object"; do
    if grep -q -- "$forbidden" "$skill"; then
      echo "html_report_preview_skill_sync_test: forbidden '$forbidden' in $skill" >&2
      exit 1
    fi
  done
done

echo "html_report_preview_skill_sync_test: PASS"
