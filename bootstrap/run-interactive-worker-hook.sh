#!/usr/bin/env bash
# Shared Claude/Codex SessionStart/SessionEnd and interactive-worker CLI wrapper.
# It deliberately loads the same gitignored environment as bridge/run.sh so a
# developer opening the repository does not need a second control-plane token.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
# shellcheck disable=SC1091
source "$script_dir/runtime-config.sh"
jarvis_load_runtime_config || exit $?
# Interactive workers never call /admin/**. Keep an operator credential loaded
# by the shared runtime config out of this process and every detached sidecar.
unset JARVIS_CONTROL_PLANE_ADMIN_TOKEN

if [ -n "${JARVIS_INTERACTIVE_WORKER_PYTHON:-}" ]; then
  python_bin="$JARVIS_INTERACTIVE_WORKER_PYTHON"
elif [ -x /usr/bin/python3 ]; then
  # The macOS system runtime uses the platform trust store. Framework/Homebrew
  # Python installations may not trust internal control-plane certificates,
  # which would silently stop the detached heartbeat sidecar.
  python_bin="/usr/bin/python3"
else
  python_bin="python3"
fi
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
