#!/usr/bin/env python3
"""Unit: PrWatchScheduler open-PR CI-failure auto re-dispatch (post-PR 生命周期清单 #1).

Gap (workflow assess-post-pr-loop → escalation/cap-post-pr-lifecycle.md):
提交 PR 后的 open 窗口内 CI 转红时，没有任何调度器把 jarvis 重新唤起来修——PrWatchScheduler
旧逻辑对 open PR 只 `return`，`_gh_pr_state` 只查 state/mergedAt 不查 CI。会话一结束，CI 由绿转
红（rebase/base 冲突/flaky）永远靠人工发现。

本测试锁住新增的两块逻辑：
  · _gh_pr_ci 解析：CheckRun.conclusion / StatusContext.state 里哪些算失败（vs 绿/pending）。
  · _maybe_dispatch_ci_fix：失败即 force 重派 + per-head 去重 + 累计上限 escalate。

Standalone: `python3 bridge/test_prwatch_ci_fix.py`. 无 gh/a1/网络（monkeypatch + fake pool）。
"""
import importlib.util
import json
import os
import sys
import tempfile
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)

TID = "84251052"
PR = "https://github.com/aliyun/terraform-provider-alicloud/pull/9972"
PROJ = "528766"


class FakeHandler:
    def __init__(self):
        self.broadcasts = []
        self.dispatched = []

    def _broadcast(self, text):
        self.broadcasts.append(text)

    def dispatch_item(self, item_id, prompt, sid, resume, notify, target, ttype, **kw):
        self.dispatched.append({"id": item_id, "kind": kw.get("kind"), "prompt": prompt})
        return "done"


class FakePool:
    """submit 记录 (key, force, kind)；accept=True 时同步跑 work（触发 dispatch_item）。"""

    def __init__(self, accept=True):
        self.accept = accept
        self.submitted = []

    def submit(self, key, work, **kw):
        self.submitted.append({"key": key, "force": kw.get("force"), "kind": kw.get("kind")})
        if self.accept:
            work()
            return True, "dispatched"
        return False, "queue_full"


def _fake_proc(rc, out):
    class _P:
        returncode = rc
        stdout = out
        stderr = ""
    return _P()


class GhPrCiParseTest(unittest.TestCase):
    def setUp(self):
        self.sched = bot.PrWatchScheduler(FakeHandler(), FakePool())
        self._orig_run = bot.subprocess.run

    def tearDown(self):
        bot.subprocess.run = self._orig_run

    def _patch(self, payload):
        bot.subprocess.run = lambda *a, **k: _fake_proc(0, json.dumps(payload))

    def test_failure_checkrun_detected(self):
        self._patch({"headRefOid": "abc123def456", "statusCheckRollup": [
            {"name": "Compile", "conclusion": "SUCCESS"},
            {"name": "TestingCoverageRate", "conclusion": "FAILURE"}]})
        head, failing = self.sched._gh_pr_ci(PR)
        self.assertEqual(head, "abc123def456")
        self.assertEqual(failing, ["TestingCoverageRate"])

    def test_status_context_failure_detected(self):
        self._patch({"headRefOid": "d1", "statusCheckRollup": [
            {"context": "ci/jenkins", "state": "FAILURE"}]})
        _, failing = self.sched._gh_pr_ci(PR)
        self.assertEqual(failing, ["ci/jenkins"])

    def test_green_and_pending_not_failing(self):
        self._patch({"headRefOid": "e1", "statusCheckRollup": [
            {"name": "a", "conclusion": "SUCCESS"},
            {"name": "b", "status": "IN_PROGRESS"},
            {"context": "c", "state": "PENDING"}]})
        _, failing = self.sched._gh_pr_ci(PR)
        self.assertEqual(failing, [])

    def test_timed_out_and_cancelled_are_failing(self):
        self._patch({"headRefOid": "f1", "statusCheckRollup": [
            {"name": "x", "conclusion": "TIMED_OUT"},
            {"name": "y", "conclusion": "CANCELLED"}]})
        _, failing = self.sched._gh_pr_ci(PR)
        self.assertEqual(sorted(failing), ["x", "y"])

    def test_query_failure_returns_none(self):
        bot.subprocess.run = lambda *a, **k: _fake_proc(1, "")
        self.assertEqual(self.sched._gh_pr_ci(PR), (None, None))


class MaybeDispatchTest(unittest.TestCase):
    def setUp(self):
        tf = tempfile.NamedTemporaryFile(suffix=".json", delete=False, mode="w")
        tf.write(json.dumps({TID: {"pr_url": PR, "project": PROJ, "submitted_at": "x"}}))
        tf.close()
        self.tmp = tf.name
        self._orig_path = bot.PRWATCH_PATH
        self._orig_bt = bot.broadcast_target
        self._orig_by = bot.broadcast_type
        bot.PRWATCH_PATH = Path(self.tmp)
        bot.broadcast_target = lambda: "t"
        bot.broadcast_type = lambda: "ty"
        self.handler = FakeHandler()
        self.pool = FakePool()
        self.sched = bot.PrWatchScheduler(self.handler, self.pool)
        # 隔离善后 subprocess（wrap.sh / log.sh）——只验派发决策，不真连 a1。
        self.sched._comment = lambda *a, **k: None
        self.sched._escalate = lambda *a, **k: None
        self._ci = (None, None)
        self.sched._gh_pr_ci = lambda url: self._ci

    def tearDown(self):
        bot.PRWATCH_PATH = self._orig_path
        bot.broadcast_target = self._orig_bt
        bot.broadcast_type = self._orig_by
        os.unlink(self.tmp)

    def _entry(self):
        # 每轮从登记表重读 entry（复刻 _tick 的 per-tick 重读，去重字段才能生效）。
        return bot._prwatch_list()[TID]

    def test_failing_ci_dispatches_and_records(self):
        self._ci = ("sha1", ["Compile"])
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 1)
        self.assertTrue(self.pool.submitted[0]["force"], "CI-fix 应 force=True 越过 24h 去重台账")
        self.assertEqual(self.pool.submitted[0]["kind"], "pr_ci_fix")
        self.assertEqual(self.handler.dispatched[0]["id"], TID)
        self.assertIn("CI", self.handler.dispatched[0]["prompt"])
        e = self._entry()
        self.assertEqual(e["ci_fix_sha"], "sha1")
        self.assertEqual(e["ci_fix_attempts"], 1)

    def test_same_head_not_redispatched(self):
        self._ci = ("sha1", ["Compile"])
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())  # 同一失败 head
        self.assertEqual(len(self.pool.submitted), 1, "同 head 已派过，不应重复刷屏")

    def test_new_head_redispatches(self):
        self._ci = ("sha1", ["Compile"])
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self._ci = ("sha2", ["Compile"])  # 修复推了新 commit，head 变了 → 再派
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 2)
        self.assertEqual(self._entry()["ci_fix_attempts"], 2)

    def test_green_ci_no_dispatch(self):
        self._ci = ("sha1", [])
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0)

    def test_query_fail_no_dispatch(self):
        self._ci = (None, None)
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0)

    def test_exhaust_attempts_escalates_once_then_stops(self):
        for sha in ("s1", "s2", "s3", "s4"):
            self._ci = (sha, ["Compile"])
            self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 3, "自动修复上限 3 次")
        self.assertTrue(self._entry().get("ci_fix_escalated"), "超限应置 escalated")
        # escalated 后即使再来新失败 head 也不再自动派（仍看守合并）。
        self._ci = ("s5", ["Compile"])
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 3)


if __name__ == "__main__":
    unittest.main(verbosity=2)
