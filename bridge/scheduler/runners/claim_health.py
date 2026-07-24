"""Runner for the Aone claim-health reconciliation job."""

from __future__ import annotations

from datetime import datetime
from pathlib import Path

from bridge.aone_workitems import AoneRuntime
from bridge.aone_events import _aone_event_flush, _dingtalk_event_flush
from .helpers import claude_bin
from bridge.jarvis_execution_runtime import DEFAULT_EXECUTION_RUNTIME
from bridge.jarvis_field_repair import FieldRepairWorker
from bridge.jarvis_task_router import ExecutionRouter
from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "claim_health"


class ClaimHealthRunner:
    def __init__(self, runtime: AoneRuntime, repo_root: Path, logger) -> None:
        self.runtime, self.repo_root, self.logger = runtime, repo_root, logger

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != "aone.claim-health" or not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="aone.claim-health runner received an invalid slot")
        if (self.repo_root / ".my-day" / "bridge" / "pause").exists():
            self.logger.info("aone.claim-health: pause flag present; skip this slot")
            return JobResult(JobResultStatus.SUCCEEDED)
        self.runtime._claim_health_activity_cache = {}
        snapshot = self.runtime._scan_claimed()
        if snapshot is not None:
            self.runtime._reconcile_stale_claims(snapshot)
        _aone_event_flush()
        _dingtalk_event_flush()
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger, task_client, repo_root):
    runtime = AoneRuntime(
        handler=None, pool=None,
        execution_router=ExecutionRouter(client=task_client, logger=logger),
        field_repair_worker=FieldRepairWorker(
            repo_root=repo_root, client=task_client,
            runtime=DEFAULT_EXECUTION_RUNTIME, claude_bin=claude_bin()),
    )
    return ClaimHealthRunner(runtime, Path(repo_root), logger)
