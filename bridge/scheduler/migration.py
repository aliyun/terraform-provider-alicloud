"""Single-source routing between legacy loops and ``SchedulerEngine``.

``JARVIS_SCHEDULER_NEW_JOBS`` is the only ownership switch.  A listed job is
owned by the new Engine; every other known job remains legacy-owned.  Business
enablement is deliberately independent and continues to use each definition's
``enabled_env`` switch.
"""

from __future__ import annotations

import os
from typing import Iterable, Mapping


SCHEDULER_NEW_JOBS_ENV = "JARVIS_SCHEDULER_NEW_JOBS"
_DEPRECATED_ENABLE_ENV = "JARVIS_SCHEDULER_ENABLE"
_DEPRECATED_JOB_ENV_PREFIX = "JARVIS_SCHEDULER_JOB_"


class SchedulerMigrationError(ValueError):
    """The operator supplied an ambiguous or unsupported cutover configuration."""


def requested_new_jobs(
    job_keys: Iterable[str], *, environ: Mapping[str, str] | None = None,
) -> frozenset[str]:
    """Return the jobs selected for the new Engine, failing closed on typos.

    An empty value means every job remains on its legacy loop.  Deprecated
    per-job/global switches are rejected instead of being silently combined
    with the single-list contract.
    """

    env = os.environ if environ is None else environ
    _reject_deprecated_routes(env)
    known = frozenset(str(key) for key in job_keys)
    raw = str(env.get(SCHEDULER_NEW_JOBS_ENV, "")).strip()
    if not raw:
        return frozenset()
    values = [value.strip() for value in raw.split(",")]
    if any(not value for value in values):
        raise SchedulerMigrationError(
            "%s must be a comma-separated list of nonblank job keys"
            % SCHEDULER_NEW_JOBS_ENV)
    seen = set()
    duplicates = set()
    for value in values:
        if value in seen:
            duplicates.add(value)
        seen.add(value)
    if duplicates:
        raise SchedulerMigrationError(
            "%s contains duplicate scheduled jobs: %s"
            % (SCHEDULER_NEW_JOBS_ENV, ", ".join(sorted(duplicates))))
    unknown = sorted(set(values).difference(known))
    if unknown:
        raise SchedulerMigrationError(
            "%s selects unknown scheduled jobs: %s"
            % (SCHEDULER_NEW_JOBS_ENV, ", ".join(unknown)))
    return frozenset(values)


def uses_new_engine(job_key: str, *, job_keys: Iterable[str],
                    environ: Mapping[str, str] | None = None) -> bool:
    """Whether this process must suppress the legacy owner for ``job_key``."""

    return job_key in requested_new_jobs(job_keys, environ=environ)


def business_job_enabled(definition: object, *,
                         environ: Mapping[str, str] | None = None) -> bool:
    """Evaluate one definition's independent business enable switch.

    Definitions without ``enabled_env`` are enabled.  Configured switches use
    an explicit 0/1 contract so a malformed value cannot accidentally enable a
    newly routed job.
    """

    env = os.environ if environ is None else environ
    name = getattr(definition, "enabled_env", None)
    if name is None:
        return True
    value = str(env.get(name, "1")).strip()
    if value == "1":
        return True
    if value == "0":
        return False
    raise SchedulerMigrationError("%s must be 0 or 1" % name)


def _reject_deprecated_routes(environ: Mapping[str, str]) -> None:
    deprecated = []
    if str(environ.get(_DEPRECATED_ENABLE_ENV, "")).strip():
        deprecated.append(_DEPRECATED_ENABLE_ENV)
    deprecated.extend(sorted(
        str(name) for name, value in environ.items()
        if str(name).startswith(_DEPRECATED_JOB_ENV_PREFIX)
        and str(value).strip()
    ))
    if deprecated:
        raise SchedulerMigrationError(
            "deprecated scheduler route variables are not supported; use %s: %s"
            % (SCHEDULER_NEW_JOBS_ENV, ", ".join(deprecated)))


__all__ = [
    "SCHEDULER_NEW_JOBS_ENV", "SchedulerMigrationError",
    "business_job_enabled", "requested_new_jobs", "uses_new_engine",
]
