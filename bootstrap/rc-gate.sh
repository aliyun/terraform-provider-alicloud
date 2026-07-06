#!/usr/bin/env bash
# bootstrap/rc-gate.sh — F3 RC 门禁线:发版前全量探测过闸
#
# 用法:
#   bootstrap/rc-gate.sh <provider-dir> [--quick]
#
#   <provider-dir>  本地 terraform-provider-alicloud 仓(传给 probe.sh tier0 的
#                   JARVIS_PROBE_PROVIDER_DIR;需含 website/docs/r + alicloud)。
#   --quick         快扫模式:tier-1 场景改为 plan 为止(probe.sh run --dry,零创建),
#                   不真实 apply。快但覆盖弱 → 整体至多判黄(未跑 apply)。
#
# 三步(全部调用现有 probe.sh 子命令,不复制其逻辑):
#   ① tier-0 机械层全量:probe.sh tier0 --all --limit 200
#        降级容忍——mech=degraded 记黄不记红(静态扫描仍出,只是 OpenAPI 侧转人工队列)。
#   ② tier-1 全场景生命周期:probe.sh list 枚举 → 逐场景 probe.sh run(顺跑,并发 1)。
#        --quick → run --dry(plan 为止);完整模式 → run(真实 apply);总超时预算可配。
#   ③ 汇总 runs/rc-gate/<date>-report.md 并按红/黄/绿判定退码。
#
# 判定与退码:
#   🔴 RED  (退 1,阻断):tier-0 api_gap 严重度 S3+ / 场景 run 退 1(provider finding)
#                        / 场景 run 退 3(destroy 失败或 state 残留)。
#   🟡 YELLOW(退 0,放行但显著标注):tier-0 机械层降级(mech=degraded)/ judgment_queue 激增
#                        / tier-0 有非 S3+ finding / 场景 run 退 2(env 阻断,非 provider bug)
#                        / quick 模式(未跑 apply)/ 无场景可跑 / 超时预算耗尽跳过。
#   🟢 GREEN(退 0):以上皆无。
#   ⚪ CANNOT_CERTIFY(退 2):tier-0 无法运行(probe.sh tier0 退 2),门禁不完整,修环境后重跑。
#   优先级:RED > CANNOT_CERTIFY > YELLOW > GREEN。
#
# 环境变量(多为测试/调优用):
#   RC_GATE_PROBE          — probe.sh 路径覆盖(默认 PATH 上的 probe.sh,兜底同目录 sibling)
#   RC_GATE_AUDIT_DIR      — 报告落盘目录(默认 <jarvis 根>/runs/rc-gate)
#   RC_GATE_QUEUE_YELLOW   — judgment_queue 激增黄线阈值(默认 40)
#   RC_GATE_TOTAL_TIMEOUT  — tier-1 阶段总超时预算秒(默认 0 = 不限);超预算后剩余场景跳过记黄
#   RC_GATE_START_EPOCH    — 测试专用:覆盖 tier-1 阶段起点 epoch(用于确定性触发超时分支)
#   JARVIS_PROBE_PROVIDER_DIR — 若已设,<provider-dir> 位置参数仍以命令行为准
#
# 退出码:0=绿/黄(放行);1=红(阻断);2=不可判 / 用法错误。

set -uo pipefail

_rcg_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/lib.sh
source "$_rcg_dir/lib.sh"

QUEUE_YELLOW="${RC_GATE_QUEUE_YELLOW:-40}"
TOTAL_TIMEOUT="${RC_GATE_TOTAL_TIMEOUT:-0}"

_usage() {
    sed -n '2,20p' "$0" >&2
}

# probe.sh 解析:显式 override > PATH(便于 PATH 桩测试)> 同目录 sibling(生产兜底)。
_resolve_probe() {
    if [ -n "${RC_GATE_PROBE:-}" ]; then echo "$RC_GATE_PROBE"; return 0; fi
    if command -v probe.sh >/dev/null 2>&1; then command -v probe.sh; return 0; fi
    echo "$_rcg_dir/probe.sh"
}

# ── 参数解析 ─────────────────────────────────────────────────────────
PROVIDER_DIR=""
QUICK=0
while [ $# -gt 0 ]; do
    case "$1" in
        --quick) QUICK=1; shift ;;
        -h|--help) _usage; exit 2 ;;
        -*) echo "rc-gate: 未知参数 '$1'" >&2; _usage; exit 2 ;;
        *) if [ -z "$PROVIDER_DIR" ]; then PROVIDER_DIR="$1"; shift
           else echo "rc-gate: 多余参数 '$1'" >&2; _usage; exit 2; fi ;;
    esac
