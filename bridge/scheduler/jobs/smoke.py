"""Side-effect-free runner used to verify the Scheduler control path."""

from __future__ import annotations

from datetime import datetime
from typing import Any

from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


JOB_PREFIX = "smoke."
RUNNER_KEY = "scheduler.smoke"


class SchedulerSmokeRunner:
    """Prove admission reaches a runner without creating business work.

    The successful result exercises the normal ``start → execute → complete``
    path. It deliberately creates no Task, Event, Aone request, subprocess, or
    headless session; the control-plane timestamps and this log line are the
    observable test result.
    """

    def __init__(self, *, logger: Any) -> None:
        self._log = logger

    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        if not definition.id.startswith(JOB_PREFIX):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="scheduler smoke runner received mismatched definition")
        if not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="scheduler smoke runner requires an aware scheduled time")
        self._log.info(
            "SchedulerSmokeRunner: completed job=%s revision=%s schedule=%s scheduled_for=%s",
            definition.id, definition.revision, type(definition.schedule).__name__,
            scheduled_for.isoformat(),
        )
        return JobResult(JobResultStatus.SUCCEEDED)


__all__ = ["JOB_PREFIX", "RUNNER_KEY", "SchedulerSmokeRunner"]
