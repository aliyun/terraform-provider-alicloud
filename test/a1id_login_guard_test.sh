#!/usr/bin/env bash
# test/a1id_login_guard_test.sh - P0 护栏:bin/a1id login 后 whoami 必须匹配 label 的期望账号
#
# Run: bash test/a1id_login_guard_test.sh
# Prints PASS/FAIL per assertion; exits 0 on all-pass, 1 on any failure.
#
# Stub a1 二进制:
#   `a1 auth login --buc` 模拟 SSO,把 live auth.yaml 写成 $FAKE_A1_LOGIN_ACCOUNT 的凭据
#     (可通过 $FAKE_A1_LOGIN_STATUS 非零模拟登录失败)
#   `a1 auth whoami` 读 live auth.yaml,回显 `Account: <name>`

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
A1ID="$(cd "$SCRIPT_DIR/../bin" && pwd)/a1id"

if [ ! -f "$A1ID" ]; then
    echo "FAIL: missing script under test: $A1ID" >&2
    exit 1
fi

TMP_DIR="$(mktemp -d)"
FAKE_BIN_DIR="$TMP_DIR/bin"
A1_CFG_DIR="$TMP_DIR/a1cfg"
LAST_OUT="$TMP_DIR/out"
LAST_ERR="$TMP_DIR/err"
ORIGINAL_PATH="$PATH"

mkdir -p "$FAKE_BIN_DIR" "$A1_CFG_DIR"

cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

# --- 假 a1 二进制 ------------------------------------------------------------
cat > "$FAKE_BIN_DIR/a1" <<'A1_STUB'
#!/usr/bin/env bash
set -uo pipefail

if [ "${1:-}" = "auth" ] && [ "${2:-}" = "login" ]; then
    # 模拟 SSO 成功:把 live auth.yaml 写成 FAKE_A1_LOGIN_ACCOUNT 的凭据
    if [ "${FAKE_A1_LOGIN_STATUS:-0}" != "0" ]; then
        echo "fake a1 login failure" >&2
        exit "${FAKE_A1_LOGIN_STATUS}"
    fi
    mkdir -p "${A1_CONFIG_DIR:-$HOME/.config/a1}"
    cat > "${A1_CONFIG_DIR:-$HOME/.config/a1}/auth.yaml" <<EOF
version: 1
current:
    user:
        account: ${FAKE_A1_LOGIN_ACCOUNT:-WORKER_1782379562571}
        user: ${FAKE_A1_LOGIN_ACCOUNT:-WORKER_1782379562571}
EOF
    exit 0
fi

if [ "${1:-}" = "auth" ] && [ "${2:-}" = "whoami" ]; then
    # 显式让 whoami 输出无 Account: 行(模拟 a1 输出格式变化 / 未登录)
    if [ "${FAKE_A1_WHOAMI_EMPTY:-0}" = "1" ]; then
        echo "not logged in" >&2
        exit 1
    fi
    live="${A1_CONFIG_DIR:-$HOME/.config/a1}/auth.yaml"
    acct=""
    if [ -f "$live" ]; then
        acct=$(awk '/^ *account:/ {print $2; exit}' "$live" 2>/dev/null || echo "")
    fi
    if [ -z "$acct" ]; then
        echo "not logged in" >&2
        exit 1
    fi
    cat <<EOF
Account:  $acct
Name:     Fake User
Emp ID:   000000
Email:    $acct@example.com
Source:   fake
EOF
    exit 0
fi

echo "fake a1: unhandled args: $*" >&2
exit 64
A1_STUB
chmod +x "$FAKE_BIN_DIR/a1"

export PATH="$FAKE_BIN_DIR:$ORIGINAL_PATH"
export A1_CONFIG_DIR="$A1_CFG_DIR"

pass=0
fail=0

