"""Registered in-process runners for Scheduler-owned jobs."""

from __future__ import annotations

from pathlib import Path
from typing import Any

if (__package__ or "").startswith("bridge."):
    from bridge.headless.jarvis_adapter import (
        jarvis_transcript_exists,
        run_jarvis_attempt,
    )
    from bridge.headless.runtime import HeadlessRuntime
else:  # bridge/main.py imports scheduler as a top-level package.
    from headless.jarvis_adapter import (
        jarvis_transcript_exists,
        run_jarvis_attempt,
    )
    from headless.runtime import HeadlessRuntime

from .daily_probe import (
    RUNNER_KEY as DAILY_PROBE_RUNNER_KEY,
    DailyProbeRunner,
)
from .smoke import RUNNER_KEY as SCHEDULER_SMOKE_RUNNER_KEY, SchedulerSmokeRunner


HANDLER_KEYS = frozenset({SCHEDULER_SMOKE_RUNNER_KEY})
HEADLESS_BUILDER_PROTOCOLS = frozenset({DAILY_PROBE_RUNNER_KEY})
RUNNER_KEYS = HANDLER_KEYS | HEADLESS_BUILDER_PROTOCOLS


def build_runners(
    *,
    logger: Any,
    headless_runtime: HeadlessRuntime | None = None,
    summary_root: Path | None = None,
) -> dict[object, object]:
    """Build the complete runner catalogue for this Scheduler process."""

    runtime = headless_runtime or HeadlessRuntime(
        run_jarvis_attempt,
        transcript_exists=jarvis_transcript_exists,
        logger=logger,
    )
    root = summary_root or (
        Path(__file__).resolve().parents[3] / "runs" / "probe")
    return {
        SCHEDULER_SMOKE_RUNNER_KEY: SchedulerSmokeRunner(logger=logger),
        DAILY_PROBE_RUNNER_KEY: DailyProbeRunner(
            runtime=runtime,
            summary_root=root,
            logger=logger,
        ),
    }


__all__ = [
    "DAILY_PROBE_RUNNER_KEY",
    "DailyProbeRunner",
    "HANDLER_KEYS",
    "HEADLESS_BUILDER_PROTOCOLS",
    "RUNNER_KEYS",
    "SCHEDULER_SMOKE_RUNNER_KEY",
    "SchedulerSmokeRunner",
    "build_runners",
]
