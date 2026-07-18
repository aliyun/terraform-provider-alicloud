#!/usr/bin/env python3
"""Unit: PrWatchScheduler open-PR 全生命周期看守（post-PR 生命周期清单 #1/#2/#3）.

Gap (workflow assess-post-pr-loop → escalation/cap-post-pr-lifecycle.md):
提交 PR 后单次 headless 会话撑不住多日合并窗口——CI 转红没人修、reviewer 在 GitHub 上评论没人
回应、轮询频率跟不上。本测试锁住 PrWatchScheduler 升级为「全生命周期看守器」后的三块逻辑：
  · #1 _gh_pr_ci 解析 (failing/pending) + _maybe_dispatch_ci_fix 派发/去重/上限/active 信号；
  · #2 _gh_pr_comments 解析(排除自己/bot) + _maybe_dispatch_comment_reply baseline+新评论派发；
  · #3 active 信号驱动双档轮询（_maybe_dispatch_ci_fix 返回 True/False）。

Standalone: `python3 bridge/test_prwatch_ci_fix.py`. 无 gh/a1/网络（monkeypatch + fake control plane）。
"""
import importlib.util
import json
import os
import sys
import tempfile
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)
from jarvis_task_router import EnqueueResult

TID = "84251052"
PR = "https://github.com/aliyun/terraform-provider-alicloud/pull/9972"
PROJ = "528766"
TITLE = "Aone workitem title"


class FakeHandler:
    def __init__(self):
        self.broadcasts = []
        self.dispatched = []
        self.execution_router = FakeRouter()

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
        raise AssertionError("recoverable PR-watch Task must not enter EphemeralExecutor")


class FakeRouter:
    def __init__(self, sink=None):
        self.sink = sink

    def enqueue(self, envelope, local_submit=None):
        if self.sink is not None:
            self.sink.submitted.append({
                "key": envelope.task_key,
                "force": True,
                "kind": envelope.task_type,
                "terraform": bool(envelope.payload.get("terraform")),
                "envelope": envelope,
            })
            if not self.sink.accept:
                return EnqueueResult(False, "control_plane_rejected")
        return EnqueueResult(True, "task_persisted")


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
        self._orig_root = bot.REPO_ROOT
        self._tmp = tempfile.TemporaryDirectory()
        bot.REPO_ROOT = Path(self._tmp.name)
        (bot.REPO_ROOT / "config").mkdir()
        (bot.REPO_ROOT / "config" / "contacts.json").write_text(
            json.dumps({"contacts": [
                {"name": "ours", "id": "1", "github": "TeamMember"},
            ]}))

    def tearDown(self):
        bot.subprocess.run = self._orig_run
        bot.REPO_ROOT = self._orig_root
        self._tmp.cleanup()

    def _patch(self, issues=None, lines=None, reviews=None):
        payloads = {
            "/issues/": issues or [],
            "/pulls/9972/comments": lines or [],
            "/pulls/9972/reviews": reviews or [],
        }

        def fake(cmd, *a, **k):
            path = str(cmd[-1])
            for marker, payload in payloads.items():
                if marker in path:
                    return _fake_proc(0, json.dumps(payload))
            return _fake_proc(1, "")

        bot.subprocess.run = fake

    def test_latest_reviewer_comment_picked(self):
        self._patch(
            issues=[
                {"id": 10, "user": {"login": "reviewer1"}, "body": "please fix X",
                 "created_at": "2026-07-01T00:00:00Z"}],
            lines=[
                {"id": 20, "user": {"login": "reviewer2"}, "body": "and Y",
                 "created_at": "2026-07-01T02:00:00Z"}],
            reviews=[
                {"id": 30, "user": {"login": "reviewer3"}, "body": "earlier review",
                 "submitted_at": "2026-07-01T01:00:00Z"}])
        key, author, snippet = self.sched._gh_pr_comments(PR)
        self.assertEqual((key, author), ("pr-20", "reviewer2"))
        self.assertIn("Y", snippet)

    def test_three_routes_and_dynamic_self_logins_are_filtered(self):
        self._patch(
            issues=[
                {"id": 1, "user": {"login": "reviewer1"}, "body": "real",
                 "created_at": "2026-07-01T00:00:00Z"},
                {"id": 2, "user": {"login": "TEAMMEMBER"}, "body": "our reply",
                 "created_at": "2026-07-01T04:00:00Z"}],
            lines=[
                {"id": 3, "user": {"login": "github-actions[bot]"}, "body": "ci",
                 "created_at": "2026-07-01T05:00:00Z"},
                {"id": 4, "user": {"login": "api-tool-agent"}, "body": "self",
                 "created_at": "2026-07-01T06:00:00Z"}],
            reviews=[
                {"id": 5, "user": {"login": "reviewer2"}, "body": "review body",
                 "submitted_at": "2026-07-01T03:00:00Z"},
                {"id": 6, "user": {"login": "reviewer3"}, "body": "   ",
                 "submitted_at": "2026-07-01T07:00:00Z"}])
        key, author, _ = self.sched._gh_pr_comments(PR)
        self.assertEqual((key, author), ("review-5", "reviewer2"),
                         "三路都要查；动态 contacts/self/bot/空 review body 都要过滤")

    def test_no_reviewer_comment_returns_none(self):
        self._patch(issues=[
            {"id": 1, "user": {"login": "api-tool-agent"}, "body": "x",
             "created_at": "2026-07-01T00:00:00Z"}])
        self.assertEqual(self.sched._gh_pr_comments(PR), (None, None, None))

    def test_query_failure_returns_none(self):
        bot.subprocess.run = lambda *a, **k: _fake_proc(1, "")
        self.assertEqual(self.sched._gh_pr_comments(PR), (None, None, None))


