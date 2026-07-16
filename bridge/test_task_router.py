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
    TaskRouter,
)


@contextmanager
def replaced_environ(values):
    """Minimal mock.patch.dict(os.environ, ..., clear=True) without unittest.mock.

    Some managed Macs block Python's optional _posixsubprocess extension, which
    unittest.mock imports indirectly through asyncio.  Environment isolation
    itself needs no such dependency.
    """
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

    def upsert_desired_task(self, env, *, execution_mode, request_id=None):
        self.calls.append((env, execution_mode, request_id))
        if self.error:
            raise self.error
        return self.response


class LegacyRecorder:
    def __init__(self, result=(True, "dispatched")):
        self.calls = 0
        self.result = result

    def __call__(self):
        self.calls += 1
        return self.result


class TaskRouterTest(unittest.TestCase):
    def test_formal_name_and_compatibility_alias_share_one_implementation(self):
        self.assertIs(TaskRouter, ExecutionRouter)

    def test_route_exposes_only_recovery_domain_kind(self):
        router = ExecutionRouter("legacy", task_types={"ticket"})
        self.assertEqual(router.route(envelope("ticket")),
                         ExecutionRoute(TASK, True))
        self.assertEqual(router.route(envelope("probe")),
                         ExecutionRoute(EPHEMERAL_JOB, False))

    def test_classify_maps_only_needs_recovery_to_domain_kind(self):
        self.assertEqual(ExecutionRouter.classify(True), (TASK, True))
        self.assertEqual(ExecutionRouter.classify(False), (EPHEMERAL_JOB, False))
        with self.assertRaises(TypeError):
            ExecutionRouter.classify("yes")

    def test_formal_migration_policy_constructor_name(self):
        router = ExecutionRouter(migration_policy="shadow", task_types={"ticket"})
        self.assertEqual(router.migration_policy, "shadow")
        self.assertEqual(router.mode, "shadow")

    def test_domain_route_is_independent_of_migration_policy(self):
        routes = {
            ExecutionRouter(policy, task_types={"ticket"}).route(envelope("ticket"))
            for policy in ("legacy", "shadow", "managed")
        }
        self.assertEqual(routes, {ExecutionRoute(TASK, True)})

    def test_is_task_is_domain_predicate_is_managed_is_compatibility_ownership(self):
        legacy = ExecutionRouter("legacy", task_types={"ticket"})
        managed = ExecutionRouter("managed", task_types={"ticket"})
        self.assertTrue(legacy.is_task(envelope("ticket")))
        self.assertFalse(legacy.is_managed(envelope("ticket")))
        self.assertTrue(managed.is_task(envelope("ticket")))
        self.assertTrue(managed.is_managed(envelope("ticket")))

    def test_legacy_never_touches_control_plane(self):
        client, legacy = FakeClient(), LegacyRecorder()
        result = TaskRouter("legacy", client=client).enqueue(envelope(), legacy)
        self.assertEqual(result, (True, "dispatched"))
        self.assertEqual(client.calls, [])
        self.assertEqual(legacy.calls, 1)

    def test_observe_legacy_is_explicit_noop(self):
        client = FakeClient()
        result = TaskRouter(
            "legacy", client=client, task_types={"ticket"}).observe(
                envelope("ticket", "modified:1"))
        self.assertEqual(result, (True, "legacy_noop"))
        self.assertEqual(client.calls, [])

    def test_shadow_task_is_observed_but_ephemeral_job_never_is(self):
        shadow = TaskRouter(
            "shadow", client=FakeClient(), task_types={"ticket"})
        result = shadow.observe(envelope("ticket", "modified:2"))
        self.assertEqual(result, (True, "shadowed"))
        self.assertEqual(shadow.client.calls[0][1], "SHADOW")

        managed = TaskRouter(
            "managed", client=FakeClient(), task_types={"probe"})
        result = managed.observe(envelope("ticket", "modified:2"))
        self.assertEqual(result, (True, "ephemeral_noop"))
        self.assertEqual(managed.client.calls, [])

    def test_observe_managed_is_fail_closed(self):
        router = TaskRouter(
            "managed", client=FakeClient(error=ControlPlaneUnavailable("down")),
            managed_task_types={"ticket"},
        )
        result = router.observe(envelope("ticket", "modified:2"))
        self.assertEqual(result, (False, "control_plane_unavailable"))

    def test_observe_refreshes_revision_for_same_task_identity(self):
        client = FakeClient()
        router = TaskRouter("shadow", client=client, task_types={"ticket"})
        router.observe(envelope("ticket", "modified:1"))
        router.observe(envelope("ticket", "modified:2"))
        self.assertEqual([call[0].task_key for call in client.calls],
                         ["aone:2100304:84345050", "aone:2100304:84345050"])
        self.assertEqual([call[0].desired_revision for call in client.calls],
                         ["modified:1", "modified:2"])
        self.assertNotEqual(client.calls[0][2], client.calls[1][2])

    def test_shadow_upserts_then_runs_legacy(self):
        client, legacy = FakeClient(), LegacyRecorder()
        result = TaskRouter(
            "shadow", client=client, task_types={"ticket"}).enqueue(
                envelope("ticket", "modified:1"), legacy)
        self.assertTrue(result.accepted)
        self.assertEqual(len(client.calls), 1, "enqueue must not duplicate observe upserts")
        self.assertEqual(client.calls[0][1], "SHADOW")
        self.assertEqual(legacy.calls, 1)

    def test_shadow_failure_does_not_affect_legacy(self):
        client = FakeClient(error=ControlPlaneUnavailable("down"))
        legacy = LegacyRecorder()
        result = TaskRouter(
            "shadow", client=client, task_types={"ticket"},
            logger=logging.getLogger("shadow-test")).enqueue(
                envelope("ticket", "modified:1"), legacy)
        self.assertEqual(result, (True, "dispatched"))
        self.assertEqual(legacy.calls, 1)

    def test_managed_allowlisted_kind_never_calls_legacy(self):
        client, legacy = FakeClient(), LegacyRecorder()
        router = TaskRouter("managed", client=client, managed_task_types={"probe"})
        result = router.enqueue(envelope("probe"), legacy)
        self.assertEqual(result, (True, "managed"))
        self.assertEqual(client.calls[0][1], "MANAGED")
        self.assertEqual(legacy.calls, 0)

    def test_managed_is_fail_closed_on_unavailable_or_rejection(self):
        for error, expected in (
            (ControlPlaneUnavailable("down"), "control_plane_unavailable"),
            (ControlPlaneConflict("busy", code="ACTIVE"), "control_plane_rejected:ACTIVE"),
            (RuntimeError("bug"), "control_plane_error"),
        ):
            legacy = LegacyRecorder()
            result = TaskRouter(
                "managed", client=FakeClient(error=error), managed_task_types={"probe"},
                logger=logging.getLogger("managed-test"),
            ).enqueue(envelope("probe"), legacy)
            self.assertEqual(result, (False, expected))
            self.assertEqual(legacy.calls, 0, "managed errors must never fall back locally")

    def test_managed_without_client_is_fail_closed(self):
        legacy = LegacyRecorder()
        result = TaskRouter("managed", managed_task_types={"probe"}).enqueue(
            envelope("probe"), legacy)
        self.assertEqual(result, (False, "control_plane_unconfigured"))
        self.assertEqual(legacy.calls, 0)

    def test_ephemeral_kind_never_upserts_and_uses_local_executor(self):
        client, legacy = FakeClient(), LegacyRecorder()
        router = TaskRouter("managed", client=client, managed_task_types={"probe"})
        result = router.enqueue(envelope("ticket", "modified:1"), legacy)
        self.assertEqual(result, (True, "dispatched"))
        self.assertEqual(client.calls, [])
        self.assertEqual(legacy.calls, 1)

    def test_empty_allowlist_manages_nothing_and_wildcard_manages_all(self):
        client, legacy = FakeClient(), LegacyRecorder()
        TaskRouter("managed", client=client).enqueue(envelope("ticket", "r1"), legacy)
        self.assertEqual(client.calls, [])
        self.assertEqual(legacy.calls, 1)

        client2, legacy2 = FakeClient(), LegacyRecorder()
        TaskRouter("managed", client=client2, managed_task_types={"*"}).enqueue(
            envelope("ticket", "r1"), legacy2)
        self.assertEqual(client2.calls[-1][1], "MANAGED")
        self.assertEqual(legacy2.calls, 0)

    def test_remote_explicit_rejection_is_not_accepted(self):
        legacy = LegacyRecorder()
        client = FakeClient(response={"accepted": False, "reason": "paused"})
        result = TaskRouter("managed", client=client, managed_task_types={"probe"}).enqueue(
            envelope(), legacy)
        self.assertEqual(result, (False, "paused"))
        self.assertEqual(legacy.calls, 0)

    def test_router_passes_stable_idempotency_request_id(self):
        client = FakeClient()
        router = TaskRouter("managed", client=client, managed_task_types={"probe"})
        router.enqueue(envelope(), LegacyRecorder())
        router.enqueue(envelope(), LegacyRecorder())
        self.assertEqual(client.calls[0][2], client.calls[1][2])

    def test_from_env_defaults_legacy_and_parses_allowlist(self):
        with replaced_environ({}):
            router = TaskRouter.from_env()
            self.assertEqual(router.mode, "legacy")
            self.assertIn("ticket", router.task_types)
            self.assertNotIn("probe", router.task_types)
        with replaced_environ({
            "JARVIS_TASK_MODE": "managed",
            "JARVIS_MANAGED_TASK_TYPES": "probe, ticket ,",
        }):
            router = TaskRouter.from_env(client=FakeClient())
            self.assertEqual(router.mode, "managed")
            self.assertEqual(router.managed_task_types, {"probe", "ticket"})

    def test_new_task_types_env_takes_precedence_over_compatibility_env(self):
        with replaced_environ({
            "JARVIS_TASK_MODE": "managed",
            "JARVIS_TASK_TYPES": "ticket",
            "JARVIS_MANAGED_TASK_TYPES": "probe",
        }):
            router = ExecutionRouter.from_env(client=FakeClient())
            self.assertEqual(router.task_types, {"ticket"})
            self.assertEqual(router.managed_task_types, {"ticket"})
            self.assertEqual(router.route(envelope("ticket")), (TASK, True))
            self.assertEqual(router.route(envelope("probe")), (EPHEMERAL_JOB, False))

    def test_new_migration_policy_env_takes_precedence_over_mode_alias(self):
        with replaced_environ({
            "JARVIS_TASK_MIGRATION_POLICY": "shadow",
            "JARVIS_TASK_MODE": "managed",
        }):
            self.assertEqual(ExecutionRouter.from_env().migration_policy, "shadow")

    def test_compatibility_allowlist_property_remains_writable(self):
        router = ExecutionRouter("managed", task_types={"ticket"})
        router.managed_task_types = {"probe"}
        self.assertEqual(router.task_types, {"probe"})
        self.assertFalse(router.is_managed(envelope("ticket")))
        self.assertTrue(router.is_managed(envelope("probe")))

    def test_rejects_invalid_mode_and_non_envelope(self):
        with self.assertRaises(ValueError):
            TaskRouter("dual-run")
        with self.assertRaises(TypeError):
            TaskRouter("legacy").enqueue({})
        with self.assertRaises(TypeError):
            ExecutionRouter("legacy").route({})
        with self.assertRaises(ValueError):
            ExecutionRouter("managed", task_types={"probe"},
                            managed_task_types={"ticket"})
        with self.assertRaises(ValueError):
            ExecutionRouter("legacy", migration_policy="managed")


if __name__ == "__main__":
    unittest.main(verbosity=2)