reset_env() {
    unset FAKE_A1_LOGIN_ACCOUNT FAKE_A1_LOGIN_STATUS FAKE_A1_WHOAMI_EMPTY
    rm -rf "$A1_CFG_DIR"
    mkdir -p "$A1_CFG_DIR/identities"
}

# 预置一个"已存的合法 jarvis 凭据"到 store 和 live,加 marker 便于验证未被覆盖。
# 也顺带绕过 a1id L42 的"首跑迁移"(store_for jarvis 已存在时该行 no-op)。
prime_existing_jarvis() {
    local marker="${1:-PRE_EXISTING_JARVIS}"
    cat > "$A1_CFG_DIR/identities/jarvis.auth.yaml" <<EOF
version: 1
current:
    user:
        account: WORKER_1782379562571
        user: WORKER_1782379562571
        marker: $marker
EOF
    cp "$A1_CFG_DIR/identities/jarvis.auth.yaml" "$A1_CFG_DIR/auth.yaml"
    echo jarvis > "$A1_CFG_DIR/identities/.active"
}

run_a1id() {
    : > "$LAST_OUT"
    : > "$LAST_ERR"
    bash "$A1ID" "$@" >"$LAST_OUT" 2>"$LAST_ERR"
    LAST_STATUS=$?
    LAST_STDOUT="$(cat "$LAST_OUT")"
    LAST_STDERR="$(cat "$LAST_ERR")"
}

assert_eq() {
    local desc="$1" want="$2" got="$3"
    if [ "$got" = "$want" ]; then echo "PASS $desc"; pass=$((pass + 1))
    else echo "FAIL $desc: got='$got' want='$want'"; fail=$((fail + 1)); fi
}
assert_ne() {
    local desc="$1" nope="$2" got="$3"
    if [ "$got" != "$nope" ]; then echo "PASS $desc"; pass=$((pass + 1))
    else echo "FAIL $desc: expected != '$nope'"; fail=$((fail + 1)); fi
}
assert_contains() {
    local desc="$1" haystack="$2" needle="$3"
    if printf '%s' "$haystack" | grep -Fq "$needle"; then echo "PASS $desc"; pass=$((pass + 1))
    else echo "FAIL $desc: missing '$needle' in '$haystack'"; fail=$((fail + 1)); fi
}
assert_not_contains() {
    local desc="$1" haystack="$2" needle="$3"
    if printf '%s' "$haystack" | grep -Fq "$needle"; then echo "FAIL $desc: unexpectedly contains '$needle'"; fail=$((fail + 1))
    else echo "PASS $desc"; pass=$((pass + 1)); fi
}
assert_file_exists() {
    local desc="$1" path="$2"
    if [ -f "$path" ]; then echo "PASS $desc"; pass=$((pass + 1))
    else echo "FAIL $desc: file missing: $path"; fail=$((fail + 1)); fi
}
assert_file_absent() {
    local desc="$1" path="$2"
    if [ ! -f "$path" ]; then echo "PASS $desc"; pass=$((pass + 1))
    else echo "FAIL $desc: file unexpectedly exists: $path"; fail=$((fail + 1)); fi
}
assert_file_contains() {
    local desc="$1" path="$2" needle="$3"
    if [ -f "$path" ] && grep -Fq "$needle" "$path"; then echo "PASS $desc"; pass=$((pass + 1))
    else echo "FAIL $desc: needle '$needle' missing from $path"; fail=$((fail + 1)); fi
}

echo "Test 1: happy path — 期望账号与 whoami 匹配,落盘成功"
reset_env
export FAKE_A1_LOGIN_ACCOUNT="WORKER_1782379562571"
run_a1id login jarvis
assert_eq "exit 0" "0" "$LAST_STATUS"
assert_contains "stdout 报成功并带账号" "$LAST_STDOUT" "'jarvis'(WORKER_1782379562571) 登录已保存"
assert_file_exists "jarvis.auth.yaml 已创建" "$A1_CFG_DIR/identities/jarvis.auth.yaml"
assert_file_contains "jarvis.auth.yaml 内容为 WORKER_1782379562571" "$A1_CFG_DIR/identities/jarvis.auth.yaml" "WORKER_1782379562571"
assert_eq ".active == jarvis" "jarvis" "$(cat "$A1_CFG_DIR/identities/.active" 2>/dev/null)"

