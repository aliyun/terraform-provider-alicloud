#!/usr/bin/env python3
"""Single-host Task executor for the persistent Jarvis control plane.

The module intentionally knows nothing about the DingTalk bridge. Integration
code injects one shared ``CapacityManager``, an
``execute(lease, session_controller)`` callback, and a
``stop_process(session_controller, reason)`` hook. The two polling loops are
also exposed as ``run_once`` and ``heartbeat_sessions_once`` so behavior can be
tested without clocks, sockets, or real processes.
"""

from __future__ import annotations

import hashlib
import logging
import os
import socket
import threading
import time
import uuid
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Callable, Dict, Mapping, Optional, Tuple

from jarvis_capacity import CapacityManager, CapacityPermit
from jarvis_task_client import (
    ControlPlaneClient,
    ControlPlaneConflict,
    ControlPlaneError,
    ControlPlaneUnavailable,
    StaleFence,
)


class LeaseProtocolError(ValueError):
    """A successful lease response did not contain the required shape."""


def _nonblank(value: Any, name: str) -> str:
    text = str(value or "").strip()
    if not text:
        raise ValueError("%s must not be empty" % name)
    return text


def _component(value: Any, name: str) -> str:
    # Keep the documented host:boot:process shape unambiguous in logs.  The
    # HTTP client still URL-quotes the complete key when placing it in a path.
    return _nonblank(value, name).replace(":", "_").replace("/", "_")


def _default_boot_id(host: str) -> str:
    configured = os.environ.get("JARVIS_BOOT_ID", "").strip()
    if configured:
        return configured
    linux_boot_id = Path("/proc/sys/kernel/random/boot_id")
    try:
        value = linux_boot_id.read_text().strip()
        if value:
            return value
    except OSError:
        pass
    # Portable fallback for macOS: wall time minus monotonic time approximates
    # the boot epoch.  The process UUID remains the uniqueness boundary.
    boot_epoch = int(time.time() - time.monotonic())
    return uuid.uuid5(uuid.NAMESPACE_DNS, "%s:%s" % (host, boot_epoch)).hex


def make_worker_key(host: Optional[str] = None, boot_id: Optional[str] = None,
                    process_uuid: Optional[str] = None) -> str:
    """Build one stable key for this worker process: ``host:boot:process``."""
    host_value = _component(host or socket.gethostname(), "host")
    boot_value = _component(boot_id or _default_boot_id(host_value), "boot_id")
    process_value = _component(process_uuid or uuid.uuid4().hex, "process_uuid")
    return "%s:%s:%s" % (host_value, boot_value, process_value)


def _new_runtime_session_id() -> str:
    return str(uuid.uuid4())


def _mapping(value: Any, name: str) -> Mapping[str, Any]:
    if not isinstance(value, Mapping):
        raise LeaseProtocolError("%s must be an object" % name)
    return value


def _field(value: Mapping[str, Any], *names: str) -> Any:
    for name in names:
        if name in value:
            return value[name]
    return None


def parse_lease_response(response: Any) -> Optional[Dict[str, Any]]:
    """Return the canonical lease object or ``None`` for an empty poll."""
    response_map = _mapping(response, "lease response")
    if not response_map:
        return None
    # Frozen v1 contract returns {task, session, resumed}; accepting the early
    # {lease:{...}} draft costs nothing and eases mixed-version rollout.
    lease = response_map.get("lease") if "lease" in response_map else response_map
    if lease is None or lease == {}:
        return None
    lease_map = _mapping(lease, "lease")
    task = _mapping(lease_map.get("task"), "lease.task")
    session = _mapping(lease_map.get("session"), "lease.session")
    session_id = _field(session, "id", "sessionId", "session_id")
    fence_token = _field(session, "fenceToken", "fence_token")
    _nonblank(session_id, "lease.session.id")
    if fence_token is None or str(fence_token).strip() == "":
        raise LeaseProtocolError("lease.session.fenceToken must not be empty")
    canonical = dict(lease_map)
    canonical["task"] = dict(task)
    canonical["session"] = dict(session)
    return canonical


StopProcess = Callable[["SessionController", str], None]
NetworkFailure = Callable[[BaseException], None]


