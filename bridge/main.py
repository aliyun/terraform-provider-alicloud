#!/usr/bin/env python3
"""Standalone composition root for Scheduler-owned Jobs.

This process contains no DingTalk stream or Persistent Worker. ``bridge/run.sh``
supervises it alongside those independent entry points on the scheduler host.
"""

from __future__ import annotations

import logging
import os
import signal
import sys
import threading
import time
from typing import Any, Mapping

from bridge.jarvis_task_client import ControlPlaneClient
from bridge.scheduler.jobs import JOBS
from bridge.scheduler.runners import build_runners
from bridge.scheduler.service import SCHEDULER_WORKER_KEY, SchedulerService


LOG = logging.getLogger("jarvis-scheduler")


def require_scheduler_role(environ: Mapping[str, str] | None = None) -> None:
    role = (os.environ if environ is None else environ).get(
        "JARVIS_BRIDGE_ROLE", "scheduler").strip()
    if role != "scheduler":
        raise RuntimeError(
            "bridge/run.sh requires JARVIS_BRIDGE_ROLE=scheduler for periodic jobs; "
            "worker hosts run bridge/run.sh only")


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
    require_scheduler_role()
    task_client = _task_client_from_env()
    return SchedulerService(
        task_client=task_client,
        runners=build_runners(
            logger=LOG,
            task_client=task_client,
            worker_key=SCHEDULER_WORKER_KEY,
            repo_root=os.path.dirname(os.path.dirname(__file__)),
        ),
        definitions=JOBS,
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
        time.sleep(5)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
