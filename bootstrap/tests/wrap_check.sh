#!/usr/bin/env bash
# bootstrap/tests/wrap_check.sh — unit tests for wrap-check.sh
#
# Scenarios:
#   1. done:false id with no runs/ file → exit 2
#   2. done:false id with a runs/ file → exit 0
#   3. done:true id with no runs/ file → exit 0 (done:true is exempt)
#   4. no claims file → exit 0
#
# Run: bash bootstrap/tests/wrap_check.sh
# Prints PASS and exits 0 on success; prints FAIL and exits 1 on any failure.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

WRAP_CHECK="$BOOTSTRAP_DIR/wrap-check.sh"

# ---------------------------------------------------------------------------
# Fake a1 stub (must be in PATH to prevent any accidental a1 call)
# ---------------------------------------------------------------------------
FAKE_BIN_DIR="$(mktemp -d)"
cat > "$FAKE_BIN_DIR/a1" <<'EOF'
#!/usr/bin/env bash
echo "a1: should not be called in wrap-check tests" >&2
exit 99
EOF
chmod +x "$FAKE_BIN_DIR/a1"
export PATH="$FAKE_BIN_DIR:$PATH"

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
pass_count=0
fail_count=0

assert_exit_code() {
    local desc="$1"
    local expected="$2"
    shift 2
    local actual=0
    "$@" >/dev/null 2>&1 || actual=$?
    if [ "$actual" -eq "$expected" ]; then
        echo "  PASS: $desc (exit $actual)"
        pass_count=$((pass_count + 1))
    else
        echo "  FAIL: $desc (expected=$expected actual=$actual)"
        fail_count=$((fail_count + 1))
    fi
}

make_jarvis_root() {
    local root
    root="$(mktemp -d)"
    mkdir -p "$root/.my-day" "$root/runs"
    # minimal config/pools.json so log.sh/wrap-check.sh can find it
    mkdir -p "$root/config"
    cp "$BOOTSTRAP_DIR/../config/pools.json" "$root/config/pools.json"
    echo "$root"
}

today="$(date -u +%F)"

# ---------------------------------------------------------------------------
# Declare all temp roots up front so the trap never references an unset variable
# ---------------------------------------------------------------------------
ROOT1=""
ROOT2=""
ROOT3=""
ROOT4=""
trap 'rm -rf "$FAKE_BIN_DIR" "$ROOT1" "$ROOT2" "$ROOT3" "$ROOT4" 2>/dev/null; exit' INT TERM EXIT

# ---------------------------------------------------------------------------
# Test 1: done:false id, no runs/ file → exit 2
# ---------------------------------------------------------------------------
echo "Test 1: done:false + no run file → exit 2"
ROOT1="$(make_jarvis_root)"

claims_file="$ROOT1/.my-day/claims-${today}.json"
printf '[{"id":"WI-T1","done":false}]\n' > "$claims_file"

assert_exit_code "missing run file → exit 2" 2 \
    env JARVIS_ROOT="$ROOT1" JARVIS_RUNS_DIR="$ROOT1/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 2: done:false id, runs/ file present → exit 0
# ---------------------------------------------------------------------------
echo "Test 2: done:false + run file present → exit 0"
ROOT2="$(make_jarvis_root)"
claims_file2="$ROOT2/.my-day/claims-${today}.json"
printf '[{"id":"WI-T2","done":false}]\n' > "$claims_file2"
# Create the run file that log.sh seen would match
touch "$ROOT2/runs/${today}-WI-T2.md"

assert_exit_code "run file present → exit 0" 0 \
    env JARVIS_ROOT="$ROOT2" JARVIS_RUNS_DIR="$ROOT2/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 3: done:true id, no runs/ file → still exit 0 (done items exempt)
# ---------------------------------------------------------------------------
echo "Test 3: done:true + no run file → exit 0"
ROOT3="$(make_jarvis_root)"
claims_file3="$ROOT3/.my-day/claims-${today}.json"
printf '[{"id":"WI-T3","done":true}]\n' > "$claims_file3"

assert_exit_code "done:true exempt → exit 0" 0 \
    env JARVIS_ROOT="$ROOT3" JARVIS_RUNS_DIR="$ROOT3/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Test 4: no claims file at all → exit 0
# ---------------------------------------------------------------------------
echo "Test 4: no claims file → exit 0"
ROOT4="$(make_jarvis_root)"
# don't write any claims file

assert_exit_code "no claims file → exit 0" 0 \
    env JARVIS_ROOT="$ROOT4" JARVIS_RUNS_DIR="$ROOT4/runs" bash "$WRAP_CHECK"

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo ""
echo "Results: $pass_count passed, $fail_count failed"
if [ "$fail_count" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    exit 1
fi
