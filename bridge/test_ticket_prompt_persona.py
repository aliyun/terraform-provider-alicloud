#!/usr/bin/env python3
"""Tests for terraform-line detection + persona-chain dispatch prompt.

Standalone: `python3 bridge/test_ticket_prompt_persona.py`. The bot module imports
cleanly without the dingtalk SDK (import is try/excepted at top), so we load it by
path. These are pure-function tests — no subprocess / a1 / Aone calls.

Covers the改造 that routes terraform tickets through terraform-pd→rd→qa instead of
jarvis-inline triage (see loops/persona-collab.md, _ticket_prompt).
"""
import importlib.util
import sys
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)


class PoolLineTest(unittest.TestCase):
    def test_terraform_pools_are_terraform_line(self):
        self.assertEqual(bot._pool_line("tf_customer"), "terraform_provider")
        self.assertEqual(bot._pool_line("tf_provider"), "terraform_provider")

    def test_non_terraform_pools(self):
        self.assertEqual(bot._pool_line("mcp_server"), "agentic_tools")

    def test_unknown_pool_returns_empty(self):
        self.assertEqual(bot._pool_line("bogus"), "")
        self.assertEqual(bot._pool_line(""), "")
        self.assertEqual(bot._pool_line(None), "")


class IsTerraformTicketTest(unittest.TestCase):
    def test_by_pool_line(self):
        self.assertTrue(bot._is_terraform_ticket("tf_customer", "任意标题"))
        self.assertTrue(bot._is_terraform_ticket("tf_provider", "任意标题"))

    def test_non_terraform_pool_plain_title(self):
        self.assertFalse(bot._is_terraform_ticket("mcp_server", "agent 门户某 bug"))
        self.assertFalse(bot._is_terraform_ticket("api_toolkit", "工具需求"))

    def test_wrong_pool_but_title_keyword_fallback(self):
        # 落错池但标题命中 alicloud_ → 仍判定 terraform
        self.assertTrue(
            bot._is_terraform_ticket("mcp_server", "接入 alicloud_oss_bucket_inventory"))
        self.assertTrue(
            bot._is_terraform_ticket("api_toolkit", "terraform-provider 报错"))

    def test_none_title_safe(self):
        self.assertFalse(bot._is_terraform_ticket("mcp_server", None))


class TicketPromptTest(unittest.TestCase):
    def test_terraform_prompt_is_persona_chain(self):
        p = bot._ticket_prompt("84215653", "[TF] alicloud_x", "tf_customer", "1086837")
        for tok in ("terraform-pd", "terraform-rd", "terraform-qa", "哨兵",
                    "identity_fallback"):
            self.assertIn(tok, p, "terraform prompt 缺 %r" % tok)
        # 编排语义:不自己写码/验收
        self.assertIn("只做编排", p)

    def test_terraform_prompt_has_pr_ci_gate(self):
        # RD 交 QA 前须 PR CI 全绿(红 CI 的 PR 不该丢给 QA 空跑)
        p = bot._ticket_prompt("84215653", "[TF] alicloud_x", "tf_provider", "528766")
        self.assertIn("PR CI", p)
        self.assertIn("gh pr checks", p)

    def test_terraform_bookend_is_silent(self):
        # terraform 线:编排层收尾只记 runs + 打标签,不发评论
        # (prompt 里"wrap.sh done"只出现在"严禁…"禁令里,故断言意图而非裸子串缺失)
        p = bot._ticket_prompt("84215653", "[TF] alicloud_x", "tf_provider", "528766")
        self.assertIn("log.sh run_done", p)
        self.assertIn("claim.sh finish", p)
        self.assertIn("claim.sh release", p)
        self.assertIn("编排层不发评论", p)
        self.assertIn("不再落任何总结评论", p)
        # wrap.sh done/sync 只应作为禁令出现,不应作为收尾指令
        self.assertIn("严禁 wrap.sh done/sync", p)

    def test_non_terraform_prompt_stays_inline(self):
        p = bot._ticket_prompt("999", "agent 门户 bug", "mcp_server", "2124589")
        self.assertNotIn("terraform-pd", p)
        self.assertIn("aone-triage", p)
        # 非 terraform 单:编排层就是执行者,仍走 wrap.sh done 发评论收尾
        self.assertIn("wrap.sh done", p)

    def test_both_branches_preserve_dedup_claim_suspend(self):
        # 回归护栏:去重/认领/挂起语义两分支都不能丢
        for pool, proj, title in (("tf_customer", "1086837", "[TF] x"),
                                   ("mcp_server", "2124589", "agent bug")):
            p = bot._ticket_prompt("42", title, pool, proj)
            for tok in ("log.sh seen", "claim.sh claim", "[[SUSPEND:"):
                self.assertIn(tok, p, "%s 分支缺 %r" % (pool, tok))


if __name__ == "__main__":
    unittest.main(verbosity=2)
