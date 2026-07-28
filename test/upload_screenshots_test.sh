#!/usr/bin/env bash
# Hermetic contract tests for the screenshot upload client.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
script="$repo_root/.agents/skills/screenshot-evidence/scripts/upload-screenshots.sh"
tmpdir="$(mktemp -d)"
tmpbin="$tmpdir/bin"
shots="$tmpdir/shots"
curl_log="$tmpdir/curl.log"
stdout_file="$tmpdir/stdout"
stderr_file="$tmpdir/stderr"
mkdir -p "$tmpbin" "$shots"
trap 'rm -rf "$tmpdir"' EXIT

PASS=0
FAIL=0

pass() { printf 'PASS: %s\n' "$1"; PASS=$((PASS + 1)); }
fail() { printf 'FAIL: %s\n' "$1"; FAIL=$((FAIL + 1)); }

assert_exit() {
  local expected="$1" actual="$2" label="$3"
  if [ "$actual" -eq "$expected" ]; then
    pass "$label"
  else
    fail "$label (expected $expected, got $actual)"
  fi
}

assert_no_secret() {
  local label="$1"
  if grep -Fq -- "unit-test-bearer-secret" "$stdout_file" "$stderr_file"; then
    fail "$label"
  else
    pass "$label"
  fi
}

assert_stderr_exact() {
  local expected="$1" label="$2" actual
  actual="$(cat "$stderr_file")"
  if [ "$actual" = "$expected" ]; then
    pass "$label"
  else
    fail "$label (got: $actual)"
  fi
}

cat >"$tmpbin/curl" <<'STUB'
#!/usr/bin/env bash
response_file=""
upload_file=""
url=""
http_code="${FAKE_HTTP_CODE:-200}"
bearer_seen=0

while [ "$#" -gt 0 ]; do
  case "$1" in
    -o|--output)
      response_file="${2:-}"
      shift 2
      ;;
    -w|--write-out|-H|--header|-F|--form|-X|--request)
      value="${2:-}"
      case "$1" in
        -H|--header)
          case "$value" in
            "Authorization: Bearer "*) bearer_seen=1 ;;
          esac
          ;;
        -F|--form)
          case "$value" in
            file=@*) upload_file="${value#file=@}" ;;
          esac
          ;;
      esac
      shift 2
      ;;
    -s|-S|-sS|--silent|--show-error)
      shift
      ;;
    http://*|https://*)
      url="$1"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

[ -n "$response_file" ] || exit 91
[ -n "$upload_file" ] || exit 92
[ "$bearer_seen" -eq 1 ] || exit 93

printf 'url=%s file=%s bearer=present\n' "$url" "$(basename "$upload_file")" >>"$CURL_LOG"
if [ "${FAKE_RESPONSE_MODE:-success}" = "reflected_failure" ]; then
  printf '%s\n' '{"success":false,"message":"Authorization: Bearer unit-test-bearer-secret"}' >"$response_file"
else
  file_name="$(basename "$upload_file")"
  key_name="${file_name%.*}"
  printf '{"success":true,"code":"SUCCESS","data":{"originalFilename":"%s","objectKey":"reports/aone/123/images/%s","signedUrl":"https://jarvis-upload-files.example/%s?signature=fake"}}\n' \
    "$file_name" "$file_name" "$key_name" >"$response_file"
fi
printf '%s' "$http_code"
STUB
chmod +x "$tmpbin/curl"

printf 'png' >"$shots/openapi.png"
printf 'jpg' >"$shots/cloudspec.jpg"
printf 'jpeg' >"$shots/provider.jpeg"
printf 'upper' >"$shots/schema.PNG"
printf 'ignore' >"$shots/notes.txt"

