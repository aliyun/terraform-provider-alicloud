"""Role fence shared by the standalone Scheduler entry points."""

from __future__ import annotations

import os
from typing import Mapping


def require_scheduler_role(environ: Mapping[str, str] | None = None) -> None:
    """Refuse to admit scheduled work on an executor-only host."""

    role = (os.environ if environ is None else environ).get(
        "JARVIS_BRIDGE_ROLE", "scheduler").strip()
    if role != "scheduler":
        raise RuntimeError(
            "bridge/scheduler.sh requires JARVIS_BRIDGE_ROLE=scheduler; "
            "worker hosts run bridge/run.sh only")


__all__ = ["require_scheduler_role"]
