#!/usr/bin/env python3
"""Hermetic tests for bridge/jarvis_task_client.py."""

import io
import json
import socket
import sys
import unittest
from pathlib import Path
from urllib.error import HTTPError, URLError

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from jarvis_task_client import (  # noqa: E402
    ControlPlaneClient,
    ControlPlaneConflict,
    ControlPlaneUnavailable,
    HandoffRequested,
    InvalidResponse,
    StaleFence,
    TaskEnvelope,
)


class FakeResponse:
    def __init__(self, payload=None, status=200, raw=None):
        self.status = status
        self._raw = raw if raw is not None else json.dumps(payload or {}).encode()

    def __enter__(self):
        return self

    def __exit__(self, *_args):
        return False

    def getcode(self):
        return self.status

    def read(self):
        return self._raw


class RecordingOpener:
    def __init__(self, responses=None, error=None):
        self.responses = list(responses or [FakeResponse({"accepted": True})])
        self.error = error
        self.calls = []

    def __call__(self, req, timeout=None):
        self.calls.append((req, timeout))
        if self.error:
            raise self.error
        return self.responses.pop(0)


def envelope():
    return TaskEnvelope(
        task_key="aone:2100304:84345050",
        source_type="AONE",
        source_ref={"project_id": "2100304", "workitem_id": "84345050"},
        task_type="ticket",
        desired_revision="modified:2026-07-15T01:02:03Z",
        trigger_mask=["SCAN"],
        payload={"item_id": "84345050", "prompt": "do work"},
        recovery_policy="RESUME_ONLY",
        priority="high",
        source_status="处理中",
    )


def headers(req):
    return {k.lower(): v for k, v in req.header_items()}


def body(req):
    return json.loads(req.data.decode())


class TaskEnvelopeTest(unittest.TestCase):
    def test_serializes_and_omits_empty_optional_fields(self):
        env = envelope()
        data = env.to_dict()
        self.assertEqual(data["taskKey"], "aone:2100304:84345050")
        self.assertNotIn("executionMode", data)
        self.assertNotIn("shadow", data)
        self.assertEqual(data["sourceRef"]["projectId"], "2100304")
        self.assertEqual(data["payload"]["itemId"], "84345050")
        self.assertEqual(data["priority"], "high")
        self.assertEqual(data["sourceStatus"], "处理中")
        self.assertNotIn("persona", data)

    def test_blank_source_status_is_omitted_from_body(self):
        # An empty/whitespace source_status normalizes to None and is omitted from
        # the serialized body (so the control plane first-insert does not persist
        # NULL). The interactive direct claim must therefore carry a non-empty
        # pre-claim status to avoid the "需要处理" misclassification.
        env = TaskEnvelope(
            task_key="aone:2100304:84345050",
            source_type="AONE",
            source_ref={"project_id": "2100304", "workitem_id": "84345050"},
            task_type="ticket",
            desired_revision="r1",
            payload={"item_id": "84345050"},
        )
        self.assertIsNone(env.source_status)
        self.assertNotIn("sourceStatus", env.to_dict())
        env_blank = TaskEnvelope(
            task_key="aone:2100304:84345050",
            source_type="AONE",
            source_ref={"project_id": "2100304", "workitem_id": "84345050"},
            task_type="ticket",
            desired_revision="r1",
            payload={"item_id": "84345050"},
            source_status="   ",
        )
        self.assertIsNone(env_blank.source_status)
        self.assertNotIn("sourceStatus", env_blank.to_dict())

    def test_request_id_is_stable(self):
        env = envelope()
        first = env.request_id("upsert")
        self.assertEqual(first, env.request_id("upsert"))
        self.assertNotEqual(first, env.request_id("direct-claim"))
        self.assertTrue(first.startswith("jarvis-upsert-"))

    def test_canonical_aone_key_carries_aone_id_for_non_aone_source(self):
        env = TaskEnvelope(
            task_key="aone:2100304:84345050",
            source_type="GITHUB",
            source_ref={"prUrl": "https://example.test/pull/1"},
            task_type="pr_ci_fix",
            desired_revision="pr-ci:abc",
            payload={"itemId": "84345050", "prompt": "fix ci"},
        )
        self.assertEqual(env.aone_id, "84345050")
        self.assertEqual(env.to_dict()["aoneId"], "84345050")

    def test_canonical_aone_key_rejects_conflicting_explicit_aone_id(self):
        with self.assertRaisesRegex(ValueError, "conflicts with canonical task_key"):
            TaskEnvelope(
                task_key="aone:2100304:84345050",
                source_type="GITHUB",
                source_ref={},
                task_type="pr_comment_reply",
                desired_revision="pr-comment:1",
                payload={"itemId": "84345050", "prompt": "reply"},
                aone_id="84399999",
            )

    def test_noncanonical_local_key_does_not_infer_aone_id(self):
        env = TaskEnvelope(
            task_key="aone:unknown:adhoc",
            source_type="LOCAL",
            source_ref={"localId": "adhoc"},
            task_type="ticket",
            desired_revision="local:1",
            payload={"itemId": "adhoc", "prompt": "run"},
        )
        self.assertIsNone(env.aone_id)
        self.assertNotIn("aoneId", env.to_dict())

    def test_rejects_missing_identity_or_non_mapping_payload(self):
        with self.assertRaises(ValueError):
            TaskEnvelope("", "AONE", {}, "ticket", "r1", {})
        with self.assertRaises(TypeError):
            TaskEnvelope("k", "AONE", {}, "ticket", "r1", [])


