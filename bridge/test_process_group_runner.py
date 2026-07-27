#!/usr/bin/env python3
"""Fault-injection coverage for bounded subprocess process trees."""

from __future__ import annotations

import ast
import os
from pathlib import Path
import signal
import subprocess
import sys
import tempfile
import time
import unittest
from unittest import mock

from bridge import process_group_runner as runner


def _pid_is_running(pid: int) -> bool:
    try:
        state = subprocess.check_output(
            ["ps", "-o", "stat=", "-p", str(pid)],
            text=True,
            stderr=subprocess.DEVNULL,
        ).strip()
    except subprocess.CalledProcessError:
        return False
    return bool(state) and not state.startswith("Z")


class ProcessGroupRunnerTest(unittest.TestCase):
    def test_timeout_kills_grandchild_after_leader_has_exited(self):
        with tempfile.TemporaryDirectory() as tmp:
            pid_file = Path(tmp) / "grandchild.pid"
            grandchild = (
                "import os,pathlib,signal,sys,time\n"
                "signal.signal(signal.SIGTERM, signal.SIG_IGN)\n"
                "pathlib.Path(sys.argv[1]).write_text(str(os.getpid()))\n"
                "time.sleep(60)\n"
            )
            leader = (
                "import pathlib,subprocess,sys,time\n"
                "subprocess.Popen([sys.executable,'-c',sys.argv[2],sys.argv[1]])\n"
                "deadline=time.monotonic()+2\n"
                "while not pathlib.Path(sys.argv[1]).exists():\n"
                "    if time.monotonic() >= deadline: raise SystemExit(2)\n"
                "    time.sleep(0.01)\n"
            )
            with self.assertRaises(subprocess.TimeoutExpired):
                runner.run_process_group(
                    [sys.executable, "-c", leader, str(pid_file), grandchild],
                    timeout=0.3,
                    capture_output=True,
                    text=True,
                    term_grace=0.2,
                    kill_grace=1.0,
                )
            grandchild = int(pid_file.read_text())
            deadline = time.monotonic() + 2
            while time.monotonic() < deadline and _pid_is_running(grandchild):
                time.sleep(0.05)
            self.assertFalse(
                _pid_is_running(grandchild),
                "timeout left the leader's grandchild running",
            )

    def test_timeout_kill_releases_inherited_a1_file_gate(self):
        gate = (
            Path(__file__).resolve().parent.parent
            / "bootstrap"
            / "a1-file-gate.py"
        )
        with tempfile.TemporaryDirectory() as tmp:
            started = Path(tmp) / "shared-started"
            exclusive_entered = Path(tmp) / "exclusive-entered"
            holder = (
                "import pathlib,signal,sys,time\n"
                "signal.signal(signal.SIGTERM, signal.SIG_IGN)\n"
                "pathlib.Path(sys.argv[1]).write_text('started')\n"
                "time.sleep(60)\n"
            )
            with self.assertRaises(subprocess.TimeoutExpired):
                runner.run_process_group(
                    [
                        sys.executable,
                        str(gate),
                        "--root",
                        tmp,
                        "shared",
                        "--",
                        sys.executable,
                        "-c",
                        holder,
                        str(started),
                    ],
                    timeout=0.3,
                    capture_output=True,
                    text=True,
                    term_grace=0.2,
                    kill_grace=1.0,
                )
            self.assertTrue(started.exists(), "shared gate holder never entered")

            writer = (
                "import pathlib,sys\n"
                "pathlib.Path(sys.argv[1]).write_text('exclusive')\n"
            )
            result = subprocess.run(
                [
                    sys.executable,
                    str(gate),
                    "--root",
                    tmp,
                    "exclusive",
                    "--",
                    sys.executable,
                    "-c",
                    writer,
                    str(exclusive_entered),
                ],
                capture_output=True,
                text=True,
                timeout=2,
            )
            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertTrue(
                exclusive_entered.exists(),
                "SIGKILL left the inherited shared gate locked",
            )

    def test_terminate_uses_captured_pgid_without_leader_lookup(self):
        process = mock.Mock()
        process.pid = 777
        process.wait.return_value = -signal.SIGTERM
        with mock.patch.object(
                runner, "_wait_group_empty", side_effect=[False, True]), \
             mock.patch.object(runner.os, "killpg") as killpg, \
             mock.patch.object(
                 runner.os,
                 "getpgid",
                 side_effect=AssertionError("must not derive PGID from dead leader"),
                 create=True):
            self.assertTrue(
                runner.terminate_process_group(
                    process, pgid=777, term_grace=0, kill_grace=0))
        self.assertEqual(
            killpg.call_args_list,
            [mock.call(777, signal.SIGTERM), mock.call(777, signal.SIGKILL)],
        )

    def test_terminate_supports_adopted_group_without_popen_handle(self):
        with mock.patch.object(
                runner, "_wait_group_empty", return_value=True), \
             mock.patch.object(runner.os, "killpg") as killpg:
            self.assertTrue(
                runner.terminate_process_group(
                    None,
                    pgid=778,
                    term_grace=0,
                    kill_grace=0,
                )
            )
        killpg.assert_called_once_with(778, signal.SIGTERM)

    def test_command_classifier_matches_only_a1id_basename(self):
        self.assertTrue(runner.command_reaches_a1id(["/repo/bin/a1id", "--"]))
        self.assertFalse(runner.command_reaches_a1id(["a1", "auth", "whoami"]))
        self.assertFalse(runner.command_reaches_a1id(["/tmp/not-a1id", "--"]))

    def test_bridge_has_no_direct_subprocess_run_of_a1id(self):
        bridge_root = Path(__file__).resolve().parent
        violations = []
        a1_entrypoints = (
            "a1id", "aone-fields.sh", "aone-get.sh", "claim.sh", "wrap.sh",
        )
        for path in bridge_root.rglob("*.py"):
            if path.name.startswith("test_"):
                continue
            source = path.read_text(encoding="utf-8")
            tree = ast.parse(source, filename=str(path))
            for node in ast.walk(tree):
                if not isinstance(node, ast.Call):
                    continue
                function = node.func
                if not (
                    isinstance(function, ast.Attribute)
                    and function.attr == "run"
                    and isinstance(function.value, ast.Name)
                    and function.value.id == "subprocess"
                ):
                    continue
                segment = ast.get_source_segment(source, node) or ""
                if any(entrypoint in segment for entrypoint in a1_entrypoints):
                    violations.append(
                        "%s:%s" % (path.relative_to(bridge_root), node.lineno))
        self.assertEqual(
            violations,
            [],
            "bin/a1id must use run_process_group, not subprocess.run(timeout)",
        )


if __name__ == "__main__":
    unittest.main()
