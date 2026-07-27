#!/usr/bin/env python3
"""Focused tests for the independent durable Persistent Worker composition."""

import os
import sys
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

sys.path.insert(0, str(Path(__file__).resolve().parent))

from bridge import persistent_worker


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


class PersistentWorkerTest(unittest.TestCase):
    def setUp(self):
        self.old_role = os.environ.get("JARVIS_BRIDGE_ROLE")
        os.environ["JARVIS_BRIDGE_ROLE"] = "worker"

    def tearDown(self):
        if self.old_role is None:
            os.environ.pop("JARVIS_BRIDGE_ROLE", None)
        else:
            os.environ["JARVIS_BRIDGE_ROLE"] = self.old_role

    def _runtime(self):
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
        runtime = self._runtime()
        worker = persistent_worker.PersistentWorker(
            runtime, executor_factory=_Executor)

        self.assertIs(worker.executor.args[0], runtime.task_client)
        self.assertIs(worker.executor.args[1], runtime.capacity_manager)
        self.assertIs(worker.executor.args[2], self.execute)
        self.assertEqual(worker.executor.kwargs["capabilities"], {
            "kinds": ["ticket", "wake"],
            "bridgeRole": "worker",
            "workerMode": "PERSISTENT",
            "client": "bridge",
        })
        beacon_paths = worker.executor.kwargs["heartbeat_beacon_paths"]
        self.assertEqual(set(beacon_paths), {"worker", "lease", "session"})
        self.assertTrue(all(
            str(path).endswith(".%s.epoch" % name)
            for name, path in beacon_paths.items()))

    def test_stop_owns_executor_and_runtime_subprocess_cleanup(self):
        runtime = self._runtime()
        release = mock.Mock()
        worker = persistent_worker.PersistentWorker(
            runtime, executor_factory=_Executor, release_claim=release)
        self.assertTrue(worker.stop(drain=True, timeout=7))

        self.assertEqual(worker.executor.stop_calls, [(True, 7)])
        runtime.ephemeral_executor.terminate_all.assert_called_once_with(release_fn=release)
        runtime.ephemeral_executor.shutdown.assert_called_once_with(
            wait=False, cancel_futures=True)

    def test_runtime_composition_does_not_construct_bot_or_schedulers(self):
        client = object()
        pool = mock.Mock()
        pool_factory = mock.Mock(return_value=pool)
        field_worker = mock.Mock()
        field_factory = mock.Mock(return_value=field_worker)
        with mock.patch.object(
                persistent_worker, "claude_bin", return_value="/bin/claude"):
            runtime = persistent_worker.PersistentTaskRuntime(
                task_client=client,
                field_repair_worker_factory=field_factory,
                ephemeral_executor_factory=pool_factory,
            )

        self.assertIs(runtime.task_client, client)
        self.assertIs(runtime.field_repair_worker, field_worker)
        self.assertIs(runtime.ephemeral_executor, pool)
        self.assertFalse(hasattr(runtime, "scanner"))
        self.assertFalse(hasattr(runtime, "daily"))
        self.assertFalse(hasattr(runtime, "prwatch"))
        self.assertEqual(
            runtime.execution_router.task_types,
            {
                "ticket", "revisit", "persona", "wake", "pr_ci_fix",
                "pr_comment_reply",
            },
        )
        self.assertEqual(runtime.capacity_manager.capacity, 3)
        self.assertIs(
            runtime.persistent_task_execution._field_repair_worker,
            field_worker,
        )
        self.assertIs(
            runtime.persistent_task_execution._dispatch_item.__self__,
            runtime,
        )


if __name__ == "__main__":
    unittest.main()
