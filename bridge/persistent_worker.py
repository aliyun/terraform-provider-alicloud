#!/usr/bin/env python3
"""Independent composition root for durable Task execution.

The worker deliberately runs outside the DingTalk/Scheduler process. An
independent Scheduler failure can therefore recover without terminating an
already leased Task Session. An explicit bridge restart stops this Worker too:
active Sessions are relinquished and then resumed by the replacement Worker.
"""

from __future__ import annotations

import logging
import os
import signal
import sys
import threading
import subprocess
import json
from pathlib import Path
from typing import Any

from bridge.persistent_tasks import (
    POST_PR_HEADLESS_KINDS,
    TASK_BOOKEND_KINDS,
    PersistentTaskExecution,
    TaskAoneBookend,
    _aone_event_sanitize_text,
    dispatch_item,
    extract_suspend,
    normalized_failure_subtype,
    release_claim,
    stop_task_process,
    terraform_rd_ready,
)
from bridge.task_policy import HEADLESS_POLICY_REVISION
from bridge.aone_tasks import master_staff
from bridge.helpers.dingtalk import _dingtalk_event_enqueue
from bridge.jarvis_capacity import CapacityManager
from bridge.jarvis_execution_runtime import (
    DEFAULT_EXECUTION_RUNTIME,
    EphemeralExecutor,
    claude_bin,
    jarvis_root,
    session_progress_excerpt,
)
from bridge.jarvis_field_repair import FieldRepairWorker
from bridge.jarvis_persistence_executor import PersistenceExecutor
from bridge.jarvis_task_client import ControlPlaneClient
from bridge.jarvis_task_router import ExecutionRouter
from bridge.process_group_runner import run_process_group


LOG = logging.getLogger("jarvis-persistent-worker")
REPO_ROOT = Path(__file__).resolve().parents[1]


def _task_client_from_env() -> ControlPlaneClient:
    """Build the worker's mandatory control-plane client."""
    base_url = (
        os.environ.get("JARVIS_CONTROL_PLANE_BASE_URL", "").strip()
        or os.environ.get("JARVIS_HTML_REPORT_BASE_URL", "").strip()
        or "https://pre-agent.aliyun-inc.com"
    )
    token = (
        os.environ.get("JARVIS_CONTROL_PLANE_TOKEN", "").strip()
        or os.environ.get("JARVIS_HTML_REPORT_TOKEN", "").strip()
    )
    if not token:
        raise RuntimeError(
            "Persistent worker requires JARVIS_CONTROL_PLANE_TOKEN or "
            "JARVIS_HTML_REPORT_TOKEN")
    return ControlPlaneClient(
        base_url,
        token,
        timeout=float(os.environ.get("JARVIS_CONTROL_PLANE_TIMEOUT", "10")),
    )


def _routine_notice(text: str) -> None:
    """Keep routine worker lifecycle notices local to the worker log."""
    LOG.info("[ROUTINE] %s", (text or "").replace("\n", " | ")[:1000])


def _quick_card(_target: str, text: str, _target_type: str = "user") -> None:
    """Worker-side notice sink; the worker never initializes DingTalk."""
    LOG.info("[BROADCAST] %s", (text or "").replace("\n", " | ")[:1000])


