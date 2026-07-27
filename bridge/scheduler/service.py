"""Lifecycle service for the fixed, fenced scheduled-job worker.

This module does not import legacy Bridge loops.  The caller supplies only the
job runners it has explicitly migrated, which makes a mixed legacy/new rollout
safe by construction.  Scheduler activation is fail-closed on Worker
acknowledgement or any job without a runner mapping.
"""

from __future__ import annotations

from datetime import datetime, timezone
import logging
import os
import socket
import threading
import time
import uuid
from typing import Any, Mapping, Optional, Sequence

from bridge.jarvis_persistence_executor import _default_boot_id
from bridge.jarvis_task_client import ControlPlaneClient, StaleFence

from .control_plane_client import HttpScheduledJobControlPlane
from .engine import (
    JobRunner,
    RunnerDispatcher,
    SchedulerEngine,
)
from .jobs import JOBS, runner_key
from .model import ScheduledJobDefinition


SCHEDULER_WORKER_KEY = "bridge-scheduler"


class SchedulerServiceError(RuntimeError):
    """A Scheduler startup/shutdown precondition was not safely satisfied."""


class SchedulerService:
    """Own the Scheduler Worker lifecycle and the single Engine polling thread."""

    def __init__(self, *, task_client: ControlPlaneClient,
                 runners: Mapping[object, JobRunner],
                 definitions: Sequence[ScheduledJobDefinition] = JOBS,
                 environ: Optional[Mapping[str, str]] = None,
                 logger: Optional[logging.Logger] = None) -> None:
        self._task_client = task_client
        self._runners = dict(runners)
        self._definitions = tuple(definitions)
        self._environ = os.environ if environ is None else environ
        self._log = logger or logging.getLogger(__name__)
        self._engine: Optional[SchedulerEngine] = None
        self._thread: Optional[threading.Thread] = None
        self._heartbeat_thread: Optional[threading.Thread] = None
        self._stop = threading.Event()
        self._heartbeat_stop = threading.Event()
        self._ready = False
        self._registered = False
        self._worker_status = "OFFLINE"
        self._lock = threading.RLock()
        self._worker_rpc_lock = threading.Lock()
        self.process_uuid = uuid.uuid4().hex
        self.worker_key = SCHEDULER_WORKER_KEY
        self.host_id = str(socket.gethostname() or socket.getfqdn()).strip()
        if not self.host_id:
            raise SchedulerServiceError("Scheduler host identity is empty")
        self.boot_id = _scheduler_boot_id(self._environ, self.host_id)
        # Parse operational knobs only after the explicit enable gate.  A
        # dormant scheduler must not change legacy Bridge startup behavior.
        self._heartbeat_interval = 30.0
        self._poll_interval = 5.0
        self._max_concurrency = 4

    @property
    def ready(self) -> bool:
        with self._lock:
            return self._ready

    @property
    def running(self) -> bool:
        """Whether this process still owns a live Scheduler Engine."""

        with self._lock:
            return self._engine is not None

    def start(self) -> bool:
        """Register/verify the fixed Worker, register jobs, recover, then poll.

        ``False`` means the checked-in registry has no Scheduler-owned jobs.
        Any configured-but-unmapped runner raises before a legacy loop can be
        suppressed, so the operator cannot lose a job during cutover.
        """

        definitions = self._definitions
        if not definitions:
            return False
        self._heartbeat_interval = _positive_float(
            self._environ.get("JARVIS_SCHEDULER_HEARTBEAT_SEC", "30"),
            "JARVIS_SCHEDULER_HEARTBEAT_SEC")
        self._poll_interval = _positive_float(
            self._environ.get("JARVIS_SCHEDULER_POLL_SEC", "5"),
            "JARVIS_SCHEDULER_POLL_SEC")
        self._max_concurrency = _positive_int(
            self._environ.get("JARVIS_SCHEDULER_MAX_CONCURRENCY", "4"),
            "JARVIS_SCHEDULER_MAX_CONCURRENCY")
        configured_runners = frozenset(runner_key(item) for item in definitions)
        missing = sorted((
            runner_name for runner_name in configured_runners
            if runner_name not in self._runners), key=_runner_label)
        if missing:
            raise SchedulerServiceError(
                "Scheduler-owned jobs have no registered runner: %s"
                % ", ".join(_runner_label(item) for item in missing))
        unused = sorted(
            set(self._runners).difference(configured_runners),
            key=_runner_label,
        )
        if unused:
            raise SchedulerServiceError(
                "registered runners are not declared by scheduler registry: %s"
                % ", ".join(_runner_label(item) for item in unused))
        with self._lock:
            if self._thread is not None:
                return True
            self._stop.clear()
            self._heartbeat_stop.clear()
            worker_client = self._task_client
            try:
                with self._worker_rpc_lock:
                    self._worker_status = "ACTIVE"
                    self._register_worker(worker_client, "ACTIVE")
                control_plane = HttpScheduledJobControlPlane(
                    self._task_client.base_url,
                    self._task_client.token,
                    timeout=self._task_client.timeout,
                    retry_delay_seconds=float(self._environ.get(
                        "JARVIS_SCHEDULER_CONTROL_RETRY_SEC", "5")),
                    worker_key=self.worker_key,
                    process_uuid=self.process_uuid,
                )
                engine = SchedulerEngine(
                    definitions,
                    control_plane=control_plane,
                    runtime=RunnerDispatcher(self._runners),
                    max_workers=self._max_concurrency,
                )
                engine.register(datetime.now(timezone.utc))
                engine.recover_interrupted()
            except Exception:
                # A Worker that could register but could not complete the
                # startup protocol must not remain ACTIVE and fence a retry.
                try:
                    with self._worker_rpc_lock:
                        self._worker_status = "OFFLINE"
                        self._register_worker(worker_client, "OFFLINE")
                except Exception as offline_error:
                    self._log.warning(
                        "Scheduler Worker OFFLINE cleanup after failed startup failed: %s",
                        type(offline_error).__name__)
                self._registered = False
                raise
            self._engine = engine
            self._ready = True
            self._heartbeat_thread = threading.Thread(
                target=self._heartbeat_loop,
                name="bridge-scheduler-heartbeat",
                daemon=True,
            )
            self._thread = threading.Thread(
                target=self._loop, name="bridge-scheduler-engine", daemon=True)
            self._heartbeat_thread.start()
            self._thread.start()
        self._log.info(
            "SchedulerEngine READY worker=%s hostId=%s processUuid=%s jobs=%s",
            self.worker_key, self.host_id, self.process_uuid,
            ",".join(sorted(definition.id for definition in definitions)))
        return True

    def stop(self, *, timeout: Optional[float] = None) -> bool:
        """Quiesce admissions and mark OFFLINE, or cancel a timed-out restart."""

        with self._lock:
            engine = self._engine
            thread = self._thread
            heartbeat_thread = self._heartbeat_thread
            registered = self._registered
            if engine is None:
                return True
            engine.stop()
            self._ready = False
            self._stop.set()
        drain_timeout = _stop_timeout(timeout, self._environ)
        deadline = time.monotonic() + drain_timeout
        if registered:
            try:
                with self._worker_rpc_lock:
                    self._worker_status = "DRAINING"
                    self._register_worker(self._task_client, "DRAINING")
            except Exception as exc:  # preserve fail-closed ownership; do not lie OFFLINE
                self._log.warning("Scheduler Worker DRAINING heartbeat failed: %s", type(exc).__name__)
        if thread is not None:
            thread.join(timeout=max(0.0, deadline - time.monotonic()))
            if thread.is_alive():
                self._log.warning(
                    "Scheduler Engine did not drain before timeout; "
                    "cancelling the planned restart")
                self._cancel_drain(engine, thread, registered)
                return False
        if not engine.wait_for_active(
                timeout=max(0.0, deadline - time.monotonic())):
            self._log.warning(
                "Scheduler admitted jobs did not drain before timeout; "
                "cancelling the planned restart")
            self._cancel_drain(engine, thread, registered)
            return False
        self._heartbeat_stop.set()
        if heartbeat_thread is not None:
            heartbeat_thread.join(timeout=self._task_client.timeout + 1)
            if heartbeat_thread.is_alive():
                self._log.warning(
                    "Scheduler heartbeat did not stop; leaving Worker non-OFFLINE")
                return False
        offline = True
        if registered:
            try:
                with self._worker_rpc_lock:
                    self._worker_status = "OFFLINE"
                    self._register_worker(self._task_client, "OFFLINE")
            except Exception as exc:
                self._log.warning("Scheduler Worker OFFLINE update failed: %s", type(exc).__name__)
                offline = False
        if not offline:
            # Keep the drained Engine/registration state so the caller can
            # retry the OFFLINE transition without losing the Worker fence.
            return False
        with self._lock:
            engine.shutdown()
            self._thread = None
            self._heartbeat_thread = None
            self._engine = None
            self._registered = False
        return True

    def _cancel_drain(
        self,
        engine: SchedulerEngine,
        thread: threading.Thread,
        registered: bool,
    ) -> None:
        """Restore ACTIVE scheduling after the restart drain budget expires."""

        if registered:
            try:
                with self._worker_rpc_lock:
                    self._worker_status = "ACTIVE"
                    self._register_worker(self._task_client, "ACTIVE")
            except Exception as exc:
                # Keep ACTIVE as the desired heartbeat state. Until the
                # control plane acknowledges it, new starts remain fail-closed;
                # the independent heartbeat loop will retry the transition.
                self._log.warning(
                    "Scheduler restart cancelled; ACTIVE restore is pending: %s",
                    type(exc).__name__)
        engine.resume()
        self._stop.clear()
        # If the old loop observed the stop event immediately before it was
        # cleared, give it one scheduling turn to exit before deciding whether
        # a replacement loop is needed.
        thread.join(timeout=0.05)
        with self._lock:
            self._ready = True
            if not thread.is_alive():
                thread = threading.Thread(
                    target=self._loop,
                    name="bridge-scheduler-engine",
                    daemon=True,
                )
                self._thread = thread
                thread.start()
        self._log.info("Scheduler planned restart cancelled; ACTIVE scheduling resumed")

    def _loop(self) -> None:
        while not self._stop.is_set():
            try:
                engine = self._engine
                if engine is not None:
                    engine.tick(datetime.now(timezone.utc))
            except Exception as exc:  # network/protocol faults close execution for this tick
                self._log.exception("Scheduler Engine tick failed closed: %s", type(exc).__name__)
            self._stop.wait(self._poll_interval)

    def _heartbeat_loop(self) -> None:
        """Keep the Scheduler fence alive independently of long-running jobs."""

        while not self._heartbeat_stop.wait(self._heartbeat_interval):
            try:
                with self._worker_rpc_lock:
                    if self._worker_status not in ("ACTIVE", "DRAINING"):
                        return
                    self._register_worker(
                        self._task_client, self._worker_status)
            except StaleFence:
                self._log.error(
                    "Scheduler heartbeat lost the current process fence; "
                    "stopping new admissions")
                with self._lock:
                    engine = self._engine
                    self._ready = False
                    self._stop.set()
                if engine is not None:
                    engine.stop()
                return
            except Exception as exc:
                self._log.warning(
                    "Scheduler Worker %s heartbeat failed: %s",
                    self._worker_status, type(exc).__name__)

    def _worker_payload(self, status: str) -> dict[str, Any]:
        with self._lock:
            active = self._engine.active_count if self._engine is not None else 0
        return {
            "workerKey": self.worker_key,
            "hostId": self.host_id,
            # ``host`` keeps the existing Worker endpoint compatible while C5
            # adds the canonical hostId check server-side.
            "host": self.host_id,
            "bootId": self.boot_id,
            "processUuid": self.process_uuid,
            "status": status,
            "maxSlots": self._max_concurrency,
            "activeSessions": active,
            "freeSlots": max(0, self._max_concurrency - active),
            "capabilities": {
                "role": "scheduler",
                "dispatch": {"pull": False},
                "scheduledJobApi": "v1",
                "jobKeys": [definition.id for definition in self._definitions],
            },
        }

    def _register_worker(self, client: ControlPlaneClient, status: str) -> None:
        request_id = "jarvis-scheduler-worker-%s-%s" % (
            status.lower(), uuid.uuid4().hex[:16])
        if status == "ACTIVE" and not self._registered:
            response = client.register_worker(self._worker_payload(status), request_id=request_id)
        else:
            response = client.heartbeat_worker(
                self.worker_key, self._worker_payload(status),
                process_uuid=self.process_uuid, request_id=request_id)
        if status == "ACTIVE":
            _require_active_worker(
                response, self.worker_key, self.host_id, self.boot_id, self.process_uuid)
            self._registered = True


