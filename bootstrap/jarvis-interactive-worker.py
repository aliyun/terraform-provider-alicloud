#!/usr/bin/env python3
"""Interactive Codex/Claude worker registration and direct-task lifecycle.

The ordinary bridge worker pulls MANAGED work.  An interactive worker must never
do that: it registers itself and heartbeats, then attaches only to the exact Aone
task requested by the user through ``POST /tasks/claim``.  A small detached
sidecar keeps both the worker and its current fenced session alive while the real
Codex/Claude host process exists.

State is local coordination data, not an authority.  It is written atomically
with mode 0600 and deliberately excludes credentials; the database fence remains
the ownership authority.
"""

from __future__ import annotations

import argparse
import contextlib
import fcntl
import hashlib
import json
import os
import re
import shlex
import socket
import subprocess
import sys
import tempfile
import time
from pathlib import Path
from typing import Any, Callable, Dict, Iterator, Mapping, Optional, Tuple, TypeVar


HERE = Path(__file__).resolve().parent
REPO_ROOT = HERE.parent
BRIDGE_DIR = REPO_ROOT / "bridge"
sys.path.insert(0, str(BRIDGE_DIR))

from jarvis_local_worker import _default_boot_id, make_worker_key  # noqa: E402
from jarvis_task_client import (  # noqa: E402
    AutomationAgentTaskClient,
    ControlPlaneConflict,
    ControlPlaneError,
    ControlPlaneUnavailable,
    StaleFence,
    TaskEnvelope,
)


CONFLICT_EXIT = 10
UNAVAILABLE_EXIT = 11
STATE_ERROR_EXIT = 12
TRANSITION_EXIT = 13
T = TypeVar("T")


def _nonblank(value: Any, name: str) -> str:
    text = str(value or "").strip()
    if not text:
        raise ValueError("%s must not be empty" % name)
    return text


def _field(value: Mapping[str, Any], *names: str) -> Any:
    for name in names:
        if name in value:
            return value[name]
    return None


def _safe_name(value: str) -> str:
    cleaned = re.sub(r"[^A-Za-z0-9_.-]+", "_", value).strip("._")
    digest = hashlib.sha256(value.encode("utf-8")).hexdigest()[:12]
    return "%s-%s" % ((cleaned[:48] or "session"), digest)


def _state_root() -> Path:
    configured = os.environ.get("JARVIS_INTERACTIVE_STATE_DIR", "").strip()
    if configured:
        return Path(configured).expanduser()
    home = Path(os.environ.get("HOME") or str(Path.home()))
    return home / ".cache" / "jarvis" / "interactive-workers"


def _state_path(client_name: str, session_id: str) -> Path:
    return _state_root() / ("%s-%s.json" % (
        _safe_name(client_name), _safe_name(session_id)))


def _ensure_private_dir(path: Path) -> None:
    path.mkdir(parents=True, exist_ok=True, mode=0o700)
    try:
        path.chmod(0o700)
    except OSError:
        pass


class StateStore:
    def __init__(self, path: Path):
        self.path = Path(path)
        self.lock_path = self.path.with_suffix(self.path.suffix + ".lock")

    @contextlib.contextmanager
    def locked(self) -> Iterator[None]:
        _ensure_private_dir(self.path.parent)
        fd = os.open(str(self.lock_path), os.O_CREAT | os.O_RDWR, 0o600)
        try:
            os.chmod(self.lock_path, 0o600)
            with os.fdopen(fd, "a+") as lock_file:
                fcntl.flock(lock_file.fileno(), fcntl.LOCK_EX)
                yield
                fcntl.flock(lock_file.fileno(), fcntl.LOCK_UN)
        except BaseException:
            try:
                os.close(fd)
            except OSError:
                pass
            raise

    def load_unlocked(self) -> Dict[str, Any]:
        try:
            raw = self.path.read_text(encoding="utf-8")
        except FileNotFoundError:
            return {}
        try:
            parsed = json.loads(raw)
        except (TypeError, ValueError) as exc:
            raise RuntimeError("invalid interactive worker state: %s" % self.path) from exc
        if not isinstance(parsed, dict):
            raise RuntimeError("interactive worker state must be an object")
        return parsed

    def load(self) -> Dict[str, Any]:
        with self.locked():
            return self.load_unlocked()

    def save_unlocked(self, state: Mapping[str, Any]) -> None:
        _ensure_private_dir(self.path.parent)
        fd, tmp_name = tempfile.mkstemp(
            prefix=".%s." % self.path.name, dir=str(self.path.parent))
        try:
            os.fchmod(fd, 0o600)
            with os.fdopen(fd, "w", encoding="utf-8") as handle:
                json.dump(dict(state), handle, ensure_ascii=False, sort_keys=True,
                          separators=(",", ":"))
                handle.write("\n")
                handle.flush()
                os.fsync(handle.fileno())
            os.replace(tmp_name, self.path)
            self.path.chmod(0o600)
        finally:
            try:
                os.unlink(tmp_name)
            except FileNotFoundError:
                pass

    def save(self, state: Mapping[str, Any]) -> None:
        with self.locked():
            self.save_unlocked(state)


