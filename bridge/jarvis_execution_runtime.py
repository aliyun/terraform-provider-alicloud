#!/usr/bin/env python3
"""Shared subprocess runtime for Task and ephemeral execution.

The runtime deliberately has no Task/Session semantics.  A caller may request
the fenced ``ProcessGuardian`` launch path, but ownership transitions remain in
``PersistenceExecutor``/``SessionController``.  Ephemeral jobs use the same timeout,
process-group cleanup, stdout/stderr capture, and spawn observability without
creating control-plane state.
"""

from __future__ import annotations

import os
import hashlib
import json
import logging
import re
import signal
import socket
import subprocess
import sys
import tempfile
import threading
import time
from concurrent.futures import ThreadPoolExecutor
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable, Mapping, Optional, Sequence, Tuple
from urllib.request import Request, urlopen

from bridge.jarvis_capacity import CapacityManager


SpawnCallback = Callable[[Any], None]
GuardedSpawn = Callable[
    [Sequence[str], Path, SpawnCallback, Optional[Mapping[str, str]]],
    Tuple[Any, Optional[int]],
]

_HEADLESS_AMP_BROKER_ENV = "JARVIS_HEADLESS_AMP_BROKER"


def _process_start_identity(pid: int) -> str:
    try:
        completed = subprocess.run(
            ["/bin/ps", "-o", "lstart=", "-p", str(int(pid))],
            capture_output=True, text=True, timeout=2, check=False)
    except (OSError, subprocess.SubprocessError, TypeError, ValueError):
        return ""
    return completed.stdout.strip()


def _process_parent(pid: int) -> int:
    try:
        completed = subprocess.run(
            ["/bin/ps", "-o", "ppid=", "-p", str(int(pid))],
            capture_output=True, text=True, timeout=2, check=False)
        return int(completed.stdout.strip() or 0)
    except (OSError, subprocess.SubprocessError, TypeError, ValueError):
        return 0


def _process_command(pid: int) -> str:
    try:
        completed = subprocess.run(
            ["/bin/ps", "-o", "command=", "-p", str(int(pid))],
            capture_output=True, text=True, timeout=2, check=False)
    except (OSError, subprocess.SubprocessError, TypeError, ValueError):
        return ""
    return completed.stdout.strip()


def _descends_from(pid: int, ancestor_pid: int,
                   ancestor_started_at: str) -> bool:
    current = int(pid)
    for _depth in range(64):
        if current <= 1:
            return False
        if current == int(ancestor_pid):
            return bool(
                ancestor_started_at
                and _process_start_identity(current) == ancestor_started_at)
        parent = _process_parent(current)
        if parent <= 0 or parent == current:
            return False
        current = parent
    return False


@dataclass(frozen=True)
class AmpBrokerDecision:
    allowed: bool
    reason: str


