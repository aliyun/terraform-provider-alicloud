"""The sole explicit registry of production scheduled-job definitions.

Definitions intentionally name handler/protocol keys only.  U3 supplies the runtime
implementations; importing this registry does not import the legacy Bridge loops.
"""

from __future__ import annotations

from typing import Mapping

from ..model import (
    AdaptiveSchedule, CapabilityValidationContext, DailySchedule, HandlerRunner,
    HeadlessRunner, IntervalSchedule, JobPurpose, MisfirePolicy, ReplayPolicy,
    ScheduledJobDefinition, definition_digest, validate_registry,
)


JOBS: tuple[ScheduledJobDefinition, ...] = (
    ScheduledJobDefinition("aone.scan", 1, "Discover Aone work and upsert ticket Tasks", JobPurpose.DISCOVERY, IntervalSchedule(1800, True), HandlerRunner("aone.scan"), MisfirePolicy.COALESCE, 300, 30, ReplayPolicy.TASK_UPSERT_IDEMPOTENT, "JARVIS_SCAN_ENABLE"),
    ScheduledJobDefinition("aone.stale_claim", 1, "Publish stale-claim maintenance events", JobPurpose.MAINTENANCE, IntervalSchedule(7200, False), HandlerRunner("aone.stale_claim"), MisfirePolicy.COALESCE, 300, 60, ReplayPolicy.EVENT_LEDGER_IDEMPOTENT),
    ScheduledJobDefinition("daily.nudge", 1, "Publish daily idle-work reminders through event ledgers", JobPurpose.MAINTENANCE, DailySchedule(9, 0), HandlerRunner("daily.nudge"), MisfirePolicy.CURRENT_DAY, 900, 300, ReplayPolicy.EVENT_LEDGER_IDEMPOTENT, "JARVIS_REVISIT_SCHED"),
    ScheduledJobDefinition("daily.probe", 1, "Run the bounded daily probe as an Ephemeral job", JobPurpose.DISCOVERY, DailySchedule(10, 0), HeadlessRunner("bridge.daily_probe", "bridge.scheduler.v1", "EPHEMERAL", lane="terraform"), MisfirePolicy.CURRENT_DAY, 3600, 300, ReplayPolicy.EPHEMERAL, "JARVIS_PROBE_SCHED"),
    ScheduledJobDefinition("aone.reply", 1, "Discover durable Aone replies and upsert wake Tasks", JobPurpose.DISCOVERY, AdaptiveSchedule(5, 30, 900, True), HandlerRunner("aone.reply"), MisfirePolicy.WAIT_FOR_COMPLETION, 300, 15, ReplayPolicy.TASK_UPSERT_IDEMPOTENT),
    ScheduledJobDefinition("pr.watch", 1, "Discover PR updates and write durable lifecycle intent", JobPurpose.DISCOVERY, AdaptiveSchedule(600, 3600, 3600, False), HandlerRunner("pr.watch"), MisfirePolicy.WAIT_FOR_COMPLETION, 900, 60, ReplayPolicy.TASK_UPSERT_IDEMPOTENT, "JARVIS_PRWATCH_ENABLE"),
    ScheduledJobDefinition("pr.lifecycle", 1, "Consume durable PR lifecycle intent through event ledgers", JobPurpose.MAINTENANCE, AdaptiveSchedule(30, 60, 600, True), HandlerRunner("pr.lifecycle"), MisfirePolicy.WAIT_FOR_COMPLETION, 900, 30, ReplayPolicy.EVENT_LEDGER_IDEMPOTENT),
)

JOB_CAPABILITIES = CapabilityValidationContext(
    headless_builder_protocols={("bridge.daily_probe", "bridge.scheduler.v1")},
    handler_keys={"aone.scan", "aone.stale_claim", "daily.nudge", "aone.reply", "pr.watch", "pr.lifecycle"},
)
JOB_DIGESTS: Mapping[str, str] = {definition.id: definition_digest(definition) for definition in JOBS}


def load_jobs() -> tuple[ScheduledJobDefinition, ...]:
    return validate_registry(JOBS, context=JOB_CAPABILITIES, expected_digests=JOB_DIGESTS)


__all__ = ["JOBS", "JOB_CAPABILITIES", "JOB_DIGESTS", "load_jobs"]
