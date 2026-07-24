"""Small stateless helpers shared by Scheduler runner adapters.

Business rules live with their owning top-level domain modules; this package only
contains adapters selected by ``jobs.yaml``.
"""

from __future__ import annotations

from bridge.aone_workitems import claude_bin

__all__ = ["claude_bin"]