def _client() -> AutomationAgentTaskClient:
    base_url = (os.environ.get("JARVIS_CONTROL_PLANE_BASE_URL", "").strip()
                or os.environ.get("JARVIS_HTML_REPORT_BASE_URL", "").strip()
                or "https://pre-agent.aliyun-inc.com")
    token = (os.environ.get("JARVIS_CONTROL_PLANE_TOKEN", "").strip()
             or os.environ.get("JARVIS_HTML_REPORT_TOKEN", "").strip())
    timeout = float(os.environ.get("JARVIS_CONTROL_PLANE_TIMEOUT", "10"))
    return AutomationAgentTaskClient(base_url, token, timeout=timeout)


def _retry_unavailable(call: Callable[[], T]) -> T:
    """Retry transactional/idempotent calls after a transient 5xx or transport loss."""
    attempts = max(1, int(os.environ.get("JARVIS_INTERACTIVE_RETRY_ATTEMPTS", "3")))
    delay = max(0.0, float(os.environ.get("JARVIS_INTERACTIVE_RETRY_DELAY", "0.2")))
    for attempt in range(attempts):
        try:
            return call()
        except ControlPlaneUnavailable:
            if attempt + 1 >= attempts:
                raise
            time.sleep(delay * (2 ** attempt))
    raise AssertionError("unreachable")


def _pid_alive(pid: Any) -> bool:
    try:
        value = int(pid)
        if value <= 0:
            return False
        os.kill(value, 0)
        return True
    except (TypeError, ValueError, OSError):
        return False


def _process_info(pid: int) -> Tuple[int, str]:
    try:
        result = subprocess.run(
            ["ps", "-o", "ppid=", "-o", "command=", "-p", str(pid)],
            capture_output=True, text=True, timeout=2, check=False)
    except (OSError, subprocess.SubprocessError):
        return 0, ""
    line = result.stdout.strip()
    if not line:
        return 0, ""
    parts = line.split(None, 1)
    try:
        parent = int(parts[0])
    except (ValueError, IndexError):
        parent = 0
    command = parts[1] if len(parts) > 1 else ""
    return parent, command


def _process_start_identity(pid: int) -> str:
    """Return a stable birth identity so PID reuse cannot extend an old lease."""
    try:
        result = subprocess.run(
            ["ps", "-o", "lstart=", "-p", str(pid)],
            capture_output=True, text=True, timeout=2, check=False)
    except (OSError, subprocess.SubprocessError):
        return ""
    return " ".join(result.stdout.split())


def _find_host_pid(client_name: str) -> Tuple[int, bool]:
    configured = os.environ.get("JARVIS_INTERACTIVE_HOST_PID", "").strip()
    if configured:
        configured_pid = int(configured)
        _parent, command = _process_info(configured_pid)
        return configured_pid, client_name.lower() in command.lower()
    needle = client_name.lower()
    pid = os.getppid()
    fallback = pid
    for depth in range(12):
        if pid <= 1:
            break
        parent, command = _process_info(pid)
        lower = command.lower()
        helper = ("jarvis-interactive-worker" in lower
                  or "run-interactive-worker-hook" in lower)
        if not helper and (re.search(r"(^|[/ ])claude([ /]|$)", lower)
                           if needle == "claude" else "codex" in lower):
            return pid, True
        if depth >= 2 and fallback == os.getppid():
            fallback = pid
        pid = parent
    # A short-lived hook shell is normally two ancestors below the real host.
    # The daemon also has the server-side 90s stale-worker reaper as a backstop.
    return fallback, False


def _host_alive(state: Mapping[str, Any]) -> bool:
    pid = state.get("hostPid")
    if not _pid_alive(pid):
        return False
    if not state.get("verifyHostCommand"):
        return False
    _parent, command = _process_info(int(pid))
    expected_start = str(state.get("hostProcessStartedAt") or "")
    return (str(state.get("client") or "").lower() in command.lower()
            and bool(expected_start)
            and _process_start_identity(int(pid)) == expected_start)


def _runtime_context() -> Tuple[str, str]:
    codex = os.environ.get("CODEX_THREAD_ID", "").strip()
    if codex:
        return "codex", codex
    claude = os.environ.get("CLAUDE_CODE_SESSION_ID", "").strip()
    if claude:
        return "claude", claude
    persisted_client = os.environ.get("JARVIS_INTERACTIVE_CLIENT", "").strip().lower()
    persisted_session = os.environ.get("JARVIS_INTERACTIVE_SESSION_ID", "").strip()
    if persisted_client in ("claude", "codex") and persisted_session:
        return persisted_client, persisted_session
    raise RuntimeError("no Claude/Codex interactive session context")


def _persist_claude_context(session_id: str) -> None:
    """Export hook identity to subsequent Claude tool subprocesses."""
    env_file = os.environ.get("CLAUDE_ENV_FILE", "").strip()
    if not env_file:
        return
    payload = (
        "export JARVIS_INTERACTIVE_CLIENT=%s\n"
        "export JARVIS_INTERACTIVE_SESSION_ID=%s\n" %
        (shlex.quote("claude"), shlex.quote(session_id))
    ).encode("utf-8")
    fd = os.open(env_file, os.O_WRONLY | os.O_APPEND | os.O_CREAT, 0o600)
    try:
        os.write(fd, payload)
    finally:
        os.close(fd)


