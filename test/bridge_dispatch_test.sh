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
          "JARVIS_REVISIT_SCHED", "JARVIS_REVISIT_HOUR", "JARVIS_REVISIT_MAX",
          "JARVIS_BACKLOG_DRAIN"):
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
        sc = b.ScanScheduler(handler=None, pool=pool)
        # 这些用例专测 tag/dedup 决策，解除灰度范围默认(否则无 pool/created 的 mock item
        # 会被 out_of_scope 提前 skip)；scope 本身由 ScopeGateTest 覆盖。
        sc.dispatch_pools = set()
        sc.dispatch_created_before = ""
        # 默认 mock 掉 _human_touched(否则 idle 分支会真发 a1 activity 查询)；
        # 默认「无人工介入」，需要人工介入的用例各自覆盖为 True。
        sc._human_touched = lambda iid: False
        sc._human_commented = lambda iid: False
        return sc

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
        # idle 且无人工介入(_human_touched=False) → 不重启实例，等每日 Revisit
        self.assertEqual(d["id"], ("skip", "idle_no_human"))
        self.assertEqual(d["fr"], ("dispatch", "new"))
        self.assertEqual(d["dd"], ("skip", "deduped"))

    def test_updated_no_tag_dispatches(self):
        # 更新单(无 jarvis 标签，pool 允许) → dispatch，force 默认 False
        sc = self._scanner()
        items = [{"id": "u1", "title": "updated fresh", "tag": []}]
        d = sc._decide(items)[0]
        self.assertEqual((d["action"], d["reason"]), ("dispatch", "new"))
        self.assertFalse(d["force"], "无 jarvis 标签的派发 force 默认 False")

    def test_idle_human_touched_dispatches_force(self):
        # jarvis-idle 且人工(辰羿)在 jarvis 上轮动作后介入 → 重新派发，且 force=True
        sc = self._scanner()
        sc._human_touched = lambda iid: True
        items = [{"id": "ih", "title": "idle human", "tag": ["jarvis-idle"]}]
        d = sc._decide(items)[0]
        self.assertEqual(d["action"], "dispatch")
        self.assertTrue(d["force"], "idle+人工介入必须 force 覆盖去重台账")

    def test_idle_human_comment_dispatches_force(self):
        # jarvis-idle 且最新评论来自人工 → 重新派发，且 force=True
        sc = self._scanner()
        sc._human_touched = lambda iid: False
        sc._human_commented = lambda iid: True
        items = [{"id": "ic", "title": "idle comment", "tag": ["jarvis-idle"]}]
        d = sc._decide(items)[0]
        self.assertEqual(d["action"], "dispatch")
        self.assertTrue(d["force"], "idle+人工评论必须 force 覆盖去重台账")

    def test_idle_jarvis_self_update_skips(self):
        # jarvis-idle 但最近动作是 jarvis 自身(_human_touched=False) → 跳过
        sc = self._scanner()
        sc._human_touched = lambda iid: False
        sc._human_commented = lambda iid: False
        items = [{"id": "is", "title": "idle self", "tag": ["jarvis-idle"]}]
        d = sc._decide(items)[0]
        self.assertEqual((d["action"], d["reason"]), ("skip", "idle_no_human"))

    def test_npe_skipped(self):
        # 人工标 jarvis-npe(路由不明) → 一律 skip npe，不自动派发
        sc = self._scanner()
        items = [{"id": "np", "title": "route unclear", "tag": ["jarvis-npe"]}]
        d = sc._decide(items)[0]
        self.assertEqual((d["action"], d["reason"]), ("skip", "npe"))

    def test_npe_overrides_idle_human_gate(self):
        # 关键防回归：idle+npe 即便人工介入(_human_touched/commented 都 True)也不重派——
        # npe 判定排在 idle 门之前，压过人工门，直到人工摘标签放行。
        sc = self._scanner()
        sc._human_touched = lambda iid: True
        sc._human_commented = lambda iid: True
        items = [{"id": "ni", "title": "idle+npe human touched",
                  "tag": ["jarvis-idle", "jarvis-npe"]}]
        d = sc._decide(items)[0]
        self.assertEqual((d["action"], d["reason"]), ("skip", "npe"),
                         "npe 必须压过 idle 人工门，有人评论也不重派")
        self.assertFalse(d["force"], "npe skip 不 force")

    def test_npe_after_claimed(self):
        # 顺序锁定：claimed 在 npe 之前——在跑实例如实报 claimed，收尾后下轮才被 npe 拦
        sc = self._scanner()
        items = [{"id": "nc", "title": "claimed+npe", "tag": ["jarvis-claimed", "jarvis-npe"]}]
        d = sc._decide(items)[0]
        self.assertEqual((d["action"], d["reason"]), ("skip", "claimed"))

    def test_terminal_status_skipped(self):
        # 终态工单(如「已发布」) → 直接 skip terminal，不派实例
        sc = self._scanner()
        items = [{"id": "t1", "title": "released", "tag": [], "status": "已发布"}]
        d = sc._decide(items)[0]
        self.assertEqual((d["action"], d["reason"]), ("skip", "terminal"))

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
    """Cold-start semantics: a (re)start only seeds the baseline snapshot and dispatches
    NOTHING — regardless of auto/supervised. A restart never re-consumes the existing
    backlog; only genuinely new / externally-updated tickets on later ticks trigger action."""

    def _scanner(self, auto):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=_ledger(tmp))
        sc = b.ScanScheduler(handler=None, pool=pool)
        sc.auto = auto
        sc._human_touched = lambda iid: False
        return sc

    def test_cold_auto_seeds_only_no_dispatch(self):
        sc = self._scanner(auto=True)
        items = [{"id": "1", "title": "a", "tag": []}, {"id": "2", "title": "b", "tag": []}]
        sc._scan = lambda: items
        calls = []
        sc._tick_auto = lambda ni, ui=None: calls.append(("auto", {i["id"] for i in ni}))
        sc._tick_supervised = lambda ni, ui=None: calls.append(("sup", None))
        self.assertTrue(sc._cold)
        sc._tick()  # cold tick → seed baseline only
        self.assertFalse(sc._cold)
        self.assertEqual(calls, [], "cold start must seed baseline only, dispatch nothing")
        self.assertEqual(set(sc._prev_snapshot.keys()), {"1", "2"})

    def test_cold_supervised_seeds_only(self):
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

    def test_new_item_after_cold_dispatches(self):
        # After the baseline seed, a genuinely new ticket on a later tick IS dispatched.
        sc = self._scanner(auto=True)
        base = [{"id": "1", "title": "a", "tag": []}, {"id": "2", "title": "b", "tag": []}]
        calls = []
        sc._tick_auto = lambda ni, ui=None: calls.append(("auto", {i["id"] for i in ni}))
        sc._scan = lambda: base
        sc._tick()  # cold → seed {1,2}, no dispatch
        self.assertEqual(calls, [])
        sc._scan = lambda: base + [{"id": "3", "title": "c", "tag": []}]
        sc._tick()  # #3 is new → auto path fires for just that one
        self.assertEqual(calls, [("auto", {"3"})])

    def test_external_update_surfaced_then_dispatched(self):
        # An existing ticket whose modified (gmtModified) time changes is classified as
        # "updated" and now flows into the dispatch decision (not new, but a candidate).
        sc = self._scanner(auto=True)
        v1 = [{"id": "7", "title": "t", "tag": [], "modified": "2026-07-05T10:00:00"}]
        v2 = [{"id": "7", "title": "t", "tag": [], "modified": "2026-07-05T11:00:00"}]
        seq = [v1, v2]
        sc._scan = lambda: seq.pop(0)
        got = []
        sc._tick_auto = lambda ni, ui=None: got.append(
            (sorted(i["id"] for i in ni), sorted((ui or {}).keys())))
        sc._tick()  # cold → seed {7}, no dispatch (_tick_auto not called)
        self.assertEqual(got, [], "cold tick seeds only, no _tick_auto call")
        sc._tick()  # #7 modified changed → updated candidate
        self.assertEqual(got, [([], ["7"])], "modified change → updated (fed to dispatch decision)")


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

    def test_revisit_query_pool_carries_tag(self):
        # _query_pool 构造的行必须带上原始 tag，供 _query 过滤 jarvis-npe。
        tmp = tempfile.mkdtemp()
        r = b.RevisitScheduler(handler=None, enabled=True, hour=9,
                               state_file=os.path.join(tmp, "revisit.last"))

        class R:
            returncode = 0
            stdout = json.dumps([{"identifier": "83", "subject": "[probe] x",
                                  "status": "问题解决中", "tag": ["jarvis-idle", "jarvis-npe"]}])
            stderr = ""
        orig = b.subprocess.run
        b.subprocess.run = lambda *a, **k: R()
        try:
            rows = r._query_pool("tf_provider", "528766")
        finally:
            b.subprocess.run = orig
        self.assertEqual(len(rows), 1)
        self.assertEqual(b._tagset(rows[0]), {"jarvis-idle", "jarvis-npe"},
                         "_query_pool 行必须带原始 tag")

    def test_revisit_query_excludes_npe(self):
        # idle+npe 与 idle 干净两条(都满足 _is_revisit_candidate) → 只有干净条进 cands。
        tmp = tempfile.mkdtemp()
        r = b.RevisitScheduler(handler=None, enabled=True, hour=9,
                               state_file=os.path.join(tmp, "revisit.last"))
        r._pool_projects = lambda: [("tf_provider", "528766")]
        r._query_pool = lambda key, project: [
            {"id": "clean", "title": "[probe] clean idle", "pool": key,
             "pool_project": project, "tag": "jarvis-idle", "description": ""},
            {"id": "npe", "title": "[probe] idle but npe", "pool": key,
             "pool_project": project, "tag": ["jarvis-idle", "jarvis-npe"], "description": ""},
        ]
        ids = [c["id"] for c in r._query()]
        self.assertIn("clean", ids, "idle 干净单进每日复查")
        self.assertNotIn("npe", ids, "idle+jarvis-npe 不投每日复查(等人工摘标签)")


