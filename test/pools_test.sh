#!/usr/bin/env bash
# test/pools_test.sh – per-pool inspection test for bootstrap/pools.sh
# Stubs a1 to return a 2-item array; uses a tmp pools.json with 2 pools.
# Asserts: two "active=2" lines appear, each with correct name+project; exit 0.
# Also asserts: one pool whose a1 call fails prints ERR and the script still exits 0.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

tmpbin=$(mktemp -d)
tmpconfig=$(mktemp -d)
trap 'rm -rf "$tmpbin" "$tmpconfig"' EXIT

PASS=0
FAIL=0

assert_pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
assert_fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# ---------------------------------------------------------------------------
# Build stub a1: always returns a 2-item JSON array
# ---------------------------------------------------------------------------
cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
if [ "$1" = "project" ] && [ "$2" = "workitem" ] && [ "$3" = "list" ]; then
    echo '[{"id":1},{"id":2}]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

# ---------------------------------------------------------------------------
# Build tmp pools.json with 2 pools
# ---------------------------------------------------------------------------
mkdir -p "$tmpconfig/config"
cat > "$tmpconfig/config/pools.json" << 'JSON'
{
  "pools": {
    "alpha": { "project": 1111111 },
    "beta":  { "project": 2222222 }
  }
}
JSON

export PATH="$tmpbin:$PATH"

# ---------------------------------------------------------------------------
# Test 1: happy path – 2 pools both succeed → 2 "active=2" lines + names + exit 0
# ---------------------------------------------------------------------------
echo "=== Test 1: two pools, both succeed ==="
output=$(JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/pools.sh" 2>&1)
exit_code=$?

echo "Output:"
echo "$output"
echo "Exit code: $exit_code"
echo ""

if [ "$exit_code" -eq 0 ]; then
    assert_pass "exit code 0"
else
    assert_fail "exit code should be 0, got $exit_code"
fi

if echo "$output" | grep -q "alpha.*1111111.*active=2 total=2"; then
    assert_pass "alpha pool line contains name, project, active=2"
else
    assert_fail "missing alpha pool line with active=2 (got: $output)"
fi

if echo "$output" | grep -q "beta.*2222222.*active=2 total=2"; then
    assert_pass "beta pool line contains name, project, active=2"
else
    assert_fail "missing beta pool line with active=2 (got: $output)"
fi

# Total line should show 4 (2 + 2)
if echo "$output" | grep -q "TOTAL active=4 total=4"; then
    assert_pass "TOTAL active=4 total=4 line present"
else
    assert_fail "TOTAL active=4 total=4 line missing (got: $output)"
fi

# ---------------------------------------------------------------------------
# Test 2: one pool's a1 call fails → prints ERR for that pool, still exits 0
# ---------------------------------------------------------------------------
echo "=== Test 2: one pool fails → ERR, script still exits 0 ==="

# Stub: project 1111111 succeeds, project 2222222 fails
cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
if [[ "$*" == *"--project 1111111"* ]]; then
    echo '[{"id":1},{"id":2}]'; exit 0
fi
echo "Error: network failure" >&2; exit 1
STUB
chmod +x "$tmpbin/a1"

output2=$(JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/pools.sh" 2>&1)
exit_code2=$?

echo "Output:"
echo "$output2"
echo "Exit code: $exit_code2"
echo ""

if [ "$exit_code2" -eq 0 ]; then
    assert_pass "exit code 0 when one pool fails"
else
    assert_fail "exit code should be 0 even when one pool fails, got $exit_code2"
fi

if echo "$output2" | grep -q "alpha.*1111111.*active=2 total=2"; then
    assert_pass "alpha pool still shows active=2"
else
    assert_fail "alpha pool active=2 missing when beta fails (got: $output2)"
fi

if echo "$output2" | grep -qE "beta.*2222222.*ERR"; then
    assert_pass "beta pool prints ERR"
else
    assert_fail "beta pool should print ERR (got: $output2)"
fi

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo "=== Summary ==="
echo "PASS: $PASS  FAIL: $FAIL"

if [ "$FAIL" -gt 0 ]; then
    echo "TESTS FAILED"
    exit 1
fi

echo "All tests passed"
exit 0