def _current_store() -> StateStore:
    client_name, session_id = _runtime_context()
    store = StateStore(_state_path(client_name, session_id))
    if not store.path.exists():
        raise RuntimeError(
            "interactive worker is not registered; SessionStart hook may be untrusted or failed")
    return store


def _capabilities(state: Mapping[str, Any]) -> Dict[str, Any]:
    return {
        "dispatch": {"pull": False},
        "workerMode": "INTERACTIVE",
        "workerKey": state.get("workerKey"),
        "client": state.get("client"),
        "clientSessionId": state.get("clientSessionId"),
        "hostPid": state.get("hostPid"),
        "kinds": ["interactive-aone"],
    }


def _worker_payload(state: Mapping[str, Any], status: str) -> Dict[str, Any]:
    active = 1 if isinstance(state.get("current"), Mapping) else 0
    return {
        "workerKey": state["workerKey"],
        "host": state["host"],
        "bootId": state["bootId"],
        "processUuid": state["processUuid"],
        "version": state.get("version") or "interactive-v1",
        "status": status,
        "maxSlots": 1,
        "activeSessions": active,
        "freeSlots": 0 if active else 1,
        "capabilities": _capabilities(state),
    }


def _heartbeat_worker(client: AutomationAgentTaskClient,
                      state: Mapping[str, Any], status: str = "ACTIVE") -> None:
    payload = _worker_payload(state, status)
    # Worker heartbeat DTO does not repeat immutable identity fields.
    heartbeat = {
        "status": status,
        "capabilities": payload["capabilities"],
        "activeSessions": payload["activeSessions"],
        "freeSlots": payload["freeSlots"],
    }
    client.heartbeat_worker(
        state["workerKey"], heartbeat,
        request_id="jarvis-interactive-worker-heartbeat-%s" %
        hashlib.sha256((state["workerKey"] + str(time.time_ns())).encode()).hexdigest()[:24])


def _register(client: AutomationAgentTaskClient, state: Mapping[str, Any],
              status: str = "ACTIVE") -> None:
    client.register_worker(
        _worker_payload(state, status),
        request_id="jarvis-interactive-worker-register-%s" %
        hashlib.sha256(state["workerKey"].encode()).hexdigest()[:24])


def _mark_old_offline(client: AutomationAgentTaskClient,
                      old_state: Mapping[str, Any]) -> None:
    if not old_state.get("workerKey"):
        return
    try:
        _heartbeat_worker(client, old_state, "OFFLINE")
    except ControlPlaneError:
        pass


def _spawn_daemon(store: StateStore, expected_worker_key: str) -> int:
    log_path = store.path.with_suffix(".log")
    _ensure_private_dir(log_path.parent)
    log_fd = os.open(str(log_path), os.O_CREAT | os.O_APPEND | os.O_WRONLY, 0o600)
    env = os.environ.copy()
    try:
        process = subprocess.Popen(
            [sys.executable, str(Path(__file__).resolve()), "daemon",
             "--state", str(store.path), "--worker-key", expected_worker_key],
            stdin=subprocess.DEVNULL, stdout=log_fd, stderr=log_fd,
            start_new_session=True, close_fds=True, env=env)
    finally:
        os.close(log_fd)
    return int(process.pid)


def _ensure_daemon(store: StateStore, expected_worker_key: str) -> None:
    with store.locked():
        state = store.load_unlocked()
        if state.get("workerKey") != expected_worker_key:
            return
        # An arbitrary live ancestor is not a safe liveness authority.  Without
        # a verified Codex/Claude host process, a detached sidecar could renew a
        # task forever after the real client has exited (or after PID reuse).
        if not state.get("verifyHostCommand"):
            return
        if _pid_alive(state.get("daemonPid")):
            return
        state["daemonPid"] = _spawn_daemon(store, expected_worker_key)
        state["daemonStartedAt"] = int(time.time())
        store.save_unlocked(state)


def _hook_output(event_name: str, message: Optional[str] = None) -> None:
    output: Dict[str, Any] = {}
    if event_name == "SessionStart" and message:
        output["hookSpecificOutput"] = {
            "hookEventName": "SessionStart",
            "additionalContext": message,
        }
    print(json.dumps(output, ensure_ascii=False, separators=(",", ":")))


