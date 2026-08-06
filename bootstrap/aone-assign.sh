#!/usr/bin/env bash
# bootstrap/aone-assign.sh — the ONLY writer for an Aone workitem's assignee.
#
# Why this exists: nothing in this repo used to write `--assignee`. The route
# phase told the model to "幂等同步源单 assignee" and the model called a1
# directly, so the rule was unconditional overwrite with no notion of a human
# having taken the ticket over. Observed damage (2026-07/08, tf_customer 1086837):
#   #84486902  辰羿 took it 07-28 → jarvis reassigned to 过载 07-31 → 辰羿 took it
#              back 08-05 → jarvis re-claimed 16 minutes later
#   #84955165  jarvis 辰羿→新山 08-05 14:24 → 辰羿 took it back 23:48
#   #85115148  jarvis 过载→新山 17:18 → jarvis 新山→过载 19:10, no human involved
#              at all (the urgent/non-urgent classification is not stable across
#              runs, so the digital worker flip-flopped two humans by itself)
#
# Policy: a ticket already held by a real human is not reassigned. The first
# write wins and later overwrites are refused, which also kills the
# machine-vs-machine flapping above. Agent ids (WORKER_*) and people outside
# both rosters stay reassignable, so initial routing is unaffected.
#
# Protected = EITHER of two rosters in config/contacts.json:
#   1. `contacts` — the API tool team: id is not WORKER_* and the entry is not
#      legacy_inbound_only (those are stale historical owners that must not keep
#      holding new work — see contacts.json agent_fallbacks_desc).
#   2. `product_maintainers` — the 云产品专属维护名单 behind tf_customer 分支 A.
#      These product-team owners are deliberately absent from `contacts`
#      (that array doubles as the DingTalk delegation whitelist), so roster 1
#      alone left every branch-A owner unprotected. Observed damage:
#        #84363256  ACK 死分支单 — open-jarvis routed it to 若即 (ACK 专属维护人)
#                   on 07-16 per 分支 A; on 08-06 16:49 the RD reclassified it as
#                   「D 手写非紧急」 and reassigned 若即 → 过载, then DM'd 过载
#                   「已按手写 resource 非紧急路由给你」. Branch A outranks D in the
#                   decision tree, so this was a pure routing regression that the
#                   contacts-only check could not see.
#
# Usage: aone-assign.sh [--check] <workitem-id> <staff-id>
# Exit:  0  written, or already the target (idempotent no-op)
#        1  usage error
#        3  refused (human owner, or ownership could not be verified)
#
# --check runs the identical decision and reports it without writing, so callers
# that need "would this reassignment be allowed?" cannot drift from the policy.
# bridge/terraform_route_notify.py uses it to skip a route DM whose recipient is
# not the protected owner — the DM recipient is derived from the subtype, not
# from the ticket, so an unskipped DM tells the wrong person to start work.
#
# Fail-closed: if the current assignee cannot be read, this refuses rather than
# reassigning blind — an unverifiable owner is treated as a human owner.
#
# Escape hatch: JARVIS_ASSIGN_OK=1 skips the protection check, mirroring
# JARVIS_MASTER_OK=1 in worktree-guard.sh. Only for a per-turn instruction from
# the repo owner; never wire it into automation.
#
# Enforcement: bootstrap/a1_command_guard.py denies a raw `project workitem
# update ... --assignee` from the Bash tool and points here. The a1 call this
# script makes is a child process, which PreToolUse does not see.
set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$script_dir/lib.sh"
jarvis_root="$(jarvis_root)"
contacts_cfg="${JARVIS_CONTACTS:-$jarvis_root/config/contacts.json}"
A1="${JARVIS_A1:-${A1_BIN:-$jarvis_root/bin/a1id --}}"

check_only=""
if [ "${1:-}" = "--check" ]; then
    check_only="1"
    shift
fi

workitem_id="${1:-}"
target="${2:-}"

if [ -z "$workitem_id" ] || [ -z "$target" ]; then
    echo "Usage: aone-assign.sh [--check] <workitem-id> <staff-id>" >&2
    exit 1
fi
case "$workitem_id" in ""|*[![:digit:]]*)
    echo "aone-assign.sh: workitem id must be numeric, got '$workitem_id'" >&2
    exit 1 ;;
esac
case "$target" in ""|*[![:alnum:]_]*)
    echo "aone-assign.sh: staff id must be alphanumeric/underscore, got '$target'" >&2
    exit 1 ;;
esac

