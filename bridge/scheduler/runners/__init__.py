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
from .legacy import (
    AONE_CLAIM_HEALTH_RUNNER_KEY,
    AONE_REPLY_RUNNER_KEY,
    AONE_SCAN_RUNNER_KEY,
    DAILY_NUDGE_RUNNER_KEY,
    EXTERNAL_RECOVERY_RUNNER_KEY,
    PR_WATCH_RUNNER_KEY,
    build_legacy_runners,
)
HANDLER_KEYS = frozenset({
    AONE_SCAN_RUNNER_KEY,
    AONE_CLAIM_HEALTH_RUNNER_KEY,
    DAILY_NUDGE_RUNNER_KEY,
    AONE_REPLY_RUNNER_KEY,
    PR_WATCH_RUNNER_KEY,
    EXTERNAL_RECOVERY_RUNNER_KEY,
})
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
    runners = {
        DAILY_PROBE_RUNNER_KEY: DailyProbeRunner(
            runtime=runtime,
            summary_root=root,
            logger=logger,
        ),
    }
    runners.update(build_legacy_runners(logger=logger))
    return runners


__all__ = [
    "DAILY_PROBE_RUNNER_KEY",
    "AONE_CLAIM_HEALTH_RUNNER_KEY",
    "AONE_REPLY_RUNNER_KEY",
    "AONE_SCAN_RUNNER_KEY",
    "DAILY_NUDGE_RUNNER_KEY",
    "DailyProbeRunner",
    "EXTERNAL_RECOVERY_RUNNER_KEY",
    "PR_WATCH_RUNNER_KEY",
    "HANDLER_KEYS",
    "HEADLESS_BUILDER_PROTOCOLS",
    "RUNNER_KEYS",
    "build_runners",
]
