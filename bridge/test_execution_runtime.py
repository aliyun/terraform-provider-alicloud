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

from datetime import datetime as real_dt, timezone as tz_module  # noqa: E402
tz_utc = tz_module.utc

from bridge.jarvis_execution_runtime import (  # noqa: E402
    ExecutionRuntime, ProcessGuardian, _select_provider_settings)


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
        with mock.patch("bridge.jarvis_execution_runtime.subprocess.Popen",
                        return_value=process), \
                mock.patch("bridge.jarvis_execution_runtime.os.write",
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

        with mock.patch("bridge.jarvis_execution_runtime.subprocess.Popen",
                        return_value=process) as popen:
            result = ExecutionRuntime().run_buffered(
                ["tool", "arg"], HERE, timeout=10,
                on_spawn=lambda value: spawned.append(value.pid),
                env={"TEST_EXECUTION_ENV": "present"})

        self.assertEqual((result.stdout, result.stderr, result.returncode),
                         ("out", "err", 7))
        self.assertFalse(result.timed_out)
        self.assertEqual(spawned, [42])
        self.assertEqual(popen.call_args.args[0], ["tool", "arg"])
        self.assertEqual(
            popen.call_args.kwargs["env"], {"TEST_EXECUTION_ENV": "present"})

    def test_timeout_kills_process_group(self):
        process = mock.Mock()
        process.pid = 99
        process.returncode = -9
        process.communicate.side_effect = [
            subprocess.TimeoutExpired("tool", 1),
            ("partial", "timeout"),
        ]

        with mock.patch("bridge.jarvis_execution_runtime.subprocess.Popen",
                        return_value=process), \
                mock.patch("bridge.jarvis_execution_runtime.os.getpgid", return_value=99), \
                mock.patch("bridge.jarvis_execution_runtime.os.killpg") as killpg:
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

        def guarded_spawn(argv, cwd, on_spawn, env):
            captured["argv"] = list(argv)
            captured["cwd"] = cwd
            captured["env"] = env
            on_spawn(process)
            return process, sentinel_write

        result = ExecutionRuntime().run_buffered(
            ["tool"], HERE, timeout=1, guarded=True,
            guarded_spawn=guarded_spawn, on_spawn=lambda _process: None,
            env={"TASK_ENV": "fenced"})

        self.assertEqual(result.stdout, "ok")
        self.assertEqual(captured, {
            "argv": ["tool"], "cwd": HERE, "env": {"TASK_ENV": "fenced"}})


class ProviderResumeFailoverTest(unittest.TestCase):
    """Resume no longer dies on a transient provider blip."""

    MEM = "/path/a.json,/path/b.json,/path/c.json"
    SID = "sess-failover-1"

    def setUp(self):
        self.tmp = self.tmpdir()
        self.addCleanup(self.tmp.cleanup)
        self.route_dir = self.tmp.name
        self._patch_env = mock.patch.dict(
            os.environ, {"JARVIS_PROVIDER_ROUTE_DIR": self.route_dir})
        self._patch_env.start()
        self.addCleanup(self._patch_env.stop)
        # model stub: a.json/b.json are same family, c.json is a different one
        self._models = {"/path/a.json": "qwen-3.7", "/path/b.json": "qwen-3.7-max",
                        "/path/c.json": "glm-5.2"}
        self.probes = {}  # path -> bool; default True

    def tmpdir(self):
        import tempfile
        return tempfile.TemporaryDirectory()

    def _probe(self, path):
        return self.probes.get(path, True)

    def _settings_model(self, path):
        return self._models.get(path)

    def _route_record(self, path, **extra):
        return {"schemaVersion": 1, "sessionId": self.SID, "lane": "terraform",
                "settingsPath": path, "model": self._models.get(path),
                "selectedAt": "2026-07-01T00:00:00+00:00", **extra}

    def _write_route(self, **extra):
        from bridge.jarvis_execution_runtime import _provider_route_file
        import json
        f = _provider_route_file(self.SID)
        f.parent.mkdir(parents=True, exist_ok=True)
        record = self._route_record("/path/a.json", **extra)
        f.write_text(json.dumps(record))

    def test_resume_healthy_provider_unchanged(self):
        self._write_route()
        selected = _select_provider_settings(
            self.MEM, self.SID, True, True, probe_settings=self._probe)
        self.assertEqual(selected, "/path/a.json")

    def test_first_probe_miss_returns_original_and_arms(self):
        """A blip gets one backoff window before failover, not a RuntimeError."""
        self.probes = {"/path/a.json": False}
        self._write_route()
        selected = _select_provider_settings(
            self.MEM, self.SID, True, True, probe_settings=self._probe)
        self.assertEqual(selected, "/path/a.json", "original provider first")
        # firstFailedAt must now be pinned so the backoff clock starts
        from bridge.jarvis_execution_runtime import _load_route_record
        record = _load_route_record(self.SID, "terraform")
        self.assertIn("firstFailedAt", record)

    def test_backoff_expired_fails_over_to_same_family(self):
        self.probes = {"/path/a.json": False, "/path/b.json": True}
        # firstFailedAt 90s ago — past the 60s window
        self._write_route(firstFailedAt=(
            "2026-07-01T00:00:00+00:00"))
        with mock.patch(
                "bridge.jarvis_execution_runtime.datetime") as dt:
            dt.now.return_value = real_dt(2026, 7, 1, 0, 1, 30, tzinfo=tz_utc)
            dt.fromisoformat = real_dt.fromisoformat
            selected = _select_provider_settings(
                self.MEM, self.SID, True, True, probe_settings=self._probe)
        self.assertEqual(selected, "/path/b.json", "same family preferred")
        # route file updated to the new provider, failoverFrom recorded
        from bridge.jarvis_execution_runtime import _load_route_record
        record = _load_route_record(self.SID, "terraform")
        self.assertEqual(record["settingsPath"], "/path/b.json")
        self.assertEqual(record["failoverFrom"]["settingsPath"], "/path/a.json")

    def test_no_healthy_candidate_raises(self):
        self.probes = {"/path/a.json": False, "/path/b.json": False,
                       "/path/c.json": False}
        self._write_route(firstFailedAt="2026-07-01T00:00:00+00:00")
        with mock.patch(
                "bridge.jarvis_execution_runtime.datetime") as dt:
            dt.now.return_value = real_dt(2026, 7, 1, 0, 5, 0, tzinfo=tz_utc)
            dt.fromisoformat = real_dt.fromisoformat
            with self.assertRaises(RuntimeError) as ctx:
                _select_provider_settings(
                    self.MEM, self.SID, True, True, probe_settings=self._probe)
        self.assertIn("no healthy candidate", str(ctx.exception))

    def test_new_session_uses_resolve_settings(self):
        """The new-session path is unchanged: probe and pick a healthy one."""
        self.probes = {"/path/a.json": False, "/path/b.json": True, "/path/c.json": False}
        with mock.patch.multiple(
                "bridge.jarvis_execution_runtime",
                _probe_settings=self._probe,
                _persist_provider_route=mock.Mock(return_value=True)):
            selected = _select_provider_settings(
                self.MEM, self.SID, True, False, probe_settings=self._probe)
        self.assertEqual(selected, "/path/b.json")


if __name__ == "__main__":
    unittest.main()
