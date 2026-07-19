#!/usr/bin/env python3
"""Production bridge recovery tests for uncertain post-PR Aone receipts."""

import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock
import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))

import jarvis_dingtalk_bot as bot


def candidate(action="claim", operation_status="UNKNOWN"):
    operation_type = "AONE_CLAIM" if action == "claim" else "AONE_RELEASE"
    return {
        "recoveryAction": action.upper(),
        "task": {
            "id": "task-1",
            "aoneId": "84362517",
            "taskType": "pr_ci_fix",
            "status": "RECOVERY_REQUIRED",
            "generation": 7,
            "payload": {
                "itemId": "84362517",
                "project": "2100304",
                "terraform": True,
            },
        },
        "operation": {
            "id": "operation-1",
            "generation": 7,
            "operationKey": "post-pr:%s:attempt-1" % action,
            "operationType": operation_type,
            "target": "84362517",
            "required": True,
            "status": operation_status,
        },
    }


class PostPrOperationRecoveryTest(unittest.TestCase):
    def scheduler(self):
        client = mock.Mock()
        handler = SimpleNamespace(
            task_client=client,
            persistence_executor=SimpleNamespace(worker_key="host:boot:process"),
            ephemeral_executor=SimpleNamespace(active_ids=lambda: []),
        )
        return bot.PostPrRecoverySensor(handler), client

    def test_exact_tag_state_requires_one_side_and_rejects_terminal_tags(self):
        states = [
            {"jarvis-claimed"},
            {"jarvis-idle"},
            {"jarvis-claimed", "jarvis-idle"},
            set(),
        ]
        for tags in states:
            for action in ("claim", "release"):
                expected = ((action == "claim" and tags == {"jarvis-claimed"})
                            or (action == "release" and tags == {"jarvis-idle"}))
                with self.subTest(tags=tags, action=action), mock.patch.object(
                        bot, "_post_pr_tag_snapshot",
                        return_value={"tags": tags, "pairs": []}):
                    self.assertEqual(
                        bot._post_pr_target_visible("84362517", action), expected)
        for terminal in ("jarvis-done", "jarvis-npe"):
            with mock.patch.object(
                    bot, "_post_pr_tag_snapshot",
                    return_value={"tags": {terminal}, "pairs": []}):
                with self.assertRaisesRegex(RuntimeError, "terminal"):
                    bot._post_pr_target_visible("84362517", "claim")

    def test_candidate_repairs_then_reconciles_found_true_with_stable_token(self):
        scheduler, client = self.scheduler()
        client.lease_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "CLAIM"}
        client.renew_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "CLAIM"}

        def repair(_item_id, _action, terraform=True,
                   before_external_step=None):
            self.assertTrue(terraform)
            before_external_step("repair-read")
            before_external_step("repair-write")

        with mock.patch.object(
                bot, "_post_pr_target_visible", side_effect=[False, True]), \
                mock.patch.object(bot, "_repair_post_pr_tags",
                                  side_effect=repair) as repair_mock:
            self.assertTrue(scheduler._recover_post_pr_candidate(candidate("claim")))

        repair_mock.assert_called_once()
        lease_request = client.lease_operation_recovery.call_args.args[0]
        reconcile_request = client.reconcile_operation.call_args.args[0]
        self.assertEqual(lease_request["recoveryToken"],
                         reconcile_request["recoveryToken"])
        self.assertTrue(reconcile_request["found"])
        self.assertFalse(reconcile_request["retryAllowed"])
        self.assertEqual(client.renew_operation_recovery.call_count, 4)
        for renew_call in client.renew_operation_recovery.call_args_list:
            self.assertEqual(renew_call.args[0], lease_request)
        client.release_operation_recovery.assert_not_called()

    def test_scheduler_renews_before_initial_and_final_read(self):
        scheduler, client = self.scheduler()
        client.lease_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "CLAIM"}
        events = []

        def renew(_request, request_id=None):
            phase = "initial" if "initial-read" in request_id else "final"
            events.append("renew-" + phase)
            return {"proceed": True, "recoveryAction": "CLAIM"}

        def readback(*_args, **_kwargs):
            events.append("read")
            return True

        client.renew_operation_recovery.side_effect = renew
        with mock.patch.object(
                bot, "_post_pr_target_visible", side_effect=readback):
            self.assertTrue(scheduler._recover_post_pr_candidate(candidate("claim")))
        self.assertEqual(events, [
            "renew-initial", "read", "renew-final", "read"])

    def test_read_or_reconcile_failure_releases_lease_without_guessing_ack(self):
        scheduler, client = self.scheduler()
        client.lease_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "RELEASE"}
        client.renew_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "RELEASE"}
        with mock.patch.object(
                bot, "_post_pr_target_visible",
                side_effect=RuntimeError("readback unavailable")), \
                mock.patch.object(bot, "_repair_post_pr_tags") as repair:
            with self.assertRaisesRegex(RuntimeError, "readback unavailable"):
                scheduler._recover_post_pr_candidate(candidate("release"))
        repair.assert_not_called()
        client.reconcile_operation.assert_not_called()
        client.release_operation_recovery.assert_called_once()

        client.reset_mock()
        client.lease_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "RELEASE"}
        client.renew_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "RELEASE"}
        client.reconcile_operation.side_effect = RuntimeError("response lost")
        with mock.patch.object(bot, "_post_pr_target_visible", return_value=True), \
                mock.patch.object(bot, "_repair_post_pr_tags") as repair:
            with self.assertRaisesRegex(RuntimeError, "response lost"):
                scheduler._recover_post_pr_candidate(candidate("release"))
        repair.assert_not_called()
        client.release_operation_recovery.assert_called_once()

    def test_competing_or_later_blocked_lease_never_reads_or_writes_aone(self):
        scheduler, client = self.scheduler()
        client.lease_operation_recovery.return_value = {"proceed": False}
        with mock.patch.object(bot, "_post_pr_target_visible") as readback, \
                mock.patch.object(bot, "_repair_post_pr_tags") as repair:
            self.assertFalse(scheduler._recover_post_pr_candidate(candidate()))
        readback.assert_not_called()
        repair.assert_not_called()
        client.reconcile_operation.assert_not_called()

    def test_due_retry_wait_candidate_is_allowed_to_recover(self):
        scheduler, client = self.scheduler()
        client.lease_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "RELEASE"}
        client.renew_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "RELEASE"}
        with mock.patch.object(
                bot, "_post_pr_target_visible", return_value=True):
            self.assertTrue(scheduler._recover_post_pr_candidate(
                candidate("release", operation_status="RETRY_WAIT")))
        client.lease_operation_recovery.assert_called_once()
        self.assertEqual(client.renew_operation_recovery.call_count, 2)
        client.reconcile_operation.assert_called_once()

    def test_candidate_and_lease_canonical_action_mismatch_releases_without_aone(self):
        scheduler, client = self.scheduler()
        client.lease_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "RELEASE"}
        with mock.patch.object(bot, "_post_pr_target_visible") as readback, \
                mock.patch.object(bot, "_repair_post_pr_tags") as repair:
            with self.assertRaisesRegex(RuntimeError, "canonical actions differ"):
                scheduler._recover_post_pr_candidate(candidate("claim"))
        readback.assert_not_called()
        repair.assert_not_called()
        client.renew_operation_recovery.assert_not_called()
        client.reconcile_operation.assert_not_called()
        client.release_operation_recovery.assert_called_once()

    def test_renew_failure_before_repair_write_aborts_write_and_releases(self):
        scheduler, client = self.scheduler()
        client.lease_operation_recovery.return_value = {
            "proceed": True, "recoveryAction": "CLAIM"}
        renewed = {"proceed": True, "recoveryAction": "CLAIM"}
        client.renew_operation_recovery.side_effect = [
            renewed, renewed, RuntimeError("renew unavailable")]
        snapshot = {"tags": {"other"}, "pairs": [("other", "77")]}
        with mock.patch.object(
                bot, "_post_pr_tag_snapshot", return_value=snapshot), \
                mock.patch.object(bot.subprocess, "run") as run:
            with self.assertRaisesRegex(RuntimeError, "renew unavailable"):
                scheduler._recover_post_pr_candidate(candidate("claim"))
        self.assertEqual(client.renew_operation_recovery.call_count, 3)
        run.assert_not_called()
        client.reconcile_operation.assert_not_called()
        client.release_operation_recovery.assert_called_once()

    def test_untrusted_candidate_is_rejected_before_lease_or_aone_access(self):
        scheduler, client = self.scheduler()
        bad = candidate()
        bad["task"]["payload"]["terraform"] = False
        with mock.patch.object(bot, "_post_pr_target_visible") as readback:
            with self.assertRaisesRegex(RuntimeError, "defense-in-depth"):
                scheduler._recover_post_pr_candidate(bad)
        client.lease_operation_recovery.assert_not_called()
        readback.assert_not_called()

    def test_terminal_tag_repair_is_fail_closed(self):
        with mock.patch.object(
                bot, "_post_pr_tag_snapshot",
                return_value={
                    "tags": {"jarvis-done", "other"},
                    "pairs": [("jarvis-done", "11"), ("other", "12")],
                }), mock.patch.object(bot.subprocess, "run") as run:
            with self.assertRaisesRegex(RuntimeError, "terminal"):
                bot._repair_post_pr_tags("84362517", "claim")
        run.assert_not_called()

    def test_tag_repair_renews_before_read_and_write(self):
        events = []
        snapshot = {"tags": {"other"}, "pairs": [("other", "77")]}

        def heartbeat(phase):
            events.append(phase)

        def readback(*_args, **_kwargs):
            events.append("read")
            return snapshot

        def write(*_args, **_kwargs):
            events.append("write")
            return SimpleNamespace(returncode=0, stdout="", stderr="")

        with mock.patch.object(
                bot, "_post_pr_tag_snapshot", side_effect=readback), \
                mock.patch.object(bot.subprocess, "run", side_effect=write):
            bot._repair_post_pr_tags(
                "84362517", "claim", before_external_step=heartbeat)
        self.assertEqual(events, [
            "repair-read", "read", "repair-write", "write"])

    def test_tag_repair_preserves_unrelated_id_and_makes_claim_idle_mutually_exclusive(self):
        snapshot = {
            "tags": {"other", "jarvis-claimed", "jarvis-idle"},
            "pairs": [
                ("other", "77"),
                ("jarvis-claimed", "88"),
                ("jarvis-idle", "99"),
            ],
        }
        for action, desired in (("claim", "jarvis-claimed"),
                                ("release", "jarvis-idle")):
            with self.subTest(action=action), mock.patch.object(
                    bot, "_post_pr_tag_snapshot", return_value=snapshot), \
                    mock.patch.object(
                        bot.subprocess, "run",
                        return_value=SimpleNamespace(
                            returncode=0, stdout="", stderr="")) as run:
                bot._repair_post_pr_tags("84362517", action)
            command = run.call_args.args[0]
            tag_value = command[command.index("--tag") + 1]
            self.assertEqual(tag_value, "77,%s" % desired)
            self.assertNotIn("88", tag_value)
            self.assertNotIn("99", tag_value)

    def test_tag_snapshot_rejects_unaligned_empty_names_and_nonempty_ids(self):
        payload = '{"fields":[{"identifier":"tag","displayValue":"","value":"77"}]}'
        with mock.patch.object(
                bot.subprocess, "run",
                return_value=SimpleNamespace(
                    returncode=0, stdout=payload, stderr="")):
            with self.assertRaisesRegex(RuntimeError, "invalid JSON"):
                bot._post_pr_tag_snapshot("84362517")


if __name__ == "__main__":
    unittest.main(verbosity=2)