class TicketPromptTest(unittest.TestCase):
    def test_prompt_carries_id_and_project(self):
        prompt = b._ticket_prompt("83902495", "F2 dispatcher", "tf_provider", "528766")
        self.assertIn("83902495", prompt)
        self.assertIn("528766", prompt)
        self.assertIn("claim", prompt)
        self.assertIn("SUSPEND", prompt)


class ScopeGateTest(unittest.TestCase):
    """灰度安全阀：_in_scope 只放行池白名单内、created 早于 cutoff 的单；_decide 对
    out-of-scope 返回 skip；pause flag 让 _tick 在 _scan 前短路。默认(未配)不限=现状。"""

    def _scanner(self):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=_ledger(tmp))
        sc = b.ScanScheduler(handler=None, pool=pool)
        sc.auto = True
        sc._human_touched = lambda iid: False
        return sc

    def test_default_is_open(self):
        # 默认(未配 env)已放开：pool 白名单空 + created 上限空 → _in_scope 对所有池/所有时间放行。
        sc = self._scanner()
        self.assertEqual(sc.dispatch_pools, set(), "default pool allowlist empty = all pools")
        self.assertEqual(sc.dispatch_created_before, "", "default created cutoff empty = no time limit")
        for pool in ("tf_provider", "tf_customer", "cloudspec", "provider_dev", "some_new_pool"):
            self.assertTrue(sc._in_scope({"pool": pool, "created": "2026-07-01 10:00"}),
                            "%s must be in scope by default" % pool)
        self.assertTrue(sc._in_scope({"pool": "tf_provider"}), "missing created ok when no cutoff")
        self.assertTrue(sc._in_scope({}), "empty item in scope when both gates open")

    def test_pool_allowlist(self):
        sc = self._scanner()
        sc.dispatch_pools = {"tf_provider"}
        sc.dispatch_created_before = ""   # 单测 pool 维度，解除 created 默认
        self.assertTrue(sc._in_scope({"pool": "tf_provider", "created": "2023-01-01 00:00"}))
        self.assertFalse(sc._in_scope({"pool": "tf_customer", "created": "2023-01-01 00:00"}))

    def test_created_before(self):
        sc = self._scanner()
        sc.dispatch_pools = set()          # 单测 created 维度，解除 pool 默认
        sc.dispatch_created_before = "2024-01-01"
        self.assertTrue(sc._in_scope({"pool": "x", "created": "2023-12-31 23:59"}))
        self.assertFalse(sc._in_scope({"pool": "x", "created": "2024-01-01 00:00"}), "cutoff day excluded")
        self.assertFalse(sc._in_scope({"pool": "x", "created": "2026-07-01 10:00"}))
        self.assertFalse(sc._in_scope({"pool": "x", "created": None}), "missing created → not in scope")
        self.assertFalse(sc._in_scope({"pool": "x"}), "no created key → not in scope")

    def test_combined_gray_release(self):
        sc = self._scanner()
        sc.dispatch_pools = {"tf_provider"}
        sc.dispatch_created_before = "2024-01-01"
        self.assertTrue(sc._in_scope({"pool": "tf_provider", "created": "2023-06-01 12:00"}))
        self.assertFalse(sc._in_scope({"pool": "tf_provider", "created": "2025-01-01 00:00"}))
        self.assertFalse(sc._in_scope({"pool": "tf_customer", "created": "2023-06-01 12:00"}))

    def test_decide_skips_out_of_scope(self):
        sc = self._scanner()
        sc.dispatch_pools = {"tf_provider"}
        sc.dispatch_created_before = "2024-01-01"
        items = [
            {"id": "1", "title": "in", "tag": [], "pool": "tf_provider", "created": "2023-01-01 00:00"},
            {"id": "2", "title": "wrongpool", "tag": [], "pool": "tf_customer", "created": "2023-01-01 00:00"},
            {"id": "3", "title": "toonew", "tag": [], "pool": "tf_provider", "created": "2025-01-01 00:00"},
        ]
        decided = {d["id"]: (d["action"], d["reason"]) for d in sc._decide(items)}
        self.assertEqual(decided["1"][0], "dispatch", "in-scope old tf_provider ticket dispatches")
        self.assertEqual(decided["2"], ("skip", "out_of_scope"), "wrong pool skipped")
        self.assertEqual(decided["3"], ("skip", "out_of_scope"), "too-new skipped")

    def test_pause_flag_short_circuits_tick(self):
        sc = self._scanner()
        scanned = []
        sc._scan = lambda: (scanned.append(True) or [])
        orig = b.REPO_ROOT
        tmp = tempfile.mkdtemp()
        try:
            b.REPO_ROOT = orig.__class__(tmp)
            pause = b.REPO_ROOT / ".my-day" / "bridge" / "pause"
            pause.parent.mkdir(parents=True, exist_ok=True)
            pause.write_text("")
            sc._tick()
            self.assertEqual(scanned, [], "pause flag must short-circuit before _scan")
            pause.unlink()
            sc._tick()
            self.assertEqual(scanned, [True], "no pause → _scan runs")
        finally:
            b.REPO_ROOT = orig


