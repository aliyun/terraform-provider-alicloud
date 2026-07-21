"""Bridge composition root for the fixed, fenced scheduled-job worker.

This module does not import legacy Bridge loops.  The caller supplies only the
job runners it has explicitly migrated, which makes a mixed legacy/new rollout
safe by construction.  Scheduler activation is fail-closed on the dedicated
credential, host identity, Worker acknowledgement, or any job without a
runner mapping.
"""

from __future__ import annotations

from datetime import datetime, timezone
import hashlib
import logging
import os
import socket
import threading
import time
import uuid
from typing import Any, Mapping, Optional

try:  # bridge/run.sh executes from bridge/, package users import bridge.scheduler
    from jarvis_task_client import ControlPlaneClient
except ModuleNotFoundError:  # pragma: no cover - import path depends on composition root
    from bridge.jarvis_task_client import ControlPlaneClient

from .control_plane_client import HttpScheduledJobControlPlane
from .engine import DurableResultPublisher, SchedulerEngine
from .jobs import JOBS, load_jobs
from .migration import SchedulerMigrationError, requested_new_jobs, scheduler_enabled
from .model import JobResult, ScheduledJobDefinition
from .runtime import JobRunner, ScannerRuntime


SCHEDULER_WORKER_KEY = "bridge-scheduler"
SCHEDULER_HOST_ID = "AgenticTools-Macmini.local"


class SchedulerCompositionError(RuntimeError):
    """A Scheduler startup/shutdown precondition was not safely satisfied."""


class EmptyResultPublisher(DurableResultPublisher):
    """Accept only jobs whose runner has no separate durable observation payload.

    The first compatible job, ``daily.probe``, performs its bounded enqueue in
    its runner and returns no observations.  A later Task/Event producing
    migration must replace this publisher with its durable acknowledgement
    adapter rather than silently dropping observations.
    """

    def publish(self, definition: ScheduledJobDefinition, result: JobResult,
                scheduled_for: datetime) -> None:
        del scheduled_for
        if result.observations:
            raise SchedulerCompositionError(
                "%s returned observations but has no durable publisher" % definition.id)


