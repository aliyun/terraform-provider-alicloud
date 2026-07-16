#!/usr/bin/env python3
"""Unit: RecoveryScheduler 死任务重委派（工单 84386065 需求三，方案 A′）。

Gap：交互 Codex/Claude worker 定向接单（dispatch.pull=false）后会话死亡 → lease 过期 →
reaper 收敛为 SHADOW+RESUMABLE / SHADOW+无 current session（CORRUPTED 归档）——交互单被
hasInteractiveLineage 永久挡在通用 Task 队列外，scan 又被 jarvis-claimed 标签 skip，
此前无人接手。恢复正门 = spawn headless jarvis（EphemeralJob 只是进程外壳）经 claim.sh
的 fenced targeted claimTask 接管。本测试锁住：候选枚举（STALE/OFFLINE worker 可见残留
assignment + 台账生前记忆；纯 scanner shadow 观察不进候选）、get_task_by_aone/timeline
佐证决策表（含 SHADOW 无 session 可重派）、recovery_policy 三分流（REPLAY_SAFE 重派 /
RESUME_ONLY·MANUAL 只播报）、JARVIS_RECOVERY_REDISPATCH=0 减档、recovery.json 幂等防抖
（dedup TTL / MAX_ROUNDS 升级 / 告警按原因只播一次）、pause 复用与 client 缺失自动禁用。

Standalone: `python3 bridge/test_recovery_scheduler.py`. 无控制面/网络（fake client/pool）。
"""
import importlib.util
import json
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock

HERE = Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location(
    "jarvis_dingtalk_bot", HERE / "jarvis_dingtalk_bot.py")
bot = importlib.util.module_from_spec(spec)
sys.modules["jarvis_dingtalk_bot"] = bot
spec.loader.exec_module(bot)

AONE = "84386065"
PROJ = "2100304"
KEY = "aone:%s:%s" % (PROJ, AONE)
WKEY = "interactive:codex:macmini:dead0001"


class FakeHandler:
    def __init__(self):
        self.broadcasts = []
        self.dispatched = []

    def _broadcast(self, text):
        self.broadcasts.append(text)

    def dispatch_item(self, item_id, prompt, sid, resume, notify, target, ttype, **kw):
        self.dispatched.append({"id": item_id, "kind": kw.get("kind"),
                                "project": kw.get("project"), "prompt": prompt,
                                "terraform": kw.get("terraform")})
        return "done"


class FakePool:
    def __init__(self, accept=True):
        self.accept = accept
        self.submitted = []

    def submit(self, key, work, **kw):
        self.submitted.append({"key": key, "force": kw.get("force"),
                               "kind": kw.get("kind"), "project": kw.get("project")})
        if self.accept:
            work()
            return True, "dispatched"
        return False, "queue_full"

    def set_proc(self, key, proc):
        pass


class FakeClient:
    """服务端契约 fixture（对齐 AutomationAgent WorkerStateResponse/TaskView/Timeline）：
    /workers → [{worker, activityStatus, assignments:[{task, session}]}]；
    /tasks/by-aone/{id} → TaskView 数组；/tasks/{id}/timeline → {task, currentWorker,
    sessions, events, operations}。"""

    def __init__(self, workers=None, tasks=None, timelines=None):
        self.workers = workers if workers is not None else []
        self.tasks = tasks or {}          # aone_id -> [task dict]
        self.timelines = timelines or {}  # task_id str -> timeline dict
        self.calls = []

    def list_workers(self):
        self.calls.append(("list_workers",))
        return self.workers

    def get_task_by_aone(self, aone_id):
        self.calls.append(("by_aone", str(aone_id)))
        return self.tasks.get(str(aone_id), [])

    def get_task_timeline(self, task_id):
        self.calls.append(("timeline", str(task_id)))
        return self.timelines.get(str(task_id), {})


def _task(status="SHADOW", task_id=11, session_id=7, key=KEY, payload=None,
          policy=None):
    task = {"id": task_id, "taskKey": key, "aoneId": AONE, "status": status,
            "executionMode": "SHADOW", "currentSessionId": session_id,
            "payload": payload if payload is not None else {
                "itemId": AONE, "project": PROJ, "kind": "ticket",
                "title": "客户工单标题", "poolKey": "tf_customer"}}
    if policy is not None:
        task["recoveryPolicy"] = policy
    return task


