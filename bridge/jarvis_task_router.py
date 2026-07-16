#!/usr/bin/env python3
"""Execution-mode router for persistent Jarvis task enqueue.

The router is the migration seam between the current in-process DispatchPool
and AutomationAgent.  It intentionally accepts a ``legacy_submit`` callback so
the existing schedulers can be wired one at a time without importing the large
bridge module here.
"""

from __future__ import annotations

import logging
import os
from typing import Callable, Iterable, NamedTuple, Optional, Tuple

from jarvis_task_client import (
    AutomationAgentTaskClient,
    ControlPlaneError,
    ControlPlaneUnavailable,
    TaskEnvelope,
)


class EnqueueResult(NamedTuple):
    accepted: bool
    reason: str


LegacySubmit = Callable[[], Tuple[bool, str]]


class TaskRouter:
    """Route one envelope through legacy, shadow, or managed ownership.

    Modes:
      * ``legacy``: local submit only.
      * ``shadow``: best-effort SHADOW upsert, then local submit regardless of
        control-plane health.  SHADOW tasks must never be leaseable server-side.
      * ``managed``: explicitly allow-listed task types are owned only by the
        data plane.  A failed upsert is fail-closed and never falls back locally.
        Non-allow-listed task types shadow-upsert and continue through legacy;
        an empty allow-list therefore manages nothing.  Use ``*`` to manage all.
    """

    VALID_MODES = {"legacy", "shadow", "managed"}

    def __init__(self, mode: str = "legacy", *,
                 client: Optional[AutomationAgentTaskClient] = None,
                 managed_task_types: Optional[Iterable[str]] = None,
                 logger: Optional[logging.Logger] = None):
        normalized = str(mode or "legacy").strip().lower()
        if normalized not in self.VALID_MODES:
            raise ValueError("invalid task mode %r (expected legacy/shadow/managed)" % mode)
        self.mode = normalized
        self.client = client
        self.managed_task_types = {
            str(value).strip().lower() for value in (managed_task_types or [])
            if str(value).strip()
        }
        self.log = logger or logging.getLogger(__name__)

    @classmethod
    def from_env(cls, *, client: Optional[AutomationAgentTaskClient] = None,
                 logger: Optional[logging.Logger] = None) -> "TaskRouter":
        raw_types = os.environ.get("JARVIS_MANAGED_TASK_TYPES", "")
        return cls(
            os.environ.get("JARVIS_TASK_MODE", "legacy"),
            client=client,
            managed_task_types=raw_types.split(","),
            logger=logger,
        )

    def _is_managed(self, envelope: TaskEnvelope) -> bool:
        if self.mode != "managed":
            return False
        task_type = envelope.task_type.strip().lower()
        return "*" in self.managed_task_types or task_type in self.managed_task_types

    def is_managed(self, envelope: TaskEnvelope) -> bool:
        """Public ownership predicate for scheduler lifecycle decisions."""
        if not isinstance(envelope, TaskEnvelope):
            raise TypeError("envelope must be a TaskEnvelope")
        return self._is_managed(envelope)

    @staticmethod
    def _remote_result(response: object, default_reason: str) -> EnqueueResult:
        if isinstance(response, dict) and response.get("accepted") is False:
            return EnqueueResult(False, str(response.get("reason") or "control_plane_rejected"))
        reason = default_reason
        if isinstance(response, dict) and response.get("reason"):
            reason = str(response["reason"])
        return EnqueueResult(True, reason)

    def _upsert(self, envelope: TaskEnvelope, execution_mode: str) -> EnqueueResult:
        if self.client is None:
            return EnqueueResult(False, "control_plane_unconfigured")
        request_id = envelope.request_id("upsert", execution_mode=execution_mode)
        response = self.client.upsert_desired_task(
            envelope, execution_mode=execution_mode, request_id=request_id)
        return self._remote_result(response, "managed" if execution_mode == "MANAGED" else "shadowed")

    @staticmethod
    def _legacy(legacy_submit: Optional[LegacySubmit]) -> EnqueueResult:
        if legacy_submit is None:
            return EnqueueResult(False, "legacy_submit_missing")
        accepted, reason = legacy_submit()
        return EnqueueResult(bool(accepted), str(reason))

    def _observe_shadow(self, envelope: TaskEnvelope) -> EnqueueResult:
        """Best-effort observation which can never block the legacy owner."""
        if self.client is None:
            return EnqueueResult(True, "shadow_unconfigured")
        try:
            remote = self._upsert(envelope, "SHADOW")
            if not remote.accepted:
                self.log.warning("shadow upsert rejected task=%s reason=%s",
                                 envelope.task_key, remote.reason)
                return EnqueueResult(True, "shadow_rejected:" + remote.reason)
            return remote
        except Exception as exc:  # shadow must never affect the legacy production path
            self.log.warning("shadow upsert failed task=%s error=%s",
                             envelope.task_key, type(exc).__name__)
            return EnqueueResult(True, "shadow_failed")

    def _observe_managed(self, envelope: TaskEnvelope) -> EnqueueResult:
        """Persist a managed observation, mapping every error fail-closed."""
        if self.client is None:
            return EnqueueResult(False, "control_plane_unconfigured")
        try:
            return self._upsert(envelope, "MANAGED")
        except ControlPlaneUnavailable:
            return EnqueueResult(False, "control_plane_unavailable")
        except ControlPlaneError as exc:
            reason = "control_plane_rejected"
            if getattr(exc, "code", ""):
                reason += ":" + str(exc.code)
            return EnqueueResult(False, reason)
        except Exception as exc:  # fail closed on client/programming surprises too
            self.log.exception("managed upsert failed closed task=%s error=%s",
                               envelope.task_key, type(exc).__name__)
            return EnqueueResult(False, "control_plane_error")

    def observe(self, envelope: TaskEnvelope) -> EnqueueResult:
        """Persist the newest desired revision without starting local work.

        Scan sensors call this before local claimed/active filters.  That keeps
        the desired state current even when the legacy scheduler intentionally
        skips dispatching the item.
        """
        if not isinstance(envelope, TaskEnvelope):
            raise TypeError("envelope must be a TaskEnvelope")
        if self.mode == "legacy":
            return EnqueueResult(True, "legacy_noop")
        if self._is_managed(envelope):
            return self._observe_managed(envelope)
        return self._observe_shadow(envelope)

    def enqueue(self, envelope: TaskEnvelope,
                legacy_submit: Optional[LegacySubmit] = None) -> EnqueueResult:
        observed = self.observe(envelope)
        if self._is_managed(envelope):
            # Managed ownership is fail-closed.  Never call legacy_submit from
            # this branch, even on timeout/5xx/invalid response/conflict.
            return observed
        # Legacy, shadow, and non-allow-listed managed kinds remain locally
        # owned.  ``observe`` has already done the one (and only one) upsert.
        return self._legacy(legacy_submit)
