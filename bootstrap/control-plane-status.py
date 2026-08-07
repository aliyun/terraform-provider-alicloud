#!/usr/bin/env python3
"""Jarvis 控制面诊断与显式恢复 CLI（bootstrap/control-plane-status.sh 的实现体）。

查询 AutomationAgent Jarvis 数据面，供人工按 Aone ID 全链路排查和受保护恢复：
  workers           worker 总览表（key/client/activityStatus/assignment aone id）
  ready             READY 任务及 eligibleWorkerCount=0 的具体原因
  task <aone_id>    单工单全链路（task 状态/current session/fence/最近 5 条 event/
                    operations 回执状态）
  operation <id>    单个 operation + Task/Session/fence/readbackSpec（只读）

凭证由 wrapper 经共享 machine runtime loader 加载。普通诊断/恢复命令使用
JARVIS_CONTROL_PLANE_TOKEN（可由 JARVIS_HTML_REPORT_TOKEN 回填）；
``unresolvable-source-cleanup`` 只使用独立的
JARVIS_CONTROL_PLANE_ADMIN_TOKEN。``legacy-cleanup`` 是零网络墓碑命令。
控制面地址可显式覆盖，默认指向生产控制面。
本文件只读环境变量、不读 env 文件。``discard-resume``、``force-release`` 和
``force-redispatch`` 是写操作，且必须显式 ``--yes``；
``unresolvable-source-cleanup`` 只有带 ``--yes`` 时才写，并额外要求非空 reason。

退出码：0=成功；1=控制面无该工单任务；2=参数/凭证问题；3=控制面请求失败；
4=Task 输入不满足跨 Runtime 契约。
"""

import argparse
import hashlib
import os
import sys
from pathlib import Path

# 追加（而非前插）bridge 目录：保持 PYTHONPATH 可以覆盖 jarvis_task_client，
# 供 test/control_plane_status_test.sh 注入 stub client 做离线测试。
sys.path.append(str(Path(__file__).resolve().parent.parent / "bridge"))

from jarvis_task_client import (  # noqa: E402
    ControlPlaneClient,
    ControlPlaneError,
    ControlPlaneUnavailable,
    DEFAULT_CONTROL_PLANE_BASE_URL,
    InvalidResponse,
)
from task_input_contract import (  # noqa: E402
    TaskInputContractError,
    portable_replacement_for_redispatch,
)

EVENT_TAIL = 5  # task 视图只列最近 N 条 event（全量看服务端 timeline）
MAX_SIGNED_INT64 = (1 << 63) - 1
CLEANUP_SNAPSHOT_COUNT_FIELDS = (
    "tasks",
    "sessions",
    "events",
    "operations",
    "pendingRequiredOperations",
    "blockingRequiredOperations",
)


def _client(command):
    base = (os.environ.get("JARVIS_CONTROL_PLANE_BASE_URL", "").strip()
            or DEFAULT_CONTROL_PLANE_BASE_URL)
    token = os.environ.get("JARVIS_CONTROL_PLANE_TOKEN", "").strip()
    admin_token = os.environ.get("JARVIS_CONTROL_PLANE_ADMIN_TOKEN", "").strip()
    timeout = float(os.environ.get("JARVIS_CONTROL_PLANE_TIMEOUT", "10"))
    if command == "unresolvable-source-cleanup":
        if not admin_token:
            sys.stderr.write(
                "error: JARVIS_CONTROL_PLANE_ADMIN_TOKEN is not configured; "
                "unresolvable-source-cleanup never falls back to the worker token\n")
            raise SystemExit(2)
        # This command only calls /admin/**. Do not retain the ordinary worker
        # token in its client, even when both credentials exist in the runtime.
        return ControlPlaneClient(
            base, admin_token=admin_token, timeout=timeout)
    if not token:
        sys.stderr.write(
            "error: JARVIS_CONTROL_PLANE_TOKEN is not configured "
            "(JARVIS_HTML_REPORT_TOKEN is reused as fallback)\n")
        raise SystemExit(2)
    # Non-admin commands do not retain the privileged credential. Their HTTP
    # requests can only authenticate with the ordinary control-plane token.
    return ControlPlaneClient(base, token, timeout=timeout)


