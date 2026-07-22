from __future__ import annotations

from datetime import datetime, timedelta, timezone
import unittest

from bridge.scheduler import (
    ExecutionDisposition, HandlerRunner, IntervalSchedule, JobPurpose, JobResult,
    JobResultStatus, MisfirePolicy, ReplayPolicy, RetryableJobError,
    ScheduledJobDefinition, ScheduledJobState, ScheduledJobStatus, SchedulerEngine,
    ScannerRuntime, TriggerPlanner, build_registrations, plan_due_slots,
)


UTC = timezone.utc


def at(hour: int, minute: int = 0) -> datetime:
    return datetime(2026, 7, 21, hour, minute, tzinfo=UTC)


def definition(job_id: str = "aone.scan") -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        job_id, 1, "test job", JobPurpose.DISCOVERY, IntervalSchedule(60, True),
        HandlerRunner(job_id), MisfirePolicy.COALESCE, 30, 5,
        ReplayPolicy.TASK_UPSERT_IDEMPOTENT,
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


class FakePublisher:
    def __init__(self):
        self.calls = []

    def publish(self, item, result, scheduled_for):
        self.calls.append((item.id, result, scheduled_for))


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
            (definition(),), planner=TriggerPlanner(),
            now=at(9), is_enabled=lambda _: False,
        )
        self.assertEqual(len(registrations), 1)
        self.assertFalse(registrations[0].enabled)
        self.assertIsNone(registrations[0].next_run_at)
        self.assertEqual(registrations[0].definition["id"], "aone.scan")

    def test_stale_or_duplicate_control_plane_rejection_never_runs_or_reports_terminal_state(self):
        due = at(8)
        control = FakeControlPlane((ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),), start_result=False)
        runner = FakeRunner()
        publisher = FakePublisher()
        engine = SchedulerEngine(
            (definition(),), control_plane=control, runtime=ScannerRuntime(runner),
            publisher=publisher, clock=lambda: at(9))

        outcomes = engine.tick(at(9))

        self.assertEqual([item.disposition for item in outcomes], [ExecutionDisposition.START_REJECTED])
        self.assertEqual(control.starts, [("aone.scan", due, at(9, 1))])
        self.assertEqual(runner.calls, [])
        self.assertEqual(publisher.calls, [])
        self.assertEqual(control.completions, [])
        self.assertEqual(control.failures, [])

    def test_runner_failure_has_one_retryable_status_report(self):
        due = at(8)
        control = FakeControlPlane((ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),))
        runner = FakeRunner(error=RetryableJobError("temporary network issue"))
        engine = SchedulerEngine(
            (definition(),), control_plane=control, runtime=ScannerRuntime(runner),
            publisher=FakePublisher(), clock=lambda: at(9))

        outcomes = engine.tick(at(9))

        self.assertEqual([item.disposition for item in outcomes], [ExecutionDisposition.RETRYABLE_FAILURE])
        self.assertEqual(len(control.failures), 1)
        job_key, scheduled_for, retryable, retry_at, error = control.failures[0]
        self.assertEqual((job_key, scheduled_for, retryable, retry_at), ("aone.scan", due, True, at(9) + timedelta(seconds=5)))
        self.assertIn("RetryableJobError", error)
        self.assertEqual(control.completions, [])

    def test_success_publishes_before_complete_and_advances_interval_from_slot_series(self):
        due = at(8)
        control = FakeControlPlane((ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),))
        runner = FakeRunner(JobResult(JobResultStatus.SUCCEEDED))
        publisher = FakePublisher()
        times = iter((at(9, 1), at(9, 3)))
        engine = SchedulerEngine(
            (definition(),), control_plane=control, runtime=ScannerRuntime(runner),
            publisher=publisher, clock=lambda: next(times))

        outcomes = engine.tick(at(9, 1))

        self.assertEqual([item.disposition for item in outcomes], [ExecutionDisposition.COMPLETED])
        self.assertEqual(len(publisher.calls), 1)
        self.assertEqual(control.starts, [("aone.scan", due, at(9, 2))])
        self.assertEqual(control.completions, [("aone.scan", due, at(9, 4))])

    def test_stop_closes_admission_and_does_not_start_a_new_scanner(self):
        due = at(8)
        control = FakeControlPlane((ScheduledJobState("aone.scan", ScheduledJobStatus.IDLE, due),))
        runner = FakeRunner()
        engine = SchedulerEngine(
            (definition(),), control_plane=control, runtime=ScannerRuntime(runner),
            publisher=FakePublisher(), clock=lambda: at(9))

        engine.stop()

        self.assertFalse(engine.accepting)
        self.assertEqual(engine.tick(at(9)), ())
        self.assertEqual(control.starts, [])
        self.assertEqual(runner.calls, [])

    def test_startup_recovery_is_explicit_and_preserves_the_reserved_successor(self):
        reserved = at(10)
        control = FakeControlPlane((
            ScheduledJobState("aone.scan", ScheduledJobStatus.RUNNING, reserved),))
        runner = FakeRunner()
        engine = SchedulerEngine(
            (definition(),), control_plane=control, runtime=ScannerRuntime(runner),
            publisher=FakePublisher(), clock=lambda: at(9))

        self.assertEqual(engine.register(at(9))[0].status, ScheduledJobStatus.RUNNING)
        self.assertEqual(plan_due_slots((definition(),), control.list_jobs(), now=at(9)), ())
        self.assertEqual(engine.tick(at(9)), ())
        self.assertEqual(control.recoveries, 0)

        recovered = engine.recover_interrupted()
        planned = plan_due_slots((definition(),), control.list_jobs(), now=at(9))
        outcomes = engine.tick(at(9))

        self.assertEqual([(state.job_key, state.status, state.next_run_at) for state in recovered], [
            ("aone.scan", ScheduledJobStatus.IDLE, reserved),
        ])
        self.assertEqual(control.recoveries, 1)
        self.assertEqual(planned, ())
        self.assertEqual(outcomes, ())
        self.assertEqual(runner.calls, [])

    def test_runtime_stop_rejects_future_runner_invocations(self):
        runtime = ScannerRuntime(FakeRunner())
        runtime.stop()
        with self.assertRaisesRegex(RuntimeError, "refusing"):
            runtime.execute(definition(), at(9))


if __name__ == "__main__":
    unittest.main()
