#!/usr/bin/env python3
"""Regression tests for bridge-owned post-PR Task bookends."""

import unittest
from unittest import mock
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))

import jarvis_dingtalk_bot as bot


class _OperationClient:
    def __init__(self, *, claim_proceed=True):
        self.claim_proceed = claim_proceed
        self.calls = []

    def begin_operation(self, request, request_id=None):
        action = "claim" if request["operationType"] == "AONE_CLAIM" else "release"
        self.calls.append(("begin_" + action, request, request_id))
        return {
            "operation": {"id": "operation-" + action, "status": "SENDING"},
            "proceed": self.claim_proceed if action == "claim" else True,
        }

    def ack_operation(self, request, request_id=None):
        self.calls.append(("ack", request, request_id))
        return {}


class _Controller:
    def __init__(self, client):
        self.client = client
        self.task = {"id": "task-1", "generation": 7}
        self.session = {"id": "session-1", "generation": 7}
        self.session_id = "session-1"
        self.worker_key = "host:boot:worker"
        self.fence_token = "fence-7"
        self.runtime_session_id = "runtime-session-1"
        self.resumed = False
        self.bound = []
        self.heartbeats = []

    def bind_process(self, process):
        self.bound.append(process)
        return process

    def heartbeat(self, detail=None):
        self.heartbeats.append(detail)
        return True


class PostPrTaskBookendTest(unittest.TestCase):
    def test_task_lease_wires_bookend_as_the_only_spawn_binder(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler.execution_router = mock.Mock(
            task_types={"pr_ci_fix", "pr_comment_reply"})
        handler._broadcast = mock.Mock()
        handler.dispatch_item = mock.Mock(return_value="done")
        controller = _Controller(_OperationClient())
        lease = {
            "task": {"id": "task-1", "generation": 7,
                     "taskType": "pr_ci_fix"},
            "session": {
                "id": "session-1", "generation": 7,
                "inputPayload": {
                    "itemId": "84362517",
                    "project": "2100304",
                    "kind": "pr_ci_fix",
                    "prompt": "fix CI",
                    "terraform": True,
                    "policyRevision": bot.HEADLESS_POLICY_REVISION,
                },
            },
        }
        bookend = mock.Mock()
        with mock.patch.object(bot, "_PostPrTaskBookend",
                               return_value=bookend) as factory:
            result = handler._execute_task_lease(lease, controller)

        self.assertEqual(result, "done")
        factory.assert_called_once_with(
            controller, "84362517", "2100304", "pr_ci_fix")
        kwargs = handler.dispatch_item.call_args.kwargs
        self.assertIs(kwargs["on_spawn"], bookend.bind_process)
        self.assertIs(kwargs["post_pr_bookend"], bookend)
        self.assertIs(kwargs["session_controller"], controller)

    def test_claim_binds_pid_then_acks_required_receipt_and_release(self):
        client = _OperationClient()
        controller = _Controller(client)
        bookend = bot._PostPrTaskBookend(
            controller, "84362517", "2100304", "pr_ci_fix")
        claimed = {"value": False}

        def claim(*_args, **_kwargs):
            claimed["value"] = True

        def release(*_args, **_kwargs):
            claimed["value"] = False

        process = mock.Mock(pid=4321)
        with mock.patch.object(bot, "_claim_workitem", side_effect=claim) as claim_call, \
                mock.patch.object(bot, "_release_post_pr_claim",
                                  side_effect=release) as release_call, \
                mock.patch.object(bot, "_post_pr_claim_visible",
                                  side_effect=lambda *_a, **_k: claimed["value"]):
            bookend.bind_process(process)
            self.assertEqual(controller.bound, [process])
            self.assertTrue(claimed["value"])
            bookend.release()

        claim_call.assert_called_once_with("84362517", "2100304", terraform=True)
        release_call.assert_called_once_with("84362517", "2100304", terraform=True)
        operation_types = [
            call[1]["operationType"] for call in client.calls
            if call[0].startswith("begin_")]
        self.assertEqual(operation_types, ["AONE_CLAIM", "AONE_RELEASE"])
        self.assertEqual(
            [call[1]["operationId"] for call in client.calls if call[0] == "ack"],
            ["operation-claim", "operation-release"])
        self.assertGreaterEqual(len(controller.heartbeats), 4)

    def test_sending_claim_recovery_point_reads_before_ack(self):
        client = _OperationClient(claim_proceed=False)
        controller = _Controller(client)
        bookend = bot._PostPrTaskBookend(
            controller, "84362517", "2100304", "pr_comment_reply")
        with mock.patch.object(bot, "_post_pr_claim_visible", return_value=True), \
                mock.patch.object(bot, "_claim_workitem") as claim_call:
            bookend.bind_process(mock.Mock(pid=4322))

        claim_call.assert_not_called()
        self.assertEqual(client.calls[-1][1]["operationId"], "operation-claim")

    def test_lost_fence_blocks_aone_write_and_guard_open(self):
        client = _OperationClient()
        controller = _Controller(client)
        controller.heartbeat = mock.Mock(return_value=False)
        bookend = bot._PostPrTaskBookend(
            controller, "84362517", "2100304", "pr_ci_fix")
        with mock.patch.object(bot, "_claim_workitem") as claim_call:
            with self.assertRaisesRegex(RuntimeError, "Task fence lost"):
                bookend.bind_process(mock.Mock(pid=4323))

        self.assertEqual(len(controller.bound), 1, "PID binding must precede Aone claim")
        claim_call.assert_not_called()

    def test_dispatch_release_failure_keeps_task_retryable(self):
        handler = bot.JarvisHandler.__new__(bot.JarvisHandler)
        handler._maybe_suspend = mock.Mock(return_value=None)
        handler._completion_broadcast = mock.Mock(return_value="must-not-broadcast")
        handler._dispatch_failed = mock.Mock()
        bookend = mock.Mock()
        bookend.lineage_policy.return_value = {
            "policyRevision": bot.HEADLESS_POLICY_REVISION,
            "aoneWritePolicy": bot.POST_PR_AONE_WRITE_POLICY,
            "kind": "pr_ci_fix",
            "aoneId": "84362517",
            "projectId": "2100304",
            "claimAttemptId": "attempt-1",
        }
        bookend.release.side_effect = RuntimeError("release receipt unavailable")
        notifications = []

        with mock.patch.object(
                bot, "run_claude_buffered",
                return_value=bot.ClaudeResult("ok", False, "success")) as run:
            with self.assertRaisesRegex(RuntimeError, "release receipt unavailable"):
                handler.dispatch_item(
                    "84362517", "prompt", "runtime-session", False,
                    notifications.append, "target", "staff",
                    on_spawn=bookend.bind_process, project="2100304",
                    kind="pr_ci_fix", terraform=True,
                    session_controller=mock.Mock(), post_pr_bookend=bookend)

        self.assertEqual(
            run.call_args.kwargs["aone_write_policy"],
            bot.POST_PR_AONE_WRITE_POLICY)
        self.assertTrue(run.call_args.kwargs["guarded"])
        handler._completion_broadcast.assert_not_called()
        self.assertEqual(notifications, [])


if __name__ == "__main__":
    unittest.main(verbosity=2)
