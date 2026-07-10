#!/usr/bin/env python3
"""Regression tests for the Tata->Jarvis delegation gate whitelist (api_tool_staff).

Standalone: `python3 bridge/test_delegation_gate.py`. The bot module imports cleanly
without the dingtalk SDK (import is try/excepted at top), so we load it by path.
"""
import importlib.util
import os
import sys
import unittest
from pathlib import Path

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)

MASTER = "320687"      # 辰羿 (default JARVIS_MASTER_STAFF)
CONTACT = "270513"     # 秋雯 / 陈闻菲 — a real entry in config/contacts.json
FAKE = "111111"        # deliberately NOT in contacts.json
ADDED = "999888"       # via JARVIS_API_TOOL_STAFF


class DelegationGateTest(unittest.TestCase):
    def setUp(self):
        # snapshot mutable module/env state so tests never leak into each other
        self._orig_repo_root = bot.REPO_ROOT
        self._orig_master = os.environ.get("JARVIS_MASTER_STAFF")
        self._orig_extra = os.environ.get("JARVIS_API_TOOL_STAFF")
        os.environ.pop("JARVIS_API_TOOL_STAFF", None)

    def tearDown(self):
        bot.REPO_ROOT = self._orig_repo_root
        for k, v in (("JARVIS_MASTER_STAFF", self._orig_master),
                     ("JARVIS_API_TOOL_STAFF", self._orig_extra)):
            if v is None:
                os.environ.pop(k, None)
            else:
                os.environ[k] = v

    def test_includes_master_and_real_contact(self):
        ids = bot.api_tool_staff()
        self.assertIn(MASTER, ids, "master (320687) must be whitelisted")
        self.assertIn(CONTACT, ids, "real contact 270513 must be whitelisted from contacts.json")

    def test_master_present_even_when_contacts_unreadable(self):
        bot.REPO_ROOT = Path("/nonexistent-xyz")
        try:
            ids = bot.api_tool_staff()
        finally:
            bot.REPO_ROOT = self._orig_repo_root
        self.assertIn(bot.master_staff(), ids,
                      "master must survive missing/broken contacts.json")

    def test_env_additive(self):
        os.environ["JARVIS_API_TOOL_STAFF"] = ADDED
        try:
            ids = bot.api_tool_staff()
        finally:
            os.environ.pop("JARVIS_API_TOOL_STAFF", None)
        self.assertIn(ADDED, ids, "JARVIS_API_TOOL_STAFF entries must be added")
        # additive, not replacing: master + contacts still present
        self.assertIn(MASTER, ids)
        self.assertIn(CONTACT, ids)

    def test_gate_decision(self):
        ids = bot.api_tool_staff()
        # this is exactly the runtime gate: `staff in api_tool_staff()`
        self.assertTrue(CONTACT in ids, "whitelisted staff must pass the gate")
        self.assertFalse(FAKE in ids, "non-contact staff must NOT pass the gate")


if __name__ == "__main__":
    unittest.main(verbosity=2)