class _DispatchBase(unittest.TestCase):
    def setUp(self):
        tf = tempfile.NamedTemporaryFile(suffix=".json", delete=False, mode="w")
        tf.write(json.dumps({TID: {"pr_url": PR, "project": PROJ,
                                  "title": TITLE, "submitted_at": "x"}}))
        tf.close()
        self.tmp = tf.name
        self._orig_path = bot.PRWATCH_PATH
        self._orig_event_path = bot.AONE_EVENT_PATH
        self._orig_bt = bot.broadcast_target
        self._orig_by = bot.broadcast_type
        self._orig_publish = bot._aone_event_publish
        bot.PRWATCH_PATH = Path(self.tmp)
        bot.AONE_EVENT_PATH = Path(self.tmp + ".events")
        bot.broadcast_target = lambda: "t"
        bot.broadcast_type = lambda: "ty"
        self.events = []
        bot._aone_event_publish = lambda *a: self.events.append(a) or True
        self.handler = FakeHandler()
        self.pool = FakePool()
        self.handler.execution_router = FakeRouter(self.pool)
        self.sched = bot.PrWatchScheduler(self.handler, self.pool)
        self.sched._comment = lambda *a, **k: 0
        self.sched._escalate = lambda *a, **k: None

    def tearDown(self):
        bot.PRWATCH_PATH = self._orig_path
        bot.AONE_EVENT_PATH = self._orig_event_path
        bot.broadcast_target = self._orig_bt
        bot.broadcast_type = self._orig_by
        bot._aone_event_publish = self._orig_publish
        os.unlink(self.tmp)
        try:
            os.unlink(self.tmp + ".events")
        except FileNotFoundError:
            pass

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
        self.assertTrue(self.pool.submitted[0]["terraform"])
        self.assertEqual(
            self.pool.submitted[0]["envelope"].aone_id, TID,
            "GITHUB PR CI trigger must retain the canonical Aone association")
        self.assertEqual(
            self.pool.submitted[0]["envelope"].source_ref["title"], TITLE)
        self.assertNotIn(TITLE, self.pool.submitted[0]["envelope"].desired_revision)
        e = self._entry()
        self.assertEqual((e["ci_fix_sha"], e["ci_fix_attempts"]), ("sha1", 1))
        self.assertEqual(self.events, [], "单次 CI 修复派发不更新 Aone")

    def test_check_migrates_legacy_registry_title_before_pr_dispatch(self):
        legacy = self._entry()
        legacy.pop("title", None)
        bot._prwatch_write({TID: legacy})
        self.sched._ticket_metadata = lambda tid: (PROJ, "Backfilled Aone title")
        self.sched._gh_pr_state = lambda _url: ("OPEN", None)
        self._ci = ("sha1", ["Compile"], False)
        self.sched._gh_pr_comments = lambda _url: (None, None, None)

        self.sched._check_one(TID, self._entry())

        self.assertEqual(self._entry()["title"], "Backfilled Aone title")
        self.assertEqual(
            self.pool.submitted[0]["envelope"].source_ref["title"],
            "Backfilled Aone title")

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
        self.assertEqual(self.events, [], "new head / 单次 retry 仍静默")

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
        self.assertEqual(self.events, [], "CI pending 不更新 Aone")

    def test_query_fail_no_dispatch_not_active(self):
        self._ci = (None, None, False)
        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0)
        self.assertFalse(active)

    def test_exhaust_attempts_escalates_once_then_stops(self):
        comments = []
        self.sched._comment = lambda *a, **k: comments.append(a) or 0
        for sha in ("s1", "s2", "s3", "s4"):
            self._ci = (sha, ["Compile"], False)
            self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 3, "自动修复上限 3 次")
        self.assertTrue(self._entry().get("ci_fix_escalated"))
        self.assertEqual(comments, [], "Terraform 重要事件必须走统一 publisher，不走 legacy _comment")
        self.assertEqual(len(self.events), 1)
        self.assertIn("ci-exhausted:s4:3", self.events[0][2])
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
        self._c = ("issue-1", "reviewer1", "please fix")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0, "首次观察只 baseline，不回应既有评论")
        self.assertEqual(self._entry().get("last_seen_comment"), "issue-1")

    def test_new_comment_dispatches(self):
        self._c = ("issue-1", "reviewer1", "c1")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # baseline
        self._c = ("pr-2", "reviewer2", "please change Y")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # new
        self.assertEqual(len(self.pool.submitted), 1)
        self.assertEqual(self.pool.submitted[0]["kind"], "pr_comment_reply")
        self.assertTrue(self.pool.submitted[0]["force"])
        self.assertTrue(self.pool.submitted[0]["terraform"])
        self.assertEqual(
            self.pool.submitted[0]["envelope"].aone_id, TID,
            "GITHUB PR comment trigger must retain the canonical Aone association")
        self.assertEqual(
            self.pool.submitted[0]["envelope"].source_ref["title"], TITLE)
        self.assertEqual(self._entry().get("last_seen_comment"), "pr-2")
        self.assertEqual(self.events, [], "普通 reviewer comment 仅在 GitHub 内处理")

    def test_same_comment_not_redispatched(self):
        self._c = ("issue-1", "reviewer1", "c1")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # baseline u1
        self._c = ("review-2", "reviewer2", "c2")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # dispatch u2
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())  # u2 again → skip
        self.assertEqual(len(self.pool.submitted), 1, "同一条最新评论不重复回应")

    def test_legacy_issue_url_baseline_is_silently_upgraded(self):
        bot._prwatch_update(TID, last_seen_comment=(
            "https://github.com/aliyun/terraform-provider-alicloud/pull/9972"
            "#issuecomment-123"))
        self._c = ("issue-123", "reviewer1", "same old comment")
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0)
        self.assertEqual(self._entry().get("last_seen_comment"), "issue-123")

    def test_no_reviewer_comment_no_dispatch(self):
        self._c = (None, None, None)
        self.sched._maybe_dispatch_comment_reply(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 0)
        self.assertIsNone(self._entry().get("last_seen_comment"))


