#!/usr/bin/env python3
"""Hermetic tests for bridge/jarvis_task_router.py."""

import logging
import os
import sys
import unittest
from contextlib import contextmanager
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from jarvis_task_client import (  # noqa: E402
    ControlPlaneConflict,
    ControlPlaneUnavailable,
    TaskEnvelope,
)
from jarvis_task_router import (  # noqa: E402
    EPHEMERAL_JOB,
    TASK,
    ExecutionRoute,
    ExecutionRouter,
)


@contextmanager
def replaced_environ(values):
    old = dict(os.environ)
    os.environ.clear()
    os.environ.update(values)
    try:
        yield
    finally:
        os.environ.clear()
        os.environ.update(old)


def envelope(task_type="probe", revision="probe:2026-07-15"):
    key = "probe:2026-07-15" if task_type == "probe" else "aone:2100304:84345050"
    return TaskEnvelope(
        task_key=key,
        source_type="TIMER" if task_type == "probe" else "AONE",
        source_ref={"id": key},
        task_type=task_type,
        desired_revision=revision,
        trigger_mask=[task_type.upper()],
        payload={"kind": task_type, "prompt": "work"},
    )


class FakeClient:
    def __init__(self, response=None, error=None):
        self.response = response if response is not None else {"accepted": True}
        self.error = error
        self.calls = []

    def upsert_desired_task(self, env, *, request_id=None):
        self.calls.append((env, request_id))
        if self.error:
            raise self.error
        return self.response


class LocalRecorder:
    def __init__(self, result=(True, "dispatched")):
        self.calls = 0
        self.result = result

    def __call__(self):
        self.calls += 1
        return self.result


class ExecutionRouterTest(unittest.TestCase):
    def test_router_surface_contains_only_current_task_model(self):
        router = ExecutionRouter(task_types={"ticket"})
        self.assertEqual(set(vars(router)), {"client", "task_types", "log"})

    def test_classify_maps_only_needs_recovery_to_domain_kind(self):
        self.assertEqual(ExecutionRouter.classify(True),
                         ExecutionRoute(TASK, True))
        self.assertEqual(ExecutionRouter.classify(False),
                         ExecutionRoute(EPHEMERAL_JOB, False))
        with self.assertRaises(TypeError):
            ExecutionRouter.classify("yes")

    def test_route_uses_recovery_task_types_and_wildcard(self):
        router = ExecutionRouter(task_types={"ticket"})
        self.assertEqual(router.route(envelope("ticket")), (TASK, True))
        self.assertEqual(router.route(envelope("probe")),
                         (EPHEMERAL_JOB, False))
        self.assertEqual(
            ExecutionRouter(task_types={"*"}).route(envelope("probe")),
            (TASK, True),
        )

    def test_default_task_classification_is_fixed(self):
        with replaced_environ({"JARVIS_UNUSED_SETTING": "probe"}):
            router = ExecutionRouter(client=FakeClient())
        self.assertIn("ticket", router.task_types)
        self.assertIn("wake", router.task_types)
        self.assertNotIn("probe", router.task_types)
        self.assertEqual(router.route(envelope("probe")),
                         (EPHEMERAL_JOB, False))

    def test_task_observe_always_persists(self):
        client = FakeClient()
        router = ExecutionRouter(client=client, task_types={"ticket"})
        result = router.observe(envelope("ticket", "modified:2"))
        self.assertEqual(result, (True, "task_persisted"))
        self.assertEqual(len(client.calls), 1)

    def test_task_enqueue_persists_and_never_runs_local_submit(self):
        client, local = FakeClient(), LocalRecorder()
        router = ExecutionRouter(client=client, task_types={"ticket"})
        result = router.enqueue(envelope("ticket"), local)
        self.assertEqual(result, (True, "task_persisted"))
        self.assertEqual(len(client.calls), 1)
        self.assertEqual(local.calls, 0)

    def test_task_is_fail_closed_on_unconfigured_control_plane(self):
        local = LocalRecorder()
        result = ExecutionRouter(task_types={"ticket"}).enqueue(
            envelope("ticket"), local)
        self.assertEqual(result, (False, "control_plane_unconfigured"))
        self.assertEqual(local.calls, 0)

    def test_task_is_fail_closed_on_unavailable_rejection_or_error(self):
        for error, expected in (
            (ControlPlaneUnavailable("down"), "control_plane_unavailable"),
            (ControlPlaneConflict("busy", code="ACTIVE"),
             "control_plane_rejected:ACTIVE"),
            (RuntimeError("bug"), "control_plane_error"),
        ):
            local = LocalRecorder()
            result = ExecutionRouter(
                client=FakeClient(error=error),
                task_types={"ticket"},
                logger=logging.getLogger("task-test"),
            ).enqueue(envelope("ticket"), local)
            self.assertEqual(result, (False, expected))
            self.assertEqual(local.calls, 0)

    def test_task_explicit_remote_rejection_is_fail_closed(self):
        client = FakeClient(response={"accepted": False, "reason": "paused"})
        local = LocalRecorder()
        result = ExecutionRouter(
            client=client, task_types={"ticket"}).enqueue(
                envelope("ticket"), local)
        self.assertEqual(result, (False, "paused"))
        self.assertEqual(local.calls, 0)

    def test_task_request_id_is_stable(self):
        client = FakeClient()
        router = ExecutionRouter(client=client, task_types={"ticket"})
        router.enqueue(envelope("ticket"), LocalRecorder())
        router.enqueue(envelope("ticket"), LocalRecorder())
        self.assertEqual(client.calls[0][1], client.calls[1][1])

    def test_ephemeral_observe_never_touches_control_plane(self):
        client = FakeClient()
        result = ExecutionRouter(
            client=client, task_types={"ticket"}).observe(envelope("probe"))
        self.assertEqual(result, (True, "ephemeral_noop"))
        self.assertEqual(client.calls, [])

    def test_ephemeral_enqueue_only_runs_local_submit(self):
        client, local = FakeClient(), LocalRecorder()
        result = ExecutionRouter(
            client=client, task_types={"ticket"}).enqueue(
                envelope("probe"), local)
        self.assertEqual(result, (True, "dispatched"))
        self.assertEqual(local.calls, 1)
        self.assertEqual(client.calls, [])

    def test_ephemeral_enqueue_requires_local_submit(self):
        result = ExecutionRouter().enqueue(envelope("probe"))
        self.assertEqual(result, (False, "local_submit_missing"))

    def test_rejects_unknown_constructor_options_and_non_envelope(self):
        with self.assertRaises(TypeError):
            ExecutionRouter(unexpected_option=True)
        with self.assertRaises(TypeError):
            ExecutionRouter().route({})
        with self.assertRaises(TypeError):
            ExecutionRouter().enqueue({})


if __name__ == "__main__":
    unittest.main(verbosity=2)
