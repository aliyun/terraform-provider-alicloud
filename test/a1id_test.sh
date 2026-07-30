#!/usr/bin/env bash
# test/a1id_test.sh — TDD for bin/a1id v2 并发多身份切换器
#
# 所有测试通过 A1ID_ROOT=<tmp> 隔离,绝不触碰真实 ~/.config/a1;
# A1_BIN 指向一个记录 argv+A1_CONFIG_DIR 的 stub。
#
# 用例覆盖(见 SUMMARY 一节):
#   1  旧 terraform-pd label 仅作为兼容别名命中 terraform-rd dir并告警
#   2  as <label> 未登录 → die + 报 'login <id>'
#   3  JARVIS_A1_IDENTITY=<label> 且已登录 → -- 走该 dir
#   4  旧 terraform-qa auth 不再命中；公共 terraform-rd 未登录时禁止回退 jarvis
#   4b strict 同样 die
#   5  v1→v2 迁移:identities/<label>.auth.yaml 存在 → 首跑生成 identities/<label>/auth.yaml
#   6  ready 已/未登录退码 0/1
#   7  别名:pd/qa/terraform-pd/terraform-qa 全部收口到 terraform-rd
#   8  login 账号不匹配 → die 且清 auth.yaml
#   9  并发冒烟:两 as 不同身份后台并跑无串扰
#   10 JARVIS_A1_IDENTITY=guozai/shanye(个人身份)→ 警告但可执行

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"
A1ID="$proj_root/bin/a1id"

if [ ! -x "$A1ID" ]; then
    echo "FATAL: $A1ID 不可执行" >&2
    exit 1
fi

PASS=0; FAIL=0
pass(){ echo "PASS: $1"; PASS=$((PASS+1)); }
fail(){ echo "FAIL: $1"; FAIL=$((FAIL+1)); }

# 每 case 一个新根,彻底隔离
new_root(){ mktemp -d; }

# stub a1 二进制:
#   auth login  → 在 A1_CONFIG_DIR 内写 auth.yaml(内容 account=STUB_WHOAMI_ACCOUNT);
#                 STUB_LOGIN_FAIL=1 → 不写 auth.yaml 且退 42(模拟 a1 login 自身失败)
#   auth whoami → 打印 Account: <STUB_WHOAMI_ACCOUNT>;STUB_WHOAMI_EMPTY=1 → 退非零无输出
#   其他        → 把 A1_CONFIG_DIR + argv 一行写入 STUB_CAPTURE(TAB 分隔)
make_stub(){
    local tmpbin="$1"
    cat > "$tmpbin/a1" <<'STUB'
#!/usr/bin/env bash
if [ -n "${STUB_CAPTURE:-}" ]; then
    printf 'A1_CONFIG_DIR=%s\tARGS=%s\n' "${A1_CONFIG_DIR:-}" "$*" >> "$STUB_CAPTURE"
fi
if [ "${1:-}" = "auth" ] && [ "${2:-}" = "login" ]; then
    if [ "${STUB_LOGIN_FAIL:-0}" = "1" ]; then
        echo "stub a1 auth login: simulated failure" >&2
        exit 42
    fi
    mkdir -p "${A1_CONFIG_DIR:-.}"
    cat > "${A1_CONFIG_DIR}/auth.yaml" <<EOF
version: 1
current:
    user:
        account: ${STUB_WHOAMI_ACCOUNT:-unknown}
        user: ${STUB_WHOAMI_ACCOUNT:-unknown}
EOF
    exit 0
fi
if [ "${1:-}" = "auth" ] && [ "${2:-}" = "whoami" ]; then
    if [ "${STUB_WHOAMI_EMPTY:-0}" = "1" ]; then
        echo "not logged in" >&2
        exit 1
    fi
    cat <<EOF
Account:  ${STUB_WHOAMI_ACCOUNT:-unknown}
Name:     stub
Emp ID:   000000
Email:    stub@example.com
EOF
    exit 0
fi
exit 0
STUB
    chmod +x "$tmpbin/a1"
}

# 预置一个身份为「已登录」(直接写 auth.yaml,不走 a1 login)
seed_login(){
    local root="$1" label="$2"
    mkdir -p "$root/identities/$label"
    cat > "$root/identities/$label/auth.yaml" <<EOF
version: 1
current:
    user:
        account: seeded-$label
EOF
}

# ===========================================================================
# Test 1: as terraform-pd -- <args> 只走 rd 目录并告警,argv 正确
# ===========================================================================
echo "=== Test 1: legacy terraform-pd → terraform-rd dir + warning ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "terraform-rd"
CAP=$(mktemp); ERR=$(mktemp)
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    bash "$A1ID" as terraform-pd -- project foo bar >/dev/null 2>"$ERR"
rc=$?
line=$(cat "$CAP")
[ "$rc" = "0" ] && pass "as terraform-pd 退 0" || fail "as terraform-pd exit=$rc"
if grep -qF "A1_CONFIG_DIR=$ROOT/identities/terraform-rd	" "$CAP"; then
    pass "legacy terraform-pd 只命中 terraform-rd dir"
else
    fail "legacy terraform-pd 未收口到 terraform-rd: $line"
fi
grep -qF "兼容别名" "$ERR" && pass "legacy terraform-pd 发出兼容别名告警" \
    || fail "legacy terraform-pd 缺兼容告警: $(cat "$ERR")"
if grep -qF "ARGS=project foo bar" "$CAP"; then
    pass "stub 收到 argv 完整"
else
    fail "argv 传错: $line"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 2: as terraform-rd 未登录 → die + 报错含 'login terraform-rd'
