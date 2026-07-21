#!/bin/bash
# bootstrap/log.sh — audit + dedup ledger for jarvis triage-loop
#
# Provides two functions:
#   run_done <id> <summary> [state]  — write runs/<UTCdate>-<id>.md
#                                       state ∈ {pending,merged} default pending
#   seen <id>                — exit 0 if runs/ already has a file for id, else exit 1
#
# Repo-root-relative paths (via git rev-parse).
# Override paths via env for testing:
#   JARVIS_RUNS_DIR       — defaults to <repo_root>/runs
#
# Guard: callable both sourced and as direct script (bash log.sh <fn> <args>).

# ---------------------------------------------------------------------------
# Path resolution (P1.e:走 lib.sh helper 消除 _log_repo_root/_log_*_dir 三样板)
# ---------------------------------------------------------------------------
# shellcheck source=bootstrap/lib.sh
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

# ---------------------------------------------------------------------------
# run_done <id> <summary> [state]
# Write runs/<UTCdate>-<id>.md with id, summary, state, and timestamp.
# state defaults to "pending" (审核中); "merged" = 已完成. 2-arg callers unbroken.
# ---------------------------------------------------------------------------
run_done() {
    local id="$1"
    local summary="$2"
    local state="${3:-}"; state="${state:-pending}"
    local runs_dir
    runs_dir="$(lib_runs_dir)"
    local utc_date
    utc_date="$(date -u +%F)"
    local filepath="$runs_dir/${utc_date}-${id}.md"

    mkdir -p "$runs_dir"
    cat > "$filepath" <<EOF
# Run record: $id

**id:** $id
**summary:** $summary
**state:** $state
**timestamp:** $(date -u +%Y-%m-%dT%H:%M:%SZ)
EOF
}

# ---------------------------------------------------------------------------
# seen <id>
# Exit 0 if any runs/ file for this id exists, else exit 1.
# ---------------------------------------------------------------------------
seen() {
    local id="$1"
    local runs_dir
    runs_dir="$(lib_runs_dir)"
    # Use glob; if any file matches *-<id>.md, return success
    local matches
    matches=$(ls "$runs_dir"/*-"${id}".md 2>/dev/null | wc -l | tr -d ' ')
    if [ "$matches" -ge 1 ]; then
        return 0
    else
        return 1
    fi
}

# ---------------------------------------------------------------------------
# Direct dispatch guard: if run as script (not sourced), dispatch by $1
# ---------------------------------------------------------------------------
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    cmd="${1:-}"
    case "$cmd" in
        run_done)
            run_done "${2:-}" "${3:-}" "${4:-}"
            ;;
        seen)
            seen "${2:-}"
            ;;
        *)
            echo "Usage: log.sh {run_done|seen} [args...]" >&2
            exit 1
            ;;
    esac
fi
