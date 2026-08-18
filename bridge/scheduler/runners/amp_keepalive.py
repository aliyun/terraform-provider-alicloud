"""Read-only AMP query used to keep the Scheduler host session active."""

from __future__ import annotations

from datetime import datetime
import math
import os
from pathlib import Path
import re
import subprocess
from typing import Any, Callable, Mapping

from bridge.process_group_runner import run_process_group

from ..model import (
    IntervalSchedule,
    JobResult,
    JobResultStatus,
    ScheduledJobDefinition,
)


JOB_KEY = "amp.keepalive"
RUNNER_KEY = "amp_keepalive"
PROJECT_ID_ENV = "JARVIS_AMP_KEEPALIVE_PROJECT_ID"
TIMEOUT_ENV = "JARVIS_AMP_KEEPALIVE_TIMEOUT"
DEFAULT_TIMEOUT_SECONDS = 30.0
_PROJECT_ID_RE = re.compile(r"^[1-9][0-9]{0,18}$")

CommandRunner = Callable[..., Any]


class AmpKeepaliveRunner:
    """Run one bounded, non-mutating AMP business query."""

    def __init__(
        self,
        *,
        repo_root: Path,
        logger: Any,
        environ: Mapping[str, str] | None = None,
        command_runner: CommandRunner = run_process_group,
    ) -> None:
        self._repo_root = Path(repo_root).resolve()
        self._log = logger
        self._environ = os.environ if environ is None else environ
        self._command_runner = command_runner

    def run(
        self,
        definition: ScheduledJobDefinition,
        scheduled_for: datetime,
    ) -> JobResult:
        del scheduled_for
        definition_error = _definition_error(definition)
        if definition_error is not None:
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error=definition_error,
            )

        project_id = self._environ.get(PROJECT_ID_ENV)
        if project_id is None or project_id == "":
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="amp keepalive project id is missing",
            )
        if (
            not isinstance(project_id, str)
            or _PROJECT_ID_RE.fullmatch(project_id) is None
        ):
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="amp keepalive project id is invalid",
            )

        try:
            timeout = _positive_timeout(
                self._environ.get(TIMEOUT_ENV, str(DEFAULT_TIMEOUT_SECONDS)))
        except (TypeError, ValueError):
            return JobResult(
                JobResultStatus.PERMANENT_FAILURE,
                error="amp keepalive timeout is invalid",
            )

        command = [
            "/usr/bin/python3",
            "-I",
            str(self._repo_root / "bootstrap" / "amp_safe.py"),
            "--repo-root",
            str(self._repo_root),
            "branch",
            "list",
            "--project-id",
            project_id,
        ]
        try:
            completed = self._command_runner(
                command,
                cwd=str(self._repo_root),
                capture_output=True,
                text=True,
                timeout=timeout,
                stdin=subprocess.DEVNULL,
            )
        except subprocess.TimeoutExpired:
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error="amp keepalive query failed: timeout",
            )
        except Exception:  # noqa: BLE001 - never expose process or AMP details
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error="amp keepalive query failed: process_error",
            )

        if completed.returncode != 0:
            return JobResult(
                JobResultStatus.RETRYABLE_FAILURE,
                error="amp keepalive query failed: nonzero_exit",
            )
        self._log.info("AmpKeepaliveRunner: read-only query succeeded")
        return JobResult(JobResultStatus.SUCCEEDED)


def _definition_error(definition: ScheduledJobDefinition) -> str | None:
    if definition.id != JOB_KEY:
        return "amp keepalive runner received mismatched definition"
    if not isinstance(definition.schedule, IntervalSchedule):
        return "amp keepalive requires an interval schedule"
    return None


def _positive_timeout(value: object) -> float:
    if not isinstance(value, str) or not value:
        raise ValueError("timeout must be a non-empty string")
    timeout = float(value)
    if not math.isfinite(timeout) or timeout <= 0:
        raise ValueError("timeout must be finite and positive")
    return timeout


__all__ = ["AmpKeepaliveRunner", "JOB_KEY", "RUNNER_KEY"]
