"""Bridge composition root for the fixed, fenced scheduled-job worker.

This module does not import legacy Bridge loops.  The caller supplies only the
job runners it has explicitly migrated, which makes a mixed legacy/new rollout
safe by construction.  Scheduler activation is fail-closed on Worker
acknowledgement or any job without a runner mapping.
"""

from __future__ import annotations

from datetime import datetime, timezone
import hashlib
import logging
import os
import socket
import threading
import uuid
from typing import Any, Mapping, Optional

try:  # bridge/run.sh executes from bridge/, package users import bridge.scheduler
    from jarvis_persistence_executor import _default_boot_id
    from jarvis_task_client import ControlPlaneClient, StaleFence
except ModuleNotFoundError:  # pragma: no cover - import path depends on composition root
    from bridge.jarvis_persistence_executor import _default_boot_id
    from bridge.jarvis_task_client import ControlPlaneClient, StaleFence

from .control_plane_client import HttpScheduledJobControlPlane
from .engine import DurableResultPublisher, SchedulerEngine
from .jobs import JOBS, load_jobs
from .migration import (
    SchedulerMigrationError, business_job_enabled,
)
from .model import JobResult, ScheduledJobDefinition
from .registry import SchedulerRegistry
from .runtime import JobRunner, ScannerRuntime


