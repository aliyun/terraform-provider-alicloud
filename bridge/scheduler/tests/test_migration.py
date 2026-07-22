from __future__ import annotations

import unittest
from bridge.scheduler.jobs import JOBS, REGISTRY


class SchedulerMigrationTests(unittest.TestCase):
    def test_checked_in_registry_contains_only_migrated_daily_probe(self):
        registry = REGISTRY
        self.assertEqual(tuple(route.job_key for route in registry.routes), ("daily.probe",))
        self.assertEqual(tuple(job.id for job in JOBS), ("daily.probe",))
        self.assertEqual(registry.scheduler_job_keys(), frozenset({"daily.probe"}))

    def test_business_enable_is_part_of_the_definition(self):
        self.assertTrue(JOBS[0].enabled)


if __name__ == "__main__":
    unittest.main()
