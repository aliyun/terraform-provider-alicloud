#!/usr/bin/env python3
"""Unit: pr_watch runner open-PR 全生命周期看守（post-PR 生命周期清单 #1/#2/#3）.

Gap (workflow assess-post-pr-loop → PR lifecycle handling):
提交 PR 后单次 headless 会话撑不住多日合并窗口——CI 转红没人修、reviewer 在 GitHub 上评论没人
回应、轮询频率跟不上。本测试锁住 pr_watch runner 作为「全生命周期看守器」的三块逻辑：
  · #1 _gh_pr_ci 解析 (failing/pending) + _maybe_dispatch_ci_fix 派发/去重/上限/active 信号；
  · #2 _gh_pr_comments 解析(排除自己/bot) + _maybe_dispatch_comment_reply baseline+新评论派发；
  · #3 active 信号驱动双档轮询（_maybe_dispatch_ci_fix 返回 True/False）。

Standalone: `python3 bridge/test_prwatch_ci_fix.py`. 无 gh/a1/网络（monkeypatch + fake control plane）。
"""
import json
import os
import sys
import tempfile
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE.parent))
from bridge import jarvis_task_router as router  # noqa: E402
from bridge.jarvis_task_router import EnqueueResult  # noqa: E402
from bridge.helpers import aone as events  # noqa: E402
from bridge.scheduler.runners import pr_watch as bot  # noqa: E402
from bridge.task_policy import policy_desired_revision  # noqa: E402

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
        self.sched = bot.PrWatchRuntime(FakeHandler(), FakePool())
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

    def test_empty_rollup_is_pending_unknown(self):
        self._patch({"headRefOid": "new-head", "statusCheckRollup": []})

        self.assertEqual(
            self.sched._gh_pr_ci(PR), ("new-head", [], True),
            "a head with no populated checks is not evidence that CI is green")

    def test_rollup_without_recognizable_check_is_pending_unknown(self):
        self._patch({"headRefOid": "new-head", "statusCheckRollup": [
            None, {"name": "metadata-only"}, "unexpected"]})

        self.assertEqual(
            self.sched._gh_pr_ci(PR), ("new-head", [], True))

    def test_query_failure_returns_none(self):
        bot.subprocess.run = lambda *a, **k: _fake_proc(1, "")
        self.assertEqual(self.sched._gh_pr_ci(PR), (None, None, False))


class GhPrCommentsParseTest(unittest.TestCase):
    def setUp(self):
        self.sched = bot.PrWatchRuntime(FakeHandler(), FakePool())
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
        tf.close()
        self.tmp = tf.name
        os.unlink(self.tmp)  # only the ".events" sibling path is used
        self._orig_prwatch_path = bot.PRWATCH_PATH
        bot.PRWATCH_PATH = Path(self.tmp + ".prwatch")
        bot._prwatch_add(TID, PR, PROJ, TITLE)
        self._orig_event_path = events.AONE_EVENT_PATH
        self._orig_bt = bot.broadcast_target
        self._orig_by = bot.broadcast_type
        self._orig_publish = events._aone_event_publish
        events.AONE_EVENT_PATH = Path(self.tmp + ".events")
        bot.broadcast_target = lambda: "t"
        bot.broadcast_type = lambda: "ty"
        self.events = []
        events._aone_event_publish = lambda *a, **k: self.events.append(a) or True
        self.handler = FakeHandler()
        self.pool = FakePool()
        self.handler.execution_router = FakeRouter(self.pool)
        self.sched = bot.PrWatchRuntime(self.handler, self.pool)
        self.sched._comment = lambda *a, **k: 0
        self.sched._escalate = lambda *a, **k: None

    def tearDown(self):
        bot.PRWATCH_PATH = self._orig_prwatch_path
        events.AONE_EVENT_PATH = self._orig_event_path
        bot.broadcast_target = self._orig_bt
        bot.broadcast_type = self._orig_by
        events._aone_event_publish = self._orig_publish
        try:
            os.unlink(self.tmp + ".events")
        except FileNotFoundError:
            pass
        try:
            os.unlink(self.tmp + ".prwatch")
        except FileNotFoundError:
            pass

    def _entry(self):
        return bot._prwatch_list()[TID]


