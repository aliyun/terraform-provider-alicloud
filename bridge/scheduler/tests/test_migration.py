from __future__ import annotations

import unittest

from bridge.scheduler.jobs import JOBS, REGISTRY
from bridge.scheduler.migration import (
    SchedulerMigrationError, business_job_enabled,
)


JOB_KEYS = tuple(definition.id for definition in JOBS)


class SchedulerMigrationTests(unittest.TestCase):
    def test_checked_in_registry_is_the_complete_ownership_snapshot(self):
        registry = REGISTRY
        self.assertEqual(
            frozenset(route.job_key for route in registry.routes),
            frozenset(JOB_KEYS),
        )
        self.assertEqual(registry.scheduler_job_keys(), frozenset())

    def test_business_enable_is_independent_from_route(self):
        probe = next(item for item in JOBS if item.id == "daily.probe")
        reply = next(item for item in JOBS if item.id == "aone.reply")
        self.assertFalse(business_job_enabled(
            probe, environ={"JARVIS_PROBE_SCHED": "0"}))
        self.assertTrue(business_job_enabled(
            probe, environ={"JARVIS_PROBE_SCHED": "1"}))
        self.assertTrue(business_job_enabled(reply, environ={}))
        with self.assertRaisesRegex(SchedulerMigrationError, "must be 0 or 1"):
            business_job_enabled(
                probe, environ={"JARVIS_PROBE_SCHED": "sometimes"})


if __name__ == "__main__":
    unittest.main()
