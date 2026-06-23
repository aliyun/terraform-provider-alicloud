#!/bin/bash

# Test harness for scan.sh inbox scanner
# Stubs a1 to return fixed data, asserts JSON output shape

set -u

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

tmpbin=$(mktemp -d)
trap 'rm -rf "$tmpbin"' EXIT

PASS=0
FAIL=0

assert_pass() {
    echo "PASS: $1"
    PASS=$((PASS + 1))
}

assert_fail() {
    echo "FAIL: $1"
    FAIL=$((FAIL + 1))
}

# ---------------------------------------------------------------------------
# Test 1: stub a1 returning 1-item list → output is JSON array containing id
# ---------------------------------------------------------------------------

cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
# Stub a1: handles whoami and workitem list
if [ "$1" = "auth" ] && [ "$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    echo "Name:     Chen Han Zhang"
    exit 0
fi
if [ "$1" = "project" ] && [ "$2" = "workitem" ] && [ "$3" = "list" ]; then
    echo '[{"identifier":"WI-001","title":"Test item","categoryIdentifier":"bug","status":"open","other":"ignored"}]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

export PATH="$tmpbin:$PATH"

output=$(bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?

echo "=== Test 1: single-item list ==="
echo "Output: $output"
echo "Exit code: $exit_code"
echo ""

if [ "$exit_code" -eq 0 ]; then
    assert_pass "exit code 0"
else
    assert_fail "exit code should be 0, got $exit_code"
fi

# Validate JSON array
if echo "$output" | jq -e 'type == "array"' > /dev/null 2>&1; then
    assert_pass "output is JSON array"
else
    assert_fail "output is not a JSON array"
fi

# Validate id field present and correct
if echo "$output" | jq -e '.[0].id == "WI-001"' > /dev/null 2>&1; then
    assert_pass "id field is WI-001"
else
    assert_fail "id field missing or wrong (expected WI-001)"
fi

# Validate title field
if echo "$output" | jq -e '.[0].title == "Test item"' > /dev/null 2>&1; then
    assert_pass "title field correct"
else
    assert_fail "title field missing or wrong"
fi

# Validate type field mapped from categoryIdentifier
if echo "$output" | jq -e '.[0].type == "bug"' > /dev/null 2>&1; then
    assert_pass "type field mapped from categoryIdentifier"
else
    assert_fail "type field missing or wrong"
fi

# Validate status field
if echo "$output" | jq -e '.[0].status == "open"' > /dev/null 2>&1; then
    assert_pass "status field correct"
else
    assert_fail "status field missing or wrong"
fi

# Validate no extra fields (only id, title, type, status)
field_count=$(echo "$output" | jq -e '.[0] | keys | length' 2>/dev/null)
if [ "$field_count" -eq 4 ]; then
    assert_pass "output object has exactly 4 fields"
else
    assert_fail "output object has $field_count fields (expected 4)"
fi

# ---------------------------------------------------------------------------
# Test 2: stub a1 returning empty list → output is [] and exit 0
# ---------------------------------------------------------------------------

cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
if [ "$1" = "auth" ] && [ "$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "$1" = "project" ] && [ "$2" = "workitem" ] && [ "$3" = "list" ]; then
    echo '[]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

echo "=== Test 2: empty inbox ==="
output=$(bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?
echo "Output: $output"
echo "Exit code: $exit_code"
echo ""

if [ "$exit_code" -eq 0 ]; then
    assert_pass "exit code 0 for empty inbox"
else
    assert_fail "exit code should be 0 for empty inbox, got $exit_code"
fi

trimmed=$(echo "$output" | tr -d ' \n')
if [ "$trimmed" = "[]" ]; then
    assert_pass "empty inbox outputs []"
else
    assert_fail "empty inbox should output [], got: $output"
fi

# ---------------------------------------------------------------------------
# Test 3: a1 list fails → scan.sh exits non-zero
# ---------------------------------------------------------------------------

cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
if [ "$1" = "auth" ] && [ "$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "$1" = "project" ] && [ "$2" = "workitem" ] && [ "$3" = "list" ]; then
    echo 'Error: not authenticated' >&2
    exit 1
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

echo "=== Test 3: a1 list failure exits non-zero ==="
output=$(bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?
echo "Exit code: $exit_code"
echo ""

if [ "$exit_code" -ne 0 ]; then
    assert_pass "non-zero exit on a1 failure"
else
    assert_fail "should exit non-zero when a1 fails, got 0"
fi

# ---------------------------------------------------------------------------
# Test 4: whoami stub exits 1 → scan.sh exits non-zero
# ---------------------------------------------------------------------------

cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
if [ "$1" = "auth" ] && [ "$2" = "whoami" ]; then
    echo "Error: authentication failed" >&2
    exit 1
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

echo "=== Test 4: whoami failure exits non-zero ==="
output=$(bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?
echo "Exit code: $exit_code"
echo ""

if [ "$exit_code" -ne 0 ]; then
    assert_pass "non-zero exit on whoami failure"
else
    assert_fail "should exit non-zero when whoami fails, got 0"
fi

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo "=== Summary ==="
echo "PASS: $PASS  FAIL: $FAIL"
echo ""

if [ "$FAIL" -gt 0 ]; then
    echo "TESTS FAILED"
    exit 1
fi

echo "All tests passed"

# ---------------------------------------------------------------------------
# Live run: print real inbox count (informational, non-fatal)
# ---------------------------------------------------------------------------
echo ""
echo "=== Live inbox (real a1) ==="
if command -v a1 > /dev/null 2>&1; then
    real_account=$(a1 auth whoami 2>/dev/null | awk '/Account:/{print $2}')
    if [ -n "$real_account" ]; then
        live_output=$(bash "$proj_root/bootstrap/scan.sh" 2>&1) || true
        count=$(echo "$live_output" | jq 'length' 2>/dev/null || echo "unknown")
        echo "Live inbox for $real_account: $count items"
    else
        echo "(could not determine account)"
    fi
else
    echo "(a1 not available in environment)"
fi

exit 0