class FakeAttentionClient:
    def __init__(self, notify=None):
        self.notify = notify
        self.current = {}
        self.upserts = []
        self.clears = []
        self.fail_clear = False

    def get_task_by_aone(self, aone_id):
        self.last_aone_id = aone_id
        return {"items": [
            {"id": 17, "generation": 1},
            {"taskId": 42, "generation": 2},
        ]}

    def upsert_task_attention(self, task_id, owner, event_key, payload):
        previous = self.current.get(str(task_id))
        projection = (str(owner), str(event_key))
        self.current[str(task_id)] = projection
        self.upserts.append((task_id, owner, event_key, payload))
        if isinstance(self.notify, list):
            value = self.notify.pop(0)
        elif self.notify is not None:
            value = self.notify
        else:
            value = previous != projection
        return {"notify": value, "task": {"id": task_id}}

    def clear_task_attention(self, task_id, *, event_key_prefix=None):
        if self.fail_clear:
            raise RuntimeError("control plane unavailable")
        self.clears.append((task_id, event_key_prefix))
        previous = self.current.get(str(task_id))
        if (event_key_prefix is None
                or (previous and previous[1].startswith(event_key_prefix))):
            self.current.pop(str(task_id), None)
        return {"notify": False, "task": {"id": task_id}}


