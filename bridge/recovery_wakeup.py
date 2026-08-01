"""Wake RECOVERY_REQUIRED + retry-exhausted idle Tasks via force-redispatch.

scan 的 ``upsert_desired_task`` 只写 ``desiredRevision``（期望态），不翻转
execution ``status``；``force=True`` 在控制面 Task 路径是空操作（只作用于
EphemeralExecutor 本地软去重）。因此一个 Task 一旦因 ``field_repair apply_timeout``
等失败进入 ``RECOVERY_REQUIRED`` 且 ``retryCount > maxRetries``，就卡死不在 READY
队列，``lease_task`` 永远拉不到——jarvis-idle 工单上的人工评论无人回复。

本 helper 封装 force-redispatch（CAS snapshot + portable input 重建 + request_id
SHA-256 幂等），由 scan / claim_health 在检测到 ``jarvis-idle + 新人工评论`` 且
Task 处于 retry-exhausted RECOVERY_REQUIRED 时调用，把 Task 推回 READY 让
persistent worker 重新 lease。详见 ``docs/execution-architecture.md`` §
「RECOVERY_REQUIRED and the desired-state upsert boundary」。

开关 ``JARVIS_RECOVERY_WAKEUP_ENABLE`` 默认关，灰度后再开。
"""

import logging
import os
from typing import Any, Dict, Mapping, Optional

from bridge.task_input_contract import portable_replacement_for_redispatch

log = logging.getLogger(__name__)


def enabled() -> bool:
    """True only when the wakeup operator opt-in is set (default off)."""
    return os.environ.get("JARVIS_RECOVERY_WAKEUP_ENABLE", "0").strip() not in (
        "0", "", "false", "False", "no", "No")


def should_wake(task: Mapping[str, Any]) -> bool:
    """True if ``task`` is RECOVERY_REQUIRED with retry budget exhausted.

    Mirrors the exhausted判定 in ``claim_health.py`` (RECOVERY_REQUIRED + retryCount
    > maxRetries): a retry-exhausted Task never self-recovers (server keeps it here
    rather than converging), so it blackholes bound to a dead RESUMABLE session.
    """
    if not isinstance(task, Mapping):
        return False
    if str(task.get("status") or "").strip() != "RECOVERY_REQUIRED":
        return False
    retry = task.get("retryCount")
    max_retries = task.get("maxRetries")
    try:
        return retry is not None and max_retries is not None and int(retry) > int(max_retries)
    except (TypeError, ValueError):
        return False


def epoch_key(task: Mapping[str, Any], aone_id: str = "") -> str:
    """Per-aone-id wake epoch: changes when desiredRevision or retryCount advances.

    Callers keep ``{aone_id: epoch_key}`` and skip a new wake while the epoch is
    unchanged, so the same dead Task is not force-redispatched every tick. A retry
    bump or a new desired revision opens a fresh epoch and allows one more attempt.
    """
    return "%s|retry=%s|desired=%s" % (
        aone_id or str(task.get("aoneId") or task.get("id") or ""),
        task.get("retryCount"),
        task.get("desiredRevision") or "")


