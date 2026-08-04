"""Guards that decide when a source item is permanently absent.

The point of these tests is that quarantining is *hard* to trigger: stamping a
terminal status on a live work item would hide real work, so every guard is
tested in isolation for its ability to veto.
"""

import tempfile
import unittest
from pathlib import Path

from bridge.scheduler.runners.source_poison import (
    SOURCE_NOT_FOUND, SourcePoisonLedger,
)

DAY = 86400.0


class SourceNotFoundSentinelTest(unittest.TestCase):
    def test_sentinel_is_falsy_so_unhandled_paths_behave_as_before(self):
        # Any existing branch testing the resolved status for truthiness keeps
        # today's behaviour if it has not been taught about the sentinel.
        self.assertFalse(SOURCE_NOT_FOUND)


class QuarantineGuardTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.path = Path(self.tmp.name) / "source-poison-health.json"
        self.addCleanup(self.tmp.cleanup)

    def _ledger(self, now):
        return SourcePoisonLedger(self.path, clock=lambda: now)

    def _twice_over_a_day(self, aone_id):
        first = self._ledger(1000.0)
        first.record("4451", aone_id)
        first.save()
        later = self._ledger(1000.0 + 2 * DAY)
        return later, later.record("4451", aone_id)

    def test_all_guards_satisfied_quarantines(self):
        led, ep = self._twice_over_a_day("779")
        self.assertTrue(led.should_quarantine(ep, "779", project_read_ok=True))

    def test_guard1_plausible_eight_digit_id_is_never_quarantined(self):
        led, ep = self._twice_over_a_day("84574563")
        self.assertFalse(
            led.should_quarantine(ep, "84574563", project_read_ok=True))

    def test_guard2_single_observation_is_not_enough(self):
        led = self._ledger(1000.0)
        ep = led.record("4451", "779")
        self.assertFalse(led.should_quarantine(ep, "779", project_read_ok=True))

    def test_guard2_two_observations_inside_24h_is_not_enough(self):
        first = self._ledger(1000.0)
        first.record("4451", "779")
        first.save()
        led = self._ledger(1000.0 + 3600.0)
        ep = led.record("4451", "779")
        self.assertFalse(led.should_quarantine(ep, "779", project_read_ok=True))

    def test_guard3_outage_blocks_quarantine(self):
        led, ep = self._twice_over_a_day("779")
        self.assertFalse(
            led.should_quarantine(ep, "779", project_read_ok=False))

    def test_non_numeric_id_is_not_quarantined(self):
        led, ep = self._twice_over_a_day("not-an-id")
        self.assertFalse(
            led.should_quarantine(ep, "not-an-id", project_read_ok=True))

    def test_mark_alerted_is_recorded_and_persisted(self):
        led, _ = self._twice_over_a_day("779")
        led.mark_alerted("4451")
        led.save()
        self.assertTrue(
            self._ledger(9999.0).episode("4451").get("lastAlertAt"))

    def test_alerted_is_false_before_marking(self):
        led, _ = self._twice_over_a_day("779")
        self.assertFalse(led.alerted("4451"))

    def test_forget_drops_the_episode_so_observations_cannot_accrue_forever(self):
        led, _ = self._twice_over_a_day("779")
        led.forget("4451")
        led.save()
        reloaded = self._ledger(9999.0)
        self.assertEqual(reloaded.episode("4451"), {})
        # A single fresh observation must not immediately re-qualify.
        ep = reloaded.record("4451", "779")
        self.assertFalse(
            reloaded.should_quarantine(ep, "779", project_read_ok=True))

    def test_corrupt_state_file_is_tolerated(self):
        self.path.write_text("{not json")
        led = self._ledger(1000.0)
        self.assertEqual(led.episode("4451"), {})


class ClassifyPointReadTest(unittest.TestCase):
    def test_404_maps_to_sentinel(self):
        from bridge.scheduler.runners.scan import ScanRunner
        self.assertIs(
            ScanRunner._classify_point_read_failure(
                "Error: workitem get failed (404): 工作项不存在"),
            SOURCE_NOT_FOUND)

    def test_403_is_not_permanent_absence(self):
        from bridge.scheduler.runners.scan import ScanRunner
        self.assertIsNone(
            ScanRunner._classify_point_read_failure(
                "Error: workitem get failed (403): no read permission"))

    def test_timeout_is_not_permanent_absence(self):
        from bridge.scheduler.runners.scan import ScanRunner
        self.assertIsNone(
            ScanRunner._classify_point_read_failure("timed out"))

    def test_empty_stderr_is_not_permanent_absence(self):
        from bridge.scheduler.runners.scan import ScanRunner
        self.assertIsNone(ScanRunner._classify_point_read_failure(""))


