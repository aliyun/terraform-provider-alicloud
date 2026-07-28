#!/usr/bin/env python3
"""Bridge process supervisor.

The supervisor owns lifecycle only.  Scheduler, durable Task execution, and
DingTalk inbound handling remain independent process entrypoints so one
component can be restarted without propagating failure to healthy components.
"""

from __future__ import annotations

from dataclasses import dataclass
import logging
import os
from pathlib import Path
import signal
import subprocess
import sys
import threading
import time
from typing import Any, Callable, Mapping, Optional, Sequence

from bridge.process_group_runner import (
    _signal_process_group,
    _wait_group_empty,
    terminate_process_group,
)
from bridge.process_identity import (
    pid_exists as _pid_exists,
    process_command as _process_command,
    process_start_identity as _process_start_identity,
)


LOG = logging.getLogger("jarvis-bridge-supervisor")
REPO_ROOT = Path(__file__).resolve().parents[1]


@dataclass(frozen=True)
class ComponentSpec:
    name: str
    module: str
    ready_marker: str
    pidfile: str
    override: str

    def command(self, environ: Mapping[str, str]) -> list[str]:
        overridden = environ.get(self.override, "").strip()
        if overridden:
            return [sys.executable, overridden]
        return [sys.executable, "-m", self.module]


SCHEDULER = ComponentSpec(
    "scheduler",
    "bridge.scheduler.scheduler",
    "Scheduler READY",
    "scheduler.pid",
    "JARVIS_SCHEDULER_ENTRY",
)
PERSISTENT_WORKER = ComponentSpec(
    "persistent-worker",
    "bridge.persistent_worker",
    "Persistent worker READY",
    "persistent-worker.pid",
    "JARVIS_PERSISTENT_WORKER",
)
DINGTALK_BOT = ComponentSpec(
    "dingtalk-bot",
    "bridge.jarvis_dingtalk_bot",
    "Bridge READY",
    "dingtalk-bot.pid",
    "JARVIS_BRIDGE_BOT",
)


def validate_runtime(environ: Mapping[str, str] | None = None) -> str:
    """Validate the selected role before any component can lease work."""

    env = os.environ if environ is None else environ
    role = env.get("JARVIS_BRIDGE_ROLE", "scheduler").strip()
    if role not in ("scheduler", "worker"):
        raise RuntimeError(
            "unsupported JARVIS_BRIDGE_ROLE=%s (accept: scheduler|worker)"
            % role)
    token = (
        env.get("JARVIS_CONTROL_PLANE_TOKEN", "").strip()
        or env.get("JARVIS_HTML_REPORT_TOKEN", "").strip()
    )
    if not token:
        raise RuntimeError(
            "bridge requires JARVIS_CONTROL_PLANE_TOKEN or "
            "JARVIS_HTML_REPORT_TOKEN")
    if role == "scheduler":
        # Lazy import keeps worker-only hosts independent of Scheduler/PyYAML.
        from bridge.scheduler.scheduler import validate as validate_scheduler
        validate_scheduler(env)
    return role