SCHEDULER_WORKER_KEY = "bridge-scheduler"


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
                 runners: Mapping[str, JobRunner], registry: SchedulerRegistry,
                 environ: Optional[Mapping[str, str]] = None,
                 logger: Optional[logging.Logger] = None) -> None:
        self._task_client = task_client
        # ``runners`` is an implementation catalogue keyed by runner name;
        # ownership is defined only by the checked-in YAML registry.
        self._runners = dict(runners)
        self._registry = registry
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
            raise SchedulerCompositionError("Scheduler host identity is empty")
        self.boot_id = _scheduler_boot_id(self._environ, self.host_id)
        # Parse operational knobs only after the explicit enable gate.  A
        # dormant scheduler must not change legacy Bridge startup behavior.
        self._heartbeat_interval = 30.0
        self._poll_interval = 5.0

    @property
    def ready(self) -> bool:
        with self._lock:
            return self._ready

    @property
    def running(self) -> bool:
        """Whether this process still owns a live Scheduler Engine."""

        with self._lock:
            return self._engine is not None

    @property
    def enabled(self) -> bool:
        return bool(self._registry.scheduler_job_keys())

    def start(self) -> bool:
        """Register/verify the fixed Worker, register jobs, recover, then poll.

        ``False`` means the checked-in registry has no Scheduler-owned jobs.
        Any configured-but-unmapped runner raises before a legacy loop can be
        suppressed, so the operator cannot lose a job during cutover.
        """

        definitions = load_jobs()
        migrated = self._registry.scheduler_job_keys()
        if not migrated:
            return False
        if "daily.probe" in migrated:
            probe = next(definition for definition in definitions
                         if definition.id == "daily.probe")
            configured_hour = self._environ.get("JARVIS_PROBE_HOUR")
            registry_hour = str(getattr(probe.schedule, "hour", ""))
            if configured_hour is not None and configured_hour.strip() != registry_hour:
                raise SchedulerCompositionError(
                    "migrated daily.probe uses registry hour %s; "
                    "remove JARVIS_PROBE_HOUR or set it to %s"
                    % (registry_hour, registry_hour))
        self._heartbeat_interval = _positive_float(
            self._environ.get("JARVIS_SCHEDULER_HEARTBEAT_SEC", "30"),
            "JARVIS_SCHEDULER_HEARTBEAT_SEC")
        self._poll_interval = _positive_float(
            self._environ.get("JARVIS_SCHEDULER_POLL_SEC", "5"),
            "JARVIS_SCHEDULER_POLL_SEC")
        configured_runners = {
            job_key: self._registry.runner_for(job_key)
            for job_key in migrated
        }
        missing = sorted(
            job_key for job_key, runner_name in configured_runners.items()
            if runner_name not in self._runners)
        if missing:
            raise SchedulerCompositionError(
                "Scheduler-owned jobs have no Bridge runner mapping: %s"
                % ", ".join(missing))
        used_runners = frozenset(configured_runners.values())
        unused = sorted(set(self._runners).difference(used_runners))
        if unused:
            raise SchedulerCompositionError(
                "Bridge runner mappings are not declared by scheduler registry: %s"
                % ", ".join(unused))
        job_runners = {
            job_key: self._runners[runner_name]
            for job_key, runner_name in configured_runners.items()
        }
        with self._lock:
            if self._thread is not None:
                return True
            self._stop.clear()
            self._heartbeat_stop.clear()
            worker_client = self._worker_client()
            try:
                with self._worker_rpc_lock:
                    self._worker_status = "ACTIVE"
                    self._register_worker(worker_client, "ACTIVE")
                control_plane = HttpScheduledJobControlPlane(
                    self._task_client.base_url,
                    self._task_client.token,
                    timeout=self._task_client.timeout,
                    worker_key=self.worker_key,
                    process_uuid=self.process_uuid,
                )
                engine = SchedulerEngine(
                    definitions,
                    control_plane=control_plane,
                    runtime=ScannerRuntime(_RoutedRunner(job_runners)),
                    publisher=EmptyResultPublisher(),
                )
                # Register the complete immutable registry; unmigrated jobs are
                # explicitly DISABLED on the new control plane and remain legacy-owned.
                engine.register(
                    datetime.now(timezone.utc),
                    is_enabled=lambda definition: (
                        definition.id in migrated
                        and business_job_enabled(definition, environ=self._environ)
                    ))
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
            self.worker_key, self.host_id, self.process_uuid, ",".join(sorted(migrated)))
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
        if registered:
            try:
                with self._worker_rpc_lock:
                    self._worker_status = "DRAINING"
                    self._register_worker(self._worker_client(), "DRAINING")
            except Exception as exc:  # preserve fail-closed ownership; do not lie OFFLINE
                self._log.warning("Scheduler Worker DRAINING heartbeat failed: %s", type(exc).__name__)
        if thread is not None:
            thread.join(timeout=_stop_timeout(timeout, self._environ))
            if thread.is_alive():
                self._log.warning(
                    "Scheduler Engine did not drain before timeout; "
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
                    self._register_worker(self._worker_client(), "OFFLINE")
            except Exception as exc:
                self._log.warning("Scheduler Worker OFFLINE update failed: %s", type(exc).__name__)
                offline = False
        with self._lock:
            self._thread = None
            self._heartbeat_thread = None
            self._engine = None
            self._registered = False
        return offline

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
                    self._register_worker(self._worker_client(), "ACTIVE")
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
                        self._worker_client(), self._worker_status)
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

    def _worker_client(self) -> ControlPlaneClient:
        return self._task_client

    def _worker_payload(self, status: str) -> dict[str, Any]:
        return {
            "workerKey": self.worker_key,
            "hostId": self.host_id,
            # ``host`` keeps the existing Worker endpoint compatible while C5
            # adds the canonical hostId check server-side.
            "host": self.host_id,
            "bootId": self.boot_id,
            "processUuid": self.process_uuid,
            "status": status,
            "maxSlots": 1,
            "activeSessions": 0,
            "freeSlots": 1,
            "capabilities": {
                "role": "scheduler",
                "dispatch": {"pull": False},
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
            _require_active_worker(
                response, self.worker_key, self.host_id, self.boot_id, self.process_uuid)
            self._registered = True


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


def _scheduler_boot_id(environ: Mapping[str, str], host_id: str) -> str:
    configured = str(environ.get("JARVIS_BOOT_ID", "")).strip()
    return configured or _default_boot_id(host_id)


def _require_active_worker(response: Any, worker_key: str, host_id: str,
                           boot_id: str, process_uuid: str) -> None:
    if not isinstance(response, Mapping):
        raise SchedulerCompositionError("Scheduler Worker register response must be an object")
    worker = response.get("worker", response)
    if not isinstance(worker, Mapping):
        raise SchedulerCompositionError("Scheduler Worker register response.worker must be an object")
    expected = {
        "workerKey": worker_key,
        "hostId": host_id,
        "bootId": boot_id,
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
    dispatch = capabilities.get("dispatch")
    if not isinstance(dispatch, Mapping) or dispatch.get("pull") is not False:
        raise SchedulerCompositionError(
            "Scheduler Worker register acknowledgement lacks capabilities.dispatch.pull=false")


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
        value = environ.get(
            "JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS",
            environ.get(
                "JARVIS_BRIDGE_STOP_WAIT",
                environ.get(
                    "JARVIS_STOP_GRACE",
                    environ.get("JARVIS_SCHEDULER_DRAIN_TIMEOUT", "600"),
                ),
            ),
        )
    return max(0.0, float(value))


__all__ = [
    "EmptyResultPublisher", "SCHEDULER_WORKER_KEY",
    "SchedulerComposition", "SchedulerCompositionError", "SchedulerMigrationError",
]
