#!/usr/bin/env python3

import hashlib
import sys
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from bridge.jarvis_external_recovery import (  # noqa: E402
    ExternalOperationRecoveryScheduler,
    RecoveryInconclusive,
)


def candidate(operation_id=11, operation_type="AONE_COMMENT", payload=None):
    return {
        "task": {"id": 1, "aoneId": "84345050"},
        "operation": {
            "id": operation_id,
            "operationType": operation_type,
            "target": "84345050",
            "createdAt": "2026-07-21T00:00:00Z",
        },
        "readbackSpec": payload or {
            "digest": hashlib.sha256(b"hello").hexdigest(), "preview": "hello"},
    }


class FakeClient:
    def __init__(self, pages=None):
        self.pages = list(pages or [{"items": [], "hasMore": False}])
        self.calls = []

    def list_external_operation_recovery_candidates(self, **kwargs):
        self.calls.append(("list", kwargs))
        return self.pages.pop(0)

    def lease_operation_recovery(self, *args, **kwargs):
        self.calls.append(("lease", args, kwargs))
        return {"proceed": True}

    def renew_operation_recovery(self, *args, **kwargs):
        self.calls.append(("renew", args, kwargs))
        return {"proceed": True}

    def release_operation_recovery(self, *args, **kwargs):
        self.calls.append(("release", args, kwargs))
        return {"status": "UNKNOWN"}

    def reconcile_operation(self, *args, **kwargs):
        self.calls.append(("reconcile", args, kwargs))
        return {"status": "ACKED"}


