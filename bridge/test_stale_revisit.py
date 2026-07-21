#!/usr/bin/env python3
"""Hermetic tests for Terraform >=8-day stale reminder and dual-channel delivery."""
import importlib.util
import json
import os
import sys
import tempfile
import unittest
from unittest import mock
from datetime import datetime
from pathlib import Path

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)

PROJECT = "528766"
TICKET = "90000001"


def epoch(text):
    return datetime.fromisoformat(text).replace(tzinfo=bot._SHANGHAI_TZ).timestamp()


def comment(cid, when, author, content, **extra):
    value = {
        "id": cid, "createdAt": when, "author": author, "content": content,
    }
    value.update(extra)
    return value


def item(created="2026-07-01 09:00:00", assignee="484483", iid=TICKET):
    return {
        "id": iid,
        "title": "alicloud_x stale",
        "pool": "tf_provider",
        "pool_project": PROJECT,
        "tag": ["jarvis-idle"],
        "status": "开发中",
        "assignee": assignee,
        "created": created,
        "modified": created,
        "terraform": True,
    }


class ClassifierTest(unittest.TestCase):
    def test_noise_does_not_reset_anchor(self):
        probes = [
            comment(1, "2026-07-08 09:00:00", "Terraform-研发数字人",
                    "已验证通过，PR #9"),
            comment(2, "2026-07-08 09:00:00", "过载",
                    "### 进度跟进 · 已 8 天无实质进展"),
            comment(3, "2026-07-08 09:00:00", "过载", "@新山(521957)"),
            comment(4, "2026-07-08 09:00:00", "过载", "收到，暂无更新"),
            comment(5, "2026-07-08 09:00:00", "WORKER_1783326253279",
                    "结论：已修复"),
            comment(6, "2026-07-08 09:00:00", "过载",
                    '[[PERSONA-HANDOFF:{"from":"pd","to":"rd"}]]'),
            comment(7, "2026-07-08 09:00:00", "过载",
                    "根因待定位，暂无结论"),
            comment(8, "2026-07-08 09:00:00", "过载",
                    "结论：尚未确认"),
            comment(9, "2026-07-08 09:00:00", "过载",
                    "当前阻塞上游 API，后续再看"),
            comment(10, "2026-07-08 09:00:00", "过载",
                    "阻塞待解决"),
            comment(11, "2026-07-08 09:00:00", "过载",
                    "有风险，需要关注"),
            comment(12, "2026-07-08 09:00:00", "过载",
                    "卡在接口，下一步未知"),
            comment(13, "2026-07-08 09:00:00", "过载",
                    "根因：需要进一步定位"),
            comment(14, "2026-07-08 09:00:00", "过载",
                    "Root cause: under investigation"),
            comment(15, "2026-07-08 09:00:00", "过载",
                    "建议调整 provider schema"),
            comment(16, "2026-07-08 09:00:00", "过载",
                    "计划新增字段 foo"),
            comment(17, "2026-07-08 09:00:00", "过载",
                    "准备调整 provider schema"),
            comment(18, "2026-07-08 09:00:00", "过载",
                    "拟删除旧逻辑"),
            comment(19, "2026-07-08 09:00:00", "过载",
                    "需要新增字段 foo"),
            comment(20, "2026-07-08 09:00:00", "过载",
                    "待调整 provider schema"),
            comment(21, "2026-07-08 09:00:00", "过载",
                    "将新增字段 foo"),
            comment(22, "2026-07-08 09:00:00", "过载",
                    "后续调整 provider schema"),
            comment(23, "2026-07-08 09:00:00", "过载",
                    "TODO: 新增字段 foo"),
            comment(24, "2026-07-08 09:00:00", "过载",
                    "方案是新增字段 foo"),
        ]
        for probe in probes:
            self.assertFalse(bot._is_substantial_progress(probe), probe)

    def test_questions_and_requests_are_not_progress(self):
        probes = [
            "请确认根因是什么？",
            "请确认 provider 字段是否支持？",
            "是否有阻塞？下一步是什么？",
            "麻烦看下 PR #123 是否已经提交？",
            "Could you confirm whether the provider schema supports foo?",
        ]
        for i, text in enumerate(probes):
            self.assertFalse(
                bot._is_substantial_progress(
                    comment(i + 20, "2026-07-08 09:00:00", "过载", text)),
                text)

    def test_system_and_persona_author_variants_are_not_progress(self):
        authors = [
            "Aone System", "系统机器人", "Automation Bot", "RD", "PD", "QA",
            "terraform-rd", "Terraform 研发数字人", "open-jarvis",
            "WORKER_1782379562571", "WORKER_1783582458263",
        ]
        for i, author in enumerate(authors):
            self.assertFalse(
                bot._is_substantial_progress(
                    comment(i + 40, "2026-07-08 09:00:00", author,
                            "根因已确认：provider schema 缺少字段，PR #123 已提交")),
                author)

    def test_noise_clause_does_not_hide_real_progress(self):
        probes = [
            "已认领，根因已确认：schema 缺少字段，PR #123 已提交。",
            "claim completed; root cause confirmed: provider schema misses foo; "
            "PR #124 submitted.",
        ]
        for i, text in enumerate(probes):
            self.assertTrue(
                bot._is_substantial_progress(
                    comment(i + 60, "2026-07-08 09:00:00", "过载", text)),
                text)

    def test_supported_progress_shapes_reset_anchor(self):
        probes = [
            "验证结果：回归测试通过。",
            "根因已确认：provider schema 缺少字段。",
            "PR https://github.com/aliyun/x/pull/9 已提交。",
            "当前阻塞在上游 API，下一步等待接口发布。",
            "计划 2026-07-20 完成修复。",
            "计划 7/20 完成修复。",
            "预计 07-20 完成修复。",
            "Blocked by the upstream API; next step: retry on July 20, 2026.",
            "已发布",
            "已合并",
            "已支持",
            "提交了 PR #123",
            "定位到 provider schema 缺少字段",
            "发现接口返回值为空",
            "确认原因是参数未透传",
            "代码改为重试三次",
            "新增字段 foo",
            "删除旧逻辑",
            "调整 provider schema",
            "日志显示调用了 DescribeX",
            "Logs show DescribeX returned an empty value.",
            "Blocker: upstream API; next step: provider retry",
            "阻塞上游 API，下一步由云产品修复",
            "依赖上游接口，下一步等发布",
            "代码已经改成重试三次",
            "字段增加了 foo",
            "补充了字段 foo",
            "修改了 provider schema",
            "修正了接口参数",
            "修了这个错误",
            "定位出 schema 缺字段",
            "查到原因是参数未透传",
            "确认是参数未透传",
            "结果表明接口返回空",
            "测试已跑通",
            "回归已过",
        ]
        for i, text in enumerate(probes):
            self.assertTrue(
                bot._is_substantial_progress(
                    comment(i + 1, "2026-07-08 09:00:00", "过载", text)),
                text)

    def test_owner_change_after_comment_starts_new_epoch(self):
        anchor = bot._stale_anchor(
            item(),
            [comment(9, "2026-07-05 09:00:00", "过载", "验证结果：测试通过")],
            [{"id": 10, "property": "指派给", "eventTime": "2026-07-07 09:00:00",
              "newValue": "新山"}])
        self.assertEqual(anchor["kind"], "owner")
        self.assertEqual(anchor["id"], "10")


