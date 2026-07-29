"""Detect RESUME_ONLY Tasks stranded after their original owner disappears."""

from __future__ import annotations

from datetime import datetime
import hashlib
import json
import os
from pathlib import Path
import time
from typing import Any, Mapping

from bridge.aone_tasks import master_staff
from bridge.helpers.dingtalk import _dingtalk_event_enqueue, _dingtalk_event_flush
from .claim_health import ClaimHealthRunner
from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


JOB_KEY = "task.owner-health"
RUNNER_KEY = "owner_health"
BLOCKING_REASONS = frozenset({
    "RESUME_OWNER_UNAVAILABLE",
    "RESUME_OWNER_NOT_REGISTERED",
    "RESUME_OWNER_NOT_QUEUE_PULLING",
})
ONE_SHOT_REASONS = frozenset({
    "RESUME_OWNER_NOT_QUEUE_PULLING",
})
MIGRATION_TOKENS = (
    "RESUME_OWNER_NOT_QUEUE_PULLING",
    "RESUME_OWNER_UNAVAILABLE",
    "RESUME_OWNER_NOT_REGISTERED",
    "OWNER_NOT_QUEUE_PULLING",
    "OWNER_UNAVAILABLE",
    "OWNER_NOT_REGISTERED",
)


class OwnerHealthRunner:
    """Page READY diagnostics plus RECOVERY_REQUIRED migration outcomes.

    The second inventory closes the control-plane race where a READY diagnostic
    is migrated to RECOVERY_REQUIRED between two reads and therefore vanishes
    from ``ready-diagnostics`` before the alerting pass can observe it.
    """

    def __init__(
        self,
        *,
        task_client: Any,
        logger: Any,
        repo_root: Path | str,
        state_path: Path | str | None = None,
        clock: Any = time.time,
    ) -> None:
        self.client = task_client
        self.logger = logger
        self.repo_root = Path(repo_root)
        self.page_size = max(1, min(500, int(os.environ.get(
            "JARVIS_OWNER_HEALTH_PAGE_SIZE", "100"))))
        self.max_pages = max(1, int(os.environ.get(
            "JARVIS_OWNER_HEALTH_MAX_PAGES", "100")))
        self.repeat_seconds = max(60, int(os.environ.get(
            "JARVIS_OWNER_HEALTH_REPEAT_SEC", "3600")))
        self.state_path = Path(state_path or os.environ.get(
            "JARVIS_OWNER_HEALTH_STATE",
            str(self.repo_root / ".my-day" / "bridge"
                / "owner-unavailable-health.json"),
        ))
        self.clock = clock

    @staticmethod
    def _task_id(task: Mapping[str, Any]) -> int | None:
        value = task.get("id")
        if value is None:
            value = task.get("taskId")
        try:
            number = int(value)
        except (TypeError, ValueError):
            return None
        return number if number > 0 else None

    @classmethod
    def _page(cls, response: Any) -> tuple[list[dict[str, Any]], int | None]:
        if isinstance(response, list):
            rows = response
            next_cursor = None
        elif isinstance(response, Mapping):
            rows = next((
                response.get(key) for key in ("items", "data", "records", "result")
                if isinstance(response.get(key), list)
            ), [])
            next_cursor = (
                response.get("nextAfterTaskId")
                if response.get("nextAfterTaskId") is not None
                else response.get("nextCursor"))
            try:
                next_cursor = int(next_cursor) if next_cursor is not None else None
            except (TypeError, ValueError):
                next_cursor = None
        else:
            raise ValueError("control-plane page must be a list or object")
        return [dict(row) for row in rows if isinstance(row, Mapping)], next_cursor

    def _paged(self, method_name: str) -> list[dict[str, Any]]:
        rows: list[dict[str, Any]] = []
        cursor = 0
        for _page_number in range(self.max_pages):
            method = getattr(self.client, method_name)
            response = method(after_task_id=cursor, limit=self.page_size)
            page, explicit_next = self._page(response)
            rows.extend(page)
            ids = [
                task_id for task_id in (
                    self._task_id(
                        row.get("task")
                        if isinstance(row.get("task"), Mapping) else row)
                    for row in page)
                if task_id is not None
            ]
            next_cursor = explicit_next
            if next_cursor is None and len(page) >= self.page_size and ids:
                next_cursor = max(ids)
            if (not page or next_cursor is None or next_cursor <= cursor):
                break
            cursor = next_cursor
        else:
            self.logger.warning(
                "owner-health %s pagination capped at %s pages",
                method_name, self.max_pages)
        return rows

    @staticmethod
    def _timestamp(value: Any) -> float | None:
        return ClaimHealthRunner._parse_control_time(value)

    @classmethod
    def _event_evidence(
        cls, timeline: Mapping[str, Any],
    ) -> tuple[str | None, float | None]:
        reason = None
        occurred = None
        for event in timeline.get("events") or []:
            if not isinstance(event, Mapping):
                continue
            material = json.dumps(
                event, ensure_ascii=False, sort_keys=True, default=str).upper()
            token = next(
                (candidate for candidate in MIGRATION_TOKENS
                 if candidate in material),
                None,
            )
            if token is None:
                continue
            reason = (
                token if token.startswith("RESUME_")
                else "RESUME_%s" % token)
            event_time = cls._timestamp(
                event.get("occurredAt")
                or event.get("createdAt")
                or event.get("eventTime"))
            if event_time is not None:
                occurred = (
                    event_time if occurred is None else min(occurred, event_time))
        return reason, occurred

    @classmethod
    def _first_seen_hint(
        cls, task: Mapping[str, Any], entry: Mapping[str, Any],
    ) -> float | None:
        values = []
        for source in (entry, task):
            for key in (
                "ownerUnavailableSince", "reasonSince", "migratedAt",
            ):
                parsed = cls._timestamp(source.get(key))
                if parsed is not None:
                    values.append(parsed)
        return min(values) if values else None

    def _ready_blockers(self) -> dict[int, dict[str, Any]]:
        blockers: dict[int, dict[str, Any]] = {}
        for entry in self._paged("list_ready_task_diagnostics"):
            task = entry.get("task")
            if not isinstance(task, Mapping):
                continue
            task_id = self._task_id(task)
            reason = str(entry.get("reasonCode") or "").strip().upper()
            if task_id is None or reason not in BLOCKING_REASONS:
                continue
            blockers[task_id] = {
                "task_id": task_id,
                "aone_id": str(task.get("aoneId") or "").strip(),
                "status": str(task.get("status") or "READY").upper(),
                "reason": reason,
                "required_worker": str(
                    entry.get("requiredWorkerKey") or "-"),
                "required_worker_status": str(
                    entry.get("requiredWorkerActivityStatus") or "-"),
                "session_id": str(task.get("currentSessionId") or ""),
                "first_seen_hint": self._first_seen_hint(task, entry),
                "source": "ready-diagnostics",
            }
        return blockers

    def _recovery_blockers(self) -> dict[int, dict[str, Any]]:
        blockers: dict[int, dict[str, Any]] = {}
        for entry in self._paged("list_source_status_candidates"):
            # list_source_status_candidates only returns Aone source status (4 fields:
            # taskId/sourceProjectKey/aoneId/sourceStatus), not the Task control-plane
            # status/recoveryPolicy/currentSessionId we need for RESUME_ONLY stranded
            # detection. Point-read the full Task to get them; best-effort skip on
            # failure so one bad aone id does not poison the whole pass.
            aone_id = str(entry.get("aoneId") or "").strip()
            if not aone_id:
                continue
            try:
                response = self.client.get_task_by_aone(aone_id)
            except Exception as exc:  # noqa: BLE001
                self.logger.warning(
                    "owner-health task point-read aone=%s failed: %s",
                    aone_id, type(exc).__name__)
                continue
            rows = ClaimHealthRunner._task_rows(response)
            if not rows:
                continue
            task = ClaimHealthRunner._current_task(rows)
            if str(task.get("status") or "").strip().upper() != "RECOVERY_REQUIRED":
                continue
            task_id = self._task_id(task)
            if task_id is None:
                continue
            try:
                timeline = self.client.get_task_timeline(str(task_id))
            except Exception as exc:  # noqa: BLE001
                self.logger.warning(
                    "owner-health timeline task=%s failed: %s",
                    task_id, type(exc).__name__)
                continue
            if not isinstance(timeline, Mapping):
                continue
            reason, migrated_at = self._event_evidence(timeline)
            sessions = [
                session for session in (timeline.get("sessions") or [])
                if isinstance(session, Mapping)
            ]
            current_id = task.get("currentSessionId")
            session = next((
                candidate for candidate in sessions
                if str(candidate.get("id")) == str(current_id)
            ), None)
            resumable = (
                isinstance(session, Mapping)
                and str(session.get("status") or "").upper() == "RESUMABLE")
            recovery_policy = str(task.get("recoveryPolicy") or "").upper()
            if reason is None and not (
                    resumable and recovery_policy == "RESUME_ONLY"):
                continue
            reason = reason or "RESUME_OWNER_UNAVAILABLE"
            current_worker = timeline.get("currentWorker")
            if not isinstance(current_worker, Mapping):
                current_worker = {}
            worker = str(
                (session or {}).get("currentWorkerKey")
                or (session or {}).get("workerKey")
                or current_worker.get("workerKey")
                or "-")
            blocker = {
                "task_id": task_id,
                "aone_id": str(task.get("aoneId") or "").strip(),
                "status": "RECOVERY_REQUIRED",
                "reason": reason,
                "required_worker": worker,
                "required_worker_status": str(
                    current_worker.get("activityStatus")
                    or current_worker.get("status")
                    or "unavailable"),
                "session_id": str(current_id or ""),
                "first_seen_hint": (
                    migrated_at or self._first_seen_hint(task, timeline)),
                "source": "recovery-migration",
            }
            blockers[task_id] = blocker
        return blockers

    def collect_blockers(self) -> dict[int, dict[str, Any]]:
        # READY→RECOVERY_REQUIRED is covered by the middle inventory;
        # RECOVERY_REQUIRED→READY is covered by the final diagnostic pass.
        blockers = self._ready_blockers()
        blockers.update(self._recovery_blockers())
        blockers.update(self._ready_blockers())
        return blockers

    def _load_state(self) -> dict[str, Any]:
        try:
            value = json.loads(self.state_path.read_text(encoding="utf-8"))
        except (OSError, ValueError):
            return {"version": 1, "episodes": {}}
        episodes = value.get("episodes") if isinstance(value, Mapping) else None
        return {
            "version": 1,
            "episodes": dict(episodes) if isinstance(episodes, Mapping) else {},
        }

    def _write_state(self, value: Mapping[str, Any]) -> None:
        self.state_path.parent.mkdir(parents=True, exist_ok=True)
        tmp = self.state_path.with_name(
            "%s.tmp.%s" % (self.state_path.name, os.getpid()))
        tmp.write_text(
            json.dumps(value, ensure_ascii=False, indent=2, sort_keys=True) + "\n",
            encoding="utf-8",
        )
        tmp.replace(self.state_path)

    @staticmethod
    def _duration(seconds: float) -> str:
        total = max(0, int(seconds))
        days, remainder = divmod(total, 86400)
        hours, remainder = divmod(remainder, 3600)
        minutes = remainder // 60
        if days:
            return "%d天%d小时" % (days, hours)
        if hours:
            return "%d小时%d分钟" % (hours, minutes)
        return "%d分钟" % minutes

    def reconcile(
        self, blockers: Mapping[int, Mapping[str, Any]],
    ) -> int:
        now = float(self.clock())
        old = self._load_state().get("episodes", {})
        current: dict[str, Any] = {}
        queued = 0
        for task_id, raw in sorted(blockers.items()):
            item = dict(raw)
            # READY -> RECOVERY_REQUIRED is one continuous unavailable-owner
            # episode.  Status is evidence, not identity: keeping it out of the
            # signature preserves the original dwell clock across migration.
            signature = "%s:%s" % (
                item.get("reason") or "unknown",
                item.get("session_id") or "none",
            )
            previous = old.get(str(task_id))
            if not isinstance(previous, Mapping) \
                    or previous.get("signature") != signature:
                hinted = item.get("first_seen_hint")
                try:
                    first_seen = min(now, float(hinted)) if hinted else now
                except (TypeError, ValueError):
                    first_seen = now
                previous = {
                    "signature": signature,
                    "firstSeenAt": first_seen,
                    "lastAlertAt": 0,
                }
            episode = dict(previous)
            episode.update({
                "signature": signature,
                "lastSeenAt": now,
                "status": item.get("status"),
                "reason": item.get("reason"),
                "aoneId": item.get("aone_id"),
                "source": item.get("source"),
            })
            first_seen = float(episode["firstSeenAt"])
            last_alert = float(episode.get("lastAlertAt") or 0)
            reason = str(item.get("reason") or "")
            repeat_due = (
                reason not in ONE_SHOT_REASONS
                and now - last_alert >= self.repeat_seconds)
            if not last_alert or repeat_due:
                dwell = self._duration(now - first_seen)
                aone_id = str(item.get("aone_id") or "").strip()
                anchor = aone_id if aone_id.isdigit() else str(task_id)
                digest = hashlib.sha256(
                    signature.encode("utf-8")).hexdigest()[:12]
                alert_index = int(max(0, now - first_seen) // self.repeat_seconds)
                event_key = "owner-health:%s:%s:%s" % (
                    task_id, digest, alert_index)
                aone_line = ("Aone：#%s" % aone_id) if aone_id else "Aone：-"
                if reason == "RESUME_OWNER_NOT_QUEUE_PULLING":
                    next_step = (
                        "请复核 Task/Session；确认旧 owner 已不再执行后，人工运行 "
                        "`bootstrap/control-plane-status.sh force-release %s %s "
                        "--reason TEXT --yes`。该告警不会自动释放 ownership。"
                        % (task_id, item.get("session_id") or "[session_id]"))
                else:
                    next_step = (
                        "请复核 Task/Session；RECOVERY_REQUIRED 且旧 RESUMABLE "
                        "上下文不可恢复时，按受保护流程执行 discard-resume。")
                body = (
                    "Jarvis 控制面检测到 RESUME_ONLY Task 的原 owner 不可用。\n\n"
                    "- Task：#%s\n"
                    "- %s\n"
                    "- 当前状态：%s\n"
                    "- 原因：%s\n"
                    "- 原 owner：%s（%s）\n"
                    "- 已滞留：%s\n"
                    "- 证据源：%s\n\n"
                    "%s"
                    % (
                        task_id,
                        aone_line,
                        item.get("status") or "-",
                        item.get("reason") or "-",
                        item.get("required_worker") or "-",
                        item.get("required_worker_status") or "-",
                        dwell,
                        item.get("source") or "-",
                        next_step,
                    ))
                if _dingtalk_event_enqueue(
                        anchor, "jarvis-fleet", event_key, master_staff(),
                        "Jarvis owner-unavailable 告警", body,
                        allow_non_tf=True):
                    queued += 1
                    episode["lastAlertAt"] = now
            current[str(task_id)] = episode
        self._write_state({"version": 1, "episodes": current})
        return queued

    def run(
        self, definition: ScheduledJobDefinition, scheduled_for: datetime,
    ) -> JobResult:
        if definition.id != JOB_KEY or not is_aware(scheduled_for):
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="task.owner-health runner received an invalid slot",
            )
        try:
            queued = self.reconcile(self.collect_blockers())
            _dingtalk_event_flush()
            self.logger.info("task.owner-health queued=%s", queued)
            return JobResult(JobResultStatus.SUCCEEDED)
        except Exception as exc:  # noqa: BLE001
            self.logger.exception(
                "task.owner-health failed closed: %s", type(exc).__name__)
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error=type(exc).__name__,
            )


def build(*, logger: Any, task_client: Any, repo_root: Path | str):
    return OwnerHealthRunner(
        logger=logger, task_client=task_client, repo_root=repo_root)


__all__ = [
    "BLOCKING_REASONS",
    "JOB_KEY",
    "ONE_SHOT_REASONS",
    "OwnerHealthRunner",
    "RUNNER_KEY",
    "build",
]
