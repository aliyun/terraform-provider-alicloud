#!/usr/bin/env bash
# bootstrap/probe-meta.sh — tier-0 元数据获取层(薄封装 amp-resource-metadata skill 的
# get_api_definition.py)。probe.sh 经此取 OpenAPI 侧 API 定义,不直接调 python。
#
# 子命令:
#   fetch <product> <version> <action>         — 取 API 定义原始 JSON(body);网络/脚本/凭证不可用→退非零+明确提示
#   cached-fetch <product> <version> <action>  — 同 fetch,套 cache.sh TTL 缓存(默认 7d,全量巡检才跑得起)
#   clear [<product> <version> <action>]       — 清缓存(给参清该 key;无参清全部 probemeta_* key)
#   available                                  — 探能力(python/脚本/凭证或自定义 fetcher);可用退 0,否则退非零+原因
#
# 设计:精度与降级并重。任一环节(无 venv python / 无脚本 / 无凭证 / 网络失败 / 空响应)→ 干净退非零,
#   调用方(probe.sh tier0)据此自动降级为纯 doc↔source + 全 judgment_queue。
#
# 环境变量:
#   AMP_SKILL_DIR       — amp-resource-metadata skill 目录(默认 <jarvis_root>/.claude/skills/amp-resource-metadata)
#   PROBE_META_PYTHON   — 覆盖底层 python 解释器(默认 venv python > 系统 python3);测试/power-user 用作 fetcher 桩
#   PROBE_META_TTL      — 缓存 TTL 秒(默认 604800 = 7 天)
#   AMP_ACCESS_KEY_ID / AMP_ACCESS_KEY_SECRET(优先)或 ALIBABA_CLOUD_ACCESS_KEY_ID / _SECRET — API 凭证(仅判 set,绝不打印值)
#
# 被 source 时不执行 main(便于单测)。
set -uo pipefail

_pm_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/lib.sh
source "$_pm_dir/lib.sh"

pm_skill_dir()  { echo "${AMP_SKILL_DIR:-$(jarvis_root)/.claude/skills/amp-resource-metadata}"; }
pm_script()     { echo "$(pm_skill_dir)/scripts/get_api_definition.py"; }
pm_venv_py()    { echo "$(pm_skill_dir)/scripts/.venv/bin/python3"; }
pm_ttl()        { echo "${PROBE_META_TTL:-604800}"; }

# 解释器解析:显式覆盖 > venv python > 系统 python3。返回空表示无可用解释器。
pm_python() {
    if [ -n "${PROBE_META_PYTHON:-}" ]; then echo "$PROBE_META_PYTHON"; return; fi
    local venv; venv="$(pm_venv_py)"
    if [ -x "$venv" ]; then echo "$venv"; return; fi
    command -v python3 >/dev/null 2>&1 && { echo python3; return; }
    echo ""
}

_pm_have_creds() {
    { [ -n "${AMP_ACCESS_KEY_ID:-}" ] && [ -n "${AMP_ACCESS_KEY_SECRET:-}" ]; } && return 0
    { [ -n "${ALIBABA_CLOUD_ACCESS_KEY_ID:-}" ] && [ -n "${ALIBABA_CLOUD_ACCESS_KEY_SECRET:-}" ]; } && return 0
    return 1
}

pm_cache_key() { echo "probemeta_${1}_${2}_${3}"; }

# available — 静态能力探测(供 doctor / tier0 决策降级)。不发网络请求。
_pm_available() {
    # 自定义 fetcher(PROBE_META_PYTHON):视为可用(桩/power-user 自带 auth)
    if [ -n "${PROBE_META_PYTHON:-}" ]; then
        [ -x "$PROBE_META_PYTHON" ] || command -v "$PROBE_META_PYTHON" >/dev/null 2>&1 || {
            echo "probe-meta: PROBE_META_PYTHON='$PROBE_META_PYTHON' 不可执行" >&2; return 1; }
        return 0
    fi
    local venv; venv="$(pm_venv_py)"
    if [ ! -x "$venv" ]; then
        echo "probe-meta: 无 venv python($venv);先跑 amp-resource-metadata/scripts/setup.sh" >&2
        return 1
    fi
    if [ ! -f "$(pm_script)" ]; then
        echo "probe-meta: 缺 get_api_definition.py($(pm_script))" >&2
        return 1
    fi
    if ! _pm_have_creds; then
        echo "probe-meta: 缺 AMP_/ALIBABA_CLOUD_ 凭证(仅判 set;白名单见 skill SKILL.md)" >&2
        return 1
    fi
    return 0
}

# fetch — 取原始 API 定义 JSON(get_api_definition.py --output json 的 body)。
_pm_fetch() {
    local product="${1:-}" version="${2:-}" action="${3:-}"
    if [ -z "$product" ] || [ -z "$version" ] || [ -z "$action" ]; then
        echo "probe-meta fetch: 用法 fetch <product> <version> <action>" >&2; return 2
    fi
    local py; py="$(pm_python)"
    if [ -z "$py" ]; then echo "probe-meta: 无可用 python 解释器" >&2; return 4; fi
    local script; script="$(pm_script)"
    if [ ! -f "$script" ] && [ -z "${PROBE_META_PYTHON:-}" ]; then
        echo "probe-meta: 缺脚本 $script" >&2; return 4
    fi
    local rc tmp_out tmp_err
    tmp_out="$(mktemp)"; tmp_err="$(mktemp)"
    "$py" "$script" --product "$product" --api-version "$version" --api "$action" --output json >"$tmp_out" 2>"$tmp_err"
    rc=$?
    if [ "$rc" -ne 0 ]; then
        sed 's/^/probe-meta(fetch err): /' "$tmp_err" >&2
        rm -f "$tmp_out" "$tmp_err"; return "$rc"
    fi
    if [ ! -s "$tmp_out" ]; then
        echo "probe-meta: $product/$version/$action 空响应" >&2
        rm -f "$tmp_out" "$tmp_err"; return 3
    fi
    cat "$tmp_out"
    rm -f "$tmp_out" "$tmp_err"
    return 0
}

_pm_cached_fetch() {
    local product="${1:-}" version="${2:-}" action="${3:-}"
    if [ -z "$product" ] || [ -z "$version" ] || [ -z "$action" ]; then
        echo "probe-meta cached-fetch: 用法 cached-fetch <product> <version> <action>" >&2; return 2
    fi
    local key; key="$(pm_cache_key "$product" "$version" "$action")"
    bash "$_pm_dir/cache.sh" get "$key" "$(pm_ttl)" -- bash "$_pm_dir/probe-meta.sh" fetch "$product" "$version" "$action"
}

_pm_clear() {
    if [ "$#" -ge 3 ]; then
        bash "$_pm_dir/cache.sh" bust "$(pm_cache_key "$1" "$2" "$3")"
        return 0
    fi
    # 无参:清全部 probemeta_* 缓存
    local cdir="${JARVIS_CACHE_DIR:-$(jarvis_root)/.my-day/cache}"
    rm -f "$cdir"/probemeta_* 2>/dev/null
    return 0
}

main() {
    local cmd="${1:-}"; shift 2>/dev/null || true
    case "$cmd" in
        fetch)        _pm_fetch "$@" ;;
        cached-fetch) _pm_cached_fetch "$@" ;;
        clear)        _pm_clear "$@" ;;
        available)    _pm_available ;;
        -h|--help|"") sed -n '2,32p' "$0" ;;
        *) echo "probe-meta: 未知命令 '$cmd'" >&2; return 2 ;;
    esac
}

[[ "${BASH_SOURCE[0]}" == "$0" ]] && main "$@"
