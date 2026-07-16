#!/usr/bin/env python3
"""Hermetic integration tests for the bridge/control-plane seam."""

import os
import signal
import subprocess
import sys
import threading
import time
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

sys.path.insert(0, str(Path(__file__).resolve().parent))

import jarvis_dingtalk_bot as bot
from jarvis_persistence_executor import SessionController
from jarvis_task_router import EnqueueResult


class _Starter:
    def __init__(self, name, calls):
        self.name = name
        self.calls = calls

    def start(self):
        self.calls.append(self.name)


class _FakePersistenceExecutor:
    instances = []

    def __init__(self, *args, **kwargs):
        self.args = args
        self.kwargs = kwargs
        self.stopped = False
        self.stop_calls = []
        self.__class__.instances.append(self)

    def start(self):
        return self

    def stop(self, *, drain=False, timeout=None):
        self.stop_calls.append((drain, timeout))
        self.stopped = True
        return True


class HandlerWiringTest(unittest.TestCase):
    ENV_KEYS = (
        "JARVIS_CONTROL_PLANE_BASE_URL", "JARVIS_CONTROL_PLANE_TOKEN",
        "JARVIS_HTML_REPORT_BASE_URL",
        "JARVIS_HTML_REPORT_TOKEN",
    )

    def setUp(self):
        self.old_env = {key: os.environ.get(key) for key in self.ENV_KEYS}
        _FakePersistenceExecutor.instances.clear()

    def tearDown(self):
        for key, value in self.old_env.items():
            if value is None:
                os.environ.pop(key, None)
            else:
                os.environ[key] = value

    def test_worker_capability_matches_fixed_task_types(self):
        client = object()
        with mock.patch.object(bot, "_task_client_from_env", return_value=client), \
                mock.patch.object(bot, "PersistenceExecutor", _FakePersistenceExecutor):
            handler = bot.JarvisHandler(no_dingtalk=True)
        self.assertIn("ticket", handler.execution_router.task_types)
        self.assertIn("wake", handler.execution_router.task_types)
        self.assertNotIn("probe", handler.execution_router.task_types)
        self.assertIsNotNone(handler.persistence_executor)
        self.assertEqual(_FakePersistenceExecutor.instances[-1].kwargs["capabilities"],
                         {"kinds": sorted(handler.execution_router.task_types)})
        self.assertEqual(
            _FakePersistenceExecutor.instances[-1].kwargs["lease_safety_margin"], 90)
        self.assertIs(_FakePersistenceExecutor.instances[-1].args[1],
                      handler.ephemeral_executor.capacity_manager)
        self.assertIs(handler.execution_runtime,
                      handler.ephemeral_executor.execution_runtime)
        for obsolete in ("task_router", "local_worker", "dispatch_pool"):
            self.assertFalse(hasattr(handler, obsolete))

    def test_worker_starts_before_every_sensor(self):
        calls = []
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.persistence_executor = _Starter("worker", calls)
        for name in ("scanner", "reconciler", "board", "prober", "reviser",
                     "watcher", "personawatch", "prwatch"):
            setattr(handler, name, _Starter(name, calls))
        handler.start_schedulers()
        self.assertEqual(calls[0], "worker")
        self.assertEqual(calls[-1], "prwatch")

    def test_stop_helper_forwards_drain_policy(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        worker = _FakePersistenceExecutor()
        handler.persistence_executor = worker
        self.assertTrue(handler.stop_persistence_executor(drain=True, timeout=7))
        self.assertEqual(worker.stop_calls, [(True, 7)])

    def test_task_client_defaults_to_pre_and_requires_token(self):
        for key in self.ENV_KEYS:
            os.environ.pop(key, None)
        with self.assertRaisesRegex(RuntimeError, "token is required"):
            bot._task_client_from_env()
        os.environ["JARVIS_HTML_REPORT_TOKEN"] = "shared-token"
        client = bot._task_client_from_env()
        self.assertEqual(client.base_url, "https://pre-agent.aliyun-inc.com")
        self.assertEqual(client.token, "shared-token")


class TaskExecutionTest(unittest.TestCase):
    def test_process_spawned_after_fence_loss_is_immediately_force_killed(self):
        lifecycle = SessionController(
            object(), "mac:boot:process",
            {"task": {"id": "task-1"},
             "session": {"id": "late-bind-session", "fenceToken": 7}},
            stop_process=bot.JarvisHandler._stop_task_process)
        # First stop attempt happens before on_spawn has supplied a process.
        lifecycle._lose_ownership("stale_fence:heartbeat")

        proc = subprocess.Popen(
            [sys.executable, "-c", (
                "import signal,time; "
                "signal.signal(signal.SIGTERM, signal.SIG_IGN); "
                "print('ready', flush=True); time.sleep(60)"
            )],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            start_new_session=True,
        )
        try:
            self.assertEqual(proc.stdout.readline().strip(), "ready")
            with mock.patch.dict(
                    os.environ, {"JARVIS_TASK_STOP_GRACE_SEC": "0.05"}):
                lifecycle.bind_process(proc)
            self.assertEqual(proc.poll(), -signal.SIGKILL)
        finally:
            if proc.poll() is None:
                os.killpg(os.getpgid(proc.pid), signal.SIGKILL)
                proc.wait(timeout=5)
            if proc.stdout is not None:
                proc.stdout.close()
            if proc.stderr is not None:
                proc.stderr.close()

    def test_fence_loss_force_kills_process_group_that_ignores_sigterm(self):
        proc = subprocess.Popen(
            [sys.executable, "-c", (
                "import signal,time; "
                "signal.signal(signal.SIGTERM, signal.SIG_IGN); "
                "print('ready', flush=True); time.sleep(60)"
            )],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            start_new_session=True,
        )
        try:
            self.assertEqual(proc.stdout.readline().strip(), "ready")
            lifecycle = SimpleNamespace(process=proc, session_id="fenced-session")
            with mock.patch.dict(
                    os.environ, {"JARVIS_TASK_STOP_GRACE_SEC": "0.05"}):
                bot.JarvisHandler._stop_task_process(lifecycle, "stale_fence:test")
            self.assertEqual(proc.poll(), -signal.SIGKILL)
        finally:
            if proc.poll() is None:
                os.killpg(os.getpgid(proc.pid), signal.SIGKILL)
                proc.wait(timeout=5)
            if proc.stdout is not None:
                proc.stdout.close()
            if proc.stderr is not None:
                proc.stderr.close()

    def test_lease_uses_persisted_runtime_resume_and_binds_process(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._broadcast = lambda _text: None
        captured = {}

        def dispatch(*args, **kwargs):
            captured["args"] = args
            captured["kwargs"] = kwargs
            return "done"

        handler.dispatch_item = dispatch
        handler.execution_router = SimpleNamespace(task_types={"ticket"})

        class Lifecycle:
            runtime_session_id = "runtime-stable"
            resumed = True

            def bind_process(self, process):
                self.process = process

        lifecycle = Lifecycle()
        lease = {
            "task": {"taskType": "ticket", "payload": {
                "itemId": "84345050", "kind": "ticket",
                "prompt": "newer-current-task-input",
            }},
            "session": {"inputPayload": {
                "itemId": "84345050", "kind": "ticket", "prompt": "go",
                "terraform": True, "target": "group-1", "targetType": "group",
            }},
        }
        self.assertEqual(handler._execute_task_lease(lease, lifecycle), "done")
        self.assertEqual(captured["args"][1], "go")
        self.assertEqual(captured["args"][2], "runtime-stable")
        self.assertTrue(captured["args"][3])
        self.assertEqual(captured["kwargs"]["on_spawn"], lifecycle.bind_process)
        self.assertIs(captured["kwargs"]["session_controller"], lifecycle)

    def test_malformed_session_input_snapshot_never_falls_forward_to_current_task(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._broadcast = lambda _text: None
        handler.dispatch_item = mock.Mock()
        handler.execution_router = SimpleNamespace(task_types={"ticket"})
        lifecycle = SimpleNamespace(runtime_session_id="r", resumed=True,
                                    bind_process=lambda _p: None)
        lease = {
            "task": {"taskType": "probe", "payload": {
                "itemId": "probe-new", "kind": "probe", "prompt": "new-input"}},
            "session": {"inputPayload": "not-json"},
        }
        with self.assertRaisesRegex(ValueError, "payload must be JSON object"):
            handler._execute_task_lease(lease, lifecycle)
        handler.dispatch_item.assert_not_called()

    def test_null_session_input_snapshot_never_falls_forward_to_newer_task(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._broadcast = lambda _text: None
        handler.dispatch_item = mock.Mock()
        handler.execution_router = SimpleNamespace(task_types={"ticket"})
        lifecycle = SimpleNamespace(runtime_session_id="r", resumed=True,
                                    bind_process=lambda _p: None)
        lease = {
            "task": {"taskType": "probe", "payload": {
                "itemId": "probe-new", "kind": "probe", "prompt": "new-input"}},
            "session": {"inputPayload": None, "processingRevision": "old-revision"},
        }
        with self.assertRaisesRegex(ValueError, "input snapshot is null"):
            handler._execute_task_lease(lease, lifecycle)
        handler.dispatch_item.assert_not_called()

    def test_task_dispatch_never_writes_local_inflight_current_state(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.ephemeral_executor = SimpleNamespace(_closed=False)
        handler._maybe_suspend = lambda *args, **kwargs: None
        handler._completion_broadcast = lambda _item_id: "done"
        result = SimpleNamespace(text="ok", is_error=False, subtype="success")
        notices = []
        with mock.patch.object(bot, "run_claude_buffered", return_value=result), \
                mock.patch.object(bot, "_inflight_add") as add, \
                mock.patch.object(bot, "_inflight_remove") as remove:
            outcome = handler.dispatch_item(
                "task-probe", "go", "runtime-1", False, notices.append,
                "target", "group", kind="probe", session_controller=object())
        self.assertEqual(outcome, "done")
        add.assert_not_called()
        remove.assert_not_called()

    def test_ephemeral_probe_lease_is_rejected_before_execution(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._broadcast = lambda _text: None
        handler.dispatch_item = mock.Mock()
        handler.execution_router = SimpleNamespace(task_types={"ticket"})
        lifecycle = SimpleNamespace(runtime_session_id="r", resumed=False,
                                    bind_process=lambda _p: None)
        with self.assertRaisesRegex(ValueError, "not enabled"):
            handler._execute_task_lease(
                {"task": {"taskType": "probe", "payload": {
                    "itemId": "probe-1", "kind": "probe", "prompt": "go"}},
                 "session": {"inputPayload": {
                     "itemId": "probe-1", "kind": "probe", "prompt": "go"}}},
                lifecycle)
        handler.dispatch_item.assert_not_called()

    def test_task_lease_without_session_snapshot_is_rejected(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._broadcast = lambda _text: None
        handler.dispatch_item = mock.Mock()
        handler.execution_router = SimpleNamespace(task_types={"ticket"})
        lifecycle = SimpleNamespace(runtime_session_id="r", resumed=False,
                                    bind_process=lambda _p: None)
        with self.assertRaisesRegex(ValueError, "input snapshot is missing"):
            handler._execute_task_lease(
                {"task": {"taskType": "ticket", "payload": {
                    "itemId": "84345050", "kind": "ticket", "prompt": "go"}},
                 "session": {}},
                lifecycle)
        handler.dispatch_item.assert_not_called()

    def test_task_suspend_returns_central_wait_state_without_local_watcher(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.ephemeral_executor = SimpleNamespace(_closed=False)
        handler.watcher = SimpleNamespace(suspend=mock.Mock())
        handler._last_comment_id = lambda _iid: 41
        handler._workitem_line = lambda iid: "#%s" % iid
        result = SimpleNamespace(
            text='waiting [[SUSPEND:{"aone_id":"843","wait_for":"320687"}]]',
            is_error=False, subtype="success")
        with mock.patch.object(bot, "run_claude_buffered", return_value=result), \
                mock.patch.object(bot, "_inflight_add") as add:
            outcome = handler.dispatch_item(
                "task-probe", "go", "runtime-1", False, lambda _text: None,
                "target", "group", kind="probe", session_controller=object())
        self.assertEqual(outcome["status"], "suspended")
        self.assertEqual(outcome["waitType"], "AONE_REPLY")
        self.assertEqual(outcome["waitKey"], "843")
        self.assertEqual(outcome["waitCursor"], "41")
        self.assertIn("waitExpireAt", outcome)
        handler.watcher.suspend.assert_not_called()
        add.assert_not_called()


class TaskAoneAssociationTest(unittest.TestCase):
    def test_association_survives_every_trigger_source_transition(self):
        cases = (
            ("ticket", "AONE", "SCAN"),
            ("pr_comment_reply", "GITHUB", "PR_COMMENT"),
            ("pr_ci_fix", "GITHUB", "PR_CI_FAILED"),
            ("wake", "AONE", "WAKE"),
            ("persona", "AONE", "PERSONA"),
            ("revisit", "AONE", "REVISIT"),
        )
        envelopes = [
            bot._task_envelope(
                item_id="84345050",
                project="2100304",
                task_type=task_type,
                source_type=source_type,
                source_ref={"sequence": index},
                desired_revision="sequence:%d" % index,
                trigger=trigger,
                prompt="step %d" % index,
            )
            for index, (task_type, source_type, trigger) in enumerate(cases)
        ]
        self.assertEqual(
            {envelope.task_key for envelope in envelopes},
            {"aone:2100304:84345050"})
        self.assertEqual(
            [envelope.aone_id for envelope in envelopes],
            ["84345050"] * len(cases))
        self.assertTrue(all(
            envelope.to_dict().get("aoneId") == "84345050"
            for envelope in envelopes))


class WakeRoutingTest(unittest.TestCase):
    def _handler(self, accepted=True):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.no_dingtalk = True
        handler._broadcast = lambda _text: None
        handler._workitem_line = lambda _iid: "#843"
        handler._workitem_project = lambda _iid: "2100304"
        handler._quick_card = mock.Mock()
        handler.ephemeral_executor = SimpleNamespace(
            submit=lambda *args, **kwargs: (accepted, "dispatched" if accepted else "full"))
        captured = {}

        class Router:
            def enqueue(self, envelope, local_submit=None):
                del local_submit
                captured["envelope"] = envelope
                return EnqueueResult(
                    accepted, "persisted" if accepted else "rejected")

        handler.execution_router = Router()
        return handler, captured

    def test_wake_uses_unified_project_key_and_returns_accepted(self):
        handler, captured = self._handler(True)
        task = {"target": "g", "target_type": "group", "session_id": "s"}
        accepted = handler._wake("843", task, [
            {"id": 8, "creator": "a", "content": "x"},
            {"id": 9, "creator": "b", "content": "y"},
        ])
        self.assertTrue(accepted)
        envelope = captured["envelope"]
        self.assertEqual(envelope.task_key, "aone:2100304:843")
        self.assertEqual(envelope.aone_id, "843")
        self.assertEqual(envelope.desired_revision, "comment:9")
        handler._quick_card.assert_called_once()

    def test_rejected_wake_is_not_announced_as_started(self):
        handler, _captured = self._handler(False)
        task = {"target": "g", "target_type": "group", "session_id": "s",
                "project": "2100304"}
        self.assertFalse(handler._wake(
            "843", task, [{"id": 9, "creator": "a", "content": "x"}]))
        handler._quick_card.assert_not_called()

    def test_wait_record_is_removed_only_after_accepted_wake(self):
        task = {"last_comment_id": 1, "last_poll": 0,
                "suspended_at": time.time() - 1000,
                "target": "g", "target_type": "group"}
        handler = SimpleNamespace(_wake=mock.Mock(return_value=False))
        watcher = bot.WaitWatcher.__new__(bot.WaitWatcher)
        watcher.handler = handler
        watcher.suspended = {"843": task}
        watcher._lock = threading.Lock()
        watcher._fetch_comments = lambda _iid: [
            {"id": 2, "creator": "human", "content": "go"}]
        watcher._remove_persisted = mock.Mock()
        watcher._tick()
        self.assertIn("843", watcher.suspended)
        watcher._remove_persisted.assert_not_called()

        task["last_poll"] = 0
        handler._wake.return_value = True
        watcher._tick()
        self.assertNotIn("843", watcher.suspended)
        watcher._remove_persisted.assert_called_once_with("843")


if __name__ == "__main__":
    unittest.main()