class ExternalOperationRecoverySchedulerTest(unittest.TestCase):
    def scheduler(self, client, observer):
        return ExternalOperationRecoveryScheduler(
            client, "worker-1", repo_root=HERE.parent, enabled=True,
            observer=observer, interval=1)

    def test_found_readback_is_token_leased_renewed_and_reconciled(self):
        client = FakeClient()
        scheduler = self.scheduler(client, lambda _candidate: (True, "aone:comment:12"))

        scheduler._recover(candidate())

        names = [call[0] for call in client.calls]
        self.assertEqual(names, ["lease", "renew", "renew", "reconcile"])
        payload = next(call[1][0] for call in client.calls if call[0] == "reconcile")
        self.assertTrue(payload["found"])
        self.assertEqual(payload["externalRef"], "aone:comment:12")
        self.assertIn("recoveryToken", payload)

        first_token = client.calls[0][1][2]
        client.calls.clear()
        scheduler._recover(candidate())
        self.assertEqual(client.calls[0][1][2], first_token)

    def test_not_found_is_definite_and_inconclusive_releases_fail_closed(self):
        client = FakeClient()
        scheduler = self.scheduler(client, lambda _candidate: (False, None))
        scheduler._recover(candidate())
        payload = next(call[1][0] for call in client.calls if call[0] == "reconcile")
        self.assertFalse(payload["found"])
        self.assertTrue(payload["retryAllowed"])
        self.assertNotIn("release", [call[0] for call in client.calls])

        unavailable = FakeClient()
        scheduler = self.scheduler(
            unavailable, lambda _candidate: (_ for _ in ()).throw(
                RecoveryInconclusive("Aone unavailable")))
        scheduler._recover(candidate())
        self.assertEqual([call[0] for call in unavailable.calls],
                         ["lease", "renew", "release"])

    def test_tick_pages_by_control_plane_cursor(self):
        client = FakeClient([
            {"items": [candidate(11)], "hasMore": True, "nextAfterOperationId": 11},
            {"items": [candidate(12)], "hasMore": False,
             "nextAfterOperationId": None},
        ])
        scheduler = self.scheduler(client, lambda _candidate: (True, "ref"))

        scheduler._tick()

        list_calls = [call for call in client.calls if call[0] == "list"]
        self.assertEqual(list_calls[0][1]["after_operation_id"], 0)
        self.assertEqual(list_calls[1][1]["after_operation_id"], 11)

    def test_lease_competition_never_reads_aone(self):
        client = FakeClient()
        client.lease_operation_recovery = mock.Mock(return_value={"proceed": False})
        observer = mock.Mock()
        scheduler = self.scheduler(client, observer)

        scheduler._recover(candidate())

        observer.assert_not_called()
        client.lease_operation_recovery.assert_called_once()

    def test_aone_readback_covers_comment_status_and_tags(self):
        scheduler = self.scheduler(FakeClient(), lambda _candidate: (False, None))
        with mock.patch.object(scheduler, "_a1_json", return_value=[
                {"id": 91, "content": "hello", "createdAt": "2026-07-21T00:01:00Z"}]):
            self.assertEqual(scheduler._observe_aone(candidate()),
                             (True, "aone:84345050:comment:91"))

        status_candidate = candidate(
            12, "AONE_STATUS", {"kind": "status", "material": "已完成"})
        with mock.patch.object(scheduler, "_a1_json", return_value={"status": "已完成"}):
            self.assertEqual(scheduler._observe_aone(status_candidate),
                             (True, "aone:84345050:status:已完成"))

        tag_candidate = candidate(
            13, "AONE_RELEASE", {"kind": "release-tag", "material": "idle"})
        with mock.patch.object(scheduler, "_a1_json", return_value={"tag": "jarvis-idle,x"}):
            self.assertEqual(scheduler._observe_aone(tag_candidate),
                             (True, "aone:84345050:tag:jarvis-idle"))

        claim_candidate = candidate(14, "AONE_CLAIM", {
            "addTag": "jarvis-claimed", "removeTag": ["jarvis-idle", "old"],
        })
        structured = {"fields": [{
            "identifier": "tag",
            "displayValue": [{"name": "jarvis-claimed"}, {"value": "other"}],
        }]}
        with mock.patch.object(scheduler, "_a1_json", return_value=structured):
            self.assertEqual(scheduler._observe_aone(claim_candidate),
                             (True, "aone:84345050:tag:jarvis-claimed"))
        structured["fields"][0]["displayValue"].append({"name": "jarvis-idle"})
        with mock.patch.object(scheduler, "_a1_json", return_value=structured):
            with self.assertRaises(RecoveryInconclusive):
                scheduler._observe_aone(claim_candidate)

    def test_readback_only_reports_absence_when_it_is_authoritative(self):
        scheduler = self.scheduler(FakeClient(), lambda _candidate: (False, None))
        with mock.patch.object(scheduler, "_a1_json", return_value=[
                {"id": 80, "content": "hello", "createdAt": "2026-07-20T23:59:00Z"}]):
            self.assertEqual(scheduler._observe_aone(candidate()), (False, None))

        status_candidate = candidate(
            12, "AONE_STATUS", {"kind": "status", "material": "已完成"})
        with mock.patch.object(scheduler, "_a1_json", return_value={"status": "开发中"}):
            with self.assertRaises(RecoveryInconclusive):
                scheduler._observe_aone(status_candidate)

        release_candidate = candidate(
            13, "AONE_RELEASE", {"kind": "release-tag", "material": "idle"})
        with mock.patch.object(scheduler, "_a1_json", return_value={
                "tag": "jarvis-idle,jarvis-claimed"}):
            with self.assertRaises(RecoveryInconclusive):
                scheduler._observe_aone(release_candidate)

    def test_cursor_rotates_across_ticks_when_page_budget_is_reached(self):
        client = FakeClient([
            {"items": [candidate(11)], "hasMore": True, "nextAfterOperationId": 11},
            {"items": [candidate(12)], "hasMore": False, "nextAfterOperationId": None},
        ])
        scheduler = ExternalOperationRecoveryScheduler(
            client, "worker-1", repo_root=HERE.parent, enabled=True,
            observer=lambda _candidate: (True, "ref"), interval=1, max_pages=1)

        scheduler._tick()
        scheduler._tick()

        list_calls = [call for call in client.calls if call[0] == "list"]
        self.assertEqual([call[1]["after_operation_id"] for call in list_calls], [0, 11])
        self.assertEqual(scheduler._after_operation_id, 0)

    def test_stop_after_lease_releases_without_observing(self):
        client = FakeClient()
        scheduler = self.scheduler(client, mock.Mock())

        def lease_and_stop(*_args, **_kwargs):
            scheduler._stop.set()
            return {"proceed": True}

        client.lease_operation_recovery = lease_and_stop
        scheduler._recover(candidate())

        self.assertEqual([call[0] for call in client.calls], ["release"])
        scheduler.observer.assert_not_called()

    def test_untrusted_candidate_and_ambiguous_comment_fail_closed(self):
        scheduler = self.scheduler(FakeClient(), lambda _candidate: (False, None))
        wrong_target = candidate()
        wrong_target["operation"]["target"] = "84399999"
        with self.assertRaises(RecoveryInconclusive):
            scheduler._observe_aone(wrong_target)

        bad_payload = candidate()
        bad_payload["readbackSpec"] = {"preview": "no digest"}
        with self.assertRaises(RecoveryInconclusive):
            scheduler._observe_aone(bad_payload)

        with mock.patch.object(scheduler, "_a1_json", return_value=[{
                "content": "hello", "createdAt": "2026-07-21T00:01:00Z"}]):
            with self.assertRaises(RecoveryInconclusive):
                scheduler._observe_aone(candidate())

        with mock.patch(
                "bridge.jarvis_external_recovery.run_process_group",
                return_value=SimpleNamespace(returncode=0, stdout="not-json")):
            with self.assertRaises(RecoveryInconclusive):
                scheduler._a1_json(["project", "workitem", "get", "84345050"])


if __name__ == "__main__":
    unittest.main(verbosity=2)
