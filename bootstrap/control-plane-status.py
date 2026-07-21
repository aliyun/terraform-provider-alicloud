#!/usr/bin/env python3
"""Jarvis 控制面诊断与显式恢复 CLI（bootstrap/control-plane-status.sh 的实现体）。

查询 AutomationAgent Jarvis 数据面，供人工按 Aone ID 全链路排查和受保护恢复：
  workers           worker 总览表（key/client/activityStatus/assignment aone id）
  ready             READY 任务及 eligibleWorkerCount=0 的具体原因
  task <aone_id>    单工单全链路（task 状态/current session/fence/最近 5 条 event/
                    operations 回执状态）
  operation <id>    单个 operation + Task/Session/fence/readbackSpec（只读）

凭证由 wrapper 经共享 machine runtime loader 加载（token 回退
JARVIS_HTML_REPORT_TOKEN）；控制面地址可显式覆盖，默认指向预发。
本文件只读环境变量、不读 env 文件。只有 ``discard-resume`` 是写操作，且必须显式 ``--yes``。

退出码：0=成功；1=控制面无该工单任务；2=缺 token 配置；3=控制面请求失败。
"""

import argparse
import os
import sys
from pathlib import Path

# 追加（而非前插）bridge 目录：保持 PYTHONPATH 可以覆盖 jarvis_task_client，
# 供 test/control_plane_status_test.sh 注入 stub client 做离线测试。
sys.path.append(str(Path(__file__).resolve().parent.parent / "bridge"))

from jarvis_task_client import ControlPlaneClient, ControlPlaneError  # noqa: E402

EVENT_TAIL = 5  # task 视图只列最近 N 条 event（全量看服务端 timeline）
DEFAULT_CONTROL_PLANE_BASE_URL = "https://pre-agent.aliyun-inc.com"


def _client():
    base = (os.environ.get("JARVIS_CONTROL_PLANE_BASE_URL", "").strip()
            or DEFAULT_CONTROL_PLANE_BASE_URL)
    token = os.environ.get("JARVIS_CONTROL_PLANE_TOKEN", "").strip()
    if not token:
        sys.stderr.write(
            "error: JARVIS_CONTROL_PLANE_TOKEN is not configured "
            "(JARVIS_HTML_REPORT_TOKEN is reused as fallback)\n")
        raise SystemExit(2)
    timeout = float(os.environ.get("JARVIS_CONTROL_PLANE_TIMEOUT", "10"))
    return ControlPlaneClient(base, token, timeout=timeout)


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
    args = parser.parse_args(argv)
    client = _client()
    try:
        if args.cmd == "workers":
            return cmd_workers(client)
        if args.cmd == "ready":
            return cmd_ready(client, args.limit)
        if args.cmd == "task":
            return cmd_task(client, args.aone_id)
        if args.cmd == "operation":
            return cmd_operation(client, args.operation_id)
        return cmd_discard_resume(
            client, args.task_id, args.session_id, args.reason, args.yes)
    except ControlPlaneError as e:
        sys.stderr.write("error: %s\n" % e)
        return 3


if __name__ == "__main__":
    sys.exit(main())
