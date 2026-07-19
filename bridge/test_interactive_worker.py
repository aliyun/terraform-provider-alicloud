#!/usr/bin/env python3
"""Hermetic lifecycle tests for the interactive Claude/Codex worker."""

import contextlib
import importlib.util
import io
import os
import stat
import tempfile
import threading
import unittest
from pathlib import Path
from unittest import mock

HERE = Path(__file__).resolve().parent
MODULE_PATH = HERE.parent / "bootstrap" / "jarvis-interactive-worker.py"
SPEC = importlib.util.spec_from_file_location("jarvis_interactive_worker", MODULE_PATH)
assert SPEC and SPEC.loader
worker = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(worker)


class FakeClient:
    def __init__(self):
        self.calls = []
        self.claim_error = None
        self.begin_error = None
        self.worker_heartbeat_error = None
        self.heartbeat_error = None
        self.suspend_error = None
        self.claim_results = []
        self.begin_results = [{
            "proceed": True,
            "operation": {"id": "op-1", "status": "SENDING"},
        }]

    def _record(self, name, *args, **kwargs):
        self.calls.append((name, args, kwargs))
        return {"ok": True}

    def register_worker(self, *args, **kwargs):
        return self._record("register_worker", *args, **kwargs)

    def heartbeat_worker(self, *args, **kwargs):
        self._record("heartbeat_worker", *args, **kwargs)
        if self.worker_heartbeat_error:
            raise self.worker_heartbeat_error
        return {"ok": True}

    def heartbeat_session(self, *args, **kwargs):
        self._record("heartbeat_session", *args, **kwargs)
        if self.heartbeat_error:
            raise self.heartbeat_error
        detail = args[3] if len(args) > 3 and isinstance(args[3], dict) else {}
        return {
            "session": {
                "id": args[0],
                "fenceToken": args[2],
                "status": "RUNNING",
                "leaseExpireAt": worker.time.time() + int(
                    detail.get("leaseSeconds") or 300),
            },
        }

    def claim_task(self, *args, **kwargs):
        self._record("claim_task", *args, **kwargs)
        if self.claim_error:
            raise self.claim_error
        if self.claim_results:
            return self.claim_results.pop(0)
        return {
            "task": {"id": "task-1", "generation": 4},
            "session": {"id": "session-1", "generation": 4,
                        "fenceToken": 9, "attemptNo": 2},
        }

    def start_session(self, *args, **kwargs):
        self._record("start_session", *args, **kwargs)
        detail = args[3] if len(args) > 3 and isinstance(args[3], dict) else {}
        return {
            "session": {
                "id": args[0],
                "fenceToken": args[2],
                "status": "RUNNING",
                "leaseExpireAt": worker.time.time() + int(
                    detail.get("leaseSeconds") or 300),
            },
        }

    def begin_operation(self, *args, **kwargs):
        self._record("begin_operation", *args, **kwargs)
        if self.begin_error:
            raise self.begin_error
        return self.begin_results.pop(0)

    def ack_operation(self, *args, **kwargs):
        return self._record("ack_operation", *args, **kwargs)

    def fail_operation(self, *args, **kwargs):
        return self._record("fail_operation", *args, **kwargs)

    def reconcile_operation(self, *args, **kwargs):
        return self._record("reconcile_operation", *args, **kwargs)

    def fail_session(self, *args, **kwargs):
        return self._record("fail_session", *args, **kwargs)

    def suspend_session(self, *args, **kwargs):
        self._record("suspend_session", *args, **kwargs)
        if self.suspend_error:
            raise self.suspend_error
        return {"ok": True}

    def complete_session(self, *args, **kwargs):
        return self._record("complete_session", *args, **kwargs)

    def get_worker_state(self, *args, **kwargs):
        self._record("get_worker_state", *args, **kwargs)
        return {"worker": {"workerKey": args[0]}}


class InteractiveWorkerTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.env = mock.patch.dict(os.environ, {
            "HOME": self.temp.name,
            "JARVIS_INTERACTIVE_STATE_DIR": self.temp.name,
            "JARVIS_INTERACTIVE_RETRY_DELAY": "0",
            "CODEX_THREAD_ID": "native-thread-1",
        }, clear=False)
        self.env.start()
        os.environ.pop("CLAUDE_CODE_SESSION_ID", None)

    def tearDown(self):
        self.env.stop()
        self.temp.cleanup()

    def _store(self):
        return worker.StateStore(worker._state_path("codex", "native-thread-1"))

    def _claude_store(self, session_id="headless-session-1"):
        return worker.StateStore(worker._state_path("claude", session_id))

    def _child_store(self, session_id="agent-child-1"):
        return worker.StateStore(worker._state_path("codex", session_id))

    def _seed(self):
        state = {
            "schemaVersion": 1,
            "client": "codex",
            "clientSessionId": "native-thread-1",
            "workerKey": "host:boot:process",
            "host": "host",
            "bootId": "boot",
            "processUuid": "process",
            "hostPid": os.getpid(),
            "hostProcessStartedAt": worker._process_start_identity(os.getpid()),
            "verifyHostCommand": True,
            "cwd": "/workspace",
            "transcriptPath": "/tmp/transcript.jsonl",
            "branch": "worktree-ticket",
            "version": "test",
            "claimCounter": 0,
            "current": None,
            "pendingClaim": None,
            "pendingOperation": None,
            "pendingSuspend": None,
            "stopped": False,
            "turnActive": True,
            "activeTurnId": None,
            "daemonPid": os.getpid(),
            "daemonStartedAt": worker.time.time(),
            "sidecarHeartbeatAt": worker.time.time(),
        }
        self._store().save(state)
        return state

    def _add_permit(self, state, *, now=None, status="RUNNING",
                    lease_seconds=300):
        current = state.get("current")
        self.assertIsInstance(current, dict)
        issued_at = worker.time.time() if now is None else float(now)
        state["daemonPid"] = os.getpid()
        state["sidecarHeartbeatAt"] = issued_at
        state["sessionPermit"] = worker._session_permit(
            state, current, {
                "session": {
                    "status": status,
                    "leaseExpireAt": issued_at + lease_seconds,
                },
            }, source="test", now=issued_at)
        return state

    def test_session_hook_registers_private_non_pulling_worker_and_offlines(self):
        fake = FakeClient()
        event = {
            "hook_event_name": "SessionStart",
            "session_id": "native-thread-1",
            "cwd": self.temp.name,
            "source": "startup",
        }
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid", return_value=(os.getpid(), True)), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"), \
                mock.patch.object(worker, "_ensure_daemon") as daemon, \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", event), 0)
            self.assertEqual(worker.hook("codex", event), 0)
            worker.hook("codex", {**event, "hook_event_name": "SessionEnd"})

        state = self._store().load()
        self.assertTrue(state["stopped"])
        self.assertEqual(stat.S_IMODE(self._store().path.stat().st_mode), 0o600)
        self.assertNotIn("token", str(state).lower())
        registrations = [c for c in fake.calls if c[0] == "register_worker"]
        self.assertEqual(len(registrations), 2)
        payload = registrations[0][1][0]
        self.assertFalse(payload["capabilities"]["dispatch"]["pull"])
        self.assertEqual(payload["maxSlots"], 1)
        self.assertNotIn("lease_task", [c[0] for c in fake.calls])
        self.assertEqual(daemon.call_count, 2)
        self.assertEqual(fake.calls[-1][1][1]["status"], "OFFLINE")

    def test_failed_first_session_registration_persists_fail_closed_tombstone(self):
        fake = FakeClient()
        fake.register_worker = mock.Mock(
            side_effect=worker.ControlPlaneUnavailable("pre unavailable"))
        event = {
            "hook_event_name": "SessionStart",
            "session_id": "native-thread-1",
            "cwd": self.temp.name,
            "source": "startup",
        }
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid", return_value=(os.getpid(), True)), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.hook("codex", event), 0)
        state = self._store().load()
        self.assertTrue(state["stopped"])
        self.assertTrue(state["registrationFailed"])
        self.assertEqual(state["clientSessionId"], "native-thread-1")

    def test_subagent_start_binds_exact_spawn_and_followup_reauthorizes_epoch(self):
        parent_transcript = Path(self.temp.name) / "parent.jsonl"
        child_transcript = Path(self.temp.name) / "child.jsonl"
        parent_transcript.write_text(
            "\n".join((
                '{"type":"response_item","payload":{"type":"function_call",'
                '"name":"spawn_agent","call_id":"call-spawn",'
                '"internal_chat_message_metadata_passthrough":{"turn_id":"root-turn"}}}',
                '{"type":"event_msg","payload":{"type":"sub_agent_activity",'
                '"kind":"started","agent_thread_id":"agent-child-1",'
                '"event_id":"call-spawn"}}',
            )) + "\n", encoding="utf-8")
        child_transcript.write_text(
            '{"type":"session_meta","payload":{"id":"agent-child-1",'
            '"session_id":"native-thread-1",'
            '"cwd":"/workspace","source":{"subagent":{"thread_spawn":{'
            '"parent_thread_id":"native-thread-1",'
            '"agent_path":"/root/child"}}}}}\n', encoding="utf-8")
        state = self._seed()
        state["activeTurnId"] = "root-turn"
        state["transcriptPath"] = str(parent_transcript)
        self._store().save(state)
        spawn_event = {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "root-turn",
            "tool_use_id": "call-spawn",
            "tool_name": "spawn_agent",
            "tool_input": {"task_name": "child", "message": "review"},
        }
        with mock.patch.object(worker, "_calling_process_matches", return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._store(), "codex", spawn_event))
        self.assertIn("call-spawn", self._store().load()["subagentSpawnPermits"])

        start_event = {
            "hook_event_name": "SubagentStart",
            "session_id": "native-thread-1",
            "turn_id": "child-turn-1",
            "agent_id": "agent-child-1",
            "agent_type": "subagent",
            "cwd": "/workspace",
            "transcript_path": str(child_transcript),
        }
        with contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", start_event), 0)
        child = self._child_store().load()
        self.assertEqual(child["subagentBinding"]["rootSessionId"], "native-thread-1")
        root = self._store().load()
        self.assertEqual(
            root["subagentRegistry"]["/root/child"]["agentId"], "agent-child-1")
        self.assertEqual(root["subagentEpochs"]["agent-child-1"]["rootTurnId"],
                         "root-turn")

        root["activeTurnId"] = "root-next"
        self._store().save(root)
        child_tool = {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "child-turn-2",
            "agent_id": "agent-child-1",
            "agent_type": "subagent",
            "tool_use_id": "tool-child-next",
            "tool_name": "apply_patch",
            "tool_input": {},
        }
        with mock.patch.object(worker, "_calling_process_matches", return_value=True):
            self.assertIn("不属于当前根 turn/任务 fence",
                          worker._guard_pre_tool_use(
                              self._child_store(), "codex", child_tool))

        followup_event = {
            "session_id": "native-thread-1",
            "turn_id": "root-next",
            "tool_use_id": "call-followup",
            "tool_name": "collaborationfollowup_task",
            "tool_input": {"target": "/root/child", "message": "continue"},
        }
        with mock.patch.object(worker, "_calling_process_matches", return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._store(), "codex",
                {**followup_event, "hook_event_name": "PreToolUse"}))
        with mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                **followup_event,
                "hook_event_name": "PostToolUse",
                "tool_response": {"ok": True},
            }), 0)
        with mock.patch.object(worker, "_calling_process_matches", return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._child_store(), "codex", child_tool))

        root = self._store().load()
        root["current"] = {
            "aoneId": "84399999", "projectId": "2100304", "taskId": "task-new",
            "sessionId": "session-new", "fenceToken": 17, "generation": 8,
            "cycle": 1, "runtimeSessionId": "interactive:new",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(root)
        with mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                mock.patch.object(worker, "_client",
                                  side_effect=AssertionError("must block before heartbeat")):
            self.assertIn("不属于当前根 turn/任务 fence",
                          worker._guard_pre_tool_use(
                              self._child_store(), "codex", child_tool))

    def test_subagent_stop_blocks_delayed_tool_until_successful_followup(self):
        state = self._seed()
        state["activeTurnId"] = "root-turn"
        state["subagentRevision"] = 1
        state["subagentRegistry"] = {
            "/root/child": {
                "agentId": "agent-child-1",
                "bindingRevision": 1,
                "sourceToolUseId": "call-spawn",
            },
        }
        state["subagentLifecycles"] = {
            "agent-child-1": {"status": "ACTIVE", "revision": 1},
        }
        state["subagentEpochs"] = {
            "agent-child-1": {
                "rootTurnId": "root-turn",
                "assignmentEpoch": worker._assignment_epoch(state),
                "authorizedRevision": 1,
            },
        }
        self._store().save(state)
        self._child_store().save({
            "schemaVersion": 1,
            "client": "codex",
            "clientSessionId": "native-thread-1",
            "transcriptPath": "/tmp/child.jsonl",
            "subagentBinding": {
                "agentId": "agent-child-1",
                "agentType": "subagent",
                "agentPath": "/root/child",
                "parentThreadId": "native-thread-1",
                "rootSessionId": "native-thread-1",
                "bindingRevision": 1,
            },
            "stopped": False,
        })
        child_tool = {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "child-turn",
            "agent_id": "agent-child-1",
            "agent_type": "subagent",
            "tool_use_id": "tool-delayed",
            "tool_name": "apply_patch",
            "tool_input": {},
        }
        with contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "SubagentStop",
                "session_id": "native-thread-1",
                "turn_id": "child-turn",
                "agent_id": "agent-child-1",
                "agent_type": "subagent",
                "agent_transcript_path": "/tmp/child.jsonl",
            }), 0)
        stopped = self._child_store().load()
        self.assertTrue(stopped["stopped"])
        self.assertEqual(stopped["transcriptPath"], "/tmp/child.jsonl")
        self.assertNotIn(
            "agent-child-1", self._store().load()["subagentEpochs"])

        stderr = io.StringIO()
        with mock.patch.object(
                worker, "_client",
                side_effect=AssertionError("must block before heartbeat")), \
                mock.patch.object(worker, "_calling_process_matches",
                                  return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", child_tool), 2)
        self.assertIn("不属于当前根 turn/任务 fence", stderr.getvalue())

        followup = {
            "session_id": "native-thread-1",
            "turn_id": "root-turn",
            "tool_use_id": "call-followup-after-stop",
            "tool_name": "collaborationfollowup_task",
            "tool_input": {"target": "/root/child", "message": "continue"},
        }
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._store(), "codex",
                {**followup, "hook_event_name": "PreToolUse"}))
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                **followup,
                "hook_event_name": "PostToolUse",
                "tool_response": {"ok": True},
            }), 0)
        self.assertFalse(self._child_store().load()["stopped"])
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._child_store(), "codex",
                {**child_tool, "tool_use_id": "tool-after-followup"}))

    def test_followup_post_tool_cannot_authorize_across_assignment_epoch(self):
        state = self._seed()
        state["activeTurnId"] = "root-turn"
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._add_permit(state)
        state["subagentRevision"] = 2
        state["subagentRegistry"] = {
            "/root/child": {
                "agentId": "agent-child-1",
                "bindingRevision": 1,
            },
        }
        state["subagentLifecycles"] = {
            "agent-child-1": {"status": "STOPPED", "revision": 2},
        }
        self._store().save(state)
        self._child_store().save({
            "schemaVersion": 1,
            "client": "codex",
            "clientSessionId": "native-thread-1",
            "subagentBinding": {
                "agentId": "agent-child-1",
                "agentPath": "/root/child",
                "parentThreadId": "native-thread-1",
                "rootSessionId": "native-thread-1",
                "bindingRevision": 1,
            },
            "stopped": True,
            "subagentStoppedRevision": 2,
        })
        followup = {
            "session_id": "native-thread-1",
            "turn_id": "root-turn",
            "tool_use_id": "call-old-followup",
            "tool_name": "collaborationsend_message",
            "tool_input": {"target": "/root/child", "message": "continue"},
        }
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_calling_process_matches",
                                  return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._store(), "codex",
                {**followup, "hook_event_name": "PreToolUse"}))
        changed = self._store().load()
        changed["current"] = {
            **changed["current"],
            "fenceToken": 10,
        }
        self._store().save(changed)
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                **followup,
                "hook_event_name": "PostToolUse",
                "tool_response": {"ok": True},
            }), 0)
        latest = self._store().load()
        self.assertNotIn("agent-child-1", latest.get("subagentEpochs", {}))
        self.assertEqual(
            latest["subagentLifecycles"]["agent-child-1"]["status"], "STOPPED")
        self.assertTrue(self._child_store().load()["stopped"])

    def test_delayed_old_bind_cannot_reclaim_reused_agent_path(self):
        parent_transcript = Path(self.temp.name) / "parent-rebind.jsonl"
        parent_transcript.write_text(
            "\n".join((
                '{"type":"response_item","payload":{"type":"function_call",'
                '"name":"spawn_agent","call_id":"call-a",'
                '"internal_chat_message_metadata_passthrough":{"turn_id":"root-turn"}}}',
                '{"type":"event_msg","payload":{"type":"sub_agent_activity",'
                '"kind":"started","agent_thread_id":"agent-a",'
                '"event_id":"call-a"}}',
                '{"type":"response_item","payload":{"type":"function_call",'
                '"name":"spawn_agent","call_id":"call-b",'
                '"internal_chat_message_metadata_passthrough":{"turn_id":"root-turn"}}}',
                '{"type":"event_msg","payload":{"type":"sub_agent_activity",'
                '"kind":"started","agent_thread_id":"agent-b",'
                '"event_id":"call-b"}}',
            )) + "\n", encoding="utf-8")
        child_paths = {}
        for agent_id in ("agent-a", "agent-b"):
            path = Path(self.temp.name) / ("%s.jsonl" % agent_id)
            path.write_text(
                '{"type":"session_meta","payload":{"id":"%s",'
                '"session_id":"native-thread-1","cwd":"/workspace",'
                '"source":{"subagent":{"thread_spawn":{'
                '"parent_thread_id":"native-thread-1",'
                '"agent_path":"/root/reused"}}}}}\n' % agent_id,
                encoding="utf-8")
            child_paths[agent_id] = path
        state = self._seed()
        state["activeTurnId"] = "root-turn"
        state["transcriptPath"] = str(parent_transcript)
        self._store().save(state)
        for call_id in ("call-a", "call-b"):
            with mock.patch.object(worker, "_calling_process_matches",
                                   return_value=True):
                self.assertIsNone(worker._guard_pre_tool_use(
                    self._store(), "codex", {
                        "hook_event_name": "PreToolUse",
                        "session_id": "native-thread-1",
                        "turn_id": "root-turn",
                        "tool_use_id": call_id,
                        "tool_name": "collaborationspawn_agent",
                        "tool_input": {
                            "task_name": call_id,
                            "message": "review",
                        },
                    }))

        def start(agent_id):
            with contextlib.redirect_stdout(io.StringIO()):
                self.assertEqual(worker.hook("codex", {
                    "hook_event_name": "SubagentStart",
                    "session_id": "native-thread-1",
                    "turn_id": "%s-turn" % agent_id,
                    "agent_id": agent_id,
                    "agent_type": "subagent",
                    "cwd": "/workspace",
                    "transcript_path": str(child_paths[agent_id]),
                }), 0)

        start("agent-a")
        start("agent-b")
        start("agent-a")
        root = self._store().load()
        self.assertEqual(
            root["subagentRegistry"]["/root/reused"]["agentId"], "agent-b")
        self.assertEqual(
            root["subagentLifecycles"]["agent-a"]["status"], "STOPPED")

        followup = {
            "session_id": "native-thread-1",
            "turn_id": "root-turn",
            "tool_use_id": "call-followup-reused",
            "tool_name": "collaborationfollowup_task",
            "tool_input": {"target": "/root/reused", "message": "continue"},
        }
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._store(), "codex",
                {**followup, "hook_event_name": "PreToolUse"}))
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                **followup,
                "hook_event_name": "PostToolUse",
                "tool_response": {"ok": True},
            }), 0)
        root = self._store().load()
        self.assertNotIn("agent-a", root["subagentEpochs"])
        self.assertIn("agent-b", root["subagentEpochs"])

    def test_session_start_serializes_with_concurrent_assignment_mutation(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        state["schemaVersion"] = 1
        state.pop("sidecarHeartbeatAt", None)
        self._store().save(state)
        fake = FakeClient()
        register_entered = threading.Event()
        allow_register = threading.Event()
        mutation_done = threading.Event()

        def blocking_register(*args, **kwargs):
            register_entered.set()
            self.assertTrue(allow_register.wait(3))
            return fake._record("register_worker", *args, **kwargs)

        fake.register_worker = blocking_register
        hook_result = []

        def run_session_start():
            hook_result.append(worker.hook("codex", {
                "hook_event_name": "SessionStart",
                "session_id": "native-thread-1",
                "cwd": self.temp.name,
                "source": "resume",
            }))

        def write_new_assignment():
            with self._store().locked():
                latest = self._store().load_unlocked()
                latest["current"]["fenceToken"] = 10
                latest["pendingOperation"] = {
                    "operationKey": "new-op", "fenceToken": 10,
                }
                self._store().save_unlocked(latest)
            mutation_done.set()

        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid",
                                  return_value=(os.getpid(), True)), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="birth-serialized"), \
                mock.patch.object(worker, "make_worker_key",
                                  return_value=state["workerKey"]), \
                mock.patch.object(worker, "_ensure_daemon"), \
                mock.patch.object(worker, "_hook_output"):
            session_thread = threading.Thread(target=run_session_start)
            session_thread.start()
            self.assertTrue(register_entered.wait(3))
            mutation_thread = threading.Thread(target=write_new_assignment)
            mutation_thread.start()
            self.assertFalse(mutation_done.wait(0.1))
            allow_register.set()
            session_thread.join(3)
            mutation_thread.join(3)

        self.assertFalse(session_thread.is_alive())
        self.assertFalse(mutation_thread.is_alive())
        self.assertEqual(hook_result, [0])
        after = self._store().load()
        self.assertEqual(after["current"]["fenceToken"], 10)
        self.assertEqual(after["pendingOperation"]["operationKey"], "new-op")

    def test_codex_turn_hooks_are_local_and_ignore_delayed_stop(self):
        self._seed()
        with mock.patch.object(worker, "_client",
                               side_effect=AssertionError("network must not be used")), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "UserPromptSubmit",
                "session_id": "native-thread-1",
                "turn_id": "turn-new",
            }), 0)
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "Stop",
                "session_id": "native-thread-1",
                "turn_id": "turn-old",
            }), 0)
        active = self._store().load()
        self.assertTrue(active["turnActive"])
        self.assertEqual(active["activeTurnId"], "turn-new")
        self.assertNotIn("turnStoppedAt", active)

        with mock.patch.object(worker, "_client",
                               side_effect=AssertionError("network must not be used")), \
                mock.patch.object(worker.time, "time", return_value=1234), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "Stop",
                "session_id": "native-thread-1",
                "turn_id": "turn-new",
            }), 0)
        stopped = self._store().load()
        self.assertFalse(stopped["turnActive"])
        self.assertEqual(stopped["turnStoppedAt"], 1234)

    def test_session_heartbeat_uses_codex_turn_grace_but_not_for_claude(self):
        state = self._seed()
        state["turnActive"] = False
        state["turnStoppedAt"] = 100
        with mock.patch.dict(os.environ, {
            "JARVIS_INTERACTIVE_TURN_GRACE_SEC": "600",
        }, clear=False):
            self.assertTrue(worker._session_heartbeat_allowed(state, now=699.9))
            self.assertFalse(worker._session_heartbeat_allowed(state, now=700))
        state["client"] = "claude"
        self.assertTrue(worker._session_heartbeat_allowed(state, now=10000))

    def test_interactive_timing_defaults_and_explicit_turn_grace(self):
        keys = (
            "JARVIS_INTERACTIVE_LEASE_SECONDS",
            "JARVIS_INTERACTIVE_HEARTBEAT_SEC",
            "JARVIS_INTERACTIVE_TURN_GRACE_SEC",
            "JARVIS_INTERACTIVE_LEASE_SAFETY_MARGIN_SEC",
            "JARVIS_INTERACTIVE_AFFINITY_SEC",
        )
        with mock.patch.dict(os.environ, {}, clear=False):
            old = {key: os.environ.pop(key, None) for key in keys}
            try:
                self.assertEqual(worker._interactive_lease_seconds(), 300)
                self.assertEqual(worker._interactive_heartbeat_seconds(), 30)
                self.assertEqual(worker._interactive_turn_grace_seconds(), 600)
                self.assertEqual(
                    worker._interactive_lease_safety_margin_seconds(), 90)
                self.assertEqual(worker._interactive_affinity_seconds(), 7200)

                os.environ["JARVIS_INTERACTIVE_TURN_GRACE_SEC"] = "90"
                self.assertEqual(worker._interactive_turn_grace_seconds(), 90)
                os.environ[
                    "JARVIS_INTERACTIVE_LEASE_SAFETY_MARGIN_SEC"] = "75"
                self.assertEqual(
                    worker._interactive_lease_safety_margin_seconds(), 75)
                os.environ["JARVIS_INTERACTIVE_AFFINITY_SEC"] = "3600"
                self.assertEqual(worker._interactive_affinity_seconds(), 3600)
            finally:
                for key in keys:
                    os.environ.pop(key, None)
                for key, value in old.items():
                    if value is not None:
                        os.environ[key] = value

    def test_session_permit_accepts_iso_or_epoch_millis_and_has_safe_fallback(self):
        state = self._seed()
        state["current"] = {
            "taskId": "task-1", "sessionId": "session-1",
            "fenceToken": 9, "generation": 4,
            "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 300,
        }
        iso = worker._session_permit(
            state, state["current"], {
                "session": {
                    "status": "RUNNING",
                    "leaseExpireAt": "2026-07-16T12:00:00Z",
                },
            }, source="test", now=100)
        self.assertEqual(
            iso["leaseExpireAt"],
            worker._epoch_seconds("2026-07-16T12:00:00Z"))
        millis = worker._session_permit(
            state, state["current"], {
                "session": {
                    "status": "RUNNING",
                    "leaseExpireAt": 1_800_000_000_000,
                },
            }, source="test", now=100)
        self.assertEqual(millis["leaseExpireAt"], 1_800_000_000)
        fallback = worker._session_permit(
            state, state["current"], {"ok": True},
            source="test", now=100)
        self.assertEqual(fallback["leaseExpireAt"], 400)
        self.assertEqual(fallback["sessionStatus"], "RUNNING")

    def test_missing_stop_has_active_turn_ttl_before_auto_suspend(self):
        state = self._seed()
        state["turnActive"] = True
        state["lastTurnActivityAt"] = 100
        with mock.patch.dict(os.environ, {
            "JARVIS_INTERACTIVE_ACTIVE_TURN_TTL_SEC": "60",
        }, clear=False):
            self.assertTrue(worker._session_heartbeat_allowed(state, now=159.9))
            self.assertFalse(worker._session_heartbeat_allowed(state, now=160))
        with mock.patch.dict(os.environ, {
            "JARVIS_INTERACTIVE_ACTIVE_TURN_TTL_SEC": "0",
        }, clear=False):
            self.assertTrue(worker._session_heartbeat_allowed(state, now=10000))

        self._store().save(state)
        with mock.patch.object(worker.time, "time", return_value=150), \
                mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "PostToolUse",
                "session_id": "native-thread-1",
                "turn_id": "long-turn",
                "tool_name": "exec_command",
            }), 0)
        touched = self._store().load()
        self.assertEqual(touched["lastTurnActivityAt"], 150)
        with mock.patch.dict(os.environ, {
            "JARVIS_INTERACTIVE_ACTIVE_TURN_TTL_SEC": "60",
        }, clear=False):
            self.assertTrue(worker._session_heartbeat_allowed(touched, now=209.9))
            self.assertFalse(worker._session_heartbeat_allowed(touched, now=210))

        state["lastTurnActivityAt"] = 100
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.dict(os.environ, {
            "JARVIS_INTERACTIVE_ACTIVE_TURN_TTL_SEC": "60",
        }, clear=False), \
                mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", side_effect=[True, False]), \
                mock.patch.object(worker.time, "time", return_value=1000), \
                mock.patch.object(worker.time, "sleep", return_value=None):
            self.assertEqual(worker.daemon(
                self._store().path, state["workerKey"]), 0)
        self.assertIn("suspend_session", [c[0] for c in fake.calls])
        self.assertIsNone(self._store().load()["current"])

    def test_user_prompt_restarts_sidecar_and_uses_local_permit(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._add_permit(state)
        self._store().save(state)
        fake = FakeClient()
        event = {
            "hook_event_name": "UserPromptSubmit",
            "session_id": "native-thread-1",
            "turn_id": "turn-resume",
        }
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                mock.patch.object(worker, "_ensure_daemon") as ensure, \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", event), 0)
        ensure.assert_called_once()
        self.assertEqual(ensure.call_args.args[0].path, self._store().path)
        self.assertEqual(ensure.call_args.args[1], state["workerKey"])
        self.assertNotIn("heartbeat_worker", [c[0] for c in fake.calls])
        self.assertNotIn("heartbeat_session", [c[0] for c in fake.calls])
        self.assertTrue(self._store().load()["turnActive"])

    def test_sidecar_clears_stale_fence_before_work(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._add_permit(state)
        self._store().save(state)
        fake = FakeClient()
        fake.heartbeat_error = worker.StaleFence("reassigned")
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive",
                                  side_effect=[True, False]), \
                mock.patch.object(worker.time, "sleep", return_value=None):
            self.assertEqual(worker.daemon(
                self._store().path, state["workerKey"]), 0)
        after = self._store().load()
        self.assertIsNone(after["current"])
        self.assertIsNone(after["pendingClaim"])
        self.assertEqual(after["lostOwnership"]["aoneId"], "84345050")

    def test_worker_404_tombstones_old_assignment_and_standard_claim_rebuilds(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-old",
            "sessionId": "session-old", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 300, "heartbeatEnabled": True,
        }
        self._add_permit(state)
        self._store().save(state)
        fake = FakeClient()
        fake.worker_heartbeat_error = worker.ControlPlaneError(
            "Worker not found", status=404, code="NotFound.Worker")
        stderr = io.StringIO()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.daemon(
                self._store().path, state["workerKey"]), 0)

        cleaned = self._store().load()
        self.assertTrue(cleaned["stopped"])
        self.assertEqual(cleaned["offlineReason"], "admin_cleanup")
        self.assertIsNone(cleaned["current"])
        self.assertNotIn("sessionPermit", cleaned)
        self.assertNotIn("daemonPid", cleaned)
        self.assertEqual(
            cleaned["lostOwnership"]["cause"], "ADMIN_CLEANUP")
        self.assertEqual(
            cleaned["lostOwnership"]["missingResource"], "Worker")
        self.assertIn("admin cleanup", cleaned["lostOwnership"]["reason"])
        self.assertNotIn(
            "heartbeat_session", [call[0] for call in fake.calls])
        self.assertIn("admin cleanup removed Worker", stderr.getvalue())

        fake.worker_heartbeat_error = None
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(
                    worker, "_find_host_pid", return_value=(os.getpid(), True)), \
                mock.patch.object(
                    worker, "_process_start_identity",
                    return_value=state["hostProcessStartedAt"]), \
                mock.patch.object(
                    worker, "_default_boot_id", return_value=state["bootId"]), \
                mock.patch.object(
                    worker, "make_worker_key", return_value=state["workerKey"]), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "SessionStart",
                "session_id": "native-thread-1",
                "cwd": self.temp.name,
                "source": "resume",
            }), 0)
            rebuilt = worker.prepare_claim("84345050", "2100304")

        after = self._store().load()
        self.assertFalse(after["stopped"])
        self.assertNotIn("lostOwnership", after)
        self.assertEqual(after["current"]["sessionId"], rebuilt["sessionId"])
        self.assertEqual(after["current"]["fenceToken"], rebuilt["fenceToken"])

    def test_session_404_tombstones_old_assignment_and_stops_sidecar(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-old",
            "sessionId": "session-old", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 300, "heartbeatEnabled": True,
        }
        self._add_permit(state)
        self._store().save(state)
        fake = FakeClient()
        fake.heartbeat_error = worker.ControlPlaneError(
            "Session not found", status=404, code="NotFound.Session")
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.daemon(
                self._store().path, state["workerKey"]), 0)

        cleaned = self._store().load()
        self.assertTrue(cleaned["stopped"])
        self.assertIsNone(cleaned["current"])
        self.assertNotIn("sessionPermit", cleaned)
        self.assertEqual(
            cleaned["lostOwnership"]["missingResource"], "Session")
        self.assertEqual(
            [call[0] for call in fake.calls].count("heartbeat_session"), 1)

    def test_pre_tool_use_uses_local_permit_and_blocks_mismatch(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._add_permit(state)
        self._store().save(state)
        event = {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "turn-tool",
            "tool_use_id": "tool-1",
            "tool_name": "Bash",
            "tool_input": {"command": "bin/a1id -- project workitem get 84345050"},
        }
        stdout = io.StringIO()
        stderr = io.StringIO()
        with mock.patch.object(
                worker, "_client",
                side_effect=AssertionError("PreToolUse must stay local")), \
                mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(stdout), contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", event), 0)
        self.assertEqual(stdout.getvalue().strip(), "{}")
        self.assertEqual(stderr.getvalue(), "")

        changed = self._store().load()
        changed["sessionPermit"]["fenceToken"] = 8
        self._store().save(changed)
        stdout = io.StringIO()
        stderr = io.StringIO()
        with mock.patch.object(
                worker, "_client",
                side_effect=AssertionError("PreToolUse must stay local")), \
                mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(stdout), contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", {**event, "tool_use_id": "tool-2"}), 2)
        after = self._store().load()
        self.assertEqual(stdout.getvalue(), "")
        self.assertIn("当前工具调用已阻断", stderr.getvalue())
        self.assertEqual(after["current"]["fenceToken"], 9)
        self.assertNotIn("lostOwnership", after)

    def test_local_permit_fails_closed_on_margin_status_and_sidecar_health(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 300, "heartbeatEnabled": True,
        }
        event = {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "turn-tool",
            "tool_use_id": "tool-local-proof",
            "tool_name": "Bash",
            "tool_input": {"command": "git status --short"},
        }
        for label, mutate, expected in (
                ("safety-margin",
                 lambda value: value["sessionPermit"].update(
                     {"leaseExpireAt": 1090}),
                 "安全边界"),
                ("session-status",
                 lambda value: value["sessionPermit"].update(
                     {"sessionStatus": "SUSPENDED"}),
                 "SUSPENDED"),
                ("sidecar-dead",
                 lambda value: value.update({"daemonPid": 99999999}),
                 "sidecar 未运行")):
            with self.subTest(label=label):
                candidate = dict(state)
                candidate["current"] = dict(state["current"])
                self._add_permit(candidate, now=1000, lease_seconds=300)
                candidate["sessionPermit"] = dict(candidate["sessionPermit"])
                mutate(candidate)
                self._store().save(candidate)
                stderr = io.StringIO()
                with mock.patch.object(worker.time, "time", return_value=1000), \
                        mock.patch.object(
                            worker, "_client",
                            side_effect=AssertionError(
                                "PreToolUse must stay local")), \
                        mock.patch.object(
                            worker, "_calling_process_matches",
                            return_value=True), \
                        contextlib.redirect_stdout(io.StringIO()), \
                        contextlib.redirect_stderr(stderr):
                    self.assertEqual(worker.hook("codex", event), 2)
                self.assertIn(expected, stderr.getvalue())

    def test_claim_start_permit_allows_first_tool_without_remote_heartbeat(self):
        self._seed()
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            claimed = worker.prepare_claim("84345050", "2100304")
            worker.acknowledge_claim(
                "84345050", "aone:2100304:84345050:tag")
        state = self._store().load()
        self.assertEqual(
            state["sessionPermit"]["sessionId"], claimed["sessionId"])
        self.assertEqual(state["sessionPermit"]["source"], "session-start")
        before = len(fake.calls)
        with mock.patch.object(
                worker, "_client",
                side_effect=AssertionError("PreToolUse must stay local")), \
                mock.patch.object(
                    worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "PreToolUse",
                "session_id": "native-thread-1",
                "turn_id": "turn-tool",
                "tool_use_id": "tool-first-after-claim",
                "tool_name": "Bash",
                "tool_input": {"command": "git status --short"},
            }), 0)
        self.assertEqual(len(fake.calls), before)

    def test_pre_tool_use_blocks_lost_claude_worker_too(self):
        state = self._seed()
        state["client"] = "claude"
        state["clientSessionId"] = "claude-session-1"
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        store = worker.StateStore(worker._state_path("claude", "claude-session-1"))
        store.save(state)
        worker._clear_lost_current(
            store, state["workerKey"], "session-1", 9, "reassigned")
        stdout = io.StringIO()
        stderr = io.StringIO()
        with mock.patch.object(worker, "_client",
                               side_effect=AssertionError("network must not be used")), \
                mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(stdout), contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("claude", {
                "hook_event_name": "PreToolUse",
                "session_id": "claude-session-1",
                "tool_use_id": "tool-claude",
                "tool_name": "Edit",
                "tool_input": {"file_path": "/workspace/main.py"},
            }), 2)
        self.assertEqual(stdout.getvalue(), "")
        self.assertIn("已失去 Aone 84345050", stderr.getvalue())

    def test_claude_subagent_tool_uses_root_worker_and_assignment(self):
        store = worker.StateStore(
            worker._state_path("claude", "claude-session-1"))
        state = self._seed()
        state["client"] = "claude"
        state["clientSessionId"] = "claude-session-1"
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:claude:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._add_permit(state)
        store.save(state)
        with mock.patch.object(
                worker, "_client",
                side_effect=AssertionError("PreToolUse must stay local")), \
                mock.patch.object(worker, "_calling_process_matches",
                                  return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.hook("claude", {
                "hook_event_name": "PreToolUse",
                "session_id": "claude-session-1",
                "agent_id": "claude-child-1",
                "agent_type": "jarvis-hook-probe",
                "tool_use_id": "tool-child-bash",
                "tool_name": "Bash",
                "tool_input": {"command": "git status --short"},
            }), 0)
        self.assertFalse(
            worker._state_path("claude", "claude-child-1").exists())

    def test_pre_tool_use_only_allows_exact_standard_claim_recovery(self):
        state = self._seed()
        state["lostOwnership"] = {
            "aoneId": "84345050", "projectId": "2100304",
            "sessionId": "session-old", "reason": "stale fence",
        }
        self._store().save(state)
        exact_claim = {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "turn-recover",
            "tool_use_id": "tool-claim",
            "tool_name": "Bash",
            "tool_input": {"command": "%s %s claim 84345050 2100304" % (
                "/bin/bash", worker.REPO_ROOT / "bootstrap" / "claim.sh")},
        }
        stdout = io.StringIO()
        with mock.patch.object(worker, "_client",
                               side_effect=AssertionError("claim gate is local")), \
                mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(stdout):
            self.assertEqual(worker.hook("codex", exact_claim), 0)
        self.assertEqual(stdout.getvalue().strip(), "{}")

        claim_script = worker.REPO_ROOT / "bootstrap" / "claim.sh"
        rejected = (
            ("Bash", "/bin/bash %s claim 84345050 2100304; a1 update" % claim_script),
            ("Bash", "/bin/bash %s claim 84345050 2100304 && a1 update" % claim_script),
            ("Bash", "/bin/bash %s claim 84345050 2100304 | tee /tmp/x" % claim_script),
            ("Bash", "/bin/bash %s claim 84345050 2100304\na1 update" % claim_script),
            ("Bash", "/bin/bash %s claim 84345050 2100304 $(a1 update)" % claim_script),
            ("Bash", "/bin/bash %s claim 84345051 2100304" % claim_script),
            ("Bash", "env BASH_ENV=/tmp/pwn /bin/bash %s claim 84345050 2100304" % claim_script),
            ("Bash", "JARVIS_INTERACTIVE_WORKER_MANAGER=/tmp/fake /bin/bash %s claim 84345050 2100304" % claim_script),
            ("Bash", "/tmp/bash %s claim 84345050 2100304" % claim_script),
            ("Bash", "/bin/bash /tmp/evil/bootstrap/claim.sh claim 84345050 2100304"),
            ("mcp__evil__execute", "/bin/bash %s claim 84345050 2100304" % claim_script),
        )
        for index, (tool_name, command) in enumerate(rejected):
            stderr = io.StringIO()
            with self.subTest(command=command), \
                    mock.patch.object(worker, "_calling_process_matches",
                                      return_value=True), \
                    contextlib.redirect_stdout(io.StringIO()), \
                    contextlib.redirect_stderr(stderr):
                self.assertEqual(worker.hook("codex", {
                    **exact_claim,
                    "tool_use_id": "tool-rejected-%d" % index,
                    "tool_name": tool_name,
                    "tool_input": {"command": command},
                }), 2)
            self.assertIn("工具调用已阻断", stderr.getvalue())

        state = self._store().load()
        state.pop("lostOwnership", None)
        state["pendingClaim"] = {
            "aoneId": "84345050", "projectId": "2100304",
            "phase": "CLAIMING", "receiptUnknown": True,
        }
        self._store().save(state)
        stderr = io.StringIO()
        with mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", {
                **exact_claim,
                "tool_use_id": "tool-edit",
                "tool_name": "Write",
                "tool_input": {"file_path": "/workspace/main.py"},
            }), 2)
        self.assertIn("回执处于 UNKNOWN", stderr.getvalue())

    def test_pre_tool_use_blocks_old_turn_process_and_uncertain_state(self):
        state = self._seed()
        state["activeTurnId"] = "turn-new"
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        event = {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "turn-new",
            "tool_use_id": "tool-guard",
            "tool_name": "apply_patch",
            "tool_input": {"patch": "*** Begin Patch"},
        }
        for label, process_matches, turn_id in (
                ("old-process", False, "turn-new"),
                ("old-turn", True, "turn-old")):
            stderr = io.StringIO()
            with self.subTest(label=label), \
                    mock.patch.object(worker, "_calling_process_matches",
                                      return_value=process_matches), \
                    mock.patch.object(worker, "_client",
                                      side_effect=AssertionError("must block locally")), \
                    contextlib.redirect_stdout(io.StringIO()), \
                    contextlib.redirect_stderr(stderr):
                self.assertEqual(worker.hook("codex", {
                    **event, "turn_id": turn_id,
                }), 2)
            self.assertIn("已阻断", stderr.getvalue())

        state = self._store().load()
        state["turnActive"] = False
        state["pendingSuspend"] = {
            "sessionId": "session-1", "fenceToken": 9,
        }
        self._store().save(state)
        stderr = io.StringIO()
        with mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", event), 2)
        self.assertIn("turn 已停止", stderr.getvalue())

    def test_pre_tool_use_allows_subagent_turn_but_rechecks_root_stop(self):
        state = self._seed()
        state["activeTurnId"] = "root-turn"
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._add_permit(state)
        state["subagentEpochs"] = {
            "agent-child-1": {
                "rootTurnId": "root-turn",
                "assignmentEpoch": worker._assignment_epoch(state),
                "authorizedRevision": 1,
            },
        }
        state["subagentLifecycles"] = {
            "agent-child-1": {"status": "ACTIVE", "revision": 1},
        }
        self._store().save(state)
        self._child_store().save({
            "schemaVersion": 1,
            "client": "codex",
            "clientSessionId": "native-thread-1",
            "transcriptPath": "/tmp/child.jsonl",
            "subagentBinding": {
                "agentId": "agent-child-1",
                "agentType": "subagent",
                "agentPath": "/root/child",
                "parentThreadId": "native-thread-1",
                "rootSessionId": "native-thread-1",
                "bindingRevision": 1,
            },
            "stopped": False,
        })
        subagent_event = {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "child-turn-context",
            "agent_id": "agent-child-1",
            "agent_type": "subagent",
            "tool_use_id": "tool-child",
            "tool_name": "Bash",
            "tool_input": {"command": "git status --short"},
        }
        with mock.patch.object(
                worker, "_client",
                side_effect=AssertionError("PreToolUse must stay local")), \
                mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.hook("codex", subagent_event), 0)

        def stop_after_first_permit_check(*args, **kwargs):
            with self._store().locked():
                latest = self._store().load_unlocked()
                latest["turnActive"] = False
                latest["turnStoppedAt"] = 1234
                self._store().save_unlocked(latest)
            return None

        stderr = io.StringIO()
        with mock.patch.object(
                worker, "_session_permit_block_reason",
                side_effect=stop_after_first_permit_check), \
                mock.patch.object(
                    worker, "_client",
                    side_effect=AssertionError("PreToolUse must stay local")), \
                mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", {
                **subagent_event,
                "tool_use_id": "tool-racing-stop",
            }), 2)
        self.assertIn("turn 已停止", stderr.getvalue())

    def test_pre_tool_use_blocks_unregistered_session_and_old_subagent_epoch(self):
        event = {
            "hook_event_name": "PreToolUse",
            "session_id": "missing-session",
            "turn_id": "turn-1",
            "tool_use_id": "tool-unregistered",
            "tool_name": "a1-update",
            "tool_input": {},
        }
        stderr = io.StringIO()
        with mock.patch.object(worker, "_client",
                               side_effect=AssertionError("must block locally")), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", event), 2)
        self.assertIn("SessionStart", stderr.getvalue())

        state = self._seed()
        state["activeTurnId"] = "root-new"
        state["current"] = {
            "aoneId": "84399999", "projectId": "2100304", "taskId": "task-new",
            "sessionId": "session-new", "fenceToken": 17, "generation": 8,
            "cycle": 2, "runtimeSessionId": "interactive:cycle:2",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        state["subagentEpochs"] = {
            "agent-child-1": {
                "rootTurnId": "root-old",
                "assignmentEpoch": "session:session-old:9:4",
                "authorizedRevision": 1,
            },
        }
        state["subagentLifecycles"] = {
            "agent-child-1": {"status": "ACTIVE", "revision": 1},
        }
        self._store().save(state)
        self._child_store().save({
            "schemaVersion": 1,
            "client": "codex",
            "clientSessionId": "native-thread-1",
            "transcriptPath": "/tmp/child.jsonl",
            "subagentBinding": {
                "agentId": "agent-child-1",
                "agentType": "subagent",
                "agentPath": "/root/child",
                "parentThreadId": "native-thread-1",
                "rootSessionId": "native-thread-1",
                "bindingRevision": 1,
            },
            "stopped": False,
        })
        stderr = io.StringIO()
        with mock.patch.object(worker, "_calling_process_matches", return_value=True), \
                mock.patch.object(worker, "_client",
                                  side_effect=AssertionError("must block before heartbeat")), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", {
                **event,
                "session_id": "native-thread-1",
                "turn_id": "child-from-root-old",
                "agent_id": "agent-child-1",
                "agent_type": "subagent",
                "tool_use_id": "tool-old-child",
                "tool_name": "apply_patch",
            }), 2)
        self.assertIn("不属于当前根 turn/任务 fence", stderr.getvalue())

    def test_legacy_current_without_permit_blocks_locally_until_sidecar_refresh(self):
        state = self._seed()
        state["activeTurnId"] = "turn-tool"
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304",
            "taskId": "task-1", "sessionId": "session-1",
            "fenceToken": 10, "generation": 4, "cycle": 1,
            "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        stderr = io.StringIO()
        with mock.patch.object(
                worker, "_client",
                side_effect=AssertionError("PreToolUse must stay local")), \
                mock.patch.object(
                    worker, "_calling_process_matches", return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(stderr):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "PreToolUse",
                "session_id": "native-thread-1",
                "turn_id": "turn-tool",
                "tool_use_id": "tool-legacy-state",
                "tool_name": "McpWrite",
                "tool_input": {},
            }), 2)
        after = self._store().load()
        self.assertEqual(after["current"]["fenceToken"], 10)
        self.assertNotIn("lostOwnership", after)
        self.assertIn("Lease Proof", stderr.getvalue())

        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(
                    worker, "_host_alive", side_effect=[True, False]), \
                mock.patch.object(worker.time, "sleep", return_value=None):
            self.assertEqual(worker.daemon(
                self._store().path, state["workerKey"]), 0)
        upgraded = self._store().load()
        self.assertIn("sidecarHeartbeatAt", upgraded)
        self.assertEqual(
            upgraded["sessionPermit"]["source"], "sidecar-heartbeat")

    def test_delayed_old_fence_cannot_clear_new_same_session_assignment(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 10, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        worker._clear_lost_current(
            self._store(), state["workerKey"], "session-1", 9,
            "delayed stale heartbeat")
        after = self._store().load()
        self.assertEqual(after["current"]["fenceToken"], 10)
        self.assertNotIn("lostOwnership", after)

    def test_session_start_preserves_lost_ownership_across_incarnations(self):
        fake = FakeClient()
        event = {
            "hook_event_name": "SessionStart",
            "session_id": "native-thread-1",
            "cwd": self.temp.name,
            "source": "resume",
        }
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid", return_value=(os.getpid(), True)), \
                mock.patch.object(worker, "_process_start_identity", return_value="birth-1"), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot-1"), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", event), 0)
            state = self._store().load()
            state["lostOwnership"] = {
                "aoneId": "84345050",
                "projectId": "2100304",
                "runtimeSessionId": "interactive:cycle:1",
                "reason": "stale fence",
            }
            self._store().save(state)
            self.assertEqual(worker.hook("codex", event), 0)
        same = self._store().load()
        self.assertEqual(same["workerKey"], state["workerKey"])
        self.assertEqual(same["lostOwnership"]["aoneId"], "84345050")

        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid", return_value=(os.getpid(), True)), \
                mock.patch.object(worker, "_process_start_identity", return_value="birth-2"), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot-2"), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", event), 0)
        replacement = self._store().load()
        self.assertNotEqual(replacement["workerKey"], same["workerKey"])
        self.assertEqual(replacement["lostOwnership"]["aoneId"], "84345050")

        output = io.StringIO()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(output):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "UserPromptSubmit",
                "session_id": "native-thread-1",
                "turn_id": "turn-lost-after-restart",
                "prompt": "继续",
            }), 0)
        self.assertIn("严禁继续", output.getvalue())

    def test_idle_codex_worker_offlines_and_next_prompt_reactivates_it(self):
        state = self._seed()
        state["lastTurnActivityAt"] = 100
        state["turnActive"] = False
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.dict(os.environ, {
            "JARVIS_INTERACTIVE_IDLE_TTL_SEC": "60",
        }, clear=False), \
                mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                mock.patch.object(worker.time, "time", return_value=200):
            self.assertEqual(worker.daemon(self._store().path, state["workerKey"]), 0)
        idle = self._store().load()
        self.assertTrue(idle["stopped"])
        self.assertEqual(idle["offlineReason"], "idle_ttl")

        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                mock.patch.object(worker, "_ensure_daemon") as ensure, \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "UserPromptSubmit",
                "session_id": "native-thread-1",
                "turn_id": "turn-after-idle",
            }), 0)
        reactivated = self._store().load()
        self.assertFalse(reactivated["stopped"])
        self.assertNotIn("offlineReason", reactivated)
        ensure.assert_called_once()

    def test_stop_state_write_failure_blocks_stop(self):
        with mock.patch.object(worker, "_record_codex_turn",
                               side_effect=OSError("disk full")), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "Stop",
                "session_id": "native-thread-1",
                "turn_id": "turn-stop",
            }), worker.STATE_ERROR_EXIT)

    def test_unverified_host_registers_offline_without_sidecar_and_cannot_claim(self):
        fake = FakeClient()
        event = {
            "hook_event_name": "SessionStart",
            "session_id": "native-thread-1",
            "cwd": self.temp.name,
            "source": "startup",
        }
        output = io.StringIO()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid", return_value=(os.getpid(), False)), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"), \
                mock.patch.object(worker, "_ensure_daemon") as daemon, \
                contextlib.redirect_stdout(output):
            self.assertEqual(worker.hook("codex", event), 0)

        state = self._store().load()
        self.assertTrue(state["stopped"])
        self.assertFalse(state["verifyHostCommand"])
        daemon.assert_not_called()
        registration = [c for c in fake.calls if c[0] == "register_worker"][-1]
        self.assertEqual(registration[1][0]["status"], "OFFLINE")
        self.assertIn("OFFLINE", output.getvalue())
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaisesRegex(RuntimeError, "not active"):
                worker.prepare_claim("84345050", "2100304")
        self.assertNotIn("claim_task", [c[0] for c in fake.calls])

    def test_delayed_session_end_does_not_offline_replacement_incarnation(self):
        fake = FakeClient()
        event = {
            "hook_event_name": "SessionStart",
            "session_id": "native-thread-1",
            "cwd": self.temp.name,
            "source": "startup",
        }
        old_pid = 101
        new_pid = 202
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid",
                                  side_effect=[(old_pid, True), (new_pid, True),
                                               (old_pid, True)]), \
                mock.patch.object(worker, "_process_start_identity",
                                  side_effect=["old-birth", "new-birth", "old-birth"]), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", event), 0)
            old_worker_key = self._store().load()["workerKey"]
            self.assertEqual(worker.hook("codex", {**event, "source": "resume"}), 0)
            replacement = self._store().load()
            self.assertNotEqual(replacement["workerKey"], old_worker_key)
            self.assertEqual(worker.hook(
                "codex", {**event, "hook_event_name": "SessionEnd"}), 0)

        after_delayed_end = self._store().load()
        self.assertEqual(after_delayed_end["workerKey"], replacement["workerKey"])
        self.assertFalse(after_delayed_end["stopped"])
        offline_heartbeats = [
            c for c in fake.calls
            if c[0] == "heartbeat_worker" and c[1][1]["status"] == "OFFLINE"
        ]
        self.assertEqual(len(offline_heartbeats), 1)
        self.assertEqual(offline_heartbeats[0][1][0], old_worker_key)

    def test_headless_preexec_and_session_start_preserve_recovery_fence(self):
        session_id = "headless-session-1"
        store = self._claude_store(session_id)
        old_state = {
            "schemaVersion": 1,
            "client": "claude",
            "clientSessionId": session_id,
            "workerKey": "host:boot:old-process",
            "host": "host",
            "bootId": "boot",
            "processUuid": "old-process",
            "hostPid": 101,
            "hostProcessStartedAt": "old-birth",
            "verifyHostCommand": True,
            "cwd": self.temp.name,
            "branch": "worktree-ticket",
            "claimCounter": 7,
            "current": {
                "aoneId": "84382166",
                "projectId": "2100304",
                "taskId": "task-old",
                "sessionId": "session-old",
                "fenceToken": 9,
                "generation": 4,
                "cycle": 7,
                "runtimeSessionId": "interactive:cycle:7",
                "leaseSeconds": 120,
                "heartbeatEnabled": False,
            },
            "pendingClaim": None,
            "pendingOperation": {
                "operationId": "operation-old",
                "operationKey": "claim-operation-old",
                "aoneId": "84382166",
                "status": "SENDING",
            },
            "pendingSuspend": None,
            "lostOwnership": {
                "aoneId": "84382166",
                "projectId": "2100304",
                "reason": "stale fence",
            },
            "stopped": True,
            "stoppedAt": 100,
            "offlineReason": "host_stopped",
            "turnActive": True,
        }
        store.save(old_state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_pid_alive", return_value=True), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="new-birth"), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"):
            registered = worker.register_headless(
                session_id, 202, client_name="claude")

        prestarted = store.load()
        self.assertFalse(registered["sameIncarnation"])
        self.assertIsNone(prestarted["current"])
        self.assertIsNone(prestarted["pendingOperation"])
        self.assertEqual(prestarted["pendingClaim"]["phase"], "READY_TO_RECOVER")
        self.assertEqual(prestarted["pendingClaim"]["runtimeSessionId"],
                         "interactive:cycle:7")
        self.assertEqual(prestarted["pendingClaim"]["operationKey"],
                         "claim-operation-old")
        self.assertTrue(prestarted["pendingClaim"]["receiptUnknown"])
        self.assertEqual(prestarted["recoveryPending"]["oldWorkerKey"],
                         old_state["workerKey"])
        self.assertEqual(prestarted["lostOwnership"]["reason"], "stale fence")

        start_event = {
            "hook_event_name": "SessionStart",
            "session_id": session_id,
            "cwd": self.temp.name,
            "source": "resume",
        }
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid", return_value=(202, True)), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="new-birth"), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("claude", start_event), 0)

        refreshed = store.load()
        self.assertEqual(refreshed["workerKey"], prestarted["workerKey"])
        self.assertEqual(refreshed["pendingClaim"], prestarted["pendingClaim"])
        self.assertEqual(refreshed["recoveryPending"],
                         prestarted["recoveryPending"])
        self.assertEqual(refreshed["lostOwnership"],
                         prestarted["lostOwnership"])

        ordinary = {
            "hook_event_name": "PreToolUse",
            "session_id": session_id,
            "tool_use_id": "ordinary-tool",
            "tool_name": "Bash",
            "tool_input": {"command": "git status --short"},
        }
        claim_script = str((worker.REPO_ROOT / "bootstrap" / "claim.sh").resolve())
        exact_claim = {
            **ordinary,
            "tool_use_id": "exact-claim",
            "tool_input": {
                "command": "/bin/bash %s claim 84382166 2100304" % claim_script,
            },
        }
        wrong_claim = {
            **exact_claim,
            "tool_use_id": "wrong-claim",
            "tool_input": {
                "command": "/bin/bash %s claim 84399999 2100304" % claim_script,
            },
        }
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.hook("claude", ordinary),
                             worker.HOOK_BLOCK_EXIT)
            self.assertEqual(worker.hook("claude", exact_claim), 0)
            self.assertEqual(worker.hook("claude", wrong_claim),
                             worker.HOOK_BLOCK_EXIT)

    def test_headless_remote_registration_runs_after_local_lock_is_released(self):
        session_id = "headless-lock-order"
        store = self._claude_store(session_id)
        fake = FakeClient()
        lock_acquired = threading.Event()

        def register_worker(*args, **kwargs):
            def probe_lock():
                with store.locked():
                    lock_acquired.set()

            probe = threading.Thread(target=probe_lock)
            probe.start()
            self.assertTrue(lock_acquired.wait(1))
            probe.join(timeout=1)
            return fake._record("register_worker", *args, **kwargs)

        fake.register_worker = register_worker
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_pid_alive", return_value=True), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="wrapper-birth"), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"):
            result = worker.register_headless(
                session_id, 303, client_name="claude")

        self.assertTrue(lock_acquired.is_set())
        self.assertTrue(result["remoteRegistered"])
        self.assertTrue(store.path.exists())

    def test_headless_partial_recovery_or_receipt_tombstone_never_becomes_idle(self):
        fake = FakeClient()
        session_id = "headless-partial-lineage"
        store = self._claude_store(session_id)
        store.save({
            "schemaVersion": 1,
            "client": "claude",
            "clientSessionId": session_id,
            "workerKey": "host:boot:old-process",
            "claimCounter": 5,
            "recoveryPending": {
                "aoneId": "84382166",
                "projectId": "2100304",
                "runtimeSessionId": "interactive:cycle:5",
                "oldWorkerKey": "host:boot:older-process",
            },
            "lostOwnership": {
                "aoneId": "84382166",
                "projectId": "2100304",
                "reason": "old owner uncertain",
            },
            "pendingOperation": {
                "operationId": "operation-partial",
                "operationKey": "operation-key-partial",
                "aoneId": "84382166",
                "status": "SENDING",
            },
            "stopped": True,
        })
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_pid_alive", return_value=True), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="partial-birth"), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"):
            worker.register_headless(session_id, 404, client_name="claude")

        recovered = store.load()
        self.assertEqual(recovered["pendingClaim"]["phase"],
                         "READY_TO_RECOVER")
        self.assertEqual(recovered["pendingClaim"]["cycle"], 5)
        self.assertEqual(recovered["pendingClaim"]["operationKey"],
                         "operation-key-partial")
        self.assertTrue(recovered["pendingClaim"]["receiptUnknown"])
        self.assertEqual(recovered["recoveryPending"]["oldWorkerKey"],
                         "host:boot:old-process")
        self.assertEqual(recovered["lostOwnership"]["reason"],
                         "old owner uncertain")
        self.assertIsNone(recovered["pendingOperation"])

        orphan_session = "headless-orphan-receipt"
        orphan_store = self._claude_store(orphan_session)
        orphan_store.save({
            "schemaVersion": 1,
            "client": "claude",
            "clientSessionId": orphan_session,
            "workerKey": "host:boot:orphan-old",
            "claimCounter": 0,
            "pendingOperation": {
                "operationId": "operation-orphan",
                "operationKey": "operation-key-orphan",
                "aoneId": "84382166",
                "status": "SENDING",
            },
            "stopped": True,
        })
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_pid_alive", return_value=True), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="orphan-birth"), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"):
            worker.register_headless(orphan_session, 405, client_name="claude")

        orphan = orphan_store.load()
        self.assertIsNone(orphan["current"])
        self.assertIsNone(orphan["pendingClaim"])
        self.assertEqual(orphan["pendingOperation"]["operationId"],
                         "operation-orphan")
        event = {
            "hook_event_name": "PreToolUse",
            "session_id": orphan_session,
            "tool_use_id": "orphan-tool",
            "tool_name": "Bash",
            "tool_input": {"command": "git status --short"},
        }
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.hook("claude", event),
                             worker.HOOK_BLOCK_EXIT)

        suspend_session = "headless-orphan-suspend"
        suspend_store = self._claude_store(suspend_session)
        suspend_store.save({
            "schemaVersion": 1,
            "client": "claude",
            "clientSessionId": suspend_session,
            "workerKey": "host:boot:suspend-old",
            "claimCounter": 0,
            "pendingSuspend": {
                "sessionId": "session-unknown",
                "fenceToken": 19,
                "request": {"waitType": "UNKNOWN"},
                "requestId": "suspend-unknown",
            },
            "lastAutoSuspended": {
                "aoneId": "84382166",
                "projectId": "2100304",
            },
            "stopped": True,
        })
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_pid_alive", return_value=True), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="suspend-birth"), \
                mock.patch.object(worker, "_default_boot_id", return_value="boot"):
            worker.register_headless(suspend_session, 406, client_name="claude")

        suspend_tombstone = suspend_store.load()
        self.assertEqual(
            suspend_tombstone["pendingSuspend"]["requestId"],
            "suspend-unknown")
        self.assertEqual(
            suspend_tombstone["lastAutoSuspended"]["aoneId"],
            "84382166")
        self.assertIsNone(suspend_tombstone["current"])
        self.assertIsNone(suspend_tombstone["pendingClaim"])
        suspend_event = {
            "hook_event_name": "PreToolUse",
            "session_id": suspend_session,
            "tool_use_id": "suspend-orphan-tool",
            "tool_name": "Bash",
            "tool_input": {"command": "git status --short"},
        }
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True), \
                contextlib.redirect_stdout(io.StringIO()), \
                contextlib.redirect_stderr(io.StringIO()):
            self.assertEqual(worker.hook("claude", suspend_event),
                             worker.HOOK_BLOCK_EXIT)

    def test_exec_headless_registers_before_exec_with_same_pid_context(self):
        order = []

        def register(session_id, host_pid, client_name, headless_policy=None):
            order.append(("register", session_id, host_pid, client_name))
            self.assertIsNone(headless_policy)
            return {"verifyHostCommand": True}

        def execvpe(executable, command, env):
            order.append(("exec", executable, list(command),
                          env["JARVIS_INTERACTIVE_SESSION_ID"]))
            raise RuntimeError("exec intercepted")

        with mock.patch.object(worker, "register_headless",
                               side_effect=register), \
                mock.patch.object(worker.os, "execvpe", side_effect=execvpe):
            with self.assertRaisesRegex(RuntimeError, "exec intercepted"):
                worker.exec_headless(
                    "headless-session-1", ["/bin/echo", "ready"],
                    client_name="claude")

        self.assertEqual(order[0],
                         ("register", "headless-session-1", os.getpid(), "claude"))
        self.assertEqual(order[1],
                         ("exec", "/bin/echo", ["/bin/echo", "ready"],
                          "headless-session-1"))

    def test_pid_reuse_cannot_keep_an_old_worker_alive(self):
        state = self._seed()
        state["hostProcessStartedAt"] = "old-birth"
        with mock.patch.object(worker, "_pid_alive", return_value=True), \
                mock.patch.object(worker, "_process_info",
                                  return_value=(1, "/Applications/Codex.app/Contents/MacOS/Codex")), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="replacement-birth"):
            self.assertFalse(worker._host_alive(state))

    def test_next_prompt_can_restart_a_dead_sidecar(self):
        state = self._seed()
        state["daemonPid"] = 99999
        self._store().save(state)
        with mock.patch.object(worker, "_pid_alive", return_value=False), \
                mock.patch.object(worker, "_spawn_daemon", return_value=43210) as spawn:
            worker._ensure_daemon(self._store(), state["workerKey"])
        spawn.assert_called_once()
        self.assertEqual(self._store().load()["daemonPid"], 43210)

    def test_claim_retry_keeps_cycle_revision_and_resumes_local_sending_receipt(self):
        self._seed()
        fake = FakeClient()
        fake.claim_error = worker.ControlPlaneUnavailable("lost response")
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaises(worker.ControlPlaneUnavailable):
                worker.prepare_claim("84345050", "2100304", "Original Aone title")
        pending = self._store().load()["pendingClaim"]
        self.assertEqual(pending["cycle"], 1)
        self.assertEqual(pending["title"], "Original Aone title")
        first_runtime = pending["runtimeSessionId"]

        fake.claim_error = None
        with mock.patch.object(worker, "_client", return_value=fake):
            first = worker.prepare_claim("84345050", "2100304", "Renamed Aone title")
        claim_calls = [c for c in fake.calls if c[0] == "claim_task"]
        self.assertEqual(claim_calls[0][2]["request_id"],
                         claim_calls[1][2]["request_id"])
        self.assertEqual(claim_calls[0][1][1].source_ref["title"],
                         "Original Aone title")
        self.assertEqual(claim_calls[1][1][1].source_ref["title"],
                         "Original Aone title")
        first_claim = [c for c in fake.calls if c[0] == "claim_task"][-1]
        first_envelope = first_claim[1][1]
        self.assertEqual(first["runtimeSessionId"], first_runtime)
        self.assertEqual(first_envelope.recovery_policy, "REPLAY_SAFE")
        self.assertNotIn("clientSessionId", first_envelope.payload)
        begin_body = [c for c in fake.calls if c[0] == "begin_operation"][-1][1][0]
        self.assertNotIn("requestDigest", begin_body)

        fake.begin_results.append({
            "proceed": False,
            "operation": {"id": "op-1", "status": "SENDING"},
        })
        with mock.patch.object(worker, "_client", return_value=fake):
            retry = worker.prepare_claim("84345050", "2100304", "Third Aone title")
        retry_claim = [c for c in fake.calls if c[0] == "claim_task"][-1]
        retry_envelope = retry_claim[1][1]
        self.assertTrue(retry["proceed"])
        self.assertEqual(retry["runtimeSessionId"], first_runtime)
        self.assertEqual(retry_envelope.desired_revision,
                         first_envelope.desired_revision)
        self.assertEqual(retry_claim[2]["free_slots"], 1)

        with mock.patch.object(worker, "_client", return_value=fake):
            worker.acknowledge_claim("84345050", "aone:2100304:84345050:tag")
            worker.transition("84345050", "suspend", "handoff")
        suspend = [c for c in fake.calls if c[0] == "suspend_session"][-1]
        self.assertEqual(suspend[1][3]["waitType"], "MANUAL")
        self.assertEqual(suspend[1][3]["waitKey"], "aone:84345050")

        fake.begin_results.append({
            "proceed": True,
            "operation": {"id": "op-2", "status": "SENDING"},
        })
        with mock.patch.object(worker, "_client", return_value=fake):
            next_claim = worker.prepare_claim("84345050", "2100304", "Renamed Aone title")
        self.assertNotEqual(next_claim["runtimeSessionId"], first_runtime)
        self.assertIn("cycle:2", next_claim["runtimeSessionId"])

    def test_blank_title_is_omitted_without_changing_interactive_revision(self):
        self._seed()
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            claim = worker.prepare_claim("84345050", "2100304", "   ")
        envelope = [c for c in fake.calls if c[0] == "claim_task"][-1][1][1]
        self.assertEqual(envelope.source_ref,
                         {"aoneId": "84345050", "projectId": "2100304"})
        expected = "interactive:%s" % worker.hashlib.sha256(
            claim["runtimeSessionId"].encode()).hexdigest()[:32]
        self.assertEqual(envelope.desired_revision, expected)
        self.assertEqual(self._store().load()["current"]["title"], "")

    def test_different_target_cannot_overwrite_inflight_local_claim(self):
        state = self._seed()
        state["pendingClaim"] = {
            "aoneId": "84345050", "projectId": "2100304",
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "phase": "CLAIMING", "claimRequestId": "claim-inflight",
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaises(worker.ControlPlaneConflict):
                worker.prepare_claim("84399999", "2100304")
        self.assertNotIn("claim_task", [c[0] for c in fake.calls])
        self.assertEqual(self._store().load()["pendingClaim"]["aoneId"],
                         "84345050")

    def test_receipt_recovery_claim_reports_capacity_for_expired_local_session(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": False,
        }
        state["pendingOperation"] = {
            "operationId": None, "operationKey": "aone-claim:task-1:4:1",
            "aoneId": "84345050", "proceed": False, "status": "BEGINNING",
        }
        self._store().save(state)
        fake = FakeClient()
        fake.begin_results = [{
            "proceed": False,
            "operation": {"id": "op-1", "status": "SENDING"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            resumed = worker.prepare_claim("84345050", "2100304")
        claim = [c for c in fake.calls if c[0] == "claim_task"][-1]
        self.assertEqual(claim[2]["free_slots"], 1)
        self.assertTrue(resumed["proceed"])

    def test_old_ack_response_cannot_mutate_replacement_incarnation(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-old",
            "sessionId": "session-old", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": False,
        }
        state["pendingOperation"] = {
            "operationId": "op-old", "operationKey": "claim-old",
            "aoneId": "84345050", "proceed": True, "status": "SENDING",
        }
        self._store().save(state)
        fake = FakeClient()

        def ack_after_replacement(*_args, **_kwargs):
            replacement = self._store().load()
            replacement["workerKey"] = "host:boot:replacement"
            replacement["current"] = {
                "aoneId": "84399999", "projectId": "2100304",
                "taskId": "task-new", "sessionId": "session-new",
                "fenceToken": 20, "heartbeatEnabled": False,
            }
            replacement["pendingOperation"] = {
                "operationId": "op-new", "operationKey": "claim-new",
                "aoneId": "84399999", "status": "SENDING",
            }
            self._store().save(replacement)
            return {"ok": True}

        fake.ack_operation = ack_after_replacement
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaises(worker.ControlPlaneConflict):
                worker.acknowledge_claim(
                    "84345050", "aone:2100304:84345050:tag")
        replacement = self._store().load()
        self.assertEqual(replacement["pendingOperation"]["operationId"], "op-new")
        self.assertFalse(replacement["current"]["heartbeatEnabled"])

    def test_lost_begin_response_uses_durable_intent_to_resume_safely(self):
        self._seed()
        fake = FakeClient()
        fake.begin_error = worker.ControlPlaneUnavailable("lost begin response")
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaises(worker.ControlPlaneUnavailable):
                worker.prepare_claim("84345050", "2100304")

        state = self._store().load()
        self.assertEqual(state["pendingOperation"]["status"], "BEGINNING")
        self.assertIsNone(state["pendingOperation"]["operationId"])
        self.assertFalse(state["current"]["heartbeatEnabled"])

        fake.begin_error = None
        fake.begin_results = [{
            "proceed": False,
            "operation": {"id": "op-1", "status": "SENDING"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            resumed = worker.prepare_claim("84345050", "2100304")
        self.assertTrue(resumed["proceed"])
        pending = self._store().load()["pendingOperation"]
        self.assertEqual(pending["operationId"], "op-1")
        self.assertEqual(pending["status"], "SENDING")

    def test_known_operation_failure_reuses_cycle_runtime_and_operation_key(self):
        self._seed()
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            first = worker.prepare_claim("84345050", "2100304")
            first_begin = [c for c in fake.calls if c[0] == "begin_operation"][-1][1][0]
            worker.fail_claim("84345050", "known Aone rejection")
        pending = self._store().load()["pendingClaim"]
        self.assertEqual(pending["cycle"], 1)
        self.assertEqual(pending["runtimeSessionId"], first["runtimeSessionId"])
        self.assertEqual(pending["operationKey"], first_begin["operationKey"])

        fake.begin_results.append({
            "proceed": True,
            "operation": {"id": "op-1", "status": "SENDING"},
        })
        with mock.patch.object(worker, "_client", return_value=fake):
            retried = worker.prepare_claim("84345050", "2100304")
        retry_begin = [c for c in fake.calls if c[0] == "begin_operation"][-1][1][0]
        self.assertEqual(retried["runtimeSessionId"], first["runtimeSessionId"])
        self.assertEqual(retry_begin["operationKey"], first_begin["operationKey"])

    def test_claude_env_file_exports_safe_fallback_runtime_context(self):
        env_file = Path(self.temp.name) / "claude-env.sh"
        with mock.patch.dict(os.environ, {
            "CLAUDE_ENV_FILE": str(env_file),
        }, clear=False):
            worker._persist_claude_context("native-'quoted")
        content = env_file.read_text(encoding="utf-8")
        self.assertIn("JARVIS_INTERACTIVE_CLIENT=claude", content)
        self.assertIn("JARVIS_INTERACTIVE_SESSION_ID=", content)
        self.assertIn("'\"'\"'", content)
        with mock.patch.dict(os.environ, {
            "CODEX_THREAD_ID": "",
            "CLAUDE_CODE_SESSION_ID": "",
            "JARVIS_INTERACTIVE_CLIENT": "claude",
            "JARVIS_INTERACTIVE_SESSION_ID": "native-'quoted",
        }, clear=False):
            self.assertEqual(worker._runtime_context(),
                             ("claude", "native-'quoted"))

    def test_native_claude_context_wins_over_stale_persisted_context(self):
        with mock.patch.dict(os.environ, {
            "CODEX_THREAD_ID": "",
            "CLAUDE_CODE_SESSION_ID": "native-claude-current",
            "JARVIS_INTERACTIVE_CLIENT": "claude",
            "JARVIS_INTERACTIVE_SESSION_ID": "persisted-claude-old",
        }, clear=False):
            self.assertEqual(worker._runtime_context(),
                             ("claude", "native-claude-current"))

    def test_nested_idea_session_does_not_reuse_outer_codex_worker(self):
        with mock.patch.dict(os.environ, {
            "CODEX_THREAD_ID": "outer-codex-thread",
            "CLAUDE_CODE_SESSION_ID": "",
            "JARVIS_INTERACTIVE_CLIENT": "claude",
            "JARVIS_INTERACTIVE_SESSION_ID": "inner-claude-session",
        }, clear=False), \
                mock.patch.object(worker, "_nearest_runtime_client",
                                  return_value="claude"):
            self.assertEqual(worker._runtime_context(),
                             ("claude", "inner-claude-session"))

    def test_nested_codex_session_does_not_reuse_outer_claude_worker(self):
        with mock.patch.dict(os.environ, {
            "CODEX_THREAD_ID": "inner-codex-thread",
            "CLAUDE_CODE_SESSION_ID": "outer-claude-session",
            "JARVIS_INTERACTIVE_CLIENT": "claude",
            "JARVIS_INTERACTIVE_SESSION_ID": "outer-claude-session",
        }, clear=False), \
                mock.patch.object(worker, "_nearest_runtime_client",
                                  return_value="codex"):
            self.assertEqual(worker._runtime_context(),
                             ("codex", "inner-codex-thread"))

    def test_control_token_falls_back_to_existing_html_report_token(self):
        with mock.patch.dict(os.environ, {
            "JARVIS_CONTROL_PLANE_BASE_URL": "https://pre.example",
            "JARVIS_CONTROL_PLANE_TOKEN": "",
            "JARVIS_HTML_REPORT_TOKEN": "shared-report-token",
        }, clear=False):
            client = worker._client()
        self.assertEqual(client.base_url, "https://pre.example")
        self.assertEqual(client.token, "shared-report-token")

    def test_unknown_sending_receipt_fails_closed_without_session_heartbeat(self):
        self._seed()
        fake = FakeClient()
        fake.begin_results = [{
            "proceed": False,
            "operation": {"id": "op-1", "status": "SENDING"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaises(worker.ControlPlaneConflict):
                worker.prepare_claim("84345050", "2100304")
        current = self._store().load()["current"]
        self.assertFalse(current["heartbeatEnabled"])

    def test_sidecar_heartbeats_active_session_then_marks_worker_offline(self):
        state = self._seed()
        state["current"] = {
            "sessionId": "session-1", "fenceToken": 9,
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", side_effect=[True, False]), \
                mock.patch.object(worker.time, "sleep", return_value=None):
            self.assertEqual(worker.daemon(self._store().path, state["workerKey"]), 0)
        names = [c[0] for c in fake.calls]
        self.assertIn("heartbeat_session", names)
        refreshed = self._store().load()["sessionPermit"]
        self.assertEqual(refreshed["source"], "sidecar-heartbeat")
        self.assertEqual(refreshed["sessionStatus"], "RUNNING")
        offline = [c for c in fake.calls if c[0] == "heartbeat_worker"][-1]
        self.assertEqual(offline[1][1]["status"], "OFFLINE")

    def test_sidecar_keeps_worker_alive_but_stops_expired_turn_session_heartbeat(self):
        state = self._seed()
        state["turnActive"] = False
        state["turnStoppedAt"] = 100
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", side_effect=[True, False]), \
                mock.patch.object(worker.time, "time", return_value=1000), \
                mock.patch.object(worker.time, "sleep", return_value=None):
            self.assertEqual(worker.daemon(self._store().path, state["workerKey"]), 0)
        names = [c[0] for c in fake.calls]
        self.assertIn("heartbeat_worker", names)
        self.assertNotIn("heartbeat_session", names)
        self.assertIn("suspend_session", names)
        after = self._store().load()
        self.assertIsNone(after["current"])
        self.assertEqual(after["pendingClaim"]["cycle"], 1)
        self.assertEqual(after["pendingClaim"]["runtimeSessionId"],
                         "interactive:cycle:1")
        suspend = [c for c in fake.calls if c[0] == "suspend_session"][0]
        self.assertEqual(
            suspend[1][3]["waitExpireAt"],
            worker._utc_timestamp(1000 + 7200))
        self.assertEqual(
            after["lastAutoSuspended"]["affinityWorkerKey"],
            state["workerKey"])
        self.assertEqual(
            after["lastAutoSuspended"]["affinityExpireAt"], 8200)
        offline = [c for c in fake.calls if c[0] == "heartbeat_worker"][-1]
        self.assertEqual(offline[1][1]["status"], "OFFLINE")

    def test_lost_suspend_response_is_replayed_before_next_turn(self):
        state = self._seed()
        state["turnActive"] = False
        state["turnStoppedAt"] = 100
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 3, "runtimeSessionId": "interactive:cycle:3",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        fake = FakeClient()
        fake.suspend_error = worker.ControlPlaneUnavailable("lost response")
        with mock.patch.object(worker.time, "time", return_value=1000):
            with self.assertRaises(worker.ControlPlaneUnavailable):
                worker._auto_suspend_idle_session(
                    self._store(), fake, state["workerKey"])
        uncertain = self._store().load()
        self.assertIsNotNone(uncertain["current"])
        self.assertIsNotNone(uncertain["pendingSuspend"])
        first_request = [c for c in fake.calls if c[0] == "suspend_session"][-1]

        fake.suspend_error = None
        fake.begin_results = [{
            "proceed": False,
            "operation": {"id": "op-1", "status": "ACKED"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "UserPromptSubmit",
                "session_id": "native-thread-1",
                "turn_id": "turn-next",
            }), 0)
        second_request = [c for c in fake.calls if c[0] == "suspend_session"][-1]
        self.assertEqual(first_request[1], second_request[1])
        self.assertEqual(first_request[2], second_request[2])
        recovered = self._store().load()
        self.assertIsNotNone(recovered["current"])
        self.assertIsNone(recovered["pendingSuspend"])
        self.assertIsNone(recovered["pendingClaim"])
        self.assertEqual(recovered["current"]["cycle"], 3)
        self.assertEqual(recovered["current"]["runtimeSessionId"],
                         "interactive:cycle:3")
        self.assertTrue(recovered["current"]["heartbeatEnabled"])

    def test_pending_external_receipt_is_not_auto_suspended(self):
        state = self._seed()
        state["turnActive"] = False
        state["turnStoppedAt"] = 100
        state["pendingOperation"] = {"operationId": "op-1"}
        state["current"] = {
            "sessionId": "session-1", "fenceToken": 9,
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker.time, "time", return_value=1000):
            self.assertFalse(worker._auto_suspend_idle_session(
                self._store(), fake, state["workerKey"]))
        self.assertNotIn("suspend_session", [c[0] for c in fake.calls])
        self.assertIsNotNone(self._store().load()["current"])

    def test_standard_claim_resumes_auto_suspended_cycle_and_runtime(self):
        state = self._seed()
        state["claimCounter"] = 5
        state["turnActive"] = False
        state["turnStoppedAt"] = 100
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 5, "runtimeSessionId": "interactive:cycle:5",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker.time, "time", return_value=1000):
            self.assertTrue(worker._auto_suspend_idle_session(
                self._store(), fake, state["workerKey"]))
        with mock.patch.object(worker, "_client", return_value=fake):
            resumed = worker.prepare_claim("84345050", "2100304")
        self.assertEqual(resumed["runtimeSessionId"], "interactive:cycle:5")
        claim = [c for c in fake.calls if c[0] == "claim_task"][-1]
        self.assertEqual(claim[2]["runtime_session_id"], "interactive:cycle:5")
        self.assertEqual(self._store().load()["claimCounter"], 5)

    def test_next_prompt_auto_resumes_same_fenced_session_with_fresh_ids(self):
        self._seed()
        fake = FakeClient()
        fake.claim_results = [
            {
                "task": {"id": "task-1", "generation": 4},
                "session": {"id": "session-1", "generation": 4,
                            "fenceToken": 9, "attemptNo": 1},
            },
            {
                "task": {"id": "task-1", "generation": 4},
                "session": {"id": "session-1", "generation": 4,
                            "fenceToken": 10, "attemptNo": 2},
            },
        ]
        with mock.patch.object(worker, "_client", return_value=fake):
            first = worker.prepare_claim("84345050", "2100304")
            worker.acknowledge_claim(
                "84345050", "aone:2100304:84345050:tag")
        state = self._store().load()
        state["lastTurnActivityAt"] = 100
        state["turnActive"] = False
        state["turnStoppedAt"] = 100
        self._store().save(state)
        with mock.patch.object(worker.time, "time", return_value=1000):
            self.assertTrue(worker._auto_suspend_idle_session(
                self._store(), fake, state["workerKey"]))

        fake.begin_results.append({
            "proceed": False,
            "operation": {"id": "op-1", "status": "ACKED"},
        })
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "UserPromptSubmit",
                "session_id": "native-thread-1",
                "turn_id": "turn-continue",
                "prompt": "继续处理刚才的单",
            }), 0)

        claims = [c for c in fake.calls if c[0] == "claim_task"]
        starts = [c for c in fake.calls if c[0] == "start_session"]
        self.assertEqual(len(claims), 2)
        self.assertNotEqual(claims[0][2]["request_id"],
                            claims[1][2]["request_id"])
        self.assertNotEqual(starts[0][2]["request_id"],
                            starts[1][2]["request_id"])
        current = self._store().load()["current"]
        self.assertEqual(current["runtimeSessionId"], first["runtimeSessionId"])
        self.assertEqual(current["cycle"], 1)
        self.assertEqual(current["fenceToken"], 10)
        self.assertTrue(current["heartbeatEnabled"])

    def test_prompt_for_different_aone_does_not_reacquire_suspended_task(self):
        state = self._seed()
        state["turnActive"] = False
        state["turnStoppedAt"] = 100
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker.time, "time", return_value=1000):
            worker._auto_suspend_idle_session(self._store(), fake, state["workerKey"])
        output = io.StringIO()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(output):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "UserPromptSubmit",
                "session_id": "native-thread-1",
                "turn_id": "turn-other",
                "prompt": "处理 84399999",
            }), 0)
        self.assertNotIn("claim_task", [c[0] for c in fake.calls])
        self.assertIsNone(self._store().load()["current"])
        self.assertIn("未自动抢回旧单", output.getvalue())

    def test_codex_restart_preserves_lineage_but_never_old_fence_or_receipt(self):
        state = self._seed()
        state["current"] = {
            "aoneId": "84345050", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-old", "fenceToken": 9, "generation": 4,
            "cycle": 7, "runtimeSessionId": "interactive:cycle:7",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        state["pendingOperation"] = {"operationId": "op-old"}
        self._store().save(state)
        fake = FakeClient()
        event = {
            "hook_event_name": "SessionStart",
            "session_id": "native-thread-1",
            "cwd": self.temp.name,
            "source": "resume",
        }
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid", return_value=(os.getpid(), True)), \
                mock.patch.object(worker, "_default_boot_id", return_value="new-boot"), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", event), 0)
        recovered = self._store().load()
        self.assertIsNone(recovered["current"])
        self.assertIsNone(recovered["pendingOperation"])
        self.assertEqual(recovered["pendingClaim"]["phase"], "READY_TO_RECOVER")
        self.assertEqual(recovered["pendingClaim"]["runtimeSessionId"],
                         "interactive:cycle:7")
        self.assertNotEqual(recovered["workerKey"], state["workerKey"])

        fake.claim_error = worker.ControlPlaneConflict("old lease still active")
        output = io.StringIO()
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_host_alive", return_value=True), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(output):
            self.assertEqual(worker.hook("codex", {
                "hook_event_name": "UserPromptSubmit",
                "session_id": "native-thread-1",
                "turn_id": "turn-after-restart",
                "prompt": "继续",
            }), 0)
        blocked = self._store().load()
        self.assertIsNone(blocked["current"])
        self.assertEqual(blocked["pendingClaim"]["phase"], "READY_TO_RECOVER")
        self.assertNotIn("claimRequestId", blocked["pendingClaim"])
        self.assertIn("不得继续该单", output.getvalue())

    # --- mid-task external-write operation receipts (operation-begin/abort/reconcile) ---

    def _seed_with_current(self, aone_id="84386065"):
        state = self._seed()
        state["current"] = {
            "aoneId": aone_id, "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-1", "fenceToken": 9, "generation": 4,
            "cycle": 1, "runtimeSessionId": "interactive:cycle:1",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        return state

    def _operation_key(self, kind, material, session_id="session-1", fence=9):
        import hashlib
        attempt = hashlib.sha256(
            ("%s:%s" % (session_id, fence)).encode("utf-8")).hexdigest()[:12]
        return "%s:task-1:4:%s:%s" % (
            kind, attempt,
            hashlib.sha256(material.encode("utf-8")).hexdigest()[:12])

    def test_operation_begin_proceeds_and_keeps_heartbeat_untouched(self):
        self._seed_with_current()
        fake = FakeClient()
        fake.begin_results = [{
            "proceed": True,
            "operation": {"id": "op-9", "status": "SENDING"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_begin("84386065", "comment", "digest-abc")
        self.assertTrue(result["proceed"])
        self.assertFalse(result["needsReadback"])
        begin_body = [c for c in fake.calls if c[0] == "begin_operation"][-1][1][0]
        self.assertEqual(begin_body["operationType"], "AONE_COMMENT")
        self.assertEqual(begin_body["operationKey"],
                         self._operation_key("comment", "digest-abc"))
        self.assertTrue(begin_body["required"])
        state = self._store().load()
        pending = state["pendingOperation"]
        self.assertEqual(pending["operationId"], "op-9")
        self.assertEqual(pending["operationKey"],
                         self._operation_key("comment", "digest-abc"))
        self.assertTrue(pending["proceed"])
        # mid-task 回执不动 heartbeatEnabled（区别于 claim 期间的续租闸门）。
        self.assertTrue(state["current"]["heartbeatEnabled"])

    def test_operation_begin_acked_skips_send_and_clears_slot(self):
        self._seed_with_current()
        fake = FakeClient()
        fake.begin_results = [{
            "proceed": False,
            "operation": {"id": "op-9", "status": "ACKED"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_begin("84386065", "comment", "digest-abc")
        self.assertFalse(result["proceed"])
        self.assertFalse(result["needsReadback"])
        self.assertEqual(result["operationStatus"], "ACKED")
        self.assertIsNone(self._store().load()["pendingOperation"])

    def test_operation_begin_sending_replay_safe_replays_local_intent(self):
        state = self._seed_with_current()
        key = self._operation_key("release-tag", "idle")
        state["pendingOperation"] = {
            "operationId": "op-9", "operationKey": key,
            "aoneId": "84386065", "kind": "release-tag",
            "proceed": False, "status": "SENDING",
        }
        self._store().save(state)
        fake = FakeClient()
        fake.begin_results = [{
            "proceed": False,
            "operation": {"id": "op-9", "status": "SENDING"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_begin(
                "84386065", "release-tag", "idle", replay_safe=True)
        self.assertTrue(result["proceed"])
        self.assertFalse(result["needsReadback"])
        begin_body = [c for c in fake.calls if c[0] == "begin_operation"][-1][1][0]
        self.assertEqual(begin_body["operationType"], "AONE_RELEASE")

    def test_operation_begin_sending_comment_requires_readback(self):
        state = self._seed_with_current()
        key = self._operation_key("comment", "digest-abc")
        state["pendingOperation"] = {
            "operationId": None, "operationKey": key,
            "aoneId": "84386065", "kind": "comment",
            "proceed": False, "status": "BEGINNING",
        }
        self._store().save(state)
        fake = FakeClient()
        fake.begin_results = [{
            "proceed": False,
            "operation": {"id": "op-9", "status": "SENDING"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_begin("84386065", "comment", "digest-abc")
        self.assertFalse(result["proceed"])
        self.assertTrue(result["needsReadback"])
        pending = self._store().load()["pendingOperation"]
        self.assertEqual(pending["operationId"], "op-9")
        self.assertEqual(pending["status"], "SENDING")

    def test_operation_begin_sending_without_local_intent_fails_closed(self):
        self._seed_with_current()
        fake = FakeClient()
        fake.begin_results = [{
            "proceed": False,
            "operation": {"id": "op-9", "status": "SENDING"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaises(worker.ControlPlaneConflict):
                worker.operation_begin("84386065", "comment", "digest-abc")

    def test_operation_begin_conflicts_with_different_pending_key(self):
        state = self._seed_with_current()
        state["pendingOperation"] = {
            "operationId": "op-other",
            "operationKey": self._operation_key("comment", "other-material"),
            "aoneId": "84386065", "kind": "comment",
            "proceed": True, "status": "SENDING",
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaises(worker.ControlPlaneConflict):
                worker.operation_begin("84386065", "comment", "digest-abc")
        self.assertNotIn("begin_operation", [c[0] for c in fake.calls])

    def test_operation_begin_unknown_slot_short_circuits_to_readback(self):
        state = self._seed_with_current()
        key = self._operation_key("comment", "digest-abc")
        state["pendingOperation"] = {
            "operationId": "op-7", "operationKey": key,
            "aoneId": "84386065", "kind": "comment",
            "proceed": False, "status": "UNKNOWN",
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_begin("84386065", "comment", "digest-abc")
        self.assertFalse(result["proceed"])
        self.assertTrue(result["needsReadback"])
        self.assertEqual(result["operationId"], "op-7")
        self.assertEqual(result["operationStatus"], "UNKNOWN")
        self.assertNotIn("begin_operation", [c[0] for c in fake.calls])
        self.assertEqual(
            self._store().load()["pendingOperation"]["status"], "UNKNOWN")

    def test_operation_abort_clears_slot_but_keeps_session(self):
        state = self._seed_with_current()
        state["pendingOperation"] = {
            "operationId": "op-9",
            "operationKey": self._operation_key("comment", "digest-abc"),
            "aoneId": "84386065", "kind": "comment",
            "proceed": True, "status": "SENDING",
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_abort("84386065", "a1 comment failed")
        self.assertFalse(result["unknown"])
        fail_body = [c for c in fake.calls if c[0] == "fail_operation"][-1][1][0]
        self.assertFalse(fail_body["unknown"])
        self.assertTrue(fail_body["retryAllowed"])
        self.assertNotIn("fail_session", [c[0] for c in fake.calls])
        state = self._store().load()
        self.assertIsNone(state["pendingOperation"])
        self.assertIsNone(state["pendingClaim"])
        self.assertEqual(state["current"]["aoneId"], "84386065")

    def test_operation_abort_unknown_freezes_local_slot(self):
        state = self._seed_with_current()
        key = self._operation_key("comment", "digest-abc")
        state["pendingOperation"] = {
            "operationId": "op-9", "operationKey": key,
            "aoneId": "84386065", "kind": "comment",
            "proceed": True, "status": "SENDING",
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_abort(
                "84386065", "comment result indeterminate", unknown=True)
        self.assertTrue(result["unknown"])
        fail_body = [c for c in fake.calls if c[0] == "fail_operation"][-1][1][0]
        self.assertTrue(fail_body["unknown"])
        self.assertNotIn("fail_session", [c[0] for c in fake.calls])
        pending = self._store().load()["pendingOperation"]
        self.assertEqual(pending["status"], "UNKNOWN")
        self.assertEqual(pending["operationKey"], key)
        self.assertIsNotNone(self._store().load()["current"])

    def test_operation_reconcile_found_clears_slot(self):
        state = self._seed_with_current()
        state["pendingOperation"] = {
            "operationId": "op-9",
            "operationKey": self._operation_key("comment", "digest-abc"),
            "aoneId": "84386065", "kind": "comment",
            "proceed": False, "status": "UNKNOWN",
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_reconcile(
                "84386065", found=True,
                external_ref="aone:84386065:comment:555")
        self.assertFalse(result["proceed"])
        reconcile_body = [
            c for c in fake.calls if c[0] == "reconcile_operation"][-1][1][0]
        self.assertTrue(reconcile_body["found"])
        self.assertEqual(reconcile_body["externalRef"],
                         "aone:84386065:comment:555")
        self.assertIsNone(self._store().load()["pendingOperation"])

    def test_operation_reconcile_not_found_keeps_key_for_retry(self):
        state = self._seed_with_current()
        key = self._operation_key("comment", "digest-abc")
        state["pendingOperation"] = {
            "operationId": "op-9", "operationKey": key,
            "aoneId": "84386065", "kind": "comment",
            "proceed": False, "status": "UNKNOWN",
        }
        self._store().save(state)
        fake = FakeClient()
        with mock.patch.object(worker, "_client", return_value=fake):
            result = worker.operation_reconcile("84386065", found=False)
        self.assertFalse(result["proceed"])
        self.assertTrue(result["retryScheduled"])
        reconcile_body = [
            c for c in fake.calls if c[0] == "reconcile_operation"][-1][1][0]
        self.assertFalse(reconcile_body["found"])
        self.assertTrue(reconcile_body["retryAllowed"])
        pending = self._store().load()["pendingOperation"]
        self.assertEqual(pending["operationKey"], key)
        self.assertEqual(pending["status"], "RETRY_WAIT")

    def test_release_receipt_is_isolated_per_lease_attempt(self):
        # 同 generation、不同 Session/fence 的两轮 claim→release：第二轮绝不能
        # 命中第一轮的 ACKED 回执而跳过标签写（否则任务同代复活后 Aone 卡在
        # jarvis-claimed，控制面却已 SUSPENDED）。
        fake = FakeClient()
        key_by_op = {}
        acked_keys = set()

        def begin_operation(request, request_id=None):
            fake.calls.append(
                ("begin_operation", (request,), {"request_id": request_id}))
            key = request["operationKey"]
            if key in acked_keys:
                op_id = next(op for op, k in key_by_op.items() if k == key)
                return {"proceed": False,
                        "operation": {"id": op_id, "status": "ACKED"}}
            op_id = "op-%d" % (len(key_by_op) + 1)
            key_by_op[op_id] = key
            return {"proceed": True,
                    "operation": {"id": op_id, "status": "SENDING"}}

        def ack_operation(request, request_id=None):
            fake.calls.append(
                ("ack_operation", (request,), {"request_id": request_id}))
            acked_keys.add(key_by_op[str(request["operationId"])])
            return {"ok": True}

        fake.begin_operation = begin_operation
        fake.ack_operation = ack_operation

        # attempt A：release-tag 写成功并 ACK。
        self._seed_with_current()
        with mock.patch.object(worker, "_client", return_value=fake):
            first = worker.operation_begin(
                "84386065", "release-tag", "idle", replay_safe=True)
            self.assertTrue(first["proceed"])
            worker.acknowledge_claim(
                "84386065", "aone:2100304:84386065:tag:jarvis-idle")

        # 任务同代复活：attempt B = 新 Session/fence，同 task/generation 再 release。
        state = self._store().load()
        state["current"] = {
            "aoneId": "84386065", "projectId": "2100304", "taskId": "task-1",
            "sessionId": "session-2", "fenceToken": 20, "generation": 4,
            "cycle": 2, "runtimeSessionId": "interactive:cycle:2",
            "leaseSeconds": 120, "heartbeatEnabled": True,
        }
        self._store().save(state)
        with mock.patch.object(worker, "_client", return_value=fake):
            second = worker.operation_begin(
                "84386065", "release-tag", "idle", replay_safe=True)
        # attempt 分量隔离：第二轮拿到 proceed=True 真执行标签写。
        self.assertTrue(second["proceed"])
        keys = [c[1][0]["operationKey"] for c in fake.calls
                if c[0] == "begin_operation"]
        self.assertEqual(len(keys), 2)
        self.assertNotEqual(keys[0], keys[1])
        self.assertEqual(keys[0], self._operation_key(
            "release-tag", "idle", session_id="session-1", fence=9))
        self.assertEqual(keys[1], self._operation_key(
            "release-tag", "idle", session_id="session-2", fence=20))

    def test_restart_orphans_midtask_receipt_and_new_claim_gets_fresh_key(self):
        # UNKNOWN 的 mid-task 回执 + 进程重启（replacement 路径）：回执坐标转入
        # orphanOperations、pendingClaim 不继承其 operationKey；服务端收敛后新
        # prepare_claim 用全新 aone-claim key 成功接单，而不是复用 mid-task key
        # 被服务端按「operationKey reused with a different request」拒绝。
        state = self._seed_with_current()
        midtask_key = self._operation_key("comment", "digest-abc")
        state["pendingOperation"] = {
            "operationId": "op-mid", "operationKey": midtask_key,
            "aoneId": "84386065", "kind": "comment",
            "proceed": False, "status": "UNKNOWN",
        }
        self._store().save(state)
        fake = FakeClient()
        event = {
            "hook_event_name": "SessionStart",
            "session_id": "native-thread-1",
            "cwd": self.temp.name,
            "source": "resume",
        }
        with mock.patch.object(worker, "_client", return_value=fake), \
                mock.patch.object(worker, "_find_host_pid",
                                  return_value=(os.getpid(), True)), \
                mock.patch.object(worker, "_default_boot_id",
                                  return_value="new-boot"), \
                mock.patch.object(worker, "_ensure_daemon"), \
                contextlib.redirect_stdout(io.StringIO()):
            self.assertEqual(worker.hook("codex", event), 0)

        recovered = self._store().load()
        self.assertIsNone(recovered["pendingOperation"])
        orphan = recovered["orphanOperations"][-1]
        self.assertEqual(orphan["operationId"], "op-mid")
        self.assertEqual(orphan["operationKey"], midtask_key)
        self.assertEqual(orphan["kind"], "comment")
        self.assertEqual(orphan["aoneId"], "84386065")
        self.assertEqual(orphan["status"], "UNKNOWN")
        pending_claim = recovered["pendingClaim"]
        self.assertEqual(pending_claim["phase"], "READY_TO_RECOVER")
        self.assertNotIn("operationKey", pending_claim)
        self.assertNotIn("receiptUnknown", pending_claim)

        # 服务端 op-mid 已被外部 reconcile 收敛为 ACKED（本地无需感知）；新 claim
        # 的 begin 携全新 key，服务端按新操作受理。
        fake.begin_results = [{
            "proceed": True,
            "operation": {"id": "op-claim-new", "status": "SENDING"},
        }]
        with mock.patch.object(worker, "_client", return_value=fake):
            claimed = worker.prepare_claim("84386065", "2100304")
        self.assertTrue(claimed["accepted"])
        self.assertTrue(claimed["proceed"])
        begin_body = [c for c in fake.calls
                      if c[0] == "begin_operation"][-1][1][0]
        self.assertEqual(begin_body["operationType"], "AONE_CLAIM")
        self.assertTrue(begin_body["operationKey"].startswith("aone-claim:"))
        self.assertNotEqual(begin_body["operationKey"], midtask_key)
        # orphan 记录跨 claim 存续，留给生产态 Reconciler/人工按 operationId 收敛。
        self.assertEqual(
            self._store().load()["orphanOperations"][-1]["operationId"],
            "op-mid")

    # --- PreToolUse recovery whitelist for frozen external-write receipts ---

    def _receipt_recovery_event(self, command, index=0):
        return {
            "hook_event_name": "PreToolUse",
            "session_id": "native-thread-1",
            "turn_id": "turn-recover",
            "tool_use_id": "tool-receipt-%d" % index,
            "tool_name": "Bash",
            "tool_input": {"command": command},
        }

    def _seed_frozen_receipt(self, status="UNKNOWN", slot_aone="84386065"):
        state = self._seed_with_current()
        state["activeTurnId"] = "turn-recover"
        state["pendingOperation"] = {
            "operationId": "op-9",
            "operationKey": self._operation_key("comment", "digest-abc"),
            "aoneId": slot_aone, "kind": "comment",
            "proceed": False, "status": status,
        }
        self._store().save(state)
        return state

    def test_frozen_external_receipt_allows_exact_recovery_commands(self):
        self._seed_frozen_receipt(status="UNKNOWN")
        wrap_script = worker.REPO_ROOT / "bootstrap" / "wrap.sh"
        claim_script = worker.REPO_ROOT / "bootstrap" / "claim.sh"
        allowed = (
            "/bin/bash %s sync 84386065 --summary-file /tmp/x.md" % wrap_script,
            "/bin/bash %s done 84386065 收敛完成 已完成" % wrap_script,
            "/bin/bash %s done-no-status 84386065 收敛完成" % wrap_script,
            "/bin/bash %s release 84386065 2100304" % claim_script,
            "/bin/bash %s finish 84386065 2100304" % claim_script,
        )
        for index, command in enumerate(allowed):
            with self.subTest(command=command), \
                    mock.patch.object(worker, "_client",
                                      side_effect=AssertionError("recovery gate is local")), \
                    mock.patch.object(worker, "_calling_process_matches",
                                      return_value=True):
                self.assertIsNone(worker._guard_pre_tool_use(
                    self._store(), "codex",
                    self._receipt_recovery_event(command, index)))
        # RETRY_WAIT 槽（reconcile not-found 后进程丢失）同样允许收敛重跑。
        state = self._store().load()
        state["pendingOperation"]["status"] = "RETRY_WAIT"
        self._store().save(state)
        with mock.patch.object(worker, "_client",
                               side_effect=AssertionError("recovery gate is local")), \
                mock.patch.object(worker, "_calling_process_matches",
                                  return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._store(), "codex", self._receipt_recovery_event(
                    "/bin/bash %s sync 84386065 progress" % wrap_script, 99)))

    def test_frozen_external_receipt_blocks_non_exact_recovery_shapes(self):
        self._seed_frozen_receipt(status="UNKNOWN")
        wrap_script = worker.REPO_ROOT / "bootstrap" / "wrap.sh"
        claim_script = worker.REPO_ROOT / "bootstrap" / "claim.sh"
        rejected = (
            "/bin/bash %s sync 84386065 x; a1 update" % wrap_script,
            "/bin/bash %s sync 84386065 x && echo ok" % wrap_script,
            "/bin/bash %s sync 84386065 x | tee /tmp/x" % wrap_script,
            "/bin/bash %s sync 84399999 --summary-file /tmp/x.md" % wrap_script,
            "bash %s sync 84386065 x" % wrap_script,
            "/bin/bash bootstrap/wrap.sh sync 84386065 x",
            "/bin/bash /tmp/evil/bootstrap/wrap.sh sync 84386065 x",
            "env FOO=1 /bin/bash %s sync 84386065 x" % wrap_script,
            "JARVIS_A1=/tmp/fake /bin/bash %s sync 84386065 x" % wrap_script,
            "/bin/bash %s update 84386065 x" % wrap_script,
            "/bin/bash %s finish 84386065 2100304 已完成" % claim_script,
            "/bin/bash %s release 84399999 2100304" % claim_script,
        )
        for index, command in enumerate(rejected):
            with self.subTest(command=command), \
                    mock.patch.object(worker, "_calling_process_matches",
                                      return_value=True):
                reason = worker._guard_pre_tool_use(
                    self._store(), "codex",
                    self._receipt_recovery_event(command, index))
            self.assertIsNotNone(reason)
            self.assertIn("已阻断", reason)
        # 普通工具仍被阻断，且提示的是收敛命令而非 claim.sh claim。
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True):
            reason = worker._guard_pre_tool_use(
                self._store(), "codex",
                self._receipt_recovery_event("git status --short", 90))
        self.assertIn("外部写回执", reason)
        self.assertIn("wrap.sh", reason)

    def test_frozen_receipt_for_other_aone_never_unlocks_recovery(self):
        # 槽的 aoneId 与 current assignment 不一致（异常/跨单残留）→ 不开白名单。
        self._seed_frozen_receipt(status="UNKNOWN", slot_aone="84399999")
        wrap_script = worker.REPO_ROOT / "bootstrap" / "wrap.sh"
        for index, aone in enumerate(("84386065", "84399999")):
            with self.subTest(aone=aone), \
                    mock.patch.object(worker, "_calling_process_matches",
                                      return_value=True):
                reason = worker._guard_pre_tool_use(
                    self._store(), "codex", self._receipt_recovery_event(
                        "/bin/bash %s sync %s progress" % (wrap_script, aone),
                        index))
            self.assertIsNotNone(reason)

    def test_claim_receipt_slot_keeps_claim_only_escape_hatch(self):
        # AONE_CLAIM 槽（无 kind 字段）不受外部写白名单影响：仍只放行标准 claim.sh。
        state = self._seed_with_current(aone_id="84345050")
        state["activeTurnId"] = "turn-recover"
        state["pendingOperation"] = {
            "operationId": "op-1", "operationKey": "aone-claim:task-1:4:1",
            "aoneId": "84345050", "proceed": False, "status": "SENDING",
        }
        self._store().save(state)
        wrap_script = worker.REPO_ROOT / "bootstrap" / "wrap.sh"
        claim_script = worker.REPO_ROOT / "bootstrap" / "claim.sh"
        with mock.patch.object(worker, "_calling_process_matches",
                               return_value=True):
            reason = worker._guard_pre_tool_use(
                self._store(), "codex", self._receipt_recovery_event(
                    "/bin/bash %s sync 84345050 progress" % wrap_script, 0))
        self.assertIsNotNone(reason)
        self.assertIn("claim.sh", reason)
        with mock.patch.object(worker, "_client",
                               side_effect=AssertionError("claim gate is local")), \
                mock.patch.object(worker, "_calling_process_matches",
                                  return_value=True):
            self.assertIsNone(worker._guard_pre_tool_use(
                self._store(), "codex", self._receipt_recovery_event(
                    "/bin/bash %s claim 84345050 2100304" % claim_script, 1)))

    def test_idle_ttl_rechecks_active_and_inflight_state_under_lock(self):
        state = self._seed()
        state["lastTurnActivityAt"] = 100
        fake = FakeClient()
        with mock.patch.dict(os.environ, {
            "JARVIS_INTERACTIVE_IDLE_TTL_SEC": "60",
        }, clear=False), mock.patch.object(worker.time, "time", return_value=200):
            self.assertFalse(worker._worker_idle_expired(state))
            self.assertFalse(worker.offline(
                self._store(), fake, state["workerKey"], "idle_ttl"))
        self.assertFalse(self._store().load()["stopped"])
        self.assertNotIn("heartbeat_worker", [c[0] for c in fake.calls])

        state = self._store().load()
        state["lastTurnActivityAt"] = 100
        state["turnActive"] = False
        state["pendingClaim"] = {"phase": "CLAIMING"}
        self._store().save(state)
        with mock.patch.dict(os.environ, {
            "JARVIS_INTERACTIVE_IDLE_TTL_SEC": "60",
        }, clear=False):
            self.assertFalse(worker._worker_idle_expired(state, now=200))
            state["pendingClaim"]["phase"] = "READY_TO_RESUME"
            self.assertTrue(worker._worker_idle_expired(state, now=200))


class StopCheckTest(unittest.TestCase):
    """Tests for the stop-check subcommand (control-plane session exit gate)."""

    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.env = mock.patch.dict(os.environ, {
            "HOME": self.temp.name,
            "JARVIS_INTERACTIVE_STATE_DIR": self.temp.name,
            "JARVIS_INTERACTIVE_RETRY_DELAY": "0",
            "CODEX_THREAD_ID": "native-thread-1",
            "CLAUDE_CODE_SESSION_ID": "",
        }, clear=False)
        self.env.start()
        self._nearest_patch = mock.patch.object(
            worker, "_nearest_runtime_client", return_value="")
        self._nearest_patch.start()

    def tearDown(self):
        self._nearest_patch.stop()
        self.env.stop()
        self.temp.cleanup()

    def _store(self):
        return worker.StateStore(worker._state_path("codex", "native-thread-1"))

    def _seed(self, **overrides):
        state = {
            "schemaVersion": 1,
            "client": "codex",
            "clientSessionId": "native-thread-1",
            "workerKey": "host:boot:process",
            "host": "host",
            "bootId": "boot",
            "processUuid": "process",
            "hostPid": os.getpid(),
            "hostProcessStartedAt": worker._process_start_identity(os.getpid()),
            "verifyHostCommand": True,
            "cwd": "/workspace",
            "transcriptPath": "/tmp/transcript.jsonl",
            "branch": "worktree-ticket",
            "version": "test",
            "claimCounter": 0,
            "current": None,
            "pendingClaim": None,
            "pendingOperation": None,
            "pendingSuspend": None,
            "stopped": False,
            "turnActive": True,
            "activeTurnId": None,
            "daemonPid": os.getpid(),
            "daemonStartedAt": worker.time.time(),
            "sidecarHeartbeatAt": worker.time.time(),
        }
        state.update(overrides)
        self._store().save(state)
        return state

    def test_stop_check_no_current_passes(self):
        """No active task → exit 0 (safe to stop)."""
        self._seed(current=None)
        self.assertEqual(worker.stop_check(), 0)

    def test_stop_check_empty_state_passes(self):
        """Empty state dict → exit 0."""
        self._store().save({})
        self.assertEqual(worker.stop_check(), 0)

    def test_stop_check_terminal_session_completed_passes(self):
        """Session COMPLETED → exit 0."""
        self._seed(
            current={"aoneId": "123", "sessionId": "s-1"},
            sessionPermit={"sessionStatus": "COMPLETED", "issuedAt": 1},
        )
        self.assertEqual(worker.stop_check(), 0)

    def test_stop_check_terminal_session_failed_passes(self):
        """Session FAILED → exit 0."""
        self._seed(
            current={"aoneId": "123", "sessionId": "s-1"},
            sessionPermit={"sessionStatus": "FAILED", "issuedAt": 1},
        )
        self.assertEqual(worker.stop_check(), 0)

    def test_stop_check_terminal_session_suspended_passes(self):
        """Session SUSPENDED → exit 0."""
        self._seed(
            current={"aoneId": "123", "sessionId": "s-1"},
            sessionPermit={"sessionStatus": "SUSPENDED", "issuedAt": 1},
        )
        self.assertEqual(worker.stop_check(), 0)

    def test_stop_check_lost_ownership_passes(self):
        """Ownership lost → exit 0 (not our problem anymore)."""
        self._seed(
            current={"aoneId": "123", "sessionId": "s-1"},
            sessionPermit={"sessionStatus": "RUNNING", "issuedAt": 1},
            lostOwnership={"aoneId": "123", "reason": "stale_fence"},
        )
        self.assertEqual(worker.stop_check(), 0)

    def test_stop_check_stopped_passes(self):
        """Worker already stopped → exit 0."""
        self._seed(
            current={"aoneId": "123", "sessionId": "s-1"},
            sessionPermit={"sessionStatus": "RUNNING", "issuedAt": 1},
            stopped=True,
        )
        self.assertEqual(worker.stop_check(), 0)

    def test_stop_check_active_session_blocks(self):
        """Active session with RUNNING status → exit 2 (blocked)."""
        self._seed(
            current={"aoneId": "456", "sessionId": "s-2"},
            sessionPermit={"sessionStatus": "RUNNING", "issuedAt": 1},
        )
        self.assertEqual(worker.stop_check(), 2)

    def test_stop_check_active_session_no_permit_blocks(self):
        """Active session with no permit → exit 2 (status unknown, block)."""
        self._seed(
            current={"aoneId": "789", "sessionId": "s-3"},
        )
        self.assertEqual(worker.stop_check(), 2)

    def test_stop_check_pending_operation_blocks(self):
        """Active session with pending operation → exit 2."""
        self._seed(
            current={"aoneId": "111", "sessionId": "s-4"},
            sessionPermit={"sessionStatus": "RUNNING", "issuedAt": 1},
            pendingOperation={"kind": "AONE_COMMENT", "status": "SENDING"},
        )
        self.assertEqual(worker.stop_check(), 2)

    def test_stop_check_pending_claim_blocks(self):
        """Active session with pending claim → exit 2."""
        self._seed(
            current={"aoneId": "222", "sessionId": "s-5"},
            sessionPermit={"sessionStatus": "RUNNING", "issuedAt": 1},
            pendingClaim={"phase": "CLAIMING"},
        )
        self.assertEqual(worker.stop_check(), 2)

    def test_stop_check_pending_suspend_blocks(self):
        """Active session with pending suspend → exit 2."""
        self._seed(
            current={"aoneId": "333", "sessionId": "s-6"},
            sessionPermit={"sessionStatus": "RUNNING", "issuedAt": 1},
            pendingSuspend={"detail": "waiting"},
        )
        self.assertEqual(worker.stop_check(), 2)

    def test_stop_check_state_missing_returns_fallback(self):
        """No state file → exit 1 (trigger wrap-check fallback)."""
        # Don't seed anything — state file doesn't exist
        self.assertEqual(worker.stop_check(), 1)


if __name__ == "__main__":
    unittest.main(verbosity=2)
