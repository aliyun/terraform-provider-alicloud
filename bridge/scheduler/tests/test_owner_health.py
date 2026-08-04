from __future__ import annotations

import logging
import tempfile
import unittest
from pathlib import Path
from typing import Any
from unittest import mock

from bridge.scheduler.runners import owner_health

_ABSENT = object()


class FakeClient:
    def __init__(self, *, now: float) -> None:
        self.now = now
        self.ready_cursors: list[int] = []
        self.source_cursors: list[int] = []
        # Default session carries a transcript, so it is NOT auto-releasable:
        # there is real work to preserve and a human must decide. Tests that
        # exercise the auto-release path clear these refs explicitly.
        self.session_refs: dict[str, Any] = {
            "transcriptUri": "/tmp/session-33.jsonl",
            "branchRef": "worktree-84550003",
            "checkpointUri": None,
        }
        self.release_calls: list[dict[str, Any]] = []
        self.release_action = "RELEASED"
        self.payload_policy: Any = owner_health.HEADLESS_POLICY_REVISION
        # desired == processing exercises the discard-resume branch; a pending
        # newer desired revision exercises force-release. Default keeps the
        # original force-release assertions intact.
        self.desired_revision: Any = "rev-new"
        # Aone status of the source work item. Non-terminal by default so the
        # existing stranded-detection assertions keep exercising the full path.
        self.source_status: str = "Open"
        self.point_read_aone_ids: list[str] = []

    def list_ready_task_diagnostics(
            self, *, after_task_id: int = 0, limit: int = 100):
        self.ready_cursors.append(after_task_id)
        rows = {
            0: [{
                "task": {
                    "id": 1,
                    "aoneId": "84550001",
                    "status": "READY",
                    "currentSessionId": 11,
                },
                "reasonCode": "RESUME_OWNER_UNAVAILABLE",
                "requiredWorkerKey": "worker-old-1",
                "requiredWorkerActivityStatus": "OFFLINE",
                "reasonSince": self.now - 3600,
            }],
            1: [{
                "task": {
                    "id": 2,
                    "aoneId": "84550002",
                    "status": "READY",
                    "currentSessionId": 22,
                },
                "reasonCode": "RESUME_OWNER_NOT_REGISTERED",
                "requiredWorkerKey": "worker-old-2",
                "requiredWorkerActivityStatus": "NOT_REGISTERED",
            }],
            2: [],
        }
        return {"items": rows.get(after_task_id, [])}

    def list_source_status_candidates(
            self, *, after_task_id: int = 0, limit: int = 100):
        # Real API only returns 4 fields (Aone source status): taskId /
        # sourceProjectKey / aoneId / sourceStatus. NOT the Task control-plane
        # status / recoveryPolicy / currentSessionId — those require a separate
        # get_task_by_aone point-read. This fake mirrors that shape so the test
        # exercises the production code path.
        self.source_cursors.append(after_task_id)
        if after_task_id == 0:
            return {
                "items": [{
                    "taskId": 3,
                    "sourceProjectKey": "2124589",
                    "aoneId": "84550003",
                    "sourceStatus": self.source_status,
                }],
                "nextAfterTaskId": 3,
            }
        return {"items": []}

    def get_task_by_aone(self, aone_id: str):
        self.point_read_aone_ids.append(aone_id)
        if aone_id == "84550003":
            task = {
                "id": 3,
                "aoneId": "84550003",
                "generation": 1,
                "status": "RECOVERY_REQUIRED",
                "recoveryPolicy": "RESUME_ONLY",
                "currentSessionId": 33,
            }
            # The real point-read carries the frozen payload; the executor
            # refuses a Task whose policyRevision is not the current one, so
            # auto-release has to see it. `None` models an API that stopped
            # returning it at all.
            if self.payload_policy is not _ABSENT:
                task["payload"] = {"policyRevision": self.payload_policy}
            return [task]
        return []

    def get_task_timeline(self, task_id: str):
        assert task_id == "3"
        session = {
            "id": 33,
            "status": "RESUMABLE",
            "workerKey": "worker-old-3",
            "historicalWorkerKey": "worker-old-3",
            "historicalWorkerId": 565,
            "historicalWorkerProcessUuid": "proc-old-3",
            "fenceToken": 4,
        }
        session.update(self.session_refs)
        return {
            # The real endpoint returns the Task snapshot alongside the sessions;
            # force-release reads its CAS fields from here.
            "task": {
                "id": 3,
                "aoneId": "84550003",
                "status": "RECOVERY_REQUIRED",
                "recoveryPolicy": "RESUME_ONLY",
                "currentSessionId": 33,
                "generation": 1,
                "stateVersion": 22,
                "retryCount": 4,
                "maxRetries": 3,
                "desiredRevision": self.desired_revision,
                "processingRevision": "rev-a",
            },
            "sessions": [session],
            "currentWorker": None,
            "events": [{
                "eventType": "TASK_OWNER_MIGRATION",
                "reasonCode": "RESUME_OWNER_UNAVAILABLE",
                "occurredAt": self.now - 7200,
            }],
        }

    def force_release_task(self, task_id: str, **kwargs):
        self.release_calls.append(dict(kwargs, task_id=task_id, via="force_release"))
        # BLOCKED means the Task stays RECOVERY_REQUIRED; RELEASED means READY.
        status = "RECOVERY_REQUIRED" if self.release_action == "BLOCKED" else "READY"
        return {"action": self.release_action, "status": status,
                "message": "fake release"}

    def discard_resume_context(self, task_id: str, expected_session_id: int,
                               reason: str = "", request_id=None):
        self.release_calls.append({
            "task_id": task_id, "expected_session_id": expected_session_id,
            "reason": reason, "via": "discard_resume"})
        # discard-resume does not return an `action`; it returns the Task
        # already at READY (verified on task 3312).
        return {"status": "READY", "message": "resume context discarded"}


