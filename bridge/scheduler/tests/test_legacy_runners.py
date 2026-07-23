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
from bridge.scheduler.runners.legacy import (
    AoneClaimHealthRunner, AoneScanRunner, LegacyBridgeContext, PrWatchRunner,
)


UTC = timezone.utc


def interval_definition(job_id: str) -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        job_id, 1, "test", IntervalSchedule(60, True), HandlerRunner(job_id),
        MisfirePolicy.COALESCE, 5,
    )


class LegacyRunnerTests(unittest.TestCase):
    def test_aone_scan_runs_one_tick_without_starting_a_legacy_loop(self):
        scanner = SimpleNamespace(_tick_calls=0)
        scanner._tick = lambda: setattr(scanner, "_tick_calls", scanner._tick_calls + 1)
        handler = SimpleNamespace(scanner=scanner)
        context = LegacyBridgeContext(handler_factory=lambda: handler)
        runner = AoneScanRunner(job_id="aone.scan", context=context,
                                logger=logging.getLogger(__name__))

        result = runner.run(interval_definition("aone.scan"), datetime.now(UTC))

        self.assertEqual(scanner._tick_calls, 1)
        self.assertEqual(result.status.value, "SUCCEEDED")

    def test_claim_health_keeps_flush_after_reconciliation(self):
        with TemporaryDirectory() as directory:
            order = []
            scanner = SimpleNamespace(_claim_health_activity_cache={})
            scanner._scan_claimed = lambda: order.append("scan") or {"one": {}}
            scanner._reconcile_stale_claims = lambda snapshot: order.append(
                "reconcile:%s" % sorted(snapshot))
            handler = SimpleNamespace(scanner=scanner)
            module = SimpleNamespace(
                REPO_ROOT=Path(directory),
                _aone_event_flush=lambda: order.append("aone-flush"),
                _dingtalk_event_flush=lambda: order.append("dingtalk-flush"),
            )
            context = LegacyBridgeContext(handler_factory=lambda: handler,
                                          legacy_module=module)
            runner = AoneClaimHealthRunner(
                job_id="aone.claim-health", context=context,
                logger=logging.getLogger(__name__))

            result = runner.run(interval_definition("aone.claim-health"), datetime.now(UTC))

        self.assertEqual(result.status.value, "SUCCEEDED")
        self.assertEqual(order, ["scan", "reconcile:['one']", "aone-flush", "dingtalk-flush"])

    def test_pr_watch_returns_legacy_active_backoff_as_next_due(self):
        watcher = SimpleNamespace(enabled=True, _active_interval=10, interval=300)
        watcher._tick = lambda: True
        handler = SimpleNamespace(prwatch=watcher)
        module = SimpleNamespace(
            _aone_event_flush=lambda: None,
            _dingtalk_event_flush=lambda: None,
        )
        context = LegacyBridgeContext(handler_factory=lambda: handler,
                                      legacy_module=module)
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
