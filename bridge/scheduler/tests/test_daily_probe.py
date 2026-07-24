from __future__ import annotations

from datetime import datetime
from pathlib import Path
from tempfile import TemporaryDirectory
from types import SimpleNamespace
import unittest
from unittest import mock
from zoneinfo import ZoneInfo

from bridge.headless_runtime import HeadlessResult, Lane, SessionPolicy
from bridge.scheduler.model import (
    DailySchedule,
    HandlerRunner,
    JobResultStatus,
    MisfirePolicy,
    ScheduledJobDefinition,
)
from bridge.scheduler.runners.daily_probe import (
    DailyProbeRunner,
    probe_prompt,
)


SHANGHAI = ZoneInfo("Asia/Shanghai")


def definition() -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        "daily.probe",
        3,
        "probe",
        DailySchedule(10, 0, "Asia/Shanghai"),
        HandlerRunner("daily_probe"),
        MisfirePolicy.CURRENT_DAY,
        300,
        True,
    )


def success_text(summary: str = "本轮完成") -> str:
    return (
        "human readable\n"
        '[[PROBE_RESULT:{"outcome":"success","cleanup_ok":true,'
        '"summary":"%s"}]]' % summary
    )


class FakeRuntime:
    def __init__(self, result: HeadlessResult):
        self.result = result
        self.requests = []

    def execute(self, request):
        self.requests.append(request)
        return self.result


class DailyProbeRunnerTests(unittest.TestCase):
    def _runner(self, runtime, root, **kwargs):
        return DailyProbeRunner(
            runtime=runtime,
            summary_root=Path(root),
            environ={
                "JARVIS_DISPATCH_TIMEOUT": "60",
                "JARVIS_DISPATCH_RETRY_MAX": "1",
                "JARVIS_DISPATCH_RETRY_BACKOFF": "0",
                "JARVIS_PROBE_SUMMARY_WRITE_RETRIES": "2",
            },
            logger=SimpleNamespace(info=lambda *a, **k: None,
                                   warning=lambda *a, **k: None),
            sleeper=lambda _seconds: None,
            **kwargs,
        )

    def test_success_uses_scheduled_boundary_and_writes_validated_summary(self):
        runtime = FakeRuntime(
            HeadlessResult(success_text("摘要内容"), False, "success", 1))
        with TemporaryDirectory() as root:
            result = self._runner(runtime, root).run(
                definition(),
                datetime(2026, 7, 23, 10, 0, tzinfo=SHANGHAI),
            )
            summary = Path(root, "probe-2026-07-23-summary.md").read_text()
        self.assertIs(result.status, JobResultStatus.SUCCEEDED)
        self.assertEqual(summary, "摘要内容")
        request = runtime.requests[0]
        self.assertIs(request.session_policy, SessionPolicy.NEW)
        self.assertIs(request.lane, Lane.TERRAFORM)
        self.assertTrue(request.guarded)
        self.assertIsNotNone(request.on_spawn)
        self.assertIn("probe-2026-07-23", request.prompt)
        self.assertIn("[[PROBE_RESULT:", request.prompt)

    def test_retry_after_midnight_keeps_previous_daily_round(self):
        runtime = FakeRuntime(
            HeadlessResult(success_text(), False, "success", 1))
        with TemporaryDirectory() as root:
            result = self._runner(runtime, root).run(
                definition(),
                datetime(2026, 7, 24, 0, 5, tzinfo=SHANGHAI),
            )
            names = [path.name for path in Path(root).iterdir()]
        self.assertIs(result.status, JobResultStatus.SUCCEEDED)
        self.assertEqual(names, ["probe-2026-07-23-summary.md"])
        self.assertIn("probe-2026-07-23", runtime.requests[0].prompt)

    def test_execution_error_and_invalid_protocol_are_retryable(self):
        cases = (
            HeadlessResult("timeout", True, "timeout", 1),
            HeadlessResult("", False, "no_result", 1),
            HeadlessResult("plain text", False, "success", 1),
            HeadlessResult(
                '[[PROBE_RESULT:{"outcome":"success","cleanup_ok":false,'
                '"summary":"dirty"}]]',
                False,
                "success",
                1,
            ),
            HeadlessResult(
                '[[PROBE_RESULT:{"outcome":"success","cleanup_ok":true,'
                '"summary":""}]]',
                False,
                "success",
                1,
            ),
        )
        for headless_result in cases:
            with self.subTest(result=headless_result):
                runtime = FakeRuntime(headless_result)
                with TemporaryDirectory() as root:
                    result = self._runner(runtime, root).run(
                        definition(),
                        datetime(2026, 7, 23, 10, 0, tzinfo=SHANGHAI),
                    )
                    self.assertEqual(list(Path(root).iterdir()), [])
                self.assertIs(
                    result.status, JobResultStatus.RETRYABLE_FAILURE)

    def test_summary_write_exhaustion_is_permanent_after_finite_retries(self):
        runtime = FakeRuntime(
            HeadlessResult(success_text(), False, "success", 1))
        with TemporaryDirectory() as root:
            runner = self._runner(runtime, root)
            with mock.patch.object(
                    runner, "_write_summary",
                    side_effect=OSError("disk full")) as write:
                result = runner.run(
                    definition(),
                    datetime(2026, 7, 23, 10, 0, tzinfo=SHANGHAI),
                )
        self.assertIs(result.status, JobResultStatus.PERMANENT_FAILURE)
        self.assertEqual(write.call_count, 3)
        self.assertEqual(len(runtime.requests), 1)

    def test_prompt_requires_direct_aone_and_machine_result(self):
        prompt = probe_prompt("probe-2026-07-23")
        self.assertIn("直接创建 Aone", prompt)
        self.assertIn('"cleanup_ok":true', prompt)
        self.assertNotIn("files drafts", prompt)


if __name__ == "__main__":
    unittest.main()