def _scheduler_boot_id(environ: Mapping[str, str], host_id: str) -> str:
    configured = str(environ.get("JARVIS_BOOT_ID", "")).strip()
    return configured or _default_boot_id(host_id)


def _require_active_worker(response: Any, worker_key: str, host_id: str,
                           boot_id: str, process_uuid: str) -> None:
    if not isinstance(response, Mapping):
        raise SchedulerServiceError("Scheduler Worker register response must be an object")
    worker = response.get("worker", response)
    if not isinstance(worker, Mapping):
        raise SchedulerServiceError("Scheduler Worker register response.worker must be an object")
    expected = {
        "workerKey": worker_key,
        "hostId": host_id,
        "bootId": boot_id,
        "processUuid": process_uuid,
        "status": "ACTIVE",
    }
    mismatched = [name for name, value in expected.items() if worker.get(name) != value]
    if mismatched:
        raise SchedulerServiceError(
            "Scheduler Worker register acknowledgement rejected/mismatched fields: %s"
            % ", ".join(mismatched))
    capabilities = worker.get("capabilities")
    if not isinstance(capabilities, Mapping) or capabilities.get("role") != "scheduler":
        raise SchedulerServiceError(
            "Scheduler Worker register acknowledgement lacks capabilities.role=scheduler")
    dispatch = capabilities.get("dispatch")
    if not isinstance(dispatch, Mapping) or dispatch.get("pull") is not False:
        raise SchedulerServiceError(
            "Scheduler Worker register acknowledgement lacks capabilities.dispatch.pull=false")


