#!/usr/bin/env bash
# bootstrap/probe.sh — tf-customer-probe 分层探测 runner
#
# 像真实客户一样,参考 provider website/docs 写 HCL、跑 terraform 生命周期,
# 主动发现 terraform-provider-alicloud 潜在且未暴露的问题。
#
# 子命令:
#   doctor                 — 环境预检(terraform/jq/凭证 set-unset/config 可解析);缺 terraform 退非零
#   list                   — 扫 probes/scenarios/*/scenario.yaml 输出表格
#   run <id> [opts]        — 核心:分层跑一个场景,落 verdict.json + runs/probe 审计
#     --tier N   请求 tier(默认取场景声明 tier,再被配置开关封顶)
#     --dry      只解析+打印步骤计划,不需 terraform 二进制(测试主路径)
#     --keep     tier-1 跑完不 destroy(留资源人工排查;慎用)
#   sweep                  — 扫 .my-day/probe/*/terraform.tfstate 报告残留 state
#
# 退出码约定(run):0=无 findings;1=有 findings;2=runner 自身错误/env 阻断;3=清理失败(最高优先级人工介入)。
#
# 分流纪律(建单准确率命门):findings=provider 疑似 bug;env_issues=环境问题(凭证/网络/allowlist/降级)。
#   鉴权/网络类错误永远归 env_issues,绝不进 findings。凭证值绝不落日志/verdict。
#
# 环境变量(多数仅测试用):
#   JARVIS_ROOT           — repo root(见 lib.sh)
#   PROBE_CONFIG          — config/probe.json 路径(默认 <root>/config/probe.json)
#   PROBE_SCENARIOS_DIR   — 场景目录(默认 <root>/probes/scenarios)
#   PROBE_WORKDIR         — 工作/state 目录(默认 <root>/.my-day/probe)
#   PROBE_AUDIT_DIR       — 审计落盘目录(默认 <root>/runs/probe)
#   ALICLOUD_REGION       — 运行 region(缺省用 config.region_fallback)
#
# 被 source 时不执行 main(便于单测内部函数)。
set -uo pipefail

_probe_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/lib.sh
source "$_probe_dir/lib.sh"

# ── 路径 / 配置访问器 ────────────────────────────────────────────────
probe_root()          { jarvis_root; }
probe_config()        { echo "${PROBE_CONFIG:-$(probe_root)/config/probe.json}"; }
probe_scenarios_dir() { echo "${PROBE_SCENARIOS_DIR:-$(probe_root)/probes/scenarios}"; }
probe_workdir_base()  { echo "${PROBE_WORKDIR:-$(probe_root)/.my-day/probe}"; }
probe_audit_dir()     { echo "${PROBE_AUDIT_DIR:-$(probe_root)/runs/probe}"; }

# cfg <jq-filter> — 读一个 config 值(-r 原始输出)
cfg() { jq -r "$1" "$(probe_config)" 2>/dev/null; }

probe_region() { echo "${ALICLOUD_REGION:-$(cfg '.region_fallback')}"; }

# ── scenario.yaml 扁平解析(grep/sed,不引入 python/yq) ───────────────
# _yaml_get <file> <key> — 打印标量值(去首尾空白);缺失打印空串
_yaml_get() {
    local file="$1" key="$2"
    sed -n "s/^${key}:[[:space:]]*//p" "$file" 2>/dev/null | head -1 | sed 's/[[:space:]]*$//'
}

# ── env 判定 ────────────────────────────────────────────────────────
_have_terraform() { command -v terraform >/dev/null 2>&1; }
# 只判断凭证是否 set,绝不读/打印其值
_have_creds() { [ -n "${ALICLOUD_ACCESS_KEY:-}" ] && [ -n "${ALICLOUD_SECRET_KEY:-}" ]; }

# ── allowlist 硬门 ──────────────────────────────────────────────────
# _plan_resource_types <plan.json> — 从 `terraform show -json` plan 抽 managed 资源类型
#   只取 mode==managed(排除 data 源,data 源不创建云资源,不该被 allowlist 拦)。
_plan_resource_types() {
    jq -r '[.resource_changes[]? | select(.mode=="managed") | .type] | unique | .[]' "$1" 2>/dev/null
}

