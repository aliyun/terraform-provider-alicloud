#!/usr/bin/env bash
# test/github_identity_test.sh - unit tests for bootstrap/github-identity.sh
#
# Run: bash test/github_identity_test.sh
# Prints PASS/FAIL per assertion; exits 0 on all-pass, 1 on any failure.

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BOOTSTRAP_DIR="$(cd "$SCRIPT_DIR/../bootstrap" && pwd)"
GITHUB_IDENTITY="$BOOTSTRAP_DIR/github-identity.sh"

if [ ! -f "$GITHUB_IDENTITY" ]; then
    echo "FAIL: missing script under test: $GITHUB_IDENTITY" >&2
    exit 1
fi

TMP_DIR="$(mktemp -d)"
FAKE_BIN_DIR="$TMP_DIR/bin"
FAKE_GH_LOG="$TMP_DIR/gh.log"
LAST_OUT="$TMP_DIR/out"
LAST_ERR="$TMP_DIR/err"
ORIGINAL_PATH="$PATH"

mkdir -p "$FAKE_BIN_DIR"
: > "$FAKE_GH_LOG"

cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

cat > "$FAKE_BIN_DIR/gh" <<'GH_STUB'
#!/usr/bin/env bash
set -uo pipefail

{
    printf 'GH_TOKEN=%s\n' "${GH_TOKEN:-}"
    printf 'ARGS=%s\n' "$*"
} >> "$FAKE_GH_LOG"

if [ "$#" -ge 2 ] && [ "$1" = "api" ] && [ "$2" = "user" ]; then
    if [ "${FAKE_GH_API_FAIL:-}" = "1" ]; then
        echo "fake gh api failure" >&2
        exit "${FAKE_GH_API_STATUS:-42}"
    fi
    if [ "$*" != "api user --jq .login" ]; then
        echo "unexpected gh api args: $*" >&2
        exit 64
    fi
    printf '%s\n' "${FAKE_GH_LOGIN:-api-tool-agent}"
    exit 0
fi

if [ "${FAKE_GH_COMMAND_FAIL:-}" = "1" ]; then
    echo "fake gh command failure" >&2
    exit "${FAKE_GH_COMMAND_STATUS:-43}"
fi

printf 'fake gh ok: %s\n' "$*"
GH_STUB
chmod +x "$FAKE_BIN_DIR/gh"

cat > "$FAKE_BIN_DIR/git" <<'GIT_STUB'
#!/usr/bin/env bash
set -uo pipefail

{
    printf 'GIT_TERMINAL_PROMPT=%s\n' "${GIT_TERMINAL_PROMPT:-}"
    printf 'GIT_ARGS=%s\n' "$*"
    if [ -n "${GIT_ASKPASS:-}" ]; then
        printf 'ASKPASS_USER=%s\n' "$("$GIT_ASKPASS" "Username for https://github.com")"
        printf 'ASKPASS_PASS=%s\n' "$("$GIT_ASKPASS" "Password for https://github.com")"
    fi
} >> "$FAKE_GH_LOG"

if [ "${FAKE_GIT_FAIL:-}" = "1" ]; then
    echo "fake git failure" >&2
    exit "${FAKE_GIT_STATUS:-44}"
fi

exit 0
GIT_STUB
chmod +x "$FAKE_BIN_DIR/git"

export FAKE_GH_LOG
export PATH="$FAKE_BIN_DIR:$ORIGINAL_PATH"

pass=0
fail=0

reset_env() {
    : > "$FAKE_GH_LOG"
    unset JARVIS_GITHUB_TOKEN JARVIS_GITHUB_LOGIN
    unset FAKE_GH_LOGIN FAKE_GH_API_FAIL FAKE_GH_API_STATUS
    unset FAKE_GH_COMMAND_FAIL FAKE_GH_COMMAND_STATUS
    unset FAKE_GIT_FAIL FAKE_GIT_STATUS
}

run_cmd() {
    : > "$LAST_OUT"
    : > "$LAST_ERR"
    "$@" >"$LAST_OUT" 2>"$LAST_ERR"
    LAST_STATUS=$?
    LAST_STDOUT="$(cat "$LAST_OUT")"
    LAST_STDERR="$(cat "$LAST_ERR")"
}

assert_eq() {
    local desc="$1"
    local want="$2"
    local got="$3"
    if [ "$got" = "$want" ]; then
        echo "PASS $desc"
        pass=$((pass + 1))
    else
        echo "FAIL $desc: got='$got' want='$want'"
        fail=$((fail + 1))
    fi
}

assert_contains() {
    local desc="$1"
    local haystack="$2"
    local needle="$3"
    if printf '%s' "$haystack" | grep -Fq "$needle"; then
        echo "PASS $desc"
        pass=$((pass + 1))
    else
        echo "FAIL $desc: missing '$needle' in '$haystack'"
        fail=$((fail + 1))
    fi
}

assert_nonzero() {
    local desc="$1"
    local got="$2"
    if [ "$got" -ne 0 ]; then
        echo "PASS $desc"
        pass=$((pass + 1))
    else
        echo "FAIL $desc: expected non-zero exit"
        fail=$((fail + 1))
    fi
}

