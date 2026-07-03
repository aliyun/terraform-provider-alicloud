#!/bin/bash

# Test harness for verify.sh's GitHub token daily check.
# An invalid/missing JARVIS_GITHUB_TOKEN must be a WARN (non-fatal) — it only
# blocks the GitHub escalation path, not Aone-only work — and must drop an
# idempotent escalation note. Stubs a1/aliyun/cloudspec to succeed and gh to
# fail so github-identity.sh check fails while every other check passes, making
# the overall exit code deterministic.

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

tmpbin=$(mktemp -d)
esc_dir=$(mktemp -d)
verify_output=$(mktemp)
trap "rm -rf $tmpbin $esc_dir; rm -f $verify_output" EXIT

# Stubs: creds succeed so only the github-token check can fail.
for c in a1 aliyun cloudspec; do
    cat > "$tmpbin/$c" << 'EOF'
#!/bin/bash
exit 0
EOF
    chmod +x "$tmpbin/$c"
done
# gh fails → github-identity.sh check fails (401-style) → WARN path.
cat > "$tmpbin/gh" << 'EOF'
#!/bin/bash
exit 1
EOF
chmod +x "$tmpbin/gh"

export PATH="$tmpbin:$PATH"
# Force the invalid-token path deterministically (models a 401 stale token).
export JARVIS_GITHUB_TOKEN="dummy-invalid-token"
# Redirect the escalation note away from the real repo.
export JARVIS_ESCALATION_DIR="$esc_dir"

# --- Run 1 -----------------------------------------------------------------
bash "$proj_root/bootstrap/verify.sh" > "$verify_output" 2>&1
exit_code=$?
output=$(cat "$verify_output")

echo "=== Test Output (run 1) ==="
echo "$output"
echo ""
echo "=== Exit Code ==="
echo "$exit_code"
echo ""

fail=0
echo "=== Assertions ==="

# Regression guard: no dependence on ambient gh auth.
if echo "$output" | grep -q "gh-cred"; then
    echo "FAIL: Output should not depend on ambient gh auth"; fail=1
else
    echo "PASS: Output does not contain ambient gh auth check"
fi

# Token check is WARN, not FAIL.
if echo "$output" | grep -q "WARN jarvis-github-token"; then
    echo "PASS: Output contains 'WARN jarvis-github-token'"
else
    echo "FAIL: Output does NOT contain 'WARN jarvis-github-token'"; fail=1
fi
if echo "$output" | grep -q "FAIL jarvis-github-token"; then
    echo "FAIL: token check must be WARN, not FAIL"; fail=1
else
    echo "PASS: Output does not FAIL on the github token"
fi

# WARN must not fail the overall run (all other checks pass via stubs).
if [ "$exit_code" -eq 0 ]; then
    echo "PASS: Exit code is 0 (WARN does not block preflight)"
else
    echo "FAIL: Exit code is $exit_code (WARN should not fail overall)"; fail=1
fi

# Escalation note dropped.
esc_glob=("$esc_dir"/github-token-invalid-*.md)
if [ -f "${esc_glob[0]}" ]; then
    echo "PASS: escalation note created (${esc_glob[0]##*/})"
else
    echo "FAIL: escalation note not created in $esc_dir"; fail=1
fi
if grep -q "修复步骤" "${esc_glob[0]}" 2>/dev/null; then
    echo "PASS: escalation note has 修复步骤 section"
else
    echo "FAIL: escalation note missing 修复步骤 section"; fail=1
fi

# --- Run 2: idempotency ----------------------------------------------------
before_count=$(ls "$esc_dir"/github-token-invalid-*.md 2>/dev/null | wc -l | tr -d ' ')
before_hash=$(cat "$esc_dir"/github-token-invalid-*.md 2>/dev/null | shasum | awk '{print $1}')

bash "$proj_root/bootstrap/verify.sh" > "$verify_output" 2>&1

after_count=$(ls "$esc_dir"/github-token-invalid-*.md 2>/dev/null | wc -l | tr -d ' ')
after_hash=$(cat "$esc_dir"/github-token-invalid-*.md 2>/dev/null | shasum | awk '{print $1}')

if [ "$before_count" = "$after_count" ] && [ "$before_count" = "1" ]; then
    echo "PASS: escalation note idempotent (still 1 file after re-run)"
else
    echo "FAIL: escalation note not idempotent (before=$before_count after=$after_count)"; fail=1
fi
if [ "$before_hash" = "$after_hash" ]; then
    echo "PASS: escalation note content unchanged on re-run"
else
    echo "FAIL: escalation note content changed on re-run"; fail=1
fi

echo ""
if [ "$fail" -eq 0 ]; then
    echo "✓ All test assertions passed"
    exit 0
else
    echo "✗ Some assertions failed"
    exit 1
fi
