#!/usr/bin/env python3
"""Hermetic integration tests for the bridge/control-plane seam."""

import json
import os
import signal
import subprocess
import sys
import tempfile
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
        self.assertTrue(callable(
            _FakePersistenceExecutor.instances[-1].kwargs["progress"]))
        for obsolete in ("task_router", "local_worker", "dispatch_pool"):
            self.assertFalse(hasattr(handler, obsolete))

    def test_exposes_assistant_text_and_tool_name_without_tool_payloads(self):
        with tempfile.TemporaryDirectory() as directory:
            transcript = Path(directory) / "session.jsonl"
            transcript.write_text("\n".join([
                json.dumps({
                    "message": {"role": "assistant", "content": [
                        {"type": "text", "text": "正在检查控制面状态"},
                        {"type": "tool_use", "name": "exec_command",
                         "input": {"token": "must-not-leak"}},
                    ]}
                }, ensure_ascii=False),
                json.dumps({
                    "message": {"role": "user", "content": [
                        {"type": "tool_result", "content": "secret-result"}
                    ]}
                }, ensure_ascii=False),
            ]), encoding="utf-8")
            with mock.patch.object(bot, "_session_file", return_value=transcript):
                excerpt = bot._session_progress_excerpt("runtime-1")

        self.assertIn("正在检查控制面状态", excerpt)
        self.assertIn("[执行工具 exec_command]", excerpt)
        self.assertNotIn("must-not-leak", excerpt)
        self.assertNotIn("secret-result", excerpt)

    def test_exposes_latest_subagent_progress_in_timestamp_order(self):
        with tempfile.TemporaryDirectory() as directory:
            transcript = Path(directory) / "session.jsonl"
            transcript.write_text(json.dumps({
                "timestamp": "2026-07-20T09:13:13Z",
                "message": {"role": "assistant", "content": [
                    {"type": "text", "text": "开始 terraform-rd"},
                ]},
            }, ensure_ascii=False), encoding="utf-8")
            subagents = transcript.with_suffix("") / "subagents"
            subagents.mkdir(parents=True)
            (subagents / "agent-rd.jsonl").write_text("\n".join([
                json.dumps({
                    "timestamp": "2026-07-20T09:14:00Z",
                    "message": {"role": "assistant", "content": [
                        {"type": "text", "text": "正在检查评审意见"},
                    ]},
                }, ensure_ascii=False),
                json.dumps({
                    "timestamp": "2026-07-20T09:15:00Z",
                    "message": {"role": "assistant", "content": [
                        {"type": "tool_use", "name": "Bash",
                         "input": {"command": "secret command"}},
                    ]},
                }, ensure_ascii=False),
            ]), encoding="utf-8")
            with mock.patch.object(bot, "_session_file", return_value=transcript):
                excerpt = bot._session_progress_excerpt("runtime-1")

        self.assertLess(excerpt.index("开始 terraform-rd"),
                        excerpt.index("正在检查评审意见"))
        self.assertIn("[执行工具 Bash]", excerpt)
        self.assertNotIn("secret command", excerpt)

    def test_redacts_sensitive_values_repeated_by_assistant(self):
        with tempfile.TemporaryDirectory() as directory:
            transcript = Path(directory) / "session.jsonl"
            transcript.write_text(json.dumps({
                "timestamp": "2026-07-20T09:14:00Z",
                "message": {"role": "assistant", "content": [
                    {"type": "text",
                     "text": "access_key_secret=super-secret-value"},
                ]},
            }), encoding="utf-8")
            with mock.patch.object(bot, "_session_file", return_value=transcript):
                excerpt = bot._session_progress_excerpt("runtime-1")

        self.assertIn("access_key_secret=[REDACTED]", excerpt)
        self.assertNotIn("super-secret-value", excerpt)

    def test_worker_starts_before_every_sensor(self):
        calls = []
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.persistence_executor = _Starter("worker", calls)
        for name in ("scanner", "daily", "aone_reply_scheduler", "prwatch"):
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


