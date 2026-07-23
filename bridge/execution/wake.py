"""Durable persistence for an Aone reply that resumes a suspended Task."""

from __future__ import annotations

import hashlib
from typing import Any, Callable, Iterable

try:  # Package tests import ``bridge``; bridge/run.sh uses top-level modules.
    from bridge.jarvis_task_client import TaskEnvelope
except ModuleNotFoundError:  # pragma: no cover - exercised by bridge/run.sh.
    from jarvis_task_client import TaskEnvelope


def _aone_task_key(project: object, item_id: object) -> str:
    value = str(project or "unknown").strip() or "unknown"
    return "aone:%s:%s" % (value, str(item_id))


class WakePersistence:
    """Build and durably upsert the ``wake`` Task for an observed Aone reply.

    The caller supplies the narrow environment-specific lookups (Aone title/project,
    prompt policy, and optional routine notice).  No local executor is needed: a wake
    is always a control-plane Task and is later leased by ``PersistenceExecutor``.
    """

    def __init__(
        self,
        *,
        execution_router: Any,
        result_instructions: Callable[[str, bool], str],
        policy_revision: str,
        title_for: Callable[[str], str] | None = None,
        project_for: Callable[[str], str] | None = None,
        line_for: Callable[[str], object] | None = None,
        routine_notice: Callable[[str], None] | None = None,
    ) -> None:
        self._router = execution_router
        self._result_instructions = result_instructions
        self._policy_revision = policy_revision
        self._title_for = title_for or (lambda _aone_id: "")
        self._project_for = project_for or (lambda _aone_id: "")
        self._line_for = line_for or (lambda aone_id: "#%s" % aone_id)
        self._routine_notice = routine_notice

    def enqueue(self, aone_id: object, task: dict[str, Any],
                new_comments: Iterable[dict[str, Any]]) -> bool:
        """Persist a resumable wake Task and return whether it was accepted."""
        aone_id = str(aone_id)
        comments = list(new_comments)
        reply_text = "\n".join(
            "@%s: %s" % (comment.get("creator", "?"), comment.get("content", ""))
            for comment in comments)
        terraform = bool(task.get("terraform"))
        prompt = (
            "工单 #%s 收到新回复:\n%s\n\n请继续处理。\n\n%s"
            % (aone_id, reply_text, self._result_instructions(aone_id, terraform)))
        project = str(task.get("project") or self._project_for(aone_id) or "")
        comment_ids: list[int] = []
        for comment in comments:
            try:
                comment_ids.append(int(comment.get("id")))
            except (AttributeError, TypeError, ValueError):
                pass
        if comment_ids:
            cursor = max(comment_ids)
            revision = "comment:%s" % cursor
        else:
            cursor = None
            revision = "comments:%s" % hashlib.sha256(
                reply_text.encode("utf-8")).hexdigest()[:20]
        title = str(task.get("title") or "").strip() or self._title_for(aone_id)
        source_ref: dict[str, Any] = {"aoneId": aone_id, "projectId": project}
        if str(title or "").strip():
            source_ref["title"] = str(title).strip()
        envelope = TaskEnvelope(
            task_key=_aone_task_key(project, aone_id),
            source_type="AONE",
            source_ref=source_ref,
            task_type="wake",
            desired_revision=revision,
            trigger_mask=["WAKE"],
            payload={
                "itemId": aone_id,
                "project": project,
                "kind": "wake",
                "prompt": prompt,
                "policyRevision": self._policy_revision,
                "priorRuntimeSessionId": task.get("session_id"),
                "terraform": terraform,
                "target": task["target"],
                "targetType": task["target_type"],
            },
            recovery_policy="RESUME_ONLY",
            comment_cursor=cursor,
            source_status=task.get("sourceStatus"),
        )
        result = self._router.enqueue(envelope)
        if not result.accepted:
            return False
        if self._routine_notice is not None:
            line = self._line_for(aone_id)
            if isinstance(line, tuple):
                line = line[0]
            self._routine_notice("工单收到回复，Task 已进入唤醒队列: %s" % line)
        return True


__all__ = ["WakePersistence"]