class HumanTouchedTest(unittest.TestCase):
    """_human_touched: 白名单模式 — 只有 config/contacts.json 登记人员(name/flower/id)
    的 operator 才算人工介入。Kelude/机器人等未登记身份不触发重派。
    绝不真发 a1 —— 全部 monkeypatch b.subprocess.run。"""

    def _scanner(self):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=_ledger(tmp))
        sc = b.ScanScheduler(handler=None, pool=pool)
        # Mock whitelist: simulate contacts.json with known human operators
        sc._human_operators = {"辰羿", "陈汉璋", "320687", "马裘", "逄磊", "404709"}
        return sc

    def _fake_run(self, rc, stdout, stderr=""):
        class R:
            returncode = rc
        R.stdout = stdout
        R.stderr = stderr
        return lambda *a, **k: R()

    def test_human_operator_is_touched(self):
        sc = self._scanner()
        orig = b.subprocess.run
        b.subprocess.run = self._fake_run(0, json.dumps([{"operator": "辰羿", "eventTime": 2}]))
        try:
            self.assertTrue(sc._human_touched("1"), "白名单内 operator → 人工介入")
        finally:
            b.subprocess.run = orig

    def test_jarvis_operator_not_touched(self):
        sc = self._scanner()
        orig = b.subprocess.run
        for op in ("open-jarvis", "WORKER_1782379562571"):
            sc._human_cache = {}
            b.subprocess.run = self._fake_run(0, json.dumps([{"operator": op, "eventTime": 2}]))
            try:
                self.assertFalse(sc._human_touched("1"), "%s 不在白名单 → 非人工" % op)
            finally:
                b.subprocess.run = orig

    def test_kelude_not_touched(self):
        """Kelude 等系统机器人不在白名单 → 不算人工介入(修复 #83899246 误派)。"""
        sc = self._scanner()
        orig = b.subprocess.run
        b.subprocess.run = self._fake_run(0, json.dumps([{"operator": "Kelude", "eventTime": 2}]))
        try:
            self.assertFalse(sc._human_touched("1"), "Kelude 不在白名单 → 非人工介入")
        finally:
            b.subprocess.run = orig

    def test_empty_whitelist_returns_false(self):
        """白名单为空集 → 任何 operator 都返回 False(保守不误派)。"""
        sc = self._scanner()
        sc._human_operators = set()  # simulate contacts.json load failure
        orig = b.subprocess.run
        b.subprocess.run = self._fake_run(0, json.dumps([{"operator": "辰羿", "eventTime": 2}]))
        try:
            self.assertFalse(sc._human_touched("1"), "白名单为空 → 一律不算人工介入")
        finally:
            b.subprocess.run = orig

    def test_failure_returns_false_conservatively(self):
        sc = self._scanner()
        orig = b.subprocess.run

        def boom(*a, **k):
            raise RuntimeError("network down")
        b.subprocess.run = boom
        try:
            self.assertFalse(sc._human_touched("9"), "查询异常一律保守返回 False")
        finally:
            b.subprocess.run = orig

    def test_nonzero_rc_returns_false(self):
        sc = self._scanner()
        orig = b.subprocess.run
        b.subprocess.run = self._fake_run(1, "", "boom")
        try:
            self.assertFalse(sc._human_touched("9"), "rc!=0 → 保守 False")
        finally:
            b.subprocess.run = orig

    def test_result_cached_within_tick(self):
        sc = self._scanner()
        calls = [0]
        orig = b.subprocess.run

        def counting(*a, **k):
            calls[0] += 1
            class R:
                returncode = 0
                stdout = json.dumps([{"operator": "辰羿"}])
                stderr = ""
            return R()
        b.subprocess.run = counting
        try:
            sc._human_touched("5")
            sc._human_touched("5")
            self.assertEqual(calls[0], 1, "同一 tick 内同 id 只查一次")
        finally:
            b.subprocess.run = orig


class HumanCommentedTest(unittest.TestCase):
    """_human_commented: idle 更新时补看 Aone 最新评论。
    最新评论来自人工 → 触发重派；最新评论来自 open-jarvis/系统账号 → 保守跳过。
    绝不真发 a1 —— 全部 monkeypatch b.subprocess.run。"""

    def _scanner(self):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(dedup_ttl=86400, ledger_path=_ledger(tmp))
        return b.ScanScheduler(handler=None, pool=pool)

    def _fake_run(self, rc, stdout, stderr=""):
        class R:
            returncode = rc
        R.stdout = stdout
        R.stderr = stderr
        return lambda *a, **k: R()

    def test_latest_human_comment_is_touched(self):
        sc = self._scanner()
        comments = [
            {"id": 1, "author": "open-jarvis", "content": "jarvis-claim x"},
            {"id": 2, "author": "辰羿", "content": "@open-jarvis 继续处理"},
        ]
        orig = b.subprocess.run
        b.subprocess.run = self._fake_run(0, json.dumps(comments))
        try:
            self.assertTrue(sc._human_commented("1"), "最新评论来自人工 → 需要重派")
        finally:
            b.subprocess.run = orig

    def test_latest_jarvis_comment_not_touched(self):
        sc = self._scanner()
        comments = [
            {"id": 1, "author": "辰羿", "content": "@open-jarvis 继续处理"},
            {"id": 2, "author": "open-jarvis", "content": "jarvis-claim x"},
        ]
        orig = b.subprocess.run
        b.subprocess.run = self._fake_run(0, json.dumps(comments))
        try:
            self.assertFalse(sc._human_commented("1"), "最新评论来自 Jarvis → 不重派")
        finally:
            b.subprocess.run = orig

    def test_human_comment_after_idle_counts_even_if_latest_is_jarvis(self):
        sc = self._scanner()
        activities = [
            {"property": "标签", "oldValue": "jarvis-claimed",
             "newValue": "jarvis-idle", "eventTime": "2026-07-06 10:00"},
        ]
        comments = [
            {"id": 1, "author": "辰羿", "createdAt": "2026-07-06 09:59:59",
             "content": "@open-jarvis 旧评论"},
            {"id": 2, "author": "辰羿", "createdAt": "2026-07-06 10:01:00",
             "content": "@open-jarvis 继续处理"},
            {"id": 3, "author": "open-jarvis", "createdAt": "2026-07-06 10:02:00",
             "content": "jarvis-claim x"},
        ]
        orig = b.subprocess.run

        def fake_run(args, **kwargs):
            class R:
                returncode = 0
                stderr = ""
            R.stdout = json.dumps(activities if "activity" in args else comments)
            return R()

        b.subprocess.run = fake_run
        try:
            self.assertTrue(sc._human_commented("1"),
                            "只要 idle 后存在人工评论，即使最新评论是 Jarvis 也需要重派")
        finally:
            b.subprocess.run = orig

    def test_human_comment_before_idle_does_not_count(self):
        sc = self._scanner()
        activities = [
            {"property": "标签", "oldValue": "jarvis-claimed",
             "newValue": "jarvis-idle", "eventTime": "2026-07-06 10:00"},
        ]
        comments = [
            {"id": 1, "author": "辰羿", "createdAt": "2026-07-06 09:59:59",
             "content": "@open-jarvis 旧评论"},
        ]
        orig = b.subprocess.run

        def fake_run(args, **kwargs):
            class R:
                returncode = 0
                stderr = ""
            R.stdout = json.dumps(activities if "activity" in args else comments)
            return R()

        b.subprocess.run = fake_run
        try:
            self.assertFalse(sc._human_commented("1"),
                             "idle 前的历史人工评论不能触发重派")
        finally:
            b.subprocess.run = orig

    def test_latest_system_comment_not_touched(self):
        sc = self._scanner()
        comments = [{"id": 1, "author": "Kelude", "content": "sync"}]
        orig = b.subprocess.run
        b.subprocess.run = self._fake_run(0, json.dumps(comments))
        try:
            self.assertFalse(sc._human_commented("1"), "系统账号评论 → 不重派")
        finally:
            b.subprocess.run = orig

    def test_comment_query_failure_returns_false(self):
        sc = self._scanner()
        orig = b.subprocess.run
        b.subprocess.run = self._fake_run(1, "", "boom")
        try:
            self.assertFalse(sc._human_commented("9"), "rc!=0 → 保守 False")
        finally:
            b.subprocess.run = orig

    def test_comment_result_cached_within_tick(self):
        sc = self._scanner()
        comment_calls = [0]
        activity_calls = [0]
        orig = b.subprocess.run

        def counting(args, **kwargs):
            class R:
                returncode = 0
                stderr = ""
            if "activity" in args:
                activity_calls[0] += 1
                R.stdout = "[]"
            else:
                comment_calls[0] += 1
                R.stdout = json.dumps([{"id": 1, "author": "辰羿", "content": "继续"}])
            return R()
        b.subprocess.run = counting
        try:
            sc._human_commented("5")
            sc._human_commented("5")
            self.assertEqual(comment_calls[0], 1, "同一 tick 内同 id 评论只查一次")
            self.assertEqual(activity_calls[0], 1, "同一 tick 内同 id activity 只查一次")
        finally:
            b.subprocess.run = orig