class PrWatchWriteIdentityAndOutcomeTest(_DispatchBase):
    def setUp(self):
        super().setUp()
        self._orig_run = bot.subprocess.run

    def tearDown(self):
        bot.subprocess.run = self._orig_run
        super().tearDown()

    def test_finish_uses_rd_but_terraform_comment_is_suppressed(self):
        calls = []

        def fake(cmd, *a, **kw):
            calls.append((list(cmd), dict(kw.get("env") or {})))
            return _fake_proc(0, "")

        bot.subprocess.run = fake
        self.sched._comment = bot.PrWatchScheduler._comment.__get__(
            self.sched, bot.PrWatchScheduler)
        self.sched._finish(TID, PROJ, "已完成")
        self.sched._comment(TID, PROJ, "done")
        self.sched._comment(TID, "2124589", "non tf")
        self.assertEqual(len(calls), 2, "Terraform _comment must not spawn wrap.sh")
        self.assertEqual(calls[0][1].get("JARVIS_A1_IDENTITY"), "terraform-rd")
        self.assertEqual(calls[0][1].get("JARVIS_A1_STRICT"), "1")
        self.assertNotEqual(calls[1][1].get("JARVIS_A1_IDENTITY"), "terraform-rd")
        self.assertIn("wrap.sh", calls[1][0][0], "非 Terraform 仍保留 wrap comment")

    def test_finish_any_nonzero_keeps_watch_without_success_broadcast(self):
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._ticket_guard = lambda _tid: "ok"
        self.sched._finish = lambda *_a: 1
        comments = []
        self.sched._comment = lambda *_a: comments.append(_a) or 0
        self.sched._check_one(TID, self._entry())
        self.assertTrue(bot._prwatch_has(TID))
        self.assertEqual(comments, [])
        self.assertEqual(self.handler.broadcasts, [])

    def test_terraform_finish_publishes_one_merged_event_and_removes(self):
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._ticket_guard = lambda _tid: "ok"
        self.sched._finish = lambda *_a: 0
        comments = []
        self.sched._comment = lambda *_a: comments.append(_a) or 1
        self.sched._check_one(TID, self._entry())
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(comments, [])
        self.assertEqual(len(self.events), 1)
        self.assertIn(":merged:", self.events[0][2])
        self.assertEqual(len(self.handler.broadcasts), 1)

    def test_terraform_finish_keeps_watch_if_event_not_durable(self):
        bot._aone_event_publish = lambda *_a: False
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._ticket_guard = lambda _tid: "ok"
        self.sched._finish = lambda *_a: 0
        self.sched._check_one(TID, self._entry())
        self.assertTrue(bot._prwatch_has(TID))
        self.assertTrue(self._entry().get("finish_succeeded"),
                        "finish 已落地时必须写补偿态，避免下轮重复 finish")
        self.assertEqual(self.handler.broadcasts, [])

    def test_legacy_terraform_finish_succeeded_state_publishes_merged_event(self):
        bot._prwatch_update(TID, finish_succeeded=True)
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._ticket_guard = lambda _tid: self.fail("pending comment must bypass guard")
        comments = []
        self.sched._comment = lambda *_a: comments.append(_a) or 0
        self.sched._check_one(TID, self._entry())
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(comments, [])
        self.assertEqual(len(self.events), 1)
        self.assertEqual(len(self.handler.broadcasts), 1)

    def test_terraform_closed_pr_publishes_important_event(self):
        self.sched._gh_pr_state = lambda _url: ("CLOSED", None)
        comments = []
        self.sched._comment = lambda *_a: comments.append(_a) or 1
        escalated = []
        self.sched._escalate = lambda *_a: escalated.append(_a)
        self.sched._check_one(TID, self._entry())
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(comments, [])
        self.assertEqual(len(self.events), 1)
        self.assertTrue(self.events[0][2].endswith(":closed"))
        self.assertEqual(len(escalated), 1)

    def test_terraform_merged_with_npe_publishes_distinct_event(self):
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._ticket_guard = lambda _tid: "npe"
        escalated = []
        self.sched._escalate = lambda *_a: escalated.append(_a)
        self.sched._check_one(TID, self._entry())
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(len(self.events), 1)
        self.assertIn(":merged-npe:", self.events[0][2])
        self.assertEqual(len(escalated), 1)

    def test_non_terraform_merged_keeps_comment_compensation(self):
        bot._prwatch_update(TID, project="2124589")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._ticket_guard = lambda _tid: "ok"
        self.sched._finish = lambda *_a: 0
        comments = []
        self.sched._comment = lambda *_a: comments.append(_a) or 0
        self.sched._check_one(TID, self._entry())
        self.assertEqual(len(comments), 1)
        self.assertEqual(self.events, [])
        self.assertFalse(bot._prwatch_has(TID))