def _control_plane_task_id(value):
    try:
        task_id = int(value)
    except (TypeError, ValueError) as exc:
        raise argparse.ArgumentTypeError(
            "TASK_ID must be a positive control-plane Task ID (not an Aone ID)"
        ) from exc
    if task_id <= 0 or task_id > MAX_SIGNED_INT64:
        raise argparse.ArgumentTypeError(
            "TASK_ID must be a positive signed int64 control-plane Task ID")
    return task_id


def _validate_cleanup_arguments(task_ids, reason, yes):
    if len(task_ids) > 200:
        return "at most 200 control-plane Task IDs may be cleaned at once"
    if len(set(task_ids)) != len(task_ids):
        return "control-plane Task IDs must be unique"
    if yes and not str(reason or "").strip():
        return "--yes requires a nonblank --reason TEXT"
    return ""


def _validate_cleanup_preview(preview, requested_task_ids):
    """Return the complete CAS tuple or a reason why deletion must not proceed."""
    if not isinstance(preview, dict):
        raise ValueError("preview response must be an object")
    required_root = ("snapshot", "activeTasks", "activeSessions", "executable", "taskIds")
    missing_root = [field for field in required_root if field not in preview]
    if missing_root:
        raise ValueError(
            "preview response is incomplete; missing %s" % ", ".join(missing_root))
    snapshot = preview["snapshot"]
    if not isinstance(snapshot, dict):
        raise ValueError("preview.snapshot must be an object")
    missing_snapshot = [
        field for field in CLEANUP_SNAPSHOT_COUNT_FIELDS
        if field not in snapshot
    ]
    if "taskIdsDigest" not in snapshot:
        missing_snapshot.append("taskIdsDigest")
    if missing_snapshot:
        raise ValueError(
            "preview.snapshot is incomplete; missing %s"
            % ", ".join(missing_snapshot))
    for field in CLEANUP_SNAPSHOT_COUNT_FIELDS:
        value = snapshot[field]
        if isinstance(value, bool) or not isinstance(value, int) or value < 0:
            raise ValueError(
                "preview.snapshot.%s must be a non-negative integer" % field)
    digest = snapshot["taskIdsDigest"]
    if not isinstance(digest, str) or not digest.strip():
        raise ValueError("preview.snapshot.taskIdsDigest must not be blank")
    for field in ("activeTasks", "activeSessions"):
        value = preview[field]
        if isinstance(value, bool) or not isinstance(value, int) or value < 0:
            raise ValueError("preview.%s must be a non-negative integer" % field)
    if not isinstance(preview["executable"], bool):
        raise ValueError("preview.executable must be a boolean")
    task_ids = preview["taskIds"]
    if not isinstance(task_ids, list) or not task_ids:
        raise ValueError("preview.taskIds must be a nonempty array")
    if any(isinstance(value, bool) or not isinstance(value, int)
           or value <= 0 or value > MAX_SIGNED_INT64 for value in task_ids):
        raise ValueError("preview.taskIds must contain positive signed int64 values")
    if len(set(task_ids)) != len(task_ids) or task_ids != sorted(task_ids):
        raise ValueError("preview.taskIds must be sorted and unique")
    if sorted(requested_task_ids) != task_ids:
        raise ValueError(
            "preview.taskIds does not match the requested control-plane Task IDs")
    if snapshot["tasks"] != len(task_ids):
        raise ValueError("preview.snapshot.tasks does not match preview.taskIds")
    canonical = ",".join(str(task_id) for task_id in task_ids).encode("utf-8")
    expected_digest = hashlib.sha256(canonical).hexdigest()
    if snapshot["taskIdsDigest"] != expected_digest:
        raise ValueError(
            "preview.snapshot.taskIdsDigest does not match canonical preview.taskIds")
    return snapshot, task_ids


