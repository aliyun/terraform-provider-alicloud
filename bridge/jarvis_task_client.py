#!/usr/bin/env python3
"""HTTP client primitives for the Jarvis control plane.

This module deliberately contains no scheduler or worker-loop policy.  Sensors
serialize :class:`TaskEnvelope` values through it and PersistenceExecutor uses the
Task/Session/Worker/Operation methods.  Keeping the HTTP boundary small also
keeps rollout policy outside the control-plane contract.
"""

from __future__ import annotations

import hashlib
import json
import uuid
from dataclasses import dataclass, field
from typing import Any, Callable, Dict, Mapping, Optional
from urllib.error import HTTPError, URLError
from urllib.parse import quote, urlencode
from urllib.request import Request, urlopen


class ControlPlaneError(RuntimeError):
    """Base class for a rejected or unavailable control-plane operation."""

    def __init__(self, message: str, *, status: Optional[int] = None,
                 code: str = "", response: Optional[Mapping[str, Any]] = None):
        super().__init__(message)
        self.status = status
        self.code = code
        self.response = dict(response or {})


class ControlPlaneUnavailable(ControlPlaneError):
    """The data plane could not be reached or returned a transient 5xx."""


class ControlPlaneConflict(ControlPlaneError):
    """The request conflicts with the current task/session state (HTTP 409)."""


class StaleFence(ControlPlaneError):
    """The caller no longer owns the session lease (HTTP 412/STALE_FENCE)."""


class HandoffRequested(StaleFence):
    """The lease was fenced for an operator-requested two-phase handoff."""


class InvalidResponse(ControlPlaneError):
    """The server returned a successful but malformed response."""


class ControlPlaneRejected(ControlPlaneError):
    """A non-retryable HTTP rejection not covered by a more specific class."""


def _nonblank(value: Any, name: str) -> str:
    text = str(value or "").strip()
    if not text:
        raise ValueError("%s must not be empty" % name)
    return text


def _aone_id_from_task_key(task_key: Any) -> Optional[str]:
    """Return the work-item id encoded by a canonical Aone task key.

    Aone project and work-item ids are numeric.  Requiring both components keeps
    local keys such as ``aone:unknown:adhoc`` from being misclassified merely
    because they share the prefix.
    """
    parts = str(task_key or "").split(":")
    if (len(parts) == 3 and parts[0] == "aone"
            and parts[1].isdigit() and parts[2].isdigit()):
        return parts[2]
    return None


def _camel_key(key: Any) -> Any:
    """Translate Python-style request field names to the API's camelCase."""
    if not isinstance(key, str) or "_" not in key:
        return key
    head, *tail = key.split("_")
    return head + "".join(part[:1].upper() + part[1:] for part in tail if part)


def _camelize(value: Any) -> Any:
    if isinstance(value, Mapping):
        return {_camel_key(key): _camelize(item) for key, item in value.items()}
    if isinstance(value, (list, tuple)):
        return [_camelize(item) for item in value]
    return value


