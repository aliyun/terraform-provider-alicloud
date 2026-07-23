"""Compatibility marker for the retired legacy scheduler-runner adapter.

All periodic jobs now have dedicated modules under :mod:`bridge.scheduler.runners`.
Kept as an import-safe module until downstream callers have moved off this path.
"""

__all__: list[str] = []
