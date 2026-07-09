#!/usr/bin/env bash
# bootstrap/probe.sh — tf-customer-probe 分层探测 runner
#
# 分层(2026-07-03 重定义):
#   tier-0 = 静态三方一致性扫描(TF 文档 ↔ OpenAPI 文档 ↔ provider 源码),不跑 terraform。
#            机械部分只做本地 文档↔源码 diff;OpenAPI 一侧留 judgment_queue 交给 skill 层查证。
#            范围红线:只测 provider 已接入(TF 已支持)的面,云产品未接入 TF 的资源/参数一律不报 gap。
#   tier-1 = 真实 apply 全生命周期探测(默认开启)。撤销免费白名单成本门,换 PrePaid/Subscription 销毁性守门。
#
# 子命令:
#   doctor                          — 环境预检(terraform/jq/凭证/config/本地 provider 仓/场景根)
#   list                            — 扫 <playground>/<product>/<id>/scenario.yaml 输出表格(含 PRODUCT 列)
#   tier0 [--no-mech] [--all] [--limit N] [--rotate N] [alicloud_xxx ...] [--dry]
#                                   — 静态三方一致性扫描(OpenAPI 机械三方 diff 预筛 + judgment_queue;
#                                     默认扫场景 resources 并集;--all=website/docs/r 全量轮换巡检)
#   run <scenario-id> [--region r] [--dry] [--keep] — tier-1 真实 apply 生命周期探测
#   sweep                           — 扫 .my-day/probe/*/terraform.tfstate 报告残留 state
#   archive [--dry]                 — 幂等归档:draft filed/rejected → drafts_archived/;
#                                     verdict retention → runs/probe/archive/<YYYYMM>/(排 ledger/summary);
#                                     .my-day/probe/<ts-sid> 空 state 且过期 → rm(排 .plugin-cache/manual-*/索引文件);
#                                     plugin-cache 陌生版本报体积;pending drafts + _quarantine + origin:generated 待办清单
#
# 退出码(run/tier0):0=无 findings;1=有 findings;2=runner 自身错误/env 阻断;3=清理失败(run 专属,最高优先级人工介入)。
#
# 分流纪律(建单准确率命门):findings=provider 疑似 bug;env_issues=环境问题(凭证/网络/prepaid/plan-only)。
#   鉴权/网络类错误永远归 env_issues,绝不进 findings。凭证值绝不落日志/verdict。测试账号边界:只用环境注入的测试 AK/SK。
#
# 场景根(playground):场景语料库外置在 jarvis 仓外,按云产品维度两级归档
#   <root>/<product>/<id>/scenario.yaml (product = 一级目录名, e.g. vpc/oss/ram)。
#   解析优先级(env > config > 默认约定):
#     1. $JARVIS_TF_PLAYGROUND 非空且目录存在
#     2. config/probe.json 的 .paths.playground_dir 非 null 且目录存在(应为绝对路径)
#     3. 默认 <jarvis 根目录的父目录>/terraform_playground
#
# 环境变量(多数仅测试用):
#   JARVIS_ROOT              — repo root(见 lib.sh)
#   JARVIS_TF_PLAYGROUND     — 场景根(playground)覆盖,优先级最高(非空且目录存在才生效)
#   PROBE_CONFIG             — config/probe.json 路径(默认 <root>/config/probe.json)
#   PROBE_WORKDIR            — 工作/state 目录(默认 <root>/.my-day/probe)
#   PROBE_AUDIT_DIR          — 审计落盘目录(默认 <root>/runs/probe)
#   JARVIS_PROBE_PROVIDER_DIR — 本地 provider 仓(默认 bootstrap/workspace.sh dir terraform_provider)
#   ALICLOUD_REGION          — region 兜底(优先级最低)
#
# 被 source 时不执行 main(便于单测内部函数)。
set -uo pipefail

_probe_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/lib.sh
source "$_probe_dir/lib.sh"

# ── 路径 / 配置访问器 ────────────────────────────────────────────────
probe_root()          { jarvis_root; }
probe_config()        { echo "${PROBE_CONFIG:-$(probe_root)/config/probe.json}"; }
probe_workdir_base()  { echo "${PROBE_WORKDIR:-$(probe_root)/.my-day/probe}"; }
probe_audit_dir()     { echo "${PROBE_AUDIT_DIR:-$(probe_root)/runs/probe}"; }

# cfg <jq-filter> — 读一个 config 值(-r 原始输出)
cfg() { jq -r "$1" "$(probe_config)" 2>/dev/null; }

# cfg_or <jq-filter> <default> — 读 config,遇 null / 空 / 解析失败回 default。
#   config 分裂防御:所有新键必须有代码内默认值,防止老 PROBE_CONFIG 缺键时 runner 报空。
cfg_or() {
    local v; v="$(jq -r "$1" "$(probe_config)" 2>/dev/null)"
    if [ -z "$v" ] || [ "$v" = "null" ]; then echo "$2"; else echo "$v"; fi
}

# ── A2:台账 + recency 索引 ────────────────────────────────────────
# ledger.jsonl(append-only,本机加速索引;真源在 Aone)。tier0/tier1 finalize + archive 各追加一行。
probe_ledger_file() { echo "$(probe_audit_dir)/ledger.jsonl"; }

# _ledger_append <jq-filter> <--argjson/--arg> ... — 追加一条 JSONL(带 ts+kind)。调用方给 filter+参数。
#   例:_ledger_append '{ts:$ts, kind:"tier0", resources:$r}' --argjson r "$reslist_json"
_ledger_append() {
    local filter="$1"; shift
    local f; f="$(probe_ledger_file)"
    mkdir -p "$(dirname "$f")" 2>/dev/null
    jq -nc --arg ts "$(date -u +%FT%TZ)" "$@" "$filter" >> "$f" 2>/dev/null
}

# .my-day/probe/t1-last-run.json — scenario→epoch;LRU 用。仅在真实执行推进到 plan 完成及以后才更新。
_t1_last_run_file() { echo "$(probe_workdir_base)/t1-last-run.json"; }

# _t1_last_run_update <sid> — 追加 sid→now(合并式;文件不存在则初始化为 {})
_t1_last_run_update() {
    local sid="$1" sf; sf="$(_t1_last_run_file)"
    [ -n "$sid" ] || return 0
    mkdir -p "$(dirname "$sf")" 2>/dev/null
    [ -s "$sf" ] || echo '{}' > "$sf"
    local tmp; tmp="$(mktemp)"
    jq --arg s "$sid" --argjson t "$(date +%s)" '.[$s]=$t' "$sf" > "$tmp" 2>/dev/null && mv -f "$tmp" "$sf" || rm -f "$tmp"
}

# _t1_last_run_get <sid> — 打印 sid 最近 epoch;索引无则扫 runs/probe 文件名回退(新旧两种正则,sid 含连字符禁 naive cut)。
#   返回值:纯整数 epoch(0 = 无)。
_t1_last_run_get() {
    local sid="$1" sf epoch=0 v
    sf="$(_t1_last_run_file)"
    if [ -s "$sf" ]; then
        v="$(jq -r --arg s "$sid" '.[$s] // 0' "$sf" 2>/dev/null)"
        [ -n "$v" ] && [ "$v" != "null" ] && epoch="$v"
    fi
    if [ "$epoch" = "0" ] || [ -z "$epoch" ]; then
        # 文件名回退:兼容 <YYYYMMDD>-<sid>.json 和 <YYYYMMDD>-<HHMMSS>-<sid>.json
        # sid 含连字符,禁 naive cut;直接 basename 剥前缀正则,尾部 -<sid>.json 匹配。
        local audit; audit="$(probe_audit_dir)"
        local latest=0 f base d
        shopt -s nullglob
        for f in "$audit"/*"-${sid}.json"; do
            [ -f "$f" ] || continue
            base="$(basename "$f")"
            # 只接受新旧两种前缀正则(避免其它 sid 前缀误命中):新 8+HHMMSS+"-"+sid.json 或旧 8+"-"+sid.json
            if [[ "$base" =~ ^[0-9]{8}-[0-9]{6}-${sid}\.json$ ]] || [[ "$base" =~ ^[0-9]{8}-${sid}\.json$ ]]; then
                d="$(jq -r '.started_at // empty' "$f" 2>/dev/null)"
                if [ -n "$d" ]; then
                    v="$(date -u -j -f "%Y-%m-%dT%H:%M:%SZ" "$d" +%s 2>/dev/null || date -u -d "$d" +%s 2>/dev/null || echo 0)"
                else
                    v="$(stat -f %m "$f" 2>/dev/null || stat -c %Y "$f" 2>/dev/null || echo 0)"
                fi
                [ -n "$v" ] && [ "$v" -gt "$latest" ] && latest="$v"
            fi
        done
        epoch="$latest"
    fi
    echo "${epoch:-0}"
}

# 场景根(playground)解析。语料库外置 jarvis 仓外(已 git 化:数据仓 tf_playground,直推 master + 工单报备),
#   按云产品维度两级归档 <root>/<product>/<id>/scenario.yaml。
#   优先级:env JARVIS_TF_PLAYGROUND > config paths.playground_dir > workspace.sh dir tf_playground > 默认约定。
probe_playground_dir() {
    if [ -n "${JARVIS_TF_PLAYGROUND:-}" ] && [ -d "$JARVIS_TF_PLAYGROUND" ]; then
        echo "$JARVIS_TF_PLAYGROUND"; return
    fi
    local cfg_dir; cfg_dir="$(cfg '.paths.playground_dir')"
    if [ -n "$cfg_dir" ] && [ "$cfg_dir" != "null" ] && [ -d "$cfg_dir" ]; then
        echo "$cfg_dir"; return
    fi
    # 数据仓 tf_playground 登记(多机 clone 后零配置):workspace.sh 解析且目录在则用之
    local ws_dir; ws_dir="$(bash "$_probe_dir/workspace.sh" dir tf_playground 2>/dev/null | head -1)"
    if [ -n "$ws_dir" ] && [ -d "$ws_dir" ]; then
        echo "$ws_dir"; return
    fi
    echo "$(dirname "$(probe_root)")/terraform_playground"
}

# _find_scenario_dirs <id> — 跨 product 目录按 id 检索场景目录(两级 <root>/<product>/<id>)
#   打印所有命中目录(每行一个,无尾斜杠;0/1/多行由调用方判 冲突)。
_find_scenario_dirs() {
    local id="$1" root d
    root="$(probe_playground_dir)"
    shopt -s nullglob
    for d in "$root"/*/"$id"/; do
        [ -f "$d/scenario.yaml" ] && echo "${d%/}"
    done
}

# _playground_has_scenario <root> — root 下存在任一 <product>/<id>/scenario.yaml 则 0(纯 glob,不依赖 find)
_playground_has_scenario() {
    local root="$1" d
    shopt -s nullglob
    for d in "$root"/*/*/scenario.yaml; do
        [ -f "$d" ] && return 0
    done
    return 1
}

# provider 仓路径:env 覆盖 > workspace.sh 解析
probe_provider_dir() {
    if [ -n "${JARVIS_PROBE_PROVIDER_DIR:-}" ]; then echo "$JARVIS_PROBE_PROVIDER_DIR"; return; fi
    bash "$_probe_dir/workspace.sh" dir terraform_provider 2>/dev/null
}

# region 解析优先级:--region > scenario.yaml region > config regions.focus > 环境 ALICLOUD_REGION
_resolve_region() { # cli_region scenario_region
    local cli="$1" scn="$2" focus
    [ -n "$cli" ] && { echo "$cli"; return; }
    [ -n "$scn" ] && { echo "$scn"; return; }
    focus="$(cfg '.regions.focus')"
    [ -n "$focus" ] && [ "$focus" != "null" ] && { echo "$focus"; return; }
    echo "${ALICLOUD_REGION:-}"
}

# ── scenario.yaml 扁平解析(grep/sed,不引入 python/yq) ───────────────
# _yaml_get <file> <key> — 打印标量值(去首尾空白);缺失打印空串
_yaml_get() {
    local file="$1" key="$2"
    sed -n "s/^${key}:[[:space:]]*//p" "$file" 2>/dev/null | head -1 | sed 's/[[:space:]]*$//'
}

# _apply_disabled <yaml> — scenario.yaml 显式 apply:false(成本安全门/大件资源)→ 0(禁用 apply);
#   缺键或非 false → 非 0(默认放行,apply 是可选键、默认 true)。probe-corpus.sh 命中风险门时写入。
_apply_disabled() { [ "$(_yaml_get "$1" apply)" = "false" ]; }

# ── env 判定 ────────────────────────────────────────────────────────
_have_terraform() { command -v terraform >/dev/null 2>&1; }
# 只判断凭证是否 set,绝不读/打印其值
_have_creds() { [ -n "${ALICLOUD_ACCESS_KEY:-}" ] && [ -n "${ALICLOUD_SECRET_KEY:-}" ]; }

