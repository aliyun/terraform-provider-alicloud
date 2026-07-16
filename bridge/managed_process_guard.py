#!/usr/bin/env python3
"""Keep a managed subprocess tied to the bridge that owns its lease.

The bridge passes two inherited pipes:

* ``gate_fd`` stays closed until the process-group leader has been durably bound
  to ``jarvis_session``.  EOF before the single-byte grant means the command is
  never started.
* ``sentinel_fd`` is held open only by the bridge.  EOF means the bridge died,
  including SIGKILL, so this guard terminates the complete managed process group.

The guard is the process-group/session leader.  The actual command and all of its
descendants inherit that group, which gives the live bridge and this crash guard
one stable PID/PGID to fence.
"""

import argparse
import os
import select
import signal
import subprocess
import sys
import time


POLL_SECONDS = 0.05
EXIT_NOT_GRANTED = 125


def _pipe_readable(fd, timeout):
    try:
        return bool(select.select([fd], [], [], timeout)[0])
    except (OSError, ValueError):
        return True


def _read_grant(gate_fd, sentinel_fd, should_stop):
    """Wait for a grant while refusing to launch after the owner disappears."""
    while True:
        if should_stop():
            return False
        try:
            ready, _, _ = select.select([gate_fd, sentinel_fd], [], [], POLL_SECONDS)
        except (OSError, ValueError):
            return False
        if sentinel_fd in ready:
            try:
                if os.read(sentinel_fd, 1) == b"":
                    return False
            except OSError:
                return False
        if gate_fd in ready:
            try:
                return os.read(gate_fd, 1) == b"G"
            except OSError:
                return False


def _owner_gone(sentinel_fd):
    if not _pipe_readable(sentinel_fd, 0):
        return False
    try:
        return os.read(sentinel_fd, 1) == b""
    except OSError:
        return True


def _signal_group(signum):
    """Signal descendants in our group without terminating the guard first."""
    previous = signal.signal(signum, signal.SIG_IGN)
    try:
        os.killpg(os.getpgrp(), signum)
    except ProcessLookupError:
        pass
    finally:
        signal.signal(signum, previous)


def _stop_group(child, grace_seconds):
    _signal_group(signal.SIGTERM)
    deadline = time.monotonic() + max(0.0, grace_seconds)
    while child.poll() is None and time.monotonic() < deadline:
        time.sleep(POLL_SECONDS)
    if child.poll() is None:
        # This intentionally kills the guard as well.  Its owner is already gone
        # (or is synchronously fencing us), and SIGKILL is the only reliable group
        # boundary for a command that ignores SIGTERM.
        os.killpg(os.getpgrp(), signal.SIGKILL)


def _group_members():
    """Return other live members of this private managed process group."""
    pgid = os.getpgrp()
    current = os.getpid()
    try:
        with subprocess.Popen(
                ["ps", "-axo", "pid=,pgid="], text=True,
                stdout=subprocess.PIPE, stderr=subprocess.DEVNULL) as probe:
            output, _ = probe.communicate(timeout=2)
            probe_pid = probe.pid
    except (OSError, subprocess.SubprocessError):
        return []
    members = []
    for line in output.splitlines():
        fields = line.split()
        if len(fields) != 2:
            continue
        try:
            pid, process_group = int(fields[0]), int(fields[1])
        except ValueError:
            continue
        if process_group == pgid and pid not in {current, probe_pid}:
            members.append(pid)
    return members


def _signal_members(pids, signum):
    for pid in pids:
        try:
            os.kill(pid, signum)
        except ProcessLookupError:
            pass


def _cleanup_residual_group(grace_seconds):
    """A completed CLI must not leave tools/background grandchildren running."""
    residual = _group_members()
    if not residual:
        return
    _signal_members(residual, signal.SIGTERM)
    deadline = time.monotonic() + max(0.0, grace_seconds)
    while time.monotonic() < deadline:
        residual = _group_members()
        if not residual:
            return
        time.sleep(POLL_SECONDS)
    _signal_members(_group_members(), signal.SIGKILL)
    deadline = time.monotonic() + 2.0
    while _group_members() and time.monotonic() < deadline:
        time.sleep(POLL_SECONDS)


def _exit_code(returncode):
    if returncode is None:
        return 1
    if returncode < 0:
        return 128 - returncode
    return returncode


def run(gate_fd, sentinel_fd, command, grace_seconds):
    stop_requested = False

    def request_stop(_signum, _frame):
        nonlocal stop_requested
        stop_requested = True

    # Install handlers before the grant wait and child spawn.  This closes the
    # narrow stop race where the guard leader could previously exit before its
    # TERM-deaf child had a handler-owned cleanup path.
    signal.signal(signal.SIGTERM, request_stop)
    signal.signal(signal.SIGINT, request_stop)

    if not _read_grant(gate_fd, sentinel_fd, lambda: stop_requested):
        return EXIT_NOT_GRANTED
    try:
        os.close(gate_fd)
    except OSError:
        pass

    # Close the race between reading the grant and spawning the command.  If the
    # bridge is already gone, no external code is executed.
    if stop_requested or _owner_gone(sentinel_fd):
        return EXIT_NOT_GRANTED

    child = subprocess.Popen(command, stdin=subprocess.DEVNULL)

    while child.poll() is None:
        if stop_requested or _owner_gone(sentinel_fd):
            _stop_group(child, grace_seconds)
            break
        time.sleep(POLL_SECONDS)
    returncode = child.wait()
    _cleanup_residual_group(grace_seconds)
    return _exit_code(returncode)


def main(argv=None):
    parser = argparse.ArgumentParser()
    parser.add_argument("--gate-fd", required=True, type=int)
    parser.add_argument("--sentinel-fd", required=True, type=int)
    parser.add_argument("--grace-seconds", type=float, default=2.0)
    parser.add_argument("command", nargs=argparse.REMAINDER)
    args = parser.parse_args(argv)
    command = args.command[1:] if args.command[:1] == ["--"] else args.command
    if not command:
        parser.error("a command is required after --")
    return run(args.gate_fd, args.sentinel_fd, command, args.grace_seconds)


if __name__ == "__main__":
    raise SystemExit(main())
