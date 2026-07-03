#!/usr/bin/env bash
# bootstrap/claim.sh – claim/release/finish primitive for concurrent multi-machine agents.
# Usage:
#   claim.sh claim   <workitem_id> <project_id>
#   claim.sh release <workitem_id> <project_id>
#   claim.sh finish  <workitem_id> <project_id>
#
# claim:   tags the workitem jarvis-claimed, freezes a `jarvis-claim <host> <utc>` prefix
#          into .my-day/claim-prefix-<id>.txt (consumed by wrap.sh sync/done into the
#          first business comment to avoid an independent claim-only comment), then
#          point-reads the workitem to verify our claimed tag landed. Exits 0 on success,
#          1 if the claimed tag is not visible in the readback (lost race; prefix cleaned).
# release: tags the workitem jarvis-idle —— 本轮 jarvis 处理完毕、释放锁、等待人或下一个 jarvis
#          接手；不动 Aone status。cleans up any leftover prefix. Exits 0.
# finish:  tags the workitem jarvis-done + 改 Aone status 为 .claim.done_status（默认"已发布
#          待需求排期"）—— jarvis 判断工单真完成。cleans up any leftover prefix. Exits 0.
#
# Tag writes preserve pre-existing tags: `a1 ... update --tag` is a whole-set replace
# (comma-separated names), so we point-read the current tags, apply the migration set
# op, and write the full expected list in one update. This keeps unrelated tags such as
# jarvis-probe across the claim→release→finish lifecycle.
#
# Reads tag names from config/pools.json (.claim.tag / .claim.idle_tag / .claim.done_tag /
# .claim.done_status). Respects JARVIS_ROOT env override for the repo root.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"
pools_cfg="$jarvis_root/config/pools.json"

if [ ! -f "$pools_cfg" ]; then
    echo "claim.sh: config/pools.json not found at $pools_cfg" >&2
    exit 1
fi

# Read tag names from config
CLAIM_TAG=$(jq -r '.claim.tag' "$pools_cfg")
IDLE_TAG=$(jq -r '.claim.idle_tag' "$pools_cfg")
DONE_TAG=$(jq -r '.claim.done_tag' "$pools_cfg")
DONE_STATUS=$(jq -r '.claim.done_status' "$pools_cfg")

cmd="${1:-}"
workitem_id="${2:-}"
project_id="${3:-}"

if [ -z "$cmd" ] || [ -z "$workitem_id" ] || [ -z "$project_id" ]; then
    echo "Usage: claim.sh <claim|release|finish> <workitem_id> <project_id>" >&2
    exit 1
fi

myday_dir="$jarvis_root/.my-day"
today_date="$(date -u +%F)"
ledger_file="$myday_dir/claims-${today_date}.json"

_claim_prefix_path() { echo "$myday_dir/claim-prefix-$1.txt"; }

# Atomically append/update entry in the claims ledger.
# Usage: _ledger_upsert <id> <done_value>   (done_value: true|false)
_ledger_upsert() {
    local id="$1"
    local done_val="$2"
    mkdir -p "$myday_dir"
    local tmp
    tmp="$(mktemp "$myday_dir/.claims-tmp.XXXXXX")"
    # Seed tmp with existing ledger or empty array; never create a partial ledger_file first
    if [ -f "$ledger_file" ]; then
        jq --arg id "$id" --argjson done "$done_val" \
            'if any(.[]; .id == $id) then
                 map(if .id == $id then .done = $done else . end)
             else
                 . + [{"id": $id, "done": $done}]
             end' \
            "$ledger_file" > "$tmp" && mv "$tmp" "$ledger_file"
    else
        printf '[{"id":"%s","done":%s}]' "$id" "$done_val" > "$tmp" && mv "$tmp" "$ledger_file"
    fi
}

# ---------------------------------------------------------------------------
# Tag-preserving primitives
# ---------------------------------------------------------------------------

