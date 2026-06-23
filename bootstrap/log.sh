#!/bin/bash
# bootstrap/log.sh — audit + dedup ledger for jarvis triage-loop
#
# Provides three functions:
#   run_done <id> <summary>  — write runs/<UTCdate>-<id>.md
#   escalate <id> <reason>   — write escalation/<id>.md
#   seen <id>                — exit 0 if runs/ already has a file for id, else exit 1
#
# Repo-root-relative paths (via git rev-parse).
# Override paths via env for testing:
#   JARVIS_RUNS_DIR       — defaults to <repo_root>/runs
#   JARVIS_ESCALATION_DIR — defaults to <repo_root>/escalation
#
# Guard: callable both sourced and as direct script (bash log.sh <fn> <args>).

# ---------------------------------------------------------------------------
# Path resolution
# ---------------------------------------------------------------------------
_log_repo_root() {
    git -C "$(dirname "${BASH_SOURCE[0]}")" rev-parse --show-toplevel 2>/dev/null \
        || git rev-parse --show-toplevel 2>/dev/null
}

_log_runs_dir() {
    echo "${JARVIS_RUNS_DIR:-$(_log_repo_root)/runs}"
}

_log_escalation_dir() {
    echo "${JARVIS_ESCALATION_DIR:-$(_log_repo_root)/escalation}"
}

# ---------------------------------------------------------------------------
# run_done <id> <summary>
# Write runs/<UTCdate>-<id>.md with id, summary, and timestamp.
# ---------------------------------------------------------------------------
run_done() {
    local id="$1"
    local summary="$2"
    local runs_dir
    runs_dir="$(_log_runs_dir)"
    local utc_date
    utc_date="$(date -u +%F)"
    local filepath="$runs_dir/${utc_date}-${id}.md"

    mkdir -p "$runs_dir"
    cat > "$filepath" <<EOF
# Run record: $id

**id:** $id
**summary:** $summary
**timestamp:** $(date -u +%Y-%m-%dT%H:%M:%SZ)
EOF
}

# ---------------------------------------------------------------------------
# escalate <id> <reason>
# Write escalation/<id>.md with id and reason.
# ---------------------------------------------------------------------------
escalate() {
    local id="$1"
    local reason="$2"
    local esc_dir
    esc_dir="$(_log_escalation_dir)"
    local filepath="$esc_dir/${id}.md"

    mkdir -p "$esc_dir"
    cat > "$filepath" <<EOF
# Escalation: $id

**id:** $id
**reason:** $reason
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
    runs_dir="$(_log_runs_dir)"
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
            run_done "${2:-}" "${3:-}"
            ;;
        escalate)
            escalate "${2:-}" "${3:-}"
            ;;
        seen)
            seen "${2:-}"
            ;;
        *)
            echo "Usage: log.sh {run_done|escalate|seen} [args...]" >&2
            exit 1
            ;;
    esac
fi
