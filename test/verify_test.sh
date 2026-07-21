#!/bin/bash
# Verify degraded credentials remain WARN-only without writing local escalation files.

set -u

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"
tmpbin="$(mktemp -d)"
verify_output="$(mktemp)"
trap 'rm -rf "$tmpbin"; rm -f "$verify_output"' EXIT

for c in a1 aliyun cloudspec; do
    printf '#!/bin/bash\nexit 0\n' > "$tmpbin/$c"
    chmod +x "$tmpbin/$c"
done
printf '#!/bin/bash\necho "HTTP 401: Bad credentials" >&2\nexit 1\n' > "$tmpbin/gh"
chmod +x "$tmpbin/gh"
cat > "$tmpbin/a1id" <<'EOF'
#!/bin/bash
case "${FAKE_A1ID_STATE:-normal}" in
  expired) exit 1 ;;
  halfdead) printf 'Account:  \nName: Jarvis Worker\nEmp ID: 1782379562571\n' ;;
  *) printf 'Account:  WORKER_1782379562571\nName: Jarvis Worker\nEmp ID: 1782379562571\n' ;;
esac
EOF
chmod +x "$tmpbin/a1id"

export PATH="$tmpbin:$PATH"
export JARVIS_GITHUB_TOKEN="dummy-invalid-token"
export JARVIS_A1ID="$tmpbin/a1id"

fail=0
run_verify() {
    FAKE_A1ID_STATE="$1" bash "$proj_root/bootstrap/verify.sh" > "$verify_output" 2>&1
}

run_verify normal || { echo "FAIL: healthy a1 + invalid GitHub token should stay non-fatal"; fail=1; }
grep -q 'WARN jarvis-github-token' "$verify_output" || { echo "FAIL: missing GitHub token WARN"; fail=1; }
grep -q 'PASS jarvis-a1-session' "$verify_output" || { echo "FAIL: healthy a1 session not recognized"; fail=1; }

for state in expired halfdead; do
    run_verify "$state" || { echo "FAIL: $state a1 session WARN should stay non-fatal"; fail=1; }
    grep -q 'WARN jarvis-a1-session' "$verify_output" || { echo "FAIL: missing a1 WARN for $state"; fail=1; }
    grep -q 'bin/a1id login jarvis' "$verify_output" || { echo "FAIL: missing repair guidance for $state"; fail=1; }
done

if [ "$fail" -eq 0 ]; then
    echo "verify_test: PASS"
    exit 0
fi
exit 1
