"""Durable PR-watch registry shared by inbound commands and the Scheduler."""

from __future__ import annotations

import json
import logging
import os
from pathlib import Path
import threading
import time

from bridge.aone_workitems import REPO_ROOT


log = logging.getLogger("jarvis-pr-watch")
PRWATCH_PATH = Path(REPO_ROOT) / ".my-day/bridge/pr-watch.json"
_prwatch_lock = threading.Lock()


def _prwatch_load():
    """Load the shared registry; a missing/corrupt file is an empty view."""
    try:
        raw = json.loads(PRWATCH_PATH.read_text())
        return ({str(key): dict(value) for key, value in raw.items()
                 if isinstance(value, dict)} if isinstance(raw, dict) else {})
    except FileNotFoundError:
        return {}
    except Exception as exc:  # noqa: BLE001
        log.warning("prwatch: could not load %s: %s", PRWATCH_PATH, exc)
        return {}


def _prwatch_has(ticket):
    with _prwatch_lock:
        return str(ticket) in _prwatch_load()


def _prwatch_write(records):
    try:
        PRWATCH_PATH.parent.mkdir(parents=True, exist_ok=True)
        temporary = PRWATCH_PATH.with_name(PRWATCH_PATH.name + ".tmp")
        temporary.write_text(json.dumps(records, ensure_ascii=False, sort_keys=True))
        os.replace(str(temporary), str(PRWATCH_PATH))
        return True
    except Exception as exc:  # noqa: BLE001
        log.warning("prwatch: could not persist %s: %s", PRWATCH_PATH, exc)
        return False


def _prwatch_acquire_file_lock():
    lock_path = PRWATCH_PATH.parent / ".pr-watch.lock"
    deadline = time.time() + 5
    while time.time() < deadline:
        try:
            PRWATCH_PATH.parent.mkdir(parents=True, exist_ok=True)
            lock_path.mkdir()
            return lock_path
        except FileExistsError:
            try:
                if time.time() - lock_path.stat().st_mtime > 10:
                    lock_path.rmdir()
                    continue
            except (FileNotFoundError, OSError):
                pass
            time.sleep(0.1)
        except OSError:
            break
    log.warning("prwatch: file lock busy; continuing with atomic best-effort write")
    return None


def _prwatch_release_file_lock(lock_path):
    if lock_path is not None:
        try:
            lock_path.rmdir()
        except OSError:
            pass


def _prwatch_add(ticket, pr_url, project, title=""):
    with _prwatch_lock:
        file_lock = _prwatch_acquire_file_lock()
        try:
            records = _prwatch_load()
            existing = records.get(str(ticket))
            entry = dict(existing) if isinstance(existing, dict) else {}
            entry.update({
                "pr_url": pr_url,
                "project": project,
                "title": str(entry.get("title") or title or "").strip(),
                "submitted_at": entry.get("submitted_at")
                or time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
            })
            records[str(ticket)] = entry
            return _prwatch_write(records)
        finally:
            _prwatch_release_file_lock(file_lock)


def _prwatch_remove(ticket):
    with _prwatch_lock:
        file_lock = _prwatch_acquire_file_lock()
        try:
            records = _prwatch_load()
            records.pop(str(ticket), None)
            return _prwatch_write(records)
        finally:
            _prwatch_release_file_lock(file_lock)


def _prwatch_update(ticket, **fields):
    with _prwatch_lock:
        file_lock = _prwatch_acquire_file_lock()
        try:
            records = _prwatch_load()
            entry = records.get(str(ticket))
            if not isinstance(entry, dict):
                return False
            records[str(ticket)].update(fields)
            return _prwatch_write(records)
        finally:
            _prwatch_release_file_lock(file_lock)


def _prwatch_list():
    with _prwatch_lock:
        return _prwatch_load()


__all__ = [
    "PRWATCH_PATH", "_prwatch_lock", "_prwatch_add", "_prwatch_has",
    "_prwatch_list", "_prwatch_load", "_prwatch_remove", "_prwatch_update",
]
