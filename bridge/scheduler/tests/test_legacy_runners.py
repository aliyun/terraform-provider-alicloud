from __future__ import annotations

from datetime import datetime, timezone
import logging
from pathlib import Path
from tempfile import TemporaryDirectory
from types import SimpleNamespace
import unittest

from bridge.scheduler import (
    AdaptiveSchedule, HandlerRunner, IntervalSchedule, MisfirePolicy,
    ScheduledJobDefinition,
)
from bridge.scheduler.runners.aone import (
    AoneClaimHealthRunner, AoneRuntimeContext, AoneScanRunner,
)
from bridge.scheduler.runners.pr import PrWatchRunner, PrWatchRuntimeContext
from bridge.scheduler.runners.nudge import DailyNudgeRunner, NudgeRuntimeContext
from bridge.scheduler.runners import build_runners


UTC = timezone.utc


class RecordingTaskClient:
    def __init__(self):
        self.upserts = []

    def upsert_desired_task(self, envelope, *, request_id):
        self.upserts.append((envelope, request_id))
        return {"accepted": True, "reason": "task_persisted"}


def interval_definition(job_id: str) -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        job_id, 1, "test", IntervalSchedule(60, True), HandlerRunner(job_id),
        MisfirePolicy.COALESCE, 5,
    )


