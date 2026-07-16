#!/usr/bin/env bash
# bootstrap/reconcile.sh — 收敛族统一入口(P1.c 合并了 sweep + watchdog + 原 reconcile)
#
# 子命令:
#   stale     — 清扫 jarvis-claimed 超 TTL 的工单 → escalate(原 sweep.sh)
#   orphan    — owner_instance 已死的 task → escalate(原 watchdog.sh)
#   drift     — 台账 vs Aone 对账,claim+seen 但缺 release → 补 release(原 reconcile.sh)
#   donecheck — jarvis-done 标签 vs Aone 状态一致性对账:标签 done 但状态落在合法完成态集合
#               外(finish 状态被拒 / 人工打回返工后标签滞留)→ escalate 告警(纯只读,不动标签/状态)
#   all       — 顺跑 stale → orphan → drift → donecheck(默认)
#
# 用法:
#   reconcile.sh                  # 默认 all
#   reconcile.sh {stale|orphan|drift|donecheck|all}
#
# 环境变量:
#   JARVIS_ROOT           — repo root(默认 git rev-parse --show-toplevel)
#   JARVIS_ESCALATION_DIR — escalation dir(默认 <root>/escalation)
#   JARVIS_RUNS_DIR       — runs dir(默认 <root>/runs)
#   RECONCILE_CLAIM_CMD   — 覆盖 claim.sh 路径(测试用)
#   JARVIS_RECONCILE_SKIP_IDS — comma-separated IDs to skip (active EphemeralExecutor jobs)
#
# Read-mostly;仅 drift 分支可能触发 claim.sh release。

set -uo pipefail

_reconcile_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=bootstrap/lib.sh
source "$_reconcile_dir/lib.sh"
JARVIS_ROOT="$(jarvis_root)"
export JARVIS_ROOT
# shellcheck source=bootstrap/log.sh
source "$_reconcile_dir/log.sh"

POOLS_JSON="$JARVIS_ROOT/config/pools.json"

# Active EphemeralExecutor job IDs — reconcile must not touch these.
# bash 3.2 兼容:不用关联数组(declare -A / arr["k"]),把 skip 列表规整为逗号包裹字符串
# `,id1,id2,`,判定时 case 模式匹 `*",$id,"*`;语义保持:逗号分隔、容忍空白、空值=无跳过。
_SKIP_IDS_RAW="${JARVIS_RECONCILE_SKIP_IDS:-}"
_SKIP_IDS_NORM=""
if [ -n "$_SKIP_IDS_RAW" ]; then
    _clean="$(printf '%s' "$_SKIP_IDS_RAW" | tr -d '[:space:]')"
    if [ -n "$_clean" ]; then
        _SKIP_IDS_NORM=",${_clean},"
        # 折叠空段(应对 ",,id,," 类脏输入),避免误判空 id 命中
        while [ "$_SKIP_IDS_NORM" != "${_SKIP_IDS_NORM/,,/,}" ]; do
            _SKIP_IDS_NORM="${_SKIP_IDS_NORM/,,/,}"
        done
    fi
fi

_is_active_dispatch() {
    [ -n "$_SKIP_IDS_NORM" ] || return 1
    case "$_SKIP_IDS_NORM" in
        *",$1,"*) return 0 ;;
        *)        return 1 ;;
    esac
}

# a1 via bin/a1id → act as the jarvis identity regardless of ambient login (CLAUDE.md #6).
# Overridable via JARVIS_A1 (tests point it at a stubbed `a1` on PATH).
A1="${JARVIS_A1:-$JARVIS_ROOT/bin/a1id --}"

# ==== 共享 helpers ====

_read_claim_tags() {
    CLAIM_TAG="$(jq -r '.claim.tag' "$POOLS_JSON")"
    IDLE_TAG="$(jq -r '.claim.idle_tag' "$POOLS_JSON")"
    DONE_TAG="$(jq -r '.claim.done_tag' "$POOLS_JSON")"
}

_projects() {
    jq -r '.pools[].project | select(. != null)' "$POOLS_JSON"
}

_claimed_items() {
    local project="$1"
    $A1 project workitem list \
        --project "$project" \
        --tag "$CLAIM_TAG" \
        --filter "NOT tag=$IDLE_TAG AND NOT tag=$DONE_TAG" \
        -f json 2>/dev/null || echo "[]"
}

_extract_ids() {
    python3 -c "
import json, sys
try:
    items = json.loads(sys.stdin.read())
    for item in items:
        item_id = item.get('identifier', item.get('id', ''))
        if item_id: print(str(item_id))
except Exception:
    pass
" 2>/dev/null || true
}

