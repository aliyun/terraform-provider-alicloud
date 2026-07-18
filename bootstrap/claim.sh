#!/usr/bin/env bash
# bootstrap/claim.sh – claim/release/finish primitive for concurrent multi-machine agents.
# Usage:
#   claim.sh claim   <workitem_id> <project_id>
#   claim.sh release <workitem_id> <project_id>
#   claim.sh finish  <workitem_id> <project_id> [status_override]
#
# claim:   tags the workitem jarvis-claimed, freezes a `jarvis-claim <host> <utc>` prefix
#          into .my-day/claim-prefix-<id>.txt (consumed by wrap.sh sync/done into the
#          first business comment to avoid an independent claim-only comment), then
#          point-reads the workitem to verify our claimed tag landed. Exits 0 on success,
#          1 if the claimed tag is not visible in the readback (lost race; prefix cleaned),
#          3 if Aone rejects the tag update because a required field is empty. Other tag
#          update failures are returned immediately and never enter the lost-race readback.
#          On a successful claim it ALSO advances the Aone status to the pool's in-progress
#          status (.progress_status, e.g. 处理中/开发中/问题解决中/Open) — but ONLY when the
#          ticket's CURRENT status is a starting state (.claim.start_statuses), so a
#          mid/late-flight ticket is never dragged backward and the write is idempotent.
#          This is best-effort/non-fatal: a rejected status is a WARN and never affects
#          claim's exit code or the lock semantics (tag + ledger stay the sole truth).
#          Set JARVIS_CLAIM_PROGRESS=0 to disable the status advance entirely.
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
# .claim.done_status / .claim.progress_status / .claim.start_statuses). Respects JARVIS_ROOT
# env override for the repo root.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
source "$script_dir/post-pr-context.sh"
jarvis_root="$(jarvis_root)"
pools_cfg="$jarvis_root/config/pools.json"

if jarvis_post_pr_context_active; then
    echo "claim.sh: post-PR 子代理不得认领、释放或完成 Aone 工单；bookend 由 bridge 在模型进程外执行" >&2
    exit 73
fi

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
PROGRESS_STATUS=$(jq -r '.claim.progress_status // empty' "$pools_cfg")

