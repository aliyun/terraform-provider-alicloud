#!/usr/bin/env bash
# test/rc_gate_test.sh — hermetic tests for bootstrap/rc-gate.sh (F3 RC 门禁线)
#
# 不依赖 terraform / 网络 / 凭证 / 真实 playground / 真实 provider 仓。
#   - probe.sh 全程 PATH 桩(可控 stdout/verdict/退码),rc-gate 经 `command -v probe.sh` 命中桩。
#   - 断言:红/黄/绿判定、退码(红1/黄0/绿0/不可判2)、报告生成、quick 模式分支(run --dry)、
#     tier0 参数(--all --limit 200)、场景 fail/destroy 残留=红、env 阻断=黄、总超时=黄跳过。
#
# Run: bash test/rc_gate_test.sh  → 打印 PASS/FAIL 汇总,全过退 0,任一失败退 1。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJ_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
RC_GATE="$PROJ_ROOT/bootstrap/rc-gate.sh"

pass=0; fail=0
ok()  { echo "  PASS: $1"; pass=$((pass+1)); }
bad() { echo "  FAIL: $1${2:+ — $2}"; fail=$((fail+1)); }

# 清理清单
TMPS=()
cleanup() { for d in "${TMPS[@]:-}"; do [ -n "$d" ] && rm -rf "$d" 2>/dev/null; done; }
trap cleanup INT TERM EXIT

DAY="$(date -u +%Y%m%d)"

# mk_env — 建一次性桩环境:PATH 桩 probe.sh + provider-dir + audit-dir + calls.log。
# 打印四行:STUBDIR PROVDIR AUDITDIR CALLS(供调用方 read)。
mk_env() {
    local base; base="$(mktemp -d)"; TMPS+=("$base")
    local stubdir="$base/bin" provdir="$base/prov" auditdir="$base/audit" work="$base/work"
    mkdir -p "$stubdir" "$provdir" "$auditdir" "$work"
    # provider-dir 只需存在(桩 tier0 不校验其内部布局)
    cat > "$stubdir/probe.sh" <<'STUB'
#!/usr/bin/env bash
# 桩 probe.sh:数据驱动(全部经 env 控制),记录每次调用到 STUB_CALLS_LOG。
echo "probe.sh $*" >> "${STUB_CALLS_LOG:-/dev/null}"
cmd="${1:-}"; shift 2>/dev/null || true
case "$cmd" in
  doctor) exit 0 ;;
  tier0)
    if [ "${STUB_TIER0_RC:-0}" = "2" ]; then
      echo "tier0: 本地 provider 仓不可用(桩)" >&2; exit 2
    fi
    findings="${STUB_TIER0_FINDINGS:-[]}"
    qn="${STUB_TIER0_QUEUE:-0}"
    mech="${STUB_TIER0_MECH:-on}"
    vj="${STUB_WORK}/tier0-verdict.json"
    jq -n --arg mech "$mech" --argjson f "$findings" --argjson qn "$qn" \
      '{schema_version:2, mode:"tier0", mech:$mech, findings:$f,
        judgment_queue:[range(0;$qn)|{resource:("r"+(.|tostring)), reason:"meta_unavailable"}],
        stats:{findings:($f|length), api_findings:([$f[]|select(.code|startswith("api_gap"))]|length)}}' > "$vj"
    api_nf="$(echo "$findings" | jq '[.[]|select(.code|startswith("api_gap"))]|length')"
    echo "verdict: $vj"
    echo "mech: $mech / findings: $(echo "$findings"|jq length) (api_gap_* $api_nf) / resources: 1 / judgment_queue: $qn / suppressed: 0"
    exit "${STUB_TIER0_RC:-0}"
    ;;
  list)
    printf '%-10s %-18s %-9s %-40s %s\n' PRODUCT ID PERSONA RESOURCES DETECT
    for s in ${STUB_SCENARIOS:-}; do
      printf '%-10s %-18s %-9s %-40s %s\n' prod "$s" persona alicloud_x detect
    done
    exit 0
    ;;
  run)
    sid=""; dry=0
    for a in "$@"; do
      case "$a" in --dry) dry=1 ;; -*) : ;; *) [ -z "$sid" ] && sid="$a" ;; esac
    done
    # --dry(quick):probe 真实语义恒退 0(plan 预览,不 apply)
    if [ "$dry" -eq 1 ]; then echo "probe plan: scenario=$sid"; exit 0; fi
    rc=0
    for pair in ${STUB_RUN_RCS:-}; do
      k="${pair%%:*}"; v="${pair##*:}"; [ "$k" = "$sid" ] && rc="$v"
    done
    echo "verdict: ${STUB_WORK}/run-$sid.json"
    echo "audit:   ${STUB_WORK}/run-$sid.json"
    exit "$rc"
    ;;
  *) echo "probe(桩): 未知命令 '$cmd'" >&2; exit 2 ;;