class SchedulerComposition:
    """Own the Scheduler Worker lifecycle and the single Engine polling thread."""

    def __init__(self, *, task_client: ControlPlaneClient,
                 runners: Mapping[str, JobRunner],
                 environ: Optional[Mapping[str, str]] = None,
                 logger: Optional[logging.Logger] = None) -> None:
        self._task_client = task_client
        self._runners = dict(runners)
        self._environ = os.environ if environ is None else environ
        self._log = logger or logging.getLogger(__name__)
        self._engine: Optional[SchedulerEngine] = None
        self._thread: Optional[threading.Thread] = None
        self._stop = threading.Event()
        self._ready = False
        self._registered = False
        self._lock = threading.RLock()
        self.process_uuid = uuid.uuid4().hex
        self.worker_key = SCHEDULER_WORKER_KEY
        self.host_id = SCHEDULER_HOST_ID
        # Parse operational knobs only after the explicit enable gate.  A
        # dormant scheduler must not change legacy Bridge startup behavior.
        self._heartbeat_interval = 30.0
        self._poll_interval = 5.0
        self._next_heartbeat = 0.0

    @property
    def ready(self) -> bool:
        with self._lock:
            return self._ready

    @property
    def enabled(self) -> bool:
        return scheduler_enabled(self._environ)

    def start(self) -> bool:
        """Register/verify the fixed Worker, register jobs, recover, then poll.

        ``False`` means the scheduler global gate is disabled or no job is
        selected.  Any selected-but-unmapped job raises before a legacy loop
        can be suppressed, so the operator cannot lose a job during cutover.
        """

        if not self.enabled:
            return False
        self._heartbeat_interval = _positive_float(
            self._environ.get("JARVIS_SCHEDULER_HEARTBEAT_SEC", "30"),
            "JARVIS_SCHEDULER_HEARTBEAT_SEC")
        self._poll_interval = _positive_float(
            self._environ.get("JARVIS_SCHEDULER_POLL_SEC", "5"),
            "JARVIS_SCHEDULER_POLL_SEC")
        definitions = load_jobs()
        migrated = requested_new_jobs(
            (definition.id for definition in definitions), environ=self._environ)
        if not migrated:
            self._log.info("SchedulerEngine disabled: no job has route=new")
            return False
        missing = sorted(migrated.difference(self._runners))
        if missing:
            raise SchedulerCompositionError(
                "selected new Scheduler jobs have no Bridge runner mapping: %s"
                % ", ".join(missing))
        self._require_scheduler_host()
        with self._lock:
            if self._thread is not None:
                return True
            self._stop.clear()
            worker_client = self._worker_client()
            self._register_worker(worker_client, "ACTIVE")
            try:
                control_plane = HttpScheduledJobControlPlane(
                    self._task_client.base_url,
                    _scheduler_token(self._environ),
                    timeout=self._task_client.timeout,
                    worker_key=self.worker_key,
                    process_uuid=self.process_uuid,
                )
                engine = SchedulerEngine(
                    definitions,
                    control_plane=control_plane,
                    runtime=ScannerRuntime(_RoutedRunner(self._runners)),
                    publisher=EmptyResultPublisher(),
                )
                # Register the complete immutable registry; unmigrated jobs are
                # explicitly DISABLED on the new control plane and remain legacy-owned.
                engine.register(
                    datetime.now(timezone.utc),
                    is_enabled=lambda definition: definition.id in migrated)
                engine.recover_interrupted()
            except Exception:
                # A Worker that could register but could not complete the
                # startup protocol must not remain ACTIVE and fence a retry.
                try:
                    self._register_worker(worker_client, "OFFLINE")
                except Exception as offline_error:
                    self._log.warning(
                        "Scheduler Worker OFFLINE cleanup after failed startup failed: %s",
                        type(offline_error).__name__)
                self._registered = False
                raise
            self._engine = engine
            self._ready = True
            self._thread = threading.Thread(
                target=self._loop, name="bridge-scheduler-engine", daemon=True)
            self._thread.start()
        self._log.info(
            "SchedulerEngine READY worker=%s hostId=%s processUuid=%s jobs=%s",
            self.worker_key, self.host_id, self.process_uuid, ",".join(sorted(migrated)))
        return True

    def stop(self, *, timeout: Optional[float] = None) -> bool:
        """Quiesce new admissions, drain the Engine loop, then mark OFFLINE."""

        with self._lock:
            engine = self._engine
            thread = self._thread
            registered = self._registered
            if engine is None:
                return True
            engine.stop()
            self._ready = False
            self._stop.set()
        if registered:
            try:
                self._register_worker(self._worker_client(), "DRAINING")
            except Exception as exc:  # preserve fail-closed ownership; do not lie OFFLINE
                self._log.warning("Scheduler Worker DRAINING heartbeat failed: %s", type(exc).__name__)
        if thread is not None:
            thread.join(timeout=_stop_timeout(timeout, self._environ))
            if thread.is_alive():
                self._log.warning("Scheduler Engine did not drain before timeout; leaving Worker non-OFFLINE")
                return False
        offline = True
        if registered:
            try:
                self._register_worker(self._worker_client(), "OFFLINE")
            except Exception as exc:
                self._log.warning("Scheduler Worker OFFLINE update failed: %s", type(exc).__name__)
                offline = False
        with self._lock:
            self._thread = None
            self._engine = None
            self._registered = False
        return offline

    def _loop(self) -> None:
        while not self._stop.is_set():
            try:
                self._heartbeat_if_due()
                engine = self._engine
                if engine is not None:
                    engine.tick(datetime.now(timezone.utc))
            except Exception as exc:  # network/protocol faults close execution for this tick
                self._log.exception("Scheduler Engine tick failed closed: %s", type(exc).__name__)
            self._stop.wait(self._poll_interval)

    def _require_scheduler_host(self) -> None:
        observed = {socket.gethostname().strip(), socket.getfqdn().strip()}
        if self.host_id not in observed:
            raise SchedulerCompositionError(
                "SchedulerEngine may run only on %s; observed hostnames=%s"
                % (self.host_id, ",".join(sorted(value for value in observed if value))))

    def _worker_client(self) -> ControlPlaneClient:
        return ControlPlaneClient(
            self._task_client.base_url,
            _scheduler_token(self._environ),
            timeout=self._task_client.timeout,
        )

    def _worker_payload(self, status: str) -> dict[str, Any]:
        return {
            "workerKey": self.worker_key,
            "hostId": self.host_id,
            # ``host`` keeps the existing Worker endpoint compatible while C5
            # adds the canonical hostId check server-side.
            "host": self.host_id,
            "processUuid": self.process_uuid,
            "status": status,
            "maxSlots": 1,
            "activeSessions": 0,
            "freeSlots": 1,
            "capabilities": {
                "role": "scheduler",
                "scheduledJobApi": "v1",
                "jobKeys": [definition.id for definition in JOBS],
            },
        }

    def _register_worker(self, client: ControlPlaneClient, status: str) -> None:
        request_id = "jarvis-scheduler-worker-%s-%s" % (
            status.lower(), hashlib.sha256(
                (self.process_uuid + uuid.uuid4().hex).encode("utf-8")).hexdigest()[:16])
        if status == "ACTIVE" and not self._registered:
            response = client.register_worker(self._worker_payload(status), request_id=request_id)
        else:
            response = client.heartbeat_worker(
                self.worker_key, self._worker_payload(status),
                process_uuid=self.process_uuid, request_id=request_id)
        if status == "ACTIVE":
            _require_active_worker(response, self.worker_key, self.host_id, self.process_uuid)
            self._registered = True
            self._next_heartbeat = time.monotonic() + self._heartbeat_interval

    def _heartbeat_if_due(self) -> None:
        if not self._registered or time.monotonic() < self._next_heartbeat:
            return
        self._register_worker(self._worker_client(), "ACTIVE")
        self._next_heartbeat = time.monotonic() + self._heartbeat_interval


