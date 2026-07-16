#!/usr/bin/env python3
"""Shared subprocess runtime for Task and ephemeral execution.

The runtime deliberately has no Task/Session semantics.  A caller may request
the fenced ``ProcessGuardian`` launch path, but ownership transitions remain in
``TaskExecutor``/``SessionLifecycle``.  Ephemeral jobs use the same timeout,
process-group cleanup, stdout/stderr capture, and spawn observability without
creating control-plane state.
"""

from __future__ import annotations

import os
import signal
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Callable, Optional, Sequence, Tuple


SpawnCallback = Callable[[Any], None]
GuardedSpawn = Callable[[Sequence[str], Path, SpawnCallback], Tuple[Any, Optional[int]]]


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
            "managed_process_guard.py")

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
              on_spawn: SpawnCallback) -> Tuple[Any, int]:
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
                os.environ.get("JARVIS_MANAGED_GUARD_GRACE_SEC", "2"),
                "--",
            ] + list(argv)
            process = subprocess.Popen(
                command, cwd=cwd, text=True, stdin=subprocess.DEVNULL,
                stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                start_new_session=True, pass_fds=(gate_read, sentinel_read))
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
            os.write(gate_write, b"G")
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
                     guarded_spawn: Optional[GuardedSpawn] = None) -> ExecutionResult:
        sentinel_write = None
        if guarded:
            spawn = guarded_spawn or self.guardian.spawn
            process, sentinel_write = spawn(argv, cwd, on_spawn)
        else:
            process = subprocess.Popen(
                list(argv), cwd=cwd, text=True, stdin=subprocess.DEVNULL,
                stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                start_new_session=True)
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