# ===========================================================================
echo "=== Test 2: as terraform-rd 未登录 → die ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
CAP=$(mktemp); ERR=$(mktemp)
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    bash "$A1ID" as terraform-rd -- x >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" != "0" ] && pass "未登录 as 退非零" || fail "未登录 as 应退非零, got=$rc"
if printf '%s' "$err" | grep -qF "login terraform-rd"; then
    pass "错误消息含修复指引 'login terraform-rd'"
else
    fail "错误消息缺 'login terraform-rd': $err"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 3: JARVIS_A1_IDENTITY=terraform-qa 兼容映射到 rd dir
# ===========================================================================
echo "=== Test 3: JARVIS_A1_IDENTITY=terraform-qa → terraform-rd dir ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "terraform-rd"
seed_login "$ROOT" "jarvis"
CAP=$(mktemp); ERR=$(mktemp)
JARVIS_A1_IDENTITY=terraform-qa A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" -- y >/dev/null 2>"$ERR"
rc=$?
[ "$rc" = "0" ] && pass "-- with qa 退 0" || fail "exit=$rc"
if grep -qF "A1_CONFIG_DIR=$ROOT/identities/terraform-rd	" "$CAP"; then
    pass "terraform-qa env 兼容别名 → rd dir"
else
    fail "terraform-qa env 未走 rd 目录: $(cat "$CAP")"
fi
grep -qF "兼容别名" "$ERR" && pass "terraform-qa env 发出兼容别名告警" \
    || fail "terraform-qa env 缺兼容告警: $(cat "$ERR")"
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 4: 只有旧 qa auth + jarvis auth、无 rd auth → die,不得命中旧 auth/回退 jarvis
# ===========================================================================
echo "=== Test 4: legacy qa auth 不命中,公共 rd 缺登录时不回退 jarvis ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "jarvis"
seed_login "$ROOT" "terraform-qa"
CAP=$(mktemp); ERR=$(mktemp)
JARVIS_A1_IDENTITY=terraform-qa A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" -- z >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" != "0" ] && pass "公共 rd 未登录时拒绝外写" || fail "应拒绝, exit=$rc"
if printf '%s' "$err" | grep -qF "terraform-rd"; then
    pass "stderr 指向唯一公共身份 terraform-rd"
else
    fail "stderr 未指向 terraform-rd: $err"
fi
if [ ! -s "$CAP" ]; then
    pass "未命中旧 qa auth、未回退 jarvis、未 exec a1"
else
    fail "不应 exec a1: $(cat "$CAP")"
fi

echo "=== Test 4b: 同场景 STRICT=1 → die,不回退 ==="
: > "$CAP"; : > "$ERR"
JARVIS_A1_IDENTITY=terraform-qa JARVIS_A1_STRICT=1 A1ID_ROOT="$ROOT" \
    A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" bash "$A1ID" -- z >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" != "0" ] && pass "STRICT=1 未登录 die" || fail "STRICT=1 应 die, exit=$rc"
if printf '%s' "$err" | grep -qF "STRICT" || printf '%s' "$err" | grep -qF "未登录"; then
    pass "STRICT 错误消息合理"
else
    fail "STRICT 错误缺关键字: $err"
fi
if [ ! -s "$CAP" ]; then
    pass "STRICT die 后未 exec a1(capture 空)"
else
    fail "STRICT die 后不应 exec a1, capture=$(cat "$CAP")"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 5: v1 → v2 迁移
# ===========================================================================
echo "=== Test 5: v1→v2 迁移(jarvis) ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
mkdir -p "$ROOT/identities"
cat > "$ROOT/identities/jarvis.auth.yaml" <<EOF
version: 1
current:
    user:
        account: WORKER_1782379562571
EOF
CAP=$(mktemp); ERR=$(mktemp)
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    bash "$A1ID" status >/dev/null 2>"$ERR"
rc=$?
[ -f "$ROOT/identities/jarvis/auth.yaml" ] && pass "v2 dir 已生成" || fail "v2 dir 未生成"
[ -f "$ROOT/identities/jarvis.auth.yaml" ] && pass "v1 文件保留(未删)" || fail "v1 文件被删了"
if diff "$ROOT/identities/jarvis.auth.yaml" "$ROOT/identities/jarvis/auth.yaml" >/dev/null 2>&1; then
    pass "v1/v2 内容一致(cp 而非 mv)"
else
    fail "v1/v2 内容不一致"
fi
err=$(cat "$ERR")
if printf '%s' "$err" | grep -qF "v1→v2 迁移"; then
    pass "stderr 有迁移提示"
else
    fail "stderr 缺迁移提示: $err"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 6: ready 两态
# ===========================================================================
echo "=== Test 6: ready 已/未登录退码 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "jarvis"
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" bash "$A1ID" ready jarvis >/dev/null 2>&1
rc=$?
[ "$rc" = "0" ] && pass "ready jarvis(已登录)退 0" || fail "ready jarvis 应退 0, got=$rc"
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" bash "$A1ID" ready terraform-pd >/dev/null 2>&1
rc=$?
[ "$rc" != "0" ] && pass "ready terraform-pd(未登录)退非零" || fail "ready terraform-pd 应退非零, got=$rc"
rm -rf "$ROOT" "$BIN"

# ===========================================================================
# Test 7: 短别名 pd → terraform-rd
# ===========================================================================
echo "=== Test 7: 别名 pd → terraform-rd ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "terraform-rd"
CAP=$(mktemp); ERR=$(mktemp)
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    bash "$A1ID" as pd -- z >/dev/null 2>"$ERR"
rc=$?
[ "$rc" = "0" ] && pass "as pd 退 0" || fail "as pd exit=$rc"
if grep -qF "A1_CONFIG_DIR=$ROOT/identities/terraform-rd	" "$CAP"; then
    pass "别名 pd 映射到 terraform-rd dir"