class _ThreadFuture:
    """Tiny Future subset used to avoid coupling this module to a process pool."""

    def __init__(self):
        self._done = threading.Event()
        self._exception: Optional[BaseException] = None

    def run(self, fn: Callable[..., Any], *args: Any) -> None:
        try:
            fn(*args)
        except BaseException as exc:  # the worker callback owns error reporting
            self._exception = exc
        finally:
            self._done.set()

    def done(self) -> bool:
        return self._done.is_set()


class _ThreadExecutor:
    """Immediate daemon-thread executor with the small API PersistenceExecutor needs."""

    def __init__(self):
        self._lock = threading.Lock()
        self._closed = False
        self._threads = []

    def submit(self, fn: Callable[..., Any], *args: Any) -> _ThreadFuture:
        with self._lock:
            if self._closed:
                raise RuntimeError("executor is closed")
            self._threads = [thread for thread in self._threads if thread.is_alive()]
            future = _ThreadFuture()
            thread = threading.Thread(
                target=future.run, args=(fn,) + args,
                name="jarvis-local-execute", daemon=True)
            self._threads.append(thread)
        thread.start()
        return future

    def shutdown(self, wait: bool = False, cancel_futures: bool = False) -> None:
        del cancel_futures  # running daemon threads are stopped through stop_process
        with self._lock:
            self._closed = True
            threads = list(self._threads)
        if wait:
            for thread in threads:
                thread.join()


