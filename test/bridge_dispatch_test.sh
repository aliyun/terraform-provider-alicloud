#!/usr/bin/env bash
# Hermetic regression entrypoint for Bridge dispatch and Scheduler runners.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"
python_bin="${BRIDGE_PYTHON:-$repo_root/.venv/bridge/bin/python}"

if [ ! -x "$python_bin" ]; then
  python_bin="$(command -v python3 || true)"
fi
if [ -z "$python_bin" ]; then
  echo "SKIP bridge_dispatch_test: python3 not available"
  exit 0
fi

PYTHONPATH="$repo_root" "$python_bin" -m unittest \
  bridge.test_aone_event_publisher \
  bridge.test_bot_control_plane \
  bridge.test_capacity \
  bridge.test_ephemeral_executor \
  bridge.test_external_recovery \
  bridge.test_field_repair \
  bridge.test_headless_runtime \
  bridge.test_persistent_task_execution \
  bridge.test_prwatch_ci_fix \
  bridge.test_stale_revisit \
  bridge.test_task_router \
  bridge.scheduler.tests.test_runners
rc=$?

if [ "$rc" -eq 0 ]; then
  echo "bridge_dispatch_test: PASS"
else
  echo "bridge_dispatch_test: FAIL (rc=$rc)"
fi
exit "$rc"