class OwnerHealthRunnerTest(unittest.TestCase):
    def setUp(self):
        self.now = 1_800_000_000.0
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.client = FakeClient(now=self.now)
        patcher = mock.patch.dict(
            owner_health.os.environ,
            {
                "JARVIS_OWNER_HEALTH_PAGE_SIZE": "1",
                "JARVIS_OWNER_HEALTH_REPEAT_SEC": "3600",
            },
        )
        patcher.start()
        self.addCleanup(patcher.stop)
        self.runner = owner_health.OwnerHealthRunner(
            task_client=self.client,
            logger=logging.getLogger(__name__),
            repo_root=self.tmp.name,
            state_path=Path(self.tmp.name) / "owner-health.json",
            clock=lambda: self.now,
        )

    def test_pages_ready_and_recovery_inventories(self):
        blockers = self.runner.collect_blockers()

        self.assertEqual(set(blockers), {1, 2, 3})
        self.assertEqual(
            self.client.ready_cursors,
            [0, 1, 2, 0, 1, 2],
            "READY inventory is paged before and after migration inventory",
        )
        self.assertEqual(self.client.source_cursors, [0, 3])
        self.assertEqual(blockers[3]["status"], "RECOVERY_REQUIRED")
        self.assertEqual(blockers[3]["source"], "recovery-migration")

    def test_terminal_source_item_is_not_a_recovery_blocker(self):
        """源工单已终态的 Task 无任何可执行的恢复动作，不该反复告警。

        服务端对这类 Task 直接拒绝 force-release（"tasks whose source item is
        terminal cannot be force released"），而 desired != processing 时
        discard-resume 也不合法，于是它会永久卡在 RECOVERY_REQUIRED 并按小时告警，
        人却无从处理。源工单既已完结，本身也无事可做。
        """
        self.client.source_status = "已完成"

        blockers = self.runner.collect_blockers()

        self.assertNotIn(3, blockers, "终态源工单不应产生 recovery-migration blocker")
        self.assertEqual(set(blockers), {1, 2}, "READY 诊断链路不受影响")
        self.assertNotIn(
            "84550003", self.client.point_read_aone_ids,
            "应在 point-read 之前就跳过，不浪费一次 Task 查询")

    def test_terminal_source_skip_logs_once_not_per_entry(self):
        """每个终态候选一行 info 会盖过它省下的告警。

        skip 分支对**每个**终态源候选都命中，而它消掉的噪声只是其中
        RECOVERY_REQUIRED 的子集。线上实测 161 行/轮（300s 间隔 ≈ 1900 行/小时），
        而被替掉的告警只有约 124 行/小时——净增约 15 倍。聚合成一行保留规模可观测性，
        明细降到 debug 供定位单条时取回。
        """
        self.client.source_status = "已完成"
        self.runner.logger = mock.Mock()

        self.runner.collect_blockers()

        summaries = [call for call in self.runner.logger.info.call_args_list
                     if "terminal-source" in str(call)]
        self.assertEqual(len(summaries), 1, "每轮只应有一行汇总")
        self.assertIn(1, summaries[0].args, "汇总必须带被跳过的条数")
        self.assertFalse(
            any("skip terminal source aone=" in str(call)
                for call in self.runner.logger.info.call_args_list),
            "逐条明细不得留在 info")
        self.assertTrue(
            any("skip terminal source aone=" in str(call)
                for call in self.runner.logger.debug.call_args_list),
            "逐条明细应仍可通过 debug 取回")

    def test_no_summary_line_when_nothing_was_skipped(self):
        self.client.source_status = "Open"
        self.runner.logger = mock.Mock()

        self.runner.collect_blockers()

        self.assertFalse(
            any("terminal-source" in str(call)
                for call in self.runner.logger.info.call_args_list))

    def test_every_terminal_status_is_skipped(self):
        for status in sorted(owner_health.TERMINAL_STATUSES):
            with self.subTest(status=status):
                self.client.source_status = status
                self.client.source_cursors.clear()
                self.assertNotIn(3, self.runner.collect_blockers())

    def test_non_terminal_source_item_still_pages(self):
        """反向保护：源工单未终态时仍按原样识别为阻塞项。"""
        for status in ("Open", "开发中", "待处理", ""):
            with self.subTest(status=status):
                self.client.source_status = status
                blockers = self.runner.collect_blockers()
                self.assertIn(3, blockers)
                self.assertEqual(blockers[3]["status"], "RECOVERY_REQUIRED")
                self.assertEqual(blockers[3]["source"], "recovery-migration")

    def test_recovery_migration_alert_records_dwell_duration(self):
        blockers = self.runner.collect_blockers()
        calls = []
        with mock.patch.object(
                owner_health, "_dingtalk_event_enqueue",
                side_effect=lambda *args, **kwargs: calls.append(
                    (args, kwargs)) or True):
            queued = self.runner.reconcile({3: blockers[3]})

        self.assertEqual(queued, 1)
        args, kwargs = calls[0]
        self.assertEqual(args[0], "84550003")
        self.assertIn("RECOVERY_REQUIRED", args[5])
        self.assertIn("已滞留：2小时0分钟", args[5])
        self.assertIn("recovery-migration", args[5])
        self.assertTrue(kwargs["allow_non_tf"])
        state = self.runner._load_state()["episodes"]["3"]
        self.assertEqual(state["firstSeenAt"], self.now - 7200)
        self.assertEqual(state["lastSeenAt"], self.now)

    def test_repeat_alert_uses_persisted_episode_age(self):
        blocker = {
            "task_id": 9,
            "aone_id": "84550009",
            "status": "READY",
            "reason": "RESUME_OWNER_UNAVAILABLE",
            "required_worker": "old",
            "required_worker_status": "OFFLINE",
            "session_id": "99",
            "source": "ready-diagnostics",
        }
        bodies = []
        with mock.patch.object(
                owner_health, "_dingtalk_event_enqueue",
                side_effect=lambda *args, **_kwargs: bodies.append(args[5]) or True):
            self.runner.reconcile({9: blocker})
            self.now += 3700
            self.runner.reconcile({9: blocker})

        self.assertEqual(len(bodies), 2)
        self.assertIn("已滞留：1小时1分钟", bodies[-1])

    def test_ready_to_recovery_migration_preserves_episode_dwell(self):
        ready = {
            "task_id": 9,
            "aone_id": "84550009",
            "status": "READY",
            "reason": "RESUME_OWNER_UNAVAILABLE",
            "required_worker": "old",
            "required_worker_status": "OFFLINE",
            "session_id": "99",
            "source": "ready-diagnostics",
        }
        recovery = dict(
            ready, status="RECOVERY_REQUIRED", source="recovery-migration")
        with mock.patch.object(
                owner_health, "_dingtalk_event_enqueue", return_value=True):
            self.runner.reconcile({9: ready})
            self.now += 3700
            self.runner.reconcile({9: recovery})

        episode = self.runner._load_state()["episodes"]["9"]
        self.assertEqual(episode["firstSeenAt"], self.now - 3700)
        self.assertEqual(episode["status"], "RECOVERY_REQUIRED")

    def test_non_queue_pulling_owner_warns_once_with_manual_force_release(self):
        self.assertIn(
            "RESUME_OWNER_NOT_QUEUE_PULLING", owner_health.BLOCKING_REASONS)
        blocker = {
            "task_id": 12,
            "aone_id": "84550012",
            "status": "READY",
            "reason": "RESUME_OWNER_NOT_QUEUE_PULLING",
            "required_worker": "interactive:codex:old",
            "required_worker_status": "ACTIVE",
            "session_id": "112",
            "source": "ready-diagnostics",
        }
        bodies = []
        with mock.patch.object(
                owner_health, "_dingtalk_event_enqueue",
                side_effect=lambda *args, **_kwargs: bodies.append(args[5]) or True):
            self.runner.reconcile({12: blocker})
            self.now += 7200
            self.runner.reconcile({12: blocker})

        self.assertEqual(len(bodies), 1)
        self.assertIn(
            "control-plane-status.sh force-release 12 112", bodies[0])
        self.assertIn("不会自动释放 ownership", bodies[0])

    def test_recovery_blockers_handles_point_read_failure(self):
        # get_task_by_aone raises on one entry — must be skipped best-effort,
        # not propagate, and other entries (ready-diagnostics blockers 1, 2)
        # still processed.
        real_get = self.client.get_task_by_aone

        def flaky(aone_id):
            if aone_id == "84550003":
                raise RuntimeError("control-plane 503")
            return real_get(aone_id)

        with mock.patch.object(
                self.client, "get_task_by_aone", side_effect=flaky):
            blockers = self.runner.collect_blockers()
        self.assertNotIn(3, blockers)
        self.assertEqual(set(blockers), {1, 2})

    def test_recovery_blockers_skips_non_recovery_required(self):
        # get_task_by_aone returns a non-RECOVERY_REQUIRED Task — the recovery
        # pass filters it out; only ready-diagnostics blockers 1, 2 remain.
        with mock.patch.object(
                self.client, "get_task_by_aone",
                return_value=[{
                    "id": 3, "aoneId": "84550003", "generation": 1,
                    "status": "READY", "recoveryPolicy": "RESUME_ONLY",
                    "currentSessionId": 33,
                }]):
            blockers = self.runner.collect_blockers()
        self.assertNotIn(3, blockers)
        self.assertEqual(set(blockers), {1, 2})

    def test_recovery_blockers_skips_entry_without_aone_id(self):
        # An entry missing aoneId must be silently skipped without calling
        # get_task_by_aone(""); the real entry still flows through.
        original = self.client.list_source_status_candidates
        probe_calls: list[str] = []

        def with_blank(*args, **kwargs):
            result = original(*args, **kwargs)
            if result.get("items"):
                result["items"].append({
                    "taskId": 999,
                    "sourceProjectKey": "2124589",
                    "sourceStatus": "Open",
                })
            return result

        real_get = self.client.get_task_by_aone

        def tracking(aone_id):
            probe_calls.append(aone_id)
            return real_get(aone_id)

        with mock.patch.object(
                self.client, "list_source_status_candidates",
                side_effect=with_blank), \
                mock.patch.object(
                    self.client, "get_task_by_aone", side_effect=tracking):
            blockers = self.runner.collect_blockers()
        self.assertNotIn("", probe_calls)
        self.assertIn("84550003", probe_calls)
        self.assertIn(3, blockers)