class NoDingtalkDegradedTest(unittest.TestCase):
    """无钉钉降级模式(JARVIS_NO_DINGTALK=1, 缺凭证): 调度器照起、TataPool 跳过、
    卡片/播报降级为 [BROADCAST] 日志行、唤醒走 headless 池 —— 全程绝不触钉钉。"""

    class _BoomSM:
        """任何被当作 streaming 模块访问的属性都炸 —— 证明降级路径根本没碰钉钉 SDK。"""
        def __getattr__(self, name):
            raise AssertionError("DingTalk streaming touched in no-dingtalk mode: .%s" % name)

    def _handler(self):
        # 无 DINGTALK_APP_KEY/SECRET 也能构造; 只建调度器, 不 start()(不起线程/子进程)。
        for k in ("DINGTALK_APP_KEY", "DINGTALK_APP_SECRET", "DINGTALK_TEMPLATE_ID"):
            os.environ.pop(k, None)
        h = b.JarvisHandler(no_dingtalk=True)
        h.sm = self._BoomSM()   # 任何真钉钉调用即触发 AssertionError
        return h

    def test_schedulers_built_tatapool_skipped(self):
        h = self._handler()
        self.assertTrue(h.no_dingtalk)
        self.assertIsNone(h.pool, "TataPool must be skipped in no-dingtalk mode")
        for name in ("dispatch_pool", "scanner", "reconciler", "board",
                     "prober", "reviser", "watcher"):
            self.assertIsNotNone(getattr(h, name, None), "%s must be constructed" % name)
        self.assertTrue(h.scanner.auto, "auto-dispatch stays on in degraded mode")

    def test_broadcast_and_card_land_in_log_no_dingtalk_call(self):
        h = self._handler()
        with self.assertLogs("jarvis-bot", level="INFO") as cm:
            h._broadcast("已自动派发 #123 测试单\n- 第二行")   # 经 _quick_card 降级
            h._quick_card("someuser", "定向卡片也降级", "user")
        joined = "\n".join(cm.output)
        self.assertIn("[BROADCAST]", joined)
        self.assertIn("已自动派发 #123", joined)
        self.assertIn("定向卡片也降级", joined)
        # _BoomSM 未抛异常即证明没走真钉钉(get_access_token 等未被调用)。

    def test_scheduler_broadcast_path_no_dingtalk(self):
        # 调度器 dispatch 汇总播报回调 = handler._broadcast → 降级落 [BROADCAST] 日志, 不触钉钉。
        # mock submit 免真起 claude(新语义: new/updated 单都过 _decide 派发, 派发汇总即播报)。
        h = self._handler()
        h.dispatch_pool.submit = lambda *a, **k: (True, "dispatched")  # 不真跑 work
        with self.assertLogs("jarvis-bot", level="INFO") as cm:
            h.scanner._tick_auto([{"id": "77", "title": "新单派发", "tag": [],
                                   "pool": "tf_provider", "pool_project": "528766"}], {})
        self.assertIn("[BROADCAST]", "\n".join(cm.output))

    def test_wake_resumes_via_headless_pool_not_card(self):
        h = self._handler()
        captured = {}

        def fake_submit(item_id, work, *, notify=None, force=False, kind="ticket"):
            captured.update(id=str(item_id), force=force, kind=kind)
            return True, "dispatched"   # 不真跑 work(不起 claude)

        h.dispatch_pool.submit = fake_submit
        # 若降级唤醒误走卡片路(_submit_card→_stream_round→get_access_token), _BoomSM 会炸。
        task = {"target": "grp", "target_type": "group", "session_id": "sess-1"}
        with self.assertLogs("jarvis-bot", level="INFO") as cm:
            h._wake("83929676", task, [{"creator": "someone", "content": "继续"}])
        self.assertEqual(captured.get("id"), "83929676")
        self.assertTrue(captured.get("force"), "wake must force past the 24h dedup ledger")
        self.assertEqual(captured.get("kind"), "wake",
                         "degraded wake must use the headless pool path, not the card path")
        self.assertIn("[BROADCAST]", "\n".join(cm.output))  # 「收到回复,唤醒中」通知降级为日志


class TataResidentModeTest(unittest.TestCase):
    """Tata 默认不保留常驻 claude 子进程；需要保温时必须显式打开。"""

    def setUp(self):
        for k in ("JARVIS_TATA_RESIDENT", "JARVIS_TATA_PREWARM"):
            os.environ.pop(k, None)

    def tearDown(self):
        for k in ("JARVIS_TATA_RESIDENT", "JARVIS_TATA_PREWARM"):
            os.environ.pop(k, None)

    def test_default_tata_runner_is_one_shot_without_pool(self):
        h = b.JarvisHandler(no_dingtalk=False)
        self.assertIsNone(h.pool, "default live bridge must not retain a TataPool")

        calls = []
        orig = b.run_tata_stream

        def fake_run_tata_stream(text, session_id, resume):
            calls.append((text, session_id, resume))
            yield "one-shot reply"

        b.run_tata_stream = fake_run_tata_stream
        try:
            self.assertEqual(list(h._tata_runner("hi", "staff-1", False)), ["one-shot reply"])
            self.assertEqual(calls, [("hi", "staff-1", False)])
        finally:
            b.run_tata_stream = orig
            h.dispatch_pool.shutdown(wait=False, cancel_futures=True)

    def test_resident_mode_requires_explicit_env(self):
        os.environ["JARVIS_TATA_RESIDENT"] = "1"
        h = b.JarvisHandler(no_dingtalk=False)
        try:
            self.assertIsInstance(h.pool, b.TataPool)
        finally:
            h.dispatch_pool.shutdown(wait=False, cancel_futures=True)
            h.pool.shutdown()


class BoardEnabledTest(unittest.TestCase):
    """board.enabled 隐雷: live main() 引用 handler.board.enabled — 该属性必须存在且为 bool,
    否则配上真钉钉凭证走 live 路径会 AttributeError。"""

    def test_enabled_attribute_present_and_bool(self):
        board = b.BoardScheduler(None)
        self.assertTrue(hasattr(board, "enabled"), "BoardScheduler must expose .enabled")
        self.assertIsInstance(board.enabled, bool)
        # 复刻 main()/降级路径的判定表达式: 不得抛 AttributeError; base_url 默认非空 → started。
        self.assertEqual("started" if board.enabled else "disabled", "started")
        # 降级路径同款判定也不得抛
        self.assertEqual("on" if board.enabled else "off", "on")


class BoardProbeMergeTest(unittest.TestCase):
    """BoardScheduler 并入 board.sh probe 健康度段(E 线 MR-5)。probe 成功→并入 "probe" 键;
    probe 缺失/失败→绝不影响主 items 同步(容错依赖 MR-5)。全程桩 board.sh, 不依赖 MR-5 是否合并。"""

    def test_probe_merged_into_payload_when_available(self):
        board = b.BoardScheduler(None)
        board._probe_section = lambda: {"pass": 3, "fail": 1, "warn": 0}
        p = board._build_payload('[{"id": "1", "state": "pool"}]')
        self.assertEqual(p["items"], [{"id": "1", "state": "pool"}])
        self.assertIn("ts", p)
        self.assertEqual(p["probe"], {"pass": 3, "fail": 1, "warn": 0})

    def test_items_sync_unaffected_when_probe_absent(self):
        board = b.BoardScheduler(None)
        board._probe_section = lambda: None
        p = board._build_payload('[{"id": "1"}]')
        self.assertNotIn("probe", p, "no probe key when probe section unavailable")
        self.assertEqual(p["items"], [{"id": "1"}])  # 主同步不受影响
        self.assertIn("ts", p)

    def test_probe_section_two_states(self):
        board = b.BoardScheduler(None)
        board._board_supports_probe = lambda: True   # 桩: 假装 board.sh 已有 probe 子命令(MR-5 后)

        class _R:
            def __init__(self, rc, out):
                self.returncode = rc; self.stdout = out; self.stderr = ""

        orig = b.subprocess.run
        try:
            b.subprocess.run = lambda *a, **k: _R(0, '{"pass": 2, "fail": 0}')  # 成功: 对象
            self.assertEqual(board._probe_section(), {"pass": 2, "fail": 0})
            b.subprocess.run = lambda *a, **k: _R(1, "")                        # 失败: rc!=0
            self.assertIsNone(board._probe_section())
            b.subprocess.run = lambda *a, **k: _R(0, '[{"id": "1"}]')           # 非对象(items 回退)
            self.assertIsNone(board._probe_section(), "array output must not be treated as probe")
        finally:
            b.subprocess.run = orig

    def test_probe_skipped_without_invoking_board_when_subcommand_absent(self):
        # 主线 board.sh 无 probe 子命令 → 存在性探测 False → 不调 subprocess, 直接 None(不污染 payload)。
        board = b.BoardScheduler(None)
        board._board_supports_probe = lambda: False
        called = []
        orig = b.subprocess.run
        try:
            b.subprocess.run = lambda *a, **k: called.append(1)
            self.assertIsNone(board._probe_section())
            self.assertEqual(called, [], "must not run board.sh when probe subcommand absent")
        finally:
            b.subprocess.run = orig


