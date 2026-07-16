#!/usr/bin/env bash
# Shared Claude/Codex SessionStart/SessionEnd and interactive-worker CLI wrapper.
# It deliberately loads the same gitignored environment as bridge/run.sh so a
# developer opening the repository does not need a second control-plane token.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"

# Worktrees share the main repository's gitignored bootstrap/.env and
# bridge/jarvis.env.  Explicit path overrides keep tests and unusual installs
# deterministic.
main_root="$repo_root"
git_common="$(git -C "$repo_root" rev-parse --path-format=absolute --git-common-dir 2>/dev/null || true)"
case "$git_common" in
  */.git) main_root="${git_common%/.git}" ;;
esac
bootstrap_env="${JARVIS_INTERACTIVE_BOOTSTRAP_ENV:-$main_root/bootstrap/.env}"
bridge_env="${JARVIS_INTERACTIVE_BRIDGE_ENV:-$main_root/bridge/jarvis.env}"

set -a
set +u
# shellcheck disable=SC1090
[ -f "$bootstrap_env" ] && . "$bootstrap_env"
# shellcheck disable=SC1090
[ -f "$bridge_env" ] && . "$bridge_env"
set -u
set +a

# The control plane intentionally reuses the established report credential.
# Keep the fallback in-process only; the Python state file never persists it.
if [ -z "${JARVIS_CONTROL_PLANE_TOKEN:-}" ] && [ -n "${JARVIS_HTML_REPORT_TOKEN:-}" ]; then
  export JARVIS_CONTROL_PLANE_TOKEN="$JARVIS_HTML_REPORT_TOKEN"
fi
if [ -z "${JARVIS_CONTROL_PLANE_BASE_URL:-}" ] && [ -n "${JARVIS_HTML_REPORT_BASE_URL:-}" ]; then
  export JARVIS_CONTROL_PLANE_BASE_URL="$JARVIS_HTML_REPORT_BASE_URL"
fi

python_bin="${JARVIS_INTERACTIVE_WORKER_PYTHON:-python3}"
manager="${JARVIS_INTERACTIVE_WORKER_MANAGER:-$script_dir/jarvis-interactive-worker.py}"
mode="${1:-}"

case "$mode" in
  claude|codex)
    exec "$python_bin" "$manager" hook "$mode"
    ;;
  cli)
    shift
    exec "$python_bin" "$manager" "$@"
    ;;
  *)
    echo "usage: run-interactive-worker-hook.sh <claude|codex|cli> [args...]" >&2
    exit 64
    ;;
esac
