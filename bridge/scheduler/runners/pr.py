"""Scheduler-owned single tick for registered pull-request lifecycle watch.

``PrWatchScheduler`` remains the legacy domain implementation during this
cutover.  This module owns its SchedulerEngine integration and deliberately
supplies only the synchronous control-plane collaborators the legacy class
needs.  In particular, it never creates ``JarvisHandler`` or starts the
legacy scheduler thread.
"""

from __future__ import annotations

from datetime import datetime, timedelta, timezone
import importlib
from pathlib import Path
from threading import RLock
from typing import Any

from ..model import (
    JobResult,
    JobResultStatus,
    ScheduledJobDefinition,
    is_aware,
)


PR_WATCH_RUNNER_KEY = "pr.watch"


def _load_pr_module() -> Any:
    """Import the legacy PR domain implementation only when its slot is due."""

    try:
        return importlib.import_module("bridge.jarvis_dingtalk_bot")
    except ModuleNotFoundError:
        return importlib.import_module("jarvis_dingtalk_bot")


class PrWatchRuntimeContext:
    """Explicit, scheduler-safe dependencies required by ``PrWatchScheduler``.

    The legacy class uses the control-plane client for attention Tasks.  Its
    handler/pool paths are intentionally unavailable in the scheduler process:
    persistent execution is represented by ``ExecutionRouter`` and no local
    executor, DingTalk client, background handler, or worker watchdog is
    constructed here.
    """

    def __init__(self, *, task_client: Any, logger: Any,
                 legacy_module: Any | None = None) -> None:
        self._task_client = task_client
        self._log = logger
        self._module: Any | None = legacy_module
        self._watcher: Any | None = None
        self._lock = RLock()

    @property
    def module(self) -> Any:
        with self._lock:
            if self._module is None:
                self._module = _load_pr_module()
            return self._module

    def _build_runtime_dependencies(self) -> None:
        """Build only synchronous control-plane collaborators on first use."""

        with self._lock:
            if hasattr(self, "task_client"):
                return
            module = self.module
            self.task_client = self._task_client
            self.execution_router = module.ExecutionRouter(
                client=self.task_client, logger=self._log)
            self.ephemeral_executor = None
            self.no_dingtalk = True

    @property
    def watcher(self) -> Any:
        with self._lock:
            if self._watcher is None:
                self._build_runtime_dependencies()
                self._watcher = self.module.PrWatchScheduler(self, pool=None)
            return self._watcher


class PrWatchRunner:
    """Preserve PR watch's active/idle back-off as an adaptive Job result."""

    def __init__(self, *, job_id: str, context: PrWatchRuntimeContext,
                 logger: Any) -> None:
        self._job_id = job_id
        self._context = context
        self._log = logger

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != self._job_id:
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="pr watch runner received mismatched definition")
        if not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="pr watch runner requires an aware scheduled time")
        watcher = self._context.watcher
        if not watcher.enabled:
            self._log.info("pr.watch disabled by JARVIS_PRWATCH_ENABLE")
            return JobResult(JobResultStatus.SUCCEEDED)
        module = self._context.module
        module._aone_event_flush()
        module._dingtalk_event_flush()
        active = watcher._tick()
        delay = watcher._active_interval if active else watcher.interval
        return JobResult(
            JobResultStatus.SUCCEEDED,
            next_due_at=datetime.now(timezone.utc) + timedelta(seconds=delay),
        )


def build_pr_watch_runners(*, logger: Any, task_client: Any,
                           repo_root: Path) -> dict[str, object]:
    context = PrWatchRuntimeContext(task_client=task_client, logger=logger)
    return {
        PR_WATCH_RUNNER_KEY: PrWatchRunner(
            job_id=PR_WATCH_RUNNER_KEY, context=context, logger=logger),
    }


__all__ = [
    "PR_WATCH_RUNNER_KEY", "PrWatchRunner", "PrWatchRuntimeContext",
    "build_pr_watch_runners",
]