class CompletionBroadcastTest(unittest.TestCase):
    """_completion_broadcast: headless 完成播报按最终 tag 区分状态 + 附可点击 Aone 链接。
    绝不真发 a1 —— 全部 monkeypatch b.subprocess.run；非数字 id / 查询失败 → 回退纯文本。"""

    @staticmethod
    def _workitem_json(tag_display, title="支持 X 参数", project="1086837", priority="高"):
        return json.dumps({
            "id": "83884678",
            "title": title,
            "url": "https://project.aone.alibaba-inc.com/v2/project/%s/req/83884678" % project,
            "fields": [
                {"identifier": "status", "value": "155741", "displayValue": "问题解决中"},
                {"identifier": "priority", "value": "95", "displayValue": priority},
                {"identifier": "space", "value": project, "displayValue": "Terraform - 客户问题"},
                {"identifier": "tag", "value": "568724", "displayValue": tag_display},
            ],
        })

    def _fake_run(self, rc, stdout, stderr=""):
        class R:
            returncode = rc
        R.stdout = stdout
        R.stderr = stderr
        return lambda *a, **k: R()

    def _broadcast(self, item_id, run_fn):
        orig = b.subprocess.run
        b.subprocess.run = run_fn
        try:
            return b.JarvisHandler._completion_broadcast(None, item_id)
        finally:
            b.subprocess.run = orig

    def test_done_tag_has_prefix_link_title_priority(self):
        msg = self._broadcast("83884678", self._fake_run(0, self._workitem_json("jarvis-done")))
        self.assertIn("✅ 工单处理完成", msg)
        self.assertIn(
            "[#83884678](https://project.aone.alibaba-inc.com/v2/project/1086837/req/83884678)", msg)
        self.assertIn("支持 X 参数", msg)
        self.assertIn("[高]", msg)

    def test_idle_tag_pending_human(self):
        msg = self._broadcast("83884678", self._fake_run(0, self._workitem_json("jarvis-idle")))
        self.assertIn("⏸️ 工单阶段完成·待人工接手", msg)
        self.assertIn(
            "[#83884678](https://project.aone.alibaba-inc.com/v2/project/1086837/req/83884678)", msg)

    def test_claimed_tag_unfinished(self):
        msg = self._broadcast("83884678", self._fake_run(0, self._workitem_json("jarvis-claimed")))
        self.assertIn("⚠️", msg)
        self.assertIn("未收尾", msg)

    def test_non_numeric_id_fallback_without_a1(self):
        called = [0]

        def spy(*a, **k):
            called[0] += 1
            raise AssertionError("a1 must not be called for non-numeric id")
        msg = self._broadcast("probe-2026-07-06", spy)
        self.assertEqual(called[0], 0, "非数字 id 不得调用 a1")
        self.assertIn("✅ 任务 #probe-2026-07-06 处理完成", msg)

    def test_query_nonzero_rc_fallback(self):
        msg = self._broadcast("83884678", self._fake_run(1, "", "boom"))
        self.assertEqual(msg, "✅ 工单 #83884678 处理完成（headless）")

    def test_query_exception_fallback(self):
        def boom(*a, **k):
            raise RuntimeError("network down")
        msg = self._broadcast("83884678", boom)
        self.assertEqual(msg, "✅ 工单 #83884678 处理完成（headless）")


class FreeSlotsTest(unittest.TestCase):
    """DispatchPool.free_slots() = max_workers - 在飞任务数(含排队)，clamp 到 >=0。
    free_slots>0 蕴含队列必空(在飞数<并发上限时不可能排队)——backlog drain 的核心不变量:
    只填真正空闲的运行槽，绝不把积压单排到未来新单前面。"""

    def test_free_slots_tracks_active(self):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(max_workers=2, queue_max=5, ledger_path=_ledger(tmp))
        self.assertEqual(pool.free_slots(), 2, "empty pool = max_workers free")
        pool._active["a"] = {"started": time.time()}
        self.assertEqual(pool.free_slots(), 1, "one in flight = max_workers-1")
        pool._active["b"] = {"started": time.time()}
        self.assertEqual(pool.free_slots(), 0, "at concurrency cap = 0 free")
        pool._active["c"] = {"started": time.time()}  # queued beyond max_workers
        self.assertEqual(pool.free_slots(), 0, "overfull (queued) clamps to 0, never negative")


class PriorityRankTest(unittest.TestCase):
    """_priority_rank: 紧急<高<中<低，未知/空值排最后(9)——只影响空闲槽<积压单数时先派谁。"""

    def test_known_ranks_strictly_ordered(self):
        self.assertLess(b._priority_rank("紧急"), b._priority_rank("高"))
        self.assertLess(b._priority_rank("高"), b._priority_rank("中"))
        self.assertLess(b._priority_rank("中"), b._priority_rank("低"))

    def test_unknown_and_blank_rank_last(self):
        self.assertEqual(b._priority_rank(""), 9)
        self.assertEqual(b._priority_rank(None), 9)
        self.assertGreater(b._priority_rank("P0"), b._priority_rank("低"))
        self.assertGreater(b._priority_rank("whatever"), b._priority_rank("低"))

    def test_whitespace_tolerated(self):
        self.assertEqual(b._priority_rank(" 高 "), b._priority_rank("高"))