class PersistentTaskRuntime:
    """Minimal dependency container for leased Task execution.

    The headless dispatch state machine is still shared with the legacy bridge
    through unbound methods and pure callbacks, but this composition root
    deliberately does not construct or retain ``JarvisHandler``.  In particular
    it owns no DingTalk streaming state and no Aone/Daily/PR Scheduler objects.
    """

    def __init__(
        self,
        *,
        task_client: Any = None,
        execution_router_factory: Any = ExecutionRouter,
        capacity_manager_factory: Any = CapacityManager,
        field_repair_worker_factory: Any = FieldRepairWorker,
        ephemeral_executor_factory: Any = EphemeralExecutor,
    ) -> None:
        self.task_client = (
            task_client if task_client is not None else _task_client_from_env())
        self.execution_router = execution_router_factory(
            client=self.task_client, logger=LOG)
        max_slots = int(os.environ.get("JARVIS_DISPATCH_MAX", "3"))
        self.execution_runtime = DEFAULT_EXECUTION_RUNTIME
        self.capacity_manager = capacity_manager_factory(max_slots)
        self.field_repair_worker = field_repair_worker_factory(
            repo_root=REPO_ROOT,
            client=self.task_client,
            runtime=self.execution_runtime,
            claude_bin=claude_bin(),
        )
        self.ephemeral_executor = ephemeral_executor_factory(
            max_workers=max_slots,
            capacity_manager=self.capacity_manager,
            execution_runtime=self.execution_runtime,
        )
        self.persistent_task_execution = PersistentTaskExecution(
            enabled_kinds=lambda: self.execution_router.task_types,
            dispatch_item=self.dispatch_item,
            task_bookend=TaskAoneBookend,
            terraform_rd_ready=terraform_rd_ready,
            routine_notice=_routine_notice,
            quick_card=_quick_card,
            field_repair_worker=self.field_repair_worker,
            task_bookend_kinds=TASK_BOOKEND_KINDS,
            post_pr_headless_kinds=POST_PR_HEADLESS_KINDS,
            broadcast_target=lambda: os.environ.get(
                "JARVIS_BROADCAST_TARGET", ""),
            broadcast_type=lambda: os.environ.get(
                "JARVIS_BROADCAST_TYPE", "group"),
            policy_revision=HEADLESS_POLICY_REVISION,
        )

    def dispatch_item(self, *args, **kwargs):
        return dispatch_item(self, *args, **kwargs)

    @staticmethod
    def _last_comment_id(item_id):
        try:
            result = run_process_group(
                [str(REPO_ROOT / "bin" / "a1id"), "--", "project",
                 "workitem", "comment", "list", str(item_id), "-f", "json"],
                capture_output=True, text=True, timeout=30,
                cwd=str(REPO_ROOT))
            if result.returncode == 0:
                return max(
                    (comment.get("id", 0)
                     for comment in json.loads(result.stdout)),
                    default=0)
        except Exception:  # noqa: BLE001
            pass
        return 0

    def _maybe_suspend(self, final_text, _sid, _target, _target_type,
                       **_kwargs):
        info = extract_suspend(final_text or "")[1]
        if info:
            info = dict(info)
            info["wait_cursor"] = self._last_comment_id(info["aone_id"])
        return info

    @staticmethod
    def _workitem_line(item_id):
        item_id = str(item_id)
        if not item_id.isdigit():
            return "#%s" % item_id
        try:
            result = run_process_group(
                [str(REPO_ROOT / "bin" / "a1id"), "--", "project",
                 "workitem", "get", item_id, "-f", "json"],
                capture_output=True, text=True, timeout=30,
                cwd=str(REPO_ROOT))
            if result.returncode:
                return "#%s" % item_id
            item = json.loads(result.stdout)
            fields = {
                field.get("identifier"): field
                for field in item.get("fields", [])
                if isinstance(field, dict)
            }
            display = lambda key: (
                (fields.get(key) or {}).get("displayValue")
                or (fields.get(key) or {}).get("value") or "")
            project = (fields.get("space") or {}).get("value") or ""
            priority = display("priority")
            return (
                "- [#%s](https://project.aone.alibaba-inc.com/v2/project/%s/"
                "req/%s) %s%s" % (
                    item_id, project, item_id,
                    str(item.get("title") or "").strip(),
                    " [%s]" % priority if priority else ""),
                display("tag"),
            )
        except Exception:  # noqa: BLE001
            return "#%s" % item_id

    def _completion_broadcast(self, item_id):
        result = self._workitem_line(item_id)
        if not isinstance(result, tuple):
            return ("✅ 工单 #%s 处理完成（headless）" % item_id
                    if str(item_id).isdigit()
                    else "✅ 任务 #%s 处理完成" % item_id)
        line, tags = result
        prefix = next((
            text for tag, text in (
                ("jarvis-done", "✅ 工单处理完成"),
                ("jarvis-idle", "⏸️ 工单阶段完成·待人工接手"),
                ("jarvis-claimed", "⚠️ 工单处理结束但未收尾（仍 claimed）"),
            ) if tag in tags), "⚠️ 本轮未处理该工单（未获认领）")
        return prefix + "\n" + line

    @staticmethod
    def _post_death_cause(item_id, cause, terraform=False):
        """Record execution failure locally only; never write the Aone ticket."""
        if not str(item_id).isdigit():
            return
        prefix = "terraform " if terraform else ""
        LOG.warning("%sTask #%s death cause: %s", prefix, item_id,
                    str(cause).replace("\n", " | ")[:500])

    def _dispatch_failed(self, item_id, result, notify, project,
                         terraform=False, kind="ticket",
                         sid="unknown-session", attempts=None, **_kwargs):
        tail = (result.text or "").strip()
        subtype = normalized_failure_subtype(
            tail, getattr(result, "subtype", ""),
            bool(getattr(result, "is_error", True)))
        cause = "headless 派发失败\nsubtype: %s\n---\n%s" % (
            subtype, tail[-800:] or "(无输出)")
        self._post_death_cause(item_id, cause, terraform=terraform)
        release_state = "无需释放"
        if (project and str(item_id).isdigit()
                and kind not in POST_PR_HEADLESS_KINDS):
            release_state = "已释放"
            try:
                release_claim(item_id, project, terraform=terraform)
            except Exception as exc:  # noqa: BLE001
                release_state = "释放失败"
                LOG.warning("Task #%s claim release failed: %s", item_id, exc)
        if (terraform and subtype == "model_provider_error"
                and project and str(item_id).isdigit()):
            event_key = "dispatch-model-provider:%s:%s" % (
                str(kind or "ticket")[:64],
                str(sid or "unknown-session")[:96],
            )
            count = max(1, int(attempts or 1))
            text = (
                "Jarvis Terraform 自动处理因模型提供方故障停止。\n\n"
                "- 工单：#%s\n"
                "- 任务类型：%s\n"
                "- 失败原因：model_provider_error（模型提供方或网关请求失败）\n"
                "- 尝试次数：%s\n"
                "- 认领释放：%s\n"
                "- 下一步：请检查模型服务恢复情况，恢复后重新派发。"
                % (item_id, kind, count, release_state)
            )
            if not _dingtalk_event_enqueue(
                    item_id, project, event_key, master_staff(),
                    "Jarvis 模型提供方故障", text):
                LOG.error(
                    "Task #%s model-provider event could not be queued",
                    item_id)
        notify("⚠️ #%s 后台处理失败（%s）" % (item_id, subtype))