class OwnerResolutionTest(unittest.TestCase):
    def test_worker_owner_uses_human_fallback(self):
        owner = bot._resolve_stale_owner(item(assignee="WORKER_1782379562571"))
        self.assertEqual(owner["staff_id"], "484483")
        self.assertEqual(owner["mention"], "@过载(484483)")
        self.assertEqual(owner["source_agent"], "WORKER_1782379562571")

    def test_human_owner_resolves_flower_and_id(self):
        owner = bot._resolve_stale_owner(item(assignee={"id": "521957", "name": "张旭"}))
        self.assertEqual(owner["staff_id"], "521957")
        self.assertEqual(owner["mention"], "@新山(521957)")


class SchedulerBoundaryTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.scheduler = bot._NudgeJob(
            handler=None, pool=None, stale_days=8, max_n=5)
        self.scheduler._ticket_timeline = lambda _item: ([], [])
        self.orig_aone = bot._aone_event_enqueue
        self.orig_dm = bot._dingtalk_event_enqueue
        self.aone_keys = []
        self.dm_keys = []
        bot._aone_event_enqueue = (
            lambda _t, _p, key, _text, **_kw:
            self.aone_keys.append(key) or True)
        bot._dingtalk_event_enqueue = (
            lambda _t, _p, key, _staff, _title, _text, **_kw:
            self.dm_keys.append(key) or True)

    def tearDown(self):
        bot._aone_event_enqueue = self.orig_aone
        bot._dingtalk_event_enqueue = self.orig_dm
        self.tmp.cleanup()

    def test_7d23h59_is_noop_exact_8d_dual_channel(self):
        anchor = epoch("2026-07-01T09:00:00")
        self.assertEqual(
            self.scheduler._remind_if_stale(
                item(), now=anchor + 8 * 86400 - 60),
            "not_due")
        self.assertEqual(self.aone_keys, [])
        self.assertEqual(self.dm_keys, [])
        self.assertEqual(
            self.scheduler._remind_if_stale(item(), now=anchor + 8 * 86400),
            "reminded")
        self.assertEqual(len(self.aone_keys), 1)
        self.assertEqual(self.aone_keys, self.dm_keys)

    def test_new_progress_and_owner_change_produce_new_event_keys(self):
        base = item()
        self.scheduler._ticket_timeline = lambda _item: (
            [comment(90, "2026-07-02 09:00:00", "过载",
                     "根因已确认：schema 缺少字段")], [])
        self.scheduler._remind_if_stale(
            base, now=epoch("2026-07-10T09:00:00"))
        first = self.aone_keys[-1]
        self.scheduler._ticket_timeline = lambda _item: (
            [comment(91, "2026-07-03 09:00:00", "过载",
                     "验证结果：修复已通过")], [])
        self.scheduler._remind_if_stale(
            base, now=epoch("2026-07-11T09:00:00"))
        second = self.aone_keys[-1]
        self.assertNotEqual(first, second)
        self.scheduler._ticket_timeline = lambda _item: (
            [], [{"id": 92, "property": "指派给",
                  "eventTime": "2026-07-04 09:00:00", "newValue": "新山"}])
        changed = item(assignee="521957")
        self.scheduler._remind_if_stale(
            changed, now=epoch("2026-07-12T09:00:00"))
        self.assertNotEqual(second, self.aone_keys[-1])

    def test_fixed_payload_contains_no_original_comment(self):
        self.scheduler._ticket_timeline = lambda _item: (
            [comment(1, "2026-07-01 09:00:00", "过载",
                     "根因已确认：SECRET-ORIGINAL-CONTENT")], [])
        captured = {}
        bot._aone_event_enqueue = (
            lambda _t, _p, _key, text, **_kw:
            captured.setdefault("aone", text) or True)
        bot._dingtalk_event_enqueue = (
            lambda _t, _p, _key, _staff, _title, text, **_kw:
            captured.setdefault("dm", text) or True)
        self.scheduler._remind_if_stale(
            item(), now=epoch("2026-07-09T09:00:00"))
        self.assertNotIn("SECRET-ORIGINAL-CONTENT", captured["aone"])
        self.assertNotIn("SECRET-ORIGINAL-CONTENT", captured["dm"])
        self.assertIn("@过载(484483)", captured["aone"])
        self.assertIn("[#%s]" % TICKET, captured["aone"])


class CandidateFairnessTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.scheduler = bot._NudgeJob(handler=None, pool=None, max_n=5)

    def tearDown(self):
        self.tmp.cleanup()

    def test_ordinary_terraform_idle_is_candidate_non_tf_stays_legacy(self):
        self.scheduler._pool_projects = lambda: [
            ("tf_provider", PROJECT), ("mcp_server", "2124589")]
        self.scheduler._query_pool = lambda key, project: [
            {
                "id": "1" if key == "tf_provider" else "2",
                "title": "ordinary ticket",
                "pool": key,
                "pool_project": project,
                "tag": ["jarvis-idle"],
                "status": "Open",
                "modified": "2026-07-01",
            }]
        ids = {row["id"] for row in self.scheduler._query()}
        self.assertIn("1", ids)
        self.assertNotIn("2", ids)

    def test_terraform_title_in_non_tf_pool_is_candidate(self):
        self.scheduler._pool_projects = lambda: [("mcp_server", "2124589")]
        self.scheduler._query_pool = lambda key, project: [{
            "id": "3",
            "title": "terraform-provider-alicloud 属性支持",
            "pool": key,
            "pool_project": project,
            "tag": ["jarvis-idle"],
            "status": "Open",
            "modified": "2026-07-01",
        }]
        rows = self.scheduler._query()
        self.assertEqual([row["id"] for row in rows], ["3"])
        self.assertTrue(rows[0]["terraform"])

    def test_customer_pr_merged_status_is_not_a_nudge_candidate(self):
        self.scheduler._pool_projects = lambda: [("tf_customer", "1086837")]
        self.scheduler._query_pool = lambda key, project: [{
            "id": "4", "title": "merged customer request", "pool": key,
            "pool_project": project,
            "tag": ["jarvis-idle"], "type": "3",
            "status": "已合入主线", "modified": "2026-07-01",
        }]
        self.assertEqual(self.scheduler._query(), [])

    def test_customer_idle_query_excludes_pr_merged_status(self):
        completed = type("Proc", (), {"returncode": 0, "stdout": "[]", "stderr": ""})()
        with mock.patch.object(bot.subprocess, "run", return_value=completed) as run:
            self.scheduler._query_pool("tf_customer", "1086837")
        cmd = run.call_args.args[0]
        self.assertEqual(cmd[cmd.index("--filter") + 1],
                         "tag=jarvis-idle AND NOT status=已合入主线")
        self.assertIn("type", cmd[cmd.index("--columns") + 1])

    def test_more_than_25_rows_rotate_without_starvation(self):
        rows = [
            {"id": str(i), "modified": "2026-07-%02d" % ((i % 28) + 1)}
            for i in range(30)]
        first = {it["id"] for it in self.scheduler._select_fair(rows, now=1000)}
        second = {it["id"] for it in self.scheduler._select_fair(rows, now=1001)}
        self.assertEqual(len(first), 5)
        self.assertEqual(len(second), 5)
        self.assertTrue(first.isdisjoint(second))

    def test_overdue_cold_row_beats_continuously_changed_hot_rows(self):
        now = 2_000_000.0
        rows = [{"id": "cold", "modified": "stable"}] + [
            {"id": "hot-%d" % i, "modified": "day-2"} for i in range(6)
        ]
        self.scheduler._write_index({
            "tickets": {
                "cold": {
                    "modified": "stable",
                    "next_check": now - 14 * 86400,
                    "last_inspected": now - 14 * 86400,
                },
                **{
                    "hot-%d" % i: {
                        "modified": "day-1",
                        "next_check": now,
                        "last_inspected": now - 86400,
                    }
                    for i in range(6)
                },
            },
        })
        chosen = self.scheduler._select_fair(rows, now=now)
        self.assertEqual(len(chosen), 5)
        self.assertIn("cold", {row["id"] for row in chosen})

    def test_query_pool_uses_page_size_1000_and_paginates(self):
        calls = []
        orig = bot.subprocess.run

        class Proc:
            returncode = 0
            stderr = ""

            def __init__(self, data):
                self.stdout = json.dumps(data)

        def fake_run(cmd, **_kwargs):
            page = int(cmd[cmd.index("--page") + 1])
            calls.append(page)
            count = 1000 if page == 1 else 30
            return Proc([
                {"identifier": "%d-%d" % (page, i), "subject": "x",
                 "tag": ["jarvis-idle"], "status": "Open"}
                for i in range(count)])

        bot.subprocess.run = fake_run
        try:
            rows = self.scheduler._query_pool("tf_provider", PROJECT)
        finally:
            bot.subprocess.run = orig
        self.assertEqual(len(rows), 1030)
        self.assertEqual(calls, [1, 2])


