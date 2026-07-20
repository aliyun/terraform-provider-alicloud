#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
skill_rel="terraform-provider-live-repro"
codex_skill="$repo_root/.agents/skills/$skill_rel"
claude_skill="$repo_root/.claude/skills/$skill_rel"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$repo_root/bootstrap/skills-mirror-lib.sh"

while IFS= read -r source_file; do
  rel="${source_file#$codex_skill/}"
  expected="$tmpdir/${rel//\//__}"
  mirror_sed_codex_to_claude < "$source_file" > "$expected"
  diff -u "$expected" "$claude_skill/$rel"
done < <(find "$codex_skill" -type f | sort)

for term in \
  "terraform plan -input=false -no-color -out=create.tfplan" \
  "TF_LOG_PROVIDER=DEBUG" \
  "DescribeFileSystems" \
  "replace_paths" \
  "Never apply a drift/replacement plan" \
  "html-report-preview" \
  "destroy-after-evidence"; do
  grep -q -- "$term" "$codex_skill/SKILL.md"
done

raw_log="$tmpdir/provider.log"
cat > "$raw_log" <<'EOF'
2026-07-20T11:45:51.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: *************** CreateFileSystem Response ***************
2026-07-20T11:45:51.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: {"FileSystemId":"fs-123","RequestId":"req-create"}Domain:, Version:
2026-07-20T11:45:51.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: ***************  Request ***************
2026-07-20T11:45:51.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: "{\"AccessKeyId\":\"SHOULD-NOT-LEAK\",\"FileSystemType\":\"standard\",\"VpcId\":\"vpc-123\",\"VSwitchId\":\"vsw-123\",\"ZoneId\":\"cn-test-a\"}"
2026-07-20T11:45:56.330+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: *************** DescribeFileSystems Response ***************
2026-07-20T11:45:56.330+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: {"FileSystems":{"FileSystem":[{"FileSystemId":"fs-123","Status":"Running","VpcId":""}]},"RequestId":"req-read"}Domain:, Version:
EOF

timeline="$tmpdir/timeline.md"
python3 "$codex_skill/scripts/extract-api-timeline.py" "$raw_log" --format markdown > "$timeline"
grep -q 'req-create' "$timeline"
grep -q 'VpcId.*vpc-123' "$timeline"
grep -q 'QuorumVswId=missing' "$timeline"
if grep -q 'SHOULD-NOT-LEAK' "$timeline"; then
  echo "terraform_provider_live_repro_skill_test: parser leaked a forbidden request field" >&2
  exit 1
fi

python3 "$codex_skill/scripts/extract-api-timeline.py" "$raw_log" --format jsonl \
  | jq -e 'select(.request_id == "req-read" and (.observations | index("QuorumVswId=missing")))' >/dev/null

report_md="$tmpdir/report.md"
report_html="$tmpdir/report.html"
printf '# Live Repro\n\n| API | RequestId |\n|---|---|\n| CreateFileSystem | req-create |\n' > "$report_md"
python3 "$codex_skill/scripts/render-report-html.py" "$report_md" "$report_html" >/dev/null
grep -q '<title>Live Repro</title>' "$report_html"
grep -q '<table>' "$report_html"
if grep -qi 'data:image' "$report_html"; then
  echo "terraform_provider_live_repro_skill_test: renderer emitted base64 image data" >&2
  exit 1
fi

unsafe_md="$tmpdir/unsafe.md"
printf '# Unsafe\n\n<script>alert(1)</script>\n' > "$unsafe_md"
if python3 "$codex_skill/scripts/render-report-html.py" "$unsafe_md" "$tmpdir/unsafe.html" >/dev/null 2>&1; then
  echo "terraform_provider_live_repro_skill_test: renderer accepted executable HTML" >&2
  exit 1
fi

echo "terraform_provider_live_repro_skill_test: PASS"