class HeadlessAmpBroker:
    """Executor-owned, process-bound authorization broker for AMP mutations."""

    _ALLOWED_ACTIONS = frozenset((
        "init", "branch-create", "branch-switch",
        "context-set-branch", "publish",
    ))

    def __init__(self, controller: Any, runtime_session_id: str):
        self.controller = controller
        self.runtime_session_id = str(runtime_session_id)
        self.directory = tempfile.TemporaryDirectory(prefix="jarvis-amp-broker-")
        os.chmod(self.directory.name, 0o700)
        self.path = os.path.join(self.directory.name, "authorize.sock")
        self._socket = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self._socket.bind(self.path)
        os.chmod(self.path, 0o600)
        self._socket.listen(8)
        self._socket.settimeout(0.2)
        self._enabled = False
        self._bound_pid = 0
        self._bound_started_at = ""
        self._repo_identity: Optional[tuple[str, int, int]] = None
        self._project_identity: Optional[tuple[str, str, str]] = None
        self._branch_identity: Optional[str] = None
        self._pin_lock = threading.Lock()
        self._closed = threading.Event()
        self._thread = threading.Thread(
            target=self._serve, name="HeadlessAmpBroker", daemon=True)
        self._thread.start()

    def enable(self, process: Any) -> None:
        pid = int(getattr(process, "pid", 0))
        started_at = _process_start_identity(pid)
        if pid <= 0 or not started_at:
            raise RuntimeError("AMP broker cannot bind headless process")
        self._bound_pid = pid
        self._bound_started_at = started_at
        self._enabled = True

    def close(self) -> None:
        self._enabled = False
        self._closed.set()
        try:
            self._socket.close()
        except OSError:
            pass
        self._thread.join(timeout=1)
        self.directory.cleanup()

    def _peer_pid(self, connection: socket.socket) -> int:
        if sys.platform == "darwin":
            try:
                # Darwin LOCAL_PEERPID = 2 at SOL_LOCAL = 0.
                return int.from_bytes(
                    connection.getsockopt(0, 2, 4), sys.byteorder)
            except OSError:
                return 0
        peercred = getattr(socket, "SO_PEERCRED", None)
        if peercred is None:
            return 0
        try:
            import struct
            return int(struct.unpack(
                "3i", connection.getsockopt(socket.SOL_SOCKET, peercred, 12))[0])
        except (OSError, struct.error):
            return 0

    def _strict_fence_check(self) -> bool:
        verify = getattr(self.controller, "verify_current_fence", None)
        return bool(callable(verify) and verify({
            "reason": "headless_amp_authorize",
            "runtimeSessionId": self.runtime_session_id,
        }))

    @staticmethod
    def _trusted_wrapper_process(pid: int) -> bool:
        command = _process_command(pid)
        if not command:
            return False
        try:
            tokens = __import__("shlex").split(command)
        except ValueError:
            return False
        trusted = str(
            (Path(__file__).resolve().parents[1]
             / "bootstrap" / "amp_safe.py").resolve())
        for index, token in enumerate(tokens):
            if os.path.realpath(token) != trusted:
                continue
            return bool(index >= 1 and tokens[index - 1] == "-I")
        return False

    @staticmethod
    def _project_target(value: Any) -> Optional[tuple[str, str, str]]:
        if not isinstance(value, Mapping):
            return None
        project_id = str(value.get("projectId") or "")
        pop_code = str(value.get("popCode") or "")
        pop_version = str(value.get("popVersion") or "")
        try:
            valid_version = bool(
                re.fullmatch(r"\d{4}-\d{2}-\d{2}", pop_version)
                and datetime.strptime(pop_version, "%Y-%m-%d"))
        except ValueError:
            valid_version = False
        if (not project_id.isdecimal() or int(project_id) <= 0
                or not re.fullmatch(r"[A-Za-z][A-Za-z0-9._-]{0,127}", pop_code)
                or not valid_version):
            return None
        return project_id, pop_code.lower(), pop_version

    @staticmethod
    def _init_scope(value: Any) -> Optional[tuple[tuple[str, str], ...]]:
        if not isinstance(value, Mapping):
            return None
        keys = frozenset(value)
        if keys == {"projectId"}:
            project_id = str(value.get("projectId") or "")
            return (("projectId", project_id),) if (
                project_id.isdecimal() and int(project_id) > 0) else None
        if keys == {"popCode", "popVersion"}:
            pop_code = str(value.get("popCode") or "")
            pop_version = str(value.get("popVersion") or "")
            try:
                valid_version = bool(
                    re.fullmatch(r"\d{4}-\d{2}-\d{2}", pop_version)
                    and datetime.strptime(pop_version, "%Y-%m-%d"))
            except ValueError:
                valid_version = False
            if (re.fullmatch(r"[A-Za-z][A-Za-z0-9._-]{0,127}", pop_code)
                    and valid_version):
                return (("popCode", pop_code.lower()),
                        ("popVersion", pop_version))
        return None

    @staticmethod
    def _project_id_scope(value: Any) -> Optional[str]:
        if not isinstance(value, Mapping) or frozenset(value) != {"projectId"}:
            return None
        project_id = str(value.get("projectId") or "")
        return project_id if (project_id.isdecimal()
                              and int(project_id) > 0) else None

    def _authorize_decision(
            self, peer_pid: int, request: Mapping[str, Any],
    ) -> AmpBrokerDecision:
        if not self._enabled or peer_pid <= 0:
            return AmpBrokerDecision(False, "broker_not_enabled")
        if not _descends_from(
                peer_pid, self._bound_pid, self._bound_started_at):
            return AmpBrokerDecision(False, "peer_not_bound_descendant")
        if not self._trusted_wrapper_process(peer_pid):
            return AmpBrokerDecision(False, "peer_not_trusted_wrapper")
        if int(request.get("schemaVersion") or 0) != 2:
            return AmpBrokerDecision(False, "unsupported_request_schema")
        action = str(request.get("action") or "")
        if action not in self._ALLOWED_ACTIONS:
            return AmpBrokerDecision(False, "action_not_allowed")
        target = request.get("target")
        receipt = request.get("receipt")
        if not isinstance(target, Mapping) or not isinstance(receipt, Mapping):
            return AmpBrokerDecision(False, "target_or_receipt_missing")
        if (str(target.get("action") or "") != action
                or str(receipt.get("action") or "") != action
                or not receipt.get("allowed")):
            return AmpBrokerDecision(False, "action_receipt_mismatch")
        # A staging init establishes only local context. Its input scope has no
        # resolved project triple yet and must never establish broker pins.
        if action == "init":
            scope = self._init_scope(target.get("project"))
            if (scope is None or self._init_scope(receipt.get("project")) != scope
                    or str(target.get("repoMode") or "") != "local-context"
                    or str(receipt.get("repoMode") or "") != "local-context"):
                return AmpBrokerDecision(False, "init_scope_receipt_invalid")
            if not self._strict_fence_check():
                return AmpBrokerDecision(False, "strict_fence_denied")
            return AmpBrokerDecision(True, "authorized_local_context")

        project_identity = self._project_target(target.get("project"))
        if project_identity is None and action == "branch-create":
            project_id = self._project_id_scope(target.get("project"))
            project_identity = ((project_id, "", "")
                                if project_id is not None else None)
        if project_identity is None:
            return AmpBrokerDecision(False, "project_target_invalid")
        receipt_identity = self._project_target(receipt.get("project"))
        if receipt_identity is None and action == "branch-create":
            receipt_project_id = self._project_id_scope(receipt.get("project"))
            receipt_identity = ((receipt_project_id, "", "")
                                if receipt_project_id is not None else None)
        if receipt_identity != project_identity:
            return AmpBrokerDecision(False, "project_receipt_mismatch")
        branch = str(target.get("branch") or "")
        publish_status = (
            action == "publish" and str(target.get("publishKind") or "") == "status")
        feature_branch = re.fullmatch(
            r"feature/[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?"
            r"(?:/[A-Za-z0-9](?:[A-Za-z0-9._-]*[A-Za-z0-9])?)*", branch)
        if not publish_status and feature_branch is None:
            return AmpBrokerDecision(False, "branch_target_invalid")
        repo = str(request.get("repo") or "")
        try:
            resolved = Path(repo).expanduser().resolve(strict=True)
        except (OSError, RuntimeError, ValueError):
            return AmpBrokerDecision(False, "repo_unresolvable")
        if not resolved.is_dir():
            return AmpBrokerDecision(False, "repo_not_directory")
        repo_mode = str(target.get("repoMode") or "")
        receipt_mode = str(receipt.get("repoMode") or "")
        if repo_mode != receipt_mode:
            return AmpBrokerDecision(False, "repo_mode_mismatch")
        if publish_status:
            if repo_mode != "project":
                return AmpBrokerDecision(False, "publish_status_repo_mode_invalid")
            if not self._strict_fence_check():
                return AmpBrokerDecision(False, "strict_fence_denied")
            return AmpBrokerDecision(True, "authorized_publish_status")

        identity: Optional[tuple[str, int, int]] = None
        if repo_mode == "model":
            if (str(receipt.get("repo") or "") != os.fspath(resolved)
                    or re.fullmatch(r"[0-9a-f]{40,64}", str(
                        receipt.get("baseline") or "")) is None
                    or str(receipt.get("branch") or "") != branch):
                return AmpBrokerDecision(False, "canonical_baseline_receipt_invalid")
            try:
                stat = resolved.stat()
                identity = (os.fspath(resolved), int(stat.st_dev), int(stat.st_ino))
            except OSError:
                return AmpBrokerDecision(False, "repo_identity_unavailable")
        elif (repo_mode != "project"
              or (action != "branch-create" and not publish_status)):
            return AmpBrokerDecision(False, "repo_mode_not_allowed_for_action")

        if not self._strict_fence_check():
            return AmpBrokerDecision(False, "strict_fence_denied")
        # Validate the complete candidate before pinning. A malformed request,
        # wrong cwd, or stale fence therefore cannot poison later requests.
        with self._pin_lock:
            if not self._strict_fence_check():
                return AmpBrokerDecision(False, "strict_fence_denied_before_pin")
            if self._project_identity is not None:
                pinned = self._project_identity
                if (pinned[0] != project_identity[0]
                        or (pinned[1:] != ("", "")
                            and project_identity[1:] != ("", "")
                            and pinned[1:] != project_identity[1:])):
                    return AmpBrokerDecision(False, "project_identity_mismatch")
            if (identity is not None and self._repo_identity is not None
                    and self._repo_identity != identity):
                return AmpBrokerDecision(False, "repo_identity_mismatch")
            if (self._branch_identity is not None
                    and self._branch_identity != branch):
                return AmpBrokerDecision(False, "branch_identity_mismatch")
            if (self._project_identity is None
                    or (self._project_identity[1:] == ("", "")
                        and project_identity[1:] != ("", ""))):
                self._project_identity = project_identity
            if identity is not None and self._repo_identity is None:
                self._repo_identity = identity
            if self._branch_identity is None:
                self._branch_identity = branch
        return AmpBrokerDecision(True, "authorized")

    def _authorize(self, peer_pid: int, request: Mapping[str, Any]) -> bool:
        return self._authorize_decision(peer_pid, request).allowed

    def _serve(self) -> None:
        while not self._closed.is_set():
            try:
                connection, _ = self._socket.accept()
            except socket.timeout:
                continue
            except OSError:
                break
            with connection:
                decision = AmpBrokerDecision(False, "invalid_request")
                value: Any = None
                try:
                    raw = connection.recv(4097)
                    if not raw:
                        decision = AmpBrokerDecision(False, "peer_probe")
                    elif len(raw) <= 4096:
                        value = json.loads(raw.decode("utf-8"))
                        if isinstance(value, Mapping):
                            decision = self._authorize_decision(
                                self._peer_pid(connection), value)
                except (OSError, UnicodeError, ValueError):
                    decision = AmpBrokerDecision(False, "request_decode_failed")
                if not decision.allowed and decision.reason != "peer_probe":
                    LOG.warning(
                        "headless AMP authorization denied reason=%s action=%s repo=%s",
                        decision.reason,
                        str(value.get("action") or "") if isinstance(
                            value, Mapping) else "",
                        str(value.get("repo") or "") if isinstance(
                            value, Mapping) else "")
                try:
                    connection.sendall((json.dumps({
                        "allowed": decision.allowed,
                        "reason": decision.reason,
                    }, ensure_ascii=True, sort_keys=True,
                        separators=(",", ":")) + "\n").encode("utf-8"))
                except OSError:
                    pass