# ── prepaid 守门(替代原免费 allowlist 成本门) ──────────────────────
# 原因不是钱(测试账号)——是包年包月/订阅资源多数无法 API 销毁,会破坏「零残留」纪律。
# _prepaid_check <plan.json> — 扫 planned after 值里 *charge_type/*payment_type 字段;0=干净,1=发现 PrePaid/Subscription
_prepaid_check() {
    local v
    while IFS= read -r v; do
        [ -z "$v" ] && continue
        case "$(printf '%s' "$v" | tr '[:upper:]' '[:lower:]')" in
            prepaid|subscription) return 1 ;;
        esac
    done < <(jq -r '[ .resource_changes[]?.change.after? // {}
                      | to_entries[]?
                      | select(.key|test("(?i)(instance_charge_type|payment_type|charge_type)$"))
                      | .value | select(type=="string") ] | .[]' "$1" 2>/dev/null)
    return 0
}
# _prepaid_should_block <plan.json> <allow_prepaid> — 0=应阻断,1=放行
#   守门关(config prepaid_guard!=true)或场景 allow_prepaid=true → 放行;否则发现 prepaid → 阻断。
_prepaid_should_block() {
    [ "$(cfg '.tiers.tier1.prepaid_guard')" = "true" ] || return 1
    [ "$2" = "true" ] && return 1
    if _prepaid_check "$1"; then return 1; else return 0; fi
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

# 日志关键词分类(鉴权/网络永远归 env,不进 findings)
_is_auth_error()    { grep -qiE 'InvalidAccessKeyId|SignatureDoesNotMatch|Forbidden|NoPermission|InvalidSecurityToken|Unauthorized' "$1" 2>/dev/null; }
_is_network_error() { grep -qiE 'registry\.terraform\.io|dial tcp|no such host|i/o timeout|connection refused|TLS handshake|could not (download|query|retrieve)' "$1" 2>/dev/null; }
_has_panic()        { grep -qiE 'panic:|goroutine [0-9]+ \[running\]|runtime error' "$1" 2>/dev/null; }

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

# ── doctor ──────────────────────────────────────────────────────────
_report_env_flag() { # NAME
    local name="$1"
    if [ -n "${!name:-}" ]; then echo "OK   $name: set"; else echo "WARN $name: unset (tier-1 apply 需要)"; fi
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
    # 场景根(playground):目录存在且至少含 1 个 <product>/<id>/scenario.yaml(缺失只 WARN + 约定/覆盖提示)
    local pg; pg="$(probe_playground_dir)"
    if [ -d "$pg" ] && _playground_has_scenario "$pg"; then
        echo "OK   场景根: $pg (两级布局 <product>/<id>/,含 scenario.yaml)"
    else
        echo "WARN 场景根不可用或无场景: $pg"
        echo "       约定: <jarvis 根目录的父目录>/terraform_playground/<product>/<id>/scenario.yaml"
        echo "       覆盖: 设 JARVIS_TF_PLAYGROUND=<绝对路径> 或 config paths.playground_dir"
    fi
    # 本地 provider 仓(tier-0 静态扫描依赖;缺失只 WARN,不阻断——tier-1 不需要)
    local pdir; pdir="$(probe_provider_dir)"
    if [ -n "$pdir" ] && [ -d "$pdir/website/docs/r" ] && [ -d "$pdir/alicloud" ]; then
        echo "OK   provider 仓: $pdir (website/docs + alicloud 齐备)"
    else
        echo "WARN provider 仓不可解析或缺 website/docs(tier-0 静态扫描不可用): ${pdir:-未解析}"
    fi
    # probe-meta(tier-0 OpenAPI 机械化元数据获取层):不可用=WARN,tier0 自动降级为纯 doc↔source+全 queue
    if bash "$_probe_dir/probe-meta.sh" available >/dev/null 2>&1; then
        echo "OK   probe-meta: OpenAPI 元数据获取层可用(tier-0 机械三方 diff 开启)"
    else
        echo "WARN probe-meta: OpenAPI 元数据获取层不可用(缺 venv/凭证)→ tier-0 自动降级为纯 doc↔source + 全 judgment_queue(现行为)"
        echo "       启用: bash .claude/skills/amp-resource-metadata/scripts/setup.sh + 配 AMP_/ALIBABA_CLOUD_ 凭证(白名单见 skill SKILL.md)"
    fi
    # B2/doctor:aliyun CLI 存在性(drifter 场景 drift_cli 直执 aliyun 需要);WARN 级不阻断——drift 场景默认关。
    if command -v aliyun >/dev/null 2>&1; then
        echo "OK   aliyun CLI: $(aliyun version 2>/dev/null | head -1 || echo installed) (drifter 场景可用)"
    else
        echo "WARN aliyun CLI 未安装(drifter 场景 drift_cli 需要;drift_enabled 默认关不阻断)"
        echo "       安装: https://help.aliyun.com/document_detail/121541.html;drift 场景转正前需在 config 开启 tiers.tier1.drift_enabled"
    fi
    return "$rc"
}

# ── list ────────────────────────────────────────────────────────────
# 两级布局 <playground>/<product>/<id>/:PRODUCT 列取 <id> 的父目录名。
# LAST_RUN 列(A2):追加在**行尾**——rc-gate.sh:146 awk '{print $2}' 取 ID 列,前两列(PRODUCT/ID)不能挪。
#   值:t1-last-run.json 索引优先;缺失时回退扫 runs/probe 文件名(新旧两种正则,含连字符 sid 用 [[ =~ ]] 精确匹配)。
#   格式:ISO UTC(便于人读+机读),从未跑过打 '-'。
_cmd_list() {
    local root d y product sid epoch iso
    root="$(probe_playground_dir)"
    printf '%-10s %-18s %-9s %-64s %-72s %s\n' PRODUCT ID PERSONA RESOURCES DETECT LAST_RUN
    shopt -s nullglob
    for d in "$root"/*/*/; do
        y="$d/scenario.yaml"
        [ -f "$y" ] || continue
        product="$(basename "$(dirname "${d%/}")")"
        sid="$(_yaml_get "$y" id)"
        epoch="$(_t1_last_run_get "$sid")"
        if [ -n "$epoch" ] && [ "$epoch" != "0" ]; then
            iso="$(date -u -r "$epoch" +%FT%TZ 2>/dev/null || date -u -d "@$epoch" +%FT%TZ 2>/dev/null || echo "$epoch")"
        else
            iso="-"
        fi
        printf '%-10s %-18s %-9s %-64s %-72s %s\n' \
            "$product" \
            "$sid" \
            "$(_yaml_get "$y" persona)" \
            "$(_yaml_get "$y" resources)" \
            "$(_yaml_get "$y" detect)" \
            "$iso"
    done
}

# ════════════════════════════════════════════════════════════════════
# tier0 — 静态三方一致性扫描(文档 ↔ 源码机械 diff + OpenAPI judgment_queue)
# ════════════════════════════════════════════════════════════════════

# _parse_doc_args <docfile> — 从 ## Argument Reference 抽顶层参数(到首个 ## / ### 子节为止)
#   输出 TSV: name<TAB>req(Required|Optional|?)<TAB>forcenew(0|1)<TAB>deprecated(0|1)
_parse_doc_args() {
    awk '
    /^## Argument Reference/ { inarg=1; inattr=0; next }
    /^## Attributes? Reference/ { inarg=0; inattr=1; next }
    (inarg || inattr) && /^###? / { inarg=0; inattr=0; next }
    inarg && /^\* `/ {
        line=$0
        name=line; sub(/^\* `/,"",name); sub(/`.*/,"",name)
        flags=""
        if (match(line, /\([^)]*\)/)) flags=substr(line, RSTART+1, RLENGTH-2)
        req = (flags ~ /Required/) ? "Required" : ((flags ~ /Optional/) ? "Optional" : "?")
        fn  = (flags ~ /ForceNew/) ? 1 : 0
        dep = (flags ~ /Deprecated/) ? 1 : 0
        print name "\t" req "\t" fn "\t" dep
    }' "$1" 2>/dev/null
}

# _parse_doc_attrs <docfile> — 从 ## Attribute(s) Reference 抽导出属性名(每行一个)
_parse_doc_attrs() {
    awk '
    /^## Argument Reference/ { inarg=1; inattr=0; next }
    /^## Attributes? Reference/ { inarg=0; inattr=1; next }
    (inarg || inattr) && /^###? / { inarg=0; inattr=0; next }
    inattr && /^\* `/ {
        name=$0; sub(/^\* `/,"",name); sub(/`.*/,"",name); print name
    }' "$1" 2>/dev/null
}

# _parse_source_schema <srcfile> — 抽顶层 Schema 键 + 标志(按大括号深度追踪,只顶层)
#   局限(P0):嵌套 Elem 内层字段跳过;首个 map[string]*schema.Schema{} 视为资源顶层 schema。P1 再深挖。
#   输出 TSV: name<TAB>optional(0|1)<TAB>required(0|1)<TAB>forcenew(0|1)<TAB>computed(0|1)<TAB>deprecated(0|1)
_parse_source_schema() {
    awk '
    function pf(k, line) {
        if (line ~ /Optional:[ \t]*true/)   opt[k]=1
        if (line ~ /Required:[ \t]*true/)   reqd[k]=1
        if (line ~ /ForceNew:[ \t]*true/)   fn[k]=1
        if (line ~ /Computed:[ \t]*true/)   comp[k]=1
        if (line ~ /Deprecated:/)           dep[k]=1
    }
    {
        # 在副本上数括号,$0 不动(避免影响后续匹配)
        t=$0; o=gsub(/\{/,"\\&",t); t=$0; c=gsub(/\}/,"\\&",t)
        if (!inschema) {
            if ($0 ~ /map\[string\]\*schema\.Schema\{/) { inschema=1; depth=o-c; base=depth }
            next
        }
        sd = depth
        if (sd == base && $0 ~ /^[ \t]*"[^"]+"[ \t]*:/) {
            name=$0; sub(/^[ \t]*"/,"",name); sub(/".*/,"",name)
            curkey=name
            if (!(name in seen)) { order[++n]=name; seen[name]=1 }
            pf(name, $0)
        } else if (curkey != "" && sd == base+1) {
            pf(curkey, $0)
        }
        depth = sd + o - c
        if (depth < base) inschema=0
    }
    END {
        for (i=1;i<=n;i++){ k=order[i]
            printf "%s\t%d\t%d\t%d\t%d\t%d\n", k, opt[k]+0, reqd[k]+0, fn[k]+0, comp[k]+0, dep[k]+0 }
    }' "$1" 2>/dev/null
}

# _count_lines <file> — 非空行数(grep -c 在 0 匹配时退非零,单独封装避免 || echo 0 双打印)
_count_lines() { local n; n=$(grep -c . "$1" 2>/dev/null); echo "${n:-0}"; }

# _source_api_actions <srcfile> — best-effort 抽 source 里实际调用的 API action 名
#   两路并集:`action :?= "X"` 赋值 + RpcPost/RpcGet 第三位字面量 action。
_source_api_actions() {
    { grep -oE 'action[[:space:]]*:?=[[:space:]]*"[A-Za-z0-9]+"' "$1" 2>/dev/null | grep -oE '"[A-Za-z0-9]+"' | tr -d '"'
      grep -oE 'Rpc(Post|Get)\("[A-Za-z0-9]+", *"[0-9][0-9-]*", *"[A-Za-z0-9]+"' "$1" 2>/dev/null | awk -F'"' '{print $6}'
    } | grep . | sort -u
}

# ── T0-mech:动作三元组 + 源码约束解析 + API 规范化 + 机械 diff ──────────
# _snake_to_camel <snake> — snake_case → CamelCase(v1 精确映射规则,不做别名/缩写猜测)
_snake_to_camel() {
    local s="$1" out="" seg oldifs="$IFS"
    IFS='_'
    for seg in $s; do
        [ -z "$seg" ] && continue
        out="$out$(printf '%s' "${seg:0:1}" | tr '[:lower:]' '[:upper:]')${seg:1}"
    done
    IFS="$oldifs"
    printf '%s' "$out"
}

# _source_pv <srcfile> — 抽 RpcPost/RpcGet 的 (product, version) 唯一集,TSV product<TAB>version
_source_pv() {
    grep -oE 'Rpc(Post|Get)\("[A-Za-z0-9]+", *"[0-9][0-9-]*"' "$1" 2>/dev/null \
        | awk -F'"' '{print $2 "\t" $4}' | sort -u
}

# _source_api_triples <srcfile> — (product,version,action) 三元组 TSV。
#   仅当 (product,version) 唯一时输出(0 或 >1 → 空,交 queue,机械层不猜)。OSS 类 SDK 风格无 RpcPost → 空。
_source_api_triples() {
    local src="$1" pv npv product version a
    pv="$(_source_pv "$src")"
    npv="$(printf '%s\n' "$pv" | grep -c .)"
    [ "$npv" = "1" ] || return 0
    product="$(printf '%s' "$pv" | cut -f1)"
    version="$(printf '%s' "$pv" | cut -f2)"
    while IFS= read -r a; do
        [ -n "$a" ] && printf '%s\t%s\t%s\n' "$product" "$version" "$a"
    done < <(_source_api_actions "$src")
}

# _parse_source_constraints <srcfile> — 顶层 schema 键的类型/枚举/范围/默认值(与 _parse_source_schema 同款
#   顶层深度追踪,嵌套 Elem 内层跳过)。解析不动一律标 unknown(进 queue,不猜)。
#   输出 TSV: name  type  enum(US-joined)  enum_status(known|unknown|none)  min  max  range_status(known|unknown|none)  default  default_status(known|unknown|none)
_parse_source_constraints() {
    awk '
    function pf(k, line,   m,seg,rest,tok,e,d,arr){
        if (line ~ /Type:[ \t]*schema\.Type[A-Za-z]+/ && stype[k]==""){
            m=line; sub(/.*schema\.Type/,"",m); sub(/[^A-Za-z].*/,"",m); stype[k]=tolower(m)
        }
        if (line ~ /StringInSlice\(/){
            if (match(line, /\[\]string\{[^}]*\}/)){
                seg=substr(line,RSTART,RLENGTH); rest=seg; e=""
                while (match(rest, /"[^"]*"/)){ tok=substr(rest,RSTART+1,RLENGTH-2); e=e tok SEP; rest=substr(rest,RSTART+RLENGTH) }
                senum[k]=e; senumst[k]="known"
            } else { if (senumst[k]!="known") senumst[k]="unknown" }
        }
        if (match(line, /IntBetween\([ \t]*-?[0-9]+[ \t]*,[ \t]*-?[0-9]+[ \t]*\)/)){
            m=substr(line,RSTART,RLENGTH); gsub(/[^0-9,-]/,"",m); split(m,arr,","); smin[k]=arr[1]; smax[k]=arr[2]; srangest[k]="known"
        } else if (line ~ /IntBetween\(/){ if (srangest[k]!="known") srangest[k]="unknown" }
        if (match(line, /(^|[^A-Za-z_])Default:[ \t]*/)){
            d=line; sub(/.*Default:[ \t]*/,"",d); sub(/,[ \t]*$/,"",d); sub(/[ \t]*$/,"",d)
            if (d ~ /\(/ || d==""){ if (sdefst[k]!="known") sdefst[k]="unknown" }
            else { gsub(/^"/,"",d); gsub(/"$/,"",d); sdef[k]=d; sdefst[k]="known" }
        }
    }
    BEGIN{ SEP=sprintf("%c",31) }
    {
        t=$0; o=gsub(/\{/,"\\&",t); t=$0; c=gsub(/\}/,"\\&",t)
        if (!inschema) {
            if ($0 ~ /map\[string\]\*schema\.Schema\{/) { inschema=1; depth=o-c; base=depth }
            next
        }
        sd = depth
        if (sd == base && $0 ~ /^[ \t]*"[^"]+"[ \t]*:/) {
            name=$0; sub(/^[ \t]*"/,"",name); sub(/".*/,"",name)
            curkey=name
            if (!(name in seen)) { order[++n]=name; seen[name]=1 }
            pf(name, $0)
        } else if (curkey != "" && sd == base+1) {
            pf(curkey, $0)
        }
        depth = sd + o - c
        if (depth < base) inschema=0
    }
    END {
        for (i=1;i<=n;i++){ k=order[i]
            printf "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", k,
                (stype[k]==""?"unknown":stype[k]),
                senum[k], (senumst[k]==""?"none":senumst[k]),
                smin[k], smax[k], (srangest[k]==""?"none":srangest[k]),
                sdef[k], (sdefst[k]==""?"none":sdefst[k]) }
    }' "$1" 2>/dev/null
}

# _api_extract_params <apidef.json> — 规范化 API 参数为 TSV。防御式多 key 兜底(enum/enumValueList/enumValueTitles)。
#   输出: name  type  required(1|0|?)  enum(US-joined)  enum_status(known|none)  min  max  range_status  default  default_status  deprecated(1|0)
_api_extract_params() {
    jq -r '
      def sep: "";
      .parameters[]?
      | (.in // "query") as $in
      | select($in=="query" or $in=="body" or $in=="path" or $in=="formData")
      | . as $p | (.schema // {}) as $s
      | ( ( $s.enum
            // (($s.enumValueList // []) | map(if type=="object" then (.value // .name // .) else . end))
            // (($s.enumValueTitles // {}) | keys) ) // [] ) as $enum
      | ($s.minimum // $s.min) as $mn
      | ($s.maximum // $s.max) as $mx
      | [ ($p.name // ""),
          (($s.type // $p.type // "") | ascii_downcase),
          ( if ($p.required==true or $s.required==true) then "1"
            elif (($p|has("required")) or ($s|has("required"))) then "0" else "?" end ),
          ( $enum | map(tostring) | join(sep) ),
          ( if ($enum|length)>0 then "known" else "none" end ),
          ( if $mn==null then "" else ($mn|tostring) end ),
          ( if $mx==null then "" else ($mx|tostring) end ),
          ( if ($mn!=null or $mx!=null) then "known" else "none" end ),
          ( if $s.default==null then "" else ($s.default|tostring) end ),
          ( if $s.default!=null then "known" else "none" end ),
          ( if ($p.deprecated==true or $s.deprecated==true) then "1" else "0" end )
        ] | @tsv
    ' "$1" 2>/dev/null
}

# _api_action_deprecated <apidef.json> — action 级 deprecated 判定:true|false|unknown
_api_action_deprecated() {
    local v; v="$(jq -r 'if .deprecated==true then "true" elif .deprecated==false then "false" else "unknown" end' "$1" 2>/dev/null)"
    printf '%s' "${v:-unknown}"
}

# _tier0_mech_param_diff <api_params.tsv> <src_schema.tsv> <src_constraints.tsv> <supp> <tol>
#   逐 TF 顶层输入参数 ↔ API 参数机械 diff。输出 tagged 行(shell 路由):
#     F<TAB>code<TAB>attr<TAB>sev<TAB>summary   S<TAB>param<TAB>apiparam<TAB>rule
#     C<TAB>attr<TAB>note                        U<TAB>attr(unmapped)   E<TAB>attr(enum_unparsed)
_tier0_mech_param_diff() {
    awk -F'\t' -v apif="$1" -v schf="$2" -v conf="$3" -v supp="$4" -v tol="$5" '
    BEGIN{ SEP=sprintf("%c",31) }
    function in_set(joined,val,   arr,i,m){ m=split(joined,arr,SEP); for(i=1;i<=m;i++) if(arr[i]==val) return 1; return 0 }
    function to_camel(s,   parts,i,out,seg,m){ m=split(s,parts,"_"); out=""; for(i=1;i<=m;i++){seg=parts[i]; if(seg=="")continue; out=out toupper(substr(seg,1,1)) substr(seg,2)} return out }
    function type_ok(tf,api){
        if(tf=="string"&&api=="string")return 1
        if((tf=="int"||tf=="float")&&(api=="integer"||api=="long"||api=="number"||api=="float"))return 1
        if(tf=="bool"&&api=="boolean")return 1
        if((tf=="list"||tf=="set")&&api=="array")return 1
        if(tf=="map"&&(api=="object"||api=="map"))return 1
        if(index(tol, tf ">" api)>0)return 1
        return 0
    }
    FILENAME==apif { aname[$1]=1; atype[$1]=$2; areq[$1]=$3; aenum[$1]=$4; aenumst[$1]=$5; amin[$1]=$6; amax[$1]=$7; arangest[$1]=$8; adef[$1]=$9; adefst[$1]=$10; next }
    FILENAME==schf { sord[++sn]=$1; sopt[$1]=$2; sreqd[$1]=$3; scomp[$1]=$5; sdep[$1]=$6; next }
    FILENAME==conf { ctype[$1]=$2; cenum[$1]=$3; cenumst[$1]=$4; cmin[$1]=$5; cmax[$1]=$6; crangest[$1]=$7; cdef[$1]=$8; cdefst[$1]=$9; next }
    END{
        for(i=1;i<=sn;i++){ k=sord[i]; if((sopt[k]==1||sreqd[k]==1)&&sdep[k]!=1){ mapcount[to_camel(k)]++ } }
        for(i=1;i<=sn;i++){
            k=sord[i]
            if(!(sopt[k]==1||sreqd[k]==1)) continue
            if(sdep[k]==1) continue
            camel=to_camel(k)
            if(index(supp, ":" camel ":")>0){ print "S\t" k "\t" camel "\tsuppress_params"; continue }
            if(!(camel in aname)){ print "U\t" k; continue }
            if(mapcount[camel]>1){ print "U\t" k; continue }
            if(cenumst[k]=="known" && aenumst[camel]=="known"){
                extra=""; m=split(cenum[k],sv,SEP)
                for(j=1;j<=m;j++){ if(sv[j]=="")continue; if(!in_set(aenum[camel],sv[j])) extra=extra sv[j] "," }
                if(extra!=""){ sub(/,$/,"",extra); print "F\tapi_gap_enum_superset\t" k "\tS3\tTF 枚举放行 API 拒绝的值 {" extra "}(客户端过宽,API 必拒)" }
                miss=""; m2=split(aenum[camel],av,SEP)
                for(j=1;j<=m2;j++){ if(av[j]=="")continue; if(!in_set(cenum[k],av[j])) miss=miss av[j] "," }
                if(miss!=""){ sub(/,$/,"",miss); print "C\t" k "\tTF 枚举比 API 更严(未含 API 值 {" miss "});方向安全,仅记录" }
            } else if(cenumst[k]=="unknown"){ print "E\t" k }
            if(sopt[k]==1 && sreqd[k]==0 && scomp[k]==0 && cdefst[k]=="none" && areq[camel]=="1")
                print "F\tapi_gap_required\t" k "\tS3\tTF 标 Optional 但 API 要求必填(无 Default/Computed 兜底)"
            if(ctype[k]!="" && ctype[k]!="unknown" && atype[camel]!="" && !type_ok(ctype[k],atype[camel]))
                print "F\tapi_gap_type\t" k "\tS3\tTF 类型 " ctype[k] " 与 API 类型 " atype[camel] " 硬冲突(经容差表过滤后仍冲突)"
            if(crangest[k]=="known" && arangest[camel]=="known"){
                bad=0
                if(amin[camel]!="" && (cmin[k]+0)<(amin[camel]+0)) bad=1
                if(amax[camel]!="" && (cmax[k]+0)>(amax[camel]+0)) bad=1
                if(bad) print "F\tapi_gap_range\t" k "\tS4\tTF 范围 [" cmin[k] "," cmax[k] "] 越过 API [" amin[camel] "," amax[camel] "]"
            }
            if(cdefst[k]=="known" && adefst[camel]=="known" && cdef[k]!=adef[camel])
                print "F\tapi_gap_default\t" k "\tS4\tTF 默认值 " cdef[k] " 与 API 默认值 " adef[camel] " 冲突"
        }
    }' "$1" "$2" "$3" 2>/dev/null
}

# ── LRU 轮换(--rotate):状态落 .my-day/probe/t0mech-scanned.json(res→last epoch) ──
_rotate_state_file() { echo "${PROBE_ROTATE_STATE:-$(probe_workdir_base)/t0mech-scanned.json}"; }
# _rotate_select <state_file> <N> <res...> — 打印 N 个最久未扫资源(least-recently-scanned 优先)
_rotate_select() {
    local sf="$1" n="$2"; shift 2
    local res ts
    for res in "$@"; do
        ts="$(jq -r --arg r "$res" '.[$r] // 0' "$sf" 2>/dev/null)"; [ -z "$ts" ] && ts=0
        printf '%s\t%s\n' "$ts" "$res"
    done | sort -n -k1,1 | head -n "$n" | cut -f2
}
# _rotate_mark <state_file> <res...> — 记录本轮已扫资源时间戳
_rotate_mark() {
    local sf="$1"; shift
    local now res tmp; now="$(date +%s)"
    mkdir -p "$(dirname "$sf")"
    [ -s "$sf" ] || echo '{}' > "$sf"
    for res in "$@"; do
        tmp="$(mktemp)"
        jq --arg r "$res" --argjson t "$now" '.[$r]=$t' "$sf" > "$tmp" 2>/dev/null && mv -f "$tmp" "$sf" || rm -f "$tmp"
    done
}

# _tier0_diff <docargs.tsv> <srckeys.tsv> — 五类 gap → TSV: code<TAB>attr<TAB>severity<TAB>summary(已按 code,attr 排序)
_tier0_diff() {
    awk -F'\t' '
    FNR==NR { darg[$1]=1; dreq[$1]=$2; dfn[$1]=$3; ddep[$1]=$4; next }
    {
        sall[$1]=1; sopt[$1]=$2; sreqd[$1]=$3; sfn[$1]=$4; scomp[$1]=$5; sdep[$1]=$6
        sinput[$1] = ($2==1 || $3==1) ? 1 : 0
    }
    END {
        for (n in darg) if (!(n in sall))
            print "doc_gap_phantom\t" n "\tS3\t文档有参数但源码 schema 无(文档幻影参数)"
        for (n in sall) if (sinput[n] && !(n in darg))
            print "doc_gap_undocumented\t" n "\tS3\t源码有参数但文档未记(未文档化参数)"
        # flag/forcenew 只查活跃(非废弃)字段——废弃字段的标注差异非可行动噪声
        for (n in darg) if ((n in sall) && !sdep[n] && !ddep[n]) {
            sreq = sreqd[n] ? "Required" : (sopt[n] ? "Optional" : "Computed")
            if (dreq[n] != "?" && sreq != "Computed" && dreq[n] != sreq)
                print "doc_gap_flag_mismatch\t" n "\tS3\t文档标 " dreq[n] " 但源码标 " sreq
            if (dfn[n] != sfn[n])
                print "doc_gap_forcenew\t" n "\tS2\t文档 ForceNew=" dfn[n] " 但源码 ForceNew=" sfn[n] "(可能意外重建)"
        }
        for (n in darg) if ((n in sall) && sdep[n] && !ddep[n])
            print "doc_gap_deprecated\t" n "\tS4\t源码已 Deprecated 但文档仍作正常参数列出(未标注废弃)"
    }' "$1" "$2" 2>/dev/null | sort
}

# tier0 findings/judgment 累加器
_emit_t0finding() { # code resource attribute summary severity
    jq -nc --arg code "$1" --arg res "$2" --arg attr "$3" --arg summary "$4" --arg sev "$5" \
        '{code:$code, resource:$res, attribute:$attr, summary:$summary, severity_hint:$sev}' >> "$FINDINGS_FILE"
    echo "FINDING[$5] $1 $2.$3: $4" >&2
}

# _tier0_mech_resource <res> <srcf> <tmp> <supp> <tol> — mech-on 单资源 OpenAPI 侧机械核验。
#   动作三元组抽取 → 逐 action deprecated_action 检查 → 创建型 action 参数级机械 diff →
#   路由 finding / suppressed[] / coverage_note / judgment_queue(带 reason)。用全局累加文件。
_tier0_mech_resource() {
    local res="$1" srcf="$2" tmp="$3" supp="$4" tol="$5"
    local triples actions_json
    triples="$(_source_api_triples "$srcf")"
    actions_json="$(_source_api_actions "$srcf" | jq -R . | jq -s .)"
    if [ -z "$triples" ]; then
        local npv reason note
        npv="$(_source_pv "$srcf" | grep -c .)"
        if [ "$npv" -gt 1 ]; then
            reason="ambiguous_triple"; note="源码含多个 (product,version),无法确定 action 归属;OpenAPI 侧交人工核。"
        else
            reason="no_triple"; note="源码非 RpcPost 风格,抽不到 (product,version) 三元组(如 OSS SDK);OpenAPI 侧交人工核。"
        fi
        jq -nc --arg res "$res" --arg reason "$reason" --argjson actions "$actions_json" --arg note "$note" \
            '{resource:$res,reason:$reason,api_actions:$actions,note:$note,detail:[]}' >> "$JQUEUE_FILE"
        return
    fi
    # 逐 action:deprecated_action 检查(拉取失败记 metafail,不报 finding)
    local tp tv ta defjson dep metafail=""
    while IFS=$'\t' read -r tp tv ta; do
        [ -z "$ta" ] && continue
        defjson="$(bash "$_probe_dir/probe-meta.sh" cached-fetch "$tp" "$tv" "$ta" </dev/null 2>/dev/null)"
        if [ -z "$defjson" ]; then metafail="$metafail $ta"; continue; fi
        printf '%s' "$defjson" > "$tmp/apidef_${ta}.json"
        dep="$(_api_action_deprecated "$tmp/apidef_${ta}.json")"
        [ "$dep" = "true" ] && _emit_t0finding api_gap_deprecated_action "$res" "$ta" \
            "源码调用的 action $ta 在 OpenAPI 标 deprecated(随上游退役,TF 侧无提示)" S3
    done <<< "$triples"

    # 创建型 action(优先 Create*)→ 参数级机械 diff
    local create_action
    create_action="$(printf '%s\n' "$triples" | cut -f3 | grep -E '^Create' | head -1)"
    [ -z "$create_action" ] && create_action="$(printf '%s\n' "$triples" | cut -f3 | grep -E '^(Run|Allocate|Add|Apply|Register|Assign|Authorize|Put)' | head -1)"

    local unmapped="" enumunp="" tag a b c d
    if [ -n "$create_action" ] && [ -f "$tmp/apidef_${create_action}.json" ]; then
        _api_extract_params "$tmp/apidef_${create_action}.json" > "$tmp/api.tsv"
        _parse_source_schema "$srcf" > "$tmp/msch.tsv"
        _parse_source_constraints "$srcf" > "$tmp/mcon.tsv"
        while IFS=$'\t' read -r tag a b c d; do
            case "$tag" in
                F) _emit_t0finding "$a" "$res" "$b" "$d" "$c" ;;
                S) jq -nc --arg res "$res" --arg p "$a" --arg ap "$b" --arg rule "$c" \
                       '{resource:$res,param:$p,api_param:$ap,rule:$rule}' >> "$SUPPRESS_FILE" ;;
                C) jq -nc --arg res "$res" --arg attr "$a" --arg note "$b" \
                       '{resource:$res,attribute:$attr,note:$note}' >> "$COVERAGE_FILE" ;;
                U) unmapped="$unmapped $a" ;;
                E) enumunp="$enumunp $a" ;;
            esac
        done < <(_tier0_mech_param_diff "$tmp/api.tsv" "$tmp/msch.tsv" "$tmp/mcon.tsv" "$supp" "$tol")
    else
        metafail="$metafail ${create_action:-<无创建型 action>}"
    fi

    # queue:prose_review(残留人工核)+ 可选 unmapped/enum_unparsed/meta_unavailable
    jq -nc --arg res "$res" --argjson actions "$actions_json" \
        --arg note "机械层已核 deprecated-action/enum/required/type/range/default;此处仅需人工核 prose 约束(长度/字符集/基数)与行为一致性。范围红线=只核已接入面。" \
        '{resource:$res,reason:"prose_review",api_actions:$actions,note:$note,detail:[]}' >> "$JQUEUE_FILE"
    local dj
    if [ -n "${unmapped// /}" ]; then
        dj="$(printf '%s\n' $unmapped | grep . | jq -R . | jq -s .)"
        jq -nc --arg res "$res" --argjson detail "$dj" \
            --arg note "以下 TF 参数 snake→Camel 未命中 API 参数(convert 改名/仅 Modify 用/未接入),交人工核,机械层不猜。" \
            '{resource:$res,reason:"unmapped_params",api_actions:[],note:$note,detail:$detail}' >> "$JQUEUE_FILE"
    fi
    if [ -n "${enumunp// /}" ]; then
        dj="$(printf '%s\n' $enumunp | grep . | jq -R . | jq -s .)"
        jq -nc --arg res "$res" --argjson detail "$dj" \
            --arg note "以下参数 ValidateFunc 枚举非字面 slice,机械层不解析,交人工核。" \
            '{resource:$res,reason:"enum_unparsed",api_actions:[],note:$note,detail:$detail}' >> "$JQUEUE_FILE"
    fi
    if [ -n "${metafail// /}" ]; then
        dj="$(printf '%s\n' $metafail | grep . | jq -R . | jq -s .)"
        jq -nc --arg res "$res" --argjson detail "$dj" \
            --arg note "以下 action 元数据拉取失败,已跳过其机械核验,交人工核。" \
            '{resource:$res,reason:"meta_unavailable",api_actions:[],note:$note,detail:$detail}' >> "$JQUEUE_FILE"
    fi
}

_cmd_tier0() {
    local dry=0 nomech=0 allres=0 limit=0 rotate=0; local reslist=()
    while [ $# -gt 0 ]; do
        case "$1" in
            --dry) dry=1; shift ;;
            --no-mech) nomech=1; shift ;;
            --all) allres=1; shift ;;
            --limit) limit="${2:-0}"; shift 2 ;;
            --rotate) rotate="${2:-0}"; shift 2 ;;
            -*) echo "tier0: 未知参数 '$1'" >&2; return 2 ;;
            *) reslist+=("$1"); shift ;;
        esac
    done

    local pdir; pdir="$(probe_provider_dir)"
    if [ -z "$pdir" ] || [ ! -d "$pdir/website/docs/r" ] || [ ! -d "$pdir/alicloud" ]; then
        echo "tier0: 本地 provider 仓不可用(需 website/docs/r + alicloud): ${pdir:-未解析}" >&2
        echo "  设 JARVIS_PROBE_PROVIDER_DIR 或 bootstrap/workspace.sh dir terraform_provider 可解析后重试。" >&2
        return 2
    fi

    # 资源清单:--all(website/docs/r 全量) > 显式参数 > 场景 resources 并集
    if [ "$allres" -eq 1 ]; then
        local f n; shopt -s nullglob; local allacc=()
        for f in "$pdir"/website/docs/r/*.html.markdown; do
            n="$(basename "$f" .html.markdown)"; allacc+=("alicloud_$n")
        done
        reslist=(); local r
        while IFS= read -r r; do [ -n "$r" ] && reslist+=("$r"); done < <(printf '%s\n' "${allacc[@]}" | grep . | sort -u)
    elif [ "${#reslist[@]}" -eq 0 ]; then
        local root d y r
        root="$(probe_playground_dir)"
        shopt -s nullglob
        local acc=""
        for d in "$root"/*/*/; do
            y="$d/scenario.yaml"; [ -f "$y" ] || continue
            acc="$acc,$(_yaml_get "$y" resources)"
        done
        while IFS= read -r r; do [ -n "$r" ] && reslist+=("$r"); done \
            < <(printf '%s' "$acc" | tr ',' '\n' | sed 's/[[:space:]]//g' | grep . | sort -u)
    fi
    [ "${#reslist[@]}" -gt 0 ] || { echo "tier0: 无资源可扫(场景 resources 为空且未传参)" >&2; return 2; }

    # --rotate:least-recently-scanned 选 N(状态落 .my-day/probe/t0mech-scanned.json)
    local rotate_sf; rotate_sf="$(_rotate_state_file)"
    if [ "$rotate" -gt 0 ]; then
        local rsel=() rr
        while IFS= read -r rr; do [ -n "$rr" ] && rsel+=("$rr"); done < <(_rotate_select "$rotate_sf" "$rotate" "${reslist[@]}")
        reslist=("${rsel[@]}")
    fi
    # --limit:截断
    if [ "$limit" -gt 0 ] && [ "${#reslist[@]}" -gt "$limit" ]; then
        reslist=("${reslist[@]:0:$limit}")
    fi

    # 机械层模式:--no-mech → off;probe-meta 可用 → on;否则自动 degraded(纯 doc↔source + 全 queue)
    local mech
    if [ "$nomech" -eq 1 ]; then mech="off"
    elif bash "$_probe_dir/probe-meta.sh" available >/dev/null 2>&1; then mech="on"
    else mech="degraded"; fi

    if [ "$dry" -eq 1 ]; then
        echo "tier0 plan: provider_dir=$pdir  mech=$mech"
        echo "  将扫描资源(文档↔源码机械 diff + OpenAPI $([ "$mech" = on ] && echo '机械三方 diff' || echo 'judgment_queue')):"
        local res name docf srcf
        for res in "${reslist[@]}"; do
            name="${res#alicloud_}"
            docf="$pdir/website/docs/r/${name}.html.markdown"
            srcf="$pdir/alicloud/resource_${res}.go"
            echo "    - $res  doc:$([ -f "$docf" ] && echo ok || echo MISSING)  src:$([ -f "$srcf" ] && echo ok || echo MISSING)"
        done
        [ "$mech" != "on" ] && echo "  (mech=$mech:OpenAPI 侧全部进 judgment_queue,交 skill 层人工双层查证)"
        echo "  范围红线:只核对已接入 TF 的面;未接入资源/参数不报 gap。"
        return 0
    fi

    # 累加器
    local tmp; tmp="$(mktemp -d)"
    STEPS_FILE="$tmp/.steps.jsonl";        : > "$STEPS_FILE"
    FINDINGS_FILE="$tmp/.findings.jsonl";  : > "$FINDINGS_FILE"
    JQUEUE_FILE="$tmp/.judgment.jsonl";    : > "$JQUEUE_FILE"
    SUPPRESS_FILE="$tmp/.suppress.jsonl";  : > "$SUPPRESS_FILE"
    COVERAGE_FILE="$tmp/.coverage.jsonl";  : > "$COVERAGE_FILE"

    # 抑制表 / 容差表(冒号包裹便于 index 精确匹配;逗号连 tf>api 容差对)
    local supp tol
    supp=":$(cfg '.tier0_mech.suppress_params | join(":")'):"
    tol="$(cfg '.tier0_mech.type_tolerance | join(",")')"

    local start_epoch; start_epoch=$(date +%s)
    local started_at; started_at="$(date -u +%FT%TZ)"
    local doc_args_total=0 src_keys_total=0 res_ok=0

    local res name docf srcf
    for res in "${reslist[@]}"; do
        name="${res#alicloud_}"
        docf="$pdir/website/docs/r/${name}.html.markdown"
        srcf="$pdir/alicloud/resource_${res}.go"
        if [ ! -f "$docf" ] || [ ! -f "$srcf" ]; then
            jq -nc --arg res "$res" --argjson actions '[]' \
                --arg note "文档或源码文件缺失(doc:$([ -f "$docf" ] && echo ok || echo missing) src:$([ -f "$srcf" ] && echo ok || echo missing)),需人工确认资源名映射或是否为 data 源;不当 gap 报。" \
                '{resource:$res, reason:"missing_doc_or_src", api_actions:$actions, note:$note, detail:[]}' >> "$JQUEUE_FILE"
            continue
        fi
        _parse_doc_args "$docf" > "$tmp/doc.tsv"
        _parse_source_schema "$srcf" > "$tmp/src.tsv"
        local da; da=$(_count_lines "$tmp/doc.tsv")
        local sk; sk=$(_count_lines "$tmp/src.tsv")
        doc_args_total=$(( doc_args_total + da )); src_keys_total=$(( src_keys_total + sk )); res_ok=$(( res_ok + 1 ))

        # 本地 文档↔源码 机械 diff(五类 doc_gap,不受 mech 开关影响)
        local code attr sev summary
        while IFS=$'\t' read -r code attr sev summary; do
            [ -z "$code" ] && continue
            _emit_t0finding "$code" "$res" "$attr" "$summary" "$sev"
        done < <(_tier0_diff "$tmp/doc.tsv" "$tmp/src.tsv")

        # OpenAPI 侧:mech=on 机械三方 diff;off/degraded 走每资源一条 queue(现行为)
        if [ "$mech" = "on" ]; then
            _tier0_mech_resource "$res" "$srcf" "$tmp" "$supp" "$tol"
        else
            local actions_json reason
            actions_json="$(_source_api_actions "$srcf" | jq -R . | jq -s .)"
            reason=$([ "$mech" = "off" ] && echo "meta_off" || echo "meta_unavailable")
            jq -nc --arg res "$res" --arg reason "$reason" --argjson actions "$actions_json" \
                --arg note "tier-0 机械层 $mech:OpenAPI 侧交 skill 层人工双层查证(核对上列 action 的参数/枚举/行为与 TF 文档一致)。范围红线=只核已接入面;未接入不报 gap。" \
                '{resource:$res, reason:$reason, api_actions:$actions, note:$note, detail:[]}' >> "$JQUEUE_FILE"
        fi
    done

    # 聚合 step
    _emit_step doc_parse 0 0 "解析 $res_ok 个资源文档,共 $doc_args_total 个 Argument"
    _emit_step source_parse 0 0 "解析 $res_ok 个资源源码,共 $src_keys_total 个顶层 schema 键"
    local nf; nf=$(_count_lines "$FINDINGS_FILE")
    local api_nf; api_nf=$(jq -s '[.[]|select(.code|startswith("api_gap"))]|length' "$FINDINGS_FILE" 2>/dev/null); [ -z "$api_nf" ] && api_nf=0
    _emit_step diff 0 0 "三方一致性 diff(mech=$mech)→ $nf findings(内 api_gap_* $api_nf)"

    # --rotate:实际执行后记录本轮已扫时间戳(--dry 不记)
    [ "$rotate" -gt 0 ] && _rotate_mark "$rotate_sf" "${reslist[@]}"

    # verdict
    local audit; audit="$(probe_audit_dir)"; mkdir -p "$audit"
    local dur=$(( $(date +%s) - start_epoch ))
    local reslist_json; reslist_json="$(printf '%s\n' "${reslist[@]}" | jq -R . | jq -s .)"
    local verdict="$tmp/verdict.json"
    jq -n \
        --argjson schema 2 \
        --arg mode "tier0" \
        --arg mech "$mech" \
        --arg pv "$(cfg '.provider.version')" \
        --arg pdir "$pdir" \
        --arg started "$started_at" \
        --argjson dur "$dur" \
        --argjson resources "$reslist_json" \
        --slurpfile steps "$STEPS_FILE" \
        --slurpfile findings "$FINDINGS_FILE" \
        --slurpfile judgment "$JQUEUE_FILE" \
        --slurpfile suppressed "$SUPPRESS_FILE" \
        --slurpfile coverage "$COVERAGE_FILE" \
        --argjson stats "$(jq -nc --argjson r "$res_ok" --argjson da "$doc_args_total" --argjson sk "$src_keys_total" --argjson f "$nf" --argjson af "$api_nf" '{resources:$r, doc_args_total:$da, source_keys_total:$sk, findings:$f, api_findings:$af}')" \
        '{schema_version:$schema, mode:$mode, mech:$mech, provider_version:$pv, provider_dir:$pdir,
          started_at:$started, duration_s:$dur, resources_scanned:$resources,
          steps:$steps, findings:$findings, judgment_queue:$judgment,
          suppressed:$suppressed, coverage_notes:$coverage, stats:$stats}' \
        > "$verdict" 2>/dev/null

    # A3 verdict 同日覆盖修复:审计副本文件名带 HHMMSS,cp 目标和 echo verdict: 行必须一致(rc-gate.sh:105 依赖此行解析)。
    local day hms; day="$(date -u +%Y%m%d)"; hms="$(date -u +%H%M%S)"
    local audit_json="$audit/${day}-${hms}-tier0.json"
    local audit_md="$audit/${day}-${hms}-tier0.md"
    cp "$verdict" "$audit_json" 2>/dev/null
    _write_tier0_md "$verdict" "$audit_md"
    echo "verdict: $audit_json"
    echo "mech: $mech / findings: $nf (api_gap_* $api_nf) / resources: $res_ok / judgment_queue: $(_count_lines "$JQUEUE_FILE") / suppressed: $(_count_lines "$SUPPRESS_FILE")"

    # A2 ledger append(tier0 finalize)——本机加速索引,真源在 Aone。
    _ledger_append \
        '{ts:$ts, kind:"tier0", mech:$mech, resources:$r, findings:$f, verdict:$v}' \
        --arg mech "$mech" \
        --argjson r "$reslist_json" \
        --argjson f "$nf" \
        --arg v "$audit_json"

    rm -rf "$tmp"
    [ "$nf" -gt 0 ] && return 1 || return 0
}

_write_tier0_md() { # verdict.json out.md
    local v="$1" out="$2"
    {
        echo "# probe tier-0 三方一致性扫描(mech=$(jq -r '.mech // "off"' "$v"))"
        echo
        echo "- provider: $(jq -r '.provider_version' "$v") / dir: $(jq -r '.provider_dir' "$v")"
        echo "- started: $(jq -r '.started_at' "$v") / duration: $(jq -r '.duration_s' "$v")s"
        echo "- 扫描资源: $(jq -r '.resources_scanned | join(", ")' "$v")"
        echo
        echo "## findings ($(jq '.findings|length' "$v")) —— doc↔source $(jq '[.findings[]|select(.code|startswith("doc_gap"))]|length' "$v") / OpenAPI 机械 $(jq '[.findings[]|select(.code|startswith("api_gap"))]|length' "$v")"
        jq -r '.findings[]? | "- [\(.severity_hint)] `\(.code)` \(.resource).\(.attribute): \(.summary)"' "$v"
        echo
        echo "## judgment_queue (机械层拿不准,交 skill 层查证, $(jq '.judgment_queue|length' "$v"))"
        jq -r '.judgment_queue[]? | "- [\(.reason)] \(.resource): api_actions=[\(.api_actions|join(","))]\(if (.detail|length)>0 then " detail=["+(.detail|join(","))+"]" else "" end)\n  \(.note)"' "$v"
        echo
        echo "## suppressed ($(jq '.suppressed|length' "$v")) / coverage_notes ($(jq '.coverage_notes|length' "$v"))"
        jq -r '.suppressed[]? | "- suppress \(.resource).\(.param)→\(.api_param) (\(.rule))"' "$v"
        jq -r '.coverage_notes[]? | "- cover \(.resource).\(.attribute): \(.note)"' "$v"
        echo
        echo "来源:jarvis tf-customer-probe (tier-0, T0-mech)"
    } > "$out"
}

# ════════════════════════════════════════════════════════════════════
# run — tier-1 真实 apply 生命周期探测
# ════════════════════════════════════════════════════════════════════
PRUN_SID=""; PRUN_WD=""; PRUN_RUN_ID=""; PRUN_REGION=""
PRUN_STARTED_AT=""; PRUN_START_EPOCH=0; PRUN_TF_VERSION="unknown"
PRUN_APPLIED=false; PRUN_KEEP=0; PRUN_DESTROYED=null; PRUN_STATE_EMPTY=null
PRUN_CLEANED=0
STEP_RC=0; STEP_DURATION=0

_usage_run() { echo "用法: probe.sh run <scenario-id> [--region <r>] [--dry] [--keep]" >&2; }

# findings 计数码:destroy_fail/state_residue → 3;其它 findings → 1;无 → 0
_verdict_exit() {
    if jq -e -s 'any(.[]; .code=="destroy_fail" or .code=="state_residue")' "$FINDINGS_FILE" >/dev/null 2>&1; then
        echo 3; return
    fi
    if [ -s "$FINDINGS_FILE" ]; then echo 1; else echo 0; fi
}

# 清理:只要 apply 执行过就 destroy(除非 --keep);幂等(重入即返回)
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
        _emit_finding destroy_fail cleanup "destroy 失败,可能残留资源,需人工清理" "$PRUN_WD/destroy.log" S1
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
    local sid="" cli_region="" dry=0 keep=0
    while [ $# -gt 0 ]; do
        case "$1" in
            --region) cli_region="${2:-}"; shift 2 ;;
            --dry)  dry=1; shift ;;
            --keep) keep=1; shift ;;
            -*) echo "run: 未知参数 '$1'" >&2; _usage_run; return 2 ;;
            *) if [ -z "$sid" ]; then sid="$1"; shift; else echo "run: 多余参数 '$1'" >&2; return 2; fi ;;
        esac
    done
    [ -n "$sid" ] || { _usage_run; return 2; }

    # 两级布局:按 id 跨 product 目录检索;同 id 命中多个 product → 明确报错(退 2)
    local sdir yaml matches=()
    while IFS= read -r _m; do [ -n "$_m" ] && matches+=("$_m"); done < <(_find_scenario_dirs "$sid")
    if [ "${#matches[@]}" -eq 0 ]; then
        echo "run: 场景不存在 '$sid'(在 $(probe_playground_dir) 下未找到 <product>/$sid/scenario.yaml)" >&2
        return 2
    elif [ "${#matches[@]}" -gt 1 ]; then
        echo "run: 场景 id '$sid' 跨 product 目录重复,id 须全局唯一。命中:" >&2
        printf '  %s\n' "${matches[@]}" >&2
        return 2
    fi
    sdir="${matches[0]}"
    yaml="$sdir/scenario.yaml"

    local scn_region update_step import_check import_address import_id_output allow_prepaid apply_off
    local steps_csv expect_fail expect_err_contains provider_from drift_cli
    scn_region="$(_yaml_get "$yaml" region)"
    update_step="$(_yaml_get "$yaml" update_step)"
    import_check="$(_yaml_get "$yaml" import_check)"
    import_address="$(_yaml_get "$yaml" import_address)"
    import_id_output="$(_yaml_get "$yaml" import_id_output)"
    allow_prepaid="$(_yaml_get "$yaml" allow_prepaid)"
    apply_off=0; _apply_disabled "$yaml" && apply_off=1
    # B2 新键(全部可选;缺省=向后兼容 update_step 语义)
    steps_csv="$(_yaml_get "$yaml" steps)"
    expect_fail="$(_yaml_get "$yaml" expect_fail)"
    expect_err_contains="$(_yaml_get "$yaml" expect_error_contains)"
    provider_from="$(_yaml_get "$yaml" provider_version_from)"
    drift_cli="$(_yaml_get "$yaml" drift_cli)"
    # 向后兼容:update_step: true 等价 steps: step2
    if [ -z "$steps_csv" ] && [ "$update_step" = "true" ]; then steps_csv="step2"; fi

    local region; region="$(_resolve_region "$cli_region" "$scn_region")"
    local tier1_enabled; tier1_enabled="$(cfg '.tiers.tier1.enabled')"
    local have_creds=0; _have_creds && have_creds=1

    # ── --dry:只打印计划,不碰 terraform ──
    if [ "$dry" -eq 1 ]; then
        echo "probe plan: scenario=$sid persona=$(_yaml_get "$yaml" persona)"
        echo "  region 解析: --region='${cli_region}' scenario='${scn_region}' focus=$(cfg '.regions.focus') env='${ALICLOUD_REGION:-}' → $region"
        echo "  tier1.enabled=$tier1_enabled  prepaid_guard=$(cfg '.tiers.tier1.prepaid_guard')  allow_prepaid=${allow_prepaid:-false}"
        echo "  steps:"
        echo "    - init"
        echo "    - validate"
        if [ "$have_creds" -eq 1 ]; then echo "    - plan -out=tf.plan"; else echo "    - plan: 将跳过(未设置 ALICLOUD 凭证 → env_issue no_creds)"; fi
        if [ "$apply_off" -eq 1 ]; then
            echo "    - (场景 apply:false → 止步 plan,plan-only 封顶,不 apply;env_issue apply_disabled_by_scenario)"
        elif [ "$tier1_enabled" != "true" ]; then
            echo "    - (tier1.enabled=false → plan-only 封顶,不 apply;env_issue tier1_disabled_plan_only)"
        else
            if [ "$allow_prepaid" = "true" ] || [ "$(cfg '.tiers.tier1.prepaid_guard')" != "true" ]; then
                echo "    - prepaid-gate: 跳过(allow_prepaid 或 prepaid_guard=false)"
            else
                echo "    - prepaid-gate: plan 出现 PrePaid/Subscription 计费类型 → env_issue prepaid_block,不 apply(销毁性守门)"
            fi
            echo "    - apply -auto-approve"
            echo "    - plan -detailed-exitcode(退2 → perpetual_diff)"
            # B2 步骤:steps CSV 泛化 update_step,step<N>_expect 逐步声明期望
            if [ -n "$steps_csv" ]; then
                local _s _e
                IFS=',' read -ra _steps_arr <<< "$steps_csv"
                for _s in "${_steps_arr[@]}"; do
                    _s="$(echo "$_s" | tr -d '[:space:]')"; [ -z "$_s" ] && continue
                    _e="$(_yaml_get "$yaml" "${_s}_expect")"
                    [ -z "$_e" ] && _e="changed"
                    echo "    - $_s 覆盖 apply + re-plan(expect=$_e:changed=perpetual_diff/no_changes=refactor_replace/fail=expected_fail_missed)"
                done
            fi
            # B2 expect_fail:三态判定纯函数 _expect_fail_verdict
            [ -n "$expect_fail" ] && echo "    - expect_fail=$expect_fail contains='${expect_err_contains}' → 阶段对齐=expected;早=常规;晚=late_validation(S3);未失败=expected_fail_missed(S2)"
            # B2 provider_version_from:upgrader 场景
            [ -n "$provider_from" ] && echo "    - upgrader: pin $provider_from → init → apply → 改回 config pin → init -upgrade → plan → diff → upgrade_diff(S2;delete+create→S1)"
            # B2 drift_cli:drifter 场景(默认 drift_enabled=false → env_issue drift_disabled)
            [ -n "$drift_cli" ] && echo "    - drift: 五重护栏 tokenize + 白名单 + 占位符 + 凭证映射 + drift_enabled 门;post plan 无 diff → drift_undetected(S2)"
            [ "$import_check" = "true" ] && echo "    - state rm $import_address → import(id 取自 output $import_id_output) → plan(退2 → import_diff)"
            echo "    - destroy -auto-approve + state-empty 校验(--keep 跳过;destroy 失败/残留 → finding + 退3)"
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
    PRUN_REGION="$region"
    PRUN_KEEP="$keep"
    PRUN_APPLIED=false; PRUN_DESTROYED=null; PRUN_STATE_EMPTY=null; PRUN_CLEANED=0
    mkdir -p "$PRUN_WD"

    STEPS_FILE="$PRUN_WD/.steps.jsonl";      : > "$STEPS_FILE"
    FINDINGS_FILE="$PRUN_WD/.findings.jsonl"; : > "$FINDINGS_FILE"
    ENV_FILE="$PRUN_WD/.env.jsonl";           : > "$ENV_FILE"

    cp "$sdir"/*.tf "$PRUN_WD"/ 2>/dev/null

    # terraform 自动化环境 + region 注入(场景 tf 不写显式 region,沿用 ALICLOUD_REGION)
    export TF_IN_AUTOMATION=1 TF_INPUT=0 TF_CLI_ARGS=-no-color
    [ -n "$region" ] && export ALICLOUD_REGION="$region"
    TF_PLUGIN_CACHE_DIR="$(probe_workdir_base)/.plugin-cache"; export TF_PLUGIN_CACHE_DIR
    mkdir -p "$TF_PLUGIN_CACHE_DIR"

    PRUN_STARTED_AT="$(date -u +%FT%TZ)"; PRUN_START_EPOCH=$(date +%s)
    PRUN_TF_VERSION="$(terraform version 2>/dev/null | head -1)"

    trap _probe_on_exit EXIT

    # init → validate → (有凭证)plan
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
    if [ "$STEP_RC" -ne 0 ]; then
        # expect_fail 优先分流:validate 阶段失败 vs 声明阶段的三态判定
        if [ -n "$expect_fail" ]; then
            _emit_expect_finding "$expect_fail" validate "$expect_err_contains" "$PRUN_WD/validate.log" \
                validate_fail S3 "官方文档示例组合 validate 不通过"
        else
            _emit_finding validate_fail validate "官方文档示例组合 validate 不通过" "$PRUN_WD/validate.log" S3
        fi
    fi

    if [ "$have_creds" -eq 0 ]; then
        _emit_env no_creds "未设置 ALICLOUD 凭证,跳过 plan/apply"
        _finalize_verdict; return "$(_verdict_exit)"
    fi

    _run_step plan "$PRUN_WD/plan.log" -- terraform plan -out=tf.plan -var "run_id=$PRUN_RUN_ID"
    if [ "$STEP_RC" -ne 0 ]; then
        if _is_auth_error "$PRUN_WD/plan.log"; then
            _emit_env auth_error "plan 报鉴权类错误(归 env,不算 provider bug)"
        elif _has_panic "$PRUN_WD/plan.log"; then
            _emit_finding plan_crash plan "terraform plan 触发 provider panic" "$PRUN_WD/plan.log" S1
        elif [ -n "$expect_fail" ]; then
            _emit_expect_finding "$expect_fail" plan "$expect_err_contains" "$PRUN_WD/plan.log" \
                plan_fail S2 "合法配置 plan 失败"
        else
            _emit_finding plan_fail plan "合法配置 plan 失败" "$PRUN_WD/plan.log" S2
        fi
        _finalize_verdict; return "$(_verdict_exit)"
    fi
    # A2 LRU 索引:仅在真实执行推进到 plan 完成及以后才更新,避免被阻断场景(no_creds/prepaid_block/tier1_disabled)永远失去 LRU 优先权。
    _t1_last_run_update "$PRUN_SID"

    # 场景级 apply 门:scenario.yaml apply:false(生成器命中成本安全门写入)→ 止步 plan,不真实创建
    if [ "$apply_off" -eq 1 ]; then
        _emit_env apply_disabled_by_scenario "场景 scenario.yaml 声明 apply:false(成本安全门/付费大件资源),封顶 plan-only,不 apply"
        _finalize_verdict; return "$(_verdict_exit)"
    fi

    # tier1.enabled=false → plan-only 封顶(语义已变:不再叫降级 tier-0)
    if [ "$tier1_enabled" != "true" ]; then
        _emit_env tier1_disabled_plan_only "配置 tiers.tier1.enabled=false,封顶 plan-only,不 apply"
        _finalize_verdict; return "$(_verdict_exit)"
    fi

    # prepaid 销毁性守门(替代原免费 allowlist 成本门)
    ( cd "$PRUN_WD" && terraform show -json tf.plan > "$PRUN_WD/plan.json" 2>/dev/null )
    if _prepaid_should_block "$PRUN_WD/plan.json" "$allow_prepaid"; then
        _emit_env prepaid_block "plan 出现 PrePaid/Subscription 计费类型资源(多数无法 API 销毁,破坏零残留纪律),终止不 apply;场景可声明 allow_prepaid:true 豁免"
        _finalize_verdict; return 2
    fi

    _run_step apply "$PRUN_WD/apply.log" -- terraform apply -auto-approve -var "run_id=$PRUN_RUN_ID"
    PRUN_APPLIED=true
    if [ "$STEP_RC" -ne 0 ]; then
        if _is_auth_error "$PRUN_WD/apply.log"; then
            _emit_env auth_error "apply 报鉴权类错误"
        elif [ -n "$expect_fail" ]; then
            _emit_expect_finding "$expect_fail" apply "$expect_err_contains" "$PRUN_WD/apply.log" \
                apply_fail S2 "合法配置 apply 失败"
        else
            _emit_finding apply_fail apply "合法配置 apply 失败" "$PRUN_WD/apply.log" S2
        fi
        _probe_cleanup; _finalize_verdict; return "$(_verdict_exit)"
    fi

    # 幂等:apply 后立刻 plan 应无 diff(退2 = 永久 diff)
    _run_step plan2 "$PRUN_WD/plan2.log" -- terraform plan -detailed-exitcode -out=tf.plan2 -var "run_id=$PRUN_RUN_ID"
    [ "$STEP_RC" -eq 2 ] && _emit_finding perpetual_diff plan2 "apply 后立即 plan 仍有 diff(幂等性破坏/永久 diff)" "$PRUN_WD/plan2.log" S2
    _detect_unexpected_replace "$PRUN_WD/plan2.log" plan2

    # B2 steps CSV(泛化 update_step;update_step:true 已在前面归一为 steps=step2)
    if [ -n "$steps_csv" ]; then
        local _step _expect _plan_rc _step_apply_rc
        IFS=',' read -ra _steps_arr <<< "$steps_csv"
        for _step in "${_steps_arr[@]}"; do
            _step="$(echo "$_step" | tr -d '[:space:]')"; [ -z "$_step" ] && continue
            [ -d "$sdir/$_step" ] || { _emit_env "step_dir_missing" "steps 声明 $_step 但 $sdir/$_step 不存在,跳过"; continue; }
            _expect="$(_yaml_get "$yaml" "${_step}_expect")"; [ -z "$_expect" ] && _expect="changed"
            cp "$sdir/$_step"/*.tf "$PRUN_WD"/ 2>/dev/null
            _run_step "apply_${_step}" "$PRUN_WD/apply_${_step}.log" -- terraform apply -auto-approve -var "run_id=$PRUN_RUN_ID"
            _step_apply_rc="$STEP_RC"
            if [ "$_step_apply_rc" -ne 0 ]; then
                # apply 失败:expect=fail 视为符合;其它当 apply_fail
                if [ "$_expect" = "fail" ]; then
                    _emit_env "step_fail_expected" "步骤 $_step apply 按预期失败(expect=fail)"
                else
                    _emit_finding apply_fail "apply_${_step}" "步骤 $_step apply 失败" "$PRUN_WD/apply_${_step}.log" S2
                fi
                continue
            fi
            # apply 成功:按 expect 分流
            _run_step "plan_${_step}" "$PRUN_WD/plan_${_step}.log" -- terraform plan -detailed-exitcode -out="tf.plan_${_step}" -var "run_id=$PRUN_RUN_ID"
            _plan_rc="$STEP_RC"
            case "$_expect" in
                changed)
                    [ "$_plan_rc" -eq 2 ] && _emit_finding perpetual_diff "plan_${_step}" "步骤 $_step 更新后 plan 仍有 diff(更新不生效)" "$PRUN_WD/plan_${_step}.log" S2
                    ;;
                no_changes)
                    if [ "$_plan_rc" -eq 2 ]; then
                        # no_changes 却出现 diff → refactor_replace(delete+create → S1;否则 S2);复用 _detect_unexpected_replace 判断
                        local _pj="$PRUN_WD/plan_${_step}.json" _sev="S2"
                        ( cd "$PRUN_WD" && terraform show -json "tf.plan_${_step}" > "$_pj" 2>/dev/null )
                        if [ -s "$_pj" ] && jq -e 'any(.resource_changes[]?; (.change.actions|index("delete")) and (.change.actions|index("create")))' "$_pj" >/dev/null 2>&1; then
                            _sev="S1"
                        fi
                        _emit_finding refactor_replace "plan_${_step}" "步骤 $_step 声明 no_changes 但 plan 出现 diff(refactor 触发替换重建)" "$PRUN_WD/plan_${_step}.log" "$_sev"
                    fi
                    ;;
                fail)
                    # apply 已经成功了(fail 期望未达)→ expected_fail_missed
                    _emit_finding expected_fail_missed "apply_${_step}" "步骤 $_step 声明 expect=fail 但 apply 成功" "$PRUN_WD/apply_${_step}.log" S2
                    ;;
            esac
        done
    fi

    # B2 provider_version_from:upgrader dance —— sed 改 workdir 副本 pin → init → apply → 改回 → init -upgrade → plan;非空 diff → upgrade_diff。
    #   注:进入此分支时 apply 已跑过(target=current pin)。upgrader 需要的是"旧版本先立起来 → 升级 → 观察 diff"。
    #   本实现按契约:①现场 tf 备份 → sed pin → init → apply → sed 回 current pin → init -upgrade → plan;②非空 diff→upgrade_diff。
    #   TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE=1 防双全量下载 lock 冲突。
    if [ -n "$provider_from" ]; then
        local _cur_pin _mtf _lockbak
        _cur_pin="$(cfg '.provider.version')"
        _mtf="$PRUN_WD/main.tf"
        export TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE=1
        if [ -f "$_mtf" ] && [ -n "$_cur_pin" ]; then
            # 备份 lock,sed 改 pin
            [ -f "$PRUN_WD/.terraform.lock.hcl" ] && cp "$PRUN_WD/.terraform.lock.hcl" "$PRUN_WD/.terraform.lock.hcl.bak" 2>/dev/null
            sed -i.bak "s/version = \"${_cur_pin//./\\.}\"/version = \"${provider_from//./\\.}\"/" "$_mtf" 2>/dev/null
            _run_step upgrade_init_old "$PRUN_WD/upgrade_init_old.log" -- terraform init -input=false -upgrade
            if [ "$STEP_RC" -ne 0 ]; then
                _emit_env upgrade_init_fail "upgrader init(pin=$provider_from) 失败,见 $PRUN_WD/upgrade_init_old.log"
            else
                _run_step upgrade_apply_old "$PRUN_WD/upgrade_apply_old.log" -- terraform apply -auto-approve -var "run_id=$PRUN_RUN_ID"
                if [ "$STEP_RC" -ne 0 ]; then
                    _emit_env upgrade_apply_fail "upgrader 旧版本 apply 失败(pin=$provider_from)"
                else
                    # 改回 current pin,init -upgrade,plan
                    sed -i.bak "s/version = \"${provider_from//./\\.}\"/version = \"${_cur_pin//./\\.}\"/" "$_mtf" 2>/dev/null
                    _run_step upgrade_init_new "$PRUN_WD/upgrade_init_new.log" -- terraform init -input=false -upgrade
                    if [ "$STEP_RC" -ne 0 ]; then
                        _emit_env upgrade_init_new_fail "upgrader init(升级至 $_cur_pin)失败"
                    else
                        _run_step upgrade_plan "$PRUN_WD/upgrade_plan.log" -- terraform plan -detailed-exitcode -out=tf.plan_upgrade -var "run_id=$PRUN_RUN_ID"
                        if [ "$STEP_RC" -eq 2 ]; then
                            # 非空 diff → upgrade_diff;delete+create → S1
                            local _upj="$PRUN_WD/plan_upgrade.json" _usev="S2"
                            ( cd "$PRUN_WD" && terraform show -json tf.plan_upgrade > "$_upj" 2>/dev/null )
                            if [ -s "$_upj" ] && jq -e 'any(.resource_changes[]?; (.change.actions|index("delete")) and (.change.actions|index("create")))' "$_upj" >/dev/null 2>&1; then
                                _usev="S1"
                            fi
                            _emit_finding upgrade_diff upgrade_plan "provider $provider_from → $_cur_pin 升级后 plan 出现 diff(state 兼容性破坏)" "$PRUN_WD/upgrade_plan.log" "$_usev"
                        fi
                    fi
                fi
            fi
            # 清理 sed .bak(不影响 destroy 阶段)
            rm -f "$_mtf.bak" 2>/dev/null
        fi
    fi

    # B2 drift_cli:drifter 五重护栏(见函数注释);drift_enabled 门 + action 白名单 + tokenize + 占位符 + 凭证钉死
    if [ -n "$drift_cli" ]; then
        _probe_drift_dance "$drift_cli"
    fi

    # B2 expect_fail 收尾判定:走过完整生命周期(validate/plan/apply 均 rc=0)都没失败 → expected_fail_missed(S2)
    #   注:各失败早退分支中,若 expect_fail 命中,会用 _emit_expect_finding 分流(expected/late/early),不到这里。
    if [ -n "$expect_fail" ]; then
        _emit_finding expected_fail_missed apply "场景声明 expect_fail=$expect_fail 但全程未失败(预期错误未触发)" "$PRUN_WD/apply.log" S2
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

# B2 _probe_drift_dance <cli_string> — drifter 五重护栏 orchestrator。
#   ①无 shell 直执:_drift_tokenize 切 argv + argv[0] 字面 aliyun + 元字符黑名单
#   ②占位符:_drift_expand_placeholders 仅认 {{output.NAME}}/{{region}} 且值过白名单正则
#   ③凭证:显式导出 ALIBABA_CLOUD_ACCESS_KEY_ID/SECRET/REGION_ID(映射 ALICLOUD_*);映射不成立=env_issue drift_env_missing
#   ④action 白名单:.tiers.tier1.drift_action_allow 按 <product>:<Action> 匹配 argv[1]:argv[2]
#   ⑤drift_enabled 默认 false:关闭时 env_issue drift_disabled 不跑
#   post-drift plan 无 diff → finding drift_undetected(默认 S2)
_probe_drift_dance() {
    local cli="$1"
    # ⑤ drift_enabled 门
    if [ "$(cfg_or '.tiers.tier1.drift_enabled' 'false')" != "true" ]; then
        _emit_env drift_disabled "config .tiers.tier1.drift_enabled=false,drift 场景不跑(无人值守日轮默认关;转正开关走 MR)"
        return 0
    fi
    # ③ 凭证钉死:ALICLOUD_ACCESS_KEY/SECRET_KEY → ALIBABA_CLOUD_ACCESS_KEY_ID/SECRET;缺则拒跑
    if [ -z "${ALICLOUD_ACCESS_KEY:-}" ] || [ -z "${ALICLOUD_SECRET_KEY:-}" ] || [ -z "${PRUN_REGION:-}" ]; then
        _emit_env drift_env_missing "ALICLOUD_ACCESS_KEY/SECRET_KEY 或 region 未设,drift 拒跑(runner 显式映射不成立)"
        return 0
    fi
    export ALIBABA_CLOUD_ACCESS_KEY_ID="$ALICLOUD_ACCESS_KEY"
    export ALIBABA_CLOUD_ACCESS_KEY_SECRET="$ALICLOUD_SECRET_KEY"
    export ALIBABA_CLOUD_REGION_ID="$PRUN_REGION"
    # ① tokenize
    local -a argv=()
    while IFS= read -r _tok; do argv+=("$_tok"); done < <(_drift_tokenize "$cli" 2>>"$PRUN_WD/drift.log")
    if [ "${#argv[@]}" -lt 3 ]; then
        _emit_env drift_cli_rejected "drift_cli tokenize 失败(见 drift.log)"
        return 0
    fi
    # ④ 白名单
    if ! _drift_action_allowed "${argv[1]}" "${argv[2]}"; then
        _emit_env drift_action_denied "drift_cli action ${argv[1]}:${argv[2]} 不在 tiers.tier1.drift_action_allow 白名单"
        return 0
    fi
    # ② 占位符展开(每个 token 独立展开)
    local -a expanded=()
    local i tok exp
    for i in "${!argv[@]}"; do
        tok="${argv[$i]}"
        if [[ "$tok" == *"{{"*"}}"* ]]; then
            exp="$(_drift_expand_placeholders "$tok" "$PRUN_WD" "$PRUN_REGION" 2>>"$PRUN_WD/drift.log")"
            if [ -z "$exp" ]; then
                _emit_env drift_cli_rejected "drift_cli 占位符展开失败(见 drift.log)"
                return 0
            fi
            expanded+=("$exp")
        else
            expanded+=("$tok")
        fi
    done
    # exec:数组直执,绝不过 sh -c/eval
    { echo "### drift exec :: $(date -u +%FT%TZ)"; printf '$ '; printf '%q ' "${expanded[@]}"; echo; } >> "$PRUN_WD/drift.log"
    "${expanded[@]}" >> "$PRUN_WD/drift.log" 2>&1
    local _drift_rc=$?
    if [ "$_drift_rc" -ne 0 ]; then
        _emit_env drift_exec_fail "aliyun CLI drift 命令执行失败(rc=$_drift_rc),见 $PRUN_WD/drift.log"
        return 0
    fi
    # post-drift plan:无 diff → drift_undetected
    _run_step plan_drift "$PRUN_WD/plan_drift.log" -- terraform plan -detailed-exitcode -var "run_id=$PRUN_RUN_ID"
    if [ "$STEP_RC" -ne 2 ]; then
        _emit_finding drift_undetected plan_drift "带外改动后 plan 未检出 diff(drift detection 失效)" "$PRUN_WD/plan_drift.log" S2
    fi
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
        --arg mode "tier1" \
        --arg sid "$PRUN_SID" \
        --arg pv "$(cfg '.provider.version')" \
        --arg tv "${PRUN_TF_VERSION:-unknown}" \
        --arg region "$PRUN_REGION" \
        --arg started "$PRUN_STARTED_AT" \
        --argjson dur "$dur" \
        --slurpfile steps "$STEPS_FILE" \
        --slurpfile findings "$FINDINGS_FILE" \
        --slurpfile envs "$ENV_FILE" \
        --argjson applied "$PRUN_APPLIED" \
        --argjson destroyed "$PRUN_DESTROYED" \
        --argjson state_empty "$PRUN_STATE_EMPTY" \
        '{schema_version:$schema, mode:$mode, scenario_id:$sid,
          provider_version:$pv, terraform_version:$tv, region:$region,
          started_at:$started, duration_s:$dur,
          steps:$steps, findings:$findings, env_issues:$envs,
          cleanup:{applied:$applied, destroyed:$destroyed, state_empty:$state_empty}}' \
        > "$verdict" 2>/dev/null

    # A3 verdict 同日覆盖修复:审计副本文件名带 HHMMSS。
    local day hms; day="$(date -u +%Y%m%d)"; hms="$(date -u +%H%M%S)"
    local audit_json="$audit/${day}-${hms}-${PRUN_SID}.json"
    local audit_md="$audit/${day}-${hms}-${PRUN_SID}.md"
    cp "$verdict" "$audit_json" 2>/dev/null
    _write_summary_md "$verdict" "$audit_md"
    echo "verdict: $verdict"
    # A3 echo verdict: 一行必须与 cp 目标一致(rc-gate.sh:105 用它取 verdict 路径解析)。tier1 保留原双行契约。
    echo "verdict: $audit_json"
    echo "audit:   $audit_json"

    # A2 ledger append(tier1 finalize):findings 计数 + verdict 路径。
    local nf; nf="$(jq '.findings|length' "$verdict" 2>/dev/null)"; [ -z "$nf" ] && nf=0
    _ledger_append \
        '{ts:$ts, kind:"tier1", scenario:$sid, findings:$f, verdict:$v}' \
        --arg sid "$PRUN_SID" \
        --argjson f "$nf" \
        --arg v "$audit_json"
}

_write_summary_md() { # verdict.json out.md
    local v="$1" out="$2"
    {
        echo "# probe verdict: $(jq -r '.scenario_id' "$v") (tier-1)"
        echo
        echo "- provider: $(jq -r '.provider_version' "$v") / terraform: $(jq -r '.terraform_version' "$v")"
        echo "- region: $(jq -r '.region' "$v")"
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

# ════════════════════════════════════════════════════════════════════
# B2 纯函数:expect_fail 三态判定 + drift_cli 五重护栏
# ════════════════════════════════════════════════════════════════════
# 阶段序:validate=1 < plan=2 < apply=3。用于 expect_fail 早晚比较。
_stage_ord() { case "$1" in validate) echo 1;; plan) echo 2;; apply) echo 3;; *) echo 0;; esac; }

# _expect_fail_verdict <declared_stage> <failed_stage> <err_contains_ok(1|0)>
#   源可单测(先例 _prepaid_should_block)。stdout=码 tag:
#     expected                     — 在声明阶段失败,且 error_contains 匹配(或无 error_contains)
#     expected_but_error_mismatch  — 阶段对上但 error_contains 未命中 → S3(证据不符合;非预期错因)
#     early_failure_fallthrough    — 早于声明阶段失败 → 走现行 finding 分流(runner 不当 expect finding 处理)
#     late_validation              — 晚于声明阶段才失败 → finding late_validation S3
#     expected_fail_missed         — 全程未失败(failed_stage 为空)→ finding expected_fail_missed S2
#   declared_stage/failed_stage ∈ {validate,plan,apply,""};failed_stage="" 表示从未失败。
_expect_fail_verdict() {
    local declared="$1" failed="$2" err_ok="${3:-1}"
    if [ -z "$failed" ]; then
        echo "expected_fail_missed"; return
    fi
    local dord ford
    dord=$(_stage_ord "$declared"); ford=$(_stage_ord "$failed")
    if [ "$dord" -eq 0 ] || [ "$ford" -eq 0 ]; then
        # 无效阶段名一律按现行分流(不当 expect 特化处理)
        echo "early_failure_fallthrough"; return
    fi
    if [ "$ford" -lt "$dord" ]; then
        echo "early_failure_fallthrough"; return
    fi
    if [ "$ford" -gt "$dord" ]; then
        echo "late_validation"; return
    fi
    # 相等
    if [ "$err_ok" = "1" ]; then echo "expected"; else echo "expected_but_error_mismatch"; fi
}

# _emit_expect_finding <declared> <actual_failed> <err_contains> <errlog> <fallback_code> <fallback_sev> <fallback_summary>
#   B2 expect_fail 分流:声明阶段=actual → expected(不报) 或 expected_but_error_mismatch(S3);
#   声明阶段>actual → early_failure_fallthrough(用 fallback 常规 finding);
#   声明阶段<actual → late_validation(S3)。
_emit_expect_finding() {
    local declared="$1" actual="$2" err_contains="$3" errlog="$4"
    local fb_code="$5" fb_sev="$6" fb_sum="$7"
    local err_ok=1
    if [ -n "$err_contains" ]; then
        grep -qi -- "$err_contains" "$errlog" 2>/dev/null && err_ok=1 || err_ok=0
    fi
    local verdict; verdict="$(_expect_fail_verdict "$declared" "$actual" "$err_ok")"
    case "$verdict" in
        expected)
            _emit_env expected_failure "阶段 $actual 按预期失败(expect_fail=$declared;error_contains 命中)"
            ;;
        expected_but_error_mismatch)
            _emit_finding expected_but_error_mismatch "$actual" \
                "阶段 $actual 按预期失败但 error 未含 '$err_contains'(错因不符合)" "$errlog" S3
            ;;
        late_validation)
            _emit_finding late_validation "$actual" \
                "声明期望 $declared 失败,但直到 $actual 才失败(前置校验太宽)" "$errlog" S3
            ;;
        early_failure_fallthrough|*)
            _emit_finding "$fb_code" "$actual" "$fb_sum" "$errlog" "$fb_sev"
            ;;
    esac
}

# _drift_tokenize <cli_string> — 按空白切 argv,元字符黑名单 + argv[0] 必须字面 aliyun。
#   护栏#1(无 shell 直执):任一 token 含 `; & | $ ( ) < > \` 反引号/换行/反斜杠即拒绝;argv[0] != aliyun 拒绝。
#   stdout:每行一个 token(过白名单);non-zero 退码表示被拒(stderr 有原因)。
_drift_tokenize() {
    local s="$1" tok
    # 换行/反斜杠一律不允许
    case "$s" in
        *$'\n'*|*$'\r'*) echo "drift_cli: 含换行/回车,拒绝" >&2; return 1 ;;
        *\\*) echo "drift_cli: 含反斜杠,拒绝" >&2; return 1 ;;
    esac
    # 用 read -ra 按 IFS 空白切,不过 shell
    local -a argv=()
    # shellcheck disable=SC2162
    read -ra argv <<< "$s"
    [ "${#argv[@]}" -ge 3 ] || { echo "drift_cli: 至少需 aliyun <product> <Action> 三段" >&2; return 1; }
    [ "${argv[0]}" = "aliyun" ] || { echo "drift_cli: argv[0] 必须字面 aliyun(拒绝 ${argv[0]})" >&2; return 1; }
    for tok in "${argv[@]}"; do
        case "$tok" in
            *';'*|*'&'*|*'|'*|*'$'*|*'('*|*')'*|*'<'*|*'>'*|*'`'*)
                echo "drift_cli: token 含元字符(拒绝:'$tok')" >&2; return 1 ;;
        esac
        printf '%s\n' "$tok"
    done
}

# _drift_action_allowed <product> <action> — 0=白名单命中,1=拒绝。config `.tiers.tier1.drift_action_allow` <product>:<Action>。
#   config 分裂防御:缺键回默认 vpc:TagResources 三件套。
_drift_action_allowed() {
    local product="$1" action="$2"
    local allow key
    # cfg 会输出 null;这里手动取值,遇 null 走默认。
    allow="$(jq -r '.tiers.tier1.drift_action_allow // ["vpc:TagResources","vpc:UnTagResources","vpc:ModifyVpcAttribute"] | .[]' "$(probe_config)" 2>/dev/null)"
    key="${product}:${action}"
    while IFS= read -r a; do
        [ -n "$a" ] || continue
        [ "$a" = "$key" ] && return 0
    done <<< "$allow"
    return 1
}

# _drift_expand_placeholders <token> <workdir> <region>
#   护栏#2(占位符受限注入):仅 `{{output.<name>}}`(terraform output -raw)与 `{{region}}`;
#   注入值须过 `^[A-Za-z0-9._:/-]+$`;不满足即 stdout 打印占位符原文并退非零。
#   stdout=展开后的 token;returns:0=展开 ok(可能未含占位符,原样返回)。
_drift_expand_placeholders() {
    local tok="$1" wd="$2" region="$3"
    local out="$tok"
    # {{region}}
    if [[ "$out" == *"{{region}}"* ]]; then
        # region 必过白名单(实际是 ALIBABA_CLOUD_REGION_ID 值,已由 runner 校验非空)
        [[ "$region" =~ ^[A-Za-z0-9._:/-]+$ ]] || { echo "drift_cli: region 值不合规($region)" >&2; return 1; }
        out="${out//\{\{region\}\}/$region}"
    fi
    # {{output.NAME}}
    while [[ "$out" =~ \{\{output\.([A-Za-z0-9_]+)\}\} ]]; do
        local name="${BASH_REMATCH[1]}" val
        val="$(cd "$wd" && terraform output -raw "$name" 2>/dev/null)"
        [ -n "$val" ] || { echo "drift_cli: output '$name' 空,占位符展开失败" >&2; return 1; }
        [[ "$val" =~ ^[A-Za-z0-9._:/-]+$ ]] || { echo "drift_cli: output '$name' 值不过白名单(${val:0:80}…)" >&2; return 1; }
        # 全量替换(简单场景够用;若需多占位符也一并处理)
        out="${out//\{\{output.$name\}\}/$val}"
    done
    printf '%s' "$out"
}

# ════════════════════════════════════════════════════════════════════
# A1 archive [--dry] — 幂等归档(draft/verdict retention/workdir gc/plugin-cache 报告/待办清单)
# ════════════════════════════════════════════════════════════════════
# 判老:优先 JSON started_at,缺则文件 mtime。返回 epoch(0=解析失败,视为未老,不动)。
_archive_verdict_epoch() {
    local f="$1" d v
    d="$(jq -r '.started_at // empty' "$f" 2>/dev/null)"
    if [ -n "$d" ]; then
        v="$(date -u -j -f "%Y-%m-%dT%H:%M:%SZ" "$d" +%s 2>/dev/null || date -u -d "$d" +%s 2>/dev/null || echo 0)"
        [ -n "$v" ] && [ "$v" != "0" ] && echo "$v" && return
    fi
    stat -f %m "$f" 2>/dev/null || stat -c %Y "$f" 2>/dev/null || echo 0
}

# _draft_status <path> — 抽 frontmatter status 值(小写),缺失/坏文件回空。
_draft_status() {
    local f="$1"
    awk 'NR==1 && $0=="---" {infm=1; next}
         infm && $0=="---" {exit}
         infm && /^[[:space:]]*status[[:space:]]*:/ { sub(/^[[:space:]]*status[[:space:]]*:[[:space:]]*/,""); gsub(/"|'"'"'/,""); sub(/[[:space:]]+$/,""); print tolower($0); exit }' "$f" 2>/dev/null
}