class SessionController:
    """Fenced session transitions passed to the injected execution callback.

    A terminal transition that cannot reach AutomationAgent remains pending and
    can be retried with the same idempotency key.  A failed lease heartbeat is
    different: ownership is no longer provable, so the corresponding process is
    stopped and every terminal ACK is suppressed.
    """

    TERMINAL_ACTIONS = {"complete", "fail", "suspend"}

    def __init__(self, client: ControlPlaneClient, worker_key: str,
                 lease: Mapping[str, Any], *,
                 lease_seconds: int = 300,
                 runtime_session_id_factory: Callable[[], str] = _new_runtime_session_id,
                 stop_process: Optional[StopProcess] = None,
                 on_network_failure: Optional[NetworkFailure] = None,
                 logger: Optional[logging.Logger] = None):
        self.client = client
        self.worker_key = _nonblank(worker_key, "worker_key")
        self.lease = dict(_mapping(lease, "lease"))
        self.task = dict(_mapping(self.lease.get("task"), "lease.task"))
        self.session = dict(_mapping(self.lease.get("session"), "lease.session"))
        self.session_id = _nonblank(
            _field(self.session, "id", "sessionId", "session_id"), "session_id")
        self.fence_token = _field(self.session, "fenceToken", "fence_token")
        if self.fence_token is None or str(self.fence_token).strip() == "":
            raise LeaseProtocolError("fence_token must not be empty")
        self.lease_seconds = int(lease_seconds)
        if self.lease_seconds <= 0:
            raise ValueError("lease_seconds must be positive")
        existing_runtime_id = _field(
            self.session, "runtimeSessionId", "runtime_session_id")
        if existing_runtime_id is None or not str(existing_runtime_id).strip():
            existing_runtime_id = runtime_session_id_factory()
        self.runtime_session_id = _nonblank(
            existing_runtime_id, "runtime_session_id")
        self.session["runtimeSessionId"] = self.runtime_session_id
        self.lease["session"] = self.session
        self.resumed = bool(self.lease.get("resumed"))
        self.stop_process = stop_process
        self.on_network_failure = on_network_failure
        self.log = logger or logging.getLogger(__name__)
        self.process: Any = None

        self._lock = threading.RLock()
        # RLock lets a stop hook inspect/call controller helpers without
        # deadlocking the renewal failure path that invoked it.
        self._transition_lock = threading.RLock()
        self._started = False
        self._terminal_action: Optional[str] = None
        self._pending_terminal: Optional[Tuple[str, Dict[str, Any]]] = None
        self._ownership_lost = False
        self._stop_requested = False
        self._stop_hook_called = False
        self._stop_reason: Optional[str] = None

    @property
    def started(self) -> bool:
        with self._lock:
            return self._started

    @property
    def terminal(self) -> bool:
        with self._lock:
            return self._terminal_action is not None

    @property
    def terminal_action(self) -> Optional[str]:
        with self._lock:
            return self._terminal_action

    @property
    def pending_terminal(self) -> Optional[str]:
        with self._lock:
            return self._pending_terminal[0] if self._pending_terminal else None

    @property
    def ownership_lost(self) -> bool:
        with self._lock:
            return self._ownership_lost

    @property
    def stop_requested(self) -> bool:
        with self._lock:
            return self._stop_requested

    def bind_process(self, process: Any) -> Any:
        """Bind and durably publish the owned process-group leader.

        ``start()`` happens before the subprocess exists, so the first transition
        cannot contain a PID.  The spawn callback repeats the same fenced start API
        here to refresh ``jarvis_session.pid``.  Merely retaining the Popen in this
        process would make it undiscoverable after a bridge crash.

        A binding that cannot be committed is an ownership failure, not a best-effort
        observability loss: synchronously invoke ``stop_process`` while the concrete
        handle is available, then surface the original failure to the caller.
        """
        with self._transition_lock:
            with self._lock:
                self.process = process
                stop_after_bind = self._stop_requested or self._ownership_lost
                stop_reason = self._stop_reason or "process_bound_after_ownership_loss"
                if stop_after_bind:
                    # A fence may be lost in the narrow window after _execute_one's
                    # pre-launch check but before the subprocess callback binds its
                    # handle.  The earlier stop hook saw process=None, so permit one
                    # synchronous retry now that there is something concrete to kill.
                    self._stop_hook_called = False
            if stop_after_bind:
                self.request_stop(stop_reason)
                return process

            try:
                pid = int(getattr(process, "pid", 0))
                if pid <= 0:
                    raise ValueError("bound process must expose a positive pid")
                response = self.client.start_session(
                    self.session_id, self.worker_key, self.fence_token,
                    self._start_detail({"pid": pid}),
                    request_id=self._request_id("bind-process-%s" % pid))
            except ControlPlaneUnavailable as exc:
                self._notify_network_failure(exc)
                self._lose_ownership("control_plane_unavailable:bind_process")
                raise
            except StaleFence:
                self._lose_ownership("stale_fence:bind_process")
                raise
            except ControlPlaneError as exc:
                self._lose_ownership("process_bind_rejected")
                self.log.warning("process bind rejected session=%s code=%s",
                                 self.session_id, getattr(exc, "code", ""))
                raise
            except Exception:
                self._lose_ownership("process_bind_failed")
                raise

            returned_runtime_id = None
            if isinstance(response, Mapping):
                returned_runtime_id = _field(
                    response, "runtimeSessionId", "runtime_session_id")
                returned_session = response.get("session")
                if returned_runtime_id is None and isinstance(returned_session, Mapping):
                    returned_runtime_id = _field(
                        returned_session, "runtimeSessionId", "runtime_session_id")
            with self._lock:
                if returned_runtime_id is not None and str(returned_runtime_id).strip():
                    self.runtime_session_id = str(returned_runtime_id).strip()
                    self.session["runtimeSessionId"] = self.runtime_session_id
                self.session["pid"] = pid
                self.lease["session"] = self.session
                self._started = True
            return process

    def _request_id(self, action: str) -> str:
        raw = "%s|%s|%s|%s" % (
            self.worker_key, self.session_id, self.fence_token, action)
        digest = hashlib.sha256(raw.encode("utf-8")).hexdigest()[:32]
        return "jarvis-session-%s-%s" % (action, digest)

    @staticmethod
    def _detail(detail: Optional[Mapping[str, Any]]) -> Dict[str, Any]:
        if detail is None:
            return {}
        if not isinstance(detail, Mapping):
            raise TypeError("session detail must be a mapping")
        return dict(detail)

    def _lease_detail(self, detail: Optional[Mapping[str, Any]]) -> Dict[str, Any]:
        payload = self._detail(detail)
        payload["leaseSeconds"] = self.lease_seconds
        return payload

    def _start_detail(self, detail: Optional[Mapping[str, Any]]) -> Dict[str, Any]:
        payload = self._lease_detail(detail)
        payload["runtimeSessionId"] = self.runtime_session_id
        return payload

    def _notify_network_failure(self, exc: BaseException) -> None:
        if self.on_network_failure is not None:
            try:
                self.on_network_failure(exc)
            except Exception:
                self.log.exception("session network-failure callback crashed session=%s",
                                   self.session_id)

    def request_stop(self, reason: str) -> None:
        reason = _nonblank(reason, "reason")
        with self._lock:
            self._stop_requested = True
            if self._stop_reason is None:
                self._stop_reason = reason
            effective_reason = self._stop_reason
            if self._stop_hook_called:
                return
            self._stop_hook_called = True
        if self.stop_process is not None:
            try:
                self.stop_process(self, effective_reason)
            except Exception:
                self.log.exception("stop_process hook crashed session=%s",
                                   self.session_id)

    def _lose_ownership(self, reason: str) -> None:
        with self._lock:
            self._ownership_lost = True
            self._pending_terminal = None
        self.request_stop(reason)

    def start(self, detail: Optional[Mapping[str, Any]] = None) -> bool:
        with self._transition_lock:
            with self._lock:
                if self._started:
                    return True
                if self._ownership_lost:
                    return False
            try:
                response = self.client.start_session(
                    self.session_id, self.worker_key, self.fence_token,
                    self._start_detail(detail), request_id=self._request_id("start"))
            except ControlPlaneUnavailable as exc:
                self._notify_network_failure(exc)
                return False
            except StaleFence:
                self._lose_ownership("stale_fence:start")
                return False
            except ControlPlaneError as exc:
                self._lose_ownership("start_rejected")
                self.log.warning("session start rejected session=%s code=%s",
                                 self.session_id, getattr(exc, "code", ""))
                return False
            returned_runtime_id = None
            if isinstance(response, Mapping):
                returned_runtime_id = _field(
                    response, "runtimeSessionId", "runtime_session_id")
                returned_session = response.get("session")
                if returned_runtime_id is None and isinstance(returned_session, Mapping):
                    returned_runtime_id = _field(
                        returned_session, "runtimeSessionId", "runtime_session_id")
            with self._lock:
                if returned_runtime_id is not None and str(returned_runtime_id).strip():
                    self.runtime_session_id = str(returned_runtime_id).strip()
                    self.session["runtimeSessionId"] = self.runtime_session_id
                    self.lease["session"] = self.session
                self._started = True
            return True

    def heartbeat(self, detail: Optional[Mapping[str, Any]] = None) -> bool:
        with self._transition_lock:
            with self._lock:
                if self._terminal_action is not None:
                    return True
                if self._ownership_lost or self._pending_terminal is not None:
                    return False
            try:
                self.client.heartbeat_session(
                    self.session_id, self.worker_key, self.fence_token,
                    self._lease_detail(detail),
                    request_id="jarvis-session-heartbeat-%s" % uuid.uuid4().hex)
                return True
            except Exception as exc:
                if isinstance(exc, ControlPlaneUnavailable):
                    self._notify_network_failure(exc)
                self._lose_ownership(
                    "stale_fence:heartbeat" if isinstance(exc, StaleFence)
                    else "heartbeat_failed")
                self.log.warning("session heartbeat failed session=%s error=%s",
                                 self.session_id, type(exc).__name__)
                return False

    def _terminal(self, action: str, detail: Optional[Mapping[str, Any]], *,
                  retry: bool = False) -> bool:
        if action not in self.TERMINAL_ACTIONS:
            raise ValueError("invalid terminal action %r" % action)
        with self._transition_lock:
            with self._lock:
                if self._terminal_action is not None:
                    if self._terminal_action != action:
                        raise RuntimeError("session already terminal: %s" % self._terminal_action)
                    return True
                if self._ownership_lost:
                    return False
                if retry:
                    if self._pending_terminal is None:
                        return False
                    action, payload = self._pending_terminal
                else:
                    payload = self._detail(detail)
                    if self._pending_terminal is not None:
                        pending_action, _pending_payload = self._pending_terminal
                        if pending_action != action:
                            raise RuntimeError(
                                "session terminal transition already pending: %s" % pending_action)
                    self._pending_terminal = (action, payload)

            method = getattr(self.client, "%s_session" % action)
            try:
                method(self.session_id, self.worker_key, self.fence_token, payload,
                       request_id=self._request_id(action))
            except ControlPlaneUnavailable as exc:
                # Keep the desired terminal transition pending.  No terminal ACK
                # was committed; the worker pauses leasing and retries on recovery.
                self._notify_network_failure(exc)
                return False
            except StaleFence:
                self._lose_ownership("stale_fence:%s" % action)
                return False
            except ControlPlaneError as exc:
                self._lose_ownership("terminal_rejected:%s" % action)
                self.log.warning("terminal transition rejected session=%s action=%s code=%s",
                                 self.session_id, action, getattr(exc, "code", ""))
                return False

            with self._lock:
                self._pending_terminal = None
                self._terminal_action = action
            return True

    def retry_terminal(self) -> bool:
        with self._lock:
            pending = self._pending_terminal
        if pending is None:
            return self.terminal
        return self._terminal(pending[0], pending[1], retry=True)

    def complete(self, result: Any = None) -> bool:
        if isinstance(result, Mapping) and set(result) == {"result"}:
            payload = dict(result)
        else:
            payload = {"result": result}
        return self._terminal("complete", payload)

    def fail(self, error: Any = None, *, retry_after_seconds: Optional[int] = None) -> bool:
        if (isinstance(error, Mapping) and "error" in error and
                set(error).issubset({"error", "retryAfterSeconds"})):
            payload = dict(error)
        else:
            payload = {"error": error}
        if retry_after_seconds is not None:
            retry_after = int(retry_after_seconds)
            if retry_after < 0:
                raise ValueError("retry_after_seconds must not be negative")
            payload["retryAfterSeconds"] = retry_after
        return self._terminal("fail", payload)

    def suspend(self, detail: Optional[Mapping[str, Any]] = None) -> bool:
        return self._terminal("suspend", detail)


