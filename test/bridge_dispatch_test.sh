#!/usr/bin/env bash
# test/bridge_dispatch_test.sh — hermetic unit tests for bridge/jarvis_dingtalk_bot.py
# F2 池调度 dispatcher: DispatchPool 并发/排队/软去重, ScanScheduler auto 决策,
# Probe/Revisit 每日门判定。全部 mock(不真起 claude、不连钉钉、不打 Aone)。
#
# Wraps a python3 unittest suite. The bridge module guards its dingtalk_stream
# import so it imports without the SDK installed; every test constructs the
# dispatcher primitives directly (no live JarvisHandler / WebSocket / a1).
#
# Run: bash test/bridge_dispatch_test.sh
# Exit 0 = all pass, non-zero = failure.

set -uo pipefail

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$test_dir/.." && pwd)"

if ! command -v python3 >/dev/null 2>&1; then
  echo "SKIP bridge_dispatch_test: python3 not available"
  exit 0
fi

BRIDGE_DIR="$repo_root/bridge" python3 - "$repo_root" <<'PY'
import os
import sys
import json
import time
import tempfile
import threading
import datetime as dt
import unittest

sys.path.insert(0, os.environ["BRIDGE_DIR"])
# Neutralize env so defaults are exercised deterministically.
for k in ("JARVIS_AUTO_DISPATCH", "JARVIS_DISPATCH_MAX", "JARVIS_DISPATCH_QUEUE_MAX",
          "JARVIS_DISPATCH_DEDUP_TTL", "JARVIS_PROBE_SCHED", "JARVIS_PROBE_HOUR",
          "JARVIS_REVISIT_SCHED", "JARVIS_REVISIT_HOUR", "JARVIS_REVISIT_MAX"):
    os.environ.pop(k, None)

import jarvis_dingtalk_bot as b


def _ledger(tmp):
    return os.path.join(tmp, "dispatched.json")


class DispatchPoolConcurrencyTest(unittest.TestCase):
    def test_cap_and_queue(self):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(max_workers=2, queue_max=3, ledger_path=_ledger(tmp))
        gate = threading.Event()
        lock = threading.Lock()
        running = []
        peak = [0]
        done = []

        def work():
            with lock:
                running.append(1)
                peak[0] = max(peak[0], len(running))
            gate.wait(5)
            with lock:
                running.pop()
            done.append(1)
            return "done"

        # cap = max_workers + queue_max = 5; first 5 accepted.
        results = [pool.submit("id%d" % i, work) for i in range(5)]
        self.assertTrue(all(ok for ok, _ in results), results)
        # 6th exceeds cap → rejected as queue_full (workers still blocked on gate).
        ok, reason = pool.submit("id5", work)
        self.assertFalse(ok)
        self.assertEqual(reason, "queue_full")
        time.sleep(0.3)  # give the executor time to start its 2 threads
        with lock:
            self.assertLessEqual(peak[0], 2, "concurrency exceeded max_workers")
        gate.set()
        pool.shutdown(wait=True)
        self.assertEqual(len(done), 5)
        self.assertEqual(peak[0], 2, "both worker threads should have run concurrently")


class DispatchPoolDedupTest(unittest.TestCase):
    def test_active_dedup(self):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(max_workers=2, ledger_path=_ledger(tmp))
        gate = threading.Event()

        def work():
            gate.wait(5)
            return "done"

        ok, _ = pool.submit("x", work)
        self.assertTrue(ok)
        ok, reason = pool.submit("x", work)  # still running
        self.assertFalse(ok)
        self.assertEqual(reason, "active")
        gate.set()
        pool.shutdown(wait=True)

    def test_ledger_ttl_window(self):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=_ledger(tmp))
        pool._ledger["y"] = time.time()
        ok, reason = pool.status("y")
        self.assertFalse(ok)
        self.assertEqual(reason, "deduped")
        pool._ledger["y"] = time.time() - 90000  # older than 24h
        ok, _ = pool.status("y")
        self.assertTrue(ok)

    def test_force_overrides_ledger_not_active(self):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=_ledger(tmp))
        pool._ledger["z"] = time.time()
        self.assertFalse(pool.status("z")[0])            # deduped without force
        self.assertTrue(pool.status("z", force=True)[0])  # force ignores ledger

    def test_restart_recovery(self):
        tmp = tempfile.mkdtemp()
        p = _ledger(tmp)
        with open(p, "w") as f:
            json.dump({"recent": time.time(), "old": time.time() - 90000}, f)
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=p)
        self.assertFalse(pool.status("recent")[0], "recent entry should dedup after restart")
        self.assertTrue(pool.status("old")[0], "expired entry should not dedup")

    def test_submit_persists_ledger(self):
        tmp = tempfile.mkdtemp()
        p = _ledger(tmp)
        pool = b.DispatchPool(max_workers=1, ledger_path=p)
        gate = threading.Event()
        ok, _ = pool.submit("persisted", lambda: gate.wait(5))
        self.assertTrue(ok)
        self.assertTrue(os.path.exists(p), "ledger file should be written on submit")
        with open(p) as f:
            self.assertIn("persisted", json.load(f))
        gate.set()
        pool.shutdown(wait=True)


