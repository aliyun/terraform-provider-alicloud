from __future__ import annotations

from datetime import datetime, timezone
import logging
import unittest

from bridge.scheduler import (
    HandlerRunner,
    IntervalSchedule,
    JobResultStatus,
    MisfirePolicy,
    ScheduledJobDefinition,
)
from bridge.scheduler.runners.recovery import ExternalRecoveryRunner


UTC = timezone.utc


def definition(job_id: str = "external.recovery") -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        job_id, 1, "test", IntervalSchedule(60, True), HandlerRunner(job_id),
        MisfirePolicy.COALESCE, 5,
    )


class FakeTaskClient:
    def __init__(self) -> None:
        self.pages: list[tuple[int, int]] = []

    def list_external_operation_recovery_candidates(
        self, *, after_operation_id: int, limit: int,
    ) -> dict[str, object]:
        self.pages.append((after_operation_id, limit))
        return {"items": [], "hasMore": False}


class ExternalRecoveryRunnerTests(unittest.TestCase):
    def test_runs_one_tick_with_explicit_runtime_dependencies(self):
        client = FakeTaskClient()
        runner = ExternalRecoveryRunner(
            task_client=client,
            worker_key="scheduler-1",
            repo_root=".",
            logger=logging.getLogger(__name__),
            enabled=True,
        )

        result = runner.run(definition(), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.SUCCEEDED)
        self.assertEqual(client.pages, [(0, 100)])

    def test_disabled_runner_skips_the_control_plane_tick(self):
        client = FakeTaskClient()
        runner = ExternalRecoveryRunner(
            task_client=client,
            worker_key="scheduler-1",
            repo_root=".",
            logger=logging.getLogger(__name__),
            enabled=False,
        )

        result = runner.run(definition(), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.SUCCEEDED)
        self.assertEqual(client.pages, [])

    def test_invalid_scheduler_invocation_is_permanent(self):
        runner = ExternalRecoveryRunner(
            task_client=FakeTaskClient(),
            worker_key="scheduler-1",
            repo_root=".",
            logger=logging.getLogger(__name__),
            enabled=True,
        )

        result = runner.run(definition("aone.scan"), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.PERMANENT_FAILURE)


if __name__ == "__main__":
    unittest.main()
