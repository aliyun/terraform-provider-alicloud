#!/bin/bash
set -euo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
skill_rel="reproduce-terraform-provider-issue"
agents_skill="$repo_root/.agents/skills/$skill_rel"
claude_skill="$repo_root/.claude/skills/$skill_rel"
tmpdir="$(mktemp -d)"
server_pid=""
cleanup() {
  if [ -n "$server_pid" ]; then
    kill "$server_pid" 2>/dev/null || true
    wait "$server_pid" 2>/dev/null || true
  fi
  rm -rf "$tmpdir"
}
trap cleanup EXIT

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
  test -x "$runtime_skill/scripts/validate-report-package.py"
  grep -q '^name: reproduce-terraform-provider-issue$' "$runtime_skill/SKILL.md"
  grep -q '^description: ' "$runtime_skill/SKILL.md"
  grep -q '^# Terraform Provider 用户问题现场排查复现$' "$runtime_skill/SKILL.md"
  grep -q '禁止执行漂移或替换计划' "$runtime_skill/SKILL.md"
  grep -q 'Terraform Provider 用户问题排查' "$runtime_skill/agents/openai.yaml"
  grep -q '\$reproduce-terraform-provider-issue' "$runtime_skill/agents/openai.yaml"
  grep -q '^## 1\. 现场状态$' "$runtime_skill/assets/report-template.md"
  grep -q '^## 14\. 完整 Terraform HCL$' "$runtime_skill/assets/report-template.md"
  grep -q '^## 15\. 在线预填与交付验证$' "$runtime_skill/assets/report-template.md"
  grep -q '^# 证据与脱敏契约$' "$runtime_skill/references/evidence-contract.md"
done

for term in \
  "terraform plan -input=false -no-color -out=create.tfplan" \
  "TF_LOG_PROVIDER=DEBUG" \
  "DescribeFileSystems" \
  "replace_paths" \
  "禁止执行漂移或替换计划" \
  "html-report-preview" \
  "destroy-after-evidence" \
  "validate-report-package.py" \
  "platform_blocked" \
  "REPORT.md" \
  "template/main.tf"; do
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
base64_md="$tmpdir/base64.md"
printf '# 不安全图片\n\n![x](data:image/png;base64,AAAA)\n' > "$base64_md"
entity_bypass_md="$tmpdir/entity-bypass.md"
printf '# Entity bypass\n\n<a href="java&#x73;cript:alert(1)">x</a>\n' > "$entity_bypass_md"

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
  if python3 "$runtime_skill/scripts/render-report-html.py" "$base64_md" "$tmpdir/base64-$runtime_index.html" >/dev/null 2>&1; then
    echo "reproduce_terraform_provider_issue_skill_test: renderer accepted base64 image data" >&2
    exit 1
  fi
  if python3 "$runtime_skill/scripts/render-report-html.py" "$entity_bypass_md" "$tmpdir/entity-$runtime_index.html" >/dev/null 2>&1; then
    echo "reproduce_terraform_provider_issue_skill_test: renderer accepted an entity-encoded executable URL" >&2
    exit 1
  fi

  python3 "$runtime_skill/scripts/render-report-html.py" \
    --lang en-US "$report_md" "$tmpdir/en-$runtime_index.html" >/dev/null
  grep -q '<html lang="en-US">' "$tmpdir/en-$runtime_index.html"
  if python3 "$runtime_skill/scripts/render-report-html.py" \
    --lang '"><script>' "$report_md" "$tmpdir/bad-lang-$runtime_index.html" >/dev/null 2>&1; then
    echo "reproduce_terraform_provider_issue_skill_test: renderer accepted an unsafe language tag" >&2
    exit 1
  fi
done

fake_bin="$tmpdir/fake-bin"
runtime_tmp="$tmpdir/runtime"
mkdir -p "$fake_bin" "$runtime_tmp"
export TMPDIR="$runtime_tmp"
cat > "$fake_bin/terraform" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' "$*" >> "$TMPDIR/terraform.log"
env | sort > "$TMPDIR/terraform.env"
case "${1:-}" in
  fmt)
    [ ! -f "$TMPDIR/fmt.fail" ] || exit 1
    ;;
  init)
    if [ -f "$TMPDIR/init-external.fail" ] || [ -f "$TMPDIR/init-invalid.fail" ]; then
      if [ -f "$TMPDIR/init-external.fail" ]; then
        echo "failed to query available provider packages" >&2
      else
        echo "invalid provider configuration" >&2
      fi
      exit 1
    fi
    ;;
  validate)
    [ ! -f "$TMPDIR/validate.fail" ] || exit 1
    ;;