def _validate_cleanup_result(result, expected_snapshot, expected_task_ids):
    """Validate the server's destructive response before claiming deletion."""
    if not isinstance(result, dict):
        raise ValueError("cleanup response must be an object")
    before = result.get("before")
    after = result.get("after")
    deleted = result.get("deletedTaskIds")
    if not isinstance(before, dict) or before != expected_snapshot:
        raise ValueError("cleanup.before does not equal the preview snapshot")
    if not isinstance(after, dict):
        raise ValueError("cleanup.after must be an object")
    missing_after = [
        field for field in CLEANUP_SNAPSHOT_COUNT_FIELDS
        if field not in after
    ]
    if "taskIdsDigest" not in after:
        missing_after.append("taskIdsDigest")
    if missing_after:
        raise ValueError(
            "cleanup.after is incomplete; missing %s" % ", ".join(missing_after))
    if any(isinstance(after[field], bool) or not isinstance(after[field], int)
           or after[field] != 0 for field in CLEANUP_SNAPSHOT_COUNT_FIELDS):
        raise ValueError("cleanup.after counts did not converge to zero")
    if after["taskIdsDigest"] != expected_snapshot["taskIdsDigest"]:
        raise ValueError(
            "cleanup.after.taskIdsDigest does not preserve the canonical target digest")
    if deleted != expected_task_ids:
        raise ValueError(
            "cleanup.deletedTaskIds does not equal the normalized preview taskIds")
    return before, after, deleted


def _trunc(text, limit):
    text = str(text)
    return text if len(text) <= limit else text[: limit - 1] + "…"


def _print_table(header, rows):
    widths = [len(h) for h in header]
    for row in rows:
        for i, cell in enumerate(row):
            widths[i] = max(widths[i], len(str(cell)))
    fmt = "  ".join("%%-%ds" % w for w in widths)
    print(fmt % header)
    for row in rows:
        print(fmt % tuple(str(c) for c in row))


def cmd_workers(client):
    workers = client.list_workers()
    if not isinstance(workers, list):
        sys.stderr.write("error: unexpected /workers response type %s\n"
                         % type(workers).__name__)
        return 3
    rows = []
    for entry in workers:
        if not isinstance(entry, dict):
            continue
        worker = entry.get("worker") if isinstance(entry.get("worker"), dict) else {}
        caps = worker.get("capabilities")
        caps = caps if isinstance(caps, dict) else {}
        assigns = []
        for a in entry.get("assignments") or []:
            if not isinstance(a, dict):
                continue
            task = a.get("task") if isinstance(a.get("task"), dict) else {}
            sess = a.get("session") if isinstance(a.get("session"), dict) else {}
            label = task.get("aoneId") or task.get("taskKey") or "?"
            assigns.append("#%s(s%s f%s)" % (label, sess.get("id", "?"),
                                             sess.get("fenceToken", "?")))
        rows.append((
            _trunc(worker.get("workerKey") or "?", 46),
            str(caps.get("client") or caps.get("workerMode") or "-"),
            str(entry.get("activityStatus") or "?"),
            str(worker.get("lastHeartbeatAt") or "-"),
            " ".join(assigns) or "-",
        ))
    _print_table(("WORKER", "CLIENT", "ACTIVITY", "LAST-HEARTBEAT", "ASSIGNMENTS"), rows)
    print("%d worker(s)" % len(rows))
    return 0


def cmd_ready(client, limit):
    diagnostics = client.list_ready_task_diagnostics(limit=limit)
    if not isinstance(diagnostics, list):
        sys.stderr.write("error: unexpected /tasks/ready-diagnostics response type %s\n"
                         % type(diagnostics).__name__)
        return 3
    rows = []
    for entry in diagnostics:
        if not isinstance(entry, dict):
            continue
        task = entry.get("task") if isinstance(entry.get("task"), dict) else {}
        rows.append((
            str(task.get("id") or "?"),
            str(task.get("aoneId") or "-"),
            _trunc(task.get("taskType") or "-", 22),
            str(task.get("recoveryPolicy") or "-"),
            str(entry.get("eligibleWorkerCount", "?")),
            str(entry.get("reasonCode") or "?"),
            _trunc(entry.get("requiredWorkerKey") or "-", 38),
            str(entry.get("requiredWorkerActivityStatus") or "-"),
        ))
    _print_table(("TASK", "AONE", "TYPE", "RECOVERY", "ELIGIBLE", "REASON",
                  "REQUIRED-WORKER", "ACTIVITY"), rows)
    blocked = sum(1 for row in rows if row[4] == "0")
    print("%d READY task(s), %d without eligible worker" % (len(rows), blocked))
    return 0


