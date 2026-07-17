#!/usr/bin/env python3
"""Hermetic tests for the shared execution runtime."""

import subprocess
import sys
import unittest
from pathlib import Path
from unittest import mock

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

import os  # noqa: E402

from jarvis_execution_runtime import ExecutionRuntime, ProcessGuardian  # noqa: E402


class ProcessGuardianTest(unittest.TestCase):
    def test_gate_broken_pipe_is_swallowed_not_terminal(self):
        # The guard child can exit before the gate is granted when the bound
        # session lost ownership (stale_fence) and stop_process fired mid-spawn.
        # os.write then hits BrokenPipe — an expected handoff, not a crash: the
        # process is returned so run_buffered observes a retryable empty result
        # instead of a terminal orchestrator_exception.
        process = mock.Mock()
        process.pid = 321
        bound = []
        with mock.patch("jarvis_execution_runtime.subprocess.Popen",
                        return_value=process), \
                mock.patch("jarvis_execution_runtime.os.write",
                           side_effect=BrokenPipeError(32, "Broken pipe")), \
                mock.patch.object(ProcessGuardian, "terminate") as terminate:
            proc, sentinel_write = ProcessGuardian().spawn(
                ["tool"], HERE, on_spawn=lambda p: bound.append(p.pid))

        self.assertIs(proc, process)
        self.assertEqual(bound, [321], "fence must bind before the gate is granted")
        terminate.assert_not_called()
        os.close(sentinel_write)


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
