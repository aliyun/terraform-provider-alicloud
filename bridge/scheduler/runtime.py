"""Import-safe execution boundary for one scheduled scanner invocation.

This module deliberately does not choose a process model or import the legacy
Bridge loops.  A later composition root may adapt ``JobRunner`` to a bounded
subprocess, a headless invocation, or an in-process handler.  The scheduler
only needs the uniform result and stop boundary defined here.
"""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime
from threading import RLock
from typing import Protocol

from .model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


class RetryableJobError(Exception):
    """A runner may raise this when retrying the current scheduled slot is safe."""


class PermanentJobError(Exception):
    """A runner may raise this when the current slot cannot succeed by retrying."""


class ScannerStopped(RuntimeError):
    """Raised when a new scanner invocation is requested after stop admission closes."""


class JobRunner(Protocol):
    """Adapter for a bounded job runner; implementations must not create a scheduler loop."""

    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        """Run one slot and return its validated result."""


@dataclass(frozen=True)
class ScannerInvocation:
    """The minimal traceable identity of a runner invocation."""

    job_key: str
    scheduled_for: datetime

    def __post_init__(self) -> None:
        if not self.job_key:
            raise ValueError("job_key must be nonblank")
        if not is_aware(self.scheduled_for):
            raise ValueError("scheduled_for must be timezone-aware")


class ScannerRuntime:
    """Run one bounded scanner at a time and map adapter failures to ``JobResult``.

    ``stop`` is deliberately admission-only in this first runtime slice: it
    prevents *new* scanners while ``execute_admitted`` lets a slot already
    committed by the control plane return its result, so ``SchedulerEngine``
    can make the terminal update.
    U7 adds process-group cancellation and deadline handling around this seam.
    """

    def __init__(self, runner: JobRunner) -> None:
        if not hasattr(runner, "run"):
            raise TypeError("runner must implement run(definition, scheduled_for)")
        self._runner = runner
        self._lock = RLock()
        self._stopped = False

    @property
    def accepting(self) -> bool:
        with self._lock:
            return not self._stopped

    def stop(self) -> None:
        """Close admission for future scanner invocations without interrupting one in flight."""

        with self._lock:
            self._stopped = True

    def resume(self) -> None:
        """Re-open admission when a bounded planned restart is cancelled."""

        with self._lock:
            self._stopped = False

    def execute(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        """Execute one slot without publishing work or mutating scheduler state."""

        invocation = ScannerInvocation(definition.id, scheduled_for)
        with self._lock:
            if self._stopped:
                raise ScannerStopped(f"scheduler is stopping; refusing {invocation.job_key}")
        return self._execute_runner(definition, invocation)

    def execute_admitted(
        self, definition: ScheduledJobDefinition, scheduled_for: datetime,
    ) -> JobResult:
        """Finish a control-plane-admitted slot even if drain starts concurrently."""

        invocation = ScannerInvocation(definition.id, scheduled_for)
        return self._execute_runner(definition, invocation)

    def _execute_runner(
        self, definition: ScheduledJobDefinition, invocation: ScannerInvocation,
    ) -> JobResult:
        try:
            result = self._runner.run(definition, invocation.scheduled_for)
            if not isinstance(result, JobResult):
                raise PermanentJobError("runner returned a non-JobResult value")
            return result
        except RetryableJobError as exc:
            return JobResult(JobResultStatus.RETRYABLE_FAILURE, error=_error_summary(exc))
        except PermanentJobError as exc:
            return JobResult(JobResultStatus.PERMANENT_FAILURE, error=_error_summary(exc))
        except Exception as exc:  # Runner adapters are untrusted; unknown failures are replay-safe by default.
            return JobResult(JobResultStatus.RETRYABLE_FAILURE, error=_error_summary(exc))


def _error_summary(exc: Exception) -> str:
    """Return a bounded one-line error summary; the control plane sanitizes it again."""

    message = " ".join(str(exc).split())
    summary = type(exc).__name__ if not message else f"{type(exc).__name__}: {message}"
    return summary[:2048]
