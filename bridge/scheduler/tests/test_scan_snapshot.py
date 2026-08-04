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
import subprocess
import sys
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


class ConcurrentCommitTests(unittest.TestCase):
    """``commit`` is read-modify-write, so two writers without a lock lose a merge.

    Only the Scheduler ticks today, but all three bridge processes construct a
    ScanRunner at startup, so single-writer is currently an accident of who
    schedules the job rather than a guarantee. Real processes are used because
    ``flock`` is per-open-file-description — threads in one process would not
    exercise it.
    """

    def test_concurrent_commits_do_not_lose_entries(self):
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        path = Path(tmp.name) / "scan-snapshot.json"
        repo_root = Path(__file__).resolve().parents[3]
        writers, per_writer = 4, 25
        script = (
            "import sys\n"
            "sys.path.insert(0, %r)\n"
            "from bridge.scan_snapshot import ScanSnapshotStore\n"
            "store = ScanSnapshotStore(%r)\n"
            "tag = sys.argv[1]\n"
            "for i in range(%d):\n"
            "    store.commit({'%%s-%%d' %% (tag, i): 'm'})\n"
            % (str(repo_root), str(path), per_writer)
        )
        procs = [
            subprocess.Popen([sys.executable, "-c", script, str(writer)])
            for writer in range(writers)
        ]
        for proc in procs:
            self.assertEqual(proc.wait(timeout=180), 0, "writer subprocess failed")
        entries = ScanSnapshotStore(path).load()
        # Every key is distinct, so a lost merge shows up as a missing key.
        self.assertEqual(len(entries), writers * per_writer)


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


def _make_tick_runner(path):
    """A ScanRunner wired for a full ``_tick`` against a temp snapshot path."""
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
    r._dispatch_retry = set()
    r.dispatch_pools = set()
    r.dispatch_created_before = ""
    r._pr_merged_status_by_pool = {}
    r.pool = None
    r._envelope = mock.Mock(return_value=SimpleNamespace())
    r._reconcile_done_status_drifts_safely = mock.Mock()
    r._reconcile_source_statuses_safely = mock.Mock()
    return r


def _make_union_item(**over):
    item = {"id": "83612870", "title": "x", "status": "待处理", "priority": "",
            "tag": "", "type": "", "category": "",
            "modified": "2026-08-03 16:56:39", "created": "2026-06-26 17:15:07",
            "assignedTo": "", "pool": "tf_provider", "pool_project": "528766"}
    item.update(over)
    return item


class TickCommitTests(unittest.TestCase):
    """_tick must actually reach disk.

    The commit is wrapped in `except Exception` so snapshot IO can never fail a
    tick — which also means a missing attribute or a bad argument degrades to a
    silent no-op that every other test still passes through. This asserts the write
    end-to-end instead of trusting the call site.
    """

    _runner = staticmethod(_make_tick_runner)
    _union_item = staticmethod(_make_union_item)

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


class FailedDispatchRetryTests(unittest.TestCase):
    """A dispatch that did not land must stay a candidate on the next tick.

    ``_prev_snapshot`` is replaced wholesale *before* any decision is made, so
    without an explicit retry set a failed dispatch is simply forgotten: on the
    following tick the ticket is neither new nor modified, and it waits for an
    unrelated Aone edit to be looked at again.

    Observed in production on 2026-08-04: a momentary control-plane outage dropped
    two freshly created tickets at 10:26, seven ticks passed with no retry, and one
    of them only resurfaced because a human happened to edit it 40 minutes later.
    The durable snapshot already refused to record them — that covers a restart,
    but restarts are hours apart.
    """

    def _path(self):
        tmp = tempfile.TemporaryDirectory()
        self.addCleanup(tmp.cleanup)
        return Path(tmp.name) / "scan-snapshot.json"

    def test_next_tick_retries_a_dispatch_that_did_not_land(self):
        path = self._path()
        runner = _make_tick_runner(path)
        # The same union every tick: `modified` never changes, so the retry set is
        # the only thing that can bring this ticket back.
        runner._scan_union = lambda: [_make_union_item()]

        runner._dispatch = mock.Mock(
            return_value=(False, "control_plane_unavailable"))
        runner._tick()
        self.assertEqual(runner._dispatch.call_count, 1)
        self.assertIn("83612870", runner._dispatch_retry)
        self.assertEqual(ScanSnapshotStore(path).load(), {},
                         "失败的派发不得落盘")

        runner._dispatch = mock.Mock(return_value=(True, "new"))
        runner._tick()
        self.assertEqual(runner._dispatch.call_count, 1, "第二轮必须重试")
        self.assertNotIn("83612870", runner._dispatch_retry)
        self.assertEqual(ScanSnapshotStore(path).load(),
                         {"83612870": "2026-08-03 16:56:39"})

        runner._dispatch = mock.Mock(return_value=(True, "new"))
        runner._tick()
        self.assertEqual(runner._dispatch.call_count, 0,
                         "已经落定，不该无限重试")

    def test_conclusive_skip_clears_a_pending_dispatch_retry(self):
        runner = _make_tick_runner(self._path())
        runner._dispatch = mock.Mock(return_value=(False, "queue_full"))
        runner._tick_auto([_make_union_item()])
        self.assertIn("83612870", runner._dispatch_retry)

        # The ticket then reaches a terminal status: that decision supersedes the
        # failed enqueue, which must not be replayed.
        runner._tick_auto([_make_union_item(status="已发布")])
        self.assertNotIn("83612870", runner._dispatch_retry)

    def test_retry_set_drops_tickets_that_left_the_union(self):
        runner = _make_tick_runner(self._path())
        runner._dispatch = mock.Mock(
            return_value=(False, "control_plane_unavailable"))
        runner._scan_union = lambda: [_make_union_item()]
        runner._tick()
        self.assertIn("83612870", runner._dispatch_retry)

        runner._scan_union = lambda: [_make_union_item(id="99999999")]
        runner._tick()
        self.assertNotIn("83612870", runner._dispatch_retry,
                         "已离开 union 的工单不应留在重试集合里")


if __name__ == "__main__":
    unittest.main()