class ScanDecideTest(unittest.TestCase):
    def _scanner(self, ledger_recent=()):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=_ledger(tmp))
        for iid in ledger_recent:
            pool._ledger[str(iid)] = time.time()
        return b.ScanScheduler(handler=None, pool=pool)

    def test_tag_and_dedup_decisions(self):
        sc = self._scanner(ledger_recent=["dd"])
        items = [
            {"id": "cl", "title": "claimed one", "tag": ["jarvis-claimed"]},
            {"id": "dn", "title": "done one", "tag": ["jarvis-done"]},
            {"id": "id", "title": "idle one", "tag": ["jarvis-idle"]},
            {"id": "fr", "title": "fresh one", "tag": []},
            {"id": "dd", "title": "deduped one", "tag": []},
        ]
        d = {x["id"]: (x["action"], x["reason"]) for x in sc._decide(items)}
        self.assertEqual(d["cl"], ("skip", "claimed"))
        self.assertEqual(d["dn"], ("skip", "done"))
        self.assertEqual(d["id"], ("skip", "idle"))   # parked → RevisitScheduler owns it
        self.assertEqual(d["fr"], ("dispatch", "new"))
        self.assertEqual(d["dd"], ("skip", "deduped"))

    def test_tag_string_and_missing(self):
        sc = self._scanner()
        items = [
            {"id": "a", "title": "csv tags", "tag": "foo,jarvis-claimed,bar"},
            {"id": "b", "title": "no tag field"},
            {"id": "c", "title": "null tag", "tag": None},
        ]
        d = {x["id"]: (x["action"], x["reason"]) for x in sc._decide(items)}
        self.assertEqual(d["a"], ("skip", "claimed"))
        self.assertEqual(d["b"], ("dispatch", "new"))
        self.assertEqual(d["c"], ("dispatch", "new"))

    def test_auto_default_on(self):
        sc = self._scanner()
        self.assertTrue(sc.auto, "JARVIS_AUTO_DISPATCH defaults to on")


