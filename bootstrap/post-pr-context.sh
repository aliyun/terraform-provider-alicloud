#!/usr/bin/env bash
# Shared post-PR headless execution-context detector.
#
# Bridge registers the guarded worker's process group in a fixed per-user runtime
# directory before releasing the child command.  Aone-facing shell entrypoints walk
# their ancestor chain and reject writes when any ancestor belongs to that registered
# process group.  This keeps the read-only boundary in force even if a prompted child
# runs a command through `env -u JARVIS_AONE_WRITE_POLICY`.

jarvis_post_pr_context_dir() {
    printf '/tmp/jarvis-headless-context-%s' "$(id -u)"
}

_jarvis_post_pr_marker_valid() {
    local marker="$1" content
    [ -f "$marker" ] || return 1
    content="$(cat "$marker" 2>/dev/null)" || return 1
    case "$content" in
        *'"policy_revision":"terraform-rd-single-writer-v3"'*) ;;
        *) return 1 ;;
    esac
    case "$content" in
        *'"kind":"pr_ci_fix"'*|*'"kind":"pr_comment_reply"'*) return 0 ;;
        *) return 1 ;;
    esac
}

jarvis_post_pr_context_active() {
    case "${JARVIS_AONE_WRITE_POLICY:-}" in
        post-pr-read-only|post-pr-no-comment) return 0 ;;
    esac

    local pid="$$" ppid pgid marker depth=0
    while [ "$depth" -lt 64 ]; do
        case "$pid" in
            ''|*[!0-9]*|0|1) break ;;
        esac
        pgid="$(ps -o pgid= -p "$pid" 2>/dev/null | tr -d '[:space:]')"
        if [ -n "$pgid" ]; then
            marker="$(jarvis_post_pr_context_dir)/${pgid}.json"
            _jarvis_post_pr_marker_valid "$marker" && return 0
        fi
        ppid="$(ps -o ppid= -p "$pid" 2>/dev/null | tr -d '[:space:]')"
        [ -n "$ppid" ] && [ "$ppid" != "$pid" ] || break
        pid="$ppid"
        depth=$((depth + 1))
    done
    return 1
}
