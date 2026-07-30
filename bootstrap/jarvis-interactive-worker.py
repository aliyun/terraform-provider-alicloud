#!/usr/bin/env python3
"""Interactive Codex/Claude worker registration and direct-task lifecycle.

The ordinary PersistenceExecutor pulls matching Tasks.  An interactive worker registers
itself with a unique worker capability, then attaches only to the exact Aone Task
requested by the user through ``POST /tasks/claim``.  A small detached
sidecar keeps both the worker and its current fenced session alive while the real
Codex/Claude host process exists.

State is local coordination data, not the Aone ownership authority.  It is written
atomically with mode 0600 and deliberately excludes credentials; the database fence
remains the ownership authority.  A headless execution policy only removes
capabilities.  That policy is also mirrored to worker capabilities, and the live
process lineage retains the exec-headless arguments so deleting a cache file cannot
turn a restricted worker into an Aone writer.
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
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable, Dict, Iterator, Mapping, Optional, Tuple, TypeVar


HERE = Path(__file__).resolve().parent
REPO_ROOT = HERE.parent
BRIDGE_DIR = REPO_ROOT / "bridge"
# Package-style imports used by the shared executor need the worktree root.
# ``python -I`` intentionally ignores cwd/PYTHONPATH, so make this explicit
# before importing the legacy top-level module names below.
sys.path.insert(0, str(REPO_ROOT))
sys.path.insert(0, str(HERE))
sys.path.insert(0, str(BRIDGE_DIR))

from a1_command_guard import pretool_a1_block_reason  # noqa: E402
from jarvis_persistence_executor import _default_boot_id, make_worker_key  # noqa: E402
from jarvis_task_client import (  # noqa: E402
    ControlPlaneClient,
    ControlPlaneConflict,
    ControlPlaneError,
    ControlPlaneRejected,
    ControlPlaneUnavailable,
    StaleFence,
    TaskEnvelope,
)
from bridge.task_policy import (  # noqa: E402
    HEADLESS_POLICY_REVISION,
    policy_desired_revision,
)


CONFLICT_EXIT = 10
UNAVAILABLE_EXIT = 11
STATE_ERROR_EXIT = 12
TRANSITION_EXIT = 13
HOOK_BLOCK_EXIT = 2
POST_PR_AONE_WRITE_POLICY = "post-pr-read-only"
POST_PR_HEADLESS_KINDS = frozenset(("pr_ci_fix", "pr_comment_reply"))
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


def _client(*, timeout_override: Optional[float] = None) -> ControlPlaneClient:
    base_url = (os.environ.get("JARVIS_CONTROL_PLANE_BASE_URL", "").strip()
                or os.environ.get("JARVIS_HTML_REPORT_BASE_URL", "").strip()
                or "https://pre-agent.aliyun-inc.com")
    token = (os.environ.get("JARVIS_CONTROL_PLANE_TOKEN", "").strip()
             or os.environ.get("JARVIS_HTML_REPORT_TOKEN", "").strip())
    timeout = float(os.environ.get("JARVIS_CONTROL_PLANE_TIMEOUT", "10"))
    if timeout_override is not None:
        timeout = min(timeout, float(timeout_override))
    return ControlPlaneClient(base_url, token, timeout=timeout)


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
            ["/bin/ps", "-o", "ppid=", "-o", "command=", "-p", str(pid)],
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


def _process_command(pid: int) -> str:
    """Return the untruncated command line used for immutable lineage checks."""
    try:
        result = subprocess.run(
            ["/bin/ps", "-ww", "-o", "command=", "-p", str(pid)],
            capture_output=True, text=True, timeout=2, check=False)
    except (OSError, subprocess.SubprocessError):
        return ""
    return result.stdout.strip()


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


def _nearest_runtime_client() -> str:
    """Resolve nested Claude/Codex launches from the nearest real host process."""
    pid = os.getppid()
    for _depth in range(16):
        if pid <= 1:
            break
        parent, command = _process_info(pid)
        lower = command.lower()
        helper = ("jarvis-interactive-worker" in lower
                  or "run-interactive-worker-hook" in lower)
        if not helper:
            if re.search(r"(^|[/ ])claude([ /]|$)", lower):
                return "claude"
            if "codex" in lower:
                return "codex"
        pid = parent
    return ""


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
    claude = os.environ.get("CLAUDE_CODE_SESSION_ID", "").strip()
    persisted_client = os.environ.get("JARVIS_INTERACTIVE_CLIENT", "").strip().lower()
    persisted_session = os.environ.get("JARVIS_INTERACTIVE_SESSION_ID", "").strip()

    # A developer can start `idea` (Claude) from a Codex terminal, or vice
    # versa. Both native environment variables may then be present. Select the
    # nearest real host instead of letting the outer client steal the inner
    # client's database-fenced assignment.
    nearest = _nearest_runtime_client()
    if nearest == "claude":
        if claude:
            return "claude", claude
        if persisted_client == "claude" and persisted_session:
            return "claude", persisted_session
    if nearest == "codex" and codex:
        return "codex", codex

    if claude and not codex:
        return "claude", claude
    if codex and not claude:
        return "codex", codex
    if claude:
        return "claude", claude
    if codex:
        return "codex", codex
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
    capabilities = {
        "dispatch": {"pull": False},
        "bridgeRole": os.environ.get("JARVIS_BRIDGE_ROLE", "interactive"),
        "workerMode": "INTERACTIVE",
        "workerKey": state.get("workerKey"),
        "client": state.get("client"),
        "clientSessionId": state.get("clientSessionId"),
        "hostPid": state.get("hostPid"),
        "hostProcessStartedAt": state.get("hostProcessStartedAt"),
        "kinds": ["interactive-aone"],
    }
    policy = state.get("headlessPolicy")
    if state.get("headlessRegistered") and isinstance(policy, Mapping):
        capabilities["headlessPolicy"] = dict(policy)
    return capabilities


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


def _heartbeat_worker(client: ControlPlaneClient,
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
        process_uuid=state["processUuid"],
        request_id="jarvis-interactive-worker-heartbeat-%s" %
        hashlib.sha256((state["workerKey"] + str(time.time_ns())).encode()).hexdigest()[:24])


def _register(client: ControlPlaneClient, state: Mapping[str, Any],
              status: str = "ACTIVE") -> None:
    client.register_worker(
        _worker_payload(state, status),
        request_id="jarvis-interactive-worker-register-%s" %
        hashlib.sha256(state["workerKey"].encode()).hexdigest()[:24])


def _mark_old_offline(client: ControlPlaneClient,
                      old_state: Mapping[str, Any]) -> None:
    if not all(old_state.get(key) for key in (
            "workerKey", "host", "bootId", "processUuid")):
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
        daemon_started = _process_start_identity(int(state["daemonPid"]))
        if daemon_started:
            state["daemonProcessStartedAt"] = daemon_started
        state.pop("sidecarHeartbeatAt", None)
        store.save_unlocked(state)


def _hook_output(event_name: str, message: Optional[str] = None) -> None:
    output: Dict[str, Any] = {}
    if event_name in ("SessionStart", "SubagentStart") and message:
        output["hookSpecificOutput"] = {
            "hookEventName": event_name,
            "additionalContext": message,
        }
    elif message:
        output["systemMessage"] = message
    print(json.dumps(output, ensure_ascii=False, separators=(",", ":")))


def _assignment_epoch(state: Mapping[str, Any]) -> str:
    """Identify the exact local task attachment a subagent may inherit."""
    current = state.get("current")
    if isinstance(current, Mapping):
        return "session:%s:%s:%s" % (
            str(current.get("sessionId") or "missing"),
            str(current.get("fenceToken") or "missing"),
            str(current.get("generation") or "missing"),
        )
    return "idle:%s" % int(state.get("claimCounter") or 0)


def _session_meta(transcript_path: Any) -> Dict[str, Any]:
    path = Path(_nonblank(transcript_path, "transcript_path")).expanduser()
    try:
        with path.open("r", encoding="utf-8") as handle:
            raw = handle.readline(1024 * 1024)
    except OSError as exc:
        raise RuntimeError("subagent transcript is unavailable") from exc
    try:
        value = json.loads(raw)
    except (TypeError, ValueError) as exc:
        raise RuntimeError("subagent transcript has invalid session metadata") from exc
    if not isinstance(value, Mapping) or value.get("type") != "session_meta":
        raise RuntimeError("subagent transcript is missing session metadata")
    payload = value.get("payload")
    if not isinstance(payload, Mapping):
        raise RuntimeError("subagent transcript session metadata is invalid")
    return dict(payload)


def _spawn_origin(parent_transcript: Any,
                  child_agent_id: str) -> Tuple[str, str]:
    """Find the exact trusted parent spawn tool call for one Codex child."""
    path = Path(_nonblank(parent_transcript, "parent transcript")).expanduser()
    spawn_call_id = ""
    calls: Dict[str, str] = {}
    try:
        with path.open("r", encoding="utf-8") as handle:
            for raw in handle:
                try:
                    item = json.loads(raw)
                except (TypeError, ValueError):
                    continue
                if not isinstance(item, Mapping):
                    continue
                payload = item.get("payload")
                if not isinstance(payload, Mapping):
                    continue
                if (item.get("type") == "event_msg"
                        and payload.get("type") == "sub_agent_activity"
                        and payload.get("kind") == "started"
                        and str(payload.get("agent_thread_id") or "") == child_agent_id):
                    spawn_call_id = str(payload.get("event_id") or "").strip()
                if (item.get("type") == "response_item"
                        and payload.get("type") in ("function_call", "custom_tool_call")
                        and payload.get("name") == "spawn_agent"):
                    call_id = str(payload.get("call_id") or "").strip()
                    metadata = payload.get("internal_chat_message_metadata_passthrough")
                    turn_id = (str(metadata.get("turn_id") or "").strip()
                               if isinstance(metadata, Mapping) else "")
                    if call_id and turn_id:
                        calls[call_id] = turn_id
    except OSError as exc:
        raise RuntimeError("parent transcript is unavailable") from exc
    parent_turn_id = calls.get(spawn_call_id, "")
    if not spawn_call_id or not parent_turn_id:
        raise RuntimeError("parent spawn receipt is not yet durable")
    return spawn_call_id, parent_turn_id


def _subagent_binding(state: Mapping[str, Any]) -> Optional[Dict[str, Any]]:
    binding = state.get("subagentBinding")
    return dict(binding) if isinstance(binding, Mapping) else None


_CODEX_SPAWN_TOOL_NAMES = {
    "spawn_agent",
    "collaborationspawn_agent",
    "collaboration.spawn_agent",
}
_CODEX_INTERACTION_TOOL_NAMES = {
    "followup_task": "followup_task",
    "send_message": "send_message",
    "collaborationfollowup_task": "followup_task",
    "collaborationsend_message": "send_message",
    "collaboration.followup_task": "followup_task",
    "collaboration.send_message": "send_message",
}


def _codex_tool_kind(event: Mapping[str, Any]) -> str:
    tool_name = str(event.get("tool_name") or "")
    if tool_name in _CODEX_SPAWN_TOOL_NAMES:
        return "spawn_agent"
    return _CODEX_INTERACTION_TOOL_NAMES.get(tool_name, "")


def _next_subagent_revision(state: Dict[str, Any]) -> int:
    revision = int(state.get("subagentRevision") or 0) + 1
    state["subagentRevision"] = revision
    return revision


def _registry_entry(value: Any) -> Optional[Dict[str, Any]]:
    if isinstance(value, Mapping):
        agent_id = str(value.get("agentId") or "").strip()
        if agent_id:
            entry = dict(value)
            entry["agentId"] = agent_id
            entry["bindingRevision"] = int(
                value.get("bindingRevision") or value.get("boundAt") or 0)
            return entry
    elif value:
        # Read old local state fail-safely during rollout. Any new versioned
        # binding supersedes this revision-zero entry.
        return {"agentId": str(value), "bindingRevision": 0}
    return None


def _resolve_subagent_target(state: Mapping[str, Any],
                             target: str) -> Tuple[str, int]:
    registry = state.get("subagentRegistry")
    registry = registry if isinstance(registry, Mapping) else {}
    direct = _registry_entry(registry.get(target))
    if direct:
        return direct["agentId"], int(direct.get("bindingRevision") or 0)
    known_ids = set()
    for value in registry.values():
        entry = _registry_entry(value)
        if entry:
            known_ids.add(entry["agentId"])
    for key in ("subagentLifecycles", "subagentEpochs"):
        values = state.get(key)
        if isinstance(values, Mapping):
            known_ids.update(str(agent_id) for agent_id in values)
    return (target, 0) if target in known_ids else ("", 0)


def _record_spawn_permit_locked(state: Dict[str, Any],
                                event: Mapping[str, Any]) -> None:
    if _codex_tool_kind(event) != "spawn_agent":
        return
    tool_use_id = str(event.get("tool_use_id") or "").strip()
    root_turn_id = str(state.get("activeTurnId") or "").strip()
    if not tool_use_id or not root_turn_id:
        return
    permits = dict(state.get("subagentSpawnPermits") or {})
    if isinstance(permits.get(tool_use_id), Mapping):
        return
    tool_input = event.get("tool_input")
    revision = _next_subagent_revision(state)
    permits[tool_use_id] = {
        "parentSessionId": str(event.get("session_id") or ""),
        "parentAgentId": str(event.get("agent_id") or ""),
        "parentThreadId": str(
            event.get("agent_id") or event.get("session_id") or ""),
        "parentTurnId": str(event.get("turn_id") or ""),
        "rootTurnId": root_turn_id,
        "assignmentEpoch": _assignment_epoch(state),
        "taskName": (str(tool_input.get("task_name") or "")
                     if isinstance(tool_input, Mapping) else ""),
        "revision": revision,
        "createdAt": time.time_ns(),
    }
    if len(permits) > 128:
        ordered = sorted(
            permits.items(), key=lambda item: int(item[1].get("createdAt") or 0),
            reverse=True)[:128]
        permits = dict(ordered)
    state["subagentSpawnPermits"] = permits


def _record_subagent_interaction_permit_locked(
        state: Dict[str, Any], event: Mapping[str, Any]) -> None:
    tool_kind = _codex_tool_kind(event)
    if tool_kind not in ("followup_task", "send_message"):
        return
    tool_use_id = str(event.get("tool_use_id") or "").strip()
    tool_input = event.get("tool_input")
    if not tool_use_id or not isinstance(tool_input, Mapping):
        return
    target = str(tool_input.get("target") or "").strip()
    if not target:
        return
    agent_id, binding_revision = _resolve_subagent_target(state, target)
    if not agent_id:
        return
    root_turn_id = str(state.get("activeTurnId") or "").strip()
    if not state.get("turnActive", True) or not root_turn_id:
        return
    permits = dict(state.get("subagentInteractionPermits") or {})
    if isinstance(permits.get(tool_use_id), Mapping):
        return
    revision = _next_subagent_revision(state)
    permits[tool_use_id] = {
        "toolKind": tool_kind,
        "target": target,
        "targetAgentId": agent_id,
        "targetBindingRevision": binding_revision,
        "parentSessionId": str(event.get("session_id") or ""),
        "parentAgentId": str(event.get("agent_id") or ""),
        "parentThreadId": str(
            event.get("agent_id") or event.get("session_id") or ""),
        "parentTurnId": str(event.get("turn_id") or ""),
        "rootTurnId": root_turn_id,
        "assignmentEpoch": _assignment_epoch(state),
        "revision": revision,
        "createdAt": time.time_ns(),
    }
    if len(permits) > 128:
        ordered = sorted(
            permits.items(), key=lambda item: int(item[1].get("revision") or 0),
            reverse=True)[:128]
        permits = dict(ordered)
    state["subagentInteractionPermits"] = permits


def _consume_subagent_interaction_locked(
        state: Dict[str, Any], event: Mapping[str, Any]) -> Tuple[str, int]:
    tool_kind = _codex_tool_kind(event)
    if tool_kind not in ("followup_task", "send_message"):
        return "", 0
    tool_use_id = str(event.get("tool_use_id") or "").strip()
    permits = dict(state.get("subagentInteractionPermits") or {})
    permit = permits.pop(tool_use_id, None)
    state["subagentInteractionPermits"] = permits
    if not isinstance(permit, Mapping):
        return "", 0
    tool_input = event.get("tool_input")
    target = (str(tool_input.get("target") or "").strip()
              if isinstance(tool_input, Mapping) else "")
    agent_id, binding_revision = _resolve_subagent_target(state, target)
    revision = int(permit.get("revision") or 0)
    if (not revision
            or str(permit.get("toolKind") or "") != tool_kind
            or str(permit.get("target") or "") != target
            or str(permit.get("targetAgentId") or "") != agent_id
            or int(permit.get("targetBindingRevision") or 0) != binding_revision
            or str(permit.get("parentSessionId") or "") !=
               str(event.get("session_id") or "")
            or str(permit.get("parentAgentId") or "") !=
               str(event.get("agent_id") or "")
            or str(permit.get("parentThreadId") or "") !=
               str(event.get("agent_id") or event.get("session_id") or "")
            or str(permit.get("parentTurnId") or "") !=
               str(event.get("turn_id") or "")
            or str(permit.get("rootTurnId") or "") !=
               str(state.get("activeTurnId") or "")
            or str(permit.get("assignmentEpoch") or "") !=
               _assignment_epoch(state)):
        return "", 0
    lifecycles = dict(state.get("subagentLifecycles") or {})
    existing = lifecycles.get(agent_id)
    existing_revision = (int(existing.get("revision") or 0)
                         if isinstance(existing, Mapping) else 0)
    if revision <= existing_revision:
        return "", 0
    lifecycles[agent_id] = {
        "status": "ACTIVE",
        "revision": revision,
        "sourceToolUseId": tool_use_id,
    }
    epochs = dict(state.get("subagentEpochs") or {})
    epochs[agent_id] = {
        "rootTurnId": str(permit.get("rootTurnId") or ""),
        "assignmentEpoch": str(permit.get("assignmentEpoch") or ""),
        "authorizedRevision": revision,
        "sourceToolUseId": tool_use_id,
    }
    state["subagentLifecycles"] = lifecycles
    state["subagentEpochs"] = epochs
    return agent_id, revision


def _bind_codex_subagent(store: StateStore,
                         event: Mapping[str, Any]) -> Optional[str]:
    """Bind one Codex child to the exact root-turn/task epoch that spawned it."""
    root_session_id = _nonblank(event.get("session_id"), "session_id")
    child_agent_id = _nonblank(event.get("agent_id"), "agent_id")
    try:
        meta = _session_meta(event.get("transcript_path"))
        if str(meta.get("id") or "") != child_agent_id:
            raise RuntimeError("subagent transcript identity mismatch")
        if str(meta.get("session_id") or "") != root_session_id:
            raise RuntimeError("subagent transcript root session mismatch")
        source = meta.get("source")
        subagent = source.get("subagent") if isinstance(source, Mapping) else None
        spawn = subagent.get("thread_spawn") if isinstance(subagent, Mapping) else None
        if not isinstance(spawn, Mapping):
            raise RuntimeError("subagent transcript has no parent spawn edge")
        parent_thread_id = _nonblank(
            spawn.get("parent_thread_id"), "parent_thread_id")
        agent_path = _nonblank(spawn.get("agent_path"), "agent_path")
        root_store = StateStore(_state_path("codex", root_session_id))
        if parent_thread_id == root_session_id:
            parent_state = root_store.load()
            parent_agent_id = ""
        else:
            parent_store = StateStore(_state_path("codex", parent_thread_id))
            parent_state = parent_store.load()
            parent_binding = _subagent_binding(parent_state)
            if (parent_binding is None
                    or str(parent_binding.get("agentId") or "") != parent_thread_id
                    or str(parent_binding.get("rootSessionId") or "") != root_session_id):
                raise RuntimeError("parent subagent has no matching root Worker binding")
            parent_agent_id = parent_thread_id
        if not parent_state:
            raise RuntimeError("parent interactive Worker state is unavailable")
        spawn_call_id, parent_turn_id = _spawn_origin(
            parent_state.get("transcriptPath"), child_agent_id)
        with root_store.locked():
            root_state = root_store.load_unlocked()
            if (not root_state
                    or str(root_state.get("clientSessionId") or "") != root_session_id):
                raise RuntimeError("root interactive Worker state is unavailable")
            permits = root_state.get("subagentSpawnPermits")
            permit = (permits.get(spawn_call_id)
                      if isinstance(permits, Mapping) else None)
            if not isinstance(permit, Mapping):
                raise RuntimeError("spawn tool was not authorized by the Worker fence")
            if (str(permit.get("parentSessionId") or "") != root_session_id
                    or str(permit.get("parentAgentId") or "") != parent_agent_id
                    or str(permit.get("parentThreadId") or "") != parent_thread_id
                    or str(permit.get("parentTurnId") or "") != parent_turn_id):
                raise RuntimeError("spawn receipt does not match the parent turn")
            registry = dict(root_state.get("subagentRegistry") or {})
            existing_entry = _registry_entry(registry.get(agent_path))
            existing_revision = (int(existing_entry.get("bindingRevision") or 0)
                                 if existing_entry else 0)
            revision = int(permit.get("revision") or 0)
            if not revision:
                raise RuntimeError("spawn receipt has no authorization revision")
            if revision >= existing_revision:
                old_agent_id = (str(existing_entry.get("agentId") or "")
                                if existing_entry else "")
                registry[agent_path] = {
                    "agentId": child_agent_id,
                    "bindingRevision": revision,
                    "sourceToolUseId": spawn_call_id,
                }
                if old_agent_id and old_agent_id != child_agent_id:
                    lifecycles = dict(root_state.get("subagentLifecycles") or {})
                    old_lifecycle = lifecycles.get(old_agent_id)
                    old_revision = (int(old_lifecycle.get("revision") or 0)
                                    if isinstance(old_lifecycle, Mapping) else 0)
                    if revision >= old_revision:
                        lifecycles[old_agent_id] = {
                            "status": "STOPPED",
                            "revision": revision,
                            "reason": "agent_path_rebound",
                            "sourceToolUseId": spawn_call_id,
                        }
                        epochs = dict(root_state.get("subagentEpochs") or {})
                        epochs.pop(old_agent_id, None)
                        root_state["subagentEpochs"] = epochs
                        root_state["subagentLifecycles"] = lifecycles
            root_state["subagentRegistry"] = registry
            current_entry = _registry_entry(registry.get(agent_path))
            registry_matches = bool(
                current_entry
                and current_entry.get("agentId") == child_agent_id
                and int(current_entry.get("bindingRevision") or 0) == revision)
            permit_is_current = bool(
                root_state.get("turnActive", True)
                and str(permit.get("rootTurnId") or "") ==
                    str(root_state.get("activeTurnId") or "")
                and str(permit.get("assignmentEpoch") or "") ==
                    _assignment_epoch(root_state))
            lifecycles = dict(root_state.get("subagentLifecycles") or {})
            lifecycle = lifecycles.get(child_agent_id)
            lifecycle_revision = (int(lifecycle.get("revision") or 0)
                                  if isinstance(lifecycle, Mapping) else 0)
            if registry_matches and permit_is_current and revision >= lifecycle_revision:
                lifecycles[child_agent_id] = {
                    "status": "ACTIVE",
                    "revision": revision,
                    "sourceToolUseId": spawn_call_id,
                }
                epochs = dict(root_state.get("subagentEpochs") or {})
                epochs[child_agent_id] = {
                    "rootTurnId": str(permit.get("rootTurnId") or ""),
                    "assignmentEpoch": str(permit.get("assignmentEpoch") or ""),
                    "authorizedRevision": revision,
                    "sourceToolUseId": spawn_call_id,
                }
                root_state["subagentEpochs"] = epochs
                root_state["subagentLifecycles"] = lifecycles
            root_store.save_unlocked(root_state)
            lifecycle = (root_state.get("subagentLifecycles") or {}).get(
                child_agent_id, {})
            active_revision = (int(lifecycle.get("revision") or 0)
                               if isinstance(lifecycle, Mapping)
                               and lifecycle.get("status") == "ACTIVE" else 0)
            stopped_revision = (int(lifecycle.get("revision") or 0)
                                if isinstance(lifecycle, Mapping)
                                and lifecycle.get("status") == "STOPPED" else 0)
        with store.locked():
            child_state = store.load_unlocked()
            local_stop_revision = int(
                child_state.get("subagentStoppedRevision") or
                (1 if child_state.get("stopped") else 0))
            local_stop_revision = max(local_stop_revision, stopped_revision)
            child_state.update({
                "schemaVersion": 1,
                "client": "codex",
                "clientSessionId": root_session_id,
                "cwd": str(event.get("cwd") or meta.get("cwd") or os.getcwd()),
                "transcriptPath": str(event.get("transcript_path") or ""),
                "subagentBinding": {
                    "agentId": child_agent_id,
                    "agentType": str(event.get("agent_type") or ""),
                    "agentPath": agent_path,
                    "parentThreadId": parent_thread_id,
                    "rootSessionId": root_session_id,
                    "bindingRevision": revision,
                },
                "stopped": not bool(active_revision > local_stop_revision),
                "registeredAt": int(time.time()),
            })
            if local_stop_revision:
                child_state["subagentStoppedRevision"] = local_stop_revision
            child_state.pop("subagentBindingPending", None)
            child_state.pop("lastError", None)
            store.save_unlocked(child_state)
        return None
    except (KeyError, RuntimeError, ValueError) as exc:
        with store.locked():
            pending = store.load_unlocked()
            pending.update({
                "schemaVersion": 1,
                "client": "codex",
                "clientSessionId": root_session_id,
                "cwd": str(event.get("cwd") or os.getcwd()),
                "transcriptPath": str(event.get("transcript_path") or ""),
                "agentId": child_agent_id,
                "agentType": str(event.get("agent_type") or ""),
                "lastError": str(exc),
            })
            if not _subagent_binding(pending):
                pending["subagentBindingPending"] = True
            pending["stopped"] = bool(pending.get("stopped"))
            store.save_unlocked(pending)
        return str(exc)


def _record_codex_turn(store: StateStore, event_name: str,
                       event: Mapping[str, Any]) -> None:
    """Record turn liveness locally; turn hooks must not depend on the network."""
    now = int(time.time())
    turn_id = str(event.get("turn_id") or "").strip()
    with store.locked():
        state = store.load_unlocked()
        if (not state
                or state.get("client") != "codex"
                or str(state.get("clientSessionId")) != str(event.get("session_id"))):
            return
        active_turn_id = str(state.get("activeTurnId") or "").strip()
        if event_name == "UserPromptSubmit":
            state["turnActive"] = True
            state["activeTurnId"] = turn_id or None
            state.pop("turnStoppedAt", None)
        elif event_name == "Stop":
            # Matching turn hooks may finish out of order.  A delayed Stop from
            # an older turn must never pause a newer turn's session heartbeat.
            if turn_id and active_turn_id and turn_id != active_turn_id:
                return
            state["turnActive"] = False
            if turn_id:
                state["activeTurnId"] = turn_id
            state["turnStoppedAt"] = now
        elif event_name in ("PreToolUse", "PostToolUse"):
            # A long turn remains live while tools continue to make progress.
            # Do not toggle turnActive here: only prompt/Stop own that boundary.
            pass
        else:
            return
        state["lastTurnActivityAt"] = now
        store.save_unlocked(state)


def _standalone_bash_tokens(event: Mapping[str, Any]) -> Optional[list]:
    """Tokenize a Bash tool command iff it is a single standalone invocation.

    Rejects shell composition/substitution (newlines, backticks, ``$(``,
    unquoted ``;&|<>``) so a recovery escape hatch cannot be used to append an
    unrelated mutation. Returns None for non-Bash tools and composite shapes.
    """
    # Both Codex 0.144.x and Claude expose their real shell tool as ``Bash``.
    # Never grant an escape hatch to an MCP/custom tool that merely carries a
    # command-looking input.
    if str(event.get("tool_name") or "").strip() != "Bash":
        return None
    tool_input = event.get("tool_input")
    if not isinstance(tool_input, Mapping):
        return None
    command = tool_input.get("command")
    if command is None:
        command = tool_input.get("cmd")
    if not isinstance(command, str) or not command.strip():
        return None
    if any(fragment in command for fragment in ("\n", "\r", "`", "$(")):
        return None
    try:
        lexer = shlex.shlex(command, posix=True, punctuation_chars=";&|<>")
        lexer.whitespace_split = True
        lexer.commenters = ""
        tokens = list(lexer)
    except ValueError:
        return None
    if not tokens or any(re.fullmatch(r"[;&|<>]+", token) for token in tokens):
        return None
    return tokens


def _exact_standard_claim(event: Mapping[str, Any]) -> Optional[Tuple[str, str]]:
    """Recognize the one command allowed to recover a fenced local Worker.

    A failed-closed PreToolUse hook must still leave one escape hatch: the
    standard database-first ``claim.sh claim`` flow.  Parse a deliberately
    narrow command shape and reject shell composition/substitution so this
    exception cannot be used to append an unrelated mutation.
    """
    tokens = _standalone_bash_tokens(event)
    if tokens is None:
        return None

    expected_script = str((REPO_ROOT / "bootstrap" / "claim.sh").resolve())
    # A fixed shell path and exact absolute script path prevent PATH, BASH_ENV,
    # per-command env overrides, malicious workdirs, symlinks and lookalike
    # repositories from turning the recovery escape hatch into arbitrary code.
    if len(tokens) != 5 or tokens[0] != "/bin/bash" or tokens[1] != expected_script:
        return None
    remaining = tokens[2:]
    if (remaining[0] != "claim"
            or not remaining[1].isdigit() or not remaining[2].isdigit()):
        return None
    return remaining[1], remaining[2]


def _exact_readonly_diagnostic(event: Mapping[str, Any]) -> bool:
    """Allow only argument-validated, side-effect-free recovery diagnostics."""
    tokens = _standalone_bash_tokens(event)
    if tokens is None or not tokens or tokens[0] != "/bin/bash":
        return False
    status_script = str((REPO_ROOT / "bootstrap" / "control-plane-status.sh").resolve())
    config_script = str((REPO_ROOT / "bootstrap" / "runtime-config.sh").resolve())
    hook_script = str((REPO_ROOT / "bootstrap" / "run-interactive-worker-hook.sh").resolve())
    if tokens == ["/bin/bash", config_script, "diagnose"]:
        return True
    if tokens == ["/bin/bash", hook_script, "cli", "status"]:
        return True
    if len(tokens) == 3 and tokens[1] == status_script and tokens[2] == "workers":
        return True
    if (len(tokens) == 4 and tokens[1] == status_script
            and tokens[2] in ("task", "operation") and tokens[3].isdigit()):
        return True
    if (len(tokens) in (3, 5) and tokens[1] == status_script
            and tokens[2] == "ready"):
        return len(tokens) == 3 or (tokens[3] == "--limit" and tokens[4].isdigit())
    return False


def _external_receipt_recovery_target(
        state: Mapping[str, Any]) -> Optional[Tuple[str, str]]:
    """(aone_id, project_id) when the pending op is a frozen external-write
    receipt (kind∈comment/status/release-tag/finish-tag, UNKNOWN/RETRY_WAIT)
    on the CURRENT assignment — the only shape whose convergence rerun
    (wrap.sh / claim.sh release|finish) may pass the PreToolUse fence.
    AONE_CLAIM receipts have no ``kind`` and keep the claim-only escape hatch.
    """
    if state.get("stopped"):
        return None
    pending = state.get("pendingOperation")
    current = state.get("current")
    if not isinstance(pending, Mapping) or not isinstance(current, Mapping):
        return None
    if str(pending.get("kind") or "") not in OPERATION_TYPE_BY_KIND:
        return None
    if str(pending.get("status") or "").upper() not in ("UNKNOWN", "RETRY_WAIT"):
        return None
    aone_id = str(pending.get("aoneId") or "").strip()
    if not aone_id or aone_id != str(current.get("aoneId") or "").strip():
        return None
    return aone_id, str(current.get("projectId") or "").strip()


def _exact_receipt_recovery(
        event: Mapping[str, Any]) -> Optional[Tuple[str, str, Optional[str]]]:
    """Recognize the exact standalone wrap.sh / claim.sh commands allowed to
    converge a frozen external-write receipt.

    Same narrowness as _exact_standard_claim: fixed ``/bin/bash``, exact
    absolute script path inside THIS repository, exact subcommand and Aone id.
    wrap.sh keeps free TRAILING arguments (--summary-file/--summary-stdin/
    status …) because the comment body is data, not shell — composition is
    already rejected by _standalone_bash_tokens. Returns
    ("wrap:<subcommand>"|"claim:<subcommand>", aone_id, project_id_or_None).
    """
    tokens = _standalone_bash_tokens(event)
    if tokens is None or len(tokens) < 4 or tokens[0] != "/bin/bash":
        return None
    wrap_script = str((REPO_ROOT / "bootstrap" / "wrap.sh").resolve())
    claim_script = str((REPO_ROOT / "bootstrap" / "claim.sh").resolve())
    if tokens[1] == wrap_script:
        if tokens[2] not in ("sync", "done", "done-no-status"):
            return None
        if not tokens[3].isdigit():
            return None
        return "wrap:%s" % tokens[2], tokens[3], None
    if tokens[1] == claim_script:
        if (len(tokens) != 5 or tokens[2] not in ("release", "finish")
                or not tokens[3].isdigit() or not tokens[4].isdigit()):
            return None
        return "claim:%s" % tokens[2], tokens[3], tokens[4]
    return None


def _receipt_recovery_matches_kind(
        state: Mapping[str, Any], command: Tuple[str, str, Optional[str]]) -> bool:
    pending = state.get("pendingOperation")
    if not isinstance(pending, Mapping):
        return False
    kind = str(pending.get("kind") or "")
    action = command[0]
    if kind == "comment":
        return action in ("wrap:sync", "wrap:done", "wrap:done-no-status")
    if kind == "status":
        return action == "wrap:done"
    if kind == "release-tag":
        return action == "claim:release"
    if kind == "finish-tag":
        return action == "claim:finish"
    return False


def _exact_operation_recovery(event: Mapping[str, Any]) -> Optional[str]:
    """Recognize direct, fenced receipt abort/reconcile recovery commands."""
    tokens = _standalone_bash_tokens(event)
    hook_script = str((REPO_ROOT / "bootstrap" / "run-interactive-worker-hook.sh").resolve())
    if (tokens is None or len(tokens) < 6 or tokens[0] != "/bin/bash"
            or tokens[1] != hook_script or tokens[2] != "cli"
            or not tokens[4].isdigit()):
        return None
    if tokens[3] == "operation-abort":
        if len(tokens) == 6 or (len(tokens) == 7 and tokens[6] == "--unknown"):
            return tokens[4]
        return None
    if tokens[3] != "operation-reconcile":
        return None
    if len(tokens) == 6 and tokens[5] == "--not-found":
        return tokens[4]
    if (len(tokens) == 7
            and ((tokens[5] == "--not-found" and tokens[6] == "--no-retry")
                 or tokens[5] == "--found")):
        return tokens[4]
    return None


def _receipt_recovery_hint(state: Mapping[str, Any]) -> str:
    target = _external_receipt_recovery_target(state)
    pending = state.get("pendingOperation")
    wrap_script = shlex.quote(str((REPO_ROOT / "bootstrap" / "wrap.sh").resolve()))
    claim_script = shlex.quote(str((REPO_ROOT / "bootstrap" / "claim.sh").resolve()))
    status_script = shlex.quote(str(
        (REPO_ROOT / "bootstrap" / "control-plane-status.sh").resolve()))
    hook_script = shlex.quote(str(
        (REPO_ROOT / "bootstrap" / "run-interactive-worker-hook.sh").resolve()))
    aone = target[0] if target else "<aone-id>"
    project = target[1] if target and target[1] else "<project-id>"
    kind = (str(pending.get("kind") or "")
            if isinstance(pending, Mapping) else "")
    operation_id = (str(pending.get("operationId") or "").strip()
                    if isinstance(pending, Mapping) else "")
    rerun_by_kind = {
        "comment": "/bin/bash %s sync|done|done-no-status %s …" %
                   (wrap_script, aone),
        "status": "/bin/bash %s done %s …" % (wrap_script, aone),
        "release-tag": "/bin/bash %s release %s %s" %
                       (claim_script, aone, project),
        "finish-tag": "/bin/bash %s finish %s %s" %
                      (claim_script, aone, project),
    }
    hints = []
    if operation_id.isdigit():
        hints.append("只读 point-read：/bin/bash %s operation %s" %
                     (status_script, operation_id))
    hints.extend((
        "证据 readback 后 reconcile：/bin/bash %s cli operation-reconcile %s "
        "--found <external-ref>|--not-found [--no-retry]" %
        (hook_script, aone),
        "确认未开始时 abort-not-started：/bin/bash %s cli operation-abort %s "
        "not-started" % (hook_script, aone),
    ))
    rerun = rerun_by_kind.get(kind)
    if rerun:
        hints.append("按 %s 原操作安全重跑：%s" % (kind, rerun))
    return "；".join(hints)


def _standard_claim_hint(state: Mapping[str, Any]) -> str:
    target: Optional[Tuple[str, str]] = None
    for key in ("current", "pendingClaim", "lostOwnership", "recoveryPending"):
        candidate = state.get(key)
        if not isinstance(candidate, Mapping):
            continue
        aone_id = str(candidate.get("aoneId") or "").strip()
        project_id = str(candidate.get("projectId") or "").strip()
        if aone_id and project_id:
            target = (aone_id, project_id)
            break
    script = shlex.quote(str((REPO_ROOT / "bootstrap" / "claim.sh").resolve()))
    if target:
        return "/bin/bash %s claim %s %s" % (script, target[0], target[1])
    return "/bin/bash %s claim <aone-id> <project-id>" % script


def _local_tool_block_reason(state: Mapping[str, Any]) -> Optional[str]:
    claim_hint = _standard_claim_hint(state)
    if state.get("stopped"):
        return ("Jarvis Worker 已离线；当前工具调用已阻断。"
                "请先重新触发 SessionStart，再通过 claim.sh 接单。")
    lost = state.get("lostOwnership")
    if isinstance(lost, Mapping):
        return ("Jarvis Worker 已失去 Aone %s 的数据库 fence；当前工具调用已阻断。"
                "唯一允许的恢复命令：%s" %
                (str(lost.get("aoneId") or "unknown"), claim_hint))
    pending_claim = state.get("pendingClaim")
    if isinstance(pending_claim, Mapping):
        if pending_claim.get("receiptUnknown"):
            return ("Jarvis 接单回执处于 UNKNOWN；当前工具调用已阻断。"
                    "唯一允许的恢复命令：%s" % claim_hint)
        return ("Jarvis 存在未完成的接单意图（%s）；当前工具调用已阻断。"
                "唯一允许的恢复命令：%s" %
                (str(pending_claim.get("phase") or "UNKNOWN"), claim_hint))
    if state.get("pendingOperation"):
        if _external_receipt_recovery_target(state):
            return ("Jarvis 存在未收敛的外部写回执（UNKNOWN/RETRY_WAIT）；"
                    "当前工具调用已阻断。仅允许重跑收敛命令：%s" %
                    _receipt_recovery_hint(state))
        return ("Jarvis 存在未确定化的外部操作回执；当前工具调用已阻断。"
                "唯一允许的恢复命令：%s" % claim_hint)
    if state.get("pendingSuspend"):
        return ("Jarvis 任务挂起结果尚未确定；当前工具调用已阻断，"
                "请先恢复控制面并重新开始一轮对话。")
    if state.get("lastAutoSuspended") or state.get("recoveryPending"):
        return ("Jarvis 任务尚未重新取得数据库 fence；当前工具调用已阻断。"
                "唯一允许的恢复命令：%s" % claim_hint)
    current = state.get("current")
    if isinstance(current, Mapping) and not current.get("heartbeatEnabled", True):
        return ("Jarvis 当前任务的外部操作回执尚未 ACK；当前工具调用已阻断。"
                "唯一允许的恢复命令：%s" % claim_hint)
    return None


def _recovery_claim_targets(state: Mapping[str, Any]) -> set[Tuple[str, str]]:
    targets: set[Tuple[str, str]] = set()
    for key in ("current", "pendingClaim", "lostOwnership", "recoveryPending"):
        candidate = state.get(key)
        if not isinstance(candidate, Mapping):
            continue
        aone_id = str(candidate.get("aoneId") or "").strip()
        project_id = str(candidate.get("projectId") or "").strip()
        if aone_id and project_id:
            targets.add((aone_id, project_id))
    return targets


def _calling_process_matches(state: Mapping[str, Any], client_name: str) -> bool:
    try:
        host_pid, verified = _find_host_pid(client_name)
        if not verified or int(state.get("hostPid") or 0) != int(host_pid):
            return False
    except (OSError, TypeError, ValueError):
        return False
    expected_start = str(state.get("hostProcessStartedAt") or "")
    actual_start = _process_start_identity(host_pid)
    return bool(expected_start and actual_start and expected_start == actual_start)


def _codex_turn_block_reason(state: Mapping[str, Any],
                             event: Mapping[str, Any],
                             binding: Optional[Mapping[str, Any]] = None) -> Optional[str]:
    if not state.get("turnActive", True):
        return ("Jarvis Codex turn 已停止；延迟到达的工具调用已阻断，"
                "请先开始新一轮对话。")
    active_turn = str(state.get("activeTurnId") or "").strip()
    if binding is not None:
        agent_id = str(event.get("agent_id") or "").strip()
        expected_agent_id = str(binding.get("agentId") or "").strip()
        if not agent_id or agent_id != expected_agent_id:
            return ("Jarvis Codex subagent 身份不匹配；当前工具调用已阻断，"
                    "请由根任务重新派发该数字人。")
        lifecycles = state.get("subagentLifecycles")
        lifecycle = (lifecycles.get(agent_id)
                     if isinstance(lifecycles, Mapping) else None)
        epochs = state.get("subagentEpochs")
        epoch = epochs.get(agent_id) if isinstance(epochs, Mapping) else None
        lifecycle_revision = (int(lifecycle.get("revision") or 0)
                              if isinstance(lifecycle, Mapping) else 0)
        child_stopped_revision = int(
            binding.get("_childStoppedRevision") or 0)
        if (not isinstance(lifecycle, Mapping)
                or lifecycle.get("status") != "ACTIVE"
                or not isinstance(epoch, Mapping)
                or int(epoch.get("authorizedRevision") or 0) != lifecycle_revision
                or lifecycle_revision <= child_stopped_revision
                or str(epoch.get("rootTurnId") or "") != active_turn
                or str(epoch.get("assignmentEpoch") or "") != _assignment_epoch(state)):
            return ("Jarvis Codex subagent 不属于当前根 turn/任务 fence；"
                    "延迟或跨任务工具调用已阻断，请由当前根任务重新派发。")
        return None
    if event.get("agent_id") or event.get("agent_type"):
        return ("Jarvis Codex subagent 缺少可信父任务绑定；当前工具调用已阻断，"
                "请重新触发 SubagentStart。")
    event_turn = str(event.get("turn_id") or "").strip()
    if event_turn and active_turn and event_turn != active_turn:
        return ("Jarvis Codex turn 已被更新；旧 turn 的工具调用已阻断，"
                "请在当前轮重新发起操作。")
    return None


def _authority_context(store: StateStore, client_name: str,
                       event: Mapping[str, Any]) -> Tuple[
                           StateStore, Dict[str, Any], Optional[Dict[str, Any]],
                           Optional[str]]:
    event_state = store.load()
    if client_name == "codex" and event.get("agent_id"):
        if not event_state or event_state.get("subagentBindingPending"):
            _bind_codex_subagent(store, event)
            event_state = store.load()
        binding = _subagent_binding(event_state)
        if binding is None:
            return store, event_state, None, (
                "Jarvis Codex subagent 未完成可信父任务绑定；当前工具调用已阻断。"
                "请重新触发 SubagentStart。")
        root_session_id = str(event.get("session_id") or "").strip()
        if (str(event_state.get("clientSessionId") or "") != root_session_id):
            return store, event_state, binding, (
                "Jarvis Codex subagent 会话标识不匹配；当前工具调用已阻断。")
        if (str(binding.get("agentId") or "") !=
                str(event.get("agent_id") or "")):
            return store, event_state, binding, (
                "Jarvis Codex subagent 身份标识不匹配；当前工具调用已阻断。")
        if str(binding.get("rootSessionId") or "") != root_session_id:
            return store, event_state, binding, (
                "Jarvis Codex subagent 根 Worker 标识不匹配；当前工具调用已阻断。")
        binding["_childStoppedRevision"] = int(
            event_state.get("subagentStoppedRevision") or
            (1 if event_state.get("stopped") else 0))
        authority_store = StateStore(_state_path("codex", root_session_id))
        authority_state = authority_store.load()
        if not authority_state:
            return authority_store, authority_state, binding, (
                "Jarvis Codex 根 Worker 状态不存在；当前工具调用已阻断。")
        return authority_store, authority_state, binding, None
    if not event_state:
        return store, event_state, None, (
            "Jarvis Worker 尚未完成 SessionStart 注册；当前工具调用已阻断。"
            "请信任/重新打开当前项目并重新触发 SessionStart。")
    return store, event_state, None, None


def _guard_pre_tool_use(store: StateStore, client_name: str,
                        event: Mapping[str, Any]) -> Optional[str]:
    """Fence every tool call while an interactive task is locally attached."""
    # Direct a1 bypasses bin/a1id's CR-exit deny. Reject it before any ordinary
    # Worker permit path can authorize the Bash tool call.
    a1_reason = pretool_a1_block_reason(event)
    if a1_reason:
        return a1_reason
    authority_store, state, binding, context_error = _authority_context(
        store, client_name, event)
    if context_error:
        return context_error
    if (str(state.get("client") or "") != client_name
            or (binding is None
                and str(state.get("clientSessionId") or "") !=
                str(event.get("session_id") or ""))):
        return ("Jarvis Worker 会话标识不匹配；当前工具调用已阻断，"
                "请重新触发 SessionStart。")
    if not _calling_process_matches(state, client_name):
        return ("Jarvis Worker 进程 incarnation 已变化；旧会话的当前工具调用已阻断，"
                "请在新会话中通过 claim.sh 接单。")
    if client_name == "codex":
        turn_reason = _codex_turn_block_reason(state, event, binding)
        if turn_reason:
            return turn_reason

    # UNKNOWN and terminal-state recovery must never hide the evidence needed
    # to choose a lawful next action. These exact commands are read-only and do
    # not weaken task/session ownership fencing for mutations.
    if _exact_readonly_diagnostic(event):
        return None

    # claim.sh is itself the database-first recovery gate.  It is the only
    # command allowed through an uncertain/lost local state, and only for the
    # exact carried task. Its local CAS, server slot check and receipt still
    # fail closed. With no carried lineage it is the normal first-claim gate.
    claim_target = _exact_standard_claim(event)
    local_reason = _local_tool_block_reason(state)
    if claim_target:
        if binding is not None:
            return ("Jarvis 只允许根 Worker 执行标准 claim.sh；"
                    "subagent 当前工具调用已阻断。")
        recovery_targets = _recovery_claim_targets(state)
        if not recovery_targets and not local_reason:
            return None
        if claim_target in recovery_targets and not state.get("stopped"):
            return None
        return ("Jarvis 只允许对当前挂起/失权的同一 Aone 单独重试 claim.sh；"
                "当前 claim 目标不匹配，工具调用已阻断。")
    # Frozen external-write receipts (wrap comment/status, release/finish tag,
    # UNKNOWN/RETRY_WAIT) would otherwise brick the session: the only way to
    # converge them is rerunning wrap.sh / claim.sh release|finish, which is
    # itself a fenced Bash tool call. Allow exactly those two standalone
    # command shapes, only for the SAME Aone id as the current assignment;
    # the actual convergence is still fail-closed inside the worker CLI's
    # begin/readback/reconcile path.
    recovery_command = _exact_receipt_recovery(event)
    if recovery_command and local_reason:
        if binding is None:
            recovery_target = _external_receipt_recovery_target(state)
            if (recovery_target is not None
                    and recovery_command[1] == recovery_target[0]
                    and (recovery_command[2] is None
                         or recovery_command[2] == recovery_target[1])
                    and _receipt_recovery_matches_kind(state, recovery_command)):
                return None
        # 不匹配（subagent / 别的 aone / 无可恢复回执）→ 维持原有阻断语义。
    operation_recovery_target = _exact_operation_recovery(event)
    if operation_recovery_target and local_reason and binding is None:
        recovery_target = _external_receipt_recovery_target(state)
        if (recovery_target is not None
                and operation_recovery_target == recovery_target[0]):
            return None
    if local_reason:
        return local_reason
    current = state.get("current")
    if not isinstance(current, Mapping):
        with authority_store.locked():
            latest = authority_store.load_unlocked()
            if not _calling_process_matches(latest, client_name):
                return ("Jarvis Worker 进程在工具执行前发生变化；旧进程调用已阻断。")
            if client_name == "codex":
                turn_reason = _codex_turn_block_reason(latest, event, binding)
                if turn_reason:
                    return turn_reason
            local_reason = _local_tool_block_reason(latest)
            if local_reason:
                return local_reason
            _record_spawn_permit_locked(latest, event)
            _record_subagent_interaction_permit_locked(latest, event)
            latest["lastTurnActivityAt"] = int(time.time())
            authority_store.save_unlocked(latest)
        return None

    expected = dict(current)
    worker_key = str(state.get("workerKey") or "")
    if not worker_key:
        return ("Jarvis Worker 本地状态缺少 workerKey；当前工具调用已阻断，"
                "请重新启动会话并走 claim.sh。")
    permit_reason = _session_permit_block_reason(state, expected)
    if permit_reason:
        return permit_reason

    # PreToolUse is deliberately local-only. The detached sidecar owns all
    # remote lease renewals; this final locked recheck ensures its permit cannot
    # race with SessionStart, claim, suspend, Stop or a newer assignment.
    with authority_store.locked():
        latest = authority_store.load_unlocked()
        if not _same_current_assignment(latest, worker_key, expected):
            return ("Jarvis Worker 归属在工具执行前发生变化；当前工具调用已阻断，"
                    "请重新开始一轮并走 claim.sh。")
        if not _calling_process_matches(latest, client_name):
            return ("Jarvis Worker 进程在工具执行前发生变化；旧进程调用已阻断，"
                    "请在新会话中通过 claim.sh 接单。")
        if client_name == "codex":
            turn_reason = _codex_turn_block_reason(latest, event, binding)
            if turn_reason:
                return turn_reason
        local_reason = _local_tool_block_reason(latest)
        if local_reason:
            return local_reason
        permit_reason = _session_permit_block_reason(
            latest, expected, now=time.time())
        if permit_reason:
            return permit_reason
        _record_spawn_permit_locked(latest, event)
        _record_subagent_interaction_permit_locked(latest, event)
        latest["lastTurnActivityAt"] = int(time.time())
        authority_store.save_unlocked(latest)
    return None


def _mark_subagent_authorized(root_session_id: str, agent_id: str,
                              revision: int) -> None:
    child_store = StateStore(_state_path("codex", agent_id))
    with child_store.locked():
        state = child_store.load_unlocked()
        if not state:
            return
        binding = _subagent_binding(state)
        if (binding is None
                or str(binding.get("agentId") or "") != agent_id
                or str(binding.get("rootSessionId") or "") != root_session_id):
            return
        stopped_revision = int(state.get("subagentStoppedRevision") or 0)
        if revision <= stopped_revision:
            return
        state["stopped"] = False
        state["lastAuthorizedRevision"] = revision
        state.pop("subagentStoppedAt", None)
        state.pop("stoppedTurnId", None)
        child_store.save_unlocked(state)


def _record_codex_post_tool(store: StateStore,
                            event: Mapping[str, Any]) -> None:
    """Refresh root activity and authorize only successful explicit child sends."""
    authority_store, _state, binding, context_error = _authority_context(
        store, "codex", event)
    if context_error:
        return
    authorized_agent_id = ""
    authorized_revision = 0
    with authority_store.locked():
        latest = authority_store.load_unlocked()
        if not _calling_process_matches(latest, "codex"):
            return
        if _codex_turn_block_reason(latest, event, binding):
            return
        latest["lastTurnActivityAt"] = int(time.time())
        authorized_agent_id, authorized_revision = (
            _consume_subagent_interaction_locked(latest, event))
        authority_store.save_unlocked(latest)
    if authorized_agent_id:
        _mark_subagent_authorized(
            str(event.get("session_id") or ""), authorized_agent_id,
            authorized_revision)


def _stop_codex_subagent(store: StateStore,
                         event: Mapping[str, Any]) -> None:
    root_session_id = _nonblank(event.get("session_id"), "session_id")
    agent_id = _nonblank(event.get("agent_id"), "agent_id")
    root_store = StateStore(_state_path("codex", root_session_id))
    revision = 0
    with root_store.locked():
        root_state = root_store.load_unlocked()
        if (root_state
                and str(root_state.get("clientSessionId") or "") == root_session_id):
            revision = _next_subagent_revision(root_state)
            lifecycles = dict(root_state.get("subagentLifecycles") or {})
            lifecycles[agent_id] = {
                "status": "STOPPED",
                "revision": revision,
                "reason": "subagent_stop",
                "stoppedTurnId": str(event.get("turn_id") or ""),
            }
            epochs = dict(root_state.get("subagentEpochs") or {})
            epochs.pop(agent_id, None)
            root_state["subagentLifecycles"] = lifecycles
            root_state["subagentEpochs"] = epochs
            root_store.save_unlocked(root_state)
    with store.locked():
        state = store.load_unlocked()
        old_revision = int(state.get("subagentStoppedRevision") or 0)
        if revision and revision < old_revision:
            return
        state.update({
            "schemaVersion": 1,
            "client": "codex",
            "clientSessionId": root_session_id,
            "agentId": agent_id,
            "agentType": str(event.get("agent_type") or ""),
            "cwd": str(event.get("cwd") or state.get("cwd") or os.getcwd()),
            "transcriptPath": str(
                event.get("agent_transcript_path") or
                state.get("transcriptPath") or ""),
            "stopped": True,
            "subagentStoppedAt": int(time.time()),
            "stoppedTurnId": str(event.get("turn_id") or ""),
        })
        if revision:
            state["subagentStoppedRevision"] = revision
        store.save_unlocked(state)


def _codex_turn_live(state: Mapping[str, Any],
                     now: Optional[float] = None) -> bool:
    if not state.get("turnActive", True):
        return False
    try:
        active_ttl = max(0.0, float(os.environ.get(
            "JARVIS_INTERACTIVE_ACTIVE_TURN_TTL_SEC", "43200")))
    except ValueError:
        active_ttl = 43200.0
    if active_ttl <= 0:
        return True
    try:
        started_at = float(state["lastTurnActivityAt"])
    except (KeyError, TypeError, ValueError):
        # Compatibility for states written before turn timestamps existed.
        return True
    current_time = time.time() if now is None else float(now)
    return current_time - started_at < active_ttl


def _interactive_turn_grace_seconds() -> float:
    """Return the Codex between-turn grace."""
    configured = os.environ.get("JARVIS_INTERACTIVE_TURN_GRACE_SEC", "600")
    try:
        return max(0.0, float(configured))
    except ValueError:
        return 600.0


def _interactive_heartbeat_seconds() -> float:
    """Return the detached interactive sidecar heartbeat interval."""
    try:
        return max(2.0, float(os.environ.get(
            "JARVIS_INTERACTIVE_HEARTBEAT_SEC", "30")))
    except ValueError:
        return 30.0


def _interactive_lease_seconds() -> int:
    """Return the targeted interactive Task lease duration."""
    try:
        return max(30, int(os.environ.get(
            "JARVIS_INTERACTIVE_LEASE_SECONDS", "660")))
    except ValueError:
        return 660


def _interactive_lease_safety_margin_seconds() -> float:
    """Return the minimum remaining lease required before a tool may start."""
    try:
        return max(0.0, float(os.environ.get(
            "JARVIS_INTERACTIVE_LEASE_SAFETY_MARGIN_SEC", "60")))
    except ValueError:
        return 60.0


def _interactive_affinity_seconds() -> float:
    """Return how long an auto-suspended turn remains affined to this Worker."""
    try:
        return max(0.0, float(os.environ.get(
            "JARVIS_INTERACTIVE_AFFINITY_SEC", "7200")))
    except ValueError:
        return 7200.0


def _epoch_seconds(value: Any) -> Optional[float]:
    if value is None or isinstance(value, bool):
        return None
    if isinstance(value, (int, float)):
        numeric = float(value)
        return numeric / 1000.0 if numeric > 100_000_000_000 else numeric
    text = str(value).strip()
    if not text:
        return None
    try:
        numeric = float(text)
        return numeric / 1000.0 if numeric > 100_000_000_000 else numeric
    except ValueError:
        pass
    try:
        return datetime.fromisoformat(text.replace("Z", "+00:00")).timestamp()
    except ValueError:
        return None


def _utc_timestamp(value: float) -> str:
    return datetime.fromtimestamp(
        float(value), timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _session_response(value: Any) -> Mapping[str, Any]:
    if not isinstance(value, Mapping):
        return {}
    session = value.get("session")
    if isinstance(session, Mapping):
        return session
    data = value.get("data")
    if isinstance(data, Mapping):
        session = data.get("session")
        if isinstance(session, Mapping):
            return session
    return value


def _session_permit(state: Mapping[str, Any],
                    current: Mapping[str, Any],
                    response: Any,
                    *,
                    source: str,
                    now: Optional[float] = None) -> Dict[str, Any]:
    """Bind one successful start/heartbeat response to the exact local fence."""
    issued_at = time.time() if now is None else float(now)
    session = _session_response(response)
    lease_expire_at = _epoch_seconds(_field(
        session, "leaseExpireAt", "leaseExpiresAt", "lease_expire_at"))
    if lease_expire_at is None:
        lease_expire_at = _epoch_seconds(_field(
            current, "leaseExpireAt", "leaseExpiresAt", "lease_expire_at"))
    if lease_expire_at is None:
        # Current control-plane deployments do not consistently echo
        # leaseExpireAt. A successful fenced start/heartbeat still proves a
        # fresh lease for the requested duration, so retain a conservative
        # local deadline until the response contract is upgraded.
        lease_expire_at = issued_at + int(
            current.get("leaseSeconds") or _interactive_lease_seconds())
    status = str(_field(session, "status", "sessionStatus") or "RUNNING").upper()
    return {
        "version": 1,
        "workerKey": str(state.get("workerKey") or ""),
        "taskId": current.get("taskId"),
        "sessionId": current.get("sessionId"),
        "fenceToken": current.get("fenceToken"),
        "generation": current.get("generation"),
        "runtimeSessionId": current.get("runtimeSessionId"),
        "sessionStatus": status,
        "leaseExpireAt": lease_expire_at,
        "issuedAt": issued_at,
        "source": source,
    }


def _refresh_session_permit_locked(state: Dict[str, Any],
                                   expected: Mapping[str, Any],
                                   response: Any,
                                   *,
                                   source: str,
                                   now: Optional[float] = None) -> bool:
    if not _same_current_assignment(state, state.get("workerKey"), expected):
        return False
    state["sessionPermit"] = _session_permit(
        state, expected, response, source=source, now=now)
    return True


def _sidecar_health_block_reason(state: Mapping[str, Any],
                                 now: float) -> Optional[str]:
    daemon_pid = state.get("daemonPid")
    if not _pid_alive(daemon_pid):
        return ("Jarvis interactive sidecar 未运行；当前工具调用已阻断，"
                "请重新触发 UserPromptSubmit/SessionStart。")
    expected_start = str(state.get("daemonProcessStartedAt") or "")
    if expected_start and _process_start_identity(int(daemon_pid)) != expected_start:
        return ("Jarvis interactive sidecar incarnation 已变化；当前工具调用已阻断，"
                "请重新触发 SessionStart。")
    health_at = _epoch_seconds(
        state.get("sidecarHeartbeatAt") or state.get("daemonStartedAt"))
    max_age = max(
        _interactive_heartbeat_seconds() * 3,
        _interactive_lease_safety_margin_seconds())
    if health_at is None or now - health_at > max_age:
        return ("Jarvis interactive sidecar 健康状态已过期；当前工具调用已阻断，"
                "请等待 sidecar 恢复或重新触发 SessionStart。")
    return None


def _session_permit_block_reason(state: Mapping[str, Any],
                                 current: Mapping[str, Any],
                                 *,
                                 now: Optional[float] = None) -> Optional[str]:
    current_time = time.time() if now is None else float(now)
    permit = state.get("sessionPermit")
    if not isinstance(permit, Mapping):
        return ("Jarvis 当前 Session 尚无本地 Lease Proof；当前工具调用已阻断，"
                "请等待 sidecar 完成一次续租或重新通过 claim.sh 接单。")
    expected_fields = {
        "workerKey": state.get("workerKey"),
        "taskId": current.get("taskId"),
        "sessionId": current.get("sessionId"),
        "fenceToken": current.get("fenceToken"),
        "generation": current.get("generation"),
        "runtimeSessionId": current.get("runtimeSessionId"),
    }
    for key, expected in expected_fields.items():
        if str(permit.get(key)) != str(expected):
            return ("Jarvis 本地 Lease Proof 与当前 task/session/fence 不匹配；"
                    "当前工具调用已阻断，请重新通过 claim.sh 取得归属。")
    status = str(permit.get("sessionStatus") or "").upper()
    if status not in ("LEASED", "RUNNING"):
        return ("Jarvis 当前 Session 状态为 %s；当前工具调用已阻断。"
                % (status or "UNKNOWN"))
    lease_expire_at = _epoch_seconds(permit.get("leaseExpireAt"))
    safety_margin = _interactive_lease_safety_margin_seconds()
    if lease_expire_at is None:
        return ("Jarvis 本地 Lease Proof 缺少 leaseExpireAt；当前工具调用已阻断。")
    if lease_expire_at <= current_time + safety_margin:
        return ("Jarvis 当前 Session 剩余租约不足 %.0f 秒安全边界；"
                "当前工具调用已阻断，请等待 sidecar 续租。" % safety_margin)
    return _sidecar_health_block_reason(state, current_time)


def _session_heartbeat_allowed(state: Mapping[str, Any],
                               now: Optional[float] = None) -> bool:
    """Keep Claude process-scoped; bound Codex's global app-server lease."""
    if str(state.get("client") or "").lower() != "codex":
        return True
    if _codex_turn_live(state, now):
        return True
    try:
        stopped_at = float(state["turnStoppedAt"])
    except (KeyError, TypeError, ValueError):
        return False
    grace = _interactive_turn_grace_seconds()
    current_time = time.time() if now is None else float(now)
    elapsed = current_time - stopped_at
    return 0.0 <= elapsed < grace


def _worker_idle_expired(state: Mapping[str, Any],
                         now: Optional[float] = None) -> bool:
    """Bound idle Codex sidecars even though the shared app-server stays alive."""
    if str(state.get("client") or "").lower() != "codex":
        return False
    pending_claim = state.get("pendingClaim")
    if (_codex_turn_live(state, now)
            or isinstance(state.get("current"), Mapping)
            or state.get("pendingOperation")
            or state.get("pendingSuspend")
            or (isinstance(pending_claim, Mapping)
                and pending_claim.get("phase") == "CLAIMING")):
        return False
    try:
        ttl = float(os.environ.get("JARVIS_INTERACTIVE_IDLE_TTL_SEC", "28800"))
    except ValueError:
        ttl = 28800.0
    if ttl <= 0:
        return False
    try:
        last_activity = float(
            state.get("lastTurnActivityAt") or state.get("registeredAt"))
    except (TypeError, ValueError):
        return False
    current_time = time.time() if now is None else float(now)
    return current_time - last_activity >= ttl


def _clear_lost_current(store: StateStore, expected_worker_key: str,
                        expected_session_id: Any,
                        expected_fence_token: Any,
                        error: str) -> None:
    with store.locked():
        latest = store.load_unlocked()
        current = latest.get("current")
        if (latest.get("workerKey") == expected_worker_key
                and isinstance(current, Mapping)
                and str(current.get("sessionId")) == str(expected_session_id)
                and str(current.get("fenceToken")) == str(expected_fence_token)):
            latest["lastError"] = error
            latest["lostOwnership"] = {
                "aoneId": current.get("aoneId"),
                "projectId": current.get("projectId"),
                "taskId": current.get("taskId"),
                "sessionId": current.get("sessionId"),
                "runtimeSessionId": current.get("runtimeSessionId"),
                "lostAt": int(time.time()),
                "reason": error,
            }
            latest["current"] = None
            latest.pop("sessionPermit", None)
            latest["pendingClaim"] = None
            latest["pendingOperation"] = None
            latest["pendingSuspend"] = None
            store.save_unlocked(latest)


def _is_explicit_admin_cleanup_not_found(exc: BaseException) -> bool:
    """Recognize deletion of a Worker/Session row, not a transient outage."""
    return (isinstance(exc, ControlPlaneError)
            and exc.status == 404
            and str(exc.code or "").lower() in (
                "", "notfound.worker", "notfound.session"))


def _stop_admin_cleaned_sidecar_unlocked(
        state: Dict[str, Any], expected_worker_key: str,
        missing_resource: str, error: BaseException) -> bool:
    """Tombstone a remotely purged assignment while the daemon owns the lock."""
    if state.get("workerKey") != expected_worker_key:
        return False
    now = int(time.time())
    current = state.get("current")
    pending = state.get("pendingClaim")
    lineage = (current if isinstance(current, Mapping)
               else pending if isinstance(pending, Mapping)
               else None)
    reason = ("admin cleanup removed control-plane %s row"
              % missing_resource.lower())
    if isinstance(lineage, Mapping):
        aone_id = lineage.get("aoneId")
        project_id = lineage.get("projectId")
        if aone_id is not None and project_id is not None:
            state["lostOwnership"] = {
                "aoneId": aone_id,
                "projectId": project_id,
                "taskId": lineage.get("taskId"),
                "sessionId": lineage.get("sessionId"),
                "runtimeSessionId": lineage.get("runtimeSessionId"),
                "lostAt": now,
                "reason": reason,
                "cause": "ADMIN_CLEANUP",
                "missingResource": missing_resource,
            }
    state["adminCleanup"] = {
        "detectedAt": now,
        "workerKey": expected_worker_key,
        "missingResource": missing_resource,
        "reason": reason,
        "errorCode": str(getattr(error, "code", "") or ""),
    }
    state["lastError"] = reason
    state["current"] = None
    state.pop("sessionPermit", None)
    state["pendingClaim"] = None
    state["pendingOperation"] = None
    state["pendingSuspend"] = None
    state.pop("lastAutoSuspended", None)
    state.pop("recoveryPending", None)
    state["subagentRegistry"] = {}
    state["subagentEpochs"] = {}
    state["subagentLifecycles"] = {}
    state["subagentSpawnPermits"] = {}
    state["subagentInteractionPermits"] = {}
    state["stopped"] = True
    state["stoppedAt"] = now
    state["offlineReason"] = "admin_cleanup"
    state.pop("daemonPid", None)
    state.pop("daemonStartedAt", None)
    state.pop("daemonProcessStartedAt", None)
    state.pop("sidecarHeartbeatAt", None)
    return True


def _finalize_pending_suspend_unlocked(state: Dict[str, Any],
                                        pending: Mapping[str, Any]) -> None:
    current = state.get("current")
    if (not isinstance(current, Mapping)
            or str(current.get("sessionId")) != str(pending.get("sessionId"))
            or str(current.get("fenceToken")) != str(pending.get("fenceToken"))):
        return
    if (current.get("aoneId") is not None
            and current.get("projectId") is not None
            and current.get("cycle") is not None
            and current.get("runtimeSessionId")):
        # Automatic turn idling is not a user release. Preserve the exact
        # runtime so the next explicit claim resumes this suspended session.
        state["pendingClaim"] = {
            "aoneId": str(current["aoneId"]),
            "projectId": str(current["projectId"]),
            "title": str(current.get("title") or ""),
            "cycle": int(current["cycle"]),
            "runtimeSessionId": str(current["runtimeSessionId"]),
            "phase": "READY_TO_RESUME",
        }
    state["lastAutoSuspended"] = {
        "aoneId": current.get("aoneId"),
        "taskId": current.get("taskId"),
        "sessionId": current.get("sessionId"),
        "runtimeSessionId": current.get("runtimeSessionId"),
        "suspendedAt": int(time.time()),
        "reason": "INTERACTIVE_TURN_IDLE",
        "affinityWorkerKey": state.get("workerKey"),
        "affinityExpireAt": pending.get("affinityExpireAt"),
    }
    state["current"] = None
    state.pop("sessionPermit", None)
    state["pendingOperation"] = None
    state["pendingSuspend"] = None


def _send_pending_suspend_locked(state: Dict[str, Any],
                                 cp: ControlPlaneClient) -> bool:
    pending = state.get("pendingSuspend")
    if not isinstance(pending, Mapping):
        return False
    current = state.get("current")
    if (not isinstance(current, Mapping)
            or str(current.get("sessionId")) != str(pending.get("sessionId"))
            or str(current.get("fenceToken")) != str(pending.get("fenceToken"))):
        state["pendingSuspend"] = None
        return False
    cp.suspend_session(
        str(pending["sessionId"]), state["workerKey"], pending["fenceToken"],
        dict(pending["request"]), process_uuid=state["processUuid"],
        request_id=str(pending["requestId"]))
    _finalize_pending_suspend_unlocked(state, pending)
    return True


def _replay_pending_suspend(store: StateStore,
                            cp: ControlPlaneClient,
                            expected_worker_key: str) -> Optional[str]:
    """Determine a persisted suspend before any daemon can renew its old fence."""
    state = store.load()
    if not isinstance(state.get("pendingSuspend"), Mapping):
        return None
    pending = dict(state["pendingSuspend"])
    try:
        with store.locked():
            latest = store.load_unlocked()
            if latest.get("workerKey") != expected_worker_key:
                return "Jarvis Worker incarnation changed while resolving turn suspend."
            _send_pending_suspend_locked(latest, cp)
            store.save_unlocked(latest)
    except StaleFence:
        _clear_lost_current(store, expected_worker_key,
                            pending.get("sessionId"),
                            pending.get("fenceToken"),
                            "pending turn suspend lost its fence")
        return ("Jarvis 上一轮自动挂起的 fence 已失效，已清除旧归属；"
                "处理该单前必须重新运行 claim.sh。")
    except ControlPlaneConflict:
        return ("Jarvis 无法确认上一轮自动挂起状态；本轮必须先通过标准 claim/回执流程恢复，"
                "不得直接执行 Aone 写操作。")
    except ControlPlaneError:
        return ("Jarvis 自动挂起结果暂时无法确认；控制面恢复前不得执行 Aone 或代码写操作。")
    return None


def _prompt_aone_ids(event: Mapping[str, Any]) -> set[str]:
    prompt = str(event.get("prompt") or event.get("user_prompt") or "")
    return set(re.findall(r"(?<!\d)\d{8}(?!\d)", prompt))


def _recovery_pending_claim(old_state: Mapping[str, Any]) -> Optional[Dict[str, Any]]:
    """Convert a dead Codex incarnation into lineage, never reuse its fence."""
    candidates = (
        ("current", old_state.get("current")),
        ("pendingClaim", old_state.get("pendingClaim")),
        ("recoveryPending", old_state.get("recoveryPending")),
    )
    for source, candidate in candidates:
        cycle = (candidate.get("cycle")
                 if isinstance(candidate, Mapping) else None)
        if cycle is None and source == "recoveryPending":
            cycle = old_state.get("claimCounter")
        try:
            cycle_number = int(cycle) if cycle is not None else 0
        except (TypeError, ValueError):
            cycle_number = 0
        if (isinstance(candidate, Mapping)
                and candidate.get("aoneId") is not None
                and candidate.get("projectId") is not None
                and cycle_number > 0
                and candidate.get("runtimeSessionId")):
            recovered = {
                "aoneId": str(candidate["aoneId"]),
                "projectId": str(candidate["projectId"]),
                "cycle": cycle_number,
                "runtimeSessionId": str(candidate["runtimeSessionId"]),
                "phase": "READY_TO_RECOVER",
            }
            if candidate.get("operationKey"):
                recovered["operationKey"] = str(candidate["operationKey"])
            if candidate.get("receiptUnknown"):
                recovered["receiptUnknown"] = True
            pending_operation = old_state.get("pendingOperation")
            # Only an AONE_CLAIM receipt slot (no ``kind``) may fold its
            # operationKey/uncertainty into the recovery claim.  A mid-task
            # external-write slot (comment/status/release-tag/finish-tag)
            # belongs to a different operation type+payload: reusing its key
            # for the next AONE_CLAIM begin would be rejected server-side as
            # an operationKey reuse and trap the task in a recovery loop.
            # Those slots become durable orphan records instead (see
            # _orphaned_operation in _build_incarnation_state).
            if (isinstance(pending_operation, Mapping)
                    and not pending_operation.get("kind")
                    and str(pending_operation.get("aoneId") or "") ==
                    recovered["aoneId"]):
                if pending_operation.get("operationKey"):
                    recovered["operationKey"] = str(
                        pending_operation["operationKey"])
                if str(pending_operation.get("status") or "").upper() != "ACKED":
                    recovered["receiptUnknown"] = True
            if (source == "current"
                    and not candidate.get("heartbeatEnabled", True)):
                recovered["receiptUnknown"] = True
            return recovered
    return None


def _orphaned_operation(old_state: Mapping[str, Any],
                        timestamp: int) -> Optional[Dict[str, Any]]:
    """Freeze a dead incarnation's mid-task receipt as an orphan record.

    A replacement process has no session context left to reconcile the
    operation with, and the slot's operationKey must never leak into the
    next AONE_CLAIM (different type+payload → server rejects the key
    reuse).  Persist the coordinates for the production reconciler /
    manual convergence instead; records are carried across every later
    incarnation and never participate in the PreToolUse guard (the
    pendingOperation slot itself is cleared).
    """
    pending = old_state.get("pendingOperation")
    if not isinstance(pending, Mapping) or not pending.get("kind"):
        return None
    return {
        "operationId": pending.get("operationId"),
        "operationKey": pending.get("operationKey"),
        "kind": pending.get("kind"),
        "aoneId": pending.get("aoneId"),
        "status": pending.get("status"),
        "orphanedAt": timestamp,
    }


def _build_incarnation_state(
        old_state: Mapping[str, Any], *,
        client_name: str, session_id: str, host_pid: int,
        host_process_started_at: str, verify_command: bool,
        cwd: str, transcript_path: Any = None, branch: str = "",
        source: Any = None, headless: bool = False,
        headless_policy: Optional[Mapping[str, Any]] = None,
        now: Optional[int] = None) -> Tuple[Dict[str, Any], bool]:
    """Build the one canonical local state shape for a host incarnation.

    A replacement process may inherit task lineage, but never the previous
    process's database fence or in-flight external-operation receipt.  Both a
    native SessionStart and the bridge's pre-exec headless registration use this
    constructor so a later SessionStart is an idempotent refresh instead of a
    second, subtly different recovery transition.
    """
    host = socket.gethostname()
    same_host_process = bool(
        verify_command
        and old_state.get("verifyHostCommand")
        and str(old_state.get("client") or "") == client_name
        and str(old_state.get("clientSessionId") or "") == session_id
        and str(old_state.get("host") or "") == host
        and str(old_state.get("hostPid") or "") == str(host_pid)
        and host_process_started_at
        and str(old_state.get("hostProcessStartedAt") or "") ==
        host_process_started_at)
    # A model switch can replay SessionStart without replacing the Codex host
    # process.  Keep the already-proven boot identity for that exact process so
    # a transient boot probe failure cannot mint another worker incarnation.
    # An explicit override change remains authoritative.
    boot_id = (
        str(old_state["bootId"])
        if same_host_process
        and old_state.get("bootId")
        and not os.environ.get("JARVIS_BOOT_ID", "").strip()
        else _default_boot_id(host))
    process_uuid = hashlib.sha256(
        ("%s|%s|%s|%s|%s" %
         (client_name, session_id, boot_id, host_pid,
          host_process_started_at)).encode()
    ).hexdigest()[:40]
    worker_key = make_worker_key(host, boot_id, process_uuid)
    same_incarnation = old_state.get("workerKey") == worker_key
    recovery_claim = (
        _recovery_pending_claim(old_state)
        if old_state and not same_incarnation else None)
    timestamp = int(time.time()) if now is None else int(now)
    orphan_operations = [
        dict(record)
        for record in (old_state.get("orphanOperations") or [])
        if isinstance(record, Mapping)] if old_state else []
    # Only a reconstructable recovery claim orphans the mid-task slot: with
    # no claim lineage the slot stays put as the fail-closed tombstone below.
    orphaned = (_orphaned_operation(old_state, timestamp)
                if recovery_claim is not None else None)
    if orphaned is not None:
        orphan_operations.append(orphaned)
    is_compact = bool(
        client_name == "codex"
        and str(source or "") == "compact"
        and same_incarnation)
    state: Dict[str, Any] = {
        "schemaVersion": 2,
        "client": client_name,
        "clientSessionId": session_id,
        "workerKey": worker_key,
        "host": host,
        "bootId": boot_id,
        "processUuid": process_uuid,
        "hostPid": host_pid,
        "hostProcessStartedAt": host_process_started_at,
        "verifyHostCommand": bool(verify_command),
        "cwd": str(cwd),
        "transcriptPath": transcript_path,
        "branch": branch,
        "source": source,
        "version": os.environ.get(
            "JARVIS_INTERACTIVE_WORKER_VERSION", "interactive-v1"),
        "claimCounter": int(old_state.get("claimCounter") or 0),
        "current": old_state.get("current") if same_incarnation else None,
        "sessionPermit": (
            old_state.get("sessionPermit") if same_incarnation else None),
        "pendingClaim": (old_state.get("pendingClaim")
                         if same_incarnation else recovery_claim),
        # Deterministic claim rejections are audit history, never an ownership
        # fence.  Preserve them across incarnations without promoting them back
        # into READY_TO_RECOVER.
        "claimBlocks": [dict(record)
                        for record in (old_state.get("claimBlocks") or [])
                        if isinstance(record, Mapping)][-20:],
        "lastClaimBlocked": old_state.get("lastClaimBlocked"),
        # A complete recovery claim absorbs an uncertain AONE_CLAIM slot into
        # receiptUnknown/operationKey and detaches a mid-task slot into an
        # orphanOperations record.  If lineage is too partial to construct
        # that claim, retain the operation as an explicit fail-closed tombstone
        # instead of silently turning the replacement into an idle worker.
        "pendingOperation": (
            old_state.get("pendingOperation")
            if same_incarnation or recovery_claim is None else None),
        # As with an orphan operation, an incomplete suspend is a durable
        # uncertainty boundary.  A reconstructable recovery claim supersedes
        # it; otherwise retain it so the replacement cannot appear idle.
        "pendingSuspend": (
            old_state.get("pendingSuspend")
            if same_incarnation or recovery_claim is None else None),
        "subagentRevision": (int(old_state.get("subagentRevision") or 0)
                             if same_incarnation else 0),
        "subagentRegistry": (old_state.get("subagentRegistry", {})
                             if same_incarnation else {}),
        "subagentEpochs": (old_state.get("subagentEpochs", {})
                           if same_incarnation else {}),
        "subagentLifecycles": (old_state.get("subagentLifecycles", {})
                               if same_incarnation else {}),
        "subagentSpawnPermits": (old_state.get("subagentSpawnPermits", {})
                                 if same_incarnation else {}),
        "subagentInteractionPermits": (
            old_state.get("subagentInteractionPermits", {})
            if same_incarnation else {}),
        "stopped": not verify_command,
        "turnActive": (bool(old_state.get("turnActive", True))
                       if is_compact else True),
        "activeTurnId": (old_state.get("activeTurnId")
                         if is_compact else None),
        "lastTurnActivityAt": timestamp,
        "registeredAt": timestamp,
    }
    inherited_headless = bool(
        same_incarnation and old_state.get("headlessRegistered"))
    effective_headless_policy = (
        dict(headless_policy)
        if isinstance(headless_policy, Mapping)
        else (dict(old_state["headlessPolicy"])
              if inherited_headless
              and isinstance(old_state.get("headlessPolicy"), Mapping)
              else None))
    if headless or inherited_headless:
        state["headlessRegistered"] = True
    if effective_headless_policy is not None:
        state["headlessPolicy"] = effective_headless_policy
    if is_compact and old_state.get("turnStoppedAt"):
        state["turnStoppedAt"] = old_state.get("turnStoppedAt")
    if ((same_incarnation or recovery_claim is None)
            and old_state.get("lastAutoSuspended")):
        state["lastAutoSuspended"] = old_state.get("lastAutoSuspended")
    # Losing ownership is itself durable fail-closed lineage.  It must survive
    # every replacement incarnation until an exact standard claim reconciles it.
    if old_state.get("lostOwnership"):
        state["lostOwnership"] = old_state.get("lostOwnership")
    # Orphaned mid-task receipts are durable convergence pointers (reconcile
    # by operationId out of band); they survive every incarnation but never
    # gate tools — the pendingOperation slot they came from is cleared.
    if orphan_operations:
        state["orphanOperations"] = orphan_operations
    if same_incarnation and old_state.get("recoveryPending"):
        state["recoveryPending"] = old_state.get("recoveryPending")
    elif recovery_claim:
        state["recoveryPending"] = {
            "aoneId": recovery_claim["aoneId"],
            "projectId": recovery_claim["projectId"],
            "cycle": recovery_claim["cycle"],
            "runtimeSessionId": recovery_claim["runtimeSessionId"],
            "oldWorkerKey": old_state.get("workerKey"),
        }
    if not verify_command:
        state["stoppedAt"] = timestamp
        state["lastError"] = (
            "headless host pid could not be verified"
            if headless else "Codex/Claude host process could not be verified")
    if same_incarnation and old_state.get("daemonPid"):
        state["daemonPid"] = old_state.get("daemonPid")
        state["daemonStartedAt"] = old_state.get("daemonStartedAt")
        state["daemonProcessStartedAt"] = old_state.get(
            "daemonProcessStartedAt")
        state["sidecarHeartbeatAt"] = old_state.get("sidecarHeartbeatAt")
    return state, same_incarnation


def _resume_codex_turn(store: StateStore,
                       event: Mapping[str, Any]) -> Optional[str]:
    """Restart the sidecar and validate its local proof for any carried task."""
    state = store.load()
    if not state:
        return "Jarvis Worker 尚未注册；本轮处理 Aone 前必须先恢复 SessionStart 并走 claim.sh。"
    if not _host_alive(state):
        return "Jarvis Worker 宿主进程校验失败；本轮不得绕过 claim.sh 直接处理 Aone。"

    cp: Optional[ControlPlaneClient] = None
    if state.get("stopped"):
        if state.get("offlineReason") != "idle_ttl":
            return "Jarvis Worker 已离线；请恢复会话注册后再通过 claim.sh 接单。"
        cp = _client()
        try:
            _retry_unavailable(lambda: _register(cp, state, "ACTIVE"))
        except ControlPlaneError:
            return "Jarvis Worker 重新注册失败；控制面恢复前不得直接处理 Aone。"
        with store.locked():
            latest = store.load_unlocked()
            if latest.get("workerKey") == state.get("workerKey"):
                latest["stopped"] = False
                latest.pop("stoppedAt", None)
                latest.pop("offlineReason", None)
                store.save_unlocked(latest)
        state = store.load()

    cp = cp or _client()
    suspend_message = _replay_pending_suspend(
        store, cp, str(state["workerKey"]))
    if suspend_message:
        return suspend_message

    state = store.load()
    lost_ownership = state.get("lostOwnership")
    if isinstance(lost_ownership, Mapping):
        requested_ids = _prompt_aone_ids(event)
        lost_aone = str(lost_ownership.get("aoneId") or "")
        _ensure_daemon(store, str(state["workerKey"]))
        if requested_ids and lost_aone not in requested_ids:
            return ("Jarvis 已失去上一单 %s 的数据库 fence，不能继续旧单；"
                    "本轮新单必须先走对应 claim.sh。" % (lost_aone or "unknown"))
        return ("Jarvis 已失去上一单 %s 的数据库 fence；在重新运行该单 claim.sh 成功前，"
                "严禁继续代码或 Aone 写操作。" % (lost_aone or "unknown"))

    pending_claim = state.get("pendingClaim")
    last_suspended = state.get("lastAutoSuspended")
    recovery_pending = state.get("recoveryPending")
    pending_matches_suspend = (
        isinstance(last_suspended, Mapping)
        and str(pending_claim.get("aoneId")) == str(last_suspended.get("aoneId"))
        and str(pending_claim.get("runtimeSessionId")) ==
        str(last_suspended.get("runtimeSessionId"))) if isinstance(
            pending_claim, Mapping) else False
    pending_matches_recovery = (
        isinstance(recovery_pending, Mapping)
        and str(pending_claim.get("aoneId")) == str(recovery_pending.get("aoneId"))
        and str(pending_claim.get("runtimeSessionId")) ==
        str(recovery_pending.get("runtimeSessionId"))) if isinstance(
            pending_claim, Mapping) else False
    auto_resumable = (
        not isinstance(state.get("current"), Mapping)
        and isinstance(pending_claim, Mapping)
        and pending_claim.get("phase") in (
            "READY_TO_RESUME", "READY_TO_RECOVER", "CLAIMING")
        and (pending_matches_suspend or pending_matches_recovery))
    requested_ids = _prompt_aone_ids(event)
    if (auto_resumable and requested_ids
            and str(pending_claim.get("aoneId")) not in requested_ids):
        _ensure_daemon(store, str(state["workerKey"]))
        return ("Jarvis 上一单 %s 已安全挂起；检测到本轮指向其它 Aone，未自动抢回旧单。"
                "处理新单前必须运行对应的 claim.sh。" % pending_claim.get("aoneId"))
    if auto_resumable:
        try:
            resumed = prepare_claim(
                str(pending_claim["aoneId"]), str(pending_claim["projectId"]),
                str(pending_claim.get("title") or ""),
                str(pending_claim.get("sourceStatus") or ""))
        except ControlPlaneConflict:
            with store.locked():
                latest = store.load_unlocked()
                latest_pending = latest.get("pendingClaim")
                if (isinstance(latest_pending, Mapping)
                        and str(latest_pending.get("aoneId")) ==
                        str(pending_claim.get("aoneId"))
                        and str(latest_pending.get("runtimeSessionId")) ==
                        str(pending_claim.get("runtimeSessionId"))):
                    latest_pending.pop("claimRequestId", None)
                    latest_pending["phase"] = (
                        "READY_TO_RECOVER" if pending_matches_recovery
                        else "READY_TO_RESUME")
                    store.save_unlocked(latest)
            _ensure_daemon(store, str(state["workerKey"]))
            return ("Jarvis 上一轮任务仍由旧 Worker fence 持有或已被其它 Worker 接管；"
                    "本轮不得继续该单，待租约回收后可重试 claim.sh。")
        except (ControlPlaneError, RuntimeError, ValueError, KeyError):
            _ensure_daemon(store, str(state["workerKey"]))
            return ("Jarvis 未能重新取得上一轮任务 fence；控制面恢复并重新通过 claim.sh 前，"
                    "不得继续该单。")
        if (bool(resumed.get("proceed"))
                or str(resumed.get("operationStatus") or "").upper() != "ACKED"):
            _ensure_daemon(store, str(state["workerKey"]))
            return ("Jarvis 已重新取得数据库 fence，但 Aone 接单回执尚未 ACK；"
                    "必须先运行同一单的 claim.sh 完成标准外部回执，再继续工作。")
        with store.locked():
            latest = store.load_unlocked()
            latest.pop("lastAutoSuspended", None)
            latest.pop("recoveryPending", None)
            store.save_unlocked(latest)
        state = store.load()

    _ensure_daemon(store, str(state["workerKey"]))
    state = store.load()
    current = state.get("current")
    if not isinstance(current, Mapping):
        if isinstance(state.get("pendingClaim"), Mapping):
            pending = state["pendingClaim"]
            if pending.get("receiptUnknown"):
                return ("Jarvis 接单回执处于 UNKNOWN，必须先完成 operation receipt 对账；"
                        "不得直接继续任务或重新写 Aone。")
            return ("Jarvis 存在未完成的接单意图（%s）；必须先重试对应 claim.sh，"
                    "不得绕过数据库 fence 继续工作。" %
                    str(pending.get("phase") or "UNKNOWN"))
        if state.get("pendingOperation") or state.get("pendingSuspend"):
            return ("Jarvis 存在未确定化的控制面操作；完成回执/挂起恢复前不得继续任务。")
        if state.get("lastAutoSuspended") or state.get("recoveryPending"):
            return ("Jarvis 存在不完整的任务恢复标记；重新通过 claim.sh 确认数据库 fence 前"
                    "不得继续任务。")
        return None
    if not current.get("heartbeatEnabled", True):
        return ("Jarvis 当前任务的外部操作回执尚未确认；本轮必须先通过标准 claim/回执流程恢复，"
                "不得直接执行 Aone 写操作。")
    permit_reason = _session_permit_block_reason(state, current)
    if permit_reason:
        return permit_reason
    return None


def hook(client_name: str, event: Mapping[str, Any]) -> int:
    event_name = str(event.get("hook_event_name") or "")
    session_id = _nonblank(event.get("session_id"), "session_id")
    state_id = (str(event.get("agent_id") or "").strip()
                if client_name == "codex" and event.get("agent_id")
                else session_id)
    store = StateStore(_state_path(client_name, state_id))

    if client_name == "codex" and event_name == "SubagentStart":
        message = _bind_codex_subagent(store, event)
        _hook_output(
            event_name,
            ("Jarvis Codex subagent 已绑定到父任务的 root-turn/task fence。"
             if not message else
             "Jarvis Codex subagent 绑定尚未完成；任何工具都会 fail-closed，"
             "直到父 spawn 回执可验证：%s" % message))
        return 0

    if client_name == "codex" and event_name == "SubagentStop":
        _stop_codex_subagent(store, event)
        _hook_output(event_name)
        return 0

    if event_name == "PreToolUse":
        # Both Codex and Claude must stop the concrete tool, not merely warn
        # the model, once local ownership becomes uncertain.  Exit 2 + stderr
        # is the blocking contract understood by both clients.
        try:
            if client_name == "codex" and not event.get("agent_id"):
                _record_codex_turn(store, event_name, event)
            message = _guard_pre_tool_use(store, client_name, event)
        except Exception as exc:
            print("interactive worker tool-fence error: %s" % type(exc).__name__,
                  file=sys.stderr)
            message = ("Jarvis Worker 无法完成工具前 fence 校验；当前工具调用已阻断，"
                       "控制面恢复后请重新走 claim.sh。")
        if message:
            print(message, file=sys.stderr)
            return HOOK_BLOCK_EXIT
        _hook_output(event_name)
        return 0

    if client_name == "codex" and event_name in (
            "UserPromptSubmit", "Stop", "PostToolUse"):
        # These hooks bracket every Codex turn. Stop is local and ordered behind
        # wrap-check; prompt start also verifies any carried fence synchronously.
        try:
            if event_name == "PostToolUse":
                _record_codex_post_tool(store, event)
            elif event.get("agent_id"):
                _hook_output(event_name)
                return 0
            else:
                _record_codex_turn(store, event_name, event)
            message = (_resume_codex_turn(store, event)
                       if event_name == "UserPromptSubmit" else None)
        except Exception as exc:
            print("interactive worker turn warning: %s" % type(exc).__name__,
                  file=sys.stderr)
            message = ("Jarvis Worker 转轮校验失败；本轮处理 Aone 前必须重新走 claim.sh，"
                       "不得沿用旧任务归属。")
        _hook_output(event_name, message)
        return STATE_ERROR_EXIT if event_name == "Stop" and message else 0

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
                offline(store, _client(), expected_worker_key)
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

    old_state: Dict[str, Any] = {}
    state: Dict[str, Any] = {}
    same_incarnation = False
    verify_command = False
    worker_key = ""
    cp: Optional[ControlPlaneClient] = None
    registration_error: Optional[BaseException] = None
    # The registration RPC and both success/failure publication happen under
    # one state lock. A failed older SessionStart can therefore never write a
    # tombstone over a newer successful registration.
    with store.locked():
        old_state = store.load_unlocked()
        try:
            if client_name == "claude":
                _persist_claude_context(session_id)
            host_pid, verify_command_value = _find_host_pid(client_name)
            host_process_started_at = (
                _process_start_identity(host_pid) if verify_command_value else "")
            verify_command = bool(
                verify_command_value and host_process_started_at)
            branch = ""
            try:
                branch = subprocess.run(
                    ["git", "-C", str(event.get("cwd") or REPO_ROOT),
                     "rev-parse", "--abbrev-ref", "HEAD"],
                    capture_output=True, text=True, timeout=2,
                    check=False).stdout.strip()
            except (OSError, subprocess.SubprocessError):
                pass
            state, same_incarnation = _build_incarnation_state(
                old_state,
                client_name=client_name,
                session_id=session_id,
                host_pid=host_pid,
                host_process_started_at=host_process_started_at,
                verify_command=verify_command,
                cwd=str(event.get("cwd") or os.getcwd()),
                transcript_path=event.get("transcript_path"),
                branch=branch,
                source=event.get("source"))
            worker_key = str(state["workerKey"])
            cp = _client()
            _retry_unavailable(lambda: _register(
                cp, state, "ACTIVE" if verify_command else "OFFLINE"))
            store.save_unlocked(state)
        except Exception as exc:
            registration_error = exc
            # A same-incarnation refresh failure must not destroy an already
            # published ACTIVE Worker. Tool-level database heartbeats still
            # fail closed while the control plane is unavailable.
            if (same_incarnation and old_state.get("registeredAt")
                    and not old_state.get("registrationFailed")):
                failed_state = dict(old_state)
                failed_state["lastRegistrationAttemptAt"] = int(time.time())
                failed_state["lastRegistrationError"] = type(exc).__name__
            else:
                failed_state = dict(state or old_state)
                failed_state.update({
                    "schemaVersion": 1,
                    "client": client_name,
                    "clientSessionId": session_id,
                    "cwd": str(event.get("cwd") or os.getcwd()),
                    "transcriptPath": event.get("transcript_path"),
                    "stopped": True,
                    "stoppedAt": int(time.time()),
                    "registrationFailed": True,
                    "lastError": "interactive worker registration failed: %s" %
                                 type(exc).__name__,
                })
            store.save_unlocked(failed_state)

    if registration_error is not None:
        print("interactive worker registration warning: %s" %
              type(registration_error).__name__,
              file=sys.stderr)
        _hook_output(
            event_name,
            "警告：Jarvis interactive worker 注册失败。不要直接修改 Aone；"
            "bootstrap/claim.sh 会 fail-closed。请确认 Codex 项目 hook 已信任以及控制面配置可用。")
        return 0

    if old_state and not same_incarnation and cp is not None:
        # Publish the replacement locally before the final remote OFFLINE. Any
        # cleanup failure is post-publication and must not tombstone the winner.
        try:
            _mark_old_offline(cp, old_state)
        except Exception as exc:
            print("interactive worker old-offline warning: %s" %
                  type(exc).__name__, file=sys.stderr)
    daemon_error: Optional[BaseException] = None
    if verify_command:
        try:
            _ensure_daemon(store, worker_key)
        except Exception as exc:
            daemon_error = exc
            print("interactive worker daemon warning: %s" % type(exc).__name__,
                  file=sys.stderr)
        _hook_output(
            event_name,
            "Jarvis interactive worker 已注册（仅定向接单，不拉取公共队列）。"
            "处理 Aone 前必须先运行 bootstrap/claim.sh；控制面异常时接单会 fail-closed。" +
            (" sidecar 启动失败，将在下一轮自动重试。" if daemon_error else ""))
    else:
        _hook_output(
            event_name,
            "警告：无法校验 Codex/Claude 宿主进程，Worker 已按 OFFLINE 注册且不会启动 sidecar。"
            "为避免幽灵续租，bootstrap/claim.sh 将 fail-closed。")
    return 0


def _normalize_headless_policy(
        *, policy_revision: Any = None, aone_write_policy: Any = None,
        headless_kind: Any = None, aone_id: Any = None,
        project_id: Any = None,
        claim_attempt_id: Any = None) -> Optional[Dict[str, str]]:
    """Validate the bridge-owned policy persisted with a headless incarnation."""
    values = {
        "policyRevision": str(policy_revision or "").strip(),
        "aoneWritePolicy": str(aone_write_policy or "").strip(),
        "kind": str(headless_kind or "").strip(),
        "aoneId": str(aone_id or "").strip(),
        "projectId": str(project_id or "").strip(),
        "claimAttemptId": str(claim_attempt_id or "").strip(),
    }
    if not any(values.values()):
        return None
    if values["policyRevision"] != HEADLESS_POLICY_REVISION:
        raise ValueError("unsupported headless policy revision")
    if values["aoneWritePolicy"] != POST_PR_AONE_WRITE_POLICY:
        raise ValueError("unsupported headless Aone write policy")
    if values["kind"] not in POST_PR_HEADLESS_KINDS:
        raise ValueError("unsupported post-PR headless kind")
    if (not values["aoneId"] or not values["projectId"]
            or not values["claimAttemptId"]):
        raise ValueError("post-PR headless policy requires Aone lineage")
    return values


def _is_post_pr_policy(value: Any) -> bool:
    if not isinstance(value, Mapping):
        return False
    return (
        str(value.get("policyRevision") or "") == HEADLESS_POLICY_REVISION
        and str(value.get("aoneWritePolicy") or "") ==
        POST_PR_AONE_WRITE_POLICY
        and str(value.get("kind") or "") in POST_PR_HEADLESS_KINDS
        and bool(str(value.get("aoneId") or "").strip())
        and bool(str(value.get("projectId") or "").strip())
        and bool(str(value.get("claimAttemptId") or "").strip()))


def _post_pr_exec_lineage(command: str) -> bool:
    """Recognize the fixed manager's policy-bearing argv in a live ancestor."""
    try:
        tokens = shlex.split(str(command or ""))
    except ValueError:
        return False
    manager = str(Path(__file__).resolve())
    try:
        manager_index = tokens.index(manager)
    except ValueError:
        return False
    if tokens[manager_index + 1:manager_index + 2] != ["exec-headless"]:
        return False
    options: Dict[str, str] = {}
    index = manager_index + 2
    while index < len(tokens):
        token = tokens[index]
        if token == "--":
            break
        if token.startswith("--") and index + 1 < len(tokens):
            options[token] = tokens[index + 1]
            index += 2
            continue
        index += 1
    return _is_post_pr_policy({
        "policyRevision": options.get("--policy-revision"),
        "aoneWritePolicy": options.get("--aone-write-policy"),
        "kind": options.get("--headless-kind"),
        "aoneId": options.get("--aone-id"),
        "projectId": options.get("--project-id"),
        "claimAttemptId": options.get("--claim-attempt-id"),
    })


def _calling_ancestors(pid: int, max_depth: int = 64) -> Dict[int, str]:
    ancestors: Dict[int, str] = {}
    current = int(pid)
    for _depth in range(max_depth):
        if current <= 1 or current in ancestors:
            break
        started_at = _process_start_identity(current)
        if started_at:
            ancestors[current] = started_at
        parent, _command = _process_info(current)
        if parent <= 0 or parent == current:
            break
        current = parent
    return ancestors


def post_pr_context_active(calling_pid: int) -> bool:
    """Whether ``calling_pid`` belongs to a live restricted headless lineage.

    The immutable live argv is checked first.  Canonical worker state is the durable
    fallback and control-plane mirror source.  The legacy ``/tmp`` process-group
    marker is intentionally not consulted here.
    """
    ancestors = _calling_ancestors(calling_pid)
    if not ancestors:
        return False
    for ancestor_pid in ancestors:
        if _post_pr_exec_lineage(_process_command(ancestor_pid)):
            return True
    root = _state_root()
    try:
        paths = list(root.glob("*.json"))
    except OSError:
        return False
    for path in paths:
        try:
            state = json.loads(path.read_text(encoding="utf-8"))
        except (OSError, TypeError, ValueError):
            continue
        if (not isinstance(state, Mapping)
                or not state.get("headlessRegistered")
                or not _is_post_pr_policy(state.get("headlessPolicy"))):
            continue
        try:
            host_pid = int(state.get("hostPid") or 0)
        except (TypeError, ValueError):
            continue
        expected_start = str(state.get("hostProcessStartedAt") or "")
        if (host_pid in ancestors and expected_start
                and ancestors[host_pid] == expected_start):
            return True
    return False


def register_headless(session_id: str, host_pid: int,
                      client_name: str = "claude",
                      headless_policy: Optional[Mapping[str, Any]] = None) -> Dict[str, Any]:
    """Pre-register one trusted wrapper incarnation before it execs Claude.

    SessionEnd keeps a stopped tombstone rather than deleting the state file.
    A retry/resume must carry that tombstone's task lineage into
    READY_TO_RECOVER, while dropping the dead process's database fence and
    uncertain external-operation state.

    The local atomic write is authoritative for the tool fence and always occurs
    before any network call.  Remote worker registration is deliberately a
    short-timeout, best-effort hint outside the state lock; native SessionStart
    performs the authoritative retried registration after Claude starts.
    """
    session_id = _nonblank(session_id, "session_id")
    host_pid = int(host_pid)
    if host_pid <= 0:
        raise ValueError("host_pid must be a positive integer")
    if client_name not in ("claude", "codex"):
        raise ValueError("client must be claude or codex")
    if headless_policy is not None:
        headless_policy = _normalize_headless_policy(
            policy_revision=headless_policy.get("policyRevision"),
            aone_write_policy=headless_policy.get("aoneWritePolicy"),
            headless_kind=headless_policy.get("kind"),
            aone_id=headless_policy.get("aoneId"),
            project_id=headless_policy.get("projectId"),
            claim_attempt_id=headless_policy.get("claimAttemptId"))
    store = StateStore(_state_path(client_name, session_id))
    old_state: Dict[str, Any] = {}
    with store.locked():
        old_state = store.load_unlocked()
        host_process_started_at = _process_start_identity(host_pid)
        verify_command = bool(_pid_alive(host_pid) and host_process_started_at)
        state, same_incarnation = _build_incarnation_state(
            old_state,
            client_name=client_name,
            session_id=session_id,
            host_pid=host_pid,
            host_process_started_at=host_process_started_at,
            verify_command=verify_command,
            cwd=str(old_state.get("cwd") or os.getcwd()),
            transcript_path=old_state.get("transcriptPath"),
            branch=str(old_state.get("branch") or ""),
            source="headless",
            headless=True,
            headless_policy=headless_policy)
        # Land local lineage FIRST.  No control-plane call may hold this lock.
        store.save_unlocked(state)

    remote_registered = False
    if verify_command:
        try:
            timeout = max(0.05, float(os.environ.get(
                "JARVIS_HEADLESS_REMOTE_REGISTER_TIMEOUT", "0.5")))
            cp = _client(timeout_override=timeout)
            _register(cp, state, "ACTIVE")
            remote_registered = True
            if old_state and not same_incarnation:
                _mark_old_offline(cp, old_state)
        except Exception:
            # Best-effort only.  The local recovery fence remains in force and
            # SessionStart retries the authoritative registration.
            pass
    result = {
        "workerKey": state["workerKey"],
        "hostPid": host_pid,
        "verifyHostCommand": verify_command,
        "sameIncarnation": same_incarnation,
        "remoteRegistered": remote_registered,
        "statePath": str(store.path),
    }
    return result


def exec_headless(session_id: str, command: list[str],
                  client_name: str = "claude",
                  headless_policy: Optional[Mapping[str, Any]] = None) -> None:
    """Atomically publish the fence state, then replace this process with Claude.

    ``exec`` preserves PID and process start time, so the pre-registered
    incarnation exactly matches the later SessionStart and every PreToolUse.
    """
    if not command:
        raise ValueError("headless command must not be empty")
    result = register_headless(
        session_id, os.getpid(), client_name=client_name,
        headless_policy=headless_policy)
    if not result.get("verifyHostCommand"):
        raise RuntimeError("headless wrapper process could not be verified")
    env = os.environ.copy()
    env["JARVIS_INTERACTIVE_CLIENT"] = client_name
    env["JARVIS_INTERACTIVE_SESSION_ID"] = session_id
    os.execvpe(command[0], list(command), env)


def offline(store: StateStore, client: Optional[ControlPlaneClient] = None,
            expected_worker_key: Optional[str] = None,
            reason: str = "host_stopped") -> bool:
    cp = client or _client()
    with store.locked():
        state = store.load_unlocked()
        if not state:
            return False
        if expected_worker_key and state.get("workerKey") != expected_worker_key:
            return False
        # The daemon decides from an unlocked snapshot. Recheck under the same
        # lock used by prompt/claim so a new turn cannot be offlined by a stale
        # TTL decision. Keep the remote OFFLINE heartbeat ordered before a
        # waiting prompt can re-register ACTIVE.
        if reason == "idle_ttl" and not _worker_idle_expired(state):
            return False
        state["stopped"] = True
        state["stoppedAt"] = int(time.time())
        state["offlineReason"] = reason
        if state.get("daemonPid") == os.getpid():
            state.pop("daemonPid", None)
            state.pop("daemonStartedAt", None)
        store.save_unlocked(state)
        try:
            _heartbeat_worker(cp, state, "OFFLINE")
        except ControlPlaneError:
            # The server reaper marks stale workers OFFLINE after heartbeats cease.
            pass
    return True


def _auto_suspend_idle_session(store: StateStore,
                               cp: ControlPlaneClient,
                               expected_worker_key: str) -> bool:
    """Suspend a settled Codex task after Stop without counting it as a crash."""
    with store.locked():
        state = store.load_unlocked()
        current = state.get("current")
        if (state.get("workerKey") != expected_worker_key
                or state.get("stopped")
                or not isinstance(current, Mapping)
                or not current.get("heartbeatEnabled", True)
                or state.get("pendingOperation")
                or _session_heartbeat_allowed(state)):
            return False
        if not isinstance(state.get("pendingSuspend"), Mapping):
            stopped_at = str(state.get("turnStoppedAt") or "unknown")
            affinity_expire_at = time.time() + _interactive_affinity_seconds()
            state["pendingSuspend"] = {
                "sessionId": current["sessionId"],
                "fenceToken": current["fenceToken"],
                "affinityExpireAt": affinity_expire_at,
                "request": {
                    "waitType": "INTERACTIVE_TURN_IDLE",
                    "waitKey": "codex-turn:%s" % hashlib.sha256(
                        (str(state.get("clientSessionId")) + "|" + stopped_at).encode()
                    ).hexdigest()[:24],
                    "waitExpireAt": _utc_timestamp(affinity_expire_at),
                    "transcriptUri": state.get("transcriptPath"),
                    "branchRef": state.get("branch"),
                    "logUri": str(store.path.with_suffix(".log")),
                },
                "requestId": "jarvis-interactive-turn-suspend-%s" %
                hashlib.sha256((str(current["sessionId"]) + "|" +
                                str(current["fenceToken"])).encode()).hexdigest()[:24],
            }
            # Persist the intent before the request. A lost response can then be
            # replayed safely before the next turn is allowed to reuse the task.
            store.save_unlocked(state)
        _send_pending_suspend_locked(state, cp)
        store.save_unlocked(state)
        return True


def daemon(state_path: Path, expected_worker_key: str) -> int:
    store = StateStore(state_path)
    interval = _interactive_heartbeat_seconds()
    cp = _client()
    while True:
        state = store.load()
        if not state or state.get("workerKey") != expected_worker_key:
            return 0
        if state.get("stopped") or not _host_alive(state):
            offline(store, cp, expected_worker_key, "host_stopped")
            return 0
        if isinstance(state.get("pendingSuspend"), Mapping):
            message = _replay_pending_suspend(store, cp, expected_worker_key)
            if message:
                print("interactive suspend replay pending", file=sys.stderr)
            time.sleep(interval)
            continue
        if _worker_idle_expired(state):
            if offline(store, cp, expected_worker_key, "idle_ttl"):
                return 0
            continue
        heartbeat_current: Optional[Mapping[str, Any]] = None
        should_suspend = False
        try:
            # Keep local incarnation replacement and remote heartbeats ordered.
            # A new SessionStart cannot publish a replacement worker in between
            # this recheck and an old daemon's final ACTIVE/session heartbeat.
            with store.locked():
                latest = store.load_unlocked()
                if (latest.get("workerKey") != expected_worker_key
                        or latest.get("stopped")):
                    return 0
                latest["daemonPid"] = os.getpid()
                daemon_started = _process_start_identity(os.getpid())
                if daemon_started:
                    latest["daemonProcessStartedAt"] = daemon_started
                # This is local sidecar liveness, not proof that the remote
                # heartbeat succeeded. Persist it before the request so a
                # control-plane 503 does not make PreToolUse mistake a live
                # sidecar for a dead one; the independently expiring session
                # permit still fences new tools at the lease safety boundary.
                heartbeat_at = time.time()
                latest["sidecarHeartbeatAt"] = heartbeat_at
                store.save_unlocked(latest)
                try:
                    _heartbeat_worker(cp, latest, "ACTIVE")
                except ControlPlaneError as exc:
                    if not _is_explicit_admin_cleanup_not_found(exc):
                        raise
                    _stop_admin_cleaned_sidecar_unlocked(
                        latest, expected_worker_key, "Worker", exc)
                    store.save_unlocked(latest)
                    print("interactive heartbeat stopped: admin cleanup removed Worker",
                          file=sys.stderr)
                    return 0
                current = latest.get("current")
                if (isinstance(current, Mapping)
                        and current.get("heartbeatEnabled", True)):
                    heartbeat_current = dict(current)
                    if _session_heartbeat_allowed(latest):
                        try:
                            # The fenced Session heartbeat is also the durable
                            # checkpoint.  Keep recovery references in the
                            # control plane instead of recreating a
                            # machine-local task ledger.
                            checkpoint = {
                                "leaseSeconds": int(
                                    current.get("leaseSeconds")
                                    or _interactive_lease_seconds()),
                                "transcriptUri": latest.get("transcriptPath"),
                                "workspaceRef": latest.get("cwd"),
                                "branchRef": latest.get("branch"),
                                "logUri": str(store.path.with_suffix(".log")),
                            }
                            response = cp.heartbeat_session(
                                str(current["sessionId"]), latest["workerKey"],
                                current["fenceToken"], checkpoint,
                                process_uuid=latest["processUuid"],
                                request_id="jarvis-interactive-session-heartbeat-%s" %
                                hashlib.sha256((str(current["sessionId"]) +
                                                str(time.time_ns())).encode()).hexdigest()[:24])
                        except ControlPlaneError as exc:
                            if not _is_explicit_admin_cleanup_not_found(exc):
                                raise
                            _stop_admin_cleaned_sidecar_unlocked(
                                latest, expected_worker_key, "Session", exc)
                            store.save_unlocked(latest)
                            print(
                                "interactive heartbeat stopped: "
                                "admin cleanup removed Session",
                                file=sys.stderr)
                            return 0
                        _refresh_session_permit_locked(
                            latest, current, response,
                            source="sidecar-heartbeat", now=heartbeat_at)
                    else:
                        should_suspend = True
                store.save_unlocked(latest)
            if should_suspend:
                _auto_suspend_idle_session(store, cp, expected_worker_key)
        except StaleFence:
            _clear_lost_current(
                store, expected_worker_key,
                (heartbeat_current or {}).get("sessionId"),
                (heartbeat_current or {}).get("fenceToken"),
                "session ownership lost")
        except ControlPlaneConflict as exc:
            print("interactive suspend conflict: %s" % type(exc).__name__, file=sys.stderr)
        except ControlPlaneError as exc:
            # Fail closed for mutations; heartbeat transport failures simply retry.
            print("interactive heartbeat warning: %s" % type(exc).__name__, file=sys.stderr)
        time.sleep(interval)


def _claim_cycle(state: Dict[str, Any], aone_id: str, project_id: str) -> Tuple[int, str]:
    pending = state.get("pendingClaim")
    current = state.get("current")
    recovery = state.get("recoveryPending")
    if (isinstance(pending, Mapping)
            and str(pending.get("aoneId")) == aone_id
            and str(pending.get("projectId")) == project_id):
        return int(pending["cycle"]), str(pending["runtimeSessionId"])
    if (isinstance(current, Mapping)
            and str(current.get("aoneId")) == aone_id
            and str(current.get("projectId")) == project_id):
        return int(current["cycle"]), str(current["runtimeSessionId"])
    if (isinstance(recovery, Mapping)
            and str(recovery.get("aoneId")) == aone_id
            and str(recovery.get("projectId")) == project_id
            and recovery.get("cycle") is not None
            and recovery.get("runtimeSessionId")):
        return int(recovery["cycle"]), str(recovery["runtimeSessionId"])
    cycle = int(state.get("claimCounter") or 0) + 1
    session_hash = hashlib.sha256(
        str(state["clientSessionId"]).encode("utf-8")).hexdigest()[:16]
    runtime_id = "interactive:%s:%s:aone:%s:%s:cycle:%d" % (
        state["client"], session_hash, project_id, aone_id, cycle)
    return cycle, runtime_id[:191]


def _claim_block_value(response: Mapping[str, Any], *keys: str) -> Any:
    """Read one allow-listed diagnostic field from a conflict response."""
    containers = [response]
    data = response.get("data") if isinstance(response, Mapping) else None
    if isinstance(data, Mapping):
        containers.append(data)
    for container in containers:
        for key in keys:
            value = container.get(key)
            if value is not None and value != "":
                return value
    return None


def _record_deterministic_claim_block(
        store: StateStore, expected_worker_key: str,
        aone_id: str, project_id: str, claim_request_id: str,
        error: ControlPlaneConflict) -> None:
    """Move an acknowledged 409 out of the global fail-closed claim slot.

    A transport failure must retain ``pendingClaim`` because the server may have
    accepted the request.  A received HTTP 409 is different: the claim was
    definitively rejected, so keeping ``CLAIMING`` would make the only allowed
    recovery command repeat the same impossible request forever.  Persist a
    bounded audit record using a request-id CAS, then release the blocking slot.
    """
    response = error.response if isinstance(error.response, Mapping) else {}
    with store.locked():
        latest = store.load_unlocked()
        pending = latest.get("pendingClaim")
        if (latest.get("workerKey") != expected_worker_key
                or not isinstance(pending, Mapping)
                or str(pending.get("aoneId")) != aone_id
                or str(pending.get("projectId")) != project_id
                or str(pending.get("claimRequestId")) != claim_request_id):
            return
        block: Dict[str, Any] = {
            "aoneId": aone_id,
            "projectId": project_id,
            "phase": "CLAIM_BLOCKED",
            "claimRequestId": claim_request_id,
            "status": error.status,
            "code": str(error.code or "Conflict.State")[:128],
            "message": str(error)[:500],
            "blockedAt": int(time.time()),
        }
        safe_fields = {
            "taskId": ("taskId", "task_id"),
            "sessionId": ("sessionId", "session_id"),
            "currentState": ("currentState", "taskState", "state"),
            "owner": ("owner", "workerKey", "leaseOwner"),
            "workerHost": ("workerHost", "hostId"),
            "leaseExpireAt": ("leaseExpireAt", "leaseExpiresAt"),
            "lastHeartbeatAt": ("lastHeartbeatAt", "workerHeartbeatAt"),
            "allowedAction": ("allowedAction", "nextAction"),
            "allowedActions": ("allowedActions",),
        }
        for target, keys in safe_fields.items():
            value = _claim_block_value(response, *keys)
            if value is not None:
                block[target] = value
        blocks = [dict(item) for item in (latest.get("claimBlocks") or [])
                  if isinstance(item, Mapping)
                  and str(item.get("claimRequestId")) != claim_request_id]
        blocks.append(block)
        latest["claimBlocks"] = blocks[-20:]
        latest["lastClaimBlocked"] = block
        # A deterministic 409 must not hot-loop on every prompt.  Clear the
        # active slot; recoveryPending retains the exact cycle/runtime lineage
        # so a later explicit claim can resume it after predecessor lease expiry.
        latest["pendingClaim"] = None
        store.save_unlocked(latest)


def _clear_claim_blocks(state: Dict[str, Any], aone_id: str,
                        project_id: str) -> None:
    blocks = [dict(item) for item in (state.get("claimBlocks") or [])
              if isinstance(item, Mapping)
              and not (str(item.get("aoneId")) == aone_id
                       and str(item.get("projectId")) == project_id)]
    state["claimBlocks"] = blocks[-20:]
    last = state.get("lastClaimBlocked")
    if (isinstance(last, Mapping)
            and str(last.get("aoneId")) == aone_id
            and str(last.get("projectId")) == project_id):
        state.pop("lastClaimBlocked", None)


def _matching_current(state: Mapping[str, Any], aone_id: str) -> Mapping[str, Any]:
    current = state.get("current")
    if not isinstance(current, Mapping) or str(current.get("aoneId")) != str(aone_id):
        raise RuntimeError("current interactive session does not own Aone %s" % aone_id)
    return current


def _same_current_assignment(state: Mapping[str, Any], worker_key: Any,
                             expected: Mapping[str, Any]) -> bool:
    current = state.get("current")
    return (state.get("workerKey") == worker_key
            and isinstance(current, Mapping)
            and str(current.get("sessionId")) == str(expected.get("sessionId"))
            and str(current.get("fenceToken")) == str(expected.get("fenceToken")))


def prepare_claim(aone_id: str, project_id: str, title: str = "",
                  source_status: str = "") -> Dict[str, Any]:
    aone_id = _nonblank(aone_id, "aone_id")
    project_id = _nonblank(project_id, "project_id")
    title = str(title or "").strip()
    source_status = str(source_status or "").strip()
    store = _current_store()
    cp = _client()
    with store.locked():
        state = store.load_unlocked()
        if not state or state.get("stopped"):
            raise RuntimeError("interactive worker is not active")
        if not state.get("verifyHostCommand"):
            raise RuntimeError("interactive worker host process is not verified")
        if state.get("pendingSuspend"):
            raise RuntimeError("interactive worker has an unresolved turn suspend")
        current_before_claim = state.get("current")
        if (isinstance(current_before_claim, Mapping)
                and (str(current_before_claim.get("aoneId")) != aone_id
                     or str(current_before_claim.get("projectId")) != project_id)):
            raise ControlPlaneConflict(
                "interactive worker already owns Aone %s" %
                current_before_claim.get("aoneId"))
        existing_claim = state.get("pendingClaim")
        if (isinstance(existing_claim, Mapping)
                and existing_claim.get("phase") == "CLAIMING"
                and (str(existing_claim.get("aoneId")) != aone_id
                     or str(existing_claim.get("projectId")) != project_id)):
            raise ControlPlaneConflict(
                "interactive worker is already claiming Aone %s" %
                existing_claim.get("aoneId"))
        cycle, runtime_id = _claim_cycle(state, aone_id, project_id)
        same_inflight = (
            isinstance(existing_claim, Mapping)
            and str(existing_claim.get("aoneId")) == aone_id
            and str(existing_claim.get("projectId")) == project_id
            and str(existing_claim.get("runtimeSessionId")) == runtime_id
            and existing_claim.get("phase") == "CLAIMING"
            and existing_claim.get("claimRequestId"))
        # Freeze the first observed value (including an unavailable/blank read) for
        # this request identifier. Retrying a lost response must send the same body.
        stable_title = (
            str(existing_claim.get("title") or "")
            if same_inflight and "title" in existing_claim else title)
        # source_status follows the same freeze contract: claim.sh reads the Aone
        # workitem status BEFORE advancing it, so the first non-blank observation is
        # the pre-claim state. A lost-response retry must resend the identical body
        # (and the claim request id is already frozen for the same inflight), so the
        # retry reuses the frozen sourceStatus instead of a freshly read value that
        # may have moved (e.g. status advanced to 处理中 after the first claim).
        stable_source_status = (
            str(existing_claim.get("sourceStatus") or "")
            if same_inflight and "sourceStatus" in existing_claim else source_status)
        claim_request_id = (
            str(existing_claim["claimRequestId"]) if same_inflight else
            "jarvis-interactive-claim-%s" % hashlib.sha256(
                (state["workerKey"] + "|" + runtime_id + "|" +
                 str(time.time_ns())).encode()).hexdigest()[:24])
        state["claimCounter"] = max(int(state.get("claimCounter") or 0), cycle)
        state["pendingClaim"] = {
            "aoneId": aone_id,
            "projectId": project_id,
            "title": stable_title,
            "sourceStatus": stable_source_status,
            "cycle": cycle,
            "runtimeSessionId": runtime_id,
            "phase": "CLAIMING",
            "claimRequestId": claim_request_id,
        }
        if isinstance(existing_claim, Mapping):
            for key in ("operationKey", "receiptUnknown"):
                if key in existing_claim:
                    state["pendingClaim"][key] = existing_claim[key]
        store.save_unlocked(state)

    _retry_unavailable(lambda: _register(cp, state))
    revision = "interactive:%s" % hashlib.sha256(runtime_id.encode()).hexdigest()[:32]
    source_ref = {"aoneId": aone_id, "projectId": project_id}
    if stable_title:
        source_ref["title"] = stable_title
    payload = {
        "kind": "ticket",
        "itemId": aone_id,
        "project": project_id,
        "trigger": "INTERACTIVE",
        "policyRevision": HEADLESS_POLICY_REVISION,
    }
    envelope = TaskEnvelope(
        task_key="aone:%s:%s" % (project_id, aone_id),
        source_type="AONE",
        source_ref=source_ref,
        task_type="ticket",
        desired_revision=policy_desired_revision(revision, payload),
        trigger_mask=["INTERACTIVE"],
        payload=payload,
        recovery_policy="REPLAY_SAFE",
        aone_id=aone_id,
        source_status=stable_source_status,
        required_capabilities={"workerKey": state["workerKey"]},
    )
    lease_seconds = _interactive_lease_seconds()
    # This endpoint is targeted to one exact task. The server checks an active
    # same-runtime assignment before this hint and then enforces real occupied
    # slots transactionally in assignTask. Reporting one lets an expired local
    # receipt-recovery session be resumed without permitting a second task.
    free_slots = 1
    try:
        lease = _retry_unavailable(lambda: cp.claim_task(
            state["workerKey"], envelope, runtime_session_id=runtime_id,
            lease_seconds=lease_seconds, free_slots=free_slots,
            process_uuid=state["processUuid"],
            request_id=claim_request_id))
    except ControlPlaneConflict as exc:
        _record_deterministic_claim_block(
            store, str(state["workerKey"]), aone_id, project_id,
            claim_request_id, exc)
        raise
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
    start_result = _retry_unavailable(lambda: cp.start_session(
        str(session_id), state["workerKey"], fence, start_detail,
        process_uuid=state["processUuid"],
        request_id="jarvis-interactive-session-start-%s" %
        hashlib.sha256((str(session_id) + "|" + str(fence)).encode()).hexdigest()[:24]))

    response_session = _session_response(start_result)
    lease_expire_at = (
        _epoch_seconds(_field(
            response_session, "leaseExpireAt", "leaseExpiresAt", "lease_expire_at"))
        or _epoch_seconds(_field(
            session, "leaseExpireAt", "leaseExpiresAt", "lease_expire_at")))
    current = {
        "aoneId": aone_id,
        "projectId": project_id,
        "title": stable_title,
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
        # The sidecar may only renew after the external-operation receipt is
        # ACKED. A process loss before that point must age out safely.
        "heartbeatEnabled": False,
    }
    if lease_expire_at is not None:
        current["leaseExpireAt"] = lease_expire_at
    with store.locked():
        latest = store.load_unlocked()
        pending_latest = latest.get("pendingClaim")
        if (latest.get("workerKey") != state.get("workerKey")
                or latest.get("stopped")
                or not isinstance(pending_latest, Mapping)
                or str(pending_latest.get("claimRequestId")) != claim_request_id):
            raise ControlPlaneConflict(
                "interactive worker incarnation changed while claiming task")
        latest["current"] = current
        latest["sessionPermit"] = _session_permit(
            latest, current, start_result,
            source="session-start", now=time.time())
        latest["pendingClaim"] = None
        _clear_claim_blocks(latest, aone_id, project_id)
        latest.pop("lostOwnership", None)
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
    if saved_operation_key and not saved_operation_key.startswith("aone-claim:"):
        # Belt-and-suspenders against pre-fix lineage: a mid-task receipt key
        # (comment/status/release-tag/finish-tag) must never be replayed as an
        # AONE_CLAIM key — the server rejects an operationKey reused with a
        # different request; derive a fresh claim key instead.
        saved_operation_key = ""
    operation_key = (saved_operation_key or
                     "aone-claim:%s:%s:%s" % (task_id, generation, cycle))
    operation_request = {
        "taskId": task_id,
        "sessionId": session_id,
        "generation": generation,
        "workerKey": state["workerKey"],
        "processUuid": state["processUuid"],
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
        if not _same_current_assignment(latest, state["workerKey"], current):
            raise ControlPlaneConflict(
                "interactive worker changed before operation receipt")
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
        if (not _same_current_assignment(latest, state["workerKey"], current)
                or not isinstance(pending, Mapping)
                or str(pending.get("operationKey")) != operation_key):
            raise ControlPlaneConflict(
                "interactive worker changed while beginning operation receipt")
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
        pending_latest = latest.get("pendingOperation")
        if (not _same_current_assignment(latest, state["workerKey"], current)
                or not isinstance(pending_latest, Mapping)
                or str(pending_latest.get("operationId")) != str(operation_id)
                or str(pending_latest.get("operationKey")) != operation_key):
            raise ControlPlaneConflict(
                "interactive worker changed while committing operation receipt")
        latest["pendingOperation"] = ({
            "operationId": operation_id,
            "operationKey": operation_key,
            "aoneId": aone_id,
            "proceed": proceed,
            "status": operation_status or "SENDING",
        } if proceed else None)
        if (operation_status == "ACKED"
                and isinstance(latest.get("current"), Mapping)):
            latest["current"]["heartbeatEnabled"] = True
            latest.pop("lastAutoSuspended", None)
            latest.pop("recoveryPending", None)
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
        "processUuid": state["processUuid"],
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
                "processUuid": state["processUuid"],
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
            latest_pending = latest.get("pendingOperation")
            if (_same_current_assignment(latest, state["workerKey"], current)
                    and isinstance(latest_pending, Mapping)
                    and str(latest_pending.get("operationId")) ==
                    str(pending.get("operationId"))):
                latest["current"]["heartbeatEnabled"] = False
                latest["lastError"] = "AONE_CLAIM acknowledgement is ambiguous"
                store.save_unlocked(latest)
        raise
    with store.locked():
        latest = store.load_unlocked()
        latest_pending = latest.get("pendingOperation")
        if (not _same_current_assignment(latest, state["workerKey"], current)
                or not isinstance(latest_pending, Mapping)
                or str(latest_pending.get("operationId")) !=
                str(pending.get("operationId"))):
            raise ControlPlaneConflict(
                "interactive worker changed while acknowledging claim receipt")
        latest["pendingOperation"] = None
        latest["current"]["heartbeatEnabled"] = True
        latest.pop("lastAutoSuspended", None)
        latest.pop("recoveryPending", None)
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
                "processUuid": state["processUuid"],
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
            process_uuid=state["processUuid"],
            request_id="jarvis-interactive-session-fail-%s" %
            hashlib.sha256(str(current["runtimeSessionId"]).encode()).hexdigest()[:24])
    except BaseException as exc:
        if first_error is None:
            first_error = exc
    if first_error is None:
        with store.locked():
            latest = store.load_unlocked()
            if not _same_current_assignment(latest, state["workerKey"], current):
                raise ControlPlaneConflict(
                    "interactive worker changed while failing claim")
            latest["current"] = None
            latest.pop("sessionPermit", None)
            # A failed required receipt remains part of this exact logical
            # claim.  Reuse its cycle/runtime/operation key after RETRY_WAIT;
            # UNKNOWN likewise stays fenced until explicit reconciliation.
            latest["pendingClaim"] = {
                "aoneId": str(current["aoneId"]),
                "projectId": str(current["projectId"]),
                "title": str(current.get("title") or ""),
                "cycle": int(current["cycle"]),
                "runtimeSessionId": str(current["runtimeSessionId"]),
                "phase": "READY_TO_CLAIM",
                "operationKey": (str(pending.get("operationKey"))
                                 if isinstance(pending, Mapping) else ""),
                "receiptUnknown": bool(unknown),
            }
            latest["pendingOperation"] = None
            store.save_unlocked(latest)
    if first_error is not None:
        raise first_error


# Generalized mid-task external-write receipts (docs/aone-operation-receipts.md).
# Unlike prepare-claim these never create/resume a session: they attach a receipt
# to the CURRENT fenced assignment for wrap.sh comment/status writes and
# claim.sh release/finish tag writes.
OPERATION_TYPE_BY_KIND = {
    "comment": "AONE_COMMENT",
    "status": "AONE_STATUS",
    "release-tag": "AONE_RELEASE",
    "finish-tag": "AONE_RELEASE",
}


def _pending_operation(state: Mapping[str, Any]) -> Mapping[str, Any]:
    pending = state.get("pendingOperation")
    if not isinstance(pending, Mapping):
        raise RuntimeError("no pending external operation receipt")
    return pending


def operation_begin(aone_id: str, kind: str, material: str, *,
                    payload_json: Optional[str] = None,
                    required: bool = True,
                    replay_safe: bool = False) -> Dict[str, Any]:
    aone_id = _nonblank(aone_id, "aone_id")
    material = _nonblank(material, "material")
    operation_type = OPERATION_TYPE_BY_KIND.get(kind)
    if operation_type is None:
        raise ValueError("unknown operation kind %s" % kind)
    if payload_json is not None:
        payload = json.loads(payload_json)
        if not isinstance(payload, Mapping):
            raise ValueError("--payload-json must be a JSON object")
        payload = dict(payload)
    else:
        payload = {"kind": kind, "material": material}
    store = _current_store()
    cp = _client()
    state = store.load()
    if state.get("stopped"):
        raise RuntimeError("interactive worker is not active")
    current = _matching_current(state, aone_id)
    task_id = current["taskId"]
    generation = current["generation"]
    # Lease-attempt component: the same task generation can be revived and
    # re-claimed by a new Session/fence after e.g. an ACKED release, and the
    # new attempt's identical write (release-tag idle …) must NOT dedup
    # against the previous attempt's ACKED receipt — that would skip the tag
    # write and strand Aone in the old state.  Same attempt (fence unchanged)
    # retries still reuse the same key; a new Session/fence derives a new one
    # (cf. bridge _PostPrTaskBookend claimAttemptId).
    attempt12 = hashlib.sha256(
        ("%s:%s" % (current["sessionId"], current["fenceToken"]))
        .encode("utf-8")).hexdigest()[:12]
    operation_key = "%s:%s:%s:%s:%s" % (
        kind, task_id, generation, attempt12,
        hashlib.sha256(material.encode("utf-8")).hexdigest()[:12])
    operation_request = {
        "taskId": task_id,
        "sessionId": current["sessionId"],
        "generation": generation,
        "workerKey": state["workerKey"],
        "processUuid": state["processUuid"],
        "fenceToken": current["fenceToken"],
        "operationKey": operation_key,
        "operationType": operation_type,
        "target": aone_id,
        "requestPayload": payload,
        "required": bool(required),
        "maxRetries": 3,
    }
    # Persist intent before sending begin — the same durable-resume proof as the
    # claim receipt: a lost begin response leaves a local record that authorizes
    # a safe retry with the same operationKey.
    unknown_short_circuit: Optional[Dict[str, Any]] = None
    with store.locked():
        latest = store.load_unlocked()
        if not _same_current_assignment(latest, state["workerKey"], current):
            raise ControlPlaneConflict(
                "interactive worker changed before operation receipt")
        existing_pending = latest.get("pendingOperation")
        had_local_intent = isinstance(existing_pending, Mapping)
        if (isinstance(existing_pending, Mapping)
                and str(existing_pending.get("operationKey")) != operation_key):
            # 单槽约束：任一时刻至多一个在途外部写回执。
            raise ControlPlaneConflict(
                "another external operation is pending for this worker")
        if (isinstance(existing_pending, Mapping)
                and str(existing_pending.get("status") or "").upper() == "UNKNOWN"):
            # 上轮 abort --unknown 留下的槽：副作用不可证，绝不盲 begin/重放，
            # 短路到 readback→reconcile 收敛路径（服务端 UNKNOWN begin 也会 409）。
            unknown_short_circuit = {
                "accepted": True,
                "proceed": False,
                "needsReadback": True,
                "operationId": existing_pending.get("operationId"),
                "operationStatus": "UNKNOWN",
            }
        elif not isinstance(existing_pending, Mapping):
            latest["pendingOperation"] = {
                "operationId": None,
                "operationKey": operation_key,
                "aoneId": aone_id,
                "kind": kind,
                "proceed": False,
                "status": "BEGINNING",
            }
            store.save_unlocked(latest)
    if unknown_short_circuit is not None:
        return unknown_short_circuit
    try:
        begun = cp.begin_operation(
            operation_request,
            request_id="jarvis-interactive-operation-begin-%s" %
            hashlib.sha256(operation_key.encode()).hexdigest()[:24])
    except ControlPlaneRejected as exc:
        # 400/401/403 are authoritative HTTP responses: the begin transaction
        # did not run, therefore the Aone side effect cannot have started.
        if exc.status in (400, 401, 403):
            with store.locked():
                latest = store.load_unlocked()
                pending = latest.get("pendingOperation")
                if (isinstance(pending, Mapping)
                        and str(pending.get("operationKey")) == operation_key):
                    latest["pendingOperation"] = None
                    store.save_unlocked(latest)
        raise
    except StaleFence:
        # A 412 proves this process no longer owns the assignment.  Keeping a
        # receipt slot or suggesting another claim with the old fence would
        # turn a deterministic ownership loss into an unrecoverable loop.
        _clear_lost_current(
            store, state["workerKey"], current["sessionId"],
            current["fenceToken"], "stale_fence:operation_begin")
        raise
    except ControlPlaneConflict as exc:
        # A 409 may mean the idempotent begin committed previously.  Preserve
        # the key (and a returned operation id when available) so point-read
        # can decide whether the external write is safe to resume.
        with store.locked():
            latest = store.load_unlocked()
            pending = latest.get("pendingOperation")
            if (isinstance(pending, Mapping)
                    and str(pending.get("operationKey")) == operation_key):
                frozen = dict(pending)
                frozen["status"] = "UNKNOWN"
                frozen["proceed"] = False
                operation_id = _field(
                    exc.response, "operationId", "operation_id", "id")
                if operation_id is not None:
                    frozen["operationId"] = operation_id
                latest["pendingOperation"] = frozen
                store.save_unlocked(latest)
        raise
    except ControlPlaneUnavailable:
        # No response (including an indeterminate 5xx) cannot prove whether the
        # transaction committed. Keep the key and expose point-read recovery.
        with store.locked():
            latest = store.load_unlocked()
            pending = latest.get("pendingOperation")
            if (isinstance(pending, Mapping)
                    and str(pending.get("operationKey")) == operation_key):
                frozen = dict(pending)
                frozen["status"] = "UNKNOWN"
                frozen["proceed"] = False
                latest["pendingOperation"] = frozen
                store.save_unlocked(latest)
        raise
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
        if (not _same_current_assignment(latest, state["workerKey"], current)
                or not isinstance(pending, Mapping)
                or str(pending.get("operationKey")) != operation_key):
            raise ControlPlaneConflict(
                "interactive worker changed while beginning operation receipt")
        latest["pendingOperation"] = {
            "operationId": operation_id,
            "operationKey": operation_key,
            "aoneId": aone_id,
            "kind": kind,
            "proceed": False,
            "status": operation_status or "SENDING",
        }
        store.save_unlocked(latest)
    needs_readback = False
    if operation_status == "ACKED":
        proceed = False
    elif server_proceed:
        proceed = True
    elif operation_status == "SENDING" and local_retry and replay_safe:
        # 幂等写（tag/status）：上一进程已持久化意图，直接重放后 ACK 是安全的。
        proceed = True
    elif operation_status == "SENDING" and local_retry:
        # 不幂等写（comment）：必须 readback 后走 reconcile 收敛，不能盲重放。
        proceed = False
        needs_readback = True
    else:
        # A SENDING receipt without a local intent record can be a lost begin
        # response from another incarnation; the external effect cannot be
        # proven either way, so fail closed.
        raise ControlPlaneConflict(
            "%s receipt is %s without resumable local intent" %
            (operation_type, operation_status or "UNKNOWN"))

    with store.locked():
        latest = store.load_unlocked()
        pending_latest = latest.get("pendingOperation")
        if (not _same_current_assignment(latest, state["workerKey"], current)
                or not isinstance(pending_latest, Mapping)
                or str(pending_latest.get("operationId")) != str(operation_id)
                or str(pending_latest.get("operationKey")) != operation_key):
            raise ControlPlaneConflict(
                "interactive worker changed while committing operation receipt")
        # mid-task 回执不动 heartbeatEnabled：session 本来就在跑（区别于 claim 期
        # 间"ACK 前不许续租"），刻意为之，见 docs/aone-operation-receipts.md。
        latest["pendingOperation"] = (None if operation_status == "ACKED" else {
            "operationId": operation_id,
            "operationKey": operation_key,
            "aoneId": aone_id,
            "kind": kind,
            "proceed": proceed,
            "status": operation_status or "SENDING",
        })
        store.save_unlocked(latest)
    return {
        "accepted": True,
        "proceed": proceed,
        "needsReadback": needs_readback,
        "operationId": operation_id,
        "operationStatus": operation_status,
    }


def operation_abort(aone_id: str, message: str, *,
                    unknown: bool = False) -> Dict[str, Any]:
    """终结当前外部写回执，但不 fail session、不写 pendingClaim（区别于 fail_claim）。"""
    store = _current_store()
    cp = _client()
    state = store.load()
    current = _matching_current(state, aone_id)
    pending = _pending_operation(state)
    error = {"errorType": "AoneOperationFailed", "message": str(message)[:500]}
    if pending.get("operationId") is not None:
        request = {
            "operationId": pending["operationId"],
            "workerKey": state["workerKey"],
            "processUuid": state["processUuid"],
            "fenceToken": current["fenceToken"],
            "error": error,
            "unknown": bool(unknown),
            "retryAllowed": not unknown,
            "retryAfterSeconds": 0,
        }
        if unknown:
            cp.fail_operation(
                request, request_id="jarvis-interactive-operation-abort-%s" %
                hashlib.sha256(str(pending["operationId"]).encode()).hexdigest()[:24])
        else:
            cp.mark_operation_not_started(
                {key: value for key, value in request.items()
                 if key not in ("error", "unknown", "retryAllowed", "retryAfterSeconds")}
                | {"reason": str(message)[:500]},
                request_id="jarvis-interactive-operation-not-started-%s" %
                hashlib.sha256(str(pending["operationId"]).encode()).hexdigest()[:24])
    with store.locked():
        latest = store.load_unlocked()
        latest_pending = latest.get("pendingOperation")
        if (not _same_current_assignment(latest, state["workerKey"], current)
                or not isinstance(latest_pending, Mapping)
                or str(latest_pending.get("operationKey")) !=
                str(pending.get("operationKey"))):
            raise ControlPlaneConflict(
                "interactive worker changed while aborting operation receipt")
        if unknown:
            # 保留本地槽（status=UNKNOWN），供下轮 operation-begin 短路到
            # readback/reconcile；UNKNOWN 永不盲重放。
            frozen = dict(latest_pending)
            frozen["proceed"] = False
            frozen["status"] = "UNKNOWN"
            latest["pendingOperation"] = frozen
        else:
            # 服务端确认副作用未开始；本地 intent 可安全清理，下次 begin 可重试。
            latest["pendingOperation"] = None
        store.save_unlocked(latest)
    return {"aborted": True, "unknown": bool(unknown)}


def operation_reconcile(aone_id: str, *, found: bool,
                        external_ref: Optional[str] = None,
                        retry_allowed: bool = True) -> Dict[str, Any]:
    store = _current_store()
    cp = _client()
    state = store.load()
    current = _matching_current(state, aone_id)
    pending = _pending_operation(state)
    if pending.get("operationId") is None:
        try:
            point = cp.get_operation_by_key(
                current["taskId"], current["generation"], pending["operationKey"])
        except ControlPlaneRejected as exc:
            if exc.status != 404:
                raise
            # Point-read absence is authoritative: begin never committed, hence
            # no external write was authorized. Clear the local tombstone.
            with store.locked():
                latest = store.load_unlocked()
                if (isinstance(latest.get("pendingOperation"), Mapping)
                        and str(latest["pendingOperation"].get("operationKey")) ==
                        str(pending.get("operationKey"))):
                    latest["pendingOperation"] = None
                    store.save_unlocked(latest)
            return {"proceed": False, "found": False, "notStarted": True}
        operation = point.get("operation") if isinstance(point, Mapping) else None
        operation_id = _field(operation, "id", "operationId") if isinstance(operation, Mapping) else None
        if operation_id is None:
            raise RuntimeError("operation point-read returned no receipt")
        pending = dict(pending)
        pending["operationId"] = operation_id
        with store.locked():
            latest = store.load_unlocked()
            latest_pending = latest.get("pendingOperation")
            if (not isinstance(latest_pending, Mapping)
                    or str(latest_pending.get("operationKey")) != str(pending.get("operationKey"))):
                raise ControlPlaneConflict("interactive worker changed during operation point-read")
            latest["pendingOperation"] = pending
            store.save_unlocked(latest)
    request: Dict[str, Any] = {
        "operationId": pending["operationId"],
        "workerKey": state["workerKey"],
        "fenceToken": current["fenceToken"],
        "found": bool(found),
        "retryAllowed": bool(retry_allowed) and not found,
        "retryAfterSeconds": 0,
    }
    if found:
        request["externalRef"] = _nonblank(external_ref, "external_ref")
    _retry_unavailable(lambda: cp.reconcile_operation(
        request, request_id="jarvis-interactive-operation-reconcile-%s" %
        hashlib.sha256(str(pending["operationId"]).encode()).hexdigest()[:24]))
    with store.locked():
        latest = store.load_unlocked()
        latest_pending = latest.get("pendingOperation")
        if (not _same_current_assignment(latest, state["workerKey"], current)
                or not isinstance(latest_pending, Mapping)
                or str(latest_pending.get("operationId")) !=
                str(pending.get("operationId"))):
            raise ControlPlaneConflict(
                "interactive worker changed while reconciling operation receipt")
        if found or not retry_allowed:
            # found：副作用已在，服务端 ACKED，恰好一次收敛；no-retry：服务端 DEAD。
            latest["pendingOperation"] = None
        else:
            # not-found+retry：保留 key（status=RETRY_WAIT），调用方随后重新
            # operation-begin 拿 proceed 重发。
            kept = dict(latest_pending)
            kept["proceed"] = False
            kept["status"] = "RETRY_WAIT"
            latest["pendingOperation"] = kept
        store.save_unlocked(latest)
    result: Dict[str, Any] = {"proceed": False, "found": bool(found)}
    if not found and retry_allowed:
        result["retryScheduled"] = True
    return result


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
            process_uuid=state["processUuid"],
            request_id="jarvis-interactive-session-suspend-%s" %
            hashlib.sha256(str(current["runtimeSessionId"]).encode()).hexdigest()[:24]))
    elif action == "complete":
        payload = {"result": {"aoneId": str(aone_id), "summary": str(detail or "completed")}}
        result = _retry_unavailable(lambda: cp.complete_session(
            str(current["sessionId"]), state["workerKey"], current["fenceToken"], payload,
            process_uuid=state["processUuid"],
            request_id="jarvis-interactive-session-complete-%s" %
            hashlib.sha256(str(current["runtimeSessionId"]).encode()).hexdigest()[:24]))
    else:
        raise ValueError("unknown transition %s" % action)
    with store.locked():
        latest = store.load_unlocked()
        if not _same_current_assignment(latest, state["workerKey"], current):
            raise ControlPlaneConflict(
                "interactive worker changed while committing task transition")
        latest["current"] = None
        latest.pop("sessionPermit", None)
        latest["pendingClaim"] = None
        latest["pendingOperation"] = None
        latest["pendingSuspend"] = None
        latest.pop("lastAutoSuspended", None)
        latest.pop("recoveryPending", None)
        latest.pop("lostOwnership", None)
        store.save_unlocked(latest)
    return dict(result)