class SchedulerRuntimeRunnerTests(unittest.TestCase):
    def test_runtime_context_constructs_no_handler_or_local_executor(self):
        scanner = object()
        observed = {}

        class Router:
            def __init__(self, *, client, logger):
                observed["router"] = (client, logger)

        class FieldRepair:
            def __init__(self, **kwargs):
                observed["field_repair"] = kwargs

        module = SimpleNamespace(
            ExecutionRouter=Router,
            FieldRepairWorker=FieldRepair,
            DEFAULT_EXECUTION_RUNTIME=object(),
            claude_bin=lambda: "claude",
            AoneScheduler=lambda context, pool: observed.update(
                context=context, pool=pool) or scanner,
        )
        client = object()
        context = AoneRuntimeContext(
            task_client=client,
            repo_root=Path("/repo"),
            logger=logging.getLogger(__name__),
            legacy_module=module,
        )

        self.assertIs(context.scanner, scanner)
        self.assertIs(context.task_client, client)
        self.assertIsNone(context.ephemeral_executor)
        self.assertTrue(context.no_dingtalk)
        self.assertIsNone(observed["pool"])
        self.assertIs(observed["context"], context)

    def test_catalogue_registers_dedicated_aone_runners(self):
        runners = build_runners(
            logger=logging.getLogger(__name__),
            task_client=object(),
            worker_key="scheduler-test",
            repo_root=Path("/repo"),
        )

        self.assertIsInstance(runners["aone.scan"], AoneScanRunner)
        self.assertIsInstance(runners["aone.claim-health"], AoneClaimHealthRunner)

    def test_nudge_context_runs_without_constructing_a_handler(self):
        observed = {}

        class NudgeJob:
            enabled = True

            def __init__(self, context):
                observed["context"] = context

            def run(self):
                observed["runs"] = observed.get("runs", 0) + 1

        context = NudgeRuntimeContext(
            legacy_module=SimpleNamespace(_NudgeJob=NudgeJob))
        runner = DailyNudgeRunner(
            job_id="daily.nudge", context=context,
            logger=logging.getLogger(__name__))

        result = runner.run(
            interval_definition("daily.nudge"), datetime.now(UTC))

        self.assertEqual(result.status.value, "SUCCEEDED")
        self.assertEqual(observed["runs"], 1)
        self.assertIs(observed["context"], context)
        self.assertIsNone(context.ephemeral_executor)
        self.assertIsNone(context.execution_router)

    def test_catalogue_registers_dedicated_nudge_runner(self):
        runners = build_runners(
            logger=logging.getLogger(__name__),
            task_client=object(),
            worker_key="scheduler-test",
            repo_root=Path("/repo"),
        )

        self.assertIsInstance(runners["daily.nudge"], DailyNudgeRunner)

    def test_pr_watch_context_constructs_no_handler_or_local_executor(self):
        watcher = object()
        observed = {}

        class Router:
            def __init__(self, *, client, logger):
                observed["router"] = (client, logger)

        module = SimpleNamespace(
            ExecutionRouter=Router,
            PrWatchScheduler=lambda context, pool: observed.update(
                context=context, pool=pool) or watcher,
        )
        client = object()
        context = PrWatchRuntimeContext(
            task_client=client,
            logger=logging.getLogger(__name__),
            legacy_module=module,
        )

        self.assertIs(context.watcher, watcher)
        self.assertIs(context.task_client, client)
        self.assertIsNone(context.ephemeral_executor)
        self.assertTrue(context.no_dingtalk)
        self.assertFalse(hasattr(context, "handler"))
        self.assertIsNone(observed["pool"])
        self.assertIs(observed["context"], context)

    def test_pr_watch_context_persists_ci_fix_without_local_pool(self):
        from bridge import jarvis_dingtalk_bot as module

        client = RecordingTaskClient()
        with TemporaryDirectory() as directory:
            original_path = module.PRWATCH_PATH
            module.PRWATCH_PATH = Path(directory) / "pr-watch.json"
            try:
                module._prwatch_add(
                    "test-item", "https://example.test/pull/1", "test-project",
                    "Scheduler refactor")
                context = PrWatchRuntimeContext(
                    task_client=client,
                    logger=logging.getLogger(__name__),
                    legacy_module=module,
                )
                watcher = context.watcher
                watcher._gh_pr_ci = lambda _url: (
                    "head-1", ["unit-tests"], False)

                active = watcher._maybe_dispatch_ci_fix(
                    "test-item", module._prwatch_list()["test-item"])
                persisted = module._prwatch_list()["test-item"]
            finally:
                module.PRWATCH_PATH = original_path

        self.assertTrue(active)
        self.assertIsNone(watcher.pool)
        self.assertEqual(len(client.upserts), 1)
        envelope, request_id = client.upserts[0]
        self.assertEqual(envelope.task_type, "pr_ci_fix")
        self.assertEqual(envelope.source_ref["head"], "head-1")
        self.assertEqual(persisted["ci_fix_sha"], "head-1")
        self.assertEqual(persisted["ci_fix_attempts"], 1)
        self.assertEqual(request_id, envelope.request_id("upsert"))

    def test_pr_watch_context_persists_comment_reply_without_local_pool(self):
        from bridge import jarvis_dingtalk_bot as module

        client = RecordingTaskClient()
        with TemporaryDirectory() as directory:
            original_path = module.PRWATCH_PATH
            module.PRWATCH_PATH = Path(directory) / "pr-watch.json"
            try:
                module._prwatch_add(
                    "test-item", "https://example.test/pull/1", "test-project",
                    "Scheduler refactor")
                module._prwatch_update(
                    "test-item", last_seen_comment="issue-1")
                context = PrWatchRuntimeContext(
                    task_client=client,
                    logger=logging.getLogger(__name__),
                    legacy_module=module,
                )
                watcher = context.watcher
                watcher._gh_pr_comments = lambda _url: (
                    "review-2", "reviewer", "please update")

                watcher._maybe_dispatch_comment_reply(
                    "test-item", module._prwatch_list()["test-item"])
                persisted = module._prwatch_list()["test-item"]
            finally:
                module.PRWATCH_PATH = original_path

        self.assertIsNone(watcher.pool)
        self.assertEqual(len(client.upserts), 1)
        envelope, request_id = client.upserts[0]
        self.assertEqual(envelope.task_type, "pr_comment_reply")
        self.assertEqual(envelope.source_ref["commentKey"], "review-2")
        self.assertEqual(persisted["last_seen_comment"], "review-2")
        self.assertEqual(request_id, envelope.request_id("upsert"))

    def test_catalogue_registers_dedicated_pr_watch_runner(self):
        runners = build_runners(
            logger=logging.getLogger(__name__),
            task_client=object(),
            worker_key="scheduler-test",
            repo_root=Path("/repo"),
        )

        self.assertIsInstance(runners["pr.watch"], PrWatchRunner)

    def test_aone_scan_runs_one_tick_without_starting_a_legacy_loop(self):
        scanner = SimpleNamespace(_tick_calls=0)
        scanner._tick = lambda: setattr(scanner, "_tick_calls", scanner._tick_calls + 1)
        context = SimpleNamespace(scanner=scanner)
        runner = AoneScanRunner(job_id="aone.scan", context=context,
                                logger=logging.getLogger(__name__))

        result = runner.run(interval_definition("aone.scan"), datetime.now(UTC))

        self.assertEqual(scanner._tick_calls, 1)
        self.assertEqual(result.status.value, "SUCCEEDED")

    def test_claim_health_keeps_flush_after_reconciliation(self):
        order = []
        scanner = SimpleNamespace(_claim_health_activity_cache={})
        scanner._scan_claimed = lambda: order.append("scan") or {"one": {}}
        scanner._reconcile_stale_claims = lambda snapshot: order.append(
            "reconcile:%s" % sorted(snapshot))
        with TemporaryDirectory() as directory:
            module = SimpleNamespace(
                REPO_ROOT=Path(directory),
                _aone_event_flush=lambda: order.append("aone-flush"),
                _dingtalk_event_flush=lambda: order.append("dingtalk-flush"),
            )
            context = SimpleNamespace(scanner=scanner, module=module)
            runner = AoneClaimHealthRunner(
                job_id="aone.claim-health", context=context,
                logger=logging.getLogger(__name__))

            result = runner.run(interval_definition("aone.claim-health"), datetime.now(UTC))

        self.assertEqual(result.status.value, "SUCCEEDED")
        self.assertEqual(order, ["scan", "reconcile:['one']", "aone-flush", "dingtalk-flush"])

    def test_pr_watch_returns_legacy_active_backoff_as_next_due(self):
        watcher = SimpleNamespace(enabled=True, _active_interval=10, interval=300)
        watcher._tick = lambda: True
        module = SimpleNamespace(
            _aone_event_flush=lambda: None,
            _dingtalk_event_flush=lambda: None,
        )
        context = SimpleNamespace(watcher=watcher, module=module)
        runner = PrWatchRunner(job_id="pr.watch", context=context,
                               logger=logging.getLogger(__name__))
        definition = ScheduledJobDefinition(
            "pr.watch", 1, "test", AdaptiveSchedule(5, 300, 3600, False),
            HandlerRunner("pr.watch"), MisfirePolicy.WAIT_FOR_COMPLETION, 5,
        )

        result = runner.run(definition, datetime.now(UTC))

        self.assertEqual(result.status.value, "SUCCEEDED")
        self.assertIsNotNone(result.next_due_at)


if __name__ == "__main__":
    unittest.main()
