"""Quarantine bookkeeping for source items Aone permanently cannot resolve.

A candidate whose Aone work item does not exist is re-read every scan cycle
forever: the point-read collapses a 404 into the same ``None`` a timeout
returns, so no status is ever reported, the entry never reaches a terminal
state, and it is therefore still a candidate next cycle. Worse, an uncached
read failure aborts the entire ownership snapshot by contract, so one bad id
can keep the board's ownership projection stale indefinitely.

This module holds only the cross-cycle evidence needed to tell "permanently
absent" apart from "Aone was unhappy for a moment". The decision to stamp a
terminal status and to notify stays in the scan runner, which owns the
control-plane call.

Quarantining is deliberately hard to trigger. Stamping a terminal status on a
live work item would hide real work, and that is far worse than continuing to
waste a read, so each guard can veto on its own.
"""

from __future__ import annotations

import json
import os
import tempfile
import time
from pathlib import Path
from typing import Any, Callable, Dict

MIN_OBSERVATIONS = 2
MIN_AGE_SECONDS = 86400.0
# Real Aone work item ids are 8+ digits. A shorter all-digit id is not a
# plausible work item, which is the strongest permanent-absence signal
# available without trusting a single API response.
PLAUSIBLE_ID_DIGITS = 8


class _SourceNotFound:
    """Sentinel for "the Aone item does not exist", as opposed to a transient
    read failure.

    Deliberately falsy. Every existing branch that tests a resolved status for
    truthiness keeps its current behaviour if it has not been taught about this
    sentinel, so forgetting to handle it degrades to today's semantics rather
    than reporting a bogus status upstream.
    """

    __slots__ = ()

    def __bool__(self) -> bool:
        return False

    def __repr__(self) -> str:
        return "SOURCE_NOT_FOUND"


SOURCE_NOT_FOUND = _SourceNotFound()


class SourcePoisonLedger:
    """Durable per-task episodes of permanent source-absence observations."""

    def __init__(self, state_path: Path | str,
                 clock: Callable[[], float] = time.time) -> None:
        self._path = Path(state_path)
        self._clock = clock
        self._episodes: Dict[str, Dict[str, Any]] = {}
        try:
            raw = json.loads(self._path.read_text())
            loaded = raw.get("episodes")
            if isinstance(loaded, dict):
                self._episodes = {
                    str(key): dict(value)
                    for key, value in loaded.items()
                    if isinstance(value, dict)
                }
        except (OSError, ValueError, AttributeError):
            # A missing or corrupt ledger must not break the scan: the worst
            # case is that evidence restarts and quarantine is delayed.
            self._episodes = {}

    def episode(self, task_id: str) -> Dict[str, Any]:
        return dict(self._episodes.get(str(task_id)) or {})

    def record(self, task_id: str, aone_id: str) -> Dict[str, Any]:
        """Note one permanent-absence observation and return the episode."""
        now = float(self._clock())
        episode = self._episodes.setdefault(str(task_id), {
            "aoneId": str(aone_id),
            "firstSeenAt": now,
            "count": 0,
            "lastAlertAt": 0,
        })
        episode["aoneId"] = str(aone_id)
        episode["lastSeenAt"] = now
        episode["count"] = int(episode.get("count") or 0) + 1
        return dict(episode)

    def forget(self, task_id: str) -> None:
        """Drop the episode once an id resolves again.

        Without this, observations from unrelated outages weeks apart would
        accumulate until they satisfied the age and count guards together.
        """
        self._episodes.pop(str(task_id), None)

    @staticmethod
    def implausible_id(aone_id: str) -> bool:
        text = str(aone_id or "").strip()
        return text.isdigit() and len(text) < PLAUSIBLE_ID_DIGITS

    def should_quarantine(self, episode: Dict[str, Any], aone_id: str, *,
                          project_read_ok: bool) -> bool:
        """True only when all three guards agree the item is really gone.

        ``project_read_ok`` must mean "some other id in the same project
        resolved during this same pass", which is what rules out a broad Aone
        outage rather than a genuinely missing item.
        """
        if not project_read_ok:
            return False
        if not self.implausible_id(aone_id):
            return False
        if int(episode.get("count") or 0) < MIN_OBSERVATIONS:
            return False
        first = float(episode.get("firstSeenAt") or 0)
        last = float(episode.get("lastSeenAt") or 0)
        return (last - first) >= MIN_AGE_SECONDS

    def mark_alerted(self, task_id: str) -> None:
        episode = self._episodes.get(str(task_id))
        if episode is not None:
            episode["lastAlertAt"] = float(self._clock())

    def alerted(self, task_id: str) -> bool:
        return bool((self._episodes.get(str(task_id)) or {}).get("lastAlertAt"))

    def save(self) -> None:
        payload = json.dumps(
            {"version": 1, "episodes": self._episodes},
            ensure_ascii=False, indent=2)
        self._path.parent.mkdir(parents=True, exist_ok=True)
        handle, tmp = tempfile.mkstemp(dir=str(self._path.parent))
        try:
            with os.fdopen(handle, "w") as stream:
                stream.write(payload)
            os.replace(tmp, self._path)
        except BaseException:
            try:
                os.unlink(tmp)
            except OSError:
                pass
            raise
