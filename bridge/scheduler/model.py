"""Small immutable contracts for Bridge scheduled jobs."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
import hashlib
import json
import math
import re
from typing import Any, Iterable, Mapping, Optional, Union


class JobResultStatus(str, Enum):
    SUCCEEDED = "SUCCEEDED"
    RETRYABLE_FAILURE = "RETRYABLE_FAILURE"
    PERMANENT_FAILURE = "PERMANENT_FAILURE"


class MisfirePolicy(str, Enum):
    COALESCE = "COALESCE"
    CURRENT_DAY = "CURRENT_DAY"
    WAIT_FOR_COMPLETION = "WAIT_FOR_COMPLETION"


_JOB_ID_RE = re.compile(r"^[a-z][a-z0-9_-]*\.[a-z][a-z0-9_-]*$")


class _FrozenDict(dict):
    """JSON-compatible mapping that remains immutable after construction."""

    def _immutable(self, *args: Any, **kwargs: Any) -> None:
        del args, kwargs
        raise TypeError("JSON value is immutable")

    __setitem__ = _immutable
    __delitem__ = _immutable
    clear = _immutable
    pop = _immutable
    popitem = _immutable
    setdefault = _immutable
    update = _immutable
    __ior__ = _immutable


def _require_nonblank(value: Any, name: str) -> None:
    if not isinstance(value, str) or not value.strip():
        raise ValueError(f"{name} must be a nonblank string")


def _require_positive_number(value: Any, name: str) -> None:
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise ValueError(f"{name} must be a finite positive number")
    if not math.isfinite(float(value)) or value <= 0:
        raise ValueError(f"{name} must be a finite positive number")


def _freeze_json(value: Any, name: str = "value") -> Any:
    if value is None or isinstance(value, (str, bool, int)):
        return value
    if isinstance(value, float):
        if not math.isfinite(value):
            raise ValueError(f"{name} must be JSON serializable")
        return value
    if isinstance(value, Mapping):
        items = []
        for key, item in value.items():
            if not isinstance(key, str):
                raise TypeError(f"{name} must have string object keys")
            items.append((key, _freeze_json(item, name)))
        return _FrozenDict(items)
    if isinstance(value, (list, tuple)):
        return tuple(_freeze_json(item, name) for item in value)
    raise TypeError(f"{name} must be JSON serializable")


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
class HeadlessRunner:
    """Fail-closed declaration for one registered headless builder protocol."""

    builder_ref: str
    protocol: str
    session_policy: str
    lane: Optional[str] = None
    model: Optional[str] = None

    def __post_init__(self) -> None:
        for name in ("builder_ref", "protocol", "session_policy"):
            _require_nonblank(getattr(self, name), name)
        for name in ("lane", "model"):
            value = getattr(self, name)
            if value is not None:
                _require_nonblank(value, name)
        if self.session_policy not in ("NEW", "RESUME"):
            raise ValueError("session_policy must be NEW or RESUME")
        if self.lane not in (None, "default", "terraform"):
            raise ValueError("lane must be default or terraform")
        if self.model is not None:
            raise ValueError(
                "model override is not supported by the Jarvis headless adapter")


@dataclass(frozen=True)
class HandlerRunner:
    handler_key: str
    payload: Any = field(default_factory=dict)

    def __post_init__(self) -> None:
        _require_nonblank(self.handler_key, "handler_key")
        object.__setattr__(self, "payload", _freeze_json(self.payload, "payload"))


ScheduleDefinition = Union[IntervalSchedule, DailySchedule, AdaptiveSchedule]
RunnerDefinition = Union[HeadlessRunner, HandlerRunner]


@dataclass(frozen=True)
class ScheduledJobDefinition:
    id: str
    revision: int
    description: str
    schedule: ScheduleDefinition
    runner: RunnerDefinition
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


@dataclass(frozen=True)
class CapabilityValidationContext:
    headless_builder_protocols: frozenset[tuple[str, str]] = frozenset()
    handler_keys: frozenset[str] = frozenset()

    def __post_init__(self) -> None:
        pairs = set()
        for pair in self.headless_builder_protocols:
            if (
                isinstance(pair, (str, bytes))
                or not isinstance(pair, (tuple, list))
                or len(pair) != 2
            ):
                raise TypeError(
                    "headless_builder_protocols must contain (builder_ref, protocol) pairs")
            builder_ref, protocol = pair
            _require_nonblank(builder_ref, "builder_ref")
            _require_nonblank(protocol, "protocol")
            pairs.add((builder_ref, protocol))
        handlers = frozenset(self.handler_keys)
        for handler_key in handlers:
            _require_nonblank(handler_key, "handler_keys")
        object.__setattr__(self, "headless_builder_protocols", frozenset(pairs))
        object.__setattr__(self, "handler_keys", handlers)


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
    if type(definition.runner) not in (HeadlessRunner, HandlerRunner):
        raise TypeError("runner must be exactly one supported runner dataclass")
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


def _snapshot_json(value: Any) -> Any:
    if isinstance(value, Mapping):
        return {key: _snapshot_json(item) for key, item in value.items()}
    if isinstance(value, tuple):
        return [_snapshot_json(item) for item in value]
    return value


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
    runner = definition.runner
    if isinstance(runner, HeadlessRunner):
        runner_snapshot = {
            "kind": "headless",
            "builder_ref": runner.builder_ref,
            "protocol": runner.protocol,
            "session_policy": runner.session_policy,
            "lane": runner.lane,
            "model": runner.model,
        }
    else:
        runner_snapshot = {
            "kind": "handler",
            "handler_key": runner.handler_key,
            "payload": _snapshot_json(runner.payload),
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


def definition_digest(definition: ScheduledJobDefinition) -> str:
    serialized = json.dumps(
        definition_snapshot(definition),
        sort_keys=True,
        separators=(",", ":"),
        allow_nan=False,
    )
    return hashlib.sha256(serialized.encode("utf-8")).hexdigest()


def validate_registry(
    definitions: Iterable[ScheduledJobDefinition],
    *,
    context: CapabilityValidationContext,
    expected_digests: Optional[Mapping[str, str]] = None,
) -> tuple[ScheduledJobDefinition, ...]:
    if not isinstance(context, CapabilityValidationContext):
        raise TypeError("context must be a CapabilityValidationContext")
    try:
        registry = tuple(definitions)
    except TypeError as exc:
        raise TypeError(
            "definitions must be an iterable of ScheduledJobDefinition") from exc
    digests = dict(expected_digests) if expected_digests is not None else None
    seen_ids: set[str] = set()
    for definition in registry:
        validate_job_definition(definition)
        if definition.id in seen_ids:
            raise ValueError(f"duplicate scheduled job id: {definition.id}")
        seen_ids.add(definition.id)
        if (
            digests is not None
            and digests.get(definition.id) != definition_digest(definition)
        ):
            raise ValueError(f"definition digest mismatch: {definition.id}")
        runner = definition.runner
        if (
            isinstance(runner, HeadlessRunner)
            and (runner.builder_ref, runner.protocol)
            not in context.headless_builder_protocols
        ):
            raise ValueError("unregistered headless builder protocol")
        if (
            isinstance(runner, HandlerRunner)
            and runner.handler_key not in context.handler_keys
        ):
            raise ValueError(f"unregistered handler key: {runner.handler_key}")
    if digests is not None and set(digests) != seen_ids:
        raise ValueError("definition digest registry keys must exactly match job ids")
    return registry


__all__ = [
    "AdaptiveSchedule",
    "CapabilityValidationContext",
    "DailySchedule",
    "HandlerRunner",
    "HeadlessRunner",
    "IntervalSchedule",
    "JobResult",
    "JobResultStatus",
    "MisfirePolicy",
    "RunnerDefinition",
    "ScheduleDefinition",
    "ScheduledJobDefinition",
    "definition_digest",
    "definition_snapshot",
    "is_aware",
    "validate_job_definition",
    "validate_registry",
]