@dataclass(frozen=True)
class ExecutionResult:
    stdout: str
    stderr: str
    returncode: int
    timed_out: bool = False


class ProcessGuardian:
    """Launch and terminate a process group behind a fenced start gate."""

    def __init__(self, guard_script: Optional[Path] = None):
        self.guard_script = guard_script or Path(__file__).with_name(
            "task_process_guard.py")

    @staticmethod
    def terminate(process: Any, *, wait_seconds: float = 5.0,
                  close_streams: bool = True) -> None:
        try:
            os.killpg(os.getpgid(process.pid), signal.SIGKILL)
        except (ProcessLookupError, OSError):
            try:
                process.kill()
            except Exception:  # noqa: BLE001
                pass
        try:
            process.wait(timeout=wait_seconds)
        except Exception:  # noqa: BLE001
            pass
        if close_streams:
            for stream in (getattr(process, "stdout", None),
                           getattr(process, "stderr", None)):
                if stream is not None:
                    try:
                        stream.close()
                    except Exception:  # noqa: BLE001
                        pass

    def spawn(self, argv: Sequence[str], cwd: Path,
              on_spawn: SpawnCallback,
              env: Optional[Mapping[str, str]] = None) -> Tuple[Any, int]:
        if on_spawn is None:
            raise ValueError("guarded process requires an on_spawn binder")
        gate_read, gate_write = os.pipe()
        sentinel_read, sentinel_write = os.pipe()
        process = None
        try:
            command = [
                sys.executable,
                str(self.guard_script),
                "--gate-fd", str(gate_read),
                "--sentinel-fd", str(sentinel_read),
                "--grace-seconds",
                os.environ.get("JARVIS_TASK_GUARD_GRACE_SEC", "2"),
                "--",
            ] + list(argv)
            process = subprocess.Popen(
                command, cwd=cwd, text=True, stdin=subprocess.DEVNULL,
                stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                start_new_session=True, pass_fds=(gate_read, sentinel_read),
                env=env)
        except Exception:
            for fd in (gate_read, gate_write, sentinel_read, sentinel_write):
                try:
                    os.close(fd)
                except OSError:
                    pass
            raise
        finally:
            if process is not None:
                for fd in (gate_read, sentinel_read):
                    try:
                        os.close(fd)
                    except OSError:
                        pass

        try:
            # The real command remains blocked until the fenced process binding
            # succeeds.  Closing the gate without a grant makes the guard exit.
            on_spawn(process)
            try:
                os.write(gate_write, b"G")
            except BrokenPipeError:
                # The guard child already exited before the gate was granted —
                # this happens when the bound session lost ownership
                # (stale_fence) and stop_process fired mid-spawn.  The real
                # command never execed, so let run_buffered observe the guard's
                # empty exit (classified as a retryable no_result) instead of
                # surfacing a terminal orchestrator_exception for what is an
                # expected ownership handoff.
                pass
        except Exception:
            for fd in (gate_write, sentinel_write):
                try:
                    os.close(fd)
                except OSError:
                    pass
            self.terminate(process)
            raise
        finally:
            try:
                os.close(gate_write)
            except OSError:
                pass
        return process, sentinel_write


