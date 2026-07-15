#!/usr/bin/env python3
"""Unit: PrWatchScheduler open-PR 全生命周期看守（post-PR 生命周期清单 #1/#2/#3）.

Gap (workflow assess-post-pr-loop → escalation/cap-post-pr-lifecycle.md):
提交 PR 后单次 headless 会话撑不住多日合并窗口——CI 转红没人修、reviewer 在 GitHub 上评论没人
回应、轮询频率跟不上。本测试锁住 PrWatchScheduler 升级为「全生命周期看守器」后的三块逻辑：
  · #1 _gh_pr_ci 解析 (failing/pending) + _maybe_dispatch_ci_fix 派发/去重/上限/active 信号；
  · #2 _gh_pr_comments 解析(排除自己/bot) + _maybe_dispatch_comment_reply baseline+新评论派发；
  · #3 active 信号驱动双档轮询（_maybe_dispatch_ci_fix 返回 True/False）。

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
        head, failing, pending = self.sched._gh_pr_ci(PR)
        self.assertEqual(head, "abc123def456")
        self.assertEqual(failing, ["TestingCoverageRate"])
        self.assertFalse(pending)

    def test_status_context_failure_detected(self):
        self._patch({"headRefOid": "d1", "statusCheckRollup": [
            {"context": "ci/jenkins", "state": "FAILURE"}]})
        _, failing, _ = self.sched._gh_pr_ci(PR)
        self.assertEqual(failing, ["ci/jenkins"])

    def test_pending_flagged_not_failing(self):
        self._patch({"headRefOid": "e1", "statusCheckRollup": [
            {"name": "a", "conclusion": "SUCCESS"},
            {"name": "b", "status": "IN_PROGRESS"},
            {"context": "c", "state": "PENDING"}]})
        _, failing, pending = self.sched._gh_pr_ci(PR)
        self.assertEqual(failing, [])
        self.assertTrue(pending, "in-progress/pending 应置 pending=True → 快档轮询")

    def test_all_green_not_pending(self):
        self._patch({"headRefOid": "g1", "statusCheckRollup": [
            {"name": "a", "conclusion": "SUCCESS"},
            {"name": "b", "conclusion": "NEUTRAL"},
            {"name": "c", "conclusion": "SKIPPED"}]})
        _, failing, pending = self.sched._gh_pr_ci(PR)
        self.assertEqual((failing, pending), ([], False))

    def test_timed_out_and_cancelled_are_failing(self):
        self._patch({"headRefOid": "f1", "statusCheckRollup": [
            {"name": "x", "conclusion": "TIMED_OUT"},
            {"name": "y", "conclusion": "CANCELLED"}]})
        _, failing, _ = self.sched._gh_pr_ci(PR)
        self.assertEqual(sorted(failing), ["x", "y"])

    def test_query_failure_returns_none(self):
        bot.subprocess.run = lambda *a, **k: _fake_proc(1, "")
        self.assertEqual(self.sched._gh_pr_ci(PR), (None, None, False))


class GhPrCommentsParseTest(unittest.TestCase):
    def setUp(self):
        self.sched = bot.PrWatchScheduler(FakeHandler(), FakePool())
        self._orig_run = bot.subprocess.run

    def tearDown(self):
        bot.subprocess.run = self._orig_run

    def _patch(self, comments):
        bot.subprocess.run = lambda *a, **k: _fake_proc(0, json.dumps({"comments": comments}))

    def test_latest_reviewer_comment_picked(self):
        self._patch([
            {"author": {"login": "api-tool-agent"}, "body": "mine", "url": "u0"},
            {"author": {"login": "reviewer1"}, "body": "please fix X", "url": "u1"},
            {"author": {"login": "reviewer2"}, "body": "and Y", "url": "u2"}])
        key, author, snippet = self.sched._gh_pr_comments(PR)
        self.assertEqual((key, author), ("u2", "reviewer2"))
        self.assertIn("Y", snippet)

    def test_self_and_bot_excluded(self):
        self._patch([
            {"author": {"login": "reviewer1"}, "body": "real", "url": "u1"},
            {"author": {"login": "github-actions[bot]"}, "body": "ci", "url": "u2"},
            {"author": {"login": "api-tool-agent"}, "body": "self reply", "url": "u3"}])
        key, author, _ = self.sched._gh_pr_comments(PR)
        self.assertEqual((key, author), ("u1", "reviewer1"),
                         "最新的非我方/非 bot 评论才算——bot/self 都要跳过")

    def test_no_reviewer_comment_returns_none(self):
        self._patch([{"author": {"login": "api-tool-agent"}, "body": "x", "url": "u1"}])
        self.assertEqual(self.sched._gh_pr_comments(PR), (None, None, None))

    def test_query_failure_returns_none(self):
        bot.subprocess.run = lambda *a, **k: _fake_proc(1, "")
        self.assertEqual(self.sched._gh_pr_comments(PR), (None, None, None))


class _DispatchBase(unittest.TestCase):
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
        self.sched._comment = lambda *a, **k: None
        self.sched._escalate = lambda *a, **k: None

    def tearDown(self):
        bot.PRWATCH_PATH = self._orig_path
        bot.broadcast_target = self._orig_bt
        bot.broadcast_type = self._orig_by
        os.unlink(self.tmp)

    def _entry(self):
        return bot._prwatch_list()[TID]


class MaybeDispatchCiFixTest(_DispatchBase):
    def setUp(self):
        super().setUp()
        self._ci = (None, None, False)
        self.sched._gh_pr_ci = lambda url: self._ci

    def test_failing_dispatches_records_and_is_active(self):
        self._ci = ("sha1", ["Compile"], False)
        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertTrue(active, "CI 失败 → active(快档)")
        self.assertEqual(len(self.pool.submitted), 1)
        self.assertTrue(self.pool.submitted[0]["force"], "CI-fix 应 force=True 越过 24h 去重")
        self.assertEqual(self.pool.submitted[0]["kind"], "pr_ci_fix")
        e = self._entry()
        self.assertEqual((e["ci_fix_sha"], e["ci_fix_attempts"]), ("sha1", 1))

    def test_same_head_not_redispatched_but_still_active(self):
        self._ci = ("sha1", ["Compile"], False)
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 1, "同 head 不重复派")
        self.assertTrue(active, "仍失败中 → 保持快档")

    def test_new_head_redispatches(self):
        self._ci = ("sha1", ["Compile"], False)
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self._ci = ("sha2", ["Compile"], False)
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 2)
        self.assertEqual(self._entry()["ci_fix_attempts"], 2)

    def test_green_no_dispatch_not_active(self):
        self._ci = ("sha1", [], False)
        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0)
        self.assertFalse(active, "全绿 → 慢档等合并")

    def test_pending_no_dispatch_but_active(self):
        self._ci = ("sha1", [], True)
        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0, "pending 不派修复")
        self.assertTrue(active, "pending → 快档等结果")

    def test_query_fail_no_dispatch_not_active(self):
        self._ci = (None, None, False)
        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0)
        self.assertFalse(active)

    def test_exhaust_attempts_escalates_once_then_stops(self):
        for sha in ("s1", "s2", "s3", "s4"):
            self._ci = (sha, ["Compile"], False)
            self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 3, "自动修复上限 3 次")
        self.assertTrue(self._entry().get("ci_fix_escalated"))
        self._ci = ("s5", ["Compile"], False)
        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 3)
        self.assertFalse(active, "escalated 后转人工 → 不再快档")


class MaybeDispatchCommentReplyTest(_DispatchBase):
    def setUp(self):
        super().setUp()
        self._c = (None, None, None)
        self.sched._gh_pr_comments = lambda url: self._c

    def test_first_seen_seeds_baseline_no_dispatch(self):
        self._c = ("u1", "reviewer1", "please fix")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0, "首次观察只 baseline，不回应既有评论")
        self.assertEqual(self._entry().get("last_seen_comment"), "u1")

    def test_new_comment_dispatches(self):
        self._c = ("u1", "reviewer1", "c1")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # baseline
        self._c = ("u2", "reviewer2", "please change Y")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # new
        self.assertEqual(len(self.pool.submitted), 1)
        self.assertEqual(self.pool.submitted[0]["kind"], "pr_comment_reply")
        self.assertTrue(self.pool.submitted[0]["force"])
        self.assertEqual(self._entry().get("last_seen_comment"), "u2")

    def test_same_comment_not_redispatched(self):
        self._c = ("u1", "reviewer1", "c1")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # baseline u1
        self._c = ("u2", "reviewer2", "c2")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # dispatch u2
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # u2 again → skip
        self.assertEqual(len(self.pool.submitted), 1, "同一条最新评论不重复回应")

    def test_no_reviewer_comment_no_dispatch(self):
        self._c = (None, None, None)
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0)
        self.assertIsNone(self._entry().get("last_seen_comment"))


class AutoRegisterTest(_DispatchBase):
    """#6 兜底发现：漏登记的 open PR 自动补登记（分支解析工单号 + aone-get 校验）。"""

    def setUp(self):
        super().setUp()
        bot._prwatch_write({})  # 清空 base 的 TID entry，从零测漏登发现
        self._prs = []
        self._proj = "528766"
        self.sched._gh_open_prs = lambda: self._prs
        self.sched._ticket_project = lambda tid: self._proj

    @staticmethod
    def _pr(n, branch):
        return {"number": n, "headRefName": branch,
                "url": "https://github.com/aliyun/terraform-provider-alicloud/pull/%d" % n}

    def test_branch_encoded_ticket_autoregisters(self):
        self._prs = [self._pr(9972, "feat/84291978-tair")]
        self.sched._maybe_autoregister_open_prs()
        reg = bot._prwatch_list()
        self.assertIn("84291978", reg)
        self.assertEqual(reg["84291978"]["project"], "528766")
        self.assertEqual(reg["84291978"]["pr_url"], self._prs[0]["url"])

    def test_already_watched_skipped(self):
        pr = self._pr(9972, "feat/84291978-tair")
        bot._prwatch_add("84291978", pr["url"], "528766")
        self._prs = [pr]
        called = []
        self.sched._ticket_project = lambda tid: called.append(tid) or "528766"
        self.sched._maybe_autoregister_open_prs()
        self.assertEqual(called, [], "已登记 PR 不应再查 project/重登")

    def test_unparseable_branch_not_registered(self):
        self._prs = [self._pr(100, "feat/gpdb-api-key")]  # 无工单号
        self.sched._maybe_autoregister_open_prs()
        self.assertEqual(bot._prwatch_list(), {}, "分支无工单号 → 不瞎登")

    def test_project_lookup_fail_not_registered(self):
        self._proj = None
        self._prs = [self._pr(101, "feat/12345678-x")]  # 有号但校验失败
        self.sched._maybe_autoregister_open_prs()
        self.assertEqual(bot._prwatch_list(), {}, "工单校验失败 → 不登记")

    def test_throttle_blocks_second_call(self):
        self._prs = [self._pr(9972, "feat/84291978-x")]
        self.sched._maybe_autoregister_open_prs()
        self._prs = [self._pr(5, "feat/55555555-y")]
        self.sched._maybe_autoregister_open_prs()  # 同 interval 内 → 节流
        reg = bot._prwatch_list()
        self.assertIn("84291978", reg)
        self.assertNotIn("55555555", reg, "同 interval 内第二次应被节流")

    def test_disabled_noop(self):
        self.sched._autoreg = False
        self._prs = [self._pr(9972, "feat/84291978-x")]
        self.sched._maybe_autoregister_open_prs()
        self.assertEqual(bot._prwatch_list(), {})


if __name__ == "__main__":
    unittest.main(verbosity=2)
