#!/usr/bin/env bash
# test/html_report_preview_test.sh - tests for AutomationAgent HTML preview upload helper.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

tmpdir="$(mktemp -d)"
tmpbin="$tmpdir/bin"
curl_log="$tmpdir/curl.log"
a1_log="$tmpdir/a1.log"
mkdir -p "$tmpbin"
trap 'rm -rf "$tmpdir"' EXIT
# Keep tests hermetic: do not source the machine-level Jarvis runtime token.
export JARVIS_RUNTIME_CONFIG_LOADED=1
export JARVIS_ROOT="$proj_root"

PASS=0
FAIL=0

assert_pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
assert_fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

assert_contains() {
    local haystack="$1"
    local needle="$2"
    local label="$3"
    if grep -Fq -- "$needle" <<<"$haystack"; then
        assert_pass "$label"
    else
        assert_fail "$label (missing: $needle; got: $haystack)"
    fi
}

assert_file_contains() {
    local file="$1"
    local needle="$2"
    local label="$3"
    if grep -Fq -- "$needle" "$file"; then
        assert_pass "$label"
    else
        assert_fail "$label (missing: $needle; log: $(cat "$file" 2>/dev/null))"
    fi
}

make_fake_curl() {
    cat > "$tmpbin/curl" <<'STUB'
#!/usr/bin/env bash
echo "$@" >> "$CURL_LOG"
target="${@: -1}"
case "$target" in
  */api/reports/aone/*)
    if [ "${FAKE_REFLECT_TOKEN:-0}" = "1" ]; then
      echo '{"success":false,"detail":"Authorization: Bearer test-token"}'
    elif [ "${FAKE_BAD_VIEW_ROUTE:-0}" = "1" ]; then
      echo '{"success":true,"code":"SUCCESS","message":"成功","data":{"aoneId":"FAKE_AONE_ID","reportId":"rid-123","objectKey":"reports/aone/FAKE_AONE_ID/rid-123.html","viewUrl":"/buc/reports/aone/FAKE_AONE_ID/rid-123/view","size":42}}'
    else
      echo '{"success":true,"code":"SUCCESS","message":"成功","data":{"aoneId":"FAKE_AONE_ID","reportId":"rid-123","objectKey":"reports/aone/FAKE_AONE_ID/rid-123.html","viewUrl":"/reports/aone/FAKE_AONE_ID/rid-123/view","size":42}}'
    fi
    ;;
  *)
    echo "unexpected curl target: $target" >&2
    exit 22
    ;;
esac
STUB
    chmod +x "$tmpbin/curl"
}

make_fake_a1id() {
    cat > "$tmpbin/a1id" <<'STUB'
#!/usr/bin/env bash
[ "${1:-}" = "--" ] && shift
echo "$@" >> "$A1_LOG"
if [ "$1 $2 $3 $4" = "project workitem attachment list" ]; then
  cat <<JSON
[
  {"id":111,"name":"old-report.zip","created":"2026-07-03T10:00:00+08:00"},
  {"id":222,"name":"new-report.zip","created":"2026-07-03T12:00:00+08:00"}
]
JSON
  exit 0
fi
if [ "$1 $2 $3 $4" = "project workitem attachment download" ]; then
  out=""
  while [ "$#" -gt 0 ]; do
    if [ "${1:-}" = "-o" ]; then
      out="${2:-}"
      break
    fi
    shift
  done
  cp "$A1_DOWNLOAD_FIXTURE" "$out"
  exit 0
fi
if [ "$1 $2 $3 $4" = "project workitem comment create" ]; then
  exit 0
fi
echo "unexpected a1id call: $*" >&2
exit 2
STUB
    chmod +x "$tmpbin/a1id"
}

make_zip() {
    local zip_path="$1"
    python3 - "$zip_path" <<'PY'
import sys, zipfile
with zipfile.ZipFile(sys.argv[1], "w") as zf:
    zf.writestr("report-a.html", "<html>A</html>")
    zf.writestr("nested/report-b.htm", "<html>B</html>")
    zf.writestr("ignored.txt", "ignore")
PY
}

echo "=== Test 1: upload local HTML and print absolute preview URL ==="
make_fake_curl
make_fake_a1id
html="$tmpdir/report.html"
printf '<html>ok</html>' > "$html"
: > "$curl_log"
: > "$a1_log"
output=$(CURL_LOG="$curl_log" A1_LOG="$a1_log" JARVIS_CURL_BIN="$tmpbin/curl" \
    JARVIS_A1_BIN="$tmpbin/a1id" JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$html" --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then assert_pass "upload exits 0"; else assert_fail "upload should exit 0, got $exit_code"; fi
assert_contains "$output" "https://pre.example/reports/aone/FAKE_AONE_ID/rid-123/view" "upload prints absolute view URL"
assert_file_contains "$curl_log" "file=@$html;type=text/html" "upload sends multipart file field"
assert_file_contains "$curl_log" "https://pre.example/api/reports/aone/FAKE_AONE_ID" "upload posts to AutomationAgent server-token endpoint"
if [ ! -s "$a1_log" ]; then
    assert_pass "upload without --comment performs no Aone call"
else
    assert_fail "upload without --comment must not call Aone (log: $(cat "$a1_log"))"
fi

echo "=== Test 2: default base URL is pre-agent and JSONL records success state ==="
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    JARVIS_HTML_REPORT_BASE_URL="" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$html" --format jsonl 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then assert_pass "token upload exits 0"; else assert_fail "token upload should exit 0, got $exit_code"; fi
assert_file_contains "$curl_log" "Authorization: Bearer test-token" "upload sends Authorization bearer token"
assert_file_contains "$curl_log" "https://pre-agent.aliyun-inc.com/api/reports/aone/FAKE_AONE_ID" "default upload targets pre-agent"
if jq -e '.success == true and .status == "uploaded" and .reportId == "rid-123" and .viewUrl == "/reports/aone/FAKE_AONE_ID/rid-123/view"' >/dev/null <<<"$output"; then
    assert_pass "JSONL exposes uploaded status and readonly route"
else
    assert_fail "JSONL must expose success/uploaded/reportId/viewUrl (got: $output)"
fi

echo "=== Test 3: missing token blocks before curl or Aone ==="
: > "$curl_log"
: > "$a1_log"
output=$(CURL_LOG="$curl_log" A1_LOG="$a1_log" JARVIS_CURL_BIN="$tmpbin/curl" \
    JARVIS_A1_BIN="$tmpbin/a1id" JARVIS_HTML_REPORT_TOKEN="" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$html" --format jsonl 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 3 ]; then assert_pass "missing token exits 3"; else assert_fail "missing token should exit 3, got $exit_code"; fi
if jq -e '.success == false and .status == "blocked" and .code == "missing_token"' >/dev/null <<<"$output"; then
    assert_pass "missing token emits structured blocked record"
else
    assert_fail "missing token JSONL contract mismatch (got: $output)"
fi
if [ ! -s "$curl_log" ] && [ ! -s "$a1_log" ]; then
    assert_pass "missing token performs no curl or Aone call"
else
    assert_fail "missing token must not call curl/Aone"
fi

echo "=== Test 4: reject non-HTML before curl ==="
txt="$tmpdir/report.txt"
printf 'not html' > "$txt"
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$txt" --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -ne 0 ]; then assert_pass "non-html exits nonzero"; else assert_fail "non-html should fail"; fi
if [ ! -s "$curl_log" ]; then assert_pass "non-html does not call curl"; else assert_fail "non-html should not call curl (log: $(cat "$curl_log"))"; fi

echo "=== Test 5: zip upload extracts all HTML members ==="
zip_path="$tmpdir/reports.zip"
make_zip "$zip_path"
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$zip_path" --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then assert_pass "zip upload exits 0"; else assert_fail "zip upload should exit 0, got $exit_code"; fi
upload_calls=$(grep -c "api/reports/aone/FAKE_AONE_ID" "$curl_log" || true)
if [ "$upload_calls" -eq 2 ]; then assert_pass "zip uploads two HTML files"; else assert_fail "zip should upload 2 HTML files, got $upload_calls"; fi
assert_contains "$output" "report-a.html" "zip output includes first member label"
assert_contains "$output" "report-b.htm" "zip output includes nested member label"

echo "=== Test 6: from-aone uses latest attachment and can comment links back ==="
make_fake_a1id
: > "$curl_log"
: > "$a1_log"
output=$(A1_DOWNLOAD_FIXTURE="$zip_path" A1_LOG="$a1_log" CURL_LOG="$curl_log" \
JARVIS_A1_BIN="$tmpbin/a1id" JARVIS_CURL_BIN="$tmpbin/curl" \
JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" from-aone FAKE_AONE_ID --base-url https://pre.example --comment 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then assert_pass "from-aone exits 0"; else assert_fail "from-aone should exit 0, got $exit_code"; fi
assert_file_contains "$a1_log" "project workitem attachment list FAKE_AONE_ID -f json" "from-aone lists attachments"
assert_file_contains "$a1_log" "project workitem attachment download FAKE_AONE_ID 222" "from-aone downloads latest attachment"
assert_file_contains "$a1_log" "project workitem comment create FAKE_AONE_ID -m" "from-aone comments when requested"
assert_contains "$output" "https://pre.example/reports/aone/FAKE_AONE_ID/rid-123/view" "from-aone prints absolute preview URL"
# Aone 评论渲染 quirk 契约:评论区按 markdown 渲染,可点击链接唯一格式 = [text](url)
# (84307546 评论 124870464 四格式对照实测:裸 URL/HTML 锚均为死文本)
if grep -Fq -- "](https://" "$a1_log"; then
    assert_pass "comment uses markdown [text](url) links (the only clickable form)"
else
    assert_fail "comment must use markdown [text](url) links (a1 log: $(cat "$a1_log" 2>/dev/null))"
fi

echo "=== Test 7: reject a non-readonly viewer route ==="
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    FAKE_BAD_VIEW_ROUTE=1 \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$html" --base-url https://pre.example --format jsonl 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -ne 0 ]; then
    assert_pass "wrong viewer route is rejected"
else
    assert_fail "wrong viewer route must fail"
fi

echo "=== Test 8: sanitize a rejected response that reflects credentials ==="
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    FAKE_REFLECT_TOKEN=1 \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$html" \
    --base-url https://pre.example --format jsonl 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -ne 0 ]; then
    assert_pass "reflected failure exits nonzero"
else
    assert_fail "reflected failure must fail"
fi
if jq -e '.success == false and .status == "failed" and .code == "upload_rejected"' >/dev/null <<<"$output"; then
    assert_pass "reflected failure emits fixed sanitized JSON"
else
    assert_fail "reflected failure JSON contract mismatch (got: $output)"
fi
if grep -Fq -- "test-token" <<<"$output"; then
    assert_fail "reflected failure leaked the token"
else
    assert_pass "reflected failure does not leak the token"
fi

echo "=== Test 9: accept only absolute HTTPS report image references ==="
safe_images="$tmpdir/safe-images.html"
cat >"$safe_images" <<'HTML'
<html><body>
  <img src="https://images.example/openapi.png?Expires=1&amp;Signature=ok"
       srcset="https://images.example/openapi.png 1x, https://images.example/openapi@2x.png 2x">
  <picture>
    <source srcset="https://images.example/cloudspec.webp 1x, https://images.example/cloudspec@2x.webp 2x">
    <img src="https://images.example/provider.png">
  </picture>
  <source srcset="relative-non-picture-media.mp4">
</body></html>
HTML
: >"$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$safe_images" \
    --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then
    assert_pass "absolute HTTPS img/picture references are accepted"
else
    assert_fail "absolute HTTPS img/picture references should upload (got $exit_code)"
fi

assert_rejected_image_reference() {
    local name="$1"
    local markup="$2"
    local fixture="$tmpdir/reject-$name.html"
    printf '%s\n' "$markup" >"$fixture"
    : >"$curl_log"
    output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" \
        JARVIS_HTML_REPORT_TOKEN="test-token" \
        bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$fixture" \
        --base-url https://pre.example 2>&1)
    exit_code=$?
    if [ "$exit_code" -ne 0 ]; then
        assert_pass "$name image reference is rejected"
    else
        assert_fail "$name image reference must be rejected"
    fi
    assert_contains "$output" "invalid_image_reference" "$name uses a fixed fail-closed diagnostic"
    if [ ! -s "$curl_log" ]; then
        assert_pass "$name is rejected before curl"
    else
        assert_fail "$name must not call curl (log: $(cat "$curl_log"))"
    fi
}

assert_rejected_image_reference relative-src '<img src="screenshots/openapi.png">'
assert_rejected_image_reference file-src '<img src="file:///tmp/openapi.png">'
assert_rejected_image_reference data-src '<img src="data:image/png;base64,AAAA">'
assert_rejected_image_reference protocol-relative-src '<img src="//images.example/openapi.png">'
assert_rejected_image_reference http-src '<img src="http://images.example/openapi.png">'
assert_rejected_image_reference relative-img-srcset \
    '<img src="https://images.example/openapi.png" srcset="https://images.example/openapi.png 1x, screenshots/openapi.png 2x">'
assert_rejected_image_reference relative-picture-srcset \
    '<picture><source srcset="screenshots/openapi.webp 1x"><img src="https://images.example/openapi.png"></picture>'

echo "=== Test 10: preflight the whole ZIP before uploading any member ==="
unsafe_zip="$tmpdir/unsafe-reports.zip"
python3 - "$unsafe_zip" <<'PY'
import sys
import zipfile

with zipfile.ZipFile(sys.argv[1], "w") as archive:
    archive.writestr("001-safe.html", '<img src="https://images.example/openapi.png">')
    archive.writestr("002-unsafe.html", '<img src="screenshots/cloudspec.png">')
PY
: >"$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$unsafe_zip" \
    --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -ne 0 ]; then
    assert_pass "unsafe ZIP is rejected"
else
    assert_fail "unsafe ZIP must be rejected"
fi
assert_contains "$output" "invalid_image_reference" "unsafe ZIP uses fixed fail-closed diagnostic"
if [ ! -s "$curl_log" ]; then
    assert_pass "unsafe ZIP performs no partial uploads"
else
    assert_fail "unsafe ZIP must be fully validated before curl (log: $(cat "$curl_log"))"
fi

echo "=== Test 11: reject screenshots named in prose but never embedded ==="
assert_rejected_prose_screenshot() {
    local name="$1"
    local markup="$2"
    local fixture="$tmpdir/prose-$name.html"
    printf '%s\n' "$markup" >"$fixture"
    : >"$curl_log"
    output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" \
        JARVIS_HTML_REPORT_TOKEN="test-token" \
        bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$fixture" \
        --base-url https://pre.example 2>&1)
    exit_code=$?
    if [ "$exit_code" -ne 0 ]; then
        assert_pass "$name prose-only screenshot is rejected"
    else
        assert_fail "$name prose-only screenshot must be rejected"
    fi
    assert_contains "$output" "screenshot_prose_without_img" \
        "$name uses the fixed prose-without-img diagnostic"
    if [ ! -s "$curl_log" ]; then
        assert_pass "$name is rejected before curl"
    else
        assert_fail "$name must not call curl (log: $(cat "$curl_log"))"
    fi
}

# The exact shape that shipped a screenshot-less "visual evidence" report.
assert_rejected_prose_screenshot cn-caption \
    '<html><body><p>截图：openapi_create_artifact_subscription_rule.png（见 manifest）</p></body></html>'
assert_rejected_prose_screenshot en-caption \
    '<html><body><td>screenshot: provider_cr_resources.jpg</td></body></html>'

assert_accepted_html() {
    local name="$1"
    local markup="$2"
    local fixture="$tmpdir/accept-$name.html"
    printf '%s\n' "$markup" >"$fixture"
    : >"$curl_log"
    output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" \
        JARVIS_HTML_REPORT_TOKEN="test-token" \
        bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$fixture" \
        --base-url https://pre.example 2>&1)
    exit_code=$?
    if [ "$exit_code" -eq 0 ]; then
        assert_pass "$name is accepted"
    else
        assert_fail "$name must upload (got $exit_code: $output)"
    fi
}

# No false positives: a genuinely text-only report, and the documented
# missing_capability degradation, must both still upload.
assert_accepted_html text-only \
    '<html><body><p>三层查证均通过，无需截图。</p></body></html>'
assert_accepted_html degraded-missing-capability \
    '<html><body><p>截图降级：openapi.png 未生成，missing_capability: 无可用浏览器通道</p></body></html>'
assert_accepted_html prose-plus-real-img \
    '<html><body><p>截图：openapi.png</p><img src="https://images.example/openapi.png"></body></html>'

echo "=== Test 12: reject inputs whose sibling images would be silently dropped ==="
img_dir="$tmpdir/report-with-images"
mkdir -p "$img_dir"
printf '<html><body><p>见下方截图</p></body></html>\n' >"$img_dir/report.html"
printf 'fake-png-bytes' >"$img_dir/openapi_shot.png"
: >"$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$img_dir" \
    --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -ne 0 ]; then
    assert_pass "directory carrying images is rejected"
else
    assert_fail "directory carrying images must be rejected"
fi
assert_contains "$output" "image_files_dropped" "directory uses the fixed dropped-image diagnostic"
if [ ! -s "$curl_log" ]; then
    assert_pass "dropped-image directory performs no uploads"
else
    assert_fail "dropped-image directory must not call curl (log: $(cat "$curl_log"))"
fi

img_zip="$tmpdir/report-with-images.zip"
python3 - "$img_zip" <<'PY'
import sys
import zipfile

with zipfile.ZipFile(sys.argv[1], "w") as archive:
    archive.writestr("report.html", "<html><body><p>见下方截图</p></body></html>")
    archive.writestr("openapi_shot.png", "fake-png-bytes")
PY
: >"$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload FAKE_AONE_ID "$img_zip" \
    --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -ne 0 ]; then
    assert_pass "ZIP carrying images is rejected"
else
    assert_fail "ZIP carrying images must be rejected"
fi
assert_contains "$output" "image_files_dropped" "ZIP uses the fixed dropped-image diagnostic"

echo ""
echo "=== Summary ==="
echo "PASS: $PASS  FAIL: $FAIL"

if [ "$FAIL" -gt 0 ]; then
    echo "TESTS FAILED"
    exit 1
fi

echo "All tests passed"
exit 0
