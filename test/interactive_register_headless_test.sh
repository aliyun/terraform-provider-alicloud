#!/usr/bin/env bash
# test/interactive_register_headless_test.sh — hermetic tests for the manager's
# `register-headless` CLI/registration path (bridge pre-registers each headless
# claude spawn so the worker-fence never trips on a retry/resume attempt).
#
# Proves: local state lands, idempotent, control-plane-unavailable tolerated
# (local file still written), a registered idle worker passes the PreToolUse
# fence, and — critically — the fence is NOT degraded (an UNregistered session
# still fails closed).
#
# Run: bash test/interactive_register_headless_test.sh

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

if ! command -v python3 >/dev/null 2>&1; then
  echo "SKIP interactive_register_headless_test: python3 not available"
  exit 0
fi

REPO_ROOT="$repo_root" python3 - <<'PY'
import os
import sys
import json
import tempfile
import importlib.util
import unittest
from pathlib import Path

REPO = Path(os.environ["REPO_ROOT"])
MANAGER = REPO / "bootstrap" / "jarvis-interactive-worker.py"

spec = importlib.util.spec_from_file_location("jarvis_interactive_worker_mod", MANAGER)
m = importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)


def _fresh_state_dir():
    d = tempfile.mkdtemp()
    os.environ["JARVIS_INTERACTIVE_STATE_DIR"] = d
    return d


class RegisterHeadlessTest(unittest.TestCase):
    def setUp(self):
        _fresh_state_dir()
        # Keep every test off the network: record register attempts instead of
        # hitting the control plane.
        self._orig_register = m._register
        self.register_calls = []

        def _rec(client, state, status="ACTIVE"):
            self.register_calls.append((status, dict(state)))

        m._register = _rec
        self._orig_client = m._client
        m._client = lambda **k: object()

    def tearDown(self):
        m._register = self._orig_register
        m._client = self._orig_client
        for var in ("CLAUDE_CODE_SESSION_ID", "CODEX_THREAD_ID",
                    "JARVIS_INTERACTIVE_CLIENT", "JARVIS_INTERACTIVE_SESSION_ID"):
            os.environ.pop(var, None)

    def test_local_state_lands(self):
        sid = "sess-alpha"
        pid = os.getpid()
        out = m.register_headless(sid, pid)
        path = m._state_path("claude", sid)
        self.assertTrue(path.exists(), "state file must be written")
        state = json.loads(path.read_text())
        self.assertEqual(state["client"], "claude")
        self.assertEqual(state["clientSessionId"], sid)
        self.assertEqual(int(state["hostPid"]), pid)
        self.assertFalse(state["stopped"], "verified live host → not stopped")
        self.assertIsNone(state["current"], "fresh worker carries no task")
        self.assertTrue(state["workerKey"])
        self.assertTrue(state["verifyHostCommand"])
        self.assertEqual(out["workerKey"], state["workerKey"])
        # A best-effort remote register was attempted.
        self.assertEqual(len(self.register_calls), 1)
        self.assertEqual(self.register_calls[0][0], "ACTIVE")

    def test_idempotent_same_pid(self):
        sid = "sess-beta"
        pid = os.getpid()
        a = m.register_headless(sid, pid)
        b = m.register_headless(sid, pid)  # must not raise
        self.assertEqual(a["workerKey"], b["workerKey"])
        self.assertTrue(b["sameIncarnation"])

    def test_remote_unavailable_still_lands_local(self):
        def _boom(client, state, status="ACTIVE"):
            raise m.ControlPlaneUnavailable("down")

        m._register = _boom
        sid = "sess-gamma"
        pid = os.getpid()
        # Must NOT raise even though the remote register is unreachable.
        m.register_headless(sid, pid)
        path = m._state_path("claude", sid)
        self.assertTrue(path.exists(), "local state must land despite remote outage")
        state = json.loads(path.read_text())
        self.assertFalse(state["stopped"])

    def _patch_host(self, pid, monkeystart="STARTID"):
        m._find_host_pid = lambda client: (pid, True)
        m._process_start_identity = lambda p: monkeystart
        m._process_info = lambda p: (1, "claude -p prompt")
        m._nearest_runtime_client = lambda: "claude"
        m._pid_alive = lambda p: True  # a freshly-spawned host is live

    def test_registered_idle_worker_passes_pre_tool_fence(self):
        orig = (m._find_host_pid, m._process_start_identity,
                m._process_info, m._nearest_runtime_client, m._pid_alive)
        try:
            pid = 424242
            self._patch_host(pid)
            sid = "sess-delta"
            m.register_headless(sid, pid)
            os.environ["CLAUDE_CODE_SESSION_ID"] = sid
            event = {
                "hook_event_name": "PreToolUse",
                "session_id": sid,
                "tool_name": "Bash",
                "tool_input": {"command": "echo hi"},
                "cwd": str(REPO),
            }
            rc = m.hook("claude", event)
            self.assertEqual(rc, 0, "registered idle worker must be allowed")
        finally:
            (m._find_host_pid, m._process_start_identity,
             m._process_info, m._nearest_runtime_client, m._pid_alive) = orig

    def test_unregistered_session_still_fails_closed(self):
        orig = (m._find_host_pid, m._process_start_identity,
                m._process_info, m._nearest_runtime_client, m._pid_alive)
        try:
            pid = 515151
            self._patch_host(pid)
            sid = "sess-unregistered"
            os.environ["CLAUDE_CODE_SESSION_ID"] = sid
            event = {
                "hook_event_name": "PreToolUse",
                "session_id": sid,
                "tool_name": "Bash",
                "tool_input": {"command": "echo hi"},
                "cwd": str(REPO),
            }
            rc = m.hook("claude", event)
            self.assertEqual(rc, m.HOOK_BLOCK_EXIT,
                             "fence must NOT be degraded: unknown session blocks")
        finally:
            (m._find_host_pid, m._process_start_identity,
             m._process_info, m._nearest_runtime_client, m._pid_alive) = orig


if __name__ == "__main__":
    res = unittest.main(argv=["x", "-v"], exit=False).result
    sys.exit(0 if res.wasSuccessful() else 1)
PY
rc=$?
if [ "$rc" -eq 0 ]; then
  echo "interactive_register_headless_test: PASS"
else
  echo "interactive_register_headless_test: FAIL"
fi
exit "$rc"
