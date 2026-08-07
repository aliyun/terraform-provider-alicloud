#!/usr/bin/env python3
"""Contract tests for durable Terraform D-route DingTalk notifications."""

import json
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest import mock

from bridge.helpers import dingtalk
from bridge import terraform_route_notify


class TerraformRouteNotifyTest(unittest.TestCase):
    def setUp(self):
        self.tempdir = tempfile.TemporaryDirectory()
        self.ledger = Path(self.tempdir.name) / "dingtalk-ledger.json"
        self.path_patch = mock.patch.object(
            dingtalk, "DINGTALK_EVENT_PATH", self.ledger)
        self.path_patch.start()
        # The owner gate shells out to aone-assign.sh --check, which reads the
        # live work item. Delivery semantics below are owner-independent, so the
        # gate is stubbed open here and exercised on its own further down.
        self.owner_patch = mock.patch.object(
            terraform_route_notify, "_recipient_owns_or_may_own",
            return_value=(True, "stubbed open"))
        self.owner_patch.start()
        dingtalk._event_inflight.clear()

    def tearDown(self):
        dingtalk._event_inflight.clear()
        self.owner_patch.stop()
        self.path_patch.stop()
        self.tempdir.cleanup()

    @staticmethod
    def _sent(command, **_kwargs):
        receipt = command[command.index("--out-track-id") + 1]
        return subprocess.CompletedProcess(
            command, 0,
            stdout=json.dumps({"status": "sent", "receipt": receipt}) + "\n",
            stderr="")

    def test_typed_subtypes_map_to_owner_and_stable_event_key(self):
        expected = {
            "handwritten-urgent": "521957",
            "handwritten-normal": "484483",
            "generated": "429768",
        }
        with mock.patch.object(
                dingtalk.subprocess, "run", side_effect=self._sent) as sender:
            for subtype, owner in expected.items():
                result = terraform_route_notify.enqueue_d_route_notification(
                    "84872924", subtype)
                self.assertTrue(result["durable"])
                self.assertTrue(result["notification_complete"])
                self.assertEqual(result["state"], "posted")
                self.assertEqual(result["staff_id"], owner)
                self.assertEqual(
                    result["event_key"],
                    f"terraform-route:d:{subtype}:owner:{owner}")
            self.assertEqual(sender.call_count, 3)

    def test_same_ticket_subtype_and_owner_retries_only_once(self):
        with mock.patch.object(
                dingtalk.subprocess, "run", side_effect=self._sent) as sender:
            first = terraform_route_notify.enqueue_d_route_notification(
                "84872924", "handwritten-urgent")
            second = terraform_route_notify.enqueue_d_route_notification(
                "84872924", "handwritten-urgent")
        self.assertEqual(sender.call_count, 1)
        self.assertEqual(first["ledger_id"], second["ledger_id"])
        self.assertEqual(first["receipt"], second["receipt"])
        self.assertEqual(second["state"], "posted")

    def test_owner_or_subtype_change_creates_a_new_event(self):
        with mock.patch.object(
                dingtalk.subprocess, "run", side_effect=self._sent) as sender:
            urgent = terraform_route_notify.enqueue_d_route_notification(
                "84872924", "handwritten-urgent")
            normal = terraform_route_notify.enqueue_d_route_notification(
                "84872924", "handwritten-normal")
        self.assertEqual(sender.call_count, 2)
        self.assertNotEqual(urgent["ledger_id"], normal["ledger_id"])
        self.assertNotEqual(urgent["receipt"], normal["receipt"])

    def test_durable_pending_and_post_uncertain_keep_the_same_receipt(self):
        uncertain = subprocess.CompletedProcess([], 0, stdout="", stderr="")
        with mock.patch.object(
                dingtalk.subprocess, "run", return_value=uncertain) as sender:
            first = terraform_route_notify.enqueue_d_route_notification(
                "84872924", "generated")
            second = terraform_route_notify.enqueue_d_route_notification(
                "84872924", "generated")
        self.assertEqual(sender.call_count, 1)
        self.assertTrue(first["durable"])
        self.assertFalse(first["notification_complete"])
        self.assertEqual(first["state"], "post_uncertain")
        self.assertEqual(first["receipt"], second["receipt"])
        self.assertEqual(second["state"], "post_uncertain")

    def test_unpersisted_ledger_never_claims_notification_complete(self):
        with mock.patch.object(
                dingtalk, "_dingtalk_event_write", return_value=False), \
                mock.patch.object(dingtalk.subprocess, "run") as sender:
            result = terraform_route_notify.enqueue_d_route_notification(
                "84872924", "handwritten-normal")
        sender.assert_not_called()
        self.assertFalse(result["durable"])
        self.assertFalse(result["notification_complete"])
        self.assertEqual(result["state"], "unpersisted")

    def test_g_is_not_a_valid_d_notification_subtype(self):
        with self.assertRaises(ValueError):
            terraform_route_notify.enqueue_d_route_notification(
                "84872924", "g")