class Proc:
    def __init__(self, rc=0, out="", err=""):
        self.returncode = rc
        self.stdout = out
        self.stderr = err


class DualChannelLedgerTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.orig_aone_path = bot.AONE_EVENT_PATH
        self.orig_dm_path = bot.DINGTALK_EVENT_PATH
        self.orig_tf = bot._is_terraform_project
        self.orig_run = bot.subprocess.run
        bot.AONE_EVENT_PATH = Path(self.tmp.name) / "aone.json"
        bot.DINGTALK_EVENT_PATH = Path(self.tmp.name) / "dm.json"
        bot._is_terraform_project = lambda project: str(project) == PROJECT
        bot._aone_event_inflight.clear()
        bot._dingtalk_event_inflight.clear()
        self.aone_fail = False
        self.aone_uncertain = False
        self.dm_fail = False
        self.dm_raise = False
        self.aone_create = 0
        self.dm_send = 0
        self.remote = []

        def fake_run(cmd, **_kwargs):
            argv = list(cmd)
            if "notify-dingtalk.sh" in argv[0]:
                self.dm_send += 1
                if self.dm_raise:
                    raise TimeoutError("transport timeout")
                if self.dm_fail:
                    return Proc(0, json.dumps({
                        "status": "failed", "reason": "network",
                        "receipt": argv[argv.index("--out-track-id") + 1]}))
                return Proc(0, json.dumps({
                    "status": "sent", "reason": "delivered",
                    "receipt": argv[argv.index("--out-track-id") + 1]}))
            if "list" in argv:
                return Proc(0, json.dumps([
                    {"id": i + 1, "content": text}
                    for i, text in enumerate(self.remote)]))
            if "create" in argv:
                self.aone_create += 1
                if self.aone_fail:
                    return Proc(1, "", "aone down")
                body = argv[argv.index("-m") + 1]
                if not self.aone_uncertain:
                    self.remote.append(body)
                    return Proc(0, '{"id":99}')
                return Proc(0, "{}")
            return Proc(1, "", "unexpected")

        bot.subprocess.run = fake_run

    def tearDown(self):
        bot.AONE_EVENT_PATH = self.orig_aone_path
        bot.DINGTALK_EVENT_PATH = self.orig_dm_path
        bot._is_terraform_project = self.orig_tf
        bot.subprocess.run = self.orig_run
        bot._aone_event_inflight.clear()
        bot._dingtalk_event_inflight.clear()
        self.tmp.cleanup()

    def event(self):
        return "revisit:stale:%s:created:1:1:484483" % TICKET

    def test_aone_success_dm_failure_only_retries_dm(self):
        self.dm_fail = True
        self.assertTrue(bot._aone_event_enqueue(TICKET, PROJECT, self.event(), "aone"))
        self.assertTrue(bot._dingtalk_event_enqueue(
            TICKET, PROJECT, self.event(), "484483", "title", "dm"))
        self.assertEqual((self.aone_create, self.dm_send), (1, 1))
        ledger = bot._dingtalk_event_load()
        lid = bot._aone_event_ledger_id(TICKET, self.event())
        ledger["pending"][lid]["not_before"] = 0
        bot._dingtalk_event_write(ledger)
        self.dm_fail = False
        bot._aone_event_enqueue(TICKET, PROJECT, self.event(), "aone")
        bot._dingtalk_event_flush()
        self.assertEqual((self.aone_create, self.dm_send), (1, 2))

    def test_dm_success_aone_failure_only_retries_aone(self):
        self.aone_fail = True
        self.assertTrue(bot._aone_event_enqueue(TICKET, PROJECT, self.event(), "aone"))
        self.assertTrue(bot._dingtalk_event_enqueue(
            TICKET, PROJECT, self.event(), "484483", "title", "dm"))
        self.assertEqual((self.aone_create, self.dm_send), (1, 1))
        self.aone_fail = False
        bot._aone_event_flush()
        bot._dingtalk_event_enqueue(
            TICKET, PROJECT, self.event(), "484483", "title", "dm")
        self.assertEqual((self.aone_create, self.dm_send), (2, 1))

    def test_aone_uncertain_never_reposts_and_dm_does_not_repeat(self):
        self.aone_uncertain = True
        self.assertTrue(bot._aone_event_enqueue(TICKET, PROJECT, self.event(), "aone"))
        self.assertTrue(bot._dingtalk_event_enqueue(
            TICKET, PROJECT, self.event(), "484483", "title", "dm"))
        lid = bot._aone_event_ledger_id(TICKET, self.event())
        ledger = bot._aone_event_load()
        ledger["pending"][lid]["not_before"] = 0
        bot._aone_event_write(ledger)
        bot._aone_event_flush()
        bot._dingtalk_event_enqueue(
            TICKET, PROJECT, self.event(), "484483", "title", "dm")
        self.assertEqual((self.aone_create, self.dm_send), (1, 1))
        self.assertEqual(bot._aone_event_load()["pending"][lid]["state"], "post_uncertain")

    def test_dm_transport_uncertain_keeps_stable_receipt_for_retry(self):
        self.dm_raise = True
        self.assertTrue(bot._dingtalk_event_enqueue(
            TICKET, PROJECT, self.event(), "484483", "title", "dm"))
        lid = bot._aone_event_ledger_id(TICKET, self.event())
        ledger = bot._dingtalk_event_load()
        pending = ledger["pending"][lid]
        receipt = pending["receipt"]
        self.assertEqual(pending["state"], "post_uncertain")
        pending["not_before"] = 0
        bot._dingtalk_event_write(ledger)
        self.dm_raise = False
        bot._dingtalk_event_flush()
        self.assertEqual(bot._dingtalk_event_load()["posted"][lid]["receipt"], receipt)
        self.assertEqual(self.dm_send, 2)


if __name__ == "__main__":
    unittest.main(verbosity=2)
