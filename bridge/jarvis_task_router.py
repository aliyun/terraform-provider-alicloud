#!/usr/bin/env python3
"""Route Jarvis work by recovery semantics, with rollout compatibility.

The stable domain decision is deliberately small:

* ``TASK`` for work that must survive interruption and resume elsewhere.
* ``EPHEMERAL_JOB`` for disposable local work.

``legacy``, ``shadow``, and ``managed`` are only migration policies controlling
where a ``TASK`` is observed/executed during rollout.  They are not task kinds.
The compatibility ``legacy_submit`` callback lets existing schedulers migrate
without importing the large bridge module here.
"""

from __future__ import annotations

import logging
import os
from typing import Callable, Iterable, NamedTuple, Optional, Tuple

from jarvis_task_client import (
    ControlPlaneClient,
    ControlPlaneError,
    ControlPlaneUnavailable,
    TaskEnvelope,
)


class EnqueueResult(NamedTuple):
    """Compatibility result for the current enqueue integration seam."""

    accepted: bool
    reason: str


TASK = "TASK"
EPHEMERAL_JOB = "EPHEMERAL_JOB"
DEFAULT_TASK_TYPES = (
    "ticket",
    "revisit",
    "persona",
    "wake",
    "pr_ci_fix",
    "pr_comment_reply",
)


class ExecutionRoute(NamedTuple):
    """Stable domain classification, independent of rollout policy."""

    execution_kind: str
    needs_recovery: bool


LegacySubmit = Callable[[], Tuple[bool, str]]


class ExecutionRouter:
    """Classify work, then apply a temporary migration policy.

    Domain classification:
      * configured ``task_types`` become ``TASK`` (``needs_recovery=True``).
      * all other kinds become ``EPHEMERAL_JOB``.

    Migration policies:
      * ``legacy``: local submit only.
      * ``shadow``: best-effort SHADOW upsert, then local submit regardless of
        control-plane health.  SHADOW tasks must never be leaseable server-side.
      * ``managed``: ``TASK`` work is owned only by the data plane.  A failed
        upsert is fail-closed and never falls back locally.  ``EPHEMERAL_JOB``
        work stays process-local and never creates control-plane state.

    ``mode``, ``managed_task_types``, ``is_managed()``, and ``TaskRouter`` remain
    compatibility names while call sites migrate to ``migration_policy``,
    ``task_types``, ``is_task()``, and ``ExecutionRouter``.
    """

    VALID_MIGRATION_POLICIES = {"legacy", "shadow", "managed"}
    VALID_MODES = VALID_MIGRATION_POLICIES

    def __init__(self, mode: Optional[str] = None, *,
                 migration_policy: Optional[str] = None,
                 client: Optional[ControlPlaneClient] = None,
                 task_types: Optional[Iterable[str]] = None,
                 managed_task_types: Optional[Iterable[str]] = None,
                 logger: Optional[logging.Logger] = None):
        if mode is not None and migration_policy is not None:
            raise ValueError("mode and migration_policy are aliases; provide only one")
        configured_policy = migration_policy if migration_policy is not None else mode
        normalized = str(configured_policy or "legacy").strip().lower()
        if normalized not in self.VALID_MIGRATION_POLICIES:
            raise ValueError(
                "invalid migration policy %r (expected legacy/shadow/managed)"
                % configured_policy)
        if task_types is not None and managed_task_types is not None:
            raise ValueError("task_types and managed_task_types are aliases; provide only one")
        configured_types = task_types if task_types is not None else managed_task_types
        self.migration_policy = normalized
        self.client = client
        self.task_types = {
            str(value).strip().lower() for value in (configured_types or [])
            if str(value).strip()
        }
        self.log = logger or logging.getLogger(__name__)

    @property
    def mode(self) -> str:
        """Compatibility alias for ``migration_policy``."""
        return self.migration_policy

    @property
    def managed_task_types(self) -> set:
        """Compatibility alias for the recoverable ``task_types`` allow-list."""
        return self.task_types

    @managed_task_types.setter
    def managed_task_types(self, values: Iterable[str]) -> None:
        self.task_types = {
            str(value).strip().lower() for value in (values or [])
            if str(value).strip()
        }

    @classmethod
    def from_env(cls, *, client: Optional[ControlPlaneClient] = None,
                 logger: Optional[logging.Logger] = None) -> "ExecutionRouter":
        raw_types = os.environ.get("JARVIS_TASK_TYPES")
        if raw_types is None:
            raw_types = os.environ.get("JARVIS_MANAGED_TASK_TYPES")
        if raw_types is None:
            raw_types = ",".join(DEFAULT_TASK_TYPES)
        migration_policy = os.environ.get("JARVIS_TASK_MIGRATION_POLICY")
        if migration_policy is None:
            migration_policy = os.environ.get("JARVIS_TASK_MODE", "legacy")
        return cls(
            migration_policy=migration_policy,
            client=client,
            task_types=raw_types.split(","),
            logger=logger,
        )

    @staticmethod
    def classify(needs_recovery: bool) -> ExecutionRoute:
        """Map the sole domain decision to its stable execution kind."""
        if not isinstance(needs_recovery, bool):
            raise TypeError("needs_recovery must be a bool")
        return ExecutionRoute(TASK if needs_recovery else EPHEMERAL_JOB,
                              needs_recovery)

    def route(self, envelope: TaskEnvelope) -> ExecutionRoute:
        """Return only the durable domain kind for one work envelope."""
        if not isinstance(envelope, TaskEnvelope):
            raise TypeError("envelope must be a TaskEnvelope")
        task_type = envelope.task_type.strip().lower()
        needs_recovery = "*" in self.task_types or task_type in self.task_types
        return self.classify(needs_recovery)

    def is_task(self, envelope: TaskEnvelope) -> bool:
        """Whether the work is a recoverable ``TASK``."""
        return self.route(envelope).needs_recovery

    def _is_managed(self, envelope: TaskEnvelope) -> bool:
        return self.migration_policy == "managed" and self.is_task(envelope)

    def is_managed(self, envelope: TaskEnvelope) -> bool:
        """Compatibility predicate for central ownership during migration."""
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
        desired state current even when the compatibility scheduler skips local
        dispatch.  Rollout policy affects observation, never domain kind.
        """
        if not isinstance(envelope, TaskEnvelope):
            raise TypeError("envelope must be a TaskEnvelope")
        if not self.is_task(envelope):
            # EphemeralJob is intentionally process-local: no Task row, no
            # Session, no lease/fence, and no shadow observation.
            return EnqueueResult(True, "ephemeral_noop")
        if self.migration_policy == "legacy":
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
        # Legacy/shadow TASK work and every EPHEMERAL_JOB remain locally owned.
        # ``observe`` has already performed the only allowed TASK upsert.
        return self._legacy(legacy_submit)


# One-release compatibility alias.  New integrations should import
# ``ExecutionRouter`` and use ``route()/is_task()`` for domain decisions.
TaskRouter = ExecutionRouter