# _allowlist_check <plan.json> — 全部 managed 资源类型 ⊆ tier1_allowlist 则退 0;否则退 1 并 stderr 列出越界项
_allowlist_check() {
    local plan_json="$1" allow t offenders=()
    allow="$(cfg '.tier1_allowlist[]')"
    while IFS= read -r t; do
        [ -z "$t" ] && continue
        if ! grep -qxF "$t" <<<"$allow"; then
            offenders+=("$t")
        fi
    done < <(_plan_resource_types "$plan_json")
    if [ "${#offenders[@]}" -gt 0 ]; then
        echo "allowlist_block: 越界资源 ${offenders[*]}" >&2
        return 1
    fi
    return 0
}

# ── 超时兜底 ────────────────────────────────────────────────────────
# _with_timeout <cmd...> — 有 timeout/gtimeout 则限时,否则跳过限时并 warn 一次(macOS 默认无 timeout)
_with_timeout() {
    local secs
    secs=$(( $(cfg '.limits.step_timeout_min') * 60 ))
    if command -v timeout >/dev/null 2>&1; then
        timeout "${secs}s" "$@"
    elif command -v gtimeout >/dev/null 2>&1; then
        gtimeout "${secs}s" "$@"
    else
        if [ -z "${_PROBE_TIMEOUT_WARNED:-}" ]; then
            echo "WARN: 无 timeout 二进制,步骤超时控制禁用(brew install coreutils 可得 gtimeout)" >&2
            _PROBE_TIMEOUT_WARNED=1
        fi
        "$@"
    fi
}

# ── tier 判定 ───────────────────────────────────────────────────────
# 配置允许的最高 tier:tier1_enabled=true → 1;否则 0。P0 绝不放行 tier-2(硬封顶)。
_config_max_tier() {
    [ "$(cfg '.tiers.tier1_enabled')" = "true" ] && echo 1 || echo 0
}
_min() { [ "$1" -le "$2" ] && echo "$1" || echo "$2"; }

# ── verdict 累加器(jsonl → 末尾 slurp 成数组) ──────────────────────
_emit_step() { # name exit dur log
    jq -nc --arg name "$1" --argjson exit "$2" --argjson dur "$3" --arg log "$4" \
        '{name:$name, exit_code:$exit, duration_s:$dur, log:$log}' >> "$STEPS_FILE"
}
_emit_finding() { # code stage summary evidence severity
    jq -nc --arg code "$1" --arg stage "$2" --arg summary "$3" --arg evidence "$4" --arg sev "$5" \
        '{code:$code, stage:$stage, summary:$summary, evidence:$evidence, severity_hint:$sev}' >> "$FINDINGS_FILE"
    echo "FINDING[$5] $1 @$2: $3" >&2
}
_emit_env() { # code detail
    jq -nc --arg code "$1" --arg detail "$2" '{code:$code, detail:$detail}' >> "$ENV_FILE"
    echo "ENV_ISSUE $1: $2" >&2
}

# ── 单步执行(在 PRUN_WD 内跑命令,输出 append 到 step 日志,回填全局 STEP_RC/STEP_DURATION) ──
# 注意:cd 只包住命令本身的子 shell,STEP_RC=$? 在函数作用域赋值,能传回调用方(全局变量)。
_run_step() { # name logfile -- cmd...
    local name="$1" log="$2"; shift 2; [ "${1:-}" = "--" ] && shift
    local start end
    start=$(date +%s)
    { echo "### probe step: $name :: $(date -u +%FT%TZ)"; echo "\$ $*"; } >> "$log"
    ( cd "$PRUN_WD" && _with_timeout "$@" ) >> "$log" 2>&1
    STEP_RC=$?
    end=$(date +%s)
    STEP_DURATION=$(( end - start ))
    _emit_step "$name" "$STEP_RC" "$STEP_DURATION" "$log"
    echo "[$name] exit=$STEP_RC (${STEP_DURATION}s)"
    return "$STEP_RC"
}

# 日志关键词分类(鉴权/网络永远归 env,不进 findings)
_is_auth_error()    { grep -qiE 'InvalidAccessKeyId|SignatureDoesNotMatch|Forbidden|NoPermission|InvalidSecurityToken|Unauthorized' "$1" 2>/dev/null; }
_is_network_error() { grep -qiE 'registry\.terraform\.io|dial tcp|no such host|i/o timeout|connection refused|TLS handshake|could not (download|query|retrieve)' "$1" 2>/dev/null; }
_has_panic()        { grep -qiE 'panic:|goroutine [0-9]+ \[running\]|runtime error' "$1" 2>/dev/null; }