done

if [ -z "$PROVIDER_DIR" ]; then
    echo "rc-gate: 缺少 <provider-dir>" >&2
    _usage
    exit 2
fi
if [ ! -d "$PROVIDER_DIR" ]; then
    echo "rc-gate: <provider-dir> 不是目录: $PROVIDER_DIR" >&2
    exit 2
fi
PROVIDER_DIR="$(cd "$PROVIDER_DIR" && pwd)"

PROBE="$(_resolve_probe)"
MODE="full"; [ "$QUICK" -eq 1 ] && MODE="quick"
DAY="$(date -u +%Y%m%d)"
STARTED_AT="$(date -u +%FT%TZ)"
GATE_START_EPOCH="$(date +%s)"
AUDIT_DIR="${RC_GATE_AUDIT_DIR:-$(jarvis_root)/runs/rc-gate}"
REPORT="$AUDIT_DIR/${DAY}-report.md"
mkdir -p "$AUDIT_DIR"

# tier0 用 JARVIS_PROBE_PROVIDER_DIR 定位 provider 仓
export JARVIS_PROBE_PROVIDER_DIR="$PROVIDER_DIR"

# ── 判定累加器 ───────────────────────────────────────────────────────
RED_REASONS=()
YELLOW_REASONS=()

add_red()    { RED_REASONS+=("$1"); }
add_yellow() { YELLOW_REASONS+=("$1"); }

# ── ① tier-0 全量机械层 ─────────────────────────────────────────────
echo "rc-gate: [1/3] tier-0 三方一致性全量扫描(probe.sh tier0 --all --limit 200)…" >&2
T0_OUT="$("$PROBE" tier0 --all --limit 200 2>&1)"; T0_RC=$?
T0_VERDICT="$(printf '%s\n' "$T0_OUT" | sed -n 's/^verdict: //p' | head -1)"

T0_MECH="unknown"; T0_FINDINGS=0; T0_API_S3PLUS=0; T0_QUEUE=0
T0_API_S3PLUS_LINES=""
T0_FINDING_LINES=""
if [ "$T0_RC" -eq 2 ]; then
    add_red "__CANNOT_CERTIFY__:tier-0 无法运行(probe.sh tier0 退 2:$(printf '%s' "$T0_OUT" | tail -1))"
