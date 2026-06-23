#!/usr/bin/env bash
# test/sweep_test.sh – TDD tests for bootstrap/sweep.sh
# Stubs a1; asserts stale jarvis-claimed items produce escalation/<id>.md,
# and fresh claims do not.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

PASS=0
FAIL=0

assert_pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
assert_fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# ---------------------------------------------------------------------------
# Shared setup: temp dirs + minimal pools.json
# ---------------------------------------------------------------------------
tmpbin=$(mktemp -d)
tmproot=$(mktemp -d)
tmpesc=$(mktemp -d)
trap 'rm -rf "$tmpbin" "$tmproot" "$tmpesc"' EXIT

mkdir -p "$tmproot/config"
cat > "$tmproot/config/pools.json" << 'JSON'
{
  "pools": {
    "test_pool": { "project": 9999 }
  },
  "claim": { "tag": "jarvis-claimed", "ttl_min": 45 }
}
JSON

# ---------------------------------------------------------------------------
# Test 1: stale claimed item (comment > 45 min old) → escalation/<id>.md created
# ---------------------------------------------------------------------------
echo "=== Test 1: stale claim → escalation file created ==="

STALE_ID="42001"

# Timestamp 60 minutes in the past
OLD_TS=$(date -u -v -60M +%Y-%m-%dT%H:%M:%SZ 2>/dev/null \
    || date -u -d "60 minutes ago" +%Y-%m-%dT%H:%M:%SZ)

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo '[{"id":"$STALE_ID"}]'
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "comment" ] && [ "\$4" = "list" ]; then
    echo '[{"content":"jarvis-claim host-A $OLD_TS","createdAt":"$OLD_TS"}]'
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"

output1=$(PATH="$tmpbin:$PATH" \
    JARVIS_ROOT="$tmproot" \
    JARVIS_ESCALATION_DIR="$tmpesc" \
    bash "$proj_root/bootstrap/sweep.sh" 2>&1)
exit1=$?

echo "Output: $output1"
echo "Exit code: $exit1"

if [ -f "$tmpesc/${STALE_ID}.md" ]; then
    assert_pass "stale claim: escalation/$STALE_ID.md created"
else
    assert_fail "stale claim: escalation/$STALE_ID.md not found (tmpesc: $tmpesc, files: $(ls "$tmpesc" 2>&1))"
fi

if [ "$exit1" -eq 0 ]; then
    assert_pass "sweep exits 0"
else
    assert_fail "sweep should exit 0, got $exit1"
fi

# stale id should appear in stdout
if echo "$output1" | grep -q "$STALE_ID"; then
    assert_pass "stale id printed to stdout"
else
    assert_fail "stale id not found in stdout: $output1"
fi

# ---------------------------------------------------------------------------
# Test 2: fresh claimed item (comment < 45 min old) → no escalation file
# ---------------------------------------------------------------------------
echo ""
echo "=== Test 2: fresh claim → no escalation file ==="

FRESH_ID="42002"
rm -f "$tmpesc/${FRESH_ID}.md"

# Timestamp 5 minutes in the past (fresh)
NEW_TS=$(date -u -v -5M +%Y-%m-%dT%H:%M:%SZ 2>/dev/null \
    || date -u -d "5 minutes ago" +%Y-%m-%dT%H:%M:%SZ)

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo '[{"id":"$FRESH_ID"}]'
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "comment" ] && [ "\$4" = "list" ]; then
    echo '[{"content":"jarvis-claim host-B $NEW_TS","createdAt":"$NEW_TS"}]'
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"

output2=$(PATH="$tmpbin:$PATH" \
    JARVIS_ROOT="$tmproot" \
    JARVIS_ESCALATION_DIR="$tmpesc" \
    bash "$proj_root/bootstrap/sweep.sh" 2>&1)
exit2=$?

echo "Output: $output2"
echo "Exit code: $exit2"

if [ ! -f "$tmpesc/${FRESH_ID}.md" ]; then
    assert_pass "fresh claim: no escalation/$FRESH_ID.md (correct)"
