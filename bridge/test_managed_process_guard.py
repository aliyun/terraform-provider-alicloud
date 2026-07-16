#!/usr/bin/env python3
"""Crash-boundary tests for managed Claude process guarding."""

import os
import signal
import subprocess
import sys
import tempfile
import time
import unittest
from pathlib import Path
from unittest import mock

sys.path.insert(0, str(Path(__file__).resolve().parent))

import jarvis_dingtalk_bot as bot


BRIDGE_DIR = Path(__file__).resolve().parent
GUARD = BRIDGE_DIR / "managed_process_guard.py"


def _pid_alive(pid):
    try:
        os.kill(pid, 0)
        return True
    except ProcessLookupError:
        return False


class ManagedProcessGuardTest(unittest.TestCase):
    def test_command_starts_only_after_fenced_bind_returns(self):
        with tempfile.TemporaryDirectory() as directory:
            marker = Path(directory) / "started"
            script = (
                "from pathlib import Path; import json; "
                "Path(%r).write_text('started'); "
                "print(json.dumps({'type':'result','result':'ok',"
                "'is_error':False,'subtype':'success'}), flush=True)"
            ) % str(marker)
            bound = []

            def bind(process):
                self.assertFalse(marker.exists())
                self.assertIsNone(process.poll())
                bound.append(process.pid)

            with mock.patch.object(bot, "jarvis_cmd",
                                   return_value=[sys.executable, "-c", script]), \
                    mock.patch.object(bot, "jarvis_root", return_value=Path(directory)), \
                    mock.patch.object(bot, "_headless_exec_command",
                                      side_effect=lambda _sid, command: command):
                result = bot.run_claude_buffered(
                    "prompt", "session", False, timeout=5,
                    on_spawn=bind, guarded=True)

            self.assertEqual(result.text, "ok")
            self.assertFalse(result.is_error)
            self.assertEqual(len(bound), 1)
            self.assertTrue(marker.exists())

    def test_failed_fenced_bind_never_launches_command(self):
        with tempfile.TemporaryDirectory() as directory:
            marker = Path(directory) / "must-not-start"
            script = "from pathlib import Path; Path(%r).write_text('bad')" % str(marker)

            def reject(_process):
                raise RuntimeError("stale fence")

            with mock.patch.object(bot, "jarvis_cmd",
                                   return_value=[sys.executable, "-c", script]), \
                    mock.patch.object(bot, "jarvis_root", return_value=Path(directory)), \
                    mock.patch.object(bot, "_headless_exec_command",
                                      side_effect=lambda _sid, command: command):
                with self.assertRaisesRegex(RuntimeError, "stale fence"):
                    bot.run_claude_buffered(
                        "prompt", "session", False, timeout=5,
                        on_spawn=reject, guarded=True)

            time.sleep(0.1)
            self.assertFalse(marker.exists())

    def test_post_pr_claim_failure_never_launches_or_advances_cursor(self):
        with tempfile.TemporaryDirectory() as directory:
            marker = Path(directory) / "must-not-start-post-pr"
            cursor = Path(directory) / "cursor"
            script = "from pathlib import Path; Path(%r).write_text('bad')" % str(marker)

            class Pool:
                process = None

                def set_proc(self, _item_id, process):
                    self.process = process

            pool = Pool()
            binder = bot._post_pr_process_binder(
                pool, "84362517", "pr_comment_reply", "528766",
                "session-post-pr", "prompt",
                on_claimed=lambda: cursor.write_text("advanced"))
            with mock.patch.object(bot, "jarvis_cmd",
                                   return_value=[sys.executable, "-c", script]), \
                    mock.patch.object(bot, "jarvis_root", return_value=Path(directory)), \
                    mock.patch.object(bot, "_headless_exec_command",
                                      side_effect=lambda _sid, command: command), \
                    mock.patch.object(bot, "_claim_workitem",
                                      side_effect=RuntimeError("claim lost")), \
                    mock.patch.object(bot, "_release_post_pr_claim") as release, \
                    mock.patch.object(bot, "_inflight_add"), \
                    mock.patch.object(bot, "_inflight_remove"):
                with self.assertRaisesRegex(RuntimeError, "claim lost"):
                    bot.run_claude_buffered(
                        "prompt", "session-post-pr", False, timeout=5,
                        on_spawn=binder, guarded=True, terraform=True,
                        aone_write_policy=bot.POST_PR_AONE_WRITE_POLICY)

            time.sleep(0.1)
            self.assertIsNotNone(pool.process, "set_proc must happen before claim")
            self.assertFalse(marker.exists(), "guard must kill before command exec")
            self.assertFalse(cursor.exists(), "claim failure must not advance PR cursor")
            release.assert_not_called()

    def test_legacy_on_spawn_is_also_gated_before_command(self):
        with tempfile.TemporaryDirectory() as directory:
            marker = Path(directory) / "must-not-start-legacy"
            script = "from pathlib import Path; Path(%r).write_text('bad')" % str(marker)

            def reject(_process):
                raise RuntimeError("pool registration failed")

            with mock.patch.object(bot, "jarvis_cmd",
                                   return_value=[sys.executable, "-c", script]), \
                    mock.patch.object(bot, "jarvis_root", return_value=Path(directory)), \
                    mock.patch.object(bot, "_headless_exec_command",
                                      side_effect=lambda _sid, command: command):
                with self.assertRaisesRegex(RuntimeError, "pool registration failed"):
                    bot.run_claude_buffered(
                        "prompt", "session", False, timeout=5,
                        on_spawn=reject)

            time.sleep(0.1)
            self.assertFalse(marker.exists())

    def test_normal_command_exit_kills_background_grandchildren(self):
        with tempfile.TemporaryDirectory() as directory:
            grandchild_pid_file = Path(directory) / "grandchild.pid"
            grandchild = (
                "import signal,time; "
                "signal.signal(signal.SIGTERM, signal.SIG_IGN); "
                "time.sleep(60)"
            )
            command = (
                "import json,subprocess,sys; from pathlib import Path; "
                "p=subprocess.Popen([sys.executable,'-c',%r],"
                "stdin=subprocess.DEVNULL,stdout=subprocess.DEVNULL,"
                "stderr=subprocess.DEVNULL); "
                "Path(%r).write_text(str(p.pid)); "
                "print(json.dumps({'type':'result','result':'ok',"
                "'is_error':False,'subtype':'success'}),flush=True)"
            ) % (grandchild, str(grandchild_pid_file))
            with mock.patch.object(bot, "jarvis_cmd",
                                   return_value=[sys.executable, "-c", command]), \
                    mock.patch.object(bot, "jarvis_root", return_value=Path(directory)), \
                    mock.patch.object(bot, "_headless_exec_command",
                                      side_effect=lambda _sid, argv: argv), \
                    mock.patch.dict(os.environ,
                                    {"JARVIS_MANAGED_GUARD_GRACE_SEC": "0.1"}):
                result = bot.run_claude_buffered(
                    "prompt", "session", False, timeout=5,
                    on_spawn=lambda _process: None, guarded=True)

            self.assertEqual(result.text, "ok")
            self.assertTrue(grandchild_pid_file.exists())
            grandchild_pid = int(grandchild_pid_file.read_text())
            deadline = time.monotonic() + 3
            while _pid_alive(grandchild_pid) and time.monotonic() < deadline:
                time.sleep(0.05)
            self.assertFalse(_pid_alive(grandchild_pid))

    def test_sigkill_of_bridge_parent_kills_term_deaf_child_group(self):
        with tempfile.TemporaryDirectory() as directory:
            child_pid_file = Path(directory) / "child.pid"
            child = (
                "import os,signal,time; from pathlib import Path; "
                "signal.signal(signal.SIGTERM, signal.SIG_IGN); "
                "Path(%r).write_text(str(os.getpid())); "
                "time.sleep(60)"
            ) % str(child_pid_file)
            parent = (
                "import os,subprocess,sys,time; "
                "gate_r,gate_w=os.pipe(); sent_r,sent_w=os.pipe(); "
                "p=subprocess.Popen([sys.executable,%r,'--gate-fd',str(gate_r),"
                "'--sentinel-fd',str(sent_r),'--grace-seconds','0.1','--',"
                "sys.executable,'-c',%r],pass_fds=(gate_r,sent_r),"
                "stdin=subprocess.DEVNULL,stdout=subprocess.DEVNULL,"
                "stderr=subprocess.DEVNULL,start_new_session=True); "
                "os.close(gate_r); os.close(sent_r); os.write(gate_w,b'G'); "
                "os.close(gate_w); print(p.pid,flush=True); time.sleep(60)"
            ) % (str(GUARD), child)
            owner = subprocess.Popen(
                [sys.executable, "-c", parent], stdout=subprocess.PIPE,
                stderr=subprocess.PIPE, text=True, start_new_session=True)
            guard_pid = None
            child_pid = None
            try:
                guard_pid = int(owner.stdout.readline().strip())
                deadline = time.monotonic() + 5
                while not child_pid_file.exists() and time.monotonic() < deadline:
                    time.sleep(0.02)
                self.assertTrue(child_pid_file.exists())
                child_pid = int(child_pid_file.read_text())
                self.assertTrue(_pid_alive(guard_pid))
                self.assertTrue(_pid_alive(child_pid))

                os.kill(owner.pid, signal.SIGKILL)
                owner.wait(timeout=5)
                deadline = time.monotonic() + 5
                while ((_pid_alive(guard_pid) or _pid_alive(child_pid))
                       and time.monotonic() < deadline):
                    time.sleep(0.05)
                self.assertFalse(_pid_alive(guard_pid))
                self.assertFalse(_pid_alive(child_pid))
            finally:
                for pid in (child_pid, guard_pid):
                    if pid and _pid_alive(pid):
                        try:
                            os.kill(pid, signal.SIGKILL)
                        except ProcessLookupError:
                            pass
                if owner.poll() is None:
                    os.kill(owner.pid, signal.SIGKILL)
                    owner.wait(timeout=5)
                if owner.stdout is not None:
                    owner.stdout.close()
                if owner.stderr is not None:
                    owner.stderr.close()


if __name__ == "__main__":
    unittest.main()
