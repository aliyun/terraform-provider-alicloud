#!/usr/bin/env python3
"""Focused shared-capacity tests for EphemeralExecutor."""

import sys
import tempfile
import threading
import time
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from jarvis_capacity import CapacityManager  # noqa: E402
from jarvis_dingtalk_bot import DispatchPool, EphemeralExecutor  # noqa: E402


class EphemeralExecutorTest(unittest.TestCase):
    def make(self, manager, directory):
        return EphemeralExecutor(
            max_workers=1, queue_max=2, capacity_manager=manager,
            ledger_path=Path(directory) / "ledger.json")

    def test_formal_name_keeps_compatibility_alias(self):
        self.assertIs(DispatchPool, EphemeralExecutor)

    def test_waits_for_shared_capacity_without_starting_work(self):
        manager = CapacityManager(1)
        task_permit = manager.acquire("task:busy")
        started = threading.Event()
        finished = threading.Event()

        with tempfile.TemporaryDirectory() as directory:
            executor = self.make(manager, directory)
            accepted, reason = executor.submit(
                "job-1", lambda: (started.set(), finished.set()))
            self.assertTrue(accepted, reason)
            self.assertFalse(started.wait(0.1))
            self.assertEqual(executor.free_slots(), 0)

            task_permit.release()
            self.assertTrue(started.wait(2))
            self.assertTrue(finished.wait(2))
            deadline = time.monotonic() + 2
            while executor.active_count() and time.monotonic() < deadline:
                time.sleep(0.01)
            self.assertEqual(manager.running_count(), 0)
            executor.shutdown(wait=True)

    def test_watchdog_does_not_age_queued_capacity_wait(self):
        manager = CapacityManager(1)
        task_permit = manager.acquire("task:busy")
        with tempfile.TemporaryDirectory() as directory:
            executor = self.make(manager, directory)
            executor._watchdog_threshold = 0
            accepted, _ = executor.submit("job-1", lambda: None)
            self.assertTrue(accepted)
            time.sleep(0.05)
            executor._watchdog_tick()
            self.assertEqual(executor.active_ids(), ["job-1"])
            task_permit.release()
            executor.shutdown(wait=True)

    def test_watchdog_releases_shared_capacity_for_stuck_worker(self):
        manager = CapacityManager(1)
        started = threading.Event()
        unblock = threading.Event()

        def work():
            started.set()
            unblock.wait(2)

        with tempfile.TemporaryDirectory() as directory:
            executor = self.make(manager, directory)
            accepted, _ = executor.submit("job-1", work)
            self.assertTrue(accepted)
            self.assertTrue(started.wait(2))
            self.assertEqual(manager.running_count(), 1)

            executor._watchdog_threshold = 0
            executor._watchdog_tick()
            self.assertEqual(manager.running_count(), 0)
            self.assertEqual(executor.active_ids(), [])

            unblock.set()
            executor.shutdown(wait=True)


if __name__ == "__main__":
    unittest.main()
