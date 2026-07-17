#!/usr/bin/env bash
# Shared post-PR headless execution-context detector.
#
# The fixed exec-headless manager persists the policy into canonical worker lineage
# and mirrors it to control-plane worker capabilities.  Aone-facing entrypoints ask
# that manager to match the current process against the live lineage.  The old /tmp
# process-group marker remains only a rollout/cache fallback; deleting or corrupting
# it cannot make a registered post-PR worker writable.

jarvis_post_pr_context_dir() {
    printf '/tmp/jarvis-headless-context-%s' "$(id -u)"
}

_jarvis_post_pr_marker_valid() {
    local marker="$1" expected_pid="$2" expected_pgid="$3" content
    [ -f "$marker" ] || return 1
    content="$(cat "$marker" 2>/dev/null)" || return 1
    case "$content" in
        *'"pid":'"$expected_pid"','*'"pgid":'"$expected_pgid"','*) ;;
        *) return 1 ;;
    esac
    case "$content" in
        *'"policy_revision":"terraform-rd-single-writer-v4"'*|\
        *'"policy_revision":"terraform-rd-single-writer-v3"'*) ;;
        *) return 1 ;;
    esac
    case "$content" in
        *'"kind":"pr_ci_fix"'*|*'"kind":"pr_comment_reply"'*) return 0 ;;
        *) return 1 ;;
    esac
}

jarvis_post_pr_lineage_active() {
    local helper_dir manager
    helper_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    manager="$helper_dir/jarvis-interactive-worker.py"
    [ -f "$manager" ] || return 1
    /usr/bin/python3 -I "$manager" post-pr-context --pid "$$" \
        >/dev/null 2>&1
}

jarvis_post_pr_exec_lineage_active() {
    local pid="$$" ppid command depth=0
    while [ "$depth" -lt 64 ]; do
        case "$pid" in
            ''|*[!0-9]*|0|1) break ;;
        esac
        command="$(/bin/ps -ww -o command= -p "$pid" 2>/dev/null)" || command=""
        case "$command" in
            *"/bootstrap/jarvis-interactive-worker.py exec-headless"*\
*"--policy-revision terraform-rd-single-writer-v4"*\
*"--aone-write-policy post-pr-read-only"*\
*"--headless-kind pr_ci_fix"*\
*"--aone-id "*\
*"--project-id "*) return 0 ;;
            *"/bootstrap/jarvis-interactive-worker.py exec-headless"*\
*"--policy-revision terraform-rd-single-writer-v4"*\
*"--aone-write-policy post-pr-read-only"*\
*"--headless-kind pr_comment_reply"*\
*"--aone-id "*\
*"--project-id "*) return 0 ;;
        esac
        ppid="$(/bin/ps -o ppid= -p "$pid" 2>/dev/null | tr -d '[:space:]')"
        [ -n "$ppid" ] && [ "$ppid" != "$pid" ] || break
        pid="$ppid"
        depth=$((depth + 1))
    done
    return 1
}

jarvis_post_pr_context_active() {
    case "${JARVIS_AONE_WRITE_POLICY:-}" in
        post-pr-read-only|post-pr-no-comment) return 0 ;;
    esac

    # Live argv is kernel-owned for the lifetime of the managed guard, so it is the
    # fail-closed fallback when the canonical-state helper is missing or unhealthy.
    jarvis_post_pr_exec_lineage_active && return 0
    jarvis_post_pr_lineage_active && return 0

    local pid="$$" ppid pgid marker depth=0
    while [ "$depth" -lt 64 ]; do
        case "$pid" in
            ''|*[!0-9]*|0|1) break ;;
        esac
        pgid="$(/bin/ps -o pgid= -p "$pid" 2>/dev/null | tr -d '[:space:]')"
        if [ -n "$pgid" ]; then
            marker="$(jarvis_post_pr_context_dir)/${pgid}.json"
            _jarvis_post_pr_marker_valid "$marker" "$pid" "$pgid" && return 0
        fi
        ppid="$(/bin/ps -o ppid= -p "$pid" 2>/dev/null | tr -d '[:space:]')"
        [ -n "$ppid" ] && [ "$ppid" != "$pid" ] || break
        pid="$ppid"
        depth=$((depth + 1))
    done
    return 1
}