esac
STUB
    chmod +x "$stubdir/probe.sh"
    printf '%s\n%s\n%s\n%s\n%s\n' "$stubdir" "$provdir" "$auditdir" "$base/calls.log" "$work"
}

# run_gate — 在桩环境跑 rc-gate;调用前 export 好 STUB_* 控制变量。
#   参数:全部原样传给 rc-gate。使用外部变量 STUBDIR/PROVDIR/AUDITDIR/CALLS/WORK。
#   捕获 OUT(stdout+stderr)与 RC。
run_gate() {
    OUT="$(PATH="$STUBDIR:$PATH" \
        STUB_CALLS_LOG="$CALLS" STUB_WORK="$WORK" \
        RC_GATE_AUDIT_DIR="$AUDITDIR" \
        bash "$RC_GATE" "$@" 2>&1)"
    RC=$?
}

# 每个用例开头调用:重置环境 + 清 STUB_* 控制变量。
fresh() {
    { read -r STUBDIR; read -r PROVDIR; read -r AUDITDIR; read -r CALLS; read -r WORK; } < <(mk_env)
    unset STUB_TIER0_MECH STUB_TIER0_FINDINGS STUB_TIER0_QUEUE STUB_TIER0_RC \
          STUB_SCENARIOS STUB_RUN_RCS RC_GATE_QUEUE_YELLOW RC_GATE_TOTAL_TIMEOUT \
          RC_GATE_START_EPOCH 2>/dev/null || true
    export STUB_TIER0_MECH="on" STUB_TIER0_FINDINGS="[]" STUB_TIER0_QUEUE="0"
    REPORT="$AUDITDIR/${DAY}-report.md"
}

report_has() { grep -qE "$1" "$REPORT" 2>/dev/null; }

# ===========================================================================
echo "Test 1: 缺 provider-dir → 退 2 + usage"
fresh
run_gate
[ "$RC" -eq 2 ] && ok "无参退 2" || bad "无参应退 2" "got $RC"
echo "$OUT" | grep -qiE "用法|usage|provider-dir" && ok "打印 usage" || bad "应打印 usage" "$OUT"

# ===========================================================================
echo "Test 2: tier0 api_gap S3 → 红 + 退 1 + 报告 VERDICT: RED + tier0 参数 --all --limit 200"
fresh
export STUB_TIER0_FINDINGS='[{"code":"api_gap_required","resource":"alicloud_vpc","attribute":"name","severity_hint":"S3","summary":"TF Optional 但 API 必填"}]'
export STUB_TIER0_RC=1
export STUB_SCENARIOS=""
run_gate "$PROVDIR"
[ "$RC" -eq 1 ] && ok "api_gap S3 → 退 1" || bad "应退 1" "got $RC"
echo "$OUT" | grep -q "VERDICT: RED" && ok "stdout VERDICT: RED" || bad "stdout 应 RED" "$OUT"
report_has "VERDICT: RED" && ok "报告 VERDICT: RED" || bad "报告应 RED"
report_has "api_gap" && ok "报告含 api_gap 明细" || bad "报告应含 api_gap"
grep -q "tier0 --all --limit 200" "$CALLS" && ok "tier0 以 --all --limit 200 调用" || bad "tier0 参数应 --all --limit 200" "$(cat "$CALLS")"

# ===========================================================================
echo "Test 3: tier0 仅 api_gap S4(低于 S3+)→ 黄 + 退 0(非阻断 finding)"
fresh
export STUB_TIER0_FINDINGS='[{"code":"api_gap_range","resource":"alicloud_vpc","attribute":"cidr","severity_hint":"S4","summary":"范围越界"}]'
export STUB_TIER0_RC=1
export STUB_SCENARIOS="s-green"; export STUB_RUN_RCS="s-green:0"
run_gate "$PROVDIR"
[ "$RC" -eq 0 ] && ok "S4-only → 退 0" || bad "应退 0" "got $RC"
echo "$OUT" | grep -q "VERDICT: YELLOW" && ok "VERDICT: YELLOW" || bad "应 YELLOW" "$OUT"

