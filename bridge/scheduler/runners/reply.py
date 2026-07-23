"""Scheduler-owned sensor for Aone replies that resume suspended Tasks."""

from __future__ import annotations

import importlib
import json
import os
from datetime import datetime
from pathlib import Path
import subprocess
import time
from typing import Any

try:
    from bridge.execution.wake import WakePersistence
except ModuleNotFoundError:  # pragma: no cover - bridge/main.py top-level import.
    from execution.wake import WakePersistence

from ..model import JobResult, JobResultStatus, ScheduledJobDefinition, is_aware


RUNNER_KEY = "aone.reply"
WAIT_TIERS = ((30 * 60, 120), (2 * 3600, 600), (float("inf"), 1800))


def _load_legacy_module() -> Any:
    try:
        return importlib.import_module("bridge.jarvis_dingtalk_bot")
    except ModuleNotFoundError:  # pragma: no cover - bridge/main.py path.
        return importlib.import_module("jarvis_dingtalk_bot")


class ReplyRunner:
    """Perform one control-plane wait scan without constructing ``JarvisHandler``."""

    def __init__(self, *, task_client: Any, logger: Any,
                 legacy_module: Any | None = None) -> None:
        self._task_client = task_client
        self._logger = logger
        self._module = legacy_module
        self._page_size = max(1, min(500, int(os.environ.get(
            "JARVIS_MANAGED_WAIT_PAGE_SIZE", "100"))))
        self._poll_state: dict[str, dict[str, float]] = {}
        self._wake: WakePersistence | None = None

    @property
    def _legacy(self) -> Any:
        if self._module is None:
            self._module = _load_legacy_module()
        return self._module

    @property
    def _wake_persistence(self) -> WakePersistence:
        if self._wake is None:
            module = self._legacy
            router = module.ExecutionRouter(client=self._task_client, logger=self._logger)
            self._wake = WakePersistence(
                execution_router=router,
                result_instructions=module._task_result_instructions,
                policy_revision=module.HEADLESS_POLICY_REVISION,
                title_for=module.JarvisHandler._workitem_title,
                project_for=module.JarvisHandler._workitem_project,
            )
        return self._wake

    @staticmethod
    def _poll_interval(first_seen: float) -> int:
        age = time.time() - first_seen
        for threshold, interval in WAIT_TIERS:
            if age < threshold:
                return interval
        return WAIT_TIERS[-1][1]

    def _fetch_comments(self, aone_id: str) -> list[dict[str, Any]] | None:
        try:
            result = subprocess.run(
                [str(Path(self._legacy.REPO_ROOT) / "bin" / "a1id"), "--",
                 "project", "workitem", "comment", "list", str(aone_id), "-f", "json"],
                capture_output=True, text=True, timeout=30,
                cwd=str(self._legacy.REPO_ROOT))
            return json.loads(result.stdout) if result.returncode == 0 else None
        except Exception:  # noqa: BLE001 - a failed poll remains retryable.
            return None

    def _list_waits(self):
        after = 0
        while True:
            page = self._task_client.list_pending_aone_reply_waits(
                after_session_id=after, limit=self._page_size)
            if not isinstance(page, dict) or not isinstance(page.get("items"), list):
                raise ValueError("control plane pending wait page is invalid")
            yield from (item for item in page["items"] if isinstance(item, dict))
            if not page.get("hasMore"):
                return
            try:
                next_after = int(page.get("nextAfterSessionId"))
            except (TypeError, ValueError) as exc:
                raise ValueError("control plane pending wait cursor is invalid") from exc
            if next_after <= after:
                raise ValueError("control plane pending wait cursor did not advance")
            after = next_after

    @staticmethod
    def _comment_id(comment: object) -> int | None:
        try:
            return int(comment.get("id"))  # type: ignore[union-attr]
        except (AttributeError, TypeError, ValueError):
            return None

    def _is_human_comment(self, creator: str, content: str) -> bool:
        return self._legacy.AoneScheduler._is_human_comment(creator, content)

    def _tick(self) -> None:
        now = time.time()
        seen: set[str] = set()
        try:
            waits = list(self._list_waits())
        except Exception as exc:  # noqa: BLE001 - next slot rebuilds the snapshot.
            self._logger.warning("aone.reply list failed; will retry: %s", exc)
            return
        for item in waits:
            task = item.get("task") if isinstance(item.get("task"), dict) else {}
            session = item.get("session") if isinstance(item.get("session"), dict) else {}
            session_id = str(session.get("id") or "").strip()
            aone_id = str(session.get("waitKey") or task.get("aoneId") or "").strip()
            if not session_id or not aone_id.isdigit():
                self._logger.warning("aone.reply ignored invalid wait task=%s session=%s",
                                     task.get("taskKey"), session_id or "<empty>")
                continue
            seen.add(session_id)
            state = self._poll_state.setdefault(
                session_id, {"first_seen": now, "last_poll": 0})
            if now - state["last_poll"] < self._poll_interval(state["first_seen"]):
                continue
            state["last_poll"] = now
            comments = self._fetch_comments(aone_id)
            if comments is None:
                continue
            try:
                baseline = int(session.get("waitCursor") or 0)
            except (TypeError, ValueError):
                baseline = 0
            new_comments = [
                comment for comment in comments
                if self._comment_id(comment) is not None
                and self._comment_id(comment) > baseline
                and self._is_human_comment(
                    str(comment.get("creator") or "").strip(),
                    str(comment.get("content") or "").strip())
            ]
            if not new_comments:
                continue
            new_comments.sort(key=lambda comment: self._comment_id(comment) or 0)
            frozen = session.get("inputPayload")
            if not isinstance(frozen, dict):
                frozen = task.get("payload") if isinstance(task.get("payload"), dict) else {}
            wake_context = {
                "session_id": session.get("runtimeSessionId"),
                "terraform": bool(frozen.get("terraform")),
                "project": str(frozen.get("project") or
                               (task.get("sourceRef") or {}).get("projectId") or ""),
                "target": str(frozen.get("target") or self._legacy.broadcast_target()),
                "target_type": str(frozen.get("targetType") or self._legacy.broadcast_type()),
                "title": str(frozen.get("title") or
                             (task.get("sourceRef") or {}).get("title") or ""),
            }
            self._logger.info("aone.reply: #%s session=%s got %d reply comment(s)",
                              aone_id, session_id, len(new_comments))
            if not self._wake_persistence.enqueue(aone_id, wake_context, new_comments):
                self._logger.warning("aone.reply: wake #%s not durably accepted; will retry",
                                     aone_id)
        for session_id in set(self._poll_state) - seen:
            self._poll_state.pop(session_id, None)

    def run(self, definition: ScheduledJobDefinition,
            scheduled_for: datetime) -> JobResult:
        if definition.id != RUNNER_KEY:
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="aone.reply runner received mismatched definition")
        if not is_aware(scheduled_for):
            return JobResult(JobResultStatus.PERMANENT_FAILURE,
                             error="aone.reply runner requires an aware scheduled time")
        self._tick()
        return JobResult(JobResultStatus.SUCCEEDED)


__all__ = ["ReplyRunner", "RUNNER_KEY", "WAIT_TIERS"]
