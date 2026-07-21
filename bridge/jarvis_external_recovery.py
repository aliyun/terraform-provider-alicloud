#!/usr/bin/env python3
"""Control-plane driven recovery for orphaned Aone operation receipts."""

from __future__ import annotations

import hashlib
import json
import logging
import os
import subprocess
import threading
import time
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable, Mapping, Optional, Tuple


class RecoveryInconclusive(RuntimeError):
    """Authoritative Aone readback was unavailable or malformed."""


def _normalized_digest(text: Any) -> str:
    value = str(text or "").replace("\r", "")
    value = "\n".join(line.rstrip() for line in value.split("\n")).rstrip()
    return hashlib.sha256(value.encode("utf-8")).hexdigest()


def _timestamp(value: Any) -> float:
    if isinstance(value, (int, float)) and not isinstance(value, bool):
        number = float(value)
        return number / 1000.0 if abs(number) >= 100_000_000_000 else number
    text = str(value or "").strip()
    if not text:
        raise ValueError("timestamp is empty")
    try:
        number = float(text)
    except ValueError:
        parsed = datetime.fromisoformat(text[:-1] + "+00:00" if text.endswith("Z") else text)
        if parsed.tzinfo is None:
            parsed = parsed.replace(tzinfo=timezone.utc)
        return parsed.timestamp()
    return number / 1000.0 if abs(number) >= 100_000_000_000 else number


def _field_value(document: Mapping[str, Any], identifier: str) -> Any:
    fields = document.get("fields") or []
    if isinstance(fields, Mapping):
        fields = list(fields.values())
    for field in fields if isinstance(fields, list) else []:
        if not isinstance(field, Mapping):
            continue
        if str(field.get("identifier") or field.get("fieldIdentifier") or "") != identifier:
            continue
        return field.get("displayValue") if field.get("displayValue") is not None else field.get("value")
    return None


def _display_text(value: Any) -> str:
    if isinstance(value, Mapping):
        value = value.get("displayValue") or value.get("name") or value.get("value")
    return str(value or "").strip()


def _tag_names(value: Any) -> set[str]:
    if isinstance(value, str):
        return {part.strip() for part in value.split(",") if part.strip()}
    if isinstance(value, Mapping):
        nested = value.get("displayValue") or value.get("name") or value.get("value")
        return _tag_names(nested)
    if isinstance(value, (list, tuple, set)):
        names = set()
        for item in value:
            names.update(_tag_names(item))
        return names
    text = str(value or "").strip()
    return {text} if text else set()


def _ticket_status_and_tags(document: Mapping[str, Any]) -> Tuple[str, set[str]]:
    status = _display_text(_field_value(document, "status"))
    if not status:
        value = document.get("status") or document.get("statusName") or ""
        if isinstance(value, Mapping):
            value = value.get("displayValue") or value.get("name") or value.get("value")
        status = str(value or "").strip()
    raw_tags = (_field_value(document, "tag") or document.get("tag")
                or document.get("tags") or document.get("labels") or [])
    return status, _tag_names(raw_tags)