@dataclass(frozen=True)
class TaskEnvelope:
    """Persistent desired-task input shared by every Jarvis sensor.

    ``task_key`` is the long-lived logical mutex (for example
    ``aone:2100304:84308492``).  ``desired_revision`` distinguishes subsequent
    facts observed for that task without creating another concurrent identity.
    The payload must contain only resumable execution inputs; credentials never
    belong in it.
    """

    task_key: str
    source_type: str
    source_ref: Mapping[str, Any]
    task_type: str
    desired_revision: str
    payload: Mapping[str, Any]
    trigger_mask: Any = field(default_factory=list)
    recovery_policy: str = "REPLAY_SAFE"
    persona: Optional[str] = None
    priority: Optional[Any] = None
    aone_id: Optional[str] = None
    comment_cursor: Optional[Any] = None
    required_capabilities: Optional[Any] = None
    max_retries: Optional[int] = None
    source_status: Optional[str] = None

    def __post_init__(self) -> None:
        object.__setattr__(self, "task_key", _nonblank(self.task_key, "task_key"))
        object.__setattr__(self, "source_type", _nonblank(self.source_type, "source_type"))
        object.__setattr__(self, "task_type", _nonblank(self.task_type, "task_type"))
        object.__setattr__(self, "desired_revision",
                           _nonblank(self.desired_revision, "desired_revision"))
        object.__setattr__(self, "recovery_policy",
                           _nonblank(self.recovery_policy, "recovery_policy"))
        if not isinstance(self.source_ref, Mapping):
            raise TypeError("source_ref must be a mapping")
        if not isinstance(self.payload, Mapping):
            raise TypeError("payload must be a mapping")
        canonical_aone_id = _aone_id_from_task_key(self.task_key)
        explicit_aone_id = str(self.aone_id or "").strip() or None
        if (canonical_aone_id and explicit_aone_id
                and explicit_aone_id != canonical_aone_id):
            raise ValueError(
                "aone_id %s conflicts with canonical task_key %s"
                % (explicit_aone_id, self.task_key))
        if canonical_aone_id and not explicit_aone_id:
            object.__setattr__(self, "aone_id", canonical_aone_id)
        source_status = str(self.source_status or "").strip() or None
        if source_status and len(source_status) > 64:
            raise ValueError("source_status must not exceed 64 characters")
        object.__setattr__(self, "source_status", source_status)

    def to_dict(self) -> Dict[str, Any]:
        data: Dict[str, Any] = {
            "taskKey": self.task_key,
            "sourceType": self.source_type,
            "sourceRef": _camelize(self.source_ref),
            "taskType": self.task_type,
            "desiredRevision": self.desired_revision,
            "triggerMask": _camelize(self.trigger_mask),
            "payload": dict(self.payload),
            "recoveryPolicy": self.recovery_policy,
        }
        if self.persona:
            data["persona"] = self.persona
        if self.priority:
            data["priority"] = self.priority
        if self.aone_id:
            data["aoneId"] = str(self.aone_id)
        if self.comment_cursor is not None:
            data["commentCursor"] = self.comment_cursor
        if self.required_capabilities is not None:
            data["requiredCapabilities"] = _camelize(self.required_capabilities)
        if self.max_retries is not None:
            data["maxRetries"] = int(self.max_retries)
        if self.source_status:
            data["sourceStatus"] = self.source_status
        return _camelize(data)

    def request_id(self, operation: str) -> str:
        """Return a stable idempotency key for one logical envelope operation."""
        material = {
            "operation": _nonblank(operation, "operation"),
            "task": self.to_dict(),
        }
        raw = json.dumps(material, ensure_ascii=False, sort_keys=True,
                         separators=(",", ":"), default=str).encode("utf-8")
        return "jarvis-%s-%s" % (operation, hashlib.sha256(raw).hexdigest()[:32])


