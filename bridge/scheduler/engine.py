"""Single-loop admission and state-reporting skeleton for Bridge scheduled jobs.

The control plane remains the authoritative current state.  This module owns
only local admission for an Engine process and never assumes an HTTP client;
an application adapter supplies the protocol below once the AutomationAgent
endpoint is deployed.
"""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from enum import Enum
import logging
from threading import RLock
import time
from typing import Callable, Iterable, Mapping, Optional, Protocol, Sequence

from .model import JobResult, JobResultStatus, ScheduledJobDefinition, definition_snapshot, is_aware
from .planner import TriggerPlanner
from .runtime import ScannerRuntime


_LOG = logging.getLogger(__name__)


def _utc_now() -> datetime:
    return datetime.now(timezone.utc)


class ScheduledJobStatus(str, Enum):
    IDLE = "IDLE"
    RUNNING = "RUNNING"
    ERROR = "ERROR"
    DISABLED = "DISABLED"


@dataclass(frozen=True)
class ScheduledJobState:
    """Control-plane fields needed for candidate planning and slot fencing."""

    job_key: str
    status: ScheduledJobStatus
    next_run_at: Optional[datetime]

    def __post_init__(self) -> None:
        if not self.job_key:
            raise ValueError("job_key must be nonblank")
        if not isinstance(self.status, ScheduledJobStatus):
            raise TypeError("status must be a ScheduledJobStatus")
        if self.next_run_at is not None and not is_aware(self.next_run_at):
            raise ValueError("next_run_at must be timezone-aware when present")


@dataclass(frozen=True)
class JobRegistration:
    """The complete one-row registration payload; an adapter serializes it for the API."""

    job_key: str
    job_name: str
    definition: Mapping[str, object]
    next_run_at: Optional[datetime]
    enabled: bool


@dataclass(frozen=True)
class ScheduledSlot:
    definition: ScheduledJobDefinition
    scheduled_for: datetime

    def __post_init__(self) -> None:
        if not is_aware(self.scheduled_for):
            raise ValueError("scheduled_for must be timezone-aware")


class ExecutionDisposition(str, Enum):
    START_REJECTED = "START_REJECTED"
    COMPLETED = "COMPLETED"
    RETRYABLE_FAILURE = "RETRYABLE_FAILURE"
    PERMANENT_FAILURE = "PERMANENT_FAILURE"


@dataclass(frozen=True)
class ExecutionOutcome:
    job_key: str
    scheduled_for: datetime
    disposition: ExecutionDisposition


class ScheduledJobControlPlane(Protocol):
    """Adapter seam matching AutomationAgent's register/list/start/complete/fail lifecycle.

    ``start`` returns ``False`` for a stale or already-admitted ``scheduled_for``
    (for example an HTTP conflict).  Terminal calls are intentionally void: a
    transport failure must surface to the composition root instead of pretending
    that a durable acknowledgement happened.
    """

    def register(self, registrations: Sequence[JobRegistration]) -> Sequence[ScheduledJobState]:
        """Upsert the authoritative registry and return its current states."""

    def list_jobs(self) -> Sequence[ScheduledJobState]:
        """Return current job states from the control plane."""

    def recover_interrupted(self) -> Sequence[ScheduledJobState]:
        """Return interrupted ``RUNNING`` jobs to their pre-reserved successors.

        This is a startup/composition operation, deliberately not part of
        :meth:`SchedulerEngine.tick`.  The caller must complete registration,
        recovery, and then normal planning in that order.
        """

    def start(
        self, job_key: str, scheduled_for: datetime, next_run_at: datetime,
    ) -> bool:
        """Admit one slot and atomically reserve its fallback successor."""

    def complete(self, job_key: str, scheduled_for: datetime, next_run_at: datetime) -> None:
        """Record a successful slot after all durable publications are acknowledged."""

    def fail(
        self,
        job_key: str,
        scheduled_for: datetime,
        *,
        retryable: bool,
        next_run_at: Optional[datetime],
        error: str,
    ) -> None:
        """Record exactly one terminal failure for the admitted slot."""


class DurableResultPublisher(Protocol):
    """Publish successful observations only after their result protocol is validated."""

    def publish(self, definition: ScheduledJobDefinition, result: JobResult, scheduled_for: datetime) -> None:
        """Return only after Task/event durable acknowledgement is known."""