# ===========================================================================
echo "Test 4: tier0 机械层 degraded → 黄 + 退 0 + 报告标注降级"
fresh
export STUB_TIER0_MECH="degraded"; export STUB_TIER0_RC=0
export STUB_SCENARIOS="s-green"; export STUB_RUN_RCS="s-green:0"
run_gate "$PROVDIR"
[ "$RC" -eq 0 ] && ok "degraded → 退 0" || bad "应退 0" "got $RC"
echo "$OUT" | grep -q "VERDICT: YELLOW" && ok "VERDICT: YELLOW" || bad "应 YELLOW" "$OUT"
report_has "degraded|降级" && ok "报告标注机械层降级" || bad "报告应标注降级"

# ===========================================================================
echo "Test 5: tier0 judgment_queue 激增(> 阈值)→ 黄 + 退 0"
fresh
export STUB_TIER0_MECH="on"; export STUB_TIER0_QUEUE=99; export STUB_TIER0_RC=0
export RC_GATE_QUEUE_YELLOW=10
export STUB_SCENARIOS="s-green"; export STUB_RUN_RCS="s-green:0"
run_gate "$PROVDIR"
[ "$RC" -eq 0 ] && ok "queue 激增 → 退 0" || bad "应退 0" "got $RC"
echo "$OUT" | grep -q "VERDICT: YELLOW" && ok "VERDICT: YELLOW" || bad "应 YELLOW" "$OUT"
report_has "queue|队列" && ok "报告标注 queue" || bad "报告应标注 queue"

# ===========================================================================
echo "Test 6: 全绿(tier0 净 + 场景全过, full 模式)→ 绿 + 退 0"
fresh
export STUB_TIER0_MECH="on"; export STUB_TIER0_QUEUE=0; export STUB_TIER0_RC=0
export STUB_SCENARIOS="s1 s2"; export STUB_RUN_RCS="s1:0 s2:0"
run_gate "$PROVDIR"
[ "$RC" -eq 0 ] && ok "全绿 → 退 0" || bad "应退 0" "got $RC"
echo "$OUT" | grep -q "VERDICT: GREEN" && ok "VERDICT: GREEN" || bad "应 GREEN" "$OUT"
# full 模式:run 不带 --dry
grep -qE "run s1( |$)" "$CALLS" && ! grep -q "run s1 --dry" "$CALLS" && ok "full 模式真实 run(无 --dry)" || bad "full 模式应无 --dry" "$(cat "$CALLS")"

# ===========================================================================
echo "Test 7: 场景 run 退 1(provider finding)→ 红 + 退 1"
fresh
export STUB_TIER0_RC=0
export STUB_SCENARIOS="s1 s2"; export STUB_RUN_RCS="s1:0 s2:1"
run_gate "$PROVDIR"
[ "$RC" -eq 1 ] && ok "场景 fail → 退 1" || bad "应退 1" "got $RC"
echo "$OUT" | grep -q "VERDICT: RED" && ok "VERDICT: RED" || bad "应 RED" "$OUT"
report_has "s2" && ok "报告点名失败场景 s2" || bad "报告应点名 s2"

# ===========================================================================
echo "Test 8: 场景 run 退 3(destroy 残留)→ 红 + 退 1 + 报告标 destroy/清理"
fresh
export STUB_TIER0_RC=0
export STUB_SCENARIOS="s1"; export STUB_RUN_RCS="s1:3"
run_gate "$PROVDIR"
[ "$RC" -eq 1 ] && ok "destroy 残留 → 退 1" || bad "应退 1" "got $RC"
echo "$OUT" | grep -q "VERDICT: RED" && ok "VERDICT: RED" || bad "应 RED" "$OUT"
report_has "destroy|残留|清理" && ok "报告标注 destroy 残留" || bad "报告应标注残留"

# ===========================================================================
echo "Test 9: 场景 run 退 2(env 阻断)→ 黄 + 退 0(非 provider bug)"
fresh
export STUB_TIER0_RC=0
export STUB_SCENARIOS="s1"; export STUB_RUN_RCS="s1:2"
run_gate "$PROVDIR"
[ "$RC" -eq 0 ] && ok "env 阻断 → 退 0" || bad "应退 0" "got $RC"
echo "$OUT" | grep -q "VERDICT: YELLOW" && ok "VERDICT: YELLOW" || bad "应 YELLOW" "$OUT"

# ===========================================================================
echo "Test 10: --quick → run --dry + 黄 + 退 0 + 报告标 quick 未跑 apply"
fresh
export STUB_TIER0_RC=0
export STUB_SCENARIOS="s1 s2"; export STUB_RUN_RCS="s1:0 s2:0"
run_gate "$PROVDIR" --quick
[ "$RC" -eq 0 ] && ok "quick → 退 0" || bad "应退 0" "got $RC"
echo "$OUT" | grep -q "VERDICT: YELLOW" && ok "VERDICT: YELLOW(quick 未跑 apply)" || bad "应 YELLOW" "$OUT"
grep -q "run s1 --dry" "$CALLS" && ok "quick 模式用 run --dry" || bad "quick 应 run --dry" "$(cat "$CALLS")"
report_has "quick|未跑 apply|plan 为止" && ok "报告标注 quick 语义" || bad "报告应标注 quick"

