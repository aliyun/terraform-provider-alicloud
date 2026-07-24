from __future__ import annotations

from datetime import datetime, timezone
import logging
from pathlib import Path
from tempfile import TemporaryDirectory
from types import SimpleNamespace
import unittest
from unittest import mock

from bridge.jarvis_task_router import ExecutionRouter
from bridge.scheduler.model import (
    HandlerRunner, IntervalSchedule, MisfirePolicy, ScheduledJobDefinition,
)
from bridge.scheduler.runners import build_runners
from bridge.scheduler.runners import claim_health, daily_nudge, scan
from bridge.scheduler.runners import pr_watch as pr


UTC = timezone.utc


class RecordingTaskClient:
    def __init__(self):
        self.upserts = []

    def upsert_desired_task(self, envelope, *, request_id):
        self.upserts.append((envelope, request_id))
        return {"accepted": True, "reason": "task_persisted"}


def definition(job_id: str) -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        job_id, 1, "test", IntervalSchedule(60, True), HandlerRunner(job_id),
        MisfirePolicy.COALESCE, 5,
    )


class SchedulerRuntimeRunnerTests(unittest.TestCase):
    def test_catalogue_uses_direct_domain_runners(self):
        runners = build_runners(
            logger=logging.getLogger(__name__),
            task_client=object(),
            worker_key="scheduler-test",
            repo_root=Path("/repo"),
        )

        self.assertIsInstance(runners["scan"], scan.ScanRunner)
        self.assertIsInstance(
            runners["claim_health"], claim_health.ClaimHealthRunner)
        self.assertIsInstance(runners["daily_nudge"], daily_nudge.DailyNudgeRunner)
        self.assertIsInstance(runners["pr_watch"], pr.PrWatchRunner)

    def test_nudge_runner_invokes_one_job(self):
        job = SimpleNamespace(enabled=True, calls=0)
        job.run = lambda: setattr(job, "calls", job.calls + 1)
        runner = daily_nudge.DailyNudgeRunner(job, logging.getLogger(__name__))

        result = runner.run(definition("daily.nudge"), datetime.now(UTC))

        self.assertEqual(result.status.value, "SUCCEEDED")
        self.assertEqual(job.calls, 1)

    def test_aone_scan_runs_one_tick(self):
        runtime = SimpleNamespace(calls=0)
        runtime._tick = lambda: setattr(runtime, "calls", runtime.calls + 1)
        runner = scan.ScanRunner(runtime)

        result = runner.run(definition("aone.scan"), datetime.now(UTC))

        self.assertEqual(result.status.value, "SUCCEEDED")
        self.assertEqual(runtime.calls, 1)

    def test_claim_health_reconciles_then_flushes(self):
        order = []
        runtime = SimpleNamespace(_claim_health_activity_cache={})
        runtime._scan_claimed = lambda: order.append("scan") or {"one": {}}
        runtime._reconcile_stale_claims = lambda snapshot: order.append(
            "reconcile:%s" % sorted(snapshot))
        with TemporaryDirectory() as directory, \
             mock.patch.object(
                 claim_health, "_aone_event_flush",
                 side_effect=lambda: order.append("aone-flush")), \
             mock.patch.object(
                 claim_health, "_dingtalk_event_flush",
                 side_effect=lambda: order.append("dingtalk-flush")):
            runner = claim_health.ClaimHealthRunner(
                runtime, Path(directory), logging.getLogger(__name__))
            result = runner.run(
                definition("aone.claim-health"), datetime.now(UTC))

        self.assertEqual(result.status.value, "SUCCEEDED")
        self.assertEqual(order, [
            "scan", "reconcile:['one']", "aone-flush", "dingtalk-flush"])

    def _watcher(self, client):
        router = ExecutionRouter(client=client, logger=logging.getLogger(__name__))
        return pr.PrWatchRuntime(
            execution_router=router, task_client=client)

    def test_pr_watch_persists_ci_fix_without_local_pool(self):
        client = RecordingTaskClient()
        with TemporaryDirectory() as directory, \
             mock.patch.object(pr, "PRWATCH_PATH",
                               Path(directory) / "pr-watch.json"):
            pr._prwatch_add(
                "test-item", "https://example.test/pull/1", "test-project",
                "Scheduler refactor")
            watcher = self._watcher(client)
            watcher._gh_pr_ci = lambda _url: (
                "head-1", ["unit-tests"], False)

            active = watcher._maybe_dispatch_ci_fix(
                "test-item", pr._prwatch_list()["test-item"])
            persisted = pr._prwatch_list()["test-item"]

        self.assertTrue(active)
        self.assertIsNone(watcher.pool)
        self.assertEqual(len(client.upserts), 1)
        envelope, request_id = client.upserts[0]
        self.assertEqual(envelope.task_type, "pr_ci_fix")
        self.assertEqual(envelope.source_ref["head"], "head-1")
        self.assertEqual(persisted["ci_fix_sha"], "head-1")
        self.assertEqual(persisted["ci_fix_attempts"], 1)
        self.assertEqual(request_id, envelope.request_id("upsert"))

    def test_pr_watch_persists_comment_reply_without_local_pool(self):
        client = RecordingTaskClient()
        with TemporaryDirectory() as directory, \
             mock.patch.object(pr, "PRWATCH_PATH",
                               Path(directory) / "pr-watch.json"):
            pr._prwatch_add(
                "test-item", "https://example.test/pull/1", "test-project",
                "Scheduler refactor")
            pr._prwatch_update("test-item", last_seen_comment="issue-1")
            watcher = self._watcher(client)
            watcher._gh_pr_comments = lambda _url: (
                "review-2", "reviewer", "please update")

            watcher._maybe_dispatch_comment_reply(
                "test-item", pr._prwatch_list()["test-item"])
            persisted = pr._prwatch_list()["test-item"]

        self.assertEqual(len(client.upserts), 1)
        envelope, request_id = client.upserts[0]
        self.assertEqual(envelope.task_type, "pr_comment_reply")
        self.assertEqual(envelope.source_ref["commentKey"], "review-2")
        self.assertEqual(persisted["last_seen_comment"], "review-2")
        self.assertEqual(request_id, envelope.request_id("upsert"))


if __name__ == "__main__":
    unittest.main()
