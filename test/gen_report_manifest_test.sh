#!/usr/bin/env bash
# test/gen_report_manifest_test.sh - gen-report.py --manifest bridges
# evidence-manifest.md local paths to signed URLs, and fails closed on gaps.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"
gen="$proj_root/.claude/skills/html-report-preview/scripts/gen-report.py"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

PASS=0
FAIL=0
assert_pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
assert_fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

assert_contains() {
    if grep -Fq -- "$2" <<<"$1"; then
        assert_pass "$3"
    else
        assert_fail "$3 (missing: $2; got: $1)"
    fi
}

shots="$tmpdir/screenshots"
mkdir -p "$shots"
: >"$shots/openapi_create_rule.png"
: >"$shots/provider_no_resource.png"

# Real manifest shape: absolute local paths in the `screenshot` column.
manifest="$tmpdir/evidence-manifest.md"
cat >"$manifest" <<MD
# Visual Evidence Manifest — Aone 74939817

## 三层截图

| layer | result | screenshot | source | note |
|---|---|---|---|---|
| OpenAPI | pass | $shots/openapi_create_rule.png | https://next.api.example/doc | CRUD 齐全 |
| CloudSpec/ACube | n-a | N/A | https://acube.example/q | missing_capability: 无浏览器通道 |
| Provider | fail | $shots/provider_no_resource.png | https://github.example/x.go | 零命中 |
MD

urls="$tmpdir/image-urls.txt"
cat >"$urls" <<'TXT'
openapi_create_rule|https://img.example/a.png?Expires=1&Signature=x
provider_no_resource|https://img.example/b.png?Expires=1&Signature=y
TXT

echo "=== Test 1: manifest + image-urls embeds signed URLs ==="
out="$tmpdir/report.html"
if python3 "$gen" --title "T" --manifest "$manifest" --image-urls "$urls" >"$out" 2>"$tmpdir/err"; then
    assert_pass "manifest mode succeeds"
else
    assert_fail "manifest mode should succeed ($(cat "$tmpdir/err"))"
fi
report="$(cat "$out")"
assert_contains "$report" 'src="https://img.example/a.png' "OpenAPI screenshot is embedded"
assert_contains "$report" 'src="https://img.example/b.png' "Provider screenshot is embedded"
assert_contains "$report" "CRUD 齐全" "manifest note is carried over"
assert_contains "$report" "https://next.api.example/doc" "manifest source is linked"
img_count="$(grep -c '<img' "$out")"
if [ "$img_count" -eq 2 ]; then
    assert_pass "exactly the 2 captured layers are embedded (n-a stays N/A)"
else
    assert_fail "expected 2 <img>, got $img_count"
fi

echo "=== Test 2: the generated report passes the upload guard ==="
stub="$tmpdir/bin"
mkdir -p "$stub"
printf '#!/usr/bin/env bash\nexit 22\n' >"$stub/curl"
chmod +x "$stub/curl"
guard_out="$(JARVIS_RUNTIME_CONFIG_LOADED=1 JARVIS_CURL_BIN="$stub/curl" \
    JARVIS_HTML_REPORT_TOKEN=t bash "$proj_root/bootstrap/html-report-preview.sh" \
    upload FAKE "$out" --base-url https://pre.example 2>&1)"
if grep -Fq 'screenshot_prose_without_img' <<<"$guard_out"; then
    assert_fail "generated report must not trip the prose-without-img guard"
else
    assert_pass "generated report clears the prose-without-img guard"
fi

echo "=== Test 3: fail closed when a screenshot has no signed URL ==="
partial="$tmpdir/partial-urls.txt"
printf 'openapi_create_rule|https://img.example/a.png\n' >"$partial"
if python3 "$gen" --title "T" --manifest "$manifest" --image-urls "$partial" \
    >"$tmpdir/bad.html" 2>"$tmpdir/bad.err"; then
    assert_fail "missing signed URL must not produce a report"
else
    assert_pass "missing signed URL fails closed"
fi
assert_contains "$(cat "$tmpdir/bad.err")" "provider_no_resource" \
    "error names the unresolved screenshot"
if [ -s "$tmpdir/bad.html" ]; then
    assert_fail "no partial report may be written"
else
    assert_pass "no partial report is written"
fi

echo "=== Test 4: uploading requires an Aone id ==="
if python3 "$gen" --title "T" --manifest "$manifest" >"$tmpdir/noid.html" 2>"$tmpdir/noid.err"; then
    assert_fail "manifest with local shots must not silently drop them"
else
    assert_pass "manifest without --image-urls/--aone-id fails closed"
fi
assert_contains "$(cat "$tmpdir/noid.err")" "--aone-id" "error explains how to proceed"

echo "=== Test 5: manifest with no screenshots at all still renders ==="
none="$tmpdir/none.md"
cat >"$none" <<'MD'
| layer | result | screenshot | source | note |
|---|---|---|---|---|
| OpenAPI | n-a | N/A | https://x.example | missing_capability: 无通道 |
MD
if python3 "$gen" --title "T" --manifest "$none" >"$tmpdir/none.html" 2>"$tmpdir/none.err"; then
    assert_pass "screenshot-free manifest renders a degraded report"
else
    assert_fail "screenshot-free manifest should render ($(cat "$tmpdir/none.err"))"
fi

echo "=== Test 6: .claude and .agents copies stay in mirror sync ==="
# shellcheck source=../bootstrap/skills-mirror-lib.sh
source "$proj_root/bootstrap/skills-mirror-lib.sh"
mirror_sed_codex_to_claude \
    <"$proj_root/.agents/skills/html-report-preview/scripts/gen-report.py" \
    >"$tmpdir/mirrored.py"
if diff -u "$tmpdir/mirrored.py" "$gen" >"$tmpdir/mirror.diff"; then
    assert_pass "gen-report.py mirror is in sync"
else
    assert_fail "gen-report.py mirror drift: $(head -6 "$tmpdir/mirror.diff")"
fi

echo "=== Test 7: uploader resolves without a hardcoded skills root ==="
uploader_line="$(grep -n 'screenshot-evidence' "$gen" | head -1)"
if grep -q 'parents\[2\]' "$gen"; then
    assert_pass "uploader is resolved relative to this script"
else
    assert_fail "uploader path should not depend on the skills root ($uploader_line)"
fi

echo ""
echo "=== Summary ==="
echo "PASS: $PASS  FAIL: $FAIL"
[ "$FAIL" -gt 0 ] && { echo "TESTS FAILED"; exit 1; }
echo "All tests passed"
exit 0
