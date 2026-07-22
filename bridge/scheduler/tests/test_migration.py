from __future__ import annotations

import unittest
from types import SimpleNamespace

from bridge.scheduler.jobs import JOBS, REGISTRY
from bridge.scheduler.migration import (
    SchedulerMigrationError, business_job_enabled,
)


class SchedulerMigrationTests(unittest.TestCase):
    def test_checked_in_registry_is_empty_until_a_job_is_migrated(self):
        registry = REGISTRY
        self.assertEqual(tuple(route.job_key for route in registry.routes), ())
        self.assertEqual(JOBS, ())
        self.assertEqual(registry.scheduler_job_keys(), frozenset())

    def test_business_enable_is_independent_from_route(self):
        probe = SimpleNamespace(enabled_env="JARVIS_PROBE_SCHED")
        reply = SimpleNamespace(enabled_env=None)
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