def build_registrations(
    definitions: Iterable[ScheduledJobDefinition],
    *,
    planner: TriggerPlanner,
    now: datetime,
    is_enabled: Callable[[ScheduledJobDefinition], bool] = lambda definition: True,
) -> tuple[JobRegistration, ...]:
    """Purely construct the complete registration snapshot for one Engine start."""

    _require_aware(now, "now")
    registrations = []
    for definition in definitions:
        enabled = bool(is_enabled(definition))
        registrations.append(JobRegistration(
            job_key=definition.id,
            job_name=definition.id,
            definition=definition_snapshot(definition),
            next_run_at=planner.initial_due(definition, now) if enabled else None,
            enabled=enabled,
        ))
    return tuple(registrations)


def plan_due_slots(
    definitions: Iterable[ScheduledJobDefinition],
    states: Iterable[ScheduledJobState],
    *,
    now: datetime,
) -> tuple[ScheduledSlot, ...]:
    """Pure candidate planning from registered definitions and authoritative state."""

    _require_aware(now, "now")
    definitions_by_key = {definition.id: definition for definition in definitions}
    candidates = []
    seen_slots: set[tuple[str, datetime]] = set()
    for state in states:
        definition = definitions_by_key.get(state.job_key)
        if definition is None or state.status not in (ScheduledJobStatus.IDLE, ScheduledJobStatus.ERROR):
            continue
        if state.next_run_at is None or state.next_run_at > now:
            continue
        key = (definition.id, state.next_run_at)
        if key in seen_slots:
            continue
        seen_slots.add(key)
        candidates.append(ScheduledSlot(definition, state.next_run_at))
    return tuple(sorted(candidates, key=lambda item: (item.scheduled_for, item.definition.id)))


