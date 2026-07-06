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
esc_dir_a1=$(mktemp -d)
esc_dir_a1_hd=$(mktemp -d)
verify_output=$(mktemp)
trap "rm -rf $tmpbin $esc_dir $esc_dir_a1 $esc_dir_a1_hd; rm -f $verify_output" EXIT

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

# Stub bin/a1id for verify.sh's a1 login-state check (invoked as `a1id -- auth whoami`).
# FAKE_A1ID_STATE drives the three states the check must distinguish:
#   normal   → jarvis alive: Account=WORKER_1782379562571
#   expired  → dead session: error to stderr, non-zero exit, no stdout
#   halfdead → session half-alive: Emp ID present but Account field empty
cat > "$tmpbin/a1id" << 'EOF'
#!/bin/bash
case "${FAKE_A1ID_STATE:-normal}" in
  expired)
    echo "a1id: 身份 'jarvis' 未登录,先 'bin/a1id login jarvis'" >&2
    exit 1 ;;
  halfdead)
    printf 'Account:  \nName:     Jarvis Worker\nEmp ID:   1782379562571\nEmail:    open_jarvis@example.com\n' ;;
  *)
    printf 'Account:  WORKER_1782379562571\nName:     Jarvis Worker\nEmp ID:   1782379562571\nEmail:    open_jarvis@example.com\n' ;;
esac
EOF
chmod +x "$tmpbin/a1id"

export PATH="$tmpbin:$PATH"
# Force the invalid-token path deterministically (models a 401 stale token).
export JARVIS_GITHUB_TOKEN="dummy-invalid-token"
# Redirect the escalation note away from the real repo.
export JARVIS_ESCALATION_DIR="$esc_dir"
# a1 login-state check calls bin/a1id; point it at the stub. Default state = normal
# so Run 1/2 (github-focused) see a healthy a1 session (PASS jarvis-a1-session).
export JARVIS_A1ID="$tmpbin/a1id"

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

# a1 login-state 正常态: default stub → Account=WORKER_… → PASS (covers 正常).
if echo "$output" | grep -q "PASS jarvis-a1-session"; then
    echo "PASS: a1 login-state OK when Account=WORKER_1782379562571 (正常态)"
else
    echo "FAIL: expected 'PASS jarvis-a1-session' for healthy a1 session"; fail=1
fi
# 正常态不应落 a1 过期 escalation note.
if ls "$esc_dir"/a1-session-expired-*.md >/dev/null 2>&1; then
    echo "FAIL: 正常态不应落 a1-session-expired escalation note"; fail=1
else
    echo "PASS: 正常态无 a1-session escalation note"
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

# --- a1 login-state check: 过期态 (expired) --------------------------------
echo ""
echo "=== a1 login-state check (过期 / 半死) ==="

a1_out=$(FAKE_A1ID_STATE=expired JARVIS_ESCALATION_DIR="$esc_dir_a1" \
    bash "$proj_root/bootstrap/verify.sh" 2>&1)
a1_rc=$?

# WARN, not FAIL.
if echo "$a1_out" | grep -q "WARN jarvis-a1-session"; then
    echo "PASS: 过期态 → WARN jarvis-a1-session"
else
    echo "FAIL: 过期态未 WARN jarvis-a1-session"; fail=1
fi
if echo "$a1_out" | grep -q "FAIL jarvis-a1-session"; then
    echo "FAIL: a1 check must be WARN, not FAIL"; fail=1
else
    echo "PASS: 过期态 a1 check 不 FAIL"
fi
# WARN must not fail the overall run.
if [ "$a1_rc" -eq 0 ]; then
    echo "PASS: 过期态 WARN 不阻断 preflight (exit 0)"
else
    echo "FAIL: 过期态 exit=$a1_rc (WARN 不应 fail overall)"; fail=1
fi
# Fix guidance present.
if echo "$a1_out" | grep -q "bin/a1id login jarvis"; then
    echo "PASS: 过期态含修复指引 (bin/a1id login jarvis)"
else
    echo "FAIL: 过期态缺修复指引 bin/a1id login jarvis"; fail=1
fi
# Escalation note dropped.
a1_esc=("$esc_dir_a1"/a1-session-expired-*.md)
if [ -f "${a1_esc[0]}" ]; then
    echo "PASS: a1 escalation note created (${a1_esc[0]##*/})"
else
    echo "FAIL: a1 escalation note not created in $esc_dir_a1"; fail=1
fi
if grep -q "修复步骤" "${a1_esc[0]}" 2>/dev/null; then
    echo "PASS: a1 escalation note has 修复步骤 section"
else
    echo "FAIL: a1 escalation note missing 修复步骤 section"; fail=1
fi

# 过期态幂等: 再跑一次,note 仍 1 份且内容不变.
a1_before_count=$(ls "$esc_dir_a1"/a1-session-expired-*.md 2>/dev/null | wc -l | tr -d ' ')
a1_before_hash=$(cat "$esc_dir_a1"/a1-session-expired-*.md 2>/dev/null | shasum | awk '{print $1}')
FAKE_A1ID_STATE=expired JARVIS_ESCALATION_DIR="$esc_dir_a1" \
    bash "$proj_root/bootstrap/verify.sh" >/dev/null 2>&1
a1_after_count=$(ls "$esc_dir_a1"/a1-session-expired-*.md 2>/dev/null | wc -l | tr -d ' ')
a1_after_hash=$(cat "$esc_dir_a1"/a1-session-expired-*.md 2>/dev/null | shasum | awk '{print $1}')
if [ "$a1_before_count" = "$a1_after_count" ] && [ "$a1_before_count" = "1" ]; then
    echo "PASS: a1 escalation note idempotent (still 1 file after re-run)"
else
    echo "FAIL: a1 escalation note not idempotent (before=$a1_before_count after=$a1_after_count)"; fail=1
fi
if [ "$a1_before_hash" = "$a1_after_hash" ]; then
    echo "PASS: a1 escalation note content unchanged on re-run"
else
    echo "FAIL: a1 escalation note content changed on re-run"; fail=1
fi

# --- a1 login-state check: 半死态 (Account 空, EmpID 在) --------------------
hd_out=$(FAKE_A1ID_STATE=halfdead JARVIS_ESCALATION_DIR="$esc_dir_a1_hd" \
    bash "$proj_root/bootstrap/verify.sh" 2>&1)
# Account 空 → 仍判 WARN (not OK).
if echo "$hd_out" | grep -q "WARN jarvis-a1-session"; then
    echo "PASS: 半死态(Account 空 EmpID 在) → WARN jarvis-a1-session"
else
    echo "FAIL: 半死态未 WARN jarvis-a1-session"; fail=1
fi
if echo "$hd_out" | grep -q "PASS jarvis-a1-session"; then
    echo "FAIL: 半死态不应判 OK (Account 空)"; fail=1
else
    echo "PASS: 半死态未误判 OK"
fi
# 半死态亦落 note (独立 esc dir,证明半死也触发 escalation).
hd_esc=("$esc_dir_a1_hd"/a1-session-expired-*.md)
if [ -f "${hd_esc[0]}" ]; then
    echo "PASS: 半死态亦落 a1 escalation note (${hd_esc[0]##*/})"
else
    echo "FAIL: 半死态未落 a1 escalation note"; fail=1
fi

echo ""
if [ "$fail" -eq 0 ]; then
    echo "✓ All test assertions passed"
    exit 0
else
    echo "✗ Some assertions failed"
    exit 1
fi
