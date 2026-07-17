#!/usr/bin/env python3
"""Regression: PersonaScheduler advances its ledger only after durable acceptance.

Bug (工单 84297352, 原根 @Terraform-研发数字人 两次都无人理):
_apply_decisions 旧逻辑在 self.pool.submit(...) 返回 (accepted=True) 后就同步把评论写进
ledger.processed。但 submit 只表示"入队接受"——EphemeralExecutor 槽满时 future 仅排队。此刻
bridge 换 token 重启，terminate_all 的 shutdown(cancel_futures=True) 丢弃未启动的排队
future（worker 从未跑、评论从未回），而 ledger 已标 processed → PersonaScheduler 永久
skip(reason=processed) 不再重派，ScanScheduler 冷启动又把它当存量积压 → 双重盲区烂尾。

Task-only 架构中，控制面成功持久化就等价于可靠接管；未持久化时 ledger 必须保持干净，
下一 tick 才能自然重试。EphemeralJob 仍以本地 worker on_start 作为确认边界。

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


class FakeTaskClient:
    def upsert_desired_task(self, envelope, request_id=None):
        return {"accepted": True, "reason": "task_persisted"}


def _make_scheduler(pool, task_client=None):
    tmp = tempfile.NamedTemporaryFile(suffix=".json", delete=False, mode="w")
    tmp.write('{"tickets": {}}')
    tmp.close()
    # handler=None → _work 里 self.handler is None，落盘后直接返回 "done"，不碰 dispatch_item。
    scheduler = bot.PersonaScheduler(handler=None, pool=pool, enabled=False,
                                     ledger_path=tmp.name)
    scheduler.execution_router = bot.ExecutionRouter(client=task_client)
    return scheduler


def _item():
    return {"id": IID, "title": "TDE breaks tair", "pool_project": "1086837"}


def _dispatch_decision():
    return {"action": "dispatch", "internal_role": "terraform-rd",
            "public_identity": "terraform-rd", "reason": "human_mention",
            "comment_id": CID, "handoff": {"round": 1}}


class LedgerOnStartTest(unittest.TestCase):
    def test_task_persistence_failure_leaves_ledger_clean(self):
        """Task 未持久化时不推进传感器游标，下一 tick 仍可重试。"""
        sched = _make_scheduler(FakePool(start_worker=False), task_client=None)
        sched._apply_decisions(_item(), [], [_dispatch_decision()])
        st = sched._get_ticket_state(IID)
        self.assertEqual(st.get("processed") or set(), set(),
                         "Task 未持久化却把评论标 processed")
        self.assertEqual(int(st.get("last_seen") or 0), 0)

    def test_durable_task_acceptance_commits_ledger(self):
        """Task 已持久化即推进游标；数据库负责后续 lease/retry。"""
        sched = _make_scheduler(
            FakePool(start_worker=False), task_client=FakeTaskClient())
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
