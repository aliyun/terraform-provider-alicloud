from __future__ import annotations

from datetime import datetime, timedelta, timezone
from pathlib import Path
import tempfile
import unittest
from zoneinfo import ZoneInfo

from bridge.scheduler.scheduler import require_scheduler_role
from bridge.scheduler.engine import TriggerPlanner
from bridge.scheduler.model import (
    AdaptiveSchedule, DailySchedule, HandlerRunner, IntervalSchedule, JobResult,
    JobResultStatus, MisfirePolicy, ScheduledJobDefinition, definition_snapshot,
)
from bridge.scheduler.runners import RUNNER_KEYS as IMPLEMENTED_RUNNER_KEYS
import bridge.scheduler.jobs as jobs


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
    def test_registry_contains_all_migrated_jobs(self):
        self.assertEqual(tuple(item.id for item in jobs.JOBS),
                         ("daily.probe", "aone.scan", "aone.claim-health",
                          "task.owner-health", "daily.nudge",
                          "aone.weekly-comment-participation",
                          "aone.reply", "pr.watch", "external.recovery"))
        self.assertEqual(jobs.RUNNER_KEYS, IMPLEMENTED_RUNNER_KEYS)
        self.assertEqual(
            ("daily_probe", "scan", "claim_health", "owner_health", "daily_nudge",
             "weekly_comment_participation", "reply", "pr_watch", "recovery"),
            tuple(item.runner.handler_key for item in jobs.JOBS),
        )
        self.assertTrue(all(isinstance(item.runner, HandlerRunner)
                            for item in jobs.JOBS))
        self.assertTrue(jobs.JOBS[0].enabled)
        self.assertEqual(jobs.JOBS[0].revision, 3)
        self.assertEqual(set(definition_snapshot(jobs.JOBS[0])), {
            "id", "revision", "description", "schedule", "runner", "misfire",
            "retry_delay_seconds", "enabled",
        })

    def test_registry_is_loaded_from_yaml(self):
        content = """\
version: 1
jobs:
  - key: test.interval
    runner: test_interval
    revision: 3
    description: loaded from yaml
    schedule: {kind: interval, interval_seconds: 17, run_immediately: false}
    misfire: COALESCE
    retry_delay_seconds: 4
    enabled: true
"""
        with tempfile.TemporaryDirectory() as temp_dir:
            path = Path(temp_dir) / "jobs.yaml"
            path.write_text(content, encoding="utf-8")
            loaded = jobs.load_jobs(path)
        self.assertEqual(len(loaded), 1)
        self.assertEqual(loaded[0].id, "test.interval")
        self.assertEqual(loaded[0].runner.handler_key, "test_interval")
        self.assertEqual(loaded[0].schedule, IntervalSchedule(17, False))

    def test_registry_rejects_unknown_and_missing_fields(self):
        cases = {
            "unknown": """\
version: 1
jobs: []
extra: true
""",
            "missing": """\
version: 1
jobs:
  - key: test.interval
    runner: test_interval
    revision: 1
    description: missing enabled
    schedule: {kind: interval, interval_seconds: 17, run_immediately: false}
    misfire: COALESCE
    retry_delay_seconds: 4
""",
            "schedule": """\
version: 1
jobs:
  - key: test.interval
    runner: test_interval
    revision: 1
    description: unknown schedule field
    schedule: {kind: interval, interval_seconds: 17, run_immediately: false, jitter: 1}
    misfire: COALESCE
    retry_delay_seconds: 4
    enabled: true
""",
        }
        with tempfile.TemporaryDirectory() as temp_dir:
            for name, content in cases.items():
                with self.subTest(name=name):
                    path = Path(temp_dir) / ("%s.yaml" % name)
                    path.write_text(content, encoding="utf-8")
                    with self.assertRaises(jobs.SchedulerJobsError):
                        jobs.load_jobs(path)

    def test_registry_rejects_duplicate_yaml_and_job_keys(self):
        cases = {
            "mapping": """\
version: 1
version: 1
jobs: []
""",
            "job": """\
version: 1
jobs:
  - key: test.interval
    runner: test_interval
    revision: 1
    description: first
    schedule: {kind: interval, interval_seconds: 17, run_immediately: false}
    misfire: COALESCE
    retry_delay_seconds: 4
    enabled: true
  - key: test.interval
    runner: second_interval
    revision: 2
    description: second
    schedule: {kind: interval, interval_seconds: 19, run_immediately: true}
    misfire: COALESCE
    retry_delay_seconds: 5
    enabled: true
""",
        }
        with tempfile.TemporaryDirectory() as temp_dir:
            for name, content in cases.items():
                with self.subTest(name=name):
                    path = Path(temp_dir) / ("%s.yaml" % name)
                    path.write_text(content, encoding="utf-8")
                    with self.assertRaises(jobs.SchedulerJobsError):
                        jobs.load_jobs(path)

    def test_only_scheduler_role_can_start_periodic_jobs(self):
        require_scheduler_role({"JARVIS_BRIDGE_ROLE": "scheduler"})
        for role in ("worker", "unknown"):
            with self.subTest(role=role), self.assertRaisesRegex(
                    RuntimeError, "JARVIS_BRIDGE_ROLE=scheduler"):
                require_scheduler_role({"JARVIS_BRIDGE_ROLE": role})

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