echo ""
echo "Test 2: mismatch — 请求 jarvis 但 SSO 登了 guozai,不落盘 + 回滚 live"
reset_env
prime_existing_jarvis "TEST2_MARKER"
# 假 SSO:浏览器 BUC 会话是 guozai(与请求 label 期望不符)
export FAKE_A1_LOGIN_ACCOUNT="guozai.gzl"
run_a1id login jarvis
assert_ne "exit 非零" "0" "$LAST_STATUS"
assert_contains "stderr 报身份不匹配" "$LAST_STDERR" "身份不匹配"
assert_contains "stderr 显示期望账号 WORKER_1782379562571" "$LAST_STDERR" "WORKER_1782379562571"
assert_contains "stderr 显示实际账号 guozai.gzl" "$LAST_STDERR" "guozai.gzl"
assert_contains "stderr 给恢复步骤" "$LAST_STDERR" "登出 BUC"
# 核心不变量:jarvis.auth.yaml 未被 guozai 凭据覆盖(marker 保留)
assert_file_contains "jarvis.auth.yaml 保留原 marker(未被覆盖)" "$A1_CFG_DIR/identities/jarvis.auth.yaml" "TEST2_MARKER"
assert_not_contains "jarvis.auth.yaml 未泄漏 guozai.gzl" "$(cat "$A1_CFG_DIR/identities/jarvis.auth.yaml")" "guozai.gzl"
# live 也回滚干净
assert_file_contains "live auth.yaml 已回滚到 WORKER_1782379562571" "$A1_CFG_DIR/auth.yaml" "WORKER_1782379562571"
assert_not_contains "live auth.yaml 不含 guozai.gzl 泄漏" "$(cat "$A1_CFG_DIR/auth.yaml")" "guozai.gzl"

echo ""
echo "Test 3: 未知 label — 早期 die,不动 live"
reset_env
cat > "$A1_CFG_DIR/auth.yaml" <<EOF
version: 1
current:
    user:
        account: WORKER_1782379562571
        user: WORKER_1782379562571
EOF
run_a1id login unknown
assert_ne "exit 非零" "0" "$LAST_STATUS"
assert_contains "stderr 提示用法" "$LAST_STDERR" "用法"
assert_file_contains "live 未被动过" "$A1_CFG_DIR/auth.yaml" "WORKER_1782379562571"

echo ""
echo "Test 4: a1 auth login 本身失败 — trap 回滚 live,不落盘"
reset_env
prime_existing_jarvis "TEST4_MARKER"
export FAKE_A1_LOGIN_STATUS="42"
run_a1id login jarvis
assert_ne "exit 非零" "0" "$LAST_STATUS"
assert_file_contains "jarvis.auth.yaml marker 保留" "$A1_CFG_DIR/identities/jarvis.auth.yaml" "TEST4_MARKER"
assert_file_contains "live auth.yaml 保持原样(WORKER_1782379562571)" "$A1_CFG_DIR/auth.yaml" "WORKER_1782379562571"

echo ""
echo "Test 5: guozai 身份 happy path — 映射表覆盖到非默认身份"
reset_env
export FAKE_A1_LOGIN_ACCOUNT="guozai.gzl"
run_a1id login guozai
assert_eq "exit 0" "0" "$LAST_STATUS"
assert_file_exists "guozai.auth.yaml 已创建" "$A1_CFG_DIR/identities/guozai.auth.yaml"
assert_file_contains "guozai.auth.yaml 内容为 guozai.gzl" "$A1_CFG_DIR/identities/guozai.auth.yaml" "guozai.gzl"

