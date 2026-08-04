"""Tests for ScanRunner pool-union discovery and dedup.

Covers the ak.issue.member (participants) source added alongside workitem.tracker
(抄送) so tickets discoverable only via participants are no longer missed when the
tracker query fails, and that overlapping results across sources deduplicate by id.
"""
import os
import unittest
from unittest import mock

from bridge.scheduler.runners import scan
from bridge.aone_tasks import DIGITAL_WORKER_IDS


def _item(iid, title="t", status="New"):
    return {"id": str(iid), "title": title, "status": status,
            "priority": "", "tag": "", "type": "", "category": "",
            "modified": "", "created": "", "assignedTo": ""}


class UnionFiltersTests(unittest.TestCase):
    def _runner(self):
        with mock.patch.object(scan, "PendingDispatchRegistry"):
            return scan.ScanRunner(
                logger=mock.Mock(), task_client=mock.Mock(),
                repo_root=os.getcwd())

    def test_union_filters_include_participants_source(self):
        runner = self._runner()
        filters = runner._union_filters(exclude_status=["Closed"])
        # Five sources, and ak.issue.member present to backstop the flaky tracker.
        self.assertEqual(len(filters), 5)
        self.assertTrue(any("ak.issue.member=" in expr for expr, _excl in filters))
        self.assertTrue(any("workitem.tracker=" in expr for expr, _excl in filters))

    def test_union_filters_participants_uses_full_digital_worker_csv(self):
        runner = self._runner()
        filters = runner._union_filters(exclude_status=["Closed"])
        member_filter = next(expr for expr, _excl in filters
                             if "ak.issue.member=" in expr)
        for wid in DIGITAL_WORKER_IDS:
            self.assertIn(wid, member_filter)


class PoolUnionDedupTests(unittest.TestCase):
    def _runner(self, per_filter_results):
        with mock.patch.object(scan, "PendingDispatchRegistry"):
            runner = scan.ScanRunner(
                logger=mock.Mock(), task_client=mock.Mock(),
                repo_root=os.getcwd())

        def fake_a1_list(_cls, _project, filter_expr):
            for needle, results in per_filter_results.items():
                if needle in filter_expr:
                    return list(results)
            return []

        runner._a1_list = fake_a1_list.__get__(scan.ScanRunner)
        return runner

    def test_participants_and_tracker_dedup_by_id(self):
        # Ticket 83818138 appears in BOTH tracker and ak.issue.member results; ticket
        # 99900001 only in participants. The union must keep each once (no duplicates).
        runner = self._runner({
            "workitem.tracker=": [_item("83818138", "shared")],
            "ak.issue.member=": [_item("83818138", "shared"),
                                 _item("99900001", "only-participant")],
        })
        rows = runner._query_pool_union("mcp_server", "2124589", ["Closed"])
        ids = [r["id"] for r in rows]
        self.assertEqual(ids, ["83818138", "99900001"])
        self.assertEqual(len(rows), 2)  # 83818138 deduped, not 3

    def test_participants_source_alone_discovers_ticket_when_tracker_empty(self):
        # The real-world gap: tracker query fails/empty, but participants catches it.
        runner = self._runner({
            "ak.issue.member=": [_item("83818138", "ValidateModule")],
        })
        rows = runner._query_pool_union("mcp_server", "2124589", ["Closed"])
        self.assertEqual([r["id"] for r in rows], ["83818138"])
        self.assertEqual(rows[0]["pool"], "mcp_server")
        self.assertEqual(rows[0]["pool_project"], "2124589")


class UnionStatusFilterTests(unittest.TestCase):
    """状态排除下沉到客户端后的契约。

    Aone 的 workitem list 每多一个 ``AND NOT status=`` 子句约多花 3 秒（同池同页实测
    0/1/4/8/14 子句 = 4/5/15/27/38s）。每轮 5 池 × 5 源 = 25 个并发查询打同一接口，
    会把单条查询顶过 ``_a1_list`` 的 90s/页超时；该池当轮零行返回，池内所有工单既不
    判定也不告警。排除因此改为本地按集合过滤：结果集与服务端过滤等价，但不付那份延迟。
    """

    def _runner(self, per_filter_results=None):
        with mock.patch.object(scan, "PendingDispatchRegistry"):
            runner = scan.ScanRunner(
                logger=mock.Mock(), task_client=mock.Mock(),
                repo_root=os.getcwd())
        if per_filter_results is None:
            return runner

        def fake_a1_list(_cls, _project, filter_expr):
            for needle, results in per_filter_results.items():
                if needle in filter_expr:
                    return list(results)
            return []

        runner._a1_list = fake_a1_list.__get__(scan.ScanRunner)
        return runner

    # 与 config/pools.json tf_customer.pr_merged_status 同形。
    PR_MERGED = {"type": "3", "type_name": "需求问题",
                 "name": "已合入主线", "id": "626904"}

    def _specs(self):
        specs = self._runner()._union_filters(
            exclude_status=["Closed", "Fixed"], pr_merged_status=self.PR_MERGED)
        for spec in specs:
            self.assertIsInstance(
                spec, tuple, "每源应返回 (filter_expr, excluded_statuses) 二元组")
        return specs

    def test_union_filters_carry_no_status_clauses(self):
        for expr, _excluded in self._specs():
            self.assertNotIn("NOT status=", expr)

    def test_pool_exclusions_apply_to_every_discovery_source(self):
        for expr, excluded in self._specs():
            if expr == "tag=jarvis-done":
                continue
            self.assertEqual(excluded, frozenset({"Closed", "Fixed"}), expr)

    def test_done_source_excludes_only_the_pr_merged_status(self):
        # done watch 必须看得见处于终态的 done 单，否则发现不了「完成后的人工追评」。
        done = next(excluded for expr, excluded in self._specs()
                    if expr == "tag=jarvis-done")
        self.assertEqual(done, frozenset({"已合入主线"}))

    def test_terminal_ticket_from_assignee_source_is_dropped(self):
        runner = self._runner({"assignedTo=": [_item("500", status="Closed")]})
        rows = runner._query_pool_union("mcp_server", "2124589", ["Closed"])
        self.assertEqual(rows, [])

    def test_done_ticket_in_terminal_status_still_enters_union(self):
        # 同一条单被 assignee 源取到又丢弃，必须仍能经 done 源进入并集：丢弃绝不能
        # 触碰去重表，否则 new_comment_after_done 会整条失效。
        runner = self._runner({
            "assignedTo=": [_item("777", status="Fixed")],
            "tag=jarvis-done": [_item("777", status="Fixed")],
        })
        rows = runner._query_pool_union("mcp_server", "2124589", ["Fixed"])
        self.assertEqual([r["id"] for r in rows], ["777"])

    def test_empty_exclude_status_keeps_every_row(self):
        runner = self._runner({"assignedTo=": [_item("500", status="Closed")]})
        rows = runner._query_pool_union("mcp_server", "2124589", [])
        self.assertEqual([r["id"] for r in rows], ["500"])


if __name__ == "__main__":
    unittest.main()