class SubprocessComponent:
    """One independently restartable bridge entrypoint."""

    def __init__(
        self,
        spec: ComponentSpec,
        *,
        environ: Mapping[str, str],
        state_dir: Path,
    ) -> None:
        self.spec = spec
        self.environ = dict(environ)
        self.pidfile = state_dir / spec.pidfile
        self.identity_file = state_dir / ("%s.identity" % spec.pidfile)
        self.process: Optional[subprocess.Popen[str]] = None
        self.process_group_id: Optional[int] = None
        self.process_start_identity: Optional[str] = None
        self.external_pid: Optional[int] = None
        self._returncode: Optional[int] = None
        self.ready = threading.Event()
        self._pump: Optional[threading.Thread] = None

    @property
    def pid(self) -> Optional[int]:
        if self.process is not None:
            return self.process.pid
        return self.external_pid

    def _abort_unverified_spawn(self, message: str) -> None:
        process = self.process
        pgid = self.process_group_id
        if process is not None and pgid is not None:
            terminate_process_group(
                process,
                pgid=pgid,
                term_grace=0.0,
                kill_grace=1.0,
            )
            if process.stdout is not None:
                process.stdout.close()
        self.process = None
        self.process_group_id = None
        self.process_start_identity = None
        raise RuntimeError(message)

    def start(self, *, adopt: bool = False) -> None:
        self._returncode = None
        if adopt:
            pid = self._read_live_pid()
            if pid is not None:
                self.external_pid = pid
                # Components created by this supervisor always start a private
                # session, so their persisted spawn PID is also the stable
                # PGID.  Preserve that identity across controlled adoption.
                self.process_group_id = pid
                self.process_start_identity = self._read_identity(pid)
                self.ready.set()
                LOG.info("adopted %s pid=%s", self.spec.name, pid)
                return
        self.pidfile.unlink(missing_ok=True)
        self.identity_file.unlink(missing_ok=True)
        env = dict(self.environ)
        env["PYTHONUNBUFFERED"] = "1"
        pythonpath = env.get("PYTHONPATH", "")
        env["PYTHONPATH"] = (
            str(REPO_ROOT) + (os.pathsep + pythonpath if pythonpath else ""))
        self.process = subprocess.Popen(
            self.spec.command(env),
            cwd=str(REPO_ROOT),
            env=env,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1,
            errors="replace",
            start_new_session=True,
        )
        # start_new_session performs setsid before exec, so the child PID is the
        # stable PGID even if the component leader exits before a watchdog sweep.
        self.process_group_id = self.process.pid
        self.process_start_identity = _process_start_identity(self.process.pid)
        if self.process_start_identity is None:
            self._abort_unverified_spawn(
                "cannot capture %s process start identity"
                % self.spec.name
            )
        try:
            current_pgid = os.getpgid(self.process.pid)
        except OSError:
            self._abort_unverified_spawn(
                "cannot verify %s private process group"
                % self.spec.name
            )
        if current_pgid != self.process.pid:
            self._abort_unverified_spawn(
                "%s did not start as a private process group"
                % self.spec.name
            )
        self.pidfile.parent.mkdir(parents=True, exist_ok=True)
        self._write_identity(
            self.process.pid,
            self.process.pid,
            self.process_start_identity,
        )
        self.pidfile.write_text("%s\n" % self.process.pid, encoding="utf-8")
        self._pump = threading.Thread(
            target=self._pump_output,
            name="bridge-%s-log" % self.spec.name,
            daemon=True,
        )
        self._pump.start()

    def _read_live_pid(self) -> Optional[int]:
        try:
            pid = int(self.pidfile.read_text(encoding="utf-8").strip())
            os.kill(pid, 0)
        except (OSError, ValueError):
            return None
        start_identity = _process_start_identity(pid)
        try:
            process_group_id = os.getpgid(pid)
        except OSError as exc:
            raise RuntimeError(
                "cannot verify adopted %s pid=%s" % (self.spec.name, pid)
            ) from exc
        if start_identity is None or process_group_id != pid:
            raise RuntimeError(
                "refusing unverified adopted %s pid=%s pgid=%s"
                % (self.spec.name, pid, process_group_id)
            )

        persisted = self._read_identity(pid)
        if persisted is not None:
            if persisted != start_identity:
                raise RuntimeError(
                    "refusing reused adopted %s pid=%s"
                    % (self.spec.name, pid)
                )
            return pid

        # One-release compatibility for components started before identity
        # companions existed.  A private group plus the expected command proves
        # enough ownership to adopt it and persist its observed birth identity.
        command = _process_command(pid)
        expected = self.environ.get(self.spec.override, "").strip()
        owner_token = expected or self.spec.module
        if not command or owner_token not in command:
            raise RuntimeError(
                "refusing legacy adopted %s pid=%s without ownership proof"
                % (self.spec.name, pid)
            )
        self._write_identity(pid, pid, start_identity)
        return pid

    def _write_identity(self, pid: int, pgid: int, start_identity: str) -> None:
        value = "%s|%s|%s\n" % (int(pid), int(pgid), start_identity)
        temporary = self.identity_file.with_name(
            "%s.tmp.%s" % (self.identity_file.name, os.getpid())
        )
        temporary.write_text(value, encoding="utf-8")
        os.replace(temporary, self.identity_file)

    def _read_identity(self, expected_pid: int) -> Optional[str]:
        try:
            fields = self.identity_file.read_text(
                encoding="utf-8").strip().split("|", 2)
            pid, pgid = int(fields[0]), int(fields[1])
            start_identity = fields[2]
        except FileNotFoundError:
            return None
        except (OSError, ValueError, IndexError) as exc:
            raise RuntimeError(
                "invalid %s process identity companion"
                % self.spec.name
            ) from exc
        if (
            pid != int(expected_pid)
            or pgid != int(expected_pid)
            or not start_identity
        ):
            raise RuntimeError(
                "invalid %s process identity companion"
                % self.spec.name
            )
        return start_identity

    def _verify_group_owner_before_signal(self, pgid: int) -> None:
        # If the original group leader remains visible, both its birth identity
        # and its private PGID must still match.  A different birth token proves
        # numeric PID/PGID reuse.  When the leader is gone, the original group
        # may legitimately persist through descendants.
        if not _pid_exists(pgid):
            return
        current = _process_start_identity(pgid)
        if (
            current is None
            or self.process_start_identity is None
            or current != self.process_start_identity
        ):
            raise RuntimeError(
                "refusing reused/uninspectable %s pgid=%s"
                % (self.spec.name, pgid)
            )
        try:
            current_pgid = os.getpgid(pgid)
        except OSError as exc:
            raise RuntimeError(
                "cannot verify %s pgid=%s" % (self.spec.name, pgid)
            ) from exc
        if current_pgid != pgid:
            raise RuntimeError(
                "refusing non-private %s pid=%s pgid=%s"
                % (self.spec.name, pgid, current_pgid)
            )

    def _pump_output(self) -> None:
        assert self.process is not None
        assert self.process.stdout is not None
        for line in self.process.stdout:
            sys.stderr.write(line)
            sys.stderr.flush()
            if self.spec.ready_marker in line:
                self.ready.set()

    def wait_ready(self, timeout: float) -> bool:
        deadline = time.monotonic() + timeout
        while time.monotonic() < deadline:
            if self.ready.wait(min(0.1, max(0.0, deadline - time.monotonic()))):
                return True
            if not self.is_alive():
                return False
        return self.ready.is_set()

    def is_alive(self) -> bool:
        if self.process is not None:
            return self.process.poll() is None
        if self.external_pid is None:
            return False
        try:
            os.kill(self.external_pid, 0)
            return True
        except OSError:
            return False

    def _term_grace_seconds(self) -> float:
        common = self.environ.get("JARVIS_BRIDGE_COMPONENT_TERM_GRACE")
        if common is not None:
            return max(0.0, float(common))
        if self.spec is SCHEDULER:
            value = self.environ.get(
                "JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS",
                self.environ.get(
                    "JARVIS_BRIDGE_STOP_WAIT",
                    self.environ.get(
                        "JARVIS_STOP_GRACE",
                        self.environ.get("JARVIS_SCHEDULER_DRAIN_TIMEOUT", "30"),
                    ),
                ),
            )
        elif self.spec is PERSISTENT_WORKER:
            value = self.environ.get("JARVIS_WORKER_DRAIN_TIMEOUT", "30")
        else:
            value = self.environ.get(
                "JARVIS_BRIDGE_STOP_WAIT",
                self.environ.get("JARVIS_STOP_GRACE", "30"),
            )
        return max(0.0, float(value))

    def _remove_pidfile_if_owned(self, pid: int) -> None:
        try:
            current = int(self.pidfile.read_text(encoding="utf-8").strip())
        except (OSError, ValueError):
            return
        if current == pid:
            self.pidfile.unlink(missing_ok=True)
            self.identity_file.unlink(missing_ok=True)

    def _finalize_stopped(self, pid: Optional[int]) -> None:
        if self.process is not None:
            self.process.wait()
            self._returncode = self.process.returncode
        if pid is not None:
            self._remove_pidfile_if_owned(pid)
        if self._pump is not None:
            self._pump.join(timeout=1.0)
        if self.process is not None and self.process.stdout is not None:
            self.process.stdout.close()
        self.process = None
        self.process_group_id = None
        self.process_start_identity = None
        self.external_pid = None
        self._pump = None
        self.ready.clear()

    def request_stop(self) -> None:
        pid = self.pid
        pgid = self.process_group_id
        if pid is None and pgid is None:
            self.pidfile.unlink(missing_ok=True)
            self.identity_file.unlink(missing_ok=True)
            return

        if pgid is not None:
            self._verify_group_owner_before_signal(pgid)
            _signal_process_group(pgid, signal.SIGTERM)
        elif pid is not None:
            # Compatibility for a component created before private sessions
            # were mandatory.  New and adopted components always take the
            # process-group branch above.
            try:
                os.kill(pid, signal.SIGTERM)
            except ProcessLookupError:
                pass

    def wait(self, timeout: float) -> bool:
        pid = self.pid
        pgid = self.process_group_id
        if pid is None and pgid is None:
            self.pidfile.unlink(missing_ok=True)
            self.identity_file.unlink(missing_ok=True)
            return True
        if pgid is not None:
            if not _wait_group_empty(pgid, timeout):
                return False
        else:
            deadline = time.monotonic() + max(0.0, timeout)
            while self.is_alive() and time.monotonic() < deadline:
                time.sleep(min(0.05, max(0.0, deadline - time.monotonic())))
            if self.is_alive():
                return False
        self._finalize_stopped(pid)
        return True

    def force_stop(self) -> None:
        pid = self.pid
        pgid = self.process_group_id
        if pid is None and pgid is None:
            self.pidfile.unlink(missing_ok=True)
            self.identity_file.unlink(missing_ok=True)
            return
        if pgid is not None:
            self._verify_group_owner_before_signal(pgid)
            _signal_process_group(pgid, signal.SIGKILL)
        elif pid is not None:
            try:
                os.kill(pid, signal.SIGKILL)
            except ProcessLookupError:
                pass

    @property
    def returncode(self) -> Optional[int]:
        if self.process is not None:
            return self.process.poll()
        return self._returncode

    def stop(self, timeout: Optional[float] = None) -> bool:
        """Compatibility helper for a single failed-start component."""
        term_grace = (
            self._term_grace_seconds()
            if timeout is None else max(0.0, float(timeout))
        )
        pgid = self.process_group_id
        if pgid is not None:
            self._verify_group_owner_before_signal(pgid)
            kill_grace = max(
                0.0,
                float(self.environ.get(
                    "JARVIS_BRIDGE_COMPONENT_KILL_GRACE", "2")),
            )
            if not terminate_process_group(
                    self.process,
                    pgid=pgid,
                    term_grace=term_grace,
                    kill_grace=kill_grace):
                raise RuntimeError(
                    "%s process group %s did not drain"
                    % (self.spec.name, pgid)
                )
            self._finalize_stopped(self.pid)
            return True

        self.request_stop()
        if self.wait(term_grace):
            return True
        self.force_stop()
        kill_grace = max(
            0.0,
            float(self.environ.get(
                "JARVIS_BRIDGE_COMPONENT_KILL_GRACE", "2")),
        )
        if self.wait(kill_grace):
            return True
        raise RuntimeError(
            "%s process group %s did not drain"
            % (self.spec.name, self.process_group_id or self.pid)
        )


