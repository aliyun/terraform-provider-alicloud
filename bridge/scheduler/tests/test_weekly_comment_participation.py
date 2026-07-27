from __future__ import annotations

from datetime import datetime, timedelta, timezone
import io
import json
import logging
from pathlib import Path
from types import SimpleNamespace
import unittest
from unittest import mock
from zoneinfo import ZoneInfo

from bridge.scheduler.model import (
    DailySchedule, HandlerRunner, JobResultStatus, MisfirePolicy,
    ScheduledJobDefinition,
)
from bridge.scheduler.runners import weekly_comment_participation as wcp


SHANGHAI = ZoneInfo("Asia/Shanghai")
UTC = timezone.utc


def definition() -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        wcp.JOB_KEY, 1, "weekly comment participation",
        DailySchedule(6, 42, "Asia/Shanghai"),
        HandlerRunner(wcp.RUNNER_KEY), MisfirePolicy.CURRENT_DAY, 300, True,
    )


class FakeTaskClient:
    def __init__(self, base_url="https://control-plane.test", token="t0k"):
        self.base_url = base_url
        self.token = token
        self.timeout = 5.0


def silent_logger() -> logging.Logger:
    logger = logging.getLogger("test-weekly-comment")
    logger.addHandler(logging.NullHandler())
    return logger


class WeeklyCommentParticipationAggregationTests(unittest.TestCase):
    """口径单测：人/数字人分类、窗口过滤、系统噪声排除、workitemCount 聚合。"""

    def _runner(self, repo_root: Path) -> wcp.WeeklyCommentParticipationRunner:
        return wcp.WeeklyCommentParticipationRunner(
            task_client=FakeTaskClient(), repo_root=repo_root,
            logger=silent_logger())

    def _fake_contact_directory(self):
        by_token = {
            "夏节": {"id": "373108", "name": "夏节", "flower": "夏节"},
            "373108": {"id": "373108", "name": "夏节", "flower": "夏节"},
        }
        return by_token, {}

    def _canned_requirements(self):
        now = datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI)
        in_window = now - timedelta(days=2)  # 2026-07-25, inside [2026-07-20, 2026-07-27]
        return [{
            "id": "1001",
            "project": "1086837",
            "title": "需求A",
            "type": "需求问题",
            "modified": in_window.timestamp(),
        }]

    def _canned_comments(self):
        in_window = "2026-07-23 10:00:00"  # Shanghai wall time, inside the window
        out_of_window = "2026-07-10 10:00:00"  # before window start (2026-07-20)
        return [
            # human, in window — counts twice
            {"author": "夏节", "content": "确认根因是 X", "createdAt": in_window},
            {"author": "夏节", "content": "已修复", "createdAt": in_window},
            # digital, in window
            {"author": "terraform-rd", "content": "已提交 PR #1", "createdAt": in_window},
            # system (kelude) — excluded
            {"author": "kelude", "content": "状态流转", "createdAt": in_window},
            # jarvis-claim bookkeeping — excluded
            {"author": "jarvis", "content": "jarvis-claim claimed #1001", "createdAt": in_window},
            # human but out of window — excluded by window filter
            {"author": "张三", "content": "旧评论", "createdAt": out_of_window},
        ]

    def test_aggregate_classifies_and_filters_by_window(self):
        with mock.patch.object(wcp, "_contact_directory",
                               return_value=self._fake_contact_directory()):
            runner = self._runner(Path("/repo"))
            runner._tf_pools = lambda: [("tf_customer", "1086837", "需求问题")]
            runner._list_active_requirements = (
                lambda project, req_type, ws: self._canned_requirements())
            runner._list_comments = lambda iid: (
                self._canned_comments() if iid == "1001" else [])

            scheduled_for = datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI)
            payload = runner._aggregate(scheduled_for)

        self.assertEqual(payload["totalComments"], 3)  # 夏节×2 + terraform-rd×1
        self.assertEqual(payload["ticketsTouched"], 1)
        self.assertEqual(payload["requirementsCovered"], 1)
        # window boundaries are ISO-8601 with offset, parseable by the frontend.
        self.assertTrue(
            payload["windowStart"].startswith("2026-07-20T06:42:00+08:00"))
        self.assertTrue(
            payload["windowEnd"].startswith("2026-07-27T06:42:00+08:00"))

        participants = payload["participants"]
        by_name = {p["name"]: p for p in participants}
        self.assertEqual(by_name["夏节"]["commentCount"], 2)
        self.assertEqual(by_name["夏节"]["workitemCount"], 1)
        self.assertFalse(by_name["夏节"]["digital"])
        self.assertEqual(by_name["Terraform RD"]["commentCount"], 1)
        self.assertEqual(by_name["Terraform RD"]["workitemCount"], 1)
        self.assertTrue(by_name["Terraform RD"]["digital"])
        # kelude / jarvis-claim / out-of-window never appear.
        self.assertNotIn("kelude", by_name)
        self.assertNotIn("Jarvis", by_name)  # jarvis-claim bookkeeping excluded
        self.assertNotIn("张三", by_name)    # out of window
        # sorted desc by commentCount.
        self.assertEqual([p["name"] for p in participants], ["夏节", "Terraform RD"])

    def test_aggregate_continues_past_a_failing_pool(self):
        with mock.patch.object(wcp, "_contact_directory",
                               return_value=self._fake_contact_directory()):
            runner = self._runner(Path("/repo"))
            runner._tf_pools = lambda: [
                ("tf_customer", "1086837", "需求问题"),
                ("tf_provider", "528766", "产品类需求"),
            ]
            runner._list_active_requirements = (
                lambda project, req_type, ws: None  # query failure for every pool
                if project == "1086837" else self._canned_requirements())
            runner._list_comments = lambda iid: (
                self._canned_comments() if iid == "1001" else [])

            payload = runner._aggregate(
                datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))

        # tf_customer failed (skipped), tf_provider still aggregated.
        self.assertEqual(payload["requirementsCovered"], 1)
        self.assertEqual(payload["totalComments"], 3)

    def test_classify_author_excludes_system_and_bookkeeping(self):
        self.assertIsNone(wcp._classify_author("kelude", "状态流转"))
        self.assertIsNone(wcp._classify_author("云知道平台公共账号", "x"))
        self.assertIsNone(wcp._classify_author("", "x"))
        self.assertIsNone(wcp._classify_author("jarvis", "jarvis-claim released #1"))

    def test_classify_author_marks_digital_workers(self):
        with mock.patch.object(wcp, "_contact_directory",
                               return_value=({}, {})):
            human = wcp._classify_author("过载", "已定位根因")
            self.assertEqual(human[0], "human")
            self.assertFalse(human[3])
            rd = wcp._classify_author("terraform-rd", "已提交 PR")
            self.assertEqual(rd[0], "digital")
            self.assertEqual(rd[2], "Terraform RD")
            self.assertTrue(rd[3])


