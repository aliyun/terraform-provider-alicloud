#!/usr/bin/env python3
"""Regression: PersonaScheduler._query_candidates must also surface tickets where
a persona 数字人 sits on 抄送 (workitem.tracker), not only jarvis-idle-tagged ones.

Bug: 一条工单指派给别人、无 jarvis 标签，客户/同事在评论区 @ 了某数字人（Aone 自动把该
数字人加入抄送）。旧 _query_candidates 只查 tag=jarvis-idle，这类单永远不进候选集，评论区
@数字人 就无人响应（实例：工单 84240297，原根 @Terraform-PD数字人 关单无人理）。
Fix: union 一条 workitem.tracker=<persona workers> 查询；派发/判定链路原样复用。

Standalone: `python3 bridge/test_persona_tracker_scan.py`. No a1/Aone/DingTalk calls
(monkeypatch _query_pool_filter 返回 canned 结果)。
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

PERSONA_WID = sorted(bot.PERSONA_WORKER_IDS)[0]  # 任一 persona worker id
TERMINAL = next(iter(bot.TERMINAL_STATUSES))     # 任一终态值


def _make_scheduler():
    """构造一个不启线程、不碰真 ledger 的 PersonaScheduler。"""
    tmp = tempfile.NamedTemporaryFile(suffix=".json", delete=False, mode="w")
    tmp.write('{"tickets": {}}')
    tmp.close()
    return bot.PersonaScheduler(handler=None, pool=None, enabled=False,
                                ledger_path=tmp.name)


def _stub_query(idle=None, tracker=None, only_project="1086837"):
    """返回一个替换 _query_pool_filter 的 staticmethod：按 filter_expr 分派 canned 结果，
    仅对 only_project 生效（其它池回空，避免跨 4 池叠加干扰断言）。"""
    idle = idle or []
    tracker = tracker or []

    def fake(key, project, filter_expr):
        if str(project) != only_project:
            return []
        if filter_expr.startswith("tag=jarvis-idle"):
            return [dict(x) for x in idle]
        if filter_expr.startswith("workitem.tracker="):
            return [dict(x) for x in tracker]
        return []
    return staticmethod(fake)


class TrackerScanTest(unittest.TestCase):
    def setUp(self):
        self._orig = bot.PersonaScheduler._query_pool_filter
        os.environ.pop("JARVIS_PERSONA_TRACKER_SCAN", None)

    def tearDown(self):
        bot.PersonaScheduler._query_pool_filter = self._orig
        os.environ.pop("JARVIS_PERSONA_TRACKER_SCAN", None)

    def _ids(self, sched):
        return {str(it.get("id")) for it in sched._query_candidates()}

    def test_tracker_ticket_without_jarvis_tag_included(self):
        """无 jarvis 标签、只在抄送里挂 persona 的单，应进候选集。"""
        bot.PersonaScheduler._query_pool_filter = _stub_query(
            idle=[{"id": "111", "title": "idle one", "tag": "jarvis-idle", "status": "处理中"}],
            tracker=[{"id": "84240297", "title": "被@的测试单", "tag": "", "status": "待处理"}],
        )
        ids = self._ids(_make_scheduler())
        self.assertIn("84240297", ids)   # tracker 命中，新增覆盖
        self.assertIn("111", ids)        # jarvis-idle 既有行为不丢

    def test_overlap_deduped(self):
        """同一单既在 jarvis-idle 又在 tracker 结果里 → 去重只留一条。"""
        row = {"id": "222", "title": "dual", "tag": "jarvis-idle", "status": "Open"}
        bot.PersonaScheduler._query_pool_filter = _stub_query(idle=[row], tracker=[dict(row)])
        cands = _make_scheduler()._query_candidates()
        self.assertEqual([c for c in cands if str(c.get("id")) == "222"].__len__(), 1)

    def test_terminal_status_filtered(self):
        """tracker 命中的终态单仍被 S9 过滤掉。"""
        bot.PersonaScheduler._query_pool_filter = _stub_query(
            tracker=[{"id": "333", "title": "done", "tag": "", "status": TERMINAL}],
        )
        self.assertNotIn("333", self._ids(_make_scheduler()))

    def test_kill_switch_disables_tracker_scan(self):
        """JARVIS_PERSONA_TRACKER_SCAN=0 → 退化为只扫 jarvis-idle，tracker 单不进候选。"""
        os.environ["JARVIS_PERSONA_TRACKER_SCAN"] = "0"
        bot.PersonaScheduler._query_pool_filter = _stub_query(
            idle=[{"id": "111", "title": "idle one", "tag": "jarvis-idle", "status": "处理中"}],
            tracker=[{"id": "84240297", "title": "被@的测试单", "tag": "", "status": "待处理"}],
        )
        ids = self._ids(_make_scheduler())
        self.assertIn("111", ids)
        self.assertNotIn("84240297", ids)

    def test_tracker_filter_uses_all_persona_workers(self):
        """tracker 过滤表达式覆盖全部 3 个 persona worker（多值 OR）。"""
        captured = {}

        def fake(key, project, filter_expr):
            if filter_expr.startswith("workitem.tracker="):
                captured["expr"] = filter_expr
            return []
        bot.PersonaScheduler._query_pool_filter = staticmethod(fake)
        _make_scheduler()._query_candidates()
        self.assertIn("expr", captured)
        for wid in bot.PERSONA_WORKER_IDS:
            self.assertIn(wid, captured["expr"])


if __name__ == "__main__":
    unittest.main(verbosity=2)
