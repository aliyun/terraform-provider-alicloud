#!/usr/bin/env python3
"""Regression: 明确关单请求由内部角色核验，最终 TerraformRD 聚合回复一次请人工关单。

旧公开接力只保留入站触发；新 prompt 不生成阶段评论/哨兵/钉钉通知，normal/close/escalate
都最多一次最终 RD wrap。

Standalone: `python3 bridge/test_persona_close_handoff.py`. No a1/Aone/DingTalk calls.
"""
import importlib.util
import os
import sys
import tempfile
import unittest
from pathlib import Path

os.environ.setdefault("JARVIS_PERSONA_NICKS", "")

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)

P = bot.PersonaScheduler


def _sched():
    tmp = tempfile.NamedTemporaryFile(suffix=".json", delete=False, mode="w")
    tmp.write('{"tickets": {}}')
    tmp.close()
    return bot.PersonaScheduler(handler=None, pool=None, enabled=False, ledger_path=tmp.name)


def _state():
    return {"last_seen": 0, "processed": set(), "dispatch_count": 0, "escalated": False}


class DetectCloseRequestTest(unittest.TestCase):
    def test_positive(self):
        for s in ("@Terraform-PD数字人 关闭这个aone", "请关单", "麻烦关掉",
                  "销单吧", "please close this ticket", "结单"):
            self.assertTrue(P._detect_close_request(s), s)

    def test_negative(self):
        for s in ("这个属性支持吗？", "@terraform-pd 看下根因", "什么时候能接入", ""):
            self.assertFalse(P._detect_close_request(s), s)


class DecideInjectsCloseContextTest(unittest.TestCase):
    def _decide_one(self, content):
        item = {"id": "84240297", "tag": ""}
        comments = [{"id": 999, "author": "原根", "content": content, "createdAt": ""}]
        return _sched()._decide_persona(item, comments, state=_state())

    def test_human_close_request_flags_dispatch(self):
        ds = self._decide_one("@terraform-rd 关闭这个aone")
        d = next(x for x in ds if x["action"] == "dispatch")
        self.assertEqual(d["reason"], "human_mention")
        self.assertEqual(d["internal_role"], "terraform-pd")
        self.assertEqual(d["public_identity"], "terraform-rd")
        self.assertTrue(d["handoff"]["close_request"])
        self.assertEqual(d["handoff"]["requester"], "原根")
        self.assertFalse(d["handoff"]["requester_is_digital"])

    def test_human_mention_without_close_is_not_flagged(self):
        ds = self._decide_one("@terraform-rd 帮忙看下根因")
        d = next(x for x in ds if x["action"] == "dispatch")
        self.assertEqual(d["reason"], "human_mention")
        self.assertFalse(d["handoff"]["close_request"])


class DetectJarvisMentionTest(unittest.TestCase):
    def test_positive(self):
        for s in ("@jarvis 关闭", "@open-jarvis 关单", "@open_jarvis 结单",
                  "@Name(WORKER_1782379562571) 关单"):
            self.assertTrue(P._detect_jarvis_mention(s), s)

    def test_negative(self):
        for s in ("@terraform-pd 关闭", "这个jarvisxx", "关单吧", ""):
            self.assertFalse(P._detect_jarvis_mention(s), s)


class JarvisCloseMentionTest(unittest.TestCase):
    """@jarvis + 关单请求 → 复用关单 handoff，由 terraform-pd 代为核验+催人工关单。"""
    def _decide_one(self, content):
        item = {"id": "84240297", "tag": ""}
        comments = [{"id": 999, "author": "原根", "content": content, "createdAt": ""}]
        return _sched()._decide_persona(item, comments, state=_state())

    def test_jarvis_close_dispatches_via_pd(self):
        ds = self._decide_one("@jarvis 关闭这个aone")
        d = next(x for x in ds if x["action"] == "dispatch")
        self.assertEqual(d["internal_role"], "terraform-pd")   # jarvis 无子代理 → 交 pd 代办
        self.assertTrue(d["handoff"]["close_request"])
        self.assertEqual(d["handoff"]["requester"], "原根")

    def test_jarvis_mention_without_close_is_ignored(self):
        ds = self._decide_one("@jarvis 看下这个")
        self.assertFalse([x for x in ds if x["action"] == "dispatch"])  # 一般 @jarvis 不触发

    def test_tracker_watch_includes_jarvis_worker(self):
        watch = bot.PERSONA_WORKER_IDS | {bot.JARVIS_ORCH_WORKER}
        self.assertIn("WORKER_1782379562571", watch)
        self.assertEqual(bot.PERSONA_WORKER_IDS, {bot.PERSONA_PUBLIC_WORKER})
        self.assertNotIn("WORKER_1783582374386", watch)
        self.assertNotIn("WORKER_1783582593461", watch)


class PromptCloseHandoffTest(unittest.TestCase):
    def test_human_requester_pings_requester_not_finish(self):
        p = bot._persona_prompt("84240297", "terraform-pd", "respond", "", 1, "snip",
                                project="1086837", close_request=True,
                                requester="原根", requester_is_digital=False)
        self.assertIn("关单请求", p)
        self.assertIn("原根", p)
        self.assertIn("release", p)
        self.assertIn("不 finish", p)          # 关单仍人工
        self.assertEqual(p.count("wrap.sh done"), 1)
        self.assertNotIn("notify-dingtalk", p)
        self.assertNotIn("wrap.sh sync", p)
        self.assertNotIn("PERSONA-HANDOFF", p)

    def test_digital_requester_escalates_to_humans(self):
        p = bot._persona_prompt("84240297", "terraform-pd", "respond", "", 1, "snip",
                                project="1086837", close_request=True,
                                requester="open-jarvis", requester_is_digital=True)
        self.assertIn("辰羿(320687)", p)
        self.assertIn("过载(484483)", p)
        self.assertIn("数字人不能授权关单", p)
        self.assertIn("不 finish", p)
        self.assertEqual(p.count("wrap.sh done"), 1)

    def test_non_close_prompt_unchanged(self):
        p = bot._persona_prompt("111", "terraform-pd", "respond", "", 1, "snip",
                                close_request=False)
        self.assertIn("迁移补位", p)
        self.assertEqual(p.count("wrap.sh done"), 1)
        self.assertNotIn("关单请求", p)

    def test_prompt_uses_single_public_identity(self):
        p = bot._persona_prompt("111", "terraform-qa", "acc_verify", "", 1, "snip")
        self.assertIn("internal_role=terraform-qa", p)
        self.assertIn("public_identity=terraform-rd", p)
        self.assertNotIn("as terraform-qa", p)
        self.assertNotIn("identity_fallback", p)
        self.assertNotIn("[QA验收]", p)

    def test_normal_close_escalate_each_have_at_most_one_rd_reply(self):
        prompts = (
            bot._persona_prompt("111", "terraform-pd", "respond", "", 1, "snip"),
            bot._persona_prompt("111", "terraform-pd", "respond", "", 1, "snip",
                                close_request=True, requester="原根"),
            bot._persona_prompt("111", "terraform-rd", "dev", "", 7, "snip",
                                escalated=True),
        )
        for p in prompts:
            self.assertEqual(p.count("wrap.sh done"), 1)
            for forbidden in ("[PD分诊]", "[RD开发]", "[QA验收]",
                              "PERSONA-HANDOFF", "wrap.sh sync",
                              "comment create", "notify-dingtalk"):
                self.assertNotIn(forbidden, p)


if __name__ == "__main__":
    unittest.main(verbosity=2)
