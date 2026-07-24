#!/usr/bin/env python3
"""Route Jarvis work by recovery semantics.

There are only two execution kinds:

* ``TASK`` survives interruption and is always persisted in the control plane.
* ``EPHEMERAL_JOB`` is disposable and may execute locally.

There is no execution-mode switch: a Task is always control-plane owned.
"""

from __future__ import annotations

import logging
import os
import json
import re
import subprocess
import hashlib
from typing import Any, Callable, Iterable, NamedTuple, Optional, Tuple

from bridge.jarvis_task_client import (
    ControlPlaneClient,
    ControlPlaneError,
    ControlPlaneUnavailable,
    TaskEnvelope,
)
from bridge.helpers.aone import REPO_ROOT


TASK = "TASK"
EPHEMERAL_JOB = "EPHEMERAL_JOB"
DEFAULT_TASK_TYPES = (
    "field_repair",
    "ticket",
    "revisit",
    "persona",
    "wake",
    "pr_ci_fix",
    "pr_comment_reply",
)
HEADLESS_POLICY_REVISION = "terraform-rd-single-writer-v4"
LOG = logging.getLogger(__name__)


def broadcast_target() -> str:
    return (os.environ.get("JARVIS_BROADCAST_TARGET")
            or os.environ.get("JARVIS_NOTIFY_GROUP")
            or "cidy1mv+qvMEybkqTXcsXTOeQ==")


def broadcast_type() -> str:
    return os.environ.get("JARVIS_BROADCAST_TYPE", "group")


def _aone_task_key(project: object, item_id: object) -> str:
    value = str(project or "unknown").strip() or "unknown"
    return "aone:%s:%s" % (value, str(item_id))


def _task_envelope(*, item_id, project, task_type, source_type, source_ref,
                   desired_revision, trigger, prompt, recovery_policy="RESUME_ONLY",
                   persona=None, priority=None, comment_cursor=None,
                   required_capabilities=None, max_retries=None,
                   source_status=None, **payload) -> TaskEnvelope:
    """Build a control-plane Task without importing a worker implementation."""
    body = {"itemId": str(item_id), "project": str(project or ""),
            "kind": str(task_type), "prompt": prompt,
            "policyRevision": HEADLESS_POLICY_REVISION}
    body.update(payload)
    return TaskEnvelope(
        task_key=(str(item_id) if str(task_type).lower() == "probe"
                  else _aone_task_key(project, item_id)),
        source_type=source_type, source_ref=source_ref, task_type=task_type,
        desired_revision=desired_revision, trigger_mask=[trigger], payload=body,
        recovery_policy=recovery_policy, persona=persona, priority=priority,
        comment_cursor=comment_cursor, required_capabilities=required_capabilities,
        max_retries=max_retries, source_status=source_status)


def _source_ref_with_title(source_ref, title):
    result = dict(source_ref)
    if str(title or "").strip():
        result["title"] = str(title).strip()
    return result


def _attention_owner_staff_id(value):
    raw = str(value or "").strip()
    if raw and not raw.upper().startswith("WORKER_"):
        if re.fullmatch(r"(?:\d+|WB\d+)", raw, re.IGNORECASE):
            return raw
        try:
            contacts = json.loads((REPO_ROOT / "config/contacts.json").read_text())
            for contact in contacts.get("contacts", []):
                if isinstance(contact, dict) and raw.lower() in {
                        str(contact.get(key) or "").lower()
                        for key in ("id", "name", "flower")}:
                    staff_id = str(contact.get("id") or "").strip()
                    if staff_id and not staff_id.upper().startswith("WORKER_"):
                        return staff_id
        except Exception:  # noqa: BLE001
            pass
    return os.environ.get("JARVIS_MASTER_STAFF", "320687").strip()