class ExternalOperationRecoveryScheduler:
    """Reconcile externally visible Aone effects without replaying writes.

    The scheduler is started only by the bridge scheduler role.  Every candidate is
    protected by the control-plane recovery token.  Unavailable or ambiguous Aone
    reads release that token and leave the receipt UNKNOWN.
    """

    def __init__(self, client: Any, worker_key: str, *, repo_root: Any,
                 interval: Optional[float] = None, page_size: int = 100,
                 max_pages: int = 10, enabled: Optional[bool] = None,
                 observer: Optional[Callable[[Mapping[str, Any]], Tuple[bool, Optional[str]]]] = None,
                 logger: Optional[logging.Logger] = None):
        self.client = client
        self.worker_key = str(worker_key or "").strip()
        if not self.worker_key:
            raise ValueError("worker_key must not be empty")
        self.repo_root = Path(repo_root)
        self.interval = float(interval if interval is not None else os.environ.get(
            "JARVIS_EXTERNAL_RECOVERY_INTERVAL", "300"))
        self.page_size = max(1, min(int(page_size), 500))
        self.max_pages = max(1, int(max_pages))
        self.enabled = (str(os.environ.get("JARVIS_EXTERNAL_RECOVERY_ENABLE", "1")) == "1"
                        if enabled is None else bool(enabled))
        self.observer = observer or self._observe_aone
        self.log = logger or logging.getLogger(__name__)
        self._stop = threading.Event()
        self._thread = None
        self._after_operation_id = 0

    def start(self) -> None:
        if not self.enabled:
            self.log.info("ExternalOperationRecoveryScheduler disabled")
            return
        if self._thread and self._thread.is_alive():
            return
        self._thread = threading.Thread(
            target=self._loop, name="external-operation-recovery", daemon=True)
        self._thread.start()

    def stop(self) -> None:
        self._stop.set()
        thread = self._thread
        if thread and thread.is_alive() and thread is not threading.current_thread():
            thread.join(timeout=min(max(self.interval, 1.0), 10.0))

    def _loop(self) -> None:
        while not self._stop.is_set():
            try:
                self._tick()
            except Exception:  # noqa: BLE001 - daemon must survive transient control-plane faults
                self.log.exception("external operation recovery tick failed")
            self._stop.wait(max(1.0, self.interval))

    def _tick(self) -> None:
        after = self._after_operation_id
        for _page_number in range(self.max_pages):
            if self._stop.is_set():
                return
            page = self.client.list_external_operation_recovery_candidates(
                after_operation_id=after, limit=self.page_size)
            if not isinstance(page, Mapping) or not isinstance(page.get("items"), list):
                raise ValueError("control plane external recovery page is invalid")
            for candidate in page["items"]:
                if self._stop.is_set():
                    return
                if isinstance(candidate, Mapping):
                    self._recover(candidate)
            if not page.get("hasMore"):
                self._after_operation_id = 0
                return
            try:
                next_after = int(page.get("nextAfterOperationId"))
            except (TypeError, ValueError) as exc:
                raise ValueError("control plane external recovery cursor is invalid") from exc
            if next_after <= after:
                raise ValueError("control plane external recovery cursor did not advance")
            after = next_after
            self._after_operation_id = after
        self.log.warning("external operation recovery reached max_pages=%d", self.max_pages)

    def _recover(self, candidate: Mapping[str, Any]) -> None:
        operation = candidate.get("operation") or {}
        operation_id = operation.get("id") if isinstance(operation, Mapping) else None
        if operation_id is None:
            self.log.warning("external recovery candidate missing operation id")
            return
        if self._stop.is_set():
            return
        token = hashlib.sha256(
            (self.worker_key + ":" + str(operation_id)).encode("utf-8")).hexdigest()
        leased = False
        try:
            lease = self.client.lease_operation_recovery(
                operation_id, self.worker_key, token,
                request_id="external-recovery-lease:%s:%s" % (operation_id, token))
            if not isinstance(lease, Mapping) or not lease.get("proceed"):
                return
            leased = True
            if self._stop.is_set():
                raise RecoveryInconclusive("scheduler is stopping")
            self._renew(operation_id, token)
            found, external_ref = self.observer(candidate)
            if not isinstance(found, bool):
                raise RecoveryInconclusive("observer did not return a definite verdict")
            if self._stop.is_set():
                raise RecoveryInconclusive("scheduler is stopping")
            self._renew(operation_id, token)
            self.client.reconcile_operation({
                "operationId": operation_id,
                "workerKey": self.worker_key,
                "found": found,
                "externalRef": external_ref if found else None,
                "retryAllowed": not found,
                "retryAfterSeconds": 30 if not found else None,
                "recoveryToken": token,
            }, request_id="external-recovery-reconcile:%s:%s" % (operation_id, token))
            leased = False
            self.log.info("external operation reconciled operation=%s found=%s",
                          operation_id, found)
        except RecoveryInconclusive as exc:
            self.log.warning("external operation readback inconclusive operation=%s: %s",
                             operation_id, exc)
        except Exception as exc:  # noqa: BLE001 - fail closed and release below
            self.log.warning("external operation recovery failed operation=%s: %s",
                             operation_id, type(exc).__name__)
        finally:
            if leased:
                try:
                    self.client.release_operation_recovery(
                        operation_id, self.worker_key, token,
                        request_id="external-recovery-release:%s:%s" % (operation_id, token))
                except Exception as exc:  # noqa: BLE001 - lease expiry remains the final fallback
                    self.log.warning("external recovery release failed operation=%s: %s",
                                     operation_id, type(exc).__name__)

    def _renew(self, operation_id: Any, token: str) -> None:
        renewed = self.client.renew_operation_recovery(
            operation_id, self.worker_key, token,
            request_id="external-recovery-renew:%s:%s:%s" % (
                operation_id, token, time.monotonic_ns()))
        if not isinstance(renewed, Mapping) or not renewed.get("proceed"):
            raise RecoveryInconclusive("recovery lease was not renewed")

    def _a1_json(self, args: list[str]) -> Any:
        command = [str(self.repo_root / "bin" / "a1id"), "--"] + args + ["-f", "json"]
        try:
            result = subprocess.run(
                command, cwd=str(self.repo_root), capture_output=True, text=True,
                timeout=float(os.environ.get("JARVIS_EXTERNAL_RECOVERY_AONE_TIMEOUT", "30")))
        except (OSError, subprocess.TimeoutExpired) as exc:
            raise RecoveryInconclusive(type(exc).__name__) from exc
        if result.returncode != 0:
            raise RecoveryInconclusive("a1 read failed rc=%d" % result.returncode)
        try:
            return json.loads(result.stdout)
        except (TypeError, ValueError) as exc:
            raise RecoveryInconclusive("a1 returned invalid JSON") from exc

    def _observe_aone(self, candidate: Mapping[str, Any]) -> Tuple[bool, Optional[str]]:
        task = candidate.get("task") or {}
        operation = candidate.get("operation") or {}
        payload = candidate.get("readbackSpec") or {}
        if not all(isinstance(value, Mapping) for value in (task, operation, payload)):
            raise RecoveryInconclusive("candidate shape is invalid")
        aone_id = str(task.get("aoneId") or "").strip()
        target = str(operation.get("target") or "").removeprefix("aone:").strip()
        if not aone_id or target != aone_id:
            raise RecoveryInconclusive("candidate Aone identity mismatch")
        operation_type = str(operation.get("operationType") or "")

        if operation_type == "AONE_COMMENT":
            digest = str(payload.get("digest") or "").strip()
            if len(digest) != 64:
                raise RecoveryInconclusive("comment digest is invalid")
            try:
                receipt_created_at = _timestamp(operation.get("createdAt"))
            except (TypeError, ValueError, OverflowError) as exc:
                raise RecoveryInconclusive("candidate createdAt is invalid") from exc
            comments = self._a1_json(
                ["project", "workitem", "comment", "list", aone_id])
            if not isinstance(comments, list):
                raise RecoveryInconclusive("comment readback is not a complete list")
            for comment in comments:
                if not isinstance(comment, Mapping):
                    continue
                if _normalized_digest(comment.get("content")) == digest:
                    try:
                        comment_created_at = _timestamp(comment.get("createdAt"))
                    except (TypeError, ValueError, OverflowError) as exc:
                        raise RecoveryInconclusive(
                            "matched comment createdAt is invalid") from exc
                    if comment_created_at < receipt_created_at:
                        continue
                    comment_id = str(comment.get("id") or "").strip()
                    if not comment_id:
                        raise RecoveryInconclusive("matched comment has no id")
                    return True, "aone:%s:comment:%s" % (aone_id, comment_id)
            return False, None

        document = self._a1_json(["project", "workitem", "get", aone_id])
        if not isinstance(document, Mapping):
            raise RecoveryInconclusive("workitem readback is not an object")
        status, tags = _ticket_status_and_tags(document)
        if operation_type == "AONE_STATUS":
            expected = str(payload.get("material") or "").strip()
            if payload.get("kind") != "status" or not expected or not status:
                raise RecoveryInconclusive("status candidate/readback is invalid")
            if status != expected:
                raise RecoveryInconclusive("current status does not prove the write was absent")
            return True, "aone:%s:status:%s" % (aone_id, expected)
        if operation_type == "AONE_CLAIM":
            added = _tag_names(payload.get("addTag"))
            if len(added) != 1:
                raise RecoveryInconclusive("claim candidate addTag is invalid")
            expected = next(iter(added))
            removed = _tag_names(payload.get("removeTag"))
        elif operation_type == "AONE_RELEASE" and payload.get("kind") == "release-tag":
            expected = "jarvis-idle"
            removed = {"jarvis-claimed"}
        elif operation_type == "AONE_RELEASE" and payload.get("kind") == "finish-tag":
            expected = "jarvis-done"
            removed = {"jarvis-claimed", "jarvis-idle"}
        else:
            raise RecoveryInconclusive("unsupported external operation candidate")
        if not expected:
            raise RecoveryInconclusive("tag candidate has no expected tag")
        found = expected in tags and tags.isdisjoint(removed)
        if not found:
            raise RecoveryInconclusive("current tags do not prove the write was absent")
        return True, "aone:%s:tag:%s" % (aone_id, expected)
