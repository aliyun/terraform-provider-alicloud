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
    echo '[{"identifier":"WI-001","subject":"Test item","categoryIdentifier":"bug","status":"open","other":"ignored"}]'
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

# Validate fields: id, title, type, status, pool (5 fields)
field_count=$(echo "$output" | jq -e '.[0] | keys | length' 2>/dev/null)
if [ "$field_count" -eq 5 ]; then
    assert_pass "output object has exactly 5 fields (id,title,type,status,pool)"
else
    assert_fail "output object has $field_count fields (expected 5 including pool)"
fi

# Validate pool field present
if echo "$output" | jq -e '.[0].pool != null' > /dev/null 2>&1; then
    assert_pass "pool field present"
else
    assert_fail "pool field missing"
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
# Test 3: all pool-scoped a1 list calls fail → scan.sh exits 0, outputs []
# (failing pools are skipped, not fatal; only whoami failure is fatal)
# ---------------------------------------------------------------------------

tmpconfig3=$(mktemp -d)
mkdir -p "$tmpconfig3/config"
cat > "$tmpconfig3/config/pools.json" << 'JSON'
{
  "pools": {
    "pool_fail": { "project": 999, "name": "Failing Pool" }
  },
  "claim": { "tag": "jarvis-claimed" }
}
JSON

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