elif [ -n "$T0_VERDICT" ] && [ -f "$T0_VERDICT" ] && jq -e . "$T0_VERDICT" >/dev/null 2>&1; then
    T0_MECH="$(jq -r '.mech // "unknown"' "$T0_VERDICT")"
    T0_FINDINGS="$(jq '(.findings // [])|length' "$T0_VERDICT")"
    T0_QUEUE="$(jq '(.judgment_queue // [])|length' "$T0_VERDICT")"
    T0_API_S3PLUS="$(jq '[(.findings // [])[] | select((.code//"")|startswith("api_gap")) | select(.severity_hint=="S1" or .severity_hint=="S2" or .severity_hint=="S3")] | length' "$T0_VERDICT")"
    T0_API_S3PLUS_LINES="$(jq -r '(.findings // [])[] | select((.code//"")|startswith("api_gap")) | select(.severity_hint=="S1" or .severity_hint=="S2" or .severity_hint=="S3") | "  - [\(.severity_hint)] `\(.code)` \(.resource // "?").\(.attribute // "?"): \(.summary // "")"' "$T0_VERDICT" 2>/dev/null)"
    T0_FINDING_LINES="$(jq -r '(.findings // [])[] | "  - [\(.severity_hint // "?")] `\(.code // "?")` \(.resource // "?").\(.attribute // "?"): \(.summary // "")"' "$T0_VERDICT" 2>/dev/null)"

    if [ "${T0_API_S3PLUS:-0}" -gt 0 ]; then
        add_red "tier-0 api_gap S3+ ×${T0_API_S3PLUS}(客户端/接入面与 OpenAPI 硬冲突)"
    fi
    if [ "$T0_MECH" = "degraded" ]; then
        add_yellow "tier-0 机械层降级(mech=degraded):OpenAPI 侧转 judgment_queue 人工双层查证"
    fi
    if [ "${T0_QUEUE:-0}" -gt "$QUEUE_YELLOW" ]; then
        add_yellow "tier-0 judgment_queue 激增(${T0_QUEUE} > 阈值 ${QUEUE_YELLOW}):人工待判面骤增"
    fi
    # 有 finding 但无 S3+ 红项 → 记黄(doc_gap / api_gap S4 等非阻断项)
    if [ "${T0_FINDINGS:-0}" -gt 0 ] && [ "${T0_API_S3PLUS:-0}" -eq 0 ]; then
        add_yellow "tier-0 非 S3+ finding ×${T0_FINDINGS}(doc_gap / 低severity api_gap,发版前建议清)"
    fi
else
    # rc 0/1 但 verdict 不可解析:保守记黄,不臆断严重度为红
    if [ "$T0_RC" -eq 1 ]; then
        add_yellow "tier-0 有 finding 但 verdict 不可解析(退码 1),严重度未知,人工复核"
    fi
fi

# ── ② tier-1 全场景生命周期 ─────────────────────────────────────────
echo "rc-gate: [2/3] tier-1 全场景枚举(probe.sh list)…" >&2
LIST_OUT="$("$PROBE" list 2>/dev/null)"
SCENARIOS=()
while IFS= read -r _sid; do
    [ -n "$_sid" ] && SCENARIOS+=("$_sid")
done < <(printf '%s\n' "$LIST_OUT" | awk 'NR>1 && $1!="PRODUCT" && $2!="" {print $2}')

SCEN_ROWS=""   # 报告表格行累加
TIER1_START_EPOCH="${RC_GATE_START_EPOCH:-$(date +%s)}"
n_pass=0; n_fail=0; n_destroy=0; n_env=0; n_planonly=0; n_timeout=0

if [ "${#SCENARIOS[@]}" -eq 0 ]; then
    add_yellow "tier-1 零覆盖:probe.sh list 未发现任何场景"
else
    echo "rc-gate: [2/3] tier-1 场景生命周期(mode=$MODE, 并发 1, 共 ${#SCENARIOS[@]} 场景)…" >&2
    for sid in "${SCENARIOS[@]}"; do
        # 总超时预算:超预算后剩余场景一律跳过记黄(不再发起 run)
        if [ "$TOTAL_TIMEOUT" -gt 0 ]; then
            _now="$(date +%s)"; _elapsed=$(( _now - TIER1_START_EPOCH ))
            if [ "$_elapsed" -ge "$TOTAL_TIMEOUT" ]; then
                n_timeout=$(( n_timeout + 1 ))
                SCEN_ROWS="${SCEN_ROWS}| \`${sid}\` | - | ⏭️ 超时跳过(预算 ${TOTAL_TIMEOUT}s 耗尽) |
"
                continue
            fi
        fi

        if [ "$QUICK" -eq 1 ]; then
            "$PROBE" run "$sid" --dry >/dev/null 2>&1; _rrc=$?
            n_planonly=$(( n_planonly + 1 ))
            SCEN_ROWS="${SCEN_ROWS}| \`${sid}\` | ${_rrc} | 🟡 plan 为止(quick,未 apply) |
"
            continue
        fi

        "$PROBE" run "$sid" >/dev/null 2>&1; _rrc=$?
        case "$_rrc" in
            0) n_pass=$(( n_pass + 1 ))
               SCEN_ROWS="${SCEN_ROWS}| \`${sid}\` | 0 | 🟢 全生命周期通过 |
" ;;
            1) n_fail=$(( n_fail + 1 ))
               add_red "tier-1 场景 \`${sid}\` 退 1(provider finding:validate/plan/apply/幂等/import 等)"
               SCEN_ROWS="${SCEN_ROWS}| \`${sid}\` | 1 | 🔴 场景 fail(provider finding) |
" ;;
            3) n_destroy=$(( n_destroy + 1 ))
               add_red "tier-1 场景 \`${sid}\` 退 3(destroy 失败 / state 残留,零残留纪律破坏)"
               SCEN_ROWS="${SCEN_ROWS}| \`${sid}\` | 3 | 🔴 destroy 残留/清理失败 |
" ;;
            2) n_env=$(( n_env + 1 ))
               add_yellow "tier-1 场景 \`${sid}\` 退 2(env 阻断:凭证/网络/prepaid/plan-only,非 provider bug)"
               SCEN_ROWS="${SCEN_ROWS}| \`${sid}\` | 2 | 🟡 env 阻断(跳过,非 provider bug) |
