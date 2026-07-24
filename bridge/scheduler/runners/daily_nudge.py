"""Runner for the daily stale-workitem reminder job."""

from __future__ import annotations

from datetime import datetime

from bridge.aone_workitems import DailyNudge
from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "daily_nudge"


class DailyNudgeRunner:
    def __init__(self, job: DailyNudge, logger) -> None:
        self.job, self.logger = job, logger

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != "daily.nudge" or not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="daily.nudge runner received an invalid slot")
        if self.job.enabled:
            self.job.run()
        else:
            self.logger.info("daily.nudge disabled by JARVIS_REVISIT_SCHED")
        return JobResult(JobResultStatus.SUCCEEDED)


def build(*, logger, task_client, repo_root):
    del task_client, repo_root
    return DailyNudgeRunner(DailyNudge(), logger)
