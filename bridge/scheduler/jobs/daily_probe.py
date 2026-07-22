"""New-Scheduler runner for the daily Terraform probe."""

from __future__ import annotations

from datetime import datetime
import threading
import uuid
from typing import Any, Callable
from zoneinfo import ZoneInfo

from ..model import JobResult, JobResultStatus, ScheduledJobDefinition


JOB_KEY = "daily.probe"
RUNNER_KEY = "daily.probe"
_SHANGHAI = ZoneInfo("Asia/Shanghai")


class DailyProbeRunner:
    """Execute one Scheduler-admitted probe round without a legacy daily loop."""

    def __init__(self, *, handler: Any, pool: Any, execution_router: Any,
                 build_context: Callable[[str], tuple[Any, str, Any, Any, Any]],
                 logger: Any) -> None:
        self._handler = handler
        self._pool = pool
        self._router = execution_router
        self._build_context = build_context
        self._log = logger

    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        if definition.id != JOB_KEY:
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="daily probe runner received mismatched definition")
        if self._handler is None or self._pool is None:
            return JobResult(JobResultStatus.RETRYABLE_FAILURE,
                             error="daily probe runtime is unavailable")
        local_time = scheduled_for.astimezone(_SHANGHAI).replace(tzinfo=None)
        round_id = "probe-%s" % local_time.date().isoformat()
        envelope, prompt, notify, target, target_type = self._build_context(round_id)
        completion = threading.Event()
        execution: dict[str, Any] = {}

        def local_submit() -> tuple[bool, str]:
            session_id = str(uuid.uuid4())

            def work() -> Any:
                try:
                    result = self._handler.dispatch_item(
                        round_id, prompt, session_id, False, notify, target, target_type,
                        kind="probe", terraform=True)
                    execution["result"] = result
                    return result
                except Exception as exc:  # Scheduler maps this to a retryable result.
                    execution["error"] = exc
                    raise
                finally:
                    completion.set()

            return self._pool.submit(round_id, work, notify=notify, kind="probe")

        route = self._router.route(envelope)
        accepted, reason = self._router.enqueue(envelope, local_submit=local_submit)
        if not accepted:
            retryable = route.needs_recovery or reason == "queue_full"
            if retryable:
                return JobResult(JobResultStatus.RETRYABLE_FAILURE,
                                 error="daily probe not submitted: %s" % reason)
            return JobResult(JobResultStatus.SUCCEEDED)
        self._log.info("DailyProbeRunner: round %s accepted (durable=%s)",
                       round_id, route.needs_recovery)
        if route.needs_recovery:
            return JobResult(JobResultStatus.SUCCEEDED)
        completion.wait()
        if execution.get("error") is not None or execution.get("result") in ("error", "cancelled"):
            return JobResult(JobResultStatus.RETRYABLE_FAILURE,
                             error="daily probe execution failed")
        return JobResult(JobResultStatus.SUCCEEDED)


__all__ = ["JOB_KEY", "RUNNER_KEY", "DailyProbeRunner"]
