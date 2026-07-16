#!/usr/bin/env python3
"""Shared, thread-safe execution capacity for all local executors.

``TaskExecutor`` and ``EphemeralExecutor`` must acquire from the same manager
before starting work.  A successful acquire returns an idempotent permit; when
no slot is available it returns ``None`` without queuing or blocking.
"""

from __future__ import annotations

import threading
from dataclasses import dataclass
from typing import Dict, Optional


@dataclass(frozen=True)
class CapacitySnapshot:
    """One internally consistent view of the shared execution capacity."""

    capacity: int
    running: int
    available: int


class CapacityPermit:
    """Ownership of one running slot.

    Callers should normally use the permit as a context manager.  ``release``
    is idempotent so cleanup paths may safely converge on it.
    """

    def __init__(self, manager: "CapacityManager", permit_id: int,
                 owner: Optional[str]):
        self._manager = manager
        self._permit_id = permit_id
        self._owner = owner
        self._lock = threading.Lock()
        self._released = False

    @property
    def owner(self) -> Optional[str]:
        return self._owner

    @property
    def released(self) -> bool:
        with self._lock:
            return self._released

    def release(self) -> bool:
        """Release the slot once; return ``True`` only for the first release."""
        with self._lock:
            if self._released:
                return False
            self._manager._release(self._permit_id)
            self._released = True
            return True

    def __enter__(self) -> "CapacityPermit":
        return self

    def __exit__(self, _exc_type, _exc_value, _traceback) -> bool:
        self.release()
        return False


class CapacityManager:
    """Atomically allocates a fixed number of process execution slots."""

    def __init__(self, capacity: int):
        if isinstance(capacity, bool):
            raise ValueError("capacity must be a non-negative integer")
        try:
            value = int(capacity)
        except (TypeError, ValueError):
            raise ValueError("capacity must be a non-negative integer")
        if value < 0 or value != capacity:
            raise ValueError("capacity must be a non-negative integer")
        self._capacity = value
        self._lock = threading.Lock()
        self._next_permit_id = 1
        self._active: Dict[int, Optional[str]] = {}

    @property
    def capacity(self) -> int:
        return self._capacity

    def acquire(self, owner: Optional[str] = None) -> Optional[CapacityPermit]:
        """Try to reserve one slot atomically; return ``None`` when full."""
        with self._lock:
            if len(self._active) >= self._capacity:
                return None
            permit_id = self._next_permit_id
            self._next_permit_id += 1
            self._active[permit_id] = owner
        return CapacityPermit(self, permit_id, owner)

    def release(self, permit: CapacityPermit) -> bool:
        """Convenience wrapper for callers that do not own cleanup directly."""
        if not isinstance(permit, CapacityPermit) or permit._manager is not self:
            raise ValueError("permit does not belong to this capacity manager")
        return permit.release()

    def snapshot(self) -> CapacitySnapshot:
        """Return running and available counts from one locked observation."""
        with self._lock:
            running = len(self._active)
            return CapacitySnapshot(
                capacity=self._capacity,
                running=running,
                available=self._capacity - running,
            )

    def running_count(self) -> int:
        return self.snapshot().running

    def available_slots(self) -> int:
        return self.snapshot().available

    def _release(self, permit_id: int) -> None:
        with self._lock:
            if permit_id not in self._active:
                raise RuntimeError("capacity permit is not active")
            del self._active[permit_id]
