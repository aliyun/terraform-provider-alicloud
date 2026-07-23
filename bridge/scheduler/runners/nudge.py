"""Scheduler-owned single tick for the daily stale-work-item reminder.

The remaining nudge domain logic lives in ``jarvis_dingtalk_bot._NudgeJob``
during the migration, but it only needs a tiny synchronous context.  In
particular, this runner must never construct ``JarvisHandler`` or a local
executor: SchedulerEngine owns its cadence and admission.
"""

from __future__ import annotations

from datetime import datetime
import importlib
from threading import RLock
from typing import Any

from ..model import (
    JobResult,
    JobResultStatus,
    ScheduledJobDefinition,
    is_aware,
)


DAILY_NUDGE_RUNNER_KEY = "daily.nudge"


def _load_nudge_module() -> Any:
    """Import the legacy nudge implementation only when its slot is due."""

    try:
        return importlib.import_module("bridge.jarvis_dingtalk_bot")
    except ModuleNotFoundError:
        return importlib.import_module("jarvis_dingtalk_bot")


class NudgeRuntimeContext:
    """The intentionally small legacy dependency surface for ``_NudgeJob``.

    ``_NudgeJob`` currently inspects these two attributes in its constructor.
    Keeping them explicit makes the scheduler process incapable of gaining a
    DingTalk client, handler-owned worker pool, or background loop by accident.
    """

    ephemeral_executor = None
    execution_router = None
    no_dingtalk = True

    def __init__(self, *, legacy_module: Any | None = None) -> None:
        self._module = legacy_module
        self._lock = RLock()

    @property
    def module(self) -> Any:
        with self._lock:
            if self._module is None:
                self._module = _load_nudge_module()
            return self._module


class DailyNudgeRunner:
    """Run one enabled ``_NudgeJob`` tick under SchedulerEngine control."""

    def __init__(self, *, job_id: str, context: NudgeRuntimeContext,
                 logger: Any) -> None:
        self._job_id = job_id
        self._context = context
        self._log = logger

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != self._job_id:
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="daily nudge runner received mismatched definition")
        if not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="daily nudge runner requires an aware scheduled time")
        # The legacy feature flag remains a kill switch.  YAML owns the cadence;
        # this flag only determines whether the business action is performed.
        job = self._context.module._NudgeJob(self._context)
        if not job.enabled:
            self._log.info("daily.nudge disabled by JARVIS_REVISIT_SCHED")
            return JobResult(JobResultStatus.SUCCEEDED)
        job.run()
        return JobResult(JobResultStatus.SUCCEEDED)


def build_nudge_runners(*, logger: Any) -> dict[str, object]:
    context = NudgeRuntimeContext()
    return {
        DAILY_NUDGE_RUNNER_KEY: DailyNudgeRunner(
            job_id=DAILY_NUDGE_RUNNER_KEY, context=context, logger=logger),
    }


__all__ = [
    "DAILY_NUDGE_RUNNER_KEY", "DailyNudgeRunner", "NudgeRuntimeContext",
    "build_nudge_runners",
]