else
    fail "pd 别名未生效: $(cat "$CAP")"
fi
grep -qF "兼容别名" "$ERR" && pass "pd 发出兼容告警" || fail "pd 缺兼容告警"
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 8: login 账号不匹配 → die + 清 auth.yaml
# ===========================================================================
echo "=== Test 8: login jarvis 但浏览器 BUC = guozai → die + 清盘 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
CAP=$(mktemp); ERR=$(mktemp)
STUB_WHOAMI_ACCOUNT="guozai.gzl" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" login jarvis >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" != "0" ] && pass "不匹配 login 退非零" || fail "不匹配 login 应退非零 rc=$rc"
if printf '%s' "$err" | grep -qF "身份不匹配"; then
    pass "报错含 '身份不匹配'"
else
    fail "缺 '身份不匹配': $err"
fi
if [ ! -f "$ROOT/identities/jarvis/auth.yaml" ]; then
    pass "不匹配后该身份 auth.yaml 已清(未污染)"
else
    fail "不匹配后 auth.yaml 未清: $(cat "$ROOT/identities/jarvis/auth.yaml")"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 9: 并发冒烟 — legacy pd/rd 并行都只读同一 rd auth,argv 无串扰
# ===========================================================================
echo "=== Test 9: 并发冒烟(pd + rd 后台并行) ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "terraform-rd"
CAP_PD=$(mktemp); CAP_RD=$(mktemp)
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP_PD" \
    bash "$A1ID" as terraform-pd -- concurrent-pd >/dev/null 2>&1 &
pid_pd=$!
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP_RD" \
    bash "$A1ID" as terraform-rd -- concurrent-rd >/dev/null 2>&1 &
pid_rd=$!
wait $pid_pd; rc_pd=$?
wait $pid_rd; rc_rd=$?
if [ "$rc_pd" = "0" ] && [ "$rc_rd" = "0" ]; then
    pass "两并发进程都退 0"
else
    fail "并发退码错 pd=$rc_pd rd=$rc_rd"
fi
if grep -qF "A1_CONFIG_DIR=$ROOT/identities/terraform-rd	" "$CAP_PD" \
    && grep -qF "concurrent-pd" "$CAP_PD"; then
    pass "legacy pd capture 使用 rd auth"
else
    fail "pd capture 错: $(cat "$CAP_PD")"
fi
if grep -qF "A1_CONFIG_DIR=$ROOT/identities/terraform-rd	" "$CAP_RD" \
    && grep -qF "concurrent-rd" "$CAP_RD"; then
    pass "rd capture 正确"
else
    fail "rd capture 错: $(cat "$CAP_RD")"
fi
if ! grep -qF "concurrent-rd" "$CAP_PD" && ! grep -qF "concurrent-pd" "$CAP_RD"; then
    pass "两并发进程 capture 无串扰"
else
    fail "并发串扰: pd=$(cat "$CAP_PD"); rd=$(cat "$CAP_RD")"
fi
rm -rf "$ROOT" "$BIN" "$CAP_PD" "$CAP_RD"

# ===========================================================================
# Test 10a: status 在无任何身份登录时仍退 0 且显示表格
#           (回归护栏:resolve_default die 不得短路 status 的 || echo 兜底)
# ===========================================================================
echo "=== Test 10a: status 无登录退 0(die 不短路 || echo) ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
OUT=$(mktemp)
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" bash "$A1ID" status >"$OUT" 2>&1
rc=$?
out=$(cat "$OUT")
[ "$rc" = "0" ] && pass "status(无登录)退 0" || fail "status 应退 0, got=$rc; out=$out"
if printf '%s' "$out" | grep -qF "身份表" && printf '%s' "$out" | grep -qF "terraform-rd" \
    && ! printf '%s' "$out" | grep -Eq '^  terraform-(pd|qa) '; then
    pass "status 只登记 terraform-rd 公共身份"
else
    fail "status 输出缺表格: $out"
fi
rm -rf "$ROOT" "$BIN" "$OUT"

# ===========================================================================
# Test 10: JARVIS_A1_IDENTITY=guozai(个人身份,已登录)→ 警告但可执行
# ===========================================================================
echo "=== Test 10: JARVIS_A1_IDENTITY=guozai → 个人身份纪律告警 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "guozai"
CAP=$(mktemp); ERR=$(mktemp)
JARVIS_A1_IDENTITY=guozai A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" -- personal-run >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" = "0" ] && pass "个人身份可执行(退 0)" || fail "个人身份执行失败 rc=$rc"
if printf '%s' "$err" | grep -qF "个人身份"; then
    pass "stderr 有个人身份纪律告警"
else
    fail "stderr 缺个人身份告警: $err"
fi
if grep -qF "A1_CONFIG_DIR=$ROOT/identities/guozai	" "$CAP"; then
    pass "确实以 guozai dir 跑"
else
    fail "capture 未走 guozai: $(cat "$CAP")"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 10b: shanye 完整登记 — env 个人身份告警 + 独立目录 + status 注册表
# ===========================================================================
echo "=== Test 10b: JARVIS_A1_IDENTITY=shanye → 个人身份纪律告警 + 独立目录 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "shanye"
CAP=$(mktemp); ERR=$(mktemp); OUT=$(mktemp)
JARVIS_A1_IDENTITY=shanye A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" -- personal-run >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" = "0" ] && pass "shanye 个人身份可执行(退 0)" || fail "shanye 个人身份执行失败 rc=$rc"
if printf '%s' "$err" | grep -qF "个人身份"; then
    pass "shanye stderr 有个人身份纪律告警"
