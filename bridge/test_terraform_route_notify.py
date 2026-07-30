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
        dingtalk._event_inflight.clear()

    def tearDown(self):
        dingtalk._event_inflight.clear()
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


if __name__ == "__main__":
    unittest.main()
