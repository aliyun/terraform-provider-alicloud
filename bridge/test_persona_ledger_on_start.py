#!/usr/bin/env python3
"""Regression: PersonaScheduler must commit the ledger (last_seen/processed/count)
only when the dispatched worker ACTUALLY STARTS, never at submit-accept time.

Bug (工单 84297352, 原根 @Terraform-研发数字人 两次都无人理):
_apply_decisions 旧逻辑在 self.pool.submit(...) 返回 (accepted=True) 后就同步把评论写进
ledger.processed。但 submit 只表示"入队接受"——EphemeralExecutor 槽满时 future 仅排队。此刻
bridge 换 token 重启，terminate_all 的 shutdown(cancel_futures=True) 丢弃未启动的排队
future（worker 从未跑、评论从未回），而 ledger 已标 processed → PersonaScheduler 永久
skip(reason=processed) 不再重派，ScanScheduler 冷启动又把它当存量积压 → 双重盲区烂尾。

Fix: 落盘挪进传给 _dispatch_persona 的 on_start 回调，仅当 _work 真正执行（拿到槽位）
才触发。排队被取消 → 回调不触发 → ledger 干净 → 下一 tick 由 _iid_in_flight 放行后自然重派。

Standalone: `python3 bridge/test_persona_ledger_on_start.py`. 无 a1/Aone/DingTalk（fake pool）。
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

IID = "84297352"
CID = 124851623


class FakePool:
    """最小 EphemeralExecutor 替身。start_worker=True 时 submit 会同步执行 work（模拟拿到槽位
    立即开跑）；False 时只入队不执行（模拟槽满排队，随后被重启 cancel_futures 丢弃）。"""

    def __init__(self, start_worker):
        self.start_worker = start_worker
        self.submitted = []

    def active_ids(self):
        return []

    def set_proc(self, *_a, **_k):
        pass

    def submit(self, key, work, **_kw):
        self.submitted.append(key)
        if self.start_worker:
            work()          # worker 真正开跑 → 触发 on_start 落盘
        return True, "dispatched"


def _make_scheduler(pool):
    tmp = tempfile.NamedTemporaryFile(suffix=".json", delete=False, mode="w")
    tmp.write('{"tickets": {}}')
    tmp.close()
    # handler=None → _work 里 self.handler is None，落盘后直接返回 "done"，不碰 dispatch_item。
    return bot.PersonaScheduler(handler=None, pool=pool, enabled=False,
                                ledger_path=tmp.name)


def _item():
    return {"id": IID, "title": "TDE breaks tair", "pool_project": "1086837"}


def _dispatch_decision():
    return {"action": "dispatch", "internal_role": "terraform-rd",
            "public_identity": "terraform-rd", "reason": "human_mention",
            "comment_id": CID, "handoff": {"round": 1}}


class LedgerOnStartTest(unittest.TestCase):
    def test_queued_but_not_started_leaves_ledger_clean(self):
        """submit 接受但 worker 未启动（排队后被重启丢弃）→ processed 不写、last_seen 不推，
        下一 tick 才能自然重派。这是修复的核心断言。"""
        sched = _make_scheduler(FakePool(start_worker=False))
        sched._apply_decisions(_item(), [], [_dispatch_decision()])
        st = sched._get_ticket_state(IID)
        self.assertEqual(st.get("processed") or set(), set(),
                         "排队未启动却把评论标 processed = 复现 84297352 黑洞")
        self.assertEqual(int(st.get("last_seen") or 0), 0)

    def test_worker_started_commits_ledger(self):
        """worker 真正开跑（on_start 触发）→ processed/last_seen/count 落盘，防同评论重复派。"""
        sched = _make_scheduler(FakePool(start_worker=True))
        sched._apply_decisions(_item(), [], [_dispatch_decision()])
        st = sched._get_ticket_state(IID)
        self.assertIn(CID, st.get("processed") or set())
        self.assertEqual(int(st.get("last_seen") or 0), CID)
        # human_mention → 计数重置为 1、escalated 清零
        self.assertEqual(int(st.get("dispatch_count") or 0), 1)
        self.assertFalse(st.get("escalated"))

    def test_skip_still_advances_last_seen(self):
        """skip 决策仍即时推 last_seen（原逻辑不变，不受本次改动影响）。"""
        sched = _make_scheduler(FakePool(start_worker=False))
        skip = {"action": "skip", "reason": "persona_no_sentinel", "comment_id": CID}
        sched._apply_decisions(_item(), [], [skip])
        st = sched._get_ticket_state(IID)
        self.assertEqual(int(st.get("last_seen") or 0), CID)
        self.assertEqual(st.get("processed") or set(), set())


if __name__ == "__main__":
    unittest.main(verbosity=2)