class AttentionProjectionTest(_DispatchBase):
    def setUp(self):
        super().setUp()
        self.client = FakeAttentionClient()
        self.handler.task_client = self.client
        self.sched = bot.PrWatchRuntime(self.handler, self.pool)
        self.sched._gh_pr_state = lambda _url: ("OPEN", None)
        self.sched._gh_pr_comments = lambda _url: (None, None, None)
        self.notices = []
        self.sched._notify_attention = (
            lambda owner, payload: self.notices.append((owner, payload)))

    def test_green_pr_persists_attention_and_notifies_only_on_notify_true(self):
        self.sched._gh_pr_ci = lambda _url: ("green-head", [], False)

        self.sched._check_one(TID, self._entry())
        self.sched._check_one(TID, self._entry())

        self.assertEqual(self.client.last_aone_id, TID)
        self.assertEqual(len(self.client.upserts), 2)
        task_id, owner, event_key, payload = self.client.upserts[0]
        self.assertEqual(task_id, "42", "latest Task generation must own attention")
        self.assertEqual(owner, bot.master_staff())
        self.assertTrue(event_key.startswith("pr-review-merge:"))
        self.assertEqual(payload["kind"], "PR_REVIEW_MERGE")
        self.assertEqual(payload["head"], "green-head")
        self.assertEqual(
            payload["aoneUrl"],
            "https://project.aone.alibaba-inc.com/v2/project/%s/workitem/%s"
            % (PROJ, TID))
        self.assertEqual(payload["prUrl"], PR)
        self.assertIn("review", payload["action"])
        self.assertEqual(len(self.notices), 1,
                         "bridge must obey control-plane notify and keep no local ledger")
        self.assertEqual(self.events, [], "attention must not create an Aone comment event")

    def test_pending_or_failing_ci_clears_stale_review_attention(self):
        for ci in (("pending-head", [], True), ("failed-head", ["Compile"], False)):
            with self.subTest(ci=ci):
                self.sched._gh_pr_ci = lambda _url, value=ci: value
                self.sched._check_one(TID, self._entry())
        self.assertEqual(
            self.client.clears, [("42", "pr-"), ("42", "pr-")])
        self.assertEqual(self.client.upserts, [])
        self.assertEqual(self.notices, [])

    def test_escalated_pending_clears_attention_without_review_or_failure_notice(self):
        bot._prwatch_update(
            TID, ci_fix_sha="failed-head", ci_fix_attempts=3,
            ci_fix_escalated=True, ci_fix_escalated_head="failed-head")
        self.sched._gh_pr_ci = lambda _url: ("pending-head", [], True)

        active = self.sched._check_one(TID, self._entry())

        self.assertTrue(active, "pending must keep the fast polling interval")
        self.assertEqual(self.client.clears, [("42", "pr-")])
        self.assertEqual(self.client.upserts, [])
        self.assertEqual(self.notices, [])
        self.assertTrue(self._entry()["ci_fix_escalated"],
                        "pending is still part of the same unresolved CI epoch")

    def test_escalated_empty_rollup_keeps_epoch_and_fast_polling(self):
        bot._prwatch_update(
            TID, ci_fix_sha="failed-head", ci_fix_attempts=3,
            ci_fix_escalated=True, ci_fix_escalated_head="failed-head")
        self.client.current["42"] = (
            bot.master_staff(), "pr-review-merge:stale-attention")
        original_run = bot.subprocess.run
        self.addCleanup(setattr, bot.subprocess, "run", original_run)
        bot.subprocess.run = lambda *a, **k: _fake_proc(0, json.dumps({
            "headRefOid": "new-head", "statusCheckRollup": []}))

        active = self.sched._check_one(TID, self._entry())

        self.assertTrue(active, "an empty rollup must stay on the fast polling interval")
        self.assertEqual(self.client.clears, [("42", "pr-")],
                         "pending/unknown CI must clear stale PR attention")
        self.assertNotIn("42", self.client.current)
        self.assertEqual(self.client.upserts, [],
                         "unknown CI must project neither review nor failure attention")
        self.assertEqual(self.notices, [])
        self.assertEqual(self.events, [])
        entry = self._entry()
        self.assertEqual(entry["ci_fix_sha"], "failed-head")
        self.assertEqual(entry["ci_fix_attempts"], 3)
        self.assertTrue(entry["ci_fix_escalated"])
        self.assertEqual(entry["ci_fix_escalated_head"], "failed-head")

    def test_failure_after_pending_notifies_once_for_new_attention_epoch(self):
        bot._prwatch_update(
            TID, ci_fix_sha="s3", ci_fix_attempts=3,
            ci_fix_escalated=True, ci_fix_escalated_head="s4")

        self.sched._gh_pr_ci = lambda _url: ("s4", ["Compile"], False)
        self.sched._check_one(TID, self._entry())
        self.sched._gh_pr_ci = lambda _url: ("s4", [], True)
        self.sched._check_one(TID, self._entry())
        self.sched._gh_pr_ci = lambda _url: ("s4", ["Compile"], False)
        self.sched._check_one(TID, self._entry())
        self.sched._check_one(TID, self._entry())

        self.assertEqual(len(self.notices), 2,
                         "clear makes the returning failure a new attention epoch")
        self.assertEqual(len(self.client.clears), 1)
        failed_keys = [call[2] for call in self.client.upserts]
        self.assertEqual(len(set(failed_keys)), 1,
                         "the durable CI failure epoch key remains stable")

    def test_escalated_green_resets_epoch_and_projects_review_merge(self):
        bot._prwatch_update(
            TID, ci_fix_sha="failed-head", ci_fix_attempts=3,
            ci_fix_escalated=True, ci_fix_escalated_head="failed-head")
        self.sched._gh_pr_ci = lambda _url: ("green-head", [], False)

        active = self.sched._check_one(TID, self._entry())

        self.assertFalse(active)
        self.assertEqual(self.client.clears, [])
        self.assertEqual(len(self.client.upserts), 1)
        _task_id, _owner, key, payload = self.client.upserts[0]
        self.assertTrue(key.startswith("pr-review-merge:"))
        self.assertEqual(payload["kind"], "PR_REVIEW_MERGE")
        self.assertEqual(payload["head"], "green-head")
        entry = self._entry()
        self.assertEqual(entry["ci_fix_attempts"], 0)
        self.assertIsNone(entry["ci_fix_sha"])
        self.assertFalse(entry["ci_fix_escalated"])
        self.assertIsNone(entry["ci_fix_escalated_head"])

    def test_real_failure_at_limit_notifies_once_across_restart(self):
        bot._prwatch_update(TID, ci_fix_sha="s3", ci_fix_attempts=3)
        self.sched._gh_pr_ci = lambda _url: ("s4", ["Compile"], False)

        self.sched._check_one(TID, self._entry())
        persisted = self._entry()
        self.assertTrue(persisted["ci_fix_escalated"])
        self.assertEqual(persisted["ci_fix_escalated_head"], "s4")

        # Recreate the scheduler to prove the JSON registry state is sufficient to
        # restore dedup after a process restart. Control plane returns notify=False for
        # the same semantic event key on the second projection.
        restarted = bot.PrWatchRuntime(self.handler, self.pool)
        restarted._gh_pr_state = lambda _url: ("OPEN", None)
        restarted._gh_pr_ci = lambda _url: ("s4", ["Compile"], False)
        restarted._gh_pr_comments = lambda _url: (None, None, None)
        restarted._notify_attention = self.sched._notify_attention
        restarted._check_one(TID, self._entry())

        self.assertEqual(len(self.pool.submitted), 0)
        self.assertEqual(len(self.client.upserts), 2)
        self.assertEqual(self.client.upserts[0][2], self.client.upserts[1][2],
                         "same failure epoch must reuse one attention event key")
        self.assertEqual(
            [call[3]["failingChecks"] for call in self.client.upserts],
            [["Compile"], ["Compile"]])
        self.assertEqual(len(self.notices), 1)
        self.assertEqual(len(self.events), 1,
                         "exhaustion Aone event must also remain idempotent")

    def test_persisted_legacy_escalation_only_alerts_on_real_failure(self):
        bot._prwatch_update(
            TID, ci_fix_sha="old-head", ci_fix_attempts=3,
            ci_fix_escalated=True)
        self.sched._gh_pr_ci = lambda _url: ("current-failing-head", ["Test"], False)

        self.sched._check_one(TID, self._entry())

        self.assertEqual(len(self.client.upserts), 1)
        _task_id, _owner, _key, payload = self.client.upserts[0]
        self.assertEqual(payload["kind"], "PR_CI_FAILED")
        self.assertEqual(payload["failingChecks"], ["Test"])
        self.assertEqual(
            self._entry()["ci_fix_escalated_head"], "current-failing-head")

    def test_merged_clear_failure_keeps_watch_for_next_tick(self):
        self.client.fail_clear = True
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._finish = lambda *_a: self.fail(
            "lifecycle must wait until attention clear converges")

        self.sched._check_one(TID, self._entry())

        self.assertTrue(bot._prwatch_has(TID))

    def test_closed_pr_converts_attention_before_unwatch(self):
        self.sched._gh_pr_state = lambda _url: ("CLOSED", None)
        self.sched._check_one(TID, self._entry())

        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(len(self.client.upserts), 1)
        _task_id, _owner, key, payload = self.client.upserts[0]
        self.assertTrue(key.startswith("pr-closed:"))
        self.assertEqual(payload["kind"], "PR_CLOSED_DECISION")
        self.assertEqual(len(self.notices), 1)


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
        envelope = self.pool.submitted[0]["envelope"]
        self.assertEqual(
            envelope.desired_revision,
            policy_desired_revision("pr-ci:sha1", envelope.payload))
        self.assertNotIn(TITLE, envelope.desired_revision)
        e = self._entry()
        self.assertEqual((e["ci_fix_sha"], e["ci_fix_attempts"]), ("sha1", 1))
        self.assertEqual(self.events, [], "单次 CI 修复派发不更新 Aone")

    def test_check_migrates_legacy_registry_title_before_pr_dispatch(self):
        legacy = self._entry()
        legacy.pop("title", None)
        bot.PRWATCH_PATH.write_text(json.dumps({TID: legacy}))
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
        self.assertEqual(self._entry().get("ci_fix_escalated_head"), "s4")
        self.assertEqual(comments, [], "Terraform 重要事件必须走统一 publisher，不走 legacy _comment")
        self.assertEqual(len(self.events), 1)
        self.assertIn("ci-exhausted:s4:3", self.events[0][2])
        self._ci = ("s5", ["Compile"], False)
        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 3)
        self.assertFalse(active, "escalated 后转人工 → 不再快档")

    def test_green_resets_persisted_failure_epoch_and_retry_budget(self):
        bot._prwatch_update(
            TID, ci_fix_sha="old-head", ci_fix_attempts=3,
            ci_fix_escalated=True, ci_fix_escalated_head="failed-head")
        self._ci = ("green-head", [], False)

        active = self.sched._maybe_dispatch_ci_fix(TID, self._entry())

        self.assertFalse(active)
        entry = self._entry()
        self.assertEqual(entry["ci_fix_attempts"], 0)
        self.assertIsNone(entry["ci_fix_sha"])
        self.assertFalse(entry["ci_fix_escalated"])
        self.assertIsNone(entry["ci_fix_escalated_head"])

        self._ci = ("future-failure", ["Compile"], False)
        self.sched._maybe_dispatch_ci_fix(TID, self._entry())
        self.assertEqual(len(self.pool.submitted), 1)
        self.assertEqual(self._entry()["ci_fix_attempts"], 1)