def hook(client_name: str, event: Mapping[str, Any]) -> int:
    event_name = str(event.get("hook_event_name") or "")
    session_id = _nonblank(event.get("session_id"), "session_id")
    store = StateStore(_state_path(client_name, session_id))
    cp = _client()

    if event_name == "SessionEnd":
        try:
            # SessionEnd can arrive after the same conversation has already been
            # resumed by a new client process.  Compute the ending process's
            # incarnation and only offline the matching state; never let a stale
            # event close a replacement worker.
            host_pid, verify_command = _find_host_pid(client_name)
            host_process_started_at = (
                _process_start_identity(host_pid) if verify_command else "")
            if verify_command and host_process_started_at:
                host = socket.gethostname()
                boot_id = _default_boot_id(host)
                process_uuid = hashlib.sha256(
                    ("%s|%s|%s|%s|%s" %
                     (client_name, session_id, boot_id, host_pid,
                      host_process_started_at)).encode()
                ).hexdigest()[:40]
                expected_worker_key = make_worker_key(host, boot_id, process_uuid)
                offline(store, cp, expected_worker_key)
            else:
                print("interactive worker SessionEnd warning: unverified host incarnation",
                      file=sys.stderr)
        except Exception as exc:  # SessionEnd must never trap the user in the client.
            print("interactive worker SessionEnd warning: %s" % type(exc).__name__,
                  file=sys.stderr)
        _hook_output(event_name)
        return 0

    if event_name != "SessionStart":
        _hook_output(event_name)
        return 0

    try:
        if client_name == "claude":
            _persist_claude_context(session_id)
        host_pid, verify_command = _find_host_pid(client_name)
        host_process_started_at = (
            _process_start_identity(host_pid) if verify_command else "")
        verify_command = bool(verify_command and host_process_started_at)
        host = socket.gethostname()
        boot_id = _default_boot_id(host)
        process_uuid = hashlib.sha256(
            ("%s|%s|%s|%s|%s" %
             (client_name, session_id, boot_id, host_pid,
              host_process_started_at)).encode()
        ).hexdigest()[:40]
        worker_key = make_worker_key(host, boot_id, process_uuid)
        old_state = store.load()
        same_incarnation = old_state.get("workerKey") == worker_key
        if old_state and not same_incarnation:
            _mark_old_offline(cp, old_state)
        branch = ""
        try:
            branch = subprocess.run(
                ["git", "-C", str(event.get("cwd") or REPO_ROOT),
                 "rev-parse", "--abbrev-ref", "HEAD"],
                capture_output=True, text=True, timeout=2, check=False).stdout.strip()
        except (OSError, subprocess.SubprocessError):
            pass
        state: Dict[str, Any] = {
            "schemaVersion": 1,
            "client": client_name,
            "clientSessionId": session_id,
            "workerKey": worker_key,
            "host": host,
            "bootId": boot_id,
            "processUuid": process_uuid,
            "hostPid": host_pid,
            "hostProcessStartedAt": host_process_started_at,
            "verifyHostCommand": verify_command,
            "cwd": str(event.get("cwd") or os.getcwd()),
            "transcriptPath": event.get("transcript_path"),
            "branch": branch,
            "source": event.get("source"),
            "version": os.environ.get("JARVIS_INTERACTIVE_WORKER_VERSION", "interactive-v1"),
            "claimCounter": int(old_state.get("claimCounter") or 0),
            "current": old_state.get("current") if same_incarnation else None,
            "pendingClaim": old_state.get("pendingClaim") if same_incarnation else None,
            "pendingOperation": old_state.get("pendingOperation") if same_incarnation else None,
            "stopped": not verify_command,
            "registeredAt": int(time.time()),
        }
        if not verify_command:
            state["stoppedAt"] = int(time.time())
            state["lastError"] = "Codex/Claude host process could not be verified"
        if same_incarnation and old_state.get("daemonPid"):
            state["daemonPid"] = old_state.get("daemonPid")
            state["daemonStartedAt"] = old_state.get("daemonStartedAt")
        _retry_unavailable(
            lambda: _register(cp, state, "ACTIVE" if verify_command else "OFFLINE"))
        store.save(state)
        if verify_command:
            _ensure_daemon(store, worker_key)
            _hook_output(
                event_name,
                "Jarvis interactive worker 已注册（仅定向接单，不拉取公共队列）。"
                "处理 Aone 前必须先运行 bootstrap/claim.sh；控制面异常时接单会 fail-closed。")
        else:
            _hook_output(
                event_name,
                "警告：无法校验 Codex/Claude 宿主进程，Worker 已按 OFFLINE 注册且不会启动 sidecar。"
                "为避免幽灵续租，bootstrap/claim.sh 将 fail-closed。")
    except Exception as exc:
        print("interactive worker registration warning: %s" % type(exc).__name__,
              file=sys.stderr)
        _hook_output(
            event_name,
            "警告：Jarvis interactive worker 注册失败。不要直接修改 Aone；"
            "bootstrap/claim.sh 会 fail-closed。请确认 Codex 项目 hook 已信任以及控制面配置可用。")
    return 0


def offline(store: StateStore, client: Optional[AutomationAgentTaskClient] = None,
            expected_worker_key: Optional[str] = None) -> None:
    cp = client or _client()
    with store.locked():
        state = store.load_unlocked()
        if not state:
            return
        if expected_worker_key and state.get("workerKey") != expected_worker_key:
            return
        state["stopped"] = True
        state["stoppedAt"] = int(time.time())
        store.save_unlocked(state)
    try:
        _heartbeat_worker(cp, state, "OFFLINE")
    except ControlPlaneError:
        # The server reaper marks stale workers OFFLINE after heartbeats cease.
        pass


