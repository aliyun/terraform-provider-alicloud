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
    echo '{"success":true,"code":"SUCCESS","message":"成功","data":{"aoneId":"83843879","reportId":"rid-123","objectKey":"reports/aone/83843879/rid-123.html","viewUrl":"/reports/aone/83843879/rid-123/view","size":42}}'
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
html="$tmpdir/report.html"
printf '<html>ok</html>' > "$html"
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload 83843879 "$html" --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then assert_pass "upload exits 0"; else assert_fail "upload should exit 0, got $exit_code"; fi
assert_contains "$output" "https://pre.example/reports/aone/83843879/rid-123/view" "upload prints absolute view URL"
assert_file_contains "$curl_log" "file=@$html;type=text/html" "upload sends multipart file field"
assert_file_contains "$curl_log" "https://pre.example/api/reports/aone/83843879" "upload posts to AutomationAgent server-token endpoint"

echo "=== Test 2: upload sends server token when configured ==="
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" JARVIS_HTML_REPORT_TOKEN="test-token" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload 83843879 "$html" --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then assert_pass "token upload exits 0"; else assert_fail "token upload should exit 0, got $exit_code"; fi
assert_file_contains "$curl_log" "Authorization: Bearer test-token" "upload sends Authorization bearer token"

echo "=== Test 3: reject non-HTML before curl ==="
txt="$tmpdir/report.txt"
printf 'not html' > "$txt"
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload 83843879 "$txt" --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -ne 0 ]; then assert_pass "non-html exits nonzero"; else assert_fail "non-html should fail"; fi
if [ ! -s "$curl_log" ]; then assert_pass "non-html does not call curl"; else assert_fail "non-html should not call curl (log: $(cat "$curl_log"))"; fi

echo "=== Test 4: zip upload extracts all HTML members ==="
zip_path="$tmpdir/reports.zip"
make_zip "$zip_path"
: > "$curl_log"
output=$(CURL_LOG="$curl_log" JARVIS_CURL_BIN="$tmpbin/curl" \
    bash "$proj_root/bootstrap/html-report-preview.sh" upload 83843879 "$zip_path" --base-url https://pre.example 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then assert_pass "zip upload exits 0"; else assert_fail "zip upload should exit 0, got $exit_code"; fi
upload_calls=$(grep -c "api/reports/aone/83843879" "$curl_log" || true)
if [ "$upload_calls" -eq 2 ]; then assert_pass "zip uploads two HTML files"; else assert_fail "zip should upload 2 HTML files, got $upload_calls"; fi
assert_contains "$output" "report-a.html" "zip output includes first member label"
assert_contains "$output" "report-b.htm" "zip output includes nested member label"

echo "=== Test 5: from-aone uses latest attachment and can comment links back ==="
make_fake_a1id
: > "$curl_log"
: > "$a1_log"
output=$(A1_DOWNLOAD_FIXTURE="$zip_path" A1_LOG="$a1_log" CURL_LOG="$curl_log" \
JARVIS_A1_BIN="$tmpbin/a1id" JARVIS_CURL_BIN="$tmpbin/curl" \
    bash "$proj_root/bootstrap/html-report-preview.sh" from-aone 83843879 --base-url https://pre.example --comment 2>&1)
exit_code=$?
echo "$output"
if [ "$exit_code" -eq 0 ]; then assert_pass "from-aone exits 0"; else assert_fail "from-aone should exit 0, got $exit_code"; fi
assert_file_contains "$a1_log" "project workitem attachment list 83843879 -f json" "from-aone lists attachments"
assert_file_contains "$a1_log" "project workitem attachment download 83843879 222" "from-aone downloads latest attachment"
assert_file_contains "$a1_log" "project workitem comment create 83843879 -m" "from-aone comments when requested"
assert_contains "$output" "https://pre.example/reports/aone/83843879/rid-123/view" "from-aone prints absolute preview URL"

echo ""
echo "=== Summary ==="
echo "PASS: $PASS  FAIL: $FAIL"

if [ "$FAIL" -gt 0 ]; then
    echo "TESTS FAILED"
    exit 1
fi

echo "All tests passed"
exit 0
