"""Independent business enablement for Scheduler job definitions.

Ownership routing lives exclusively in ``config/scheduler-jobs.yaml``.  This
module intentionally contains no route-selection configuration.
"""

from __future__ import annotations

import os
from typing import Mapping


class SchedulerMigrationError(ValueError):
    """The operator supplied an ambiguous or unsupported cutover configuration."""


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
__all__ = [
    "SchedulerMigrationError", "business_job_enabled",
]