def _cas_snapshot(timeline: Mapping[str, Any], task_id: str) -> Dict[str, Any]:
    """Extract force-redispatch CAS fields from a Task timeline.

    Mirrors ``bootstrap/control-plane-status.py::_force_release_snapshot`` so the
    CLI script and this helper stay byte-identical on the CAS contract; keep them
    in sync if the server changes the expected_* schema.
    """
    if not isinstance(timeline, Mapping):
        raise ValueError("unexpected task timeline response for %s" % task_id)
    task = timeline.get("task")
    if not isinstance(task, Mapping):
        raise ValueError("task timeline does not contain a task snapshot for %s" % task_id)
    observed = task.get("id")
    if observed is not None and str(observed) != str(task_id):
        raise ValueError(
            "task timeline returned task %s while %s was requested" % (observed, task_id))

    selected_session_id = task.get("currentSessionId")
    selected_session = None
    if selected_session_id is not None:
        selected_session = next((
            s for s in (timeline.get("sessions") or [])
            if isinstance(s, Mapping) and str(s.get("id")) == str(selected_session_id)
        ), None)
        if selected_session is None:
            raise ValueError(
                "task timeline does not contain session %s" % selected_session_id)

    required = {
        "status": task.get("status"),
        "generation": task.get("generation"),
        "stateVersion": task.get("stateVersion"),
        "retryCount": task.get("retryCount"),
    }
    missing = [k for k, v in required.items() if v is None]
    if missing:
        raise ValueError(
            "task timeline is missing force-redispatch CAS field(s): %s"
            % ", ".join(missing))

    return {
        "expected_task_status": required["status"],
        "expected_session_id": selected_session_id,
        "expected_session_status": (
            selected_session.get("status") if selected_session is not None else None),
        "expected_generation": required["generation"],
        "expected_state_version": required["stateVersion"],
        "expected_fence_token": (
            selected_session.get("fenceToken") if selected_session is not None else None),
        "expected_retry_count": required["retryCount"],
        "expected_desired_revision": task.get("desiredRevision"),
        "expected_processing_revision": task.get("processingRevision"),
        "expected_worker_key": (
            selected_session.get("historicalWorkerKey") if selected_session is not None else None),
        "expected_worker_id": (
            selected_session.get("historicalWorkerId") if selected_session is not None else None),
        "expected_worker_process_uuid": (
            selected_session.get("historicalWorkerProcessUuid")
            if selected_session is not None else None),
    }


def wake_redispatch(
        client: Any, task_id: str, aone_id: str, reason: str,
        *, target_runtime: str = "PERSISTENT") -> Dict[str, Any]:
    """Force-redispatch one retry-exhausted RECOVERY_REQUIRED Task to an online worker.

    ``target_worker_key`` / ``target_host_id`` left None so the server selects an
    eligible online queue-pull worker (fresh heartbeat/host/capacity/capacity checks
    server-side). ``request_id`` is auto-generated as SHA-256(task_id+payload) inside
    ``ControlPlaneClient.force_redispatch_task`` for server-side idempotency; callers
    still keep a per-aone-id epoch ledger to avoid re-calling every tick.

    Returns the server response (``action`` / ``previousGeneration`` / ``generation``
    / ``releasedSessionId`` / ``status`` / ``message``). Raises on CAS mismatch,
    transport failure, or a non-mapping timeline.
    """
    timeline = client.get_task_timeline(str(task_id))
    if not isinstance(timeline, Mapping):
        raise ValueError(
            "get_task_timeline(%s) returned non-mapping %s" % (task_id, type(timeline).__name__))
    snapshot = _cas_snapshot(timeline, task_id)
    replacement = portable_replacement_for_redispatch(
        timeline, target_runtime=target_runtime,
        selected_session_id=snapshot["expected_session_id"])
    result = client.force_redispatch_task(
        str(task_id),
        target_worker_key=None,
        target_host_id=None,
        target_runtime=target_runtime,
        portable_input_replacement=replacement,
        reason=reason,
        **snapshot)
    new_gen = (result.get("task") or {}).get("generation", result.get("generation")) if isinstance(result, dict) else None
    final_status = (result.get("task") or {}).get("status", result.get("status")) if isinstance(result, dict) else None
    log.info(
        "recovery_wakeup aone=%s task=%s action=%s gen=%s→%s status=%s",
        aone_id, task_id, result.get("action") if isinstance(result, dict) else None,
        snapshot["expected_generation"], new_gen, final_status)
    return result if isinstance(result, dict) else {"raw": result}


RESUME_CONTEXT_REFS = ("transcriptUri", "branchRef", "checkpointUri")


def resume_context_empty(session: Any) -> bool:
    """True when a RESUMABLE session has no externalized context to resume.

    A Session that died before writing a transcript, a branch or a checkpoint is
    ``RESUMABLE`` in name only — there is literally nothing to resume from. That
    is the common shape for a Task that failed during field repair or claim,
    before the agent ever started, and it is what makes releasing such a Session
    lossless.
    """
    if not isinstance(session, Mapping):
        return False
    return not any(
        str(session.get(ref) or "").strip() for ref in RESUME_CONTEXT_REFS)


