"""Scheduler-owned runner for external operation recovery.

The recovery scheduler still owns the control-plane lease protocol and Aone
read-back logic.  This adapter gives that protocol one SchedulerEngine-admitted
slot at a time; it deliberately never creates the legacy Bridge handler or its
background loop.
"""

from __future__ import annotations

from datetime import datetime
from pathlib import Path
from typing import Any

try:  # Package tests import ``bridge.scheduler``; bridge/main.py imports ``scheduler``.
    from bridge.jarvis_external_recovery import ExternalOperationRecoveryScheduler
except ModuleNotFoundError:  # pragma: no cover - exercised by bridge/main.py.
    from jarvis_external_recovery import ExternalOperationRecoveryScheduler

from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "external.recovery"


class ExternalRecoveryRunner:
    """Run exactly one external-recovery tick for an admitted Scheduler slot."""

    def __init__(
        self,
        *,
        task_client: Any,
        worker_key: str,
        repo_root: Path | str,
        logger: Any,
        enabled: bool | None = None,
    ) -> None:
        self._logger = logger
        self._recovery = ExternalOperationRecoveryScheduler(
            task_client,
            worker_key,
            repo_root=repo_root,
            enabled=enabled,
            logger=logger,
        )

    def run(
        self, definition: ScheduledJobDefinition, scheduled_for: datetime,
    ) -> JobResult:
        """Invoke the bounded recovery protocol, or skip if its kill switch is off."""
        invalid = self._definition_error(definition, scheduled_for)
        if invalid is not None:
            return invalid
        if not self._recovery.enabled:
            self._logger.info("external.recovery disabled by JARVIS_EXTERNAL_RECOVERY_ENABLE")
            return JobResult(JobResultStatus.SUCCEEDED)
        self._recovery._tick()
        return JobResult(JobResultStatus.SUCCEEDED)

    @staticmethod
    def _definition_error(
        definition: ScheduledJobDefinition, scheduled_for: datetime,
    ) -> JobResult | None:
        if definition.id != RUNNER_KEY:
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="external.recovery runner received mismatched definition",
            )
        if not is_aware(scheduled_for):
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="external.recovery runner requires an aware scheduled time",
            )
        return None


__all__ = ["ExternalRecoveryRunner", "RUNNER_KEY"]