else
    fail "shanye stderr 缺个人身份告警: $err"
fi
if grep -qF "A1_CONFIG_DIR=$ROOT/identities/shanye	" "$CAP"; then
    pass "shanye 使用独立 identities/shanye 目录"
else
    fail "capture 未走 shanye 独立目录: $(cat "$CAP")"
fi
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" bash "$A1ID" status >"$OUT" 2>&1
out=$(cat "$OUT")
if printf '%s' "$out" | grep -Eq '^  shanye +shanye\.xzq +yes$'; then
    pass "status 身份表包含 shanye 及期望 BUC 账号"
else
    fail "status 身份表缺 shanye 登记: $out"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR" "$OUT"

# ===========================================================================
# B2 移植:login 生命周期护栏(从被删的 a1id_login_guard_test.sh 搬到 v2 布局)
# ===========================================================================

# ---------------------------------------------------------------------------
# Test 11: login happy path — whoami 匹配期望账号 → auth.yaml 落盘 + 退 0
# ---------------------------------------------------------------------------
echo "=== Test 11: login happy path(whoami == 期望账号)==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
CAP=$(mktemp); ERR=$(mktemp)
STUB_WHOAMI_ACCOUNT="WORKER_1782379562571" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" login jarvis >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" = "0" ] && pass "login jarvis 匹配期望账号 → 退 0" || fail "login exit=$rc err=$err"
if [ -s "$ROOT/identities/jarvis/auth.yaml" ]; then
    pass "login 成功后 identities/jarvis/auth.yaml 已落盘"
else
    fail "login 成功但 auth.yaml 未落盘"
fi
if grep -qF "WORKER_1782379562571" "$ROOT/identities/jarvis/auth.yaml"; then
    pass "auth.yaml 内容含期望账号"
else
    fail "auth.yaml 内容错: $(cat "$ROOT/identities/jarvis/auth.yaml")"
fi
# 成功路径必须清 trap(否则退出时会误删)——间接验证:紧跟一条 ready 应仍成立
A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" bash "$A1ID" ready jarvis >/dev/null 2>&1
[ $? -eq 0 ] && pass "登录成功后 ready jarvis 退 0(trap 已清 + 未误删)" \
             || fail "登录成功后 ready jarvis 应 0"
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ---------------------------------------------------------------------------
# Test 11b: 数字人 worker 两种 whoami 呈现都可登录(账号名 / empId),错账号仍拒
# (实测:新开 worker 的 whoami Account=账号名 'terraform-rd',非 WORKER_empId)
# ---------------------------------------------------------------------------
echo "=== Test 11b: login 账号集合匹配(账号名/empId 双形态)==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
CAP=$(mktemp); ERR=$(mktemp)
STUB_WHOAMI_ACCOUNT="terraform-rd" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" login terraform-rd >/dev/null 2>"$ERR"
[ $? -eq 0 ] && pass "whoami=账号名 'terraform-rd' → login 放行" \
             || fail "账号名形态应放行 err=$(cat "$ERR")"
[ -s "$ROOT/identities/terraform-rd/auth.yaml" ] \
    && pass "账号名形态 auth.yaml 已落盘" || fail "账号名形态 auth.yaml 未落盘"
rm -rf "$ROOT/identities/terraform-rd"
STUB_WHOAMI_ACCOUNT="WORKER_1783582458263" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" login terraform-rd >/dev/null 2>"$ERR"
[ $? -eq 0 ] && pass "whoami=empId 'WORKER_1783582458263' → login 放行" \
             || fail "empId 形态应放行 err=$(cat "$ERR")"
STUB_WHOAMI_ACCOUNT="guozai.gzl" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" login terraform-rd >/dev/null 2>"$ERR"
rc=$?
[ "$rc" != "0" ] && pass "whoami=别人账号 → 仍拒(退非零)" || fail "错账号应拒 got=$rc"
grep -qF "身份不匹配" "$ERR" && pass "错账号 stderr 报'身份不匹配'" \
                             || fail "错账号应报身份不匹配: $(cat "$ERR")"
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ---------------------------------------------------------------------------
# Test 12: login mismatch 时既有凭据不被破坏(marker 内容原样保留)
# ---------------------------------------------------------------------------
echo "=== Test 12: mismatch 时既有 auth.yaml 被 EXIT trap 回滚(marker 保留)==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
mkdir -p "$ROOT/identities/jarvis"
cat > "$ROOT/identities/jarvis/auth.yaml" <<EOF
version: 1
current:
    user:
        account: WORKER_1782379562571
        marker: PRE_EXISTING_JARVIS_MARKER
EOF
CAP=$(mktemp); ERR=$(mktemp)
# 浏览器 BUC 会话是 guozai(与请求 label 期望不符)—— login 触发 mismatch die + trap 回滚
STUB_WHOAMI_ACCOUNT="guozai.gzl" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" login jarvis >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" != "0" ] && pass "mismatch login 退非零" || fail "mismatch login 应退非零 got=$rc"
if printf '%s' "$err" | grep -qF "身份不匹配"; then
    pass "stderr 报'身份不匹配'"
else
    fail "stderr 缺'身份不匹配': $err"
fi
# 核心不变量:预置的 marker 仍在(EXIT trap 从备份回滚)
if grep -qF "PRE_EXISTING_JARVIS_MARKER" "$ROOT/identities/jarvis/auth.yaml"; then
    pass "既有 auth.yaml 的 marker 已由 trap 回滚保留"
