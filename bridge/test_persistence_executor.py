#!/usr/bin/env python3
"""Hermetic tests for bridge/jarvis_persistence_executor.py."""

import logging
import sys
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

from jarvis_persistence_executor import (  # noqa: E402
    LeaseProtocolError,
    SessionController,
    PersistenceExecutor,
    make_worker_key,
    parse_lease_response,
)
from jarvis_capacity import CapacityManager  # noqa: E402
from jarvis_task_client import (  # noqa: E402
    ControlPlaneConflict,
    ControlPlaneUnavailable,
    StaleFence,
)


LOG = logging.getLogger("local-worker-test")
LOG.handlers = [logging.NullHandler()]
LOG.propagate = False


def lease_response(session_id="s1", fence_token=7, resumed=False,
                   runtime_session_id=None):
    response = {
        "task": {"id": "t1", "taskType": "ticket"},
        "session": {"id": session_id, "fenceToken": fence_token},
        "resumed": resumed,
    }
    if runtime_session_id is not None:
        response["session"]["runtimeSessionId"] = runtime_session_id
    return response


class FakeClient:
    def __init__(self, **outcomes):
        self.outcomes = {name: list(values) for name, values in outcomes.items()}
        self.calls = []

    def _call(self, name, *args, **kwargs):
        self.calls.append({"name": name, "args": args, "kwargs": kwargs})
        values = self.outcomes.get(name)
        if values:
            value = values.pop(0)
        elif name == "lease_task":
            value = {}
        else:
            value = {"accepted": True}
        if isinstance(value, BaseException):
            raise value
        return value

    def named(self, name):
        return [call for call in self.calls if call["name"] == name]

    def register_worker(self, *args, **kwargs):
        return self._call("register_worker", *args, **kwargs)

    def heartbeat_worker(self, *args, **kwargs):
        return self._call("heartbeat_worker", *args, **kwargs)

    def lease_task(self, *args, **kwargs):
        return self._call("lease_task", *args, **kwargs)

    def start_session(self, *args, **kwargs):
        return self._call("start_session", *args, **kwargs)

    def heartbeat_session(self, *args, **kwargs):
        return self._call("heartbeat_session", *args, **kwargs)

    def complete_session(self, *args, **kwargs):
        return self._call("complete_session", *args, **kwargs)

    def fail_session(self, *args, **kwargs):
        return self._call("fail_session", *args, **kwargs)

    def suspend_session(self, *args, **kwargs):
        return self._call("suspend_session", *args, **kwargs)


class ManualFuture:
    def __init__(self, fn, args):
        self.fn = fn
        self.args = args
        self._done = False
        self.error = None

    def run(self):
        try:
            self.fn(*self.args)
        except BaseException as exc:
            self.error = exc
        self._done = True

    def done(self):
        return self._done


class ManualExecutor:
    def __init__(self, inline=False):
        self.inline = inline
        self.futures = []

    def submit(self, fn, *args):
        future = ManualFuture(fn, args)
        self.futures.append(future)
        if self.inline:
            future.run()
        return future

    def run_all(self):
        for future in list(self.futures):
            if not future.done():
                future.run()


class FakeClock:
    def __init__(self):
        self.now = 0.0

    def __call__(self):
        return self.now

    def advance(self, seconds):
        self.now += seconds


class FakeProcess:
    def __init__(self, pid):
        self.pid = pid


class WorkerKeyAndLeaseTest(unittest.TestCase):
    def test_worker_key_is_stable_and_injectable(self):
        first = make_worker_key("mac-mini", "boot-1", "process-1")
        second = make_worker_key("mac-mini", "boot-1", "process-1")
        self.assertEqual(first, "mac-mini:boot-1:process-1")
        self.assertEqual(first, second)

    def test_frozen_direct_lease_empty_204_and_draft_wrapper(self):
        direct = lease_response(resumed=True)
        self.assertEqual(parse_lease_response(direct)["resumed"], True)
        self.assertIsNone(parse_lease_response({}))
        self.assertIsNone(parse_lease_response({"lease": None}))
        wrapped = parse_lease_response({"lease": lease_response()})
        self.assertEqual(wrapped["session"]["id"], "s1")
        with self.assertRaises(LeaseProtocolError):
            parse_lease_response({"task": {"id": "t1"}})


