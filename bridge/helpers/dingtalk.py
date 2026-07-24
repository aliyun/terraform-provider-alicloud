"""Durable DingTalk lifecycle-event delivery."""

from __future__ import annotations

import json
import logging
import os
from pathlib import Path
import subprocess
import threading
import time
import uuid

from . import aone


log = logging.getLogger("jarvis-dingtalk-event")
DINGTALK_EVENT_PATH = Path(aone.REPO_ROOT) / ".my-day/bridge/dingtalk-event-ledger.json"
_event_lock = threading.RLock()
_event_inflight = set()


def _dingtalk_event_load():
    """Load the DingTalk ledger independently from the Aone ledger."""
    empty = {"pending": {}, "posted": {}, "suppressed": {}}
    try:
        if not DINGTALK_EVENT_PATH.exists():
            return empty
        raw = json.loads(DINGTALK_EVENT_PATH.read_text())
        return {
            name: raw.get(name) if isinstance(raw, dict) and isinstance(raw.get(name), dict) else {}
            for name in empty
        }
    except Exception as exc:  # noqa: BLE001
        log.warning("dingtalk-event: could not load %s: %s", DINGTALK_EVENT_PATH, exc)
        return empty


def _dingtalk_event_write(ledger):
    try:
        DINGTALK_EVENT_PATH.parent.mkdir(parents=True, exist_ok=True)
        tmp = DINGTALK_EVENT_PATH.with_name(DINGTALK_EVENT_PATH.name + ".tmp")
        tmp.write_text(json.dumps(ledger, ensure_ascii=False, default=str))
        os.replace(str(tmp), str(DINGTALK_EVENT_PATH))
        return True
    except Exception as exc:  # noqa: BLE001
        log.warning("dingtalk-event: could not persist %s: %s", DINGTALK_EVENT_PATH, exc)
        return False


def _dingtalk_event_out_track_id(event_digest, staff_id):
    return str(uuid.uuid5(uuid.NAMESPACE_URL, "jarvis-dingtalk:%s\x00%s" % (event_digest, staff_id)))


def _dingtalk_event_mark(ledger_id, record, bucket, state):
    with _event_lock:
        ledger = _dingtalk_event_load()
        ledger["pending"].pop(ledger_id, None)
        done = dict(record, state=state, completed_at=time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()))
        ledger[bucket][ledger_id] = done
        return _dingtalk_event_write(ledger)


def _dingtalk_result(stdout):
    for line in reversed(str(stdout or "").splitlines()):
        try:
            value = json.loads(line)
        except Exception:  # noqa: BLE001
            continue
        if isinstance(value, dict) and value.get("status") in ("sent", "skipped", "failed"):
            return value
    return None


