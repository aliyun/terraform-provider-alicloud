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


def _item(iid, title="t"):
    return {"id": str(iid), "title": title, "status": "New",
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
        self.assertTrue(any("ak.issue.member=" in f for f in filters))
        self.assertTrue(any("workitem.tracker=" in f for f in filters))

    def test_union_filters_participants_uses_full_digital_worker_csv(self):
        runner = self._runner()
        filters = runner._union_filters(exclude_status=["Closed"])
        member_filter = next(f for f in filters if "ak.issue.member=" in f)
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


if __name__ == "__main__":
    unittest.main()
