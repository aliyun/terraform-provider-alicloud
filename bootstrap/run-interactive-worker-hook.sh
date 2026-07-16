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
expected_event="${2:-}"

case "$mode" in
  claude|codex)
    if [ "$expected_event" = "PreToolUse" ]; then
      # Codex only treats exit 2 + non-empty stderr as a blocking PreTool
      # result. Normalize missing Python/manager crashes/malformed input and
      # every other nonzero into that contract for both clients.
      payload="$(mktemp)" || {
        echo "worker-fence-hook: cannot allocate hook payload" >&2
        exit 2
      }
      output="$(mktemp)" || {
        rm -f "$payload"
        echo "worker-fence-hook: cannot allocate hook output" >&2
        exit 2
      }
      trap 'rm -f "$payload" "$output"' EXIT
      if ! cat > "$payload"; then
        echo "worker-fence-hook: cannot read PreToolUse payload" >&2
        exit 2
      fi
      trusted_python="/usr/bin/python3"
      trusted_manager="$script_dir/jarvis-interactive-worker.py"
      if [ ! -x "$trusted_python" ] || [ ! -f "$trusted_manager" ]; then
        echo "worker-fence-hook: trusted PreToolUse manager is unavailable" >&2
        exit 2
      fi
      # PreTool fencing is a security boundary: runtime env files may select a
      # manager for ordinary lifecycle/testing, but must not replace the code
      # that decides whether a concrete tool is allowed to execute.
      "$trusted_python" -I "$trusted_manager" hook "$mode" \
        --expected-event PreToolUse < "$payload" > "$output"
      rc=$?
      if [ "$rc" -eq 0 ]; then
        cat "$output"
        exit 0
      fi
      [ ! -s "$output" ] || cat "$output" >&2
      if [ "$rc" -ne 2 ]; then
        echo "worker-fence-hook: manager failed with rc=$rc; tool blocked" >&2
      fi
      exit 2
    fi
    if [ -n "$expected_event" ]; then
      exec "$python_bin" "$manager" hook "$mode" \
        --expected-event "$expected_event"
    fi
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
