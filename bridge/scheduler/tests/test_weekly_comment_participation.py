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


def definition() -> ScheduledJobDefinition:
    return ScheduledJobDefinition(
        wcp.JOB_KEY, 1, "weekly comment participation",
        DailySchedule(6, 42, "Asia/Shanghai"), HandlerRunner(wcp.RUNNER_KEY),
        MisfirePolicy.CURRENT_DAY, 300, True,
    )


class FakeTaskClient:
    def __init__(self, base_url="https://control-plane.test", token="machine"):
        self.base_url = base_url
        self.token = token
        self.timeout = 5.0


def runner() -> wcp.WeeklyCommentParticipationRunner:
    instance = wcp.WeeklyCommentParticipationRunner(
        task_client=FakeTaskClient(), repo_root=Path("/repo"),
        logger=logging.getLogger("test-weekly-comment"))
    instance._tf_pools = lambda: [
        ("tf_customer", "1086837", "需求问题"),
        ("tf_provider", "528766", "产品类需求"),
    ]
    instance._delivery_status_definitions = lambda: wcp._DELIVERY_STATUS_DEFINITIONS
    return instance


NOW = datetime(2026, 7, 27, 6, 42, tzinfo=SHANGHAI)


def source(item_id: str, *, pool="1086837", status="开发中",
           created="2026-07-01 10:00:00", modified=None, title="item"):
    return {
        "id": item_id, "project": pool, "title": title, "type": "需求问题",
        "status": status, "createdAt": created,
        "modified": (modified or (NOW - timedelta(hours=1))).timestamp(),
    }


def previous_snapshot(*items, cursor="2026-07-26T06:42:00+08:00"):
    return {
        "manifest": {
            "snapshotId": "old", "poolCursors": {
                "tf_customer": cursor, "tf_provider": cursor,
            },
        },
        "workitems": list(items),
    }


