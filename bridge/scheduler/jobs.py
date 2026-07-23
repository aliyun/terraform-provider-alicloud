"""Explicit, import-time validated registry of all Bridge scheduled jobs."""

from __future__ import annotations

from .model import (
    AdaptiveSchedule,
    DailySchedule,
    HandlerRunner,
    IntervalSchedule,
    MisfirePolicy,
    ScheduledJobDefinition,
)


def _job(
    key: str,
    description: str,
    schedule: object,
    runner: object,
    misfire: MisfirePolicy,
    retry: float,
    *,
    revision: int = 1,
    enabled: bool = True,
) -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        key, revision, description, schedule, runner, misfire, retry, enabled)


JOBS = (
    _job(
        "daily.probe",
        "每日启动独立 Terraform Headless 会话执行合成客户探测并归档摘要",
        DailySchedule(10, 0),
        HandlerRunner("daily.probe"),
        MisfirePolicy.CURRENT_DAY,
        300,
        revision=2,
        enabled=False,
    ),
    _job(
        "aone.scan",
        "扫描 Aone 指派、跟踪与待处理任务并写入控制面任务队列",
        IntervalSchedule(1800, True),
        HandlerRunner("aone.scan"),
        MisfirePolicy.COALESCE,
        60,
    ),
    _job(
        "aone.claim-health",
        "核验已认领工单与控制面会话状态，并补偿生命周期事件",
        IntervalSchedule(300, True),
        HandlerRunner("aone.claim-health"),
        MisfirePolicy.COALESCE,
        30,
    ),
    _job(
        "daily.nudge",
        "每日检查 Terraform 待处理任务并对停滞事项发送双通道催办",
        DailySchedule(9, 0),
        HandlerRunner("daily.nudge"),
        MisfirePolicy.CURRENT_DAY,
        300,
    ),
    _job(
        "aone.reply",
        "轮询控制面中等待 Aone 人工回复的会话，并在收到回复后唤醒",
        IntervalSchedule(30, True),
        HandlerRunner("aone.reply"),
        MisfirePolicy.COALESCE,
        10,
    ),
    _job(
        "pr.watch",
        "看守已登记 PR 的评审、CI 与合并生命周期，并按活跃度动态调整检查间隔",
        AdaptiveSchedule(5, 3600, 604800, False),
        HandlerRunner("pr.watch"),
        MisfirePolicy.WAIT_FOR_COMPLETION,
        60,
    ),
    _job(
        "external.recovery",
        "核验控制面中未知的 Aone 外部操作回执，确认后安全收敛",
        IntervalSchedule(300, True),
        HandlerRunner("external.recovery"),
        MisfirePolicy.COALESCE,
        30,
    ),
)


def runner_key(definition: ScheduledJobDefinition) -> object:
    return definition.runner.handler_key


JOB_KEYS = frozenset(definition.id for definition in JOBS)
if len(JOB_KEYS) != len(JOBS):
    raise ValueError("scheduled job ids must be unique")
RUNNER_KEYS = frozenset(runner_key(definition) for definition in JOBS)


__all__ = ["JOBS", "JOB_KEYS", "RUNNER_KEYS", "runner_key"]
