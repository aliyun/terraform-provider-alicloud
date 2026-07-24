"""Runner for the Aone discovery job."""

from __future__ import annotations

from datetime import datetime

from bridge.aone_workitems import AoneRuntime
from .helpers import claude_bin
from bridge.jarvis_execution_runtime import DEFAULT_EXECUTION_RUNTIME
from bridge.jarvis_field_repair import FieldRepairWorker
from bridge.jarvis_task_router import ExecutionRouter
from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "scan"


class ScanRunner:
    def __init__(self, runtime: AoneRuntime) -> None:
        self.runtime = runtime

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != "aone.scan" or not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="aone.scan runner received an invalid slot")
        self.runtime._tick()
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger, task_client, repo_root):
    runtime = AoneRuntime(
        handler=None, pool=None,
        execution_router=ExecutionRouter(client=task_client, logger=logger),
        field_repair_worker=FieldRepairWorker(
            repo_root=repo_root, client=task_client,
            runtime=DEFAULT_EXECUTION_RUNTIME, claude_bin=claude_bin()),
    )
    return ScanRunner(runtime)
