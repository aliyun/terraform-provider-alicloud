"""Tests for bridge.worker_offline: identity persistence + best-effort OFFLINE."""

from __future__ import annotations

import json
import tempfile
import unittest
from pathlib import Path

from bridge import worker_offline


class _RecordingClient:
    """Minimal control-plane stub recording OFFLINE heartbeats."""

    def __init__(self, *, fail: bool = False) -> None:
        self.calls: list[tuple[str, dict, str]] = []
        self.fail = fail

    def heartbeat_worker(self, worker_key, heartbeat=None, *,
                         process_uuid=None, request_id=None):
        if self.fail:
            raise RuntimeError("control plane 409")
        self.calls.append((worker_key, dict(heartbeat or {}), process_uuid))
        return {"status": "OFFLINE"}


class WorkerOfflineTest(unittest.TestCase):
    def setUp(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.dir = Path(self._tmp.name)
        self.addCleanup(self._tmp.cleanup)

    def test_write_then_clear_identity_roundtrip(self):
        worker_offline.write_identity(
            "scheduler", worker_key="bridge-scheduler",
            process_uuid="uuid-1", host_id="h1", boot_id="b1",
            directory=self.dir)
        path = self.dir / ("scheduler%s" % worker_offline.CP_IDENTITY_SUFFIX)
        self.assertTrue(path.exists())
        record = json.loads(path.read_text(encoding="utf-8"))
        self.assertEqual(record["workerKey"], "bridge-scheduler")
        self.assertEqual(record["processUuid"], "uuid-1")

        worker_offline.clear_identity("scheduler", directory=self.dir)
        self.assertFalse(path.exists())
        # idempotent: clearing a missing identity does not raise
        worker_offline.clear_identity("scheduler", directory=self.dir)

    def test_offline_all_posts_offline_with_persisted_uuid_and_removes_file(self):
        worker_offline.write_identity(
            "scheduler", worker_key="bridge-scheduler",
            process_uuid="dead-uuid", host_id="host-a", boot_id="boot-a",
            directory=self.dir)
        worker_offline.write_identity(
            "persistent-worker", worker_key="pw-key",
            process_uuid="pw-uuid", host_id="host-a", boot_id="boot-a",
            directory=self.dir)
        client = _RecordingClient()

        count = worker_offline.offline_all(client=client, directory=self.dir)

        self.assertEqual(count, 2)
        keys = {c[0] for c in client.calls}
        self.assertEqual(keys, {"bridge-scheduler", "pw-key"})
        by_key = {c[0]: c for c in client.calls}
        # OFFLINE carries the persisted (dead owner's) processUuid + status.
        self.assertEqual(by_key["bridge-scheduler"][2], "dead-uuid")
        self.assertEqual(by_key["bridge-scheduler"][1]["status"], "OFFLINE")
        # identity files removed after a successful OFFLINE
        self.assertEqual(list(self.dir.glob("*%s" % worker_offline.CP_IDENTITY_SUFFIX)), [])

    def test_offline_all_keeps_file_and_never_raises_when_client_rejects(self):
        worker_offline.write_identity(
            "scheduler", worker_key="bridge-scheduler",
            process_uuid="dead-uuid", host_id="host-a", boot_id="boot-a",
            directory=self.dir)
        client = _RecordingClient(fail=True)

        count = worker_offline.offline_all(client=client, directory=self.dir)

        self.assertEqual(count, 0)
        # file left in place so the reaper / a retry remains the fallback
        path = self.dir / ("scheduler%s" % worker_offline.CP_IDENTITY_SUFFIX)
        self.assertTrue(path.exists())

    def test_offline_all_no_files_is_noop(self):
        self.assertEqual(
            worker_offline.offline_all(client=_RecordingClient(), directory=self.dir),
            0)

    def test_offline_all_skips_identity_missing_uuid(self):
        path = self.dir / ("scheduler%s" % worker_offline.CP_IDENTITY_SUFFIX)
        path.write_text(json.dumps({"workerKey": "bridge-scheduler"}),
                        encoding="utf-8")
        client = _RecordingClient()
        count = worker_offline.offline_all(client=client, directory=self.dir)
        self.assertEqual(count, 0)
        self.assertEqual(client.calls, [])

    def test_offline_all_without_token_returns_zero(self):
        worker_offline.write_identity(
            "scheduler", worker_key="bridge-scheduler",
            process_uuid="dead-uuid", host_id="host-a", boot_id="boot-a",
            directory=self.dir)
        # No client passed and no token in env -> builds no client, no raise.
        count = worker_offline.offline_all(
            client=None, environ={}, directory=self.dir)
        self.assertEqual(count, 0)

    def test_offline_helper_client_never_retains_admin_token(self):
        client = worker_offline._client_from_env({
            "JARVIS_CONTROL_PLANE_TOKEN": "worker-token",
            "JARVIS_CONTROL_PLANE_ADMIN_TOKEN": "operator-only",
        }, 3.0)

        self.assertEqual(client.token, "worker-token")
        self.assertEqual(client.admin_token, "")


if __name__ == "__main__":
    unittest.main()