class AoneReplySchedulerTest(unittest.TestCase):
    def _handler(self, pages, comments, wake_result=True):
        client = mock.Mock()
        client.list_pending_aone_reply_waits.side_effect = pages
        handler = SimpleNamespace(
            task_client=client,
            _wake=mock.Mock(return_value=wake_result),
        )
        return handler

    def _sensor(self, handler, comments):
        sensor = bot.AoneReplyScheduler(handler)
        sensor._fetch_comments = mock.Mock(return_value=comments)
        return sensor

    @staticmethod
    def _wait(session_id=10, cursor="40"):
        return {
            "task": {
                "taskKey": "aone:2100304:84345050",
                "aoneId": "84345050",
                "sourceRef": {"projectId": "2100304", "title": "Managed title"},
                "payload": {"project": "2100304"},
            },
            "session": {
                "id": session_id,
                "runtimeSessionId": "runtime-1",
                "waitType": "AONE_REPLY",
                "waitKey": "84345050",
                "waitCursor": cursor,
                "inputPayload": {
                    "project": "2100304", "terraform": True,
                    "target": "staff-1", "targetType": "user",
                },
            },
        }

    def test_rebuilds_wait_from_control_plane_and_wakes_only_new_external_comments(self):
        page = {"items": [self._wait()], "nextAfterSessionId": 10, "hasMore": False}
        comments = [
            {"id": "39", "creator": "human", "content": "old"},
            {"id": "41", "creator": "open-jarvis", "content": "self"},
            {"id": "43", "creator": "human", "content": "reply"},
            {"id": "42", "creator": "human-2", "content": "earlier reply"},
        ]
        handler = self._handler([page], comments)
        sensor = self._sensor(handler, comments)

        sensor._tick()

        handler._wake.assert_called_once()
        aone_id, context, observed = handler._wake.call_args.args
        self.assertEqual(aone_id, "84345050")
        self.assertEqual(context["session_id"], "runtime-1")
        self.assertTrue(context["terraform"])
        self.assertEqual(context["title"], "Managed title")
        self.assertEqual([int(c["id"]) for c in observed], [42, 43])

    def test_list_failure_keeps_local_throttle_and_retries_without_wake(self):
        handler = self._handler([RuntimeError("503")], [])
        sensor = self._sensor(handler, [])
        sensor._poll_state["10"] = {"first_seen": 1, "last_poll": 2}

        sensor._tick()

        self.assertIn("10", sensor._poll_state)
        handler._wake.assert_not_called()

    def test_failed_wake_keeps_wait_discovery_retryable(self):
        page = {"items": [self._wait()], "nextAfterSessionId": 10, "hasMore": False}
        comments = [{"id": 41, "creator": "human", "content": "reply"}]
        handler = self._handler([page], comments, wake_result=False)
        sensor = self._sensor(handler, comments)

        sensor._tick()

        self.assertIn("10", sensor._poll_state)
        handler._wake.assert_called_once()

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
        handler._routine_notice = mock.Mock()
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
        with mock.patch.object(bot, "_terraform_rd_ready", return_value=True):
            self.assertEqual(handler._execute_task_lease(lease, lifecycle), "done")
        self.assertEqual(captured["args"][1], "go")
        self.assertEqual(captured["args"][2], "runtime-stable")
        self.assertTrue(captured["args"][3])
        # ticket is a TASK_BOOKEND_KIND: the executor owns the Aone bookend, so on_spawn
        # is the bookend's bind_process (which still binds the controller's process).
        bookend = captured["kwargs"]["task_bookend"]
        self.assertIsInstance(bookend, bot._TaskAoneBookend)
        self.assertIs(bookend.controller, lifecycle)
        self.assertEqual(captured["kwargs"]["on_spawn"], bookend.bind_process)
        self.assertIs(captured["kwargs"]["session_controller"], lifecycle)
        self.assertIs(captured["args"][4], handler._routine_notice)

    def test_adhoc_task_replies_to_original_conversation(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._quick_card = mock.Mock()
        captured = {}

        def dispatch(*args, **kwargs):
            captured["notify"] = args[4]
            return "done"

        handler.dispatch_item = dispatch
        handler.execution_router = SimpleNamespace(task_types={"adhoc"})
        lifecycle = SimpleNamespace(runtime_session_id="r", resumed=False,
                                    bind_process=lambda _p: None)
        lease = {
            "task": {"taskType": "adhoc"},
            "session": {"inputPayload": {
                "itemId": "handoff-1", "kind": "adhoc", "prompt": "go",
                "target": "conversation-1", "targetType": "group",
            }},
        }

        self.assertEqual(handler._execute_task_lease(lease, lifecycle), "done")
        captured["notify"]("finished")
        handler._quick_card.assert_called_once_with(
            "conversation-1", "finished", "group")

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
        # The on-disk inflight registry is gone entirely: live-process tracking is owned
        # by EphemeralExecutor._active (in-memory). This is now a structural guarantee.
        self.assertFalse(hasattr(bot, "_inflight_add"))
        self.assertFalse(hasattr(bot, "INFLIGHT_PATH"))
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.ephemeral_executor = SimpleNamespace(_closed=False)
        handler._maybe_suspend = lambda *args, **kwargs: None
        handler._completion_broadcast = lambda _item_id: "done"
        result = SimpleNamespace(text="ok", is_error=False, subtype="success")
        notices = []
        with mock.patch.object(bot, "run_claude_buffered", return_value=result):
            outcome = handler.dispatch_item(
                "task-probe", "go", "runtime-1", False, notices.append,
                "target", "group", kind="probe", session_controller=object())
        self.assertEqual(outcome, "done")

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
        with mock.patch.object(bot, "run_claude_buffered", return_value=result):
            outcome = handler.dispatch_item(
                "task-probe", "go", "runtime-1", False, lambda _text: None,
                "target", "group", kind="probe", session_controller=object())
        self.assertEqual(outcome["status"], "suspended")
        self.assertEqual(outcome["waitType"], "AONE_REPLY")
        self.assertEqual(outcome["waitKey"], "843")
        self.assertEqual(outcome["waitCursor"], "41")
        self.assertIn("waitExpireAt", outcome)
        handler.watcher.suspend.assert_not_called()


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

    def test_card_submit_reports_aone_title_only_in_source_ref(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        captured = {}

        class Router:
            def enqueue(self, envelope, local_submit=None):
                del local_submit
                captured["envelope"] = envelope
                return EnqueueResult(True, "task_persisted")

        handler.execution_router = Router()
        handler.ephemeral_executor = object()
        handler._quick_card = mock.Mock()
        handler._submit_card(
            "84345050", "group", "group", "go", "runtime", False,
            project="2100304", title="Aone card title")
        envelope = captured["envelope"]
        self.assertEqual(envelope.source_ref, {
            "aoneId": "84345050", "projectId": "2100304",
            "title": "Aone card title",
        })
        self.assertNotIn("title", envelope.payload)
        self.assertNotIn("Aone card title", envelope.desired_revision)


class WakeRoutingTest(unittest.TestCase):
    def _handler(self, accepted=True):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.no_dingtalk = True
        handler._broadcast = lambda _text: None
        handler._workitem_line = lambda _iid: "#843"
        handler._workitem_project = lambda _iid: "2100304"
        handler._workitem_title = lambda _iid: "Point-read title"
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
        self.assertEqual(envelope.source_ref["title"], "Point-read title")
        handler._quick_card.assert_not_called()

    def test_rejected_wake_is_not_announced_as_started(self):
        handler, _captured = self._handler(False)
        task = {"target": "g", "target_type": "group", "session_id": "s",
                "project": "2100304"}
        self.assertFalse(handler._wake(
            "843", task, [{"id": 9, "creator": "a", "content": "x"}]))
        handler._quick_card.assert_not_called()


class AoneSchedulerUnionTest(unittest.TestCase):
    """AoneScheduler 统一探测：每池 assignee∪tracker∪idle 并集去重，取代 scan.sh 单一 assignee
    出数据，消除「指派给人 / 抄送数字人」盲区（黑洞成因）。"""

    def _scanner(self):
        s = bot.AoneScheduler.__new__(bot.AoneScheduler)
        return s

    def test_query_pool_union_merges_three_sources_and_dedups(self):
        s = self._scanner()
        seen_filters = []

        def fake_a1_list(project, flt):
            seen_filters.append(flt)
            if flt.startswith("assignedTo="):
                return [{"id": "1", "title": "指派单"}]
            if flt.startswith("workitem.tracker="):
                return [{"id": "1", "title": "重复(抄送同一单)"},
                        {"id": "2", "title": "抄送数字人单"}]
            if flt.startswith("tag=jarvis-idle"):
                return [{"id": "3", "title": "idle 单"}]
            return []

        s._a1_list = fake_a1_list
        rows = s._query_pool_union("tf_customer", "1086837", ["Closed", "已发布"])
        ids = sorted(r["id"] for r in rows)
        self.assertEqual(ids, ["1", "2", "3"], "三源并集按 id 去重（#1 只保留一次）")
        # 三源查询并行发出 → seen_filters 顺序不定，按集合断言（顺序无关）
        worker_csv = ",".join(sorted(bot.DIGITAL_WORKER_IDS))
        self.assertEqual(len(seen_filters), 3, "assignee/tracker/idle 三源各查一次")
        # 每源都叠加 pools.json 状态排除
        self.assertTrue(all("NOT status=Closed" in f and "NOT status=已发布" in f
                            for f in seen_filters),
                        "三源过滤须叠加 exclude_status")
        # 数字人 id 单一真源
        self.assertTrue(any("assignedTo=%s" % worker_csv in f for f in seen_filters))
        self.assertTrue(any("workitem.tracker=%s" % worker_csv in f for f in seen_filters))
        self.assertTrue(any("tag=jarvis-idle" in f for f in seen_filters))
        # 每源都带 pool/pool_project 戳
        self.assertTrue(all(r.get("pool") == "tf_customer"
                            and r.get("pool_project") == "1086837" for r in rows))

    def test_scan_union_iterates_pools(self):
        s = self._scanner()
        with mock.patch.object(bot.AoneScheduler, "_read_pools",
                               return_value=[("tf_customer", "1086837", []),
                                             ("tf_provider", "528766", [])]):
            calls = []

            def fake_union(key, project, excl):
                calls.append(key)
                return [{"id": "%s-x" % key, "pool": key, "pool_project": project}]
            s._query_pool_union = fake_union
            items = s._scan_union()
        self.assertEqual(calls, ["tf_customer", "tf_provider"])
        self.assertEqual({it["id"] for it in items},
                         {"tf_customer-x", "tf_provider-x"})

    def test_scan_union_no_pools_returns_none(self):
        s = self._scanner()
        with mock.patch.object(bot.AoneScheduler, "_read_pools", return_value=[]):
            self.assertIsNone(s._scan_union())

    def test_digital_worker_ids_single_source(self):
        # 单一真源含编排层 + 公开 RD + 旧 PD/QA 兼容 worker
        self.assertIn(bot.JARVIS_ORCH_WORKER, bot.DIGITAL_WORKER_IDS)
        self.assertIn(bot.PERSONA_PUBLIC_WORKER, bot.DIGITAL_WORKER_IDS)
        self.assertTrue(bot.PERSONA_LEGACY_WORKER_IDS <= bot.DIGITAL_WORKER_IDS)


class TaskBookendDispatchTest(unittest.TestCase):
    """B-proper: the control-plane Task run authors a structured result; the executor
    (via _TaskAoneBookend) owns the single Aone write. The run never self-claims, so the
    self-lease-conflict is gone. Verifies commit routing (done/idle/suspend) and the
    fail-closed path when the run exits clean without a valid [[AONE_RESULT]]."""

    def _handler(self):
        h = bot.JarvisHandler.__new__(bot.JarvisHandler)
        h._broadcast = lambda _t: None
        h._completion_broadcast = lambda _iid: "done-broadcast"
        h._maybe_suspend = lambda *a, **k: None
        h._last_comment_id = lambda _iid: 12345
        return h

    def _controller(self):
        return SimpleNamespace(
            task={"id": 603, "generation": 1},
            session={"generation": 1},
            runtime_session_id="rt-1", resumed=False)

    def _run(self, final, outcome_terraform=True):
        h = self._handler()
        ctrl = self._controller()
        bookend = bot._TaskAoneBookend(ctrl, "84407231", "1086837",
                                       outcome_terraform, "persona")
        calls = {}
        h._dispatch_failed = lambda *a, **k: calls.setdefault("failed", a)
        with mock.patch.object(bot, "run_claude_buffered",
                               return_value=bot.ClaudeResult(final, False, "success")), \
             mock.patch.object(bot, "_claim_workitem",
                               side_effect=lambda *a, **k: calls.setdefault("claim", a)), \
             mock.patch.object(bot, "_aone_event_enqueue",
                               side_effect=lambda *a, **k: calls.setdefault("reply", (a, k)) or True), \
             mock.patch.object(bot, "_finish_workitem",
                               side_effect=lambda *a, **k: calls.setdefault("finish", a)), \
             mock.patch.object(bot, "_release_post_pr_claim",
                               side_effect=lambda *a, **k: calls.setdefault("release", a)):
            out = h.dispatch_item(
                "84407231", "prompt", "sid", False, lambda _t: None,
                "tgt", "group", project="1086837", kind="persona", terraform=True,
                session_controller=ctrl, task_bookend=bookend)
        return out, calls

    def test_done_writes_reply_then_finishes(self):
        out, calls = self._run(
            '[[AONE_RESULT:{"outcome":"done","reply_body":"结论 done"}]]')
        self.assertEqual(out, "done")
        self.assertIn("reply", calls)
        self.assertIn("finish", calls)
        self.assertNotIn("release", calls)
        # terraform reply → terraform-rd identity
        self.assertEqual(calls["reply"][1].get("identity"), bot.PERSONA_PUBLIC_IDENTITY)

    def test_idle_writes_reply_then_releases(self):
        out, calls = self._run(
            '[[AONE_RESULT:{"outcome":"idle","reply_body":"阶段完成"}]]')
        self.assertEqual(out, "done")
        self.assertIn("reply", calls)
        self.assertIn("release", calls)
        self.assertNotIn("finish", calls)

    def test_suspend_writes_reply_and_returns_wait_state(self):
        out, calls = self._run(
            '[[AONE_RESULT:{"outcome":"suspend","reply_body":"@新山 请确认",'
            '"suspend_wait_for":"521957"}]]')
        self.assertIsInstance(out, dict)
        self.assertEqual(out["status"], "suspended")
        self.assertEqual(out["waitKey"], "84407231")
        self.assertEqual(out["waitCursor"], "12345")
        self.assertIn("reply", calls)
        self.assertNotIn("finish", calls)
        self.assertNotIn("release", calls)

    def test_missing_result_fails_closed_without_commit(self):
        out, calls = self._run("干完了但忘了输出结构化结果")
        self.assertEqual(out["status"], "error")
        self.assertEqual(out["error"]["subtype"], "missing_task_result")
        self.assertIn("干完了但忘了输出结构化结果", out["error"]["message"])
        self.assertIn("failed", calls)
        self.assertNotIn("reply", calls)
        self.assertNotIn("finish", calls)
        self.assertNotIn("release", calls)

    def test_non_terraform_reply_uses_jarvis_identity(self):
        h = self._handler()
        ctrl = self._controller()
        bookend = bot._TaskAoneBookend(ctrl, "999", "2100304", False, "ticket")
        calls = {}
        with mock.patch.object(bot, "run_claude_buffered",
                               return_value=bot.ClaudeResult(
                                   '[[AONE_RESULT:{"outcome":"idle","reply_body":"x"}]]',
                                   False, "success")), \
             mock.patch.object(bot, "_aone_event_enqueue",
                               side_effect=lambda *a, **k: calls.setdefault("reply", k) or True), \
             mock.patch.object(bot, "_release_post_pr_claim", lambda *a, **k: None):
            out = h.dispatch_item(
                "999", "p", "sid", False, lambda _t: None, "tgt", "group",
                project="2100304", kind="ticket", terraform=False,
                session_controller=ctrl, task_bookend=bookend)
        self.assertEqual(out, "done")
        self.assertEqual(calls["reply"].get("identity"), "jarvis")
        self.assertTrue(calls["reply"].get("allow_non_tf"))


class PostPrRerouteDispatchTest(unittest.TestCase):
    """Tier-0 acceptance for the post-PR replay-safe reroute: pr_ci_fix / pr_comment_reply
    now run through _TaskAoneBookend(writes_reply=False) — idempotent claim on bind,
    release-to-idle on clean completion, NO Aone reply, NO [[AONE_RESULT]] required, and
    NO fenced-operation calls (the deleted machinery). REPLAY_SAFE means a re-run's claim/
    release are idempotent no-ops after the first."""

    def _handler(self):
        h = bot.JarvisHandler.__new__(bot.JarvisHandler)
        h._broadcast = lambda _t: None
        h._completion_broadcast = lambda _iid: "done-broadcast"
        h._maybe_suspend = lambda *a, **k: None
        return h

    def _controller(self):
        return SimpleNamespace(
            task={"id": 700, "generation": 1}, session={"generation": 1},
            runtime_session_id="rt-1", resumed=False,
            bind_process=lambda _p: None)

    def _run(self, kind, final, is_error=False):
        h = self._handler()
        ctrl = self._controller()
        bookend = bot._TaskAoneBookend(ctrl, "84448059", "528766", True, kind,
                                       writes_reply=False)
        calls = {"claim": 0, "release": 0, "reply": 0, "finish": 0}
        h._dispatch_failed = lambda *a, **k: calls.__setitem__("failed", True)
        with mock.patch.dict(os.environ, {"JARVIS_DISPATCH_RETRY_MAX": "0",
                                          "JARVIS_DISPATCH_RETRY_BACKOFF": "0"}), \
             mock.patch.object(bot, "run_claude_buffered",
                               return_value=bot.ClaudeResult(final, is_error, "success")), \
             mock.patch.object(bot, "_claim_workitem",
                               side_effect=lambda *a, **k: calls.__setitem__("claim", calls["claim"] + 1)), \
             mock.patch.object(bot, "_release_post_pr_claim",
                               side_effect=lambda *a, **k: calls.__setitem__("release", calls["release"] + 1)), \
             mock.patch.object(bot, "_finish_workitem",
                               side_effect=lambda *a, **k: calls.__setitem__("finish", calls["finish"] + 1)), \
             mock.patch.object(bot, "_aone_event_enqueue",
                               side_effect=lambda *a, **k: calls.__setitem__("reply", calls["reply"] + 1) or True):
            # bind (claim) then dispatch, mirroring _execute_task_lease on_spawn=bind_process
            bookend.bind_process(SimpleNamespace(pid=1))
            out = h.dispatch_item(
                "84448059", "prompt", "sid", False, lambda _t: None, "tgt", "group",
                project="528766", kind=kind, terraform=True,
                session_controller=ctrl, task_bookend=bookend)
        return out, calls, bookend

    def test_pr_ci_fix_clean_completion_releases_idle_no_reply(self):
        # No [[AONE_RESULT]] in the run text — post-PR runs never emit one.
        out, calls, _bk = self._run("pr_ci_fix", "CI 已修复并 force-push")
        self.assertEqual(out, "done")
        self.assertEqual(calls["claim"], 1)      # claimed once on bind
        self.assertEqual(calls["release"], 1)    # released-to-idle once on completion
        self.assertEqual(calls["reply"], 0)      # NO Aone reply comment
        self.assertEqual(calls["finish"], 0)     # never finishes (PR still open)
        self.assertNotIn("failed", calls)

    def test_pr_comment_reply_same_reroute(self):
        out, calls, _bk = self._run("pr_comment_reply", "已回复评审意见")
        self.assertEqual(out, "done")
        self.assertEqual((calls["claim"], calls["release"], calls["reply"]), (1, 1, 0))

    def test_error_does_not_release_so_replay_reclaims(self):
        # A failed run stays claimed (no release); REPLAY_SAFE re-lease re-claims.
        out, calls, _bk = self._run("pr_ci_fix", "boom", is_error=True)
        self.assertEqual(out["status"], "error")
        self.assertEqual(out["error"]["message"], "boom")
        self.assertIn("failed", calls)
        self.assertEqual(calls["release"], 0)

    def test_claim_failure_detail_is_returned_to_control_plane(self):
        h = self._handler()
        ctrl = self._controller()
        bookend = bot._TaskAoneBookend(
            ctrl, "9001", "528766", True, "pr_ci_fix", writes_reply=False)
        h._dispatch_failed = lambda *a, **k: None
        detail = ("bridge claim failed for #9001 (rc=3): "
                  "【涉及云产品】不能为空,【Terraform需求类型】不能为空")
        with mock.patch.object(bot, "run_claude_buffered", side_effect=RuntimeError(detail)):
            out = h.dispatch_item(
                "9001", "prompt", "sid", False, lambda _t: None, "tgt", "group",
                project="528766", kind="pr_ci_fix", terraform=True,
                session_controller=ctrl, task_bookend=bookend)
        self.assertEqual(out["status"], "error")
        self.assertEqual(out["error"]["errorType"], "AoneClaimFailed")
        self.assertIn("涉及云产品", out["error"]["message"])
        self.assertIn("Terraform需求类型", out["error"]["message"])

    def test_release_idle_is_idempotent_under_replay(self):
        # Simulate the completion twice (a re-lease re-runs release_idle): still one write.
        _out, _calls, bk = self._run("pr_ci_fix", "ok")
        released_before = None
        with mock.patch.object(bot, "_release_post_pr_claim") as rel:
            bk.release_idle()        # second call after the dispatch already released
            bk.release_idle()
            released_before = rel.call_count
        self.assertEqual(released_before, 0)  # already released → no further writes

    def test_writes_reply_false_bookend_has_no_fenced_operation_surface(self):
        # The reroute must not resurrect the deleted fenced-operation API.
        bk = bot._TaskAoneBookend(self._controller(), "1", "528766", True,
                                  "pr_ci_fix", writes_reply=False)
        self.assertFalse(hasattr(bk, "lineage_policy"))
        self.assertFalse(hasattr(bk, "_begin"))
        self.assertFalse(hasattr(bot, "_PostPrTaskBookend"))


class CompletionBroadcastTest(unittest.TestCase):
    """_completion_broadcast must reflect what the run left on the ticket, not
    merely that the process exited cleanly (see self-lease-conflict fix)."""

    def _broadcast(self, line, tag):
        with mock.patch.object(bot.JarvisHandler, "_workitem_line",
                               return_value=(line, tag)):
            return bot.JarvisHandler._completion_broadcast(None, "84407231")

    def test_jarvis_done_reports_completion(self):
        self.assertTrue(self._broadcast("- [#84407231](url) t", "jarvis-done")
                        .startswith("✅ 工单处理完成"))

    def test_jarvis_idle_reports_staged(self):
        self.assertTrue(self._broadcast("- [#84407231](url) t", "jarvis-idle")
                        .startswith("⏸️ 工单阶段完成·待人工接手"))

    def test_jarvis_claimed_reports_unwrapped(self):
        self.assertIn("未收尾", self._broadcast("- [#84407231](url) t", "jarvis-claimed"))

    def test_empty_tag_is_not_reported_as_completion(self):
        out = self._broadcast("- [#84407231](url) t", "")
        self.assertFalse(out.startswith("✅"))
        self.assertIn("未获认领", out)

    def test_numeric_query_failure_falls_back_to_headless(self):
        with mock.patch.object(bot.JarvisHandler, "_workitem_line",
                               return_value="#84407231"):
            out = bot.JarvisHandler._completion_broadcast(None, "84407231")
        self.assertIn("处理完成（headless）", out)

    def test_non_numeric_pseudo_id_has_no_workitem(self):
        with mock.patch.object(bot.JarvisHandler, "_workitem_line",
                               return_value="#probe-2026-07-17"):
            out = bot.JarvisHandler._completion_broadcast(None, "probe-2026-07-17")
        self.assertIn("任务 #probe-2026-07-17 处理完成", out)


class ExtractTaskResultTest(unittest.TestCase):
    """The control-plane Task-path run hands the executor a structured outcome
    instead of writing Aone itself. A missing/garbage sentinel must yield None so
    the caller fails closed rather than silently succeeding."""

    def test_valid_done_result(self):
        text = ('前言\n[[AONE_RESULT:{"outcome":"done","reply_body":"结论 X",'
                '"target_status":"已发布","mr_cr_links":["http://mr/1"]}]]\n尾')
        clean, res = bot.extract_task_result(text)
        self.assertNotIn("AONE_RESULT", clean)
        self.assertEqual(res["outcome"], "done")
        self.assertEqual(res["reply_body"], "结论 X")
        self.assertEqual(res["target_status"], "已发布")
        self.assertEqual(res["mr_cr_links"], ["http://mr/1"])

    def test_no_sentinel_returns_none(self):
        clean, res = bot.extract_task_result("just some run output, no sentinel")
        self.assertIsNone(res)
        self.assertEqual(clean, "just some run output, no sentinel")

    def test_invalid_outcome_rejected(self):
        _clean, res = bot.extract_task_result(
            '[[AONE_RESULT:{"outcome":"maybe","reply_body":"x"}]]')
        self.assertIsNone(res)

    def test_empty_reply_body_rejected(self):
        _clean, res = bot.extract_task_result(
            '[[AONE_RESULT:{"outcome":"idle","reply_body":"  "}]]')
        self.assertIsNone(res)

    def test_suspend_requires_wait_for(self):
        _clean, res = bot.extract_task_result(
            '[[AONE_RESULT:{"outcome":"suspend","reply_body":"等确认"}]]')
        self.assertIsNone(res)
        _clean, ok = bot.extract_task_result(
            '[[AONE_RESULT:{"outcome":"suspend","reply_body":"等确认",'
            '"suspend_wait_for":"320687"}]]')
        self.assertEqual(ok["outcome"], "suspend")
        self.assertEqual(ok["suspend_wait_for"], "320687")

    def test_last_span_wins_and_all_stripped(self):
        text = ('[[AONE_RESULT:{"outcome":"idle","reply_body":"first"}]]'
                'mid'
                '[[AONE_RESULT:{"outcome":"done","reply_body":"second"}]]')
        clean, res = bot.extract_task_result(text)
        self.assertEqual(res["reply_body"], "second")
        self.assertNotIn("AONE_RESULT", clean)


if __name__ == "__main__":
    unittest.main()