def _positive_float(value: Any, name: str) -> float:
    try:
        parsed = float(value)
    except (TypeError, ValueError) as exc:
        raise SchedulerServiceError("%s must be a positive number" % name) from exc
    if parsed <= 0:
        raise SchedulerServiceError("%s must be a positive number" % name)
    return parsed


def _positive_int(value: Any, name: str) -> int:
    try:
        parsed = int(value)
    except (TypeError, ValueError) as exc:
        raise SchedulerServiceError("%s must be a positive integer" % name) from exc
    if parsed <= 0 or str(parsed) != str(value).strip():
        raise SchedulerServiceError("%s must be a positive integer" % name)
    return parsed


def _runner_label(value: object) -> str:
    if isinstance(value, tuple) and len(value) == 2:
        return "%s/%s" % value
    return str(value)


def _stop_timeout(value: Optional[float], environ: Mapping[str, str]) -> float:
    if value is None:
        value = environ.get(
            "JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS",
            environ.get(
                "JARVIS_BRIDGE_STOP_WAIT",
                environ.get(
                    "JARVIS_STOP_GRACE",
                    environ.get("JARVIS_SCHEDULER_DRAIN_TIMEOUT", "30"),
                ),
            ),
        )
    return max(0.0, float(value))


__all__ = [
    "SCHEDULER_WORKER_KEY", "SchedulerService", "SchedulerServiceError",
]