_extract_title() {
    local item_id="$1"
    local json="$2"
    printf '%s' "$json" | python3 -c "
import json, sys
try:
    items = json.loads(sys.stdin.read())
    target = '$item_id'
    for item in items:
        iid = str(item.get('identifier', item.get('id', '')))
        if iid == target:
            print(item.get('subject', item.get('title', '')))
            break
except Exception:
    pass
" 2>/dev/null || true
}

# ==== stale (原 sweep.sh) ====

_parse_utc_epoch() {
    local ts="$1" epoch=""
    if epoch=$(date -u -j -f "%Y-%m-%dT%H:%M:%SZ" "$ts" +%s 2>/dev/null); then echo "$epoch"; return; fi
    if epoch=$(date -u -d "$ts" +%s 2>/dev/null); then echo "$epoch"; return; fi
    if epoch=$(python3 -c "
import sys, calendar, datetime
try:
    dt = datetime.datetime.strptime('$ts', '%Y-%m-%dT%H:%M:%SZ')
    print(int(calendar.timegm(dt.timetuple())))
except Exception:
    sys.exit(1)
" 2>/dev/null); then echo "$epoch"; return; fi
    echo ""
}

_cmd_stale() {
    _read_claim_tags
    local TTL_MIN
    TTL_MIN="$(jq -r '.claim.ttl_min' "$POOLS_JSON")"
    local TTL_SECS=$((TTL_MIN * 60))
    local NOW_EPOCH
    NOW_EPOCH=$(date -u +%s)

    local STALE_IDS=()

    while IFS= read -r project; do
        [ -z "$project" ] && continue
        local claimed_json item_ids
        claimed_json=$(_claimed_items "$project")
        item_ids=$(printf '%s' "$claimed_json" | _extract_ids)

        while IFS= read -r item_id; do
            [ -z "$item_id" ] && continue
            # Skip IDs actively being processed by EphemeralExecutor
            if _is_active_dispatch "$item_id"; then
                echo "SKIP(active): $item_id"
                continue
            fi
            # 幂等：已 escalate 过的工单不再重复追加
            local _esc_dir="${JARVIS_ESCALATION_DIR:-$JARVIS_ROOT/escalation}"
            [ -f "$_esc_dir/$item_id.md" ] && continue
            local comments_json latest_ts claim_epoch age_secs age_min
            comments_json=$($A1 project workitem comment list "$item_id" -f json 2>/dev/null || echo "[]")
            latest_ts=$(printf '%s' "$comments_json" | python3 -c "
import json, sys, re
try:
    comments = json.loads(sys.stdin.read())
    pattern = re.compile(r'jarvis-claim\s+\S+\s+(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)')
    timestamps = []
    for c in comments:
        content = c.get('content', '')
        m = pattern.search(content)
        if m: timestamps.append(m.group(1))
    if timestamps: print(sorted(timestamps)[-1])
except Exception:
    pass
" 2>/dev/null || true)
            [ -z "$latest_ts" ] && continue
            claim_epoch=$(_parse_utc_epoch "$latest_ts")
            if [ -z "$claim_epoch" ]; then
                echo "WARN: could not parse claim timestamp '$latest_ts' for item $item_id, skipping" >&2
                continue
            fi
            age_secs=$(( NOW_EPOCH - claim_epoch ))
            if [ "$age_secs" -gt "$TTL_SECS" ]; then
                age_min=$(( age_secs / 60 ))
                local item_title
                item_title=$(_extract_title "$item_id" "$claimed_json")
                escalate "$item_id" "stale claim >${TTL_MIN}min (age: ${age_min}min)" "$item_title"
                STALE_IDS+=("$item_id")
            fi
        done <<< "$item_ids"
    done <<< "$(_projects)"

    if [ "${#STALE_IDS[@]}" -gt 0 ]; then
        echo "STALE_CLAIMS: ${STALE_IDS[*]}"
    else
        echo "STALE_CLAIMS: none"
    fi
}

# ==== orphan (原 watchdog.sh) ====

_cmd_orphan() {
    local _esc_dir="${JARVIS_ESCALATION_DIR:-$JARVIS_ROOT/escalation}"
    local escaped=0
    while IFS= read -r aid; do
        [ -z "$aid" ] && continue
        [ -f "$_esc_dir/$aid.md" ] && continue
        escalate "$aid" "owner dead, awaiting adopt"
        escaped=1
    done < <(bash "$_reconcile_dir/coord.sh" list-orphans)
    [ "$escaped" -eq 0 ] && echo "ORPHANS: none"
}

# ==== drift (原 reconcile.sh) ====

_cmd_drift() {
    _read_claim_tags
    local RECONCILED_IDS=()

    while IFS= read -r project; do
        [ -z "$project" ] && continue
        local claimed_json item_ids
        claimed_json=$(_claimed_items "$project")
        item_ids=$(printf '%s' "$claimed_json" | _extract_ids)

        while IFS= read -r item_id; do
            [ -z "$item_id" ] && continue
            # Skip IDs actively being processed by EphemeralExecutor
            if _is_active_dispatch "$item_id"; then
                echo "SKIP(active): $item_id"
                continue
            fi
            if seen "$item_id"; then
                local _claim_cmd="${RECONCILE_CLAIM_CMD:-$_reconcile_dir/claim.sh}"
                $_claim_cmd release "$item_id" "$project"
                echo "RECONCILED: $item_id"
                RECONCILED_IDS+=("$item_id")
            fi
            # 无 run 文件 → 留给 stale 处理
        done <<< "$item_ids"
    done <<< "$(_projects)"

    if [ "${#RECONCILED_IDS[@]}" -eq 0 ]; then
        echo "RECONCILED: none"
    fi
}

# ==== donecheck (jarvis-done vs Aone status 一致性对账) ====

# 合法完成态全集 = .claim.done_statuses ∪ 各池 .done_status 叶子值(string 或 object 的 value)。
# 后者含 tf_provider 的「待发布」等非终态但合法的 jarvis-done 停靠态 —— 只在此集合外才算漂移。
_legit_done_statuses() {
    # LC_ALL=C 关键：默认/CJK locale 下 `sort -u` 会把 collation 相等的中文状态(如 已完成/已拒绝)
    # 误判为重复并丢行,导致合法完成态漏进集合 → 误报漂移。字节序去重才安全。
    jq -r '
      ((.claim.done_statuses // [])[]),
      (.pools[]?.done_status | if type=="object" then .[] elif type=="string" then . else empty end)
    ' "$POOLS_JSON" 2>/dev/null | awk 'NF' | LC_ALL=C sort -u
}

# 点读单个工单当前 status displayValue（镜像 claim.sh _get_status，读不到回空 → 调用方跳过不误报）。
_workitem_status() {
    local id="$1" json
    json="$($A1 project workitem get "$id" -f json 2>/dev/null)" || return 0
    printf '%s' "$json" | jq -r '
        (.fields // [])[] | select(.identifier=="status") | .displayValue // empty
    ' 2>/dev/null || return 0
}

_done_tagged_items() {
    local project="$1"
    $A1 project workitem list \
        --project "$project" \
        --tag "$DONE_TAG" \
        -f json 2>/dev/null || echo "[]"
}

_cmd_donecheck() {
    _read_claim_tags
    local legit; legit="$(_legit_done_statuses)"
    local flagged=0
    while IFS= read -r project; do
        [ -z "$project" ] && continue
        local json ids
        json="$(_done_tagged_items "$project")"
        ids="$(printf '%s' "$json" | _extract_ids)"
        while IFS= read -r item_id; do
            [ -z "$item_id" ] && continue
            _is_active_dispatch "$item_id" && { echo "SKIP(active): $item_id"; continue; }
            local st; st="$(_workitem_status "$item_id")"
            # 读不到状态 → 无法判定,跳过不误报(与 claim.sh 的 graceful-degrade 一致)。
            [ -z "$st" ] && continue
            if printf '%s\n' "$legit" | grep -qxF "$st"; then
                continue   # 合法完成态,标签与真源一致
            fi
            echo "DRIFT(done_status): $item_id status='$st' not in legit-done set (project=$project)"
            local title; title="$(_extract_title "$item_id" "$json")"
            escalate "$item_id" \
                "done_status_drift: 标签 jarvis-done 但 Aone 状态='$st' 不在合法完成态集合(finish 状态被拒 或 人工打回返工后标签滞留)——请人工核对状态或摘掉 jarvis-done 标签" \
                "$title"
            flagged=$((flagged + 1))
        done <<< "$ids"
    done <<< "$(_projects)"
    if [ "$flagged" -eq 0 ]; then echo "DONECHECK: none"; else echo "DONECHECK: flagged $flagged"; fi
}

# ==== CLI dispatch ====

_help() { sed -n '2,21p' "$0"; }

cmd="${1:-all}"
case "$cmd" in
    stale)     _cmd_stale ;;
    orphan)    _cmd_orphan ;;
    drift)     _cmd_drift ;;
    donecheck) _cmd_donecheck ;;
    all)       _cmd_stale; _cmd_orphan; _cmd_drift; _cmd_donecheck ;;
    -h|--help) _help ;;
    *)
        echo "reconcile: unknown command '$cmd'" >&2
        _help
        exit 1 ;;
esac

exit 0
