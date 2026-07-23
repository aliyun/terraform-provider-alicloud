#!/usr/bin/env python3
"""Independent composition root for durable Task execution.

The worker deliberately runs outside the DingTalk/Scheduler process.  A
scheduler restart can therefore replace its stream listener and periodic-job
runtime without fencing or terminating an already leased Task Session.
"""

from __future__ import annotations

import logging
import os
import signal
import sys
import threading
from pathlib import Path
from typing import Any

try:  # Executed by bridge/run.sh from the bridge directory.
    import jarvis_dingtalk_bot as bot
    from execution.persistent_task import PersistentTaskExecution, stop_task_process
    from jarvis_capacity import CapacityManager
    from jarvis_execution_runtime import DEFAULT_EXECUTION_RUNTIME
    from jarvis_field_repair import FIELD_REPAIR_KIND, FieldRepairWorker
    from jarvis_persistence_executor import PersistenceExecutor
    from jarvis_task_client import ControlPlaneClient
    from jarvis_task_router import ExecutionRouter
except ModuleNotFoundError:  # Package import for tests and tools.
    from bridge import jarvis_dingtalk_bot as bot
    from bridge.execution.persistent_task import (
        PersistentTaskExecution,
        stop_task_process,
    )
    from bridge.jarvis_capacity import CapacityManager
    from bridge.jarvis_execution_runtime import DEFAULT_EXECUTION_RUNTIME
    from bridge.jarvis_field_repair import FIELD_REPAIR_KIND, FieldRepairWorker
    from bridge.jarvis_persistence_executor import PersistenceExecutor
    from bridge.jarvis_task_client import ControlPlaneClient
    from bridge.jarvis_task_router import ExecutionRouter


LOG = logging.getLogger("jarvis-task-worker")
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
            "Task worker requires JARVIS_CONTROL_PLANE_TOKEN or "
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


class TaskExecutionRuntime:
    """Minimal dependency container for leased Task execution.

    The headless dispatch state machine is still shared with the legacy bridge
    through unbound methods and pure callbacks, but this composition root
    deliberately does not construct or retain ``JarvisHandler``.  In particular
    it owns no DingTalk streaming state and no Aone/Daily/PR Scheduler objects.
    """

    # These methods contain the established headless execution/bookend behavior.
    # Reuse them as unbound implementation methods until that state machine is
    # fully extracted; none of them requires JarvisHandler.__init__.
    dispatch_item = bot.JarvisHandler.dispatch_item
    _maybe_suspend = bot.JarvisHandler._maybe_suspend
    _completion_broadcast = bot.JarvisHandler._completion_broadcast
    _dispatch_failed = bot.JarvisHandler._dispatch_failed
    _post_death_cause = staticmethod(bot.JarvisHandler._post_death_cause)
    _workitem_line = staticmethod(bot.JarvisHandler._workitem_line)
    _last_comment_id = staticmethod(bot.JarvisHandler._last_comment_id)

    def __init__(
        self,
        *,
        task_client: Any = None,
        execution_router_factory: Any = ExecutionRouter,
        capacity_manager_factory: Any = CapacityManager,
        field_repair_worker_factory: Any = FieldRepairWorker,
        ephemeral_executor_factory: Any = bot.EphemeralExecutor,
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
            claude_bin=bot.claude_bin(),
        )
        self.ephemeral_executor = ephemeral_executor_factory(
            max_workers=max_slots,
            capacity_manager=self.capacity_manager,
            execution_runtime=self.execution_runtime,
        )
        self.persistent_task_execution = PersistentTaskExecution(
            enabled_kinds=lambda: self.execution_router.task_types,
            dispatch_item=self.dispatch_item,
            task_bookend=lambda *args, **kwargs: bot._TaskAoneBookend(
                *args, **kwargs),
            terraform_rd_ready=bot._terraform_rd_ready,
            routine_notice=_routine_notice,
            quick_card=_quick_card,
            field_repair_worker=self.field_repair_worker,
            field_repair_kind=FIELD_REPAIR_KIND,
            task_bookend_kinds=bot.TASK_BOOKEND_KINDS,
            post_pr_headless_kinds=bot.POST_PR_HEADLESS_KINDS,
            broadcast_target=bot.broadcast_target,
            broadcast_type=bot.broadcast_type,
        )


class TaskWorker:
    """Compose the durable executor around the bridge's task implementation."""

    def __init__(
        self,
        runtime: Any,
        *,
        executor_factory: Any = PersistenceExecutor,
        stop_process: Any = stop_task_process,
        release_claim: Any = bot._release_claim,
        progress: Any = bot._session_progress_excerpt,
    ):
        self.runtime = runtime
        self.executor = executor_factory(
            runtime.task_client,
            runtime.capacity_manager,
            runtime.persistent_task_execution.execute,
            lambda controller, reason: stop_process(controller, reason, logger=LOG),
            worker_id_file=os.environ.get(
                "JARVIS_WORKER_ID_FILE",
                os.path.join(bot.jarvis_root(), ".my-day", "bridge", "worker-id")),
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
                controller.runtime_session_id),
            logger=LOG,
        )
        self._release_claim = release_claim

    def start(self) -> "TaskWorker":
        self.executor.start()
        return self

    def stop(self, *, drain: bool = False, timeout: Any = None) -> bool:
        if not self.executor.stopped:
            if not self.executor.stop(drain=drain, timeout=timeout):
                return False
        # The runtime's ephemeral pool can own a subprocess created by a
        # persistent Task callback.
        self.runtime.ephemeral_executor.terminate_all(
            release_fn=self._release_claim)
        self.runtime.ephemeral_executor.shutdown(
            wait=False, cancel_futures=True)
        return True


def build_task_worker(*, executor_factory: Any = PersistenceExecutor) -> TaskWorker:
    """Build the durable worker without constructing a Bot or Scheduler."""
    return TaskWorker(TaskExecutionRuntime(), executor_factory=executor_factory)


def main() -> int:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s [%(threadName)s] %(message)s",
        stream=sys.stderr,
    )
    worker = build_task_worker().start()
    stop = threading.Event()

    def request_stop(signum: int, _frame: Any) -> None:
        LOG.info("Task worker signal %s received; stopping leased sessions", signum)
        stop.set()

    signal.signal(signal.SIGTERM, request_stop)
    signal.signal(signal.SIGINT, request_stop)
    LOG.info("Task worker READY pid=%s role=%s", os.getpid(),
             os.environ.get("JARVIS_BRIDGE_ROLE", "scheduler"))
    stop.wait()
    worker.stop(
        drain=False,
        timeout=float(os.environ.get("JARVIS_WORKER_DRAIN_TIMEOUT", "30")),
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
