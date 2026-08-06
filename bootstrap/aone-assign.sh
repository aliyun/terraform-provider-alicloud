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
# Policy: a ticket already held by an active API-team human is not reassigned.
# The first write wins and later overwrites are refused, which also kills the
# machine-vs-machine flapping above. Agent ids (WORKER_*) and people outside the
# team roster stay reassignable, so initial routing is unaffected.
#
# Protected = present in config/contacts.json AND id is not WORKER_* AND the
# entry is not legacy_inbound_only (those are stale historical owners that must
# not keep holding new work — see contacts.json agent_fallbacks_desc).
#
# Usage: aone-assign.sh <workitem-id> <staff-id>
# Exit:  0  written, or already the target (idempotent no-op)
#        1  usage error
#        3  refused (human owner, or ownership could not be verified)
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

workitem_id="${1:-}"
target="${2:-}"

if [ -z "$workitem_id" ] || [ -z "$target" ]; then
    echo "Usage: aone-assign.sh <workitem-id> <staff-id>" >&2
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
    echo "aone-assign.sh: 已被 API 团队真人接管的工单不改派（CLAUDE.md 工作纪律 #10/#11）。" >&2
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
# be an error, otherwise every revisit of a correctly-routed ticket fails.
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
    protected="$(jq -r --arg id "$current" '
        (.contacts // [])
        | map(select(
            (.id == $id)
            and ((.id | startswith("WORKER_")) | not)
            and (.legacy_inbound_only != true)
          ))
        | length
    ' "$contacts_cfg" 2>/dev/null)" || protected=""
    if [ -z "$protected" ]; then
        refuse "could not evaluate the contacts roster (fail-closed)"
    fi
    if [ "$protected" != "0" ]; then
        refuse "已由真人 ${current_name:-$current}($current) 跟进中，拒绝改派到 $target"
    fi
fi

# --- write -------------------------------------------------------------------
if ! $A1 project workitem update "$workitem_id" --assignee "$target" >/dev/null 2>&1; then
    echo "aone-assign.sh: a1 assignee update failed for #$workitem_id → $target" >&2
    exit 1
fi
bash "$script_dir/cache.sh" bust "wi-$workitem_id" >/dev/null 2>&1 || true
echo "aone-assign.sh: #$workitem_id assignee ${current_name:-$current}($current) → $target"
exit 0
