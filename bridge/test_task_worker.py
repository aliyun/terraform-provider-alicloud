#!/usr/bin/env python3
"""Focused tests for the independent durable Task worker composition."""

import os
import sys
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

sys.path.insert(0, str(Path(__file__).resolve().parent))

import task_worker


class _Executor:
    def __init__(self, *args, **kwargs):
        self.args = args
        self.kwargs = kwargs
        self.stopped = False
        self.stop_calls = []

    def start(self):
        return self

    def stop(self, *, drain, timeout):
        self.stop_calls.append((drain, timeout))
        self.stopped = True
        return True


class TaskWorkerTest(unittest.TestCase):
    def setUp(self):
        self.old_role = os.environ.get("JARVIS_BRIDGE_ROLE")
        os.environ["JARVIS_BRIDGE_ROLE"] = "worker"

    def tearDown(self):
        if self.old_role is None:
            os.environ.pop("JARVIS_BRIDGE_ROLE", None)
        else:
            os.environ["JARVIS_BRIDGE_ROLE"] = self.old_role

    def _handler(self):
        self.execute = mock.Mock()
        self.pool = mock.Mock()
        return SimpleNamespace(
            task_client=object(),
            capacity_manager=object(),
            persistent_task_execution=SimpleNamespace(execute=self.execute),
            execution_router=SimpleNamespace(task_types={"ticket", "wake"}),
            ephemeral_executor=self.pool,
        )

    def test_composes_executor_outside_jarvis_handler(self):
        handler = self._handler()
        with mock.patch.object(task_worker.bot, "jarvis_root", return_value="/repo"):
            worker = task_worker.TaskWorker(handler, executor_factory=_Executor)

        self.assertIs(worker.executor.args[0], handler.task_client)
        self.assertIs(worker.executor.args[1], handler.capacity_manager)
        self.assertIs(worker.executor.args[2], self.execute)
        self.assertEqual(worker.executor.kwargs["capabilities"], {
            "kinds": ["ticket", "wake"],
            "bridgeRole": "worker",
            "workerMode": "PERSISTENT",
            "client": "bridge",
        })

    def test_stop_owns_executor_and_handler_subprocess_cleanup(self):
        handler = self._handler()
        with mock.patch.object(task_worker.bot, "jarvis_root", return_value="/repo"), \
                mock.patch.object(task_worker.bot, "_release_claim") as release:
            worker = task_worker.TaskWorker(handler, executor_factory=_Executor)
            self.assertTrue(worker.stop(drain=True, timeout=7))

        self.assertEqual(worker.executor.stop_calls, [(True, 7)])
        handler.ephemeral_executor.terminate_all.assert_called_once_with(release_fn=release)
        handler.ephemeral_executor.shutdown.assert_called_once_with(
            wait=False, cancel_futures=True)


if __name__ == "__main__":
    unittest.main()