refuse() {
    echo "aone-assign.sh: refusing to reassign #$workitem_id — $1" >&2
    echo "aone-assign.sh: 已被真人接管的工单不改派（CLAUDE.md 工作纪律 #12）——含 API 团队在册真人与云产品专属维护名单。" >&2
    echo "aone-assign.sh: 如确需改派，请人工在 Aone 上操作，或由仓库主人当轮授权后设 JARVIS_ASSIGN_OK=1。" >&2
    exit 3
}

# --- current owner (fresh read; the 3h aone-get cache would hide a takeover) --
detail="$(JARVIS_CACHE_TTL=0 bash "$script_dir/aone-get.sh" "$workitem_id" 2>/dev/null)" || detail=""
if [ -z "$detail" ]; then
    refuse "could not read the work item (fail-closed)"
fi

current="$(printf '%s' "$detail" | jq -r '
    (.fields // [])
    | map(select(.identifier == "assignedTo"))
    | (.[0].value // "")
' 2>/dev/null)" || current=""
current_name="$(printf '%s' "$detail" | jq -r '
    (.fields // [])
    | map(select(.identifier == "assignedTo"))
    | (.[0].displayValue // "")
' 2>/dev/null)" || current_name=""

if [ -z "$current" ] || [ "$current" = "null" ]; then
    refuse "could not determine the current assignee (fail-closed)"
fi

# --- idempotent no-op --------------------------------------------------------
# The route phase re-runs every dispatch; converging on the same owner must not
# be an error, otherwise every revisit of a correctly-routed ticket fails. This
# is also what makes --check usable as the route-DM gate: the recipient already
# owning the ticket is exactly the case where the DM is addressed correctly.
if [ "$current" = "$target" ]; then
    echo "aone-assign.sh: #$workitem_id already assigned to $target${current_name:+ ($current_name)}; no change"
    exit 0
fi

# --- protection --------------------------------------------------------------
if [ "${JARVIS_ASSIGN_OK:-}" = "1" ]; then
    echo "aone-assign.sh: JARVIS_ASSIGN_OK=1 — bypassing human-owner protection for #$workitem_id" >&2
else
    if [ ! -f "$contacts_cfg" ]; then
        refuse "contacts roster not found at $contacts_cfg (fail-closed)"
    fi
    # One pass over both rosters. A missing/emptied roster array is treated as
    # damage, not as "nobody is protected" — silently degrading would re-open
    # exactly the overwrite loop this script exists to stop.
    protected="$(jq -r --arg id "$current" '
        def roster($key): (.[$key] // null);
        def usable($key): (roster($key) | type) == "array" and (roster($key) | length) > 0;
        if (usable("contacts") and usable("product_maintainers")) | not then
            "roster-damaged"
        else
            ((.product_maintainers | map(select(
                (.id == $id) and ((.id | startswith("WORKER_")) | not)
              ))) as $product
            | (.contacts | map(select(
                (.id == $id)
                and ((.id | startswith("WORKER_")) | not)
                and (.legacy_inbound_only != true)
              ))) as $team
            | if ($product | length) > 0
              then "product\t" + (($product[0].product // "") | tostring)
              elif ($team | length) > 0 then "team"
              else "open" end)
        end
    ' "$contacts_cfg" 2>/dev/null)" || protected=""
    case "$protected" in
        "")
            refuse "could not evaluate the contacts roster (fail-closed)" ;;
        roster-damaged)
            refuse "contacts/product_maintainers roster missing or empty in $contacts_cfg (fail-closed)" ;;
        product*)
            product="${protected#product}"
            product="${product#$'\t'}"
            refuse "已由云产品专属维护人 ${current_name:-$current}($current) 持单（专属维护名单${product:+：$product}，tf_customer 分支 A 优先于 D/G/pure datasource），拒绝改派到 $target" ;;
        team)
            refuse "已由真人 ${current_name:-$current}($current) 跟进中，拒绝改派到 $target" ;;
        open)
            : ;;
        *)
            refuse "unexpected roster verdict '$protected' (fail-closed)" ;;
    esac
fi

# --- check mode: same decision, no write -------------------------------------
if [ -n "$check_only" ]; then
    echo "aone-assign.sh: --check #$workitem_id ${current_name:-$current}($current) → $target would be allowed"
    exit 0
fi

# --- write -------------------------------------------------------------------
if ! $A1 project workitem update "$workitem_id" --assignee "$target" >/dev/null 2>&1; then
    echo "aone-assign.sh: a1 assignee update failed for #$workitem_id → $target" >&2
    exit 1
fi
bash "$script_dir/cache.sh" bust "wi-$workitem_id" >/dev/null 2>&1 || true
echo "aone-assign.sh: #$workitem_id assignee ${current_name:-$current}($current) → $target"
exit 0
