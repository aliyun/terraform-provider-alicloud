#!/usr/bin/env python3
"""Best-effort control-plane OFFLINE for bridge workers that were force-killed.

A ``SIGKILL``'d worker cannot deregister itself, so its stable ``workerKey``
fence lingers on the control plane (~9 min until the reaper reconciles it),
blocking a replacement with ``409 WORKER_ALREADY_ACTIVE``. To close that
window, each worker persists its control-plane identity
(``{workerKey, processUuid, hostId, bootId}``) on ACTIVE registration
(:func:`write_identity`) and removes it on a clean OFFLINE
(:func:`clear_identity`). Whoever force-kills the worker — the ``bridge.main``
supervisor after its shutdown grace, or ``run.sh`` on its force path — then
invokes :func:`offline_all`, which reads the surviving identity files and posts
an OFFLINE heartbeat on the dead process's behalf.

The OFFLINE heartbeat carries the *persisted* ``processUuid`` so the control
plane accepts it as the fence owner. If a newer process already took over the
same ``workerKey`` with a fresh ``processUuid``, the server rejects the stale
OFFLINE — so a leftover identity file can never knock a live replacement
offline.

Kept independent of :mod:`bridge.scheduler` (and PyYAML) so worker-only hosts,
which run only the Persistent Worker, can call it too.
"""

from __future__ import annotations

import argparse
import json
import logging
import os
import sys
from pathlib import Path
from typing import Any, Mapping, Optional

LOG = logging.getLogger("bridge.worker_offline")

CP_IDENTITY_SUFFIX = ".worker-cp.json"


def state_dir(environ: Optional[Mapping[str, str]] = None) -> Path:
    """Resolve the bridge state directory that holds worker identity files."""
    env = os.environ if environ is None else environ
    configured = str(env.get("JARVIS_BRIDGE_STATE_DIR", "")).strip()
    if configured:
        return Path(configured)
    repo_root = Path(__file__).resolve().parent.parent
    return repo_root / ".my-day" / "bridge"


def _identity_path(name: str, environ: Optional[Mapping[str, str]] = None,
                   *, directory: Optional[Path] = None) -> Path:
    base = directory if directory is not None else state_dir(environ)
    return base / ("%s%s" % (name, CP_IDENTITY_SUFFIX))


def write_identity(name: str, *, worker_key: str, process_uuid: str,
                   host_id: str, boot_id: str,
                   environ: Optional[Mapping[str, str]] = None,
                   directory: Optional[Path] = None) -> None:
    """Atomically persist a worker's control-plane identity for later OFFLINE.

    Called right after ACTIVE registration succeeds. Best-effort: a failure to
    persist must never take down the worker, so exceptions are swallowed with a
    warning (the reaper remains the ultimate fallback).
    """
    path = _identity_path(name, environ, directory=directory)
    record = {
        "workerKey": str(worker_key),
        "processUuid": str(process_uuid),
        "hostId": str(host_id),
        "bootId": str(boot_id),
    }
    try:
        path.parent.mkdir(parents=True, exist_ok=True)
        temporary = path.with_name("%s.tmp.%s" % (path.name, os.getpid()))
        temporary.write_text(
            json.dumps(record, ensure_ascii=False), encoding="utf-8")
        os.replace(temporary, path)
    except Exception as exc:  # noqa: BLE001 - persistence is best-effort
        LOG.warning("worker-offline: could not persist identity %s: %s",
                    name, type(exc).__name__)


def clear_identity(name: str, environ: Optional[Mapping[str, str]] = None,
                   *, directory: Optional[Path] = None) -> None:
    """Remove a worker's identity file after a clean OFFLINE (idempotent)."""
    try:
        _identity_path(name, environ, directory=directory).unlink(
            missing_ok=True)
    except Exception as exc:  # noqa: BLE001 - never block a clean shutdown
        LOG.warning("worker-offline: could not clear identity %s: %s",
                    name, type(exc).__name__)


