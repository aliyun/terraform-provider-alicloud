"""The Scheduler's small, explicit runner catalogue.

The YAML registry names a runner key, and each key below has exactly one module.
"""

from __future__ import annotations

from pathlib import Path
from typing import Any

if (__package__ or "").startswith("bridge."):
    from bridge.headless_runtime import HeadlessRuntime, jarvis_transcript_exists, run_jarvis_attempt
else:
    from headless_runtime import HeadlessRuntime, jarvis_transcript_exists, run_jarvis_attempt

from . import (aone_workitem_ownership, claim_health, daily_nudge, daily_probe,
               owner_health, pr_watch, recovery, reply, scan,
               weekly_comment_participation)

HANDLER_KEYS = frozenset({
    scan.RUNNER_KEY,
    claim_health.RUNNER_KEY,
    owner_health.RUNNER_KEY,
    daily_nudge.RUNNER_KEY,
    pr_watch.RUNNER_KEY,
    reply.RUNNER_KEY,
    daily_probe.RUNNER_KEY,
    recovery.RUNNER_KEY,
    weekly_comment_participation.RUNNER_KEY,
    aone_workitem_ownership.RUNNER_KEY,
})
RUNNER_KEYS = HANDLER_KEYS


def build_runners(
    *, logger: Any, task_client: Any, worker_key: str, repo_root: Path,
    headless_runtime: HeadlessRuntime | None = None, summary_root: Path | None = None,
) -> dict[str, object]:
    runtime = headless_runtime or HeadlessRuntime(
        run_jarvis_attempt, transcript_exists=jarvis_transcript_exists, logger=logger)
    root = summary_root or (Path(__file__).resolve().parents[3] / "runs" / "probe")
    return {
        scan.RUNNER_KEY: scan.build(logger=logger, task_client=task_client, repo_root=repo_root),
        claim_health.RUNNER_KEY: claim_health.build(
            logger=logger, task_client=task_client, repo_root=repo_root),
        owner_health.RUNNER_KEY: owner_health.build(
            logger=logger, task_client=task_client, repo_root=repo_root),
        daily_nudge.RUNNER_KEY: daily_nudge.build(
            logger=logger, task_client=task_client, repo_root=repo_root),
        pr_watch.RUNNER_KEY: pr_watch.build_pr_watch_runners(
            logger=logger, task_client=task_client, repo_root=repo_root)[pr_watch.PR_WATCH_RUNNER_KEY],
        reply.RUNNER_KEY: reply.ReplyRunner(task_client=task_client, logger=logger),
        daily_probe.RUNNER_KEY: daily_probe.DailyProbeRunner(
            runtime=runtime, summary_root=root, logger=logger),
        recovery.RUNNER_KEY: recovery.ExternalRecoveryRunner(
            task_client=task_client, worker_key=worker_key, repo_root=repo_root, logger=logger),
        weekly_comment_participation.RUNNER_KEY: weekly_comment_participation.build(
            logger=logger, task_client=task_client, repo_root=repo_root),
        aone_workitem_ownership.RUNNER_KEY: aone_workitem_ownership.build(
            logger=logger, task_client=task_client, repo_root=repo_root),
    }


__all__ = ["HANDLER_KEYS", "RUNNER_KEYS", "build_runners"]
