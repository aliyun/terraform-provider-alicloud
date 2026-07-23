"""Scheduler-owned single ticks for Aone discovery and claim health.

The scheduler owns cadence and error handling.  The still-legacy
``AoneScheduler`` supplies the bounded domain operations, but it is given a
small runtime dependency container rather than a ``JarvisHandler``.  This
module deliberately has no module-level Bot import and never constructs a
handler or background loop.
"""

from __future__ import annotations

from datetime import datetime
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


AONE_SCAN_RUNNER_KEY = "aone.scan"
AONE_CLAIM_HEALTH_RUNNER_KEY = "aone.claim-health"


def _load_aone_module() -> Any:
    """Import the legacy Aone domain implementation only when it is needed."""

    try:
        return importlib.import_module("bridge.jarvis_dingtalk_bot")
    except ModuleNotFoundError:
        return importlib.import_module("jarvis_dingtalk_bot")


class AoneRuntimeContext:
    """Explicit, scheduler-safe dependencies required by ``AoneScheduler``.

    It intentionally excludes every handler-owned collaborator (DingTalk,
    local executors, worker watchdogs, and background threads).
    """

    def __init__(self, *, task_client: Any, repo_root: Path,
                 logger: Any, legacy_module: Any | None = None) -> None:
        self._task_client = task_client
        self._repo_root = repo_root
        self._log = logger
        self._module: Any | None = legacy_module
        self._scanner: Any | None = None
        self._lock = RLock()

    @property
    def module(self) -> Any:
        with self._lock:
            if self._module is None:
                self._module = _load_aone_module()
            return self._module

    def _build_runtime_dependencies(self) -> None:
        with self._lock:
            if hasattr(self, "task_client"):
                return
            module = self.module
            self.task_client = self._task_client
            self.execution_router = module.ExecutionRouter(
                client=self.task_client, logger=self._log)
            self.field_repair_worker = module.FieldRepairWorker(
                repo_root=self._repo_root,
                client=self.task_client,
                runtime=module.DEFAULT_EXECUTION_RUNTIME,
                claude_bin=module.claude_bin(),
            )
            self.ephemeral_executor = None
            self.no_dingtalk = True

    @property
    def scanner(self) -> Any:
        with self._lock:
            if self._scanner is None:
                self._build_runtime_dependencies()
                self._scanner = self.module.AoneScheduler(self, pool=None)
            return self._scanner


class _AoneRunner:
    def __init__(self, *, job_id: str, context: AoneRuntimeContext,
                 logger: Any) -> None:
        self._job_id = job_id
        self._context = context
        self._log = logger

    def _matches(self, definition: ScheduledJobDefinition,
                 scheduled_for: datetime) -> JobResult | None:
        if definition.id != self._job_id:
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="Aone scheduler runner received mismatched definition")
        if not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="Aone scheduler runner requires an aware scheduled time")
        return None

    @staticmethod
    def _success() -> JobResult:
        return JobResult(JobResultStatus.SUCCEEDED)


class AoneScanRunner(_AoneRunner):
    """Execute one discovery/dispatch tick without starting a legacy loop."""

    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        invalid = self._matches(definition, scheduled_for)
        if invalid:
            return invalid
        self._context.scanner._tick()
        return self._success()


class AoneClaimHealthRunner(_AoneRunner):
    """Execute one claim-health reconciliation and event flush."""

    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        invalid = self._matches(definition, scheduled_for)
        if invalid:
            return invalid
        module = self._context.module
        if (Path(module.REPO_ROOT) / ".my-day" / "bridge" / "pause").exists():
            self._log.info("aone.claim-health: pause flag present; skip this slot")
            return self._success()
        scanner = self._context.scanner
        scanner._claim_health_activity_cache = {}
        snapshot = scanner._scan_claimed()
        if snapshot is not None:
            scanner._reconcile_stale_claims(snapshot)
        module._aone_event_flush()
        module._dingtalk_event_flush()
        return self._success()


def build_aone_runners(*, logger: Any, task_client: Any,
                       repo_root: Path) -> dict[str, object]:
    context = AoneRuntimeContext(
        task_client=task_client, repo_root=repo_root, logger=logger)
    return {
        AONE_SCAN_RUNNER_KEY: AoneScanRunner(
            job_id=AONE_SCAN_RUNNER_KEY, context=context, logger=logger),
        AONE_CLAIM_HEALTH_RUNNER_KEY: AoneClaimHealthRunner(
            job_id=AONE_CLAIM_HEALTH_RUNNER_KEY, context=context, logger=logger),
    }


__all__ = [
    "AONE_CLAIM_HEALTH_RUNNER_KEY", "AONE_SCAN_RUNNER_KEY",
    "AoneClaimHealthRunner", "AoneRuntimeContext", "AoneScanRunner",
    "build_aone_runners",
]
