from __future__ import annotations

from datetime import datetime, timedelta, timezone
import unittest
from zoneinfo import ZoneInfo

from bridge.scheduler import (
    AdaptiveSchedule, CapabilityValidationContext, DailySchedule, HandlerRunner,
    IntervalSchedule, JobResult, JobResultStatus, MisfirePolicy,
    ScheduledJobDefinition, TriggerPlanner, definition_digest,
    definition_snapshot, validate_registry,
)
import bridge.scheduler.registry as registry


UTC = timezone.utc
SHANGHAI = ZoneInfo("Asia/Shanghai")


def at(hour: int, minute: int = 0, second: int = 0, tz=UTC) -> datetime:
    return datetime(2026, 7, 20, hour, minute, second, tzinfo=tz)


def definition(*, revision: int = 1, schedule=None) -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        "aone.scan", revision, "scan",
        schedule or IntervalSchedule(60, True), HandlerRunner("aone.scan"),
        MisfirePolicy.COALESCE, 5,
    )


class SchedulerCoreTests(unittest.TestCase):
    def test_registry_is_yaml_loaded_with_all_migrated_jobs(self):
        loaded = registry.load_jobs()
        self.assertEqual(tuple(item.id for item in loaded),
                         ("daily.probe", "aone.scan", "aone.claim-health",
                          "daily.nudge", "aone.reply", "pr.watch",
                          "external.recovery"))
        self.assertIs(loaded, registry.JOBS)

    def test_registry_materializes_generator_once_and_returns_tuple(self):
        seen = []
        item = definition()
        def source():
            seen.append("once")
            yield item
        value = validate_registry(source(), context=CapabilityValidationContext(handler_keys={"aone.scan"}))
        self.assertEqual(seen, ["once"])
        self.assertEqual(value, (item,))

    def test_registry_rejects_digest_revision_drift(self):
        item = definition(revision=2)
        with self.assertRaisesRegex(ValueError, "digest mismatch"):
            validate_registry((item,), context=CapabilityValidationContext(handler_keys={"aone.scan"}), expected_digests={"aone.scan": "0" * 64})

    def test_failure_result_cannot_commit_success_state(self):
        with self.assertRaisesRegex(ValueError, "failure results"):
            JobResult(JobResultStatus.PERMANENT_FAILURE, next_due_at=at(1), error="bad protocol")

    def test_slot_identity_includes_revision_and_utc_instant(self):
        planner = TriggerPlanner()
        first = planner.slot_key(definition(revision=1), at(2))
        changed = planner.slot_key(definition(revision=2), at(2))
        same_in_shanghai = planner.slot_key(definition(revision=1), at(10, tz=SHANGHAI))
        self.assertNotEqual(first, changed)
        self.assertEqual(first, same_in_shanghai)
        self.assertIn("@r1@", first)

    def test_interval_uses_absolute_series_not_completion_delay(self):
        planner = TriggerPlanner()
        due = at(1)
        self.assertEqual(planner.next_due(definition(), slot_due_at=due, completed_at=due + timedelta(seconds=150), result=JobResult(JobResultStatus.SUCCEEDED)), due + timedelta(seconds=180))

    def test_daily_misfire_replays_current_day_then_advances_one_day(self):
        planner = TriggerPlanner()
        daily = ScheduledJobDefinition(
            "daily.scan", 1, "daily", DailySchedule(10, 0),
            HandlerRunner("daily.scan"), MisfirePolicy.CURRENT_DAY, 300,
        )
        self.assertEqual(planner.initial_due(daily, at(3)), at(2))
        self.assertEqual(planner.next_due(daily, slot_due_at=at(2), completed_at=at(3), result=JobResult(JobResultStatus.SUCCEEDED)), at(2) + timedelta(days=1))

    def test_daily_retry_never_crosses_the_next_plan_boundary(self):
        planner = TriggerPlanner()
        daily = ScheduledJobDefinition(
            "daily.scan", 1, "daily", DailySchedule(10, 0),
            HandlerRunner("daily.scan"), MisfirePolicy.CURRENT_DAY, 300,
        )
        slot = datetime(2026, 7, 23, 10, 0, tzinfo=SHANGHAI)
        self.assertEqual(
            planner.retry_due(
                daily,
                slot_due_at=slot,
                failed_at=datetime(2026, 7, 24, 9, 54, tzinfo=SHANGHAI),
            ),
            datetime(2026, 7, 24, 9, 59, tzinfo=SHANGHAI),
        )
        self.assertIsNone(planner.retry_due(
            daily,
            slot_due_at=slot,
            failed_at=datetime(2026, 7, 24, 9, 55, tzinfo=SHANGHAI),
        ))

    def test_schedule_and_misfire_combinations_fail_closed(self):
        with self.assertRaisesRegex(ValueError, "does not support"):
            ScheduledJobDefinition(
                "daily.scan", 1, "daily", DailySchedule(10, 0),
                HandlerRunner("daily.scan"), MisfirePolicy.COALESCE, 300,
            )
        with self.assertRaisesRegex(ValueError, "does not support"):
            ScheduledJobDefinition(
                "aone.scan", 1, "interval", IntervalSchedule(60, True),
                HandlerRunner("aone.scan"), MisfirePolicy.CURRENT_DAY, 5,
            )

    def test_adaptive_out_of_bounds_is_protocol_error_not_clamped(self):
        planner = TriggerPlanner()
        adaptive = definition(schedule=AdaptiveSchedule(10, 60, 300, True))
        completed = at(1)
        with self.assertRaisesRegex(ValueError, "outside"):
            planner.next_due(adaptive, slot_due_at=completed, completed_at=completed, result=JobResult(JobResultStatus.SUCCEEDED, next_due_at=completed + timedelta(seconds=301)))
        self.assertEqual(planner.next_due(adaptive, slot_due_at=completed, completed_at=completed, result=JobResult(JobResultStatus.SUCCEEDED)), completed + timedelta(seconds=60))

    def test_failure_does_not_generate_a_new_slot(self):
        planner = TriggerPlanner()
        with self.assertRaisesRegex(ValueError, "successful"):
            planner.next_due(definition(), slot_due_at=at(1), completed_at=at(1), result=JobResult(JobResultStatus.RETRYABLE_FAILURE, error="network"))


if __name__ == "__main__":
    unittest.main()
