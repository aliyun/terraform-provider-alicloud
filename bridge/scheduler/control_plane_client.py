"""HTTP adapter for the scheduled-job control-plane contract.

The scheduler core depends only on :class:`ScheduledJobControlPlane`; this
adapter is the composition-layer implementation for AutomationAgent's JSON
API.  Importing it has no network or environment side effects.  Requests are
made only by explicit lifecycle calls and every uncertain response is rejected
so an Engine never treats an unacknowledged state transition as durable.
"""

from __future__ import annotations

from datetime import datetime, timezone
import json
import time
from typing import Any, Callable, Mapping, Optional, Sequence
from urllib.error import HTTPError, URLError
from urllib.parse import quote
from urllib.request import Request, urlopen

from .engine import JobRegistration, ScheduledJobState, ScheduledJobStatus
from .model import is_aware


class ScheduledJobControlPlaneError(RuntimeError):
    """Base error for a scheduled-job control-plane operation that did not commit."""

    terminal_retryable = False

    def __init__(self, message: str, *, status: Optional[int] = None) -> None:
        super().__init__(message)
        self.status = status


class ScheduledJobControlPlaneUnavailable(ScheduledJobControlPlaneError):
    """The control plane cannot be safely reached or returned a 5xx response."""

    terminal_retryable = True


class ScheduledJobControlPlaneRejected(ScheduledJobControlPlaneError):
    """The server rejected a request or returned a malformed successful response."""


class ScheduledJobControlPlaneProtocolError(ScheduledJobControlPlaneRejected):
    """A 2xx response did not conform to the scheduled-job API contract."""

    # A well-formed 2xx response that violates the contract (e.g. a missing or
    # non-boolean ``admitted``) is a definitive server-side error, not transport
    # uncertainty: retrying cannot resolve the ambiguity and ``start`` must fail
    # closed rather than risk a second admission attempt.  Only genuine transport
    # loss (``ScheduledJobControlPlaneUnavailable``) is retryable.
    terminal_retryable = False


def _nonblank(value: Any, name: str) -> str:
    text = str(value or "").strip()
    if not text:
        raise ValueError(f"{name} must not be empty")
    return text


def _require_aware(value: datetime, name: str) -> datetime:
    if not is_aware(value):
        raise ValueError(f"{name} must be timezone-aware")
    return value


def _to_utc_timestamp(value: datetime, name: str) -> str:
    """Encode an aware instant in the Java API's ISO-8601 UTC representation."""

    return _require_aware(value, name).astimezone(timezone.utc).isoformat(
        timespec="milliseconds")


def _from_utc_timestamp(value: Any, name: str) -> datetime:
    if not isinstance(value, str) or not value.strip():
        raise ScheduledJobControlPlaneProtocolError(
            f"control plane response {name} must be an ISO-8601 timestamp")
    try:
        parsed = datetime.fromisoformat(value.strip().replace("Z", "+00:00"))
    except ValueError as exc:
        raise ScheduledJobControlPlaneProtocolError(
            f"control plane response {name} must be an ISO-8601 timestamp") from exc
    if not is_aware(parsed):
        raise ScheduledJobControlPlaneProtocolError(
            f"control plane response {name} must include a timezone")
    return parsed.astimezone(timezone.utc)


def _same_millisecond(left: datetime, right: datetime) -> bool:
    return int(left.timestamp() * 1000) == int(right.timestamp() * 1000)