# _archive_drafts <dry> — draft 归档;stdout 打印摘要;stderr 打印每项动作。写入全局 __ARCH_MOVED_DRAFTS(计数)。
__ARCH_MOVED_DRAFTS=0
__ARCH_PENDING_DRAFTS=""
_archive_drafts() {
    local dry="$1"
    local root drafts arc dest st f base
    root="$(probe_root)"
    drafts="$root/$(cfg_or '.paths.drafts' 'escalation/probe-drafts')"
    arc="$root/$(cfg_or '.paths.drafts_archived' 'escalation/probe-drafts/archived')"
    [ -d "$drafts" ] || { echo "  drafts: 无 $drafts 目录,跳过"; return 0; }
    [ "$dry" = "0" ] && mkdir -p "$arc" 2>/dev/null
    __ARCH_MOVED_DRAFTS=0
    __ARCH_PENDING_DRAFTS=""
    shopt -s nullglob
    for f in "$drafts"/*.md; do
        base="$(basename "$f")"
        [ "$base" = ".gitkeep" ] && continue
        st="$(_draft_status "$f")"
        case "$st" in
            filed|rejected|rejected-*)
                __ARCH_MOVED_DRAFTS=$((__ARCH_MOVED_DRAFTS+1))
                if [ "$dry" = "1" ]; then
                    echo "  DRY draft mv: $base → archived/ (status=$st)" >&2
                else
                    mv -f "$f" "$arc/$base" 2>/dev/null && echo "  draft mv: $base → archived/ (status=$st)" >&2
                fi
                ;;
            pending-review|pending)
                __ARCH_PENDING_DRAFTS="$__ARCH_PENDING_DRAFTS $base"
                ;;
            *)
                # 空/未知 status 留原地,进 pending 清单
                __ARCH_PENDING_DRAFTS="$__ARCH_PENDING_DRAFTS $base"
                ;;
        esac
    done
    echo "  drafts: moved=$__ARCH_MOVED_DRAFTS pending=$(printf '%s' "$__ARCH_PENDING_DRAFTS" | wc -w | tr -d ' ')"
}

# _archive_verdicts <dry> — verdict retention;写 __ARCH_MOVED_VERDICTS。
__ARCH_MOVED_VERDICTS=0
_archive_verdicts() {
    local dry="$1"
    local audit ret_days now cutoff f base ep ym destdir
    audit="$(probe_audit_dir)"
    [ -d "$audit" ] || { echo "  verdicts: 无 $audit 目录,跳过"; return 0; }
    ret_days="$(cfg_or '.limits.audit_retention_days' '60')"
    now="$(date +%s)"
    cutoff=$(( now - ret_days * 86400 ))
    __ARCH_MOVED_VERDICTS=0
    shopt -s nullglob
    for f in "$audit"/*.json "$audit"/*.md; do
        [ -f "$f" ] || continue
        base="$(basename "$f")"
        # 排除项:ledger.jsonl 与 *-summary.md
        case "$base" in
            ledger.jsonl|*-summary.md) continue ;;
        esac
        ep="$(_archive_verdict_epoch "$f")"
        [ -z "$ep" ] || [ "$ep" = "0" ] && continue
        [ "$ep" -lt "$cutoff" ] || continue
        ym="$(date -u -r "$ep" +%Y%m 2>/dev/null || date -u -d "@$ep" +%Y%m 2>/dev/null)"
        [ -n "$ym" ] || continue
        destdir="$audit/archive/$ym"
        __ARCH_MOVED_VERDICTS=$((__ARCH_MOVED_VERDICTS+1))
        if [ "$dry" = "1" ]; then
            echo "  DRY verdict mv: $base → archive/$ym/ (age=$(( (now-ep)/86400 ))d)" >&2
        else
            mkdir -p "$destdir" 2>/dev/null
            mv -f "$f" "$destdir/$base" 2>/dev/null && echo "  verdict mv: $base → archive/$ym/" >&2
        fi
    done
    echo "  verdicts: moved=$__ARCH_MOVED_VERDICTS (retention=${ret_days}d)"
}

# _archive_workdir_gc <dry> — 工作目录 gc;仅匹配 ^[0-9]{8,}-.+$ 目录形态;
#   显式排除 .plugin-cache/、t0mech-scanned.json、t1-last-run.json、manual-*;tfstate 非空绝不删。
__ARCH_MOVED_WORKDIRS=0
_archive_workdir_gc() {
    local dry="$1"
    local base ret_days now cutoff d name mt state_n
    base="$(probe_workdir_base)"
    [ -d "$base" ] || { echo "  workdir: 无 $base 目录,跳过"; return 0; }
    ret_days="$(cfg_or '.limits.workdir_retention_days' '7')"
    now="$(date +%s)"
    cutoff=$(( now - ret_days * 86400 ))
    __ARCH_MOVED_WORKDIRS=0
    shopt -s nullglob
    for d in "$base"/*/; do
        d="${d%/}"
        name="$(basename "$d")"
        # 排除项:.plugin-cache、manual-* 前缀(非文件形态的排除)
        case "$name" in
            .plugin-cache|manual-*) continue ;;
        esac
        # 目录形态限定:契约 ^[0-9]{8,}-.+$(<ts>-<sid> 形态);实际 runner 用 YYYYMMDDTHHMMSSZ-<sid>,
        # 中间夹 T/Z 非纯数字,故放宽为 ^[0-9]{8,}[A-Za-z0-9]*-.+$ —— 覆盖两种 <ts>-<sid> 写法,
        # 仍排除 .plugin-cache / manual-* / notmatching-dir 等非 runner 目录。
        [[ "$name" =~ ^[0-9]{8,}[A-Za-z0-9]*-.+$ ]] || continue
        # tfstate 存在且 resources 非空 → 绝不删(sweep 残留即停语义)
        if [ -f "$d/terraform.tfstate" ]; then
            state_n="$(jq -r '.resources | length' "$d/terraform.tfstate" 2>/dev/null)"
            [ -z "$state_n" ] && state_n=0
            [ "$state_n" -gt 0 ] && continue
        fi
        # mtime 判老
        mt="$(stat -f %m "$d" 2>/dev/null || stat -c %Y "$d" 2>/dev/null || echo "$now")"
        [ "$mt" -lt "$cutoff" ] || continue
        __ARCH_MOVED_WORKDIRS=$((__ARCH_MOVED_WORKDIRS+1))
        if [ "$dry" = "1" ]; then
            echo "  DRY workdir rm: $name (age=$(( (now-mt)/86400 ))d)" >&2
        else
            rm -rf "$d" 2>/dev/null && echo "  workdir rm: $name" >&2
        fi
    done
    # 排除文件:t0mech-scanned.json / t1-last-run.json 不被扫到(不是目录),但显式声明纪律
    echo "  workdir: gc=$__ARCH_MOVED_WORKDIRS (retention=${ret_days}d; 排除 .plugin-cache/manual-*/tfstate 非空;不动 t0mech-scanned.json/t1-last-run.json)"
}

