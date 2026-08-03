"""Tests for the durable discovery snapshot that survives Scheduler restarts.

The invariant under test is asymmetric, and that asymmetry is the whole design:

- A **missing** snapshot entry costs one re-evaluation, which reaches the same
  conclusion — the decision is a pure function of Aone state.
- A **spurious** entry skips re-evaluation and strands the ticket until its next
  Aone modification. That is a lost ticket.

So every test here pushes on the same question: did an outcome that never durably
landed get recorded as seen?
"""
import os
import tempfile
import unittest
from datetime import datetime, timedelta
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

from bridge.scan_snapshot import ScanSnapshotStore
from bridge.scheduler.runners import scan


class ScanSnapshotStoreTests(unittest.TestCase):
    def _store(self):
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        return ScanSnapshotStore(Path(tmp.name) / "scan-snapshot.json")

    def test_missing_file_loads_empty(self):
        self.assertEqual(self._store().load(), {})

    def test_commit_round_trips(self):
        store = self._store()
        store.commit({"83612870": "2026-08-03 16:56:39"})
        self.assertEqual(store.load(), {"83612870": "2026-08-03 16:56:39"})

    def test_commit_merges_previous_entries(self):
        store = self._store()
        store.commit({"1": "m1"})
        store.commit({"2": "m2"})
        self.assertEqual(store.load(), {"1": "m1", "2": "m2"})

    def test_commit_prunes_to_keep_set(self):
        # `keep` is the live union: entries that left it stop being tracked so the
        # file does not grow without bound.
        store = self._store()
        store.commit({"1": "m1", "2": "m2"})
        store.commit({"3": "m3"}, keep=["1", "3"])
        self.assertEqual(store.load(), {"1": "m1", "3": "m3"})

    def test_corrupt_file_degrades_to_cold_start(self):
        store = self._store()
        store.path.parent.mkdir(parents=True, exist_ok=True)
        store.path.write_text("{ not json", encoding="utf-8")
        # Everything looks new again (one expensive tick) rather than raising.
        self.assertEqual(store.load(), {})

    def test_entries_without_modified_are_dropped(self):
        # A blank `modified` cannot drive the diff, so it must not masquerade as a
        # recorded observation.
        store = self._store()
        store.commit({"1": "m1", "2": ""})
        self.assertEqual(store.load(), {"1": "m1"})

    def test_write_leaves_no_temp_files(self):
        store = self._store()
        store.commit({"1": "m1"})
        leftovers = sorted(p.name for p in store.path.parent.iterdir()
                           if p.name.endswith(".tmp"))
        self.assertEqual(leftovers, [])


class ScanRunnerHydrationTests(unittest.TestCase):
    def test_restart_hydrates_prev_snapshot_from_disk(self):
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        path = Path(tmp.name) / "scan-snapshot.json"
        ScanSnapshotStore(path).commit({"83612870": "2026-08-03 16:56:39"})
        with mock.patch.dict(os.environ,
                             {"JARVIS_SCAN_SNAPSHOT_PATH": str(path)}), \
                mock.patch.object(scan, "PendingDispatchRegistry"):
            runner = scan.ScanRunner(
                logger=mock.Mock(), task_client=mock.Mock(),
                repo_root=os.getcwd())
        # The point of the whole change: a restart must not replay this as new.
        self.assertIn("83612870", runner._prev_snapshot)
        self.assertEqual(
            runner._prev_snapshot["83612870"]["modified"],
            "2026-08-03 16:56:39")

    def test_missing_snapshot_starts_cold_without_raising(self):
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        path = Path(tmp.name) / "absent.json"
        with mock.patch.dict(os.environ,
                             {"JARVIS_SCAN_SNAPSHOT_PATH": str(path)}), \
                mock.patch.object(scan, "PendingDispatchRegistry"):
            runner = scan.ScanRunner(
                logger=mock.Mock(), task_client=mock.Mock(),
                repo_root=os.getcwd())
        self.assertEqual(runner._prev_snapshot, {})


