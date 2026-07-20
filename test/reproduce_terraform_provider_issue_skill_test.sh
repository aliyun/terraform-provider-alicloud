#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
skill_rel="reproduce-terraform-provider-issue"
agents_skill="$repo_root/.agents/skills/$skill_rel"
claude_skill="$repo_root/.claude/skills/$skill_rel"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$repo_root/bootstrap/skills-mirror-lib.sh"

while IFS= read -r source_file; do
  rel="${source_file#$agents_skill/}"
  expected="$tmpdir/${rel//\//__}"
  mirror_sed_codex_to_claude < "$source_file" > "$expected"
  diff -u "$expected" "$claude_skill/$rel"
done < <(find "$agents_skill" -type f | sort)

for runtime_skill in "$agents_skill" "$claude_skill"; do
  test -f "$runtime_skill/SKILL.md"
  test -x "$runtime_skill/scripts/extract-api-timeline.py"
  test -x "$runtime_skill/scripts/render-report-html.py"
  grep -q '^name: reproduce-terraform-provider-issue$' "$runtime_skill/SKILL.md"
  grep -q '^description: ' "$runtime_skill/SKILL.md"
  grep -q '^# Terraform Provider 用户问题现场排查复现$' "$runtime_skill/SKILL.md"
  grep -q '禁止执行漂移或替换计划' "$runtime_skill/SKILL.md"
  grep -q 'Terraform Provider 用户问题排查' "$runtime_skill/agents/openai.yaml"
  grep -q '^## 1\. 现场状态$' "$runtime_skill/assets/report-template.md"
  grep -q '^# 证据与脱敏契约$' "$runtime_skill/references/evidence-contract.md"
done

for term in \
  "terraform plan -input=false -no-color -out=create.tfplan" \
  "TF_LOG_PROVIDER=DEBUG" \
  "DescribeFileSystems" \
  "replace_paths" \
  "禁止执行漂移或替换计划" \
  "html-report-preview" \
  "destroy-after-evidence"; do
  grep -q -- "$term" "$agents_skill/SKILL.md"
done

raw_log="$tmpdir/provider.log"
cat > "$raw_log" <<'EOF'
2026-07-20T11:45:51.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: *************** CreateFileSystem Response ***************
2026-07-20T11:45:51.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: {"FileSystemId":"fs-123","RequestId":"req-create"}Domain:, Version:
2026-07-20T11:45:51.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: ***************  Request ***************
2026-07-20T11:45:51.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: "{\"AccessKeyId\":\"SHOULD-NOT-LEAK\",\"FileSystemType\":\"standard\",\"VpcId\":\"vpc-123\",\"VSwitchId\":\"vsw-123\",\"ZoneId\":\"cn-test-a\"}"
2026-07-20T11:45:56.330+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: *************** DescribeFileSystems Response ***************
2026-07-20T11:45:56.330+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: {"FileSystems":{"FileSystem":[{"FileSystemId":"fs-123","Status":"Running","VpcId":""}]},"RequestId":"req-read"}Domain:, Version:
2026-07-20T11:46:01.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: *************** CreateDBInstance Response ***************
2026-07-20T11:46:01.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: {"DBInstanceId":"rm-123","RequestId":"req-db"}Domain:, Version:
2026-07-20T11:46:01.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: ***************  Request ***************
2026-07-20T11:46:01.109+0800 [DEBUG] provider.terraform-provider-alicloud_v1.285.0: "{\"AccessKeyId\":\"SHOULD-NOT-LEAK\",\"DBInstanceClass\":\"rds.mysql.s1.small\",\"Engine\":\"MySQL\"}"
EOF

report_md="$tmpdir/report.md"
printf '# 现场复现报告\n\n| API | RequestId |\n|---|---|\n| CreateFileSystem | req-create |\n' > "$report_md"
unsafe_md="$tmpdir/unsafe.md"
printf '# 不安全内容\n\n<script>alert(1)</script>\n' > "$unsafe_md"

runtime_index=0
timeline_args=(
  --request-field FileSystemType
  --request-field VpcId
  --request-field VSwitchId
  --request-field ZoneId
  --request-field DBInstanceClass
  --request-field Engine
  --target-field FileSystemId
  --target-field DBInstanceId
  --observe-field DescribeFileSystems:VpcId
  --observe-field DescribeFileSystems:QuorumVswId
)
for runtime_skill in "$agents_skill" "$claude_skill"; do
  runtime_index=$((runtime_index + 1))
  timeline="$tmpdir/timeline-$runtime_index.md"
  report_html="$tmpdir/report-$runtime_index.html"

  python3 "$runtime_skill/scripts/extract-api-timeline.py" "$raw_log" \
    "${timeline_args[@]}" --format markdown > "$timeline"
  grep -q 'req-create' "$timeline"
  grep -q 'VpcId.*vpc-123' "$timeline"
  grep -q 'QuorumVswId=missing' "$timeline"
  grep -q 'DBInstanceClass.*rds.mysql.s1.small' "$timeline"
  grep -q 'DBInstanceId.*rm-123' "$timeline"
  if grep -q 'SHOULD-NOT-LEAK' "$timeline"; then
    echo "reproduce_terraform_provider_issue_skill_test: parser leaked a forbidden request field" >&2
    exit 1
  fi

  default_timeline="$tmpdir/default-timeline-$runtime_index.md"
  python3 "$runtime_skill/scripts/extract-api-timeline.py" "$raw_log" > "$default_timeline"
  if grep -q 'FileSystemType\|DBInstanceClass\|fs-123\|rm-123' "$default_timeline"; then
    echo "reproduce_terraform_provider_issue_skill_test: parser emitted unselected fields" >&2
    exit 1
  fi

  if python3 "$runtime_skill/scripts/extract-api-timeline.py" "$raw_log" \
    --request-field AccessKeyId >/dev/null 2>&1; then
    echo "reproduce_terraform_provider_issue_skill_test: parser accepted a sensitive field" >&2
    exit 1
  fi

  python3 "$runtime_skill/scripts/extract-api-timeline.py" "$raw_log" \
    "${timeline_args[@]}" --format jsonl \
    | jq -e 'select(.request_id == "req-read" and (.observations | index("QuorumVswId=missing")))' >/dev/null

  python3 "$runtime_skill/scripts/render-report-html.py" "$report_md" "$report_html" >/dev/null
  grep -q '<title>现场复现报告</title>' "$report_html"
  grep -q '<table>' "$report_html"
  if grep -qi 'data:image' "$report_html"; then
    echo "reproduce_terraform_provider_issue_skill_test: renderer emitted base64 image data" >&2
    exit 1
  fi

  if python3 "$runtime_skill/scripts/render-report-html.py" "$unsafe_md" "$tmpdir/unsafe-$runtime_index.html" >/dev/null 2>&1; then
    echo "reproduce_terraform_provider_issue_skill_test: renderer accepted executable HTML" >&2
    exit 1
  fi
done

echo "reproduce_terraform_provider_issue_skill_test: PASS"
