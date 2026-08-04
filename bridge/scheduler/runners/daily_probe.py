"""Scheduler-owned daily Terraform probe built on the generic HeadlessRuntime."""

from __future__ import annotations

from datetime import datetime, time as wall_time, timedelta
import json
import os
from pathlib import Path
import re
import time
import uuid
from typing import Any, Mapping
from zoneinfo import ZoneInfo

from bridge.headless_runtime import (
    HeadlessRequest,
    HeadlessRuntime,
    Lane,
    SessionPolicy,
)

from ..model import (
    DailySchedule,
    JobResult,
    JobResultStatus,
    ScheduledJobDefinition,
)


JOB_KEY = "daily.probe"
RUNNER_KEY = "daily_probe"
PROTOCOL = "probe-result-v1"
_SHANGHAI = ZoneInfo("Asia/Shanghai")
_RESULT_RE = re.compile(r"\[\[PROBE_RESULT:(.*?)\]\]", re.DOTALL)


def probe_prompt(round_id: str) -> str:
    """Build the business prompt for one pure discovery round."""

    return (
        "【headless 探测轮 %s】你是 Jarvis headless 实例，按 loops/tf-probe.md 跑一轮合成客户探测：\n"
        "1) bootstrap/probe.sh doctor 预检（不绿则记缺口后退出）。\n"
        "2) tier-0：bootstrap/probe.sh tier0 增量扫本轮涉及资源，judgment_queue 走双层查证。\n"
        "3) tier-1：bootstrap/probe.sh list 挑最久未跑的 ≤ config.limits.max_scenarios_per_run 个场景，"
        "逐个 bootstrap/probe.sh run <id>（region 默认 focus）。\n"
        "4) findings 处置严格按 .claude/skills/tf-customer-probe SKILL.md Step C/D 与 "
        "config/probe.json ticket.mode 执行：有效 finding 去重后直接创建 Aone，不写本地中间工单文件"
        "（去重+日上限纪律见 skill）。\n"
        "5) bootstrap/probe.sh sweep 清残留（残留退 1 即停并升级）。\n"
        "6) bootstrap/probe.sh archive 归档已建单 finding / 超期 verdict / 工作目录。\n"
        "7) 按 .claude/skills/tf-customer-probe/references/knowledge-distillation.md 契约把本轮学到的"
        "产品级知识蒸馏进 playground <product>/KNOWLEDGE.md，并在轮次汇报列出。\n"
        "这是纯探测轮，不持有工单、免 bookend；结束把轮次摘要"
        "（tier0 资源数/findings、tier1 场景数/Aone 新建或去重数/env 数、归档件数、"
        "蒸馏条目数）汇报，并在最终输出末尾严格附带单行机器结果：\n"
        '[[PROBE_RESULT:{"outcome":"success","cleanup_ok":true,'
        '"summary":"<非空的本轮摘要>"}]]\n'
        "只有所有步骤完成且 sweep 确认无残留时才允许 outcome=success、cleanup_ok=true；"
        "否则 outcome=failed 或 cleanup_ok=false，并在 summary 写明失败阶段。"
        % round_id
    )


