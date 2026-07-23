"""Load and validate the Scheduler job registry and its runner capabilities."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
from typing import Any, Mapping, Optional

try:
    import yaml
except ImportError as exc:  # fail closed with an actionable install instruction
    raise RuntimeError(
        "PyYAML is required for Scheduler registry loading; install bridge/requirements.txt"
    ) from exc


class SchedulerRegistryError(ValueError):
    """The checked-in scheduler ownership manifest is invalid."""


@dataclass(frozen=True)
class SchedulerRegistry:
    """Validated complete ownership snapshot for the standalone Scheduler."""

    definitions: tuple["ScheduledJobDefinition", ...]

    def scheduler_job_keys(self) -> frozenset[str]:
        return frozenset(definition.id for definition in self.definitions)

    def handler_keys(self) -> frozenset[str]:
        from .model import HandlerRunner

        return frozenset(
            definition.runner.handler_key
            for definition in self.definitions
            if isinstance(definition.runner, HandlerRunner)
        )

    def headless_builder_protocols(self) -> frozenset[tuple[str, str]]:
        from .model import HeadlessRunner

        return frozenset(
            (definition.runner.builder_ref, definition.runner.protocol)
            for definition in self.definitions
            if isinstance(definition.runner, HeadlessRunner)
        )

    def runner_keys(self) -> frozenset[object]:
        return self.handler_keys() | self.headless_builder_protocols()


def default_registry_path() -> Path:
    return Path(__file__).resolve().parent / "jobs.yaml"


def load_scheduler_registry(
    *, path: Optional[Path] = None,
) -> SchedulerRegistry:
    """Load and validate the complete checked-in scheduler registry."""

    target = path or default_registry_path()
    try:
        raw = yaml.safe_load(target.read_text(encoding="utf-8"))
    except OSError as exc:
        raise SchedulerRegistryError("cannot read scheduler registry %s: %s" % (target, exc)) from exc
    except yaml.YAMLError as exc:
        raise SchedulerRegistryError("invalid scheduler registry YAML %s: %s" % (target, exc)) from exc
    if not isinstance(raw, Mapping) or set(raw) != {"version", "jobs"} or raw["version"] != 1:
        raise SchedulerRegistryError("scheduler registry requires only version: 1 and jobs: %s" % target)
    rows = raw["jobs"]
    if not isinstance(rows, list):
        raise SchedulerRegistryError("scheduler registry jobs must be a list: %s" % target)
    definitions = []
    seen = set()
    for index, row in enumerate(rows, start=1):
        if not isinstance(row, Mapping):
            raise SchedulerRegistryError("invalid job entry #%s in %s" % (index, target))
        key = _required_text(row, "key", index, target)
        if key in seen:
            raise SchedulerRegistryError("duplicate job key in %s: %s" % (target, key))
        seen.add(key)
        definition = _definition_from_row(row, index, target)
        if definition.id != key:
            raise SchedulerRegistryError("job key/id mismatch in %s: %s" % (target, key))
        definitions.append(definition)
    return SchedulerRegistry(tuple(definitions))


def _required_text(row: Mapping[str, object], field: str, index: int, path: Path) -> str:
    value = _optional_text(row.get(field), index, path)
    if not value:
        raise SchedulerRegistryError("job #%s is missing %s in %s" % (index, field, path))
    return value


def _optional_text(value: object, index: int, path: Path) -> Optional[str]:
    if value is None:
        return None
    if not isinstance(value, str) or not value.strip():
        raise SchedulerRegistryError("job #%s has invalid text value in %s" % (index, path))
    return value.strip()


def _definition_from_row(
    row: Mapping[str, Any], index: int, path: Path,
) -> "ScheduledJobDefinition":
    from .model import MisfirePolicy, ScheduledJobDefinition
    expected = {
        "key", "revision", "description", "schedule", "runner", "misfire",
        "retry_delay_seconds", "enabled",
    }
    unknown = set(row).difference(expected)
    if unknown:
        raise SchedulerRegistryError("unknown job fields in %s #%s: %s" % (
            path, index, ",".join(sorted(str(name) for name in unknown))))
    try:
        return ScheduledJobDefinition(
            id=_required_text(row, "key", index, path),
            revision=_required_int(row, "revision", index, path),
            description=_required_text(row, "description", index, path),
            schedule=_schedule_from_yaml(row.get("schedule"), index, path),
            runner=_runner_from_yaml(row.get("runner"), index, path),
            misfire=MisfirePolicy(_required_text(row, "misfire", index, path)),
            retry_delay_seconds=_required_number(row, "retry_delay_seconds", index, path),
            enabled=_required_bool(row, "enabled", index, path),
        )
    except (TypeError, ValueError) as exc:
        raise SchedulerRegistryError("invalid job definition #%s in %s: %s" % (
            index, path, exc)) from exc


def _schedule_from_yaml(value: object, index: int, path: Path) -> object:
    from .model import AdaptiveSchedule, DailySchedule, IntervalSchedule
    if not isinstance(value, Mapping):
        raise SchedulerRegistryError("job #%s schedule must be a mapping in %s" % (index, path))
    kind = _required_text(value, "kind", index, path)
    if kind == "interval":
        _require_exact_fields(value, {"kind", "interval_seconds", "run_immediately"}, "schedule", index, path)
        return IntervalSchedule(value["interval_seconds"], value["run_immediately"])
    if kind == "daily":
        _require_exact_fields(value, {"kind", "hour", "minute", "timezone"}, "schedule", index, path)
        return DailySchedule(value["hour"], value["minute"], value["timezone"])
    if kind == "adaptive":
        _require_exact_fields(value, {"kind", "min_delay_seconds", "default_delay_seconds", "max_delay_seconds", "run_immediately"}, "schedule", index, path)
        return AdaptiveSchedule(value["min_delay_seconds"], value["default_delay_seconds"], value["max_delay_seconds"], value["run_immediately"])
    raise SchedulerRegistryError("job #%s has unsupported schedule kind %s" % (index, kind))


def _runner_from_yaml(value: object, index: int, path: Path) -> object:
    from .model import HandlerRunner, HeadlessRunner
    if not isinstance(value, Mapping):
        raise SchedulerRegistryError("job #%s runner must be a mapping in %s" % (index, path))
    kind = _required_text(value, "kind", index, path)
    if kind == "handler":
        _require_exact_fields(value, {"kind", "handler_key", "payload"}, "runner", index, path, optional={"payload"})
        return HandlerRunner(value["handler_key"], value.get("payload", {}))
    if kind == "headless":
        _require_exact_fields(value, {"kind", "builder_ref", "protocol", "session_policy", "lane", "model"}, "runner", index, path, optional={"lane", "model"})
        return HeadlessRunner(value["builder_ref"], value["protocol"], value["session_policy"], value.get("lane"), value.get("model"))
    raise SchedulerRegistryError("job #%s has unsupported runner kind %s" % (index, kind))


def _require_exact_fields(value: Mapping[str, Any], required: set[str], name: str,
                          index: int, path: Path, optional: set[str] | None = None) -> None:
    optional = optional or set()
    observed = set(value)
    missing = required.difference(optional).difference(observed)
    unknown = observed.difference(required)
    if missing or unknown:
        raise SchedulerRegistryError("job #%s invalid %s fields in %s" % (index, name, path))


def _required_int(row: Mapping[str, Any], field: str, index: int, path: Path) -> int:
    value = row.get(field)
    if isinstance(value, bool) or not isinstance(value, int):
        raise SchedulerRegistryError("job #%s requires integer %s in %s" % (index, field, path))
    return value


def _required_number(row: Mapping[str, Any], field: str, index: int, path: Path) -> float | int:
    value = row.get(field)
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise SchedulerRegistryError("job #%s requires number %s in %s" % (index, field, path))
    return value


def _required_bool(row: Mapping[str, Any], field: str, index: int, path: Path) -> bool:
    value = row.get(field)
    if not isinstance(value, bool):
        raise SchedulerRegistryError("job #%s requires boolean %s in %s" % (index, field, path))
    return value


from .model import CapabilityValidationContext, ScheduledJobDefinition, validate_registry
from .runners import HANDLER_KEYS, HEADLESS_BUILDER_PROTOCOLS


REGISTRY = load_scheduler_registry()
JOBS: tuple[ScheduledJobDefinition, ...] = validate_registry(
    REGISTRY.definitions,
    context=CapabilityValidationContext(
        handler_keys=HANDLER_KEYS,
        headless_builder_protocols=HEADLESS_BUILDER_PROTOCOLS,
    ),
)


def load_jobs() -> tuple[ScheduledJobDefinition, ...]:
    """Return the already validated immutable registry snapshot."""

    return JOBS


__all__ = [
    "JOBS",
    "REGISTRY",
    "SchedulerRegistry",
    "SchedulerRegistryError",
    "default_registry_path",
    "load_jobs",
    "load_scheduler_registry",
]
