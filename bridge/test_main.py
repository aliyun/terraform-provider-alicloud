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
    def __init__(self, spec, *, events, generation, stop_event=None):
        self.spec = spec
        self.events = events
        self.generation = generation
        self.stop_event = stop_event
        self.pid = generation
        self.alive = True

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

    def stop(self):
        self.events.append(("stop", self.spec.name, self.generation))
        self.alive = False


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
        worker_stop = events.index(("stop", "persistent-worker", 1))
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
            events.index(("stop", "persistent-worker", 1)),
        )
        self.assertFalse(any(
            event[0] == "start" and event[1] == "dingtalk-bot"
            for event in events
        ))

    def test_full_stop_quiesces_scheduler_first_and_worker_last(self):
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

        stops = [event[1] for event in events if event[0] == "stop"]
        self.assertEqual(
            stops, ["scheduler", "dingtalk-bot", "persistent-worker"])

    def test_controlled_restart_preserves_worker_and_consumes_marker(self):
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

        stops = [event[1] for event in events if event[0] == "stop"]
        self.assertEqual(stops, ["scheduler", "dingtalk-bot"])

    def test_stale_restart_marker_is_consumed_without_preserving_worker(self):
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
            marker.write_text("not-this-supervisor\n", encoding="utf-8")
            supervisor._shutdown()
            self.assertFalse(marker.exists())

        stops = [event[1] for event in events if event[0] == "stop"]
        self.assertEqual(
            stops, ["scheduler", "dingtalk-bot", "persistent-worker"])


if __name__ == "__main__":
    unittest.main()
