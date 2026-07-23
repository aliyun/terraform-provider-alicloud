#!/usr/bin/env python3
"""Route Jarvis work by recovery semantics.

There are only two execution kinds:

* ``TASK`` survives interruption and is always persisted in the control plane.
* ``EPHEMERAL_JOB`` is disposable and may execute locally.

There is no execution-mode switch: a Task is always control-plane owned.
"""

from __future__ import annotations

import logging
from typing import Callable, Iterable, NamedTuple, Optional, Tuple

from bridge.jarvis_task_client import (
    ControlPlaneClient,
    ControlPlaneError,
    ControlPlaneUnavailable,
    TaskEnvelope,
)


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
