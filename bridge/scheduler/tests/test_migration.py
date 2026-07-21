from __future__ import annotations

import unittest

from bridge.scheduler.jobs import JOBS
from bridge.scheduler.migration import (
    SchedulerMigrationError, business_job_enabled, requested_new_jobs,
    uses_new_engine,
)


JOB_KEYS = tuple(definition.id for definition in JOBS)


class SchedulerMigrationTests(unittest.TestCase):
    def test_empty_list_keeps_every_job_on_legacy(self):
        self.assertEqual(
            requested_new_jobs(JOB_KEYS, environ={}),
            frozenset(),
        )
        self.assertFalse(
            uses_new_engine("daily.probe", job_keys=JOB_KEYS, environ={}))

    def test_one_list_is_the_complete_new_engine_ownership_set(self):
        environ = {
            "JARVIS_SCHEDULER_NEW_JOBS": " daily.probe, aone.scan ",
        }
        self.assertEqual(
            requested_new_jobs(JOB_KEYS, environ=environ),
            frozenset(("daily.probe", "aone.scan")),
        )
        self.assertTrue(
            uses_new_engine("daily.probe", job_keys=JOB_KEYS, environ=environ))
        self.assertFalse(
            uses_new_engine("pr.watch", job_keys=JOB_KEYS, environ=environ))

    def test_unknown_or_blank_job_key_fails_closed(self):
        for value in (
            "daily.probe,unknown.job", "daily.probe,", ",daily.probe",
            "daily.probe,daily.probe",
        ):
            with self.subTest(value=value), self.assertRaises(SchedulerMigrationError):
                requested_new_jobs(
                    JOB_KEYS, environ={"JARVIS_SCHEDULER_NEW_JOBS": value})

    def test_deprecated_route_variables_cannot_be_combined_with_new_contract(self):
        for environ in (
            {"JARVIS_SCHEDULER_ENABLE": "1"},
            {"JARVIS_SCHEDULER_JOB_DAILY_PROBE": "new"},
        ):
            with self.subTest(environ=environ), self.assertRaisesRegex(
                    SchedulerMigrationError, "deprecated"):
                requested_new_jobs(JOB_KEYS, environ=environ)

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