class SessionControllerTest(unittest.TestCase):
    def make(self, client=None, **kwargs):
        return SessionController(
            client or FakeClient(), "mac:boot:proc", lease_response(),
            lease_seconds=45, lease_safety_margin=10, logger=LOG, **kwargs)

    def test_lease_safety_margin_must_be_smaller_than_ttl(self):
        with self.assertRaisesRegex(ValueError, "less than lease_seconds"):
            SessionController(
                FakeClient(), "mac:boot:proc", lease_response(),
                lease_seconds=300, lease_safety_margin=300, logger=LOG)

    def test_start_and_heartbeat_renew_the_fenced_lease(self):
        client = FakeClient()
        controller = self.make(client, runtime_session_id_factory=lambda: "runtime-new")
        self.assertTrue(controller.start())
        self.assertTrue(controller.heartbeat())
        self.assertTrue(controller.heartbeat())
        self.assertEqual(client.named("start_session")[0]["args"][3],
                         {"leaseSeconds": 45, "runtimeSessionId": "runtime-new"})
        self.assertEqual(client.named("heartbeat_session")[0]["args"][3],
                         {"leaseSeconds": 45})
        heartbeats = client.named("heartbeat_session")
        self.assertNotEqual(heartbeats[0]["kwargs"]["request_id"],
                            heartbeats[1]["kwargs"]["request_id"])
        self.assertEqual(controller.lease["session"]["runtimeSessionId"],
                         "runtime-new")

    def test_resumed_lease_reuses_persisted_runtime_session_id(self):
        lease = lease_response(resumed=True, runtime_session_id="runtime-existing")

        def must_not_generate():
            raise AssertionError("resumed lease must not mint another runtime session")

        controller = SessionController(
            FakeClient(), "mac:boot:proc", lease, lease_seconds=45,
            lease_safety_margin=10,
            runtime_session_id_factory=must_not_generate, logger=LOG)
        self.assertTrue(controller.resumed)
        self.assertEqual(controller.runtime_session_id, "runtime-existing")
        self.assertTrue(controller.start())

    def test_terminal_network_failure_is_pending_then_retried_idempotently(self):
        client = FakeClient(complete_session=[
            ControlPlaneUnavailable("offline"), {"accepted": True},
        ])
        failures = []
        lifecycle = self.make(client, on_network_failure=failures.append)
        self.assertTrue(lifecycle.start())
        self.assertFalse(lifecycle.complete({"answer": 42}))
        self.assertEqual(lifecycle.pending_terminal, "complete")
        self.assertFalse(lifecycle.terminal)
        self.assertEqual(len(failures), 1)

        self.assertTrue(lifecycle.retry_terminal())
        self.assertEqual(lifecycle.terminal_action, "complete")
        calls = client.named("complete_session")
        self.assertEqual(calls[0]["args"][3], {"result": {"answer": 42}})
        self.assertEqual(calls[0]["kwargs"]["request_id"],
                         calls[1]["kwargs"]["request_id"])

    def test_fail_uses_error_and_retry_after_contract(self):
        client = FakeClient()
        lifecycle = self.make(client)
        lifecycle.start()
        self.assertTrue(lifecycle.fail({"type": "TRANSIENT"}, retry_after_seconds=30))
        payload = client.named("fail_session")[0]["args"][3]
        self.assertEqual(payload, {
            "error": {"type": "TRANSIENT"},
            "retryAfterSeconds": 30,
        })

    def test_stale_fence_stops_the_bound_process_without_terminal_ack(self):
        client = FakeClient(heartbeat_session=[StaleFence("old")])
        stopped = []
        lifecycle = self.make(
            client, stop_process=lambda current, reason: stopped.append(
                (current.process, reason)))
        lifecycle.start()
        process = FakeProcess(4321)
        lifecycle.bind_process(process)
        self.assertFalse(lifecycle.heartbeat())
        self.assertTrue(lifecycle.ownership_lost)
        self.assertEqual(stopped, [(process, "stale_fence:heartbeat")])
        self.assertEqual(client.named("complete_session"), [])
        self.assertEqual(client.named("fail_session"), [])

    def test_transient_heartbeat_503_keeps_running_then_renews(self):
        clock = FakeClock()
        client = FakeClient(heartbeat_session=[
            ControlPlaneUnavailable("deploying"),
            {"leaseExpireAt": "2026-07-17T06:00:00Z"},
        ])
        stopped = []
        lifecycle = SessionController(
            client, "mac:boot:proc", lease_response(),
            lease_seconds=300, lease_safety_margin=90, clock=clock,
            stop_process=lambda current, reason: stopped.append(
                (current.process, reason)), logger=LOG)
        self.assertTrue(lifecycle.start())
        process = FakeProcess(4321)
        lifecycle.bind_process(process)
        first_deadline = lifecycle.lease_deadline

        clock.advance(30)
        self.assertTrue(lifecycle.heartbeat())
        self.assertFalse(lifecycle.ownership_lost)
        self.assertEqual(stopped, [])
        self.assertEqual(lifecycle.lease_deadline, first_deadline)

        clock.advance(30)
        self.assertTrue(lifecycle.heartbeat())
        self.assertFalse(lifecycle.ownership_lost)
        self.assertEqual(lifecycle.lease_deadline, 360)
        self.assertEqual(lifecycle.session["leaseExpireAt"],
                         "2026-07-17T06:00:00Z")

    def test_sustained_heartbeat_503_stops_at_safety_boundary(self):
        clock = FakeClock()
        client = FakeClient(heartbeat_session=[
            ControlPlaneUnavailable("deploying") for _ in range(7)
        ])
        stopped = []
        lifecycle = SessionController(
            client, "mac:boot:proc", lease_response(),
            lease_seconds=300, lease_safety_margin=90, clock=clock,
            stop_process=lambda current, reason: stopped.append(
                (current.process, reason)), logger=LOG)
        self.assertTrue(lifecycle.start())
        process = FakeProcess(4321)
        lifecycle.bind_process(process)

        for _ in range(6):
            clock.advance(30)
            self.assertTrue(lifecycle.heartbeat())
        self.assertEqual(clock(), 180)
        self.assertEqual(stopped, [])

        clock.advance(30)
        self.assertFalse(lifecycle.heartbeat())
        self.assertTrue(lifecycle.ownership_lost)
        self.assertEqual(stopped, [
            (process, "lease_proof_expiring:heartbeat")])

    def test_bind_process_persists_pid_with_a_distinct_fenced_start(self):
        client = FakeClient()
        lifecycle = self.make(client, runtime_session_id_factory=lambda: "runtime-new")
        self.assertTrue(lifecycle.start())
        process = FakeProcess(4321)

        self.assertIs(lifecycle.bind_process(process), process)

        starts = client.named("start_session")
        self.assertEqual(len(starts), 2)
        self.assertEqual(starts[1]["args"][0:3], ("s1", "mac:boot:proc", 7))
        self.assertEqual(starts[1]["args"][3], {
            "pid": 4321,
            "leaseSeconds": 45,
            "runtimeSessionId": "runtime-new",
        })
        self.assertNotEqual(starts[0]["kwargs"]["request_id"],
                            starts[1]["kwargs"]["request_id"])
        self.assertIn("bind-process-4321", starts[1]["kwargs"]["request_id"])
        self.assertEqual(lifecycle.session["pid"], 4321)

    def test_unavailable_or_rejected_process_binding_stops_process_fail_closed(self):
        for failure, reason in (
                (ControlPlaneUnavailable("offline"),
                 "control_plane_unavailable:bind_process"),
                (ControlPlaneConflict("rejected"), "process_bind_rejected"),
                (StaleFence("old"), "stale_fence:bind_process")):
            with self.subTest(failure=type(failure).__name__):
                client = FakeClient(start_session=[{"accepted": True}, failure])
                stopped = []
                failures = []
                lifecycle = self.make(
                    client,
                    stop_process=lambda current, why: stopped.append(
                        (current.process, why)),
                    on_network_failure=failures.append)
                self.assertTrue(lifecycle.start())
                process = FakeProcess(4321)

                with self.assertRaises(type(failure)):
                    lifecycle.bind_process(process)

                self.assertTrue(lifecycle.ownership_lost)
                self.assertTrue(lifecycle.stop_requested)
                self.assertEqual(stopped, [(process, reason)])
                self.assertEqual(failures, [failure] if isinstance(
                    failure, ControlPlaneUnavailable) else [])
                self.assertNotIn("pid", lifecycle.session)

    def test_process_bound_after_fence_loss_is_stopped_synchronously(self):
        stopped = []
        lifecycle = self.make(
            FakeClient(),
            stop_process=lambda current, reason: stopped.append(
                (current.process, reason)))

        lifecycle._lose_ownership("stale_fence:heartbeat")
        self.assertEqual(stopped, [(None, "stale_fence:heartbeat")])

        lifecycle.bind_process("late-process")

        self.assertEqual(stopped, [
            (None, "stale_fence:heartbeat"),
            ("late-process", "stale_fence:heartbeat"),
        ])
        self.assertTrue(lifecycle.ownership_lost)
        self.assertTrue(lifecycle.stop_requested)


