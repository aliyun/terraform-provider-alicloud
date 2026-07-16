#!/usr/bin/env python3
"""Hermetic tests for the shared execution runtime."""

import subprocess
import sys
import unittest
from pathlib import Path
from unittest import mock

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from jarvis_execution_runtime import ExecutionRuntime  # noqa: E402


class ExecutionRuntimeTest(unittest.TestCase):
    def test_captures_process_result_and_reports_spawn(self):
        process = mock.Mock()
        process.pid = 42
        process.returncode = 7
        process.communicate.return_value = ("out", "err")
        spawned = []

        with mock.patch("jarvis_execution_runtime.subprocess.Popen",
                        return_value=process) as popen:
            result = ExecutionRuntime().run_buffered(
                ["tool", "arg"], HERE, timeout=10,
                on_spawn=lambda value: spawned.append(value.pid))

        self.assertEqual((result.stdout, result.stderr, result.returncode),
                         ("out", "err", 7))
        self.assertFalse(result.timed_out)
        self.assertEqual(spawned, [42])
        self.assertEqual(popen.call_args.args[0], ["tool", "arg"])

    def test_timeout_kills_process_group(self):
        process = mock.Mock()
        process.pid = 99
        process.returncode = -9
        process.communicate.side_effect = [
            subprocess.TimeoutExpired("tool", 1),
            ("partial", "timeout"),
        ]

        with mock.patch("jarvis_execution_runtime.subprocess.Popen",
                        return_value=process), \
                mock.patch("jarvis_execution_runtime.os.getpgid", return_value=99), \
                mock.patch("jarvis_execution_runtime.os.killpg") as killpg:
            result = ExecutionRuntime().run_buffered(
                ["tool"], HERE, timeout=1)

        self.assertTrue(result.timed_out)
        self.assertEqual(result.stdout, "partial")
        killpg.assert_called_once()

    def test_guarded_spawn_is_injected_and_sentinel_is_closed(self):
        process = mock.Mock()
        process.pid = 123
        process.returncode = 0
        process.communicate.return_value = ("ok", "")
        sentinel_read, sentinel_write = __import__("os").pipe()
        __import__("os").close(sentinel_read)
        captured = {}

        def guarded_spawn(argv, cwd, on_spawn):
            captured["argv"] = list(argv)
            captured["cwd"] = cwd
            on_spawn(process)
            return process, sentinel_write

        result = ExecutionRuntime().run_buffered(
            ["tool"], HERE, timeout=1, guarded=True,
            guarded_spawn=guarded_spawn, on_spawn=lambda _process: None)

        self.assertEqual(result.stdout, "ok")
        self.assertEqual(captured, {"argv": ["tool"], "cwd": HERE})


if __name__ == "__main__":
    unittest.main()