class ControlPlaneClient:
    """Small JSON-over-HTTP client for Jarvis control-plane APIs."""

    DEFAULT_PREFIX = "/api/jarvis/v1"
    PATHS = {
        "task_upsert": "tasks/upsert",
        "task_lease": "tasks/lease",
        "task_claim": "tasks/claim",
        "worker_register": "workers/register",
        "operation_begin": "operations/begin",
        "operation_ack": "operations/ack",
        "operation_fail": "operations/fail",
        "operation_not_started": "operations/not-started",
        "operation_reconcile": "operations/reconcile",
    }
    WORKER_HEARTBEAT_PATH = "workers/{worker_key}/heartbeat"
    WORKER_HANDOFF_REQUESTS_PATH = "workers/{worker_key}/handoff-requests"
    WORKERS_PATH = "workers"
    WORKER_STATE_PATH = "workers/{worker_key}/state"
    SESSION_ACTION_PATH = "sessions/{session_id}/{action}"
    TASK_BY_AONE_PATH = "tasks/by-aone/{aone_id}"
    TASK_TIMELINE_PATH = "tasks/{task_id}/timeline"
    DISCARD_RESUME_CONTEXT_PATH = "tasks/{task_id}/discard-resume-context"
    READY_TASK_DIAGNOSTICS_PATH = "tasks/ready-diagnostics"
    SOURCE_STATUS_CANDIDATES_PATH = "tasks/source-status-candidates"
    SOURCE_STATUS_PATH = "tasks/{task_id}/source-status"
    PENDING_AONE_WAITS_PATH = "sessions/waits/aone-reply"
    EXTERNAL_OPERATION_RECOVERY_CANDIDATES_PATH = (
        "operations/external-recovery-candidates")
    OPERATION_RECOVERY_LEASE_PATH = "operations/recovery/lease"
    OPERATION_RECOVERY_RENEW_PATH = "operations/recovery/renew"
    OPERATION_RECOVERY_RELEASE_PATH = "operations/recovery/release"
    OPERATION_PATH = "operations/{operation_id}"
    OPERATION_BY_KEY_PATH = "operations/by-key"

    def __init__(self, base_url: str, token: str = "", *, timeout: float = 10.0,
                 api_prefix: str = DEFAULT_PREFIX,
                 opener: Optional[Callable[..., Any]] = None):
        self.base_url = _nonblank(base_url, "base_url").rstrip("/")
        self.token = str(token or "")
        self.timeout = float(timeout)
        if self.timeout <= 0:
            raise ValueError("timeout must be positive")
        prefix = "/" + str(api_prefix or self.DEFAULT_PREFIX).strip("/")
        self.api_prefix = prefix.rstrip("/")
        self._opener = opener or urlopen

    @staticmethod
    def new_request_id() -> str:
        return "jarvis-%s" % uuid.uuid4()

    def _endpoint(self, suffix: str) -> str:
        return self.base_url + self.api_prefix + "/" + suffix.lstrip("/")

    @staticmethod
    def _path_segment(value: Any, name: str) -> str:
        return quote(_nonblank(value, name), safe="")

    @staticmethod
    def _response_detail(raw: bytes) -> Dict[str, Any]:
        if not raw:
            return {}
        try:
            parsed = json.loads(raw.decode("utf-8", errors="replace"))
        except (TypeError, ValueError):
            return {}
        return parsed if isinstance(parsed, dict) else {}

    @staticmethod
    def _error_message(status: Optional[int], detail: Mapping[str, Any]) -> str:
        code = str(detail.get("code") or detail.get("errorCode")
                   or detail.get("error_code") or "")
        message = str(detail.get("message") or detail.get("error") or "request rejected")
        prefix = "control plane"
        if status is not None:
            prefix += " HTTP %s" % status
        if code:
            prefix += " %s" % code
        return "%s: %s" % (prefix, message[:300])

    def _raise_http(self, status: int, raw: bytes) -> None:
        detail = self._response_detail(raw)
        code = str(detail.get("code") or detail.get("errorCode")
                   or detail.get("error_code") or "")
        message = self._error_message(status, detail)
        kwargs = {"status": status, "code": code, "response": detail}
        if code.lower() == "preconditionfailed.handoffrequested":
            raise HandoffRequested(message, **kwargs)
        if status == 412 or code.upper() == "STALE_FENCE":
            raise StaleFence(message, **kwargs)
        if status == 409:
            raise ControlPlaneConflict(message, **kwargs)
        if status >= 500:
            raise ControlPlaneUnavailable(message, **kwargs)
        raise ControlPlaneRejected(message, **kwargs)

    def _request(self, method: str, suffix: str, *, payload: Optional[Mapping[str, Any]] = None,
                 request_id: Optional[str] = None) -> Any:
        method = _nonblank(method, "method").upper()
        if payload is not None and not isinstance(payload, Mapping):
            raise TypeError("payload must be a mapping")
        body = None
        if payload is not None:
            body = json.dumps(_camelize(payload), ensure_ascii=False, separators=(",", ":"),
                              default=str).encode("utf-8")
        headers = {
            "Accept": "application/json",
            "User-Agent": "jarvis-control-plane/1",
        }
        if method != "GET":
            headers["Content-Type"] = "application/json"
            headers["Idempotency-Key"] = _nonblank(
                request_id or self.new_request_id(), "request_id")
        if self.token:
            headers["Authorization"] = "Bearer " + self.token
        req = Request(self._endpoint(suffix), data=body, method=method, headers=headers)
        try:
            with self._opener(req, timeout=self.timeout) as resp:
                status_value = getattr(resp, "status", None)
                if status_value is None:
                    status_value = resp.getcode()
                status = int(status_value)
                raw = resp.read() or b""
        except HTTPError as exc:
            try:
                raw = exc.read() or b""
            except Exception:  # pragma: no cover - defensive for broken HTTPError fakes
                raw = b""
            self._raise_http(int(exc.code), raw)
            raise AssertionError("unreachable")
        except (URLError, TimeoutError, OSError) as exc:
            # Do not include request headers/body here: they may carry credentials or payload data.
            raise ControlPlaneUnavailable(
                "control plane unavailable: %s" % type(exc).__name__) from exc

        if status < 200 or status >= 300:
            self._raise_http(status, raw)
        if not raw:
            return {}
        try:
            parsed = json.loads(raw.decode("utf-8", errors="strict"))
        except (UnicodeDecodeError, ValueError, TypeError) as exc:
            raise InvalidResponse("control plane returned invalid JSON", status=status) from exc
        if not isinstance(parsed, (dict, list)):
            raise InvalidResponse(
                "control plane response must be an object or array", status=status)
        return parsed

    def _post(self, suffix: str, payload: Mapping[str, Any], *,
              request_id: Optional[str] = None) -> Dict[str, Any]:
        parsed = self._request("POST", suffix, payload=payload, request_id=request_id)
        if not isinstance(parsed, dict):
            raise InvalidResponse("control plane mutation response must be an object")
        return parsed

    def _get(self, suffix: str) -> Any:
        return self._request("GET", suffix)

    def upsert_task(self, envelope: TaskEnvelope, *,
                    request_id: Optional[str] = None) -> Dict[str, Any]:
        rid = request_id or envelope.request_id("upsert")
        return self._post(self.PATHS["task_upsert"],
                          envelope.to_dict(), request_id=rid)

    def upsert_desired_task(self, envelope: TaskEnvelope, *,
                            request_id: Optional[str] = None) -> Dict[str, Any]:
        """Persist the latest desired revision for one recoverable Task."""
        return self.upsert_task(envelope, request_id=request_id)

    def register_worker(self, worker: Mapping[str, Any], *,
                        request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(self.PATHS["worker_register"], worker, request_id=request_id)

    def heartbeat_worker(self, worker_key: str, heartbeat: Optional[Mapping[str, Any]] = None,
                         *, process_uuid: Optional[str] = None,
                         request_id: Optional[str] = None) -> Dict[str, Any]:
        payload = dict(heartbeat or {})
        if process_uuid is not None:
            payload["processUuid"] = _nonblank(process_uuid, "process_uuid")
        path = self.WORKER_HEARTBEAT_PATH.format(
            worker_key=self._path_segment(worker_key, "worker_key"))
        return self._post(path, payload, request_id=request_id)

    def list_force_handoff_requests(
            self, worker_key: str, *, process_uuid: str) -> Any:
        path = self.WORKER_HANDOFF_REQUESTS_PATH.format(
            worker_key=self._path_segment(worker_key, "worker_key"))
        query = urlencode({"processUuid": _nonblank(process_uuid, "process_uuid")})
        try:
            result = self._get(path + "?" + query)
        except ControlPlaneRejected as exc:
            # Rolling deployment compatibility: workers may update before the
            # control plane exposes handoff directives. An older 404 means there
            # cannot yet be a durable request to miss.
            if exc.status == 404:
                return []
            raise
        if not isinstance(result, list):
            raise InvalidResponse("force handoff requests response must be an array")
        return result

    def lease_task(self, worker_key: str, *, lease_seconds: int = 300,
                   capabilities: Optional[Any] = None,
                   process_uuid: Optional[str] = None,
                   request_id: Optional[str] = None) -> Dict[str, Any]:
        ttl = int(lease_seconds)
        if ttl <= 0:
            raise ValueError("lease_seconds must be positive")
        payload: Dict[str, Any] = {
            "workerKey": _nonblank(worker_key, "worker_key"),
            "leaseSeconds": ttl,
        }
        if capabilities is not None:
            payload["capabilities"] = capabilities
        if process_uuid is not None:
            payload["processUuid"] = _nonblank(process_uuid, "process_uuid")
        return self._post(self.PATHS["task_lease"], payload, request_id=request_id)

    def claim_task(self, worker_key: str, envelope: TaskEnvelope, *,
                   runtime_session_id: str, lease_seconds: int = 300,
                   free_slots: int = 1, process_uuid: Optional[str] = None,
                   request_id: Optional[str] = None) -> Dict[str, Any]:
        """Atomically attach ``worker_key`` to one explicitly named Task.

        Interactive workers use this endpoint instead of :meth:`lease_task`: the
        latter pulls the next matching Task, while this request carries the complete
        canonical envelope and can only claim that exact Task.  The envelope's
        ``requiredCapabilities.workerKey`` keeps targeted interactive work out of
        unrelated workers without introducing another Task mode.
        """
        if not isinstance(envelope, TaskEnvelope):
            raise TypeError("envelope must be a TaskEnvelope")
        ttl = int(lease_seconds)
        if ttl <= 0:
            raise ValueError("lease_seconds must be positive")
        slots = int(free_slots)
        if slots < 0:
            raise ValueError("free_slots must not be negative")
        payload = {
            "workerKey": _nonblank(worker_key, "worker_key"),
            "runtimeSessionId": _nonblank(runtime_session_id, "runtime_session_id"),
            "leaseSeconds": ttl,
            "freeSlots": slots,
            "task": envelope.to_dict(),
        }
        if process_uuid is not None:
            payload["processUuid"] = _nonblank(process_uuid, "process_uuid")
        rid = request_id or envelope.request_id("direct-claim")
        return self._post(self.PATHS["task_claim"], payload, request_id=rid)

    @staticmethod
    def _session_payload(worker_key: str, fence_token: Any,
                         detail: Optional[Mapping[str, Any]],
                         process_uuid: Optional[str] = None) -> Dict[str, Any]:
        payload = dict(detail or {})
        payload.update({
            "workerKey": _nonblank(worker_key, "worker_key"),
            "fenceToken": fence_token,
        })
        if fence_token is None or str(fence_token).strip() == "":
            raise ValueError("fence_token must not be empty")
        if process_uuid is not None:
            payload["processUuid"] = _nonblank(process_uuid, "process_uuid")
        return payload

    def _session_transition(self, action: str, session_id: str, worker_key: str,
                            fence_token: Any, detail: Optional[Mapping[str, Any]],
                            request_id: Optional[str],
                            process_uuid: Optional[str] = None) -> Dict[str, Any]:
        path = self.SESSION_ACTION_PATH.format(
            session_id=self._path_segment(session_id, "session_id"),
            action=self._path_segment(action, "action"),
        )
        payload = self._session_payload(worker_key, fence_token, detail, process_uuid)
        return self._post(path, payload, request_id=request_id)

    def start_session(self, session_id: str, worker_key: str, fence_token: Any,
                      detail: Optional[Mapping[str, Any]] = None, *,
                      process_uuid: Optional[str] = None,
                      request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._session_transition("start", session_id, worker_key, fence_token,
                                        detail, request_id, process_uuid)

    def heartbeat_session(self, session_id: str, worker_key: str, fence_token: Any,
                          detail: Optional[Mapping[str, Any]] = None, *,
                          process_uuid: Optional[str] = None,
                          request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._session_transition("heartbeat", session_id, worker_key, fence_token,
                                        detail, request_id, process_uuid)

    def acknowledge_force_handoff(
            self, session_id: str, worker_key: str, fence_token: Any, *,
            process_uuid: Optional[str] = None,
            request_id: Optional[str] = None) -> Dict[str, Any]:
        """Confirm that the fenced local process has stopped.

        The control plane keeps a forced handoff non-dispatchable until this ACK
        arrives. ``fence_token`` is the worker's last owned token; the server
        validates that the pending handoff advanced it exactly once.
        """
        if fence_token is None or str(fence_token).strip() == "":
            raise ValueError("fence_token must not be empty")
        path = self.SESSION_ACTION_PATH.format(
            session_id=self._path_segment(session_id, "session_id"),
            action="handoff-ack")
        payload = {
            "workerKey": _nonblank(worker_key, "worker_key"),
            "oldFenceToken": fence_token,
        }
        if process_uuid is not None:
            payload["processUuid"] = _nonblank(process_uuid, "process_uuid")
        return self._post(path, payload, request_id=request_id)

    def suspend_session(self, session_id: str, worker_key: str, fence_token: Any,
                        detail: Optional[Mapping[str, Any]] = None, *,
                        process_uuid: Optional[str] = None,
                        request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._session_transition("suspend", session_id, worker_key, fence_token,
                                        detail, request_id, process_uuid)

    def complete_session(self, session_id: str, worker_key: str, fence_token: Any,
                         detail: Optional[Mapping[str, Any]] = None, *,
                         process_uuid: Optional[str] = None,
                         request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._session_transition("complete", session_id, worker_key, fence_token,
                                        detail, request_id, process_uuid)

    def fail_session(self, session_id: str, worker_key: str, fence_token: Any,
                     detail: Optional[Mapping[str, Any]] = None, *,
                     process_uuid: Optional[str] = None,
                     request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._session_transition("fail", session_id, worker_key, fence_token,
                                        detail, request_id, process_uuid)

    def relinquish_session(self, session_id: str, worker_key: str, fence_token: Any,
                           detail: Optional[Mapping[str, Any]] = None, *,
                           process_uuid: Optional[str] = None,
                           request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._session_transition("relinquish", session_id, worker_key, fence_token,
                                        detail, request_id, process_uuid)

    def begin_operation(self, operation: Mapping[str, Any], *,
                        request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(self.PATHS["operation_begin"], operation,
                          request_id=request_id)

    def ack_operation(self, operation: Mapping[str, Any], *,
                      request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(self.PATHS["operation_ack"], operation,
                          request_id=request_id)

    def fail_operation(self, operation: Mapping[str, Any], *,
                       request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(self.PATHS["operation_fail"], operation,
                          request_id=request_id)

    def mark_operation_not_started(self, operation: Mapping[str, Any], *,
                                   request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(self.PATHS["operation_not_started"], operation,
                          request_id=request_id)

    def reconcile_operation(self, operation: Mapping[str, Any], *,
                            request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(self.PATHS["operation_reconcile"], operation,
                          request_id=request_id)

    def get_operation(self, operation_id: Any) -> Any:
        path = self.OPERATION_PATH.format(
            operation_id=self._path_segment(operation_id, "operation_id"))
        return self._get(path)

    def get_operation_by_key(self, task_id: Any, generation: Any,
                             operation_key: str) -> Any:
        if task_id is None or str(task_id).strip() == "":
            raise ValueError("task_id must not be empty")
        if generation is None or str(generation).strip() == "":
            raise ValueError("generation must not be empty")
        query = urlencode({
            "taskId": task_id,
            "generation": generation,
            "operationKey": _nonblank(operation_key, "operation_key"),
        })
        return self._get("%s?%s" % (self.OPERATION_BY_KEY_PATH, query))

    def list_external_operation_recovery_candidates(
            self, *, after_operation_id: int = 0, limit: int = 100) -> Any:
        after = int(after_operation_id)
        page_size = int(limit)
        if after < 0:
            raise ValueError("after_operation_id must not be negative")
        if page_size <= 0 or page_size > 500:
            raise ValueError("limit must be between 1 and 500")
        return self._get("%s?%s" % (
            self.EXTERNAL_OPERATION_RECOVERY_CANDIDATES_PATH,
            urlencode({"afterOperationId": after, "limit": page_size})))

    @staticmethod
    def _operation_recovery_payload(operation_id: Any, worker_key: str,
                                    recovery_token: str) -> Dict[str, Any]:
        if operation_id is None or str(operation_id).strip() == "":
            raise ValueError("operation_id must not be empty")
        return {
            "operationId": operation_id,
            "workerKey": _nonblank(worker_key, "worker_key"),
            "recoveryToken": _nonblank(recovery_token, "recovery_token"),
        }

    def lease_operation_recovery(self, operation_id: Any, worker_key: str,
                                 recovery_token: str, *,
                                 request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(
            self.OPERATION_RECOVERY_LEASE_PATH,
            self._operation_recovery_payload(operation_id, worker_key, recovery_token),
            request_id=request_id)

    def renew_operation_recovery(self, operation_id: Any, worker_key: str,
                                 recovery_token: str, *,
                                 request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(
            self.OPERATION_RECOVERY_RENEW_PATH,
            self._operation_recovery_payload(operation_id, worker_key, recovery_token),
            request_id=request_id)

    def release_operation_recovery(self, operation_id: Any, worker_key: str,
                                   recovery_token: str, *,
                                   request_id: Optional[str] = None) -> Dict[str, Any]:
        return self._post(
            self.OPERATION_RECOVERY_RELEASE_PATH,
            self._operation_recovery_payload(operation_id, worker_key, recovery_token),
            request_id=request_id)

    def get_task_by_aone(self, aone_id: str) -> Any:
        path = self.TASK_BY_AONE_PATH.format(
            aone_id=self._path_segment(aone_id, "aone_id"))
        return self._get(path)

    def list_pending_aone_reply_waits(self, *, after_session_id: int = 0,
                                      limit: int = 100) -> Any:
        after = int(after_session_id)
        page_size = int(limit)
        if after < 0:
            raise ValueError("after_session_id must not be negative")
        if page_size <= 0 or page_size > 500:
            raise ValueError("limit must be between 1 and 500")
        query = urlencode({"afterSessionId": after, "limit": page_size})
        return self._get("%s?%s" % (self.PENDING_AONE_WAITS_PATH, query))

    def list_workers(self) -> Any:
        """Return every registered worker and its persisted heartbeat state."""
        return self._get(self.WORKERS_PATH)

    def list_ready_task_diagnostics(self, *, limit: int = 100) -> Any:
        """Return READY tasks with the reason each one can or cannot be leased."""
        page_size = int(limit)
        if page_size <= 0 or page_size > 500:
            raise ValueError("limit must be between 1 and 500")
        return self._get("%s?%s" % (
            self.READY_TASK_DIAGNOSTICS_PATH,
            urlencode({"limit": page_size})))

    def get_worker_state(self, worker_key: str) -> Any:
        """Return one worker plus its current task/session assignments."""
        path = self.WORKER_STATE_PATH.format(
            worker_key=self._path_segment(worker_key, "worker_key"))
        return self._get(path)

    def get_task_timeline(self, task_id: str) -> Any:
        path = self.TASK_TIMELINE_PATH.format(
            task_id=self._path_segment(task_id, "task_id"))
        return self._get(path)

    def discard_resume_context(self, task_id: str, expected_session_id: int, reason: str,
                               *, request_id: Optional[str] = None) -> Dict[str, Any]:
        """Explicitly discard one reviewed legacy resume context."""
        path = self.DISCARD_RESUME_CONTEXT_PATH.format(
            task_id=self._path_segment(task_id, "task_id"))
        reason = _nonblank(reason, "reason")
        return self._post(path, {
            "expectedSessionId": int(expected_session_id),
            "reason": reason,
        }, request_id=request_id)

    def list_source_status_candidates(self, *, after_task_id: int = 0,
                                      limit: int = 100) -> Any:
        """Page through Aone-backed Tasks independently from dispatch eligibility."""
        after = int(after_task_id)
        page_size = int(limit)
        if after < 0:
            raise ValueError("after_task_id must not be negative")
        if page_size <= 0 or page_size > 500:
            raise ValueError("limit must be between 1 and 500")
        return self._get("%s?%s" % (
            self.SOURCE_STATUS_CANDIDATES_PATH,
            urlencode({"afterTaskId": after, "limit": page_size})))

    def update_source_status(self, task_id: str, aone_id: str, source_status: str,
                             *, request_id: Optional[str] = None) -> Dict[str, Any]:
        """Report source lifecycle metadata without upserting or waking a Task."""
        path = self.SOURCE_STATUS_PATH.format(
            task_id=self._path_segment(task_id, "task_id"))
        aone_id = _nonblank(aone_id, "aone_id")
        source_status = _nonblank(source_status, "source_status")
        return self._post(path, {
            "aoneId": aone_id,
            "sourceStatus": source_status,
        }, request_id=request_id)