def _client_from_env(environ: Mapping[str, str], timeout: float):
    """Build a ControlPlaneClient from env, or return None when unconfigured.

    Mirrors ``bridge.scheduler.scheduler._task_client_from_env`` but stays
    independent of the Scheduler module (and PyYAML) and never raises: a
    missing token just means "cannot offline remotely", which is best-effort.
    """
    from bridge.jarvis_task_client import ControlPlaneClient

    base_url = (
        str(environ.get("JARVIS_CONTROL_PLANE_BASE_URL", "")).strip()
        or str(environ.get("JARVIS_HTML_REPORT_BASE_URL", "")).strip()
        or "https://pre-agent.aliyun-inc.com"
    )
    token = (
        str(environ.get("JARVIS_CONTROL_PLANE_TOKEN", "")).strip()
        or str(environ.get("JARVIS_HTML_REPORT_TOKEN", "")).strip()
    )
    if not token:
        LOG.warning("worker-offline: no control-plane token; skipping remote OFFLINE")
        return None
    return ControlPlaneClient(base_url, token, timeout=timeout)


def _offline_one(client: Any, record: Mapping[str, Any]) -> bool:
    """Post one OFFLINE heartbeat using the persisted owner processUuid."""
    worker_key = str(record.get("workerKey") or "").strip()
    process_uuid = str(record.get("processUuid") or "").strip()
    if not worker_key or not process_uuid:
        LOG.warning("worker-offline: skipping identity without workerKey/processUuid")
        return False
    host_id = str(record.get("hostId") or "").strip()
    payload = {
        "workerKey": worker_key,
        "hostId": host_id,
        "host": host_id,
        "bootId": str(record.get("bootId") or "").strip(),
        "status": "OFFLINE",
    }
    client.heartbeat_worker(
        worker_key, payload, process_uuid=process_uuid,
        request_id="jarvis-worker-offline-%s" % process_uuid[:16])
    return True


def offline_all(*, environ: Optional[Mapping[str, str]] = None,
                client: Any = None, timeout: float = 10.0,
                directory: Optional[Path] = None) -> int:
    """Offline every worker with a surviving identity file. Best-effort.

    Reads each ``*.worker-cp.json`` under the state dir, posts an OFFLINE
    heartbeat with its persisted processUuid, and removes the file on success.
    Never raises; returns the count of workers successfully taken OFFLINE.
    """
    env = os.environ if environ is None else environ
    base = directory if directory is not None else state_dir(env)
    try:
        files = sorted(base.glob("*%s" % CP_IDENTITY_SUFFIX))
    except Exception as exc:  # noqa: BLE001 - missing dir etc.
        LOG.warning("worker-offline: cannot list %s: %s", base, type(exc).__name__)
        return 0
    if not files:
        return 0
    if client is None:
        client = _client_from_env(env, timeout)
        if client is None:
            return 0
    offlined = 0
    for path in files:
        try:
            record = json.loads(path.read_text(encoding="utf-8"))
        except Exception as exc:  # noqa: BLE001 - corrupt/partial file
            LOG.warning("worker-offline: unreadable identity %s: %s",
                        path.name, type(exc).__name__)
            continue
        try:
            if _offline_one(client, record):
                offlined += 1
                path.unlink(missing_ok=True)
                LOG.info("worker-offline: released fence for %s",
                         record.get("workerKey"))
        except Exception as exc:  # noqa: BLE001 - control-plane reject/network
            # Leave the file in place; the reaper remains the fallback and a
            # later attempt (or a takeover by a fresh process) can still win.
            LOG.warning("worker-offline: OFFLINE failed for %s: %s",
                        record.get("workerKey"), type(exc).__name__)
    return offlined


def main(argv: Optional[list] = None) -> int:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s [%(threadName)s] %(message)s",
        stream=sys.stderr,
    )
    parser = argparse.ArgumentParser(
        description="Best-effort control-plane OFFLINE for force-killed bridge workers")
    parser.add_argument(
        "--all", action="store_true",
        help="offline every worker with a surviving identity file")
    parser.add_argument(
        "--timeout", type=float, default=10.0,
        help="per-request control-plane timeout seconds (default 10)")
    args = parser.parse_args(argv)
    # --all is the only mode today; keep the flag explicit for forward clarity.
    count = offline_all(timeout=args.timeout)
    LOG.info("worker-offline: %d worker(s) taken OFFLINE", count)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