class AutoRegisterTest(_DispatchBase):
    """#6 兜底发现：漏登记的 open PR 自动补登记（分支解析工单号 + aone-get 校验）。"""

    def setUp(self):
        super().setUp()
        bot._prwatch_write({})  # 清空 base 的 TID entry，从零测漏登发现
        self._prs = []
        self._proj = "528766"
        self.sched._gh_open_prs = lambda: self._prs
        self.sched._ticket_metadata = lambda tid: (self._proj, TITLE)

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
        self.assertEqual(reg["84291978"]["title"], TITLE)

    def test_already_watched_skipped(self):
        pr = self._pr(9972, "feat/84291978-tair")
        bot._prwatch_add("84291978", pr["url"], "528766")
        self._prs = [pr]
        called = []
        self.sched._ticket_metadata = lambda tid: (called.append(tid) or "528766", TITLE)
        self.sched._maybe_autoregister_open_prs()
        self.assertEqual(called, [], "已登记 PR 不应再查 project/重登")

    def test_registry_freezes_first_aone_title_and_never_uses_pr_title(self):
        pr = self._pr(9972, "feat/84291978-tair")
        bot._prwatch_add("84291978", pr["url"], "528766", "Original Aone title")
        bot._prwatch_add("84291978", pr["url"], "528766", "Renamed Aone title")
        self.assertEqual(bot._prwatch_list()["84291978"]["title"],
                         "Original Aone title")
        self.assertNotEqual(bot._prwatch_list()["84291978"]["title"],
                            "PR 9972 title")

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