# _archive_plugin_cache_report — 报告与 provider.version + 各场景 provider_version_from 都不符的版本(只报体积不删)。
_archive_plugin_cache_report() {
    local base cache_dir pv from set_of pgroot d y v hits=0 dir
    base="$(probe_workdir_base)"
    cache_dir="$base/.plugin-cache"
    [ -d "$cache_dir" ] || { echo "  plugin-cache: 无 $cache_dir 目录,跳过"; return 0; }
    pv="$(cfg '.provider.version')"
    # 场景声明的 provider_version_from 并集
    pgroot="$(probe_playground_dir)"
    from=""
    shopt -s nullglob
    for d in "$pgroot"/*/*/; do
        y="$d/scenario.yaml"; [ -f "$y" ] || continue
        v="$(_yaml_get "$y" provider_version_from)"
        [ -n "$v" ] && from="$from $v"
    done
    set_of=" $pv $from "  # 空格包夹便于 grep -F 精确匹配
    # registry.terraform.io/aliyun/alicloud/<version>/... 结构下,version 为目录名
    for dir in "$cache_dir"/registry.terraform.io/aliyun/alicloud/*/; do
        [ -d "$dir" ] || continue
        v="$(basename "${dir%/}")"
        case "$set_of" in
            *" $v "*) continue ;;
        esac
        local sz; sz="$(du -sh "$dir" 2>/dev/null | awk '{print $1}')"
        echo "  plugin-cache 陌生版本: $v ($sz) 与 provider.version=$pv 及任何场景 provider_version_from 都不符" >&2
        hits=$((hits+1))
    done
    echo "  plugin-cache: 陌生版本=$hits(只报体积不删,人工评估)"
}

# _archive_todos — 待办清单(pending drafts / _quarantine / origin: generated 未校订)
_archive_todos() {
    local root pgroot d y n qc qdir base
    root="$(probe_root)"
    pgroot="$(probe_playground_dir)"
    # pending-review drafts
    local pend
    pend="$(echo "$__ARCH_PENDING_DRAFTS" | tr ' ' '\n' | grep . || true)"
    if [ -n "$pend" ]; then
        echo "  TODO drafts pending-review:"
        while IFS= read -r base; do [ -n "$base" ] && echo "    - $base"; done <<< "$pend"
    fi
    # playground _quarantine
    qdir="$pgroot/_quarantine"
    qc=0
    if [ -d "$qdir" ]; then
        shopt -s nullglob
        for d in "$qdir"/*/*/; do
            [ -d "$d" ] || continue
            qc=$((qc+1))
            echo "    - _quarantine: ${d#$pgroot/}" >&2
        done
    fi
    echo "  TODO _quarantine: $qc"
    # origin: generated 未校订(scenario.yaml 有 origin: generated)
    n=0
    shopt -s nullglob
    for d in "$pgroot"/*/*/; do
        y="$d/scenario.yaml"; [ -f "$y" ] || continue
        [ "$(_yaml_get "$y" origin)" = "generated" ] && n=$((n+1))
    done
    echo "  TODO origin=generated 未校订: $n"
}