class PersistentWorker:
    """Compose the durable executor around the bridge's task implementation."""

    def __init__(
        self,
        runtime: Any,
        *,
        executor_factory: Any = PersistenceExecutor,
        stop_process: Any = stop_task_process,
        release_claim: Any = release_claim,
        progress: Any = session_progress_excerpt,
    ):
        self.runtime = runtime
        beacon_prefix = Path(os.environ.get(
            "JARVIS_EXECUTOR_BEACON_PREFIX",
            os.path.join(
                jarvis_root(), ".my-day", "bridge",
                "heartbeat.persistent-worker"),
        ))
        self.executor = executor_factory(
            runtime.task_client,
            runtime.capacity_manager,
            runtime.persistent_task_execution.execute,
            lambda controller, reason: stop_process(controller, reason, logger=LOG),
            worker_id_file=os.environ.get(
                "JARVIS_WORKER_ID_FILE",
                os.path.join(jarvis_root(), ".my-day", "bridge", "worker-id")),
            capabilities={
                "kinds": sorted(runtime.execution_router.task_types),
                "bridgeRole": os.environ.get("JARVIS_BRIDGE_ROLE", "scheduler"),
                "workerMode": "PERSISTENT",
                "client": "bridge",
            },
            lease_seconds=int(os.environ.get("JARVIS_LEASE_SECONDS", "660")),
            lease_safety_margin=float(
                os.environ.get("JARVIS_LEASE_SAFETY_MARGIN_SEC", "60")),
            lease_interval=float(os.environ.get("JARVIS_LEASE_POLL_SEC", "2")),
            worker_heartbeat_interval=float(
                os.environ.get("JARVIS_WORKER_HEARTBEAT_SEC", "30")),
            session_heartbeat_interval=float(
                os.environ.get("JARVIS_SESSION_HEARTBEAT_SEC", "30")),
            retry_interval=float(os.environ.get("JARVIS_CONTROL_PLANE_RETRY_SEC", "5")),
            progress=lambda _lease, controller: progress(
                controller.runtime_session_id,
                sanitizer=_aone_event_sanitize_text),
            heartbeat_beacon_paths={
                name: Path("%s.%s.epoch" % (beacon_prefix, name))
                for name in ("worker", "lease", "session")
            },
            logger=LOG,
        )
        self._release_claim = release_claim

    def start(self) -> "PersistentWorker":
        self.executor.start()
        return self

    def request_stop(self) -> None:
        """Stop new leases immediately; terminal handoff follows in ``stop``."""
        self.executor.begin_shutdown()

    def stop(self, *, drain: bool = False, timeout: Any = None) -> bool:
        self.request_stop()
        handoff_complete = True
        if not self.executor.stopped:
            if not self.executor.stop(drain=drain, timeout=timeout):
                handoff_complete = False
                LOG.warning(
                    "Persistent Worker stopped locally with incomplete control-plane "
                    "handoff; lease expiry/recovery will reconcile remaining Sessions")
        # The runtime's ephemeral pool can own a subprocess created by a
        # persistent Task callback. Always terminate it even when a control-plane
        # relinquish failed so restart cannot leave an unmanaged local process.
        self.runtime.ephemeral_executor.terminate_all(
            release_fn=self._release_claim)
        self.runtime.ephemeral_executor.shutdown(
            wait=False, cancel_futures=True)
        return handoff_complete


def build_persistent_worker(
        *, executor_factory: Any = PersistenceExecutor) -> PersistentWorker:
    """Build the durable worker without constructing a Bot or Scheduler."""
    return PersistentWorker(
        PersistentTaskRuntime(), executor_factory=executor_factory)


def main() -> int:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s [%(threadName)s] %(message)s",
        stream=sys.stderr,
    )
    worker = build_persistent_worker().start()
    stop = threading.Event()

    def request_stop(signum: int, _frame: Any) -> None:
        LOG.info("Persistent worker signal %s received; stopping leased sessions", signum)
        worker.request_stop()
        stop.set()

    signal.signal(signal.SIGTERM, request_stop)
    signal.signal(signal.SIGINT, request_stop)
    LOG.info("Persistent worker READY pid=%s role=%s", os.getpid(),
             os.environ.get("JARVIS_BRIDGE_ROLE", "scheduler"))
    stop.wait()
    stopped_cleanly = worker.stop(
        drain=False,
        timeout=float(os.environ.get("JARVIS_WORKER_DRAIN_TIMEOUT", "30")),
    )
    if not stopped_cleanly:
        LOG.error(
            "Persistent Worker exit has incomplete Session handoff; "
            "replacement must rely on lease expiry/recovery")
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