def daemon(state_path: Path, expected_worker_key: str) -> int:
    store = StateStore(state_path)
    interval = max(2.0, float(os.environ.get(
        "JARVIS_INTERACTIVE_HEARTBEAT_SEC", "20")))
    cp = _client()
    while True:
        state = store.load()
        if not state or state.get("workerKey") != expected_worker_key:
            return 0
        if state.get("stopped") or not _host_alive(state):
            offline(store, cp, expected_worker_key)
            return 0
        try:
            _heartbeat_worker(cp, state, "ACTIVE")
            current = state.get("current")
            if (isinstance(current, Mapping)
                    and current.get("heartbeatEnabled", True)):
                cp.heartbeat_session(
                    str(current["sessionId"]), state["workerKey"],
                    current["fenceToken"],
                    {"leaseSeconds": int(current.get("leaseSeconds") or 120)},
                    request_id="jarvis-interactive-session-heartbeat-%s" %
                    hashlib.sha256((str(current["sessionId"]) +
                                    str(time.time_ns())).encode()).hexdigest()[:24])
        except (StaleFence, ControlPlaneConflict):
            with store.locked():
                latest = store.load_unlocked()
                if (latest.get("workerKey") == expected_worker_key
                        and isinstance(latest.get("current"), Mapping)
                        and latest["current"].get("sessionId") ==
                        (state.get("current") or {}).get("sessionId")):
                    latest["lastError"] = "session ownership lost"
                    latest["current"] = None
                    latest["pendingOperation"] = None
                    store.save_unlocked(latest)
        except ControlPlaneError as exc:
            # Fail closed for mutations; heartbeat transport failures simply retry.
            print("interactive heartbeat warning: %s" % type(exc).__name__, file=sys.stderr)
        time.sleep(interval)


def _claim_cycle(state: Dict[str, Any], aone_id: str, project_id: str) -> Tuple[int, str]:
    pending = state.get("pendingClaim")
    current = state.get("current")
    if (isinstance(pending, Mapping)
            and str(pending.get("aoneId")) == aone_id
            and str(pending.get("projectId")) == project_id):
        return int(pending["cycle"]), str(pending["runtimeSessionId"])
    if (isinstance(current, Mapping)
            and str(current.get("aoneId")) == aone_id
            and str(current.get("projectId")) == project_id):
        return int(current["cycle"]), str(current["runtimeSessionId"])
    cycle = int(state.get("claimCounter") or 0) + 1
    session_hash = hashlib.sha256(
        str(state["clientSessionId"]).encode("utf-8")).hexdigest()[:16]
    runtime_id = "interactive:%s:%s:aone:%s:%s:cycle:%d" % (
        state["client"], session_hash, project_id, aone_id, cycle)
    return cycle, runtime_id[:191]


def _matching_current(state: Mapping[str, Any], aone_id: str) -> Mapping[str, Any]:
    current = state.get("current")
    if not isinstance(current, Mapping) or str(current.get("aoneId")) != str(aone_id):
        raise RuntimeError("current interactive session does not own Aone %s" % aone_id)
    return current