class RouteDmOwnerGateTest(unittest.TestCase):
    """The DM recipient comes from the subtype, so a protected owner must gate it.

    Regression for #84363256: an ACK ticket owned by 若即 (专属维护名单, 分支 A)
    was reclassified as 「D 手写非紧急」, reassigned to 过载, and 过载 was DM'd
    「已按手写 resource 非紧急路由给你」. Refusing only the assignee write would
    still leave the DM pointing at the wrong person.
    """

    def setUp(self):
        self.tempdir = tempfile.TemporaryDirectory()
        self.ledger = Path(self.tempdir.name) / "dingtalk-ledger.json"
        self.path_patch = mock.patch.object(
            dingtalk, "DINGTALK_EVENT_PATH", self.ledger)
        self.path_patch.start()
        dingtalk._event_inflight.clear()

    def tearDown(self):
        dingtalk._event_inflight.clear()
        self.path_patch.stop()
        self.tempdir.cleanup()

    def _enqueue_with_gate(self, allowed, detail="refused"):
        with mock.patch.object(
                terraform_route_notify, "_recipient_owns_or_may_own",
                return_value=(allowed, detail)) as gate, \
                mock.patch.object(dingtalk, "_dingtalk_event_enqueue") as sender:
            result = terraform_route_notify.enqueue_d_route_notification(
                "84363256", "handwritten-normal")
        return result, gate, sender

    def test_protected_owner_skips_the_dm_without_enqueueing(self):
        result, gate, sender = self._enqueue_with_gate(
            False, "已由云产品专属维护人 若即(377376) 持单")
        gate.assert_called_once_with("84363256", "484483")
        sender.assert_not_called()
        self.assertTrue(result["skipped"])
        self.assertEqual(result["state"], "skipped_owner_protected")
        self.assertFalse(result["notification_complete"])
        self.assertFalse(result["durable"])
        self.assertEqual(result["ledger_id"], "")
        self.assertIn("若即", result["skip_reason"])

    def test_a_skip_is_not_reported_as_a_delivery_failure(self):
        # rc=0 so the finalizer does not retry; the payload still says not-sent.
        with mock.patch.object(
                terraform_route_notify, "_recipient_owns_or_may_own",
                return_value=(False, "已由真人 辰羿(320687) 跟进中")), \
                mock.patch.object(dingtalk, "_dingtalk_event_enqueue"):
            rc = terraform_route_notify.main(
                ["--ticket", "84363256", "--subtype", "handwritten-normal"])
        self.assertEqual(rc, 0)

    def test_open_owner_still_enqueues(self):
        result, _gate, sender = self._enqueue_with_gate(True, "would be allowed")
        sender.assert_called_once()
        self.assertFalse(result["skipped"])

    def test_unverifiable_ownership_skips_rather_than_paging_someone(self):
        # aone-assign.sh --check is fail-closed, so a read failure lands here as
        # not-allowed. The assignee sync failed the same way; staying quiet keeps
        # both halves of the route step consistent.
        with mock.patch.object(
                terraform_route_notify.subprocess, "run",
                side_effect=OSError("no bash")), \
                mock.patch.object(dingtalk, "_dingtalk_event_enqueue") as sender:
            result = terraform_route_notify.enqueue_d_route_notification(
                "84363256", "generated")
        sender.assert_not_called()
        self.assertTrue(result["skipped"])
        self.assertIn("owner check failed to run", result["skip_reason"])

    def test_gate_reads_the_policy_from_aone_assign_check(self):
        with mock.patch.object(
                terraform_route_notify.subprocess, "run") as runner:
            runner.return_value = subprocess.CompletedProcess(
                [], 3, stdout="", stderr="aone-assign.sh: refusing to reassign\n")
            allowed, detail = terraform_route_notify._recipient_owns_or_may_own(
                "84363256", "484483")
        argv = runner.call_args[0][0]
        self.assertEqual(argv[0], "bash")
        self.assertTrue(argv[1].endswith("bootstrap/aone-assign.sh"))
        self.assertEqual(argv[2:], ["--check", "84363256", "484483"])
        self.assertFalse(allowed)
        self.assertIn("refusing to reassign", detail)


if __name__ == "__main__":
    unittest.main()
