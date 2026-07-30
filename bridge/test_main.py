#!/usr/bin/env python3
"""Tests for the Bridge process supervisor lifecycle contract."""

import os
from pathlib import Path
import subprocess
import sys
import tempfile
import threading
import time
import unittest
from unittest import mock

from bridge import main


def _pid_is_running(pid):
    try:
        state = subprocess.check_output(
            ["ps", "-o", "stat=", "-p", str(pid)],
            text=True,
            stderr=subprocess.DEVNULL,
        ).strip()
    except subprocess.CalledProcessError:
        return False
    return bool(state) and not state.startswith("Z")


class _FakeComponent:
    def __init__(
            self, spec, *, events, generation, stop_event=None,
            graceful=True, wait_delay=0.0, exit_status=0):
        self.spec = spec
        self.events = events
        self.generation = generation
        self.stop_event = stop_event
        self.pid = generation
        self.alive = True
        self.graceful = graceful
        self.wait_delay = wait_delay
        self.exit_status = exit_status
        self.returncode = None

    def start(self, *, adopt=False):
        self.events.append(("start", self.spec.name, self.generation, adopt))

    def wait_ready(self, _timeout):
        self.events.append(("ready", self.spec.name, self.generation))
        if (self.spec is main.SCHEDULER and self.generation == 2
                and self.stop_event is not None):
            self.stop_event.set()
        return True

    def is_alive(self):
        # The first Scheduler dies only after initial READY. Its replacement
        # stays alive long enough for the test to request supervisor shutdown.
        if self.spec is main.SCHEDULER and self.generation == 1:
            self.alive = False
        return self.alive

    def request_stop(self):
        self.events.append(("request-stop", self.spec.name, self.generation))
        if self.graceful:
            self.alive = False
            self.returncode = self.exit_status

    def wait(self, timeout):
        self.events.append(("wait", self.spec.name, timeout))
        if self.alive and self.wait_delay:
            time.sleep(min(self.wait_delay, timeout))
        return not self.alive

    def force_stop(self):
        self.events.append(("force-stop", self.spec.name, self.generation))
        self.alive = False
        self.returncode = -9

    def stop(self, timeout=5.0):
        self.request_stop()
        if self.wait(timeout):
            return True
        self.force_stop()
        return self.wait(2.0)