# Resolve the done status for a project: a per-pool .done_status (matched by project id)
# overrides the global .claim.done_status. Different Aone projects have different status
# enums (e.g. 2100304 only allows 处理中/已完成/已取消, not the global "已发布待需求排期"),
# so finish must use the target project's own completion status. Echoes "" if neither set.
#
# The per-pool .done_status may be EITHER a string (one completion status for the whole
# pool) OR an object keyed by workitem type displayValue (e.g. {"产品类需求":"已发布","功能缺陷":
# "Fixed"}) — because a single Aone project's status enum differs by workitem category
# (产品类需求 完成态=已发布, but 功能缺陷 枚举无「已发布」, 完成态=Fixed; 需求问题 完成态=已发布待需求方验收). When the
# per-pool value is an object, select by <wtype> (2nd arg); an unknown/empty wtype falls
# back to the object's first value to avoid an empty status. When it is a (non-empty)
# string, return it verbatim (backward compatible). If no per-pool value, fall back to the
# global .claim.done_status (always a string). Echoes "" if neither set.
_pool_done_status() {
    local project="$1" wtype="${2:-}" ps
    ps="$(jq -r --arg p "$project" --arg w "$wtype" \
        '(.pools[]? | select((.project|tostring) == $p) | .done_status)
         | if . == null then empty
           elif type == "object" then (.[$w] // (to_entries[0].value // empty))
           else . end' \
        "$pools_cfg" 2>/dev/null)"
    if [ -n "$ps" ]; then echo "$ps"; else jq -r '.claim.done_status // empty' "$pools_cfg" 2>/dev/null; fi
}

# Resolve the in-progress status for a project — an exact clone of _pool_done_status but
# resolving .progress_status per-pool and falling back to the global .claim.progress_status.
# Same string/object-by-workitemType/first-value-fallback semantics: different Aone projects
# name their "in progress" state differently (处理中 / 开发中 / 问题解决中 / Open / In Progress),
# and within a single project the enum differs by workitem category (需求 vs 功能缺陷), so a
# per-pool .progress_status may be EITHER a string OR an object keyed by workitem type
# displayValue. Object → select by <wtype> (2nd arg), unknown/empty wtype falls back to the
# object's first value. Non-empty string → verbatim. No per-pool value → global
# .claim.progress_status. Echoes "" if neither set.
_pool_progress_status() {
    local project="$1" wtype="${2:-}" ps
    ps="$(jq -r --arg p "$project" --arg w "$wtype" \
        '(.pools[]? | select((.project|tostring) == $p) | .progress_status)
         | if . == null then empty
           elif type == "object" then (.[$w] // (to_entries[0].value // empty))
           else . end' \
        "$pools_cfg" 2>/dev/null)"
    if [ -n "$ps" ]; then echo "$ps"; else jq -r '.claim.progress_status // empty' "$pools_cfg" 2>/dev/null; fi
}

# Point-read a workitem's type displayValue (e.g. 需求 / 功能缺陷 / 任务), used to pick the
# per-category done_status. The top-level categoryIdentifier/workitemType are null; the
# real value lives in fields[] under identifier=="workitemType". Echoes "" on any failure
# (get call fails / field absent) so callers degrade gracefully.
_get_wtype() {
    local id="$1" json
    json="$($A1 project workitem get "$id" -f json 2>/dev/null)" || return 0
    printf '%s' "$json" | jq -r '
        (.fields // [])[] | select(.identifier=="workitemType") | .displayValue // empty
    ' 2>/dev/null || return 0
}

# Point-read a workitem's CURRENT status displayValue (e.g. 待处理 / 处理中 / 已发布), used to
# decide whether a claim may advance the status (only from a starting state). Same fields[]
# shape as _get_wtype: identifier=="status". Echoes "" on any failure (get call fails / field
# absent) so callers degrade gracefully — an empty status means "can't verify starting state",
# and _advance_status then skips the write rather than risk a backward move.
_get_status() {
    local id="$1" json
    json="$($A1 project workitem get "$id" -f json 2>/dev/null)" || return 0
    printf '%s' "$json" | jq -r '
        (.fields // [])[] | select(.identifier=="status") | .displayValue // empty
    ' 2>/dev/null || return 0
}

# True (exit 0) if <cur> is a member of .claim.start_statuses, else 1. The start set (待处理/
# 新建/New/待认领/Reopen…) is the only set from which claim advances the status: it excludes
# every in-progress/mid/late/terminal state, so a ticket already moving is never dragged back,
# and — since the target progress_status is never itself a start status — the advance is
# idempotent across repeated claims.
_in_start_statuses() {
    local cur="$1"
    jq -e -n --arg cur "$cur" --slurpfile cfg "$pools_cfg" \
        '($cfg[0].claim.start_statuses // []) | index($cur) != null' >/dev/null 2>&1
}

# True (exit 0) if <cur> is a member of .claim.done_statuses — the set of terminal statuses
# (已发布/已完成/已拒绝…) at which work is already concluded. finish uses this to skip the
# status write when a ticket is ALREADY terminal: forcing done_status onto an already-done
# ticket is at best a no-op WARN and at worst a backward move (e.g. dragging 已完成/已拒绝
# back to 已发布). Empty list → never matches (preserves old always-write behavior).
_in_done_statuses() {
    local cur="$1"
    jq -e -n --arg cur "$cur" --slurpfile cfg "$pools_cfg" \
        '($cfg[0].claim.done_statuses // []) | index($cur) != null' >/dev/null 2>&1
}

# Best-effort, non-fatal advance of a claimed workitem's Aone status to the pool's in-progress
# status. Mirrors finish's tag/status separation: this NEVER exits non-zero or otherwise fails
# the claim — the lock's truth is the jarvis-claimed tag + ledger, not the status write.
# Guards, in order (any → skip silently, no backward move, no noise):
#   - env gate JARVIS_CLAIM_PROGRESS (default 1); "0" disables the whole advance.
#   - empty current status → can't verify the starting state, don't risk a backward move.
#   - current status NOT in .claim.start_statuses → already in-progress/mid-flight/terminal.
#   - empty resolved progress_status → nothing configured for this pool/category.
# Only then write; a rejected status (enum mismatch / invalid transition) is a WARN, not fatal.
_advance_status() {
    local id="$1" project="$2" cur wtype prog
    [ "${JARVIS_CLAIM_PROGRESS:-1}" = "0" ] && return 0
    cur="$(_get_status "$id")"
    if [ -z "$cur" ]; then
        echo "claim.sh: note: could not read current status for $id; skipping status advance" >&2
        return 0
    fi
    _in_start_statuses "$cur" || return 0   # not a starting state → leave as-is (no backward move)
    wtype="$(_get_wtype "$id")"
    prog="$(_pool_progress_status "$project" "$wtype")"
    [ -n "$prog" ] || return 0
    if $A1 project workitem update "$id" --status "$prog" >/dev/null 2>&1; then
        echo "claim.sh: advanced $id status → $prog"
    else
        echo "claim.sh: warning: claimed but status '$prog' rejected for $id (enum mismatch? not a valid transition from '$cur') — set a valid per-pool progress_status in pools.json" >&2
    fi
    return 0
}

# Check for unmerged MRs/PRs associated with a workitem.
# Returns 0 if unmerged MRs found, 1 if none found.
# Checks: (1) active worktree branches containing the workitem ID,
#         (2) open GitHub PRs from api-tool-agent mentioning the ID.
_has_unmerged_mr() {
    local id="$1"

    # Check 1: active worktree branches containing this workitem ID.
    # The parent dir of jarvis_root holds all worktrees (e.g. ../jarvis-worktree-*).
    local parent_dir
    parent_dir="$(cd "$jarvis_root/.." && pwd)"
    if git -C "$parent_dir" worktree list --porcelain 2>/dev/null | grep -q "worktree-.*${id}"; then
        return 0
    fi

    # Check 2: open GitHub PRs from api-tool-agent mentioning this ID.
    if [ -n "${JARVIS_GITHUB_TOKEN:-}" ]; then
        local prs
        prs="$(GH_TOKEN="$JARVIS_GITHUB_TOKEN" gh search prs \
            --author api-tool-agent --state open "$id" \
            --json repository,number,title --limit 5 2>/dev/null || echo "[]")"
        if [ "$prs" != "[]" ] && [ -n "$prs" ]; then
            local count
            count="$(echo "$prs" | jq 'length' 2>/dev/null || echo 0)"
            if [ "$count" -gt 0 ] 2>/dev/null; then
                return 0
            fi
        fi
    fi

    return 1
}

cmd="${1:-}"
workitem_id="${2:-}"
project_id="${3:-}"
# finish only: optional 4th arg overrides the resolved per-pool done_status (e.g. PrWatch
# passes 已完成). Empty/absent → behavior unchanged (full backward compatibility).
STATUS_OVERRIDE="${4:-}"

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

# Best-effort point-read by itemId only. A title read failure must not block the
# existing claim flow; the interactive worker simply omits sourceRef.title.
_get_workitem_title() {
    local id="$1" json title
    json="$($A1 project workitem get "$id" -f json 2>/dev/null)" || return 1
    title="$(printf '%s' "$json" | jq -r '
        (.subject // .title //
         ((.fields // []) | map(select(.identifier=="subject" or .identifier=="title"))
          | (.[0].displayValue // "")))
        | tostring | gsub("^\\s+|\\s+$"; "")
    ' 2>/dev/null)" || return 1
    printf '%s' "$title"
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

# Return current tags as name<TAB>id pairs, one per line, TAB-separated. Empty tag
# → no output. Returns non-zero on get failure. Uses value (numeric IDs) so cross-
# project / whitelist-external tags (e.g. "企业级能力-资源化-OpenAPI(Terraform)"
# only exists at Aone-global level, not in project field options) survive downstream
# updates — a1 --tag validates names but accepts IDs, so preserving by ID never trips
# the "tag not found" gate.
_get_tag_pairs() {
    local id="$1" json
    json="$($A1 project workitem get "$id" -f json 2>/dev/null)" || return 1
    printf '%s' "$json" | jq -r '
        (.fields // []) | map(select(.identifier=="tag")) | .[0] |
        [(.displayValue // "" | split(", ")), (.value // "" | split(", "))]
        | if (.[0] == [""] or .[0] == []) then empty
          else (transpose | .[] | @tsv) end
    ' 2>/dev/null || return 1
}

# Run an a1 tag update while preserving its output/exit code. Aone rejects updates to some
# legacy tickets with an error shaped exactly like `【<name>】(<numeric-id>)不能为空`; classify
# only that structured validation error as rc=3 so orchestration can query legal options and
# repair the missing field. Looser messages containing merely `不能为空` are not enough: they
# may describe unrelated validation and must retain the original failure code.
_run_tag_update() {
    local out rc parsed field_id field_name
    if out="$("$@" 2>&1)"; then
        return 0
    else
        rc=$?
    fi

    parsed="$(printf '%s\n' "$out" | sed -n \
        's/.*【\([^】][^】]*\)】[[:space:]]*[(（]\([0-9][0-9]*\)[)）][[:space:]]*不能为空.*/\2\t\1/p' \
        | head -n 1)"
    if [ -n "$parsed" ]; then
        field_id="${parsed%%$'\t'*}"
        field_name="${parsed#*$'\t'}"
        printf '%s\n' "$out" >&2
        printf 'claim.sh: blocked missing_required_field %s %s\n' "$field_id" "$field_name" >&2
        return 3
    fi

    printf '%s\n' "$out" >&2
    return "$rc"
}

# Write the full expected tag set in one update, preserving pre-existing tags by ID
# (so cross-project / whitelist-external tags survive) and adding new tags by name.
# Args: <id> <add_csv> <remove_csv> [extra a1 update flags...].
# Degrades to a legacy bare `--tag <add>` write (may drop other tags) with a stderr
# warning if the current tags can't be read, so claim/release/finish never hard-fail
# on a transient get error.
_update_tags_merged() {
    local id="$1" add="$2" remove="$3"; shift 3
    local pairs kept_ids to_write
    if ! pairs="$(_get_tag_pairs "$id")"; then
        echo "claim.sh: warning: could not read existing tags for $id; writing only '$add' (pre-existing tags may be lost)" >&2
        _run_tag_update $A1 project workitem update "$id" --tag "$add" "$@"
        return
    fi
    # existing (name,id) pairs whose name NOT in (remove ∪ add) → keep by numeric ID
    kept_ids="$(printf '%s' "$pairs" | jq -rRs --arg rm "$remove" --arg add "$add" '
        def clean: split(",") | map(gsub("^\\s+|\\s+$";"")) | map(select(length>0));
        ($rm|clean) as $r | ($add|clean) as $a |
        split("\n") | map(select(length>0)) | map(split("\t")) |
        map(select(([.[0]]|inside($r+$a)|not))) |
        map(.[1]) | map(select(length>0)) | join(",")
    ' 2>/dev/null || true)"
    to_write="$kept_ids"
    [ -n "$add" ] && to_write="${to_write:+${to_write},}${add}"
    if [ -z "$to_write" ]; then
        # nothing to write — clearing all tags is not what callers expect
        return 0
    fi
    _run_tag_update $A1 project workitem update "$id" --tag "$to_write" "$@"
}

# Interactive Claude/Codex sessions use the database fence as the sole mutex.
# The wrapper loads the same gitignored control-plane environment as bridge/run.sh.
_is_interactive_context() {
    [ -n "${CLAUDE_CODE_SESSION_ID:-}" ] || [ -n "${CODEX_THREAD_ID:-}" ] || \
        { case "${JARVIS_INTERACTIVE_CLIENT:-}" in claude|codex) true ;; *) false ;; esac \
          && [ -n "${JARVIS_INTERACTIVE_SESSION_ID:-}" ]; }
}

_interactive_worker() {
    local runner="${JARVIS_INTERACTIVE_WORKER_RUNNER:-$script_dir/run-interactive-worker-hook.sh}"
    bash "$runner" cli "$@"
}

_interactive_fail_claim() {
    local id="$1" message="$2" unknown="${3:-0}"
    if [ "$unknown" = "1" ]; then
        _interactive_worker operation-fail "$id" "$message" --unknown >/dev/null 2>&1 || \
            echo "claim.sh: warning: failed to persist interactive claim failure for $id" >&2
    else
        _interactive_worker operation-fail "$id" "$message" >/dev/null 2>&1 || \
            echo "claim.sh: warning: failed to persist interactive claim failure for $id" >&2
    fi
}

# release/finish tag 回执的 UNKNOWN 收敛分支（docs/aone-operation-receipts.md）：
# 上轮 abort --unknown 后 operation-begin 会短路 needsReadback，此处 point-read 现
# 有 tags——命中目标 tag → reconcile --found（副作用已在，跳过重写）；未命中 →
# reconcile --not-found 后重新 begin 拿 proceed。stdout 输出 skip|write；point-read
# 不可达/回执调用失败返回非零（槽保持 UNKNOWN/RETRY_WAIT，fail-closed）。
_tag_receipt_readback() {
    local id="$1" kind="$2" material="$3" tag="$4" ref="$5"
    local tags again
    if ! tags="$(_get_tags "$id")"; then
        echo "claim.sh: $kind receipt readback unreachable for $id; receipt stays UNKNOWN" >&2
        return 2
    fi
    if _csv_contains "$tags" "$tag"; then
        _interactive_worker operation-reconcile "$id" --found "$ref" >/dev/null || return 2
        echo "skip"
        return 0
    fi
    _interactive_worker operation-reconcile "$id" --not-found >/dev/null || return 2
    again="$(_interactive_worker operation-begin "$id" "$kind" "$material" --replay-safe)" || return 2
    printf '%s' "$again" | jq -e '.proceed == true' >/dev/null 2>&1 || return 2
    echo "write"
}

case "$cmd" in
    claim)
        interactive_claim=0
        interactive_prepare=""
        if _is_interactive_context; then
            interactive_title="$(_get_workitem_title "$workitem_id" || true)"
            interactive_prepare="$(_interactive_worker prepare-claim "$workitem_id" "$project_id" "$interactive_title")"
            interactive_rc=$?
            if [ "$interactive_rc" -ne 0 ]; then
                if [ "$interactive_rc" -eq 10 ]; then
                    echo "claim.sh: workitem $workitem_id is already owned by another worker" >&2
                    exit 1
                fi
                echo "claim.sh: interactive control-plane claim failed closed for $workitem_id (rc=$interactive_rc)" >&2
                exit 2
            fi
            if ! printf '%s' "$interactive_prepare" | jq -e '
                (.accepted == true) and (.proceed == true or .operationStatus == "ACKED")
            ' >/dev/null 2>&1; then
                echo "claim.sh: invalid interactive claim receipt for $workitem_id" >&2
                exit 2
            fi
            interactive_claim=1
            if [ "$(printf '%s' "$interactive_prepare" | jq -r '.proceed')" != "true" ]; then
                # ACKED means a prior invocation already completed the exact Aone
                # effect. Rebuild only local bookends; never write Aone again.
                utcnow=$(date -u +%Y-%m-%dT%H:%M:%SZ)
                host="${JARVIS_SELF_HOST:-$(scutil --get ComputerName 2>/dev/null || hostname -s 2>/dev/null || hostname)}"
                mkdir -p "$myday_dir"
                [ -f "$(_claim_prefix_path "$workitem_id")" ] || \
                    printf 'jarvis-claim %s %s' "$host" "$utcnow" > "$(_claim_prefix_path "$workitem_id")"
                _advance_status "$workitem_id" "$project_id"
                _ledger_upsert "$workitem_id" false
                echo "claim.sh: interactive claim for workitem $workitem_id already acknowledged"
                exit 0
            fi
        fi

        # 1. Tag as claimed, preserving any pre-existing tags: existing ∪ {claimed} − {idle}
        _update_tags_merged "$workitem_id" "$CLAIM_TAG" "$IDLE_TAG"
        tag_rc=$?
        if [ "$tag_rc" -eq 3 ]; then
            # The tag never landed, so no prefix/ledger cleanup is needed. The caller must
            # inspect aone-fields.sh missing, choose a legal value, fill it, and retry claim.
            [ "$interactive_claim" = "1" ] && \
                _interactive_fail_claim "$workitem_id" "Aone rejected claim tag: missing required field"
            exit 3
        elif [ "$tag_rc" -ne 0 ]; then
            # Transport/auth/other validation failures are not races. Preserve the original
            # rc and stop before readback, which could otherwise mistake a stale tag for ours.
            # rc=1 is reserved by claim for an actual lost race, so remap an a1 rc=1 to rc=2.
            [ "$interactive_claim" = "1" ] && \
                _interactive_fail_claim "$workitem_id" "Aone claim tag update failed (rc=$tag_rc)"
            [ "$tag_rc" -eq 1 ] && tag_rc=2
            exit "$tag_rc"
        fi

        # 2. 冻结 hostname + utc 到 prefix 文件；wrap.sh sync/done 第一次发业务评论时自动 prefix
        # 到顶部并消费，避免 claim 时刻额外刷一条独立留痕评论。
        # 优先 ComputerName(用户自定义显示名)，回落 hostname -s 短名，最后 hostname；避免 mac 默认带 .local 后缀
        # JARVIS_SELF_HOST 覆盖(多机分布式统一标识/测试注入)；不设=现状 scutil 链。
        utcnow=$(date -u +%Y-%m-%dT%H:%M:%SZ)
        host="${JARVIS_SELF_HOST:-$(scutil --get ComputerName 2>/dev/null || hostname -s 2>/dev/null || hostname)}"
        mkdir -p "$myday_dir"
        printf 'jarvis-claim %s %s' "$host" "$utcnow" > "$(_claim_prefix_path "$workitem_id")"

        # 3. Readback via POINT-READ (workitem get), retries for a small safety margin.
        # The point-read confirms *our own* write landed; by itself it does NOT arbitrate a
        # true race — two instances claiming near-simultaneously both see jarvis-claimed.
        # When JARVIS_CLAIM_SETTLE>0 a cross-machine arbitration round (step 4) breaks the
        # tie via the reconcile timestamp範式; when it's 0 (default) behavior is unchanged
        # and mutual exclusion still rests on the ledger + reconcile stale/drift sweep.
        claimed_ok=""
        for _ in 1 2 3; do
            if _csv_contains "$(_get_tags "$workitem_id" 2>/dev/null || echo "")" "$CLAIM_TAG"; then
                claimed_ok=1; break
            fi
            sleep "${JARVIS_CLAIM_READBACK_SLEEP:-2}"
        done
        if [ -z "$claimed_ok" ]; then
            # 抢锁失败：清理本机 prefix 痕迹（不留给别人）
            rm -f "$(_claim_prefix_path "$workitem_id")"
            [ "$interactive_claim" = "1" ] && \
                _interactive_fail_claim "$workitem_id" "Aone claim tag readback was inconclusive" 1
            echo "claim.sh: lost race for workitem $workitem_id (claimed tag not visible in point-read readback)" >&2
            exit 1
        fi

        if [ "$interactive_claim" = "1" ]; then
            external_ref="aone:$project_id:$workitem_id:tag:$CLAIM_TAG"
            if ! _interactive_worker operation-ack "$workitem_id" "$external_ref" >/dev/null; then
                echo "claim.sh: Aone tag landed but control-plane ACK failed for $workitem_id; stopping fail-closed" >&2
                exit 2
            fi
            echo "claim.sh: claimed workitem $workitem_id in project $project_id (database fenced)"
            _advance_status "$workitem_id" "$project_id"
            _ledger_upsert "$workitem_id" false
            exit 0
        fi

        # 4. Cross-machine arbitration — ONLY when JARVIS_CLAIM_SETTLE>0. settle<=0 (default)
        # keeps the historical point-read-only semantics: no comment write, no settle wait.
        settle="${JARVIS_CLAIM_SETTLE:-0}"
        if ! [ "$settle" -gt 0 ] 2>/dev/null; then
            echo "claim.sh: claimed workitem $workitem_id in project $project_id"
            # Best-effort advance Aone status to the pool's in-progress status (non-fatal).
            _advance_status "$workitem_id" "$project_id"
            _ledger_upsert "$workitem_id" false
            exit 0
        fi

        # 4a. Write our claim comment. Format `jarvis-claim <host> <utc> <nonce>` keeps
        # reconcile.sh's existing regex working (nonce trails the timestamp, non-breaking).
        nonce="$(head -c8 /dev/urandom 2>/dev/null | od -An -tx1 2>/dev/null | tr -d ' \n' | cut -c1-8)"
        [ -n "$nonce" ] || nonce="$(printf '%08x' $((RANDOM * RANDOM)))"
        $A1 project workitem comment create "$workitem_id" -m "jarvis-claim $host $utcnow $nonce" >/dev/null 2>&1 \
            || echo "claim.sh: warning: arbitration comment write failed for $workitem_id" >&2

        # 4b. Settle window, then read the whole comment stream back.
        sleep "$settle"
        arb_comments="$($A1 project workitem comment list "$workitem_id" -f json 2>/dev/null || echo "[]")"

        # 4c. Total order over (epoch, nonce); the earliest claim wins. Echo the winner host.
        # Reuses the reconcile UTC→epoch idea (calendar.timegm) inline to avoid sourcing.
        winner_host="$(printf '%s' "$arb_comments" | python3 -c "
import json, sys, re, calendar, datetime
def epoch(ts):
    try:
        return calendar.timegm(datetime.datetime.strptime(ts, '%Y-%m-%dT%H:%M:%SZ').timetuple())
    except Exception:
        return None
pat = re.compile(r'jarvis-claim\s+(\S+)\s+(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)\s+(\S+)')
try:
    cs = json.loads(sys.stdin.read())
except Exception:
    cs = []
cands = []
for c in cs:
    content = c.get('content', '') if isinstance(c, dict) else ''
    m = pat.search(content)
    if not m:
        continue
    e = epoch(m.group(2))
    if e is None:
        continue
    cands.append((e, m.group(3), m.group(1)))   # (epoch, nonce, host)
if cands:
    cands.sort(key=lambda x: (x[0], x[1]))
    print(cands[0][2])
" 2>/dev/null || echo "")"

        # 4d. Verdict. Empty (only our own comment / parse miss) or winner==self → we win.
        # Otherwise an earlier claim from another host exists → stand down and roll back.
        if [ -z "$winner_host" ] || [ "$winner_host" = "$host" ]; then
            echo "claim.sh: claimed workitem $workitem_id in project $project_id (arbitration won)"
            # Best-effort advance Aone status to the pool's in-progress status (non-fatal).
            _advance_status "$workitem_id" "$project_id"
            _ledger_upsert "$workitem_id" false
            exit 0
        fi

        # Lost arbitration: an earlier claim from another host won. Do NOT touch the shared
        # tag set — the winner legitimately holds jarvis-claimed, and Aone `--tag` is a
        # whole-set replace with no per-owner identity, so writing idle here would CLOBBER
        # the winner's claim (and the winner, having already exited 0, would never restore
        # it). Our own step-1 claimed write is idempotent with the winner's (same tag name),
        # so leaving the tag set as-is is correct. Just drop our local prefix and stand down;
        # we never ledger-upserted, so there is nothing local to release.
        rm -f "$(_claim_prefix_path "$workitem_id")"
        echo "claim.sh: lost arbitration for workitem $workitem_id to host '$winner_host' (earlier claim); standing down, winner keeps $CLAIM_TAG" >&2
        exit 1
        ;;

    release)
        interactive_release=0
        if _is_interactive_context; then
            if ! _interactive_worker has-current "$workitem_id" >/dev/null 2>&1; then
                echo "claim.sh: refusing interactive release for $workitem_id without its database-fenced session" >&2
                exit 2
            fi
            interactive_release=1
        fi
        # fenced 会话的 idle tag 写走 operation receipt：begin 失败 fail-closed，不碰
        # Aone tag、session 保持（docs/aone-operation-receipts.md）。
        release_proceed=1
        release_ref="aone:$project_id:$workitem_id:tag:$IDLE_TAG"
        if [ "$interactive_release" = "1" ]; then
            if ! release_begin="$(_interactive_worker operation-begin "$workitem_id" release-tag idle --replay-safe)"; then
                echo "claim.sh: release receipt begin failed closed for $workitem_id; Aone untouched" >&2
                exit 2
            fi
            if ! printf '%s' "$release_begin" | jq -e '.accepted == true' >/dev/null 2>&1; then
                echo "claim.sh: invalid release receipt for $workitem_id" >&2
                exit 2
            fi
            if printf '%s' "$release_begin" | jq -e '.needsReadback == true' >/dev/null 2>&1; then
                # UNKNOWN 存量回执：point-read 收敛，永不盲重放。
                verdict="$(_tag_receipt_readback "$workitem_id" release-tag idle "$IDLE_TAG" "$release_ref")" || exit 2
                if [ "$verdict" = "skip" ]; then release_proceed=0; fi
            elif ! printf '%s' "$release_begin" | jq -e '.proceed == true' >/dev/null 2>&1; then
                release_proceed=0   # ACKED：上轮已落 idle tag，恰好一次，跳过重写
            fi
        fi
        # Tag as idle, preserving other tags: existing − {claimed} ∪ {idle}
        # 本轮 jarvis 处理完，释放锁；后续等待人或下一个 jarvis 接手（不动 Aone status）
        if [ "$release_proceed" = "1" ]; then
            _update_tags_merged "$workitem_id" "$IDLE_TAG" "$CLAIM_TAG"
            release_tag_rc=$?
        else
            release_tag_rc=0
        fi
        if [ "$interactive_release" = "1" ] && [ "$release_tag_rc" -ne 0 ]; then
            _interactive_worker operation-abort "$workitem_id" "Aone release tag update failed (rc=$release_tag_rc)" >/dev/null 2>&1 || \
                echo "claim.sh: warning: failed to abort release receipt for $workitem_id" >&2
            echo "claim.sh: Aone release write failed; keeping database session active for $workitem_id" >&2
            [ "$release_tag_rc" -eq 1 ] && release_tag_rc=2
            exit "$release_tag_rc"
        fi
        if [ "$interactive_release" = "1" ] && [ "$release_proceed" = "1" ]; then
            # readback 确认 idle tag 落地后才 ACK；point-read 不可达 → 冻结 UNKNOWN。
            if ! release_tags="$(_get_tags "$workitem_id")"; then
                _interactive_worker operation-abort "$workitem_id" "release tag readback unreachable" --unknown >/dev/null 2>&1 || true
                echo "claim.sh: release tag readback inconclusive for $workitem_id; receipt frozen UNKNOWN" >&2
                exit 2
            fi
            if ! _csv_contains "$release_tags" "$IDLE_TAG"; then
                _interactive_worker operation-abort "$workitem_id" "release idle tag not visible after write" >/dev/null 2>&1 || true
                echo "claim.sh: release idle tag not visible after write for $workitem_id" >&2
                exit 2
            fi
            if ! _interactive_worker operation-ack "$workitem_id" "$release_ref" >/dev/null; then
                echo "claim.sh: Aone idle tag landed but receipt ACK failed for $workitem_id; retry release" >&2
                exit 2
            fi
        fi
        if [ "$interactive_release" = "1" ]; then
            if ! _interactive_worker suspend "$workitem_id" "released by interactive worker" >/dev/null; then
                echo "claim.sh: Aone is idle but database session suspend failed for $workitem_id; retry release" >&2
                exit 2
            fi
        fi
        # 兜底：如果整个 jarvis 周期没发过业务评论，prefix 文件仍在，此处清理
        rm -f "$(_claim_prefix_path "$workitem_id")"
        echo "claim.sh: released workitem $workitem_id in project $project_id"
        _ledger_upsert "$workitem_id" true
        exit 0
        ;;

    finish)
        # Hard gate: refuse to finish if unmerged MRs exist (CLAUDE.md: MR 未合并禁 finish).
        # The workitem should be released (jarvis-idle) so humans can merge + verify.
        if [ "${JARVIS_SKIP_MR_GATE:-0}" != "1" ] && _has_unmerged_mr "$workitem_id"; then
            echo "claim.sh: refusing to finish $workitem_id — unmerged MR(s) detected" >&2
            echo "claim.sh: use 'claim.sh release' instead to hand off for merge review," >&2
            echo "claim.sh: or set JARVIS_SKIP_MR_GATE=1 to override this gate." >&2
            exit 2
        fi

        interactive_finish=0
        if _is_interactive_context; then
            if ! _interactive_worker has-current "$workitem_id" >/dev/null 2>&1; then
                echo "claim.sh: refusing interactive finish for $workitem_id without its database-fenced session" >&2
                exit 2
            fi
            interactive_finish=1
        fi

        # Tag as done, preserving other tags: existing − {claimed, idle} ∪ {done}.
        # Tag FIRST and status SEPARATELY: previously tag+status went in one a1 update, so
        # an unsupported status (project enum mismatch) failed the whole call and left the
        # workitem stuck as idle with no done tag. Tag lands first; then we try to move status.
        #
        # 一致性闸门(jarvis-done vs Aone 状态)：jarvis-done 标签是 jarvis 侧「干完了」的标记，
        # 但只有当 Aone status 真落到合法完成态(已终态 或 本池 done_status 写成功)时，标签才
        # 与真源一致。若状态没能落地(done_status 被 workflow 拒 / 无 done_status 配置且当前非终态)，
        # 继续挂 jarvis-done 会造成「标签说 done、真源说没完」的黑洞——_decide 从此永久 skip 这单，
        # 人工催办/reopen 都不再触发重派。故此处降级：撤 jarvis-done、改打 jarvis-idle(交 idle 门 +
        # RevisitScheduler 兜底重访/重派) + escalate 提示人工设正确状态或摘标签。
        wtype="$(_get_wtype "$workitem_id")"
        eff_status="$(_pool_done_status "$project_id" "$wtype")"
        # Optional caller override (4th arg): PrWatchScheduler passes 已完成 on PR merge.
        # already-terminal short-circuit + downgrade black-hole guard below stay intact.
        [ -n "$STATUS_OVERRIDE" ] && eff_status="$STATUS_OVERRIDE"
        # fenced 会话的 done tag 写走 operation receipt（同 release，kind=finish-tag）。
        finish_proceed=1
        finish_ref="aone:$project_id:$workitem_id:tag:$DONE_TAG"
        if [ "$interactive_finish" = "1" ]; then
            if ! finish_begin="$(_interactive_worker operation-begin "$workitem_id" finish-tag done --replay-safe)"; then
                echo "claim.sh: finish receipt begin failed closed for $workitem_id; Aone untouched" >&2
                exit 2
            fi
            if ! printf '%s' "$finish_begin" | jq -e '.accepted == true' >/dev/null 2>&1; then
                echo "claim.sh: invalid finish receipt for $workitem_id" >&2
                exit 2
            fi
            if printf '%s' "$finish_begin" | jq -e '.needsReadback == true' >/dev/null 2>&1; then
                verdict="$(_tag_receipt_readback "$workitem_id" finish-tag done "$DONE_TAG" "$finish_ref")" || exit 2
                if [ "$verdict" = "skip" ]; then finish_proceed=0; fi
            elif ! printf '%s' "$finish_begin" | jq -e '.proceed == true' >/dev/null 2>&1; then
                finish_proceed=0   # ACKED：上轮已落 done tag，恰好一次，跳过重写
            fi
        fi
        if [ "$finish_proceed" = "1" ]; then
            _update_tags_merged "$workitem_id" "$DONE_TAG" "$CLAIM_TAG,$IDLE_TAG"
            finish_tag_rc=$?
        else
            finish_tag_rc=0
        fi
        if [ "$interactive_finish" = "1" ] && [ "$finish_tag_rc" -ne 0 ]; then
            _interactive_worker operation-abort "$workitem_id" "Aone finish tag update failed (rc=$finish_tag_rc)" >/dev/null 2>&1 || \
                echo "claim.sh: warning: failed to abort finish receipt for $workitem_id" >&2
            echo "claim.sh: Aone finish tag write failed; keeping database session active for $workitem_id" >&2
            [ "$finish_tag_rc" -eq 1 ] && finish_tag_rc=2
            exit "$finish_tag_rc"
        fi
        if [ "$interactive_finish" = "1" ] && [ "$finish_proceed" = "1" ]; then
            if ! finish_tags="$(_get_tags "$workitem_id")"; then
                _interactive_worker operation-abort "$workitem_id" "finish tag readback unreachable" --unknown >/dev/null 2>&1 || true
                echo "claim.sh: finish tag readback inconclusive for $workitem_id; receipt frozen UNKNOWN" >&2
                exit 2
            fi
            if ! _csv_contains "$finish_tags" "$DONE_TAG"; then
                _interactive_worker operation-abort "$workitem_id" "finish done tag not visible after write" >/dev/null 2>&1 || true
                echo "claim.sh: finish done tag not visible after write for $workitem_id" >&2
                exit 2
            fi
            if ! _interactive_worker operation-ack "$workitem_id" "$finish_ref" >/dev/null; then
                echo "claim.sh: Aone done tag landed but receipt ACK failed for $workitem_id; retry finish" >&2
                exit 2
            fi
        fi
        rm -f "$(_claim_prefix_path "$workitem_id")"
        # finish 的 status 写与下方降级路径**不**叠 operation receipt：status 落地失败
        # 已有 jarvis-idle 降级 + escalate 兜底（不留黑洞），降级 tag 写维持 best-effort，
        # 再叠一层回执只会把可自愈路径变成 fail-closed 死路（docs/aone-operation-receipts.md）。
        cur_status="$(_get_status "$workitem_id")"
        status_ok=0   # 1 = Aone status 已落到合法完成态，jarvis-done 与真源一致
        if [ -n "$cur_status" ] && { _in_done_statuses "$cur_status" || [ "$cur_status" = "$eff_status" ]; }; then
            # Already at a terminal status — 全局 .claim.done_statuses(已发布/已完成/已拒绝…)
            # 或 per-池 done_status(如 tf_provider 产品类需求"待发布")。写 done_status 会 no-op-WARN
            # 或把 已完成/已拒绝 拖回 已发布。与 donecheck 的合法完成态并集判据一致
            # (loops/aone-triage.md §三:「.claim.done_statuses ∪ 各池 .done_status,含 tf_provider 待发布」)。
            echo "claim.sh: finished workitem $workitem_id in project $project_id (already terminal: $cur_status)"
            status_ok=1
        elif [ -n "$eff_status" ]; then
            if $A1 project workitem update "$workitem_id" --status "$eff_status" >/dev/null 2>&1; then
                echo "claim.sh: finished workitem $workitem_id in project $project_id (status=$eff_status)"
                status_ok=1
            else
                echo "claim.sh: warning: status '$eff_status' was rejected for $workitem_id (project enum mismatch / workflow disallows transition)" >&2
            fi
        else
            echo "claim.sh: warning: no done_status configured for $workitem_id (project $project_id, type '${wtype:-?}'); current status '${cur_status:-?}' is not terminal" >&2
        fi

        if [ "$status_ok" != "1" ]; then
            # 状态未落到合法完成态 → 不留 jarvis-done 黑洞：降级为 jarvis-idle + escalate。
            _update_tags_merged "$workitem_id" "$IDLE_TAG" "$DONE_TAG"
            downgrade_tag_rc=$?
            if [ "$interactive_finish" = "1" ] && [ "$downgrade_tag_rc" -ne 0 ]; then
                echo "claim.sh: failed to downgrade Aone to idle; keeping database session active for $workitem_id" >&2
                [ "$downgrade_tag_rc" -eq 1 ] && downgrade_tag_rc=2
                exit "$downgrade_tag_rc"
            fi
            bash "$script_dir/log.sh" escalate "$workitem_id" \
                "finish_status_unresolved: done_status='${eff_status:-<none>}' 未落地(当前='${cur_status:-<unknown>}',workitemType='${wtype:-?}')；已降级 jarvis-done→jarvis-idle，请人工设正确完成态或摘标签" \
                "" >/dev/null 2>&1 || true
            echo "claim.sh: finish $workitem_id — 状态未落合法完成态，已降级 jarvis-done→jarvis-idle 并 escalate（交 Revisit 兜底）" >&2
        fi
        if [ "$interactive_finish" = "1" ]; then
            if [ "$status_ok" = "1" ]; then
                if ! _interactive_worker complete "$workitem_id" "Aone reached terminal state" >/dev/null; then
                    echo "claim.sh: Aone is complete but database session completion failed for $workitem_id; retry finish" >&2
                    exit 2
                fi
            elif ! _interactive_worker suspend "$workitem_id" "finish downgraded to jarvis-idle" >/dev/null; then
                echo "claim.sh: finish downgraded to idle but database session suspend failed for $workitem_id" >&2
                exit 2
            fi
        fi
        _ledger_upsert "$workitem_id" true
        exit 0
        ;;

    *)
        echo "claim.sh: unknown command '$cmd'. Use claim, release, or finish." >&2
        exit 1
        ;;
esac
