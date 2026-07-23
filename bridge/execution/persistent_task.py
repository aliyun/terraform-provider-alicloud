"""Worker-side translation of a leased persistent Task into one execution run.

This module deliberately has no DingTalk, Scheduler, or Bot dependency.  The
composition root supplies the narrow adapters needed for Aone bookends, notices,
and the headless runner.
"""

from __future__ import annotations

import json
import logging
import os
import signal
import subprocess
from typing import Any, Callable, Collection


class PersistentTaskExecution:
    """Execute frozen Task-session inputs through injected worker adapters."""

    def __init__(
        self,
        *,
        enabled_kinds: Callable[[], Collection[str]],
        dispatch_item: Callable[..., Any],
        task_bookend: Callable[..., Any],
        terraform_rd_ready: Callable[[], bool],
        routine_notice: Callable[[str], None],
        quick_card: Callable[[str, str, str], None],
        field_repair_worker: Any = None,
        field_repair_kind: str,
        task_bookend_kinds: Collection[str],
        post_pr_headless_kinds: Collection[str],
        broadcast_target: Callable[[], str],
        broadcast_type: Callable[[], str],
    ) -> None:
        self._enabled_kinds = enabled_kinds
        self._dispatch_item = dispatch_item
        self._task_bookend = task_bookend
        self._terraform_rd_ready = terraform_rd_ready
        self._routine_notice = routine_notice
        self._quick_card = quick_card
        self._field_repair_worker = field_repair_worker
        self._field_repair_kind = field_repair_kind
        self._task_bookend_kinds = frozenset(task_bookend_kinds)
        self._post_pr_headless_kinds = frozenset(post_pr_headless_kinds)
        self._broadcast_target = broadcast_target
        self._broadcast_type = broadcast_type

    def execute(self, lease: object, controller: object) -> Any:
        """Translate one immutable Task lease into the shared execution runtime."""
        task = lease.get("task") if isinstance(lease, dict) else None
        if not isinstance(task, dict):
            raise ValueError("Task lease.task must be an object")
        session = lease.get("session") if isinstance(lease, dict) else None
        if not isinstance(session, dict):
            raise ValueError("Task lease.session must be an object")
        if "inputPayload" not in session:
            raise ValueError("Task session input snapshot is missing")
        payload = session.get("inputPayload")
        if payload is None:
            raise ValueError("Task session input snapshot is null")
        if isinstance(payload, str):
            try:
                payload = json.loads(payload)
            except (TypeError, ValueError) as exc:
                raise ValueError("Task payload must be JSON object") from exc
        if not isinstance(payload, dict):
            raise ValueError("Task payload must be an object")
        kind = str(payload.get("kind") or task.get("taskType") or "").strip().lower()
        enabled = self._enabled_kinds()
        if "*" not in enabled and kind not in enabled:
            raise ValueError("TASK kind is not enabled: %s" % (kind or "<empty>"))
        item_id = str(payload.get("itemId") or task.get("aoneId") or
                      task.get("taskKey") or "").strip()
        prompt = str(payload.get("prompt") or "")
        if not item_id or not prompt:
            raise ValueError("Task requires itemId and prompt")
        target = str(payload.get("target") or self._broadcast_target())
        target_type = str(payload.get("targetType") or self._broadcast_type())
        project = str(payload.get("project") or "")
        terraform = bool(payload.get("terraform"))
        expected_comment_cursor = payload.get("expectedCommentCursor")
        if kind == self._field_repair_kind:
            if self._field_repair_worker is None:
                raise RuntimeError("field repair worker is unavailable")
            return self._field_repair_worker.execute(payload, controller)
        if kind in self._task_bookend_kinds and self._field_repair_worker is not None:
            repair_result = self._field_repair_worker.repair_only(
                item_id, project, terraform=terraform, controller=controller)
            if repair_result.get("status") != "completed":
                return repair_result
        task_bookend = None
        if kind in self._post_pr_headless_kinds:
            if not self._terraform_rd_ready():
                raise RuntimeError(
                    "terraform-rd identity not ready; refusing to run post-PR Task #%s "
                    "closed-fail (no silent SUCCEEDED)" % item_id)
            task_bookend = self._task_bookend(
                controller, item_id, project, True, kind, writes_reply=False)
        elif kind in self._task_bookend_kinds:
            if terraform and not self._terraform_rd_ready():
                raise RuntimeError(
                    "terraform-rd identity not ready; refusing to run Task #%s "
                    "closed-fail (no silent SUCCEEDED)" % item_id)
            task_bookend = self._task_bookend(
                controller, item_id, project, terraform, kind,
                expected_comment_cursor=expected_comment_cursor)
        if task_bookend is not None:
            task_bookend.capture_comment_baseline()
        on_spawn = (task_bookend.bind_process if task_bookend is not None
                    else controller.bind_process)
        notify = ((lambda text: self._quick_card(target, text, target_type))
                  if kind == "adhoc" else self._routine_notice)
        return self._dispatch_item(
            item_id, prompt, controller.runtime_session_id, controller.resumed,
            notify, target, target_type, on_spawn=on_spawn, project=project,
            kind=kind, terraform=terraform, session_controller=controller,
            task_bookend=task_bookend)


def stop_task_process(controller: object, reason: str, *, logger: logging.Logger) -> bool:
    """Fence a leased process without relying on a Bot instance."""
    proc = getattr(controller, "process", None)
    if proc is None:
        return True
    try:
        if proc.poll() is not None:
            return True
    except Exception:  # noqa: BLE001
        pass
    try:
        grace = float(os.environ.get("JARVIS_TASK_STOP_GRACE_SEC", "5"))
    except (TypeError, ValueError):
        grace = 5.0
    grace = max(0.0, min(grace, 60.0))
    try:
        os.killpg(os.getpgid(proc.pid), signal.SIGTERM)
        logger.warning("Task session %s process stopping (%s)",
                       getattr(controller, "session_id", "?"), reason)
    except (ProcessLookupError, OSError, AttributeError):
        try:
            proc.terminate()
        except Exception:  # noqa: BLE001
            pass
    try:
        proc.wait(timeout=grace)
        return True
    except subprocess.TimeoutExpired:
        pass
    except Exception:  # noqa: BLE001
        pass
    try:
        os.killpg(os.getpgid(proc.pid), signal.SIGKILL)
    except (ProcessLookupError, OSError, AttributeError):
        try:
            proc.kill()
        except Exception:  # noqa: BLE001
            pass
    try:
        proc.wait(timeout=5)
    except Exception:  # noqa: BLE001
        logger.exception("Task session %s process could not be reaped (%s)",
                         getattr(controller, "session_id", "?"), reason)
        return False
    logger.warning("Task session %s process force-killed after %.1fs (%s)",
                   getattr(controller, "session_id", "?"), grace, reason)
    return True


__all__ = ["PersistentTaskExecution", "stop_task_process"]