def _notify_task_attention(owner_staff_id, payload):
    lines = ["Jarvis 检测到一个需要你关注的工单。", "",
             "原因：%s" % str(payload.get("reason") or "需要人工关注"),
             "建议操作：%s" % str(payload.get("action") or "请打开看板查看并处理")]
    for label, key in (("Aone", "aoneUrl"), ("PR", "prUrl")):
        if payload.get(key):
            lines.append("%s：%s" % (label, payload[key]))
    try:
        subprocess.run(
            [str(REPO_ROOT / "bootstrap/notify-dingtalk.sh"), str(owner_staff_id),
             "Jarvis 工单关注提醒", "--body-stdin"], input="\n".join(lines),
            capture_output=True, text=True, timeout=30, check=False)
    except Exception as exc:  # noqa: BLE001
        LOG.warning("Task attention notification failed: %s", exc)


class _TaskAttentionPublisher:
    """Control-plane attention projection shared by task producers."""

    def __init__(self, client, notifier=None, source="task", required=False):
        self.client, self.notifier = client, notifier or _notify_task_attention
        self.source, self.required = str(source or "task"), bool(required)

    def upsert(self, task_id, owner_staff_id, event_key, payload):
        method = getattr(self.client, "upsert_task_attention", None)
        if not callable(method):
            return not self.required
        owner = _attention_owner_staff_id(owner_staff_id)
        try:
            response = method(str(task_id), owner, str(event_key), dict(payload))
        except Exception as exc:  # noqa: BLE001
            return not self.required if getattr(exc, "status", None) == 404 else False
        if not isinstance(response, dict):
            return False
        if response.get("notify") is True:
            try:
                self.notifier(owner, dict(payload))
            except Exception as exc:  # noqa: BLE001
                LOG.warning("attention[%s] notifier failed: %s", self.source, exc)
        return True

    def clear(self, task_id, event_key_prefix=None):
        method = getattr(self.client, "clear_task_attention", None)
        if not callable(method):
            return not self.required
        try:
            method(str(task_id), event_key_prefix=event_key_prefix)
            return True
        except Exception as exc:  # noqa: BLE001
            return not self.required if getattr(exc, "status", None) == 404 else False


class WakePersistence:
    """Persist the wake generation created by an observed Aone reply."""

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
        aone_id = str(aone_id)
        comments = list(new_comments)
        reply_text = "\n".join(
            "@%s: %s" % (comment.get("creator", "?"), comment.get("content", ""))
            for comment in comments)
        terraform = bool(task.get("terraform"))
        project = str(task.get("project") or self._project_for(aone_id) or "")
        prompt = (
            "工单 #%s 收到新回复:\n%s\n\n请继续处理。\n\n%s"
            % (aone_id, reply_text, self._result_instructions(aone_id, terraform)))
        comment_ids = []
        for comment in comments:
            try:
                comment_ids.append(int(comment.get("id")))
            except (AttributeError, TypeError, ValueError):
                pass
        cursor = max(comment_ids) if comment_ids else None
        revision = (
            "comment:%s" % cursor
            if cursor is not None
            else "comments:%s" % hashlib.sha256(
                reply_text.encode("utf-8")).hexdigest()[:20]
        )
        title = str(task.get("title") or "").strip() or self._title_for(aone_id)
        source_ref: dict[str, Any] = {"aoneId": aone_id, "projectId": project}
        if str(title or "").strip():
            source_ref["title"] = str(title).strip()
        result = self._router.enqueue(TaskEnvelope(
            task_key=_aone_task_key(project, aone_id), source_type="AONE",
            source_ref=source_ref, task_type="wake", desired_revision=revision,
            trigger_mask=["WAKE"], payload={
                "itemId": aone_id, "project": project, "kind": "wake", "prompt": prompt,
                "policyRevision": self._policy_revision,
                "priorRuntimeSessionId": task.get("session_id"), "terraform": terraform,
                "target": task["target"], "targetType": task["target_type"],
            }, recovery_policy="RESUME_ONLY", comment_cursor=cursor,
            source_status=task.get("sourceStatus"),
        ))
        if not result.accepted:
            return False
        if self._routine_notice is not None:
            line = self._line_for(aone_id)
            self._routine_notice(
                "工单收到回复，Task 已进入唤醒队列: %s"
                % (line[0] if isinstance(line, tuple) else line))
        return True


