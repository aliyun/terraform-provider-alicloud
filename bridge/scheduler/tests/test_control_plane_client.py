from __future__ import annotations

from datetime import datetime, timedelta, timezone
import io
import json
import socket
import unittest
from urllib.error import HTTPError, URLError

from bridge.scheduler import (
    HttpScheduledJobControlPlane, JobRegistration, ScheduledJobControlPlaneProtocolError,
    ScheduledJobControlPlaneRejected, ScheduledJobControlPlaneUnavailable, ScheduledJobStatus,
)


UTC = timezone.utc


def at(hour: int, minute: int = 0) -> datetime:
    return datetime(2026, 7, 21, hour, minute, tzinfo=UTC)


def state(*, status: str = "IDLE", next_run_at: str | None = "2026-07-21T09:00:00.000+00:00"):
    return {"jobKey": "aone.scan", "status": status, "nextRunAt": next_run_at}


class FakeResponse:
    def __init__(self, payload=None, *, status=200, raw=None):
        self.status = status
        self._raw = raw if raw is not None else json.dumps(payload).encode("utf-8")

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
        self.responses = list(responses or [FakeResponse({})])
        self.error = error
        self.calls = []

    def __call__(self, request, timeout=None):
        self.calls.append((request, timeout))
        if self.error:
            raise self.error
        return self.responses.pop(0)


def payload(request):
    return json.loads(request.data.decode("utf-8"))


def headers(request):
    return {key.lower(): value for key, value in request.header_items()}