def stop_check() -> int:
    """Control-plane-aware session exit gate.

    Returns:
        0 — safe to stop (no active session, or session already terminal)
        2 — stop blocked (active session with unfinished work)
        1 — state unavailable (caller should fall back to wrap-check.sh)
    """
    try:
        store = _current_store()
    except RuntimeError:
        return 1
    try:
        state = store.load()
    except (RuntimeError, OSError):
        return 1

    if not state:
        return 0

    current = state.get("current")
    if not isinstance(current, Mapping):
        return 0

    permit = state.get("sessionPermit") or {}
    session_status = str(permit.get("sessionStatus") or "")
    if session_status in ("COMPLETED", "FAILED", "SUSPENDED"):
        return 0

    if state.get("lostOwnership") or state.get("stopped"):
        return 0

    aone_id = current.get("aoneId") or "?"
    problems: list[str] = []

    if state.get("pendingOperation"):
        problems.append("  %s (外部操作回执未收敛)" % aone_id)

    if state.get("pendingClaim"):
        problems.append("  %s (接单意图未完成)" % aone_id)

    if state.get("pendingSuspend"):
        problems.append("  %s (挂起结果未确定)" % aone_id)

    if problems:
        print("stop-check: 活跃 session 存在未收敛状态:", file=sys.stderr)
        for p in problems:
            print(p, file=sys.stderr)
        return HOOK_BLOCK_EXIT

    print("stop-check: session 仍活跃 (aone=%s, status=%s); "
          "请先 claim.sh release 或 wrap.sh done 收尾" % (aone_id, session_status or "ACTIVE"),
          file=sys.stderr)
    return HOOK_BLOCK_EXIT


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
    hook_parser.add_argument("--expected-event", choices=(
        "SessionStart", "SessionEnd", "UserPromptSubmit", "SubagentStart",
        "SubagentStop", "PreToolUse", "PostToolUse", "Stop"))
    daemon_parser = sub.add_parser("daemon")
    daemon_parser.add_argument("--state", required=True)
    daemon_parser.add_argument("--worker-key", required=True)
    claim_parser = sub.add_parser("prepare-claim")
    claim_parser.add_argument("aone_id")
    claim_parser.add_argument("project_id")
    claim_parser.add_argument("title", nargs="?", default="")
    claim_parser.add_argument("source_status", nargs="?", default="")
    ack_parser = sub.add_parser("operation-ack")
    ack_parser.add_argument("aone_id")
    ack_parser.add_argument("external_ref")
    fail_parser = sub.add_parser("operation-fail")
    fail_parser.add_argument("aone_id")
    fail_parser.add_argument("message")
    fail_parser.add_argument("--unknown", action="store_true")
    begin_parser = sub.add_parser("operation-begin")
    begin_parser.add_argument("aone_id")
    begin_parser.add_argument("kind", choices=sorted(OPERATION_TYPE_BY_KIND))
    begin_parser.add_argument("material")
    begin_parser.add_argument("--payload-json")
    begin_parser.add_argument("--not-required", action="store_true")
    begin_parser.add_argument("--replay-safe", action="store_true")
    abort_parser = sub.add_parser("operation-abort")
    abort_parser.add_argument("aone_id")
    abort_parser.add_argument("message")
    abort_parser.add_argument("--unknown", action="store_true")
    reconcile_parser = sub.add_parser("operation-reconcile")
    reconcile_parser.add_argument("aone_id")
    reconcile_group = reconcile_parser.add_mutually_exclusive_group(required=True)
    reconcile_group.add_argument("--found", metavar="external_ref")
    reconcile_group.add_argument("--not-found", action="store_true")
    reconcile_parser.add_argument("--no-retry", action="store_true")
    current_parser = sub.add_parser("has-current")
    current_parser.add_argument("aone_id")
    suspend_parser = sub.add_parser("suspend")
    suspend_parser.add_argument("aone_id")
    suspend_parser.add_argument("detail", nargs="?", default="released")
    complete_parser = sub.add_parser("complete")
    complete_parser.add_argument("aone_id")
    complete_parser.add_argument("detail", nargs="?", default="completed")
    register_parser = sub.add_parser("register-headless")
    register_parser.add_argument("--session-id", required=True)
    register_parser.add_argument("--pid", required=True, type=int)
    register_parser.add_argument("--client", choices=("claude", "codex"),
                                 default="claude")
    register_parser.add_argument("--policy-revision")
    register_parser.add_argument("--aone-write-policy")
    register_parser.add_argument("--headless-kind")
    register_parser.add_argument("--aone-id")
    register_parser.add_argument("--project-id")
    register_parser.add_argument("--claim-attempt-id")
    exec_parser = sub.add_parser("exec-headless")
    exec_parser.add_argument("--session-id", required=True)
    exec_parser.add_argument("--client", choices=("claude", "codex"),
                             default="claude")
    exec_parser.add_argument("--policy-revision")
    exec_parser.add_argument("--aone-write-policy")
    exec_parser.add_argument("--headless-kind")
    exec_parser.add_argument("--aone-id")
    exec_parser.add_argument("--project-id")
    exec_parser.add_argument("--claim-attempt-id")
    exec_parser.add_argument("headless_command", nargs=argparse.REMAINDER)
    context_parser = sub.add_parser("post-pr-context")
    context_parser.add_argument("--pid", type=int, default=os.getppid())
    sub.add_parser("status")
    sub.add_parser("stop-check")
    return parser


