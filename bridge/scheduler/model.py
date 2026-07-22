"""Immutable contracts for the Bridge scheduled-job subsystem.

This module deliberately has no Bridge runtime imports.  It is safe to import in
tests, command-line tools, and the eventual Scheduler Worker before any network
or process resources are available.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
import hashlib
import json
import math
import re
from typing import Any, Iterable, Mapping, Optional, Union


class ScheduleKind(str, Enum):
    INTERVAL = "interval"
    DAILY = "daily"
    ADAPTIVE = "adaptive"


class RunnerKind(str, Enum):
    COMMAND = "command"
    HEADLESS = "headless"
    HANDLER = "handler"


class JobPurpose(str, Enum):
    DISCOVERY = "DISCOVERY"
    MAINTENANCE = "MAINTENANCE"


class ReplayPolicy(str, Enum):
    TASK_UPSERT_IDEMPOTENT = "TASK_UPSERT_IDEMPOTENT"
    EVENT_LEDGER_IDEMPOTENT = "EVENT_LEDGER_IDEMPOTENT"
    EPHEMERAL = "EPHEMERAL"


class JobResultStatus(str, Enum):
    SUCCEEDED = "SUCCEEDED"
    RETRYABLE_FAILURE = "RETRYABLE_FAILURE"
    PERMANENT_FAILURE = "PERMANENT_FAILURE"


class MisfirePolicy(str, Enum):
    COALESCE = "COALESCE"
    CURRENT_DAY = "CURRENT_DAY"
    WAIT_FOR_COMPLETION = "WAIT_FOR_COMPLETION"


class CheckpointUpgradePolicy(str, Enum):
    """How a persisted checkpoint is handled after a definition revision change."""

    RESET_FULL = "RESET_FULL"
    RESET_OVERLAP = "RESET_OVERLAP"
    MIGRATE = "MIGRATE"


_JOB_ID_RE = re.compile(r"^[a-z][a-z0-9_-]*\.[a-z][a-z0-9_-]*$")
_WORKSPACE_KEY_RE = re.compile(r"^[a-z][a-z0-9_]*$")
_ENV_NAME_RE = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*$")


class _FrozenDict(dict):
    """JSON-compatible mapping that remains immutable after construction."""

    def _immutable(self, *args: Any, **kwargs: Any) -> None:
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


def _require_enum(value: Any, enum_type: type[Enum], name: str) -> None:
    if not isinstance(value, enum_type):
        raise ValueError(f"{name} must be a {enum_type.__name__}")


def _require_positive_number(value: Any, name: str) -> None:
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise ValueError(f"{name} must be a finite positive number")
    if not math.isfinite(float(value)) or value <= 0:
        raise ValueError(f"{name} must be a finite positive number")


def _require_positive_int(value: Any, name: str) -> None:
    if isinstance(value, bool) or not isinstance(value, int) or value <= 0:
        raise ValueError(f"{name} must be a positive integer")


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


def _freeze_strings(values: Iterable[str], name: str) -> frozenset[str]:
    if isinstance(values, (str, bytes)):
        raise TypeError(f"{name} must be an iterable of nonblank strings")
    try:
        materialized = tuple(values)
    except TypeError as exc:
        raise TypeError(f"{name} must be an iterable of nonblank strings") from exc
    for value in materialized:
        _require_nonblank(value, name)
    return frozenset(materialized)


def is_aware(value: datetime) -> bool:
    return isinstance(value, datetime) and value.tzinfo is not None and value.utcoffset() is not None


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
        if isinstance(self.hour, bool) or not isinstance(self.hour, int) or not 0 <= self.hour <= 23:
            raise ValueError("hour must be an integer from 0 through 23")
        if isinstance(self.minute, bool) or not isinstance(self.minute, int) or not 0 <= self.minute <= 59:
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
        if not self.min_delay_seconds <= self.default_delay_seconds <= self.max_delay_seconds:
            raise ValueError("adaptive delays must satisfy min <= default <= max")
        if not isinstance(self.run_immediately, bool):
            raise ValueError("run_immediately must be a bool")


@dataclass(frozen=True)
class CommandRunner:
    argv: tuple[str, ...]
    workspace_key: str
    env_allowlist: tuple[str, ...] = ()

    def __post_init__(self) -> None:
        if isinstance(self.argv, (str, bytes)):
            raise TypeError("argv must be an argument sequence, not a shell string")
        try:
            argv = tuple(self.argv)
        except TypeError as exc:
            raise TypeError("argv must be an argument sequence") from exc
        if not argv or any(not isinstance(arg, str) or not arg for arg in argv):
            raise ValueError("argv must contain nonempty string arguments")
        if not isinstance(self.workspace_key, str) or not _WORKSPACE_KEY_RE.fullmatch(self.workspace_key):
            raise ValueError("workspace_key must be a workspace-key-shaped string")
        if isinstance(self.env_allowlist, (str, bytes)):
            raise TypeError("env_allowlist must be an iterable of environment names")
        env_allowlist = tuple(self.env_allowlist)
        if any(not isinstance(name, str) or not _ENV_NAME_RE.fullmatch(name) for name in env_allowlist):
            raise ValueError("env_allowlist must contain valid environment variable names")
        object.__setattr__(self, "argv", argv)
        object.__setattr__(self, "env_allowlist", env_allowlist)


@dataclass(frozen=True)
class HeadlessRunner:
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


@dataclass(frozen=True)
class HandlerRunner:
    handler_key: str
    payload: Any = field(default_factory=dict)

    def __post_init__(self) -> None:
        _require_nonblank(self.handler_key, "handler_key")
        object.__setattr__(self, "payload", _freeze_json(self.payload, "payload"))


ScheduleDefinition = Union[IntervalSchedule, DailySchedule, AdaptiveSchedule]
RunnerDefinition = Union[CommandRunner, HeadlessRunner, HandlerRunner]


@dataclass(frozen=True)
class ScheduledJobDefinition:
    id: str
    revision: int
    description: str
    purpose: JobPurpose
    schedule: ScheduleDefinition
    runner: RunnerDefinition
    misfire: MisfirePolicy
    timeout_seconds: float
    retry_delay_seconds: float
    replay_policy: ReplayPolicy
    enabled: bool = True
    checkpoint_upgrade: CheckpointUpgradePolicy = CheckpointUpgradePolicy.RESET_FULL
    checkpoint_upgrader_key: Optional[str] = None

    def __post_init__(self) -> None:
        validate_job_definition(self)


@dataclass(frozen=True)
class TaskObservation:
    task_key: str
    source_type: str
    source_ref: Mapping[str, Any]
    task_type: str
    desired_revision: str
    trigger_mask: tuple[str, ...]
    payload: Mapping[str, Any]
    recovery_policy: str = "REPLAY_SAFE"
    persona: Optional[str] = None
    priority: Any = None
    aone_id: Optional[str] = None
    comment_cursor: Any = None
    required_capabilities: Any = None
    max_retries: Optional[int] = None
    source_status: Optional[str] = None

    def __post_init__(self) -> None:
        for name in ("task_key", "source_type", "task_type", "desired_revision", "recovery_policy"):
            _require_nonblank(getattr(self, name), name)
        if not isinstance(self.source_ref, Mapping):
            raise TypeError("source_ref must be a mapping")
        if not isinstance(self.payload, Mapping):
            raise TypeError("payload must be a mapping")
        if isinstance(self.trigger_mask, (str, bytes)):
            raise TypeError("trigger_mask must be an iterable of nonblank strings")
        try:
            trigger_mask = tuple(self.trigger_mask)
        except TypeError as exc:
            raise TypeError("trigger_mask must be an iterable of nonblank strings") from exc
        for trigger in trigger_mask:
            _require_nonblank(trigger, "trigger_mask")
        if self.persona is not None:
            _require_nonblank(self.persona, "persona")
        if self.aone_id is not None:
            _require_nonblank(self.aone_id, "aone_id")
        if self.max_retries is not None:
            if isinstance(self.max_retries, bool) or not isinstance(self.max_retries, int) or self.max_retries < 0:
                raise ValueError("max_retries must be a non-negative integer")
        source_status = None if self.source_status is None else str(self.source_status).strip()
        if source_status == "":
            source_status = None
        if source_status is not None and len(source_status) > 64:
            raise ValueError("source_status must not exceed 64 characters")
        for name in ("source_ref", "payload", "priority", "comment_cursor", "required_capabilities"):
            object.__setattr__(self, name, _freeze_json(getattr(self, name), name))
        object.__setattr__(self, "trigger_mask", trigger_mask)
        object.__setattr__(self, "source_status", source_status)


@dataclass(frozen=True)
class JobResult:
    status: JobResultStatus
    observations: tuple[TaskObservation, ...] = ()
    checkpoint: Any = None
    next_due_at: Optional[datetime] = None
    error: Optional[str] = None

    def __post_init__(self) -> None:
        _require_enum(self.status, JobResultStatus, "status")
        if isinstance(self.observations, (str, bytes)):
            raise TypeError("observations must be an iterable of TaskObservation")
        observations = tuple(self.observations)
        if not all(isinstance(item, TaskObservation) for item in observations):
            raise TypeError("observations must contain TaskObservation values")
        checkpoint = _freeze_json(self.checkpoint, "checkpoint")
        if self.next_due_at is not None and not is_aware(self.next_due_at):
            raise ValueError("next_due_at must be timezone-aware")
        if self.status is JobResultStatus.SUCCEEDED:
            if self.error is not None:
                raise ValueError("SUCCEEDED results must not have an error")
        else:
            _require_nonblank(self.error, "error")
            if observations or checkpoint is not None or self.next_due_at is not None:
                raise ValueError("failure results must not carry observations, checkpoint, or next_due_at")
        object.__setattr__(self, "observations", observations)
        object.__setattr__(self, "checkpoint", checkpoint)


@dataclass(frozen=True)
class CapabilityValidationContext:
    registered_workspace_keys: frozenset[str] = frozenset()
    headless_builder_protocols: frozenset[tuple[str, str]] = frozenset()
    handler_keys: frozenset[str] = frozenset()
    checkpoint_upgrader_keys: frozenset[str] = frozenset()

    def __post_init__(self) -> None:
        workspace_keys = _freeze_strings(self.registered_workspace_keys, "registered_workspace_keys")
        if any(not _WORKSPACE_KEY_RE.fullmatch(key) for key in workspace_keys):
            raise ValueError("registered_workspace_keys must contain workspace keys")
        pairs = tuple(self.headless_builder_protocols)
        normalized_pairs = set()
        for pair in pairs:
            if isinstance(pair, (str, bytes)) or not isinstance(pair, (tuple, list)) or len(pair) != 2:
                raise TypeError("headless_builder_protocols must contain (builder_ref, protocol) pairs")
            builder_ref, protocol = pair
            _require_nonblank(builder_ref, "builder_ref")
            _require_nonblank(protocol, "protocol")
            normalized_pairs.add((builder_ref, protocol))
        object.__setattr__(self, "registered_workspace_keys", workspace_keys)
        object.__setattr__(self, "headless_builder_protocols", frozenset(normalized_pairs))
        object.__setattr__(self, "handler_keys", _freeze_strings(self.handler_keys, "handler_keys"))
        object.__setattr__(
            self,
            "checkpoint_upgrader_keys",
            _freeze_strings(self.checkpoint_upgrader_keys, "checkpoint_upgrader_keys"),
        )


def validate_job_definition(definition: ScheduledJobDefinition) -> None:
    if type(definition) is not ScheduledJobDefinition:
        raise TypeError("definition must be a ScheduledJobDefinition")
    if not isinstance(definition.id, str) or not _JOB_ID_RE.fullmatch(definition.id):
        raise ValueError("id must match '<domain>.<name>'")
    _require_positive_int(definition.revision, "revision")
    _require_nonblank(definition.description, "description")
    _require_enum(definition.purpose, JobPurpose, "purpose")
    _require_enum(definition.misfire, MisfirePolicy, "misfire")
    _require_positive_number(definition.timeout_seconds, "timeout_seconds")
    _require_positive_number(definition.retry_delay_seconds, "retry_delay_seconds")
    _require_enum(definition.replay_policy, ReplayPolicy, "replay_policy")
    _require_enum(definition.checkpoint_upgrade, CheckpointUpgradePolicy, "checkpoint_upgrade")
    if not isinstance(definition.enabled, bool):
        raise ValueError("enabled must be a bool")
    if definition.checkpoint_upgrade in (
        CheckpointUpgradePolicy.RESET_OVERLAP,
        CheckpointUpgradePolicy.MIGRATE,
    ):
        _require_nonblank(definition.checkpoint_upgrader_key, "checkpoint_upgrader_key")
    elif definition.checkpoint_upgrader_key is not None:
        raise ValueError(
            "checkpoint_upgrader_key is valid only with RESET_OVERLAP or MIGRATE")
    if type(definition.schedule) not in (IntervalSchedule, DailySchedule, AdaptiveSchedule):
        raise TypeError("schedule must be exactly one supported schedule dataclass")
    if type(definition.runner) not in (CommandRunner, HeadlessRunner, HandlerRunner):
        raise TypeError("runner must be exactly one supported runner dataclass")
    if definition.purpose is JobPurpose.DISCOVERY and definition.replay_policy is ReplayPolicy.EVENT_LEDGER_IDEMPOTENT:
        raise ValueError("discovery jobs cannot use EVENT_LEDGER_IDEMPOTENT")
    if definition.purpose is JobPurpose.MAINTENANCE and definition.replay_policy is ReplayPolicy.TASK_UPSERT_IDEMPOTENT:
        raise ValueError("maintenance jobs cannot use TASK_UPSERT_IDEMPOTENT")


def _snapshot_json(value: Any) -> Any:
    if isinstance(value, Mapping):
        return {key: _snapshot_json(item) for key, item in value.items()}
    if isinstance(value, tuple):
        return [_snapshot_json(item) for item in value]
    return value


def definition_snapshot(definition: ScheduledJobDefinition) -> dict[str, Any]:
    validate_job_definition(definition)
    schedule = definition.schedule
    if type(schedule) is IntervalSchedule:
        schedule_snapshot = {"kind": ScheduleKind.INTERVAL.value, "interval_seconds": schedule.interval_seconds, "run_immediately": schedule.run_immediately}
    elif type(schedule) is DailySchedule:
        schedule_snapshot = {"kind": ScheduleKind.DAILY.value, "hour": schedule.hour, "minute": schedule.minute, "timezone": schedule.timezone}
    else:
        schedule_snapshot = {"kind": ScheduleKind.ADAPTIVE.value, "min_delay_seconds": schedule.min_delay_seconds, "default_delay_seconds": schedule.default_delay_seconds, "max_delay_seconds": schedule.max_delay_seconds, "run_immediately": schedule.run_immediately}
    runner = definition.runner
    if type(runner) is CommandRunner:
        runner_snapshot = {"kind": RunnerKind.COMMAND.value, "argv": list(runner.argv), "workspace_key": runner.workspace_key, "env_allowlist": list(runner.env_allowlist)}
    elif type(runner) is HeadlessRunner:
        runner_snapshot = {"kind": RunnerKind.HEADLESS.value, "builder_ref": runner.builder_ref, "protocol": runner.protocol, "session_policy": runner.session_policy, "lane": runner.lane, "model": runner.model}
    else:
        runner_snapshot = {"kind": RunnerKind.HANDLER.value, "handler_key": runner.handler_key, "payload": _snapshot_json(runner.payload)}
    return {
        "id": definition.id,
        "revision": definition.revision,
        "description": definition.description,
        "purpose": definition.purpose.value,
        "schedule": schedule_snapshot,
        "runner": runner_snapshot,
        "misfire": definition.misfire.value,
        "timeout_seconds": definition.timeout_seconds,
        "retry_delay_seconds": definition.retry_delay_seconds,
        "replay_policy": definition.replay_policy.value,
        "enabled": definition.enabled,
        "checkpoint_upgrade": definition.checkpoint_upgrade.value,
        "checkpoint_upgrader_key": definition.checkpoint_upgrader_key,
    }


def definition_digest(definition: ScheduledJobDefinition) -> str:
    serialized = json.dumps(definition_snapshot(definition), sort_keys=True, separators=(",", ":"), allow_nan=False)
    return hashlib.sha256(serialized.encode("utf-8")).hexdigest()


def validate_registry(
    definitions: Iterable[ScheduledJobDefinition], *, context: CapabilityValidationContext,
    expected_digests: Optional[Mapping[str, str]] = None,
) -> tuple[ScheduledJobDefinition, ...]:
    """Materialize once, validate in declared order, and return that tuple."""

    if not isinstance(context, CapabilityValidationContext):
        raise TypeError("context must be a CapabilityValidationContext")
    try:
        registry = tuple(definitions)
    except TypeError as exc:
        raise TypeError("definitions must be an iterable of ScheduledJobDefinition") from exc
    digests = dict(expected_digests) if expected_digests is not None else None
    seen_ids: set[str] = set()
    for definition in registry:
        validate_job_definition(definition)
        if definition.id in seen_ids:
            raise ValueError(f"duplicate scheduled job id: {definition.id}")
        seen_ids.add(definition.id)
        if digests is not None and digests.get(definition.id) != definition_digest(definition):
            raise ValueError(f"definition digest mismatch: {definition.id}")
        runner = definition.runner
        if isinstance(runner, CommandRunner) and runner.workspace_key not in context.registered_workspace_keys:
            raise ValueError(f"unregistered workspace key: {runner.workspace_key}")
        if isinstance(runner, HeadlessRunner) and (runner.builder_ref, runner.protocol) not in context.headless_builder_protocols:
            raise ValueError("unregistered headless builder protocol")
        if isinstance(runner, HandlerRunner) and runner.handler_key not in context.handler_keys:
            raise ValueError(f"unregistered handler key: {runner.handler_key}")
        if (
            definition.checkpoint_upgrade is not CheckpointUpgradePolicy.RESET_FULL
            and definition.checkpoint_upgrader_key not in context.checkpoint_upgrader_keys
        ):
            raise ValueError(
                f"unregistered checkpoint upgrader key: {definition.checkpoint_upgrader_key}")
    if digests is not None and set(digests) != seen_ids:
        raise ValueError("definition digest registry keys must exactly match job ids")
    return registry