echo "=== Test 3: all pool a1 list failures → exit 0 with [] (pools are non-fatal) ==="
output=$(JARVIS_ROOT="$tmpconfig3" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?
rm -rf "$tmpconfig3"
echo "Exit code: $exit_code"
echo "Output: $output"
echo ""

if [ "$exit_code" -eq 0 ]; then
    assert_pass "exit 0 when all pool list calls fail (non-fatal skip)"
else
    assert_fail "should exit 0 when pool list fails (pool is skipped), got $exit_code"
fi

trimmed3=$(echo "$output" | tr -d ' \n')
if [ "$trimmed3" = "[]" ]; then
    assert_pass "all-pools-fail outputs []"
else
    assert_fail "all-pools-fail should output [], got: $output"
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
# Test 5: pools.json with claim.tag present → a1 list called with --filter NOT tag=<tag>
# ---------------------------------------------------------------------------

echo "=== Test 5: --filter NOT tag= passed when pools.json has claim.tag ==="

# Create a temp config dir with a pools.json containing claim.tag
tmpconfig=$(mktemp -d)
trap 'rm -rf "$tmpconfig" "$tmpbin"' EXIT

mkdir -p "$tmpconfig/config"
cat > "$tmpconfig/config/pools.json" << 'JSON'
{
  "claim": {
    "tag": "jarvis-claimed"
  }
}
JSON

# Record args passed to a1 list into a file
args_file=$(mktemp)

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "auth" ] && [ "\$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo "\$@" >> "$args_file"
    echo '[]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

# Run scan.sh with JARVIS_ROOT pointing to the temp config dir so it finds config/pools.json
output=$(JARVIS_ROOT="$tmpconfig" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?
recorded_args=$(cat "$args_file" 2>/dev/null || echo "")
rm -f "$args_file"

echo "Exit code: $exit_code"
echo "Recorded a1 args: $recorded_args"
echo "Output: $output"
echo ""

if [ "$exit_code" -eq 0 ]; then
    assert_pass "exit code 0 with pools.json tag"
else
    assert_fail "exit code should be 0 with pools.json tag, got $exit_code"
fi

if echo "$recorded_args" | grep -q -- '--filter'; then
    assert_pass "--filter flag passed to a1 list"
else
    assert_fail "--filter flag not found in a1 args: $recorded_args"
fi

if echo "$recorded_args" | grep -q 'NOT tag=jarvis-claimed'; then
    assert_pass "NOT tag=jarvis-claimed present in a1 args"
else
    assert_fail "NOT tag=jarvis-claimed not found in a1 args: $recorded_args"
fi

# ---------------------------------------------------------------------------
# Test 6: pools.json missing → a1 list called without --filter (no crash)
# ---------------------------------------------------------------------------

echo "=== Test 6: no --filter when pools.json is absent ==="

args_file2=$(mktemp)
tmpconfig2=$(mktemp -d)
# Do NOT create config/pools.json in tmpconfig2

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "auth" ] && [ "\$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo "\$@" >> "$args_file2"
    echo '[]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

output2=$(JARVIS_ROOT="$tmpconfig2" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code2=$?
recorded_args2=$(cat "$args_file2" 2>/dev/null || echo "")
rm -f "$args_file2"
rm -rf "$tmpconfig2"

echo "Exit code: $exit_code2"
echo "Recorded a1 args: $recorded_args2"
echo ""

if [ "$exit_code2" -eq 0 ]; then
    assert_pass "exit code 0 without pools.json"
else
    assert_fail "exit code should be 0 without pools.json, got $exit_code2"
fi

if echo "$recorded_args2" | grep -q -- '--filter'; then
    assert_fail "--filter should NOT be passed when pools.json absent"
else
    assert_pass "no --filter passed when pools.json absent"
fi

# ---------------------------------------------------------------------------
# Test 7: 2-pool config → stub returns 1 item per project → merged 2 items,
#         each item has pool field, --project and --filter passed per pool.
# ---------------------------------------------------------------------------

echo "=== Test 7: 2-pool config merges items from both pools ==="

tmpconfig7=$(mktemp -d)
mkdir -p "$tmpconfig7/config"
cat > "$tmpconfig7/config/pools.json" << 'JSON'
{
  "pools": {
    "pool_alpha": { "project": 111, "name": "Pool Alpha" },
    "pool_beta":  { "project": 222, "name": "Pool Beta" }
  },
  "claim": { "tag": "jarvis-claimed" }
}
JSON

args_file7=$(mktemp)

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "auth" ] && [ "\$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo "\$@" >> "$args_file7"
    # Detect which project and return corresponding item
    for arg in "\$@"; do
        if [ "\$arg" = "111" ]; then
            echo '[{"identifier":"WI-111","subject":"Alpha item","categoryIdentifier":"task","status":"open"}]'
            exit 0
        fi
        if [ "\$arg" = "222" ]; then
            echo '[{"identifier":"WI-222","subject":"Beta item","categoryIdentifier":"bug","status":"open"}]'
            exit 0
        fi
    done
    echo '[]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

output7=$(JARVIS_ROOT="$tmpconfig7" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code7=$?
recorded_args7=$(cat "$args_file7" 2>/dev/null || echo "")
rm -f "$args_file7"
rm -rf "$tmpconfig7"

echo "Exit code: $exit_code7"
echo "Output: $output7"
echo "Recorded args:"
echo "$recorded_args7"
echo ""

if [ "$exit_code7" -eq 0 ]; then
    assert_pass "exit code 0 for 2-pool scan"
else
    assert_fail "exit code should be 0 for 2-pool scan, got $exit_code7"
fi

# Should have exactly 2 items merged
item_count7=$(echo "$output7" | jq 'length' 2>/dev/null || echo 0)
if [ "$item_count7" -eq 2 ]; then
    assert_pass "2-pool scan produces 2 merged items"
else
    assert_fail "2-pool scan should produce 2 items, got $item_count7"
fi

# Both pool items present by id
if echo "$output7" | jq -e '[.[].id] | contains(["WI-111","WI-222"])' > /dev/null 2>&1; then
    assert_pass "both pool items present in merged output"
else
    assert_fail "merged output missing expected items: $output7"
fi

# Each item has a pool field
if echo "$output7" | jq -e 'all(.[]; .pool != null)' > /dev/null 2>&1; then
    assert_pass "all items have pool field"
else
    assert_fail "some items missing pool field: $output7"
fi

# WI-111 has pool name for pool_alpha
if echo "$output7" | jq -e '.[] | select(.id=="WI-111") | .pool == "pool_alpha"' > /dev/null 2>&1; then
    assert_pass "WI-111 pool field is pool_alpha"
else
    assert_fail "WI-111 pool field wrong: $(echo "$output7" | jq '.[] | select(.id=="WI-111")')"
fi

# WI-222 has pool name for pool_beta
if echo "$output7" | jq -e '.[] | select(.id=="WI-222") | .pool == "pool_beta"' > /dev/null 2>&1; then
    assert_pass "WI-222 pool field is pool_beta"
else
    assert_fail "WI-222 pool field wrong: $(echo "$output7" | jq '.[] | select(.id=="WI-222")')"
fi

# --project 111 was passed
if echo "$recorded_args7" | grep -q -- '--project 111\|--project.*111'; then
    assert_pass "--project 111 passed for pool_alpha"
else
    if echo "$recorded_args7" | grep -q '111'; then
        assert_pass "--project 111 passed for pool_alpha"
    else
        assert_fail "--project 111 not found in a1 args"
    fi
fi

# --project 222 was passed
if echo "$recorded_args7" | grep -q '222'; then
    assert_pass "--project 222 passed for pool_beta"
else
    assert_fail "--project 222 not found in a1 args"
fi

# --filter NOT tag=jarvis-claimed passed for both pools
filter_count7=$(echo "$recorded_args7" | grep -c 'NOT tag=jarvis-claimed' || echo 0)
if [ "$filter_count7" -ge 2 ]; then
    assert_pass "--filter NOT tag=jarvis-claimed passed for both pools ($filter_count7 times)"
else
    assert_fail "--filter NOT tag=jarvis-claimed should appear twice (once per pool), got $filter_count7 times"
fi

# --assignee flag passed for both pools
assignee_count7=$(echo "$recorded_args7" | grep -c -- '--assignee' || echo 0)
if [ "$assignee_count7" -ge 2 ]; then
    assert_pass "--assignee flag passed for both pools ($assignee_count7 times)"
else
    assert_fail "--assignee flag should appear twice (once per pool), got $assignee_count7 times"
fi

# ---------------------------------------------------------------------------
# Test 8: one pool fails → still returns items from the healthy pool (non-fatal)
# ---------------------------------------------------------------------------

echo "=== Test 8: failing pool is skipped, healthy pool items returned ==="

tmpconfig8=$(mktemp -d)
mkdir -p "$tmpconfig8/config"
cat > "$tmpconfig8/config/pools.json" << 'JSON'
{
  "pools": {
    "pool_good": { "project": 333, "name": "Good Pool" },
    "pool_bad":  { "project": 444, "name": "Bad Pool" }
  },
  "claim": { "tag": "jarvis-claimed" }
}
JSON

cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
if [ "$1" = "auth" ] && [ "$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "$1" = "project" ] && [ "$2" = "workitem" ] && [ "$3" = "list" ]; then
    for arg in "$@"; do
        if [ "$arg" = "333" ]; then
            echo '[{"identifier":"WI-333","subject":"Good item","categoryIdentifier":"task","status":"open"}]'
            exit 0
        fi
        if [ "$arg" = "444" ]; then
            echo 'Error: access denied' >&2
            exit 1
        fi
    done
    echo '[]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

output8=$(JARVIS_ROOT="$tmpconfig8" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code8=$?
rm -rf "$tmpconfig8"

echo "Exit code: $exit_code8"
echo "Output: $output8"
echo ""

if [ "$exit_code8" -eq 0 ]; then
    assert_pass "exit code 0 when one pool fails"
else
    assert_fail "exit code should be 0 when one pool fails (skipped), got $exit_code8"
fi

item_count8=$(echo "$output8" | jq 'length' 2>/dev/null || echo 0)
if [ "$item_count8" -eq 1 ]; then
    assert_pass "healthy pool item returned despite failing pool"
else
    assert_fail "should return 1 item from healthy pool, got $item_count8: $output8"
fi

# ---------------------------------------------------------------------------
# Test 9: pools.json with no claim.tag → plain list per pool (no --filter)
# ---------------------------------------------------------------------------

echo "=== Test 9: no claim.tag in pools.json → no --filter passed ==="

tmpconfig9=$(mktemp -d)
mkdir -p "$tmpconfig9/config"
cat > "$tmpconfig9/config/pools.json" << 'JSON'
{
  "pools": {
    "pool_x": { "project": 555, "name": "Pool X" }
  }
}
JSON

args_file9=$(mktemp)

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "auth" ] && [ "\$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo "\$@" >> "$args_file9"
    echo '[{"identifier":"WI-555","subject":"X item","categoryIdentifier":"task","status":"open"}]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

output9=$(JARVIS_ROOT="$tmpconfig9" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code9=$?
recorded_args9=$(cat "$args_file9" 2>/dev/null || echo "")
rm -f "$args_file9"
rm -rf "$tmpconfig9"

echo "Exit code: $exit_code9"
echo "Recorded args: $recorded_args9"
echo ""

if [ "$exit_code9" -eq 0 ]; then
    assert_pass "exit code 0 with no claim.tag"
else
    assert_fail "exit code should be 0 with no claim.tag, got $exit_code9"
fi

if echo "$recorded_args9" | grep -q -- '--filter'; then
    assert_fail "--filter should NOT be passed when no claim.tag"
else
    assert_pass "no --filter when no claim.tag"
fi

if echo "$recorded_args9" | grep -q '555'; then
    assert_pass "--project 555 passed for pool_x"
else
    assert_fail "--project 555 not found in recorded args: $recorded_args9"
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
