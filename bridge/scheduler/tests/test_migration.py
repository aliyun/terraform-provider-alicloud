from __future__ import annotations

import unittest
from bridge.scheduler.model import HandlerRunner, HeadlessRunner, definition_snapshot
from bridge.scheduler.registry import JOBS, REGISTRY, default_registry_path


class SchedulerMigrationTests(unittest.TestCase):
    def test_checked_in_registry_contains_all_migrated_jobs(self):
        registry = REGISTRY
        self.assertEqual(tuple(job.id for job in JOBS),
                         ("daily.probe", "aone.scan", "aone.claim-health",
                          "daily.nudge", "aone.reply", "pr.watch",
                          "external.recovery"))
        self.assertEqual(registry.scheduler_job_keys(),
                         frozenset({"daily.probe", "aone.scan",
                                    "aone.claim-health", "daily.nudge", "aone.reply",
                                    "pr.watch", "external.recovery"}))
        self.assertEqual(registry.handler_keys(), frozenset({
            "aone.scan", "aone.claim-health", "daily.nudge", "aone.reply",
            "pr.watch", "external.recovery",
        }))
        self.assertEqual(
            registry.headless_builder_protocols(),
            frozenset({("probe.daily", "probe-result-v1")}),
        )
        self.assertIsInstance(JOBS[0].runner, HeadlessRunner)
        self.assertTrue(all(isinstance(job.runner, HandlerRunner)
                            for job in JOBS[1:]))
        self.assertEqual(set(definition_snapshot(JOBS[1])), {
            "id", "revision", "description", "schedule", "runner", "misfire",
            "retry_delay_seconds", "enabled",
        })
        self.assertEqual(default_registry_path().name, "jobs.yaml")
        self.assertEqual(default_registry_path().parent.name, "scheduler")

    def test_business_enable_is_part_of_the_definition(self):
        self.assertFalse(JOBS[0].enabled)
        self.assertEqual(JOBS[0].revision, 2)
        self.assertTrue(all(job.enabled for job in JOBS[1:]))


if __name__ == "__main__":
    unittest.main()
