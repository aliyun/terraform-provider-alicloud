"""Cross-process registry for tickets awaiting supervised dispatch approval."""

from __future__ import annotations

import fcntl
import json
import os
from pathlib import Path
import threading
import time
from typing import Any, Mapping


DEFAULT_PATH = (
    Path(__file__).resolve().parents[1]
    / ".my-day"
    / "bridge"
    / "pending-dispatch.json"
)


class PendingDispatchRegistry:
    """Small locked JSON registry shared by Scheduler and DingTalk Bot."""

    def __init__(self, path: str | os.PathLike[str] | None = None) -> None:
        configured = path or os.environ.get("JARVIS_PENDING_DISPATCH_PATH")
        self.path = Path(configured) if configured else DEFAULT_PATH
        self.lock_path = self.path.with_name(self.path.name + ".lock")

    def _read(self) -> dict[str, Any]:
        if not self.path.exists():
            return {"version": 1, "pending": {}}
        raw = json.loads(self.path.read_text(encoding="utf-8"))
        pending = raw.get("pending") if isinstance(raw, dict) else None
        if not isinstance(pending, dict):
            raise ValueError("pending dispatch registry is invalid")
        return {"version": 1, "pending": pending}

    def _write(self, value: Mapping[str, Any]) -> None:
        self.path.parent.mkdir(parents=True, exist_ok=True)
        tmp = self.path.with_name(
            ".%s.%s.%s.tmp"
            % (self.path.name, os.getpid(), threading.get_ident()))
        try:
            tmp.write_text(
                json.dumps(value, ensure_ascii=False, sort_keys=True),
                encoding="utf-8",
            )
            os.chmod(tmp, 0o600)
            os.replace(tmp, self.path)
        finally:
            tmp.unlink(missing_ok=True)

    def _locked(self, callback):
        self.lock_path.parent.mkdir(parents=True, exist_ok=True)
        with self.lock_path.open("a+", encoding="utf-8") as lock_file:
            os.chmod(self.lock_path, 0o600)
            fcntl.flock(lock_file.fileno(), fcntl.LOCK_EX)
            try:
                return callback()
            finally:
                fcntl.flock(lock_file.fileno(), fcntl.LOCK_UN)

    def stage(
        self,
        item: Mapping[str, Any],
        dispatch_context: Mapping[str, Any],
        *,
        force: bool = False,
    ) -> bool:
        """Persist a new approval candidate; existing candidates stay unchanged."""

        item_id = str(item.get("id") or "").strip()
        revision = str(dispatch_context.get("revision") or "").strip()
        if not item_id.isdigit() or not revision:
            raise ValueError("pending dispatch requires numeric id and revision")

        def update() -> bool:
            value = self._read()
            if item_id in value["pending"]:
                return False
            value["pending"][item_id] = {
                "item": dict(item),
                "dispatchContext": dict(dispatch_context),
                "force": bool(force),
                "stagedAt": time.strftime(
                    "%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
            }
            self._write(value)
            return True

        return bool(self._locked(update))

    def get(self, item_id: object) -> dict[str, Any] | None:
        def read():
            record = self._read()["pending"].get(str(item_id))
            return dict(record) if isinstance(record, dict) else None

        return self._locked(read)

    def list(self) -> list[dict[str, Any]]:
        def read():
            records = self._read()["pending"]
            return [
                dict(records[item_id])
                for item_id in sorted(
                    records, key=lambda value: (
                        0, int(value)) if str(value).isdigit()
                    else (1, str(value)))
                if isinstance(records[item_id], dict)
            ]

        return self._locked(read)

    def clear(self) -> None:
        """Discard stale approvals when automatic dispatch is actually active."""

        self._locked(
            lambda: self._write({"version": 1, "pending": {}}))

    def remove(self, item_id: object, *, expected_revision: str) -> bool:
        """Remove only the generation that was successfully enqueued."""

        item_id = str(item_id)

        def update() -> bool:
            value = self._read()
            record = value["pending"].get(item_id)
            context = record.get("dispatchContext") if isinstance(record, dict) else None
            if (not isinstance(context, dict)
                    or str(context.get("revision") or "") != expected_revision):
                return False
            value["pending"].pop(item_id, None)
            self._write(value)
            return True

        return bool(self._locked(update))
