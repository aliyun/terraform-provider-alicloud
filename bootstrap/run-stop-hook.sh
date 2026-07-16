#!/usr/bin/env bash
# Resolve and run the project Stop hook. Codex also records turn completion,
# but only after wrap-check has allowed the turn to stop.

set -uo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
wrap_check="$script_dir/wrap-check.sh"

if [ ! -f "$wrap_check" ]; then
    echo "stop-hook: wrap-check.sh not found at $wrap_check" >&2
    exit 2
fi

mode="${1:-wrap-only}"
if [ "$mode" != "codex" ]; then
    exec bash "$wrap_check"
fi

payload="$(mktemp)"
wrap_output="$(mktemp)"
worker_output="$(mktemp)"
trap 'rm -f "$payload" "$wrap_output" "$worker_output"' EXIT
cat > "$payload"

# Codex launches multiple handlers for one event concurrently. Keep this as a
# single ordered handler so a blocking wrap-check cannot race with a Stop signal
# that would pause a still-running task lease.
bash "$wrap_check" < "$payload" > "$wrap_output"
rc=$?
[ ! -s "$wrap_output" ] || cat "$wrap_output" >&2
if [ "$rc" -ne 0 ]; then
    [ "$rc" -eq 2 ] || echo "stop-hook: wrap-check failed with rc=$rc" >&2
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