class BacklogDrainTest(unittest.TestCase):
    """空闲机会式消化积压(JARVIS_BACKLOG_DRAIN): 仅 auto+开关开、无新单/更新单、池有空闲运行槽
    的 tick 才涓流派发从未被 jarvis 碰过的积压单；新单永远优先；冷启动那轮不消化。"""

    def _scanner(self, max_workers=2):
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(max_workers=max_workers, dedup_ttl=86400, ledger_path=_ledger(tmp))
        sc = b.ScanScheduler(handler=None, pool=pool)
        sc.auto = True
        sc.backlog_drain = True
        sc.dispatch_pools = set()          # 解除灰度范围默认(mock item 无 pool/created)
        sc.dispatch_created_before = ""
        sc._human_touched = lambda iid: False
        return sc, pool

    def _capture_dispatch(self, sc):
        """替换 _dispatch 记录 (id, force)，不真起 pool/claude；返回记录列表。"""
        got = []
        sc._dispatch = lambda it, force=False: (
            got.append((str(it.get("id")), force)) or (True, "ok"))
        return got

    def test_idle_tick_drains_untouched_backlog(self):
        sc, _pool = self._scanner()
        got = self._capture_dispatch(sc)
        item = {"id": "bk", "title": "backlog one", "tag": []}
        sc._scan = lambda: [item]
        sc._tick()   # cold → seed baseline only
        self.assertEqual(got, [], "cold tick drains nothing")
        sc._tick()   # no new/updated + free slot → drain
        self.assertEqual(got, [("bk", False)], "idle tick drains the untouched backlog item")

    def test_new_item_tick_skips_drain(self):
        # 同一轮既有新单又有积压候选 → 只派新单，积压本轮不派(新单优先)。
        sc, _pool = self._scanner()
        auto_calls, drain_calls = [], []
        sc._tick_auto = lambda ni, ui=None: auto_calls.append({str(i["id"]) for i in ni})
        sc._tick_backlog = lambda snap: drain_calls.append(set(snap))
        backlog = {"id": "bk", "title": "old backlog", "tag": []}
        sc._scan = lambda: [backlog]
        sc._tick()   # cold seed {bk}
        sc._scan = lambda: [backlog, {"id": "nw", "title": "new", "tag": []}]
        sc._tick()   # nw is new → auto path, drain must NOT run
        self.assertEqual(auto_calls, [{"nw"}], "only the new item flows to auto dispatch")
        self.assertEqual(drain_calls, [], "a tick with new items must not run backlog drain")

    def test_full_pool_no_drain(self):
        sc, pool = self._scanner(max_workers=2)
        got = self._capture_dispatch(sc)
        pool._active = {"busy1": {"started": time.time()}, "busy2": {"started": time.time()}}
        sc._scan = lambda: [{"id": "bk", "title": "backlog", "tag": []}]
        sc._tick()   # cold seed
        sc._tick()   # idle tick but free_slots==0 → no drain
        self.assertEqual(got, [], "no free running slot → backlog not drained")

    def test_queued_pool_no_drain(self):
        # 3 在飞 with max_workers=2 → 有排队；free_slots=max(0,2-3)=0 → 不消化。
        sc, pool = self._scanner(max_workers=2)
        got = self._capture_dispatch(sc)
        pool._active = {"a": {}, "b": {}, "c": {}}
        sc._scan = lambda: [{"id": "bk", "title": "backlog", "tag": []}]
        sc._tick(); sc._tick()
        self.assertEqual(pool.free_slots(), 0, "queue non-empty ⇒ free_slots 0")
        self.assertEqual(got, [], "queued jobs mean no free slot → never queue backlog ahead of new")

    def test_cold_start_no_drain(self):
        sc, _pool = self._scanner()
        drain_calls, auto_calls = [], []
        sc._tick_backlog = lambda snap: drain_calls.append(1)
        sc._tick_auto = lambda ni, ui=None: auto_calls.append(1)
        sc._scan = lambda: [{"id": "bk", "title": "backlog", "tag": []}]
        self.assertTrue(sc._cold)
        sc._tick()   # cold tick
        self.assertFalse(sc._cold)
        self.assertEqual(drain_calls, [], "cold start must not drain backlog (restart no-storm)")
        self.assertEqual(auto_calls, [], "cold start dispatches nothing")
        self.assertEqual(set(sc._prev_snapshot), {"bk"}, "cold start seeds baseline")

    def test_drain_respects_dedup_ledger(self):
        # 真 submit + FakeHandler: 消化过一条后同快照再 tick，被去重台账/active 拦下不重复派。
        tmp = tempfile.mkdtemp()
        pool = b.DispatchPool(max_workers=2, dedup_ttl=86400, ledger_path=_ledger(tmp))

        class FakeHandler:
            def _broadcast(self, text):
                pass

            def dispatch_item(self, *a, **k):
                return "ok"

        sc = b.ScanScheduler(handler=FakeHandler(), pool=pool)
        sc.auto = True
        sc.backlog_drain = True
        sc.dispatch_pools = set()
        sc.dispatch_created_before = ""
        sc._human_touched = lambda iid: False

        submits = []
        real_submit = pool.submit

        def spy(iid, work, **k):
            r = real_submit(iid, work, **k)
            submits.append((str(iid), r[0]))
            return r
        pool.submit = spy

        item = {"id": "bk", "title": "backlog", "tag": [],
                "pool": "tf_provider", "pool_project": "528766"}
        sc._scan = lambda: [item]
        sc._tick()   # cold seed
        sc._tick()   # drain #1 → bk dispatched, ledger records it
        sc._tick()   # drain #2 → bk deduped/active, NOT re-dispatched
        pool.shutdown(wait=True)

        bk = [ok for iid, ok in submits if iid == "bk"]
        self.assertEqual(bk, [True, False],
                         "first drain dispatches bk; second rejected by dedup ledger/active")
        self.assertIn("bk", pool._ledger, "dispatched backlog id recorded in dedup ledger")

    def test_is_backlog_filters_tags_terminal_scope(self):
        sc, _pool = self._scanner()
        self.assertTrue(sc._is_backlog({"id": "1", "tag": []}), "clean untouched ticket = backlog")
        self.assertFalse(sc._is_backlog({"id": "2", "tag": ["jarvis-claimed"]}), "claimed ≠ backlog")
        self.assertFalse(sc._is_backlog({"id": "3", "tag": ["jarvis-idle"]}), "idle ≠ backlog")
        self.assertFalse(sc._is_backlog({"id": "4", "tag": ["jarvis-done"]}), "done ≠ backlog")
        self.assertFalse(sc._is_backlog({"id": "5", "tag": [], "status": "已发布"}), "terminal ≠ backlog")
        sc.dispatch_pools = {"tf_provider"}   # 收窄灰度范围
        self.assertTrue(sc._is_backlog({"id": "6", "tag": [], "pool": "tf_provider"}))
        self.assertFalse(sc._is_backlog({"id": "7", "tag": [], "pool": "tf_customer"}), "out of scope ≠ backlog")

    def test_is_backlog_excludes_npe_but_keeps_probe(self):
        sc, _pool = self._scanner()
        self.assertFalse(sc._is_backlog({"id": "n", "tag": ["jarvis-npe"]}), "npe ≠ backlog (路由不明)")
        # 摘掉 npe → 恢复为可消化积压
        self.assertTrue(sc._is_backlog({"id": "n", "tag": []}), "npe 摘掉后恢复为积压")
        # 不做 jarvis-* 前缀匹配：jarvis-probe 仍是合法派发对象
        self.assertTrue(sc._is_backlog({"id": "p", "tag": ["jarvis-probe"]}),
                        "jarvis-probe 不在排除集，仍算积压(勿改成前缀匹配)")

    def test_drain_skips_npe_item(self):
        sc, _pool = self._scanner()
        got = self._capture_dispatch(sc)
        snapshot = [
            {"id": "clean", "title": "clean", "tag": []},
            {"id": "npe", "title": "route unclear", "tag": ["jarvis-npe"]},
        ]
        sc._scan = lambda: snapshot
        sc._tick()   # cold seed
        sc._tick()   # drain
        self.assertEqual([iid for iid, _ in got], ["clean"], "空闲消化只派干净单，跳过 jarvis-npe")

    def test_drain_only_picks_backlog_from_mixed_snapshot(self):
        sc, _pool = self._scanner()
        got = self._capture_dispatch(sc)
        snapshot = [
            {"id": "clean", "title": "clean", "tag": []},
            {"id": "clm", "title": "claimed", "tag": ["jarvis-claimed"]},
            {"id": "idl", "title": "idle", "tag": ["jarvis-idle"]},
            {"id": "dn", "title": "done", "tag": ["jarvis-done"]},
            {"id": "term", "title": "released", "tag": [], "status": "已发布"},
        ]
        sc._scan = lambda: snapshot
        sc._tick()   # cold seed
        sc._tick()   # drain
        self.assertEqual([iid for iid, _ in got], ["clean"],
                         "only the untagged non-terminal in-scope item is drained")

    def test_switch_off_restores_pure_new_behavior(self):
        sc, _pool = self._scanner()
        sc.backlog_drain = False
        drain_calls = []
        sc._tick_backlog = lambda snap: drain_calls.append(1)
        got = self._capture_dispatch(sc)
        sc._scan = lambda: [{"id": "bk", "title": "backlog", "tag": []}]
        sc._tick()   # cold seed
        sc._tick()   # idle tick, but drain disabled
        self.assertEqual(drain_calls, [], "JARVIS_BACKLOG_DRAIN off → _tick_backlog never called")
        self.assertEqual(got, [], "switch off → an idle tick dispatches nothing (old behavior)")

    def test_drain_sorts_priority_high_first(self):
        # free=1, 两条不同优先级 → 先派高优先级(低优先级列在前证明 sort 生效)。
        sc, _pool = self._scanner(max_workers=1)
        got = self._capture_dispatch(sc)
        items = [
            {"id": "lo", "title": "low", "tag": [], "priority": "低", "created": "2026-07-01 09:00"},
            {"id": "hi", "title": "high", "tag": [], "priority": "高", "created": "2026-07-05 09:00"},
        ]
        sc._scan = lambda: items
        sc._tick(); sc._tick()
        self.assertEqual([iid for iid, _ in got], ["hi"], "free_slots=1 → higher priority drained first")

    def test_drain_created_tiebreak_older_first(self):
        # 同优先级 → 先派 created 更早的(旧单先)。
        sc, _pool = self._scanner(max_workers=1)
        got = self._capture_dispatch(sc)
        items = [
            {"id": "newer", "title": "n", "tag": [], "priority": "中", "created": "2026-07-05 09:00"},
            {"id": "older", "title": "o", "tag": [], "priority": "中", "created": "2026-07-01 09:00"},
        ]
        sc._scan = lambda: items
        sc._tick(); sc._tick()
        self.assertEqual([iid for iid, _ in got], ["older"], "same priority → earlier gmtCreate first")