else
    assert_fail "fresh claim: escalation/$FRESH_ID.md should NOT exist"
fi

if [ "$exit2" -eq 0 ]; then
    assert_pass "sweep exits 0 on fresh claim"
else
    assert_fail "sweep should exit 0 on fresh claim, got $exit2"
fi

# ---------------------------------------------------------------------------
# Test 3: item with no jarvis-claim comment → no escalation (skip gracefully)
# ---------------------------------------------------------------------------
echo ""
echo "=== Test 3: no jarvis-claim comment → skipped gracefully ==="

NO_CLAIM_ID="42003"
rm -f "$tmpesc/${NO_CLAIM_ID}.md"

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo '[{"id":"$NO_CLAIM_ID"}]'
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "comment" ] && [ "\$4" = "list" ]; then
    echo '[]'
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"

output3=$(PATH="$tmpbin:$PATH" \
    JARVIS_ROOT="$tmproot" \
    JARVIS_ESCALATION_DIR="$tmpesc" \
    bash "$proj_root/bootstrap/sweep.sh" 2>&1)
exit3=$?

echo "Output: $output3"
echo "Exit code: $exit3"

if [ ! -f "$tmpesc/${NO_CLAIM_ID}.md" ]; then
    assert_pass "no-comment item: no escalation file created"
else
    assert_fail "no-comment item: escalation should NOT be created"
fi

if [ "$exit3" -eq 0 ]; then
    assert_pass "sweep exits 0 when no claim comment found"
else
    assert_fail "sweep should exit 0 when no claim comment found, got $exit3"
fi

# ---------------------------------------------------------------------------
# Test 4: item with malformed claim timestamp → no crash, no escalation, exit 0
# ---------------------------------------------------------------------------
echo ""
echo "=== Test 4: malformed timestamp in claim comment → graceful skip ==="

MALFORMED_ID="42004"
rm -f "$tmpesc/${MALFORMED_ID}.md"

cat > "$tmpbin/a1" << STUB
#!/bin/bash
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "list" ]; then
    echo '[{"id":"$MALFORMED_ID"}]'
    exit 0
fi
if [ "\$1" = "project" ] && [ "\$2" = "workitem" ] && [ "\$3" = "comment" ] && [ "\$4" = "list" ]; then
    echo '[{"content":"jarvis-claim host-C 9999-99-99T99:99:99Z","createdAt":"2099-99-99T99:99:99Z"}]'
    exit 0
fi
exit 0
STUB
chmod +x "$tmpbin/a1"

output4=$(PATH="$tmpbin:$PATH" \
    JARVIS_ROOT="$tmproot" \
    JARVIS_ESCALATION_DIR="$tmpesc" \
    bash "$proj_root/bootstrap/sweep.sh" 2>&1)
exit4=$?
# Capture stderr separately for the WARN check
err4=$(PATH="$tmpbin:$PATH" \
    JARVIS_ROOT="$tmproot" \
    JARVIS_ESCALATION_DIR="$tmpesc" \
    bash "$proj_root/bootstrap/sweep.sh" 2>&1 1>/dev/null)
output4_full="$output4"$'\n'"$err4"

echo "Output: $output4"
echo "Exit code: $exit4"

if [ ! -f "$tmpesc/${MALFORMED_ID}.md" ]; then
    assert_pass "malformed timestamp: no escalation file created"
else
    assert_fail "malformed timestamp: escalation should NOT be created"
fi

if [ "$exit4" -eq 0 ]; then
    assert_pass "sweep exits 0 on malformed timestamp (graceful)"
else
    assert_fail "sweep should exit 0 on malformed timestamp, got $exit4"
fi

# Should see a WARN message about parse failure
if echo "$output4_full" | grep -q "WARN.*could not parse"; then
    assert_pass "malformed timestamp: WARN message logged"
else
    assert_fail "malformed timestamp: WARN message not found in output"
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