class ClientContractTest(unittest.TestCase):
    def make(self, opener):
        return ControlPlaneClient(
            "https://pre-agent.example/", "super-secret", timeout=4.5, opener=opener)

    def test_authorization_idempotency_and_upsert_body(self):
        opener = RecordingOpener()
        client = self.make(opener)
        client.upsert_desired_task(envelope(), request_id="req-123")
        req, timeout = opener.calls[0]
        self.assertEqual(req.full_url,
                         "https://pre-agent.example/api/jarvis/v1/tasks/upsert")
        self.assertEqual(timeout, 4.5)
        hs = headers(req)
        self.assertEqual(hs["authorization"], "Bearer super-secret")
        self.assertEqual(hs["idempotency-key"], "req-123")
        self.assertNotIn("x-request-id", hs)
        self.assertNotIn("executionMode", body(req))
        self.assertNotIn("shadow", body(req))
        self.assertEqual(body(req)["desiredRevision"],
                         "modified:2026-07-15T01:02:03Z")

    def test_all_worker_and_session_endpoints(self):
        opener = RecordingOpener(responses=[
            FakeResponse({}), FakeResponse({}), FakeResponse(raw=b"[]"),
            *[FakeResponse({}) for _ in range(8)],
        ])
        c = self.make(opener)
        c.register_worker({"worker_key": "mac:boot:pid"}, request_id="r1")
        c.heartbeat_worker("w1", {"status": "ACTIVE"}, process_uuid="p1", request_id="r2")
        self.assertEqual(
            c.list_force_handoff_requests("w1", process_uuid="p1"), [])
        c.lease_task("w1", lease_seconds=45,
                     capabilities={"kinds": ["probe"]}, process_uuid="p1",
                     request_id="r3")
        c.start_session("s1", "w1", 7, {"pid": 99}, process_uuid="p1", request_id="r4")
        c.heartbeat_session("s1", "w1", 7, {"progress": "running"},
                            process_uuid="p1", request_id="r5")
        c.acknowledge_force_handoff(
            "s1", "w1", 7, process_uuid="p1", request_id="r5-ack")
        c.suspend_session("s1", "w1", 7, {"wait_for": "320687"},
                          process_uuid="p1", request_id="r6")
        c.complete_session("s1", "w1", 7, {"result": "done"},
                           process_uuid="p1", request_id="r7")
        c.fail_session("s1", "w1", 7, {"error": "boom"},
                       process_uuid="p1", request_id="r8")
        c.relinquish_session("s1", "w1", 7, {"reason": "restart"},
                             process_uuid="p1", request_id="r9")

        paths = [call[0].full_url.rsplit("/api/jarvis/v1/", 1)[1]
                 for call in opener.calls]
        self.assertEqual(paths, [
            "workers/register", "workers/w1/heartbeat",
            "workers/w1/handoff-requests?processUuid=p1", "tasks/lease",
            "sessions/s1/start", "sessions/s1/heartbeat",
            "sessions/s1/handoff-ack", "sessions/s1/suspend",
            "sessions/s1/complete", "sessions/s1/fail", "sessions/s1/relinquish",
        ])
        self.assertEqual(body(opener.calls[0][0])["workerKey"], "mac:boot:pid")
        lease = body(opener.calls[3][0])
        self.assertEqual(lease["workerKey"], "w1")
        self.assertEqual(lease["leaseSeconds"], 45)
        self.assertEqual(body(opener.calls[1][0])["processUuid"], "p1")
        handoff_ack = body(opener.calls[6][0])
        self.assertEqual(handoff_ack["oldFenceToken"], 7)
        self.assertNotIn("fenceToken", handoff_ack)
        for req, _timeout in opener.calls[4:6] + opener.calls[7:]:
            payload = body(req)
            self.assertNotIn("sessionId", payload)
            self.assertEqual(payload["workerKey"], "w1")
            self.assertEqual(payload["fenceToken"], 7)
            self.assertEqual(payload["processUuid"], "p1")

    def test_operations_and_queries_use_contract_paths(self):
        opener = RecordingOpener(responses=[
            FakeResponse({"ok": True}), FakeResponse({"ok": True}),
            FakeResponse({"ok": True}), FakeResponse({"ok": True}),
            FakeResponse({"taskId": "t1"}),
            FakeResponse([{"eventType": "LEASED"}]),
            FakeResponse([{"workerKey": "w1"}]),
            FakeResponse({"worker": {"workerKey": "host/one"}}),
            FakeResponse([{"eligibleWorkerCount": 0}]),
            FakeResponse({"items": [], "nextAfterSessionId": None, "hasMore": False}),
        ])
        c = self.make(opener)
        operation = {"operation_id": "op1", "fence_token": 7}
        # begin/ack/fail/reconcile are retained for the interactive-worker receipt engine.
        c.begin_operation(operation, request_id="o1")
        c.ack_operation(operation, request_id="o2")
        c.fail_operation(operation, request_id="o3")
        c.reconcile_operation(operation, request_id="o4")
        task = c.get_task_by_aone("84345050")
        timeline = c.get_task_timeline("task/1")
        workers = c.list_workers()
        worker_state = c.get_worker_state("host/one")
        ready = c.list_ready_task_diagnostics(limit=9)
        waits = c.list_pending_aone_reply_waits(after_session_id=7, limit=25)

        paths = [call[0].full_url.rsplit("/api/jarvis/v1/", 1)[1]
                 for call in opener.calls]
        self.assertEqual(paths, [
            "operations/begin", "operations/ack", "operations/fail",
            "operations/reconcile",
            "tasks/by-aone/84345050",
            "tasks/task%2F1/timeline", "workers", "workers/host%2Fone/state",
            "tasks/ready-diagnostics?limit=9",
            "sessions/waits/aone-reply?afterSessionId=7&limit=25",
        ])
        self.assertEqual(body(opener.calls[0][0]),
                         {"operationId": "op1", "fenceToken": 7})
        self.assertEqual(task["taskId"], "t1")
        self.assertEqual(timeline[0]["eventType"], "LEASED")
        self.assertEqual(workers[0]["workerKey"], "w1")
        self.assertEqual(worker_state["worker"]["workerKey"], "host/one")
        self.assertEqual(ready[0]["eligibleWorkerCount"], 0)
        self.assertEqual(waits["items"], [])
        for req, _timeout in opener.calls[-6:]:
            self.assertEqual(req.get_method(), "GET")
            self.assertIsNone(req.data)
            self.assertNotIn("idempotency-key", headers(req))

    def test_source_status_observation_uses_metadata_only_endpoints(self):
        opener = RecordingOpener(responses=[
            FakeResponse({"items": [], "nextAfterTaskId": None, "hasMore": False}),
            FakeResponse({"id": 42, "sourceStatus": "已发布"}),
        ])
        client = self.make(opener)

        page = client.list_source_status_candidates(after_task_id=12, limit=30)
        updated = client.update_source_status(
            "42", "84386065", "已发布", request_id="source-status-42")

        self.assertFalse(page["hasMore"])
        self.assertEqual(updated["sourceStatus"], "已发布")
        self.assertEqual(opener.calls[0][0].full_url,
                         "https://pre-agent.example/api/jarvis/v1/"
                         "tasks/source-status-candidates?afterTaskId=12&limit=30")
        self.assertEqual(opener.calls[1][0].full_url,
                         "https://pre-agent.example/api/jarvis/v1/tasks/42/source-status")
        self.assertEqual(body(opener.calls[1][0]),
                         {"aoneId": "84386065", "sourceStatus": "已发布"})
        self.assertEqual(headers(opener.calls[1][0])["idempotency-key"],
                         "source-status-42")

    def test_task_attention_uses_control_plane_dedup_contract(self):
        opener = RecordingOpener(responses=[
            FakeResponse({"notify": True, "task": {"id": 42}}),
            FakeResponse({"notify": False, "task": {"id": 42}}),
        ])
        client = self.make(opener)

        stored = client.upsert_task_attention(
            "task/42", "320687", "pr-review:abc", {
                "reason": "CI is green",
                "action": "Review and merge",
                "pr_url": "https://example.test/pull/1",
            }, request_id="attention-task-42")
        cleared = client.clear_task_attention(
            "task/42", request_id="attention-clear-task-42")

        put, delete = [call[0] for call in opener.calls]
        self.assertTrue(put.full_url.endswith(
            "/api/jarvis/v1/tasks/task%2F42/attention"))
        self.assertEqual(put.get_method(), "PUT")
        self.assertEqual(body(put), {
            "ownerStaffId": "320687",
            "eventKey": "pr-review:abc",
            "payload": {
                "reason": "CI is green",
                "action": "Review and merge",
                "prUrl": "https://example.test/pull/1",
            },
        })
        self.assertEqual(headers(put)["idempotency-key"], "attention-task-42")
        self.assertTrue(stored["notify"])
        self.assertEqual(delete.full_url, put.full_url)
        self.assertEqual(delete.get_method(), "DELETE")
        self.assertIsNone(delete.data)
        self.assertEqual(headers(delete)["idempotency-key"],
                         "attention-clear-task-42")
        self.assertFalse(cleared["notify"])

        with self.assertRaises(TypeError):
            client.upsert_task_attention("42", "320687", "key", [])
        with self.assertRaises(ValueError):
            client.upsert_task_attention("42", "", "key", {})

    def test_external_operation_recovery_uses_dedicated_page_and_token_paths(self):
        opener = RecordingOpener(responses=[
            FakeResponse({"items": [], "nextAfterOperationId": None, "hasMore": False}),
            FakeResponse({"proceed": True}), FakeResponse({"proceed": True}),
            FakeResponse({"status": "UNKNOWN"}),
        ])
        client = self.make(opener)

        page = client.list_external_operation_recovery_candidates(
            after_operation_id=12, limit=30)
        client.lease_operation_recovery(44, "worker-1", "token-1", request_id="lease")
        client.renew_operation_recovery(44, "worker-1", "token-1", request_id="renew")
        client.release_operation_recovery(44, "worker-1", "token-1", request_id="release")

        self.assertFalse(page["hasMore"])
        paths = [call[0].full_url.rsplit("/api/jarvis/v1/", 1)[1]
                 for call in opener.calls]
        self.assertEqual(paths, [
            "operations/external-recovery-candidates?afterOperationId=12&limit=30",
            "operations/recovery/lease", "operations/recovery/renew",
            "operations/recovery/release",
        ])
        for request, _timeout in opener.calls[1:]:
            self.assertEqual(body(request), {
                "operationId": 44,
                "workerKey": "worker-1",
                "recoveryToken": "token-1",
            })

        with self.assertRaises(ValueError):
            client.list_external_operation_recovery_candidates(after_operation_id=-1)
        with self.assertRaises(ValueError):
            client.list_external_operation_recovery_candidates(limit=501)

    def test_operation_point_read_and_not_started_use_distinct_contracts(self):
        opener = RecordingOpener(responses=[
            FakeResponse({"operation": {"id": 44}}),
            FakeResponse({"operation": {"id": 44}}),
            FakeResponse({"id": 44, "status": "FAILED_NOT_STARTED"}),
        ])
        client = self.make(opener)

        client.get_operation(44)
        client.get_operation_by_key(9, 3, "comment:key/one")
        client.mark_operation_not_started({
            "operation_id": 44, "worker_key": "worker-1", "fence_token": 7,
            "process_uuid": "process-1", "reason": "rejected before send",
        }, request_id="not-started-44")

        paths = [call[0].full_url.rsplit("/api/jarvis/v1/", 1)[1]
                 for call in opener.calls]
        self.assertEqual(paths, [
            "operations/44",
            "operations/by-key?taskId=9&generation=3&operationKey=comment%3Akey%2Fone",
            "operations/not-started",
        ])
        self.assertEqual(body(opener.calls[2][0]), {
            "operationId": 44, "workerKey": "worker-1", "fenceToken": 7,
            "processUuid": "process-1", "reason": "rejected before send",
        })

    def test_pending_wait_query_validates_keyset_bounds(self):
        c = self.make(RecordingOpener())
        with self.assertRaises(ValueError):
            c.list_pending_aone_reply_waits(after_session_id=-1)
        with self.assertRaises(ValueError):
            c.list_pending_aone_reply_waits(limit=501)
        with self.assertRaises(ValueError):
            c.list_ready_task_diagnostics(limit=0)
        with self.assertRaises(ValueError):
            c.list_ready_task_diagnostics(limit=501)

    def test_discard_resume_context_posts_exact_session_and_reason(self):
        opener = RecordingOpener(responses=[FakeResponse({"id": 42, "status": "READY"})])
        c = self.make(opener)

        result = c.discard_resume_context(
            "task/42", 7, "original worker retired", request_id="discard-42-7")

        req, _timeout = opener.calls[0]
        self.assertTrue(req.full_url.endswith(
            "/api/jarvis/v1/tasks/task%2F42/discard-resume-context"))
        self.assertEqual(body(req), {
            "expectedSessionId": 7,
            "reason": "original worker retired",
        })
        self.assertEqual(headers(req)["idempotency-key"], "discard-42-7")
        self.assertEqual(result["status"], "READY")

    def test_direct_claim_is_targeted_task_and_allows_zero_free_slots(self):
        opener = RecordingOpener()
        c = self.make(opener)
        c.claim_task("worker:interactive", envelope(),
                     runtime_session_id="interactive:codex:s:aone:p:i:cycle:1",
                     lease_seconds=120, free_slots=0, request_id="claim-1")
        req, _timeout = opener.calls[0]
        self.assertTrue(req.full_url.endswith("/api/jarvis/v1/tasks/claim"))
        payload = body(req)
        self.assertEqual(payload["workerKey"], "worker:interactive")
        self.assertEqual(payload["runtimeSessionId"],
                         "interactive:codex:s:aone:p:i:cycle:1")
        self.assertEqual(payload["freeSlots"], 0)
        self.assertNotIn("executionMode", payload["task"])
        self.assertNotIn("shadow", payload["task"])
        self.assertEqual(headers(req)["idempotency-key"], "claim-1")

        with self.assertRaises(ValueError):
            c.claim_task("w", envelope(), runtime_session_id="r", free_slots=-1)

    def test_bearer_token_is_optional(self):
        opener = RecordingOpener()
        c = ControlPlaneClient("https://pre-agent.example", opener=opener)
        c.register_worker({"workerKey": "w1"}, request_id="r1")
        self.assertNotIn("authorization", headers(opener.calls[0][0]))

    def test_409_maps_conflict(self):
        err = HTTPError("https://x", 409, "Conflict", {},
                        io.BytesIO(b'{"code":"ACTIVE_SESSION","message":"busy"}'))
        with self.assertRaises(ControlPlaneConflict) as got:
            self.make(RecordingOpener(error=err)).lease_task("w1", request_id="r")
        self.assertEqual(got.exception.status, 409)
        self.assertEqual(got.exception.code, "ACTIVE_SESSION")

    def test_412_and_stale_fence_code_map_stale_fence(self):
        for status in (409, 412):
            err = HTTPError("https://x", status, "Precondition", {},
                            io.BytesIO(b'{"code":"STALE_FENCE","message":"old"}'))
            with self.assertRaises(StaleFence):
                self.make(RecordingOpener(error=err)).heartbeat_session(
                    "s", "w", 1, request_id="r")

    def test_handoff_requested_code_maps_specific_stale_fence(self):
        err = HTTPError(
            "https://x", 412, "Precondition", {},
            io.BytesIO(b'{"code":"PreconditionFailed.HandoffRequested",'
                         b'"message":"stop and acknowledge"}'))

        with self.assertRaises(HandoffRequested):
            self.make(RecordingOpener(error=err)).heartbeat_session(
                "s", "w", 1, request_id="r")

    def test_handoff_request_poll_tolerates_older_control_plane_404(self):
        err = HTTPError(
            "https://x", 404, "Not Found", {},
            io.BytesIO(b'{"message":"not found"}'))

        requests = self.make(RecordingOpener(error=err)).list_force_handoff_requests(
            "w", process_uuid="p")

        self.assertEqual(requests, [])

    def test_network_and_5xx_are_unavailable_without_token_leak(self):
        with self.assertRaises(ControlPlaneUnavailable) as got:
            self.make(RecordingOpener(error=URLError(socket.timeout()))).lease_task(
                "w1", request_id="r")
        self.assertNotIn("super-secret", str(got.exception))

        err = HTTPError("https://x", 503, "Unavailable", {},
                        io.BytesIO(b'{"message":"down"}'))
        with self.assertRaises(ControlPlaneUnavailable):
            self.make(RecordingOpener(error=err)).register_worker(
                {"worker_key": "w"}, request_id="r")

    def test_invalid_success_json_is_rejected(self):
        opener = RecordingOpener(responses=[FakeResponse(raw=b"not-json")])
        with self.assertRaises(InvalidResponse):
            self.make(opener).register_worker({"worker_key": "w"}, request_id="r")

    def test_session_requires_fence_and_lease_rejects_nonpositive_ttl(self):
        c = self.make(RecordingOpener())
        with self.assertRaises(ValueError):
            c.start_session("s", "w", None)
        with self.assertRaises(ValueError):
            c.lease_task("w", lease_seconds=0)


if __name__ == "__main__":
    unittest.main(verbosity=2)
