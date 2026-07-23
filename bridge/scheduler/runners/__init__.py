"""Registered in-process runners for Scheduler-owned jobs."""

from __future__ import annotations

from typing import Any

from .smoke import RUNNER_KEY as SCHEDULER_SMOKE_RUNNER_KEY, SchedulerSmokeRunner


RUNNER_KEYS = frozenset({SCHEDULER_SMOKE_RUNNER_KEY})


def build_runners(*, logger: Any) -> dict[str, SchedulerSmokeRunner]:
    """Build the complete runner catalogue for this Scheduler process."""

    return {
        SCHEDULER_SMOKE_RUNNER_KEY: SchedulerSmokeRunner(logger=logger),
    }


__all__ = [
    "RUNNER_KEYS",
    "SCHEDULER_SMOKE_RUNNER_KEY",
    "SchedulerSmokeRunner",
    "build_runners",
]