class ExecutionRuntime:
    """Run one buffered process with uniform timeout and cleanup behavior."""

    def __init__(self, guardian: Optional[ProcessGuardian] = None):
        self.guardian = guardian or ProcessGuardian()

    def run_buffered(self, argv: Sequence[str], cwd: Path, *, timeout: float,
                     on_spawn: Optional[SpawnCallback] = None,
                     guarded: bool = False,
                     guarded_spawn: Optional[GuardedSpawn] = None,
                     env: Optional[Mapping[str, str]] = None) -> ExecutionResult:
        sentinel_write = None
        if guarded:
            spawn = guarded_spawn or self.guardian.spawn
            process, sentinel_write = spawn(argv, cwd, on_spawn, env)
        else:
            process = subprocess.Popen(
                list(argv), cwd=cwd, text=True, stdin=subprocess.DEVNULL,
                stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                start_new_session=True, env=env)
            if on_spawn is not None:
                try:
                    on_spawn(process)
                except Exception:  # noqa: BLE001 - ephemeral observability is best effort
                    pass

        try:
            try:
                stdout, stderr = process.communicate(timeout=timeout)
            except subprocess.TimeoutExpired:
                self.guardian.terminate(process, close_streams=False)
                try:
                    stdout, stderr = process.communicate(timeout=5)
                except Exception:  # noqa: BLE001
                    stdout, stderr = "", ""
                return ExecutionResult(
                    stdout or "", stderr or "",
                    int(getattr(process, "returncode", -signal.SIGKILL) or 0),
                    timed_out=True)
            return ExecutionResult(
                stdout or "", stderr or "",
                int(getattr(process, "returncode", 0) or 0))
        finally:
            if sentinel_write is not None:
                try:
                    os.close(sentinel_write)
                except OSError:
                    pass


DEFAULT_PROCESS_GUARDIAN = ProcessGuardian()
DEFAULT_EXECUTION_RUNTIME = ExecutionRuntime(DEFAULT_PROCESS_GUARDIAN)


LOG = logging.getLogger(__name__)
REPO_ROOT = Path(__file__).resolve().parents[1]
ClaudeResult = __import__("collections").namedtuple(
    "ClaudeResult", "text is_error subtype")


def claude_bin() -> str:
    configured = os.environ.get("CLAUDE_BIN")
    if configured:
        return configured
    local = Path.home() / ".local" / "bin" / "claude"
    return str(local) if local.exists() else "claude"


def jarvis_root() -> str:
    return os.environ.get("JARVIS_ROOT") or str(REPO_ROOT)


def a1_command_env(terraform: bool = False) -> dict[str, str]:
    env = os.environ.copy()
    for key in ("JARVIS_A1_IDENTITY", "JARVIS_A1_STRICT",
                "JARVIS_AONE_WRITE_POLICY"):
        env.pop(key, None)
    if terraform:
        env["JARVIS_A1_IDENTITY"] = "terraform-rd"
        env["JARVIS_A1_STRICT"] = "1"
    return env


def _probe_settings(path: str, timeout: float = 5) -> bool:
    if not os.path.isfile(path):
        return False
    try:
        with open(path, encoding="utf-8") as stream:
            env = json.load(stream)["env"]
        request = Request(
            env["ANTHROPIC_BASE_URL"].rstrip("/") + "/v1/messages",
            data=json.dumps({
                "model": env["ANTHROPIC_MODEL"].split("[")[0],
                "max_tokens": 1,
                "messages": [{"role": "user", "content": "."}],
            }).encode(),
            headers={
                "authorization": "Bearer " + env["ANTHROPIC_AUTH_TOKEN"],
                "anthropic-version": "2023-06-01",
                "user-agent": "claude-cli/1.0",
                "content-type": "application/json",
            },
        )
        with urlopen(request, timeout=timeout) as response:
            return response.status // 100 == 2
    except Exception:  # noqa: BLE001
        return False


def _settings_candidates(member: str) -> list[str]:
    return [os.path.expanduser(value.strip())
            for value in member.split(",") if value.strip()]


def _resolve_settings(
    chain: str,
    probe_settings: Callable[[str], bool] = _probe_settings,
) -> str:
    candidates = _settings_candidates(chain)
    if len(candidates) < 2:
        return candidates[0] if candidates else chain
    return next((path for path in candidates if probe_settings(path)),
                candidates[-1])


def _provider_route_file(session_id: str) -> Path:
    root = Path(os.environ.get("JARVIS_PROVIDER_ROUTE_DIR")
                or Path.home() / ".cache" / "jarvis" / "provider-routes")
    return root / ("%s.json" % hashlib.sha256(
        str(session_id).encode()).hexdigest())


def _settings_model(path: str) -> Optional[str]:
    try:
        with open(path, encoding="utf-8") as stream:
            model = (json.load(stream).get("env") or {}).get("ANTHROPIC_MODEL")
        return str(model).split("[")[0] if model else None
    except Exception:  # noqa: BLE001
        return None