def release_stranded(client: Any, task_id: str, reason: str) -> Dict[str, Any]:
    """Unstrand a RESUME_ONLY Task whose RESUMABLE Session has nothing to resume.

    ``force_redispatch_task`` is not usable here: the server requires
    ``REPLAY_SAFE`` and every stranded Task of this shape is ``RESUME_ONLY``, so
    redispatch always answers BLOCKED. Two APIs that *do* accept ``RESUME_ONLY``
    split the work between them, and which one fits depends on whether a newer
    desired revision is already queued:

    - ``desired != processing`` → ``force_release_task``. It fences and cancels
      the dead Session and requeues the pending revision in a fresh generation,
      so the Task returns to READY with a reset retry budget.
    - ``desired == processing`` → ``discard_resume_context``. There is no newer
      revision for force-release to requeue, so it answers
      ``OWNERSHIP_CLEARED_RECOVERY_REQUIRED`` — the dead Session is gone but the
      Task stays parked. discard-resume, by contrast, cancels the Session and
      moves the Task straight back to READY (verified on task 3312). Pointing
      the operator at a discard-resume that needs a Session after we already
      cleared it is what made four tasks stay stranded while still paging.

    The emptiness check is deliberately re-done against this fresh timeline
    rather than trusting the caller's snapshot: releasing a Session that has
    since externalized a transcript would throw away real work. Raises
    ``ValueError`` when the fresh read no longer matches, so the caller keeps
    alerting instead of silently doing nothing.

    Returns the server response (``action`` / ``status`` / ``message``).
    """
    timeline = client.get_task_timeline(str(task_id))
    if not isinstance(timeline, Mapping):
        raise ValueError(
            "get_task_timeline(%s) returned non-mapping %s"
            % (task_id, type(timeline).__name__))
    task = timeline.get("task")
    if not isinstance(task, Mapping):
        raise ValueError("task timeline does not contain a task snapshot for %s" % task_id)
    if str(task.get("recoveryPolicy") or "").upper() != "RESUME_ONLY":
        raise ValueError("task %s is no longer RESUME_ONLY" % task_id)
    session_id = task.get("currentSessionId")
    session = next((
        candidate for candidate in (timeline.get("sessions") or [])
        if isinstance(candidate, Mapping)
        and str(candidate.get("id")) == str(session_id)
    ), None)
    if not isinstance(session, Mapping):
        # No live Session to discard: nothing this path can do. The caller must
        # stop recommending discard-resume, since that API needs a Session.
        raise ValueError("task %s no longer owns session %s" % (task_id, session_id))
    if str(session.get("status") or "").upper() != "RESUMABLE":
        raise ValueError(
            "session %s is %s, not RESUMABLE" % (session_id, session.get("status")))
    if not resume_context_empty(session):
        raise ValueError(
            "session %s now has resume context; refusing to release" % session_id)

    desired = str(task.get("desiredRevision") or "")
    processing = str(task.get("processingRevision") or "")
    if desired and desired != processing:
        # A newer revision is already queued; force-release requeues it in a
        # fresh generation and the Task returns to READY.
        snapshot = _cas_snapshot(timeline, task_id)
        result = client.force_release_task(str(task_id), reason=reason, **snapshot)
    elif session_id is not None:
        # desired == processing: no revision to requeue. force-release would
        # leave the Task parked RECOVERY_REQUIRED; discard-resume cancels the
        # empty Session and moves it to READY (task 3312, 2026-08-01).
        result = client.discard_resume_context(
            str(task_id), int(session_id), reason=reason)
    else:
        raise ValueError(
            "task %s has no session and no pending revision; nothing to release"
            % task_id)
    action = (result or {}).get("action") if isinstance(result, dict) else None
    log.info(
        "recovery_release task=%s session=%s action=%s status=%s desired_pending=%s",
        task_id, session_id, action,
        (result or {}).get("status") if isinstance(result, dict) else None,
        desired != processing)
    return result if isinstance(result, dict) else {"raw": result}


__all__ = [
    "enabled", "should_wake", "epoch_key", "wake_redispatch",
    "resume_context_empty", "release_stranded", "RELEASE_SUCCESS_ACTIONS",
]

# force-release outcomes that mean the dead Session is gone. RELEASED also
# requeued a pending revision in a fresh generation; OWNERSHIP_CLEARED_* only
# cleared ownership because there was no newer revision to queue — the Task stays
# RECOVERY_REQUIRED but stops being an owner-unavailable blocker either way.
RELEASE_SUCCESS_ACTIONS = frozenset({
    "RELEASED",
    "OWNERSHIP_CLEARED_RECOVERY_REQUIRED",
})