def main(argv: Optional[list[str]] = None) -> int:
    args = _parser().parse_args(argv)
    if args.command == "hook":
        try:
            event = json.load(sys.stdin)
            if not isinstance(event, Mapping):
                raise ValueError("hook input must be an object")
            if (args.expected_event
                    and str(event.get("hook_event_name") or "") !=
                    args.expected_event):
                raise ValueError("hook event does not match expected event")
        except Exception as exc:
            print("interactive hook input error: %s" % type(exc).__name__, file=sys.stderr)
            if args.expected_event in ("PreToolUse", "Stop"):
                print("Jarvis Worker 无法验证工具/停止事件；已 fail-closed。",
                      file=sys.stderr)
                return HOOK_BLOCK_EXIT
            _hook_output("SessionStart", "Jarvis worker hook input 无效；接单将 fail-closed。")
            return 0
        try:
            return hook(args.client, event)
        except Exception as exc:
            print("interactive hook execution error: %s" % type(exc).__name__,
                  file=sys.stderr)
            if args.expected_event in ("PreToolUse", "Stop"):
                print("Jarvis Worker hook 执行失败；已 fail-closed。", file=sys.stderr)
                return HOOK_BLOCK_EXIT
            _hook_output(args.expected_event or "SessionStart",
                         "Jarvis Worker hook 执行失败；接单将 fail-closed。")
            return 0
    try:
        if args.command == "daemon":
            return daemon(Path(args.state), args.worker_key)
        if args.command == "prepare-claim":
            _print_json(prepare_claim(args.aone_id, args.project_id, args.title,
                                      args.source_status))
        elif args.command == "operation-ack":
            _print_json(acknowledge_claim(args.aone_id, args.external_ref))
        elif args.command == "operation-fail":
            fail_claim(args.aone_id, args.message, unknown=args.unknown)
            _print_json({"failed": True})
        elif args.command == "operation-begin":
            _print_json(operation_begin(
                args.aone_id, args.kind, args.material,
                payload_json=args.payload_json,
                required=not args.not_required,
                replay_safe=args.replay_safe))
        elif args.command == "operation-abort":
            _print_json(operation_abort(
                args.aone_id, args.message, unknown=args.unknown))
        elif args.command == "operation-reconcile":
            _print_json(operation_reconcile(
                args.aone_id,
                found=args.found is not None,
                external_ref=args.found,
                retry_allowed=not args.no_retry))
        elif args.command == "has-current":
            return 0 if has_current(args.aone_id) else 1
        elif args.command == "suspend":
            _print_json(transition(args.aone_id, "suspend", args.detail))
        elif args.command == "complete":
            _print_json(transition(args.aone_id, "complete", args.detail))
        elif args.command == "register-headless":
            policy = _normalize_headless_policy(
                policy_revision=args.policy_revision,
                aone_write_policy=args.aone_write_policy,
                headless_kind=args.headless_kind,
                aone_id=args.aone_id,
                project_id=args.project_id,
                claim_attempt_id=args.claim_attempt_id)
            _print_json(register_headless(
                args.session_id, args.pid, client_name=args.client,
                headless_policy=policy))
        elif args.command == "exec-headless":
            command = list(args.headless_command)
            if command[:1] == ["--"]:
                command = command[1:]
            policy = _normalize_headless_policy(
                policy_revision=args.policy_revision,
                aone_write_policy=args.aone_write_policy,
                headless_kind=args.headless_kind,
                aone_id=args.aone_id,
                project_id=args.project_id,
                claim_attempt_id=args.claim_attempt_id)
            exec_headless(
                args.session_id, command, client_name=args.client,
                headless_policy=policy)
        elif args.command == "post-pr-context":
            return 0 if post_pr_context_active(args.pid) else 1
        elif args.command == "status":
            _print_json(worker_status())
        elif args.command == "stop-check":
            return stop_check()
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
