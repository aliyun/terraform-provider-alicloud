#!/usr/bin/env python3
"""Regression tests for the Tata->Jarvis handoff-to-ticket path (#84141740).

Standalone: `python3 bridge/test_handoff_ticket.py`. The bot module imports cleanly
without the dingtalk SDK (import is try/excepted at top), so we load it by path.
subprocess.run is mocked — no real a1 / Aone calls.
"""
import importlib.util
import os
import sys
import unittest
from pathlib import Path
from unittest import mock

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)


class _FakeProc:
    def __init__(self, returncode=0, stdout="", stderr=""):
        self.returncode = returncode
        self.stdout = stdout
        self.stderr = stderr


def _make_handler():
    """Build a JarvisHandler without running __init__, wired with a recording
    _quick_card so tests can inspect the receipt text without touching DingTalk."""
    h = bot.JarvisHandler.__new__(bot.JarvisHandler)
    h.no_dingtalk = True
    h.cards = []

    def _rec(target, text, target_type="user"):
        h.cards.append((target, text, target_type))

    h._quick_card = _rec
    return h


class HandoffModeTest(unittest.TestCase):
    def setUp(self):
        self._orig = os.environ.get("JARVIS_HANDOFF_MODE")
        os.environ.pop("JARVIS_HANDOFF_MODE", None)

    def tearDown(self):
        if self._orig is None:
            os.environ.pop("JARVIS_HANDOFF_MODE", None)
        else:
            os.environ["JARVIS_HANDOFF_MODE"] = self._orig

    def test_handoff_mode_default_ticket(self):
        self.assertEqual(bot.handoff_mode(), "ticket")
        os.environ["JARVIS_HANDOFF_MODE"] = "exec"
        self.assertEqual(bot.handoff_mode(), "exec")
        os.environ["JARVIS_HANDOFF_MODE"] = "  EXEC  "
        self.assertEqual(bot.handoff_mode(), "exec")


class HandoffPoolTest(unittest.TestCase):
    def setUp(self):
        self._orig = os.environ.get("JARVIS_HANDOFF_POOL")
        os.environ.pop("JARVIS_HANDOFF_POOL", None)

    def tearDown(self):
        if self._orig is None:
            os.environ.pop("JARVIS_HANDOFF_POOL", None)
        else:
            os.environ["JARVIS_HANDOFF_POOL"] = self._orig

    def test_handoff_pool_default(self):
        key, proj, cfs = bot.handoff_pool()
        self.assertEqual(key, "api_toolkit")
        self.assertEqual(proj, "2100304")
        self.assertEqual(cfs, "107239=906688")

    def test_handoff_pool_override_fallback(self):
        os.environ["JARVIS_HANDOFF_POOL"] = "nonexistent_xyz"
        key, proj, _cfs = bot.handoff_pool()
        # unknown key falls back to api_toolkit
        self.assertEqual(key, "api_toolkit")
        self.assertEqual(proj, "2100304")


class HandoffToTicketTest(unittest.TestCase):
    def setUp(self):
        for k in ("JARVIS_HANDOFF_POOL", "JARVIS_HANDOFF_MODE"):
            os.environ.pop(k, None)

    def test_handoff_creates_ticket(self):
        h = _make_handler()
        proc = _FakeProc(returncode=0, stdout="8888 sometitle 待处理 jarvis\n")
        with mock.patch.object(bot.subprocess, "run", return_value=proc) as m:
            new_id = h._handoff_to_ticket("测试委派任务内容", "320687", "target", "user")
        self.assertEqual(new_id, "8888")
        self.assertTrue(m.called, "subprocess.run must be invoked to create the ticket")
        # receipt card fired, containing #id and the aone url
        joined = " ".join(c[1] for c in h.cards)
        self.assertIn("#8888", joined)
        self.assertIn("project.aone.alibaba-inc.com", joined)
        self.assertIn("8888", joined.split("req/")[-1])

    def test_handoff_create_failure(self):
        h = _make_handler()
        proc = _FakeProc(returncode=1, stdout="", stderr="boom")
        with mock.patch.object(bot.subprocess, "run", return_value=proc):
            new_id = h._handoff_to_ticket("任务", "320687", "target", "user")
        self.assertIsNone(new_id)
        joined = " ".join(c[1] for c in h.cards)
        self.assertTrue("失败" in joined or "异常" in joined,
                        "failure path must surface a 失败/异常 receipt")

    def test_handoff_exception_path(self):
        h = _make_handler()
        with mock.patch.object(bot.subprocess, "run", side_effect=RuntimeError("x")):
            new_id = h._handoff_to_ticket("任务", "320687", "target", "user")
        self.assertIsNone(new_id)
        joined = " ".join(c[1] for c in h.cards)
        self.assertTrue("失败" in joined or "异常" in joined)


if __name__ == "__main__":
    unittest.main(verbosity=2)