def _print_task(client, task):
    print("== task %s %s ==" % (task.get("id", "?"), task.get("taskKey") or "?"))
    print("  status=%s  mode=%s  generation=%s  retry=%s/%s"
          % (task.get("status"), task.get("executionMode"), task.get("generation"),
             task.get("retryCount"), task.get("maxRetries")))
    print("  desired=%s  processing=%s  processed=%s"
          % (task.get("desiredRevision"), task.get("processingRevision"),
             task.get("processedRevision")))
    if task.get("lastError"):
        print("  lastError=%s" % _trunc(task.get("lastError"), 160))

    tid = task.get("id")
    tl = client.get_task_timeline(str(tid)) if tid is not None else None
    if not isinstance(tl, dict):
        print("  (timeline unavailable)")
        return
    cw = tl.get("currentWorker")
    sessions = [s for s in (tl.get("sessions") or []) if isinstance(s, dict)]
    current = next((s for s in sessions
                    if s.get("id") == task.get("currentSessionId")), None)
    if current is not None:
        print("  current session: id=%s status=%s fence=%s attempt=%s resume=%s"
              % (current.get("id"), current.get("status"), current.get("fenceToken"),
                 current.get("attemptNo"), current.get("resumeCount")))
        print("    lease=%s worker=%s runtime=%s"
              % (current.get("leaseExpireAt") or "-",
                 (cw or {}).get("workerKey") if isinstance(cw, dict) else "-",
                 _trunc(current.get("runtimeSessionId") or "-", 60)))
        if current.get("waitType"):
            print("    wait: type=%s key=%s cursor=%s expire=%s"
                  % (current.get("waitType"), current.get("waitKey"),
                     current.get("waitCursor"), current.get("waitExpireAt")))
    else:
        print("  current session: none (%d session(s) total)" % len(sessions))

    events = [e for e in (tl.get("events") or []) if isinstance(e, dict)]
    tail = events[-EVENT_TAIL:]
    print("  recent events (last %d of %d):" % (len(tail), len(events)))
    for e in tail:
        print("    %s %s %s→%s actor=%s%s"
              % (e.get("occurredAt") or "-", e.get("eventType"),
                 e.get("fromStatus") or "-", e.get("toStatus") or "-",
                 e.get("actor") or "-",
                 (" detail=%s" % _trunc(e.get("detail"), 80)) if e.get("detail") else ""))

    operations = [o for o in (tl.get("operations") or []) if isinstance(o, dict)]
    print("  operations (%d):" % len(operations))
    for o in operations:
        print("    op %s %s %s target=%s%s%s"
              % (o.get("id"), o.get("operationType"), o.get("status"),
                 o.get("target") or "-",
                 (" ref=%s" % o.get("externalRef")) if o.get("externalRef") else "",
                 (" lastError=%s" % _trunc(o.get("lastError"), 80))
                 if o.get("lastError") else ""))


def cmd_task(client, aone_id):
    tasks = client.get_task_by_aone(str(aone_id))
    if isinstance(tasks, dict):
        tasks = [tasks]
    tasks = [t for t in (tasks or []) if isinstance(t, dict)]
    if not tasks:
        print("no control-plane task for aone #%s" % aone_id)
        return 1
    for task in tasks:
        _print_task(client, task)
    return 0