def prepare_claim(aone_id: str, project_id: str) -> Dict[str, Any]:
    aone_id = _nonblank(aone_id, "aone_id")
    project_id = _nonblank(project_id, "project_id")
    store = _current_store()
    cp = _client()
    with store.locked():
        state = store.load_unlocked()
        if not state or state.get("stopped"):
            raise RuntimeError("interactive worker is not active")
        if not state.get("verifyHostCommand"):
            raise RuntimeError("interactive worker host process is not verified")
        current_before_claim = state.get("current")
        if (isinstance(current_before_claim, Mapping)
                and (str(current_before_claim.get("aoneId")) != aone_id
                     or str(current_before_claim.get("projectId")) != project_id)):
            raise ControlPlaneConflict(
                "interactive worker already owns Aone %s" %
                current_before_claim.get("aoneId"))
        cycle, runtime_id = _claim_cycle(state, aone_id, project_id)
        state["claimCounter"] = max(int(state.get("claimCounter") or 0), cycle)
        state["pendingClaim"] = {
            "aoneId": aone_id,
            "projectId": project_id,
            "cycle": cycle,
            "runtimeSessionId": runtime_id,
        }
        store.save_unlocked(state)

    _retry_unavailable(lambda: _register(cp, state))
    revision = "interactive:%s" % hashlib.sha256(runtime_id.encode()).hexdigest()[:32]
    envelope = TaskEnvelope(
        task_key="aone:%s:%s" % (project_id, aone_id),
        source_type="AONE",
        source_ref={"aoneId": aone_id, "projectId": project_id},
        task_type="ticket",
        desired_revision=revision,
        trigger_mask=["INTERACTIVE"],
        payload={
            "kind": "ticket",
            "itemId": aone_id,
            "project": project_id,
            "trigger": "INTERACTIVE",
        },
        recovery_policy="REPLAY_SAFE",
        aone_id=aone_id,
        required_capabilities={"workerKey": state["workerKey"]},
    )
    lease_seconds = max(30, int(os.environ.get(
        "JARVIS_INTERACTIVE_LEASE_SECONDS", "120")))
    # Retries of a locally active task report zero free capacity.  The direct
    # endpoint may resume that exact fenced task, but this one-slot worker must
    # never look available for a second assignment.
    free_slots = 0 if isinstance(current_before_claim, Mapping) else 1
    lease = _retry_unavailable(lambda: cp.claim_task(
        state["workerKey"], envelope, runtime_session_id=runtime_id,
        lease_seconds=lease_seconds, free_slots=free_slots,
        request_id="jarvis-interactive-claim-%s" %
        hashlib.sha256((state["workerKey"] + runtime_id).encode()).hexdigest()[:24]))
    task = lease.get("task") if isinstance(lease, Mapping) else None
    session = lease.get("session") if isinstance(lease, Mapping) else None
    if not isinstance(task, Mapping) or not isinstance(session, Mapping):
        raise RuntimeError("direct claim returned no task/session")
    task_id = _field(task, "id", "taskId", "task_id")
    session_id = _field(session, "id", "sessionId", "session_id")
    generation = _field(session, "generation") or _field(task, "generation")
    fence = _field(session, "fenceToken", "fence_token")
    if task_id is None or session_id is None or generation is None or fence is None:
        raise RuntimeError("direct claim response is incomplete")
    start_detail = {
        "runtimeSessionId": runtime_id,
        "pid": int(state["hostPid"]),
        "transcriptUri": state.get("transcriptPath"),
        "workspaceRef": state.get("cwd"),
        "branchRef": state.get("branch"),
        "logUri": str(store.path.with_suffix(".log")),
        "leaseSeconds": lease_seconds,
    }
    _retry_unavailable(lambda: cp.start_session(
        str(session_id), state["workerKey"], fence, start_detail,
        request_id="jarvis-interactive-session-start-%s" %
        hashlib.sha256(runtime_id.encode()).hexdigest()[:24]))

    current = {
        "aoneId": aone_id,
        "projectId": project_id,
        "taskId": task_id,
        "sessionId": session_id,
        "generation": generation,
        "fenceToken": fence,
        "attemptNo": _field(session, "attemptNo", "attempt_no") or 1,
        "cycle": cycle,
        "runtimeSessionId": runtime_id,
        "leaseSeconds": lease_seconds,
        "startedAt": (_field(current_before_claim, "startedAt")
                      if isinstance(current_before_claim, Mapping)
                      else int(time.time())),
        # The sidecar may only renew after the external-operation receipt has
        # been durably recorded.  A lost begin response must age out safely.
        "heartbeatEnabled": bool(
            isinstance(current_before_claim, Mapping)
            and current_before_claim.get("heartbeatEnabled")),
    }
    with store.locked():
        latest = store.load_unlocked()
        latest["current"] = current
        latest["pendingClaim"] = None
        store.save_unlocked(latest)

    request_payload = {
        "aoneId": aone_id,
        "projectId": project_id,
        "addTag": "jarvis-claimed",
        "removeTag": "jarvis-idle",
    }
    retry_claim = state.get("pendingClaim")
    saved_operation_key = (
        str(retry_claim.get("operationKey"))
        if isinstance(retry_claim, Mapping) and retry_claim.get("operationKey")
        else ""
    )
    operation_key = (saved_operation_key or
                     "aone-claim:%s:%s:%s" % (task_id, generation, cycle))
    operation_request = {
        "taskId": task_id,
        "sessionId": session_id,
        "generation": generation,
        "workerKey": state["workerKey"],
        "fenceToken": fence,
        "operationKey": operation_key,
        "operationType": "AONE_CLAIM",
        "target": aone_id,
        "requestPayload": request_payload,
        "required": True,
        "maxRetries": 3,
    }
    # Persist intent before sending begin.  If the HTTP response is lost, this
    # proves that the caller could not yet have returned to claim.sh to perform
    # the Aone write, so the same idempotent tag effect may be resumed safely.
    with store.locked():
        latest = store.load_unlocked()
        existing_pending = latest.get("pendingOperation")
        had_local_intent = isinstance(existing_pending, Mapping)
        if (isinstance(existing_pending, Mapping)
                and str(existing_pending.get("operationKey")) != operation_key):
            raise ControlPlaneConflict(
                "another external operation is pending for this worker")
        if not isinstance(existing_pending, Mapping):
            existing_pending = {
                "operationId": None,
                "operationKey": operation_key,
                "aoneId": aone_id,
                "proceed": False,
                "status": "BEGINNING",
            }
            latest["pendingOperation"] = existing_pending
            store.save_unlocked(latest)
    begun = cp.begin_operation(
        operation_request,
        request_id="jarvis-interactive-operation-begin-%s" %
        hashlib.sha256(operation_key.encode()).hexdigest()[:24])
    operation = begun.get("operation") if isinstance(begun, Mapping) else None
    if not isinstance(operation, Mapping):
        raise RuntimeError("operation begin returned no receipt")
    operation_id = _field(operation, "id", "operationId", "operation_id")
    if operation_id is None:
        raise RuntimeError("operation receipt has no id")
    operation_status = str(
        _field(operation, "status", "operationStatus") or "").upper()
    server_proceed = bool(begun.get("proceed", True))
    local_retry = (
        had_local_intent
        and isinstance(existing_pending, Mapping)
        and str(existing_pending.get("operationKey")) == operation_key
        and (existing_pending.get("operationId") is None
             or str(existing_pending.get("operationId")) == str(operation_id))
    )
    with store.locked():
        latest = store.load_unlocked()
        pending = latest.get("pendingOperation")
        if (isinstance(pending, Mapping)
                and str(pending.get("operationKey")) == operation_key):
            latest["pendingOperation"] = {
                "operationId": operation_id,
                "operationKey": operation_key,
                "aoneId": aone_id,
                "proceed": False,
                "status": operation_status or "SENDING",
            }
            store.save_unlocked(latest)
    if operation_status == "ACKED":
        proceed = False
    elif server_proceed:
        proceed = True
    elif operation_status == "SENDING" and local_retry:
        # The prior process persisted ownership before dying, so replaying the
        # idempotent Aone tag replacement and then ACKing is safe.
        proceed = True
    else:
        # A SENDING receipt without a local record can be a lost begin response.
        # We cannot prove whether an external effect started, so fail closed.
        raise ControlPlaneConflict(
            "AONE_CLAIM receipt is %s without resumable local ownership" %
            (operation_status or "UNKNOWN"))

    with store.locked():
        latest = store.load_unlocked()
        latest["pendingOperation"] = ({
            "operationId": operation_id,
            "operationKey": operation_key,
            "aoneId": aone_id,
            "proceed": proceed,
            "status": operation_status or "SENDING",
        } if proceed else None)
        if isinstance(latest.get("current"), Mapping):
            latest["current"]["heartbeatEnabled"] = True
        store.save_unlocked(latest)
    return {
        "accepted": True,
        "proceed": proceed,
        "workerKey": state["workerKey"],
        "taskId": task_id,
        "sessionId": session_id,
        "generation": generation,
        "fenceToken": fence,
        "operationId": operation_id,
        "operationStatus": operation_status,
        "runtimeSessionId": runtime_id,
    }


