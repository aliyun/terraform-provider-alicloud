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
import fcntl
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


def load_or_create_worker_key(path: Optional[Any] = None) -> str:
    """Return the installation-scoped Worker id stored on this machine."""
    configured = path or os.environ.get("JARVIS_WORKER_ID_FILE")
    worker_path = Path(configured or ".my-day/bridge/worker-id").expanduser()
    worker_path.parent.mkdir(parents=True, exist_ok=True)
    lock_path = worker_path.with_name(worker_path.name + ".lock")
    with lock_path.open("a+") as lock_file:
        os.chmod(lock_path, 0o600)
        fcntl.flock(lock_file.fileno(), fcntl.LOCK_EX)
        if worker_path.exists():
            value = _nonblank(worker_path.read_text().strip(), "worker_id_file")
            os.chmod(worker_path, 0o600)
            return _component(value, "worker_id_file")
        value = "worker-" + uuid.uuid4().hex
        temp_path = worker_path.with_name(
            "%s.tmp.%s.%s" % (worker_path.name, os.getpid(), uuid.uuid4().hex))
        fd = os.open(str(temp_path), os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
        try:
            with os.fdopen(fd, "w") as stream:
                stream.write(value + "\n")
                stream.flush()
                os.fsync(stream.fileno())
            os.replace(temp_path, worker_path)
        finally:
            try:
                temp_path.unlink()
            except FileNotFoundError:
                pass
        return value


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
    can be retried with the same idempotency key.  A stale fence immediately
    revokes ownership.  A transient heartbeat transport failure may keep running
    only while the last successful fenced renewal still proves enough local lease
    margin; crossing that boundary stops the process fail-closed.
    """

    TERMINAL_ACTIONS = {"complete", "fail", "suspend", "relinquish"}

    def __init__(self, client: ControlPlaneClient, worker_key: str,
                 lease: Mapping[str, Any], *,
                 process_uuid: Optional[str] = None,
                 lease_seconds: int = 300,
                 lease_safety_margin: float = 90.0,
                 runtime_session_id_factory: Callable[[], str] = _new_runtime_session_id,
                 stop_process: Optional[StopProcess] = None,
                 on_network_failure: Optional[NetworkFailure] = None,
                 clock: Callable[[], float] = time.monotonic,
                 logger: Optional[logging.Logger] = None):
        self.client = client
        self.worker_key = _nonblank(worker_key, "worker_key")
        self.process_uuid = _nonblank(process_uuid or "unknown", "process_uuid")
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
        self.lease_safety_margin = float(lease_safety_margin)
        if self.lease_safety_margin < 0:
            raise ValueError("lease_safety_margin must not be negative")
        if self.lease_safety_margin >= self.lease_seconds:
            raise ValueError("lease_safety_margin must be less than lease_seconds")
        if not callable(clock):
            raise TypeError("clock must be callable")
        self.clock = clock
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
        # Monotonic local proof, refreshed only after a successful fenced
        # start/heartbeat.  The server also returns leaseExpireAt; it is retained
        # in ``session`` for diagnostics, while the deadline uses request-start +
        # TTL so wall-clock skew cannot extend ownership.
        self._lease_deadline: Optional[float] = None

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

    @property
    def lease_deadline(self) -> Optional[float]:
        with self._lock:
            return self._lease_deadline

    @staticmethod
    def _response_session(response: Any) -> Optional[Mapping[str, Any]]:
        if not isinstance(response, Mapping):
            return None
        nested = response.get("session")
        return nested if isinstance(nested, Mapping) else response

    def _refresh_lease_proof(self, response: Any, request_started_at: float) -> None:
        response_session = self._response_session(response)
        lease_expire_at = (_field(
            response_session, "leaseExpireAt", "leaseExpiresAt", "lease_expire_at")
            if response_session is not None else None)
        with self._lock:
            self._lease_deadline = float(request_started_at) + self.lease_seconds
            if lease_expire_at is not None:
                self.session["leaseExpireAt"] = lease_expire_at
                self.lease["session"] = self.session

    def _lease_proof_has_margin(self) -> bool:
        with self._lock:
            deadline = self._lease_deadline
        return (deadline is not None
                and self.clock() < deadline - self.lease_safety_margin)

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
                request_started_at = self.clock()
                pid = int(getattr(process, "pid", 0))
                if pid <= 0:
                    raise ValueError("bound process must expose a positive pid")
                response = self.client.start_session(
                    self.session_id, self.worker_key, self.fence_token,
                    self._start_detail({"pid": pid}),
                    process_uuid=self.process_uuid,
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
            self._refresh_lease_proof(response, request_started_at)
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
                request_started_at = self.clock()
                response = self.client.start_session(
                    self.session_id, self.worker_key, self.fence_token,
                    self._start_detail(detail), process_uuid=self.process_uuid,
                    request_id=self._request_id("start"))
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
            self._refresh_lease_proof(response, request_started_at)
            return True

    def heartbeat(self, detail: Optional[Mapping[str, Any]] = None) -> bool:
        with self._transition_lock:
            with self._lock:
                if self._terminal_action is not None:
                    return True
                if self._ownership_lost or self._pending_terminal is not None:
                    return False
            try:
                request_started_at = self.clock()
                response = self.client.heartbeat_session(
                    self.session_id, self.worker_key, self.fence_token,
                    self._lease_detail(detail),
                    process_uuid=self.process_uuid,
                    request_id="jarvis-session-heartbeat-%s" % uuid.uuid4().hex)
                self._refresh_lease_proof(response, request_started_at)
                return True
            except StaleFence:
                self._lose_ownership("stale_fence:heartbeat")
                return False
            except ControlPlaneUnavailable as exc:
                self._notify_network_failure(exc)
                if self._lease_proof_has_margin():
                    self.log.warning(
                        "session heartbeat temporarily unavailable session=%s; "
                        "retaining ownership within lease proof",
                        self.session_id)
                    return True
                self._lose_ownership("lease_proof_expiring:heartbeat")
                self.log.warning("session heartbeat failed session=%s error=%s",
                                 self.session_id, type(exc).__name__)
                return False
            except ControlPlaneError as exc:
                self._lose_ownership("heartbeat_rejected")
                self.log.warning("session heartbeat rejected session=%s code=%s",
                                 self.session_id, getattr(exc, "code", ""))
                return False
            except Exception as exc:
                self._lose_ownership("heartbeat_failed")
                self.log.warning("session heartbeat failed session=%s error=%s",
                                 self.session_id, type(exc).__name__)
                return False

    def adopt_lease(self, lease: Mapping[str, Any]) -> bool:
        """Adopt the fence token from a re-issued lease for this same session.

        The control plane may re-offer a session this worker is *already*
        running (e.g. after the ownership heartbeat is delayed past the lease
        window during a network wobble).  Every lease grant rotates the fence
        token, so silently discarding the re-lease would leave this running
        controller holding a fence the server has already superseded — the very
        next ``heartbeat`` would then be rejected 412 STALE_FENCE and the bound
        process killed mid-flight.  Adopting the new fence keeps ownership
        continuous so the pending heartbeat renews against the current token.

        Returns True when the fence was adopted, False when the session is
        already terminal / ownership was lost (a re-lease must never resurrect
        a controller that has stopped or is committing a terminal transition).
        """
        incoming = _mapping(lease, "lease")
        incoming_session = _mapping(incoming.get("session"), "lease.session")
        incoming_id = _nonblank(
            _field(incoming_session, "id", "sessionId", "session_id"),
            "session_id")
        if incoming_id != self.session_id:
            raise LeaseProtocolError(
                "adopt_lease session mismatch: %s != %s" % (
                    incoming_id, self.session_id))
        new_fence = _field(incoming_session, "fenceToken", "fence_token")
        if new_fence is None or str(new_fence).strip() == "":
            raise LeaseProtocolError("fence_token must not be empty")
        with self._transition_lock:
            with self._lock:
                if (self._terminal_action is not None
                        or self._ownership_lost
                        or self._pending_terminal is not None):
                    return False
                # Swap the fence only; the deadline is left to the next
                # successful heartbeat's proof refresh so wall-clock skew can
                # never extend ownership from an unproven local timestamp.
                self.fence_token = new_fence
                self.session["fenceToken"] = new_fence
                self.lease["session"] = self.session
                return True

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
                       process_uuid=self.process_uuid,
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

    def relinquish(self, reason: str = "worker_stopping") -> bool:
        return self._terminal("relinquish", {"reason": str(reason)})


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
                 worker_id_file: Optional[Any] = None,
                 capabilities: Optional[Mapping[str, Any]] = None,
                 lease_seconds: int = 300,
                 lease_safety_margin: float = 90.0,
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
            worker_key or load_or_create_worker_key(worker_id_file), "worker_key")
        legacy_parts = self.worker_key.split(":", 2) if worker_key else []
        self.host = _component(
            host or (legacy_parts[0] if len(legacy_parts) == 3 else socket.gethostname()),
            "host")
        self.boot_id = _component(
            boot_id or (legacy_parts[1] if len(legacy_parts) == 3
                        else _default_boot_id(self.host)), "boot_id")
        self.process_uuid = _component(
            process_uuid or (legacy_parts[2] if len(legacy_parts) == 3
                             else uuid.uuid4().hex), "process_uuid")
        self.capabilities = dict(capabilities or {})
        self.lease_seconds = int(lease_seconds)
        if self.lease_seconds <= 0:
            raise ValueError("lease_seconds must be positive")
        self.lease_safety_margin = float(lease_safety_margin)
        if self.lease_safety_margin < 0:
            raise ValueError("lease_safety_margin must not be negative")
        if self.lease_safety_margin >= self.lease_seconds:
            raise ValueError("lease_safety_margin must be less than lease_seconds")
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
                process_uuid=self.process_uuid,
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
                process_uuid=self.process_uuid,
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
            process_uuid=self.process_uuid,
            lease_seconds=self.lease_seconds,
            lease_safety_margin=self.lease_safety_margin,
            runtime_session_id_factory=self._runtime_session_id_factory,
            stop_process=self.stop_process,
            on_network_failure=self._mark_network_failure,
            clock=self.clock,
            logger=self.log,
        )
        with self._lock:
            if self._draining or self._stopped:
                controller.start()
                self._fail_worker_stopping(controller, "worker_not_accepting")
                if capacity_permit is not None:
                    capacity_permit.release()
                return False
            existing = self._sessions.get(controller.session_id)
            if existing is not None:
                # A re-lease for a session we already run: adopt the rotated
                # fence into the live controller instead of discarding it, so
                # the pending heartbeat renews against the current token rather
                # than losing ownership to a self-inflicted STALE_FENCE.
                adopted = existing.controller.adopt_lease(lease)
                if adopted:
                    self.log.warning(
                        "re-lease fence adopted session=%s", controller.session_id)
                else:
                    self.log.warning(
                        "duplicate lease ignored session=%s (controller terminal)",
                        controller.session_id)
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
        """Stop the local process and relinquish without consuming task retry budget."""
        controller.request_stop("worker_stopping")
        if controller.pending_terminal:
            # The work already has a stable complete/fail/suspend intent.  Never
            # replace it with WorkerStopping; retry the same idempotent transition
            # once and otherwise let lease expiry preserve the frozen attempt.
            controller.retry_terminal()
            return
        if not controller.terminal and not controller.ownership_lost:
            controller.relinquish(reason)

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
        """Stop loops, relinquish active Sessions, then mark this Worker OFFLINE."""
        with self._lock:
            registered = self._registered
        if drain:
            drained = self.drain(timeout)
        else:
            with self._lock:
                self._draining = True
            if registered:
                self._heartbeat_worker_if_due(force=True)
            drained = self.active_count() == 0
        self._stop_event.set()
        if not drained:
            with self._lock:
                controllers = [record.controller for record in self._sessions.values()]
            for controller in controllers:
                self._fail_worker_stopping(controller, "worker_stopping")
                if controller.terminal or controller.ownership_lost:
                    self._remove_session(controller.session_id)
            drained = self.active_count() == 0
        offline = not registered
        if registered and drained:
            try:
                self.client.heartbeat_worker(
                    self.worker_key, self._worker_payload("OFFLINE"),
                    process_uuid=self.process_uuid,
                    request_id="jarvis-worker-offline-%s" % uuid.uuid4().hex)
                offline = True
            except Exception as exc:
                self._mark_network_failure(exc)
                self.log.warning("worker offline failed worker=%s error=%s",
                                 self.worker_key, type(exc).__name__)
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
        return drained and offline
