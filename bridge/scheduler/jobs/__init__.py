"""Load the complete production Job registry from the adjacent ``jobs.yaml``."""

from __future__ import annotations

from typing import Mapping

from ..model import CapabilityValidationContext, ScheduledJobDefinition, definition_digest, validate_registry
from ..registry import SchedulerRegistry, load_scheduler_registry
from .smoke import RUNNER_KEY as SCHEDULER_SMOKE_RUNNER_KEY, SchedulerSmokeRunner


REGISTRY: SchedulerRegistry = load_scheduler_registry()
JOBS: tuple[ScheduledJobDefinition, ...] = REGISTRY.definitions

JOB_CAPABILITIES = CapabilityValidationContext(
    headless_builder_protocols=set(),
    handler_keys={"aone.scan", "aone.stale_claim", "daily.nudge", "aone.reply", "pr.watch", "pr.lifecycle", "scheduler.smoke"},
)
JOB_DIGESTS: Mapping[str, str] = {definition.id: definition_digest(definition) for definition in JOBS}


def load_jobs() -> tuple[ScheduledJobDefinition, ...]:
    return validate_registry(JOBS, context=JOB_CAPABILITIES, expected_digests=JOB_DIGESTS)


__all__ = [
    "REGISTRY", "JOBS", "JOB_CAPABILITIES", "JOB_DIGESTS", "load_jobs",
    "SCHEDULER_SMOKE_RUNNER_KEY", "SchedulerSmokeRunner",
]