class BridgeSupervisorTest(unittest.TestCase):
    def _env(self, state_dir, role="scheduler"):
        return {
            "JARVIS_BRIDGE_ROLE": role,
            "JARVIS_CONTROL_PLANE_TOKEN": "test-token",
            "JARVIS_BRIDGE_STATE_DIR": str(state_dir),
            "JARVIS_BRIDGE_COMPONENT_RESTART_SEC": "0",
            "JARVIS_BRIDGE_COMPONENT_READY_WAIT": "1",
        }

    def test_auto_dispatch_zero_is_supported_for_both_roles(self):
        scheduler = main.validate_runtime({
            "JARVIS_BRIDGE_ROLE": "scheduler",
            "JARVIS_CONTROL_PLANE_TOKEN": "test-token",
            "JARVIS_AUTO_DISPATCH": "0",
            "JARVIS_BOOT_ID": "test-boot",
        })
        role = main.validate_runtime({
            "JARVIS_BRIDGE_ROLE": "worker",
            "JARVIS_CONTROL_PLANE_TOKEN": "test-token",
            "JARVIS_AUTO_DISPATCH": "0",
            "JARVIS_BOOT_ID": "test-boot",
        })
        self.assertEqual(scheduler, "scheduler")
        self.assertEqual(role, "worker")

    def test_worker_role_contains_only_persistent_worker(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir), role="worker"))
        self.assertEqual(supervisor.specs, (main.PERSISTENT_WORKER,))

    def test_component_starts_in_private_session_and_keeps_spawn_pgid(self):
        process = mock.Mock()
        process.pid = 4321
        process.stdout = []
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch.object(
                main.subprocess, "Popen", return_value=process) as popen, \
             mock.patch.object(
                 main, "_process_start_identity", return_value="birth-4321"), \
             mock.patch.object(main.os, "getpgid", return_value=4321):
            component = main.SubprocessComponent(
                main.PERSISTENT_WORKER,
                environ=self._env(Path(temp_dir), role="worker"),
                state_dir=Path(temp_dir),
            )
            component.start()
            component._pump.join(timeout=1)

        self.assertTrue(popen.call_args.kwargs["start_new_session"])
        self.assertEqual(component.process_group_id, 4321)

    def test_component_stop_kills_grandchild_after_leader_exits(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            state_dir = Path(temp_dir)
            grandchild_pidfile = state_dir / "grandchild.pid"
            entrypoint = state_dir / "leader.py"
            child = (
                "import os,pathlib,signal,time\n"
                "signal.signal(signal.SIGTERM, signal.SIG_IGN)\n"
                "pathlib.Path(os.environ['FAULT_GRANDCHILD_PIDFILE'])"
                ".write_text(str(os.getpid()))\n"
                "time.sleep(60)\n"
            )
            entrypoint.write_text(
                "import os,pathlib,subprocess,sys,time\n"
                "subprocess.Popen([sys.executable, '-c', %r])\n"
                "target=pathlib.Path(os.environ['FAULT_GRANDCHILD_PIDFILE'])\n"
                "deadline=time.monotonic()+2\n"
                "while not target.exists():\n"
                "    if time.monotonic() >= deadline: raise SystemExit(2)\n"
                "    time.sleep(0.01)\n" % child,
                encoding="utf-8",
            )
            spec = main.ComponentSpec(
                "fault-component",
                "unused",
                "unused",
                "fault-component.pid",
                "FAULT_COMPONENT_ENTRY",
            )
            environ = self._env(state_dir, role="worker")
            environ.update({
                "FAULT_COMPONENT_ENTRY": str(entrypoint),
                "FAULT_GRANDCHILD_PIDFILE": str(grandchild_pidfile),
                "JARVIS_BRIDGE_COMPONENT_TERM_GRACE": "0.1",
                "JARVIS_BRIDGE_COMPONENT_KILL_GRACE": "1",
            })
            component = main.SubprocessComponent(
                spec,
                environ=environ,
                state_dir=state_dir,
            )
            component.start()
            assert component.process is not None
            component.process.wait(timeout=3)
            grandchild = int(grandchild_pidfile.read_text(encoding="utf-8"))
            self.assertTrue(_pid_is_running(grandchild))

            component.stop()
            deadline = time.monotonic() + 2
            while time.monotonic() < deadline and _pid_is_running(grandchild):
                time.sleep(0.05)
            self.assertFalse(
                _pid_is_running(grandchild),
                "component stop left a TERM-ignoring grandchild alive",
            )
            self.assertFalse(component.pidfile.exists())
            self.assertIsNone(component.process_group_id)

    def test_adopted_component_stops_by_persisted_group_without_popen(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            state_dir = Path(temp_dir)
            component = main.SubprocessComponent(
                main.PERSISTENT_WORKER,
                environ=self._env(state_dir, role="worker"),
                state_dir=state_dir,
            )
            component.pidfile.write_text(
                "%s\n" % os.getpid(),
                encoding="utf-8",
            )
            component.identity_file.write_text(
                "%s|%s|test-birth\n" % (os.getpid(), os.getpid()),
                encoding="utf-8",
            )
            with mock.patch.object(
                    main, "_process_start_identity", return_value="test-birth"), \
                 mock.patch.object(
                     main.os, "getpgid", return_value=os.getpid()), \
                 mock.patch.object(
                    main,
                    "terminate_process_group",
                    return_value=True) as terminate:
                component.start(adopt=True)
                self.assertIsNone(component.process)
                self.assertEqual(component.external_pid, os.getpid())
                self.assertEqual(component.process_group_id, os.getpid())
                component.stop()

            terminate.assert_called_once()
            self.assertIsNone(terminate.call_args.args[0])
            self.assertEqual(
                terminate.call_args.kwargs["pgid"],
                os.getpid(),
            )
            self.assertFalse(component.pidfile.exists())
            self.assertFalse(component.identity_file.exists())
            self.assertIsNone(component.external_pid)
            self.assertIsNone(component.process_group_id)

    def test_component_stop_fails_closed_on_reused_group_leader(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            state_dir = Path(temp_dir)
            component = main.SubprocessComponent(
                main.PERSISTENT_WORKER,
                environ=self._env(state_dir, role="worker"),
                state_dir=state_dir,
            )
            component.external_pid = 777
            component.process_group_id = 777
            component.process_start_identity = "original-birth"
            component.pidfile.write_text("777\n", encoding="utf-8")
            component.identity_file.write_text(
                "777|777|original-birth\n",
                encoding="utf-8",
            )
            with mock.patch.object(main, "_pid_exists", return_value=True), \
                 mock.patch.object(
                     main, "_process_start_identity", return_value="reused-birth"), \
                 mock.patch.object(
                     main, "terminate_process_group") as terminate:
                with self.assertRaisesRegex(RuntimeError, "reused/uninspectable"):
                    component.stop()
            terminate.assert_not_called()
            self.assertTrue(component.pidfile.exists())
            self.assertTrue(component.identity_file.exists())
            self.assertEqual(component.process_group_id, 777)

    def test_scheduler_failure_restarts_without_stopping_ready_worker(self):
        events = []
        generations = {}
        stop_event = threading.Event()

        def factory(spec, **_kwargs):
            generations[spec.name] = generations.get(spec.name, 0) + 1
            return _FakeComponent(
                spec,
                events=events,
                generation=generations[spec.name],
                stop_event=stop_event,
            )

        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"):
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir)),
                component_factory=factory,
                stop_event=stop_event,
            )
            self.assertEqual(supervisor.run(), 0)

        scheduler_restart = events.index(("start", "scheduler", 2, False))
        worker_stop = events.index(("request-stop", "persistent-worker", 1))
        self.assertLess(scheduler_restart, worker_stop)
        self.assertEqual(generations["persistent-worker"], 1)

    def test_outer_ready_is_reported_before_bot_startup(self):
        events = []
        generations = {}
        stop_event = threading.Event()

        def factory(spec, **_kwargs):
            generations[spec.name] = generations.get(spec.name, 0) + 1
            return _FakeComponent(
                spec,
                events=events,
                generation=generations[spec.name],
                stop_event=stop_event,
            )

        def observe_ready(message, *args):
            if message.startswith("Bridge READY"):
                events.append(("supervisor-ready",))
                stop_event.set()

        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"), mock.patch.object(
                main.LOG, "info", side_effect=observe_ready):
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir)),
                component_factory=factory,
                stop_event=stop_event,
            )
            self.assertEqual(supervisor.run(), 0)

        self.assertLess(
            events.index(("supervisor-ready",)),
            events.index(("request-stop", "persistent-worker", 1)),
        )
        self.assertFalse(any(
            event[0] == "start" and event[1] == "dingtalk-bot"
            for event in events
        ))

    def test_full_stop_quiesces_every_component_before_waiting(self):
        events = []
        generations = {}

        def factory(spec, **_kwargs):
            generations[spec.name] = 1
            return _FakeComponent(spec, events=events, generation=1)

        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"):
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir)),
                component_factory=factory,
            )
            for spec in supervisor.specs:
                component = factory(spec)
                supervisor.components[spec.name] = component
            supervisor._shutdown()

        stops = [event[1] for event in events if event[0] == "request-stop"]
        self.assertEqual(
            stops, ["scheduler", "dingtalk-bot", "persistent-worker"])
        first_wait = next(
            index for index, event in enumerate(events) if event[0] == "wait")
        self.assertTrue(all(
            index < first_wait
            for index, event in enumerate(events)
            if event[0] == "request-stop"
        ))

    def test_shutdown_uses_one_shared_deadline_and_force_stops_survivors(self):
        events = []

        def factory(spec, **_kwargs):
            return _FakeComponent(
                spec,
                events=events,
                generation=1,
                graceful=False,
                wait_delay=0.02,
            )

        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"):
            env = self._env(Path(temp_dir))
            env["JARVIS_BRIDGE_STOP_WAIT"] = "0.03"
            env["JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS"] = "0.03"
            supervisor = main.BridgeSupervisor(
                environ=env,
                component_factory=factory,
            )
            self.assertEqual(supervisor.shutdown_timeout, 0.03)
            self.assertEqual(supervisor.shutdown_grace_timeout, 0.0)
            self.assertEqual(supervisor.shutdown_force_timeout, 0.03)
            for spec in supervisor.specs:
                supervisor.components[spec.name] = factory(spec)
            started = time.monotonic()
            self.assertFalse(supervisor._shutdown())
            elapsed = time.monotonic() - started

        self.assertLess(elapsed, 0.09, "deadline must not multiply per component")
        requests = [
            index for index, event in enumerate(events)
            if event[0] == "request-stop"
        ]
        waits = [
            index for index, event in enumerate(events)
            if event[0] == "wait"
        ]
        self.assertLess(max(requests), min(waits))
        self.assertEqual(
            [event[1] for event in events if event[0] == "force-stop"],
            ["scheduler", "dingtalk-bot", "persistent-worker"],
        )

    def test_forced_shutdown_flags_survivors_for_offline(self):
        def factory(spec, **_kwargs):
            return _FakeComponent(
                spec, events=[], generation=1, graceful=False, wait_delay=0.0)

        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"):
            env = self._env(Path(temp_dir))
            env["JARVIS_BRIDGE_STOP_WAIT"] = "0.03"
            env["JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS"] = "0.03"
            supervisor = main.BridgeSupervisor(
                environ=env, component_factory=factory)
            for spec in supervisor.specs:
                supervisor.components[spec.name] = factory(spec)
            supervisor._shutdown()

        # Force-kill happened → run() must release the fences afterwards.
        self.assertTrue(supervisor._shutdown_forced)

    def test_clean_shutdown_does_not_flag_offline(self):
        def factory(spec, **_kwargs):
            return _FakeComponent(spec, events=[], generation=1)

        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"):
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir)), component_factory=factory)
            for spec in supervisor.specs:
                supervisor.components[spec.name] = factory(spec)
            supervisor._shutdown()

        self.assertFalse(supervisor._shutdown_forced)

    def test_offline_helper_invokes_worker_offline_module(self):
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"), mock.patch(
                "bridge.main.subprocess.run") as run:
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir)))
            supervisor._offline_forced_workers()

        run.assert_called_once()
        argv = run.call_args.args[0]
        self.assertIn("-m", argv)
        self.assertIn("bridge.worker_offline", argv)
        self.assertIn("--all", argv)

    def test_offline_helper_never_raises_when_subprocess_fails(self):
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"), mock.patch(
                "bridge.main.subprocess.run",
                side_effect=OSError("boom")):
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir)))
            # Must swallow the error — a failed fence release cannot break stop.
            supervisor._offline_forced_workers()

    def test_default_shutdown_budget_reserves_force_reap_inside_30_seconds(self):
        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"):
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir)),
            )

        self.assertEqual(supervisor.shutdown_timeout, 30.0)
        self.assertEqual(supervisor.shutdown_grace_timeout, 28.0)
        self.assertEqual(supervisor.shutdown_force_timeout, 2.0)

    def test_worker_handoff_failure_is_not_silently_reported_as_clean_shutdown(self):
        events = []

        def factory(spec, **_kwargs):
            return _FakeComponent(
                spec,
                events=events,
                generation=1,
                exit_status=1 if spec is main.PERSISTENT_WORKER else 0,
            )

        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"), self.assertLogs(
                    "jarvis-bridge-supervisor", level="WARNING") as logs:
            supervisor = main.BridgeSupervisor(
                environ=self._env(Path(temp_dir)),
                component_factory=factory,
            )
            for spec in supervisor.specs:
                supervisor.components[spec.name] = factory(spec)
            self.assertFalse(supervisor._shutdown())

        self.assertIn(
            "persistent-worker shutdown reported exit status 1",
            "\n".join(logs.output),
        )

    def test_legacy_restart_marker_does_not_preserve_worker(self):
        events = []

        def factory(spec, **_kwargs):
            return _FakeComponent(spec, events=events, generation=1)

        with tempfile.TemporaryDirectory() as temp_dir, mock.patch(
                "bridge.scheduler.scheduler.validate"):
            state_dir = Path(temp_dir)
            supervisor = main.BridgeSupervisor(
                environ=self._env(state_dir),
                component_factory=factory,
            )
            for spec in supervisor.specs:
                supervisor.components[spec.name] = factory(spec)
            marker = state_dir / "preserve-persistent-worker-once"
            marker.write_text("%s\n" % os.getpid(), encoding="utf-8")
            supervisor._shutdown()
            self.assertFalse(marker.exists())

        stops = [event[1] for event in events if event[0] == "request-stop"]
        self.assertEqual(
            stops, ["scheduler", "dingtalk-bot", "persistent-worker"])


if __name__ == "__main__":
    unittest.main()
