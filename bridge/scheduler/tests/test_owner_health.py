from __future__ import annotations

import logging
import tempfile
import unittest
from pathlib import Path
from unittest import mock

from bridge.scheduler.runners import owner_health


class FakeClient:
    def __init__(self, *, now: float) -> None:
        self.now = now
        self.ready_cursors: list[int] = []
        self.source_cursors: list[int] = []

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
        self.source_cursors.append(after_task_id)
        if after_task_id == 0:
            return {
                "items": [{
                    "taskId": 3,
                    "aoneId": "84550003",
                    "status": "RECOVERY_REQUIRED",
                    "recoveryPolicy": "RESUME_ONLY",
                    "currentSessionId": 33,
                    "updatedAt": self.now - 7200,
                }],
                "nextAfterTaskId": 3,
            }
        return {"items": []}

    def get_task_timeline(self, task_id: str):
        assert task_id == "3"
        return {
            "sessions": [{
                "id": 33,
                "status": "RESUMABLE",
                "workerKey": "worker-old-3",
            }],
            "currentWorker": None,
            "events": [{
                "eventType": "TASK_OWNER_MIGRATION",
                "reasonCode": "RESUME_OWNER_UNAVAILABLE",
                "occurredAt": self.now - 7200,
            }],
        }


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


if __name__ == "__main__":
    unittest.main()
