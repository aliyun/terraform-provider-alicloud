#!/usr/bin/env python3
"""Hermetic lifecycle tests for the interactive Claude/Codex worker."""

import contextlib
import importlib.util
import io
import os
import stat
import tempfile
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
        return self._record("heartbeat_worker", *args, **kwargs)

    def heartbeat_session(self, *args, **kwargs):
        return self._record("heartbeat_session", *args, **kwargs)

    def claim_task(self, *args, **kwargs):
        self._record("claim_task", *args, **kwargs)
        if self.claim_error:
            raise self.claim_error
        return {
            "task": {"id": "task-1", "generation": 4},
            "session": {"id": "session-1", "generation": 4,
                        "fenceToken": 9, "attemptNo": 2},
        }

    def start_session(self, *args, **kwargs):
        return self._record("start_session", *args, **kwargs)

    def begin_operation(self, *args, **kwargs):
        self._record("begin_operation", *args, **kwargs)
        if self.begin_error:
            raise self.begin_error
        return self.begin_results.pop(0)

    def ack_operation(self, *args, **kwargs):
        return self._record("ack_operation", *args, **kwargs)

    def fail_operation(self, *args, **kwargs):
        return self._record("fail_operation", *args, **kwargs)

    def fail_session(self, *args, **kwargs):
        return self._record("fail_session", *args, **kwargs)

    def suspend_session(self, *args, **kwargs):
        return self._record("suspend_session", *args, **kwargs)

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
            "stopped": False,
        }
        self._store().save(state)
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

    def test_pid_reuse_cannot_keep_an_old_worker_alive(self):
        state = self._seed()
        state["hostProcessStartedAt"] = "old-birth"
        with mock.patch.object(worker, "_pid_alive", return_value=True), \
                mock.patch.object(worker, "_process_info",
                                  return_value=(1, "/Applications/Codex.app/Contents/MacOS/Codex")), \
                mock.patch.object(worker, "_process_start_identity",
                                  return_value="replacement-birth"):
            self.assertFalse(worker._host_alive(state))

    def test_claim_retry_keeps_cycle_revision_and_resumes_local_sending_receipt(self):
        self._seed()
        fake = FakeClient()
        fake.claim_error = worker.ControlPlaneUnavailable("lost response")
        with mock.patch.object(worker, "_client", return_value=fake):
            with self.assertRaises(worker.ControlPlaneUnavailable):
                worker.prepare_claim("84345050", "2100304")
        pending = self._store().load()["pendingClaim"]
        self.assertEqual(pending["cycle"], 1)
        first_runtime = pending["runtimeSessionId"]

        fake.claim_error = None
        with mock.patch.object(worker, "_client", return_value=fake):
            first = worker.prepare_claim("84345050", "2100304")
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
            retry = worker.prepare_claim("84345050", "2100304")
        retry_claim = [c for c in fake.calls if c[0] == "claim_task"][-1]
        retry_envelope = retry_claim[1][1]
        self.assertTrue(retry["proceed"])
        self.assertEqual(retry["runtimeSessionId"], first_runtime)
        self.assertEqual(retry_envelope.desired_revision,
                         first_envelope.desired_revision)
        self.assertEqual(retry_claim[2]["free_slots"], 0)

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
            next_claim = worker.prepare_claim("84345050", "2100304")
        self.assertNotEqual(next_claim["runtimeSessionId"], first_runtime)
        self.assertIn("cycle:2", next_claim["runtimeSessionId"])

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
        offline = [c for c in fake.calls if c[0] == "heartbeat_worker"][-1]
        self.assertEqual(offline[1][1]["status"], "OFFLINE")


if __name__ == "__main__":
    unittest.main(verbosity=2)
