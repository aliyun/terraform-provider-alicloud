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

# a1 via bin/a1id → act as the jarvis identity regardless of ambient login (CLAUDE.md #6).
# Overridable via JARVIS_A1 (tests point it at a stubbed `a1` on PATH).
A1="${JARVIS_A1:-$jarvis_root/bin/a1id --}"

if [ ! -f "$pools_cfg" ]; then
    echo "claim.sh: config/pools.json not found at $pools_cfg" >&2
    exit 1
fi

# Read tag names from config
CLAIM_TAG=$(jq -r '.claim.tag' "$pools_cfg")
IDLE_TAG=$(jq -r '.claim.idle_tag' "$pools_cfg")
DONE_TAG=$(jq -r '.claim.done_tag' "$pools_cfg")
DONE_STATUS=$(jq -r '.claim.done_status' "$pools_cfg")

# Resolve the done status for a project: a per-pool .done_status (matched by project id)
# overrides the global .claim.done_status. Different Aone projects have different status
# enums (e.g. 2100304 only allows 处理中/已完成/已取消, not the global "已发布待需求排期"),
# so finish must use the target project's own completion status. Echoes "" if neither set.
_pool_done_status() {
    local project="$1" ps
    ps="$(jq -r --arg p "$project" \
        '(.pools[]? | select((.project|tostring) == $p) | .done_status) // empty' \
        "$pools_cfg" 2>/dev/null)"
    if [ -n "$ps" ]; then echo "$ps"; else jq -r '.claim.done_status // empty' "$pools_cfg" 2>/dev/null; fi
}

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
#
# Each entry records the owning instance: {"id","done","owner"}. owner comes from
# coord_self() (lib.sh): COORD_ID if set, else cc-<CLAUDE_CODE_SESSION_ID> for an
# interactive/headless Claude session, else "". wrap-check.sh scopes its Stop block by
# owner (via the same coord_self) so a session only hard-blocks on its own (or
# ownerless legacy) claims — and, since D2, two different interactive sessions get
# distinct cc-<sid> owners and stop blocking each other.
#
# INSERT writes owner. UPDATE (release/finish calling _ledger_upsert id true)
# COALESCE-preserves the original owner: only backfills the current owner when the
# existing owner is empty/missing, so a different instance closing a claim never
# overwrites the original claimer. Legacy entries lacking an owner field read as "".
#
# Concurrency (D1): the read-modify-write is guarded by a mkdir lock (macOS has no
# flock). mv is atomic per-write, but without the lock two processes could each read
# the pre-update ledger and the later mv would clobber the earlier one (lost update).
# A stale lock (>10s, holder likely died mid-write) is stolen; after ~5s of contention
# we proceed unlocked rather than hang a claim.
_ledger_upsert() {
    local id="$1"
    local done_val="$2"
    local owner
    owner="$(coord_self)"
    mkdir -p "$myday_dir"

    local lock="$myday_dir/.claims.lock" acquired="" waited=0
    while [ "$waited" -lt 50 ]; do
        if mkdir "$lock" 2>/dev/null; then acquired=1; break; fi
        # Steal only a provably-old lock (holder likely died mid-write). If stat can't
        # read the mtime (lock just vanished / transient), DO NOT steal — a failed stat
        # must not read as "ancient" or we'd yank a live holder's lock and lose updates.
        local lock_mtime
        lock_mtime="$(stat -f %m "$lock" 2>/dev/null || stat -c %Y "$lock" 2>/dev/null || echo "")"
        if [ -n "$lock_mtime" ] && [ "$(( $(date +%s) - lock_mtime ))" -gt 10 ]; then
            rm -rf "$lock" 2>/dev/null
        fi
        sleep 0.1; waited=$((waited + 1))
    done

    local tmp
    tmp="$(mktemp "$myday_dir/.claims-tmp.XXXXXX")"
    # Seed tmp with existing ledger or empty array; never create a partial ledger_file first
    if [ -f "$ledger_file" ]; then
        jq --arg id "$id" --argjson done "$done_val" --arg owner "$owner" \
            'if any(.[]; .id == $id) then
                 map(if .id == $id then (.done = $done | (if ((.owner // "") == "") then .owner = $owner else . end)) else . end)
             else
                 . + [{"id": $id, "done": $done, "owner": $owner}]
             end' \
            "$ledger_file" > "$tmp" && mv "$tmp" "$ledger_file"
    else
        jq -n --arg id "$id" --argjson done "$done_val" --arg owner "$owner" \
            '[{"id": $id, "done": $done, "owner": $owner}]' > "$tmp" && mv "$tmp" "$ledger_file"
    fi

    [ -n "$acquired" ] && rmdir "$lock" 2>/dev/null
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
    json="$($A1 project workitem get "$id" -f json 2>/dev/null)" || return 1
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
        $A1 project workitem update "$id" --tag "$new" "$@"
    else
        echo "claim.sh: warning: could not read existing tags for $id; writing only '$add' (pre-existing tags may be lost)" >&2
        $A1 project workitem update "$id" --tag "$add" "$@"
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
        # Tag as done, preserving other tags: existing − {claimed, idle} ∪ {done}.
        # Tag FIRST and status SEPARATELY: previously tag+status went in one a1 update, so
        # an unsupported status (project enum mismatch) failed the whole call and left the
        # workitem stuck as idle with no done tag. Now the done tag always lands; a rejected
        # status is a WARN, not fatal.
        eff_status="$(_pool_done_status "$project_id")"
        _update_tags_merged "$workitem_id" "$DONE_TAG" "$CLAIM_TAG,$IDLE_TAG"
        rm -f "$(_claim_prefix_path "$workitem_id")"
        if [ -n "$eff_status" ]; then
            if $A1 project workitem update "$workitem_id" --status "$eff_status" >/dev/null 2>&1; then
                echo "claim.sh: finished workitem $workitem_id in project $project_id (status=$eff_status)"
            else
                echo "claim.sh: warning: finished (tagged $DONE_TAG) but status '$eff_status' was rejected for $workitem_id — set a valid status manually or add a per-pool done_status in pools.json" >&2
                echo "claim.sh: finished workitem $workitem_id in project $project_id (status unchanged)"
            fi
        else
            echo "claim.sh: finished workitem $workitem_id in project $project_id (no done_status configured)"
        fi
        _ledger_upsert "$workitem_id" true
        exit 0
        ;;

    *)
        echo "claim.sh: unknown command '$cmd'. Use claim, release, or finish." >&2
        exit 1
        ;;
esac