class EscalateMessageTest(_DispatchBase):
    """_escalate 消息格式：pr_url 非空时工单号渲染成指向对应 PR 的 markdown 链接（钉钉 AI
    卡片按 markdown 渲染，[text](url) 可点）；缺省回退纯文本 #<tid>。诉求：escalate 消息
    "🚩 需人工介入 #<id>" 里的 #<id> 须带可点击链接指向对应 PR，而非仅展示 id 号。"""

    def setUp(self):
        super().setUp()
        # _DispatchBase.setUp 把 _escalate mock 成 noop；恢复真实实现以校验消息格式。
        self.sched._escalate = bot.PrWatchRuntime._escalate.__get__(self.sched)

    def test_pr_url_renders_clickable_link_to_pr(self):
        self.sched._escalate(TID, "PR CI 反复失败超过自动修复上限(3)，请人工介入", PR)
        self.assertEqual(len(self.handler.broadcasts), 1)
        msg = self.handler.broadcasts[0]
        self.assertIn("[#%s](%s)" % (TID, PR), msg,
                      "工单号须渲染成指向对应 PR 的 markdown 可点击链接")
        self.assertIn("请人工介入", msg, "reason 正文保留")

    def test_missing_pr_url_falls_back_to_plain_id(self):
        self.sched._escalate(TID, "PR 未合并即关闭，请人工确认工单去向")
        self.assertEqual(len(self.handler.broadcasts), 1)
        msg = self.handler.broadcasts[0]
        self.assertIn("#%s" % TID, msg, "缺省仍展示工单号")
        self.assertNotIn("](http", msg, "无 pr_url 时回退纯文本，不渲染 markdown 链接")


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
        envelope = self.pool.submitted[0]["envelope"]
        self.assertEqual(
            envelope.desired_revision,
            policy_desired_revision("pr-comment:pr-2", envelope.payload))
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
        self._orig_run = bot.run_process_group

    def tearDown(self):
        bot.run_process_group = self._orig_run
        super().tearDown()

    def test_finish_uses_rd_but_terraform_comment_is_suppressed(self):
        calls = []

        def fake(cmd, *a, **kw):
            calls.append((list(cmd), dict(kw.get("env") or {})))
            return _fake_proc(0, "")

        bot.run_process_group = fake
        self.sched._comment = bot.PrWatchRuntime._comment.__get__(
            self.sched, bot.PrWatchRuntime)
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

    @staticmethod
    def _status_snapshot(project="1086837", workitem_type="3",
                         status="问题解决中", status_id="155741", tags=()):
        return {
            "project": str(project),
            "workitem_type": str(workitem_type),
            "workitem_type_name": "需求问题" if str(workitem_type) == "3" else "其它类型",
            "status": status,
            "status_id": str(status_id),
            "tags": set(tags),
        }

    def test_customer_type3_merged_updates_status_and_uses_rd(self):
        bot._prwatch_update(TID, project="1086837")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        snapshots = iter((
            self._status_snapshot(),
            self._status_snapshot(status="已合入主线", status_id="626904"),
        ))
        self.sched._workitem_snapshot = lambda _tid: next(snapshots)
        self.sched._finish = lambda *_a: self.fail("customer type 3 must not finish")
        calls = []

        def fake(cmd, *a, **kw):
            calls.append((list(cmd), dict(kw.get("env") or {})))
            return _fake_proc(0, "")

        bot.run_process_group = fake
        self.sched._check_one(TID, self._entry())

        update, env = next(call for call in calls if "update" in call[0])
        self.assertEqual(update[update.index("--status") + 1], "已合入主线")
        self.assertNotIn("--tag", update)
        self.assertEqual(env.get("JARVIS_A1_IDENTITY"), "terraform-rd")
        self.assertEqual(env.get("JARVIS_A1_STRICT"), "1")
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(len(self.events), 1)
        self.assertIn(":merged:", self.events[0][2])

    def test_customer_point_read_parses_authoritative_scope_type_status_and_tags(self):
        payload = {"fields": [
            {"identifier": "space", "value": "1086837",
             "displayValue": "Terraform - 客户需求"},
            {"identifier": "workitemType", "value": "3",
             "displayValue": "需求问题"},
            {"identifier": "status", "value": "626904",
             "displayValue": "已合入主线"},
            {"identifier": "tag", "value": "11,22",
             "displayValue": "jarvis-idle,customer-tag"},
        ]}
        calls = []

        def fake(cmd, *a, **kw):
            calls.append((list(cmd), dict(kw.get("env") or {})))
            return _fake_proc(0, json.dumps(payload))

        bot.run_process_group = fake
        snapshot = bot.PrWatchRuntime._workitem_snapshot(self.sched, TID)
        self.assertEqual(snapshot["project"], "1086837")
        self.assertEqual(snapshot["workitem_type"], "3")
        self.assertEqual(snapshot["workitem_type_name"], "需求问题")
        self.assertEqual((snapshot["status"], snapshot["status_id"]),
                         ("已合入主线", "626904"))
        self.assertEqual(snapshot["tags"], {"jarvis-idle", "customer-tag"})
        self.assertEqual(calls[0][1].get("JARVIS_A1_IDENTITY"), "terraform-rd")
        self.assertEqual(calls[0][1].get("JARVIS_A1_STRICT"), "1")

    def test_customer_type3_already_merged_status_is_idempotent(self):
        bot._prwatch_update(TID, project="1086837")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._workitem_snapshot = lambda _tid: self._status_snapshot(
            status="已合入主线", status_id="626904")
        self.sched._finish = lambda *_a: self.fail("customer type 3 must not finish")
        bot.run_process_group = lambda *_a, **_k: self.fail(
            "already-confirmed merged status must not update")
        self.sched._check_one(TID, self._entry())
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(len(self.events), 1)

    def test_customer_initial_read_failure_keeps_watch(self):
        bot._prwatch_update(TID, project="1086837")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._workitem_snapshot = lambda _tid: (_ for _ in ()).throw(
            RuntimeError("read failed"))
        self.sched._finish = lambda *_a: self.fail("read failure must not finish")
        self.sched._check_one(TID, self._entry())
        self.assertTrue(bot._prwatch_has(TID))
        self.assertEqual(self.events, [])

    def test_customer_status_update_failure_keeps_watch(self):
        bot._prwatch_update(TID, project="1086837")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._workitem_snapshot = lambda _tid: self._status_snapshot()
        self.sched._finish = lambda *_a: self.fail("update failure must not finish")
        bot.run_process_group = lambda *_a, **_k: _fake_proc(1, "update failed")
        self.sched._check_one(TID, self._entry())
        self.assertTrue(bot._prwatch_has(TID))
        self.assertEqual(self.events, [])

    def test_customer_status_readback_failure_keeps_watch(self):
        bot._prwatch_update(TID, project="1086837")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        snapshots = iter((self._status_snapshot(), RuntimeError("readback failed")))

        def snapshot(_tid):
            value = next(snapshots)
            if isinstance(value, Exception):
                raise value
            return value

        self.sched._workitem_snapshot = snapshot
        self.sched._finish = lambda *_a: self.fail("readback failure must not finish")
        bot.run_process_group = lambda *_a, **_k: _fake_proc(0, "")
        self.sched._check_one(TID, self._entry())
        self.assertTrue(bot._prwatch_has(TID))
        self.assertEqual(self.events, [])

    def test_customer_status_readback_requires_exact_scope_type_and_status(self):
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        mismatches = (
            self._status_snapshot(project="528766", status="已合入主线",
                                  status_id="626904"),
            self._status_snapshot(workitem_type="36", status="已合入主线",
                                  status_id="626904"),
            self._status_snapshot(status="已合入主线", status_id="wrong"),
            self._status_snapshot(status="问题解决中", status_id="155741"),
        )
        for index, after in enumerate(mismatches):
            with self.subTest(after=after):
                if not bot._prwatch_has(TID):
                    bot._prwatch_add(TID, PR, "1086837", TITLE)
                bot._prwatch_update(TID, project="1086837")
                self.events.clear()
                snapshots = iter((self._status_snapshot(), after))
                self.sched._workitem_snapshot = lambda _tid, it=snapshots: next(it)
                self.sched._finish = lambda *_a: self.fail(
                    "invalid readback must not finish")
                bot.run_process_group = lambda *_a, **_k: _fake_proc(0, "")
                self.sched._check_one(TID, self._entry())
                self.assertTrue(bot._prwatch_has(TID), index)
                self.assertEqual(self.events, [])

    def test_customer_merged_event_failure_keeps_watch_after_status(self):
        bot._prwatch_update(TID, project="1086837")
        events._aone_event_publish = lambda *_a, **_k: False
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._workitem_snapshot = lambda _tid: self._status_snapshot(
            status="已合入主线", status_id="626904")
        self.sched._finish = lambda *_a: self.fail("customer type 3 must not finish")
        self.sched._check_one(TID, self._entry())
        self.assertTrue(bot._prwatch_has(TID))

    def test_customer_legacy_finish_state_still_converges_status_first(self):
        bot._prwatch_update(TID, project="1086837", finish_succeeded=True)
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._workitem_snapshot = lambda _tid: self._status_snapshot(
            status="已合入主线", status_id="626904")
        self.sched._ticket_guard = lambda *_a: self.fail("status path must bypass guard")
        self.sched._check_one(TID, self._entry())
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(len(self.events), 1)

    def test_customer_non_type3_preserves_legacy_finish_for_all_other_types(self):
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._ticket_guard = lambda _tid: "ok"
        finishes = []
        self.sched._finish = lambda *args: finishes.append(args) or 0
        for item_type in ("36", "38", "27", "349"):
            with self.subTest(item_type=item_type):
                if not bot._prwatch_has(TID):
                    bot._prwatch_add(TID, PR, "1086837", TITLE)
                bot._prwatch_update(TID, project="1086837")
                self.events.clear()
                self.sched._workitem_snapshot = lambda _tid, t=item_type: (
                    self._status_snapshot(workitem_type=t))
                self.sched._check_one(TID, self._entry())
                self.assertFalse(bot._prwatch_has(TID))
                self.assertEqual(finishes[-1], (TID, "1086837", "已完成"))
                self.assertEqual(len(self.events), 1)

    def test_customer_registry_with_non_customer_point_read_preserves_legacy_path(self):
        bot._prwatch_update(TID, project="1086837")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._workitem_snapshot = lambda _tid: self._status_snapshot(
            project="528766")
        self.sched._ticket_guard = lambda _tid: "ok"
        finishes = []
        self.sched._finish = lambda *args: finishes.append(args) or 0
        self.sched._check_one(TID, self._entry())
        self.assertEqual(finishes, [(TID, "1086837", "已完成")])
        self.assertFalse(bot._prwatch_has(TID))

    def test_customer_type3_terminal_guard_does_not_move_status(self):
        bot._prwatch_update(TID, project="1086837")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._workitem_snapshot = lambda _tid: self._status_snapshot(
            status="已发布", status_id="terminal")
        self.sched._finish = lambda *_a: self.fail("terminal customer must not finish")
        bot.run_process_group = lambda *_a, **_k: self.fail(
            "terminal customer must not update status")
        self.sched._check_one(TID, self._entry())
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(len(self.events), 1)

    def test_customer_type3_npe_guard_does_not_move_status(self):
        bot._prwatch_update(TID, project="1086837")
        self.sched._gh_pr_state = lambda _url: ("MERGED", "2026-07-01T00:00:00Z")
        self.sched._workitem_snapshot = lambda _tid: self._status_snapshot(
            tags=("jarvis-npe",))
        self.sched._finish = lambda *_a: self.fail("npe customer must not finish")
        bot.run_process_group = lambda *_a, **_k: self.fail(
            "npe customer must not update status")
        escalated = []
        self.sched._escalate = lambda *_a: escalated.append(_a)
        self.sched._check_one(TID, self._entry())
        self.assertFalse(bot._prwatch_has(TID))
        self.assertEqual(len(self.events), 1)
        self.assertIn(":merged-npe:", self.events[0][2])
        self.assertEqual(len(escalated), 1)


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
        self.assertEqual(self.handler.broadcasts, [],
                         "merged event 已持久化到 Aone，不再群播 routine 状态")

    def test_terraform_finish_keeps_watch_if_event_not_durable(self):
        events._aone_event_publish = lambda *_a, **_k: False
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
        self.assertEqual(self.handler.broadcasts, [],
                         "补偿收尾成功后不再重复群播")

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
        bot._prwatch_remove(TID)  # 清空 base 的 TID entry，从零测漏登发现
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


