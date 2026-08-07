#!/usr/bin/env bash
# Stop hook: control-plane session exit gate with wrap-check.sh fallback.
#
# Primary path: delegates to jarvis-interactive-worker.py stop-check which reads
# the session state cached by the 30s heartbeat sidecar. This replaces the older
# approach of scanning .my-day/ ledger files.
#
# Fallback: if the worker state is unavailable (exit 1), falls back to wrap-check.sh
# for backward compatibility and degraded-mode safety.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wrap_check="$script_dir/wrap-check.sh"

# --- Resolve Python + manager through the shared machine runtime config ---
# shellcheck disable=SC1091
source "$script_dir/runtime-config.sh"
jarvis_load_runtime_config || exit $?
# stop-check is an ordinary worker endpoint; never retain or forward an admin
# credential that happens to be present on an operator workstation.
unset JARVIS_CONTROL_PLANE_ADMIN_TOKEN

python_bin="${JARVIS_INTERACTIVE_WORKER_PYTHON:-}"
if [ -z "$python_bin" ]; then
  if [ -x /usr/bin/python3 ]; then
    python_bin="/usr/bin/python3"
  else
    python_bin="python3"
  fi
fi
manager="${JARVIS_INTERACTIVE_WORKER_MANAGER:-$script_dir/jarvis-interactive-worker.py}"

mode="${1:-wrap-only}"

# --- Primary path: control-plane session exit gate ---
_run_stop_check() {
    "$python_bin" "$manager" stop-check 2>&1
}

stop_output="$(_run_stop_check)"
rc=$?

if [ "$rc" -eq 0 ]; then
    # Control-plane path: safe to stop.
    if [ "$mode" = "codex" ]; then
        # Codex: also need to record the turn Stop event via the worker hook.
        payload="$(mktemp)"
        worker_output="$(mktemp)"
        trap 'rm -f "$payload" "$worker_output"' EXIT
        cat > "$payload"

        worker_hook="$script_dir/run-interactive-worker-hook.sh"
        if [ ! -x "$worker_hook" ]; then
            echo "stop-hook: interactive worker hook not found at $worker_hook" >&2
            exit 2
        fi

        bash "$worker_hook" codex Stop < "$payload" > "$worker_output"
        wrc=$?
        if [ "$wrc" -ne 0 ]; then
            [ ! -s "$worker_output" ] || cat "$worker_output" >&2
            echo "stop-hook: interactive worker Stop signal failed with rc=$wrc" >&2
            exit 2
        fi
        printf '{}\n'
    fi
    exit 0
fi

if [ "$rc" -eq 2 ]; then
    # Control-plane path: explicitly blocked.
    [ -n "$stop_output" ] && echo "$stop_output" >&2
    exit 2
fi

# --- Fallback: wrap-check.sh (state unavailable, Python crash, etc.) ---
if [ ! -f "$wrap_check" ]; then
    echo "stop-hook: worker state unavailable and wrap-check.sh not found" >&2
    exit 2
fi

if [ "$mode" != "codex" ]; then
    exec bash "$wrap_check"
fi

# Codex fallback: same ordered sequence as before.
payload="$(mktemp)"
wrap_output="$(mktemp)"
worker_output="$(mktemp)"
trap 'rm -f "$payload" "$wrap_output" "$worker_output"' EXIT
cat > "$payload"

bash "$wrap_check" < "$payload" > "$wrap_output"
rc=$?
[ ! -s "$wrap_output" ] || cat "$wrap_output" >&2
if [ "$rc" -ne 0 ]; then
    [ "$rc" -eq 2 ] || echo "stop-hook: wrap-check failed with rc=$rc (degraded mode)" >&2
    exit 2
fi

worker_hook="$script_dir/run-interactive-worker-hook.sh"
if [ ! -x "$worker_hook" ]; then
    echo "stop-hook: interactive worker hook not found at $worker_hook" >&2
    exit 2
fi

bash "$worker_hook" codex Stop < "$payload" > "$worker_output"
rc=$?
if [ "$rc" -ne 0 ]; then
    [ ! -s "$worker_output" ] || cat "$worker_output" >&2
    echo "stop-hook: interactive worker Stop signal failed with rc=$rc" >&2
    exit 2
fi
printf '{}\n'