def acknowledge_claim(aone_id: str, external_ref: str) -> Dict[str, Any]:
    store = _current_store()
    cp = _client()
    state = store.load()
    current = _matching_current(state, aone_id)
    pending = state.get("pendingOperation")
    if not isinstance(pending, Mapping):
        raise RuntimeError("no pending AONE_CLAIM operation")
    request = {
        "operationId": pending["operationId"],
        "workerKey": state["workerKey"],
        "fenceToken": current["fenceToken"],
        "externalRef": _nonblank(external_ref, "external_ref"),
    }
    try:
        result = _retry_unavailable(lambda: cp.ack_operation(
            request, request_id="jarvis-interactive-operation-ack-%s" %
            hashlib.sha256(str(pending["operationId"]).encode()).hexdigest()[:24]))
    except Exception as exc:
        # If ACK delivery itself is ambiguous, try to freeze the receipt as UNKNOWN.
        try:
            cp.fail_operation({
                "operationId": pending["operationId"],
                "workerKey": state["workerKey"],
                "fenceToken": current["fenceToken"],
                "error": {"errorType": type(exc).__name__,
                          "message": "Aone claim write succeeded but ACK failed"},
                "unknown": True,
                "retryAfterSeconds": 0,
            }, request_id="jarvis-interactive-operation-ack-unknown-%s" %
            hashlib.sha256(str(pending["operationId"]).encode()).hexdigest()[:24])
        except Exception:
            pass
        with store.locked():
            latest = store.load_unlocked()
            if isinstance(latest.get("current"), Mapping):
                latest["current"]["heartbeatEnabled"] = False
                latest["lastError"] = "AONE_CLAIM acknowledgement is ambiguous"
                store.save_unlocked(latest)
        raise
    with store.locked():
        latest = store.load_unlocked()
        latest["pendingOperation"] = None
        store.save_unlocked(latest)
    return dict(result)


def fail_claim(aone_id: str, message: str, *, unknown: bool = False) -> None:
    store = _current_store()
    cp = _client()
    state = store.load()
    current = _matching_current(state, aone_id)
    pending = state.get("pendingOperation")
    error = {"errorType": "AoneClaimFailed", "message": str(message)[:500]}
    first_error: Optional[BaseException] = None
    if isinstance(pending, Mapping):
        try:
            cp.fail_operation({
                "operationId": pending["operationId"],
                "workerKey": state["workerKey"],
                "fenceToken": current["fenceToken"],
                "error": error,
                "unknown": bool(unknown),
                "retryAfterSeconds": 0,
            }, request_id="jarvis-interactive-operation-fail-%s" %
            hashlib.sha256(str(pending["operationId"]).encode()).hexdigest()[:24])
        except BaseException as exc:
            first_error = exc
    try:
        cp.fail_session(
            str(current["sessionId"]), state["workerKey"], current["fenceToken"],
            {"error": error, "retryAfterSeconds": 0},
            request_id="jarvis-interactive-session-fail-%s" %
            hashlib.sha256(str(current["runtimeSessionId"]).encode()).hexdigest()[:24])
    except BaseException as exc:
        if first_error is None:
            first_error = exc
    if first_error is None:
        with store.locked():
            latest = store.load_unlocked()
            latest["current"] = None
            # A failed required receipt remains part of this exact logical
            # claim.  Reuse its cycle/runtime/operation key after RETRY_WAIT;
            # UNKNOWN likewise stays fenced until explicit reconciliation.
            latest["pendingClaim"] = {
                "aoneId": str(current["aoneId"]),
                "projectId": str(current["projectId"]),
                "cycle": int(current["cycle"]),
                "runtimeSessionId": str(current["runtimeSessionId"]),
                "operationKey": (str(pending.get("operationKey"))
                                 if isinstance(pending, Mapping) else ""),
                "receiptUnknown": bool(unknown),
            }
            latest["pendingOperation"] = None
            store.save_unlocked(latest)
    if first_error is not None:
        raise first_error


def has_current(aone_id: str) -> bool:
    try:
        state = _current_store().load()
    except RuntimeError:
        return False
    current = state.get("current")
    return isinstance(current, Mapping) and str(current.get("aoneId")) == str(aone_id)