@dataclass
class _ActiveSession:
    controller: SessionController
    future: Any = None
    capacity_permit: Optional[CapacityPermit] = None


Execute = Callable[[Mapping[str, Any], SessionController], Any]


class PersistenceExecutor:
    """One-machine lease/execute/heartbeat loop with fail-closed networking."""

    def __init__(self, client: ControlPlaneClient, capacity_manager: CapacityManager,
                 execute: Execute, stop_process: StopProcess, *,
                 host: Optional[str] = None, boot_id: Optional[str] = None,
                 process_uuid: Optional[str] = None, worker_key: Optional[str] = None,
                 capabilities: Optional[Mapping[str, Any]] = None,
                 lease_seconds: int = 300,
                 runtime_session_id_factory: Callable[[], str] = _new_runtime_session_id,
                 lease_interval: float = 2.0, worker_heartbeat_interval: float = 30.0,
                 session_heartbeat_interval: float = 30.0, retry_interval: float = 5.0,
                 clock: Callable[[], float] = time.monotonic,
                 executor: Optional[Any] = None,
                 logger: Optional[logging.Logger] = None):
        if not callable(execute):
            raise TypeError("execute must be callable")
        if not callable(stop_process):
            raise TypeError("stop_process must be callable")
        if not isinstance(capacity_manager, CapacityManager):
            raise TypeError("capacity_manager must be a CapacityManager")
        self.capacity_manager = capacity_manager
        self.max_slots = capacity_manager.capacity
        if self.max_slots <= 0:
            raise ValueError("capacity_manager.capacity must be positive")
        for name, value in (
            ("lease_interval", lease_interval),
            ("worker_heartbeat_interval", worker_heartbeat_interval),
            ("session_heartbeat_interval", session_heartbeat_interval),
            ("retry_interval", retry_interval),
        ):
            if float(value) <= 0:
                raise ValueError("%s must be positive" % name)

        self.client = client
        self.execute = execute
        self.stop_process = stop_process
        self.worker_key = _nonblank(
            worker_key or make_worker_key(host, boot_id, process_uuid), "worker_key")
        parts = self.worker_key.split(":", 2)
        self.host = parts[0] if len(parts) == 3 else (host or socket.gethostname())
        self.boot_id = parts[1] if len(parts) == 3 else (boot_id or "unknown")
        self.process_uuid = parts[2] if len(parts) == 3 else (process_uuid or "unknown")
        self.capabilities = dict(capabilities or {})
        self.lease_seconds = int(lease_seconds)
        if self.lease_seconds <= 0:
            raise ValueError("lease_seconds must be positive")
        if not callable(runtime_session_id_factory):
            raise TypeError("runtime_session_id_factory must be callable")
        self._runtime_session_id_factory = runtime_session_id_factory
        self.lease_interval = float(lease_interval)
        self.worker_heartbeat_interval = float(worker_heartbeat_interval)
        self.session_heartbeat_interval = float(session_heartbeat_interval)
        self.retry_interval = float(retry_interval)
        self.clock = clock
        self.log = logger or logging.getLogger(__name__)

        self._executor = executor or _ThreadExecutor()
        self._owns_executor = executor is None
        self._lock = threading.RLock()
        self._sessions: Dict[str, _ActiveSession] = {}
        self._registered = False
        self._network_healthy = False
        self._last_worker_heartbeat: Optional[float] = None
        self._last_error: Optional[str] = None
        self._draining = False
        self._stopped = False
        self._stop_event = threading.Event()
        self._threads = []

    @property
    def network_healthy(self) -> bool:
        with self._lock:
            return self._network_healthy

    @property
    def draining(self) -> bool:
        with self._lock:
            return self._draining

    @property
    def stopped(self) -> bool:
        with self._lock:
            return self._stopped

    @property
    def last_error(self) -> Optional[str]:
        with self._lock:
            return self._last_error

    def active_session_ids(self):
        with self._lock:
            return sorted(self._sessions)

    def active_count(self) -> int:
        with self._lock:
            return len(self._sessions)

    def available_slots(self) -> int:
        return self.capacity_manager.available_slots()

    def _mark_network_failure(self, exc: BaseException) -> None:
        with self._lock:
            self._network_healthy = False
            self._registered = False
            self._last_error = type(exc).__name__

    def _mark_network_healthy(self) -> None:
        with self._lock:
            self._network_healthy = True
            self._last_error = None

    def _worker_payload(self, status: str) -> Dict[str, Any]:
        return {
            "workerKey": self.worker_key,
            "host": self.host,
            "bootId": self.boot_id,
            "processUuid": self.process_uuid,
            "status": status,
            "maxSlots": self.max_slots,
            "activeSessions": self.active_count(),
            "freeSlots": self.available_slots(),
            "capabilities": dict(self.capabilities),
        }

    def _register(self) -> bool:
        try:
            self.client.register_worker(
                self._worker_payload("DRAINING" if self.draining else "ACTIVE"),
                request_id="jarvis-worker-register-" + hashlib.sha256(
                    self.worker_key.encode("utf-8")).hexdigest()[:32])
        except Exception as exc:
            self._mark_network_failure(exc)
            self.log.warning("worker register failed worker=%s error=%s",
                             self.worker_key, type(exc).__name__)
            return False
        with self._lock:
            self._registered = True
        self._mark_network_healthy()
        return True

    def _heartbeat_worker_if_due(self, *, force: bool = False) -> bool:
        now = self.clock()
        with self._lock:
            due = (force or self._last_worker_heartbeat is None or
                   now - self._last_worker_heartbeat >= self.worker_heartbeat_interval)
        if not due:
            return True
        try:
            self.client.heartbeat_worker(
                self.worker_key,
                self._worker_payload("DRAINING" if self.draining else "ACTIVE"),
                request_id="jarvis-worker-heartbeat-%s" % uuid.uuid4().hex)
        except Exception as exc:
            self._mark_network_failure(exc)
            self.log.warning("worker heartbeat failed worker=%s error=%s",
                             self.worker_key, type(exc).__name__)
            return False
        with self._lock:
            self._last_worker_heartbeat = now
        self._mark_network_healthy()
        return True

    def run_once(self) -> bool:
        """Register/heartbeat and perform at most one lease poll."""
        with self._lock:
            if self._stopped:
                return False
            registered = self._registered
        if not registered and not self._register():
            return False
        if not self._heartbeat_worker_if_due(force=not self.network_healthy):
            return False
        if self.draining:
            return True
        with self._lock:
            has_pending_terminal = any(
                record.controller.pending_terminal for record in self._sessions.values())
        if has_pending_terminal:
            # Recover and commit already-finished work before admitting another
            # lease after a partial control-plane outage.
            return True
        capacity_permit = self.capacity_manager.acquire(
            "task-executor:%s" % self.worker_key)
        if capacity_permit is None:
            return True
        try:
            response = self.client.lease_task(
                self.worker_key, capabilities=self.capabilities,
                lease_seconds=self.lease_seconds,
                request_id="jarvis-worker-lease-%s" % uuid.uuid4().hex)
            lease = parse_lease_response(response)
        except ControlPlaneConflict as exc:
            if capacity_permit is not None:
                capacity_permit.release()
            # A concurrent poll or backend max-slot fence is capacity pressure,
            # not a network outage.  Keep the registration healthy and retry on
            # the normal poll cadence.
            self._mark_network_healthy()
            self.log.warning("task lease conflict worker=%s code=%s",
                             self.worker_key, getattr(exc, "code", ""))
            return True
        except Exception as exc:
            if capacity_permit is not None:
                capacity_permit.release()
            self._mark_network_failure(exc)
            self.log.warning("task lease failed worker=%s error=%s",
                             self.worker_key, type(exc).__name__)
            return False
        self._mark_network_healthy()
        if lease is None:
            if capacity_permit is not None:
                capacity_permit.release()
            return True
        try:
            return self._accept_lease(lease, capacity_permit=capacity_permit)
        except Exception as exc:
            if capacity_permit is not None:
                capacity_permit.release()
            # A controller/client integration bug must fail closed instead of
            # killing the long-running lease-loop thread.
            self._mark_network_failure(exc)
            self.log.exception("leased task setup failed worker=%s error=%s",
                               self.worker_key, type(exc).__name__)
            return False

    def _accept_lease(self, lease: Mapping[str, Any], *,
                      capacity_permit: Optional[CapacityPermit] = None) -> bool:
        if capacity_permit is None:
            capacity_permit = self.capacity_manager.acquire(
                "task-executor:%s" % self.worker_key)
            if capacity_permit is None:
                return False
        controller = SessionController(
            self.client, self.worker_key, lease,
            lease_seconds=self.lease_seconds,
            runtime_session_id_factory=self._runtime_session_id_factory,
            stop_process=self.stop_process,
            on_network_failure=self._mark_network_failure,
            logger=self.log,
        )
        with self._lock:
            if self._draining or self._stopped:
                self._fail_worker_stopping(controller, "worker_not_accepting")
                if capacity_permit is not None:
                    capacity_permit.release()
                return False
            if controller.session_id in self._sessions:
                self.log.warning("duplicate lease ignored session=%s", controller.session_id)
                if capacity_permit is not None:
                    capacity_permit.release()
                return False
        if not controller.start():
            if capacity_permit is not None:
                capacity_permit.release()
            return False
        record = _ActiveSession(controller, capacity_permit=capacity_permit)
        with self._lock:
            if self._draining or self._stopped:
                self._fail_worker_stopping(controller, "worker_not_accepting")
                if capacity_permit is not None:
                    capacity_permit.release()
                return False
            self._sessions[controller.session_id] = record
        try:
            future = self._executor.submit(self._execute_one, controller)
        except Exception as exc:
            if not (controller.terminal or controller.ownership_lost or
                    controller.stop_requested):
                controller.fail({"errorType": type(exc).__name__,
                                "message": "execution submission failed"})
            if controller.terminal or controller.ownership_lost:
                self._remove_session(controller.session_id)
            elif capacity_permit is not None:
                # Submission failed before a terminal ACK could be accepted.
                # The Session remains server-owned until lease expiry, but this
                # process no longer consumes a local execution slot.
                capacity_permit.release()
            return False
        with self._lock:
            current = self._sessions.get(controller.session_id)
            if current is not None:
                current.future = future
        return True

    @staticmethod
    def _fail_worker_stopping(controller: SessionController, reason: str) -> None:
        """Return an accidentally-held lease to retryable state, never SUSPENDED.

        A worker shutdown is not an external wait condition.  Encoding it as a
        suspend without ``waitExpireAt`` leaves a permanent task that no sensor can
        wake, so use the normal retryable failure transition instead.  Stop the
        owned process synchronously before making the task retryable; otherwise a
        fast reaper/new worker could start the replacement while the old process is
        still inside its graceful-termination window.
        """
        controller.request_stop("worker_stopping")
        if controller.pending_terminal:
            # The work already has a stable complete/fail/suspend intent.  Never
            # replace it with WorkerStopping; retry the same idempotent transition
            # once and otherwise let lease expiry preserve the frozen attempt.
            controller.retry_terminal()
            return
        if not controller.terminal and not controller.ownership_lost:
            controller.fail({
                "error": {
                    "errorType": "WorkerStopping",
                    "message": str(reason),
                },
                "retryAfterSeconds": 0,
            })

    @staticmethod
    def _result_detail(result: Any) -> Dict[str, Any]:
        if result is None:
            return {}
        if isinstance(result, Mapping):
            return dict(result)
        return {"result": str(result)}

    def _execute_one(self, controller: SessionController) -> None:
        try:
            # A queued callback may not start until after a heartbeat has already
            # invalidated its fence.  Never launch work we no longer own.
            if controller.ownership_lost or controller.stop_requested:
                return
            result = self.execute(controller.lease, controller)
            if controller.terminal or controller.pending_terminal or controller.stop_requested:
                return
            detail = self._result_detail(result)
            status = str(detail.get("status") or detail.get("result") or "").lower()
            if status in {"suspend", "suspended"}:
                controller.suspend({key: value for key, value in detail.items()
                                   if key not in {"status", "result"}})
            elif status in {"error", "fail", "failed"} or result is False:
                retry_after = detail.get("retryAfterSeconds")
                if retry_after is None:
                    retry_after = detail.get("retry_after_seconds")
                error = detail.get("error")
                if error is None:
                    error = {key: value for key, value in detail.items()
                             if key not in {"status", "retryAfterSeconds",
                                            "retry_after_seconds"}}
                controller.fail(error, retry_after_seconds=retry_after)
            else:
                controller.complete(detail)
        except Exception as exc:
            if not (controller.terminal or controller.pending_terminal or
                    controller.ownership_lost or controller.stop_requested):
                controller.fail({
                    "errorType": type(exc).__name__,
                    "message": str(exc)[:500],
                })
        finally:
            if controller.terminal or controller.ownership_lost:
                self._remove_session(controller.session_id)

    def _remove_session(self, session_id: str) -> None:
        with self._lock:
            record = self._sessions.pop(str(session_id), None)
        if record is not None and record.capacity_permit is not None:
            record.capacity_permit.release()

    @staticmethod
    def _future_done(record: _ActiveSession) -> bool:
        future = record.future
        if future is None:
            return True
        try:
            return bool(future.done())
        except Exception:
            return False

    def heartbeat_sessions_once(self) -> None:
        """Renew every active session independently from worker/lease polling."""
        with self._lock:
            records = list(self._sessions.values())
        for record in records:
            controller = record.controller
            if controller.pending_terminal:
                controller.retry_terminal()
            elif not controller.terminal and not controller.ownership_lost:
                controller.heartbeat()
            if ((controller.terminal or controller.ownership_lost) and
                    self._future_done(record)):
                self._remove_session(controller.session_id)

    def _lease_loop(self) -> None:
        while not self._stop_event.is_set():
            self.run_once()
            interval = self.lease_interval if self.network_healthy else self.retry_interval
            self._stop_event.wait(interval)

    def _session_loop(self) -> None:
        while not self._stop_event.wait(self.session_heartbeat_interval):
            self.heartbeat_sessions_once()

    def start(self) -> "PersistenceExecutor":
        with self._lock:
            if self._stopped:
                raise RuntimeError("worker is stopped")
            if self._threads:
                return self
            self._threads = [
                threading.Thread(target=self._lease_loop, name="jarvis-lease-loop", daemon=True),
                threading.Thread(target=self._session_loop, name="jarvis-session-loop", daemon=True),
            ]
            threads = list(self._threads)
        for thread in threads:
            thread.start()
        return self

    def run_forever(self) -> None:
        self.start()
        while not self._stop_event.wait(1.0):
            pass

    def drain(self, timeout: Optional[float] = None) -> bool:
        """Stop leasing and wait until all active executions leave the worker."""
        with self._lock:
            self._draining = True
            registered = self._registered
        if registered:
            self._heartbeat_worker_if_due(force=True)
        deadline = None if timeout is None else self.clock() + max(0.0, float(timeout))
        while self.active_count() > 0:
            if deadline is not None and self.clock() >= deadline:
                return False
            self._stop_event.wait(0.05)
        return True

    def stop(self, *, drain: bool = False, timeout: Optional[float] = None) -> bool:
        """Stop loops; optionally drain, otherwise stop and retry-fail active work."""
        if drain:
            drained = self.drain(timeout)
        else:
            with self._lock:
                self._draining = True
                registered = self._registered
            if registered:
                self._heartbeat_worker_if_due(force=True)
            drained = self.active_count() == 0
        self._stop_event.set()
        if not drained:
            with self._lock:
                controllers = [record.controller for record in self._sessions.values()]
            for controller in controllers:
                self._fail_worker_stopping(controller, "worker_stopping")
        with self._lock:
            self._stopped = True
            threads = list(self._threads)
        current = threading.current_thread()
        for thread in threads:
            if thread is not current:
                thread.join(timeout=1.0)
        if self._owns_executor:
            try:
                self._executor.shutdown(wait=drained, cancel_futures=not drained)
            except TypeError:  # pragma: no cover - Python <3.9
                self._executor.shutdown(wait=drained)
        return drained