class NotifyAttentionImportTest(unittest.TestCase):
    """Regression: `_notify_attention` must resolve the notifier it calls.

    The scheduler runner split left `pr_watch._notify_attention` calling
    `_notify_task_attention` without importing it, so every attention notice
    raised NameError.  `_TaskAttentionPublisher.upsert` catches notifier
    exceptions best-effort, so the board flipped to 「转人工」 while Aone and
    DingTalk stayed silent and the Task still reported success.

    Every other attention test replaces `_notify_attention` in setUp, so the
    real body never executed and the missing name survived.  This test calls it
    directly and only stubs the outbound subprocess, keeping the name
    resolution itself under test.
    """

    def test_notify_attention_reaches_the_real_notifier(self):
        argv_seen = []

        class _RecordingSubprocess:
            @staticmethod
            def run(argv, **kwargs):
                argv_seen.append([str(part) for part in argv])
                return None

        original = router.subprocess
        router.subprocess = _RecordingSubprocess
        try:
            bot.PrWatchRuntime._notify_attention(
                "320687",
                {"reason": "CI 已通过，等待人工 review/merge。",
                 "action": "请 review 代码；确认无误后合并 PR。",
                 "aoneUrl": "https://example.invalid/workitem",
                 "prUrl": PR})
        finally:
            router.subprocess = original

        self.assertEqual(
            len(argv_seen), 1,
            "the real notifier must be reached, not silently skipped")
        self.assertIn("notify-dingtalk.sh", " ".join(argv_seen[0]))
        self.assertIn("320687", argv_seen[0])


if __name__ == "__main__":
    unittest.main(verbosity=2)