class PersistenceExecutorTest(unittest.TestCase):
    def make(self, client, execute, stop_process, *, executor=None, clock=None,
             capacity=2, capacity_manager=None, **kwargs):
        manager = capacity_manager or CapacityManager(capacity)
        lease_safety_margin = kwargs.pop("lease_safety_margin", 10)
        return PersistenceExecutor(
            client, manager, execute, stop_process,
            worker_key="mac:boot:proc", lease_seconds=45,
            lease_safety_margin=lease_safety_margin,
            executor=executor or ManualExecutor(inline=True),
            clock=clock or FakeClock(), logger=LOG, **kwargs)

    def test_shared_capacity_is_reserved_before_lease_and_released_after_run(self):
        manager = CapacityManager(1)
        ephemeral = manager.acquire("ephemeral:busy")
        client = FakeClient(lease_task=[lease_response()])
        executor = ManualExecutor()
        worker = self.make(
            client, lambda *_args: "done", lambda *_args: None,
            capacity_manager=manager, executor=executor)

        self.assertTrue(worker.run_once())
        self.assertEqual(client.named("lease_task"), [],
                         "PersistenceExecutor must reserve capacity before leasing")

        ephemeral.release()
        self.assertTrue(worker.run_once())
        self.assertEqual(manager.running_count(), 1)
        self.assertEqual(worker.active_session_ids(), ["s1"])
        executor.run_all()
        self.assertEqual(manager.running_count(), 0)
        self.assertEqual(worker.active_count(), 0)

    def test_register_lease_start_execute_and_complete(self):
        client = FakeClient(lease_task=[lease_response(resumed=True)])
        executed = []

        def execute(lease, lifecycle):
            executed.append((lease, lifecycle))
            return {"status": "done", "value": 42}

        worker = self.make(
            client, execute, lambda *_args: None,
            runtime_session_id_factory=lambda: "runtime-new")
        self.assertTrue(worker.run_once())
        self.assertEqual(worker.active_count(), 0)
        self.assertEqual(len(executed), 1)
        self.assertTrue(executed[0][0]["resumed"])
        self.assertTrue(executed[0][1].started)
        self.assertEqual(executed[0][1].runtime_session_id, "runtime-new")

        register = client.named("register_worker")[0]
        self.assertEqual(register["args"][0]["workerKey"], "mac:boot:proc")
        self.assertEqual(register["args"][0]["maxSlots"], 2)
        lease_call = client.named("lease_task")[0]
        self.assertNotIn("free_slots", lease_call["kwargs"])
        self.assertEqual(lease_call["kwargs"]["lease_seconds"], 45)
        start_payload = client.named("start_session")[0]["args"][3]
        self.assertEqual(start_payload["leaseSeconds"], 45)
        self.assertEqual(start_payload["runtimeSessionId"], "runtime-new")
        completed = client.named("complete_session")[0]["args"][3]
        self.assertEqual(completed, {
            "result": {"status": "done", "value": 42},
        })

    def test_execute_result_maps_suspend_and_fail(self):
        failed = {"status": "failed", "error": {"type": "TRANSIENT"},
                  "retryAfterSeconds": 20}
        for result, expected in (("suspended", "suspend_session"),
                                 (failed, "fail_session")):
            with self.subTest(result=result):
                client = FakeClient(lease_task=[lease_response()])
                worker = self.make(client, lambda *_args: result, lambda *_args: None)
                self.assertTrue(worker.run_once())
                self.assertEqual(len(client.named(expected)), 1)
                self.assertEqual(len(client.named("complete_session")), 0)
                if expected == "fail_session":
                    self.assertEqual(client.named(expected)[0]["args"][3], {
                        "error": {"type": "TRANSIENT"},
                        "retryAfterSeconds": 20,
                    })

    def test_managed_wait_state_is_forwarded_without_runner_status_fields(self):
        client = FakeClient(lease_task=[lease_response()])
        result = {
            "status": "suspended",
            "waitType": "AONE_REPLY",
            "waitKey": "843",
            "waitCursor": "41",
            "waitExpireAt": "2026-07-17T00:00:00Z",
        }
        worker = self.make(client, lambda *_args: result, lambda *_args: None)
        self.assertTrue(worker.run_once())
        self.assertEqual(client.named("suspend_session")[0]["args"][3], {
            "waitType": "AONE_REPLY",
            "waitKey": "843",
            "waitCursor": "41",
            "waitExpireAt": "2026-07-17T00:00:00Z",
        })

    def test_session_stale_fence_stops_only_its_lease_immediately(self):
        client = FakeClient(
            lease_task=[lease_response()],
            heartbeat_session=[{}, StaleFence("old")])
        executor = ManualExecutor()
        executed, stopped = [], []
        worker = self.make(
            client, lambda *_args: executed.append(True),
            lambda lifecycle, reason: stopped.append(
                (lifecycle.session_id, reason)), executor=executor)
        self.assertTrue(worker.run_once())
        self.assertEqual(worker.active_session_ids(), ["s1"])

        worker.heartbeat_sessions_once()
        self.assertEqual(stopped, [])
        worker.heartbeat_sessions_once()
        self.assertEqual(stopped, [("s1", "stale_fence:heartbeat")])
        self.assertEqual(executed, [])
        terminal_names = {"complete_session", "fail_session", "suspend_session"}
        self.assertFalse(any(call["name"] in terminal_names
                             for call in client.calls))

        executor.run_all()
        self.assertEqual(executed, [], "lost queued work must never start")
        self.assertEqual(worker.active_count(), 0)

    def test_network_failure_pauses_leasing_until_register_recovers(self):
        client = FakeClient(
            register_worker=[{}, ControlPlaneUnavailable("still offline"), {}],
            lease_task=[ControlPlaneUnavailable("offline"), {}],
        )
        executed = []
        worker = self.make(client, lambda *_args: executed.append(True),
                           lambda *_args: None)

        self.assertFalse(worker.run_once())
        self.assertFalse(worker.network_healthy)
        self.assertEqual(len(client.named("lease_task")), 1)
        self.assertFalse(worker.run_once())
        self.assertEqual(len(client.named("lease_task")), 1,
                         "failed recovery registration must not lease")
        self.assertTrue(worker.run_once())
        self.assertTrue(worker.network_healthy)
        self.assertEqual(len(client.named("lease_task")), 2)
        self.assertEqual(executed, [])
        self.assertFalse(any(call["name"] in {
            "complete_session", "fail_session", "suspend_session"
        } for call in client.calls))

    def test_pending_terminal_is_retried_before_another_lease(self):
        client = FakeClient(
            lease_task=[lease_response("s1"), lease_response("s2")],
            complete_session=[ControlPlaneUnavailable("offline"), {}, {}],
        )
        worker = self.make(client, lambda *_args: "done", lambda *_args: None,
                           capacity=2)
        self.assertTrue(worker.run_once())
        self.assertFalse(worker.network_healthy)
        self.assertEqual(worker.active_session_ids(), ["s1"])

        self.assertTrue(worker.run_once())
        self.assertEqual(len(client.named("lease_task")), 1,
                         "pending terminal must block new admission")
        worker.heartbeat_sessions_once()
        self.assertEqual(worker.active_count(), 0)
        self.assertTrue(worker.run_once())
        self.assertEqual(len(client.named("lease_task")), 2)

    def test_capacity_manager_is_the_local_gate_and_intervals_are_injectable(self):
        clock = FakeClock()
        client = FakeClient(lease_task=[{}, {}, {}])
        manager = CapacityManager(3)
        worker = self.make(
            client, lambda *_args: None, lambda *_args: None,
            capacity_manager=manager, clock=clock, worker_heartbeat_interval=10)
        worker.run_once()
        worker.run_once()
        self.assertEqual(len(client.named("heartbeat_worker")), 1)
        clock.advance(11)
        worker.run_once()
        self.assertEqual(len(client.named("heartbeat_worker")), 2)
        self.assertEqual(len(client.named("lease_task")), 3)
        self.assertTrue(all("free_slots" not in call["kwargs"]
                            for call in client.named("lease_task")))
        self.assertEqual(client.named("register_worker")[0]["args"][0]["maxSlots"], 3)

        no_capacity = FakeClient()
        blocked_manager = CapacityManager(1)
        ephemeral = blocked_manager.acquire("ephemeral:busy")
        gated = self.make(
            no_capacity, lambda *_args: None, lambda *_args: None,
            capacity_manager=blocked_manager)
        self.assertTrue(gated.run_once())
        self.assertEqual(no_capacity.named("lease_task"), [])
        ephemeral.release()

    def test_task_and_ephemeral_execution_share_one_capacity_manager(self):
        manager = CapacityManager(3)
        ephemeral = manager.acquire("ephemeral:busy")
        client = FakeClient(lease_task=[
            lease_response("s1"), lease_response("s2"), lease_response("s3")])
        worker = self.make(
            client, lambda *_args: "done", lambda *_args: None,
            capacity_manager=manager, executor=ManualExecutor())
        self.assertTrue(worker.run_once())
        self.assertTrue(worker.run_once())
        self.assertTrue(worker.run_once())
        self.assertEqual(worker.active_count(), 2)
        self.assertEqual(len(client.named("lease_task")), 2,
                         "Task plus EphemeralJob work must stay within maxSlots")
        self.assertEqual(client.named("register_worker")[0]["args"][0]["maxSlots"], 3)
        ephemeral.release()

    def test_backend_capacity_conflict_is_not_a_network_failure(self):
        client = FakeClient(lease_task=[
            ControlPlaneConflict("full", code="WORKER_FULL"), {},
        ])
        worker = self.make(client, lambda *_args: None, lambda *_args: None)
        self.assertTrue(worker.run_once())
        self.assertTrue(worker.network_healthy)
        self.assertTrue(worker.run_once())
        self.assertEqual(len(client.named("register_worker")), 1)
        self.assertEqual(len(client.named("lease_task")), 2)

    def test_drain_stops_new_leases_and_immediate_stop_retries_active(self):
        client = FakeClient(lease_task=[lease_response(), lease_response("s2")])
        executor = ManualExecutor()
        stopped = []

        def stop_before_release(lifecycle, reason):
            stopped.append((lifecycle.session_id, reason,
                            len(client.named("fail_session"))))

        worker = self.make(
            client, lambda *_args: "done",
            stop_before_release,
            executor=executor)
        self.assertTrue(worker.run_once())
        self.assertFalse(worker.drain(timeout=0))
        self.assertTrue(worker.run_once())
        self.assertEqual(len(client.named("lease_task")), 1)

        self.assertFalse(worker.stop())
        self.assertTrue(worker.stopped)
        self.assertEqual(len(client.named("suspend_session")), 0)
        self.assertEqual(client.named("fail_session")[0]["args"][3], {
            "error": {
                "errorType": "WorkerStopping",
                "message": "worker_stopping",
            },
            "retryAfterSeconds": 0,
        })
        self.assertEqual(stopped, [("s1", "worker_stopping", 0)],
                         "the process must stop before the task becomes retryable")
        executor.run_all()
        self.assertEqual(worker.active_count(), 0)

    def test_lease_arriving_after_drain_is_failed_retryable_not_suspended(self):
        client = FakeClient()
        worker = self.make(client, lambda *_args: "done", lambda *_args: None,
                           executor=ManualExecutor())
        worker._draining = True
        self.assertFalse(worker._accept_lease(lease_response()))
        self.assertEqual(len(client.named("start_session")), 0)
        self.assertEqual(len(client.named("suspend_session")), 0)
        self.assertEqual(client.named("fail_session")[0]["args"][3], {
            "error": {
                "errorType": "WorkerStopping",
                "message": "worker_not_accepting",
            },
            "retryAfterSeconds": 0,
        })

    def test_stop_preserves_pending_terminal_intent_instead_of_replacing_it(self):
        client = FakeClient(
            lease_task=[lease_response()],
            complete_session=[
                ControlPlaneUnavailable("offline"),
                ControlPlaneUnavailable("still offline"),
            ],
        )
        stopped = []
        worker = self.make(
            client, lambda *_args: "done",
            lambda lifecycle, reason: stopped.append((lifecycle.session_id, reason)),
            executor=ManualExecutor(inline=True))

        self.assertTrue(worker.run_once())
        self.assertEqual(worker.active_count(), 1)
        self.assertFalse(worker.stop())

        self.assertEqual(stopped, [("s1", "worker_stopping")])
        self.assertEqual(len(client.named("complete_session")), 2)
        self.assertEqual(client.named("fail_session"), [])
        controller = worker._sessions["s1"].controller
        self.assertEqual(controller.pending_terminal, "complete")


if __name__ == "__main__":
    unittest.main(verbosity=2)
