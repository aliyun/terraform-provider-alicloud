"""Pure due-time and slot-identity calculations for scheduled jobs."""

from __future__ import annotations

from datetime import datetime, time, timedelta, timezone
from math import floor
from zoneinfo import ZoneInfo

from .model import AdaptiveSchedule, DailySchedule, IntervalSchedule, JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


_SHANGHAI = ZoneInfo("Asia/Shanghai")
_UTC = timezone.utc


def _require_aware(value: datetime, name: str) -> datetime:
    if not is_aware(value):
        raise ValueError(f"{name} must be a timezone-aware datetime")
    return value


class TriggerPlanner:
    """Calculate due times only; retries and persistence remain an Engine concern."""

    def initial_due(self, definition: ScheduledJobDefinition, now: datetime) -> datetime:
        _require_aware(now, "now")
        schedule = definition.schedule
        if isinstance(schedule, IntervalSchedule):
            return now if schedule.run_immediately else now + timedelta(seconds=schedule.interval_seconds)
        if isinstance(schedule, DailySchedule):
            local_now = now.astimezone(_SHANGHAI)
            return datetime.combine(local_now.date(), time(schedule.hour, schedule.minute), tzinfo=_SHANGHAI)
        if isinstance(schedule, AdaptiveSchedule):
            return now if schedule.run_immediately else now + timedelta(seconds=schedule.default_delay_seconds)
        raise TypeError("definition.schedule must be supported")

    def next_due(self, definition: ScheduledJobDefinition, *, slot_due_at: datetime, completed_at: datetime, result: JobResult) -> datetime:
        _require_aware(slot_due_at, "slot_due_at")
        _require_aware(completed_at, "completed_at")
        if not isinstance(result, JobResult) or result.status is not JobResultStatus.SUCCEEDED:
            raise ValueError("only a successful result may create a new slot")
        schedule = definition.schedule
        if isinstance(schedule, IntervalSchedule):
            delta = (completed_at - slot_due_at).total_seconds()
            steps = 1 if delta < 0 else max(1, floor(delta / schedule.interval_seconds) + 1)
            return slot_due_at + timedelta(seconds=steps * schedule.interval_seconds)
        if isinstance(schedule, DailySchedule):
            local_slot = slot_due_at.astimezone(_SHANGHAI)
            return datetime.combine(local_slot.date() + timedelta(days=1), time(schedule.hour, schedule.minute), tzinfo=_SHANGHAI)
        if isinstance(schedule, AdaptiveSchedule):
            due = result.next_due_at or completed_at + timedelta(seconds=schedule.default_delay_seconds)
            lower = completed_at + timedelta(seconds=schedule.min_delay_seconds)
            upper = completed_at + timedelta(seconds=schedule.max_delay_seconds)
            if due < lower or due > upper:
                raise ValueError("adaptive next_due_at is outside the declared bounds")
            return due
        raise TypeError("definition.schedule must be supported")

    def slot_key(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> str:
        """Return the stable identity, including the definition revision."""

        _require_aware(scheduled_for, "scheduled_for")
        utc = scheduled_for.astimezone(_UTC)
        return f"{definition.id}@r{definition.revision}@{utc.strftime('%Y-%m-%dT%H:%M:%S.%fZ')}"
