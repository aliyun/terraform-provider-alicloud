#!/usr/bin/env python3
"""Pure SchedulerService process entrypoint.

The bridge process supervisor lives in :mod:`bridge.main`.  This module owns
only Scheduler validation, construction, signal handling, and bounded drain.
"""

from __future__ import annotations

import logging
import os
import signal
import sys
import threading
import time
from typing import Any, Mapping

from bridge.jarvis_task_client import (
    ControlPlaneClient, DEFAULT_CONTROL_PLANE_BASE_URL,
)
from bridge.scheduler.jobs import JOBS, runner_key
from bridge.scheduler.runners import build_runners
from bridge.scheduler.service import SCHEDULER_WORKER_KEY, SchedulerService


LOG = logging.getLogger("jarvis-scheduler")


def require_scheduler_role(environ: Mapping[str, str] | None = None) -> None:
    env = os.environ if environ is None else environ
    role = env.get("JARVIS_BRIDGE_ROLE", "scheduler").strip()
    if role != "scheduler":
        raise RuntimeError(
            "Scheduler requires JARVIS_BRIDGE_ROLE=scheduler; "
            "worker hosts run only the Persistent Worker")
def _task_client_from_env(
        environ: Mapping[str, str] | None = None) -> ControlPlaneClient:
    env = os.environ if environ is None else environ
    base_url = (
        env.get("JARVIS_CONTROL_PLANE_BASE_URL", "").strip()
        or DEFAULT_CONTROL_PLANE_BASE_URL
    )
    token = (
        env.get("JARVIS_CONTROL_PLANE_TOKEN", "").strip()
        or env.get("JARVIS_HTML_REPORT_TOKEN", "").strip()
    )
    if not token:
        raise RuntimeError(
            "Scheduler requires JARVIS_CONTROL_PLANE_TOKEN or "
            "JARVIS_HTML_REPORT_TOKEN")
    return ControlPlaneClient(
        base_url,
        token,
        timeout=float(env.get("JARVIS_CONTROL_PLANE_TIMEOUT", "10")),
    )


def build_scheduler(
        environ: Mapping[str, str] | None = None) -> SchedulerService:
    """Validate and build the Scheduler without starting control-plane work."""

    env = os.environ if environ is None else environ
    require_scheduler_role(env)
    task_client = _task_client_from_env(env)
    runners = build_runners(
        logger=LOG,
        task_client=task_client,
        worker_key=SCHEDULER_WORKER_KEY,
        repo_root=os.path.dirname(os.path.dirname(os.path.dirname(__file__))),
    )
    declared = {runner_key(definition) for definition in JOBS}
    missing = declared.difference(runners)
    unused = set(runners).difference(declared)
    if missing or unused:
        raise RuntimeError(
            "Scheduler runner registry mismatch: missing=%s unused=%s"
            % (
                ",".join(sorted(map(str, missing))) or "-",
                ",".join(sorted(map(str, unused))) or "-",
            ))
    return SchedulerService(
        task_client=task_client,
        runners=runners,
        definitions=JOBS,
        environ=env,
        logger=LOG,
    )


def validate(environ: Mapping[str, str] | None = None) -> None:
    """Fail closed on role, token, PyYAML/job registry, or runner drift."""

    build_scheduler(environ)


def main() -> int:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s [%(threadName)s] %(message)s",
        stream=sys.stderr,
    )
    if sys.argv[1:] == ["--validate"]:
        validate()
        LOG.info("Scheduler validation OK")
        return 0
    if sys.argv[1:]:
        raise SystemExit("usage: python -m bridge.scheduler.scheduler [--validate]")

    scheduler = build_scheduler()
    stop = threading.Event()

    def request_stop(signum: int, _frame: Any) -> None:
        LOG.info("Scheduler signal %s received; draining admitted jobs", signum)
        stop.set()

    signal.signal(signal.SIGTERM, request_stop)
    signal.signal(signal.SIGINT, request_stop)
    scheduler.start()
    LOG.info(
        "Scheduler READY pid=%s worker=%s jobs=%s",
        os.getpid(),
        scheduler.worker_key,
        ",".join(definition.id for definition in JOBS),
    )
    stop.wait()
    while not scheduler.stop():
        LOG.error(
            "Scheduler drain incomplete; remaining DRAINING until supervisor "
            "force-stop or a retry completes")
        time.sleep(5)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