def _dingtalk_event_publish_digest(ticket, project, event_digest, staff_id, title, text,
                                  allow_non_tf=False):
    ticket, project, staff_id = str(ticket or ""), str(project or ""), str(staff_id or "").strip()
    event_digest = str(event_digest or "").strip()
    title, text = aone._aone_event_sanitize_text(title, limit=120), aone._aone_event_sanitize_text(text)
    if (not ticket.isdigit() or not project or not staff_id or staff_id.startswith("WORKER_")
            or not aone._AONE_EVENT_DIGEST_RE.fullmatch(event_digest) or not title or not text):
        return False
    if not aone._is_terraform_project(project) and not allow_non_tf:
        return False
    ledger_id = aone._aone_event_ledger_id_from_digest(ticket, event_digest)
    now = time.time()
    with _event_lock:
        ledger = _dingtalk_event_load()
        if ledger_id in ledger["posted"] or ledger_id in ledger["suppressed"]:
            return True
        if ledger_id in _event_inflight:
            return False
        record = ledger["pending"].get(ledger_id)
        if not isinstance(record, dict):
            record = {"ticket": ticket, "project": project, "event_digest": event_digest,
                      "staff_id": staff_id, "title": title, "text": text, "state": "pending",
                      "attempts": 0, "receipt": _dingtalk_event_out_track_id(event_digest, staff_id),
                      "allow_non_tf": bool(allow_non_tf),
                      "created_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())}
        else:
            record.update(ticket=ticket, project=project, event_digest=event_digest,
                          staff_id=staff_id, title=title, text=text,
                          allow_non_tf=bool(record.get("allow_non_tf") or allow_non_tf))
            record.setdefault("receipt", _dingtalk_event_out_track_id(event_digest, staff_id))
        try:
            not_before = float(record.get("not_before") or 0)
        except (TypeError, ValueError):
            not_before = 0
        ledger["pending"][ledger_id] = record
        if not _dingtalk_event_write(ledger) or not_before > now:
            return False
        _event_inflight.add(ledger_id)
    try:
        with _event_lock:
            ledger = _dingtalk_event_load()
            record = ledger["pending"].get(ledger_id, record)
            record.update(state="posting", last_attempt_at=time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()))
            ledger["pending"][ledger_id] = record
            if not _dingtalk_event_write(ledger):
                return False
        try:
            proc = subprocess.run(
                [str(aone.REPO_ROOT / "bootstrap/notify-dingtalk.sh"), "--result-json",
                 "--out-track-id", record["receipt"], staff_id, title, text],
                capture_output=True, text=True, cwd=str(aone.REPO_ROOT), timeout=120)
            result = _dingtalk_result(proc.stdout)
            error = (proc.stderr or "").strip()[:300] if proc.returncode or result is None else ""
        except Exception as exc:  # noqa: BLE001
            proc, result, error = None, None, "notify transport exception: %s" % exc
        with _event_lock:
            ledger = _dingtalk_event_load()
            record = ledger["pending"].get(ledger_id, record)
            record["attempts"] = int(record.get("attempts") or 0) + 1
            if result and result.get("receipt"):
                record["receipt"] = str(result["receipt"])
            ledger["pending"][ledger_id] = record
            _dingtalk_event_write(ledger)
        if result and result.get("status") == "sent":
            return _dingtalk_event_mark(ledger_id, record, "posted", "posted")
        if result and result.get("status") == "skipped":
            record["reason"] = str(result.get("reason") or "suppressed")[:120]
            return _dingtalk_event_mark(ledger_id, record, "suppressed", "suppressed")
        with _event_lock:
            ledger = _dingtalk_event_load()
            record = ledger["pending"].get(ledger_id, record)
            record["state"] = "post_uncertain" if proc is None or (proc.returncode == 0 and result is None) else "failed"
            record["error"] = str((result or {}).get("reason") or error or "notify failed")[:300]
            record["not_before"] = time.time() + min(86400, 300 * (2 ** min(int(record["attempts"]) - 1, 8)))
            ledger["pending"][ledger_id] = record
            _dingtalk_event_write(ledger)
        return False
    finally:
        with _event_lock:
            _event_inflight.discard(ledger_id)


def _dingtalk_event_publish(ticket, project, event_key, staff_id, title, text, allow_non_tf=False):
    return _dingtalk_event_publish_digest(ticket, project, aone._aone_event_digest(event_key),
                                          staff_id, title, text, allow_non_tf=allow_non_tf)


def _dingtalk_event_enqueue(ticket, project, event_key, staff_id, title, text, allow_non_tf=False):
    if _dingtalk_event_publish(ticket, project, event_key, staff_id, title, text, allow_non_tf=allow_non_tf):
        return True
    ledger_id = aone._aone_event_ledger_id(ticket, event_key)
    with _event_lock:
        ledger = _dingtalk_event_load()
        return any(ledger_id in ledger[name] for name in ("pending", "posted", "suppressed"))


def _dingtalk_event_flush(limit=20):
    with _event_lock:
        pending = list(_dingtalk_event_load()["pending"].values())[:max(0, int(limit))]
    return sum(bool(isinstance(rec, dict)
                    and _dingtalk_event_publish_digest(
                        rec.get("ticket"), rec.get("project"), rec.get("event_digest"),
                        rec.get("staff_id"), rec.get("title"), rec.get("text"),
                        allow_non_tf=bool(rec.get("allow_non_tf"))))
               for rec in pending)
