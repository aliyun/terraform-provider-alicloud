"""Durable discovery snapshot for the Aone scan runner.

The scan runner diffs the current pool union against the previous tick to decide
which tickets are new or externally updated.  That diff lived only in process
memory, so every Scheduler restart made all ~385 union members look new: one very
expensive cold tick (measured ~33 minutes), plus a long tail of idle tickets
re-entering the lag-window retry path.  This store keeps the minimum the diff
needs — ``{aoneId: modified}`` — across restarts.

**Correctness rule** (the reason this file stays small and boring): the snapshot
decides only *whether a ticket gets re-evaluated*, never *what the decision is* —
the decision is a pure function of Aone state (tags, status, comments).  So the
two ways of being wrong are not symmetric:

- A **missing** entry costs one re-evaluation, which reaches the same conclusion.
- A **spurious** entry skips re-evaluation, silently stranding that ticket until
  its next Aone modification.  That is a lost ticket.

Therefore the runner commits an entry only after that ticket's decision reached a
conclusive resting state.  Inconclusive outcomes — Aone query failure, a pending
lag window, a failed dispatch — are deliberately left out so the next tick (or
the next process) evaluates them again.

A missing or unreadable file degrades to the pre-existing behaviour: everything
looks new, one expensive tick, then the file is warm again.  Failure direction is
cost, never correctness.

``commit`` is read-modify-write (load, merge, prune, replace), so concurrent
writers silently lose one side's merge even though ``os.replace`` keeps the file
itself intact.  Measured without the lock: 4 processes × 25 commits of distinct
keys landed 81 of 100 entries — a 19% loss, not a theoretical race.  Today only
the Scheduler ticks, but all three bridge processes construct a ``ScanRunner`` at
startup, so single-writer is an accident of who happens to schedule the job
rather than something the code guarantees.  ``commit`` therefore takes an
exclusive ``flock`` (same pattern as ``bridge/pending_dispatch.py``, which needs
it because Bot and Scheduler share that file).  ``load`` stays lock-free:
``os.replace`` is atomic, so a reader always sees a complete file.
"""

from __future__ import annotations

import fcntl
import json
import os
from pathlib import Path
import threading
from typing import Any, Callable, Iterable, Mapping


DEFAULT_PATH = (
    Path(__file__).resolve().parents[1]
    / ".my-day"
    / "bridge"
    / "scan-snapshot.json"
)

# Bounds the file against unbounded growth if pruning ever misses (e.g. a caller
# that never passes `keep`). Far above the real union size (~400).
MAX_ENTRIES = 20000


class ScanSnapshotStore:
    """Atomic ``{aoneId: modified}`` store for cross-restart discovery diffing."""

    def __init__(self, path: str | os.PathLike[str] | None = None) -> None:
        configured = path or os.environ.get("JARVIS_SCAN_SNAPSHOT_PATH")
        self.path = Path(configured) if configured else DEFAULT_PATH
        self.lock_path = self.path.with_name(self.path.name + ".lock")

    def load(self) -> dict[str, str]:
        """Return the persisted ``{aoneId: modified}`` map.

        Best-effort by design: a missing, truncated, or malformed file returns an
        empty map, which reproduces the old all-new cold start rather than
        blocking the tick.
        """
        try:
            raw = json.loads(self.path.read_text(encoding="utf-8"))
        except FileNotFoundError:
            return {}
        except Exception:  # noqa: BLE001 — corrupt file must not stall discovery
            return {}
        entries = raw.get("entries") if isinstance(raw, dict) else None
        if not isinstance(entries, dict):
            return {}
        out = {}
        for key, value in entries.items():
            aone_id = str(key or "").strip()
            modified = str(value or "").strip()
            if aone_id and modified:
                out[aone_id] = modified
        return out

    def commit(
        self,
        conclusive: Mapping[str, str],
        *,
        keep: Iterable[str] | None = None,
    ) -> dict[str, str]:
        """Merge conclusive observations, prune to ``keep``, and write atomically.

        ``conclusive`` holds only tickets whose decision this tick reached a
        resting state (see module docstring).  ``keep`` is the current union id
        set; entries outside it are dropped so the file tracks the live union.

        A partial union failure therefore prunes that pool's entries and costs one
        re-evaluation pass next tick — cost, not correctness. Returns the map that
        was written so callers can keep their in-memory view aligned.

        Holds an exclusive lock for the whole read-modify-write so a second writer
        cannot drop this merge (see module docstring).
        """
        return self._locked(lambda: self._commit_locked(conclusive, keep))

    def _commit_locked(
        self,
        conclusive: Mapping[str, str],
        keep: Iterable[str] | None,
    ) -> dict[str, str]:
        merged = self.load()
        for key, value in (conclusive or {}).items():
            aone_id = str(key or "").strip()
            modified = str(value or "").strip()
            if aone_id and modified:
                merged[aone_id] = modified
        if keep is not None:
            keep_ids = {str(k or "").strip() for k in keep}
            merged = {k: v for k, v in merged.items() if k in keep_ids}
        if len(merged) > MAX_ENTRIES:
            # Deterministic trim; the union is orders of magnitude smaller, so this
            # only ever fires on a caller bug.
            merged = dict(sorted(merged.items())[:MAX_ENTRIES])
        self._write({"version": 1, "entries": merged})
        return merged

    def _locked(self, callback: Callable[[], Any]) -> Any:
        self.lock_path.parent.mkdir(parents=True, exist_ok=True)
        with self.lock_path.open("a+", encoding="utf-8") as lock_file:
            os.chmod(self.lock_path, 0o600)
            fcntl.flock(lock_file.fileno(), fcntl.LOCK_EX)
            try:
                return callback()
            finally:
                fcntl.flock(lock_file.fileno(), fcntl.LOCK_UN)

    def _write(self, value: Mapping[str, object]) -> None:
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
