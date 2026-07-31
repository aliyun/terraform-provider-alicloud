"""Tests for bridge.recovery_wakeup (should_wake + CAS force-redispatch + epoch).

Covers the shared helper used by scan (方案 A) and claim_health (方案 B):
- should_wake: RECOVERY_REQUIRED + retryCount > maxRetries only.
- wake_redispatch: get_task_timeline → CAS snapshot → portable_replacement →
  force_redispatch_task(target PERSISTENT, auto-target).
- epoch_key: changes on retry/desired advance (caller ledger dedup).
- enabled: JARVIS_RECOVERY_WAKEUP_ENABLE default off.
"""
import os
import unittest
from unittest import mock

from bridge import recovery_wakeup
from bridge.recovery_wakeup import (
    enabled, epoch_key, should_wake, wake_redispatch,
)


def _timeline(retry=4, max_retries=3, status="RECOVERY_REQUIRED",
              gen=1, desired="d1", sess_id=1210, fence=4, task_id=471):
    return {
        "task": {
            "id": task_id, "status": status, "generation": gen,
            "stateVersion": 9, "retryCount": retry, "maxRetries": max_retries,
            "currentSessionId": sess_id, "desiredRevision": desired,
            "processingRevision": "p1",
        },
        "sessions": [{
            "id": sess_id, "status": "RESUMABLE", "fenceToken": fence,
            "historicalWorkerKey": "worker-x",
            "historicalWorkerId": 1,
            "historicalWorkerProcessUuid": "p1",
        }],
    }


class FakeClient:
    def __init__(self, timeline, force_result):
        self._timeline = timeline
        self.force_result = force_result
        self.force_calls = []

    def get_task_timeline(self, task_id):
        return self._timeline

    def force_redispatch_task(self, task_id, **kwargs):
        self.force_calls.append({"task_id": task_id, **kwargs})
        return self.force_result


class ShouldWakeTests(unittest.TestCase):
    def test_retry_exhausted(self):
        self.assertTrue(should_wake(
            {"status": "RECOVERY_REQUIRED", "retryCount": 4, "maxRetries": 3}))

    def test_retry_within_budget(self):
        self.assertFalse(should_wake(
            {"status": "RECOVERY_REQUIRED", "retryCount": 2, "maxRetries": 3}))

    def test_equal_retry_not_exhausted(self):
        self.assertFalse(should_wake(
            {"status": "RECOVERY_REQUIRED", "retryCount": 3, "maxRetries": 3}))

    def test_other_status(self):
        self.assertFalse(should_wake(
            {"status": "RUNNING", "retryCount": 4, "maxRetries": 3}))

    def test_missing_fields(self):
        self.assertFalse(should_wake({"status": "RECOVERY_REQUIRED"}))

    def test_non_mapping(self):
        self.assertFalse(should_wake(None))

    def test_non_numeric_retry(self):
        self.assertFalse(should_wake(
            {"status": "RECOVERY_REQUIRED", "retryCount": "x", "maxRetries": 3}))


class EpochKeyTests(unittest.TestCase):
    def test_changes_on_retry(self):
        base = {"retryCount": 4, "desiredRevision": "d1"}
        self.assertNotEqual(
            epoch_key(base, aone_id="111"),
            epoch_key({**base, "retryCount": 5}, aone_id="111"))

    def test_changes_on_desired(self):
        base = {"retryCount": 4, "desiredRevision": "d1"}
        self.assertNotEqual(
            epoch_key(base, aone_id="111"),
            epoch_key({**base, "desiredRevision": "d2"}, aone_id="111"))

    def test_stable_when_unchanged(self):
        base = {"retryCount": 4, "desiredRevision": "d1"}
        self.assertEqual(
            epoch_key(base, aone_id="111"),
            epoch_key(base, aone_id="111"))


class WakeRedispatchTests(unittest.TestCase):
    def setUp(self):
        # portable_replacement_for_redispatch is heavy (input contract) and not
        # under test here; stub it to None (portable input, no rebuild).
        self._orig = recovery_wakeup.portable_replacement_for_redispatch
        recovery_wakeup.portable_replacement_for_redispatch = lambda *a, **k: None

    def tearDown(self):
        recovery_wakeup.portable_replacement_for_redispatch = self._orig

    def test_calls_force_with_cas(self):
        client = FakeClient(
            _timeline(),
            {"action": "REDISPATCHED",
             "task": {"generation": 2, "status": "READY"},
             "generation": 2})
        result = wake_redispatch(client, "471", "84103836", "test reason")
        self.assertEqual(result["action"], "REDISPATCHED")
        call = client.force_calls[0]
        self.assertEqual(call["task_id"], "471")
        self.assertEqual(call["expected_task_status"], "RECOVERY_REQUIRED")
        self.assertEqual(call["expected_generation"], 1)
        self.assertEqual(call["expected_state_version"], 9)
        self.assertEqual(call["expected_retry_count"], 4)
        self.assertEqual(call["expected_session_id"], 1210)
        self.assertEqual(call["expected_session_status"], "RESUMABLE")
        self.assertEqual(call["expected_fence_token"], 4)
        self.assertEqual(call["expected_worker_key"], "worker-x")
        self.assertIsNone(call["target_worker_key"])
        self.assertIsNone(call["target_host_id"])
        self.assertEqual(call["target_runtime"], "PERSISTENT")
        self.assertEqual(call["reason"], "test reason")
        self.assertIsNone(call["portable_input_replacement"])

    def test_raises_on_missing_cas(self):
        client = FakeClient(
            {"task": {"id": 471, "status": "RECOVERY_REQUIRED"}, "sessions": []}, {})
        with self.assertRaises(ValueError):
            wake_redispatch(client, "471", "111", "x")

    def test_raises_on_wrong_task_id(self):
        tl = {"task": {"id": 999, "status": "RECOVERY_REQUIRED", "generation": 1,
                       "stateVersion": 1, "retryCount": 4, "maxRetries": 3,
                       "currentSessionId": 1},
              "sessions": [{"id": 1, "status": "RESUMABLE", "fenceToken": 1}]}
        client = FakeClient(tl, {})
        with self.assertRaises(ValueError):
            wake_redispatch(client, "471", "111", "x")

    def test_propagates_force_error(self):
        client = FakeClient(_timeline(), None)
        client.force_redispatch_task = mock.Mock(side_effect=RuntimeError("boom"))
        with self.assertRaises(RuntimeError):
            wake_redispatch(client, "471", "111", "x")


class EnabledTests(unittest.TestCase):
    def test_default_off(self):
        with mock.patch.dict(os.environ, {}, clear=False):
            os.environ.pop("JARVIS_RECOVERY_WAKEUP_ENABLE", None)
            self.assertFalse(enabled())

    def test_on(self):
        with mock.patch.dict(os.environ, {"JARVIS_RECOVERY_WAKEUP_ENABLE": "1"}):
            self.assertTrue(enabled())

    def test_off_values(self):
        for value in ("0", "", "false", "False", "no"):
            with mock.patch.dict(os.environ, {"JARVIS_RECOVERY_WAKEUP_ENABLE": value}):
                self.assertFalse(enabled())


if __name__ == "__main__":
    unittest.main()