class BridgeSupervisor:
    """Start, monitor, and independently restart the selected entrypoints."""

    def __init__(
        self,
        *,
        environ: Mapping[str, str] | None = None,
        component_factory: Callable[..., Any] = SubprocessComponent,
        stop_event: Optional[threading.Event] = None,
    ) -> None:
        self.environ = dict(os.environ if environ is None else environ)
        self.role = validate_runtime(self.environ)
        self.state_dir = Path(self.environ.get(
            "JARVIS_BRIDGE_STATE_DIR",
            str(REPO_ROOT / ".my-day" / "bridge"),
        ))
        self.state_dir.mkdir(parents=True, exist_ok=True)
        self.component_factory = component_factory
        self.stop_event = stop_event or threading.Event()
        self.components: dict[str, Any] = {}
        self.ready_timeout = float(self.environ.get(
            "JARVIS_BRIDGE_COMPONENT_READY_WAIT",
            self.environ.get("JARVIS_SCHEDULER_READY_WAIT", "30"),
        ))
        self.restart_delay = float(self.environ.get(
            "JARVIS_BRIDGE_COMPONENT_RESTART_SEC", "5"))
        stop_wait = float(self.environ.get(
            "JARVIS_BRIDGE_STOP_WAIT",
            self.environ.get("JARVIS_STOP_GRACE", "30"),
        ))
        restart_wait = float(self.environ.get(
            "JARVIS_SCHEDULER_DRAIN_TIMEOUT_SECONDS", "30"))
        # run.sh owns the mechanical outer deadline. Use no more than either
        # compatibility timeout so stop and restart both leave enough time for
        # the supervisor to force-reap its exact children before the shell
        # escalates the supervisor itself.
        self.shutdown_timeout = max(0.0, min(stop_wait, restart_wait))
        # Reserve mechanical time inside the same total budget so the
        # supervisor can SIGKILL and reap exact children before run.sh/launchd
        # reaches its matching outer deadline.
        self.shutdown_force_timeout = min(2.0, self.shutdown_timeout)
        self.shutdown_grace_timeout = max(
            0.0, self.shutdown_timeout - self.shutdown_force_timeout)

    @property
    def specs(self) -> Sequence[ComponentSpec]:
        if self.role == "worker":
            return (PERSISTENT_WORKER,)
        # Scheduler must be validated and READY before the worker may lease.
        return (SCHEDULER, PERSISTENT_WORKER, DINGTALK_BOT)

    def request_stop(self) -> None:
        self.stop_event.set()

    def _new_component(self, spec: ComponentSpec) -> Any:
        return self.component_factory(
            spec, environ=self.environ, state_dir=self.state_dir)

    def _start_until_ready(
        self,
        spec: ComponentSpec,
        *,
        adopt: bool = False,
    ) -> bool:
        first = True
        while not self.stop_event.is_set():
            if spec is SCHEDULER:
                from bridge.scheduler.scheduler import validate
                validate(self.environ)
            component = self._new_component(spec)
            component.start(adopt=adopt and first)
            first = False
            if component.wait_ready(self.ready_timeout):
                self.components[spec.name] = component
                LOG.info("%s supervised READY pid=%s", spec.name, component.pid)
                return True
            LOG.error(
                "%s failed before READY; retrying independently in %.1fs",
                spec.name,
                self.restart_delay,
            )
            # A dead leader may still have live descendants in its private
            # group.  Always drain the captured PGID before retrying.
            component.stop(timeout=min(self.ready_timeout, 5.0))
            if self.stop_event.wait(self.restart_delay):
                break
        return False

    def run(self) -> int:
        # Scheduler validation/READY is the admission gate and the Worker must
        # be alive before the outer process manager considers this supervisor
        # healthy.  The Bot is deliberately not part of that readiness gate:
        # a credential/stream startup failure must never make run.sh/launchd
        # roll back an already-READY (or adopted) Persistent Worker.
        required_specs = (
            self.specs
            if self.role == "worker"
            else (SCHEDULER, PERSISTENT_WORKER)
        )
        for spec in required_specs:
            if not self._start_until_ready(
                    spec, adopt=spec is PERSISTENT_WORKER):
                self._shutdown()
                return 1
        LOG.info(
            "Bridge READY pid=%s role=supervisor required=%s",
            os.getpid(),
            ",".join(spec.name for spec in required_specs),
        )
        if self.role == "scheduler":
            # Retry forever while Scheduler and Worker stay alive.  This call
            # may block on repeated Bot startup failures, but the outer
            # readiness contract has already committed to preserving Worker.
            self._start_until_ready(DINGTALK_BOT)
        while not self.stop_event.wait(0.2):
            for spec in self.specs:
                component = self.components.get(spec.name)
                if component is not None and component.is_alive():
                    continue
                LOG.error("%s exited; healthy components stay running", spec.name)
                component = self.components.pop(spec.name, None)
                if component is not None:
                    component.stop()
                self._start_until_ready(spec)
                if self.stop_event.is_set():
                    break
        return 0 if self._shutdown() else 1

    def _remove_legacy_restart_marker(self) -> None:
        marker = self.state_dir / "preserve-persistent-worker-once"
        # Older bridge versions used this marker to keep the Persistent Worker
        # alive across an explicit restart. Restarts now stop every component;
        # remove any leftover marker so a downgrade cannot revive that
        # ambiguous behavior.
        marker.unlink(missing_ok=True)

    def _shutdown(self) -> bool:
        self._remove_legacy_restart_marker()
        # Request every component to quiesce before waiting for any one of them:
        # Scheduler stops admission, Bot stops inbound work, and Worker stops
        # leasing immediately. They then share one global grace deadline.
        shutdown_order = (
            (SCHEDULER, DINGTALK_BOT, PERSISTENT_WORKER)
            if self.role == "scheduler" else (PERSISTENT_WORKER,)
        )
        components = []
        for spec in shutdown_order:
            component = self.components.pop(spec.name, None)
            if component is not None:
                components.append(component)
                component.request_stop()

        shutdown_started = time.monotonic()
        deadline = shutdown_started + self.shutdown_grace_timeout
        survivors = []
        clean = True
        for component in components:
            remaining = max(0.0, deadline - time.monotonic())
            if component.wait(remaining):
                returncode = getattr(component, "returncode", None)
                if returncode not in (None, 0):
                    clean = False
                    LOG.warning(
                        "%s shutdown reported exit status %s",
                        component.spec.name,
                        returncode,
                    )
            else:
                survivors.append(component)

        if survivors:
            clean = False
            LOG.warning(
                "Bridge shutdown grace %.1fs expired; force-stopping: %s",
                self.shutdown_grace_timeout,
                ",".join(component.spec.name for component in survivors),
            )
            # Signal every survivor before waiting for any one of them, again
            # avoiding a per-component timeout multiplier.
            for component in survivors:
                component.force_stop()
            force_deadline = shutdown_started + self.shutdown_timeout
            for component in survivors:
                remaining = max(0.0, force_deadline - time.monotonic())
                if not component.wait(remaining):
                    LOG.error(
                        "managed component still alive after SIGKILL: %s pid=%s",
                        component.spec.name,
                        component.pid,
                    )
        return clean


def main() -> int:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s [%(threadName)s] %(message)s",
        stream=sys.stderr,
    )
    if sys.argv[1:] == ["--validate"]:
        validate_runtime()
        LOG.info("Bridge validation OK")
        return 0
    if sys.argv[1:]:
        raise SystemExit("usage: python -m bridge.main [--validate]")
    supervisor = BridgeSupervisor()

    def request_stop(signum: int, _frame: Any) -> None:
        LOG.info("Bridge supervisor signal %s received", signum)
        supervisor.request_stop()

    signal.signal(signal.SIGTERM, request_stop)
    signal.signal(signal.SIGINT, request_stop)
    return supervisor.run()


if __name__ == "__main__":
    raise SystemExit(main())