else
    fail "既有凭据被覆盖(trap 回滚失效): $(cat "$ROOT/identities/jarvis/auth.yaml" 2>/dev/null)"
fi
# 且不含污染的 guozai.gzl 账号(login 阶段 stub 写入的内容)
if ! grep -qF "guozai.gzl" "$ROOT/identities/jarvis/auth.yaml"; then
    pass "回滚后不含污染 guozai.gzl"
else
    fail "既有 auth.yaml 泄漏了 guozai.gzl: $(cat "$ROOT/identities/jarvis/auth.yaml")"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ---------------------------------------------------------------------------
# Test 13: a1 auth login 自身失败(退非零)→ trap 回滚 + 不落盘
# ---------------------------------------------------------------------------
echo "=== Test 13: a1 auth login 退非零 → trap 回滚 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
mkdir -p "$ROOT/identities/jarvis"
cat > "$ROOT/identities/jarvis/auth.yaml" <<EOF
version: 1
current:
    user:
        account: WORKER_1782379562571
        marker: PRE_TEST13_MARKER
EOF
CAP=$(mktemp); ERR=$(mktemp)
# STUB_LOGIN_FAIL=1 让 a1 auth login 直接退 42;a1id 不再走到 whoami
STUB_LOGIN_FAIL=1 STUB_WHOAMI_ACCOUNT="WORKER_1782379562571" A1ID_ROOT="$ROOT" \
    A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" bash "$A1ID" login jarvis >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" != "0" ] && pass "a1 login 失败 → login 命令退非零" || fail "应退非零 got=$rc"
# 核心不变量:预置 marker 仍在(trap 从备份回滚)
if grep -qF "PRE_TEST13_MARKER" "$ROOT/identities/jarvis/auth.yaml"; then
    pass "a1 login 失败后既有 marker 保留(trap 回滚)"
else
    fail "既有凭据被破坏: $(cat "$ROOT/identities/jarvis/auth.yaml" 2>/dev/null)"
fi
# 无备份场景:预置删除,再触发 login fail,断言 dir 空(未落盘)
rm -rf "$ROOT/identities/jarvis"
STUB_LOGIN_FAIL=1 STUB_WHOAMI_ACCOUNT="WORKER_1782379562571" A1ID_ROOT="$ROOT" \
    A1_BIN="$BIN/a1" bash "$A1ID" login jarvis >/dev/null 2>&1
if [ ! -f "$ROOT/identities/jarvis/auth.yaml" ]; then
    pass "无备份时 a1 login 失败 → 不落盘(trap 清理)"
else
    fail "无备份时不应落盘: $(cat "$ROOT/identities/jarvis/auth.yaml")"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ---------------------------------------------------------------------------
# Test 14: whoami 返空视为不匹配 → 同样触发 trap 回滚
# ---------------------------------------------------------------------------
echo "=== Test 14: whoami 返空 → 判不匹配 + trap 回滚 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
mkdir -p "$ROOT/identities/jarvis"
cat > "$ROOT/identities/jarvis/auth.yaml" <<EOF
version: 1
current:
    user:
        account: WORKER_1782379562571
        marker: PRE_TEST14_MARKER
EOF
CAP=$(mktemp); ERR=$(mktemp)
STUB_WHOAMI_EMPTY=1 STUB_WHOAMI_ACCOUNT="WORKER_1782379562571" A1ID_ROOT="$ROOT" \
    A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" bash "$A1ID" login jarvis >/dev/null 2>"$ERR"
rc=$?
err=$(cat "$ERR")
[ "$rc" != "0" ] && pass "whoami 返空 → login 退非零" || fail "应退非零 got=$rc"
if printf '%s' "$err" | grep -qF "<空>"; then
    pass "stderr 提到 <空>(视为不匹配)"
else
    fail "stderr 缺 <空>: $err"
fi
if grep -qF "PRE_TEST14_MARKER" "$ROOT/identities/jarvis/auth.yaml"; then
    pass "whoami 返空 → 既有 marker 保留(trap 回滚)"
else
    fail "既有凭据被 whoami-empty 路径破坏: $(cat "$ROOT/identities/jarvis/auth.yaml" 2>/dev/null)"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# B2 移植:live 首跑收编门四件套(A1ID_ROOT/auth.yaml 是 live)
# ===========================================================================

# ---------------------------------------------------------------------------
# Test 15: live whoami == jarvis 期望账号 → 收编进 identities/jarvis/auth.yaml
# ---------------------------------------------------------------------------
echo "=== Test 15: live 收编 — whoami 匹配 jarvis 期望 → 收编 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
# 预置 live(=$A1ID_ROOT/auth.yaml)含带 marker 的 yaml
cat > "$ROOT/auth.yaml" <<EOF
version: 1
current:
    user:
        account: open_jarvis
        marker: LIVE_JARVIS_CONTENT
EOF
CAP=$(mktemp); ERR=$(mktemp)
# stub whoami 声称是 jarvis 期望账号(BUC 视角 vs yaml account 字段的 alias 不对称)
STUB_WHOAMI_ACCOUNT="WORKER_1782379562571" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" status >/dev/null 2>"$ERR"
err=$(cat "$ERR")
if [ -s "$ROOT/identities/jarvis/auth.yaml" ]; then
    pass "收编成功:identities/jarvis/auth.yaml 已生成"
else
    fail "未收编: $err"
fi
if grep -qF "LIVE_JARVIS_CONTENT" "$ROOT/identities/jarvis/auth.yaml" 2>/dev/null; then
    pass "收编内容 = live 原样(marker 保留)"