def _load_provider_route(session_id: str, lane: str) -> Optional[str]:
    """Return the pinned settings path, or None. See :func:`_load_route_record`."""
    record = _load_route_record(session_id, lane)
    return record.get("settingsPath") if record else None


def _load_route_record(session_id: str, lane: str) -> Optional[dict]:
    """Return the full route record, including any failover bookkeeping."""
    if not session_id:
        return None
    try:
        route = json.loads(
            _provider_route_file(session_id).read_text(encoding="utf-8"))
        if (route.get("schemaVersion") == 1
                and route.get("sessionId") == str(session_id)
                and route.get("lane") == lane
                and isinstance(route.get("settingsPath"), str)
                and route.get("settingsPath")):
            return route
    except (OSError, TypeError, ValueError):
        pass
    return None


def _persist_provider_route(
    session_id: str, lane: str, selected: str, *,
    failover: bool = False) -> bool:
    """Pin or update the provider route for one session.

    A new selection writes a fresh record. ``failover=True`` updates an existing
    record in place: it swaps the settings path and model, stamps a new
    ``selectedAt`` and records ``failoverFrom`` so the change is auditable, while
    clearing the ``firstFailedAt`` that armed the failover. The route file is the
    per-session source of truth, so updating it is how a resumed session learns
    to use the new provider.
    """
    if not session_id or not selected:
        return True
    target = _provider_route_file(session_id)
    try:
        target.parent.mkdir(parents=True, exist_ok=True, mode=0o700)
        record = {
            "schemaVersion": 1,
            "sessionId": str(session_id),
            "lane": lane,
            "settingsPath": str(selected),
            "model": _settings_model(selected),
            "selectedAt": datetime.now(timezone.utc).isoformat(),
        }
        if failover:
            previous = _load_route_record(session_id, lane) or {}
            if previous.get("settingsPath") and previous.get("settingsPath") != selected:
                record["failoverFrom"] = {
                    "settingsPath": previous.get("settingsPath"),
                    "model": previous.get("model"),
                }
        temporary = target.with_name(".%s.%s.tmp" % (target.name, os.getpid()))
        temporary.write_text(json.dumps(record, ensure_ascii=False),
                              encoding="utf-8")
        temporary.chmod(0o600)
        temporary.replace(target)
        return True
    except OSError as exc:
        LOG.warning("provider route persist failed session=%s: %s",
                    session_id, exc)
        return False


def _arm_provider_route_failure(session_id: str, lane: str) -> bool:
    """Record the first time the pinned provider failed health check.

    Idempotent: only writes on the first failure, so the backoff window is
    measured from the genuine first blip, not the most recent retry.
    """
    target = _provider_route_file(session_id)
    try:
        route = _load_route_record(session_id, lane) or {}
        if route.get("firstFailedAt"):
            return True
        route["firstFailedAt"] = datetime.now(timezone.utc).isoformat()
        target.parent.mkdir(parents=True, exist_ok=True, mode=0o700)
        temporary = target.with_name(".%s.%s.tmp" % (target.name, os.getpid()))
        temporary.write_text(json.dumps(route, ensure_ascii=False),
                             encoding="utf-8")
        temporary.chmod(0o600)
        temporary.replace(target)
        return True
    except OSError as exc:
        LOG.warning("provider route arm failed session=%s: %s", session_id, exc)
        return False


def _infer_provider_route(session_id: str,
                          candidates: Sequence[str]) -> Optional[str]:
    observed = None
    try:
        transcripts = list(
            (Path.home() / ".claude" / "projects").glob(
                "*/%s.jsonl" % session_id))
    except OSError:
        transcripts = []
    for transcript in transcripts:
        try:
            with transcript.open(encoding="utf-8") as stream:
                for line in stream:
                    try:
                        message = json.loads(line).get("message") or {}
                    except (TypeError, ValueError):
                        continue
                    model = message.get("model") if isinstance(
                        message, dict) else None
                    if model and model != "<synthetic>":
                        observed = str(model).split("[")[0]
        except OSError:
            continue
    matches = [path for path in candidates
               if observed and _settings_model(path) == observed]
    return matches[0] if len(matches) == 1 else None


# How long to keep retrying the original provider before a health-check failure
# is treated as something other than a transient blip. The HeadlessRuntime retry
# loop already backs off between attempts; this window lets one or two of those
# retries land before failover arms, so a few-second gateway hiccup self-heals
# instead of forcing a cross-provider resume.
_PROVIDER_FAILURE_BACKOFF_SECONDS = 60


def _failover_settings(
    member: str, original: str, probe_settings: Callable[[str], bool],
) -> Optional[str]:
    """A healthy provider to resume on, preferring the same model family.

    Same-family first keeps the resume context on a model that already produced
    it -- a qwen to qwen or glm to glm swap is safer than qwen to glm. Falls
    back to any healthy candidate when no same-family match exists, since
    stranding the Task is worse than a cross-family resume.
    """
    candidates = [c for c in _settings_candidates(member)
                  if c != original and probe_settings(c)]
    if not candidates:
        return None
    original_model = _settings_model(original) or ""
    family = original_model.split("-")[0] if original_model else ""
    same_family = ([c for c in candidates
                    if (_settings_model(c) or "").split("-")[0] == family]
                   if family else [])
    return same_family[0] if same_family else candidates[0]


def _parse_route_timestamp(value) -> Optional[datetime]:
    if not value or not isinstance(value, str):
        return None
    try:
        return datetime.fromisoformat(str(value))
    except (TypeError, ValueError):
        return None