class ColdStartTest(unittest.TestCase):
    """Fused cold-start semantics (rebase onto fc2a7ea): supervised seeds baseline + returns
    with no notification; auto DOES dispatch the qualifying backlog on the first tick
    (diff-independent — gated by tag-skip + the dedup ledger)."""

    def _scanner(self, auto):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=_ledger(tmp))
        sc = b.ScanScheduler(handler=None, pool=pool)
        sc.auto = auto
        return sc

    def test_cold_auto_dispatches_backlog(self):
        sc = self._scanner(auto=True)
        items = [{"id": "1", "title": "a", "tag": []}, {"id": "2", "title": "b", "tag": []}]
        sc._scan = lambda: items
        calls = []
        sc._tick_auto = lambda ni, ui=None: calls.append(("auto", {i["id"] for i in ni}))
        sc._tick_supervised = lambda ni, ui=None: calls.append(("sup", None))
        self.assertTrue(sc._cold)
        sc._tick()  # cold tick
        self.assertFalse(sc._cold)
        self.assertEqual(calls, [("auto", {"1", "2"})], "auto cold start must dispatch backlog")
        self.assertEqual(set(sc._prev_snapshot.keys()), {"1", "2"})

    def test_cold_supervised_seeds_only_then_diffs(self):
        sc = self._scanner(auto=False)
        items = [{"id": "1", "title": "a", "tag": []}, {"id": "2", "title": "b", "tag": []}]
        sc._scan = lambda: items
        calls = []
        sc._tick_auto = lambda ni, ui=None: calls.append("auto")
        sc._tick_supervised = lambda ni, ui=None: calls.append(("sup", {i["id"] for i in ni}))
        sc._tick()  # cold tick → seed only, no notification
        self.assertEqual(calls, [], "supervised cold start must not notify/dispatch")
        self.assertEqual(set(sc._prev_snapshot.keys()), {"1", "2"})
        self.assertFalse(sc._cold)
        # second tick with one genuinely new item → supervised path fires for just that one
        sc._scan = lambda: items + [{"id": "3", "title": "c", "tag": []}]
        sc._tick()
        self.assertEqual(calls, [("sup", {"3"})])

    def test_external_update_surfaced_not_dispatched(self):
        # Merge with master 6de3d06: an existing ticket whose modified (gmtModified) time
        # changes is classified as "updated" (sensing/notify), NOT "new" — never dispatched.
        sc = self._scanner(auto=True)
        v1 = [{"id": "7", "title": "t", "tag": [], "modified": "2026-07-05T10:00:00"}]
        v2 = [{"id": "7", "title": "t", "tag": [], "modified": "2026-07-05T11:00:00"}]
        seq = [v1, v2]
        sc._scan = lambda: seq.pop(0)
        got = []
        sc._tick_auto = lambda ni, ui=None: got.append(
            (sorted(i["id"] for i in ni), sorted((ui or {}).keys())))
        sc._tick()  # cold auto → #7 dispatched as new
        sc._tick()  # #7 modified changed → updated, not new
        self.assertEqual(got[0], (["7"], []), "first tick: #7 is new")
        self.assertEqual(got[1], ([], ["7"]), "modified change → updated (surfaced, not dispatched)")


class DailySchedulerTest(unittest.TestCase):
    def test_probe_due_gating(self):
        tmp = tempfile.mkdtemp()
        st = os.path.join(tmp, "probe.last")
        p = b.ProbeScheduler(handler=None, enabled=True, hour=10, state_file=st)
        self.assertFalse(p._due(now=dt.datetime(2026, 7, 3, 9, 0)))   # before hour
        self.assertTrue(p._due(now=dt.datetime(2026, 7, 3, 11, 0)))   # after hour, not run
        p._mark_run(when=dt.datetime(2026, 7, 3, 11, 0))
        self.assertFalse(p._due(now=dt.datetime(2026, 7, 3, 23, 0)))  # already ran today
        self.assertTrue(p._due(now=dt.datetime(2026, 7, 4, 10, 30)))  # next day

    def test_disabled_never_due(self):
        tmp = tempfile.mkdtemp()
        st = os.path.join(tmp, "revisit.last")
        r = b.RevisitScheduler(handler=None, enabled=False, hour=9, state_file=st)
        self.assertFalse(r._due(now=dt.datetime(2026, 7, 3, 12, 0)))

    def test_revisit_candidate_matching(self):
        tmp = tempfile.mkdtemp()
        r = b.RevisitScheduler(handler=None, enabled=True, hour=9,
                               state_file=os.path.join(tmp, "revisit.last"))
        self.assertTrue(r._is_revisit_candidate({"title": "[probe] dns bug"}))
        self.assertTrue(r._is_revisit_candidate({"title": "x", "description": "待续条件: PR 合并"}))
        self.assertTrue(r._is_revisit_candidate({"title": "y", "description": "等待 maintainer 合并"}))
        self.assertFalse(r._is_revisit_candidate({"title": "ordinary ticket"}))


class TicketPromptTest(unittest.TestCase):
    def test_prompt_carries_id_and_project(self):
        prompt = b._ticket_prompt("83902495", "F2 dispatcher", "tf_provider", "528766")
        self.assertIn("83902495", prompt)
        self.assertIn("528766", prompt)
        self.assertIn("claim", prompt)
        self.assertIn("SUSPEND", prompt)


def _run():
    loader = unittest.TestLoader()
    suite = loader.loadTestsFromModule(sys.modules[__name__])
    result = unittest.TextTestRunner(verbosity=2).run(suite)
    return 0 if result.wasSuccessful() else 1


sys.exit(_run())
PY
rc=$?

echo ""
if [ "$rc" -eq 0 ]; then
  echo "bridge_dispatch_test: PASS"
else
  echo "bridge_dispatch_test: FAIL (rc=$rc)"
fi
exit "$rc"
