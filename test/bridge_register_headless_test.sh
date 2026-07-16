#!/usr/bin/env bash
# test/bridge_register_headless_test.sh — hermetic tests for the bridge side of the
# headless worker-fence registration fix.
#
# Proves: the bridge helper shells out to the trusted manager CLI (best-effort,
# never raising); run_claude_buffered / run_claude_stream register the freshly
# spawned worker with (session_id, p.pid); and dispatch_item re-registers on
# EVERY retry attempt (mock counts calls == spawn count).
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
    def test_helper_invokes_manager_cli_bestexpr(self):
        calls = {}

        class _R:
            returncode = 0

        def _run(argv, **kw):
            calls["argv"] = argv
            calls["kw"] = kw
            return _R()

        orig = b.subprocess.run
        b.subprocess.run = _run
        try:
            b._register_headless_worker("sid-x", 4321)
        finally:
            b.subprocess.run = orig
        argv = calls["argv"]
        self.assertIn("run-interactive-worker-hook.sh", " ".join(str(a) for a in argv))
        self.assertIn("cli", argv)
        self.assertIn("register-headless", argv)
        joined = " ".join(str(a) for a in argv)
        self.assertIn("sid-x", joined)
        self.assertIn("4321", joined)

    def test_helper_never_raises_on_failure(self):
        def _boom(argv, **kw):
            raise OSError("no such file")

        orig = b.subprocess.run
        b.subprocess.run = _boom
        try:
            # Must swallow: registration failure must never block dispatch.
            b._register_headless_worker("sid-y", 99)
        finally:
            b.subprocess.run = orig


class BufferedRegistersTest(unittest.TestCase):
    def test_buffered_registers_spawned_worker(self):
        seen = []
        orig_reg = b._register_headless_worker
        orig_popen = b.subprocess.Popen
        b._register_headless_worker = lambda sid, pid: seen.append((sid, pid))
        b.subprocess.Popen = lambda *a, **k: FakeProc(
            7777, '{"type":"result","result":"ok","is_error":false,"subtype":"success"}')
        try:
            b.run_claude_buffered("hi", "sid-buf", False, timeout=5)
        finally:
            b._register_headless_worker = orig_reg
            b.subprocess.Popen = orig_popen
        self.assertEqual(seen, [("sid-buf", 7777)])

    def test_stream_registers_spawned_worker(self):
        seen = []
        orig_reg = b._register_headless_worker
        orig_popen = b.subprocess.Popen
        b._register_headless_worker = lambda sid, pid: seen.append((sid, pid))

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
            b._register_headless_worker = orig_reg
            b.subprocess.Popen = orig_popen
            b.parse_stream_lines = orig_parse
        self.assertEqual(seen, [("sid-str", 8888)])


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

        orig_reg = b._register_headless_worker
        orig_popen = b.subprocess.Popen
        b._register_headless_worker = lambda sid, pid: seen.append(pid)
        b.subprocess.Popen = lambda *a, **k: _P()

        h = b.JarvisHandler(no_dingtalk=True)
        h._maybe_suspend = lambda *a, **k: None
        h._completion_broadcast = lambda item_id: "ok"
        try:
            rv = h.dispatch_item("84000001", "prompt", "sid-retry", False,
                                 lambda msg: None, "t", "chat", project=None)
        finally:
            b._register_headless_worker = orig_reg
            b.subprocess.Popen = orig_popen
            os.environ.pop("JARVIS_DISPATCH_RETRY_MAX", None)
            os.environ.pop("JARVIS_DISPATCH_RETRY_BACKOFF", None)
        self.assertEqual(rv, "done")
        self.assertEqual(seen, [1001, 1002, 1003],
                         "each retry attempt must re-register the new pid")


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