_cmd_archive() {
    local dry=0
    while [ $# -gt 0 ]; do
        case "$1" in
            --dry) dry=1; shift ;;
            -*) echo "archive: 未知参数 '$1'" >&2; return 2 ;;
            *) echo "archive: 多余参数 '$1'" >&2; return 2 ;;
        esac
    done
    if [ "$dry" = "1" ]; then echo "archive plan (dry-run,不做真操作):"; else echo "archive run:"; fi
    _archive_drafts    "$dry"
    _archive_verdicts  "$dry"
    _archive_workdir_gc "$dry"
    _archive_plugin_cache_report
    _archive_todos
    # A2 ledger:archive 追加一行(dry 也追加,便于审计;dry 用 kind:"archive_dry")
    local kind; [ "$dry" = "1" ] && kind="archive_dry" || kind="archive"
    _ledger_append \
        '{ts:$ts, kind:$k, moved:{drafts:$md, verdicts:$mv, workdirs:$mw}}' \
        --arg k "$kind" \
        --argjson md "$__ARCH_MOVED_DRAFTS" \
        --argjson mv "$__ARCH_MOVED_VERDICTS" \
        --argjson mw "$__ARCH_MOVED_WORKDIRS"
    return 0
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
        doctor)  _cmd_doctor "$@" ;;
        list)    _cmd_list "$@" ;;
        tier0)   _cmd_tier0 "$@" ;;
        run)     _cmd_run "$@" ;;
        sweep)   _cmd_sweep "$@" ;;
        archive) _cmd_archive "$@" ;;
        -h|--help|"") _usage ;;
        *) echo "probe: 未知命令 '$cmd'" >&2; _usage; exit 2 ;;
    esac
}

[[ "${BASH_SOURCE[0]}" == "$0" ]] && main "$@"
