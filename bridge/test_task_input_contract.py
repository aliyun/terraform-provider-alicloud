#!/usr/bin/env python3
"""Portable input and legacy redispatch reconstruction tests."""

import unittest

from bridge.task_input_contract import (
    PORTABLE_INPUT_CONTRACT,
    TaskInputContractError,
    canonical_input_digest,
    portable_replacement_for_redispatch,
    portable_task_payload,
)


class TaskInputContractTest(unittest.TestCase):
    def test_payload_is_symmetric_across_runtime_modes(self):
        for runtime in ("INTERACTIVE", "PERSISTENT"):
            with self.subTest(runtime=runtime):
                payload = portable_task_payload(
                    item_id="84407231",
                    project="1086837",
                    kind="ticket",
                    prompt="continue",
                    origin_runtime=runtime,
                    trigger="TEST",
                )
                self.assertEqual(
                    payload["inputContract"], PORTABLE_INPUT_CONTRACT)
                self.assertEqual(
                    payload["inputContext"]["originRuntime"], runtime)

    def test_legacy_canonical_aone_input_is_rehydrated(self):
        source = {
            "kind": "ticket",
            "itemId": "84407231",
            "project": "1086837",
            "trigger": "INTERACTIVE",
            "policyRevision": "terraform-rd-single-writer-v6",
        }
        timeline = {
            "task": {
                "id": 603,
                "taskKey": "aone:1086837:84407231",
                "aoneId": "84407231",
                "sourceType": "AONE",
                "taskType": "ticket",
                "desiredRevision": "interactive:old",
                "processingRevision": "interactive:old",
                "currentSessionId": 1195,
                "sourceRef": {
                    "aoneId": "84407231",
                    "projectId": "1086837",
                    "title": "legacy interactive task",
                },
            },
            "sessions": [{
                "id": 1195,
                "inputPayload": source,
            }],
            "effectiveInputPayload": source,
            "effectiveInputDigest": canonical_input_digest(source),
        }

        replacement = portable_replacement_for_redispatch(
            timeline,
            target_runtime="PERSISTENT",
            selected_session_id=1195,
        )

        self.assertIsNotNone(replacement)
        payload = replacement["payload"]
        self.assertEqual(payload["itemId"], "84407231")
        self.assertEqual(payload["project"], "1086837")
        self.assertEqual(payload["inputContract"], "PORTABLE_V1")
        self.assertTrue(payload["inputContext"]["rehydrated"])
        self.assertIn("读取当前 Aone 详情和最新评论", payload["prompt"])
        self.assertIn("禁止重复评论", payload["prompt"])
        self.assertNotEqual(
            replacement["desiredRevision"], "interactive:old")
        self.assertEqual(
            replacement["replacementDigest"],
            canonical_input_digest(payload))
        self.assertEqual(
            replacement["desiredRevision"],
            "portable:%s" % replacement["replacementDigest"])

    def test_non_aone_legacy_input_is_typed_blocked(self):
        source = {"itemId": "local-1", "kind": "adhoc"}
        timeline = {
            "task": {
                "id": 1,
                "taskKey": "local:1",
                "sourceType": "LOCAL",
                "taskType": "adhoc",
                "currentSessionId": 2,
            },
            "effectiveInputPayload": source,
            "effectiveInputDigest": canonical_input_digest(source),
        }
        with self.assertRaises(TaskInputContractError) as raised:
            portable_replacement_for_redispatch(
                timeline, target_runtime="INTERACTIVE",
                selected_session_id=2)
        self.assertEqual(raised.exception.code, "INPUT_NOT_REHYDRATABLE")


if __name__ == "__main__":
    unittest.main()