esac
exit 0
EOF
chmod +x "$fake_bin/terraform"

make_package() {
  local package_dir="$1"
  mkdir -p "$package_dir/template"
  cat > "$package_dir/template/main.tf" <<'EOF'
terraform {
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "1.285.0"
    }
  }
}

variable "profile" {
  type     = string
  default  = null
  nullable = true
}

provider "alicloud" {
  profile = var.profile
  region  = "cn-hangzhou"
}
EOF
  cat > "$package_dir/template/README.md" <<'EOF'
# Reproduction template

Use `TF_VAR_profile` locally. Never put credentials in this package.
EOF
  python3 - "$package_dir" <<'PY'
import pathlib
import urllib.parse
import sys

root = pathlib.Path(sys.argv[1])
hcl = (root / "template/main.tf").read_text()
encoded = urllib.parse.quote_plus(hcl, safe="*-._").replace("+", "%20")
ts = "1721700000123"
url = (
    "https://api.aliyun.com/terraform?"
    "spm=XToCode.TerraformAI.QA.0&activeTab=code&source=PlayGround&"
    f"sourcePath=TerraformAI/{ts}::{ts}&params={encoded}"
)
(root / "REPORT.md").write_text(
    "# 客户问题复现报告\n\n"
    "<!-- REPORT_PACKAGE_VALIDATED -->\n\n"
    "## 完整 Terraform HCL\n\n"
    f"```hcl\n{hcl}```\n\n"
    f"[打开在线 Terraform 调试页]({url})\n",
    encoding="utf-8",
)
PY
  python3 "$agents_skill/scripts/render-report-html.py" \
    "$package_dir/REPORT.md" "$package_dir/REPORT.html" >/dev/null
}

package_dir="$tmpdir/package"
make_package "$package_dir"
terraform_log="$runtime_tmp/terraform.log"
: > "$terraform_log"
validator="$agents_skill/scripts/validate-report-package.py"
output=$(PATH="$fake_bin:$PATH" \
  TF_CLI_ARGS="-chdir=/test-only-invalid" \
  ALICLOUD_ACCESS_KEY="test-only-not-a-real-key" \
  TF_VAR_profile="test-only-profile" \
  python3 "$validator" "$package_dir" --format json)
jq -e '.success == true and .status == "validated" and .viewer_copy.status == "platform_blocked"' >/dev/null <<<"$output"
grep -Fq -- 'fmt -check -recursive' "$terraform_log"
grep -Fq -- 'init -backend=false -input=false -no-color' "$terraform_log"
grep -Fq -- 'validate -no-color' "$terraform_log"
if grep -Eq '^(TF_CLI_ARGS|TF_VAR_|ALICLOUD_|ALIBABA_CLOUD_|API_TOKEN)=' "$runtime_tmp/terraform.env"; then
  echo "reproduce_terraform_provider_issue_skill_test: terraform subprocess inherited sensitive env" >&2
  exit 1
fi

lock_package="$tmpdir/lock-package"
cp -R "$package_dir" "$lock_package"
: > "$lock_package/template/.terraform.lock.hcl"
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$lock_package" --format json \
  | jq -e '.success == true' >/dev/null
grep -Fq -- '-lockfile=readonly' "$terraform_log"

set +e
: > "$runtime_tmp/fmt.fail"
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --format json >/dev/null
fmt_rc=$?
rm -f "$runtime_tmp/fmt.fail"
: > "$runtime_tmp/init-external.fail"
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --format json >/dev/null
init_rc=$?
rm -f "$runtime_tmp/init-external.fail"
: > "$runtime_tmp/init-invalid.fail"
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --format json >/dev/null
init_invalid_rc=$?
rm -f "$runtime_tmp/init-invalid.fail"
: > "$runtime_tmp/validate.fail"
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --format json >/dev/null
validate_rc=$?
rm -f "$runtime_tmp/validate.fail"
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --require-preview --format json >/dev/null
preview_missing_rc=$?
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --require-viewer-copy --format json >/dev/null
viewer_rc=$?
set -e
test "$fmt_rc" -eq 2
test "$init_rc" -eq 3
test "$init_invalid_rc" -eq 2
test "$validate_rc" -eq 2
test "$preview_missing_rc" -eq 3
test "$viewer_rc" -eq 3

