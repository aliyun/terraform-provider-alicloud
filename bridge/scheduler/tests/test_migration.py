from __future__ import annotations

import unittest
from bridge.scheduler.model import HandlerRunner, definition_snapshot
from bridge.scheduler.registry import JOBS, REGISTRY, default_registry_path


class SchedulerMigrationTests(unittest.TestCase):
    def test_checked_in_registry_contains_all_smoke_schedule_types(self):
        # The standalone scheduler runtime (bridge/main.py) wires only the smoke
        # runner, so the checked-in registry is smoke-only for the first cut.
        # daily.probe migration is deferred until its Ephemeral runner is wired.
        registry = REGISTRY
        self.assertEqual(tuple(job.id for job in JOBS),
                         ("smoke.interval", "smoke.daily", "smoke.adaptive"))
        self.assertEqual(registry.scheduler_job_keys(),
                         frozenset({"smoke.interval", "smoke.daily", "smoke.adaptive"}))
        self.assertEqual(registry.handler_keys(), frozenset({"scheduler.smoke"}))
        self.assertTrue(all(isinstance(job.runner, HandlerRunner) for job in JOBS))
        self.assertEqual(set(definition_snapshot(JOBS[0])), {
            "id", "revision", "description", "schedule", "runner", "misfire",
            "retry_delay_seconds", "enabled",
        })
        self.assertEqual(default_registry_path().name, "jobs.yaml")
        self.assertEqual(default_registry_path().parent.name, "scheduler")

    def test_business_enable_is_part_of_the_definition(self):
        self.assertTrue(all(job.enabled for job in JOBS))


if __name__ == "__main__":
    unittest.main()