class HttpScheduledJobControlPlane:
    """Strict standard-library HTTP implementation of ``ScheduledJobControlPlane``.

    Environment resolution belongs to the composition root; construction
    itself does no I/O.
    """

    DEFAULT_PATH = "/api/jarvis/v1/scheduled-jobs"

    def __init__(
        self,
        base_url: str,
        token: str = "",
        *,
        timeout: float = 10,
        retry_delay_seconds: float = 5,
        path: str = DEFAULT_PATH,
        opener: Optional[Callable[..., Any]] = None,
        worker_key: Optional[str] = None,
        process_uuid: Optional[str] = None,
    ) -> None:
        self.base_url = _nonblank(base_url, "base_url").rstrip("/")
        self.token = str(token or "")
        self.timeout = float(timeout)
        if self.timeout <= 0:
            raise ValueError("timeout must be positive")
        self.retry_delay_seconds = float(retry_delay_seconds)
        if self.retry_delay_seconds <= 0:
            raise ValueError("retry_delay_seconds must be positive")
        self.path = "/" + _nonblank(path, "path").strip("/")
        self._opener = opener or urlopen
        if (worker_key is None) != (process_uuid is None):
            raise ValueError("worker_key and process_uuid must be configured together")
        self.worker_key = None if worker_key is None else _nonblank(worker_key, "worker_key")
        self.process_uuid = (
            None if process_uuid is None else _nonblank(process_uuid, "process_uuid"))

    def register(self, registrations: Sequence[JobRegistration]) -> Sequence[ScheduledJobState]:
        payload = self._identity_payload({
            "jobs": [self._registration_payload(registration) for registration in registrations],
        })
        return self._states(self._request("POST", "/register", payload), "register")

    def list_jobs(self) -> Sequence[ScheduledJobState]:
        return self._states(self._request("GET", "", None), "list")

    def recover_interrupted(self) -> Sequence[ScheduledJobState]:
        """Recover only the server-reported interrupted slots, fail-closed otherwise."""

        response = self._object(
            self._request("POST", "/recover-interrupted", self._identity_payload({})),
            "recover-interrupted")
        recovered = response.get("recovered")
        if type(recovered) is not int or recovered < 0:
            raise ScheduledJobControlPlaneProtocolError(
                "control plane recover-interrupted response must include non-negative integer recovered")
        if "jobs" not in response:
            raise ScheduledJobControlPlaneProtocolError(
                "control plane recover-interrupted response must include jobs")
        jobs = self._states(response["jobs"], "recover-interrupted.jobs")
        if recovered != len(jobs):
            raise ScheduledJobControlPlaneProtocolError(
                "control plane recover-interrupted recovered count must equal jobs length")
        for job in jobs:
            if job.status is not ScheduledJobStatus.IDLE or job.next_run_at is None:
                raise ScheduledJobControlPlaneProtocolError(
                    "control plane recover-interrupted jobs must be IDLE with reserved nextRunAt")
        return jobs

    def start(
        self, job_key: str, scheduled_for: datetime, next_run_at: datetime,
    ) -> bool:
        payload = self._identity_payload({
            "scheduledFor": _to_utc_timestamp(scheduled_for, "scheduled_for"),
            "nextRunAt": _to_utc_timestamp(next_run_at, "next_run_at"),
        })
        uncertain = False
        while True:
            try:
                response = self._object(self._request(
                    "POST", self._job_path(job_key, "start"), payload), "start")
                admitted = response.get("admitted")
                if type(admitted) is not bool:
                    raise ScheduledJobControlPlaneProtocolError(
                        "control plane start response must include boolean admitted")
                state = self._state(response.get("job"), "start.job")
                if admitted:
                    return True
                # This call has not launched the runner yet. If an earlier
                # uncertain attempt reserved the exact successor, the retry's
                # non-admission is proof that this process owns the admission.
                if (uncertain
                        and state.status is ScheduledJobStatus.RUNNING
                        and state.next_run_at is not None
                        and _same_millisecond(state.next_run_at, next_run_at)):
                    return True
                return False
            except Exception as exc:
                if uncertain and self._reserved_running(job_key, next_run_at):
                    return True
                if not getattr(exc, "terminal_retryable", False):
                    raise
                uncertain = True
                time.sleep(self.retry_delay_seconds)

    def _reserved_running(self, job_key: str, next_run_at: datetime) -> bool:
        """Reconcile an uncertain start without authorizing a second invocation."""

        try:
            states = self.list_jobs()
        except Exception:
            return False
        return any(
            state.job_key == job_key
            and state.status is ScheduledJobStatus.RUNNING
            and state.next_run_at is not None
            and _same_millisecond(state.next_run_at, next_run_at)
            for state in states
        )

    def complete(self, job_key: str, scheduled_for: datetime, next_run_at: datetime) -> None:
        response = self._request(
            "POST", self._job_path(job_key, "complete"),
            self._identity_payload({
                "scheduledFor": _to_utc_timestamp(scheduled_for, "scheduled_for"),
                "nextRunAt": _to_utc_timestamp(next_run_at, "next_run_at"),
            }),
        )
        self._state(response, "complete")

    def fail(
        self,
        job_key: str,
        scheduled_for: datetime,
        *,
        retryable: bool,
        next_run_at: Optional[datetime],
        error: str,
    ) -> None:
        if not isinstance(retryable, bool):
            raise ValueError("retryable must be a bool")
        if retryable and next_run_at is None:
            raise ValueError("next_run_at is required for retryable failures")
        if not retryable and next_run_at is not None:
            raise ValueError("next_run_at must be None for permanent failures")
        payload = self._identity_payload({
            "scheduledFor": _to_utc_timestamp(scheduled_for, "scheduled_for"),
            "retryable": retryable,
            "nextRunAt": (
                _to_utc_timestamp(next_run_at, "next_run_at")
                if next_run_at is not None else None
            ),
            "error": _nonblank(error, "error"),
        })
        self._state(self._request("POST", self._job_path(job_key, "fail"), payload), "fail")

    def _registration_payload(self, registration: JobRegistration) -> dict[str, Any]:
        if not isinstance(registration, JobRegistration):
            raise TypeError("registrations must contain JobRegistration values")
        if registration.enabled and registration.next_run_at is None:
            raise ValueError("next_run_at is required for enabled job registrations")
        return {
            "jobKey": _nonblank(registration.job_key, "job_key"),
            "jobName": _nonblank(registration.job_name, "job_name"),
            "definition": dict(registration.definition),
            "nextRunAt": (
                _to_utc_timestamp(registration.next_run_at, "next_run_at")
                if registration.next_run_at is not None else None
            ),
            "enabled": bool(registration.enabled),
        }

    def _job_path(self, job_key: str, action: str) -> str:
        return "/%s/%s" % (quote(_nonblank(job_key, "job_key"), safe=""), action)

    def _endpoint(self, suffix: str) -> str:
        return self.base_url + self.path + suffix

    def _identity_payload(self, payload: Mapping[str, Any]) -> dict[str, Any]:
        """Attach the fixed Worker fence to every scheduled-job mutation.

        The headers below fence GET/list as well.  Keeping the identity in the
        mutation body makes the contract visible to the Java API and avoids a
        proxy dropping non-standard headers; the server must require both and
        reject a mismatch.
        """

        result = dict(payload)
        if self.worker_key is not None:
            result["workerKey"] = self.worker_key
            result["processUuid"] = self.process_uuid
        return result

    def _request(self, method: str, suffix: str, payload: Optional[Mapping[str, Any]]) -> Any:
        body = None
        headers = {
            "Accept": "application/json",
            "User-Agent": "jarvis-scheduled-job-control-plane/1",
        }
        if self.worker_key is not None:
            headers["X-Jarvis-Worker-Key"] = self.worker_key
            headers["X-Jarvis-Process-Uuid"] = str(self.process_uuid)
        if payload:
            body = json.dumps(payload, ensure_ascii=False, separators=(",", ":")).encode("utf-8")
            headers["Content-Type"] = "application/json"
        if self.token:
            headers["Authorization"] = "Bearer " + self.token
        request = Request(self._endpoint(suffix), data=body, method=method, headers=headers)
        try:
            with self._opener(request, timeout=self.timeout) as response:
                status_value = getattr(response, "status", None)
                if status_value is None:
                    status_value = response.getcode()
                try:
                    status = int(status_value)
                except (TypeError, ValueError) as exc:
                    raise ScheduledJobControlPlaneProtocolError(
                        "scheduled-job control plane returned an invalid HTTP status") from exc
                raw = response.read() or b""
        except HTTPError as exc:
            self._raise_http(int(exc.code), "rejected")
        except (URLError, TimeoutError, OSError) as exc:
            raise ScheduledJobControlPlaneUnavailable(
                f"scheduled-job control plane unavailable: {type(exc).__name__}") from exc
        if status < 200 or status >= 300:
            self._raise_http(status, "rejected")
        if not raw:
            raise ScheduledJobControlPlaneProtocolError(
                "scheduled-job control plane returned an empty response", status=status)
        try:
            return json.loads(raw.decode("utf-8", errors="strict"))
        except (UnicodeDecodeError, TypeError, ValueError) as exc:
            raise ScheduledJobControlPlaneProtocolError(
                "scheduled-job control plane returned invalid JSON", status=status) from exc

    @staticmethod
    def _raise_http(status: int, detail: str) -> None:
        message = f"scheduled-job control plane HTTP {status}: {detail}"
        if status >= 500:
            raise ScheduledJobControlPlaneUnavailable(message, status=status)
        raise ScheduledJobControlPlaneRejected(message, status=status)

    @staticmethod
    def _object(value: Any, name: str) -> Mapping[str, Any]:
        if not isinstance(value, Mapping):
            raise ScheduledJobControlPlaneProtocolError(
                f"control plane {name} response must be an object")
        return value

    def _states(self, value: Any, name: str) -> tuple[ScheduledJobState, ...]:
        if not isinstance(value, list):
            raise ScheduledJobControlPlaneProtocolError(
                f"control plane {name} response must be an array")
        return tuple(self._state(item, f"{name}[{index}]") for index, item in enumerate(value))

    def _state(self, value: Any, name: str) -> ScheduledJobState:
        item = self._object(value, name)
        job_key = item.get("jobKey")
        status = item.get("status")
        if not isinstance(job_key, str) or not job_key.strip():
            raise ScheduledJobControlPlaneProtocolError(
                f"control plane {name}.jobKey must be a nonblank string")
        try:
            parsed_status = ScheduledJobStatus(status)
        except (TypeError, ValueError) as exc:
            raise ScheduledJobControlPlaneProtocolError(
                f"control plane {name}.status is invalid") from exc
        if "nextRunAt" not in item:
            raise ScheduledJobControlPlaneProtocolError(
                f"control plane {name}.nextRunAt is required")
        next_run_at = item["nextRunAt"]
        return ScheduledJobState(
            job_key.strip(), parsed_status,
            None if next_run_at is None else _from_utc_timestamp(next_run_at, f"{name}.nextRunAt"),
        )
