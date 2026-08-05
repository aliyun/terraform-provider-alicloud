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

    def setUp(self):
        # _aggregate now parallel-prefetches comment lists per pool; stub it so
        # these tests (which mock _list_comments directly) don't spawn real a1.
        patch = mock.patch.object(wcp, "parallel_a1_per_id", return_value={})
        patch.start()
        self.addCleanup(patch.stop)

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
            jarvis = wcp._classify_author("jarvis", "已处理")
            self.assertEqual(jarvis[0], "digital")
            self.assertTrue(jarvis[3])
            chinese = wcp._classify_author("Terraform-研发数字人", "已处理")
            self.assertEqual(chinese[0], "digital")
            self.assertTrue(chinese[3])


class DeliveryMetricsAggregationTests(unittest.TestCase):
    def test_agent_intervention_ignores_claim_and_tag_activities(self):
        created = datetime(2026, 7, 1, tzinfo=SHANGHAI).timestamp()
        closed = datetime(2026, 7, 3, tzinfo=SHANGHAI).timestamp()
        activities = [
            {
                "operator": "jarvis",
                "property": "状态",
                "newValue": "jarvis-claim claimed #1",
                "eventTime": "2026-07-01 08:00",
            },
            {
                "operator": "Terraform-研发数字人",
                "property": "标签",
                "newValue": "jarvis-claimed",
                "eventTime": "2026-07-01 08:01",
            },
            {
                "operator": "Kelude",
                "property": "状态",
                "newValue": "自动流转",
                "eventTime": "2026-07-01 08:02",
            },
        ]

        self.assertIsNone(wcp._first_agent_intervention_at(
            [], activities, created, closed))

    def setUp(self):
        patch = mock.patch.object(wcp, "parallel_a1_per_id", return_value={})
        patch.start()
        self.addCleanup(patch.stop)

    def _runner(self) -> wcp.WeeklyCommentParticipationRunner:
        return wcp.WeeklyCommentParticipationRunner(
            task_client=FakeTaskClient(), repo_root=Path("/repo"),
            logger=silent_logger())

    def test_status_mapping_never_treats_solution_or_excluded_as_closed(self):
        self.assertEqual(
            wcp._classify_delivery_status("tf_customer", "已合入主线"), "solution")
        self.assertEqual(
            wcp._classify_delivery_status("tf_customer", "验收通过"), "closed")
        self.assertEqual(
            wcp._classify_delivery_status("tf_customer", "客户未响应"), "excluded")
        self.assertEqual(
            wcp._classify_delivery_status("tf_provider", "待发布"), "solution")
        self.assertEqual(
            wcp._classify_delivery_status("tf_provider", "已发布"), "closed")
        self.assertEqual(
            wcp._classify_delivery_status("tf_provider", "长期跟进"), "external_wait")
        self.assertEqual(
            wcp._classify_delivery_status("tf_provider", "开发中"), "open")

    def test_comment_parser_maps_aone_no_comments_text_to_empty_list(self):
        self.assertEqual(wcp._parse_comment_list("No comments found\n"), [])
        self.assertIsNone(wcp._parse_comment_list("not valid json"))

    def test_delivery_aggregate_uses_exact_close_activity_and_all_comments(self):
        runner = self._runner()
        runner._tf_pools = lambda: [
            ("tf_customer", "1086837", "需求问题"),
            ("tf_provider", "528766", "产品类需求"),
        ]
        rows = {
            "1086837": [
                {
                    "id": "1001", "title": "客户需求已验收", "status": "验收通过",
                    "createdAt": "2026-07-01 10:00:00",
                    "modified": datetime(2026, 7, 25, tzinfo=SHANGHAI).timestamp(),
                },
                {
                    "id": "1002", "title": "已提供方案", "status": "已合入主线",
                    "createdAt": "2026-07-02 10:00:00",
                    "modified": datetime(2026, 7, 26, tzinfo=SHANGHAI).timestamp(),
                },
            ],
            "528766": [
                {
                    "id": "2001", "title": "待正式发布", "status": "待发布",
                    "createdAt": "2026-07-03 10:00:00",
                    "modified": datetime(2026, 7, 26, tzinfo=SHANGHAI).timestamp(),
                },
                {
                    "id": "2002", "title": "已发布资源", "status": "已发布",
                    "createdAt": "2026-07-10 10:00:00",
                    "modified": datetime(2026, 7, 26, tzinfo=SHANGHAI).timestamp(),
                },
            ],
        }
        runner._list_closed_requirements = (
            lambda project, req_type, statuses, ws: rows[project])
        activities = {
            "1001": [
                {"property": "状态", "newValue": "验收通过",
                 "eventTime": "2026-07-21 10:00:00"},
            ],
            "2002": [
                {"fieldName": "status", "toValue": {"displayValue": "已发布"},
                 "createdAt": "2026-07-22 22:00:00"},
            ],
        }
        comments = {
            "1001": [
                {"author": "夏节", "content": "人工评论一",
                 "createdAt": "2026-07-05 10:00:00"},
                {"author": "terraform-rd", "content": "数字人评论",
                 "createdAt": "2026-07-06 10:00:00"},
                {"author": "kelude", "content": "状态流转",
                 "createdAt": "2026-07-21 10:00:00"},
            ],
            "2002": [
                {"author": "过载", "content": "人工评论二",
                 "createdAt": "2026-07-20 10:00:00"},
                {"author": "过载", "content": "人工评论三",
                 "createdAt": "2026-07-21 10:00:00"},
            ],
        }
        runner._list_activities = lambda iid: activities.get(iid, [])
        runner._list_comments = lambda iid: comments.get(iid, [])

        snapshot = runner._aggregate_delivery(
            datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))

        self.assertEqual(snapshot["coverageStart"], "2023-04-01T00:00:00+08:00")
        self.assertEqual(snapshot["cohortBasis"], "created_at")
        self.assertEqual(snapshot["scope"], "closed_success_only")
        self.assertEqual(snapshot["durationBasis"], "calendar_days")
        self.assertEqual(snapshot["closedCount"], 2)
        by_id = {item["id"]: item for item in snapshot["workitems"]}
        self.assertEqual(set(by_id), {"1001", "2002"})
        self.assertEqual(by_id["1001"]["statusClass"], "closed")
        self.assertEqual(by_id["1001"]["closedAt"],
                         "2026-07-21T10:00:00+08:00")
        self.assertEqual(by_id["1001"]["deliveryDays"], 20.0)
        self.assertEqual(by_id["1001"]["firstAgentInterventionAt"],
                         "2026-07-06T10:00:00+08:00")
        self.assertEqual(by_id["1001"]["agentInterventionElapsedDays"], 15.0)
        self.assertEqual(by_id["1001"]["humanCommentCount"], 1)
        self.assertEqual(by_id["1001"]["digitalCommentCount"], 1)
        self.assertEqual(by_id["1001"]["systemCommentCount"], 1)
        self.assertEqual(by_id["1001"]["totalCommentCount"], 3)
        self.assertNotIn("1002", by_id)
        self.assertNotIn("2001", by_id)  # 待发布是 solution，不是最终关单

        summaries = {row["pool"]: row for row in snapshot["pools"]}
        self.assertEqual(summaries["tf_customer"]["closedCount"], 1)
        self.assertEqual(summaries["tf_customer"]["averageDeliveryDays"], 20.0)
        self.assertEqual(summaries["tf_customer"]["humanCommentTotal"], 1)
        self.assertEqual(summaries["tf_provider"]["closedCount"], 1)
        self.assertEqual(summaries["tf_provider"]["humanCommentTotal"], 2)
        self.assertEqual(summaries["tf_provider"]["humanCommentAverage"], 2.0)

    def test_delivery_aggregate_marks_partial_when_item_has_no_exact_status_activity(self):
        runner = self._runner()
        runner._tf_pools = lambda: [
            ("tf_provider", "528766", "产品类需求")]
        runner._list_closed_requirements = lambda project, req_type, statuses, ws: [{
            "id": "2002", "title": "已发布资源", "status": "已发布",
            "createdAt": "2026-07-10 10:00:00",
            "modified": datetime(2026, 7, 26, tzinfo=SHANGHAI).timestamp(),
        }]
        runner._list_activities = lambda iid: []
        runner._list_comments = lambda iid: []

        snapshot = runner._aggregate_delivery(
            datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))

        self.assertFalse(snapshot["complete"])
        self.assertEqual(snapshot["failedWorkitemIds"], ["2002"])
        self.assertEqual(snapshot["failedWorkitemCount"], 1)
        self.assertEqual(snapshot["closedCount"], 0)

    def test_zero_comment_closed_item_stays_in_snapshot_with_zero_humans(self):
        runner = self._runner()
        runner._tf_pools = lambda: [
            ("tf_customer", "1086837", "需求问题")]
        runner._list_closed_requirements = (
            lambda project, req_type, statuses, ws: [{
                "id": "1003", "title": "无评论关单", "status": "验收通过",
                "createdAt": "2026-07-20 10:00:00",
            }])
        runner._list_activities = lambda iid: [{
            "property": "状态", "newValue": "验收通过",
            "eventTime": "2026-07-21 10:00:00",
        }]
        runner._list_comments = lambda iid: []

        snapshot = runner._aggregate_delivery(
            datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))

        self.assertTrue(snapshot["complete"])
        self.assertEqual(snapshot["closedCount"], 1)
        self.assertEqual(snapshot["workitems"][0]["humanCommentCount"], 0)
        self.assertEqual(snapshot["workitems"][0]["totalCommentCount"], 0)
        self.assertIsNone(snapshot["workitems"][0]["firstAgentInterventionAt"])
        self.assertIsNone(snapshot["workitems"][0]["agentInterventionElapsedDays"])

    def test_delivery_candidate_scan_starts_at_fy24_coverage_boundary(self):
        runner = self._runner()
        runner._tf_pools = lambda: [
            ("tf_customer", "1086837", "需求问题")]
        candidate_starts = []

        def list_candidates(project, req_type, statuses, window_start_epoch):
            candidate_starts.append(window_start_epoch)
            return []

        runner._list_closed_requirements = list_candidates

        runner._aggregate_delivery(
            datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))

        self.assertEqual(
            candidate_starts,
            [datetime(2023, 4, 1, 0, 0, tzinfo=SHANGHAI).timestamp()])

    def test_closed_candidate_scan_is_created_at_sorted_and_stops_before_fy24(self):
        runner = self._runner()
        recent = {
            "identifier": "1001", "subject": "FY24 item", "status": "验收通过",
            "gmtCreate": "2023-04-01 00:00:00",
        }
        old = {
            "identifier": "999", "subject": "old item", "status": "验收通过",
            "gmtCreate": "2023-03-31 23:59:59",
        }
        result = SimpleNamespace(
            returncode=0, stdout=json.dumps([recent, old], ensure_ascii=False),
            stderr="")
        with mock.patch.object(wcp, "run_process_group", return_value=result) as run:
            rows = runner._list_closed_requirements(
                "1086837", "需求问题", ("验收通过",),
                datetime(2023, 4, 1, tzinfo=SHANGHAI).timestamp())

        self.assertEqual([row["id"] for row in rows], ["1001"])
        command = run.call_args.args[0]
        self.assertIn("gmtCreate:desc", command)
        self.assertNotIn("finishTime:desc", command)


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

    def test_publish_delivery_writes_pages_then_commits_snapshot(self):
        runner = self._runner()
        snapshot = {
            "snapshotId": "20260727T064200+0800",
            "generatedAt": "2026-07-27T06:42:00+08:00",
            "windowStart": "2023-04-01T00:00:00+08:00",
            "windowEnd": "2026-07-27T06:42:00+08:00",
            "coverageStart": "2023-04-01T00:00:00+08:00",
            "coverageEnd": "2026-07-27T06:42:00+08:00",
            "cohortBasis": "created_at",
            "scope": "closed_success_only",
            "durationBasis": "calendar_days",
            "closedCount": 3,
            "pools": [],
            "workitems": [{"id": str(index), "title": "x" * 20}
                          for index in range(5)],
        }
        captured = []

        def open_request(request, timeout):
            captured.append(request)
            response = SimpleNamespace(getcode=lambda: 200, read=lambda: b"{}")
            return mock.MagicMock(
                __enter__=lambda self: response,
                __exit__=lambda self, *args: False)

        with mock.patch.object(wcp, "DELIVERY_PAGE_ITEMS", 2), \
             mock.patch.object(wcp.urllib.request, "urlopen",
                               side_effect=open_request):
            runner._publish_delivery(snapshot)

        self.assertEqual(len(captured), 4)
        self.assertEqual([request.get_method() for request in captured],
                         ["PUT", "PUT", "PUT", "POST"])
        self.assertTrue(captured[0].full_url.endswith(
            "/api/jarvis/v1/board/delivery-metrics/snapshots/"
            "20260727T064200+0800/pages/0"))
        self.assertTrue(captured[-1].full_url.endswith(
            "/api/jarvis/v1/board/delivery-metrics/snapshots/"
            "20260727T064200+0800/commit"))
        self.assertNotIn("%", captured[0].full_url)
        first_page = json.loads(captured[0].data.decode("utf-8"))
        self.assertEqual(first_page["items"][0]["id"], "0")
        commit = json.loads(captured[-1].data.decode("utf-8"))
        self.assertEqual(commit["pageCount"], 3)
        self.assertEqual(commit["itemCount"], 5)
        self.assertNotIn("workitems", commit)


class WeeklyCommentParticipationRunTests(unittest.TestCase):
    def _runner(self) -> wcp.WeeklyCommentParticipationRunner:
        runner = wcp.WeeklyCommentParticipationRunner(
            task_client=FakeTaskClient(), repo_root=Path("/repo"),
            logger=silent_logger())
        runner._aggregate = lambda scheduled_for: {
            "totalComments": 0, "participants": [],
            "ticketsTouched": 0, "requirementsCovered": 0}
        runner._publish = lambda payload: None
        runner._aggregate_delivery = lambda scheduled_for: {
            "snapshotId": "20260727T064200+0800",
            "closedCount": 0, "pools": [], "workitems": []}
        runner._publish_delivery = lambda payload: None
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

    def test_delivery_failure_does_not_block_legacy_weekly_publish(self):
        runner = self._runner()
        published = []
        runner._publish = published.append
        runner._aggregate_delivery = lambda scheduled_for: (
            (_ for _ in ()).throw(RuntimeError("pool list failed")))

        result = runner.run(
            definition(), datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI))

        self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
        self.assertEqual(len(published), 1)


if __name__ == "__main__":
    unittest.main()