def _select_provider_settings(
    member: str,
    session_id: str,
    terraform: bool,
    resume: bool,
    probe_settings: Callable[[str], bool] = _probe_settings,
) -> str:
    lane = "terraform" if terraform else "default"
    record = _load_route_record(session_id, lane)
    selected = record.get("settingsPath") if record else None
    if selected:
        if not resume or probe_settings(selected):
            return selected
        # Resume with a provider that just failed its health probe. The old
        # code raised immediately, so a multi-second gateway blip burned the
        # whole retry budget into RECOVERY_REQUIRED because every retry hit
        # the same raise before the process even started. Give the original
        # provider one backoff window first -- the probe is a 5s /v1/messages
        # call and the real CLI has its own connection retry -- then fail
        # over to a healthy candidate and pin it so subsequent resumes use it.
        first_failed = _parse_route_timestamp(record.get("firstFailedAt"))
        if first_failed is None:
            _arm_provider_route_failure(session_id, lane)
            return selected
        if ((datetime.now(timezone.utc) - first_failed).total_seconds()
                < _PROVIDER_FAILURE_BACKOFF_SECONDS):
            return selected
        replacement = _failover_settings(member, selected, probe_settings)
        if replacement:
            _persist_provider_route(session_id, lane, replacement, failover=True)
            LOG.warning(
                "provider route failover session=%s %s -> %s (model %s -> %s)",
                session_id, selected, replacement,
                _settings_model(selected), _settings_model(replacement))
            return replacement
        raise RuntimeError(
            "model_provider_error: original provider route failed health "
            "check and no healthy candidate is available for failover")
    candidates = _settings_candidates(member)
    if resume:
        selected = (candidates[0] if len(candidates) == 1
                    else _infer_provider_route(session_id, candidates))
        if not selected:
            raise RuntimeError(
                "model_provider_error: original provider route is unknown; "
                "wait for an operator to start a new session")
    else:
        selected = _resolve_settings(member, probe_settings)
    if not _persist_provider_route(session_id, lane, selected):
        raise RuntimeError(
            "model_provider_error: original provider route could not be pinned; "
            "refusing an unsafe resumable launch")
    return selected


def jarvis_cmd(session_id: str = "", terraform: bool = False,
               resume: bool = False,
               probe_settings: Callable[[str], bool] = _probe_settings) -> list[str]:
    override = os.environ.get("JARVIS_CC")
    if override:
        return [override]
    default = str(Path.home() / ".claude" / "idea_settings.json")
    raw = ((os.environ.get("JARVIS_SETTINGS_TF")
            or os.environ.get("JARVIS_SETTINGS") or default)
           if terraform else os.environ.get("JARVIS_SETTINGS") or default)
    pool = [value.strip() for value in raw.split(":") if value.strip()]
    member = (pool[int(hashlib.md5(session_id.encode()).hexdigest(), 16)
                   % len(pool)]
              if len(pool) > 1 and session_id else pool[0])
    return [
        claude_bin(), "--settings",
        _select_provider_settings(
            member, session_id, terraform, resume, probe_settings),
        "--permission-mode", "bypassPermissions",
    ]


def session_file(session_id: str) -> Path:
    slug = re.sub(r"[^a-zA-Z0-9]", "-", os.path.realpath(jarvis_root()))
    return (Path.home() / ".claude" / "projects" / slug
            / ("%s.jsonl" % session_id))


def session_file_exists(session_id: str) -> bool:
    try:
        return session_file(session_id).exists()
    except Exception:  # noqa: BLE001
        return False


