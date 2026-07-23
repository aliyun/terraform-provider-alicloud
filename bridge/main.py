#!/usr/bin/env python3
"""Standalone composition root for Scheduler-owned Jobs.

This process deliberately contains no DingTalk stream, legacy scanner, PR/reply
watcher, or Task queue worker.  It is the sole runtime behind
``bridge/scheduler.sh`` while the legacy bridge remains in
``jarvis_dingtalk_bot.py``.
"""

from __future__ import annotations

import logging
import os
import signal
import sys
import threading
from typing import Any

try:  # Executed by bridge/scheduler.sh from the bridge directory.
    from jarvis_task_client import ControlPlaneClient
    from scheduler.registry import JOBS, REGISTRY
    from scheduler.runners import build_runners
    from scheduler.service import SchedulerService
except ModuleNotFoundError:  # Package import for tests and tools.
    from bridge.jarvis_task_client import ControlPlaneClient
    from bridge.scheduler.registry import JOBS, REGISTRY
    from bridge.scheduler.runners import build_runners
    from bridge.scheduler.service import SchedulerService


LOG = logging.getLogger("jarvis-scheduler")
def _task_client_from_env() -> ControlPlaneClient:
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
            "Scheduler requires JARVIS_CONTROL_PLANE_TOKEN or JARVIS_HTML_REPORT_TOKEN")
    return ControlPlaneClient(
        base_url, token,
        timeout=float(os.environ.get("JARVIS_CONTROL_PLANE_TIMEOUT", "10")),
    )


def build_scheduler() -> SchedulerService:
    return SchedulerService(
        task_client=_task_client_from_env(),
        runners=build_runners(logger=LOG),
        registry=REGISTRY,
        logger=LOG,
    )


def main() -> int:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s [%(threadName)s] %(message)s",
        stream=sys.stderr,
    )
    scheduler = build_scheduler()
    stop = threading.Event()

    def request_stop(signum: int, _frame: Any) -> None:
        LOG.info("Scheduler signal %s received; draining admitted jobs", signum)
        stop.set()

    signal.signal(signal.SIGTERM, request_stop)
    signal.signal(signal.SIGINT, request_stop)
    scheduler.start()
    LOG.info("Scheduler READY pid=%s worker=%s jobs=%s", os.getpid(),
             scheduler.worker_key, ",".join(definition.id for definition in JOBS))
    stop.wait()
    # A planned shutdown must not kill an admitted job midway through its
    # control-plane completion update.  Keep the process alive and retry the
    # bounded drain until it finishes; the supervisor therefore cannot start a
    # second Scheduler while the first still owns the job.
    while not scheduler.stop():
        LOG.error("Scheduler drain timed out; keeping this process alive and retrying")
        stop.wait(5)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