else
    fail "收编内容错: $(cat "$ROOT/identities/jarvis/auth.yaml" 2>/dev/null)"
fi
if printf '%s' "$err" | grep -qF "首跑收编"; then
    pass "stderr 有收编提示"
else
    fail "stderr 缺收编提示: $err"
fi
if [ "$(cat "$ROOT/identities/.active" 2>/dev/null)" = "jarvis" ]; then
    pass "收编后 .active = jarvis"
else
    fail ".active 未设 jarvis(got=$(cat "$ROOT/identities/.active" 2>/dev/null))"
fi
# S1 锚定验证:收编 whoami 时 stub 必须收到 A1_CONFIG_DIR=$A1ID_ROOT(与 live 同源)
if grep -qF "A1_CONFIG_DIR=$ROOT	ARGS=auth whoami" "$CAP"; then
    pass "收编 whoami 锚定 A1_CONFIG_DIR=\$A1ID_ROOT(S1 修复生效)"
else
    fail "收编 whoami 未锚定 A1ID_ROOT: $(cat "$CAP")"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ---------------------------------------------------------------------------
# Test 16: live whoami 是别人 → 跳过收编 + stderr 提示
# ---------------------------------------------------------------------------
echo "=== Test 16: live 收编 — whoami 非 jarvis 期望 → 跳过 + 提示 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
cat > "$ROOT/auth.yaml" <<EOF
version: 1
current:
    user:
        account: guozai.gzl
EOF
CAP=$(mktemp); ERR=$(mktemp)
STUB_WHOAMI_ACCOUNT="guozai.gzl" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" status >/dev/null 2>"$ERR"
err=$(cat "$ERR")
if [ ! -f "$ROOT/identities/jarvis/auth.yaml" ]; then
    pass "非 jarvis live → 不收编,identities/jarvis/auth.yaml 未生成"
else
    fail "错误收编了非 jarvis live: $(cat "$ROOT/identities/jarvis/auth.yaml")"
fi
if printf '%s' "$err" | grep -qF "跳过首跑迁移"; then
    pass "stderr 有'跳过首跑迁移'提示"
else
    fail "stderr 缺'跳过首跑迁移': $err"
fi
if printf '%s' "$err" | grep -qF "guozai.gzl"; then
    pass "stderr 显示 live 的实际账号 guozai.gzl"
else
    fail "stderr 缺 live 实际账号: $err"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ---------------------------------------------------------------------------
# Test 17: jarvis dir 已存在 → 不触发收编(即便 live 存在)
# ---------------------------------------------------------------------------
echo "=== Test 17: live 收编 — jarvis dir 已存在 → 不触发 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
cat > "$ROOT/auth.yaml" <<EOF
version: 1
current:
    user:
        account: open_jarvis
        marker: LIVE_UNUSED
EOF
seed_login "$ROOT" "jarvis"   # jarvis dir 里已经有 auth.yaml
JARVIS_PRE=$(cat "$ROOT/identities/jarvis/auth.yaml")
CAP=$(mktemp); ERR=$(mktemp)
STUB_WHOAMI_ACCOUNT="WORKER_1782379562571" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" status >/dev/null 2>"$ERR"
err=$(cat "$ERR")
if ! printf '%s' "$err" | grep -qF "首跑收编"; then
    pass "jarvis dir 已存在 → 不触发收编(stderr 无收编提示)"
else
    fail "不该触发收编却触发了: $err"
fi
# 且预置内容未被 live 覆盖
if [ "$(cat "$ROOT/identities/jarvis/auth.yaml")" = "$JARVIS_PRE" ]; then
    pass "已存在的 jarvis/auth.yaml 未被 live 覆盖"
else
    fail "预置 jarvis/auth.yaml 被覆盖了"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ---------------------------------------------------------------------------
# Test 18: v1 jarvis.auth.yaml 文件存在 → 不触发收编(与 dir 同权)
# ---------------------------------------------------------------------------
echo "=== Test 18: live 收编 — v1 jarvis.auth.yaml 存在 → 不触发 ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
cat > "$ROOT/auth.yaml" <<EOF
version: 1
current:
    user:
        account: open_jarvis
EOF
mkdir -p "$ROOT/identities"
cat > "$ROOT/identities/jarvis.auth.yaml" <<EOF
version: 1
current:
    user:
        account: v1_jarvis
EOF
CAP=$(mktemp); ERR=$(mktemp)
STUB_WHOAMI_ACCOUNT="WORKER_1782379562571" A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" \
    STUB_CAPTURE="$CAP" bash "$A1ID" status >/dev/null 2>"$ERR"
err=$(cat "$ERR")
# 这里预期:v1→v2 迁移会把 jarvis.auth.yaml 复制到 jarvis/auth.yaml,但收编门要跳过。
# 该门条件为「jarvis dir + v1 file 都不存在」——v1 存在故不触发收编;stderr 应看到迁移
# 提示而非"首跑收编"提示。
if ! printf '%s' "$err" | grep -qF "首跑收编"; then
    pass "v1 jarvis.auth.yaml 存在 → 收编门跳过(未触发 v2 首跑收编)"
else
    fail "收编门错误触发: $err"
fi
if printf '%s' "$err" | grep -qF "v1→v2 迁移"; then
    pass "但 v1→v2 迁移仍照常触发(v1 复制到 v2 布局)"
else
    fail "缺 v1→v2 迁移提示: $err"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR"

