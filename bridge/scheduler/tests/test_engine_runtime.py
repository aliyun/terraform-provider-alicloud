from __future__ import annotations

from dataclasses import replace
from datetime import datetime, timedelta, timezone
import threading
import unittest

from bridge.scheduler.engine import (
    RetryableJobError, RunnerDispatcher, ScheduledJobState, ScheduledJobStatus,
    SchedulerEngine, TriggerPlanner, build_registrations, plan_due_slots,
)
from bridge.scheduler.model import (
    DailySchedule, HandlerRunner, IntervalSchedule, JobResult, JobResultStatus,
    MisfirePolicy, ScheduledJobDefinition,
)


UTC = timezone.utc
SHANGHAI = timezone(timedelta(hours=8))


def at(hour: int, minute: int = 0) -> datetime:
    return datetime(2026, 7, 21, hour, minute, tzinfo=UTC)


def definition(job_id: str = "aone.scan") -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        job_id, 1, "test job", IntervalSchedule(60, True),
        HandlerRunner(job_id), MisfirePolicy.COALESCE, 5,
    )


class FakeControlPlane:
    def __init__(self, states, *, start_result=True):
        self.states = tuple(states)
        self.start_result = start_result
        self.registrations = []
        self.starts = []
        self.completions = []
        self.failures = []
        self.recoveries = 0

    def register(self, registrations):
        self.registrations.append(tuple(registrations))
        return self.states

    def list_jobs(self):
        return self.states

    def recover_interrupted(self):
        self.recoveries += 1
        self.states = tuple(
            ScheduledJobState(
                state.job_key,
                ScheduledJobStatus.IDLE if state.status is ScheduledJobStatus.RUNNING else state.status,
                state.next_run_at,
            )
            for state in self.states
        )
        return tuple(state for state in self.states if state.status is ScheduledJobStatus.IDLE)

    def start(self, job_key, scheduled_for, next_run_at):
        self.starts.append((job_key, scheduled_for, next_run_at))
        return self.start_result

    def complete(self, job_key, scheduled_for, next_run_at):
        self.completions.append((job_key, scheduled_for, next_run_at))

    def fail(self, job_key, scheduled_for, *, retryable, next_run_at, error):
        self.failures.append((job_key, scheduled_for, retryable, next_run_at, error))


class FakeRunner:
    def __init__(self, result=None, error=None):
        self.result = result or JobResult(JobResultStatus.SUCCEEDED)
        self.error = error
        self.calls = []

    def run(self, item, scheduled_for):
        self.calls.append((item.id, scheduled_for))
        if self.error:
            raise self.error
        return self.result


class BlockingRunner(FakeRunner):
    def __init__(self, release: threading.Event, started: threading.Event):
        super().__init__()
        self.release = release
        self.started = started

    def run(self, item, scheduled_for):
        self.calls.append((item.id, scheduled_for))
        self.started.set()
        self.release.wait(5)
        return self.result