def _assignment(task=None, session_status="RUNNING", session_id=7, fence=3):
    return {"task": task if task is not None else _task(status="RUNNING"),
            "session": {"id": session_id, "taskId": 11, "fenceToken": fence,
                        "status": session_status}}


def _worker_entry(act, assignments=(), wkey=WKEY):
    return {"worker": {"workerKey": wkey, "capabilities": {"client": "codex"}},
            "activityStatus": act,
            "assignments": list(assignments)}


def _timeline(session_status="RESUMABLE", session_id=7, current_worker=None,
              events=(), operations=()):
    return {"task": {"id": 11},
            "currentWorker": current_worker,
            "sessions": [{"id": session_id, "status": session_status,
                          "fenceToken": 4, "attemptNo": 2}],
            "events": list(events), "operations": list(operations)}


class RecoverySchedulerTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.ledger = Path(self.tmp.name) / "recovery.json"
        # 隔离 pause 闸门与任何 .my-day 路径到临时目录（REPO_ROOT 在 _tick 调用时读取）。
        self._orig_root = bot.REPO_ROOT
        bot.REPO_ROOT = Path(self.tmp.name)
        self.handler = FakeHandler()
        self.pool = FakePool()

    def tearDown(self):
        bot.REPO_ROOT = self._orig_root
        self.tmp.cleanup()

    def _sched(self, client, **env):
        with mock.patch.dict(bot.os.environ, env, clear=False):
            sched = bot.RecoveryScheduler(self.handler, self.pool, client=client,
                                          ledger_path=self.ledger)
        return sched

    def _ledger_data(self):
        return json.loads(self.ledger.read_text())

    # -- 候选发现与重派 --------------------------------------------------------

    def test_stale_worker_dead_assignment_redispatches_and_records_ledger(self):
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW")]},
            timelines={"11": _timeline("RESUMABLE")})
        self._sched(client)._tick()
        self.assertEqual(len(self.pool.submitted), 1)
        sub = self.pool.submitted[0]
        self.assertEqual(sub["key"], AONE)
        self.assertTrue(sub["force"], "重派必须 force 越过 24h 派发去重台账")
        self.assertEqual(sub["kind"], "ticket")
        self.assertEqual(sub["project"], PROJ)
        self.assertEqual(self.handler.dispatched[0]["id"], AONE)
        data = self._ledger_data()
        self.assertEqual(data["pending"][AONE]["count"], 1)
        self.assertIn("last_ts", data["pending"][AONE])
        self.assertTrue(any("已重委派 #%s" % AONE in b for b in self.handler.broadcasts))

    def test_ready_task_is_redispatched_without_timeline_lookup(self):
        client = FakeClient(
            workers=[_worker_entry("OFFLINE", [_assignment()])],
            tasks={AONE: [_task(status="READY")]})
        self._sched(client)._tick()
        self.assertEqual(len(self.pool.submitted), 1)
        self.assertNotIn(("timeline", "11"), client.calls,
                         "READY 直接可派，不需要 timeline 佐证")

    def test_memory_covers_assignment_hidden_after_lease_expiry(self):
        """核心适配：/workers 的 assignments 只含 lease 未过期的条目——worker 死后
        assignment 从响应消失，靠台账「生前记忆」仍能追查重派。"""
        client = FakeClient(
            workers=[_worker_entry("BUSY", [_assignment()])],
            tasks={AONE: [_task(status="RUNNING")]},
            timelines={"11": _timeline("RUNNING",
                                       current_worker={"workerKey": WKEY})})
        sched = self._sched(client)
        sched._tick()  # 活 worker：只记忆，不派
        self.assertEqual(self.pool.submitted, [])
        self.assertIn(AONE, self._ledger_data()["workers"][WKEY]["assign"])
        # worker 死透：assignments 已被服务端过滤为空，task 已被 reaper 收敛。
        client.workers = [_worker_entry("STALE", [])]
        client.tasks = {AONE: [_task(status="SHADOW")]}
        client.timelines = {"11": _timeline("RESUMABLE")}
        sched._tick()
        self.assertEqual(len(self.pool.submitted), 1)
        self.assertEqual(self.pool.submitted[0]["key"], AONE)

    # -- 不重派的路径 ----------------------------------------------------------

    def test_recovery_required_announces_once_never_dispatches(self):
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="RECOVERY_REQUIRED")]})
        sched = self._sched(client)
        sched._tick()
        self.assertEqual(self.pool.submitted, [],
                         "RECOVERY_REQUIRED 严禁自动重派（服务端 claim 会 409）")
        alerts = [b for b in self.handler.broadcasts if "RECOVERY_REQUIRED" in b]
        self.assertEqual(len(alerts), 1)
        self.assertIn("readback", alerts[0])
        sched._tick()  # 二次 tick 只播一次（announced 台账去重）
        alerts = [b for b in self.handler.broadcasts if "RECOVERY_REQUIRED" in b]
        self.assertEqual(len(alerts), 1)

    def test_healthy_worker_is_untouched(self):
        client = FakeClient(
            workers=[_worker_entry("BUSY", [_assignment()]),
                     _worker_entry("IDLE", [], wkey="interactive:claude:ok:1")])
        self._sched(client)._tick()
        self.assertEqual(self.pool.submitted, [])
        self.assertEqual(self.handler.broadcasts, [])
        self.assertNotIn(("by_aone", AONE), client.calls,
                         "健康 worker 不触发佐证查询")

    def test_lease_not_yet_expired_waits_for_next_tick(self):
        """死 worker 但 lease 未过期（task 仍 RUNNING、currentWorker 还是它）→ 本轮不动。"""
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="RUNNING")]},
            timelines={"11": _timeline("RUNNING",
                                       current_worker={"workerKey": WKEY})})
        self._sched(client)._tick()
        self.assertEqual(self.pool.submitted, [])
        # 候选保留在台账（生前记忆），下轮 reaper 收敛后可重派。
        self.assertIn(AONE, self._ledger_data()["workers"][WKEY]["assign"])

    def test_task_taken_over_by_another_worker_resolves(self):
        client = FakeClient(
            workers=[_worker_entry("OFFLINE", [_assignment()])],
            tasks={AONE: [_task(status="RUNNING")]},
            timelines={"11": _timeline(
                "RUNNING", current_worker={"workerKey": "interactive:other:host:2"})})
        self._sched(client)._tick()
        self.assertEqual(self.pool.submitted, [])
        data = self._ledger_data()
        self.assertEqual(data["pending"], {})
        self.assertEqual(data["workers"], {}, "被接管后生前记忆应摘除")

    def test_shadow_without_current_session_is_redispatched(self):
        """REPLAY_SAFE 死亡的 CORRUPTED 归档路径：reaper 剥离 current session →
        SHADOW 裸态，targeted claim 可直接开新 fenced Session → 可重派。"""
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW", session_id=None)]})
        self._sched(client)._tick()
        self.assertEqual(len(self.pool.submitted), 1)
        self.assertEqual(self.pool.submitted[0]["key"], AONE)
        self.assertNotIn(("timeline", "11"), client.calls,
                         "无 current session 不需要 timeline 佐证")

    def test_suspended_task_resolves(self):
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SUSPENDED")]})
        self._sched(client)._tick()
        self.assertEqual(self.pool.submitted, [],
                         "SUSPENDED（挂起等人·亲和期）不该重派")
        self.assertEqual(self._ledger_data()["workers"], {}, "应摘除生前记忆")

    def test_pure_scanner_shadow_task_never_becomes_candidate(self):
        """控制面里存在的 SHADOW 观察任务（从未被任何 worker claim）不进候选：
        候选只来自死 worker 的 assignment/生前记忆。"""
        client = FakeClient(
            workers=[_worker_entry("OFFLINE", []),
                     _worker_entry("IDLE", [], wkey="interactive:claude:ok:1")],
            tasks={AONE: [_task(status="SHADOW", session_id=None)]})
        self._sched(client)._tick()
        self.assertEqual(self.pool.submitted, [])
        self.assertNotIn(("by_aone", AONE), client.calls,
                         "无 assignment/记忆 → 不得触发佐证查询")

    # -- recovery_policy 分流 ----------------------------------------------------

    def test_resume_only_policy_announces_once_no_dispatch(self):
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW", policy="RESUME_ONLY")]},
            timelines={"11": _timeline("RESUMABLE")})
        sched = self._sched(client)
        sched._tick()
        self.assertEqual(self.pool.submitted, [],
                         "RESUME_ONLY 仅原 runtime 可续跑，不得重派")
        alerts = [b for b in self.handler.broadcasts if "RESUME_ONLY" in b]
        self.assertEqual(len(alerts), 1)
        self.assertEqual(self._ledger_data()["pending"][AONE]["announced"],
                         "resume_only")
        sched._tick()  # 同因只播一次
        self.assertEqual(
            len([b for b in self.handler.broadcasts if "RESUME_ONLY" in b]), 1)

    def test_manual_policy_announces_once_no_dispatch(self):
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW", policy="MANUAL")]},
            timelines={"11": _timeline("RESUMABLE")})
        sched = self._sched(client)
        sched._tick()
        sched._tick()
        self.assertEqual(self.pool.submitted, [], "MANUAL 按策略只转人工")
        alerts = [b for b in self.handler.broadcasts if "MANUAL" in b]
        self.assertEqual(len(alerts), 1)
        self.assertIn("转人工", alerts[0])

    def test_explicit_replay_safe_policy_redispatches(self):
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW", policy="REPLAY_SAFE")]},
            timelines={"11": _timeline("RESUMABLE")})
        self._sched(client)._tick()
        self.assertEqual(len(self.pool.submitted), 1)

    def test_redispatch_off_alerts_only(self):
        """JARVIS_RECOVERY_REDISPATCH=0 减档：探测/佐证/播报/台账照常但不 spawn。"""
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW")]},
            timelines={"11": _timeline("RESUMABLE")})
        sched = self._sched(client, JARVIS_RECOVERY_REDISPATCH="0")
        sched._tick()
        self.assertEqual(self.pool.submitted, [], "减档模式不得 spawn")
        self.assertEqual(self.handler.dispatched, [])
        alerts = [b for b in self.handler.broadcasts
                  if "JARVIS_RECOVERY_REDISPATCH=0" in b]
        self.assertEqual(len(alerts), 1)
        data = self._ledger_data()
        self.assertEqual(data["pending"][AONE]["announced"], "alert_only")
        self.assertNotIn("count", data["pending"][AONE],
                         "减档告警不得消耗重派轮次预算")
        self.assertIn(AONE, data["workers"][WKEY]["assign"],
                      "生前记忆照常维护（探测不受减档影响）")
        sched._tick()  # 同因只播一次
        self.assertEqual(
            len([b for b in self.handler.broadcasts
                 if "JARVIS_RECOVERY_REDISPATCH=0" in b]), 1)

    # -- 幂等防抖 ---------------------------------------------------------------

    def test_second_tick_within_dedup_ttl_does_not_redispatch(self):
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW")]},
            timelines={"11": _timeline("RESUMABLE")})
        sched = self._sched(client)
        sched._tick()
        sched._tick()
        self.assertEqual(len(self.pool.submitted), 1,
                         "dedup TTL 内同单不得重复重派")
        self.assertEqual(self._ledger_data()["pending"][AONE]["count"], 1)

    def test_expired_cooldown_allows_next_round(self):
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW")]},
            timelines={"11": _timeline("RESUMABLE")})
        sched = self._sched(client, JARVIS_RECOVERY_DEDUP_TTL="0")
        sched._tick()
        sched._tick()
        self.assertEqual(len(self.pool.submitted), 2)
        self.assertEqual(self._ledger_data()["pending"][AONE]["count"], 2)

    def test_max_rounds_escalates_once_and_stops_auto(self):
        self.ledger.write_text(json.dumps(
            {"pending": {AONE: {"count": 3, "last_ts": 0}}, "workers": {}}))
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW")]},
            timelines={"11": _timeline("RESUMABLE")})
        sched = self._sched(client, JARVIS_RECOVERY_MAX_ROUNDS="3")
        sched._tick()
        self.assertEqual(self.pool.submitted, [], "超 MAX_ROUNDS 不再自动重派")
        esc = [b for b in self.handler.broadcasts if "超上限" in b]
        self.assertEqual(len(esc), 1)
        self.assertTrue(self._ledger_data()["pending"][AONE]["escalated"])
        sched._tick()  # 升级播报只发一次
        esc = [b for b in self.handler.broadcasts if "超上限" in b]
        self.assertEqual(len(esc), 1)

    def test_max_per_tick_caps_dispatches(self):
        entries, tasks, timelines = [], {}, {}
        for i in range(3):
            aone = "9000%d" % i
            t = _task(status="SHADOW", task_id=20 + i,
                      key="aone:%s:%s" % (PROJ, aone))
            t["aoneId"] = aone
            entries.append(_assignment(task=t, session_id=7))
            tasks[aone] = [t]
            timelines[str(20 + i)] = _timeline("RESUMABLE")
        client = FakeClient(
            workers=[_worker_entry("OFFLINE", entries)],
            tasks=tasks, timelines=timelines)
        self._sched(client, JARVIS_RECOVERY_MAX_PER_TICK="2")._tick()
        self.assertEqual(len(self.pool.submitted), 2)

    def test_pool_rejection_does_not_consume_round_budget(self):
        self.pool = FakePool(accept=False)
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment()])],
            tasks={AONE: [_task(status="SHADOW")]},
            timelines={"11": _timeline("RESUMABLE")})
        self._sched(client)._tick()
        self.assertNotIn(AONE, self._ledger_data()["pending"],
                         "提交被拒不该记 count/last_ts（下轮原样重试）")

    # -- 开关 / 暂停 / 容错 ------------------------------------------------------

    def test_missing_client_disables_scheduler(self):
        sched = bot.RecoveryScheduler(self.handler, self.pool, client=None,
                                      ledger_path=self.ledger)
        self.assertFalse(sched.enabled)
        sched.start()
        self.assertIsNone(sched._thread)

    def test_env_zero_disables_even_with_client(self):
        sched = self._sched(FakeClient(), JARVIS_RECOVERY_SCHED="0")
        self.assertFalse(sched.enabled)
        sched.start()
        self.assertIsNone(sched._thread)

    def test_pause_flag_skips_tick_entirely(self):
        pause = Path(self.tmp.name) / ".my-day" / "bridge" / "pause"
        pause.parent.mkdir(parents=True, exist_ok=True)
        pause.write_text("")
        client = FakeClient(workers=[_worker_entry("STALE", [_assignment()])])
        self._sched(client)._tick()
        self.assertEqual(client.calls, [], "pause 期间不得触碰控制面")
        self.assertEqual(self.pool.submitted, [])

    def test_list_workers_failure_is_swallowed(self):
        class BoomClient(FakeClient):
            def list_workers(self):
                raise RuntimeError("boom")

        sched = self._sched(BoomClient())
        sched._tick()  # 不抛异常即通过
        self.assertEqual(self.pool.submitted, [])

    def test_non_canonical_task_key_is_ignored(self):
        t = _task(status="SHADOW", key="probe-2026-07-16")
        client = FakeClient(workers=[_worker_entry("STALE", [_assignment(task=t)])])
        self._sched(client)._tick()
        self.assertEqual(self.pool.submitted, [])
        self.assertNotIn(("by_aone", AONE), client.calls)

    # -- 派发内容 ---------------------------------------------------------------

    def test_dispatch_uses_payload_prompt_when_present(self):
        payload = {"itemId": AONE, "project": PROJ, "kind": "ticket",
                   "title": "T", "poolKey": "tf_customer", "prompt": "resume-me",
                   "terraform": True}
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment(task=_task(payload=payload))])],
            tasks={AONE: [_task(status="SHADOW", payload=payload)]},
            timelines={"11": _timeline("RESUMABLE")})
        self._sched(client)._tick()
        d = self.handler.dispatched[0]
        self.assertEqual(d["prompt"], "resume-me")
        self.assertTrue(d["terraform"])

    def test_dispatch_synthesizes_ticket_prompt_for_sparse_interactive_payload(self):
        """交互 claim 的 payload 只有 itemId/project/kind/trigger → 用 _ticket_prompt 合成。"""
        sparse = {"itemId": AONE, "project": PROJ, "kind": "ticket",
                  "trigger": "INTERACTIVE"}
        client = FakeClient(
            workers=[_worker_entry("STALE", [_assignment(task=_task(payload=sparse))])],
            tasks={AONE: [_task(status="SHADOW", payload=sparse)]},
            timelines={"11": _timeline("RESUMABLE")})
        self._sched(client)._tick()
        d = self.handler.dispatched[0]
        self.assertIn(AONE, d["prompt"])
        self.assertIn("claim.sh claim %s %s" % (AONE, PROJ), d["prompt"])


if __name__ == "__main__":
    unittest.main()
