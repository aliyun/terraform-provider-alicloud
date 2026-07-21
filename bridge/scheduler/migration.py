"""Explicit per-job cutover gates between legacy loops and ``SchedulerEngine``.

The migration is deliberately opt-in.  A job remains on its legacy loop until
both the global scheduler gate and that job's ``new`` route are set.  This
keeps mixed deployment safe: the old loop is still the default and a job never
has two owners in one Bridge process.
"""

from __future__ import annotations

import os
from typing import Iterable, Mapping


SCHEDULER_ENABLE_ENV = "JARVIS_SCHEDULER_ENABLE"
JOB_ENV_PREFIX = "JARVIS_SCHEDULER_JOB_"


class SchedulerMigrationError(ValueError):
    """The operator supplied an ambiguous or unsupported job cutover route."""


def scheduler_enabled(environ: Mapping[str, str] | None = None) -> bool:
    env = os.environ if environ is None else environ
    value = str(env.get(SCHEDULER_ENABLE_ENV, "0")).strip()
    if value in ("", "0"):
        return False
    if value == "1":
        return True
    raise SchedulerMigrationError(
        "%s must be 0 or 1" % SCHEDULER_ENABLE_ENV)


def job_env_name(job_key: str) -> str:
    return JOB_ENV_PREFIX + job_key.upper().replace(".", "_").replace("-", "_")


def requested_new_jobs(
    job_keys: Iterable[str], *, environ: Mapping[str, str] | None = None,
) -> frozenset[str]:
    """Return the explicitly migrated job keys, rejecting unknown route values."""

    env = os.environ if environ is None else environ
    known = frozenset(str(key) for key in job_keys)
    requested = set()
    for key in known:
        value = str(env.get(job_env_name(key), "legacy")).strip().lower()
        if value == "legacy":
            continue
        if value == "new":
            requested.add(key)
            continue
        raise SchedulerMigrationError(
            "%s must be legacy or new" % job_env_name(key))
    for name, value in env.items():
        if not name.startswith(JOB_ENV_PREFIX):
            continue
        if name in {job_env_name(key) for key in known}:
            continue
        if str(value).strip().lower() == "new":
            raise SchedulerMigrationError(
                "%s selects an unknown scheduled job" % name)
    return frozenset(requested)


def uses_new_engine(job_key: str, *, job_keys: Iterable[str],
                    environ: Mapping[str, str] | None = None) -> bool:
    """Whether this process must suppress the legacy owner for ``job_key``."""

    return scheduler_enabled(environ) and job_key in requested_new_jobs(
        job_keys, environ=environ)


__all__ = [
    "JOB_ENV_PREFIX", "SCHEDULER_ENABLE_ENV", "SchedulerMigrationError",
    "job_env_name", "requested_new_jobs", "scheduler_enabled", "uses_new_engine",
]
