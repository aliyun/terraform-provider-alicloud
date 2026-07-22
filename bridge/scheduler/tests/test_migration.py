from __future__ import annotations

import unittest
from bridge.scheduler.jobs import JOBS, REGISTRY


class SchedulerMigrationTests(unittest.TestCase):
    def test_checked_in_registry_contains_probe_and_all_smoke_schedule_types(self):
        registry = REGISTRY
        self.assertEqual(tuple(route.job_key for route in registry.routes),
                         ("daily.probe", "smoke.interval", "smoke.daily", "smoke.adaptive"))
        self.assertEqual(tuple(job.id for job in JOBS),
                         ("daily.probe", "smoke.interval", "smoke.daily", "smoke.adaptive"))
        self.assertEqual(registry.scheduler_job_keys(),
                         frozenset({"daily.probe", "smoke.interval", "smoke.daily", "smoke.adaptive"}))

    def test_business_enable_is_part_of_the_definition(self):
        self.assertFalse(JOBS[0].enabled)
        self.assertTrue(all(job.enabled for job in JOBS[1:]))


if __name__ == "__main__":
    unittest.main()
