#!/usr/bin/env python3
"""Regression: persona (terraform-pd/rd/qa) comments must NOT count as human
intervention in ScanScheduler._is_human_comment.

Bug: personas are jarvis-driven collaboration instances. Treating their stage
comments as "human touched" made every persona comment on a jarvis-idle ticket
trigger a force re-dispatch → redundant headless instances → duplicate bookend
comment churn. Fix: exclude persona authors via _author_role.

Standalone: `python3 bridge/test_scan_persona_churn.py`. No a1/Aone/DingTalk calls.
"""
import importlib.util
import os
import sys
import unittest
from pathlib import Path

# Persona display names are only resolvable via the nick map (rd/qa 中文名 + pd CJK
# 边界),so mirror the bridge env before loading the module.
os.environ["JARVIS_PERSONA_NICKS"] = (
    "terraform-pd=Terraform-PD数字人,"
    "terraform-rd=Terraform-研发数字人,"
    "terraform-qa=Terraform-质量保障数字人"
)

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)

_ishuman = bot.ScanScheduler._is_human_comment


class PersonaNotHumanTest(unittest.TestCase):
    def test_personas_are_not_human(self):
        for name in ("Terraform-PD数字人", "Terraform-研发数字人",
                     "Terraform-质量保障数字人"):
            self.assertFalse(_ishuman(name),
                             "%s 应被判为非人工(persona)" % name)

    def test_jarvis_and_bots_not_human(self):
        for name in ("open-jarvis", "WORKER_1782379562571", "Kelude",
                     "云知道平台公共账号"):
            self.assertFalse(_ishuman(name))

    def test_jarvis_claim_content_not_human(self):
        # 认领系统评论不算人工
        self.assertFalse(_ishuman("过载", "jarvis-claim AgenticTools-Macmini"))

    def test_real_humans_still_human(self):
        # 真人（未登记 persona/jarvis）仍应触发人工介入门
        for name in ("过载", "辰羿", "某位客户"):
            self.assertTrue(_ishuman(name), "%s 应判为人工" % name)

    def test_empty_author_not_human(self):
        self.assertFalse(_ishuman(""))
        self.assertFalse(_ishuman(None))


if __name__ == "__main__":
    unittest.main(verbosity=2)
