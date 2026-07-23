#!/usr/bin/env python3
"""Hermetic tests for bridge/jarvis_capacity.py."""

import sys
import threading
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from bridge.jarvis_capacity import CapacityManager, CapacitySnapshot  # noqa: E402


class CapacityManagerTest(unittest.TestCase):
    def test_rejects_invalid_capacity_and_allows_zero(self):
        for value in (-1, True, 1.5, "not-an-int", None):
            with self.subTest(value=value), self.assertRaises(ValueError):
                CapacityManager(value)

        manager = CapacityManager(0)
        self.assertIsNone(manager.acquire("disabled"))
        self.assertEqual(manager.snapshot(), CapacitySnapshot(0, 0, 0))

    def test_acquire_is_bounded_and_snapshot_is_consistent(self):
        manager = CapacityManager(2)

        first = manager.acquire("task:t1")
        second = manager.acquire("ephemeral:j1")

        self.assertIsNotNone(first)
        self.assertIsNotNone(second)
        self.assertEqual(first.owner, "task:t1")
        self.assertEqual(second.owner, "ephemeral:j1")
        self.assertIsNone(manager.acquire("task:t2"))
        self.assertEqual(manager.snapshot(), CapacitySnapshot(2, 2, 0))
        self.assertEqual(manager.running_count(), 2)
        self.assertEqual(manager.available_slots(), 0)

    def test_release_is_idempotent_and_reopens_capacity(self):
        manager = CapacityManager(1)
        permit = manager.acquire("task:t1")
        self.assertIsNotNone(permit)
        self.assertFalse(permit.released)

        self.assertTrue(permit.release())
        self.assertTrue(permit.released)
        self.assertFalse(permit.release())
        self.assertEqual(manager.snapshot(), CapacitySnapshot(1, 0, 1))

        replacement = manager.acquire("ephemeral:j1")
        self.assertIsNotNone(replacement)
        self.assertTrue(manager.release(replacement))
        self.assertFalse(manager.release(replacement))

    def test_context_manager_releases_after_success_and_exception(self):
        manager = CapacityManager(1)
        with manager.acquire("task:t1") as permit:
            self.assertFalse(permit.released)
            self.assertEqual(manager.running_count(), 1)
        self.assertTrue(permit.released)
        self.assertEqual(manager.running_count(), 0)

        with self.assertRaisesRegex(RuntimeError, "boom"):
            with manager.acquire("ephemeral:j1"):
                raise RuntimeError("boom")
        self.assertEqual(manager.snapshot(), CapacitySnapshot(1, 0, 1))

    def test_manager_rejects_foreign_permit(self):
        first = CapacityManager(1)
        second = CapacityManager(1)
        permit = first.acquire("task:t1")

        with self.assertRaisesRegex(ValueError, "does not belong"):
            second.release(permit)
        self.assertEqual(first.running_count(), 1)
        self.assertTrue(permit.release())

    def test_concurrent_release_is_idempotent(self):
        manager = CapacityManager(1)
        permit = manager.acquire("task:t1")
        start = threading.Barrier(9)
        result_lock = threading.Lock()
        results = []

        def release():
            start.wait()
            result = permit.release()
            with result_lock:
                results.append(result)

        threads = [threading.Thread(target=release) for _ in range(8)]
        for thread in threads:
            thread.start()
        start.wait()
        for thread in threads:
            thread.join(timeout=5)
            self.assertFalse(thread.is_alive())

        self.assertEqual(results.count(True), 1)
        self.assertEqual(results.count(False), 7)
        self.assertEqual(manager.snapshot(), CapacitySnapshot(1, 0, 1))

    def test_concurrent_acquire_never_oversells(self):
        capacity = 3
        worker_count = 40
        manager = CapacityManager(capacity)
        start = threading.Barrier(worker_count + 1)
        release = threading.Event()
        all_attempted = threading.Event()
        result_lock = threading.Lock()
        results = []

        def worker(index):
            start.wait()
            permit = manager.acquire("worker:%d" % index)
            with result_lock:
                results.append(permit is not None)
                if len(results) == worker_count:
                    all_attempted.set()
            if permit is not None:
                release.wait(timeout=5)
                permit.release()

        threads = [
            threading.Thread(target=worker, args=(index,))
            for index in range(worker_count)
        ]
        for thread in threads:
            thread.start()
        start.wait()
        self.assertTrue(all_attempted.wait(timeout=5))

        self.assertEqual(sum(results), capacity)
        self.assertEqual(manager.snapshot(),
                         CapacitySnapshot(capacity, capacity, 0))

        release.set()
        for thread in threads:
            thread.join(timeout=5)
            self.assertFalse(thread.is_alive())
        self.assertEqual(manager.snapshot(),
                         CapacitySnapshot(capacity, 0, capacity))


if __name__ == "__main__":
    unittest.main()
