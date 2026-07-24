"""Load and strictly validate the checked-in Scheduler job registry."""

from __future__ import annotations

from pathlib import Path
from typing import Any, Mapping

try:
    import yaml
except ImportError as exc:
    raise RuntimeError(
        "PyYAML is required for Scheduler registry loading; "
        "run bootstrap/bridge-python.sh"
    ) from exc

from .model import (
    AdaptiveSchedule,
    DailySchedule,
    HandlerRunner,
    IntervalSchedule,
    MisfirePolicy,
    ScheduledJobDefinition,
)


JOBS_PATH = Path(__file__).with_name("jobs.yaml")
_TOP_LEVEL_FIELDS = frozenset({"version", "jobs"})
_JOB_FIELDS = frozenset({
    "key",
    "runner",
    "revision",
    "description",
    "schedule",
    "misfire",
    "retry_delay_seconds",
    "enabled",
})
_SCHEDULE_SPECS = {
    "interval": (
        IntervalSchedule,
        ("interval_seconds", "run_immediately"),
    ),
    "daily": (
        DailySchedule,
        ("hour", "minute", "timezone"),
    ),
    "adaptive": (
        AdaptiveSchedule,
        (
            "min_delay_seconds",
            "default_delay_seconds",
            "max_delay_seconds",
            "run_immediately",
        ),
    ),
}


class SchedulerJobsError(ValueError):
    """The checked-in Scheduler job registry is invalid."""


class _UniqueKeyLoader(yaml.SafeLoader):
    """Safe YAML loader that rejects duplicate mapping keys."""


def _construct_unique_mapping(
    loader: _UniqueKeyLoader,
    node: yaml.nodes.MappingNode,
    deep: bool = False,
) -> dict[object, object]:
    loader.flatten_mapping(node)
    result: dict[object, object] = {}
    for key_node, value_node in node.value:
        key = loader.construct_object(key_node, deep=deep)
        try:
            duplicate = key in result
        except TypeError as exc:
            raise yaml.constructor.ConstructorError(
                "while constructing a mapping",
                node.start_mark,
                "found an unhashable mapping key",
                key_node.start_mark,
            ) from exc
        if duplicate:
            raise yaml.constructor.ConstructorError(
                "while constructing a mapping",
                node.start_mark,
                "found duplicate key %r" % (key,),
                key_node.start_mark,
            )
        result[key] = loader.construct_object(value_node, deep=deep)
    return result


_UniqueKeyLoader.add_constructor(
    yaml.resolver.BaseResolver.DEFAULT_MAPPING_TAG,
    _construct_unique_mapping,
)


def load_jobs(path: Path = JOBS_PATH) -> tuple[ScheduledJobDefinition, ...]:
    """Return the immutable definitions loaded from ``jobs.yaml``."""

    try:
        raw = yaml.load(
            path.read_text(encoding="utf-8"),
            Loader=_UniqueKeyLoader,
        )
    except OSError as exc:
        raise SchedulerJobsError(
            "cannot read scheduler jobs %s: %s" % (path, exc)
        ) from exc
    except yaml.YAMLError as exc:
        raise SchedulerJobsError(
            "invalid scheduler jobs YAML %s: %s" % (path, exc)
        ) from exc

    root = _mapping(raw, "scheduler jobs root")
    _require_exact_fields(root, _TOP_LEVEL_FIELDS, "scheduler jobs root")
    if type(root["version"]) is not int or root["version"] != 1:
        raise SchedulerJobsError("scheduler jobs version must be exactly 1")

    rows = root["jobs"]
    if not isinstance(rows, list) or not rows:
        raise SchedulerJobsError("scheduler jobs must be a non-empty list")

    definitions: list[ScheduledJobDefinition] = []
    seen_keys: set[str] = set()
    for index, value in enumerate(rows, start=1):
        row = _mapping(value, "job #%d" % index)
        _require_exact_fields(row, _JOB_FIELDS, "job #%d" % index)
        key = _text(row["key"], "job #%d key" % index)
        if key in seen_keys:
            raise SchedulerJobsError("duplicate scheduler job key: %s" % key)
        seen_keys.add(key)
        definitions.append(_definition(row, index, key))

    handler_keys = tuple(item.runner.handler_key for item in definitions)
    if len(set(handler_keys)) != len(handler_keys):
        raise SchedulerJobsError("scheduler handler keys must be unique")
    return tuple(definitions)


def _definition(
    row: Mapping[str, Any],
    index: int,
    key: str,
) -> ScheduledJobDefinition:
    try:
        return ScheduledJobDefinition(
            id=key,
            revision=row["revision"],
            description=row["description"],
            schedule=_schedule(row["schedule"], index),
            runner=HandlerRunner(
                _text(row["runner"], "job #%d runner" % index)),
            misfire=MisfirePolicy(
                _text(row["misfire"], "job #%d misfire" % index)),
            retry_delay_seconds=row["retry_delay_seconds"],
            enabled=row["enabled"],
        )
    except SchedulerJobsError:
        raise
    except (TypeError, ValueError) as exc:
        raise SchedulerJobsError(
            "invalid scheduler job #%d (%s): %s" % (index, key, exc)
        ) from exc


def _schedule(value: object, index: int) -> object:
    schedule = _mapping(value, "job #%d schedule" % index)
    kind = _text(schedule.get("kind"), "job #%d schedule kind" % index)
    spec = _SCHEDULE_SPECS.get(kind)
    if spec is None:
        raise SchedulerJobsError(
            "job #%d has unsupported schedule kind: %s" % (index, kind)
        )
    constructor, fields = spec
    _require_exact_fields(
        schedule,
        frozenset(("kind", *fields)),
        "job #%d %s schedule" % (index, kind),
    )
    return constructor(*(schedule[field] for field in fields))


def _mapping(value: object, name: str) -> Mapping[str, Any]:
    if not isinstance(value, Mapping):
        raise SchedulerJobsError("%s must be a mapping" % name)
    if any(not isinstance(key, str) for key in value):
        raise SchedulerJobsError("%s keys must be strings" % name)
    return value


def _require_exact_fields(
    value: Mapping[str, Any],
    expected: frozenset[str],
    name: str,
) -> None:
    observed = set(value)
    missing = expected.difference(observed)
    unknown = observed.difference(expected)
    if missing or unknown:
        details = []
        if missing:
            details.append("missing %s" % ", ".join(sorted(missing)))
        if unknown:
            details.append("unknown %s" % ", ".join(sorted(unknown)))
        raise SchedulerJobsError(
            "%s has invalid fields: %s" % (name, "; ".join(details))
        )


def _text(value: object, name: str) -> str:
    if not isinstance(value, str) or not value.strip():
        raise SchedulerJobsError("%s must be a nonblank string" % name)
    return value.strip()


def runner_key(definition: ScheduledJobDefinition) -> str:
    return definition.runner.handler_key


JOBS = load_jobs()
JOB_KEYS = frozenset(definition.id for definition in JOBS)
RUNNER_KEYS = frozenset(runner_key(definition) for definition in JOBS)


__all__ = [
    "JOBS",
    "JOBS_PATH",
    "JOB_KEYS",
    "RUNNER_KEYS",
    "SchedulerJobsError",
    "load_jobs",
    "runner_key",
]