" ;;
            *) n_env=$(( n_env + 1 ))
               add_yellow "tier-1 场景 \`${sid}\` 退 ${_rrc}(异常退码,人工复核)"
               SCEN_ROWS="${SCEN_ROWS}| \`${sid}\` | ${_rrc} | 🟡 异常退码 |
" ;;
        esac
    done
fi

# 超时预算耗尽:跳过的场景 = 覆盖缺口,整体标黄
if [ "$n_timeout" -gt 0 ]; then
    add_yellow "tier-1 总超时预算(${TOTAL_TIMEOUT}s)耗尽,跳过 ${n_timeout} 场景(未跑,覆盖缺口)"
fi

# quick 模式:整体标黄(未跑 apply,覆盖弱)
if [ "$QUICK" -eq 1 ] && [ "${#SCENARIOS[@]}" -gt 0 ]; then
    add_yellow "quick 模式:tier-1 仅 plan 为止(未 apply),覆盖弱,发版前建议完整模式复跑"
fi

# ── 判定收敛 ─────────────────────────────────────────────────────────
# CANNOT_CERTIFY 以特殊 red reason 编码,拆出来单独判优先级。
CANNOT_CERTIFY=0
REAL_RED=()
if [ "${#RED_REASONS[@]}" -gt 0 ]; then
    for r in "${RED_REASONS[@]}"; do
        case "$r" in
            __CANNOT_CERTIFY__:*) CANNOT_CERTIFY=1 ;;
            *) REAL_RED+=("$r") ;;
        esac
    done
fi

VERDICT="GREEN"; EXIT_CODE=0
if [ "${#REAL_RED[@]}" -gt 0 ]; then
    VERDICT="RED"; EXIT_CODE=1
elif [ "$CANNOT_CERTIFY" -eq 1 ]; then
    VERDICT="CANNOT_CERTIFY"; EXIT_CODE=2
elif [ "${#YELLOW_REASONS[@]}" -gt 0 ]; then
    VERDICT="YELLOW"; EXIT_CODE=0
else
    VERDICT="GREEN"; EXIT_CODE=0
fi

DUR=$(( $(date +%s) - GATE_START_EPOCH ))

# ── ③ 报告落盘 ──────────────────────────────────────────────────────
_emoji="🟢"
case "$VERDICT" in
    RED) _emoji="🔴" ;;
    YELLOW) _emoji="🟡" ;;
    CANNOT_CERTIFY) _emoji="⚪" ;;
esac

