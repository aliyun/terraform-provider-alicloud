#!/usr/bin/env bash
# test/aone_image_extract_test.sh — tests for bootstrap/aone-image-extract.sh
#
# 覆盖:
#   T1 无附件工单 → 输出 images=0
#   T2 混合附件(图 + 非图) → 只下图,manifest 列全,skipped_non_image 计数正确
#   T3 二次运行 + summary.md 更新 → cache=true 且回显 summary 内容
#   T4 a1 download 失败(退非零) → download_failed 计数,不阻断

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"
tmpdir="$(mktemp -d)"
tmpbin="$tmpdir/bin"
a1_log="$tmpdir/a1.log"
mkdir -p "$tmpbin"
trap 'rm -rf "$tmpdir"' EXIT

PASS=0; FAIL=0
assert_pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
assert_fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }
assert_contains() {
    local hay="$1" needle="$2" label="$3"
    if grep -Fq -- "$needle" <<<"$hay"; then assert_pass "$label"
    else assert_fail "$label (missing: $needle; got: $hay)"; fi
}
assert_not_contains() {
    local hay="$1" needle="$2" label="$3"
    if grep -Fq -- "$needle" <<<"$hay"; then assert_fail "$label (unexpectedly found: $needle)"
    else assert_pass "$label"; fi
}

# stub a1id: default variant returns empty attachment list; variants set via A1_MODE env
make_stub_a1() {
    cat > "$tmpbin/a1id" <<'STUB'
#!/usr/bin/env bash
[ "${1:-}" = "--" ] && shift
echo "$@" >> "$A1_LOG"
sub="$1 $2 $3 $4"
case "$A1_MODE:$sub" in
  empty:"project workitem attachment list")
    echo "No attachments found"; exit 0 ;;
  empty_json:"project workitem attachment list")
    echo "[]"; exit 0 ;;
  mixed:"project workitem attachment list")
    cat <<'JSON'
[
  {"id":501,"name":"error-screen.png","created":"2026-07-01T10:00:00+08:00"},
  {"id":502,"name":"stack-trace.log","created":"2026-07-01T10:05:00+08:00"},
  {"id":503,"name":"console screenshot.jpg","created":"2026-07-01T10:10:00+08:00"},
  {"id":504,"name":"api-payload.json","created":"2026-07-01T10:15:00+08:00"},
  {"id":505,"name":"cli.JPEG","created":"2026-07-01T10:20:00+08:00"}
]
JSON
    exit 0 ;;
  failed:"project workitem attachment list")
    cat <<'JSON'
[
  {"id":601,"name":"broken.png","created":"2026-07-01T10:00:00+08:00"}
]
JSON
    exit 0 ;;
esac
if [ "$sub" = "project workitem attachment download" ]; then
    out=""
    while [ "$#" -gt 0 ]; do
        [ "${1:-}" = "-o" ] && { out="${2:-}"; break; }
        shift
    done
    if [ "$A1_MODE" = "failed" ]; then
        echo "simulated download failure" >&2; exit 33
    fi
    printf 'FAKE_IMG_BYTES' > "$out"
    exit 0
fi
echo "unexpected a1id call: $*" >&2; exit 2
STUB
    chmod +x "$tmpbin/a1id"
}

run_script() {
    local id="$1" mode="$2"
    A1_MODE="$mode" A1_LOG="$a1_log" \
        JARVIS_A1="$tmpbin/a1id --" \
        JARVIS_IMAGE_OCR_DIR="$tmpdir/ocr" \
        bash "$proj_root/bootstrap/aone-image-extract.sh" "$id" 2>&1
}

make_stub_a1

echo "=== T1: no attachments → images=0 ==="
: > "$a1_log"
out=$(run_script 99991 empty)
assert_contains "$out" "id=99991" "T1: id echoed"
assert_contains "$out" "images=0" "T1: images=0"
assert_not_contains "$out" "## Local paths" "T1: no Local paths section"

echo "=== T1b: empty JSON array (a1 shape variant) → images=0 ==="
: > "$a1_log"
out=$(run_script 99991 empty_json)
assert_contains "$out" "images=0" "T1b: images=0 on empty array"

echo "=== T2: mixed attachments → only images downloaded, correct manifest ==="
: > "$a1_log"
out=$(run_script 12345 mixed)
assert_contains "$out" "id=12345" "T2: id echoed"
assert_contains "$out" "images=3" "T2: images=3 (png, jpg, JPEG)"
assert_contains "$out" "downloaded=3" "T2: downloaded=3"
assert_contains "$out" "skipped_non_image=2" "T2: skipped=2 (log + json)"
assert_contains "$out" "att_id=501" "T2: png in manifest"
assert_contains "$out" "att_id=503" "T2: jpg in manifest"
assert_contains "$out" "att_id=505" "T2: JPEG in manifest"
assert_not_contains "$out" "att_id=502" "T2: log skipped from manifest"
assert_not_contains "$out" "att_id=504" "T2: json skipped from manifest"
assert_contains "$out" "cache=false" "T2: cache=false on first run"
assert_contains "$out" "## Next step" "T2: next-step block for skill agent"

# check downloads landed and space in filename got sanitized
for f in "$tmpdir/ocr/12345/501-error-screen.png" "$tmpdir/ocr/12345/503-console_screenshot.jpg" "$tmpdir/ocr/12345/505-cli.JPEG"; do
    if [ -s "$f" ]; then assert_pass "T2: downloaded $f"; else assert_fail "T2: missing $f"; fi
done

echo "=== T2b: re-run, files already local → no re-download ==="
: > "$a1_log"
out=$(run_script 12345 mixed)
assert_contains "$out" "downloaded=0" "T2b: no re-download when files already present"
assert_contains "$out" "images=3" "T2b: manifest still lists 3 images"

echo "=== T3: cached summary.md → cache=true, replay ==="
mkdir -p "$tmpdir/ocr/12345"
cat > "$tmpdir/ocr/12345/summary.md" <<'SUM'
### 图片 1 (error-screen.png)
- 错误消息: ErrorCode=InvalidParam Message="foo bar"
SUM
: > "$a1_log"
out=$(run_script 12345 mixed)
assert_contains "$out" "cache=true" "T3: cache=true when summary.md present"
assert_contains "$out" "## Cached summary" "T3: cached summary block emitted"
assert_contains "$out" "ErrorCode=InvalidParam" "T3: cached summary content replayed"
assert_not_contains "$out" "## Next step" "T3: next-step suppressed when cached"

echo "=== T4: download failure → download_failed counted, others still emit ==="
: > "$a1_log"
out=$(run_script 22222 failed)
assert_contains "$out" "download_failed=1" "T4: failure counted"
assert_contains "$out" "images=0" "T4: no successful image in manifest"
# stale zero-byte file must be cleaned up
if [ ! -e "$tmpdir/ocr/22222/601-broken.png" ]; then
    assert_pass "T4: no partial file left behind"
else
    assert_fail "T4: partial file leaked at 22222/601-broken.png"
fi

echo ""
echo "=== SUMMARY ==="
echo "PASS: $PASS"
echo "FAIL: $FAIL"
[ "$FAIL" -eq 0 ]
