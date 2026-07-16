#!/usr/bin/env bash
# test/bridge_register_headless_test.sh — hermetic tests for the bridge side of the
# headless worker-fence registration fix.
#
# Proves: the bridge always launches through the fixed isolated manager wrapper;
# run_claude_buffered / run_claude_stream wrap the complete Claude command; and
# dispatch_item creates a fresh pre-exec wrapper on EVERY retry attempt.
#
# Run: bash test/bridge_register_headless_test.sh

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

if ! command -v python3 >/dev/null 2>&1; then
  echo "SKIP bridge_register_headless_test: python3 not available"
  exit 0
fi

BRIDGE_DIR="$repo_root/bridge" python3 - "$repo_root" <<'PY'
import os
import sys
import json
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, os.environ["BRIDGE_DIR"])
import jarvis_dingtalk_bot as b

b.INFLIGHT_PATH = Path(tempfile.mkdtemp()) / "inflight.json"


class FakeProc:
    def __init__(self, pid, out):
        self.pid = pid
        self._out = out
        self.returncode = 0

    def communicate(self, timeout=None):
        return self._out, ""

    def wait(self, timeout=None):
        return 0

    def kill(self):
        pass


class HelperTest(unittest.TestCase):
    def test_helper_uses_fixed_isolated_manager(self):
        argv = b._headless_exec_command(
            "sid-x", ["/opt/claude", "--resume", "sid-x"])
        self.assertEqual(argv[:2], ["/usr/bin/python3", "-I"])
        self.assertTrue(argv[2].endswith(
            "/bootstrap/jarvis-interactive-worker.py"))
        self.assertEqual(argv[3:8], [
            "exec-headless", "--session-id", "sid-x", "--client", "claude"])
        self.assertEqual(argv[8], "--")
        self.assertEqual(argv[9:], ["/opt/claude", "--resume", "sid-x"])

    def test_helper_rejects_missing_identity_or_command(self):
        with self.assertRaises(ValueError):
            b._headless_exec_command("", ["/opt/claude"])
        with self.assertRaises(ValueError):
            b._headless_exec_command("sid-y", [])


class BufferedRegistersTest(unittest.TestCase):
    def test_buffered_registers_spawned_worker(self):
        seen = []
        orig_wrap = b._headless_exec_command
        orig_popen = b.subprocess.Popen
        def _wrap(sid, command):
            seen.append((sid, list(command)))
            return list(command)
        b._headless_exec_command = _wrap
        b.subprocess.Popen = lambda *a, **k: FakeProc(
            7777, '{"type":"result","result":"ok","is_error":false,"subtype":"success"}')
        try:
            b.run_claude_buffered("hi", "sid-buf", False, timeout=5)
        finally:
            b._headless_exec_command = orig_wrap
            b.subprocess.Popen = orig_popen
        self.assertEqual(len(seen), 1)
        self.assertEqual(seen[0][0], "sid-buf")
        self.assertIn("--session-id", seen[0][1])

    def test_stream_registers_spawned_worker(self):
        seen = []
        orig_wrap = b._headless_exec_command
        orig_popen = b.subprocess.Popen
        def _wrap(sid, command):
            seen.append((sid, list(command)))
            return list(command)
        b._headless_exec_command = _wrap

        class _StreamProc(FakeProc):
            def __init__(self):
                super().__init__(8888, "")
                self.stdout = iter(())
                self.stderr = None

        b.subprocess.Popen = lambda *a, **k: _StreamProc()
        orig_parse = b.parse_stream_lines
        b.parse_stream_lines = lambda stream: iter(())
        try:
            list(b.run_claude_stream("hi", "sid-str", False, timeout=5))
        finally:
            b._headless_exec_command = orig_wrap
            b.subprocess.Popen = orig_popen
            b.parse_stream_lines = orig_parse
        self.assertEqual(len(seen), 1)
        self.assertEqual(seen[0][0], "sid-str")
        self.assertIn("--session-id", seen[0][1])


class DispatchRetryRegistersTest(unittest.TestCase):
    def test_registration_happens_every_attempt(self):
        os.environ["JARVIS_DISPATCH_RETRY_MAX"] = "2"
        os.environ["JARVIS_DISPATCH_RETRY_BACKOFF"] = "0"
        seen = []
        pids = iter([1001, 1002, 1003])
        # First two attempts return a transient error → retry; third succeeds.
        outs = iter([
            '{"type":"result","result":"boom","is_error":true,"subtype":"error_during_execution"}',
            '{"type":"result","result":"boom","is_error":true,"subtype":"error_during_execution"}',
            '{"type":"result","result":"done","is_error":false,"subtype":"success"}',
        ])

        class _P(FakeProc):
            def __init__(self):
                super().__init__(next(pids), next(outs))
                self.returncode = 0 if "success" in self._out else 1

        orig_wrap = b._headless_exec_command
        orig_popen = b.subprocess.Popen
        def _wrap(sid, command):
            seen.append((sid, list(command)))
            return list(command)
        b._headless_exec_command = _wrap
        b.subprocess.Popen = lambda *a, **k: _P()

        h = b.JarvisHandler(no_dingtalk=True)
        h._maybe_suspend = lambda *a, **k: None
        h._completion_broadcast = lambda item_id: "ok"
        try:
            rv = h.dispatch_item("84000001", "prompt", "sid-retry", False,
                                 lambda msg: None, "t", "chat", project=None)
        finally:
            b._headless_exec_command = orig_wrap
            b.subprocess.Popen = orig_popen
            os.environ.pop("JARVIS_DISPATCH_RETRY_MAX", None)
            os.environ.pop("JARVIS_DISPATCH_RETRY_BACKOFF", None)
        self.assertEqual(rv, "done")
        self.assertEqual(len(seen), 3,
                         "each retry attempt must create a fresh pre-exec wrapper")
        self.assertTrue(all(sid == "sid-retry" for sid, _ in seen))
        self.assertIn("--session-id", seen[0][1])
        self.assertIn("--resume", seen[1][1])
        self.assertIn("--resume", seen[2][1])


if __name__ == "__main__":
    res = unittest.main(argv=["x", "-v"], exit=False).result
    sys.exit(0 if res.wasSuccessful() else 1)
PY
rc=$?
if [ "$rc" -eq 0 ]; then
  echo "bridge_register_headless_test: PASS"
else
  echo "bridge_register_headless_test: FAIL"
fi
exit "$rc"