class SchedulerEngine:
    """One synchronous Engine tick; callers choose the outer wait/clock mechanism.

    Keeping the tick synchronous in U3 makes job admission and status commits
    testable.  It also makes the one-job overlap invariant explicit; U7 can add
    bounded concurrent Scanner workers at the same protocol boundary.
    """

    def __init__(
        self,
        definitions: Iterable[ScheduledJobDefinition],
        *,
        control_plane: ScheduledJobControlPlane,
        runtime: ScannerRuntime,
        publisher: DurableResultPublisher,
        planner: Optional[TriggerPlanner] = None,
        clock: Callable[[], datetime] = _utc_now,
        terminal_retry_delay_seconds: float = 5.0,
        sleeper: Callable[[float], None] = time.sleep,
    ) -> None:
        self._definitions = tuple(definitions)
        self._definitions_by_key = {definition.id: definition for definition in self._definitions}
        if len(self._definitions_by_key) != len(self._definitions):
            raise ValueError("definitions must have unique job ids")
        self._control_plane = control_plane
        self._runtime = runtime
        self._publisher = publisher
        self._planner = planner or TriggerPlanner()
        self._clock = clock
        self._terminal_retry_delay_seconds = float(terminal_retry_delay_seconds)
        if self._terminal_retry_delay_seconds <= 0:
            raise ValueError("terminal_retry_delay_seconds must be positive")
        self._sleeper = sleeper
        self._lock = RLock()
        self._stopped = False
        self._active_slots: set[tuple[str, datetime]] = set()

    @property
    def accepting(self) -> bool:
        with self._lock:
            return not self._stopped

    def register(
        self,
        now: datetime,
        *,
        is_enabled: Callable[[ScheduledJobDefinition], bool] = lambda definition: True,
    ) -> tuple[ScheduledJobState, ...]:
        """Register all definitions before any candidate admission is attempted."""

        registrations = build_registrations(
            self._definitions, planner=self._planner, now=now, is_enabled=is_enabled)
        return tuple(self._control_plane.register(registrations))

    def recover_interrupted(self) -> tuple[ScheduledJobState, ...]:
        """Perform the one explicit recovery step before normal slot planning.

        Recovery is intentionally never implicit in :meth:`tick`: a caller
        that has not completed startup coordination must not mutate an
        authoritative ``RUNNING`` slot merely by polling for due work.
        """

        return tuple(self._control_plane.recover_interrupted())

    def stop(self) -> None:
        """Stop accepting new jobs while allowing an admitted slot to finish its commit chain."""

        with self._lock:
            self._stopped = True
        self._runtime.stop()

    def resume(self) -> None:
        """Resume admission after a planned restart exceeded its drain budget."""

        self._runtime.resume()
        with self._lock:
            self._stopped = False

    def tick(self, now: datetime) -> tuple[ExecutionOutcome, ...]:
        """Fetch authoritative state, run each due slot once, and report its terminal state."""

        _require_aware(now, "now")
        if not self.accepting:
            return ()
        states = self._control_plane.list_jobs()
        outcomes = []
        for slot in plan_due_slots(self._definitions, states, now=now):
            if not self._claim_local(slot):
                continue
            try:
                started_at = self._now()
                reserved_next_run_at = self._planner.reserve_next_due(
                    slot.definition,
                    slot_due_at=slot.scheduled_for,
                    started_at=started_at,
                )
                if not self.accepting or not self._control_plane.start(
                    slot.definition.id, slot.scheduled_for, reserved_next_run_at,
                ):
                    outcomes.append(ExecutionOutcome(slot.definition.id, slot.scheduled_for, ExecutionDisposition.START_REJECTED))
                    continue
                outcomes.append(self._run_admitted(slot))
            finally:
                self._release_local(slot)
        return tuple(outcomes)

    def _run_admitted(self, slot: ScheduledSlot) -> ExecutionOutcome:
        result = self._runtime.execute_admitted(
            slot.definition, slot.scheduled_for)

        finished_at = self._now()

        if result.status is JobResultStatus.SUCCEEDED:
            try:
                next_run_at = self._planner.next_due(
                    slot.definition,
                    slot_due_at=slot.scheduled_for,
                    completed_at=finished_at,
                    result=result,
                )
            except Exception as exc:
                return self._report_failure(
                    slot, finished_at, _permanent_failure_summary(exc))
            try:
                self._publisher.publish(slot.definition, result, slot.scheduled_for)
            except Exception as exc:
                # A Task/event acknowledgement is unknown.  Do not complete the
                # slot or advance its business cursor; replay through the stable
                # identity on the retry path instead.
                return self._report_failure(
                    slot, finished_at, _retryable_failure_summary(exc))
            self._commit_terminal(
                slot,
                lambda: self._control_plane.complete(
                    slot.definition.id, slot.scheduled_for, next_run_at),
            )
            return ExecutionOutcome(slot.definition.id, slot.scheduled_for, ExecutionDisposition.COMPLETED)

        return self._report_failure(slot, finished_at, result)

    def _report_failure(
        self,
        slot: ScheduledSlot,
        now: datetime,
        result: JobResult,
    ) -> ExecutionOutcome:
        retryable = result.status is JobResultStatus.RETRYABLE_FAILURE
        retry_at = now + timedelta(seconds=slot.definition.retry_delay_seconds) if retryable else None
        self._commit_terminal(
            slot,
            lambda: self._control_plane.fail(
                slot.definition.id,
                slot.scheduled_for,
                retryable=retryable,
                next_run_at=retry_at,
                error=result.error or "scheduled job failed without an error summary",
            ),
        )
        disposition = ExecutionDisposition.RETRYABLE_FAILURE if retryable else ExecutionDisposition.PERMANENT_FAILURE
        return ExecutionOutcome(slot.definition.id, slot.scheduled_for, disposition)

    def _commit_terminal(self, slot: ScheduledSlot, operation: Callable[[], None]) -> None:
        """Retry an ambiguous terminal transport failure without re-running the job."""

        while True:
            try:
                operation()
                return
            except Exception as exc:
                if not getattr(exc, "terminal_retryable", False):
                    raise
                _LOG.warning(
                    "scheduled job terminal update will retry job=%s slot=%s error=%s",
                    slot.definition.id,
                    slot.scheduled_for.isoformat(),
                    type(exc).__name__,
                )
                self._sleeper(self._terminal_retry_delay_seconds)

    def _now(self) -> datetime:
        value = self._clock()
        _require_aware(value, "clock result")
        return value

    def _claim_local(self, slot: ScheduledSlot) -> bool:
        key = (slot.definition.id, slot.scheduled_for)
        with self._lock:
            if self._stopped or key in self._active_slots:
                return False
            self._active_slots.add(key)
            return True

    def _release_local(self, slot: ScheduledSlot) -> None:
        with self._lock:
            self._active_slots.discard((slot.definition.id, slot.scheduled_for))


def _require_aware(value: datetime, name: str) -> None:
    if not is_aware(value):
        raise ValueError(f"{name} must be a timezone-aware datetime")


def _permanent_failure_summary(exc: Exception) -> JobResult:
    return _failure_summary(exc, JobResultStatus.PERMANENT_FAILURE)


def _retryable_failure_summary(exc: Exception) -> JobResult:
    return _failure_summary(exc, JobResultStatus.RETRYABLE_FAILURE)


def _failure_summary(exc: Exception, status: JobResultStatus) -> JobResult:
    message = " ".join(str(exc).split())
    summary = type(exc).__name__ if not message else f"{type(exc).__name__}: {message}"
    return JobResult(status, error=summary[:2048])
