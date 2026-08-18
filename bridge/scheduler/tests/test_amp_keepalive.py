from __future__ import annotations

import logging
from datetime import datetime, timezone
from pathlib import Path
import subprocess
from types import SimpleNamespace
import unittest

from bridge.scheduler.model import (
    HandlerRunner,
    IntervalSchedule,
    JobResultStatus,
    MisfirePolicy,
    ScheduledJobDefinition,
)
from bridge.scheduler.runners.amp_keepalive import AmpKeepaliveRunner


UTC = timezone.utc


def definition(
    *, job_id: str = "amp.keepalive", schedule: object | None = None,
) -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        job_id,
        1,
        "keep AMP authentication warm",
        schedule or IntervalSchedule(600, True),
        HandlerRunner("amp_keepalive"),
        MisfirePolicy.COALESCE,
        60,
    )


class RecordingCommandRunner:
    def __init__(self, result=None, error: BaseException | None = None):
        self.result = result or SimpleNamespace(
            returncode=0, stdout="branch response", stderr="")
        self.error = error
        self.calls: list[tuple[list[str], dict[str, object]]] = []

    def __call__(self, command, **kwargs):
        self.calls.append((list(command), dict(kwargs)))
        if self.error is not None:
            raise self.error
        return self.result


class AmpKeepaliveRunnerTests(unittest.TestCase):
    def runner(self, command_runner, environ):
        return AmpKeepaliveRunner(
            repo_root=Path("/repo"),
            logger=logging.getLogger(__name__),
            environ=environ,
            command_runner=command_runner,
        )

    def test_success_runs_real_branch_query_through_isolated_safe_wrapper(self):
        command_runner = RecordingCommandRunner()
        runner = self.runner(command_runner, {
            "JARVIS_AMP_KEEPALIVE_PROJECT_ID": "123456",
        })

        result = runner.run(definition(), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.SUCCEEDED)
        self.assertEqual(len(command_runner.calls), 1)
        command, kwargs = command_runner.calls[0]
        self.assertEqual(command, [
            "/usr/bin/python3",
            "-I",
            "/repo/bootstrap/amp_safe.py",
            "--repo-root",
            "/repo",
            "branch",
            "list",
            "--project-id",
            "123456",
        ])
        self.assertEqual(kwargs["cwd"], "/repo")
        self.assertEqual(kwargs["timeout"], 30.0)
        self.assertTrue(kwargs["capture_output"])
        self.assertTrue(kwargs["text"])
        self.assertIs(kwargs["stdin"], subprocess.DEVNULL)

    def test_configured_timeout_is_applied(self):
        command_runner = RecordingCommandRunner()
        runner = self.runner(command_runner, {
            "JARVIS_AMP_KEEPALIVE_PROJECT_ID": "9",
            "JARVIS_AMP_KEEPALIVE_TIMEOUT": "4.5",
        })

        result = runner.run(definition(), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.SUCCEEDED)
        self.assertEqual(command_runner.calls[0][1]["timeout"], 4.5)

    def test_missing_project_id_is_permanent_and_does_not_spawn(self):
        command_runner = RecordingCommandRunner()
        result = self.runner(command_runner, {}).run(
            definition(), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.PERMANENT_FAILURE)
        self.assertEqual(result.error, "amp keepalive project id is missing")
        self.assertEqual(command_runner.calls, [])

    def test_project_id_must_be_a_strict_positive_decimal(self):
        for project_id in (
            "0", "-1", "+1", "01", "1.0", "abc", " 9", "1" * 20,
        ):
            with self.subTest(project_id=project_id):
                command_runner = RecordingCommandRunner()
                result = self.runner(command_runner, {
                    "JARVIS_AMP_KEEPALIVE_PROJECT_ID": project_id,
                }).run(definition(), datetime.now(UTC))

                self.assertIs(
                    result.status, JobResultStatus.PERMANENT_FAILURE)
                self.assertEqual(
                    result.error, "amp keepalive project id is invalid")
                self.assertEqual(command_runner.calls, [])

    def test_invalid_timeout_is_permanent_and_does_not_spawn(self):
        for timeout in ("0", "-1", "nan", "inf", "not-a-number"):
            with self.subTest(timeout=timeout):
                command_runner = RecordingCommandRunner()
                result = self.runner(command_runner, {
                    "JARVIS_AMP_KEEPALIVE_PROJECT_ID": "9",
                    "JARVIS_AMP_KEEPALIVE_TIMEOUT": timeout,
                }).run(definition(), datetime.now(UTC))

                self.assertIs(
                    result.status, JobResultStatus.PERMANENT_FAILURE)
                self.assertEqual(
                    result.error, "amp keepalive timeout is invalid")
                self.assertEqual(command_runner.calls, [])

    def test_nonzero_exit_is_retryable_without_captured_output(self):
        secret = "https://oauth.example/login?token=super-secret"
        command_runner = RecordingCommandRunner(SimpleNamespace(
            returncode=17,
            stdout='{"access_token":"super-secret"}',
            stderr=secret,
        ))
        result = self.runner(command_runner, {
            "JARVIS_AMP_KEEPALIVE_PROJECT_ID": "9",
        }).run(definition(), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
        self.assertEqual(result.error, "amp keepalive query failed: nonzero_exit")
        self.assertNotIn("secret", result.error)
        self.assertNotIn("oauth", result.error)
        self.assertLessEqual(len(result.error), 80)

    def test_timeout_is_retryable_without_captured_output(self):
        command_runner = RecordingCommandRunner(error=subprocess.TimeoutExpired(
            ["safe-wrapper"], 30,
            output=b'access_token="super-secret"',
            stderr=b"https://oauth.example/login",
        ))
        result = self.runner(command_runner, {
            "JARVIS_AMP_KEEPALIVE_PROJECT_ID": "9",
        }).run(definition(), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
        self.assertEqual(result.error, "amp keepalive query failed: timeout")
        self.assertNotIn("secret", result.error)
        self.assertNotIn("oauth", result.error)

    def test_os_error_is_retryable_without_exception_detail(self):
        command_runner = RecordingCommandRunner(
            error=OSError("token=super-secret https://oauth.example/login"))
        result = self.runner(command_runner, {
            "JARVIS_AMP_KEEPALIVE_PROJECT_ID": "9",
        }).run(definition(), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
        self.assertEqual(
            result.error, "amp keepalive query failed: process_error")
        self.assertNotIn("secret", result.error)
        self.assertNotIn("oauth", result.error)

    def test_mismatched_definition_is_permanent(self):
        command_runner = RecordingCommandRunner()
        result = self.runner(command_runner, {
            "JARVIS_AMP_KEEPALIVE_PROJECT_ID": "9",
        }).run(definition(job_id="aone.scan"), datetime.now(UTC))

        self.assertIs(result.status, JobResultStatus.PERMANENT_FAILURE)
        self.assertEqual(
            result.error, "amp keepalive runner received mismatched definition")
        self.assertEqual(command_runner.calls, [])


if __name__ == "__main__":
    unittest.main()