def cmd_operation(client, operation_id):
    point = client.get_operation(operation_id)
    if not isinstance(point, dict):
        sys.stderr.write("error: unexpected operation point-read response\n")
        return 3
    operation = point.get("operation") if isinstance(point.get("operation"), dict) else {}
    task = point.get("task") if isinstance(point.get("task"), dict) else {}
    session = point.get("session") if isinstance(point.get("session"), dict) else {}
    worker = point.get("worker") if isinstance(point.get("worker"), dict) else {}
    readback = point.get("readbackSpec") if isinstance(point.get("readbackSpec"), dict) else {}
    print("operation=%s type=%s status=%s target=%s" % (
        operation.get("id", operation_id), operation.get("operationType"),
        operation.get("status"), operation.get("target") or "-"))
    print("task=%s state=%s generation=%s session=%s sessionState=%s fence=%s" % (
        task.get("id", "-"), task.get("status", "-"), task.get("generation", "-"),
        session.get("id", "-"), session.get("status", "-"),
        session.get("fenceToken", "-")))
    print("worker=%s lease=%s heartbeat=%s" % (
        worker.get("workerKey", "-"), session.get("leaseExpireAt", "-"),
        session.get("lastHeartbeatAt", "-")))
    print("readbackSpec=%s" % readback)
    return 0


def cmd_discard_resume(client, task_id, session_id, reason, yes):
    if not yes:
        sys.stderr.write(
            "error: discard-resume cancels the exact resumable session; "
            "review `task <aone_id>` first and pass --yes\n")
        return 2
    result = client.discard_resume_context(
        str(task_id), int(session_id), reason,
        request_id="discard-resume:%s:%s" % (task_id, session_id))
    print("discarded resume context: task=%s session=%s status=%s"
          % (result.get("id", task_id), session_id, result.get("status", "?")))
    return 0