# Point-read current tag NAMES for a workitem, echoed comma-joined (may be empty).
# Uses `workitem get` (strong point-read) instead of `workitem list --tag ...` (a
# search-index query whose lag can exceed the retry window → the deterministic
# false-negative "lost race" this fix removes). Returns non-zero only if the get
# call itself fails, so callers can degrade to a legacy bare-tag write.
_get_tags() {
    local id="$1" json
    json="$(a1 project workitem get "$id" -f json 2>/dev/null)" || return 1
    # get -f json exposes fields[]; the tag field is identifier=="tag", whose
    # displayValue is the comma-joined human tag names (format multiList). Empty/
    # single/multi all parse: null → "", "a" → "a", "a,b" → "a,b".
    printf '%s' "$json" | jq -r '
        (.fields // []) | map(select(.identifier=="tag")) | (.[0].displayValue // "")
    ' 2>/dev/null || return 1
}

# Echo comma-joined result of (existing ∪ add) − remove over tag-name sets.
# Args: <existing_csv> <add_csv> <remove_csv> (any may be empty; comma-separated).
_compute_tags() {
    jq -rn --arg ex "$1" --arg add "$2" --arg rm "$3" '
        def sc: if . == "" then []
                else split(",") | map(gsub("^\\s+|\\s+$";"")) | map(select(length>0)) end;
        ((($ex|sc) + ($add|sc)) | unique) - ($rm|sc) | join(",")
    '
}

# True (exit 0) if comma-joined <csv> contains tag-name <needle>.
_csv_contains() {
    jq -e -n --arg csv "$1" --arg n "$2" \
        '($csv | split(",") | map(gsub("^\\s+|\\s+$";""))) | index($n) != null' \
        >/dev/null 2>&1
}

# Write the full expected tag set in one update, preserving pre-existing tags.
# Args: <id> <add_csv> <remove_csv> [extra a1 update flags...].
# Degrades to a legacy bare `--tag <add>` write (may drop other tags) with a stderr
# warning if the current tags can't be read, so claim/release/finish never hard-fail
# on a transient get error.
_update_tags_merged() {
    local id="$1" add="$2" remove="$3"; shift 3
    local existing new
    if existing="$(_get_tags "$id")"; then
        new="$(_compute_tags "$existing" "$add" "$remove")"
        a1 project workitem update "$id" --tag "$new" "$@"
    else
        echo "claim.sh: warning: could not read existing tags for $id; writing only '$add' (pre-existing tags may be lost)" >&2
        a1 project workitem update "$id" --tag "$add" "$@"
    fi
}

case "$cmd" in
    claim)
        # 1. Tag as claimed, preserving any pre-existing tags: existing ∪ {claimed} − {idle}
        _update_tags_merged "$workitem_id" "$CLAIM_TAG" "$IDLE_TAG"

        # 2. 冻结 hostname + utc 到 prefix 文件；wrap.sh sync/done 第一次发业务评论时自动 prefix
        # 到顶部并消费，避免 claim 时刻额外刷一条独立留痕评论。
        # 优先 ComputerName(用户自定义显示名)，回落 hostname -s 短名，最后 hostname；避免 mac 默认带 .local 后缀
        utcnow=$(date -u +%Y-%m-%dT%H:%M:%SZ)
        host=$(scutil --get ComputerName 2>/dev/null || hostname -s 2>/dev/null || hostname)
        mkdir -p "$myday_dir"
        printf 'jarvis-claim %s %s' "$host" "$utcnow" > "$(_claim_prefix_path "$workitem_id")"

        # 3. Readback via POINT-READ (workitem get), retries for a small safety margin.
        # Concurrency semantics: the point-read confirms *our own* write landed; it does
        # NOT arbitrate a true race — two instances claiming near-simultaneously will both
        # see jarvis-claimed and both exit 0 here. Mutual exclusion still rests on the
        # ledger + reconcile stale/drift sweep (unchanged; not a regression vs the old
        # list-based readback, which likewise could not distinguish owners — it only added
        # a deterministic false-negative from search-index lag).
        for _ in 1 2 3; do
            if _csv_contains "$(_get_tags "$workitem_id" 2>/dev/null || echo "")" "$CLAIM_TAG"; then
                echo "claim.sh: claimed workitem $workitem_id in project $project_id"
                _ledger_upsert "$workitem_id" false
                exit 0
            fi
            sleep "${JARVIS_CLAIM_READBACK_SLEEP:-2}"
        done
        # 抢锁失败：清理本机 prefix 痕迹（不留给别人）
        rm -f "$(_claim_prefix_path "$workitem_id")"
        echo "claim.sh: lost race for workitem $workitem_id (claimed tag not visible in point-read readback)" >&2
        exit 1
        ;;

    release)
        # Tag as idle, preserving other tags: existing − {claimed} ∪ {idle}
        # 本轮 jarvis 处理完，释放锁；后续等待人或下一个 jarvis 接手（不动 Aone status）
        _update_tags_merged "$workitem_id" "$IDLE_TAG" "$CLAIM_TAG"
        # 兜底：如果整个 jarvis 周期没发过业务评论，prefix 文件仍在，此处清理
        rm -f "$(_claim_prefix_path "$workitem_id")"
        echo "claim.sh: released workitem $workitem_id in project $project_id"
        _ledger_upsert "$workitem_id" true
        exit 0
        ;;

    finish)
        # Tag as done, preserving other tags: existing − {claimed, idle} ∪ {done}；同时改 status
        _update_tags_merged "$workitem_id" "$DONE_TAG" "$CLAIM_TAG,$IDLE_TAG" --status "$DONE_STATUS"
        rm -f "$(_claim_prefix_path "$workitem_id")"
        echo "claim.sh: finished workitem $workitem_id in project $project_id (status=$DONE_STATUS)"
        _ledger_upsert "$workitem_id" true
        exit 0
        ;;

    *)
        echo "claim.sh: unknown command '$cmd'. Use claim, release, or finish." >&2
        exit 1
        ;;
esac