unsafe_package="$tmpdir/unsafe-package"
cp -R "$package_dir" "$unsafe_package"
printf '\n<script>alert(1)</script>\n' >> "$unsafe_package/REPORT.html"
set +e
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$unsafe_package" --format json >/dev/null
unsafe_rc=$?
set -e
test "$unsafe_rc" -eq 2

entity_package="$tmpdir/entity-package"
cp -R "$package_dir" "$entity_package"
printf '\n<a href="java&#x73;cript:alert(1)">x</a>\n' >> "$entity_package/REPORT.html"
set +e
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$entity_package" --format json >/dev/null
entity_rc=$?
set -e
test "$entity_rc" -eq 2

credential_package="$tmpdir/credential-package"
cp -R "$package_dir" "$credential_package"
cat > "$credential_package/actual-values.txt" <<'EOF'
{"accessKeyId":"AKID_TEST_ONLY_123456","accessKeySecret":"test-only-secret-value","api_token":"test-only-token-value"}
EOF
set +e
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$credential_package" --format json >/dev/null
credential_rc=$?
set -e
test "$credential_rc" -eq 2

field_names_package="$tmpdir/field-names-package"
cp -R "$package_dir" "$field_names_package"
cat > "$field_names_package/field-names.txt" <<'EOF'
The API documentation names accessKeyId, accessKeySecret, and api_token fields.
EOF
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$field_names_package" --format json \
  | jq -e '.success == true' >/dev/null

raw_log_package="$tmpdir/raw-log-package"
cp -R "$package_dir" "$raw_log_package"
cat > "$raw_log_package/provider-debug.txt" <<'EOF'
2026-01-02T03:04:05.000+0000 [DEBUG] provider.terraform-provider-example: raw response body
EOF
set +e
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$raw_log_package" --format json >/dev/null
raw_log_rc=$?
set -e
test "$raw_log_rc" -eq 2

hcl_mismatch_package="$tmpdir/hcl-mismatch-package"
cp -R "$package_dir" "$hcl_mismatch_package"
printf '\n# byte mismatch\n' >> "$hcl_mismatch_package/template/main.tf"
set +e
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$hcl_mismatch_package" --format json >/dev/null
hcl_mismatch_rc=$?
set -e
test "$hcl_mismatch_rc" -eq 2

double_package="$tmpdir/double-package"
cp -R "$package_dir" "$double_package"
python3 - "$double_package" <<'PY'
import pathlib
import sys

root = pathlib.Path(sys.argv[1])
for name in ("REPORT.md", "REPORT.html"):
    path = root / name
    text = path.read_text()
    marker = "params="
    start = text.index(marker) + len(marker)
    end_candidates = [i for i in (text.find('"', start), text.find(")", start)) if i >= 0]
    end = min(end_candidates)
    encoded = text[start:end].replace("%", "%25")
    path.write_text(text[:start] + encoded + text[end:])
PY
set +e
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$double_package" --format json >/dev/null
double_rc=$?
set -e
test "$double_rc" -eq 2

serve_dir="$tmpdir/serve"
mkdir -p "$serve_dir"
cp "$package_dir/REPORT.html" "$serve_dir/report.html"
port_file="$tmpdir/http.port"
header_file="$tmpdir/http.headers"
cat > "$tmpdir/http_server.py" <<'PY'
import http.server
import pathlib
import socketserver
import sys