class TataSessionIdTest(unittest.TestCase):
    """Regression: Tata 会话 id 必须是合法 UUID(而非 staffId)。

    一次性冷起模式(默认 JARVIS_TATA_RESIDENT unset → pool None)下，claude CLI 对
    --session-id/--resume 强校验；曾直传 staffId → "Invalid session ID. Must be a
    valid UUID." 门面整条 Tata 链报错。见 _tata_session()。"""

    def _handler(self):
        os.environ.pop("JARVIS_TATA_RESIDENT", None)  # 默认一次性冷起 → pool None
        h = b.JarvisHandler(no_dingtalk=True)
        self.assertIsNone(h.pool, "test 依赖一次性冷起模式(pool None)")
        return h

    def test_session_id_is_uuid_not_staffid(self):
        import uuid as _uuid
        h = self._handler()
        staff = "320687"
        sid, resume = h._tata_session(staff)
        _uuid.UUID(sid)  # 非法 UUID 会抛 → 失败
        self.assertNotEqual(sid, staff, "session id 不能是 staffId")
        self.assertFalse(resume, "首轮应 --session-id 建会话, 非 --resume")

    def test_session_stable_and_resumes_after_first_turn(self):
        h = self._handler()
        staff = "320687"
        sid1, r1 = h._tata_session(staff)
        sid2, r2 = h._tata_session(staff)
        self.assertEqual(sid1, sid2, "同 staff 会话 id 需稳定(续聊)")
        self.assertFalse(r1, "首轮 session")
        self.assertTrue(r2, "后续轮 resume")

    def test_distinct_staff_distinct_session(self):
        h = self._handler()
        sid_a, _ = h._tata_session("111")
        sid_b, r_b = h._tata_session("222")
        self.assertNotEqual(sid_a, sid_b, "不同 staff 会话隔离")
        self.assertFalse(r_b, "新 staff 首轮 session")


# ==========================================================================
# Headless dispatch resilience (feat/bridge-dispatch-resilience)
# run_claude_buffered + _classify_result + bounded resume retry + failure→Aone.
# All mocked: no live claude / a1 / wrap.sh / DingTalk.
# ==========================================================================

RESUME_HINT = "从中断处继续"


class _RetrySelf:
    """Minimal fake ``self`` for exercising JarvisHandler.dispatch_item in isolation.
    Stubs the collaborators dispatch_item touches so no live claude/a1/DingTalk runs."""

    def __init__(self, suspend=None, completion="✅ done"):
        self._suspend = suspend
        self._completion = completion
        self.failed_calls = []

    def _maybe_suspend(self, final, sid, target, target_type):
        return self._suspend

    def _completion_broadcast(self, item_id):
        return self._completion

    def _workitem_line(self, aone_id):
        return "#%s" % aone_id

    def _dispatch_failed(self, item_id, res, notify, project):
        self.failed_calls.append((item_id, res, project))


def _scripted_runner(results):
    """(fake_run_claude_buffered, calls): each entry records prompt/sid/resume so the
    retry semantics (same sid, resume/续跑 vs fresh) can be asserted per attempt."""
    calls = []
    it = iter(results)

    def fake(text, session_id, resume, timeout=None, on_spawn=None):
        calls.append({"text": text, "sid": session_id, "resume": resume})
        return next(it)

    return fake, calls


class ClassifyResultTest(unittest.TestCase):
    """_classify_result: pure parse of `claude --output-format json`. Reads
    is_error/subtype/result from the LAST result object; rc!=0 always overrides
    is_error True (the stream-path bug it fixes); no result + rc!=0 → stderr tail."""

    def test_success(self):
        out = json.dumps({"type": "result", "is_error": False,
                          "subtype": "success", "result": "hi there"})
        self.assertEqual(b._classify_result(out, "", 0),
                         b.ClaudeResult("hi there", False, "success"))

    def test_is_error_true(self):
        out = json.dumps({"type": "result", "is_error": True,
                          "subtype": "error_during_execution", "result": "partial"})
        r = b._classify_result(out, "", 0)
        self.assertTrue(r.is_error)
        self.assertEqual(r.subtype, "error_during_execution")

    def test_rc_override_forces_error(self):
        out = json.dumps({"type": "result", "is_error": False,
                          "subtype": "success", "result": "ok"})
        r = b._classify_result(out, "", 1)
        self.assertTrue(r.is_error, "rc!=0 must force is_error True even if json says False")
        self.assertEqual(r.text, "ok")

    def test_leading_banner_tolerated_takes_last(self):
        lines = "\n".join([
            "Some non-JSON banner line the gateway printed",
            json.dumps({"type": "system", "subtype": "init"}),
            json.dumps({"type": "result", "is_error": False,
                        "subtype": "success", "result": "first"}),
            json.dumps({"type": "result", "is_error": False,
                        "subtype": "success", "result": "last"}),
        ])
        r = b._classify_result(lines, "", 0)
        self.assertEqual(r.text, "last", "must take the LAST result object")
        self.assertFalse(r.is_error)

    def test_no_result_rc_nonzero_uses_stderr(self):
        r = b._classify_result("garbage not json\n\n",
                               "boom line1\nfatal: gateway 400 模型提供方错误", 1)
        self.assertEqual(
            r, b.ClaudeResult("fatal: gateway 400 模型提供方错误", True, "no_result"))


class BufferedRunnerCmdTest(unittest.TestCase):
    """run_claude_buffered: argv shape (json, NOT stream-json/partial/verbose),
    session/resume flag selection, and hard-timeout killpg. Popen fully faked."""

    def setUp(self):
        self._orig_popen = b.subprocess.Popen
        self._orig_cmd = b.jarvis_cmd
        b.jarvis_cmd = lambda sid=None: ["claude", "--settings", "/fake/idea.json"]

    def tearDown(self):
        b.subprocess.Popen = self._orig_popen
        b.jarvis_cmd = self._orig_cmd

    def _fake_popen(self, captured, out):
        class FakeP:
            def __init__(s, argv, **kw):
                captured["argv"] = argv
                s.pid = 4242
                s.returncode = 0

            def communicate(s, timeout=None):
                return (out, "")
        return FakeP

    def test_argv_json_no_stream_flags_session(self):
        captured = {}
        out = json.dumps({"type": "result", "is_error": False,
                          "subtype": "success", "result": "ok"})
        b.subprocess.Popen = self._fake_popen(captured, out)
        res = b.run_claude_buffered("hello", "sid-1", False, timeout=10)
        argv = captured["argv"]
        self.assertEqual(argv[argv.index("--output-format") + 1], "json")
        self.assertNotIn("--include-partial-messages", argv)
        self.assertNotIn("--verbose", argv)
        self.assertIn("--session-id", argv)
        self.assertEqual(argv[argv.index("--session-id") + 1], "sid-1")
        self.assertNotIn("--resume", argv)
        self.assertEqual(res, b.ClaudeResult("ok", False, "success"))

    def test_resume_flag(self):
        captured = {}
        out = json.dumps({"type": "result", "is_error": False,
                          "subtype": "success", "result": "ok"})
        b.subprocess.Popen = self._fake_popen(captured, out)
        b.run_claude_buffered("hi", "sid-2", True, timeout=10)
        argv = captured["argv"]
        self.assertIn("--resume", argv)
        self.assertEqual(argv[argv.index("--resume") + 1], "sid-2")
        self.assertNotIn("--session-id", argv)

    def test_timeout_kills_group(self):
        killed = {}

        class FakeP:
            def __init__(s, argv, **kw):
                s.pid = 9999

            def communicate(s, timeout=None):
                raise b.subprocess.TimeoutExpired(cmd="claude", timeout=timeout or 1)
        b.subprocess.Popen = FakeP
        orig_killpg, orig_getpgid = b.os.killpg, b.os.getpgid
        b.os.getpgid = lambda pid: pid
        b.os.killpg = lambda pgid, sig: killed.setdefault("pgid", pgid)
        try:
            res = b.run_claude_buffered("hi", "sid-3", False, timeout=1)
        finally:
            b.os.killpg, b.os.getpgid = orig_killpg, orig_getpgid
        self.assertTrue(res.is_error)
        self.assertEqual(res.subtype, "timeout")
        self.assertEqual(killed.get("pgid"), 9999, "must attempt killpg on the group")