# ── doctor ──────────────────────────────────────────────────────────
_report_env_flag() { # NAME
    local name="$1"
    if [ -n "${!name:-}" ]; then echo "OK   $name: set"; else echo "WARN $name: unset (tier-0 可跑,tier-1 需要)"; fi
}
_cmd_doctor() {
    local rc=0
    if _have_terraform; then
        echo "OK   terraform: $(terraform version 2>/dev/null | head -1)"
    else
        echo "MISS terraform: 未安装 — brew install hashicorp/tap/terraform"
        rc=1
    fi
    if command -v jq >/dev/null 2>&1; then echo "OK   jq: $(jq --version 2>/dev/null)"; else echo "MISS jq: 未安装"; rc=1; fi
    _report_env_flag ALICLOUD_ACCESS_KEY
    _report_env_flag ALICLOUD_SECRET_KEY
    if jq -e . "$(probe_config)" >/dev/null 2>&1; then
        echo "OK   config: $(probe_config) 可解析"
    else
        echo "MISS config: $(probe_config) 不可解析"
        rc=1
    fi
    return "$rc"
}

# ── list ────────────────────────────────────────────────────────────
_cmd_list() {
    local sdir d y
    sdir="$(probe_scenarios_dir)"
    printf '%-18s %-4s %-9s %-70s %s\n' ID TIER PERSONA RESOURCES DETECT
    shopt -s nullglob
    for d in "$sdir"/*/; do
        y="$d/scenario.yaml"
        [ -f "$y" ] || continue
        printf '%-18s %-4s %-9s %-70s %s\n' \
            "$(_yaml_get "$y" id)" \
            "$(_yaml_get "$y" tier)" \
            "$(_yaml_get "$y" persona)" \
            "$(_yaml_get "$y" resources)" \
            "$(_yaml_get "$y" detect)"
    done
}

# ── run:全局运行态(供 cleanup/finalize 在 trap 后仍可见) ───────────
PRUN_SID=""; PRUN_WD=""; PRUN_RUN_ID=""
PRUN_TIER_REQ=0; PRUN_TIER_EFF=0
PRUN_STARTED_AT=""; PRUN_START_EPOCH=0; PRUN_TF_VERSION="unknown"
PRUN_APPLIED=false; PRUN_KEEP=0; PRUN_DESTROYED=null; PRUN_STATE_EMPTY=null
PRUN_CLEANED=0
STEP_RC=0; STEP_DURATION=0

_usage_run() { echo "用法: probe.sh run <scenario-id> [--tier N] [--dry] [--keep]" >&2; }

# findings 计数码:destroy_fail/state_residue → 3;其它 findings → 1;无 → 0
_verdict_exit() {
    if jq -e -s 'any(.[]; .code=="destroy_fail" or .code=="state_residue")' "$FINDINGS_FILE" >/dev/null 2>&1; then
        echo 3; return
    fi
    if [ -s "$FINDINGS_FILE" ]; then echo 1; else echo 0; fi
}

# tier-1 清理:只要 apply 执行过就 destroy(除非 --keep);幂等(重入即返回)
_probe_cleanup() {
    [ "$PRUN_CLEANED" -eq 1 ] && return 0
    PRUN_CLEANED=1
    [ "$PRUN_APPLIED" = "true" ] || return 0
    if [ "$PRUN_KEEP" -eq 1 ]; then
        echo "--keep: 跳过 destroy,资源保留在 $PRUN_WD(记得手动 destroy)" >&2
        return 0
    fi
    _run_step destroy "$PRUN_WD/destroy.log" -- terraform destroy -auto-approve -var "run_id=$PRUN_RUN_ID"
    if [ "$STEP_RC" -ne 0 ]; then
        PRUN_DESTROYED=false
        _emit_finding destroy_fail cleanup "destroy 失败,可能残留计费资源,需人工清理" "$PRUN_WD/destroy.log" S1
        return 0
    fi
    PRUN_DESTROYED=true
    local remaining
    remaining="$(cd "$PRUN_WD" && terraform state list 2>/dev/null | grep -c . )"
    [ -z "$remaining" ] && remaining=0
    if [ "$remaining" -gt 0 ]; then
        PRUN_STATE_EMPTY=false
        _emit_finding state_residue cleanup "destroy 后 state 仍残留 $remaining 个资源" "$PRUN_WD/destroy.log" S1
    else
        PRUN_STATE_EMPTY=true
    fi
}

# EXIT 安全网:异常退出也兜底 destroy(幂等,不重复 finalize)
_probe_on_exit() { _probe_cleanup; }

_cmd_run() {
    local sid="" tier_req="" dry=0 keep=0
    while [ $# -gt 0 ]; do
        case "$1" in
            --tier) tier_req="${2:-}"; shift 2 ;;
            --dry)  dry=1; shift ;;
            --keep) keep=1; shift ;;
            -*) echo "run: 未知参数 '$1'" >&2; _usage_run; return 2 ;;
            *) if [ -z "$sid" ]; then sid="$1"; shift; else echo "run: 多余参数 '$1'" >&2; return 2; fi ;;
        esac
    done
    [ -n "$sid" ] || { _usage_run; return 2; }

    local sdir yaml
    sdir="$(probe_scenarios_dir)/$sid"
    yaml="$sdir/scenario.yaml"
    if [ ! -f "$yaml" ]; then
        echo "run: 场景不存在 $yaml" >&2
        return 2
    fi

    # tier 判定:effective = min(场景 tier, 请求 tier, 配置允许最高 tier)
    local scenario_tier config_max intent effective downgraded=0
    scenario_tier="$(_yaml_get "$yaml" tier)"; scenario_tier="${scenario_tier:-0}"
    [ -n "$tier_req" ] || tier_req="$scenario_tier"
    config_max="$(_config_max_tier)"
    intent="$(_min "$scenario_tier" "$tier_req")"
    effective="$(_min "$intent" "$config_max")"
    [ "$effective" -lt "$intent" ] && downgraded=1

    local update_step import_check import_address import_id_output
    update_step="$(_yaml_get "$yaml" update_step)"
    import_check="$(_yaml_get "$yaml" import_check)"
    import_address="$(_yaml_get "$yaml" import_address)"
    import_id_output="$(_yaml_get "$yaml" import_id_output)"

    local have_creds=0; _have_creds && have_creds=1

    # ── --dry:只打印计划,不碰 terraform ──
    if [ "$dry" -eq 1 ]; then
        echo "probe plan: scenario=$sid persona=$(_yaml_get "$yaml" persona)"
        echo "  tier: scenario=$scenario_tier requested=$tier_req config_max=$config_max → effective=tier-$effective"
        if [ "$downgraded" -eq 1 ]; then
            echo "  tier_downgraded: 配置 tiers.tier1_enabled=false,intent=tier-$intent 封顶为 tier-$effective"
        fi
        echo "  region: $(probe_region)"
        echo "  steps:"
        echo "    - init"
        echo "    - validate"
        if [ "$have_creds" -eq 1 ]; then
            echo "    - plan -out=tf.plan"
        else
            echo "    - plan: 将跳过(未设置 ALICLOUD 凭证 → env_issue no_creds)"
        fi
        if [ "$effective" -ge 1 ]; then
            echo "    - allowlist-gate: plan managed 资源 ⊆ config.tier1_allowlist(越界 → env_issue allowlist_block,终止)"
            echo "    - apply -auto-approve"
            echo "    - plan -detailed-exitcode(退2 → finding perpetual_diff)"
            [ "$update_step" = "true" ] && echo "    - step2 覆盖 apply + re-plan(抓 更新不生效)"
            [ "$import_check" = "true" ] && echo "    - state rm $import_address → import(id 取自 output $import_id_output) → plan(退2 → finding import_diff)"
            echo "    - destroy -auto-approve + state-empty 校验(--keep 时跳过;destroy 失败/残留 → finding + 退3)"
        else
            echo "  (tier-0: 只静态校验,不创建任何云资源)"
        fi
        return 0
    fi

    # ── 真实执行 ──
    if ! _have_terraform; then
        echo "no_terraform: 本机无 terraform 二进制,无法执行 run(--dry 可无 terraform)。brew install hashicorp/tap/terraform" >&2
        return 2
    fi

    local ts short
    ts="$(date -u +%Y%m%dT%H%M%SZ)"
    short="$(date +%s | tail -c 6)"
    PRUN_SID="$sid"
    PRUN_WD="$(probe_workdir_base)/${ts}-${sid}"
    PRUN_RUN_ID="$(cfg '.name_prefix')-${sid}-${short}"
    PRUN_TIER_REQ="$tier_req"
    PRUN_TIER_EFF="$effective"
    PRUN_KEEP="$keep"
    PRUN_APPLIED=false; PRUN_DESTROYED=null; PRUN_STATE_EMPTY=null; PRUN_CLEANED=0
    mkdir -p "$PRUN_WD"

    # verdict 累加器文件
    STEPS_FILE="$PRUN_WD/.steps.jsonl";      : > "$STEPS_FILE"
    FINDINGS_FILE="$PRUN_WD/.findings.jsonl"; : > "$FINDINGS_FILE"
    ENV_FILE="$PRUN_WD/.env.jsonl";           : > "$ENV_FILE"

    # 拷贝场景 tf(绝不在 probes/ 场景目录里跑 terraform)
    cp "$sdir"/*.tf "$PRUN_WD"/ 2>/dev/null

    # terraform 自动化环境
    export TF_IN_AUTOMATION=1 TF_INPUT=0 TF_CLI_ARGS=-no-color
    TF_PLUGIN_CACHE_DIR="$(probe_workdir_base)/.plugin-cache"; export TF_PLUGIN_CACHE_DIR
    mkdir -p "$TF_PLUGIN_CACHE_DIR"

    [ "$downgraded" -eq 1 ] && _emit_env tier_downgraded "配置 tier1_enabled=false,intent=tier-$intent 封顶为 tier-$effective"

    PRUN_STARTED_AT="$(date -u +%FT%TZ)"; PRUN_START_EPOCH=$(date +%s)
    PRUN_TF_VERSION="$(terraform version 2>/dev/null | head -1)"

    trap _probe_on_exit EXIT

    # ── tier-0 链:init → validate → (有凭证)plan ──
    _run_step init "$PRUN_WD/init.log" -- terraform init -input=false
    if [ "$STEP_RC" -ne 0 ]; then
        if _is_network_error "$PRUN_WD/init.log"; then
            _emit_env init_network_fail "terraform init 失败,日志含 registry/网络关键词"
        else
            _emit_env init_fail "terraform init 失败(非网络),见 $PRUN_WD/init.log"
        fi
        _finalize_verdict; return 2
    fi

    _run_step validate "$PRUN_WD/validate.log" -- terraform validate
    [ "$STEP_RC" -ne 0 ] && _emit_finding validate_fail validate "官方文档示例组合 validate 不通过" "$PRUN_WD/validate.log" S3

    if [ "$have_creds" -eq 0 ]; then
        _emit_env no_creds "未设置 ALICLOUD 凭证,跳过 plan/apply(tier-0 静态校验完成)"
        _finalize_verdict; return "$(_verdict_exit)"
    fi

    _run_step plan "$PRUN_WD/plan.log" -- terraform plan -out=tf.plan -var "run_id=$PRUN_RUN_ID"
    if [ "$STEP_RC" -ne 0 ]; then
        if _is_auth_error "$PRUN_WD/plan.log"; then
            _emit_env auth_error "plan 报鉴权类错误(归 env,不算 provider bug)"
        elif _has_panic "$PRUN_WD/plan.log"; then
            _emit_finding plan_crash plan "terraform plan 触发 provider panic" "$PRUN_WD/plan.log" S1
        else
            _emit_finding plan_fail plan "合法配置 plan 失败" "$PRUN_WD/plan.log" S2
        fi
        _finalize_verdict; return "$(_verdict_exit)"
    fi

    # ── tier-0 到此为止 ──
    if [ "$effective" -lt 1 ]; then
        _finalize_verdict; return "$(_verdict_exit)"
    fi

    # ── tier-1 链 ──
    ( cd "$PRUN_WD" && terraform show -json tf.plan > "$PRUN_WD/plan.json" 2>/dev/null )
    if ! _allowlist_check "$PRUN_WD/plan.json"; then
        _emit_env allowlist_block "plan 含 tier1_allowlist 之外的资源,终止(不 apply)"
        _finalize_verdict; return 2
    fi

    _run_step apply "$PRUN_WD/apply.log" -- terraform apply -auto-approve -var "run_id=$PRUN_RUN_ID"
    PRUN_APPLIED=true
    if [ "$STEP_RC" -ne 0 ]; then
        if _is_auth_error "$PRUN_WD/apply.log"; then
            _emit_env auth_error "apply 报鉴权类错误"
        else
            _emit_finding apply_fail apply "合法配置 apply 失败" "$PRUN_WD/apply.log" S2
        fi
        _probe_cleanup; _finalize_verdict; return "$(_verdict_exit)"
    fi

    # 幂等:apply 后立刻 plan 应无 diff(退2 = 永久 diff)
    _run_step plan2 "$PRUN_WD/plan2.log" -- terraform plan -detailed-exitcode -out=tf.plan2 -var "run_id=$PRUN_RUN_ID"
    [ "$STEP_RC" -eq 2 ] && _emit_finding perpetual_diff plan2 "apply 后立即 plan 仍有 diff(幂等性破坏/永久 diff)" "$PRUN_WD/plan2.log" S2
    _detect_unexpected_replace "$PRUN_WD/plan2.log" plan2

    # step2 覆盖(更新场景)
    if [ "$update_step" = "true" ] && [ -d "$sdir/step2" ]; then
        cp "$sdir"/step2/*.tf "$PRUN_WD"/ 2>/dev/null
        _run_step apply_step2 "$PRUN_WD/apply_step2.log" -- terraform apply -auto-approve -var "run_id=$PRUN_RUN_ID"
        if [ "$STEP_RC" -ne 0 ]; then
            _emit_finding apply_fail apply_step2 "更新(step2) apply 失败" "$PRUN_WD/apply_step2.log" S2
        else
            _run_step plan_step2 "$PRUN_WD/plan_step2.log" -- terraform plan -detailed-exitcode -var "run_id=$PRUN_RUN_ID"
            [ "$STEP_RC" -eq 2 ] && _emit_finding perpetual_diff plan_step2 "step2 更新后 plan 仍有 diff(更新不生效)" "$PRUN_WD/plan_step2.log" S2
        fi
    fi

    # import 断链检查
    if [ "$import_check" = "true" ] && [ -n "$import_address" ] && [ -n "$import_id_output" ]; then
        local imp_id
        imp_id="$(cd "$PRUN_WD" && terraform output -raw "$import_id_output" 2>/dev/null)"
        if [ -n "$imp_id" ]; then
            ( cd "$PRUN_WD" && terraform state rm "$import_address" >> "$PRUN_WD/import.log" 2>&1 )
            _run_step import "$PRUN_WD/import.log" -- terraform import -var "run_id=$PRUN_RUN_ID" "$import_address" "$imp_id"
            if [ "$STEP_RC" -ne 0 ]; then
                _emit_finding import_diff import "terraform import 失败(import 断链)" "$PRUN_WD/import.log" S2
            else
                _run_step plan_import "$PRUN_WD/plan_import.log" -- terraform plan -detailed-exitcode -var "run_id=$PRUN_RUN_ID"
                [ "$STEP_RC" -eq 2 ] && _emit_finding import_diff plan_import "import 后 plan 有 diff(import 未还原完整状态)" "$PRUN_WD/plan_import.log" S2
            fi
        else
            _emit_env import_id_missing "output $import_id_output 为空,跳过 import 检查"
        fi
    fi

    _probe_cleanup
    _finalize_verdict
    return "$(_verdict_exit)"
}

# apply 后 plan JSON 若出现 delete+create(替换重建)→ unexpected_replace(S1)
_detect_unexpected_replace() { # planlog stage
    local stage="$2" pj="$PRUN_WD/plan2.json"
    ( cd "$PRUN_WD" && terraform show -json tf.plan2 > "$pj" 2>/dev/null ) || return 0
    [ -s "$pj" ] || return 0
    if jq -e 'any(.resource_changes[]?; (.change.actions|index("delete")) and (.change.actions|index("create")))' "$pj" >/dev/null 2>&1; then
        _emit_finding unexpected_replace "$stage" "apply 后 plan 出现 delete+create(意外替换重建)" "$pj" S1
    fi
}

# verdict.json 落工作目录 + 复制到 runs/probe + 生成人读 md
_finalize_verdict() {
    local audit; audit="$(probe_audit_dir)"; mkdir -p "$audit"
    local dur=$(( $(date +%s) - PRUN_START_EPOCH ))
    local verdict="$PRUN_WD/verdict.json"
    jq -n \
        --argjson schema 1 \
        --arg sid "$PRUN_SID" \
        --argjson treq "$PRUN_TIER_REQ" \
        --argjson teff "$PRUN_TIER_EFF" \
        --arg pv "$(cfg '.provider.version')" \
        --arg tv "${PRUN_TF_VERSION:-unknown}" \
        --arg region "$(probe_region)" \
        --arg started "$PRUN_STARTED_AT" \
        --argjson dur "$dur" \
        --slurpfile steps "$STEPS_FILE" \
        --slurpfile findings "$FINDINGS_FILE" \
        --slurpfile envs "$ENV_FILE" \
        --argjson applied "$PRUN_APPLIED" \
        --argjson destroyed "$PRUN_DESTROYED" \
        --argjson state_empty "$PRUN_STATE_EMPTY" \
        '{schema_version:$schema, scenario_id:$sid, tier_requested:$treq, tier_effective:$teff,
          provider_version:$pv, terraform_version:$tv, region:$region,
          started_at:$started, duration_s:$dur,
          steps:$steps, findings:$findings, env_issues:$envs,
          cleanup:{applied:$applied, destroyed:$destroyed, state_empty:$state_empty}}' \
        > "$verdict" 2>/dev/null

    local day; day="$(date -u +%Y%m%d)"
    cp "$verdict" "$audit/${day}-${PRUN_SID}.json" 2>/dev/null
    _write_summary_md "$verdict" "$audit/${day}-${PRUN_SID}.md"
    echo "verdict: $verdict"
    echo "audit:   $audit/${day}-${PRUN_SID}.json"
}

_write_summary_md() { # verdict.json out.md
    local v="$1" out="$2"
    {
        echo "# probe verdict: $(jq -r '.scenario_id' "$v")"
        echo
        echo "- provider: $(jq -r '.provider_version' "$v") / terraform: $(jq -r '.terraform_version' "$v")"
        echo "- tier: requested=$(jq -r '.tier_requested' "$v") effective=$(jq -r '.tier_effective' "$v") / region: $(jq -r '.region' "$v")"
        echo "- started: $(jq -r '.started_at' "$v") / duration: $(jq -r '.duration_s' "$v")s"
        echo
        echo "## findings ($(jq '.findings|length' "$v"))"
        jq -r '.findings[]? | "- [\(.severity_hint)] `\(.code)` @\(.stage): \(.summary)\n  evidence: \(.evidence)"' "$v"
        echo
        echo "## env_issues ($(jq '.env_issues|length' "$v"))"
        jq -r '.env_issues[]? | "- `\(.code)`: \(.detail)"' "$v"
        echo
        echo "## cleanup"
        jq -r '.cleanup | "- applied=\(.applied) destroyed=\(.destroyed) state_empty=\(.state_empty)"' "$v"
        echo
        echo "来源:jarvis tf-customer-probe"
    } > "$out"
}

# ── sweep ───────────────────────────────────────────────────────────
_cmd_sweep() {
    local base f n residue=0
    base="$(probe_workdir_base)"
    shopt -s nullglob
    for f in "$base"/*/terraform.tfstate; do
        n="$(jq -r '.resources | length' "$f" 2>/dev/null)"
        [ -z "$n" ] && n=0
        if [ "$n" -gt 0 ]; then
            echo "RESIDUE: $f ($n resources)"
            residue=1
        fi
    done
    if [ "$residue" -eq 1 ]; then
        echo "人工清理:cd 到对应工作目录跑 'terraform destroy',或按 managed_by=jarvis-probe 标签用 aliyun CLI 扫真实孤儿资源。"
        return 1
    fi
    echo "sweep: 无残留 state"
    return 0
}

# ── CLI ─────────────────────────────────────────────────────────────
_usage() { sed -n '2,30p' "$0"; }

main() {
    local cmd="${1:-}"; shift 2>/dev/null || true
    case "$cmd" in
        doctor) _cmd_doctor "$@" ;;
        list)   _cmd_list "$@" ;;
        run)    _cmd_run "$@" ;;
        sweep)  _cmd_sweep "$@" ;;
        -h|--help|"") _usage ;;
        *) echo "probe: 未知命令 '$cmd'" >&2; _usage; exit 2 ;;
    esac
}

[[ "${BASH_SOURCE[0]}" == "$0" ]] && main "$@"