echo "=== success: upload every PNG/JPG and emit name|signed_url ==="
: >"$curl_log"
JARVIS_RUNTIME_CONFIG_LOADED=1 \
JARVIS_CURL_BIN="$tmpbin/curl" \
JARVIS_HTML_REPORT_TOKEN="unit-test-bearer-secret" \
JARVIS_HTML_REPORT_BASE_URL="https://pre.example/" \
CURL_LOG="$curl_log" \
bash "$script" 123 "$shots" >"$stdout_file" 2>"$stderr_file"
rc=$?
assert_exit 0 "$rc" "successful upload exits 0"
if [ "$(wc -l <"$stdout_file" | tr -d ' ')" -eq 4 ]; then
  pass "all PNG/JPG variants produce output"
else
  fail "expected four uploaded image records (got: $(tr '\n' ';' <"$stdout_file"))"
fi
for expected in \
  "openapi|https://jarvis-upload-files.example/openapi?signature=fake" \
  "cloudspec|https://jarvis-upload-files.example/cloudspec?signature=fake" \
  "provider|https://jarvis-upload-files.example/provider?signature=fake" \
  "schema|https://jarvis-upload-files.example/schema?signature=fake"; do
  if grep -Fxq -- "$expected" "$stdout_file"; then
    pass "output contains $expected"
  else
    fail "output missing $expected"
  fi
done
if [ "$(grep -c 'api/reports/aone/123/images' "$curl_log" || true)" -eq 4 ]; then
  pass "client posts each image to the server-token image API"
else
  fail "unexpected upload endpoint/call count"
fi
assert_no_secret "bearer token is absent from stdout/stderr"

echo "=== missing token: fail before curl ==="
: >"$curl_log"
JARVIS_RUNTIME_CONFIG_LOADED=1 \
JARVIS_CURL_BIN="$tmpbin/curl" \
JARVIS_HTML_REPORT_TOKEN="" \
JARVIS_HTML_REPORT_BASE_URL="https://pre.example" \
CURL_LOG="$curl_log" \
bash "$script" 123 "$shots" >"$stdout_file" 2>"$stderr_file"
rc=$?
assert_exit 3 "$rc" "missing token exits 3"
assert_stderr_exact "upload-screenshots: missing_token" "missing token uses fixed sanitized diagnostic"
if [ ! -s "$curl_log" ]; then
  pass "missing token performs no upload"
else
  fail "missing token called curl"
fi
assert_no_secret "missing-token diagnostics contain no bearer token"

echo "=== server failure: fail closed and sanitize reflected response ==="
: >"$curl_log"
JARVIS_RUNTIME_CONFIG_LOADED=1 \
JARVIS_CURL_BIN="$tmpbin/curl" \
JARVIS_HTML_REPORT_TOKEN="unit-test-bearer-secret" \
JARVIS_HTML_REPORT_BASE_URL="https://pre.example" \
FAKE_HTTP_CODE=500 \
FAKE_RESPONSE_MODE=reflected_failure \
CURL_LOG="$curl_log" \
bash "$script" 123 "$shots" >"$stdout_file" 2>"$stderr_file"
rc=$?
if [ "$rc" -ne 0 ]; then
  pass "server failure exits nonzero"
else
  fail "server failure must exit nonzero"
fi
if [ ! -s "$stdout_file" ]; then
  pass "server failure emits no partial URL records"
else
  fail "server failure emitted partial output"
fi
assert_stderr_exact "upload-screenshots: upload_failed" "server failure uses fixed sanitized diagnostic"
assert_no_secret "server response cannot reflect bearer token to stdout/stderr"

echo "=== implementation dependency contract ==="
if grep -Eq '(^|[^[:alnum:]_])aliyun([[:space:]]|$)' "$script"; then
  fail "upload client must not invoke aliyun CLI"
else
  pass "upload client has no aliyun CLI dependency"
fi
for required in \
  "bootstrap/runtime-config.sh" \
  "jarvis_load_runtime_config" \
  "JARVIS_CURL_BIN" \
  "/api/reports/aone/"; do
  if grep -Fq -- "$required" "$script"; then
    pass "script includes $required"
  else
    fail "script missing $required"
  fi
done

printf '\nPASS: %d  FAIL: %d\n' "$PASS" "$FAIL"
[ "$FAIL" -eq 0 ]
