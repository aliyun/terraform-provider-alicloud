#!/usr/bin/env python3
"""Scan must not re-publish a desired revision the control plane already holds.

An upsert is a desired-state write, but on a **resting** Task (SUCCEEDED /
CANCELED / FAILED) the control plane also re-arms it to READY — observed on
Task 665 as `TASK_UPSERTED SUCCEEDED→READY` seven times and
`CANCELED→READY` once. So an information-free re-dispatch (same
desiredRevision) resurrects an already-finished Task, a Persistent Worker
leases it, and a full execution slot burns 20-40 minutes re-reaching the same
conclusion. Five such zombies filled worker-bba9's capacity and starved new
tickets with NO_FREE_SLOTS.

`TaskEnvelope.request_id` hashes the whole envelope, not just the revision, so
any unrelated drift (sourceStatus, title, prompt) mints a fresh idempotency key
and the control plane's own dedup cannot catch these either.
"""

import unittest
from unittest import mock

from bridge.scheduler.runners import scan


ITEM = {
    "id": "84202521",
    "title": "scan must not resurrect a finished task",
    "pool": "tf_customer",
    "pool_project": "1086837",
    "modified": "2026-07-29 18:22",
    "status": "处理中",
    "tag": ["jarvis-idle"],
}


def _runner(rows):
    """A ScanRunner wired to a control-plane router whose read returns `rows`."""
    runner = scan.ScanRunner.__new__(scan.ScanRunner)
    runner.pool = None
    runner.handler = None
    runner._woken_epochs = {}
    client = mock.Mock()
    client.get_task_by_aone.side_effect = (
        rows if callable(rows) else lambda _aone_id: rows)
    router = mock.Mock()
    router.client = client
    router.is_task.return_value = True
    router.enqueue.return_value = scan.EnqueueResult(True, "task_persisted")
    runner.execution_router = router
    return runner, router


def _row(runner, item, **overrides):
    """A live Task row already carrying the exact revision this item would send."""
    envelope = runner._envelope(dict(item))
    row = {
        "id": "665",
        "status": "SUCCEEDED",
        "sourceType": "AONE",
        "taskType": "ticket",
        "recoveryPolicy": "RESUME_ONLY",
        "generation": 12,
        "desiredRevision": envelope.desired_revision,
        "processingRevision": envelope.desired_revision,
        "processedRevision": "comment:125161724",
    }
    row.update(overrides)
    return row


class ScanDesiredRevisionNoopTest(unittest.TestCase):
    def test_identical_desired_revision_is_not_re_upserted(self):
        runner, router = _runner([])
        row = _row(runner, ITEM)
        router.client.get_task_by_aone.side_effect = lambda _aone_id: [row]

        ok, reason = runner._dispatch(dict(ITEM), force=True)

        router.enqueue.assert_not_called()
        self.assertTrue(ok, "a no-op is a resting state, not a failed dispatch")
        self.assertEqual(reason, "noop_desired_revision_unchanged")

    def test_resting_task_is_not_resurrected(self):
        for status in ("SUCCEEDED", "CANCELED", "CANCELLED", "FAILED"):
            with self.subTest(status=status):
                runner, router = _runner([])
                row = _row(runner, ITEM, status=status)
                router.client.get_task_by_aone.side_effect = (
                    lambda _aone_id: [row])

                runner._dispatch(dict(ITEM), force=True)

                router.enqueue.assert_not_called()

    def test_advanced_desired_revision_still_upserts(self):
        """A real new signal must always reach the control plane."""
        runner, router = _runner([])
        row = _row(
            runner, ITEM,
            desiredRevision="modified:2026-07-01 09:00|policy:v6|input:deadbeef")
        router.client.get_task_by_aone.side_effect = lambda _aone_id: [row]

        ok, _reason = runner._dispatch(dict(ITEM), force=True)

        router.enqueue.assert_called_once()
        self.assertTrue(ok)

    def test_unreadable_task_falls_through_to_upsert(self):
        """Read failure leans toward dispatching: a lost comment is unrecoverable."""
        runner, router = _runner([])
        router.client.get_task_by_aone.side_effect = RuntimeError("control plane down")

        runner._dispatch(dict(ITEM), force=True)

        router.enqueue.assert_called_once()

    def test_absent_task_still_upserts(self):
        """First dispatch of a ticket has no live Task to compare against."""
        runner, router = _runner([])

        runner._dispatch(dict(ITEM), force=True)

        router.enqueue.assert_called_once()


class ScanNoopAccountingTest(unittest.TestCase):
    """A no-op must not be booked as a dispatch, but must still come to rest."""

    def _runner(self, dispatch_reason):
        runner = scan.ScanRunner.__new__(scan.ScanRunner)
        runner._prefetch_idle_claimed_a1 = mock.Mock()
        runner._done_watch_retry = set()
        runner._done_watch_confirm = set()
        runner._signal_query_retry = set()
        runner._dispatch_retry = set()
        runner._decide = lambda candidates: [{
            "id": ITEM["id"], "title": ITEM["title"], "item": dict(ITEM),
            "action": "dispatch", "reason": "new", "force": True,
            "conclusive": False, "dispatch_context": None,
        }]
        runner._dispatch = mock.Mock(
            return_value=scan.EnqueueResult(True, dispatch_reason))
        return runner

    def test_noop_is_not_reported_as_a_dispatch(self):
        runner = self._runner(scan.NOOP_DESIRED_REVISION)

        with mock.patch.object(scan.log, "info") as info:
            conclusive = runner._tick_auto([dict(ITEM)])

        logged = " ".join(str(call.args) for call in info.call_args_list)
        self.assertIn("no upsert", logged)
        self.assertNotIn("dispatched", logged)
        self.assertNotIn("persisted", logged)
        # Nothing left to do, so the decision rests and is not replayed next tick.
        self.assertEqual(conclusive, {ITEM["id"]: ITEM["modified"]})
        self.assertNotIn(ITEM["id"], runner._dispatch_retry)

    def test_real_dispatch_is_still_reported(self):
        runner = self._runner("task_persisted")

        with mock.patch.object(scan.log, "info") as info:
            conclusive = runner._tick_auto([dict(ITEM)])

        logged = " ".join(str(call.args) for call in info.call_args_list)
        self.assertIn("dispatched", logged)
        self.assertEqual(conclusive, {ITEM["id"]: ITEM["modified"]})


if __name__ == "__main__":
    unittest.main()