# ===========================================================================
# Test 19: post-PR 子代理对 Aone 严格只读，且 env -u 不能绕过进程上下文
# ===========================================================================
echo "=== Test 19: post-PR Aone read-only guard ==="
ROOT=$(new_root); BIN=$(mktemp -d); make_stub "$BIN"
seed_login "$ROOT" "terraform-rd"
CAP=$(mktemp); ERR=$(mktemp)

assert_post_pr_blocked() {
    local label="$1"; shift
    : > "$CAP"; : > "$ERR"
    JARVIS_A1_IDENTITY=terraform-rd JARVIS_A1_STRICT=1 \
        JARVIS_AONE_WRITE_POLICY=post-pr-read-only \
        A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
        bash "$A1ID" -- "$@" >/dev/null 2>"$ERR"
    local rc=$?
    if [ "$rc" != "0" ] && [ ! -s "$CAP" ]; then
        pass "$label 被 a1id fail-closed"
    else
        fail "$label 未阻断: rc=$rc capture=$(cat "$CAP")"
    fi
}

assert_post_pr_allowed() {
    local label="$1"; shift
    : > "$CAP"; : > "$ERR"
    JARVIS_A1_IDENTITY=terraform-rd JARVIS_A1_STRICT=1 \
        JARVIS_AONE_WRITE_POLICY=post-pr-read-only \
        A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
        bash "$A1ID" -- "$@" >/dev/null 2>"$ERR"
    local rc=$?
    if [ "$rc" = "0" ] && grep -qF "ARGS=$*" "$CAP"; then
        pass "$label 明确读操作放行"
    else
        fail "$label 被误拦: rc=$rc err=$(cat "$ERR")"
    fi
}

assert_post_pr_as_blocked() {
    local label="$1"; shift
    : > "$CAP"; : > "$ERR"
    JARVIS_AONE_WRITE_POLICY=post-pr-read-only \
        A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
        bash "$A1ID" as terraform-rd -- "$@" >/dev/null 2>"$ERR"
    local rc=$?
    if [ "$rc" != "0" ] && [ ! -s "$CAP" ]; then
        pass "$label 被显式身份入口 fail-closed"
    else
        fail "$label 未阻断: rc=$rc capture=$(cat "$CAP")"
    fi
}

assert_post_pr_blocked "update status" project workitem update 123 --status Closed
assert_post_pr_blocked "update tag" project workitem update 123 --tag post-pr-test-tag
assert_post_pr_blocked "create" project workitem create --project 528766 --title x
assert_post_pr_blocked "delete" project workitem delete 123
assert_post_pr_blocked "comment create" project workitem comment create 123 -m no
assert_post_pr_blocked "伪造 jarvis-claim marker" project workitem comment create 123 \
    -m "jarvis-claim host-a 2026-07-16T12:34:56Z deadbeef"
assert_post_pr_blocked "relation add" project workitem relation add 123 relate:456
assert_post_pr_blocked "relation remove" project workitem relation remove 123 relate:456
assert_post_pr_blocked "attachment upload" project workitem attachment upload 123 file.txt
assert_post_pr_blocked "attachment delete" project workitem attachment delete 123 9
assert_post_pr_blocked "global format prefix create" \
    --format json project workitem create --project 528766 --title x
assert_post_pr_blocked "interleaved global format relation" \
    project --format json workitem relation add 123 relate:456
assert_post_pr_as_blocked "显式 terraform-rd update" \
    project workitem update 123 --status Closed

assert_post_pr_allowed "get" project workitem get 123
assert_post_pr_allowed "list" project workitem list --project 528766
assert_post_pr_allowed "activity" project workitem activity 123
assert_post_pr_allowed "comment list" project workitem comment list 123
assert_post_pr_allowed "relation list" project workitem relation list 123
assert_post_pr_allowed "attachment list" project workitem attachment list 123
assert_post_pr_allowed "attachment download" project workitem attachment download 123 9
assert_post_pr_allowed "field list" project workitem field list --project 528766
assert_post_pr_allowed "global format prefix get" \
    --format json project workitem get 123

# A shell-level PreToolUse parser cannot prove arbitrary Python subprocess
# behavior.  The a1id runtime fence remains authoritative after that wrapper.
: > "$CAP"; : > "$ERR"
JARVIS_A1_IDENTITY=terraform-rd JARVIS_A1_STRICT=1 \
    JARVIS_AONE_WRITE_POLICY=post-pr-read-only \
    A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    python3 - "$A1ID" > /dev/null 2>"$ERR" <<'PY'
import os
import subprocess
import sys

result = subprocess.run([
    "bash", sys.argv[1], "--", "--format", "json",
    "project", "workitem", "create", "--project", "528766", "--title", "x",
], env=os.environ.copy(), check=False)
raise SystemExit(result.returncode)
PY
rc=$?
if [ "$rc" != "0" ] && [ ! -s "$CAP" ]; then
    pass "python subprocess 包装仍被 a1id runtime fence 阻断"
else
    fail "python subprocess 绕过 runtime fence: rc=$rc capture=$(cat "$CAP")"
fi
assert_post_pr_allowed "field options" project workitem field options status --project 528766

context_dir="/tmp/jarvis-headless-context-$(id -u)"
context_pgid="$(ps -o pgid= -p $$ | tr -d '[:space:]')"
context_marker="$context_dir/$context_pgid.json"
lineage_dir="$(mktemp -d)"
manager="$proj_root/bootstrap/jarvis-interactive-worker.py"
mkdir -p "$context_dir"
rm -f "$context_marker"
JARVIS_INTERACTIVE_STATE_DIR="$lineage_dir" \
JARVIS_CONTROL_PLANE_BASE_URL="http://127.0.0.1:1" \
JARVIS_HEADLESS_REMOTE_REGISTER_TIMEOUT="0.05" \
    /usr/bin/python3 -I "$manager" register-headless \
    --session-id "a1id-lineage-$$" --pid "$$" --client claude \
    --policy-revision terraform-rd-single-writer-v6 \
    --aone-write-policy post-pr-read-only \
    --headless-kind pr_comment_reply --aone-id 123 --project-id 528766 \
    --claim-attempt-id a1id-lineage-attempt \
    >/dev/null
