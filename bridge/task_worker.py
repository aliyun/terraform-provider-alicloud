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
from typing import Any

try:  # Executed by bridge/run.sh from the bridge directory.
    import jarvis_dingtalk_bot as bot
    from jarvis_persistence_executor import PersistenceExecutor
except ModuleNotFoundError:  # Package import for tests and tools.
    from bridge import jarvis_dingtalk_bot as bot
    from bridge.jarvis_persistence_executor import PersistenceExecutor


LOG = logging.getLogger("jarvis-task-worker")


class TaskWorker:
    """Compose the durable executor around the bridge's task implementation."""

    def __init__(self, handler: Any, *, executor_factory: Any = PersistenceExecutor):
        self.handler = handler
        self.executor = executor_factory(
            handler.task_client,
            handler.capacity_manager,
            handler.persistent_task_execution.execute,
            lambda controller, reason: bot.stop_task_process(controller, reason, logger=LOG),
            worker_id_file=os.environ.get(
                "JARVIS_WORKER_ID_FILE",
                os.path.join(bot.jarvis_root(), ".my-day", "bridge", "worker-id")),
            capabilities={
                "kinds": sorted(handler.execution_router.task_types),
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
            progress=lambda _lease, controller: bot._session_progress_excerpt(
                controller.runtime_session_id),
            logger=LOG,
        )

    def start(self) -> "TaskWorker":
        self.executor.start()
        return self

    def stop(self, *, drain: bool = False, timeout: Any = None) -> bool:
        if not self.executor.stopped:
            if not self.executor.stop(drain=drain, timeout=timeout):
                return False
        # The handler is only a dependency container here, but its ephemeral
        # pool can own a subprocess created by a persistent Task callback.
        self.handler.ephemeral_executor.terminate_all(release_fn=bot._release_claim)
        self.handler.ephemeral_executor.shutdown(wait=False, cancel_futures=True)
        return True


def build_task_worker() -> TaskWorker:
    """Build a no-stream handler solely as the Task execution implementation."""
    return TaskWorker(bot.JarvisHandler(no_dingtalk=True))


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
