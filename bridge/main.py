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
        self.process: Optional[subprocess.Popen[str]] = None
        self.external_pid: Optional[int] = None
        self.ready = threading.Event()
        self._pump: Optional[threading.Thread] = None

    @property
    def pid(self) -> Optional[int]:
        if self.process is not None:
            return self.process.pid
        return self.external_pid

    def start(self, *, adopt: bool = False) -> None:
        if adopt:
            pid = self._read_live_pid()
            if pid is not None:
                self.external_pid = pid
                self.ready.set()
                LOG.info("adopted %s pid=%s", self.spec.name, pid)
                return
        self.pidfile.unlink(missing_ok=True)
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
        )
        self.pidfile.parent.mkdir(parents=True, exist_ok=True)
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
            return pid
        except (OSError, ValueError):
            return None

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

    def stop(self) -> None:
        pid = self.pid
        if pid is None:
            self.pidfile.unlink(missing_ok=True)
            return
        try:
            os.kill(pid, signal.SIGTERM)
        except ProcessLookupError:
            pass
        # Scheduler owns its drain retry loop.  Never replace a still-draining
        # process or escalate to SIGKILL.
        while self.is_alive():
            time.sleep(0.1)
        if self.process is not None:
            self.process.wait()
        self.pidfile.unlink(missing_ok=True)


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
            if component.is_alive():
                component.stop()
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
                self.components.pop(spec.name, None)
                self._start_until_ready(spec)
                if self.stop_event.is_set():
                    break
        self._shutdown()
        return 0

    def _preserve_worker_requested(self) -> bool:
        marker = self.state_dir / "preserve-persistent-worker-once"
        try:
            owner = marker.read_text(encoding="utf-8").strip()
        except OSError:
            return False
        # Consume stale/mismatched markers too. A later PID reuse must never
        # turn an abandoned one-shot handshake into permission to preserve a
        # Worker during an explicit full stop.
        marker.unlink(missing_ok=True)
        if owner != str(os.getpid()):
            return False
        return True

    def _shutdown(self) -> None:
        preserve_worker = self._preserve_worker_requested()
        # Quiesce scheduled admissions before stopping either inbound work or
        # the lease executor.  The Worker is always last so already-admitted
        # Scheduler work and inbound persistence can settle first.
        shutdown_order = (
            (SCHEDULER, DINGTALK_BOT, PERSISTENT_WORKER)
            if self.role == "scheduler" else (PERSISTENT_WORKER,)
        )
        for spec in shutdown_order:
            if preserve_worker and spec is PERSISTENT_WORKER:
                LOG.info(
                    "controlled restart: preserving Persistent Worker pid=%s",
                    self.components.get(spec.name).pid
                    if self.components.get(spec.name) is not None else "?",
                )
                continue
            component = self.components.pop(spec.name, None)
            if component is not None:
                component.stop()


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