class IncrementalCollectionTests(unittest.TestCase):
    def test_bootstrap_keeps_all_lifecycle_facts_but_enriches_only_recent(self):
        instance = runner()
        old = source("1000", modified=NOW - timedelta(days=30))
        recent = source("1001", status="验收通过")
        instance._list_incremental_requirements = mock.Mock(
            side_effect=[[old, recent], []])
        instance._list_comments = mock.Mock(return_value=[])
        instance._list_activities = mock.Mock(return_value=[{
            "property": "状态", "newValue": "验收通过",
            "eventTime": "2026-07-26 10:00:00",
        }])

        result = instance._aggregate_delivery(NOW, None)

        self.assertEqual(result["scope"], "lifecycle_fact_cache")
        self.assertEqual(result["collectionMode"], "bootstrap")
        self.assertEqual({row["id"] for row in result["workitems"]}, {"1000", "1001"})
        self.assertEqual(result["collectionStats"], {
            "listCalls": 2, "commentCalls": 1, "activityCalls": 1,
            "changedWorkitemCount": 2,
        })
        by_id = {row["id"]: row for row in result["workitems"]}
        self.assertEqual(by_id["1000"]["commentFreshness"], "unknown")
        self.assertEqual(by_id["1000"]["statusClass"], "open")
        self.assertEqual(by_id["1001"]["closedAt"], "2026-07-26T10:00:00+08:00")

    def test_incremental_budget_is_two_lists_plus_changed_comments_and_transitions(self):
        instance = runner()
        unchanged = {
            **source("900", modified=NOW - timedelta(days=20)), "pool": "tf_customer",
            "poolId": "1086837", "statusClass": "open", "modifiedAt":
            "2026-07-01T10:00:00+08:00", "commentEvents": [],
        }
        changed_open = source("1001")
        changed_closed = source("2001", pool="528766", status="已发布")
        cutoffs = []

        def list_changed(project, req_type, modified_after, coverage_start):
            cutoffs.append(modified_after)
            return [changed_open] if project == "1086837" else [changed_closed]

        instance._list_incremental_requirements = mock.Mock(side_effect=list_changed)
        instance._list_comments = mock.Mock(return_value=[])
        instance._list_activities = mock.Mock(return_value=[{
            "fieldName": "status", "toValue": "已发布",
            "createdAt": "2026-07-27 05:00:00",
        }])

        result = instance._aggregate_delivery(NOW, previous_snapshot(unchanged))

        self.assertEqual(result["collectionMode"], "incremental")
        self.assertEqual(result["collectionStats"]["listCalls"], 2)
        self.assertEqual(result["collectionStats"]["commentCalls"], 2)
        self.assertEqual(result["collectionStats"]["activityCalls"], 1)
        expected = datetime(2026, 7, 26, 6, 32, tzinfo=SHANGHAI).timestamp()
        self.assertEqual(cutoffs, [expected, expected])
        self.assertIn("900", {row["id"] for row in result["workitems"]})

    def test_reclosed_item_refreshes_to_latest_terminal_transition(self):
        instance = runner()
        prior = {
            **source("2001", pool="528766", status="已发布"),
            "pool": "tf_provider", "poolId": "528766", "statusClass": "closed",
            "modifiedAt": "2026-07-20T10:00:00+08:00",
            "closedAt": "2026-07-10T10:00:00+08:00", "commentEvents": [],
        }
        instance._list_incremental_requirements = mock.Mock(
            side_effect=[[], [source("2001", pool="528766", status="已发布")]])
        instance._list_comments = mock.Mock(return_value=[])
        instance._list_activities = mock.Mock(return_value=[
            {"property": "状态", "newValue": "已发布", "eventTime": "2026-07-10 10:00:00"},
            {"property": "状态", "newValue": "已发布", "eventTime": "2026-07-26 10:00:00"},
        ])

        result = instance._aggregate_delivery(NOW, previous_snapshot(prior))
        item = next(row for row in result["workitems"] if row["id"] == "2001")
        self.assertEqual(item["closedAt"], "2026-07-26T10:00:00+08:00")
        self.assertEqual(instance._list_activities.call_count, 1)

    def test_comment_failure_retains_cache_and_marks_stale(self):
        instance = runner()
        event = {"createdAt": "2026-07-25T10:00:00+08:00",
                 "authorToken": "1", "authorName": "人", "digital": False}
        prior = {
            **source("1001"), "pool": "tf_customer", "poolId": "1086837",
            "statusClass": "open", "modifiedAt": "2026-07-25T10:00:00+08:00",
            "commentEvents": [event], "humanCommentCount": 1,
            "digitalCommentCount": 0, "systemCommentCount": 0,
            "totalCommentCount": 1, "commentFreshness": "fresh",
        }
        instance._list_incremental_requirements = mock.Mock(
            side_effect=[[source("1001")], []])
        instance._list_comments = mock.Mock(return_value=None)

        result = instance._aggregate_delivery(NOW, previous_snapshot(prior))
        item = next(row for row in result["workitems"] if row["id"] == "1001")
        self.assertEqual(item["commentEvents"], [event])
        self.assertEqual(item["humanCommentCount"], 1)
        self.assertEqual(item["commentFreshness"], "stale")
        self.assertEqual(result["freshness"]["staleCommentWorkitemCount"], 1)

    def test_legacy_upgrade_clears_all_aone_inferred_agent_metrics(self):
        instance = runner()
        prior = {
            **source("1001", modified=NOW - timedelta(days=20)),
            "pool": "tf_customer", "poolId": "1086837",
            "statusClass": "open", "modifiedAt": "2026-07-01T10:00:00+08:00",
            "commentEvents": [], "firstAgentInterventionAt":
            "2026-07-02T10:00:00+08:00", "agentInterventionElapsedDays": 4.5,
            "agentParticipated": True, "agentActiveDays": 3,
        }
        instance._list_incremental_requirements = mock.Mock(side_effect=[[], []])

        result = instance._aggregate_delivery(NOW, previous_snapshot(prior))

        item = next(row for row in result["workitems"] if row["id"] == "1001")
        self.assertIsNone(item["firstAgentInterventionAt"])
        self.assertIsNone(item["agentInterventionElapsedDays"])
        self.assertNotIn("agentParticipated", item)
        self.assertNotIn("agentActiveDays", item)

    def test_activity_failure_excludes_close_and_is_retried_without_comment_call(self):
        instance = runner()
        closing = source("2001", pool="528766", status="已发布")
        instance._list_incremental_requirements = mock.Mock(side_effect=[[], [closing]])
        instance._list_comments = mock.Mock(return_value=[])
        instance._list_activities = mock.Mock(return_value=None)

        first = instance._aggregate_delivery(NOW, previous_snapshot())
        item = next(row for row in first["workitems"] if row["id"] == "2001")
        self.assertIsNone(item["closedAt"])
        self.assertEqual(item["pendingTransition"], "closedAt")
        self.assertEqual(first["closedCount"], 0)

        instance._list_incremental_requirements = mock.Mock(side_effect=[[], []])
        instance._list_comments.reset_mock()
        instance._list_activities = mock.Mock(return_value=[{
            "property": "状态", "newValue": "已发布",
            "eventTime": "2026-07-27 05:00:00",
        }])
        second = instance._aggregate_delivery(NOW + timedelta(days=1), first)
        retried = next(row for row in second["workitems"] if row["id"] == "2001")
        self.assertIsNotNone(retried["closedAt"])
        self.assertNotIn("pendingTransition", retried)
        instance._list_comments.assert_not_called()
        self.assertEqual(instance._list_activities.call_count, 1)

    def test_list_failure_fails_closed_before_any_publish(self):
        instance = runner()
        instance._load_delivery_snapshot = lambda: previous_snapshot()
        instance._list_incremental_requirements = mock.Mock(side_effect=[[], None])
        instance._publish_delivery = mock.Mock()
        instance._publish = mock.Mock()

        result = instance.run(definition(), NOW)

        self.assertIs(result.status, JobResultStatus.RETRYABLE_FAILURE)
        instance._publish_delivery.assert_not_called()
        instance._publish.assert_not_called()

    def test_list_uses_exact_modified_filter_and_counts_real_pages(self):
        instance = runner()
        response = SimpleNamespace(returncode=0, stdout="[]", stderr="")
        with mock.patch.object(wcp, "run_process_group", return_value=response) as run:
            result = instance._list_incremental_requirements(
                "1086837", "需求问题",
                datetime(2026, 7, 26, 6, 32, tzinfo=SHANGHAI).timestamp(),
                datetime(2023, 4, 1, tzinfo=SHANGHAI).timestamp())
        self.assertEqual(result, [])
        argv = run.call_args.args[0]
        index = argv.index("--filter")
        self.assertEqual(
            argv[index + 1],
            "modified>'2026-07-26 06:32:00' AND "
            "created>='2023-04-01 00:00:00'")
        self.assertEqual(instance._last_list_pages, 1)

    def test_bootstrap_list_pushes_fy24_coverage_into_filter(self):
        instance = runner()
        response = SimpleNamespace(returncode=0, stdout="[]", stderr="")
        with mock.patch.object(wcp, "run_process_group", return_value=response) as run:
            result = instance._list_incremental_requirements(
                "1086837", "需求问题", None,
                datetime(2023, 4, 1, tzinfo=SHANGHAI).timestamp())
        self.assertEqual(result, [])
        argv = run.call_args.args[0]
        index = argv.index("--filter")
        self.assertEqual(argv[index + 1], "created>='2023-04-01 00:00:00'")


