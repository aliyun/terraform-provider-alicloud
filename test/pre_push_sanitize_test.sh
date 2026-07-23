#!/usr/bin/env bash
# test/pre_push_sanitize_test.sh - regression tests for public push sanitization.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"
sanitizer="$proj_root/bootstrap/pre-push-sanitize.sh"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

PASS=0
FAIL=0

pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

new_repo() {
    local repo="$1"
    mkdir -p "$repo"
    git -C "$repo" init -q -b master
    git -C "$repo" config user.name "Sanitizer Test"
    git -C "$repo" config user.email "sanitizer@example.invalid"
    git -C "$repo" config core.hooksPath /dev/null
}

commit_all() {
    local repo="$1"
    local message="$2"
    git -C "$repo" add -A
    git -C "$repo" commit -q -m "$message"
}

run_sanitizer() {
    local repo="$1"
    (
        cd "$repo"
        bash "$sanitizer" base
    ) 2>&1
}

assert_exit() {
    local expected="$1"
    local actual="$2"
    local label="$3"
    if [ "$actual" -eq "$expected" ]; then
        pass "$label"
    else
        fail "$label (expected exit $expected, got $actual)"
    fi
}

echo "=== Test 1: deleted and context-only generic RequestId do not block ==="
repo="$tmpdir/context"
new_repo "$repo"
cat > "$repo/example.md" <<'EOF'
OpenAPI RequestId is a field name.
Keep this RequestId documentation as context.
old value
EOF
commit_all "$repo" "base"
git -C "$repo" tag base
cat > "$repo/example.md" <<'EOF'
Keep this RequestId documentation as context.
new value
EOF
commit_all "$repo" "update nearby prose"
output="$(run_sanitizer "$repo")"
code=$?
echo "$output"
assert_exit 0 "$code" "deleted and context-only generic field names pass"

echo "=== Test 2: newly added generic RequestId without a value does not block ==="
repo="$tmpdir/generic"
new_repo "$repo"
printf 'base\n' > "$repo/example.md"
commit_all "$repo" "base"
git -C "$repo" tag base
printf 'base\nOpenAPI RequestId is recorded for diagnostics.\n' > "$repo/example.md"
commit_all "$repo" "document the API field"
output="$(run_sanitizer "$repo")"
code=$?
echo "$output"
assert_exit 0 "$code" "new generic field name passes"

echo "=== Test 3: newly added concrete RequestId value blocks ==="
repo="$tmpdir/request-id-value"
new_repo "$repo"
printf 'base\n' > "$repo/example.json"
commit_all "$repo" "base"
git -C "$repo" tag base
request_key="Request""Id"
request_value="01234567-89ab-cdef-0123-456789abcdef"
printf '{"%s":"%s"}\n' "$request_key" "$request_value" > "$repo/example.json"
commit_all "$repo" "add response fixture"
output="$(run_sanitizer "$repo")"
code=$?
echo "$output"
assert_exit 1 "$code" "concrete RequestId value is blocked"
if grep -Fq "RequestId 值" <<<"$output"; then
    pass "RequestId value reports the expected category"
else
    fail "RequestId value should report its category"
fi

echo "=== Test 4: commit messages remain fully scanned ==="
repo="$tmpdir/message"
new_repo "$repo"
printf 'base\n' > "$repo/example.md"
commit_all "$repo" "base"
git -C "$repo" tag base
printf 'safe change\n' > "$repo/example.md"
internal_host="project.aone.""alibaba-inc.com"
commit_all "$repo" "See $internal_host before release"
output="$(run_sanitizer "$repo")"
code=$?
echo "$output"
assert_exit 1 "$code" "sensitive commit message is blocked"
if grep -Fq "Aone URL" <<<"$output"; then
    pass "commit-message violation reports the expected category"
else
    fail "commit-message violation should report its category"
fi

echo "=== Test 5: patch metadata and rename-only changes are not content ==="
repo="$tmpdir/rename"
new_repo "$repo"
printf 'safe\n' > "$repo/RequestId-old-name.md"
commit_all "$repo" "base"
git -C "$repo" tag base
git -C "$repo" mv RequestId-old-name.md RequestId-new-name.md
commit_all "$repo" "rename documentation"
output="$(run_sanitizer "$repo")"
code=$?
echo "$output"
assert_exit 0 "$code" "rename metadata is excluded from content scanning"

echo "=== Test 6: other newly added public-artifact violations still block ==="
repo="$tmpdir/public-violations"
new_repo "$repo"
printf 'base\n' > "$repo/example.txt"
commit_all "$repo" "base"
git -C "$repo" tag base
{
    printf '%s\n' "$internal_host/workitem/FAKE"
    printf '#%s\n' "1234567"
    printf 'r-%s\n' "0123456789ab"
    printf 'i-%s\n' "0123456789abcdefgh"
    printf 'lb-%s\n' "0123456789ab"
    printf 's-%s\n' "0123456789ab"
    printf 'Co%sAuthored-By: %s\n' "-" "Codex"
    printf '@%s(%s)\n' "Example" "ABCD1234"
} > "$repo/example.txt"
commit_all "$repo" "add unsafe public content"
output="$(run_sanitizer "$repo")"
code=$?
echo "$output"
assert_exit 1 "$code" "newly added sensitive content is blocked"
for category in \
    "Aone URL" \
    "7 位以上工单号" \
    "Redis/RDS 实例 ID" \
    "ECS 实例 ID" \
    "SLB 实例 ID" \
    "OSS 实例 ID (s-)" \
    "AI 署名 (Co-Authored-By)" \
    "花名+工号引用"; do
    if grep -Fq "命中: $category" <<<"$output"; then
        pass "sensitive additions retain category: $category"
    else
        fail "missing sensitive category: $category"
    fi
done

echo "=== Test 7: binary changes are explicitly handed to manual review ==="
repo="$tmpdir/binary"
new_repo "$repo"
printf 'base\n' > "$repo/example.md"
commit_all "$repo" "base"
git -C "$repo" tag base
printf '\000\001\002\003' > "$repo/artifact.bin"
commit_all "$repo" "add binary artifact"
output="$(run_sanitizer "$repo")"
code=$?
echo "$output"
assert_exit 0 "$code" "binary-only change does not scan patch metadata"
if grep -Fq "二进制变更" <<<"$output" && grep -Fq "artifact.bin" <<<"$output"; then
    pass "binary-only change emits explicit manual-review warning"
else
    fail "binary-only change should emit a manual-review warning"
fi

echo ""
echo "=== Summary ==="
echo "PASS: $PASS  FAIL: $FAIL"
if [ "$FAIL" -ne 0 ]; then
    exit 1
fi
echo "All tests passed"
