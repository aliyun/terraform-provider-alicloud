#!/usr/bin/env python3
"""Regression tests for bridge-spawned Claude worker pre-exec fencing."""

import io
import json
import os
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest import mock

import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))

import jarvis_dingtalk_bot as bot


class _BufferedProcess:
    _next_pid = 4000

    def __init__(self, result):
        self.result = result
        self.returncode = 0
        self.stdout = None
        self.stderr = None
        self.pid = self.__class__._next_pid
        self.__class__._next_pid += 1

    def communicate(self, timeout=None):
        return json.dumps({
            "type": "result",
            "result": self.result.text,
            "is_error": self.result.is_error,
            "subtype": self.result.subtype,
        }), ""

    def poll(self):
        return None


class _StreamProcess:
    def __init__(self):
        self.pid = 5001
        self.stdout = iter([
            json.dumps({
                "type": "result",
                "result": "stream-ok",
                "is_error": False,
                "subtype": "success",
            }) + "\n",
        ])
        self.stderr = io.StringIO("")
        self.returncode = 0

    def wait(self, timeout=None):
        return self.returncode

    def kill(self):
        self.returncode = -9


class HeadlessWorkerFenceTest(unittest.TestCase):
    def test_wrapper_uses_fixed_isolated_manager(self):
        wrapped = bot._headless_exec_command(
            "runtime-session", ["/opt/claude", "--resume", "runtime-session"])
        self.assertEqual(wrapped[:2], ["/usr/bin/python3", "-I"])
        self.assertEqual(Path(wrapped[2]).resolve(),
                         (Path(bot.__file__).resolve().parents[1] /
                          "bootstrap" / "jarvis-interactive-worker.py"))
        self.assertEqual(wrapped[3:8], [
            "exec-headless", "--session-id", "runtime-session",
            "--client", "claude",
        ])
        self.assertEqual(wrapped[8], "--")
        self.assertEqual(wrapped[9:],
                         ["/opt/claude", "--resume", "runtime-session"])

    def test_real_wrapper_publishes_same_pid_state_before_exec(self):
        with tempfile.TemporaryDirectory() as directory:
            script = (
                "import glob,json,os; "
                "paths=glob.glob(%r+'/*.json'); "
                "state=json.load(open(paths[0])); "
                "assert state['hostPid']==os.getpid(); "
                "assert state['headlessRegistered']; "
                "print(str(os.getpid())+' '+state['workerKey'])"
            ) % directory
            command = bot._headless_exec_command(
                "wrapper-integration",
                ["/usr/bin/python3", "-c", script])
            env = os.environ.copy()
            env.update({
                "JARVIS_INTERACTIVE_STATE_DIR": directory,
                "JARVIS_CONTROL_PLANE_BASE_URL": "http://127.0.0.1:1",
                "JARVIS_HEADLESS_REMOTE_REGISTER_TIMEOUT": "0.05",
            })
            completed = subprocess.run(
                command, cwd=directory, env=env,
                capture_output=True, text=True, timeout=5)

            self.assertEqual(completed.returncode, 0, completed.stderr)
            pid, worker_key = completed.stdout.strip().split(" ", 1)
            state_path = next(Path(directory).glob("*.json"))
            state = json.loads(state_path.read_text())
            self.assertEqual(state["hostPid"], int(pid))
            self.assertEqual(state["workerKey"], worker_key)

    def test_stream_path_spawns_only_the_preexec_wrapper(self):
        process = _StreamProcess()
        with mock.patch.object(bot, "jarvis_cmd",
                               return_value=["/opt/claude"]), \
                mock.patch.object(bot.subprocess, "Popen",
                                  return_value=process) as popen:
            output = list(bot.run_claude_stream(
                "prompt", "runtime-session", True, timeout=5))

        self.assertEqual(output[-1], "stream-ok")
        command = popen.call_args.args[0]
        self.assertEqual(command[:4], [
            "/usr/bin/python3", "-I",
            str(Path(bot.__file__).resolve().parents[1] /
                "bootstrap" / "jarvis-interactive-worker.py"),
            "exec-headless",
        ])
        self.assertIn("--resume", command)

    def test_every_dispatch_retry_spawns_a_fresh_preexec_wrapper(self):
        processes = [
            _BufferedProcess(bot.ClaudeResult(
                "transient output", True, "overloaded_error")),
            _BufferedProcess(bot.ClaudeResult(
                "completed", False, "success")),
        ]
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._maybe_suspend = mock.Mock(return_value=None)
        handler._completion_broadcast = mock.Mock(return_value="done")
        handler._dispatch_failed = mock.Mock()
        notifications = []
        with mock.patch.dict(os.environ, {
            "JARVIS_DISPATCH_RETRY_MAX": "1",
            "JARVIS_DISPATCH_RETRY_BACKOFF": "0",
        }, clear=False), \
                mock.patch.object(bot, "jarvis_cmd",
                                  return_value=["/opt/claude"]), \
                mock.patch.object(bot.subprocess, "Popen",
                                  side_effect=processes) as popen, \
                mock.patch.object(bot, "_inflight_add"), \
                mock.patch.object(bot, "_inflight_remove"), \
                mock.patch.object(bot.time, "sleep"):
            result = handler.dispatch_item(
                "84382166", "prompt", "runtime-session", False,
                notifications.append, "target", "staff",
                project="2100304")

        self.assertEqual(result, "done")
        self.assertEqual(popen.call_count, 2)
        commands = [call.args[0] for call in popen.call_args_list]
        self.assertTrue(all(command[3] == "exec-headless"
                            for command in commands))
        self.assertIn("--session-id", commands[0])
        self.assertIn("--resume", commands[1])
        self.assertEqual(commands[0][7], "claude")
        self.assertEqual(commands[1][7], "claude")

    def test_managed_guard_binds_guard_but_wraps_the_claude_child(self):
        captured = {}
        guard = _BufferedProcess(bot.ClaudeResult("managed-ok", False, "success"))
        bound = []

        def spawn(argv, cwd, on_spawn):
            captured["argv"] = list(argv)
            on_spawn(guard)
            return guard, None

        with mock.patch.object(bot, "jarvis_cmd",
                               return_value=["/opt/claude"]), \
                mock.patch.object(bot, "_spawn_guarded_managed_process",
                                  side_effect=spawn):
            result = bot.run_claude_buffered(
                "prompt", "managed-runtime", True, timeout=5,
                on_spawn=lambda process: bound.append(process.pid),
                guarded=True)

        self.assertEqual(result.text, "managed-ok")
        self.assertEqual(bound, [guard.pid])
        self.assertEqual(captured["argv"][3], "exec-headless")
        self.assertIn("/opt/claude", captured["argv"])
        self.assertNotEqual(bound[0], os.getpid())


if __name__ == "__main__":
    unittest.main(verbosity=2)