class WeeklyCommentParticipationPublishTests(unittest.TestCase):
    def _runner(self) -> wcp.WeeklyCommentParticipationRunner:
        return wcp.WeeklyCommentParticipationRunner(
            task_client=FakeTaskClient(token="secret-token"),
            repo_root=Path("/repo"), logger=silent_logger())

    def _capture_request(self, status_code=200, body=b'{"success":true}'):
        response = SimpleNamespace(
            getcode=lambda: status_code, read=lambda: body)
        opener = mock.MagicMock()
        opener.return_value = mock.MagicMock(
            __enter__=lambda self: response,
            __exit__=lambda self, *a: False)
        return opener

    def test_publish_puts_json_to_board_stats_endpoint_with_bearer(self):
        runner = self._runner()
        payload = {"totalComments": 3, "participants": []}
        opener = self._capture_request()
        with mock.patch.object(wcp.urllib.request, "urlopen", opener):
            runner._publish(payload)
        request = opener.call_args.args[0]
        self.assertEqual(request.get_method(), "PUT")
        self.assertTrue(
            request.full_url.endswith(
                "/api/jarvis/v1/board/stats/tf-weekly-comment-participation"),
            request.full_url)
        self.assertEqual(request.headers.get("Authorization"), "Bearer secret-token")
        self.assertEqual(request.headers.get("Content-type"), "application/json")
        sent = json.loads(request.data.decode("utf-8"))
        self.assertEqual(sent["totalComments"], 3)

    def test_publish_raises_on_http_error(self):
        runner = self._runner()
        err = wcp.urllib.error.HTTPError(
            "https://control-plane.test/api/jarvis/v1/board/stats/x",
            500, "Server Error", {}, io.BytesIO(b'{"err":1}'))
        with mock.patch.object(wcp.urllib.request, "urlopen", side_effect=err):
            with self.assertRaises(RuntimeError):
                runner._publish({"totalComments": 1})

    def test_publish_raises_when_base_url_missing(self):
        runner = wcp.WeeklyCommentParticipationRunner(
            task_client=FakeTaskClient(base_url=""), repo_root=Path("/repo"),
            logger=silent_logger())
        with self.assertRaises(RuntimeError):
            runner._publish({"totalComments": 1})


class WeeklyCommentParticipationRunTests(unittest.TestCase):
    def _runner(self) -> wcp.WeeklyCommentParticipationRunner:
        runner = wcp.WeeklyCommentParticipationRunner(
            task_client=FakeTaskClient(), repo_root=Path("/repo"),
            logger=silent_logger())
        runner._aggregate = lambda scheduled_for: {
            "totalComments": 0, "participants": [],
            "ticketsTouched": 0, "requirementsCovered": 0}
        runner._publish = lambda payload: None
        return runner

    def test_success_returns_succeeded(self):
        result = self._runner().run(definition(), datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))
        self.assertIs(result.status, JobResultStatus.SUCCEEDED)

    def test_invalid_slot_is_permanent_failure(self):
        runner = self._runner()
        wrong = ScheduledJobDefinition(
            "aone.scan", 1, "x", DailySchedule(6, 42, "Asia/Shanghai"),
            HandlerRunner("scan"), MisfirePolicy.CURRENT_DAY, 300, True)
        result = runner.run(wrong, datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))
        self.assertIs(result.status, JobResultStatus.PERMANENT_FAILURE)
        # naive scheduled_for is also an invalid slot
        result = runner.run(definition(), datetime(2026, 7, 27, 6, 42))
        self.assertIs(result.status, JobResultStatus.PERMANENT_FAILURE)

    def test_publish_failure_is_retryable(self):
        runner = self._runner()
        runner._publish = lambda payload: (_ for _ in ()).throw(
            RuntimeError("board stats PUT HTTP 500"))
        result = runner.run(definition(), datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))
        self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)


if __name__ == "__main__":
    unittest.main()