: > "$CAP"; : > "$ERR"
env -u JARVIS_AONE_WRITE_POLICY \
    JARVIS_A1_IDENTITY=terraform-rd JARVIS_A1_STRICT=1 \
    JARVIS_INTERACTIVE_STATE_DIR="$lineage_dir" \
    A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    bash "$A1ID" -- project workitem update 123 --status Closed \
    >/dev/null 2>"$ERR"
rc=$?
if [ "$rc" != "0" ] && [ ! -s "$CAP" ]; then
    pass "marker 缺失 + env -u 后仍由 canonical worker lineage 阻断 Aone 写"
else
    fail "marker 缺失时 lineage 未阻断写: rc=$rc capture=$(cat "$CAP")"
fi

printf '%s' '{broken-json' > "$context_marker"
: > "$CAP"; : > "$ERR"
env -u JARVIS_AONE_WRITE_POLICY \
    JARVIS_A1_IDENTITY=terraform-rd JARVIS_A1_STRICT=1 \
    JARVIS_INTERACTIVE_STATE_DIR="$lineage_dir" \
    A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    bash "$A1ID" -- project workitem comment create 123 -m no \
    >/dev/null 2>"$ERR"
rc=$?
if [ "$rc" != "0" ] && [ ! -s "$CAP" ]; then
    pass "marker 损坏 + env 清空后 post-PR Aone 写仍被拒绝"
else
    fail "损坏 marker 绕过 lineage: rc=$rc capture=$(cat "$CAP")"
fi

: > "$CAP"; : > "$ERR"
env -u JARVIS_AONE_WRITE_POLICY \
    JARVIS_A1_IDENTITY=terraform-rd JARVIS_A1_STRICT=1 \
    JARVIS_INTERACTIVE_STATE_DIR="$lineage_dir" \
    A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    bash "$A1ID" -- project workitem get 123 \
    >/dev/null 2>"$ERR"
rc=$?
if [ "$rc" = "0" ] && grep -qF "ARGS=project workitem get 123" "$CAP"; then
    pass "marker 损坏 + env 清空后明确 Aone 读仍允许"
else
    fail "post-PR lineage 误拦只读: rc=$rc err=$(cat "$ERR")"
fi

rm -f "$context_marker" "$lineage_dir"/*.json
sleep 30 &
ended_pid=$!
JARVIS_INTERACTIVE_STATE_DIR="$lineage_dir" \
JARVIS_CONTROL_PLANE_BASE_URL="http://127.0.0.1:1" \
JARVIS_HEADLESS_REMOTE_REGISTER_TIMEOUT="0.05" \
    /usr/bin/python3 -I "$manager" register-headless \
    --session-id "a1id-ended-$$" --pid "$ended_pid" --client claude \
    --policy-revision terraform-rd-single-writer-v6 \
    --aone-write-policy post-pr-read-only \
    --headless-kind pr_ci_fix --aone-id 123 --project-id 528766 \
    --claim-attempt-id a1id-ended-attempt \
    >/dev/null
kill "$ended_pid" 2>/dev/null || true
wait "$ended_pid" 2>/dev/null || true
: > "$CAP"; : > "$ERR"
env -u JARVIS_AONE_WRITE_POLICY \
    JARVIS_A1_IDENTITY=terraform-rd JARVIS_A1_STRICT=1 \
    JARVIS_INTERACTIVE_STATE_DIR="$lineage_dir" \
    A1ID_ROOT="$ROOT" A1_BIN="$BIN/a1" STUB_CAPTURE="$CAP" \
    bash "$A1ID" -- project workitem update 123 --status Closed \
    >/dev/null 2>"$ERR"
rc=$?
if [ "$rc" = "0" ] && grep -qF \
        "ARGS=project workitem update 123 --status Closed" "$CAP"; then
    pass "已结束 worker 的 tombstone 不误伤正常交互写"
else
    fail "已结束 worker 仍误拦当前进程: rc=$rc err=$(cat "$ERR")"
fi
rm -rf "$ROOT" "$BIN" "$CAP" "$ERR" "$lineage_dir"

# ===========================================================================
# Test 20: 两份活跃 aone-triage Skill 的个人身份纪律枚举保持一致
# ===========================================================================
echo "=== Test 20: aone-triage Skill 个人身份枚举包含 shanye ==="
for skill in \
    "$proj_root/.agents/skills/aone-triage/SKILL.md" \
    "$proj_root/.claude/skills/aone-triage/SKILL.md"; do
    rel="${skill#"$proj_root/"}"
    count=$(grep -cF "chenyi/guozai/linjun/shanye" "$skill" || true)
    if [ "$count" = "2" ]; then
        pass "$rel 两处个人身份枚举一致"
    else
        fail "$rel 应有 2 处完整枚举,实际 $count 处"
    fi
done

# ---------------------------------------------------------------------------
echo ""
echo "=== SUMMARY ==="
echo "PASS: $PASS  FAIL: $FAIL"
if [ "$FAIL" -gt 0 ]; then
    echo "TESTS FAILED"
    exit 1
fi
echo "All tests passed"
exit 0