class ConclusiveCommitTests(unittest.TestCase):
    """Only durably-landed outcomes may be persisted."""

    def _runner(self):
        r = scan.ScanRunner.__new__(scan.ScanRunner)
        r._prefetch_idle_claimed_a1 = mock.Mock()
        r.dispatch_pools = set()
        r.dispatch_created_before = ""
        r._pr_merged_status_by_pool = {}
        r._envelope = mock.Mock(return_value=SimpleNamespace())
        r.pool = None
        r._signal_query_retry = set()
        r._done_watch_retry = set()
        return r

    @staticmethod
    def _untagged_item():
        return {"id": "83612870", "tag": "", "status": "待处理", "title": "x",
                "modified": "2026-08-03 16:56:39"}

    def test_successful_dispatch_is_persisted(self):
        runner = self._runner()
        runner._dispatch = mock.Mock(return_value=(True, "new"))
        self.assertEqual(
            runner._tick_auto([self._untagged_item()]),
            {"83612870": "2026-08-03 16:56:39"})

    def test_failed_dispatch_is_not_persisted_so_a_restart_retries(self):
        # The lost-ticket guard. Today's in-memory snapshot advances before any
        # decision, so an enqueue failure waits for the next Aone modification to be
        # looked at again. Persisting that same point would make the gap survive
        # restarts; leaving it out means the next process sees the ticket as new.
        runner = self._runner()
        runner._dispatch = mock.Mock(return_value=(False, "enqueue_failed"))
        self.assertEqual(runner._tick_auto([self._untagged_item()]), {})

    def test_terminal_skip_is_persisted(self):
        runner = self._runner()
        runner._dispatch = mock.Mock(return_value=(True, "new"))
        item = dict(self._untagged_item(), status="已发布")
        conclusive = runner._tick_auto([item])
        self.assertIn("83612870", conclusive)
        runner._dispatch.assert_not_called()

    def test_claimed_failed_comment_read_is_not_persisted(self):
        # The claimed branch has the same shape as idle: _claimed_human_comment
        # returns None both when there is genuinely no new comment and when the read
        # failed. In-memory that was survivable because a restart re-evaluated
        # everything; persisting it would swallow a human comment across restarts.
        runner = self._runner()
        runner._dispatch = mock.Mock(return_value=(True, "new"))
        runner._claimed_human_comment = mock.Mock(return_value=None)
        runner._activity_cache = {}
        runner._human_comment_cache = {"claimed:84621304": False}  # 查询失败标记
        item = {"id": "84621304", "tag": "jarvis-claimed", "status": "评估中",
                "title": "y", "modified": "2026-08-03 16:56:39"}
        self.assertEqual(runner._tick_auto([item]), {})
        self.assertIn("84621304", runner._signal_query_retry)

    def test_claimed_clean_read_is_persisted(self):
        runner = self._runner()
        runner._dispatch = mock.Mock(return_value=(True, "new"))
        runner._claimed_human_comment = mock.Mock(return_value=None)
        runner._activity_cache = {}
        runner._human_comment_cache = {}
        item = {"id": "84621304", "tag": "jarvis-claimed", "status": "评估中",
                "title": "y", "modified": "2026-08-03 16:56:39"}
        self.assertIn("84621304", runner._tick_auto([item]))
        self.assertNotIn("84621304", runner._signal_query_retry)

    def test_lag_window_skip_is_not_persisted(self):
        runner = self._runner()
        runner._last_idle_at = mock.Mock(
            return_value=datetime.now() - timedelta(days=27))
        runner._human_comment = mock.Mock(return_value=None)
        runner._human_touched = mock.Mock(return_value=False)
        runner._activity_cache = {}
        runner._human_comment_cache = {}
        item = {"id": "84621304", "tag": "jarvis-idle", "status": "评估中",
                "title": "y",
                "modified": (datetime.now() - timedelta(minutes=5)).strftime(
                    "%Y-%m-%d %H:%M:%S")}
        self.assertEqual(runner._tick_auto([item]), {})


class TickCommitTests(unittest.TestCase):
    """_tick must actually reach disk.

    The commit is wrapped in `except Exception` so snapshot IO can never fail a
    tick — which also means a missing attribute or a bad argument degrades to a
    silent no-op that every other test still passes through. This asserts the write
    end-to-end instead of trusting the call site.
    """

    def _runner(self, path):
        r = scan.ScanRunner.__new__(scan.ScanRunner)
        r.snapshot_store = ScanSnapshotStore(path)
        r.auto = True
        r._prev_snapshot = {}
        r._prefetch_idle_claimed_a1 = mock.Mock()
        r._load_human_operators = lambda: set()
        r._human_cache = {}
        r._human_comment_cache = {}
        r._activity_cache = {}
        r._raw_comment_cache = {}
        r._done_watch_retry = set()
        r._done_watch_confirm = set()
        r._done_drift_retry = set()
        r._signal_query_retry = set()
        r.dispatch_pools = set()
        r.dispatch_created_before = ""
        r._pr_merged_status_by_pool = {}
        r.pool = None
        r._envelope = mock.Mock(return_value=SimpleNamespace())
        r._reconcile_done_status_drifts_safely = mock.Mock()
        r._reconcile_source_statuses_safely = mock.Mock()
        return r

    @staticmethod
    def _union_item(**over):
        item = {"id": "83612870", "title": "x", "status": "待处理", "priority": "",
                "tag": "", "type": "", "category": "",
                "modified": "2026-08-03 16:56:39", "created": "2026-06-26 17:15:07",
                "assignedTo": "", "pool": "tf_provider", "pool_project": "528766"}
        item.update(over)
        return item

    def test_tick_commits_conclusive_entries_to_disk(self):
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        path = Path(tmp.name) / "scan-snapshot.json"
        runner = self._runner(path)
        runner._dispatch = mock.Mock(return_value=(True, "new"))
        runner._scan_union = lambda: [self._union_item()]

        runner._tick()

        self.assertEqual(ScanSnapshotStore(path).load(),
                         {"83612870": "2026-08-03 16:56:39"})

    def test_tick_does_not_persist_a_failed_dispatch(self):
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        path = Path(tmp.name) / "scan-snapshot.json"
        runner = self._runner(path)
        runner._dispatch = mock.Mock(return_value=(False, "enqueue_failed"))
        runner._scan_union = lambda: [self._union_item()]

        runner._tick()

        self.assertEqual(ScanSnapshotStore(path).load(), {})

    def test_tick_keeps_existing_file_when_every_pool_query_comes_back_empty(self):
        # An empty union means every pool query failed this tick. Pruning to it would
        # wipe the file and force a cold scan on the next restart for no reason.
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        path = Path(tmp.name) / "scan-snapshot.json"
        ScanSnapshotStore(path).commit({"1": "m1"})
        runner = self._runner(path)
        runner._dispatch = mock.Mock(return_value=(True, "new"))
        runner._scan_union = lambda: []

        runner._tick()

        self.assertEqual(ScanSnapshotStore(path).load(), {"1": "m1"})


if __name__ == "__main__":
    unittest.main()
