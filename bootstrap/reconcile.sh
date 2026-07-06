#!/usr/bin/env bash
# bootstrap/reconcile.sh — 收敛族统一入口(P1.c 合并了 sweep + watchdog + 原 reconcile)
#
# 子命令:
#   stale   — 清扫 jarvis-claimed 超 TTL 的工单 → escalate(原 sweep.sh)
#   orphan  — owner_instance 已死的 task → escalate(原 watchdog.sh)
#   drift   — 台账 vs Aone 对账,claim+seen 但缺 release → 补 release(原 reconcile.sh)
#   all     — 顺跑 stale → orphan → drift(默认)
#
# 用法:
#   reconcile.sh                  # 默认 all
#   reconcile.sh {stale|orphan|drift|all}
#
# 环境变量:
#   JARVIS_ROOT           — repo root(默认 git rev-parse --show-toplevel)
#   JARVIS_ESCALATION_DIR — escalation dir(默认 <root>/escalation)
#   JARVIS_RUNS_DIR       — runs dir(默认 <root>/runs)
#   RECONCILE_CLAIM_CMD   — 覆盖 claim.sh 路径(测试用)
#   JARVIS_RECONCILE_SKIP_IDS — comma-separated IDs to skip (active DispatchPool workers)
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

# Active DispatchPool worker IDs — reconcile must not touch these.
_SKIP_IDS_RAW="${JARVIS_RECONCILE_SKIP_IDS:-}"
declare -A _SKIP_ID_MAP
if [ -n "$_SKIP_IDS_RAW" ]; then
    IFS=',' read -ra _skip_arr <<< "$_SKIP_IDS_RAW"
    for _sid in "${_skip_arr[@]}"; do
        _sid="$(echo "$_sid" | tr -d '[:space:]')"
        [ -n "$_sid" ] && _SKIP_ID_MAP["$_sid"]=1
    done
fi

_is_active_dispatch() {
    [ "${_SKIP_ID_MAP[$1]+_}" ]
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
            # Skip IDs actively being processed by DispatchPool
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
            # Skip IDs actively being processed by DispatchPool
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

# ==== CLI dispatch ====

_help() { sed -n '2,21p' "$0"; }

cmd="${1:-all}"
case "$cmd" in
    stale)  _cmd_stale ;;
    orphan) _cmd_orphan ;;
    drift)  _cmd_drift ;;
    all)    _cmd_stale; _cmd_orphan; _cmd_drift ;;
    -h|--help) _help ;;
    *)
        echo "reconcile: unknown command '$cmd'" >&2
        _help
        exit 1 ;;
esac

exit 0
