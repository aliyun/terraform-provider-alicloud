"""Versioned YAML registry for Scheduler job ownership and runner routing."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, Mapping, Optional

try:
    import yaml
except ImportError as exc:  # fail closed with an actionable install instruction
    raise RuntimeError(
        "PyYAML is required for Scheduler registry loading; install bridge/requirements.txt"
    ) from exc


class SchedulerRegistryError(ValueError):
    """The checked-in scheduler ownership manifest is invalid."""


@dataclass(frozen=True)
class SchedulerJobRoute:
    job_key: str
    owner: str
    runner: Optional[str] = None

    @property
    def scheduler_owned(self) -> bool:
        return self.owner == "scheduler"


@dataclass(frozen=True)
class SchedulerRegistry:
    """Validated complete ownership snapshot used by both Bridge paths."""

    routes: tuple[SchedulerJobRoute, ...]

    def route_for(self, job_key: str) -> SchedulerJobRoute:
        for route in self.routes:
            if route.job_key == job_key:
                return route
        raise SchedulerRegistryError("job is not declared in scheduler registry: %s" % job_key)

    def scheduler_job_keys(self) -> frozenset[str]:
        return frozenset(route.job_key for route in self.routes if route.scheduler_owned)

    def runner_for(self, job_key: str) -> str:
        route = self.route_for(job_key)
        if not route.scheduler_owned or not route.runner:
            raise SchedulerRegistryError("job is not scheduler-owned: %s" % job_key)
        return route.runner


def default_registry_path() -> Path:
    return Path(__file__).resolve().parents[2] / "config" / "scheduler-jobs.yaml"


def load_scheduler_registry(
    *, path: Optional[Path] = None, known_job_keys: Iterable[str],
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
    if not isinstance(rows, list) or not rows:
        raise SchedulerRegistryError("scheduler registry jobs must be a nonempty list: %s" % target)
    routes = []
    seen = set()
    for index, row in enumerate(rows, start=1):
        if not isinstance(row, Mapping) or set(row).difference({"key", "owner", "runner"}):
            raise SchedulerRegistryError("invalid job entry #%s in %s" % (index, target))
        key = _required_text(row, "key", index, target)
        owner = _required_text(row, "owner", index, target)
        runner = _optional_text(row.get("runner"), index, target)
        if key in seen:
            raise SchedulerRegistryError("duplicate job key in %s: %s" % (target, key))
        seen.add(key)
        if owner not in {"legacy", "scheduler"}:
            raise SchedulerRegistryError("job %s owner must be legacy or scheduler" % key)
        if owner == "scheduler" and not runner:
            raise SchedulerRegistryError("scheduler-owned job requires runner: %s" % key)
        if owner == "legacy" and runner is not None:
            raise SchedulerRegistryError("legacy job cannot declare runner: %s" % key)
        routes.append(SchedulerJobRoute(key, owner, runner))
    known = frozenset(str(key) for key in known_job_keys)
    declared = frozenset(route.job_key for route in routes)
    missing = sorted(known.difference(declared))
    unknown = sorted(declared.difference(known))
    if missing or unknown:
        details = []
        if missing:
            details.append("missing=" + ",".join(missing))
        if unknown:
            details.append("unknown=" + ",".join(unknown))
        raise SchedulerRegistryError("scheduler registry must declare every Job: " + " ".join(details))
    return SchedulerRegistry(tuple(routes))


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


__all__ = [
    "SchedulerJobRoute", "SchedulerRegistry", "SchedulerRegistryError",
    "default_registry_path", "load_scheduler_registry",
]
