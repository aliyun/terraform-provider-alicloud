#!/usr/bin/env python3
"""Portable shared/exclusive gate for a1 lock-file maintenance.

macOS does not ship the ``flock(1)`` command, but both macOS and Linux expose
BSD advisory locks through Python's standard-library ``fcntl.flock``.  The
locked descriptor is deliberately inherited across ``exec`` so it covers the
complete a1id/a1 process lifetime.  The kernel releases it on normal exit,
SIGTERM, SIGKILL, or an exec failure; there is no PID marker to reap and no
PID-reuse or check/delete race.

Shared mode waits for a cleaner already holding the exclusive lock.  Exclusive
mode is non-blocking: if any a1id invocation is active, cleanup safely skips
instead of delaying bridge startup.
"""

from __future__ import annotations

import argparse
import errno
import fcntl
import os
from pathlib import Path
import sys


GATE_NAME = ".jarvis-a1-file-gate"


def _close_inherited_gate() -> None:
    """Drop a mismatched parent gate before acquiring this invocation's mode."""

    value = os.environ.get("JARVIS_A1_GATE_FD", "")
    if value.isdigit():
        try:
            os.close(int(value))
        except OSError:
            pass
    for name in (
        "JARVIS_A1_GATE_FD",
        "JARVIS_A1_GATE_MODE",
        "JARVIS_A1_GATE_ROOT",
    ):
        os.environ.pop(name, None)


def _open_gate(root: str) -> tuple[int, str]:
    root_path = Path(root).expanduser()
    root_path.mkdir(mode=0o700, parents=True, exist_ok=True)
    # Keep the caller-visible spelling (notably macOS /var vs /private/var)
    # because a1id exposes A1_CONFIG_DIR and existing callers compare it.  open
    # still resolves aliases to the same gate inode, so mutual exclusion does
    # not require rewriting A1ID_ROOT through realpath.
    canonical_root = os.path.abspath(str(root_path))
    flags = os.O_CREAT | os.O_RDWR
    if hasattr(os, "O_NOFOLLOW"):
        flags |= os.O_NOFOLLOW
    descriptor = os.open(
        os.path.join(canonical_root, GATE_NAME),
        flags,
        0o600,
    )
    os.set_inheritable(descriptor, True)
    return descriptor, canonical_root


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("mode", choices=("shared", "exclusive"))
    parser.add_argument("--root", required=True)
    parser.add_argument("command", nargs=argparse.REMAINDER)
    args = parser.parse_args(argv)
    command = list(args.command)
    if command[:1] == ["--"]:
        command = command[1:]
    if not command:
        parser.error("a command is required after --")

    _close_inherited_gate()
    descriptor, canonical_root = _open_gate(args.root)
    operation = fcntl.LOCK_SH
    if args.mode == "exclusive":
        operation = fcntl.LOCK_EX | fcntl.LOCK_NB
    try:
        fcntl.flock(descriptor, operation)
    except BlockingIOError:
        os.close(descriptor)
        # Cleanup is opportunistic.  An active shared owner proves a1 may own
        # one of the files, so skipping the complete pass is the safe result.
        print(
            "a1-locks-clean: skip (active a1id gate owner detected)",
            file=sys.stderr,
        )
        return 0
    except OSError as exc:
        os.close(descriptor)
        if args.mode == "exclusive" and exc.errno in (errno.EACCES, errno.EAGAIN):
            print(
                "a1-locks-clean: skip (active a1id gate owner detected)",
                file=sys.stderr,
            )
            return 0
        raise

    environment = os.environ.copy()
    environment["A1ID_ROOT"] = canonical_root
    environment["JARVIS_A1_GATE_FD"] = str(descriptor)
    environment["JARVIS_A1_GATE_MODE"] = args.mode
    environment["JARVIS_A1_GATE_ROOT"] = canonical_root
    os.execvpe(command[0], command, environment)
    return 127


if __name__ == "__main__":
    raise SystemExit(main())