class ProjectionAndTransportTests(unittest.TestCase):
    def test_metrics_publish_defaults_to_preprod_without_mutating_task_client(self):
        client = FakeTaskClient(base_url="https://agent.aliyun-inc.com")
        instance = wcp.WeeklyCommentParticipationRunner(
            task_client=client, repo_root=Path("/repo"),
            logger=logging.getLogger("test-weekly-comment"), environ={})
        response = SimpleNamespace(getcode=lambda: 200, read=lambda: b"{}")
        opener = mock.MagicMock()
        opener.return_value = mock.MagicMock(
            __enter__=lambda self: response, __exit__=lambda self, *args: False)

        with mock.patch.object(wcp.urllib.request, "urlopen", opener):
            instance._publish({"totalComments": 0, "participants": []})

        self.assertEqual(opener.call_args.args[0].full_url,
                         "https://pre-agent.aliyun-inc.com/api/jarvis/v1/board/"
                         "stats/tf-weekly-comment-participation")
        self.assertEqual(client.base_url, "https://agent.aliyun-inc.com")

    def test_metrics_endpoint_has_dedicated_config_override(self):
        instance = wcp.WeeklyCommentParticipationRunner(
            task_client=FakeTaskClient(base_url="https://agent.aliyun-inc.com"),
            repo_root=Path("/repo"), logger=logging.getLogger("test-weekly-comment"),
            environ={"JARVIS_METRICS_BASE_URL": "https://metrics.example/"})
        response = SimpleNamespace(getcode=lambda: 200, read=lambda: b"{}")
        opener = mock.MagicMock()
        opener.return_value = mock.MagicMock(
            __enter__=lambda self: response, __exit__=lambda self, *args: False)

        with mock.patch.object(wcp.urllib.request, "urlopen", opener):
            instance._get_json("/metrics-test")

        self.assertEqual(opener.call_args.args[0].full_url,
                         "https://metrics.example/metrics-test")

    def test_missing_metrics_base_uses_metrics_specific_error(self):
        instance = runner()
        instance._metrics_base_url = ""
        with self.assertRaisesRegex(RuntimeError, "metrics base_url"):
            instance._publish({"totalComments": 0, "participants": []})

    def test_weekly_projection_uses_cached_comment_events_only(self):
        instance = runner()
        instance._list_active_requirements = mock.Mock(side_effect=AssertionError)
        snapshot = {"workitems": [{
            "id": "1", "commentEvents": [
                {"createdAt": "2026-07-25T10:00:00+08:00", "authorToken": "h1",
                 "authorName": "夏节", "digital": False},
                {"createdAt": "2026-07-24T10:00:00+08:00", "authorToken": "rd",
                 "authorName": "Terraform RD", "digital": True},
                {"createdAt": "2026-07-01T10:00:00+08:00", "authorToken": "old",
                 "authorName": "旧", "digital": False},
            ],
        }]}
        payload = instance._aggregate(NOW, snapshot)
        self.assertEqual(payload["totalComments"], 2)
        self.assertEqual(payload["ticketsTouched"], 1)
        self.assertEqual({row["name"] for row in payload["participants"]},
                         {"夏节", "Terraform RD"})
        instance._list_active_requirements.assert_not_called()

    def test_current_snapshot_full_contract_and_machine_bearer(self):
        instance = runner()
        response = SimpleNamespace(getcode=lambda: 200, read=lambda: json.dumps({
            "manifest": {"snapshotId": "s1", "poolCursors": {}},
            "items": [{"id": "1"}],
        }).encode())
        opener = mock.MagicMock()
        opener.return_value = mock.MagicMock(
            __enter__=lambda self: response, __exit__=lambda self, *args: False)
        with mock.patch.object(wcp.urllib.request, "urlopen", opener):
            current = instance._load_delivery_snapshot()
        request = opener.call_args.args[0]
        self.assertEqual(request.get_method(), "GET")
        self.assertEqual(request.full_url,
            "https://pre-agent.aliyun-inc.com/api/jarvis/v1/board/"
            "delivery-metrics/snapshots/current")
        self.assertEqual(request.headers["Authorization"], "Bearer machine")
        self.assertEqual(current["workitems"], [{"id": "1"}])

    def test_current_snapshot_404_means_bootstrap_but_409_fails_closed(self):
        instance = runner()
        not_found = wcp.urllib.error.HTTPError("x", 404, "", {}, io.BytesIO())
        with mock.patch.object(wcp.urllib.request, "urlopen", side_effect=not_found):
            self.assertIsNone(instance._load_delivery_snapshot())
        conflict = wcp.urllib.error.HTTPError("x", 409, "", {}, io.BytesIO(b"bad"))
        with mock.patch.object(wcp.urllib.request, "urlopen", side_effect=conflict):
            with self.assertRaises(RuntimeError):
                instance._load_delivery_snapshot()

    def test_publish_delivery_keeps_manifest_fields_and_chunks_items(self):
        instance = runner()
        snapshot = {
            "snapshotId": "20260727T064200+0800", "scope": "lifecycle_fact_cache",
            "poolCursors": {"tf_customer": "x"}, "collectionMode": "incremental",
            "collectionStats": {"listCalls": 2}, "freshness": {"collectedAt": "x"},
            "workitems": [{"id": str(index)} for index in range(3)],
        }
        captured = []

        def urlopen(request, timeout):
            captured.append(request)
            response = SimpleNamespace(getcode=lambda: 200, read=lambda: b"{}")
            return mock.MagicMock(
                __enter__=lambda self: response, __exit__=lambda self, *args: False)

        with mock.patch.object(wcp, "DELIVERY_PAGE_ITEMS", 2), \
             mock.patch.object(wcp.urllib.request, "urlopen", side_effect=urlopen):
            instance._publish_delivery(snapshot)
        self.assertEqual([req.get_method() for req in captured], ["PUT", "PUT", "POST"])
        commit = json.loads(captured[-1].data)
        self.assertEqual(commit["scope"], "lifecycle_fact_cache")
        self.assertEqual(commit["collectionStats"]["listCalls"], 2)
        self.assertEqual(commit["pageCount"], 2)
        self.assertNotIn("workitems", commit)


class SmallHelperTests(unittest.TestCase):
    def test_status_and_comment_parsers(self):
        self.assertEqual(wcp._classify_delivery_status(
            "tf_customer", "已发布待需求方验收"), "solution")
        self.assertEqual(wcp._classify_delivery_status(
            "tf_provider", "已发布"), "closed")
        self.assertEqual(wcp._parse_comment_list("No comments found\n"), [])
        self.assertIsNone(wcp._parse_comment_list("bad-json"))

    def test_exact_close_uses_latest_transition(self):
        activities = [
            {"property": "状态", "newValue": "已发布", "eventTime": "2026-07-01"},
            {"property": "状态", "newValue": "已发布", "eventTime": "2026-07-03"},
        ]
        value = wcp._exact_closed_at(activities, ("已发布",))
        self.assertEqual(value, datetime(2026, 7, 3, tzinfo=SHANGHAI).timestamp())


if __name__ == "__main__":
    unittest.main()