class _RoutedRunner(JobRunner):
    def __init__(self, runners: Mapping[str, JobRunner]) -> None:
        self._runners = dict(runners)

    def run(self, definition: ScheduledJobDefinition, scheduled_for: datetime) -> JobResult:
        try:
            runner = self._runners[definition.id]
        except KeyError as exc:
            raise SchedulerCompositionError(
                "no runner mapping for admitted job %s" % definition.id) from exc
        return runner.run(definition, scheduled_for)


def _scheduler_token(environ: Mapping[str, str]) -> str:
    token = str(environ.get("JARVIS_SCHEDULER_CONTROL_PLANE_TOKEN", "")).strip()
    if not token:
        raise SchedulerCompositionError(
            "JARVIS_SCHEDULER_CONTROL_PLANE_TOKEN is required for SchedulerEngine; "
            "do not reuse the general Task Worker credential")
    return token


def _require_active_worker(response: Any, worker_key: str, host_id: str,
                           process_uuid: str) -> None:
    if not isinstance(response, Mapping):
        raise SchedulerCompositionError("Scheduler Worker register response must be an object")
    worker = response.get("worker", response)
    if not isinstance(worker, Mapping):
        raise SchedulerCompositionError("Scheduler Worker register response.worker must be an object")
    expected = {
        "workerKey": worker_key,
        "hostId": host_id,
        "processUuid": process_uuid,
        "status": "ACTIVE",
    }
    mismatched = [name for name, value in expected.items() if worker.get(name) != value]
    if mismatched:
        raise SchedulerCompositionError(
            "Scheduler Worker register acknowledgement rejected/mismatched fields: %s"
            % ", ".join(mismatched))
    capabilities = worker.get("capabilities")
    if not isinstance(capabilities, Mapping) or capabilities.get("role") != "scheduler":
        raise SchedulerCompositionError(
            "Scheduler Worker register acknowledgement lacks capabilities.role=scheduler")


def _positive_float(value: Any, name: str) -> float:
    try:
        parsed = float(value)
    except (TypeError, ValueError) as exc:
        raise SchedulerCompositionError("%s must be a positive number" % name) from exc
    if parsed <= 0:
        raise SchedulerCompositionError("%s must be a positive number" % name)
    return parsed


def _stop_timeout(value: Optional[float], environ: Mapping[str, str]) -> float:
    if value is None:
        value = environ.get("JARVIS_SCHEDULER_DRAIN_TIMEOUT", "30")
    return max(0.0, float(value))


__all__ = [
    "EmptyResultPublisher", "SCHEDULER_HOST_ID", "SCHEDULER_WORKER_KEY",
    "SchedulerComposition", "SchedulerCompositionError", "SchedulerMigrationError",
]