def transition(aone_id: str, action: str, detail: Optional[str] = None) -> Dict[str, Any]:
    store = _current_store()
    cp = _client()
    state = store.load()
    current = _matching_current(state, aone_id)
    if state.get("pendingOperation"):
        raise RuntimeError("cannot close session while AONE_CLAIM receipt is pending")
    if action == "suspend":
        payload: Dict[str, Any] = {
            "waitType": "MANUAL",
            "waitKey": "aone:%s" % aone_id,
            "transcriptUri": state.get("transcriptPath"),
            "branchRef": state.get("branch"),
            "logUri": str(store.path.with_suffix(".log")),
        }
        if detail:
            payload["waitCursor"] = hashlib.sha256(detail.encode()).hexdigest()[:32]
        result = _retry_unavailable(lambda: cp.suspend_session(
            str(current["sessionId"]), state["workerKey"], current["fenceToken"], payload,
            request_id="jarvis-interactive-session-suspend-%s" %
            hashlib.sha256(str(current["runtimeSessionId"]).encode()).hexdigest()[:24]))
    elif action == "complete":
        payload = {"result": {"aoneId": str(aone_id), "summary": str(detail or "completed")}}
        result = _retry_unavailable(lambda: cp.complete_session(
            str(current["sessionId"]), state["workerKey"], current["fenceToken"], payload,
            request_id="jarvis-interactive-session-complete-%s" %
            hashlib.sha256(str(current["runtimeSessionId"]).encode()).hexdigest()[:24]))
    else:
        raise ValueError("unknown transition %s" % action)
    with store.locked():
        latest = store.load_unlocked()
        latest["current"] = None
        latest["pendingClaim"] = None
        latest["pendingOperation"] = None
        store.save_unlocked(latest)
    return dict(result)


def worker_status() -> Any:
    state = _current_store().load()
    return _client().get_worker_state(state["workerKey"])


def _print_json(value: Any) -> None:
    print(json.dumps(value, ensure_ascii=False, sort_keys=True,
                     separators=(",", ":"), default=str))


def _parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="command", required=True)
    hook_parser = sub.add_parser("hook")
    hook_parser.add_argument("client", choices=("claude", "codex"))
    daemon_parser = sub.add_parser("daemon")
    daemon_parser.add_argument("--state", required=True)
    daemon_parser.add_argument("--worker-key", required=True)
    claim_parser = sub.add_parser("prepare-claim")
    claim_parser.add_argument("aone_id")
    claim_parser.add_argument("project_id")
    ack_parser = sub.add_parser("operation-ack")
    ack_parser.add_argument("aone_id")
    ack_parser.add_argument("external_ref")
    fail_parser = sub.add_parser("operation-fail")
    fail_parser.add_argument("aone_id")
    fail_parser.add_argument("message")
    fail_parser.add_argument("--unknown", action="store_true")
    current_parser = sub.add_parser("has-current")
    current_parser.add_argument("aone_id")
    suspend_parser = sub.add_parser("suspend")
    suspend_parser.add_argument("aone_id")
    suspend_parser.add_argument("detail", nargs="?", default="released")
    complete_parser = sub.add_parser("complete")
    complete_parser.add_argument("aone_id")
    complete_parser.add_argument("detail", nargs="?", default="completed")
    sub.add_parser("status")
    return parser


def main(argv: Optional[list[str]] = None) -> int:
    args = _parser().parse_args(argv)
    if args.command == "hook":
        try:
            event = json.load(sys.stdin)
            if not isinstance(event, Mapping):
                raise ValueError("hook input must be an object")
        except Exception as exc:
            print("interactive hook input error: %s" % type(exc).__name__, file=sys.stderr)
            _hook_output("SessionStart", "Jarvis worker hook input 无效；接单将 fail-closed。")
            return 0
        return hook(args.client, event)
    try:
        if args.command == "daemon":
            return daemon(Path(args.state), args.worker_key)
        if args.command == "prepare-claim":
            _print_json(prepare_claim(args.aone_id, args.project_id))
        elif args.command == "operation-ack":
            _print_json(acknowledge_claim(args.aone_id, args.external_ref))
        elif args.command == "operation-fail":
            fail_claim(args.aone_id, args.message, unknown=args.unknown)
            _print_json({"failed": True})
        elif args.command == "has-current":
            return 0 if has_current(args.aone_id) else 1
        elif args.command == "suspend":
            _print_json(transition(args.aone_id, "suspend", args.detail))
        elif args.command == "complete":
            _print_json(transition(args.aone_id, "complete", args.detail))
        elif args.command == "status":
            _print_json(worker_status())
        return 0
    except ControlPlaneConflict as exc:
        print("interactive claim conflict: %s" % exc, file=sys.stderr)
        return CONFLICT_EXIT
    except (ControlPlaneUnavailable, StaleFence) as exc:
        print("interactive control plane unavailable/fenced: %s" % exc, file=sys.stderr)
        return UNAVAILABLE_EXIT
    except ControlPlaneError as exc:
        print("interactive control plane rejected request: %s" % exc, file=sys.stderr)
        return TRANSITION_EXIT
    except (RuntimeError, ValueError, KeyError) as exc:
        print("interactive worker state error: %s" % exc, file=sys.stderr)
        return STATE_ERROR_EXIT


if __name__ == "__main__":
    raise SystemExit(main())
