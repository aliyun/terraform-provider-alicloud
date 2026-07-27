"""Bounded subprocess execution with complete process-group cleanup.

Every bridge command that can reach ``bin/a1id`` runs through this module.
``start_new_session=True`` makes the child PID the process-group ID; retaining
that value at spawn time is important because the command leader can exit while
an ``a1`` grandchild remains alive and holds a filesystem lock.
"""

from __future__ import annotations

import errno
import os
import signal
import subprocess
import time
from typing import Any, Iterable, Mapping, Optional, Sequence


DEFAULT_TERM_GRACE_SECONDS = 1.0
DEFAULT_KILL_GRACE_SECONDS = 2.0


def _process_group_members(pgid: int) -> list[tuple[int, str]]:
    """Return live (non-zombie) members of ``pgid``.

    ``killpg(..., 0)`` treats zombies as present even though they can no longer
    hold a lock or execute work.  Reading the process table lets cleanup confirm
    that every runnable descendant is gone without waiting forever for an
    unrelated init process to reap a zombie.
    """

    try:
        snapshot = subprocess.check_output(
            ["ps", "-axo", "pid=,pgid=,stat="],
            text=True,
            stderr=subprocess.DEVNULL,
        )
    except (OSError, subprocess.SubprocessError):
        try:
            os.killpg(int(pgid), 0)
        except ProcessLookupError:
            return []
        except OSError as exc:
            if exc.errno == errno.ESRCH:
                return []
        # The process table is unavailable but the group still exists.  Use a
        # sentinel member so callers fail closed and escalate to SIGKILL.
        return [(-1, "?")]

    members: list[tuple[int, str]] = []
    for line in snapshot.splitlines():
        fields = line.split(None, 2)
        if len(fields) != 3:
            continue
        try:
            pid = int(fields[0])
            member_pgid = int(fields[1])
        except ValueError:
            continue
        state = fields[2].strip()
        if member_pgid == int(pgid) and not state.startswith("Z"):
            members.append((pid, state))
    return members


def _wait_group_empty(pgid: int, timeout: float) -> bool:
    deadline = time.monotonic() + max(0.0, float(timeout))
    while True:
        if not _process_group_members(pgid):
            return True
        remaining = deadline - time.monotonic()
        if remaining <= 0:
            return False
        time.sleep(min(0.05, remaining))


def _signal_process_group(pgid: int, signum: int) -> None:
    """Signal a group, falling back to its captured members on EPERM.

    Some constrained launch environments permit signaling owned child PIDs but
    reject a group-wide signal.  Production still takes the atomic ``killpg``
    path; the fallback preserves cleanup in those constrained environments.
    """

    try:
        os.killpg(int(pgid), signum)
        return
    except ProcessLookupError:
        return
    except OSError as exc:
        if exc.errno == errno.ESRCH:
            return
        if exc.errno != errno.EPERM:
            raise
    for pid, _state in _process_group_members(pgid):
        if pid <= 0:
            continue
        try:
            os.kill(pid, signum)
        except ProcessLookupError:
            pass


def terminate_process_group(
    process: Optional[subprocess.Popen[Any]],
    *,
    pgid: int,
    term_grace: float = DEFAULT_TERM_GRACE_SECONDS,
    kill_grace: float = DEFAULT_KILL_GRACE_SECONDS,
) -> bool:
    """Terminate and verify a previously isolated process group.

    ``pgid`` must be the value captured when the child was spawned.  This
    function deliberately never derives it from ``process.pid`` after a timeout:
    the leader may already have exited while a grandchild remains in the group.
    """

    known_pgid = int(pgid)
    _signal_process_group(known_pgid, signal.SIGTERM)

    drained = _wait_group_empty(known_pgid, term_grace)
    if not drained:
        _signal_process_group(known_pgid, signal.SIGKILL)
        drained = _wait_group_empty(known_pgid, kill_grace)

    if process is not None:
        try:
            process.wait(timeout=max(0.0, float(kill_grace)))
        except (ChildProcessError, ProcessLookupError):
            pass
        except subprocess.TimeoutExpired:
            drained = False
    return drained


def run_process_group(
    args: Sequence[str] | str,
    *,
    timeout: Optional[float] = None,
    input: Any = None,
    capture_output: bool = False,
    check: bool = False,
    cwd: Optional[str] = None,
    env: Optional[Mapping[str, str]] = None,
    text: Optional[bool] = None,
    encoding: Optional[str] = None,
    errors: Optional[str] = None,
    stdin: Any = None,
    stdout: Any = None,
    stderr: Any = None,
    term_grace: float = DEFAULT_TERM_GRACE_SECONDS,
    kill_grace: float = DEFAULT_KILL_GRACE_SECONDS,
    **popen_kwargs: Any,
) -> subprocess.CompletedProcess[Any]:
    """A ``subprocess.run``-compatible subset with timeout tree cleanup."""

    if capture_output:
        if stdout is not None or stderr is not None:
            raise ValueError(
                "stdout and stderr arguments may not be used with capture_output")
        stdout = subprocess.PIPE
        stderr = subprocess.PIPE
    if input is not None:
        if stdin is not None:
            raise ValueError("stdin and input arguments may not both be used")
        stdin = subprocess.PIPE
    if "start_new_session" in popen_kwargs:
        raise ValueError("run_process_group always owns start_new_session")

    process = subprocess.Popen(
        args,
        stdin=stdin,
        stdout=stdout,
        stderr=stderr,
        cwd=cwd,
        env=env,
        text=text,
        encoding=encoding,
        errors=errors,
        start_new_session=True,
        **popen_kwargs,
    )
    # POSIX setsid(2), requested by start_new_session, makes PID == PGID before
    # exec.  Keep this known value; never look it up after the leader can exit.
    pgid = process.pid
    try:
        stdout_data, stderr_data = process.communicate(input=input, timeout=timeout)
    except subprocess.TimeoutExpired as exc:
        terminate_process_group(
            process,
            pgid=pgid,
            term_grace=term_grace,
            kill_grace=kill_grace,
        )
        try:
            stdout_data, stderr_data = process.communicate(timeout=kill_grace)
        except subprocess.TimeoutExpired:
            stdout_data, stderr_data = exc.output, exc.stderr
        raise subprocess.TimeoutExpired(
            exc.cmd,
            exc.timeout,
            output=stdout_data if stdout_data is not None else exc.output,
            stderr=stderr_data if stderr_data is not None else exc.stderr,
        ) from exc

    completed = subprocess.CompletedProcess(
        args, process.returncode, stdout_data, stderr_data)
    if check:
        completed.check_returncode()
    return completed


def command_reaches_a1id(args: Sequence[Any] | str) -> bool:
    """Return whether a command directly invokes the repository a1 wrapper."""

    values: Iterable[Any]
    if isinstance(args, str):
        values = (args,)
    else:
        values = args
    return any(os.path.basename(str(value)) == "a1id" for value in values)


__all__ = [
    "command_reaches_a1id",
    "run_process_group",
    "terminate_process_group",
]