class HttpScheduledJobControlPlaneTests(unittest.TestCase):
    def client(self, opener):
        return HttpScheduledJobControlPlane(
            "https://pre-agent.example/", "test-token", timeout=4.5, opener=opener)

    def test_register_serializes_contract_and_normalizes_to_utc(self):
        opener = RecordingOpener([FakeResponse([state()])])
        registration = JobRegistration(
            "aone.scan", "Aone scan", {"revision": 1},
            datetime(2026, 7, 21, 17, tzinfo=timezone(timedelta(hours=8))), True,
        )

        states = self.client(opener).register([registration])

        request, timeout = opener.calls[0]
        self.assertEqual(timeout, 4.5)
        self.assertEqual(request.full_url, "https://pre-agent.example/api/jarvis/v1/scheduled-jobs/register")
        self.assertEqual(headers(request)["authorization"], "Bearer test-token")
        self.assertEqual(payload(request), {"jobs": [{
            "jobKey": "aone.scan", "jobName": "Aone scan", "definition": {"revision": 1},
            "nextRunAt": "2026-07-21T09:00:00.000+00:00", "enabled": True,
        }]})
        self.assertEqual(states[0].status, ScheduledJobStatus.IDLE)
        self.assertEqual(states[0].next_run_at, at(9))

    def test_list_parses_utc_state_and_start_requires_explicit_admitted_boolean(self):
        opener = RecordingOpener([
            FakeResponse([state(status="ERROR", next_run_at="2026-07-21T17:00:00.000+08:00")]),
            FakeResponse({"admitted": False, "job": state(status="RUNNING")}),
        ])
        client = self.client(opener)

        states = client.list_jobs()
        admitted = client.start("aone.scan", at(9), at(10))

        self.assertEqual(states[0].next_run_at, at(9))
        self.assertFalse(admitted)
        request, _ = opener.calls[1]
        self.assertEqual(payload(request), {
            "scheduledFor": "2026-07-21T09:00:00.000+00:00",
            "nextRunAt": "2026-07-21T10:00:00.000+00:00",
        })
        self.assertTrue(request.full_url.endswith("/aone.scan/start"))

    def test_recover_interrupted_posts_without_body_and_strictly_parses_response(self):
        opener = RecordingOpener([FakeResponse({
            "recovered": 1,
            "jobs": [state(status="IDLE", next_run_at="2026-07-21T09:00:00.000+00:00")],
        })])

        recovered = self.client(opener).recover_interrupted()

        self.assertEqual([(job.job_key, job.status, job.next_run_at) for job in recovered], [
            ("aone.scan", ScheduledJobStatus.IDLE, at(9)),
        ])
        request, _ = opener.calls[0]
        self.assertEqual(request.full_url, "https://pre-agent.example/api/jarvis/v1/scheduled-jobs/recover-interrupted")
        self.assertEqual(request.get_method(), "POST")
        self.assertIsNone(request.data)

    def test_recover_interrupted_response_and_transport_fail_closed(self):
        malformed = (
            {},
            {"recovered": True, "jobs": []},
            {"recovered": 1, "jobs": []},
            {"recovered": 1, "jobs": [state(status="RUNNING")]},
        )
        for response in malformed:
            with self.subTest(response=response):
                with self.assertRaises(ScheduledJobControlPlaneProtocolError):
                    self.client(RecordingOpener([FakeResponse(response)])).recover_interrupted()
        with self.assertRaises(ScheduledJobControlPlaneRejected):
            self.client(RecordingOpener(error=HTTPError(
                "https://x", 409, "conflict", {}, io.BytesIO(b"conflict")))).recover_interrupted()
        with self.assertRaises(ScheduledJobControlPlaneProtocolError):
            self.client(RecordingOpener([FakeResponse(raw=b"not-json")])).recover_interrupted()

    def test_complete_and_fail_serialize_exact_slot_and_next_run(self):
        opener = RecordingOpener([
            FakeResponse(state(next_run_at="2026-07-21T10:00:00.000+00:00")),
            FakeResponse(state(status="ERROR", next_run_at="2026-07-21T09:00:05.000+00:00")),
            FakeResponse(state(status="ERROR", next_run_at=None)),
        ])
        client = self.client(opener)

        client.complete("aone.scan", at(9), at(10))
        client.fail("aone.scan", at(9), retryable=True, next_run_at=at(9, 0) + timedelta(seconds=5), error="retry")
        client.fail("aone.scan", at(9), retryable=False, next_run_at=None, error="permanent")

        self.assertEqual(payload(opener.calls[0][0]), {
            "scheduledFor": "2026-07-21T09:00:00.000+00:00",
            "nextRunAt": "2026-07-21T10:00:00.000+00:00",
        })
        self.assertEqual(payload(opener.calls[1][0]), {
            "scheduledFor": "2026-07-21T09:00:00.000+00:00", "retryable": True,
            "nextRunAt": "2026-07-21T09:00:05.000+00:00", "error": "retry",
        })
        self.assertIsNone(payload(opener.calls[2][0])["nextRunAt"])

    def test_http_network_and_invalid_json_fail_closed(self):
        unavailable = HTTPError("https://x", 503, "unavailable", {}, io.BytesIO(b"down"))
        with self.assertRaises(ScheduledJobControlPlaneUnavailable):
            self.client(RecordingOpener(error=unavailable)).list_jobs()
        rejected = HTTPError("https://x", 409, "conflict", {}, io.BytesIO(b"conflict"))
        with self.assertRaises(ScheduledJobControlPlaneRejected):
            self.client(RecordingOpener(error=rejected)).list_jobs()
        with self.assertRaises(ScheduledJobControlPlaneUnavailable):
            self.client(RecordingOpener(error=URLError(socket.timeout()))).list_jobs()
        with self.assertRaises(ScheduledJobControlPlaneProtocolError):
            self.client(RecordingOpener([FakeResponse(raw=b"not-json")])).list_jobs()

    def test_missing_or_non_boolean_admitted_is_a_protocol_error(self):
        for response in ({"job": state(status="RUNNING")}, {"admitted": "true", "job": state(status="RUNNING")}):
            with self.subTest(response=response):
                with self.assertRaises(ScheduledJobControlPlaneProtocolError):
                    self.client(RecordingOpener([FakeResponse(response)])).start(
                        "aone.scan", at(9), at(10))

    def test_constructor_reads_existing_control_plane_environment_only_when_constructed(self):
        client = HttpScheduledJobControlPlane(environ={
            "JARVIS_CONTROL_PLANE_BASE_URL": "https://control.example/",
            "JARVIS_CONTROL_PLANE_TOKEN": "from-env",
            "JARVIS_CONTROL_PLANE_TIMEOUT": "7.5",
        })
        self.assertEqual(client.base_url, "https://control.example")
        self.assertEqual(client.token, "from-env")
        self.assertEqual(client.timeout, 7.5)


if __name__ == "__main__":
    unittest.main()