{
    echo "# RC 门禁报告 — $(basename "$PROVIDER_DIR")"
    echo
    echo "- provider_dir: \`$PROVIDER_DIR\`"
    echo "- mode: \`$MODE\`$([ "$QUICK" -eq 1 ] && echo '（quick=plan 为止,未跑 apply,零创建）')"
    echo "- probe: \`$PROBE\`"
    echo "- started(UTC): $STARTED_AT / duration: ${DUR}s"
    echo "- 场景并发: 1（顺跑）/ 总超时预算: ${TOTAL_TIMEOUT}s（0=不限）"
    echo
    echo "## VERDICT: $VERDICT $_emoji  (exit $EXIT_CODE)"
    echo
    echo "> 退码语义:🔴 RED=1 阻断 · 🟡 YELLOW=0 放行但显著标注 · 🟢 GREEN=0 · ⚪ CANNOT_CERTIFY=2 门禁不完整。"
    echo

    echo "## 判定摘要"
    if [ "${#REAL_RED[@]}" -gt 0 ]; then
        echo "### 🔴 RED（阻断,退 1）×${#REAL_RED[@]}"
        for r in "${REAL_RED[@]}"; do echo "- $r"; done
    fi
    if [ "$CANNOT_CERTIFY" -eq 1 ]; then
        echo "### ⚪ CANNOT_CERTIFY（门禁不完整,退 2）"
        for r in "${RED_REASONS[@]}"; do
            case "$r" in __CANNOT_CERTIFY__:*) echo "- ${r#__CANNOT_CERTIFY__:}" ;; esac
        done
    fi
    if [ "${#YELLOW_REASONS[@]}" -gt 0 ]; then
        echo "### 🟡 YELLOW（放行但标注,退 0）×${#YELLOW_REASONS[@]}"
        for y in "${YELLOW_REASONS[@]}"; do echo "- $y"; done
    fi
    if [ "$VERDICT" = "GREEN" ]; then
        echo "### 🟢 GREEN — tier-0 净 + tier-1 全场景过闸,无红无黄。"
    fi
    echo

    echo "## tier-0 三方一致性（probe.sh tier0 --all --limit 200）"
    echo
    echo "- exit: $T0_RC / mech: \`$T0_MECH\` / findings: $T0_FINDINGS / api_gap(S3+): $T0_API_S3PLUS / judgment_queue: ${T0_QUEUE}（黄线阈值 ${QUEUE_YELLOW}）"
    if [ "$T0_RC" -eq 2 ]; then
        echo "- ⚪ tier-0 无法运行(provider 仓不可用),静态扫描缺失 → 门禁不完整。"
    fi
    if [ -n "$T0_API_S3PLUS_LINES" ]; then
        echo
        echo "### 🔴 api_gap S3+ 明细"
        printf '%s\n' "$T0_API_S3PLUS_LINES"
    fi
    if [ -n "$T0_FINDING_LINES" ] && [ "${T0_FINDINGS:-0}" -gt 0 ]; then
        echo
        echo "### findings 全量（${T0_FINDINGS}）"
        printf '%s\n' "$T0_FINDING_LINES"
    fi
    echo
    echo "> 降级容忍:mech=degraded 只记黄不记红——静态 doc↔source diff 仍出,OpenAPI 侧转 judgment_queue 交人工双层查证。"
    echo

    echo "## tier-1 场景生命周期（mode=$MODE, 并发 1）"
    echo
    echo "- 场景数: ${#SCENARIOS[@]} / 通过: $n_pass / 场景fail: $n_fail / destroy残留: $n_destroy / env跳过: $n_env / plan为止: $n_planonly / 超时跳过: $n_timeout"
    if [ "${#SCENARIOS[@]}" -gt 0 ]; then
        echo
        echo "| 场景 | 退码 | 判定 |"
        echo "|------|------|------|"
        printf '%s' "$SCEN_ROWS"
    else
        echo
        echo "- （probe.sh list 未发现任何场景 → tier-1 零覆盖,记黄）"
    fi
    echo

    echo "## 结论与衔接"
    case "$VERDICT" in
        RED) echo "- 🔴 **阻断发版**:先修上列红项(provider finding / destroy 残留 / api_gap S3+),复跑门禁转绿/黄再进入 terraform-provider-release SOP。" ;;
        CANNOT_CERTIFY) echo "- ⚪ **门禁不完整**:tier-0 无法运行,修 provider 仓/环境(website/docs/r + alicloud + probe-meta)后重跑,勿把不可判当绿放行。" ;;
        YELLOW) echo "- 🟡 **可放行但需知情**:上列黄项非硬阻断,建议发版前尽量清;quick 模式请完整模式复跑一遍再发。" ;;
        GREEN) echo "- 🟢 **过闸**:可进入 terraform-provider-release SOP 的 PR/ACC/合并环节。" ;;
    esac
    echo "- 门禁读法 / 触发时机 / 与 terraform-provider-release SOP 衔接位置见 \`loops/tf-probe.md\`「RC 门禁」节。"
    echo
    echo "来源:jarvis F3 RC 门禁线(bootstrap/rc-gate.sh)"
} > "$REPORT"

# ── stdout 摘要 ─────────────────────────────────────────────────────
echo "rc-gate: [3/3] 报告 → $REPORT" >&2
echo "rc-gate: tier-0 mech=$T0_MECH findings=$T0_FINDINGS api_gap(S3+)=$T0_API_S3PLUS queue=$T0_QUEUE (exit $T0_RC)"
echo "rc-gate: tier-1 场景=${#SCENARIOS[@]} pass=$n_pass fail=$n_fail destroy=$n_destroy env=$n_env planonly=$n_planonly timeout=$n_timeout (mode=$MODE)"
if [ "${#REAL_RED[@]}" -gt 0 ]; then
    for r in "${REAL_RED[@]}"; do echo "rc-gate:   RED  · $r"; done
fi
if [ "${#YELLOW_REASONS[@]}" -gt 0 ]; then
    for y in "${YELLOW_REASONS[@]}"; do echo "rc-gate:   YELLOW · $y"; done
fi
echo "VERDICT: $VERDICT"
exit "$EXIT_CODE"
