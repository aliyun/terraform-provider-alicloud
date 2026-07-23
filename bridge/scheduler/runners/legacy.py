"""Adapters that move legacy Bridge periodic ticks into SchedulerEngine.

The business logic intentionally stays in :mod:`jarvis_dingtalk_bot` for this
cutover.  These adapters own cadence, admission, retries, and observability;
the legacy classes supply one bounded unit of work only.  No adapter calls a
legacy ``.start()`` method, so SchedulerEngine is the sole loop owner.
"""

from __future__ import annotations

from datetime import datetime, timedelta, timezone
import importlib
from pathlib import Path
from threading import RLock
from typing import Any, Callable

from ..model import (
    JobResult,
    JobResultStatus,
    ScheduledJobDefinition,
    is_aware,
)


AONE_SCAN_RUNNER_KEY = "aone.scan"
AONE_CLAIM_HEALTH_RUNNER_KEY = "aone.claim-health"
DAILY_NUDGE_RUNNER_KEY = "daily.nudge"
AONE_REPLY_RUNNER_KEY = "aone.reply"
PR_WATCH_RUNNER_KEY = "pr.watch"
EXTERNAL_RECOVERY_RUNNER_KEY = "external.recovery"


def _load_legacy_module() -> Any:
    """Import the same module for package tests and ``bridge/main.py``."""

    try:
        return importlib.import_module("bridge.jarvis_dingtalk_bot")
    except ModuleNotFoundError:
        return importlib.import_module("jarvis_dingtalk_bot")


class LegacyBridgeContext:
    """Lazily construct one no-DingTalk handler shared by all adapters."""

    def __init__(self, *, handler_factory: Callable[[], Any] | None = None,
                 legacy_module: Any | None = None) -> None:
        self._handler_factory = handler_factory
        self._handler: Any | None = None
        self._module: Any | None = legacy_module
        self._lock = RLock()

    @property
    def module(self) -> Any:
        with self._lock:
            if self._module is None:
                self._module = _load_legacy_module()
            return self._module

    @property
    def handler(self) -> Any:
        with self._lock:
            if self._handler is None:
                factory = self._handler_factory
                self._handler = factory() if factory else self.module.JarvisHandler(no_dingtalk=True)
            return self._handler


class _LegacyRunner:
    """Shared identity and result validation for migrated legacy ticks."""

    def __init__(self, *, job_id: str, context: LegacyBridgeContext, logger: Any) -> None:
        self._job_id = job_id
        self._context = context
        self._log = logger

    def _matches(self, definition: ScheduledJobDefinition,
                 scheduled_for: datetime) -> JobResult | None:
        if definition.id != self._job_id:
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="legacy scheduler runner received mismatched definition")
        if not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="legacy scheduler runner requires an aware scheduled time")
        return None

    @staticmethod
    def _success(*, next_due_at: datetime | None = None) -> JobResult:
        return JobResult(JobResultStatus.SUCCEEDED, next_due_at=next_due_at)


class AoneScanRunner(_LegacyRunner):
    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        invalid = self._matches(definition, scheduled_for)
        if invalid:
            return invalid
        self._context.handler.scanner._tick()
        return self._success()


class AoneClaimHealthRunner(_LegacyRunner):
    """The former independent claim-health loop, executed once per slot."""

    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        invalid = self._matches(definition, scheduled_for)
        if invalid:
            return invalid
        module = self._context.module
        if (Path(module.REPO_ROOT) / ".my-day" / "bridge" / "pause").exists():
            self._log.info("aone.claim-health: pause flag present; skip this slot")
            return self._success()
        scanner = self._context.handler.scanner
        scanner._claim_health_activity_cache = {}
        snapshot = scanner._scan_claimed()
        if snapshot is not None:
            scanner._reconcile_stale_claims(snapshot)
        module._aone_event_flush()
        module._dingtalk_event_flush()
        return self._success()


class DailyNudgeRunner(_LegacyRunner):
    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        invalid = self._matches(definition, scheduled_for)
        if invalid:
            return invalid
        # The legacy feature flag remains a kill switch.  YAML owns the cadence;
        # this flag only determines whether the business action is performed.
        job = self._context.module._NudgeJob(self._context.handler)
        if not job.enabled:
            self._log.info("daily.nudge disabled by JARVIS_REVISIT_SCHED")
            return self._success()
        job.run()
        return self._success()


class AoneReplyRunner(_LegacyRunner):
    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        invalid = self._matches(definition, scheduled_for)
        if invalid:
            return invalid
        self._context.handler.aone_reply_scheduler._tick()
        return self._success()


class PrWatchRunner(_LegacyRunner):
    """Preserve PR watch's active/idle back-off as an adaptive Job result."""

    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        invalid = self._matches(definition, scheduled_for)
        if invalid:
            return invalid
        watcher = self._context.handler.prwatch
        if not watcher.enabled:
            self._log.info("pr.watch disabled by JARVIS_PRWATCH_ENABLE")
            return self._success()
        module = self._context.module
        module._aone_event_flush()
        module._dingtalk_event_flush()
        active = watcher._tick()
        delay = watcher._active_interval if active else watcher.interval
        return self._success(next_due_at=datetime.now(timezone.utc) + timedelta(seconds=delay))


class ExternalRecoveryRunner(_LegacyRunner):
    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        invalid = self._matches(definition, scheduled_for)
        if invalid:
            return invalid
        recovery = self._context.handler.external_operation_recovery
        if not recovery.enabled:
            self._log.info("external.recovery disabled by JARVIS_EXTERNAL_RECOVERY_ENABLE")
            return self._success()
        recovery._tick()
        return self._success()


def build_legacy_runners(*, logger: Any) -> dict[str, object]:
    """Build all migrated business-loop adapters around one lazy handler."""

    context = LegacyBridgeContext()
    return {
        AONE_SCAN_RUNNER_KEY: AoneScanRunner(
            job_id="aone.scan", context=context, logger=logger),
        AONE_CLAIM_HEALTH_RUNNER_KEY: AoneClaimHealthRunner(
            job_id="aone.claim-health", context=context, logger=logger),
        DAILY_NUDGE_RUNNER_KEY: DailyNudgeRunner(
            job_id="daily.nudge", context=context, logger=logger),
        AONE_REPLY_RUNNER_KEY: AoneReplyRunner(
            job_id="aone.reply", context=context, logger=logger),
        PR_WATCH_RUNNER_KEY: PrWatchRunner(
            job_id="pr.watch", context=context, logger=logger),
        EXTERNAL_RECOVERY_RUNNER_KEY: ExternalRecoveryRunner(
            job_id="external.recovery", context=context, logger=logger),
    }


__all__ = [
    "AONE_CLAIM_HEALTH_RUNNER_KEY", "AONE_REPLY_RUNNER_KEY", "AONE_SCAN_RUNNER_KEY",
    "DAILY_NUDGE_RUNNER_KEY", "EXTERNAL_RECOVERY_RUNNER_KEY", "PR_WATCH_RUNNER_KEY",
    "AoneClaimHealthRunner", "AoneReplyRunner", "AoneScanRunner", "DailyNudgeRunner",
    "ExternalRecoveryRunner", "LegacyBridgeContext", "PrWatchRunner", "build_legacy_runners",
]