echo "Test 1: check requires JARVIS_GITHUB_TOKEN"
reset_env
run_cmd bash "$GITHUB_IDENTITY" check
assert_nonzero "missing token exits non-zero" "$LAST_STATUS"
assert_contains "missing token explains env var" "$LAST_STDERR" "JARVIS_GITHUB_TOKEN"

echo "Test 2: check uses token and default login"
reset_env
export JARVIS_GITHUB_TOKEN="token-default"
export FAKE_GH_LOGIN="api-tool-agent"
run_cmd bash "$GITHUB_IDENTITY" check
assert_eq "default check exits 0" "0" "$LAST_STATUS"
assert_contains "check exports token to gh" "$(cat "$FAKE_GH_LOG")" "GH_TOKEN=token-default"
assert_contains "check calls gh api user" "$(cat "$FAKE_GH_LOG")" "ARGS=api user --jq .login"

echo "Test 3: check rejects JARVIS_GITHUB_LOGIN override attempts"
reset_env
export JARVIS_GITHUB_TOKEN="token-custom"
export JARVIS_GITHUB_LOGIN="custom-agent"
export FAKE_GH_LOGIN="custom-agent"
run_cmd bash "$GITHUB_IDENTITY" check
assert_nonzero "custom login override exits non-zero" "$LAST_STATUS"
assert_contains "custom login override still expects api-tool-agent" "$LAST_STDERR" "api-tool-agent"

echo "Test 4: check rejects login mismatch"
reset_env
export JARVIS_GITHUB_TOKEN="token-mismatch"
export FAKE_GH_LOGIN="wrong-agent"
run_cmd bash "$GITHUB_IDENTITY" check
assert_nonzero "login mismatch exits non-zero" "$LAST_STATUS"
assert_contains "login mismatch names expected login" "$LAST_STDERR" "api-tool-agent"
assert_contains "login mismatch names actual login" "$LAST_STDERR" "wrong-agent"

echo "Test 5: check reports gh api failure"
reset_env
export JARVIS_GITHUB_TOKEN="token-api-fail"
export FAKE_GH_API_FAIL="1"
run_cmd bash "$GITHUB_IDENTITY" check
assert_nonzero "gh api failure exits non-zero" "$LAST_STATUS"
assert_contains "gh api failure message is clear" "$LAST_STDERR" "gh api user failed"
assert_contains "gh api stderr is preserved" "$LAST_STDERR" "fake gh api failure"

echo "Test 6: gh mode checks identity then forwards command with token"
reset_env
export JARVIS_GITHUB_TOKEN="token-forward"
export FAKE_GH_LOGIN="api-tool-agent"
run_cmd bash "$GITHUB_IDENTITY" gh pr view 123 --json url
assert_eq "gh mode exits 0" "0" "$LAST_STATUS"
assert_contains "gh mode forwards command" "$(cat "$FAKE_GH_LOG")" "ARGS=pr view 123 --json url"
assert_contains "gh mode uses token for forwarded command" "$(cat "$FAKE_GH_LOG")" "GH_TOKEN=token-forward"

echo "Test 7: push mode checks identity and uses token-backed askpass"
reset_env
export JARVIS_GITHUB_TOKEN="token-push"
export FAKE_GH_LOGIN="api-tool-agent"
run_cmd bash "$GITHUB_IDENTITY" push api-tool-agent/terraform-provider-alicloud HEAD feature/branch
assert_eq "push mode exits 0" "0" "$LAST_STATUS"
assert_contains "push mode invokes git push" "$(cat "$FAKE_GH_LOG")" "GIT_ARGS=push https://github.com/api-tool-agent/terraform-provider-alicloud.git HEAD:feature/branch"
assert_contains "push askpass supplies username" "$(cat "$FAKE_GH_LOG")" "ASKPASS_USER=x-access-token"
assert_contains "push askpass supplies token" "$(cat "$FAKE_GH_LOG")" "ASKPASS_PASS=token-push"

echo "Test 8: push mode rejects invalid repo"
reset_env
export JARVIS_GITHUB_TOKEN="token-push-invalid"
export FAKE_GH_LOGIN="api-tool-agent"
run_cmd bash "$GITHUB_IDENTITY" push https://github.com/api-tool-agent/terraform-provider-alicloud HEAD feature/branch
assert_nonzero "invalid repo exits non-zero" "$LAST_STATUS"
assert_contains "invalid repo message is clear" "$LAST_STDERR" "owner/repo"

echo "Test 9: gh mode reports forwarded gh failure"
reset_env
export JARVIS_GITHUB_TOKEN="token-command-fail"
export FAKE_GH_LOGIN="api-tool-agent"
export FAKE_GH_COMMAND_FAIL="1"
export FAKE_GH_COMMAND_STATUS="43"
run_cmd bash "$GITHUB_IDENTITY" gh pr create --draft
assert_eq "forwarded gh status is preserved" "43" "$LAST_STATUS"
assert_contains "forwarded gh stderr is preserved" "$LAST_STDERR" "fake gh command failure"
assert_contains "forwarded gh failure message is clear" "$LAST_STDERR" "gh command failed"

echo ""
echo "Results: $pass passed, $fail failed"
if [ "$fail" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    exit 1
fi
