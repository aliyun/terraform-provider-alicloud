#!/usr/bin/env python3
"""Cross-platform process birth identity for PID-reuse-safe supervision."""

from __future__ import annotations

import ctypes
import hashlib
import os
from pathlib import Path
import subprocess
import sys
from typing import Optional


def _linux_start_identity(pid: int) -> Optional[str]:
    stat_path = Path("/proc") / str(int(pid)) / "stat"
    try:
        value = stat_path.read_text(encoding="utf-8")
        closing = value.rfind(")")
        # Fields after comm start at proc field 3; starttime is field 22.
        start_ticks = value[closing + 2:].split()[19]
        boot_id = Path("/proc/sys/kernel/random/boot_id").read_text(
            encoding="utf-8").strip()
    except (OSError, IndexError, ValueError):
        return None
    return "linux:%s:%s" % (boot_id, start_ticks)


def _darwin_start_identity(pid: int) -> Optional[str]:
    if sys.platform != "darwin":
        return None

    class ProcBsdInfo(ctypes.Structure):
        _fields_ = [
            ("pbi_flags", ctypes.c_uint32),
            ("pbi_status", ctypes.c_uint32),
            ("pbi_xstatus", ctypes.c_uint32),
            ("pbi_pid", ctypes.c_uint32),
            ("pbi_ppid", ctypes.c_uint32),
            ("pbi_uid", ctypes.c_uint32),
            ("pbi_gid", ctypes.c_uint32),
            ("pbi_ruid", ctypes.c_uint32),
            ("pbi_rgid", ctypes.c_uint32),
            ("pbi_svuid", ctypes.c_uint32),
            ("pbi_svgid", ctypes.c_uint32),
            ("pbi_rfu_1", ctypes.c_uint32),
            ("pbi_comm", ctypes.c_char * 16),
            ("pbi_name", ctypes.c_char * 32),
            ("pbi_nfiles", ctypes.c_uint32),
            ("pbi_pgid", ctypes.c_uint32),
            ("pbi_pjobc", ctypes.c_uint32),
            ("e_tdev", ctypes.c_uint32),
            ("e_tpgid", ctypes.c_uint32),
            ("pbi_nice", ctypes.c_int32),
            ("pbi_start_tvsec", ctypes.c_uint64),
            ("pbi_start_tvusec", ctypes.c_uint64),
        ]

    try:
        libproc = ctypes.CDLL("/usr/lib/libproc.dylib", use_errno=True)
        proc_pidinfo = libproc.proc_pidinfo
        proc_pidinfo.argtypes = [
            ctypes.c_int,
            ctypes.c_int,
            ctypes.c_uint64,
            ctypes.c_void_p,
            ctypes.c_int,
        ]
        proc_pidinfo.restype = ctypes.c_int
        info = ProcBsdInfo()
        size = ctypes.sizeof(info)
        result = proc_pidinfo(
            int(pid),
            3,  # PROC_PIDTBSDINFO
            0,
            ctypes.byref(info),
            size,
        )
    except (OSError, ValueError):
        return None
    if result != size or info.pbi_pid != int(pid):
        return None
    return "darwin:%s:%s" % (
        info.pbi_start_tvsec,
        info.pbi_start_tvusec,
    )


def process_start_identity(pid: int) -> Optional[str]:
    """Return boot-scoped process birth identity with subsecond/tick precision."""

    identity = _linux_start_identity(pid) or _darwin_start_identity(pid)
    if identity is not None:
        return identity
    # Conservative fallback for other POSIX hosts.  Command is included to
    # make accidental equality stricter, though Linux/Darwin use native birth
    # metadata above.
    try:
        started = subprocess.check_output(
            ["ps", "-o", "lstart=", "-p", str(int(pid))],
            text=True,
            stderr=subprocess.DEVNULL,
        )
        command = process_command(pid) or ""
    except (OSError, subprocess.SubprocessError, ValueError):
        return None
    token = "".join(started.split())
    command_digest = hashlib.sha256(command.encode("utf-8")).hexdigest()[:16]
    return "posix:%s:%s" % (token, command_digest) if token else None


def process_command(pid: int) -> Optional[str]:
    try:
        value = subprocess.check_output(
            ["ps", "-o", "command=", "-p", str(int(pid))],
            text=True,
            stderr=subprocess.DEVNULL,
        ).strip()
    except (OSError, subprocess.SubprocessError, ValueError):
        return None
    return value or None


def pid_exists(pid: int) -> bool:
    try:
        os.kill(int(pid), 0)
        return True
    except ProcessLookupError:
        return False
    except PermissionError:
        return True


def process_identity_record(pid: int) -> Optional[str]:
    identity = process_start_identity(pid)
    if identity is None:
        return None
    try:
        pgid = os.getpgid(int(pid))
    except OSError:
        return None
    return "%s|%s|%s" % (int(pid), pgid, identity)


def main(argv: list[str]) -> int:
    if len(argv) != 2:
        print("usage: process_identity.py <pid>", file=sys.stderr)
        return 2
    try:
        pid = int(argv[1])
    except ValueError:
        return 2
    record = process_identity_record(pid)
    if record is None:
        return 1
    print(record)
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
