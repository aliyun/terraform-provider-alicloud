#!/bin/bash
# bootstrap/lib.sh — shared utilities for jarvis bootstrap scripts
#
# Source this file to get common functions. Safe to source multiple times.

# Resolve the MAIN repo root, even from inside a git worktree.
# Uses git-common-dir which always points to the main repo's .git,
# then strips the trailing /.git to get the repo root.
# Respects JARVIS_ROOT env override (used in tests and explicit calls).
jarvis_root() {
    if [ -n "${JARVIS_ROOT:-}" ]; then
        echo "$JARVIS_ROOT"
        return
    fi
    local git_common
    git_common="$(git -C "$(dirname "${BASH_SOURCE[1]:-${BASH_SOURCE[0]}}")" \
        rev-parse --path-format=absolute --git-common-dir 2>/dev/null)" \
        || git_common="$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null)" \
        || { echo "lib.sh: not in a git repo" >&2; return 1; }
    # git-common-dir returns /path/to/repo/.git — strip trailing /.git
    echo "${git_common%/.git}"
}

# Standard runs/ dir(P1.e 抽出:plan.sh + log.sh + wrap-check.sh 曾各自重复实现)。
# 尊重 JARVIS_RUNS_DIR env override(测试用)。
lib_runs_dir() {
    echo "${JARVIS_RUNS_DIR:-$(jarvis_root)/runs}"
}

# Standard escalation/ dir(P1.e 抽出:log.sh + sweep/watchdog 曾各自重复实现)。
# 尊重 JARVIS_ESCALATION_DIR env override(测试用)。
lib_escalation_dir() {
    echo "${JARVIS_ESCALATION_DIR:-$(jarvis_root)/escalation}"
}

# Effective owner/self for claim-ledger attribution (cap Phase 2, D2).
# Resolution order:
#   1. COORD_ID           — coord.sh instance id (triage-one / bridge loops export it)
#   2. cc-<session-id>    — interactive/headless Claude session; CLAUDE_CODE_SESSION_ID
#                           is stable across all Bash tool subprocesses of one session,
#                           so two different sessions get distinct owners and no longer
#                           block each other in wrap-check.
#   3. codex-<thread-id>   — interactive Codex thread; CODEX_THREAD_ID has the same
#                           cross-tool stability as Claude's native session id.
#   4. persisted hook id   — Claude's CLAUDE_ENV_FILE fallback when the native id is
#                           absent from a later tool subprocess.
#   5. ""                 — no stable source → ownerless (legacy behavior: wrap-check
#                           still holds you to account, no regression).
# claim.sh (writes owner) and wrap-check.sh (computes self) MUST both use this so the
# values match within a session.
coord_self() {
    if [ -n "${COORD_ID:-}" ]; then
        echo "$COORD_ID"
    elif [ -n "${CLAUDE_CODE_SESSION_ID:-}" ]; then
        echo "cc-$CLAUDE_CODE_SESSION_ID"
    elif [ -n "${CODEX_THREAD_ID:-}" ]; then
        echo "codex-$CODEX_THREAD_ID"
    elif [ "${JARVIS_INTERACTIVE_CLIENT:-}" = "claude" ] \
            && [ -n "${JARVIS_INTERACTIVE_SESSION_ID:-}" ]; then
        echo "cc-$JARVIS_INTERACTIVE_SESSION_ID"
    elif [ "${JARVIS_INTERACTIVE_CLIENT:-}" = "codex" ] \
            && [ -n "${JARVIS_INTERACTIVE_SESSION_ID:-}" ]; then
        echo "codex-$JARVIS_INTERACTIVE_SESSION_ID"
    else
        echo ""
    fi
}
