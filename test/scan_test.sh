#!/bin/bash

# Test harness for scan.sh inbox scanner
# Stubs a1 to return fixed data, asserts JSON output shape

set -u

# Bypass scan.sh 30min TTL cache in tests (P1.d cache infra) — 避免用例串扰读旧 scan.json
export JARVIS_SCAN_TTL=0

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
    echo '[{"identifier":"WI-001","subject":"Test item","categoryIdentifier":"bug","status":"open","priority":"high","tag":["p0","urgent"],"other":"ignored"}]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

export PATH="$tmpbin:$PATH"
export JARVIS_A1=a1   # 走 PATH 上的 a1 stub,不经真 bin/a1id

# Test 1 走 no-pools fallback path:tmpdir 无 config/pools.json → global assignee list
tmpconfig1=$(mktemp -d)
output=$(JARVIS_ROOT="$tmpconfig1" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?
rm -rf "$tmpconfig1"

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

# Validate priority field
if echo "$output" | jq -e '.[0].priority == "high"' > /dev/null 2>&1; then
    assert_pass "priority field correct"
else
    assert_fail "priority field missing or wrong (expected \"high\")"
fi

# Validate tag field
if echo "$output" | jq -e '.[0].tag == ["p0","urgent"]' > /dev/null 2>&1; then
    assert_pass "tag field correct"
else
    assert_fail "tag field missing or wrong (expected [\"p0\",\"urgent\"])"
fi

# Validate fields: no-pools fallback outputs id,title,type,status,priority,tag,category,modified,created
# (9 字段,无 pool)。modified 由 F2 加入、created 由灰度范围限定加入(供 JARVIS_DISPATCH_CREATED_BEFORE 判定)。
field_count=$(echo "$output" | jq -e '.[0] | keys | length' 2>/dev/null)
if [ "$field_count" -eq 9 ]; then
    assert_pass "output object has exactly 9 fields (…,category,modified,created)"
else
    assert_fail "output object has $field_count fields (expected 9: id,title,type,status,priority,tag,category,modified,created)"
fi

# Validate category field present (no-pools fallback stamps category:null)
if echo "$output" | jq -e '.[0] | has("category")' > /dev/null 2>&1; then
    assert_pass "category field present (no-pools fallback stamps null)"
else
    assert_fail "category field missing in no-pools output"
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
# 独立 tmpdir:no-pools fallback,避免读主 repo pools.json 或 cache
tmpconfig2_empty=$(mktemp -d)
output=$(JARVIS_ROOT="$tmpconfig2_empty" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?
rm -rf "$tmpconfig2_empty"
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
output=$(JARVIS_ROOT="$tmpconfig3" bash "$proj_root/bootstrap/scan.sh" 2>/dev/null)  # stdout(JSON)only; WARN 走 stderr
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
# 独立 tmpdir:避免读主 repo pools.json;whoami fail → scan.sh 应 exit 1
tmpconfig4=$(mktemp -d)
output=$(JARVIS_ROOT="$tmpconfig4" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code=$?
rm -rf "$tmpconfig4"
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

# pools.json 只有 claim.tag 无 pools block → scan.sh 走 no-pools fallback,不生成 --filter
# (P1 refactor 后 scan.sh 也从 pool-scoped path 移除 claim-tag exclusion,dedup 下推到 triage)
if echo "$recorded_args" | grep -q -- '--filter'; then
    assert_fail "--filter should NOT be passed (no-pools fallback + P1 refactor): $recorded_args"
else
    assert_pass "no --filter passed (no-pools fallback + claim-tag exclusion removed)"
fi

if echo "$recorded_args" | grep -q 'NOT tag=jarvis-claimed'; then
    assert_fail "NOT tag=jarvis-claimed should NOT be in args (P1 refactor removed claim-tag exclusion): $recorded_args"
else
    assert_pass "no NOT tag=jarvis-claimed in a1 args (claim-tag exclusion removed)"
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

# scan.sh 展开 3 categories (req/bug/task) per pool;stub 无视 category 每次都返同一 item,
# 所以 2 pool × 3 category = 6 items(stub 简化;真实场景每 category 应返不同 items)
item_count7=$(echo "$output7" | jq 'length' 2>/dev/null || echo 0)
if [ "$item_count7" -eq 6 ]; then
    assert_pass "2-pool × 3-category scan produces 6 items (stub returns same item per category)"
else
    assert_fail "2-pool × 3-category scan should produce 6 items, got $item_count7"
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

# scan.sh:P1 refactor 后无 claim-tag exclusion(dedup 下推到 triage);不应有 NOT tag=jarvis-claimed
# 注:`grep -c` 无匹配 exit 1 但 stdout 已是 "0",不加 `|| echo 0`(避免追加 "0\n0" 让整数比较失败)
filter_count7=$(echo "$recorded_args7" | grep -c 'NOT tag=jarvis-claimed')
if [ "$filter_count7" -eq 0 ]; then
    assert_pass "no NOT tag=jarvis-claimed (P1 refactor removed claim-tag exclusion from scan)"
else
    assert_fail "NOT tag=jarvis-claimed should NOT appear (removed by refactor), got $filter_count7 times"
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

output8=$(JARVIS_ROOT="$tmpconfig8" bash "$proj_root/bootstrap/scan.sh" 2>/dev/null)  # stdout(JSON)only; WARN 走 stderr
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

# scan.sh 展开 3 categories,healthy pool 每 category 返同一 stub item = 3 items;failing pool 跳过
item_count8=$(echo "$output8" | jq 'length' 2>/dev/null || echo 0)
if [ "$item_count8" -eq 3 ]; then
    assert_pass "healthy pool × 3 categories = 3 items despite failing pool"
else
    assert_fail "should return 3 items from healthy pool (3 categories), got $item_count8: $output8"
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
# Test 10: --columns id,title,status,priority,tag passed to a1 list (pool-scoped)
# ---------------------------------------------------------------------------

echo "=== Test 10: --columns id,title,status,priority,tag passed to a1 list (pool-scoped) ==="

tmpconfig10=$(mktemp -d)
mkdir -p "$tmpconfig10/config"
cat > "$tmpconfig10/config/pools.json" << 'JSON'
{
  "pools": {
    "pool_cols": { "project": 777, "name": "Cols Pool" }
  },
  "claim": { "tag": "jarvis-claimed" }
}
JSON

args_file10=$(mktemp)

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "auth" ] && [ "\$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo "\$@" >> "$args_file10"
    echo '[{"identifier":"WI-777","subject":"Cols item","categoryIdentifier":"task","status":"open","priority":"medium","tag":["qa"]}]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

output10=$(JARVIS_ROOT="$tmpconfig10" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code10=$?
recorded_args10=$(cat "$args_file10" 2>/dev/null || echo "")
rm -f "$args_file10"
rm -rf "$tmpconfig10"

echo "Exit code: $exit_code10"
echo "Recorded a1 args: $recorded_args10"
echo "Output: $output10"
echo ""

if [ "$exit_code10" -eq 0 ]; then
    assert_pass "exit code 0 with --columns test"
else
    assert_fail "exit code should be 0 with --columns test, got $exit_code10"
fi

if echo "$recorded_args10" | grep -q -- '--columns'; then
    assert_pass "--columns flag passed to a1 list"
else
    assert_fail "--columns flag not found in a1 args: $recorded_args10"
fi

if echo "$recorded_args10" | grep -q 'id,title,status,priority,tag'; then
    assert_pass "--columns id,title,status,priority,tag present in a1 args"
else
    assert_fail "--columns id,title,status,priority,tag not found in a1 args: $recorded_args10"
fi

# Output object should have priority and tag fields
if echo "$output10" | jq -e '.[0].priority == "medium"' > /dev/null 2>&1; then
    assert_pass "pool-scoped output carries priority field"
else
    assert_fail "pool-scoped output missing priority field: $output10"
fi

if echo "$output10" | jq -e '.[0].tag == ["qa"]' > /dev/null 2>&1; then
    assert_pass "pool-scoped output carries tag field"
else
    assert_fail "pool-scoped output missing tag field: $output10"
fi

# ---------------------------------------------------------------------------
# Test 11: --columns id,title,status,priority,tag passed to a1 list (no-pools fallback)
# ---------------------------------------------------------------------------

echo "=== Test 11: --columns passed in no-pools fallback path ==="

tmpconfig11=$(mktemp -d)
# No pools.json at all → falls back to global assignee list
args_file11=$(mktemp)

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "auth" ] && [ "\$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo "\$@" >> "$args_file11"
    echo '[{"identifier":"WI-888","subject":"Global item","categoryIdentifier":"req","status":"open","priority":"low","tag":[]}]'
    exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

output11=$(JARVIS_ROOT="$tmpconfig11" bash "$proj_root/bootstrap/scan.sh" 2>&1)
exit_code11=$?
recorded_args11=$(cat "$args_file11" 2>/dev/null || echo "")
rm -f "$args_file11"
rm -rf "$tmpconfig11"

echo "Exit code: $exit_code11"
echo "Recorded a1 args: $recorded_args11"
echo "Output: $output11"
echo ""

if [ "$exit_code11" -eq 0 ]; then
    assert_pass "exit code 0 for no-pools fallback with --columns"
else
    assert_fail "exit code should be 0 for no-pools fallback, got $exit_code11"
fi

if echo "$recorded_args11" | grep -q -- '--columns'; then
    assert_pass "--columns flag passed in no-pools fallback"
else
    assert_fail "--columns flag not found in no-pools fallback args: $recorded_args11"
fi

if echo "$recorded_args11" | grep -q 'id,title,status,priority,tag'; then
    assert_pass "--columns id,title,status,priority,tag in no-pools fallback"
else
    assert_fail "--columns id,title,status,priority,tag not found in no-pools fallback args: $recorded_args11"
fi

if echo "$output11" | jq -e '.[0].priority == "low"' > /dev/null 2>&1; then
    assert_pass "no-pools fallback output carries priority field"
else
    assert_fail "no-pools fallback output missing priority field: $output11"
fi

if echo "$output11" | jq -e '.[0].tag == []' > /dev/null 2>&1; then
    assert_pass "no-pools fallback output carries tag field (empty array)"
else
    assert_fail "no-pools fallback output missing tag field: $output11"
fi

# ---------------------------------------------------------------------------
# Test 12: JARVIS_SCAN_ASSIGNEE overrides the scan target (decoupled from whoami)
# ---------------------------------------------------------------------------
cat > "$tmpbin/a1" << 'STUB'
#!/bin/bash
if [ "$1" = "auth" ] && [ "$2" = "whoami" ]; then
    echo "Account:  chenhanzhang.chz"; exit 0
fi
if [ "$1" = "project" ] && [ "$2" = "workitem" ] && [ "$3" = "list" ]; then
    prev=""
    for a in "$@"; do
        [ "$prev" = "--assignee" ] && echo "$a" >> "$ASSIGNEE_CAP"
        prev="$a"
    done
    echo '[]'; exit 0
fi
exit 1
STUB
chmod +x "$tmpbin/a1"

tmpcfgA=$(mktemp -d); mkdir -p "$tmpcfgA/config"
cat > "$tmpcfgA/config/pools.json" << 'JSON'
{ "pools": { "p": { "project": 12345, "name": "P" } }, "claim": { "tag": "jarvis-claimed" } }
JSON

echo "=== Test 12: JARVIS_SCAN_ASSIGNEE override vs default whoami ==="
capA=$(mktemp)
JARVIS_ROOT="$tmpcfgA" JARVIS_SCAN_TTL=0 JARVIS_SCAN_ASSIGNEE=320687 ASSIGNEE_CAP="$capA" \
    bash "$proj_root/bootstrap/scan.sh" >/dev/null 2>&1
if grep -qx "320687" "$capA"; then
    assert_pass "JARVIS_SCAN_ASSIGNEE overrides --assignee (320687)"
else
    assert_fail "expected --assignee 320687, got: $(sort -u "$capA" | tr '\n' ',')"
fi
if ! grep -qx "chenhanzhang.chz" "$capA"; then
    assert_pass "override suppresses the whoami assignee"
else
    assert_fail "whoami assignee leaked despite override"
fi

: > "$capA"
JARVIS_ROOT="$tmpcfgA" JARVIS_SCAN_TTL=0 ASSIGNEE_CAP="$capA" \
    bash "$proj_root/bootstrap/scan.sh" >/dev/null 2>&1
if grep -qx "chenhanzhang.chz" "$capA"; then
    assert_pass "default assignee falls back to whoami account"
else
    assert_fail "expected whoami assignee by default, got: $(sort -u "$capA" | tr '\n' ',')"
fi
rm -rf "$tmpcfgA"; rm -f "$capA"

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
