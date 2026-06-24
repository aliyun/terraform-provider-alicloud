#!/usr/bin/env bash
# bootstrap/tests/wrap_done.sh — unit tests for wrap.sh done subcommand
# Tests: missing status→nonzero, full args+a1 ok→0, A1_FAIL=1→nonzero
# Run: bash bootstrap/tests/wrap_done.sh
# Prints PASS and exits 0 on success; prints FAIL and exits 1 on any failure.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
JARVIS_ROOT="$(cd "$BOOTSTRAP_DIR/.." && pwd)"

# ---------------------------------------------------------------------------
# Fake a1 binary: exits 0 normally, exits 1 when A1_FAIL=1
# ---------------------------------------------------------------------------
FAKE_BIN_DIR="$(mktemp -d)"
trap 'rm -rf "$FAKE_BIN_DIR" "$RUNS_DIR" "$POOLS_TMP"' EXIT

cat > "$FAKE_BIN_DIR/a1" <<'EOF'
#!/usr/bin/env bash
if [ "${A1_FAIL:-0}" = "1" ]; then
    echo "a1: simulated failure" >&2
    exit 1
fi
exit 0
EOF
chmod +x "$FAKE_BIN_DIR/a1"

# Fake cache.sh stub (so wrap.sh can call it without side effects)
cat > "$FAKE_BIN_DIR/cache.sh" <<'EOF'
#!/usr/bin/env bash
exit 0
EOF
chmod +x "$FAKE_BIN_DIR/cache.sh"

# Prepend fake bin to PATH so wrap.sh picks up fake a1
export PATH="$FAKE_BIN_DIR:$PATH"

# ---------------------------------------------------------------------------
# Temp runs dir so log.sh run_done writes there
# ---------------------------------------------------------------------------
RUNS_DIR="$(mktemp -d)"
export JARVIS_RUNS_DIR="$RUNS_DIR"

# ---------------------------------------------------------------------------
# Temp pools.json so wrap.sh doesn't need the real config
# ---------------------------------------------------------------------------
POOLS_TMP="$(mktemp)"
cat > "$POOLS_TMP" <<'EOF'
{
  "claim": {
    "label": "jarvis-claimed",
    "done_label": "jarvis-done",
    "ttl_min": 45
  }
}
EOF

# Override JARVIS_ROOT so wrap.sh finds our temp pools.json
export JARVIS_ROOT
# We patch pools_cfg by providing a symlink from config/pools.json
# Actually, wrap.sh resolves pools_cfg from $JARVIS_ROOT/config/pools.json
# So we need to make a temporary jarvis_root with config/pools.json
FAKE_JARVIS_ROOT="$(mktemp -d)"
mkdir -p "$FAKE_JARVIS_ROOT/config"
cp "$POOLS_TMP" "$FAKE_JARVIS_ROOT/config/pools.json"
# Also symlink bootstrap/ so wrap.sh can call log.sh and cache.sh
mkdir -p "$FAKE_JARVIS_ROOT/bootstrap"
# Provide a real log.sh in the fake root that uses JARVIS_RUNS_DIR
cp "$BOOTSTRAP_DIR/log.sh" "$FAKE_JARVIS_ROOT/bootstrap/log.sh"
# Provide a stub cache.sh in the fake root
cat > "$FAKE_JARVIS_ROOT/bootstrap/cache.sh" <<'CEOF'
#!/usr/bin/env bash
exit 0
CEOF
chmod +x "$FAKE_JARVIS_ROOT/bootstrap/cache.sh"
export JARVIS_ROOT="$FAKE_JARVIS_ROOT"

WRAP="$BOOTSTRAP_DIR/wrap.sh"

# ---------------------------------------------------------------------------
# Helper
# ---------------------------------------------------------------------------
pass_count=0
fail_count=0

assert_exit() {
    local desc="$1"
    local expected="$2"   # "nonzero" or "zero"
    shift 2
    local actual_exit=0
    "$@" >/dev/null 2>&1 || actual_exit=$?
    if [ "$expected" = "zero" ] && [ "$actual_exit" -eq 0 ]; then
        echo "  PASS: $desc"
        pass_count=$((pass_count + 1))
    elif [ "$expected" = "nonzero" ] && [ "$actual_exit" -ne 0 ]; then
        echo "  PASS: $desc"
        pass_count=$((pass_count + 1))
    else
        echo "  FAIL: $desc (expected=$expected actual_exit=$actual_exit)"
        fail_count=$((fail_count + 1))
    fi
}

# ---------------------------------------------------------------------------
# Test 1: missing status arg → nonzero exit
# ---------------------------------------------------------------------------
echo "Test 1: missing status → nonzero"
assert_exit "wrap.sh done missing status" nonzero \
    bash "$WRAP" done "WI-001" "some summary"

# ---------------------------------------------------------------------------
# Test 2: full args + a1 succeeds → exit 0, run_done file written
# ---------------------------------------------------------------------------
echo "Test 2: full args + a1 ok → zero"
A1_FAIL=0 assert_exit "wrap.sh done full args ok" zero \
    bash "$WRAP" done "WI-002" "all good" "done"

# Verify run_done file was written
run_file_count=$(ls "$RUNS_DIR"/*-WI-002.md 2>/dev/null | wc -l | tr -d ' ')
if [ "$run_file_count" -ge 1 ]; then
    echo "  PASS: run_done file written for WI-002"
    pass_count=$((pass_count + 1))
else
    echo "  FAIL: run_done file NOT written for WI-002"
    fail_count=$((fail_count + 1))
fi

# ---------------------------------------------------------------------------
# Test 3: A1_FAIL=1 → nonzero exit (a1 failures now exit 1)
# ---------------------------------------------------------------------------
echo "Test 3: A1_FAIL=1 → nonzero"
assert_exit "wrap.sh done a1 failure exits nonzero" nonzero \
    env A1_FAIL=1 bash "$WRAP" done "WI-003" "some summary" "done"

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