class _FakeClient:
    def __init__(self, fail=False):
        self.calls = []
        self.fail = fail

    def update_source_status(self, task_id, aone_id, source_status, *,
                             request_id=None):
        if self.fail:
            raise RuntimeError("control plane rejected")
        self.calls.append((task_id, aone_id, source_status))
        return {}


class QuarantineWiringTest(unittest.TestCase):
    """End-to-end behaviour of the scan runner's quarantine step."""

    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.path = Path(self.tmp.name) / "source-poison-health.json"
        self.addCleanup(self.tmp.cleanup)
        self.sent = []

        from bridge.scheduler.runners import scan as scan_mod
        self.scan_mod = scan_mod
        original = scan_mod._dingtalk_event_enqueue
        self.addCleanup(
            lambda: setattr(scan_mod, "_dingtalk_event_enqueue", original))

        def fake_enqueue(ticket, project, event_key, staff_id, title, text,
                         allow_non_tf=False):
            self.sent.append(event_key)
            return True

        scan_mod._dingtalk_event_enqueue = fake_enqueue

    def _ledger(self, now):
        return SourcePoisonLedger(self.path, clock=lambda: now)

    def _runner(self):
        return self.scan_mod.ScanRunner.__new__(self.scan_mod.ScanRunner)

    _TASK = {"taskId": "4451", "aoneId": "779", "sourceProjectKey": "1086837"}

    def _seed_one_prior_observation(self):
        led = self._ledger(1000.0)
        led.record("4451", "779")
        led.save()

    def test_quarantines_and_notifies_once_when_guards_pass(self):
        self._seed_one_prior_observation()
        client = _FakeClient()
        led = self._ledger(1000.0 + 2 * DAY)
        count = self._runner()._quarantine_absent_sources(
            client, [dict(self._TASK)], led, {"1086837"})
        self.assertEqual(count, 1)
        self.assertEqual(client.calls, [("4451", "779", "Invalid")])
        self.assertEqual(self.sent, ["source-poison:not-found:4451"])

        # A later pass must neither re-stamp nor re-notify.
        client2 = _FakeClient()
        self.sent.clear()
        led2 = self._ledger(1000.0 + 3 * DAY)
        led2._episodes["4451"]["lastAlertAt"] = 1.0
        self._runner()._quarantine_absent_sources(
            client2, [dict(self._TASK)], led2, {"1086837"})
        self.assertEqual(self.sent, [])

    def test_outage_blocks_stamp_and_notification(self):
        self._seed_one_prior_observation()
        client = _FakeClient()
        led = self._ledger(1000.0 + 2 * DAY)
        count = self._runner()._quarantine_absent_sources(
            client, [dict(self._TASK)], led, set())
        self.assertEqual(count, 0)
        self.assertEqual(client.calls, [])
        self.assertEqual(self.sent, [])

    def test_first_sighting_records_evidence_but_does_not_stamp(self):
        client = _FakeClient()
        led = self._ledger(1000.0)
        count = self._runner()._quarantine_absent_sources(
            client, [dict(self._TASK)], led, {"1086837"})
        self.assertEqual(count, 0)
        self.assertEqual(client.calls, [])
        # Evidence must survive for the next pass to build on.
        self.assertEqual(self._ledger(2000.0).episode("4451").get("count"), 1)

    def test_stamp_failure_does_not_mark_as_alerted(self):
        self._seed_one_prior_observation()
        led = self._ledger(1000.0 + 2 * DAY)
        count = self._runner()._quarantine_absent_sources(
            _FakeClient(fail=True), [dict(self._TASK)], led, {"1086837"})
        self.assertEqual(count, 0)
        self.assertEqual(self.sent, [])
        self.assertFalse(self._ledger(3000.0).alerted("4451"))

    def test_plausible_id_is_left_alone_entirely(self):
        task = {"taskId": "2091", "aoneId": "84574563",
                "sourceProjectKey": "709564"}
        led = self._ledger(1000.0)
        led.record("2091", "84574563")
        led.save()
        client = _FakeClient()
        count = self._runner()._quarantine_absent_sources(
            client, [task], self._ledger(1000.0 + 2 * DAY), {"709564"})
        self.assertEqual(count, 0)
        self.assertEqual(client.calls, [])


if __name__ == "__main__":
    unittest.main()