# ===========================================================================
echo "Test 11: --quick 但 tier0 红 → 红压过 quick 黄 + 退 1"
fresh
export STUB_TIER0_FINDINGS='[{"code":"api_gap_type","resource":"alicloud_vpc","attribute":"x","severity_hint":"S3","summary":"类型冲突"}]'
export STUB_TIER0_RC=1
export STUB_SCENARIOS="s1"; export STUB_RUN_RCS="s1:0"
run_gate "$PROVDIR" --quick
[ "$RC" -eq 1 ] && ok "quick+tier0 红 → 退 1" || bad "应退 1" "got $RC"
echo "$OUT" | grep -q "VERDICT: RED" && ok "VERDICT: RED" || bad "应 RED" "$OUT"

# ===========================================================================
echo "Test 12: tier0 runner 错误(退 2)→ 不可判 + 退 2"
fresh
export STUB_TIER0_RC=2
export STUB_SCENARIOS="s1"; export STUB_RUN_RCS="s1:0"
run_gate "$PROVDIR"
[ "$RC" -eq 2 ] && ok "tier0 退 2 → 门禁退 2" || bad "应退 2" "got $RC"
echo "$OUT" | grep -qE "VERDICT: CANNOT_CERTIFY|不可判" && ok "VERDICT: CANNOT_CERTIFY" || bad "应 CANNOT_CERTIFY" "$OUT"

# ===========================================================================
echo "Test 13: 无场景(tier0 净, list 空)→ 黄 + 退 0(tier-1 零覆盖)"
fresh
export STUB_TIER0_RC=0; export STUB_SCENARIOS=""
run_gate "$PROVDIR"
[ "$RC" -eq 0 ] && ok "无场景 → 退 0" || bad "应退 0" "got $RC"
echo "$OUT" | grep -q "VERDICT: YELLOW" && ok "VERDICT: YELLOW(零覆盖)" || bad "应 YELLOW" "$OUT"

# ===========================================================================
echo "Test 14: 报告落盘 <date>-report.md + 含 tier-0 与 tier-1 两节"
fresh
export STUB_TIER0_RC=0; export STUB_SCENARIOS="s1"; export STUB_RUN_RCS="s1:0"
run_gate "$PROVDIR"
[ -f "$REPORT" ] && ok "报告生成于 $AUDITDIR/${DAY}-report.md" || bad "报告未生成" "$REPORT"
report_has "tier-0|tier0|三方一致性" && ok "报告含 tier-0 节" || bad "报告缺 tier-0 节"
report_has "tier-1|tier1|场景" && ok "报告含 tier-1 节" || bad "报告缺 tier-1 节"

# ===========================================================================
echo "Test 15: 总超时预算耗尽 → 剩余场景跳过(黄) + 未真正 run"
fresh
export STUB_TIER0_RC=0; export STUB_SCENARIOS="s1 s2"; export STUB_RUN_RCS="s1:0 s2:0"
export RC_GATE_TOTAL_TIMEOUT=1
# 注入极早的起点,使 tier-1 阶段一进入即判超时(测试专用 override)
export RC_GATE_START_EPOCH=1
OUT="$(PATH="$STUBDIR:$PATH" STUB_CALLS_LOG="$CALLS" STUB_WORK="$WORK" \
    RC_GATE_AUDIT_DIR="$AUDITDIR" RC_GATE_TOTAL_TIMEOUT=1 RC_GATE_START_EPOCH=1 \
    bash "$RC_GATE" "$PROVDIR" 2>&1)"; RC=$?
[ "$RC" -eq 0 ] && ok "超时跳过 → 退 0(黄)" || bad "应退 0" "got $RC"
echo "$OUT" | grep -q "VERDICT: YELLOW" && ok "VERDICT: YELLOW(超时)" || bad "应 YELLOW" "$OUT"
report_has "超时|timeout|预算" && ok "报告标注超时" || bad "报告应标注超时"
! grep -qE "^probe.sh run " "$CALLS" && ok "超时预算耗尽:未真正 run 场景" || bad "不应发起 run" "$(cat "$CALLS")"

# ===========================================================================
echo ""
echo "==================================================="
echo "rc_gate_test: PASS=$pass FAIL=$fail"
echo "==================================================="
[ "$fail" -eq 0 ] || exit 1
