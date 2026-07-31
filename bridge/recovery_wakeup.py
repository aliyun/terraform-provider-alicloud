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


__all__ = ["enabled", "should_wake", "epoch_key", "wake_redispatch"]