class DispatchRetryTest(unittest.TestCase):
    """dispatch_item bounded resume retry loop: same sid, retry transient, fast-fail
    terminal, SUSPEND not treated as failure. run_claude_buffered scripted, sleep noop."""

    def setUp(self):
        self._orig_rcb = b.run_claude_buffered
        self._orig_sleep = b.time.sleep
        b.time.sleep = lambda *a, **k: None
        os.environ["JARVIS_DISPATCH_RETRY_MAX"] = "2"

    def tearDown(self):
        b.run_claude_buffered = self._orig_rcb
        b.time.sleep = self._orig_sleep
        os.environ.pop("JARVIS_DISPATCH_RETRY_MAX", None)

    def _dispatch(self, fs, results, sid="sid-fixed", project=None):
        runner, calls = _scripted_runner(results)
        b.run_claude_buffered = runner
        notifies = []
        outcome = b.JarvisHandler.dispatch_item(
            fs, "123", "orig prompt", sid, False,
            notifies.append, "grp", "group", project=project)
        return outcome, calls, notifies

    def test_retry_then_success(self):
        ERR = b.ClaudeResult("partial", True, "error")
        OK = b.ClaudeResult("done text", False, "success")
        fs = _RetrySelf()
        outcome, calls, _ = self._dispatch(fs, [ERR, ERR, OK])
        self.assertEqual(outcome, "done")
        self.assertEqual(len(calls), 3, "1 initial + 2 retries")
        self.assertEqual({c["sid"] for c in calls}, {"sid-fixed"}, "same sid every attempt")
        for c in calls[1:]:
            self.assertTrue(c["resume"], "retry attempts resume=True")
            self.assertIn(RESUME_HINT, c["text"], "retry attempts use 续跑 prompt")
        self.assertEqual(fs.failed_calls, [], "clean finish → no _dispatch_failed")

    def test_all_error_bounded(self):
        ERR = b.ClaudeResult("partial", True, "error")
        fs = _RetrySelf()
        outcome, calls, _ = self._dispatch(fs, [ERR] * 6)
        self.assertEqual(outcome, "error")
        self.assertEqual(len(calls), 3, "MAX=2 → exactly 3 calls, no runaway")
        self.assertEqual(len(fs.failed_calls), 1)

    def test_first_clean_no_retry(self):
        OK = b.ClaudeResult("done", False, "success")
        fs = _RetrySelf()
        outcome, calls, _ = self._dispatch(fs, [OK])
        self.assertEqual(outcome, "done")
        self.assertEqual(len(calls), 1)

    def test_suspend_not_retried(self):
        SUS = b.ClaudeResult("blah [[SUSPEND:{...}]]", False, "success")
        fs = _RetrySelf(suspend={"aone_id": "123", "wait_for": "someone"})
        outcome, calls, notifies = self._dispatch(fs, [SUS])
        self.assertEqual(outcome, "suspended")
        self.assertEqual(len(calls), 1, "clean SUSPEND → no retry")
        self.assertEqual(fs.failed_calls, [], "SUSPEND is not a failure")

    def test_terminal_timeout_not_retried(self):
        TO = b.ClaudeResult("", True, "timeout")
        fs = _RetrySelf()
        outcome, calls, _ = self._dispatch(fs, [TO])
        self.assertEqual(outcome, "error")
        self.assertEqual(len(calls), 1, "terminal timeout → no retry")
        self.assertEqual(len(fs.failed_calls), 1)


class ResumeFallbackTest(unittest.TestCase):
    """Retry continuation choice: last attempt produced output → resume=True + 续跑;
    no output → fall back to fresh (resume=False + original prompt)."""

    def setUp(self):
        self._orig_rcb = b.run_claude_buffered
        self._orig_sleep = b.time.sleep
        b.time.sleep = lambda *a, **k: None
        os.environ["JARVIS_DISPATCH_RETRY_MAX"] = "2"

    def tearDown(self):
        b.run_claude_buffered = self._orig_rcb
        b.time.sleep = self._orig_sleep
        os.environ.pop("JARVIS_DISPATCH_RETRY_MAX", None)

    def _calls(self, results):
        runner, calls = _scripted_runner(results)
        b.run_claude_buffered = runner
        b.JarvisHandler.dispatch_item(
            _RetrySelf(), "123", "ORIGINAL", "sid-x", False,
            (lambda t: None), "grp", "group")
        return calls

    def test_empty_text_falls_back_to_fresh(self):
        ERR = b.ClaudeResult("", True, "error")
        OK = b.ClaudeResult("done", False, "success")
        calls = self._calls([ERR, OK])
        self.assertFalse(calls[1]["resume"], "no output → fresh restart (resume=False)")
        self.assertEqual(calls[1]["text"], "ORIGINAL", "fresh restart reuses original prompt")

    def test_nonempty_text_resumes(self):
        ERR = b.ClaudeResult("some progress", True, "error")
        OK = b.ClaudeResult("done", False, "success")
        calls = self._calls([ERR, OK])
        self.assertTrue(calls[1]["resume"], "had output → resume=True")
        self.assertIn(RESUME_HINT, calls[1]["text"])


class _DFSelf:
    _post_death_cause = staticmethod(b.JarvisHandler._post_death_cause)


class DispatchFailedTest(unittest.TestCase):
    """_dispatch_failed: post death cause via wrap.sh sync (numeric id only) + release
    claim (project set only) + notify. wrap.sh & _release_claim fully stubbed."""

    def setUp(self):
        self._orig_run = b.subprocess.run
        self._orig_release = b._release_claim
        self.run_calls = []
        self.release_calls = []

        class R:
            returncode = 0
            stdout = ""
            stderr = ""

        def fake_run(*a, **k):
            self.run_calls.append((a, k))
            return R()

        b.subprocess.run = fake_run
        b._release_claim = lambda iid, project: self.release_calls.append((iid, project))

    def tearDown(self):
        b.subprocess.run = self._orig_run
        b._release_claim = self._orig_release

    def _fail(self, item_id, project, text="internal tail output"):
        res = b.ClaudeResult(text, True, "error")
        notifies = []
        b.JarvisHandler._dispatch_failed(_DFSelf(), item_id, res, notifies.append, project)
        return notifies

    def test_numeric_id_posts_and_releases(self):
        notifies = self._fail("84080203", "528766")
        self.assertEqual(len(self.run_calls), 1, "one wrap.sh sync")
        argv = self.run_calls[0][0][0]
        self.assertIn(str(b.REPO_ROOT / "bootstrap" / "wrap.sh"), argv)
        self.assertIn("sync", argv)
        self.assertIn("84080203", argv)
        self.assertIn("--summary-stdin", argv)
        stdin = self.run_calls[0][1].get("input", "")
        self.assertIn("subtype: error", stdin, "death cause carries the subtype")
        self.assertEqual(self.release_calls, [("84080203", "528766")])
        self.assertTrue(any("失败" in n for n in notifies), "notify carries a failure line")

    def test_project_none_posts_but_no_release(self):
        notifies = self._fail("84080203", None)
        self.assertEqual(len(self.run_calls), 1, "death cause still posted")
        self.assertEqual(self.release_calls, [], "no project → no release")
        self.assertTrue(any("失败" in n for n in notifies))

    def test_non_numeric_id_no_subprocess(self):
        self._fail("probe-2026-07-08", "528766")
        self.assertEqual(len(self.run_calls), 0, "pseudo-id has no workitem → no wrap.sh")
        self.assertEqual(self.release_calls, [], "non-numeric id → no release")


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
