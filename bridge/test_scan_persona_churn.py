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


class HumanOperatorWhitelistTest(unittest.TestCase):
    """_load_human_operators 与评论路径同一不变量:jarvis 自身 + 三数字人不入
    「人工介入门」白名单,即使日后被补录进 contacts.json;外部 agent(镇元agent)
    与真人保留。"""

    def _load_with_contacts(self, contacts):
        import json as _json
        import tempfile
        with tempfile.TemporaryDirectory() as td:
            cfg_dir = Path(td) / "config"
            cfg_dir.mkdir()
            (cfg_dir / "contacts.json").write_text(
                _json.dumps({"contacts": contacts}, ensure_ascii=False))
            orig = bot.REPO_ROOT
            bot.REPO_ROOT = Path(td)
            try:
                return bot.ScanScheduler._load_human_operators(None)
            finally:
                bot.REPO_ROOT = orig

    def test_personas_excluded_even_if_registered(self):
        ops = self._load_with_contacts([
            {"name": "open-jarvis", "flower": None, "id": "WORKER_1782379562571"},
            {"name": "terraform-pd", "flower": "Terraform-PD数字人",
             "id": "WORKER_1783582374386"},
            {"name": "terraform-rd", "flower": "Terraform-研发数字人",
             "id": "WORKER_1783582458263"},
            {"name": "terraform-qa", "flower": "Terraform-质量保障数字人",
             "id": "WORKER_1783582593461"},
            {"name": "镇元agent", "flower": None, "id": "WORKER_1783326253279"},
            {"name": "过载", "flower": "过载", "id": "484483"},
        ])
        for banned in ("open-jarvis", "WORKER_1782379562571",
                       "terraform-pd", "WORKER_1783582374386",
                       "Terraform-PD数字人", "terraform-rd",
                       "WORKER_1783582458263", "terraform-qa",
                       "WORKER_1783582593461"):
            self.assertNotIn(banned, ops, "%s 不应入人工门白名单" % banned)
        for kept in ("镇元agent", "WORKER_1783326253279", "过载", "484483"):
            self.assertIn(kept, ops, "%s 应保留在人工门白名单" % kept)


if __name__ == "__main__":
    unittest.main(verbosity=2)