echo ""
echo "Test 6: whoami 返空(a1 输出格式变化) — 视为不匹配,回滚"
reset_env
prime_existing_jarvis "TEST6_MARKER"
# SSO 部分正常,但 whoami 拿不到 Account: 行(模拟 a1 输出结构变了)
export FAKE_A1_LOGIN_ACCOUNT="WORKER_1782379562571"
export FAKE_A1_WHOAMI_EMPTY="1"
run_a1id login jarvis
assert_ne "exit 非零" "0" "$LAST_STATUS"
assert_contains "stderr 提到 <空>" "$LAST_STDERR" "<空>"
assert_file_contains "jarvis.auth.yaml marker 保留" "$A1_CFG_DIR/identities/jarvis.auth.yaml" "TEST6_MARKER"

echo ""
echo "Test 7: migration happy — live 是 WORKER_1782379562571 且 store 空,首跑迁移应成功"
reset_env
# 预置:live 是合法的 jarvis 会话,store 里还没有 jarvis.auth.yaml
cat > "$A1_CFG_DIR/auth.yaml" <<EOF
version: 1
current:
    user:
        account: WORKER_1782379562571
        user: WORKER_1782379562571
EOF
# 触发迁移(status 不 activate,只跑顶层迁移代码路径,再走到 case);
# 用 -- auth whoami 端到端跑一次,顺带验证迁移后 activate 能通
run_a1id -- auth whoami
assert_eq "exit 0" "0" "$LAST_STATUS"
assert_file_exists "jarvis.auth.yaml 迁移成功" "$A1_CFG_DIR/identities/jarvis.auth.yaml"
assert_file_contains "jarvis.auth.yaml 内容是 WORKER_1782379562571" "$A1_CFG_DIR/identities/jarvis.auth.yaml" "WORKER_1782379562571"
assert_eq ".active == jarvis" "jarvis" "$(cat "$A1_CFG_DIR/identities/.active" 2>/dev/null)"

echo ""
echo "Test 8: migration guard — live 是 guozai.gzl,首跑迁移应跳过 + warn"
reset_env
# 预置:live 是残留的 guozai 会话(实际就是本次踩坑的复现)
cat > "$A1_CFG_DIR/auth.yaml" <<EOF
version: 1
current:
    user:
        account: guozai.gzl
        user: guozai.gzl
EOF
# 任何 a1id 命令都会跑迁移;这里用 -- 端到端跑,期望 activate 因 store 空而 die
run_a1id -- auth whoami
assert_ne "exit 非零" "0" "$LAST_STATUS"
assert_contains "stderr 报跳过迁移" "$LAST_STDERR" "跳过首跑迁移"
assert_contains "stderr 显示 live 的实际账号 guozai.gzl" "$LAST_STDERR" "guozai.gzl"
assert_contains "stderr 追加提示要 login jarvis" "$LAST_STDERR" "bin/a1id login jarvis"
assert_contains "activate 后续 die 提示未登录" "$LAST_STDERR" "身份 'jarvis' 未登录"
assert_file_absent "jarvis.auth.yaml 未被 guozai 污染" "$A1_CFG_DIR/identities/jarvis.auth.yaml"
assert_not_contains "live auth.yaml 未被 activate 覆盖(还是 guozai)" "$(cat "$A1_CFG_DIR/auth.yaml")" "WORKER_1782379562571"

echo ""
echo "Test 9: migration 边界 — live 缺失,不迁移不 warn(纯静默)"
reset_env
# store 空 + live 也不存在:干净初始状态
run_a1id status
assert_eq "exit 0" "0" "$LAST_STATUS"
assert_file_absent "jarvis.auth.yaml 未被创建" "$A1_CFG_DIR/identities/jarvis.auth.yaml"
assert_not_contains "stderr 不应该有跳过迁移 warn" "$LAST_STDERR" "跳过首跑迁移"

echo ""
echo "Results: $pass passed, $fail failed"
if [ "$fail" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL"
    exit 1
fi
