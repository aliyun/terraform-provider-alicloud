from __future__ import annotations

from datetime import datetime, timedelta, timezone
import unittest
from zoneinfo import ZoneInfo

from bridge.scheduler import (
    AdaptiveSchedule, CapabilityValidationContext, CheckpointUpgradePolicy, DailySchedule, HandlerRunner,
    IntervalSchedule, JobPurpose, JobResult, JobResultStatus, MisfirePolicy,
    ReplayPolicy, ScheduledJobDefinition, TaskObservation, TriggerPlanner, definition_digest,
    definition_snapshot, validate_registry,
)
import bridge.scheduler.jobs as jobs


UTC = timezone.utc
SHANGHAI = ZoneInfo("Asia/Shanghai")


def at(hour: int, minute: int = 0, second: int = 0, tz=UTC) -> datetime:
    return datetime(2026, 7, 20, hour, minute, second, tzinfo=tz)


def definition(*, revision: int = 1, schedule=None) -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        "aone.scan", revision, "scan", JobPurpose.DISCOVERY,
        schedule or IntervalSchedule(60, True), HandlerRunner("aone.scan"),
        MisfirePolicy.COALESCE, 30, 5, ReplayPolicy.TASK_UPSERT_IDEMPOTENT,
    )


class SchedulerCoreTests(unittest.TestCase):
    def test_registry_is_yaml_loaded_with_all_smoke_schedule_types(self):
        loaded = jobs.load_jobs()
        self.assertEqual(tuple(item.id for item in loaded),
                         ("smoke.interval", "smoke.daily", "smoke.adaptive"))
        self.assertIs(loaded, jobs.JOBS)

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

    def test_checkpoint_upgrade_contract_is_in_definition_digest(self):
        migrated = ScheduledJobDefinition(
            "aone.scan", 2, "scan", JobPurpose.DISCOVERY, IntervalSchedule(60, True),
            HandlerRunner("aone.scan"), MisfirePolicy.COALESCE, 30, 5,
            ReplayPolicy.TASK_UPSERT_IDEMPOTENT,
            checkpoint_upgrade=CheckpointUpgradePolicy.MIGRATE,
            checkpoint_upgrader_key="aone.scan.v1_to_v2",
        )
        snapshot = definition_snapshot(migrated)
        self.assertEqual(snapshot["checkpoint_upgrade"], "MIGRATE")
        self.assertEqual(snapshot["checkpoint_upgrader_key"], "aone.scan.v1_to_v2")
        self.assertNotEqual(definition_digest(migrated), definition_digest(definition(revision=2)))
        with self.assertRaisesRegex(ValueError, "checkpoint_upgrader_key"):
            ScheduledJobDefinition("aone.scan", 2, "scan", JobPurpose.DISCOVERY, IntervalSchedule(60, True), HandlerRunner("aone.scan"), MisfirePolicy.COALESCE, 30, 5, ReplayPolicy.TASK_UPSERT_IDEMPOTENT, checkpoint_upgrade=CheckpointUpgradePolicy.MIGRATE)
        overlap = ScheduledJobDefinition(
            "aone.scan", 2, "scan", JobPurpose.DISCOVERY, IntervalSchedule(60, True),
            HandlerRunner("aone.scan"), MisfirePolicy.COALESCE, 30, 5,
            ReplayPolicy.TASK_UPSERT_IDEMPOTENT,
            checkpoint_upgrade=CheckpointUpgradePolicy.RESET_OVERLAP,
            checkpoint_upgrader_key="aone.scan.overlap.v2",
        )
        self.assertEqual(definition_snapshot(overlap)["checkpoint_upgrader_key"],
                         "aone.scan.overlap.v2")
        with self.assertRaisesRegex(ValueError, "RESET_OVERLAP or MIGRATE"):
            ScheduledJobDefinition("aone.scan", 2, "scan", JobPurpose.DISCOVERY, IntervalSchedule(60, True), HandlerRunner("aone.scan"), MisfirePolicy.COALESCE, 30, 5, ReplayPolicy.TASK_UPSERT_IDEMPOTENT, checkpoint_upgrader_key="wrong")
        with self.assertRaisesRegex(ValueError, "unregistered checkpoint upgrader"):
            validate_registry(
                (overlap,),
                context=CapabilityValidationContext(handler_keys={"aone.scan"}),
            )
        self.assertEqual(
            validate_registry(
                (overlap,),
                context=CapabilityValidationContext(
                    handler_keys={"aone.scan"},
                    checkpoint_upgrader_keys={"aone.scan.overlap.v2"},
                ),
            ),
            (overlap,),
        )

    def test_task_observation_keeps_v14_optional_fields_immutable(self):
        observation = TaskObservation(
            "aone:2100304:84352841", "AONE", {"aoneId": "84352841"}, "ticket",
            "updated:1", ["SCAN", "TRACKER"], {"prompt": ["inspect"]},
            recovery_policy="RESUME_ONLY", persona="terraform-rd", priority={"rank": 2},
            aone_id="84352841", comment_cursor={"id": 3},
            required_capabilities={"kinds": ["terraform"]}, max_retries=2,
            source_status="开发中",
        )
        self.assertEqual(observation.trigger_mask, ("SCAN", "TRACKER"))
        self.assertIsInstance(observation.payload["prompt"], tuple)
        self.assertIsInstance(observation.required_capabilities["kinds"], tuple)
        self.assertEqual(observation.source_status, "开发中")
        with self.assertRaises(TypeError):
            observation.payload["prompt"] = ()
        for kwargs in (
            {"trigger_mask": "SCAN"}, {"max_retries": -1}, {"max_retries": True},
            {"persona": " "}, {"aone_id": " "}, {"source_status": "x" * 65},
        ):
            values = dict(task_key="key", source_type="AONE", source_ref={}, task_type="ticket",
                          desired_revision="1", trigger_mask=["SCAN"], payload={})
            values.update(kwargs)
            with self.subTest(kwargs=kwargs), self.assertRaises((TypeError, ValueError)):
                TaskObservation(**values)

    def test_failure_result_cannot_commit_success_state(self):
        with self.assertRaisesRegex(ValueError, "failure results"):
            JobResult(JobResultStatus.RETRYABLE_FAILURE, checkpoint={"cursor": 1}, error="temporary")
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
        daily = definition(schedule=DailySchedule(10, 0))
        self.assertEqual(planner.initial_due(daily, at(3)), at(2))
        self.assertEqual(planner.next_due(daily, slot_due_at=at(2), completed_at=at(3), result=JobResult(JobResultStatus.SUCCEEDED)), at(2) + timedelta(days=1))

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
