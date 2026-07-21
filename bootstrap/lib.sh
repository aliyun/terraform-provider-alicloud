#!/bin/bash
# bootstrap/lib.sh — shared utilities for jarvis bootstrap scripts
#
# Source this file to get common functions. Safe to source multiple times.

_jarvis_lib_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "$_jarvis_lib_dir/runtime-config.sh"
jarvis_load_runtime_config || return $?

# Resolve the MAIN repo root, even from inside a git worktree.
# Uses git-common-dir which always points to the main repo's .git,
# then strips the trailing /.git to get the repo root.
# Respects JARVIS_ROOT env override (used in tests and explicit calls).
jarvis_root() {
    if [ -n "${JARVIS_ROOT:-}" ]; then
        echo "$JARVIS_ROOT"
        return
    fi
    # Anchor at the CALLING script's dir (BASH_SOURCE[1] = whoever sourced us)
    # so resolution targets the jarvis repo regardless of process CWD; fall
    # back to CWD if that dir isn't a repo.
    #
    # NO `--path-format=absolute` here: that flag needs git ≥ 2.31, and older
    # git (AliOS/RHEL yum builds) echoes the unknown flag as literal rev-parse
    # output — "--path-format=absolute\n.." ended up inside generated cache paths
    # on the first Linux worker, breaking the preflight stamp. Instead
    # absolutize the (possibly relative) --git-common-dir via cd + pwd, which
    # works on every git that has worktrees at all.
    local anchor base git_common
    anchor="$(cd "$(dirname "${BASH_SOURCE[1]:-${BASH_SOURCE[0]}}")" 2>/dev/null && pwd)"
    base="${anchor:-$PWD}"
    git_common="$(cd "$base" 2>/dev/null && git rev-parse --git-common-dir 2>/dev/null)" \
        || { base="$PWD"; git_common="$(git rev-parse --git-common-dir 2>/dev/null)"; } \
        || { echo "lib.sh: not in a git repo" >&2; return 1; }
    case "$git_common" in
        /*) ;;
        *) git_common="$base/$git_common" ;;
    esac
    # Canonicalize ../ segments; for worktrees --git-common-dir points at the
    # MAIN repo's .git, which is exactly the semantic we want.
    git_common="$(cd "$git_common" 2>/dev/null && pwd)" \
        || { echo "lib.sh: git common dir unresolvable: $git_common" >&2; return 1; }
    # git-common-dir returns /path/to/repo/.git — strip trailing /.git
    echo "${git_common%/.git}"
}

# Standard runs/ dir(P1.e 抽出:plan.sh + log.sh + wrap-check.sh 曾各自重复实现)。
# 尊重 JARVIS_RUNS_DIR env override(测试用)。
lib_runs_dir() {
    echo "${JARVIS_RUNS_DIR:-$(jarvis_root)/runs}"
}

# Effective owner/self for claim-ledger attribution.
# Resolution order:
#   1. JARVIS_CLAIM_OWNER — explicit stable control-plane worker/session attribution
#   2. cc-<session-id>    — interactive/headless Claude session; CLAUDE_CODE_SESSION_ID
#                           is stable across all Bash tool subprocesses of one session,
#                           so two different sessions get distinct owners and no longer
#                           block each other in wrap-check.
#   3. codex-<thread-id>   — interactive Codex thread; CODEX_THREAD_ID has the same
#                           cross-tool stability as Claude's native session id.
#   4. persisted hook id   — hook fallback when the native id is
#                           absent from a later tool subprocess.
#   5. ""                 — no stable source → ownerless (legacy behavior: wrap-check
#                           still holds you to account, no regression).
# claim.sh (writes owner) and wrap-check.sh (computes self) MUST both use this so the
# values match within a session.
claim_owner() {
    if [ -n "${JARVIS_CLAIM_OWNER:-}" ]; then
        echo "$JARVIS_CLAIM_OWNER"
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

# Interactive database-fenced context detection — same predicate as claim.sh's
# local _is_interactive_context (kept there verbatim to avoid churn against
# parallel branches). New callers (wrap.sh receipts) use this shared copy.
jarvis_interactive_context() {
    [ -n "${CLAUDE_CODE_SESSION_ID:-}" ] || [ -n "${CODEX_THREAD_ID:-}" ] || \
        { case "${JARVIS_INTERACTIVE_CLIENT:-}" in claude|codex) true ;; *) false ;; esac \
          && [ -n "${JARVIS_INTERACTIVE_SESSION_ID:-}" ]; }
}

# Invoke the interactive worker CLI (jarvis-interactive-worker.py via the env-
# loading runner). JARVIS_INTERACTIVE_WORKER_RUNNER lets tests inject a stub.
jarvis_interactive_worker_cli() {
    local lib_dir
    lib_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    bash "${JARVIS_INTERACTIVE_WORKER_RUNNER:-$lib_dir/run-interactive-worker-hook.sh}" cli "$@"
}