def _force_release_snapshot(timeline, task_id, session_id):
    if not isinstance(timeline, dict):
        raise ValueError("unexpected task timeline response")
    task = timeline.get("task")
    if not isinstance(task, dict):
        raise ValueError("task timeline does not contain a task snapshot")
    observed_task_id = task.get("id")
    if observed_task_id is not None and str(observed_task_id) != str(task_id):
        raise ValueError(
            "task timeline returned task %s while %s was requested"
            % (observed_task_id, task_id))

    selected_session_id = (
        session_id if session_id is not None else task.get("currentSessionId"))
    selected_session = None
    if selected_session_id is not None:
        selected_session = next((
            candidate for candidate in (timeline.get("sessions") or [])
            if isinstance(candidate, dict)
            and str(candidate.get("id")) == str(selected_session_id)
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
    missing = [key for key, value in required.items() if value is None]
    if missing:
        raise ValueError(
            "task timeline is missing force-release CAS field(s): %s"
            % ", ".join(missing))

    return {
        "expected_task_status": required["status"],
        "expected_session_id": selected_session_id,
        "expected_session_status": (
            selected_session.get("status")
            if selected_session is not None else None),
        "expected_generation": required["generation"],
        "expected_state_version": required["stateVersion"],
        "expected_fence_token": (
            selected_session.get("fenceToken")
            if selected_session is not None else None),
        "expected_retry_count": required["retryCount"],
        "expected_desired_revision": task.get("desiredRevision"),
        "expected_processing_revision": task.get("processingRevision"),
        "expected_worker_key": (
            selected_session.get("historicalWorkerKey")
            if selected_session is not None else None),
        "expected_worker_id": (
            selected_session.get("historicalWorkerId")
            if selected_session is not None else None),
        "expected_worker_process_uuid": (
            selected_session.get("historicalWorkerProcessUuid")
            if selected_session is not None else None),
    }


def cmd_force_release(client, task_id, session_id, reason, yes):
    if not yes:
        sys.stderr.write(
            "error: force-release fences and detaches the exact reviewed ownership; "
            "pass --yes to fetch a fresh CAS snapshot and continue\n")
        return 2
    try:
        timeline = client.get_task_timeline(str(task_id))
        snapshot = _force_release_snapshot(timeline, task_id, session_id)
        result = client.force_release_task(
            str(task_id), reason=reason, **snapshot)
    except ControlPlaneError as exc:
        if getattr(exc, "status", None) == 404:
            sys.stderr.write(
                "error: force-release endpoint is unavailable (HTTP 404); "
                "the control-plane server version has not been deployed\n")
            return 3
        raise
    except (TypeError, ValueError) as exc:
        sys.stderr.write("error: cannot build force-release CAS: %s\n" % exc)
        return 3

    task = result.get("task") if isinstance(result.get("task"), dict) else {}
    old_generation = result.get(
        "previousGeneration", snapshot["expected_generation"])
    new_generation = task.get("generation", result.get("generation", "?"))
    old_session = result.get(
        "releasedSessionId", snapshot["expected_session_id"])
    final_status = task.get("status", result.get("status", "?"))
    print(
        "force release: action=%s oldGeneration=%s newGeneration=%s "
        "oldSession=%s finalStatus=%s"
        % (
            result.get("action", "?"),
            old_generation,
            new_generation,
            old_session if old_session is not None else "-",
            final_status,
        ))
    if result.get("message"):
        print("  message=%s" % result.get("message"))
    return 0


def cmd_force_redispatch(
        client, task_id, session_id, reason, yes, target_worker_key,
        target_host_id, target_runtime):
    if not yes:
        sys.stderr.write(
            "error: force-redispatch fences the reviewed ownership and targets "
            "another online queue worker; pass --yes to fetch a fresh CAS "
            "snapshot and continue\n")
        return 2
    try:
        timeline = client.get_task_timeline(str(task_id))
        snapshot = _force_release_snapshot(timeline, task_id, session_id)
        replacement = portable_replacement_for_redispatch(
            timeline,
            target_runtime=target_runtime,
            selected_session_id=snapshot["expected_session_id"],
        )
        result = client.force_redispatch_task(
            str(task_id),
            target_worker_key=target_worker_key,
            target_host_id=target_host_id,
            target_runtime=target_runtime,
            portable_input_replacement=replacement,
            reason=reason,
            **snapshot)
    except ControlPlaneError as exc:
        if getattr(exc, "status", None) == 404:
            sys.stderr.write(
                "error: force-redispatch endpoint is unavailable (HTTP 404); "
                "the control-plane server version has not been deployed\n")
            return 3
        raise
    except TaskInputContractError as exc:
        sys.stderr.write(
            "error: force-redispatch blocked by Task input contract: "
            "%s\n" % exc)
        return 4
    except (TypeError, ValueError) as exc:
        sys.stderr.write(
            "error: cannot build force-redispatch CAS: %s\n" % exc)
        return 3

    task = result.get("task") if isinstance(result.get("task"), dict) else {}
    target = (
        result.get("targetWorker")
        if isinstance(result.get("targetWorker"), dict) else {})
    old_generation = result.get(
        "previousGeneration", snapshot["expected_generation"])
    new_generation = task.get("generation", result.get("generation", "?"))
    old_session = result.get(
        "releasedSessionId", snapshot["expected_session_id"])
    final_status = task.get("status", result.get("status", "?"))
    print(
        "cross-machine redispatch: action=%s targetRuntime=%s oldGeneration=%s "
        "newGeneration=%s oldSession=%s finalStatus=%s"
        % (
            result.get("action", "?"),
            target_runtime,
            old_generation,
            new_generation,
            old_session if old_session is not None else "-",
            final_status,
        ))
    print(
        "  targetWorkerId=%s targetWorkerKey=%s hostId=%s processUuid=%s "
        "activity=%s"
        % (
            target.get("id", "-"),
            target.get("workerKey", "-"),
            target.get("hostId", "-"),
            target.get("processUuid", "-"),
            target.get("activityStatus", target.get("status", "-")),
        ))
    print(
        "  dispatch=targeted/queued; READY does not mean the target worker "
        "has started a new session")
    if replacement is not None:
        print(
            "  input=rehydrated PORTABLE_V1 sourceDigest=%s replacementDigest=%s"
            % (
                replacement["expectedSourceInputDigest"][:16],
                replacement["replacementDigest"][:16],
            ))
    if result.get("message"):
        print("  message=%s" % result.get("message"))
    return 0


def cmd_legacy_cleanup():
    sys.stderr.write(
        "error: legacy-cleanup was retired and sent no HTTP request; use "
        "unresolvable-source-cleanup CONTROL_PLANE_TASK_ID... to preview, then "
        "repeat with --reason TEXT --yes\n")
    return 2


def cmd_unresolvable_source_cleanup(client, task_ids, reason, yes):
    preview = client.preview_unresolvable_source_cleanup(task_ids)
    try:
        snapshot, normalized_ids = _validate_cleanup_preview(preview, task_ids)
    except ValueError as exc:
        sys.stderr.write(
            "error: malformed unresolvable-source cleanup preview; deletion was not "
            "attempted: %s\n" % exc)
        return 3

    print("unresolvable-source cleanup preview:")
    print("  controlPlaneTaskIds=%s (control-plane Task IDs, not Aone IDs)" %
          ",".join(str(task_id) for task_id in normalized_ids))
    print(
        "  tasks=%s sessions=%s events=%s operations=%s "
        "pendingRequiredOperations=%s blockingRequiredOperations=%s"
        % tuple(snapshot[field] for field in CLEANUP_SNAPSHOT_COUNT_FIELDS))
    print("  activeTasks=%s activeSessions=%s executable=%s digest=%s" % (
        preview["activeTasks"], preview["activeSessions"],
        preview["executable"], snapshot["taskIdsDigest"]))

    if not yes:
        print(
            "preview only; no rows deleted. After review, repeat with "
            "--reason TEXT --yes.")
        return 0

    blockers = (
        preview["activeTasks"],
        preview["activeSessions"],
        snapshot["blockingRequiredOperations"],
    )
    if not preview["executable"] or any(blockers):
        sys.stderr.write(
            "error: cleanup blocked; requires executable=true, activeTasks=0, "
            "activeSessions=0, and blockingRequiredOperations=0\n")
        return 3

    try:
        result = client.cleanup_unresolvable_source_tasks(
            task_ids=normalized_ids,
            expected_snapshot=snapshot,
            reason=str(reason).strip(),
        )
    except (ControlPlaneUnavailable, InvalidResponse) as exc:
        sys.stderr.write(
            "error: cleanup outcome is post-uncertain (%s); do not replay because "
            "the server has no idempotent receipt—inspect the Task IDs first\n" % exc)
        return 3
    try:
        before, after, deleted = _validate_cleanup_result(
            result, snapshot, normalized_ids)
    except ValueError as exc:
        sys.stderr.write(
            "error: cleanup response is post-uncertain (%s); do not replay because "
            "the server has no idempotent receipt—inspect the Task IDs first\n" % exc)
        return 3
    print("deleted control-plane Task IDs: %s" % (
        ",".join(str(task_id) for task_id in deleted) or "-"))
    print("  before.tasks=%s after.tasks=%s" % (
        before.get("tasks", "?"), after.get("tasks", "?")))
    return 0


def main(argv=None):
    parser = argparse.ArgumentParser(
        prog="control-plane-status",
        description="Jarvis 控制面排查与受保护恢复 CLI（workers / READY 派发诊断 / 单工单全链路）")
    sub = parser.add_subparsers(dest="cmd", required=True)
    sub.add_parser("workers", help="list every registered worker with assignments")
    p_ready = sub.add_parser("ready", help="list READY tasks with dispatch eligibility reasons")
    p_ready.add_argument("--limit", type=int, default=100,
                         help="maximum READY tasks to return (1-500, default: 100)")
    p_task = sub.add_parser("task", help="full chain for one Aone work item")
    p_task.add_argument("aone_id", help="Aone work item id, e.g. 84386065")
    p_operation = sub.add_parser("operation", help="point-read one operation")
    p_operation.add_argument("operation_id", type=int)
    p_discard = sub.add_parser(
        "discard-resume",
        help="cancel one exact legacy resumable session after operator review")
    p_discard.add_argument("task_id", type=int, help="control-plane Task id")
    p_discard.add_argument("session_id", type=int, help="expected current Session id")
    p_discard.add_argument("--reason", required=True, help="auditable recovery reason")
    p_discard.add_argument("--yes", action="store_true", help="confirm context discard")
    p_release = sub.add_parser(
        "force-release",
        help="force-release one exact ownership snapshot after fresh timeline review")
    p_release.add_argument("task_id", type=int, help="control-plane Task id")
    p_release.add_argument(
        "session_id", type=int, nargs="?",
        help="expected Session id (default: fresh task.currentSessionId)")
    p_release.add_argument("--reason", required=True, help="auditable release reason")
    p_release.add_argument("--yes", action="store_true", help="confirm ownership release")
    p_redispatch = sub.add_parser(
        "force-redispatch",
        help="release one ownership snapshot and target an online worker on another host")
    p_redispatch.add_argument("task_id", type=int, help="control-plane Task id")
    p_redispatch.add_argument(
        "session_id", type=int, nargs="?",
        help="expected Session id (default: fresh task.currentSessionId)")
    target = p_redispatch.add_mutually_exclusive_group(required=True)
    target.add_argument(
        "--auto-target", action="store_true",
        help="let the server select an eligible online worker on another host")
    target.add_argument(
        "--target-worker", metavar="WORKER_KEY",
        help="require this exact online Worker key")
    target.add_argument(
        "--target-host", metavar="HOST_ID",
        help="select an eligible online Worker on this exact machine")
    p_redispatch.add_argument(
        "--target-runtime", required=True,
        choices=("INTERACTIVE", "PERSISTENT"),
        help="runtime mode that must receive the portable Task")
    p_redispatch.add_argument(
        "--reason", required=True, help="auditable redispatch reason")
    p_redispatch.add_argument(
        "--yes", action="store_true", help="confirm cross-machine redispatch")
    p_legacy = sub.add_parser(
        "legacy-cleanup",
        help="retired cleanup command; prints migration usage and sends no HTTP")
    p_legacy.add_argument("--yes", action="store_true",
                          help=argparse.SUPPRESS)
    p_cleanup = sub.add_parser(
        "unresolvable-source-cleanup",
        help="preview/delete explicitly selected control-plane Tasks whose source is unresolvable")
    p_cleanup.add_argument(
        "task_ids", nargs="+", type=_control_plane_task_id, metavar="TASK_ID",
        help="positive control-plane Task ID (not an Aone work-item ID; max 200)")
    p_cleanup.add_argument(
        "--reason", help="nonblank audit reason (required with --yes)")
    p_cleanup.add_argument(
        "--yes", action="store_true",
        help="delete only the complete, executable preview snapshot")
    args = parser.parse_args(argv)
    if args.cmd == "legacy-cleanup":
        return cmd_legacy_cleanup()
    if args.cmd == "unresolvable-source-cleanup":
        argument_error = _validate_cleanup_arguments(
            args.task_ids, args.reason, args.yes)
        if argument_error:
            sys.stderr.write("error: %s\n" % argument_error)
            return 2
    client = _client(args.cmd)
    try:
        if args.cmd == "workers":
            return cmd_workers(client)
        if args.cmd == "ready":
            return cmd_ready(client, args.limit)
        if args.cmd == "task":
            return cmd_task(client, args.aone_id)
        if args.cmd == "operation":
            return cmd_operation(client, args.operation_id)
        if args.cmd == "unresolvable-source-cleanup":
            return cmd_unresolvable_source_cleanup(
                client, args.task_ids, args.reason, args.yes)
        if args.cmd == "force-release":
            return cmd_force_release(
                client, args.task_id, args.session_id, args.reason, args.yes)
        if args.cmd == "force-redispatch":
            return cmd_force_redispatch(
                client, args.task_id, args.session_id, args.reason, args.yes,
                args.target_worker, args.target_host, args.target_runtime)
        return cmd_discard_resume(
            client, args.task_id, args.session_id, args.reason, args.yes)
    except (TypeError, ValueError) as e:
        sys.stderr.write("error: %s\n" % e)
        return 2
    except ControlPlaneError as e:
        sys.stderr.write("error: %s\n" % e)
        return 3


if __name__ == "__main__":
    sys.exit(main())