root = pathlib.Path(sys.argv[1])
port_file = pathlib.Path(sys.argv[2])
header_file = pathlib.Path(sys.argv[3])
class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=str(root), **kwargs)
    def do_GET(self):
        header_file.write_text(
            "".join(f"{key}: {value}\n" for key, value in self.headers.items())
        )
        if (root / "redirect.enabled").exists():
            self.send_response(302)
            self.send_header("Location", "/different-path.html")
            self.end_headers()
            return
        body = (root / "report.html").read_bytes()
        self.send_response(200)
        self.send_header("Content-Type", "text/html; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)
    def log_message(self, *_args):
        pass
with socketserver.TCPServer(("127.0.0.1", 0), Handler) as server:
    port_file.write_text(str(server.server_address[1]))
    server.serve_forever()
PY
python3 "$tmpdir/http_server.py" "$serve_dir" "$port_file" "$header_file" &
server_pid=$!
for _ in $(seq 1 50); do
  test -s "$port_file" && break
  sleep 0.02
done
test -s "$port_file"
port="$(cat "$port_file")"
fake_aone_id="FAKE_AONE_ID"
readonly_path="/reports/aone/$fake_aone_id/rid-123/view"
local_origin="http://127.0.0.1:$port"

assert_preview_record_invalid() {
  local label="$1"
  local absolute_url="$2"
  local expected_origin="${3:-$local_origin}"
  local view_url="${4:-$readonly_path}"
  local record="$tmpdir/invalid-preview-$label.json"
  cat > "$record" <<EOF
{"success":true,"status":"uploaded","reportId":"rid-123","viewUrl":"$view_url","url":"$absolute_url"}
EOF
  set +e
  PATH="$fake_bin:$PATH" \
    python3 "$validator" "$package_dir" --require-preview \
    --preview-json "$record" --preview-origin "$expected_origin" \
    --preview-marker "客户问题复现报告" --format json >/dev/null
  local rc=$?
  set -e
  if [ "$rc" -ne 2 ]; then
    echo "reproduce_terraform_provider_issue_skill_test: preview bypass '$label' returned $rc, want 2" >&2
    exit 1
  fi
}

assert_preview_record_invalid "path" "$local_origin/report.html"
assert_preview_record_invalid "host" "http://localhost:$port$readonly_path"
assert_preview_record_invalid "scheme" "$local_origin$readonly_path" "https://127.0.0.1:$port"
assert_preview_record_invalid "query" "$local_origin$readonly_path?marker=客户问题复现报告"
assert_preview_record_invalid "fragment" "$local_origin$readonly_path#客户问题复现报告"
assert_preview_record_invalid \
  "percent-path" "$local_origin/reports/aone/$fake_aone_id/rid-123/%76iew"
assert_preview_record_invalid \
  "percent-view-url" "$local_origin/reports/aone/$fake_aone_id/rid-123/%76iew" \
  "$local_origin" "/reports/aone/$fake_aone_id/rid-123/%76iew"

preview_json="$tmpdir/preview.json"
cat > "$preview_json" <<EOF
{"success":true,"status":"uploaded","reportId":"rid-123","viewUrl":"$readonly_path","url":"http://127.0.0.1:$port$readonly_path"}
EOF
set +e
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --require-preview \
  --preview-json "$preview_json" --preview-origin "$local_origin" --format json >/dev/null
no_marker_rc=$?
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --require-preview \
  --preview-json "$preview_json" --preview-origin "$local_origin" \
  --preview-marker "" --format json >/dev/null
empty_marker_rc=$?
set -e
test "$no_marker_rc" -eq 2
test "$empty_marker_rc" -eq 2

: > "$serve_dir/redirect.enabled"
set +e
PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --require-preview \
  --preview-json "$preview_json" --preview-origin "$local_origin" \
  --preview-marker "客户问题复现报告" --format json >/dev/null
redirect_rc=$?
set -e
rm -f "$serve_dir/redirect.enabled"
test "$redirect_rc" -eq 3

set +e
preview_output=$(PATH="$fake_bin:$PATH" \
  python3 "$validator" "$package_dir" --require-preview \
  --preview-json "$preview_json" --preview-origin "$local_origin" \
  --preview-marker "客户问题复现报告" --format json)
preview_rc=$?
set -e
if [ "$preview_rc" -ne 0 ] ||
   ! jq -e '.success == true and .preview.status == "verified"' >/dev/null <<<"$preview_output"; then
  echo "reproduce_terraform_provider_issue_skill_test: preview validation failed: $preview_output" >&2
  exit 1
fi
if grep -Eiq '^(Authorization|Cookie):' "$header_file"; then
  echo "reproduce_terraform_provider_issue_skill_test: anonymous preview GET leaked auth headers" >&2
  exit 1
fi
kill "$server_pid"
wait "$server_pid" 2>/dev/null || true
server_pid=""

echo "reproduce_terraform_provider_issue_skill_test: PASS"
