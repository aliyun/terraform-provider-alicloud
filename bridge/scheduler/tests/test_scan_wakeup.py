"""Tests for ScanRunner._maybe_wake (方案 A: force-redispatch 唤醒 idle Task).

_maybe_wake is the scan-side integration of recovery_wakeup: on a force=True
dispatch (jarvis-idle/claimed + new human comment), if the control-plane Task is
RECOVERY_REQUIRED + retry-exhausted, force-redispatch it back to READY instead of
the inert desired-state upsert. Idempotent per epoch (retry/desired advance).
"""
import os
import unittest
from unittest import mock

from bridge import recovery_wakeup
from bridge.scheduler.runners import scan


def _task_row(status="RECOVERY_REQUIRED", retry=4, max_retries=3,
              task_id=471, desired="d1", gen=1):
    return {"id": task_id, "status": status, "retryCount": retry,
            "maxRetries": max_retries, "desiredRevision": desired,
            "generation": gen}


def _timeline(task_id=471, gen=1, retry=4, desired="d1", sess_id=1210, fence=4):
    return {
        "task": {"id": task_id, "status": "RECOVERY_REQUIRED", "generation": gen,
                 "stateVersion": 9, "retryCount": retry, "maxRetries": 3,
                 "currentSessionId": sess_id, "desiredRevision": desired,
                 "processingRevision": "p1"},
        "sessions": [{"id": sess_id, "status": "RESUMABLE", "fenceToken": fence,
                      "historicalWorkerKey": "w", "historicalWorkerId": 1,
                      "historicalWorkerProcessUuid": "p"}],
    }


class MaybeWakeTests(unittest.TestCase):
    def setUp(self):
        self._orig = recovery_wakeup.portable_replacement_for_redispatch
        recovery_wakeup.portable_replacement_for_redispatch = lambda *a, **k: None
        os.environ["JARVIS_RECOVERY_WAKEUP_ENABLE"] = "1"

    def tearDown(self):
        recovery_wakeup.portable_replacement_for_redispatch = self._orig
        os.environ.pop("JARVIS_RECOVERY_WAKEUP_ENABLE", None)

    def _runner(self, get_by_aone_rows, force_result, timeline=None):
        runner = scan.ScanRunner(
            logger=mock.Mock(), task_client=mock.Mock(), repo_root=mock.Mock())
        runner.execution_router.client.get_task_by_aone = mock.Mock(return_value=get_by_aone_rows)
        runner.execution_router.client.get_task_timeline = mock.Mock(
            return_value=timeline or _timeline())
        runner.execution_router.client.force_redispatch_task = mock.Mock(return_value=force_result)
        return runner

    def test_wakes_recovery_required_retry_exhausted(self):
        runner = self._runner([_task_row()], {"action": "REDISPATCHED",
                                              "task": {"generation": 2}})
        self.assertTrue(runner._maybe_wake("84103836"))
        self.assertIn("84103836", runner._woken_epochs)
        runner.execution_router.client.force_redispatch_task.assert_called_once()

    def test_skips_when_not_recovery_required(self):
        runner = self._runner([_task_row(status="RUNNING")], {"action": "REDISPATCHED"})
        self.assertFalse(runner._maybe_wake("84103836"))
        self.assertNotIn("84103836", runner._woken_epochs)
        runner.execution_router.client.force_redispatch_task.assert_not_called()

    def test_skips_when_retry_within_budget(self):
        runner = self._runner([_task_row(retry=2, max_retries=3)], {"action": "REDISPATCHED"})
        self.assertFalse(runner._maybe_wake("84103836"))
        runner.execution_router.client.force_redispatch_task.assert_not_called()

    def test_idempotent_same_epoch(self):
        runner = self._runner([_task_row()], {"action": "REDISPATCHED",
                                              "task": {"generation": 2}})
        self.assertTrue(runner._maybe_wake("84103836"))
        runner.execution_router.client.force_redispatch_task.reset_mock()
        # second call same epoch: returns True (skip enqueue) without a new force call
        self.assertTrue(runner._maybe_wake("84103836"))
        runner.execution_router.client.force_redispatch_task.assert_not_called()

    def test_new_epoch_after_retry_bump_wakes_again(self):
        runner = self._runner([_task_row(retry=4)], {"action": "REDISPATCHED",
                                                     "task": {"generation": 2}})
        self.assertTrue(runner._maybe_wake("84103836"))
        runner.execution_router.client.get_task_by_aone = mock.Mock(
            return_value=[_task_row(retry=5)])
        runner.execution_router.client.get_task_timeline = mock.Mock(
            return_value=_timeline(retry=5))
        runner.execution_router.client.force_redispatch_task.reset_mock()
        self.assertTrue(runner._maybe_wake("84103836"))
        runner.execution_router.client.force_redispatch_task.assert_called_once()

    def test_non_redispatch_action_does_not_mark_woken(self):
        runner = self._runner([_task_row()], {"action": "BLOCKED"})
        self.assertFalse(runner._maybe_wake("84103836"))
        self.assertNotIn("84103836", runner._woken_epochs)

    def test_get_task_by_aone_failure_returns_false(self):
        runner = self._runner([_task_row()], {"action": "REDISPATCHED"})
        runner.execution_router.client.get_task_by_aone = mock.Mock(
            side_effect=RuntimeError("net"))
        self.assertFalse(runner._maybe_wake("84103836"))

    def test_empty_task_rows_returns_false(self):
        runner = self._runner([], {"action": "REDISPATCHED"})
        self.assertFalse(runner._maybe_wake("84103836"))

    def test_dict_response_handled(self):
        # get_task_by_aone may return a single dict instead of a list
        runner = self._runner(_task_row(), {"action": "REDISPATCHED",
                                            "task": {"generation": 2}})
        self.assertTrue(runner._maybe_wake("84103836"))


if __name__ == "__main__":
    unittest.main()