class DailyProbeRunner:
    """Build and execute one Scheduler-admitted, non-persistent Probe session."""

    def __init__(
        self,
        *,
        runtime: HeadlessRuntime,
        summary_root: Path,
        environ: Mapping[str, str] | None = None,
        logger: Any,
        sleeper: Any = time.sleep,
    ) -> None:
        self._runtime = runtime
        self._summary_root = Path(summary_root)
        self._environ = os.environ if environ is None else environ
        self._log = logger
        self._sleep = sleeper

    def run(
        self,
        definition: ScheduledJobDefinition,
        scheduled_for: datetime,
    ) -> JobResult:
        error = _definition_error(definition)
        if error:
            return JobResult(JobResultStatus.PERMANENT_FAILURE, error=error)
        try:
            request = HeadlessRequest(
                prompt=probe_prompt(
                    _round_id(definition.schedule, scheduled_for)),
                session_id=str(uuid.uuid4()),
                session_policy=SessionPolicy.NEW,
                lane=Lane.TERRAFORM,
                timeout_seconds=_positive_float(
                    self._environ.get("JARVIS_DISPATCH_TIMEOUT", "86400"),
                    "JARVIS_DISPATCH_TIMEOUT",
                ),
                max_retries=_nonnegative_int(
                    self._environ.get("JARVIS_DISPATCH_RETRY_MAX", "2"),
                    "JARVIS_DISPATCH_RETRY_MAX",
                ),
                retry_backoff_seconds=_nonnegative_float(
                    self._environ.get("JARVIS_DISPATCH_RETRY_BACKOFF", "30"),
                    "JARVIS_DISPATCH_RETRY_BACKOFF",
                ),
                # The guardian only provides parent-death cleanup here; it does
                # not create a Task/Session fence. Planned stop keeps the
                # Scheduler alive and waits, while an abnormal parent exit
                # closes the sentinel and terminates the orphan process group.
                on_spawn=lambda _process: None,
                guarded=True,
            )
            summary_write_retries = _nonnegative_int(
                self._environ.get("JARVIS_PROBE_SUMMARY_WRITE_RETRIES", "2"),
                "JARVIS_PROBE_SUMMARY_WRITE_RETRIES",
            )
        except (TypeError, ValueError) as exc:
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="invalid daily probe headless contract: %s" % exc,
            )
        result = self._runtime.execute(request)
        if result.is_error:
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error="headless %s after %s attempt(s): %s"
                % (result.subtype, result.attempts, _bounded(result.text)),
            )
        if result.subtype == "no_result":
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error="headless execution returned no terminal result",
            )
        protocol_result, protocol_error = _parse_probe_result(result.text)
        if protocol_error:
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error=protocol_error,
            )
        assert protocol_result is not None
        round_id = _round_id(definition.schedule, scheduled_for)
        summary_error = self._write_summary_with_retry(
            round_id, protocol_result["summary"], summary_write_retries)
        if summary_error is not None:
            self._log.warning(
                "DailyProbeRunner: summary write failed round=%s error=%s",
                round_id,
                type(summary_error).__name__,
            )
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="daily probe summary write failed: %s"
                % type(summary_error).__name__,
            )
        self._log.info(
            "DailyProbeRunner: completed round=%s attempts=%s",
            round_id,
            result.attempts,
        )
        return JobResult(JobResultStatus.SUCCEEDED)

    def _write_summary_with_retry(
        self, round_id: str, text: str, retries: int,
    ) -> OSError | None:
        for attempt in range(retries + 1):
            try:
                self._write_summary(round_id, text)
                return None
            except OSError as exc:
                if attempt >= retries:
                    return exc
                self._sleep(min(0.1 * (attempt + 1), 1.0))
        raise AssertionError("unreachable")

    def _write_summary(self, round_id: str, text: str) -> None:
        self._summary_root.mkdir(parents=True, exist_ok=True)
        target = self._summary_root / ("%s-summary.md" % round_id)
        temporary = target.with_suffix(target.suffix + ".tmp")
        temporary.write_text(text or "", encoding="utf-8")
        os.replace(temporary, target)


def _definition_error(definition: ScheduledJobDefinition) -> str | None:
    if definition.id != JOB_KEY:
        return "daily probe runner received mismatched definition"
    if not isinstance(definition.schedule, DailySchedule):
        return "daily probe requires a daily schedule"
    return None


def _parse_probe_result(text: str) -> tuple[dict[str, Any] | None, str | None]:
    matches = _RESULT_RE.findall(str(text or ""))
    if not matches:
        return None, "probe result sentinel is missing"
    try:
        value = json.loads(matches[-1])
    except (TypeError, ValueError):
        return None, "probe result sentinel is invalid JSON"
    if not isinstance(value, dict):
        return None, "probe result sentinel must be an object"
    outcome = value.get("outcome")
    cleanup_ok = value.get("cleanup_ok")
    summary = value.get("summary")
    if outcome != "success":
        return None, "probe reported unsuccessful outcome"
    if cleanup_ok is not True:
        return None, "probe cleanup was not confirmed"
    if not isinstance(summary, str) or not summary.strip():
        return None, "probe result summary is empty"
    return {
        "outcome": outcome,
        "cleanup_ok": cleanup_ok,
        "summary": summary.strip(),
    }, None


def _round_id(schedule: DailySchedule, scheduled_for: datetime) -> str:
    """Name the latest daily plan boundary at/before a possibly retried instant."""

    local = scheduled_for.astimezone(_SHANGHAI)
    boundary = datetime.combine(
        local.date(),
        wall_time(schedule.hour, schedule.minute),
        tzinfo=_SHANGHAI,
    )
    if local < boundary:
        boundary -= timedelta(days=1)
    return "probe-%s" % boundary.date().isoformat()


def _positive_float(value: object, name: str) -> float:
    try:
        parsed = float(value)
    except (TypeError, ValueError) as exc:
        raise ValueError("%s must be positive" % name) from exc
    if parsed <= 0:
        raise ValueError("%s must be positive" % name)
    return parsed


def _nonnegative_float(value: object, name: str) -> float:
    try:
        parsed = float(value)
    except (TypeError, ValueError) as exc:
        raise ValueError("%s must be non-negative" % name) from exc
    if parsed < 0:
        raise ValueError("%s must be non-negative" % name)
    return parsed


def _nonnegative_int(value: object, name: str) -> int:
    try:
        parsed = int(value)
    except (TypeError, ValueError) as exc:
        raise ValueError("%s must be a non-negative integer" % name) from exc
    if parsed < 0 or str(parsed) != str(value).strip():
        raise ValueError("%s must be a non-negative integer" % name)
    return parsed


def _bounded(value: str, limit: int = 500) -> str:
    text = str(value or "").strip().replace("\n", " | ")
    return text[:limit] or "no result text"


__all__ = [
    "DailyProbeRunner",
    "JOB_KEY",
    "PROTOCOL",
    "probe_prompt",
]