def session_progress_excerpt(
    session_id: str,
    max_bytes: int = 512 * 1024,
    max_chars: int = 12000,
    sanitizer: Optional[Callable[..., str]] = None,
    transcript_for: Callable[[str], Path] = session_file,
) -> Optional[str]:
    root = transcript_for(session_id)
    paths = [root]
    try:
        subagents = list((root.with_suffix("") / "subagents").glob("*.jsonl"))
        subagents.sort(key=lambda path: path.stat().st_mtime, reverse=True)
        paths.extend(subagents[:8])
    except OSError:
        pass
    entries = []
    per_file_bytes = max(4096, int(max_bytes) // len(paths))
    for path_index, path in enumerate(paths):
        try:
            with path.open("rb") as handle:
                size = handle.seek(0, os.SEEK_END)
                handle.seek(max(0, size - per_file_bytes))
                raw = handle.read().decode("utf-8", errors="replace")
        except (OSError, TypeError, ValueError):
            continue
        if size > per_file_bytes:
            raw = raw.split("\n", 1)[-1]
        for line_index, line in enumerate(raw.splitlines()):
            try:
                event = json.loads(line)
            except (TypeError, ValueError):
                continue
            message = event.get("message") if isinstance(event, dict) else None
            if not isinstance(message, dict) or message.get("role") != "assistant":
                continue
            content = message.get("content")
            parts = []
            if isinstance(content, str) and content.strip():
                parts.append(content.strip())
            elif isinstance(content, list):
                for block in content:
                    if not isinstance(block, dict):
                        continue
                    if (block.get("type") == "text"
                            and str(block.get("text") or "").strip()):
                        parts.append(str(block["text"]).strip())
                    elif block.get("type") == "tool_use":
                        name = re.sub(
                            r"[^a-zA-Z0-9_.:-]", "",
                            str(block.get("name") or "tool"))[:80]
                        parts.append("[执行工具 %s]" % (name or "tool"))
            if parts:
                entries.append((str(event.get("timestamp") or ""), path_index,
                                line_index, "\n".join(parts)))
    if not entries:
        return None
    entries.sort(key=lambda entry: entry[:3])
    excerpt = "\n\n".join(entry[3] for entry in entries[-24:])
    excerpt = re.sub(
        r"\x1b(?:\[[0-?]*[ -/]*[@-~]|\][^\x07]*(?:\x07|\x1b\\))",
        "", excerpt)
    excerpt = "".join(char for char in excerpt
                      if char in "\n\r\t" or ord(char) >= 32)
    if len(excerpt) > max_chars:
        excerpt = "…" + excerpt[-max_chars + 1:]
    if sanitizer is not None:
        excerpt = sanitizer(excerpt, limit=max_chars)
    return excerpt or None


def _headless_exec_command(session_id: str,
                           command: Sequence[str]) -> list[str]:
    if not session_id or not command:
        raise ValueError("headless session_id and command are required")
    manager = REPO_ROOT / "bootstrap" / "jarvis-interactive-worker.py"
    return [
        "/usr/bin/python3", "-I", str(manager), "exec-headless",
        "--session-id", str(session_id), "--client", "claude", "--",
    ] + list(command)


def _normalized_failure_subtype(
    text: str,
    subtype: str,
    is_error: bool = True,
) -> str:
    current = str(subtype or "").strip() or ("error" if is_error else "success")
    lowered = str(text or "").lower()
    provider_failure = (
        "model_provider_error" in lowered
        or "original provider route" in lowered
        or "cross-provider resume" in lowered
        or bool(re.search(
            r"(?:模型(?:提供方|供应商|网关).{0,24}(?:错误|失败|异常|不可用)"
            r"|(?:model|llm|claude)[ _-]*(?:provider|gateway).{0,32}"
            r"(?:error|failed|failure|unavailable|invalid|timeout))",
            str(text or ""), re.IGNORECASE)))
    if provider_failure:
        return "model_provider_error"
    if current.lower() == "success" and is_error:
        return "execution_error"
    return current


def classify_result(out: str, err: str, returncode: int) -> ClaudeResult:
    last = None
    for raw in (out or "").splitlines():
        try:
            value = json.loads(raw.strip())
        except (TypeError, ValueError):
            continue
        if isinstance(value, dict) and value.get("type") == "result":
            last = value
    if last is not None:
        text = last.get("result")
        text = text if isinstance(text, str) else ""
        is_error = bool(last.get("is_error")) or returncode != 0
        subtype = last.get("subtype") or (
            "error" if is_error else "success")
        return ClaudeResult(
            text, is_error,
            _normalized_failure_subtype(text, subtype, is_error))
    if returncode != 0:
        text = ((err or "").strip().splitlines()[-1:] or ["unknown"])[0]
        return ClaudeResult(
            text, True, _normalized_failure_subtype(text, "no_result"))
    return ClaudeResult("", False, "no_result")


def run_claude_buffered(
    text: str,
    session_id: str,
    resume: bool,
    timeout: Optional[float] = None,
    on_spawn: Optional[SpawnCallback] = None,
    terraform: bool = False,
    guarded: bool = False,
    execution_runtime: Optional[ExecutionRuntime] = None,
    command_builder: Callable[..., list[str]] = jarvis_cmd,
    headless_wrapper: Callable[[str, Sequence[str]], list[str]] = (
        _headless_exec_command),
) -> ClaudeResult:
    timeout = timeout or int(
        os.environ.get("JARVIS_DISPATCH_TIMEOUT", "86400"))
    argv = command_builder(
        session_id, terraform=terraform, resume=resume
    ) + ["-p", text, "--output-format", "json"]
    argv += ["--resume", session_id] if resume else [
        "--session-id", session_id]
    environment = a1_command_env(terraform=terraform)
    broker: Optional[HeadlessAmpBroker] = None
    spawn_callback = on_spawn
    if guarded and on_spawn is not None:
        callback_owner = getattr(on_spawn, "__self__", None)
        controller = getattr(callback_owner, "controller", callback_owner)
        if callable(getattr(controller, "verify_current_fence", None)):
            broker = HeadlessAmpBroker(controller, session_id)
            environment[_HEADLESS_AMP_BROKER_ENV] = broker.path

            def bind_and_enable(process: Any) -> None:
                on_spawn(process)
                # For TaskAoneBookend this is deliberately after both the
                # control-plane PID bind and the Aone claimed-tag write.
                broker.enable(process)

            spawn_callback = bind_and_enable
    try:
        execution = (execution_runtime or DEFAULT_EXECUTION_RUNTIME).run_buffered(
            headless_wrapper(session_id, argv),
            Path(jarvis_root()), timeout=timeout, on_spawn=spawn_callback,
            guarded=guarded, env=environment)
    finally:
        if broker is not None:
            broker.close()
    if execution.timed_out:
        return ClaudeResult(execution.stdout or "", True, "timeout")
    return classify_result(
        execution.stdout, execution.stderr, execution.returncode)


class EphemeralExecutor:
    """Bounded disposable-job executor sharing capacity with durable Tasks."""

    def __init__(self, max_workers=None, queue_max=None,
                 capacity_manager=None, execution_runtime=None):
        self.max_workers = int(max_workers if max_workers is not None
                               else os.environ.get("JARVIS_DISPATCH_MAX", "3"))
        self.queue_max = int(queue_max if queue_max is not None
                             else os.environ.get("JARVIS_DISPATCH_QUEUE_MAX", "20"))
        self._executor = ThreadPoolExecutor(
            max_workers=max(1, self.max_workers),
            thread_name_prefix="dispatch")
        self.capacity_manager = capacity_manager or CapacityManager(self.max_workers)
        if not all(callable(getattr(self.capacity_manager, name, None))
                   for name in ("acquire", "available_slots")):
            raise TypeError("capacity_manager must provide acquire/available_slots")
        self.execution_runtime = execution_runtime or DEFAULT_EXECUTION_RUNTIME
        self._active = {}
        self._lock = threading.Lock()
        self._closed = False
        self._watchdog_interval = float(
            os.environ.get("JARVIS_DISPATCH_WATCHDOG_INTERVAL", "60"))
        self._watchdog_threshold = float(
            os.environ.get("JARVIS_DISPATCH_WATCHDOG_THRESHOLD", "3600"))
        self._watchdog_thread = threading.Thread(
            target=self._watchdog_loop, daemon=True,
            name="EphemeralExecutorWatchdog")
        self._watchdog_thread.start()

    def status(self, item_id, force=False):
        with self._lock:
            if str(item_id) in self._active:
                return False, "active"
            if len(self._active) >= self.max_workers + self.queue_max:
                return False, "queue_full"
        return True, "ok"

    def active_count(self):
        with self._lock:
            return len(self._active)

    def active_ids(self):
        with self._lock:
            return sorted(self._active)

    def free_slots(self):
        with self._lock:
            queued = any(entry.get("started") is None
                         for entry in self._active.values())
            running = sum(entry.get("started") is not None
                          for entry in self._active.values())
        return 0 if queued else max(0, min(
            self.max_workers - running,
            self.capacity_manager.available_slots()))

    def submit(self, item_id, work, *, notify=None, force=False, kind="ticket",
               project=None, terraform=False):
        item_id = str(item_id)
        with self._lock:
            if self._closed:
                return False, "closing"
            if item_id in self._active:
                return False, "active"
            if len(self._active) >= self.max_workers + self.queue_max:
                return False, "queue_full"
            self._active[item_id] = {
                "queuedAt": time.time(), "started": None, "kind": kind,
                "future": None, "project": project, "proc": None,
                "permit": None, "terraform": bool(terraform),
            }

        def run():
            permit = None
            try:
                while permit is None:
                    with self._lock:
                        if self._closed or item_id not in self._active:
                            return "cancelled"
                    permit = self.capacity_manager.acquire(
                        "ephemeral:%s" % item_id)
                    if permit is None:
                        time.sleep(0.05)
                with self._lock:
                    entry = self._active.get(item_id)
                    if entry is None or self._closed:
                        return "cancelled"
                    entry.update(permit=permit, started=time.time())
                return work()
            except Exception as exc:  # noqa: BLE001
                LOG.exception("EphemeralExecutor job #%s crashed: %s",
                              item_id, exc)
                if notify:
                    try:
                        notify("⚠️ #%s 后台处理异常: %s" % (item_id, exc))
                    except Exception:  # noqa: BLE001
                        pass
                return "error"
            finally:
                if permit is not None:
                    permit.release()
                self._terminate_worker(item_id)
                with self._lock:
                    self._active.pop(item_id, None)

        future = self._executor.submit(run)
        with self._lock:
            if item_id in self._active:
                self._active[item_id]["future"] = future
        return True, "dispatched"

    def set_proc(self, item_id, process):
        with self._lock:
            entry = self._active.get(str(item_id))
            if entry is not None:
                entry["proc"] = process

    def _terminate_worker(self, item_id):
        with self._lock:
            entry = self._active.get(str(item_id))
            process = entry.get("proc") if entry else None
        if process is not None:
            try:
                os.killpg(os.getpgid(process.pid), signal.SIGKILL)
            except (ProcessLookupError, PermissionError, OSError):
                pass

    def _watchdog_loop(self):
        while not self._closed:
            try:
                time.sleep(self._watchdog_interval)
                if not self._closed:
                    self._watchdog_tick()
            except Exception:  # noqa: BLE001
                if not self._closed:
                    LOG.exception("EphemeralExecutor watchdog tick failed")

    def _watchdog_tick(self):
        now = time.time()
        with self._lock:
            victims = [
                (item_id, entry, now - entry["started"])
                for item_id, entry in self._active.items()
                if entry.get("started") is not None
                and now - entry["started"] > self._watchdog_threshold
            ]
        for item_id, entry, age in victims:
            LOG.warning(
                "EphemeralExecutor slot #%s zombie (age=%ds, kind=%s)",
                item_id, int(age), entry.get("kind", "?"))
            process = entry.get("proc")
            if process is not None:
                try:
                    os.killpg(os.getpgid(process.pid), signal.SIGKILL)
                except (ProcessLookupError, PermissionError, OSError):
                    pass
            future = entry.get("future")
            if future is not None:
                future.cancel()
            permit = entry.get("permit")
            if permit is not None:
                permit.release()
            with self._lock:
                self._active.pop(item_id, None)

    def terminate_all(self, release_fn=None, grace=3):
        with self._lock:
            self._closed = True
        self.shutdown(wait=False, cancel_futures=True)

        def sweep(sig):
            with self._lock:
                snapshot = [
                    (item_id, entry.get("proc"), entry.get("project"),
                     entry.get("kind"), bool(entry.get("terraform")))
                    for item_id, entry in self._active.items()
                ]
            for _item_id, process, _project, _kind, _terraform in snapshot:
                if process is not None:
                    try:
                        os.killpg(os.getpgid(process.pid), sig)
                    except (ProcessLookupError, OSError):
                        pass
            return snapshot

        first = sweep(signal.SIGTERM)
        time.sleep(grace)
        second = sweep(signal.SIGKILL)
        claims = {}
        for item_id, process, project, kind, terraform in first + second:
            if process is not None and kind == "ticket" and project:
                claims[item_id] = (project, terraform)
        if release_fn:
            for item_id, (project, terraform) in claims.items():
                try:
                    release_fn(item_id, project, terraform=terraform)
                except Exception as exc:  # noqa: BLE001
                    LOG.warning("terminate_all: release #%s failed: %s",
                                item_id, exc)
        return sorted({
            item_id for item_id, process, *_rest in first + second
            if process is not None
        })

    def shutdown(self, wait=False, cancel_futures=False):
        try:
            self._executor.shutdown(
                wait=wait, cancel_futures=cancel_futures)
        except TypeError:  # pragma: no cover
            self._executor.shutdown(wait=wait)


__all__ = [
    "ClaudeResult", "DEFAULT_EXECUTION_RUNTIME", "DEFAULT_PROCESS_GUARDIAN",
    "EphemeralExecutor", "ExecutionResult", "ExecutionRuntime", "ProcessGuardian",
    "a1_command_env", "claude_bin", "classify_result", "jarvis_cmd",
    "jarvis_root", "run_claude_buffered", "session_file",
    "session_file_exists", "session_progress_excerpt",
]