class SchedulerEngineTests(unittest.TestCase):
    def test_planning_is_pure_and_only_admits_due_idle_or_retryable_error_states(self):
        first = definition("aone.scan")
        second = definition("aone.reply")
        states = (
            ScheduledJobState("aone.reply", ScheduledJobStatus.ERROR, at(8)),
            ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, at(8, 30)),
            ScheduledJobState("unknown.job", ScheduledJobStatus.IDLE, at(7)),
            ScheduledJobState("aone.scan", ScheduledJobStatus.RUNNING, at(7)),
        )
        planned = plan_due_slots((first, second), states, now=at(9))
        self.assertEqual([(slot.definition.id, slot.scheduled_for) for slot in planned], [
            ("aone.reply", at(8)), ("aone.scan", at(8, 30)),
        ])

    def test_registration_plans_enabled_initial_due_without_control_plane_io(self):
        registrations = build_registrations(
            (replace(definition(), enabled=False),),
            planner=TriggerPlanner(), now=at(9),
        )
        self.assertEqual(len(registrations), 1)
        self.assertFalse(registrations[0].enabled)
        self.assertIsNone(registrations[0].next_run_at)
        self.assertEqual(registrations[0].definition["id"], "aone.scan")

    def test_stale_or_duplicate_control_plane_rejection_never_runs_or_reports_terminal_state(self):
        due = at(8)
        control = FakeControlPlane((ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),), start_result=False)
        runner = FakeRunner()
        engine = SchedulerEngine(
            (definition(),), control_plane=control,
            runtime=RunnerDispatcher({"aone.scan": runner}),
            clock=lambda: at(9))

        engine.tick(at(9))

        self.assertEqual(control.starts, [("aone.scan", due, at(9, 1))])
        self.assertEqual(runner.calls, [])
        self.assertEqual(control.completions, [])
        self.assertEqual(control.failures, [])

    def test_runner_failure_has_one_retryable_status_report(self):
        due = at(8)
        control = FakeControlPlane((ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),))
        runner = FakeRunner(error=RetryableJobError("temporary network issue"))
        engine = SchedulerEngine(
            (definition(),), control_plane=control,
            runtime=RunnerDispatcher({"aone.scan": runner}),
            clock=lambda: at(9))

        engine.tick(at(9))

        self.assertEqual(len(control.failures), 1)
        job_key, scheduled_for, retryable, retry_at, error = control.failures[0]
        self.assertEqual((job_key, scheduled_for, retryable, retry_at), ("aone.scan", due, True, at(9) + timedelta(seconds=5)))
        self.assertIn("RetryableJobError", error)
        self.assertEqual(control.completions, [])

    def test_daily_failure_at_next_boundary_becomes_permanent_without_renaming_slot(self):
        due = datetime(2026, 7, 23, 10, 0, tzinfo=SHANGHAI)
        failed_at = datetime(2026, 7, 24, 9, 55, tzinfo=SHANGHAI)
        daily = ScheduledJobDefinition(
            "daily.probe", 2, "probe", DailySchedule(10, 0),
            HandlerRunner("daily.probe"), MisfirePolicy.CURRENT_DAY, 300,
        )
        control = FakeControlPlane((
            ScheduledJobState("daily.probe", ScheduledJobStatus.IDLE, due),))
        runner = FakeRunner(error=RetryableJobError("temporary"))
        engine = SchedulerEngine(
            (daily,),
            control_plane=control,
            runtime=RunnerDispatcher({"daily.probe": runner}),
            clock=lambda: failed_at,
        )

        engine.tick(failed_at)

        job_key, scheduled_for, retryable, retry_at, error = control.failures[0]
        self.assertEqual(
            (job_key, scheduled_for, retryable, retry_at),
            ("daily.probe", due, False, None),
        )
        self.assertIn("retry window exhausted", error)

    def test_success_publishes_before_complete_and_advances_interval_from_slot_series(self):
        due = at(8)
        control = FakeControlPlane((ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),))
        runner = FakeRunner(JobResult(JobResultStatus.SUCCEEDED))
        times = iter((at(9, 1), at(9, 3)))
        engine = SchedulerEngine(
            (definition(),), control_plane=control,
            runtime=RunnerDispatcher({"aone.scan": runner}),
            clock=lambda: next(times))

        engine.tick(at(9, 1))

        self.assertEqual(control.starts, [("aone.scan", due, at(9, 2))])
        self.assertEqual(control.completions, [("aone.scan", due, at(9, 4))])

    def test_stop_closes_admission_and_does_not_start_a_new_scanner(self):
        due = at(8)
        control = FakeControlPlane((ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),))
        runner = FakeRunner()
        engine = SchedulerEngine(
            (definition(),), control_plane=control,
            runtime=RunnerDispatcher({"aone.scan": runner}),
            clock=lambda: at(9))

        engine.stop()

        self.assertFalse(engine.accepting)
        self.assertIsNone(engine.tick(at(9)))
        self.assertEqual(control.starts, [])
        self.assertEqual(runner.calls, [])

    def test_startup_recovery_is_explicit_and_preserves_the_reserved_successor(self):
        reserved = at(10)
        control = FakeControlPlane((
            ScheduledJobState("aone.scan", ScheduledJobStatus.RUNNING, reserved),))
        runner = FakeRunner()
        engine = SchedulerEngine(
            (definition(),), control_plane=control,
            runtime=RunnerDispatcher({"aone.scan": runner}),
            clock=lambda: at(9))

        self.assertEqual(engine.register(at(9))[0].status, ScheduledJobStatus.RUNNING)
        self.assertEqual(plan_due_slots((definition(),), control.list_jobs(), now=at(9)), ())
        self.assertIsNone(engine.tick(at(9)))
        self.assertEqual(control.recoveries, 0)

        recovered = engine.recover_interrupted()
        planned = plan_due_slots((definition(),), control.list_jobs(), now=at(9))
        engine.tick(at(9))

        self.assertEqual([(state.job_key, state.status, state.next_run_at) for state in recovered], [
            ("aone.scan", ScheduledJobStatus.IDLE, reserved),
        ])
        self.assertEqual(control.recoveries, 1)
        self.assertEqual(planned, ())
        self.assertEqual(runner.calls, [])

    def test_runner_definition_is_the_only_handler_routing_source(self):
        runner = FakeRunner()
        item = definition("aone.scan")
        runtime = RunnerDispatcher({"aone.scan": runner})

        result = runtime.run(item, at(9))

        self.assertEqual(result.status, JobResultStatus.SUCCEEDED)
        self.assertEqual(runner.calls, [("aone.scan", at(9))])

    def test_long_job_does_not_block_another_job_with_bounded_workers(self):
        release = threading.Event()
        first_started = threading.Event()
        second_done = threading.Event()
        first = BlockingRunner(release, first_started)

        class SignallingRunner(FakeRunner):
            def run(self, item, scheduled_for):
                value = super().run(item, scheduled_for)
                second_done.set()
                return value

        second = SignallingRunner()
        due = at(8)
        definitions = (definition("aone.scan"), definition("aone.reply"))
        control = FakeControlPlane((
            ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),
            ScheduledJobState("aone.reply", ScheduledJobStatus.IDLE, due),
        ))
        engine = SchedulerEngine(
            definitions,
            control_plane=control,
            runtime=RunnerDispatcher({
                "aone.scan": first,
                "aone.reply": second,
            }),
            clock=lambda: at(9),
            max_workers=2,
        )
        engine.tick(at(9))
        self.assertTrue(first_started.wait(1))
        self.assertTrue(second_done.wait(1))
        self.assertEqual(engine.active_count, 1)
        release.set()
        self.assertTrue(engine.wait_for_active(timeout=1))
        engine.shutdown()
        self.assertCountEqual(
            [key for key, _due, _next in control.completions],
            ["aone.scan", "aone.reply"],
        )

    def test_same_job_does_not_overlap_and_stop_waits_without_killing(self):
        release = threading.Event()
        started = threading.Event()
        runner = BlockingRunner(release, started)
        control = FakeControlPlane((
            ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, at(8)),
            ScheduledJobState("aone.scan", ScheduledJobStatus.ERROR, at(8, 1)),
        ))
        engine = SchedulerEngine(
            (definition("aone.scan"),),
            control_plane=control,
            runtime=RunnerDispatcher({"aone.scan": runner}),
            clock=lambda: at(9),
            max_workers=2,
        )
        engine.tick(at(9))
        self.assertTrue(started.wait(1))
        self.assertEqual(len(runner.calls), 1)
        self.assertEqual(len(control.starts), 1)
        engine.stop()
        self.assertFalse(engine.wait_for_active(timeout=0.01))
        self.assertEqual(len(runner.calls), 1)
        release.set()
        self.assertTrue(engine.wait_for_active(timeout=1))
        engine.shutdown()

    def test_worker_bound_does_not_pre_admit_unbounded_slots(self):
        releases = [threading.Event(), threading.Event()]
        starts = [threading.Event(), threading.Event()]
        runners = {
            "aone.scan": BlockingRunner(releases[0], starts[0]),
            "aone.reply": BlockingRunner(releases[1], starts[1]),
            "aone.third": FakeRunner(),
        }
        due = at(8)
        definitions = tuple(definition(key) for key in runners)
        control = FakeControlPlane(tuple(
            ScheduledJobState(key, ScheduledJobStatus.IDLE, due)
            for key in runners))
        engine = SchedulerEngine(
            definitions,
            control_plane=control,
            runtime=RunnerDispatcher(runners),
            clock=lambda: at(9),
            max_workers=2,
        )
        engine.tick(at(9))
        self.assertTrue(all(event.wait(1) for event in starts))
        self.assertEqual(engine.active_count, 2)
        self.assertEqual(len(control.starts), 2)
        self.assertEqual(runners["aone.third"].calls, [])
        for event in releases:
            event.set()
        self.assertTrue(engine.wait_for_active(timeout=1))
        engine.shutdown()


if __name__ == "__main__":
    unittest.main()