class ExecutionRoute(NamedTuple):
    """Stable domain classification for one work item."""

    execution_kind: str
    needs_recovery: bool


class EnqueueResult(NamedTuple):
    accepted: bool
    reason: str


LocalSubmit = Callable[[], Tuple[bool, str]]


class ExecutionRouter:
    """Classify work and send it to its only valid executor.

    ``task_types`` contains work kinds that need recovery.  ``*`` classifies
    every envelope as a ``TASK``; an empty set classifies all work as
    ``EPHEMERAL_JOB``.
    """

    def __init__(self, *, client: Optional[ControlPlaneClient] = None,
                 task_types: Optional[Iterable[str]] = None,
                 logger: Optional[logging.Logger] = None):
        configured_types = DEFAULT_TASK_TYPES if task_types is None else task_types
        self.client = client
        self.task_types = {
            str(value).strip().lower() for value in configured_types
            if str(value).strip()
        }
        self.log = logger or logging.getLogger(__name__)

    @staticmethod
    def classify(needs_recovery: bool) -> ExecutionRoute:
        """Map the sole domain decision to its stable execution kind."""
        if not isinstance(needs_recovery, bool):
            raise TypeError("needs_recovery must be a bool")
        return ExecutionRoute(TASK if needs_recovery else EPHEMERAL_JOB,
                              needs_recovery)

    def route(self, envelope: TaskEnvelope) -> ExecutionRoute:
        if not isinstance(envelope, TaskEnvelope):
            raise TypeError("envelope must be a TaskEnvelope")
        task_type = envelope.task_type.strip().lower()
        needs_recovery = "*" in self.task_types or task_type in self.task_types
        return self.classify(needs_recovery)

    def is_task(self, envelope: TaskEnvelope) -> bool:
        return self.route(envelope).needs_recovery

    @staticmethod
    def _remote_result(response: object) -> EnqueueResult:
        if isinstance(response, dict) and response.get("accepted") is False:
            return EnqueueResult(
                False, str(response.get("reason") or "control_plane_rejected"))
        reason = "task_persisted"
        if isinstance(response, dict) and response.get("reason"):
            reason = str(response["reason"])
        return EnqueueResult(True, reason)

    def _persist_task(self, envelope: TaskEnvelope) -> EnqueueResult:
        """Persist a recoverable Task, mapping every failure fail-closed."""
        if self.client is None:
            return EnqueueResult(False, "control_plane_unconfigured")
        try:
            request_id = envelope.request_id("upsert")
            response = self.client.upsert_desired_task(
                envelope,
                request_id=request_id,
            )
            return self._remote_result(response)
        except ControlPlaneUnavailable:
            return EnqueueResult(False, "control_plane_unavailable")
        except ControlPlaneError as exc:
            reason = "control_plane_rejected"
            if getattr(exc, "code", ""):
                reason += ":" + str(exc.code)
            return EnqueueResult(False, reason)
        except Exception as exc:
            self.log.exception(
                "task persistence failed closed task=%s error=%s",
                envelope.task_key, type(exc).__name__)
            return EnqueueResult(False, "control_plane_error")

    @staticmethod
    def _run_ephemeral(local_submit: Optional[LocalSubmit]) -> EnqueueResult:
        if local_submit is None:
            return EnqueueResult(False, "local_submit_missing")
        accepted, reason = local_submit()
        return EnqueueResult(bool(accepted), str(reason))

    def observe(self, envelope: TaskEnvelope) -> EnqueueResult:
        """Persist Task desired state; EphemeralJob has no durable observation."""
        route = self.route(envelope)
        if route.needs_recovery:
            return self._persist_task(envelope)
        return EnqueueResult(True, "ephemeral_noop")

    def enqueue(self, envelope: TaskEnvelope,
                local_submit: Optional[LocalSubmit] = None) -> EnqueueResult:
        """Persist Task or execute EphemeralJob locally, never both."""
        route = self.route(envelope)
        if route.needs_recovery:
            return self._persist_task(envelope)
        return self._run_ephemeral(local_submit)
