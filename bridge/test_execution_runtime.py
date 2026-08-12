#!/usr/bin/env python3
"""Hermetic tests for the shared execution runtime."""

import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

import os  # noqa: E402

from datetime import datetime as real_dt, timezone as tz_module  # noqa: E402
tz_utc = tz_module.utc

from bridge.jarvis_execution_runtime import (  # noqa: E402
    ExecutionRuntime, HeadlessAmpBroker, ProcessGuardian,
    _select_provider_settings)


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


class HeadlessAmpBrokerTest(unittest.TestCase):
    def setUp(self):
        self.controller = mock.Mock()
        self.controller.verify_current_fence.return_value = True
        self.broker = HeadlessAmpBroker(self.controller, "runtime-1")
        self.addCleanup(self.broker.close)
        self.process = mock.Mock(pid=os.getpid())

    def _request(self, repo=None, action="publish", peer=None, *,
                 project_id="3065873", repo_mode="model",
                 branch="feature/guard", publish_kind="pre"):
        resolved = Path(repo or HERE).resolve()
        project = {"projectId": project_id, "popCode": "eventbridge",
                   "popVersion": "2020-04-01"}
        target = {
            "schemaVersion": 1, "action": action, "repoMode": repo_mode,
            "project": project, "branch": branch,
            "publishKind": publish_kind if action == "publish" else None,
        }
        receipt = {
            "allowed": True, "action": action, "repoMode": repo_mode,
            "project": project,
        }
        if repo_mode == "model":
            receipt.update(repo=str(resolved), baseline="a" * 40,
                           branch=branch)
        return self.broker._authorize(
            os.getpid() if peer is None else peer,
            {"schemaVersion": 2, "action": action, "repo": str(resolved),
             "target": target, "receipt": receipt})

    def test_disabled_until_full_bind_callback_enables_it(self):
        with mock.patch("bridge.jarvis_execution_runtime._descends_from",
                        return_value=True), \
                mock.patch.object(self.broker, "_trusted_wrapper_process",
                                  return_value=True):
            self.assertFalse(self._request())
            self.broker.enable(self.process)
            self.assertTrue(self._request())

    def test_requires_descendant_and_exact_trusted_wrapper(self):
        self.broker.enable(self.process)
        with mock.patch("bridge.jarvis_execution_runtime._descends_from",
                        return_value=False), \
                mock.patch.object(self.broker, "_trusted_wrapper_process",
                                  return_value=True):
            self.assertFalse(self._request())
        with mock.patch("bridge.jarvis_execution_runtime._descends_from",
                        return_value=True), \
                mock.patch.object(self.broker, "_trusted_wrapper_process",
                                  return_value=False):
            self.assertFalse(self._request())

    def test_binds_first_repository_inode_and_rechecks_fence(self):
        self.broker.enable(self.process)
        with tempfile.TemporaryDirectory() as first, \
                tempfile.TemporaryDirectory() as second, \
                mock.patch("bridge.jarvis_execution_runtime._descends_from",
                           return_value=True), \
                mock.patch.object(self.broker, "_trusted_wrapper_process",
                                  return_value=True):
            self.assertTrue(self._request(first))
            self.assertFalse(self._request(second))
            self.controller.verify_current_fence.return_value = False
            self.assertFalse(self._request(first))

    def test_invalid_first_repository_does_not_poison_later_canonical_pin(self):
        self.broker.enable(self.process)
        with tempfile.TemporaryDirectory() as wrong, \
                tempfile.TemporaryDirectory() as canonical, \
                mock.patch("bridge.jarvis_execution_runtime._descends_from",
                           return_value=True), \
                mock.patch.object(self.broker, "_trusted_wrapper_process",
                                  return_value=True):
            request = {
                "schemaVersion": 2, "action": "publish", "repo": wrong,
                "target": {
                    "schemaVersion": 1, "action": "publish",
                    "repoMode": "model", "branch": "feature/guard",
                    "publishKind": "pre",
                    "project": {"projectId": "3065873",
                                "popCode": "eventbridge",
                                "popVersion": "2020-04-01"}},
                "receipt": {
                    "allowed": True, "action": "publish", "repoMode": "model",
                    "repo": canonical, "baseline": "a" * 40,
                    "branch": "feature/guard",
                    "project": {"projectId": "3065873",
                                "popCode": "eventbridge",
                                "popVersion": "2020-04-01"}},
            }
            decision = self.broker._authorize_decision(os.getpid(), request)
            self.assertFalse(decision.allowed)
            self.assertEqual(decision.reason,
                             "canonical_baseline_receipt_invalid")
            self.assertIsNone(self.broker._repo_identity)
            self.assertTrue(self._request(canonical))

    def test_cross_project_remote_mutation_is_denied_after_project_pin(self):
        self.broker.enable(self.process)
        with tempfile.TemporaryDirectory() as repo, \
                mock.patch("bridge.jarvis_execution_runtime._descends_from",
                           return_value=True), \
                mock.patch.object(self.broker, "_trusted_wrapper_process",
                                  return_value=True):
            self.assertTrue(self._request(repo))
            decision = self.broker._authorize_decision(
                os.getpid(), self._request_payload(
                    repo, project_id="3033394"))
            self.assertFalse(decision.allowed)
            self.assertEqual(decision.reason, "project_identity_mismatch")

    def test_mutations_pin_feature_branch_for_the_session(self):
        self.broker.enable(self.process)
        with tempfile.TemporaryDirectory() as repo, \
                mock.patch("bridge.jarvis_execution_runtime._descends_from",
                           return_value=True), \
                mock.patch.object(self.broker, "_trusted_wrapper_process",
                                  return_value=True):
            self.assertTrue(self._request(repo, branch="feature/one"))
            decision = self.broker._authorize_decision(
                os.getpid(), self._request_payload(repo, branch="feature/two"))
            self.assertFalse(decision.allowed)
            self.assertEqual(decision.reason, "branch_identity_mismatch")
            self.assertEqual(self.broker._branch_identity, "feature/one")

    def test_init_and_publish_status_never_establish_mutation_pins(self):
        self.broker.enable(self.process)
        with tempfile.TemporaryDirectory() as repo, \
                mock.patch("bridge.jarvis_execution_runtime._descends_from",
                           return_value=True), \
                mock.patch.object(self.broker, "_trusted_wrapper_process",
                                  return_value=True):
            resolved = str(Path(repo).resolve())
            init_project = {"popCode": "eventbridge",
                            "popVersion": "2020-04-01"}
            init_request = {
                "schemaVersion": 2, "action": "init", "repo": resolved,
                "target": {"schemaVersion": 1, "action": "init",
                           "repoMode": "local-context",
                           "project": init_project, "branch": None},
                "receipt": {"allowed": True, "action": "init",
                            "repoMode": "local-context",
                            "project": init_project},
            }
            self.assertTrue(self.broker._authorize(
                os.getpid(), init_request))
            self.assertTrue(self._request(
                repo, publish_kind="status", repo_mode="project", branch=""))
            self.assertIsNone(self.broker._repo_identity)
            self.assertIsNone(self.broker._project_identity)
            self.assertIsNone(self.broker._branch_identity)

    def _request_payload(self, repo, *, project_id="3065873",
                         branch="feature/guard"):
        resolved = str(Path(repo).resolve())
        project = {"projectId": project_id, "popCode": "eventbridge",
                   "popVersion": "2020-04-01"}
        return {
            "schemaVersion": 2, "action": "publish", "repo": resolved,
            "target": {"schemaVersion": 1, "action": "publish",
                       "repoMode": "model", "branch": branch,
                       "publishKind": "pre", "project": project},
            "receipt": {"allowed": True, "action": "publish",
                        "repoMode": "model", "repo": resolved,
                        "baseline": "a" * 40, "branch": branch,
                        "project": project},
        }


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
