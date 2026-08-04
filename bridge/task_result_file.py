#!/usr/bin/env python3
"""File transport for a Task run's structured result.

The stdout sentinel ``[[AONE_RESULT:{...}]]`` shares one channel with everything
else the run prints, and that channel is bounded: a long ``reply_body`` plus
``evidence``/``mr_cr_links`` arrays can hit the output limit mid-JSON, so the
closing ``]]`` never arrives and a finished run is scored as having produced
nothing. Measured on this bridge's own log, 18 of 72 ``missing_task_result``
failures ended with a tail that *is* the middle of the result JSON.

So the result also gets a file of its own:

    <repo>/.my-day/task-result/<item-id>.json

Writing it through ``bootstrap/task-result.sh`` validates the object against the
same contract the executor applies and fails loudly *inside the run*, where the
agent can still fix it — instead of after the fact, where the only remedy is a
full retry. The executor prefers this file and keeps the sentinel as a fallback,
so neither transport alone is a single point of failure.

Freshness is by deletion, not by timestamp: the executor removes any prior file
before spawning, so a file that exists afterwards was written by this run. Two
runs never race for one item — the Task key mutex and the Aone claim serialize
them.
"""

from __future__ import annotations

import json
import logging
import os
import re
import sys
from pathlib import Path
from typing import Any, Optional, Tuple

from bridge.helpers.aone import REPO_ROOT

log = logging.getLogger(__name__)

TASK_RESULT_DIR = Path(REPO_ROOT) / ".my-day" / "task-result"

# Item ids are Aone work-item numbers. Keep the character class narrow so a
# crafted id can never escape the directory.
_ITEM_ID_RE = re.compile(r"[A-Za-z0-9][A-Za-z0-9._-]{0,63}")


def _safe_item_id(item_id: Any) -> str:
    text = str(item_id or "").strip()
    if not _ITEM_ID_RE.fullmatch(text):
        raise ValueError("invalid task-result item id: %r" % (item_id,))
    return text


def task_result_path(item_id: Any) -> Path:
    """Absolute path of one item's result file."""
    return TASK_RESULT_DIR / ("%s.json" % _safe_item_id(item_id))


def clear_task_result(item_id: Any) -> bool:
    """Drop any prior result file. Returns True when one was removed."""
    try:
        path = task_result_path(item_id)
    except ValueError:
        return False
    try:
        path.unlink()
        return True
    except FileNotFoundError:
        return False
    except OSError as exc:
        log.warning("task-result: could not clear %s: %s", path, exc)
        return False


def write_task_result(item_id: Any, payload: Any) -> Path:
    """Atomically persist one result object. Raises on unusable input or I/O."""
    path = task_result_path(item_id)
    if not isinstance(payload, dict):
        raise ValueError("task result must be a JSON object")
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_name(path.name + ".tmp")
    tmp.write_text(json.dumps(payload, ensure_ascii=False), encoding="utf-8")
    os.replace(str(tmp), str(path))
    return path


def read_task_result(item_id: Any) -> Tuple[Optional[dict], str]:
    """Return ``(payload, note)``. ``payload`` is None when there is nothing usable.

    ``note`` is empty on a clean read and otherwise describes why the file was
    ignored, so a corrupt file degrades to the sentinel fallback with a reason in
    the log rather than silently.
    """
    try:
        path = task_result_path(item_id)
    except ValueError as exc:
        return None, str(exc)
    try:
        raw = path.read_text(encoding="utf-8")
    except FileNotFoundError:
        return None, ""
    except OSError as exc:
        return None, "unreadable: %s" % exc
    try:
        payload = json.loads(raw)
    except ValueError as exc:
        return None, "not valid JSON: %s" % exc
    if not isinstance(payload, dict):
        return None, "not a JSON object"
    return payload, ""


def _validate(payload: Any):
    """Apply the executor's own contract. Imported late to avoid an import cycle."""
    from bridge.persistent_tasks import task_result_correction, validate_task_result
    result, reason = validate_task_result(payload)
    return result, reason, task_result_correction(reason)


def main(argv=None) -> int:
    """``write <item-id>`` — validate JSON on stdin, then persist it."""
    args = list(sys.argv[1:] if argv is None else argv)
    if len(args) != 2 or args[0] != "write":
        print("usage: python3 -m bridge.task_result_file write <item-id>",
              file=sys.stderr)
        return 64
    item_id = args[1]
    raw = sys.stdin.read()
    if not raw.strip():
        print("task-result: refusing to write an empty result", file=sys.stderr)
        return 65
    try:
        payload = json.loads(raw)
    except ValueError as exc:
        print("task-result: stdin is not valid JSON: %s" % exc, file=sys.stderr)
        return 65
    try:
        _result, reason, correction = _validate(payload)
    except Exception as exc:  # noqa: BLE001 — never mask a validator defect as valid
        print("task-result: could not validate: %s" % exc, file=sys.stderr)
        return 70
    if reason:
        print("task-result: rejected (%s)" % reason, file=sys.stderr)
        if correction:
            print(correction, file=sys.stderr)
        return 65
    try:
        path = write_task_result(item_id, payload)
    except (OSError, ValueError) as exc:
        print("task-result: could not persist: %s" % exc, file=sys.stderr)
        return 74
    print(str(path))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
