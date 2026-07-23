"""Small immutable contracts for Bridge scheduled jobs."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
import math
import re
from typing import Any, Mapping, Optional, Union


class JobResultStatus(str, Enum):
    SUCCEEDED = "SUCCEEDED"
    RETRYABLE_FAILURE = "RETRYABLE_FAILURE"
    PERMANENT_FAILURE = "PERMANENT_FAILURE"


class MisfirePolicy(str, Enum):
    COALESCE = "COALESCE"
    CURRENT_DAY = "CURRENT_DAY"
    WAIT_FOR_COMPLETION = "WAIT_FOR_COMPLETION"


_JOB_ID_RE = re.compile(r"^[a-z][a-z0-9_-]*\.[a-z][a-z0-9_-]*$")


def _require_nonblank(value: Any, name: str) -> None:
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{name} must be a nonblank string")


def _require_positive_number(value: Any, name: str) -> None:
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise ValueError(f"{name} must be a finite positive number")
    if not math.isfinite(float(value)) or value <= 0:
        raise ValueError(f"{name} must be a finite positive number")


def is_aware(value: datetime) -> bool:
    return (
        isinstance(value, datetime)
        and value.tzinfo is not None
        and value.utcoffset() is not None
    )


@dataclass(frozen=True)
class IntervalSchedule:
    interval_seconds: float
    run_immediately: bool

    def __post_init__(self) -> None:
        _require_positive_number(self.interval_seconds, "interval_seconds")
        if not isinstance(self.run_immediately, bool):
            raise ValueError("run_immediately must be a bool")


@dataclass(frozen=True)
class DailySchedule:
    hour: int
    minute: int
    timezone: str = "Asia/Shanghai"

    def __post_init__(self) -> None:
        if (
            isinstance(self.hour, bool)
            or not isinstance(self.hour, int)
            or not 0 <= self.hour <= 23
        ):
            raise ValueError("hour must be an integer from 0 through 23")
        if (
            isinstance(self.minute, bool)
            or not isinstance(self.minute, int)
            or not 0 <= self.minute <= 59
        ):
            raise ValueError("minute must be an integer from 0 through 59")
        if self.timezone != "Asia/Shanghai":
            raise ValueError("timezone must be exactly Asia/Shanghai")


@dataclass(frozen=True)
class AdaptiveSchedule:
    min_delay_seconds: float
    default_delay_seconds: float
    max_delay_seconds: float
    run_immediately: bool

    def __post_init__(self) -> None:
        _require_positive_number(self.min_delay_seconds, "min_delay_seconds")
        _require_positive_number(self.default_delay_seconds, "default_delay_seconds")
        _require_positive_number(self.max_delay_seconds, "max_delay_seconds")
        if not (
            self.min_delay_seconds
            <= self.default_delay_seconds
            <= self.max_delay_seconds
        ):
            raise ValueError("adaptive delays must satisfy min <= default <= max")
        if not isinstance(self.run_immediately, bool):
            raise ValueError("run_immediately must be a bool")


@dataclass(frozen=True)
class HandlerRunner:
    handler_key: str
    payload: Mapping[str, Any] = field(default_factory=dict)

    def __post_init__(self) -> None:
        _require_nonblank(self.handler_key, "handler_key")
        if not isinstance(self.payload, Mapping):
            raise TypeError("payload must be a mapping")
        object.__setattr__(self, "payload", dict(self.payload))


ScheduleDefinition = Union[IntervalSchedule, DailySchedule, AdaptiveSchedule]


@dataclass(frozen=True)
class ScheduledJobDefinition:
    id: str
    revision: int
    description: str
    schedule: ScheduleDefinition
    runner: HandlerRunner
    misfire: MisfirePolicy
    retry_delay_seconds: float
    enabled: bool = True

    def __post_init__(self) -> None:
        validate_job_definition(self)


@dataclass(frozen=True)
class JobResult:
    status: JobResultStatus
    next_due_at: Optional[datetime] = None
    error: Optional[str] = None

    def __post_init__(self) -> None:
        if not isinstance(self.status, JobResultStatus):
            raise ValueError("status must be a JobResultStatus")
        if self.next_due_at is not None and not is_aware(self.next_due_at):
            raise ValueError("next_due_at must be timezone-aware")
        if self.status is JobResultStatus.SUCCEEDED:
            if self.error is not None:
                raise ValueError("SUCCEEDED results must not have an error")
        else:
            _require_nonblank(self.error, "error")
            if self.next_due_at is not None:
                raise ValueError("failure results must not carry next_due_at")


def validate_job_definition(definition: ScheduledJobDefinition) -> None:
    if type(definition) is not ScheduledJobDefinition:
        raise TypeError("definition must be a ScheduledJobDefinition")
    if not isinstance(definition.id, str) or not _JOB_ID_RE.fullmatch(definition.id):
        raise ValueError("id must match '<domain>.<name>'")
    if (
        isinstance(definition.revision, bool)
        or not isinstance(definition.revision, int)
        or definition.revision <= 0
    ):
        raise ValueError("revision must be a positive integer")
    _require_nonblank(definition.description, "description")
    if type(definition.schedule) not in (
        IntervalSchedule,
        DailySchedule,
        AdaptiveSchedule,
    ):
        raise TypeError("schedule must be exactly one supported schedule dataclass")
    if type(definition.runner) is not HandlerRunner:
        raise TypeError("runner must be a HandlerRunner")
    if not isinstance(definition.misfire, MisfirePolicy):
        raise ValueError("misfire must be a MisfirePolicy")
    allowed_misfires = {
        IntervalSchedule: frozenset({MisfirePolicy.COALESCE}),
        DailySchedule: frozenset({MisfirePolicy.CURRENT_DAY}),
        AdaptiveSchedule: frozenset({
            MisfirePolicy.COALESCE,
            MisfirePolicy.WAIT_FOR_COMPLETION,
        }),
    }
    if definition.misfire not in allowed_misfires[type(definition.schedule)]:
        raise ValueError(
            "%s schedule does not support misfire %s"
            % (type(definition.schedule).__name__, definition.misfire.value))
    _require_positive_number(
        definition.retry_delay_seconds, "retry_delay_seconds")
    if not isinstance(definition.enabled, bool):
        raise ValueError("enabled must be a bool")


def definition_snapshot(definition: ScheduledJobDefinition) -> dict[str, Any]:
    validate_job_definition(definition)
    schedule = definition.schedule
    if isinstance(schedule, IntervalSchedule):
        schedule_snapshot = {
            "kind": "interval",
            "interval_seconds": schedule.interval_seconds,
            "run_immediately": schedule.run_immediately,
        }
    elif isinstance(schedule, DailySchedule):
        schedule_snapshot = {
            "kind": "daily",
            "hour": schedule.hour,
            "minute": schedule.minute,
            "timezone": schedule.timezone,
        }
    else:
        schedule_snapshot = {
            "kind": "adaptive",
            "min_delay_seconds": schedule.min_delay_seconds,
            "default_delay_seconds": schedule.default_delay_seconds,
            "max_delay_seconds": schedule.max_delay_seconds,
            "run_immediately": schedule.run_immediately,
        }
    runner_snapshot = {
        "kind": "handler",
        "handler_key": definition.runner.handler_key,
        "payload": dict(definition.runner.payload),
    }
    return {
        "id": definition.id,
        "revision": definition.revision,
        "description": definition.description,
        "schedule": schedule_snapshot,
        "runner": runner_snapshot,
        "misfire": definition.misfire.value,
        "retry_delay_seconds": definition.retry_delay_seconds,
        "enabled": definition.enabled,
    }
__all__ = [
    "AdaptiveSchedule",
    "DailySchedule",
    "HandlerRunner",
    "IntervalSchedule",
    "JobResult",
    "JobResultStatus",
    "MisfirePolicy",
    "ScheduleDefinition",
    "ScheduledJobDefinition",
    "definition_snapshot",
    "is_aware",
    "validate_job_definition",
]