class OwnerHealthAutoReleaseTest(unittest.TestCase):
    """A RESUMABLE session with nothing to resume is converged, not paged about.

    RESUME_ONLY pins a Task to one worker incarnation, so every bridge restart
    strands another batch. Before this the runner only alerted, so each stranded
    Task paged a human every repeat window forever while staying unrunnable.
    """

    def setUp(self):
        self.now = 1_800_000_000.0
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.client = FakeClient(now=self.now)
        # No transcript, branch or checkpoint: nothing to lose by releasing.
        self.client.session_refs = {
            "transcriptUri": None, "branchRef": None, "checkpointUri": None}
        patcher = mock.patch.dict(
            owner_health.os.environ,
            {"JARVIS_OWNER_HEALTH_PAGE_SIZE": "1",
             "JARVIS_OWNER_HEALTH_REPEAT_SEC": "3600"},
        )
        patcher.start()
        self.addCleanup(patcher.stop)
        self.runner = self._runner()

    def _runner(self):
        return owner_health.OwnerHealthRunner(
            task_client=self.client,
            logger=logging.getLogger(__name__),
            repo_root=self.tmp.name,
            state_path=Path(self.tmp.name) / "owner-health.json",
            clock=lambda: self.now,
        )

    def _blocker(self):
        return self.runner.collect_blockers()[3]

    def _reconcile(self, blocker):
        bodies = []
        with mock.patch.object(
                owner_health, "_dingtalk_event_enqueue",
                side_effect=lambda *args, **_kw: bodies.append(args[5]) or True):
            queued = self.runner.reconcile({3: blocker})
        return queued, bodies

    def test_empty_resume_context_is_released_instead_of_alerted(self):
        blocker = self._blocker()
        self.assertTrue(blocker["releasable"])
        queued, bodies = self._reconcile(blocker)

        self.assertEqual(queued, 0, "converged blockers must not page")
        self.assertEqual(bodies, [])
        self.assertEqual(len(self.client.release_calls), 1)
        call = self.client.release_calls[0]
        self.assertEqual(call["task_id"], "3")
        self.assertEqual(call["via"], "force_release")
        # Full CAS evidence, same snapshot shape the force-redispatch path uses.
        self.assertEqual(call["expected_session_id"], 33)
        self.assertEqual(call["expected_session_status"], "RESUMABLE")
        self.assertEqual(call["expected_generation"], 1)
        self.assertEqual(call["expected_state_version"], 22)
        self.assertEqual(call["expected_retry_count"], 4)
        self.assertEqual(call["expected_fence_token"], 4)
        self.assertIn("no transcript/branch/checkpoint", call["reason"])
        # Episode dropped, so a genuinely new blocker restarts its dwell clock.
        self.assertEqual(self.runner._load_state()["episodes"], {})

    def test_session_with_a_transcript_is_never_released(self):
        self.client.session_refs = {
            "transcriptUri": "/tmp/session-33.jsonl",
            "branchRef": None, "checkpointUri": None}
        blocker = self._blocker()
        self.assertFalse(blocker["releasable"])
        queued, bodies = self._reconcile(blocker)

        self.assertEqual(queued, 1)
        self.assertEqual(self.client.release_calls, [])
        self.assertIn("RECOVERY_REQUIRED", bodies[0])

    def test_release_waits_out_the_dwell_threshold(self):
        blocker = self._blocker()
        with mock.patch.dict(
                owner_health.os.environ,
                {"JARVIS_OWNER_HEALTH_RELEASE_AFTER_SEC": "86400"}):
            runner = self._runner()
            with mock.patch.object(
                    owner_health, "_dingtalk_event_enqueue", return_value=True):
                queued = runner.reconcile({3: blocker})

        # Blocker is only 2h old; a slow legitimate recovery may still land.
        self.assertEqual(self.client.release_calls, [])
        self.assertEqual(queued, 1)

    def test_kill_switch_restores_alert_only_behaviour(self):
        blocker = self._blocker()
        with mock.patch.dict(
                owner_health.os.environ,
                {"JARVIS_OWNER_HEALTH_AUTO_RELEASE": "0"}):
            runner = self._runner()
            with mock.patch.object(
                    owner_health, "_dingtalk_event_enqueue", return_value=True):
                queued = runner.reconcile({3: blocker})

        self.assertEqual(self.client.release_calls, [])
        self.assertEqual(queued, 1)

    def test_blocked_release_still_pages_and_does_not_hot_loop(self):
        # e.g. "tasks whose source item is terminal cannot be force released".
        self.client.release_action = "BLOCKED"
        blocker = self._blocker()
        queued, bodies = self._reconcile(blocker)
        self.assertEqual(queued, 1)
        self.assertEqual(len(bodies), 1)
        self.assertEqual(len(self.client.release_calls), 1)

        # Same window: no second attempt, and no second page.
        self.now += 60
        queued2, bodies2 = self._reconcile(blocker)
        self.assertEqual(len(self.client.release_calls), 1)
        self.assertEqual(queued2, 0)
        self.assertEqual(bodies2, [])

        # Next window: one more attempt is allowed.
        self.now += 3600
        self._reconcile(blocker)
        self.assertEqual(len(self.client.release_calls), 2)

    def test_session_already_cleared_pages_once_not_every_window(self):
        """A prior force-release cleared the Session; collect sees none.

        Such a blocker is not releasable (no live Session to discard) and the
        default alert path recommended a discard-resume that now has no
        Session to act on — so it must be flagged unactionable and paged once,
        not every repeat window.
        """
        blocker = self._blocker()
        blocker["session_id"] = ""  # collect found no live Session
        blocker["releasable"] = False
        queued, bodies = self._reconcile(blocker)

        self.assertEqual(self.client.release_calls, [])
        self.assertEqual(queued, 1, "page once with the real next step")
        self.assertIn("幻影 Session 已被本 runner 清除", bodies[0])
        self.assertIn("留一条人工评论", bodies[0])

        # Second window: unactionable, so no repeat.
        self.now += 7200
        queued2, bodies2 = self._reconcile(blocker)
        self.assertEqual(queued2, 0)
        self.assertEqual(bodies2, [])

    def test_no_pending_revision_uses_discard_resume_to_revive(self):
        """desired==processing: discard-resume moves the Task to READY."""
        self.client.desired_revision = "rev-a"  # == processing
        blocker = self._blocker()
        self.assertTrue(blocker["releasable"])
        queued, bodies = self._reconcile(blocker)

        self.assertEqual(queued, 0, "discarded blockers must not page")
        self.assertEqual(bodies, [])
        calls = self.client.release_calls
        self.assertEqual(len(calls), 1)
        self.assertEqual(calls[0]["via"], "discard_resume")
        self.assertEqual(calls[0]["task_id"], "3")
        self.assertEqual(calls[0]["expected_session_id"], 33)
        self.assertIn("no transcript/branch/checkpoint", calls[0]["reason"])
        self.assertEqual(self.runner._load_state()["episodes"], {})

    def test_no_session_and_no_pending_revision_is_unactionable(self):
        """Prior force-release cleared the Session; nothing can revive it."""
        self.client.desired_revision = "rev-a"  # == processing, no pending
        blocker = self._blocker()
        self.assertTrue(blocker["releasable"])  # collect saw a live session
        # But by the time reconcile runs, the fresh timeline shows the Session
        # is already gone (a previous force-release pass cleared it).
        original = self.client.get_task_timeline
        self.client.get_task_timeline = lambda tid: {
            "task": {**original("3")["task"], "currentSessionId": None},
            "sessions": [], "currentWorker": None, "events": []}
        queued, bodies = self._reconcile(blocker)

        self.assertEqual(self.client.release_calls, [])
        self.assertEqual(queued, 1, "page once with the real next step")
        self.assertIn("幻影 Session 已被本 runner 清除", bodies[0])
        self.assertIn("留一条人工评论", bodies[0])
        # Not repeat: same episode, second window.
        self.now += 7200
        queued2, bodies2 = self._reconcile(blocker)
        self.assertEqual(queued2, 0)
        self.assertEqual(bodies2, [])

    def test_stale_frozen_policy_is_never_released(self):
        """Reviving it would fail StaleTaskPolicyError and re-strand every window."""
        self.client.payload_policy = "terraform-rd-single-writer-v4"
        blocker = self._blocker()
        self.assertFalse(blocker["releasable"])
        self.assertTrue(blocker["stale_policy"])
        queued, bodies = self._reconcile(blocker)

        self.assertEqual(self.client.release_calls, [])
        self.assertEqual(queued, 1)
        # The page must name the real blocker, not send the operator after a
        # discard-resume that cannot help.
        self.assertIn("payload 冻结在旧策略版本", bodies[0])
        self.assertIn("stale_task_policy_revision", bodies[0])
        self.assertNotIn("discard-resume", bodies[0])

    def test_stale_frozen_policy_pages_once_not_every_window(self):
        self.client.payload_policy = "terraform-rd-single-writer-v4"
        blocker = self._blocker()
        _queued, bodies = self._reconcile(blocker)
        self.assertEqual(len(bodies), 1)

        self.now += 7200  # two repeat windows later
        queued2, bodies2 = self._reconcile(blocker)
        self.assertEqual(queued2, 0)
        self.assertEqual(bodies2, [])

    def test_absent_payload_is_treated_as_unverifiable_not_releasable(self):
        self.client.payload_policy = _ABSENT
        blocker = self._blocker()
        self.assertFalse(blocker["releasable"])
        # Unknown is not the same as stale: do not claim a cause we cannot see.
        self.assertFalse(blocker["stale_policy"])
        queued, bodies = self._reconcile(blocker)

        self.assertEqual(self.client.release_calls, [])
        self.assertEqual(queued, 1)
        self.assertIn("discard-resume", bodies[0])

    def test_context_appearing_before_the_write_aborts_the_release(self):
        """The emptiness check is re-done on the fresh timeline, not the snapshot."""
        blocker = self._blocker()
        self.assertTrue(blocker["releasable"])
        # Session externalized a transcript between collect and reconcile.
        self.client.session_refs = {
            "transcriptUri": "/tmp/late.jsonl",
            "branchRef": None, "checkpointUri": None}
        queued, bodies = self._reconcile(blocker)

        self.assertEqual(self.client.release_calls, [])
        self.assertEqual(queued, 1, "falls back to paging instead of guessing")
        self.assertEqual(len(bodies), 1)


if __name__ == "__main__":
    unittest.main()
